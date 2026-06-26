package store

import (
	"context"
	"math"
	"time"
)

// LearnerStats is a snapshot of one learner's activity for their dashboard.
type LearnerStats struct {
	TotalSubmissions int
	Pending          int // status = submitted (awaiting grading)
	Graded           int
	Returned         int // sent back for rework
	AvgScore         int // rounded average of graded scores
	BestScore        int
	TotalPoints      int // sum of graded scores
	ActiveDays       int // distinct days the learner submitted
	Competitions     int // competitions entered
	CoursesAvailable int // published courses
	LessonsAvailable int
}

// LearnerStats aggregates the headline numbers for a learner's dashboard.
func (s *Store) LearnerStats(ctx context.Context, userID int64) (LearnerStats, error) {
	var st LearnerStats

	rows, err := s.db.QueryContext(ctx,
		`SELECT status, COUNT(*) FROM submissions WHERE user_id = ? GROUP BY status`, userID)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			rows.Close()
			return st, err
		}
		switch status {
		case "submitted":
			st.Pending = n
		case "graded":
			st.Graded = n
		case "returned":
			st.Returned = n
		}
		st.TotalSubmissions += n
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return st, err
	}

	var avg float64
	_ = s.db.QueryRowContext(ctx,
		`SELECT COALESCE(AVG(score),0), COALESCE(MAX(score),0), COALESCE(SUM(score),0)
		 FROM submissions WHERE user_id = ? AND status = 'graded' AND score IS NOT NULL`, userID).
		Scan(&avg, &st.BestScore, &st.TotalPoints)
	st.AvgScore = int(math.Round(avg))

	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT date(created_at)) FROM submissions WHERE user_id = ?`, userID).Scan(&st.ActiveDays)
	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM competition_entries WHERE user_id = ?`, userID).Scan(&st.Competitions)
	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM courses WHERE published = 1`).Scan(&st.CoursesAvailable)
	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM lessons l JOIN courses c ON c.id = l.course_id WHERE c.published = 1`).
		Scan(&st.LessonsAvailable)

	return st, nil
}

// LearnerSubmissionSeries returns daily submission counts for the last `days`
// days (inclusive of today), with empty days filled in.
func (s *Store) LearnerSubmissionSeries(ctx context.Context, userID int64, days int) ([]DayCount, error) {
	if days < 1 {
		days = 14
	}
	counts := map[string]int{}
	rows, err := s.db.QueryContext(ctx,
		`SELECT date(created_at) d, COUNT(*) c FROM submissions
		 WHERE user_id = ? AND created_at >= date('now', ?) GROUP BY d`,
		userID, signedDays(-(days - 1)))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d string
		var c int
		if err := rows.Scan(&d, &c); err != nil {
			rows.Close()
			return nil, err
		}
		counts[d] = c
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]DayCount, 0, days)
	today := time.Now().UTC()
	for i := days - 1; i >= 0; i-- {
		day := today.AddDate(0, 0, -i)
		key := day.Format("2006-01-02")
		out = append(out, DayCount{Date: key, Label: day.Format("Jan 2"), Count: counts[key]})
	}
	return out, nil
}

// RecentSubmissionsByUser lists a learner's most recent submissions (with the
// assignment title) for the dashboard activity table.
func (s *Store) RecentSubmissionsByUser(ctx context.Context, userID int64, limit int) ([]Submission, error) {
	if limit < 1 {
		limit = 6
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT sub.id, sub.status, sub.score, sub.created_at, a.title, a.slug, a.language, a.max_points
		FROM submissions sub
		JOIN assignments a ON a.id = sub.assignment_id
		WHERE sub.user_id = ?
		ORDER BY sub.created_at DESC
		LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Submission
	for rows.Next() {
		var sub Submission
		if err := rows.Scan(&sub.ID, &sub.Status, &sub.Score, &sub.CreatedAt,
			&sub.Assignment.Title, &sub.Assignment.Slug, &sub.Assignment.Language, &sub.Assignment.MaxPoints); err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

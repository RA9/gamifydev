package store

import (
	"context"
	"strconv"
	"time"
)

// AdminStats is a snapshot of platform-wide numbers for the admin dashboard.
type AdminStats struct {
	TotalUsers       int
	Learners         int
	Graders          int
	Admins           int
	NewUsers7d       int
	NewUsers30d      int
	ActiveSessions   int
	TotalCourses     int
	PublishedCourses int
	DraftCourses     int
	TotalLessons     int
}

// DayCount is one bucket in a per-day time series.
type DayCount struct {
	Date  string // YYYY-MM-DD
	Label string // short display label, e.g. "Jun 24"
	Count int
}

// AdminStats aggregates the headline numbers shown on the admin overview.
func (s *Store) AdminStats(ctx context.Context) (AdminStats, error) {
	var st AdminStats

	rows, err := s.db.QueryContext(ctx, `SELECT role, COUNT(*) FROM users GROUP BY role`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var role string
		var n int
		if err := rows.Scan(&role, &n); err != nil {
			rows.Close()
			return st, err
		}
		switch role {
		case "learner":
			st.Learners = n
		case "grader":
			st.Graders = n
		case "admin":
			st.Admins = n
		}
		st.TotalUsers += n
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return st, err
	}

	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE created_at >= datetime('now','-7 days')`).Scan(&st.NewUsers7d)
	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE created_at >= datetime('now','-30 days')`).Scan(&st.NewUsers30d)

	now := time.Now().UTC().Format(time.RFC3339)
	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT user_id) FROM sessions WHERE expires_at > ?`, now).Scan(&st.ActiveSessions)

	_ = s.db.QueryRowContext(ctx,
		`SELECT COUNT(*), COALESCE(SUM(published),0) FROM courses`).Scan(&st.TotalCourses, &st.PublishedCourses)
	st.DraftCourses = st.TotalCourses - st.PublishedCourses

	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM lessons`).Scan(&st.TotalLessons)

	return st, nil
}

// SignupSeries returns daily signup counts for the last `days` days (inclusive
// of today), with empty days filled in so the series is contiguous.
func (s *Store) SignupSeries(ctx context.Context, days int) ([]DayCount, error) {
	if days < 1 {
		days = 14
	}
	counts := map[string]int{}
	rows, err := s.db.QueryContext(ctx,
		`SELECT date(created_at) d, COUNT(*) c FROM users
		 WHERE created_at >= date('now', ?) GROUP BY d`,
		signedDays(-(days - 1)))
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
		out = append(out, DayCount{
			Date:  key,
			Label: day.Format("Jan 2"),
			Count: counts[key],
		})
	}
	return out, nil
}

// RecentUsers returns the most recently registered users.
func (s *Store) RecentUsers(ctx context.Context, limit int) ([]User, error) {
	if limit < 1 {
		limit = 8
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+userCols+` FROM users ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.Bio, &u.AvatarURL, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// signedDays renders an integer day offset as a SQLite modifier like "-13 days".
func signedDays(n int) string {
	return strconv.Itoa(n) + " days"
}

// PublicStats are the few figures the landing page states about the
// curriculum.
//
// Queried rather than written into the template, because the previous landing
// page hardcoded "7 learning tracks" against a database holding four, and a
// number nobody can update is a number that will be wrong.
type PublicStats struct {
	Paths   int
	Courses int
	Lessons int
}

// PublicCurriculumStats counts only published material — what a visitor could
// actually go and read.
func (s *Store) PublicCurriculumStats(ctx context.Context) (PublicStats, error) {
	var st PublicStats
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM paths WHERE published = 1`).Scan(&st.Paths); err != nil {
		return st, err
	}
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM courses WHERE published = 1`).Scan(&st.Courses); err != nil {
		return st, err
	}
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lessons l
		JOIN courses c ON c.id = l.course_id
		WHERE c.published = 1`).Scan(&st.Lessons); err != nil {
		return st, err
	}
	return st, nil
}

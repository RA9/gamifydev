package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Assignment struct {
	ID         int64
	CourseID   sql.NullInt64
	Slug       string
	Title      string
	Language   string
	Prompt     string
	Starter    string
	MaxPoints  int
	Published  bool
	Sort       int
	CourseName string // joined for display
}

type Submission struct {
	ID         int64
	Assignment Assignment
	UserID     int64
	UserName   string
	UserEmail  string
	Code       string
	Note       string
	Status     string
	Score      sql.NullInt64
	Feedback   string
	GraderName string
	GradedAt   string
	CreatedAt  string
}

func (s Submission) IsGraded() bool { return s.Status == "graded" }

// UpsertAssignment inserts/updates an assignment by slug (used by the seeder).
func (s *Store) UpsertAssignment(ctx context.Context, a Assignment) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO assignments (course_id, slug, title, language, prompt, starter, max_points, published, sort, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(slug) DO UPDATE SET
			course_id=excluded.course_id, title=excluded.title, language=excluded.language,
			prompt=excluded.prompt, starter=excluded.starter, max_points=excluded.max_points,
			published=excluded.published, sort=excluded.sort, updated_at=datetime('now')`,
		a.CourseID, a.Slug, a.Title, a.Language, a.Prompt, a.Starter, a.MaxPoints, boolToInt(a.Published), a.Sort)
	return err
}

const assignmentCols = `a.id, a.course_id, a.slug, a.title, a.language, a.prompt, a.starter, a.max_points, a.published, a.sort, COALESCE(c.title,'')`

func scanAssignment(sc interface{ Scan(...any) error }) (*Assignment, error) {
	var a Assignment
	var pub int
	if err := sc.Scan(&a.ID, &a.CourseID, &a.Slug, &a.Title, &a.Language, &a.Prompt, &a.Starter, &a.MaxPoints, &pub, &a.Sort, &a.CourseName); err != nil {
		return nil, err
	}
	a.Published = pub == 1
	return &a, nil
}

func (s *Store) ListAssignments(ctx context.Context, includeUnpublished bool) ([]Assignment, error) {
	q := `SELECT ` + assignmentCols + ` FROM assignments a LEFT JOIN courses c ON c.id = a.course_id`
	if !includeUnpublished {
		q += ` WHERE a.published = 1`
	}
	q += ` ORDER BY a.sort, a.title`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Assignment
	for rows.Next() {
		a, err := scanAssignment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

// ListAssignmentsByCourse returns a course's published assignments in order —
// the "checkpoint" items shown within the course.
func (s *Store) ListAssignmentsByCourse(ctx context.Context, courseID int64) ([]Assignment, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+assignmentCols+` FROM assignments a LEFT JOIN courses c ON c.id = a.course_id
		 WHERE a.course_id = ? AND a.published = 1 ORDER BY a.sort, a.title`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Assignment
	for rows.Next() {
		a, err := scanAssignment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

func (s *Store) GetAssignmentBySlug(ctx context.Context, slug string) (*Assignment, error) {
	a, err := scanAssignment(s.db.QueryRowContext(ctx,
		`SELECT `+assignmentCols+` FROM assignments a LEFT JOIN courses c ON c.id = a.course_id WHERE a.slug = ?`, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

// --- Submissions ------------------------------------------------------------

func (s *Store) CreateSubmission(ctx context.Context, assignmentID, userID int64, code, note string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO submissions (assignment_id, user_id, code, note, status) VALUES (?, ?, ?, ?, 'submitted') RETURNING id`,
		assignmentID, userID, code, note).Scan(&id)
	return id, err
}

// LatestSubmission returns a learner's most recent submission for an assignment.
func (s *Store) LatestSubmission(ctx context.Context, assignmentID, userID int64) (*Submission, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, code, note, status, score, feedback, graded_at, created_at
		FROM submissions WHERE assignment_id = ? AND user_id = ? ORDER BY created_at DESC LIMIT 1`,
		assignmentID, userID)
	var sub Submission
	err := row.Scan(&sub.ID, &sub.Code, &sub.Note, &sub.Status, &sub.Score, &sub.Feedback, &sub.GradedAt, &sub.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &sub, err
}

// CountPendingSubmissions counts submissions awaiting grading (for badges).
func (s *Store) CountPendingSubmissions(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM submissions WHERE status != 'graded'`).Scan(&n)
	return n, err
}

// GradingQueue lists submissions for the grader view (pending first, newest first).
func (s *Store) GradingQueue(ctx context.Context, onlyPending bool) ([]Submission, error) {
	q := `
		SELECT sub.id, sub.status, sub.score, sub.created_at,
		       u.name, u.email, a.title, a.language
		FROM submissions sub
		JOIN users u ON u.id = sub.user_id
		JOIN assignments a ON a.id = sub.assignment_id`
	if onlyPending {
		q += ` WHERE sub.status != 'graded'`
	}
	q += ` ORDER BY (sub.status = 'graded'), sub.created_at DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Submission
	for rows.Next() {
		var sub Submission
		if err := rows.Scan(&sub.ID, &sub.Status, &sub.Score, &sub.CreatedAt,
			&sub.UserName, &sub.UserEmail, &sub.Assignment.Title, &sub.Assignment.Language); err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

// GetSubmission returns a submission with its assignment + learner + grader.
func (s *Store) GetSubmission(ctx context.Context, id int64) (*Submission, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT sub.id, sub.user_id, sub.code, sub.note, sub.status, sub.score, sub.feedback,
		       sub.graded_at, sub.created_at, u.name, u.email,
		       a.id, a.slug, a.title, a.language, a.prompt, a.max_points,
		       COALESCE(g.name,'')
		FROM submissions sub
		JOIN users u ON u.id = sub.user_id
		JOIN assignments a ON a.id = sub.assignment_id
		LEFT JOIN users g ON g.id = sub.graded_by
		WHERE sub.id = ?`, id)
	var sub Submission
	err := row.Scan(&sub.ID, &sub.UserID, &sub.Code, &sub.Note, &sub.Status, &sub.Score, &sub.Feedback,
		&sub.GradedAt, &sub.CreatedAt, &sub.UserName, &sub.UserEmail,
		&sub.Assignment.ID, &sub.Assignment.Slug, &sub.Assignment.Title, &sub.Assignment.Language,
		&sub.Assignment.Prompt, &sub.Assignment.MaxPoints, &sub.GraderName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &sub, err
}

// GradeSubmission records a grade + feedback and marks the submission graded
// (or returned for rework).
func (s *Store) GradeSubmission(ctx context.Context, id, graderID int64, score int, feedback, status string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE submissions SET score = ?, feedback = ?, status = ?, graded_by = ?, graded_at = ?
		WHERE id = ?`,
		score, feedback, status, graderID, time.Now().UTC().Format(time.RFC3339), id)
	return err
}

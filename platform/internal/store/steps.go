package store

import (
	"context"
	"database/sql"
	"errors"
)

// Step is one interactive step of a step-based lesson.
type Step struct {
	ID          int64
	LessonID    int64
	Sort        int
	Instruction string
	Starter     string
	Checks      string // raw JSON array of {text,test}
	Lang        string // html | js
	Scaffold    string // optional base HTML rendered before a js step's code
}

// langOr defaults an empty/unknown language to html.
func langOr(l string) string {
	switch l {
	case "js", "python":
		return l
	default:
		return "html"
	}
}

// ListSteps returns a lesson's steps in order.
func (s *Store) ListSteps(ctx context.Context, lessonID int64) ([]Step, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, lesson_id, sort, instruction, starter, checks, lang, scaffold FROM lesson_steps
		 WHERE lesson_id = ? ORDER BY sort`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Step
	for rows.Next() {
		var st Step
		if err := rows.Scan(&st.ID, &st.LessonID, &st.Sort, &st.Instruction, &st.Starter, &st.Checks, &st.Lang, &st.Scaffold); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

// GetStep returns a single step by id.
func (s *Store) GetStep(ctx context.Context, id int64) (*Step, error) {
	var st Step
	err := s.db.QueryRowContext(ctx,
		`SELECT id, lesson_id, sort, instruction, starter, checks, lang, scaffold FROM lesson_steps WHERE id = ?`, id).
		Scan(&st.ID, &st.LessonID, &st.Sort, &st.Instruction, &st.Starter, &st.Checks, &st.Lang, &st.Scaffold)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &st, err
}

// ReplaceLessonSteps swaps a lesson's steps for the given set (used by the
// seeder; idempotent).
func (s *Store) ReplaceLessonSteps(ctx context.Context, lessonID int64, steps []Step) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `DELETE FROM lesson_steps WHERE lesson_id = ?`, lessonID); err != nil {
		return err
	}
	for i, st := range steps {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO lesson_steps (lesson_id, sort, instruction, starter, checks, lang, scaffold) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			lessonID, i, st.Instruction, st.Starter, st.Checks, langOr(st.Lang), st.Scaffold); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// CreateStep appends a step to a lesson (sort = next index) and returns its id.
func (s *Store) CreateStep(ctx context.Context, st Step) (int64, error) {
	var sort int
	_ = s.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(sort)+1, 0) FROM lesson_steps WHERE lesson_id = ?`, st.LessonID).Scan(&sort)
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO lesson_steps (lesson_id, sort, instruction, starter, checks, lang, scaffold)
		 VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		st.LessonID, sort, st.Instruction, st.Starter, st.Checks, langOr(st.Lang), st.Scaffold).Scan(&id)
	return id, err
}

// UpdateStep saves a step's content.
func (s *Store) UpdateStep(ctx context.Context, id int64, st Step) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE lesson_steps SET instruction=?, starter=?, checks=?, lang=?, scaffold=?, updated_at=datetime('now') WHERE id=?`,
		st.Instruction, st.Starter, st.Checks, langOr(st.Lang), st.Scaffold, id)
	return err
}

// DeleteStep removes a step and any progress against it.
func (s *Store) DeleteStep(ctx context.Context, id int64) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM step_progress WHERE step_id = ?`, id); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `DELETE FROM lesson_steps WHERE id = ?`, id)
	return err
}

// MoveStep swaps a step's order with its neighbour (dir -1 = up, +1 = down).
func (s *Store) MoveStep(ctx context.Context, id int64, dir int) error {
	cur, err := s.GetStep(ctx, id)
	if err != nil {
		return err
	}
	var nbID int64
	var nbSort int
	q := `SELECT id, sort FROM lesson_steps WHERE lesson_id = ? AND sort < ? ORDER BY sort DESC LIMIT 1`
	if dir > 0 {
		q = `SELECT id, sort FROM lesson_steps WHERE lesson_id = ? AND sort > ? ORDER BY sort ASC LIMIT 1`
	}
	if err := s.db.QueryRowContext(ctx, q, cur.LessonID, cur.Sort).Scan(&nbID, &nbSort); err != nil {
		return nil // no neighbour — nothing to do
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if _, err := tx.ExecContext(ctx, `UPDATE lesson_steps SET sort=? WHERE id=?`, nbSort, cur.ID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE lesson_steps SET sort=? WHERE id=?`, cur.Sort, nbID); err != nil {
		return err
	}
	return tx.Commit()
}

// MarkStepComplete records that a user finished a step (idempotent).
func (s *Store) MarkStepComplete(ctx context.Context, userID, stepID int64) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO step_progress (user_id, step_id) VALUES (?, ?)
		 ON CONFLICT(user_id, step_id) DO NOTHING`, userID, stepID)
	return err
}

// CompletedStepIDs returns the set of a lesson's steps the user has completed.
func (s *Store) CompletedStepIDs(ctx context.Context, userID, lessonID int64) (map[int64]bool, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.step_id FROM step_progress p
		JOIN lesson_steps s ON s.id = p.step_id
		WHERE p.user_id = ? AND s.lesson_id = ?`, userID, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]bool{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

// StepProgress is a {done,total} pair for a step-based lesson.
type StepProgress struct{ Done, Total int }

// StepProgressByLesson returns, per lesson in a course that HAS steps, how many
// steps exist and how many the user has completed (userID 0 = none done).
func (s *Store) StepProgressByLesson(ctx context.Context, userID, courseID int64) (map[int64]StepProgress, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.lesson_id, COUNT(*) AS total,
		       COUNT(p.step_id) AS done
		FROM lesson_steps s
		JOIN lessons l ON l.id = s.lesson_id
		LEFT JOIN step_progress p ON p.step_id = s.id AND p.user_id = ?
		WHERE l.course_id = ?
		GROUP BY s.lesson_id`, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]StepProgress{}
	for rows.Next() {
		var lid int64
		var sp StepProgress
		if err := rows.Scan(&lid, &sp.Total, &sp.Done); err != nil {
			return nil, err
		}
		out[lid] = sp
	}
	return out, rows.Err()
}

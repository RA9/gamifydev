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
}

// ListSteps returns a lesson's steps in order.
func (s *Store) ListSteps(ctx context.Context, lessonID int64) ([]Step, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, lesson_id, sort, instruction, starter, checks FROM lesson_steps
		 WHERE lesson_id = ? ORDER BY sort`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Step
	for rows.Next() {
		var st Step
		if err := rows.Scan(&st.ID, &st.LessonID, &st.Sort, &st.Instruction, &st.Starter, &st.Checks); err != nil {
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
		`SELECT id, lesson_id, sort, instruction, starter, checks FROM lesson_steps WHERE id = ?`, id).
		Scan(&st.ID, &st.LessonID, &st.Sort, &st.Instruction, &st.Starter, &st.Checks)
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
			`INSERT INTO lesson_steps (lesson_id, sort, instruction, starter, checks) VALUES (?, ?, ?, ?, ?)`,
			lessonID, i, st.Instruction, st.Starter, st.Checks); err != nil {
			return err
		}
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

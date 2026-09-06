package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrAccountRestricted = errors.New("account is not allowed to participate")
	ErrWorkLocked        = errors.New("scheduled work is not available yet")
)

type accountQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func effectiveAccountState(ctx context.Context, q accountQuerier, userID int64) (string, error) {
	var state string
	var expires sql.NullString
	err := q.QueryRowContext(ctx, `SELECT state, expires_at FROM account_status WHERE user_id = ?`, userID).Scan(&state, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return AccountActive, nil
	}
	if err != nil {
		return "", err
	}
	if state != AccountActive && expires.Valid {
		var elapsed int
		if err := q.QueryRowContext(ctx, `SELECT CASE WHEN datetime(?) <= datetime('now') THEN 1 ELSE 0 END`, expires.String).Scan(&elapsed); err != nil {
			return "", err
		}
		if elapsed == 1 {
			return AccountActive, nil
		}
	}
	return state, nil
}

func requireActiveAccount(ctx context.Context, q accountQuerier, userID int64) error {
	state, err := effectiveAccountState(ctx, q, userID)
	if err != nil {
		return err
	}
	if state != AccountActive {
		return ErrAccountRestricted
	}
	return nil
}

func (s *Store) RequireActiveAccount(ctx context.Context, userID int64) error {
	return requireActiveAccount(ctx, s.db, userID)
}

type scheduledResource struct {
	kind string
	id   int64
}

func (s *Store) RequireLessonAvailable(ctx context.Context, userID, lessonID int64) error {
	return s.requireScheduledWork(ctx, userID, scheduledResource{kind: "lesson", id: lessonID})
}

func (s *Store) RequireAssignmentAvailable(ctx context.Context, userID, assignmentID int64) error {
	return s.requireScheduledWork(ctx, userID, scheduledResource{kind: "checkpoint", id: assignmentID})
}

func (s *Store) RequireProblemAvailable(ctx context.Context, userID, problemID int64) error {
	return s.requireScheduledWork(ctx, userID, scheduledResource{kind: "problem", id: problemID})
}

func (s *Store) requireScheduledWork(ctx context.Context, userID int64, resource scheduledResource) error {
	if err := s.RequireActiveAccount(ctx, userID); err != nil {
		return err
	}
	var scheduleID int64
	var mode string
	var available, exempt int
	err := s.db.QueryRowContext(ctx, `
		SELECT cs.id, COALESCE(l.mode, a.mode, p.mode, 'theoretical'),
		       CASE WHEN date(cs.due_on) <= date('now') THEN 1 ELSE 0 END,
		       CASE WHEN EXISTS (
		         SELECT 1 FROM json_each(e.placement_exemptions) ex
		         JOIN courses ec ON ec.slug = CAST(ex.value AS TEXT)
		         WHERE ec.id = cs.course_id
		       ) THEN 1 ELSE 0 END
		FROM enrollments e
		JOIN cohort_schedule cs ON cs.cohort_id = e.cohort_id
		LEFT JOIN lessons l ON l.id = cs.lesson_id
		LEFT JOIN assignments a ON a.id = cs.assignment_id
		LEFT JOIN problems p ON p.id = cs.problem_id
		WHERE e.user_id = ? AND e.state = 'active'
		  AND ((? = 'lesson' AND cs.lesson_id = ?)
		    OR (? = 'checkpoint' AND cs.assignment_id = ?)
		    OR (? = 'problem' AND cs.problem_id = ?))
		LIMIT 1`, userID, resource.kind, resource.id, resource.kind, resource.id, resource.kind, resource.id).
		Scan(&scheduleID, &mode, &available, &exempt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	if exempt == 1 || mode == "theoretical" {
		return nil
	}
	if available == 0 {
		return ErrWorkLocked
	}
	enrollment, err := s.LiveEnrollment(ctx, userID)
	if err != nil {
		return err
	}
	if enrollment == nil || !enrollment.CohortID.Valid {
		return ErrWorkLocked
	}
	items, err := s.CohortSchedule(ctx, enrollment.CohortID.Int64, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == scheduleID {
			return nil
		}
		if !item.Exempt && item.Practical() && !item.Done {
			return ErrWorkLocked
		}
	}
	return ErrWorkLocked
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RA9/gamifydev/platform/internal/attendance"
)

// Attendance resolution, sanctions, and appeals.
//
// Everything here is written to be safe to run before the policy is trusted:
// resolution is idempotent, and sanctions carry a `shadow` flag that defaults to
// on. A shadow sanction is a record of what the rule *would* have done.

// --- Absence notices ---------------------------------------------------------

// AbsenceNotice is a learner telling their cohort they'll be away.
type AbsenceNotice struct {
	ID        int64
	UserID    int64
	FromDay   string
	ToDay     string
	Reason    string
	CreatedAt string
}

// FileAbsence records an absence notice. Filing after the fact is allowed: the
// behaviour being taught is telling your team, not predicting your life.
func (s *Store) FileAbsence(ctx context.Context, userID, cohortID int64, from, to, reason string) error {
	var cohort any
	if cohortID != 0 {
		cohort = cohortID
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO absence_notices (user_id, cohort_id, from_day, to_day, reason)
		VALUES (?, ?, date(?), date(?), ?)`, userID, cohort, from, to, reason)
	return err
}

// AbsenceNotices lists a learner's notices, newest first.
func (s *Store) AbsenceNotices(ctx context.Context, userID int64, limit int) ([]AbsenceNotice, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, from_day, to_day, reason, created_at
		FROM absence_notices WHERE user_id = ?
		ORDER BY from_day DESC, id DESC LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AbsenceNotice
	for rows.Next() {
		var n AbsenceNotice
		if err := rows.Scan(&n.ID, &n.UserID, &n.FromDay, &n.ToDay, &n.Reason, &n.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// --- Resolution --------------------------------------------------------------

// ResolveAttendance records daily state for every active learner in a cohort,
// for each scheduled day that has passed.
//
// Only scheduled days are considered — weekends and rest days carry no work, so
// counting them would punish the rest the schedule deliberately builds in.
//
// Idempotent: re-resolving a day overwrites it with the same answer, so the job
// can safely re-scan a trailing window.
func (s *Store) ResolveAttendance(ctx context.Context, cohortID int64, since string) (int, error) {
	// Scheduled days in range that have already passed.
	dayRows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT due_on FROM cohort_schedule
		WHERE cohort_id = ? AND date(due_on) >= date(?) AND date(due_on) < date('now')
		ORDER BY due_on`, cohortID, since)
	if err != nil {
		return 0, err
	}
	var days []string
	for dayRows.Next() {
		var d string
		if err := dayRows.Scan(&d); err != nil {
			dayRows.Close()
			return 0, err
		}
		days = append(days, d)
	}
	dayRows.Close()
	if err := dayRows.Err(); err != nil {
		return 0, err
	}
	if len(days) == 0 {
		return 0, nil
	}

	memberRows, err := s.db.QueryContext(ctx, `
		SELECT user_id FROM cohort_members
		WHERE cohort_id = ? AND left_at IS NULL AND role = 'learner'`, cohortID)
	if err != nil {
		return 0, err
	}
	var members []int64
	for memberRows.Next() {
		var uid int64
		if err := memberRows.Scan(&uid); err != nil {
			memberRows.Close()
			return 0, err
		}
		members = append(members, uid)
	}
	memberRows.Close()
	if err := memberRows.Err(); err != nil {
		return 0, err
	}

	n := 0
	for _, uid := range members {
		for _, day := range days {
			var excused, active int
			if err := s.db.QueryRowContext(ctx, `
				SELECT
				  (SELECT COUNT(*) FROM absence_notices
				     WHERE user_id = ? AND date(?) BETWEEN date(from_day) AND date(to_day)),
				  (SELECT COUNT(*) FROM activity_days WHERE user_id = ? AND day = date(?))`,
				uid, day, uid, day).Scan(&excused, &active); err != nil {
				return 0, err
			}
			state := attendance.Resolve(excused > 0, active > 0)
			if _, err := s.db.ExecContext(ctx, `
				INSERT INTO attendance_marks (user_id, cohort_id, day, state, resolved_at)
				VALUES (?, ?, date(?), ?, datetime('now'))
				ON CONFLICT(user_id, day) DO UPDATE SET
					state = excluded.state, cohort_id = excluded.cohort_id,
					resolved_at = excluded.resolved_at`,
				uid, cohortID, day, state); err != nil {
				return 0, err
			}
			n++
		}
	}
	return n, nil
}

// AttendanceDays returns a learner's resolved days in chronological order.
func (s *Store) AttendanceDays(ctx context.Context, userID int64) ([]attendance.Day, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT day, state FROM attendance_marks WHERE user_id = ? ORDER BY day`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []attendance.Day
	for rows.Next() {
		var d attendance.Day
		if err := rows.Scan(&d.Day, &d.State); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// --- Sanctions ---------------------------------------------------------------

// Sanction is a recorded consequence — real or, by default, hypothetical.
type Sanction struct {
	ID         int64
	UserID     int64
	UserName   string
	CohortID   sql.NullInt64
	Kind       string
	Reason     string
	WindowFrom string
	WindowTo   string
	Shadow     bool
	AppliedAt  sql.NullString
	ExpiresAt  sql.NullString
	CreatedAt  string
	// Populated by listings.
	AppealState string
}

// RecordSanction stores a sanction for one absence run. The unique index on
// (user, window) makes re-evaluation idempotent — the same stretch of silence
// never produces two records.
//
// `shadow` false means the sanction is actually applied; the caller is
// responsible for having checked that enforcement is enabled.
func (s *Store) RecordSanction(ctx context.Context, sc Sanction, shadow bool) (int64, bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback() //nolint:errcheck

	shadowV := 1
	var enrollmentID any
	if !shadow {
		shadowV = 0
		if sc.Kind == attendance.KindDrop {
			enrollment, err := liveEnrollment(ctx, tx, sc.UserID)
			if err != nil {
				return 0, false, err
			}
			if enrollment == nil || (sc.CohortID.Valid && (!enrollment.CohortID.Valid || enrollment.CohortID.Int64 != sc.CohortID.Int64)) {
				return 0, false, ErrNotFound
			}
			enrollmentID = enrollment.ID
		}
	}
	var cohort any
	if sc.CohortID.Valid {
		cohort = sc.CohortID.Int64
	}
	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO sanctions
			(user_id, cohort_id, enrollment_id, kind, reason, window_from, window_to, shadow)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, window_from, window_to) DO NOTHING
		RETURNING id`,
		sc.UserID, cohort, enrollmentID, sc.Kind, sc.Reason, sc.WindowFrom, sc.WindowTo, shadowV).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil // already recorded for this window
	}
	if err != nil {
		return 0, false, err
	}
	if !shadow {
		if err := applySanctionTx(ctx, tx, id); err != nil {
			return 0, false, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, false, err
	}
	return id, true, nil
}

// ApplySanction carries out a previously recorded real sanction. Real sanctions
// created by RecordSanction are applied in that same transaction; this method is
// retained as an idempotent operator/recovery entry point.
func (s *Store) ApplySanction(ctx context.Context, sanctionID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck
	if err := applySanctionTx(ctx, tx, sanctionID); err != nil {
		return err
	}
	return tx.Commit()
}

func applySanctionTx(ctx context.Context, tx *sql.Tx, sanctionID int64) error {
	var userID int64
	var cohortID, enrollmentID sql.NullInt64
	var appliedAt sql.NullString
	var kind string
	if err := tx.QueryRowContext(ctx, `
		SELECT user_id, cohort_id, enrollment_id, kind, applied_at
		FROM sanctions WHERE id = ? AND shadow = 0`, sanctionID).
		Scan(&userID, &cohortID, &enrollmentID, &kind, &appliedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if appliedAt.Valid {
		return nil
	}
	if kind == attendance.KindDrop {
		if !enrollmentID.Valid {
			return ErrNotFound
		}
		e, err := scanEnrollment(tx.QueryRowContext(ctx, `
			SELECT `+enrollCols+` FROM enrollments
			WHERE id = ? AND user_id = ? AND state IN ('placed','active','paused')`,
			enrollmentID.Int64, userID))
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if cohortID.Valid && (!e.CohortID.Valid || e.CohortID.Int64 != cohortID.Int64) {
			return ErrNotFound
		}
		allowed := false
		for _, candidate := range enrollTransitions[e.State] {
			if candidate == EnrollDropped {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("%w: %s -> %s", ErrBadTransition, e.State, EnrollDropped)
		}
		if cohortID.Valid {
			if _, err := tx.ExecContext(ctx, `
				UPDATE cohort_members SET left_at = datetime('now')
				WHERE cohort_id = ? AND user_id = ? AND left_at IS NULL`, cohortID.Int64, userID); err != nil {
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE enrollments
			SET state = 'dropped', reason = 'attendance', ended_at = datetime('now'),
			    updated_at = datetime('now')
			WHERE id = ?`, e.ID); err != nil {
			return err
		}
	}
	_, err := tx.ExecContext(ctx,
		`UPDATE sanctions SET applied_at = datetime('now') WHERE id = ? AND applied_at IS NULL`, sanctionID)
	return err
}

// LatestSanction returns a learner's most recent sanction, or nil.
func (s *Store) LatestSanction(ctx context.Context, userID int64) (*Sanction, error) {
	var sc Sanction
	err := s.db.QueryRowContext(ctx, `
		SELECT s.id, s.user_id, s.cohort_id, s.kind, s.reason, s.window_from, s.window_to,
		       s.shadow, s.applied_at, s.expires_at, s.created_at,
		       COALESCE(a.state, '')
		FROM sanctions s LEFT JOIN appeals a ON a.sanction_id = s.id
		WHERE s.user_id = ? ORDER BY s.created_at DESC, s.id DESC LIMIT 1`, userID).
		Scan(&sc.ID, &sc.UserID, &sc.CohortID, &sc.Kind, &sc.Reason, &sc.WindowFrom, &sc.WindowTo,
			&sc.Shadow, &sc.AppliedAt, &sc.ExpiresAt, &sc.CreatedAt, &sc.AppealState)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &sc, err
}

// ListSanctions returns recent sanctions for the admin view. `shadowOnly`
// restricts to hypothetical ones — the calibration view.
func (s *Store) ListSanctions(ctx context.Context, shadowOnly bool, limit int) ([]Sanction, error) {
	q := `
		SELECT s.id, s.user_id, u.name, s.cohort_id, s.kind, s.reason,
		       s.window_from, s.window_to, s.shadow, s.applied_at, s.expires_at, s.created_at,
		       COALESCE(a.state, '')
		FROM sanctions s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN appeals a ON a.sanction_id = s.id`
	if shadowOnly {
		q += ` WHERE s.shadow = 1`
	}
	q += ` ORDER BY s.created_at DESC, s.id DESC LIMIT ?`
	rows, err := s.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Sanction
	for rows.Next() {
		var sc Sanction
		if err := rows.Scan(&sc.ID, &sc.UserID, &sc.UserName, &sc.CohortID, &sc.Kind, &sc.Reason,
			&sc.WindowFrom, &sc.WindowTo, &sc.Shadow, &sc.AppliedAt, &sc.ExpiresAt,
			&sc.CreatedAt, &sc.AppealState); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

// --- Appeals -----------------------------------------------------------------

// Appeal is a learner contesting a sanction.
type Appeal struct {
	ID           int64
	SanctionID   int64
	UserID       int64
	UserName     string
	Body         string
	State        string
	DecisionNote string
	DecidedAt    sql.NullString
	CreatedAt    string
	// Context for the reviewer.
	SanctionKind   string
	SanctionReason string
	WindowFrom     string
	WindowTo       string
	Shadow         bool
}

// FileAppeal opens an appeal against a sanction. One per sanction.
func (s *Store) FileAppeal(ctx context.Context, sanctionID, userID int64, body string) error {
	var owner int64
	if err := s.db.QueryRowContext(ctx,
		`SELECT user_id FROM sanctions WHERE id = ?`, sanctionID).Scan(&owner); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if owner != userID {
		return ErrNotFound // don't confirm someone else's sanction exists
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO appeals (sanction_id, user_id, body) VALUES (?, ?, ?)
		ON CONFLICT(sanction_id) DO UPDATE SET body = excluded.body`,
		sanctionID, userID, body)
	return err
}

// OpenAppeals lists appeals awaiting a human decision.
func (s *Store) OpenAppeals(ctx context.Context) ([]Appeal, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.id, a.sanction_id, a.user_id, u.name, a.body, a.state,
		       a.decision_note, a.decided_at, a.created_at,
		       s.kind, s.reason, s.window_from, s.window_to, s.shadow
		FROM appeals a
		JOIN users u ON u.id = a.user_id
		JOIN sanctions s ON s.id = a.sanction_id
		WHERE a.state = 'open' ORDER BY a.created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Appeal
	for rows.Next() {
		var a Appeal
		if err := rows.Scan(&a.ID, &a.SanctionID, &a.UserID, &a.UserName, &a.Body, &a.State,
			&a.DecisionNote, &a.DecidedAt, &a.CreatedAt,
			&a.SanctionKind, &a.SanctionReason, &a.WindowFrom, &a.WindowTo, &a.Shadow); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// DecideAppeal records a human's decision.
//
// Granting an appeal reverses an applied sanction: the learner's enrollment is
// restored to active. Upholding leaves it in place. Either way a person decided,
// which is the property that matters.
func (s *Store) DecideAppeal(ctx context.Context, appealID, deciderID int64, grant bool, note string) error {
	state := "upheld"
	if grant {
		state = "granted"
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	var userID int64
	var cohortID, enrollmentID sql.NullInt64
	var shadow bool
	if err := tx.QueryRowContext(ctx, `
		SELECT a.user_id, s.cohort_id, s.shadow, s.enrollment_id
		FROM appeals a JOIN sanctions s ON s.id = a.sanction_id
		WHERE a.id = ? AND a.state = 'open'`, appealID).
		Scan(&userID, &cohortID, &shadow, &enrollmentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if grant && !shadow {
		if !enrollmentID.Valid {
			return ErrNotFound
		}
		// A blank unplaced row can be created passively while viewing placement. It
		// carries no application decision and must not block a valid appeal. Any
		// meaningful newer enrollment remains a conflict for an operator to resolve.
		live, err := liveEnrollment(ctx, tx, userID)
		if err != nil {
			return err
		}
		if live != nil {
			blank := live.State == EnrollUnplaced && !live.PathID.Valid &&
				!live.PlacementResultID.Valid && !live.CohortID.Valid && live.TZBand == ""
			if !blank {
				return ErrEnrollmentLocked
			}
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM enrollments WHERE id = ? AND state = 'unplaced'`, live.ID); err != nil {
				return err
			}
		}

		// Restore only the exact enrollment ended by this sanction. A fresh row
		// preserves terminal history while copying the generated binding.
		var droppedID int64
		err = tx.QueryRowContext(ctx, `
			SELECT id FROM enrollments
			WHERE id = ? AND user_id = ? AND state = 'dropped'
			  AND (? IS NULL OR cohort_id = ?)`,
			enrollmentID.Int64, userID, cohortID, cohortID).Scan(&droppedID)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO enrollments
				(user_id, path_id, placement_result_id, state, reason, placed_at,
				 started_at, ended_at, cohort_id, tz_band, placement_exemptions)
			SELECT user_id, path_id, placement_result_id, 'active', 'appeal granted',
			       placed_at, COALESCE(started_at, datetime('now')), NULL,
			       cohort_id, tz_band, placement_exemptions
			FROM enrollments WHERE id = ?`, droppedID); err != nil {
			return err
		}
		if cohortID.Valid {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO cohort_members (cohort_id, user_id, role) VALUES (?, ?, ?)
				ON CONFLICT(cohort_id, user_id) WHERE left_at IS NULL DO NOTHING`,
				cohortID.Int64, userID, RoleLearner); err != nil {
				return err
			}
		}
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE appeals
		SET state = ?, decided_by = ?, decided_at = datetime('now'), decision_note = ?
		WHERE id = ? AND state = 'open'`, state, deciderID, note, appealID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return ErrNotFound
	}
	return tx.Commit()
}

// AttendanceOverview summarises the program for the admin calibration view.
type AttendanceOverview struct {
	Resolved, Present, Excused, Absent int
	ShadowSanctions, RealSanctions     int
	OpenAppeals                        int
}

// AttendanceStats aggregates resolved marks and sanction counts.
func (s *Store) AttendanceStats(ctx context.Context) (AttendanceOverview, error) {
	var o AttendanceOverview
	err := s.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM attendance_marks),
		  (SELECT COUNT(*) FROM attendance_marks WHERE state = 'present'),
		  (SELECT COUNT(*) FROM attendance_marks WHERE state = 'excused'),
		  (SELECT COUNT(*) FROM attendance_marks WHERE state = 'absent'),
		  (SELECT COUNT(*) FROM sanctions WHERE shadow = 1),
		  (SELECT COUNT(*) FROM sanctions WHERE shadow = 0),
		  (SELECT COUNT(*) FROM appeals WHERE state = 'open')`).
		Scan(&o.Resolved, &o.Present, &o.Excused, &o.Absent,
			&o.ShadowSanctions, &o.RealSanctions, &o.OpenAppeals)
	return o, err
}

// LearnersToEvaluate lists active cohort learners for the evaluator.
func (s *Store) LearnersToEvaluate(ctx context.Context) ([]struct {
	UserID, CohortID int64
}, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT m.user_id, m.cohort_id FROM cohort_members m
		JOIN cohorts c ON c.id = m.cohort_id
		WHERE m.left_at IS NULL AND m.role = 'learner' AND c.state = 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct{ UserID, CohortID int64 }
	for rows.Next() {
		var r struct{ UserID, CohortID int64 }
		if err := rows.Scan(&r.UserID, &r.CohortID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

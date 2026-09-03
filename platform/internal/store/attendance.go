package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	shadowV := 1
	var appliedAt any
	if !shadow {
		shadowV = 0
		appliedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	}
	var cohort any
	if sc.CohortID.Valid {
		cohort = sc.CohortID.Int64
	}
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO sanctions (user_id, cohort_id, kind, reason, window_from, window_to, shadow, applied_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, window_from, window_to) DO NOTHING`,
		sc.UserID, cohort, sc.Kind, sc.Reason, sc.WindowFrom, sc.WindowTo, shadowV, appliedAt)
	if err != nil {
		return 0, false, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return 0, false, nil // already recorded for this window
	}
	id, _ := res.LastInsertId()
	return id, true, nil
}

// ApplySanction carries out a real (non-shadow) drop: the learner leaves the
// cohort and their enrollment ends.
//
// Deliberately narrow — it drops, it does not ban. A dropped learner keeps their
// account and may reapply.
func (s *Store) ApplySanction(ctx context.Context, sanctionID int64) error {
	var userID int64
	var cohortID sql.NullInt64
	var kind string
	if err := s.db.QueryRowContext(ctx,
		`SELECT user_id, cohort_id, kind FROM sanctions WHERE id = ? AND shadow = 0`,
		sanctionID).Scan(&userID, &cohortID, &kind); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if kind != attendance.KindDrop {
		return nil
	}
	if cohortID.Valid {
		if err := s.RemoveMember(ctx, cohortID.Int64, userID); err != nil {
			return err
		}
	}
	e, err := s.LiveEnrollment(ctx, userID)
	if err != nil || e == nil {
		return err
	}
	return s.TransitionEnrollment(ctx, e.ID, EnrollDropped, "attendance")
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
	var sanctionID, userID int64
	if err := s.db.QueryRowContext(ctx,
		`SELECT sanction_id, user_id FROM appeals WHERE id = ?`, appealID).
		Scan(&sanctionID, &userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if _, err := s.db.ExecContext(ctx, `
		UPDATE appeals SET state = ?, decided_by = ?, decided_at = datetime('now'), decision_note = ?
		WHERE id = ?`, state, deciderID, note, appealID); err != nil {
		return err
	}
	if !grant {
		return nil
	}
	// Reverse the sanction: rejoin the cohort if it still exists, and revive the
	// enrollment.
	var cohortID sql.NullInt64
	var shadow bool
	_ = s.db.QueryRowContext(ctx,
		`SELECT cohort_id, shadow FROM sanctions WHERE id = ?`, sanctionID).Scan(&cohortID, &shadow)
	if shadow {
		return nil // nothing was ever applied
	}
	if cohortID.Valid {
		if err := s.AddMember(ctx, cohortID.Int64, userID, RoleLearner); err != nil {
			return err
		}
	}
	// A dropped enrollment is terminal, so restoration means a fresh one placed
	// back on the same path — history is preserved rather than rewritten.
	var pathID sql.NullInt64
	_ = s.db.QueryRowContext(ctx,
		`SELECT path_id FROM enrollments WHERE user_id = ? ORDER BY id DESC LIMIT 1`, userID).Scan(&pathID)
	e, err := s.EnsureEnrollment(ctx, userID)
	if err != nil {
		return err
	}
	if e.State == EnrollUnplaced && pathID.Valid {
		if err := s.PlaceEnrollment(ctx, e.ID, pathID.Int64, "appeal granted"); err != nil {
			return err
		}
	}
	if e2, _ := s.LiveEnrollment(ctx, userID); e2 != nil && e2.State == EnrollPlaced {
		if cohortID.Valid {
			_, _ = s.db.ExecContext(ctx,
				`UPDATE enrollments SET cohort_id = ? WHERE id = ?`, cohortID.Int64, e2.ID)
		}
		return s.TransitionEnrollment(ctx, e2.ID, EnrollActive, "appeal granted")
	}
	return nil
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

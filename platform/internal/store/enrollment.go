package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// --- Account participation ---------------------------------------------------

// Account states. A missing account_status row means AccountActive, so existing
// users need no backfill.
const (
	AccountActive    = "active"
	AccountSuspended = "suspended"
	AccountBanned    = "banned"
)

// AccountState reports whether a user may currently participate. A suspension
// whose expires_at has passed reads as active — sanction:expire tidies the row
// up later, but the read path must never depend on that job having run.
func (s *Store) AccountState(ctx context.Context, userID int64) (string, error) {
	var state string
	var expires sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT state, expires_at FROM account_status WHERE user_id = ?`, userID).
		Scan(&state, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return AccountActive, nil
	}
	if err != nil {
		return AccountActive, err
	}
	if state != AccountActive && expires.Valid {
		var elapsed int
		if err := s.db.QueryRowContext(ctx,
			`SELECT CASE WHEN datetime(?) <= datetime('now') THEN 1 ELSE 0 END`, expires.String).
			Scan(&elapsed); err == nil && elapsed == 1 {
			return AccountActive, nil
		}
	}
	return state, nil
}

// SetAccountState records an account-level participation change.
func (s *Store) SetAccountState(ctx context.Context, userID int64, state, reason string, expiresAt *string) error {
	switch state {
	case AccountActive, AccountSuspended, AccountBanned:
	default:
		return fmt.Errorf("unknown account state %q", state)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO account_status (user_id, state, reason, effective_at, expires_at, updated_at)
		VALUES (?, ?, ?, datetime('now'), ?, datetime('now'))
		ON CONFLICT(user_id) DO UPDATE SET
			state=excluded.state, reason=excluded.reason,
			effective_at=excluded.effective_at, expires_at=excluded.expires_at,
			updated_at=excluded.updated_at`,
		userID, state, reason, expiresAt)
	return err
}

// --- Enrollment --------------------------------------------------------------

// Enrollment states.
const (
	EnrollUnplaced  = "unplaced"  // registered, diagnostic not yet taken
	EnrollPlaced    = "placed"    // path decided, waiting for a cohort
	EnrollActive    = "active"    // cohort running
	EnrollPaused    = "paused"    // deferred; may rejoin a later cohort
	EnrollDropped   = "dropped"   // terminal
	EnrollCompleted = "completed" // terminal
)

// enrollTransitions is the whole state machine. Keeping it as data (rather than
// scattered if-statements at call sites) is what makes the rule enforceable:
// every change goes through TransitionEnrollment.
var enrollTransitions = map[string][]string{
	EnrollUnplaced: {EnrollPlaced},
	EnrollPlaced:   {EnrollActive, EnrollPaused, EnrollDropped},
	EnrollActive:   {EnrollPaused, EnrollDropped, EnrollCompleted},
	EnrollPaused:   {EnrollActive, EnrollPlaced, EnrollDropped},
	// terminal
	EnrollDropped:   {},
	EnrollCompleted: {},
}

var (
	// ErrBadTransition is returned when a caller asks for a move the state
	// machine forbids (e.g. reviving a dropped enrollment).
	ErrBadTransition = errors.New("illegal enrollment transition")
	// ErrEnrollmentLocked prevents an existing queued, active, or paused
	// enrollment from being rebound to a different placement decision.
	ErrEnrollmentLocked = errors.New("enrollment is already bound to a placement result")
	// ErrEnrollmentPlacementInvalid means the selected result cannot enroll this
	// user: it is missing, failed, stale, belongs to someone else, or has no
	// published generated path.
	ErrEnrollmentPlacementInvalid = errors.New("placement result is not eligible for enrollment")
)

// EnrollmentTerminal reports whether a state ends the enrollment's life.
func EnrollmentTerminal(state string) bool {
	return state == EnrollDropped || state == EnrollCompleted
}

// Enrollment is a learner's relationship to a path.
type Enrollment struct {
	ID                int64
	UserID            int64
	PathID            sql.NullInt64
	PlacementResultID sql.NullInt64
	State             string
	Reason            string
	PlacedAt          sql.NullString
	StartedAt         sql.NullString
	EndedAt           sql.NullString
	// Set once cohort formation places the learner (phase 3).
	CohortID sql.NullInt64
	TZBand   string
	// Exemptions is copied from PlacementResultID when the learner enrolls. It is
	// deliberately a snapshot: later placement activity cannot rewrite a running
	// cohort schedule.
	Exemptions []string
}

const enrollCols = `id, user_id, path_id, placement_result_id, state, reason, placed_at, started_at, ended_at, cohort_id, tz_band, placement_exemptions`

func scanEnrollment(row interface{ Scan(...any) error }) (*Enrollment, error) {
	var e Enrollment
	var exemptions string
	err := row.Scan(&e.ID, &e.UserID, &e.PathID, &e.PlacementResultID, &e.State, &e.Reason,
		&e.PlacedAt, &e.StartedAt, &e.EndedAt, &e.CohortID, &e.TZBand, &exemptions)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(exemptions), &e.Exemptions); err != nil {
		return nil, fmt.Errorf("decode enrollment exemptions: %w", err)
	}
	return &e, nil
}

// LiveEnrollment returns the user's current non-terminal enrollment, or nil if
// they have none (never registered for a path, or previously dropped).
type enrollmentQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func liveEnrollment(ctx context.Context, q enrollmentQuerier, userID int64) (*Enrollment, error) {
	e, err := scanEnrollment(q.QueryRowContext(ctx,
		`SELECT `+enrollCols+` FROM enrollments
		 WHERE user_id = ? AND state IN ('unplaced','placed','active','paused')`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return e, err
}

func (s *Store) LiveEnrollment(ctx context.Context, userID int64) (*Enrollment, error) {
	return liveEnrollment(ctx, s.db, userID)
}

// enrollmentAfterRace turns a concurrent duplicate enrollment into the same
// idempotent result as a sequential retry after the losing transaction releases
// its snapshot. A competing different binding remains locked.
func (s *Store) enrollmentAfterRace(ctx context.Context, userID, resultID int64, tzBand string, cause error) (*Enrollment, error) {
	for attempt := 0; attempt < 5; attempt++ {
		enrollment, err := s.LiveEnrollment(ctx, userID)
		if err != nil {
			return nil, cause
		}
		if enrollment != nil && enrollment.State != EnrollUnplaced {
			if enrollment.PlacementResultID.Valid && enrollment.PlacementResultID.Int64 == resultID && enrollment.TZBand == tzBand {
				return enrollment, nil
			}
			return nil, ErrEnrollmentLocked
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
	return nil, cause
}

// EnsureEnrollment returns the user's live enrollment, creating an `unplaced`
// one if they have none. Idempotent: safe to call on every request.
func (s *Store) EnsureEnrollment(ctx context.Context, userID int64) (*Enrollment, error) {
	if e, err := s.LiveEnrollment(ctx, userID); err != nil || e != nil {
		return e, err
	}
	// The partial unique index makes the concurrent case a conflict rather than
	// a duplicate; fall back to re-reading if we lost the race.
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO enrollments (user_id, state) VALUES (?, ?)`, userID, EnrollUnplaced)
	if err != nil {
		if e, rerr := s.LiveEnrollment(ctx, userID); rerr == nil && e != nil {
			return e, nil
		}
		return nil, err
	}
	return s.LiveEnrollment(ctx, userID)
}

// EnrollFromPlacement atomically binds a learner's live enrollment to the exact
// latest passing placement result, its generated path, its exemption snapshot,
// and the selected timezone band. Once placed, the binding is immutable; a
// repeated submission of the same result and band is an idempotent no-op.
func (s *Store) EnrollFromPlacement(ctx context.Context, userID, resultID int64, tzBand string) (*Enrollment, error) {
	tzBand = strings.TrimSpace(tzBand)
	if userID <= 0 || resultID <= 0 || tzBand == "" {
		return nil, ErrEnrollmentPlacementInvalid
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	enrollment, err := liveEnrollment(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	if enrollment != nil && enrollment.State != EnrollUnplaced {
		if enrollment.PlacementResultID.Valid && enrollment.PlacementResultID.Int64 == resultID && enrollment.TZBand == tzBand {
			return enrollment, nil
		}
		return nil, ErrEnrollmentLocked
	}

	var pathID, attemptID int64
	var exemptions string
	err = tx.QueryRowContext(ctx, `
		SELECT pr.recommended_path_id, pr.attempt_id, pr.exemptions
		FROM placement_results pr
		JOIN assessment_attempts a ON a.id = pr.attempt_id
		JOIN assessments bank ON bank.id = a.assessment_id
		JOIN paths p ON p.id = pr.recommended_path_id
		WHERE pr.id = ? AND pr.user_id = ? AND pr.guest_id IS NULL
		  AND pr.passed = 1 AND a.user_id = ? AND a.guest_id IS NULL
		  AND a.submitted_at IS NOT NULL AND bank.kind = 'placement'
		  AND p.published = 1
		  AND pr.id = (
			SELECT latest.id FROM placement_results latest
			WHERE latest.user_id = ?
			ORDER BY datetime(latest.created_at) DESC, latest.id DESC LIMIT 1
		  )`, resultID, userID, userID, userID).Scan(&pathID, &attemptID, &exemptions)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEnrollmentPlacementInvalid
	}
	if err != nil {
		return nil, err
	}
	var decoded []string
	if err := json.Unmarshal([]byte(exemptions), &decoded); err != nil {
		return nil, fmt.Errorf("decode placement exemptions: %w", err)
	}

	if enrollment == nil {
		enrollment, err = scanEnrollment(tx.QueryRowContext(ctx, `
			INSERT INTO enrollments
				(user_id, path_id, placement_result_id, state, reason, placed_at, tz_band, placement_exemptions)
			VALUES (?, ?, ?, 'placed', 'placement result', datetime('now'), ?, ?)
			RETURNING `+enrollCols, userID, pathID, resultID, tzBand, exemptions))
	} else {
		enrollment, err = scanEnrollment(tx.QueryRowContext(ctx, `
			UPDATE enrollments
			SET path_id = ?, placement_result_id = ?, state = 'placed',
			    reason = 'placement result', placed_at = datetime('now'),
			    tz_band = ?, placement_exemptions = ?, updated_at = datetime('now')
			WHERE id = ? AND state = 'unplaced'
			RETURNING `+enrollCols, pathID, resultID, tzBand, exemptions, enrollment.ID))
	}
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrEnrollmentLocked
		}
		return s.enrollmentAfterRace(ctx, userID, resultID, tzBand, err)
	}

	// Close any placement sitting that raced with enrollment. StartAttemptFor
	// checks the same live-enrollment invariant; this update closes the opposite
	// ordering, where the sitting obtained its write lock first.
	if _, err := tx.ExecContext(ctx, `
		UPDATE assessment_attempts
		SET abandoned = 1
		WHERE user_id = ? AND id != ? AND submitted_at IS NULL AND abandoned = 0
		  AND assessment_id IN (SELECT id FROM assessments WHERE kind = 'placement')`,
		userID, attemptID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return s.enrollmentAfterRace(ctx, userID, resultID, tzBand, err)
	}
	return enrollment, nil
}

// TransitionEnrollment moves an enrollment to `to`, rejecting any move the state
// machine forbids. It stamps the timestamp column that corresponds to the new
// state and is the only supported way to change enrollment state.
func (s *Store) TransitionEnrollment(ctx context.Context, id int64, to, reason string) error {
	if _, ok := enrollTransitions[to]; !ok {
		return fmt.Errorf("%w: unknown state %q", ErrBadTransition, to)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	var from string
	if err := tx.QueryRowContext(ctx, `SELECT state FROM enrollments WHERE id = ?`, id).Scan(&from); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	allowed := false
	for _, cand := range enrollTransitions[from] {
		if cand == to {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("%w: %s -> %s", ErrBadTransition, from, to)
	}

	// Stamp the timestamp this state is defined by.
	set := "state = ?, reason = ?, updated_at = datetime('now')"
	switch to {
	case EnrollPlaced:
		set += ", placed_at = COALESCE(placed_at, datetime('now'))"
	case EnrollActive:
		set += ", started_at = COALESCE(started_at, datetime('now')), ended_at = NULL"
	}
	if EnrollmentTerminal(to) {
		set += ", ended_at = datetime('now')"
	}
	if _, err := tx.ExecContext(ctx, `UPDATE enrollments SET `+set+` WHERE id = ?`, to, reason, id); err != nil {
		return err
	}
	return tx.Commit()
}

// PlaceEnrollment is the legacy/operator placement path. Generated learner
// enrollment must use EnrollFromPlacement so the result, path, exemptions, and
// timezone are committed together. This method still updates path and state in
// one transaction so a rejected transition cannot partially mutate the row. It
// cannot rewrite an enrollment already bound to a generated placement result.
func (s *Store) PlaceEnrollment(ctx context.Context, id, pathID int64, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	var from string
	var placementResultID sql.NullInt64
	if err := tx.QueryRowContext(ctx,
		`SELECT state, placement_result_id FROM enrollments WHERE id = ?`, id).
		Scan(&from, &placementResultID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if placementResultID.Valid {
		return ErrEnrollmentLocked
	}
	allowed := false
	for _, candidate := range enrollTransitions[from] {
		if candidate == EnrollPlaced {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("%w: %s -> %s", ErrBadTransition, from, EnrollPlaced)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE enrollments
		SET path_id = ?, state = 'placed', reason = ?, cohort_id = NULL,
		    placed_at = COALESCE(placed_at, datetime('now')), updated_at = datetime('now')
		WHERE id = ?`, pathID, reason, id); err != nil {
		return err
	}
	return tx.Commit()
}

// --- Activity ledger ---------------------------------------------------------

// RollupActivity materializes activity_days from the signals the platform
// already records, for days on or after `since` (a YYYY-MM-DD string).
//
// Signals counted as "showed up and did something":
//   - completing a lesson step  (step_progress.completed_at)
//   - submitting an assignment  (submissions.created_at)
//   - finishing a lesson        (lesson_progress.completed_at)
//
// Standup posts are written straight to the ledger by PostStandup rather than
// rolled up, since they have no other durable home.
//
// Returns the number of (user, day) rows inserted. Idempotent — re-running over
// the same window inserts nothing new, so the job can safely overlap windows.
func (s *Store) RollupActivity(ctx context.Context, since string) (int, error) {
	// Each source contributes (user_id, day, tag); the union is grouped so a day
	// with several kinds of activity records all of them in `sources`.
	const q = `
		INSERT INTO activity_days (user_id, day, sources)
		SELECT user_id, day, GROUP_CONCAT(DISTINCT tag)
		FROM (
			SELECT p.user_id AS user_id, date(p.completed_at) AS day, 'step' AS tag
			FROM step_progress p
			WHERE p.completed_at IS NOT NULL AND date(p.completed_at) >= date(?)
			UNION ALL
			SELECT sub.user_id, date(sub.created_at), 'submission'
			FROM submissions sub
			WHERE date(sub.created_at) >= date(?)
			UNION ALL
			SELECT lp.user_id, date(lp.completed_at), 'lesson'
			FROM lesson_progress lp
			WHERE date(lp.completed_at) >= date(?)
		)
		GROUP BY user_id, day
		ON CONFLICT(user_id, day) DO NOTHING`
	res, err := s.db.ExecContext(ctx, q, since, since, since)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

// ActiveDaysBetween returns the days (YYYY-MM-DD, ascending) a user was active
// within an inclusive date range. Phase 5's attendance window reads this.
func (s *Store) ActiveDaysBetween(ctx context.Context, userID int64, from, to string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT day FROM activity_days
		 WHERE user_id = ? AND day >= date(?) AND day <= date(?) ORDER BY day`,
		userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// MarkActive records activity for a user on a given day directly. Used by code
// paths that are real engagement but leave no other durable trace (posting a
// standup, in a later phase). `day` is YYYY-MM-DD.
func (s *Store) MarkActive(ctx context.Context, userID int64, day, source string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO activity_days (user_id, day, sources) VALUES (?, date(?), ?)
		ON CONFLICT(user_id, day) DO UPDATE SET
			sources = CASE
				WHEN instr(',' || activity_days.sources || ',', ',' || excluded.sources || ',') > 0
					THEN activity_days.sources
				ELSE activity_days.sources || ',' || excluded.sources
			END`,
		userID, day, source)
	return err
}

// NormalizeSources keeps a sources list stable and de-duplicated for display.
func NormalizeSources(s string) string {
	seen := map[string]bool{}
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	sort.Strings(out)
	return strings.Join(out, ",")
}

// ExpireSanctions clears account suspensions whose expiry has passed, returning
// how many accounts were restored. Bans (which have no expiry) are untouched.
//
// AccountState already treats an elapsed suspension as active, so this is
// bookkeeping: it keeps stored state and effective state in agreement.
func (s *Store) ExpireSanctions(ctx context.Context) (int, error) {
	res, err := s.db.ExecContext(ctx, `
		UPDATE account_status
		SET state = 'active', reason = '', expires_at = NULL, updated_at = datetime('now')
		WHERE state = 'suspended'
		  AND expires_at IS NOT NULL
		  AND datetime(expires_at) <= datetime('now')`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

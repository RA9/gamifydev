package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
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

// ErrBadTransition is returned when a caller asks for a move the state machine
// does not allow (e.g. reviving a dropped enrollment).
var ErrBadTransition = errors.New("illegal enrollment transition")

// EnrollmentTerminal reports whether a state ends the enrollment's life.
func EnrollmentTerminal(state string) bool {
	return state == EnrollDropped || state == EnrollCompleted
}

// Enrollment is a learner's relationship to a path.
type Enrollment struct {
	ID        int64
	UserID    int64
	PathID    sql.NullInt64
	State     string
	Reason    string
	PlacedAt  sql.NullString
	StartedAt sql.NullString
	EndedAt   sql.NullString
}

const enrollCols = `id, user_id, path_id, state, reason, placed_at, started_at, ended_at`

func scanEnrollment(row interface{ Scan(...any) error }) (*Enrollment, error) {
	var e Enrollment
	err := row.Scan(&e.ID, &e.UserID, &e.PathID, &e.State, &e.Reason, &e.PlacedAt, &e.StartedAt, &e.EndedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// LiveEnrollment returns the user's current non-terminal enrollment, or nil if
// they have none (never registered for a path, or previously dropped).
func (s *Store) LiveEnrollment(ctx context.Context, userID int64) (*Enrollment, error) {
	e, err := scanEnrollment(s.db.QueryRowContext(ctx,
		`SELECT `+enrollCols+` FROM enrollments
		 WHERE user_id = ? AND state IN ('unplaced','placed','active','paused')`, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return e, err
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

// PlaceEnrollment records the outcome of the diagnostic: which path the learner
// is routed to. Moves unplaced -> placed.
func (s *Store) PlaceEnrollment(ctx context.Context, id, pathID int64, reason string) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE enrollments SET path_id = ? WHERE id = ?`, pathID, id); err != nil {
		return err
	}
	return s.TransitionEnrollment(ctx, id, EnrollPlaced, reason)
}

// --- Activity ledger ---------------------------------------------------------

// RollupActivity materializes activity_days from the signals the platform
// already records, for days on or after `since` (a YYYY-MM-DD string).
//
// Signals counted as "showed up and did something":
//   - completing a lesson step  (step_progress.completed_at)
//   - submitting an assignment  (submissions.created_at)
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
		)
		GROUP BY user_id, day
		ON CONFLICT(user_id, day) DO NOTHING`
	res, err := s.db.ExecContext(ctx, q, since, since)
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

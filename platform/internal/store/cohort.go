package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Cohorts, membership, and the async daily standup.

const (
	CohortActive    = "active"
	CohortCompleted = "completed"
	CohortArchived  = "archived"

	RoleLearner = "learner"
	RoleMentor  = "mentor"
)

// Cohort is a group running one path together.
type Cohort struct {
	ID        int64
	PathID    int64
	Name      string
	State     string
	TZBand    string
	StartsOn  string
	CreatedAt string
	// Populated by listing queries.
	ActiveMembers int
	PathTitle     string
	PathSlug      string
}

// Member is someone in a cohort.
type Member struct {
	UserID   int64
	Name     string
	Email    string
	Role     string
	IsBot    bool
	JoinedAt string
	// Populated where the caller asked for standup state.
	PostedToday bool
}

// IsMentor reports whether this member facilitates rather than learns.
func (m Member) IsMentor() bool { return m.Role == RoleMentor }

// CreateCohort opens a cohort.
func (s *Store) CreateCohort(ctx context.Context, c Cohort) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO cohorts (path_id, name, state, tz_band, starts_on)
		VALUES (?, ?, ?, ?, ?) RETURNING id`,
		c.PathID, c.Name, CohortActive, c.TZBand, c.StartsOn).Scan(&id)
	return id, err
}

// GetCohort loads one cohort with its path and active-member count.
func (s *Store) GetCohort(ctx context.Context, id int64) (*Cohort, error) {
	var c Cohort
	err := s.db.QueryRowContext(ctx, `
		SELECT c.id, c.path_id, c.name, c.state, c.tz_band, c.starts_on, c.created_at,
		       COALESCE(p.title,''), COALESCE(p.slug,''),
		       (SELECT COUNT(*) FROM cohort_members m
		          WHERE m.cohort_id = c.id AND m.left_at IS NULL AND m.role = 'learner')
		FROM cohorts c LEFT JOIN paths p ON p.id = c.path_id
		WHERE c.id = ?`, id).
		Scan(&c.ID, &c.PathID, &c.Name, &c.State, &c.TZBand, &c.StartsOn, &c.CreatedAt,
			&c.PathTitle, &c.PathSlug, &c.ActiveMembers)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &c, err
}

// CohortForUser returns the cohort a user is actively a member of, or nil.
func (s *Store) CohortForUser(ctx context.Context, userID int64) (*Cohort, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		SELECT m.cohort_id FROM cohort_members m
		JOIN cohorts c ON c.id = m.cohort_id
		WHERE m.user_id = ? AND m.left_at IS NULL AND c.state = 'active'
		ORDER BY m.joined_at DESC LIMIT 1`, userID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.GetCohort(ctx, id)
}

// AddMember puts a user in a cohort. Idempotent for an existing active member.
func (s *Store) AddMember(ctx context.Context, cohortID, userID int64, role string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cohort_members (cohort_id, user_id, role) VALUES (?, ?, ?)
		ON CONFLICT(cohort_id, user_id) WHERE left_at IS NULL DO NOTHING`,
		cohortID, userID, role)
	return err
}

// RemoveMember marks a membership ended without deleting it, so the member's
// standup history stays attributable.
func (s *Store) RemoveMember(ctx context.Context, cohortID, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE cohort_members SET left_at = datetime('now')
		 WHERE cohort_id = ? AND user_id = ? AND left_at IS NULL`, cohortID, userID)
	return err
}

// CohortMembers lists active members, mentors first.
func (s *Store) CohortMembers(ctx context.Context, cohortID int64) ([]Member, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT u.id, u.name, u.email, m.role, u.is_bot, m.joined_at
		FROM cohort_members m JOIN users u ON u.id = m.user_id
		WHERE m.cohort_id = ? AND m.left_at IS NULL
		ORDER BY CASE m.role WHEN 'mentor' THEN 0 ELSE 1 END, m.joined_at`, cohortID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Member
	for rows.Next() {
		var m Member
		var bot int
		if err := rows.Scan(&m.UserID, &m.Name, &m.Email, &m.Role, &bot, &m.JoinedAt); err != nil {
			return nil, err
		}
		m.IsBot = bot == 1
		out = append(out, m)
	}
	return out, rows.Err()
}

// --- Formation ---------------------------------------------------------------

// QueueGroup is a set of learners waiting for the same path and band.
type QueueGroup struct {
	PathID    int64
	PathTitle string
	TZBand    string
	Waiting   int
	OldestAge time.Duration
}

// PlacementQueue reports learners in state 'placed' with no cohort, grouped by
// the two things cohort formation matches on.
func (s *Store) PlacementQueue(ctx context.Context) ([]QueueGroup, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT e.path_id, COALESCE(p.title,''), e.tz_band, COUNT(*),
		       CAST((julianday('now') - julianday(MIN(e.placed_at))) * 86400 AS INTEGER)
		FROM enrollments e LEFT JOIN paths p ON p.id = e.path_id
		WHERE e.state = 'placed' AND e.cohort_id IS NULL AND e.path_id IS NOT NULL
		GROUP BY e.path_id, e.tz_band
		ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []QueueGroup
	for rows.Next() {
		var g QueueGroup
		var secs sql.NullInt64
		if err := rows.Scan(&g.PathID, &g.PathTitle, &g.TZBand, &g.Waiting, &secs); err != nil {
			return nil, err
		}
		if secs.Valid && secs.Int64 > 0 {
			g.OldestAge = time.Duration(secs.Int64) * time.Second
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// TakeFromQueue claims up to `limit` waiting enrollments for a path/band, moves
// them into `cohortID`, and activates them.
//
// Done in one transaction so a learner cannot be placed into two cohorts by
// concurrent formation runs.
func (s *Store) TakeFromQueue(ctx context.Context, cohortID, pathID int64, band string, limit int) ([]int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	rows, err := tx.QueryContext(ctx, `
		SELECT id, user_id FROM enrollments
		WHERE state = 'placed' AND cohort_id IS NULL AND path_id = ? AND tz_band = ?
		ORDER BY placed_at LIMIT ?`, pathID, band, limit)
	if err != nil {
		return nil, err
	}
	type row struct{ enrollID, userID int64 }
	var claimed []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.enrollID, &r.userID); err != nil {
			rows.Close()
			return nil, err
		}
		claimed = append(claimed, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var users []int64
	for _, r := range claimed {
		if _, err := tx.ExecContext(ctx, `
			UPDATE enrollments
			SET cohort_id = ?, state = 'active',
			    started_at = COALESCE(started_at, datetime('now')),
			    updated_at = datetime('now')
			WHERE id = ? AND state = 'placed' AND cohort_id IS NULL`, cohortID, r.enrollID); err != nil {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO cohort_members (cohort_id, user_id, role) VALUES (?, ?, 'learner')
			ON CONFLICT(cohort_id, user_id) WHERE left_at IS NULL DO NOTHING`,
			cohortID, r.userID); err != nil {
			return nil, err
		}
		users = append(users, r.userID)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return users, nil
}

// --- Bot mentor --------------------------------------------------------------

// EnsureBotMentor returns the id of the shared bot mentor account, creating it
// on first use. It has no usable password, so it can author posts but never
// sign in.
func (s *Store) EnsureBotMentor(ctx context.Context) (int64, error) {
	const email = "mentor@bot.gamifydev.local"
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT id FROM users WHERE email = ?`, email).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, name, role, is_bot, bio)
		VALUES (?, '', 'Pixel', 'learner', 1, 'Your cohort mentor. I run standup and flag blockers.')
		RETURNING id`, email).Scan(&id)
	return id, err
}

// --- Repack ------------------------------------------------------------------

// ThinCohort is an active cohort that has decayed below the viable size.
type ThinCohort struct {
	ID      int64
	PathID  int64
	TZBand  string
	Members int
}

// ThinCohorts lists active cohorts with fewer than `min` active learners.
func (s *Store) ThinCohorts(ctx context.Context, min int) ([]ThinCohort, error) {
	// The member count is a correlated subquery, not an aggregate, so it must be
	// filtered in an outer SELECT rather than with HAVING.
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, path_id, tz_band, n FROM (
			SELECT c.id, c.path_id, c.tz_band,
			       (SELECT COUNT(*) FROM cohort_members m
			          WHERE m.cohort_id = c.id AND m.left_at IS NULL AND m.role = 'learner') AS n
			FROM cohorts c WHERE c.state = 'active'
		)
		WHERE n < ? ORDER BY n`, min)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ThinCohort
	for rows.Next() {
		var t ThinCohort
		if err := rows.Scan(&t.ID, &t.PathID, &t.TZBand, &t.Members); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// MergeTarget finds another active cohort on the same path and band with room
// for `incoming` more learners. Returns 0 when there is none.
func (s *Store) MergeTarget(ctx context.Context, fromCohortID, pathID int64, band string, incoming, maxSize int) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
		SELECT c.id FROM cohorts c
		WHERE c.state = 'active' AND c.id != ? AND c.path_id = ? AND c.tz_band = ?
		  AND (SELECT COUNT(*) FROM cohort_members m
		         WHERE m.cohort_id = c.id AND m.left_at IS NULL AND m.role = 'learner') + ? <= ?
		ORDER BY (SELECT COUNT(*) FROM cohort_members m
		            WHERE m.cohort_id = c.id AND m.left_at IS NULL AND m.role = 'learner') DESC
		LIMIT 1`, fromCohortID, pathID, band, incoming, maxSize).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return id, err
}

// MergeCohorts moves every active learner from `from` into `into`, then archives
// the emptied cohort. Mentors are not moved; the target has its own.
func (s *Store) MergeCohorts(ctx context.Context, from, into int64) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck

	rows, err := tx.QueryContext(ctx, `
		SELECT user_id FROM cohort_members
		WHERE cohort_id = ? AND left_at IS NULL AND role = 'learner'`, from)
	if err != nil {
		return 0, err
	}
	var users []int64
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			rows.Close()
			return 0, err
		}
		users = append(users, uid)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	for _, uid := range users {
		if _, err := tx.ExecContext(ctx,
			`UPDATE cohort_members SET left_at = datetime('now')
			 WHERE cohort_id = ? AND user_id = ? AND left_at IS NULL`, from, uid); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO cohort_members (cohort_id, user_id, role) VALUES (?, ?, 'learner')
			ON CONFLICT(cohort_id, user_id) WHERE left_at IS NULL DO NOTHING`, into, uid); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE enrollments SET cohort_id = ?, updated_at = datetime('now')
			 WHERE user_id = ? AND cohort_id = ?`, into, uid, from); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE cohorts SET state = 'archived' WHERE id = ?`, from); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return len(users), nil
}

// --- Standups ----------------------------------------------------------------

// Standup is one cohort-day.
type Standup struct {
	ID       int64
	CohortID int64
	Day      string
	OpensAt  string
	ClosesAt string
	Prompt   string
	Closed   bool
	IsOpen   bool // computed against the database clock
}

// StandupEntry is one member's post.
type StandupEntry struct {
	ID        int64
	StandupID int64
	UserID    int64
	UserName  string
	IsBot     bool
	Yesterday string
	Today     string
	Blockers  string
	PostedAt  string
	Late      bool
}

// EnsureStandup creates a cohort's standup for a local day if it doesn't exist.
// Idempotent — the unique index makes a concurrent second call a no-op.
func (s *Store) EnsureStandup(ctx context.Context, cohortID int64, day, opensAt, closesAt, prompt string) (int64, bool, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO standups (cohort_id, day, opens_at, closes_at, prompt)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(cohort_id, day) DO NOTHING`,
		cohortID, day, opensAt, closesAt, prompt)
	if err != nil {
		return 0, false, err
	}
	n, _ := res.RowsAffected()
	var id int64
	if err := s.db.QueryRowContext(ctx,
		`SELECT id FROM standups WHERE cohort_id = ? AND day = ?`, cohortID, day).Scan(&id); err != nil {
		return 0, false, err
	}
	return id, n > 0, nil
}

const standupCols = `id, cohort_id, day, opens_at, closes_at, prompt, closed,
	CASE WHEN closed = 0 AND datetime('now') BETWEEN datetime(opens_at) AND datetime(closes_at)
	     THEN 1 ELSE 0 END`

func scanStandup(row interface{ Scan(...any) error }) (*Standup, error) {
	var s Standup
	var closed, open int
	if err := row.Scan(&s.ID, &s.CohortID, &s.Day, &s.OpensAt, &s.ClosesAt, &s.Prompt, &closed, &open); err != nil {
		return nil, err
	}
	s.Closed, s.IsOpen = closed == 1, open == 1
	return &s, nil
}

// TodayStandup returns a cohort's standup for the given local day, or nil.
func (s *Store) TodayStandup(ctx context.Context, cohortID int64, day string) (*Standup, error) {
	st, err := scanStandup(s.db.QueryRowContext(ctx,
		`SELECT `+standupCols+` FROM standups WHERE cohort_id = ? AND day = ?`, cohortID, day))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return st, err
}

// RecentStandups returns a cohort's most recent standups, newest first.
func (s *Store) RecentStandups(ctx context.Context, cohortID int64, limit int) ([]Standup, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+standupCols+` FROM standups WHERE cohort_id = ? ORDER BY day DESC LIMIT ?`,
		cohortID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Standup
	for rows.Next() {
		st, err := scanStandup(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *st)
	}
	return out, rows.Err()
}

// StandupEntries returns the posts for a standup, oldest first.
func (s *Store) StandupEntries(ctx context.Context, standupID int64) ([]StandupEntry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT e.id, e.standup_id, e.user_id, u.name, u.is_bot,
		       e.yesterday, e.today, e.blockers, e.posted_at, e.late
		FROM standup_entries e JOIN users u ON u.id = e.user_id
		WHERE e.standup_id = ? ORDER BY e.posted_at`, standupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StandupEntry
	for rows.Next() {
		var e StandupEntry
		var bot, late int
		if err := rows.Scan(&e.ID, &e.StandupID, &e.UserID, &e.UserName, &bot,
			&e.Yesterday, &e.Today, &e.Blockers, &e.PostedAt, &late); err != nil {
			return nil, err
		}
		e.IsBot, e.Late = bot == 1, late == 1
		out = append(out, e)
	}
	return out, rows.Err()
}

// ErrStandupClosed is returned when a post arrives after the window shut.
var ErrStandupClosed = errors.New("standup is closed")

// PostStandup records a member's entry. Re-posting on the same day updates the
// existing entry rather than creating a second one.
//
// Posting also records the day in the activity ledger: showing up for standup is
// participation, and phase 5's attendance window reads that ledger.
func (s *Store) PostStandup(ctx context.Context, standupID, userID int64, yesterday, today, blockers string) error {
	var closesAt string
	var closed int
	if err := s.db.QueryRowContext(ctx,
		`SELECT closes_at, closed FROM standups WHERE id = ?`, standupID).Scan(&closesAt, &closed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if closed == 1 {
		return ErrStandupClosed
	}
	// A post inside the window is on time; the column exists so phase 5 can tell
	// "late but present" from "absent" without re-deriving it.
	var late int
	if err := s.db.QueryRowContext(ctx,
		`SELECT CASE WHEN datetime('now') > datetime(?) THEN 1 ELSE 0 END`, closesAt).Scan(&late); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO standup_entries (standup_id, user_id, yesterday, today, blockers, late)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(standup_id, user_id) DO UPDATE SET
			yesterday = excluded.yesterday, today = excluded.today,
			blockers = excluded.blockers, posted_at = datetime('now')`,
		standupID, userID, yesterday, today, blockers, late); err != nil {
		return err
	}
	day := time.Now().UTC().Format("2006-01-02")
	return s.MarkActive(ctx, userID, day, "standup")
}

// HasPosted reports whether a user has an entry for a standup.
func (s *Store) HasPosted(ctx context.Context, standupID, userID int64) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM standup_entries WHERE standup_id = ? AND user_id = ?`,
		standupID, userID).Scan(&n)
	return n > 0, err
}

// OpenCohorts lists active cohorts, for the jobs that fan out over them.
func (s *Store) OpenCohorts(ctx context.Context) ([]Cohort, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.path_id, c.name, c.state, c.tz_band, c.starts_on, c.created_at,
		       COALESCE(p.title,''), COALESCE(p.slug,'')
		FROM cohorts c LEFT JOIN paths p ON p.id = c.path_id
		WHERE c.state = 'active' ORDER BY c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Cohort
	for rows.Next() {
		var c Cohort
		if err := rows.Scan(&c.ID, &c.PathID, &c.Name, &c.State, &c.TZBand, &c.StartsOn,
			&c.CreatedAt, &c.PathTitle, &c.PathSlug); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CloseExpiredStandups shuts windows whose close time has passed.
func (s *Store) CloseExpiredStandups(ctx context.Context) (int, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE standups SET closed = 1
		 WHERE closed = 0 AND datetime(closes_at) <= datetime('now')`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

// SetEnrollmentBand records the timezone band a learner chose at enrollment.
func (s *Store) SetEnrollmentBand(ctx context.Context, enrollmentID int64, band string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE enrollments SET tz_band = ?, updated_at = datetime('now') WHERE id = ?`,
		band, enrollmentID)
	return err
}

// CohortStats summarises the program for the admin view.
type CohortStats struct {
	Active, Archived, Learners, Queued int
}

// CohortOverview aggregates cohort counts.
func (s *Store) CohortOverview(ctx context.Context) (CohortStats, error) {
	var st CohortStats
	err := s.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM cohorts WHERE state='active'),
		  (SELECT COUNT(*) FROM cohorts WHERE state='archived'),
		  (SELECT COUNT(*) FROM cohort_members m JOIN cohorts c ON c.id=m.cohort_id
		     WHERE m.left_at IS NULL AND m.role='learner' AND c.state='active'),
		  (SELECT COUNT(*) FROM enrollments WHERE state='placed' AND cohort_id IS NULL)`).
		Scan(&st.Active, &st.Archived, &st.Learners, &st.Queued)
	return st, err
}

// ListCohorts returns cohorts for the admin table.
func (s *Store) ListCohorts(ctx context.Context, limit int) ([]Cohort, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT c.id, c.path_id, c.name, c.state, c.tz_band, c.starts_on, c.created_at,
		       COALESCE(p.title,''), COALESCE(p.slug,''),
		       (SELECT COUNT(*) FROM cohort_members m
		          WHERE m.cohort_id = c.id AND m.left_at IS NULL AND m.role='learner')
		FROM cohorts c LEFT JOIN paths p ON p.id = c.path_id
		ORDER BY CASE c.state WHEN 'active' THEN 0 ELSE 1 END, c.id DESC LIMIT %d`, limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Cohort
	for rows.Next() {
		var c Cohort
		if err := rows.Scan(&c.ID, &c.PathID, &c.Name, &c.State, &c.TZBand, &c.StartsOn,
			&c.CreatedAt, &c.PathTitle, &c.PathSlug, &c.ActiveMembers); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

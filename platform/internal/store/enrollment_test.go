package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

// newTestStore opens a migrated store backed by a throwaway file. A file (rather
// than :memory:) is used because the store opens a pool and in-memory SQLite
// gives each connection its own private database.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return st
}

func mkUser(t *testing.T, st *Store, email string) int64 {
	t.Helper()
	var id int64
	err := st.db.QueryRow(
		`INSERT INTO users (email, name) VALUES (?, ?) RETURNING id`, email, "Test").Scan(&id)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return id
}

func TestEnrollmentLifecycle(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "learner@example.com")

	e, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("ensure: %v", err)
	}
	if e.State != EnrollUnplaced {
		t.Fatalf("new enrollment state = %q, want %q", e.State, EnrollUnplaced)
	}

	// Idempotent: a second call must not create a second live enrollment. The
	// partial unique index is what guarantees this.
	again, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("ensure twice: %v", err)
	}
	if again.ID != e.ID {
		t.Fatalf("EnsureEnrollment created a second enrollment (%d then %d)", e.ID, again.ID)
	}

	// unplaced -> active is not a legal move; only unplaced -> placed is.
	if err := st.TransitionEnrollment(ctx, e.ID, EnrollActive, ""); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("unplaced->active err = %v, want ErrBadTransition", err)
	}

	var pathID int64
	if err := st.db.QueryRow(
		`INSERT INTO paths (slug, title, published) VALUES ('p','P',1) RETURNING id`).Scan(&pathID); err != nil {
		t.Fatalf("create path: %v", err)
	}
	if err := st.PlaceEnrollment(ctx, e.ID, pathID, "diagnostic"); err != nil {
		t.Fatalf("place: %v", err)
	}
	got, _ := st.LiveEnrollment(ctx, uid)
	if got.State != EnrollPlaced || !got.PlacedAt.Valid {
		t.Fatalf("after place: state=%q placed_at.valid=%v", got.State, got.PlacedAt.Valid)
	}
	if !got.PathID.Valid || got.PathID.Int64 != pathID {
		t.Fatalf("path not recorded: %+v", got.PathID)
	}

	// placed -> active -> paused -> active is the deferral round trip.
	for _, to := range []string{EnrollActive, EnrollPaused, EnrollActive} {
		if err := st.TransitionEnrollment(ctx, e.ID, to, ""); err != nil {
			t.Fatalf("transition to %s: %v", to, err)
		}
	}
	got, _ = st.LiveEnrollment(ctx, uid)
	if !got.StartedAt.Valid {
		t.Fatal("started_at not stamped on active")
	}
	if got.EndedAt.Valid {
		t.Fatal("ended_at should be cleared when returning to active")
	}

	// Terminal, and terminal really is terminal.
	if err := st.TransitionEnrollment(ctx, e.ID, EnrollDropped, "went dark"); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := st.TransitionEnrollment(ctx, e.ID, EnrollActive, ""); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("dropped->active err = %v, want ErrBadTransition", err)
	}

	// A dropped enrollment is not "live", so the learner can enroll again —
	// that's the reapplication path, and it must produce a new row.
	if live, _ := st.LiveEnrollment(ctx, uid); live != nil {
		t.Fatalf("dropped enrollment still reads as live: %+v", live)
	}
	fresh, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("re-enroll: %v", err)
	}
	if fresh.ID == e.ID {
		t.Fatal("re-enrollment reused the dropped row instead of creating a new one")
	}
}

func TestAccountStateExpiry(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "suspended@example.com")

	// No row at all must read as active, so existing users need no backfill.
	if got, _ := st.AccountState(ctx, uid); got != AccountActive {
		t.Fatalf("default state = %q, want active", got)
	}

	past := time.Now().Add(-time.Hour).UTC().Format("2006-01-02 15:04:05")
	future := time.Now().Add(24 * time.Hour).UTC().Format("2006-01-02 15:04:05")

	if err := st.SetAccountState(ctx, uid, AccountSuspended, "test", &future); err != nil {
		t.Fatalf("suspend: %v", err)
	}
	if got, _ := st.AccountState(ctx, uid); got != AccountSuspended {
		t.Fatalf("state = %q, want suspended", got)
	}

	// An elapsed suspension must read as active even before the expiry job runs —
	// the read path cannot depend on a job having fired.
	if err := st.SetAccountState(ctx, uid, AccountSuspended, "test", &past); err != nil {
		t.Fatalf("re-suspend: %v", err)
	}
	if got, _ := st.AccountState(ctx, uid); got != AccountActive {
		t.Fatalf("elapsed suspension reads as %q, want active", got)
	}

	// ExpireSanctions then makes stored state agree with effective state.
	n, err := st.ExpireSanctions(ctx)
	if err != nil || n != 1 {
		t.Fatalf("ExpireSanctions = %d, %v; want 1, nil", n, err)
	}

	// A ban has no expiry and must survive the sweep.
	if err := st.SetAccountState(ctx, uid, AccountBanned, "conduct", nil); err != nil {
		t.Fatalf("ban: %v", err)
	}
	if _, err := st.ExpireSanctions(ctx); err != nil {
		t.Fatalf("expire: %v", err)
	}
	if got, _ := st.AccountState(ctx, uid); got != AccountBanned {
		t.Fatalf("ban was cleared by the expiry sweep, got %q", got)
	}
}

func TestRollupActivityIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "active@example.com")

	// Two step completions on the same day plus a submission on another day.
	mustExec(t, st, `INSERT INTO lessons (course_id, slug, title) VALUES (1,'l','L')`)
	mustExec(t, st, `INSERT INTO lesson_steps (id, lesson_id, sort) VALUES (1, 1, 0), (2, 1, 1)`)
	mustExec(t, st, `INSERT INTO step_progress (user_id, step_id, completed_at) VALUES
		(?, 1, '2026-01-10 09:00:00'), (?, 2, '2026-01-10 17:00:00')`, uid, uid)
	mustExec(t, st, `INSERT INTO assignments (slug, title) VALUES ('a','A')`)
	mustExec(t, st, `INSERT INTO submissions (assignment_id, user_id, created_at)
		VALUES (1, ?, '2026-01-12 12:00:00')`, uid)

	n, err := st.RollupActivity(ctx, "2026-01-01")
	if err != nil {
		t.Fatalf("rollup: %v", err)
	}
	if n != 2 {
		t.Fatalf("first rollup inserted %d days, want 2 (Jan 10 and Jan 12)", n)
	}

	// Re-running over the same window must insert nothing — the job overlaps its
	// window deliberately, so this property is load-bearing.
	n, err = st.RollupActivity(ctx, "2026-01-01")
	if err != nil {
		t.Fatalf("rollup twice: %v", err)
	}
	if n != 0 {
		t.Fatalf("second rollup inserted %d rows, want 0 (not idempotent)", n)
	}

	days, err := st.ActiveDaysBetween(ctx, uid, "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("active days: %v", err)
	}
	if len(days) != 2 || days[0] != "2026-01-10" || days[1] != "2026-01-12" {
		t.Fatalf("active days = %v, want [2026-01-10 2026-01-12]", days)
	}

	// MarkActive records a day that leaves no other trace, and must not duplicate
	// a source that is already listed.
	if err := st.MarkActive(ctx, uid, "2026-01-10", "standup"); err != nil {
		t.Fatalf("mark active: %v", err)
	}
	if err := st.MarkActive(ctx, uid, "2026-01-10", "standup"); err != nil {
		t.Fatalf("mark active twice: %v", err)
	}
	var sources string
	if err := st.db.QueryRow(
		`SELECT sources FROM activity_days WHERE user_id = ? AND day = '2026-01-10'`, uid).Scan(&sources); err != nil {
		t.Fatalf("read sources: %v", err)
	}
	if got := NormalizeSources(sources); got != "standup,step" {
		t.Fatalf("sources = %q (raw %q), want %q", got, sources, "standup,step")
	}
}

func mustExec(t *testing.T, st *Store, q string, args ...any) {
	t.Helper()
	if _, err := st.db.Exec(q, args...); err != nil {
		t.Fatalf("exec %s: %v", q, err)
	}
}

package store

import (
	"context"
	"database/sql"
	"encoding/json"
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

func placementResultFixture(t *testing.T, st *Store, userID int64, pathSlug string, published, passed bool, exemptions []string) (int64, int64) {
	t.Helper()
	ctx := context.Background()
	assessmentID, err := st.UpsertAssessment(ctx, Assessment{
		Slug: "placement", Title: "Placement", Kind: "placement", TimeLimitS: 1800, PerTopic: 1,
	})
	if err != nil {
		t.Fatalf("upsert placement assessment: %v", err)
	}
	var pathID int64
	if err := st.db.QueryRow(
		`INSERT INTO paths (slug, title, published) VALUES (?, ?, ?) RETURNING id`,
		pathSlug, pathSlug, boolToInt(published)).Scan(&pathID); err != nil {
		t.Fatalf("create path: %v", err)
	}
	var attemptID int64
	if err := st.db.QueryRow(`
		INSERT INTO assessment_attempts
			(user_id, assessment_id, expires_at, submitted_at, score, topic_scores)
		VALUES (?, ?, datetime('now', '+1 hour'), datetime('now'), 100, '{}')
		RETURNING id`, userID, assessmentID).Scan(&attemptID); err != nil {
		t.Fatalf("create submitted placement attempt: %v", err)
	}
	raw, err := json.Marshal(exemptions)
	if err != nil {
		t.Fatalf("encode exemptions: %v", err)
	}
	var resultID int64
	if err := st.db.QueryRow(`
		INSERT INTO placement_results
			(attempt_id, user_id, recommended_path_id, foundations_required, passed, exemptions)
		VALUES (?, ?, ?, 0, ?, ?) RETURNING id`,
		attemptID, userID, pathID, boolToInt(passed), string(raw)).Scan(&resultID); err != nil {
		t.Fatalf("create placement result: %v", err)
	}
	return pathID, resultID
}

func TestEnrollFromPlacementBindsExactGeneratedSnapshot(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "generated@example.com")
	original, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("ensure enrollment: %v", err)
	}
	pathID, resultID := placementResultFixture(t, st, uid, "generated-path", true, true, []string{"c", "data-structures"})

	enrollment, err := st.EnrollFromPlacement(ctx, uid, resultID, "americas")
	if err != nil {
		t.Fatalf("enroll from placement: %v", err)
	}
	if enrollment.ID != original.ID {
		t.Fatalf("enrollment id = %d, want existing unplaced row %d", enrollment.ID, original.ID)
	}
	if enrollment.State != EnrollPlaced || !enrollment.PlacedAt.Valid {
		t.Fatalf("state = %q, placed_at = %+v; want placed with timestamp", enrollment.State, enrollment.PlacedAt)
	}
	if !enrollment.PathID.Valid || enrollment.PathID.Int64 != pathID {
		t.Fatalf("path = %+v, want %d", enrollment.PathID, pathID)
	}
	if !enrollment.PlacementResultID.Valid || enrollment.PlacementResultID.Int64 != resultID {
		t.Fatalf("placement result = %+v, want %d", enrollment.PlacementResultID, resultID)
	}
	if enrollment.TZBand != "americas" {
		t.Fatalf("timezone band = %q, want americas", enrollment.TZBand)
	}
	if len(enrollment.Exemptions) != 2 || enrollment.Exemptions[0] != "c" || enrollment.Exemptions[1] != "data-structures" {
		t.Fatalf("exemptions = %v, want exact placement snapshot", enrollment.Exemptions)
	}

	// Idempotency depends on the immutable enrollment binding, not mutable path
	// publication state.
	mustExec(t, st, `UPDATE paths SET published = 0 WHERE id = ?`, pathID)
	again, err := st.EnrollFromPlacement(ctx, uid, resultID, "americas")
	if err != nil {
		t.Fatalf("repeat enrollment: %v", err)
	}
	if again.ID != enrollment.ID || again.PlacementResultID.Int64 != resultID {
		t.Fatalf("repeat enrollment changed binding: before=%+v after=%+v", enrollment, again)
	}

	_, newerResultID := placementResultFixture(t, st, uid, "newer-path", true, true, []string{"algorithms"})
	if _, err := st.EnrollFromPlacement(ctx, uid, newerResultID, "americas"); !errors.Is(err, ErrEnrollmentLocked) {
		t.Fatalf("rebind err = %v, want ErrEnrollmentLocked", err)
	}
	locked, err := st.LiveEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("reload locked enrollment: %v", err)
	}
	if locked.PathID.Int64 != pathID || locked.PlacementResultID.Int64 != resultID || locked.TZBand != "americas" || len(locked.Exemptions) != 2 {
		t.Fatalf("locked enrollment mutated: %+v", locked)
	}
}

func TestConcurrentDuplicateEnrollmentIsIdempotent(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "concurrent@example.com")
	if _, err := st.EnsureEnrollment(ctx, uid); err != nil {
		t.Fatalf("ensure enrollment: %v", err)
	}
	_, resultID := placementResultFixture(t, st, uid, "concurrent-path", true, true, []string{"c"})

	start := make(chan struct{})
	type outcome struct {
		enrollment *Enrollment
		err        error
	}
	outcomes := make(chan outcome, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			e, err := st.EnrollFromPlacement(ctx, uid, resultID, "europe_africa")
			outcomes <- outcome{enrollment: e, err: err}
		}()
	}
	close(start)
	first, second := <-outcomes, <-outcomes
	for i, got := range []outcome{first, second} {
		if got.err != nil || got.enrollment == nil {
			t.Fatalf("concurrent enrollment %d = %+v, %v", i+1, got.enrollment, got.err)
		}
	}
	if first.enrollment.ID != second.enrollment.ID || first.enrollment.PlacementResultID.Int64 != resultID || second.enrollment.PlacementResultID.Int64 != resultID {
		t.Fatalf("concurrent enrollment produced different bindings: %+v / %+v", first.enrollment, second.enrollment)
	}
}

func TestLegacyPlacementCannotRewriteGeneratedBinding(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "immutable@example.com")
	pathID, resultID := placementResultFixture(t, st, uid, "immutable-path", true, true, []string{"c"})
	enrollment, err := st.EnrollFromPlacement(ctx, uid, resultID, "europe_africa")
	if err != nil {
		t.Fatalf("generated enrollment: %v", err)
	}
	if err := st.TransitionEnrollment(ctx, enrollment.ID, EnrollActive, "test"); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if err := st.TransitionEnrollment(ctx, enrollment.ID, EnrollPaused, "test"); err != nil {
		t.Fatalf("pause: %v", err)
	}
	otherPath := mkPath(t, st, "operator-rewrite")
	if err := st.PlaceEnrollment(ctx, enrollment.ID, otherPath, "operator"); !errors.Is(err, ErrEnrollmentLocked) {
		t.Fatalf("operator rewrite err = %v, want ErrEnrollmentLocked", err)
	}
	got, err := st.LiveEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("reload enrollment: %v", err)
	}
	if got.State != EnrollPaused || got.PathID.Int64 != pathID || got.PlacementResultID.Int64 != resultID || len(got.Exemptions) != 1 || got.Exemptions[0] != "c" {
		t.Fatalf("generated binding was rewritten: %+v", got)
	}
}

func TestGeneratedEnrollmentMigrationPreservesLegacyExemptions(t *testing.T) {
	ctx := context.Background()
	st, err := Open(filepath.Join(t.TempDir(), "upgrade.db"))
	if err != nil {
		t.Fatalf("open upgrade store: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	mustExec(t, st, `CREATE TABLE placement_results (
		id INTEGER PRIMARY KEY, user_id INTEGER, passed INTEGER,
		recommended_path_id INTEGER, exemptions TEXT, created_at TEXT)`)
	mustExec(t, st, `CREATE TABLE enrollments (
		id INTEGER PRIMARY KEY, user_id INTEGER, path_id INTEGER, state TEXT,
		placed_at TEXT, ended_at TEXT, cohort_id INTEGER)`)
	mustExec(t, st, `CREATE TABLE sanctions (
		id INTEGER PRIMARY KEY, user_id INTEGER, cohort_id INTEGER, kind TEXT,
		shadow INTEGER)`)
	mustExec(t, st, `INSERT INTO placement_results
		(id, user_id, passed, recommended_path_id, exemptions, created_at)
		VALUES (1, 42, 1, 10, '["c"]', '2026-01-01 00:00:00')`)
	mustExec(t, st, `INSERT INTO enrollments (id, user_id, path_id, placed_at)
		VALUES (1, 42, 99, '2026-01-02 00:00:00')`)
	mustExec(t, st, `INSERT INTO enrollments
		(id, user_id, path_id, state, placed_at, ended_at, cohort_id) VALUES
		(2, 7, 10, 'dropped', '2026-01-01 00:00:00', '2026-02-01 10:00:02', 5),
		(3, 7, 11, 'dropped', '2026-03-01 00:00:00', '2026-04-01 10:00:02', 5)`)
	mustExec(t, st, `INSERT INTO sanctions
		(id, user_id, cohort_id, kind, shadow) VALUES
		(1, 7, 5, 'drop', 0), (2, 7, 5, 'drop', 0)`)
	mustExec(t, st, `ALTER TABLE sanctions ADD COLUMN applied_at TEXT`)
	mustExec(t, st, `UPDATE sanctions SET applied_at = CASE id
		WHEN 1 THEN '2026-02-01 10:00:00' ELSE '2026-04-01 10:00:00' END`)
	migration, err := migrationsFS.ReadFile("migrations/0027_generated_enrollment.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, string(migration)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}
	var boundResultID sql.NullInt64
	var exemptions string
	if err := st.db.QueryRow(`SELECT placement_result_id, placement_exemptions FROM enrollments WHERE id = 1`).
		Scan(&boundResultID, &exemptions); err != nil {
		t.Fatalf("read migrated enrollment: %v", err)
	}
	if boundResultID.Valid {
		t.Fatalf("migration falsely bound mismatched result %d", boundResultID.Int64)
	}
	if exemptions != `["c"]` {
		t.Fatalf("legacy exemptions = %q, want [\"c\"]", exemptions)
	}
	rows, err := st.db.Query(`SELECT id, enrollment_id FROM sanctions ORDER BY id`)
	if err != nil {
		t.Fatalf("read migrated sanctions: %v", err)
	}
	defer rows.Close()
	want := map[int64]int64{1: 2, 2: 3}
	for rows.Next() {
		var sanctionID, enrollmentID int64
		if err := rows.Scan(&sanctionID, &enrollmentID); err != nil {
			t.Fatalf("scan migrated sanction: %v", err)
		}
		if enrollmentID != want[sanctionID] {
			t.Fatalf("sanction %d enrollment = %d, want %d", sanctionID, enrollmentID, want[sanctionID])
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("migrated sanctions: %v", err)
	}
}

func TestEnrollFromPlacementRejectsIneligibleResultsWithoutMutation(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T, *Store, int64) int64
	}{
		{
			name: "failed",
			setup: func(t *testing.T, st *Store, uid int64) int64 {
				_, resultID := placementResultFixture(t, st, uid, "failed-path", true, false, nil)
				return resultID
			},
		},
		{
			name: "stale",
			setup: func(t *testing.T, st *Store, uid int64) int64 {
				_, stale := placementResultFixture(t, st, uid, "stale-path", true, true, nil)
				placementResultFixture(t, st, uid, "current-path", true, true, nil)
				return stale
			},
		},
		{
			name: "foreign",
			setup: func(t *testing.T, st *Store, _ int64) int64 {
				owner := mkUser(t, st, "owner@example.com")
				_, resultID := placementResultFixture(t, st, owner, "foreign-path", true, true, nil)
				return resultID
			},
		},
		{
			name: "unpublished path",
			setup: func(t *testing.T, st *Store, uid int64) int64 {
				_, resultID := placementResultFixture(t, st, uid, "hidden-path", false, true, nil)
				return resultID
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			st := newTestStore(t)
			uid := mkUser(t, st, "learner@example.com")
			original, err := st.EnsureEnrollment(ctx, uid)
			if err != nil {
				t.Fatalf("ensure enrollment: %v", err)
			}
			resultID := tc.setup(t, st, uid)

			if _, err := st.EnrollFromPlacement(ctx, uid, resultID, "asia_pacific"); !errors.Is(err, ErrEnrollmentPlacementInvalid) {
				t.Fatalf("enroll err = %v, want ErrEnrollmentPlacementInvalid", err)
			}
			got, err := st.LiveEnrollment(ctx, uid)
			if err != nil {
				t.Fatalf("reload enrollment: %v", err)
			}
			if got.ID != original.ID || got.State != EnrollUnplaced || got.PathID.Valid || got.PlacementResultID.Valid || got.TZBand != "" || len(got.Exemptions) != 0 {
				t.Fatalf("invalid placement partially mutated enrollment: %+v", got)
			}
		})
	}
}

func TestInvalidPlacementDoesNotCreateEnrollment(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "no-enrollment@example.com")
	_, resultID := placementResultFixture(t, st, uid, "failed-without-enrollment", true, false, nil)

	if _, err := st.EnrollFromPlacement(ctx, uid, resultID, "europe_africa"); !errors.Is(err, ErrEnrollmentPlacementInvalid) {
		t.Fatalf("enroll err = %v, want ErrEnrollmentPlacementInvalid", err)
	}
	if enrollment, err := st.LiveEnrollment(ctx, uid); err != nil || enrollment != nil {
		t.Fatalf("invalid placement created enrollment: %+v, %v", enrollment, err)
	}
}

func TestEnrollFromPlacementRollsBackMalformedSnapshot(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "malformed@example.com")
	original, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("ensure enrollment: %v", err)
	}
	_, resultID := placementResultFixture(t, st, uid, "malformed-path", true, true, nil)
	mustExec(t, st, `UPDATE placement_results SET exemptions = '{' WHERE id = ?`, resultID)

	if _, err := st.EnrollFromPlacement(ctx, uid, resultID, "europe_africa"); err == nil {
		t.Fatal("malformed exemption snapshot unexpectedly enrolled")
	}
	got, err := st.LiveEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("reload enrollment: %v", err)
	}
	if got.ID != original.ID || got.State != EnrollUnplaced || got.PathID.Valid || got.PlacementResultID.Valid || got.TZBand != "" {
		t.Fatalf("failed transaction partially mutated enrollment: %+v", got)
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

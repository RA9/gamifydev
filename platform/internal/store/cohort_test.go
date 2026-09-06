package store

import (
	"context"
	"errors"
	"testing"
)

// enroll creates a user placed on a path in a band, ready for cohort formation.
func enroll(t *testing.T, st *Store, email string, pathID int64, band string) int64 {
	t.Helper()
	ctx := context.Background()
	uid := mkUser(t, st, email)
	e, err := st.EnsureEnrollment(ctx, uid)
	if err != nil {
		t.Fatalf("enroll: %v", err)
	}
	if err := st.PlaceEnrollment(ctx, e.ID, pathID, "test"); err != nil {
		t.Fatalf("place: %v", err)
	}
	if err := st.SetEnrollmentBand(ctx, e.ID, band); err != nil {
		t.Fatalf("band: %v", err)
	}
	return uid
}

func mkPath(t *testing.T, st *Store, slug string) int64 {
	t.Helper()
	var id int64
	if err := st.db.QueryRow(
		`INSERT INTO paths (slug, title, published) VALUES (?, ?, 1) RETURNING id`,
		slug, slug).Scan(&id); err != nil {
		t.Fatalf("create path: %v", err)
	}
	return id
}

func TestPlacementQueueGroupsByPathAndBand(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	p1, p2 := mkPath(t, st, "frontend"), mkPath(t, st, "backend")

	enroll(t, st, "a@x.com", p1, "europe_africa")
	enroll(t, st, "b@x.com", p1, "europe_africa")
	enroll(t, st, "c@x.com", p1, "asia_pacific")
	enroll(t, st, "d@x.com", p2, "europe_africa")

	groups, err := st.PlacementQueue(ctx)
	if err != nil {
		t.Fatalf("queue: %v", err)
	}
	// Three distinct (path, band) buckets — grouping on both is what keeps a
	// cohort's standup window workable.
	if len(groups) != 3 {
		t.Fatalf("got %d groups, want 3: %+v", len(groups), groups)
	}
	for _, g := range groups {
		want := 1
		if g.PathID == p1 && g.TZBand == "europe_africa" {
			want = 2
		}
		if g.Waiting != want {
			t.Fatalf("group %+v: waiting = %d, want %d", g, g.Waiting, want)
		}
	}
}

func TestTakeFromQueueActivatesAndAssigns(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	u1 := enroll(t, st, "a@x.com", pid, "europe_africa")
	u2 := enroll(t, st, "b@x.com", pid, "europe_africa")

	cid, err := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "C1", TZBand: "europe_africa", StartsOn: "2026-09-03"})
	if err != nil {
		t.Fatalf("create cohort: %v", err)
	}
	users, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 10)
	if err != nil {
		t.Fatalf("take: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("claimed %d learners, want 2", len(users))
	}

	// Each learner is now active, in the cohort, and off the queue.
	for _, uid := range []int64{u1, u2} {
		e, _ := st.LiveEnrollment(ctx, uid)
		if e.State != EnrollActive {
			t.Fatalf("user %d state = %q, want active", uid, e.State)
		}
		c, _ := st.CohortForUser(ctx, uid)
		if c == nil || c.ID != cid {
			t.Fatalf("user %d not in the cohort", uid)
		}
	}
	groups, _ := st.PlacementQueue(ctx)
	if len(groups) != 0 {
		t.Fatalf("queue still has %d group(s) after placement", len(groups))
	}

	// A second run must claim nobody — the guard against double-placing a
	// learner when two formation runs overlap.
	again, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 10)
	if err != nil {
		t.Fatalf("second take: %v", err)
	}
	if len(again) != 0 {
		t.Fatalf("second take claimed %d learners, want 0", len(again))
	}
}

func TestTakeFromQueueRespectsLimit(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	for _, e := range []string{"a@x.com", "b@x.com", "c@x.com", "d@x.com"} {
		enroll(t, st, e, pid, "europe_africa")
	}
	cid, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "C", TZBand: "europe_africa", StartsOn: "2026-09-03"})
	users, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 2)
	if err != nil {
		t.Fatalf("take: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("claimed %d, want 2", len(users))
	}
	// The rest stay queued for the next cohort.
	groups, _ := st.PlacementQueue(ctx)
	if len(groups) != 1 || groups[0].Waiting != 2 {
		t.Fatalf("remaining queue = %+v, want 2 waiting", groups)
	}
}

func TestBotMentorIsStableAndCannotSignIn(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	id1, err := st.EnsureBotMentor(ctx)
	if err != nil {
		t.Fatalf("ensure bot: %v", err)
	}
	id2, _ := st.EnsureBotMentor(ctx)
	if id1 != id2 {
		t.Fatalf("EnsureBotMentor created two accounts (%d, %d)", id1, id2)
	}
	var hash string
	var isBot int
	if err := st.db.QueryRow(`SELECT password_hash, is_bot FROM users WHERE id = ?`, id1).
		Scan(&hash, &isBot); err != nil {
		t.Fatalf("read bot: %v", err)
	}
	if hash != "" {
		t.Fatal("the bot mentor has a password hash; it must not be able to sign in")
	}
	if isBot != 1 {
		t.Fatal("bot mentor is not flagged as a bot")
	}
}

func TestThinCohortsAndMerge(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")

	// A thin cohort of one and a healthy one on the same path and band.
	thinID, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "thin", TZBand: "europe_africa", StartsOn: "2026-09-01"})
	bigID, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "big", TZBand: "europe_africa", StartsOn: "2026-09-02"})
	lonely := enroll(t, st, "lonely@x.com", pid, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, thinID, pid, "europe_africa", 1); err != nil {
		t.Fatalf("fill thin: %v", err)
	}
	for _, e := range []string{"p@x.com", "q@x.com", "r@x.com"} {
		enroll(t, st, e, pid, "europe_africa")
	}
	if _, err := st.TakeFromQueue(ctx, bigID, pid, "europe_africa", 3); err != nil {
		t.Fatalf("fill big: %v", err)
	}

	thin, err := st.ThinCohorts(ctx, 3)
	if err != nil {
		t.Fatalf("thin: %v", err)
	}
	if len(thin) != 1 || thin[0].ID != thinID {
		t.Fatalf("thin cohorts = %+v, want just the one-member cohort", thin)
	}

	target, err := st.MergeTarget(ctx, thinID, pid, "europe_africa", 1, 8)
	if err != nil {
		t.Fatalf("merge target: %v", err)
	}
	if target != bigID {
		t.Fatalf("merge target = %d, want %d", target, bigID)
	}
	moved, err := st.MergeCohorts(ctx, thinID, target)
	if err != nil || moved != 1 {
		t.Fatalf("merge = %d, %v; want 1, nil", moved, err)
	}

	// The learner is in the surviving cohort, and their enrollment followed.
	c, _ := st.CohortForUser(ctx, lonely)
	if c == nil || c.ID != bigID {
		t.Fatalf("merged learner is in %+v, want cohort %d", c, bigID)
	}
	e, _ := st.LiveEnrollment(ctx, lonely)
	if !e.CohortID.Valid || e.CohortID.Int64 != bigID {
		t.Fatalf("enrollment cohort = %+v, want %d", e.CohortID, bigID)
	}
	// The emptied cohort is archived, not deleted, so its standup history stays.
	emptied, _ := st.GetCohort(ctx, thinID)
	if emptied.State != CohortArchived {
		t.Fatalf("emptied cohort state = %q, want archived", emptied.State)
	}
}

func TestMergeTargetRespectsMaxSize(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	thinID, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "thin", TZBand: "europe_africa", StartsOn: "2026-09-01"})
	fullID, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "full", TZBand: "europe_africa", StartsOn: "2026-09-02"})
	for i := 0; i < 8; i++ {
		enroll(t, st, string(rune('a'+i))+"@full.com", pid, "europe_africa")
	}
	if _, err := st.TakeFromQueue(ctx, fullID, pid, "europe_africa", 8); err != nil {
		t.Fatalf("fill: %v", err)
	}
	// Merging 2 more into a cohort of 8 would exceed the cap, so there is no
	// target — better a thin cohort than an unreadable standup.
	target, err := st.MergeTarget(ctx, thinID, pid, "europe_africa", 2, 8)
	if err != nil {
		t.Fatalf("merge target: %v", err)
	}
	if target != 0 {
		t.Fatalf("merge target = %d, want 0 (no room)", target)
	}
}

func TestMergeTargetWontCrossPathOrBand(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	p1, p2 := mkPath(t, st, "frontend"), mkPath(t, st, "backend")
	thinID, _ := st.CreateCohort(ctx, Cohort{PathID: p1, Name: "thin", TZBand: "europe_africa", StartsOn: "2026-09-01"})
	// Same band, different path.
	st.CreateCohort(ctx, Cohort{PathID: p2, Name: "other-path", TZBand: "europe_africa", StartsOn: "2026-09-02"})
	// Same path, different band.
	st.CreateCohort(ctx, Cohort{PathID: p1, Name: "other-band", TZBand: "asia_pacific", StartsOn: "2026-09-02"})

	target, err := st.MergeTarget(ctx, thinID, p1, "europe_africa", 1, 8)
	if err != nil {
		t.Fatalf("merge target: %v", err)
	}
	if target != 0 {
		t.Fatalf("merged across path or band (target %d)", target)
	}
}

func TestStandupLifecycle(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	cid, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "C", TZBand: "europe_africa", StartsOn: "2026-09-03"})
	uid := enroll(t, st, "learner@x.com", pid, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 1); err != nil {
		t.Fatalf("fill: %v", err)
	}

	// Open a window that is currently live.
	id, created, err := st.EnsureStandup(ctx, cid, "2026-09-03",
		"2000-01-01 00:00:00", "2999-01-01 00:00:00", "What's today?")
	if err != nil || !created {
		t.Fatalf("EnsureStandup = %v, created=%v", err, created)
	}
	// Idempotent: a second call is a no-op, which matters because standup:open
	// runs every 30 minutes.
	id2, created2, err := st.EnsureStandup(ctx, cid, "2026-09-03", "x", "y", "z")
	if err != nil {
		t.Fatalf("second ensure: %v", err)
	}
	if created2 || id2 != id {
		t.Fatalf("second EnsureStandup created a duplicate (%d vs %d, created=%v)", id2, id, created2)
	}

	su, _ := st.TodayStandup(ctx, cid, "2026-09-03")
	if su == nil || !su.IsOpen {
		t.Fatalf("standup not open: %+v", su)
	}

	if posted, _ := st.HasPosted(ctx, id, uid); posted {
		t.Fatal("HasPosted true before posting")
	}
	if err := st.PostStandup(ctx, id, uid, "did a thing", "doing another", "stuck on X"); err != nil {
		t.Fatalf("post: %v", err)
	}
	if posted, _ := st.HasPosted(ctx, id, uid); !posted {
		t.Fatal("HasPosted false after posting")
	}

	// Posting also records the day in the general activity ledger.
	days, _ := st.ActiveDaysBetween(ctx, uid, "2000-01-01", "2999-01-01")
	if len(days) != 1 {
		t.Fatalf("activity days = %v, want one from the standup", days)
	}

	// Re-posting edits rather than duplicating.
	if err := st.PostStandup(ctx, id, uid, "updated", "revised", ""); err != nil {
		t.Fatalf("repost: %v", err)
	}
	entries, _ := st.StandupEntries(ctx, id)
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1 after editing", len(entries))
	}
	if entries[0].Today != "revised" || entries[0].Blockers != "" {
		t.Fatalf("edit did not take: %+v", entries[0])
	}
}

func TestPostAfterCloseIsRejected(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	cid, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "C", TZBand: "europe_africa", StartsOn: "2026-09-03"})
	uid := enroll(t, st, "late@x.com", pid, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 1); err != nil {
		t.Fatalf("fill: %v", err)
	}

	id, _, _ := st.EnsureStandup(ctx, cid, "2026-09-03",
		"2000-01-01 00:00:00", "2000-01-02 00:00:00", "old")
	n, err := st.CloseExpiredStandups(ctx)
	if err != nil || n != 1 {
		t.Fatalf("CloseExpiredStandups = %d, %v; want 1", n, err)
	}
	if err := st.PostStandup(ctx, id, uid, "", "too late", ""); !errors.Is(err, ErrStandupClosed) {
		t.Fatalf("post to closed standup err = %v, want ErrStandupClosed", err)
	}
	// Idempotent sweep.
	if n, _ := st.CloseExpiredStandups(ctx); n != 0 {
		t.Fatalf("second close affected %d rows, want 0", n)
	}
}

func TestMembersListMentorsFirstAndRemovalKeepsHistory(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	pid := mkPath(t, st, "frontend")
	cid, _ := st.CreateCohort(ctx, Cohort{PathID: pid, Name: "C", TZBand: "europe_africa", StartsOn: "2026-09-03"})
	bot, _ := st.EnsureBotMentor(ctx)
	if err := st.AddMember(ctx, cid, bot, RoleMentor); err != nil {
		t.Fatalf("add mentor: %v", err)
	}
	uid := enroll(t, st, "m@x.com", pid, "europe_africa")
	if _, err := st.TakeFromQueue(ctx, cid, pid, "europe_africa", 1); err != nil {
		t.Fatalf("fill: %v", err)
	}

	members, _ := st.CohortMembers(ctx, cid)
	if len(members) != 2 {
		t.Fatalf("got %d members, want 2", len(members))
	}
	if !members[0].IsMentor() {
		t.Fatalf("mentor is not listed first: %+v", members)
	}
	// Adding twice must not duplicate an active membership.
	if err := st.AddMember(ctx, cid, uid, RoleLearner); err != nil {
		t.Fatalf("re-add: %v", err)
	}
	members, _ = st.CohortMembers(ctx, cid)
	if len(members) != 2 {
		t.Fatalf("re-adding duplicated a membership: %d members", len(members))
	}

	// Removal ends the membership without deleting the row, so posts stay
	// attributable.
	if err := st.RemoveMember(ctx, cid, uid); err != nil {
		t.Fatalf("remove: %v", err)
	}
	members, _ = st.CohortMembers(ctx, cid)
	if len(members) != 1 {
		t.Fatalf("after removal: %d active members, want 1", len(members))
	}
	var rows int
	_ = st.db.QueryRow(`SELECT COUNT(*) FROM cohort_members WHERE cohort_id = ?`, cid).Scan(&rows)
	if rows != 2 {
		t.Fatalf("removal deleted history: %d rows remain, want 2", rows)
	}
	if c, _ := st.CohortForUser(ctx, uid); c != nil {
		t.Fatal("a removed member still reads as being in the cohort")
	}
}

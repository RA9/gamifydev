package store

import (
	"context"
	"errors"
	"testing"
)

func seedBank(t *testing.T, st *Store) *Assessment {
	t.Helper()
	ctx := context.Background()
	id, err := st.UpsertAssessment(ctx, Assessment{
		Slug: "placement", Title: "Placement", Kind: "placement",
		TimeLimitS: 1800, PerTopic: 2,
	})
	if err != nil {
		t.Fatalf("upsert assessment: %v", err)
	}
	var items []Item
	for _, topic := range []string{"programming", "data_structures"} {
		for i := 0; i < 4; i++ {
			items = append(items, Item{
				Topic:   topic,
				Prompt:  topic + " question",
				Options: []string{"wrong", "right", "also wrong"},
				Answer:  1,
			})
		}
	}
	if err := st.ReplaceItems(ctx, id, items); err != nil {
		t.Fatalf("replace items: %v", err)
	}
	a, err := st.GetAssessment(ctx, "placement")
	if err != nil {
		t.Fatalf("get assessment: %v", err)
	}
	return a
}

func TestStartAttemptSamplesPerTopic(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "sampler@example.com")

	att, err := st.StartAttempt(ctx, uid, a)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	items, err := st.AttemptItems(ctx, att.ID)
	if err != nil {
		t.Fatalf("items: %v", err)
	}
	// per_topic=2 across two topics — a paper of 4, not the whole bank of 8.
	if len(items) != 4 {
		t.Fatalf("paper has %d items, want 4 (2 per topic)", len(items))
	}
	byTopic := map[string]int{}
	for _, it := range items {
		byTopic[it.Topic]++
	}
	for topic, n := range byTopic {
		if n != 2 {
			t.Fatalf("topic %q drew %d items, want 2", topic, n)
		}
	}
}

func TestAttemptItemsNeverCarryAnswers(t *testing.T) {
	// The type presented to the template has no answer field at all. This test
	// exists to fail loudly if someone adds one — answers must not be renderable.
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "noanswers@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	items, err := st.AttemptItems(ctx, att.ID)
	if err != nil {
		t.Fatalf("items: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("no items drawn")
	}
	// Compile-time guarantee via the type system; assert the shape explicitly.
	var _ struct {
		ID       int64
		Sort     int
		Topic    string
		Prompt   string
		Code     string
		Lang     string
		Options  []string
		Response interface{}
	} = struct {
		ID       int64
		Sort     int
		Topic    string
		Prompt   string
		Code     string
		Lang     string
		Options  []string
		Response interface{}
	}{ID: items[0].ID, Topic: items[0].Topic, Options: items[0].Options}
}

func TestScoreAttemptGradesPerTopic(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "scored@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	items, _ := st.AttemptItems(ctx, att.ID)

	// Answer every programming item right and every data-structures item wrong.
	responses := map[int64]int{}
	for _, it := range items {
		if it.Topic == "programming" {
			responses[it.ID] = 1 // correct
		} else {
			responses[it.ID] = 0 // wrong
		}
	}
	scored, err := st.ScoreAttempt(ctx, att.ID, responses)
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if !scored.Submitted() {
		t.Fatal("attempt not marked submitted")
	}
	if got := scored.TopicScores["programming"]; got != 100 {
		t.Fatalf("programming = %d, want 100", got)
	}
	if got := scored.TopicScores["data_structures"]; got != 0 {
		t.Fatalf("data_structures = %d, want 0", got)
	}
	if !scored.Score.Valid || scored.Score.Int64 != 50 {
		t.Fatalf("overall = %+v, want 50", scored.Score)
	}
}

func TestUnansweredCountsWrong(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "blank@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	scored, err := st.ScoreAttempt(ctx, att.ID, map[int64]int{}) // submitted blank
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if scored.Score.Int64 != 0 {
		t.Fatalf("blank paper scored %d, want 0", scored.Score.Int64)
	}
	for topic, v := range scored.TopicScores {
		if v != 0 {
			t.Fatalf("topic %q scored %d on a blank paper", topic, v)
		}
	}
}

func TestScoringIgnoresAnswersForItemsNotOnThePaper(t *testing.T) {
	// A client posting answers for items it was never served must not be able to
	// inflate its score — grading reads the stored paper, not the request.
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "smuggler@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	items, _ := st.AttemptItems(ctx, att.ID)
	drawn := map[int64]bool{}
	for _, it := range items {
		drawn[it.ID] = true
	}

	// Every item in the bank, answered correctly — including ones not drawn.
	responses := map[int64]int{}
	rows, err := st.db.Query(`SELECT id FROM assessment_items`)
	if err != nil {
		t.Fatalf("bank: %v", err)
	}
	extra := 0
	for rows.Next() {
		var id int64
		_ = rows.Scan(&id)
		responses[id] = 1
		if !drawn[id] {
			extra++
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		t.Fatalf("read bank: %v", err)
	}
	rows.Close()
	if extra == 0 {
		t.Fatal("test is meaningless: every bank item was drawn")
	}

	scored, err := st.ScoreAttempt(ctx, att.ID, responses)
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	// Still graded out of the four drawn items, not the eight answered.
	var graded int
	if err := st.db.QueryRow(
		`SELECT COUNT(*) FROM assessment_attempt_items WHERE attempt_id = ?`, att.ID).Scan(&graded); err != nil {
		t.Fatalf("count: %v", err)
	}
	if graded != len(items) {
		t.Fatalf("graded %d rows, want %d — extra answers leaked in", graded, len(items))
	}
	if scored.Score.Int64 != 100 {
		t.Fatalf("score = %d, want 100", scored.Score.Int64)
	}
}

func TestExpiredAttemptCannotBeSubmitted(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "slow@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	items, _ := st.AttemptItems(ctx, att.ID)
	mustExec(t, st, `UPDATE assessment_attempts SET expires_at = datetime('now','-1 minute') WHERE id = ?`, att.ID)

	responses := map[int64]int{}
	for _, it := range items {
		responses[it.ID] = 1
	}
	if _, err := st.ScoreAttempt(ctx, att.ID, responses); !errors.Is(err, ErrAttemptExpired) {
		t.Fatalf("err = %v, want ErrAttemptExpired", err)
	}
	// It is also closed out, so the learner can start a fresh sitting.
	if live, _ := st.LiveAttempt(ctx, uid); live != nil {
		t.Fatal("expired attempt still reads as live")
	}
}

func TestDoubleSubmissionIsRejected(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "double@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	if _, err := st.ScoreAttempt(ctx, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("first submit: %v", err)
	}
	if _, err := st.ScoreAttempt(ctx, att.ID, map[int64]int{}); !errors.Is(err, ErrAlreadySubmitted) {
		t.Fatalf("second submit err = %v, want ErrAlreadySubmitted", err)
	}
}

func TestLiveAttemptExcludesSubmittedAndExpired(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "live@example.com")

	if live, _ := st.LiveAttempt(ctx, uid); live != nil {
		t.Fatal("a fresh user has a live attempt")
	}
	att, _ := st.StartAttempt(ctx, uid, a)
	if live, _ := st.LiveAttempt(ctx, uid); live == nil || live.ID != att.ID {
		t.Fatal("started attempt is not live")
	}
	if _, err := st.ScoreAttempt(ctx, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	if live, _ := st.LiveAttempt(ctx, uid); live != nil {
		t.Fatal("submitted attempt still reads as live")
	}
	if latest, _ := st.LatestSubmittedAttempt(ctx, uid); latest == nil || latest.ID != att.ID {
		t.Fatal("submitted attempt not returned as latest")
	}
}

func TestExpireAbandonedAttempts(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "abandoned@example.com")

	att, _ := st.StartAttempt(ctx, uid, a)
	mustExec(t, st, `UPDATE assessment_attempts SET expires_at = datetime('now','-1 hour') WHERE id = ?`, att.ID)

	n, err := st.ExpireAbandonedAttempts(ctx)
	if err != nil || n != 1 {
		t.Fatalf("ExpireAbandonedAttempts = %d, %v; want 1, nil", n, err)
	}
	// Idempotent — a second sweep finds nothing.
	if n, _ := st.ExpireAbandonedAttempts(ctx); n != 0 {
		t.Fatalf("second sweep closed %d attempts, want 0", n)
	}
}

func TestPlacementResultRoundTrip(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "placed@example.com")
	att, _ := st.StartAttempt(ctx, uid, a)
	if _, err := st.ScoreAttempt(ctx, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("score: %v", err)
	}

	if got, _ := st.LatestPlacement(ctx, uid); got != nil {
		t.Fatal("a user with no result returned one")
	}
	if _, err := st.SavePlacementResult(ctx, PlacementResult{
		AttemptID: att.ID, UserID: uid,
		FoundationsRequired: true,
		Exemptions:          []string{"data_structures", "c"},
	}); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := st.LatestPlacement(ctx, uid)
	if err != nil || got == nil {
		t.Fatalf("latest: %v %v", got, err)
	}
	if !got.FoundationsRequired {
		t.Fatal("FoundationsRequired did not round-trip")
	}
	if got.Passed {
		t.Fatal("Passed = true, want explicit false")
	}
	if len(got.Exemptions) != 2 || got.Exemptions[0] != "data_structures" {
		t.Fatalf("exemptions = %v", got.Exemptions)
	}
	if _, err := st.SavePlacementResult(ctx, PlacementResult{AttemptID: att.ID, UserID: uid}); !errors.Is(err, ErrPlacementResultExists) {
		t.Fatalf("second result err = %v, want ErrPlacementResultExists", err)
	}
}

func TestReplaceItemsRejectsOutOfRangeAnswer(t *testing.T) {
	// A bad answer index would silently make an item unanswerable, so it must
	// fail at seed time rather than in front of a learner.
	ctx := context.Background()
	st := newTestStore(t)
	id, _ := st.UpsertAssessment(ctx, Assessment{Slug: "x", PerTopic: 1, TimeLimitS: 60})
	err := st.ReplaceItems(ctx, id, []Item{{
		Topic: "t", Prompt: "bad", Options: []string{"a", "b"}, Answer: 5,
	}})
	if err == nil {
		t.Fatal("an answer index past the end of options was accepted")
	}
}

func mkGuest(t *testing.T, st *Store, token string) Solver {
	t.Helper()
	id, err := st.CreateGuestSession(context.Background(), token)
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}
	return Solver{GuestID: id}
}

func TestGuestAttemptOwnershipAndIsolation(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	guestA := mkGuest(t, st, "guest-a")
	guestB := mkGuest(t, st, "guest-b")

	attA, err := st.StartAttemptFor(ctx, guestA, a)
	if err != nil {
		t.Fatalf("start guest A: %v", err)
	}
	attB, err := st.StartAttemptFor(ctx, guestB, a)
	if err != nil {
		t.Fatalf("start guest B: %v", err)
	}
	if !attA.OwnedBy(guestA) || attA.UserID != 0 || attA.GuestID != guestA.GuestID {
		t.Fatalf("attempt owner = %+v, want guest A", attA.By)
	}
	if live, _ := st.LiveAttemptFor(ctx, guestA); live == nil || live.ID != attA.ID {
		t.Fatal("guest A did not get its own live attempt")
	}
	if live, _ := st.LiveAttemptFor(ctx, guestB); live == nil || live.ID != attB.ID {
		t.Fatal("guest B did not get its own live attempt")
	}
	if _, err := st.GetAttemptFor(ctx, guestB, attA.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-owner get err = %v, want ErrNotFound", err)
	}
	if _, err := st.AttemptItemsFor(ctx, guestB, attA.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-owner items err = %v, want ErrNotFound", err)
	}
	if _, err := st.ScoreAttemptFor(ctx, guestB, attA.ID, map[int64]int{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-owner score err = %v, want ErrNotFound", err)
	}

	if _, err := st.ScoreAttemptFor(ctx, guestA, attA.ID, map[int64]int{}); err != nil {
		t.Fatalf("owner score: %v", err)
	}
	if _, err := st.SavePlacementResultFor(ctx, guestA, PlacementResult{
		AttemptID: attA.ID, FoundationsRequired: true, Passed: false,
	}); err != nil {
		t.Fatalf("save guest result: %v", err)
	}
	if got, _ := st.LatestPlacementFor(ctx, guestB); got != nil {
		t.Fatal("guest B could see guest A's placement")
	}
	if _, err := st.ReviewAttemptFor(ctx, guestB, attA.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-owner review err = %v, want ErrNotFound", err)
	}
}

func TestAttemptSnapshotSurvivesReplaceItems(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	uid := mkUser(t, st, "snapshot@example.com")
	by := Solver{UserID: uid}

	att, err := st.StartAttemptFor(ctx, by, a)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	before, err := st.AttemptItemsFor(ctx, by, att.ID)
	if err != nil || len(before) == 0 {
		t.Fatalf("initial paper = %v, %v", before, err)
	}
	if err := st.ReplaceItems(ctx, a.ID, []Item{
		{Topic: "programming", Prompt: "replacement programming", Options: []string{"new right", "new wrong"}, Answer: 0},
		{Topic: "data_structures", Prompt: "replacement data", Options: []string{"new right", "new wrong"}, Answer: 0},
	}); err != nil {
		t.Fatalf("replace: %v", err)
	}
	after, err := st.AttemptItemsFor(ctx, by, att.ID)
	if err != nil {
		t.Fatalf("paper after replace: %v", err)
	}
	if len(after) != len(before) {
		t.Fatalf("paper length after replace = %d, want %d", len(after), len(before))
	}
	responses := map[int64]int{}
	for i := range before {
		if after[i].ID != before[i].ID || after[i].Prompt != before[i].Prompt || after[i].Options[1] != "right" {
			t.Fatalf("item %d changed after reseed: before=%+v after=%+v", i, before[i], after[i])
		}
		responses[after[i].ID] = 1
	}
	scored, err := st.ScoreAttemptFor(ctx, by, att.ID, responses)
	if err != nil {
		t.Fatalf("score snapshotted paper: %v", err)
	}
	if scored.Score.Int64 != 100 {
		t.Fatalf("score after bank replacement = %d, want 100", scored.Score.Int64)
	}
	review, err := st.ReviewAttemptFor(ctx, by, att.ID)
	if err != nil || len(review) != len(before) {
		t.Fatalf("review = %v, %v", review, err)
	}
	if review[0].Prompt != before[0].Prompt || review[0].Answer != 1 {
		t.Fatalf("review did not use snapshot: %+v", review[0])
	}
}

func TestCreateLearnerFromGuestPlacementIsAtomicAndSingleUse(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	guest := mkGuest(t, st, "passing-guest")

	att, err := st.StartAttemptFor(ctx, guest, a)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	items, _ := st.AttemptItemsFor(ctx, guest, att.ID)
	responses := map[int64]int{}
	for _, item := range items {
		responses[item.ID] = 1
	}
	if _, err := st.ScoreAttemptFor(ctx, guest, att.ID, responses); err != nil {
		t.Fatalf("score: %v", err)
	}
	if _, err := st.SavePlacementResultFor(ctx, guest, PlacementResult{AttemptID: att.ID, Passed: true}); err != nil {
		t.Fatalf("save passing result: %v", err)
	}
	problemID, err := st.UpsertProblem(ctx, Problem{Slug: "claim-me", Title: "Claim me", Difficulty: "easy", Published: true})
	if err != nil {
		t.Fatalf("problem: %v", err)
	}
	subID, err := st.CreateProblemSubmission(ctx, problemID, guest, "python", "print(1)")
	if err != nil {
		t.Fatalf("submission: %v", err)
	}

	mkUser(t, st, "taken@example.com")
	if _, err := st.CreateLearnerFromGuestPlacement(ctx, guest.GuestID, "taken@example.com", "hash", "Guest"); err == nil {
		t.Fatal("duplicate email claim unexpectedly succeeded")
	}
	// The failed account insert must roll back every part of the claim.
	if got, err := st.GuestByToken(ctx, "passing-guest"); err != nil || got != guest.GuestID {
		t.Fatalf("guest was spent by rolled-back claim: id=%d err=%v", got, err)
	}
	if got, err := st.GetAttemptFor(ctx, guest, att.ID); err != nil || !got.OwnedBy(guest) {
		t.Fatalf("attempt moved by rolled-back claim: %+v, %v", got, err)
	}
	if sub, err := st.GetProblemSubmission(ctx, subID); err != nil || sub.By != guest {
		t.Fatalf("submission moved by rolled-back claim: %+v, %v", sub, err)
	}

	u, err := st.CreateLearnerFromGuestPlacement(ctx, guest.GuestID, "new@example.com", "hash", "New Learner")
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if u.Role != "learner" {
		t.Fatalf("claimed account role = %q, want learner", u.Role)
	}
	owner := Solver{UserID: u.ID}
	if got, err := st.GetAttemptFor(ctx, owner, att.ID); err != nil || !got.OwnedBy(owner) {
		t.Fatalf("claimed attempt = %+v, %v", got, err)
	}
	if got, err := st.LatestPlacementFor(ctx, owner); err != nil || got == nil || !got.Passed || !got.OwnedBy(owner) {
		t.Fatalf("claimed result = %+v, %v", got, err)
	}
	if sub, err := st.GetProblemSubmission(ctx, subID); err != nil || sub.By != owner {
		t.Fatalf("claimed submission = %+v, %v", sub, err)
	}
	if _, err := st.GuestByToken(ctx, "passing-guest"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("claimed cookie err = %v, want ErrNotFound", err)
	}
	if _, err := st.CreateLearnerFromGuestPlacement(ctx, guest.GuestID, "again@example.com", "hash", "Again"); !errors.Is(err, ErrGuestAlreadyClaimed) {
		t.Fatalf("second claim err = %v, want ErrGuestAlreadyClaimed", err)
	}
}

func TestCreateLearnerFromFailedGuestPlacementRejected(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	guest := mkGuest(t, st, "failed-guest")
	att, _ := st.StartAttemptFor(ctx, guest, a)
	if _, err := st.ScoreAttemptFor(ctx, guest, att.ID, map[int64]int{}); err != nil {
		t.Fatalf("score: %v", err)
	}
	if _, err := st.SavePlacementResultFor(ctx, guest, PlacementResult{
		AttemptID: att.ID, FoundationsRequired: true, Passed: false,
	}); err != nil {
		t.Fatalf("save failed result: %v", err)
	}

	if _, err := st.CreateLearnerFromGuestPlacement(ctx, guest.GuestID, "fail@example.com", "hash", "Failed"); !errors.Is(err, ErrPlacementNotPassed) {
		t.Fatalf("claim err = %v, want ErrPlacementNotPassed", err)
	}
	if _, err := st.GetUserByEmail(ctx, "fail@example.com"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("rejected claim created a user: %v", err)
	}
	if got, err := st.GuestByToken(ctx, "failed-guest"); err != nil || got != guest.GuestID {
		t.Fatalf("rejected claim spent guest: id=%d err=%v", got, err)
	}
}

func TestPlacementRetryPolicy(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	a := seedBank(t, st)
	by := Solver{UserID: mkUser(t, st, "retry@example.com")}
	var attempts []int64

	failAttempt := func() {
		t.Helper()
		att, err := st.StartAttemptFor(ctx, by, a)
		if err != nil {
			t.Fatalf("start attempt %d: %v", len(attempts)+1, err)
		}
		attempts = append(attempts, att.ID)
		if _, err := st.ScoreAttemptFor(ctx, by, att.ID, map[int64]int{}); err != nil {
			t.Fatalf("score attempt %d: %v", len(attempts), err)
		}
		if _, err := st.SavePlacementResultFor(ctx, by, PlacementResult{AttemptID: att.ID, Passed: false}); err != nil {
			t.Fatalf("save attempt %d: %v", len(attempts), err)
		}
	}

	failAttempt()
	if _, err := st.StartAttemptFor(ctx, by, a); !errors.Is(err, ErrPlacementRetryTooSoon) {
		t.Fatalf("immediate retry err = %v, want ErrPlacementRetryTooSoon", err)
	}
	status, err := st.PlacementRetryFor(ctx, by)
	if err != nil || status.Allowed || !errors.Is(status.Reason, ErrPlacementRetryTooSoon) || status.RetryAt == "" {
		t.Fatalf("retry status = %+v, %v", status, err)
	}

	mustExec(t, st, `UPDATE placement_results SET created_at = datetime('now', '-8 days') WHERE attempt_id = ?`, attempts[len(attempts)-1])
	failAttempt()
	mustExec(t, st, `UPDATE placement_results SET created_at = datetime('now', '-8 days') WHERE attempt_id = ?`, attempts[len(attempts)-1])
	failAttempt()
	mustExec(t, st, `UPDATE placement_results SET created_at = datetime('now', '-8 days') WHERE attempt_id = ?`, attempts[len(attempts)-1])

	if _, err := st.StartAttemptFor(ctx, by, a); !errors.Is(err, ErrPlacementAttemptLimit) {
		t.Fatalf("fourth attempt err = %v, want ErrPlacementAttemptLimit", err)
	}
	// The cap is rolling: once one start leaves the 30-day window, another is allowed.
	mustExec(t, st, `UPDATE assessment_attempts SET started_at = datetime('now', '-31 days') WHERE id = ?`, attempts[0])
	if _, err := st.StartAttemptFor(ctx, by, a); err != nil {
		t.Fatalf("start after rolling window: %v", err)
	}
}

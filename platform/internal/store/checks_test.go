package store

import (
	"context"
	"testing"
)

// mkAssignment creates an assignment, optionally Python and required.
func mkAssignment(t *testing.T, st *Store, slug, lang string, passPoints int) int64 {
	t.Helper()
	var id int64
	if err := st.db.QueryRow(`
		INSERT INTO assignments (slug, title, language, max_points, pass_points, required, published)
		VALUES (?, ?, ?, 100, ?, 1, 1) RETURNING id`, slug, slug, lang, passPoints).Scan(&id); err != nil {
		t.Fatalf("assignment: %v", err)
	}
	return id
}

func mkSubmission(t *testing.T, st *Store, assignmentID, userID int64, code string) int64 {
	t.Helper()
	var id int64
	if err := st.db.QueryRow(
		`INSERT INTO submissions (assignment_id, user_id, code) VALUES (?, ?, ?) RETURNING id`,
		assignmentID, userID, code).Scan(&id); err != nil {
		t.Fatalf("submission: %v", err)
	}
	return id
}

func TestAutoGradableRequiresPythonAndChecks(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)

	py := mkAssignment(t, st, "py", "python", 0)
	js := mkAssignment(t, st, "js", "javascript", 0)

	// No checks yet: nothing to run.
	if ok, _ := st.AutoGradable(ctx, py); ok {
		t.Fatal("an assignment with no checks reported auto-gradable")
	}

	checks := []Check{{Label: "returns 3", Test: "f() == 3", Points: 1}}
	if err := st.ReplaceChecks(ctx, py, checks); err != nil {
		t.Fatalf("replace: %v", err)
	}
	if ok, _ := st.AutoGradable(ctx, py); !ok {
		t.Fatal("a Python assignment with checks is not auto-gradable")
	}

	// The sandbox runs Python only. Checks on another language may exist as
	// documentation but must never claim they can be executed.
	if err := st.ReplaceChecks(ctx, js, checks); err != nil {
		t.Fatalf("replace js: %v", err)
	}
	if ok, _ := st.AutoGradable(ctx, js); ok {
		t.Fatal("a JavaScript assignment reported auto-gradable; the sandbox cannot run it")
	}

	// Removing every check takes it back out of the auto-graded pool.
	if err := st.ReplaceChecks(ctx, py, nil); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if ok, _ := st.AutoGradable(ctx, py); ok {
		t.Fatal("still auto-gradable after its checks were removed")
	}
}

func TestAllWeightedChecksMustPassToGrade(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "learner@x.com")
	aid := mkAssignment(t, st, "a", "python", 0)
	if err := st.ReplaceChecks(ctx, aid, []Check{
		{Label: "one", Test: "True", Points: 1},
		{Label: "two", Test: "True", Points: 1},
	}); err != nil {
		t.Fatalf("checks: %v", err)
	}
	list, _ := st.ListChecks(ctx, aid)

	// One of two passing must NOT grade. A checkpoint is a gate, not a curve —
	// partial credit would let a broken solution through.
	sub := mkSubmission(t, st, aid, uid, "code")
	passed, err := st.RecordCheckRun(ctx, sub, map[int64]bool{list[0].ID: true}, "", 100, 0)
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if passed {
		t.Fatal("graded with only half the checks passing")
	}
	got, _ := st.LatestSubmission(ctx, aid, uid)
	if got.Status == "graded" {
		t.Fatalf("status = %q after a partial pass, want it left ungraded", got.Status)
	}
	// Crucially it is not "returned" either: that means a human sent it back,
	// and nobody looked at this.
	if got.Status != "submitted" {
		t.Fatalf("status = %q, want submitted — 'returned' would imply a human reviewed it", got.Status)
	}
	if got.ChecksPassed != 1 || got.ChecksTotal != 2 {
		t.Fatalf("recorded %d/%d, want 1/2", got.ChecksPassed, got.ChecksTotal)
	}

	// All passing grades it, and the score is full marks.
	sub2 := mkSubmission(t, st, aid, uid, "better code")
	passed, err = st.RecordCheckRun(ctx, sub2,
		map[int64]bool{list[0].ID: true, list[1].ID: true}, "", 100, 0)
	if err != nil || !passed {
		t.Fatalf("full pass = %v, %v; want true, nil", passed, err)
	}
	got, _ = st.LatestSubmission(ctx, aid, uid)
	if got.Status != "graded" || !got.Score.Valid || got.Score.Int64 != 100 {
		t.Fatalf("after a full pass: status=%q score=%+v", got.Status, got.Score)
	}
}

func TestAdvisoryChecksDoNotBlock(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "adv@x.com")
	aid := mkAssignment(t, st, "a", "python", 0)
	if err := st.ReplaceChecks(ctx, aid, []Check{
		{Label: "must pass", Test: "True", Points: 1},
		{Label: "nice to have", Test: "True", Points: 0}, // advisory
	}); err != nil {
		t.Fatalf("checks: %v", err)
	}
	list, _ := st.ListChecks(ctx, aid)

	// The weighted check passes, the advisory one fails: still a pass.
	sub := mkSubmission(t, st, aid, uid, "code")
	passed, err := st.RecordCheckRun(ctx, sub,
		map[int64]bool{list[0].ID: true, list[1].ID: false}, "", 100, 0)
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if !passed {
		t.Fatal("a failing advisory check blocked the gate")
	}
	// But it is still reported, so the learner sees the suggestion.
	results, _ := st.SubmissionCheckResults(ctx, sub)
	if len(results) != 2 {
		t.Fatalf("%d results, want 2", len(results))
	}
}

func TestGateSeesAutoGradedSubmission(t *testing.T) {
	// The whole point: auto-grading works by setting the same status/score the
	// human path sets, so the existing progression gate needs no changes.
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "gate@x.com")
	var courseID int64
	if err := st.db.QueryRow(
		`INSERT INTO courses (slug, title, published) VALUES ('c','C',1) RETURNING id`).Scan(&courseID); err != nil {
		t.Fatalf("course: %v", err)
	}
	aid := mkAssignment(t, st, "chk", "python", 80)
	mustExec(t, st, `UPDATE assignments SET course_id = ? WHERE id = ?`, courseID, aid)
	if err := st.ReplaceChecks(ctx, aid, []Check{{Label: "works", Test: "True", Points: 1}}); err != nil {
		t.Fatalf("checks: %v", err)
	}
	list, _ := st.ListChecks(ctx, aid)

	// Before: the gate is shut.
	state, _ := st.AssignmentPassState(ctx, uid, courseID)
	if state[aid] {
		t.Fatal("gate open before any submission")
	}

	sub := mkSubmission(t, st, aid, uid, "code")
	if _, err := st.RecordCheckRun(ctx, sub, map[int64]bool{list[0].ID: true}, "", 100, 80); err != nil {
		t.Fatalf("record: %v", err)
	}

	// After: open, with no human involved.
	state, _ = st.AssignmentPassState(ctx, uid, courseID)
	if !state[aid] {
		t.Fatal("gate still shut after every check passed — auto-grading did not reach the gate")
	}
}

func TestHiddenLabelsWithheldUntilPassed(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "hidden@x.com")
	aid := mkAssignment(t, st, "a", "python", 0)
	if err := st.ReplaceChecks(ctx, aid, []Check{
		{Label: "public one", Test: "True", Points: 1},
		{Label: "handles the empty list", Test: "True", Points: 1, Hidden: true},
	}); err != nil {
		t.Fatalf("checks: %v", err)
	}
	list, _ := st.ListChecks(ctx, aid)

	// Failing: the hidden label must not leak, or the spec can be read off the
	// failure list one resubmission at a time.
	sub := mkSubmission(t, st, aid, uid, "code")
	if _, err := st.RecordCheckRun(ctx, sub,
		map[int64]bool{list[0].ID: true, list[1].ID: false}, "", 100, 0); err != nil {
		t.Fatalf("record: %v", err)
	}
	results, _ := st.SubmissionCheckResults(ctx, sub)
	for _, r := range results {
		if r.Hidden && r.Label == "handles the empty list" {
			t.Fatal("a hidden check's label leaked while the submission was failing")
		}
	}

	// Passing: it is safe to reveal, and useful feedback.
	sub2 := mkSubmission(t, st, aid, uid, "fixed")
	if _, err := st.RecordCheckRun(ctx, sub2,
		map[int64]bool{list[0].ID: true, list[1].ID: true}, "", 100, 0); err != nil {
		t.Fatalf("record: %v", err)
	}
	results, _ = st.SubmissionCheckResults(ctx, sub2)
	var revealed bool
	for _, r := range results {
		if r.Hidden && r.Label == "handles the empty list" {
			revealed = true
		}
	}
	if !revealed {
		t.Fatal("hidden label still withheld after passing")
	}
}

func TestTestExpressionsNeverReachTheLearner(t *testing.T) {
	// The test expression is the answer key.
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "spy@x.com")
	aid := mkAssignment(t, st, "a", "python", 0)
	if err := st.ReplaceChecks(ctx, aid, []Check{
		{Label: "works", Test: "solve([1,2,3]) == 6", Points: 1},
	}); err != nil {
		t.Fatalf("checks: %v", err)
	}
	list, _ := st.ListChecks(ctx, aid)
	sub := mkSubmission(t, st, aid, uid, "code")
	if _, err := st.RecordCheckRun(ctx, sub, map[int64]bool{list[0].ID: false}, "", 100, 0); err != nil {
		t.Fatalf("record: %v", err)
	}
	results, _ := st.SubmissionCheckResults(ctx, sub)
	for _, r := range results {
		if r.Test != "" {
			t.Fatalf("check expression %q was returned to the learner-facing view", r.Test)
		}
	}
}

func TestPendingRunsOnlyIncludeAutoGradable(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	uid := mkUser(t, st, "pending@x.com")

	py := mkAssignment(t, st, "py", "python", 0)
	js := mkAssignment(t, st, "js", "javascript", 0)
	_ = st.ReplaceChecks(ctx, py, []Check{{Label: "x", Test: "True", Points: 1}})
	_ = st.ReplaceChecks(ctx, js, []Check{{Label: "x", Test: "True", Points: 1}})
	pySub := mkSubmission(t, st, py, uid, "code")
	mkSubmission(t, st, js, uid, "code")

	pending, err := st.PendingCheckRuns(ctx, 10)
	if err != nil {
		t.Fatalf("pending: %v", err)
	}
	if len(pending) != 1 || pending[0] != pySub {
		t.Fatalf("pending = %v, want just the Python submission %d", pending, pySub)
	}

	// Once run, it drops out — the job must not re-grade forever.
	list, _ := st.ListChecks(ctx, py)
	if _, err := st.RecordCheckRun(ctx, pySub, map[int64]bool{list[0].ID: true}, "", 100, 0); err != nil {
		t.Fatalf("record: %v", err)
	}
	pending, _ = st.PendingCheckRuns(ctx, 10)
	if len(pending) != 0 {
		t.Fatalf("pending = %v after running, want empty", pending)
	}
}

func TestCheckHealthCountsFailures(t *testing.T) {
	// The signal that a checkpoint is badly worded rather than learners being
	// wrong: one check failing far more than its siblings.
	ctx := context.Background()
	st := newTestStore(t)
	aid := mkAssignment(t, st, "a", "python", 0)
	_ = st.ReplaceChecks(ctx, aid, []Check{
		{Label: "easy", Test: "True", Points: 1},
		{Label: "confusing", Test: "True", Points: 1},
	})
	list, _ := st.ListChecks(ctx, aid)
	for i := 0; i < 3; i++ {
		u := mkUser(t, st, "u"+string(rune('a'+i))+"@x.com")
		sub := mkSubmission(t, st, aid, u, "code")
		_, _ = st.RecordCheckRun(ctx, sub,
			map[int64]bool{list[0].ID: true, list[1].ID: false}, "", 100, 0)
	}
	health, err := st.AssignmentCheckHealth(ctx, aid)
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if len(health) != 2 {
		t.Fatalf("%d rows, want 2", len(health))
	}
	if health[0].Fails != 0 || health[1].Fails != 3 {
		t.Fatalf("fails = %d / %d, want 0 / 3", health[0].Fails, health[1].Fails)
	}
}

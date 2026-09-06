package store

import (
	"context"
	"errors"
	"testing"
)

// live starts a sitting for a fresh candidate and returns who they are and the
// attempt they are sitting.
func live(t *testing.T, st *Store, email string) (Solver, *Attempt) {
	t.Helper()
	ctx := context.Background()
	a := seedBank(t, st)
	by := Solver{UserID: mkUser(t, st, email)}
	att, err := st.StartAttemptFor(ctx, by, a)
	if err != nil {
		t.Fatalf("start attempt: %v", err)
	}
	return by, att
}

// The rule the whole feature exists for: one warning, then the paper closes.
func TestTheFirstViolationWarnsAndTheSecondClosesThePaper(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "one@example.com")

	first, err := st.RecordViolation(ctx, by, att.ID, "copy")
	if err != nil {
		t.Fatalf("first violation: %v", err)
	}
	if first.Strikes != 1 || first.Voided {
		t.Fatalf("first violation = %+v, want one strike and no void", first)
	}
	// Still sittable: a warning must not quietly end the test.
	if got, err := st.LiveAttemptFor(ctx, by); err != nil || got == nil {
		t.Fatalf("the sitting ended on the first warning")
	}

	second, err := st.RecordViolation(ctx, by, att.ID, "hidden")
	if err != nil {
		t.Fatalf("second violation: %v", err)
	}
	if second.Strikes != 2 || !second.Voided {
		t.Fatalf("second violation = %+v, want two strikes and a void", second)
	}
	if got, _ := st.LiveAttemptFor(ctx, by); got != nil {
		t.Error("the sitting is still live after being voided")
	}
}

// Voiding has to be a real failure, not a flag on an otherwise open paper —
// scored zero, closed, and tagged with why.
func TestAVoidedSittingIsClosedScoredZeroAndTagged(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "void@example.com")

	for i := 0; i < ProctorStrikeLimit; i++ {
		if _, err := st.RecordViolation(ctx, by, att.ID, "copy"); err != nil {
			t.Fatalf("violation %d: %v", i, err)
		}
	}
	got, err := st.GetAttemptFor(ctx, by, att.ID)
	if err != nil {
		t.Fatalf("get attempt: %v", err)
	}
	if !got.Voided() || got.VoidedReason != VoidCheating {
		t.Errorf("voided reason = %q, want %q", got.VoidedReason, VoidCheating)
	}
	if !got.Submitted() {
		t.Error("a voided sitting is not closed, so it would still be sittable")
	}
	if !got.Score.Valid || got.Score.Int64 != 0 {
		t.Errorf("score = %v, want 0", got.Score)
	}
}

// The candidate must not be able to submit their way out of a closed paper —
// and must be told the real reason rather than "already submitted".
func TestAVoidedSittingCannotBeScored(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "resubmit@example.com")

	for i := 0; i < ProctorStrikeLimit; i++ {
		if _, err := st.RecordViolation(ctx, by, att.ID, "blur"); err != nil {
			t.Fatalf("violation: %v", err)
		}
	}
	_, err := st.ScoreAttemptFor(ctx, by, att.ID, map[int64]int{})
	if !errors.Is(err, ErrAttemptVoided) {
		t.Errorf("score of a voided attempt = %v, want ErrAttemptVoided", err)
	}
}

// A sitting is not somebody else's to end. Without this, the id of an attempt
// would be enough to fail a stranger's admissions test.
func TestOnlyTheCandidateCanCollectStrikesOnTheirOwnSitting(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	_, att := live(t, st, "mine@example.com")
	stranger := Solver{UserID: mkUser(t, st, "stranger@example.com")}

	if _, err := st.RecordViolation(ctx, stranger, att.ID, "copy"); !errors.Is(err, ErrNotFound) {
		t.Errorf("a stranger reporting on someone else's sitting = %v, want ErrNotFound", err)
	}
	if _, err := st.RecordViolation(ctx, Solver{}, att.ID, "copy"); !errors.Is(err, ErrNoSolver) {
		t.Errorf("an unidentified reporter = %v, want ErrNoSolver", err)
	}
}

// The log is what an administrator reads when somebody appeals, so it has to
// record what was actually seen rather than only how many times.
func TestTheViolationLogKeepsTheEvidence(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "evidence@example.com")

	if _, err := st.RecordViolation(ctx, by, att.ID, "copy"); err != nil {
		t.Fatalf("violation: %v", err)
	}
	if _, err := st.RecordViolation(ctx, by, att.ID, "capture"); err != nil {
		t.Fatalf("violation: %v", err)
	}
	log, err := st.AttemptViolations(ctx, att.ID)
	if err != nil {
		t.Fatalf("violations: %v", err)
	}
	if len(log) != 2 || log[0].Kind != "copy" || log[1].Kind != "capture" {
		t.Errorf("log = %+v, want copy then capture in order", log)
	}
}

// The vocabulary is fixed. A client inventing kinds would fill the log with
// text nobody can act on — and could spend a candidate's strikes on nothing.
func TestAnUnknownViolationKindIsRejected(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "unknown@example.com")

	if _, err := st.RecordViolation(ctx, by, att.ID, "sneezed"); !errors.Is(err, ErrUnknownViolation) {
		t.Errorf("unknown kind = %v, want ErrUnknownViolation", err)
	}
	got, _ := st.GetAttemptFor(ctx, by, att.ID)
	if got.Violations != 0 {
		t.Errorf("a rejected kind still cost a strike (violations = %d)", got.Violations)
	}
}

// A report arriving after the paper is closed must not deepen the hole or
// error: a page shutting down can easily fire one last event.
func TestReportingAfterTheEndIsHarmless(t *testing.T) {
	ctx := context.Background()
	st := newTestStore(t)
	by, att := live(t, st, "late@example.com")

	for i := 0; i < ProctorStrikeLimit; i++ {
		if _, err := st.RecordViolation(ctx, by, att.ID, "copy"); err != nil {
			t.Fatalf("violation: %v", err)
		}
	}
	out, err := st.RecordViolation(ctx, by, att.ID, "copy")
	if err != nil {
		t.Fatalf("late violation: %v", err)
	}
	if out.Strikes != ProctorStrikeLimit || !out.Voided {
		t.Errorf("late report = %+v, want the settled state unchanged", out)
	}
	log, _ := st.AttemptViolations(ctx, att.ID)
	if len(log) != ProctorStrikeLimit {
		t.Errorf("the log grew to %d after the paper closed", len(log))
	}
}

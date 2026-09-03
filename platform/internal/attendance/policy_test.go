package attendance

import "testing"

func days(states ...string) []Day {
	out := make([]Day, len(states))
	for i, s := range states {
		// Day labels only need to be ordered and distinct.
		out[i] = Day{Day: "2026-09-" + pad(i+1), State: s}
	}
	return out
}

func pad(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func TestResolvePrefersExcused(t *testing.T) {
	// Filing a notice must never be punished, even if the learner shows up
	// anyway — otherwise telling your team is a worse move than staying quiet.
	if got := Resolve(true, true); got != Excused {
		t.Fatalf("excused+active = %q, want %q", got, Excused)
	}
	if got := Resolve(true, false); got != Excused {
		t.Fatalf("excused = %q, want %q", got, Excused)
	}
	if got := Resolve(false, true); got != Present {
		t.Fatalf("active = %q, want %q", got, Present)
	}
	if got := Resolve(false, false); got != Absent {
		t.Fatalf("silence = %q, want %q", got, Absent)
	}
}

func TestEvaluateNeedsAConsecutiveRun(t *testing.T) {
	// Below the limit: nothing happens.
	d := Evaluate(days(Absent, Absent))
	if d.Sanction {
		t.Fatalf("sanctioned after 2 absences (limit is %d)", UnexcusedLimit)
	}
	// At the limit: a drop, never a ban.
	d = Evaluate(days(Absent, Absent, Absent))
	if !d.Sanction {
		t.Fatalf("no sanction after %d consecutive absences", UnexcusedLimit)
	}
	if d.Kind != KindDrop {
		t.Fatalf("kind = %q, want %q — attendance never bans", d.Kind, KindDrop)
	}
	if d.Run.Length != 3 {
		t.Fatalf("run length = %d, want 3", d.Run.Length)
	}
}

func TestExcusedBreaksTheRun(t *testing.T) {
	// This is the whole point of the policy: silence is the offense, not absence.
	if d := Evaluate(days(Absent, Excused, Absent, Absent)); d.Sanction {
		t.Fatal("a filed notice failed to break the absence run")
	}
	// The same shape without the notice does trigger.
	if d := Evaluate(days(Absent, Absent, Absent, Absent)); !d.Sanction {
		t.Fatal("four silent days did not trigger")
	}
}

func TestOnlyTheCurrentRunCounts(t *testing.T) {
	// Missed three days, then came back. The rule catches people who have gone,
	// not people who had a bad week and recovered.
	d := Evaluate(days(Absent, Absent, Absent, Present, Present))
	if d.Sanction {
		t.Fatal("sanctioned for a past run the learner already recovered from")
	}
	if d.Run.Length != 0 {
		t.Fatalf("current run = %d, want 0 after returning", d.Run.Length)
	}
	// But the history is still visible for an operator.
	if got := LongestRun(days(Absent, Absent, Absent, Present)).Length; got != 3 {
		t.Fatalf("LongestRun = %d, want 3", got)
	}
}

func TestPresentBreaksTheRunToo(t *testing.T) {
	if d := Evaluate(days(Absent, Absent, Present, Absent, Absent)); d.Sanction {
		t.Fatal("activity mid-run did not reset the count")
	}
}

func TestEmptyAndAllPresent(t *testing.T) {
	if d := Evaluate(nil); d.Sanction {
		t.Fatal("sanctioned a learner with no scheduled days yet")
	}
	if d := Evaluate(days(Present, Present, Present, Present)); d.Sanction {
		t.Fatal("sanctioned a learner who attended everything")
	}
}

func TestRunBoundsAreReported(t *testing.T) {
	// The window is recorded on the sanction so it can be appealed against
	// something specific rather than a vague accusation.
	d := Evaluate(days(Present, Absent, Absent, Absent))
	if d.Run.From != "2026-09-02" || d.Run.To != "2026-09-04" {
		t.Fatalf("run bounds = %s..%s, want 2026-09-02..2026-09-04", d.Run.From, d.Run.To)
	}
	if d.Reason == "" {
		t.Fatal("no reason recorded")
	}
}

func TestExcusedBudget(t *testing.T) {
	d := days(Excused, Excused, Present)
	if got := ExcusedUsed(d); got != 2 {
		t.Fatalf("used = %d, want 2", got)
	}
	if got := BudgetLeft(d); got != ExcusedBudget-2 {
		t.Fatalf("left = %d, want %d", got, ExcusedBudget-2)
	}
	// Overspending clamps at zero rather than going negative.
	var many []Day
	for i := 0; i < ExcusedBudget+5; i++ {
		many = append(many, Day{Day: "d", State: Excused})
	}
	if got := BudgetLeft(many); got != 0 {
		t.Fatalf("left = %d, want 0", got)
	}
}

func TestRateCountsExcusedAsAttended(t *testing.T) {
	// A learner who told us they'd be out did what was asked; their rate should
	// not read like someone who vanished.
	if got := Rate(days(Present, Excused, Present, Absent)); got != 75 {
		t.Fatalf("rate = %d, want 75", got)
	}
	if got := Rate(nil); got != 100 {
		t.Fatalf("rate with no days = %d, want 100", got)
	}
	if got := Rate(days(Absent, Absent)); got != 0 {
		t.Fatalf("rate = %d, want 0", got)
	}
}

func TestPolicyIsInternallyConsistent(t *testing.T) {
	// Guardrails on the provisional numbers. If someone tunes these, these
	// assertions are what stop the policy becoming incoherent.
	if UnexcusedLimit < 2 {
		t.Fatalf("UnexcusedLimit of %d would drop learners for a single miss", UnexcusedLimit)
	}
	if ExcusedBudget < UnexcusedLimit {
		t.Fatalf("ExcusedBudget (%d) below UnexcusedLimit (%d): a learner could not "+
			"excuse their way through a stretch the rule would drop them for",
			ExcusedBudget, UnexcusedLimit)
	}
}

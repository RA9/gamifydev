// Package attendance decides what a learner's silence means.
//
// Pure, and deliberately small: these constants are the entire enforcement
// policy, and the PRD is explicit that they cannot be chosen correctly in
// advance (§07). They must be calibrated against a real cohort's attendance
// before enforcement is switched on — which is why every sanction this package
// justifies is recorded in shadow mode by default.
//
// The governing decision (§03): the offense is going dark, not being absent.
// A learner who files a notice is excused. Removal is a drop, not a ban.
package attendance

import "time"

// Daily states, resolved per scheduled day.
const (
	// Present — posted standup, or completed scheduled work, that day.
	Present = "present"
	// Excused — covered by an absence notice.
	Excused = "excused"
	// Absent — silence on a day work was due.
	Absent = "absent"
)

// Sanction kinds. Note what is absent: there is no ban. Attendance can remove
// you from a cohort; only conduct gets you banned.
const (
	KindDrop    = "drop"
	KindSuspend = "suspend"
)

const (
	// UnexcusedLimit is how many consecutive SCHEDULED days of silence trigger a
	// drop. Counting scheduled days rather than calendar days is what keeps
	// weekends and rest days from burning a learner's margin.
	//
	// Provisional. Calibrate before enforcing.
	UnexcusedLimit = 3

	// ExcusedBudget is how many days a learner may excuse per cohort. Generous
	// enough for illness, exams and travel; finite so that "excused" still
	// means something.
	//
	// Provisional. Calibrate before enforcing.
	ExcusedBudget = 8

	// ReapplyAfter is how long a dropped learner waits before reapplying. Short,
	// because the point is a fresh cohort rather than a punishment.
	ReapplyAfter = 14 * 24 * time.Hour
)

// Day is one resolved scheduled day for a learner.
type Day struct {
	Day   string // YYYY-MM-DD
	State string
}

// Resolve decides a single day's state from the signals available.
//
// Order matters: a filed notice wins even if the learner also showed up, so
// telling your team you'll be out is never punished if you turn up anyway.
func Resolve(excused, active bool) string {
	switch {
	case excused:
		return Excused
	case active:
		return Present
	default:
		return Absent
	}
}

// Run is a stretch of consecutive unexcused absences.
type Run struct {
	From, To string
	Length   int
}

// LongestRun finds the longest stretch of consecutive absent days in a
// chronologically ordered slice.
//
// "Consecutive" means consecutive *scheduled* days — the caller supplies only
// scheduled days, so a weekend between two absences does not break the run, and
// neither does it count toward it.
func LongestRun(days []Day) Run {
	var best, cur Run
	for _, d := range days {
		if d.State != Absent {
			cur = Run{}
			continue
		}
		if cur.Length == 0 {
			cur.From = d.Day
		}
		cur.To = d.Day
		cur.Length++
		if cur.Length > best.Length {
			best = cur
		}
	}
	return best
}

// CurrentRun is the run of absences ending at the most recent day, which is what
// enforcement acts on — a learner who went quiet and came back is not sanctioned
// for the stretch they already recovered from.
func CurrentRun(days []Day) Run {
	var cur Run
	for _, d := range days {
		if d.State != Absent {
			cur = Run{}
			continue
		}
		if cur.Length == 0 {
			cur.From = d.Day
		}
		cur.To = d.Day
		cur.Length++
	}
	return cur
}

// Decision is what the policy would do about a learner right now.
type Decision struct {
	Sanction bool
	Kind     string
	Run      Run
	Reason   string
}

// Evaluate applies the rolling window to a learner's resolved days.
//
// Only the run ending at the most recent scheduled day counts. Someone who
// missed three days a month ago and has attended since is not dropped for it —
// the rule exists to catch people who have gone, not to keep a ledger of sins.
func Evaluate(days []Day) Decision {
	run := CurrentRun(days)
	if run.Length < UnexcusedLimit {
		return Decision{Run: run}
	}
	return Decision{
		Sanction: true,
		Kind:     KindDrop,
		Run:      run,
		Reason: "no standup or activity on " + itoa(run.Length) +
			" consecutive scheduled days, and no absence notice",
	}
}

// ExcusedUsed counts how much of the budget a learner has spent.
func ExcusedUsed(days []Day) int {
	n := 0
	for _, d := range days {
		if d.State == Excused {
			n++
		}
	}
	return n
}

// BudgetLeft is how many excused days a learner has remaining.
func BudgetLeft(days []Day) int {
	left := ExcusedBudget - ExcusedUsed(days)
	if left < 0 {
		return 0
	}
	return left
}

// Rate is the share of scheduled days attended (present or excused), 0-100.
// Excused days count as attended: the learner did what was asked of them.
func Rate(days []Day) int {
	if len(days) == 0 {
		return 100
	}
	ok := 0
	for _, d := range days {
		if d.State != Absent {
			ok++
		}
	}
	return ok * 100 / len(days)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

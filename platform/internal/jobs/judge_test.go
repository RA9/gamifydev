package jobs

import (
	"testing"

	"github.com/RA9/gamifydev/platform/internal/store"
)

func TestVerdictNamesWhyTheProgramFailedBeforeHowManyTestsItFailed(t *testing.T) {
	tests := []struct {
		name        string
		res         GradeResult
		passed, tot int
		want        string
	}{
		{
			// The whole point of the ordering: a program that didn't build
			// fails every test, and "0/12 wrong" would send the learner
			// hunting for a logic bug in code that never ran.
			name: "a build failure outranks the failed count",
			res:  GradeResult{CompileFailed: true, Crashed: true}, passed: 0, tot: 12,
			want: store.VerdictCompileError,
		},
		{
			name: "all tests passing is accepted",
			res:  GradeResult{}, passed: 12, tot: 12,
			want: store.VerdictAccepted,
		},
		{
			// Killed on the way out, after every answer was already correct.
			// Failing this submission would be punishing a solved problem for
			// the sandbox's shutdown timing.
			name: "a timeout after the last correct answer is still accepted",
			res:  GradeResult{TimedOut: true}, passed: 5, tot: 5,
			want: store.VerdictAccepted,
		},
		{
			name: "too slow outranks the tests it never reached",
			res:  GradeResult{TimedOut: true}, passed: 3, tot: 10,
			want: store.VerdictTimeLimit,
		},
		{
			name: "a crash is named as one, not as a wrong answer",
			res:  GradeResult{Crashed: true}, passed: 0, tot: 10,
			want: store.VerdictRuntimeError,
		},
		{
			name: "a program that ran and finished is judged on its answers",
			res:  GradeResult{}, passed: 9, tot: 10,
			want: store.VerdictWrongAnswer,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := verdictFor(tc.res, tc.passed, tc.tot); got != tc.want {
				t.Errorf("verdictFor() = %q, want %q", got, tc.want)
			}
		})
	}
}

// A submission is only accepted when every test passed. Partial credit is fine
// for coursework, where the point is to measure learning; a problem is either
// solved or it isn't, and a bank that called 9/10 "solved" would be lying on
// the one screen learners trust to tell them where they stand.
func TestNothingShortOfEveryTestIsAccepted(t *testing.T) {
	for passed := 0; passed < 10; passed++ {
		if got := verdictFor(GradeResult{}, passed, 10); got == store.VerdictAccepted {
			t.Errorf("%d/10 was accepted", passed)
		}
	}
}

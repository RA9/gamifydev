package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// registerJudgeJobs adds the practice-problem judge.
//
// The judge is deliberately the same machinery as checkpoint grading — one
// grader, one sandbox, one set of semantics. A problem bank with its own
// private notion of "passed" would eventually disagree with the coursework
// about whether the same program is correct.
func registerJudgeJobs(r *Runner, st *store.Store, exec runner.Executor) {
	if exec == nil || !exec.Enabled() {
		// Without a sandbox nothing can be judged, and a job that reports
		// healthy while leaving every submission pending is worse than none.
		return
	}
	r.Register(Job{
		Name: "judge:run",
		// Faster than checkpoint grading: a problem is something a visitor is
		// sitting and watching, not homework they submit and come back to.
		Every:   15 * time.Second,
		Timeout: 2 * time.Minute,
		Run: func(ctx context.Context) (string, error) {
			ids, err := st.PendingJudgeRuns(ctx, 10)
			if err != nil {
				return "", err
			}
			judged, accepted := 0, 0
			for _, id := range ids {
				verdict, err := Judge(ctx, st, exec, id)
				if err != nil {
					// Leaving the row pending would make one pathological
					// program retry forever and block the queue behind it, so
					// the failure is recorded on the submission itself.
					if rerr := st.RecordVerdict(ctx, id, store.Verdict{
						Name:   store.VerdictError,
						Output: "The judge couldn't run your program. This is our fault, not yours — try again.",
					}); rerr != nil {
						return fmt.Sprintf("judged %d, then submission %d failed: %v", judged, id, err), rerr
					}
					judged++
					continue
				}
				judged++
				if verdict == store.VerdictAccepted {
					accepted++
				}
			}
			if judged == 0 {
				return "", nil
			}
			return fmt.Sprintf("judged %d submission(s), %d accepted", judged, accepted), nil
		},
	})
}

// Judge grades one submission and records its verdict, returning the verdict it
// recorded.
//
// Exported so the submit handler can run it inline while the visitor is still
// watching, instead of making them wait out a poll interval for a program that
// finishes in under a second. The job calls exactly the same function, and
// remains the backstop for anything the inline attempt didn't get to.
func Judge(ctx context.Context, st *store.Store, exec runner.Executor, subID int64) (string, error) {
	sub, err := st.GetProblemSubmission(ctx, subID)
	if err != nil {
		return "", err
	}
	tests, err := st.ProblemTests(ctx, sub.ProblemID)
	if err != nil {
		return "", err
	}
	if len(tests) == 0 {
		// An unfinished problem must not read as solved.
		return store.VerdictError, st.RecordVerdict(ctx, subID, store.Verdict{
			Name:   store.VerdictError,
			Output: "This problem has no tests yet, so nothing can be judged.",
		})
	}

	res, err := GradeWithin(ctx, exec, sub.Language, sub.Code, nil, tests, sub.TimeLimitMs)
	if err != nil {
		return "", err
	}

	passed, failedLabel := 0, ""
	for _, t := range tests {
		if res.Passed[t.ID] {
			passed++
			continue
		}
		// Name the first requirement that broke — but only if the author meant
		// it to be visible. A hidden test's label is the answer key.
		if failedLabel == "" && !t.Hidden {
			failedLabel = t.Label
		}
	}
	verdict := verdictFor(res, passed, len(tests))
	// A program that never ran didn't fail a particular requirement, it failed
	// all of them for one reason — and that reason is already the headline.
	// Pointing at the first test would send someone reading a traceback off to
	// study an assertion that was never evaluated.
	if verdict != store.VerdictWrongAnswer && verdict != store.VerdictTimeLimit {
		failedLabel = ""
	}
	return verdict, st.RecordVerdict(ctx, subID, store.Verdict{
		Name:        verdict,
		Passed:      passed,
		Total:       len(tests),
		RuntimeMs:   int(res.DurationMs),
		Output:      res.Output,
		FailedLabel: failedLabel,
	})
}

// verdictFor names what went wrong, in the order the learner most needs to hear
// it.
//
// Why the program failed outranks how many tests it failed: a submission that
// didn't compile fails everything, and telling someone "0/12 wrong answers"
// when the real news is a missing semicolon sends them hunting for a bug in
// logic that never ran. Only once the program built, ran, and finished does the
// count of failed tests become the honest headline.
func verdictFor(res GradeResult, passed, total int) string {
	switch {
	case res.CompileFailed:
		return store.VerdictCompileError
	case passed == total:
		// Checked before the timeout: a run that answered everything correctly
		// and was killed while shutting down still solved the problem.
		return store.VerdictAccepted
	case res.TimedOut:
		return store.VerdictTimeLimit
	case res.Crashed:
		return store.VerdictRuntimeError
	default:
		return store.VerdictWrongAnswer
	}
}

package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RA9/gamifydev/platform/internal/pyharness"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// registerCheckJobs adds automated checkpoint grading.
//
// Runs as a job rather than inline on submit for two reasons: a sandbox run
// takes seconds and should not block the learner's request, and a burst of
// submissions (a whole cohort hitting the same checkpoint on the same evening,
// which the shared schedule guarantees) needs to be metered rather than
// stampede the runner.
func registerCheckJobs(r *Runner, st *store.Store, exec runner.Executor) {
	if exec == nil || !exec.Enabled() {
		// Registering a job that cannot work would report healthy while doing
		// nothing. Say why, and leave it unregistered.
		return
	}
	r.Register(Job{
		Name:    "checks:run",
		Every:   time.Minute,
		Timeout: 5 * time.Minute,
		Run: func(ctx context.Context) (string, error) {
			// A bounded batch keeps one slow submission from starving the rest
			// and keeps the job inside its timeout.
			ids, err := st.PendingCheckRuns(ctx, 20)
			if err != nil {
				return "", err
			}
			ran, passed := 0, 0
			for _, id := range ids {
				ok, err := runChecks(ctx, st, exec, id)
				if err != nil {
					// One bad submission must not abandon the batch — record it
					// and move on, or a single pathological program blocks every
					// other learner's result.
					return fmt.Sprintf("ran %d, then submission %d failed: %v", ran, id, err), nil
				}
				ran++
				if ok {
					passed++
				}
			}
			if ran == 0 {
				return "", nil
			}
			return fmt.Sprintf("graded %d submission(s), %d passed", ran, passed), nil
		},
	})
}

// runChecks executes one submission's checks in the sandbox and records the
// outcome. Returns whether every weighted check passed.
func runChecks(ctx context.Context, st *store.Store, exec runner.Executor, subID int64) (bool, error) {
	sub, err := st.GetSubmissionForChecks(ctx, subID)
	if err != nil {
		return false, err
	}
	checks, err := st.ListChecks(ctx, sub.AssignmentID)
	if err != nil {
		return false, err
	}
	if len(checks) == 0 {
		return false, nil
	}

	tests := make([]string, len(checks))
	for i, c := range checks {
		tests[i] = c.Test
	}

	res, err := exec.Run(ctx, runner.Request{Code: pyharness.Build(sub.Code, tests)})
	if err != nil {
		return false, err
	}
	results, logs := pyharness.ParseOutput(res.Stdout, len(tests))

	passedIDs := make(map[int64]bool, len(checks))
	for i, c := range checks {
		if i < len(results) {
			passedIDs[c.ID] = results[i]
		}
	}

	// What the learner sees. Their own printed output first, then the reason the
	// run itself failed if it did — a timeout or traceback explains a wall of
	// failed checks far better than the checks do.
	var out strings.Builder
	if len(logs) > 0 {
		out.WriteString(strings.Join(logs, "\n"))
	}
	switch {
	case res.TimedOut:
		out.WriteString("\n\n⚠ Your program didn't finish in time. Check for a loop that never ends.")
	case strings.TrimSpace(res.Stderr) != "":
		out.WriteString("\n\n" + strings.TrimSpace(res.Stderr))
	}

	return st.RecordCheckRun(ctx, subID, passedIDs, strings.TrimSpace(out.String()),
		sub.MaxPoints, sub.PassPoints)
}

package jobs

import (
	"context"
	"fmt"
	"strconv"
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

	var passedIDs map[int64]bool
	var output string
	switch sub.Language {
	case runner.LangC:
		passedIDs, output, err = runC(ctx, exec, sub.Code, checks)
	default:
		passedIDs, output, err = runPython(ctx, exec, sub.Code, checks)
	}
	if err != nil {
		return false, err
	}
	return st.RecordCheckRun(ctx, subID, passedIDs, output, sub.MaxPoints, sub.PassPoints)
}

// runPython verifies an interpreted submission: one run, with every check
// evaluated against the namespace the learner's program left behind.
func runPython(ctx context.Context, exec runner.Executor, code string, checks []store.Check) (map[int64]bool, string, error) {
	tests := make([]string, len(checks))
	for i, c := range checks {
		tests[i] = c.Test
	}
	res, err := exec.Run(ctx, runner.Request{Lang: runner.LangPython, Code: pyharness.Build(code, tests)})
	if err != nil {
		return nil, "", err
	}
	results, logs := pyharness.ParseOutput(res.Stdout, len(tests))
	passed := make(map[int64]bool, len(checks))
	for i, c := range checks {
		if i < len(results) {
			passed[c.ID] = results[i]
		}
	}
	return passed, learnerOutput(strings.Join(logs, "\n"), res), nil
}

// runC verifies a compiled submission.
//
// A C program leaves no namespace to inspect, so each check asserts over what
// the program printed for a given input. Checks are grouped by stdin so a
// checkpoint whose assertions share one input compiles once rather than once
// per check.
func runC(ctx context.Context, exec runner.Executor, code string, checks []store.Check) (map[int64]bool, string, error) {
	// Preserve author order within each group, and group order by first
	// appearance, so the output a learner sees is stable between runs.
	var order []string
	groups := map[string][]store.Check{}
	for _, c := range checks {
		if _, seen := groups[c.Stdin]; !seen {
			order = append(order, c.Stdin)
		}
		groups[c.Stdin] = append(groups[c.Stdin], c)
	}

	passed := make(map[int64]bool, len(checks))
	var out strings.Builder
	for _, stdin := range order {
		group := groups[stdin]
		res, err := exec.Run(ctx, runner.Request{Lang: runner.LangC, Code: code, Stdin: stdin})
		if err != nil {
			return nil, "", err
		}
		if res.CompileFailed {
			// Nothing ran, so every check fails — but the compiler's message is
			// the only useful thing to show, not a wall of failed assertions.
			return passed, "Your program didn't compile:\n\n" + strings.TrimSpace(res.Stderr), nil
		}

		tests := make([]string, len(group))
		for i, c := range group {
			tests[i] = c.Test
		}
		prog := pyharness.BuildAssertions(
			map[string]string{"_out": res.Stdout, "_err": res.Stderr, "_code": code, "_in": stdin},
			map[string]int{"_exit": res.ExitCode},
			tests)
		ares, err := exec.Run(ctx, runner.Request{Lang: runner.LangPython, Code: prog})
		if err != nil {
			return nil, "", err
		}
		results, _ := pyharness.ParseOutput(ares.Stdout, len(tests))
		for i, c := range group {
			if i < len(results) {
				passed[c.ID] = results[i]
			}
		}

		// Show the learner what their program actually printed for each input —
		// for a stdin-driven exercise that is the whole debugging story.
		if len(order) > 1 && strings.TrimSpace(stdin) != "" {
			out.WriteString("$ echo " + strconv.Quote(strings.TrimSpace(stdin)) + " | ./prog\n")
		}
		out.WriteString(strings.TrimRight(res.Stdout, "\n"))
		if s := strings.TrimSpace(res.Stderr); s != "" {
			out.WriteString("\n" + s)
		}
		if res.TimedOut {
			out.WriteString("\n⚠ Timed out — check for a loop that never ends.")
		}
		out.WriteString("\n")
	}
	return passed, strings.TrimSpace(out.String()), nil
}

// learnerOutput assembles what to show beneath the check list: the program's
// own output, then the reason the run itself failed. A timeout or traceback
// explains a wall of failed checks far better than the checks do.
func learnerOutput(logs string, res runner.Result) string {
	var out strings.Builder
	out.WriteString(logs)
	switch {
	case res.TimedOut:
		out.WriteString("\n\n⚠ Your program didn't finish in time. Check for a loop that never ends.")
	case strings.TrimSpace(res.Stderr) != "":
		out.WriteString("\n\n" + strings.TrimSpace(res.Stderr))
	}
	return strings.TrimSpace(out.String())
}

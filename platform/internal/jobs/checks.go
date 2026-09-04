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

	// Fixture files the program is meant to operate on, if the author set any.
	files, err := st.FileMap(ctx, sub.AssignmentID)
	if err != nil {
		return false, err
	}

	res, err := Grade(ctx, exec, sub.Language, sub.Code, files, checks)
	if err != nil {
		return false, err
	}
	return st.RecordCheckRun(ctx, subID, res.Passed, res.Output, sub.MaxPoints, sub.PassPoints)
}

// GradeResult is what one graded program yields.
//
// It carries why a run went wrong, not just which checks failed, because the
// two callers want different things from it: a checkpoint shows the learner a
// transcript, while the problem judge has to name a verdict — and "didn't
// compile" and "ran too long" are different verdicts, indistinguishable from a
// list of failed checks alone.
type GradeResult struct {
	Passed        map[int64]bool
	Output        string
	CompileFailed bool
	TimedOut      bool
	// Crashed means the program itself died before it could be judged, rather
	// than answering wrongly. Only claimed when we actually know: a non-zero
	// exit for a compiled program, or a Python run that never reached the first
	// check because it raised on the way there.
	Crashed    bool
	DurationMs int64
}

// Grade runs a set of checks against one program.
//
// Separated from the job loop and exported so a checkpoint can be graded
// without a database or a submission behind it — which is what lets a test
// prove that an authored checkpoint is actually solvable, rather than waiting
// for a learner to discover that it isn't.
func Grade(ctx context.Context, exec runner.Executor, lang, code string, files map[string]string, checks []store.Check) (GradeResult, error) {
	return GradeWithin(ctx, exec, lang, code, files, checks, 0)
}

// GradeWithin is Grade with an explicit wall-clock budget for the learner's
// program, in milliseconds; 0 means the runner's own default.
//
// A checkpoint doesn't need this — it asks whether the learner understood
// something, and the runner's default is a generous backstop. A practice
// problem does: "fast enough" is part of the exercise, so each problem sets the
// limit its intended solution comfortably fits and a brute force does not.
func GradeWithin(ctx context.Context, exec runner.Executor, lang, code string, files map[string]string, checks []store.Check, timeoutMs int) (GradeResult, error) {
	switch lang {
	case runner.LangC, runner.LangShell:
		// Neither leaves a namespace to inspect, so both are checked by
		// asserting over what the program printed for a given input.
		return runByOutput(ctx, exec, lang, code, files, checks, timeoutMs)
	default:
		return runPython(ctx, exec, code, files, checks, timeoutMs)
	}
}

// runPython verifies an interpreted submission: one run, with every check
// evaluated against the namespace the learner's program left behind.
func runPython(ctx context.Context, exec runner.Executor, code string, files map[string]string, checks []store.Check, timeoutMs int) (GradeResult, error) {
	tests := make([]string, len(checks))
	for i, c := range checks {
		tests[i] = c.Test
	}
	res, err := exec.Run(ctx, runner.Request{Lang: runner.LangPython, Code: pyharness.Build(code, tests), Files: files, TimeoutMs: timeoutMs})
	if err != nil {
		return GradeResult{}, err
	}
	results, logs := pyharness.ParseOutput(res.Stdout, len(tests))
	passed := make(map[int64]bool, len(checks))
	for i, c := range checks {
		if i < len(results) {
			passed[c.ID] = results[i]
		}
	}
	return GradeResult{
		Passed:   passed,
		Output:   learnerOutput(strings.Join(logs, "\n"), res),
		TimedOut: res.TimedOut,
		// Not one check reported, yet the program had something to say on
		// stderr: it raised before the harness ever ran.
		Crashed:    !res.TimedOut && !strings.Contains(res.Stdout, pyharness.ResultMarker) && strings.TrimSpace(res.Stderr) != "",
		DurationMs: res.DurationMs,
	}, nil
}

// runByOutput verifies a submission that leaves no inspectable namespace —
// a compiled program or a shell script.
//
// Each check asserts over what the program printed for a given input. Checks
// are grouped by stdin so a checkpoint whose assertions share one input runs
// once rather than once per check.
func runByOutput(ctx context.Context, exec runner.Executor, lang, code string, files map[string]string, checks []store.Check, timeoutMs int) (GradeResult, error) {
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

	out := GradeResult{Passed: make(map[int64]bool, len(checks))}
	var transcript strings.Builder
	for _, stdin := range order {
		group := groups[stdin]
		res, err := exec.Run(ctx, runner.Request{Lang: lang, Code: code, Stdin: stdin, Files: files, TimeoutMs: timeoutMs})
		if err != nil {
			return GradeResult{}, err
		}
		// Only the learner's own program counts toward the reported runtime —
		// the assertion pass is our overhead, not theirs.
		out.DurationMs += res.DurationMs
		out.TimedOut = out.TimedOut || res.TimedOut
		out.Crashed = out.Crashed || (!res.TimedOut && !res.CompileFailed && res.ExitCode != 0)
		if res.CompileFailed {
			// Nothing ran, so every check fails — but the compiler's message is
			// the only useful thing to show, not a wall of failed assertions.
			out.CompileFailed = true
			out.Output = "Your program didn't compile:\n\n" + strings.TrimSpace(res.Stderr)
			return out, nil
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
			return GradeResult{}, err
		}
		results, _ := pyharness.ParseOutput(ares.Stdout, len(tests))
		for i, c := range group {
			if i < len(results) {
				out.Passed[c.ID] = results[i]
			}
		}

		// Show the learner what their program actually printed for each input —
		// for a stdin-driven exercise that is the whole debugging story.
		if len(order) > 1 && strings.TrimSpace(stdin) != "" {
			transcript.WriteString("$ echo " + strconv.Quote(strings.TrimSpace(stdin)) + " | " + progName(lang) + "\n")
		}
		transcript.WriteString(strings.TrimRight(res.Stdout, "\n"))
		if s := strings.TrimSpace(res.Stderr); s != "" {
			transcript.WriteString("\n" + s)
		}
		if res.TimedOut {
			transcript.WriteString("\n⚠ Timed out — check for a loop that never ends.")
		}
		transcript.WriteString("\n")
	}
	out.Output = strings.TrimSpace(transcript.String())
	return out, nil
}

// progName is how the learner's program is referred to in the transcript shown
// beneath the checks.
func progName(lang string) string {
	if lang == runner.LangShell {
		return "./main.sh"
	}
	return "./prog"
}

// learnerOutput assembles what to show beneath the check list: the program's
// own output, then the reason the run itself failed. A timeout or traceback
// explains a wall of failed checks far better than the checks do.
func learnerOutput(logs string, res runner.Result) string {
	var out strings.Builder
	out.WriteString(logs)
	switch {
	case res.TimedOut:
		// Checks report one at a time, so the ones above this line did run. What
		// stopped is everything from the first unticked check onward — usually an
		// algorithm too slow for the input it was handed, or a loop with no exit.
		out.WriteString("\n\n⚠ Your program ran out of time. Every check below the last one it reached is untested — look for an approach that's too slow for the input, or a loop that never ends.")
	case strings.TrimSpace(res.Stderr) != "":
		out.WriteString("\n\n" + strings.TrimSpace(res.Stderr))
	}
	return strings.TrimSpace(out.String())
}

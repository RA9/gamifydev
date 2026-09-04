package seed

import (
	"os/exec"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/runner"
)

func newPythonExecutor(t *testing.T) runner.Executor {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3 on PATH")
	}
	e, err := runner.New(runner.Config{Mode: "local", PythonPath: "python3", AllowUnsafe: true})
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	return e
}

// Every problem ships with an answer, and the answer has to pass. A problem
// nobody can solve is worse than no problem at all: the learner has no way to
// tell an impossible exercise from their own mistake, and the bank's whole
// claim is that a green result means you got it right.
func TestEverySeededProblemIsSolvableByItsOwnSolution(t *testing.T) {
	e := newPythonExecutor(t)
	for _, p := range seedProblems {
		t.Run(p.slug, func(t *testing.T) {
			t.Parallel()
			if p.solution == "" {
				t.Fatalf("problem %q ships no reference solution, so nothing proves it can be solved", p.slug)
			}
			// Give the tests the ids the store would have; Grade reports by id.
			tests := p.tests
			for i := range tests {
				tests[i].ID = int64(i + 1)
			}
			res, err := jobs.GradeWithin(t.Context(), e, "python", p.solution, nil, tests, p.timeLimitMs)
			if err != nil {
				t.Fatalf("grade: %v", err)
			}
			for _, c := range tests {
				if !res.Passed[c.ID] {
					t.Errorf("the reference solution failed %q\noutput:\n%s", c.Label, res.Output)
				}
			}
			if res.TimedOut {
				t.Errorf("the reference solution ran out of %dms — the limit leaves no "+
					"headroom for a learner writing the same algorithm slightly less tersely",
					p.timeLimitMs)
			}
		})
	}
}

// Two problems tell the learner outright that the obvious approach is too slow.
// That is a promise about the time limit, and a limit generous enough to let a
// quadratic solution through turns the lesson into a lie.
func TestTheProblemsThatPromiseBruteForceIsTooSlowMeanIt(t *testing.T) {
	e := newPythonExecutor(t)
	bruteForce := map[string]string{
		"two-sum": `def two_sum(nums, target):
    for i in range(len(nums)):
        for j in range(i + 1, len(nums)):
            if nums[i] + nums[j] == target:
                return (i, j)
    return None
`,
		"group-anagrams": `def group_anagrams(words):
    groups = []
    for w in words:
        for g in groups:
            if sorted(g[0]) == sorted(w):
                g.append(w)
                break
        else:
            groups.append([w])
    return groups
`,
	}
	for slug, code := range bruteForce {
		t.Run(slug, func(t *testing.T) {
			t.Parallel()
			var p seedProblem
			for _, sp := range seedProblems {
				if sp.slug == slug {
					p = sp
				}
			}
			if p.slug == "" {
				t.Fatalf("no seeded problem %q — it was renamed or removed", slug)
			}
			tests := p.tests
			for i := range tests {
				tests[i].ID = int64(i + 1)
			}
			res, err := jobs.GradeWithin(t.Context(), e, "python", code, nil, tests, p.timeLimitMs)
			if err != nil {
				t.Fatalf("grade: %v", err)
			}
			passedAll := true
			for _, c := range tests {
				if !res.Passed[c.ID] {
					passedAll = false
				}
			}
			if passedAll {
				t.Errorf("the brute-force solution passed everything inside %dms, so the "+
					"statement's claim that it is too slow is not true", p.timeLimitMs)
			}
			// It should be the speed that stopped it, not a wrong answer —
			// brute force here is correct, only slow.
			if !res.TimedOut {
				t.Errorf("brute force failed for some reason other than the time limit\noutput:\n%s", res.Output)
			}
		})
	}
}

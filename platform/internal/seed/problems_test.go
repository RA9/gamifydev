package seed

import (
	"os/exec"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func newProblemExecutor(t *testing.T) runner.Executor {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("no python3 on PATH")
	}
	if _, err := exec.LookPath("cc"); err != nil {
		t.Skip("no C compiler on PATH")
	}
	e, err := runner.New(runner.Config{Mode: "local", PythonPath: "python3", CCPath: "cc", AllowUnsafe: true})
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
	e := newProblemExecutor(t)
	for _, p := range seedProblems {
		t.Run(p.slug, func(t *testing.T) {
			if p.solutionLang != "c" {
				t.Parallel()
			}
			if p.solution == "" {
				t.Fatalf("problem %q ships no reference solution, so nothing proves it can be solved", p.slug)
			}
			lang := p.solutionLang
			if lang == "" {
				lang = "python"
			}
			res, err := jobs.GradeWithin(t.Context(), e, lang, p.solution, nil, withIDs(p.tests), p.timeLimitMs)
			if err != nil {
				t.Fatalf("grade: %v", err)
			}
			for _, c := range withIDs(p.tests) {
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

// A problem whose statement says the obvious approach is too slow declares that
// approach beside it, and this holds the claim to account. A limit generous
// enough to let the slow approach through turns the lesson into a lie, and the
// learner who writes the clever solution never finds out it mattered.
func TestTheApproachesAProblemCallsTooSlowReallyAre(t *testing.T) {
	e := newProblemExecutor(t)
	claimed := 0
	for _, p := range seedProblems {
		if p.tooSlow == "" {
			continue
		}
		claimed++
		t.Run(p.slug, func(t *testing.T) {
			if p.solutionLang != "c" {
				t.Parallel()
			}
			tests := withIDs(p.tests)
			lang := p.solutionLang
			if lang == "" {
				lang = "python"
			}
			res, err := jobs.GradeWithin(t.Context(), e, lang, p.tooSlow, nil, tests, p.timeLimitMs)
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
				t.Errorf("the slow approach passed everything inside %dms, so the "+
					"statement's claim that it will not finish is not true", p.timeLimitMs)
			}
			// And it should be the clock that stopped it, not a wrong answer:
			// the approach the statement warns about is correct, only slow.
			if !res.TimedOut {
				t.Errorf("the slow approach failed for some reason other than the time "+
					"limit, so it is not the example the statement thinks it is\noutput:\n%s",
					res.Output)
			}
		})
	}
	if claimed == 0 {
		t.Error("no problem in the bank declares an approach that is too slow, so " +
			"nothing here is testing the time limits at all")
	}
}

// withIDs gives the checks the ids the store would have assigned, since Grade
// reports results by id.
func withIDs(checks []store.Check) []store.Check {
	out := make([]store.Check, len(checks))
	copy(out, checks)
	for i := range out {
		out[i].ID = int64(i + 1)
	}
	return out
}

// Slugs are the bank's public URLs and its upsert key. Two problems sharing one
// would mean the second silently overwrites the first at seed time — no error,
// just a problem that quietly stopped existing.
func TestEveryProblemHasItsOwnSlug(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range seedProblems {
		if seen[p.slug] {
			t.Errorf("two problems share the slug %q", p.slug)
		}
		seen[p.slug] = true
	}
}

// Content rules the bank relies on, checked here rather than discovered by a
// learner opening a half-written problem.
func TestEveryProblemIsCompletelyAuthored(t *testing.T) {
	for _, p := range seedProblems {
		t.Run(p.slug, func(t *testing.T) {
			if p.title == "" || p.statement == "" || p.topic == "" {
				t.Error("missing a title, statement or topic")
			}
			switch p.difficulty {
			case "easy", "medium", "hard":
			default:
				t.Errorf("difficulty %q is not one the database accepts", p.difficulty)
			}
			if p.timeLimitMs <= 0 {
				t.Error("no time limit, so the judge would use the runner's default")
			}
			if len(p.starters) == 0 {
				t.Error("no starter, so the editor would offer no language at all")
			}
			if len(p.tests) < 3 {
				t.Errorf("only %d test(s) — too few to distinguish a solution from a guess", len(p.tests))
			}
			hidden := 0
			for _, c := range p.tests {
				if c.Hidden {
					hidden++
				}
			}
			if hidden == 0 {
				t.Error("every test is visible, so a solution can be written to match the " +
					"examples without solving anything")
			}
			if hidden == len(p.tests) {
				t.Error("every test is hidden, so a learner is told nothing about what is checked")
			}
		})
	}
}

// The bank has to be big enough, and spread widely enough, to be worth opening
// twice. A single accidental deletion that halved it, or quietly left only the
// hard ones, would still pass every other test in this file.
func TestTheBankIsBigAndVariedEnoughToBrowse(t *testing.T) {
	byDifficulty := map[string]int{}
	byTopic := map[string]int{}
	for _, p := range seedProblems {
		byDifficulty[p.difficulty]++
		byTopic[p.topic]++
	}
	if len(seedProblems) < 80 {
		t.Errorf("only %d problems — too thin a bank to come back to", len(seedProblems))
	}
	// Somebody on their first week and somebody preparing for interviews have
	// to both find something here, or the bank serves neither.
	for _, d := range []string{"easy", "medium", "hard"} {
		if byDifficulty[d] < 10 {
			t.Errorf("only %d %s problems — not enough for anyone working at that level", byDifficulty[d], d)
		}
	}
	if len(byTopic) < 10 {
		t.Errorf("only %d topics — the bank teaches too narrow a range", len(byTopic))
	}
}

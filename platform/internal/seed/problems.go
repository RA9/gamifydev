package seed

import (
	"context"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// The practice problem bank.
//
// Problems remain open practice that gates nothing and can be attempted without
// signing in. A curated subset is also linked from courses so practical work sits
// beside the theory it reinforces; the public bank remains available to everyone.
//
// **They are also deliberately not the canon.** "Two sum", "valid parentheses"
// and "group anagrams" have a worked solution on the first page of every search
// and in the training data of every assistant, so setting them measures whether
// someone recognised the problem, not whether they can solve one.
//
// So the problems are written rather than borrowed: each states a small rule
// system, and the rules are the specification. Most of the bank has no
// counterpart to look up at all.
//
// A dozen or so are the exception, and honestly so. Some techniques cannot be
// taught around — breadth-first search, topological order, a monotonic stack, a
// sweep over interval endpoints, Dijkstra — and inventing an unrecognisable
// dressing for them would make the problem worse, not more original. Those are
// varied instead of disguised: clusters join diagonally, the metro searches over
// lines rather than stations, tickets escalate at 48 hours, brackets inside
// quotes are text. Each variation is covered by a hidden test, so a canonical
// solution pasted in unread fails on the clause it does not know about — which
// is the outcome worth engineering, since it is also what happens to somebody
// who skims the statement.
//
// Difficulty is set by what a problem asks you to notice, not by how much code
// it takes. A problem whose statement warns that the obvious approach is too
// slow declares that approach in tooSlow, and a test proves the time limit
// actually rejects it.
type seedProblem struct {
	slug, title, difficulty, topic string
	statement                      string
	timeLimitMs                    int
	starters                       map[string]string
	tests                          []store.Check
	workloadMinutes                int
	mode, language, solutionLang   string
	// solution is a correct answer, kept beside the problem so a test can prove
	// the problem is solvable rather than leaving a learner to discover it
	// isn't. Never served to anyone.
	solution string
	// solutions is the same guarantee for a problem offered in more than one
	// language: one reference answer per language, each proven by the suite. A
	// problem that claims to accept C, Go and Python without anybody having
	// written all three is a promise nobody checked.
	solutions map[string]string
	// tooSlow is an approach the statement claims won't finish in time. Set it
	// on any problem that makes that claim, and a test holds the claim to
	// account — a limit generous enough to let the slow approach through turns
	// the lesson into a lie.
	tooSlow string
	// tooSlowIn is tooSlow per language. The same limit does not mean the same
	// thing in C as in Python — a quadratic loop that dies in one finishes
	// comfortably in the other — so the claim is made, and checked, only for
	// the languages where it actually holds.
	tooSlowIn map[string]string
}

// Test constructors. The bank is thousands of lines of content, and
// store.Check{Label: ..., Test: ..., Points: 1} repeated six hundred times
// buries the problems in punctuation.
//
// vis is a test the learner can see on the problem page; hid is one they can't.
// Hidden tests are what stop a solution that pattern-matches the examples from
// passing as an understanding of the rules.
func vis(label, test string) store.Check {
	return store.Check{Label: label, Test: test, Points: 1}
}

func hid(label, test string) store.Check {
	return store.Check{Label: label, Test: test, Points: 1, Hidden: true}
}

// py and cStarter declare the languages a problem accepts.
func py(starter string) map[string]string {
	return map[string]string{"python": starter}
}

// tri is the starter set for a problem answerable in any of the three languages
// the bank teaches. Each starter carries the input parsing and the printing, so
// what the learner writes is the part the problem is actually about.
func tri(c, golang, python string) map[string]string {
	return map[string]string{"c": c, "go": golang, "python": python}
}

// tokens asserts on what the program printed, compared as whitespace-separated
// tokens.
//
// Comparing tokens rather than exact bytes is deliberate. Three languages print
// the same answer with different habits — a trailing newline, numbers joined by
// a space or by newlines — and none of those differences is what is being
// taught. What has to match is the answer.
func tokens(label, stdin, want string) store.Check {
	return store.Check{
		Label: label, Stdin: stdin, Points: 1,
		Test: "_out.split() == " + pyList(strings.Fields(want)),
	}
}

// hid2 is tokens for a case the learner cannot see.
func hid2(label, stdin, want string) store.Check {
	c := tokens(label, stdin, want)
	c.Hidden = true
	return c
}

// exact asserts on the whole printed text, for the answers where the spacing is
// part of the answer and tokens would quietly forgive losing it.
func exact(label, stdin, want string) store.Check {
	return store.Check{
		Label: label, Stdin: stdin, Points: 1,
		Test: "_out.rstrip(\"\\n\") == " + strconv.Quote(want),
	}
}

// exactHid is exact for a case the learner cannot see.
func exactHid(label, stdin, want string) store.Check {
	c := exact(label, stdin, want)
	c.Hidden = true
	return c
}

// pyList renders strings as a Python list literal for a check expression.
func pyList(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = strconv.Quote(s)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func cStarter(starter string) map[string]string {
	return map[string]string{"c": starter}
}

// seedProblems is the whole bank, assembled from the themed groups.
//
// Split by the skill a problem exercises rather than by difficulty, because
// that is how they are written and reviewed: a gap in the sliding-window group
// is obvious when the sliding-window problems sit together, and invisible when
// they are scattered through a list sorted by how hard they are.
var seedProblems = func() []seedProblem {
	var all []seedProblem
	for _, group := range [][]seedProblem{
		basicProblems,
		stringProblems,
		tableProblems,
		orderProblems,
		windowProblems,
		stackProblems,
		intervalProblems,
		graphProblems,
		foundationProblems,
	} {
		all = append(all, group...)
	}
	return all
}()

// inlineCode turns the bank's ~name~ marker into markdown inline code.
//
// Go raw strings cannot contain a backtick, and a hundred statements written as
// quoted strings spliced around `+ "`" +` would be unreadable — which, for a
// file that is almost entirely prose a learner will read, is the thing most
// worth avoiding. Single tildes carry no meaning in GFM (strikethrough needs
// two), and every statement here is ours, so the substitution is unambiguous.
func inlineCode(s string) string {
	return strings.ReplaceAll(s, "~", "`")
}

// seedProblemBank upserts the bank, returning how many problems it wrote.
//
// Idempotent, like the rest of the seed: content is replaced wholesale by slug,
// so editing a problem here and re-running picks up the edit instead of
// stacking a second copy beside it.
func seedProblemBank(ctx context.Context, st *store.Store) (int, error) {
	for i, p := range seedProblems {
		language := p.language
		if language == "" {
			language = "python"
		}
		id, err := st.UpsertProblem(ctx, store.Problem{
			Slug: p.slug, Title: p.title, Difficulty: p.difficulty, Topic: p.topic,
			Statement: inlineCode(p.statement), TimeLimitMs: p.timeLimitMs, Sort: i, Published: true,
			WorkloadMinutes: p.workloadMinutes, Mode: p.mode, Language: language,
		})
		if err != nil {
			return 0, err
		}
		if err := st.ReplaceStarters(ctx, id, p.starters); err != nil {
			return 0, err
		}
		if err := st.ReplaceProblemTests(ctx, id, p.tests); err != nil {
			return 0, err
		}
	}
	return len(seedProblems), nil
}

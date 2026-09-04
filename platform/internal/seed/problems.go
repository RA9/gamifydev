package seed

import (
	"context"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// The practice problem bank.
//
// These are deliberately not coursework. A problem gates nothing, belongs to no
// path, and can be attempted by someone who has never signed in — it exists so
// a visitor can find out in ninety seconds whether they like writing code here,
// and so a learner between lessons has something to sharpen on.
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
	// solution is a correct answer, kept beside the problem so a test can prove
	// the problem is solvable rather than leaving a learner to discover it
	// isn't. Never served to anyone.
	solution string
	// tooSlow is an approach the statement claims won't finish in time. Set it
	// on any problem that makes that claim, and a test holds the claim to
	// account — a limit generous enough to let the slow approach through turns
	// the lesson into a lie.
	tooSlow string
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

// py is the starter set for a Python-only problem, which is all of them today.
func py(starter string) map[string]string {
	return map[string]string{"python": starter}
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
		id, err := st.UpsertProblem(ctx, store.Problem{
			Slug: p.slug, Title: p.title, Difficulty: p.difficulty, Topic: p.topic,
			Statement: inlineCode(p.statement), TimeLimitMs: p.timeLimitMs, Sort: i, Published: true,
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

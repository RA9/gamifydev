package seed

import (
	"context"

	"github.com/RA9/gamifydev/platform/internal/store"
)

// The practice problem bank.
//
// These are deliberately not coursework. A problem gates nothing, belongs to no
// path, and can be attempted by someone who has never signed in — it exists so
// a visitor can find out in ninety seconds whether they like writing code here,
// and so a learner between lessons has something to sharpen on.
//
// Difficulty is set by what the problem asks you to notice, not by how much
// code it takes. "Two sum" is medium not because the loop is long but because
// the honest answer requires seeing why the obvious one is too slow — and its
// time limit is set so that the obvious one isn't.
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
}

var seedProblems = []seedProblem{
	{
		slug: "sum-a-list", title: "Sum a list", difficulty: "easy", topic: "Arrays",
		statement: `Write a function ` + "`total(nums)`" + ` that returns the sum of a list of
whole numbers.

An empty list sums to ` + "`0`" + ` — that isn't a special case to code around, it's
what "add up nothing" means.

    total([1, 2, 3])  ->  6
    total([])         ->  0`,
		timeLimitMs: 5000,
		starters: map[string]string{
			"python": "def total(nums):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Adds up a short list", Test: "total([1, 2, 3]) == 6", Points: 1},
			{Label: "An empty list sums to zero", Test: "total([]) == 0", Points: 1},
			{Label: "Handles negative numbers", Test: "total([-4, 10, -6]) == 0", Points: 1},
			{Label: "Handles a long list", Test: "total(list(range(1001))) == 500500", Points: 1, Hidden: true},
		},
		solution: "def total(nums):\n    return sum(nums)\n",
	},
	{
		slug: "reverse-words", title: "Reverse the words", difficulty: "easy", topic: "Strings",
		statement: `Write ` + "`reverse_words(s)`" + `, which returns the words of ` + "`s`" + ` in the
opposite order, separated by single spaces.

    reverse_words("the quick brown fox")  ->  "fox brown quick the"

Runs of whitespace collapse, and leading or trailing whitespace disappears —
"the words" is a list, not a slice of the original string.`,
		timeLimitMs: 5000,
		starters: map[string]string{
			"python": "def reverse_words(s):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Reverses a simple sentence", Test: `reverse_words("the quick brown fox") == "fox brown quick the"`, Points: 1},
			{Label: "A single word is unchanged", Test: `reverse_words("hello") == "hello"`, Points: 1},
			{Label: "An empty string stays empty", Test: `reverse_words("") == ""`, Points: 1},
			{Label: "Collapses extra whitespace", Test: `reverse_words("  a   b  ") == "b a"`, Points: 1, Hidden: true},
		},
		solution: "def reverse_words(s):\n    return ' '.join(reversed(s.split()))\n",
	},
	{
		slug: "count-characters", title: "Count the characters", difficulty: "easy", topic: "Dictionaries",
		statement: `Write ` + "`counts(s)`" + `, returning a dictionary mapping each character in
` + "`s`" + ` to the number of times it appears.

    counts("aab")  ->  {"a": 2, "b": 1}

Characters that never appear are simply absent from the result — a key with a
count of zero would be a claim about a character that isn't there.`,
		timeLimitMs: 5000,
		starters: map[string]string{
			"python": "def counts(s):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Counts repeated characters", Test: `counts("aab") == {"a": 2, "b": 1}`, Points: 1},
			{Label: "An empty string has no counts", Test: `counts("") == {}`, Points: 1},
			{Label: "Counts spaces too", Test: `counts("a b")[" "] == 1`, Points: 1},
			{Label: "Absent characters have no key", Test: `"z" not in counts("aab")`, Points: 1, Hidden: true},
		},
		solution: "def counts(s):\n    out = {}\n    for ch in s:\n        out[ch] = out.get(ch, 0) + 1\n    return out\n",
	},
	{
		slug: "two-sum", title: "Two sum", difficulty: "medium", topic: "Hash maps",
		statement: `Write ` + "`two_sum(nums, target)`" + `, returning the indices of the two
numbers in ` + "`nums`" + ` that add up to ` + "`target`" + `, as a tuple ` + "`(i, j)`" + ` with
` + "`i < j`" + `. Exactly one such pair exists.

    two_sum([2, 7, 11, 15], 9)  ->  (0, 1)

The obvious answer — try every pair — is correct and too slow: the last test
feeds it a large list, and the time limit is set so that checking every pair
won't finish. What you need is a way to ask "have I already seen the number
that would complete this one?" without searching for it.`,
		timeLimitMs: 4000,
		starters: map[string]string{
			"python": "def two_sum(nums, target):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Finds the pair in a small list", Test: "two_sum([2, 7, 11, 15], 9) == (0, 1)", Points: 1},
			{Label: "Finds a pair that isn't at the start", Test: "two_sum([3, 2, 4], 6) == (1, 2)", Points: 1},
			{Label: "Handles a repeated number", Test: "two_sum([3, 3], 6) == (0, 1)", Points: 1},
			{
				// The whole point of the problem. Named so that a learner whose
				// nested loop times out is told which requirement it was.
				Label:  "Fast enough on a large list",
				Test:   "two_sum(list(range(200000)), 399997) == (199998, 199999)",
				Points: 1,
			},
		},
		solution: "def two_sum(nums, target):\n    seen = {}\n    for i, n in enumerate(nums):\n        if target - n in seen:\n            return (seen[target - n], i)\n        seen[n] = i\n    return None\n",
	},
	{
		slug: "balanced-brackets", title: "Balanced brackets", difficulty: "medium", topic: "Stacks",
		statement: `Write ` + "`balanced(s)`" + `, returning ` + "`True`" + ` if every bracket in
` + "`s`" + ` is closed by the matching kind, in the right order.

    balanced("({[]})")  ->  True
    balanced("([)]")    ->  False

Counting brackets is not enough: ` + "`([)]`" + ` has two of each and is still wrong.
What matters is which bracket is still waiting to be closed, and that is always
the most recent one.`,
		timeLimitMs: 5000,
		starters: map[string]string{
			"python": "def balanced(s):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Accepts properly nested brackets", Test: `balanced("({[]})") is True`, Points: 1},
			{Label: "Rejects brackets closed in the wrong order", Test: `balanced("([)]") is False`, Points: 1},
			{Label: "An empty string is balanced", Test: `balanced("") is True`, Points: 1},
			{Label: "Rejects an unclosed bracket", Test: `balanced("(((") is False`, Points: 1},
			{Label: "Rejects a stray closing bracket", Test: `balanced(")") is False`, Points: 1, Hidden: true},
		},
		solution: "def balanced(s):\n    pairs = {')': '(', ']': '[', '}': '{'}\n    stack = []\n    for ch in s:\n        if ch in '([{':\n            stack.append(ch)\n        elif ch in pairs:\n            if not stack or stack.pop() != pairs[ch]:\n                return False\n    return not stack\n",
	},
	{
		slug: "longest-run", title: "Longest run", difficulty: "medium", topic: "Arrays",
		statement: `Write ` + "`longest_run(nums)`" + `, returning the length of the longest
stretch of equal values sitting next to each other.

    longest_run([1, 1, 2, 2, 2, 3])  ->  3
    longest_run([])                  ->  0

One pass is enough. If you find yourself starting over from each position, you
are re-counting stretches you have already walked.`,
		timeLimitMs: 4000,
		starters: map[string]string{
			"python": "def longest_run(nums):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Finds a run in the middle", Test: "longest_run([1, 1, 2, 2, 2, 3]) == 3", Points: 1},
			{Label: "An empty list has no run", Test: "longest_run([]) == 0", Points: 1},
			{Label: "All-distinct values give a run of one", Test: "longest_run([1, 2, 3]) == 1", Points: 1},
			{Label: "Finds a run at the very end", Test: "longest_run([5, 1, 1, 1]) == 3", Points: 1},
			{Label: "Fast enough on a long list", Test: "longest_run([0] * 300000) == 300000", Points: 1, Hidden: true},
		},
		solution: "def longest_run(nums):\n    best = run = 0\n    prev = object()\n    for n in nums:\n        run = run + 1 if n == prev else 1\n        prev = n\n        if run > best:\n            best = run\n    return best\n",
	},
	{
		slug: "group-anagrams", title: "Group anagrams", difficulty: "hard", topic: "Hash maps",
		statement: `Write ` + "`group_anagrams(words)`" + `, which groups words that are
rearrangements of each other.

Return a list of groups. Each group is a list of words in the order they
appeared in the input, and the groups themselves are in the order their first
word appeared.

    group_anagrams(["eat", "tea", "tan", "ate", "nat", "bat"])
      ->  [["eat", "tea", "ate"], ["tan", "nat"], ["bat"]]

Comparing every word against every other word is too slow for the last test.
The trick is to find something that is *the same* for any two anagrams, and use
it as a key.`,
		timeLimitMs: 5000,
		starters: map[string]string{
			"python": "def group_anagrams(words):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{
				Label:  "Groups the classic example",
				Test:   `group_anagrams(["eat","tea","tan","ate","nat","bat"]) == [["eat","tea","ate"],["tan","nat"],["bat"]]`,
				Points: 1,
			},
			{Label: "No words, no groups", Test: "group_anagrams([]) == []", Points: 1},
			{Label: "Words with no anagram stand alone", Test: `group_anagrams(["abc","def"]) == [["abc"],["def"]]`, Points: 1},
			{
				Label:  "Keeps groups in first-appearance order",
				Test:   `group_anagrams(["bat","tab","cat"]) == [["bat","tab"],["cat"]]`,
				Points: 1, Hidden: true,
			},
			{
				// 12000 distinct words that fall into 880 anagram groups. Both
				// numbers matter: a solution that compares each word against
				// every group it has built so far does ~12000x880 comparisons
				// here and takes several times the limit, while grouping by a
				// key takes milliseconds.
				//
				// An input of all-identical words would not test this at all —
				// the scanning solution would match on the first group every
				// time and sail through.
				Label:  "Fast enough on many words",
				Test:   `len(group_anagrams(["w%05d" % i for i in range(12000)])) == 880`,
				Points: 1,
			},
		},
		solution: "def group_anagrams(words):\n    groups = {}\n    for w in words:\n        groups.setdefault(''.join(sorted(w)), []).append(w)\n    return list(groups.values())\n",
	},
	{
		slug: "run-length-encode", title: "Run-length encoding", difficulty: "hard", topic: "Strings",
		statement: `Write ` + "`encode(s)`" + `, compressing runs of the same character into the
character followed by its count — but only when doing so is actually shorter.

    encode("aaabbc")  ->  "a3b2c"
    encode("abc")     ->  "abc"

A run of one is written as the bare character: ` + "`c1`" + ` is longer than ` + "`c`" + `,
and a compressor that grows its input isn't one.`,
		timeLimitMs: 4000,
		starters: map[string]string{
			"python": "def encode(s):\n    # your code here\n    pass\n",
		},
		tests: []store.Check{
			{Label: "Compresses runs", Test: `encode("aaabbc") == "a3b2c"`, Points: 1},
			{Label: "Leaves single characters bare", Test: `encode("abc") == "abc"`, Points: 1},
			{Label: "An empty string encodes to nothing", Test: `encode("") == ""`, Points: 1},
			{Label: "Handles one long run", Test: `encode("a" * 12) == "a12"`, Points: 1},
			{Label: "Handles a run at the end", Test: `encode("abbb") == "ab3"`, Points: 1, Hidden: true},
			{Label: "Fast enough on a long string", Test: `encode("ab" * 100000) == "ab" * 100000`, Points: 1, Hidden: true},
		},
		solution: "def encode(s):\n    out = []\n    i = 0\n    while i < len(s):\n        j = i\n        while j < len(s) and s[j] == s[i]:\n            j += 1\n        n = j - i\n        out.append(s[i] if n == 1 else s[i] + str(n))\n        i = j\n    return ''.join(out)\n",
	},
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
			Statement: p.statement, TimeLimitMs: p.timeLimitMs, Sort: i, Published: true,
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

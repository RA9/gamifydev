package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// stringProblems: reading a string one character at a time, and the state you
// have to carry while you do it.
var stringProblems = []seedProblem{
	{
		slug: "emphasis-pass", title: "Emphasis pass", difficulty: "easy", topic: "Strings",
		statement: `A chat client renders emphasis by case. A word ending in ~!~ is shouted
(uppercased); a word ending in ~?~ is muttered (lowercased); everything else is
left exactly as it was.

Write ~emphasise(text)~. Words are separated by single spaces, and the
punctuation stays on the word.

    emphasise("Look out! is it Safe?")  ->  "Look OUT! is it safe?"`,
		timeLimitMs: 5000,
		starters:    py("def emphasise(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Shouts and mutters in one pass", `emphasise("Look out! is it Safe?") == "Look OUT! is it safe?"`),
			vis("Plain words are untouched", `emphasise("nothing to see") == "nothing to see"`),
			vis("An empty message stays empty", `emphasise("") == ""`),
			hid("Punctuation in the middle of a word does not count", `emphasise("wait!for it") == "wait!for it"`),
			hid("A word that is only punctuation still works", `emphasise("well ! ?") == "well ! ?"`),
		},
		solution: "def emphasise(text):\n    out = []\n    for w in text.split(' '):\n        if w.endswith('!'):\n            out.append(w.upper())\n        elif w.endswith('?'):\n            out.append(w.lower())\n        else:\n            out.append(w)\n    return ' '.join(out)\n",
	},
	{
		slug: "redact-long-numbers", title: "Redact long numbers", difficulty: "medium", topic: "Strings",
		statement: `A support tool has to hide anything that looks like an account number
before a transcript is shared. A run of **six or more** digits is replaced by
the same number of ~#~ characters. Shorter runs — years, quantities, door
numbers — are left alone.

Write ~redact(text)~.

    redact("call 5551234 or ext 12345")  ->  "call ####### or ext 12345"

The replacement is the same length as what it hides, so the transcript still
lines up.`,
		timeLimitMs: 5000,
		starters:    py("def redact(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Hides a long run and keeps a short one", `redact("call 5551234 or ext 12345") == "call ####### or ext 12345"`),
			vis("A year is not an account number", `redact("since 1999") == "since 1999"`),
			vis("Nothing to redact", `redact("no digits here") == "no digits here"`),
			vis("Exactly six digits is long enough", `redact("id 123456") == "id ######"`),
			hid("Five digits is not", `redact("id 12345") == "id 12345"`),
			hid("A run at the very end is caught", `redact("ref 9876543") == "ref #######"`),
			hid("Two long runs in one line", `redact("1234567 and 7654321") == "####### and #######"`),
		},
		solution: "def redact(text):\n    out = []\n    i = 0\n    while i < len(text):\n        if text[i].isdigit():\n            j = i\n            while j < len(text) and text[j].isdigit():\n                j += 1\n            run = j - i\n            out.append('#' * run if run >= 6 else text[i:j])\n            i = j\n        else:\n            out.append(text[i])\n            i += 1\n    return ''.join(out)\n",
	},
	{
		slug: "badge-initials", title: "Badge initials", difficulty: "medium", topic: "Strings",
		statement: `A conference badge shows a person's initials. Name particles — ~van~,
~der~, ~de~, ~den~, ~da~, ~di~, ~bin~, ~al~ — are skipped, but only when they
are written in lower case. A capitalised ~Van~ is somebody's actual name and
gets an initial like any other word.

Write ~initials(name)~ returning the uppercase initials with nothing between
them.

    initials("ada van der lovelace")  ->  "AL"
    initials("Vincent Van Gogh")      ->  "VVG"

Words are separated by whitespace, and there may be more than one space between
them.`,
		timeLimitMs: 5000,
		starters:    py("def initials(name):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Skips lower-case particles", `initials("ada van der lovelace") == "AL"`),
			vis("A capitalised particle is a name", `initials("Vincent Van Gogh") == "VVG"`),
			vis("An ordinary name is untouched", `initials("Grace Hopper") == "GH"`),
			vis("Nobody has no initials", `initials("") == ""`),
			hid("Extra spaces do not create empty initials", `initials("  Ada   Lovelace ") == "AL"`),
			hid("A particle-looking word that is not one", `initials("Alan Turing") == "AT"`),
		},
		solution: "PARTICLES = {'van', 'der', 'de', 'den', 'da', 'di', 'bin', 'al'}\n\ndef initials(name):\n    out = []\n    for w in name.split():\n        if w.islower() and w in PARTICLES:\n            continue\n        out.append(w[0].upper())\n    return ''.join(out)\n",
	},
	{
		slug: "ticker-window", title: "Ticker window", difficulty: "medium", topic: "Strings",
		statement: `A departure board scrolls a message right to left through a narrow
window. The message repeats forever with exactly three spaces between the end
of one pass and the start of the next.

Write ~visible(text, width, t)~ returning the ~width~ characters showing at
tick ~t~. At tick 0 the window starts at the first character of the message;
each tick moves the message one place left.

    visible("HELLO", 3, 0)  ->  "HEL"
    visible("HELLO", 3, 6)  ->  "  H"

~t~ can be enormous, and the window can be wider than the message itself.`,
		timeLimitMs: 5000,
		starters:    py("def visible(text, width, t):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("The start of the message", `visible("HELLO", 3, 0) == "HEL"`),
			vis("Scrolled into the gap and round again", `visible("HELLO", 3, 6) == "  H"`),
			vis("A window wider than the message", `visible("AB", 6, 0) == "AB   A"`),
			hid("A tick far in the future", `visible("HELLO", 5, 8 * 10 ** 9) == "HELLO"`),
			hid("A window of nothing shows nothing", `visible("HELLO", 0, 3) == ""`),
		},
		solution: "def visible(text, width, t):\n    loop = text + '   '\n    return ''.join(loop[(t + i) % len(loop)] for i in range(width))\n",
	},
	{
		slug: "message-segments", title: "Message segments", difficulty: "medium", topic: "Strings",
		statement: `A messaging gateway bills by segment. Most characters cost one unit, but
any character in the ~special~ string costs two. A message of 160 units or
fewer is a single segment; anything longer is split into segments of 153 units,
and a part-filled final segment still counts.

Write ~segments(text, special)~ returning how many segments are billed. An
empty message is not sent at all, so it costs nothing.

    segments("hello", "")     ->  1
    segments("a" * 161, "")   ->  2`,
		timeLimitMs: 5000,
		starters:    py("def segments(text, special):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A short message is one segment", `segments("hello", "") == 1`),
			vis("One unit over splits it", `segments("a" * 161, "") == 2`),
			vis("An empty message is never sent", `segments("", "") == 0`),
			vis("Special characters cost double", `segments("€" * 81, "€") == 2`),
			hid("Exactly 160 units still fits in one", `segments("a" * 160, "") == 1`),
			hid("A part-filled final segment counts", `segments("a" * 307, "") == 3`),
			hid("Only the listed characters cost double", `segments("€$", "$") == 3 - 2`),
		},
		solution: "def segments(text, special):\n    n = sum(2 if c in special else 1 for c in text)\n    if n == 0:\n        return 0\n    if n <= 160:\n        return 1\n    return -(-n // 153)\n",
	},
	{
		slug: "sentence-caps", title: "Sentence caps", difficulty: "medium", topic: "Strings",
		statement: `Write ~sentence_caps(text)~, which capitalises the first letter of every
sentence and changes nothing else — not the middle of words, not the case of
anything already written.

A sentence ends at ~.~, ~!~ or ~?~. The next letter after one of those, however
far away it is, starts the next sentence.

    sentence_caps("hello there. how are you? fine!")
      ->  "Hello there. How are you? Fine!"

Characters that are not letters never get capitalised, and never stop the
search for the letter that does.`,
		timeLimitMs: 5000,
		starters:    py("def sentence_caps(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Capitalises each sentence", `sentence_caps("hello there. how are you? fine!") == "Hello there. How are you? Fine!"`),
			vis("Leaves the rest of the words alone", `sentence_caps("iPhone and iPad.") == "IPhone and iPad."`),
			vis("An empty string stays empty", `sentence_caps("") == ""`),
			hid("Skips past spaces and quotes to find the letter", `sentence_caps("stop.  'go' now") == "Stop.  'Go' now"`),
			hid("An already capital letter is left as it is", `sentence_caps("One. Two.") == "One. Two."`),
			hid("Trailing punctuation with nothing after it", `sentence_caps("done.") == "Done."`),
		},
		solution: "def sentence_caps(text):\n    out = list(text)\n    start = True\n    for i, ch in enumerate(out):\n        if start and ch.isalpha():\n            out[i] = ch.upper()\n            start = False\n        elif ch in '.!?':\n            start = True\n    return ''.join(out)\n",
	},
	{
		slug: "strip-markup", title: "Strip markup", difficulty: "medium", topic: "Strings",
		statement: `Write ~strip_tags(text)~, removing every ~<...>~ span along with its angle
brackets.

The catch is the stray bracket. A ~<~ with no ~>~ anywhere after it is not the
start of a tag — it is a less-than sign somebody typed, and it stays.

    strip_tags("a <b>bold</b> move")  ->  "a bold move"
    strip_tags("5 < 6")               ->  "5 < 6"`,
		timeLimitMs: 4000,
		starters:    py("def strip_tags(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Removes tags", `strip_tags("a <b>bold</b> move") == "a bold move"`),
			vis("A stray bracket stays", `strip_tags("5 < 6") == "5 < 6"`),
			vis("Nothing to strip", `strip_tags("plain") == "plain"`),
			vis("The whole string is one tag", `strip_tags("<hr>") == ""`),
			hid("A tag ends at the first closing bracket after it", `strip_tags("a < b <i>c</i>") == "a c"`),
			hid("A stray bracket after a real tag", `strip_tags("a <b> c < d") == "a  c < d"`),
			hid("A closing bracket on its own is ordinary text", `strip_tags("3 > 2") == "3 > 2"`),
		},
		solution: "def strip_tags(text):\n    out = []\n    i = 0\n    while i < len(text):\n        if text[i] == '<':\n            j = text.find('>', i)\n            if j == -1:\n                out.append(text[i])\n                i += 1\n            else:\n                i = j + 1\n        else:\n            out.append(text[i])\n            i += 1\n    return ''.join(out)\n",
	},
	{
		slug: "one-transposition", title: "One transposition", difficulty: "medium", topic: "Strings",
		statement: `A spell checker wants to know whether a typo is just two neighbouring
letters typed the wrong way round.

Write ~one_swap(a, b)~ returning ~True~ only if swapping exactly one pair of
**adjacent** characters in ~a~ produces ~b~.

    one_swap("form", "from")  ->  True
    one_swap("form", "form")  ->  False

A word is not a transposition of itself: the swap has to actually change
something.`,
		timeLimitMs: 5000,
		starters:    py("def one_swap(a, b):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Two neighbours swapped", `one_swap("form", "from") is True`),
			vis("A word is not a swap of itself", `one_swap("form", "form") is False`),
			vis("Different lengths cannot be a swap", `one_swap("form", "forms") is False`),
			vis("Two changes that are not a swap", `one_swap("cat", "dog") is False`),
			hid("Non-adjacent letters swapped does not count", `one_swap("abcd", "dbca") is False`),
			hid("Swapping two identical letters changes nothing", `one_swap("aab", "aab") is False`),
			hid("A swap at the very end", `one_swap("teh", "the") is True`),
		},
		solution: "def one_swap(a, b):\n    if len(a) != len(b) or a == b:\n        return False\n    diff = [i for i in range(len(a)) if a[i] != b[i]]\n    return (len(diff) == 2 and diff[1] == diff[0] + 1\n            and a[diff[0]] == b[diff[1]] and a[diff[1]] == b[diff[0]])\n",
	},
	{
		slug: "slug-registry", title: "Slug registry", difficulty: "hard", topic: "Strings",
		statement: `A blog turns every post title into a URL slug. The first post to claim a
slug keeps it; each later post claiming the same one gets ~-2~, then ~-3~, and
so on appended.

Write ~assign(names)~ returning the slug each post ends up with, in order.

    assign(["post", "post", "post"])  ->  ["post", "post-2", "post-3"]

The hard part is the collision you did not cause. If a post is *called*
~post-2~, that slug is already taken by the time the second ~post~ needs it —
and whatever you hand out has to be free too.

    assign(["post", "post-2", "post"])  ->  ["post", "post-2", "post-2"]

is wrong. No slug may ever be issued twice.`,
		timeLimitMs: 5000,
		starters:    py("def assign(names):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Numbers the repeats", `assign(["post", "post", "post"]) == ["post", "post-2", "post-3"]`),
			vis("Distinct titles keep their slugs", `assign(["a", "b"]) == ["a", "b"]`),
			vis("Nothing to assign", `assign([]) == []`),
			vis("Steps over a slug somebody already took", `assign(["post", "post-2", "post"]) == ["post", "post-2", "post-3"]`),
			hid("A taken slug pushes the next one further", `assign(["p", "p-2", "p-3", "p"]) == ["p", "p-2", "p-3", "p-4"]`),
			hid("A collision on the generated slug too", `assign(["p", "p", "p-2"]) == ["p", "p-2", "p-2-2"]`),
		},
		solution: "def assign(names):\n    used = set()\n    counts = {}\n    out = []\n    for n in names:\n        if n not in used:\n            used.add(n)\n            counts[n] = 1\n            out.append(n)\n            continue\n        k = counts.get(n, 1)\n        while True:\n            k += 1\n            cand = n + '-' + str(k)\n            if cand not in used:\n                break\n        counts[n] = k\n        used.add(cand)\n        out.append(cand)\n    return out\n",
	},
	{
		slug: "letter-drift", title: "Letter drift", difficulty: "medium", topic: "Strings",
		statement: `A toy cipher shifts each letter further than the last: the first letter
moves 0 places along the alphabet, the second 1, the third 2, and so on,
wrapping round from ~z~ to ~a~. Case is kept.

Characters that are not letters pass through untouched **and do not advance the
shift** — the count follows the letters, not the positions.

Write ~drift(text)~.

    drift("abc")     ->  "ace"
    drift("a-bc")    ->  "a-ce"`,
		timeLimitMs: 5000,
		starters:    py("def drift(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Each letter drifts one further", `drift("abc") == "ace"`),
			vis("Punctuation does not advance the shift", `drift("a-bc") == "a-ce"`),
			vis("An empty string stays empty", `drift("") == ""`),
			vis("Case is kept", `drift("AB") == "AC"`),
			hid("Wraps round the end of the alphabet", `drift("zz") == "za"`),
			hid("The first letter never moves", `drift("q") == "q"`),
		},
		solution: "def drift(text):\n    out = []\n    k = 0\n    for ch in text:\n        if ch.isalpha():\n            base = ord('A') if ch.isupper() else ord('a')\n            out.append(chr((ord(ch) - base + k) % 26 + base))\n            k += 1\n        else:\n            out.append(ch)\n    return ''.join(out)\n",
	},
	{
		slug: "readable-list", title: "Readable list", difficulty: "easy", topic: "Strings",
		statement: `Write ~phrase(names)~, joining names the way a sentence would: commas
between all but the last two, and ~and~ before the last.

    phrase(["Ada"])                 ->  "Ada"
    phrase(["Ada", "Grace"])        ->  "Ada and Grace"
    phrase(["Ada", "Grace", "Kay"]) ->  "Ada, Grace and Kay"

An empty list has nothing to say.`,
		timeLimitMs: 5000,
		starters:    py("def phrase(names):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("One name stands alone", `phrase(["Ada"]) == "Ada"`),
			vis("Two names take an and", `phrase(["Ada", "Grace"]) == "Ada and Grace"`),
			vis("Three names take commas and an and", `phrase(["Ada", "Grace", "Kay"]) == "Ada, Grace and Kay"`),
			vis("Nobody says nothing", `phrase([]) == ""`),
			hid("A long list keeps one and", `phrase(["a", "b", "c", "d"]) == "a, b, c and d"`),
		},
		solution: "def phrase(names):\n    if not names:\n        return ''\n    if len(names) == 1:\n        return names[0]\n    return ', '.join(names[:-1]) + ' and ' + names[-1]\n",
	},
	{
		slug: "column-widths", title: "Column widths", difficulty: "medium", topic: "Strings",
		statement: `Before a table can be printed, each column needs a width: the length of
the longest cell in it.

Write ~widths(rows)~, where ~rows~ is a list of lists of strings. Rows may be
ragged — a row that stops early simply has nothing in the later columns, and
contributes no width to them.

    widths([["name", "id"], ["ada"]])  ->  [4, 2]

The result has one entry per column, counting the widest row.`,
		timeLimitMs: 5000,
		starters:    py("def widths(rows):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Takes the longest cell per column", `widths([["name", "id"], ["ada"]]) == [4, 2]`),
			vis("A single row is its own widths", `widths([["abc", "d"]]) == [3, 1]`),
			vis("No rows, no columns", "widths([]) == []"),
			hid("A later row can be the wider one", `widths([["a"], ["bbb"]]) == [3]`),
			hid("An empty cell has width zero", `widths([["", "xx"]]) == [0, 2]`),
			hid("A ragged row adds columns", `widths([["a"], ["b", "cc"]]) == [1, 2]`),
		},
		solution: "def widths(rows):\n    n = max((len(r) for r in rows), default=0)\n    return [max((len(r[i]) for r in rows if i < len(r)), default=0) for i in range(n)]\n",
	},
	{
		slug: "keysmash-filter", title: "Keysmash filter", difficulty: "medium", topic: "Strings",
		statement: `A signup form rejects names that look like somebody leaned on the
keyboard. The rule, exactly:

* four or more consonants in a row, or
* the same character three or more times in a row

Anything else is accepted. Comparisons ignore case, and the vowels are
~a e i o u~.

Write ~keysmash(s)~ returning ~True~ when the text breaks either rule.

    keysmash("asdfgh")  ->  True
    keysmash("hello")   ->  False

The rule is the specification, not a judgement: ~rhythms~ trips it, and that is
the correct answer.`,
		timeLimitMs: 5000,
		starters:    py("def keysmash(s):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Four consonants in a row", `keysmash("asdfgh") is True`),
			vis("An ordinary word passes", `keysmash("hello") is False`),
			vis("Three of the same character", `keysmash("aaa") is True`),
			vis("An empty string is fine", `keysmash("") is False`),
			hid("Case does not matter", `keysmash("ASDFGH") is True`),
			hid("Three consonants are not enough", `keysmash("string") is False`),
			hid("A repeated non-letter counts too", `keysmash("a!!!b") is True`),
		},
		solution: "def keysmash(s):\n    low = s.lower()\n    run = 0\n    for ch in low:\n        if ch.isalpha() and ch not in 'aeiou':\n            run += 1\n            if run >= 4:\n                return True\n        else:\n            run = 0\n    rep = 1\n    for a, b in zip(low, low[1:]):\n        rep = rep + 1 if a == b else 1\n        if rep >= 3:\n            return True\n    return False\n",
	},
	{
		slug: "headline-case", title: "Headline case", difficulty: "medium", topic: "Strings",
		statement: `A headline capitalises every word except the small ones — ~a an the of
and in on to for~ — which stay lower case. The first and last word are always
capitalised whatever they are.

Every word is otherwise normalised: first letter up, the rest down.

Write ~headline(s)~.

    headline("the lord of the rings")  ->  "The Lord of the Rings"
    headline("gone WITH the wind")     ->  "Gone With the Wind"`,
		timeLimitMs: 5000,
		starters:    py("def headline(s):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Small words stay small", `headline("the lord of the rings") == "The Lord of the Rings"`),
			vis("Shouting is normalised", `headline("gone WITH the wind") == "Gone With the Wind"`),
			vis("One word is always capitalised", `headline("of") == "Of"`),
			hid("A small word at the end is capitalised", `headline("what are you waiting for") == "What Are You Waiting For"`),
			hid("A small word at the start is capitalised", `headline("a tale of two cities") == "A Tale of Two Cities"`),
		},
		solution: "SMALL = {'a', 'an', 'the', 'of', 'and', 'in', 'on', 'to', 'for'}\n\ndef headline(s):\n    words = s.split()\n    out = []\n    for i, w in enumerate(words):\n        if 0 < i < len(words) - 1 and w.lower() in SMALL:\n            out.append(w.lower())\n        else:\n            out.append(w[0].upper() + w[1:].lower())\n    return ' '.join(out)\n",
	},
}

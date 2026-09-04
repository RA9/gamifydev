package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// stackProblems: last-in-first-out thinking, and the small parsers and state
// machines that fall out of it.
var stackProblems = []seedProblem{
	{
		slug: "undo-history", title: "Undo history", difficulty: "medium", topic: "Stacks",
		statement: `An editor records three kinds of command:

* ~do:X~ — perform action ~X~
* ~undo~ — take back the most recent action
* ~redo~ — put back the most recently undone action

Doing something new throws away everything waiting to be redone: once you have
branched off, the old future is gone. Undo with nothing to undo, and redo with
nothing to redo, both do nothing.

Write ~history(commands)~ returning the list of actions still in effect, oldest
first.

    history(["do:a", "do:b", "undo"])           ->  ["a"]
    history(["do:a", "undo", "do:b", "redo"])   ->  ["b"]`,
		timeLimitMs: 5000,
		starters:    py("def history(commands):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Undo takes back the last action", `history(["do:a", "do:b", "undo"]) == ["a"]`),
			vis("A new action discards the redo stack", `history(["do:a", "undo", "do:b", "redo"]) == ["b"]`),
			vis("Redo puts it back", `history(["do:a", "undo", "redo"]) == ["a"]`),
			vis("Nothing done at all", "history([]) == []"),
			hid("Undo with nothing to undo is harmless", `history(["undo", "do:a"]) == ["a"]`),
			hid("Redo with nothing to redo is harmless", `history(["do:a", "redo"]) == ["a"]`),
			hid("Several undos then several redos", `history(["do:a", "do:b", "undo", "undo", "redo", "redo"]) == ["a", "b"]`),
		},
		solution: "def history(commands):\n    done = []\n    undone = []\n    for c in commands:\n        if c == 'undo':\n            if done:\n                undone.append(done.pop())\n        elif c == 'redo':\n            if undone:\n                done.append(undone.pop())\n        else:\n            done.append(c[3:])\n            undone.clear()\n    return done\n",
	},
	{
		slug: "nesting-depth", title: "Nesting depth", difficulty: "hard", topic: "Stacks",
		statement: `Write ~depth(text)~ returning how deeply the brackets in ~text~ nest.
Brackets come in three kinds — ~()~, ~[]~ and ~{}~ — and each must be closed by
its own kind, in order.

Anything inside double quotes is **text, not structure**: brackets in there are
just characters and are ignored completely. Quotes toggle in and out.

If the brackets do not balance, or a quote is left open at the end, return ~-1~.

    depth("a(b[c]{d})")  ->  2
    depth('("(")')       ->  1

In the second, the inner bracket is quoted, so only the outer pair counts.`,
		timeLimitMs: 5000,
		starters:    py("def depth(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts the deepest nesting", `depth("a(b[c]{d})") == 2`),
			vis("Quoted brackets are just text", `depth('("(")') == 1`),
			vis("No brackets at all", `depth("plain") == 0`),
			vis("An unclosed bracket is invalid", `depth("(") == -1`),
			vis("The wrong closing kind is invalid", `depth("(]") == -1`),
			hid("An unclosed quote is invalid", `depth('"') == -1`),
			hid("A closing bracket with nothing open", `depth(")") == -1`),
			hid("Quoted text can hide an imbalance entirely", `depth('")))"') == 0`),
			hid("Three levels deep", `depth("([{}])") == 3`),
			hid("Six levels deep", `depth("(([[{{}}]]))") == 6`),
		},
		solution: "CLOSERS = {')': '(', ']': '[', '}': '{'}\n\ndef depth(text):\n    stack = []\n    best = 0\n    quoted = False\n    for ch in text:\n        if ch == '\"':\n            quoted = not quoted\n            continue\n        if quoted:\n            continue\n        if ch in '([{':\n            stack.append(ch)\n            if len(stack) > best:\n                best = len(stack)\n        elif ch in CLOSERS:\n            if not stack or stack.pop() != CLOSERS[ch]:\n                return -1\n    if stack or quoted:\n        return -1\n    return best\n",
	},
	{
		slug: "dice-notation", title: "Dice notation", difficulty: "medium", topic: "Parsing",
		statement: `Tabletop games write dice rolls as ~2d6+3~: roll two six-sided dice and add
three. The count may be left off (~d20~ means one), and the modifier may be
missing, positive or negative.

Write ~parse_roll(s)~ returning ~(count, sides, modifier)~.

    parse_roll("2d6+3")  ->  (2, 6, 3)
    parse_roll("d20")    ->  (1, 20, 0)
    parse_roll("3d8-2")  ->  (3, 8, -2)

The input is always well formed, and may have spaces around it.`,
		timeLimitMs: 5000,
		starters:    py("def parse_roll(s):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A full roll", `parse_roll("2d6+3") == (2, 6, 3)`),
			vis("The count defaults to one", `parse_roll("d20") == (1, 20, 0)`),
			vis("A negative modifier", `parse_roll("3d8-2") == (3, 8, -2)`),
			hid("Surrounding spaces are ignored", `parse_roll("  d4  ") == (1, 4, 0)`),
			hid("Multi-digit everything", `parse_roll("10d100+25") == (10, 100, 25)`),
			hid("A modifier of zero is written out", `parse_roll("1d6+0") == (1, 6, 0)`),
		},
		solution: "def parse_roll(s):\n    s = s.strip()\n    mod = 0\n    for i, ch in enumerate(s):\n        if ch in '+-' and i > 0:\n            mod = int(s[i:])\n            s = s[:i]\n            break\n    count, _, sides = s.partition('d')\n    return (int(count) if count else 1, int(sides), mod)\n",
	},
	{
		slug: "bracket-repair", title: "Bracket repair", difficulty: "medium", topic: "Stacks",
		statement: `Write ~repairs(s)~ returning the smallest number of brackets you would have
to **add** to a string of ~(~ and ~)~ to make it balanced. You may add them
anywhere; you may not remove or move anything.

    repairs("(()")   ->  1
    repairs("())(")  ->  2

An empty string needs nothing. You do not need to build the repaired string —
only count.`,
		timeLimitMs: 5000,
		starters:    py("def repairs(s):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("One closing bracket missing", `repairs("(()") == 1`),
			vis("One of each missing", `repairs("())(") == 2`),
			vis("Already balanced", `repairs("()()") == 0`),
			vis("Nothing to repair", `repairs("") == 0`),
			hid("All opening brackets", `repairs("(((") == 3`),
			hid("All closing brackets", `repairs(")))") == 3`),
			hid("Balanced but deeply nested", `repairs("((()))") == 0`),
		},
		solution: "def repairs(s):\n    need_open = 0\n    unmatched = 0\n    for ch in s:\n        if ch == '(':\n            unmatched += 1\n        elif unmatched:\n            unmatched -= 1\n        else:\n            need_open += 1\n    return need_open + unmatched\n",
	},
	{
		slug: "csv-fields", title: "CSV fields", difficulty: "hard", topic: "Parsing",
		statement: `Split one line of CSV into its fields.

* Fields are separated by commas.
* A field may be wrapped in double quotes, and a comma inside quotes is part of
  the field, not a separator.
* Inside a quoted field, two double quotes in a row mean one literal quote.
* The quotes themselves are never part of the value.

Write ~fields(line)~.

    fields('a,"b,c",d')        ->  ["a", "b,c", "d"]
    fields('"he said ""hi"""') ->  ['he said "hi"']

An empty line is one empty field, and ~a,,b~ has an empty field in the middle.`,
		timeLimitMs: 5000,
		starters:    py("def fields(line):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A quoted comma is not a separator", `fields('a,"b,c",d') == ["a", "b,c", "d"]`),
			vis("Doubled quotes are one literal quote", `fields('"he said ""hi"""') == ['he said "hi"']`),
			vis("A plain line", `fields("a,b") == ["a", "b"]`),
			vis("An empty line is one empty field", `fields("") == [""]`),
			hid("An empty field in the middle", `fields("a,,b") == ["a", "", "b"]`),
			hid("A trailing comma leaves an empty field", `fields("a,") == ["a", ""]`),
			hid("An empty quoted field", `fields('a,"",b') == ["a", "", "b"]`),
			hid("Quotes disappear from the value", `fields('"plain"') == ["plain"]`),
		},
		solution: "def fields(line):\n    out = []\n    cur = []\n    quoted = False\n    i = 0\n    while i < len(line):\n        ch = line[i]\n        if quoted:\n            if ch == '\"':\n                if i + 1 < len(line) and line[i + 1] == '\"':\n                    cur.append('\"')\n                    i += 2\n                    continue\n                quoted = False\n            else:\n                cur.append(ch)\n        elif ch == '\"':\n            quoted = True\n        elif ch == ',':\n            out.append(''.join(cur))\n            cur = []\n        else:\n            cur.append(ch)\n        i += 1\n    out.append(''.join(cur))\n    return out\n",
	},
	{
		slug: "indent-parents", title: "Indent parents", difficulty: "hard", topic: "Stacks",
		statement: `An outline is written with two spaces of indent per level. Write
~parents(lines)~ returning, for each line, the index of the line it hangs from —
the nearest line above it at one less level. A line at level 0 hangs from
nothing: report ~-1~.

    parents(["a", "  b", "  c", "    d", "e"])  ->  [-1, 0, 0, 2, -1]

Levels never skip on the way in: a line is at most one level deeper than the one
above it. They can skip on the way out — ~d~ is followed by ~e~ two levels
shallower — and that is the case worth thinking about.`,
		timeLimitMs: 5000,
		starters:    py("def parents(lines):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Finds each line's parent", `parents(["a", "  b", "  c", "    d", "e"]) == [-1, 0, 0, 2, -1]`),
			vis("A flat list has no parents", `parents(["a", "b"]) == [-1, -1]`),
			vis("Nothing to parse", "parents([]) == []"),
			vis("A single chain", `parents(["a", "  b", "    c"]) == [-1, 0, 1]`),
			hid("Dropping several levels at once", `parents(["a", "  b", "    c", "  d"]) == [-1, 0, 1, 0]`),
			hid("Two independent trees", `parents(["a", "  b", "c", "  d"]) == [-1, 0, -1, 2]`),
		},
		solution: "def parents(lines):\n    stack = []\n    out = []\n    for i, line in enumerate(lines):\n        level = (len(line) - len(line.lstrip(' '))) // 2\n        while stack and stack[-1][0] >= level:\n            stack.pop()\n        out.append(stack[-1][1] if stack else -1)\n        stack.append((level, i))\n    return out\n",
	},
	{
		slug: "tag-nesting", title: "Tag nesting", difficulty: "medium", topic: "Stacks",
		statement: `A markup document is given as a list of tag names in the order they appear:
~"a"~ opens a tag, ~"/a"~ closes it, and a name ending in ~/~ (like ~"br/"~) is
self-closing and opens nothing.

Write ~well_formed(tags)~ returning whether every tag is closed, in the right
order, by a closer of the same name.

    well_formed(["a", "b", "/b", "/a"])  ->  True
    well_formed(["a", "b", "/a", "/b"])  ->  False`,
		timeLimitMs: 5000,
		starters:    py("def well_formed(tags):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Properly nested", `well_formed(["a", "b", "/b", "/a"]) is True`),
			vis("Crossed tags are not nesting", `well_formed(["a", "b", "/a", "/b"]) is False`),
			vis("An empty document is well formed", "well_formed([]) is True"),
			vis("A self-closing tag opens nothing", `well_formed(["br/"]) is True`),
			hid("An unclosed tag", `well_formed(["a"]) is False`),
			hid("A closer with nothing open", `well_formed(["/a"]) is False`),
			hid("The wrong name closing", `well_formed(["a", "/b"]) is False`),
			hid("Self-closing tags among real ones", `well_formed(["a", "br/", "/a"]) is True`),
		},
		solution: "def well_formed(tags):\n    stack = []\n    for t in tags:\n        if t.endswith('/'):\n            continue\n        if t.startswith('/'):\n            if not stack or stack.pop() != t[1:]:\n                return False\n        else:\n            stack.append(t)\n    return not stack\n",
	},
	{
		slug: "stack-machine", title: "Stack machine", difficulty: "medium", topic: "Stacks",
		statement: `A tiny machine runs four instructions against a stack:

* ~push n~ — put the number ~n~ on top
* ~add~ — take the top two, put back their sum
* ~dup~ — copy the top
* ~drop~ — throw away the top

Write ~run(program)~ returning the stack when the program finishes, bottom
first. If an instruction needs more than is on the stack, the machine faults:
return ~None~ immediately.

    run(["push 2", "push 3", "add"])  ->  [5]
    run(["add"])                      ->  None`,
		timeLimitMs: 5000,
		starters:    py("def run(program):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Adds two numbers", `run(["push 2", "push 3", "add"]) == [5]`),
			vis("Adding on an empty stack faults", `run(["add"]) is None`),
			vis("An empty program leaves an empty stack", "run([]) == []"),
			vis("Duplicating the top", `run(["push 1", "dup"]) == [1, 1]`),
			hid("Dropping from an empty stack faults", `run(["drop"]) is None`),
			hid("Adding with only one value faults", `run(["push 1", "add"]) is None`),
			hid("Negative numbers push fine", `run(["push -5", "push 5", "add"]) == [0]`),
			hid("A fault stops the program there", `run(["push 1", "drop", "drop", "push 9"]) is None`),
		},
		solution: "def run(program):\n    stack = []\n    for ins in program:\n        if ins.startswith('push '):\n            stack.append(int(ins[5:]))\n        elif ins == 'add':\n            if len(stack) < 2:\n                return None\n            b = stack.pop()\n            a = stack.pop()\n            stack.append(a + b)\n        elif ins == 'dup':\n            if not stack:\n                return None\n            stack.append(stack[-1])\n        elif ins == 'drop':\n            if not stack:\n                return None\n            stack.pop()\n    return stack\n",
	},
	{
		slug: "days-until-taller", title: "Days until taller", difficulty: "hard", topic: "Stacks",
		statement: `Standing at each building along a street, how far along is the first
building **at least as tall** as this one? Write ~wait(heights)~ returning that
distance for each building, or 0 when nothing ahead is tall enough.

    wait([3, 1, 4])  ->  [2, 1, 0]

Streets run to two hundred thousand buildings, so looking ahead from every
building is too slow. The way through: walk the street once, keeping the
buildings still waiting for an answer. Each new building answers all the
shorter ones behind it at once — and each building is answered only once, so the
work is bounded even though the inner loop looks like it isn't.`,
		timeLimitMs: 4000,
		starters:    py("def wait(heights):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Distances to the next tall enough building", "wait([3, 1, 4]) == [2, 1, 0]"),
			vis("Nothing ahead is tall enough", "wait([9, 1, 1]) == [0, 1, 0]"),
			vis("An empty street", "wait([]) == []"),
			vis("Equal heights count as tall enough", "wait([2, 2]) == [1, 0]"),
			vis("Fast enough for a long street", "wait(list(range(200000, 0, -1))) == [0] * 200000"),
			hid("An ascending street answers every building next door", "wait(list(range(1000))) == [1] * 999 + [0]"),
			hid("A descending street answers nobody", "wait([5, 4, 3]) == [0, 0, 0]"),
			hid("One tall building answers everything behind it", "wait([1, 1, 1, 9]) == [1, 1, 1, 0]"),
		},
		solution: "def wait(heights):\n    out = [0] * len(heights)\n    stack = []\n    for i, h in enumerate(heights):\n        while stack and heights[stack[-1]] <= h:\n            j = stack.pop()\n            out[j] = i - j\n        stack.append(i)\n    return out\n",
		tooSlow:  "def wait(heights):\n    out = []\n    for i, h in enumerate(heights):\n        d = 0\n        for j in range(i + 1, len(heights)):\n            if heights[j] >= h:\n                d = j - i\n                break\n        out.append(d)\n    return out\n",
	},
	{
		slug: "macro-expand", title: "Macro expansion", difficulty: "hard", topic: "Recursion",
		statement: `A build file defines macros: ~macros~ maps a name to a body, and a body is
words separated by single spaces. A word that is itself a macro name expands in
turn; a word that is not is left as it is.

Write ~expand(macros, name)~ returning the fully expanded text.

    expand({"a": "b c", "b": "hello"}, "a")  ->  "hello c"

A macro that reaches itself, however indirectly, can never finish expanding:
return ~None~ rather than looping forever.

    expand({"a": "b", "b": "a"}, "a")  ->  None`,
		timeLimitMs: 5000,
		starters:    py("def expand(macros, name):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Expands through another macro", `expand({"a": "b c", "b": "hello"}, "a") == "hello c"`),
			vis("A cycle never finishes", `expand({"a": "b", "b": "a"}, "a") is None`),
			vis("An undefined name is its own text", `expand({}, "x") == "x"`),
			vis("A macro that expands to nothing else", `expand({"a": "plain text"}, "a") == "plain text"`),
			hid("A macro that names itself", `expand({"a": "a"}, "a") is None`),
			hid("The same macro used twice is not a cycle", `expand({"a": "b b", "b": "x"}, "a") == "x x"`),
			hid("A cycle further down still counts", `expand({"a": "b", "b": "c", "c": "b"}, "a") is None`),
		},
		solution: "def expand(macros, name):\n    active = set()\n\n    def go(n):\n        if n in active:\n            return None\n        if n not in macros:\n            return n\n        active.add(n)\n        parts = []\n        for w in macros[n].split(' '):\n            r = go(w)\n            if r is None:\n                return None\n            parts.append(r)\n        active.discard(n)\n        return ' '.join(parts)\n\n    return go(name)\n",
	},
	{
		slug: "text-edits", title: "Text edits", difficulty: "medium", topic: "Stacks",
		statement: `A terminal treats two characters as commands rather than text:

* ~#~ deletes the character before it
* ~@~ deletes everything back to the nearest space, leaving the space itself

Either command on an empty line does nothing. Write ~apply_edits(text)~
returning what is actually on screen.

    apply_edits("ab#c")     ->  "ac"
    apply_edits("hi the@")  ->  "hi "`,
		timeLimitMs: 5000,
		starters:    py("def apply_edits(text):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Backspace removes one character", `apply_edits("ab#c") == "ac"`),
			vis("The word command removes back to the space", `apply_edits("hi the@") == "hi "`),
			vis("Nothing typed", `apply_edits("") == ""`),
			vis("Deleting from an empty line does nothing", `apply_edits("###a") == "a"`),
			hid("The word command on the first word clears the line", `apply_edits("hello@") == ""`),
			hid("The word command on an empty line", `apply_edits("@") == ""`),
			hid("Backspacing over a space", `apply_edits("a b#") == "a "`),
			hid("Several commands in a row", `apply_edits("abc##d") == "ad"`),
		},
		solution: "def apply_edits(text):\n    out = []\n    for ch in text:\n        if ch == '#':\n            if out:\n                out.pop()\n        elif ch == '@':\n            while out and out[-1] != ' ':\n                out.pop()\n        else:\n            out.append(ch)\n    return ''.join(out)\n",
	},
	{
		slug: "expand-runs", title: "Expand runs", difficulty: "medium", topic: "Parsing",
		statement: `A compressed string writes a character followed by how many times it
repeats. A character with no number after it appears once, and the count can be
more than one digit.

Write ~decode(s)~ returning the original text.

    decode("a3b2")  ->  "aaabb"
    decode("abc")   ->  "abc"
    decode("a12")   ->  "aaaaaaaaaaaa"

Digits only ever follow a character, never start the string.`,
		timeLimitMs: 5000,
		starters:    py("def decode(s):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Expands the runs", `decode("a3b2") == "aaabb"`),
			vis("No numbers means no repeats", `decode("abc") == "abc"`),
			vis("A multi-digit count", `decode("a12") == "a" * 12`),
			vis("Nothing to decode", `decode("") == ""`),
			hid("A count of one is allowed to be written", `decode("a1b1") == "ab"`),
			hid("Mixed counted and uncounted", `decode("ab3c") == "abbbc"`),
			hid("Spaces are characters too", `decode(" 3") == "   "`),
		},
		solution: "def decode(s):\n    out = []\n    i = 0\n    while i < len(s):\n        ch = s[i]\n        i += 1\n        j = i\n        while j < len(s) and s[j].isdigit():\n            j += 1\n        n = int(s[i:j]) if j > i else 1\n        out.append(ch * n)\n        i = j\n    return ''.join(out)\n",
	},
}

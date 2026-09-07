package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// stackProblems: last-in-first-out thinking, and the small parsers and state
// machines that fall out of it.
var stackProblems = []seedProblem{
	{
		slug: "undo-history", title: "Undo history", difficulty: "medium", topic: "Stacks",
		statement: `An editor records three kinds of command: ~do:X~ performs action ~X~,
~undo~ takes back the most recent action, ~redo~ puts back the most recently
undone one.

Doing something new throws away everything waiting to be redone. Undo with
nothing to undo, and redo with nothing to redo, both do nothing.

**Input.** A count ~n~, then ~n~ commands, one per line.

**Output.** The actions still in effect, oldest first, space-separated.

    Input        Output
    3            a
    do:a
    do:b
    undo`,
		timeLimitMs: 4000,
		starters:    triWords("keep a done stack and an undone stack"),
		tests: []store.Check{
			tokens("Undo takes back the last action", "3\ndo:a do:b undo\n", "a"),
			tokens("A new action discards the redo stack", "4\ndo:a undo do:b redo\n", "b"),
			tokens("Redo puts it back", "3\ndo:a undo redo\n", "a"),
			tokens("Nothing done at all", "0\n", ""),
			hid2("Undo with nothing to undo is harmless", "2\nundo do:a\n", "a"),
			hid2("Redo with nothing to redo is harmless", "2\ndo:a redo\n", "a"),
			hid2("Several undos then several redos", "6\ndo:a do:b undo undo redo redo\n", "a b"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char done[100000][32], undone[100000][32];
    int d = 0, u = 0;
    for (int i = 0; i < n; i++) {
        char cmd[40];
        if (scanf("%39s", cmd) != 1) return 1;
        if (strcmp(cmd, "undo") == 0) { if (d > 0) strcpy(undone[u++], done[--d]); }
        else if (strcmp(cmd, "redo") == 0) { if (u > 0) strcpy(done[d++], undone[--u]); }
        else { strcpy(done[d++], cmd + 3); u = 0; }
    }
    for (int i = 0; i < d; i++) printf(i ? " %s" : "%s", done[i]);
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var done, undone []string
	for i := 0; i < n; i++ {
		var cmd string
		fmt.Fscan(in, &cmd)
		switch {
		case cmd == "undo":
			if len(done) > 0 {
				undone = append(undone, done[len(done)-1])
				done = done[:len(done)-1]
			}
		case cmd == "redo":
			if len(undone) > 0 {
				done = append(done, undone[len(undone)-1])
				undone = undone[:len(undone)-1]
			}
		default:
			done = append(done, cmd[3:])
			undone = undone[:0]
		}
	}
	fmt.Println(strings.Join(done, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    done, undone = [], []
    for cmd in data[1:1 + n]:
        if cmd == "undo":
            if done:
                undone.append(done.pop())
        elif cmd == "redo":
            if undone:
                done.append(undone.pop())
        else:
            done.append(cmd[3:])
            undone.clear()
    print(" ".join(done))

main()
`},
	},
	{
		slug: "nesting-depth", title: "Nesting depth", difficulty: "hard", topic: "Stacks",
		statement: `Brackets come in three kinds — ~()~, ~[]~ and ~{}~ — and each must be
closed by its own kind, in order.

Anything inside double quotes is **text, not structure**: brackets in there are
just characters and are ignored completely. Quotes toggle in and out.

**Input.** One line.

**Output.** How deeply the brackets nest. If they do not balance, or a quote is
left open at the end, print ~-1~.

    Input          Output
    a(b[c]{d})     2

    Input          Output
    ("(")          1

In the second, the inner bracket is quoted, so only the outer pair counts.`,
		timeLimitMs: 4000,
		starters:    triLine("track the open brackets, ignoring anything inside quotes"),
		tests: []store.Check{
			tokens("Counts the deepest nesting", "a(b[c]{d})\n", "2"),
			tokens("Quoted brackets are just text", "(\"(\")\n", "1"),
			tokens("No brackets at all", "plain\n", "0"),
			tokens("An unclosed bracket is invalid", "(\n", "-1"),
			tokens("The wrong closing kind is invalid", "(]\n", "-1"),
			hid2("An unclosed quote is invalid", "\"\n", "-1"),
			hid2("A closing bracket with nothing open", ")\n", "-1"),
			hid2("Quoted text can hide an imbalance entirely", "\")))\"\n", "0"),
			hid2("Six levels deep", "(([[{{}}]]))\n", "6"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005], stack[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int top = 0, best = 0, quoted = 0;
    for (int i = 0; line[i]; i++) {
        char c = line[i];
        if (c == '"') { quoted = !quoted; continue; }
        if (quoted) continue;
        if (c == '(' || c == '[' || c == '{') {
            stack[top++] = c;
            if (top > best) best = top;
        } else if (c == ')' || c == ']' || c == '}') {
            char want = c == ')' ? '(' : (c == ']' ? '[' : '{');
            if (top == 0 || stack[--top] != want) { printf("-1\n"); return 0; }
        }
    }
    printf("%d\n", (top || quoted) ? -1 : best);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	closers := map[byte]byte{')': '(', ']': '[', '}': '{'}
	var stack []byte
	best, quoted := 0, false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if c == '"' {
			quoted = !quoted
			continue
		}
		if quoted {
			continue
		}
		if c == '(' || c == '[' || c == '{' {
			stack = append(stack, c)
			if len(stack) > best {
				best = len(stack)
			}
		} else if want, ok := closers[c]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != want {
				fmt.Println(-1)
				return
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 || quoted {
		fmt.Println(-1)
		return
	}
	fmt.Println(best)
}
`,
			"python": `import sys

CLOSERS = {")": "(", "]": "[", "}": "{"}

def main():
    line = sys.stdin.readline().rstrip("\n")
    stack, best, quoted = [], 0, False
    for ch in line:
        if ch == '"':
            quoted = not quoted
            continue
        if quoted:
            continue
        if ch in "([{":
            stack.append(ch)
            best = max(best, len(stack))
        elif ch in CLOSERS:
            if not stack or stack.pop() != CLOSERS[ch]:
                print(-1)
                return
    print(-1 if stack or quoted else best)

main()
`},
	},
	{
		slug: "dice-notation", title: "Dice notation", difficulty: "medium", topic: "Parsing",
		statement: `Tabletop games write dice rolls as ~2d6+3~: roll two six-sided dice and add
three. The count may be left off (~d20~ means one), and the modifier may be
missing, positive or negative.

**Input.** One roll, as a single word.

**Output.** ~count sides modifier~, space-separated.

    Input     Output
    2d6+3     2 6 3

    Input     Output
    d20       1 20 0`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    char s[64];
    if (scanf("%63s", s) != 1) return 1;
    /* split on d, then on any + or - */
    printf("%d %d %d\n", 0, 0, 0);
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	// split on d, then on any + or -
	fmt.Println(0, 0, 0)
}
`, `import sys

def main():
    s = sys.stdin.read().split()[0]
    # split on d, then on any + or -
    print(0, 0, 0)

main()
`),
		tests: []store.Check{
			tokens("A full roll", "2d6+3\n", "2 6 3"),
			tokens("The count defaults to one", "d20\n", "1 20 0"),
			tokens("A negative modifier", "3d8-2\n", "3 8 -2"),
			hid2("Multi-digit everything", "10d100+25\n", "10 100 25"),
			hid2("A modifier of zero is written out", "1d6+0\n", "1 6 0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    char s[64];
    if (scanf("%63s", s) != 1) return 1;
    long mod = 0;
    for (int i = 1; s[i]; i++) {
        if (s[i] == '+' || s[i] == '-') { mod = strtol(s + i, NULL, 10); s[i] = '\0'; break; }
    }
    char *d = strchr(s, 'd');
    long count = (d == s) ? 1 : strtol(s, NULL, 10);
    long sides = strtol(d + 1, NULL, 10);
    printf("%ld %ld %ld\n", count, sides, mod);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	mod := 0
	if i := strings.IndexAny(s[1:], "+-"); i >= 0 {
		mod, _ = strconv.Atoi(s[1+i:])
		s = s[:1+i]
	}
	parts := strings.SplitN(s, "d", 2)
	count := 1
	if parts[0] != "" {
		count, _ = strconv.Atoi(parts[0])
	}
	sides, _ := strconv.Atoi(parts[1])
	fmt.Println(count, sides, mod)
}
`,
			"python": `import sys

def main():
    s = sys.stdin.read().split()[0]
    mod = 0
    for i, ch in enumerate(s):
        if ch in "+-" and i > 0:
            mod = int(s[i:])
            s = s[:i]
            break
    count, _, sides = s.partition("d")
    print(int(count) if count else 1, int(sides), mod)

main()
`},
	},
	{
		slug: "bracket-repair", title: "Bracket repair", difficulty: "medium", topic: "Stacks",
		statement: `**Input.** One line of ~(~ and ~)~. An empty line is written as ~-~.

**Output.** The smallest number of brackets you would have to **add** to make it
balanced. You may add them anywhere; you may not remove or move anything.

    Input     Output
    (()       1

    Input     Output
    ())(      2`,
		timeLimitMs: 4000,
		starters:    triLine("count the closers with nothing open, and the openers left over"),
		tests: []store.Check{
			tokens("One closing bracket missing", "(()\n", "1"),
			tokens("One of each missing", "())(\n", "2"),
			tokens("Already balanced", "()()\n", "0"),
			tokens("Nothing to repair", "-\n", "0"),
			hid2("All opening brackets", "(((\n", "3"),
			hid2("All closing brackets", ")))\n", "3"),
			hid2("Balanced but deeply nested", "((()))\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    int open = 0, need = 0;
    for (int i = 0; line[i]; i++) {
        if (line[i] == '(') open++;
        else if (line[i] == ')') { if (open) open--; else need++; }
    }
    printf("%d\n", need + open);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	line, _ := in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	open, need := 0, 0
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '(':
			open++
		case ')':
			if open > 0 {
				open--
			} else {
				need++
			}
		}
	}
	fmt.Println(need + open)
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    open_count = need = 0
    for ch in line:
        if ch == "(":
            open_count += 1
        elif ch == ")":
            if open_count:
                open_count -= 1
            else:
                need += 1
    print(need + open_count)

main()
`},
	},
	{
		slug: "csv-fields", title: "CSV fields", difficulty: "hard", topic: "Parsing",
		statement: `Split one line of CSV into its fields.

* Fields are separated by commas.
* A field may be wrapped in double quotes, and a comma inside quotes is part of
  the field, not a separator.
* Inside a quoted field, two double quotes in a row mean one literal quote.
* The quotes themselves are never part of the value.

**Input.** One line. An empty line is written as ~-~ and has one empty field.

**Output.** One field per line. An empty field prints an empty line.

    Input           Output
    a,"b,c",d       a
                    b,c
                    d`,
		timeLimitMs: 4000,
		starters:    triLine("walk the line, tracking whether you are inside quotes"),
		tests: []store.Check{
			exact("A quoted comma is not a separator", "a,\"b,c\",d\n", "a\nb,c\nd"),
			exact("Doubled quotes are one literal quote", "\"he said \"\"hi\"\"\"\n", "he said \"hi\""),
			exact("A plain line", "a,b\n", "a\nb"),
			exact("An empty line is one empty field", "-\n", ""),
			exactHid("An empty field in the middle", "a,,b\n", "a\n\nb"),
			exactHid("A trailing comma leaves an empty field", "a,\n", "a"),
			exactHid("An empty quoted field", "a,\"\",b\n", "a\n\nb"),
			exactHid("Quotes disappear from the value", "\"plain\"\n", "plain"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    if (strcmp(line, "-") == 0) { printf("\n"); return 0; }
    int quoted = 0, i = 0;
    while (line[i]) {
        char c = line[i];
        if (quoted) {
            if (c == '"') {
                if (line[i + 1] == '"') { putchar('"'); i += 2; continue; }
                quoted = 0;
            } else putchar(c);
        } else if (c == '"') quoted = 1;
        else if (c == ',') putchar('\n');
        else putchar(c);
        i++;
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	raw, _ := in.ReadString('\n')
	line := strings.TrimRight(raw, "\r\n")
	if line == "-" {
		fmt.Println("")
		return
	}
	var b strings.Builder
	quoted := false
	for i := 0; i < len(line); {
		c := line[i]
		switch {
		case quoted && c == '"':
			if i+1 < len(line) && line[i+1] == '"' {
				b.WriteByte('"')
				i += 2
				continue
			}
			quoted = false
		case quoted:
			b.WriteByte(c)
		case c == '"':
			quoted = true
		case c == ',':
			b.WriteByte('\n')
		default:
			b.WriteByte(c)
		}
		i++
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    if line == "-":
        print("")
        return
    out, quoted, i = [], False, 0
    while i < len(line):
        c = line[i]
        if quoted:
            if c == '"':
                if i + 1 < len(line) and line[i + 1] == '"':
                    out.append('"')
                    i += 2
                    continue
                quoted = False
            else:
                out.append(c)
        elif c == '"':
            quoted = True
        elif c == ",":
            out.append("\n")
        else:
            out.append(c)
        i += 1
    print("".join(out))

main()
`},
	},
	{
		slug: "indent-parents", title: "Indent parents", difficulty: "hard", topic: "Stacks",
		statement: `An outline is written with two spaces of indent per level.

**Input.** A count ~n~, then ~n~ lines. Each is a depth ~d~ followed by the
line's text as one word.

**Output.** For each line, the index of the line it hangs from — the nearest line
above it at one less level — space-separated. A line at level 0 hangs from
nothing: report ~-1~.

    Input      Output
    5          -1 0 0 2 -1
    0 a
    1 b
    1 c
    2 d
    0 e

Levels never skip on the way in. They can skip on the way out — ~d~ is followed
by ~e~ two levels shallower — and that is the case worth thinking about.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static int level[100000], stackLevel[100000], stackIdx[100000];
    int top = 0;
    /* pop back to the parent's level, then report and push */
    printf("\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	// pop back to the parent's level, then report and push
	_ = n
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    levels = [int(data[1 + 2 * i]) for i in range(n)]
    # pop back to the parent's level, then report and push
    print("")

main()
`),
		tests: []store.Check{
			tokens("Finds each line's parent", "5\n0 a\n1 b\n1 c\n2 d\n0 e\n", "-1 0 0 2 -1"),
			tokens("A flat list has no parents", "2\n0 a\n0 b\n", "-1 -1"),
			tokens("Nothing to parse", "0\n", ""),
			tokens("A single chain", "3\n0 a\n1 b\n2 c\n", "-1 0 1"),
			hid2("Dropping several levels at once", "4\n0 a\n1 b\n2 c\n1 d\n", "-1 0 1 0"),
			hid2("Two independent trees", "4\n0 a\n1 b\n0 c\n1 d\n", "-1 0 -1 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static int stackLevel[100000], stackIdx[100000];
    int top = 0, first = 1;
    for (int i = 0; i < n; i++) {
        int level;
        char text[64];
        if (scanf("%d %63s", &level, text) != 2) return 1;
        while (top > 0 && stackLevel[top - 1] >= level) top--;
        int parent = top > 0 ? stackIdx[top - 1] : -1;
        printf(first ? "%d" : " %d", parent);
        first = 0;
        stackLevel[top] = level;
        stackIdx[top] = i;
        top++;
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type frame struct{ level, idx int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var stack []frame
	out := make([]string, n)
	for i := 0; i < n; i++ {
		var level int
		var text string
		fmt.Fscan(in, &level, &text)
		for len(stack) > 0 && stack[len(stack)-1].level >= level {
			stack = stack[:len(stack)-1]
		}
		parent := -1
		if len(stack) > 0 {
			parent = stack[len(stack)-1].idx
		}
		out[i] = strconv.Itoa(parent)
		stack = append(stack, frame{level, i})
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    stack, out = [], []
    for i in range(n):
        level = int(data[1 + 2 * i])
        while stack and stack[-1][0] >= level:
            stack.pop()
        out.append(str(stack[-1][1] if stack else -1))
        stack.append((level, i))
    print(" ".join(out))

main()
`},
	},
	{
		slug: "tag-nesting", title: "Tag nesting", difficulty: "medium", topic: "Stacks",
		statement: `A markup document is given as tag names in the order they appear: ~a~ opens
a tag, ~/a~ closes it, and a name ending in ~/~ (like ~br/~) is self-closing and
opens nothing.

**Input.** A count ~n~, then ~n~ tag names.

**Output.** ~yes~ if every tag is closed, in the right order, by a closer of the
same name; ~no~ otherwise.

    Input        Output
    4            yes
    a b /b /a`,
		timeLimitMs: 4000,
		starters:    triWords("push openers, match closers, and check nothing is left"),
		tests: []store.Check{
			tokens("Properly nested", "4\na b /b /a\n", "yes"),
			tokens("Crossed tags are not nesting", "4\na b /a /b\n", "no"),
			tokens("An empty document is well formed", "0\n", "yes"),
			tokens("A self-closing tag opens nothing", "1\nbr/\n", "yes"),
			hid2("An unclosed tag", "1\na\n", "no"),
			hid2("A closer with nothing open", "1\n/a\n", "no"),
			hid2("The wrong name closing", "2\na /b\n", "no"),
			hid2("Self-closing tags among real ones", "3\na br/ /a\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char stack[100000][32];
    int top = 0, ok = 1;
    for (int i = 0; i < n; i++) {
        char t[32];
        if (scanf("%31s", t) != 1) return 1;
        size_t len = strlen(t);
        if (len && t[len - 1] == '/') continue;
        if (t[0] == '/') {
            if (top == 0 || strcmp(stack[top - 1], t + 1) != 0) { ok = 0; break; }
            top--;
        } else {
            strcpy(stack[top++], t);
        }
    }
    printf("%s\n", (ok && top == 0) ? "yes" : "no");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var stack []string
	ok := true
	for i := 0; i < n; i++ {
		var t string
		fmt.Fscan(in, &t)
		switch {
		case strings.HasSuffix(t, "/"):
		case strings.HasPrefix(t, "/"):
			if len(stack) == 0 || stack[len(stack)-1] != t[1:] {
				ok = false
			} else {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, t)
		}
		if !ok {
			break
		}
	}
	if ok && len(stack) == 0 {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    stack, ok = [], True
    for t in data[1:1 + n]:
        if t.endswith("/"):
            continue
        if t.startswith("/"):
            if not stack or stack.pop() != t[1:]:
                ok = False
                break
        else:
            stack.append(t)
    print("yes" if ok and not stack else "no")

main()
`},
	},
	{
		slug: "stack-machine", title: "Stack machine", difficulty: "medium", topic: "Stacks",
		statement: `A tiny machine runs four instructions against a stack: ~push n~ puts ~n~ on
top, ~add~ replaces the top two with their sum, ~dup~ copies the top, ~drop~
throws the top away.

**Input.** A count ~n~, then ~n~ instructions, one per line.

**Output.** The stack when the program finishes, bottom first, space-separated.
If an instruction needs more than is on the stack, the machine faults: print
~fault~ and stop.

    Input        Output
    3            5
    push 2
    push 3
    add`,
		timeLimitMs: 4000,
		starters:    triLines("run each instruction, faulting on an empty stack"),
		tests: []store.Check{
			tokens("Adds two numbers", "3\npush 2\npush 3\nadd\n", "5"),
			tokens("Adding on an empty stack faults", "1\nadd\n", "fault"),
			tokens("An empty program leaves an empty stack", "0\n", ""),
			tokens("Duplicating the top", "2\npush 1\ndup\n", "1 1"),
			hid2("Dropping from an empty stack faults", "1\ndrop\n", "fault"),
			hid2("Adding with only one value faults", "2\npush 1\nadd\n", "fault"),
			hid2("Negative numbers push fine", "3\npush -5\npush 5\nadd\n", "0"),
			hid2("A fault stops the program there", "4\npush 1\ndrop\ndrop\npush 9\n", "fault"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    getchar();
    static long long stack[100000];
    int top = 0;
    for (int i = 0; i < n; i++) {
        char line[64];
        if (!fgets(line, sizeof line, stdin)) break;
        line[strcspn(line, "\n")] = '\0';
        if (strncmp(line, "push ", 5) == 0) stack[top++] = strtoll(line + 5, NULL, 10);
        else if (strcmp(line, "add") == 0) {
            if (top < 2) { printf("fault\n"); return 0; }
            long long b = stack[--top], a = stack[--top];
            stack[top++] = a + b;
        } else if (strcmp(line, "dup") == 0) {
            if (top < 1) { printf("fault\n"); return 0; }
            stack[top] = stack[top - 1]; top++;
        } else if (strcmp(line, "drop") == 0) {
            if (top < 1) { printf("fault\n"); return 0; }
            top--;
        }
    }
    for (int i = 0; i < top; i++) printf(i ? " %lld" : "%lld", stack[i]);
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	var stack []int
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		line := strings.TrimRight(raw, "\r\n")
		switch {
		case strings.HasPrefix(line, "push "):
			v, _ := strconv.Atoi(line[5:])
			stack = append(stack, v)
		case line == "add":
			if len(stack) < 2 {
				fmt.Println("fault")
				return
			}
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			stack = append(stack[:len(stack)-2], a+b)
		case line == "dup":
			if len(stack) < 1 {
				fmt.Println("fault")
				return
			}
			stack = append(stack, stack[len(stack)-1])
		case line == "drop":
			if len(stack) < 1 {
				fmt.Println("fault")
				return
			}
			stack = stack[:len(stack)-1]
		}
	}
	out := make([]string, len(stack))
	for i, v := range stack {
		out[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    stack = []
    for line in raw[1:1 + n]:
        line = line.strip()
        if line.startswith("push "):
            stack.append(int(line[5:]))
        elif line == "add":
            if len(stack) < 2:
                print("fault"); return
            b, a = stack.pop(), stack.pop()
            stack.append(a + b)
        elif line == "dup":
            if not stack:
                print("fault"); return
            stack.append(stack[-1])
        elif line == "drop":
            if not stack:
                print("fault"); return
            stack.pop()
    print(" ".join(str(v) for v in stack))

main()
`},
	},
	{
		slug: "days-until-taller", title: "Days until taller", difficulty: "hard", topic: "Stacks",
		statement: `Standing at each building along a street, how far along is the first
building **at least as tall** as this one?

**Input.** A count ~n~, then ~n~ heights.

**Output.** That distance for each building, space-separated, or 0 when nothing
ahead is tall enough.

    Input      Output
    3          2 1 0
    3 1 4

Looking ahead from every building is quadratic. The way through: walk the street
once, keeping the buildings still waiting for an answer. Each new building
answers all the shorter ones behind it at once — and each is answered only once,
so the work is bounded even though the inner loop looks like it isn't.`,
		timeLimitMs: 4000,
		starters:    triCount("keep the unanswered buildings on a stack"),
		tests: []store.Check{
			tokens("Distances to the next tall enough building", "3\n3 1 4\n", "2 1 0"),
			tokens("Nothing ahead is tall enough", "3\n9 1 1\n", "0 1 0"),
			tokens("An empty street", "0\n", ""),
			tokens("Equal heights count as tall enough", "2\n2 2\n", "1 0"),
			hid2("A descending street answers nobody", "3\n5 4 3\n", "0 0 0"),
			hid2("One tall building answers everything behind it", "4\n1 1 1 9\n", "1 1 1 0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long *h = malloc((size_t)(n ? n : 1) * sizeof *h);
    int *out = calloc((size_t)(n ? n : 1), sizeof *out);
    int *stack = malloc((size_t)(n ? n : 1) * sizeof *stack);
    int top = 0;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld", &h[i]) != 1) return 1;
        while (top > 0 && h[stack[top - 1]] <= h[i]) { int j = stack[--top]; out[j] = i - j; }
        stack[top++] = i;
    }
    for (int i = 0; i < n; i++) printf(i ? " %d" : "%d", out[i]);
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	h := make([]int, n)
	out := make([]int, n)
	var stack []int
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
		for len(stack) > 0 && h[stack[len(stack)-1]] <= h[i] {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			out[j] = i - j
		}
		stack = append(stack, i)
	}
	parts := make([]string, n)
	for i, v := range out {
		parts[i] = strconv.Itoa(v)
	}
	fmt.Println(strings.Join(parts, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    h = [int(x) for x in data[1:1 + n]]
    out, stack = [0] * n, []
    for i, v in enumerate(h):
        while stack and h[stack[-1]] <= v:
            j = stack.pop()
            out[j] = i - j
        stack.append(i)
    print(" ".join(str(x) for x in out))

main()
`},
	},
	{
		slug: "macro-expand", title: "Macro expansion", difficulty: "hard", topic: "Recursion",
		statement: `A build file defines macros. A body is words separated by single spaces; a
word that is itself a macro name expands in turn, and a word that is not is left
as it is.

**Input.** A count ~n~, then ~n~ lines of ~name body...~ (the rest of the line is
the body). The last line names the macro to expand.

**Output.** The fully expanded text. A macro that reaches itself, however
indirectly, can never finish expanding: print ~cycle~.

    Input        Output
    2            hello c
    a b c
    b hello
    a`,
		timeLimitMs: 4000,
		starters:    triLines("expand each word, watching for a macro that reaches itself"),
		tests: []store.Check{
			exact("Expands through another macro", "2\na b c\nb hello\na\n", "hello c"),
			exact("A cycle never finishes", "2\na b\nb a\na\n", "cycle"),
			exact("An undefined name is its own text", "0\nx\n", "x"),
			exact("A macro that expands to nothing else", "1\na plain text\na\n", "plain text"),
			exactHid("A macro that names itself", "1\na a\na\n", "cycle"),
			exactHid("The same macro used twice is not a cycle", "2\na b b\nb x\na\n", "x x"),
			exactHid("A cycle further down still counts", "3\na b\nb c\nc b\na\n", "cycle"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

static char name[1000][32], body[1000][1005];
static int macros;
static int active[1000];
static int failed;

static int find(const char *w) {
    for (int i = 0; i < macros; i++) if (strcmp(name[i], w) == 0) return i;
    return -1;
}

static void expand(const char *w, int first) {
    int at = find(w);
    if (at < 0) { printf(first ? "%s" : " %s", w); return; }
    if (active[at]) { failed = 1; return; }
    active[at] = 1;
    /* Tokenised by hand: strtok keeps global state, so a nested expansion
       would destroy the scan its own caller is in the middle of. */
    char copy[1005];
    strcpy(copy, body[at]);
    int f = first, start = 0;
    for (int i = 0; !failed; i++) {
        if (copy[i] != ' ' && copy[i] != '\0') continue;
        char save = copy[i];
        copy[i] = '\0';
        if (i > start) { expand(copy + start, f); f = 0; }
        start = i + 1;
        if (save == '\0') break;
    }
    active[at] = 0;
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    getchar();
    for (int i = 0; i < n; i++) {
        char line[1100];
        if (!fgets(line, sizeof line, stdin)) break;
        line[strcspn(line, "\n")] = '\0';
        char *sp = strchr(line, ' ');
        if (!sp) continue;
        *sp = '\0';
        strcpy(name[macros], line);
        strcpy(body[macros], sp + 1);
        macros++;
    }
    char target[64];
    if (scanf("%63s", target) != 1) return 1;
    expand(target, 1);
    if (failed) printf("cycle");
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var macros = map[string]string{}
var active = map[string]bool{}

func expand(w string) (string, bool) {
	body, ok := macros[w]
	if !ok {
		return w, true
	}
	if active[w] {
		return "", false
	}
	active[w] = true
	var parts []string
	for _, p := range strings.Split(body, " ") {
		r, ok := expand(p)
		if !ok {
			return "", false
		}
		parts = append(parts, r)
	}
	delete(active, w)
	return strings.Join(parts, " "), true
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		line := strings.TrimRight(raw, "\r\n")
		if k := strings.Index(line, " "); k > 0 {
			macros[line[:k]] = line[k+1:]
		}
	}
	var target string
	fmt.Fscan(in, &target)
	if out, ok := expand(target); ok {
		fmt.Println(out)
	} else {
		fmt.Println("cycle")
	}
}
`,
			"python": `import sys

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    macros = {}
    for line in raw[1:1 + n]:
        if " " in line:
            k, v = line.split(" ", 1)
            macros[k] = v
    target = raw[1 + n].strip()
    active = set()

    def expand(w):
        if w in active:
            return None
        if w not in macros:
            return w
        active.add(w)
        parts = []
        for p in macros[w].split(" "):
            r = expand(p)
            if r is None:
                return None
            parts.append(r)
        active.discard(w)
        return " ".join(parts)

    out = expand(target)
    print("cycle" if out is None else out)

main()
`},
	},
	{
		slug: "text-edits", title: "Text edits", difficulty: "medium", topic: "Stacks",
		statement: `A terminal treats two characters as commands rather than text: ~#~ deletes
the character before it, and ~@~ deletes everything back to the nearest space,
leaving the space itself. Either command on an empty line does nothing.

**Input.** One line. An empty line is written as ~-~.

**Output.** What is actually on screen.

    Input      Output
    ab#c       ac

    Input      Output
    hi the@    hi`,
		timeLimitMs: 4000,
		starters:    triLine("build the screen a character at a time"),
		tests: []store.Check{
			exact("Backspace removes one character", "ab#c\n", "ac"),
			exact("The word command removes back to the space", "hi the@\n", "hi "),
			exact("Nothing typed", "-\n", ""),
			exact("Deleting from an empty line does nothing", "###a\n", "a"),
			exactHid("The word command on the first word clears the line", "hello@\n", ""),
			exactHid("The word command on an empty line", "@\n", ""),
			exactHid("Backspacing over a space", "a b#\n", "a "),
			exactHid("Several commands in a row", "abc##d\n", "ad"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    static char line[100005], out[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    if (strcmp(line, "-") == 0) { printf("\n"); return 0; }
    int top = 0;
    for (int i = 0; line[i]; i++) {
        char c = line[i];
        if (c == '#') { if (top > 0) top--; }
        else if (c == '@') { while (top > 0 && out[top - 1] != ' ') top--; }
        else out[top++] = c;
    }
    out[top] = '\0';
    printf("%s\n", out);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	raw, _ := in.ReadString('\n')
	line := strings.TrimRight(raw, "\r\n")
	if line == "-" {
		fmt.Println("")
		return
	}
	var out []byte
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '#':
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
		case '@':
			for len(out) > 0 && out[len(out)-1] != ' ' {
				out = out[:len(out)-1]
			}
		default:
			out = append(out, line[i])
		}
	}
	fmt.Println(string(out))
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    if line == "-":
        print("")
        return
    out = []
    for ch in line:
        if ch == "#":
            if out:
                out.pop()
        elif ch == "@":
            while out and out[-1] != " ":
                out.pop()
        else:
            out.append(ch)
    print("".join(out))

main()
`},
	},
	{
		slug: "expand-runs", title: "Expand runs", difficulty: "medium", topic: "Parsing",
		statement: `A compressed string writes a character followed by how many times it
repeats. A character with no number after it appears once, and the count can be
more than one digit. Digits only ever follow a character, never start the string.

**Input.** One line. An empty line is written as ~-~.

**Output.** The original text.

    Input      Output
    a3b2       aaabb

    Input      Output
    abc        abc`,
		timeLimitMs: 4000,
		starters:    triLine("read a character, then any digits that follow it"),
		tests: []store.Check{
			exact("Expands the runs", "a3b2\n", "aaabb"),
			exact("No numbers means no repeats", "abc\n", "abc"),
			exact("A multi-digit count", "a12\n", "aaaaaaaaaaaa"),
			exact("Nothing to decode", "-\n", ""),
			exactHid("A count of one is allowed to be written", "a1b1\n", "ab"),
			exactHid("Mixed counted and uncounted", "ab3c\n", "abbbc"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    if (strcmp(line, "-") == 0) { printf("\n"); return 0; }
    int n = strlen(line), i = 0;
    while (i < n) {
        char c = line[i++];
        int j = i;
        while (j < n && isdigit((unsigned char)line[j])) j++;
        long count = 1;
        if (j > i) { count = 0; for (int k = i; k < j; k++) count = count * 10 + (line[k] - '0'); }
        for (long k = 0; k < count; k++) putchar(c);
        i = j;
    }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	raw, _ := in.ReadString('\n')
	line := strings.TrimRight(raw, "\r\n")
	if line == "-" {
		fmt.Println("")
		return
	}
	var b strings.Builder
	for i := 0; i < len(line); {
		c := line[i]
		i++
		j := i
		for j < len(line) && line[j] >= '0' && line[j] <= '9' {
			j++
		}
		count := 1
		if j > i {
			count, _ = strconv.Atoi(line[i:j])
		}
		b.WriteString(strings.Repeat(string(c), count))
		i = j
	}
	fmt.Println(b.String())
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    if line == "-":
        print("")
        return
    out, i = [], 0
    while i < len(line):
        c = line[i]
        i += 1
        j = i
        while j < len(line) and line[j].isdigit():
            j += 1
        count = int(line[i:j]) if j > i else 1
        out.append(c * count)
        i = j
    print("".join(out))

main()
`},
	},
}

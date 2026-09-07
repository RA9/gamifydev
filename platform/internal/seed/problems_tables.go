package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// tableProblems: counting, looking up, and the tie-breaks that decide what "the
// most" means. Inputs are bounded so the C answer can use a fixed array and a
// sort rather than a hash map written from scratch — the lesson is the rule,
// not the container.
var tableProblems = []seedProblem{
	{
		slug: "loot-tally", title: "Loot tally", difficulty: "easy", topic: "Dictionaries",
		statement: `After a raid you have a pile of drops.

**Input.** ~n threshold~ on the first line, then ~n~ item names (single words,
at most 1000 drops).

**Output.** Every item that dropped at least ~threshold~ times, as ~item count~
on its own line. Ordered by count, highest first; items with the same count are
ordered alphabetically.

    Input                    Output
    5 2                      ore 3
    ore gem ore ore gem      gem 2`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n, threshold;
    if (scanf("%d %d", &n, &threshold) != 2) return 1;
    static char name[1000][32];
    static int count[1000];
    int kinds = 0;
    for (int i = 0; i < n; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        /* count this drop */
    }
    /* print the kinds at or above the threshold, best first */
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, threshold int
	fmt.Fscan(in, &n, &threshold)
	count := map[string]int{}
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		_ = w
		// count this drop
	}
	_ = sort.Strings
	// print the kinds at or above the threshold, best first
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n, threshold = int(data[0]), int(data[1])
    drops = data[2:2 + n]
    # count the drops, then print the ones at or above the threshold
    
main()
`),
		tests: []store.Check{
			tokens("Orders by count", "5 2\nore gem ore ore gem\n", "ore 3 gem 2"),
			tokens("Drops below the threshold are left out", "3 2\nore gem ore\n", "ore 2"),
			tokens("Nothing dropped", "0 1\n", ""),
			hid2("A tie is broken alphabetically", "2 1\nb a\n", "a 1 b 1"),
			hid2("A threshold nothing reaches", "2 5\na a\n", ""),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n, threshold;
    if (scanf("%d %d", &n, &threshold) != 2) return 1;
    static char name[1000][32];
    static int count[1000];
    int kinds = 0;
    for (int i = 0; i < n; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        int at = -1;
        for (int j = 0; j < kinds; j++) if (strcmp(name[j], w) == 0) at = j;
        if (at < 0) { at = kinds++; strcpy(name[at], w); count[at] = 0; }
        count[at]++;
    }
    for (;;) {
        int best = -1;
        for (int j = 0; j < kinds; j++) {
            if (count[j] < threshold) continue;
            if (best < 0 || count[j] > count[best] ||
                (count[j] == count[best] && strcmp(name[j], name[best]) < 0)) best = j;
        }
        if (best < 0) break;
        printf("%s %d\n", name[best], count[best]);
        count[best] = -1;
    }
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, threshold int
	fmt.Fscan(in, &n, &threshold)
	count := map[string]int{}
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		count[w]++
	}
	var kinds []string
	for k, v := range count {
		if v >= threshold {
			kinds = append(kinds, k)
		}
	}
	sort.Slice(kinds, func(i, j int) bool {
		if count[kinds[i]] != count[kinds[j]] {
			return count[kinds[i]] > count[kinds[j]]
		}
		return kinds[i] < kinds[j]
	})
	for _, k := range kinds {
		fmt.Fprintln(out, k, count[k])
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, threshold = int(data[0]), int(data[1])
    counts = {}
    for d in data[2:2 + n]:
        counts[d] = counts.get(d, 0) + 1
    kinds = sorted((k for k, v in counts.items() if v >= threshold),
                   key=lambda k: (-counts[k], k))
    print("\n".join(f"{k} {counts[k]}" for k in kinds))

main()
`},
	},
	{
		slug: "unmatched-item", title: "The unmatched item", difficulty: "easy", topic: "Dictionaries",
		statement: `A warehouse scans every crate on the way in and on the way out. Exactly
one item does not balance: it appears a different number of times in the two
lists.

**Input.** ~a b~ on the first line — how many scans in, and how many out. Then
~a~ inbound names, then ~b~ outbound names. At most 1000 of each.

**Output.** The item that does not balance.

    Input      Output
    3 2        b
    a b c
    a c`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int a, b;
    if (scanf("%d %d", &a, &b) != 2) return 1;
    static char name[2000][32];
    static int delta[2000];
    int kinds = 0;
    /* add one for each inbound, take one for each outbound, print the odd one */
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
	var a, b int
	fmt.Fscan(in, &a, &b)
	delta := map[string]int{}
	// add one for each inbound, take one for each outbound, print the odd one
	_ = delta
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    a, b = int(data[0]), int(data[1])
    inbound = data[2:2 + a]
    outbound = data[2 + a:2 + a + b]
    # add one for each inbound, take one for each outbound, print the odd one
    print("")

main()
`),
		tests: []store.Check{
			tokens("Finds the item that never left", "3 2\na b c\na c\n", "b"),
			tokens("Order does not matter", "2 1\nx y\ny\n", "x"),
			tokens("An item can go out without coming in", "0 1\nz\n", "z"),
			hid2("Repeats count, not just presence", "3 2\na a b\na b\n", "a"),
			hid2("The odd one out is buried in the middle", "4 3\np q r s\ns p q\n", "r"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int a, b;
    if (scanf("%d %d", &a, &b) != 2) return 1;
    static char name[2000][32];
    static int delta[2000];
    int kinds = 0;
    for (int i = 0; i < a + b; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        int at = -1;
        for (int j = 0; j < kinds; j++) if (strcmp(name[j], w) == 0) at = j;
        if (at < 0) { at = kinds++; strcpy(name[at], w); delta[at] = 0; }
        delta[at] += (i < a) ? 1 : -1;
    }
    for (int j = 0; j < kinds; j++) if (delta[j] != 0) { printf("%s\n", name[j]); return 0; }
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	fmt.Fscan(in, &a, &b)
	delta := map[string]int{}
	order := []string{}
	for i := 0; i < a+b; i++ {
		var w string
		fmt.Fscan(in, &w)
		if _, seen := delta[w]; !seen {
			order = append(order, w)
		}
		if i < a {
			delta[w]++
		} else {
			delta[w]--
		}
	}
	for _, w := range order {
		if delta[w] != 0 {
			fmt.Println(w)
			return
		}
	}
	fmt.Println("")
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    a, b = int(data[0]), int(data[1])
    delta = {}
    for x in data[2:2 + a]:
        delta[x] = delta.get(x, 0) + 1
    for x in data[2 + a:2 + a + b]:
        delta[x] = delta.get(x, 0) - 1
    for k, v in delta.items():
        if v != 0:
            print(k)
            return
    print("")

main()
`},
	},
	{
		slug: "complete-kits", title: "Complete kits", difficulty: "easy", topic: "Dictionaries",
		statement: `A kit is a recipe of parts.

**Input.** ~p k~ on the first line — how many part types you hold, and how many
the kit needs. Then ~p~ lines of ~part count~ (what you have), then ~k~ lines of
~part needed~ (what one kit takes).

**Output.** How many complete kits you can build.

    Input          Output
    2 2            4
    bolt 10
    nut 4
    bolt 2
    nut 1

You run out of nuts first. A part you have none of stops you completely, and a
kit that requires nothing is not a kit — print 0 for it.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int p, k;
    if (scanf("%d %d", &p, &k) != 2) return 1;
    static char have[1000][32];
    static int stock[1000];
    for (int i = 0; i < p; i++) if (scanf("%31s %d", have[i], &stock[i]) != 2) return 1;
    /* for each needed part, how many kits does the stock allow? */
    printf("%d\n", 0);
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
	var p, k int
	fmt.Fscan(in, &p, &k)
	stock := map[string]int{}
	for i := 0; i < p; i++ {
		var name string
		var n int
		fmt.Fscan(in, &name, &n)
		stock[name] = n
	}
	// for each needed part, how many kits does the stock allow?
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    p, k = int(data[0]), int(data[1])
    stock = {data[2 + 2 * i]: int(data[3 + 2 * i]) for i in range(p)}
    base = 2 + 2 * p
    # for each needed part, how many kits does the stock allow?
    print(0)

main()
`),
		tests: []store.Check{
			tokens("The scarcest part decides", "2 2\nbolt 10\nnut 4\nbolt 2\nnut 1\n", "4"),
			tokens("A missing part means none at all", "1 2\nbolt 10\nbolt 1\nnut 1\n", "0"),
			tokens("A kit of nothing is not a kit", "1 0\nbolt 10\n", "0"),
			hid2("Spare parts you do not need are ignored", "2 1\nbolt 4\nglue 99\nbolt 2\n", "2"),
			hid2("Not quite enough for one", "1 1\nbolt 1\nbolt 2\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int p, k;
    if (scanf("%d %d", &p, &k) != 2) return 1;
    static char have[1000][32];
    static int stock[1000];
    for (int i = 0; i < p; i++) if (scanf("%31s %d", have[i], &stock[i]) != 2) return 1;
    if (k == 0) { printf("0\n"); return 0; }
    int best = -1;
    for (int i = 0; i < k; i++) {
        char part[32];
        int need;
        if (scanf("%31s %d", part, &need) != 2) return 1;
        int got = 0;
        for (int j = 0; j < p; j++) if (strcmp(have[j], part) == 0) got = stock[j];
        int kits = got / need;
        if (best < 0 || kits < best) best = kits;
    }
    printf("%d\n", best);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, k int
	fmt.Fscan(in, &p, &k)
	stock := map[string]int{}
	for i := 0; i < p; i++ {
		var name string
		var n int
		fmt.Fscan(in, &name, &n)
		stock[name] = n
	}
	if k == 0 {
		fmt.Println(0)
		return
	}
	best := -1
	for i := 0; i < k; i++ {
		var part string
		var need int
		fmt.Fscan(in, &part, &need)
		kits := stock[part] / need
		if best < 0 || kits < best {
			best = kits
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    p, k = int(data[0]), int(data[1])
    stock = {data[2 + 2 * i]: int(data[3 + 2 * i]) for i in range(p)}
    base = 2 + 2 * p
    if k == 0:
        print(0)
        return
    best = None
    for i in range(k):
        part, need = data[base + 2 * i], int(data[base + 1 + 2 * i])
        kits = stock.get(part, 0) // need
        best = kits if best is None else min(best, kits)
    print(best)

main()
`},
	},
	{
		slug: "attendance-streaks", title: "Attendance streaks", difficulty: "easy", topic: "Dictionaries",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~name days~ where ~days~ is a
string of ~P~ (present) and ~A~ (absent), one character per day in order. A
student with no recorded days has ~-~ instead.

**Output.** One line per student, in the order given: ~name streak~, where
~streak~ is their longest run of consecutive days present.

    Input        Output
    1            ada 2
    ada PPAP`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char name[32], days[1005];
        if (scanf("%31s %1000s", name, days) != 2) return 1;
        /* longest run of P, printing name and streak */
        printf("%s %d\n", name, 0);
    }
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var name, days string
		fmt.Fscan(in, &name, &days)
		// longest run of P
		fmt.Fprintln(out, name, 0)
	}
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    for i in range(n):
        name, days = data[1 + 2 * i], data[2 + 2 * i]
        # longest run of P
        print(name, 0)

main()
`),
		tests: []store.Check{
			tokens("Finds the longest run", "1\nada PPAP\n", "ada 2"),
			tokens("Never present is a streak of nothing", "1\nbob AA\n", "bob 0"),
			tokens("No students at all", "0\n", ""),
			hid2("Every student is reported", "2\na P\nb -\n", "a 1 b 0"),
			hid2("A run that reaches the last day counts", "1\nc APP\n", "c 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char name[32], days[1005];
        if (scanf("%31s %1000s", name, days) != 2) return 1;
        int best = 0, run = 0;
        for (int j = 0; days[j]; j++) {
            run = days[j] == 'P' ? run + 1 : 0;
            if (run > best) best = run;
        }
        printf("%s %d\n", name, best);
    }
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var name, days string
		fmt.Fscan(in, &name, &days)
		best, run := 0, 0
		for _, c := range days {
			if c == 'P' {
				run++
			} else {
				run = 0
			}
			if run > best {
				best = run
			}
		}
		fmt.Fprintln(out, name, best)
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    for i in range(n):
        name, days = data[1 + 2 * i], data[2 + 2 * i]
        best = run = 0
        for c in days:
            run = run + 1 if c == "P" else 0
            best = max(best, run)
        print(name, best)

main()
`},
	},
	{
		slug: "hold-queue", title: "Hold queue", difficulty: "medium", topic: "Dictionaries",
		statement: `A library records holds in the order they were placed, across every title
at once.

**Input.** A count ~n~, then ~n~ lines of ~title person~. The last line names
the ~title~ and ~person~ being asked about.

**Output.** That person's place in the queue for that title, counting from 1, or
0 if they are not waiting for it.

    Input              Output
    3                  2
    Dune ada
    Emma bob
    Dune bob
    Dune bob

Nobody queues twice for the same book: a repeat request from someone already in
that queue is ignored and does not move anyone.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char t[1000][32], p[1000][32];
    for (int i = 0; i < n; i++) if (scanf("%31s %31s", t[i], p[i]) != 2) return 1;
    char title[32], person[32];
    if (scanf("%31s %31s", title, person) != 2) return 1;
    /* walk the holds for that title, skipping repeats */
    printf("%d\n", 0);
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
	titles := make([]string, n)
	people := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &titles[i], &people[i])
	}
	var title, person string
	fmt.Fscan(in, &title, &person)
	// walk the holds for that title, skipping repeats
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    holds = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    title, person = data[1 + 2 * n], data[2 + 2 * n]
    # walk the holds for that title, skipping repeats
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Counts only holds on that title", "3\nDune ada\nEmma bob\nDune bob\nDune bob\n", "2"),
			tokens("First in line", "1\nDune ada\nDune ada\n", "1"),
			tokens("Not waiting at all", "1\nDune ada\nDune zoe\n", "0"),
			tokens("Nobody is waiting for anything", "0\nDune ada\n", "0"),
			hid2("A repeat hold does not move the queue", "3\nD a\nD a\nD b\nD b\n", "2"),
			hid2("A repeat hold keeps the original place", "3\nD a\nD b\nD a\nD a\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char t[1000][32], p[1000][32];
    for (int i = 0; i < n; i++) if (scanf("%31s %31s", t[i], p[i]) != 2) return 1;
    char title[32], person[32];
    if (scanf("%31s %31s", title, person) != 2) return 1;
    static char seen[1000][32];
    int seenN = 0, place = 0;
    for (int i = 0; i < n; i++) {
        if (strcmp(t[i], title) != 0) continue;
        int dup = 0;
        for (int j = 0; j < seenN; j++) if (strcmp(seen[j], p[i]) == 0) dup = 1;
        if (dup) continue;
        strcpy(seen[seenN++], p[i]);
        place++;
        if (strcmp(p[i], person) == 0) { printf("%d\n", place); return 0; }
    }
    printf("0\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	titles := make([]string, n)
	people := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &titles[i], &people[i])
	}
	var title, person string
	fmt.Fscan(in, &title, &person)
	seen := map[string]bool{}
	place := 0
	for i := 0; i < n; i++ {
		if titles[i] != title || seen[people[i]] {
			continue
		}
		seen[people[i]] = true
		place++
		if people[i] == person {
			fmt.Println(place)
			return
		}
	}
	fmt.Println(0)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    holds = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    title, person = data[1 + 2 * n], data[2 + 2 * n]
    seen, place = set(), 0
    for t, p in holds:
        if t != title or p in seen:
            continue
        seen.add(p)
        place += 1
        if p == person:
            print(place)
            return
    print(0)

main()
`},
	},
	{
		slug: "tag-cloud-sizes", title: "Tag cloud sizes", difficulty: "medium", topic: "Dictionaries",
		statement: `A tag cloud draws each tag at one of five sizes. The rarest tag is size 1,
the most common is size 5, and everything between is spread evenly:

    size = 1 + (count - lowest) * 4 / (highest - lowest)

using whole-number division. If every tag has the same count there is no range
to spread across, and every tag is drawn at size 3.

**Input.** A count ~n~, then ~n~ lines of ~tag count~.

**Output.** ~tag size~ per line, in the order given.

    Input       Output
    2           go 1
    go 1        python 5
    python 5`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char tag[1000][32];
    static int count[1000];
    for (int i = 0; i < n; i++) if (scanf("%31s %d", tag[i], &count[i]) != 2) return 1;
    /* find the range, then size each tag */
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	tags := make([]string, n)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &tags[i], &counts[i])
	}
	// find the range, then size each tag
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    tags = [(data[1 + 2 * i], int(data[2 + 2 * i])) for i in range(n)]
    # find the range, then size each tag

main()
`),
		tests: []store.Check{
			tokens("The extremes take the extreme sizes", "2\ngo 1\npython 5\n", "go 1 python 5"),
			tokens("All equal means all middling", "2\na 3\nb 3\n", "a 3 b 3"),
			tokens("No tags, no sizes", "0\n", ""),
			hid2("A lone tag has no range either", "1\nsolo 9\n", "solo 3"),
			hid2("The middle lands where the formula says", "3\na 0\nb 2\nc 4\n", "a 1 b 3 c 5"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char tag[1000][32];
    static int count[1000];
    for (int i = 0; i < n; i++) if (scanf("%31s %d", tag[i], &count[i]) != 2) return 1;
    if (n == 0) return 0;
    int lo = count[0], hi = count[0];
    for (int i = 1; i < n; i++) {
        if (count[i] < lo) lo = count[i];
        if (count[i] > hi) hi = count[i];
    }
    for (int i = 0; i < n; i++) {
        int size = (lo == hi) ? 3 : 1 + (count[i] - lo) * 4 / (hi - lo);
        printf("%s %d\n", tag[i], size);
    }
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	tags := make([]string, n)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &tags[i], &counts[i])
	}
	if n == 0 {
		return
	}
	lo, hi := counts[0], counts[0]
	for _, c := range counts {
		if c < lo {
			lo = c
		}
		if c > hi {
			hi = c
		}
	}
	for i := range tags {
		size := 3
		if lo != hi {
			size = 1 + (counts[i]-lo)*4/(hi-lo)
		}
		fmt.Fprintln(out, tags[i], size)
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    tags = [(data[1 + 2 * i], int(data[2 + 2 * i])) for i in range(n)]
    if not tags:
        return
    lo = min(c for _, c in tags)
    hi = max(c for _, c in tags)
    for name, c in tags:
        size = 3 if lo == hi else 1 + (c - lo) * 4 // (hi - lo)
        print(name, size)

main()
`},
	},
	{
		slug: "keep-the-last", title: "Keep the last", difficulty: "medium", topic: "Dictionaries",
		statement: `A playlist has picked up duplicates. Remove them — but keep the **last**
time each track appears, not the first, and leave the survivors in the order
those last appearances happen.

**Input.** A count ~n~, then ~n~ track names.

**Output.** The deduplicated playlist, space-separated on one line.

    Input        Output
    4            b a c
    a b a c

The first ~a~ goes; the second one stays where it is.`,
		timeLimitMs: 4000,
		starters:    triWords("keep only the last appearance of each track"),
		tests: []store.Check{
			tokens("Keeps the later copy in its own place", "4\na b a c\n", "b a c"),
			tokens("Nothing to remove", "2\nx y\n", "x y"),
			tokens("An empty playlist", "0\n", ""),
			hid2("All the same leaves one", "3\na a a\n", "a"),
			hid2("Several duplicates at once", "4\na b b a\n", "b a"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char w[100000][32];
    for (int i = 0; i < n; i++) if (scanf("%31s", w[i]) != 1) return 1;
    int first = 1;
    for (int i = 0; i < n; i++) {
        int later = 0;
        for (int j = i + 1; j < n; j++) if (strcmp(w[i], w[j]) == 0) later = 1;
        if (later) continue;
        printf(first ? "%s" : " %s", w[i]);
        first = 0;
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
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	w := make([]string, n)
	for i := range w {
		fmt.Fscan(in, &w[i])
	}
	last := map[string]int{}
	for i, s := range w {
		last[s] = i
	}
	var out []string
	for i, s := range w {
		if last[s] == i {
			out = append(out, s)
		}
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    tracks = data[1:1 + n]
    last = {t: i for i, t in enumerate(tracks)}
    print(" ".join(t for i, t in enumerate(tracks) if last[t] == i))

main()
`},
	},
	{
		slug: "roster-diff", title: "Roster diff", difficulty: "easy", topic: "Sets",
		statement: `Compare two team rosters.

**Input.** ~a b~ on the first line, then ~a~ names for the old roster and ~b~
for the new.

**Output.** Three lines: who joined, who left, and who stayed — each sorted
alphabetically and space-separated. An empty group prints an empty line.

    Input        Output
    2 2          cy
    ada bob      ada
    bob cy       bob

A name repeated in a roster is still just one person.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static int cmp(const void *x, const void *y) { return strcmp((const char *)x, (const char *)y); }

int main(void) {
    int a, b;
    if (scanf("%d %d", &a, &b) != 2) return 1;
    static char before[1000][32], after[1000][32];
    for (int i = 0; i < a; i++) if (scanf("%31s", before[i]) != 1) return 1;
    for (int i = 0; i < b; i++) if (scanf("%31s", after[i]) != 1) return 1;
    /* print joined, left, stayed — each sorted */
    printf("\n\n\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	fmt.Fscan(in, &a, &b)
	before := map[string]bool{}
	after := map[string]bool{}
	for i := 0; i < a; i++ {
		var s string
		fmt.Fscan(in, &s)
		before[s] = true
	}
	for i := 0; i < b; i++ {
		var s string
		fmt.Fscan(in, &s)
		after[s] = true
	}
	_ = sort.Strings
	_ = strings.Join
	// print joined, left, stayed — each sorted
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    a, b = int(data[0]), int(data[1])
    before = set(data[2:2 + a])
    after = set(data[2 + a:2 + a + b])
    # print joined, left, stayed — each sorted
    print("")
    print("")
    print("")

main()
`),
		tests: []store.Check{
			exact("Joined, left and stayed", "2 2\nada bob\nbob cy\n", "cy\nada\nbob"),
			exact("Nobody moved", "1 1\na\na\n", "\n\na"),
			exact("A team from nothing", "0 2\na b\n", "a b"),
			exactHid("Results are sorted, not in roster order", "2 2\nz y\ny z\n", "\n\ny z"),
			exactHid("A repeated name is one person", "2 0\na a\n", "\na"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char before[1000][32], after[1000][32];
static int na, nb;

static int has(char list[1000][32], int n, const char *w) {
    for (int i = 0; i < n; i++) if (strcmp(list[i], w) == 0) return 1;
    return 0;
}

static int cmpstr(const void *x, const void *y) {
    return strcmp((const char *)x, (const char *)y);
}

static void emit(char list[1000][32], int n, char other[1000][32], int m, int want) {
    static char picked[1000][32];
    int k = 0;
    for (int i = 0; i < n; i++) {
        if (has(list, i, list[i])) continue;
        if (has(other, m, list[i]) != want) continue;
        strcpy(picked[k++], list[i]);
    }
    qsort(picked, k, 32, cmpstr);
    for (int i = 0; i < k; i++) printf(i ? " %s" : "%s", picked[i]);
    printf("\n");
}

int main(void) {
    if (scanf("%d %d", &na, &nb) != 2) return 1;
    for (int i = 0; i < na; i++) if (scanf("%31s", before[i]) != 1) return 1;
    for (int i = 0; i < nb; i++) if (scanf("%31s", after[i]) != 1) return 1;
    emit(after, nb, before, na, 0);
    emit(before, na, after, nb, 0);
    emit(after, nb, before, na, 1);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func readSet(in *bufio.Reader, n int) map[string]bool {
	s := map[string]bool{}
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		s[w] = true
	}
	return s
}

func emit(from, other map[string]bool, want bool) {
	var out []string
	for k := range from {
		if other[k] == want {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	fmt.Println(strings.Join(out, " "))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b int
	fmt.Fscan(in, &a, &b)
	before := readSet(in, a)
	after := readSet(in, b)
	emit(after, before, false)
	emit(before, after, false)
	emit(after, before, true)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    a, b = int(data[0]), int(data[1])
    before = set(data[2:2 + a])
    after = set(data[2 + a:2 + a + b])
    print(" ".join(sorted(after - before)))
    print(" ".join(sorted(before - after)))
    print(" ".join(sorted(after & before)))

main()
`},
	},
	{
		slug: "restock-list", title: "Restock list", difficulty: "easy", topic: "Dictionaries",
		statement: `**Input.** ~p s level~ on the first line. Then ~p~ lines of ~item count~
(what is in stock), then ~s~ item names (one sale each).

**Output.** The sorted names of everything left at or below ~level~,
space-separated.

    Input          Output
    2 2 0          jam
    tea 3
    jam 1
    tea jam

An item can go negative — that is a backorder, not an error — and an item sold
that was never in stock starts from zero.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int p, s, level;
    if (scanf("%d %d %d", &p, &s, &level) != 3) return 1;
    static char name[2000][32];
    static int left[2000];
    int kinds = 0;
    /* read stock, apply sales, print what is at or below the level */
    printf("\n");
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, s, level int
	fmt.Fscan(in, &p, &s, &level)
	left := map[string]int{}
	// read stock, apply sales, print what is at or below the level
	_, _ = sort.Strings, strings.Join
	_ = left
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    p, s, level = int(data[0]), int(data[1]), int(data[2])
    # read stock, apply sales, print what is at or below the level
    print("")

main()
`),
		tests: []store.Check{
			exact("Reports what ran down to the level", "2 2 0\ntea 3\njam 1\ntea jam\n", "jam"),
			exact("Nothing sold, nothing to restock", "1 0 0\ntea 3\n", ""),
			exact("Everything is low", "2 2 0\na 1\nb 1\na b\n", "a b"),
			exactHid("An item can go into backorder", "1 1 0\ntea 0\ntea\n", "tea"),
			exactHid("Selling something never stocked", "0 1 0\nghost\n", "ghost"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char name[2000][32];
static int left[2000];
static int kinds;

static int slot(const char *w) {
    for (int i = 0; i < kinds; i++) if (strcmp(name[i], w) == 0) return i;
    strcpy(name[kinds], w);
    left[kinds] = 0;
    return kinds++;
}

static int cmpstr(const void *x, const void *y) { return strcmp((const char *)x, (const char *)y); }

int main(void) {
    int p, s, level;
    if (scanf("%d %d %d", &p, &s, &level) != 3) return 1;
    for (int i = 0; i < p; i++) {
        char w[32]; int c;
        if (scanf("%31s %d", w, &c) != 2) return 1;
        left[slot(w)] = c;
    }
    for (int i = 0; i < s; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        left[slot(w)]--;
    }
    static char picked[2000][32];
    int k = 0;
    for (int i = 0; i < kinds; i++) if (left[i] <= level) strcpy(picked[k++], name[i]);
    qsort(picked, k, 32, cmpstr);
    for (int i = 0; i < k; i++) printf(i ? " %s" : "%s", picked[i]);
    printf("\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, s, level int
	fmt.Fscan(in, &p, &s, &level)
	left := map[string]int{}
	for i := 0; i < p; i++ {
		var w string
		var c int
		fmt.Fscan(in, &w, &c)
		left[w] = c
	}
	for i := 0; i < s; i++ {
		var w string
		fmt.Fscan(in, &w)
		left[w]--
	}
	var out []string
	for k, v := range left {
		if v <= level {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    p, s, level = int(data[0]), int(data[1]), int(data[2])
    left = {}
    for i in range(p):
        left[data[3 + 2 * i]] = int(data[4 + 2 * i])
    base = 3 + 2 * p
    for item in data[base:base + s]:
        left[item] = left.get(item, 0) - 1
    print(" ".join(sorted(k for k, v in left.items() if v <= level)))

main()
`},
	},
	{
		slug: "best-shared-day", title: "Best shared day", difficulty: "medium", topic: "Sets",
		statement: `**Input.** ~d p~ on the first line — how many candidate days, and how many
people. Then ~d~ day names in the order they should be considered. Then ~p~
lines, each a count ~k~ followed by ~k~ day names that person can make.

**Output.** The day the most people can make. If several tie, the one earliest
in the candidate list. If there are no days at all, print an empty line.

    Input              Output
    2 2                tue
    mon tue
    1 tue
    2 mon tue

Days nobody mentions are still candidates — with nobody available. Days outside
the candidate list are ignored.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int d, p;
    if (scanf("%d %d", &d, &p) != 2) return 1;
    static char day[1000][32];
    static int votes[1000];
    for (int i = 0; i < d; i++) if (scanf("%31s", day[i]) != 1) return 1;
    /* tally availability, then take the best day, ties going to the earlier one */
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
	var d, p int
	fmt.Fscan(in, &d, &p)
	days := make([]string, d)
	for i := range days {
		fmt.Fscan(in, &days[i])
	}
	// tally availability, then take the best day, ties going to the earlier one
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    d, p = int(data[0]), int(data[1])
    days = data[2:2 + d]
    # tally availability, then take the best day, ties going to the earlier one
    print("")

main()
`),
		tests: []store.Check{
			exact("Picks the most popular day", "2 2\nmon tue\n1 tue\n2 mon tue\n", "tue"),
			exact("A tie goes to the earlier day", "2 1\nmon tue\n2 mon tue\n", "mon"),
			exact("No candidate days", "0 0\n", ""),
			exactHid("A day nobody can make can still win", "2 0\nmon tue\n", "mon"),
			exactHid("Days outside the candidate list are ignored", "1 2\nmon\n1 sun\n1 mon\n", "mon"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int d, p;
    if (scanf("%d %d", &d, &p) != 2) return 1;
    static char day[1000][32];
    static int votes[1000];
    for (int i = 0; i < d; i++) { if (scanf("%31s", day[i]) != 1) return 1; votes[i] = 0; }
    for (int i = 0; i < p; i++) {
        int k;
        if (scanf("%d", &k) != 1) return 1;
        for (int j = 0; j < k; j++) {
            char w[32];
            if (scanf("%31s", w) != 1) return 1;
            for (int t = 0; t < d; t++) if (strcmp(day[t], w) == 0) votes[t]++;
        }
    }
    if (d == 0) { printf("\n"); return 0; }
    int best = 0;
    for (int i = 1; i < d; i++) if (votes[i] > votes[best]) best = i;
    printf("%s\n", day[best]);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var d, p int
	fmt.Fscan(in, &d, &p)
	days := make([]string, d)
	index := map[string]int{}
	for i := range days {
		fmt.Fscan(in, &days[i])
		index[days[i]] = i
	}
	votes := make([]int, d)
	for i := 0; i < p; i++ {
		var k int
		fmt.Fscan(in, &k)
		for j := 0; j < k; j++ {
			var w string
			fmt.Fscan(in, &w)
			if t, ok := index[w]; ok {
				votes[t]++
			}
		}
	}
	if d == 0 {
		fmt.Println("")
		return
	}
	best := 0
	for i := 1; i < d; i++ {
		if votes[i] > votes[best] {
			best = i
		}
	}
	fmt.Println(days[best])
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    d, p = int(data[0]), int(data[1])
    days = data[2:2 + d]
    votes = {day: 0 for day in days}
    pos = 2 + d
    for _ in range(p):
        k = int(data[pos]); pos += 1
        for w in data[pos:pos + k]:
            if w in votes:
                votes[w] += 1
        pos += k
    if not days:
        print("")
        return
    best = max(days, key=lambda day: (votes[day], -days.index(day)))
    print(best)

main()
`},
	},
	{
		slug: "pair-up", title: "Pair up", difficulty: "medium", topic: "Dictionaries",
		statement: `A buddy scheme pairs people off.

**Input.** A count ~n~, then ~n~ lines of ~a b~.

**Output.** ~invalid~ if any name turns up in more than one pair or is paired
with themselves. Otherwise every name and its buddy, one ~name buddy~ per line,
sorted by name.

    Input      Output
    1          ada bob
    ada bob    bob ada

    Input      Output
    2          invalid
    ada bob
    ada cy`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char who[2000][32], buddy[2000][32];
    int people = 0;
    /* build the pairing, or print invalid */
    return 0;
}
`, `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	buddy := map[string]string{}
	// build the pairing, or print invalid
	_, _ = buddy, sort.Strings
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    # build the pairing, or print invalid

main()
`),
		tests: []store.Check{
			tokens("Maps both directions", "1\nada bob\n", "ada bob bob ada"),
			tokens("Somebody in two pairs is invalid", "2\nada bob\nada cy\n", "invalid"),
			tokens("Nobody to pair is fine", "0\n", ""),
			hid2("Pairing with yourself is invalid", "1\nada ada\n", "invalid"),
			hid2("A clash on the second name is caught too", "2\na b\nc b\n", "invalid"),
			hid2("Several valid pairs", "2\na b\nc d\n", "a b b a c d d c"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char who[2000][32], buddy[2000][32];
static int people;

static int find(const char *w) {
    for (int i = 0; i < people; i++) if (strcmp(who[i], w) == 0) return i;
    return -1;
}

static int cmp(const void *x, const void *y) {
    return strcmp(*(const char *const *)x, *(const char *const *)y);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char a[32], b[32];
        if (scanf("%31s %31s", a, b) != 2) return 1;
        if (strcmp(a, b) == 0 || find(a) >= 0 || find(b) >= 0) { printf("invalid\n"); return 0; }
        strcpy(who[people], a); strcpy(buddy[people], b); people++;
        strcpy(who[people], b); strcpy(buddy[people], a); people++;
    }
    int idx[2000];
    for (int i = 0; i < people; i++) idx[i] = i;
    for (int i = 0; i < people; i++)
        for (int j = i + 1; j < people; j++)
            if (strcmp(who[idx[i]], who[idx[j]]) > 0) { int t = idx[i]; idx[i] = idx[j]; idx[j] = t; }
    for (int i = 0; i < people; i++) printf("%s %s\n", who[idx[i]], buddy[idx[i]]);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	buddy := map[string]string{}
	for i := 0; i < n; i++ {
		var a, b string
		fmt.Fscan(in, &a, &b)
		_, ha := buddy[a]
		_, hb := buddy[b]
		if a == b || ha || hb {
			fmt.Fprintln(out, "invalid")
			return
		}
		buddy[a], buddy[b] = b, a
	}
	names := make([]string, 0, len(buddy))
	for k := range buddy {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		fmt.Fprintln(out, k, buddy[k])
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    buddy = {}
    for i in range(n):
        a, b = data[1 + 2 * i], data[2 + 2 * i]
        if a == b or a in buddy or b in buddy:
            print("invalid")
            return
        buddy[a], buddy[b] = b, a
    for k in sorted(buddy):
        print(k, buddy[k])

main()
`},
	},
	{
		slug: "poll-winner", title: "Poll winner", difficulty: "medium", topic: "Dictionaries",
		statement: `**Input.** A count ~n~, then ~n~ votes as single words.

**Output.** The option with the most votes. When two options tie, the winner is
the one that was **voted for first** — the option that got on the board earliest,
not the one that finished the count first. An empty poll has no winner: print
~none~.

    Input                  Output
    4                      blue
    blue red red blue

Both have two votes, and blue was voted for first.`,
		timeLimitMs: 4000,
		starters:    triWords("count the votes and break ties by who was voted for first"),
		tests: []store.Check{
			tokens("The clear majority wins", "3\na b a\n", "a"),
			tokens("A tie goes to whoever was voted for first", "4\nblue red red blue\n", "blue"),
			tokens("Nobody voted", "0\n", "none"),
			tokens("A single vote decides it", "1\nonly\n", "only"),
			hid2("First vote, not first to reach the count", "4\nx y y x\n", "x"),
			hid2("A three-way tie still resolves", "3\nc b a\n", "c"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char name[100000][32];
    static int count[100000];
    int kinds = 0;
    for (int i = 0; i < n; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        int at = -1;
        for (int j = 0; j < kinds; j++) if (strcmp(name[j], w) == 0) at = j;
        if (at < 0) { at = kinds++; strcpy(name[at], w); count[at] = 0; }
        count[at]++;
    }
    if (kinds == 0) { printf("none\n"); return 0; }
    int best = 0;
    for (int i = 1; i < kinds; i++) if (count[i] > count[best]) best = i;
    printf("%s\n", name[best]);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	count := map[string]int{}
	var order []string
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		if _, seen := count[w]; !seen {
			order = append(order, w)
		}
		count[w]++
	}
	if len(order) == 0 {
		fmt.Println("none")
		return
	}
	best := order[0]
	for _, w := range order {
		if count[w] > count[best] {
			best = w
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    counts, order = {}, []
    for v in data[1:1 + n]:
        if v not in counts:
            counts[v] = 0
            order.append(v)
        counts[v] += 1
    if not order:
        print("none")
        return
    print(max(order, key=lambda k: (counts[k], -order.index(k))))

main()
`},
	},
	{
		slug: "wasted-bytes", title: "Wasted bytes", difficulty: "easy", topic: "Sets",
		statement: `A backup holds the same file under many names.

**Input.** A count ~n~, then ~n~ lines of ~name size digest~. Two uploads with
the same digest hold identical bytes.

**Output.** How many bytes are spent on copies — everything except the first
upload of each digest.

    Input             Output
    3                 10
    a.txt 10 x
    b.txt 10 x
    c.txt 5 y`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char seen[100000][40];
    int seenN = 0;
    long long wasted = 0;
    /* add the size of every upload whose digest was already seen */
    printf("%lld\n", wasted);
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
	seen := map[string]bool{}
	wasted := 0
	// add the size of every upload whose digest was already seen
	_ = seen
	fmt.Println(wasted)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    # add the size of every upload whose digest was already seen
    print(0)

main()
`),
		tests: []store.Check{
			tokens("A second copy is waste", "3\na.txt 10 x\nb.txt 10 x\nc.txt 5 y\n", "10"),
			tokens("All distinct files waste nothing", "2\na 1 p\nb 2 q\n", "0"),
			tokens("Nothing uploaded", "0\n", "0"),
			hid2("Three copies waste two of them", "3\na 4 x\nb 4 x\nc 4 x\n", "8"),
			hid2("The name has nothing to do with it", "2\nsame 7 p\nsame 7 q\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char seen[100000][40];
    int seenN = 0;
    long long wasted = 0;
    for (int i = 0; i < n; i++) {
        char name[64], digest[40];
        long long size;
        if (scanf("%63s %lld %39s", name, &size, digest) != 3) return 1;
        int dup = 0;
        for (int j = 0; j < seenN; j++) if (strcmp(seen[j], digest) == 0) dup = 1;
        if (dup) wasted += size; else strcpy(seen[seenN++], digest);
    }
    printf("%lld\n", wasted);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	seen := map[string]bool{}
	wasted := 0
	for i := 0; i < n; i++ {
		var name, digest string
		var size int
		fmt.Fscan(in, &name, &size, &digest)
		if seen[digest] {
			wasted += size
		} else {
			seen[digest] = true
		}
	}
	fmt.Println(wasted)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    seen, wasted = set(), 0
    for i in range(n):
        size, digest = int(data[2 + 3 * i]), data[3 + 3 * i]
        if digest in seen:
            wasted += size
        else:
            seen.add(digest)
    print(wasted)

main()
`},
	},
	{
		slug: "recipe-cost", title: "Recipe cost", difficulty: "easy", topic: "Dictionaries",
		statement: `**Input.** ~r p~ on the first line — how many ingredients the recipe needs,
and how many have prices. Then ~r~ lines of ~ingredient quantity~, then ~p~ lines
of ~ingredient price~.

**Output.** What the recipe costs. If any ingredient has no price you cannot cost
the recipe at all — print ~-1~ rather than quietly leaving it out.

    Input            Output
    2 2              65
    flour 2
    salt 1
    flour 30
    salt 5`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int r, p;
    if (scanf("%d %d", &r, &p) != 2) return 1;
    static char item[1000][32];
    static long long qty[1000];
    for (int i = 0; i < r; i++) if (scanf("%31s %lld", item[i], &qty[i]) != 2) return 1;
    /* look up each price; print -1 if one is missing */
    printf("%d\n", 0);
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
	var r, p int
	fmt.Fscan(in, &r, &p)
	// look up each price; print -1 if one is missing
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    r, p = int(data[0]), int(data[1])
    # look up each price; print -1 if one is missing
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Costs a full recipe", "2 2\nflour 2\nsalt 1\nflour 30\nsalt 5\n", "65"),
			tokens("An unpriced ingredient spoils the total", "1 1\nsaffron 1\nflour 30\n", "-1"),
			tokens("An empty recipe costs nothing", "0 1\nflour 30\n", "0"),
			hid2("A missing price is not treated as free", "2 1\nflour 1\ngold 1\nflour 30\n", "-1"),
			hid2("Prices for things you do not need are ignored", "1 2\nsalt 3\nsalt 2\npepper 99\n", "6"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int r, p;
    if (scanf("%d %d", &r, &p) != 2) return 1;
    static char item[1000][32], priced[1000][32];
    static long long qty[1000], price[1000];
    for (int i = 0; i < r; i++) if (scanf("%31s %lld", item[i], &qty[i]) != 2) return 1;
    for (int i = 0; i < p; i++) if (scanf("%31s %lld", priced[i], &price[i]) != 2) return 1;
    long long total = 0;
    for (int i = 0; i < r; i++) {
        int at = -1;
        for (int j = 0; j < p; j++) if (strcmp(priced[j], item[i]) == 0) at = j;
        if (at < 0) { printf("-1\n"); return 0; }
        total += price[at] * qty[i];
    }
    printf("%lld\n", total);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var r, p int
	fmt.Fscan(in, &r, &p)
	items := make([]string, r)
	qty := make([]int, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(in, &items[i], &qty[i])
	}
	price := map[string]int{}
	for i := 0; i < p; i++ {
		var k string
		var v int
		fmt.Fscan(in, &k, &v)
		price[k] = v
	}
	total := 0
	for i, it := range items {
		v, ok := price[it]
		if !ok {
			fmt.Println(-1)
			return
		}
		total += v * qty[i]
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    r, p = int(data[0]), int(data[1])
    recipe = [(data[2 + 2 * i], int(data[3 + 2 * i])) for i in range(r)]
    base = 2 + 2 * r
    price = {data[base + 2 * i]: int(data[base + 1 + 2 * i]) for i in range(p)}
    total = 0
    for item, q in recipe:
        if item not in price:
            print(-1)
            return
        total += price[item] * q
    print(total)

main()
`},
	},
}

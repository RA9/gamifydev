package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// basicProblems: loops, conditionals and arithmetic. The on-ramp.
//
// Every problem here reads its input and prints its answer, which is the one
// contract C, Go and Python all share — there is no namespace to inspect in a
// compiled binary. The starters carry the reading and the printing so what the
// learner writes is the part the problem is about.
var basicProblems = []seedProblem{
	{
		slug: "double-tiles", title: "Double tiles", difficulty: "easy", topic: "Loops",
		statement: `A board game scores each tile you land on, and some tiles are gold —
a gold tile is worth twice its face value.

**Input.** A count ~n~, then ~n~ lines of ~value gold~, where ~gold~ is 1 for a
gold tile and 0 otherwise.

**Output.** The total score.

    Input        Output
    3            13
    3 0
    5 1
    0 0

An empty board scores 0.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long total = 0;
    for (int i = 0; i < n; i++) {
        int value, gold;
        if (scanf("%d %d", &value, &gold) != 2) return 1;
        /* add this tile to the total */
    }
    printf("%lld\n", total);
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
	total := 0
	for i := 0; i < n; i++ {
		var value, gold int
		fmt.Fscan(in, &value, &gold)
		// add this tile to the total
	}
	fmt.Println(total)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    total = 0
    for i in range(n):
        value, gold = int(data[1 + 2 * i]), int(data[2 + 2 * i])
        # add this tile to the total
    print(total)

main()
`),
		tests: []store.Check{
			tokens("Doubles the gold tiles", "3\n3 0\n5 1\n0 0\n", "13"),
			tokens("No gold means the plain total", "3\n1 0\n2 0\n3 0\n", "6"),
			tokens("An empty board scores nothing", "0\n", "0"),
			hid2("All gold doubles everything", "2\n4 1\n6 1\n", "20"),
			hid2("Handles negative tiles", "2\n-3 1\n4 0\n", "-2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long total = 0;
    for (int i = 0; i < n; i++) {
        int value, gold;
        if (scanf("%d %d", &value, &gold) != 2) return 1;
        total += gold ? 2LL * value : value;
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
	var n int
	fmt.Fscan(in, &n)
	total := 0
	for i := 0; i < n; i++ {
		var value, gold int
		fmt.Fscan(in, &value, &gold)
		if gold == 1 {
			total += 2 * value
		} else {
			total += value
		}
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    total = 0
    for i in range(n):
        value, gold = int(data[1 + 2 * i]), int(data[2 + 2 * i])
        total += value * 2 if gold else value
    print(total)

main()
`},
	},
	{
		slug: "tide-crossings", title: "Tide crossings", difficulty: "easy", topic: "Loops",
		statement: `A sensor records the water level once an hour. A crossing is any hour
where the level was strictly below the flood mark and the next reading is at or
above it.

**Input.** ~n mark~ on the first line, then ~n~ levels.

**Output.** How many crossings there were.

    Input        Output
    5 5          2
    1 2 5 4 6

Note what this does not count: staying above the mark isn't a crossing, and
falling back below it isn't either. Only the moment it goes up through.`,
		timeLimitMs: 4000,
		starters:    triCountWith("mark", "count the rises through the mark"),
		tests: []store.Check{
			tokens("Counts each rise through the mark", "5 5\n1 2 5 4 6\n", "2"),
			tokens("Staying above is not a crossing", "3 5\n6 7 8\n", "0"),
			tokens("A single reading cannot cross", "1 5\n9\n", "0"),
			tokens("No readings, no crossings", "0 5\n", "0"),
			hid2("Landing exactly on the mark counts", "2 5\n4 5\n", "1"),
			hid2("Falling back below does not count", "3 5\n6 1 6\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n, mark;
    if (scanf("%d %d", &n, &mark) != 2) return 1;
    int prev = 0, count = 0;
    for (int i = 0; i < n; i++) {
        int v;
        if (scanf("%d", &v) != 1) return 1;
        if (i > 0 && prev < mark && v >= mark) count++;
        prev = v;
    }
    printf("%d\n", count);
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
	var n, mark int
	fmt.Fscan(in, &n, &mark)
	values := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	count := 0
	for i := 1; i < n; i++ {
		if values[i-1] < mark && values[i] >= mark {
			count++
		}
	}
	fmt.Println(count)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, mark = int(data[0]), int(data[1])
    values = [int(x) for x in data[2:2 + n]]
    print(sum(1 for a, b in zip(values, values[1:]) if a < mark <= b))

main()
`},
	},
	{
		slug: "quiet-hours", title: "Quiet hours", difficulty: "easy", topic: "Loops",
		statement: `A noise monitor logs a reading each hour. An hour counts as restful only
if it belongs to a stretch of **three or more** consecutive hours that are all
at or below the limit — a single quiet hour between two loud ones is not rest.

**Input.** ~n limit~ on the first line, then ~n~ readings.

**Output.** How many hours qualify.

    Input          Output
    6 3            3
    1 1 1 9 2 2

The first three hours form a stretch of three, so all three count. The last two
are quiet but only two long, so neither does.`,
		timeLimitMs: 4000,
		starters:    triCountWith("limit", "count hours inside a run of three or more"),
		tests: []store.Check{
			tokens("Counts a stretch of three", "6 3\n1 1 1 9 2 2\n", "3"),
			tokens("A short stretch counts for nothing", "5 3\n1 1 9 1 1\n", "0"),
			tokens("An empty log has no restful hours", "0 3\n", "0"),
			tokens("A long stretch counts in full", "5 3\n0 0 0 0 0\n", "5"),
			hid2("A stretch running to the end still counts", "4 3\n9 1 1 1\n", "3"),
			hid2("Two separate stretches both count", "7 3\n1 1 1 9 1 1 1\n", "6"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n, limit;
    if (scanf("%d %d", &n, &limit) != 2) return 1;
    int run = 0, total = 0;
    for (int i = 0; i <= n; i++) {
        int v = limit + 1;
        if (i < n && scanf("%d", &v) != 1) return 1;
        if (v <= limit) {
            run++;
        } else {
            if (run >= 3) total += run;
            run = 0;
        }
    }
    printf("%d\n", total);
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
	var n, limit int
	fmt.Fscan(in, &n, &limit)
	run, total := 0, 0
	for i := 0; i <= n; i++ {
		v := limit + 1
		if i < n {
			fmt.Fscan(in, &v)
		}
		if v <= limit {
			run++
		} else {
			if run >= 3 {
				total += run
			}
			run = 0
		}
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, limit = int(data[0]), int(data[1])
    values = [int(x) for x in data[2:2 + n]]
    total = run = 0
    for r in values + [limit + 1]:
        if r <= limit:
            run += 1
        else:
            if run >= 3:
                total += run
            run = 0
    print(total)

main()
`},
	},
	{
		slug: "mind-changes", title: "Mind changes", difficulty: "easy", topic: "Loops",
		statement: `A committee votes over and over until it settles.

**Input.** A count ~n~, then ~n~ votes as single words.

**Output.** How many times a vote differed from the one before it.

    Input               Output
    4                   2
    yes yes no yes

Fewer than two votes means no changes.`,
		timeLimitMs: 4000,
		starters:    triWords("count how often a vote differs from the one before"),
		tests: []store.Check{
			tokens("Counts each switch", "4\nyes yes no yes\n", "2"),
			tokens("A unanimous run never changes", "3\nno no no\n", "0"),
			tokens("One vote cannot change", "1\nyes\n", "0"),
			tokens("No votes at all", "0\n", "0"),
			hid2("Every vote different", "4\na b c d\n", "3"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    char prev[32] = "", cur[32];
    int changes = 0;
    for (int i = 0; i < n; i++) {
        if (scanf("%31s", cur) != 1) return 1;
        if (i > 0 && strcmp(prev, cur) != 0) changes++;
        strcpy(prev, cur);
    }
    printf("%d\n", changes);
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
	words := make([]string, n)
	for i := range words {
		fmt.Fscan(in, &words[i])
	}
	changes := 0
	for i := 1; i < n; i++ {
		if words[i] != words[i-1] {
			changes++
		}
	}
	fmt.Println(changes)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    words = data[1:1 + n]
    print(sum(1 for a, b in zip(words, words[1:]) if a != b))

main()
`},
	},
	{
		slug: "stamp-card", title: "Stamp card", difficulty: "easy", topic: "Arithmetic",
		statement: `A coffee shop stamps your card for every drink. Every fifth drink is
free — the 5th, the 10th, the 15th, and so on, counting from the start.

**Input.** A count ~n~, then ~n~ prices in order.

**Output.** What you actually pay.

    Input          Output
    5              12
    3 3 3 3 3

The fifth drink is free whatever it costs, so ordering the expensive one fifth
is a strategy. That is not your problem to solve, only to charge for.`,
		timeLimitMs: 4000,
		starters:    triCount("charge every drink except each fifth one"),
		tests: []store.Check{
			tokens("The fifth drink is free", "5\n3 3 3 3 3\n", "12"),
			tokens("Fewer than five drinks are all charged", "2\n2 4\n", "6"),
			tokens("An empty order costs nothing", "0\n", "0"),
			tokens("The tenth is free too", "10\n1 1 1 1 1 1 1 1 1 1\n", "8"),
			hid2("The free drink is the fifth, whatever it costs", "5\n1 1 1 1 99\n", "4"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long total = 0;
    for (int i = 1; i <= n; i++) {
        int p;
        if (scanf("%d", &p) != 1) return 1;
        if (i % 5 != 0) total += p;
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
	var n int
	fmt.Fscan(in, &n)
	total := 0
	for i := 1; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		if i%5 != 0 {
			total += p
		}
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    prices = [int(x) for x in data[1:1 + n]]
    print(sum(p for i, p in enumerate(prices, 1) if i % 5))

main()
`},
	},
	{
		slug: "fare-cap", title: "Fare cap", difficulty: "easy", topic: "Arithmetic",
		statement: `A tram charges 2 for any trip, plus 1 for every kilometre beyond the
third. No single trip costs more than 8, and no day costs more than 20 however
much you ride.

**Input.** A count ~n~, then ~n~ whole-kilometre trips.

**Output.** What the day costs.

    Input        Output
    3            14
    2 5 20

That is 2 + 4 + 8: the last trip hits the per-trip cap.`,
		timeLimitMs: 4000,
		starters:    triCount("charge this trip, capped at 8, and cap the day at 20"),
		tests: []store.Check{
			tokens("Charges the base, the extra, and the trip cap", "3\n2 5 20\n", "14"),
			tokens("A short trip is just the base fare", "1\n1\n", "2"),
			tokens("No trips, no fare", "0\n", "0"),
			tokens("The daily cap holds", "10\n20 20 20 20 20 20 20 20 20 20\n", "20"),
			hid2("Exactly three kilometres adds nothing", "1\n3\n", "2"),
			hid2("The daily cap applies to the total, not each trip", "3\n8 8 8\n", "20"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int total = 0;
    for (int i = 0; i < n; i++) {
        int d;
        if (scanf("%d", &d) != 1) return 1;
        int fare = 2 + (d > 3 ? d - 3 : 0);
        if (fare > 8) fare = 8;
        total += fare;
    }
    if (total > 20) total = 20;
    printf("%d\n", total);
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
	total := 0
	for i := 0; i < n; i++ {
		var d int
		fmt.Fscan(in, &d)
		fare := 2
		if d > 3 {
			fare += d - 3
		}
		if fare > 8 {
			fare = 8
		}
		total += fare
	}
	if total > 20 {
		total = 20
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    trips = [int(x) for x in data[1:1 + n]]
    print(min(20, sum(min(8, 2 + max(0, d - 3)) for d in trips)))

main()
`},
	},
	{
		slug: "battery-log", title: "Battery log", difficulty: "easy", topic: "Simulation",
		statement: `A phone starts the day at 50%. The log is a list of whole-number
changes: negative for drain, positive for charge. The battery never goes below
0 or above 100 — an event that would push it past either end just leaves it
there.

**Input.** A count ~n~, then ~n~ signed changes.

**Output.** The level at the end of the day.

    Input          Output
    3              60
    -30 -40 60

The second event would take it to -20, so it stops at 0, and the charge lifts it
from there.`,
		timeLimitMs: 4000,
		starters:    triCount("apply this change and keep the level between 0 and 100"),
		tests: []store.Check{
			tokens("Clamps at empty before charging again", "3\n-30 -40 60\n", "60"),
			tokens("A quiet day changes nothing", "0\n", "50"),
			tokens("Cannot charge past full", "2\n80 20\n", "100"),
			hid2("Cannot drain past empty", "2\n-90 -90\n", "0"),
			hid2("Clamping happens per event, not at the end", "3\n-90 10 10\n", "20"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int level = 50;
    for (int i = 0; i < n; i++) {
        int e;
        if (scanf("%d", &e) != 1) return 1;
        level += e;
        if (level < 0) level = 0;
        if (level > 100) level = 100;
    }
    printf("%d\n", level);
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
	level := 50
	for i := 0; i < n; i++ {
		var e int
		fmt.Fscan(in, &e)
		level += e
		if level < 0 {
			level = 0
		}
		if level > 100 {
			level = 100
		}
	}
	fmt.Println(level)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    level = 50
    for x in data[1:1 + n]:
        level = max(0, min(100, level + int(x)))
    print(level)

main()
`},
	},
	{
		slug: "shift-pay", title: "Shift pay", difficulty: "easy", topic: "Arithmetic",
		statement: `A workshop pays 10 an hour for the first 40 hours of the week, 20 an
hour for the next 10, and 30 an hour beyond 50.

**Input.** A single number: the hours worked. Never negative.

**Output.** The week's wage.

    Input   Output
    45      500

    Input   Output
    52      660`,
		timeLimitMs: 4000,
		starters:    triScalars([]string{"hours"}, "work out the wage across the three bands"),
		tests: []store.Check{
			tokens("Pays the middle band", "45\n", "500"),
			tokens("Pays the top band", "52\n", "660"),
			tokens("A short week is all base rate", "10\n", "100"),
			tokens("No hours, no pay", "0\n", "0"),
			hid2("Exactly 40 hours stays on the base rate", "40\n", "400"),
			hid2("Exactly 50 hours fills the middle band", "50\n", "600"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long hours;
    if (scanf("%lld", &hours) != 1) return 1;
    long long base = hours < 40 ? hours : 40;
    long long mid = hours > 40 ? (hours < 50 ? hours : 50) - 40 : 0;
    long long top = hours > 50 ? hours - 50 : 0;
    printf("%lld\n", 10 * base + 20 * mid + 30 * top);
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var hours int
	fmt.Fscan(in, &hours)
	fmt.Println(10*min(hours, 40) + 20*max(0, min(hours, 50)-40) + 30*max(0, hours-50))
}
`,
			"python": `import sys

def main():
    hours = int(sys.stdin.read().split()[0])
    print(10 * min(hours, 40) + 20 * max(0, min(hours, 50) - 40) + 30 * max(0, hours - 50))

main()
`},
	},
	{
		slug: "matching-dice", title: "Matching dice", difficulty: "easy", topic: "Loops",
		statement: `A game is played with two dice.

**Input.** A count ~n~, then ~n~ lines of ~a b~ — the two dice for that roll.

**Output.** The length of the longest unbroken run of rolls where both dice
showed the same number.

    Input        Output
    4            2
    1 1
    2 2
    3 4
    5 5

No rolls means no run.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int best = 0, run = 0;
    for (int i = 0; i < n; i++) {
        int a, b;
        if (scanf("%d %d", &a, &b) != 2) return 1;
        /* extend or reset the run */
    }
    printf("%d\n", best);
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
	best, run := 0, 0
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		// extend or reset the run
	}
	fmt.Println(best)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    best = run = 0
    for i in range(n):
        a, b = int(data[1 + 2 * i]), int(data[2 + 2 * i])
        # extend or reset the run
    print(best)

main()
`),
		tests: []store.Check{
			tokens("Finds the longer of two runs", "4\n1 1\n2 2\n3 4\n5 5\n", "2"),
			tokens("No matches at all", "2\n1 2\n3 4\n", "0"),
			tokens("An empty game has no run", "0\n", "0"),
			hid2("A run at the very end counts", "4\n1 2\n6 6\n6 6\n6 6\n", "3"),
			hid2("Every roll matching is one long run", "5\n4 4\n4 4\n4 4\n4 4\n4 4\n", "5"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int best = 0, run = 0;
    for (int i = 0; i < n; i++) {
        int a, b;
        if (scanf("%d %d", &a, &b) != 2) return 1;
        run = (a == b) ? run + 1 : 0;
        if (run > best) best = run;
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
	var n int
	fmt.Fscan(in, &n)
	best, run := 0, 0
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a == b {
			run++
		} else {
			run = 0
		}
		if run > best {
			best = run
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    best = run = 0
    for i in range(n):
        a, b = int(data[1 + 2 * i]), int(data[2 + 2 * i])
        run = run + 1 if a == b else 0
        best = max(best, run)
    print(best)

main()
`},
	},
	{
		slug: "row-of-seedlings", title: "Row of seedlings", difficulty: "easy", topic: "Arithmetic",
		statement: `You plant the first seedling at position 0 of a bed, then one every
~gap~ centimetres, for as long as the next one still fits inside a bed
~length~ centimetres long. A seedling exactly at the far end fits.

**Input.** ~length gap~ on one line. The bed can be enormous.

**Output.** How many seedlings you plant.

    Input     Output
    10 3      4

Those sit at 0, 3, 6 and 9. A bed of negative length holds nothing.`,
		timeLimitMs: 4000,
		starters:    triScalars([]string{"length", "gap"}, "work out how many fit, without walking the bed"),
		tests: []store.Check{
			tokens("Fits four in a ten-centimetre bed", "10 3\n", "4"),
			tokens("A bed with no room still holds the first", "0 3\n", "1"),
			tokens("A negative bed holds nothing", "-1 3\n", "0"),
			hid2("A seedling exactly at the end fits", "9 3\n", "4"),
			hid2("A huge bed does not need a loop", "1000000000000 7\n", "142857142858"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long length, gap;
    if (scanf("%lld %lld", &length, &gap) != 2) return 1;
    printf("%lld\n", length < 0 ? 0 : length / gap + 1);
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
	var length, gap int
	fmt.Fscan(in, &length, &gap)
	if length < 0 {
		fmt.Println(0)
		return
	}
	fmt.Println(length/gap + 1)
}
`,
			"python": `import sys

def main():
    length, gap = (int(x) for x in sys.stdin.read().split()[:2])
    print(0 if length < 0 else length // gap + 1)

main()
`},
	},
	{
		slug: "working-days-ahead", title: "Working days ahead", difficulty: "medium", topic: "Arithmetic",
		statement: `A workshop is open every day except Sunday. Days of the week are
numbered 0 for Monday through 6 for Sunday.

**Input.** ~start days~ on one line: the day you begin on, and how many
**working** days to move forward. Sundays are stepped over and never counted.

**Output.** The day you land on.

    Input     Output
    0 6       0

Six working days from Monday is Tuesday, Wednesday, Thursday, Friday, Saturday,
then Monday again — Sunday is skipped.

~days~ can run past a trillion, so stepping one day at a time will not finish in
any language, and the count will not fit in a 32-bit integer. Six working days
is exactly one week.`,
		timeLimitMs: 4000,
		starters:    triScalars([]string{"start", "days"}, "jump whole weeks first, then walk the remainder"),
		tests: []store.Check{
			tokens("Six working days is a whole week", "0 6\n", "0"),
			tokens("One day forward from Monday", "0 1\n", "1"),
			tokens("Sunday is stepped over", "5 1\n", "0"),
			tokens("Going nowhere lands where you started", "3 0\n", "3"),
			tokens("Fast enough for a trillion days", "0 1000000000000\n", "4"),
			hid2("Starting on a Sunday still works", "6 1\n", "0"),
			hid2("Five days from Monday reaches Saturday", "0 5\n", "5"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long start, days;
    if (scanf("%lld %lld", &start, &days) != 2) return 1;
    long long d = (start + (days / 6) * 7) % 7;
    for (long long i = 0; i < days % 6; i++) {
        d = (d + 1) % 7;
        if (d == 6) d = 0;
    }
    printf("%lld\n", d);
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
	var start, days int
	fmt.Fscan(in, &start, &days)
	d := (start + (days/6)*7) % 7
	for i := 0; i < days%6; i++ {
		d = (d + 1) % 7
		if d == 6 {
			d = 0
		}
	}
	fmt.Println(d)
}
`,
			"python": `import sys

def main():
    start, days = (int(x) for x in sys.stdin.read().split()[:2])
    d = (start + (days // 6) * 7) % 7
    for _ in range(days % 6):
        d = (d + 1) % 7
        if d == 6:
            d = 0
    print(d)

main()
`},
		tooSlowIn: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long start, days;
    if (scanf("%lld %lld", &start, &days) != 2) return 1;
    long long d = start;
    while (days > 0) {
        d = (d + 1) % 7;
        if (d != 6) days--;
    }
    printf("%lld\n", d);
    return 0;
}
`,
			"python": `import sys

def main():
    start, days = (int(x) for x in sys.stdin.read().split()[:2])
    d = start
    while days > 0:
        d = (d + 1) % 7
        if d != 6:
            days -= 1
    print(d)

main()
`},
	},
	{
		slug: "biggest-swing", title: "Biggest swing", difficulty: "easy", topic: "Loops",
		statement: `**Input.** A count ~n~, then ~n~ readings.

**Output.** The largest jump between two readings that sit next to each other —
the size of the jump, never its direction.

    Input      Output
    3          7
    3 10 4

Fewer than two readings means there is no jump at all, which is 0.`,
		timeLimitMs: 4000,
		starters:    triCount("track the largest gap between neighbours"),
		tests: []store.Check{
			tokens("Finds the biggest jump", "3\n3 10 4\n", "7"),
			tokens("A drop counts the same as a rise", "2\n10 1\n", "9"),
			tokens("One reading has no jump", "1\n5\n", "0"),
			tokens("No readings at all", "0\n", "0"),
			hid2("The biggest jump may be at the end", "4\n1 2 3 40\n", "37"),
			hid2("A flat run swings by nothing", "3\n7 7 7\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int prev = 0, best = 0;
    for (int i = 0; i < n; i++) {
        int v;
        if (scanf("%d", &v) != 1) return 1;
        if (i > 0 && abs(v - prev) > best) best = abs(v - prev);
        prev = v;
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
	var n int
	fmt.Fscan(in, &n)
	prev, best := 0, 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if i > 0 {
			d := v - prev
			if d < 0 {
				d = -d
			}
			if d > best {
				best = d
			}
		}
		prev = v
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    values = [int(x) for x in data[1:1 + n]]
    print(max((abs(a - b) for a, b in zip(values, values[1:])), default=0))

main()
`},
	},
	{
		slug: "shortest-queue", title: "Shortest queue", difficulty: "easy", topic: "Loops",
		statement: `Each checkout lane holds some baskets, and each basket is a number of
items.

**Input.** A count ~lanes~. Then one line per lane: how many baskets it holds,
followed by that many basket sizes.

**Output.** The index of the lane with the fewest items in total, counting from
0. If two lanes tie, take the one further left.

    Input        Output
    3            2
    2 3 4
    1 10
    3 1 1 1

There is always at least one lane. A lane with nobody in it holds zero items,
which is as short as a lane gets.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int lanes;
    if (scanf("%d", &lanes) != 1) return 1;
    int best = 0;
    for (int i = 0; i < lanes; i++) {
        int baskets, total = 0;
        if (scanf("%d", &baskets) != 1) return 1;
        for (int j = 0; j < baskets; j++) {
            int items;
            if (scanf("%d", &items) != 1) return 1;
            total += items;
        }
        /* keep the smallest total, ties going to the earlier lane */
    }
    printf("%d\n", best);
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
	var lanes int
	fmt.Fscan(in, &lanes)
	best := 0
	for i := 0; i < lanes; i++ {
		var baskets, total int
		fmt.Fscan(in, &baskets)
		for j := 0; j < baskets; j++ {
			var items int
			fmt.Fscan(in, &items)
			total += items
		}
		// keep the smallest total, ties going to the earlier lane
	}
	fmt.Println(best)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    pos = 1
    best = 0
    for i in range(data[0]):
        baskets = data[pos]
        total = sum(data[pos + 1:pos + 1 + baskets])
        pos += 1 + baskets
        # keep the smallest total, ties going to the earlier lane
    print(best)

main()
`),
		tests: []store.Check{
			tokens("Picks the lane with fewest items", "3\n2 3 4\n1 10\n3 1 1 1\n", "2"),
			tokens("Counts items, not people", "2\n1 9\n2 1 1\n", "1"),
			tokens("An empty lane wins", "2\n1 5\n0\n", "1"),
			tokens("A single lane is the answer", "1\n1 8\n", "0"),
			hid2("A tie goes to the leftmost lane", "3\n2 2 2\n1 4\n2 1 3\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int lanes;
    if (scanf("%d", &lanes) != 1) return 1;
    int best = 0, bestTotal = 0;
    for (int i = 0; i < lanes; i++) {
        int baskets, total = 0;
        if (scanf("%d", &baskets) != 1) return 1;
        for (int j = 0; j < baskets; j++) {
            int items;
            if (scanf("%d", &items) != 1) return 1;
            total += items;
        }
        if (i == 0 || total < bestTotal) { best = i; bestTotal = total; }
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
	var lanes int
	fmt.Fscan(in, &lanes)
	best, bestTotal := 0, 0
	for i := 0; i < lanes; i++ {
		var baskets, total int
		fmt.Fscan(in, &baskets)
		for j := 0; j < baskets; j++ {
			var items int
			fmt.Fscan(in, &items)
			total += items
		}
		if i == 0 || total < bestTotal {
			best, bestTotal = i, total
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    pos, best, best_total = 1, 0, None
    for i in range(data[0]):
        baskets = data[pos]
        total = sum(data[pos + 1:pos + 1 + baskets])
        pos += 1 + baskets
        if best_total is None or total < best_total:
            best, best_total = i, total
    print(best)

main()
`},
	},
	{
		slug: "parking-meter", title: "Parking meter", difficulty: "easy", topic: "Arithmetic",
		statement: `The first 30 minutes are free. After that you pay 50 cents for every
15-minute block you have **started** — one minute into a block costs the whole
block. A day never costs more than 900 cents.

**Input.** A single number: the minutes parked.

**Output.** The charge in cents.

    Input   Output
    31      50

    Input   Output
    46      100

Watch the word *started*: 45 minutes is exactly one block past the free half
hour, and 46 minutes is into the second.`,
		timeLimitMs: 4000,
		starters:    triScalars([]string{"minutes"}, "charge each started block after the free half hour"),
		tests: []store.Check{
			tokens("The free half hour costs nothing", "30\n", "0"),
			tokens("One minute over starts a block", "31\n", "50"),
			tokens("A started second block is charged in full", "46\n", "100"),
			tokens("A long stay hits the daily cap", "10000\n", "900"),
			hid2("Exactly one block past the free time", "45\n", "50"),
			hid2("Arriving and leaving costs nothing", "0\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long minutes;
    if (scanf("%lld", &minutes) != 1) return 1;
    long long paid = minutes > 30 ? minutes - 30 : 0;
    long long blocks = (paid + 14) / 15;
    long long cost = 50 * blocks;
    printf("%lld\n", cost > 900 ? 900 : cost);
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
	var minutes int
	fmt.Fscan(in, &minutes)
	paid := 0
	if minutes > 30 {
		paid = minutes - 30
	}
	cost := 50 * ((paid + 14) / 15)
	if cost > 900 {
		cost = 900
	}
	fmt.Println(cost)
}
`,
			"python": `import sys

def main():
    minutes = int(sys.stdin.read().split()[0])
    paid = max(0, minutes - 30)
    print(min(900, 50 * -(-paid // 15)))

main()
`},
	},
}

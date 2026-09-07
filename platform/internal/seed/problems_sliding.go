package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// windowProblems: two pointers, sliding windows and prefix sums — the family of
// tricks whose whole point is not redoing work you have already done.
//
// The statements say where the naive approach stops scaling but do not promise
// a timeout, because discriminating a quadratic solution from a linear one over
// stdin needs megabytes of test data. The guidance is honest; the clock is not
// asked to enforce it.
var windowProblems = []seedProblem{
	{
		slug: "busiest-window", title: "Busiest window", difficulty: "medium", topic: "Sliding window",
		statement: `A turnstile counts arrivals each minute.

**Input.** ~n k~ on the first line, then ~n~ counts.

**Output.** The starting minute of the ~k~-minute stretch with the most
arrivals, counting from 0. If two stretches tie, the earlier one. If ~k~ is
bigger than the log, or not positive, there is no such stretch: print ~-1~.

    Input        Output
    4 2          1
    1 5 2 3

Adding up each window from scratch does the same additions over and over — a
window is the last one with one number swapped out and one swapped in.`,
		timeLimitMs: 4000,
		starters:    triCountWith("k", "slide the window, adjusting the running total"),
		tests: []store.Check{
			tokens("Finds the busiest stretch", "4 2\n1 5 2 3\n", "1"),
			tokens("A tie goes to the earlier stretch", "3 1\n2 2 2\n", "0"),
			tokens("A window longer than the log", "2 5\n1 2\n", "-1"),
			tokens("A window of nothing", "2 0\n1 2\n", "-1"),
			hid2("The whole log is one window", "2 2\n4 1\n", "0"),
			hid2("The busiest stretch is at the end", "3 1\n1 1 9\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, k;
    if (scanf("%d %d", &n, &k) != 2) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    for (int i = 0; i < n; i++) if (scanf("%lld", &v[i]) != 1) return 1;
    if (k <= 0 || k > n) { printf("-1\n"); return 0; }
    long long total = 0;
    for (int i = 0; i < k; i++) total += v[i];
    long long best = total;
    int at = 0;
    for (int i = k; i < n; i++) {
        total += v[i] - v[i - k];
        if (total > best) { best = total; at = i - k + 1; }
    }
    printf("%d\n", at);
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
	var n, k int
	fmt.Fscan(in, &n, &k)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	if k <= 0 || k > n {
		fmt.Println(-1)
		return
	}
	total := 0
	for i := 0; i < k; i++ {
		total += v[i]
	}
	best, at := total, 0
	for i := k; i < n; i++ {
		total += v[i] - v[i-k]
		if total > best {
			best, at = total, i-k+1
		}
	}
	fmt.Println(at)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, k = int(data[0]), int(data[1])
    v = [int(x) for x in data[2:2 + n]]
    if k <= 0 or k > n:
        print(-1)
        return
    total = sum(v[:k])
    best, at = total, 0
    for i in range(k, n):
        total += v[i] - v[i - k]
        if total > best:
            best, at = total, i - k + 1
    print(at)

main()
`},
	},
	{
		slug: "calm-stretch", title: "Calm stretch", difficulty: "easy", topic: "Sliding window",
		statement: `An instrument is calm while it never jumps by more than ~tol~ from one
reading to the next.

**Input.** ~n tol~ on the first line, then ~n~ readings.

**Output.** The length of the longest calm stretch. A single reading is a stretch
of one — it has not jumped anywhere. No readings at all is a stretch of zero.

    Input        Output
    3 3          2
    1 2 10`,
		timeLimitMs: 4000,
		starters:    triCountWith("tol", "extend the run while the jump is within tolerance"),
		tests: []store.Check{
			tokens("Finds the calm run", "3 3\n1 2 10\n", "2"),
			tokens("One reading is a stretch of one", "1 1\n7\n", "1"),
			tokens("No readings", "0 1\n", "0"),
			tokens("Everything is calm", "4 1\n1 2 3 4\n", "4"),
			hid2("A jump exactly at the tolerance is still calm", "2 3\n1 4\n", "2"),
			hid2("A drop counts as much as a rise", "3 3\n10 1 2\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, tol;
    if (scanf("%d %d", &n, &tol) != 2) return 1;
    long long prev = 0;
    int best = 0, run = 0;
    for (int i = 0; i < n; i++) {
        long long v;
        if (scanf("%lld", &v) != 1) return 1;
        run = (i > 0 && llabs(v - prev) <= tol) ? run + 1 : 1;
        if (run > best) best = run;
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
	var n, tol int
	fmt.Fscan(in, &n, &tol)
	best, run, prev := 0, 0, 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		d := v - prev
		if d < 0 {
			d = -d
		}
		if i > 0 && d <= tol {
			run++
		} else {
			run = 1
		}
		if run > best {
			best = run
		}
		prev = v
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, tol = int(data[0]), int(data[1])
    v = [int(x) for x in data[2:2 + n]]
    if not v:
        print(0)
        return
    best = run = 1
    for a, b in zip(v, v[1:]):
        run = run + 1 if abs(b - a) <= tol else 1
        best = max(best, run)
    print(best)

main()
`},
	},
	{
		slug: "balanced-stretch", title: "Balanced stretch", difficulty: "hard", topic: "Prefix sums",
		statement: `A market log records each day as ~U~ (up) or ~D~ (down).

**Input.** A count ~n~, then a string of ~n~ characters, each ~U~ or ~D~. A log
of no days is written as ~-~.

**Output.** The length of the longest stretch with exactly as many up days as
down days.

    Input     Output
    3         2
    UDU

Checking every stretch is quadratic. The way in: keep a running score, +1 for up
and -1 for down. Two days with the **same** running score have a balanced stretch
between them — so all you need to remember is the first day each score was seen.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char days[200005];
    if (scanf("%200000s", days) != 1) return 1;
    /* running score, remembering where each score was first seen */
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
	var days string
	fmt.Fscan(in, &n, &days)
	// running score, remembering where each score was first seen
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n, days = int(data[0]), data[1]
    # running score, remembering where each score was first seen
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Finds the balanced stretch", "3\nUDU\n", "2"),
			tokens("Nothing balances", "2\nUU\n", "0"),
			tokens("An empty log", "0\n-\n", "0"),
			tokens("The whole log balances", "4\nUUDD\n", "4"),
			hid2("The stretch need not start at the beginning", "3\nUUD\n", "2"),
			hid2("A long imbalance around a short balance", "5\nUUUDU\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char days[200005];
    if (scanf("%200000s", days) != 1) return 1;
    static int first[400005];
    for (int i = 0; i < 400005; i++) first[i] = -2;
    int base = 200002;
    first[base] = -1;
    int score = 0, best = 0;
    for (int i = 0; i < n; i++) {
        score += days[i] == 'U' ? 1 : -1;
        int at = base + score;
        if (first[at] == -2) first[at] = i;
        else if (i - first[at] > best) best = i - first[at];
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
	var days string
	fmt.Fscan(in, &n, &days)
	first := map[int]int{0: -1}
	score, best := 0, 0
	for i := 0; i < n; i++ {
		if days[i] == 'U' {
			score++
		} else {
			score--
		}
		if at, ok := first[score]; ok {
			if i-at > best {
				best = i - at
			}
		} else {
			first[score] = i
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, days = int(data[0]), data[1]
    first = {0: -1}
    score = best = 0
    for i in range(n):
        score += 1 if days[i] == "U" else -1
        if score in first:
            best = max(best, i - first[score])
        else:
            first[score] = i
    print(best)

main()
`},
	},
	{
		slug: "average-alarm", title: "Average alarm", difficulty: "easy", topic: "Sliding window",
		statement: `A monitor raises an alarm whenever the average of the last ~k~ readings is
**strictly above** ~limit~.

**Input.** ~n k limit~ on the first line, then ~n~ readings.

**Output.** How many windows of exactly ~k~ readings raise an alarm. If the log
is shorter than ~k~, no window exists, so no alarms.

    Input        Output
    3 2 2        2
    1 5 1

Both windows average 3. Work in whole numbers: comparing the sum against
~limit * k~ says the same thing as comparing the average against ~limit~, with
no rounding to argue about.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, k, limit;
    if (scanf("%d %d %d", &n, &k, &limit) != 3) return 1;
    /* slide a window of k and count the sums above limit * k */
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
	var n, k, limit int
	fmt.Fscan(in, &n, &k, &limit)
	// slide a window of k and count the sums above limit * k
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n, k, limit = int(data[0]), int(data[1]), int(data[2])
    v = [int(x) for x in data[3:3 + n]]
    # slide a window of k and count the sums above limit * k
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Both windows trip it", "3 2 2\n1 5 1\n", "2"),
			tokens("A quiet log raises nothing", "3 2 5\n1 1 1\n", "0"),
			tokens("Too short for a window", "1 2 0\n9\n", "0"),
			hid2("Exactly at the limit is not above it", "2 2 2\n2 2\n", "0"),
			hid2("Only some windows trip it", "4 2 2\n9 0 0 0\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, k, limit;
    if (scanf("%d %d %d", &n, &k, &limit) != 3) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    for (int i = 0; i < n; i++) if (scanf("%lld", &v[i]) != 1) return 1;
    if (k <= 0 || k > n) { printf("0\n"); return 0; }
    long long total = 0, bar = (long long)limit * k;
    for (int i = 0; i < k; i++) total += v[i];
    int count = total > bar ? 1 : 0;
    for (int i = k; i < n; i++) {
        total += v[i] - v[i - k];
        if (total > bar) count++;
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
	var n, k, limit int
	fmt.Fscan(in, &n, &k, &limit)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	if k <= 0 || k > n {
		fmt.Println(0)
		return
	}
	total, bar := 0, limit*k
	for i := 0; i < k; i++ {
		total += v[i]
	}
	count := 0
	if total > bar {
		count++
	}
	for i := k; i < n; i++ {
		total += v[i] - v[i-k]
		if total > bar {
			count++
		}
	}
	fmt.Println(count)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, k, limit = int(data[0]), int(data[1]), int(data[2])
    v = [int(x) for x in data[3:3 + n]]
    if k <= 0 or k > n:
        print(0)
        return
    total, bar = sum(v[:k]), limit * k
    count = 1 if total > bar else 0
    for i in range(k, n):
        total += v[i] - v[i - k]
        if total > bar:
            count += 1
    print(count)

main()
`},
	},
	{
		slug: "affordable-pairs", title: "Affordable pairs", difficulty: "medium", topic: "Two pointers",
		statement: `A price list is already sorted, cheapest first.

**Input.** ~n budget~ on the first line, then ~n~ prices in ascending order.

**Output.** How many pairs of two **different** items cost ~budget~ or less
together.

    Input        Output
    4 5          4
    1 2 3 4

Those are 1+2, 1+3, 1+4 and 2+3. Checking every pair is quadratic — but the list
is sorted, and that is the whole gift: if the cheapest and the dearest fit
together, so does the cheapest with everything in between.`,
		timeLimitMs: 4000,
		starters:    triCountWith("budget", "close in from both ends of the sorted list"),
		tests: []store.Check{
			tokens("Counts the affordable pairs", "4 5\n1 2 3 4\n", "4"),
			tokens("Nothing is affordable", "2 5\n9 9\n", "0"),
			tokens("One item makes no pair", "1 100\n1\n", "0"),
			tokens("Everything is affordable", "3 5\n1 1 1\n", "3"),
			hid2("Exactly on budget counts", "2 5\n2 3\n", "1"),
			hid2("An empty list", "0 5\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    long long budget;
    if (scanf("%d %lld", &n, &budget) != 2) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    for (int i = 0; i < n; i++) if (scanf("%lld", &v[i]) != 1) return 1;
    long long count = 0;
    int i = 0, j = n - 1;
    while (i < j) {
        if (v[i] + v[j] <= budget) { count += j - i; i++; }
        else j--;
    }
    printf("%lld\n", count);
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
	var n, budget int
	fmt.Fscan(in, &n, &budget)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	count, i, j := 0, 0, n-1
	for i < j {
		if v[i]+v[j] <= budget {
			count += j - i
			i++
		} else {
			j--
		}
	}
	fmt.Println(count)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, budget = int(data[0]), int(data[1])
    v = [int(x) for x in data[2:2 + n]]
    count, i, j = 0, 0, n - 1
    while i < j:
        if v[i] + v[j] <= budget:
            count += j - i
            i += 1
        else:
            j -= 1
    print(count)

main()
`},
	},
	{
		slug: "range-totals", title: "Range totals", difficulty: "medium", topic: "Prefix sums",
		statement: `**Input.** ~n q~ on the first line, then ~n~ values, then ~q~ lines of
~lo hi~ asking for the sum of ~values[lo]~ through ~values[hi]~ **inclusive**.

**Output.** One answer per query, space-separated, in order.

    Input        Output
    3 2          6 2
    1 2 3
    0 2
    1 1

Answering each query by adding up its range does the same additions again and
again; answering all of them from one pass over the values does not.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, q;
    if (scanf("%d %d", &n, &q) != 2) return 1;
    /* build running totals, then answer each query in one step */
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
	var n, q int
	fmt.Fscan(in, &n, &q)
	// build running totals, then answer each query in one step
	fmt.Println("")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, q = data[0], data[1]
    values = data[2:2 + n]
    # build running totals, then answer each query in one step
    print("")

main()
`),
		tests: []store.Check{
			tokens("Answers each query", "3 2\n1 2 3\n0 2\n1 1\n", "6 2"),
			tokens("No queries", "2 0\n1 2\n", ""),
			tokens("A single-element range", "1 1\n5\n0 0\n", "5"),
			hid2("Overlapping ranges", "4 3\n1 1 1 1\n0 1\n1 3\n0 3\n", "2 3 4"),
			hid2("Negative values sum too", "2 1\n-1 5\n0 1\n", "4"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, q;
    if (scanf("%d %d", &n, &q) != 2) return 1;
    long long *pre = malloc((size_t)(n + 1) * sizeof *pre);
    pre[0] = 0;
    for (int i = 0; i < n; i++) {
        long long v;
        if (scanf("%lld", &v) != 1) return 1;
        pre[i + 1] = pre[i] + v;
    }
    for (int i = 0; i < q; i++) {
        int lo, hi;
        if (scanf("%d %d", &lo, &hi) != 2) return 1;
        printf(i ? " %lld" : "%lld", pre[hi + 1] - pre[lo]);
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
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &q)
	pre := make([]int, n+1)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		pre[i+1] = pre[i] + v
	}
	out := make([]string, q)
	for i := 0; i < q; i++ {
		var lo, hi int
		fmt.Fscan(in, &lo, &hi)
		out[i] = strconv.Itoa(pre[hi+1] - pre[lo])
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, q = data[0], data[1]
    values = data[2:2 + n]
    pre = [0]
    for v in values:
        pre.append(pre[-1] + v)
    base = 2 + n
    out = []
    for i in range(q):
        lo, hi = data[base + 2 * i], data[base + 1 + 2 * i]
        out.append(str(pre[hi + 1] - pre[lo]))
    print(" ".join(out))

main()
`},
	},
	{
		slug: "worst-drawdown", title: "Worst drawdown", difficulty: "medium", topic: "Prefix sums",
		statement: `A drawdown is how far a price has fallen from the highest point it had
reached **before** it.

**Input.** A count ~n~, then ~n~ prices.

**Output.** The worst drawdown. A price that only ever rises has a drawdown of 0,
and so does an empty list.

    Input        Output
    4            4
    5 3 6 2

The 2 comes after a peak of 6.`,
		timeLimitMs: 4000,
		starters:    triCount("track the running peak and the worst fall from it"),
		tests: []store.Check{
			tokens("Measures from the peak before it", "4\n5 3 6 2\n", "4"),
			tokens("A rising market never falls", "3\n1 2 3\n", "0"),
			tokens("No prices at all", "0\n", "0"),
			tokens("A single price", "1\n7\n", "0"),
			hid2("A later peak does not count backwards", "3\n9 1 100\n", "8"),
			hid2("The worst fall is the first one", "4\n10 1 5 4\n", "9"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long peak = 0, worst = 0;
    int started = 0;
    for (int i = 0; i < n; i++) {
        long long p;
        if (scanf("%lld", &p) != 1) return 1;
        if (!started || p > peak) { peak = p; started = 1; }
        if (peak - p > worst) worst = peak - p;
    }
    printf("%lld\n", worst);
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
	peak, worst, started := 0, 0, false
	for i := 0; i < n; i++ {
		var p int
		fmt.Fscan(in, &p)
		if !started || p > peak {
			peak, started = p, true
		}
		if peak-p > worst {
			worst = peak - p
		}
	}
	fmt.Println(worst)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    peak, worst = None, 0
    for x in data[1:1 + n]:
        p = int(x)
        if peak is None or p > peak:
            peak = p
        worst = max(worst, peak - p)
    print(worst)

main()
`},
	},
	{
		slug: "pivot-points", title: "Pivot points", difficulty: "medium", topic: "Prefix sums",
		statement: `A pivot is a position where everything to its left weighs exactly as much
as everything to its right. The value at the pivot itself belongs to neither
side.

**Input.** A count ~n~, then ~n~ weights.

**Output.** Every pivot position, space-separated, in order. An end position is a
pivot when the other side sums to zero — including when that side is empty.

    Input        Output
    4            2
    1 2 3 3

Re-adding each side at every position is quadratic; one running total is not.`,
		timeLimitMs: 4000,
		starters:    triCount("carry the left total and derive the right from the whole"),
		tests: []store.Check{
			tokens("Finds the balance point", "4\n1 2 3 3\n", "2"),
			tokens("A single item balances nothing against nothing", "1\n9\n", "0"),
			tokens("No items, no pivots", "0\n", ""),
			hid2("There can be more than one", "3\n0 0 0\n", "0 1 2"),
			hid2("There can be none", "2\n1 9\n", ""),
			hid2("Negative weights can balance", "5\n1 -1 0 3 -3\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long *w = malloc((size_t)(n ? n : 1) * sizeof *w);
    long long total = 0;
    for (int i = 0; i < n; i++) { if (scanf("%lld", &w[i]) != 1) return 1; total += w[i]; }
    long long left = 0;
    int first = 1;
    for (int i = 0; i < n; i++) {
        if (left == total - left - w[i]) { printf(first ? "%d" : " %d", i); first = 0; }
        left += w[i];
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
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	w := make([]int, n)
	total := 0
	for i := range w {
		fmt.Fscan(in, &w[i])
		total += w[i]
	}
	var out []string
	left := 0
	for i, x := range w {
		if left == total-left-x {
			out = append(out, strconv.Itoa(i))
		}
		left += x
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    w = [int(x) for x in data[1:1 + n]]
    total = sum(w)
    left, out = 0, []
    for i, x in enumerate(w):
        if left == total - left - x:
            out.append(str(i))
        left += x
    print(" ".join(out))

main()
`},
	},
	{
		slug: "log-throttle", title: "Log throttle", difficulty: "medium", topic: "Sliding window",
		statement: `A logger accepts at most ~cap~ messages in any ~window~ of time. Messages
arrive in time order. When one arrives and ~cap~ messages have already been
accepted within the last ~window~ ticks, it is dropped — and a dropped message
does not count against the limit either, because it was never logged.

A message at time ~t~ counts against one at time ~u~ only while ~u > t - window~.

**Input.** ~n window cap~ on the first line, then ~n~ arrival times.

**Output.** How many messages are lost.

    Input          Output
    3 10 2         1
    1 2 3`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n, window, cap;
    if (scanf("%d %d %d", &n, &window, &cap) != 3) return 1;
    /* keep the accepted times still inside the window */
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
	var n, window, cap int
	fmt.Fscan(in, &n, &window, &cap)
	// keep the accepted times still inside the window
	fmt.Println(0)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, window, cap = data[0], data[1], data[2]
    times = data[3:3 + n]
    # keep the accepted times still inside the window
    print(0)

main()
`),
		tests: []store.Check{
			tokens("The third message is over the cap", "3 10 2\n1 2 3\n", "1"),
			tokens("Spread out enough to all get through", "3 10 1\n1 20 40\n", "0"),
			tokens("Nothing logged", "0 10 1\n", "0"),
			hid2("A message ageing out frees a slot", "3 10 2\n1 2 11\n", "0"),
			hid2("A dropped message does not hold a slot", "4 10 1\n1 1 1 30\n", "2"),
			hid2("A cap of nothing drops everything", "2 10 0\n1 2\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, window, cap;
    if (scanf("%d %d %d", &n, &window, &cap) != 3) return 1;
    long long *kept = malloc((size_t)(n ? n : 1) * sizeof *kept);
    int head = 0, tail = 0, lost = 0;
    for (int i = 0; i < n; i++) {
        long long t;
        if (scanf("%lld", &t) != 1) return 1;
        while (head < tail && kept[head] <= t - window) head++;
        if (tail - head < cap) kept[tail++] = t; else lost++;
    }
    printf("%d\n", lost);
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
	var n, window, capacity int
	fmt.Fscan(in, &n, &window, &capacity)
	var kept []int
	lost := 0
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		for len(kept) > 0 && kept[0] <= t-window {
			kept = kept[1:]
		}
		if len(kept) < capacity {
			kept = append(kept, t)
		} else {
			lost++
		}
	}
	fmt.Println(lost)
}
`,
			"python": `import sys
from collections import deque

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, window, cap = data[0], data[1], data[2]
    kept, lost = deque(), 0
    for t in data[3:3 + n]:
        while kept and kept[0] <= t - window:
            kept.popleft()
        if len(kept) < cap:
            kept.append(t)
        else:
            lost += 1
    print(lost)

main()
`},
	},
	{
		slug: "cooldown-check", title: "Cooldown check", difficulty: "easy", topic: "Dictionaries",
		statement: `A game gives each ability a cooldown: once used, it cannot be used again
for ~cooldown~ turns. Measure the gap as the distance between the two turns —
turn 0 and turn 3 are 3 apart, so a cooldown of 3 permits that and a cooldown of
4 does not.

**Input.** ~n cooldown~ on the first line, then ~n~ ability names, one per turn.

**Output.** ~yes~ if the whole sequence is legal, ~no~ otherwise.

    Input        Output
    3 2          yes
    a b a`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n, cooldown;
    if (scanf("%d %d", &n, &cooldown) != 2) return 1;
    /* remember the turn each ability was last used */
    printf("yes\n");
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
	var n, cooldown int
	fmt.Fscan(in, &n, &cooldown)
	// remember the turn each ability was last used
	fmt.Println("yes")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n, cooldown = int(data[0]), int(data[1])
    actions = data[2:2 + n]
    # remember the turn each ability was last used
    print("yes")

main()
`),
		tests: []store.Check{
			tokens("Far enough apart", "3 2\na b a\n", "yes"),
			tokens("Used again too soon", "2 2\na a\n", "no"),
			tokens("A gap exactly equal to the cooldown is legal", "2 1\na a\n", "yes"),
			tokens("Nothing to check", "0 5\n", "yes"),
			hid2("Exactly at the cooldown is allowed", "4 3\na b c a\n", "yes"),
			hid2("One turn short is not", "4 4\na b c a\n", "no"),
			hid2("A cooldown of zero permits anything", "2 0\na a\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n, cooldown;
    if (scanf("%d %d", &n, &cooldown) != 2) return 1;
    static char name[100000][32];
    static int last[100000];
    int kinds = 0, ok = 1;
    for (int i = 0; i < n; i++) {
        char w[32];
        if (scanf("%31s", w) != 1) return 1;
        int at = -1;
        for (int j = 0; j < kinds; j++) if (strcmp(name[j], w) == 0) at = j;
        if (at < 0) { at = kinds++; strcpy(name[at], w); }
        else if (i - last[at] < cooldown) ok = 0;
        last[at] = i;
    }
    printf("%s\n", ok ? "yes" : "no");
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
	var n, cooldown int
	fmt.Fscan(in, &n, &cooldown)
	last := map[string]int{}
	ok := true
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		if at, seen := last[w]; seen && i-at < cooldown {
			ok = false
		}
		last[w] = i
	}
	if ok {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n, cooldown = int(data[0]), int(data[1])
    last, ok = {}, True
    for i, a in enumerate(data[2:2 + n]):
        if a in last and i - last[a] < cooldown:
            ok = False
        last[a] = i
    print("yes" if ok else "no")

main()
`},
	},
	{
		slug: "neighbour-smoothing", title: "Neighbour smoothing", difficulty: "easy", topic: "Sliding window",
		statement: `Smooth a noisy signal: replace each reading with the whole-number average
of itself and its immediate neighbours, rounding **down**.

Readings at the ends have only one neighbour, so they average over two values
rather than three. Every average is taken from the **original** readings, not
from the smoothed ones you are producing.

**Input.** A count ~n~, then ~n~ readings.

**Output.** The smoothed signal, space-separated.

    Input      Output
    3          1 2 2
    1 2 3`,
		timeLimitMs: 4000,
		starters:    triCount("average each reading with the neighbours it has"),
		tests: []store.Check{
			tokens("Smooths a short signal", "3\n1 2 3\n", "1 2 2"),
			tokens("A single reading is its own average", "1\n5\n", "5"),
			tokens("Nothing to smooth", "0\n", ""),
			tokens("A flat signal stays flat", "3\n4 4 4\n", "4 4 4"),
			hid2("Averages come from the original readings", "5\n0 9 0 9 0\n", "4 3 6 3 4"),
			hid2("Rounding is always down", "2\n0 1\n", "0 0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    for (int i = 0; i < n; i++) if (scanf("%lld", &v[i]) != 1) return 1;
    for (int i = 0; i < n; i++) {
        long long sum = v[i];
        int k = 1;
        if (i > 0) { sum += v[i - 1]; k++; }
        if (i + 1 < n) { sum += v[i + 1]; k++; }
        long long q = sum / k;
        if (sum % k != 0 && ((sum < 0) != (k < 0))) q--;
        printf(i ? " %lld" : "%lld", q);
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
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	out := make([]string, n)
	for i := range v {
		sum, k := v[i], 1
		if i > 0 {
			sum += v[i-1]
			k++
		}
		if i+1 < n {
			sum += v[i+1]
			k++
		}
		q := sum / k
		if sum%k != 0 && (sum < 0) != (k < 0) {
			q--
		}
		out[i] = strconv.Itoa(q)
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    v = [int(x) for x in data[1:1 + n]]
    out = []
    for i in range(n):
        window = v[max(0, i - 1):min(n, i + 2)]
        out.append(str(sum(window) // len(window)))
    print(" ".join(out))

main()
`},
	},
	{
		slug: "peak-hours", title: "Peak hours", difficulty: "easy", topic: "Sliding window",
		statement: `**Input.** A count ~n~, then ~n~ counts.

**Output.** The positions that are strictly higher than every neighbour they
have, space-separated. A position at either end has only one neighbour and only
has to beat that one. A lone reading has no neighbours at all, so it is a peak.

    Input      Output
    3          1
    1 3 2`,
		timeLimitMs: 4000,
		starters:    triCount("report every position that beats the neighbours it has"),
		tests: []store.Check{
			tokens("Finds the peak", "3\n1 3 2\n", "1"),
			tokens("A lone reading is a peak", "1\n5\n", "0"),
			tokens("Nothing at all", "0\n", ""),
			tokens("A peak at the end", "2\n1 2\n", "1"),
			hid2("A plateau is not a peak", "4\n1 2 2 1\n", ""),
			hid2("Two separate peaks", "4\n3 1 3 1\n", "0 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    for (int i = 0; i < n; i++) if (scanf("%lld", &v[i]) != 1) return 1;
    int first = 1;
    for (int i = 0; i < n; i++) {
        int left = (i == 0) || v[i] > v[i - 1];
        int right = (i == n - 1) || v[i] > v[i + 1];
        if (left && right) { printf(first ? "%d" : " %d", i); first = 0; }
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
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	var out []string
	for i := range v {
		if (i == 0 || v[i] > v[i-1]) && (i == n-1 || v[i] > v[i+1]) {
			out = append(out, strconv.Itoa(i))
		}
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    v = [int(x) for x in data[1:1 + n]]
    out = [str(i) for i in range(n)
           if (i == 0 or v[i] > v[i - 1]) and (i == n - 1 or v[i] > v[i + 1])]
    print(" ".join(out))

main()
`},
	},
}

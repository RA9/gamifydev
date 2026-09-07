package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// orderProblems: sorting, ranking, and the tie-breaks and derived keys that make
// an ordering mean something.
var orderProblems = []seedProblem{
	{
		slug: "competition-ranks", title: "Competition ranks", difficulty: "medium", topic: "Sorting",
		statement: `A race is scored the way races are: the highest score is 1st, everyone
tied shares a rank, and the next rank after a tie **skips**. Two people in 2nd
place means nobody is 3rd.

**Input.** A count ~n~, then ~n~ scores.

**Output.** Each competitor's rank, space-separated, in the order the scores were
given.

    Input          Output
    4              1 2 2 4
    10 8 8 5

Counting how many people beat each runner, one runner at a time, is quadratic:
fine for the examples here, hopeless for a field of two hundred thousand. Sort
once and every rank falls out of the same pass.`,
		timeLimitMs: 4000,
		starters:    triCount("rank the scores without comparing every pair"),
		tests: []store.Check{
			tokens("A tie shares a rank and skips the next", "4\n10 8 8 5\n", "1 2 2 4"),
			tokens("Everyone distinct", "3\n3 1 2\n", "1 3 2"),
			tokens("Nobody raced", "0\n", ""),
			tokens("Everyone tied for first", "3\n5 5 5\n", "1 1 1"),
			hid2("A tie at the bottom", "3\n9 1 1\n", "1 2 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static int desc(const void *a, const void *b) {
    long long x = *(const long long *)a, y = *(const long long *)b;
    return (x < y) - (x > y);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    long long *v = malloc((size_t)(n ? n : 1) * sizeof *v);
    long long *s = malloc((size_t)(n ? n : 1) * sizeof *s);
    for (int i = 0; i < n; i++) { if (scanf("%lld", &v[i]) != 1) return 1; s[i] = v[i]; }
    qsort(s, n, sizeof *s, desc);
    for (int i = 0; i < n; i++) {
        int lo = 0, hi = n - 1, at = n;
        while (lo <= hi) {
            int mid = (lo + hi) / 2;
            if (s[mid] <= v[i]) { at = mid; hi = mid - 1; } else lo = mid + 1;
        }
        printf(i ? " %d" : "%d", at + 1);
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
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	v := make([]int, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	count := map[int]int{}
	for _, x := range v {
		count[x]++
	}
	distinct := make([]int, 0, len(count))
	for k := range count {
		distinct = append(distinct, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(distinct)))
	rank := map[int]int{}
	place := 1
	for _, k := range distinct {
		rank[k] = place
		place += count[k]
	}
	for i, x := range v {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, rank[x])
	}
	fmt.Fprintln(out)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    scores = [int(x) for x in data[1:1 + n]]
    counts = {}
    for s in scores:
        counts[s] = counts.get(s, 0) + 1
    rank, place = {}, 1
    for s in sorted(counts, reverse=True):
        rank[s] = place
        place += counts[s]
    print(" ".join(str(rank[s]) for s in scores))

main()
`},
	},
	{
		slug: "version-order", title: "Version order", difficulty: "medium", topic: "Sorting",
		statement: `Version strings are dotted numbers, and they do not sort like text: ~1.10~
comes **after** ~1.2~, because ten is more than two.

**Input.** A count ~n~, then ~n~ version strings.

**Output.** The versions in ascending order, space-separated. Compare segment by
segment as numbers; a version with fewer segments is padded with zeros, so ~1.2~
and ~1.2.0~ are the same version and keep the order they were given in.

    Input               Output
    3                   1.2 1.2.1 1.10
    1.10 1.2 1.2.1`,
		timeLimitMs: 4000,
		starters:    triWords("sort by the numeric value of each segment"),
		tests: []store.Check{
			tokens("Ten sorts after two", "3\n1.10 1.2 1.2.1\n", "1.2 1.2.1 1.10"),
			tokens("Nothing to sort", "0\n", ""),
			tokens("Leading zeros are just numbers", "2\n1.09 1.9\n", "1.09 1.9"),
			hid2("A short version pads with zeros", "2\n1.2.1 1.2\n", "1.2 1.2.1"),
			hid2("Equal versions keep their given order", "2\n1.2.0 1.2\n", "1.2.0 1.2"),
			hid2("Major version wins over everything", "2\n2.0 1.99.99\n", "1.99.99 2.0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define SEG 8
static char raw[100000][64];
static long long key[100000][SEG];
static int order[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    for (int i = 0; i < SEG; i++) {
        if (key[x][i] != key[y][i]) return key[x][i] < key[y][i] ? -1 : 1;
    }
    return (x > y) - (x < y);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%63s", raw[i]) != 1) return 1;
        order[i] = i;
        for (int j = 0; j < SEG; j++) key[i][j] = 0;
        int seg = 0;
        long long acc = 0;
        for (int p = 0; ; p++) {
            char c = raw[i][p];
            if (c == '.' || c == '\0') {
                if (seg < SEG) key[i][seg++] = acc;
                acc = 0;
                if (c == '\0') break;
            } else {
                acc = acc * 10 + (c - '0');
            }
        }
    }
    qsort(order, n, sizeof *order, cmp);
    for (int i = 0; i < n; i++) printf(i ? " %s" : "%s", raw[order[i]]);
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
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	v := make([]string, n)
	for i := range v {
		fmt.Fscan(in, &v[i])
	}
	key := func(s string) [8]int {
		var k [8]int
		for i, part := range strings.Split(s, ".") {
			if i < 8 {
				k[i], _ = strconv.Atoi(part)
			}
		}
		return k
	}
	sort.SliceStable(v, func(i, j int) bool {
		a, b := key(v[i]), key(v[j])
		for x := 0; x < 8; x++ {
			if a[x] != b[x] {
				return a[x] < b[x]
			}
		}
		return false
	})
	fmt.Println(strings.Join(v, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    versions = data[1:1 + n]
    def key(v):
        parts = [int(p) for p in v.split(".")]
        return tuple(parts + [0] * (8 - len(parts)))
    print(" ".join(sorted(versions, key=key)))

main()
`},
	},
	{
		slug: "ticket-order", title: "Ticket order", difficulty: "medium", topic: "Sorting",
		statement: `A support queue holds tickets, severity 1 being the most urgent.

A ticket that has waited **more than 48 hours** is escalated: treat it as one
severity more urgent, though nothing goes above 1.

**Input.** A count ~n~, then ~n~ lines of ~id severity hours~.

**Output.** The ids, space-separated, ordered by effective severity, then by the
longest wait, then by id.

    Input          Output
    3              c b a
    a 2 10
    b 1 1
    c 2 50

Ticket ~c~ escalates to severity 1 and has waited longest, so it goes first.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char id[100000][32];
    static int sev[100000], hours[100000], order[100000];
    for (int i = 0; i < n; i++) {
        if (scanf("%31s %d %d", id[i], &sev[i], &hours[i]) != 3) return 1;
        order[i] = i;
    }
    /* sort by effective severity, then wait, then id */
    for (int i = 0; i < n; i++) printf(i ? " %s" : "%s", id[order[i]]);
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

type ticket struct {
	id          string
	sev, hours  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	ts := make([]ticket, n)
	for i := range ts {
		fmt.Fscan(in, &ts[i].id, &ts[i].sev, &ts[i].hours)
	}
	_ = sort.Slice
	// sort by effective severity, then wait, then id
	out := make([]string, n)
	for i, t := range ts {
		out[i] = t.id
	}
	fmt.Println(strings.Join(out, " "))
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    tickets = [(data[1 + 3 * i], int(data[2 + 3 * i]), int(data[3 + 3 * i])) for i in range(n)]
    # sort by effective severity, then wait, then id
    print(" ".join(t[0] for t in tickets))

main()
`),
		tests: []store.Check{
			tokens("Escalation reorders the queue", "3\na 2 10\nb 1 1\nc 2 50\n", "c b a"),
			tokens("Without escalation, severity decides", "2\na 2 1\nb 1 1\n", "b a"),
			tokens("An empty queue", "0\n", ""),
			hid2("Exactly 48 hours does not escalate", "2\na 1 1\nb 2 48\n", "a b"),
			hid2("Severity 1 cannot escalate further", "2\na 1 100\nb 1 99\n", "a b"),
			hid2("Ids break a complete tie", "2\nz 1 5\na 1 5\n", "a z"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char id[100000][32];
static int sev[100000], hours[100000];

static int eff(int i) {
    int e = hours[i] > 48 ? sev[i] - 1 : sev[i];
    return e < 1 ? 1 : e;
}

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (eff(x) != eff(y)) return eff(x) - eff(y);
    if (hours[x] != hours[y]) return hours[y] - hours[x];
    return strcmp(id[x], id[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static int order[100000];
    for (int i = 0; i < n; i++) {
        if (scanf("%31s %d %d", id[i], &sev[i], &hours[i]) != 3) return 1;
        order[i] = i;
    }
    qsort(order, n, sizeof *order, cmp);
    for (int i = 0; i < n; i++) printf(i ? " %s" : "%s", id[order[i]]);
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

type ticket struct {
	id         string
	sev, hours int
}

func (t ticket) eff() int {
	if t.hours > 48 && t.sev > 1 {
		return t.sev - 1
	}
	return t.sev
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	ts := make([]ticket, n)
	for i := range ts {
		fmt.Fscan(in, &ts[i].id, &ts[i].sev, &ts[i].hours)
	}
	sort.Slice(ts, func(i, j int) bool {
		a, b := ts[i], ts[j]
		if a.eff() != b.eff() {
			return a.eff() < b.eff()
		}
		if a.hours != b.hours {
			return a.hours > b.hours
		}
		return a.id < b.id
	})
	out := make([]string, n)
	for i, t := range ts {
		out[i] = t.id
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    tickets = [(data[1 + 3 * i], int(data[2 + 3 * i]), int(data[3 + 3 * i])) for i in range(n)]
    def key(t):
        tid, sev, hours = t
        eff = max(1, sev - 1) if hours > 48 else sev
        return (eff, -hours, tid)
    print(" ".join(t[0] for t in sorted(tickets, key=key)))

main()
`},
	},
	{
		slug: "bracket-seeding", title: "Bracket seeding", difficulty: "hard", topic: "Sorting",
		statement: `A knockout bracket is seeded so the two best players can only meet in the
final, the top four only in the semis, and so on.

You build it by doubling: start with seed 1 alone. To go from a bracket of
~size~ to one of ~2 * size~, replace every seed ~s~ with the pair
~s, 2 * size + 1 - s~, keeping them in place.

    1                    (size 1)
    1 2                  (size 2)
    1 4 2 3              (size 4)
    1 8 4 5 2 7 3 6      (size 8)

**Input.** A single power of two, ~n~.

**Output.** The first-round matches, one ~a b~ per line, read left to right off
that final row. A bracket of one player has no matches and prints nothing.

    Input   Output
    4       1 4
            2 3`,
		timeLimitMs: 4000,
		starters:    triScalars([]string{"n"}, "double the bracket until it holds n, then read off the pairs"),
		tests: []store.Check{
			tokens("A bracket of four", "4\n", "1 4 2 3"),
			tokens("A bracket of two is one match", "2\n", "1 2"),
			tokens("One player plays nobody", "1\n", ""),
			tokens("A bracket of eight", "8\n", "1 8 4 5 2 7 3 6"),
			hid2("The top seed meets the bottom seed", "16\n", "1 16 8 9 4 13 5 12 2 15 7 10 3 14 6 11"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    long long n;
    if (scanf("%lld", &n) != 1) return 1;
    static long long cur[1 << 16], nxt[1 << 16];
    long long len = 1, size = 1;
    cur[0] = 1;
    while (size < n) {
        size *= 2;
        long long k = 0;
        for (long long i = 0; i < len; i++) {
            nxt[k++] = cur[i];
            nxt[k++] = size + 1 - cur[i];
        }
        for (long long i = 0; i < k; i++) cur[i] = nxt[i];
        len = k;
    }
    for (long long i = 0; i + 1 < len; i += 2) printf("%lld %lld\n", cur[i], cur[i + 1]);
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
	order := []int{1}
	for size := 1; size < n; {
		size *= 2
		next := make([]int, 0, len(order)*2)
		for _, s := range order {
			next = append(next, s, size+1-s)
		}
		order = next
	}
	for i := 0; i+1 < len(order); i += 2 {
		fmt.Fprintln(out, order[i], order[i+1])
	}
}
`,
			"python": `import sys

def main():
    n = int(sys.stdin.read().split()[0])
    order, size = [1], 1
    while size < n:
        size *= 2
        nxt = []
        for s in order:
            nxt.append(s)
            nxt.append(size + 1 - s)
        order = nxt
    for i in range(0, len(order) - 1, 2):
        print(order[i], order[i + 1])

main()
`},
	},
	{
		slug: "merge-streams", title: "Merge streams", difficulty: "medium", topic: "Sorting",
		statement: `Several machines each write a log, and each log is already in time order.

**Input.** A count ~k~, then ~k~ lines. Each line is a count ~m~ followed by
~m~ timestamps.

**Output.** Every entry in one ordered list, space-separated. Entries with the
same timestamp are ordered by which stream they came from — the earlier stream
first — and two entries from the same stream keep their order.

    Input      Output
    2          1 2 3 3
    2 1 3
    2 2 3`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int k;
    if (scanf("%d", &k) != 1) return 1;
    /* read every stream, then emit all entries in order */
    printf("\n");
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
	var k int
	fmt.Fscan(in, &k)
	_ = sort.SliceStable
	// read every stream, then emit all entries in order
	fmt.Println("")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    k = data[0]
    # read every stream, then emit all entries in order
    print("")

main()
`),
		tests: []store.Check{
			tokens("Interleaves two streams", "2\n2 1 3\n2 2 3\n", "1 2 3 3"),
			tokens("One stream comes back as it was", "1\n2 5 6\n", "5 6"),
			tokens("No streams at all", "0\n", ""),
			tokens("Empty streams contribute nothing", "3\n0\n1 1\n0\n", "1"),
			hid2("A tie goes to the earlier stream", "3\n1 2\n1 1\n1 1\n", "1 1 2"),
			hid2("Three streams merge cleanly", "3\n2 1 4\n2 2 5\n2 3 6\n", "1 2 3 4 5 6"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long val[200000];
static int stream[200000], seq[200000], idx[200000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (val[x] != val[y]) return val[x] < val[y] ? -1 : 1;
    if (stream[x] != stream[y]) return stream[x] - stream[y];
    return seq[x] - seq[y];
}

int main(void) {
    int k;
    if (scanf("%d", &k) != 1) return 1;
    int total = 0;
    for (int s = 0; s < k; s++) {
        int m;
        if (scanf("%d", &m) != 1) return 1;
        for (int j = 0; j < m; j++) {
            if (scanf("%lld", &val[total]) != 1) return 1;
            stream[total] = s;
            seq[total] = j;
            idx[total] = total;
            total++;
        }
    }
    qsort(idx, total, sizeof *idx, cmp);
    for (int i = 0; i < total; i++) printf(i ? " %lld" : "%lld", val[idx[i]]);
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
	"strconv"
	"strings"
)

type entry struct {
	val, stream int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	fmt.Fscan(in, &k)
	var all []entry
	for s := 0; s < k; s++ {
		var m int
		fmt.Fscan(in, &m)
		for j := 0; j < m; j++ {
			var v int
			fmt.Fscan(in, &v)
			all = append(all, entry{v, s})
		}
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].val != all[j].val {
			return all[i].val < all[j].val
		}
		return all[i].stream < all[j].stream
	})
	out := make([]string, len(all))
	for i, e := range all {
		out[i] = strconv.Itoa(e.val)
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    k, pos = data[0], 1
    tagged = []
    for s in range(k):
        m = data[pos]; pos += 1
        for v in data[pos:pos + m]:
            tagged.append((v, s))
        pos += m
    tagged.sort(key=lambda p: (p[0], p[1]))
    print(" ".join(str(v) for v, _ in tagged))

main()
`},
	},
	{
		slug: "fair-split", title: "Fair split", difficulty: "hard", topic: "Dynamic programming",
		statement: `Two removal vans have to share a load.

**Input.** A count ~n~ (at most 30), then ~n~ whole-number weights, none
negative and none above 1000.

**Output.** The smallest possible difference between the two vans' totals.

    Input        Output
    3            0
    1 2 3

Every item goes in one van or the other; a van may end up empty. Thirty items is
a realistic load, and trying all 2^30 ways to divide them is not. What you can
afford to track is every total one van could possibly end up with.`,
		timeLimitMs: 4000,
		starters:    triCount("track which totals are reachable, then pick the closest to half"),
		tests: []store.Check{
			tokens("An even split", "3\n1 2 3\n", "0"),
			tokens("A load that cannot be evened out", "2\n8 1\n", "7"),
			tokens("Nothing to move", "0\n", "0"),
			tokens("One item goes in one van", "1\n5\n", "5"),
			tokens("Fast enough for thirty items", "30\n1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1\n", "0"),
			hid2("A near miss", "3\n1 2 4\n", "1"),
			hid2("Zero-weight items change nothing", "3\n0 0 5\n", "5"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static int w[40];
    int total = 0;
    for (int i = 0; i < n; i++) { if (scanf("%d", &w[i]) != 1) return 1; total += w[i]; }
    static char reach[30001];
    memset(reach, 0, sizeof reach);
    reach[0] = 1;
    for (int i = 0; i < n; i++)
        for (int s = total; s >= w[i]; s--)
            if (reach[s - w[i]]) reach[s] = 1;
    int best = total;
    for (int s = 0; s <= total; s++) {
        if (!reach[s]) continue;
        int d = total - 2 * s;
        if (d < 0) d = -d;
        if (d < best) best = d;
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
	w := make([]int, n)
	total := 0
	for i := range w {
		fmt.Fscan(in, &w[i])
		total += w[i]
	}
	reach := make([]bool, total+1)
	reach[0] = true
	for _, x := range w {
		for s := total; s >= x; s-- {
			if reach[s-x] {
				reach[s] = true
			}
		}
	}
	best := total
	for s := 0; s <= total; s++ {
		if !reach[s] {
			continue
		}
		d := total - 2*s
		if d < 0 {
			d = -d
		}
		if d < best {
			best = d
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    weights = [int(x) for x in data[1:1 + n]]
    total = sum(weights)
    reach = {0}
    for w in weights:
        reach |= {r + w for r in reach}
    print(min(abs(total - 2 * r) for r in reach))

main()
`},
		tooSlowIn: map[string]string{
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    weights = [int(x) for x in data[1:1 + n]]
    total = sum(weights)
    best = total
    for mask in range(1 << n):
        s = 0
        for i in range(n):
            if mask >> i & 1:
                s += weights[i]
        d = abs(total - 2 * s)
        if d < best:
            best = d
    print(best)

main()
`},
	},
	{
		slug: "group-leaders", title: "Group leaders", difficulty: "easy", topic: "Sorting",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~group name score~.

**Output.** One ~group name~ per line, sorted by group, naming whoever has the
highest score in that group. If two names in a group tie, the one earlier in the
alphabet leads.

    Input               Output
    3                   blue cy
    red ada 5           red bob
    red bob 7
    blue cy 1`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char grp[1000][32], best[1000][32];
    static long long bestScore[1000];
    int groups = 0;
    /* keep the leader per group, then print sorted by group */
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
	_ = sort.Strings
	// keep the leader per group, then print sorted by group
	_ = n
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    # keep the leader per group, then print sorted by group

main()
`),
		tests: []store.Check{
			tokens("The highest score leads each group", "3\nred ada 5\nred bob 7\nblue cy 1\n", "blue cy red bob"),
			tokens("No entries, no groups", "0\n", ""),
			tokens("A group of one", "1\nsolo ada 3\n", "solo ada"),
			hid2("A tie goes to the earlier name", "2\ng zoe 5\ng ada 5\n", "g ada"),
			hid2("Negative scores still rank", "2\ng a -5\ng b -1\n", "g b"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char grp[1000][32], best[1000][32];
static long long bestScore[1000];
static int groups;

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    return strcmp(grp[x], grp[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char g[32], name[32];
        long long score;
        if (scanf("%31s %31s %lld", g, name, &score) != 3) return 1;
        int at = -1;
        for (int j = 0; j < groups; j++) if (strcmp(grp[j], g) == 0) at = j;
        if (at < 0) {
            at = groups++;
            strcpy(grp[at], g); strcpy(best[at], name); bestScore[at] = score;
            continue;
        }
        if (score > bestScore[at] || (score == bestScore[at] && strcmp(name, best[at]) < 0)) {
            strcpy(best[at], name); bestScore[at] = score;
        }
    }
    static int order[1000];
    for (int i = 0; i < groups; i++) order[i] = i;
    qsort(order, groups, sizeof *order, cmp);
    for (int i = 0; i < groups; i++) printf("%s %s\n", grp[order[i]], best[order[i]]);
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

type leader struct {
	name  string
	score int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	top := map[string]leader{}
	for i := 0; i < n; i++ {
		var g, name string
		var score int
		fmt.Fscan(in, &g, &name, &score)
		cur, ok := top[g]
		if !ok || score > cur.score || (score == cur.score && name < cur.name) {
			top[g] = leader{name, score}
		}
	}
	groups := make([]string, 0, len(top))
	for g := range top {
		groups = append(groups, g)
	}
	sort.Strings(groups)
	for _, g := range groups {
		fmt.Fprintln(out, g, top[g].name)
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    top = {}
    for i in range(n):
        g, name, score = data[1 + 3 * i], data[2 + 3 * i], int(data[3 + 3 * i])
        cur = top.get(g)
        if cur is None or score > cur[0] or (score == cur[0] and name < cur[1]):
            top[g] = (score, name)
    for g in sorted(top):
        print(g, top[g][1])

main()
`},
	},
	{
		slug: "foreign-alphabet", title: "A foreign alphabet", difficulty: "medium", topic: "Sorting",
		statement: `A dictionary in another language uses the same letters in a different
order.

**Input.** The alphabet on the first line, then a count ~n~, then ~n~ words.

**Output.** The words in that language's order, space-separated. Compare letter
by letter under the given order; when one word is a prefix of another, the
shorter comes first. Every letter in every word appears in the alphabet.

    Input     Output
    ba        ba ab
    2
    ba ab`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    char alphabet[64];
    if (scanf("%63s", alphabet) != 1) return 1;
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* rank each letter, then sort the words by that ranking */
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
	var alphabet string
	fmt.Fscan(in, &alphabet)
	var n int
	fmt.Fscan(in, &n)
	_, _ = sort.Slice, strings.Join
	// rank each letter, then sort the words by that ranking
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    alphabet, n = data[0], int(data[1])
    words = data[2:2 + n]
    # rank each letter, then sort the words by that ranking
    print("")

main()
`),
		tests: []store.Check{
			tokens("Sorts by the given order", "ba\n2\nba ab\n", "ba ab"),
			tokens("A prefix comes first", "ab\n2\nab a\n", "a ab"),
			tokens("Nothing to sort", "abc\n0\n", ""),
			hid2("The usual alphabet is just one order", "abcdefghijklmnopqrstuvwxyz\n2\nb a\n", "a b"),
			hid2("The disagreement can be late in the word", "bazyxwvutsrqponmlkjihgfedc\n2\nzzb zza\n", "zzb zza"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static int rank[256];
static char word[100000][64];

static int cmp(const void *a, const void *b) {
    const char *x = (const char *)a, *y = (const char *)b;
    for (int i = 0; ; i++) {
        if (!x[i] && !y[i]) return 0;
        if (!x[i]) return -1;
        if (!y[i]) return 1;
        if (rank[(unsigned char)x[i]] != rank[(unsigned char)y[i]])
            return rank[(unsigned char)x[i]] - rank[(unsigned char)y[i]];
    }
}

int main(void) {
    char alphabet[64];
    if (scanf("%63s", alphabet) != 1) return 1;
    for (int i = 0; alphabet[i]; i++) rank[(unsigned char)alphabet[i]] = i;
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) if (scanf("%63s", word[i]) != 1) return 1;
    qsort(word, n, 64, cmp);
    for (int i = 0; i < n; i++) printf(i ? " %s" : "%s", word[i]);
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
	var alphabet string
	fmt.Fscan(in, &alphabet)
	rank := map[byte]int{}
	for i := 0; i < len(alphabet); i++ {
		rank[alphabet[i]] = i
	}
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := range words {
		fmt.Fscan(in, &words[i])
	}
	sort.Slice(words, func(i, j int) bool {
		a, b := words[i], words[j]
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return rank[a[k]] < rank[b[k]]
			}
		}
		return len(a) < len(b)
	})
	fmt.Println(strings.Join(words, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    alphabet, n = data[0], int(data[1])
    words = data[2:2 + n]
    rank = {c: i for i, c in enumerate(alphabet)}
    print(" ".join(sorted(words, key=lambda w: [rank[c] for c in w])))

main()
`},
	},
	{
		slug: "round-robin-feed", title: "Round-robin feed", difficulty: "easy", topic: "Sorting",
		statement: `A social feed shows one post from each account in turn, then goes round
again. When an account runs out of posts it is simply skipped; the others carry
on.

**Input.** A count ~a~, then ~a~ lines. Each is a count ~m~ followed by ~m~ post
ids.

**Output.** The feed in order, space-separated.

    Input      Output
    2          1 3 2
    2 1 2
    1 3`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int a;
    if (scanf("%d", &a) != 1) return 1;
    /* take the first post of each account, then the second, and so on */
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
	var a int
	fmt.Fscan(in, &a)
	// take the first post of each account, then the second, and so on
	fmt.Println("")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    a = data[0]
    # take the first post of each account, then the second, and so on
    print("")

main()
`),
		tests: []store.Check{
			tokens("Takes one from each in turn", "2\n2 1 2\n1 3\n", "1 3 2"),
			tokens("A single account posts in order", "1\n3 1 2 3\n", "1 2 3"),
			tokens("No accounts, no feed", "0\n", ""),
			hid2("An exhausted account is skipped, not padded", "2\n1 1\n3 2 3 4\n", "1 2 3 4"),
			hid2("Empty accounts contribute nothing", "3\n0\n1 1\n0\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static int post[1000][1000], len[1000];

int main(void) {
    int a;
    if (scanf("%d", &a) != 1) return 1;
    int longest = 0;
    for (int i = 0; i < a; i++) {
        if (scanf("%d", &len[i]) != 1) return 1;
        if (len[i] > longest) longest = len[i];
        for (int j = 0; j < len[i]; j++) if (scanf("%d", &post[i][j]) != 1) return 1;
    }
    int first = 1;
    for (int r = 0; r < longest; r++)
        for (int i = 0; i < a; i++)
            if (r < len[i]) { printf(first ? "%d" : " %d", post[i][r]); first = 0; }
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
	var a int
	fmt.Fscan(in, &a)
	accounts := make([][]int, a)
	longest := 0
	for i := range accounts {
		var m int
		fmt.Fscan(in, &m)
		accounts[i] = make([]int, m)
		for j := range accounts[i] {
			fmt.Fscan(in, &accounts[i][j])
		}
		if m > longest {
			longest = m
		}
	}
	var out []string
	for r := 0; r < longest; r++ {
		for _, acc := range accounts {
			if r < len(acc) {
				out = append(out, strconv.Itoa(acc[r]))
			}
		}
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    a, pos = data[0], 1
    accounts = []
    for _ in range(a):
        m = data[pos]; pos += 1
        accounts.append(data[pos:pos + m])
        pos += m
    out = []
    for r in range(max((len(x) for x in accounts), default=0)):
        for acc in accounts:
            if r < len(acc):
                out.append(str(acc[r]))
    print(" ".join(out))

main()
`},
	},
	{
		slug: "most-improved", title: "Most improved", difficulty: "medium", topic: "Sorting",
		statement: `Two leaderboards list names best-first.

**Input.** ~b a~ on the first line, then ~b~ names for the old board and ~a~ for
the new.

**Output.** Everyone on the new board, space-separated, ordered by how many
places they climbed, biggest climb first, ties broken alphabetically. Somebody
who was not on the old board is treated as having been just off the end of it —
position ~b~.

    Input        Output
    3 3          c a b
    a b c
    c a b

Only people on the new board are returned; whoever dropped off it is gone.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int b, a;
    if (scanf("%d %d", &b, &a) != 2) return 1;
    /* work out each climb, then order by it */
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
	var b, a int
	fmt.Fscan(in, &b, &a)
	_, _ = sort.Slice, strings.Join
	// work out each climb, then order by it
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    b, a = int(data[0]), int(data[1])
    before = data[2:2 + b]
    after = data[2 + b:2 + b + a]
    # work out each climb, then order by it
    print("")

main()
`),
		tests: []store.Check{
			tokens("The biggest climb leads", "3 3\na b c\nc a b\n", "c a b"),
			tokens("A newcomer climbs from just off the board", "1 2\na\nz a\n", "z a"),
			tokens("An empty new board", "1 0\na\n", ""),
			hid2("Nobody moved, so alphabetical", "2 2\na b\na b\n", "a b"),
			hid2("Someone who dropped off is not reported", "2 1\na b\nb\n", "b"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char before[100000][32], after[100000][32];
static int gain[100000], order[100000];

static int cmp(const void *x, const void *y) {
    int i = *(const int *)x, j = *(const int *)y;
    if (gain[i] != gain[j]) return gain[j] - gain[i];
    return strcmp(after[i], after[j]);
}

int main(void) {
    int b, a;
    if (scanf("%d %d", &b, &a) != 2) return 1;
    for (int i = 0; i < b; i++) if (scanf("%31s", before[i]) != 1) return 1;
    for (int i = 0; i < a; i++) if (scanf("%31s", after[i]) != 1) return 1;
    for (int i = 0; i < a; i++) {
        int was = b;
        for (int j = 0; j < b; j++) if (strcmp(before[j], after[i]) == 0) { was = j; break; }
        gain[i] = was - i;
        order[i] = i;
    }
    qsort(order, a, sizeof *order, cmp);
    for (int i = 0; i < a; i++) printf(i ? " %s" : "%s", after[order[i]]);
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
	var b, a int
	fmt.Fscan(in, &b, &a)
	was := map[string]int{}
	for i := 0; i < b; i++ {
		var s string
		fmt.Fscan(in, &s)
		was[s] = i
	}
	after := make([]string, a)
	for i := range after {
		fmt.Fscan(in, &after[i])
	}
	gain := map[string]int{}
	for i, name := range after {
		old, ok := was[name]
		if !ok {
			old = b
		}
		gain[name] = old - i
	}
	sort.Slice(after, func(i, j int) bool {
		if gain[after[i]] != gain[after[j]] {
			return gain[after[i]] > gain[after[j]]
		}
		return after[i] < after[j]
	})
	fmt.Println(strings.Join(after, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    b, a = int(data[0]), int(data[1])
    before = data[2:2 + b]
    after = data[2 + b:2 + b + a]
    was = {n: i for i, n in enumerate(before)}
    gain = {n: was.get(n, b) - i for i, n in enumerate(after)}
    print(" ".join(sorted(after, key=lambda n: (-gain[n], n))))

main()
`},
	},
	{
		slug: "seat-overflow", title: "Seat overflow", difficulty: "medium", topic: "Simulation",
		statement: `A hall seats people row by row. Requests are honoured in order. If the
preferred row is full, the person takes the next row with space; if there is no
such row, the hall cannot seat everyone and the whole plan fails.

**Input.** ~r n~ on the first line, then ~r~ row capacities, then ~n~ lines of
~name preferred_row~ (rows numbered from 0).

**Output.** One ~name row~ per line in request order, or the single word
~impossible~.

    Input      Output
    2 2        a 0
    1 1        b 1
    a 0
    b 0

Nobody moves backwards: overflow only ever goes to a later row.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int r, n;
    if (scanf("%d %d", &r, &n) != 2) return 1;
    static int left[1000];
    for (int i = 0; i < r; i++) if (scanf("%d", &left[i]) != 1) return 1;
    /* seat each request, or print impossible */
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
	var r, n int
	fmt.Fscan(in, &r, &n)
	left := make([]int, r)
	for i := range left {
		fmt.Fscan(in, &left[i])
	}
	// seat each request, or print impossible
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    r, n = int(data[0]), int(data[1])
    left = [int(x) for x in data[2:2 + r]]
    # seat each request, or print impossible

main()
`),
		tests: []store.Check{
			tokens("Overflow moves to the next row", "2 2\n1 1\na 0\nb 0\n", "a 0 b 1"),
			tokens("Everyone gets their preference", "2 1\n1 1\na 1\n", "a 1"),
			tokens("Nobody to seat", "1 0\n1\n", ""),
			tokens("A hall that cannot hold them", "1 1\n0\na 0\n", "impossible"),
			hid2("Overflow never goes backwards", "2 2\n5 1\na 1\nb 1\n", "impossible"),
			hid2("Overflow can skip a full row", "3 3\n1 0 2\na 0\nb 0\nc 0\n", "a 0 b 2 c 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int r, n;
    if (scanf("%d %d", &r, &n) != 2) return 1;
    static int left[1000];
    static char name[1000][32];
    static int row[1000];
    for (int i = 0; i < r; i++) if (scanf("%d", &left[i]) != 1) return 1;
    for (int i = 0; i < n; i++) {
        int pref;
        if (scanf("%31s %d", name[i], &pref) != 2) return 1;
        int at = pref;
        while (at < r && left[at] == 0) at++;
        if (at >= r) { printf("impossible\n"); return 0; }
        left[at]--;
        row[i] = at;
    }
    for (int i = 0; i < n; i++) printf("%s %d\n", name[i], row[i]);
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
	var r, n int
	fmt.Fscan(in, &r, &n)
	left := make([]int, r)
	for i := range left {
		fmt.Fscan(in, &left[i])
	}
	names := make([]string, n)
	rows := make([]int, n)
	for i := 0; i < n; i++ {
		var pref int
		fmt.Fscan(in, &names[i], &pref)
		at := pref
		for at < r && left[at] == 0 {
			at++
		}
		if at >= r {
			fmt.Fprintln(out, "impossible")
			return
		}
		left[at]--
		rows[i] = at
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, names[i], rows[i])
	}
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    r, n = int(data[0]), int(data[1])
    left = [int(x) for x in data[2:2 + r]]
    base = 2 + r
    out = []
    for i in range(n):
        name, pref = data[base + 2 * i], int(data[base + 1 + 2 * i])
        at = pref
        while at < r and left[at] == 0:
            at += 1
        if at >= r:
            print("impossible")
            return
        left[at] -= 1
        out.append(f"{name} {at}")
    print("\n".join(out))

main()
`},
	},
	{
		slug: "league-table", title: "League table", difficulty: "medium", topic: "Sorting",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~home away home_goals away_goals~.

**Output.** The team names, space-separated, ordered by points, then goal
difference, then alphabetically — each highest first except the name.

A win is 3 points, a draw is 1 point each, a loss is nothing. Goal difference is
goals scored minus goals conceded across every match. Every team that appears in
a result is in the table, even if it lost everything.

    Input        Output
    1            A B
    A B 2 1`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static char team[2000][32];
    static int pts[2000], gd[2000];
    int teams = 0;
    /* score every match, then order the table */
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
	var n int
	fmt.Fscan(in, &n)
	_, _ = sort.Slice, strings.Join
	// score every match, then order the table
	_ = n
	fmt.Println("")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    # score every match, then order the table
    print("")

main()
`),
		tests: []store.Check{
			tokens("A win puts you top", "1\nA B 2 1\n", "A B"),
			tokens("A draw is settled by name", "1\nB A 1 1\n", "A B"),
			tokens("No matches, no table", "0\n", ""),
			hid2("Goal difference separates equal points", "2\nA C 5 0\nB C 1 0\n", "A B C"),
			hid2("Three points beats a better goal difference", "3\nA B 1 0\nC D 9 0\nC A 0 1\n", "A C B D"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char team[2000][32];
static int pts[2000], gd[2000];
static int teams;

static int slot(const char *w) {
    for (int i = 0; i < teams; i++) if (strcmp(team[i], w) == 0) return i;
    strcpy(team[teams], w);
    pts[teams] = 0; gd[teams] = 0;
    return teams++;
}

static int cmp(const void *x, const void *y) {
    int i = *(const int *)x, j = *(const int *)y;
    if (pts[i] != pts[j]) return pts[j] - pts[i];
    if (gd[i] != gd[j]) return gd[j] - gd[i];
    return strcmp(team[i], team[j]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char h[32], a[32];
        int hg, ag;
        if (scanf("%31s %31s %d %d", h, a, &hg, &ag) != 4) return 1;
        int x = slot(h), y = slot(a);
        gd[x] += hg - ag; gd[y] += ag - hg;
        if (hg > ag) pts[x] += 3;
        else if (ag > hg) pts[y] += 3;
        else { pts[x]++; pts[y]++; }
    }
    static int order[2000];
    for (int i = 0; i < teams; i++) order[i] = i;
    qsort(order, teams, sizeof *order, cmp);
    for (int i = 0; i < teams; i++) printf(i ? " %s" : "%s", team[order[i]]);
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
	var n int
	fmt.Fscan(in, &n)
	pts, gd := map[string]int{}, map[string]int{}
	for i := 0; i < n; i++ {
		var h, a string
		var hg, ag int
		fmt.Fscan(in, &h, &a, &hg, &ag)
		for _, t := range []string{h, a} {
			if _, ok := pts[t]; !ok {
				pts[t], gd[t] = 0, 0
			}
		}
		gd[h] += hg - ag
		gd[a] += ag - hg
		switch {
		case hg > ag:
			pts[h] += 3
		case ag > hg:
			pts[a] += 3
		default:
			pts[h]++
			pts[a]++
		}
	}
	teams := make([]string, 0, len(pts))
	for t := range pts {
		teams = append(teams, t)
	}
	sort.Slice(teams, func(i, j int) bool {
		a, b := teams[i], teams[j]
		if pts[a] != pts[b] {
			return pts[a] > pts[b]
		}
		if gd[a] != gd[b] {
			return gd[a] > gd[b]
		}
		return a < b
	})
	fmt.Println(strings.Join(teams, " "))
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    pts, gd = {}, {}
    for i in range(n):
        h, a = data[1 + 4 * i], data[2 + 4 * i]
        hg, ag = int(data[3 + 4 * i]), int(data[4 + 4 * i])
        for t in (h, a):
            pts.setdefault(t, 0)
            gd.setdefault(t, 0)
        gd[h] += hg - ag
        gd[a] += ag - hg
        if hg > ag:
            pts[h] += 3
        elif ag > hg:
            pts[a] += 3
        else:
            pts[h] += 1
            pts[a] += 1
    print(" ".join(sorted(pts, key=lambda t: (-pts[t], -gd[t], t))))

main()
`},
	},
}

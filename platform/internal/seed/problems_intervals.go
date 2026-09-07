package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// intervalProblems: things with a start and an end, and the sweeps and greedy
// choices that make sense of a pile of them.
var intervalProblems = []seedProblem{
	{
		slug: "merge-bookings", title: "Merge bookings", difficulty: "medium", topic: "Intervals",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~start end~, in no particular
order.

**Output.** The fewest bookings covering the same time, sorted by start, one
~start end~ per line. Bookings that merely touch — one ending exactly where the
next begins — become one booking.

    Input        Output
    3            1 4
    1 3          6 8
    2 4
    6 8`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* sort by start, then absorb each booking into the last */
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
	_ = sort.Slice
	// sort by start, then absorb each booking into the last
	_ = n
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    spans = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    # sort by start, then absorb each booking into the last

main()
`),
		tests: []store.Check{
			tokens("Overlapping bookings become one", "3\n1 3\n2 4\n6 8\n", "1 4 6 8"),
			tokens("Touching bookings merge too", "2\n1 2\n2 3\n", "1 3"),
			tokens("Nothing booked", "0\n", ""),
			tokens("Separate bookings stay separate", "2\n1 2\n5 6\n", "1 2 5 6"),
			hid2("Order of arrival does not matter", "3\n6 8\n1 3\n2 4\n", "1 4 6 8"),
			hid2("A booking swallowed whole by another", "2\n1 10\n2 3\n", "1 10"),
			hid2("A chain of overlaps collapses to one", "3\n1 3\n2 5\n4 9\n", "1 9"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (lo[x] != lo[y]) return lo[x] < lo[y] ? -1 : 1;
    return (hi[x] > hi[y]) - (hi[x] < hi[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    long long curLo = 0, curHi = 0;
    int open = 0;
    for (int i = 0; i < n; i++) {
        long long s = lo[idx[i]], e = hi[idx[i]];
        if (open && s <= curHi) { if (e > curHi) curHi = e; continue; }
        if (open) printf("%lld %lld\n", curLo, curHi);
        curLo = s; curHi = e; open = 1;
    }
    if (open) printf("%lld %lld\n", curLo, curHi);
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

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	spans := make([]span, n)
	for i := range spans {
		fmt.Fscan(in, &spans[i].lo, &spans[i].hi)
	}
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].lo != spans[j].lo {
			return spans[i].lo < spans[j].lo
		}
		return spans[i].hi < spans[j].hi
	})
	var merged []span
	for _, s := range spans {
		if k := len(merged); k > 0 && s.lo <= merged[k-1].hi {
			if s.hi > merged[k-1].hi {
				merged[k-1].hi = s.hi
			}
			continue
		}
		merged = append(merged, s)
	}
	for _, s := range merged {
		fmt.Fprintln(out, s.lo, s.hi)
	}
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    spans = sorted((data[1 + 2 * i], data[2 + 2 * i]) for i in range(n))
    out = []
    for s, e in spans:
        if out and s <= out[-1][1]:
            out[-1] = (out[-1][0], max(out[-1][1], e))
        else:
            out.append((s, e))
    print("\n".join(f"{s} {e}" for s, e in out))

main()
`},
	},
	{
		slug: "double-booked", title: "Double booked", difficulty: "medium", topic: "Intervals",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~start end~.

**Output.** The first two bookings that genuinely overlap, as ~s1 e1 s2 e2~ with
the earlier-starting one first, or ~none~ if the diary is clean. A booking that
ends exactly when another begins is not a clash. "First" means the clash whose
earlier booking starts soonest.

    Input        Output
    2            1 3 2 4
    1 3
    2 4`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* sort by start; a clash, if any, is between neighbours */
    printf("none\n");
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
	_ = sort.Slice
	// sort by start; a clash, if any, is between neighbours
	fmt.Println("none")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    spans = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    # sort by start; a clash, if any, is between neighbours
    print("none")

main()
`),
		tests: []store.Check{
			tokens("Finds the overlap", "2\n1 3\n2 4\n", "1 3 2 4"),
			tokens("Touching is not clashing", "2\n1 2\n2 3\n", "none"),
			tokens("An empty diary", "0\n", "none"),
			tokens("A single booking clashes with nothing", "1\n1 5\n", "none"),
			hid2("Order of arrival does not matter", "2\n2 4\n1 3\n", "1 3 2 4"),
			hid2("A booking inside another is a clash", "2\n1 10\n2 3\n", "1 10 2 3"),
			hid2("The earliest clash is reported", "3\n1 9\n2 3\n5 6\n", "1 9 2 3"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (lo[x] != lo[y]) return lo[x] < lo[y] ? -1 : 1;
    return (hi[x] > hi[y]) - (hi[x] < hi[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    for (int i = 0; i + 1 < n; i++) {
        int a = idx[i], b = idx[i + 1];
        if (lo[b] < hi[a]) {
            printf("%lld %lld %lld %lld\n", lo[a], hi[a], lo[b], hi[b]);
            return 0;
        }
    }
    printf("none\n");
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

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	spans := make([]span, n)
	for i := range spans {
		fmt.Fscan(in, &spans[i].lo, &spans[i].hi)
	}
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].lo != spans[j].lo {
			return spans[i].lo < spans[j].lo
		}
		return spans[i].hi < spans[j].hi
	})
	for i := 0; i+1 < n; i++ {
		if spans[i+1].lo < spans[i].hi {
			fmt.Println(spans[i].lo, spans[i].hi, spans[i+1].lo, spans[i+1].hi)
			return
		}
	}
	fmt.Println("none")
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    spans = sorted((data[1 + 2 * i], data[2 + 2 * i]) for i in range(n))
    for a, b in zip(spans, spans[1:]):
        if b[0] < a[1]:
            print(a[0], a[1], b[0], b[1])
            return
    print("none")

main()
`},
	},
	{
		slug: "rooms-needed", title: "Rooms needed", difficulty: "hard", topic: "Intervals",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~start end~.

**Output.** The fewest rooms that could hold every booking — the largest number
happening at the same moment. A booking that ends exactly when another starts can
reuse the room.

    Input        Output
    3            2
    1 3
    2 4
    3 5

Counting how many are live at each booking's start, one booking at a time, is
quadratic. Think of the starts and ends as one sequence of events and walk it in
order.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* build one list of events and sweep it */
    printf("%d\n", 0);
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
	_ = sort.Slice
	// build one list of events and sweep it
	fmt.Println(0)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    # build one list of events and sweep it
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Two at once at the busiest moment", "3\n1 3\n2 4\n3 5\n", "2"),
			tokens("An empty diary needs no rooms", "0\n", "0"),
			tokens("Bookings back to back share a room", "3\n1 2\n2 3\n3 4\n", "1"),
			tokens("Everything at once", "3\n1 9\n1 9\n1 9\n", "3"),
			hid2("A long booking spanning many short ones", "3\n0 100\n1 2\n3 4\n", "2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long at[200000];
static int delta[200000];
static int idx[200000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (at[x] != at[y]) return at[x] < at[y] ? -1 : 1;
    return delta[x] - delta[y];
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    int m = 0;
    for (int i = 0; i < n; i++) {
        long long s, e;
        if (scanf("%lld %lld", &s, &e) != 2) return 1;
        at[m] = s; delta[m] = 1; idx[m] = m; m++;
        at[m] = e; delta[m] = -1; idx[m] = m; m++;
    }
    qsort(idx, m, sizeof *idx, cmp);
    int live = 0, best = 0;
    for (int i = 0; i < m; i++) {
        live += delta[idx[i]];
        if (live > best) best = live;
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
	"sort"
)

type event struct{ at, delta int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	events := make([]event, 0, 2*n)
	for i := 0; i < n; i++ {
		var s, e int
		fmt.Fscan(in, &s, &e)
		events = append(events, event{s, 1}, event{e, -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].at != events[j].at {
			return events[i].at < events[j].at
		}
		return events[i].delta < events[j].delta
	})
	live, best := 0, 0
	for _, ev := range events {
		live += ev.delta
		if live > best {
			best = live
		}
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    events = []
    for i in range(n):
        events.append((data[1 + 2 * i], 1))
        events.append((data[2 + 2 * i], -1))
    events.sort()
    live = best = 0
    for _, d in events:
        live += d
        best = max(best, live)
    print(best)

main()
`},
	},
	{
		slug: "free-slots", title: "Free slots", difficulty: "medium", topic: "Intervals",
		statement: `A day runs from 0 to ~day_end~.

**Input.** ~n day_end min_len~ on the first line, then ~n~ busy periods as
~start end~, unsorted and possibly overlapping.

**Output.** The gaps of at least ~min_len~ where nothing is booked, one
~start end~ per line, in order. The gap before the first booking and after the
last both count.

    Input          Output
    1 10 2         0 2
    2 4            4 10`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    long long dayEnd, minLen;
    if (scanf("%d %lld %lld", &n, &dayEnd, &minLen) != 3) return 1;
    /* sort the busy periods, then report the gaps between them */
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
	var n, dayEnd, minLen int
	fmt.Fscan(in, &n, &dayEnd, &minLen)
	_ = sort.Slice
	// sort the busy periods, then report the gaps between them
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, day_end, min_len = data[0], data[1], data[2]
    busy = sorted((data[3 + 2 * i], data[4 + 2 * i]) for i in range(n))
    # report the gaps between them

main()
`),
		tests: []store.Check{
			tokens("Gaps either side of a booking", "1 10 2\n2 4\n", "0 2 4 10"),
			tokens("A completely empty day is one gap", "0 8 1\n", "0 8"),
			tokens("Short gaps are not offered", "1 3 2\n1 2\n", ""),
			hid2("Overlapping busy periods count once", "2 10 2\n2 6\n3 4\n", "0 2 6 10"),
			hid2("A fully booked day has no gaps", "1 10 1\n0 10\n", ""),
			hid2("A gap exactly the minimum length counts", "1 7 2\n0 5\n", "5 7"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (lo[x] != lo[y]) return lo[x] < lo[y] ? -1 : 1;
    return (hi[x] > hi[y]) - (hi[x] < hi[y]);
}

int main(void) {
    int n;
    long long dayEnd, minLen;
    if (scanf("%d %lld %lld", &n, &dayEnd, &minLen) != 3) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    long long t = 0;
    for (int i = 0; i < n; i++) {
        long long s = lo[idx[i]], e = hi[idx[i]];
        if (s - t >= minLen) printf("%lld %lld\n", t, s);
        if (e > t) t = e;
    }
    if (dayEnd - t >= minLen) printf("%lld %lld\n", t, dayEnd);
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

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, dayEnd, minLen int
	fmt.Fscan(in, &n, &dayEnd, &minLen)
	busy := make([]span, n)
	for i := range busy {
		fmt.Fscan(in, &busy[i].lo, &busy[i].hi)
	}
	sort.Slice(busy, func(i, j int) bool {
		if busy[i].lo != busy[j].lo {
			return busy[i].lo < busy[j].lo
		}
		return busy[i].hi < busy[j].hi
	})
	t := 0
	for _, b := range busy {
		if b.lo-t >= minLen {
			fmt.Fprintln(out, t, b.lo)
		}
		if b.hi > t {
			t = b.hi
		}
	}
	if dayEnd-t >= minLen {
		fmt.Fprintln(out, t, dayEnd)
	}
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, day_end, min_len = data[0], data[1], data[2]
    busy = sorted((data[3 + 2 * i], data[4 + 2 * i]) for i in range(n))
    t, out = 0, []
    for s, e in busy:
        if s - t >= min_len:
            out.append(f"{t} {s}")
        t = max(t, e)
    if day_end - t >= min_len:
        out.append(f"{t} {day_end}")
    print("\n".join(out))

main()
`},
	},
	{
		slug: "talks-attended", title: "Talks attended", difficulty: "hard", topic: "Greedy",
		statement: `A conference runs talks at overlapping times.

**Input.** A count ~n~, then ~n~ lines of ~start end~.

**Output.** The most talks you could sit through, given you must stay for a whole
talk and cannot be in two places at once. A talk starting exactly when another
ends is fine.

    Input        Output
    3            2
    1 3
    2 4
    3 5

The instinct is to take the talk that starts soonest, or the shortest one.
Neither is right. What you want at every moment is to be free again as early as
possible.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* sort by end time, then take every talk that still fits */
    printf("%d\n", 0);
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
	_ = sort.Slice
	// sort by end time, then take every talk that still fits
	fmt.Println(0)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    talks = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    # sort by end time, then take every talk that still fits
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Two of the three fit", "3\n1 3\n2 4\n3 5\n", "2"),
			tokens("Nothing on", "0\n", "0"),
			tokens("Talks back to back all fit", "3\n1 2\n2 3\n3 4\n", "3"),
			tokens("Everything clashes with everything", "3\n1 9\n2 8\n3 7\n", "1"),
			hid2("Earliest start is the wrong choice", "3\n1 10\n2 3\n4 5\n", "2"),
			hid2("Order of arrival does not matter", "3\n3 5\n1 3\n2 4\n", "2"),
			hid2("A single talk", "1\n0 1\n", "1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (hi[x] != hi[y]) return hi[x] < hi[y] ? -1 : 1;
    return (lo[x] > lo[y]) - (lo[x] < lo[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    int count = 0, started = 0;
    long long last = 0;
    for (int i = 0; i < n; i++) {
        int t = idx[i];
        if (!started || lo[t] >= last) { count++; last = hi[t]; started = 1; }
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
	"sort"
)

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	talks := make([]span, n)
	for i := range talks {
		fmt.Fscan(in, &talks[i].lo, &talks[i].hi)
	}
	sort.Slice(talks, func(i, j int) bool { return talks[i].hi < talks[j].hi })
	count, last, started := 0, 0, false
	for _, t := range talks {
		if !started || t.lo >= last {
			count++
			last = t.hi
			started = true
		}
	}
	fmt.Println(count)
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    talks = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    count, last = 0, None
    for s, e in sorted(talks, key=lambda t: t[1]):
        if last is None or s >= last:
            count += 1
            last = e
    print(count)

main()
`},
	},
	{
		slug: "cover-the-day", title: "Cover the day", difficulty: "hard", topic: "Greedy",
		statement: `Volunteers each offer a shift.

**Input.** ~n day_end~ on the first line, then ~n~ lines of ~start end~, in any
order and possibly overlapping.

**Output.** ~yes~ if the shifts together cover every moment from 0 to ~day_end~
with nobody unattended, ~no~ otherwise. A day of length 0 is covered by nobody at
all.

    Input        Output
    2 10         yes
    0 5
    4 10

The greedy move: from wherever you are covered up to, take the shift that reaches
furthest among those starting at or before that point.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    long long dayEnd;
    if (scanf("%d %lld", &n, &dayEnd) != 2) return 1;
    /* sort by start, then extend the covered reach as far as possible */
    printf("no\n");
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
	var n, dayEnd int
	fmt.Fscan(in, &n, &dayEnd)
	_ = sort.Slice
	// sort by start, then extend the covered reach as far as possible
	fmt.Println("no")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, day_end = data[0], data[1]
    shifts = sorted((data[2 + 2 * i], data[3 + 2 * i]) for i in range(n))
    # extend the covered reach as far as possible
    print("no")

main()
`),
		tests: []store.Check{
			tokens("Overlapping shifts cover the day", "2 10\n0 5\n4 10\n", "yes"),
			tokens("A gap in the middle", "2 10\n0 4\n5 10\n", "no"),
			tokens("A day of no length needs nobody", "0 0\n", "yes"),
			tokens("Nobody covers a real day", "0 5\n", "no"),
			hid2("The day must be covered from the very start", "1 10\n1 10\n", "no"),
			hid2("Reaching short of the end fails", "1 10\n0 9\n", "no"),
			hid2("The longest reach is not always the first shift", "3 10\n0 2\n1 3\n1 10\n", "yes"),
			hid2("Shifts touching exactly still cover", "2 10\n0 5\n5 10\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (lo[x] != lo[y]) return lo[x] < lo[y] ? -1 : 1;
    return (hi[x] > hi[y]) - (hi[x] < hi[y]);
}

int main(void) {
    int n;
    long long dayEnd;
    if (scanf("%d %lld", &n, &dayEnd) != 2) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    long long reach = 0;
    int i = 0;
    while (reach < dayEnd) {
        long long best = reach;
        while (i < n && lo[idx[i]] <= reach) {
            if (hi[idx[i]] > best) best = hi[idx[i]];
            i++;
        }
        if (best == reach) { printf("no\n"); return 0; }
        reach = best;
    }
    printf("yes\n");
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

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, dayEnd int
	fmt.Fscan(in, &n, &dayEnd)
	shifts := make([]span, n)
	for i := range shifts {
		fmt.Fscan(in, &shifts[i].lo, &shifts[i].hi)
	}
	sort.Slice(shifts, func(i, j int) bool { return shifts[i].lo < shifts[j].lo })
	reach, i := 0, 0
	for reach < dayEnd {
		best := reach
		for i < len(shifts) && shifts[i].lo <= reach {
			if shifts[i].hi > best {
				best = shifts[i].hi
			}
			i++
		}
		if best == reach {
			fmt.Println("no")
			return
		}
		reach = best
	}
	fmt.Println("yes")
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n, day_end = data[0], data[1]
    shifts = sorted((data[2 + 2 * i], data[3 + 2 * i]) for i in range(n))
    reach, i = 0, 0
    while reach < day_end:
        best = reach
        while i < len(shifts) and shifts[i][0] <= reach:
            best = max(best, shifts[i][1])
            i += 1
        if best == reach:
            print("no")
            return
        reach = best
    print("yes")

main()
`},
	},
	{
		slug: "longest-free-gap", title: "Longest free gap", difficulty: "easy", topic: "Intervals",
		statement: `**Input.** A count ~n~, then ~n~ meetings as ~start end~, unsorted and
possibly overlapping.

**Output.** The longest stretch of free time **between** meetings — from the
moment you are free until the next one begins. Time before the first meeting and
after the last does not count. Fewer than two meetings means no gap at all.

    Input        Output
    2            3
    1 2
    5 6`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* sort by start, then measure each gap from the running end */
    printf("%d\n", 0);
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
	_ = sort.Slice
	// sort by start, then measure each gap from the running end
	fmt.Println(0)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    meetings = sorted((data[1 + 2 * i], data[2 + 2 * i]) for i in range(n))
    # measure each gap from the running end
    print(0)

main()
`),
		tests: []store.Check{
			tokens("The gap between two meetings", "2\n1 2\n5 6\n", "3"),
			tokens("One meeting has no gap", "1\n1 2\n", "0"),
			tokens("No meetings at all", "0\n", "0"),
			tokens("Back to back leaves no gap", "2\n1 2\n2 3\n", "0"),
			hid2("A meeting inside another leaves no gap", "2\n1 10\n2 3\n", "0"),
			hid2("The longest of several gaps", "3\n0 1\n2 3\n9 10\n", "6"),
			hid2("Order of arrival does not matter", "2\n5 6\n1 2\n", "3"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long lo[100000], hi[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    if (lo[x] != lo[y]) return lo[x] < lo[y] ? -1 : 1;
    return (hi[x] > hi[y]) - (hi[x] < hi[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &lo[i], &hi[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    long long best = 0, end = 0;
    int started = 0;
    for (int i = 0; i < n; i++) {
        long long s = lo[idx[i]], e = hi[idx[i]];
        if (started && s - end > best) best = s - end;
        end = started ? (e > end ? e : end) : e;
        started = 1;
    }
    printf("%lld\n", best);
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

type span struct{ lo, hi int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	m := make([]span, n)
	for i := range m {
		fmt.Fscan(in, &m[i].lo, &m[i].hi)
	}
	sort.Slice(m, func(i, j int) bool { return m[i].lo < m[j].lo })
	best, end, started := 0, 0, false
	for _, s := range m {
		if started && s.lo-end > best {
			best = s.lo - end
		}
		if !started || s.hi > end {
			end = s.hi
		}
		started = true
	}
	fmt.Println(best)
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    meetings = sorted((data[1 + 2 * i], data[2 + 2 * i]) for i in range(n))
    best, end = 0, None
    for s, e in meetings:
        if end is not None and s - end > best:
            best = s - end
        end = e if end is None else max(end, e)
    print(best)

main()
`},
	},
	{
		slug: "shared-availability", title: "Shared availability", difficulty: "medium", topic: "Two pointers",
		statement: `Two people each give their free time as sorted, non-overlapping ranges.

**Input.** ~a b~ on the first line, then ~a~ ranges for the first person and ~b~
for the second, each as ~start end~.

**Output.** When they are both free, one ~start end~ per line, in order. A range
of zero length is not a meeting: ranges that merely touch produce nothing.

    Input        Output
    1 1          3 5
    1 5
    3 7

Because both lists are already sorted you never need to look backwards —
whichever range ends first can be retired.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int a, b;
    if (scanf("%d %d", &a, &b) != 2) return 1;
    /* walk both lists together, retiring whichever ends first */
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
	// walk both lists together, retiring whichever ends first
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    a, b = data[0], data[1]
    first = [(data[2 + 2 * i], data[3 + 2 * i]) for i in range(a)]
    base = 2 + 2 * a
    second = [(data[base + 2 * i], data[base + 1 + 2 * i]) for i in range(b)]
    # walk both lists together, retiring whichever ends first

main()
`),
		tests: []store.Check{
			tokens("The overlapping part", "1 1\n1 5\n3 7\n", "3 5"),
			tokens("Touching is not overlapping", "1 1\n1 3\n3 5\n", ""),
			tokens("One person is never free", "1 0\n1 5\n", ""),
			tokens("Neither is ever free", "0 0\n", ""),
			hid2("One range can overlap several", "1 2\n0 10\n1 2\n3 4\n", "1 2 3 4"),
			hid2("Several overlaps in order", "2 1\n1 4\n6 9\n2 7\n", "2 4 6 7"),
			hid2("No overlap at all", "1 1\n1 2\n5 6\n", ""),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static long long al[100000], ah[100000], bl[100000], bh[100000];

int main(void) {
    int a, b;
    if (scanf("%d %d", &a, &b) != 2) return 1;
    for (int i = 0; i < a; i++) if (scanf("%lld %lld", &al[i], &ah[i]) != 2) return 1;
    for (int i = 0; i < b; i++) if (scanf("%lld %lld", &bl[i], &bh[i]) != 2) return 1;
    int i = 0, j = 0;
    while (i < a && j < b) {
        long long lo = al[i] > bl[j] ? al[i] : bl[j];
        long long hi = ah[i] < bh[j] ? ah[i] : bh[j];
        if (lo < hi) printf("%lld %lld\n", lo, hi);
        if (ah[i] < bh[j]) i++; else j++;
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

type span struct{ lo, hi int }

func read(in *bufio.Reader, n int) []span {
	s := make([]span, n)
	for i := range s {
		fmt.Fscan(in, &s[i].lo, &s[i].hi)
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var a, b int
	fmt.Fscan(in, &a, &b)
	first, second := read(in, a), read(in, b)
	i, j := 0, 0
	for i < len(first) && j < len(second) {
		lo, hi := first[i].lo, first[i].hi
		if second[j].lo > lo {
			lo = second[j].lo
		}
		if second[j].hi < hi {
			hi = second[j].hi
		}
		if lo < hi {
			fmt.Fprintln(out, lo, hi)
		}
		if first[i].hi < second[j].hi {
			i++
		} else {
			j++
		}
	}
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    a, b = data[0], data[1]
    first = [(data[2 + 2 * i], data[3 + 2 * i]) for i in range(a)]
    base = 2 + 2 * a
    second = [(data[base + 2 * i], data[base + 1 + 2 * i]) for i in range(b)]
    i = j = 0
    out = []
    while i < len(first) and j < len(second):
        lo = max(first[i][0], second[j][0])
        hi = min(first[i][1], second[j][1])
        if lo < hi:
            out.append(f"{lo} {hi}")
        if first[i][1] < second[j][1]:
            i += 1
        else:
            j += 1
    print("\n".join(out))

main()
`},
	},
	{
		slug: "worst-lateness", title: "Worst lateness", difficulty: "hard", topic: "Greedy",
		statement: `A workshop has jobs. One job runs at a time, starting at time 0, with no
gaps. A job's lateness is when it finishes minus its deadline — negative if it
finishes early.

**Input.** A count ~n~, then ~n~ lines of ~duration deadline~.

**Output.** The smallest possible value of the **largest** lateness, across all
orderings. No jobs at all has a worst lateness of 0.

    Input        Output
    2            2
    4 2
    1 3

Running the short job first finishes both sooner, but the tight deadline slips
further. One ordering rule is always optimal here, and it is worth convincing
yourself why swapping any two adjacent jobs out of that order never helps.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* choose the order, then run the clock forward */
    printf("%d\n", 0);
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
	_ = sort.Slice
	// choose the order, then run the clock forward
	fmt.Println(0)
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    jobs = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    # choose the order, then run the clock forward
    print(0)

main()
`),
		tests: []store.Check{
			tokens("The tight deadline goes first", "2\n4 2\n1 3\n", "2"),
			tokens("No jobs, no lateness", "0\n", "0"),
			tokens("Finishing early gives negative lateness", "1\n1 10\n", "-9"),
			hid2("Everything comfortably on time", "2\n1 5\n1 6\n", "-4"),
			hid2("Order genuinely matters — the other way round gives 3", "2\n3 3\n1 1\n", "1"),
			hid2("Several jobs sharing a deadline", "3\n1 4\n1 4\n1 4\n", "-1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <stdlib.h>

static long long dur[100000], due[100000];
static int idx[100000];

static int cmp(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    return (due[x] > due[y]) - (due[x] < due[y]);
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        if (scanf("%lld %lld", &dur[i], &due[i]) != 2) return 1;
        idx[i] = i;
    }
    qsort(idx, n, sizeof *idx, cmp);
    long long t = 0, worst = 0;
    int started = 0;
    for (int i = 0; i < n; i++) {
        t += dur[idx[i]];
        long long late = t - due[idx[i]];
        if (!started || late > worst) { worst = late; started = 1; }
    }
    printf("%lld\n", started ? worst : 0);
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

type job struct{ dur, due int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	jobs := make([]job, n)
	for i := range jobs {
		fmt.Fscan(in, &jobs[i].dur, &jobs[i].due)
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].due < jobs[j].due })
	t, worst, started := 0, 0, false
	for _, j := range jobs {
		t += j.dur
		if late := t - j.due; !started || late > worst {
			worst, started = late, true
		}
	}
	if !started {
		fmt.Println(0)
		return
	}
	fmt.Println(worst)
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    jobs = [(data[1 + 2 * i], data[2 + 2 * i]) for i in range(n)]
    t, worst = 0, None
    for dur, due in sorted(jobs, key=lambda j: j[1]):
        t += dur
        late = t - due
        worst = late if worst is None else max(worst, late)
    print(0 if worst is None else worst)

main()
`},
	},
	{
		slug: "room-numbers", title: "Room numbers", difficulty: "medium", topic: "Intervals",
		statement: `Bookings arrive in start order and each takes the **lowest-numbered** room
that is free. Rooms are numbered from 0, and a new room is opened only when every
existing one is busy. A room is free again the moment its booking ends.

**Input.** A count ~n~, then ~n~ lines of ~start end~ in start order.

**Output.** The room each booking gets, space-separated.

    Input        Output
    3            0 1 0
    1 3
    2 4
    3 5

The third booking starts exactly when the first ends, so room 0 is free again.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static long long freeAt[100000];
    int rooms = 0;
    /* find the lowest room free by this start time */
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
	var freeAt []int
	_ = freeAt
	// find the lowest room free by this start time
	fmt.Println("")
}
`, `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    free_at = []
    # find the lowest room free by this start time
    print("")

main()
`),
		tests: []store.Check{
			tokens("Reuses the lowest free room", "3\n1 3\n2 4\n3 5\n", "0 1 0"),
			tokens("Nothing booked", "0\n", ""),
			tokens("Back-to-back bookings share one room", "3\n1 2\n2 3\n3 4\n", "0 0 0"),
			tokens("Everything at once opens a room each", "3\n1 9\n1 9\n1 9\n", "0 1 2"),
			hid2("A freed low room beats a freed high one", "3\n0 5\n1 2\n3 9\n", "0 1 1"),
			hid2("One long booking holds its room throughout", "3\n0 100\n1 2\n3 4\n", "0 1 1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    static long long freeAt[100000];
    int rooms = 0, first = 1;
    for (int i = 0; i < n; i++) {
        long long s, e;
        if (scanf("%lld %lld", &s, &e) != 2) return 1;
        int at = -1;
        for (int r = 0; r < rooms; r++) if (freeAt[r] <= s) { at = r; break; }
        if (at < 0) { at = rooms++; }
        freeAt[at] = e;
        printf(first ? "%d" : " %d", at);
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
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var freeAt []int
	out := make([]string, n)
	for i := 0; i < n; i++ {
		var s, e int
		fmt.Fscan(in, &s, &e)
		at := -1
		for r, f := range freeAt {
			if f <= s {
				at = r
				break
			}
		}
		if at < 0 {
			freeAt = append(freeAt, e)
			at = len(freeAt) - 1
		} else {
			freeAt[at] = e
		}
		out[i] = strconv.Itoa(at)
	}
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    data = [int(x) for x in sys.stdin.read().split()]
    n = data[0]
    free_at, out = [], []
    for i in range(n):
        s, e = data[1 + 2 * i], data[2 + 2 * i]
        at = None
        for r, f in enumerate(free_at):
            if f <= s:
                at = r
                break
        if at is None:
            free_at.append(e)
            at = len(free_at) - 1
        else:
            free_at[at] = e
        out.append(str(at))
    print(" ".join(out))

main()
`},
	},
}

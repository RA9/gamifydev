package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// graphProblems: grids, graphs and recursion — problems where the shape of the
// data is the problem, and the answer comes from walking it.
var graphProblems = []seedProblem{
	{
		slug: "flooded-region", title: "Flooded region", difficulty: "medium", topic: "Grids",
		statement: `A map is a grid of ~#~ for land and ~.~ for water.

**Input.** ~rows cols r c~ on the first line, then ~rows~ lines of the map.

**Output.** How many land squares are joined to the one at row ~r~, column ~c~ —
counting squares that touch edge to edge, not corner to corner. If the starting
square is water, the region is empty: print 0.

    Input        Output
    3 3 0 0      3
    ##.
    #..
    ..#`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int rows, cols, r, c;
    if (scanf("%d %d %d %d", &rows, &cols, &r, &c) != 4) return 1;
    static char grid[1000][1005];
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    /* walk outward from (r, c) over land only */
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
	var rows, cols, r, c int
	fmt.Fscan(in, &rows, &cols, &r, &c)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	// walk outward from (r, c) over land only
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols, r, c = (int(x) for x in data[:4])
    grid = data[4:4 + rows]
    # walk outward from (r, c) over land only
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Counts the joined land", "3 3 0 0\n##.\n#..\n..#\n", "3"),
			tokens("Starting on water", "3 3 0 2\n##.\n#..\n..#\n", "0"),
			tokens("A lone square", "1 3 0 2\n..#\n", "1"),
			hid2("Corners do not join", "2 2 0 0\n#.\n.#\n", "1"),
			hid2("The whole map is one region", "2 2 1 1\n##\n##\n", "4"),
			hid2("A map with no land at all", "2 2 0 0\n..\n..\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static char grid[1000][1005];
static int seen[1000][1005];
static int qr[1000000], qc[1000000];

int main(void) {
    int rows, cols, r, c;
    if (scanf("%d %d %d %d", &rows, &cols, &r, &c) != 4) return 1;
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    if (grid[r][c] != '#') { printf("0\n"); return 0; }
    int head = 0, tail = 0, count = 0;
    qr[tail] = r; qc[tail] = c; tail++; seen[r][c] = 1;
    int dr[4] = {1, -1, 0, 0}, dc[4] = {0, 0, 1, -1};
    while (head < tail) {
        int y = qr[head], x = qc[head];
        head++; count++;
        for (int k = 0; k < 4; k++) {
            int ny = y + dr[k], nx = x + dc[k];
            if (ny < 0 || nx < 0 || ny >= rows || nx >= cols) continue;
            if (seen[ny][nx] || grid[ny][nx] != '#') continue;
            seen[ny][nx] = 1;
            qr[tail] = ny; qc[tail] = nx; tail++;
        }
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
	var rows, cols, r, c int
	fmt.Fscan(in, &rows, &cols, &r, &c)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	if grid[r][c] != '#' {
		fmt.Println(0)
		return
	}
	type cell struct{ y, x int }
	seen := map[cell]bool{{r, c}: true}
	queue := []cell{{r, c}}
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++
		for _, d := range []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := cell{cur.y + d.y, cur.x + d.x}
			if n.y < 0 || n.x < 0 || n.y >= rows || n.x >= cols {
				continue
			}
			if seen[n] || grid[n.y][n.x] != '#' {
				continue
			}
			seen[n] = true
			queue = append(queue, n)
		}
	}
	fmt.Println(count)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols, r, c = (int(x) for x in data[:4])
    grid = data[4:4 + rows]
    if grid[r][c] != "#":
        print(0)
        return
    seen = {(r, c)}
    stack = [(r, c)]
    count = 0
    while stack:
        y, x = stack.pop()
        count += 1
        for dy, dx in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            ny, nx = y + dy, x + dx
            if 0 <= ny < rows and 0 <= nx < cols and (ny, nx) not in seen and grid[ny][nx] == "#":
                seen.add((ny, nx))
                stack.append((ny, nx))
    print(count)

main()
`},
	},
	{
		slug: "count-clusters", title: "Count clusters", difficulty: "medium", topic: "Grids",
		statement: `On the same kind of map, count the separate groups of land.

Here squares join **including diagonally** — corner to corner counts. That one
word changes the answer.

**Input.** ~rows cols~ on the first line, then ~rows~ lines of the map.

**Output.** How many clusters there are.

    Input        Output
    3 3          1
    #.#
    .#.
    #.#

Every square touches the middle one at a corner, so it is all one cluster.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    static char grid[1000][1005];
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    /* count the groups, joining diagonally too */
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
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	// count the groups, joining diagonally too
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    # count the groups, joining diagonally too
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Diagonals join it all up", "3 3\n#.#\n.#.\n#.#\n", "1"),
			tokens("Land too far apart to join", "2 3\n#..\n..#\n", "2"),
			tokens("No land at all", "2 2\n..\n..\n", "0"),
			tokens("One solid block", "2 2\n##\n##\n", "1"),
			hid2("Diagonal neighbours are one cluster", "2 2\n#.\n.#\n", "1"),
			hid2("Three separate clusters", "1 5\n#.#.#\n", "3"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static char grid[1000][1005];
static int seen[1000][1005];
static int qr[1000000], qc[1000000];

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    int clusters = 0;
    for (int r = 0; r < rows; r++) for (int c = 0; c < cols; c++) {
        if (grid[r][c] != '#' || seen[r][c]) continue;
        clusters++;
        int head = 0, tail = 0;
        qr[tail] = r; qc[tail] = c; tail++; seen[r][c] = 1;
        while (head < tail) {
            int y = qr[head], x = qc[head];
            head++;
            for (int dy = -1; dy <= 1; dy++) for (int dx = -1; dx <= 1; dx++) {
                int ny = y + dy, nx = x + dx;
                if (ny < 0 || nx < 0 || ny >= rows || nx >= cols) continue;
                if (seen[ny][nx] || grid[ny][nx] != '#') continue;
                seen[ny][nx] = 1;
                qr[tail] = ny; qc[tail] = nx; tail++;
            }
        }
    }
    printf("%d\n", clusters);
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
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	type cell struct{ y, x int }
	seen := map[cell]bool{}
	clusters := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] != '#' || seen[cell{r, c}] {
				continue
			}
			clusters++
			seen[cell{r, c}] = true
			stack := []cell{{r, c}}
			for len(stack) > 0 {
				cur := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						n := cell{cur.y + dy, cur.x + dx}
						if n.y < 0 || n.x < 0 || n.y >= rows || n.x >= cols {
							continue
						}
						if seen[n] || grid[n.y][n.x] != '#' {
							continue
						}
						seen[n] = true
						stack = append(stack, n)
					}
				}
			}
		}
	}
	fmt.Println(clusters)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    seen, clusters = set(), 0
    for r in range(rows):
        for c in range(cols):
            if grid[r][c] != "#" or (r, c) in seen:
                continue
            clusters += 1
            seen.add((r, c))
            stack = [(r, c)]
            while stack:
                y, x = stack.pop()
                for dy in (-1, 0, 1):
                    for dx in (-1, 0, 1):
                        ny, nx = y + dy, x + dx
                        if (0 <= ny < rows and 0 <= nx < cols
                                and (ny, nx) not in seen and grid[ny][nx] == "#"):
                            seen.add((ny, nx))
                            stack.append((ny, nx))
    print(clusters)

main()
`},
	},
	{
		slug: "warehouse-route", title: "Warehouse route", difficulty: "hard", topic: "Graphs",
		statement: `A robot starts at the top-left of a grid and must reach the bottom-right,
moving one square at a time up, down, left or right, never entering a ~#~.

**Input.** ~rows cols~ on the first line, then ~rows~ lines of the grid.

**Output.** The fewest moves, or ~-1~ if it cannot get there. Standing on the
target already is 0 moves.

    Input      Output
    2 2        2
    ..
    ..

"Fewest" rules out following one path and hoping: you need to reach every square
one move away before any square two away.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    static char grid[1000][1005];
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    /* explore in rings out from the start */
    printf("%d\n", -1);
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
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	// explore in rings out from the start
	fmt.Println(-1)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    # explore in rings out from the start
    print(-1)

main()
`),
		tests: []store.Check{
			tokens("Across a clear grid", "2 2\n..\n..\n", "2"),
			tokens("Starting on a blocked square", "1 2\n#.\n", "-1"),
			tokens("Already there", "1 1\n.\n", "0"),
			tokens("Walled off entirely", "3 2\n..\n##\n..\n", "-1"),
			hid2("Around an obstacle", "3 3\n..#\n#..\n...\n", "4"),
			hid2("The target itself is blocked", "2 2\n..\n.#\n", "-1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static char grid[1000][1005];
static int dist[1000][1005];
static int qr[1000000], qc[1000000];

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    if (grid[0][0] == '#' || grid[rows - 1][cols - 1] == '#') { printf("-1\n"); return 0; }
    for (int i = 0; i < rows; i++) for (int j = 0; j < cols; j++) dist[i][j] = -1;
    int head = 0, tail = 0;
    qr[tail] = 0; qc[tail] = 0; tail++; dist[0][0] = 0;
    int dr[4] = {1, -1, 0, 0}, dc[4] = {0, 0, 1, -1};
    while (head < tail) {
        int y = qr[head], x = qc[head];
        head++;
        if (y == rows - 1 && x == cols - 1) { printf("%d\n", dist[y][x]); return 0; }
        for (int k = 0; k < 4; k++) {
            int ny = y + dr[k], nx = x + dc[k];
            if (ny < 0 || nx < 0 || ny >= rows || nx >= cols) continue;
            if (dist[ny][nx] >= 0 || grid[ny][nx] == '#') continue;
            dist[ny][nx] = dist[y][x] + 1;
            qr[tail] = ny; qc[tail] = nx; tail++;
        }
    }
    printf("-1\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"fmt"
	"os"
)

type cell struct{ y, x, d int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	if grid[0][0] == '#' || grid[rows-1][cols-1] == '#' {
		fmt.Println(-1)
		return
	}
	seen := make([][]bool, rows)
	for i := range seen {
		seen[i] = make([]bool, cols)
	}
	seen[0][0] = true
	queue := []cell{{0, 0, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.y == rows-1 && cur.x == cols-1 {
			fmt.Println(cur.d)
			return
		}
		for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			ny, nx := cur.y+d[0], cur.x+d[1]
			if ny < 0 || nx < 0 || ny >= rows || nx >= cols || seen[ny][nx] || grid[ny][nx] == '#' {
				continue
			}
			seen[ny][nx] = true
			queue = append(queue, cell{ny, nx, cur.d + 1})
		}
	}
	fmt.Println(-1)
}
`,
			"python": `import sys
from collections import deque

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    if grid[0][0] == "#" or grid[rows - 1][cols - 1] == "#":
        print(-1)
        return
    seen = {(0, 0)}
    q = deque([(0, 0, 0)])
    while q:
        y, x, d = q.popleft()
        if (y, x) == (rows - 1, cols - 1):
            print(d)
            return
        for dy, dx in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            ny, nx = y + dy, x + dx
            if (0 <= ny < rows and 0 <= nx < cols and (ny, nx) not in seen
                    and grid[ny][nx] != "#"):
                seen.add((ny, nx))
                q.append((ny, nx, d + 1))
    print(-1)

main()
`},
	},
	{
		slug: "everything-needed", title: "Everything needed", difficulty: "medium", topic: "Graphs",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~package dep1 dep2 ...~ (a package
with no dependencies is written as just its name). The last line names the
package to install.

**Output.** The sorted names of everything that must be installed for it,
space-separated, following dependencies all the way down. The package itself is
not listed unless something it needs depends back on it. Dependencies can form
cycles; walking one forever is not an answer.

    Input          Output
    3              core lib
    app lib
    lib core
    core
    app`,
		timeLimitMs: 4000,
		starters:    triLines("follow the dependencies, remembering what you have already seen"),
		tests: []store.Check{
			tokens("Follows the chain down", "3\napp lib\nlib core\ncore\napp\n", "core lib"),
			tokens("Nothing needed", "1\napp\napp\n", ""),
			tokens("An unknown package needs nothing", "0\nghost\n", ""),
			hid2("A cycle back to the start includes it", "2\na b\nb a\na\n", "a b"),
			hid2("A shared dependency is listed once", "4\na b c\nb d\nc d\nd\na\n", "b c d"),
			hid2("A cycle further down still terminates", "3\na b\nb c\nc b\na\n", "b c"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static char name[1000][32], deps[1000][1005];
static int macros;
static char found[1000][32];
static int foundN;

static int findMacro(const char *w) {
    for (int i = 0; i < macros; i++) if (strcmp(name[i], w) == 0) return i;
    return -1;
}

static int seenAlready(const char *w) {
    for (int i = 0; i < foundN; i++) if (strcmp(found[i], w) == 0) return 1;
    return 0;
}

static int cmpstr(const void *a, const void *b) { return strcmp((const char *)a, (const char *)b); }

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    getchar();
    for (int i = 0; i < n; i++) {
        char line[1100];
        if (!fgets(line, sizeof line, stdin)) break;
        line[strcspn(line, "\n")] = '\0';
        char *sp = strchr(line, ' ');
        if (sp) { *sp = '\0'; strcpy(deps[macros], sp + 1); }
        else deps[macros][0] = '\0';
        strcpy(name[macros], line);
        macros++;
    }
    char target[64];
    if (scanf("%63s", target) != 1) return 1;
    static char stack[100000][32];
    int top = 0;
    int at = findMacro(target);
    if (at >= 0) {
        char copy[1005];
        strcpy(copy, deps[at]);
        char *tok = strtok(copy, " ");
        while (tok) { strcpy(stack[top++], tok); tok = strtok(NULL, " "); }
    }
    while (top > 0) {
        char cur[32];
        strcpy(cur, stack[--top]);
        if (seenAlready(cur)) continue;
        strcpy(found[foundN++], cur);
        int m = findMacro(cur);
        if (m < 0) continue;
        char copy[1005];
        strcpy(copy, deps[m]);
        char *tok = strtok(copy, " ");
        while (tok) { strcpy(stack[top++], tok); tok = strtok(NULL, " "); }
    }
    qsort(found, foundN, 32, cmpstr);
    for (int i = 0; i < foundN; i++) printf(i ? " %s" : "%s", found[i]);
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
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	deps := map[string][]string{}
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		fields := strings.Fields(strings.TrimRight(raw, "\r\n"))
		if len(fields) == 0 {
			continue
		}
		deps[fields[0]] = fields[1:]
	}
	var target string
	fmt.Fscan(in, &target)
	seen := map[string]bool{}
	stack := append([]string{}, deps[target]...)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[cur] {
			continue
		}
		seen[cur] = true
		stack = append(stack, deps[cur]...)
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    deps = {}
    for line in raw[1:1 + n]:
        parts = line.split()
        if parts:
            deps[parts[0]] = parts[1:]
    target = raw[1 + n].strip()
    seen, stack = set(), list(deps.get(target, []))
    while stack:
        cur = stack.pop()
        if cur in seen:
            continue
        seen.add(cur)
        stack.extend(deps.get(cur, []))
    print(" ".join(sorted(seen)))

main()
`},
	},
	{
		slug: "install-order", title: "Install order", difficulty: "hard", topic: "Graphs",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~package dep1 dep2 ...~.

**Output.** An order to install every package such that nothing is installed
before something it needs, space-separated. Where several could go next, install
the alphabetically first — so there is exactly one right answer. If the
dependencies form a cycle, no order exists: print ~cycle~.

    Input        Output
    2            lib app
    app lib
    lib

Packages that only ever appear as somebody else's dependency are still packages
and still get installed.`,
		timeLimitMs: 4000,
		starters:    triLines("count what each package waits on, then release them in order"),
		tests: []store.Check{
			tokens("Dependencies come first", "2\napp lib\nlib\napp\n", "lib app"),
			tokens("A cycle has no order", "2\na b\nb a\n", "cycle"),
			tokens("Nothing to install", "0\n", ""),
			tokens("Independent packages go alphabetically", "2\nb\na\n", "a b"),
			hid2("A package only named as a dependency is installed", "1\napp lib\n", "lib app"),
			hid2("Alphabetical among everything currently possible", "3\nz\na z\nb\n", "b z a"),
			hid2("A package depending on itself is a cycle", "1\na a\n", "cycle"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define MAXN 2000
static char name[MAXN][32];
static int total;
static int indeg[MAXN];
static int edges[MAXN][MAXN];
static int edgeN[MAXN];
static int done[MAXN];

static int slot(const char *w) {
    for (int i = 0; i < total; i++) if (strcmp(name[i], w) == 0) return i;
    strcpy(name[total], w);
    return total++;
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    getchar();
    for (int i = 0; i < n; i++) {
        char line[1100];
        if (!fgets(line, sizeof line, stdin)) break;
        line[strcspn(line, "\n")] = '\0';
        char *tok = strtok(line, " ");
        if (!tok) continue;
        int me = slot(tok);
        while ((tok = strtok(NULL, " "))) {
            int dep = slot(tok);
            edges[dep][edgeN[dep]++] = me;
            indeg[me]++;
        }
    }
    int emitted = 0, first = 1;
    while (emitted < total) {
        int pick = -1;
        for (int i = 0; i < total; i++) {
            if (done[i] || indeg[i] != 0) continue;
            if (pick < 0 || strcmp(name[i], name[pick]) < 0) pick = i;
        }
        if (pick < 0) { printf("cycle\n"); return 0; }
        done[pick] = 1;
        emitted++;
        printf(first ? "%s" : " %s", name[pick]);
        first = 0;
        for (int k = 0; k < edgeN[pick]; k++) indeg[edges[pick][k]]--;
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
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	indeg := map[string]int{}
	unlocks := map[string][]string{}
	names := map[string]bool{}
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		f := strings.Fields(strings.TrimRight(raw, "\r\n"))
		if len(f) == 0 {
			continue
		}
		names[f[0]] = true
		if _, ok := indeg[f[0]]; !ok {
			indeg[f[0]] = 0
		}
		for _, d := range f[1:] {
			names[d] = true
			if _, ok := indeg[d]; !ok {
				indeg[d] = 0
			}
			unlocks[d] = append(unlocks[d], f[0])
			indeg[f[0]]++
		}
	}
	var out []string
	done := map[string]bool{}
	for len(out) < len(names) {
		var pick string
		for k := range names {
			if done[k] || indeg[k] != 0 {
				continue
			}
			if pick == "" || k < pick {
				pick = k
			}
		}
		if pick == "" {
			fmt.Println("cycle")
			return
		}
		done[pick] = true
		out = append(out, pick)
		for _, m := range unlocks[pick] {
			indeg[m]--
		}
	}
	_ = sort.Strings
	fmt.Println(strings.Join(out, " "))
}
`,
			"python": `import sys
import heapq

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    indeg, unlocks, names = {}, {}, set()
    for line in raw[1:1 + n]:
        f = line.split()
        if not f:
            continue
        names.add(f[0])
        indeg.setdefault(f[0], 0)
        for d in f[1:]:
            names.add(d)
            indeg.setdefault(d, 0)
            unlocks.setdefault(d, []).append(f[0])
            indeg[f[0]] += 1
    ready = [k for k in names if indeg[k] == 0]
    heapq.heapify(ready)
    out = []
    while ready:
        k = heapq.heappop(ready)
        out.append(k)
        for m in unlocks.get(k, []):
            indeg[m] -= 1
            if indeg[m] == 0:
                heapq.heappush(ready, m)
    print(" ".join(out) if len(out) == len(names) else "cycle")

main()
`},
	},
	{
		slug: "friend-groups", title: "Friend groups", difficulty: "medium", topic: "Graphs",
		statement: `**Input.** ~p f~ on the first line, then ~p~ names, then ~f~ lines of
~a b~ meaning those two are friends.

**Output.** Two numbers: how many separate groups of friends there are, and how
many people are in the largest. Friendship goes both ways, and a friend of a
friend is in the same group. Somebody with no friends is a group of one.

    Input        Output
    3 1          2 2
    a b c
    a b

Nobody at all is ~0 0~.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int p, f;
    if (scanf("%d %d", &p, &f) != 2) return 1;
    /* join friends into groups, then measure them */
    printf("%d %d\n", 0, 0);
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
	var p, f int
	fmt.Fscan(in, &p, &f)
	// join friends into groups, then measure them
	fmt.Println(0, 0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    p, f = int(data[0]), int(data[1])
    people = data[2:2 + p]
    # join friends into groups, then measure them
    print(0, 0)

main()
`),
		tests: []store.Check{
			tokens("A pair and a loner", "3 1\na b c\na b\n", "2 2"),
			tokens("Nobody at all", "0 0\n", "0 0"),
			tokens("Everyone alone", "2 0\na b\n", "2 1"),
			tokens("Everyone connected", "3 2\na b c\na b\nb c\n", "1 3"),
			hid2("A friend of a friend joins the group", "4 3\na b c d\na b\nc d\nb c\n", "1 4"),
			hid2("Two groups of two", "4 2\na b c d\na b\nc d\n", "2 2"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

static char name[1000][32];
static int parent[1000], size[1000];
static int total;

static int slot(const char *w) {
    for (int i = 0; i < total; i++) if (strcmp(name[i], w) == 0) return i;
    strcpy(name[total], w);
    parent[total] = total;
    size[total] = 1;
    return total++;
}

static int find(int x) { while (parent[x] != x) { parent[x] = parent[parent[x]]; x = parent[x]; } return x; }

int main(void) {
    int p, f;
    if (scanf("%d %d", &p, &f) != 2) return 1;
    for (int i = 0; i < p; i++) { char w[32]; if (scanf("%31s", w) != 1) return 1; slot(w); }
    for (int i = 0; i < f; i++) {
        char a[32], b[32];
        if (scanf("%31s %31s", a, b) != 2) return 1;
        int x = find(slot(a)), y = find(slot(b));
        if (x != y) { parent[x] = y; size[y] += size[x]; }
    }
    int groups = 0, biggest = 0;
    for (int i = 0; i < total; i++) {
        if (find(i) != i) continue;
        groups++;
        if (size[i] > biggest) biggest = size[i];
    }
    printf("%d %d\n", groups, biggest);
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
	var p, f int
	fmt.Fscan(in, &p, &f)
	adj := map[string][]string{}
	people := make([]string, p)
	for i := range people {
		fmt.Fscan(in, &people[i])
		adj[people[i]] = nil
	}
	for i := 0; i < f; i++ {
		var a, b string
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	seen := map[string]bool{}
	groups, biggest := 0, 0
	for _, start := range people {
		if seen[start] {
			continue
		}
		groups++
		seen[start] = true
		stack, size := []string{start}, 0
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			size++
			for _, n := range adj[cur] {
				if !seen[n] {
					seen[n] = true
					stack = append(stack, n)
				}
			}
		}
		if size > biggest {
			biggest = size
		}
	}
	fmt.Println(groups, biggest)
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    p, f = int(data[0]), int(data[1])
    people = data[2:2 + p]
    adj = {name: [] for name in people}
    base = 2 + p
    for i in range(f):
        a, b = data[base + 2 * i], data[base + 1 + 2 * i]
        adj.setdefault(a, []).append(b)
        adj.setdefault(b, []).append(a)
    seen, groups, biggest = set(), 0, 0
    for start in people:
        if start in seen:
            continue
        groups += 1
        seen.add(start)
        stack, size = [start], 0
        while stack:
            cur = stack.pop()
            size += 1
            for n in adj.get(cur, []):
                if n not in seen:
                    seen.add(n)
                    stack.append(n)
        biggest = max(biggest, size)
    print(groups, biggest)

main()
`},
	},
	{
		slug: "fewest-changes", title: "Fewest changes", difficulty: "hard", topic: "Graphs",
		statement: `A metro has lines; you can ride one as far as you like for free, and what
costs you is changing lines.

**Input.** A count ~n~, then ~n~ lines of ~lineName station1 station2 ...~. The
last line is ~start end~.

**Output.** The fewest line changes needed, ~-1~ if there is no route, and 0 if
you are already there.

    Input                 Output
    2                     1
    red A B C
    blue C D
    A D

The thing worth noticing: what you are searching over is lines, not stations.`,
		timeLimitMs: 4000,
		starters:    triLines("search over lines rather than stations"),
		tests: []store.Check{
			tokens("One line reaches both", "2\nred A B C\nblue C D\nA C\n", "0"),
			tokens("One change of line", "2\nred A B C\nblue C D\nA D\n", "1"),
			tokens("Already there", "1\nred A\nA A\n", "0"),
			tokens("No route at all", "2\nred A\nblue X\nA X\n", "-1"),
			hid2("An unknown station", "1\nred A\nA Z\n", "-1"),
			hid2("Two changes", "3\nr A B\ng B C\nb C D\nA D\n", "2"),
			hid2("The shorter of two routes wins", "3\nr A B\ng B Z\nx A Z\nA Z\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

#define MAXL 500
#define MAXS 60
static char lineName[MAXL][32];
static char station[MAXL][MAXS][32];
static int stops[MAXL];
static int lines;

static int onLine(int l, const char *s) {
    for (int i = 0; i < stops[l]; i++) if (strcmp(station[l][i], s) == 0) return 1;
    return 0;
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    getchar();
    for (int i = 0; i < n; i++) {
        char raw[2000];
        if (!fgets(raw, sizeof raw, stdin)) break;
        raw[strcspn(raw, "\n")] = '\0';
        char *tok = strtok(raw, " ");
        if (!tok) continue;
        strcpy(lineName[lines], tok);
        stops[lines] = 0;
        while ((tok = strtok(NULL, " ")) && stops[lines] < MAXS)
            strcpy(station[lines][stops[lines]++], tok);
        lines++;
    }
    char start[32], end[32];
    if (scanf("%31s %31s", start, end) != 2) return 1;
    if (strcmp(start, end) == 0) { printf("0\n"); return 0; }
    static int cost[MAXL], queue[MAXL];
    for (int i = 0; i < lines; i++) cost[i] = -1;
    int head = 0, tail = 0;
    for (int i = 0; i < lines; i++) if (onLine(i, start)) { cost[i] = 0; queue[tail++] = i; }
    while (head < tail) {
        int l = queue[head++];
        if (onLine(l, end)) { printf("%d\n", cost[l]); return 0; }
        for (int s = 0; s < stops[l]; s++)
            for (int m = 0; m < lines; m++)
                if (cost[m] < 0 && onLine(m, station[l][s])) { cost[m] = cost[l] + 1; queue[tail++] = m; }
    }
    printf("-1\n");
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
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	lines := map[string][]string{}
	at := map[string][]string{}
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		f := strings.Fields(strings.TrimRight(raw, "\r\n"))
		if len(f) == 0 {
			continue
		}
		lines[f[0]] = f[1:]
		for _, s := range f[1:] {
			at[s] = append(at[s], f[0])
		}
	}
	var start, end string
	fmt.Fscan(in, &start, &end)
	if start == end {
		fmt.Println(0)
		return
	}
	if len(at[start]) == 0 || len(at[end]) == 0 {
		fmt.Println(-1)
		return
	}
	type step struct {
		line string
		cost int
	}
	seen := map[string]bool{}
	var queue []step
	for _, l := range at[start] {
		seen[l] = true
		queue = append(queue, step{l, 0})
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, s := range lines[cur.line] {
			if s == end {
				fmt.Println(cur.cost)
				return
			}
			for _, nxt := range at[s] {
				if !seen[nxt] {
					seen[nxt] = true
					queue = append(queue, step{nxt, cur.cost + 1})
				}
			}
		}
	}
	fmt.Println(-1)
}
`,
			"python": `import sys
from collections import deque

def main():
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    lines, at = {}, {}
    for line in raw[1:1 + n]:
        f = line.split()
        if not f:
            continue
        lines[f[0]] = f[1:]
        for s in f[1:]:
            at.setdefault(s, []).append(f[0])
    start, end = raw[1 + n].split()
    if start == end:
        print(0)
        return
    if start not in at or end not in at:
        print(-1)
        return
    seen = set(at[start])
    q = deque((l, 0) for l in at[start])
    while q:
        l, cost = q.popleft()
        if end in lines[l]:
            print(cost)
            return
        for s in lines[l]:
            for nxt in at[s]:
                if nxt not in seen:
                    seen.add(nxt)
                    q.append((nxt, cost + 1))
    print(-1)

main()
`},
	},
	{
		slug: "count-routes", title: "Count routes", difficulty: "medium", topic: "Dynamic programming",
		statement: `A delivery van starts at the top-left of a grid and must reach the
bottom-right, moving only **right or down**, never entering a ~#~.

**Input.** ~rows cols~ on the first line, then ~rows~ lines of the grid.

**Output.** How many different routes there are. There may be none, which is 0.

    Input      Output
    2 2        2
    ..
    ..

Each square can only be reached from the one above it and the one to its left,
which is the whole of the answer.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    static char grid[1000][1005];
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    /* fill a table from the top-left corner outward */
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
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	// fill a table from the top-left corner outward
	fmt.Println(0)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    # fill a table from the top-left corner outward
    print(0)

main()
`),
		tests: []store.Check{
			tokens("Two ways across a small grid", "2 2\n..\n..\n", "2"),
			tokens("A wall closes one of them", "2 2\n..\n#.\n", "1"),
			tokens("A single blocked square", "1 1\n#\n", "0"),
			tokens("Nowhere to go but stay", "1 1\n.\n", "1"),
			hid2("Completely walled off", "3 2\n..\n##\n..\n", "0"),
			hid2("A wider grid has many routes", "3 3\n...\n...\n...\n", "6"),
			hid2("The target itself is blocked", "2 2\n..\n.#\n", "0"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>

static char grid[1000][1005];
static long long table[1000][1005];

int main(void) {
    int rows, cols;
    if (scanf("%d %d", &rows, &cols) != 2) return 1;
    for (int i = 0; i < rows; i++) if (scanf("%1004s", grid[i]) != 1) return 1;
    for (int r = 0; r < rows; r++) for (int c = 0; c < cols; c++) {
        if (grid[r][c] == '#') { table[r][c] = 0; continue; }
        if (r == 0 && c == 0) { table[r][c] = 1; continue; }
        table[r][c] = (r ? table[r - 1][c] : 0) + (c ? table[r][c - 1] : 0);
    }
    printf("%lld\n", table[rows - 1][cols - 1]);
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
	var rows, cols int
	fmt.Fscan(in, &rows, &cols)
	grid := make([]string, rows)
	for i := range grid {
		fmt.Fscan(in, &grid[i])
	}
	table := make([][]int, rows)
	for i := range table {
		table[i] = make([]int, cols)
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '#' {
				continue
			}
			if r == 0 && c == 0 {
				table[r][c] = 1
				continue
			}
			if r > 0 {
				table[r][c] += table[r-1][c]
			}
			if c > 0 {
				table[r][c] += table[r][c-1]
			}
		}
	}
	fmt.Println(table[rows-1][cols-1])
}
`,
			"python": `import sys

def main():
    data = sys.stdin.read().split()
    rows, cols = int(data[0]), int(data[1])
    grid = data[2:2 + rows]
    table = [[0] * cols for _ in range(rows)]
    for r in range(rows):
        for c in range(cols):
            if grid[r][c] == "#":
                continue
            if r == 0 and c == 0:
                table[r][c] = 1
            else:
                table[r][c] = (table[r - 1][c] if r else 0) + (table[r][c - 1] if c else 0)
    print(table[rows - 1][cols - 1])

main()
`},
	},
	{
		slug: "weighted-nesting", title: "Weighted nesting", difficulty: "medium", topic: "Recursion",
		statement: `A score sheet is a bracketed list that may contain numbers and other lists.
A number counts for its value times its depth: numbers at the top level count
once, numbers one list down count twice, and so on.

**Input.** One line, a bracketed list using ~[~ and ~]~ with numbers separated by
spaces. An empty sheet is ~[]~.

**Output.** The total.

    Input        Output
    [1 [2]]      5

That is 1x1 + 2x2.`,
		timeLimitMs: 4000,
		starters:    triLine("track the depth as you walk the brackets"),
		tests: []store.Check{
			tokens("Deeper numbers count for more", "[1 [2]]\n", "5"),
			tokens("Three levels down", "[[[3]]]\n", "9"),
			tokens("An empty sheet", "[]\n", "0"),
			tokens("A flat sheet is a plain sum", "[1 2 3]\n", "6"),
			hid2("Empty lists inside contribute nothing", "[[] [[]]]\n", "0"),
			hid2("Negative numbers weigh too", "[[-1]]\n", "-2"),
			hid2("Mixed depths in one list", "[1 [1 [1]]]\n", "6"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <stdlib.h>

int main(void) {
    static char line[100005];
    if (!fgets(line, sizeof line, stdin)) line[0] = '\0';
    line[strcspn(line, "\n")] = '\0';
    long long total = 0;
    int depth = 0;
    for (int i = 0; line[i]; ) {
        char c = line[i];
        if (c == '[') { depth++; i++; }
        else if (c == ']') { depth--; i++; }
        else if (c == '-' || isdigit((unsigned char)c)) {
            char *end;
            long long v = strtoll(line + i, &end, 10);
            total += v * depth;
            i = (int)(end - line);
        } else i++;
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
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	raw, _ := in.ReadString('\n')
	line := strings.TrimRight(raw, "\r\n")
	total, depth := 0, 0
	for i := 0; i < len(line); {
		switch c := line[i]; {
		case c == '[':
			depth++
			i++
		case c == ']':
			depth--
			i++
		case c == '-' || (c >= '0' && c <= '9'):
			j := i + 1
			for j < len(line) && line[j] >= '0' && line[j] <= '9' {
				j++
			}
			v, _ := strconv.Atoi(line[i:j])
			total += v * depth
			i = j
		default:
			i++
		}
	}
	fmt.Println(total)
}
`,
			"python": `import sys

def main():
    line = sys.stdin.readline().rstrip("\n")
    total, depth, i = 0, 0, 0
    while i < len(line):
        c = line[i]
        if c == "[":
            depth += 1
            i += 1
        elif c == "]":
            depth -= 1
            i += 1
        elif c == "-" or c.isdigit():
            j = i + 1
            while j < len(line) and line[j].isdigit():
                j += 1
            total += int(line[i:j]) * depth
            i = j
        else:
            i += 1
    print(total)

main()
`},
	},
	{
		slug: "org-depth", title: "Org depth", difficulty: "medium", topic: "Recursion",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~manager report1 report2 ...~ (a
manager with nobody reporting to them is written as just their name). The last
line names the root.

**Output.** How many levels deep the chart goes from that root, counting the root
as level 1. Somebody with nobody reporting to them is a chart of depth 1, and
somebody not in the chart at all is also 1. A chart where somebody reports,
however indirectly, to themselves is broken: print ~-1~.

    Input        Output
    2            2
    ceo vp
    vp
    ceo`,
		timeLimitMs: 4000,
		starters:    triLines("walk down each branch, watching for a loop"),
		tests: []store.Check{
			tokens("Two levels", "2\nceo vp\nvp\nceo\n", "2"),
			tokens("Nobody reports to them", "1\nsolo\nsolo\n", "1"),
			tokens("A cycle is broken", "2\na b\nb a\na\n", "-1"),
			tokens("Someone not in the chart at all", "0\nghost\n", "1"),
			hid2("The deepest branch decides", "4\na b c\nb\nc d\nd\na\n", "3"),
			hid2("Reporting to yourself is a cycle", "1\na a\na\n", "-1"),
			hid2("A cycle in a branch is still broken", "3\na b\nb c\nc b\na\n", "-1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

#define MAXN 1000
static char name[MAXN][32], reports[MAXN][1005];
static int total;
static int active[MAXN];
static int broken;

static int find(const char *w) {
    for (int i = 0; i < total; i++) if (strcmp(name[i], w) == 0) return i;
    return -1;
}

static int depth(const char *w) {
    int at = find(w);
    if (at < 0) return 1;
    if (active[at]) { broken = 1; return 1; }
    active[at] = 1;
    char copy[1005];
    strcpy(copy, reports[at]);
    int best = 0, start = 0;
    for (int i = 0; !broken; i++) {
        if (copy[i] != ' ' && copy[i] != '\0') continue;
        char save = copy[i];
        copy[i] = '\0';
        if (i > start) { int d = depth(copy + start); if (d > best) best = d; }
        start = i + 1;
        if (save == '\0') break;
    }
    active[at] = 0;
    return best + 1;
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
        if (sp) { *sp = '\0'; strcpy(reports[total], sp + 1); }
        else reports[total][0] = '\0';
        strcpy(name[total], line);
        total++;
    }
    char root[64];
    if (scanf("%63s", root) != 1) return 1;
    int d = depth(root);
    printf("%d\n", broken ? -1 : d);
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

var reports = map[string][]string{}
var active = map[string]bool{}
var broken bool

func depth(w string) int {
	kids, ok := reports[w]
	if !ok {
		return 1
	}
	if active[w] {
		broken = true
		return 1
	}
	active[w] = true
	best := 0
	for _, k := range kids {
		if d := depth(k); d > best {
			best = d
		}
	}
	delete(active, w)
	return best + 1
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	var n int
	fmt.Fscan(in, &n)
	in.ReadString('\n')
	for i := 0; i < n; i++ {
		raw, _ := in.ReadString('\n')
		f := strings.Fields(strings.TrimRight(raw, "\r\n"))
		if len(f) == 0 {
			continue
		}
		reports[f[0]] = f[1:]
	}
	var root string
	fmt.Fscan(in, &root)
	d := depth(root)
	if broken {
		fmt.Println(-1)
		return
	}
	fmt.Println(d)
}
`,
			"python": `import sys

def main():
    sys.setrecursionlimit(10000)
    raw = sys.stdin.read().split("\n")
    n = int(raw[0])
    reports = {}
    for line in raw[1:1 + n]:
        f = line.split()
        if f:
            reports[f[0]] = f[1:]
    root = raw[1 + n].strip()
    active = set()

    def depth(w):
        if w not in reports:
            return 1
        if w in active:
            return None
        active.add(w)
        best = 0
        for k in reports[w]:
            d = depth(k)
            if d is None:
                return None
            best = max(best, d)
        active.discard(w)
        return best + 1

    d = depth(root)
    print(-1 if d is None else d)

main()
`},
	},
	{
		slug: "two-teams", title: "Two teams", difficulty: "hard", topic: "Graphs",
		statement: `Some people refuse to be on the same team, and the refusal is mutual.

**Input.** ~p r~ on the first line, then ~p~ names, then ~r~ lines of ~a b~
meaning those two refuse each other.

**Output.** ~yes~ if everybody can be split into exactly two teams with no pair
of rivals on the same one, ~no~ otherwise. Teams may be any size, and one may be
empty.

    Input        Output
    3 3          no
    a b c
    a b
    b c
    c a

Three people who all refuse each other cannot be split in two, however you try —
and the shape of that failure is the thing to look for.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int p, r;
    if (scanf("%d %d", &p, &r) != 2) return 1;
    /* colour each group two ways and look for a clash */
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
	var p, r int
	fmt.Fscan(in, &p, &r)
	// colour each group two ways and look for a clash
	fmt.Println("yes")
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    p, r = int(data[0]), int(data[1])
    people = data[2:2 + p]
    # colour each group two ways and look for a clash
    print("yes")

main()
`),
		tests: []store.Check{
			tokens("A single rivalry splits fine", "2 1\na b\na b\n", "yes"),
			tokens("A triangle cannot be split", "3 3\na b c\na b\nb c\nc a\n", "no"),
			tokens("Nobody at all", "0 0\n", "yes"),
			tokens("Nobody objects to anybody", "2 0\na b\n", "yes"),
			hid2("A four-way ring splits", "4 4\na b c d\na b\nb c\nc d\nd a\n", "yes"),
			hid2("Separate groups are checked separately", "5 4\na b c d e\na b\nc d\nd e\ne c\n", "no"),
			hid2("A chain of any length splits", "3 2\na b c\na b\nb c\n", "yes"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

#define MAXP 1000
static char name[MAXP][32];
static int adj[MAXP][MAXP], deg[MAXP];
static int team[MAXP], total;

static int slot(const char *w) {
    for (int i = 0; i < total; i++) if (strcmp(name[i], w) == 0) return i;
    strcpy(name[total], w);
    team[total] = -1;
    return total++;
}

int main(void) {
    int p, r;
    if (scanf("%d %d", &p, &r) != 2) return 1;
    for (int i = 0; i < p; i++) { char w[32]; if (scanf("%31s", w) != 1) return 1; slot(w); }
    for (int i = 0; i < r; i++) {
        char a[32], b[32];
        if (scanf("%31s %31s", a, b) != 2) return 1;
        int x = slot(a), y = slot(b);
        adj[x][deg[x]++] = y;
        adj[y][deg[y]++] = x;
    }
    static int queue[MAXP];
    for (int s = 0; s < total; s++) {
        if (team[s] >= 0) continue;
        int head = 0, tail = 0;
        team[s] = 0; queue[tail++] = s;
        while (head < tail) {
            int cur = queue[head++];
            for (int k = 0; k < deg[cur]; k++) {
                int nb = adj[cur][k];
                if (team[nb] < 0) { team[nb] = 1 - team[cur]; queue[tail++] = nb; }
                else if (team[nb] == team[cur]) { printf("no\n"); return 0; }
            }
        }
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
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, r int
	fmt.Fscan(in, &p, &r)
	adj := map[string][]string{}
	people := make([]string, p)
	for i := range people {
		fmt.Fscan(in, &people[i])
		adj[people[i]] = nil
	}
	for i := 0; i < r; i++ {
		var a, b string
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	team := map[string]int{}
	for _, s := range people {
		if _, ok := team[s]; ok {
			continue
		}
		team[s] = 0
		queue := []string{s}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, nb := range adj[cur] {
				if t, ok := team[nb]; !ok {
					team[nb] = 1 - team[cur]
					queue = append(queue, nb)
				} else if t == team[cur] {
					fmt.Println("no")
					return
				}
			}
		}
	}
	fmt.Println("yes")
}
`,
			"python": `import sys
from collections import deque

def main():
    data = sys.stdin.read().split()
    p, r = int(data[0]), int(data[1])
    people = data[2:2 + p]
    adj = {name: [] for name in people}
    base = 2 + p
    for i in range(r):
        a, b = data[base + 2 * i], data[base + 1 + 2 * i]
        adj.setdefault(a, []).append(b)
        adj.setdefault(b, []).append(a)
    team = {}
    for s in people:
        if s in team:
            continue
        team[s] = 0
        q = deque([s])
        while q:
            cur = q.popleft()
            for nb in adj.get(cur, []):
                if nb not in team:
                    team[nb] = 1 - team[cur]
                    q.append(nb)
                elif team[nb] == team[cur]:
                    print("no")
                    return
    print("yes")

main()
`},
	},
	{
		slug: "cheapest-route", title: "Cheapest route", difficulty: "hard", topic: "Graphs",
		statement: `**Input.** A count ~n~, then ~n~ lines of ~a b cost~ — roads, drivable
either way, with positive costs. The last line is ~start end~.

**Output.** The lowest total cost from ~start~ to ~end~, ~-1~ if there is no way
through, and 0 if you are already there.

    Input        Output
    2            7
    a b 5
    b c 2
    a c

Taking the cheapest road out of each town in turn is not the same as the cheapest
route — what you want is to always extend from the town you can currently reach
most cheaply, wherever it is.`,
		timeLimitMs: 4000,
		starters: tri(`#include <stdio.h>
#include <string.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* always extend from the cheapest town reached so far */
    printf("%d\n", -1);
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
	// always extend from the cheapest town reached so far
	fmt.Println(-1)
}
`, `import sys

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    # always extend from the cheapest town reached so far
    print(-1)

main()
`),
		tests: []store.Check{
			tokens("Adds up the route", "2\na b 5\nb c 2\na c\n", "7"),
			tokens("Already there", "0\na a\n", "0"),
			tokens("No way through", "1\na b 1\na z\n", "-1"),
			tokens("Roads go both ways", "1\na b 4\nb a\n", "4"),
			hid2("The long way round can be cheaper", "3\na b 10\na c 1\nc b 1\na b\n", "2"),
			hid2("The greedy first step is not always right", "4\na b 1\nb z 50\na c 5\nc z 5\na z\n", "10"),
			hid2("A town with no roads at all", "1\na b 1\nz a\n", "-1"),
		},
		solutions: map[string]string{
			"c": `#include <stdio.h>
#include <string.h>

#define MAXT 2000
static char name[MAXT][32];
static int total;
static long long dist[MAXT];
static int done[MAXT];
static int to[MAXT][64];
static long long weight[MAXT][64];
static int deg[MAXT];

static int slot(const char *w) {
    for (int i = 0; i < total; i++) if (strcmp(name[i], w) == 0) return i;
    strcpy(name[total], w);
    return total++;
}

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    for (int i = 0; i < n; i++) {
        char a[32], b[32];
        long long cost;
        if (scanf("%31s %31s %lld", a, b, &cost) != 3) return 1;
        int x = slot(a), y = slot(b);
        to[x][deg[x]] = y; weight[x][deg[x]++] = cost;
        to[y][deg[y]] = x; weight[y][deg[y]++] = cost;
    }
    char start[32], end[32];
    if (scanf("%31s %31s", start, end) != 2) return 1;
    if (strcmp(start, end) == 0) { printf("0\n"); return 0; }
    int s = -1, e = -1;
    for (int i = 0; i < total; i++) {
        if (strcmp(name[i], start) == 0) s = i;
        if (strcmp(name[i], end) == 0) e = i;
    }
    if (s < 0 || e < 0) { printf("-1\n"); return 0; }
    for (int i = 0; i < total; i++) dist[i] = -1;
    dist[s] = 0;
    for (;;) {
        int cur = -1;
        for (int i = 0; i < total; i++)
            if (!done[i] && dist[i] >= 0 && (cur < 0 || dist[i] < dist[cur])) cur = i;
        if (cur < 0) break;
        done[cur] = 1;
        if (cur == e) { printf("%lld\n", dist[cur]); return 0; }
        for (int k = 0; k < deg[cur]; k++) {
            int nb = to[cur][k];
            long long step = dist[cur] + weight[cur][k];
            if (dist[nb] < 0 || step < dist[nb]) dist[nb] = step;
        }
    }
    printf("-1\n");
    return 0;
}
`,
			"go": `package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type item struct {
	town string
	cost int
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].cost < p[j].cost }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	last := old[len(old)-1]
	*p = old[:len(old)-1]
	return last
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	adj := map[string][]item{}
	for i := 0; i < n; i++ {
		var a, b string
		var cost int
		fmt.Fscan(in, &a, &b, &cost)
		adj[a] = append(adj[a], item{b, cost})
		adj[b] = append(adj[b], item{a, cost})
	}
	var start, end string
	fmt.Fscan(in, &start, &end)
	if start == end {
		fmt.Println(0)
		return
	}
	best := map[string]int{start: 0}
	q := &pq{{start, 0}}
	heap.Init(q)
	for q.Len() > 0 {
		cur := heap.Pop(q).(item)
		if cur.town == end {
			fmt.Println(cur.cost)
			return
		}
		if b, ok := best[cur.town]; ok && cur.cost > b {
			continue
		}
		for _, nb := range adj[cur.town] {
			step := cur.cost + nb.cost
			if b, ok := best[nb.town]; !ok || step < b {
				best[nb.town] = step
				heap.Push(q, item{nb.town, step})
			}
		}
	}
	fmt.Println(-1)
}
`,
			"python": `import sys
import heapq

def main():
    data = sys.stdin.read().split()
    n = int(data[0])
    adj = {}
    for i in range(n):
        a, b, cost = data[1 + 3 * i], data[2 + 3 * i], int(data[3 + 3 * i])
        adj.setdefault(a, []).append((b, cost))
        adj.setdefault(b, []).append((a, cost))
    start, end = data[1 + 3 * n], data[2 + 3 * n]
    if start == end:
        print(0)
        return
    best = {start: 0}
    heap = [(0, start)]
    while heap:
        cost, town = heapq.heappop(heap)
        if town == end:
            print(cost)
            return
        if cost > best.get(town, cost):
            continue
        for nxt, road in adj.get(town, []):
            step = cost + road
            if step < best.get(nxt, step + 1):
                best[nxt] = step
                heapq.heappush(heap, (step, nxt))
    print(-1)

main()
`},
	},
}

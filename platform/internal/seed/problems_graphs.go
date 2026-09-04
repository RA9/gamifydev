package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// graphProblems: grids, graphs and recursion — problems where the shape of the
// data is the problem, and the answer comes from walking it.
var graphProblems = []seedProblem{
	{
		slug: "flooded-region", title: "Flooded region", difficulty: "medium", topic: "Grids",
		statement: `A map is a list of equal-length strings, ~#~ for land and ~.~ for water.
Write ~region(grid, r, c)~ returning how many land squares are joined to the one
at row ~r~, column ~c~ — counting squares that touch edge to edge, not corner
to corner.

    region(["##.", "#..", "..#"], 0, 0)  ->  3

If the starting square is water, the region is empty: return 0.`,
		timeLimitMs: 5000,
		starters:    py("def region(grid, r, c):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Counts the joined land", `region(["##.", "#..", "..#"], 0, 0) == 3`),
			vis("Starting on water", `region(["##.", "#..", "..#"], 0, 2) == 0`),
			vis("A lone square", `region(["..#"], 0, 2) == 1`),
			hid("Corners do not join", `region(["#.", ".#"], 0, 0) == 1`),
			hid("The whole map is one region", `region(["##", "##"], 1, 1) == 4`),
			hid("A map with no land at all", `region(["..", ".."], 0, 0) == 0`),
		},
		solution: "def region(grid, r, c):\n    if grid[r][c] != '#':\n        return 0\n    rows, cols = len(grid), len(grid[0])\n    seen = {(r, c)}\n    stack = [(r, c)]\n    n = 0\n    while stack:\n        y, x = stack.pop()\n        n += 1\n        for dy, dx in ((1, 0), (-1, 0), (0, 1), (0, -1)):\n            ny, nx = y + dy, x + dx\n            if 0 <= ny < rows and 0 <= nx < cols and (ny, nx) not in seen and grid[ny][nx] == '#':\n                seen.add((ny, nx))\n                stack.append((ny, nx))\n    return n\n",
	},
	{
		slug: "count-clusters", title: "Count clusters", difficulty: "medium", topic: "Grids",
		statement: `On the same kind of map, write ~clusters(grid)~ returning how many separate
groups of land there are.

Here squares join **including diagonally** — corner to corner counts. That one
word changes the answer:

    clusters(["#.#", ".#.", "#.#"])  ->  1

Every square touches the middle one at a corner, so it is all one cluster.`,
		timeLimitMs: 5000,
		starters:    py("def clusters(grid):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Diagonals join it all up", `clusters(["#.#", ".#.", "#.#"]) == 1`),
			vis("Land too far apart to join", `clusters(["#..", "..#"]) == 2`),
			vis("No land at all", `clusters(["..", ".."]) == 0`),
			vis("One solid block", `clusters(["##", "##"]) == 1`),
			hid("Diagonal neighbours are one cluster", `clusters(["#.", ".#"]) == 1`),
			hid("Three separate clusters", `clusters(["#.#.#"]) == 3`),
		},
		solution: "def clusters(grid):\n    rows, cols = len(grid), len(grid[0])\n    seen = set()\n    n = 0\n    for r in range(rows):\n        for c in range(cols):\n            if grid[r][c] != '#' or (r, c) in seen:\n                continue\n            n += 1\n            seen.add((r, c))\n            stack = [(r, c)]\n            while stack:\n                y, x = stack.pop()\n                for dy in (-1, 0, 1):\n                    for dx in (-1, 0, 1):\n                        ny, nx = y + dy, x + dx\n                        if (0 <= ny < rows and 0 <= nx < cols\n                                and (ny, nx) not in seen and grid[ny][nx] == '#'):\n                            seen.add((ny, nx))\n                            stack.append((ny, nx))\n    return n\n",
	},
	{
		slug: "warehouse-route", title: "Warehouse route", difficulty: "hard", topic: "Graphs",
		statement: `A robot starts at the top-left of a grid and must reach the bottom-right.
It moves one square at a time, up, down, left or right, and cannot enter a
square marked ~#~.

Write ~steps(grid)~ returning the fewest moves, or ~-1~ if it cannot get there.
Standing on the target already is 0 moves.

    steps(["..", ".."])   ->  2
    steps(["#."])         ->  -1

Notice what "fewest" rules out: following one path to the end and hoping. You
need to reach every square that is one move away before any square that is two.`,
		timeLimitMs: 5000,
		starters:    py("def steps(grid):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Across a clear grid", `steps(["..", ".."]) == 2`),
			vis("Starting on a blocked square", `steps(["#."]) == -1`),
			vis("Already there", `steps(["."]) == 0`),
			vis("Walled off entirely", `steps(["..", "##", ".."]) == -1`),
			hid("Around an obstacle", `steps(["..#", "#..", "..."]) == 4`),
			hid("The target itself is blocked", `steps(["..", ".#"]) == -1`),
			hid("A long corridor", `steps(["." * 50]) == 49`),
		},
		solution: "from collections import deque\n\ndef steps(grid):\n    rows, cols = len(grid), len(grid[0])\n    if grid[0][0] == '#' or grid[rows - 1][cols - 1] == '#':\n        return -1\n    seen = {(0, 0)}\n    q = deque([(0, 0, 0)])\n    while q:\n        r, c, d = q.popleft()\n        if (r, c) == (rows - 1, cols - 1):\n            return d\n        for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):\n            nr, nc = r + dr, c + dc\n            if 0 <= nr < rows and 0 <= nc < cols and (nr, nc) not in seen and grid[nr][nc] != '#':\n                seen.add((nr, nc))\n                q.append((nr, nc, d + 1))\n    return -1\n",
	},
	{
		slug: "everything-needed", title: "Everything needed", difficulty: "medium", topic: "Graphs",
		statement: `A project's dependencies are a dictionary from a package to the list of
packages it needs directly. Write ~needed(deps, start)~ returning the sorted
names of everything that must be installed for ~start~ — following dependencies
all the way down.

    needed({"app": ["lib"], "lib": ["core"], "core": []}, "app")
      ->  ["core", "lib"]

The package itself is not in the list, unless something it needs depends back on
it. Dependencies can and do form cycles; walking one forever is not an answer.`,
		timeLimitMs: 5000,
		starters:    py("def needed(deps, start):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Follows the chain down", `needed({"app": ["lib"], "lib": ["core"], "core": []}, "app") == ["core", "lib"]`),
			vis("Nothing needed", `needed({"app": []}, "app") == []`),
			vis("An unknown package needs nothing", `needed({}, "ghost") == []`),
			hid("A cycle back to the start includes it", `needed({"a": ["b"], "b": ["a"]}, "a") == ["a", "b"]`),
			hid("A shared dependency is listed once", `needed({"a": ["b", "c"], "b": ["d"], "c": ["d"], "d": []}, "a") == ["b", "c", "d"]`),
			hid("A cycle further down still terminates", `needed({"a": ["b"], "b": ["c"], "c": ["b"]}, "a") == ["b", "c"]`),
		},
		solution: "def needed(deps, start):\n    out = set()\n    stack = list(deps.get(start, []))\n    while stack:\n        n = stack.pop()\n        if n in out:\n            continue\n        out.add(n)\n        stack.extend(deps.get(n, []))\n    return sorted(out)\n",
	},
	{
		slug: "install-order", title: "Install order", difficulty: "hard", topic: "Graphs",
		statement: `Given the same dependency dictionary, write ~install_order(deps)~ returning
an order to install every package in, such that nothing is installed before
something it needs.

Where several packages could go next, install the alphabetically first — so
there is exactly one right answer. If the dependencies form a cycle, no order
exists: return ~None~.

    install_order({"app": ["lib"], "lib": []})  ->  ["lib", "app"]

Packages that only ever appear as somebody else's dependency are still packages
and still get installed.`,
		timeLimitMs: 5000,
		starters:    py("def install_order(deps):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Dependencies come first", `install_order({"app": ["lib"], "lib": []}) == ["lib", "app"]`),
			vis("A cycle has no order", `install_order({"a": ["b"], "b": ["a"]}) is None`),
			vis("Nothing to install", "install_order({}) == []"),
			vis("Independent packages go alphabetically", `install_order({"b": [], "a": []}) == ["a", "b"]`),
			hid("A package only named as a dependency is installed", `install_order({"app": ["lib"]}) == ["lib", "app"]`),
			hid("Alphabetical among everything currently possible", `install_order({"z": [], "a": ["z"], "b": []}) == ["b", "z", "a"]`),
			hid("A package depending on itself is a cycle", `install_order({"a": ["a"]}) is None`),
		},
		solution: "import heapq\n\ndef install_order(deps):\n    names = set(deps)\n    for ds in deps.values():\n        names.update(ds)\n    indeg = {n: 0 for n in names}\n    unlocks = {n: [] for n in names}\n    for n, ds in deps.items():\n        for d in ds:\n            unlocks[d].append(n)\n            indeg[n] += 1\n    ready = [n for n in names if indeg[n] == 0]\n    heapq.heapify(ready)\n    out = []\n    while ready:\n        n = heapq.heappop(ready)\n        out.append(n)\n        for m in unlocks[n]:\n            indeg[m] -= 1\n            if indeg[m] == 0:\n                heapq.heappush(ready, m)\n    return out if len(out) == len(names) else None\n",
	},
	{
		slug: "friend-groups", title: "Friend groups", difficulty: "medium", topic: "Graphs",
		statement: `Write ~groups(people, friendships)~ returning a tuple: how many separate
groups of friends there are, and how many people are in the largest.

Friendship goes both ways, and a friend of a friend is in the same group.
Somebody with no friends is a group of one.

    groups(["a", "b", "c"], [("a", "b")])  ->  (2, 2)

Nobody at all is ~(0, 0)~.`,
		timeLimitMs: 5000,
		starters:    py("def groups(people, friendships):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A pair and a loner", `groups(["a", "b", "c"], [("a", "b")]) == (2, 2)`),
			vis("Nobody at all", "groups([], []) == (0, 0)"),
			vis("Everyone alone", `groups(["a", "b"], []) == (2, 1)`),
			vis("Everyone connected", `groups(["a", "b", "c"], [("a", "b"), ("b", "c")]) == (1, 3)`),
			hid("A friend of a friend joins the group", `groups(["a", "b", "c", "d"], [("a", "b"), ("c", "d"), ("b", "c")]) == (1, 4)`),
			hid("Two groups of two", `groups(["a", "b", "c", "d"], [("a", "b"), ("c", "d")]) == (2, 2)`),
		},
		solution: "def groups(people, friendships):\n    adj = {p: [] for p in people}\n    for a, b in friendships:\n        adj[a].append(b)\n        adj[b].append(a)\n    seen = set()\n    count = 0\n    biggest = 0\n    for p in people:\n        if p in seen:\n            continue\n        count += 1\n        seen.add(p)\n        stack = [p]\n        size = 0\n        while stack:\n            x = stack.pop()\n            size += 1\n            for y in adj[x]:\n                if y not in seen:\n                    seen.add(y)\n                    stack.append(y)\n        if size > biggest:\n            biggest = size\n    return (count, biggest)\n",
	},
	{
		slug: "fewest-changes", title: "Fewest changes", difficulty: "hard", topic: "Graphs",
		statement: `A metro is a dictionary from line name to the list of stations on it. You
can ride a line as far as you like for free; what costs you is changing lines.

Write ~changes(lines, start, end)~ returning the fewest line changes needed. One
line that reaches both stations is 0 changes. Return ~-1~ if there is no route,
and 0 if you are already there.

    changes({"red": ["A", "B", "C"], "blue": ["C", "D"]}, "A", "C")  ->  0
    changes({"red": ["A", "B", "C"], "blue": ["C", "D"]}, "A", "D")  ->  1

The thing worth noticing: what you are searching over is lines, not stations.`,
		timeLimitMs: 5000,
		starters:    py("def changes(lines, start, end):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("One line reaches both", `changes({"red": ["A", "B", "C"], "blue": ["C", "D"]}, "A", "C") == 0`),
			vis("One change of line", `changes({"red": ["A", "B", "C"], "blue": ["C", "D"]}, "A", "D") == 1`),
			vis("Already there", `changes({"red": ["A"]}, "A", "A") == 0`),
			vis("No route at all", `changes({"red": ["A"], "blue": ["X"]}, "A", "X") == -1`),
			hid("An unknown station", `changes({"red": ["A"]}, "A", "Z") == -1`),
			hid("Two changes", `changes({"r": ["A", "B"], "g": ["B", "C"], "b": ["C", "D"]}, "A", "D") == 2`),
			hid("The shorter of two routes wins", `changes({"r": ["A", "B"], "g": ["B", "Z"], "x": ["A", "Z"]}, "A", "Z") == 0`),
		},
		solution: "from collections import deque\n\ndef changes(lines, start, end):\n    if start == end:\n        return 0\n    at = {}\n    for name, stations in lines.items():\n        for s in stations:\n            at.setdefault(s, []).append(name)\n    if start not in at or end not in at:\n        return -1\n    seen = set(at[start])\n    q = deque((name, 0) for name in at[start])\n    while q:\n        name, cost = q.popleft()\n        if end in lines[name]:\n            return cost\n        for s in lines[name]:\n            for nxt in at[s]:\n                if nxt not in seen:\n                    seen.add(nxt)\n                    q.append((nxt, cost + 1))\n    return -1\n",
	},
	{
		slug: "count-routes", title: "Count routes", difficulty: "medium", topic: "Dynamic programming",
		statement: `A delivery van starts at the top-left of a grid and must reach the
bottom-right, moving only **right or down**. Squares marked ~#~ are blocked.

Write ~routes(grid)~ returning how many different routes there are.

    routes(["..", ".."])   ->  2
    routes(["..", "#."])   ->  1

There may be no route at all, which is 0. Each square can only be reached from
the one above it and the one to its left, which is the whole of the answer.`,
		timeLimitMs: 5000,
		starters:    py("def routes(grid):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Two ways across a small grid", `routes(["..", ".."]) == 2`),
			vis("A wall closes one of them", `routes(["..", "#."]) == 1`),
			vis("A single blocked square", `routes(["#"]) == 0`),
			vis("Nowhere to go but stay", `routes(["."]) == 1`),
			hid("Completely walled off", `routes(["..", "##", ".."]) == 0`),
			hid("A wider grid has many routes", `routes(["...", "...", "..."]) == 6`),
			hid("The target itself is blocked", `routes(["..", ".#"]) == 0`),
		},
		solution: "def routes(grid):\n    rows, cols = len(grid), len(grid[0])\n    table = [[0] * cols for _ in range(rows)]\n    for r in range(rows):\n        for c in range(cols):\n            if grid[r][c] == '#':\n                continue\n            if r == 0 and c == 0:\n                table[r][c] = 1\n            else:\n                table[r][c] = (table[r - 1][c] if r else 0) + (table[r][c - 1] if c else 0)\n    return table[rows - 1][cols - 1]\n",
	},
	{
		slug: "weighted-nesting", title: "Weighted nesting", difficulty: "medium", topic: "Recursion",
		statement: `A score sheet is a list that may contain numbers and other lists, nested
as deeply as you like. A number counts for its value times its depth: numbers at
the top level count once, numbers one list down count twice, and so on.

Write ~weighted(items)~ returning the total.

    weighted([1, [2]])   ->  5
    weighted([[[3]]])    ->  9

That first one is 1x1 + 2x2. An empty sheet scores 0.`,
		timeLimitMs: 5000,
		starters:    py("def weighted(items):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Deeper numbers count for more", "weighted([1, [2]]) == 5"),
			vis("Three levels down", "weighted([[[3]]]) == 9"),
			vis("An empty sheet", "weighted([]) == 0"),
			vis("A flat sheet is a plain sum", "weighted([1, 2, 3]) == 6"),
			hid("Empty lists inside contribute nothing", "weighted([[], [[]]]) == 0"),
			hid("Negative numbers weigh too", "weighted([[-1]]) == -2"),
			hid("Mixed depths in one list", "weighted([1, [1, [1]]]) == 6"),
		},
		solution: "def weighted(items, depth=1):\n    total = 0\n    for x in items:\n        if isinstance(x, list):\n            total += weighted(x, depth + 1)\n        else:\n            total += x * depth\n    return total\n",
	},
	{
		slug: "org-depth", title: "Org depth", difficulty: "medium", topic: "Recursion",
		statement: `An org chart maps each person to the list of people who report to them.
Write ~depth(reports, root)~ returning how many levels deep the chart goes
starting from ~root~, counting the root as level 1.

    depth({"ceo": ["vp"], "vp": []}, "ceo")  ->  2

Somebody with nobody reporting to them is a chart of depth 1. A chart where
somebody reports, however indirectly, to themselves is broken: return ~-1~.`,
		timeLimitMs: 5000,
		starters:    py("def depth(reports, root):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Two levels", `depth({"ceo": ["vp"], "vp": []}, "ceo") == 2`),
			vis("Nobody reports to them", `depth({"solo": []}, "solo") == 1`),
			vis("A cycle is broken", `depth({"a": ["b"], "b": ["a"]}, "a") == -1`),
			vis("Someone not in the chart at all", `depth({}, "ghost") == 1`),
			hid("The deepest branch decides", `depth({"a": ["b", "c"], "b": [], "c": ["d"], "d": []}, "a") == 3`),
			hid("Reporting to yourself is a cycle", `depth({"a": ["a"]}, "a") == -1`),
			hid("A cycle in a branch is still broken", `depth({"a": ["b"], "b": ["c"], "c": ["b"]}, "a") == -1`),
		},
		solution: "def depth(reports, root):\n    active = set()\n\n    def go(n):\n        if n in active:\n            return -1\n        active.add(n)\n        best = 0\n        for r in reports.get(n, []):\n            d = go(r)\n            if d == -1:\n                return -1\n            if d > best:\n                best = d\n        active.discard(n)\n        return best + 1\n\n    return go(root)\n",
	},
	{
		slug: "two-teams", title: "Two teams", difficulty: "hard", topic: "Graphs",
		statement: `Some people refuse to be on the same team. ~rivals~ maps each person to the
list of people they will not play with, and the refusal is mutual.

Write ~two_teams(rivals)~ returning whether everybody can be split into exactly
two teams with no pair of rivals on the same one. Teams may be any size, and one
may be empty.

    two_teams({"a": ["b"], "b": ["a"]})                        ->  True
    two_teams({"a": ["b"], "b": ["c"], "c": ["a"]})            ->  False

Three people who all refuse each other cannot be split in two, however you try
— and the shape of that failure is the thing to look for.`,
		timeLimitMs: 5000,
		starters:    py("def two_teams(rivals):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("A single rivalry splits fine", `two_teams({"a": ["b"], "b": ["a"]}) is True`),
			vis("A triangle cannot be split", `two_teams({"a": ["b"], "b": ["c"], "c": ["a"]}) is False`),
			vis("Nobody at all", "two_teams({}) is True"),
			vis("Nobody objects to anybody", `two_teams({"a": [], "b": []}) is True`),
			hid("A four-way ring splits", `two_teams({"a": ["b", "d"], "b": ["a", "c"], "c": ["b", "d"], "d": ["c", "a"]}) is True`),
			hid("Separate groups are checked separately", `two_teams({"a": ["b"], "b": ["a"], "c": ["d"], "d": ["e"], "e": ["c"]}) is False`),
			hid("A chain of any length splits", `two_teams({"a": ["b"], "b": ["a", "c"], "c": ["b"]}) is True`),
		},
		solution: "from collections import deque\n\ndef two_teams(rivals):\n    team = {}\n    for start in rivals:\n        if start in team:\n            continue\n        team[start] = 0\n        q = deque([start])\n        while q:\n            p = q.popleft()\n            for other in rivals.get(p, []):\n                if other not in team:\n                    team[other] = 1 - team[p]\n                    q.append(other)\n                elif team[other] == team[p]:\n                    return False\n    return True\n",
	},
	{
		slug: "cheapest-route", title: "Cheapest route", difficulty: "hard", topic: "Graphs",
		statement: `Roads are given as ~(a, b, cost)~ triples and can be driven either way.
Write ~cheapest(roads, start, end)~ returning the lowest total cost from
~start~ to ~end~, or ~-1~ if there is no way through. Being already there costs
0.

    cheapest([("a", "b", 5), ("b", "c", 2)], "a", "c")  ->  7

Costs are positive. Taking the cheapest road out of each town in turn is not the
same as the cheapest route — what you want is to always extend from the town you
can currently reach most cheaply, wherever it is.`,
		timeLimitMs: 5000,
		starters:    py("def cheapest(roads, start, end):\n    # your code here\n    pass\n"),
		tests: []store.Check{
			vis("Adds up the route", `cheapest([("a", "b", 5), ("b", "c", 2)], "a", "c") == 7`),
			vis("Already there", `cheapest([], "a", "a") == 0`),
			vis("No way through", `cheapest([("a", "b", 1)], "a", "z") == -1`),
			vis("Roads go both ways", `cheapest([("a", "b", 4)], "b", "a") == 4`),
			hid("The long way round can be cheaper", `cheapest([("a", "b", 10), ("a", "c", 1), ("c", "b", 1)], "a", "b") == 2`),
			hid("The greedy first step is not always right", `cheapest([("a", "b", 1), ("b", "z", 50), ("a", "c", 5), ("c", "z", 5)], "a", "z") == 10`),
			hid("A town with no roads at all", `cheapest([("a", "b", 1)], "z", "a") == -1`),
		},
		solution: "import heapq\n\ndef cheapest(roads, start, end):\n    adj = {}\n    for a, b, cost in roads:\n        adj.setdefault(a, []).append((b, cost))\n        adj.setdefault(b, []).append((a, cost))\n    best = {start: 0}\n    heap = [(0, start)]\n    while heap:\n        cost, town = heapq.heappop(heap)\n        if town == end:\n            return cost\n        if cost > best.get(town, cost):\n            continue\n        for nxt, road in adj.get(town, []):\n            step = cost + road\n            if step < best.get(nxt, step + 1):\n                best[nxt] = step\n                heapq.heappush(heap, (step, nxt))\n    return -1\n",
	},
}

# A-Star and Heuristic Search

Dijkstra's algorithm is correct and reasonably fast, but it has one obvious inefficiency: it explores outward in *every* direction equally. If you're routing from London to Edinburgh, Dijkstra will happily explore roads heading towards Cornwall first, because they're closer. It has no idea where the goal is.

**A\*** (pronounced "A-star") fixes exactly that. It adds a *guess* about how far each node still is from the goal, and uses that guess to lean the search in the right direction. It's the standard algorithm for game pathfinding and route planning, and — as you'll see — Dijkstra turns out to be a special case of it.

## Dijkstra's blind spot

Picture a grid where you're pathfinding from the left edge to the right edge. Dijkstra expands nodes in order of distance from the start, which draws an expanding circle:

```text
Dijkstra:  expands in all directions      A*: leans toward the goal

    . . o o o . . . . .                    . . . . . . . . . .
    . o o o o o . . . .                    . . o o o . . . . .
    o o o S o o o . . G                    o o o S o o o o o G
    . o o o o o . . . .                    . . o o o . . . . .
    . . o o o . . . . .                    . . . . . . . . . .

    o = node expanded before reaching G
```

Both find the same optimal path. A* just does far less work getting there, because it stops wasting effort on nodes that are obviously heading away from the goal.

:::key
Dijkstra knows how far you've come. A* also estimates how far you have left to go — and that one extra piece of information is what lets it aim.
:::

## g, h, and f

A* tracks two numbers for every node `n`:

- **g(n)** — the *actual* cost of the best known path from the start to `n`. This is exactly Dijkstra's `dist`. It's a fact.
- **h(n)** — the **heuristic**: an *estimate* of the remaining cost from `n` to the goal. It's a guess, computed without exploring anything.

And it orders its priority queue by their sum:

**f(n) = g(n) + h(n)** — the estimated total cost of a route that goes through `n`.

Dijkstra pops the node with the smallest `g`. A* pops the node with the smallest `f`. That is the entire difference between the two algorithms.

```text
node   g (cost so far)   h (estimate left)   f = g + h
-----------------------------------------------------
 X            5                  12              17
 Y            9                   3              12    <- A* expands Y first
                                                          (Dijkstra would pick X)
```

Node X is closer to the start, so Dijkstra grabs it. A* looks at the totals and sees that Y is on a much more promising route.

:::analogy
Planning a drive, `g` is the miles already on your odometer and `h` is the straight-line distance still shown on the map. You don't choose the next town by which is nearest to home — you choose by which gives the shortest *total* trip. That's `f`.
:::

## Admissible and consistent heuristics

A* only guarantees the optimal path if the heuristic behaves. Two properties matter.

**Admissible** — h(n) **never overestimates** the true remaining cost. It may be too low, or exactly right, but never too high.

Admissibility is what guarantees optimality, and the reason is worth spelling out. Suppose A* is about to pop the goal with total cost f = C. Any other node `n` still waiting in the queue has f(n) = g(n) + h(n) ≤ g(n) + (true remaining cost) = the true cost of the best route through `n`. Since f(n) ≥ C (or it would have been popped first), every unexplored route costs at least C. So the path A* found is optimal.

Break admissibility and that argument dies. An overestimating heuristic can make a genuinely good route *look* expensive, so A* never expands it and returns a worse path — quickly, and with no warning.

**Consistent** (or monotone) — a stronger condition: for every edge from `n` to `m` with cost w, `h(n) ≤ w + h(m)`. In words, the estimate never drops by more than the cost of the step you took. Consistency implies admissibility, and it gives you something practical: with a consistent heuristic, a node's `g` is final the first time you pop it, so you never need to re-expand a node — exactly Dijkstra's guarantee.

:::warning
An **inadmissible** heuristic (one that overestimates) makes A* fast but no longer optimal. That's sometimes a deliberate choice in games — "good enough, computed instantly" — but you must know you're making it.
:::

## Heuristics on a grid

On a grid, the natural heuristics are geometric distances that ignore all obstacles. Ignoring obstacles is precisely what makes them admissible: the real path can only be longer than the unobstructed one, never shorter.

**Manhattan distance** — for a grid where you can only move up, down, left and right:

```python
def manhattan(a, b):
    return abs(a[0] - b[0]) + abs(a[1] - b[1])

print(manhattan((1, 2), (4, 6)))   # 7
```

Three steps across plus four steps down is 7 moves minimum, walls or no walls. Never an overestimate — admissible.

**Euclidean distance** — for movement in any direction:

```python
import math

def euclidean(a, b):
    return math.hypot(a[0] - b[0], a[1] - b[1])

print(round(euclidean((1, 2), (4, 6)), 2))   # 5.0
```

The straight line is the shortest possible route between two points, so it can never overestimate either.

:::warning
Using **Euclidean** distance on a grid that only allows four-way movement is admissible but weak — it always underestimates, so A* explores more than it needs to. Using **Manhattan** distance on a grid that allows diagonal movement is a real bug: diagonals make the true cost lower than Manhattan predicts, so the heuristic overestimates and optimality is lost.
:::

## A* in code

The implementation is Dijkstra with `f` in the heap instead of `g`.

```python
import heapq

def a_star(neighbours, start, goal, h):
    """neighbours(n) -> list of (node, cost);  h(n) -> estimate to goal."""
    g = {start: 0}
    parent = {start: None}
    heap = [(h(start), start)]                     # (f, node)
    while heap:
        _, current = heapq.heappop(heap)
        if current == goal:
            path = [current]                       # rebuild by walking parents
            while parent[path[-1]] is not None:
                path.append(parent[path[-1]])
            return list(reversed(path))
        for nxt, cost in neighbours(current):
            tentative = g[current] + cost          # relaxation, as before
            if nxt not in g or tentative < g[nxt]:
                g[nxt] = tentative
                parent[nxt] = current
                heapq.heappush(heap, (tentative + h(nxt), nxt))
    return None                                    # goal unreachable
```

And a small runnable example on a 4×4 grid with one wall:

```python
walls = {(1, 1), (1, 2)}
def grid_neighbours(cell):
    r, c = cell
    out = []
    for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
        n = (r + dr, c + dc)
        if 0 <= n[0] < 4 and 0 <= n[1] < 4 and n not in walls:
            out.append((n, 1))
    return out

path = a_star(grid_neighbours, (0, 0), (3, 3), lambda n: manhattan(n, (3, 3)))
print(len(path))   # 7  (6 moves: the Manhattan distance, no detour needed)
```

**Complexity.** In the worst case A* degenerates to Dijkstra and is O((V + E) log V) — that happens when the heuristic gives no useful information. In practice a good heuristic can cut the number of expanded nodes enormously. Big O doesn't capture that improvement at all, because the heuristic changes the constant and the explored fraction, not the growth class.

## The whole family, in one line

Change the ordering rule and you get three different algorithms:

```text
order by  g        ->  Dijkstra          optimal, explores everywhere
order by  g + h    ->  A*                optimal (if h is admissible), aimed
order by  h        ->  Greedy best-first fast, often NOT optimal
```

**A\* with h(n) = 0 is exactly Dijkstra.** A zero heuristic is trivially admissible — it never overestimates anything — so A* stays optimal, but with nothing to aim at, `f = g` and you're back to expanding in circles. Dijkstra is A* that has given up guessing.

At the other extreme, **greedy best-first search** orders by `h` alone, ignoring how much you've already spent. It charges straight at the goal and is often blisteringly fast — but it will happily commit to a route that heads towards the goal and then hits a long wall, because it never weighs the cost already sunk. It is not optimal.

:::key
`f = g + h` is a balance. All `g` and you're thorough but blind (Dijkstra). All `h` and you're fast but reckless (greedy best-first). A* keeps both terms, which is why it's both aimed and optimal.
:::

## Check Your Understanding

:::quiz
Q: What does the heuristic h(n) in A* represent?
- The exact distance already travelled from the start
- An estimate of the remaining cost from n to the goal *
- The number of nodes expanded so far
- The weight of the cheapest edge leaving n
E: `g` is the actual cost so far and `h` is the estimated cost still to go. A* orders its queue by their sum, f = g + h.
:::

:::quiz
Q: What must be true of a heuristic for A* to be guaranteed to find the optimal path?
- It must be exactly correct
- It must never underestimate the remaining cost
- It must never overestimate the remaining cost *
- It must always return zero
E: That property is called admissibility. Underestimating is safe (it just means more exploring); overestimating can hide the genuinely best route and break optimality.
:::

:::match
Q: Match each ordering rule to the algorithm it produces.
- Order by g only | Dijkstra's algorithm
- Order by g + h | A* search
- Order by h only | Greedy best-first search
E: All three are the same loop with a different priority. Dropping h loses the aim; dropping g loses the optimality guarantee.
:::

## Talk about it

> An admissible heuristic must never overestimate — being too optimistic is safe, being too pessimistic is not. Explain in your own words why underestimating only costs you extra searching, while overestimating can make A* miss the best route entirely.

## What's next

You've now spent three lessons finding the cheapest way from one point to another. Next up is a different question about a graph entirely: **Minimum Spanning Trees: Prim and Kruskal** — how do you connect *every* node using the least total edge weight? It's the problem of laying cable to every house in a town, and the two classic answers are both greedy, both provably correct, and built from tools you already have.

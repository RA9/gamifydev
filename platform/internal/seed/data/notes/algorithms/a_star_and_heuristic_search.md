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

```c
#include <stdio.h>
#include <stdlib.h>

typedef struct { int row, col; } Cell;

int manhattan(Cell a, Cell b) {
    return abs(a.row - b.row) + abs(a.col - b.col);
}

int main(void) {
    printf("%d\n", manhattan((Cell){1, 2}, (Cell){4, 6})); /* 7 */
    return 0;
}
```

Three steps across plus four steps down is 7 moves minimum, walls or no walls. Never an overestimate — admissible.

**Euclidean distance** — for movement in any direction:

```c
#include <math.h>
#include <stdio.h>

double euclidean(Cell a, Cell b) {
    return hypot((double)a.row - b.row, (double)a.col - b.col);
}

int main(void) {
    printf("%.2f\n", euclidean((Cell){1, 2}, (Cell){4, 6})); /* 5.00 */
    return 0;
}
```

The straight line is the shortest possible route between two points, so it can never overestimate either.

:::warning
Using **Euclidean** distance on a grid that only allows four-way movement is admissible but weak — it always underestimates, so A* explores more than it needs to. Using **Manhattan** distance on a grid that allows diagonal movement is a real bug: diagonals make the true cost lower than Manhattan predicts, so the heuristic overestimates and optimality is lost.
:::

## A* in code

The implementation is Dijkstra with `f` in the heap instead of `g`.

```c
#include <limits.h>
#include <stdbool.h>
#include <stddef.h>

typedef struct { size_t node; int cost; } Edge;
typedef struct { size_t node; int priority; } HeapEntry;
typedef size_t (*NeighbourFn)(size_t node, Edge out[], size_t capacity);
typedef int (*HeuristicFn)(size_t node);

static bool heap_push(HeapEntry heap[], size_t *size, size_t capacity, HeapEntry item) {
    if (*size == capacity) return false;
    size_t i = (*size)++;
    while (i > 0) {
        size_t parent = (i - 1) / 2;
        if (heap[parent].priority <= item.priority) break;
        heap[i] = heap[parent];
        i = parent;
    }
    heap[i] = item;
    return true;
}

static HeapEntry heap_pop(HeapEntry heap[], size_t *size) {
    HeapEntry result = heap[0], last = heap[--*size];
    size_t i = 0;
    while (2 * i + 1 < *size) {
        size_t child = 2 * i + 1;
        if (child + 1 < *size && heap[child + 1].priority < heap[child].priority) ++child;
        if (last.priority <= heap[child].priority) break;
        heap[i] = heap[child];
        i = child;
    }
    if (*size > 0) heap[i] = last;
    return result;
}

bool a_star(size_t node_count, size_t start, size_t goal, NeighbourFn neighbours,
            HeuristicFn h, size_t path[], size_t path_capacity, size_t *path_length) {
    int g[64]; size_t parent[64];
    HeapEntry heap[256]; size_t heap_size = 0;
    if (node_count > 64) return false;
    for (size_t i = 0; i < node_count; ++i) g[i] = INT_MAX;
    g[start] = 0; parent[start] = start;
    if (!heap_push(heap, &heap_size, 256, (HeapEntry){start, h(start)})) return false;

    while (heap_size > 0) {
        HeapEntry entry = heap_pop(heap, &heap_size);
        size_t current = entry.node;
        if (entry.priority != g[current] + h(current)) continue; /* Stale entry. */
        if (current == goal) {
            size_t length = 0;
            for (size_t node = goal;; node = parent[node]) {
                if (length == path_capacity) return false;
                path[length++] = node;
                if (node == start) break;
            }
            for (size_t i = 0; i < length / 2; ++i) {
                size_t tmp = path[i]; path[i] = path[length - 1 - i]; path[length - 1 - i] = tmp;
            }
            *path_length = length;
            return true;
        }
        Edge edges[8];
        size_t count = neighbours(current, edges, 8);
        for (size_t i = 0; i < count; ++i) {
            size_t next = edges[i].node;
            if (next < node_count && edges[i].cost >= 0 &&
                g[current] <= INT_MAX - edges[i].cost) {
                int tentative = g[current] + edges[i].cost;
                if (tentative < g[next]) {
                    g[next] = tentative; parent[next] = current;
                    if (!heap_push(heap, &heap_size, 256, (HeapEntry){next, tentative + h(next)})) return false;
                }
            }
        }
    }
    return false;
}
```

And a small runnable example on a 4×4 grid with one wall:

```c
#include <stddef.h>
#include <stdio.h>

static const bool walls[4][4] = {
    {false, false, false, false},
    {false, true,  true,  false},
    {false, false, false, false},
    {false, false, false, false}
};

size_t grid_neighbours(size_t node, Edge out[], size_t capacity) {
    static const int directions[4][2] = {{1,0}, {-1,0}, {0,1}, {0,-1}};
    int row = (int)(node / 4), col = (int)(node % 4); size_t count = 0;
    for (size_t i = 0; i < 4; ++i) {
        int r = row + directions[i][0], c = col + directions[i][1];
        if (r >= 0 && r < 4 && c >= 0 && c < 4 && !walls[r][c] && count < capacity) {
            out[count++] = (Edge){(size_t)(r * 4 + c), 1};
        }
    }
    return count;
}

int grid_heuristic(size_t node) {
    Cell cell = {(int)(node / 4), (int)(node % 4)};
    return manhattan(cell, (Cell){3, 3});
}

int main(void) {
    size_t path[16], length = 0;
    if (!a_star(16, 0, 15, grid_neighbours, grid_heuristic, path, 16, &length)) return 1;
    printf("%zu\n", length); /* 7 nodes, or 6 moves. */
    return 0;
}
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

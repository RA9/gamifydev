# Minimum Spanning Trees: Prim and Kruskal

Shortest-path algorithms answer "what's the cheapest way from A to B?" This lesson asks something different: **what's the cheapest way to connect everything?**

Imagine you're laying fibre to every house in a village. You don't need a direct cable between every pair of houses — you just need the network to be joined up, so that any house can reach any other. And you want to use as little cable as possible. That's the **minimum spanning tree** problem, and it has two beautiful greedy solutions.

## Spanning trees

A **spanning tree** of a connected graph is a subset of its edges that:

- **touches every vertex** (that's the "spanning" part), and
- **contains no cycles** (that's the "tree" part).

Those two conditions together force a specific edge count. To connect V vertices you need at least V-1 edges, and any more than V-1 must create a cycle. So **every spanning tree of a V-vertex graph has exactly V-1 edges** — no exceptions.

A graph usually has many spanning trees. A **minimum spanning tree** (MST) is one whose edge weights sum to the smallest possible total.

```text
the graph            A --1-- B
                     | \     |
                     4   2   5
                     |     \ |
                     C --3-- D

all five edges:   A-B 1   A-D 2   C-D 3   A-C 4   B-D 5

one spanning tree:   A-B (1), A-C (4), B-D (5)   total = 10
another:             A-B (1), A-D (2), C-D (3)   total =  6   <- minimum
```

Both of those choices use 3 edges for 4 vertices, both connect everything, and neither contains a cycle. Only one of them is minimum.

:::key
A spanning tree connects all V vertices with exactly **V-1 edges** and no cycles. The *minimum* spanning tree is the one whose total weight is smallest.
:::

## The cycle rule that makes greedy work

Both algorithms in this lesson are **greedy**: they repeatedly grab the cheapest edge that doesn't break the rules, and never reconsider. For most problems that's a recipe for a wrong answer. Here it's provably optimal, thanks to the **cut property**.

Split the vertices into any two groups. The cheapest edge crossing between those groups is always safe to include in *some* minimum spanning tree — because any spanning tree must cross that divide somewhere, and if it crosses using a more expensive edge, you could swap in the cheap one and get a lighter tree. Prim and Kruskal are just two different ways of choosing which cut to look at, and that shared foundation is why both are correct.

:::analogy
Think of the cut property as a rule about bridges. However you divide a country into two halves, the cheapest bridge across that divide is one you'd always want to build. Any plan that connects the halves with a pricier bridge instead can be improved by swapping.
:::

## Prim's algorithm: grow one tree

Prim's starts from a single vertex and grows one connected blob outwards, at each step adding the cheapest edge that connects a vertex *inside* the tree to a vertex *outside* it. Finding that cheapest crossing edge is exactly a job for a **priority queue** — the same min-heap Dijkstra used.

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

enum { A, B, C, D, VERTICES };
typedef struct { size_t from, to; int weight; } Edge;
typedef struct { const Edge *items; size_t count; } Adjacency;

static const Edge from_a[] = {{A,B,1}, {A,C,4}, {A,D,2}};
static const Edge from_b[] = {{B,A,1}, {B,D,5}};
static const Edge from_c[] = {{C,A,4}, {C,D,3}};
static const Edge from_d[] = {{D,A,2}, {D,B,5}, {D,C,3}};
static const Adjacency graph[] = {{from_a,3}, {from_b,2}, {from_c,2}, {from_d,3}};

static void push(Edge heap[], size_t *n, Edge edge) {
    size_t i = (*n)++;
    while (i > 0) {
        size_t parent = (i - 1) / 2;
        if (heap[parent].weight <= edge.weight) break;
        heap[i] = heap[parent]; i = parent;
    }
    heap[i] = edge;
}

static Edge pop(Edge heap[], size_t *n) {
    Edge result = heap[0], last = heap[--*n]; size_t i = 0;
    while (2 * i + 1 < *n) {
        size_t child = 2 * i + 1;
        if (child + 1 < *n && heap[child + 1].weight < heap[child].weight) ++child;
        if (last.weight <= heap[child].weight) break;
        heap[i] = heap[child]; i = child;
    }
    if (*n > 0) heap[i] = last;
    return result;
}

size_t prim(const Adjacency graph[], size_t start, Edge mst[], int *total) {
    bool in_tree[VERTICES] = {false}; Edge heap[16];
    size_t heap_size = 0, count = 0;
    *total = 0; in_tree[start] = true;
    for (size_t i = 0; i < graph[start].count; ++i) push(heap, &heap_size, graph[start].items[i]);
    while (heap_size > 0 && count + 1 < VERTICES) {
        Edge edge = pop(heap, &heap_size);
        if (in_tree[edge.to]) continue;
        in_tree[edge.to] = true; mst[count++] = edge; *total += edge.weight;
        for (size_t i = 0; i < graph[edge.to].count; ++i) {
            Edge next = graph[edge.to].items[i];
            if (!in_tree[next.to]) push(heap, &heap_size, next);
        }
    }
    return count;
}

int main(void) {
    Edge mst[VERTICES - 1]; int total;
    size_t count = prim(graph, A, mst, &total);
    for (size_t i = 0; i < count; ++i) printf("%c-%c %d\n", (char)('A' + mst[i].from), (char)('A' + mst[i].to), mst[i].weight);
    printf("total = %d\n", total); /* 6 */
    return 0;
}
```

Trace it:

```text
tree = {A}          candidate edges: A-B 1, A-D 2, A-C 4
  take A-B (1)      tree = {A, B}     add B's edges: B-D 5
  take A-D (2)      tree = {A, B, D}  add D's edges: D-C 3
  take D-C (3)      tree = {A, B, C, D}   all vertices in — stop
total weight 6, and 3 edges for 4 vertices
```

**Complexity: O(E log V)** with a binary heap. Every edge can be pushed once, and each push and pop costs O(log V). Prim's tends to be the better choice on **dense** graphs, where E is close to V².

## Union-find, in brief

Kruskal's algorithm needs to answer one question fast: *would adding this edge create a cycle?* That's the same as asking whether its two endpoints are already connected.

**Union-find** (also called disjoint-set) is the data structure for exactly that. It keeps a collection of disjoint groups and supports two operations:

- **find(x)** — which group is x in? Returns a representative element.
- **union(x, y)** — merge x's group and y's group into one.

Adding an edge creates a cycle precisely when both endpoints already `find` to the same representative.

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct { size_t *parent, *rank; size_t count; } UnionFind;

bool uf_init(UnionFind *uf, size_t count) {
    *uf = (UnionFind){.count = count};
    if (count > SIZE_MAX / sizeof(size_t)) return false;
    uf->parent = malloc(count * sizeof *uf->parent);
    uf->rank = calloc(count, sizeof *uf->rank);
    if ((uf->parent == NULL || uf->rank == NULL) && count != 0) {
        free(uf->parent); free(uf->rank); *uf = (UnionFind){0}; return false;
    }
    for (size_t i = 0; i < count; ++i) uf->parent[i] = i;
    return true;
}

void uf_destroy(UnionFind *uf) {
    free(uf->parent); free(uf->rank); *uf = (UnionFind){0};
}

size_t uf_find(UnionFind *uf, size_t x) {
    while (uf->parent[x] != x) {
        uf->parent[x] = uf->parent[uf->parent[x]]; /* Path compression. */
        x = uf->parent[x];
    }
    return x;
}

bool uf_union(UnionFind *uf, size_t x, size_t y) {
    size_t rx = uf_find(uf, x), ry = uf_find(uf, y);
    if (rx == ry) return false;
    if (uf->rank[rx] < uf->rank[ry]) uf->parent[rx] = ry;
    else {
        uf->parent[ry] = rx;
        if (uf->rank[rx] == uf->rank[ry]) ++uf->rank[rx];
    }
    return true;
}

/* After uf_init(&uf, 3), uf_union(&uf, A, B) is true once, then false. */
```

Each group is stored as a tree of parent pointers, and the path-compression line flattens those trees as it walks them, pointing nodes straight at their grandparents. With that optimisation, find and union are effectively constant time — technically O(α(n)), where α is the inverse Ackermann function, a value below 5 for any n you will ever encounter.

:::tip
Union-find is worth knowing well beyond Kruskal's. It's the standard tool for "are these two things in the same group?" — connected components, network connectivity, and grouping in clustering algorithms all use it.
:::

## Kruskal's algorithm: sort the edges

Kruskal's takes a global view. Sort every edge by weight, then walk the sorted list taking each edge unless it would create a cycle. Unlike Prim's, the partial result is a *forest* of separate pieces that gradually merge into one tree.

```c
#include <stddef.h>
#include <stdlib.h>

static int edge_by_weight(const void *a, const void *b) {
    const Edge *x = a, *y = b;
    return (x->weight > y->weight) - (x->weight < y->weight);
}

bool kruskal(size_t vertex_count, Edge edges[], size_t edge_count,
             Edge mst[], size_t *mst_count, int *total) {
    UnionFind uf;
    if (!uf_init(&uf, vertex_count)) return false;
    qsort(edges, edge_count, sizeof *edges, edge_by_weight);
    *mst_count = 0; *total = 0;
    for (size_t i = 0; i < edge_count && *mst_count + 1 < vertex_count; ++i) {
        if (uf_union(&uf, edges[i].from, edges[i].to)) {
            mst[(*mst_count)++] = edges[i];
            *total += edges[i].weight;
        }
    }
    bool connected = vertex_count == 0 || *mst_count + 1 == vertex_count;
    uf_destroy(&uf);
    return connected;
}

Edge edges[] = {{A,B,1}, {A,D,2}, {C,D,3}, {A,C,4}, {B,D,5}};
Edge mst[VERTICES - 1]; size_t mst_count; int total;
/* kruskal(VERTICES, edges, 5, mst, &mst_count, &total) produces total 6. */
```

Trace it:

```text
sorted edges:  A-B 1,  A-D 2,  C-D 3,  A-C 4,  B-D 5

A-B 1   different groups -> take.   groups: {A,B} {C} {D}
A-D 2   different groups -> take.   groups: {A,B,D} {C}
C-D 3   different groups -> take.   groups: {A,B,C,D}
        3 edges = V-1 -> stop.  (A-C and B-D would both make cycles)

total weight 6 — the same tree Prim's found
```

Both algorithms found the same MST here. They usually do, and when a graph has ties they may find different trees of equal total weight — both are still correct.

**Complexity: O(E log E)**, dominated by sorting the edges. Since E ≤ V², log E is at most 2·log V, so this is often written O(E log V) — the same class as Prim's. Kruskal's is a natural fit for **sparse** graphs, where there aren't many edges to sort.

:::warning
Don't confuse a minimum spanning tree with shortest paths. The MST minimises the *total* weight of the whole network; it does **not** guarantee the cheapest route between any particular pair of nodes. A path through the MST can be far worse than Dijkstra's answer.
:::

## Where MSTs are used

**Network and cable layout.** The original motivation: connect every building, substation or exchange with the least total cable, pipe or road.

**Clustering.** Build the MST of your data points (weights are distances), then delete the k-1 most expensive edges. What falls apart is k clusters whose members are close to each other — single-linkage clustering, which is an MST in disguise.

**Approximation for harder problems.** The travelling salesman problem — visit every city once and return home, as cheaply as possible — has no known efficient exact solution, but an MST gives a lower bound and can be walked to build a tour provably within a factor of the optimum.

**Image segmentation and maze generation.** Treat pixels or grid cells as vertices with weighted edges, and the MST carves out regions or corridors.

:::example
To generate a random maze: build a grid graph, assign every wall a random weight, and compute the MST. The result touches every cell and contains no loops — which is exactly the definition of a perfect maze with a unique path between any two points.
:::

## Check Your Understanding

:::quiz
Q: How many edges does a spanning tree of a graph with 10 vertices have?
- 9 *
- 10
- 11
- It depends on the number of edges in the graph
E: A spanning tree always has exactly V-1 edges. Fewer and it can't connect everything; more and it must contain a cycle.
:::

:::quiz
Q: What does Kruskal's algorithm use union-find for?
- To sort the edges by weight
- To check whether adding an edge would create a cycle *
- To find the shortest path between two vertices
- To count the vertices
E: An edge creates a cycle exactly when both endpoints are already in the same group, which is one `find` call on each end.
:::

## Talk about it

> Both Prim's and Kruskal's take the cheapest available option at every step and never go back to reconsider — and for this problem, that's provably optimal. Think of a decision where "always take the cheapest option right now" would clearly lead you somewhere bad. What's different about that situation?

## What's next

You've now met four greedy algorithms — Dijkstra, A*, Prim and Kruskal — all of which take the best local option and never look back, and all of which are provably correct. Next up is **Greedy Algorithms**, where we finally examine that strategy directly: when does greed actually work, and what does failure look like? You'll see a tiny coin-change example where greedy confidently returns the wrong answer.

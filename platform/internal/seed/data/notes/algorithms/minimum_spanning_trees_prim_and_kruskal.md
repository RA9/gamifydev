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

```python
import heapq

def prim(graph, start):
    """graph: {node: [(neighbour, weight), ...]} — undirected, both directions."""
    in_tree = {start}
    edges = []                                     # (weight, from, to)
    for v, w in graph[start]:
        heapq.heappush(edges, (w, start, v))
    mst, total = [], 0
    while edges and len(in_tree) < len(graph):
        w, u, v = heapq.heappop(edges)             # cheapest crossing edge
        if v in in_tree:
            continue                               # both ends inside: a cycle
        in_tree.add(v)
        mst.append((u, v, w))
        total += w
        for nxt, w2 in graph[v]:                   # v's edges may now cross
            if nxt not in in_tree:
                heapq.heappush(edges, (w2, v, nxt))
    return mst, total

graph = {
    "A": [("B", 1), ("C", 4), ("D", 2)],
    "B": [("A", 1), ("D", 5)],
    "C": [("A", 4), ("D", 3)],
    "D": [("A", 2), ("B", 5), ("C", 3)],
}
print(prim(graph, "A"))
# ([('A', 'B', 1), ('A', 'D', 2), ('D', 'C', 3)], 6)
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

```python
class UnionFind:
    def __init__(self, items):
        self.parent = {x: x for x in items}        # each item starts alone

    def find(self, x):
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]   # path compression
            x = self.parent[x]
        return x

    def union(self, x, y):
        rx, ry = self.find(x), self.find(y)
        if rx == ry:
            return False                           # already together: a cycle
        self.parent[rx] = ry
        return True

uf = UnionFind(["A", "B", "C"])
print(uf.union("A", "B"))   # True   — two separate groups, now merged
print(uf.union("A", "B"))   # False  — already together; adding this edge = a cycle
print(uf.find("A") == uf.find("B"))   # True
```

Each group is stored as a tree of parent pointers, and the path-compression line flattens those trees as it walks them, pointing nodes straight at their grandparents. With that optimisation, find and union are effectively constant time — technically O(α(n)), where α is the inverse Ackermann function, a value below 5 for any n you will ever encounter.

:::tip
Union-find is worth knowing well beyond Kruskal's. It's the standard tool for "are these two things in the same group?" — connected components, network connectivity, and grouping in clustering algorithms all use it.
:::

## Kruskal's algorithm: sort the edges

Kruskal's takes a global view. Sort every edge by weight, then walk the sorted list taking each edge unless it would create a cycle. Unlike Prim's, the partial result is a *forest* of separate pieces that gradually merge into one tree.

```python
def kruskal(nodes, edge_list):
    """edge_list: [(weight, u, v), ...]"""
    uf = UnionFind(nodes)
    mst, total = [], 0
    for w, u, v in sorted(edge_list):              # cheapest edges first
        if uf.union(u, v):                         # False means it'd be a cycle
            mst.append((u, v, w))
            total += w
            if len(mst) == len(nodes) - 1:         # V-1 edges: done
                break
    return mst, total

edges = [(1, "A", "B"), (2, "A", "D"), (3, "C", "D"), (4, "A", "C"), (5, "B", "D")]
print(kruskal(["A", "B", "C", "D"], edges))
# ([('A', 'B', 1), ('A', 'D', 2), ('C', 'D', 3)], 6)
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

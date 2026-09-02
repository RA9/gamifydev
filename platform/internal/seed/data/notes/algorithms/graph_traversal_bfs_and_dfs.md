# Graph Traversal: BFS and DFS

A graph is a set of nodes joined by edges, and unlike a tree it can loop back on itself. That one difference changes everything about how you explore it. **Breadth-first search** (BFS) and **depth-first search** (DFS) are the two fundamental ways to visit every reachable node, and between them they underpin an enormous share of the algorithms in the rest of this course.

In this lesson you'll implement both, understand exactly why the visited set is not optional, and see the surprisingly long list of problems that are secretly just a traversal.

## The graph we'll explore

We'll use an **adjacency list** — a dictionary mapping each node to a list of its neighbours. It's the representation you'll want almost every time, because it stores only the edges that exist.

```python
graph = {
    "A": ["B", "C"],
    "B": ["A", "D"],
    "C": ["A", "D"],
    "D": ["B", "C", "E"],
    "E": ["D"],
}
```

```text
        A
       / \
      B   C
       \ /
        D
        |
        E
```

Note that A-B-D-C-A forms a **cycle**: starting at A you can walk in a loop and return to A. Trees can't do this. That loop is the thing that will bite you if you're careless.

## Breadth-first search: a queue

BFS explores in rings. Visit the start node, then all its neighbours, then everything one step beyond those, and so on outwards. The tool is a **queue** — first in, first out — exactly like level-order traversal of a tree.

```python
from collections import deque

def bfs(graph, start):
    visited = {start}                       # nodes already queued or seen
    order = []
    queue = deque([start])
    while queue:
        node = queue.popleft()              # take from the FRONT
        order.append(node)
        for neighbour in graph[node]:
            if neighbour not in visited:
                visited.add(neighbour)      # mark on ENQUEUE, not on visit
                queue.append(neighbour)
    return order

print(bfs(graph, "A"))   # ['A', 'B', 'C', 'D', 'E']
```

```text
queue: [A]           visit A, enqueue B, C
queue: [B, C]        visit B, D is new -> enqueue D
queue: [C, D]        visit C, A and D already visited -> nothing
queue: [D]           visit D, E is new -> enqueue E
queue: [E]           visit E, then queue empties — done
```

Because the queue always hands you the *oldest* pending node, BFS finishes an entire distance-ring before starting the next one. That's the property everything else depends on.

:::warning
Mark a node visited when you **enqueue** it, not when you dequeue it. Otherwise a node with two neighbours pointing at it gets queued twice and visited twice, and on a dense graph the queue can blow up.
:::

## Depth-first search: a stack, or recursion

DFS commits: it follows one path as far as it goes, backtracking only when it runs out of unvisited neighbours. The recursive version is natural, since the call stack does the remembering.

```python
def dfs(graph, node, visited=None, order=None):
    if visited is None:
        visited, order = set(), []
    visited.add(node)
    order.append(node)
    for neighbour in graph[node]:
        if neighbour not in visited:
            dfs(graph, neighbour, visited, order)
    return order

print(dfs(graph, "A"))   # ['A', 'B', 'D', 'C', 'E']
```

Follow it: A → B → D → C, and from C both neighbours are already visited, so it backs up to D and finds E. Compare that to BFS's `['A', 'B', 'C', 'D', 'E']` — same graph, same start, a completely different shape of exploration.

The iterative version swaps the queue for a **stack** — last in, first out — and is otherwise nearly identical to BFS. That similarity is the point.

```python
def dfs_iterative(graph, start):
    visited, order = set(), []
    stack = [start]
    while stack:
        node = stack.pop()                  # take from the BACK (vs popleft)
        if node in visited:
            continue
        visited.add(node)
        order.append(node)
        for neighbour in reversed(graph[node]):
            if neighbour not in visited:
                stack.append(neighbour)
    return order

print(dfs_iterative(graph, "A"))   # ['A', 'B', 'D', 'C', 'E']
```

:::key
BFS and DFS are the **same algorithm** with a different container. Take the oldest pending node (a queue) and you explore in rings; take the newest (a stack) and you plunge down one branch. Everything else is bookkeeping.
:::

## The visited set, and what traversal costs

Remove the `visited` check and watch what happens on our graph. From A you go to B; from B back to A; from A to B again — forever. The traversal never terminates, because the cycle A-B-A has no natural end.

```text
without visited:   A -> B -> A -> B -> A -> ...   never stops

with visited:      A -> B -> (A seen, skip) -> D -> ...   terminates
```

This is the single biggest difference between traversing a tree and traversing a graph: a tree has no cycles, so "keep going down" always ends at a leaf, and a graph makes no such promise. The visited set does a second job too — even on a **directed acyclic graph**, it prevents re-exploring nodes reachable by multiple paths, which without it can blow up exponentially.

:::analogy
Exploring a graph without a visited set is wandering a hedge maze with no thread and no chalk. Every junction looks new, so you re-walk the same loop forever. The visited set is the chalk mark on the ground: "been here, don't bother."
:::

That check is also what makes the cost easy to state. Let V be the number of vertices (nodes) and E the number of edges. With an adjacency list, both BFS and DFS are **O(V + E)** in time:

- Each node is dequeued/popped exactly once thanks to the visited set — O(V).
- For each node, you scan its adjacency list once. Summed over all nodes, that touches every edge — O(E) for a directed graph, O(2E) for an undirected one, which is still O(E).

Space is O(V): the visited set plus the queue or stack, each of which holds at most V nodes.

:::warning
With an **adjacency matrix** instead of a list, finding a node's neighbours means scanning a whole row of V entries — so traversal becomes O(V²) regardless of how few edges there are. On a sparse graph that's dramatically worse. Match your representation to your algorithm.
:::

## BFS finds shortest paths — in unweighted graphs

Here's BFS's headline property. Because it finishes each distance-ring before starting the next, **the first time BFS reaches a node is via a path with the fewest edges**. Track each node's parent and you get the actual route.

```python
def shortest_path(graph, start, goal):
    if start == goal:
        return [start]
    parent = {start: None}
    queue = deque([start])
    while queue:
        node = queue.popleft()
        for neighbour in graph[node]:
            if neighbour not in parent:
                parent[neighbour] = node
                if neighbour == goal:                # rebuild the path
                    path = [goal]
                    while parent[path[-1]] is not None:
                        path.append(parent[path[-1]])
                    return list(reversed(path))
                queue.append(neighbour)
    return None                                       # goal unreachable

print(shortest_path(graph, "A", "E"))   # ['A', 'B', 'D', 'E']
```

The words **unweighted** and **fewest edges** are load-bearing. BFS counts hops, treating every edge as costing the same. The moment edges carry different costs — road distances, network latencies, ticket prices — a three-hop cheap route can beat a two-hop expensive one, and BFS's guarantee collapses. That's the problem the next lesson solves.

DFS makes no shortest-path promise at all. It happily returns a long, wandering route because it commits to the first path it finds.

:::key
BFS gives the shortest path in an **unweighted** graph, because it explores in order of hop count. DFS gives *a* path, not the shortest one.
:::

## What traversal is secretly used for

An enormous number of problems reduce to "run a traversal and watch what happens."

**Connected components.** Loop over every node; if it hasn't been visited, start a fresh traversal from it. Each traversal sweeps up one connected island — friend clusters in a social graph, separate regions on a map.

**Cycle detection.** In an undirected graph, if DFS reaches an already-visited node that isn't the one you just came from, you've found a cycle. In a directed graph, you look for an edge back to a node still on the current recursion stack.

**Flood fill.** The paint-bucket tool in an image editor is DFS or BFS over a grid where each pixel's neighbours are the four adjacent pixels of the same colour. Minesweeper's cascade of empty squares is the same algorithm.

**Shortest hops.** Degrees of separation between two people, the fewest moves in a puzzle, the fewest flights between two airports — all BFS on an unweighted graph.

**Topological sort.** For a directed acyclic graph of dependencies — build steps, course prerequisites, spreadsheet cells — a DFS that records nodes in post-order and reverses the result gives a valid processing order.

:::example
Flood fill on a grid: treat each cell as a node whose neighbours are the up/down/left/right cells sharing the original colour. BFS from the clicked pixel, recolouring as you go, and the region fills outward in rings — visually exactly what the paint bucket does.
:::

## Check Your Understanding

:::quiz
Q: What happens if you run BFS or DFS on a cyclic graph without a visited set?
- It returns the wrong nodes but still finishes
- It runs forever, looping around the cycle *
- It raises an error immediately
- It works fine; the visited set is only an optimisation
E: A cycle gives the traversal an endless supply of "new" nodes to move to. The visited set is what makes graph traversal terminate at all.
:::

:::quiz
Q: What is the time complexity of BFS on a graph stored as an adjacency list?
- O(V)
- O(V + E) *
- O(V²)
- O(E log V)
E: Every vertex is processed once and every edge is examined once across the whole traversal. (With an adjacency matrix it would be O(V²) instead.)
:::

## Talk about it

> BFS and DFS differ by one line — a queue instead of a stack — yet they produce completely different exploration orders and have different uses. Pick a real search you've done (looking for a file, exploring a game map, researching a topic) and describe whether you naturally went breadth-first or depth-first, and what that cost you.

## What's next

BFS gives you the shortest path when every edge costs the same — but real maps have distances, real networks have latencies, and real journeys have prices. Next up is **Shortest Paths: Dijkstra and Bellman-Ford**, two algorithms that handle weighted edges. Dijkstra is BFS upgraded with a priority queue; Bellman-Ford is slower but handles a case Dijkstra flatly cannot.

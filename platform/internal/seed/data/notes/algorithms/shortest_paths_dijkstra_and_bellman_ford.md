# Shortest Paths: Dijkstra and Bellman-Ford

BFS finds the route with the fewest hops. But real graphs have **weights** — roads have lengths, network links have latencies, flights have prices — and there the cheapest route often takes more hops than the shortest one. Fewest hops and lowest cost are different questions.

This lesson covers the two classic answers. **Dijkstra's algorithm** is fast and is what a map app is doing under the hood. **Bellman-Ford** is slower but handles negative edge weights and can detect a situation where "shortest path" stops making sense at all.

## The weighted graph

We'll store weights alongside neighbours: each entry is a `(neighbour, cost)` pair.

```python
graph = {
    "A": [("B", 1), ("C", 4)],
    "B": [("C", 2), ("D", 6)],
    "C": [("D", 3)],
    "D": [],
}
```

```text
              1
       A -----------> B
        \            / \
       4 \        2 /   \ 6
           v      v       v
             C ---------> D
                   3
```

The fewest-hops routes from A to D are A→C→D (4 + 3 = 7) and A→B→D (1 + 6 = 7), both two edges. But A→B→C→D takes three edges and costs 1 + 2 + 3 = **6**. BFS would confidently return a worse route, because it can't see the numbers.

:::key
In a weighted graph, "shortest" means lowest total edge weight, not fewest edges. BFS optimises hop count and is simply answering a different question.
:::

## Relaxation: the shared core

Both algorithms in this lesson are built on one operation, and understanding it is most of the battle. Keep a table `dist` holding your **best known** distance from the start to every node. Initialise the start to 0 and everything else to infinity — "no route found yet." Then, for any edge `u → v` with weight `w`, ask a single question:

> Is going to `u` and then taking this edge cheaper than the best route to `v` I already know?

```python
if dist[u] + w < dist[v]:
    dist[v] = dist[u] + w
    parent[v] = u
```

That's **relaxation**. The name comes from thinking of `dist[v]` as an over-tight upper bound on the true distance that gets relaxed downwards as you discover better routes.

```text
dist[A]=0  dist[B]=inf  dist[C]=inf  dist[D]=inf

relax A->B (1):   0 + 1 < inf   ->  dist[B] = 1
relax A->C (4):   0 + 4 < inf   ->  dist[C] = 4
relax B->C (2):   1 + 2 < 4     ->  dist[C] = 3   (improved!)
relax C->D (3):   3 + 3 < inf   ->  dist[D] = 6
```

Every shortest-path algorithm is just a policy for **which edges to relax, and in what order**. Dijkstra and Bellman-Ford differ only in that policy.

:::analogy
Relaxation is comparing routes on a travel site. You have a saved best price for reaching a city; every new itinerary you check either beats it — so you save the cheaper one — or it doesn't, and you move on.
:::

## Dijkstra: greedy, with a priority queue

Dijkstra's policy is greedy: **always expand the unfinished node with the smallest known distance.** Once you expand a node, you declare its distance final and never revisit it. To always get the smallest, you need a **priority queue** — a min-heap, exactly the structure from the Data Structures course. It's BFS with the plain queue swapped for a heap ordered by distance.

```python
import heapq

def dijkstra(graph, start):
    dist = {node: float("inf") for node in graph}
    dist[start] = 0
    heap = [(0, start)]                        # (distance, node)
    done = set()
    while heap:
        d, u = heapq.heappop(heap)             # smallest known distance
        if u in done:
            continue                           # a stale, worse entry
        done.add(u)
        for v, w in graph[u]:
            if d + w < dist[v]:                # relax the edge u -> v
                dist[v] = d + w
                heapq.heappush(heap, (dist[v], v))
    return dist

print(dijkstra(graph, "A"))
# {'A': 0, 'B': 1, 'C': 3, 'D': 6}
```

```text
pop (0, A)   relax A->B: dist[B]=1     relax A->C: dist[C]=4
pop (1, B)   relax B->C: 1+2=3 < 4 -> dist[C]=3
             relax B->D: 1+6=7     -> dist[D]=7
pop (3, C)   relax C->D: 3+3=6 < 7 -> dist[D]=6
pop (4, C)   already done -> skip (a stale entry)
pop (6, D)   no outgoing edges;  pop (7, D)  already done -> skip
```

Note the "stale entry" trick: rather than updating a node's priority inside the heap, you push a new better entry and skip any old one you pop later. **Complexity: O((V + E) log V)** — each of the V nodes is finalised once, each of the E edges can push one heap entry, and every push and pop costs O(log V). That's the same in the best, average and worst case.

:::tip
To recover the route and not just the distance, record `parent[v] = u` inside the relaxation, then walk the parent chain backwards from your goal. That's the same trick BFS used to rebuild its path.
:::

## Why negative weights break Dijkstra

Dijkstra's whole correctness argument rests on one assumption: **once you pop the smallest tentative distance, nothing can ever make it smaller.** Why would that be true? Because every remaining route to that node has to go through some other unfinished node whose distance is already larger, and adding a non-negative edge weight can only make it larger still.

Remove "non-negative" and the argument collapses. Here's a concrete failure:

```text
          A
       1 / \ 5
        v   v
        B    C
        |    ^
        +----+
          -4

edges:  A->B 1    A->C 5    B->C -4

true shortest A to C:  A -> B -> C  =  1 + (-4)  =  -3
```

Dijkstra pops A (distance 0), relaxes to get `dist[B] = 1` and `dist[C] = 5`. It then pops B — the smallest — and relaxes B→C, giving 1 + (-4) = -3. Whether it fixes `dist[C]` depends on the implementation, but the deeper problem is the one that bites in bigger graphs: **Dijkstra may finalise a node before a negative edge elsewhere has had a chance to improve it**, and once finalised, the node is never reconsidered. There's no error message. The algorithm terminates happily and hands you a plausible-looking, wrong number.

:::warning
Dijkstra **requires** non-negative edge weights. On a graph with a negative edge it doesn't crash or loop — it returns a wrong answer with total confidence. Always check your weights before reaching for it.
:::

## Bellman-Ford: relax everything, and detect negative cycles

Bellman-Ford throws away the cleverness. Its policy: **relax every edge in the graph, and repeat V-1 times.** Why V-1? Any shortest path in a graph with V nodes uses at most V-1 edges — more than that and it would have to revisit a node, meaning it contains a cycle, and dropping the cycle gives a path that's no longer. After round 1, every correct one-edge path is settled; after round 2, every two-edge path; after V-1 rounds, everything.

```python
def bellman_ford(graph, start):
    dist = {node: float("inf") for node in graph}
    dist[start] = 0
    edges = [(u, v, w) for u in graph for v, w in graph[u]]

    for _ in range(len(graph) - 1):            # V-1 rounds
        changed = False
        for u, v, w in edges:
            if dist[u] != float("inf") and dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
                changed = True
        if not changed:                        # settled early
            break

    for u, v, w in edges:                      # one extra round: still improving?
        if dist[u] != float("inf") and dist[u] + w < dist[v]:
            raise ValueError("graph contains a negative cycle")
    return dist

print(bellman_ford(graph, "A"))
# {'A': 0, 'B': 1, 'C': 3, 'D': 6}
```

Because it never commits to a node being final, negative edges are no problem — a later round simply improves the value. **Complexity: O(V · E)** — V-1 rounds over E edges, dramatically worse than Dijkstra on a large graph. The `changed` flag stops early when nothing improved, but the worst case stands.

Now look at that final loop, Bellman-Ford's other gift. A **negative cycle** is a loop whose edge weights sum to less than zero. Walk it once and your total drops; walk it again and it drops further. There is no shortest path at all, because looping once more always makes the cost lower.

```text
         -2                one lap B -> C -> B costs -2 + 1 = -1
    B --------> C          two laps cost -2, three laps -3, ...
    ^           |          no minimum exists
    |     1     |
    +-----------+
```

After V-1 rounds, a graph with no negative cycle has converged and nothing more can improve. So if **one extra round still finds an improvement**, a negative cycle must exist — exactly what the second loop checks. This isn't academic: in currency arbitrage, nodes are currencies and edge weights are the negative logarithm of exchange rates, so a negative cycle is a sequence of trades that returns more money than you started with.

:::key
After V-1 rounds Bellman-Ford has converged, *unless* a negative cycle exists. So one more round that still improves something is proof of a negative cycle — and proof that no shortest path is defined.
:::

## Which do I use?

```text
All edge weights non-negative?      ->  Dijkstra        O((V+E) log V)
Any negative edge weights?          ->  Bellman-Ford    O(V * E)
Need to detect a negative cycle?    ->  Bellman-Ford    O(V * E)
All edge weights equal (or none)?   ->  BFS             O(V + E)
```

The rule in one sentence: **use BFS if the graph is unweighted, Dijkstra if the weights are non-negative, and Bellman-Ford only when negative weights are genuinely possible.** Bellman-Ford is much slower, so don't reach for it defensively.

:::example
Road distances, latencies, ticket prices and physical lengths are all non-negative by nature — Dijkstra territory. Negative weights appear when the weight is a profit, an energy gain, or a logarithm of a ratio.
:::

## Check Your Understanding

:::quiz
Q: Why does Dijkstra's algorithm require non-negative edge weights?
- Negative numbers make the priority queue crash
- It finalises each node's distance permanently, and a negative edge could later have improved it *
- Priority queues cannot store negative values
- It would loop forever
E: Dijkstra's greedy commitment relies on the fact that extending a path can only increase its cost. A negative edge breaks that, so a node may be finalised too early and never corrected.
:::

:::quiz
Q: What does Bellman-Ford's extra relaxation round after V-1 rounds detect?
- Whether the graph is connected
- A negative cycle *
- The shortest path's hop count
- Whether the start node exists
E: After V-1 rounds the distances have converged unless something makes them improve forever. An improvement in one more round proves a negative cycle exists, so no shortest path is defined.
:::

## Talk about it

> Dijkstra is fast because it commits: once a node's distance is settled, it never looks again. Bellman-Ford is slow because it refuses to commit until the very end. Describe a decision-making process — in code or in life — where committing early makes you fast, and explain what assumption has to hold for that speed to be safe.

## What's next

Dijkstra explores outward in every direction, which is wasteful when you know roughly where you're going. Next up is **A-Star and Heuristic Search**, which adds an estimate of the remaining distance so the search leans towards the goal instead of spreading evenly. It's the algorithm behind game pathfinding and route planning, and you'll see that Dijkstra is really just A* with the estimate switched off.

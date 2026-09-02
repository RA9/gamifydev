# Graphs and How to Represent Them

A tree let each node have many children but exactly one parent, and forbade loops. A **graph** drops both restrictions. Anything can connect to anything, in any pattern, including cycles. That makes it the most general structure in this course — and the one that models the messiest and most interesting real data.

Road networks, friendships, web links, dependencies between tasks, connections between neurons: all graphs. This lesson gives you the vocabulary and, more importantly, the two ways to actually store one in memory, along with a clear rule for picking between them.

## Vertices and edges

A graph is two sets: a set of **vertices** (also called nodes) and a set of **edges** (connections between pairs of vertices). That's the whole definition — no root, no hierarchy, no rules about shape.

Here's a small social graph we'll use throughout:

```text
     A ------- B
     |       / |
     |      /  |
     C ----'   D          E
```

Five vertices: A, B, C, D, E. Four edges: A–B, A–C, B–C, B–D. Vertex E has no edges at all, which is perfectly legal.

Two counts get used constantly, so they get their own letters: **V** is the number of vertices, **E** the number of edges. Every complexity in this lesson is written in terms of those.

**Degree.** The degree of a vertex is how many edges touch it. Here A has degree 2, B has degree 3, C has degree 2, D has degree 1, and E has degree 0. (In a directed graph you separate this into *in-degree* and *out-degree* — edges arriving versus edges leaving.)

:::analogy
A tree is a company org chart: one boss at the top, one manager per person, no loops. A graph is the company's actual social life — who talks to whom, in whatever tangle that happens to be, with cliques and loops and the occasional person nobody has met.
:::

## The four questions to ask about any graph

Whenever you meet a graph, four properties tell you almost everything.

**Directed or undirected?** In an **undirected** graph an edge goes both ways: if A is friends with B, B is friends with A. In a **directed** graph (or digraph) each edge has an arrow and only works one way — think "A follows B" on social media, or a one-way street.

```text
 undirected: A --- B        A and B both connect to each other
 directed:   A --> B        A points to B, but not the reverse
```

**Weighted or unweighted?** An **unweighted** edge simply exists. A **weighted** edge carries a number — a distance, a cost, a travel time, a bandwidth. Road maps are weighted graphs; that weight is what shortest-path algorithms actually minimise.

**Cyclic or acyclic?** A **cycle** is a path that returns to where it started. Our graph has one: A → B → C → A. A graph with no cycles at all is **acyclic**. A **directed acyclic graph** — a **DAG** — is so useful it gets its own initials: it's the shape of task dependencies, build systems, course prerequisites, spreadsheet formulas and version-control history. If a dependency graph ever contains a cycle, you have a circular dependency and nothing can be ordered.

**Connected?** A graph is **connected** if you can reach every vertex from every other. Ours is not — E is stranded. The reachable clusters are called **connected components**, and our graph has two: {A, B, C, D} and {E}.

:::key
A tree is just a special graph: connected, undirected, acyclic, with exactly V − 1 edges. Everything you learned about trees is a restricted case of what you're learning now.
:::

## Representation 1: the adjacency list

The first way to store a graph is to give each vertex a list of its neighbours. In Python that's naturally a dictionary of lists.

```python
graph = {
    "A": ["B", "C"],
    "B": ["A", "C", "D"],
    "C": ["A", "B"],
    "D": ["B"],
    "E": [],
}

print(graph["B"])         # ['A', 'C', 'D']
print(len(graph["B"]))    # 3  -- the degree of B
print("D" in graph["A"])  # False -- no edge between A and D
```

Because the graph is undirected, every edge appears twice — B is in A's list and A is in B's. For a directed graph you'd store each edge only in the source vertex's list.

Space is **O(V + E)**: one entry per vertex, plus one list slot per edge endpoint. Nothing is stored for edges that don't exist, which is the key property.

For weights, swap the inner list for a dictionary mapping neighbour to weight:

```python
roads = {
    "A": {"B": 5, "C": 2},
    "B": {"A": 5, "D": 7},
    "C": {"A": 2},
    "D": {"B": 7},
}
print(roads["A"]["C"])   # 2
```

## Representation 2: the adjacency matrix

The second way is a V × V grid of 0s and 1s. Row `i`, column `j` holds 1 if there's an edge from vertex `i` to vertex `j`, and 0 otherwise.

```text
        A  B  C  D  E
    A [ 0  1  1  0  0 ]
    B [ 1  0  1  1  0 ]
    C [ 1  1  0  0  0 ]
    D [ 0  1  0  0  0 ]
    E [ 0  0  0  0  0 ]
```

For an undirected graph the matrix is symmetric across the diagonal — the top-right half mirrors the bottom-left — because every edge is recorded in both directions. A directed graph's matrix generally is not symmetric.

```python
labels = ["A", "B", "C", "D", "E"]
matrix = [
    [0, 1, 1, 0, 0],
    [1, 0, 1, 1, 0],
    [1, 1, 0, 0, 0],
    [0, 1, 0, 0, 0],
    [0, 0, 0, 0, 0],
]

b, d = labels.index("B"), labels.index("D")
print(matrix[b][d])    # 1  -- B and D are connected
print(sum(matrix[b]))  # 3  -- the degree of B
```

Checking whether a specific edge exists is one array access: **O(1)**, unbeatable. For weights, store the weight instead of 1 (using `None` or infinity for "no edge").

The catch is the empty space. Our graph has 4 edges and the matrix has 25 cells, 17 of which are zeros stored at full cost. Space is **O(V²)** whether the graph is crowded or nearly empty.

:::warning
An adjacency matrix for a million-vertex graph would need a trillion cells. Social networks, road maps and web-link graphs all have millions of vertices and only a handful of edges each — for those, a matrix is not merely wasteful, it's impossible.
:::

## Choosing between them

```text
 operation / property            adjacency list     adjacency matrix
 --------------------------------------------------------------------
 space                           O(V + E)           O(V^2)
 is there an edge u -> v?        O(deg(u))          O(1)
 list all neighbours of u        O(deg(u))          O(V)
 add an edge                     O(1)               O(1)
 remove an edge                  O(deg(u))          O(1)
 add a vertex                    O(1)               O(V^2) (rebuild)
 iterate over every edge         O(V + E)           O(V^2)
 best when                       E is small         E is close to V^2
```

The decision hinges on one word: **density**.

A graph is **sparse** when E is far smaller than V² — most possible edges are absent. This is the overwhelmingly common case in real data. Your friends number in the hundreds, not the billions; a city intersection connects to four roads, not to every other intersection.

A graph is **dense** when E approaches V² — most pairs are connected. This happens with small, complete-ish graphs: a distance table between 50 cities, a tournament where everyone plays everyone.

:::key
**Sparse graph → adjacency list. Dense graph → adjacency matrix.** Most real-world graphs are sparse, so the adjacency list is the default choice, and the one you'll see in nearly all algorithm code.
:::

Note the row for "list all neighbours". Most graph algorithms spend their whole life asking exactly that question, vertex after vertex. On an adjacency list, answering it for every vertex totals O(V + E). On a matrix it's O(V²) no matter how few edges exist — you must scan a whole row of mostly zeros each time. That single line is why the adjacency list dominates in practice.

:::tip
When you build an adjacency list for an undirected graph, remember to add the edge in *both* directions. Adding it once quietly turns your undirected graph into a directed one, and every algorithm you run on it will give subtly wrong answers.
:::

To feel the size difference, put real numbers on it.

:::example
A social network with 1,000,000 users averaging 200 friends each: an adjacency list stores about 200,000,000 entries. A matrix would need 1,000,000² = 1,000,000,000,000 cells — five thousand times more, almost all of them zero.
:::

## Check Your Understanding

:::quiz
Q: You're storing a road network with 500,000 intersections, where each intersection connects to about 4 others. Which representation should you use?
- Adjacency matrix, because edge lookup is O(1)
- Adjacency list, because the graph is sparse *
- Either one; they use the same space
- Neither; a graph can't store roads
E: With ~2,000,000 edges and 500,000 vertices, this graph is extremely sparse. A matrix would need 250,000,000,000 cells. The list needs O(V + E).
:::

:::predict
Q: What does this print?
```python
graph = {"A": ["B", "C"], "B": ["A", "C", "D"],
         "C": ["A", "B"], "D": ["B"], "E": []}
print(len(graph["B"]), len(graph["E"]), "D" in graph["C"])
```
- 3 0 False *
- 3 0 True
- 4 1 False
- 2 0 False
E: B has three neighbours, E has none, and there is no edge between C and D.
:::

:::match
Q: Match each graph term to its meaning.
- Directed | Edges have a one-way arrow
- Weighted | Edges carry a number such as distance or cost
- DAG | Directed, with no cycles anywhere
- Connected component | A cluster of vertices all reachable from one another
E: These four properties describe almost any graph you'll meet, and together they determine which algorithms apply to it.
:::

## Talk about it

> Pick a network you know well — a transit map, a group chat, the imports between files in a project. Describe it as a graph: what are the vertices, what are the edges, is it directed, is it weighted, and could it contain a cycle? Then say which representation you'd choose and why.

## What's next

You've reached the end of the Data Structures course, and it's worth pausing on how far that is. You started not knowing what a data structure *was*, and you can now explain contiguous versus linked memory, amortized growth, hashing with its honest worst case, five kinds of tree, and the density trade-off you just met. That's a real foundation.

Storing a graph is only half the story — the other half is walking one, and breadth-first and depth-first search are waiting for you in the Algorithms course, where the queue and the stack you built here turn out to be the whole difference between them. Before that, though, comes the course that makes all of these comparisons rigorous: **Complexity and Analysis**, where Big O stops being a shorthand you've been trusting and becomes a tool you can actually wield. Pixel will see you there.

# Trees and Terminology

Every structure so far has been *linear*: arrays, linked lists, stacks and queues all lay their items out in one line. A **tree** breaks that line open. Each item can point to several others, so the structure fans out — and suddenly you can model hierarchies, which is how an enormous amount of real-world data is actually shaped.

This lesson is mostly vocabulary, and that's on purpose. The next four lessons all lean on these words, so we'll fix them to one concrete example and keep coming back to it. Learn the picture and the words follow.

## Our running example

Here's a small project folder. Keep this diagram in view for the rest of the lesson.

```text
                        project/
                     /      |      \
                 src/     docs/   README.md
                /    \        \
          main.py   utils/   guide.md
                       \
                    helpers.py
```

Eight items, connected in a hierarchy. Every folder contains things; every file sits inside exactly one folder. That "exactly one" is the essence of a tree.

:::analogy
A family tree, an org chart, a table of contents, and the folders on your computer are all the same shape. Each entry has one thing directly above it and any number of things directly below.
:::

## The words

**Node.** One item in the tree. `src/` is a node; so is `guide.md`. Our tree has 8 nodes.

**Edge.** A single connection between a node and one directly below it — the lines in the diagram. Our tree has 7 edges. That's not a coincidence: a tree with n nodes always has exactly **n − 1** edges, because every node except the root is reached by exactly one edge from above.

**Root.** The single node at the top, with nothing above it. Here that's `project/`. A tree has exactly one root.

**Parent and child.** If an edge runs from `src/` down to `main.py`, then `src/` is the **parent** and `main.py` is the **child**. Every node has exactly one parent, except the root, which has none. Nodes sharing a parent are **siblings** — `main.py` and `utils/` are siblings.

**Leaf.** A node with no children. Our leaves are `main.py`, `helpers.py`, `guide.md` and `README.md` — the actual files, as it happens.

**Internal node.** Any node that isn't a leaf; it has at least one child. Here: `project/`, `src/`, `docs/`, `utils/`.

**Degree.** The number of children a node has. `project/` has degree 3, `src/` has degree 2, `docs/` has degree 1, and every leaf has degree 0.

:::key
Root at the top, leaves at the bottom, one parent per node. Everything else in tree vocabulary is a way of measuring distance in that picture.
:::

## Path, depth and height

**Path.** The sequence of nodes you walk through to get from one node to another. From the root to `helpers.py` the path is `project/ -> src/ -> utils/ -> helpers.py`. In a tree, there is exactly **one** path between any two nodes — never zero, never two.

**Depth** (of a node). How many edges lie between the root and that node. The root has depth 0.

```text
 depth 0:  project/
 depth 1:  src/    docs/    README.md
 depth 2:  main.py    utils/    guide.md
 depth 3:  helpers.py
```

**Height** (of a node). The number of edges on the longest downward path from that node to a leaf. Every leaf has height 0. `utils/` has height 1. `src/` has height 2. The **height of the tree** is the height of its root — here, **3**.

Depth counts downwards from the root; height counts upwards from the leaves. Beginners mix these up constantly, so say it once more: depth is measured from the top, height from the bottom.

:::warning
Some textbooks count depth and height in *nodes* rather than *edges*, which shifts every number by one. Neither is wrong — but before you compare answers with someone, agree on which convention you're using. This course counts edges.
:::

**Subtree.** Pick any node; that node together with everything hanging below it is a subtree. The subtree rooted at `src/` contains `src/`, `main.py`, `utils/` and `helpers.py`. This idea is quietly powerful: *a subtree is itself a tree*, which is exactly why tree code is so naturally recursive.

:::example
The subtree rooted at `docs/` is `docs/ -> guide.md`. It has its own root (`docs/`), its own leaf (`guide.md`), and a height of 1 — a complete little tree in its own right.
:::

## A tree is a connected acyclic graph

Here's the formal definition, which will make more sense after the graphs lesson but is worth meeting now.

A tree is a set of nodes and edges that is:

- **connected** — you can get from any node to any other by following edges, and
- **acyclic** — there are no loops; you can never follow edges and arrive back where you started.

Those two conditions are what force the properties above. Connected guarantees at least one path between any pair; acyclic guarantees at most one. Together: exactly one path, one parent per node, and n − 1 edges.

```text
  a tree                     NOT a tree (has a cycle)

     A                            A
    / \                          / \
   B   C                        B---C

  one path from B to C       two paths from B to C
```

## Why hierarchies are everywhere

Trees model containment and ranking, and software is full of both.

**Filesystems.** Exactly our example. Folders contain folders contain files.

**The DOM.** A web page is a tree of elements: `<html>` contains `<body>` contains `<div>` contains `<p>`. When you write CSS like `div p`, you're describing an ancestor relationship in that tree.

**Org charts and category menus.** One manager per employee; one parent category per subcategory.

**Program structure.** A compiler parses your source code into a tree — an expression like `2 * (3 + 4)` becomes a node for `*` whose children are `2` and a `+` node. Evaluating it is a walk over that tree.

## Trees in code

The natural representation mirrors the picture: a node holds its value and a list of references to its children.

```python
class TreeNode:
    def __init__(self, name):
        self.name = name
        self.children = []

root = TreeNode("project/")
src = TreeNode("src/")
docs = TreeNode("docs/")
utils = TreeNode("utils/")

root.children = [src, docs, TreeNode("README.md")]
src.children = [TreeNode("main.py"), utils]
docs.children = [TreeNode("guide.md")]
utils.children = [TreeNode("helpers.py")]

print(len(root.children))     # 3
print(root.children[0].name)  # src/
```

That's the linked-list idea again, with a *list* of next-pointers instead of a single one.

Because every subtree is a tree, tree code is almost always recursive: handle this node, then call yourself on each child. Height is a perfect example — the height of a node is one more than the tallest of its children's heights, and a leaf's height is 0.

```python
def height(node):
    if not node.children:      # base case: a leaf
        return 0
    return 1 + max(height(child) for child in node.children)

print(height(root))    # 3
print(height(src))     # 2
print(height(docs))    # 1
```

Compare that with the diagram: the longest path from `project/` down is to `helpers.py`, three edges away. The function agrees.

:::tip
When a tree problem feels hard, ask: "what's the answer for a leaf, and how do I combine my children's answers into mine?" That's the base case and the recursive case, and it solves a startling share of tree exercises.
:::

## Check Your Understanding

:::quiz
Q: In the project tree, what is the depth of `utils/` and the height of `src/`?
- depth 1, height 1
- depth 2, height 2 *
- depth 3, height 0
- depth 2, height 3
E: `utils/` is two edges below the root, so depth 2. From `src/` the longest downward path is src -> utils -> helpers.py, two edges, so height 2.
:::

:::match
Q: Match each term to its meaning.
- Leaf | A node with no children
- Root | The single node with no parent
- Degree | How many children a node has
- Subtree | A node together with everything below it
E: These four cover most tree vocabulary; depth and height then measure distance from the top and the bottom.
:::

:::fill
Q: Complete the count of edges in any tree with n nodes.
`edges = n - ___`
- 1 *
- 2
- 0
E: Every node except the root is reached by exactly one edge from its parent, so there are n - 1 edges.
:::

## Talk about it

> Pick a hierarchy from your own life — a menu in an app you use, the chapters of a book, the departments where you work or study. Describe it as a tree: what is the root, what are the leaves, roughly how deep does it go, and does anything about it break the "exactly one parent" rule?

## What's next

You now have the vocabulary, which means the interesting trees are finally reachable. Next up is **Binary Trees and Binary Search Trees**, where we limit each node to at most two children and then add one ordering rule. That single rule turns a tree into a searching machine — and shows you exactly what happens when the tree grows lopsided.

# Tree Traversals

A tree branches, so "visit every node" has no single obvious meaning — there are several sensible orders, and each one is genuinely useful for a different job. **Traversal** is the name for systematically visiting every node exactly once.

In this lesson you'll learn the three depth-first orders (pre-order, in-order, post-order), what each one is actually *for*, level-order traversal with a queue, and how to write an in-order traversal iteratively with a stack you manage yourself.

## The tree we'll use

Here's a binary search tree — remember, in a BST every value in a node's left subtree is smaller than the node and every value in its right subtree is larger.

```text
                4
              /   \
            2       6
           / \     / \
          1   3   5   7
```

And a minimal node structure. These nodes have automatic storage, so no allocation or cleanup is needed:

```c
typedef struct Node {
    int value;
    struct Node *left;
    struct Node *right;
} Node;

Node n1 = {1, NULL, NULL}, n3 = {3, NULL, NULL};
Node n5 = {5, NULL, NULL}, n7 = {7, NULL, NULL};
Node n2 = {2, &n1, &n3}, n6 = {6, &n5, &n7};
Node root_node = {4, &n2, &n6};
Node *root = &root_node;
```

Every traversal below has the same three ingredients — visit the node, go left, go right — and differs only in the **order** those three happen. That's the entire idea.

:::key
Pre-order, in-order and post-order are the same recursion with one line moved. The names describe *when the node itself is visited*: before its children, between them, or after them.
:::

## Pre-order: node, left, right

Visit the node first, then recurse into the left subtree, then the right.

```c
#include <stddef.h>

void pre_order(const Node *node, int out[], size_t *count) {
    if (node == NULL) return;             /* Empty subtree. */
    out[(*count)++] = node->value;        /* Visit before the children. */
    pre_order(node->left, out, count);
    pre_order(node->right, out, count);
}

int result[7];
size_t count = 0;
pre_order(root, result, &count);          /* 4 2 1 3 6 5 7 */
```

**What it's for: copying and serialising.** Because a parent is always emitted before its children, you can rebuild the tree by reading the sequence left to right and inserting each value as you go — the structure comes back intact. Post-order or in-order don't give you that. Pre-order is also the order you'd use to print a nested directory listing, or to render a document outline, where the heading must appear before what's under it.

## In-order: left, node, right

Recurse left, visit the node, then recurse right.

```c
#include <stddef.h>

void in_order(const Node *node, int out[], size_t *count) {
    if (node == NULL) return;
    in_order(node->left, out, count);
    out[(*count)++] = node->value;        /* Visit between the subtrees. */
    in_order(node->right, out, count);
}

int result[7];
size_t count = 0;
in_order(root, result, &count);           /* 1 2 3 4 5 6 7 */
```

Look at that output. **On a BST, in-order traversal produces the values in sorted order** — and it does so in O(n), with no sorting at all. The reason is direct: the BST property says everything left of a node is smaller and everything right is larger, so "all the smaller things, then me, then all the larger things" *is* ascending order.

That single fact is why BSTs are worth building. It gives you sorted iteration for free, plus range queries: to list every value between 3 and 6, do an in-order walk and skip subtrees that can't contain anything in range. Trace it with a finger — you always exhaust the entire left branch before saying a node's own name, which is exactly what puts the small values first.

## Post-order: left, right, node

Recurse into both subtrees, and visit the node last.

```c
#include <stddef.h>

void post_order(const Node *node, int out[], size_t *count) {
    if (node == NULL) return;
    post_order(node->left, out, count);
    post_order(node->right, out, count);
    out[(*count)++] = node->value;        /* Visit after both children. */
}

int result[7];
size_t count = 0;
post_order(root, result, &count);         /* 1 3 2 5 7 6 4 */
```

**What it's for: anything where the children must be finished before the parent.** *Deleting a tree*, because you cannot free a node before its children without losing the pointers to them. *Computing sizes* — directory totals, subtree heights, node counts — because a parent's answer is built from its children's. And *evaluating an expression tree*, where leaves are numbers and internal nodes are operators, so both operand subtrees must be evaluated before the operator can be applied:

```c
typedef enum { NUMBER, ADD, MULTIPLY } ExprKind;

typedef struct Expr {
    ExprKind kind;
    int value;
    const struct Expr *left;
    const struct Expr *right;
} Expr;

int evaluate(const Expr *node) {
    if (node->kind == NUMBER) return node->value;
    int left = evaluate(node->left);       /* Children first. */
    int right = evaluate(node->right);
    return node->kind == ADD ? left + right : left * right;
}

Expr two = {NUMBER, 2, NULL, NULL}, three = {NUMBER, 3, NULL, NULL};
Expr four = {NUMBER, 4, NULL, NULL};
Expr sum = {ADD, 0, &two, &three};
Expr tree = {MULTIPLY, 0, &sum, &four};    /* (2 + 3) * 4; evaluates to 20. */
```

```text
same tree, three orders:

pre-order    4  2  1  3  6  5  7      parent before children  -> copy, serialise
in-order     1  2  3  4  5  6  7      sorted on a BST         -> ordered output
post-order   1  3  2  5  7  6  4      children before parent  -> delete, evaluate
```

:::tip
Pick the traversal by asking *when you need the parent's information*. Need it before descending (like a path prefix)? Pre-order. Need results from below first (like a total)? Post-order. Want sorted output from a BST? In-order.
:::

## Level-order: breadth-first with a queue

The three above are all **depth-first** — they plunge to the bottom of one branch before trying the next. **Level-order** traversal does the opposite: it visits all nodes at depth 0, then all at depth 1, and so on. Recursion can't do this naturally, because the call stack is depth-first by construction. Instead you need an explicit **queue** — first in, first out. Take a node from the front, visit it, and push its children onto the back.

```c
#include <stdbool.h>
#include <stddef.h>

bool level_order(const Node *root, int out[], size_t capacity, size_t *count) {
    if (root == NULL) { *count = 0; return true; }
    const Node *queue[64];
    if (capacity > 64) capacity = 64;
    size_t front = 0, back = 0, used = 0;
    queue[back++] = root;
    while (front < back) {
        const Node *node = queue[front++];       /* Take from the front. */
        if (used == capacity) return false;
        out[used++] = node->value;
        if (node->left != NULL) {
            if (back == 64) return false;
            queue[back++] = node->left;
        }
        if (node->right != NULL) {
            if (back == 64) return false;
            queue[back++] = node->right;
        }
    }
    *count = used;
    return true;
}

int result[7];
size_t count;
level_order(root, result, 7, &count);             /* 4 2 6 1 3 5 7 */
```

The queue is what enforces the level ordering. Visiting 4 pushes 2 and 6; visiting 2 pushes 1 and 3 onto the *back*, behind 6 — so 6 is handled before either of them, and the current level always finishes first.

Level-order is what you want for "find the shallowest node matching X", for rendering a tree level by level, and for anything organisation-chart shaped. It's also exactly the algorithm you'll meet next lesson as **breadth-first search** on a graph — a tree is just a graph with no cycles.

:::analogy
Depth-first traversal is exploring a building by walking to the end of one corridor and into every room off it before coming back. Breadth-first is checking every room on floor one, then every room on floor two. Same building, very different order of discovery.
:::

## Iterative in-order with an explicit stack

Recursion is doing something specific for you: the call stack remembers which nodes you still owe a visit to. You can take that job over yourself with an array used as a **stack** — last in, first out. The pattern: walk as far left as you can, pushing every node you pass; when you can't go left any further, pop a node, visit it, and move to its right child.

```c
#include <stdbool.h>
#include <stddef.h>

bool in_order_iterative(const Node *root, int out[], size_t capacity,
                        size_t *count) {
    const Node *stack[64], *current = root;
    size_t top = 0, used = 0;
    while (top > 0 || current != NULL) {
        while (current != NULL) {              /* Dive left; remember the way back. */
            if (top == 64) return false;
            stack[top++] = current;
            current = current->left;
        }
        current = stack[--top];                 /* Nothing more to the left. */
        if (used == capacity) return false;
        out[used++] = current->value;
        current = current->right;
    }
    *count = used;
    return true;
}

int result[7];
size_t count;
in_order_iterative(root, result, 7, &count);      /* 1 2 3 4 5 6 7 */
```

Why bother, when the recursive version is four lines? On a deeply unbalanced tree — a BST built from sorted input degenerates into a chain of n nodes — recursion can exhaust the C call stack, while an explicit stack can be sized, checked, or grown deliberately. The iterative version can also be paused, which lets an iterator walk a huge tree lazily.

**All four traversals are O(n) time**, since each node is visited exactly once. Space is O(h) for the depth-first ones, where h is the tree height — O(log n) for a balanced tree, O(n) for a degenerate one. Level-order's queue can hold an entire level, which for a balanced tree is up to about n/2 nodes, so it is O(n) space.

:::warning
A recursive traversal's space cost is the tree's **height**, not its node count — but on a tree built by inserting already-sorted values, the height *is* the node count. That's the same degenerate case that makes an unbalanced BST slow, and it's why balanced trees exist.
:::

## Check Your Understanding

:::quiz
Q: Which traversal produces sorted output when run on a binary search tree?
- Pre-order
- In-order *
- Post-order
- Level-order
E: The BST property guarantees everything left of a node is smaller and everything right is larger, so visiting left-then-node-then-right yields ascending order.
:::

:::predict
Q: What does this print for the tree with root 4, children 2 and 6, and leaves 1, 3, 5, 7?
```c
#include <stddef.h>

void walk(const Node *node, int out[], size_t *count) {
    if (node == NULL) return;
    walk(node->left, out, count);
    walk(node->right, out, count);
    out[(*count)++] = node->value;
}

int result[7];
size_t count = 0;
walk(root, result, &count);
/* What values are now in result? */
```
- `{1, 3, 2, 5, 7, 6, 4}` *
- `{4, 2, 1, 3, 6, 5, 7}`
- `{1, 2, 3, 4, 5, 6, 7}`
- `{4, 2, 6, 1, 3, 5, 7}`
E: The node is appended after both recursive calls, so this is post-order: both children of every node come out before the node itself, and the root is last.
:::

## Talk about it

> Level-order traversal needs a queue, and depth-first traversal needs a stack — and recursion is really just borrowing the machine's stack. Explain in your own words why swapping a stack for a queue is enough to completely change the order in which a structure gets explored.

## What's next

Trees are the easy case: they branch, but they never loop back on themselves, so you can never arrive somewhere twice. Next up is **Graph Traversal: BFS and DFS**, where you'll apply these same two patterns to graphs — and discover that the moment cycles are possible, you need one extra piece of bookkeeping or your traversal will run forever.

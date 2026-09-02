# Binary Trees and Binary Search Trees

Now we take the general tree and add two restrictions. First: each node may have **at most two children**. Second: the values must be arranged in a specific order. The first restriction gives us a shape that's easy to reason about; the second turns that shape into a search engine.

By the end of this lesson you'll be able to find a value among a million by asking about twenty questions — and you'll also see exactly how that guarantee falls apart if you're careless. That failure is what the next lesson exists to fix.

## Binary trees

A **binary tree** is a tree where every node has at most two children, conventionally called **left** and **right**. That's the entire definition — no ordering yet.

```text
            A
           / \
          B   C
         / \    \
        D   E    F
```

The left/right distinction matters even when a node has only one child. `C` above has a right child and no left child, and that's a genuinely different tree from one where `C` has a left child and no right child.

```python
class Node:
    def __init__(self, key):
        self.key = key
        self.left = None
        self.right = None
```

Two pointers instead of a list of children. Simple, and it's what the rest of this course builds on.

## Full, complete, perfect, balanced

Four adjectives get thrown around for binary trees. They're easy to confuse and each means something specific.

**Full** — every node has either 0 or 2 children. No node has exactly one.

**Complete** — every level is entirely filled except possibly the last, and the last level fills from the left with no gaps. This one matters enormously in the next lesson, because complete trees pack perfectly into an array.

**Perfect** — every internal node has 2 children and all leaves sit at the same depth. A perfect tree of height h has exactly 2^(h+1) − 1 nodes.

**Balanced** — informally, no part of the tree is dramatically deeper than another; formally, the height stays proportional to log n. A common precise rule is that for every node, the heights of its two subtrees differ by at most 1.

```text
   full (0 or 2 children)     complete (fills left to right)

          A                            A
         / \                          / \
        B   C                        B   C
       / \                          / \  /
      D   E                        D  E F

   perfect (all leaves level)      degenerate (not balanced)

          A                          A
         / \                           \
        B   C                           B
       / \ / \                           \
      D  E F  G                           C
```

:::key
Perfect implies complete, and complete implies... nothing about fullness on its own. The one to remember for now is **complete**: all levels full except the last, which fills left to right. Heaps depend on it.
:::

## The binary search tree rule

A **binary search tree** (BST) is a binary tree with one extra promise, and this promise is the entire point:

> For every node, all keys in its **left** subtree are less than the node's key, and all keys in its **right** subtree are greater.

Note "all keys in the subtree", not just the immediate children. The rule applies at every node, all the way down.

```text
                50
               /  \
             30    70
            /  \   /  \
          20   40 60   80
```

Check it at `30`: everything left of it (20) is smaller, everything right of it (40) is bigger. Check it at `50`: the entire left subtree {20, 30, 40} is smaller, the entire right subtree {60, 70, 80} is bigger. The invariant holds everywhere.

:::analogy
A BST is the guessing game "higher or lower". Every node you visit tells you which half of the remaining possibilities to throw away. You never look at the discarded half again.
:::

## Searching and inserting

Search follows the rule mechanically: compare, then go left or right. Each comparison eliminates an entire subtree.

```python
def search(node, key):
    while node is not None:
        if key == node.key:
            return True
        node = node.left if key < node.key else node.right
    return False
```

Looking for 40 in the tree above: at 50, 40 is smaller, go left. At 30, 40 is bigger, go right. At 40 — found, three comparisons.

Insert works the same way. Walk down as if searching, and when you fall off the bottom, that empty spot is exactly where the new key belongs.

```python
def insert(node, key):
    if node is None:
        return Node(key)
    if key < node.key:
        node.left = insert(node.left, key)
    elif key > node.key:
        node.right = insert(node.right, key)
    return node          # duplicate keys are ignored here

root = None
for k in [50, 30, 70, 20, 40, 60, 80]:
    root = insert(root, k)

print(search(root, 40))   # True
print(search(root, 45))   # False
```

There's a free bonus in the ordering rule. Visit the left subtree, then the node, then the right subtree — an **in-order traversal** — and the keys come out sorted:

```python
def inorder(node, out):
    if node is None:
        return
    inorder(node.left, out)
    out.append(node.key)
    inorder(node.right, out)

out = []
inorder(root, out)
print(out)   # [20, 30, 40, 50, 60, 70, 80]
```

A hash table could never do that. This is what we traded ordering away for, and here we get it back.

## The cost is O(h), not O(log n)

Here's the sentence people get wrong. Search, insert and delete in a BST cost **O(h)**, where h is the **height** of the tree. Not O(log n). O(h). The distinction matters because each step of the walk moves down exactly one level, so the work is bounded by how deep the tree goes — and how deep it goes depends entirely on its shape.

When the tree is balanced, the height is about log₂(n), and O(h) becomes the O(log n) everybody quotes:

```text
        n items        balanced height       comparisons
             7                  2                  ~3
     1,000,000                ~19                 ~20
 1,000,000,000                ~29                 ~30
```

A billion items, about thirty comparisons. That's the promise of a balanced BST, and it's genuinely remarkable.

## The degenerate case

Now watch it break. Insert the same keys in *sorted* order:

```python
root = None
for k in [20, 30, 40, 50, 60, 70, 80]:
    root = insert(root, k)
```

Every key is bigger than the one before, so every insert goes right, and right again, and right again:

```text
   20
     \
      30
        \
         40
           \
            50 ... 60 ... 70 ... 80
```

The BST rule is still perfectly satisfied, and an in-order traversal still returns the keys sorted. But the tree has become a linked list. With n nodes the height is n − 1, so O(h) is now **O(n)**. Searching for 80 takes seven comparisons instead of three, and with a million sorted inserts it takes a million.

This is called a **degenerate** or **pathological** tree, and it's not a rare accident. Sorted or nearly-sorted input is extremely common — timestamps, auto-increment IDs, alphabetised names. A plain BST fed real-world data has a nasty habit of collapsing into a list. So never quote O(log n) for a plain BST without attaching the condition *only when balanced*: nothing in the ordering rule keeps it that way.

:::key
BST operations are **O(h)**. Balanced: h ≈ log n, so O(log n). Degenerate: h = n − 1, so O(n). The whole game is keeping h small.
:::

## Check Your Understanding

:::quiz
Q: What is the worst-case time to search a plain binary search tree with n nodes?
- O(1)
- O(log n)
- O(n) *
- O(n log n)
E: If the keys were inserted in sorted order the tree degenerates into a chain of height n - 1, so a search may visit every node.
:::

:::predict
Q: Using the `insert` function from this lesson, what does this print?
```python
root = None
for k in [8, 3, 10, 1]:
    root = insert(root, k)
print(root.left.left.key)
```
- 1 *
- 3
- 8
- 10
E: 3 goes left of 8, then 1 is smaller than both so it goes left of 3. `root.left.left` is the node holding 1.
:::

:::fill
Q: Complete the binary search tree ordering rule.
`all keys in the left subtree < node.key < all keys in the ___ subtree`
- right *
- parent
- sibling
E: Smaller keys live entirely in the left subtree, larger keys entirely in the right — at every node, not just the root.
:::

## Talk about it

> A binary search tree keeps its data sorted and searchable, but only stays fast if it stays bushy rather than stringy. Before reading on, invent your own idea: if you noticed a tree growing lopsided during an insert, what could you do to it — without breaking the ordering rule — to make it shorter again? Describe your idea in plain words.

## What's next

Before we repair the lopsided tree, there's a different use of binary trees worth meeting: one where the shape is *always* perfectly controlled and the whole thing lives inside a plain array. Next up is **Heaps and Priority Queues**, the structure that answers "what's the most urgent item right now?" in constant time — and it uses the complete-tree definition from earlier in this lesson.

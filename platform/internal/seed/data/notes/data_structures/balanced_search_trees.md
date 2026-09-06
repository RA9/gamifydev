# Balanced Search Trees

We left the binary search tree in a bad state. Fed sorted data, it collapsed into a chain and its O(log n) promise became O(n). Everything about a BST is fine except its *shape* — and shape is the one thing the ordering rule doesn't control.

A **balanced search tree** is a BST that actively maintains its own shape. After every insert and delete it checks whether it's growing lopsided and, if so, rearranges itself. The result is a hard guarantee: height stays O(log n), so search, insert and delete are O(log n) in the **worst** case, not just on average.

## The problem, one more time

Two trees, the same seven keys, built by inserting in different orders:

```c
#include <stdio.h>
#include <stdlib.h>

typedef struct Node {
    int key;
    struct Node *left;
    struct Node *right;
} Node;

Node *insert(Node *node, int key) {
    if (node == NULL) {
        Node *created = malloc(sizeof *created);
        if (created == NULL) {
            return NULL;
        }
        *created = (Node){key, NULL, NULL};
        return created;
    }
    if (key < node->key) {
        Node *left = insert(node->left, key);
        if (left != NULL) node->left = left;
    } else if (key > node->key) {
        Node *right = insert(node->right, key);
        if (right != NULL) node->right = right;
    }
    return node;
}

int height(const Node *node) {
    if (node == NULL) {
        return -1;                 // empty tree has height -1, a leaf has 0
    }
    int left = height(node->left);
    int right = height(node->right);
    return 1 + (left > right ? left : right);
}

void free_tree(Node *node) {
    if (node == NULL) return;
    free_tree(node->left);
    free_tree(node->right);
    free(node);
}

int main(void) {
    int sorted_keys[] = {1, 2, 3, 4, 5, 6, 7};
    int bushy_keys[] = {4, 2, 6, 1, 3, 5, 7};
    Node *sorted_tree = NULL;
    Node *bushy_tree = NULL;

    for (size_t i = 0; i < 7; i++) sorted_tree = insert(sorted_tree, sorted_keys[i]);
    for (size_t i = 0; i < 7; i++) bushy_tree = insert(bushy_tree, bushy_keys[i]);

    printf("%d\n", height(sorted_tree));   // 6
    printf("%d\n", height(bushy_tree));    // 2
    free_tree(sorted_tree);
    free_tree(bushy_tree);
    return 0;
}
```

Same data, same code, same ordering rule satisfied in both. One tree answers a search in up to 7 comparisons; the other in 3. Scale that to a million keys and it's a million comparisons versus twenty.

:::key
The BST rule fixes *which side* a key goes on. It says nothing about **height**. Balanced trees add a second rule that constrains height, and enforce it on every write.
:::

## Rotations: rearranging without breaking the rule

The tool for fixing shape is a **rotation**. It's a small, local rewiring of three pointers that changes the height of a section of the tree while preserving the BST ordering exactly.

```text
 right rotation at 30

        30                          20
       /  \                        /  \
     20    D        -->          10    30
    /  \                        /  \   / \
  10    C                      A    B C   D
  / \
 A   B
```

Look carefully at where the subtrees land. Before, reading the tree in order gives A, 10, B, 20, C, 30, D. After, it gives A, 10, B, 20, C, 30, D — identical. The sorted order is untouched; only the depths changed. The left side got one level shorter and the right side got one level taller.

A **left rotation** is the same move mirrored. Every balanced-tree algorithm is built from these two operations, applied at the right moments.

:::analogy
A rotation is like lifting a mobile hanging from the ceiling and re-hooking it one branch higher. The same ornaments hang in the same left-to-right order; the arrangement just sits differently.
:::

We won't write rotation code here — the pointer bookkeeping is fiddly and it's the *idea* that matters. What you need to carry forward is: rotations are O(1), they preserve the ordering rule, and they trade height between the two sides of a node.

## AVL trees: strict balance

An **AVL tree** adds this rule: for every node, the heights of its left and right subtrees differ by at most 1.

After each insert or delete the tree walks back up toward the root, checking that condition at each node. Where it fails, one or two rotations restore it. Both operations stay O(log n) because the walk-up is bounded by the height.

The rule is strict, which has a clear consequence: AVL trees stay very short — noticeably shorter than the alternatives — so **lookups are fast**. The price is that strictness triggers rebalancing more often, so **inserts and deletes do more rotation work**.

## Red-black trees: looser balance, fewer rotations

A **red-black tree** takes a different bargain. Each node is painted red or black, and a handful of colouring rules (the root is black, a red node's children are black, and every path from a node down to a leaf passes through the same number of black nodes) combine to guarantee that the longest path is at most twice the shortest.

That's a looser guarantee than AVL's — the tree can be lumpier — but it's still O(log n) height, which is all the complexity promise needs. In exchange, restoring the rules after a write usually needs only a recolouring, and never more than a small constant number of rotations.

```text
 property            AVL                    red-black
 --------------------------------------------------------------
 balance rule        subtree heights        longest path <= 2x
                     differ by <= 1         shortest path
 typical height      shorter                taller
 lookup              faster                 fast
 insert/delete       more rotations         fewer rotations
 good fit            read-heavy workloads   write-heavy workloads
 all operations      O(log n) worst case    O(log n) worst case
```

Both are O(log n) for everything. Choosing between them is choosing where you'd rather spend your time — and red-black trees' cheaper writes are why they show up so often in standard libraries (Java's `TreeMap` and many C++ `std::map` implementations are red-black trees).

:::warning
Neither AVL nor red-black keeps the tree *perfectly* balanced — that would cost too much to maintain. They keep it "balanced enough", which is precisely the point: the guarantee you need is O(log n) height, not a flawless shape.
:::

## 2-3 and 2-3-4 trees: another way to stay level

There's a completely different strategy for staying balanced: instead of rotating, let nodes hold **more than one key**.

In a **2-3 tree**, every node is either a *2-node* (one key, two children) or a *3-node* (two keys, three children).

```text
        [ 20 | 40 ]              a 3-node: two keys, three children
        /     |     \
   [10]   [25|30]   [50]

  keys < 20    20..40    > 40
```

Growth works upward instead of downward. When you insert into a full node it **splits**, pushing its middle key up into the parent — and if that fills the parent, it splits too. New levels are only ever created at the root, which means **every leaf is always at exactly the same depth**. Perfect balance, for free, with no rotations at all.

A **2-3-4 tree** is the same idea with room for up to three keys and four children per node. It's worth knowing because a 2-3-4 tree and a red-black tree are the same structure in disguise: a red-black tree is a binary encoding of a 2-3-4 tree, where red nodes represent keys glued into their black parent.

:::example
Insert 35 into the 3-node `[25|30]` above and it momentarily holds 25, 30, 35. It splits: 25 and 35 become separate nodes, and 30 is pushed up into the parent, which becomes `[20 | 30 | 40]` — and if that overflows, it splits in turn.
:::

## B-trees: why databases love wide nodes

Now push the 2-3-4 idea much further. A **B-tree** is a search tree where each node holds *hundreds* of keys and has *hundreds* of children. All leaves stay at the same depth, exactly as in a 2-3 tree.

Why would anyone want such wide nodes? Because of where the data lives. Reading from a disk or an SSD happens in fixed-size blocks — you can't fetch one key, you fetch a whole page of several kilobytes. And a single one of those fetches is thousands of times slower than any comparison the CPU does once the page is in memory.

So the cost that matters isn't comparisons. It's **how many blocks you have to fetch**, which is the height of the tree. B-trees are designed so one node fills exactly one disk block: you pay for one fetch and get hundreds of keys out of it.

```text
 storing 1,000,000 keys

 balanced binary tree   height ~ log2(1,000,000)   ~= 20 levels -> ~20 fetches
 B-tree, fan-out 100    height ~ log100(1,000,000)  =   3 levels -> ~3 fetches
```

Twenty disk reads versus three, for the same data. That's the entire reason B-trees (and the closely related B+ trees) sit underneath practically every database index and many filesystems.

:::key
High **fan-out** means low height. When each step down the tree costs a slow disk fetch, minimising the number of steps matters far more than minimising comparisons — and that's what a B-tree optimises for.
:::

That's four different balancing strategies, and you will almost certainly never write one from scratch.

:::tip
You rarely implement these yourself. What you use daily is the result: a conventional C hash table uses buckets, but a database index is almost certainly a B-tree, and many ordered-map libraries use red-black trees. Knowing which is which tells you what the thing is fast at.
:::

## Check Your Understanding

:::quiz
Q: What do balanced search trees guarantee that a plain BST does not?
- That lookups are O(1)
- That the height stays O(log n) even in the worst case *
- That the keys are stored in an array
- That no rotations are ever needed
E: A plain BST can degenerate to height n - 1. Balanced trees enforce a shape rule on every write, keeping height O(log n) and therefore all operations O(log n) in the worst case.
:::

:::quiz
Q: Why do databases use B-trees with hundreds of keys per node instead of binary trees?
- Because comparisons are expensive
- Because a wide node means a shorter tree and therefore far fewer slow disk fetches *
- Because B-trees don't need to be sorted
- Because binary trees can't store more than 1,000 keys
E: Disk reads happen in blocks and dominate the cost. High fan-out shrinks the height, so a lookup touches only a handful of blocks.
:::

:::match
Q: Match each tree to its defining idea.
- AVL | Subtree heights differ by at most 1
- Red-black | Colouring rules keep the longest path within twice the shortest
- 2-3 tree | Nodes hold one or two keys and split upward
- B-tree | Very wide nodes sized to fit one disk block
E: All four guarantee O(log n) height; they differ in how strictly they enforce it and what hardware they're tuned for.
:::

## Talk about it

> AVL trees rebalance aggressively for faster reads; red-black trees rebalance lazily for faster writes. Think of two systems you use — say a dictionary app that's written once and read constantly, versus a live leaderboard updated every second. Which tree would you choose for each, and what would go wrong if you swapped them?

## What's next

Every tree so far has compared whole keys at each node. The next one doesn't compare at all — it walks the key one character at a time, so the cost depends on the *length* of the word rather than on how many words are stored. In **Tries**, you'll see the structure behind autocomplete, and why "find everything starting with 'pix'" is a question hash tables simply cannot answer.

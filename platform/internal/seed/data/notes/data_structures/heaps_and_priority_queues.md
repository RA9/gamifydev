# Heaps and Priority Queues

A queue serves whoever arrived first. But a hospital doesn't work that way, and neither does an operating system deciding which task to run next. Sometimes what you need is: *give me the most urgent item, whenever I ask*.

That's a **priority queue**, and the structure that usually implements it is a **heap** — a binary tree with a beautifully weak ordering rule, packed into a plain array with no pointers at all. It's one of the most satisfying structures in the course.

## The heap property

A **min-heap** obeys one rule:

> Every node's key is less than or equal to the keys of its children.

That's all. Note how much weaker this is than the BST rule — it says nothing about left versus right, only about parent versus child.

```text
            1
          /   \
         3     6
        / \   /
       5   9 8
```

Check it: 1 ≤ 3 and 1 ≤ 6; 3 ≤ 5 and 3 ≤ 9; 6 ≤ 8. Valid. Notice that 5 sits to the left of 9 but 8 (smaller than 9) is over on the right — the heap simply doesn't care. It is *not* sorted.

What the rule does guarantee is the one thing we want: **the smallest key is at the root**. Reading it is one array access, so **peek is O(1)**.

A **max-heap** is the mirror image — every node is ≥ its children, so the largest key sits at the root. Everything below works the same with the comparison flipped.

:::analogy
A heap is a well-run tournament bracket. The champion is definitely at the top, and each match winner beat everyone below them. But you learn nothing about who is second best from the bracket alone — the runner-up could be anywhere in the loser's half.
:::

So a heap is deliberately *under*-organised. It knows the champion and nothing else, and that ignorance is exactly what makes it cheap to maintain.

:::key
The heap property is **parent before child**, nothing more. That weakness is a feature: it's cheap enough to restore after every change, while still keeping the minimum (or maximum) instantly available at the root.
:::

## A tree living inside an array

A heap is always kept as a **complete** binary tree — every level full except the last, which fills left to right. That's the definition from the previous lesson, and it's what makes the next trick possible.

Because there are never any gaps, you can number the nodes row by row and store them in a plain array. No node objects, no pointers, no `left` and `right` fields.

```text
            1(0)
          /     \
        3(1)    6(2)
        /  \    /
     5(3) 9(4) 8(5)

  array:  [ 1 , 3 , 6 , 5 , 9 , 8 ]
  index:    0   1   2   3   4   5
```

The parent-child relationships become arithmetic:

```text
  left child of i   =  2i + 1
  right child of i  =  2i + 2
  parent of i       =  (i - 1) // 2      (integer division)
```

Check it against the diagram. Node 3 is at index 1; its children should be at 2(1)+1 = 3 and 2(1)+2 = 4, holding 5 and 9. Correct. The parent of index 5 is (5−1)//2 = 2, holding 6. Correct.

:::key
A heap needs no pointers. Because it's a **complete** tree, the array index encodes the structure — `2i+1`, `2i+2`, `(i-1)//2`. That's why heaps are compact and cache-friendly.
:::

## Insert: append, then sift up

To insert, put the new key in the only place that keeps the tree complete — the end of the array — and then let it climb until the heap property holds again. Swapping a node with its parent while it's smaller is called **sifting up** (or bubbling up).

```python
def sift_up(heap, i):
    while i > 0:
        parent = (i - 1) // 2
        if heap[i] < heap[parent]:
            heap[i], heap[parent] = heap[parent], heap[i]
            i = parent
        else:
            break

heap = [1, 3, 6, 5, 9, 8]
heap.append(2)
sift_up(heap, len(heap) - 1)
print(heap)   # [1, 3, 2, 5, 9, 8, 6]
```

Trace it. The 2 lands at index 6; its parent is index 2, holding 6. Since 2 < 6 they swap. Now 2 is at index 2, whose parent is index 0, holding 1. Since 2 > 1 it stops. Two comparisons, one swap.

The climb can never take more steps than the tree is tall, and a complete tree of n nodes has height ⌊log₂ n⌋. So **insert is O(log n)**.

## Extract-min: swap, shrink, then sift down

Removing the root is the mirror image. You can't just delete index 0 — that would leave a hole. Instead, move the *last* element into the root (which keeps the tree complete), shrink the array by one, and let that element sink until the heap property holds. This is **sifting down**: repeatedly swap with the *smaller* of the two children.

```text
 start:  [1, 3, 6, 5, 9, 8]      extract the 1

 step 1: move last (8) to root -> [8, 3, 6, 5, 9]
 step 2: children of 8 are 3 and 6; smaller is 3 -> swap
                                 -> [3, 8, 6, 5, 9]
 step 3: children of 8 (index 1) are 5 and 9; smaller is 5 -> swap
                                 -> [3, 5, 6, 8, 9]
 done. new minimum 3 is at the root.
```

Like the climb, the sink is bounded by the height, so **extract-min is O(log n)**.

Python's standard library gives you all of this as `heapq`, which works directly on an ordinary list:

```python
import heapq

heap = []
for x in [5, 1, 9, 3]:
    heapq.heappush(heap, x)

print(heap[0])              # 1   -- peek, O(1)
print(heapq.heappop(heap))  # 1
print(heapq.heappop(heap))  # 3
print(heap)                 # [5, 9]
```

:::warning
Don't print a heap and expect a sorted list. `[1, 3, 2, 5, 9, 8, 6]` is a perfectly valid heap. Only the root is guaranteed to be in the right place — everything else is merely "not smaller than its parent".
:::

## Build-heap, and what it costs

If you already have all n items, you don't need n separate inserts. Start from the array as-is and sift *down* every internal node, working backwards from the last one. Surprisingly, this is **O(n)** overall — cheaper than the O(n log n) you'd pay for n inserts — because most nodes are near the bottom and have almost nowhere to sink.

```python
import heapq

nums = [9, 4, 7, 1, 8]
heapq.heapify(nums)    # O(n)
print(nums[0])         # 1
```

:::example
In a heap of 1,000 items, about 500 are leaves that can't sink at all, ~250 can sink at most one level, ~125 at most two. Adding that up stays proportional to n, not n log n.
:::

## Priority queue vs heap

These two names get used interchangeably, and it's worth separating them.

A **priority queue** is an *abstract data type* — a description of behaviour. It promises three operations: add an item with a priority, peek at the highest-priority item, and remove the highest-priority item. It says nothing about how.

A **heap** is a *concrete data structure* — an actual arrangement of memory. It happens to implement all three operations efficiently, which is why it's the usual choice.

```text
 implementation        insert      extract-min    peek
 ------------------------------------------------------
 unsorted array        O(1)        O(n)           O(n)
 sorted array          O(n)        O(1)           O(1)
 binary heap           O(log n)    O(log n)       O(1)
```

The unsorted array is great at inserting and terrible at extracting; the sorted array is the reverse. The heap refuses to be terrible at either. That balance is why it wins in practice.

Where you'll meet priority queues:

- **Schedulers.** An operating system picks the highest-priority runnable task, thousands of times a second.
- **Dijkstra's shortest-path algorithm.** It repeatedly needs "the unvisited node with the smallest known distance" — exactly extract-min.
- **Top-k problems.** To keep the 10 largest items from a huge stream, hold a min-heap of size 10: compare each new item to the root, and if it's bigger, replace the root and sift down. Memory stays constant no matter how long the stream is.

:::tip
`heapq` is a min-heap only. For a max-heap, push negated values (`heapq.heappush(h, -x)`) and negate again on the way out, or push tuples like `(-priority, item)` so the largest priority sorts first.
:::

## Check Your Understanding

:::quiz
Q: In an array-based binary heap, where do the children of index `i` live?
- i-1 and i+1
- 2i and 2i+1
- 2i+1 and 2i+2 *
- i/2 and i/2+1
E: With zero-based indexing the children of i are at 2i+1 and 2i+2, and the parent of i is at (i-1)//2.
:::

:::predict
Q: What does this print?
```python
import heapq
heap = []
for x in [4, 7, 2, 9]:
    heapq.heappush(heap, x)
print(heapq.heappop(heap), heap[0])
```
- 2 4 *
- 4 7
- 9 7
- 2 9
E: The heap always keeps its smallest value at the root, so `heappop` returns 2 and the next smallest, 4, rises to `heap[0]`.
:::

:::fill
Q: Complete the operation performed on a newly appended item to restore the heap property.
`sift_ ___ (heap, len(heap) - 1)`
- up *
- down
- out
E: A new item is appended at the end and climbs toward the root; sifting down is used on the root after an extraction.
:::

## Talk about it

> A heap keeps only a partial order — the root is guaranteed correct, and everything else is merely "not smaller than its parent". Explain in your own words why maintaining that weak rule is much cheaper than keeping the whole collection sorted, and name a situation where the weak rule would not be good enough.

## What's next

The heap controls its shape absolutely, which is why it never degenerates. A plain binary search tree has no such discipline — and we left it collapsed into a chain two lessons ago. Next up is **Balanced Search Trees**, where rotations, AVL and red-black trees, and the B-trees inside real databases all exist for one purpose: guaranteeing that the height stays O(log n) no matter what order the data arrives in.

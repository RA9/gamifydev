# Heap Sort

Merge sort guarantees O(n log n) but wants O(n) extra memory. Quick sort is in-place but only guarantees O(n log n) on average. Heap sort takes the remaining corner of that trade-off: **O(n log n) in the worst case, in place**. That combination is its whole reason for existing.

It's also the clearest example in this course of a data structure becoming an algorithm. You met heaps in the Data Structures course; heap sort is what happens when you notice that a heap will hand you the maximum element over and over, and that "repeatedly take the maximum" is a sort.

## Recalling the heap

A **binary heap** is a complete binary tree stored in a plain array. In a **max-heap**, every parent is greater than or equal to both of its children, so the largest value in the whole structure sits at the root — index 0.

Because the tree is complete (filled left to right, no gaps), the parent/child relationships are pure arithmetic. No pointers, no nodes:

```text
array:   [9, 7, 8, 3, 4, 2]

index i:      left child = 2i + 1      right child = 2i + 2

                  9  (index 0)
                /   \
        (1)   7       8   (2)
             / \     /
       (3)  3   4   2
           (4)     (5)
```

Check it: index 1 holds 7, its children are at indices 3 and 4 (values 3 and 4) — both smaller. Index 0 holds 9, children at 1 and 2 (7 and 8) — both smaller. That's a valid max-heap.

Notice the heap property says nothing about siblings. `7` sits to the left of `8` even though it's smaller. A heap is not a sorted array; it only promises that the maximum is on top.

:::key
A max-heap stored in an array gives you the maximum in O(1) — it's always at index 0. Everything about heap sort follows from repeatedly taking that element.
:::

## Sift down: the one operation you need

Suppose the heap property holds everywhere except at one node, whose value might be too small. **Sift down** (also called heapify-down) fixes it: compare the node with its children, swap it with the larger child if that child is bigger, and follow it down, repeating until it's in a valid spot or reaches the bottom.

```python
def sift_down(items, start, end):
    """Restore the max-heap property at index `start`, treating
    items[0:end] as the heap."""
    root = start
    while 2 * root + 1 < end:
        child = 2 * root + 1                # left child
        if child + 1 < end and items[child] < items[child + 1]:
            child += 1                      # right child is bigger, use it
        if items[root] < items[child]:
            items[root], items[child] = items[child], items[root]
            root = child                    # follow the value down
        else:
            return                          # heap property restored
```

The value travels at most the height of the tree, and a complete binary tree of n nodes has height ⌊log₂ n⌋. So **sift down is O(log n)**.

:::analogy
Sifting down is a manager who turns out to be less capable than a direct report. You swap them, then check the new team below — repeating until nobody underneath outranks the person you moved. The competent people bubble upwards; the misplaced one settles at its real level.
:::

## Building the heap — and a nice surprise

To turn an arbitrary array into a max-heap, sift down every node that has at least one child, starting from the **last** such node and working backwards to the root. Working backwards matters: by the time you sift down node `i`, both of its subtrees are already valid heaps, so a single sift-down is enough.

```python
def build_max_heap(items):
    n = len(items)
    for i in range(n // 2 - 1, -1, -1):     # last parent down to the root
        sift_down(items, i, n)
    return items

print(build_max_heap([3, 7, 2, 9, 4, 8]))   # [9, 7, 8, 3, 4, 2]
```

Now the surprise. It looks like n sift-downs at O(log n) each, so O(n log n) — but it's actually **O(n)**.

The reason is that most nodes are near the bottom, and a node near the bottom barely moves. In a complete tree, about half the nodes are leaves (height 0, no work at all), a quarter are at height 1 (at most one swap), an eighth at height 2, and so on. Summing "number of nodes at height h × h" over all heights gives a series that converges to a constant multiple of n rather than growing like n log n.

```text
height   nodes (approx)   max swaps each   total swaps
   0         n/2                0               0
   1         n/4                1              n/4
   2         n/8                2              n/4
   3        n/16                3             3n/16
   ...                                     ------------
                                    sums to less than 2n
```

Build-heap is O(n). It's a genuinely pleasant result: the expensive-looking setup is the cheap part.

:::key
Building a heap from an unsorted array is **O(n)**, not O(n log n), because the many cheap leaf-level nodes vastly outnumber the few expensive ones near the root.
:::

## The sort

With a max-heap built, the sort is a loop with a beautiful trick. The maximum is at index 0, and its correct final home is the last position. So swap them. Now the largest value is in place at the end — shrink the heap by one so you never touch it again, and sift the newly-arrived root back down.

```python
def heap_sort(items):
    n = len(items)
    for i in range(n // 2 - 1, -1, -1):     # phase 1: build the max-heap, O(n)
        sift_down(items, i, n)
    for end in range(n - 1, 0, -1):         # phase 2: n-1 extractions
        items[0], items[end] = items[end], items[0]   # max goes to its home
        sift_down(items, 0, end)            # restore the shrunken heap
    return items

print(heap_sort([3, 7, 2, 9, 4, 8]))   # [2, 3, 4, 7, 8, 9]
```

Trace phase 2 on the heap `[9, 7, 8, 3, 4, 2]`, with `|` marking the boundary between the live heap and the finished sorted tail:

```text
heap = [9, 7, 8, 3, 4, 2 |]

swap root with last  ->  [2, 7, 8, 3, 4 | 9]
sift 2 down          ->  [8, 7, 2, 3, 4 | 9]

swap root with last  ->  [4, 7, 2, 3 | 8, 9]
sift 4 down          ->  [7, 4, 2, 3 | 8, 9]

swap root with last  ->  [3, 4, 2 | 7, 8, 9]
sift 3 down          ->  [4, 3, 2 | 7, 8, 9]

swap root with last  ->  [2, 3 | 4, 7, 8, 9]
sift 2 down          ->  [3, 2 | 4, 7, 8, 9]

swap root with last  ->  [2, 3 | 4, 7, 8, 9]   heap size 1, done
final:                   [2, 3, 4, 7, 8, 9]
```

The sorted region grows from the right, and the heap shrinks from the right, sharing one array. That's why the sort is in place: the "extra" space for the output is the space the heap just vacated.

:::tip
The `end` parameter in `sift_down` is doing quiet but essential work: it's the fence between the live heap and the finished tail. Get it wrong and already-sorted values get pulled back into the heap and scrambled.
:::

## The complexity, and the honest downsides

**Time.** Phase 1 is O(n). Phase 2 does n-1 extractions, each with an O(log n) sift down, so it's O(n log n). Together: **O(n log n) in the best, average and worst case.** Like merge sort, the bound doesn't depend on the input's arrangement. Unlike merge sort, there's no adversarial input at all — heap sort has no equivalent of quick sort's sorted-array disaster.

**Space.** O(1) auxiliary. Every operation is a swap inside the original array, and both loops are iterative, so there's no recursion stack either.

**Stability.** Heap sort is **not stable**. The root-to-end swap moves values across long distances, freely jumping equal elements over each other.

**Speed in practice.** Here's the honest part: heap sort is usually *slower* than a well-implemented quick sort despite the better worst-case bound. Sift down jumps between indices `i`, `2i+1`, `4i+3` and so on — a scattered access pattern that the CPU cache handles poorly, while quick sort sweeps contiguously. Heap sort also performs more swaps than quick sort on typical data.

```text
Sort        Best        Average      Worst        Space   Stable
------------------------------------------------------------------
Merge       O(n log n)  O(n log n)   O(n log n)   O(n)    yes
Quick       O(n log n)  O(n log n)   O(n^2)       O(log n) no
Heap        O(n log n)  O(n log n)   O(n log n)   O(1)    no
```

So why does it exist? Because it's the only one of the three with a hard worst-case guarantee *and* O(1) space. That makes it the right choice in memory-constrained environments, in systems where a worst-case bound is a requirement rather than a preference, and — as you saw last lesson — as introsort's safety net when quick sort's recursion gets too deep.

:::example
An embedded controller with a few kilobytes of RAM cannot spare merge sort's O(n) buffer, and a real-time system cannot accept quick sort's small chance of O(n²). Heap sort fits both constraints exactly, which is why it turns up in places you'd never expect a "slower" sort.
:::

## Check Your Understanding

:::quiz
Q: What is heap sort's worst-case time complexity?
- O(n)
- O(n log n) *
- O(n²)
- O(log n)
E: Building the heap is O(n), then n-1 extractions each cost O(log n) for the sift down. The total is O(n log n) with no bad inputs.
:::

:::quiz
Q: Building a max-heap from an unsorted array of n elements costs:
- O(1)
- O(log n)
- O(n) *
- O(n log n)
E: About half the nodes are leaves and do no work at all, and only the few nodes near the root can sift far. Summing the work by height gives O(n), not O(n log n).
:::

:::reorder
Q: What is the first step of one iteration of heap sort's extraction phase?
- Swap the root (the maximum) with the last element of the live heap *
- Sift the new root down
- Shrink the heap boundary by one
- Rebuild the entire heap from scratch
E: Each iteration swaps the maximum into its final position first, then shrinks the boundary, then sifts the newly-promoted root down to restore the heap.
:::

## Talk about it

> Heap sort has a better worst-case guarantee than quick sort but is usually slower in practice, mostly because of how it jumps around in memory. Big O deliberately ignores that kind of hardware detail. When is Big O the right lens for choosing an algorithm, and when do you have to actually measure instead?

## What's next

That's the four classic comparison sorts done, and you've seen how the same problem admits genuinely different trade-offs. Next up is **Tree Traversals** — the algorithms for visiting every node of a tree, in several different useful orders. Heaps were trees in disguise; now you'll work with real branching structures, and you'll find that recursion does most of the work for you.

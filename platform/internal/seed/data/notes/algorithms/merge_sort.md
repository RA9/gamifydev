# Merge Sort

The simple sorts were all O(n²) because they moved data one small step at a time. Merge sort breaks out of that by using a completely different strategy: split the problem in half, solve both halves, then combine. That strategy has a name — **divide and conquer** — and it's one of the most productive ideas in all of algorithms.

Merge sort gives you O(n log n) in *every* case, not just on average, and it's stable. In this lesson you'll build it from the merge step upwards and see exactly where the log n comes from.

## Divide and conquer

Divide and conquer has three moves:

1. **Divide** — break the problem into smaller sub-problems of the same shape.
2. **Conquer** — solve each sub-problem, usually by recursing until they're trivially small.
3. **Combine** — assemble the sub-answers into an answer for the whole.

For sorting, that reads: cut the list in half, sort each half, then merge the two sorted halves into one sorted list. The base case is a list of zero or one element, which is already sorted by definition.

The clever part isn't the dividing — anyone can cut a list in half. It's that **merging two already-sorted lists is cheap**, and that's what makes the whole thing work.

:::analogy
Imagine two card players each holding a sorted hand, face up. To combine them into one sorted pile you never search: you just compare the two top cards and take the smaller one, over and over. Each comparison places exactly one card.
:::

## The merge step

This is the heart of the algorithm, so let's build it alone first. Given two sorted lists, walk a pointer along each. Repeatedly take the smaller of the two front elements. When one list runs out, everything remaining in the other is already sorted and larger, so append it wholesale.

```python
def merge(left, right):
    result = []
    i = j = 0
    while i < len(left) and j < len(right):
        if left[i] <= right[j]:        # <= keeps the sort stable
            result.append(left[i])
            i += 1
        else:
            result.append(right[j])
            j += 1
    result.extend(left[i:])            # whatever is left over
    result.extend(right[j:])
    return result

print(merge([1, 4, 9], [2, 3, 10]))   # [1, 2, 3, 4, 9, 10]
```

Trace it:

```text
left  = [1, 4, 9]      right = [2, 3, 10]
i=0 j=0   1 <= 2   take 1     result=[1]
i=1 j=0   4 >  2   take 2     result=[1,2]
i=1 j=1   4 >  3   take 3     result=[1,2,3]
i=1 j=2   4 <= 10  take 4     result=[1,2,3,4]
i=2 j=2   9 <= 10  take 9     result=[1,2,3,4,9]
i=3       left exhausted -> extend with right[2:] = [10]
                             result=[1,2,3,4,9,10]
```

Every comparison consumes exactly one element, and there are `len(left) + len(right)` elements total, so merging is **O(n)** where n is the combined length.

:::key
Merging two sorted lists of total length n takes O(n) time. That single fact is what buys merge sort its O(n log n) — the expensive-looking combine step is actually linear.
:::

## The full algorithm

With `merge` in hand, the sort itself is four lines:

```python
def merge_sort(items):
    if len(items) <= 1:                    # base case: already sorted
        return items
    mid = len(items) // 2
    left = merge_sort(items[:mid])         # conquer the left half
    right = merge_sort(items[mid:])        # conquer the right half
    return merge(left, right)              # combine

print(merge_sort([5, 2, 8, 1, 9, 3]))   # [1, 2, 3, 5, 8, 9]
```

Here's the whole call tree for `[5, 2, 8, 1]` — splitting on the way down, merging on the way back up:

```text
                [5, 2, 8, 1]
               /            \            split
        [5, 2]                [8, 1]
        /    \                /    \     split
     [5]      [2]          [8]      [1]  base cases
        \    /                \    /
        [2, 5]                [1, 8]     merge
               \            /
                [1, 2, 5, 8]             merge
```

Notice the shape: going down, nothing is compared. All the actual sorting work happens on the way back up, in the merges.

:::tip
When you're debugging a divide-and-conquer function, test the combine step in isolation first. If `merge` is right and the base case is right, the recursion almost always is too. Most merge sort bugs live in `merge`.
:::

## Why it's O(n log n) in every case

Look at that diagram again in terms of *levels*.

**How many levels?** Each level halves the sub-list size: n, n/2, n/4, ... down to 1. The number of halvings it takes to get from n to 1 is log₂(n). So there are about **log n levels**.

**How much work per level?** At every level, the merges together touch each of the n elements exactly once. So each level costs **O(n)**.

```text
level 0:   one merge of size 8          -> 8 elements touched
level 1:   two merges of size 4         -> 8 elements touched
level 2:   four merges of size 2        -> 8 elements touched
                                           ---------------
           log2(8) = 3 levels x O(n)    =  O(n log n)
```

Multiply: log n levels × O(n) per level = **O(n log n)**.

And here's the important part — nothing in that argument depended on the data. The split is always down the middle, and the merge always touches every element, no matter what order the input was in. So merge sort is O(n log n) **best, average and worst**. There is no bad input. That predictability is its signature strength, and it's exactly what quick sort cannot promise.

:::key
Merge sort is O(n log n) in the best, average *and* worst case. Its running time depends only on the size of the input, never on its arrangement.
:::

## The cost: O(n) extra space

Merge sort's weakness is memory. `merge` builds a brand-new `result` list, so at any moment you're holding a copy of the data being merged. The peak extra space is **O(n)**.

Compare that with insertion sort or heap sort, which rearrange the array in place at O(1) extra space. If you're sorting an array that barely fits in memory, that doubling is a genuine problem.

You can reduce the copying with a more careful implementation that allocates one scratch buffer up front and merges into it, rather than slicing at every level. But you can't get merge sort down to O(1) auxiliary space without giving up its simplicity — true in-place merging exists, and it is genuinely difficult.

:::warning
The version above also slices with `items[:mid]` and `items[mid:]`, which copies at every level. That's fine for learning and for modest inputs, but a production implementation passes indices into a shared array and reuses a single scratch buffer.
:::

## Where merge sort shines

**Stability.** Because `merge` takes from `left` when the two fronts are *equal* (`left[i] <= right[j]`), elements that compare equal keep their original order. Flip that `<=` to `<` and you silently lose stability — a one-character bug worth knowing about.

**Linked lists.** Merge sort is the natural sort for linked lists. Merging two sorted linked lists needs only pointer rewiring, no extra array at all, so the O(n) space penalty vanishes. And linked lists have no random access, which rules out quick sort's partitioning and heap sort's index arithmetic.

**External sorting.** When the data is far too big for memory — sorting a file of billions of records — you sort chunks that do fit, write each sorted chunk to disk, then merge the chunks together by streaming. Merging only ever needs the *front* of each run in memory, so it works beautifully on data you can only read sequentially.

:::example
To sort a 500 GB file on a machine with 8 GB of RAM: read 4 GB, sort it in memory, write it out; repeat to get 125 sorted files; then merge all 125 by repeatedly taking the smallest current line across the files. Merge sort is the only classic sort that adapts to this naturally.
:::

## Check Your Understanding

:::quiz
Q: What is merge sort's worst-case time complexity?
- O(n)
- O(n log n) *
- O(n²)
- O(log n)
E: Merge sort always splits down the middle and always merges in linear time, so its cost is log n levels × O(n) per level regardless of the input arrangement.
:::

:::predict
Q: What does this print?
```python
def merge(left, right):
    result = []
    i = j = 0
    while i < len(left) and j < len(right):
        if left[i] <= right[j]:
            result.append(left[i]); i += 1
        else:
            result.append(right[j]); j += 1
    result.extend(left[i:]); result.extend(right[j:])
    return result

print(merge([2, 6], [1, 3, 4]))
```
- [1, 2, 3, 4, 6] *
- [2, 6, 1, 3, 4]
- [1, 2, 3, 6, 4]
- [6, 4, 3, 2, 1]
E: Compare fronts: 1 < 2 take 1, then 2 <= 3 take 2, then 3 < 6 take 3, then 4 < 6 take 4, then left is what remains — append 6.
:::

:::quiz
Q: What does merge sort pay for its guaranteed O(n log n)?
- It is not stable
- It needs O(n) extra memory *
- It only works on numbers
- It requires sorted input
E: The merge step builds a new list, so merge sort uses O(n) auxiliary space — the main trade-off against in-place sorts like heap sort and quick sort.
:::

## Talk about it

> Merge sort guarantees O(n log n) but needs O(n) extra memory. Quick sort, which you'll meet next, is in-place but can degrade to O(n²) on unlucky input. Describe a system where you'd insist on the guaranteed bound even at the cost of memory — and one where you'd happily take the risk.

## What's next

You've now met divide and conquer, and you'll see its shape again and again. Next up is **Quick Sort**, which divides differently: instead of splitting blindly down the middle and doing the work in the merge, it does the work up front by partitioning around a pivot — so there's nothing to combine at all. It's in-place, it's usually faster than merge sort in practice, and it has a worst case you need to know about.

# Quick Sort

Merge sort splits the list blindly and does its real work while combining. Quick sort inverts that: it does the work up front by **partitioning** the data around a chosen value, and then there is nothing to combine at all. That inversion makes it in-place and, in practice, usually the fastest general-purpose comparison sort.

It also gives quick sort a worst case that merge sort doesn't have. This lesson covers both sides honestly: how partitioning works, why the average is O(n log n), why the worst case is O(n²), and what real implementations do about it.

## Pivot and partition

Pick one element and call it the **pivot**. Rearrange the array so that everything smaller than the pivot ends up to its left, and everything larger ends up to its right. The pivot is now sitting in its **final sorted position** — nothing will ever move it again.

```text
[7, 2, 9, 4, 1, 8]      pivot = 8 (last element)

after partitioning:
[7, 2, 4, 1, | 8 | 9]
 all < 8       ^    all > 8
            final position of 8
```

Then recurse on the left part and the right part. When both are sorted, the whole array is sorted — with no merge step, because the pieces were already in the right regions.

:::key
Partitioning places the pivot in its **final** position and guarantees everything left of it is smaller and everything right is larger. That's why quick sort needs no combine step: the recursion sorts two regions that are already correctly separated.
:::

## The Lomuto partition scheme

There are two classic partition schemes. **Lomuto's** is the easier one to read, so we'll use it. It picks the last element as pivot and sweeps left to right, maintaining a boundary index `i` that marks the end of the "smaller than pivot" region.

```python
def partition(items, low, high):
    pivot = items[high]                       # pivot = last element
    i = low - 1                               # end of the "smaller" region
    for j in range(low, high):
        if items[j] <= pivot:
            i += 1
            items[i], items[j] = items[j], items[i]
        # else: leave items[j] in the "larger" region
    items[i + 1], items[high] = items[high], items[i + 1]   # pivot into place
    return i + 1                              # pivot's final index
```

Trace it on `[7, 2, 9, 4, 1, 8]` with `low=0, high=5`, so `pivot = 8`:

```text
start   i=-1   [7, 2, 9, 4, 1, 8]
j=0  7<=8  i=0  swap items[0],items[0]  [7, 2, 9, 4, 1, 8]
j=1  2<=8  i=1  swap items[1],items[1]  [7, 2, 9, 4, 1, 8]
j=2  9>8        no swap                 [7, 2, 9, 4, 1, 8]
j=3  4<=8  i=2  swap items[2],items[3]  [7, 2, 4, 9, 1, 8]
j=4  1<=8  i=3  swap items[3],items[4]  [7, 2, 4, 1, 9, 8]
end  swap items[4], items[5]            [7, 2, 4, 1, 8, 9]
     return 4  -> 8 is final at index 4
```

The invariant that makes this work: everything in `items[low..i]` is `<= pivot`, and everything in `items[i+1..j-1]` is `> pivot`.

The other scheme, **Hoare's**, walks two pointers inwards from both ends and swaps out-of-place pairs when they meet. It performs fewer swaps on average and is what most optimised libraries use, but its index bookkeeping is fiddlier, so Lomuto is the better one to learn on.

## The full sort

```python
def quick_sort(items, low=0, high=None):
    if high is None:
        high = len(items) - 1
    if low < high:                             # base case: 0 or 1 element
        p = partition(items, low, high)
        quick_sort(items, low, p - 1)          # sort the left region
        quick_sort(items, p + 1, high)         # sort the right region
    return items

print(quick_sort([7, 2, 9, 4, 1, 8]))   # [1, 2, 4, 7, 8, 9]
print(quick_sort([3, 3, 1]))            # [1, 3, 3]
```

Note `p - 1` and `p + 1`: the pivot itself is excluded from both recursive calls, because it's already final. Forgetting that and passing `p` into a recursive call gives you a sub-problem that never shrinks — an infinite recursion.

:::warning
`quick_sort(items, low, p)` instead of `quick_sort(items, low, p - 1)` is the classic quick sort bug. The pivot never leaves the range, the sub-problem never gets smaller, and you get a `RecursionError` instead of a sort.
:::

## Average case vs worst case

The whole story of quick sort is how evenly the pivot splits the array.

**When the pivot lands near the middle**, each partition halves the problem. That's log n levels of recursion, and each level does O(n) work partitioning — the same arithmetic as merge sort. **O(n log n) average case.**

```text
good pivots (median-ish):        bad pivots (always smallest):

        n                              n
      /   \                             \
   n/2     n/2                          n-1
   / \     / \                            \
 n/4 n/4 n/4 n/4                          n-2
  log n levels                              \
  O(n log n)                                ...
                                  n levels, O(n^2)
```

**When the pivot is always the smallest or largest element**, one side gets everything and the other gets nothing. The array shrinks by only one element per level, so there are n levels of O(n) work: **O(n²) worst case.**

And here is the part that catches people out. With Lomuto's "pivot = last element" rule, the input that triggers the worst case is **an already-sorted array**. Every pivot is the maximum of its range, every partition is maximally lopsided.

```text
[1, 2, 3, 4, 5]   pivot 5 -> [1,2,3,4 | 5]      left size 4
[1, 2, 3, 4]      pivot 4 -> [1,2,3 | 4]        left size 3
[1, 2, 3]         pivot 3 -> [1,2 | 3]          left size 2
...
n levels, O(n^2) total — on the nicest-looking input imaginable
```

That's genuinely alarming: sorted or nearly-sorted data is extremely common in the real world. Recursion depth also becomes O(n), so a large sorted array can blow the call stack too.

:::warning
Naive quick sort is O(n²) on **already-sorted input**, which is one of the most common inputs a program ever sees. Never ship a fixed-position pivot.
:::

## Fixing the pivot

The fix is not to make the worst case impossible — for any deterministic pivot rule, some adversarial input exists — but to make it astronomically unlikely.

**Randomised pivot.** Pick a random index and swap it to the pivot position before partitioning.

```python
import random

def partition_random(items, low, high):
    r = random.randint(low, high)
    items[r], items[high] = items[high], items[r]   # random element becomes pivot
    return partition(items, low, high)
```

Now no particular input is bad; the worst case requires an unlucky *sequence of coin flips*, and the probability of that is vanishingly small for any reasonable n. Expected time is O(n log n) for every input.

**Median-of-three.** Look at the first, middle, and last elements, and use their median as the pivot. This costs three comparisons and makes sorted and reverse-sorted input — the two common pathological cases — into best cases instead. Most library implementations use this or a larger sample.

Real implementations go further still: **introsort** starts as quick sort, counts its recursion depth, and if the depth exceeds about 2·log n it switches to heap sort, which has a guaranteed O(n log n) bound. That gives you quick sort's speed with heap sort's guarantee.

:::tip
Combine both habits: randomise or median-of-three for the pivot, and recurse into the *smaller* partition first while looping on the larger. That caps stack depth at O(log n) even when the splits are uneven.
:::

## In-place, not stable, and fast anyway

**In-place.** Partitioning only swaps elements within the original array. The only extra memory is the recursion stack: O(log n) with good pivots, O(n) in the worst case. Merge sort's O(n) auxiliary array is always there.

**Not stable.** Partition swaps elements across long distances, which can easily reorder equal values. If you need stability, use merge sort — or sort by a tuple that includes the original index.

**Usually faster than merge sort in practice**, despite the worse bound. Three reasons: it does no allocation, so it never pays for memory management; it works entirely within one contiguous array, so its sequential sweeps make excellent use of the CPU cache; and its inner loop is a comparison and a swap, with a smaller constant factor than merge's copying.

```text
Sort         Best        Average      Worst       Space        Stable
-----------------------------------------------------------------------
Merge        O(n log n)  O(n log n)   O(n log n)  O(n)         yes
Quick        O(n log n)  O(n log n)   O(n^2)      O(log n)*    no

* stack space with good pivots; O(n) in the worst case
```

:::key
Quick sort: O(n log n) **average**, O(n²) **worst**, in-place, not stable. With a randomised or median-of-three pivot the worst case effectively disappears, and its cache behaviour makes it the fastest general sort in practice.
:::

## Check Your Understanding

:::quiz
Q: What input makes naive quick sort (pivot = last element) hit its O(n²) worst case?
- A list of all identical values, with a three-way partition
- An already-sorted list *
- A list in random order
- A list of length 1
E: With the last element as pivot, a sorted list means every pivot is the maximum of its range, so each partition peels off just one element and the recursion runs n levels deep.
:::

:::fill
Q: Complete the recursive call so the pivot is excluded from the left partition.
`quick_sort(items, low, ___)`
- p - 1 *
- p
- p + 1
- high
E: The pivot at index `p` is already in its final position, so the left region ends at `p - 1`. Passing `p` leaves the sub-problem the same size and the recursion never terminates.
:::

:::quiz
Q: Why is quick sort often faster than merge sort in real programs?
- It has a better worst-case bound
- It is stable while merge sort is not
- It sorts in place with good cache behaviour and no allocation *
- It needs fewer comparisons in the worst case
E: Merge sort's bound is better. Quick sort wins on constant factors: no auxiliary array to allocate or copy into, and tight sequential passes over one contiguous block of memory.
:::

## Talk about it

> Randomising the pivot doesn't remove quick sort's O(n²) worst case — it just makes hitting it depend on chance rather than on the input. Is an algorithm that is "almost certainly fast" acceptable in software you'd trust with something important? Describe a system where you would insist on a guaranteed bound instead.

## What's next

You now have one sort with a guarantee but extra memory (merge), and one that's in-place and fast but only guaranteed on average (quick). Next up is **Heap Sort**, which gets both: O(n log n) in the *worst* case **and** in-place. It's built directly on the max-heap you met in the Data Structures course, so it's a chance to see a data structure turn into an algorithm.

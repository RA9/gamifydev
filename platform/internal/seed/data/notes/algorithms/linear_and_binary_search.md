# Linear and Binary Search

Searching is the most common thing programs do: find the user with this id, find the score in this list, check whether this word exists. There are two fundamental ways to search a sequence, and the difference between them is one of the sharpest illustrations of why algorithms matter.

In this lesson you'll implement both, see exactly what binary search demands in exchange for its speed, meet the off-by-one bugs that trip up almost everyone, and learn when it's worth sorting your data first just so you can search it faster.

## Linear search: check everything

**Linear search** walks the sequence from one end, comparing as it goes, and stops when it finds a match or runs out of items.

```c
#include <stddef.h>
#include <stdio.h>

ptrdiff_t linear_search(const int items[], size_t n, int target) {
    for (size_t i = 0; i < n; ++i) {
        if (items[i] == target) return (ptrdiff_t)i;
    }
    return -1;
}

int main(void) {
    int nums[] = {7, 3, 9, 1, 4};
    size_t n = sizeof nums / sizeof nums[0];
    printf("%td\n", linear_search(nums, n, 9)); /* 2 */
    printf("%td\n", linear_search(nums, n, 5)); /* -1 */
    return 0;
}
```

Best case it's the first item — one comparison, O(1). Worst case the target is last or absent, so you touch every element: **O(n) worst case, O(n) average**. Space is O(1); you allocate nothing.

The great virtue of linear search is that it asks for nothing. The data doesn't need to be sorted, doesn't need to be indexable, doesn't even need to be in memory all at once — it works on a stream, a linked list, a file being read line by line. Any collection you can iterate, you can linear search.

:::analogy
Linear search is looking for your friend at a party by walking up to every person in turn. Slow if the room is packed, but it works no matter how the guests are arranged.
:::

## Binary search: halve the problem

If the sequence is **sorted**, you can do far better. Look at the middle element. If it's your target, done. If it's too big, the answer can only be in the left half. If it's too small, only in the right half. Either way you've eliminated half the data with one comparison, and then you repeat.

```c
#include <stddef.h>
#include <stdio.h>

ptrdiff_t binary_search(const int items[], size_t n, int target) {
    ptrdiff_t low = 0;
    ptrdiff_t high = (ptrdiff_t)n - 1;
    while (low <= high) {
        ptrdiff_t mid = low + (high - low) / 2;
        if (items[mid] == target) return mid;
        if (items[mid] < target) low = mid + 1;  /* Target is to the right. */
        else high = mid - 1;                     /* Target is to the left. */
    }
    return -1;
}

int main(void) {
    int nums[] = {1, 3, 4, 7, 9, 12, 15};
    size_t n = sizeof nums / sizeof nums[0];
    printf("%td\n", binary_search(nums, n, 12)); /* 5 */
    printf("%td\n", binary_search(nums, n, 2));  /* -1 */
    return 0;
}
```

Trace the search for `12`:

```text
index:   0   1   2   3   4    5    6
value:   1   3   4   7   9   12   15

low=0 high=6  mid=3  items[3]=7   7 < 12  -> go right, low=4
low=4 high=6  mid=5  items[5]=12  match!  -> return 5
```

Two comparisons instead of six. Each step throws away half the remaining range, so from `n` items you need about log₂(n) steps: **O(log n) worst case and average**, O(1) space for the loop version. For a million sorted items, that's about 20 comparisons.

:::key
Binary search is O(log n) — but only on **sorted, randomly-accessible** data. Run it on an unsorted list and it won't error; it will confidently return the wrong answer. The precondition is the whole deal.
:::

## The classic pitfalls

Binary search is famously easy to get subtly wrong. Here are the four traps, all of which live in three lines of code.

**1. The loop condition.** It must be `while low <= high`, not `<`. With `<`, a range that has narrowed to a single element (`low == high`) is never examined, so you miss targets that sit at the very edge.

**2. Not moving the bounds past mid.** Writing `low = mid` instead of `low = mid + 1` means that when `low` and `high` are adjacent, `mid` computes to `low` again and nothing changes — an infinite loop. You must always exclude the element you just checked.

**3. Midpoint overflow.** The obvious `mid = (low + high) // 2` can overflow in languages with fixed-width integers, because `low + high` may exceed the maximum int even though `mid` itself would fit. The safe form is:

```c
ptrdiff_t mid = low + (high - low) / 2;
```

C integers have fixed widths, so the overflow bug is real. This form avoids adding two potentially large positive indices and was adopted only after the obvious form had survived in widely used library code for years.

**4. Unsorted input.** No amount of careful index arithmetic saves you if the precondition is violated.

:::warning
`while low < high` and `low = mid` are the two most common binary search bugs — the first silently misses edge elements, the second hangs forever. Test every search against an array of size 1 and size 2, and against a target that is smaller than everything and larger than everything.
:::

## The recursive version

Binary search is naturally recursive: searching a half is the same problem on a smaller range.

```c
#include <stddef.h>
#include <stdio.h>

ptrdiff_t binary_search_range(const int items[], int target,
                              ptrdiff_t low, ptrdiff_t high) {
    if (low > high) return -1;            /* Base case: empty range. */
    ptrdiff_t mid = low + (high - low) / 2;
    if (items[mid] == target) return mid;
    if (items[mid] < target) {
        return binary_search_range(items, target, mid + 1, high);
    }
    return binary_search_range(items, target, low, mid - 1);
}

ptrdiff_t binary_search_rec(const int items[], size_t n, int target) {
    return binary_search_range(items, target, 0, (ptrdiff_t)n - 1);
}

int main(void) {
    int items[] = {1, 3, 4, 7, 9, 12, 15};
    size_t n = sizeof items / sizeof items[0];
    printf("%td\n", binary_search_rec(items, n, 1));  /* 0 */
    printf("%td\n", binary_search_rec(items, n, 20)); /* -1 */
    return 0;
}
```

Both recursive calls are tail calls, so this is one of those functions that converts straight back into the loop above. Stack depth is only O(log n) — about 20 frames for a million items — so unlike a linear recursion, this one is unlikely to exhaust a typical C thread's stack.

:::tip
Never allocate and copy a sub-array for each recursive call. That costs O(n) work across the copies and destroys the O(log n) advantage. Pass the original array plus index bounds, as above.
:::

## When is sorting first worth it?

Sorting costs O(n log n) — you'll see why in the next few lessons. So sorting *just* to run one search is a bad trade: O(n log n) to set up plus O(log n) to search, versus O(n) for a single linear scan.

The maths flips when you search repeatedly. For `k` searches over `n` items:

```text
linear search, k times   ->  O(k * n)
sort once, then search   ->  O(n log n + k log n)
```

With n = 1,000,000, one search favours the linear scan. A thousand searches favours sorting overwhelmingly. The rule of thumb: **sort when the data is stable and you'll query it many times**; scan when you'll look once, or when the data changes constantly.

:::example
C's standard library provides `bsearch` for searching a sorted array and `qsort` for sorting one. But if you're doing pure membership checks with no need for order, a hash set gives O(1) average lookups and beats both — reach for the right data structure before reaching for the clever algorithm.
:::

## Check Your Understanding

:::quiz
Q: What does binary search require that linear search does not?
- The list must contain only numbers
- The list must be sorted *
- The list must have an even number of items
- The target must actually be present
E: Binary search decides which half to discard by comparing against the middle value, which is only meaningful if the data is in order.
:::

:::predict
Q: What does this print?
```c
#include <stddef.h>
#include <stdio.h>

size_t binary_search_passes(const int items[], size_t n, int target) {
    ptrdiff_t low = 0;
    ptrdiff_t high = (ptrdiff_t)n - 1;
    size_t passes = 0;
    while (low <= high) {
        ++passes;
        ptrdiff_t mid = low + (high - low) / 2;
        if (items[mid] == target) return passes;
        if (items[mid] < target) low = mid + 1;
        else high = mid - 1;
    }
    return passes;
}

int main(void) {
    int items[] = {2, 4, 6, 8, 10, 12, 14, 16};
    printf("%zu\n", binary_search_passes(items, 8, 2));
    return 0;
}
```
- 1
- 3 *
- 8
- 0
E: mid=3 sees 8 (too big, high=2), mid=1 sees 4 (too big, high=0), mid=0 sees 2 — match on the third pass. Three passes for eight items, because log₂(8) = 3.
:::

:::quiz
Q: You will search the same unchanging list of 500,000 items 10,000 times. What's the better plan?
- Linear search each time
- Sort once, then binary search each time *
- Sort before every single search
- Neither, searching is impossible at that size
E: Sorting costs O(n log n) once, and each of the 10,000 searches then costs O(log n) instead of O(n). The one-off sort pays for itself many times over.
:::

## Talk about it

> Binary search trades a precondition (sorted data) for a huge speed-up. Think of something outside programming that works the same way — a system that's fast only because someone did organising work up front. What happens in that system when the organisation is wrong or out of date?

## What's next

You've now seen a precondition buy you a dramatic speed-up, which raises an obvious question: how do you get a list sorted in the first place? Next up is **Bubble, Selection and Insertion Sort**, the three simple sorts. They're all O(n²), and you'll meet faster ones soon — but understanding exactly *why* they're slow is what makes merge sort and quick sort feel inevitable rather than magical.

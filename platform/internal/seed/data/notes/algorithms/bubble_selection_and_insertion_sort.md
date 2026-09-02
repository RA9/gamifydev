# Bubble, Selection and Insertion Sort

Sorting is the workhorse of computing. It powers binary search, it groups related records together, and it makes almost every other algorithm easier. This lesson covers the three simple sorts — the ones you could invent yourself with a pencil and a row of cards.

All three are O(n²), and you'll meet much faster sorts in the next lessons. But two of them have genuine real-world uses, and understanding exactly *where* their time goes is what makes the fast sorts feel like an obvious next step rather than a magic trick.

## Bubble sort: swap neighbours until nothing moves

**Bubble sort** repeatedly walks the list comparing each pair of neighbours, swapping them if they're out of order. After one full pass, the largest value has "bubbled" all the way to the end. Repeat, and the next largest settles into place, and so on.

```python
def bubble_sort(items):
    n = len(items)
    for i in range(n - 1):
        swapped = False
        for j in range(n - 1 - i):                  # the last i are already sorted
            if items[j] > items[j + 1]:
                items[j], items[j + 1] = items[j + 1], items[j]
                swapped = True
        if not swapped:                             # early exit: already sorted
            return items
    return items

print(bubble_sort([5, 1, 4, 2]))   # [1, 2, 4, 5]
```

A trace of the first pass over `[5, 1, 4, 2]`:

```text
[5, 1, 4, 2]   compare 5,1  -> swap
[1, 5, 4, 2]   compare 5,4  -> swap
[1, 4, 5, 2]   compare 5,2  -> swap
[1, 4, 2, 5]   pass 1 done — 5 is now in its final place
```

**Cost:** O(n²) worst case, O(n²) average. The `swapped` flag gives it one redeeming feature: on an already-sorted list, the first pass performs no swaps and the function returns immediately, so the **best case is O(n)**. Space is O(1) — it sorts in place. Bubble sort is **stable**: it only swaps on a strict `>`, so equal elements never cross each other.

:::warning
Without the `swapped` flag, bubble sort is O(n²) even on already-sorted input. The flag is what earns it the O(n) best case, and plenty of textbook versions leave it out.
:::

## Selection sort: find the minimum, put it in front

**Selection sort** takes a different tack. Scan the unsorted portion for the smallest element, then swap it into the front of that portion. Repeat with the remaining tail.

```python
def selection_sort(items):
    n = len(items)
    for i in range(n - 1):
        smallest = i
        for j in range(i + 1, n):
            if items[j] < items[smallest]:
                smallest = j
        if smallest != i:
            items[i], items[smallest] = items[smallest], items[i]
    return items

print(selection_sort([5, 1, 4, 2]))   # [1, 2, 4, 5]
```

Trace it:

```text
[5, 1, 4, 2]   smallest in [5,1,4,2] is 1  -> swap with position 0
[1, 5, 4, 2]   smallest in [5,4,2]   is 2  -> swap with position 1
[1, 2, 4, 5]   smallest in [4,5]     is 4  -> already there, no swap
[1, 2, 4, 5]   sorted
```

**Cost:** O(n²) in *all* cases — best, average and worst. There's no early exit possible, because it can't know the minimum of a range without looking at every element of it. Space is O(1).

What selection sort *is* good at is swap count: exactly one swap per position, so at most n-1 swaps total. Bubble sort can perform O(n²) swaps. If writing to the storage medium is genuinely expensive — old flash memory, or records so large that moving one costs real time — minimising writes can matter more than minimising comparisons.

The version above is **not stable**. Swapping a distant minimum into place can jump it over an equal element and reverse their original order.

:::example
Sort `[3a, 3b, 1]` where `3a` and `3b` are equal values you can tell apart. Selection sort swaps `1` with position 0, giving `[1, 3b, 3a]` — the two 3s have traded places. That's exactly what "not stable" means.
:::

## Insertion sort: place each card into a sorted hand

**Insertion sort** is how most people sort a hand of playing cards. Keep the left portion sorted. Take the next card, slide it leftwards past every card bigger than it, and drop it into place.

```python
def insertion_sort(items):
    for i in range(1, len(items)):
        current = items[i]
        j = i - 1
        while j >= 0 and items[j] > current:
            items[j + 1] = items[j]      # shift the bigger value right
            j -= 1
        items[j + 1] = current           # drop current into the gap
    return items

print(insertion_sort([5, 1, 4, 2]))   # [1, 2, 4, 5]
```

The trace, with `|` marking the boundary of the sorted portion:

```text
[5 | 1, 4, 2]   take 1: shift 5 right, insert 1 at front
[1, 5 | 4, 2]   take 4: shift 5 right, insert 4
[1, 4, 5 | 2]   take 2: shift 5, shift 4, insert 2
[1, 2, 4, 5 |]  sorted
```

**Cost:** O(n²) worst case (reverse-sorted input — every element travels the full distance) and O(n²) average. But the inner `while` stops the instant it finds something smaller, so on already-sorted input each element does exactly one comparison: **best case O(n)**, and more usefully, **O(n·k) on data where no element is more than k positions from home**. Space is O(1). It is **stable**, because the shift condition is a strict `>` — an equal element stops the shifting and the newcomer settles after it.

:::analogy
Insertion sort is sorting a hand of cards as they're dealt to you. You never re-examine the whole hand; you just slide each new card into its slot. If the deck was nearly in order already, almost every card goes straight to the end of your hand.
:::

## Stability, precisely

A sort is **stable** if elements that compare equal keep their original relative order. This matters more than it first sounds, because it lets you sort by several keys in sequence.

```python
players = [("ada", 90), ("bo", 85), ("cy", 90), ("di", 85)]
# already ordered by name. Now stable-sort by score, descending:
players.sort(key=lambda p: -p[1])
print(players)
# [('ada', 90), ('cy', 90), ('bo', 85), ('di', 85)]
```

Within each score, names stayed alphabetical — because Python's `sort` is stable. An unstable sort could have produced `('cy', 90), ('ada', 90)` and silently destroyed the earlier ordering.

```text
Sort           Best      Average   Worst     Space   Stable
------------------------------------------------------------
Bubble         O(n)      O(n^2)    O(n^2)    O(1)    yes
Selection      O(n^2)    O(n^2)    O(n^2)    O(1)    no
Insertion      O(n)      O(n^2)    O(n^2)    O(1)    yes
```

:::key
Stability means equal elements never cross. Bubble and insertion sort are stable because they only move an element past something **strictly** greater. Selection sort is not, because its long-distance swap can leapfrog an equal element.
:::

## Which of these do people actually use?

Bubble sort: essentially never in production. It's a teaching tool. Its one honest use is as the simplest possible thing to write when n is tiny and you'll never look at the code again.

Selection sort: rarely, and only when swap cost dominates comparison cost.

Insertion sort: **genuinely used, every day**. Two reasons. First, it's outstanding on small arrays — the constant factors are tiny, with no recursion and no extra allocation, so for arrays of maybe 10 to 30 elements it beats the asymptotically better sorts. Second, it's near-linear on nearly-sorted data, which is extremely common in practice.

Real library sorts exploit both facts. **Hybrid sorts** run a fast O(n log n) algorithm down to small sub-arrays and then finish with insertion sort. Python's built-in sort is a hybrid built around merge sort and runs of insertion sort, and it detects already-sorted runs so that sorting nearly-ordered data is close to O(n).

:::tip
"Asymptotically worse" doesn't mean "always slower." Big O describes growth as n gets large; for n = 15 the constant factors dominate, and that's precisely the gap insertion sort fills inside every serious sorting library.
:::

## Check Your Understanding

:::quiz
Q: Which of these three sorts is O(n²) even on an already-sorted list?
- Bubble sort with the swapped flag
- Selection sort *
- Insertion sort
- All three
E: Selection sort must scan the entire unsorted portion to find its minimum, every time. It has no way to detect that the work is already done, so its best case equals its worst case.
:::

:::predict
Q: What does this print?
```python
def insertion_sort(items):
    for i in range(1, len(items)):
        current = items[i]
        j = i - 1
        while j >= 0 and items[j] > current:
            items[j + 1] = items[j]
            j -= 1
        items[j + 1] = current
    return items

print(insertion_sort([3, 1, 2]))
```
- [1, 2, 3] *
- [3, 2, 1]
- [1, 3, 2]
- [2, 1, 3]
E: Take 1: shift 3 right, insert 1 -> [1, 3, 2]. Take 2: shift 3 right, insert 2 -> [1, 2, 3].
:::

:::quiz
Q: Why do production sorting libraries still contain insertion sort?
- It is the fastest sort for large arrays
- It is fast on small and nearly-sorted arrays, so hybrid sorts finish with it *
- It uses less memory than any other sort
- It is the only stable sort
E: Insertion sort has very low constant factors and is near-linear on nearly-ordered data, which makes it the right finisher for the small sub-arrays a divide-and-conquer sort produces.
:::

## Talk about it

> Selection sort makes the fewest swaps but always does the maximum number of comparisons; bubble sort can make far more swaps but can finish in one pass on sorted data. Describe a situation where you'd care much more about the number of swaps than the number of comparisons, and explain what makes moving data expensive there.

## What's next

Every sort in this lesson compares neighbours or scans linearly, and that's exactly what pins them at O(n²). The way out is to stop working element-by-element and start splitting the problem. Next up is **Merge Sort**, which cuts the list in half, sorts each half, and merges — giving O(n log n) in *every* case, not just on average. It's your first real divide-and-conquer algorithm.

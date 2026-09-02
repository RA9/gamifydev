# The Common Runtimes

Almost every algorithm you will meet in a normal programming life lands on one of eight rungs. Learning those eight — what each one feels like, and one real algorithm that lives there — turns complexity analysis from an abstract exercise into pattern recognition. This lesson walks the ladder from best to worst, then puts real numbers next to each rung so the gaps become impossible to ignore.

## The ladder, rung by rung

**O(1) — constant.** The work doesn't change when the input grows. Looking up `arr[500]` in an array takes the same time whether the array holds 600 items or 600 million, because the machine computes the address directly. Also here: pushing and popping a stack, and a hash table lookup *on average*.

**O(log n) — logarithmic.** Each step throws away a fixed fraction of what's left, usually half. Binary search on a sorted array: 1,000,000 items, 20 comparisons. Also here: search, insert and delete in a balanced binary search tree, and push/pop on a heap.

**O(n) — linear.** You touch each item a constant number of times. Finding the largest number in an unsorted list, walking a linked list, counting how many players scored above 50. Double the data, double the work.

**O(n log n) — linearithmic.** You do a linear amount of work at each of log n levels. This is the home of good general-purpose sorting: merge sort, heapsort, and quicksort's average case. It's only slightly worse than linear, and it's the practical ceiling for "fast."

:::key
The first four rungs — O(1), O(log n), O(n), O(n log n) — are all comfortably usable on huge inputs. Once you step above O(n log n), the input size you can handle starts collapsing fast.
:::

**O(n²) — quadratic.** Every item meets every item. Bubble sort, selection sort and insertion sort in the worst case; comparing all pairs in a list; a nested loop over the rows and columns of an n×n grid. Fine for hundreds of items, painful for millions.

**O(n³) — cubic.** Three nested loops. Textbook matrix multiplication of two n×n matrices does n multiplications for each of n² output cells. The Floyd–Warshall algorithm, which finds shortest paths between *all* pairs of nodes in a graph, is also O(n³) — and is genuinely used, because n is often small.

**O(2ⁿ) — exponential.** Every extra item *doubles* the work. Enumerating all subsets of n items is the classic: n items have 2ⁿ subsets, so you can't do better if you truly must look at all of them. Naive recursive Fibonacci sits roughly here too.

**O(n!) — factorial.** Every extra item multiplies the work by n. Enumerating all orderings of n items: n! permutations. Brute-force travelling salesman — try every route through n cities — lives here.

:::analogy
Climbing this ladder is like changing the unit on a map. O(log n) is a world map: doubling the territory barely changes the paper. O(n) is a street map. O(n²) is a floor plan. O(2ⁿ) is a diagram of every possible arrangement of your furniture — add one more chair and the diagram doubles.
:::

## Spotting a rung by the shape of the code

Each rung has a code shape that produces it, and after a while you'll recognise them without thinking. Here's the cheat sheet:

```text
rung          what the code is doing
-----------   ---------------------------------------------------
O(1)          a fixed amount of work, no loop over the input
O(log n)      a value or range that halves each step
O(n)          one pass over the data
O(n log n)    a linear pass repeated at each of log n levels
O(n^2)        a loop over the data inside a loop over the data
O(n^3)        three levels of that
O(2^n)        for each item, two branches: include it or don't
O(n!)         trying every possible ordering of the items
```

The bottom two rows are worth reading as sentences rather than symbols. "For each item, include it or don't" is why subsets are 2ⁿ. "Try every ordering" is why permutations are n!. When a brute-force idea sounds like either of those, you've found an algorithm you can only afford on tiny inputs.

## The numbers that make it visceral

Here is roughly how many operations each rung costs at four input sizes. Logarithms are base 2 and rounded; the very large numbers are approximations, marked with `~`.

```text
                n = 10        n = 100          n = 1,000        n = 1,000,000
-----------  ---------  --------------  -----------------  -------------------
O(1)                 1               1                  1                    1
O(log n)             3               7                 10                   20
O(n)                10             100              1,000            1,000,000
O(n log n)          33             664              9,966           ~2 x 10^7
O(n^2)             100          10,000          1,000,000           ~1 x 10^12
O(n^3)           1,000       1,000,000      1,000,000,000           ~1 x 10^18
O(2^n)           1,024   ~1.3 x 10^30     ~1.1 x 10^301        beyond writing
O(n!)        3,628,800  ~9.3 x 10^157    ~4 x 10^2567          beyond writing
```

Sit with the bottom-left corner for a moment. At n = 10 — ten items, a tiny list — the factorial algorithm already does 3.6 million operations while the linear one does ten. And at n = 100, `2ⁿ` is a 31-digit number. That is not "slow." That is "will not finish before the sun burns out, on any hardware, ever."

Meanwhile look across the `O(log n)` row: going from ten items to a million items costs you *seventeen extra operations*. That's the whole reason balanced trees and binary search exist.

:::warning
Exponential and factorial aren't merely bad — they're a different category. A faster computer moves an O(2ⁿ) algorithm from n = 40 to maybe n = 45. You can't buy your way out; you have to change the algorithm.
:::

## What size input each rung can actually handle

Flip the table around. As a rough thought experiment, suppose a machine gets through 100 million simple operations per second — a round number for reasoning, not a benchmark. Roughly what n can each rung finish in about a second?

```text
complexity     largest n finishing in ~1 second (rough)
------------   ---------------------------------------
O(log n)       essentially unlimited
O(n)           ~100,000,000
O(n log n)     ~5,000,000
O(n^2)         ~10,000
O(n^3)         ~460
O(2^n)         ~26
O(n!)          ~11
```

These are ballparks, and real constants shift them, but the *shape* is right and it is worth memorising. It explains a lot of engineering decisions at a glance. Why does a coding problem with "n ≤ 20" often want an exponential solution? Because 2²⁰ is about a million, which is nothing. Why does "n ≤ 100,000" rule out nested loops? Because 10¹⁰ operations is far too many.

:::tip
When you see the input limits on a problem, read them backwards to guess the intended complexity. `n ≤ 10⁶` says "find something linear or n log n." `n ≤ 500` says "n³ is probably fine." `n ≤ 20` says "brute force over subsets is expected."
:::

## A note on honesty at the top of the ladder

Two small precision points, because they come up.

Naive recursive Fibonacci is usually described as O(2ⁿ). That's a correct *upper bound*, but the tight answer is a little smaller — the call count grows like about 1.62ⁿ rather than 2ⁿ. This is a real example of Big O being a ceiling rather than an exact fit, exactly as the last lesson described.

And brute-force travelling salesman is described as O(n!), which is right, though the count of genuinely distinct round trips is smaller: with a fixed starting city and symmetric distances it's `(n-1)!/2`. Dividing by 2 and shifting n by one doesn't change the category at all — it's still factorial, still hopeless past a couple of dozen cities.

:::example
Twenty cities is a small map. Brute-force TSP over 20 cities examines `19!/2` routes — about 6 × 10¹⁶. At 100 million routes per second that's roughly 600 million seconds, or close to twenty years. Nineteen cities would take about one year. That single extra city cost you roughly eighteen more.
:::

## Check Your Understanding

:::quiz
Q: Binary search on a sorted array of 1,000,000 items does roughly how many comparisons?
- 1,000,000
- 1,000
- 20 *
- 1
E: Binary search is O(log n), and log₂(1,000,000) is about 20 — each comparison halves the range that's left.
:::

:::quiz
Q: You need to process 1,000,000 records. Which complexity is the highest rung that's still comfortably practical?
- O(n log n) *
- O(n²)
- O(2ⁿ)
- O(n!)
E: At n = 1,000,000, O(n log n) is about 2 × 10⁷ operations. O(n²) would be about 10¹², which is far too slow, and the rungs above that are impossible.
:::

:::match
Q: Match each complexity to an algorithm that has it.
- O(1) | Reading arr[i] from an array
- O(log n) | Binary search on a sorted array
- O(n log n) | Merge sort
- O(n²) | Bubble sort, worst case
E: These four are the workhorses you'll recognise most often; knowing them by sight makes analysing new code much faster.
:::

## Talk about it

> Look at the table again and pick the pair of rows that surprised you most. Explain in your own words what causes that gap — what is the O(n²) algorithm *doing* that the O(n log n) one avoids? Then describe a task from your own experience where you suspect an O(n²) approach is hiding.

## What's next

You can now recognise the eight common runtimes and reason about which one a problem can afford. But recognising a runtime in someone else's finished algorithm is different from deriving it from code you're reading right now. In **Analyzing Loops and Recursion** you'll learn the practical rules — loops multiply, blocks add, halving means log — and how to reason about recursive functions without any heavy mathematics.

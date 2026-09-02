# Analyzing Loops and Recursion

Knowing the common runtimes is one thing; looking at unfamiliar code and working out which one it is, is another. Happily, there are only a handful of rules, and they cover almost everything. This lesson gives you those rules for loops, then a picture-based way to handle recursion that needs no formal mathematics at all.

## Four rules for loops

**Rule 1: a loop over n items is O(n).** One pass, constant work inside, `n` iterations.

```python
def count_wins(matches):
    wins = 0
    for m in matches:        # n iterations
        if m == "W":
            wins += 1
    return wins

print(count_wins(["W", "L", "W"]))
# 2
```

**Rule 2: nested loops multiply.** If the outer runs n times and the inner runs n times *for each* outer iteration, that's n × n.

```python
def all_pairs(players):
    pairs = []
    for a in players:            # n
        for b in players:        # n, for each a
            pairs.append((a, b))
    return pairs

print(len(all_pairs([1, 2, 3])))
# 9
```

Three players give 3 × 3 = 9 pairs. In general n², so `O(n²)`. A third level of nesting would give `O(n³)`.

**Rule 3: sequential blocks add — so take the biggest.** Code that runs one after another sums its costs, and Big O then keeps only the dominant term.

```python
def summarise(nums):
    total = 0
    for x in nums:               # O(n)
        total += x

    for i in range(len(nums)):   # O(n^2)
        for j in range(len(nums)):
            if nums[i] == nums[j]:
                pass
    return total
```

That's `n + n²` steps, which is `O(n²)`. The linear pass is free by comparison — it makes no difference at all to the answer.

**Rule 4: halving means logarithmic.** If the range shrinks by a constant *fraction* each iteration (rather than a constant *amount*), the loop count is logarithmic.

```python
def binary_search(sorted_nums, target):
    lo, hi = 0, len(sorted_nums) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if sorted_nums[mid] == target:
            return mid
        if sorted_nums[mid] < target:
            lo = mid + 1          # throw away the lower half
        else:
            hi = mid - 1          # throw away the upper half
    return -1

print(binary_search([2, 4, 6, 8, 10], 8))
# 3
```

Each pass discards half the remaining range, so the number of passes is how many times you can halve n before reaching 1 — that's `log₂ n`, hence `O(log n)`.

:::key
Loop analysis in four lines: one loop over n → **multiply by n**. Nested loops → **multiply**. Blocks in sequence → **add, then keep the biggest**. A range that halves → **log n**.
:::

## The triangle loop

Here's the pattern that trips people up, because the inner loop's length changes:

```python
def triangle(n):
    count = 0
    for i in range(n):         # i = 0, 1, 2, ..., n-1
        for j in range(i):     # runs i times
            count += 1
    return count

print(triangle(5))
# 10
```

The inner loop runs 0 times, then 1, then 2, then 3, then 4. The total is `0 + 1 + 2 + 3 + 4 = 10`. Not n², so surely it's better than quadratic?

It isn't. Adding up `0 + 1 + 2 + ... + (n-1)` gives `n(n-1)/2`. (The close cousin `1 + 2 + ... + n` equals `n(n+1)/2` — the same shape, shifted by one.) Multiply that out and you get `n²/2 - n/2`. Drop the constant and the lower-order term, and you're left with `O(n²)`.

```text
n        inner-loop total = n(n-1)/2      n^2
-----    ---------------------------      ----------
5                                 10              25
100                            4,950          10,000
1,000                        499,500       1,000,000
```

The triangle really is about half the full square — but "half of quadratic" is still quadratic, and Big O deliberately can't tell the two apart, because the factor of 2 doesn't grow.

:::warning
A common misreading is "the inner loop doesn't run n times, so it isn't O(n²)." What matters is the *total* across all iterations, not the length of any single inner pass. Half of n² is still n².
:::

## Watch for hidden loops

Not every loop looks like a `for` — some of Python's most convenient operations quietly iterate:

```python
if name in name_list:      # a list: O(n) - it scans
if name in name_set:       # a set: O(1) on average - it hashes
```

Put the first one inside a loop over n names and you've written an `O(n²)` algorithm that has only one visible loop — exactly the duplicate-usernames trap from the first lesson.

:::tip
Before analysing code, ask of every function call and operator: "what does this cost?" A list `in`, a `list.remove`, a slice like `nums[1:]`, joining strings with `+` in a loop — each is linear, and each will multiply whatever loop it sits inside.
:::

## Recursion as a tree

For recursion, draw the calls as a tree and ask two questions: **how many levels deep does it go**, and **how much work happens across each level**? Multiply them.

**Shape 1: `T(n) = T(n/2) + O(1)` → O(log n).** One recursive call, on half the input, with constant work at each step. This is binary search written recursively. The input goes n → n/2 → n/4 → ... → 1, which is `log₂ n` levels, each doing a fixed amount of work.

**Shape 2: `T(n) = 2T(n/2) + O(n)` → O(n log n).** Two calls on half the input each, plus a linear amount of work to combine the results. This is merge sort. Picture it at n = 8:

```text
level 0:  [ 8 items ]                       8 units of merging
level 1:  [4] [4]                           4 + 4  = 8
level 2:  [2] [2] [2] [2]                   2+2+2+2 = 8
level 3:  [1][1][1][1][1][1][1][1]          base case
```

Every level costs about n, because the pieces get smaller but there are proportionally more of them. And the number of levels is how many times you can halve 8: three. So the total is `n × log n` → `O(n log n)`.

**Shape 3: `T(n) = T(n-1) + O(1)` → O(n).** One call, on an input one *smaller* — not halved. That's n levels of constant work.

```python
def countdown(n):
    if n == 0:
        return 0
    return 1 + countdown(n - 1)

print(countdown(5))
# 5
```

Compare shapes 1 and 3 carefully: `n/2` gives log n, `n - 1` gives n. Subtracting is dramatically weaker than dividing.

**Shape 4: `T(n) = 2T(n-1) + O(1)` → O(2ⁿ).** Two calls, each on an input just one smaller. Now the tree *branches* n levels deep, so the number of calls doubles at every level: 1, 2, 4, 8, 16... Towers of Hanoi is the honest example — moving n discs takes exactly `2ⁿ - 1` moves, and there is no cleverer way. Twenty discs is 1,048,575 moves; forty discs is 1,099,511,627,775 — over a trillion, for twenty more discs.

:::analogy
The difference between shapes 3 and 4 is the difference between a queue and a rumour. In shape 3 each call passes the job to exactly one successor — a line of people, n long. In shape 4 each call tells *two* people, who each tell two more. After n rounds the line has n people and the rumour has reached 2ⁿ.
:::

Not every branching recursion is doomed, though — sometimes the branches overlap, and that changes everything.

:::example
Naive recursive Fibonacci is `T(n) = T(n-1) + T(n-2) + O(1)` — a branching shape, so it's exponential. Adding memoisation changes the tree completely: each value of n is computed once and then reused, so there are only n distinct calls, giving `O(n)`. Same formula, same code shape, and a change from exponential to linear.
:::

## Check Your Understanding

:::predict
Q: What is the time complexity of this function?
```python
def f(nums):
    n = len(nums)
    for i in range(n):
        for j in range(i):
            print(nums[i], nums[j])
    for k in range(n):
        print(nums[k])
```
- O(n)
- O(n log n)
- O(n²) *
- O(n³)
E: The triangle loop totals n(n-1)/2, which is O(n²). The separate loop adds O(n), and the dominant term is n².
:::

:::quiz
Q: A recursive function satisfies T(n) = T(n/2) + O(1). What is its complexity?
- O(1)
- O(log n) *
- O(n)
- O(n log n)
E: One call on half the input means log₂ n levels, with constant work at each — that's O(log n), the shape of binary search.
:::

:::fill
Q: Complete the loop condition so this runs in O(log n) rather than O(n).
`while n > 1: n = n ___ 2`
- // *
- -
- +
E: Integer-dividing by 2 halves the value each pass, giving log₂ n iterations. Subtracting 2 would give n/2 iterations, which is still O(n).
:::

## Talk about it

> Take a function you've written recently with a nested loop in it, or picture one. Work through its complexity out loud using the four rules, naming what n is for that function. Then describe what would have to change about the data structures you used for the inner operation to drop it by a whole rung.

## What's next

You've been counting time all course. But every algorithm also spends *memory*, and the same counting technique applies — including a source of memory use you may not have noticed: the call stack that recursion builds. In **Space Complexity** you'll learn to measure memory the same way, understand what "in-place" really means, and see the trade where you buy speed with memory.

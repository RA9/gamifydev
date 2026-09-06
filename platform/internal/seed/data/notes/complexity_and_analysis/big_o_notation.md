# Big O Notation

Last lesson you counted steps and got expressions like `n + 2` and `n²/2`. Big O is the agreed-upon way to write those down so that every programmer in the world means the same thing. It deliberately throws away detail — and the detail it throws away is exactly the detail that never mattered. By the end of this lesson you'll be able to look at a short piece of code and say its Big O out loud with confidence.

## An upper bound on growth

Big O answers one question: **as n gets large, what does the work grow no faster than?**

When we say an algorithm is `O(n²)`, we're saying: past some point, its step count stays below some fixed multiple of n². Not equal to n² — *below a multiple of* n². Big O is a **ceiling**, an upper bound on the growth.

That "some fixed multiple" is doing real work in the definition. We're allowed to pick any constant we like, once, and it must hold forever after. So `n²/2` is `O(n²)` (pick the multiple 1), and `500n²` is also `O(n²)` (pick 500). Both grow *like* n²; they just start from different places.

:::analogy
Big O is a speed limit sign, not a speedometer. "O(n²)" says "this algorithm's growth never exceeds the n² line." It doesn't promise the algorithm actually runs that slowly, and it certainly doesn't tell you the number on the speedometer at any given moment.
:::

## Why we drop constants and small terms

Suppose you carefully count and find an algorithm does exactly `3n² + 5n + 100` steps. Big O calls this `O(n²)`. Two things got dropped: the `3`, and the entire `5n + 100`. Here's why that's honest rather than lazy.

**Lower-order terms fade.** Watch what share of the total each piece contributes:

```text
n         3n^2            5n          100      total       3n^2 share
-----     -----------     -------     ----     ----------  ----------
10                300          50      100            450      67%
100            30,000         500      100         30,600      98%
1,000       3,000,000       5,000      100      3,005,100    99.8%
```

At n = 10 the extra terms matter. At n = 1,000 they're rounding error. Since Big O is a statement about *large* n, the n² term is the only one still standing.

**Constants are machine details.** The `3` says each "step" of the loop costs three primitive operations on this particular implementation. Rewrite it in C, or on a faster processor, and the 3 becomes a 1 or a 7. It is exactly the kind of thing that varies with the stopwatch, so Big O refuses to record it.

:::key
To go from a step count to Big O: **keep the fastest-growing term, drop everything else, drop the constant multiplier.** `3n² + 5n + 100` → `O(n²)`. `7n + 4` → `O(n)`. `12` → `O(1)`.
:::

One caution before we move on, because this simplification can be taken too far.

:::warning
Dropping constants is a statement about *growth*, not about real life. A `O(n)` algorithm with a huge constant can absolutely be slower than an `O(n²)` one on small inputs. Big O tells you who wins eventually — it doesn't tell you where "eventually" begins.
:::

## Saying it out loud

The notation has a spoken form you'll hear in every code review and interview:

```text
written        said
--------       ----------------------------------------
O(1)           "oh of one"          / "constant time"
O(log n)       "oh of log n"        / "logarithmic"
O(n)           "oh of n"            / "linear"
O(n log n)     "oh of n log n"      / "linearithmic"
O(n^2)         "oh of n squared"    / "quadratic"
O(2^n)         "oh of two to the n" / "exponential"
```

It's a capital letter O, not a zero. It's often written "big-O" to make that clear. And when you see `log n` in a complexity, it's base 2 unless someone says otherwise — though you'll see next lesson that the base barely matters here.

## Deriving Big O from code

The technique is: count what dominates, then simplify. Let's do several.

**A single pass.**

```c
#include <stddef.h>
#include <stdio.h>

int largest(const int nums[], size_t n) {
    int best = nums[0];
    for (size_t i = 0; i < n; i++) {  // runs n times
        if (nums[i] > best) {
            best = nums[i];
        }
    }
    return best;
}

int main(void) {
    const int nums[] = {4, 9, 2};
    printf("%d\n", largest(nums, 3));  // 9
    return 0;
}
```

One loop over n items, constant work inside. Roughly `n` steps plus a fixed handful → `O(n)`.

**Two loops in a row.**

```c
#include <stddef.h>

typedef struct {
    long total;
    int best;
} Summary;

Summary sum_then_max(const int nums[], size_t n) {
    long total = 0;
    for (size_t i = 0; i < n; i++) {  // n steps
        total += nums[i];
    }

    int best = nums[0];
    for (size_t i = 0; i < n; i++) {  // another n steps
        if (nums[i] > best) {
            best = nums[i];
        }
    }
    return (Summary){total, best};
}
```

That's `n + n = 2n` steps. Drop the constant → `O(n)`. Sequential work *adds*, and adding two linear passes is still linear.

**A loop inside a loop.**

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

bool any_pair_sums_to(const int nums[], size_t n, int target) {
    for (size_t i = 0; i < n; i++) {      // n times
        for (size_t j = 0; j < n; j++) {  // n times, for each i
            if (nums[i] + nums[j] == target) {
                return true;
            }
        }
    }
    return false;
}

int main(void) {
    const int nums[] = {1, 5, 9};
    printf("%s\n", any_pair_sums_to(nums, 3, 10) ? "true" : "false");  // true
    return 0;
}
```

The inner loop runs n times for each of the n outer iterations: `n × n = n²` → `O(n²)`. Nested work *multiplies*.

:::example
Note the early `return true`. Big O usually describes the **worst case** unless someone says otherwise, and the worst case here is "no pair matches," where both loops run to completion. If you want to talk about the lucky case, say so explicitly: "O(1) best case, O(n²) worst case."
:::

**A fixed number of iterations.**

```c
#include <stddef.h>
#include <stdio.h>

void first_ten(const int nums[static 10]) {
    for (size_t i = 0; i < 10; i++) {  // always 10, regardless of n
        printf("%d\n", nums[i]);
    }
}
```

Ten iterations whether the array holds 12 items or 12 million. The work doesn't grow with n at all → `O(1)`. Constant time doesn't mean *fast*; it means *unchanging as n grows*.

**Halving the range.**

```c
#include <stdio.h>

unsigned count_halvings(unsigned n) {
    unsigned steps = 0;
    while (n > 1) {
        n /= 2;  // the range shrinks by half each time
        steps++;
    }
    return steps;
}

int main(void) {
    printf("%u\n", count_halvings(1000));  // 9
    return 0;
}
```

Starting at 1,000, `n` goes 500, 250, 125, 62, 31, 15, 7, 3, 1 — nine steps. Double the input to 2,000 and you get just one more step. That's the signature of `O(log n)`.

:::tip
When you see a value being *halved* (or divided by any fixed amount) each iteration, reach for `O(log n)`. When you see it being *decremented by one*, reach for `O(n)`. That single distinction explains why binary search beats a linear scan.
:::

## It describes growth, not seconds

This is the part people quietly get wrong, so let's be blunt about it. `O(n)` is not a duration. It cannot be converted into milliseconds. It contains no information about your CPU, your language or your input.

What it *does* let you say is things like: "this is `O(n)`, so if I ten-times the input, expect roughly ten times the work." That's a prediction about *ratios*, and ratios survive the trip from your laptop to production.

```text
claim                                   legitimate?
-------------------------------------   -----------
"binary search is O(log n)"             yes
"binary search takes O(log n) seconds"  no - O has no units
"O(n) is always faster than O(n^2)"     no - only for large enough n
"O(n) grows more slowly than O(n^2)"    yes
```

:::key
Big O is a claim about the *shape of the growth curve*, not a point on it. It answers "what happens when the input gets bigger?" and refuses to answer "how long will this take?"
:::

## Check Your Understanding

:::quiz
Q: What is the Big O of an algorithm that performs exactly `6n + 300` steps?
- O(6n + 300)
- O(n) *
- O(300)
- O(n²)
E: Drop the constant multiplier 6 and the lower-order term 300. The fastest-growing term is n, so it's O(n).
:::

:::predict
Q: What Big O describes this function's running time?
```c
#include <stddef.h>

long f(size_t n, const int grid[n][n]) {
    long total = 0;
    for (size_t i = 0; i < n; i++) {
        for (size_t j = 0; j < n; j++) {
            total += grid[i][j];
        }
    }
    for (size_t i = 0; i < n; i++) {
        total += (long)i;
    }
    return total;
}
```
- O(n)
- O(n²) *
- O(n³)
- O(2ⁿ)
E: The nested loops give n × n = n² and the separate loop gives n. Sequential blocks add, so n² + n, and the dominant term is n² → O(n²).
:::

:::match
Q: Match each step count to its Big O.
- 4n + 9 | O(n)
- 2n² + 1000n | O(n²)
- 17 | O(1)
E: Keep only the fastest-growing term and drop the constant multiplier in front of it.
:::

## Talk about it

> Big O deliberately discards the constant factor, so `O(n)` code with an expensive step and `O(n)` code with a cheap step get the same label. Describe a situation where that discarded constant would matter enormously to you as an engineer, and explain how you'd decide between two algorithms that share a Big O.

## What's next

Big O gives you an upper bound — a ceiling. But a ceiling alone can't say "this algorithm's growth is *exactly* this shape," and sometimes that's what you actually mean. In the next lesson, **Big-Theta, Big Omega and the Small Variants**, you'll meet the rest of the family: the floor, the tight fit, and the strict versions — plus why almost everyone says "Big O" when they really mean something else.

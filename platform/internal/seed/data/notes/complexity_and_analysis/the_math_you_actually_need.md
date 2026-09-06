# The Math You Actually Need

Complexity analysis has a reputation for demanding heavy mathematics. It doesn't. The whole subject leans on a small toolkit — logarithms, powers of two, one summation formula, a little counting and a little probability — and every piece of it can be understood concretely. This lesson assembles that toolkit, and ties each tool to a result you've already used.

## Logarithms: how many times can I halve it?

Forget any definition involving inverse functions. For our purposes:

> `log₂ n` is **how many times you can halve n before you reach 1**.

Start at 1,000. Halve it: 500, 250, 125, 62, 31, 15, 7, 3, 1. That's nine halvings, and `log₂ 1000` is about 9.97. Round trip confirmed.

The other direction is just as useful: `log₂ n` is **the exponent you'd raise 2 to in order to get n**. Since 2¹⁰ = 1024, `log₂ 1024 = 10` exactly.

```text
n              log2 n (about)
-----------    --------------
8                           3
1,000                      10
1,000,000                  20
1,000,000,000              30
```

Stare at that table for a moment, because it's the entire reason binary search and balanced trees matter. Going from a thousand items to a *billion* — a millionfold increase — adds twenty steps.

:::key
Logarithms grow absurdly slowly. That's the whole point. Any time an algorithm discards a fixed fraction of the remaining work each step, its cost is logarithmic, and logarithmic costs barely notice how big your data gets.
:::

**Why the base rarely matters.** You'll sometimes see `log n` with no base. In Big O that's fine, because changing the base only changes a constant factor:

```text
log2 n  = log10 n / log10 2  = log10 n x 3.32...
```

Multiplying by a fixed 3.32 is exactly the sort of constant Big O throws away, so `O(log₂ n)`, `O(log₁₀ n)` and `O(ln n)` are all the same class. Computer science defaults to base 2 because we halve things so often, but the notation doesn't care.

## Powers of two

Because computers work in binary and algorithms halve things, a handful of powers of two are worth knowing on sight:

```text
2^10  = 1,024                     ~ a thousand
2^20  = 1,048,576                 ~ a million
2^30  = 1,073,741,824             ~ a billion
2^32  = 4,294,967,296             the classic 32-bit limit
2^64  = 18,446,744,073,709,551,616  ~ 1.8 x 10^19
```

The pattern "every 10 powers is roughly another factor of a thousand" lets you estimate logs in your head. A million items? About 2²⁰, so about 20 halvings. A billion? About 30.

There's one more fact hiding in this list. Add up the powers of two on the way up:

```text
1 + 2 + 4 + 8 + ... + 1024 = 2047,  which is just under 2 x 1024
```

The whole sum is barely more than its last term. That's why a dynamic array that doubles its capacity when it fills up ends up doing only about `2n` total copying work across n appends — so each append costs `O(1)` on average, even though individual appends occasionally do a big copy.

:::analogy
Doubling is so lopsided that the final step outweighs everything before it combined. Folding a piece of paper, the last fold moves more thickness than all the earlier folds put together. That's why "amortised O(1)" works: the rare expensive operation is paid for by everything cheap that came before it.
:::

## The arithmetic series

One summation formula covers most of what you'll need:

```text
1 + 2 + 3 + ... + n  =  n(n+1)/2
```

Check it at n = 4: `1 + 2 + 3 + 4 = 10`, and `4 × 5 / 2 = 10`. And at n = 100: `100 × 101 / 2 = 5050`.

Why it's true, without a proof: pair the first with the last, the second with the second-to-last. `1 + 100 = 101`. `2 + 99 = 101`. Every pair sums to `n + 1`, and there are `n/2` pairs. Multiply.

```c
#include <stdbool.h>
#include <stdio.h>

bool check(unsigned n) {
    unsigned long long sum = 0;
    for (unsigned i = 1; i <= n; i++) {
        sum += i;
    }
    return sum == (unsigned long long)n * (n + 1) / 2;
}

int main(void) {
    printf("%s\n", check(100) ? "true" : "false");  // true
    return 0;
}
```

This is the formula behind the triangle loop. `for (size_t i = 0; i < n; i++)` with `for (size_t j = 0; j < i; j++)` inside runs `0 + 1 + ... + (n-1) = n(n-1)/2` times. Multiply that out and it's `n²/2 - n/2`, which Big O flattens to `O(n²)`. Whenever you see a nested loop whose inner length depends on the outer counter, this formula is what tells you it's still quadratic.

## Counting: permutations and combinations

Two counting ideas explain the top of the runtime ladder.

**Permutations — how many orderings?** With n distinct items, you have n choices for first place, `n-1` for second, and so on. That product is **n factorial**:

```text
5!  = 120
10! = 3,628,800
13! = 6,227,020,800
20! = 2,432,902,008,176,640,000    (~2.4 x 10^18)
```

This is exactly why brute-force travelling salesman is `O(n!)`: a tour *is* an ordering of the cities, and checking every tour means checking every ordering. Twenty cities means about 2.4 × 10¹⁸ orderings before you even start trimming. It also explains the deep difference between `2ⁿ` and `n!` — each new item *doubles* an exponential cost, but *multiplies a factorial cost by n*. Factorial pulls away and never looks back.

**Combinations — how many selections, order not mattering?** Written "n choose k":

```text
C(n, k) = n! / (k! x (n-k)!)
```

`C(5, 2) = 120 / (2 × 6) = 10`: ten ways to pick 2 items from 5. A more famous one: `C(52, 5) = 2,598,960`, the number of distinct five-card poker hands.

And the related fact that explains exponential algorithms: **a set of n items has 2ⁿ subsets**, because each item is independently either in or out. Any algorithm that must examine every subset is `O(2ⁿ)` and cannot be otherwise.

:::example
Ten items have 2¹⁰ = 1,024 subsets and 10! = 3,628,800 orderings. Twenty items have about a million subsets but 2.4 × 10¹⁸ orderings. Same jump in n; wildly different consequences. When a brute-force solution is "try every subset," small n is survivable. When it's "try every ordering," it isn't.
:::

## Probability and expected value

The last tool explains why two of your favourite tools are described with the word "average."

**Expected value** is the long-run average outcome: each result multiplied by its probability, all added up. A fair six-sided die gives `(1+2+3+4+5+6)/6 = 3.5`. You'll never roll a 3.5 — expected value describes the average over many trials, not any single one.

**Why hash tables are O(1) on average.** A hash table with n items spread over m buckets has an average of `n/m` items per bucket. Keep that ratio bounded — resize when it gets too big — and an average lookup examines a constant number of items, giving `O(1)`.

The word *average* is load-bearing. If every key happens to hash to the same bucket, that bucket becomes a linked list of n items and a lookup is `O(n)`. That's the true worst case, and it's why we say hash tables are `O(1)` **average**, `O(n)` **worst**. In practice a decent hash function scatters keys well enough that the worst case doesn't occur by accident.

:::warning
"O(1) average" is a statement about typical behaviour, not a guarantee. Adversarial input designed to collide on purpose can push a hash table to its `O(n)` worst case — a real class of denial-of-service attack, and the reason serious hash implementations randomise their hashing.
:::

**Why randomised quicksort works.** Quicksort's `O(n²)` worst case happens when the pivot is consistently terrible — already-sorted input with a first-element pivot is the classic trap. Choosing the pivot *at random* doesn't make that case impossible; it makes it astronomically unlikely, because an adversary can no longer predict your choices. The expected running time becomes `O(n log n)` regardless of what input arrives.

:::tip
Notice the pattern in both cases: randomness converts "bad on certain inputs" into "bad with tiny probability on all inputs." That's a genuinely powerful engineering idea, and you'll meet it again anywhere worst cases are caused by structure in the data.
:::

## Check Your Understanding

:::quiz
Q: Roughly what is log₂ of 1,000,000?
- 6
- 20 *
- 1,000
- 500,000
E: 2²⁰ is 1,048,576, just over a million — so you can halve a million about 20 times before reaching 1.
:::

:::quiz
Q: Why do we say hash table lookup is O(1) *on average* rather than simply O(1)?
- Because hashing itself is slow
- Because if many keys collide into one bucket, a lookup degrades to O(n) *
- Because it depends on the programming language
- Because average and worst case always mean the same thing
E: With a good hash and a bounded load factor, lookups examine a constant number of items. But collisions can pile keys into one bucket, and in the worst case a lookup scans all n.
:::

:::fill
Q: Complete the closed form for the sum 1 + 2 + ... + n.
`total = n * (n + 1) ___ 2`
- / *
- +
- -
E: The arithmetic series 1 + 2 + ... + n equals n(n+1)/2 — the formula behind the triangle loop's O(n²).
:::

## Talk about it

> Of the tools in this lesson — logarithms, powers of two, the arithmetic series, factorials and combinations, expected value — pick the one that changed how you read code the most. Explain a specific piece of code, yours or anyone's, that you now understand differently, and say what you would have guessed about it before this course.

## What's next

That's the course. You can now count operations instead of seconds, express growth in Big O and know exactly how strong that claim is, recognise the eight common runtimes on sight, analyse loops and recursion from the code in front of you, budget memory as carefully as time, and speak honestly about P, NP and what nobody yet knows.

That's a genuine shift in how you see programs — you have the vocabulary that turns "this feels slow" into a specific, arguable claim. Next comes the **Algorithms** course, where you'll put all of it to work: sorting, searching, graph traversal, greedy methods, dynamic programming and more, analysing each one as you build it. You'll find that the hard part of learning a new algorithm is no longer the algorithm. Pixel will see you there.

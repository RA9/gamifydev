# Big-Theta, Big Omega and the Small Variants

Big O is one member of a family of five. Each one makes a different claim about how one function's growth compares to another's, and the differences are exactly the differences between "at most," "at least" and "exactly." You'll mostly use Big O in daily conversation — but knowing the whole family is what lets you read a paper, follow a proof sketch, or notice when someone's claim is weaker than it sounds.

## Five comparisons, one idea

Every member of the family compares your algorithm's step count to a simple reference function like `n`, `n²` or `n log n`. The only thing that changes is the *kind* of comparison:

```text
notation      spoken           comparison         plain meaning
-----------   --------------   ----------------   -----------------------------
O(g)          "big oh"         <=                 grows no faster than g
Omega(g)      "big omega"      >=                 grows no slower than g
Theta(g)      "big theta"      =                  grows exactly like g
o(g)          "little oh"      <  (strictly)      grows strictly slower than g
omega(g)      "little omega"   >  (strictly)      grows strictly faster than g
```

Every one of these ignores constant factors and only cares about large n — that part never changes. Only the direction and the strictness change.

:::analogy
Think of a stretch of motorway. **Big O** is the speed limit sign: you never exceed it. **Big Omega** is a minimum-speed sign: you never drop below it. **Big Theta** is a description of how you actually drive: between the two, hemmed in on both sides. The little variants are the strict versions — little-o means you're *always eventually* well under the limit, not merely at it.
:::

## Big O: at most (the ceiling)

`f(n) = O(g(n))` means: past some starting point, `f(n)` stays at or below some fixed multiple of `g(n)`.

Because it's only a ceiling, a true Big O claim can be wildly generous. If an algorithm does `3n` steps, all of these are *true*:

```text
3n = O(n)        true, and tight
3n = O(n log n)  true, but loose
3n = O(n^2)      true, but very loose
3n = O(2^n)      true, and almost useless
```

Nothing above is a lie. A ceiling at 30,000 metres really is a ceiling. It's just not informative.

:::warning
"This algorithm is O(n³)" does **not** mean it takes n³ time. It means it takes *no more than* about n³ time. If someone reassures you that their code is O(2ⁿ), they've told you almost nothing — every practical algorithm is O(2ⁿ).
:::

## Big Omega: at least (the floor)

`f(n) = Ω(g(n))` means: past some starting point, `f(n)` stays at or above some fixed multiple of `g(n)`. It's the mirror image of Big O.

You'll meet Big Omega most often in statements about *problems* rather than *algorithms* — claims that no possible solution can do better. The famous one:

> Any sorting algorithm that works by comparing pairs of elements needs Ω(n log n) comparisons in the worst case.

That's a floor under every comparison sort that has ever been or will ever be written. It's why merge sort's `O(n log n)` isn't just good — it's optimal for its category. (Sorts that don't compare elements, like counting sort, sidestep the floor by using extra assumptions about the data.)

:::key
Big O bounds an algorithm from **above**: "it costs no more than this." Big Omega bounds it from **below**: "it costs at least this." Omega is how you say "nobody can do better than this," which is a claim about the problem, not about one person's code.
:::

## Big Theta: exactly (both at once)

`f(n) = Θ(g(n))` means both things at once: `f` is `O(g)` *and* `f` is `Ω(g)`. It's squeezed between a fixed multiple of `g` above and a fixed multiple of `g` below, so it genuinely grows *like* g.

```text
3n^2 + 5n = O(n^2)      yes - ceiling
3n^2 + 5n = Omega(n^2)  yes - floor
3n^2 + 5n = Theta(n^2)  yes - both, so it is a tight description

3n^2 + 5n = O(n^3)      yes - true but loose
3n^2 + 5n = Theta(n^3)  NO  - it does not grow like n^3
```

Theta is the strongest and most useful of the three, because it's the only one that can be *wrong* in the direction that matters. Anyone can inflate a Big O claim; nobody can inflate a Theta claim.

## The little variants: strictly

Big O allows `f` and `g` to grow at the same rate (`3n` is `O(n)`). **Little-o** forbids that. `f(n) = o(g(n))` means `f` grows *strictly* slower than `g` — as n grows, the ratio `f(n)/g(n)` shrinks toward zero.

```text
n       = o(n^2)      yes - n/n^2 = 1/n, which shrinks toward 0
n log n = o(n^2)      yes - beaten by n^2 by an ever-widening margin
3n      = o(n)        NO  - the ratio stays at 3, it never shrinks
n^2     = o(n^2)      NO  - a function never grows strictly slower than itself
```

**Little-omega** (`ω`) is the strict version of Big Omega: `f(n) = ω(g(n))` means `f` grows strictly *faster*, so the ratio `f(n)/g(n)` grows without limit. `n²` is `ω(n)`. `2ⁿ` is `ω(n¹⁰⁰)` — exponential growth eventually leaves every polynomial behind, however large its exponent.

The little variants are rare in everyday engineering talk. They show up when someone needs to say "definitively better, not merely no worse" — for instance, "we found an algorithm that is `o(n²)`," which rules out a mere constant-factor improvement.

:::tip
A quick way to keep the five straight: Big versions allow ties (≤, ≥). Little versions don't (<, >). Theta is the one that pins you down on both sides.
:::

## Why everyone says "Big O" when they mean Theta

Listen to programmers talk and you'll hear "merge sort is O(n log n)" constantly. That's true. But what the speaker almost always means is "merge sort is **Θ**(n log n)" — it grows *exactly* like n log n, not merely no faster.

This sloppiness is usually harmless. In context, everyone understands that a stated Big O is meant to be the tightest one the speaker knows. Nobody says "my algorithm is O(n!)" about a linear scan, even though it's technically true.

It stops being harmless in two situations.

**When the bound is quietly loose.** "Our new lookup is O(n)" sounds like a linear scan. If it's actually Θ(log n), the speaker has undersold it — and if a reviewer rejects the design on that basis, the sloppiness cost something real.

**When best, worst and average cases differ.** Here's the important one, and it's the single most common confusion in this whole topic:

:::warning
Big O is **not** "the worst case" and Big Omega is **not** "the best case." Those are two completely separate axes. You first pick *which case you're describing* (best, average, worst), and then you pick *which bound* you're stating about it (O, Ω, Θ). You can absolutely say "the best case is Θ(n)" or "the worst case is Ω(n²)."
:::

Quicksort makes the point:

```text
quicksort, average case  = Theta(n log n)
quicksort, worst case    = Theta(n^2)
```

Both statements are true; they describe different inputs. Saying only "quicksort is O(n²)" is honest but leaves out why anyone uses it. Saying only "quicksort is O(n log n)" hides the pathological case. The precise habit is to name the case *and* the bound: "quicksort is Θ(n log n) on average and Θ(n²) in the worst case."

:::example
Insertion sort on an already-sorted list does one comparison per element and no shifting: that's Θ(n), its best case. On a reverse-sorted list it shifts everything, every time: Θ(n²), its worst case. The single label "O(n²)" is a correct ceiling for both, but it hides that insertion sort is genuinely excellent on nearly-sorted data — which is why real sorting libraries use it for small or almost-ordered chunks.
:::

## Check Your Understanding

:::quiz
Q: An algorithm always performs exactly 4n + 7 steps. Which statement is FALSE?
- It is O(n)
- It is Ω(n)
- It is Θ(n)
- It is Θ(n²) *
E: 4n + 7 is O(n²) (a valid ceiling) but not Ω(n²), so it cannot be Θ(n²). Theta requires the growth to match on both sides.
:::

:::quiz
Q: Which statement correctly describes the relationship between Big O and the worst case?
- Big O always means the worst case
- Big Omega always means the best case
- Case (best/average/worst) and bound (O/Ω/Θ) are independent choices *
- Big O can only be used on sorted input
E: You choose which input case you're describing, then choose which kind of bound to state about it. "Best case is Θ(n)" and "worst case is Ω(n²)" are both perfectly valid.
:::

:::match
Q: Match each notation to what it claims about growth.
- O(g) | Grows no faster than g
- Ω(g) | Grows no slower than g
- Θ(g) | Grows at the same rate as g
- o(g) | Grows strictly slower than g
E: Big O is a ceiling, Big Omega a floor, Big Theta both, and little-o the strict version of the ceiling.
:::

## Talk about it

> Someone tells you their search function is "O(n²)". Later you read the code and discover it's actually Θ(log n) — far better than advertised. Their original claim was technically true. Describe a concrete situation on a team where that technically-true statement could still cause a bad decision, and explain what you'd want them to have said instead.

## What's next

You can now state a growth rate precisely and say exactly how strong your claim is. Next it's time to get familiar with the specific growth rates you'll actually meet — the handful of shapes that cover nearly every algorithm you'll ever write. In **The Common Runtimes** you'll walk the ladder from constant to factorial, with real algorithms at each rung and a table of actual operation counts that makes the differences impossible to forget.

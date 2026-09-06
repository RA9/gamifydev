# Greedy Algorithms

You've now met four algorithms that all share a strategy: take the best-looking option right now, commit to it, and never reconsider. Dijkstra always expands the nearest node. Prim always adds the cheapest crossing edge. Kruskal always takes the cheapest safe edge. That strategy has a name — **greedy** — and this lesson examines it directly.

The interesting question isn't how greedy works; it's *when it's allowed to work*. Greedy algorithms are fast and beautifully simple, and they are also wrong on a great many problems. Knowing which side of the line you're on is the whole skill.

## The greedy choice

A greedy algorithm builds its answer one decision at a time. At each step it applies some local rule — "take the cheapest", "take the one that finishes soonest", "take the largest that fits" — commits to that choice, and moves on to a smaller version of the problem. **It never goes back.**

That refusal to reconsider is what makes greedy fast. There's no branching, no backtracking, no exploring alternatives. If sorting is needed, it's usually O(n log n); otherwise it's often a single O(n) pass.

It's also what makes greedy dangerous. A choice that looks best in isolation can quietly rule out a much better overall answer, and the algorithm will never notice.

:::analogy
Greedy is walking up a hill by always stepping in the steepest upward direction. You'll definitely reach *a* summit. Whether it's the highest summit depends entirely on the shape of the landscape — and you'll never know, because you never look back down.
:::

## When greed is provably correct

A greedy algorithm is guaranteed correct when the problem has two properties.

**The greedy-choice property.** There is always an optimal solution that includes the locally-best choice. In other words, taking the greedy option never costs you the ability to reach an optimal answer. This is what the cut property proved for minimum spanning trees.

**Optimal substructure.** After you make that choice, what remains is a smaller instance of the same problem, and solving it optimally gives you an optimal solution overall. This is the same property that lets divide-and-conquer work.

Both must hold. Optimal substructure alone isn't enough — plenty of problems have it and still defeat greedy, because the greedy choice isn't safe.

:::key
Greedy is correct only when the locally-best choice is provably part of *some* globally-optimal answer. If you can't argue that, you can't trust the algorithm — no matter how many test cases pass.
:::

## A success: interval scheduling

Here's a problem where greedy shines. You have a list of activities, each with a start and finish time, and one room. Overlapping activities can't both run. **Select the largest possible number of activities.**

The tempting rules — take the shortest activity, or the one that starts earliest — both fail. The rule that works is: **always take the activity that finishes earliest** among those that still fit.

```c
#include <limits.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

typedef struct { int start, finish; } Activity;

static int by_finish(const void *a, const void *b) {
    const Activity *x = a, *y = b;
    return (x->finish > y->finish) - (x->finish < y->finish);
}

size_t select_activities(Activity activities[], size_t n, Activity chosen[]) {
    qsort(activities, n, sizeof *activities, by_finish);
    size_t count = 0;
    int last_finish = INT_MIN;
    for (size_t i = 0; i < n; ++i) {
        if (activities[i].start >= last_finish) {
            chosen[count++] = activities[i];
            last_finish = activities[i].finish;
        }
    }
    return count;
}

int main(void) {
    Activity acts[] = {{1,4}, {3,5}, {0,6}, {5,7}, {3,9}, {5,9}, {6,10}, {8,11}};
    Activity chosen[sizeof acts / sizeof acts[0]];
    size_t count = select_activities(acts, sizeof acts / sizeof acts[0], chosen);
    for (size_t i = 0; i < count; ++i) printf("(%d,%d)%s", chosen[i].start, chosen[i].finish, i + 1 == count ? "\n" : " ");
    return 0;
}
```

```text
sorted by finish:  (1,4) (3,5) (0,6) (5,7) (3,9) (5,9) (6,10) (8,11)

take (1,4)   -> room free from 4
skip (3,5)   starts at 3, still busy
skip (0,6)   starts at 0, still busy
take (5,7)   -> room free from 7
skip (3,9), (5,9), (6,10)   all start before 7
take (8,11)

3 activities — and no other selection can do better
```

Why is "earliest finish" the safe choice? Because it leaves the maximum amount of room free for everything that follows. Any optimal schedule's first activity can be swapped for this one without making the schedule any worse — that's the greedy-choice property, made concrete. The whole algorithm is a sort plus one linear pass: **O(n log n)**.

## A success: coin change with a canonical system

Now a classic. Given coin denominations and an amount, make that amount using the **fewest coins**. The greedy rule is obvious: always take the largest coin that still fits.

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

static int descending(const void *a, const void *b) {
    int x = *(const int *)a, y = *(const int *)b;
    return (y > x) - (y < x);
}

bool greedy_change(int coins[], size_t coin_count, int amount,
                   int used[], size_t capacity, size_t *used_count) {
    qsort(coins, coin_count, sizeof *coins, descending);
    size_t count = 0;
    for (size_t i = 0; i < coin_count; ++i) {
        if (coins[i] <= 0) return false;
        while (amount >= coins[i]) {
            if (count == capacity) return false;
            used[count++] = coins[i];
            amount -= coins[i];
        }
    }
    *used_count = count;
    return amount == 0;
}

int main(void) {
    int coins[] = {1, 5, 10, 25}, used[63]; size_t count;
    if (!greedy_change(coins, 4, 63, used, 63, &count)) return EXIT_FAILURE;
    for (size_t i = 0; i < count; ++i) printf("%d%s", used[i], i + 1 == count ? "\n" : " ");
    /* 25 25 10 1 1 1 */
    return 0;
}
```

Six coins for 63 cents, and that is genuinely optimal. Try it on any amount with these denominations and greedy always wins.

Coin systems where greedy always gives the optimal answer are called **canonical**. Most real currencies are canonical, which is not an accident — they were designed so that ordinary people making change by instinct get the right answer.

:::tip
Greedy is often the right *first* thing to try: it's quick to write and quick to test. Just don't ship it until you've either proved the greedy-choice property or convinced yourself with a genuine attempt to break it.
:::

## The failure that matters

Here is the heart of this lesson. Change the coin system to **{1, 3, 4}** and ask for **6**.

Greedy takes the largest coin that fits — a 4. That leaves 2, which needs two 1s.

```text
greedy:    6  ->  take 4  (remaining 2)
                  take 1  (remaining 1)
                  take 1  (remaining 0)
           = 4 + 1 + 1  ->  THREE coins

optimal:   3 + 3          ->  TWO coins
```

Greedy returns three coins. The optimal answer is two. Watch it happen:

```c
int coins[] = {1, 3, 4}, used[6];
size_t count;
greedy_change(coins, 3, 6, used, 6, &count);
/* used contains 4, 1, 1: three coins. */
```

Nothing went wrong mechanically. The code is correct; the *strategy* is wrong. The greedy-choice property simply doesn't hold for {1, 3, 4}: taking the 4 is locally best but destroys the optimal solution, because the optimal answer doesn't contain a 4 at all. And greedy, by definition, never goes back to find out.

This is the failure mode to carry with you. It's silent. There's no error, no crash, no warning — just an answer that's a little worse than it should be, on some inputs and not others.

:::warning
Greedy failures do not announce themselves. Calling `greedy_change` with coins `{1, 3, 4}` and amount `6` returns perfectly valid coins summing to 6 — just not the fewest coins. Testing on a few inputs will not catch this; only reasoning about the greedy-choice property will.
:::

The fix, when greedy fails, is usually **dynamic programming**: instead of committing to one choice, try every choice for the last coin and keep the best result, reusing the sub-answers you've already computed.

```c
#include <limits.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

int optimal_change(const int coins[], size_t coin_count, size_t amount) {
    if (amount == SIZE_MAX || amount + 1 > SIZE_MAX / sizeof(size_t)) return -1;
    size_t *best = malloc((amount + 1) * sizeof *best);
    if (best == NULL) return -1;
    best[0] = 0;
    for (size_t i = 1; i <= amount; ++i) {
        best[i] = SIZE_MAX;
        for (size_t j = 0; j < coin_count; ++j) {
            if (coins[j] > 0 && (size_t)coins[j] <= i &&
                best[i - (size_t)coins[j]] != SIZE_MAX &&
                best[i - (size_t)coins[j]] + 1 < best[i]) {
                best[i] = best[i - (size_t)coins[j]] + 1;
            }
        }
    }
    int answer = best[amount] == SIZE_MAX || best[amount] > INT_MAX
               ? -1 : (int)best[amount];
    free(best);
    return answer;
}

/* optimal_change((int[]){1, 3, 4}, 3, 6) is 2. */
```

That's O(amount × number of coins) — slower than greedy's near-instant answer, and correct for every coin system. You pay for correctness with work, which is the trade greedy was trying to avoid.

:::key
Greedy commits to one choice and is fast. Dynamic programming tries every choice but remembers the sub-answers so it doesn't repeat itself. When the greedy-choice property fails, that's the upgrade path.
:::

## The greedy algorithms you already know

You've met more greedy algorithms than you realised, and it's worth naming why each one is safe.

**Dijkstra's algorithm.** Always expands the nearest unfinished node and finalises it. Safe *only because* edge weights are non-negative — that's precisely the condition that makes the greedy choice provably correct. Allow negative edges and the greedy-choice property fails, which is exactly why Dijkstra breaks.

**Prim's algorithm.** Always adds the cheapest edge leaving the growing tree. Safe by the cut property.

**Kruskal's algorithm.** Always adds the cheapest edge that doesn't form a cycle. Safe by the same cut property, applied to a different division of the vertices.

**Huffman coding.** A compression algorithm that builds a code where common symbols get short bit patterns. It repeatedly merges the two least-frequent symbols into a subtree — a greedy choice, provably optimal, and the reason text compresses as well as it does.

Notice a pattern: every one of these has an *actual proof* behind it. That's the difference between a greedy algorithm that's famous and one that's a bug.

:::example
Dijkstra is the clearest case study on this list. It's greedy, it's optimal, and the moment you violate its one precondition — non-negative weights — it becomes exactly like greedy coin change on {1, 3, 4}: still fast, still confident, quietly wrong.
:::

## Check Your Understanding

:::quiz
Q: With coins {1, 3, 4}, how many coins does the greedy algorithm use to make 6, and how many are actually needed?
- Greedy uses 2, optimal is 2
- Greedy uses 3, optimal is 2 *
- Greedy uses 3, optimal is 3
- Greedy uses 6, optimal is 2
E: Greedy takes the 4 first and then needs two 1s, giving three coins. The optimal answer is 3 + 3, which is two coins — but greedy never considers skipping the 4.
:::

:::quiz
Q: Which two properties must a problem have for a greedy algorithm to be provably optimal?
- Sorted input and no duplicates
- The greedy-choice property and optimal substructure *
- A finite input and a base case
- Non-negative weights and a priority queue
E: The greedy-choice property says the locally-best choice is part of some optimal solution; optimal substructure says what remains after that choice is a smaller instance of the same problem.
:::

:::predict
Q: What does this print?
```c
int coins[] = {1, 4, 5}, used[8];
size_t count;
greedy_change(coins, 3, 8, used, 8, &count);
for (size_t i = 0; i < count; ++i) printf("%d%s", used[i], i + 1 == count ? "\n" : " ");
```
- `5 1 1 1` *
- `4 4`
- `5 4`
- No solution
E: Greedy takes the 5 first, leaving 3, which no 4 fits into — so it finishes with three 1s, four coins in total. The optimal answer, 4 + 4, uses only two. Another coin system where greed fails.
:::

## Talk about it

> Greedy algorithms fail silently: they return a valid-looking answer that just isn't the best one. Describe how you'd go about *convincing yourself* that a greedy rule is safe for a new problem — and why running a few test cases isn't enough on its own.

## What's next

When greed is wrong you need to explore alternatives — but exploring every possibility is usually far too slow. Next up is **Backtracking**, which explores systematically while abandoning any partial answer the moment it can't possibly work. That pruning is what turns an impossible search into a practical one, and you'll use it to solve the N-Queens puzzle.

# Space Complexity

Time isn't the only budget an algorithm spends. Memory is finite too, and running out of it doesn't make your program slow — it makes it crash. The good news is that you already know the technique: you count memory the same way you counted steps, as a function of n. This lesson shows you how, including one source of memory use that's easy to miss entirely.

## Counting memory instead of steps

Space complexity asks: **how much memory does this algorithm need, as the input grows?**

Same notation, same rules. Drop constants, drop lower-order terms, keep the dominant growth.

```c
#include <stddef.h>

long total(const int nums[], size_t n) {
    long result = 0;                    // one variable
    for (size_t i = 0; i < n; i++) {   // one loop variable
        result += nums[i];
    }
    return result;
}
```

However long `nums` is, this function creates only a fixed number of extra variables. Two isn't "small" or "big" — it's *constant*, so the extra memory is `O(1)`.

Now compare:

```c
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

int *doubled(const int nums[], size_t n) {
    int *out = malloc(n * sizeof *out);  // grows to n items
    if (out == NULL && n != 0) {
        return NULL;
    }
    for (size_t i = 0; i < n; i++) {
        out[i] = nums[i] * 2;
    }
    return out;
}

int main(void) {
    const int nums[] = {1, 2, 3};
    int *out = doubled(nums, 3);
    if (out == NULL) {
        return EXIT_FAILURE;
    }
    printf("[%d, %d, %d]\n", out[0], out[1], out[2]);  // [2, 4, 6]
    free(out);
    return 0;
}
```

The new array holds one entry per input item, so the extra memory is `O(n)`. Ten items in, ten out; a million in, a million out.

:::key
Space complexity counts memory that **grows with n**. A fixed handful of variables is O(1) no matter how many of them there are. A structure with one slot per input item is O(n).
:::

## Auxiliary space vs total space

Two different questions hide behind "how much memory does it use?"

**Total space** counts everything, including the input itself. Since an array of n items always occupies n slots, the total space of any algorithm that takes an array is at least `O(n)`. That makes total space a slightly boring measurement.

**Auxiliary space** counts only the *extra* memory the algorithm allocates on top of the input it was handed. This is almost always the interesting number, and it's what people mean when they say "space complexity" without qualification.

```text
algorithm                          auxiliary     total
--------------------------------   -----------   ---------
sum an array with a running total     O(1)          O(n)
build a doubled copy of an array      O(n)          O(n)
merge sort an array                   O(n)          O(n)
```

:::tip
When someone quotes a space complexity, check which one they mean. "Merge sort uses O(n) space" is about auxiliary space — the temporary arrays it merges into. Saying merge sort uses O(n) *total* space would be true but would fail to distinguish it from anything else.
:::

## What "in-place" means

An algorithm is **in-place** when it rearranges the input using only `O(1)` auxiliary space — a few variables, no second copy of the data.

Reversing a list out-of-place:

```c
#include <stddef.h>
#include <stdlib.h>

int *reverse_copy(const int nums[], size_t n) {
    int *out = malloc(n * sizeof *out);  // builds a whole new array
    if (out == NULL && n != 0) {
        return NULL;
    }
    for (size_t i = 0; i < n; i++) {
        out[n - 1 - i] = nums[i];
    }
    return out;
}
```

That's `O(n)` auxiliary space. In-place, using two indices that walk toward each other:

```c
#include <stddef.h>
#include <stdio.h>

void reverse_in_place(int nums[], size_t n) {
    size_t lo = 0;
    size_t hi = n;
    while (lo < hi) {
        hi--;
        int temp = nums[lo];
        nums[lo] = nums[hi];
        nums[hi] = temp;
        lo++;
    }
}

int main(void) {
    int nums[] = {1, 2, 3, 4};
    reverse_in_place(nums, 4);
    printf("[%d, %d, %d, %d]\n", nums[0], nums[1], nums[2], nums[3]);  // [4, 3, 2, 1]
    return 0;
}
```

A fixed number of variables, regardless of array length: `O(1)` auxiliary space. Note the cost — the original order is gone. In-place algorithms destroy their input, which is sometimes exactly what you want and sometimes a bug waiting to happen.

:::warning
"In-place" is about auxiliary space, not about whether you used the word `new`. And it's applied a little loosely in practice: quicksort is universally called in-place even though its recursion needs O(log n) stack space on average. The spirit is "no second copy of the data."
:::

## The call stack is memory too

Here's the one people forget. Every function call reserves a **stack frame** — space for its parameters, its local variables, and where to return to. Frames stack up and are only released when the calls return. So **recursion depth is space**.

```c
#include <stdio.h>

unsigned sum_recursive(unsigned n) {
    if (n == 0) {
        return 0;
    }
    return n + sum_recursive(n - 1);
}

int main(void) {
    printf("%u\n", sum_recursive(5));  // 15
    return 0;
}
```

This looks like it allocates nothing. But calling `sum_recursive(900)` means nine hundred frames are open at once before the first one returns: `O(n)` space. The loop version uses `O(1)`:

```c
#include <stdio.h>

unsigned sum_loop(unsigned n) {
    unsigned total = 0;
    for (unsigned i = 1; i <= n; i++) {
        total += i;
    }
    return total;
}

int main(void) {
    printf("%u\n", sum_loop(5));  // 15
    return 0;
}
```

Same answer, same `O(n)` time, completely different space. C does not guarantee a recursion-depth guard, so sufficiently deep recursion can exhaust the process's call stack and terminate the program.

```text
recursive shape                  stack depth     space
-----------------------------    -----------     ---------
T(n) = T(n-1) + ...              n               O(n)
T(n) = T(n/2) + ...              log n           O(log n)
merge sort (2 calls on n/2)      log n           O(log n) stack + O(n) buffers
quicksort, average               log n           O(log n)
quicksort, worst case            n               O(n)
```

Notice the last two rows: quicksort's bad case is bad for *both* budgets. If every partition is maximally lopsided, you get n levels of recursion, which is O(n²) time *and* O(n) stack.

:::analogy
The call stack is a pile of sticky notes. Each call writes a note saying "I was in the middle of this, come back here when you're done" and drops it on the pile. A recursion 100,000 deep is a pile of 100,000 notes, all of which must be kept until the deepest one finally answers.
:::

## Trading time for space

Very often you can spend memory to save time, or spend time to save memory. Recognising the trade is a real engineering skill.

**Memoisation buys time with memory.** Store answers you've already computed, so you never compute them twice.

```c
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

unsigned long long fib_memo_impl(size_t n, unsigned long long cache[]) {
    if (n < 2) {
        return n;
    }
    if (cache[n] == 0) {
        cache[n] = fib_memo_impl(n - 1, cache) + fib_memo_impl(n - 2, cache);
    }
    return cache[n];
}

unsigned long long fib_memo(size_t n) {
    unsigned long long *cache = calloc(n + 1, sizeof *cache);
    if (cache == NULL) {
        fprintf(stderr, "could not allocate Fibonacci cache\n");
        exit(EXIT_FAILURE);
    }
    unsigned long long result = fib_memo_impl(n, cache);
    free(cache);
    return result;
}

int main(void) {
    printf("%llu\n", fib_memo(50));  // 12586269025
    return 0;
}
```

Naive Fibonacci is exponential in time and `O(n)` in space (the stack). Memoised, it's `O(n)` time and `O(n)` space — the cache holds one entry per value of n. You spent linear memory to delete exponential time. Excellent deal.

**Recomputing buys memory with time.** The reverse trade shows up whenever data won't fit. Graphs are the classic case, and you've already built both structures:

```text
graph with V nodes and E edges
------------------------------------------------
adjacency matrix   O(V^2) space   O(1) edge check
adjacency list     O(V + E) space O(degree) edge check
```

For a social graph with V = 100,000 people, a matrix needs 10¹⁰ cells — ten billion, almost all of them empty. The list stores only the edges that exist. You accept a slower "are these two connected?" check in exchange for a structure that actually fits in memory.

:::example
Concrete numbers help. A C `int` is typically 4 bytes, so an array of one million ints is about 4 MB — comfortable. An n×n grid of ints with n = 1,000,000 would be 10¹² cells, roughly 4 terabytes. The same `O(n)` versus `O(n²)` gap you saw for time is just as brutal for memory, and it fails harder: slow code finishes eventually, but code that runs out of memory stops.
:::

## Check Your Understanding

:::quiz
Q: What is the auxiliary space complexity of a function that reverses a list in place using two index variables?
- O(n)
- O(log n)
- O(1) *
- O(n²)
E: It allocates a fixed number of variables regardless of list length, so the extra memory doesn't grow with n.
:::

:::quiz
Q: Why does a recursive function that recurses n levels deep use O(n) space even if it allocates no arrays?
- C copies the entire array at every call
- Each open call keeps a stack frame alive until it returns *
- Recursion always duplicates its input
- Because O(n) time always implies O(n) space
E: Every pending call holds a stack frame with its locals and return address. n pending calls means n frames alive at once.
:::

:::match
Q: Match each approach to its auxiliary space.
- Summing an array with a running total | O(1)
- Merge sort's temporary buffers | O(n)
- An adjacency matrix for V nodes | O(V²)
- Recursion that halves the input each call | O(log n)
E: Count only what grows with the input, and remember that pending recursive calls occupy stack space.
:::

## Talk about it

> Memoisation trades memory for speed, and it isn't always the right trade — a cache that grows without limit is a memory leak with good intentions. Describe a situation where you'd refuse to memoise even though it would make the code faster, and explain what you'd measure to make that call.

## What's next

So far every algorithm we've analysed had a known, reasonable running time. But some problems resist that entirely — nobody has ever found a fast algorithm for them, and nobody has proved one is impossible. In **Complexity Classes P and NP** you'll meet the most famous open question in computer science, why it matters for the encryption protecting your bank details, and what these much-misquoted letters actually mean. Pixel finds this the most interesting corner of the whole subject.

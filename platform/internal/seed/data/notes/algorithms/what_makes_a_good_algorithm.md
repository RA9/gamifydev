# What Makes a Good Algorithm

You've spent two courses learning how to *store* data and how to *measure* the cost of code. This course is about the third piece: the recipes themselves. An **algorithm** is a precise recipe for turning an input into an output, and by the end of this course you'll know the classic ones by name and be able to reason about which fits your problem.

This first lesson sets the ground rules. What actually counts as an algorithm, how do we judge one, and how do we write one down so another person (or another you, six months from now) can follow it? Pixel is along for the whole course.

## What an algorithm actually is

An algorithm is a finite sequence of unambiguous steps that takes some input, always terminates, and produces the correct output. Every word in that sentence is doing work:

- **Finite** — it's written down in a limited number of steps. Not "keep trying things."
- **Unambiguous** — each step means exactly one thing. "Pick a good pivot" is not a step; "pick the element at index 0" is.
- **Terminates** — it finishes for every valid input. A loop that might spin forever is not an algorithm.
- **Correct** — for every valid input, it produces the answer the problem asked for. Not usually. Every time.

Notice what's *not* in the list: speed. A slow algorithm is still an algorithm. Correctness is the entry fee; efficiency is what you optimise for afterwards.

```text
input  ->  [ a finite list of unambiguous steps ]  ->  output
           must terminate, must be correct
```

:::analogy
An algorithm is a cooking recipe. "Boil water, add 200g pasta, wait 10 minutes, drain" is finite, unambiguous, terminating and correct. "Cook until it feels right" is a technique, not a recipe — nobody else can follow it and get the same result.
:::

## Correctness, efficiency, readability

Once your algorithm is correct, three qualities compete for your attention, and they genuinely pull in different directions.

**Correctness** is non-negotiable. A fast wrong answer is worthless. This sounds obvious until you're tempted by a clever trick that works "for all the inputs I tried."

**Efficiency** is how the cost grows with input size — the Big O work from the last course. Time and space are separate budgets, and you can often trade one for the other. A lookup table spends memory to save time; recomputing values spends time to save memory.

**Readability** is how easily a human understands the code. This is not a nicety. Code that nobody can follow is code that nobody can fix, and most real bugs live in the parts people gave up reading.

The honest ordering for almost all real work: get it correct, get it readable, then make it fast **only where measurement says it matters**.

:::key
Correct first, clear second, fast third. An algorithm that's fast but wrong is useless; one that's fast but unreadable becomes wrong the moment someone edits it.
:::

## Specifications and invariants

Before you can say an algorithm is correct, you need to say what "correct" *means*. That statement is the **specification**: what you promise about the input, and what you promise about the output.

```text
Problem:  maximum of a list

Precondition:   nums is a list of numbers with at least one element
Postcondition:  returns a value m such that
                  m is in nums, and
                  m >= every element of nums
```

That's a contract. The precondition is what you demand from the caller; the postcondition is what you guarantee back. Notice the specification says *nothing* about how you'd do it — that freedom is exactly what lets you swap in a better algorithm later without breaking anything.

To argue that a loop actually meets its postcondition, you use an **invariant**: something true before the loop starts and still true after every single pass.

```c
#include <stddef.h>
#include <stdio.h>

/* Precondition: n > 0. */
int maximum(const int nums[], size_t n) {
    int best = nums[0];
    /* Invariant: best is the largest value in nums[0..i). */
    for (size_t i = 1; i < n; ++i) {
        if (nums[i] > best) best = nums[i];
    }
    return best;
}

int main(void) {
    int nums[] = {3, 9, 2, 9, 4};
    printf("%d\n", maximum(nums, sizeof nums / sizeof nums[0])); /* 9 */
    return 0;
}
```

The invariant holds at the start (`best` is the largest of the first one element — trivially true). Each pass keeps it true: if the new element is bigger, `best` becomes it; otherwise `best` was already the largest. When the loop ends, `i` has passed every index, so `best` is the largest of the whole list. That's a proof, not a hope.

:::tip
When a loop confuses you, write its invariant as a one-line comment above it. Half the time the act of writing that sentence reveals the bug, because you discover the sentence isn't actually true.
:::

## Writing an algorithm in pseudocode

Before you write code, write the idea in **pseudocode** — plain structured English, indented like code, with no language syntax to distract you. It's for thinking and for explaining, not for running.

```text
ALGORITHM linear_maximum(nums)
  best <- nums[0]
  FOR each remaining element x in nums
      IF x > best THEN
          best <- x
  RETURN best
```

Good pseudocode names its steps, shows its control flow, and skips everything irrelevant — no type declarations, headers, allocation details, or error handling. The examples in this course use C99, but sketching the language-neutral shape first keeps syntax and resource management from obscuring the algorithm.

:::warning
Pseudocode is allowed to be informal, but it is not allowed to be *vague*. "Sort the list somehow" hides the most expensive step in the whole algorithm. If a line hides real work, expand it — otherwise your cost analysis will be wrong.
:::

## Brute force is a legitimate starting point

**Brute force** means trying every possibility. It has a bad reputation it doesn't deserve. Suppose you want to know whether any two numbers in a list add up to a target:

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>

bool has_pair_bruteforce(const int nums[], size_t n, int target) {
    for (size_t i = 0; i < n; ++i) {
        for (size_t j = i + 1; j < n; ++j) {
            if ((int64_t)nums[i] + nums[j] == target) return true;
        }
    }
    return false;
}

int main(void) {
    int nums[] = {4, 1, 9, 7};
    size_t n = sizeof nums / sizeof nums[0];
    printf("%s\n", has_pair_bruteforce(nums, n, 16) ? "true" : "false");
    printf("%s\n", has_pair_bruteforce(nums, n, 20) ? "true" : "false");
    return 0;
}
```

Two nested loops over `n` items is O(n²). But look at what this version gives you: it's obviously correct, you can read it in ten seconds, and it's now a **reference implementation** you can test any faster version against.

And here is the faster version, which uses a hash set to remember what it has already seen:

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct {
    int *keys;
    bool *used;
    size_t capacity;
} IntSet;

static size_t hash_int(int value, size_t capacity) {
    return ((uint32_t)value * UINT32_C(2654435761)) % capacity;
}

static bool set_contains(const IntSet *set, int value) {
    size_t i = hash_int(value, set->capacity);
    while (set->used[i]) {
        if (set->keys[i] == value) return true;
        i = (i + 1) % set->capacity;
    }
    return false;
}

static void set_add(IntSet *set, int value) {
    size_t i = hash_int(value, set->capacity);
    while (set->used[i] && set->keys[i] != value) i = (i + 1) % set->capacity;
    set->used[i] = true;
    set->keys[i] = value;
}

bool has_pair_fast(const int nums[], size_t n, int target) {
    if (n == 0 || n > (SIZE_MAX - 1) / 2) return false;
    IntSet seen = {.capacity = 2 * n + 1};
    if (seen.capacity > SIZE_MAX / sizeof *seen.keys) return false;
    seen.keys = malloc(seen.capacity * sizeof *seen.keys);
    seen.used = calloc(seen.capacity, sizeof *seen.used);
    if (seen.keys == NULL || seen.used == NULL) {
        free(seen.keys);
        free(seen.used);
        return false;
    }

    bool found = false;
    for (size_t i = 0; i < n && !found; ++i) {
        int64_t complement = (int64_t)target - nums[i];
        if (complement >= INT32_MIN && complement <= INT32_MAX) {
            found = set_contains(&seen, (int)complement);
        }
        set_add(&seen, nums[i]);
    }
    free(seen.keys);
    free(seen.used);
    return found;
}
```

One pass, O(n) average time thanks to O(1) average set lookups, at the cost of O(n) extra memory. That's the whole game in miniature: start from something correct, understand where the cost is, spend memory or structure to remove it.

:::example
The brute force gives you a free test oracle. Generate random small lists, run both functions, and assert they agree. If they ever disagree, you've found a bug in the clever one — and you have the exact input that triggers it.
:::

## Check Your Understanding

:::quiz
Q: Which property is NOT part of the definition of an algorithm?
- It terminates for every valid input
- Its steps are unambiguous
- It runs in under one second *
- It produces correct output
E: Speed isn't part of the definition. A correct algorithm can be slow; efficiency is something you improve after correctness is secured.
:::

:::match
Q: Match each term to its meaning.
- Precondition | What must be true of the input before the algorithm runs
- Postcondition | What the algorithm guarantees about its output
- Invariant | Something that stays true after every pass of a loop
E: The pre/postconditions form the contract; the invariant is the argument that the loop actually honours it.
:::

:::quiz
Q: Why is a brute-force solution still worth writing?
- It is always the fastest option
- It is obviously correct and becomes a reference to test faster versions against *
- It uses no memory at all
- It avoids the need for a specification
E: Brute force is easy to get right and easy to read, which makes it an excellent baseline and test oracle for the optimised version you write next.
:::

## Talk about it

> Think about a task you do by hand — following a recipe, sorting laundry, finding a name in a long contact list. Try to write it as pseudocode with a precondition and a postcondition. Where did you find yourself writing something vague like "keep going until it looks done", and how would you make that step unambiguous?

## What's next

You now have the vocabulary for the whole course: specification, invariant, brute force as a baseline, and the correct-then-clear-then-fast ordering. Next up is **Recursion and Tail Recursion**, where a function calls itself to solve a smaller copy of the same problem. Nearly every algorithm in the second half of this course — merge sort, quick sort, tree traversal, backtracking — is built on that one idea, so it's worth getting comfortable with it now.

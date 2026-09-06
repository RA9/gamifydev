# Recursion and Tail Recursion

A **recursive** function is one that calls itself to solve a smaller copy of the same problem. It's the shape behind merge sort, quick sort, tree traversal, graph search and backtracking — nearly everything in the second half of this course. So it's worth more than a passing look.

In this lesson you'll write recursion deliberately, trace it by hand, understand exactly what it costs in memory, and learn what a **tail call** is and why some compilers can turn one back into a loop. C permits that optimisation, but does not guarantee it.

## Base case and recursive case

Every correct recursive function has two parts, and missing either one breaks it.

The **base case** is the input so small you can answer it directly, with no further recursion. It's what makes the process stop.

The **recursive case** calls the function on a strictly smaller input, moving one step closer to the base case.

```c
#include <stdio.h>

unsigned long long factorial(unsigned int n) {
    if (n == 0) return 1;               /* Base case. */
    return n * factorial(n - 1);        /* Recursive case on a smaller n. */
}

int main(void) {
    printf("%llu\n", factorial(5));     /* 120 */
    return 0;
}
```

The word "smaller" is doing all the work. Because `n` shrinks by one each call, it must eventually hit `0`. If the input doesn't shrink *toward* the base case, the recursion never ends.

:::key
Two ingredients, always: a **base case** that returns without recursing, and a **recursive case** on a smaller input that provably moves toward it.
:::

## The call stack is your space cost

When a function calls another function, the computer has to remember where it was — which local variables it held and which line to return to. That bundle is a **stack frame**, and frames pile up on the **call stack**.

Recursion piles frames up on *itself*. Here's `factorial(3)` mid-flight, at the moment it reaches the base case:

```text
   +---------------------------+
   | factorial(0)  -> returns 1|   <- top of stack, base case
   +---------------------------+
   | factorial(1)  waiting     |
   +---------------------------+
   | factorial(2)  waiting     |
   +---------------------------+
   | factorial(3)  waiting     |   <- bottom, the original call
   +---------------------------+
```

Four frames alive at once. That is a real memory cost: a recursion of depth `d` uses O(d) space even if it allocates nothing else. This is why "recursion is elegant" and "recursion is free" are two very different claims.

Once the base case returns, the answers cascade back down the stack:

```text
factorial(3)
= 3 * factorial(2)
= 3 * (2 * factorial(1))
= 3 * (2 * (1 * factorial(0)))
= 3 * (2 * (1 * 1))        # base case returned
= 3 * (2 * 1)
= 3 * 2
= 6
```

Read it top to bottom: calls stack *up* on the way down, then values fold *back* on the way up. Every recursion you meet in this course follows that motion.

:::warning
C does not specify a recursion-depth limit or provide a standard exception for exhausting the call stack. A recursion that is too deep can terminate the program or cause undefined behaviour, so a recursive walk over a million-item linked list is unsafe where a loop would use constant stack space.
:::

## Tracing a small recursion

Let's trace something with a return value that isn't just multiplication. This one reverses a string:

```c
#include <stddef.h>
#include <stdio.h>
#include <string.h>

void reverse_into(const char *s, size_t length, char out[]) {
    if (length == 0) {                  /* Base case: terminate the result. */
        out[0] = '\0';
        return;
    }
    reverse_into(s + 1, length - 1, out);
    out[length - 1] = s[0];             /* Attach the first character last. */
    out[length] = '\0';
}

int main(void) {
    const char *word = "pixel";
    char reversed[6];                   /* Five characters plus '\0'. */
    reverse_into(word, strlen(word), reversed);
    printf("%s\n", reversed);           /* lexip */
    return 0;
}
```

Walk it by hand:

```text
reverse("cat")
  -> reverse("at") + "c"
       -> reverse("t") + "a"
            -> "t"                  base case
       -> "t" + "a"  = "ta"
  -> "ta" + "c" = "tac"
```

Each level peels one character off the front and re-attaches it at the *end* on the way back up. The reversal happens entirely during the unwinding.

:::tip
When a recursion confuses you, trace it on paper for an input of size 3. Write one call per line, indent each deeper call, then fill in the return values bottom-up. This removes almost all the mystery, and it's exactly how you'll debug merge sort later.
:::

## What a tail call is

Look closely at these two functions. They compute the same thing, but they are structurally different:

```c
#include <stdio.h>

unsigned long long fact_normal(unsigned int n) {
    if (n == 0) return 1;
    return n * fact_normal(n - 1);      /* Multiply after the call returns. */
}

unsigned long long fact_tail(unsigned int n, unsigned long long acc) {
    if (n == 0) return acc;
    return fact_tail(n - 1, acc * n);   /* The call is the whole return. */
}

int main(void) {
    printf("%llu %llu\n", fact_normal(5), fact_tail(5, 1)); /* 120 120 */
    return 0;
}
```

In `fact_normal`, the recursive call is not the last thing that happens — there's still a multiplication waiting for its result. The frame must stay alive to do that multiplication.

In `fact_tail`, the recursive call is the *entire* return expression. Nothing is waiting. The running total is carried forward in an extra parameter called an **accumulator**. That's a **tail call**: a recursive call in the final position, with no work left to do afterwards.

:::key
A call is a **tail call** when its result is returned directly, with no pending work in the caller. The caller's frame has nothing left to do, so in principle it could be discarded rather than kept.
:::

## Tail call optimisation (and what C guarantees)

Because a tail-calling frame has nothing left to do, a compiler can *reuse* it instead of pushing a new one. That transformation is **tail call optimisation** (TCO), and it turns a recursion of depth `n` from O(n) stack space into O(1) — effectively a loop. Scheme requires it. C compilers often apply it when optimisation is enabled, but the C99 standard does not require them to.

```c
unsigned long long fact_tail(unsigned int n, unsigned long long acc) {
    if (n == 0) return acc;
    return fact_tail(n - 1, acc * n);
}

/* Do not assume a large call such as fact_tail(5000, 1) uses constant stack. */
```

Writing it in tail form may let a C compiler optimise it, but portable code cannot rely on that. It is still worth *recognising* the shape, because a tail-recursive function is always the easiest kind to convert into a loop by hand.

:::warning
Don't assume "it's tail recursive, so it's safe." That is only true in a language or implementation that guarantees the optimisation. C99 does not; if stack usage matters, write the loop explicitly.
:::

## Converting recursion to iteration

A tail-recursive function converts to a loop mechanically. The accumulator becomes a variable, and the recursive call becomes a reassignment:

```c
#include <stdio.h>

unsigned long long fact_loop(unsigned int n) {
    unsigned long long acc = 1;
    while (n > 0) {
        acc *= n;          /* Same update as the accumulator argument. */
        --n;               /* Same update as the n argument. */
    }
    return acc;
}

int main(void) {
    printf("%llu\n", fact_loop(5)); /* 120 */
    return 0;
}
```

Non-tail recursion — where work waits on the way back up — can also be converted, but you have to supply the stack yourself with an explicit array or linked structure. That's exactly what you'll do in the tree traversal lesson.

:::example
Reversing a mutable C string iteratively needs no stack or extra allocation at all: keep one index at each end and swap those characters while the indices move inward. For `"cat"`, swapping `c` and `t` produces `"tac"`.
:::

## Check Your Understanding

:::predict
Q: What does this print?
```c
#include <stdio.h>

const char *count_down(int n) {
    if (n <= 0) return "go";
    return count_down(n - 2);
}

int main(void) {
    printf("%s\n", count_down(5));
    return 0;
}
```
- go *
- The program exhausts the call stack
- 1
- NULL
E: 5 -> 3 -> 1 -> -1, and `-1 <= 0` is true, so the base case catches it and returns `"go"`. Using `<=` instead of `== 0` is what saves it from overshooting.
:::

:::quiz
Q: What is the space cost of a recursion that reaches depth d?
- O(1), recursion uses no extra memory
- O(d), one stack frame per pending call *
- O(d²)
- It depends only on the size of the return value
E: Each pending call holds a stack frame with its own locals and return address, so depth d means O(d) stack space.
:::

:::quiz
Q: Which call is a tail call?
- `return n * f(n - 1)`
- `return f(n - 1, acc * n)` *
- `x = f(n - 1); return x + 1`
- `return f(n - 1) + f(n - 2)`
E: A tail call's result is returned directly with no pending work. The others all have a multiplication, an addition, or a second call waiting to happen.
:::

## Talk about it

> C leaves tail-call optimisation to the compiler, so portable code cannot depend on it. Describe a situation where relying on an optional optimisation would be acceptable — and one where you would insist on predictable resource usage instead.

## What's next

You can now write recursion on purpose, name what it costs, spot a tail call, and unwind a simple recursion into a loop. Next up is **Linear and Binary Search**, your first pair of real algorithms — one that works on anything, and one that's dramatically faster but demands sorted input. Binary search is also a beautiful little recursion, so you'll get to use everything from this lesson right away.

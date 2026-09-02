# Recursion and Tail Recursion

A **recursive** function is one that calls itself to solve a smaller copy of the same problem. It's the shape behind merge sort, quick sort, tree traversal, graph search and backtracking — nearly everything in the second half of this course. So it's worth more than a passing look.

In this lesson you'll write recursion deliberately, trace it by hand, understand exactly what it costs in memory, and learn what a **tail call** is and why some languages can turn one back into a loop. Python, as you'll see, is not one of those languages.

## Base case and recursive case

Every correct recursive function has two parts, and missing either one breaks it.

The **base case** is the input so small you can answer it directly, with no further recursion. It's what makes the process stop.

The **recursive case** calls the function on a strictly smaller input, moving one step closer to the base case.

```python
def factorial(n):
    if n == 0:                      # base case
        return 1
    return n * factorial(n - 1)     # recursive case, on a smaller n

print(factorial(5))   # 120
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
Python caps recursion depth at roughly 1000 frames by default and raises `RecursionError` past it. That cap is a feature — it catches runaway recursion early — but it also means a recursive walk over a million-item linked list will crash where a loop would not.
:::

## Tracing a small recursion

Let's trace something with a return value that isn't just multiplication. This one reverses a string:

```python
def reverse(s):
    if len(s) <= 1:            # base case: "" and "a" are their own reverse
        return s
    return reverse(s[1:]) + s[0]

print(reverse("pixel"))   # lexip
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

```python
def fact_normal(n):
    if n == 0:
        return 1
    return n * fact_normal(n - 1)      # multiply AFTER the call returns

def fact_tail(n, acc=1):
    if n == 0:
        return acc
    return fact_tail(n - 1, acc * n)   # the call IS the whole return

print(fact_normal(5), fact_tail(5))    # 120 120
```

In `fact_normal`, the recursive call is not the last thing that happens — there's still a multiplication waiting for its result. The frame must stay alive to do that multiplication.

In `fact_tail`, the recursive call is the *entire* return expression. Nothing is waiting. The running total is carried forward in an extra parameter called an **accumulator**. That's a **tail call**: a recursive call in the final position, with no work left to do afterwards.

:::key
A call is a **tail call** when its result is returned directly, with no pending work in the caller. The caller's frame has nothing left to do, so in principle it could be discarded rather than kept.
:::

## Tail call optimisation (and why Python has none)

Because a tail-calling frame has nothing left to do, a compiler can *reuse* it instead of pushing a new one. That transformation is **tail call optimisation** (TCO), and it turns a recursion of depth `n` from O(n) stack space into O(1) — effectively a loop. Scheme requires it. Several functional languages guarantee it, and some compilers for C-family languages apply it as an optimisation.

Python deliberately does not do this. The reasoning is that discarding those frames also discards the stack trace, and Python's designers decided that readable tracebacks are worth more than free recursion depth. So in Python:

```python
def fact_tail(n, acc=1):
    if n == 0:
        return acc
    return fact_tail(n - 1, acc * n)

# fact_tail(5000)  ->  RecursionError: maximum recursion depth exceeded
```

Writing it in tail form buys you nothing in Python. It's still worth *recognising* the shape, because a tail-recursive function is always the easiest kind to convert into a loop by hand.

:::warning
Don't assume "it's tail recursive, so it's safe." That's only true in a language that guarantees the optimisation. In Python, a tail-recursive function of depth 5000 crashes exactly like any other.
:::

## Converting recursion to iteration

A tail-recursive function converts to a loop mechanically. The accumulator becomes a variable, and the recursive call becomes a reassignment:

```python
def fact_loop(n):
    acc = 1
    while n > 0:
        acc = acc * n      # same update as the accumulator argument
        n = n - 1          # same update as the n argument
    return acc

print(fact_loop(5))       # 120
print(fact_loop(5000) == fact_loop(5000))   # True — no stack limit at all
```

Non-tail recursion — where work waits on the way back up — can also be converted, but you have to supply the stack yourself with an explicit list. That's exactly what you'll do in the tree traversal lesson.

:::example
Reversing a string iteratively needs no stack at all, because you can build the answer as you go: `out = ""` then `for ch in s: out = ch + out`. For `"cat"` that builds `"c"`, then `"ac"`, then `"tac"`.
:::

## Check Your Understanding

:::predict
Q: What does this print?
```python
def count_down(n):
    if n <= 0:
        return "go"
    return count_down(n - 2)

print(count_down(5))
```
- go *
- RecursionError
- 1
- None
E: 5 -> 3 -> 1 -> -1, and `-1 <= 0` is true, so the base case catches it and returns "go". Using `<=` instead of `== 0` is what saves it from overshooting.
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

> Python refuses to optimise tail calls so that error tracebacks stay complete and readable. That's a deliberate trade of performance for debuggability. Describe another place in programming where you'd happily give up speed to make failures easier to understand — and one where you wouldn't.

## What's next

You can now write recursion on purpose, name what it costs, spot a tail call, and unwind a simple recursion into a loop. Next up is **Linear and Binary Search**, your first pair of real algorithms — one that works on anything, and one that's dramatically faster but demands sorted input. Binary search is also a beautiful little recursion, so you'll get to use everything from this lesson right away.

# Recursion

Recursion is when a function calls *itself* to solve a problem. It sounds strange the first time — a function that uses its own name inside its body? — but it turns out to be one of the most elegant tools in a programmer's kit. Some problems are naturally shaped like a smaller copy of themselves, and recursion lets your code mirror that shape directly.

You already know loops, so you can technically solve everything without recursion. But once you learn to *see* the recursive shape in a problem, certain tasks — walking a branching dungeon, adding up a skill tree, exploring nested data — become dramatically shorter to write.

## What Recursion Is

Picture a dungeon made of rooms. Every room has a door that leads deeper into the dungeon, and eventually you hit a dead-end room with treasure. How do you explore it? "Walk into the room. If it's the treasure room, stop. Otherwise, explore the *next* room the same way." That last sentence — "explore the next room the same way" — is recursion. The instruction refers back to itself.

In code, a recursive function is simply a function that calls itself:

```python
def explore(room):
    print(f"Entered {room}")
    if room == 0:
        print("Found the treasure!")
        return
    explore(room - 1)  # the function calls itself

explore(3)
# Entered 3
# Entered 2
# Entered 1
# Entered 0
# Found the treasure!
```

Each call handles one room and then hands the rest of the job to another call of the same function, working on a slightly smaller problem.

:::analogy
Recursion is like a row of nesting dolls. To open the whole set you open one doll, and inside is the *same task* on a smaller doll — until you reach the tiny solid doll that doesn't open. That tiny doll is where the process finally stops.
:::

## The Two Essential Parts

Every correct recursive function has exactly two ingredients. Miss either one and it breaks.

**1. The base case** — the situation so simple you can answer it directly, with no more recursion. In the dungeon, that's the treasure room: you stop. The base case is what makes the process *end*.

**2. The recursive case** — where the function calls itself on a *smaller* version of the problem, moving one step closer to the base case.

```python
def countdown(n):
    if n == 0:          # base case: stop here
        print("Liftoff!")
        return
    print(n)
    countdown(n - 1)    # recursive case: smaller problem

countdown(3)
# 3
# 2
# 1
# Liftoff!
```

The magic is in "smaller." Each call to `countdown` gets an `n` that is one less than before. Because `n` keeps shrinking, it *must* eventually reach `0`, and then the base case ends everything.

:::key
Two parts, always: a **base case** that stops, and a **recursive case** that moves toward it on a smaller input. If the input doesn't shrink toward the base case, recursion never ends.
:::

## Tracing Factorial Step by Step

The classic first example is the factorial. `5!` (read "five factorial") means `5 * 4 * 3 * 2 * 1`. Notice the recursive shape hiding inside it: `5!` is just `5 * 4!`, and `4!` is `4 * 3!`, and so on. Each factorial is defined in terms of a smaller factorial.

```python
def factorial(n):
    if n == 0:           # base case: 0! is defined as 1
        return 1
    return n * factorial(n - 1)  # recursive case

print(factorial(5))  # 120
```

Let's trace exactly what happens when we call `factorial(3)`. Python has to fully evaluate the inner call before it can multiply:

```text
factorial(3)
= 3 * factorial(2)
= 3 * (2 * factorial(1))
= 3 * (2 * (1 * factorial(0)))
= 3 * (2 * (1 * 1))        # base case reached, returns 1
= 3 * (2 * 1)
= 3 * 2
= 6
```

Read it top to bottom: the calls *stack up*, each waiting on the one below it. When `factorial(0)` finally returns `1`, the answers cascade back up — `1`, then `1`, then `2`, then `6`. That "go all the way down, then come back up" motion is the heart of how recursion runs.

:::tip
When a recursive function confuses you, trace it by hand on paper for a small input like `n = 3`. Write each call on its own line and substitute the returned value back in. Seeing the stack build and unwind removes almost all the mystery.
:::

## Classic Examples

**Summing a flat list.** A list's sum is its first item plus the sum of everything else — a smaller list. The base case is the empty list, whose sum is `0`.

```python
def list_sum(nums):
    if not nums:              # base case: empty list
        return 0
    return nums[0] + list_sum(nums[1:])  # first + sum of the rest

print(list_sum([10, 20, 5, 30]))  # 65
```

Each call peels off `nums[0]` and hands the shorter `nums[1:]` to the next call, until nothing is left.

**Fibonacci.** In this sequence each number is the sum of the previous two: `0, 1, 1, 2, 3, 5, 8, 13...`. It has *two* base cases and *two* recursive calls.

```python
def fib(n):
    if n < 2:        # base cases: fib(0)=0, fib(1)=1
        return n
    return fib(n - 1) + fib(n - 2)

print([fib(i) for i in range(8)])  # [0, 1, 1, 2, 3, 5, 8, 13]
```

Fibonacci reads beautifully, but be warned: this naive version recomputes the same values over and over. `fib(30)` triggers over a million calls. It's a lovely teaching example and a terrible way to actually compute large Fibonacci numbers — a plain loop is far faster here.

:::predict
Q: What does this return?
```python
def list_sum(nums):
    if not nums:
        return 0
    return nums[0] + list_sum(nums[1:])

print(list_sum([4, 4, 2]))
```
- 10 *
- 442
- 0
- error
E: It adds the first item to the sum of the rest: 4 + (4 + (2 + 0)) = 10.
:::

## Recursion vs. Iteration

Anything you can do with recursion you can also do with a loop, and vice versa. So which should you reach for?

**Prefer a loop when** the problem is a straight, flat repetition — counting, running totals, walking a single list. Loops are usually faster and use less memory, because each recursive call takes a slot on Python's *call stack* while a loop reuses the same few variables. The factorial and countdown above are honestly clearer as loops:

```python
def factorial_loop(n):
    result = 1
    for i in range(2, n + 1):
        result *= i
    return result

print(factorial_loop(5))  # 120
```

**Prefer recursion when** the data itself branches or nests — trees, folders inside folders, a skill tree where each skill unlocks more skills, or any structure of "unknown depth." Trying to loop over something that branches unpredictably forces you to manage your own stack of "places still to visit," which is exactly the bookkeeping recursion does for you automatically.

:::key
Rule of thumb: flat and repetitive → loop. Branching, nested, or "unknown depth" → recursion. Choose whichever makes the code read like the problem.
:::

## When the Base Case Goes Wrong

Because each call adds a frame to the call stack, a recursion that never reaches its base case doesn't just spin forever like a `while True` loop — it piles up frames until Python runs out of room and raises a `RecursionError`.

The most common mistake is forgetting the base case entirely:

```python
def broken_countdown(n):
    print(n)
    broken_countdown(n - 1)   # no base case — never stops!

# broken_countdown(3)
# 3, 2, 1, 0, -1, -2, ...
# RecursionError: maximum recursion depth exceeded
```

The second common mistake is a base case that the recursion can never actually *hit* because the input isn't moving toward it:

```python
def wrong_step(n):
    if n == 0:
        return
    wrong_step(n - 2)   # skips right past 0 for odd n!

# wrong_step(5) -> 5, 3, 1, -1, -3, ... never equals 0 -> RecursionError
```

Here `5 -> 3 -> 1 -> -1` sails straight past `0` and never satisfies the base case. The fix is either a more forgiving condition (`if n <= 0`) or a step that truly converges on the base case.

:::warning
A `RecursionError` almost always means one of two things: you forgot the base case, or your recursive call isn't shrinking the input *toward* that base case. Check both before blaming Python. (Python caps recursion depth around 1000 by default, on purpose, to catch these bugs.)
:::

## Walking a Nested Structure

Here's where recursion truly earns its place. Imagine a skill tree stored as a nested list: some entries are plain skill costs, and some are sub-branches (lists) containing more costs, nested to any depth. You want the total cost of the whole tree.

A normal loop can't handle this cleanly, because you don't know ahead of time how deep the nesting goes. Recursion handles it naturally: for each item, either it's a number (add it) or it's a list (recurse into it).

```python
def deep_sum(tree):
    total = 0
    for item in tree:
        if isinstance(item, list):
            total += deep_sum(item)   # branch: sum the sub-tree
        else:
            total += item             # leaf: a plain cost
    return total

skill_tree = [10, [20, 5, [1, 1]], 30, [2, [3, 4]]]
print(deep_sum(skill_tree))  # 86
```

Follow the structure: `deep_sum` hits `10` (a leaf, add it), then `[20, 5, [1, 1]]` (a branch, so it calls itself on that sub-list), and so on all the way down. The base case here is implicit — a list with no sub-lists simply never triggers another recursive call, and the loop ends. Each level of nesting becomes one level of recursive calls, which is exactly why this fits in a handful of lines instead of a tangle of manual loops.

:::example
Swap the numbers for `isinstance(item, list)` checks in your head with `[1, [2, [3]]]`: `deep_sum` returns `1 + deep_sum([2, [3]])` = `1 + (2 + deep_sum([3]))` = `1 + (2 + 3)` = `6`. Each nested list peels open one recursive layer, just like the dungeon rooms.
:::

:::quiz
Q: Which pair of ingredients must every correct recursive function have?
- A base case and a recursive case that shrinks toward it *
- A loop and a counter
- Two base cases and no recursive call
- A global variable and a return statement
E: The base case stops the recursion; the recursive case calls the function on a smaller input that moves toward the base case. Without both, it either never stops or never recurses.
:::

## Recap

- Recursion is a function calling itself to solve a smaller version of the same problem.
- Every recursive function needs a **base case** (where it stops) and a **recursive case** (a self-call on a smaller input that moves toward the base case).
- Trace by hand for small inputs: calls stack up going down, then answers cascade back up (as in `factorial`).
- Classic patterns: countdown, summing a list, Fibonacci — though a loop is often faster for flat, repetitive work.
- Prefer loops for flat repetition; prefer recursion for branching or nested structures of unknown depth.
- A missing or unreachable base case causes a `RecursionError` — always confirm the input shrinks toward the base case.
- Recursion shines on nested data, like totaling a skill tree with `deep_sum`, where each layer of nesting maps to one layer of recursive calls.

**Next up:** exploring how recursion powers real data structures like trees and file systems — and when to reach for iteration instead.

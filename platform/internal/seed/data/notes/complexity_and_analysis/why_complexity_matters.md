# Why Complexity Matters

You've built arrays, linked lists, hash tables and trees. Now comes the question that decides which one you actually reach for: *how fast is it, and how does that answer change when the data gets bigger?* This lesson is about learning to answer that without a stopwatch. By the end you'll understand why "it ran in 2 seconds on my laptop" tells you almost nothing, and what to measure instead.

## The stopwatch lies

Imagine you write a function, time it, and report: "it takes 2 seconds." Pixel asks a fair question — 2 seconds on *what*?

That number quietly depends on things that have nothing to do with your code:

- **The machine.** A five-year-old laptop and a modern server can differ by a large factor on the exact same code.
- **What else is running.** A backup job, a browser with forty tabs, another process hogging the disk — all of it steals time.
- **The language and runtime.** The same algorithm written in C and in Python can differ enormously in raw speed.
- **The input you happened to test.** Two seconds on 1,000 records says nothing about 10 million records.

So a timing is a fact about *one run, on one machine, on one input*. It isn't a property of your algorithm. If you and a teammate report different seconds for the same code, you've both measured your computers, not your idea.

:::warning
Timings aren't useless — they're essential for finding real bottlenecks in real programs. They're just not *portable*. Never compare two algorithms by comparing two numbers measured on two different machines with two different inputs.
:::

## Count operations instead

Here's the move that makes the problem tractable. Instead of asking "how many seconds?", ask **"how many steps, as a function of the input size?"**

We call the input size **n**. It's whatever "bigger" means for your problem: the number of items in a list, the number of nodes in a tree, the number of characters in a string.

```python
def total(scores):
    result = 0                 # 1 step
    for s in scores:           # runs n times
        result = result + s    # 1 step each time
    return result              # 1 step

print(total([10, 20, 30]))
# 60
```

Count it: roughly `n` additions, plus a couple of steps at the ends. Call it `n + 2` steps. That expression is a fact about the *algorithm*. It's true on your laptop, on a server, in C, in Python. The seconds change; the shape `n + 2` does not.

:::key
Complexity analysis replaces "how long did it take?" with "how many steps does it take, expressed in terms of the input size n?" The first is a property of your afternoon. The second is a property of your algorithm.
:::

## Growth rate is what survives

Different machines change the *constant* in front of the count. A fast machine might do each step in 1 nanosecond, a slow one in 10. That's a factor of 10 — annoying, but fixed. It doesn't get worse as your data grows.

What *does* get worse as your data grows is the shape of the count. `n + 2` and `n²/2` behave completely differently once n is large, and no amount of faster hardware rescues the second one.

```text
input size n     n + 2 steps        n^2 / 2 steps
------------     -----------        -------------
        10                12                   50
       100               102                5,000
     1,000             1,002              500,000
    10,000            10,002           50,000,000
```

Multiply n by 10 and the first column multiplies by about 10. The second column multiplies by 100. That relationship — *what happens to the work when the input grows* — is called the **growth rate**, and it's the thing worth naming, comparing and arguing about.

:::analogy
Buying a faster computer is like widening a road: everything moves a bit quicker. Choosing a better growth rate is like replacing the road with a train line. The road helps until traffic doubles. The train line still works when traffic is a hundred times bigger.
:::

## A story: the duplicate usernames

You're asked to find whether a list of usernames contains any duplicates. The obvious approach is to compare every name with every later name:

```python
def has_duplicate_pairs(names):
    for i in range(len(names)):
        for j in range(i + 1, len(names)):
            if names[i] == names[j]:
                return True
    return False

print(has_duplicate_pairs(["ana", "bo", "ana"]))
# True
```

For a list of n names this does at most `n(n-1)/2` comparisons — every pair, once. That's roughly `n²/2`, so we call it a **quadratic** approach.

At n = 100 that's `100 × 99 / 2 = 4,950` comparisons. Nothing. It finishes before you lift your finger off the Enter key, and you ship it.

Then the product succeeds and the list grows to 1,000,000 names. Now it's `1,000,000 × 999,999 / 2`, which is about **500,000,000,000** comparisons — five hundred billion. As a thought experiment, suppose a machine that manages 100 million comparisons per second (a round number, not a measurement). Five hundred billion divided by a hundred million is 5,000 seconds: about **1 hour and 23 minutes**, for a check that felt instant in testing.

Now the hash-table version, using a set:

```python
def has_duplicate_set(names):
    seen = set()
    for name in names:
        if name in seen:
            return True
        seen.add(name)
    return False

print(has_duplicate_set(["ana", "bo", "ana"]))
# True
```

Each name is added and looked up once, and set lookups are fast on average — so this is roughly `n` steps, not `n²/2`. At n = 1,000,000 that's a million operations: about **0.01 seconds** at the same imaginary rate. Same problem, same machine, same language. The difference is entirely the growth rate.

:::example
Same two algorithms, three input sizes, one machine:

```text
n            pairwise (~n^2/2)      set-based (~n)
--------     -----------------      --------------
100                      4,950                 100
10,000              49,995,000              10,000
1,000,000      ~500,000,000,000           1,000,000
```
The gap isn't a constant factor. It widens every time n grows.
:::

## What "scales" actually means

People say software "scales" as if it means "is fast." It doesn't. **Scaling is about how the cost changes when the input grows**, not about how fast it is today.

A program that takes 5 seconds on 1,000 items and 10 seconds on 2,000 items scales well — it's slow, but doubling the data doubled the time. A program that takes 0.01 seconds on 1,000 items and 40 seconds on 20,000 items does not scale, no matter how snappy it felt at first.

This is why complexity analysis is a *design* tool rather than a *tuning* tool. You use it before you write the code, to choose the shape of the solution, because a bad growth rate can't be optimised away later. You can make each step 30% cheaper. You cannot make `n²` behave like `n`.

:::tip
When you're deciding between two approaches, ask "what happens at 100× the data?" before you ask "which one is faster right now?" The first question is the one that will still matter next year.
:::

## Check Your Understanding

:::quiz
Q: Why is "my function ran in 2 seconds" a poor way to describe an algorithm's efficiency?
- Seconds are too small a unit to be meaningful
- It depends on the machine, the load, the language and the specific input *
- Algorithms should always be measured in minutes
- Because timing code is technically impossible
E: A timing describes one run on one machine with one input. Change any of those and the number changes, so it isn't a property of the algorithm itself.
:::

:::quiz
Q: An algorithm does about n²/2 steps. Roughly how much more work does it do when n grows from 1,000 to 10,000?
- About 10 times more
- About 100 times more *
- About the same
- About 20 times more
E: n grew 10×, and squaring that gives 10 × 10 = 100. Quadratic work grows by the square of the input growth.
:::

:::fill
Q: Complete the term for the standard variable naming the input size.
`We count steps as a function of the input size, written ___.`
- n *
- t
- k
E: By convention `n` names the size of the input — the number of items, nodes or characters you're working with.
:::

## Talk about it

> Think of a program or app you use that felt fast when you started using it and became slow as your data grew — a photo library, a notes app, a chat history. What do you think grew, and what operation do you suspect was doing work proportional to *every* item you'd ever stored? Describe what you'd measure to find out.

## What's next

You now have the central idea: describe an algorithm by how its work *grows* with the input, not by how many seconds it took on your machine. What we've been doing informally — saying "roughly n" or "roughly n²/2" — needs a precise, agreed-upon notation so that everyone means the same thing. That notation is the subject of the next lesson, **Big O Notation**, where you'll learn exactly why `n²/2` and `3n² + 5n + 100` get the same name.

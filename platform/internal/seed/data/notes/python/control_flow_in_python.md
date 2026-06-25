# Control Flow in Python

Up to now, your programs have run straight down the page, one line after another, the same way every time. Real programs aren't like that — they *react*. They check conditions, take different paths, and repeat work without complaint. This is **control flow**: the art of deciding *what* runs and *how many times*. It's where code starts to feel intelligent.

Keep the [Code Lab](#code) open. This is a lesson you learn by running, not just reading.

## Making decisions with if

An `if` statement runs a block of code *only when* a condition is true. The condition is anything that evaluates to `True` or `False`.

```python
temperature = 30
if temperature > 25:
    print("It's warm out — wear a t-shirt.")
```

Read it like a sentence: "if the temperature is greater than 25, then print this." The indented line belongs to the `if` — remember, in Python, indentation is the structure. If the condition is `False`, Python simply skips the indented block.

Often you want a fallback. `else` covers "when the condition was *not* true," and `elif` (short for "else if") lets you check additional conditions in order:

```python
score = 72
if score >= 90:
    print("Grade: A")
elif score >= 80:
    print("Grade: B")
elif score >= 70:
    print("Grade: C")
else:
    print("Keep practicing!")
```

Python checks each condition top to bottom and runs the **first** one that's true, then skips the rest. Order matters: because `72` passes `>= 70`, it prints "Grade: C" and never reaches the `else`.

:::analogy
Think of `if`/`elif`/`else` as a fork in a hiking trail with several signposts. You read each sign in order and take the *first* path that applies to you. Once you've chosen a path, you don't double back to read the others.
:::

## Comparisons and logic

Conditions are built from **comparison operators** that ask yes/no questions:

```python
print(5 == 5)   # True   — equal to (note the double =)
print(5 != 3)   # True   — not equal to
print(5 > 3)    # True   — greater than
print(5 <= 5)   # True   — less than or equal to
```

A frequent stumble: `=` *assigns* a value, while `==` *compares* two values. Use `==` inside conditions.

To combine conditions, use the **logical operators** `and`, `or`, and `not`:

```python
age = 20
has_ticket = True

if age >= 18 and has_ticket:
    print("Welcome to the show!")

if not has_ticket:
    print("Please buy a ticket first.")
```

`and` is true only when *both* sides are true. `or` is true when *at least one* side is true. `not` flips a truth value. These let you express rich, real-world rules in a single readable line.

:::warning
`=` and `==` are not interchangeable. `if x = 5:` is a syntax error in Python (which is actually a kindness — it stops a common bug before it starts). Always use `==` to compare.
:::

## Truthiness

Here's a Python superpower: conditions don't strictly need a comparison. Many values are considered "truthy" or "falsy" on their own. The number `0`, an empty string `""`, and an empty list `[]` are all treated as **False**; most other values are treated as **True**.

```python
name = ""
if name:
    print(f"Hello, {name}")
else:
    print("No name was given.")   # this runs, because "" is falsy
```

This lets you write clean checks like "if there's something here, use it." You'll see `if some_list:` to mean "if the list isn't empty." It reads naturally once it clicks.

:::key
Empty or zero-like things are *falsy*: `0`, `0.0`, `""`, `[]`, `{}`, and `None`. Almost everything else is *truthy*. This is why `if my_list:` is the idiomatic way to ask "does my list have anything in it?"
:::

## Repeating with loops

Computers shine at doing the same thing many times without getting bored. That's what **loops** are for.

A `while` loop repeats *as long as* its condition stays true:

```python
count = 1
while count <= 3:
    print(f"Count is {count}")
    count = count + 1   # crucial: move toward making the condition false
```

That last line matters enormously. If the condition never becomes false, the loop runs forever — an *infinite loop*. Every `while` loop needs something inside it that eventually flips the condition.

A `for` loop repeats over a *sequence* of items. Paired with `range()`, it counts for you:

```python
for i in range(5):
    print(i)        # prints 0, 1, 2, 3, 4
```

Notice `range(5)` produces `0, 1, 2, 3, 4` — it starts at 0 and stops *before* 5, giving exactly five numbers. A `for` loop can also walk through the items of a list directly, which is often what you really want:

```python
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print(f"I like {fruit}.")
```

:::tip
Reach for `for` when you know what you're looping over (a fixed range, or the items in a list). Reach for `while` when you're repeating until *some condition* changes and you don't know in advance how many times that will take.
:::

## Steering loops with break and continue

Sometimes you need to bail out of a loop early, or skip just one round. Two keywords give you that control.

`break` exits the loop immediately, no matter where you are:

```python
for n in range(100):
    if n == 3:
        break
    print(n)        # prints 0, 1, 2 — then break stops the loop
```

`continue` skips the rest of the current iteration and jumps to the next one:

```python
for n in range(5):
    if n == 2:
        continue
    print(n)        # prints 0, 1, 3, 4 — 2 is skipped
```

Use them sparingly and clearly. `break` is perfect for "stop as soon as I find what I want." `continue` is handy for "ignore the cases I don't care about and keep going."

## Talk about it

> Control flow is just the logic of everyday decisions written down precisely. Think of a routine you follow — getting ready in the morning, deciding what to wear based on the weather, or repeating a chore until it's done. Describe it in plain words using "if", "otherwise", and "repeat until". You've just sketched a program. Which parts are `if`/`else`, and which parts are loops?

## What's next

Your programs can now choose paths and repeat work — that's a giant leap. But as logic grows, you'll want to *package* it into reusable, named chunks instead of repeating yourself. That's exactly what **Functions in Python** is about. Run a few loops in the [Code Lab](#code) first, then come find Pixel for the next step.

:::quiz
Q: How many numbers does `range(5)` produce, and what are they?
- Five numbers: 0, 1, 2, 3, 4 *
- Five numbers: 1, 2, 3, 4, 5
- Six numbers: 0 through 5
E: `range(5)` starts at 0 and stops before 5, yielding 0, 1, 2, 3, 4 — five numbers.
:::

:::predict
```python
for n in range(5):
    if n == 2:
        continue
    print(n)
```
- 0 1 3 4 *
- 0 1 2 3 4
- 0 1
E: `continue` skips the iteration when n is 2, so 2 is never printed; the loop continues with 3 and 4.
:::

:::fill
Q: Complete the keyword that adds a second condition between `if` and `else`.
`if x > 10:`
`    ...`
`___ x > 5:`
`    ...`
`else:`
`    ...`
- elif *
- elseif
- else if
E: Python uses `elif` (one word) for an additional condition after an `if`.
:::

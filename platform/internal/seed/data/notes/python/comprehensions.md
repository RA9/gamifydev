# Comprehensions

You already know how to build a new list with a `for` loop and `.append()`. It works, but it takes several lines to say something simple like "double every score." Python offers a shortcut that says it in *one* line: the **comprehension**. It's one of the most loved features of the language — compact, readable, and fast. In this lesson you'll learn list, dictionary, and set comprehensions, always shown right next to the plain loop they replace, so you can see exactly what's happening. You'll also learn the equally important skill of knowing when *not* to use one.

This is a lesson made for the [Code Lab](#code). Run both versions of each example side by side and confirm they produce the same result.

## From loop to comprehension

Say you have a list of scores and want a new list with each one doubled. Here's the familiar loop:

```python
scores = [10, 20, 30]

# The plain loop
doubled = []
for score in scores:
    doubled.append(score * 2)
print(doubled)      # [20, 40, 60]
```

A **list comprehension** does the same thing in one line:

```python
scores = [10, 20, 30]

# The comprehension
doubled = [score * 2 for score in scores]
print(doubled)      # [20, 40, 60]
```

Read it left to right: *"give me `score * 2` for each `score` in `scores`."* The three moving parts of the loop — create an empty list, loop, append — collapse into a single expression wrapped in square brackets.

:::analogy
A plain loop is a recipe written out step by step: "get a bowl, then for each apple, slice it, then put the slice in the bowl." A comprehension is the same recipe said as one sentence: "a bowl of sliced-apple for each apple." Same result, but the comprehension states the *outcome* directly instead of the procedure.
:::

## The shape of a comprehension

Every list comprehension follows the same skeleton:

```
[ expression   for item in collection ]
      ▲              ▲
   what to make   where the items come from
```

The `expression` is what you want in the new list — it can transform the item any way you like. The `for item in collection` part is exactly the loop header you already know.

```python
names = ["ada", "grace", "linus"]
loud = [name.upper() for name in names]
print(loud)         # ['ADA', 'GRACE', 'LINUS']

lengths = [len(name) for name in names]
print(lengths)      # [3, 5, 5]
```

:::key
A list comprehension is `[expression for item in collection]`. The expression builds each new value; the `for` clause supplies the items. It replaces the three-line "empty list + loop + append" pattern with one readable line — as long as that line stays simple.
:::

## Filtering with if

Add an `if` at the end and the comprehension keeps only the items that pass the test. Compare the loop and the comprehension:

```python
scores = [45, 90, 30, 88, 12]

# The plain loop
passing = []
for score in scores:
    if score >= 50:
        passing.append(score)
print(passing)      # [90, 88]

# The comprehension
passing = [score for score in scores if score >= 50]
print(passing)      # [90, 88]
```

Read it as: *"give me `score` for each `score` in `scores` **if** it's at least 50."* The trailing `if` decides which items make it in.

You can transform *and* filter at the same time — double only the even numbers:

```python
nums = [1, 2, 3, 4, 5, 6]
even_doubled = [n * 2 for n in nums if n % 2 == 0]
print(even_doubled)   # [4, 8, 12]
```

:::predict
```python
print([x for x in range(4) if x % 2 == 0])
```
- [0, 1, 2, 3]
- [0, 2] *
- [2, 4]
E: `range(4)` gives 0, 1, 2, 3. The `if x % 2 == 0` keeps only the even numbers, leaving [0, 2].
:::

## Dictionary comprehensions

The same idea builds dictionaries — just use curly braces and a `key: value` expression. Suppose you want to map each player's name to their score length or a computed value:

```python
names = ["ada", "grace", "linus"]

# The plain loop
name_lengths = {}
for name in names:
    name_lengths[name] = len(name)
print(name_lengths)   # {'ada': 3, 'grace': 5, 'linus': 5}

# The comprehension
name_lengths = {name: len(name) for name in names}
print(name_lengths)   # {'ada': 3, 'grace': 5, 'linus': 5}
```

The pattern is `{key: value for item in collection}`. You can filter these too. Here we keep only the high scorers from an existing dictionary:

```python
scores = {"Ada": 90, "Grace": 40, "Linus": 75}
winners = {name: pts for name, pts in scores.items() if pts >= 70}
print(winners)        # {'Ada': 90, 'Linus': 75}
```

Notice we loop over `.items()` to get both the key and value, exactly like a normal dictionary loop — just wrapped in the comprehension form.

## Set comprehensions

Curly braces with a *single* expression (no colon) build a **set** — unique values only, no duplicates. It's handy for pulling the distinct results out of some data in one step:

```python
rolls = [1, 3, 3, 5, 1, 6, 5]
distinct = {r for r in rolls}
print(distinct)       # {1, 3, 5, 6}  — duplicates dropped

# transform while deduping: which even numbers appeared?
evens_seen = {r for r in rolls if r % 2 == 0}
print(evens_seen)     # {6}
```

Same skeleton as a list comprehension, just with `{ }` — and because it's a set, any repeats collapse automatically.

## When a plain loop is clearer

Comprehensions are wonderful, but they are not always the right choice. Their whole value is *readability at a glance* — and once a comprehension gets long, complicated, or does more than "build a collection," that value disappears. Then a plain loop is the better, kinder code.

Don't cram heavy logic into a comprehension:

```python
# Hard to read — too much going on in one line
result = [x * 2 if x > 0 else x * -1 for x in nums if x != 0 and x < 100]
```

The loop version is longer but far easier to follow and debug:

```python
# Clearer — each step is on its own line
result = []
for x in nums:
    if x == 0 or x >= 100:
        continue
    if x > 0:
        result.append(x * 2)
    else:
        result.append(x * -1)
```

Also avoid a comprehension when the loop's real job is a **side effect** — like printing — rather than *building a collection*. A comprehension should produce a value; if you're only doing something each pass, use a normal `for` loop.

```python
# Don't do this — building a list of Nones just to print
[print(name) for name in names]   # works, but wasteful and confusing

# Just loop
for name in names:
    print(name)
```

:::warning
The moment you find yourself squinting at a comprehension to understand it, rewrite it as a loop. Clever one-liners feel good to write and hurt to read later. A comprehension earns its place only when it stays short and obvious — otherwise, reach back for the plain loop you already know.
:::

:::tip
A good rule of thumb: if the comprehension fits comfortably on one line and does a single clear thing — transform, filter, or both — use it. If it needs nested logic, an `else`, multiple conditions, or a side effect, use a loop. Readability wins over brevity every time.
:::

:::predict
```python
words = ["hi", "hey", "yo", "hello"]
print([w for w in words if len(w) <= 2])
```
- ['hi', 'hey', 'yo', 'hello']
- ['hi', 'yo'] *
- ['hey', 'hello']
E: The filter keeps only words with length 2 or less. "hi" and "yo" qualify; "hey" (3) and "hello" (5) are excluded.
:::

## What's next

You've learned to compress a whole build-a-collection loop into a single expressive line — list comprehensions for the everyday, dictionary and set comprehensions for labeled and unique data, and the judgment to know when a plain loop reads better. That judgment matters as much as the syntax: the goal is always clear code, not the shortest code. Try rewriting a few of your earlier loops as comprehensions in the [Code Lab](#code) — and rewrite one back into a loop when it gets too dense. Pixel's proud of you.

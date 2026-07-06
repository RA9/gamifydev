# Functional Tools

Python borrows a handful of ideas from functional programming that let you transform, filter, and combine collections without hand-writing loops. Used well, they make code shorter *and* clearer. Used carelessly, they make it cryptic. This lesson is about knowing which tool fits — and when a plain comprehension beats them all.

## map and filter (and Why Comprehensions Often Win)

`map(func, iterable)` applies a function to every item. `filter(func, iterable)` keeps only the items where the function returns truthy. Both return lazy iterators, so you usually wrap them in `list()` to see the result.

```python
scores = [10, 25, 40]

doubled = list(map(lambda s: s * 2, scores))
print(doubled)   # → [20, 50, 80]

high = list(filter(lambda s: s >= 25, scores))
print(high)      # → [25, 40]
```

Those work, but Python has a more idiomatic tool for exactly this shape of task: the **list comprehension**. Compare:

```python
# map + lambda
doubled = list(map(lambda s: s * 2, scores))

# comprehension — the Pythonic version
doubled = [s * 2 for s in scores]

# filter + lambda
high = list(filter(lambda s: s >= 25, scores))

# comprehension with a condition
high = [s for s in scores if s >= 25]
```

The comprehensions read left to right in plain English and don't need `list()` or a `lambda`. This is the community's default choice.

:::key
Reach for a comprehension first. `map`/`filter` earn their place mainly when you already have a **named function** to pass — then `map(str.upper, names)` is cleaner than `[n.upper() for n in names]`. But whenever you'd write a `lambda` inside `map`/`filter`, a comprehension is almost always clearer.
:::

```python
names = ["ada", "grace", "linus"]

# Passing an existing named function — map reads nicely here
shout = list(map(str.upper, names))
print(shout)   # → ['ADA', 'GRACE', 'LINUS']
```

## Sorting and Extremes with key

`sorted()`, `min()`, and `max()` all accept a `key=` argument: a function that turns each item into the value you want to compare by. This is where lambdas genuinely shine.

Say a leaderboard is a list of `(name, score)` tuples. Sort it by score, highest first:

```python
leaderboard = [("Ada", 90), ("Grace", 120), ("Linus", 75)]

by_score = sorted(leaderboard, key=lambda row: row[1], reverse=True)
print(by_score)
# → [('Grace', 120), ('Ada', 90), ('Linus', 75)]
```

The `key=lambda row: row[1]` says "compare rows by their score (index 1)." `reverse=True` flips it to descending — the natural order for a leaderboard.

`min()` and `max()` take the same `key`:

```python
top = max(leaderboard, key=lambda row: row[1])
print(top)   # → ('Grace', 120)

worst = min(leaderboard, key=lambda row: row[1])
print(worst) # → ('Linus', 75)
```

You can sort by anything you can compute. Sort names by length, then break ties alphabetically by returning a tuple:

```python
names = ["Ada", "Bo", "Grace", "Al"]
ordered = sorted(names, key=lambda n: (len(n), n))
print(ordered)   # → ['Bo', 'Al', 'Ada', 'Grace']
```

The key returns `(length, name)`, and tuples compare element by element: shorter names first, and equal-length names fall back to alphabetical order.

:::tip
When sorting objects or dataclasses, `operator.itemgetter` and `operator.attrgetter` are faster, named alternatives to a lambda: `sorted(players, key=attrgetter("score"))` reads more clearly than `key=lambda p: p.score` and skips the lambda entirely.
:::

:::predict
Q: What does this print?
```python
words = ["kiwi", "fig", "apple"]
print(sorted(words, key=len))
```
- ['apple', 'fig', 'kiwi']
- ['fig', 'kiwi', 'apple']*
- ['kiwi', 'fig', 'apple']
- ['apple', 'kiwi', 'fig']
E: `key=len` sorts by string length, ascending: "fig" (3), "kiwi" (4), "apple" (5). The original order and alphabetical order don't apply.
:::

## Combining Everything with reduce

`functools.reduce(func, iterable)` boils a whole sequence down to a single value by applying a two-argument function cumulatively. It carries a running result and folds in each item.

```python
from functools import reduce

nums = [1, 2, 3, 4]
product = reduce(lambda acc, n: acc * n, nums)
print(product)   # → 24   (((1*2)*3)*4)
```

`acc` is the accumulated result so far; `n` is the next item. You can pass a starting value as a third argument:

```python
total = reduce(lambda acc, n: acc + n, nums, 100)
print(total)   # → 110   (starts at 100)
```

Honesty matters here: for common jobs, Python already has a *better* built-in. Summing? Use `sum(nums)`. Largest? Use `max(nums)`. Those are clearer than `reduce`.

:::warning
`reduce` is easy to overuse. If a built-in (`sum`, `max`, `min`, `any`, `all`, `"".join`) or a simple loop expresses the intent more plainly, prefer it. Save `reduce` for genuinely custom folds that have no built-in equivalent — like multiplying a list, or merging a sequence of dicts.
:::

## Memoization with lru_cache

Some functions are expensive and get called repeatedly with the same arguments. `functools.lru_cache` remembers past results so a repeat call is instant — a technique called **memoization**. Just decorate the function.

```python
from functools import lru_cache

@lru_cache(maxsize=None)
def fib(n: int) -> int:
    if n < 2:
        return n
    return fib(n - 1) + fib(n - 2)

print(fib(35))   # instant, even though naive recursion would be brutally slow
```

Without the cache, `fib(35)` recomputes the same subproblems millions of times. With it, each `fib(n)` is computed once and stored; later calls with the same `n` return the saved answer immediately.

:::analogy
`lru_cache` is a receptionist with a notepad. The first time someone asks a hard question, she works out the answer and jots it down. Next time anyone asks the *same* question, she just reads it off the pad instead of solving it again.
:::

Imagine an expensive score calculation you call inside a leaderboard render loop:

```python
from functools import lru_cache

@lru_cache(maxsize=128)
def rank_bonus(score: int) -> int:
    # pretend this is a slow computation
    total = 0
    for i in range(score):
        total += i % 7
    return total

# Called many times with repeating scores — the cache pays off fast
print(rank_bonus(90))   # computed once
print(rank_bonus(90))   # served from cache, instantly
```

Two rules keep `lru_cache` safe:

- **Arguments must be hashable** — ints, strings, tuples work; lists and dicts don't.
- **The function must be pure** — same inputs always give the same output, with no side effects. Caching a function that reads the current time or a database would hand back stale answers.

`maxsize` caps how many results are stored, evicting the least-recently-used ones when full; `maxsize=None` means unlimited.

## A Tour of itertools

The `itertools` module is a toolbox of lazy iterator builders. They're memory-efficient (they produce items on demand) and often replace fiddly manual loops. A few you'll actually reach for:

**`chain`** — stitches multiple iterables into one stream, no copying:

```python
from itertools import chain

melee = ["sword", "axe"]
ranged = ["bow", "sling"]
for weapon in chain(melee, ranged):
    print(weapon)   # → sword, axe, bow, sling
```

**`count`** — an endless counter, handy for generating IDs:

```python
from itertools import count

ids = count(start=1)
print(next(ids), next(ids), next(ids))   # → 1 2 3
```

**`cycle`** — loops over an iterable forever, e.g. rotating turns:

```python
from itertools import cycle, islice

turns = cycle(["Ada", "Grace"])
print(list(islice(turns, 4)))   # → ['Ada', 'Grace', 'Ada', 'Grace']
```

`count` and `cycle` are *infinite*, so pair them with `islice`, `zip`, or a `break` — never a bare `for` that tries to consume them all.

**`combinations`** — every way to pick `r` items, order not mattering. Perfect for pairing up players in a tournament:

```python
from itertools import combinations

players = ["Ada", "Grace", "Linus"]
for pair in combinations(players, 2):
    print(pair)
# → ('Ada', 'Grace'), ('Ada', 'Linus'), ('Grace', 'Linus')
```

**`groupby`** — groups *consecutive* items sharing a key. Its one gotcha: it only groups runs that are already adjacent, so you almost always **sort first** by the same key:

```python
from itertools import groupby

players = [
    ("Ada", "fire"), ("Bo", "ice"),
    ("Cy", "fire"), ("Di", "ice"),
]

players.sort(key=lambda p: p[1])   # sort by element first!

for element, group in groupby(players, key=lambda p: p[1]):
    names = [name for name, _ in group]
    print(element, names)
# → fire ['Ada', 'Cy']
# → ice ['Bo', 'Di']
```

:::warning
`groupby` does **not** collect all matching items across the whole iterable — only consecutive runs. Forgetting to sort first is the classic bug: unsorted input produces several tiny groups for the same key instead of one big one.
:::

:::predict
Q: What does this print?
```python
from itertools import combinations
print(len(list(combinations([1, 2, 3, 4], 2))))
```
- 4
- 6*
- 8
- 16
E: `combinations` gives every unordered pair from 4 items: (1,2), (1,3), (1,4), (2,3), (2,4), (3,4) — 6 pairs.
:::

## Recap

- `map`/`filter` transform and select items, but a **comprehension** is usually clearer — reserve `map`/`filter` for passing an existing named function.
- `sorted`, `min`, and `max` take `key=` to compare by a computed value; lambdas (or `itemgetter`/`attrgetter`) let you sort a leaderboard by score, and `reverse=True` gives descending order.
- Returning a tuple from `key` sorts by multiple criteria with tie-breaking.
- `functools.reduce` folds a sequence to one value, but prefer built-ins like `sum`/`max` when they fit.
- `functools.lru_cache` memoizes pure, hashable-argument functions — turning repeated expensive calls into instant lookups.
- `itertools` gives lazy building blocks: `chain` (join streams), `count`/`cycle` (infinite — slice them!), `combinations` (unordered picks), and `groupby` (group consecutive runs — sort first!).
- The through-line: pick the tool that makes intent obvious. Readable beats clever.

**Next up:** Generators and Iterators — building your own lazy sequences with `yield` and understanding what powers `itertools` under the hood.

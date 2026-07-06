# Iterators and Generators

You already write `for` loops over lists, strings, and dicts every day. But *how* does a `for` loop actually walk through those items? And what do you do when the "list" is too big to fit in memory — or genuinely endless, like a stream of enemy waves that never stops? This lesson pulls back the curtain on the iterator protocol and then hands you generators: the cleanest, laziest, most memory-friendly way to produce values on demand.

## The Iterator Protocol

Every time you write `for item in something`, Python quietly does two things. First it calls `iter(something)` to get an **iterator**. Then it calls `next(...)` on that iterator over and over, pulling one value each time, until the iterator signals it's out of values.

You can drive this by hand:

```python
waves = ["goblins", "orcs", "dragon"]
it = iter(waves)

print(next(it))  # → goblins
print(next(it))  # → orcs
print(next(it))  # → dragon
print(next(it))  # → raises StopIteration
```

That final `next(it)` raises `StopIteration`. This is not a bug — it's the agreed-upon signal that means "nothing left." A `for` loop catches that exception silently and simply stops looping. You never see it because the loop handles it for you.

:::key
An **iterable** is anything you can call `iter()` on (lists, strings, dicts, sets, files). An **iterator** is what `iter()` gives you back — an object with a `__next__` method that produces the next value or raises `StopIteration`. Every iterator is also iterable; not every iterable is an iterator.
:::

So a `for` loop is really just this, spelled out:

```python
it = iter(waves)
while True:
    try:
        item = next(it)
    except StopIteration:
        break
    print(item)
```

Understanding this matters because it explains a common surprise: an iterator is **single-use**. Once you've exhausted it, it stays exhausted.

```python
it = iter([1, 2, 3])
print(list(it))  # → [1, 2, 3]
print(list(it))  # → []  (already drained!)
```

:::warning
A list can be looped over again and again — but the *iterator* from `iter(list)` is one-shot. If `list(gen)` gives you data the first time and an empty list the second, you've hit an exhausted iterator. Generators (coming up) have exactly this behavior.
:::

## What Makes Something Iterable

An object is iterable if it defines a `__iter__` method that returns an iterator. You rarely write this by hand because generators do it for you — but seeing it once demystifies the whole thing. Here's a tiny iterable that yields three spawn points:

```python
class SpawnPoints:
    def __init__(self, count):
        self.count = count

    def __iter__(self):
        n = 0
        while n < self.count:
            yield f"spawn_{n}"
            n += 1

for point in SpawnPoints(3):
    print(point)  # → spawn_0, spawn_1, spawn_2
```

Notice `__iter__` used `yield`. That single keyword turns a method into a generator, which *is* an iterator — so we get the whole protocol for free. Let's look at `yield` properly.

## Generator Functions with yield

A **generator function** looks like a normal function, except it uses `yield` instead of (or alongside) `return`. The moment Python sees `yield` anywhere in a function body, calling that function no longer runs the body — it hands you back a generator object, paused at the very top.

```python
def enemy_waves():
    print("wave 1 forming...")
    yield "goblins"
    print("wave 2 forming...")
    yield "orcs"
    print("boss incoming...")
    yield "dragon"

gen = enemy_waves()      # nothing printed yet — body hasn't run
print(next(gen))         # → wave 1 forming...  then  goblins
print(next(gen))         # → wave 2 forming...  then  orcs
```

Each `next()` runs the function *up to and including* the next `yield`, then freezes it — local variables, position, everything. The next `next()` thaws it and continues from exactly where it stopped. When the function falls off the end (or hits `return`), it raises `StopIteration`.

:::analogy
A normal function is a vending machine: you press the button, it drops all its output at once, done. A generator is a bartender pouring drinks one at a time — it pours a `yield`, then leans on the counter waiting. It remembers exactly where it was in the recipe and picks up there when you ask for the next round.
:::

Because generators are iterators, you almost never call `next()` yourself. You loop:

```python
for wave in enemy_waves():
    print(f"Spawning: {wave}")
```

:::predict
Q: What does this print?
```python
def counter():
    yield 1
    yield 2
    yield 3

print(list(counter()))
```
- [1, 2, 3]*
- [1, 1, 1]
- <generator object>
- (1, 2, 3)
E: `list()` drives the generator to exhaustion, collecting each yielded value in order → [1, 2, 3].
:::

## Why Generators Are Lazy and Memory-Efficient

The killer feature is **laziness**: a generator computes each value only when asked, and forgets it right after. It never holds the whole sequence in memory at once.

Compare building a list of a million level IDs versus generating them:

```python
# Eager: builds and stores ALL million strings in RAM at once
def all_levels_list(n):
    result = []
    for i in range(n):
        result.append(f"level_{i:06d}")
    return result

# Lazy: produces one string at a time, storing nothing extra
def all_levels_gen(n):
    for i in range(n):
        yield f"level_{i:06d}"
```

The list version allocates a million strings up front — real memory, real delay. The generator version holds essentially nothing; each `level_...` string exists only for the instant you use it, then is discarded. If you only need the first ten, the generator produces exactly ten and never touches the other 999,990.

:::tip
Reach for a generator whenever you're producing a sequence you'll consume once, top to bottom — especially if it's large, expensive to compute, or you might stop early. If you need random access (`data[500]`) or to loop over it several times, build a list instead.
:::

You can confirm the memory difference isn't hand-waving:

```python
import sys

big_list = [i for i in range(100_000)]
big_gen  = (i for i in range(100_000))

print(sys.getsizeof(big_list))  # → hundreds of thousands of bytes
print(sys.getsizeof(big_gen))   # → ~200 bytes, no matter the range
```

That second line uses a **generator expression** — our next topic.

## Generator Expressions

You already know list comprehensions: `[x*x for x in range(3)]`. Swap the square brackets for parentheses and you get a **generator expression** — same syntax, but lazy. It produces a generator instead of a list.

```python
squares_list = [x*x for x in range(5)]   # → [0, 1, 4, 9, 16]  (a real list)
squares_gen  = (x*x for x in range(5))   # → a generator, nothing computed yet

print(squares_list)          # → [0, 1, 4, 9, 16]
print(squares_gen)           # → <generator object ...>
print(list(squares_gen))     # → [0, 1, 4, 9, 16]  (now it runs)
```

They shine when feeding a function that consumes items one by one. For example, summing enemy health without ever building the list:

```python
enemies = ["goblin", "orc", "troll", "dragon"]

# The parentheses of sum() double as the generator's parentheses
total_letters = sum(len(name) for name in enemies)
print(total_letters)  # → 20
```

No intermediate list is created — `sum` pulls one length at a time. This works with `max`, `min`, `any`, `all`, `"".join(...)`, and more.

:::warning
Because a generator is single-use, listing it twice bites you. `g = (x for x in range(3)); list(g)` gives `[0, 1, 2]`, but a second `list(g)` gives `[]`. If you need the values more than once, materialize them into a list.
:::

## Infinite Generators and Consuming Them Safely

Since a generator only produces values on demand, it can describe an **endless** sequence — something impossible with a list. A `while True` with a `yield` inside will keep handing out values forever:

```python
def level_ids():
    n = 1
    while True:
        yield f"level_{n:04d}"
        n += 1
```

This is perfectly safe *as long as you never try to consume all of it*. Do NOT call `list(level_ids())` — it would loop forever and hang your program. Instead, you must always give yourself a way out. The simplest is `break`:

```python
for level in level_ids():
    print(level)
    if level == "level_0005":
        break
# → level_0001 ... level_0005, then stops
```

Another clean pattern is pairing an infinite generator with a bounded one using `zip`. `zip` stops as soon as its shortest input runs out, so a finite `range` reins in the infinite side:

```python
for i, level in zip(range(3), level_ids()):
    print(i, level)
# → 0 level_0001, 1 level_0002, 2 level_0003
```

:::analogy
An infinite generator is a spring of water — it flows forever, but you only ever fill the cup you're holding. The danger isn't the spring; it's trying to pour the whole spring into a bucket. Always bring a cup with a known size.
:::

## itertools.islice — Taking the First N

Slicing a list with `waves[:5]` is easy, but you can't slice a generator that way — generators don't support `[:5]`. The standard library's `itertools.islice` is the generator-friendly equivalent: it lazily takes a slice out of *any* iterable, including infinite ones.

```python
from itertools import islice

# Grab the first 5 level IDs from an endless generator — safely
first_five = list(islice(level_ids(), 5))
print(first_five)
# → ['level_0001', 'level_0002', 'level_0003', 'level_0004', 'level_0005']
```

`islice(iterable, stop)` takes the first `stop` items. It also accepts start and step like a normal slice: `islice(gen, 2, 10, 2)` skips the first two, then takes every other item up to the tenth. Crucially, it pulls values lazily — `islice(level_ids(), 5)` asks the infinite generator for exactly five values and no more, so it never hangs.

```python
from itertools import islice

def enemy_stream():
    kinds = ["goblin", "orc", "troll"]
    i = 0
    while True:
        yield kinds[i % len(kinds)]
        i += 1

# Take just the first 7 spawns from an endless stream
wave = list(islice(enemy_stream(), 7))
print(wave)
# → ['goblin', 'orc', 'troll', 'goblin', 'orc', 'troll', 'goblin']
```

:::key
`itertools.islice(iterable, n)` is the safe, lazy way to say "give me the first N" from any iterable — the only sane way to pull a finite chunk out of an infinite generator. Wrap it in `list(...)` when you want a concrete list back.
:::

:::predict
Q: What does this print?
```python
from itertools import islice

def naturals():
    n = 0
    while True:
        n += 1
        yield n

print(list(islice(naturals(), 4)))
```
- [1, 2, 3, 4]*
- [0, 1, 2, 3]
- hangs forever
- [4]
E: `naturals` increments before yielding, so it produces 1, 2, 3, 4... and `islice(..., 4)` lazily takes exactly the first four → [1, 2, 3, 4]. It never runs forever because islice stops after 4.
:::

## Recap

- A `for` loop is sugar for `iter()` + repeated `next()`, stopping when `StopIteration` is raised.
- An **iterable** can be handed to `iter()`; an **iterator** is what you get back — one-shot, exhausted after a full pass.
- A **generator function** uses `yield`; calling it returns a paused generator that resumes at each `yield` and remembers its state.
- Generators are **lazy** and **memory-efficient**: they compute one value at a time and store nothing extra — ideal for large or run-once sequences.
- **Generator expressions** `(x for x in ...)` are lazy list comprehensions; feed them straight into `sum`, `max`, `any`, `join`, and friends.
- Generators can be **infinite** (`while True: yield`). Never fully consume them — always `break`, `zip` with something finite, or slice.
- **`itertools.islice(iterable, n)`** safely takes the first N items from any iterable, including infinite generators.

**Next up:** Decorators — wrapping functions to add behavior like logging and validation without touching their code.

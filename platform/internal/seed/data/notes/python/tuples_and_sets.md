# Tuples and Sets

You've met lists and dictionaries — the two workhorses of Python collections. Now meet their sharper, more specialized cousins: the **tuple** and the **set**. Each solves a specific problem that lists and dictionaries handle clumsily. A tuple is a *fixed* group of values that should never change — like a coordinate. A set is a collection where *duplicates are impossible* and membership checks are lightning fast — like tracking which achievements a player has earned. Learn when to reach for each and your code gets clearer, safer, and faster.

This is a lesson made for the [Code Lab](#code). Type the examples, tweak them, and watch what happens.

## Tuples: fixed, ordered groups

A **tuple** is an ordered sequence of values, just like a list — but you write it with **parentheses** instead of square brackets, and once you make it, you *cannot change it*.

```python
point = (3, 5)
print(point)        # (3, 5)
print(point[0])     # 3  — index just like a list
print(point[1])     # 5
print(len(point))   # 2
```

Indexing and slicing work exactly as they do for lists — position `0` is first, `[-1]` is last. The one big difference: a tuple is **immutable**. Try to change an item and Python stops you cold.

```python
point = (3, 5)
point[0] = 9   # TypeError: 'tuple' object does not support item assignment
```

:::analogy
A list is like a whiteboard — you can erase and rewrite any square whenever you like. A tuple is like a stamp pressed into concrete: the values are set the moment you make it, and they stay that way forever. That permanence is a *feature*, not a limitation — some data is never meant to change.
:::

## When to use a tuple

Reach for a tuple whenever a group of values belongs together and shouldn't be edited afterward. Coordinates are the classic example: an `(x, y)` position on a game map is a single thing made of two numbers, and it would be a bug to accidentally change just one half of it.

```python
spawn_point = (0, 0)
boss_location = (128, 64)
rgb_red = (255, 0, 0)      # a color is three fixed channels
```

Because tuples can't change, they're also safe to use as **dictionary keys** (lists can't be keys, but tuples can). This lets you map a grid position straight to what's on it:

```python
grid = {}
grid[(2, 3)] = "treasure"
grid[(5, 1)] = "trap"
print(grid[(2, 3)])   # treasure
```

## Unpacking a tuple

Here's where tuples really shine. You can **unpack** a tuple into separate variables in one line — Python hands out the values in order:

```python
point = (3, 5)
x, y = point       # x gets 3, y gets 5
print(x)           # 3
print(y)           # 5
```

No brackets needed on the right side either — Python packs and unpacks automatically. This makes swapping two values a one-liner with no temporary variable:

```python
a, b = 10, 20
a, b = b, a        # swap!
print(a, b)        # 20 10
```

:::key
Tuples are **ordered** (like lists) and **immutable** (unlike lists). Write them with `(parentheses)`, index them with `[0]`, and unpack them with `x, y = point`. Use one whenever a fixed group of values — a coordinate, a color, a pair — travels together as a single unit.
:::

## Returning multiple values from a function

This is the everyday reason tuples matter. A Python function can only "return" one thing — but if that one thing is a tuple, you've effectively returned several values at once. The caller unpacks them.

```python
def get_player_stats():
    health = 100
    mana = 50
    return health, mana        # this is a tuple: (100, 50)

hp, mp = get_player_stats()    # unpack the returned tuple
print(f"HP: {hp}, MP: {mp}")   # HP: 100, MP: 50
```

You'll see this pattern constantly. A function that finds an enemy might return `(x, y)`; one that parses a name might return `(first, last)`. Returning a tuple keeps related results bundled together.

:::predict
```python
def divide(a, b):
    return a // b, a % b

q, r = divide(17, 5)
print(q, r)
```
- 3 2 *
- 2 3
- 3.4 2
E: The function returns a tuple `(3, 2)` — `17 // 5` is the whole-number quotient `3`, and `17 % 5` is the remainder `2`. Unpacking assigns `q = 3` and `r = 2`.
:::

## Sets: unique values only

A **set** is an unordered collection where every value is **unique** — duplicates simply can't exist. You write one with curly braces, like a dictionary, but with single values instead of `key: value` pairs.

```python
achievements = {"first_kill", "level_10", "speedrun"}
print(achievements)   # order is not guaranteed!
print(len(achievements))   # 3
```

Two things to notice. First, a set has **no order** — you can't index it (`achievements[0]` is an error), and Python may print the items in any order. Second, if you *try* to include a duplicate, the set silently ignores it:

```python
badges = {"gold", "gold", "silver", "gold"}
print(badges)         # {'gold', 'silver'}  — only one 'gold'
print(len(badges))    # 2
```

:::warning
Curly braces with values make a **set**, but empty curly braces `{}` make an empty **dictionary**, not a set — that's a classic trap. To make an empty set you must write `set()`:

```python
earned = set()        # empty set — correct
oops = {}             # this is an empty DICT, not a set
```
:::

## Adding, removing, and fast membership

You grow a set with `.add()` and shrink it with `.remove()` or `.discard()`. Adding a value that's already there does nothing — which is exactly what you want for "earned achievements."

```python
earned = set()
earned.add("first_kill")
earned.add("first_kill")   # ignored — already there
earned.add("boss_slain")
print(earned)              # {'first_kill', 'boss_slain'}

earned.discard("first_kill")   # remove; no error if missing
```

The superpower of sets is checking membership with `in`. For a long list, `in` has to scan item by item; for a set, the check is nearly instant no matter how big the set gets.

```python
if "boss_slain" in earned:
    print("Achievement unlocked!")
```

:::tip
If your code asks "have I seen this before?" or "does this collection contain X?" a *lot*, store your data in a set. Membership testing with `in` is dramatically faster on a set than on a list, because a set is built for exactly that question.
:::

## Combining sets: union, intersection, difference

Sets come with clean operators borrowed from math. Say two players have each earned some achievements:

```python
ada = {"first_kill", "speedrun", "level_10"}
grace = {"first_kill", "level_10", "no_damage"}

print(ada | grace)   # union: everything either earned
# {'first_kill', 'speedrun', 'level_10', 'no_damage'}

print(ada & grace)   # intersection: earned by BOTH
# {'first_kill', 'level_10'}

print(ada - grace)   # difference: Ada has, Grace doesn't
# {'speedrun'}
```

- `|` **union** — items in *either* set (merge them).
- `&` **intersection** — items in *both* sets (the overlap).
- `-` **difference** — items in the first but not the second.

These three operators replace loops full of `if` checks. "Which achievements do both players share?" is just `ada & grace`.

:::predict
```python
a = {1, 2, 3, 4}
b = {3, 4, 5, 6}
print(a & b)
```
- {1, 2, 3, 4, 5, 6}
- {3, 4} *
- {1, 2, 5, 6}
E: `&` is intersection — only the values present in *both* sets. `3` and `4` appear in each, so the result is `{3, 4}`.
:::

## Deduping a list with set()

One of the most common real uses of a set is stripping duplicates out of a list. Wrap the list in `set()` and every repeat vanishes; wrap that back in `list()` if you need a list again.

```python
scores = [10, 20, 10, 30, 20, 10]
unique = set(scores)
print(unique)          # {10, 20, 30}

unique_list = list(set(scores))
print(len(unique_list))   # 3
```

Watch out for one thing: because sets have no order, `list(set(...))` may not preserve the original ordering. If order matters, that's a sign the job wants a list, not a set. But when you just need "the distinct values," `set()` is the shortest tool in the box.

## How they all compare

You now have four collection types. Here's the quick mental map:

- **List** `[ ]` — ordered, changeable, allows duplicates. Your default for "a sequence of things."
- **Tuple** `( )` — ordered, *fixed*, allows duplicates. For a group of values that shouldn't change: coordinates, pairs, function returns.
- **Set** `{ }` (values) — *unordered*, changeable, *no duplicates*. For uniqueness and fast membership tests.
- **Dictionary** `{ }` (key: value) — labeled lookups by key.

The question that picks the winner: *Does order matter? Should it change? Do duplicates make sense?* Answer those three and the right structure is usually obvious.

:::quiz
Q: You want to track the distinct enemy types a player has defeated, and quickly check whether they've beaten a "dragon." Which structure fits best?
- a list
- a tuple
- a set *
E: You need uniqueness (each type counts once) and fast `in` checks ("have they beaten a dragon?"). That's exactly what a set is for.
:::

## What's next

You've rounded out Python's core collections: lists and dictionaries for the everyday, tuples for fixed groups, and sets for uniqueness and speed. Next up, **Loops: for and while**, where you'll learn to march through any of these collections item by item — spawning enemies, counting down timers, and totalling scores. Run these set operations in the [Code Lab](#code) first, then go loop with Pixel.

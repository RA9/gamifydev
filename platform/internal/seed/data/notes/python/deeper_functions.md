# Deeper Functions

You already know how to define a function, pass it parameters, and hand back a result with `return`. That's the foundation. Now we level up. Real game code needs functions that are *flexible* — ones that adapt to however the caller wants to use them, accept any number of inputs, and stay predictable under pressure. This lesson is about those power-ups.

As always, the [Code Lab](#code) is your workbench. Run every example and poke at it.

## Default parameter values

A **default value** lets a parameter fall back to something sensible when the caller doesn't supply it. You saw this briefly before; here it is doing real work.

```python
def spawn_enemy(name, hp=100, level=1):
    return f"{name} (HP {hp}, Lvl {level})"

print(spawn_enemy("Goblin"))              # Goblin (HP 100, Lvl 1)
print(spawn_enemy("Dragon", 500, 20))     # Dragon (HP 500, Lvl 20)
```

The rule that matters: parameters *with* defaults must come **after** parameters without them. Python needs to know which required values you're passing before it starts filling in optional ones.

```python
def spawn_enemy(hp=100, name):   # SyntaxError!
    ...
```

Defaults turn one function into a whole family of behaviours. Most enemies are level 1 with 100 HP, so callers rarely spell it out — but the option is right there when they need a boss.

## Keyword arguments

When you call a function, you can pass values by **position** (first value goes to first parameter) or by **keyword** (naming the parameter explicitly). Keywords make a call self-documenting and let you skip over defaults you don't care about.

```python
def spawn_enemy(name, hp=100, level=1, boss=False):
    tag = "BOSS " if boss else ""
    return f"{tag}{name} (HP {hp}, Lvl {level})"

# Positional — you have to count arguments in your head:
print(spawn_enemy("Slime", 30, 2, False))

# Keyword — reads like a sentence, and you can skip level:
print(spawn_enemy("Lich", hp=800, boss=True))
```

That second call sets `hp` and `boss` but lets `level` keep its default. You couldn't do that positionally — you'd have to pass a value for `level` just to reach `boss`.

:::tip
Reach for keyword arguments whenever a call has more than two or three values, or when you're passing a bare `True`/`False`. `spawn_enemy("Lich", boss=True)` tells the reader *what* that `True` means; `spawn_enemy("Lich", True)` makes them go hunt for the definition.
:::

## Accepting any number of positional arguments: `*args`

Sometimes you don't know how many values the caller will send. A `add_to_party` function might take one hero or five. The `*args` syntax collects all the extra positional arguments into a **tuple**.

```python
def total_score(*scores):
    return sum(scores)

print(total_score(10, 20))          # 30
print(total_score(5, 5, 5, 5, 5))   # 25
print(total_score())                # 0
```

Inside the function, `scores` is an ordinary tuple — you can loop over it, measure its length, index into it. The `*` is what does the collecting; the name `args` is just convention (you could write `*scores` as above, and it's clearer when you do).

```python
def announce(leader, *party):
    print(f"Leader: {leader}")
    for member in party:
        print(f"  - {member}")

announce("Ada", "Grace", "Linus", "Katherine")
```

Here `leader` grabs the first argument and `*party` sweeps up the rest. You get a required first value plus an open-ended list of the others.

## Accepting any number of keyword arguments: `**kwargs`

The keyword counterpart is `**kwargs`. It gathers any extra *keyword* arguments into a **dictionary**, mapping each name to its value.

```python
def build_character(name, **stats):
    print(f"Character: {name}")
    for stat, value in stats.items():
        print(f"  {stat}: {value}")

build_character("Warrior", strength=18, agility=12, luck=7)
```

That prints each stat and its value. Because `stats` is a plain dict, the caller can pass whatever attributes make sense — you didn't have to predict them all in advance. Together, `*args` and `**kwargs` let a function accept *anything*:

```python
def log_event(*args, **kwargs):
    print("positional:", args)
    print("keyword:", kwargs)

log_event("hit", 42, target="Goblin", crit=True)
# positional: ('hit', 42)
# keyword: {'target': 'Goblin', 'crit': True}
```

:::key
`*args` becomes a **tuple** of leftover positional values; `**kwargs` becomes a **dict** of leftover keyword values. One star for positions, two stars for named keywords. Order in the signature is always: normal params, then `*args`, then `**kwargs`.
:::

## The mutable-default trap

Here is a bug that catches nearly everyone once. It looks harmless: give a parameter an empty list as its default.

```python
def add_loot(item, bag=[]):     # DANGER
    bag.append(item)
    return bag

print(add_loot("sword"))    # ['sword']
print(add_loot("shield"))   # ['sword', 'shield']  ← the sword is still there!
```

The second call didn't get a fresh empty bag — it got the *same* bag, with the sword still in it. Why? A default value is evaluated **once**, when the function is defined, not each time it's called. That single list is created at definition and then shared by every call that relies on the default. Each `append` piles onto the same object.

The fix is a fixed habit: default to `None`, then build a fresh object inside.

```python
def add_loot(item, bag=None):
    if bag is None:
        bag = []
    bag.append(item)
    return bag

print(add_loot("sword"))    # ['sword']
print(add_loot("shield"))   # ['shield']  ← fresh bag each time
```

:::warning
Never use a mutable object — a list `[]`, a dict `{}`, or a set — as a default argument. It's created once and quietly shared across every call, producing bugs that seem impossible. Default to `None` and create the real object in the body.
:::

## Local vs. global scope, and the `global` keyword

A variable created inside a function is **local** — it lives only while the function runs. Variables at the top level of your file are **global**. A function can *read* a global, but assigning to a name inside a function creates a new *local* one instead.

```python
high_score = 0            # global

def try_beat(new):
    high_score = new     # this makes a NEW local — the global is untouched!

try_beat(500)
print(high_score)        # 0
```

That surprises people, but it's a feature: functions can't accidentally trample outer variables. If you *really* must reassign a global from inside, the `global` keyword unlocks it:

```python
high_score = 0

def try_beat(new):
    global high_score
    if new > high_score:
        high_score = new

try_beat(500)
print(high_score)        # 500
```

:::analogy
Global state is like a shared save file the whole game writes to. It works, but once ten functions can all rewrite it, tracking down *who* changed the score to a weird value becomes a nightmare. Passing values in and returning them out — hand-to-hand — keeps every function honest about what it touched.
:::

That's exactly why `global` is a tool of last resort. Prefer to *return* a new value and let the caller decide what to do with it: `high_score = try_beat(high_score, 500)`. Your functions stay self-contained and easy to reason about.

## Returning multiple values

A function can hand back several results at once by separating them with commas — Python packs them into a tuple, and you **unpack** them on the other side.

```python
def roll_stats():
    return 18, 12, 7      # strength, agility, luck

strength, agility, luck = roll_stats()
print(strength)   # 18
print(luck)       # 7
```

One `return`, three values, three variables. This reads cleanly and saves you from inventing a throwaway container just to smuggle multiple results out of a function.

## `lambda`: tiny inline functions

A **lambda** is a one-line, unnamed function. The syntax is `lambda params: expression` — no `def`, no `return`, just an expression whose value is handed back automatically.

```python
double = lambda n: n * 2
print(double(5))   # 10
```

You wouldn't normally assign a lambda to a name like that (just use `def`). Lambdas shine when a function wants a *small* function as an argument. The classic case is `sorted` with a `key`: you tell it *how* to rank each item.

```python
party = [
    {"name": "Ada", "level": 5},
    {"name": "Linus", "level": 12},
    {"name": "Grace", "level": 8},
]

by_level = sorted(party, key=lambda hero: hero["level"])
for hero in by_level:
    print(hero["name"], hero["level"])
# Ada 5
# Grace 8
# Linus 12
```

The lambda `lambda hero: hero["level"]` pulls out the value `sorted` should compare. Writing a whole named `def` for something used once, right here, would be noise — the lambda keeps the intent local and readable.

:::example
Sort the same party from strongest to weakest by adding `reverse=True`:

```python
strongest = sorted(party, key=lambda hero: hero["level"], reverse=True)
print(strongest[0]["name"])   # Linus
```

The `key` decides *what* to compare; `reverse` decides the direction.
:::

## Practice

:::predict
Q: What does the second call print?
```python
def add_loot(item, bag=[]):
    bag.append(item)
    return bag

add_loot("sword")
print(add_loot("shield"))
```
- ['shield']
- ['sword', 'shield'] *
- ['sword']
E: The default list is created once and shared, so the sword from the first call is still in the bag. This is exactly the mutable-default trap — default to `None` instead.
:::

:::quiz
Q: Inside a function, what does `**kwargs` collect the extra keyword arguments into?
- A tuple
- A dictionary *
- A list
E: `**kwargs` gathers leftover *keyword* arguments into a dict mapping names to values. `*args` is the one that makes a tuple.
:::

## Recap

- **Default values** make parameters optional; they must come after required parameters.
- **Keyword arguments** let you name values at the call site — clearer, and they let you skip defaults.
- `*args` sweeps extra positional arguments into a **tuple**; `**kwargs` sweeps extra keyword arguments into a **dict**.
- Never use a mutable default like `[]` — it's shared across calls. Default to `None` and build a fresh object inside.
- Assigning inside a function makes a **local** variable; `global` reaches out to reassign a global, but prefer returning values instead.
- Return several values at once by comma-separating them, then unpack: `a, b = f()`.
- A `lambda` is a tiny inline function, perfect as the `key` for `sorted`.

**Next up:** Modules and the Standard Library — pulling in Python's huge toolbox and splitting your own code into clean, importable files.

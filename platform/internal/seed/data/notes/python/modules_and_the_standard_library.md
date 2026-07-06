# Modules and the Standard Library

You don't have to build everything from scratch. Python ships with a giant toolbox of ready-made code — the **standard library** — covering math, randomness, dates, data structures, and far more. On top of that, you can split *your own* growing program into tidy files. Both of these happen through one idea: the **module**. A module is just a `.py` file full of code you can pull into another file with `import`.

Fire up the [Code Lab](#code) and try each import as we go.

## What `import` does

When you `import` a module, Python runs that file (once) and gives you access to the names defined inside it. Let's grab the built-in `math` module.

```python
import math

print(math.sqrt(144))     # 12.0
print(math.pi)            # 3.141592653589793
print(math.floor(4.8))    # 4
```

Notice the pattern: after `import math`, everything inside lives behind the `math.` prefix. That prefix is a feature, not a chore — `math.sqrt` tells any reader exactly where `sqrt` came from, and it keeps your own names from colliding with the library's.

:::analogy
A module is like a labelled toolbox on a shelf. `import math` carries the whole box to your bench; `math.sqrt` reaches in and pulls out the specific tool. You could dump every tool loose onto the bench, but then you'd never remember which drawer a wrench belongs to.
:::

## Three styles of import

There's more than one way to bring code in, and each has a moment.

**1. Import the whole module** — the default, safest choice. Everything stays namespaced.

```python
import random
print(random.randint(1, 6))    # a dice roll, 1–6
```

**2. Import specific names** with `from ... import ...`. Now you can use the name directly, no prefix.

```python
from random import randint
print(randint(1, 6))           # no random. prefix needed
```

**3. Rename on import** with `as`, handy for long module names or well-known nicknames.

```python
import statistics as stats
print(stats.mean([10, 20, 30]))    # 20
```

:::warning
Avoid `from module import *` — the star grabs *everything* and drops it into your file with no prefix. Two modules might both define `open` or `time`, and the second import silently clobbers the first. You lose track of where names came from. Import the module, or import the specific names you need by hand.
:::

## Splitting your own code into modules

As a project grows, one enormous file becomes hard to navigate. The cure is to break it into modules — each file a focused set of related functions — and import between them. Say you put your dice logic in `dice.py`:

```python
# dice.py
import random

def roll(sides=6):
    return random.randint(1, sides)

def roll_many(count, sides=6):
    return [roll(sides) for _ in range(count)]
```

Now, from another file in the same folder, you import it by its filename (without the `.py`):

```python
# game.py
import dice

print(dice.roll())          # one d6
print(dice.roll_many(3))    # e.g. [4, 1, 6]
```

Or pull in just the pieces you want:

```python
# game.py
from dice import roll, roll_many
print(roll(20))             # a d20 roll
```

Your program now reads like a set of chapters instead of one endless scroll, and `dice.py` can be reused in the next game you build.

## The `if __name__ == "__main__":` guard

Here's a subtlety. When you *import* `dice.py`, Python runs the whole file top to bottom to define its functions. If `dice.py` also had loose test code at the bottom — say `print(roll())` — that would fire every time someone imported it. Annoying.

Python gives each module a built-in variable, `__name__`. When a file is run **directly**, its `__name__` is the string `"__main__"`. When it's **imported**, `__name__` is the module's name instead (`"dice"`). So you guard your "run directly" code like this:

```python
# dice.py
import random

def roll(sides=6):
    return random.randint(1, sides)

if __name__ == "__main__":
    # Only runs when you execute: python dice.py
    print("Testing the dice:", roll(), roll(), roll())
```

Run `python dice.py` and the test prints. Import `dice` from `game.py` and it stays silent — you only get the functions. This guard is everywhere in real Python; it's the standard way to let a file be *both* an importable library *and* a runnable script.

:::key
`if __name__ == "__main__":` means "only do this when I'm the file being run directly, not when I'm imported." Put demos, tests, and program entry points inside it so importing your module doesn't trigger side effects.
:::

## A tour of the standard library

The standard library is huge; here are modules you'll reach for constantly.

**`random`** — dice, shuffles, and random picks. The heartbeat of any game.

```python
import random

print(random.randint(1, 6))                    # inclusive 1–6
print(random.choice(["sword", "bow", "staff"])) # a random pick
deck = [1, 2, 3, 4, 5]
random.shuffle(deck)                            # shuffles in place
print(deck)
```

**`datetime`** — timestamps, durations, and dates. Perfect for "when did the player last log in?"

```python
from datetime import datetime, timedelta

now = datetime.now()
print(now.year)                        # e.g. 2026
print(now.strftime("%Y-%m-%d %H:%M"))  # formatted timestamp

tomorrow = now + timedelta(days=1)
print(tomorrow.strftime("%A"))         # the weekday name
```

**`math`** — square roots, rounding, constants, trig.

```python
import math
print(math.ceil(4.1))       # 5
print(math.gcd(12, 18))     # 6
```

**`json`** — turn Python data into text (to save a file) and back again. Ideal for save games.

```python
import json

save = {"name": "Ada", "level": 5, "inventory": ["sword", "potion"]}
text = json.dumps(save)         # dict → JSON string
print(text)

loaded = json.loads(text)       # JSON string → dict
print(loaded["level"])          # 5
```

**`collections`** — specialised containers. `Counter` tallies things; `defaultdict` supplies a default for missing keys so you never hit a `KeyError`.

```python
from collections import Counter, defaultdict

rolls = [3, 6, 6, 1, 6, 3]
tally = Counter(rolls)
print(tally)                    # Counter({6: 3, 3: 2, 1: 1})
print(tally.most_common(1))     # [(6, 3)]  — the most frequent roll

scores = defaultdict(int)       # missing keys default to 0
scores["ada"] += 10             # no need to check if "ada" exists first
scores["ada"] += 5
print(scores["ada"])            # 15
print(scores["linus"])          # 0  — auto-created, no error
```

**`statistics`** — quick averages without pulling in a heavyweight library.

```python
import statistics
damage = [12, 15, 9, 20, 15]
print(statistics.mean(damage))    # 14.2
print(statistics.median(damage))  # 15
print(statistics.mode(damage))    # 15
```

:::tip
When you're about to write something fiddly — averaging numbers, counting duplicates, formatting a date — pause and ask "is this in the standard library?" It usually is, it's already tested, and it's faster than rolling your own.
:::

## Reaching beyond: `pip` and third-party packages

The standard library is vast, but not infinite. When you need something it doesn't cover — a web framework, a game engine, a plotting library — you install a **third-party package** from the Python Package Index (PyPI) using `pip`, Python's package installer, from your terminal: `pip install requests`. After that, `import requests` works just like a built-in module. (We'll properly cover environments and dependencies later; for now, just know that's the door to the wider ecosystem.)

## Practice

:::predict
Q: A file `dice.py` ends with `if __name__ == "__main__": print("rolling")`. What happens when another file runs `import dice`?
- It prints "rolling"
- It does NOT print "rolling" *
- It raises an ImportError
E: On import, `__name__` is `"dice"`, not `"__main__"`, so the guarded block is skipped. It would only print if you ran `python dice.py` directly.
:::

:::quiz
Q: Which import style lets you write `randint(1, 6)` with no `random.` prefix?
- `import random`
- `from random import randint` *
- `import random as r`
E: `from random import randint` pulls the name directly into your file. Plain `import random` keeps the prefix; `import random as r` would need `r.randint`.
:::

## Recap

- A **module** is a `.py` file you pull in with `import`; the standard library is a huge collection of them, ready to use.
- Three import styles: `import math` (namespaced), `from random import randint` (direct name), and `import x as y` (renamed). Avoid `import *`.
- Split large programs into your own modules and import between files by filename (no `.py`).
- `if __name__ == "__main__":` runs code only when a file is executed directly, not when it's imported — the home for tests and entry points.
- Keep `random`, `datetime`, `math`, `json`, `collections` (`Counter`, `defaultdict`), and `statistics` in your back pocket.
- For anything beyond the standard library, `pip install <package>` pulls it from PyPI.

**Next up:** Errors and Exceptions — reading tracebacks calmly and writing code that recovers gracefully when things go wrong.

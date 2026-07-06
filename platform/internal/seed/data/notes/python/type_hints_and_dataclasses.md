# Type Hints and Dataclasses

You already write functions, classes, and collections. This lesson adds two tools that make Python code easier to read, safer to edit, and far less boilerplate-heavy: **type hints** for documenting the shapes of your data, and **dataclasses** for building tidy value objects without hand-writing constructors.

## What Type Hints Are

A type hint is a note attached to a variable, parameter, or return value that says what kind of thing it holds. You write them with a colon for variables and parameters, and an arrow for return values.

```python
def add_score(current: int, bonus: int) -> int:
    return current + bonus

health: float = 100.0
name: str = "Ada"
```

Read `current: int` as "current is expected to be an int" and `-> int` as "this function returns an int." Nothing about the *behavior* changed — the function runs exactly as it would without the hints.

:::key
Python does **not** enforce type hints at runtime. They are documentation that happens to be machine-readable. Passing a string where an `int` is hinted will not raise an error on its own — Python trusts you.
:::

So why bother? Three audiences read them:

- **Other humans** (including future you) understand the code faster.
- **Editors** use them for autocomplete and inline warnings.
- **Static checkers** like `mypy` or `pyright` scan your code *before* you run it and flag mismatches.

```python
def greet(player: str) -> str:
    return "Hello, " + player

greet(42)   # runs fine and crashes with a TypeError at runtime,
            # but mypy would flag it as an error *before* you ran it
```

That last point is the real payoff. A checker catches whole categories of bugs at your desk instead of in production.

## Annotating Variables, Parameters, and Returns

The three places you annotate, side by side:

```python
level: int = 1                      # variable

def level_up(level: int) -> int:    # parameters and return
    return level + 1

def save(name: str) -> None:        # returns nothing
    print(f"Saving {name}...")
```

Use `-> None` when a function does its work through side effects (printing, writing a file) and returns nothing useful. It's an explicit "no return value here," which is clearer than leaving the return unannotated.

You don't have to annotate every local variable. Often the value on the right makes the type obvious:

```python
score = 0          # clearly an int, no annotation needed
total = score + 5  # still obviously an int
```

Annotate when the type *isn't* obvious — an empty container, a value from an external source, or anything a reader might guess wrong.

```python
scores: list[int] = []   # without the hint, is this list of what?
```

:::tip
Annotate function signatures generously — they are the contract of your code and the thing readers scan first. Be relaxed about local variables; only annotate the ones whose type isn't clear at a glance.
:::

## Typing Containers

A bare `list` tells the reader almost nothing. What's *in* the list? Since Python 3.9 you write the element type in square brackets, and in 3.12 this is the standard, clean form:

```python
scores: list[int] = [10, 20, 30]
names: list[str] = ["Ada", "Grace"]
```

Dictionaries take two types — the key type and the value type:

```python
inventory: dict[str, int] = {"potion": 3, "sword": 1}
#              key   value
```

Tuples and sets follow the same idea:

```python
point: tuple[int, int] = (5, 12)      # exactly two ints
tags: set[str] = {"boss", "fire"}
```

### When a Value Might Be Missing

Sometimes a value is *either* a real value *or* nothing. The modern way to say "an int, or None" is the `| None` union syntax:

```python
def find_player(name: str) -> str | None:
    players = {"ada", "grace"}
    if name.lower() in players:
        return name
    return None   # not found
```

You'll also see the older spelling using `Optional` from the `typing` module. `Optional[str]` means exactly the same as `str | None`:

```python
from typing import Optional

def find_player(name: str) -> Optional[str]:
    ...
```

Both are correct. In new Python 3.12 code, prefer `str | None` — it reads naturally and needs no import.

:::warning
`str | None` means the value can genuinely be `None`, so callers must handle that case. If you `return None` and the caller does `result.upper()`, you'll hit an `AttributeError` on the `None`. A checker flags this for you — which is exactly why the hint is worth writing.
:::

:::predict
Q: A checker like mypy is analyzing this. Which line would it flag?
```python
def hp(current: int) -> int:
    return current + 10

hp("full")
```
- `def hp(current: int) -> int:`
- `return current + 10`
- `hp("full")`*
- None of them
E: The call passes a `str` where an `int` is expected. mypy flags the call site `hp("full")` — even though Python itself would not stop you from writing it.
:::

## Meet the Dataclass

Suppose you want a small object to hold an item's data. The plain-class way is a lot of repetitive typing:

```python
class Item:
    def __init__(self, name: str, damage: int, price: int):
        self.name = name
        self.damage = damage
        self.price = price
```

Every field is named **three times**. And you still don't get a readable `print()` output or a way to compare two items for equality. A **dataclass** generates all of that for you. Decorate a class with `@dataclass` and just declare the fields with type hints:

```python
from dataclasses import dataclass

@dataclass
class Item:
    name: str
    damage: int
    price: int
```

That's the whole class. Behind the scenes, `@dataclass` writes three methods for you:

```python
sword = Item("Sword", damage=12, price=50)

print(sword)          # → Item(name='Sword', damage=12, price=50)   (__repr__)
print(sword.damage)   # → 12

other = Item("Sword", 12, 50)
print(sword == other) # → True   (__eq__ compares field by field)
```

- `__init__` — the constructor that takes `name`, `damage`, `price`.
- `__repr__` — a helpful printout showing every field, great for debugging.
- `__eq__` — compares two items field by field, so equal data means equal objects.

With a plain class, `print(sword)` would show something useless like `<__main__.Item object at 0x104f2>`, and `sword == other` would be `False` even with identical data. The dataclass fixes both.

:::analogy
A dataclass is like a fill-in-the-blank form. You write the field labels — `name`, `damage`, `price` — and Python prints the form, stamps a comparison rule on it, and files it away. You never handwrite the same three boxes over and over.
:::

## Fields With Defaults and Methods

Fields can have default values, just like function parameters. Give a default with `=`:

```python
from dataclasses import dataclass

@dataclass
class Player:
    name: str
    level: int = 1          # defaults
    health: int = 100
    score: int = 0

hero = Player("Ada")                 # uses all the defaults
boss = Player("Ada", level=9, health=500)
print(hero)   # → Player(name='Ada', level=1, health=100, score=0)
```

The same rule as function parameters applies: **fields with defaults must come after fields without them.** Put your required fields first, your optional ones last.

:::warning
Never use a mutable default like `items: list[str] = []` directly — every instance would secretly share the *same* list. Python even raises an error to stop you. Use `field(default_factory=list)` instead, which builds a fresh list per instance:

```python
from dataclasses import dataclass, field

@dataclass
class Player:
    name: str
    inventory: list[str] = field(default_factory=list)
```
:::

Dataclasses are still ordinary classes, so you can add your own methods. They read cleanly because the data is right there in the fields:

```python
from dataclasses import dataclass, field

@dataclass
class Player:
    name: str
    health: int = 100
    inventory: list[str] = field(default_factory=list)

    def take_damage(self, amount: int) -> None:
        self.health = max(0, self.health - amount)

    def is_alive(self) -> bool:
        return self.health > 0

    def pick_up(self, item: str) -> None:
        self.inventory.append(item)

ada = Player("Ada")
ada.take_damage(30)
ada.pick_up("potion")
print(ada.health)     # → 70
print(ada.is_alive()) # → True
print(ada.inventory)  # → ['potion']
```

The `@dataclass` decorator only generates the boilerplate (`__init__`, `__repr__`, `__eq__`). Your own methods are left exactly as you wrote them.

:::predict
Q: What does this print?
```python
from dataclasses import dataclass

@dataclass
class Item:
    name: str
    price: int = 0

a = Item("Gem", 100)
b = Item("Gem", 100)
print(a == b)
```
- True*
- False
- Item(name='Gem', price=100)
- An error
E: `@dataclass` generates an `__eq__` that compares field by field. Both items have the same `name` and `price`, so `a == b` is `True` — unlike a plain class, where it would be `False`.
:::

## When to Reach for a Dataclass

You now have three ways to group related values. Pick by what the data is *for*:

**A plain `dict`** is best for loose, dynamic, short-lived data — parsing some JSON, passing a bag of options around, keys you don't know ahead of time.

```python
config = {"volume": 0.8, "difficulty": "hard"}
```

The downside: nothing stops a typo like `config["diffculty"]`, there's no autocomplete, and every access is a string lookup.

**A dataclass** is best when the data has a *fixed, known shape* that appears throughout your program — a `Player`, an `Item`, a `Position`. You get named fields with types, editor autocomplete, a readable `repr`, free equality, and typo protection.

```python
player = Player("Ada", level=5)
player.levl   # editor/checker catches this immediately; a dict typo wouldn't
```

**A plain class** earns its keep when the object is mostly *behavior* rather than data — lots of methods, custom construction logic, no natural set of "just fields." If you're mainly bundling fields together, the dataclass saves you the boilerplate.

:::tip
A quick test: if you were about to write an `__init__` that just copies each argument onto `self`, stop — that's exactly the job `@dataclass` does for free. Reach for it whenever a class is "a name and some typed fields."
:::

## Recap

- Type hints annotate variables (`x: int`), parameters (`def f(x: int)`), and returns (`-> int`). Use `-> None` for functions that return nothing useful.
- Python does **not** enforce hints at runtime — they're for human readers, editors, and static checkers like mypy that catch mismatches before you run the code.
- Type containers with their contents: `list[int]`, `dict[str, int]`, `tuple[int, int]`, `set[str]`.
- For "a value or nothing," prefer `str | None` (modern) over `Optional[str]` (older but equivalent).
- `@dataclass` auto-generates `__init__`, `__repr__`, and `__eq__` from field declarations, killing the copy-onto-self boilerplate.
- Fields can have defaults (put them last); use `field(default_factory=list)` for mutable defaults, never a bare `[]`.
- Dataclasses are normal classes — add your own methods freely.
- Choose a dataclass for fixed, known data shapes that recur; a dict for loose dynamic data; a plain class for behavior-heavy objects.

**Next up:** Functional Tools — `map`/`filter`, sorting with `key=`, `reduce`, memoization with `lru_cache`, and a tour of `itertools`.

# Inheritance and Dunder Methods

You already know how to build a class like `Player` or `Enemy`. This lesson adds two superpowers. **Inheritance** lets one class build on another so you don't repeat yourself. **Dunder methods** — the ones with double underscores like `__str__` — let your own objects plug into Python's built-in tools: `print()`, `==`, `len()`, and `sorted()`. Together they make your objects feel like a natural part of the language.

## Inheritance: Building on What Exists

Imagine you already have an `Enemy` class. Now you want a `Boss` — which is *exactly* like an enemy, but tougher and with a special move. You *could* copy all of `Enemy`'s code into a new class. Don't. Instead, have `Boss` **inherit** from `Enemy`.

```python
class Enemy:
    def __init__(self, name, hp, power):
        self.name = name
        self.hp = hp
        self.power = power

    def take_damage(self, amount):
        self.hp = self.hp - amount

    def is_alive(self):
        return self.hp > 0


class Boss(Enemy):        # Boss inherits everything from Enemy
    pass
```

Writing `class Boss(Enemy):` means "a `Boss` is a kind of `Enemy`." The parentheses name the **parent class** (also called the base or superclass). Even with just `pass`, `Boss` already has everything `Enemy` has:

```python
dragon = Boss("Dragon", 500, power=40)
dragon.take_damage(100)
print(dragon.hp)          # → 400
print(dragon.is_alive())  # → True
```

We wrote no `__init__`, no `take_damage`, no `is_alive` on `Boss` — it inherited all of them. That's the whole point: shared behavior lives in one place.

## The "Is-A" Test

Inheritance models an **"is-a" relationship**. A `Boss` *is an* `Enemy`. A cat *is an* animal. Before you reach for inheritance, say the sentence out loud: "A Boss is an Enemy." If that sounds right, inheritance fits.

If instead you'd say "*has a*" — a Player *has an* Inventory — that's not inheritance. That's just an attribute holding another object. Don't make `Player` inherit from `Inventory`; a player isn't a kind of inventory.

:::analogy
Inheritance is a family tree, not a toolbox. A `Boss` inherits from `Enemy` the way a child inherits traits from a parent — it *is* one of them, with its own extras. If you're tempted to inherit just to grab a handy method, ask whether the "is-a" sentence is actually true. If it isn't, you want a plain attribute, not a parent.
:::

## Calling the Parent with `super().__init__(...)`

Usually a subclass wants to *add* something. A `Boss` might have an `enrage_threshold` on top of the usual name, hp, and power. You write a new `__init__` — but you don't want to re-type the lines that set `name`, `hp`, and `power`. Let the parent handle those with `super()`:

```python
class Boss(Enemy):
    def __init__(self, name, hp, power, enrage_threshold):
        super().__init__(name, hp, power)   # let Enemy set name, hp, power
        self.enrage_threshold = enrage_threshold  # Boss adds this

    def is_enraged(self):
        return self.hp <= self.enrage_threshold
```

`super()` refers to the parent class. So `super().__init__(name, hp, power)` runs `Enemy`'s `__init__`, which sets up the shared attributes. Then `Boss` adds its own.

```python
dragon = Boss("Dragon", 500, power=40, enrage_threshold=150)
print(dragon.power)          # → 40   (set by Enemy via super)
print(dragon.is_enraged())   # → False

dragon.take_damage(400)
print(dragon.hp)             # → 100
print(dragon.is_enraged())  # → True  (100 <= 150)
```

:::warning
If your subclass defines its own `__init__` and forgets to call `super().__init__(...)`, the parent's setup never runs — so attributes like `self.hp` won't exist, and the first method that uses them crashes with an `AttributeError`. When you write a subclass `__init__`, calling `super().__init__(...)` is almost always your first line.
:::

## Overriding Methods

A subclass can also *replace* a parent's method by defining one with the same name. This is called **overriding**. A `Boss` might hit harder than a normal enemy:

```python
class Boss(Enemy):
    def __init__(self, name, hp, power, enrage_threshold):
        super().__init__(name, hp, power)
        self.enrage_threshold = enrage_threshold

    def attack(self, target):
        damage = self.power
        if self.hp <= self.enrage_threshold:
            damage = self.power * 2      # enraged bosses hit twice as hard
        target.take_damage(damage)
        print(f"{self.name} strikes for {damage}!")
```

When you call `dragon.attack(zed)`, Python uses the `Boss` version, not the `Enemy` one. The subclass gets the final say on its own behavior while still inheriting everything it *didn't* override.

## Dunder Methods: Hooking Into Python

Now the second superpower. **Dunder methods** ("double underscore," like `__init__`) are special method names Python looks for automatically. You've already used one — `__init__` runs when you build an object. There are many more, and each one lets your object respond to a built-in operation. Define the right dunder and your object suddenly works with `print()`, `==`, `len()`, or `sorted()`.

You never call these directly (you don't write `zed.__str__()`). You define them, and Python calls them for you at the right moment.

## `__str__` vs `__repr__`: Printing Your Objects

By default, printing an object gives something useless like `<__main__.Player object at 0x104f2c9d0>`. Define `__str__` to control what `print()` shows:

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp

    def __str__(self):
        return f"{self.name} ({self.hp} HP)"

zed = Player("Zed", 100)
print(zed)          # → Zed (100 HP)
print(f"Fighting {zed}!")   # → Fighting Zed (100 HP)!
```

`__str__` must **return** a string (don't `print` inside it). Python calls it whenever your object needs to be shown to a human — `print()`, `str()`, and f-strings all use it.

There's a close cousin, `__repr__`. Where `__str__` is the friendly, human-facing version, `__repr__` is the developer-facing one — meant to be unambiguous, ideally showing how the object could be recreated. It's what you see in the interactive shell and inside lists.

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp

    def __repr__(self):
        return f"Player({self.name!r}, {self.hp})"

zed = Player("Zed", 100)
print([zed])        # → [Player('Zed', 100)]  ← uses __repr__
```

:::key
Rule of thumb: `__str__` is for your users, `__repr__` is for you the developer. If you only write one, make it `__repr__` — Python falls back to it when `__str__` is missing, so you get useful output everywhere.
:::

## `__eq__`: Making `==` Meaningful

By default, `==` between two objects checks whether they're the *exact same object* in memory — so two players with identical stats count as different. Define `__eq__` to compare by value instead:

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp

    def __eq__(self, other):
        return self.name == other.name and self.hp == other.hp

a = Player("Zed", 100)
b = Player("Zed", 100)

print(a == b)   # → True   (same name and hp)
print(a is b)   # → False  (still two separate objects)
```

Notice `==` and `is` now mean different things: `==` asks "are these equal?" (your rule), while `is` asks "are these literally the same object?" (identity). Your `__eq__` defines what "equal" means for a `Player`.

## `__len__`: Working With `len()`

If your object is a collection of things, `__len__` lets `len()` work on it. Picture a `Party` of players:

```python
class Party:
    def __init__(self, members):
        self.members = members

    def __len__(self):
        return len(self.members)

party = Party([Player("Zed", 100), Player("Nyx", 80)])
print(len(party))   # → 2
```

`len(party)` quietly calls `party.__len__()`. Your object now behaves like a built-in container.

## `__lt__`: Sorting With `sorted()`

To sort your objects, Python needs to know how to compare two of them. Define `__lt__` ("less than") and `sorted()` — plus `min()`, `max()`, and the `<` operator — all start working. Let's sort players by HP:

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp

    def __lt__(self, other):
        return self.hp < other.hp

    def __repr__(self):
        return f"{self.name}({self.hp})"

roster = [Player("Zed", 100), Player("Nyx", 80), Player("Kai", 120)]
print(sorted(roster))   # → [Nyx(80), Zed(100), Kai(120)]
```

`sorted()` compares pairs of players with `<`, which calls your `__lt__`. Because you said "less HP comes first," the list comes out ordered from weakest to strongest. Change the comparison and you change the sort order.

:::tip
You don't have to memorize every dunder. Just remember the pattern: a built-in operation (`print`, `==`, `len`, `sorted`) has a matching dunder (`__str__`, `__eq__`, `__len__`, `__lt__`). Define the dunder, and the built-in "just works" on your object.
:::

## Practice

:::predict
Q: What does this print?
```python
class Enemy:
    def __init__(self, name):
        self.name = name
    def describe(self):
        return f"a wild {self.name}"

class Boss(Enemy):
    def describe(self):
        return f"THE MIGHTY {self.name}"

print(Boss("Dragon").describe())
```
- a wild Dragon
- THE MIGHTY Dragon*
- Dragon
- error
E: `Boss` overrides `describe`, so its version runs instead of the inherited `Enemy` one. `Boss("Dragon")` still works because `Boss` inherits `Enemy.__init__`.
:::

:::quiz
Q: You want `sorted(enemies)` to order enemies by their `hp`. Which dunder method should you define?
- `__str__`
- `__len__`
- `__lt__`*
- `__init__`
E: `sorted()` compares items with `<`, which calls `__lt__`. Defining `__lt__` to compare `hp` tells Python how to order your objects.
:::

## Recap

- **Inheritance** (`class Boss(Enemy):`) lets a subclass reuse a parent's attributes and methods. Use it for true **"is-a"** relationships; use a plain attribute for "has-a".
- Call `super().__init__(...)` in a subclass to run the parent's setup, then add the subclass's own attributes.
- **Overriding** a method (same name in the subclass) replaces the parent's version for that class.
- **Dunder methods** hook your objects into Python's built-ins. You define them; Python calls them.
- `__str__` controls `print()` (human-friendly); `__repr__` is the developer-facing fallback. `__eq__` powers `==` by value. `__len__` powers `len()`. `__lt__` powers `sorted()`, `min()`, and `max()`.

**Next up:** the interactive lab — you'll build a `Player` class from scratch with `__init__`, `hp`, `take_damage`, `is_alive`, and `__str__`, bringing everything from these two lessons together.

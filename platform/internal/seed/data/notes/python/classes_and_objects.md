# Classes and Objects

Up to now you've worked with Python's built-in types: strings, numbers, lists, dictionaries. Classes let you build *your own* types — a `Player`, an `Enemy`, an `Item` — that bundle data and behavior together. This is the heart of object-oriented programming, and it's how nearly all real software is structured.

## What a Class Is (and What an Object Is)

A **class** is a blueprint. It describes what a certain kind of thing *has* and what it can *do*. An **object** (also called an **instance**) is one concrete thing built from that blueprint.

Think of it this way: a class is the idea of "a player." An object is *this specific* player named "Zed" with 100 HP. From one class you can stamp out as many objects as you like — one hundred players, each with their own name and health, all sharing the same rules for how a player works.

:::analogy
A class is a cookie cutter; an object is a cookie. The cutter defines the *shape* — every cookie made from it is a "player-shaped" cookie. But each cookie is its own physical thing: you can frost one and eat another. Changing one cookie doesn't touch the rest.
:::

The class defines the shape once. The objects are the real, individual things your program actually works with.

## Defining a Class

You define a class with the `class` keyword. By convention, class names use **PascalCase** — every word capitalized, no underscores.

```python
class Player:
    pass
```

That `pass` is just a placeholder that means "empty for now." This class is valid but useless — it holds no data and does nothing. Let's create an object from it anyway:

```python
class Player:
    pass

zed = Player()          # build one instance
print(zed)              # → <__main__.Player object at 0x104f2c9d0>
print(type(zed))        # → <class '__main__.Player'>
```

Calling `Player()` — the class name followed by parentheses — creates a new object. Right now `zed` has no name and no health. We fix that with `__init__`.

## `__init__` and `self`

The `__init__` method runs automatically every time you create an object. Its job is to set up that object's starting data. The name is Python's convention for "initialize" — those are **double underscores** on each side (you'll hear them called "dunder" methods, short for "double underscore").

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp
```

Two things need explaining here.

First, **`self`**. It is always the first parameter of a method, and it refers to *the specific object being worked on*. When you write `self.name = name`, you're saying "store this name on *this particular* player." Each object has its own `self`, which is how one player can be "Zed" with 100 HP while another is "Nyx" with 80.

Second, the difference between `name` and `self.name`. The plain `name` is the value passed in as an argument. The `self.name` is where you *save* it onto the object so it sticks around after `__init__` finishes.

```python
class Player:
    def __init__(self, name, hp):
        self.name = name       # save the argument onto the object
        self.hp = hp

zed = Player("Zed", 100)       # __init__ runs with name="Zed", hp=100
print(zed.name)                # → Zed
print(zed.hp)                  # → 100
```

Notice you *don't* pass anything for `self` yourself. Python fills it in automatically — `Player("Zed", 100)` becomes a behind-the-scenes call with `self` set to the new object.

:::key
`self` is not magic and it's not a keyword — it's just the conventional name for "this object." Every method that works on an instance takes `self` as its first parameter, and you access the object's own data through it: `self.hp`, `self.name`.
:::

## Instance Attributes

The values stored on an object — `self.name`, `self.hp` — are called **instance attributes**. "Instance" because each object gets its own copy. Two players made from the same class have completely independent attributes:

```python
zed = Player("Zed", 100)
nyx = Player("Nyx", 80)

print(zed.hp)    # → 100
print(nyx.hp)    # → 80

zed.hp = 50      # change one player's HP
print(zed.hp)    # → 50
print(nyx.hp)    # → 80  (unaffected — separate object)
```

You read an attribute with `object.attribute` and change it the same way. Each object minds its own business.

## Methods: Giving Objects Behavior

A **method** is a function defined inside a class. It describes something an object can *do*. Because methods live on the object, they can read and change that object's own attributes through `self`.

Let's give a `Player` the ability to take damage:

```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp

    def take_damage(self, amount):
        self.hp = self.hp - amount

    def is_alive(self):
        return self.hp > 0
```

Now the behavior travels *with* the data. To call a method, use `object.method(...)`:

```python
zed = Player("Zed", 100)

zed.take_damage(30)
print(zed.hp)          # → 70

zed.take_damage(50)
print(zed.hp)          # → 20

print(zed.is_alive())  # → True

zed.take_damage(25)
print(zed.hp)          # → -5
print(zed.is_alive())  # → False
```

Again, you don't pass `self` — `zed.take_damage(30)` automatically sends `zed` in as `self`. So `take_damage` reaches into *this* player's `hp` and lowers it.

:::tip
A quick way to read a method call: `zed.take_damage(30)` means "tell *zed* to take 30 damage." The object before the dot is the one the method acts on.
:::

## Why Classes Beat Loose Dictionaries

You could represent a player with a dictionary instead:

```python
zed = {"name": "Zed", "hp": 100}

def take_damage(player, amount):
    player["hp"] = player["hp"] - amount

take_damage(zed, 30)
print(zed["hp"])   # → 70
```

This works for a while. But as soon as *behavior* gets involved, dictionaries start to hurt. The data and the functions that operate on it drift apart — `take_damage` lives somewhere far from the dictionary it's meant to work on. Nothing stops you from mistyping `"hp"` as `"HP"`, or from creating a "player" that's missing a field entirely. And every function has to remember the exact string keys.

With a class, the data and its behavior live in one place. The blueprint guarantees every player has a `name` and an `hp`. The methods are attached, discoverable, and can't be called on the wrong shape of data. Your `Player` becomes a real, self-contained concept instead of a bag of strings you have to handle carefully.

:::warning
Dictionaries are perfect for plain data with no attached behavior — a config, a row from a file, a JSON response. Reach for a class the moment your data grows *rules* and *actions*: things it can do, invariants it must keep. "Nouns that do things" want to be classes.
:::

## Modeling an Enemy

The same pattern builds any concept in your game world. Here's an `Enemy`, which looks a lot like `Player` but adds an attack power:

```python
class Enemy:
    def __init__(self, name, hp, power):
        self.name = name
        self.hp = hp
        self.power = power

    def attack(self, target):
        target.take_damage(self.power)
        print(f"{self.name} hits {target.name} for {self.power}!")
```

Now enemies and players can interact — because `attack` calls the target's own `take_damage` method:

```python
zed = Player("Zed", 100)
goblin = Enemy("Goblin", 30, power=15)

goblin.attack(zed)     # → Goblin hits Zed for 15!
print(zed.hp)          # → 85
```

Two different classes, cooperating through methods. This is how object-oriented programs are built: small, well-defined objects that send each other messages.

## Class vs. Instance Attributes

Everything so far has used *instance* attributes — set on `self`, unique per object. There's a second kind: a **class attribute**, defined directly in the class body and *shared by every instance*.

```python
class Player:
    max_hp = 100          # class attribute — shared by all players

    def __init__(self, name):
        self.name = name  # instance attribute — unique per player
        self.hp = Player.max_hp

zed = Player("Zed")
nyx = Player("Nyx")

print(zed.max_hp)   # → 100
print(nyx.max_hp)   # → 100  (same shared value)
```

Use a class attribute for things that are the same for *every* object of that type — a cap, a default, a constant. Use an instance attribute (on `self`) for anything that differs from one object to the next, like a name or current HP. When in doubt, prefer instance attributes; shared mutable class attributes can surprise you.

## Practice

:::predict
Q: What does this print?
```python
class Player:
    def __init__(self, name, hp):
        self.name = name
        self.hp = hp
    def take_damage(self, amount):
        self.hp = self.hp - amount

zed = Player("Zed", 100)
zed.take_damage(40)
print(zed.hp)
```
- 100
- 60*
- 40
- error
E: `__init__` sets `hp` to 100, then `take_damage(40)` subtracts 40 through `self.hp`, leaving 60.
:::

:::quiz
Q: What is `self` inside a method?
- A Python keyword you must never rename
- A reference to the specific object the method is acting on*
- The class itself, shared by all instances
- The value returned by `__init__`
E: `self` is the conventional name for the current instance — the particular object the method was called on. It's how a method reaches that object's own attributes.
:::

## Recap

- A **class** is a blueprint; an **object** (instance) is a concrete thing built from it. One class, many independent objects.
- Define a class with `class Name:` (PascalCase). Create an instance by calling it: `Player("Zed", 100)`.
- `__init__` runs on creation and sets up starting data. `self` is the current object; `self.hp = hp` saves data onto it.
- **Instance attributes** live on `self` and are unique per object. **Methods** are functions in the class that act on the object through `self`.
- Prefer a class over a loose dictionary once data gains behavior and rules — it keeps data and actions together and guarantees a consistent shape.
- **Class attributes** are shared by every instance; **instance attributes** differ per object.

**Next up:** Inheritance and Dunder Methods — building specialized classes from existing ones, and teaching your objects to work with `print()`, `==`, `len()`, and `sorted()`.

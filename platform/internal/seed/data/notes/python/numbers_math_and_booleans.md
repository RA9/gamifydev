# Numbers, Math, and Booleans

Games run on numbers: damage dealt, XP earned, gold spent, health remaining, whether you're alive or defeated. This lesson covers the two number types Python uses, the full set of math operators (including a couple you may not have seen), the handy built-in math helpers, how to convert between types, and the true/false logic that drives every decision your programs make.

## Two Kinds of Numbers: int and float

Python has two everyday number types. An **int** (integer) is a whole number with no decimal point. A **float** is a number *with* a decimal point.

```python
level = 7           # int    — a whole number
gold = 250          # int
health_pct = 0.75   # float  — has a decimal point
crit_mult = 1.5     # float
```

You can check a value's type with `type()`:

```python
print(type(7))      # → <class 'int'>
print(type(1.5))    # → <class 'float'>
```

The moment a decimal enters the picture, the result becomes a float — even if the fraction is zero:

```python
print(type(10 / 2))   # → <class 'float'>   (division always gives a float)
print(10 / 2)         # → 5.0               (note the .0)
```

:::key
`int` = whole numbers (`3`, `-40`, `9001`). `float` = numbers with a decimal (`0.5`, `3.14`, `5.0`). Regular division with `/` **always** produces a float, even `4 / 2` which gives `2.0`.
:::

## Arithmetic Operators

The basic operators work like the math you already know:

```python
print(10 + 3)    # → 13   addition
print(10 - 3)    # → 7    subtraction
print(10 * 3)    # → 30   multiplication
print(10 / 3)    # → 3.333333...  division (always a float)
```

Python adds three more that are especially useful in games.

**Floor division** `//` divides and throws away the remainder, giving a whole-number result:

```python
print(10 // 3)   # → 3    (3.33 rounded down)
print(17 // 5)   # → 3    (how many full stacks of 5 fit in 17)
```

**Modulo** `%` gives the *remainder* left over after division:

```python
print(10 % 3)    # → 1    (10 = 3*3 + 1, remainder 1)
print(17 % 5)    # → 2    (2 items left over after full stacks)
```

**Exponent** `**` raises a number to a power:

```python
print(2 ** 3)    # → 8    (2 * 2 * 2)
print(5 ** 2)    # → 25   (5 squared)
```

:::analogy
Think of `//` and `%` as splitting loot into stacks. If you have 17 arrows and stacks hold 5, then `17 // 5` is `3` (three full stacks) and `17 % 5` is `2` (two arrows left over). Together they account for everything.
:::

## Operator Precedence

When several operators appear in one expression, Python follows an order — the same **PEMDAS** rules from math class. Powers first, then multiply/divide (including `//` and `%`), then add/subtract. Whatever's in **parentheses** happens first of all.

```python
print(2 + 3 * 4)       # → 14   (multiply first: 3*4=12, then +2)
print((2 + 3) * 4)     # → 20   (parentheses first: 5*4)
print(2 ** 3 * 2)      # → 16   (power first: 8, then *2)
```

Consider a damage calculation. Base damage plus a bonus, then doubled by a critical hit — parentheses make the intent unmistakable:

```python
base = 10
bonus = 5
damage = (base + bonus) * 2
print(damage)          # → 30
```

:::tip
When in doubt, add parentheses. Even when they don't change the result, they make your intent obvious to the next person reading the code — often future you.
:::

## Built-in Math Helpers

Python ships with small, ready-to-use functions for common jobs — no imports needed.

`round()` rounds to the nearest whole number (or to a given number of decimals):

```python
print(round(3.7))        # → 4
print(round(3.14159, 2)) # → 3.14   (keep 2 decimal places)
```

`abs()` gives the absolute value — distance from zero, always positive:

```python
print(abs(-15))          # → 15    (e.g. size of a health change, ignoring direction)
```

`min()` and `max()` pick the smallest or largest of several values:

```python
print(min(30, 12, 45))   # → 12
print(max(30, 12, 45))   # → 45
```

These two are perfect for clamping values. Want health that never drops below 0 or above a cap?

```python
health = 120
health = max(0, min(100, health))   # clamp into the 0–100 range
print(health)                       # → 100
```

## Converting Between Types

Sometimes you have the right value in the wrong type — a number stuck inside a string, or a float where you need a whole count. Python gives you conversion functions: `int()`, `float()`, and `str()`.

```python
print(int("42"))     # → 42     (string → int)
print(float("3.5"))  # → 3.5    (string → float)
print(str(99))       # → "99"   (int → string)
```

Converting a float to an int **truncates** — it chops off the decimal part, it does not round:

```python
print(int(3.9))      # → 3      (not 4! the .9 is simply dropped)
```

This matters a lot with user input. Text typed by a player arrives as a **string**, so you must convert it before doing math:

```python
typed = "50"                 # imagine this came from the player
xp = int(typed) + 10
print(xp)                    # → 60
```

:::warning
`"50" + 10` doesn't add — it raises a `TypeError`, because you can't add a string to an int. Convert first with `int("50")`. And `int("hello")` fails too: the conversion only works when the text really is a number.
:::

## Booleans: True and False

A **boolean** is a value that is either `True` or `False` — nothing else. It's how a program represents a yes/no fact: is the player alive? is the quest done? Booleans are their own type, and the capital letter matters.

```python
is_alive = True
quest_done = False
print(type(is_alive))    # → <class 'bool'>
```

**Comparison operators** compare two values and produce a boolean:

```python
print(5 > 3)     # → True
print(5 < 3)     # → False
print(5 >= 5)    # → True    (greater than or equal)
print(5 <= 4)    # → False
print(5 == 5)    # → True    (equal — note the DOUBLE equals)
print(5 != 3)    # → True    (not equal)
```

:::warning
Comparing for equality uses **two** equals signs: `==`. A single `=` means *assign a value to a variable*, which is a completely different thing. Writing `if hp = 0` is an error; you want `if hp == 0`.
:::

A quick game example — is the player at full health?

```python
hp = 100
max_hp = 100
print(hp == max_hp)      # → True
```

## Combining Booleans: and, or, not

Real decisions often depend on more than one fact. Python combines booleans with `and`, `or`, and `not`.

- `and` is `True` only when **both** sides are true.
- `or` is `True` when **at least one** side is true.
- `not` flips a boolean to its opposite.

```python
print(True and False)    # → False   (both must be true)
print(True or False)     # → True    (at least one is true)
print(not True)          # → False   (flipped)
```

They shine when combined with comparisons. Can the player use a special ability — enough mana *and* not stunned?

```python
mana = 40
is_stunned = False

can_cast = mana >= 30 and not is_stunned
print(can_cast)          # → True
```

Or check whether a level is over — boss defeated *or* timer ran out:

```python
boss_hp = 0
time_left = 12

level_over = boss_hp <= 0 or time_left <= 0
print(level_over)        # → True   (boss is down, so it's over regardless of the timer)
```

:::key
`and` needs **both** true. `or` needs **at least one** true. `not` **flips** a boolean. Comparisons (`>`, `<`, `==`, `!=`, `>=`, `<=`) produce the booleans you feed into them.
:::

## Practice

:::predict
Q: What does this print?
```python
print(17 % 5)
```
- 3
- 2 *
- 3.4
- 12
E: `%` is the modulo (remainder) operator. 17 divided by 5 is 3 with 2 left over, so `17 % 5` is `2`. (Floor division `17 // 5` would give the `3`.)
:::

:::predict
Q: What does this print?
```python
mana = 20
is_stunned = False
print(mana >= 30 and not is_stunned)
```
- True
- False *
- 20
- None
E: `mana >= 30` is `False` (20 is less than 30). With `and`, both sides must be true, so the whole expression is `False` — the second condition never even matters here.
:::

## Recap

- Python has **int** (whole numbers) and **float** (numbers with a decimal). Division with `/` always yields a float.
- Arithmetic: `+ - * /`, plus `//` (floor division), `%` (remainder/modulo), and `**` (exponent).
- **Precedence** follows PEMDAS; parentheses run first and make intent clear.
- Built-in helpers: `round()`, `abs()`, `min()`, `max()` — great for clamping values into a range.
- Convert types with `int()`, `float()`, `str()`; converting a float to int **truncates** (drops the decimal), and user input is always a string until you convert it.
- **Booleans** are `True`/`False`. Comparison operators (`==`, `!=`, `<`, `>`, `<=`, `>=`) produce them — and `==` (double equals) tests equality, not `=`.
- Combine booleans with `and` (both true), `or` (at least one), and `not` (flip).

**Next up:** Making Decisions — using these booleans with `if`, `elif`, and `else` to branch your programs.

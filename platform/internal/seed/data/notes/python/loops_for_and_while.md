# Loops: for and while

Computers are brilliant at doing the same thing over and over without getting bored. That's what **loops** are for. Spawn ten enemies, count a timer down from sixty, add up every score in a match — you don't write the line ten times, you write it once and tell Python to repeat it. Python gives you two loops: the `for` loop, for when you know *what* you're stepping through, and the `while` loop, for when you're repeating *until* some condition changes. Master both and you can automate almost anything.

This is a lesson made for the [Code Lab](#code). Run every example and change the numbers to see how the output shifts.

## The for loop: stepping through a collection

A `for` loop walks through the items of a collection — a list, a string, anything you can iterate — handing you one item at a time.

```python
enemies = ["goblin", "orc", "troll"]
for enemy in enemies:
    print(f"Spawning a {enemy}!")
# Spawning a goblin!
# Spawning a orc!
# Spawning a troll!
```

The variable `enemy` is fresh each time around — first it's `"goblin"`, then `"orc"`, then `"troll"`. The loop ends automatically when the list runs out. You can loop over a string the same way, one character at a time:

```python
for letter in "GO":
    print(letter)
# G
# O
```

:::analogy
A `for` loop is like a conveyor belt in a factory. Items ride past you one by one, you do the same action to each as it arrives, and when the belt is empty the work stops on its own. You never have to count how many items there are — the belt just runs until it's done.
:::

## Counting with range()

Often you want to repeat something a set number of times, not walk an existing list. That's what `range()` is for — it generates a sequence of numbers to loop over.

```python
for i in range(3):
    print(f"Wave {i}")
# Wave 0
# Wave 1
# Wave 2
```

`range(3)` gives `0, 1, 2` — it starts at 0 and **stops before** the number you give it. That "stops before" catches everyone at first, so keep it in mind. You can also give a start and a step:

```python
for i in range(1, 4):        # start at 1, stop before 4
    print(i)                 # 1, 2, 3

for i in range(0, 10, 2):    # start, stop, step
    print(i)                 # 0, 2, 4, 6, 8
```

The three forms are `range(stop)`, `range(start, stop)`, and `range(start, stop, step)`. A negative step even counts backward — perfect for a countdown:

```python
for t in range(5, 0, -1):
    print(t)                 # 5, 4, 3, 2, 1
print("Go!")
```

:::key
`range(stop)` counts from 0 up to *but not including* `stop`. Add a start (`range(start, stop)`) or a step (`range(start, stop, step)`) for more control, and use a negative step to count down. Pair `range()` with a `for` loop whenever you need to repeat something a known number of times.
:::

## enumerate(): index and value together

Sometimes you need both the item *and* its position. You could track a counter by hand, but Python's `enumerate()` does it for you — it hands back the index and the value on every pass.

```python
players = ["Ada", "Grace", "Linus"]
for index, name in enumerate(players):
    print(f"{index}: {name}")
# 0: Ada
# 1: Grace
# 2: Linus
```

Each loop, `enumerate` gives a `(index, value)` pair that you unpack into two variables. It's cleaner and less error-prone than keeping a separate `count = 0` that you have to remember to increase.

## Looping over a dictionary

To walk through a dictionary, `.items()` gives you each key and value together — the same tidy unpacking trick:

```python
scores = {"Ada": 90, "Grace": 85, "Linus": 70}
for name, score in scores.items():
    print(f"{name} scored {score}")
# Ada scored 90
# Grace scored 85
# Linus scored 70
```

This is the standard way to process every entry in a dictionary. (Looping the dictionary directly gives just the keys; `.items()` gives you both, which is usually what you want.)

## The while loop: repeat until a condition changes

A `for` loop runs a known number of times. A `while` loop runs *as long as a condition stays true* — you use it when you don't know the count ahead of time, only the stopping condition.

```python
health = 3
while health > 0:
    print(f"Health: {health}")
    health = health - 1
print("Game over!")
# Health: 3
# Health: 2
# Health: 1
# Game over!
```

Python checks the condition `health > 0` before each pass. As long as it's true, the body runs. The crucial part is that something *inside* the loop must eventually make the condition false — here, `health` drops by 1 each time until it hits 0.

## Infinite loops and how to avoid them

If the condition never becomes false, the loop runs forever — an **infinite loop** that freezes your program. This almost always happens because you forgot to update the variable in the condition.

```python
# DANGER — do not run this
health = 3
while health > 0:
    print("still going...")   # health never changes → forever!
```

:::warning
Every `while` loop needs an "off switch." Before you run one, ask: *what line makes this condition eventually false?* In the safe example above it's `health = health - 1`. If you can't point to that line, you've got an infinite loop. If one does run away in the [Code Lab](#code), stop it and check that the condition variable actually changes inside the body.
:::

## break and continue

Two keywords give you finer control inside any loop. `break` **exits the loop immediately**, and `continue` **skips to the next pass**.

```python
# break: stop as soon as we find the boss
for enemy in ["goblin", "orc", "BOSS", "troll"]:
    if enemy == "BOSS":
        print("Boss found — stop spawning!")
        break
    print(f"Spawned {enemy}")
# Spawned goblin
# Spawned orc
# Boss found — stop spawning!
```

```python
# continue: skip the fallen enemies, act on the rest
healths = [10, 0, 25, 0, 8]
for hp in healths:
    if hp == 0:
        continue          # skip this one, go to the next
    print(f"Enemy still alive with {hp} HP")
# Enemy still alive with 10 HP
# Enemy still alive with 25 HP
# Enemy still alive with 8 HP
```

`break` ends the whole loop; `continue` only abandons the current pass and moves on. A common pattern is an intentional `while True:` loop that you exit with `break` when the right moment comes.

:::predict
```python
total = 0
for n in range(1, 6):
    if n == 3:
        continue
    total = total + n
print(total)
```
- 15
- 12 *
- 6
E: The loop would add 1+2+3+4+5 = 15, but `continue` skips `n == 3`, so 3 is left out. 15 − 3 = 12.
:::

## The accumulator pattern: a running total

One of the most useful loop patterns is the **accumulator**: you start a variable at a base value, then update it a little on every pass. Summing scores is the textbook example.

```python
scores = [120, 85, 200, 45]
total = 0                    # start the accumulator
for score in scores:
    total = total + score    # add each score in turn
print(f"Total score: {total}")   # Total score: 450
```

The pattern is always the same three steps: **initialize** a variable before the loop, **update** it inside the loop, **use** it after. It works for more than sums — count how many enemies remain, find the highest score, build up a string. Once you recognize the shape, you'll spot it everywhere.

```python
# accumulate a count instead of a sum
enemies = [10, 0, 25, 0, 0, 8]
alive = 0
for hp in enemies:
    if hp > 0:
        alive = alive + 1
print(f"{alive} enemies still standing")   # 3 enemies still standing
```

:::quiz
Q: Which loop should you choose when you want to keep asking the player for input until they type "quit" — and you have no idea how many tries that will take?
- a `for` loop over `range(100)`
- a `while` loop *
- an `enumerate()` loop
E: You don't know the number of repetitions ahead of time, only the stopping condition ("until they type quit"). That's exactly when a `while` loop fits — it repeats until the condition changes.
:::

## What's next

You can now repeat work automatically: `for` loops to march through collections and ranges, `while` loops to run until a condition flips, `break` and `continue` to steer mid-loop, and the accumulator pattern to build up totals. Next up, **Comprehensions**, where you'll learn to write certain loops in a single tidy line — transforming and filtering collections faster than ever. Build the accumulator examples in the [Code Lab](#code), then loop back with Pixel.

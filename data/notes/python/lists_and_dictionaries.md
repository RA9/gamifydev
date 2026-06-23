# Lists and Dictionaries

So far you've worked mostly with single values — one number, one string, one truth. But real programs deal with *collections*: a shopping list, a class of students, a phone book, a game's high scores. Python gives you two superb tools for holding many values at once: the **list** and the **dictionary**. These two structures show up in nearly every Python program ever written. Get comfortable here and you've unlocked a huge amount of power.

This is a lesson made for the [Code Lab](#code). Build the examples, poke at them, break them on purpose.

## Lists: ordered collections

A **list** holds an ordered sequence of values inside square brackets. The items can be anything — numbers, strings, even other lists.

```python
fruits = ["apple", "banana", "cherry"]
print(fruits)          # ['apple', 'banana', 'cherry']
print(len(fruits))     # 3 — len() tells you how many items
```

Each item has a position called an **index**, and — this surprises everyone at first — counting starts at **0**. So the first item is index `0`, the second is `1`, and so on.

```python
print(fruits[0])    # apple   — the first item
print(fruits[1])    # banana
print(fruits[-1])   # cherry  — negative indices count from the end
```

That `[-1]` trick is wonderfully handy: it always gives you the *last* item, no matter how long the list is. `[-2]` gives the second-to-last, and so on.

:::analogy
A list is like a row of numbered lockers. The lockers are in a fixed order, and each has a number on it — but the first locker is number 0, not 1. To get what's inside, you give the locker number: `fruits[0]`. Negative numbers are like counting backward from the far end of the hallway.
:::

## Slicing and growing lists

You can grab a *range* of items with **slicing**, using `start:stop`. As with `range`, the stop index is *not* included.

```python
numbers = [10, 20, 30, 40, 50]
print(numbers[1:3])   # [20, 30]   — items at index 1 and 2 (3 is excluded)
print(numbers[:2])    # [10, 20]   — from the start up to index 2
print(numbers[2:])    # [30, 40, 50] — from index 2 to the end
```

Lists are **mutable**, meaning you can change them after creating them. You can add to the end with `.append()`, change an item by its index, and remove items too:

```python
tasks = ["wake up", "code"]
tasks.append("celebrate")     # add to the end
tasks[0] = "wake up early"    # change the first item
print(tasks)                  # ['wake up early', 'code', 'celebrate']
```

:::key
Lists are **ordered** (items keep their position) and **mutable** (you can change them in place). Indexing starts at 0, slicing excludes the stop index, and `.append()` adds to the end. These four facts cover most of what you'll do with lists.
:::

To visit every item, loop over the list directly with a `for` loop:

```python
for task in tasks:
    print(f"TODO: {task}")
```

## Dictionaries: labeled collections

A list is great when *position* is what matters. But often you want to look things up by a meaningful **name** rather than a number — a price by product name, a score by player name, a definition by word. That's a **dictionary**: a collection of `key → value` pairs, written with curly braces.

```python
ages = {"Ada": 36, "Linus": 54, "Grace": 85}
print(ages["Ada"])     # 36 — look up by key, not by position
```

Instead of asking "what's at position 0?", you ask "what's the value for the key `"Ada"`?" Keys are usually strings (though they can be numbers too), and each key maps to exactly one value.

:::analogy
A dictionary is exactly like a real dictionary or a contacts app. You don't flip to "page 0" to find someone's number — you look them up by *name*. The name is the key; the number is the value. That's why dictionaries are perfect whenever your data has natural labels.
:::

## Reading, adding, and updating

You add a new pair or update an existing one with the same simple syntax — assign to a key:

```python
ages = {"Ada": 36}
ages["Grace"] = 85       # add a new key
ages["Ada"] = 37         # update an existing key
print(ages)              # {'Ada': 37, 'Grace': 85}
```

To loop through a dictionary, `.items()` hands you each key and value together — clean and readable:

```python
for name, age in ages.items():
    print(f"{name} is {age} years old.")
```

:::warning
Looking up a key that doesn't exist raises a `KeyError` and stops your program:

```python
print(ages["Pixel"])   # KeyError: 'Pixel'
```

When you're not sure a key exists, use `.get()`, which returns `None` (or a default you choose) instead of crashing:

```python
print(ages.get("Pixel"))        # None
print(ages.get("Pixel", 0))     # 0 — your chosen fallback
```
:::

## Choosing the right tool

So when do you reach for a list versus a dictionary? The question to ask is: *how will I look things up?*

- Use a **list** when order matters and you access items by position, or when you simply have a sequence to process from start to finish — a queue of tasks, a series of measurements, the lines of a file.
- Use a **dictionary** when each item has a natural label and you'll look it up by that label — a username to a profile, a product to a price, a setting to its value.

Here's a small worked example that uses both together. We have a list of orders, and a dictionary mapping each item to its price. We loop through the orders, look up each price, and total it up:

```python
prices = {"coffee": 4, "tea": 3, "cake": 6}
orders = ["coffee", "cake", "coffee", "tea"]

total = 0
for item in orders:
    total = total + prices[item]   # look up the price by name

print(f"You ordered {len(orders)} items.")
print(f"Total: ${total}")
# You ordered 4 items.
# Total: $17
```

See how naturally they cooperate? The **list** holds the *sequence* of what was ordered (order matters, duplicates allowed), and the **dictionary** holds the *lookup table* of prices (each item labeled with its cost). Recognizing which structure fits which job is a real mark of growing as a programmer — and you just did it.

:::tip
A quick mental test: if you'd describe your data as "a list of things in order," use a list. If you'd describe it as "a thing *for* each name/label," use a dictionary. Lists answer "what's next?"; dictionaries answer "what's the value for this key?"
:::

## Talk about it

> Picture the information in your phone. Your photo gallery is a *list* — an ordered sequence you scroll through. Your contacts are a *dictionary* — you look people up by name, not by position. Find two more examples from your own life: one that's naturally a list, and one that's naturally a dictionary. Explaining *why* each fits is how the distinction becomes second nature.

## What's next

You can now hold and organize collections of data — and combine lists and dictionaries to model surprisingly real situations. That's the backbone of practically every useful program. In the final lesson, **Next Steps with Python**, we'll look at where to go from here: shortcuts that make this code shorter, tools that extend Python's powers, and project ideas to make it all stick. Run that café example in the [Code Lab](#code), then come celebrate with Pixel.

:::quiz
Q: In the list `colors = ["red", "green", "blue"]`, what is `colors[1]`?
- green *
- red
- blue
E: Indexing starts at 0, so index 1 is the second item, "green".
:::

:::predict
```python
data = {"a": 1, "b": 2}
data["c"] = 3
print(len(data))
```
- 3 *
- 2
- error
E: Assigning to a new key adds it, so the dictionary now has three pairs and `len` is 3.
:::

:::fill
Q: Complete the method that loops over both keys and values of a dictionary.
`for key, value in ages.___():`
- items *
- keys
- values
E: `.items()` yields each (key, value) pair, perfect for unpacking into two loop variables.
:::

:::match
Q: Match each task to the better data structure.
- A queue of songs to play in order | list
- Looking up a price by product name | dictionary
- Storing the lines of a poem in sequence | list
E: Ordered sequences fit lists; labeled lookups fit dictionaries.
:::

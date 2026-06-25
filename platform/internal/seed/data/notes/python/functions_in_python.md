# Functions in Python

As your programs grow, you'll notice yourself writing the same few lines over and over, or building one giant wall of code that's hard to follow. **Functions** are the cure. A function is a named, reusable bundle of instructions — you write it once and call it whenever you need it. They turn tangled code into clear, labeled tools. This is one of the most important ideas in all of programming, so let's take our time.

As always, the [Code Lab](#code) is your workbench. Try every example.

## Defining and calling a function

You create a function with the `def` keyword, give it a name, and indent the body underneath. Then you *call* it by writing its name with parentheses.

```python
def greet():
    print("Hello there!")

greet()    # calls the function — prints "Hello there!"
greet()    # call it again — prints it again
```

Defining a function doesn't run it; it just teaches Python what `greet` *means*. The code inside only runs when you call it. That separation is the whole point: define once, use many times.

:::analogy
A function is like a recipe card in a box. Writing the card (`def`) doesn't cook anything — it just records the steps. Cooking happens when you pull the card and follow it (the call). And once the card exists, anyone can make the dish again and again without rewriting the recipe.
:::

## Parameters and return values

Functions get powerful when they accept **inputs** and hand back **outputs**. Inputs are called *parameters*; you list them in the parentheses. To send a result back, use `return`.

```python
def add(a, b):
    return a + b

total = add(3, 4)
print(total)        # 7
```

Here `a` and `b` are parameters — placeholders for whatever values you pass in. When you call `add(3, 4)`, `a` becomes `3` and `b` becomes `4`. The `return` statement sends the result back to wherever the function was called, so `total` ends up holding `7`.

There's an important distinction between *printing* and *returning*. `print` shows something on screen; `return` hands a value back to your program so you can keep using it. A function that returns a value is far more flexible:

```python
def square(n):
    return n * n

print(square(5))           # 25
print(square(5) + square(2))  # 25 + 4 = 29 — because square returns a usable value
```

:::key
`print` is for *humans* to read; `return` is for the *rest of your program* to use. If a function only prints, you can't do math with its result. Prefer `return` when the value will be used elsewhere.
:::

## Default arguments

You can give a parameter a **default value**, used when the caller doesn't supply one. This makes functions flexible without forcing callers to spell out every detail.

```python
def greet(name, greeting="Hello"):
    return f"{greeting}, {name}!"

print(greet("Ada"))                 # Hello, Ada!
print(greet("Ada", "Welcome"))      # Welcome, Ada!
```

Because `greeting` has a default, you can leave it out and get `"Hello"` for free — but you can still override it when you want something different. Parameters with defaults must come *after* parameters without them.

## Returning more than one value

Sometimes a function naturally produces several results. Python lets you return them together as a **tuple** — just separate them with commas — and unpack them on the other side.

```python
def min_and_max(numbers):
    return min(numbers), max(numbers)

low, high = min_and_max([4, 1, 9, 2])
print(low)    # 1
print(high)   # 9
```

The function returns two values at once; `low, high = ...` neatly splits them into two variables. This is a clean, very Pythonic pattern you'll see often.

:::tip
A *tuple* is just an ordered group of values, like `(1, 9)`. When a function returns several things, it's actually returning one tuple — and `low, high = ...` unpacks it. You don't even need the parentheses when returning: `return a, b` is enough.
:::

## Scope: where names live

When you create a variable *inside* a function, it's **local** — it exists only while the function runs and is invisible outside it. Variables created at the top level of your file are **global**. This separation keeps functions tidy and prevents them from accidentally clobbering each other's data.

```python
message = "global"          # global variable

def show():
    message = "local"       # a separate, local variable
    print(message)          # local

show()
print(message)              # global — the outer one is untouched
```

When Python looks up a name, it searches outward in layers — **Local**, then any **Enclosing** function, then **Global**, then **Built-in**. (You'll hear this called the *LEGB rule*; you don't need to memorize it, just trust that Python starts close and widens the search.) The practical lesson: a function can *read* outer variables, but assigning inside a function creates a new local one by default.

This is exactly why functions are so valuable: they're little sealed rooms. What happens inside mostly stays inside, so you can reason about one function without worrying about the rest of your program.

## Why decomposition matters

The deepest reason to use functions isn't reuse — it's *thinking*. Breaking a big problem into small, named functions (this is called **decomposition**) lets you tackle one piece at a time and read your program like a table of contents.

```python
def get_price(item):
    prices = {"coffee": 4, "tea": 3}
    return prices[item]

def apply_tax(amount):
    return amount * 1.1

def checkout(item):
    price = get_price(item)
    return apply_tax(price)

print(checkout("coffee"))   # 4.4
```

Each function does one clear job, and `checkout` reads almost like a sentence describing the process. Good function names are like good headings: they tell you what's happening without making you read every detail.

A docstring — a triple-quoted string right under the `def` — documents what a function does. Future-you will be grateful:

```python
def apply_tax(amount):
    """Return the amount with 10% tax added."""
    return amount * 1.1
```

:::warning
Avoid using a *mutable* value (like a list `[]` or dict `{}`) as a default argument. It's created **once** and shared across all calls, leading to surprising bugs:

```python
def add_item(item, basket=[]):   # risky!
    basket.append(item)
    return basket

print(add_item("apple"))   # ['apple']
print(add_item("banana"))  # ['apple', 'banana'] — the list was reused!
```

The fix is to default to `None` and create a fresh list inside:

```python
def add_item(item, basket=None):
    if basket is None:
        basket = []
    basket.append(item)
    return basket
```
:::

## Talk about it

> Think about how you'd explain making a sandwich to someone who'd never done it. You'd probably name sub-steps: "toast the bread," "add the filling," "cut it in half." Each of those is a function. Pick a familiar multi-step task and break it into three or four named functions. What inputs would each need, and what would each hand back? That instinct to break things into named pieces is the heart of programming well.

## What's next

You can now package logic into reusable, well-named tools — a skill that will shape how you write every program from here on. Next, in **Lists and Dictionaries**, you'll learn Python's two workhorse data structures for holding *collections* of values. Functions plus data structures is where Python really starts to sing. Warm up in the [Code Lab](#code), then carry on.

:::quiz
Q: What's the difference between `print` and `return` inside a function?
- `return` hands a value back to the program to reuse; `print` only shows it on screen *
- They do exactly the same thing
- `print` hands a value back; `return` shows it on screen
E: `return` produces a value the rest of your program can use, while `print` merely displays text.
:::

:::predict
```python
def add(a, b=10):
    return a + b

print(add(5))
```
- 15 *
- 5
- error
E: `b` defaults to 10 when not supplied, so `add(5)` returns 5 + 10 = 15.
:::

:::reorder
Order these lines to define a `double` function and then call it.
- def double(n):
-     return n * 2
- print(double(4))
E: Define the function and its body first, then call it afterward.
:::

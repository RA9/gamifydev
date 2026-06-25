# Python Basics

You've met Python — now let's start a real conversation with it. In this lesson you'll learn how to *remember* values, work with different kinds of data, do arithmetic, and get your program talking with you. These are the building blocks under every Python program ever written, no matter how fancy. Master them and everything else gets easier.

Open the [Code Lab](#code) alongside this lesson. Reading is good; running is better.

## Variables: giving values a name

A **variable** is a name that points to a value. Instead of juggling raw numbers and text in your head, you give them labels so you can refer to them later.

```python
name = "Ada"
age = 36
print(name)
print(age)
```

The `=` here does *not* mean "equals" the way it does in math. It means **assignment**: "take the value on the right and store it under the name on the left." After `name = "Ada"`, the word `name` stands in for `"Ada"` everywhere you use it.

:::analogy
Think of a variable as a labeled box. You write a label on the box (`age`) and put something inside it (`36`). Later, when you say `age`, Python opens the box and hands you what's inside. You can swap the contents anytime — that's why it's called a *variable*: its value can vary.
:::

You can change what a variable holds simply by assigning again:

```python
score = 10
score = 25
print(score)   # 25 — the new value replaced the old one
```

## Dynamic typing and core types

Every value in Python has a **type** — a category that describes what kind of data it is. The four you'll meet first are:

- **`int`** — whole numbers, like `7`, `0`, or `-42`.
- **`float`** — numbers with a decimal point, like `3.14` or `-0.5`.
- **`str`** — text, called a *string*, written in quotes: `"hello"` or `'Pixel'`.
- **`bool`** — a truth value, either `True` or `False`.

Python is **dynamically typed**, which means *you* don't have to declare a variable's type — Python figures it out from the value you give it. Assign a number and it's a number; assign text and it's a string. You can even point the same variable at a different type later (though it's good style not to do that carelessly).

```python
x = 5          # int
x = 5.0        # now a float
x = "five"     # now a str
```

You can ask Python what type something is with the built-in `type()` function:

```python
print(type(5))        # <class 'int'>
print(type(3.14))     # <class 'float'>
print(type("hi"))     # <class 'str'>
print(type(True))     # <class 'bool'>
```

:::key
A *value* has a type; a *variable* just points at whatever value you assign. Because Python is dynamically typed, you never write the type yourself — but knowing the type still matters, because what you can *do* with a value depends on it.
:::

## Printing and f-strings

You've already met `print()` — it shows output. But the real magic comes when you mix text and variables together. The cleanest way is an **f-string**: put an `f` before the opening quote, then drop variables inside `{ }`.

```python
name = "Ada"
age = 36
print(f"{name} is {age} years old.")
# Ada is 36 years old.
```

Python replaces each `{...}` with the value inside it. You can even put expressions in the braces:

```python
price = 4
print(f"Two coffees cost {price * 2} dollars.")
# Two coffees cost 8 dollars.
```

F-strings are the friendliest way to build messages, so we'll use them constantly. They read almost like a sentence with blanks to fill in.

:::tip
Forgetting the `f` is a classic beginner slip. `print("Hi {name}")` prints the literal text `Hi {name}` — braces and all. With the `f`, `print(f"Hi {name}")` prints `Hi Ada`. If your variable name shows up literally in the output, check for the missing `f`.
:::

## Talking back: input and comments

So far your programs talk *to* you. With `input()`, they can listen too. `input()` pauses the program, waits for the user to type something, and hands back what they typed — always as a **string**.

```python
name = input("What is your name? ")
print(f"Nice to meet you, {name}!")
```

Because `input()` always returns text, if you want a number you must convert it. `int(...)` turns a string into an integer; `float(...)` turns it into a decimal number:

```python
age_text = input("How old are you? ")
age = int(age_text)
print(f"Next year you'll be {age + 1}.")
```

As your programs grow, you'll want to leave notes for yourself. A **comment** starts with `#`; Python ignores everything after it on that line. Comments explain *why* you did something — they're messages to future-you (and to teammates).

```python
# Convert the user's text into a number so we can do math with it
age = int(input("How old are you? "))
```

:::warning
Because `input()` gives you a string, `input("Number: ") + 1` will fail or behave strangely — you can't add a number to text. Always convert with `int()` or `float()` first when you need to do arithmetic.
:::

## Arithmetic that surprises beginners

Python does ordinary math exactly as you'd expect: `+`, `-`, and `*` add, subtract, and multiply. But three operators deserve special attention.

```python
print(7 + 3)    # 10
print(7 - 3)    # 4
print(7 * 3)    # 21
print(7 / 3)    # 2.3333333333333335  — regular division, always a float
print(7 // 3)   # 2   — floor division, drops the remainder
print(7 % 3)    # 1   — modulo, the remainder
print(2 ** 4)   # 16  — exponent, 2 to the power of 4
```

Notice that `/` always gives a **float**, even when the numbers divide evenly (`6 / 2` is `3.0`, not `3`). When you want a whole-number result — like "how many full boxes fit?" — use `//`, *floor division*. And `%`, the *modulo* operator, gives the remainder, which is wonderfully handy for questions like "is this number even?" (`n % 2 == 0`).

The `**` operator raises a number to a power: `2 ** 10` is `1024`. These three — `//`, `%`, and `**` — trip up almost every newcomer at least once, so try them yourself in the [Code Lab](#code) until they feel familiar.

## Talk about it

> Variables let you give meaningful names to values. Imagine describing yourself to a computer using only variables — what would you store? Pick three things about your day (a count, a price, a name, a yes/no fact) and decide which Python type each one would be: `int`, `float`, `str`, or `bool`. Naming the type is the first step to thinking like a programmer.

## What's next

You can now store data, label it, do math with it, and exchange messages with the user. Next, in **Control Flow in Python**, you'll teach your programs to make *decisions* and *repeat* work — the moment code stops being a straight line and starts feeling alive. Bring your new variables along; you're about to make them do interesting things.

:::quiz
Q: What does `type(3.0)` report?
- <class 'float'> *
- <class 'int'>
- <class 'str'>
E: A number written with a decimal point is a float, so `type(3.0)` is `<class 'float'>`.
:::

:::predict
```python
print(7 // 2)
print(7 % 2)
```
- 3 then 1 *
- 3.5 then 0
- 1 then 3
E: `7 // 2` floors the division to 3; `7 % 2` is the remainder, 1.
:::

:::fill
Q: Complete the f-string so it prints "Hi Ada".
`name = "Ada"`
`print(___"Hi {name}")`
- f *
- s
- p
E: An f-string needs the `f` prefix so Python substitutes the value of `name` into the `{ }`.
:::

:::match
Q: Match each value to its Python type.
- `42` | int
- `"hi"` | str
- `True` | bool
- `3.14` | float
E: Whole number → int, quoted text → str, True/False → bool, decimal → float.
:::

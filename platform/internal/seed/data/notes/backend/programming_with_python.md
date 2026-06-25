# Programming with Python

To build a backend, you need a language to write the logic in. **Python** is the most popular choice for beginners — and for good reason: it reads almost like plain English, so you can focus on *ideas* instead of fighting syntax.

By the end of this lesson you'll recognise Python's core building blocks and read simple Python code with confidence.

## Why Python?

Python is clean, readable, and everywhere — web backends, data science, automation, AI. The same fundamentals you met in *Intro to Programming* (variables, conditions, loops, functions) all appear here, just with friendlier syntax.

:::analogy
If some languages are like assembling furniture with a cryptic manual, Python is the one where the instructions are written in clear sentences. Same furniture — far less frustration.
:::

## Variables and types

No special keywords needed — just name a box and put a value in it:

```python
name = "Ada"        # a string (text)
age = 25            # an integer (number)
is_student = True   # a boolean (True/False)

print("Hello, " + name)
```

`print(...)` displays a value — your everyday tool for seeing what your code is doing.

## Making decisions

Python uses **indentation** (spaces) to group code, instead of curly braces:

```python
age = 18

if age >= 18:
    print("You can vote")
else:
    print("Not yet")
```

Notice there are no `{ }` — the indented lines belong to the `if`. Clean and readable.

## Loops

Repeat work without copy-paste:

```python
for i in range(3):
    print("Hello number", i)

# Hello number 0
# Hello number 1
# Hello number 2
```

:::quiz
Q: How does Python know which lines belong inside an `if` or a loop?
- Curly braces `{ }`
- Indentation (the spaces at the start of the line) *
- Semicolons
E: Python uses indentation to define blocks. Lines indented under an `if`, `for`, or `def` belong to it.
:::

## Functions

Package reusable steps with `def`:

```python
def greet(name):
    return "Hi, " + name + "!"

print(greet("Ada"))    # Hi, Ada!
print(greet("Grace"))  # Hi, Grace!
```

`def` defines a function; `return` hands a value back to the caller. This is exactly the same idea as JavaScript functions — just tidier.

## Lists and dictionaries

Two collections you'll use constantly on the backend:

```python
todos = ["Learn Python", "Build an API"]   # a list (ordered items)

user = { "name": "Ada", "age": 25 }          # a dictionary (key → value)
print(user["name"])                          # Ada
```

**Dictionaries** are especially important — they look just like the JSON data your API will send to the frontend.

:::quiz
Q: Which Python collection maps keys to values, like `"name" -> "Ada"`?
- A list
- A dictionary *
- A loop
E: A dictionary stores key–value pairs. Lists store ordered items; dictionaries store labelled ones (and they map neatly onto JSON).
:::

:::key
Python is a readable language built from the same fundamentals you already know. Key tools: **variables**, **if/else**, **for** loops, **functions** with `def`/`return`, and the **list** and **dictionary** collections.
:::

## Talk about it

Explain out loud:

> "How does Python group code without curly braces, and what's the difference between a list and a dictionary?"

When you take the quiz for this lesson, you'll test these Python ideas directly.

## What's next

You can write logic. Next, **Working with Data and Databases** — where you'll learn to store and retrieve information that outlives a single run of your program.

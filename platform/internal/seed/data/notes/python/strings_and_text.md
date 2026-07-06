# Strings and Text

Almost every program handles text: player names, chat messages, item descriptions, error messages, menu labels. In Python, a piece of text is called a **string**. This lesson shows you how to create strings, glue them together, measure them, pull them apart, reshape them, and drop values right into the middle of them.

## Creating Strings

A string is text wrapped in quotes. Python lets you use **single quotes** or **double quotes** — they work exactly the same, so pick one and stay consistent.

```python
player = "Nova"
weapon = 'Frost Blade'
```

Why two choices? Because it lets you include the *other* quote inside without any fuss:

```python
line1 = "It's a trap!"          # apostrophe is fine inside double quotes
line2 = 'She said "go now".'    # double quotes fine inside single quotes
```

For text that spans multiple lines, use **triple quotes** (three of either kind). Everything between them — including the line breaks — becomes part of the string.

```python
intro = """Welcome, adventurer.
The gate is open.
Your quest begins now."""
```

Triple-quoted strings are perfect for long messages, help text, or dialogue that runs across several lines.

:::tip
Consistency matters more than which quote you pick. Many Python projects default to double quotes. Choose a style and your code will read cleanly.
:::

## Joining Strings: Concatenation

Gluing strings together end-to-end is called **concatenation**, and you do it with `+`.

```python
first = "Frost"
last = "Blade"
name = first + last
print(name)          # → FrostBlade
```

Notice there's no space unless you add one — Python joins exactly what you give it:

```python
greeting = "Welcome, " + "Nova" + "!"
print(greeting)      # → Welcome, Nova!
```

You can also repeat a string with `*`, which is great for simple visuals:

```python
print("=" * 20)      # → ====================
print("go! " * 3)    # → go! go! go! 
```

:::warning
`+` only joins string to string. Try to add a string and a number directly — like `"Level " + 5` — and Python raises a `TypeError`. You must convert the number first with `str(5)`, or better, use an f-string (coming up soon).
:::

## Measuring and Indexing

To find out how many characters a string contains, use `len()`:

```python
name = "Nova"
print(len(name))     # → 4
```

Each character sits at a numbered position called an **index**, and Python counts from **0**, not 1. So the first character is at index `0`.

```python
name = "Nova"
print(name[0])       # → N   (first character)
print(name[1])       # → o   (second character)
print(name[3])       # → a   (fourth character)
```

You can also count from the *end* using negative indexes. `-1` is the last character, `-2` the second-to-last, and so on — very handy when you don't know the length.

```python
name = "Nova"
print(name[-1])      # → a   (last character)
print(name[-2])      # → v
```

:::key
Python indexes start at **0**. For a string of length `n`, the valid positions run from `0` to `n - 1`. Reaching past the end (like `name[4]` on a 4-letter word) raises an `IndexError`.
:::

## Slicing: Grabbing a Piece

A **slice** pulls out a range of characters using `s[start:stop]`. The slice includes the `start` position but **stops before** `stop`.

```python
message = "GameOver"
print(message[0:4])   # → Game   (positions 0,1,2,3 — not 4)
print(message[4:8])   # → Over
```

Leave out a number and Python fills in a sensible default — the beginning or the end:

```python
tag = "player_42"
print(tag[:6])        # → player   (from start up to position 6)
print(tag[7:])        # → 42       (from position 7 to the end)
```

Slicing is your everyday tool for chopping usernames, trimming prefixes, or reading part of a code.

## Strings Are Immutable

Here's an important idea: strings are **immutable**, meaning you can't change a character *inside* an existing string. Try it and Python refuses:

```python
name = "Nova"
# name[0] = "L"       # ❌ TypeError: strings can't be modified in place
```

Instead, you build a *new* string and reassign it:

```python
name = "Nova"
name = "L" + name[1:]   # make a new string, then point `name` at it
print(name)             # → Lova
```

This feels strict at first, but it makes strings safe and predictable — nobody can quietly alter your text from under you. Every string method you're about to meet works this way: it **returns a new string** and leaves the original untouched.

## Common String Methods

Methods are actions you call on a string with a dot: `text.method()`. Here are the ones you'll use constantly.

Change the case with `.upper()` and `.lower()`:

```python
name = "Nova"
print(name.upper())         # → NOVA
print("SHOUT".lower())      # → shout
```

Trim stray whitespace from the ends with `.strip()` — essential for cleaning user input:

```python
typed = "  nova  "
print(typed.strip())        # → nova   (spaces gone)
```

Swap text with `.replace(old, new)`:

```python
chat = "gg wp"
print(chat.replace("gg", "good game"))   # → good game wp
```

Break a string into a list of pieces with `.split()`, and join a list back into a string with `.join()`:

```python
csv = "sword,shield,potion"
items = csv.split(",")       # → ['sword', 'shield', 'potion']
print(items[1])              # → shield

joined = " > ".join(items)   # glue list items with " > " between them
print(joined)                # → sword > shield > potion
```

Test how a string begins or ends, and locate text inside it:

```python
cmd = "/mute Nova"
print(cmd.startswith("/"))       # → True   (it's a command)
print("report.txt".endswith(".txt"))   # → True
print("hello world".find("world"))     # → 6  (index where it starts; -1 if absent)
```

:::tip
Because methods return a *new* string, you can chain them: `"  Nova  ".strip().upper()` gives `"NOVA"`. Read left to right — strip first, then uppercase the result.
:::

## f-strings: The Modern Way to Format

Concatenating with `+` gets clumsy fast, and it breaks the moment a number sneaks in. The clean, modern solution is the **f-string**. Put an `f` right before the opening quote, then drop any value inside `{ }`.

```python
player = "Nova"
level = 7

# The old, clunky way:
print("Player " + player + " is level " + str(level))

# The f-string way — much nicer:
print(f"Player {player} is level {level}")
# → Player Nova is level 7
```

Anything Python can evaluate can go inside the braces, including math:

```python
hp = 80
max_hp = 100
print(f"Health: {hp}/{max_hp} ({hp / max_hp * 100}%)")
# → Health: 80/100 (80.0%)
```

Notice you don't need `str()` around the numbers — the f-string converts them for you. This is why f-strings are the go-to way to build messages in modern Python.

:::analogy
An f-string is a fill-in-the-blank sentence. You write the sentence once with blanks — `f"Player {player} is level {level}"` — and Python drops the real values into the blanks. No juggling `+` signs or quotes.
:::

## Escape Sequences

Some characters are hard to type directly inside a string — like a newline or a tab. For those, you use an **escape sequence**: a backslash `\` followed by a letter that stands for the special character.

The two you'll meet most are `\n` (newline — starts a new line) and `\t` (tab — inserts a tab stop):

```python
print("Line one\nLine two")
# → Line one
# → Line two

print("Name:\tNova")
# → Name:	Nova
```

And when you literally need a backslash in your text, you escape it with another backslash:

```python
print("Path: C:\\Games\\Save")   # → Path: C:\Games\Save
```

:::warning
`\n` is a *single* newline character, not the two letters "backslash n." That's why `len("a\nb")` is `3`, not `4` — the `\n` counts as one character.
:::

## Practice

:::predict
Q: What does this print?
```python
name = "Nova"
print(name[1:3])
```
- No
- ov *
- ova
- Nov
E: Slicing `[1:3]` starts at index 1 (`o`) and stops *before* index 3, so it includes positions 1 and 2 — `o` and `v` — giving `ov`.
:::

:::predict
Q: What does this print?
```python
score = 42
print(f"Score: {score}")
```
- Score: {score}
- Score: 42 *
- Score: score
- f"Score: 42"
E: An f-string replaces `{score}` with the variable's value, `42`, producing `Score: 42`. The `f` prefix is what enables the substitution.
:::

## Recap

- A **string** is text in single, double, or triple quotes; triple quotes span multiple lines.
- Join strings with `+` (concatenation) and repeat them with `*`.
- `len()` counts characters; indexing starts at **0**, and negative indexes count from the end (`s[-1]` is last).
- **Slicing** `s[start:stop]` grabs a range that includes `start` but stops before `stop`; omit either side for a default.
- Strings are **immutable** — methods return a *new* string rather than changing the original.
- Handy methods: `.upper()`, `.lower()`, `.strip()`, `.replace()`, `.split()`, `.join()`, `.startswith()`, `.endswith()`, `.find()`.
- **f-strings** (`f"...{value}..."`) are the modern, readable way to drop values into text.
- **Escape sequences** like `\n` (newline) and `\t` (tab) put special characters into strings.

**Next up:** Numbers, Math, and Booleans — doing the arithmetic and true/false logic behind damage, XP, and health.

# Setting Up Python

Welcome to Python — one of the friendliest programming languages ever made, and a favorite for games, tools, websites, data, and automation. This lesson gets you oriented: how Python actually runs your code, the two ways to run it, how to print things out, and how to leave notes for yourself. No prior experience needed. If you can read this sentence, you can learn Python.

## What Python Is and How It Runs Your Code

Python is a **programming language** — a way of writing instructions that a computer can follow. But your computer's processor doesn't natively understand Python. It understands its own low-level machine code. So something has to sit in the middle and translate.

That something is the **Python interpreter**. It reads your Python code one line at a time, figures out what each line means, and carries out the instruction right away. The interpreter is a program, and running your code really means "handing your file to the interpreter and letting it do the work."

:::analogy
Think of the interpreter as a helpful translator standing between you and the computer. You speak Python; the computer speaks machine code. You hand the translator a sentence, and it immediately turns around and tells the computer exactly what to do — one sentence at a time, top to bottom.
:::

Because Python runs line by line, you get results fast, and errors point you to the exact line that tripped up. That top-to-bottom, one-line-at-a-time behavior is a big part of why Python feels so approachable for beginners.

## Two Ways to Run Python

There are two everyday ways to run Python code, and it helps to know both from day one.

The first is the **REPL** — an interactive prompt where you type one line, press Enter, and see the result *immediately*. REPL stands for Read–Eval–Print Loop: it **R**eads what you type, **E**valuates it, **P**rints the result, and **L**oops back to wait for more. It's like a conversation with Python.

```python
>>> 2 + 2
4
>>> "hi" + "!"
'hi!'
>>> name = "Nova"
>>> name
'Nova'
```

Those `>>>` marks are the prompt Python shows you; you don't type them. The REPL is perfect for quick experiments — "what does this do?" — and you'll reach for it constantly.

The second way is running a **script**: a saved file of Python code, ending in `.py`, that runs from top to bottom all at once. This is how you build real, reusable programs.

```python
# game_intro.py
print("Loading dungeon...")
print("Spawning player...")
print("Ready. Good luck!")
```

The REPL is your sketchpad. A `.py` script is the finished drawing you keep and share.

:::key
REPL = one line at a time, instant feedback, great for trying things. Script (`.py` file) = a whole program saved and run start-to-finish. You'll use both, often in the same session.
:::

## print(): Your Window Into a Program

Computers do a lot of work silently. Variables get set, math happens, decisions are made — but none of it shows up on screen unless you *ask*. The tool that asks is `print()`.

`print()` displays whatever you put between its parentheses.

```python
print("Hello, world!")     # → Hello, world!
print(42)                  # → 42
print(3 + 4)               # → 7
```

You can print several things at once by separating them with commas. Python puts a space between them for you:

```python
print("Score:", 10, "Lives:", 3)   # → Score: 10 Lives: 3
```

In the REPL, typing an expression shows its value automatically — but inside a `.py` script nothing appears unless you `print()` it. This is one of the very first things that trips up new coders.

:::warning
Running a script that has no `print()` calls can look like "nothing happened." The program probably *did* run — it just had no reason to show you anything. When in doubt, add a `print()` to peek at what's going on.
:::

`print()` is more than decoration. It's your primary way to *see* what your program is doing — checking a value, confirming a step ran, tracing a bug. Get comfortable sprinkling it everywhere while you learn.

## Comments: Notes to Yourself

A **comment** is a note in your code that Python completely ignores. Comments start with a `#`, and everything after the `#` on that line is skipped by the interpreter.

```python
# This whole line is a comment — Python skips it.

print("Level 1")   # You can also comment at the end of a line.

# health = 100      <- a line that's "commented out" won't run
```

Comments are for humans. Use them to explain *why* something is happening, not just to restate the obvious.

```python
# Bad: says what the code already says
score = score + 10   # add 10 to score

# Better: explains the reason
score = score + 10   # bonus for clearing the level without damage
```

There's no separate "block comment" syntax in Python — you just start each line with `#`. Most editors let you comment or uncomment a whole selection with a keyboard shortcut, which is handy when you want to temporarily disable some lines while testing.

:::tip
"Commenting out" code — putting a `#` in front of a line to stop it running — is a quick way to test an idea without deleting anything. Toggle it back on when you're done.
:::

## Running Python in This Course

Here's the good news for getting started: **in this course you write and run Python right here in the browser.** Every lab runs your code through an in-browser Python sandbox. There's nothing to install, no setup to wrestle with, no "it works on my machine" — you type Python, hit run, and see the output immediately. That means you can focus entirely on *learning the language* instead of configuring tools.

So when a lab asks you to print something or do a bit of math, you'll write real Python and it'll really run — all inside your browser tab.

## Running Python On Your Own Machine

Someday you'll want Python on your own computer — to build a project, automate a chore, or make a game you can share. Here's the short version so it's not a mystery.

First, you install Python from the official site, **python.org**. Grab the latest version (Python 3), run the installer, and — on Windows especially — check the box that says "Add Python to PATH." That box makes the `python` command available everywhere in your terminal.

Then you write a script and run it from a terminal (Command Prompt, PowerShell, or Terminal on Mac/Linux):

```python
# hello.py
print("Running from my very own machine!")
```

```text
$ python hello.py
Running from my very own machine!
```

That `$` is the terminal prompt; you type `python hello.py` after it. The interpreter reads your file, runs it top to bottom, and prints the output back to the terminal. On some systems the command is `python3` instead of `python` — if one doesn't work, try the other.

You can also just type `python` on its own to drop into the REPL, exactly like the interactive prompt we saw earlier. Type `exit()` to leave it.

:::analogy
Installing Python is like buying your own kitchen. In this course, we've handed you a fully stocked kitchen ready to cook in — you just start. When you install Python yourself, you're setting up your own stove and pantry at home so you can cook whenever you like, on whatever you like.
:::

## Practice

:::quiz
Q: What does the Python interpreter do?
- Turns your Python code into instructions the computer can run, one line at a time *
- Stores your files in the cloud
- Automatically fixes bugs in your code
- Draws the graphics for your game
E: The interpreter is the translator between your Python and the computer's own language. It reads and runs your code line by line.
:::

:::predict
Q: What does this script print?
```python
# starting up
print("Boot", "sequence", "complete")
```
- Boot sequence complete *
- Bootsequencecomplete
- "Boot" "sequence" "complete"
- # starting up
E: `print()` shows the three pieces separated by commas and inserts a space between each, so you get `Boot sequence complete`. The first line is a comment, so Python ignores it entirely.
:::

## Recap

- Python code is run by the **interpreter**, which translates and executes it one line at a time, top to bottom.
- The **REPL** is an interactive prompt for quick experiments — type a line, see the result instantly. Its prompt is `>>>`.
- A **script** is a saved `.py` file that runs start to finish; it's how you build real programs.
- `print()` is your window into a program — it's the main way to *see* values and confirm what your code is doing.
- **Comments** start with `#` and are ignored by Python; use them to explain *why*.
- In **this course**, Python runs in your browser via an in-browser sandbox — nothing to install.
- On **your own machine**, install Python from **python.org** and run a file with `python file.py` in a terminal.

**Next up:** Strings and Text — creating, combining, and reshaping the words your programs work with.

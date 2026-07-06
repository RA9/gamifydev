# Code Style, Linting and Type Checking

Working code is only half the job. Code that's *readable* and *checked* is what lets you (and your teammates) keep moving without tripping over yesterday's mistakes. This lesson covers the three tools that guard that quality: a formatter, a linter, and a type checker. You'll run all of them in your own terminal — this is a workflow lesson, not a browser lab.

## PEP 8 and Why Consistency Matters

Python has an official style guide called **PEP 8**. It spells out the small stuff: four spaces of indentation, `snake_case` for variable and function names, spaces around operators, a blank line or two between functions, and so on.

None of these rules change what the code *does*. So why bother? Because you read code far more often than you write it, and a consistent style lets your eyes skim past the formatting to focus on the logic. When a whole codebase follows one style, everything looks like it was written by one careful person — even with ten authors.

:::analogy
Style guides are like traffic conventions. There's no deep truth to "drive on the right" — it's an arbitrary choice. But *everyone agreeing* on it is what keeps the road flowing. PEP 8 is Python's shared side of the road.
:::

Here's the same logic written carelessly and then to PEP 8:

```python
# Before — technically valid, but noisy
def Add( x,y ):
  return x+y
result=Add(2,3)
```

```python
# After — PEP 8
def add(x, y):
    return x + y


result = add(2, 3)
```

Same behavior, far easier to read. The good part: you don't have to apply these rules by hand.

## Auto-Formatting with black

A **formatter** rewrites your code into a consistent style automatically. The most popular is **black**, famous for being *opinionated* — it makes the style choices for you so nobody has to argue about them.

Install it into your project's virtual environment, then run it:

```bash
pip install black
black .
```

The `.` means "format every Python file in this directory and below." Black rewrites the files in place. Run it, and the messy example above becomes the clean one — spacing, blank lines, and quotes all normalized. You simply stop thinking about formatting.

To *check* without changing files (useful for seeing what would change):

```bash
black --check .
```

The newer tool **ruff** also includes a formatter that's black-compatible and extremely fast:

```bash
ruff format .
```

:::tip
Configure your editor to run the formatter on every save. Then your code is always tidy without a single manual step — you type freely and it snaps into shape the instant you hit save.
:::

## Linting: Catching Bugs and Smells

Formatting is about *looks*. **Linting** is about *problems*. A linter reads your code without running it and flags likely bugs, unused variables, unreachable code, imports you never use, and other "smells."

**ruff** is the modern favorite — one fast tool that rolls together the checks of many older ones. (The long-standing **flake8** does the same job and you'll still see it in many projects.)

```bash
pip install ruff
ruff check .
```

Consider this snippet, which runs fine but hides two issues:

```python
import os          # imported but never used
import sys

def parse(line):
    parts = line.split(",")
    return sys.argv          # 'parts' is computed, then ignored
```

`ruff check .` catches both without you running anything:

```text
app.py:1:1: F401 `os` imported but unused
app.py:5:5: F841 local variable `parts` is assigned to but never used
```

Each finding has a code (`F401`, `F841`) you can look up. Many are auto-fixable — ruff will remove the dead import and let you decide about the unused variable:

```bash
ruff check --fix .
```

:::warning
A linter reasons about your code *statically* — by reading it, never by running it. So it catches whole categories of mistakes (typos in names, unused imports, obviously unreachable branches) but it can't know a value your program computes at runtime. Treat it as a sharp-eyed proofreader, not a test suite. You still need both.
:::

## Type Checking with mypy

You met type hints in an earlier lesson — annotations like `name: str` and `-> int` that document what a function expects and returns. Here's the crucial fact those hints hinge on:

:::key
Python does **not** enforce type hints at runtime. They are documentation the interpreter ignores while running. A separate *tool* reads them and checks that your code is consistent. That tool is a type checker, and the standard one is **mypy**.
:::

Install and run it against your package:

```bash
pip install mypy
mypy src/
```

Now watch it earn its keep. This function claims to take a number, but the caller passes text:

```python
def double(n: int) -> int:
    return n * 2

double("hello")   # a string, not an int
```

Run the code and Python happily returns `"hellohello"` — the bug slips through silently. But mypy reads the hints and catches it *before* you ever run anything:

```text
error: Argument 1 to "double" has incompatible type "str"; expected "int"  [arg-type]
```

That's the whole point: the hint `n: int` was a promise, and mypy holds your code to it. As projects grow, this turns a class of "it crashed in production" bugs into "the checker complained before I committed."

:::example
You rename a function to return `None` in one place but forget a caller that still does `result + 1`. The program might run for weeks before that path executes and crashes. mypy sees the mismatch — `None + int` — the moment you run it, turning a lurking runtime crash into an instant, located error message.
:::

## Wiring It Into Your Workflow

Running three tools by hand before every commit is easy to forget. The fix is to make the tools run *automatically* — first on your machine, then again on the server.

A **pre-commit hook** runs checks the instant you try to commit. If the formatter, linter, or type checker objects, the commit is blocked until you fix it. The popular `pre-commit` tool wires this up from a small config file:

```text
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.5.0
    hooks:
      - id: ruff          # lint
      - id: ruff-format   # format
  - repo: https://github.com/pre-commit/mirrors-mypy
    rev: v1.10.0
    hooks:
      - id: mypy
```

You install the hooks once, and from then on they run on every commit:

```bash
pip install pre-commit
pre-commit install
```

The second line of defense is **CI** (Continuous Integration) — a server that re-runs these same checks on every push, so nothing slips through even if someone skipped the local hook. At a high level, a CI job just runs the very commands you already know:

```bash
ruff check .
ruff format --check .
mypy src/
pytest
```

If any command fails, the whole check fails, and the team sees a red mark before the code is merged.

:::tip
Adopt these tools gradually on an existing project. Turn on the formatter first (it's painless), then the linter, then the type checker with relaxed settings you tighten over time. Trying to satisfy all three at once on a big old codebase is discouraging; one at a time is a pleasure.
:::

## Practice

:::predict
Q: A function is annotated `def double(n: int) -> int:` and somewhere your code calls `double("hi")`. You run the program normally with `python`. What happens?
- Python raises a TypeError at runtime because the hint is violated
- Python runs the code; the hint is ignored at runtime and no error is raised *
- The file refuses to load until the type is fixed
- Python automatically converts `"hi"` into an integer
E: Type hints are not enforced when Python runs — they're documentation the interpreter ignores. `"hi" * 2` just yields `"hihi"`. Only a separate tool like mypy reads the hints and reports the mismatch, and only when you run *it*.
:::

:::quiz
Q: Which tool's main job is to catch likely bugs and code smells — like an unused import or a variable that's assigned but never used — without running the code?
- black (a formatter)
- ruff / flake8 (a linter) *
- pip (a package installer)
- twine (a publishing tool)
E: That's linting. A linter such as ruff or flake8 statically reads your code and flags probable problems. black only reformats; it doesn't hunt for bugs.
:::

## Recap

- PEP 8 is Python's shared style guide. Consistency doesn't change behavior, but it makes code far faster to read and maintain.
- A **formatter** (black, or `ruff format`) rewrites code into a consistent style automatically — run it on save and stop thinking about formatting.
- A **linter** (ruff, flake8) statically flags likely bugs and smells like unused imports and dead variables; many are auto-fixable with `ruff check --fix`.
- A **type checker** (mypy) reads your type hints and verifies consistency. Python ignores hints at runtime, so mypy is what actually enforces them — catching type mismatches before you run the code.
- Wire the tools in with a **pre-commit hook** locally and **CI** on the server, so every commit and push is checked automatically.
- On an existing project, adopt the tools one at a time — formatter, then linter, then type checker — rather than all at once.

**Next up:** with your environment, structure, and quality tools in place, you're ready to build and ship Python projects with confidence.

# Project Structure and Packaging

A single `script.py` is fine for a quick experiment. But real projects grow: more files, tests, a description, dependencies. This lesson shows you how to lay a project out so its pieces find each other cleanly and — when you're ready — can be shared with the world. You'll do all of this on your own machine, in your own terminal, not in the browser labs.

## Why Layout Matters

When everything lives in one giant file, you eventually can't find anything, and you certainly can't test one part without dragging in the rest. A good layout does three quiet but important jobs: it groups related code together, it makes *imports* between your own files predictable, and it tells tools (and teammates) exactly where things live.

:::analogy
Think of a project like a well-organized kitchen. Ingredients go in the pantry, tools in the drawers, the recipe on the counter. You *could* pile everything on one table, but the moment you need the whisk you'd be digging through the flour. Structure is just knowing where to reach.
:::

The good news: Python's conventions are simple and widely shared. Learn them once and every Python project you open will feel familiar.

## A Real Project Layout

Here's a typical small project. Take a moment to read the whole tree before we walk through it.

```text
my-package/
├── README.md
├── pyproject.toml
├── requirements.txt
├── src/
│   └── myapp/
│       ├── __init__.py
│       ├── core.py
│       └── utils.py
└── tests/
    ├── __init__.py
    ├── test_core.py
    └── test_utils.py
```

Four things at the top level, and two folders that do the real work. Let's take them in turn.

- **`README.md`** — the front door. What the project is, how to install it, how to run it. The first thing a human reads.
- **`pyproject.toml`** — the project's ID card and config file. More on this below.
- **`requirements.txt`** — the pinned dependency list from the previous lesson.
- **`src/myapp/`** — your actual code, gathered into a *package*.
- **`tests/`** — your tests, kept separate from the code they check.

## Packages and `__init__.py`

A **package** is just a folder of Python files that belong together. What turns an ordinary folder into a package is a file named `__init__.py` inside it. It can be completely empty — its mere presence says, "this folder is an importable package."

In the tree above, `src/myapp/` is a package because it contains `__init__.py`. That lets other code write `import myapp` and reach the files inside.

```python
# src/myapp/__init__.py
# Often empty. Sometimes used to expose a tidy public interface:
from myapp.core import greet

__all__ = ["greet"]
```

That optional last touch means someone can write `from myapp import greet` directly, without knowing it actually lives in `core.py`. You decide what your package presents to the outside.

:::tip
An empty `__init__.py` is perfectly normal and very common. Don't feel you need to put anything in it — its job is done just by existing.
:::

## Importing Between Your Own Modules

Each `.py` file is a **module**. Once your code is a package, one module can import another using dotted paths that start from the package name.

Say `utils.py` has a small helper and `core.py` wants to use it:

```python
# src/myapp/utils.py
def shout(text: str) -> str:
    return text.upper() + "!"
```

```python
# src/myapp/core.py
from myapp.utils import shout

def greet(name: str) -> str:
    return shout(f"hello {name}")
```

The line `from myapp.utils import shout` reads as "from the `utils` module inside the `myapp` package, bring in `shout`." Because both files sit under `myapp/`, they can reference each other by the package name.

:::key
Always import your own modules using the full package path (`from myapp.utils import shout`), not a bare `from utils import shout`. The full path is unambiguous and works no matter where the program is started from.
:::

Your tests import the same way, reaching into the package to check it:

```python
# tests/test_core.py
from myapp.core import greet

def test_greet():
    assert greet("ada") == "HELLO ADA!"
```

## The `src/` Layout Idea

You may have noticed the code lives under `src/myapp/` rather than just `myapp/` at the top level. This is the **src layout**, and it's a deliberate, recommended choice.

The reason is subtle but valuable. When your package sits at the top level, Python can accidentally import it straight from the folder you're standing in — even if you never actually *installed* it. Your tests might pass locally for the wrong reason, then fail for a user who installs the real thing. Tucking the code inside `src/` removes that accident: to import `myapp`, you must install it properly (even in editable mode, shown below). So your tests run against the package the way a real user gets it.

:::warning
Without a `src/` folder, tests can silently import your code from the working directory instead of the installed package. That hides packaging bugs until someone else hits them. The `src/` layout forces an honest install, catching those problems early.
:::

To develop against your own package, install it in **editable mode** from the project root, inside your activated virtual environment:

```bash
pip install -e .
```

The `-e` means "editable": Python links to your source folder, so edits take effect immediately without reinstalling. The `.` means "the project in this directory."

## `pyproject.toml` Basics

`pyproject.toml` is the modern, standard place to describe your project. TOML is a plain, readable config format. A minimal but complete file looks like this:

```toml
[build-system]
requires = ["setuptools>=68"]
build-backend = "setuptools.build_meta"

[project]
name = "myapp"
version = "0.1.0"
description = "A friendly greeting library."
readme = "README.md"
requires-python = ">=3.12"
dependencies = [
    "requests>=2.32",
]

[project.optional-dependencies]
dev = ["pytest>=8.0"]
```

Reading it top to bottom:

- **`[build-system]`** tells tools *how* to build your project — which tool turns your source into an installable package. `setuptools` is a common, safe default.
- **`[project]`** is the ID card: the name people will `pip install`, the version, a one-line description, and where the README lives.
- **`requires-python`** states which Python versions you support.
- **`dependencies`** lists the libraries your code needs to run. These get installed automatically when someone installs your package.
- **`[project.optional-dependencies]`** groups extras — here, `dev` tools only contributors need. A developer installs them with `pip install -e ".[dev]"`.

:::example
Because `pytest` is listed under `dev` (not the main `dependencies`), an end user who just wants to *use* your greeting library never has to download a testing framework they'll never run. It only shows up for people actually working on the code.
:::

## Building and Publishing to PyPI

Once your project has a `pyproject.toml`, you can turn it into shareable files and, if you like, publish it to PyPI so anyone can `pip install` it. This is a high-level tour — you won't need it for a while, but it's good to see the shape.

First, build the distributable files:

```bash
pip install build
python -m build
```

This creates a `dist/` folder containing two files: a `.whl` (a "wheel," the fast pre-built format pip prefers) and a `.tar.gz` (the source archive). These are the exact artifacts PyPI serves to everyone else.

Then upload them with **twine**, the standard publishing tool:

```bash
pip install twine
twine upload dist/*
```

Twine asks for your PyPI credentials and pushes your wheel and source archive to the public index. Moments later, `pip install myapp` works for anyone, anywhere.

:::tip
Before publishing to the real PyPI, practice on **TestPyPI** (`twine upload --repository testpypi dist/*`), a sandbox copy of the index. You can push, install, and verify without claiming a real package name or worrying about mistakes.
:::

## Practice

:::quiz
Q: What single thing turns an ordinary folder of `.py` files into an importable Python package?
- A `pyproject.toml` in the folder
- A file named `__init__.py` inside the folder *
- Naming the folder `src`
- Running `python -m build`
E: The presence of `__init__.py` (even an empty one) marks a folder as a package, which is what makes `import myapp` work. `pyproject.toml`, `src/`, and `build` all play other roles.
:::

:::predict
Q: Your package lives at `src/myapp/`, and `core.py` needs the `shout` function from `utils.py` in the same package. Which import is the correct, unambiguous choice?
- `from myapp.utils import shout` *
- `import shout`
- `from src.myapp.utils import shout`
- `from ...utils import shout`
E: Import by the full package path: `from myapp.utils import shout`. The `src/` folder is a layout detail that isn't part of the import path, and a bare `import shout` would fail once the code is installed as a package.
:::

## Recap

- A real project separates code (`src/myapp/`) from tests (`tests/`), with a `README.md`, `pyproject.toml`, and `requirements.txt` at the top.
- A folder becomes an importable *package* the moment it contains an `__init__.py`, which may be empty.
- Import your own modules by full package path — `from myapp.utils import shout` — so imports work no matter where the program starts.
- The `src/` layout forces an honest install, so tests run against the package as a real user would receive it. Develop against it with `pip install -e .`.
- `pyproject.toml` is your project's ID card: name, version, description, supported Python, and dependencies (with optional `dev` extras).
- To share: `python -m build` produces a wheel and source archive in `dist/`, and `twine upload dist/*` publishes them to PyPI. Practice on TestPyPI first.

**Next up:** Code Style, Linting and Type Checking — the tools that keep your growing project readable and catch mistakes before they ship.

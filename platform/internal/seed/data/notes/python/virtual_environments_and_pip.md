# Virtual Environments and pip

Every real Python project depends on other people's code — libraries you install rather than write. This lesson is about doing that cleanly, so that two projects on your machine never step on each other's toes. It's a workflow lesson: you'll run these commands in your own terminal, not in the browser labs.

## The Problem Virtual Environments Solve

When you install Python, it comes with one shared, global place to store extra libraries. The moment you type `pip install requests`, that library lands in *that one place*, available to every project you own.

That sounds convenient until you have two projects. Project A needs version 1 of some library. Project B was written last year against version 2, and version 2 changed how things work. There is only one global slot, so installing what B needs *breaks* A. You can't have both.

:::analogy
A global install is like one shared kitchen for the whole apartment building. The moment your neighbor swaps the flour for gluten-free, *your* cake recipe changes too. A virtual environment gives each project its own private kitchen — its own shelf of ingredients, at exactly the versions that recipe expects.
:::

A **virtual environment** ("venv" for short) is a self-contained folder that holds a copy of Python plus its own private set of installed libraries. Activate it, and `pip install` puts libraries *there* instead of in the global spot. Each project gets its own, and they never interfere.

:::key
One virtual environment per project. It is the single habit that prevents the most common "it works on my machine" headaches in Python.
:::

## Creating a Virtual Environment

Python ships with a built-in tool for this, called `venv`. You run it from inside your project folder. The convention is to name the environment folder `.venv` (the leading dot keeps it tidy and out of the way).

```bash
cd my-project
python -m venv .venv
```

That's it. `python -m venv` means "run the venv module," and `.venv` is the folder it creates. Inside, you'll find a private Python and an empty library shelf, waiting.

:::tip
On some systems the command is `python3` rather than `python`. If `python -m venv .venv` complains, try `python3 -m venv .venv`. On Python 3.12 either works once your environment is active.
:::

Creating the folder does **not** turn it on. It just sits there until you activate it.

## Activating and Deactivating

Activating a venv tells your current terminal session, "for now, use *this* project's private Python and libraries." The command differs by operating system, so here are both.

On **macOS and Linux**:

```bash
source .venv/bin/activate
```

On **Windows** (PowerShell):

```bash
.venv\Scripts\Activate.ps1
```

Once active, your prompt changes to show the environment name in parentheses, like a little badge:

```bash
(.venv) my-project $
```

That badge is your signal that anything you install now goes into the project, not the global pile. To confirm which Python you're actually using:

```bash
which python      # macOS / Linux
where python      # Windows
```

The path should point *inside* your `.venv` folder. When you're done working, turn it off with:

```bash
deactivate
```

The badge disappears and you're back to the system Python.

:::warning
Never commit the `.venv` folder to Git. It's large, machine-specific, and fully rebuildable from your requirements file (coming up next). Add a line reading `.venv/` to your `.gitignore`.
:::

## Installing Packages with pip

`pip` is Python's package installer. It fetches libraries from the Python Package Index (PyPI, a huge public catalog) and drops them into whichever environment is active.

```bash
pip install requests
```

You can install a specific version, or install several at once:

```bash
pip install requests==2.32.3
pip install rich pytest
```

To see what's currently installed in the active environment:

```bash
pip list
```

To remove something:

```bash
pip uninstall requests
```

Because your venv is active, all of this touches *only* this project. Deactivate, switch to another project, activate its venv, and you'll see an entirely different set of libraries.

:::example
You start a web-scraping project. You activate its venv and run `pip install requests beautifulsoup4`. Over in your data project's venv, `pip list` shows `pandas` and `matplotlib` instead — no overlap, no conflict. Two private kitchens, two very different pantries.
:::

## Recording and Reinstalling with requirements.txt

Here's the payoff. Your venv currently holds an exact set of libraries at exact versions. You want to *capture* that list so a teammate — or future you, on a new laptop — can recreate it perfectly.

`pip freeze` prints every installed library with its pinned version. Redirect that into a file, conventionally named `requirements.txt`:

```bash
pip freeze > requirements.txt
```

The file looks like this:

```text
beautifulsoup4==4.12.3
certifi==2024.7.4
requests==2.32.3
soupsieve==2.5
```

Commit *that* file to Git (it's tiny and reproducible). Now anyone can rebuild the exact environment: create a fresh venv, activate it, and install from the list.

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

The `-r` flag means "read requirements from this file." Every library, at every recorded version, reinstalls in one command.

:::key
The pattern to memorize: `pip freeze > requirements.txt` to *save* the current state, and `pip install -r requirements.txt` to *restore* it. Save when you change dependencies; restore when you set up the project somewhere new.
:::

## Understanding Version Pins

Those `==2.32.3` numbers follow **semantic versioning**, written as `MAJOR.MINOR.PATCH`:

- **MAJOR** (the first number) changes when the library breaks backward compatibility — old code may stop working.
- **MINOR** (the middle) changes when new features are added, but old code keeps working.
- **PATCH** (the last) changes for bug fixes only.

So `requests==2.32.3` means major 2, minor 32, patch 3. Knowing this lets you loosen a pin deliberately. In a `requirements.txt` you might write:

```text
requests==2.32.3          # exactly this version, nothing else
rich>=13.0                # this version or any newer one
pytest>=8.0,<9.0          # 8.x only — accept features and fixes, avoid the next major
```

That last line is a sweet spot: `>=8.0,<9.0` accepts any 8.x release (safe minor and patch updates) but refuses to jump to 9.0, where something might break.

:::tip
When in doubt, pin exactly with `==`. Exact pins make your builds *reproducible* — the same install today and six months from now. Loosen pins only when you have a reason and a way to test the result.
:::

## A Word on Newer Tools

`venv` and `pip` are built in, universal, and enough for everything in this course — learn them first. But you'll hear about faster or more full-featured alternatives: **uv** (a very fast drop-in replacement that handles environments and installs together), **Poetry**, and **pipenv** (both add lockfiles and richer dependency management). They solve the same core problem you just learned; once these fundamentals click, picking one up later is easy.

## Practice

:::quiz
Q: You run `pip install pandas` but forgot to activate your project's virtual environment first. Where does `pandas` get installed?
- Into the project's `.venv` folder anyway
- Into the global, system-wide Python *
- Nowhere — pip refuses to run without a venv
- Into a brand-new venv pip creates automatically
E: With no venv active, pip installs into the shared global Python — exactly the mess venvs exist to avoid. Always check for the `(.venv)` badge in your prompt before installing.
:::

:::predict
Q: A teammate clones your repo, which includes `requirements.txt`. Which single command reinstalls every dependency at the recorded versions?
- `pip freeze > requirements.txt`
- `pip install -r requirements.txt` *
- `pip list requirements.txt`
- `python -m venv requirements.txt`
E: `pip install -r requirements.txt` reads the file and installs each pinned library. `pip freeze >` does the opposite — it *writes* the current state out to the file.
:::

## Recap

- Global installs share one library shelf across all projects, so different version needs collide. A virtual environment gives each project its own private shelf.
- Create one with `python -m venv .venv`, then activate it: `source .venv/bin/activate` on macOS/Linux, `.venv\Scripts\Activate.ps1` on Windows. The `(.venv)` badge confirms it's on. `deactivate` turns it off.
- With a venv active, `pip install` affects only that project. Never commit the `.venv` folder — add it to `.gitignore`.
- `pip freeze > requirements.txt` saves the exact current versions; `pip install -r requirements.txt` restores them anywhere.
- Version pins follow `MAJOR.MINOR.PATCH`. Pin exactly with `==` for reproducible builds; loosen deliberately with ranges like `>=8.0,<9.0`.
- Newer tools — uv, Poetry, pipenv — solve the same problem with extra features, but venv and pip are the universal foundation.

**Next up:** Project Structure and Packaging — how to lay out a real Python project so its pieces import cleanly and can be shared.

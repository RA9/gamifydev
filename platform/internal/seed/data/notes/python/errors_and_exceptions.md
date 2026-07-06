# Errors and Exceptions

Every programmer writes code that breaks. The difference between a beginner and a pro isn't that the pro's code never fails — it's that the pro reads the error calmly, knows what it means, and writes code that *recovers* instead of crashing. In this lesson we'll decode Python's error messages, meet the exceptions you'll see most, and learn to catch them so a bad dice roll or a fat-fingered player input doesn't take the whole game down.

Keep the [Code Lab](#code) open — the fastest way to understand an error is to trigger it yourself.

## Reading a traceback (bottom line first)

When Python hits an error it can't handle, it prints a **traceback**: a report of where things went wrong. It looks intimidating, but there's a trick — **read it from the bottom up**.

```python
def get_hp(player):
    return player["hp"]

hero = {"name": "Ada"}
print(get_hp(hero))
```

That produces:

```
Traceback (most recent call last):
  File "game.py", line 5, in <module>
    print(get_hp(hero))
  File "game.py", line 2, in get_hp
    return player["hp"]
KeyError: 'hp'
```

The **last line** is the headline: `KeyError: 'hp'`. That's *what* went wrong — we asked for a key that isn't there. The lines above it are the trail of *where*, most recent last: the error happened in `get_hp` at line 2, which was called from line 5. So the drill is: read the bottom line to learn the problem, then walk up to find the exact spot. Ninety percent of debugging is just reading that bottom line properly.

:::key
Always read a traceback from the **bottom up**. The final line names the error type and message — that's your real clue. The stack above it just shows the path the code took to get there.
:::

## Syntax errors vs. runtime exceptions

There are two fundamentally different kinds of problem, and telling them apart saves you time.

A **syntax error** means the code is not valid Python — Python can't even *start* running the file. A missing colon, an unclosed bracket, a stray indent.

```python
def start_game()      # SyntaxError: missing colon
    print("Go!")
```

Nothing runs; Python refuses the whole file until you fix the grammar.

A **runtime exception** (or just "exception") happens while the program is *already running*. The syntax was fine — Python got going — but then it hit something it couldn't do, like dividing by zero or opening a missing file.

```python
print("Starting...")      # this runs fine
print(10 / 0)             # ZeroDivisionError — crashes here, mid-run
```

The distinction: syntax errors are caught *before* execution and can't be handled with `try`; exceptions happen *during* execution and *can* be caught and recovered from. This whole lesson is about that second kind.

## Common exception types

You'll meet the same handful of exceptions over and over. Learn to recognise them on sight.

```python
int("五")            # ValueError — right type (str) but wrong content for int()
"hp: " + 100         # TypeError — can't add a string and an int
{"hp": 100}["mp"]    # KeyError — that key isn't in the dict
[1, 2, 3][9]         # IndexError — no item at position 9
10 / 0               # ZeroDivisionError — division by zero
print(scoer)         # NameError — 'scoer' was never defined (a typo!)
```

- **`ValueError`** — the value is the wrong *content*: `int("hello")`.
- **`TypeError`** — the value is the wrong *type* for the operation: adding a string to a number.
- **`KeyError`** — a dictionary key that doesn't exist.
- **`IndexError`** — a list/string index past the end.
- **`ZeroDivisionError`** — dividing by zero.
- **`NameError`** — using a variable that was never defined (usually a typo).

:::tip
The name tells you the category and the message tells you the specifics. `ValueError: invalid literal for int() with base 10: 'hello'` reads as "I tried to make an int, but 'hello' isn't a number." Slow down and *read the words* — they're more helpful than they first look.
:::

## try / except: catching an exception

To stop an exception from crashing your program, wrap the risky code in a `try` block and handle the fallout in an `except` block. If the `try` code raises the named exception, Python jumps to `except` instead of crashing.

```python
raw = "not a number"

try:
    number = int(raw)
    print(f"You rolled {number}")
except ValueError:
    print("That wasn't a valid number — try again.")
```

Because `int("not a number")` raises `ValueError`, the `int(...)` line fails, the rest of `try` is skipped, and the `except` block runs its friendly message. The program keeps going instead of dying. This is the core pattern for handling anything you don't fully control — user input, files, network data.

## Catch specifically, not blindly

You can write a bare `except:` that catches *everything*, but you almost never should.

```python
# Too broad — hides real bugs:
try:
    hp = int(player_input)
except:                          # catches ANYTHING, even typos in your own code
    print("Something went wrong")
```

The problem: a bare `except` swallows *every* error, including a `NameError` from your own typo or a `KeyboardInterrupt` when the user tries to quit. You wanted to handle bad input, but you've accidentally hidden bugs that should have been loud. Name the exception you actually expect:

```python
try:
    hp = int(player_input)
except ValueError:
    print("Please enter a whole number.")
```

Now only the case you planned for is caught; a genuine bug still surfaces its traceback so you can fix it. You can handle several types by listing more `except` blocks:

```python
try:
    hp = data["hp"] / turns
except KeyError:
    print("No hp field in the save data.")
except ZeroDivisionError:
    print("Can't divide by zero turns.")
```

:::warning
Avoid the bare `except:`. It catches errors you never meant to catch — including your own bugs — and turns a clear crash into a silent mystery. Always name the specific exception(s) you expect to handle.
:::

## else and finally

Two optional clauses round out the pattern.

- **`else`** runs only if the `try` block succeeded with *no* exception. It's for the "happy path" code that should run when nothing went wrong.
- **`finally`** runs *no matter what* — success, failure, even if you return early. It's for cleanup you can't skip, like closing a file.

```python
try:
    roll = int(input_value)
except ValueError:
    print("Not a number.")
else:
    print(f"Great roll: {roll}")   # only if int() succeeded
finally:
    print("Turn complete.")        # always runs
```

If `input_value` is `"7"`, you see the `else` message and then `finally`. If it's `"oops"`, you see the `except` message and then `finally`. Either way, `finally` fires — that guarantee is what makes it perfect for tidy-up work.

## Raising your own exceptions

Sometimes *you* need to signal that something is wrong — a value your function refuses to accept. Use `raise` with an exception type and a helpful message.

```python
def set_difficulty(level):
    if level not in ("easy", "normal", "hard"):
        raise ValueError(f"Unknown difficulty: {level!r}")
    return level

set_difficulty("nightmare")   # ValueError: Unknown difficulty: 'nightmare'
```

Raising a clear exception is far kinder than letting bad data slip deeper into your program, where it'll cause a confusing crash three functions later. Fail early, fail with a message that says exactly what's wrong. Pick the type that fits: `ValueError` for a bad value, `TypeError` for a wrong type.

## EAFP: ask forgiveness, not permission

Python has a cultural style with a memorable name: **EAFP** — "Easier to Ask Forgiveness than Permission." Rather than checking whether something *will* work before you try it, you just *try it* and handle the exception if it fails.

Compare the two approaches for reading a player's HP from a dict:

```python
# LBYL — "Look Before You Leap": check first
if "hp" in save and save["hp"] > 0:
    hp = save["hp"]
else:
    hp = 100

# EAFP — the Pythonic way: try, then handle failure
try:
    hp = save["hp"]
except KeyError:
    hp = 100
```

Both work, but the EAFP version reads as "grab the hp; if it's missing, default to 100." It also avoids a subtle race where the key vanishes *between* your check and your use. In Python, reaching for `try`/`except` isn't admitting defeat — it's the idiomatic, readable way to handle the exceptional case.

:::example
Parsing player input safely, EAFP-style — keep asking until you get a valid number:

```python
def read_bet(raw):
    try:
        bet = int(raw)
    except ValueError:
        return None          # signal "invalid" to the caller
    if bet <= 0:
        return None
    return bet

print(read_bet("50"))    # 50
print(read_bet("lots"))  # None — caught the ValueError
print(read_bet("-5"))    # None — valid int, but rejected
```

The function tries the conversion, catches the one failure it expects, and hands back a clean result the caller can trust.
:::

## Practice

:::predict
Q: Which exception does this raise?
```python
d = {"hp": 100}
print(d["mp"])
```
- IndexError
- KeyError *
- ValueError
E: Looking up a dictionary key that doesn't exist raises `KeyError`. `IndexError` is for lists; `ValueError` is for wrong-content values.
:::

:::predict
Q: What does this print?
```python
try:
    n = int("7")
except ValueError:
    print("bad")
else:
    print("ok")
finally:
    print("done")
```
- bad, then done
- ok, then done *
- ok only
E: `int("7")` succeeds, so `except` is skipped, `else` runs ("ok"), and `finally` always runs ("done"). The `else` block fires precisely because no exception occurred.
:::

## Recap

- Read tracebacks **bottom-up**: the last line names the error; the stack above shows how you got there.
- **Syntax errors** stop the file before it runs; **runtime exceptions** happen mid-execution and can be caught.
- Know the usual suspects: `ValueError`, `TypeError`, `KeyError`, `IndexError`, `ZeroDivisionError`, `NameError`.
- Wrap risky code in `try` and handle it in `except`; add `else` for the success path and `finally` for guaranteed cleanup.
- Catch **specific** exceptions, never a bare `except:` — it hides your own bugs.
- Signal your own errors with `raise ValueError("clear message")`, failing early and loudly.
- Embrace **EAFP**: try the operation and handle the exception, rather than checking everything up front.

**Next up:** you've got the tools to build programs that bend instead of break. Take these patterns into the Code Lab and deliberately cause a few errors — nothing teaches exceptions like catching them yourself.

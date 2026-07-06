# Context Managers

Some things must be *undone* no matter what. You open a file — you must close it. You grant the player a temporary speed buff — you must remove it when the timer ends, even if the level crashes halfway through. Doing that cleanup by hand with `try`/`finally` works but clutters every call site. Python's `with` statement packages "set up, then guarantee tear-down" into one tidy construct. This lesson shows how `with` works, how to build your own context managers two different ways, and why they make cleanup bulletproof.

## The with Statement and Why It Matters

Consider opening a file the manual way. You must remember to close it — and close it *even if an error is thrown in between*:

```python
f = open("scores.txt")
try:
    data = f.read()
    process(data)          # if this raises, we still need to close f
finally:
    f.close()              # runs no matter what
```

That `try`/`finally` dance is correct but noisy, and it's easy to forget the `finally`. The `with` statement does it for you:

```python
with open("scores.txt") as f:
    data = f.read()
    process(data)
# f is now closed — automatically, even if process() raised
```

When the `with` block ends — whether it finished normally, hit a `return`, or blew up with an exception — Python guarantees the cleanup runs. The object after `with` is a **context manager**: something that knows how to set itself up on entry and tear itself down on exit.

:::key
`with` guarantees that cleanup happens no matter how the block ends — normal finish, early `return`, or an exception mid-way. That guarantee is the entire point: you never leak a file handle, a lock, or a leftover buff.
:::

:::analogy
A context manager is an airlock. Walking in, the outer door seals behind you (setup). Whatever you do inside, when you leave, the inner door *always* closes — you can't accidentally leave it hanging open. The airlock enforces the ritual so you don't have to remember it.
:::

## File Handling: The Classic Example

Files are the textbook case because an unclosed file is a real resource leak — the operating system only lets you hold so many open at once. `with` closes them for you the instant the block ends.

```python
# Writing a save file
with open("save.txt", "w") as f:
    f.write("level=5\n")
    f.write("hp=30\n")
# file flushed and closed here

# Reading it back
with open("save.txt") as f:
    for line in f:
        print(line.strip())
# → level=5
# → hp=30
```

:::warning
Real file input/output like `open(...)` runs on *your own machine*, where there's a filesystem. The in-browser practice labs on this platform have no real disk, so `open()` won't behave there. Read these file examples to understand the pattern, then run them in a normal Python environment on your computer.
:::

You can even open several context managers in one `with` by separating them with commas — handy for copying one save file to another:

```python
with open("save.txt") as src, open("backup.txt", "w") as dst:
    dst.write(src.read())
# both files closed automatically
```

## The __enter__ / __exit__ Protocol

Under the hood, any object works with `with` if it implements two special methods. `__enter__` runs when the block starts and its return value is what `as` binds. `__exit__` runs when the block ends and does the cleanup. Let's build a scoped speed-buff manager for a game — it applies a buff on entry and always removes it on exit.

```python
class SpeedBuff:
    def __init__(self, player, amount):
        self.player = player
        self.amount = amount

    def __enter__(self):
        self.player["speed"] += self.amount
        print(f"Buff applied: speed now {self.player['speed']}")
        return self.player            # bound to the name after `as`

    def __exit__(self, exc_type, exc_value, traceback):
        self.player["speed"] -= self.amount
        print(f"Buff removed: speed back to {self.player['speed']}")
        return False                  # don't suppress exceptions

hero = {"name": "Ada", "speed": 10}

with SpeedBuff(hero, 5) as p:
    print(f"Dashing at speed {p['speed']}!")
# → Buff applied: speed now 15
# → Dashing at speed 15!
# → Buff removed: speed back to 10
```

Trace the order: `SpeedBuff(hero, 5)` builds the manager, `__enter__` runs and applies +5, `as p` catches its return value, the block runs, then `__exit__` runs and removes the buff. The buff comes off cleanly whether the block succeeds or crashes.

Those three parameters on `__exit__` — `exc_type`, `exc_value`, `traceback` — are how the manager learns whether the block exited normally or via an exception. If the block raised, they hold the exception details; if it finished cleanly, all three are `None`. Returning `False` (or `None`) means "I did my cleanup, now let any exception continue propagating." Returning `True` would *swallow* the exception — occasionally what you want, but rarely.

```python
class SpeedBuff:
    # ... __init__ and __enter__ as above ...
    def __exit__(self, exc_type, exc_value, traceback):
        self.player["speed"] -= self.amount     # cleanup runs regardless
        if exc_type is not None:
            print(f"(cleaned up despite a {exc_type.__name__})")
        return False

hero = {"name": "Ada", "speed": 10}
try:
    with SpeedBuff(hero, 5) as p:
        raise RuntimeError("level crashed")
except RuntimeError:
    print("caught outside the with")
print(hero["speed"])   # → 10  (buff was still removed!)
```

Even though the block raised, `__exit__` ran first and restored the speed. That's the guarantee in action.

:::key
`__enter__` sets up and returns the value bound by `as`. `__exit__(exc_type, exc_value, traceback)` tears down — it runs on every exit path. Return `False` to let exceptions propagate (the normal choice); return `True` only if you deliberately want to suppress them.
:::

:::predict
Q: What does this print?
```python
class Tag:
    def __enter__(self):
        print("open")
        return self
    def __exit__(self, *args):
        print("close")

with Tag():
    print("inside")
```
- open / inside / close*
- inside / open / close
- open / close / inside
- inside / close
E: `__enter__` runs first ("open"), then the block body ("inside"), then `__exit__` on exit ("close").
:::

## contextlib.contextmanager with yield

Writing a whole class for a simple manager is a lot of ceremony. The `contextlib` module offers a shortcut: decorate a generator with `@contextmanager`, and everything *before* the `yield` becomes setup, everything *after* becomes cleanup. The value you `yield` is what `as` binds.

```python
from contextlib import contextmanager

@contextmanager
def speed_buff(player, amount):
    player["speed"] += amount               # setup (like __enter__)
    print(f"Buff applied: speed now {player['speed']}")
    try:
        yield player                        # hand control to the with-block
    finally:
        player["speed"] -= amount           # cleanup (like __exit__)
        print(f"Buff removed: speed back to {player['speed']}")

hero = {"name": "Ada", "speed": 10}
with speed_buff(hero, 5) as p:
    print(f"Dashing at speed {p['speed']}!")
# → Buff applied: speed now 15
# → Dashing at speed 15!
# → Buff removed: speed back to 10
```

The `yield` is the seam. Code runs up to `yield` when the block starts, the `with` body runs while the generator is paused there, and code after `yield` runs when the block ends. Wrapping the `yield` in `try`/`finally` is what makes cleanup fire even when the block raises — the `finally` is the generator-flavored equivalent of `__exit__` always running.

:::tip
Use `@contextmanager` for quick, one-off managers where a whole class would be overkill. Reach for the `__enter__`/`__exit__` class form when the manager holds real state, needs to inspect the exception, or you want it to be reusable and self-documenting.
:::

A favorite use is a scoped **timer** — start the clock on entry, report elapsed time on exit, so you can measure how long a level-load takes:

```python
from contextlib import contextmanager
import time

@contextmanager
def timer(label):
    start = time.perf_counter()
    try:
        yield
    finally:
        elapsed = time.perf_counter() - start
        print(f"{label} took {elapsed:.4f}s")

with timer("level load"):
    total = sum(range(1_000_000))     # some work to measure
# → level load took 0.0123s
```

Here we `yield` nothing — there's no value to hand back, we just want the timing to bracket the block. The `finally` guarantees the elapsed time prints even if the work inside throws.

:::quiz
Q: In a `@contextmanager` generator, which part runs as the cleanup (the `__exit__` equivalent)?
- Everything before the `yield`
- Everything after the `yield`*
- The `yield` line itself
- Nothing runs automatically
E: Code before `yield` is setup (like `__enter__`); code after `yield` is cleanup (like `__exit__`). Wrapping the `yield` in `try`/`finally` ensures that cleanup runs even if the block raises.
:::

:::warning
Always guard the `yield` in a `@contextmanager` with `try`/`finally` when the cleanup *must* run. If you write the cleanup after a bare `yield` with no `finally` and the block raises, the exception skips right past your cleanup — the buff never comes off, the file never closes.
:::

## Recap

- The `with` statement runs setup on entry and **guarantees** cleanup on exit — even on early `return` or an exception.
- Files are the classic example: `with open(...) as f:` closes the file automatically, no `try`/`finally` needed. (Real file I/O runs on your own machine, not the browser labs.)
- A class becomes a context manager by implementing `__enter__` (setup; its return value binds to `as`) and `__exit__` (teardown; always runs).
- `__exit__(exc_type, exc_value, traceback)` receives exception info; return `False`/`None` to let exceptions propagate, `True` to suppress them.
- `@contextlib.contextmanager` turns a generator into a context manager: setup before `yield`, cleanup after — wrap the `yield` in `try`/`finally` so cleanup always fires.
- Choose the class form for stateful, reusable managers; the `@contextmanager` form for quick ones like a scoped buff or a timer.

**Next up:** keep applying these tools — generators, decorators, and context managers together form the backbone of clean, resource-safe Python.

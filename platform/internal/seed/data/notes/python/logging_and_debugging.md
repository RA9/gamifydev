# Logging and Debugging

Every program you write will break. That is not a sign you are bad at this — it is the normal state of building software. What separates a struggling coder from a confident one is not writing bug-free code; it is knowing how to *find* the bug fast. This lesson gives you the tools: proper logging, reading tracebacks, and Python's built-in debugger.

## Why `print` Debugging Falls Short

When something goes wrong, the first instinct is to sprinkle `print()` calls everywhere.

```python
def take_damage(player, amount):
    print("got here")
    print(player)
    print(amount)
    player["hp"] -= amount
    print("hp now", player["hp"])
    return player
```

This *works*, and it is a fine first move. But it does not scale. Once you have twenty `print("got here")` lines across five files, your terminal is a wall of noise and you cannot tell which line printed what. Worse, you have to go back and *delete* every one before shipping — and you always miss a few, so `print("HELLO?????")` ends up in production.

:::analogy
`print` debugging is like leaving sticky notes all over your house to remember things. Fine for a day. But after a month the walls are covered, half are outdated, and you can't find the one that matters. Logging is a proper notebook — dated, searchable, and you can rip out a whole section at once.
:::

`print` has no *levels*, so you can't say "show me only the serious stuff." It has no timestamps, so you can't tell when something happened. And it always goes to the screen — you can't quietly redirect it to a file. Logging solves all three.

## The `logging` Module and Levels

Python ships with a `logging` module in the standard library — no install needed. Instead of `print`, you emit **log records**, and each record has a **severity level**. There are five, from least to most severe:

```text
DEBUG    → fine-grained detail, useful while developing
INFO     → normal events ("player joined", "level loaded")
WARNING  → something unexpected but survivable
ERROR    → something failed; a feature didn't work
CRITICAL → the whole program may be going down
```

The magic is that you can set a **threshold**: only records at or above that level actually appear. Turn the dial to `WARNING` and all your chatty `DEBUG` and `INFO` lines go silent — without deleting a single one.

The quickest way to get started is `logging.basicConfig`:

```python
import logging

logging.basicConfig(level=logging.DEBUG)

logging.debug("Spawning enemy at (12, 7)")
logging.info("Player entered the dungeon")
logging.warning("Save file is an old version")
logging.error("Could not load texture: goblin.png")
logging.critical("Out of memory — shutting down")
```

Run it and every line shows, because the threshold is `DEBUG` (the lowest). Change one line to `level=logging.WARNING` and the first two vanish. Nothing else about your code changes.

:::key
Log levels are a volume knob you set *once*. Write `logging.debug(...)` freely all through your code; the threshold decides what actually prints. This is the single biggest reason logging beats `print`.
:::

## Getting a Logger and Formatting Messages

Calling `logging.info(...)` directly uses the shared **root logger**. That is fine for a script, but in a real project you want a **named logger** per module so you can tell where a message came from. The convention is to name it after the module:

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    datefmt="%H:%M:%S",
)

logger = logging.getLogger(__name__)

def load_level(name):
    logger.info("Loading level: %s", name)
    if name not in ("forest", "cave", "castle"):
        logger.error("Unknown level %r — falling back to 'forest'", name)
        name = "forest"
    return name

load_level("swamp")
```

Output looks like:

```text
14:32:07 [INFO] __main__: Loading level: swamp
14:32:07 [ERROR] __main__: Unknown level 'swamp' — falling back to 'forest'
```

Notice two things. First, the `format` string with `%(asctime)s`, `%(levelname)s`, and so on lets you decide exactly what each line looks like — timestamp, level, logger name, message. Second, we pass the level name as `%s` and the value as a *separate argument* rather than building the string ourselves:

```python
logger.info("Player %s reached level %d", name, level)   # good
logger.info(f"Player {name} reached level {level}")       # avoid
```

The first form lets `logging` skip building the string entirely when the message is below threshold — a small but free performance win.

:::tip
Use `logger.exception("...")` inside an `except` block. It logs at `ERROR` level *and* automatically attaches the full traceback, so you capture what went wrong without any extra work.
:::

## Reading a Traceback (Bottom-Up)

When Python hits an error it can't handle, it prints a **traceback** and stops. Beginners see the red wall of text and panic. Don't — a traceback is a map, and the treasure is at the bottom.

```python
def apply_potion(inventory, slot):
    potion = inventory[slot]
    return potion["heal"]

def heal_player(player):
    return apply_potion(player["items"], 5)

player = {"items": ["sword", "shield"]}
heal_player(player)
```

Running this produces:

```text
Traceback (most recent call last):
  File "game.py", line 10, in <module>
    heal_player(player)
  File "game.py", line 7, in heal_player
    return apply_potion(player["items"], 5)
  File "game.py", line 3, in apply_potion
    potion = inventory[slot]
             ~~~~~~~~~^^^^^^
IndexError: list index out of range
```

Read it **bottom-up**. The very last line is the actual error: `IndexError: list index out of range`. That is *what* went wrong. The line just above it is *where* it happened: line 3, `potion = inventory[slot]`. The lines above that trace the path of calls that led there — `heal_player` called `apply_potion`, which is how you got here.

So the story is: we asked for slot `5` of a list that only has two items (`sword`, `shield`). The fix is upstream — `heal_player` passed a hard-coded `5`.

:::warning
Always read the **last line first**. The top of a traceback is where your program *started*; the bottom is where it *broke*. New coders waste time staring at the top. The error type and message at the very bottom tell you the most.
:::

## Debugging Techniques That Don't Need Tools

Half of debugging is a way of thinking, not a command. Three techniques carry you a long way.

**Rubber-duck debugging.** Explain your code, line by line, out loud, to an inanimate object — the classic is a rubber duck on your desk. It sounds silly and it works. The act of forming sentences forces you to slow down, and you catch the flawed assumption yourself, usually mid-sentence. "So this loops over the enemies, and for each one I... wait, I'm resetting `total` *inside* the loop."

**Bisecting.** When a bug appeared "sometime recently" and you don't know where, cut the search space in half. Does it break halfway through the code? Then the cause is in the first half — ignore the rest. Repeat. Ten rounds of halving narrows a thousand lines down to one. (Git can automate this across commits with `git bisect`.)

**Minimal reproduction.** Strip the problem down to the smallest snippet that still fails. Delete everything unrelated — the UI, the network call, the fancy config — until you have ten lines that break. Often the bug becomes obvious the moment the clutter is gone. And if it doesn't, you now have something tiny you can share when you ask for help.

:::example
A player's score shows `0` after a win. Instead of re-reading the whole 400-line game loop, you bisect: add one log line at the midpoint. Score is correct there — so the bug is in the *second* half. Bisect again. Within four checks you find `score = 0` sitting inside the "you win" handler, left over from a copy-paste. Found in minutes, not hours.
:::

## The Built-in Debugger: `pdb` and `breakpoint()`

Logging tells you what happened *after* the fact. Sometimes you want to freeze the program mid-run and poke around live. That is what a **debugger** does, and Python has one built in: `pdb`.

The modern way to drop into it is the `breakpoint()` function. Put it on any line, and when Python reaches it, execution pauses and hands you an interactive prompt:

```python
def calculate_damage(base, multiplier, armor):
    raw = base * multiplier
    breakpoint()          # execution pauses right here
    final = raw - armor
    return max(final, 0)

calculate_damage(10, 3, 5)
```

Run the script normally and you land at a `(Pdb)` prompt. Now you are *inside* the paused function and can inspect anything. The core commands:

```text
n   (next)      run the current line, stop at the next one
s   (step)      step INTO a function call on this line
c   (continue)  resume until the next breakpoint (or the end)
p   (print)     print a variable, e.g.  p raw
q   (quit)      abort the program
```

A typical session:

```text
(Pdb) p raw
30
(Pdb) p armor
5
(Pdb) n
> game.py(5)calculate_damage()
-> return max(final, 0)
(Pdb) p final
25
(Pdb) c
```

You printed `raw` (30) and `armor` (5), stepped one line with `n` to run `final = raw - armor`, confirmed `final` is `25`, then `c` to let the program finish. No `print` statements added, no code deleted afterward.

The difference between `n` and `s` trips people up: `n` (next) treats a function call as a single step and runs the whole thing, staying at your level. `s` (step) dives *into* that function so you can watch it run line by line. Use `s` when you suspect the bug is inside the call; use `n` when you trust it and want to skip past.

:::tip
`breakpoint()` respects an environment variable: run `PYTHONBREAKPOINT=0 python game.py` and every `breakpoint()` call is skipped entirely. So you can leave them in while developing and disable them all at once, without editing files.
:::

## Practice

:::predict
Q: In a Python traceback, which line tells you the actual error type and message?
- The first line under "Traceback (most recent call last)"
- The last line *
- The line with the highest line number
- Whichever line is indented the most
E: Tracebacks read bottom-up. The final line names the exception (like `IndexError: list index out of range`) — that is the actual error. The lines above trace the calls that led to it.
:::

:::quiz
Q: You set `logging.basicConfig(level=logging.WARNING)`. Which of these calls will actually print?
- `logging.debug("spawn coords")`
- `logging.info("player joined")`
- `logging.error("texture failed to load")` *
- `logging.debug("frame rendered")`
E: The threshold is `WARNING`, so only records at `WARNING` or higher (`WARNING`, `ERROR`, `CRITICAL`) appear. `debug` and `info` are below the line and stay silent — no code deletion needed.
:::

## Recap

- `print` debugging is fine for a quick check but has no levels, no timestamps, and must be cleaned up by hand.
- The `logging` module gives you five levels — DEBUG, INFO, WARNING, ERROR, CRITICAL — and a threshold you set once to control what prints.
- Use `logging.basicConfig` to configure format and level; grab a named logger with `logging.getLogger(__name__)` and pass values as arguments, not f-strings.
- Read tracebacks **bottom-up**: the last line is the error, the line above is where it happened, the rest is the path that got there.
- Rubber-ducking, bisecting, and minimal reproductions find bugs without any special tools.
- Drop `breakpoint()` on a line to pause live and inspect with `pdb`: `n` (next), `s` (step in), `c` (continue), `p` (print), `q` (quit).

**Next up:** Git and Version Control — how to save your progress, undo mistakes safely, and collaborate without stepping on each other's code.

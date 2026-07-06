# Project: Text Adventure Engine

This is your capstone. You've learned variables, loops, `input()`, dictionaries, functions, and maybe a class or two. Now you'll pull all of it together to build something people actually played for years before graphics existed: a **text adventure**. You type commands like `go north` and `take key`; the game describes where you are and what you can do. Behind that simple surface is a genuinely elegant program — a tiny *engine* driven entirely by data.

:::project
**Build a data-driven text adventure.** The world is a set of connected rooms, each described in plain text. The player moves between rooms by typing directions, picks up items into an inventory, and reaches a **win condition** (for example: reach the Treasure Vault while carrying the golden key). The whole thing runs locally in a file — say `adventure.py` — that you launch with `python adventure.py`. This is a *brief*, not a walkthrough: I'll give you the shape, the data structure, and the milestones, but the code is yours to write.
:::

:::tip
Like the command-line game before it, this project needs `input()`, which only works when you run Python on your **own machine**. Create `adventure.py` in a folder, and use the edit-save-run loop: change the file, save, run `python adventure.py` (or `python3`), see what happens, repeat. Build one milestone at a time and *play-test after each one* — never write the whole thing before running it.
:::

## The big idea: the world is data

The trick that makes text adventures manageable is this: **don't hard-code the map into your logic**. Instead, store the world as data — a dictionary of rooms — and write a small amount of engine code that reads that data. Adding a new room becomes editing the data, not rewriting the program. It's the same principle behind the quiz game's questions array: separate *what the world is* from *how the game runs*.

Here's a suggested structure. Each room is a dictionary; the whole map is a dictionary of rooms keyed by name:

```python
rooms = {
    "Entrance Hall": {
        "description": "A dusty hall. Cold air drifts from the north.",
        "exits": {"north": "Library", "east": "Kitchen"},
        "items": ["torch"],
    },
    "Library": {
        "description": "Towering shelves. A small brass key glints on a desk.",
        "exits": {"south": "Entrance Hall"},
        "items": ["brass key"],
    },
    "Kitchen": {
        "description": "Rusted pots hang everywhere. A locked door stands to the east.",
        "exits": {"west": "Entrance Hall", "east": "Treasure Vault"},
        "items": [],
    },
    "Treasure Vault": {
        "description": "Gold beyond counting. You win!",
        "exits": {"west": "Kitchen"},
        "items": [],
    },
}
```

Read that carefully — it's the whole game. `exits` maps a direction to the *name* of the room it leads to. Because rooms are keyed by name, following an exit is a single lookup: `rooms[destination_name]`. That's why a dictionary is the right tool.

:::key
The `exits` dictionary is a map from direction to a room *name string*, and that string is a key back into `rooms`. So `go north` becomes: look up `"north"` in the current room's exits to get a name, then look up that name in `rooms` to get the room. Two dictionary lookups and the player has moved. No `if` ladder listing every room needed.
:::

:::analogy
Think of `rooms` as a hotel directory and each room's `exits` as the little signs by the elevator: "North → Library". The sign doesn't contain the Library; it just names it. To actually get there you take the name to the directory and look it up. Your engine does exactly that, every turn.
:::

## Milestone checklist

Build in this order. Each milestone is playable, so you can test as you go:

1. **Show the current room.** Track a variable `current` holding the player's room name. Print the room's description and its available exits. Get this working with *no* movement yet.
2. **Parse commands.** Read a line with `input()`, lowercase it, and split it into words. Recognise `go <direction>`, `take <item>`, `look`, `inventory`, and `quit`.
3. **Move between rooms.** On `go north`, check the current room's `exits`. If that direction exists, update `current`. If not, tell the player they can't go that way.
4. **Pick up items into an inventory.** Keep an `inventory` list. On `take torch`, if the item is in the room, move it from the room's `items` list into `inventory`.
5. **Win condition.** When the player enters (or is standing in) the Treasure Vault while carrying the `"brass key"`, print a victory message and end the game.

## Parsing commands

Commands are just text you slice apart. Lowercase the whole line so `GO NORTH` and `go north` behave the same, then `.split()` into a list of words. The first word is the *verb*, the rest is the *target*.

```python
raw = input("> ").strip().lower()
words = raw.split()

if not words:
    continue  # player just pressed Enter; ask again

verb = words[0]
target = " ".join(words[1:])  # everything after the verb, rejoined

if verb == "go":
    # target is "north", "east", etc.
    ...
elif verb == "take":
    # target is "brass key", "torch", etc.
    ...
elif verb == "quit":
    break
```

Using `" ".join(words[1:])` for the target means multi-word items like `"brass key"` work naturally. You're building a miniature language interpreter — that's exactly what a command parser is.

:::warning
A player *will* type nonsense: `go up` into a wall, `take dragon` for an item that isn't there, or just a blank line. Handle every case with a friendly message ("You can't go that way." / "There's no dragon here."). A text adventure that crashes or silently ignores the player feels broken. Assume every input is hostile and check before you act on it.
:::

## Lean on functions (and maybe a class)

Don't cram everything into one giant loop. Break the work into functions with clear jobs — it keeps the engine readable and makes each piece testable on its own:

```python
def describe(room_name):
    room = rooms[room_name]
    print("\n" + room_name)
    print(room["description"])
    if room["items"]:
        print("You see: " + ", ".join(room["items"]))
    print("Exits: " + ", ".join(room["exits"].keys()))


def move(current, direction):
    exits = rooms[current]["exits"]
    if direction in exits:
        return exits[direction]      # the new room name
    print("You can't go that way.")
    return current                   # unchanged
```

Notice `move` returns the room name to become the new `current` — either the destination or the same room if the exit doesn't exist. Your main loop becomes a clean sequence: describe, read a command, dispatch it to the right function, repeat.

If you're comfortable with classes, wrap the state in a `Game` or `Player` class so `current`, `inventory`, and the methods that change them live together. That's a natural next step, but plain functions plus a few module-level variables work perfectly for a first version — don't over-engineer before it runs.

:::example
A tidy main loop, once your functions exist, looks like this:

```python
current = "Entrance Hall"
inventory = []

while True:
    describe(current)
    if current == "Treasure Vault" and "brass key" in inventory:
        print("The vault opens with your key. You win!")
        break

    raw = input("> ").strip().lower()
    words = raw.split()
    if not words:
        continue
    # ... dispatch verb/target to move(), take(), etc. ...
```

The win check sits right after `describe`, so the moment the player arrives holding the key, the game celebrates and stops.
:::

## Practice

:::quiz
Q: Why store rooms as a dictionary keyed by room name?
- Dictionaries are always faster than lists no matter what
- So you can look up the room a command points to by name in one step *
- Because rooms can't be stored in lists
- Because dictionaries automatically draw the map for you
E: A dict maps a room's name straight to its data, so "go north" becomes: read the exit name, then look up that room by its key. One direct lookup instead of scanning a list.
:::

:::quiz
Q: What's the main reason to split the engine into functions like `describe()` and `move()`?
- Functions make the program run in fewer lines total, always
- Each function has one clear job, so the code is readable and each piece can be tested on its own *
- Python refuses to run loops unless they call functions
- Functions store the player's inventory automatically
E: Small, single-purpose functions keep the main loop clean and let you test movement or descriptions in isolation. That's the same "separate concerns" idea behind storing the world as data.
:::

## Stretch goals

Ship the five milestones first — a working, winnable adventure. Then reach for these:

1. **Locked doors needing a key.** Give an exit a required item: the Kitchen's east door only opens if `"brass key" in inventory`. Store `"locked": "brass key"` on the exit and check it in `move()`.
2. **Enemies and combat.** Drop a monster into a room. Entering triggers a fight — reuse the HP-and-turns battle loop from the command-line game lesson. Block the exit until the monster is defeated.
3. **A `use` command.** Let items *do* things: `use torch` reveals a hidden exit in a dark room; `use potion` restores HP.
4. **Save and load with `json`.** Write the player's `current` room and `inventory` to a file so they can quit and resume. Python's `json` module makes this short:

```python
import json

def save_game(current, inventory):
    with open("save.json", "w") as f:
        json.dump({"current": current, "inventory": inventory}, f)

def load_game():
    with open("save.json") as f:
        data = json.load(f)
    return data["current"], data["inventory"]
```

Add `save` and `load` as commands. Now your adventure remembers the player between sessions — a feature real games charge money for.

5. **A richer map.** Add rooms, dead ends, and a longer path to the vault. Because the world is data, this is pure content work — no engine changes at all. That payoff *is* the lesson.

## Recap

- The world is **data**: a dictionary of rooms, each a dictionary of `description`, `exits`, and `items`.
- `exits` maps a direction to a room **name**, which is itself a key back into `rooms` — so moving is two dictionary lookups.
- **Parse commands** by lowercasing, `.split()`-ing into words, and treating the first word as the verb and the rest as the target.
- Follow the **milestones** in order — describe, parse, move, inventory, win — and play-test after each one.
- **Guard every input**; players will type things you didn't plan for, and a good game answers kindly instead of crashing.
- Split logic into **functions** (and optionally a class) so each piece is small, readable, and testable.
- **Stretch goals** — locked doors, enemies, `use`, and `json` save/load — build on the same data-driven engine without rewriting it.

**Next up:** you've built a real program with state, parsing, and persistence. That's the skeleton of countless applications — keep this engine, keep expanding the world, and make it unmistakably yours.

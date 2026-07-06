# Building a Command-Line Game

Up to now your Python has run inside the browser labs — and those are great for learning the language, but they have one hard limit: they can't stop and *wait for you to type*. A real game needs a player. It needs to ask "what's your guess?" and pause until you answer. That means it's time to graduate to running Python on your **own machine**, in a real file you can save and run again and again.

:::tip
This lesson is a *reading* lesson — you follow along and build on your own computer. Make a folder somewhere friendly (your Desktop is fine), and inside it create a file called `game.py`. To run it, open a terminal in that folder and type `python game.py` (or `python3 game.py` on macOS/Linux). Every time you change the file, save it and run that command again. That loop — edit, save, run — is the heartbeat of every programmer's day.
:::

We'll build a number-guessing game, then grow it into a tiny monster battle. You'll see the code get bigger in stages, and I'll explain *why* at each step, not just *what*.

## Step 1: Say hello and read a line

The star of local Python is `input()`. It prints a prompt, waits for the player to type something and press Enter, and hands you back whatever they typed — always as a **string**.

```python
name = input("What's your name, adventurer? ")
print("Welcome, " + name + "!")
```

Save that as `game.py` and run `python game.py`. The program stops at the question until you type your name. That pause is the thing the browser labs could never do — and it's the whole reason a game feels alive.

:::key
`input()` **always** returns a string, even if the player types `42`. You get the text `"42"`, not the number `42`. Before you can do math or comparisons with a number, you must convert it — that conversion is the source of about half of all beginner game bugs.
:::

## Step 2: A secret number and one guess

Let's pick a secret number the player has to find. Python's `random` module rolls the dice for us. We `import` it once at the top of the file, then call `random.randint(1, 100)` to get a whole number from 1 to 100 (both ends included).

```python
import random

secret = random.randint(1, 100)

guess = int(input("Guess my number (1-100): "))

if guess == secret:
    print("Incredible! You got it!")
else:
    print("Nope. I was thinking of " + str(secret) + ".")
```

Two conversions are doing quiet, important work here. `int(input(...))` turns the typed text into a real integer so `guess == secret` compares numbers, not text. And `str(secret)` turns the number back into text so we can glue it into a sentence with `+`. Numbers and strings don't mix with `+` — Python will raise an error if you try — so you convert deliberately at the boundary.

:::analogy
Think of `input()` as a mail slot: everything that comes through it arrives as an envelope of text. `int()` is you opening the envelope and pulling out the number inside. You can't do arithmetic on a sealed envelope.
:::

## Step 3: The main loop — keep guessing until you win

One guess is no fun. A game needs a **loop** that keeps asking until the player wins, with a hint each round. A `while True:` loop runs forever until we deliberately `break` out of it.

```python
import random

secret = random.randint(1, 100)
tries = 0

while True:
    guess = int(input("Guess my number (1-100): "))
    tries = tries + 1

    if guess < secret:
        print("Too low. Aim higher.")
    elif guess > secret:
        print("Too high. Come down a bit.")
    else:
        print("You got it in " + str(tries) + " tries!")
        break
```

This is the core of the whole game. Notice the shape: loop, read input, update state (`tries` counts up), give feedback (higher/lower), and only `break` when the win condition is met. Almost every game you'll ever write is some version of this loop.

:::example
Play it a few times. A smart strategy is to guess the middle (50), then the middle of whatever half is left (25 or 75), and so on. That's *binary search*, and it finds any number from 1–100 in at most 7 guesses. Games quietly teach real computer-science ideas.
:::

## Step 4: Don't let a typo crash the game

Here's a problem. Run the game and, when it asks for a guess, type `banana` instead of a number. The program crashes with a scary `ValueError`, because `int("banana")` is impossible. A polished game should never crash on bad input — it should shrug and ask again.

We wrap the risky conversion in `try` / `except`. If `int()` fails, the `except` block runs instead of the program dying, and `continue` jumps straight back to the top of the loop for another go.

```python
import random

secret = random.randint(1, 100)
tries = 0

while True:
    raw = input("Guess my number (1-100): ")

    try:
        guess = int(raw)
    except ValueError:
        print("That's not a whole number. Try again.")
        continue

    tries = tries + 1

    if guess < secret:
        print("Too low. Aim higher.")
    elif guess > secret:
        print("Too high. Come down a bit.")
    else:
        print("You got it in " + str(tries) + " tries!")
        break
```

:::warning
Validate input at the exact moment you convert it. If you let a bad value slip past `int()`, it poisons everything downstream — comparisons, scores, math. Catch it early, tell the player kindly, and loop back. A game that survives a fat-fingered `banna` feels ten times more professional than one that explodes.
:::

## Step 5: Play again?

When the player wins, offer another round. We wrap the *whole* game in an outer loop, and at the end ask a yes/no question. We normalise the answer with `.strip().lower()` so `"Yes"`, `"  y "`, and `"YES"` all count.

```python
import random

while True:  # one full game per pass
    secret = random.randint(1, 100)
    tries = 0

    while True:  # the guessing loop
        raw = input("Guess my number (1-100): ")
        try:
            guess = int(raw)
        except ValueError:
            print("That's not a whole number. Try again.")
            continue

        tries += 1
        if guess < secret:
            print("Too low. Aim higher.")
        elif guess > secret:
            print("Too high. Come down a bit.")
        else:
            print("You got it in " + str(tries) + " tries!")
            break

    again = input("Play again? (y/n) ").strip().lower()
    if again != "y":
        print("Thanks for playing!")
        break
```

Notice `tries += 1` — that's shorthand for `tries = tries + 1`, and you'll see it everywhere. The two nested loops read cleanly: the outer one is "another game?", the inner one is "another guess?". Getting comfortable with loops-inside-loops is a genuine milestone.

## Step 6: Grow it into a monster battle

The guessing game already has everything a bigger game needs: a loop, input, state, feedback, and a win condition. Let's reuse all of it to fight a monster. Now we track **HP** (hit points) for the player and the beast, and each turn the player chooses to attack or heal.

```python
import random

player_hp = 30
monster_hp = 25

print("A wild Bug-Beast appears! Defeat it before it defeats you.")

while player_hp > 0 and monster_hp > 0:
    print("\nYour HP: " + str(player_hp) + "   Monster HP: " + str(monster_hp))
    move = input("Do you [a]ttack or [h]eal? ").strip().lower()

    if move == "a":
        damage = random.randint(4, 9)
        monster_hp -= damage
        print("You strike for " + str(damage) + " damage!")
    elif move == "h":
        healed = random.randint(3, 7)
        player_hp += healed
        print("You patch yourself up for " + str(healed) + " HP.")
    else:
        print("You hesitate, wasting your turn...")

    # The monster always gets a turn (unless it's already down)
    if monster_hp > 0:
        hit = random.randint(3, 8)
        player_hp -= hit
        print("The Bug-Beast bites back for " + str(hit) + " damage!")

if player_hp > 0:
    print("\nThe Bug-Beast is defeated. Victory!")
else:
    print("\nYou have fallen. Game over.")
```

Look at the loop condition: `while player_hp > 0 and monster_hp > 0`. The fight continues only while *both* are still standing; the moment either drops to zero or below, the loop ends and we check who survived. The `random.randint` calls give every hit a little unpredictability, so no two battles feel the same. And `-=` is just like `+=` — `monster_hp -= damage` means `monster_hp = monster_hp - damage`.

:::key
State is the memory of your game. In the guesser it was `secret` and `tries`; in the battle it's `player_hp` and `monster_hp`. Every turn, the loop *reads* the state, *changes* it based on the player's choice and some randomness, and *shows* it. Master that read-change-show rhythm and you can build any turn-based game.
:::

## Practice

:::quiz
Q: Why must you wrap `int(input(...))` in a `try` / `except ValueError` block?
- Because `input()` runs faster inside a try block
- Because `int()` crashes the whole program if the player types something that isn't a number, and except lets you recover *
- Because Python requires every input to be inside try
- Because it converts the string to lowercase automatically
E: `int("banana")` raises a `ValueError`. Without `try`/`except`, that error stops the program dead. Catching it lets you print a friendly message and ask again instead of crashing.
:::

:::quiz
Q: In the monster battle, what does `while player_hp > 0 and monster_hp > 0:` guarantee?
- The loop runs exactly 30 times
- The fight keeps going only while BOTH fighters still have HP left *
- The player always wins
- The monster heals every turn
E: `and` means both conditions must be true to continue. As soon as either HP hits zero or below, the condition becomes false and the battle loop ends.
:::

## Make it yours

You've got two working games. Now make them *yours* — this is where the real learning happens:

1. **Difficulty levels.** Ask the player for easy (1–20), medium (1–100), or hard (1–500) before the guessing game starts, and pick the range from their answer.
2. **Limited guesses.** Give the guesser only 7 tries. Count down and end the game with a loss if they run out.
3. **Give the monster personality.** Store a list of monster names and pick one with `random.choice(["Bug-Beast", "Null Pointer", "Syntax Serpent"])`.
4. **A potion limit.** Let the player heal only 3 times per battle — track a `potions` counter so healing isn't infinite.
5. **A running scoreboard.** Across replays, keep a `wins` variable and print "You've won 3 battles!" as they go.
6. **Critical hits.** Give attacks a 1-in-5 chance (roll `random.randint(1, 5) == 1`) to deal double damage, with a "Critical hit!" message.

## Recap

- Local Python unlocks `input()`, which **pauses** the program and hands back what the player typed — always as a **string**.
- Convert at the boundary: `int()` to get a number in, `str()` to put a number into a sentence.
- The **main game loop** (`while True:` … `break`) is the beating heart: loop, read input, update state, give feedback, check the win condition.
- **Validate input** with `try` / `except ValueError` and `continue` so a typo can never crash the game.
- Wrap everything in a **play-again loop**, normalising the answer with `.strip().lower()`.
- **State** (`tries`, `player_hp`, `monster_hp`) is your game's memory; `+=` and `-=` update it each turn.
- The same loop that ran the guesser powers the **monster battle** — you already know how to build a game.

**Next up:** a bigger build — a data-driven text adventure with rooms, items, and an inventory you carry from place to place.

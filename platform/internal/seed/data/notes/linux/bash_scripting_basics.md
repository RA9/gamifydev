# Bash Scripting Basics

Until now you've typed commands one at a time, waiting for each to finish. That's powerful, but the real magic of the command line is this: you can write a list of commands down once, save it, and run the whole list whenever you like. That saved list is a **script**, and writing your first one is a genuine milestone. Let's automate something together.

## What is a script, really?

A **script** is just a plain text file containing commands — the very same commands you've been typing by hand. Instead of you typing them one by one, the shell reads them from the file, top to bottom, and runs each in turn.

That's the whole idea. Anything you can do at the prompt, you can put in a script. The payoff is enormous: a task that took ten careful commands becomes a single file you can run again and again, perfectly, without forgetting a step.

:::analogy
A script is a recipe card. Cooking a dish by improvising each step is like typing commands live — it works, but you have to remember everything. Writing the recipe down means anyone (including future you) can follow it exactly, every time, without missing the salt. The shell is the cook that follows your card.
:::

## The shebang: telling the system how to run it

Every Bash script starts with a special first line:

```bash
#!/bin/bash
```

This is called the **shebang** (from "hash-bang," the `#` and `!`). It tells the system, "run this file using the Bash program located at `/bin/bash`." Without it, the system might not know which interpreter to use. Make it the very first line of every script you write — it's a small piece of ceremony that pays off.

:::key
The shebang line `#!/bin/bash` must be the *first* line of your script. It points to the program that should run the file. Think of it as the script introducing which language it speaks before it starts talking.
:::

## Comments: notes for humans

Any line starting with `#` (other than the shebang) is a **comment**. The shell ignores it completely — comments are notes for *you* and anyone else reading the script later.

```bash
#!/bin/bash
# This script greets the user.
# Written while learning GamifyDev.
echo "Hello!"
```

Good comments explain *why* you did something, not just what. Your future self will thank you.

## Your first script

Let's write one. Create a file called `greet.sh` and put this inside:

```bash
#!/bin/bash
# A friendly greeting script
echo "Hello from my first script!"
echo "The command line is fun."
```

The `.sh` ending is a convention that signals "this is a shell script." It's not strictly required, but it helps you spot scripts at a glance.

## Making it executable and running it

Here's a wrinkle: you can't run your new script just yet. By default, a fresh text file isn't marked as *runnable*. You give it permission to run with `chmod +x` ("change mode, add execute"):

```bash
$ chmod +x greet.sh
```

Now run it by giving its path. Because the script is in your current folder, you point to it with `./` (remember, `.` means "right here"):

```bash
$ ./greet.sh
Hello from my first script!
The command line is fun.
```

That `./` matters. Without it, the shell looks for the program in its usual system locations and won't find your local file. The `./` says, "the script is right here in this folder."

:::tip
The two-step pattern — `chmod +x script.sh` once, then `./script.sh` to run — trips up nearly every beginner the first time. If you see "permission denied," it almost always means you forgot the `chmod +x` step. Pixel has seen it a hundred times; you're in good company.
:::

:::reorder
Order the steps to create and run a script from scratch.
- touch greet.sh
- chmod +x greet.sh
- ./greet.sh
E: First create the file, then mark it executable with `chmod +x`, then run it with `./`. Skip the middle step and you'll get "permission denied."
:::

## Variables: giving names to values

Scripts get far more useful when they can remember things. A **variable** is a named box that holds a value:

```bash
#!/bin/bash
NAME="Ada"
echo "Hello, $NAME!"
```

Running that prints `Hello, Ada!`. A few rules worth burning into memory:

- When you *set* a variable, there are **no spaces** around the `=`. Write `NAME="Ada"`, never `NAME = "Ada"`.
- When you *use* a variable, put a `$` in front of its name: `$NAME`.
- Wrapping it in double quotes — `"$NAME"` — is a safe habit that prevents surprises when the value contains spaces.

:::warning
The spacing rule catches everyone. `NAME = "Ada"` (with spaces) fails, because Bash thinks `NAME` is a command. Assignment must be tight: `NAME="Ada"`. When in doubt, no spaces around the `=`.
:::

:::predict
```bash
GREETING="Hi"
NAME="Sam"
echo "$GREETING, $NAME!"
```
- Hi, Sam! *
- $GREETING, $NAME!
- GREETING, NAME!
E: Bash replaces each `$VARIABLE` with its stored value, so the line becomes `Hi, Sam!`.
:::

## Reading input from the user

You can make a script interactive with `read`, which pauses and waits for the user to type something, then stores it in a variable:

```bash
#!/bin/bash
echo "What is your name?"
read NAME
echo "Welcome, $NAME!"
```

Run it, type your name, press Enter, and the script greets you personally. Just like that, your script holds a conversation.

## Making decisions with `if`

Scripts can choose what to do based on a condition, using `if`:

```bash
#!/bin/bash
echo "How old are you?"
read AGE
if [ "$AGE" -ge 18 ]; then
  echo "You're an adult."
else
  echo "You're young — enjoy it!"
fi
```

A few things to notice: the condition sits inside square brackets with spaces around them, `-ge` means "greater than or equal to," and the whole block closes with `fi` (that's `if` spelled backwards — a Bash quirk you'll grow fond of).

## Repeating with a `for` loop

When you want to do something several times, a `for` loop handles the repetition:

```bash
#!/bin/bash
for COLOR in red green blue; do
  echo "I like $COLOR"
done
```

This runs the body once for each item in the list, setting `COLOR` to `red`, then `green`, then `blue`. It prints three lines. The block opens with `do` and closes with `done`. Loops are how a script does in one breath what would take you many keystrokes by hand.

## A small worked script

Let's combine what you've learned into one real, useful script — a personalized greeting that reacts to the time of day:

```bash
#!/bin/bash
# Greet the user and list their favorite things.

echo "What is your name?"
read NAME

echo "Hello, $NAME! Nice to meet you."

# Print a short list using a loop
echo "Here are three things worth learning:"
for TOPIC in Linux Bash scripting; do
  echo " - $TOPIC"
done

echo "Happy coding, $NAME!"
```

Save it as `welcome.sh`, run `chmod +x welcome.sh`, then `./welcome.sh`. It asks your name, greets you, loops through a list, and signs off. Every piece of that script is something you learned in this lesson — and together they do real work.

:::fill
Q: Complete the line that makes a script runnable.
`chmod ___ welcome.sh`
- +x *
- -x
- cat
E: `chmod +x` adds execute permission so you can run the script with `./welcome.sh`.
:::

## Talk about it

> You've written your first scripts — with a shebang, comments, variables, input, an `if` decision, and a `for` loop. In your own words, why is putting commands in a script so much more powerful than typing them one at a time, and what does `chmod +x` actually do for you?

## What's next

You've crossed a real threshold — you're not just using Linux, you're *programming* it. In the final lesson you'll round out your foundation with file permissions, pipes, package managers, and a map of where to take these skills next. The adventure is almost a full loop.

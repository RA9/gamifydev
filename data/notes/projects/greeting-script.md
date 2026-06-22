# Greeting Script

Your first useful shell script: one that knows your name and the time of day. By the end you'll have a real command you can run any time you want a friendly hello from your terminal.

:::project
You'll have an executable Bash script called `greet.sh` that greets you by name and says "Good morning", "Good afternoon", or "Good evening" depending on the current hour.
:::

To follow along you need a terminal and a text editor. If you're new to the command line, warm up in the in-app [Terminal Trainer](#terminal) first, then come back.

## Step 1 — Create the file with a shebang

Make a file called `greet.sh`. The first line, the "shebang", tells the system to run it with Bash. Add an `echo` so the script does something visible.

```bash
#!/bin/bash
echo "Hello there!"
```

## Step 2 — Make it executable and run it

A new script isn't runnable by default. Give it execute permission with `chmod +x`, then run it by pointing at it with `./`. The `./` means "the file in this folder".

```bash
chmod +x greet.sh
./greet.sh
```

## Step 3 — Add a NAME variable

Greetings are nicer when they're personal. Define a variable and use it inside the message. Note there are no spaces around the `=`, and you read a variable back with a `$` in front of its name.

```bash
#!/bin/bash
NAME="Pixel"
echo "Hello, $NAME!"
```

## Step 4 — Get the current hour

Bash can ask the system for the time. The `date +%H` command prints the current hour in 24-hour format (00–23). We capture it into a variable using `$( )`, which runs a command and hands back its output.

```bash
#!/bin/bash
NAME="Pixel"
HOUR=$(date +%H)
echo "The current hour is $HOUR"
```

## Step 5 — Pick the greeting with if/elif/else

Now bring it together. Compare the hour with `if/elif/else` to choose the right greeting, then print it with the name. The `-lt` means "less than". This is the complete script.

```bash
#!/bin/bash
NAME="Pixel"
HOUR=$(date +%H)

if [ "$HOUR" -lt 12 ]; then
  GREETING="Good morning"
elif [ "$HOUR" -lt 18 ]; then
  GREETING="Good afternoon"
else
  GREETING="Good evening"
fi

echo "$GREETING, $NAME!"
```

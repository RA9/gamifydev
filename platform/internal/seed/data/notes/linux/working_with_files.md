# Working with Files

You can now walk the filesystem with confidence. Time to start changing it. This lesson is about the everyday craft of working with files: making them, reading them, writing into them, copying, moving, and — carefully — deleting them. These are the actions you'll perform countless times as a developer. Keep the [Terminal Trainer](#terminal) open and create real files as you go.

## Making folders and empty files

Let's build a little workspace. To make a new folder, use `mkdir` ("make directory"):

```bash
$ mkdir project
$ ls
project
```

To create a new, empty file, use `touch`:

```bash
$ touch notes.txt
$ ls
notes.txt  project
```

`touch` is wonderfully simple — name a file that doesn't exist yet, and it springs into being, empty and ready. (If the file already exists, `touch` just updates its timestamp without harming the contents.)

:::tip
A common first move when starting something new is `mkdir project` followed by `cd project` — make the folder, then step into it. You'll repeat that two-step pattern constantly, so let it become muscle memory.
:::

## Writing into a file with `>`

An empty file isn't very interesting. Let's put something inside. Remember `echo` from earlier — it prints text. If you add a `>` after it, that text flows *into a file* instead of onto the screen:

```bash
$ echo "Remember to water the plants" > notes.txt
```

The `>` is called **redirection**. Picture it as an arrow pointing where the text should go: instead of the screen, send it into `notes.txt`. To check that it worked, read the file back with `cat`:

```bash
$ cat notes.txt
Remember to water the plants
```

`cat` prints a file's contents straight to the screen — your quickest way to peek inside a small file.

:::warning
Be careful: a single `>` *overwrites*. It wipes whatever was in the file and replaces it with the new text. If `notes.txt` already had ten lines, `echo "hi" > notes.txt` throws them all away and leaves just `hi`. Always be sure you mean to replace before you use `>`.
:::

## Adding to a file with `>>`

When you want to *add* a line without destroying what's already there, use a double arrow, `>>`. It **appends** — it tacks the new text onto the end:

```bash
$ echo "First line" > notes.txt
$ echo "Second line" >> notes.txt
$ echo "Third line" >> notes.txt
$ cat notes.txt
First line
Second line
Third line
```

The first command, with `>`, created the file with one line. Each `>>` after it added another line on the end. One arrow replaces; two arrows append. That small difference matters enormously.

:::key
`>` writes and **overwrites** — it replaces everything. `>>` writes and **appends** — it adds to the end. When in doubt, reach for `>>` so you don't lose existing content by accident.
:::

:::predict
```bash
echo "hello" > greeting.txt
echo "world" >> greeting.txt
cat greeting.txt
```
- hello\nworld *
- world
- hello
E: `>` writes "hello" into the file, then `>>` appends "world" on a new line. `cat` prints both lines, so you see hello then world.
:::

## Copying and moving

Two more everyday actions: copying and moving. These commands aren't in the practice Trainer, but you'll use them constantly on a real system, so it's worth meeting them now.

**Copy** with `cp` ("copy"). You give it the original, then the destination:

```bash
$ cp notes.txt notes-backup.txt
```

That makes a duplicate called `notes-backup.txt` while leaving the original untouched — a quick way to keep a safety copy before you change something.

**Move or rename** with `mv` ("move"):

```bash
$ mv notes.txt project/notes.txt
```

`mv` relocates a file from one place to another. And here's a neat trick: moving a file to a new *name* in the same folder is exactly how you rename it.

```bash
$ mv notes.txt todo.txt
```

That "moved" `notes.txt` to `todo.txt` in place — effectively a rename.

:::analogy
Think of `cp` as a photocopier — the original stays on your desk and a duplicate comes out. `mv` is a courier — it picks the file up and carries it somewhere else, so it's no longer where it started. Renaming is just a very short courier trip that ends in the same room.
:::

## Deleting — with respect

To delete a file, use `rm` ("remove"):

```bash
$ rm notes-backup.txt
```

Gone. And that word — *gone* — deserves your full attention. On the Linux command line, `rm` does **not** send files to a recycle bin or trash. There's no "are you sure?" pop-up and no easy undo. When you remove a file, it's removed.

This isn't meant to scare you away from `rm` — you'll use it all the time. It's meant to build a healthy habit: read the line before you press Enter. Make sure the filename is the one you intend. Especially be cautious with wildcards like `*`, which can match far more files than you expect.

:::warning
`rm` deletes permanently — there is no trash to recover from. Before you run it, pause and read the exact filename you typed. A moment of care here saves you from a very bad afternoon. Treat `rm` with the same respect you'd give a sharp kitchen knife: useful, but never careless.
:::

## A full mini-session

Let's tie it all together with a sequence you can run in the [Terminal Trainer](#terminal) right now:

```bash
$ mkdir journal
$ cd journal
$ touch day1.txt
$ echo "Today I learned about files" > day1.txt
$ echo "It went well!" >> day1.txt
$ cat day1.txt
Today I learned about files
It went well!
$ ls
day1.txt
```

You made a folder, stepped in, created a file, wrote two lines into it, read them back, and listed what's there. That's the whole everyday loop of working with files — and you just did it.

:::quiz
Q: You run `echo "new" > log.txt` on a file that already contains five lines. What happens?
- The five lines are replaced; only "new" remains *
- "new" is added after the five lines
- The command fails with an error
- A backup of the old lines is saved automatically
E: A single `>` overwrites the whole file. The five old lines are gone and only "new" is left. Use `>>` when you want to keep the existing content and add to it.
:::

## Talk about it

> You've learned to create files with `touch`, write into them with `>`, append with `>>`, read them with `cat`, and copy, move, and delete them. In your own words, explain the difference between `>` and `>>`, and why `rm` is a command you should always double-check before running.

## What's next

You've been typing commands one at a time. But what if you could write a whole list of commands once and run them together, automatically, whenever you want? That's **scripting** — and it's where the command line truly becomes a superpower. On to your first Bash script.

# Navigating the Filesystem

You can now run commands — but a command always runs *somewhere*. Every shell sits inside a folder, and to work confidently you need to know where you are and how to move. This lesson gives you a mental map of the Linux filesystem and three commands that let you walk through it: `pwd`, `ls`, and `cd`. Keep the [Terminal Trainer](#terminal) open and try each command as you read.

## The filesystem is a tree

![The Linux filesystem is one tree starting at root /, with home, etc, and usr branching off; your home folder holds files like notes.txt and a projects folder](/images/lessons/linux-filesystem-tree.svg)

On Linux, every file and folder lives in one big upside-down tree. At the very top sits a single starting point called the **root**, written as a lone forward slash: `/`. Everything branches out from there — folders inside folders inside folders, all the way down to your files.

A few common branches you'll see:

```
/
├── home
│   └── you
│       ├── Documents
│       └── Downloads
├── etc
└── usr
```

`/home` holds the personal folders for each user, and inside it `/home/you` is *your* space. `/etc` keeps system settings, and `/usr` stores installed programs. You don't need to memorize them — just notice that it's all one connected tree growing down from `/`.

:::analogy
Think of the filesystem as a real building. The root `/` is the front entrance. Folders are hallways and rooms branching off from it. Files are the things sitting inside the rooms. To reach anything, you walk a path of hallways from the entrance — and the command line is simply how you walk.
:::

## Home, sweet `~`

You spend most of your time in your **home directory** — your personal room in the building, where your documents and projects live. Its full address is something like `/home/you`, but typing that constantly would be tedious, so Linux gives it a shortcut: the tilde, `~`.

Anywhere you could write `/home/you`, you can write `~` instead. That's why the prompt showed `~` earlier — it was telling you that you were standing in your home directory.

:::key
The squiggle `~` always means "my home directory." It's a shortcut you'll use constantly — `cd ~` takes you straight home from anywhere in the tree.
:::

## Where am I? `pwd`

Whenever you feel unsure where you're standing, ask:

```bash
$ pwd
/home/you
```

`pwd` — "print working directory" — prints your exact location as a full path from the root. It's the command-line equivalent of checking the "you are here" dot on a map. There's no harm in running it as often as you like.

## What's here? `ls`

To see what's inside your current folder, use `ls` (short for "list"):

```bash
$ ls
Documents  Downloads  notes.txt  photo.png
```

By itself, `ls` shows the names of everything in the current directory. But it has two flags worth knowing right away.

**Long format with `-l`** gives you details — sizes, dates, and permissions — one item per line:

```bash
$ ls -l
drwxr-xr-x  2 you you 4096 Jun 21 09:14 Documents
-rw-r--r--  1 you you   85 Jun 21 09:10 notes.txt
```

**Show hidden files with `-a`.** On Linux, any file whose name begins with a dot is hidden from the normal `ls`. The `-a` flag ("all") reveals them:

```bash
$ ls -a
.  ..  .bashrc  Documents  notes.txt
```

Notice the `.` and `..` at the start — we'll meet those in a moment.

:::tip
You can combine flags. `ls -la` (or `ls -al`) shows *all* files in *long* format at once. Stacking single-letter flags behind one dash is a common, time-saving habit on the command line.
:::

:::fill
Q: Complete the command to list files with full details.
`ls ___`
- -l *
- -a
- pwd
E: `ls -l` lists files in long format, showing sizes, dates, and permissions one per line.
:::

## Moving around: `cd`

To actually *move* into a folder, use `cd` — "change directory." Give it the name of a folder you want to step into:

```bash
$ pwd
/home/you
$ cd Documents
$ pwd
/home/you/Documents
```

You walked from your home into the Documents room. Notice how `pwd` confirms the move.

To go back **up** one level — toward the root — use the special `..`, which always means "the folder above me":

```bash
$ cd ..
$ pwd
/home/you
```

And to jump straight back **home** from anywhere, use `cd ~` (or simply `cd` with no argument, which also goes home):

```bash
$ cd ~
$ pwd
/home/you
```

:::key
Two special names show up everywhere: a single dot `.` means "the folder I'm in right now," and two dots `..` mean "the folder one level up." You'll use `cd ..` constantly to climb back out of folders.
:::

## Absolute vs relative paths

There are two ways to describe where something is, and the difference matters.

An **absolute path** starts from the root `/` and spells out the full route, like a complete street address: `/home/you/Documents`. It means the same thing no matter where you currently stand.

A **relative path** starts from wherever you happen to be right now, like "two doors down on the left." If you're in `/home/you`, then `Documents` is a relative path to `/home/you/Documents`. But that same relative name would point somewhere else if you ran it from a different folder.

```bash
# Absolute — works from anywhere
$ cd /home/you/Documents

# Relative — depends on where you start
$ cd Documents
```

:::analogy
An absolute path is like giving someone your full mailing address — it works whether they're across the street or across the world. A relative path is like saying "it's the next room over" — useful and quick, but only meaningful if the other person knows where you're standing right now.
:::

## Practice in the Terminal Trainer

This is the moment to make it real. Hop into the [Terminal Trainer](#terminal) and try a little tour:

```bash
$ pwd
$ ls
$ cd ..
$ ls
$ cd ~
$ pwd
```

Watch how `pwd` changes as you move and how `ls` shows different contents in each folder. Getting lost is part of learning — `cd ~` always brings you home, and Pixel is right there cheering each step.

:::reorder
Move into a folder, look around, then climb back up. Order these commands.
- cd Documents
- ls
- cd ..
E: Step into the folder, list what's inside, then use `cd ..` to climb back up one level.
:::

## Talk about it

> You've learned that the filesystem is one big tree growing from the root `/`, that `~` is your home, and that `.` and `..` mean "here" and "up." In your own words, explain the difference between an absolute path and a relative path — and why `cd ..` is such a handy command.

## What's next

Now you can find your way around. The natural next step is to actually *do* things to the files you find. In the next lesson you'll create, read, copy, move, and delete files — and you'll learn why one little command, `rm`, deserves real respect.

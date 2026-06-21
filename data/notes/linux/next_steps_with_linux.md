# Next Steps with Linux

You've come a long way. You can explain what Linux is, talk to the shell, navigate the filesystem, work with files, and write scripts that make decisions and loop. This final lesson ties up a few important loose ends — permissions, pipes, package managers, and environment variables — and points you toward where to go next. Think of it as the map you unfold at the end of the trail to see how far you've climbed and which paths lead onward.

## File permissions: who can do what

Remember `ls -l`? Take a closer look at the strange string of letters at the start of each line:

```bash
$ ls -l
-rwxr-xr--  1 you you  220 Jun 21 09:14 backup.sh
```

That cluster `-rwxr-xr--` is the file's **permissions**, and it answers a simple question: who is allowed to do what? Linux divides the world into three groups — the file's **owner**, its **group**, and **everyone else** — and for each group it tracks three abilities:

- **r** — read (look at the contents)
- **w** — write (change the contents)
- **x** — execute (run it as a program)

Read the string in three chunks of three. After the first character (which marks file type), `rwx` describes the owner, `r-x` the group, and `r--` everyone else. A dash means "not allowed." So in the example, the owner can read, write, and run; the group can read and run; everyone else can only read.

:::analogy
Permissions are like keys to a house. The owner has the full keyring — front door, every room, the garage. The group might have a key to the front door and living room only. Everyone else can peer through the window but not come in. `chmod` is how you hand out or take back keys.
:::

You already met the tool that changes these: `chmod`. When you ran `chmod +x script.sh` last lesson, you added the **x** (execute) permission so the file could run. That's the same system at work — you were quite literally granting your script permission to be run.

:::key
The `-rwxr-xr--` column from `ls -l` shows permissions for three groups — owner, group, and others — using **r**ead, **w**rite, and e**x**ecute. `chmod` adds or removes these. `chmod +x` is the one you'll reach for most, to make scripts runnable.
:::

## Pipes: connecting commands together

Here's one of the most beautiful ideas in all of Linux. The vertical bar `|`, called a **pipe**, takes the output of one command and feeds it straight into the next as input. Small commands snap together like building blocks into something bigger.

A perfect partner for pipes is `grep`, which filters text, keeping only the lines that contain a pattern you name. Combine them:

```bash
$ ls -l | grep ".sh"
-rwxr-xr--  1 you you  220 Jun 21 09:14 backup.sh
-rwxr-xr--  1 you you  118 Jun 21 09:12 greet.sh
```

Read it as a sentence: "list everything in long format, then pipe that list into grep, which keeps only the lines mentioning `.sh`." Neither command knew about the other — the pipe simply passed the result along.

:::analogy
A pipe is an assembly line. The first station produces something and slides it down the belt to the next station, which does its part and passes it on. Each worker is simple and does one job well; the power comes from chaining them. This is the heart of the Unix philosophy: small tools, combined.
:::

## Redirection recap

You already learned redirection's cousins. It's worth holding all three in one place:

- `>` sends output *into a file*, overwriting it.
- `>>` sends output into a file, appending to the end.
- `|` sends output into another *command*.

The first two save results to disk; the pipe keeps results flowing between programs. Together they let you capture, transform, and route data however you need.

:::tip
A handy combo: filter with a pipe, then save with redirection. For example, `ls -l | grep ".sh" > scripts.txt` finds every script and writes the list into a file in one line. You already know every piece of that — you just snapped them together.
:::

## Package managers: installing software

So far you've used commands that come built in. But how do you *get* new programs? On Linux you rarely download installers from random websites. Instead you use a **package manager** — a trusted system that fetches, installs, and updates software for you from official sources.

Which one you use depends on your distro. On Ubuntu and Debian it's `apt`:

```bash
$ sudo apt update
$ sudo apt install git
```

The first line refreshes the catalog of available software; the second installs a program (here, `git`). The `sudo` prefix means "do this as the superuser" — installing software is a system-wide change, so it needs elevated permission. Other distros have their own managers (Fedora uses `dnf`, Alpine uses `apk`), but the idea is identical: ask, and the manager safely fetches it for you.

:::warning
`sudo` gives a command full administrative power over the system — it can change or delete anything. Use it deliberately, only when a task genuinely needs it, and read the command twice before pressing Enter. Great power, as ever, asks for great care.
:::

## Environment variables and PATH

In the scripting lesson you made your own variables. The system keeps a set of its own, called **environment variables**, that shape how your whole session behaves. You can peek at one with `echo`:

```bash
$ echo $HOME
/home/you
$ echo $PATH
/usr/local/bin:/usr/bin:/bin
```

The most important is `PATH`. It's a list of folders, separated by colons, where the shell looks when you type a command. When you type `ls`, the shell walks through each folder in `PATH` until it finds a program called `ls`. This is exactly why you needed `./script.sh` earlier — your script wasn't in any `PATH` folder, so you had to point directly at it with `./`.

:::key
`PATH` is the list of folders the shell searches to find commands. If a program isn't in one of those folders, the shell can't find it by name — which is why you run local scripts with `./` to point right at them.
:::

## Project ideas to keep growing

The best way to cement everything is to build something small and real. A few ideas sized just right for where you are:

- **A file organizer.** Write a script that loops over files and moves them into folders by type — images into one folder, text files into another. You know `for`, `mv`, and `if`; that's all it takes.
- **A backup script.** Make a script that copies an important folder into a timestamped backup. You've used `cp`, variables, and `echo` — combine them.
- **A daily greeting.** Expand your `welcome.sh` to read input, make decisions with `if`, and report something useful, like the current folder's contents.

Start tiny. A five-line script that does one real thing teaches you more than a hundred lines you only read about.

:::example
Try this in the [Terminal Trainer](#terminal) right now as a warm-up: make a folder, drop a couple of files in it, then list and filter them.

```bash
$ mkdir practice
$ cd practice
$ touch notes.txt photo.png script.sh
$ ls
```

From here, imagine the script you'd write to sort those three files by type. You already have every command you'd need.
:::

## Where to go next

You've finished the Linux and Bash path — genuinely well done. The skills you now hold are the same ones professional developers use every single day to deploy code, manage servers, and automate the boring parts of their work.

Keep the [Terminal Trainer](#terminal) close and replay its missions until the commands feel like second nature. Then carry these habits into the rest of your journey: when you learn a new language or tool, you'll almost always meet it through a terminal, and now that terminal feels like home. Pixel is proud of how far you've come — and this is only the start of what you can build.

## Talk about it

> Look back across the whole path: operating systems, the shell, navigating folders, working with files, and scripting. Which idea clicked most for you, and which one would you like to practice more? Pick one project from this lesson and describe, in plain words, the steps you'd take to build it.

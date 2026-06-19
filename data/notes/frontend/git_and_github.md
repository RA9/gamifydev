# Git and GitHub

You've been building projects — but what happens when you make a change that breaks everything and can't undo it? Or when you want to share your work, or collaborate? The answer to all of these is **Git**, the tool professional developers use every single day.

By the end of this lesson you'll understand what Git does, the core workflow, and how GitHub fits in.

## Git vs GitHub — not the same thing

People mix these up constantly, so let's be clear:

- **Git** is a *version control system* — software that runs on your computer and tracks every change to your project over time. It works completely offline.
- **GitHub** is a *website* that hosts Git projects online, so you can back them up, share them, and collaborate.

:::analogy
Git is like the "save history" in a video game — checkpoints you can always return to. GitHub is the cloud save that lets you play from any device and share your progress with friends.
:::

## Why version control matters

Without Git, "saving" means overwriting your file and losing what was there before. With Git you create **commits** — labelled snapshots of your whole project. You can:

- go back to any previous snapshot
- see exactly what changed, and when
- experiment on a **branch** without risking your working version
- collaborate without emailing `final_v3_REALLY_final.zip` around

## The core workflow

This is the loop you'll repeat thousands of times. Your changes travel through four places:

![The Git workflow: working directory, staging, local repository, GitHub](/images/lessons/git-workflow.svg)

1. **Working directory** — the files as you're editing them right now.
2. `git add` — moves your chosen changes to the **staging area** (a "ready to save" shelf).
3. `git commit` — saves a snapshot of the staged changes into your **local repository**, with a message.
4. `git push` — uploads your commits to **GitHub**.

```bash
git add index.html        # stage one file (or: git add .)
git commit -m "Add homepage hero"   # snapshot it, with a message
git push                  # send commits up to GitHub
```

:::tip
Write commit messages that finish the sentence "This commit will…" — like "Add login form" or "Fix navbar spacing". Your future self (and teammates) will thank you when scanning the history.
:::

:::quiz
Q: Which command saves a snapshot of your staged changes to your local repository?
- `git add`
- `git commit` *
- `git push`
E: `git commit` records a snapshot (with a message). `git add` only stages changes; `git push` uploads existing commits to GitHub.
:::

## Getting started (the first time)

A typical first-project flow looks like this:

```bash
git init                  # start tracking this folder with Git
git add .                 # stage everything
git commit -m "First commit"
# create an empty repo on GitHub, then:
git remote add origin https://github.com/you/my-site.git
git push -u origin main
```

After that first setup, your day-to-day is just the `add → commit → push` loop.

:::quiz
Q: What's the difference between Git and GitHub?
- They're two names for the same thing
- Git tracks versions on your computer; GitHub hosts projects online *
- Git is online; GitHub is offline
E: Git is the version-control tool running locally; GitHub is the online service that hosts Git repositories for backup and collaboration.
:::

:::key
**Git** tracks the history of your project in snapshots called **commits**. The everyday loop is **add → commit → push**: stage changes, snapshot them locally, then upload to **GitHub** to back up and share.
:::

## Talk about it

Explain out loud:

> "What's the difference between Git and GitHub, and what do `add`, `commit`, and `push` each do?"

If you can walk a friend through the workflow, you've got the foundation every developer relies on.

## What's next

Now that you can save and share work, it's time to build something to show off: your own **portfolio website**.

# Git and GitHub

You've been building projects — but what happens when a change breaks everything and you can't undo it? Or when you want to back up your work, share it, or collaborate? The answer to all of these is **Git**, the tool professional developers use every single day.

By the end of this lesson you'll understand what version control is, the core Git workflow with real commands, and how GitHub fits in.

## What version control is — and why you want it

**Version control** is a system that records changes to your files over time so you can recall any earlier version, see exactly what changed, and work without fear of losing anything.

Without it, "saving" means overwriting a file and losing what was there before. People end up with folders like `site_final_v2_REALLY_final_v3.zip`. Version control replaces all of that with a clean, searchable history.

With Git you create **commits** — labelled snapshots of your whole project. You can:

- go back to any previous snapshot
- see exactly what changed, when, and why
- experiment on a **branch** without risking your working version
- collaborate with others without emailing files around

## Git vs GitHub — not the same thing

People mix these up constantly, so let's be clear.

- **Git** is a *version control system* — software that runs on your computer and tracks every change to your project. It works completely offline.
- **GitHub** is a *website* that hosts Git projects online so you can back them up, share them, and collaborate. (GitLab and Bitbucket are alternatives.)

:::analogy
Git is like the save-history in a video game — checkpoints you can always return to. GitHub is the cloud save that lets you play from any device and share your progress with friends.
:::

## One-time setup

Install Git from [git-scm.com](https://git-scm.com), then tell Git who you are. You only do this once per computer:

```bash
git config --global user.name "Ada Lovelace"
git config --global user.email "ada@example.com"
```

These details get stamped on every commit you make. Check that it worked:

```bash
git config --list
git --version
```

:::tip
Use the same email here that you'll use for your GitHub account. That way GitHub can link your commits to your profile.
:::

## The core local workflow

This is the loop you'll repeat thousands of times. Your changes travel through three places: your **working directory** (files as you edit them), the **staging area** (a "ready to save" shelf), and the **repository** (the saved history).

**Start tracking a project.** Inside your project folder, run:

```bash
git init
```

This creates a hidden `.git` folder — your project is now a Git repository.

**Check what's going on.** `git status` is the command you'll run most. It tells you what's changed and what's staged:

```bash
git status
```

**Stage your changes.** Move the files you want to save onto the staging shelf:

```bash
git add index.html        # stage one specific file
git add .                  # stage everything that changed
```

**Commit the snapshot.** Save the staged changes into history with a message:

```bash
git commit -m "Add homepage hero section"
```

**Review your history.** See the list of commits you've made:

```bash
git log
git log --oneline          # compact, one line per commit
```

:::example
A complete first session from scratch:

```bash
git init
git add .
git commit -m "Initial commit: basic HTML and CSS"
# ...edit some files...
git add .
git commit -m "Add navigation bar and footer"
git log --oneline
```
:::

## Writing good commit messages

A commit message explains *what changed and why*. Future-you will thank present-you.

- Write in the present tense, as a command: "Add login form," not "Added" or "Adding."
- Keep the first line short (around 50 characters) and specific.
- Describe the change, not the obvious: "Fix broken nav link on mobile" beats "update."

```bash
git commit -m "Fix footer overlap on small screens"   # good
git commit -m "stuff"                                  # useless later
```

## Ignoring files with .gitignore

Some files should never be committed — system junk, secrets, huge build folders. Create a file named `.gitignore` in your project root and list patterns to skip:

```bash
# .gitignore
node_modules/
.DS_Store
.env
dist/
*.log
```

Git will then pretend those files don't exist for tracking purposes. Add `.gitignore` *before* your first commit so the junk never gets in.

:::warning
Never commit secrets — API keys, passwords, tokens. Once something is committed and pushed, it lives in the history even if you delete it later. Put secrets in a `.env` file and ignore that file.
:::

## Branches — a quick intro

A **branch** lets you work on something new without touching your main version. The default branch is usually called `main`.

```bash
git branch                 # list branches; * marks the current one
git switch -c new-feature  # create a new branch AND switch to it
```

Now you can commit freely on `new-feature`. When the work is good, you **merge** it back into `main`:

```bash
git switch main            # go back to main
git merge new-feature      # bring the feature's commits into main
```

:::analogy
A branch is like writing on a photocopy of your document. You can scribble all over the copy; if it works out, you fold those edits back into the original. If not, you just throw the copy away.
:::

## GitHub — putting it online

A **remote** is a copy of your repository hosted elsewhere — on GitHub. Here's how to connect a local project to a new GitHub repo.

1. On GitHub, click **New repository**. Give it a name and **don't** add a README (you already have local files).
2. GitHub shows you a URL like `https://github.com/yourname/yourrepo.git`.
3. Connect your local repo to it and push your work up:

```bash
git remote add origin https://github.com/yourname/yourrepo.git
git branch -M main         # make sure your branch is named main
git push -u origin main    # push and remember this remote+branch
```

The `-u` means "set this as the default," so next time you can just type `git push`.

**Cloning** copies an existing GitHub repo down to your machine:

```bash
git clone https://github.com/someone/cool-project.git
```

**Pulling** brings down changes others (or you, from another computer) pushed:

```bash
git pull
```

## Your everyday loop

Once set up, your daily rhythm is short and sweet:

```bash
git pull                          # get the latest (if collaborating)
# ...do your work...
git status                        # see what changed
git add .                         # stage it
git commit -m "Describe the change"
git push                          # send it to GitHub
```

:::quiz
Q: What is the difference between Git and GitHub?
- They are two names for the same program
- Git is version-control software on your computer; GitHub is a website that hosts Git repos online *
- Git is online and GitHub runs locally
E: Git is the local tool that tracks changes (and works offline). GitHub is a hosting service where you store and share those Git repositories.
:::

:::quiz
Q: Which command saves a snapshot of your staged changes into the project history?
- git add
- git commit -m "message" *
- git status
E: `git add` stages changes onto the shelf, and `git commit` records those staged changes as a permanent snapshot with a message. `git status` only reports state.
:::

:::fill
To upload your local `main` branch to GitHub for the first time and remember the connection, you run `git push ___ origin main`.
- -u *
- -m
- --all
E: `git push -u origin main` sets `origin/main` as the upstream so future pushes can be just `git push`.
:::

## Recap

- **Version control** records your project's history so you can undo, compare, and collaborate safely.
- **Git** is local software; **GitHub** is an online host for Git repositories.
- Configure your identity once with `git config --global`.
- The core loop: `git init` -> edit -> `git add` -> `git commit -m` -> `git log`.
- Write clear, present-tense commit messages, and use `.gitignore` to skip junk and secrets.
- Branches (`git branch`, `git switch -c`, `git merge`) let you work safely in parallel.
- Connect to GitHub with `git remote add origin`, push with `git push -u origin main`, and use `git clone` / `git pull` to download changes.

**Next up:** Deploying Your Website — getting your project onto a real, public URL.

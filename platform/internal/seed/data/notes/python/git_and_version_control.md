# Git and Version Control

You just spent three hours getting a feature working. You "clean up" one file, and suddenly everything is broken — and you have no idea what you changed. Without version control, your only option is to remember. With it, you rewind to the exact working state in one command. This lesson teaches you Git: the safety net every developer relies on, and the tool that makes teamwork possible.

## What Version Control Is and Why You Need It

**Version control** is a system that records the history of your project over time. Every time you save a meaningful checkpoint, it remembers the *entire* state of your files at that moment — and lets you jump back to any of them later.

Think of the save system in a game. You don't just have one save slot that gets overwritten; you have a timeline of saves, and if a boss fight goes badly you reload an earlier one. Git is that, for your code.

:::analogy
Writing code without version control is like playing a hard game with no save points. One wrong move and you start the whole level over. Git gives you a save point after every room — mess up, and you reload the last good one instead of the beginning.
:::

It solves three problems at once. **History:** you can see exactly what changed, when, and why. **Collaboration:** many people can work on the same project without emailing `final_v2_REALLY_final.py` back and forth. **Safety:** experiments are free — try a wild idea on a branch, and if it fails, throw it away with the original untouched.

**Git** is the most widely used version control system. It runs on your machine, tracking a single project folder called a **repository** (or "repo").

## The Core Loop

Ninety percent of daily Git is four commands in a rhythm: check status, stage changes, commit them, review history. Let's walk through starting a fresh project.

First, turn a folder into a repo with `git init`:

```bash
mkdir dungeon-crawler
cd dungeon-crawler
git init
```

Now create a file and check where things stand with `git status`:

```bash
echo "print('Welcome, adventurer!')" > game.py
git status
```

Git reports `game.py` as **untracked** — it sees the file but isn't recording it yet. To tell Git "I want this in my next checkpoint," you **stage** it with `git add`:

```bash
git add game.py
```

Staging is like putting items in a chest before you seal it — you choose exactly what goes in. Now **commit**: seal the chest and label it. A commit is a permanent snapshot with a message describing it.

```bash
git commit -m "Add starting game script"
```

Make more changes, and the loop repeats: `status` to see what moved, `add` to stage, `commit` to snapshot. To see your history, use `git log`:

```bash
git log --oneline
```

```text
a1b2c3d Add starting game script
```

:::key
The everyday loop is: **edit → `git status` → `git add` → `git commit`**. Staging (`add`) and committing (`commit`) are two separate steps on purpose — it lets you snapshot only *some* of your changes and leave the rest for a later commit.
:::

## Writing Good Commit Messages

A commit message is a note to your future self and your teammates. Six months from now, `git log` is the only record of *why* you made a change. "fixed stuff" tells you nothing; "Fix crash when inventory is empty" tells you everything.

The widely-followed convention: write the summary line in the **imperative mood**, as if completing the sentence "This commit will...". Keep it under about 50 characters, and capitalize it with no period.

```text
Add health bar to player HUD
Fix off-by-one error in level loader
Remove unused sprite assets
```

Not:

```text
added some stuff
changes
asdfasdf
fixed the thing finally omg
```

If a change needs more explanation, leave the summary line short, add a blank line, then write a longer body:

```bash
git commit -m "Fix save corruption on quit" -m "The save was written on a background thread that could be killed mid-write. Now we flush and join the thread before exiting."
```

:::tip
A good test: read your message as "This commit will [your message]." If "This commit will fixed stuff" sounds wrong, so does your message. "This commit will fix crash on empty inventory" reads clean — that's the imperative mood.
:::

## Branching and Merging

So far everything lives on one timeline, called `main` by default. But what if you want to try a risky new feature without breaking the working game? You make a **branch** — a parallel copy of the timeline where you can experiment freely.

Create and switch to a branch in one step with `git switch -c` (the modern command; older tutorials use `git checkout -b`, which does the same thing):

```bash
git switch -c add-boss-fight
```

You are now on the `add-boss-fight` branch. Commit away — add enemies, tweak logic, break things. Meanwhile `main` sits untouched and playable. Switch back any time:

```bash
git switch main
```

When the feature is done and working, you fold it back into `main` with **merge**. From the branch you want to merge *into* (`main`), you merge the feature branch in:

```bash
git switch main
git merge add-boss-fight
```

Git combines the two timelines. If both branches changed *different* files, this is seamless. If they changed the *same lines*, Git pauses and asks you to resolve a **merge conflict** — it marks the clashing spots so you can pick which version wins. Conflicts feel scary at first but are routine; you just edit the file to the version you want and commit.

:::analogy
A branch is a "what if" save file. You copy your progress, go try the dangerous shortcut, and if it works you keep it — if it doesn't, your main save never knew it happened. Merging is bringing the successful run back into your main timeline.
:::

## Remotes and GitHub

Everything so far lives only on your computer. A **remote** is a copy of your repo hosted elsewhere — most commonly on **GitHub**. Remotes are how you back up your work and how teams share a single source of truth.

If a project already exists on GitHub, you **clone** it to get a full local copy:

```bash
git clone https://github.com/someuser/dungeon-crawler.git
```

Once connected to a remote (conventionally named `origin`), two commands move commits back and forth. **Push** sends your local commits up to GitHub:

```bash
git push origin main
```

**Pull** brings down commits others have added, merging them into your local copy:

```bash
git pull origin main
```

The rhythm on a team: `pull` before you start (so you have the latest), do your work, `commit`, then `push` to share it. Pull often to avoid drifting far from your teammates and piling up conflicts.

When you want to propose your branch's changes be merged into the shared `main`, you open a **pull request** (PR) on GitHub. A PR is a "please review and merge my branch" request — teammates comment, suggest fixes, and approve before it lands. It is where code review happens, and it's the heart of how teams collaborate on GitHub.

## `.gitignore` — Keep Junk (and Secrets) Out

Not everything in your project folder belongs in the repo. Python leaves behind `__pycache__/` folders and `.pyc` files. Your virtual environment (`venv/`) is huge and machine-specific. And some files — like passwords or API keys — must *never* be committed at all.

A `.gitignore` file lists patterns Git should pretend it doesn't see. Create one at the root of your repo:

```text
# Python cruft
__pycache__/
*.pyc

# Virtual environments
venv/
.venv/

# Secrets and local config
.env
secrets.json

# OS and editor files
.DS_Store
.vscode/
```

Anything matching these patterns stays on your disk but never gets staged or committed. Add `.gitignore` itself to the repo so everyone on the team shares the same rules.

:::warning
Never commit secrets — passwords, API keys, tokens, database URLs. Once you push a secret to GitHub, treat it as leaked *forever*: even if you delete it in a later commit, it lives on in the history where anyone can dig it out. Put secrets in a `.env` file, add `.env` to `.gitignore` on day one, and if you ever do push a key by accident, **rotate it** (generate a new one and revoke the old) immediately.
:::

Committing your `venv/` or `__pycache__/` won't leak anything dangerous, but it bloats the repo and causes noisy, meaningless diffs. Anyone who clones your project rebuilds those from your requirements file — they don't need your copy.

## Practice

:::quiz
Q: You just wrote a database password directly into `config.py` to test locally. What should you do before committing?
- Commit it now and delete it in the next commit
- Move the secret to a `.env` file and add `.env` to `.gitignore` *
- Push it so teammates can use the same password
- Nothing — private repos are safe to store secrets in
E: Never commit secrets. Deleting them later doesn't help — they stay in the history. Keep the secret in a `.env` file that `.gitignore` excludes, so it never enters the repo at all.
:::

:::predict
Q: You run `git add game.py` but forget to run `git commit`. What has happened to your changes?
- They are permanently saved to the project history
- They are staged, ready for the next commit, but not yet snapshotted *
- They are pushed to GitHub automatically
- They are lost when you close the terminal
E: `git add` only *stages* changes — it puts them in line for the next commit. Nothing becomes a permanent snapshot until you run `git commit`. Staging and committing are deliberately two separate steps.
:::

## Recap

- Version control records your project's history so you can review changes, collaborate, and safely undo mistakes — a save system for your code.
- The daily loop is **edit → `git status` → `git add` → `git commit`**; `git log` shows the history.
- Write commit messages in the imperative mood ("Fix crash on empty inventory"), short and descriptive.
- Branch with `git switch -c name` to experiment safely; `git merge` folds a finished branch back into `main`.
- A remote like GitHub backs up and shares your work: `git clone` to copy, `git push` to send commits, `git pull` to receive them; pull requests are where teammates review changes.
- Use `.gitignore` to keep `__pycache__/`, `venv/`, and OS clutter out of the repo — and **never commit secrets**; if one leaks, rotate it immediately.

**Next up:** Virtual Environments and Packaging — isolating your project's dependencies so "it works on my machine" becomes "it works everywhere."

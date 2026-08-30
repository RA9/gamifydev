# Git and GitHub

If you're serious about becoming a developer, Git stops being optional immediately.

It is not just a backup tool. It is how you:

- record the history of a project
- experiment safely
- collaborate with other people
- review work before it ships
- recover from mistakes without panic

By the end of this lesson, you should understand the real day-to-day Git workflow, not just a few commands copied from a tutorial.

## Git vs GitHub

These are related, but not the same thing.

- **Git** is the version-control software running on your machine.
- **GitHub** is a hosting platform for Git repositories.

Git tracks the history.
GitHub stores and shares that history online.

:::analogy
Git is your save system. GitHub is the shared online vault where those saves can be backed up, reviewed, and collaborated on.
:::

## Why version control matters

Without version control, mistakes feel expensive.

You end up with folders like:

- `portfolio-final`
- `portfolio-final-2`
- `portfolio-final-actual-final`

Git replaces that chaos with commits — named snapshots in a structured history.

A good Git workflow gives you:

- confidence to experiment
- a clear record of changes
- easier debugging when something breaks
- cleaner collaboration when more than one person touches the codebase

## The four core areas in Git

When people first learn Git, the hardest part is understanding where changes live.

Think in four places:

1. **Working directory** — your actual files as you edit them
2. **Staging area** — the changes you've selected for the next snapshot
3. **Local repository** — the commit history on your machine
4. **Remote repository** — the copy on GitHub

A normal workflow is:

- edit files
- check status
- stage selected changes
- commit them
- push them to GitHub

## First-time setup

On a new machine, configure your identity once:

```bash
git config --global user.name "Ada Lovelace"
git config --global user.email "ada@example.com"
```

Check it:

```bash
git config --list
git --version
```

:::tip
Use the same email address as your GitHub account if you want your commits connected to your profile.
:::

## Starting a repository

Inside a project folder:

```bash
git init
```

That creates the `.git` directory and turns the folder into a repository.

Then check what Git sees:

```bash
git status
```

`git status` is one of the most important commands you will ever learn.

## The everyday command loop

### Stage changes

```bash
git add index.html
git add css/style.css
git add .
```

Use targeted adds when you want precise commits. Use `git add .` when everything currently changed belongs together.

### Commit changes

```bash
git commit -m "Add hero section and CTA styles"
```

A commit should represent one meaningful unit of work.

### Push to GitHub

```bash
git push
```

That sends your local commits to the remote repository.

## Connecting to GitHub

After creating a new empty repo on GitHub, connect your local repo to it:

```bash
git remote add origin https://github.com/yourname/portfolio.git
git branch -M main
git push -u origin main
```

What those do:

- `remote add origin ...` connects your local repo to GitHub
- `branch -M main` ensures the branch is called `main`
- `push -u origin main` uploads and remembers the upstream branch

After that, `git push` and `git pull` are usually enough.

## Reading history

A professional workflow includes checking history, not just creating it.

```bash
git log
git log --oneline
git diff
git diff --staged
```

Use these to answer:

- what changed?
- what is staged right now?
- what did I commit earlier?

## Writing useful commit messages

Good commit messages make your project history readable.

Strong examples:

- `Add responsive navigation layout`
- `Fix quiz score reset bug`
- `Refactor task rendering into reusable function`

Weak examples:

- `update`
- `stuff`
- `changes`

A good commit message should tell a reviewer what this snapshot *does*.

## .gitignore is part of a professional setup

Some files should never be committed.

Common examples:

```text
node_modules/
.env
.DS_Store
dist/
*.log
```

Put them in `.gitignore`.

Why this matters:

- avoids giant unnecessary files
- avoids committing secrets
- keeps the repo clean and reviewable

:::warning
Never commit secrets. If an API key, token, or password reaches Git history and gets pushed, deleting the file afterward does not erase the original exposure.
:::

## Branches let you work safely

Branches are one of the most important ideas in Git.

Create and switch to a new branch:

```bash
git switch -c feature/contact-form
```

Now you can make commits without affecting `main` directly.

When the work is ready:

```bash
git switch main
git merge feature/contact-form
```

This is how real teams isolate features, bug fixes, and experiments.

## Pull requests and code review

On GitHub, the usual team workflow is:

1. create a branch locally
2. push the branch to GitHub
3. open a **pull request**
4. review the changes
5. merge when approved

A pull request is more than a merge button. It is a discussion space around code.

A strong PR includes:

- a clear title
- what changed
- why it changed
- screenshots if the UI changed
- notes for reviewers

Even on solo projects, learning this workflow is worth it because it mirrors real engineering teams.

## Merge conflicts are normal

Sometimes Git can't merge automatically because the same lines changed in different places.

That is a **merge conflict**.

The workflow is:

- open the conflicted file
- choose the correct final version
- remove Git's conflict markers
- stage the resolved file
- commit the resolution

Conflicts are not proof you broke Git. They are a normal part of collaboration.

## Practical solo-developer workflow

For a personal project, a strong routine looks like this:

```bash
git pull
# do the work
git status
git diff
git add .
git commit -m "Describe the change clearly"
git push
```

If the change is risky or larger:

```bash
git switch -c feature/dashboard-api
# work, commit, push
```

That habit alone will make you look much more professional.

## GitHub profile and repo hygiene

If you're building public work:

- add a clear README
- use meaningful repo names
- keep old experimental repos private or archived if they are noisy
- add screenshots for UI projects
- pin your strongest repos on your GitHub profile

GitHub is part of your public developer presence.

:::quiz
Q: What is the main value of a Git branch?
- It makes CSS load faster
- It lets you work on changes without risking your main version *
- It stores passwords securely
E: Branches isolate work so you can experiment, review, and merge intentionally instead of editing main directly.
:::

:::quiz
Q: What is the purpose of `git add`?
- Upload code to GitHub
- Move selected changes into the staging area *
- Delete tracked files
E: `git add` prepares selected changes for the next commit. It does not create a commit or push anything online by itself.
:::

:::fill
To inspect changes you have made but have **not staged yet**, use `git ___`.
- diff *
- log
- clone
E: `git diff` shows unstaged changes in your working directory.
:::

## What good looks like

You should now be able to:

- explain the difference between Git and GitHub
- initialize and connect a repository
- use status, add, commit, push, pull, and log confidently
- understand branches and pull requests
- write useful commit messages
- avoid common beginner mistakes with ignored files and secrets

## What's next

In **Deploying Your Website**, you'll take that version-controlled project and ship it to a real public URL. After that, you'll use Git and GitHub again to publish and maintain your portfolio.
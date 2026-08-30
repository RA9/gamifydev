# Git Best Practices

Knowing Git commands isn't enough. The difference between a developer who *uses* Git and one who uses it *well* comes down to habits — small disciplines that keep your project history clean, your secrets safe, and your team sane.

## Atomic commits — one logical change

An **atomic commit** captures exactly one logical change. Not "everything I did today" and not "half a feature."

**Good — one logical change per commit:**

```bash
git commit -m "Add responsive navigation menu"
git commit -m "Fix quiz score not resetting on restart"
git commit -m "Update footer links and copyright year"
```

**Bad — too much in one commit:**

```bash
git commit -m "Add nav, fix score bug, update footer, refactor CSS"
```

**Bad — too granular:**

```bash
git commit -m "Add nav container div"
git commit -m "Add first nav link"
git commit -m "Add second nav link"
```

:::key
An atomic commit should make sense on its own. If you need to revert it later, only one logical change is undone. If you're reviewing history, each commit tells a clear story.
:::

## Good commit messages — imperative mood, 50-char subject

Write commit messages as if you're completing the sentence: "If applied, this commit will..."

```text
If applied, this commit will... Add responsive navigation menu ✅
If applied, this commit will... Fix quiz score reset bug      ✅
If applied, this commit will... Updated the thing             ❌ (past tense)
If applied, this commit will... stuff                         ❌ (meaningless)
```

**The format:**

```text
Subject line: 50 chars or less, imperative mood, no period
                                                          (blank line)
Optional body: wrap at 72 chars, explain WHY not WHAT
```

**Strong examples:**

```bash
git commit -m "Add keyboard shortcuts for quiz options (1-4 keys)"
git commit -m "Fix hero image overflow on mobile viewports"
git commit -m "Refactor renderQuestion to accept a question object"
```

**Weak examples:**

```bash
git commit -m "updates"
git commit -m "Fixed stuff"
git commit -m "WIP"
git commit -m "asdfasdf"
```

:::tip
If you can't write a clear commit message, your commit probably does too many things. Split it into smaller, focused commits.
:::

## .gitignore essentials

A `.gitignore` file tells Git which files and folders to never track. Create it at the root of your project:

```text
# Dependencies
node_modules/

# Environment variables (secrets!)
.env
.env.local

# OS-generated files
.DS_Store
Thumbs.db

# Build output
dist/
build/

# Editor files
.vscode/settings.json
*.swp

# Logs
*.log
npm-debug.log*
```

### Why each matters

| Entry | Why ignore it |
|-------|--------------|
| `node_modules/` | Thousands of dependency files — reinstall with `npm install` |
| `.env` | Contains API keys, passwords, and secrets |
| `.DS_Store` | macOS metadata — useless to everyone else |
| `dist/` | Build output — regenerate with `npm run build` |
| `*.log` | Debug logs — not part of source code |

:::warning
Add `.gitignore` as one of the first things you do in a new project — *before* your first commit. If you commit `node_modules` or `.env` and then add them to `.gitignore`, the files are already in Git history. Ignoring them afterward doesn't erase the past.
:::

## Never commit secrets

This rule is absolute. API keys, tokens, passwords, and private credentials must never appear in a Git commit.

**Wrong — secret hard-coded in source:**

```js
const API_KEY = "sk-abc123-real-secret-key";
fetch(`https://api.example.com/data?key=${API_KEY}`);
```

**Right — secret in .env (which is gitignored):**

```text
# .env file (not committed)
API_KEY=sk-abc123-real-secret-key
```

```js
// In your code, reference it from the environment
const API_KEY = process.env.API_KEY;
```

If a secret reaches GitHub (even briefly), consider it compromised. Rotate the key immediately.

:::key
Git history is permanent. Deleting a file doesn't erase it from the commit that added it. If a secret has ever been committed and pushed, it's exposed — even if you delete the file in the next commit.
:::

## Branch naming conventions

Use prefixes to categorize branches. This makes the branch list scannable and the purpose of each branch obvious.

```bash
feature/quiz-timer        # new functionality
feature/responsive-nav    # new functionality

fix/score-reset-bug       # bug fix
fix/image-overflow        # bug fix

chore/update-readme       # maintenance, non-feature work
chore/compress-images     # cleanup

refactor/quiz-state       # code restructuring, no behavior change
```

**Rules:**
- Lowercase, hyphen-separated
- Start with a category prefix
- Keep it short but descriptive
- No spaces, no special characters

## Keeping main deployable

The `main` branch should always represent a working, shippable state. This is the golden rule of professional Git workflows.

**What this means:**
- Never commit broken code directly to `main`
- Use feature branches for all work-in-progress
- Merge only tested, reviewed code
- If something breaks on `main`, fix it immediately (highest priority)

```bash
# Bad: committing half-finished work directly to main
git switch main
git add .
git commit -m "WIP - not done yet"  # 🚫

# Good: work on a branch, merge when ready
git switch -c feature/quiz-timer
git add .
git commit -m "Add timer countdown UI"
# ... finish the feature, test it, then merge
```

:::tip
A deployable `main` means you can ship at any time. If a client, boss, or recruiter asks to see your project, the live site matches the `main` branch and it works perfectly.
:::

## git stash basics

Sometimes you're in the middle of work and need to switch branches, but you don't want to commit half-finished changes. `git stash` saves your uncommitted changes temporarily.

```bash
# Save current changes to the stash
git stash

# Your working directory is now clean — switch branches freely
git switch main
# do something on main
git switch feature/quiz-timer

# Bring your stashed changes back
git stash pop
```

**Common stash commands:**

```bash
git stash              # stash all uncommitted changes
git stash pop          # apply the most recent stash and remove it
git stash list         # see all stashed changes
git stash drop         # discard the most recent stash
git stash -m "WIP timer logic"  # stash with a description
```

:::warning
Don't use stash as long-term storage. Stashed changes are easy to forget about and can become confusing. Stash for quick context switches (minutes to hours), not days.
:::

## git diff before committing

Always review what you're about to commit:

```bash
# See unstaged changes
git diff

# See staged changes (what will be in the next commit)
git diff --staged

# See a summary of changed files
git status
```

This catches:
- Debug `console.log` statements you forgot to remove
- Unintentional whitespace changes
- Files you didn't mean to include

Make it a habit: `git diff --staged` → read it → `git commit`.

## The professional Git routine

```bash
# Start of work
git switch main
git pull
git switch -c feature/new-work

# During work
git status                    # what changed?
git diff                      # review changes
git add specific-file.js      # stage intentionally
git commit -m "Clear message" # atomic commit

# Ready to merge
git push -u origin feature/new-work
# Open PR on GitHub, review, merge

# After merge
git switch main
git pull
git branch -d feature/new-work
```

:::quiz
Q: Why should commit messages use imperative mood ("Add feature") instead of past tense ("Added feature")?
- Imperative mood matches Git's own conventions and completes the sentence "If applied, this commit will..." *
- Past tense is grammatically incorrect
- Imperative mood makes commits run faster
- Past tense confuses the terminal
E: Git's own generated messages use imperative mood ("Merge branch..."). Following this convention keeps the history consistent and each message reads as an instruction — what the commit *does* when applied.
:::

:::quiz
Q: What is the risk of committing an API key to a Git repository?
- The API will stop working
- The key is permanently in Git history even if you delete the file, making it accessible to anyone who clones the repo *
- Git will reject the commit
- The key will be automatically rotated
E: Git history is append-only. A committed secret exists in the commit that added it, even after the file is deleted in a later commit. Anyone with repo access (or who clones a public repo) can find it.
:::

## Recap

- **Atomic commits** capture one logical change — not too big, not too small.
- **Commit messages** use imperative mood, 50 characters or fewer, and describe *what* the commit does.
- **`.gitignore`** prevents `node_modules`, `.env`, `.DS_Store`, and build output from being tracked.
- **Never commit secrets.** Use `.env` files and add them to `.gitignore`.
- **Branch naming** uses prefixes: `feature/`, `fix/`, `chore/`, `refactor/`.
- **Keep `main` deployable** — only merge tested, working code.
- **`git stash`** temporarily saves uncommitted work for quick branch switches.
- **`git diff --staged`** before committing catches debug logs and accidental changes.

**Next up:** Deployment Platforms and CI Basics — shipping your code to the world and automating the process.

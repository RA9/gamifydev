# Branching Merging and Pull Requests

Branches are how professional developers work on features, fix bugs, and experiment — all without risking the stable `main` branch. This lesson covers the full branching workflow, from creating a branch to merging a pull request.

## Creating a branch with git switch

A branch is a parallel line of development. You branch off from `main`, make changes, and merge back when the work is ready.

```bash
# Create and switch to a new branch
git switch -c feature/quiz-timer
```

The `-c` flag means "create." Without it, `git switch` switches to an *existing* branch.

```bash
# Switch to an existing branch
git switch main

# Switch back
git switch feature/quiz-timer
```

:::tip
`git switch` is the modern replacement for `git checkout`. It does less (only switches branches, doesn't restore files), which makes it safer and easier to understand.
:::

## Working on a branch

Once on your branch, work normally — edit files, add, commit:

```bash
# Make changes to your code
git add script.js style.css
git commit -m "Add timer countdown to quiz questions"

# More work, more commits
git add script.js
git commit -m "Auto-advance when timer reaches zero"
```

All these commits happen on `feature/quiz-timer`. The `main` branch is untouched.

```text
main:           A ── B ── C
                          \
feature/quiz-timer:        D ── E
```

Commits D and E only exist on the feature branch. Anyone working on `main` doesn't see them.

## Switching back to main

```bash
git switch main
```

Now your working directory reflects the `main` branch — without the timer changes. Switch back to `feature/quiz-timer` and they reappear. The branches are completely independent.

:::warning
Always commit or stash your changes before switching branches. If you have uncommitted work, Git may prevent the switch or carry dirty changes to the other branch.
:::

## Merging a branch

When the feature is complete and tested, merge it into `main`:

```bash
git switch main
git merge feature/quiz-timer
```

This takes all the commits from `feature/quiz-timer` and applies them to `main`.

## Fast-forward vs merge commit

Git merges in two ways depending on the situation:

**Fast-forward merge** — when `main` hasn't changed since you branched:

```text
Before:
main:      A ── B ── C
                      \
feature:               D ── E

After fast-forward:
main:      A ── B ── C ── D ── E
```

Git simply moves `main` forward. No extra commit needed.

**Merge commit** — when `main` has new commits since you branched:

```text
Before:
main:      A ── B ── C ── F
                      \
feature:               D ── E

After merge:
main:      A ── B ── C ── F ── M (merge commit)
                      \       /
feature:               D ── E
```

Git creates a merge commit (M) that combines both lines of work.

:::key
Fast-forward merges produce a cleaner, linear history. Merge commits preserve the branching history. Both are normal. Don't stress about which type happens — Git chooses automatically based on whether `main` has moved.
:::

## Merge conflicts — the <<<< ==== >>>> markers

A **merge conflict** happens when the same lines were changed on both branches. Git can't decide which version to keep, so it asks you.

When you run `git merge` and there's a conflict, Git marks the file:

```js
function getMessage(percentage) {
<<<<<<< HEAD
  if (percentage >= 80) return "Great job!";
=======
  if (percentage >= 80) return "Excellent work!";
>>>>>>> feature/quiz-timer
}
```

The markers mean:

- `<<<<<<< HEAD` — this is what's on your current branch (`main`)
- `=======` — the divider between the two versions
- `>>>>>>> feature/quiz-timer` — this is what's on the incoming branch

## Resolving conflicts

1. **Open the file** and find the conflict markers
2. **Choose the correct version** — keep one side, combine both, or write something new
3. **Delete the markers** entirely (`<<<<<<<`, `=======`, `>>>>>>>`)
4. **Save the file**
5. **Stage and commit:**

```bash
git add script.js
git commit -m "Resolve merge conflict in getMessage"
```

```js
// After resolving — clean, no markers:
function getMessage(percentage) {
  if (percentage >= 80) return "Excellent work!";
}
```

:::warning
Never leave conflict markers in your code. The `<<<<<<<`, `=======`, and `>>>>>>>` lines will cause syntax errors. Always search the file for `<<<<` after resolving to make sure you caught everything.
:::

## Pushing branches to GitHub

To share a branch on GitHub (for a pull request):

```bash
git push -u origin feature/quiz-timer
```

The `-u` flag sets the upstream so future `git push` commands on this branch know where to go.

```bash
# After -u is set, just use:
git push
```

## Opening a pull request on GitHub

A **pull request** (PR) is a request to merge your branch into another branch on GitHub. It's the standard way teams review code before merging.

1. Push your branch to GitHub
2. Go to your repo on GitHub — you'll see a "Compare & pull request" banner
3. Click it (or go to Pull Requests → New)
4. Set the **base branch** to `main` and the **compare branch** to your feature branch
5. Write a title and description
6. Click **Create pull request**

## PR description best practices

A strong PR description answers three questions:

```markdown
## What
Added a countdown timer to quiz questions.

## Why
Users were spending too long on questions, reducing
engagement. A 15-second timer adds urgency.

## How
- Added `startTimer()` and `stopTimer()` functions
- Timer auto-advances to next question when it reaches zero
- Visual warning state when 5 seconds remain

## Screenshots
[screenshot of timer in action]

## Testing
- Verified timer resets on each question
- Tested auto-advance on timeout
- Checked that manual answer clears the timer
```

:::tip
Even on solo projects, writing PR descriptions builds a habit that will matter on your first team. Reviewers shouldn't have to read every line to understand what changed and why.
:::

## Reviewing a pull request

On the PR page, the **Files changed** tab shows a diff of every modified file. Green lines are additions; red lines are deletions.

As a reviewer (or reviewing your own work), check:

- Does the code do what the description says?
- Are there any bugs or edge cases?
- Is the code readable and well-organized?
- Are there unnecessary changes (leftover debug logs, unrelated formatting)?

## Merging and deleting branches

Once the PR is approved:

1. Click **Merge pull request** on GitHub
2. Click **Confirm merge**
3. Click **Delete branch** (GitHub offers this right after merging)

Locally, clean up:

```bash
git switch main
git pull                           # get the merged changes
git branch -d feature/quiz-timer   # delete the local branch
```

The `-d` flag safely deletes the branch (only if it's been merged). Use `-D` to force-delete an unmerged branch.

:::key
Delete branches after merging. Keeping dozens of old branches makes the repo messy and confusing. The work is preserved in the merge — the branch label is no longer needed.
:::

## The full branching workflow

```bash
# 1. Start the feature
git switch main
git pull
git switch -c feature/quiz-timer

# 2. Work and commit
git add .
git commit -m "Add timer UI and countdown logic"
git add .
git commit -m "Add auto-advance on timeout"

# 3. Push the branch
git push -u origin feature/quiz-timer

# 4. Open a PR on GitHub, get it reviewed, merge

# 5. Clean up
git switch main
git pull
git branch -d feature/quiz-timer
```

:::quiz
Q: What does `git switch -c feature/quiz-timer` do?
- Deletes the quiz-timer branch
- Creates a new branch called feature/quiz-timer and switches to it *
- Merges feature/quiz-timer into main
- Pushes the branch to GitHub
E: The `-c` flag tells `git switch` to create a new branch with the given name and immediately switch to it. Without `-c`, it would try to switch to an existing branch and fail if it doesn't exist.
:::

:::quiz
Q: What should you do when you see `<<<<<<<`, `=======`, and `>>>>>>>` markers in a file?
- Delete the file entirely
- Choose the correct code, remove all conflict markers, then stage and commit *
- Run `git merge` again
- Ignore them — the browser will skip those lines
E: Conflict markers are Git's way of showing you two competing versions of the same code. You must manually decide which version (or combination) to keep, delete all the markers, and then commit the resolved file.
:::

## Recap

- **`git switch -c`** creates and switches to a new branch.
- **Work on branches** to isolate features from the stable `main` branch.
- **`git merge`** combines a branch back into `main` — fast-forward or merge commit.
- **Merge conflicts** show competing changes with `<<<<<<<` / `=======` / `>>>>>>>` markers — resolve manually.
- **Push branches** with `git push -u origin branch-name`.
- **Pull requests** on GitHub are the standard way to review and merge code.
- **PR descriptions** explain what changed, why, and how.
- **Delete branches** after merging to keep the repo clean.

**Next up:** Git Best Practices — the habits that keep your project history clean and professional.

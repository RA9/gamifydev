# Backtracking

Greedy commits to one path and never looks back. Brute force tries every possibility to the bitter end. **Backtracking** sits between them: it builds a solution one piece at a time, and the instant a partial solution becomes impossible, it throws that whole branch away and steps back.

That "throw it away early" move is called **pruning**, and it's the entire point. Backtracking is what makes searches that look astronomically large — puzzles, schedules, mazes, board configurations — actually finish.

## Choose, explore, un-choose

Every backtracking algorithm has the same three-beat rhythm:

1. **Choose** — make the next decision, adding one piece to a partial solution.
2. **Explore** — recurse to make the following decision, given that choice.
3. **Un-choose** — undo the choice, so the next option can be tried from a clean state.

```python
def solve(partial):
    if is_complete(partial):
        return partial                # found a full solution
    for option in options(partial):
        if is_valid(partial, option):  # PRUNE: skip doomed branches
            partial.append(option)     # choose
            result = solve(partial)    # explore
            if result is not None:
                return result
            partial.pop()              # un-choose (backtrack)
    return None                        # every option failed here
```

That skeleton — with `is_valid` doing the pruning — is behind sudoku solvers, maze solvers, crossword fillers and every constraint puzzle you've ever seen.

:::key
Backtracking = choose, explore, un-choose. The `is_valid` check is what separates it from brute force: it kills a branch **before** the recursion, rather than discovering the failure at the very bottom.
:::

## The state-space tree

Picture every possible partial solution as a node in a tree. The root is the empty solution; each child adds one more decision; the leaves are complete solutions. That's the **state-space tree**, and backtracking is a depth-first traversal of it.

```text
                    ( )                      empty
              /      |      \
           (a)      (b)      (c)             first choice
          /   \    /   \    /   \
       (a,b) (a,c) X   ... ...  ...          second choice
        / \
      ...  X  <- invalid: prune here, never explore below
```

Every `X` is a branch that gets cut. And here's why that matters so much: pruning a node at depth 2 in a tree that's 8 levels deep doesn't save you one node — it saves you the **entire subtree** beneath it, which could be thousands of leaves. Cuts near the top of the tree are worth enormously more than cuts near the bottom.

:::analogy
Backtracking is solving a maze with a piece of chalk. Walk down a corridor, and the moment you hit a dead end, walk back to the last junction and mark that route as tried. You never re-walk a corridor you've already ruled out, and you never wander in circles.
:::

## N-Queens

The classic worked example. Place N chess queens on an N×N board so that no two attack each other — no shared row, column, or diagonal.

Brute force would be hopeless: on an 8×8 board there are over four billion ways to place 8 pieces on 64 squares. Backtracking makes it instant, using two observations. First, each queen must be in a different row, so we can place exactly one queen per row and only choose its *column*. Second, we can check each placement against the queens already placed, and reject immediately.

```python
def solve_n_queens(n):
    cols = []                            # cols[r] = column of the queen in row r

    def is_safe(row, col):
        for r, c in enumerate(cols):
            if c == col:                          # same column
                return False
            if abs(row - r) == abs(col - c):      # same diagonal
                return False
        return True

    def place(row):
        if row == n:                     # all n rows filled: solved
            return list(cols)
        for col in range(n):
            if is_safe(row, col):        # PRUNE
                cols.append(col)         # choose
                result = place(row + 1)  # explore
                if result is not None:
                    return result
                cols.pop()               # un-choose
        return None                      # no column works in this row

    return place(0)

print(solve_n_queens(4))   # [1, 3, 0, 2]
print(solve_n_queens(8))   # [0, 4, 7, 5, 2, 6, 1, 3]
```

The diagonal test is the neat bit: two squares are on the same diagonal exactly when the difference in their rows equals the difference in their columns.

Here's the 4-queens solution `[1, 3, 0, 2]` — row 0 has its queen in column 1, row 1 in column 3, and so on:

```text
    . Q . .        row 0, col 1
    . . . Q        row 1, col 3
    Q . . .        row 2, col 0
    . . Q .        row 3, col 2
```

And here's the search finding it, with each pruned attempt marked:

```text
row 0: try col 0  -> place
  row 1: col 0 X (same col)   col 1 X (diagonal)   col 2 -> place
    row 2: col 0 X  col 1 X  col 2 X  col 3 X   -> dead end, back up
  row 1: col 3 -> place
    row 2: col 1 -> place
      row 3: every column attacked -> dead end, back up, back up, back up
row 0: try col 1  -> place
  row 1: col 3 -> place
    row 2: col 0 -> place
      row 3: col 2 -> place.  row == 4.  SOLVED
```

Notice how quickly whole branches die. That's pruning doing its job.

:::tip
To find *all* solutions rather than the first, don't return early — collect the completed solution into a list and keep going. The un-choose step is what makes that work: after recording a solution you `pop()` and continue exploring as normal.
:::

## Other classic backtracking problems

**Maze solving.** From your current cell, try each of the four directions. If a direction is a wall, or off the grid, or already on your current path, prune it. Otherwise step there and recurse. If every direction fails, un-step and report failure to the caller. The path list *is* the partial solution.

```python
def solve_maze(grid, pos, goal, path):
    r, c = pos
    if not (0 <= r < len(grid) and 0 <= c < len(grid[0])):
        return None                       # off the grid
    if grid[r][c] == "#" or pos in path:  # wall, or already on this path
        return None
    path.append(pos)                      # choose
    if pos == goal:
        return list(path)
    for step in ((1, 0), (-1, 0), (0, 1), (0, -1)):
        found = solve_maze(grid, (r + step[0], c + step[1]), goal, path)
        if found is not None:             # explore
            return found
    path.pop()                            # un-choose
    return None

maze = ["...#", ".#..", "...."]
print(solve_maze(maze, (0, 0), (2, 3), []))
# [(0, 0), (1, 0), (2, 0), (2, 1), (2, 2), (1, 2), (1, 3), (2, 3)]
```

Note that the route it returns wanders — it steps up to row 1 near the end rather than going straight along the bottom. Backtracking returns the **first** path it finds, not the shortest one. If you need the shortest, that's BFS's job, not backtracking's.

**The knight's tour.** Move a knight around a chessboard so it visits every square exactly once. From each square there are up to eight legal moves; you try one, recurse, and back up if the tour gets stuck. Plain backtracking solves small boards easily but struggles on a full 8×8, which is why people add a heuristic: always try the move with the *fewest* onward options first. That gets you into the dead-end branches early, when backing out is cheap.

:::warning
Backtracking is not magic. Pruning shrinks the search space, sometimes dramatically, but the worst case can still be exponential. If your `is_valid` check is weak, backtracking degenerates into brute force with extra function calls.
:::

## How it differs from brute force and greedy

All three explore the same space of possible answers. What differs is how much they look at.

**Brute force** generates every complete candidate and tests each one at the end. For N-Queens, it would place all 8 queens and only then ask whether they attack each other — discovering failure after eight decisions instead of after two.

**Backtracking** tests *partial* candidates and cuts as soon as a partial one is doomed. Same search space, but whole subtrees are never generated.

**Greedy** makes one decision at each step and never revisits it. It explores a single root-to-leaf path — extremely fast, but with no recovery if that path is bad.

```text
                     explores            can recover?    typical cost
brute force    every complete candidate      n/a          exponential
backtracking   partial candidates, pruned    yes          exponential worst,
                                                          often far better
greedy         one path, no reconsidering    no           usually O(n log n)
```

So the ordering is: greedy if you can prove it correct, backtracking when you can't but the constraints prune well, brute force only as a baseline to check the others against.

:::example
For 8-queens, brute force examines over four billion placements. Backtracking finds the first valid solution after exploring only a few thousand partial boards — because the "no two queens attack" rule kills most branches within the first two or three rows.
:::

## Check Your Understanding

:::quiz
Q: What makes backtracking different from plain brute force?
- It uses less memory
- It abandons a partial solution as soon as it cannot possibly work *
- It always finds the optimal answer
- It never uses recursion
E: Both search the same space, but backtracking prunes doomed partial solutions before recursing, so entire subtrees are never generated at all.
:::

:::fill
Q: Complete the missing "un-choose" step that lets backtracking try the next option.
`partial.___()`
- pop *
- append
- clear
- sort
E: After exploring a choice fails, you remove it so the state is clean for the next option. Forgetting this leaves stale choices in the partial solution and corrupts every branch that follows.
:::

:::quiz
Q: In a state-space tree, where is pruning most valuable?
- At the leaves, where the answers are
- Near the root, because it removes an entire subtree *
- It makes no difference where you prune
- Only on the rightmost branch
E: Cutting a node removes everything beneath it. Near the root a single cut can eliminate a huge fraction of the search space; at a leaf it saves one candidate.
:::

## Talk about it

> Backtracking's power comes entirely from detecting failure *early*. Think of a task where you can tell a plan won't work long before you finish carrying it out. How much effort does noticing early save you, and what would it cost to only find out at the very end?

## What's next

You've now got a search technique that prunes intelligently. Next up is **String Search Algorithms**, where the problem is finding a small pattern inside a large text. The naive approach re-compares characters it has already matched, and three clever algorithms each find a different way to stop wasting that work — including one that scans the pattern backwards and skips whole chunks of text at a time.

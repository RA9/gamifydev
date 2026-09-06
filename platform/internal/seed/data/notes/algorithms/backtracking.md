# Backtracking

Greedy commits to one path and never looks back. Brute force tries every possibility to the bitter end. **Backtracking** sits between them: it builds a solution one piece at a time, and the instant a partial solution becomes impossible, it throws that whole branch away and steps back.

That "throw it away early" move is called **pruning**, and it's the entire point. Backtracking is what makes searches that look astronomically large — puzzles, schedules, mazes, board configurations — actually finish.

## Choose, explore, un-choose

Every backtracking algorithm has the same three-beat rhythm:

1. **Choose** — make the next decision, adding one piece to a partial solution.
2. **Explore** — recurse to make the following decision, given that choice.
3. **Un-choose** — undo the choice, so the next option can be tried from a clean state.

```c
#include <stdbool.h>
#include <stddef.h>

#define MAX_CHOICES 64

bool is_complete(const int partial[], size_t length);
size_t options(const int partial[], size_t length, int choices[], size_t capacity);
bool is_valid(const int partial[], size_t length, int option);

bool solve(int partial[], size_t *length, size_t capacity) {
    if (is_complete(partial, *length)) return true;

    int choices[MAX_CHOICES];
    size_t choice_count = options(partial, *length, choices, MAX_CHOICES);
    for (size_t i = 0; i < choice_count; ++i) {
        int option = choices[i];
        if (*length < capacity && is_valid(partial, *length, option)) {
            partial[(*length)++] = option;       /* Choose. */
            if (solve(partial, length, capacity)) return true; /* Explore. */
            --*length;                           /* Un-choose. */
        }
    }
    return false;
}
```

That skeleton — with `is_valid` doing the pruning — is behind sudoku solvers, maze solvers, crossword fillers and every constraint puzzle you've ever seen.

:::key
Backtracking = choose, explore, un-choose. The `is_valid` check is what separates it from brute force: it kills a branch **before** the recursion, rather than discovering the failure at the very bottom. In C, the partial solution's current length also makes the un-choose operation explicit.
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

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

static bool is_safe(const size_t cols[], size_t row, size_t col) {
    for (size_t r = 0; r < row; ++r) {
        size_t row_gap = row - r;
        size_t col_gap = cols[r] > col ? cols[r] - col : col - cols[r];
        if (cols[r] == col || row_gap == col_gap) return false;
    }
    return true;
}

static bool place_queen(size_t cols[], size_t n, size_t row) {
    if (row == n) return true;
    for (size_t col = 0; col < n; ++col) {
        if (is_safe(cols, row, col)) {
            cols[row] = col;                     /* Choose. */
            if (place_queen(cols, n, row + 1)) return true; /* Explore. */
            /* Un-choose: the next column overwrites cols[row]. */
        }
    }
    return false;
}

size_t *solve_n_queens(size_t n) {
    if (n == 0 || n > SIZE_MAX / sizeof(size_t)) return NULL;
    size_t *cols = malloc(n * sizeof *cols);
    if (cols == NULL) return NULL;
    if (!place_queen(cols, n, 0)) { free(cols); return NULL; }
    return cols;                                  /* Caller must free it. */
}

int main(void) {
    size_t *cols = solve_n_queens(4);
    if (cols == NULL) return EXIT_FAILURE;
    for (size_t row = 0; row < 4; ++row) printf("%zu%s", cols[row], row == 3 ? "\n" : " ");
    free(cols);                                    /* 1 3 0 2 */
    return 0;
}
```

The diagonal test is the neat bit: two squares are on the same diagonal exactly when the difference in their rows equals the difference in their columns.

Here's the 4-queens solution `{1, 3, 0, 2}` — row 0 has its queen in column 1, row 1 in column 3, and so on:

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
To find *all* solutions rather than the first, don't return early — copy each completed `cols` array into caller-owned storage and keep going. The un-choose step is what makes that work: after recording a solution, return to the previous row and continue trying columns as normal.
:::

## Other classic backtracking problems

**Maze solving.** From your current cell, try each of the four directions. If a direction is a wall, or off the grid, or already on your current path, prune it. Otherwise step there and recurse. If every direction fails, un-step and report failure to the caller. The path list *is* the partial solution.

```c
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

typedef struct { int row, col; } Cell;

static bool same_cell(Cell a, Cell b) { return a.row == b.row && a.col == b.col; }

bool solve_maze(const char *grid[], int rows, int cols, Cell pos, Cell goal,
                Cell path[], size_t *length, size_t capacity) {
    if (pos.row < 0 || pos.row >= rows || pos.col < 0 || pos.col >= cols) return false;
    if (grid[pos.row][pos.col] == '#') return false;
    for (size_t i = 0; i < *length; ++i) if (same_cell(path[i], pos)) return false;
    if (*length == capacity) return false;

    path[(*length)++] = pos;                       /* Choose. */
    if (same_cell(pos, goal)) return true;
    static const Cell steps[] = {{1, 0}, {-1, 0}, {0, 1}, {0, -1}};
    for (size_t i = 0; i < 4; ++i) {
        Cell next = {pos.row + steps[i].row, pos.col + steps[i].col};
        if (solve_maze(grid, rows, cols, next, goal, path, length, capacity)) return true;
    }
    --*length;                                      /* Un-choose. */
    return false;
}

int main(void) {
    const char *maze[] = {"...#", ".#..", "...."};
    Cell path[12]; size_t length = 0;
    if (!solve_maze(maze, 3, 4, (Cell){0, 0}, (Cell){2, 3}, path, &length, 12)) return 1;
    for (size_t i = 0; i < length; ++i) printf("(%d,%d)%s", path[i].row, path[i].col, i + 1 == length ? "\n" : " ");
    return 0;
}
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
Q: Complete the missing "un-choose" step when `length` tracks the used portion of a C array.
`___length;`
- -- *
- ++
- clear
- sort
E: After exploring a choice fails, decrement the used length so the next option overwrites that slot. Forgetting this leaves stale choices in the partial solution and corrupts every branch that follows.
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

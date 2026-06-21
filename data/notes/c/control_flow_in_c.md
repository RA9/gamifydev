# Control Flow in C

So far your programs have run straight from top to bottom, every line, every time. Real programs need to *choose* and *repeat* — skip a block when a condition isn't met, or run the same logic a hundred times. That power is called **control flow**, and it's where your code starts to feel alive. Let's give your programs a brain.

## Making decisions with if

The `if` statement runs a block of code only when a condition is true. A condition is any expression that evaluates to true or false.

```c
#include <stdio.h>

int main(void) {
    int score = 85;
    if (score >= 60) {
        printf("You passed!\n");
    }
    return 0;
}
```

If you want a fallback when the condition is false, add `else`. And to test several possibilities in order, chain `else if`:

```c
#include <stdio.h>

int main(void) {
    int score = 85;
    if (score >= 90) {
        printf("Grade: A\n");
    } else if (score >= 80) {
        printf("Grade: B\n");
    } else if (score >= 70) {
        printf("Grade: C\n");
    } else {
        printf("Keep practicing!\n");
    }
    return 0;
}
```

C checks each condition top to bottom and runs the **first** one that's true, then skips the rest.

:::analogy
An `if/else if/else` chain is like a series of forks on a road. You take the first turn whose sign matches your destination, and once you've turned, you ignore all the later forks. The final `else` is the "if nothing matched, go here" road.
:::

## Comparison and logical operators

To build conditions, you need operators that produce true/false results.

Comparison operators: `==` (equal), `!=` (not equal), `<`, `>`, `<=`, `>=`.

Logical operators let you combine conditions: `&&` (and), `||` (or), `!` (not).

```c
int age = 20;
int hasTicket = 1;

if (age >= 18 && hasTicket) {
    printf("Welcome in!\n");
}
```

:::warning
The single most infamous C bug: writing `=` when you mean `==`. `if (x = 5)` *assigns* 5 to `x` and is always true — it does not compare. To test equality, you need the double `==`. Read your conditions carefully.
:::

In C, there's no separate `true`/`false` keyword by default — `0` means false and **any non-zero value** means true. That's why `if (hasTicket)` works when `hasTicket` is `1`.

:::quiz
Q: What does `if (x = 3)` do?
- Checks whether x equals 3
- Assigns 3 to x, and the condition is always true *
- Causes a compile error
E: A single `=` is assignment. The expression's value becomes 3 (non-zero = true), so the block always runs. Use `==` to compare.
:::

## Choosing among many values with switch

When you're checking one variable against several fixed values, a `switch` is cleaner than a long `if/else if` chain.

```c
#include <stdio.h>

int main(void) {
    int day = 3;
    switch (day) {
        case 1:
            printf("Monday\n");
            break;
        case 2:
            printf("Tuesday\n");
            break;
        case 3:
            printf("Wednesday\n");
            break;
        default:
            printf("Another day\n");
    }
    return 0;
}
```

Each `case` is a possible value. The `break` stops the switch once a match runs.

:::warning
Forgetting `break` causes **fall-through**: execution keeps running into the next case below it. Sometimes that's intentional, but usually a missing `break` is an accident that prints more than you expected. Add a `break` to every case unless you truly want to fall through.
:::

## Repeating with loops

Loops let you run a block of code over and over. C gives you three flavors.

The **`while`** loop runs as long as its condition stays true:

```c
int count = 1;
while (count <= 3) {
    printf("%d\n", count);
    count++;        // count++ adds 1 — without it, this loops forever!
}
```

The **`for`** loop packs setup, condition, and update into one tidy line — ideal when you know how many times to repeat:

```c
for (int i = 1; i <= 3; i++) {
    printf("%d\n", i);
}
```

Read the `for` header as: start `i` at 1; keep going while `i <= 3`; after each pass, do `i++`.

The **`do...while`** loop runs its body *first*, then checks — so it always runs at least once:

```c
int n = 10;
do {
    printf("This prints once even though the condition is false.\n");
} while (n < 5);
```

:::key
Use `for` when you know the number of repetitions, `while` when you loop until some condition changes, and `do...while` when the body must run at least once before checking.
:::

:::predict
```c
#include <stdio.h>

int main(void) {
    int sum = 0;
    for (int i = 1; i <= 4; i++) {
        sum = sum + i;
    }
    printf("%d\n", sum);
    return 0;
}
```
- 10 *
- 4
- 6
E: The loop adds 1 + 2 + 3 + 4 = 10. After `i` reaches 4 it increments to 5, the condition fails, and the loop ends.
:::

## Steering loops with break and continue

Inside a loop you have two special controls. `break` exits the loop immediately. `continue` skips the rest of the current pass and jumps to the next one.

```c
#include <stdio.h>

int main(void) {
    for (int i = 1; i <= 10; i++) {
        if (i == 5) {
            break;          // stop the whole loop at 5
        }
        if (i % 2 == 0) {
            continue;       // skip even numbers
        }
        printf("%d\n", i);  // prints 1, then 3
    }
    return 0;
}
```

:::warning
Watch out for **infinite loops** — a loop whose condition never becomes false. The classic cause is forgetting to change the variable the condition depends on (like missing `count++`). If your program hangs, suspect a loop that never updates its counter.
:::

:::fill
Q: Complete the keyword that immediately exits a loop.
`if (found) { ___; }`
- break *
- continue
- return
E: `break` exits the loop entirely. `continue` only skips to the next pass.
:::

## Talk about it

> You met three loops — `while`, `for`, and `do...while`. Think of a real task (counting steps, asking a user until they type a valid answer, printing a calendar). Which loop fits which task best, and why? Explaining your choice out loud is how the difference sticks.

## What's next

Your programs can now decide and repeat — that's a huge leap. But as logic grows, repeating yourself gets messy. In **Functions in C**, you'll learn to package reusable chunks of logic with names, so your code stays clean and your ideas stay organized. Let's go.

# Functions in C

As programs grow, cramming everything into `main` becomes a tangled mess. The cure is **functions** — named, reusable blocks of logic you can call whenever you need them. Functions are how you break a big problem into small, understandable pieces. Master them and your code becomes cleaner, shorter, and far easier to reason about.

## What a function is

You've already been using one function the whole time: `main`. A function is a self-contained block of code that does one job. You **define** it once, then **call** it as many times as you like.

```c
#include <stdio.h>

void greet(void) {
    printf("Hello from a function!\n");
}

int main(void) {
    greet();      // call it
    greet();      // call it again — no need to rewrite the logic
    return 0;
}
```

Here `greet` is defined above `main`, and `main` calls it twice. Write the logic once, reuse it forever.

:::analogy
A function is like a recipe card. Writing the recipe is *defining* the function. Each time you cook the dish, you *call* the recipe — you don't rewrite the instructions every time. Some recipes need ingredients (parameters) and some produce a dish you take away (a return value).
:::

## Parameters and return types

Functions become powerful when they accept **inputs** (parameters) and produce an **output** (a return value).

The word before the function name is its **return type** — the kind of value it hands back. The values in parentheses are its **parameters** — the inputs it needs.

```c
#include <stdio.h>

int add(int a, int b) {
    return a + b;
}

int main(void) {
    int sum = add(3, 4);   // sum becomes 7
    printf("%d\n", sum);
    return 0;
}
```

Reading the header `int add(int a, int b)`: this function is named `add`, takes two `int` parameters, and returns an `int`. The `return` keyword sends a value back to whoever called the function and ends it immediately.

When a function returns nothing, its return type is **`void`** — that's why `greet` above was `void greet(void)`.

:::quiz
Q: In `double area(double w, double h)`, what does the first `double` mean?
- The function takes a double
- The function returns a double *
- The function is named double
E: The type written *before* the function name is the return type — the kind of value the function gives back.
:::

:::predict
```c
#include <stdio.h>

int square(int n) {
    return n * n;
}

int main(void) {
    printf("%d\n", square(5));
    return 0;
}
```
- 25 *
- 10
- 5
E: `square(5)` returns `5 * 5`, which is 25.
:::

## Function prototypes

C reads your file top to bottom. If `main` calls a function defined *below* it, the compiler hasn't seen that function yet and complains. The fix is a **prototype** (also called a declaration): a one-line preview of the function near the top, ending in a semicolon.

```c
#include <stdio.h>

int multiply(int a, int b);   // prototype — "this function exists, here's its shape"

int main(void) {
    printf("%d\n", multiply(6, 7));
    return 0;
}

int multiply(int a, int b) {  // full definition, defined later
    return a * b;
}
```

The prototype tells the compiler the function's name, parameter types, and return type — enough to verify your calls are correct before it ever reaches the full definition.

:::tip
A prototype is just the function header followed by a semicolon. You can copy the first line of the definition, add `;`, and place it near the top. In real projects, prototypes live in header (`.h`) files so many source files can share them.
:::

:::reorder
Arrange this program so it compiles cleanly with the prototype on top.
- #include <stdio.h>
- int cube(int n);
- int main(void) {
-   printf("%d", cube(3));
-   return 0;
- }
- int cube(int n) { return n * n * n; }
E: Include the header, declare the prototype, define `main` which calls `cube`, then provide the full `cube` definition below.
:::

## Scope: local vs. global

**Scope** is the region of code where a variable exists and can be used. A variable declared *inside* a function is **local** — it lives only while that function runs and is invisible to everyone else.

```c
#include <stdio.h>

void demo(void) {
    int secret = 42;     // local to demo
    printf("%d\n", secret);
}

int main(void) {
    demo();
    // printf("%d", secret);  // ERROR: secret doesn't exist out here
    return 0;
}
```

A variable declared *outside* all functions is **global** — every function can see it. Globals seem convenient but are usually a trap: when anything can change a value, bugs become very hard to track down.

:::key
Prefer **local** variables and pass data through parameters and return values. Local scope keeps each function self-contained and predictable. Reach for globals only when truly necessary.
:::

:::warning
Two functions can each have a local variable named `count`, and they won't interfere — they're separate boxes in separate scopes. Don't assume same-named variables in different functions are connected. They aren't.
:::

## Why decomposition matters

The real win of functions isn't saving keystrokes — it's **decomposition**, breaking a hard problem into small pieces you can understand one at a time.

Compare a 200-line `main` that does everything against a `main` that reads like a summary:

```c
int main(void) {
    int n = readNumber();
    int result = computeFactorial(n);
    printResult(result);
    return 0;
}
```

Even without seeing the bodies, you understand the whole program at a glance. Each function can be written, tested, and fixed on its own. This is how professionals tame complexity: many small, well-named functions, each doing one job well.

:::fill
Q: Complete the keyword that sends a value back from a function.
`int twice(int n) { ___ n * 2; }`
- return *
- void
- break
E: `return` hands a value back to the caller and ends the function.
:::

## Talk about it

> Imagine writing a program that plays a game of tic-tac-toe. Without writing code, what functions would you create to break the problem into pieces? Naming them (like `printBoard`, `checkWinner`, `getMove`) is the first real step of thinking like a programmer.

## What's next

You can now organize logic into clean, reusable functions — a skill that scales to programs of any size. Next comes the lesson C is famous for: **Pointers and Memory**. It's the part that intimidates newcomers, but with the right mental model it clicks. Pixel believes in you — let's tackle it together.

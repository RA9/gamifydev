# Control Flow in Java

So far your programs have marched straight down the page, one line after another, never pausing to think. Real programs aren't like that. They make choices ("if the player has no lives left, end the game") and they repeat work ("for each enemy on screen, move it"). This is **control flow**, the art of directing which code runs, when, and how many times. Master this, and you can build almost anything.

## Making Decisions with `if`

The simplest decision is `if`: run some code *only when* a condition is true.

```java
int score = 85;

if (score >= 80) {
    System.out.println("Great job!");
}
```

The part in parentheses, `score >= 80`, is a **condition**. It's a question that evaluates to either `true` or `false`. If it's true, the code inside the braces runs. If false, Java skips it entirely.

You can handle the "otherwise" case with `else`, and chain several possibilities with `else if`:

```java
int score = 72;

if (score >= 90) {
    System.out.println("A");
} else if (score >= 80) {
    System.out.println("B");
} else if (score >= 70) {
    System.out.println("C");
} else {
    System.out.println("Keep practicing!");
}
```

Java checks each condition from top to bottom and runs the **first** one that's true, then skips the rest. Since `72` isn't ≥ 90 or ≥ 80 but *is* ≥ 70, this prints `C`.

:::analogy
An `if/else if/else` chain is like a series of forks in a road. You walk up to the first fork: if the condition is true, you take that path and the journey ends. If not, you walk to the next fork and ask again. The final `else` is the path you take when none of the others opened.
:::

## Comparison and Logical Operators

To build conditions, you need operators that produce `true` or `false`.

**Comparison operators** compare two values:

```java
a == b   // equal to        (note: TWO equals signs)
a != b   // not equal to
a >  b   // greater than
a <  b   // less than
a >= b   // greater than or equal to
a <= b   // less than or equal to
```

:::warning
A single `=` *assigns* a value; a double `==` *compares* two values. Writing `if (x = 5)` when you meant `if (x == 5)` is a classic mistake. In Java the compiler will usually catch it, but train your eyes now: comparing always uses `==`.
:::

**Logical operators** combine conditions:

```java
boolean canEnter = (age >= 18) && (hasTicket);   // AND: both must be true
boolean getsDiscount = (isStudent) || (isSenior); // OR: at least one true
boolean isClosed = !isOpen;                       // NOT: flips true/false
```

- `&&` (AND) is true only when **both** sides are true.
- `||` (OR) is true when **at least one** side is true.
- `!` (NOT) flips a boolean: `!true` becomes `false`.

```java
int age = 20;
boolean hasTicket = true;

if (age >= 18 && hasTicket) {
    System.out.println("Welcome to the show!");
}
```

## Choosing with `switch`

When you're checking one variable against many possible values, a long `if/else if` chain gets tiring. The `switch` statement is a cleaner alternative for this exact case:

```java
int day = 3;

switch (day) {
    case 1:
        System.out.println("Monday");
        break;
    case 2:
        System.out.println("Tuesday");
        break;
    case 3:
        System.out.println("Wednesday");
        break;
    default:
        System.out.println("Another day");
}
```

Java jumps straight to the matching `case`, runs its code, and `break` tells it to stop. The `default` case runs if nothing else matched, like the final `else`.

:::warning
Don't forget the `break` at the end of each case. Without it, Java "falls through" and keeps running the next case's code too, a surprising bug. Until you specifically want fall-through behavior, always include `break`.
:::

## Repeating with Loops

Loops let you run the same block of code many times without copying it. Java gives you a few flavors.

The **`while` loop** repeats *as long as* a condition stays true:

```java
int count = 1;
while (count <= 3) {
    System.out.println("Count is " + count);
    count++;   // count++ means "add 1 to count"
}
```

This prints `Count is 1`, `2`, then `3`. Notice `count++` increasing the counter each pass, that's crucial.

:::warning
If your `while` condition never becomes false, you get an **infinite loop** and your program hangs forever. The fix is almost always making sure something inside the loop moves the condition toward `false`, like `count++` here.
:::

The **`for` loop** packs the setup, condition, and update into one tidy line. It's perfect when you know how many times to repeat:

```java
for (int i = 0; i < 5; i++) {
    System.out.println("Iteration " + i);
}
```

Read the header as three parts separated by semicolons: **start** (`int i = 0`), **keep going while** (`i < 5`), and **after each pass** (`i++`). This runs five times, printing iterations 0 through 4.

:::analogy
A `for` loop is like a recipe step that says "stir 5 times." You set a counter, check whether you've reached the target, do the work, and tick the counter up, all in one breath.
:::

## Walking Collections and Steering Loops

When you have a list of items and just want to visit each one, the **enhanced for-each loop** is the clearest tool:

```java
int[] scores = {90, 75, 88};

for (int score : scores) {
    System.out.println("Score: " + score);
}
```

Read `for (int score : scores)` as "for each score in scores." You don't manage a counter at all, Java hands you each item in turn.

Two keywords let you steer a loop from inside it:

- **`break`** exits the loop immediately.
- **`continue`** skips the rest of the current pass and jumps to the next one.

```java
for (int i = 1; i <= 10; i++) {
    if (i == 5) {
        break;       // stop the loop entirely when i hits 5
    }
    if (i % 2 == 0) {
        continue;    // skip even numbers
    }
    System.out.println(i);   // prints 1, then 3
}
```

:::key
Use `if/else if/else` and `switch` to make decisions, and `while`, `for`, and for-each to repeat work. Build conditions with comparison (`==`, `>`, `<`) and logical (`&&`, `||`, `!`) operators, and steer loops with `break` and `continue`.
:::

## Check Your Understanding

:::predict
```java
int x = 4;
if (x > 5) {
    System.out.println("big");
} else if (x > 2) {
    System.out.println("medium");
} else {
    System.out.println("small");
}
```
- medium *
- big
- small
E: `x` is 4, which is not greater than 5 but *is* greater than 2, so the first true branch (`medium`) runs and the rest are skipped.
:::

:::quiz
Q: What does the `&&` operator require to evaluate to `true`?
- At least one side is true
- Both sides are true *
- Neither side is true
- The sides are different
E: `&&` (AND) is true only when *both* conditions are true. For "at least one," you'd use `||` (OR).
:::

:::predict
```java
for (int i = 0; i < 4; i++) {
    if (i == 2) continue;
    System.out.print(i);
}
```
- 013 *
- 0123
- 012
E: The loop runs for i = 0, 1, 2, 3. When i is 2, `continue` skips the print, so you get 0, 1, then 3 with no spaces: 013.
:::

## Talk about it

> Loops and conditionals are how programs handle situations the programmer can't predict in advance, like a user who might enter 3 items or 300. Think of a small everyday task (sorting laundry, checking a to-do list). How would you describe it using "if" decisions and "repeat" loops? Talking through this trains your brain to think like a programmer.

## What's next

Your programs can now think and repeat, an enormous leap. But as they grow, you'll notice the same chunks of logic appearing again and again. Next, in **Methods in Java**, you'll learn to wrap reusable behavior into named methods you can call whenever you need them, keeping your code clean, organized, and far easier to understand. Onward, with Pixel by your side.

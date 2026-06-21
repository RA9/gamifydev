# C Basics

Every program is really just data being moved, transformed, and displayed. In this lesson you'll learn how C holds onto data using **variables**, the **types** that describe what kind of data you're storing, and how to show results on screen with `printf`. These are the building blocks for everything else. Take your time — these fundamentals pay off forever.

## Variables: named boxes for data

A **variable** is a named place in memory where your program stores a value. Before you can use one in C, you must *declare* it — tell the compiler its name and what type of value it will hold.

```c
#include <stdio.h>

int main(void) {
    int age = 25;
    printf("Age: %d\n", age);
    return 0;
}
```

That line `int age = 25;` does three things: it says the type is `int` (a whole number), names the variable `age`, and stores the value `25` in it.

:::analogy
A variable is like a labeled box on a shelf. The label is the name (`age`), and the box can only hold a certain kind of thing — the *type*. You can open the box and change what's inside, but the label and the kind of thing it holds stay fixed.
:::

You can declare first and assign later:

```c
int score;        // declare: an empty box for a whole number
score = 100;      // assign: put 100 into it
score = 150;      // reassign: now it holds 150
```

:::warning
C does **not** automatically clear a variable when you declare it. If you read a variable before assigning it a value, you get whatever random "garbage" was already sitting in that memory. Always give your variables a value before using them.
:::

## C's core types

C is **statically typed**, meaning every variable's type is fixed when you declare it. Here are the types you'll use constantly:

- **`int`** — whole numbers like `-3`, `0`, `42`.
- **`char`** — a single character like `'A'` or `'?'` (note the single quotes).
- **`float`** — a decimal number, single precision.
- **`double`** — a decimal number, double precision (more accurate; prefer this for decimals).

```c
#include <stdio.h>

int main(void) {
    int   apples   = 4;
    char  grade    = 'A';
    double price   = 2.50;
    printf("%d apples, grade %c, $%.2f each\n", apples, grade, price);
    return 0;
}
```

:::match
Q: Match each type to what it stores.
- `int` | Whole numbers
- `char` | A single character
- `double` | Decimal numbers
E: Each C type holds a specific kind of value. Choosing the right one matters for both correctness and memory.
:::

## Printing with format specifiers

`printf` doesn't just print plain text — it prints *formatted* text. Inside the quotes you place **format specifiers**, little placeholders that start with `%`. C then fills each one with a value you pass after the string.

The common specifiers:

- `%d` — an `int` (decimal integer)
- `%c` — a single `char`
- `%f` — a `float` or `double`
- `%s` — a string of text
- `%.2f` — a decimal rounded to 2 places

```c
#include <stdio.h>

int main(void) {
    int   level = 7;
    double hp   = 88.5;
    printf("Level %d with %.1f HP\n", level, hp);
    return 0;
}
```

The values after the string are matched to the specifiers *in order*: `level` fills the `%d`, `hp` fills the `%.1f`.

:::fill
Q: Complete the format specifier to print an integer.
`printf("Items: %___\n", count);`
- d *
- s
- c
E: `%d` prints a decimal integer. `%s` is for strings and `%c` is for a single character.
:::

## Reading input with scanf

The mirror image of `printf` is `scanf`, which reads typed input from the user. It uses the same format specifiers, but with one crucial difference: you must put an `&` (ampersand) before the variable, so `scanf` knows *where* to store the value.

```c
#include <stdio.h>

int main(void) {
    int age;
    printf("Enter your age: ");
    scanf("%d", &age);
    printf("Next year you'll be %d\n", age + 1);
    return 0;
}
```

:::warning
Forgetting the `&` in `scanf` is one of the most common beginner bugs. `scanf("%d", age)` (no `&`) compiles with a warning but misbehaves at run time. You'll understand exactly *why* the `&` is needed once you reach the Pointers lesson.
:::

## Constants, comments, and arithmetic

Sometimes a value should never change — like the number of days in a week. Mark it with `const`:

```c
const int DAYS_IN_WEEK = 7;
```

Now the compiler will stop you if you accidentally try to change it. Use `const` to make your intentions clear and your code safer.

**Comments** are notes for humans; the compiler ignores them:

```c
// This is a single-line comment
/* This is a
   multi-line comment */
```

C handles **arithmetic** with the operators you'd expect: `+`, `-`, `*`, `/`, and `%` (remainder). One sharp edge to know:

```c
int a = 7 / 2;       // 3, NOT 3.5 — integer division drops the decimal
int r = 7 % 2;       // 1, the remainder
double b = 7.0 / 2;  // 3.5 — a decimal in the mix gives a decimal result
```

:::key
Dividing two integers in C throws away the fractional part. `7 / 2` is `3`, not `3.5`. To keep decimals, make at least one value a `double` (e.g. `7.0`).
:::

:::predict
```c
#include <stdio.h>

int main(void) {
    int total = 10;
    int people = 4;
    printf("%d\n", total / people);
    return 0;
}
```
- 2 *
- 2.5
- 3
E: Both values are `int`, so integer division applies: `10 / 4` is `2` with the remainder discarded.
:::

:::quiz
Q: What does the `%` operator do in `13 % 5`?
- Divides and keeps the decimal
- Gives the remainder after division (here, 3) *
- Calculates a percentage
E: `%` is the modulo operator. `13 / 5` is 2 with 3 left over, so `13 % 5` is `3`.
:::

## Talk about it

> You learned that `int` division silently drops the decimal part. Why do you think C behaves this way instead of automatically giving you a decimal? When might "throw away the remainder" actually be the behavior you *want*?

## What's next

You can now store data, do math with it, and move it in and out of your program. Next up is **Control Flow in C**, where your programs start making decisions and repeating work — the moment code stops being a straight line and starts being *smart*. Onward!

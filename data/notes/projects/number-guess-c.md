# Number Guessing Game

A classic first game in C: the computer picks a secret number and you hunt for it with "too high" and "too low" hints. You'll touch input, randomness, conditions, and loops — the core of nearly every program.

:::project
You'll have a complete C program that picks a random number from 1 to 100 and lets you guess repeatedly, giving hints each time, until you find it and get congratulated.
:::

To follow along you need a C compiler such as `gcc`. Save your code as `guess.c`, then compile and run it with `gcc guess.c -o guess` followed by `./guess`.

## Step 1 — Start with a skeleton

Every C program begins with `main`. Include `<stdio.h>` for `printf`, print a welcome line, and `return 0` to signal success.

```c
#include <stdio.h>

int main(void) {
    printf("Welcome to the Number Guessing Game!\n");
    return 0;
}
```

## Step 2 — Pick a random secret number

Add `<stdlib.h>` for `rand()` and `<time.h>` for the clock. Seeding with `srand(time(NULL))` makes the number different each run. The `% 100 + 1` squeezes the result into the range 1–100.

```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main(void) {
    srand(time(NULL));
    int secret = rand() % 100 + 1;

    printf("I'm thinking of a number between 1 and 100.\n");
    return 0;
}
```

## Step 3 — Read a guess

Declare a variable for the player's guess and read it with `scanf`. The `&guess` hands `scanf` the address where it should store the number you type.

```c
int guess;
printf("Take a guess: ");
scanf("%d", &guess);
printf("You guessed %d\n", guess);
```

## Step 4 — Compare and give a hint

Use `if/else if/else` to compare the guess against the secret and print the right hint. Try this logic on its own before wrapping it in a loop.

```c
if (guess < secret) {
    printf("Too low!\n");
} else if (guess > secret) {
    printf("Too high!\n");
} else {
    printf("You got it!\n");
}
```

## Step 5 — Loop until correct (complete program)

Wrap the guessing in a `while` loop that keeps going until the guess matches. Here's the full, runnable program.

```c
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

int main(void) {
    srand(time(NULL));
    int secret = rand() % 100 + 1;
    int guess = 0;

    printf("Welcome to the Number Guessing Game!\n");
    printf("I'm thinking of a number between 1 and 100.\n");

    while (guess != secret) {
        printf("Take a guess: ");
        scanf("%d", &guess);

        if (guess < secret) {
            printf("Too low! Try again.\n");
        } else if (guess > secret) {
            printf("Too high! Try again.\n");
        } else {
            printf("Congratulations! You found it!\n");
        }
    }

    return 0;
}
```

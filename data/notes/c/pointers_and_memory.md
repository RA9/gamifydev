# Pointers and Memory

This is the lesson C is famous for. Pointers have a fearsome reputation, but that's only because most explanations skip the mental picture. Once you *see* what a pointer really is — a memory address, a "house number" — the fog lifts and C's superpower becomes clear. Breathe easy: we'll build the intuition piece by piece, and Pixel is right beside you.

## Memory is a street of houses

Your computer's memory is one enormous row of storage slots, and every slot has a unique numbered **address**, just like houses on a street. When you declare a variable, C parks its value in one of those slots.

```c
int age = 30;
```

This stores `30` somewhere in memory — say, at address `0x7ffd1234`. Normally you don't care *where* it lives; you just use the name `age`. But C lets you peek at and work with that address directly. That ability is what a **pointer** is all about.

:::analogy
Imagine a street of houses. Each house holds something (a value) and has a unique street number (an address). A *pointer* is a slip of paper with a house number written on it. The paper isn't the house — it just tells you where to find it. Follow the number, and you arrive at the value.
:::

## The two key operators: & and *

![The variable x holds 42 at address 0x1A4; the pointer p holds that address, so *p reaches the value 42](/images/lessons/c-pointer-memory.svg)

C gives you two operators to work with addresses.

The **address-of** operator `&` asks, "Where does this variable live?" It gives you the address.

The **dereference** operator `*` asks, "What's stored *at* this address?" It follows the pointer to the value.

```c
#include <stdio.h>

int main(void) {
    int age = 30;
    int *ptr = &age;     // ptr holds the ADDRESS of age

    printf("Value of age:   %d\n", age);     // 30
    printf("Address of age: %p\n", &age);    // some address
    printf("ptr points to:  %d\n", *ptr);    // 30 (follow the pointer)
    return 0;
}
```

Read `int *ptr = &age;` as: "`ptr` is a pointer to an `int`, and it stores the address of `age`." The `*` in the *declaration* marks the variable as a pointer. The `*` when *using* it means "go fetch the value there."

:::key
`&x` gives you the **address** of `x`. `*p` gives you the **value at** the address stored in `p`. The `&` packs an address; the `*` unpacks one.
:::

Because `ptr` points *at* `age`, you can even change `age` through the pointer:

```c
int age = 30;
int *ptr = &age;
*ptr = 99;            // change the value at ptr's address
printf("%d\n", age);  // 99 — age changed!
```

:::quiz
Q: If `int *p = &score;`, what does `*p` give you?
- The address of score
- The value stored in score *
- A copy of the pointer
E: Dereferencing with `*` follows the pointer and reads the value at that address — the value in `score`.
:::

This is also the answer to a mystery from an earlier lesson: `scanf("%d", &age)` passes the *address* of `age` so `scanf` knows where to store what the user types. Now it makes sense.

## The stack and the heap

C programs use memory in two main regions, and knowing the difference is essential.

The **stack** is fast, automatic, and temporary. Local variables live here. When a function is called its locals are pushed onto the stack, and when the function returns they vanish automatically. You don't manage the stack — C handles it for you.

The **heap** is a large pool of memory you manage *by hand*. You ask for a chunk when you need it and give it back when you're done. The heap is for data that must outlive a single function, or whose size you don't know until the program runs.

:::analogy
The stack is like a stack of cafeteria trays — you add a tray, use it, remove it, all in strict order, and it's cleaned up for you. The heap is like renting a storage unit: you request space, keep the key as long as you need it, and *you* are responsible for returning it. Forget to return it, and you keep paying for space you no longer use.
:::

## Asking for memory: malloc and free

To use the heap, you request memory with **`malloc`** (memory allocate) and later release it with **`free`**.

```c
#include <stdio.h>
#include <stdlib.h>   // needed for malloc and free

int main(void) {
    int *numbers = malloc(3 * sizeof(int));   // room for 3 ints

    if (numbers == NULL) {        // always check malloc succeeded
        printf("Out of memory!\n");
        return 1;
    }

    numbers[0] = 10;
    numbers[1] = 20;
    numbers[2] = 30;
    printf("%d\n", numbers[1]);   // 20

    free(numbers);                // give the memory back
    numbers = NULL;               // avoid using a freed pointer
    return 0;
}
```

A few things to notice: `sizeof(int)` asks C how many bytes an `int` takes, so `3 * sizeof(int)` is exactly enough room for three. `malloc` returns the address of the new block — or `NULL` if it failed, which you should always check.

:::tip
Set a pointer to `NULL` right after you `free` it. A `NULL` pointer is an obvious "points to nothing," which is far easier to catch than a pointer that still holds the address of memory you already released.
:::

## Why C trusts you with memory

Most modern languages hide memory management behind a "garbage collector" that cleans up automatically. C deliberately does not. It hands you direct control because that control is exactly what makes C fast and suitable for operating systems and tiny devices, where every byte and every microsecond counts.

With that power comes responsibility. Two classic dangers:

A **memory leak** happens when you `malloc` but never `free`. The memory stays reserved, unusable, until your program ends. In long-running programs, leaks pile up until you run out of memory.

A **dangling pointer** happens when you `free` memory but keep using the pointer. It now points to a slot that may hold anything — reading or writing through it causes unpredictable, hard-to-find bugs.

:::warning
The golden rule of C memory: **every `malloc` deserves exactly one `free`.** Not zero (a leak), not two (a crash). Track your allocations carefully, and free what you allocate when you're done with it.
:::

:::predict
```c
#include <stdio.h>

int main(void) {
    int x = 5;
    int *p = &x;
    *p = *p + 10;
    printf("%d\n", x);
    return 0;
}
```
- 15 *
- 5
- 10
E: `p` points to `x`. `*p = *p + 10` reads x's value (5), adds 10, and writes 15 back through the pointer into x.
:::

:::match
Q: Match each term to its meaning.
- `&` | Gets a variable's address
- `*` | Gets the value at an address
- `free` | Returns heap memory to the system
E: `&` packs an address, `*` unpacks one, and `free` releases what `malloc` gave you.
:::

## Talk about it

> Pointers gave you the ability to change a variable through its address, and the heap gave you memory you control by hand. In your own words, why might a language like C *choose* to make you manage memory manually instead of doing it for you? What do you gain, and what do you risk?

## What's next

You've faced C's most famous concept and built a real mental model for it — that's a genuine milestone. In the final lesson, **Next Steps with C**, you'll see where these fundamentals lead: arrays, strings, structs, fun project ideas, and how everything you've learned carries into the rest of your programming life. Let's finish strong.

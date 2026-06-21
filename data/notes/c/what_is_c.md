# What is C?

Welcome to the start of your C adventure. Before you write a single line, it helps to know *what* C is and *why* it has stuck around for half a century while flashier languages came and went. C is small, fast, and close to the machine — learning it is like learning how engines work before you drive. Pixel is here cheering you on, so let's dig in.

## A language that powers the world

C was created in the early 1970s by Dennis Ritchie at Bell Labs. It was built to write the Unix operating system, and that origin tells you everything about its personality: C was designed to be powerful enough to build an operating system, yet simple enough that a small team could understand the whole thing.

That decision rippled out across decades. Today, C still sits underneath an astonishing amount of the software you use every day:

- **Operating systems** — the cores (kernels) of Linux, Windows, and macOS are largely C.
- **Embedded systems** — the chip in your microwave, your car's engine controller, a smart thermostat.
- **Other languages** — Python, Ruby, and PHP were themselves originally written *in C*.
- **Performance-critical tools** — databases, game engines, and networking code.

:::analogy
Think of programming languages as vehicles. Python is a comfortable automatic car — easy and forgiving. C is a manual transmission race car: more to manage, but you feel every gear and you control the engine directly. Once you understand the manual, every other car makes more sense.
:::

When you learn C, you're not just learning a language. You're learning how computers actually work underneath the friendly surface of higher-level languages.

## Compiled vs. interpreted

Here's a core idea that shapes how C feels to use. Languages generally run in one of two ways.

An **interpreted** language (like Python or JavaScript) is read and executed line by line at run time by another program called an interpreter. You write code, you run it, the interpreter does the translating on the spot.

A **compiled** language like C works differently. Before your program can run, a tool called a **compiler** translates your entire human-readable source code into raw machine instructions the processor understands directly. That translated file is the program.

:::key
C is *compiled* ahead of time into machine code. This is the main reason C programs run so fast — there's no interpreter standing between your code and the processor while the program runs.
:::

The tradeoff: you must compile before you can run, and the compiler is strict. It will refuse to build code it doesn't understand. That strictness feels annoying at first, but it's a feature — the compiler catches whole classes of mistakes before your program ever runs.

## The compile-then-run model

Let's make this concrete. When you build a C program, two stages happen:

1. **Compile** — turn your `.c` source file into an executable program.
2. **Run** — execute that program.

The most common compiler is **gcc** (the GNU Compiler Collection). On a typical terminal, the flow looks like this:

```
gcc hello.c -o hello
./hello
```

The first line says: "Compile `hello.c` and output (`-o`) a program named `hello`." The second line runs it. If you change your code, you must compile again before the change takes effect.

:::tip
The `-o name` flag names your output program. Without it, gcc produces a file called `a.out` by default — a quirky historical name you'll see everywhere in C tutorials.
:::

## Your first C program

Tradition demands that every programmer's first program greets the world. Here is the classic, in full:

```c
#include <stdio.h>

int main(void) {
    printf("Hello, world!\n");
    return 0;
}
```

Small as it is, every line earns its place:

- `#include <stdio.h>` pulls in the **standard input/output** library so you can use `printf`.
- `int main(void)` declares the starting point. Every C program begins running at `main`.
- `printf("Hello, world!\n");` prints text. The `\n` is a newline — it moves to the next line.
- `return 0;` tells the system the program finished successfully. Zero means "all good."

Don't worry if some of this looks mysterious right now. You'll meet each piece up close in the lessons ahead.

:::predict
```c
#include <stdio.h>

int main(void) {
    printf("C is fun!\n");
    return 0;
}
```
- C is fun! *
- "C is fun!\n"
- Hello, world!
E: `printf` prints exactly the text inside the quotes. The `\n` becomes a newline (not literal characters), and the quotes aren't printed.
:::

## Where C lives today

You might wonder: with so many modern languages, why learn a language from 1972? Because C never left. It's the quiet foundation. When a company needs code that is fast, predictable, and runs on tiny hardware with almost no resources, they reach for C. Operating systems, device drivers, game engines, robotics, spacecraft software — C is there.

Learning C also makes you a stronger programmer in *any* language, because it teaches you what's really happening: how memory works, how data is stored, what "fast" actually means.

:::quiz
Q: Why do C programs typically run faster than interpreted ones?
- C code is shorter
- C is compiled to machine code ahead of time, so no interpreter runs alongside it *
- C uses more memory
E: Compilation translates C directly into machine instructions before running, removing the run-time interpreter layer.
:::

:::reorder
Arrange these into a valid minimal C program.
- #include <stdio.h>
- int main(void) {
-   printf("Hi");
-   return 0;
- }
E: Include the header so `printf` is available, open `main`, print, return success, then close the brace.
:::

## Talk about it

> You just learned that C is a compiled language sitting close to the machine. In your own words, how would you explain the difference between a *compiled* and an *interpreted* language to a friend who has never coded? What everyday analogy would you reach for?

## What's next

Now that you know what C is and why it matters, it's time to actually speak the language. In **C Basics**, you'll learn how to store information in variables, work with C's core data types, and print your own formatted output. Pixel will be right there with you — let's keep going.

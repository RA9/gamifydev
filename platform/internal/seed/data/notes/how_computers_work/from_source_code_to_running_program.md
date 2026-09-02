# From Source Code to Running Program

You've written C. You've written Java. You typed some text into a file, ran a command, and something happened. This course is about what that "something" actually is — and we start at the very beginning, with the journey from a text file you can read to a process running on real hardware. By the end of this lesson you'll be able to name every stage of that journey and say what each one produces.

## The four stages

Your source code is just text. A CPU cannot read text. So a chain of tools transforms it, one step at a time, until what's left is raw numbers the machine can obey.

```text
hello.c        ── preprocessor ──▶  hello.i     (expanded C text)
hello.i        ── compiler     ──▶  hello.s     (assembly text)
hello.s        ── assembler    ──▶  hello.o     (object code, binary)
hello.o + libs ── linker       ──▶  hello       (executable file)
hello          ── loader (OS)  ──▶  a process   (running in memory)
```

When you type `gcc hello.c -o hello`, all four tool stages run back to back and quietly delete the files in between. That's convenient, but it hides the machinery. Let's open it up.

:::analogy
Think of it like publishing a book. The preprocessor pastes in the quoted passages. The compiler translates the manuscript into another language. The assembler typesets it. The linker binds your chapters together with the reference material. The loader is the reader picking it up and starting to read.
:::

## Preprocessing: pure text substitution

The **preprocessor** runs first, and it does something surprisingly dumb: it edits your text. It doesn't understand C. It follows the lines starting with `#` and produces a new, bigger C file.

```c
#include <stdio.h>
#define GREETING "Hello, world!"

int main(void) {
    printf("%s\n", GREETING);
    return 0;
}
```

`#include <stdio.h>` means "paste the entire contents of `stdio.h` here." `#define GREETING "Hello, world!"` means "everywhere you see `GREETING` below, write that string instead." After preprocessing, there are no `#include` or `#define` lines left at all — just plain C, usually hundreds of lines longer than what you wrote.

You can see it for yourself with `gcc -E hello.c`, which stops after this stage and prints the result.

:::key
The preprocessor is a text tool, not a C tool. It expands includes and macros and strips comments. Everything after it is ordinary C source with all the pieces pasted in.
:::

## Compiling: C down to assembly

Now the real **compiler** takes over. It parses the expanded C, checks types, reports your errors, applies optimizations, and emits **assembly language** — a human-readable form of the instructions a particular CPU understands.

Running `gcc -S hello.c` stops here and leaves you a `.s` file. For a 64-bit x86 machine, `main` might come out roughly like this:

```text
main:
    push    rbp                       ; save the old frame pointer
    mov     rbp, rsp                  ; set up this function's frame
    lea     rdi, [rip + message]      ; put the string's address in rdi
    call    puts                      ; call the library print routine
    mov     eax, 0                    ; return value 0
    pop     rbp                       ; restore the frame pointer
    ret                               ; hand control back to the caller
```

Don't worry about the details — we'll unpack instructions properly in a later lesson. The point is that one line of C became several machine-level steps, and they're specific to *this* kind of CPU. Compile the same file for a phone's ARM chip and you'd get different instructions entirely.

## Assembling and linking

The **assembler** turns that assembly text into **object code**: an actual binary file (`.o`) holding the numeric encodings of those instructions. It's almost a program — but not quite. Your code calls `puts`, and the assembler has no idea where `puts` lives. It leaves a labelled hole.

The **linker** fills the holes. It takes all your object files plus the library code, matches every "I need `puts`" against every "I provide `puts`", patches in the right addresses, and writes one executable file.

This is where a very common error comes from:

:::warning
"Undefined reference to `sqrt`" is a *linker* error, not a compiler error. Your code compiled fine — the linker just couldn't find the function's actual body. On many systems that means adding `-lm` to link the math library.
:::

There are two ways to link. **Static linking** copies the library's machine code directly into your executable: the file is bigger, but it runs anywhere with no dependencies. **Dynamic linking** leaves a note saying "find `libc` at load time." The file stays small, many programs share one copy of the library in memory, and a security fix to the library helps every program at once — but the program breaks if that library is missing or the wrong version.

## Loading: how a file becomes a process

An executable is still just a file sitting on disk. When you run it, the operating system's **loader** takes over:

```text
1. Read the executable's headers.
2. Map its sections into a fresh region of memory:
     .text    the machine instructions (read-only)
     .rodata  constants, string literals
     .data    initialised globals
     .bss     zero-initialised globals (takes no file space)
3. Set up a stack for local variables, and a heap for malloc.
4. If dynamically linked, load the shared libraries and patch the addresses.
5. Set the program counter to the entry point and let the CPU go.
```

That live thing in memory — instructions, data, stack, heap, and the CPU's bookkeeping about it — is a **process**. The file is the recipe; the process is the meal being cooked.

:::tip
This is why an executable built on one operating system won't run on another even with the same CPU. The instructions might be identical, but the file format, the loader's expectations, and the system calls are all different.
:::

## Compiled, interpreted, and bytecode in between

You've now met all three of the models you've been programming in.

**Compiled ahead of time (C).** The whole pipeline above runs once, on your machine, before anyone runs the program. The result is native machine code. Fast to run, slow to build, and tied to one CPU family and OS.

**Interpreted (Python).** There's no ahead-of-time trip to machine code. An interpreter program reads your source and carries out each statement as it goes. Nothing to build, runs anywhere the interpreter runs, but there's a layer of software between your code and the CPU on every single operation — which costs speed.

**Bytecode plus a virtual machine (Java).** A middle path. `javac` compiles your `.java` to `.class` files containing **bytecode**: instructions for an imaginary CPU that doesn't physically exist. The Java Virtual Machine then executes that bytecode, and typically compiles the hot parts to real machine code while the program runs. That's why Java's slogan was "write once, run anywhere" — you ship bytecode, and each platform brings its own JVM.

:::analogy
Compiled code is a book fully translated into your language before you buy it. Interpreted code is a live interpreter standing beside you translating sentence by sentence. Bytecode is a translation into a simple shared language that any local interpreter can handle quickly.
:::

## Check Your Understanding

:::quiz
Q: Which stage of the pipeline is responsible for resolving a call to a library function like `printf` to an actual address?
- The preprocessor
- The compiler
- The linker *
- The loader
E: The compiler and assembler leave a labelled hole for `printf`. The linker matches that hole against the library that provides it and patches in the address.
:::

:::match
Q: Match each tool to what it produces.
- Preprocessor | Expanded C source with includes and macros pasted in
- Compiler | Assembly text for a specific CPU
- Assembler | Object code, a binary file with unresolved holes
- Loader | A running process in memory
E: Each stage narrows the gap between readable text and something a CPU can execute directly.
:::

:::quiz
Q: What is the main practical trade-off of dynamic linking versus static linking?
- Dynamic linking always runs faster
- Dynamic linking gives smaller files and shared library updates, but depends on the library being present *
- Static linking only works for C programs
- Static linking removes the need for a linker
E: Dynamic linking keeps executables small and lets one library update fix many programs, at the cost of a runtime dependency that can be missing or mismatched.
:::

## Talk about it

> You've now seen three ways a program can reach the CPU: compiled ahead of time, interpreted line by line, and compiled to bytecode for a virtual machine. Pick a kind of software you use every day — a game, a website, a phone app — and argue which model you'd choose to build it in, and what you'd be giving up.

## What's next

You can now follow a program from text file to living process, and you know why C, Python and Java feel so different to run. But we've been treating the CPU as a black box that "just executes" instructions. Time to open it. In **How a CPU Executes a Program**, you'll meet the fetch-decode-execute cycle, the program counter, and the clock — the heartbeat underneath everything you've ever written. Pixel is looking forward to this one.

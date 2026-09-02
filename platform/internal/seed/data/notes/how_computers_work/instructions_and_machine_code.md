# Instructions and Machine Code

The CPU fetches an instruction, decodes it, executes it. But what *is* an instruction, physically? It's not a word. It's not a line of code. It's a number — a specific pattern of bits that the control unit knows how to read. In this lesson you'll take one apart, meet assembly language as its human-readable twin, and see why the same C file becomes different bits on different machines.

## Instructions are numbers

Everything in a computer's memory is bits, and instructions are no exception. There is no separate "code memory" made of a special substance. The bytes holding your program look exactly like the bytes holding your data; the only difference is that the program counter is pointing at them.

That's a genuinely strange idea, and it's the foundation of the modern computer: **code is data**. It's what lets a compiler write a program to a file, a loader copy it into memory, and a CPU then treat those same bytes as commands.

:::analogy
Sheet music and a shopping list can be written with the same ink on the same paper. What makes one music isn't the ink — it's that a musician is reading it as music. Instructions are bytes that the CPU has been pointed at.
:::

## Opcode and operands

An instruction has to say two things: *what to do*, and *what to do it to*. So its bits are divided into fields.

The **opcode** (operation code) is the "what": add, subtract, load, store, jump. The **operands** are the "to what": which registers, which memory address, which constant.

To make this concrete, here's a small **made-up teaching instruction set** — not a real chip, just a clean example. Every instruction is exactly 16 bits, split into four 4-bit fields:

```text
 bits:  15..12    11..8     7..4      3..0
       ┌────────┬─────────┬─────────┬─────────┐
       │ opcode │  dest   │  src1   │  src2   │
       └────────┴─────────┴─────────┴─────────┘

 opcode 0001 = ADD      registers R0..R15 fit in 4 bits
 opcode 0010 = SUB
 opcode 0011 = LOAD
```

So "add R1 and R2, put the answer in R1" encodes as:

```text
  ADD   R1     R1     R2
  0001  0001   0001   0010     =  0x1112
```

The bits `0001000100010010` — the number `0x1112` — *are* the instruction. When the control unit decodes it, the first four bits switch on the adder, and the other fields select which registers to read and write.

:::key
An instruction is a number whose bit fields the CPU reads as "operation" plus "operands". Decoding is nothing more mystical than pulling those fields apart.
:::

Real instruction sets are messier than our toy. On x86-64 (the architecture in most laptops and desktops) instructions are **variable length**, from one byte up to fifteen. On AArch64 (64-bit ARM, in most phones) every instruction is exactly four bytes. Both approaches work; they trade off differently, as we'll see.

## Assembly: the same instructions, in words

Nobody writes `0x1112` by hand. **Assembly language** gives every opcode a short name (a *mnemonic*) and lets you write register names instead of bit patterns.

The crucial property is that assembly is essentially **one-to-one** with machine code. One assembly line becomes one instruction. This is not like C, where a single line might become a dozen instructions. Assembly is machine code with the numbers spelled out for humans.

```text
; x86-64, Intel syntax
mov   eax, 5      ; put the value 5 into register eax
add   eax, 3      ; add 3 to eax, leaving 8
```

The **assembler** does the translation, which is mostly a lookup and some arithmetic — much simpler than compiling. And it works in reverse: a **disassembler** turns machine code back into readable assembly, which is how people inspect programs they don't have the source for. On x86-64, for instance, the single byte `0xC3` is the `ret` instruction that returns from a function.

:::tip
You almost certainly won't write production code in assembly — compilers are better at it than humans in nearly every case. But being able to *read* a little assembly is a real skill: it's how you find out what the compiler actually did with your loop.
:::

## A tiny program in assembly

Let's take a small piece of C and write the same thing in pseudo-assembly (again, a simplified teaching notation, close in spirit to real assembly but not any specific chip).

```c
int sum = 0;
for (int i = 1; i <= 5; i++) {
    sum += i;
}
// sum is 15
```

```text
        MOV  R1, 0         ; sum  = 0
        MOV  R2, 1         ; i    = 1
loop:   CMP  R2, 5         ; compare i with 5, set the flags
        JGT  end           ; if i > 5, jump to 'end'
        ADD  R1, R1, R2    ; sum  = sum + i
        ADD  R2, R2, 1     ; i    = i + 1
        JMP  loop          ; go back to the top
end:    HALT               ; done; R1 holds 15
```

Look at what happened to the `for` loop. There is no `for` instruction on any CPU. The loop was built out of a comparison, a conditional jump, and an unconditional jump backwards. `loop:` and `end:` are **labels** — names for addresses, which the assembler replaces with real numbers.

Trace it once by hand and you'll see the whole shape of machine execution:

```text
i:   1    2    3    4    5    6
sum: 1    3    6   10   15    (i > 5, jump to end)
```

:::warning
Assembly has no types, no scope, and no safety net. `R1` is sixteen bits (or thirty-two, or sixty-four) of *something*; whether it's an integer, a character or an address exists only in the programmer's head. Everything C protects you from, assembly hands straight to you.
:::

## Instruction set architectures

The full catalogue of instructions a CPU understands — their names, their bit encodings, the registers available, how memory is addressed — is called an **instruction set architecture**, or **ISA**. x86-64, AArch64 and RISC-V are all ISAs.

An ISA is a contract. On one side, chip designers promise that any processor implementing it will run these instructions correctly. On the other, compiler writers promise to emit only instructions in the set. That contract is why a program from years ago still runs on a brand new chip, even though the silicon inside is completely redesigned.

It also explains something you may have wondered about: why the same `hello.c` produces different binaries on different machines. The C is portable; the *output* isn't. Different ISA means different opcodes, different register names, different numbers of registers, and different conventions for passing arguments to functions. A compiler for ARM and a compiler for x86-64 are aiming at genuinely different machines.

:::key
C is portable at the level of *source code*, not binaries. "Write once, compile anywhere" — recompile for each ISA. Java's bytecode took the other route: ship one binary and let each platform's virtual machine bridge the gap.
:::

## RISC and CISC

ISAs cluster into two philosophies.

**CISC** — Complex Instruction Set Computer — offers many instructions, some of them doing quite a lot of work in one go, with variable lengths and operations that can read and write memory directly. It grew up when memory was scarce and expensive, so packing more meaning into fewer bytes was a real win. x86-64 is the great survivor of this tradition.

**RISC** — Reduced Instruction Set Computer — takes the opposite bet: a smaller set of simple, fixed-width instructions, and a strict **load/store** rule where arithmetic happens only between registers and memory is touched only by dedicated load and store instructions. Simple, uniform instructions are easier to decode and easier to pipeline, so you can run them faster and burn less power. ARM and RISC-V are in this camp.

The debate is much softer than it used to be. Modern x86-64 chips decode their complex instructions internally into simpler RISC-like operations before executing them, so the two worlds have quietly borrowed from each other for a long time.

:::analogy
CISC is a kitchen appliance with a "make a whole lasagne" button. RISC is a good knife, a pan and a hob. The single button is convenient when it fits your dish exactly; the simple tools are faster and more flexible the rest of the time.
:::

## Check Your Understanding

:::quiz
Q: What is the relationship between assembly language and machine code?
- Assembly is a high-level language that compiles to many instructions per line
- Assembly is essentially a one-to-one readable notation for machine instructions *
- Assembly runs inside a virtual machine
- Assembly is machine code written in decimal instead of binary
E: Each assembly line normally corresponds to exactly one machine instruction, which is why assemblers are simple and disassembly is possible.
:::

:::fill
Q: Complete the name for the field of an instruction that says *which operation* to perform.
`An instruction is made of an ___ plus its operands.`
- opcode *
- operand
- register
- label
E: The opcode selects the operation; the operands say which registers, addresses or constants it applies to.
:::

:::quiz
Q: Why does the same C source file produce different binaries on an x86-64 laptop and an ARM phone?
- C is interpreted differently on each machine
- The two CPUs implement different instruction set architectures *
- ARM cannot run loops
- The phone compiles to Java bytecode
E: Different ISAs mean different opcodes, registers and conventions, so the compiler must generate entirely different machine code for each target.
:::

## Talk about it

> An instruction set architecture is a promise that lets hardware and software evolve independently for decades. Describe another example — from computing or anywhere else — where agreeing on a fixed interface let two sides change freely underneath. What would break if that agreement were abandoned?

## What's next

You now know what the CPU is fetching and why it looks different on different machines. But notice how much of our assembly was `LOAD` and `MOV` — shuffling data between registers and memory. That shuffling is where most real programs actually spend their time. In **Registers, RAM and the Memory Hierarchy** you'll find out why, and get an intuition for speed differences so large they're hard to believe.

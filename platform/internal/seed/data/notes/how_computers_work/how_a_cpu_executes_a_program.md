# How a CPU Executes a Program

The loader has put your program in memory and pointed the CPU at the first instruction. Now what? A processor looks like a mysterious slab of silicon, but what it does is startlingly simple and startlingly repetitive. In this lesson you'll learn the loop that every CPU runs, billions of times a second, and you'll finally have a real answer to "what does 3 GHz mean?"

## One instruction at a time

Here is the whole secret: a CPU does one small thing, then the next small thing, forever.

Not "run the program." Not even "run the line." Something much smaller: *add these two numbers*, *copy this value from memory into a register*, *if that result was zero, jump somewhere else*. Individually the steps are almost insultingly simple. There are just an enormous number of them, executed at a pace no human can picture.

:::analogy
A CPU is a fanatically obedient clerk in a tiny office. Someone hands them a card that says "add the number in box 1 to the number in box 2." They do it, put the answer in box 3, and reach for the next card. They never get bored, never improvise, and never wonder why. The intelligence is in the stack of cards, not the clerk.
:::

## The fetch-decode-execute cycle

That "reach for the next card" loop has a name — the **fetch-decode-execute cycle**, sometimes called the instruction cycle. It has three beats.

**Fetch.** The CPU reads the next instruction from memory. How does it know which one? A special register called the **program counter** (PC) holds the memory address of the next instruction. The CPU fetches from that address, then advances the PC so it points at the one after.

**Decode.** The instruction arrives as a pattern of bits. The **control unit** works out what it means: which operation is this, and which registers or memory addresses does it touch?

**Execute.** The CPU actually does it — adds, compares, loads, stores, jumps.

```text
      ┌──────────────────────────────────┐
      │                                  │
      ▼                                  │
  [ FETCH ]  read instruction at PC      │
      │      PC = PC + instruction size  │
      ▼                                  │
  [ DECODE ] what operation? which data? │
      │                                  │
      ▼                                  │
  [EXECUTE]  do it; store the result ────┘
```

Then it starts again. That's the entire job.

A **jump** or **branch** instruction is the one twist: instead of letting the PC advance normally, it *writes a new value into the PC*. That single ability is what gives you `if`, `while`, `for`, and function calls. Every loop you have ever written is, underneath, an instruction that rewinds the program counter.

:::key
The program counter is the "you are here" marker in your program. Normal instructions nudge it forward by one; branches move it somewhere else. Control flow is nothing more than editing the PC.
:::

## Inside the CPU

Three parts do most of the work.

**Registers** are the CPU's own handful of storage slots — typically a few dozen, each holding one value (often 64 bits on a modern machine). They are the *only* place the CPU can do arithmetic. Anything in RAM must be loaded into a register first. We'll come back to why they matter so much in the memory hierarchy lesson.

**The ALU** (arithmetic logic unit) is the calculator. It adds, subtracts, compares, and performs bitwise operations like AND and OR. Give it two register values and an operation code, and it produces a result plus some **flags** — small one-bit facts such as "the result was zero" or "the result overflowed." Those flags are what a conditional branch tests.

**The control unit** is the conductor. It decodes each instruction and switches on the right paths inside the chip: open this register, route it to the ALU, tell the ALU to subtract, write the answer back there.

```text
        ┌────────────── CPU ──────────────┐
        │                                 │
        │  [Control Unit] ── decodes ──┐  │
        │        │                     │  │
        │        ▼                     ▼  │
        │  [ Registers ] ◀────────▶ [ ALU ]
        │        ▲                        │
        └────────┼────────────────────────┘
                 │ load / store
                 ▼
             [   RAM   ]
```

## The clock and what GHz means

Everything inside the CPU marches to a **clock** — an electrical signal that ticks on and off at a fixed rate. Each tick is a **cycle**, and it's the beat that keeps all the parts of the chip in step: signals get a fixed window to travel and settle before the next tick captures the result.

**Hertz** means cycles per second. One gigahertz (GHz) is a billion cycles per second, so a chip described as 3 GHz ticks about three billion times each second. Modern desktop and laptop CPUs typically sit in the low single-digit GHz range.

A tempting mistake is to read that as "three billion instructions per second." It isn't. Some instructions take several cycles; meanwhile a modern CPU may finish more than one instruction per cycle by overlapping them. Cycles are the *beat*, not the *work*.

:::warning
Clock speed alone is a terrible way to compare processors. Two chips at the same GHz can differ hugely in how much work they get done per cycle. Comparing GHz across different chip designs is like comparing two engines by RPM alone.
:::

Clock rates climbed steeply for decades and then flattened out, and the reason is physics, not laziness. Roughly speaking, the power a chip burns rises with its clock rate — and every watt becomes heat that has to leave a fingernail-sized piece of silicon. Push the clock higher and you hit a **power wall**: the chip cooks itself. So designers stopped chasing speed and started chasing *width* — doing more per tick, and adding more processors.

## Pipelining: the CPU's assembly line

Here's the first way to do more per tick. Fetch, decode and execute use different parts of the chip. If you run them strictly one after another, most of the CPU sits idle at any moment.

So real CPUs **pipeline**: while instruction 1 is executing, instruction 2 is being decoded and instruction 3 is being fetched. Like a car assembly line, every station is busy on a different item.

```text
cycle:      1      2      3      4      5      6
instr 1   fetch  decode  exec
instr 2          fetch  decode  exec
instr 3                 fetch  decode  exec
instr 4                        fetch  decode  exec
```

Each individual instruction still takes three stages, but once the pipeline is full one instruction *completes* every cycle. Real pipelines have many more stages than three.

The catch is branches. When the CPU hits `if (x > 0)`, it doesn't yet know which way to go — but the pipeline needs to keep fetching *something*. So the CPU guesses, using a **branch predictor** that learns from past behaviour. Guess right and nothing is lost. Guess wrong and the partly-done work must be thrown away and the pipeline refilled, costing many cycles.

:::tip
This is one real reason predictable code can run faster than unpredictable code. A loop that branches the same way a thousand times in a row is easy to predict; one that branches randomly is not. It's rarely worth contorting your code for, but it explains some surprising benchmark results.
:::

## Cores, and what multi-core buys you

The other escape from the power wall was to stop making one processor faster and put several on the same chip. A **core** is a complete processor — its own program counter, registers, ALU and control unit — so a four-core CPU can genuinely run four streams of instructions at the same instant.

What that gives you is *throughput*, not magic speed. Four cores don't make one single-threaded program four times faster; that program still runs on one core. They help when there are four things to do at once — four programs, or one program deliberately split into parallel pieces.

And splitting is hard. Whatever part of a job must happen in order stays stubbornly serial, so it caps your speedup no matter how many cores you add. Ten cooks in a kitchen help enormously; ten cooks stirring one pot do not.

:::analogy
Making a core faster is hiring a quicker chef. Adding cores is opening more kitchens. More kitchens serve more customers — but a single dish still cooks at the speed of one kitchen.
:::

## Check Your Understanding

:::quiz
Q: What does the program counter hold?
- The number of instructions executed so far
- The address of the next instruction to fetch *
- The result of the last calculation
- How many cores are in use
E: The PC is the CPU's "you are here" marker: it stores the memory address of the next instruction, and branches work by writing a new value into it.
:::

:::fill
Q: Complete the name of the three-beat loop every CPU runs.
`fetch → ___ → execute`
- decode *
- compile
- link
- branch
E: After fetching the instruction's bits, the control unit decodes them to work out which operation to perform and on what data.
:::

:::quiz
Q: A program runs on one thread. You move it to a CPU with eight cores instead of two, same clock speed. What should you expect?
- Roughly four times faster
- Roughly eight times faster
- Roughly the same speed for that program *
- It will not run at all
E: A single-threaded program uses one core. Extra cores raise total throughput across tasks, but they can't speed up work that isn't split into parallel pieces.
:::

## Talk about it

> The fetch-decode-execute cycle is almost absurdly simple, yet it's the foundation of every piece of software you've ever used. Explain in your own words how something this simple, repeated billions of times a second, can add up to a video call or a game. Where does the apparent complexity actually come from?

## What's next

You've seen the loop and the machinery around it. But we keep saying the CPU "fetches an instruction" without ever saying what an instruction *is*. In **Instructions and Machine Code** you'll look at the actual bits — opcodes and operands — see how assembly maps one-to-one onto them, and find out why the same C file compiles into completely different instructions on a laptop and a phone.

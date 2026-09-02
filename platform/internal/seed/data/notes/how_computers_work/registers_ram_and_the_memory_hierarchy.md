# Registers, RAM and the Memory Hierarchy

When you learned about pointers, memory was one flat street of numbered houses. That picture is a good one, and it's also a polite lie. Real machines have several kinds of memory stacked on top of each other, differing in speed by factors of *millions*. This lesson gives you an intuition for those differences that will quietly improve every program you write from now on.

## Why memory has levels

A designer would love one kind of memory that is enormous, instant, and cheap. Physics and economics refuse.

Fast memory is built from circuits that hold a bit actively, using several transistors per bit. It's quick, but it's bulky and power-hungry, so you can't have much of it. Cheaper memory stores each bit as a tiny charge in a capacitor — far denser, so you can have gigabytes, but reading it takes much longer and it has to be constantly refreshed or it forgets.

There's also a limit nothing can argue with. Light travels about 30 centimetres in one nanosecond, and electrical signals in a chip travel more slowly than that. A memory chip several centimetres away from the CPU simply *cannot* answer in well under a nanosecond, no matter how much you spend. Distance costs time.

So instead of one memory, machines use many: a little bit of very fast storage close to the CPU, backed by more of something slower, backed by more of something slower still.

:::key
The memory hierarchy exists because fast, big and cheap are three things you can only pick two of. Every level trades size against speed, with the fastest sitting closest to the CPU.
:::

## Registers: the CPU's own hands

At the very top are the **registers** you met earlier. A modern CPU core has a few dozen of them, each typically holding 64 bits — a few hundred bytes in total, for the entire core.

They are not "very fast memory." They're part of the CPU itself, and they're the only place arithmetic can happen. `sum = a + b` cannot be done in RAM. The values must be loaded into registers, added there, and the result stored back. Every program is, at some level, a conversation about getting data into registers and out again.

:::analogy
Registers are your hands. The cache is the workbench in front of you. RAM is the shelf across the room. The SSD is the storage locker down the corridor. You can only actually *work* on what's in your hands.
:::

## The pyramid

Between the registers and RAM sit two or three levels of **cache** — small pools of fast memory on the CPU chip that keep copies of recently used data. Below RAM sit the storage devices, and below those, other machines.

```text
                  ▲ faster, smaller, more expensive per byte
        ┌───────────────────┐
        │     Registers     │   hundreds of bytes per core
        ├───────────────────┤
        │     L1 cache      │   tens of KB per core
        ├───────────────────┤
        │     L2 cache      │   hundreds of KB to a few MB
        ├───────────────────┤
        │     L3 cache      │   several MB, shared by cores
        ├───────────────────┤
        │       RAM         │   gigabytes
        ├───────────────────┤
        │     SSD / disk    │   hundreds of GB to terabytes
        ├───────────────────┤
        │      Network      │   effectively unlimited
        └───────────────────┘
                  ▼ slower, bigger, cheaper per byte
```

There's another line worth noticing, between RAM and the SSD. Everything above it is **volatile**: it forgets everything when the power goes. Everything below it is **persistent**. That's the real reason "save your work" exists.

## How big are the gaps?

Here are the numbers that make this whole lesson worth learning. Treat them as **rough orders of magnitude**, not facts about your machine — real figures vary a great deal between processors, memory types, drives and workloads, and they shift with every hardware generation. What matters is the *shape*, not the digits.

```text
  register access        under a nanosecond
  L1 cache hit           around 1 ns
  L2 cache hit           a few ns          (roughly 4)
  L3 cache hit           tens of ns        (roughly 10 to 40)
  main memory (RAM)      around 100 ns
  SSD read               around 100 microseconds   (100,000 ns)
  spinning disk seek     around 10 milliseconds    (10,000,000 ns)
  network round trip     milliseconds to hundreds of milliseconds
```

Read down that list slowly. Between a cache hit and a RAM access there's roughly a hundredfold gap. Between RAM and an SSD, another thousandfold. Between an SSD and a spinning disk, another hundredfold on top.

:::warning
Don't quote these numbers as measurements — they're teaching approximations. If you ever need real figures for a real decision, measure the actual machine. What you should carry away is the ratios, which have stayed roughly stable even as everything got faster.
:::

## If a register access took one second

Nanoseconds mean nothing to human intuition, so let's stretch time. Suppose one nanosecond became one second — a scale factor of a billion. Now everything happens at a pace you can feel.

```text
  L1 cache hit        1 second        glance at the note in your hand
  L2 cache hit        4 seconds       reach for the paper on your desk
  L3 cache hit        ~20 seconds     walk to the filing cabinet
  RAM access          ~2 minutes      go find it in the next room
  SSD read            ~1 day          post a letter and wait for the reply
  disk seek           ~4 months       order it from overseas by sea freight
  network round trip  years           a research expedition
```

That is the real gap. A cache miss that sends the CPU out to RAM isn't a small delay — on this scale it's the difference between glancing at your hand and walking to another room. And a program that reaches out to disk on every operation is ordering everything by sea freight.

Now the CPU's behaviour makes sense. All that machinery — pipelines, branch predictors, caches, prefetchers — exists mostly to avoid standing around waiting for memory. A modern processor spends a startling share of its life doing exactly that.

:::tip
When a program is unexpectedly slow, "it's doing too much arithmetic" is rarely the answer. "It's waiting on memory or on I/O" very often is. Learning to think about *where the data lives* is one of the highest-value habits in performance work.
:::

## What this means for the code you write

You don't control the cache directly, and you shouldn't try. But three habits follow naturally from the pyramid.

**Touch less data.** Halving the data your loop walks over often does more for speed than halving the arithmetic.

**Touch it in order.** The hierarchy is designed to reward code that reads memory in a predictable, sequential pattern — the next lesson explains exactly why.

**Batch your trips downward.** Reading a file in one large chunk beats a thousand tiny reads, for the same reason you'd rather make one trip to the shops than a hundred.

```c
/* Two loops, same arithmetic, very different memory behaviour. */

long total = 0;
for (int i = 0; i < n; i++)
    total += data[i];          /* walks straight through memory */

long total2 = 0;
for (int i = 0; i < n; i++)
    total2 += data[order[i]];  /* jumps around; every step may miss */
```

Both loops do `n` additions. On a large array the second can be many times slower, and not one line of arithmetic explains why. The difference is entirely in where the data had to come from.

## Check Your Understanding

:::quiz
Q: Why do computers use a hierarchy of memory instead of one large, fast memory?
- Operating systems require at least three levels
- Fast memory is expensive and physically limited in size, so only a small amount can sit close to the CPU *
- RAM cannot store program instructions
- It makes programs easier to write
E: Speed, capacity and cost pull against each other, and signals take time to travel, so designers stack a little fast memory on top of a lot of slow memory.
:::

:::match
Q: Match each level to its rough scale (approximate, and varies by machine).
- Registers | Hundreds of bytes, accessed in well under a nanosecond
- L1 cache | Tens of kilobytes, accessed in about a nanosecond
- RAM | Gigabytes, accessed in roughly a hundred nanoseconds
- SSD | Hundreds of gigabytes or more, accessed in microseconds
E: Each step down the pyramid buys roughly two or three orders of magnitude more capacity at a similar cost in latency.
:::

:::quiz
Q: A loop performs the same number of additions in both versions, but one reads array elements in a random order and is far slower. What is the most likely explanation?
- Random numbers are expensive to add
- The random version misses in cache and waits on RAM far more often *
- The compiler refuses to optimise random access
- The CPU runs at a lower clock speed for random data
E: Sequential access is served from cache most of the time; jumping around defeats that, so the CPU repeatedly stalls waiting for main memory.
:::

## Talk about it

> The "one second per nanosecond" scale turns invisible delays into human distances. Pick a program you use often and describe, in that stretched time, what it might be doing while you wait for it — and where you'd guess the time actually goes.

## What's next

You now know that the gap between cache and RAM is the difference between a glance and a walk to another room. In **CPU Cache** you'll see exactly how the machine tries to win that bet on your behalf, why reading a 2D array one way can be dramatically faster than reading it the other way, and why arrays so often outrun linked lists in practice even when the Big O says they shouldn't.

# CPU Cache

Last lesson ended on an uncomfortable fact: reaching into RAM costs the CPU something like a hundred times what a cache hit costs. If every memory access paid that price, modern processors would spend nearly all their time waiting. They don't — because of the cache. This lesson shows you how the cache works, and how to write code that cooperates with it instead of fighting it.

## The bet the cache makes

The **cache** is a small amount of fast memory on the CPU chip that holds copies of data from RAM. When the CPU needs a value, it checks the cache first. If the value is there — a **hit** — it comes back almost immediately. If not — a **miss** — the CPU has to go out to RAM and wait.

The cache can't hold everything; it's thousands of times smaller than RAM. So it makes a bet: *the data you're about to use is probably data you've used recently, or data that sits next to it.*

That bet sounds flimsy. In practice it's astonishingly good, because of how real programs behave. Loops run over the same variables again and again. Arrays get walked from start to finish. Functions call the same helpers. Programs are creatures of habit.

:::analogy
A cache is the small pile of books on your desk, next to a library of thousands. You keep on the desk what you're reading now and what you read five minutes ago. Almost every time you reach for something, it's already there — and when it isn't, you walk to the shelves and accept the delay.
:::

## Cache lines: the neighbours come free

Here's the first thing that surprises people. The cache never fetches a single byte. It transfers memory in fixed-size blocks called **cache lines**, commonly 64 bytes on modern processors.

Ask for one `int` and the cache brings back the entire aligned 64-byte block containing it. With 4-byte integers, that's 16 of them:

```text
you asked for a[0]; the cache fetched all of this

 ┌──────────────────── one 64-byte cache line ────────────────────┐
 │ a[0] a[1] a[2] a[3] a[4] a[5] a[6] a[7] ... a[14] a[15]        │
 └────────────────────────────────────────────────────────────────┘
    ▲
    the miss you paid for       the other 15 are now free
```

So the *next* fifteen reads, if you go in order, are hits. One slow trip bought you sixteen fast accesses. Walk the array backwards and it works just as well. Jump around by 64 bytes or more and every single read pays full price.

:::key
Memory moves in cache-line-sized blocks, not bytes. Data that lives near the data you're using effectively comes along for free — which is why *how you lay data out* can matter as much as how much of it there is.
:::

## Temporal and spatial locality

Those two habits of real programs have names, and they're worth knowing.

**Temporal locality**: if you used something, you'll probably use it again soon. A loop counter, an accumulator, the head of a list you're rebuilding. Caches exploit this by keeping recent data around.

**Spatial locality**: if you used something, you'll probably use its neighbours soon. Walking an array, reading the fields of a struct, executing consecutive instructions. Caches exploit this with cache lines, and with **prefetching** — hardware that spots a steady forward march through memory and starts fetching the next lines before you ask.

Code with good locality gets the machine's help for free. Code with poor locality gets none of it, and no amount of clever arithmetic will make up the difference.

## Hits, misses and eviction

The cache is full almost all the time. So when a new line arrives, an old one must go. That's **eviction**, and caches generally aim to throw out something that hasn't been used recently — the least useful guess about the future.

This gives programs a property that's easy to miss: performance can fall off a cliff once your **working set** — the data you're actively touching — grows past what the cache can hold. Below that size, nearly everything is a hit. Above it, lines you'll need again keep getting evicted before you get back to them, and the hit rate collapses. The code didn't change; the data outgrew the desk.

:::warning
This is why a benchmark on a small array can be wildly misleading. A loop that looks fast on ten thousand elements may behave completely differently on ten million, because the small case was living entirely in cache. Always test at realistic sizes.
:::

## Two ways to walk a grid

Now the demonstration that makes all of this concrete. In C, a 2D array is stored **row-major**: the whole of row 0 sits in memory, then the whole of row 1, and so on. `a[i][j]` and `a[i][j+1]` are adjacent; `a[i][j]` and `a[i+1][j]` are a full row apart.

```c
#define N 4096
static int a[N][N];    /* one row = 4096 ints = 16 KB */

/* Version 1: j innermost — walks straight through memory */
long sum = 0;
for (int i = 0; i < N; i++)
    for (int j = 0; j < N; j++)
        sum += a[i][j];

/* Version 2: i innermost — jumps 16 KB every step */
long sum = 0;
for (int j = 0; j < N; j++)
    for (int i = 0; i < N; i++)
        sum += a[i][j];
```

Identical arithmetic. Identical number of additions. Identical answer. Only the loop order differs — and on a large array the second version is commonly *several times* slower. Measure it yourself on a big enough grid and the gap is impossible to miss.

Here's why, in memory order:

```text
Version 1 (fast):  a[0][0] a[0][1] a[0][2] ...   all in one cache line
                   ├── 1 miss, then 15 hits ──┤

Version 2 (slow):  a[0][0] ......16 KB...... a[1][0] ......16 KB...... a[2][0]
                   ├─ miss ─┤                ├─ miss ─┤                ├─ miss ─┤
```

Version 2 fetches a whole 64-byte line, uses 4 bytes of it, and moves so far away that the line is long evicted before it comes back. It throws away fifteen sixteenths of every transfer.

:::tip
The rule of thumb: **make the innermost loop move fastest through memory.** In C and C++ that means the rightmost index varies innermost. Some other languages store arrays column-major, so check before you assume — but the principle is the same everywhere.
:::

## Laying your data out

The same thinking applies to structs. Suppose you have a million particles and a physics step that only needs their positions.

```c
/* Array of structs: fields interleaved */
struct Particle { float x, y, z; float r, g, b; };
struct Particle particles[N];

/* Struct of arrays: fields kept together by kind */
struct Particles {
    float x[N], y[N], z[N];
    float r[N], g[N], b[N];
} particles;
```

With the array of structs, reading every `x` drags all the colour data through the cache too, because it shares the cache lines. With the struct of arrays, `x[0..15]` fills a line completely and nothing is wasted. Neither layout is universally right — if you usually touch *all* the fields of one particle at a time, the array of structs wins — but it's a choice worth making deliberately rather than by accident.

This also settles an old puzzle from your data structures course. Walking an array and walking a linked list are both O(n), yet the array is often dramatically faster in practice. An array's elements are adjacent, so the prefetcher can race ahead. A linked list's nodes may be scattered anywhere in the heap, and worse, the CPU can't fetch the next node until it has *read* the current node's pointer — a chain of dependent misses with nothing to overlap them. Big O counts operations; it doesn't count how far each one has to reach.

:::key
Big O tells you how work grows. The memory hierarchy tells you what each unit of work actually costs. You need both — an O(n) algorithm with terrible locality can lose to an O(n log n) one with excellent locality on real data.
:::

## Check Your Understanding

:::quiz
Q: You read a single `int` that isn't in cache. What does the CPU actually fetch from memory?
- Exactly those 4 bytes
- An entire cache line, commonly 64 bytes, containing that int *
- The whole array
- One byte at a time until the int is complete
E: Caches transfer fixed-size lines. You pay for the whole line, which is why the neighbouring values are effectively free afterwards.
:::

:::fill
Q: Complete the index that must vary in the innermost loop to walk a C 2D array in memory order.
`for (int i = 0; i < N; i++) for (int j = 0; j < N; j++) sum += a[i][___];`
- j *
- i
- N
- 0
E: C stores rows contiguously, so the rightmost index has to vary fastest. That marches straight through memory and uses every byte of each cache line.
:::

:::quiz
Q: Traversing an array and traversing a linked list are both O(n). Why is the array usually faster in practice?
- Arrays use fewer instructions per element in every case
- Array elements are contiguous, so cache lines and prefetching help; list nodes are scattered and each step depends on the previous load *
- Linked lists are O(n²)
- Compilers refuse to optimise linked lists
E: The array gets spatial locality and prefetching; the list produces a chain of dependent, unpredictable misses that the hardware can't hide.
:::

## Talk about it

> You've now seen two loops with identical Big O and identical arithmetic differ by a large factor purely because of memory layout. Does that mean Big O analysis was a waste of time? Make the case for how the two ways of thinking fit together, and when you'd reach for each.

## What's next

You've been trusting the machine to add, compare and count this whole course. In **How Computers Calculate** we go right down to the bottom: binary place value, hexadecimal, how negative numbers are represented with two's complement, what happens when a number overflows, and how a handful of logic gates can add two numbers together at all.

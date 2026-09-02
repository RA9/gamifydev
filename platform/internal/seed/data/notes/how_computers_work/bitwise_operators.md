# Bitwise Operators

You've spent this course learning that everything is bits. Now you get the tools to reach in and work with them directly. Bitwise operators are how you pack eight yes/no answers into a single byte, pull the red channel out of a colour, or test one flag among dozens — and they do it in a single machine instruction. This is the last lesson of the course, and it's a satisfying one: pure, practical, and built entirely on what you already know.

## Working inside a byte

Ordinary operators treat a number as one whole value: `a + b` adds the numbers. **Bitwise** operators ignore the number entirely and work on each bit position independently. `5 + 3` is 8; `5 & 3` is 1, because the operator lined up their bits and compared column by column:

```text
   5  =  0101
   3  =  0011
       ────────
   &  =  0001   =  1
```

There are six operators, and each maps to a single instruction the ALU already knows how to perform. They are among the fastest operations a computer can do.

:::key
Bitwise operators work column by column on bit positions, with no carrying between columns. That's what makes them different from arithmetic — and what makes them so useful for packing and unpacking data.
:::

## AND, OR, XOR and NOT

Four operators combine bits. Their truth tables are the whole definition:

```text
   a  b │ a & b   a | b   a ^ b        a │ ~a
  ──────┼──────────────────────       ───┼────
   0  0 │   0       0       0          0 │  1
   0  1 │   0       1       1          1 │  0
   1  0 │   0       1       1
   1  1 │   1       1       0
```

**`&` AND** gives 1 only when *both* bits are 1 — a filter that keeps a bit only if the other operand permits it. **`|` OR** gives 1 when *either* bit is 1, so it turns bits on. **`^` XOR** (exclusive or) gives 1 when the bits *differ*; remember its two magic properties, `x ^ x == 0` and `x ^ 0 == x`. **`~` NOT** flips every bit, and takes one operand rather than two.

Here they all are on real values:

```c
#include <stdio.h>

int main(void) {
    int a = 12;      /* 1100 */
    int b = 10;      /* 1010 */

    printf("%d\n", a & b);   /* 8   — 1000 */
    printf("%d\n", a | b);   /* 14  — 1110 */
    printf("%d\n", a ^ b);   /* 6   — 0110 */
    printf("%d\n", ~a);      /* -13 — every bit flipped */
    return 0;
}
```

That `~12 == -13` is not a typo, and you already know why: flipping all the bits of 12 and reading the result as two's complement gives −13. Inverting a bit pattern and negating a number are related but *not* the same operation.

:::warning
Don't confuse `&` with `&&`, or `|` with `||`. The doubled versions are logical operators working on true/false as a whole; the single versions work bit by bit. `1 & 2` is `0` while `1 && 2` is true — a bug that compiles cleanly and behaves strangely.
:::

## Shifting left and right

Two more operators slide bits sideways. **`<<` left shift** moves bits toward the high end, filling with zeros on the right; **`>>` right shift** moves them toward the low end.

```text
   5 << 1                          20 >> 2
   00000101                        00010100
   ────────                        ────────
   00001010  = 10                  00000101  = 5
```

Because each column is worth twice the one to its right, shifting left by `n` multiplies by 2ⁿ and shifting right by `n` divides by 2ⁿ. `x << 3` is `x * 8`; `x >> 1` is `x / 2`.

There's a subtlety with right-shifting *signed* numbers. If the top bit is 1 — a negative number — what gets shifted in at the top? An **arithmetic** shift copies the sign bit in, preserving the sign, so `-8 >> 1` gives `-4`. A **logical** shift brings in zeros, which turns a negative number into a large positive one.

In C, right-shifting a negative signed value is implementation-defined: compilers overwhelmingly do the arithmetic shift, but the standard doesn't force it. Shift unsigned types when you're manipulating bit patterns and you sidestep the question entirely. Java made the distinction explicit with two operators — `>>` is arithmetic, `>>>` is logical.

:::tip
Write `x * 8` when you mean multiplication and `x << 3` when you mean "shift the bit pattern". Compilers have converted multiplication by powers of two into shifts for decades — you gain nothing but confusion by doing it by hand. Also note that `-7 >> 1` is `-4` while `-7 / 2` is `-3` in C, so they are not interchangeable for negatives.
:::

## Masks: set, clear, toggle, test

A **mask** is a value whose bits mark the positions you care about. The building block is `1 << n` — a value with a single 1 in position `n`, so `1 << 0` is `00000001` and `1 << 3` is `00001000`. From that come the four essential moves:

```c
x |=  (1u << n);        /* SET    bit n to 1                       */
x &= ~(1u << n);        /* CLEAR  bit n to 0                       */
x ^=  (1u << n);        /* TOGGLE bit n                            */
if (x & (1u << n)) { }  /* TEST   bit n — nonzero means it is set  */
```

Read them against the truth tables and each one is obvious. OR with a 1 forces that position on; AND with a mask of all 1s except one 0 forces it off; XOR with a 1 flips it; AND with a single 1 isolates it. In every case the other bits are untouched.

Masks also extract whole *fields*. A colour packed as `0xRRGGBB` is three 8-bit channels in one integer — shift the field you want down to the bottom, then mask off everything above it:

```c
unsigned int colour = 0xFF8800;

unsigned int r = (colour >> 16) & 0xFF;   /* 255 */
unsigned int g = (colour >>  8) & 0xFF;   /* 136 */
unsigned int b =  colour        & 0xFF;   /*   0 */
```

Shift, then mask. That pattern turns up everywhere binary data is packed: network headers, file formats, hardware registers, instruction encodings — including the toy instruction format from earlier in this course.

## Flags packed into one integer

Here's where bitwise operators earn their keep in everyday code. When something has several independent on/off properties, one integer with one bit per property is compact, fast to test, and trivial to pass around.

```c
#define PERM_READ    (1u << 0)   /* 001 */
#define PERM_WRITE   (1u << 1)   /* 010 */
#define PERM_EXEC    (1u << 2)   /* 100 */

unsigned int perms = PERM_READ | PERM_WRITE;   /* 011 */

if (perms & PERM_WRITE) puts("can write");     /* prints */
if (perms & PERM_EXEC)  puts("can execute");   /* does not print */

perms &= ~PERM_WRITE;                          /* revoke write */
if (perms & PERM_WRITE) puts("still writable");/* does not print */
```

Notice how naturally it reads once you know the idiom: `|` to grant, `& ~` to revoke, `&` to check. This is exactly how file permissions, event masks and configuration options are represented across enormous amounts of real software — you've been using these APIs without seeing the mechanism.

:::analogy
A flags integer is a row of light switches on one panel. OR flips specific switches on, AND-with-NOT flips them off, and AND lets you glance at just one. The panel is a single number you can store, copy or send in one go.
:::

## Clever tricks and their price

XOR's two properties — `x ^ x == 0` and `x ^ 0 == x` — enable some genuinely delightful tricks.

**Swapping without a temporary variable:**

```c
int a = 3, b = 7;
a ^= b;        /* a = 3^7          */
b ^= a;        /* b = 7^(3^7) = 3  */
a ^= b;        /* a = (3^7)^3 = 7  */
/* a is 7, b is 3 */
```

It works because XOR is its own inverse. It's also no faster than a temporary variable on any modern machine, and it fails catastrophically if `a` and `b` are the same variable, which zeroes it. A lovely puzzle; a poor habit.

**Finding the unique element.** Given a list where every value appears exactly twice except one, XOR everything together. The pairs cancel to zero and the loner survives:

```python
def find_unique(nums):
    result = 0
    for n in nums:
        result ^= n
    return result

print(find_unique([4, 1, 2, 1, 2]))   # 4
```

This one is genuinely excellent: O(n) time, O(1) space, and order doesn't matter because XOR is commutative and associative.

But temper the enthusiasm. Bit tricks are dense, and dense code is expensive to read. `perms & PERM_WRITE` is clear to anyone; a chain of shifts and masks replacing a simple loop is a puzzle your future self has to re-solve at the worst possible moment. Reach for these when they express the idea *directly* — packed data, flags, hardware fields — and leave them alone when they're merely a shortcut for something a compiler would optimise anyway.

:::warning
Watch the edges. `1 << 31` on a 32-bit signed `int` is undefined behaviour in C — write `1u << 31`. Shifting by an amount greater than or equal to the type's width is undefined too. When you're working with bit patterns rather than quantities, use unsigned types and you avoid most of these traps.
:::

## Check Your Understanding

:::predict
```c
#include <stdio.h>
int main(void) {
    unsigned int x = 0xA;      /* 1010 in binary */
    x |= (1u << 0);            /* set bit 0 */
    x &= ~(1u << 3);           /* clear bit 3 */
    printf("%u\n", x);
    return 0;
}
```
- 3 *
- 11
- 10
- 9
E: Start with 1010. Setting bit 0 gives 1011. Clearing bit 3 gives 0011, which is 3.
:::

:::quiz
Q: Why does XORing every number in a list where all values appear twice except one leave you with the unique value?
- XOR sorts the values as it goes
- Each pair cancels to 0 because `x ^ x == 0`, and `0 ^ y == y` leaves the loner *
- XOR adds the values and subtracts the duplicates
- It only works if the list is sorted
E: XOR is commutative and associative, so pairs cancel wherever they sit in the list, and the single unmatched value XORed with zero is itself.
:::

:::match
Q: Match each bit operation to its idiom.
- Set bit n | `x |= (1u << n)`
- Clear bit n | `x &= ~(1u << n)`
- Toggle bit n | `x ^= (1u << n)`
- Test bit n | `x & (1u << n)`
E: OR turns bits on, AND with an inverted mask turns them off, XOR flips them, and AND with a single-bit mask isolates one for testing.
:::

## Talk about it

> Bit tricks sit right on the line between elegant and unreadable. Pick one from this lesson that you'd happily put in code you share with others, and one you'd avoid, and explain what makes the difference. Is "the compiler can't do this for me" a good test?

## What's next

Take a moment here, because you've just finished something substantial: the last lesson of **How Computers Work**, and the last lesson of the whole **Computer Science Foundations** path.

Look at what you can now do. In **Data Structures** you learned how information is organised — arrays, lists, stacks, queues, hash tables, trees — and when each one earns its place. In **Complexity and Analysis** you learned to reason about cost before writing a line, and to tell an O(n log n) idea from an O(n²) one on sight. In **Algorithms** you learned the classic techniques for searching, sorting and traversing, and how to recognise the shape of a problem you've met before. And in **How Computers Work** you followed a program from text file to running process, through the fetch-decode-execute cycle, down past the cache into RAM, and all the way to the bits, the gates and the two's complement adder underneath.

That combination is what separates someone who can make code work from someone who understands *why* it works — the difference between guessing at a performance problem and reasoning about it. You've earned that.

Where next? Four directions build straight on this foundation. **Operating systems** picks up where the loader left off: processes, threads, scheduling, virtual memory, and how one machine convincingly pretends to be many. **Networking** turns this course's byte-order and encoding work into protocols and sockets. **Databases** applies the memory hierarchy at a much larger scale, where "the disk is slow" shapes every design decision. And **system design** puts it all together at the level of whole systems and real traffic.

Whichever you pick, you're not starting from scratch any more — you're extending a model of the machine that you actually built. Pixel has been beside you through every lesson of this path, and is genuinely proud of how far you've come. Go build something.

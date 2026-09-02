# How Computers Calculate

We've said "the ALU adds two numbers" as though that were obvious. It isn't. A computer has no concept of the number seven — it has voltages, high and low, on and off. Everything else is a convention built on top. In this lesson you'll learn that convention: how numbers are written in binary, how negative numbers work, what happens when a number runs out of room, and how a few logic gates can add.

## Why binary at all

Building a circuit that reliably tells apart ten different voltage levels is hard. Building one that tells apart *two* — roughly high, roughly low — is easy, cheap, and survives noise, heat and manufacturing wobble.

So computers use two symbols. One **bit** is a single high-or-low, written 0 or 1. Eight bits make a **byte**, which can hold 2⁸ = 256 different patterns.

Nothing about a bit says "number." A pattern of bits means whatever the program decides it means: a number, a letter, a colour, a pixel, an instruction. Numbers are just the first and most important convention.

:::key
Bits carry no meaning on their own. `01000001` is the number 65, the letter `A`, or part of an instruction — depending entirely on how the code chooses to read it.
:::

## Place value: counting in binary

You already know place value; you just know it in tens. In `156`, the 1 means a hundred, the 5 means fifty, the 6 means six. Each column is worth ten times the one to its right.

Binary is the same idea with two instead of ten. Each column is worth twice the one to its right.

```text
 column value: 128  64  32  16   8   4   2   1
        bits:    1   0   0   1   1   1   0   0

        128 + 16 + 8 + 4  =  156
```

To go **binary → decimal**, add up the column values wherever there's a 1. That's all.

To go **decimal → binary**, walk the columns from the largest down, taking each one that fits:

```text
156 − 128 = 28   →  1 in the 128 column
 28 − 64  ?  no  →  0 in the 64 column
 28 − 32  ?  no  →  0 in the 32 column
 28 − 16  = 12   →  1 in the 16 column
 12 − 8   =  4   →  1 in the 8 column
  4 − 4   =  0   →  1 in the 4 column
  nothing left   →  0, 0 in the 2 and 1 columns

156 = 10011100
```

Counting up works exactly like decimal, carrying when a column fills — `0, 1, 10, 11, 100, 101, 110, 111, 1000` is zero through eight.

:::tip
Learn the powers of two by sight — 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024. You'll meet them constantly: 8-bit values top out at 255, a byte has 256 patterns, and 1024 is why kilobytes are the size they are.
:::

## Hexadecimal: binary for humans

Binary is correct and unreadable. `10011100` is easy to miscount and painful to type. So programmers write bits in **hexadecimal** — base 16, digits `0`–`9` then `A`–`F` for ten through fifteen.

Hex works because **one hex digit is exactly four bits**. To convert, split the bits into groups of four and translate each independently:

```text
   1001   1100
     9      C      →  0x9C

   0x9C = 9 × 16 + 12 = 144 + 12 = 156   ✓
```

No arithmetic, no carrying between groups — just a lookup. That's why memory addresses, colours (`#FF8800`), and byte dumps are all written in hex. The `0x` prefix is how C and most other languages say "this is hexadecimal."

```c
int a = 156;      /* decimal      */
int b = 0x9C;     /* hexadecimal  */
int c = 0b10011100; /* binary — supported by many modern compilers */
/* a == b == c */
```

:::analogy
Hex is shorthand, the way "3 dozen" is shorthand for 36. It isn't a different number system for the computer — the machine still sees bits. Hex just packs four bits into one symbol your eye can take in at a glance.
:::

## Negative numbers: two's complement

Here's a real puzzle. A byte holds 256 patterns. If they're all positive, you get 0 to 255. But you need negatives too — and there's no minus sign in hardware.

The naive fix is to reserve the top bit as a sign. It sort of works, but it gives you two different zeros (`00000000` and `10000000`) and it complicates the adder, because adding a positive to a negative would need different logic from adding two positives. Real machines use **two's complement** instead, and it is one of the most elegant tricks in computing.

To negate a number: **invert every bit, then add 1.**

```text
  5   =  00000101
invert:  11111010
  +1:    11111011   =  −5
```

Now look what happens when you add 7 and −5, using nothing but ordinary binary addition:

```text
    00000111        7
  + 11111011      −5
  ───────────
  1 00000010

  the 9th bit falls off the end of the byte
  what's left:  00000010  =  2      ✓
```

The plain adder got it right, with no special case at all. **That's the whole point of two's complement**: subtraction becomes addition, and one adder circuit handles every combination of signs. The carry out of the top just falls off the edge.

Reading a two's complement byte is easy: if the top bit is 0 it's a normal positive number; if the top bit is 1, negate it (invert and add 1) to find its magnitude.

The range is `−128` to `+127`, and the asymmetry bothers people until they count the patterns:

```text
  00000000            0
  00000001 .. 01111111    +1 .. +127     (127 patterns)
  10000000 .. 11111111   −128 .. −1      (128 patterns)
```

Zero has to live somewhere, and it takes up a slot on the positive side. So there's one more negative number than positive. For 32-bit `int`, the same logic gives −2,147,483,648 to +2,147,483,647.

:::warning
That extra negative has a sharp edge: for 8 bits, `−(−128)` doesn't fit. Invert `10000000` to get `01111111`, add 1, and you're back at `10000000` — still −128. Negating the most negative value overflows at every width, which has caused real bugs in real software.
:::

## Overflow and wraparound

A fixed number of bits means a fixed range, and going past the end doesn't raise an alarm — it wraps around, like an odometer rolling over.

```c
#include <stdio.h>

int main(void) {
    unsigned char u = 255;   /* 11111111, the largest 8-bit value */
    u = u + 1;               /* carry falls off the end */
    printf("%u\n", u);
    /* 0 */
    return 0;
}
```

Unsigned arithmetic in C is *defined* to wrap like this: results are taken modulo 2ⁿ. 255 + 1 becomes 0, and 0 − 1 becomes 255.

Signed overflow is a different and nastier story. On typical hardware `127 + 1` in a byte lands on `10000000`, which reads as −128 — the biggest positive number suddenly becoming the most negative. But in C, overflowing a signed integer is **undefined behaviour**: the standard doesn't promise wraparound, and optimisers are allowed to assume it never happens. Code that relies on it can behave differently at different optimisation levels.

:::tip
If you need values that might exceed a type's range, either pick a wider type, check before you compute (`if (a > INT_MAX - b)` rather than `if (a + b < 0)`), or use unsigned types where wraparound is well defined. Never test for overflow *after* it happens in signed arithmetic.
:::

## Adding with logic gates

One last question: how does the hardware add at all? With **logic gates** — tiny circuits that combine bits.

`AND` outputs 1 only if both inputs are 1. `OR` outputs 1 if either is. `XOR` (exclusive or) outputs 1 if the inputs *differ*. Now think about adding two single bits:

```text
   A   B  │  sum  carry
  ────────┼─────────────
   0   0  │   0     0
   0   1  │   1     0
   1   0  │   1     0
   1   1  │   0     1        1 + 1 = 10 in binary
```

Compare the columns. The `sum` column is exactly `A XOR B`. The `carry` column is exactly `A AND B`. Two gates, and you have a **half adder** — correct addition of two bits.

To add multi-bit numbers you also need to accept a carry coming *in* from the column to the right. That's a **full adder**: three inputs, and `sum = A XOR B XOR carry-in`, with a carry out whenever at least two of the three inputs are 1.

Chain eight full adders together, each one's carry-out feeding the next one's carry-in, and you have a circuit that adds two bytes:

```text
  bit7   bit6   bit5   ...   bit1   bit0
   ▲      ▲      ▲            ▲      ▲
  [FA]◀──[FA]◀──[FA]◀── ... ─[FA]◀──[FA]◀─ carry in = 0
   │      │      │            │      │
  sum7   sum6   sum5   ...   sum1   sum0
```

That's the heart of the ALU. Subtraction reuses it via two's complement. Multiplication is built from shifts and additions. Every calculation your programs have ever done rests on this handful of gates.

## Check Your Understanding

:::quiz
Q: In 8-bit two's complement, how do you find the representation of −5?
- Set the top bit of 5 to get 10000101
- Invert all the bits of 5, then add 1 *
- Subtract 5 from 255
- Write 5 and mark it negative in a separate flag
E: Invert `00000101` to `11111010`, add 1, and you get `11111011`. This is what makes ordinary addition handle subtraction correctly.
:::

:::predict
```c
#include <stdio.h>
int main(void) {
    unsigned char u = 250;
    u = u + 10;
    printf("%u\n", u);
    return 0;
}
```
- 4 *
- 260
- 255
- 0
E: Unsigned arithmetic wraps modulo 256. 250 + 10 = 260, and 260 − 256 = 4.
:::

:::quiz
Q: Why does an 8-bit signed range run from −128 to 127 rather than −127 to 127?
- One pattern is reserved for errors
- Zero occupies one of the 256 patterns, leaving 127 positives and 128 negatives *
- The top bit cannot be used
- Because 128 is not a power of two
E: 256 patterns must cover zero as well, and in two's complement zero sits on the positive side, so there is one more negative value than positive.
:::

## Talk about it

> Two's complement was chosen because it lets a single adder circuit handle both addition and subtraction, with no special cases. Describe another design decision — in software or anywhere else — where the "clever" choice was really about removing special cases. Why is having fewer cases so often worth more than being obvious?

## What's next

You can now read the bits of a number and follow them all the way down to the gates. But bits get used for far more than numbers — most obviously, for text. In **Character Encodings: ASCII and Unicode** you'll find out why the same file can open as gibberish on another machine, what a code point really is, and why the length of a string is a much trickier question than it looks.

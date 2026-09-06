# Floating Point Math

Evaluate `0.1 + 0.2` in C and print enough significant digits, and you get `0.30000000000000004`. That's not a bug in C or your machine — it's the honest answer given how decimals are stored. This lesson explains where that trailing `4` comes from, and turns it from an unsettling curiosity into a set of habits that will keep real bugs out of your programs.

## Inside a float

Integers are stored exactly: a fixed pattern of bits, one pattern per value. Real numbers can't work that way — there are infinitely many of them between 0 and 1 alone, and you have 32 or 64 bits. So floating point makes a trade: cover an enormous range, at limited precision.

The **IEEE 754** standard splits the bits into three fields, in the spirit of scientific notation:

```text
  32-bit float ("single precision")

   ┌─┬────────────┬───────────────────────────────┐
   │s│  exponent  │           fraction            │
   └─┴────────────┴───────────────────────────────┘
    1     8 bits             23 bits

  64-bit double ("double precision")
    1 sign  +  11 exponent  +  52 fraction
```

The **sign** bit is 0 for positive, 1 for negative. The **exponent** says how far to slide the binary point — that's the "floating" part, and it's what buys the huge range. The **fraction** (also called the mantissa or significand) holds the digits, and it's what limits the precision.

Roughly: a `float` gives about 7 decimal digits of precision, a `double` about 15 to 16. Use `double` unless you have a specific reason not to — memory pressure, or hardware that only does 32-bit maths.

:::key
Floats trade exactness for range. A fixed number of fraction bits means a fixed number of significant digits, no matter how large or small the number is. Most values you write get rounded to the nearest representable one.
:::

## Why 0.1 cannot be stored exactly

Here's the crucial insight, and it isn't about computers at all.

In decimal, write 1/3 and you get 0.3333… forever. Ten digits, a hundred digits, a million — it never terminates. That's not a flaw in decimal; it's that 3 doesn't divide a power of 10. Cut it off anywhere and you have an approximation.

Binary has exactly the same problem, with different fractions. A binary fraction can only terminate if its denominator is a power of two. So 0.5, 0.25 and 0.75 are perfectly exact. But 0.1 is 1/10, and 10 isn't a power of two:

```text
  0.1 in binary  =  0.0001100110011001100110011001100...
                            └──── 1100 repeats forever ────┘
```

The hardware stores as much as fits and rounds the rest. In a 32-bit float, 0.1 becomes these bits:

```text
  0 01111011 10011001100110011001101      = 0x3DCCCCCD

  which is exactly  0.100000001490116119384765625
```

Not 0.1. The nearest thing to 0.1 the format can hold. A `double` gets much closer, but still not exact — and now `0.1 + 0.2` makes sense:

```c
printf("%.17g\n", 0.1 + 0.2);         /* 0.30000000000000004 */
printf("%d\n", 0.1 + 0.2 == 0.3);      /* 0 (false) */
printf("%.17g\n", 0.1 * 3.0);         /* 0.30000000000000004 */
```

Two slightly-wrong numbers were added, and the slightly-wrong sum rounded to a value that isn't the slightly-wrong stored form of 0.3. Everything behaved correctly. The result is still not what you wanted.

:::analogy
Ask someone to write down one third with a pen and two decimal places. They write 0.33. Ask for three of them added up and you get 0.99, not 1. Nobody made a mistake — the notation just couldn't hold the value.
:::

## Never compare floats with ==

The direct consequence: equality on floats is a trap.

```c
double x = 0.1 + 0.2;
if (x == 0.3) {            /* almost never true */
    puts("equal");
}
```

Instead, ask whether two values are *close enough*, within a small tolerance usually called **epsilon**:

```c
#include <math.h>

int nearly_equal(double a, double b, double eps) {
    return fabs(a - b) <= eps;
}

nearly_equal(0.1 + 0.2, 0.3, 1e-9);   /* 1 (true) */
```

A fixed tolerance like `1e-9` is fine when your numbers are around 1. It's useless when they're around a billion, where the gap between neighbouring doubles is already larger than that — and far too loose when they're around 1e-20. Serious comparisons scale the tolerance with the magnitude of the values. In C99, you can build that comparison from `<math.h>`:

```c
#include <math.h>

int nearly_equal_scaled(double a, double b,
                        double relative_tolerance,
                        double absolute_tolerance) {
    double difference = fabs(a - b);
    double scale = fmax(fabs(a), fabs(b));
    double tolerance = fmax(absolute_tolerance,
                            relative_tolerance * scale);

    return difference <= tolerance;
}

nearly_equal_scaled(0.1 + 0.2, 0.3, 1e-9, 1e-12);  /* 1 (true) */
```

:::warning
The same trap catches loop conditions. `for (double t = 0.0; t != 1.0; t += 0.1)` may never terminate, because the accumulated value steps straight past 1.0 without landing on it. Count with an integer and derive the float: `for (int i = 0; i <= 10; i++) { double t = i / 10.0; ... }`.
:::

## Errors that grow

Rounding once is harmless. The danger is rounding repeatedly, or in ways that amplify.

**Accumulated error.** Every operation rounds a little. Do it thousands of times and the little errors add up:

```c
double total = 0.0;

for (int i = 0; i < 10; ++i) {
    total += 0.1;
}
printf("%.17g\n", total);       /* 0.99999999999999989 */
printf("%d\n", total == 1.0);   /* 0 (false) */
```

**Absorption.** When one value is vastly larger than the other, the small one can vanish entirely — there simply aren't enough fraction bits to record it:

```c
double a = 1e16;
printf("%.1f\n", (a + 1.0) - a);  /* 0.0 — the 1 disappeared */
```

**Catastrophic cancellation.** This is the nastiest. Subtract two nearly equal numbers and the leading digits — the ones you were confident about — cancel out, promoting whatever rounding noise was hiding in the low bits to the front of the result:

```c
double x = 1.0000001;
double y = 1.0000000;
printf("%.17g\n", x - y);  /* 1.0000000005838672e-07, not exactly 1e-07 */
```

Both inputs were accurate to about sixteen digits. The difference is accurate to far fewer, because most of those digits agreed and destroyed each other. When a formula subtracts near-equal quantities, that's a signal to look for an algebraically equivalent form that doesn't.

:::tip
Adding many numbers of similar size is more accurate than adding them in a random order — summing from smallest to largest keeps the running total close to the values being added. If accuracy really matters, there are compensated summation algorithms that track the lost bits and add them back.
:::

## NaN, infinity and the other special values

IEEE 754 reserves some bit patterns for values that aren't ordinary numbers, so that awkward operations return something instead of crashing.

**Infinity.** Dividing a nonzero float by zero gives `+inf` or `-inf` rather than a crash. Infinities behave sensibly under further arithmetic: `inf + 1` is `inf`.

**NaN** — Not a Number. The result of an operation with no meaningful answer: `0.0/0.0`, `sqrt(-1.0)`, `inf - inf`.

NaN has one famously strange property: **it isn't equal to anything, including itself.**

```c
#include <math.h>

const double n = NAN;
printf("%d\n", n == n);  /* 0 (false) */
printf("%d\n", n != n);  /* 1 (true) */
```

That looks broken, but it's deliberate — NaN means "no valid value here", and two unknowns shouldn't be declared equal. It also gives you the classic portable test: `x != x` is true only for NaN. C99 provides the clearer `isnan` macro in `<math.h>`, and you should use it.

:::warning
NaN is contagious: almost any arithmetic involving a NaN produces a NaN. One bad value early in a pipeline can quietly turn an entire result into NaN, and because comparisons with NaN are all false, an `if (value > threshold)` guard silently lets it through. Check for it at the point where it can first appear.
:::

## When not to use floating point

Sometimes the right move is to not use floats at all.

**Money is the classic case.** Currency is exact and decimal by nature; floats are approximate and binary. Store amounts as **integers of the smallest unit** — cents, or pennies — and format for display at the very end.

```c
long price_cents  = 1999;   /* £19.99 */
long tax_cents    = price_cents * 8 / 100;
long total_cents  = price_cents + tax_cents;
/* exact integer arithmetic; divide by 100 only when printing */
```

For values with a fixed number of decimal places, use **scaled integers** so the calculation never enters binary floating point:

```c
long a_tenths = 1;  /* 0.1 */
long b_tenths = 2;  /* 0.2 */
long c_tenths = 3;  /* 0.3 */

printf("%d\n", a_tenths + b_tenths == c_tenths);  /* 1 (true) */
```

Standard C99 has no built-in decimal type. When fixed-point integers are not flexible enough, use a purpose-built decimal library and construct values from text or scaled integers — converting from a `double` would faithfully preserve the binary float's error.

Floats remain exactly right for what they were designed for: physics, graphics, audio, statistics, machine learning — measured quantities with inherent uncertainty, where huge range matters and the last digit doesn't.

## Check Your Understanding

:::quiz
Q: Why can't 0.1 be stored exactly in a binary floating point number?
- Because floats only store integers
- Because 1/10 is a repeating fraction in binary, just as 1/3 repeats in decimal *
- Because 0.1 is too small for a double
- Because of a rounding bug in most CPUs
E: A binary fraction terminates only when its denominator is a power of two. Ten isn't, so 0.1 repeats forever and must be rounded to fit.
:::

:::predict
```c
#include <math.h>
#include <stdio.h>

int main(void) {
    double n = NAN;
    printf("%d\n", n == n);
    return 0;
}
```
- 1 (true)
- 0 (false) *
- NaN
- The program fails
E: NaN compares unequal to everything, including itself, so the comparison produces C's integer false value, 0. That property is what makes `x != x` a portable NaN test.
:::

:::quiz
Q: You're building a shopping cart. How should you store prices?
- As `float`, rounding to two decimals when displaying
- As `double`, which has enough precision to be safe
- As integers counting the smallest currency unit, or as a decimal type *
- As strings, converting to float for each calculation
E: Money is exact and decimal; binary floats are neither. Integer cents or a purpose-built decimal type keeps every total exactly right.
:::

## Talk about it

> Floating point makes a deliberate trade: enormous range in exchange for approximate values. Describe a program you'd happily write with floats and one where you'd refuse, and explain what makes the difference. What question would you ask yourself to decide?

## What's next

You've now seen how numbers and text are laid out in bits. One question remains about anything larger than a single byte: in what *order* do those bytes sit in memory? In **Endianness** you'll see the two competing answers, why they still both exist, and the very specific situations where getting it wrong turns your data into gibberish.

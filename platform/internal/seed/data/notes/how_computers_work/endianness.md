# Endianness

A single byte has no ambiguity — it goes in one memory slot and that's that. But a 32-bit integer needs four slots, and now somebody has to decide which byte goes first. Two reasonable answers exist, both are in use, and the disagreement between them has broken more file parsers and network protocols than almost anything else. This lesson makes sure it never surprises you.

## The question a single byte never raises

Memory is addressed one byte at a time. Store the number `0x12345678` — four bytes: `12`, `34`, `56`, `78` — starting at address `0x100`, and the hardware must place them in some order.

**Byte order**, or **endianness**, is that choice. And notice it only ever applies to values *wider than one byte*. A `char`, a byte of ASCII, a UTF-8 code unit — none of them care. That's part of why UTF-8 travels so painlessly between machines while UTF-16 needs a marker to say which way round it is.

:::key
Endianness is about the order of *bytes* within a multi-byte value. It never changes the order of bits inside a byte, and it never affects single-byte data.
:::

## Big-endian and little-endian

**Big-endian** stores the **most significant byte first**, at the lowest address. It reads the way we write numbers on paper: the "big end" comes first.

**Little-endian** stores the **least significant byte first**. The little end leads.

```text
value: 0x12345678   (most significant byte 0x12, least significant 0x78)

        address:  0x100  0x101  0x102  0x103
                 ┌──────┬──────┬──────┬──────┐
  big-endian     │  12  │  34  │  56  │  78  │
                 └──────┴──────┴──────┴──────┘

                 ┌──────┬──────┬──────┬──────┐
  little-endian  │  78  │  56  │  34  │  12  │
                 └──────┴──────┴──────┴──────┘
```

Both hold the number `0x12345678`. Both are correct. They just disagree about which end of the number lands at the low address.

The names come from *Gulliver's Travels*, where two factions go to war over which end of a boiled egg to crack. The joke is well earned: neither order is meaningfully better, and the industry split anyway.

:::analogy
Think of writing a date. `2026-09-02` puts the most significant part first; `02/09/2026` puts the least significant first. Both encode the same day. Everything works fine until two systems exchange dates without agreeing which convention they're using — and then it fails silently, which is the worst way to fail.
:::

## Which machines use which

**Little-endian** dominates today. x86-64 processors are little-endian, and ARM chips — which are technically able to do either — run little-endian in nearly all common configurations. So your laptop and your phone are almost certainly little-endian.

**Big-endian** hasn't disappeared. Several older and specialised architectures use it, and more importantly it is baked permanently into the internet: the headers of TCP, UDP and IP all specify big-endian, which is why big-endian is also called **network byte order**.

File formats picked sides too. PNG is defined as big-endian; BMP and ZIP are little-endian. Whoever wrote the specification chose, and every reader and writer has to follow.

## Where it actually bites you

For most day-to-day programming, endianness is invisible. The compiler and CPU agree, so `int x = 5;` works and you never think about it. It becomes visible in exactly four situations — and they're all situations where bytes cross a boundary.

**Binary file formats.** You write a 4-byte length field on a little-endian machine and someone reads it on a big-endian one. They get `0x78563412` instead of `0x12345678` — a plausible-looking, completely wrong number.

**Network protocols.** Two machines exchanging raw binary data must agree on order. This is exactly why network byte order was standardised.

**Casting pointers between types.** Reinterpreting the same memory as a different width exposes the layout:

```c
#include <stdio.h>

int main(void) {
    unsigned int   v = 1;
    unsigned short *s = (unsigned short *)&v;
    printf("%u\n", s[0]);
    /* 1 on a little-endian machine, 0 on a big-endian one */
    return 0;
}
```

**Raw memory dumps.** Print bytes directly and the order is right there:

```c
#include <stdio.h>

int main(void) {
    unsigned int value = 0x12345678;
    unsigned char *bytes = (unsigned char *)&value;

    for (int i = 0; i < 4; i++)
        printf("%02X ", bytes[i]);
    putchar('\n');
    /* little-endian output:  78 56 34 12 */
    return 0;
}
```

The first time you see `78 56 34 12` in a debugger and briefly believe your data is corrupt is a rite of passage.

:::warning
Endianness bugs are quiet. Nothing crashes; you just get a wrong number that often still looks like a number. A file length of `0x78563412` instead of `0x12345678` reads as "about two billion bytes" rather than "about three hundred million" — plausible enough to pass a casual glance and cause chaos later.
:::

## Detecting endianness in code

You can ask the machine directly. Store a `1` and look at which byte it landed in:

```c
#include <stdio.h>

int main(void) {
    unsigned int x = 1;                        /* 0x00000001 */
    unsigned char *first = (unsigned char *)&x;

    if (*first == 1)
        puts("little-endian");   /* the 01 byte is at the low address */
    else
        puts("big-endian");      /* the low address holds 00 */
    return 0;
}
```

Inspecting bytes through an `unsigned char *` is the sanctioned way to examine any object's representation in C — that's exactly what it's for. Casting to some *other* type instead, as in the earlier `unsigned short` example, is useful for demonstrating layout but is not something to rely on in real code.

:::tip
Writing a detection function is a great exercise, but for production code you usually don't need one. If you always convert to a defined byte order at every I/O boundary, your program never has to know which kind of machine it's running on.
:::

## Converting: the htons family

Since network byte order is fixed, C gives you four small functions to convert between it and whatever your host uses. They live in `<arpa/inet.h>` on Unix-like systems, and their names read as an abbreviation:

```text
  htons   host to network, short  (16 bits)
  htonl   host to network, long   (32 bits)
  ntohs   network to host, short  (16 bits)
  ntohl   network to host, long   (32 bits)
```

```c
#include <stdio.h>
#include <arpa/inet.h>

int main(void) {
    unsigned short port = 8080;
    unsigned short wire = htons(port);   /* ready to send */

    unsigned short back = ntohs(wire);   /* recovered on the other side */
    printf("%u\n", back);                /* 8080 */
    return 0;
}
```

The elegant part: on a big-endian machine these functions do nothing at all, because host order already *is* network order. You call them unconditionally and the right thing happens everywhere. That's the pattern to internalise.

**Convert at the boundary, not in the middle.** The moment a multi-byte value enters your program from a file or a socket, convert it to host order. The moment it leaves, convert it back. In between, work in host order and never think about it again. Programs that scatter byte-swapping through their logic are the ones that end up double-swapping some values and not others.

Formats that avoid the problem entirely are worth noticing too. Text protocols like JSON sidestep byte order because they're sequences of single-byte UTF-8 units. UTF-16 files often begin with a **byte order mark** — the bytes `FE FF` for big-endian or `FF FE` for little-endian — so a reader can tell which it is. Explicitness beats assumption every time.

## Check Your Understanding

:::quiz
Q: A little-endian machine stores the 32-bit value `0x0000ABCD`. What are the four bytes at increasing addresses?
- 00 00 AB CD
- CD AB 00 00 *
- AB CD 00 00
- 00 00 CD AB
E: Little-endian puts the least significant byte first. The bytes of the value from most to least significant are 00 00 AB CD, so reversed they are CD AB 00 00.
:::

:::fill
Q: Complete the function that converts a 16-bit port number from host order to network order before sending it.
`unsigned short wire = ___(port);`
- htons *
- ntohs
- htonl
- ntohl
E: "Host TO Network, Short" — `htons` is the 16-bit conversion in the sending direction.
:::

:::quiz
Q: Which of these is *not* affected by endianness?
- A 32-bit integer written to a binary file
- A UTF-8 encoded string *
- A 16-bit port number sent over a socket
- A `double` reinterpreted as raw bytes
E: UTF-8 is a sequence of single-byte units, and byte order only matters for values wider than one byte. That is one of the reasons UTF-8 travels so cleanly between systems.
:::

## Talk about it

> Endianness bugs are notorious for being silent — the program keeps running and quietly produces wrong values. Describe a strategy you'd use to catch this class of bug early when writing code that reads a binary format. What could you put in place so a wrong assumption fails loudly instead of quietly?

## What's next

You've spent this lesson moving bytes around. In the final lesson of the course, **Bitwise Operators**, you'll get the tools to work *inside* a byte: AND, OR, XOR, NOT and the shifts. They're how you pack flags into a single integer, pull fields out of a packed value, and do a surprising amount of real work with almost no instructions.

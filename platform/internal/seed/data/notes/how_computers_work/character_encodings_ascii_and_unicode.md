# Character Encodings: ASCII and Unicode

You've printed strings since your very first program, and the machine has quietly cooperated. But there is no such thing as text in a computer — only bytes, plus an agreement about what those bytes mean. When two programs disagree about that agreement, you get the garbled characters everyone has seen and nobody enjoys. This lesson explains the agreement, so you'll know exactly what went wrong when you next meet `Ã©` where an `é` should be.

## Bytes are not text

A file on disk is a sequence of bytes. A string in memory is a sequence of bytes. Nothing in those bytes says "I am English text" or "I am a JPEG."

An **encoding** is the agreement that maps characters to bytes and back. Write with one encoding, read with a different one, and you get nonsense — the bytes were fine, the interpretation wasn't.

:::key
There is no "plain text." Every string of bytes is text *in some encoding*. If you don't know which encoding, you don't know what the bytes say — you're only guessing.
:::

## ASCII, and the mess that followed

The first widely adopted agreement was **ASCII**. It uses 7 bits, giving 128 possible values, and assigns them like this:

```text
   0 –  31   control codes (newline is 10, tab is 9)
      32     space
  48 –  57   the digits '0' to '9'
  65 –  90   'A' to 'Z'
  97 – 122   'a' to 'z'
     127     delete
```

So `'A'` is 65, `'a'` is 97, and `'0'` is 48. Those aren't arbitrary trivia — they're why the classic C trick `c - '0'` converts a digit character to its numeric value, and why uppercase and lowercase are exactly 32 apart.

ASCII was enough for American English and almost nothing else. No `é`, no `ß`, no `ñ`, no Greek, no Cyrillic, no Chinese, and no way to fit them: 128 slots was the whole budget.

Bytes have 8 bits, though, so there were 128 unused values above 127. Everyone filled them in — differently. Western Europe used one **code page** where byte 233 meant `é`. Eastern Europe used another where the same byte meant something else. Others existed for Greek, Cyrillic, Hebrew, and more. A document was only readable if you happened to know which code page it was written in, and nothing in the file told you.

Worse, none of these schemes could hold Chinese, Japanese or Korean, whose writing systems need thousands of characters. Those languages got their own separate multi-byte schemes, mutually incompatible with everything else. A single document couldn't mix Greek and Japanese at all.

:::warning
This isn't ancient history. Old files, old databases and old protocols still carry these encodings, and "the text looked fine on my machine" is one of the most common bug reports in software that handles user data.
:::

## Unicode: every character gets a number

**Unicode** solved the problem by separating two things that had always been tangled together.

First, give every character in every writing system its own unique number, called a **code point**. Written `U+` followed by hexadecimal:

```text
  U+0041   A          Latin capital A
  U+00E9   é          e with acute accent
  U+03A9   Ω          Greek capital omega
  U+20AC   €          euro sign
  U+4E2D   中          CJK ideograph "middle"
  U+1F600  😀         grinning face
```

Unicode is deliberately compatible at the bottom: the first 128 code points are exactly ASCII, so `U+0041` is 65, the same as always.

Second — and this is the part people miss — **Unicode does not say how to store those numbers as bytes.** It's a catalogue, not a file format. Code points go up to U+10FFFF, which needs 21 bits, and there's more than one sensible way to squeeze 21-bit numbers into bytes.

:::analogy
Unicode is the phone book: every person has a unique number. An encoding is how you write that number down — with dashes, with a country code, in words. Same person, different notation, and you have to know which notation you're reading.
:::

## UTF-8: turning numbers into bytes

**UTF-8** is the encoding that won, and deservedly. It is **variable width**: a code point takes one to four bytes depending on how big it is.

```text
  U+0000 – U+007F     1 byte    0xxxxxxx
  U+0080 – U+07FF     2 bytes   110xxxxx 10xxxxxx
  U+0800 – U+FFFF     3 bytes   1110xxxx 10xxxxxx 10xxxxxx
  U+10000 – U+10FFFF  4 bytes   11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```

Worked out for real characters:

```text
  'A'  U+0041   →  41
  'é'  U+00E9   →  C3 A9
  '€'  U+20AC   →  E2 82 AC
  '😀' U+1F600  →  F0 9F 98 80
```

Three properties make UTF-8 excellent. It is **ASCII-compatible**: any pure-ASCII file is already valid UTF-8, byte for byte, so decades of existing text and tools kept working. It is **self-synchronising**: the leading byte announces the length and every continuation byte starts `10`, so you can always find the start of a character even from the middle of a file. And it wastes no space on English while still reaching every writing system.

**UTF-16** uses two bytes for most common characters and four for the rest. **UTF-32** always uses four, which makes indexing simple but doubles or quadruples the size of ordinary text. Both also have a byte-order problem that UTF-8 doesn't — a topic for the next lesson.

## Why len() surprises you

Now the fun part. Ask "how long is this string?" and there are at least three defensible answers.

```python
s = "café"
print(len(s))                    # 4   — code points (Python 3)
print(len(s.encode("utf-8")))    # 5   — bytes: é takes two
```

C works at the byte level, so `strlen` counts bytes, not characters:

```c
#include <stdio.h>
#include <string.h>

int main(void) {
    const char *s = "café";        /* source file saved as UTF-8 */
    printf("%zu\n", strlen(s));    /* 5, not 4 */
    return 0;
}
```

It gets stranger. Some characters can be written more than one way. `é` might be the single code point U+00E9, or it might be a plain `e` followed by a **combining acute accent** (U+0301). They look identical on screen and compare as different strings:

```python
a = "é"            # U+00E9
b = "é"      # 'e' + combining accent
print(a == b)      # False
print(len(a), len(b))   # 1 2
```

And emoji push it further. What a reader calls "one character" — the thing that moves the cursor one step — is a **grapheme cluster**, which may be several code points glued together. A family emoji is built from several people joined by invisible zero-width joiners; a flag is two regional-indicator letters.

```python
print(len("👨‍👩‍👧"))   # 5  — three people plus two joiners
print(len("🇬🇧"))       # 2  — two regional indicator symbols
```

:::tip
Before slicing or truncating a string, ask which unit you actually mean: bytes (for storage and network limits), code points (for most programming), or grapheme clusters (for anything a human will read). Truncating UTF-8 at a fixed byte count is a reliable way to cut a character in half and produce invalid text.
:::

## Mojibake, and how to avoid it

**Mojibake** — from Japanese, roughly "character transformation" — is what you get when bytes written in one encoding are read as another. The classic:

```python
b = "café".encode("utf-8")     # b'caf\xc3\xa9'
print(b.decode("latin-1"))     # café
```

The two bytes `C3 A9` are one `é` in UTF-8. Read as single-byte Latin-1, they're `Ã` and `©`. Nothing is corrupted; the bytes are exactly what was written. Only the agreement broke.

Three rules keep you out of this:

**Use UTF-8 everywhere.** Source files, databases, config files, APIs, filenames. It is the sane default and the web's standard.

**Be explicit at I/O boundaries.** Every time bytes become text or text becomes bytes — opening a file, reading a socket, writing a response — state the encoding rather than accepting a platform default that may differ between machines.

```python
with open("notes.txt", encoding="utf-8") as f:
    text = f.read()
```

**Decode at the edges, work in text in the middle.** Convert incoming bytes to strings as early as you can, do all your logic on strings, and encode once on the way out. Programs that pass half-decoded bytes around are where these bugs breed.

## Check Your Understanding

:::quiz
Q: What is the difference between a Unicode code point and an encoding like UTF-8?
- They are two names for the same thing
- A code point is the number assigned to a character; an encoding decides how that number is stored as bytes *
- A code point is always one byte; an encoding is always four
- Code points apply to files, encodings apply to memory
E: Unicode assigns every character a number such as U+00E9. UTF-8, UTF-16 and UTF-32 are different ways of writing those numbers down as bytes.
:::

:::predict
```python
s = "héllo"
print(len(s), len(s.encode("utf-8")))
```
- 5 6 *
- 5 5
- 6 6
- 6 5
E: The string is five code points, but `é` needs two bytes in UTF-8, so the encoded form is six bytes long.
:::

:::quiz
Q: You open a file and see `café` where `café` should be. What most likely happened?
- The file is corrupted and the data is lost
- UTF-8 bytes were decoded as a single-byte encoding such as Latin-1 *
- The file was saved as UTF-32
- The font is missing those characters
E: This is classic mojibake: the `é` was stored as the two UTF-8 bytes C3 A9, and something read each byte as its own Latin-1 character. The bytes are intact — only the interpretation was wrong.
:::

## Talk about it

> Unicode's key move was separating *which character this is* from *how it's stored*. Explain why that separation was so powerful, and describe another situation in software where keeping "what something is" apart from "how it's represented" would make a system easier to change.

## What's next

You've just seen that multi-byte values raise a question single bytes never do: which byte goes first? For text, UTF-8 sidesteps it. For numbers, it's unavoidable. In **Endianness** you'll meet the two answers, see exactly how a four-byte integer sits in memory under each, and learn where the choice will actually bite you.

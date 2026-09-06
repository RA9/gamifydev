# String Search Algorithms

Find-in-page. Grep. A database `LIKE` query. A virus scanner matching a signature. All of them are doing the same thing: **substring search** — locating a small pattern inside a large text.

The obvious approach works, and it's slow in a specific, fixable way. This lesson covers the naive method and then three classic improvements, each of which attacks the waste from a different angle. The implementations get more involved as we go, so the goal here is that you can *explain* each one clearly and reach for the right one, not that you memorise every index.

## Naive search: try every position

Given a text `T` of length n and a pattern `P` of length m (with m ≤ n), find the position of the first occurrence of `P` in `T`, or report that there isn't one. Every algorithm in this lesson solves exactly that, differing only in how much redundant comparing it does.

```text
text:     A B A B C A B A B D
pattern:          C A B A B      -> match starts at index 4
```

The naive method lines the pattern up at position 0, compares characters until they disagree, slides right by one, and repeats.

```c
#include <stddef.h>
#include <stdio.h>
#include <string.h>

ptrdiff_t naive_search(const char *text, const char *pattern) {
    size_t n = strlen(text), m = strlen(pattern);
    if (m > n) return -1;
    for (size_t start = 0; start <= n - m; ++start) {
        size_t j = 0;
        while (j < m && text[start + j] == pattern[j]) ++j;
        if (j == m) return (ptrdiff_t)start;
    }
    return -1;
}

int main(void) {
    printf("%td\n", naive_search("ababcababd", "cabab")); /* 4 */
    printf("%td\n", naive_search("ababcababd", "xyz"));   /* -1 */
    return 0;
}
```

**Cost: O(n·m) worst case.** There are up to n-m+1 starting positions and each can cost up to m comparisons. In practice it's usually much faster — on ordinary English text the first character mismatches almost immediately, so the average is close to O(n). The killer input is one with lots of near-misses:

```text
text:     A A A A A A A A A B      start 0: match AAA, fail on B
pattern:  A A A B                  start 1: match AAA, fail on B
                                   start 2: ... same again, and again
```

Look at what's being wasted. At start 1, we re-compare the exact same characters we already examined at start 0. Every algorithm that follows is a different answer to the question: **how do we avoid re-comparing what we already know?**

:::key
Naive search is O(n·m) worst case and O(n) average on typical text. Its flaw is that after a mismatch it throws away everything it learned and restarts one position to the right.
:::

## Rabin-Karp: compare hashes, not characters

Rabin-Karp's idea: instead of comparing m characters at every position, compare a single **hash** of the window against the hash of the pattern — one integer comparison instead of m character comparisons. That would be pointless if hashing each window cost O(m). The trick is a **rolling hash**, which computes the next window's hash from the current one in O(1): as the window slides one character right, you subtract the outgoing character's contribution, shift, and add the incoming one.

```c
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

ptrdiff_t rabin_karp(const char *text, const char *pattern) {
    const uint64_t base = 256, mod = 1000003;
    size_t n = strlen(text), m = strlen(pattern);
    if (m == 0) return 0;
    if (m > n) return -1;

    uint64_t high = 1;
    for (size_t i = 1; i < m; ++i) high = (high * base) % mod;
    uint64_t p_hash = 0, t_hash = 0;
    for (size_t i = 0; i < m; ++i) {
        p_hash = (p_hash * base + (unsigned char)pattern[i]) % mod;
        t_hash = (t_hash * base + (unsigned char)text[i]) % mod;
    }

    for (size_t start = 0; start <= n - m; ++start) {
        if (p_hash == t_hash && memcmp(text + start, pattern, m) == 0) {
            return (ptrdiff_t)start;
        }
        if (start < n - m) {
            uint64_t outgoing = ((unsigned char)text[start] * high) % mod;
            t_hash = (t_hash + mod - outgoing) % mod;
            t_hash = (t_hash * base + (unsigned char)text[start + m]) % mod;
        }
    }
    return -1;
}
```

The verification step is essential: two different strings can hash to the same value — a **collision** — so a hash match is only a *candidate*. **Cost: O(n + m) average**, because collisions are rare with a good hash and a large modulus, so the expensive verification almost never runs; **O(n·m) worst case**, if extreme bad luck (or an adversary) makes every window collide.

Rabin-Karp's real strength is searching for **many patterns at once**: hash all of them into a set, and each window's single hash can be checked against every pattern in one operation. That's how plagiarism detectors and some intrusion-detection systems work.

:::analogy
Comparing hashes is like checking whether two long documents match by comparing their page counts first. A different page count means definitely different — no reading needed. The same page count means *maybe*, so you still have to read them. The cheap check filters out almost everything.
:::

## Knuth-Morris-Pratt: never re-compare a match

KMP attacks the waste directly. Its insight: when the pattern mismatches after matching k characters, you already know exactly what those k characters were — they're the first k characters of the pattern. So you can work out in advance how far to shift, without re-examining any text.

The key is the **prefix table** (also called the failure function). For each position in the pattern, it records the length of the longest **proper prefix** of the pattern that is also a **suffix** of the part matched so far.

```text
pattern:   A  B  A  B  C
index:     0  1  2  3  4
table:     0  0  1  2  0

table[3] = 2  because "ABAB" ends with "AB", which is also how it starts.
```

That number tells you how much of your progress survives a mismatch. If you matched "ABAB" and then failed, you don't restart — the "AB" at the end is already a valid start of a fresh attempt, so you resume comparing from pattern index 2.

```c
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

size_t *build_prefix_table(const char *pattern, size_t m) {
    size_t *table = calloc(m, sizeof *table);
    if (table == NULL && m != 0) return NULL;
    size_t length = 0;
    for (size_t i = 1; i < m; ++i) {
        while (length > 0 && pattern[i] != pattern[length]) {
            length = table[length - 1];
        }
        if (pattern[i] == pattern[length]) ++length;
        table[i] = length;
    }
    return table;                         /* Caller owns this allocation. */
}

int main(void) {
    const char *pattern = "ABABC";
    size_t m = strlen(pattern);
    size_t *table = build_prefix_table(pattern, m);
    if (table == NULL) return EXIT_FAILURE;
    for (size_t i = 0; i < m; ++i) printf("%zu%s", table[i], i + 1 == m ? "\n" : " ");
    free(table);
    return 0;
}
```

The search then walks the text with a single pointer that **never moves backwards**:

```c
#include <stddef.h>
#include <stdlib.h>
#include <string.h>

ptrdiff_t kmp_search(const char *text, const char *pattern) {
    size_t n = strlen(text), m = strlen(pattern);
    if (m == 0) return 0;
    size_t *table = build_prefix_table(pattern, m);
    if (table == NULL) return -1;

    size_t j = 0;
    ptrdiff_t result = -1;
    for (size_t i = 0; i < n; ++i) {       /* i never moves backwards. */
        while (j > 0 && text[i] != pattern[j]) j = table[j - 1];
        if (text[i] == pattern[j]) ++j;
        if (j == m) {
            result = (ptrdiff_t)(i - m + 1);
            break;
        }
    }
    free(table);
    return result;
}
```

**Cost: O(n + m) in the worst case** — O(m) to build the table, O(n) to scan. That's a genuine guarantee, not an average. Because `i` only ever increases and each character of the text is examined a bounded number of times, there is no bad input at all. Extra space is O(m) for the table.

:::key
KMP's guarantee comes from one fact: the text pointer never moves backwards. Everything the algorithm needs to recover from a mismatch is precomputed in the prefix table, so no text character is ever re-read.
:::

## Boyer-Moore: scan backwards and skip ahead

Boyer-Moore does something that feels wrong at first: it aligns the pattern and then compares **right to left**, from the pattern's last character back towards its first. The payoff is the **bad-character rule**. When you mismatch on some text character `c`, look at where `c` last appears in the pattern. Shift the pattern so that occurrence lines up with `c`. And if `c` doesn't appear in the pattern *at all*, you can skip the pattern's entire length in one jump — because no alignment overlapping that character could ever match.

```text
text:   A B C D E F G H I J K L E X A M P L E
index:  0 1 2 3 4 5 6 . . . . . 12  . . . . 18

step 1: E X A M P L E
                    ^  last pattern char vs text[6] = 'G'
                       'G' never occurs in "EXAMPLE"
                       -> shift the whole pattern past it: 7 places

step 2:                       E X A M P L E
                                          ^  vs text[13] = 'X'
                                             'X' is at pattern index 1
                                             -> shift 5 to line them up

step 3:                                 E X A M P L E
                                                    ^  'E' vs 'E', then
                                                       keep going left: match
```

In step 1, seven text characters were skipped after examining exactly one. That's what makes Boyer-Moore **sublinear in practice**: on a long text with a reasonably long pattern, it examines only a *fraction* of the characters. The bigger the alphabet and the longer the pattern, the better it does, because mismatched characters are more likely to be absent from the pattern entirely. Real implementations pair the bad-character rule with a second **good-suffix rule**, which uses the part that *did* match to justify an even bigger shift, and take whichever shift is larger.

**Cost:** O(n/m) best case, O(n + m) with the good-suffix rule in the worst case for a full implementation, and O(n·m) worst case for a simplified bad-character-only version. This is the algorithm at the heart of most `grep` implementations and many text editors' find function.

:::warning
Don't take "sublinear" to mean "always faster." Boyer-Moore's per-step bookkeeping is heavier than a plain comparison, so on short patterns or small alphabets — searching DNA, which has four letters — its big skips rarely materialise and a simpler algorithm can win.
:::

## Choosing, and what comes after

```text
Algorithm      Average      Worst        Extra space   Best at
------------------------------------------------------------------------
Naive          O(n)*        O(n*m)       O(1)          short text, quick code
Rabin-Karp     O(n + m)     O(n*m)       O(1)          many patterns at once
KMP            O(n + m)     O(n + m)     O(m)          guaranteed bound
Boyer-Moore    sublinear    O(n + m)**   O(m + a)      long patterns, big alphabet

*  on typical text;  ** with the good-suffix rule;  a = alphabet size
```

All four algorithms **preprocess the pattern**. If instead you're going to search the *same text* many times with different patterns, it pays to preprocess the text instead — that's what **suffix arrays** and **suffix trees** do, giving O(m) searches after an O(n) build, and a **trie** does the same for matching many patterns at once.

:::tip
In everyday C code, use a well-tested library routine when it fits: `strstr` searches null-terminated strings, while `memmem` is a common non-standard extension for byte ranges. Learn these algorithms to understand the trade-offs and to recognise when a specialised structure is what you actually need.
:::

## Check Your Understanding

:::quiz
Q: What is the worst-case time complexity of KMP?
- O(n · m)
- O(n + m) *
- O(n log m)
- O(m²)
E: KMP builds its prefix table in O(m) and then scans the text in O(n) with a pointer that never moves backwards, giving a true worst-case guarantee with no bad inputs.
:::

:::quiz
Q: Why must Rabin-Karp verify a match after the hashes agree?
- Because the rolling hash is slow
- Because two different strings can hash to the same value *
- Because the pattern might be empty
- Because the text might be unsorted
E: Equal hashes only mean the strings *might* match. A hash collision would otherwise produce a false positive, so a real character comparison confirms each candidate.
:::

## Talk about it

> KMP guarantees O(n + m) but never skips a single character, while Boyer-Moore has a worse worst case yet often reads only a fraction of the text. Explain in your own words how an algorithm can be *faster in practice* than one with a better worst-case bound, and what that means for how you choose between them.

## What's next

You've covered searching, sorting, graphs, greedy strategies and pattern matching. One lesson remains: **Caching Strategies**, about what to keep close at hand and what to throw away when you run out of room. It's an algorithmic problem hiding in every layer of every system you'll ever build, and it makes a fitting close to the course.

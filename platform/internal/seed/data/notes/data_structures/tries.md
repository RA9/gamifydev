# Tries

Every search structure so far has treated a key as one indivisible thing: hash the whole string, or compare the whole string. A **trie** (usually pronounced "try", from re*trie*val) does something different. It takes the key apart and stores it one character at a time, so shared beginnings are shared in the structure itself.

That single change makes a question possible that no hash table can answer: "give me every word that starts with 'ca'". If you've ever wondered how a search box suggests completions the instant you type, this is the lesson.

## Paths spell words

In a trie, each **edge** is labelled with a character, and each **node** represents the prefix formed by the path from the root down to it. To store a word, you walk down from the root spelling it out, creating nodes for characters that don't have an edge yet.

Here's a trie holding `cat`, `car`, `card` and `dog`. The `*` marks a node where a complete word ends.

```text
              (root)
              /    \
            c       d
            |       |
            a       o
           / \      |
          t*  r*    g*
              |
              d*
```

Read the paths: root → c → a → t is "cat". Root → c → a → r → d is "card". Notice that `cat`, `car` and `card` all share the `c`-`a` portion — stored once, not three times. That sharing is the trie's whole personality.

:::key
A trie stores keys along **paths**, not in nodes. The node itself often holds no key at all — just an array or table of outgoing character pointers and a flag saying "a word ends here".
:::

## The end-of-word flag

That flag isn't decoration; without it the trie would be wrong. Look at the node for `ca`: it's on the way to three real words, but "ca" itself is not one of them. And the node for `car` is *both* a complete word and on the way to `card`.

So every node carries a boolean. Reaching a node tells you the prefix exists. The flag tells you whether that prefix is also a stored word.

```c
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#define ALPHABET_SIZE 26

typedef struct TrieNode {
    struct TrieNode *children[ALPHABET_SIZE];
    bool is_word;
} TrieNode;

TrieNode *trie_node_create(void) {
    return calloc(1, sizeof(TrieNode));
}

bool insert(TrieNode *root, const char *word) {
    TrieNode *node = root;
    for (size_t i = 0; word[i] != '\0'; i++) {
        size_t index = (size_t)(word[i] - 'a');
        if (index >= ALPHABET_SIZE) return false;
        if (node->children[index] == NULL) {
            node->children[index] = trie_node_create();
            if (node->children[index] == NULL) return false;
        }
        node = node->children[index];
    }
    node->is_word = true;       // mark the last node
    return true;
}

void trie_free(TrieNode *node) {
    if (node == NULL) return;
    for (size_t i = 0; i < ALPHABET_SIZE; i++) trie_free(node->children[i]);
    free(node);
}

TrieNode *root = trie_node_create();
const char *words[] = {"cat", "car", "card", "dog"};
for (size_t i = 0; i < 4; i++) insert(root, words[i]);
printf("%d\n", (root->children['c' - 'a'] != NULL) +
                 (root->children['d' - 'a'] != NULL)); // 2
```

:::warning
Forgetting the end-of-word flag is the classic trie bug. Without it, a trie holding only "card" would happily report that "ca" and "car" are stored words, because those nodes exist along the path.
:::

## Lookup is O(m), not O(n)

Searching means walking the key's characters down from the root. If any character has no matching edge, the key isn't there.

```c
const TrieNode *walk(const TrieNode *root, const char *text) {
    const TrieNode *node = root;
    for (size_t i = 0; text[i] != '\0'; i++) {
        size_t index = (size_t)(text[i] - 'a');
        if (index >= ALPHABET_SIZE || node->children[index] == NULL) return NULL;
        node = node->children[index];
    }
    return node;
}

bool contains(const TrieNode *root, const char *word) {
    const TrieNode *node = walk(root, word);
    return node != NULL && node->is_word;
}

printf("%s\n", contains(root, "car") ? "true" : "false"); // true
printf("%s\n", contains(root, "ca") ? "true" : "false");  // false, prefix only
printf("%s\n", walk(root, "ca") != NULL ? "true" : "false"); // true
```

Count the work: one step per character. Looking up a 4-letter word takes 4 steps. Insert is the same walk with node creation, so also one step per character.

So lookup and insert are **O(m)**, where m is the *length of the key*. Read that carefully — there's no `n` in it. Whether the trie holds ten words or ten million, looking up `"card"` takes four steps. The number of stored keys does not appear in the cost at all.

:::analogy
A trie is the tabbed index of an encyclopedia, one tab per letter, then tabs within tabs. Finding "cardinal" means four flips: C, then A, then R, then D. Adding a thousand more volumes doesn't make those four flips any slower.
:::

## The superpower: prefix queries

Here's what tries are actually for. Every key beginning with a given prefix lives in the subtree hanging below that prefix's node. So "find everything starting with X" is: walk to X's node, then collect every word beneath it.

```c
#include <string.h>

void collect(const TrieNode *node, char path[], size_t length, size_t capacity) {
    if (node->is_word) printf("%s ", path);
    for (size_t i = 0; i < ALPHABET_SIZE && length + 1 < capacity; i++) {
        if (node->children[i] != NULL) {
            path[length] = (char)('a' + i);
            path[length + 1] = '\0';
            collect(node->children[i], path, length + 1, capacity);
        }
    }
}

void words_with_prefix(const TrieNode *root, const char *prefix) {
    const TrieNode *start = walk(root, prefix);
    if (start == NULL) return;
    char path[128];
    snprintf(path, sizeof path, "%s", prefix);
    collect(start, path, strlen(path), sizeof path);
}

words_with_prefix(root, "ca"); // car card cat
putchar('\n');
words_with_prefix(root, "z");  // prints no words
putchar('\n');
trie_free(root);
```

The walk to the prefix costs O(m), and then you visit exactly the nodes that lead to matching words — nothing else in the trie is touched. The total is O(m + k), where k is the size of what you collect. You never look at `dog`, or at the other nine million entries.

This is autocomplete. It's also spell-check candidate generation, IP routing tables (matching the longest matching address prefix), and predictive text on a phone keyboard.

:::example
Type "ca" into a search box backed by a trie of a million product names. The structure walks two edges, lands on one node, and reads off the products below it. The other 999,000-odd names are never examined.
:::

## Comparing honestly with a hash table

Tries look like a strict upgrade until you check the details. They aren't.

```text
                        hash table             trie
 ----------------------------------------------------------------
 exact lookup           O(m) average           O(m) worst case
 worst-case lookup      O(n)                   O(m)
 prefix query           O(n)  -- scan all      O(m + k)
 keys in sorted order   no                     yes (walk children
                                                in order)
 memory                 one entry per key      one node per distinct
                                                prefix + child arrays
 cache behaviour        good                   poorer (pointer chasing)
```

A few things worth being precise about.

**Exact lookup isn't actually faster in a trie.** A hash table must read the whole key to compute its hash, which is also O(m). But it does so in one tight pass over contiguous bytes, then makes one array access. A trie makes m separate hops through scattered nodes. In practice, for plain "is this key present?", the hash table usually wins.

**The trie's real edge is the worst case and the prefix query.** Its O(m) has no collision story attached, and prefix search is something a hash table simply cannot do without scanning everything.

**Memory is the trie's weak point.** Every distinct prefix needs a node, and every node needs a way to store its children. A small hash table per node is flexible but has real overhead; a fixed array of 26 (or 128) pointers per node is faster but wastes space on nodes with one child. Long keys with few shared prefixes are the bad case — you get a node per character and share almost nothing.

:::tip
A **compressed trie** (also called a radix tree) fixes the worst of that waste by collapsing chains of single-child nodes into one node holding a whole string. Storing "card" alone becomes one node labelled `card` instead of four nodes.
:::

Weighing all of that up gives a fairly clean rule of thumb.

:::key
Choose a **hash table** for plain key lookup. Choose a **trie** when you need prefix matching, sorted traversal, or a hard worst-case bound — and you can afford the memory.
:::

## Check Your Understanding

:::quiz
Q: What does the m stand for in a trie's O(m) lookup cost?
- The number of keys stored
- The length of the key being looked up *
- The size of the alphabet
- The height of the shortest word
E: A trie walks one node per character, so the cost depends on the key's length, not on how many keys are stored.
:::

:::predict
Q: A trie holds only the word "card". What does this print?
```c
TrieNode *root = trie_node_create();
insert(root, "card");
printf("%s %s\n", contains(root, "car") ? "true" : "false",
       walk(root, "car") != NULL ? "true" : "false");
trie_free(root);
```
- false true *
- true true
- false false
- true false
E: The node for "car" exists on the path to "card", so it is a valid prefix — but its `is_word` flag was never set, so it is not a stored word.
:::

:::fill
Q: Complete the field that marks a node as the end of a stored word.
`node->is_ ___ = true;`
- word *
- leaf
- root
E: The flag distinguishes "this prefix exists" from "this prefix is itself a complete key" — the node for "car" inside "card" needs exactly that distinction.
:::

## Talk about it

> A trie shares storage between keys with common beginnings, which is why it's brilliant for English words and terrible for random 32-character IDs. Describe a dataset from an app you know that would suit a trie well, and one that would waste enormous space in one — and say what it is about each dataset that decides it.

## What's next

Trees have taken us a long way, but they carry one restriction we've never questioned: every node has exactly one parent, and there are no loops. Drop both of those and you get the most general structure of all. In **Graphs and How to Represent Them**, we look at networks — roads, friendships, dependencies, web links — and at the two ways to store them, each with a very different cost profile.

# Why Data Structures Matter

You've written code that stores things — a list of names, a variable holding a score. This course is about the next question: *how* should that data be arranged? A **data structure** is a way of organising data in memory so that the operations you care about are fast. Choosing well can turn a program that crawls into one that answers instantly, without changing a single line of your logic.

By the end of this course you'll know the classic structures by name, know what each one is good at, and be able to pick the right one on purpose instead of by habit. Pixel will be along for the whole tour.

## A data structure is an arrangement, not a keyword

Beginners often meet data structures as vocabulary — "a list", "a dictionary" — and assume the difference is syntax. It isn't. The difference is *physical*: where the pieces of your data sit in memory and how they're connected to each other.

Consider storing three numbers. Here are two genuinely different arrangements:

```text
Arrangement A — side by side in one block (an array)

  addr:  100   104   108
        +-----+-----+-----+
        |  10 |  20 |  30 |
        +-----+-----+-----+

Arrangement B — scattered, each piece pointing to the next (a linked list)

  addr 340        addr 812        addr 108
  +----+----+     +----+----+     +----+------+
  | 10 | 812| --> | 20 | 108| --> | 30 | NULL |
  +----+----+     +----+----+     +----+------+
```

Same three numbers. Completely different machines. In A, the computer can jump straight to the third number by arithmetic. In B, it must start at the front and follow two arrows. That physical difference is what every lesson in this course is really about.

:::key
A data structure is a decision about **how data is laid out and linked in memory**. The syntax you type is just the surface; the layout is what determines speed.
:::

## The same data, wildly different costs

Here's the motivating example. Suppose you have a million usernames and you need to answer one question over and over: "is this name taken?"

Stored as a C array, checking membership means scanning:

```c
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

bool contains_linear(const char *names[], size_t count, const char *target) {
    for (size_t i = 0; i < count; i++) {
        if (strcmp(names[i], target) == 0) return true;
    }
    return false;
}

int main(void) {
    const char *names[] = {"ada", "grace", "linus", "margaret", "alan"};
    size_t count = sizeof names / sizeof names[0];
    printf("%s\n", contains_linear(names, count, "linus") ? "true" : "false");
    printf("%s\n", contains_linear(names, count, "pixel") ? "true" : "false");
    return 0;
}
```

The loop is simple, but look at what the machine does for `"pixel"`: it compares against `"ada"`, then `"grace"`, then `"linus"`, then `"margaret"`, then `"alan"` — every single item — before it can report `false`. With a million names, that's a million comparisons.

Now store the same data in a small open-addressed hash set (the structure you'll meet in the Hash Tables lesson):

```c
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#define SET_CAPACITY 16

typedef struct { const char *slots[SET_CAPACITY]; } StringSet;

size_t hash_string(const char *text) {
    size_t hash = 2166136261u;
    for (const unsigned char *ch = (const unsigned char *)text; *ch != '\0'; ch++) {
        hash = (hash ^ *ch) * 16777619u;
    }
    return hash;
}

void set_insert(StringSet *set, const char *value) {
    size_t index = hash_string(value) % SET_CAPACITY;
    while (set->slots[index] != NULL) index = (index + 1) % SET_CAPACITY;
    set->slots[index] = value;
}

bool set_contains(const StringSet *set, const char *value) {
    size_t index = hash_string(value) % SET_CAPACITY;
    while (set->slots[index] != NULL) {
        if (strcmp(set->slots[index], value) == 0) return true;
        index = (index + 1) % SET_CAPACITY;
    }
    return false;
}

int main(void) {
    StringSet names = {{NULL}};
    const char *values[] = {"ada", "grace", "linus", "margaret", "alan"};
    for (size_t i = 0; i < 5; i++) set_insert(&names, values[i]);
    printf("%s\n", set_contains(&names, "linus") ? "true" : "false");
    printf("%s\n", set_contains(&names, "pixel") ? "true" : "false");
    return 0;
}
```

The set needs more implementation code, but each lookup does not scan all names. It computes a number from the name and jumps more or less straight to the right spot. With a properly resized table, the average work per lookup stays roughly the same as the collection grows.

:::analogy
Searching a list is like looking for a word by reading a dictionary cover to cover. Searching a hash-based structure is like using the alphabetical thumb-tabs: you jump to the "P" section immediately. Same book, same words — a completely different amount of walking.
:::

## We compare structures by the cost of operations

To choose between structures, you need a fair way to compare them. The trick is to ignore the exact machine and the exact number of seconds — those change with every laptop — and instead ask: **as the amount of data grows, how does the work grow?**

Three operations come up again and again:

- **Lookup** (also called search or access) — find an item, or find out whether it's there.
- **Insert** — add a new item.
- **Delete** — remove an item.

We describe each cost with a shorthand called Big O. You'll study it properly in a later course; for now you only need three of its members:

```text
O(1)       constant   work doesn't grow with n at all
O(log n)   logarithmic work grows very slowly (double n -> +1 step)
O(n)       linear     work grows in step with n (double n -> double work)
```

Here `n` means "the number of items stored". So "lookup is O(n)" means: with twice as many items, expect roughly twice the work. "Lookup is O(1)" means: it costs about the same whether you have ten items or ten million.

:::example
Scanning a list of 1,000 names takes on the order of 1,000 comparisons; 1,000,000 names takes on the order of 1,000,000. That's O(n). A hash lookup takes about the same handful of steps in both cases. That's O(1).
:::

## There is no best structure, only trade-offs

It would be lovely if one structure won everything. None does. Every structure is fast at some operations because it accepted a cost somewhere else.

- Arrays give you instant access by position, but inserting into the middle means shifting everything after it.
- Linked lists splice new items in cheaply, but you can't jump to the middle — you have to walk there.
- Hash tables look things up fast, but they throw away any sense of order.
- Trees keep data sorted and searchable, but only stay fast if you keep them balanced.

So the real skill isn't memorising structures. It's asking, for *your* program: which operations happen constantly, and which happen rarely? Then pick the structure that's fast where it counts.

:::tip
Before choosing a structure, write down the one or two operations your code will perform most often. "I look up by ID thousands of times per second, and add items rarely" points at a completely different structure than "I always process the oldest item first."
:::

There's a matching danger in the other direction, though, and it's the more common mistake among people who have just learned all of this.

:::warning
Don't reach for a fancy structure before you know it's needed. A plain list is perfectly fine for twenty items, and simple code that's fast enough beats clever code that's hard to read. Reach for a specialised structure when the operation counts justify it.
:::

## Check Your Understanding

:::quiz
Q: What does it mean to say that an operation is O(n)?
- It always takes exactly n seconds
- The work grows in proportion to the number of items *
- It uses n bytes of memory
- It never takes more than one step
E: O(n) describes growth: double the items, roughly double the work. It says nothing about absolute seconds or memory.
:::

:::match
Q: Match each operation name to what it does.
- Lookup | Find an item or check whether it's present
- Insert | Add a new item to the structure
- Delete | Remove an item from the structure
E: Almost every structure in this course is judged by the cost of these three operations.
:::

## Talk about it

> Think of an app you use often — a music player, a messaging app, a game. Pick one thing it does that feels instant even though it must be searching through a lot of data. What data would it need to store, and which operation do you think it has optimised for? Describe your guess in your own words.

## What's next

You now have the frame for the whole course: data structures are arrangements in memory, and we judge them by the cost of insert, lookup and delete. Next up is the most fundamental arrangement of all — **Arrays**, where every element sits side by side in one unbroken block of memory. Once you understand why that layout gives instant access by index, every other structure becomes easier to reason about.

# Hash Tables

In the first lesson you saw a list scan through a million names to answer "is this taken?", while a set answered almost instantly. This is the lesson where that magic gets explained. A **hash table** turns a key into an array index by arithmetic, so instead of *searching* for where something lives, it *computes* where it lives.

Hash-table libraries for C, Java's `HashMap` and `HashSet`, and JavaScript's `Map` all implement this idea. It's arguably the most-used data structure in working software, so it's worth getting precisely right, including the part where it isn't fast.

## From key to bucket index

Underneath, a hash table is just an array. The slots are called **buckets**. The clever part is the **hash function**: a piece of code that takes a key of any shape and produces a number.

```c
#include <stddef.h>
#include <stdio.h>

size_t simple_hash(const char *key, size_t bucket_count) {
    size_t total = 0;
    for (const unsigned char *ch = (const unsigned char *)key; *ch != '\0'; ch++) {
        total += *ch;          // each char is converted to its numeric code
    }
    return total % bucket_count;
}

int main(void) {
    printf("%zu\n", simple_hash("cat", 8));   // 0
    printf("%zu\n", simple_hash("dog", 8));   // 2
    printf("%zu\n", simple_hash("bird", 8));  // 1
    return 0;
}
```

Two steps, always. First, turn the key into a number — here by adding up character codes. Second, squeeze that number into the valid range of bucket indices with `% num_buckets`. Now `"cat"` doesn't need to be *found*; it belongs in bucket 0, and it always will.

```text
 buckets (array of 8 slots)

  0: ("cat", 9)
  1: ("bird", 2)
  2: ("dog", 4)
  3:
  4:
  5:
  6:
  7:
```

To look up `"dog"` you hash it, get 2, and read bucket 2. One hash computation, one array access. No scanning.

:::analogy
A hash function is like a coat-check attendant who assigns your coat a hook based on your surname rather than the order you arrived. To get it back you don't search every hook — you say your name, they compute the hook, done.
:::

The hash function has to be fast, since it runs on every single operation, and it should scatter keys widely so they don't all pile into the same few buckets. Real ones are far more sophisticated than adding up character codes — but the two steps are always the same: make a number, then take it modulo the bucket count.

:::key
A hash table stores key-value pairs in an array, and uses a **hash function** to compute which slot a key belongs to. Lookup is arithmetic, not searching — that's the whole idea.
:::

## Collisions are inevitable

There are infinitely many possible keys and only a finite number of buckets, so sooner or later two different keys hash to the same index. That's a **collision**, and it isn't a bug — it's mathematically unavoidable.

Our toy hash function collides very easily, because addition doesn't care about order:

```c
printf("%zu\n", simple_hash("cat", 8));   // 0
printf("%zu\n", simple_hash("act", 8));   // 0, same letters and same sum
```

Every hash table needs a plan for this. There are two families of answers.

**Chaining.** Each bucket holds a small collection — usually a linked list — of all the pairs that landed there. On lookup, you hash to the bucket and then scan that short chain.

```text
  0: ("cat", 9) -> ("act", 1) -> NULL
  1:
  2: ("dog", 4) -> NULL
```

**Open addressing.** Every bucket holds at most one pair. On a collision, you *probe* for another slot by a fixed rule — the simplest being "try the next slot, and the next, wrapping around" (linear probing). Lookups follow the same probe sequence until they find the key or hit an empty slot.

```text
 a table holding only "cat" and "dog"; now insert ("act", 1)

  0: ("cat", 9)     <- occupied, probe onward
  1: ("act", 1)     <- landed here
  2: ("dog", 4)
```

Both work. Chaining is simpler to reason about and handles a full-ish table gracefully; open addressing keeps everything in one contiguous array, which is kinder to the cache. Real implementations vary: many compact C hash tables use open addressing, while Java's `HashMap` uses chaining and upgrades long chains into trees.

## Average O(1), worst case O(n) — say both

Here is the sentence to memorise, and to say out loud in full:

**Insert, lookup and delete in a hash table are O(1) on average, and O(n) in the worst case.**

The average case assumes keys spread out over the buckets reasonably evenly. Then each bucket holds a small, roughly constant number of items, and a lookup is a hash plus a glance at a tiny chain. Constant time.

The worst case is when *every* key collides into a single bucket. Then the table degenerates into one long list, and lookup means scanning all n items.

```text
 pathological case: every key hashes to bucket 3

  0:
  1:
  2:
  3: (k1) -> (k2) -> (k3) -> ... -> (kn)     lookup = O(n)
```

In practice this essentially never happens by accident with a decent hash function, which is why people say "hash tables are O(1)" in casual conversation. But it's a real worst case, and it can be triggered deliberately by an attacker feeding your server keys chosen to collide.

:::warning
Never write "hash lookup is O(1)" in an exam answer or a design doc without the word **average**. The worst case really is O(n), and pretending otherwise is the difference between understanding the structure and reciting it.
:::

## Load factor and resizing

How do you keep the average case actually average? Don't let the table get crowded. The **load factor** measures crowding:

```text
load factor = number of stored items / number of buckets

  4 items in 8 buckets  -> 0.5   comfortable
  7 items in 8 buckets  -> 0.875 crowded, collisions common
```

When the load factor crosses a threshold — commonly somewhere around 0.7 — the table **resizes**: allocate a bigger bucket array (typically double) and *rehash* every existing key into it. Rehashing is necessary because the bucket index depends on `% num_buckets`, and that divisor just changed.

A resize is O(n), which sounds alarming. But it's the same argument as growing a dynamic array: doubling means resizes happen exponentially less often, so the cost averages out. Insert stays **amortized O(1)**.

:::tip
If you know roughly how many items you'll store, creating the table with that capacity up front avoids several rehashes. It's a small win, but it's free — and it's why many library hash maps accept an initial-capacity argument.
:::

## Keys must be hashable

For any of this to work, a key must satisfy two rules.

**It must produce a hash value.** In C, the table's API usually accepts a hash function for the key type, or defines one internally.

**Its hash must never change while it's stored.** This is the important one. If you use an object as a key, and then mutate it so its hash changes, the table will compute a different bucket next time and simply fail to find an item that's definitely in there.

In C, the table cannot enforce immutability for you, so store stable key data or copy it into table-owned memory:

```c
#include <stdio.h>
#include <string.h>

typedef struct {
    char name[16];
    int birth_year;
} PersonKey;

size_t person_hash(const PersonKey *key, size_t bucket_count) {
    size_t hash = (size_t)key->birth_year;
    for (const unsigned char *ch = (const unsigned char *)key->name; *ch != '\0'; ch++) {
        hash = hash * 31u + *ch;
    }
    return hash % bucket_count;
}

int main(void) {
    PersonKey key = {"ada", 1815};
    PersonKey stored_key = key;  // copy the key into table-owned storage
    printf("bucket %zu\n", person_hash(&stored_key, 8));
    strcpy(key.name, "grace");  // changing the caller's copy cannot move stored_key
    printf("bucket %zu\n", person_hash(&stored_key, 8));
    return 0;
}
```

Integer values and copied strings or structs can make stable keys. A pointer to mutable caller-owned data is risky unless the table copies the pointed-to value. C leaves that contract to the implementation, so documenting ownership is essential.

:::key
Equal keys must have equal hashes, and a key's hash must not change while it's in the table. This is why hash-table keys are almost always **immutable** values.
:::

## The everyday version

You've been using all of this already:

```c
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#define CAPACITY 8

typedef enum { EMPTY, OCCUPIED, DELETED } SlotState;
typedef struct { char key[16]; int value; SlotState state; } Entry;
typedef struct { Entry entries[CAPACITY]; } HashTable;

size_t simple_hash(const char *key, size_t bucket_count);

size_t find_key(const HashTable *table, const char *key) {
    size_t index = simple_hash(key, CAPACITY);
    for (size_t probes = 0; probes < CAPACITY; probes++) {
        const Entry *entry = &table->entries[index];
        if (entry->state == EMPTY) break;
        if (entry->state == OCCUPIED && strcmp(entry->key, key) == 0) return index;
        index = (index + 1) % CAPACITY;
    }
    return CAPACITY;
}

bool set(HashTable *table, const char *key, int value) {
    size_t index = simple_hash(key, CAPACITY);
    size_t available = CAPACITY;
    for (size_t probes = 0; probes < CAPACITY; probes++) {
        Entry *entry = &table->entries[index];
        if (entry->state == OCCUPIED && strcmp(entry->key, key) == 0) {
            entry->value = value;
            return true;
        }
        if (entry->state == DELETED && available == CAPACITY) available = index;
        if (entry->state == EMPTY) {
            available = available == CAPACITY ? index : available;
            break;
        }
        index = (index + 1) % CAPACITY;
    }
    if (available == CAPACITY) return false;
    Entry *entry = &table->entries[available];
    snprintf(entry->key, sizeof entry->key, "%s", key);
    entry->value = value;
    entry->state = OCCUPIED;
    return true;
}

bool get(const HashTable *table, const char *key, int *value) {
    size_t index = find_key(table, key);
    if (index == CAPACITY) return false;
    *value = table->entries[index].value;
    return true;
}

void remove_key(HashTable *table, const char *key) {
    size_t index = find_key(table, key);
    if (index != CAPACITY) table->entries[index].state = DELETED;
}

int main(void) {
    HashTable inventory = {0};
    int count;
    set(&inventory, "potion", 3);
    set(&inventory, "sword", 1);
    set(&inventory, "shield", 2);       // insert, average O(1)
    if (get(&inventory, "potion", &count)) printf("%d\n", count);
    remove_key(&inventory, "sword");    // delete, average O(1)
    printf("%s\n", get(&inventory, "sword", &count) ? "true" : "false");
    return 0;
}
```

One thing a hash table cannot give you: order that means anything. Buckets are assigned by arithmetic, so there's no notion of "the smallest key" or "the next key after this one" without looking at everything. An implementation can maintain insertion order as separate metadata, but that is not sorted order and is not something hashing provides. When you need sorted or ranged queries, you want a tree — which is where we're heading.

## Check Your Understanding

:::quiz
Q: What is the worst-case time complexity of a lookup in a hash table?
- O(1)
- O(log n)
- O(n) *
- O(n²)
E: If every key collides into one bucket, the lookup degenerates into scanning all n items. Average case is O(1); worst case is O(n).
:::

:::predict
Q: What does this print?
```c
#include <stdio.h>

printf("%zu %zu\n", simple_hash("cat", 8), simple_hash("act", 8));
```
- 0 0 *
- 0 1
- 1 0
- 3 5
E: Adding character codes ignores order, so anagrams like "cat" and "act" produce the same sum and collide in the same bucket.
:::

:::fill
Q: Complete the term for the ratio of stored items to available buckets.
`___ factor = items / buckets`
- load *
- hash
- collision
E: The load factor measures how crowded the table is; crossing a threshold triggers a resize and rehash.
:::

## Talk about it

> Hash tables give up all sense of order to gain average O(1) lookup. Describe a feature you'd struggle to build on a hash table alone — something like "show me the ten highest scores" or "find the next appointment after 3pm" — and explain what the structure would have to do to answer it.

## What's next

Hash tables answer "is this exact key here?" brilliantly and "what comes next in order?" not at all. To get ordering back we need a shape that branches. Next up is **Trees and Terminology**, a vocabulary lesson that gives you the words — root, leaf, depth, height, subtree — that every remaining lesson in this course depends on. Pixel recommends reading it slowly; the payoff arrives in the three lessons after it.

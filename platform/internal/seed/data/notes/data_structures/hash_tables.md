# Hash Tables

In the first lesson you saw a list scan through a million names to answer "is this taken?", while a set answered almost instantly. This is the lesson where that magic gets explained. A **hash table** turns a key into an array index by arithmetic, so instead of *searching* for where something lives, it *computes* where it lives.

Python's `dict` and `set`, Java's `HashMap` and `HashSet`, JavaScript's `Map` — all hash tables. It's arguably the most-used data structure in working software, so it's worth getting precisely right, including the part where it isn't fast.

## From key to bucket index

Underneath, a hash table is just an array. The slots are called **buckets**. The clever part is the **hash function**: a piece of code that takes a key of any shape and produces a number.

```python
def simple_hash(key, num_buckets):
    total = 0
    for ch in key:
        total += ord(ch)      # ord() gives a character's numeric code
    return total % num_buckets

print(simple_hash("cat", 8))   # 0
print(simple_hash("dog", 8))   # 2
print(simple_hash("bird", 8))  # 1
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

```python
print(simple_hash("cat", 8))   # 0
print(simple_hash("act", 8))   # 0  -- same letters, same sum
```

Every hash table needs a plan for this. There are two families of answers.

**Chaining.** Each bucket holds a small collection — usually a linked list — of all the pairs that landed there. On lookup, you hash to the bucket and then scan that short chain.

```text
  0: ("cat", 9) -> ("act", 1) -> None
  1:
  2: ("dog", 4) -> None
```

**Open addressing.** Every bucket holds at most one pair. On a collision, you *probe* for another slot by a fixed rule — the simplest being "try the next slot, and the next, wrapping around" (linear probing). Lookups follow the same probe sequence until they find the key or hit an empty slot.

```text
 a table holding only "cat" and "dog"; now insert ("act", 1)

  0: ("cat", 9)     <- occupied, probe onward
  1: ("act", 1)     <- landed here
  2: ("dog", 4)
```

Both work. Chaining is simpler to reason about and handles a full-ish table gracefully; open addressing keeps everything in one contiguous array, which is kinder to the cache. Real implementations vary — Python's dict uses a form of open addressing, and Java's HashMap uses chaining (upgrading long chains into trees).

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

**It must produce a hash value.** In Python that means implementing `__hash__`; most built-in types do.

**Its hash must never change while it's stored.** This is the important one. If you use an object as a key, and then mutate it so its hash changes, the table will compute a different bucket next time and simply fail to find an item that's definitely in there.

Python enforces this by refusing to hash mutable types at all:

```python
scores = {}
scores[("ada", 1815)] = 10     # a tuple is immutable -> fine
print(scores[("ada", 1815)])   # 10

try:
    scores[["ada", 1815]] = 10          # a list is mutable
except TypeError:
    print("lists are unhashable")       # lists are unhashable
```

Strings, numbers, and tuples of immutable things are hashable. Lists, dicts and sets are not. That restriction isn't Python being fussy — it's protecting you from a bug that would be nearly impossible to find.

:::key
Equal keys must have equal hashes, and a key's hash must not change while it's in the table. This is why hash-table keys are almost always **immutable** values.
:::

## The everyday version

You've been using all of this already:

```python
inventory = {"potion": 3, "sword": 1}
inventory["shield"] = 2         # insert, average O(1)
print(inventory["potion"])      # 3, average O(1)
del inventory["sword"]          # delete, average O(1)
print("sword" in inventory)     # False, average O(1)
print(inventory)                # {'potion': 3, 'shield': 2}
```

One thing a hash table cannot give you: order that means anything. Buckets are assigned by arithmetic, so there's no notion of "the smallest key" or "the next key after this one" without looking at everything. (Python dicts do remember *insertion* order as a separate convenience, but that's not sorted order, and it isn't something hash tables give you in general.) When you need sorted or ranged queries, you want a tree — which is where we're heading.

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
```python
def simple_hash(key, num_buckets):
    total = 0
    for ch in key:
        total += ord(ch)
    return total % num_buckets

print(simple_hash("cat", 8), simple_hash("act", 8))
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

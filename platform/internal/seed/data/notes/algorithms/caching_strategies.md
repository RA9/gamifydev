# Caching Strategies

Every algorithm in this course assumed that getting at your data was free. It isn't. Reading from a CPU register, from RAM, from an SSD and from a server on another continent differ by factors so large that the *access* often costs more than the computation.

A **cache** is a small, fast store holding copies of things you expect to need again. Caching is where algorithms meet the real machine, and the interesting part is not storing things — it's deciding what to throw away when you run out of room.

## The memory hierarchy

Computer storage is a pyramid. Each level down is bigger and cheaper per byte, and dramatically slower.

```text
        CPU registers      bytes           fastest
        L1 / L2 / L3 cache KB to MB          |
        RAM                GB                |
        SSD / disk         TB                |
        network / cloud    unlimited       slowest
```

Every level in that pyramid is a cache for the one below it. Your CPU's L1 cache holds recently-used RAM; your operating system keeps recently-read disk blocks in RAM; your browser keeps downloaded images on disk; a web app keeps database results in memory. Same idea at every scale, same eviction question at every level.

It works because of a property of real programs called **locality**. *Temporal* locality: something you used recently, you'll probably use again soon. *Spatial* locality: if you used one item, you'll probably want its neighbours. Both hold for almost all real workloads, which is why a small cache can capture a large share of the accesses.

:::analogy
A cache is the desk you're working at. Everything you own is in the filing cabinet down the hall, but the handful of papers you need right now sit within arm's reach. The desk is tiny and the cabinet is enormous — the whole art is deciding which paper to put back when the desk fills up.
:::

## Hits, misses, and hit rate

Two outcomes when you look something up. A **hit** means it's in the cache — fast. A **miss** means it isn't, so you fetch it from the slow source and usually store it for next time. **Hit rate** is hits divided by total lookups, and it's the number that decides whether a cache is earning its keep.

The arithmetic is stark. Suppose a hit costs 1 unit of time and a miss costs 100:

```text
hit rate    average cost per lookup
---------------------------------------
   0%       0.00*1 + 1.00*100 = 100.0
  50%       0.50*1 + 0.50*100 =  50.5
  90%       0.90*1 + 0.10*100 =  10.9
  99%       0.99*1 + 0.01*100 =   2.0
```

Notice the shape: the jump from 90% to 99% cuts the average cost by more than five times. When misses are expensive, the last few percent of hit rate are worth more than everything before them — which is why eviction policy matters so much.

:::key
Cache performance is dominated by the **miss** cost, not the hit cost. That's why a small improvement in hit rate near the top of the range can produce a large improvement in real speed.
:::

## The eviction problem

A cache is small by definition — if it could hold everything, it would just be the storage. So when it's full and something new arrives, you must **evict** an existing entry. The theoretically perfect policy is Bélády's algorithm: evict whatever will be needed furthest in the future. It's provably optimal and completely unimplementable, because it requires knowing the future. Every real policy is a guess at it, based on the past.

Three common guesses:

- **LRU** — evict the **least recently used** item. Bets on temporal locality.
- **LFU** — evict the **least frequently used** item. Bets on long-run popularity.
- **MFU** — evict the **most frequently used** item. Sounds absurd; sometimes isn't.

## LRU: least recently used

LRU is the default choice, and usually the right one. The guess it makes — "what I haven't touched in a while, I probably won't touch soon" — matches how real programs behave.

The implementation challenge is doing it in O(1). You need two things at once: instant lookup by key, and instant identification of the oldest entry. The classic answer combines two structures you already know:

- a **hash map** from key to node, for O(1) lookup, and
- a **doubly linked list** in recency order, for O(1) move-to-front and O(1) removal from the back.

```text
hash map                doubly linked list (most recent -> least recent)

  "a" -> node ------->  [ a ] <-> [ c ] <-> [ b ] <-> [ d ]
  "b" -> node ---------------------^          ^         ^
  "c" -> node -------------^                             |
  "d" -> node --------------------------------------------

  get("b"):  map finds b's node in O(1)
             unlink it and splice it at the front in O(1)
  evict:     drop the node at the tail — "d"
```

The doubly linked part is essential: to unlink a node in O(1) you need a pointer to its *predecessor*, and only a doubly linked list gives you that.

Python's `OrderedDict` provides exactly these operations, so a working LRU cache is short:

```python
from collections import OrderedDict

class LRUCache:
    def __init__(self, capacity):
        self.capacity = capacity
        self.data = OrderedDict()
        self.hits = self.misses = 0

    def get(self, key):
        if key not in self.data:
            self.misses += 1
            return None
        self.hits += 1
        self.data.move_to_end(key)               # mark as most recently used
        return self.data[key]

    def put(self, key, value):
        if key in self.data:
            self.data.move_to_end(key)
        self.data[key] = value
        if len(self.data) > self.capacity:
            self.data.popitem(last=False)        # evict the oldest

cache = LRUCache(2)
cache.put("a", 1)
cache.put("b", 2)
print(cache.get("a"))     # 1   — this makes "a" the most recent
cache.put("c", 3)         # cache is full: evicts "b", not "a"
print(cache.get("b"))     # None — "b" was evicted
print(cache.hits, cache.misses)   # 1 1
```

**Cost: O(1) for both get and put**, with O(capacity) space.

:::warning
LRU has one classic failure: a **sequential scan** of data larger than the cache. Reading a million records once evicts everything useful and caches a million things you'll never look at again — a hit rate of zero and a cache full of junk. Real systems detect scans and bypass the cache for them.
:::

## LFU and MFU: counting accesses

LFU counts accesses and evicts the item with the lowest count. It bets on long-run popularity rather than recency, which suits workloads with a stable set of hot items — a CDN serving the same popular files for months, or a database caching a handful of heavily-read reference tables.

```python
from collections import Counter

class LFUCache:
    def __init__(self, capacity):
        self.capacity = capacity
        self.data = {}
        self.counts = Counter()

    def get(self, key):
        if key not in self.data:
            return None
        self.counts[key] += 1
        return self.data[key]

    def put(self, key, value):
        if key not in self.data and len(self.data) >= self.capacity:
            coldest = min(self.data, key=lambda k: self.counts[k])
            del self.data[coldest]
            del self.counts[coldest]
        self.data[key] = value
        self.counts[key] += 1

cache = LFUCache(2)
cache.put("x", 1)
cache.put("y", 2)
cache.get("x"); cache.get("x")     # x now has a high count
cache.put("z", 3)                  # evicts y, the least frequently used
print(cache.get("y"), cache.get("x"))   # None 1
```

That `min` scan makes eviction O(n); a production LFU keeps buckets of items grouped by count to get back to O(1).

LFU's real weakness is **cache pollution by history**. An item that was hugely popular last week has a big count and will sit there indefinitely, even though nobody wants it any more. LRU forgets naturally; LFU has to be made to forget, usually by ageing the counts downwards over time.

:::tip
Many production caches use a hybrid. One common design keeps a small LRU "admission window" in front of a frequency-based main cache, so a new item has to prove itself popular before it can evict something established. That gets LFU's stability without letting a scan wipe the cache.
:::

Now the third policy, **MFU** — evicting the *most* used item. It sounds like the worst idea in this lesson. There is a real case for it, though, and it comes from a specific access pattern: **once you've used something a lot, you're done with it.**

Consider a job that processes each record a fixed number of times and then never touches it again — a multi-pass batch conversion, a sequential scan through a table, a streaming pipeline. Under that pattern, a high access count is evidence that an item is *finished*, and the item with the lowest count is the one whose work still lies ahead. LRU and LFU both get this exactly backwards, keeping the completed items and evicting the ones about to be needed.

MFU is genuinely rare, but it makes the real lesson concrete: **no eviction policy is universally right.** Each encodes an assumption about future access, and it's only as good as that assumption is for your workload. Measure your hit rate before you defend your policy.

:::key
Every eviction policy is a prediction about the future. LRU predicts "recent means soon again", LFU predicts "popular stays popular", MFU predicts "heavily used means finished". Pick the one whose prediction matches your actual access pattern.
:::

## Invalidation and TTL, honestly

There's a well-worn joke that the two hard problems in computer science are naming things and cache invalidation. It's funny because it's true, and this is the part of caching nobody has fully solved.

The problem: a cache holds a *copy*. When the original changes, the copy is **stale** — silently wrong, and served confidently to your users.

**TTL (time to live).** Attach an expiry to each entry; after that, treat it as a miss and re-fetch. Simple, predictable, and always a compromise — a short TTL destroys your hit rate, a long one serves stale data for longer. You're literally choosing how wrong you're willing to be.

**Explicit invalidation.** When the underlying data changes, delete or update the cached copy. Correct in principle, and hard in practice, because you must find *every* cached thing that derived from the changed data. Miss one and it stays wrong forever. Distribute the cache across several machines and you now have to invalidate them all, in the presence of network failures and race conditions.

Neither is a solution; both are trade-offs. Be honest about which you've chosen and what it costs.

:::warning
The most dangerous cache bug isn't a low hit rate — it's a stale entry served as if it were fresh. A slow correct answer is recoverable; a fast wrong one may never be noticed. Before adding a cache, decide exactly how stale a value is allowed to be.
:::

## Check Your Understanding

:::quiz
Q: Which pair of data structures gives an LRU cache O(1) get and put?
- Two arrays
- A hash map and a doubly linked list *
- A binary search tree and a queue
- A heap and a hash set
E: The hash map gives O(1) lookup by key; the doubly linked list gives O(1) move-to-front and O(1) removal of the least recent entry from the tail.
:::

:::quiz
Q: A cache has a 90% hit rate. Hits cost 1 unit and misses cost 100. What is the average cost per lookup?
- 1.0
- 10.9 *
- 50.5
- 100.0
E: 0.90 × 1 + 0.10 × 100 = 0.9 + 10 = 10.9. Because misses dominate, raising the hit rate to 99% would cut this to about 2.0.
:::

## Talk about it

> Caching trades correctness risk for speed: a cached value may be out of date, and you've decided that's acceptable. Pick something you use — a social feed, a bank balance, a train timetable, a game leaderboard — and describe how stale that data could reasonably be before it becomes a real problem. What would you set the TTL to, and why?

## What's next

That's the Algorithms course complete, and it's a substantial thing to have finished. You can now recognise and reason about searching, four sorting algorithms with genuinely different trade-offs, tree and graph traversal, weighted shortest paths, heuristic search, minimum spanning trees, greedy strategies and where they break, backtracking with pruning, pattern matching, and caching. That's the toolkit professional engineers actually reach for.

Next comes **How Computers Work**, which goes the other direction — down, into the machine your algorithms have been running on. Bits and binary, the CPU and its instruction cycle, the memory hierarchy you just met from the caching side, and how a line of your code becomes something a processor executes. Pixel has been with you the whole way, and is genuinely impressed by how far you've come.

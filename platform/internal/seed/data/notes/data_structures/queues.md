# Queues

A **queue** is the polite twin of the stack. Where a stack serves the newest item first, a queue serves the *oldest* first — exactly like a line at a shop. First in, first out.

That sounds like a trivial change, and conceptually it is. But it creates a real implementation puzzle: a queue has two active ends instead of one, and the obvious array-based version turns out to be accidentally slow. Working out why, and how to fix it, teaches you more about arrays than another lesson on arrays would.

## First in, first out

The two core operations get their own names:

- **enqueue(x)** — add `x` to the back of the queue.
- **dequeue()** — remove and return the item at the front.

Some queues also offer **peek()** (look at the front without removing) and **is_empty()**. The ordering rule is **FIFO**: first in, first out.

```text
             front                        back
               |                            |
 dequeue <--  [ a ][ b ][ c ][ d ]  <-- enqueue

 dequeue() returns 'a'  ->  [ b ][ c ][ d ]
 enqueue('e')           ->  [ b ][ c ][ d ][ e ]
```

:::analogy
A queue is the line at a ticket window. New arrivals join the back; the person served is always the one who has been waiting longest. Nobody has to remember anything — the structure itself enforces fairness.
:::

Compare that with the previous lesson. A stack has one active end, so it only ever needs to know where the top is. A queue has two, and every implementation question in this lesson comes down to keeping track of both without doing unnecessary work.

:::key
Stack = LIFO, one active end. Queue = FIFO, two active ends: you add at the back and remove from the front. Everything else about queues follows from having two ends to manage.
:::

## Why the naive array queue is O(n)

Here's the tempting first attempt: use a list, `append` to enqueue, and `pop(0)` to dequeue.

```python
queue = []
queue.append("a")
queue.append("b")
queue.append("c")

print(queue.pop(0))   # a
print(queue)          # ['b', 'c']
```

It works. It's also quietly expensive. Remember from the Arrays lesson: an array must stay contiguous with no gaps, so removing the item at index 0 forces every remaining element to shift one slot left.

```text
 pop(0) on [a][b][c][d]

  [ a ][ b ][ c ][ d ]
        \    \    \      every element shifts left
  [ b ][ c ][ d ]
```

Enqueue is amortized O(1), but dequeue is **O(n)**. Process a million-item queue this way and you do on the order of a million shifts *per dequeue* early on — the whole loop becomes O(n²).

:::warning
`list.pop(0)` and `list.insert(0, x)` are the two most common accidental performance bugs in beginner Python. They look symmetric with `append` and `pop()`, but they are O(n), not O(1), because everything after index 0 has to move.
:::

## The fix: head and tail indices

The shifting is pure waste. Nothing about the *data* needs to move — we only moved it to keep the front at index 0. So stop insisting on that. Keep a `head` index that says where the front currently is, and just advance it.

```text
 [ a ][ b ][ c ][ d ]      head=0
   ^head

 dequeue -> a
 [ - ][ b ][ c ][ d ]      head=1   (nothing shifted)
        ^head
```

Now dequeue is O(1). But we've traded one problem for another: the space at the front is abandoned, and the queue crawls rightwards through the array until it falls off the end, wasting all the slots behind it.

## Circular buffers

The elegant fix is to let the queue **wrap around**. Treat the array as a ring: when an index runs past the last slot, it comes back to slot 0. The arithmetic for that is just the remainder operator, `%`.

```text
 capacity 4, wrapped around

     slot:   0     1     2     3
           +-----+-----+-----+-----+
           |  e  |  b  |  c  |  d  |
           +-----+-----+-----+-----+
              ^tail  ^head
           (next write)  (front of queue)

 order of items: b, c, d, e
```

This is called a **circular buffer** (or ring buffer). Here it is in full:

```python
class RingQueue:
    def __init__(self, capacity):
        self.data = [None] * capacity
        self.capacity = capacity
        self.head = 0     # index of the front item
        self.size = 0     # how many items are stored

    def enqueue(self, value):
        if self.size == self.capacity:
            raise IndexError("queue is full")
        tail = (self.head + self.size) % self.capacity
        self.data[tail] = value
        self.size += 1

    def dequeue(self):
        if self.size == 0:
            raise IndexError("queue is empty")
        value = self.data[self.head]
        self.head = (self.head + 1) % self.capacity
        self.size -= 1
        return value

q = RingQueue(4)
for ch in "abc":
    q.enqueue(ch)
print(q.dequeue())   # a
q.enqueue("d")
q.enqueue("e")       # wraps around into slot 0
print(q.dequeue())   # b
print(q.data)        # ['e', 'b', 'c', 'd']
```

Every operation touches exactly one slot and updates two integers. Both `enqueue` and `dequeue` are **O(1)**, with no shifting and no wasted space. Look at that final `q.data`: `e` really did land in slot 0, behind the front of the queue, and the structure handles it without blinking.

:::tip
The `% capacity` is doing all the work: it turns "one past the end" back into "the beginning". Whenever you see an index advanced with `(i + 1) % n`, someone is treating a flat array as a circle.
:::

## Deques: both ends, both directions

A **deque** (pronounced "deck", short for *double-ended queue*) lets you add and remove at *either* end, all in O(1). It's the generalisation of both a stack and a queue — restrict it to one end and it's a stack; push at one end and pop at the other and it's a queue.

Python ships one in the standard library, usually implemented as a linked structure of small blocks:

```python
from collections import deque

q = deque()
q.append("a")        # add to the right
q.append("b")
q.appendleft("z")    # add to the left
print(q)             # deque(['z', 'a', 'b'])
print(q.popleft())   # z
print(q.pop())       # b
print(list(q))       # ['a']
```

`append`, `appendleft`, `pop` and `popleft` are all O(1). If you need a queue in Python, reach for `deque` rather than a list — you get the correct complexity without writing the ring arithmetic yourself.

:::key
Use a **list** when you need index access. Use a **deque** when you need fast operations at the front. The one thing a deque gives up is O(1) access to the middle.
:::

## Where queues show up

**Task queues.** A web server receiving more requests than it can handle at once parks them in a queue and works through them in arrival order. Print spoolers, background job runners and message brokers are all queues at heart.

**Breadth-first search.** When you explore a maze, a network, or a tree level by level, you keep a queue of places to visit next. You'll meet this properly in the Algorithms course — and it's the reason so many graph algorithms start with `from collections import deque`.

**Buffering.** Audio, video and network data arrive in bursts but must be consumed at a steady rate. A circular buffer sits between producer and consumer, absorbing the difference. Fixed-capacity ring buffers are especially popular here, because they never allocate memory at an awkward moment.

:::example
A keyboard buffer is a small circular queue. Type faster than the program reads, and your keystrokes wait in order. Overflow it and the oldest or newest keystrokes get dropped — which is what a full ring buffer must decide to do.
:::

## Check Your Understanding

:::quiz
Q: Why is `queue.pop(0)` on a Python list an O(n) operation?
- Because Python lists are linked lists
- Because every remaining element must shift left to keep the array contiguous *
- Because it has to search for index 0
- Because it copies the list twice
E: A dynamic array has no gaps, so removing the front element forces all n-1 following elements to move one slot left.
:::

:::predict
Q: What does this print?
```python
from collections import deque
q = deque([1, 2, 3])
q.append(4)
q.appendleft(0)
print(q.popleft(), q.pop())
```
- 0 4 *
- 1 3
- 0 3
- 4 0
E: `appendleft(0)` puts 0 at the front and `append(4)` puts 4 at the back, so `popleft()` gives 0 and `pop()` gives 4.
:::

:::match
Q: Match each structure to its ordering rule.
- Stack | Last in, first out
- Queue | First in, first out
- Deque | Add and remove at either end
E: The ordering rule is the whole definition — the underlying array or linked list is just an implementation detail.
:::

## Talk about it

> A circular buffer with a fixed capacity eventually fills up, and then it must choose: refuse the new item, or throw away the oldest one to make room. Think of a live video stream and a bank's transaction queue. Which choice would you make for each, and what goes wrong if you pick the other one?

## What's next

Stacks and queues control the *order* things come out. They don't help you find a specific item — for that you'd still be scanning. Next comes the structure that solves the searching problem outright: **Hash Tables**, where a key is converted into an array index by arithmetic, giving average-case O(1) lookup. It's the machinery behind Python's `dict` and Java's `HashMap`, and it comes with a worst case you need to know about.

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

Here's the tempting first attempt: use an array, write at the end to enqueue, and shift left to dequeue.

```c
#include <stdio.h>
#include <string.h>

int main(void) {
    char queue[4] = {'a', 'b', 'c'};
    size_t size = 3;
    char front = queue[0];
    memmove(&queue[0], &queue[1], (size - 1) * sizeof queue[0]);
    size--;

    printf("%c\n", front);             // a
    printf("[%c, %c]\n", queue[0], queue[1]); // [b, c]
    return 0;
}
```

It works. It's also quietly expensive. Remember from the Arrays lesson: an array must stay contiguous with no gaps, so removing the item at index 0 forces every remaining element to shift one slot left.

```text
 remove index 0 from [a][b][c][d]

  [ a ][ b ][ c ][ d ]
        \    \    \      every element shifts left
  [ b ][ c ][ d ]
```

Enqueue is amortized O(1), but dequeue is **O(n)**. Process a million-item queue this way and you do on the order of a million shifts *per dequeue* early on — the whole loop becomes O(n²).

:::warning
Removing or inserting at index 0 by calling `memmove` is a common accidental performance bug in array-backed C code. It is O(n), not O(1), because everything after index 0 has to move.
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

```c
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

typedef struct {
    char *data;
    size_t capacity;
    size_t head;     // index of the front item
    size_t size;     // how many items are stored
} RingQueue;

bool queue_init(RingQueue *queue, size_t capacity) {
    queue->data = malloc(capacity * sizeof queue->data[0]);
    if (queue->data == NULL) return false;
    queue->capacity = capacity;
    queue->head = 0;
    queue->size = 0;
    return true;
}

bool enqueue(RingQueue *queue, char value) {
    if (queue->size == queue->capacity) return false;
    size_t tail = (queue->head + queue->size) % queue->capacity;
    queue->data[tail] = value;
    queue->size++;
    return true;
}

bool dequeue(RingQueue *queue, char *value) {
    if (queue->size == 0) return false;
    *value = queue->data[queue->head];
    queue->head = (queue->head + 1) % queue->capacity;
    queue->size--;
    return true;
}

int main(void) {
    RingQueue queue;
    if (!queue_init(&queue, 4)) return 1;
    enqueue(&queue, 'a'); enqueue(&queue, 'b'); enqueue(&queue, 'c');
    char value;
    dequeue(&queue, &value); printf("%c\n", value); // a
    enqueue(&queue, 'd');
    enqueue(&queue, 'e');                           // wraps into slot 0
    dequeue(&queue, &value); printf("%c\n", value); // b
    printf("[%c, %c, %c, %c]\n", queue.data[0], queue.data[1],
           queue.data[2], queue.data[3]);            // [e, b, c, d]
    free(queue.data);
    return 0;
}
```

Every operation touches exactly one slot and updates two integers. Both `enqueue` and `dequeue` are **O(1)**, with no shifting and no wasted space. Look at that final `q.data`: `e` really did land in slot 0, behind the front of the queue, and the structure handles it without blinking.

:::tip
The `% capacity` is doing all the work: it turns "one past the end" back into "the beginning". Whenever you see an index advanced with `(i + 1) % n`, someone is treating a flat array as a circle.
:::

## Deques: both ends, both directions

A **deque** (pronounced "deck", short for *double-ended queue*) lets you add and remove at *either* end, all in O(1). It's the generalisation of both a stack and a queue — restrict it to one end and it's a stack; push at one end and pop at the other and it's a queue.

A circular array can expose operations at both ends with constant-time index arithmetic:

```c
#include <stdbool.h>
#include <stdio.h>

#define DEQUE_CAPACITY 8

typedef struct { char data[DEQUE_CAPACITY]; size_t head; size_t size; } Deque;

bool push_back(Deque *deque, char value) {
    if (deque->size == DEQUE_CAPACITY) return false;
    deque->data[(deque->head + deque->size++) % DEQUE_CAPACITY] = value;
    return true;
}

bool push_front(Deque *deque, char value) {
    if (deque->size == DEQUE_CAPACITY) return false;
    deque->head = (deque->head + DEQUE_CAPACITY - 1) % DEQUE_CAPACITY;
    deque->data[deque->head] = value;
    deque->size++;
    return true;
}

bool pop_front(Deque *deque, char *value) {
    if (deque->size == 0) return false;
    *value = deque->data[deque->head];
    deque->head = (deque->head + 1) % DEQUE_CAPACITY;
    deque->size--;
    return true;
}

bool pop_back(Deque *deque, char *value) {
    if (deque->size == 0) return false;
    size_t index = (deque->head + deque->size - 1) % DEQUE_CAPACITY;
    *value = deque->data[index];
    deque->size--;
    return true;
}

int main(void) {
    Deque deque = {{0}, 0, 0};
    push_back(&deque, 'a'); push_back(&deque, 'b'); push_front(&deque, 'z');
    char value;
    pop_front(&deque, &value); printf("%c\n", value); // z
    pop_back(&deque, &value); printf("%c\n", value);  // b
    pop_front(&deque, &value); printf("%c\n", value); // a
    return 0;
}
```

`push_back`, `push_front`, `pop_back` and `pop_front` are all O(1). In C, you can implement the deque as a ring as above or use a library whose ownership and capacity rules fit your program.

:::key
Use an **array** when you need index access. Use a **deque** when you need fast operations at the front. The one thing a deque gives up is simple O(1) logical access to arbitrary middle positions.
:::

## Where queues show up

**Task queues.** A web server receiving more requests than it can handle at once parks them in a queue and works through them in arrival order. Print spoolers, background job runners and message brokers are all queues at heart.

**Breadth-first search.** When you explore a maze, a network, or a tree level by level, you keep a queue of places to visit next. You'll meet this properly in the Algorithms course — and it's why graph algorithms commonly maintain a ring queue or call a deque library.

**Buffering.** Audio, video and network data arrive in bursts but must be consumed at a steady rate. A circular buffer sits between producer and consumer, absorbing the difference. Fixed-capacity ring buffers are especially popular here, because they never allocate memory at an awkward moment.

:::example
A keyboard buffer is a small circular queue. Type faster than the program reads, and your keystrokes wait in order. Overflow it and the oldest or newest keystrokes get dropped — which is what a full ring buffer must decide to do.
:::

## Check Your Understanding

:::quiz
Q: Why is removing index 0 from an array-backed queue an O(n) operation?
- Because C arrays are linked lists
- Because every remaining element must shift left to keep the array contiguous *
- Because it has to search for index 0
- Because it copies the list twice
E: A dynamic array has no gaps, so removing the front element forces all n-1 following elements to move one slot left.
:::

:::predict
Q: What does this print?
```c
Deque deque = {{'1', '2', '3'}, 0, 3};
push_back(&deque, '4');
push_front(&deque, '0');
char front, back;
pop_front(&deque, &front);
pop_back(&deque, &back);
printf("%c %c\n", front, back);
```
- 0 4 *
- 1 3
- 0 3
- 4 0
E: `push_front` puts 0 at the front and `push_back` puts 4 at the back, so the two pop operations return 0 and 4.
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

Stacks and queues control the *order* things come out. They don't help you find a specific item — for that you'd still be scanning. Next comes the structure that solves the searching problem outright: **Hash Tables**, where a key is converted into an array index by arithmetic, giving average-case O(1) lookup. It's the machinery behind common C hash-table libraries and Java's `HashMap`, and it comes with a worst case you need to know about.

# Linked Lists

An array keeps its elements shoulder to shoulder in one block of memory. A **linked list** does the opposite: it lets each element live wherever it likes, and connects them with pointers. Each piece knows only one thing about the rest of the list — where to find the next piece.

That single change flips every cost in the array lesson on its head. Insertion becomes cheap, random access disappears, and a subtle performance cost appears that Big O won't tell you about. Let's build one and see.

## Nodes and links

The building block is a **node**: a small bundle holding a value and a reference to the next node. The list itself is nothing more than a variable pointing at the first node, traditionally called the **head**. The last node points at nothing — `None` in Python, `null` in Java, `NULL` in C — which is how you know you've reached the end.

```text
head
 |
 v
+----+----+     +----+----+     +----+------+
| 17 |  o-+---> |  3 |  o-+---> | 42 | None |
+----+----+     +----+----+     +----+------+
 addr 340        addr 812        addr 108
```

Look at the addresses: 340, 812, 108. They're in no particular order, and that's fine. The arrows, not the addresses, define the sequence.

```python
class Node:
    def __init__(self, value):
        self.value = value
        self.next = None

head = Node(17)
head.next = Node(3)
head.next.next = Node(42)

print(head.value)             # 17
print(head.next.value)        # 3
print(head.next.next.value)   # 42
```

:::analogy
An array is a row of numbered mailboxes on one wall. A linked list is a treasure hunt: each clue holds a prize and the address of the next clue. You can hide the clues anywhere in town — as long as each one points to the next, the trail still works.
:::

## Traversal is O(n), and there is no random access

To reach the third value in an array, you compute an address. To reach the third value in a linked list, you have no choice but to start at the head and follow arrows.

```python
def get(head, index):
    current = head
    for _ in range(index):
        if current is None:
            return None
        current = current.next
    return current.value if current else None

print(get(head, 2))   # 42
```

Reaching index `i` costs `i` hops, so reaching the last item of an n-item list costs n hops. Access by position is **O(n)**, and so is searching for a value. There is no formula that jumps to the middle, because the nodes aren't laid out in a pattern the computer can do arithmetic on.

:::key
A linked list has **no random access**. Every position is reached by walking from the head, so `get(i)` and `search(value)` are both O(n). Arrays win this one decisively.
:::

Walking the whole list to print it looks like this:

```python
def show(head):
    current = head
    while current is not None:
        print(current.value, end=" ")
        current = current.next
    print()

show(head)   # 17 3 42
```

## Insert and delete are O(1) — given a node reference

Here's the payoff. Splicing a new node into the middle of a linked list means rewiring two pointers. Nothing shifts. Nothing is copied. The rest of the list never learns anything happened.

```text
insert 50 after the node holding 17

 before:   [17| o-+---> [3 | o-+---> [42|None]

 step 1:   new node 50 points at 3
 step 2:   17 points at 50

 after:    [17| o-+---> [50| o-+---> [3 | o-+---> [42|None]
```

```python
def insert_after(node, value):
    new_node = Node(value)
    new_node.next = node.next   # 1. new node points where node pointed
    node.next = new_node        # 2. node now points at the new node
                                # order matters — swap these and you lose the tail

insert_after(head, 50)
show(head)   # 17 50 3 42
```

Two assignments, no matter how long the list is. That's **O(1)**.

:::warning
Read the phrase "O(1) insert" carefully: it's O(1) *once you already hold a reference to the node you're inserting after*. If you have to search for that node first, the search is O(n) and dominates. The pointer surgery is free; finding the right spot is not.
:::

Deletion works the same way — point the previous node past the doomed one:

```python
def delete_after(node):
    if node.next is not None:
        node.next = node.next.next   # skip over the next node

delete_after(head)
show(head)   # 17 3 42
```

## Singly vs doubly linked

The list we've built is **singly linked**: each node points forward only. That means you can never step backwards, and to delete a node you need a reference to the node *before* it — awkward, since the node itself can't tell you who points at it.

A **doubly linked list** fixes that by giving each node a `prev` pointer as well.

```text
None <-+ 17 +--> <--+ 3 +--> <--+ 42 +-> None
       +----+       +---+       +----+
```

```python
class DNode:
    def __init__(self, value):
        self.value = value
        self.next = None
        self.prev = None
```

With `prev` available you can walk in either direction, and you can delete a node given only that node — no hunting for its predecessor. The cost is one extra pointer of memory per node, and twice as many pointers to keep consistent when you edit. Most production list implementations (including Python's `collections.deque`) are doubly linked for exactly these reasons.

:::tip
When you write pointer-rewiring code, draw the before-and-after boxes on paper first and number the assignments. Nearly every linked-list bug is an assignment done in the wrong order, which drops a whole section of the list on the floor.
:::

## When a linked list actually wins

Comparing honestly:

```text
operation                    array (dynamic)   linked list
---------------------------------------------------------
access by index              O(1)              O(n)
search for a value           O(n)              O(n)
insert/delete at front       O(n)              O(1)
insert/delete at end         O(1) amortized    O(1) with a tail pointer
insert/delete given a node   O(n)  (shifting)  O(1)
extra memory per element     none              1-2 pointers
cache behaviour              excellent         poor
```

Linked lists are the right answer when you're constantly adding and removing at the ends or at a position you already hold, and you rarely need to jump to index `i`. That's precisely the shape of the next two lessons — stacks and queues — and it's why doubly linked lists sit underneath so many queue implementations.

But be fair about the downsides. A linked list spends extra memory on pointers, and its nodes are scattered across memory, so walking it defeats the processor's cache: each hop is likely a fresh trip to main memory. An array walk of a million items and a linked-list walk of a million items are both O(n), and the array walk is usually dramatically faster in practice.

:::key
Linked lists buy **cheap structural edits** with **expensive traversal and poor cache locality**. If your code mostly reads by index or scans in order, an array is very often the better choice even where a list looks more "elegant".
:::

## Check Your Understanding

:::quiz
Q: Why can't a linked list offer O(1) access by index the way an array can?
- Because linked lists are always sorted
- Because its nodes aren't in a predictable memory pattern, so the position must be walked to *
- Because pointers are slow to read
- Because linked lists can't store numbers
E: Array access is address arithmetic on a contiguous block. A linked list's nodes can be anywhere, so the only way to reach position i is to follow i pointers.
:::

:::predict
Q: What does this print?
```python
class Node:
    def __init__(self, value):
        self.value = value
        self.next = None

a = Node(1)
a.next = Node(2)
a.next.next = Node(3)
a.next = a.next.next
print(a.next.value)
```
- 3 *
- 2
- 1
- None
E: `a.next` is reassigned to the node holding 3, so the node holding 2 is skipped over and dropped from the list.
:::

## Talk about it

> A linked list makes insertion cheap and access expensive; an array does the reverse. Imagine you're storing the songs in a playlist that a user reorders constantly by dragging items around, versus the frames of a video that you always read straight through. Which structure fits each job, and why?

## What's next

You've now met the two foundational layouts — contiguous and linked — and everything ahead is built from them. Next we take a linked list or an array and *restrict* what you're allowed to do with it, which sounds like a downside but turns out to be enormously useful. In **Stacks**, you may only add and remove at one end, and that single rule quietly powers undo buttons, browser history, and the way your own function calls run.

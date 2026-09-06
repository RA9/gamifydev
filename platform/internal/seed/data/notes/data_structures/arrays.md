# Arrays

An **array** is the simplest data structure there is, and also the one everything else is built on top of. It stores its elements side by side in a single unbroken block of memory. That one decision — *contiguous* storage — is what gives arrays their famous superpower and their equally famous weakness.

In this lesson you'll see exactly how the computer finds `data[7]` without looking at the first seven items, how a dynamic array in C grows even though its current block has a fixed size, and why inserting into the middle is expensive.

## Contiguous memory

Picture memory as a long numbered street of byte-sized slots. When you create an array of five 4-byte integers, the computer reserves twenty bytes *in a row* and remembers the address of the first one.

```text
 array of 5 ints, each 4 bytes, starting at address 1000

 index:      0       1       2       3       4
          +-------+-------+-------+-------+-------+
 value:   |  17   |   3   |  42   |   8   |  99   |
          +-------+-------+-------+-------+-------+
 addr:     1000    1004    1008    1012    1016
```

Notice there are no arrows, no bookkeeping, no gaps. The elements are *implicitly* connected by being neighbours. The array itself is really just one number: the address where the block starts, called the **base address**.

:::analogy
An array is a row of numbered mailboxes bolted to a wall in order. You don't need directions to box 12 — you know where box 1 is, and every box is the same width, so you can walk straight to it.
:::

## The index formula gives O(1) access

Because every element is the same size and they sit in a row, the address of element `i` is pure arithmetic:

```text
address of element i  =  base_address + i * element_size
```

For our array above, element 3 lives at `1000 + 3 * 4 = 1012`. Check the diagram — that's exactly right.

This is the whole trick. Reading `data[3]` is one multiply, one add, and one memory fetch. Reading `data[999999]` is *also* one multiply, one add, and one fetch. The cost does not depend on how big the array is or where in it you're looking.

```c
#include <stdio.h>

int main(void) {
    int data[5] = {17, 3, 42, 8, 99};
    printf("%d\n", data[3]);   // 8
    return 0;
}
```

:::key
Array access by index is **O(1)** — constant time — because the computer *computes* the address with `base + i * element_size` instead of searching for it.
:::

That constant-time access by position is called **random access**: you can reach any element as cheaply as any other, in any order. Not every structure offers it, as you'll see in the next lesson.

## Fixed size vs dynamic arrays

A raw array in C has a size fixed when it's created. `int data[5]` reserves room for exactly five ints. If a sixth arrives, there's no guarantee the memory right after the block is free — a neighbouring variable may already live there. C does not perform bounds checks for you, so the program must keep every index below the array length.

```c
int scores[5] = {0};
scores[0] = 100;
/* scores[5] = 7; would be out of bounds: valid indexes are 0 through 4. */
```

In C, you can build a growable version by tracking both its length and capacity:

```c
#include <stdio.h>
#include <stdlib.h>

typedef struct {
    int *data;
    size_t length;
    size_t capacity;
} IntArray;

int append(IntArray *array, int value) {
    if (array->length == array->capacity) {
        size_t new_capacity = array->capacity == 0 ? 1 : array->capacity * 2;
        int *new_data = realloc(array->data, new_capacity * sizeof *new_data);
        if (new_data == NULL) {
            return 0;
        }
        array->data = new_data;
        array->capacity = new_capacity;
    }
    array->data[array->length++] = value;
    return 1;
}

int main(void) {
    IntArray scores = {NULL, 0, 0};
    if (!append(&scores, 100) || !append(&scores, 7) || !append(&scores, 42)) {
        free(scores.data);
        return 1;
    }
    printf("[%d, %d, %d]\n", scores.data[0], scores.data[1], scores.data[2]);
    printf("%zu\n", scores.length);
    free(scores.data);
    return 0;
}
```

This `IntArray` is a **dynamic array**: an array that grows. Under the hood it still owns one contiguous block, but the block is deliberately made bigger than needed, so there's spare room at the end for the next few appends.

:::warning
"Dynamic" doesn't mean the block magically stretches. Memory can't stretch — something else is already parked next door. Growing an array always means allocating a *new, larger* block somewhere else and copying the old contents into it.
:::

## Doubling and amortized O(1) append

If growing means copying everything, isn't `append` expensive? Sometimes, yes — but rarely enough that it doesn't matter, thanks to one clever choice: when the block fills up, don't add one slot. **Double** the capacity.

```text
capacity 1  [a]                     full -> allocate 2, copy 1 item
capacity 2  [a][b]                  full -> allocate 4, copy 2 items
capacity 4  [a][b][c][d]            full -> allocate 8, copy 4 items
capacity 8  [a][b][c][d][e][f][g][h]
```

Look at the pattern. Appending items 5, 6, 7 and 8 costs one step each — there's spare room. Only item 9 triggers a copy. As the array gets bigger, the expensive copies get *further and further apart*, because each doubling buys twice as much breathing room as the last.

Add up the total copying to reach n items and it comes to less than 2n operations overall. Spread across n appends, that's a constant amount of work each. We call this **amortized O(1)**: any single append might be O(n), but the average over many appends is O(1).

:::key
Append to a dynamic array is **amortized O(1)**. Most appends are instant; occasionally one triggers a doubling that costs O(n), and doubling makes those rare enough to average out.
:::

## Insert and delete in the middle are O(n)

Contiguity gives you free arithmetic, but it demands something in return: **no gaps**. Every element must sit exactly one slot after the one before it. So inserting anywhere but the end forces everything to the right to shuffle over.

```text
insert 50 at index 1 into [17, 3, 42, 8]

  before:  [ 17 |  3 | 42 |  8 |    ]
                  \    \    \        shift each one right
  after:   [ 17 | 50 |  3 | 42 |  8 ]
```

Deleting has the mirror problem: remove index 1 and everything after it must slide left to close the hole.

```c
#include <stdio.h>
#include <string.h>

int main(void) {
    int data[5] = {17, 3, 42, 8};
    size_t length = 4;

    memmove(&data[2], &data[1], (length - 1) * sizeof data[0]);
    data[1] = 50;             // shifts 3, 42, 8 one place right
    length++;
    printf("[%d, %d, %d, %d, %d]\n", data[0], data[1], data[2], data[3], data[4]);

    memmove(&data[0], &data[1], (length - 1) * sizeof data[0]);
    length--;                 // shifts everything left to close the gap
    printf("[%d, %d, %d, %d]\n", data[0], data[1], data[2], data[3]);
    return 0;
}
```

In the worst case — inserting or deleting at the front — all n elements move. That's **O(n)**. Deleting or appending at the *end* touches nothing else, so that stays O(1) (amortized for append).

:::tip
If you repeatedly insert or remove at index 0 by shifting an array inside a loop, that loop is quietly O(n²). Either append to the end and reverse once at the finish, or use a structure built for front operations — you'll meet deques in the Queues lesson.
:::

## Cache friendliness: the hidden bonus

There's a practical advantage that Big O notation doesn't capture. Modern processors don't fetch one value from memory at a time; they fetch a whole chunk — typically 64 bytes — into a small, very fast memory called the **cache**.

When you walk an array in order, the first fetch drags the next several elements in for free. By the time you ask for element 1, it's already sitting in the cache. Looping over an array is therefore much faster in practice than looping over the same number of values scattered around memory, even though both are O(n).

:::analogy
It's like a librarian who brings you a whole shelf instead of a single book. If the next book you want is on that shelf, you get it instantly. Arrays keep their next book on the same shelf, every time.
:::

## Check Your Understanding

:::fill
Q: Complete the formula the computer uses to find element `i` of an array.
`address = base_address + i * ___`
- element_size *
- array_length
- index
E: Every element is the same size, so multiplying the index by the element size gives the offset from the start of the block.
:::

:::predict
Q: What does this print?
```c
#include <stdio.h>
#include <string.h>

int main(void) {
    int data[5] = {1, 2, 3};
    size_t length = 3;
    memmove(&data[1], &data[0], length * sizeof data[0]);
    data[0] = 9;
    length++;
    data[length++] = 4;
    printf("[%d, %d, %d, %d, %d]\n", data[0], data[1], data[2], data[3], data[4]);
    return 0;
}
```
- [9, 1, 2, 3, 4] *
- [1, 2, 3, 9, 4]
- [9, 4, 1, 2, 3]
- [1, 2, 3, 4, 9]
E: `memmove` makes room at the front and 9 is written there; assigning at `data[length]` adds 4 to the end.
:::

:::quiz
Q: Why is appending to a dynamic array described as *amortized* O(1) rather than plain O(1)?
- Because appending is actually always O(n)
- Because occasional resizes cost O(n), but doubling makes them rare enough to average out to O(1) *
- Because C arrays always resize on every append
- Because the array must be sorted first
E: Most appends use spare capacity and cost O(1). A full block triggers an O(n) copy, but each doubling postpones the next one twice as long, so the average per append stays constant.
:::

## Talk about it

> Arrays trade cheap insertion for cheap access: you get instant reads by index, but you pay to shift elements whenever you insert or delete in the middle. Describe a program you could imagine writing where that trade is clearly a *good* deal — and one where it would clearly hurt.

## What's next

Arrays win on random access and cache behaviour, and lose on inserting anywhere but the end. That weakness is exactly what the next structure is designed to fix. In **Linked Lists**, we abandon contiguous memory entirely: each element carries a pointer to the next one, so splicing an item into the middle costs nothing at all. As always, that gift comes with a price tag — and we'll be honest about it.

# Stacks

A **stack** is a collection with one rule: you may only add and remove at the same end. The last thing you put in is the first thing you get back out. That's it — a whole data structure defined by a restriction.

Restrictions sound like a loss, but this one buys you something valuable: every operation becomes O(1), the code becomes impossible to misuse, and the "last in, first out" order turns out to match a surprising number of real problems — including the mechanism your own function calls run on.

## Last in, first out

The classic picture is a stack of plates. You add a plate to the top; you take a plate from the top. Reaching the plate at the bottom means removing every plate above it first. This ordering has a name: **LIFO**, for last in, first out.

```text
 push 'a'    push 'b'    push 'c'      pop -> 'c'    pop -> 'b'

                          +---+
              +---+       | c |  <- top
   +---+      | b |       | b |         +---+
   | a |      | a |       | a |         | b |          +---+
   +---+      +---+       +---+         | a |          | a |
                                        +---+          +---+
```

The operations have standard names:

- **push(x)** — put `x` on the top.
- **pop()** — remove and return the top item.
- **peek()** — look at the top item without removing it.
- **is_empty()** — is there anything in there?

:::analogy
A stack is the pile of unread mail on your desk. Each new letter lands on top, and when you finally sit down you read the newest one first. The letter from three weeks ago is at the bottom, waiting — which is either a feature or a problem, depending on the letter.
:::

## Every operation is O(1)

Notice what a stack never has to do: search, shift, or reorder. Both operations touch exactly one end, and that end is always known. So `push`, `pop`, `peek` and `is_empty` are all **O(1)** — constant time, regardless of how many items are stacked up.

In Python you get a stack for free by using a list and only ever touching its end:

```python
stack = []

stack.append("a")   # push
stack.append("b")
stack.append("c")

print(stack[-1])    # peek -> c
print(stack.pop())  # c
print(stack.pop())  # b
print(stack)        # ['a']
```

`append` and `pop()` with no argument both work at the end of the list, which — as you saw in the Arrays lesson — needs no shifting. That's why a dynamic array makes such a natural stack.

:::warning
`stack.pop()` on an empty list raises `IndexError`. A stack has no "top" when it's empty, and no data structure can invent one. Always check `if stack:` before popping, or wrap the pop in a guard that returns a sensible default.
:::

## Two ways to build one

You can implement a stack on top of either structure from the previous lessons, and both give O(1) operations.

**With an array (dynamic array).** The top of the stack is the end of the array. Push is an append; pop removes the last slot. Occasionally a push triggers a resize, so push is *amortized* O(1).

**With a linked list.** The top of the stack is the head. Push creates a node pointing at the old head; pop moves the head forward. No resizing ever happens, so every push is a true O(1) — at the cost of one pointer per item and worse cache behaviour.

Here's the linked-list version, written out so you can see there's no loop anywhere:

```python
class Stack:
    def __init__(self):
        self.head = None      # top of the stack

    def push(self, value):
        self.head = {"value": value, "next": self.head}

    def pop(self):
        if self.head is None:
            raise IndexError("pop from empty stack")
        node = self.head
        self.head = node["next"]
        return node["value"]

s = Stack()
s.push(1)
s.push(2)
print(s.pop())   # 2
print(s.pop())   # 1
```

:::key
A stack is an **abstract idea** — "LIFO, with push/pop/peek" — not a particular layout. An array and a linked list both implement it perfectly well. What you promise the caller is the behaviour, not the internals.
:::

## Where stacks show up

Once you know the shape, you start seeing stacks everywhere.

**Undo.** Every action you take gets pushed onto a stack. Ctrl+Z pops the most recent one and reverses it. LIFO is exactly right here: undo should reverse your *latest* action, not your first.

**Browser back button.** Each page you visit is pushed; Back pops the previous one off. (The forward button is a second stack — pages popped off the back stack get pushed onto it.)

**Balanced brackets.** Checking whether `([]{})` is properly nested is a textbook stack problem. Push every opening bracket; on a closing bracket, pop and check it matches.

```python
def balanced(text):
    pairs = {")": "(", "]": "[", "}": "{"}
    stack = []
    for ch in text:
        if ch in "([{":
            stack.append(ch)
        elif ch in pairs:
            if not stack or stack.pop() != pairs[ch]:
                return False
    return not stack

print(balanced("([]{})"))   # True
print(balanced("(]"))       # False
print(balanced("(("))       # False
```

Trace `"([]{})"`: push `(`, push `[`, then `]` pops `[` and matches, push `{`, then `}` pops `{` and matches, then `)` pops `(` and matches. The stack ends empty, so everything was closed in the right order.

**Expression evaluation.** Calculators and compilers use stacks to evaluate arithmetic while respecting precedence and parentheses — numbers on one stack, operators on another.

:::example
For `"(("` the loop pushes both open brackets and never pops. `return not stack` then reports `False`, because leftover items mean something was opened and never closed.
:::

## The call stack

Here's the one that ties back to something you've already used. When your program calls a function, the computer has to remember where to return to, plus that function's local variables. It stores all of that in a **stack frame**, and pushes the frame onto a stack — *the* stack, the one your program runs on.

```text
main() calls greet(), greet() calls shout()

    +-----------+
    |  shout()  |  <- currently running (top)
    +-----------+
    |  greet()  |  waiting
    +-----------+
    |   main()  |  waiting
    +-----------+
```

When `shout()` returns, its frame is popped and control resumes exactly where `greet()` left off. This is why functions return in the reverse order they were called — LIFO, enforced by hardware and the language runtime.

It also explains something from the recursion lesson. A recursive function that never reaches its base case keeps pushing frames without popping any. Eventually the stack runs out of room, and Python raises `RecursionError` — a stack overflow, the most famously named error in programming. It's not mysterious at all once you know the structure underneath it.

:::tip
When you're debugging and see a "stack trace" or "traceback", you're reading the call stack printed top to bottom. The innermost call is the one that failed; each line below it is the caller that's still waiting.
:::

## Check Your Understanding

:::quiz
Q: What does LIFO mean for a stack?
- The first item added is the first removed
- The last item added is the first removed *
- Items are removed in sorted order
- Items are removed at random
E: LIFO is "last in, first out" — push and pop both act on the same end, so the most recent item leaves first.
:::

:::predict
Q: What does this print?
```python
stack = []
for ch in "PIX":
    stack.append(ch)
out = ""
while stack:
    out += stack.pop()
print(out)
```
- XIP *
- PIX
- P
- X
E: The characters are pushed P, I, X and popped in reverse order — a stack reverses whatever you feed it.
:::

:::fill
Q: Complete the method name for the operation that looks at the top item without removing it.
`top = stack.___()`
- peek *
- push
- pop
E: `peek` inspects the top without changing the stack; `pop` would remove it.
:::

## Talk about it

> The call stack means a program only ever "remembers" the chain of calls that led to where it is now. Explain in your own words why a runaway recursive function causes a stack overflow rather than simply looping forever, and what that tells you about the memory a stack occupies.

## What's next

A stack serves the newest item first — perfect for undo and function calls, badly wrong for a queue at a coffee shop. Next up is the mirror image: **Queues**, where the first item in is the first item out. You'll see why the obvious array implementation is accidentally O(n), and how a neat trick called a circular buffer fixes it.

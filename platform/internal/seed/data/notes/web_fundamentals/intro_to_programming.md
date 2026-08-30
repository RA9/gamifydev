# Intro to Programming

Programming sounds intimidating, but at its core it's just giving a computer a clear, ordered list of instructions. This lesson teaches the ideas that show up in *every* programming language — once you have these, learning any specific language is mostly learning new spelling.

We'll stay language-agnostic but sprinkle in tiny JavaScript snippets to make the ideas concrete.

## What is programming?

**Programming is writing instructions that a computer can follow.** A computer is fast and precise but completely literal — it does *exactly* what you say, in the order you say it, with no common sense to fill in the gaps. Your job as a programmer is to break a goal into steps so small and so clear that even a literal machine can follow them perfectly.

:::analogy
Writing a program is like writing a recipe for someone who has never cooked and takes everything literally. "Add salt" isn't enough — they need "add 1 teaspoon of salt." Vague instructions produce surprising results.
:::

## How a computer runs instructions

By default, a computer runs your instructions **in sequence, top to bottom**, one line at a time. Each line finishes before the next begins.

```js
console.log("Step 1: wake up");
console.log("Step 2: brush teeth");
console.log("Step 3: eat breakfast");
```

That prints the three lines in order. Swap two lines and the output order changes — order matters.

## Variables: labelled boxes for data

A **variable** is a named container that stores a piece of information so you can use it later. The name is a label; the value inside can change.

:::analogy
A variable is a labelled box. The label is the name (`age`), and inside is the value (`25`). You can open the box to read what's inside, or replace it with something new.
:::

```js
let name = "Ada";   // a box labelled "name" holding the text "Ada"
let age = 25;       // a box labelled "age" holding the number 25
age = 26;           // replace what's in the box

console.log(name);  // Ada
console.log(age);   // 26
```

### Kinds of data

Most languages have a few basic **data types**:

- **Text** (called *strings*): `"hello"`, `"Frontend"` — always in quotes.
- **Numbers**: `7`, `3.14`, `-20`.
- **Booleans**: `true` or `false` — yes/no values.
- **Lists** (called *arrays*): `["red", "green", "blue"]` — many values in order.

```js
let title = "Profile Card";   // string
let score = 0;                // number
let isLoggedIn = false;       // boolean
let colors = ["red", "green", "blue"]; // array
```

:::fill
A value wrapped in quotes like `"hello"` is called a ____.
ANSWER: string
HINT: It's the data type used for text.
:::

## The universal building blocks

Almost every program is built from just four ideas. Learn these and you can read code in any language.

### 1. Sequence

Doing things in order, one after another — exactly what you saw above. This is the default.

### 2. Conditionals (making decisions)

A **conditional** lets the program choose a path based on whether something is true.

```js
let temperature = 30;

if (temperature > 25) {
  console.log("It's hot, wear shorts.");
} else {
  console.log("It's cool, bring a jacket.");
}
```

The computer checks the condition (`temperature > 25`). If it's `true`, it runs the first block; otherwise it runs the `else` block.

:::analogy
A conditional is a fork in the road. "If it's raining, take the umbrella; otherwise leave it home." The program picks a direction based on the current situation.
:::

### 3. Loops (repetition)

A **loop** repeats a block of instructions so you don't have to copy-paste them.

```js
for (let i = 1; i <= 3; i++) {
  console.log("Hello number " + i);
}
// Hello number 1
// Hello number 2
// Hello number 3
```

Instead of writing three `console.log` lines, the loop runs the same line three times, counting as it goes.

:::analogy
A loop is like telling someone "do 10 push-ups" instead of saying "do a push-up" ten separate times.
:::

### 4. Functions (reusable steps)

A **function** is a named, reusable block of instructions. You define it once, then *call* it whenever you need it — optionally passing in values (called *arguments*) and getting a result back.

```js
function greet(person) {
  return "Hi, " + person + "!";
}

console.log(greet("Ada"));   // Hi, Ada!
console.log(greet("Grace")); // Hi, Grace!
```

`greet` does one job. We wrote it once and reused it with different inputs.

:::analogy
A function is like a coffee machine: you put in a button press (the input), it runs its fixed steps inside, and out comes coffee (the output). You don't rebuild the machine each morning.
:::

:::quiz
Q: Which building block would you use to run the same code 100 times without copying it?
- A conditional
- A loop *
- A variable
E: Loops handle repetition. Conditionals make decisions; variables store data.
:::

## Syntax vs logic

Two different things can go wrong in code:

- **Syntax** is the *grammar* of the language — the exact punctuation and spelling it expects. A missing bracket or quote is a **syntax error**, and the program won't run at all.
- **Logic** is whether your steps actually solve the problem. Code can be perfectly spelled yet still do the wrong thing — that's a **logic error**.

```js
// Syntax error: missing closing quote
let name = "Ada;

// Logic error: runs fine, but the math is wrong
function double(n) {
  return n + 2;   // oops — should be n * 2
}
```

:::key
A program with a syntax error usually refuses to run. A program with a logic error runs happily but gives the wrong answer. Both are normal parts of coding.
:::

## Bugs and errors are normal

A **bug** is any mistake that makes a program behave incorrectly. Every programmer — beginner and expert — writes bugs constantly. The skill isn't avoiding them; it's *finding and fixing* them (called **debugging**). Error messages are not insults; they're hints pointing you toward the problem.

:::tip
When you hit an error, read the message slowly and look at the line number it mentions. Most of the time the fix is small: a missing bracket, a typo in a variable name, or a wrong comparison.
:::

## Thinking in pseudocode

Before writing real code, it helps to plan in **pseudocode** — plain language steps that capture the logic without worrying about exact syntax. Say you want to find the largest number in a list:

```bash
SET biggest to the first number
FOR each number in the list:
    IF the number is bigger than biggest:
        SET biggest to that number
SHOW biggest
```

Once the steps are clear, translating them into JavaScript is straightforward:

```js
let numbers = [4, 9, 2, 7];
let biggest = numbers[0];

for (let i = 1; i < numbers.length; i++) {
  if (numbers[i] > biggest) {
    biggest = numbers[i];
  }
}

console.log(biggest); // 9
```

## The problem-solving mindset

The hardest part of programming usually isn't the typing — it's breaking a big, fuzzy goal into small, exact steps. A reliable approach:

1. **Understand** the goal in one sentence.
2. **Break it down** into smaller sub-problems.
3. **Solve one piece** at a time, in pseudocode first.
4. **Translate** each piece into code.
5. **Test** it, find bugs, and fix them.

:::quiz
Q: What's the best first move when facing a big, intimidating coding task?
- Type as fast as possible and hope it works
- Break the problem into smaller, clearer steps *
- Memorize the whole language first
E: Decomposition — splitting a big problem into small steps — is the core problem-solving skill in programming. You can plan it in pseudocode before writing any real code.
:::

## Recap

- Programming is writing **clear, ordered instructions** for a literal machine.
- Code runs **top to bottom in sequence** by default.
- **Variables** are labelled boxes that store data; common types are strings, numbers, booleans, and arrays.
- The four universal building blocks: **sequence, conditionals, loops, functions**.
- **Syntax errors** stop a program; **logic errors** let it run but give wrong results — both are normal.
- **Bugs** are expected; debugging is a core skill, and error messages are helpful hints.
- Plan in **pseudocode** and **break problems into small steps** before writing real code.

**Next up:** Building a Website with HTML and CSS — you'll turn these ideas into a real, styled landing page.

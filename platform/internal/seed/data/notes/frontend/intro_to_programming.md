# Intro to Programming

Before you write a single line of HTML, it helps to understand what *programming* actually is — because once the core ideas click, every language you ever learn becomes a variation on the same handful of concepts.

By the end of this lesson you'll be able to explain, in your own words, what a program is and the building blocks every programmer uses every day.

## What is a program?

A program is a set of **instructions** that tells a computer exactly what to do, step by step. That's it. The skill of programming is really the skill of **breaking a big, fuzzy problem into small, precise steps** a computer can follow.

The catch: computers are extremely fast but extremely *literal*. They do exactly what you say — not what you *meant*.

:::analogy
Imagine giving directions to someone who takes everything literally. You say "make me a sandwich" and they stare blankly — they need every step: get bread, open the jar, spread the butter… Programming is writing instructions for the most literal helper imaginable. Vague won't cut it; precise will.
:::

## The building blocks

Nearly every programming language — JavaScript, Python, Java, C — shares the same core ideas. Learn them once and you'll recognise them everywhere.

### 1. Values and data types

A **value** is a piece of data. Values come in **types**:

- **Numbers** — `25`, `3.14`
- **Strings** (text) — `"Ada"`, `"hello world"`
- **Booleans** (yes/no) — `true` or `false`

### 2. Variables

A **variable** is a named box that stores a value so you can use it later. You put a value in, give the box a label, and refer to it by name.

![Variables are labeled boxes that store values](/images/lessons/variable-box.svg)

```js
let age = 25;
let name = "Ada";
let isStudent = true;
```

Now whenever you write `age`, the computer reads "25". Change the box's contents and everything that uses it updates — that's what makes variables so powerful.

:::quiz
Q: What does a variable do?
- Stores a value under a name so you can reuse it *
- Permanently draws something on the screen
- Makes the computer run faster
E: A variable is a labelled box: you store a value in it and refer back to it by name.
:::

### 3. Conditions — making decisions

A **condition** lets a program choose between paths: *if* something is true, do one thing; *otherwise*, do another.

```js
let age = 18;

if (age >= 18) {
  show("You can vote");
} else {
  show("Not yet");
}
```

Read it like English: "**if** age is 18 or more, show 'You can vote'; **otherwise**, show 'Not yet'." Conditions are how programs react to different situations.

### 4. Loops — repeating without copy-paste

A **loop** repeats an action many times without you writing it out over and over.

```js
for (let i = 1; i <= 3; i++) {
  show("Hello number " + i);
}
// Hello number 1
// Hello number 2
// Hello number 3
```

:::analogy
A loop is like telling someone "do 20 push-ups" instead of writing "do a push-up" twenty times. You describe the repetition once, and the computer handles the counting.
:::

### 5. Functions — reusable steps

A **function** is a named, reusable block of instructions. Define it once, then *call* it whenever you need it — possibly with different inputs.

```js
function greet(name) {
  show("Hi, " + name + "!");
}

greet("Ada");   // Hi, Ada!
greet("Grace"); // Hi, Grace!
```

Functions keep your code organised and save you from repeating yourself.

:::quiz
Q: You need to run the same 5 lines of code in three different places. What's the best tool?
- Copy and paste it three times
- Put it in a function and call it three times *
- Write it once and hope it runs everywhere
E: A function lets you write the logic once and reuse it by name — less code, fewer bugs, easier to change.
:::

## How this runs on the web

On the web, three languages divide the work — and each maps to what you've just learned:

1. **HTML** — *structure*: what content exists.
2. **CSS** — *presentation*: how it looks.
3. **JavaScript** — *behaviour*: the actual programming — variables, conditions, loops, and functions running in the browser.

:::key
Programming is breaking a problem into small, precise steps. Almost every language is built from the same five ideas: **values, variables, conditions, loops, and functions.** Master these and you can read and write code in any of them.
:::

## Talk about it

Try explaining out loud:

> "What is a program? And what's the difference between a variable, a condition, a loop, and a function?"

If you can give a one-sentence answer for each — ideally with your own analogy — you're ready to move on.

## What's next

Next up, **HTML Basics**, where you'll build the structure of a real web page. Then CSS makes it beautiful, and JavaScript brings these programming ideas to life in the browser.

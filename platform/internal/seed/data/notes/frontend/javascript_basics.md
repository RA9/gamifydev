# JavaScript Basics

HTML gives a page structure. CSS makes it look good. But both are *static* — they just sit there. **JavaScript** is what makes a page *do* things: respond to clicks, validate forms, update content without reloading. It's the **behaviour** layer of the web.

By the end of this lesson you'll understand how JavaScript reacts to what users do, and you'll be able to read a simple interactive script.

## What JavaScript adds

HTML and CSS describe what a page *is*. JavaScript describes what it *does* when someone interacts with it.

![A user event triggers JavaScript, which updates the page](/images/lessons/js-event-flow.svg)

The pattern above is the heart of front-end JavaScript: **an event happens → your code runs → the page changes.** Almost everything interactive you've ever used follows this loop.

:::analogy
If HTML is the skeleton and CSS is the skin, JavaScript is the **muscles and nervous system** — it senses what happens (a click, a keypress) and makes the body react.
:::

## Variables and types, quickly

JavaScript stores data in variables, just like the programming concepts you met earlier. You'll mostly use two keywords:

- `const` — a box whose value won't be reassigned (your default).
- `let` — a box whose value *can* change later.

```js
const name = "Ada";   // won't change
let score = 0;        // will change as the game goes on
score = score + 10;   // now 10
```

The everyday types are the same as any language: **numbers**, **strings** (text in quotes), and **booleans** (`true` / `false`).

:::predict
```js
let score = 0;
score = score + 10;
score = score + 5;
console.log(score);
```
- 15 *
- 105
- 0
- "105"
E: Each line adds to `score`: 0 → 10 → 15, so it logs the number 15.
:::

:::tip
Reach for `const` by default and only switch to `let` when you know the value needs to change. It makes your intent clear and prevents accidental reassignment bugs.
:::

## Functions: packaging behaviour

A **function** wraps up steps you can run on demand, optionally with inputs (called *parameters*).

```js
function add(a, b) {
  return a + b;
}

const total = add(3, 4); // 7
```

`return` hands a value back to whoever called the function. Functions are how you give a name to "the thing that happens when…".

:::quiz
Q: What does the `return` keyword do in a function?
- Prints text to the screen
- Sends a value back to wherever the function was called *
- Stops the whole program
E: `return` hands a value back to the caller. Here `add(3, 4)` returns `7`, which is stored in `total`.
:::

:::reorder
Arrange these lines to define a `greet` function and then call it.
- function greet(name) {
-   return "Hi, " + name;
- }
- greet("Ada");
E: Define the function first — header, body, then the closing brace — and call it afterwards.
:::

## Reacting to events

This is where JavaScript comes alive. You pick an element on the page, listen for an event, and run a function when it happens.

Say your HTML has a button and a message:

```html
<button id="cheer">Cheer me on</button>
<p id="msg">Waiting…</p>
```

JavaScript can make the button *do something* when clicked:

```js
const button = document.getElementById("cheer");
const msg = document.getElementById("msg");

button.addEventListener("click", function () {
  msg.textContent = "You've got this! 🎉";
});
```

Walk through it with the diagram in mind:

1. **Event** — the user clicks the button.
2. **JavaScript runs** — the function inside `addEventListener` executes.
3. **Page updates** — `msg.textContent` changes the paragraph's text on screen.

That trio — select an element, listen for an event, change the page — is the foundation of every interactive feature you'll build.

:::quiz
Q: In `button.addEventListener("click", ...)`, what is `"click"`?
- The text shown on the button
- The name of the event to listen for *
- A CSS class
E: `"click"` is the event type. The function you pass runs each time that event happens on the button.
:::

## A note on the DOM

When the browser loads your HTML, it builds a live, in-memory model of the page called the **DOM** (Document Object Model). `document.getElementById(...)` reaches into that model to grab an element, and changing it (like `textContent`) instantly updates what the user sees.

:::key
JavaScript is the **behaviour** layer. The core loop is **event → code → page update**. You select an element from the **DOM**, listen for an **event**, and run a **function** that changes the page.
:::

## Talk about it

Explain it in your own words:

> "What does JavaScript add that HTML and CSS can't? Describe what happens, step by step, when a user clicks a button."

If you can narrate the event → code → update loop, you understand the essence of front-end JavaScript.

## What's next

You've now met all three languages of the web: structure (HTML), style (CSS), and behaviour (JavaScript). Next you'll combine them to build **interactive websites** — and from there, real projects.

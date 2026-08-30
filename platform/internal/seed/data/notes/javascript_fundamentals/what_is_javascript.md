# What is JavaScript

Every web page you interact with — clicking buttons, submitting forms, watching animations — runs JavaScript under the hood. This lesson covers what JavaScript actually is, where it runs, and how to start writing it.

## The Three Layers of the Web

A web page is built from three technologies, each with a distinct job:

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Structure | HTML | Content and semantic meaning |
| Presentation | CSS | Visual styling and layout |
| Behavior | **JavaScript** | Interactivity and logic |

JavaScript is the *behavior layer*. It's what makes a page respond to user actions, validate input before it hits a server, update content without a full reload, and animate elements on screen.

:::key
HTML defines *what* is on the page, CSS defines *how it looks*, and JavaScript defines *what it does*.
:::

## Where JavaScript Runs

### The Browser Engine

Every modern browser ships with a built-in JavaScript engine:

- **Chrome / Edge** → V8
- **Firefox** → SpiderMonkey
- **Safari** → JavaScriptCore

When someone visits your page, *their* browser downloads your JS file and executes it locally. You don't need to install anything — if you have a browser, you already have a JavaScript runtime.

### Beyond the Browser: Node.js

In 2009, Ryan Dahl took Chrome's V8 engine and wrapped it in a standalone runtime called **Node.js**. This lets JavaScript run on servers, build tools, and command-line scripts — anywhere, not just the browser.

```
Browser JS → manipulates web pages, responds to clicks
Node.js    → runs servers, reads files, talks to databases
```

The language is the same. The environment and available APIs differ. In this course we focus on browser JavaScript, but knowing Node exists explains why JS is everywhere.

## Adding JavaScript to a Page

### Inline Script Tag

The simplest way — write JS directly inside a `<script>` tag:

```html
<!DOCTYPE html>
<html>
  <body>
    <h1>Hello</h1>

    <script>
      console.log("Page loaded!");
    </script>
  </body>
</html>
```

### External File (Recommended)

For real projects, keep JS in its own `.js` file:

```html
<body>
  <h1>Hello</h1>
  <script src="app.js"></script>
</body>
```

```js
// app.js
console.log("Code lives in its own file now.");
```

:::tip
Always use external files once a project grows beyond a few lines. It keeps HTML readable and lets the browser cache the script separately.
:::

### The `defer` Attribute

Placing `<script>` at the end of `<body>` ensures HTML is parsed first. But modern practice uses `defer` in the `<head>` instead:

```html
<head>
  <script src="app.js" defer></script>
</head>
<body>
  <h1>Hello</h1>
</body>
```

`defer` tells the browser: *download the script now, but don't run it until the HTML is fully parsed.* You get the best of both worlds — early download and safe execution order.

:::warning
Without `defer` or placement at the end of `<body>`, a script in the `<head>` blocks HTML parsing. The page will appear to hang while the script loads and runs.
:::

## console.log and the Dev Console

`console.log()` is your most important debugging tool. It prints values to the browser's **developer console**.

Open it: right-click the page → **Inspect** → **Console** tab (or press `F12`).

```js
console.log("Hello, world!");       // → Hello, world!
console.log(42);                    // → 42
console.log(3 + 4);                 // → 7
console.log("Score:", 10, "HP:", 5); // → Score: 10 HP: 5
```

Other useful console methods:

```js
console.warn("Heads up!");       // yellow warning icon
console.error("Something broke"); // red error icon
console.table([                   // renders a table
  { name: "Ada", level: 5 },
  { name: "Bob", level: 3 },
]);
```

:::tip
You can type JavaScript directly into the console and hit Enter. It's a live scratchpad — great for experimenting without creating a file.
:::

## JavaScript is Single-Threaded

JavaScript runs on **one thread**. That means it executes one instruction at a time, top to bottom. It cannot run two pieces of code simultaneously.

```js
console.log("First");
console.log("Second");
console.log("Third");
// Always prints: First, Second, Third — in that order.
```

This sounds limiting, but JS handles it with an **event loop** — a mechanism that queues up work (like network responses or timers) and processes it when the main thread is free. You'll learn the event loop in depth later, but the key takeaway now:

:::key
JavaScript is single-threaded. Long-running code blocks everything else — including the UI. Keep operations fast, and use asynchronous patterns (callbacks, promises) for slow work like network requests.
:::

## A Brief History: ES6 and Beyond

JavaScript was created in **10 days** in 1995 by Brendan Eich at Netscape. It was standardized as **ECMAScript** (ES).

For years the language barely changed. Then **ES6 (ES2015)** landed — the biggest update ever:

- `let` and `const` (replacing `var`)
- Arrow functions (`=>`)
- Template literals (backticks)
- Destructuring
- Classes
- Promises
- Modules (`import` / `export`)

Since ES6, JavaScript receives **yearly updates** (ES2016, ES2017, …). Modern features you'll see in this course — like optional chaining (`?.`), nullish coalescing (`??`), and `async/await` — all came from these annual releases.

:::tip
When you see "ES6+" in documentation or job postings, it means modern JavaScript — the version of the language you should be writing today.
:::

## Your First Program

Create an `index.html` and an `app.js` file in the same folder:

```html
<!-- index.html -->
<!DOCTYPE html>
<html>
  <head>
    <title>My First JS</title>
    <script src="app.js" defer></script>
  </head>
  <body>
    <h1 id="greeting">Hello</h1>
  </body>
</html>
```

```js
// app.js
const name = "World";
console.log(`Hello, ${name}!`);

const heading = document.getElementById("greeting");
heading.textContent = `Hello, ${name}!`;
```

Open `index.html` in your browser. The heading changes, and the console shows your message. That's JavaScript at work.

:::quiz
Q: What is JavaScript's role in a web page?
- Structure and content
- Visual styling and layout
- Behavior and interactivity *
E: JavaScript is the behavior layer. HTML handles structure, CSS handles presentation, and JS handles interactivity and logic.
:::

:::quiz
Q: What does the `defer` attribute do on a script tag?
- Prevents the script from running entirely
- Downloads the script immediately but delays execution until HTML is fully parsed *
- Makes the script run before any HTML is loaded
E: `defer` lets the browser download the script early without blocking HTML parsing. The script runs only after the document is ready.
:::

:::quiz
Q: Why is it important that JavaScript is single-threaded?
- It means JS can run multiple tasks at once
- Long-running code blocks everything else, including the UI *
- It makes JavaScript faster than other languages
E: A single thread means one task at a time. Heavy computation freezes the page until it finishes, which is why asynchronous patterns matter.
:::

## Recap

- JavaScript is the behavior layer of the web — it makes pages interactive.
- JS runs in every browser via a built-in engine (V8, SpiderMonkey, JavaScriptCore) and on servers via Node.js.
- Add JS to a page with `<script src="app.js" defer></script>` in the `<head>`, or a `<script>` tag at the end of `<body>`.
- `console.log()` prints to the dev console — your primary debugging tool.
- JavaScript is single-threaded: one instruction at a time, managed by the event loop.
- ES6 (2015) was the turning point. Modern JS receives yearly updates with powerful features.

**Next up:** Variables, Constants, and Data Types — the building blocks you'll use in every line of JavaScript you write.

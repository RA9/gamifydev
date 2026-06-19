# Building Interactive JavaScript Websites

You can build a beautiful page — but so far it just sits there. Now you'll make it **interactive**: responding to clicks, updating content on the fly, and reacting to what people type. This is the moment a web *page* becomes a web *app*.

By the end of this lesson you'll understand how JavaScript reads and changes the page through the DOM, and you'll trace a small interactive feature end to end.

## The DOM: your page as a living tree

When the browser loads your HTML, it doesn't just display it — it builds a **live model** of the page in memory called the **DOM** (Document Object Model). Every element becomes a node in a tree.

![The page is a tree — the DOM](/images/lessons/dom-tree.svg)

JavaScript works by reaching into this tree to **find** elements and **change** them. Change a node, and the browser instantly redraws that part of the page.

:::analogy
The DOM is like a puppet. The HTML you wrote is the puppet's starting pose; JavaScript is the hand inside it. Pull a string (change a node) and the puppet moves (the page updates) — instantly, without reloading.
:::

## The three moves of interactivity

Almost every interactive feature is built from the same three steps:

1. **Select** an element from the DOM.
2. **Listen** for an event on it.
3. **Update** the page in response.

```js
// 1. Select
const button = document.getElementById("like");
const count = document.getElementById("count");
let likes = 0;

// 2. Listen
button.addEventListener("click", function () {
  // 3. Update
  likes = likes + 1;
  count.textContent = likes;
});
```

Every click adds one to `likes` and rewrites the `count` element. That's a working "like" button — and the same pattern scales to entire apps.

:::quiz
Q: What is the DOM?
- A separate programming language
- The browser's live, in-memory tree of the page that JS can change *
- A CSS layout system
E: The DOM (Document Object Model) is the browser's living model of the page. JavaScript reads and changes it to make pages interactive.
:::

## Reading what users type

Interactivity isn't only clicks — you can read input, too:

```html
<input id="nameField" placeholder="Your name" />
<button id="greet">Greet me</button>
<p id="hello"></p>
```

```js
const field = document.getElementById("nameField");
const hello = document.getElementById("hello");

document.getElementById("greet").addEventListener("click", function () {
  hello.textContent = "Hello, " + field.value + "!";
});
```

`field.value` grabs whatever the user typed, and we drop it straight into the page. Read input → process it → update the DOM.

:::tip
Reach for `textContent` when you're setting plain text — it's safe. Avoid dumping untrusted text into the page as HTML; that's how cross-site scripting bugs sneak in.
:::

## Thinking in "state"

As features grow, it helps to separate your **data** (often called *state*) from how it's **shown**. In the like button, `likes` is the state; the `count` element is the view. The flow becomes a loop:

> an event changes the **state** → you **re-render** the view from that state.

This tiny idea — keep state, render from it — is the seed of every modern framework like React.

:::quiz
Q: A user clicks a button to add an item to a list. In the "three moves", what is clicking the button?
- Selecting an element
- The event you listen for *
- Updating the page
E: The click is the *event*. You select the button first, listen for its click, then update the page (add the item).
:::

:::key
JavaScript makes pages interactive by working on the **DOM** — the live tree of your page. The core pattern is **select → listen → update**, and as apps grow you separate **state** (your data) from the **view** (what's shown).
:::

## Talk about it

Explain in your own words:

> "What is the DOM, and what are the three steps behind almost every interactive feature?"

If you can describe the select → listen → update loop with an example, you're ready to build real interactive features.

## What's next

You can now build interactive pages. Next, **Git and GitHub** teaches you how to save your work properly and share it with the world.

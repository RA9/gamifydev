# The DOM: Selecting and Changing Elements

So far your JavaScript has talked to the console. Now it is time to talk to the *page itself* — finding elements, reading them, and changing them while the user watches. The tool that makes this possible is the DOM.

## What the DOM actually is

When a browser loads your HTML, it does not keep the raw text around. It parses that text into a living, in-memory model called the **DOM** — the Document Object Model. Every tag becomes an *object*, and those objects are nested exactly like your tags are nested, forming a tree.

Consider this HTML:

```html
<body>
  <h1 id="title">Welcome</h1>
  <ul class="list">
    <li>Apples</li>
    <li>Bananas</li>
  </ul>
</body>
```

The browser turns it into a tree: `body` has two children (`h1` and `ul`), and `ul` has two children (the two `li` elements). JavaScript can walk this tree, read any node, and change it. When you change a node, the browser instantly re-draws the page to match.

:::analogy
Think of the HTML file as a recipe printed on paper, and the DOM as the actual cake sitting on the counter. You can't change the cake by editing the paper after baking — but the DOM is a magic cake you *can* keep editing, and it re-bakes itself the instant you touch it.
:::

The global object `document` is your entry point into this tree. Everything below starts from `document`.

## Selecting elements

Before you can change something, you have to *find* it. There are three selectors worth knowing.

### querySelector — the everyday workhorse

`document.querySelector(cssSelector)` returns the **first** element that matches a CSS selector, or `null` if nothing matches.

```js
const title = document.querySelector("#title");   // by id
const list = document.querySelector(".list");     // by class
const firstItem = document.querySelector("li");   // by tag
const nested = document.querySelector("ul.list li"); // full CSS selector
```

Because it speaks CSS, one method covers nearly every case. The leading `#` means id, `.` means class, and a bare word means a tag name — exactly like in your stylesheets.

### querySelectorAll — many elements at once

`querySelectorAll` returns **all** matches as a **NodeList** (an array-like list).

```js
const items = document.querySelectorAll("li");
console.log(items.length); // 2
console.log(items[0].textContent); // "Apples"
```

A NodeList is not a real array, but you can loop it with `forEach`, or convert it with `Array.from(items)` if you need array methods like `map` or `filter`.

### getElementById — the classic

```js
const title = document.getElementById("title"); // note: no "#"
```

It is slightly faster and very readable, but `querySelector("#title")` does the same job, so most people just use `querySelector` everywhere for consistency.

:::warning
`querySelector` returns `null` when nothing matches. If you then write `null.textContent`, you get the dreaded "Cannot read properties of null". Double-check your selector spelling, and make sure your script runs *after* the element exists on the page.
:::

## Reading and changing content

Once you hold an element, `textContent` and `innerHTML` let you read or replace what is inside it.

```js
const title = document.querySelector("#title");

console.log(title.textContent);   // "Welcome"  (reads the text)
title.textContent = "Hello!";     // replaces the text
```

`textContent` deals in **plain text**. If you assign `"<b>Hi</b>"`, the user literally sees the characters `<b>Hi</b>` — the tags are not interpreted.

`innerHTML` deals in **HTML**. The browser parses whatever you assign:

```js
const box = document.querySelector(".list");
box.innerHTML = "<li>New item</li><li>Another</li>"; // becomes real <li> elements
```

That power comes with a serious caution.

:::warning
Never put text that came from a user (a comment, a username, a search box) directly into `innerHTML`. A malicious user could type `<img src=x onerror="steal()">`, and `innerHTML` would run that code — this is a **XSS** (cross-site scripting) attack. For user-supplied text, always use `textContent`. Reserve `innerHTML` for markup *you* wrote.
:::

There is also `value`, but that is specific to form inputs — we will get to it in a moment.

## Changing styles

You can change how an element looks in two ways. The quick way is the `style` property:

```js
const title = document.querySelector("#title");
title.style.color = "crimson";
title.style.backgroundColor = "black"; // note: camelCase, not background-color
title.style.fontSize = "32px";
```

Notice that CSS properties with dashes become camelCase in JavaScript: `background-color` → `backgroundColor`, `font-size` → `fontSize`.

Setting inline styles works, but it scatters design decisions through your JS and is hard to undo. The cleaner approach is to define classes in CSS and toggle them with **classList**.

```css
.highlight { background: yellow; color: black; }
.hidden { display: none; }
```

```js
const title = document.querySelector("#title");

title.classList.add("highlight");      // add a class
title.classList.remove("highlight");   // remove it
title.classList.toggle("hidden");      // add if missing, remove if present
const isOn = title.classList.contains("highlight"); // true / false
```

`toggle` is especially handy for things like show/hide and on/off states — one line flips it.

:::tip
Keep *appearance* in your CSS and *behavior* in your JS. Instead of setting ten style properties from JavaScript, write one CSS class and add it with `classList.add`. Your code stays readable and your designer stays happy.
:::

## Attributes and input values

Attributes are the extra settings on a tag: `href`, `src`, `alt`, `disabled`, `data-*`, and so on.

```js
const link = document.querySelector("a");

console.log(link.getAttribute("href")); // read an attribute
link.setAttribute("href", "https://gamifydev.com"); // change it
link.setAttribute("target", "_blank");
```

Form inputs are special: their *current* text lives in the `value` property, not in a `value` attribute.

```html
<input id="name" type="text" />
```

```js
const input = document.querySelector("#name");
console.log(input.value);   // whatever the user has typed
input.value = "Carlos";     // set it programmatically
```

:::key
For inputs, read and write `element.value`. For most other attributes, use `getAttribute` / `setAttribute`. Mixing these up is one of the most common beginner bugs.
:::

## Creating and inserting elements

You are not limited to elements that already exist — you can build new ones.

```js
const li = document.createElement("li"); // make a fresh <li> (not on page yet)
li.textContent = "Cherries";
li.classList.add("fruit");

const list = document.querySelector(".list");
list.append(li);   // add as the LAST child
list.prepend(li);  // ...or add as the FIRST child
```

`createElement` makes the element in memory; nothing appears until you *insert* it with `append` (end) or `prepend` (start). To delete an element, call `.remove()` on it:

```js
const first = document.querySelector(".list li");
first.remove(); // gone from the page
```

## Looping over many elements

`querySelectorAll` plus a loop lets you update a whole group at once.

```html
<ul class="todo">
  <li>Buy milk</li>
  <li>Walk dog</li>
  <li>Write code</li>
</ul>
```

```js
const items = document.querySelectorAll(".todo li");

items.forEach((item, index) => {
  item.textContent = `${index + 1}. ${item.textContent}`;
  item.classList.add("done");
});
```

After this runs, the list reads "1. Buy milk", "2. Walk dog", "3. Write code", each carrying the `done` class. One small loop, every item updated.

:::example
A "select all" effect: grab every checkbox and check it.
```js
const boxes = document.querySelectorAll("input[type=checkbox]");
boxes.forEach(box => { box.checked = true; });
```
:::

## Practice

:::predict
What does this print?
```js
const el = document.createElement("p");
el.textContent = "Hi";
console.log(document.querySelector("p"));
```
:::

:::quiz
Q: A user typed their name into a comment box. Where should you put that text to display it safely?
- el.innerHTML = userText
- el.textContent = userText *
- el.setAttribute("html", userText)
- document.write(userText)
E: User-supplied text in `innerHTML` risks an XSS attack because the browser would run any embedded tags or scripts. `textContent` treats it as plain text, so it is always safe to display.
:::

:::fill
To flip a class on or off in a single call, you use `element.classList.______("active")`.
- toggle
E: `toggle` adds the class if it is absent and removes it if it is present — perfect for on/off states.
:::

## Recap

- The **DOM** is your HTML turned into a tree of JavaScript objects; changing the tree re-draws the page instantly.
- Select with `querySelector` (first match), `querySelectorAll` (a NodeList of all matches), or `getElementById`.
- `textContent` is safe plain text; `innerHTML` parses HTML and is risky with user input (**XSS**).
- Style with `element.style.camelCase`, but prefer `classList.add/remove/toggle/contains` plus CSS classes.
- Use `getAttribute`/`setAttribute` for attributes, and `.value` for input contents.
- Build elements with `createElement`, place them with `append`/`prepend`, delete with `remove`.
- Loop a NodeList with `forEach` to update many elements at once.

**Next up:** JavaScript Events — making your page react when the user clicks, types, and submits.

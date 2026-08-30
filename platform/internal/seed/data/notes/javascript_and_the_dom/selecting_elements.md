# Selecting Elements

Before you can change anything on a page, you have to *find* it. The DOM gives you several ways to reach into the tree and grab the elements you need. This lesson covers every selector method worth knowing, the difference between the collections they return, and the habits that keep your code fast and bug-free.

## querySelector — the everyday workhorse

`document.querySelector(selector)` takes any CSS selector and returns the **first** matching element, or `null` if nothing matches.

```js
const title = document.querySelector("#page-title");   // by id
const hero = document.querySelector(".hero");          // by class
const firstParagraph = document.querySelector("p");    // by tag
const nested = document.querySelector("nav .link.active"); // compound selector
```

Because it speaks the same selector language as your CSS, one method handles nearly every case.

:::tip
If you can select it in CSS, you can select it with `querySelector`. Pseudo-classes like `:first-child`, `:not(.hidden)`, and `[data-role="admin"]` all work.
:::

## querySelectorAll — many elements at once

`querySelectorAll` returns **all** matches wrapped in a **NodeList**.

```js
const items = document.querySelectorAll(".todo-item");
console.log(items.length); // however many matched
console.log(items[0]);     // first element
```

You can loop a NodeList with `forEach`, but it is **not** a full array — methods like `map`, `filter`, and `reduce` are missing.

## NodeList vs HTMLCollection

There are two array-like types you will encounter:

| | NodeList | HTMLCollection |
|---|---|---|
| Returned by | `querySelectorAll` | `getElementsByClassName`, `children` |
| Has `forEach` | Yes | No |
| Live or static | Static (snapshot) | Live (updates automatically) |

A **live** collection changes when the DOM changes. A **static** NodeList is a snapshot taken at the moment you called `querySelectorAll` — safer and more predictable.

:::warning
`getElementsByClassName` returns a live HTMLCollection. If you add or remove elements during a loop, the collection shifts under you and you can skip items or loop forever. Prefer `querySelectorAll` unless you specifically need live behavior.
:::

## Converting to a real array with Array.from()

When you need array methods, convert with `Array.from()` or the spread operator:

```js
const items = document.querySelectorAll("li");

// Array.from
const texts = Array.from(items).map(el => el.textContent);

// spread into an array
const filtered = [...items].filter(el => !el.classList.contains("done"));
```

Both produce a true `Array` with full access to `map`, `filter`, `reduce`, `find`, and everything else.

## getElementById — the classic

```js
const title = document.getElementById("page-title"); // no "#" prefix
```

It is the oldest selector and marginally faster than `querySelector`, but it only works with ids. Most codebases just use `querySelector` everywhere for consistency.

:::key
`getElementById("title")` takes a bare id — no `#`. `querySelector("#title")` takes a CSS selector — the `#` is required. Mixing these up is a common source of silent `null` returns.
:::

## Handling null — the most common DOM bug

When nothing matches, `querySelector` and `getElementById` return `null`. Trying to read a property on `null` throws immediately:

```js
const el = document.querySelector(".typo-class");
el.textContent = "Hello"; // TypeError: Cannot read properties of null
```

Always guard against this:

```js
const el = document.querySelector(".status");
if (el) {
  el.textContent = "Ready";
}
```

Or, when the element *must* exist, treat the null as a real bug and fail loudly so you catch it during development:

```js
const el = document.querySelector("#app");
if (!el) throw new Error("Missing #app element — check your HTML");
```

## Scoped queries — searching within an element

`querySelector` and `querySelectorAll` can be called on *any* element, not just `document`. When called on an element, the search is scoped to its descendants:

```js
const sidebar = document.querySelector(".sidebar");

// only finds links INSIDE the sidebar
const sidebarLinks = sidebar.querySelectorAll("a");

const mainContent = document.querySelector("main");
// only the first <h2> inside main
const heading = mainContent.querySelector("h2");
```

This is essential when the same class appears in different sections and you need to target a specific one.

:::tip
Scoped queries keep your selectors simple. Instead of writing `document.querySelector(".sidebar .nav .link.active")`, grab the sidebar once and then query inside it. Shorter selectors are easier to read and less likely to break when the HTML structure changes.
:::

## Caching selectors in variables

Every call to `querySelector` walks the DOM tree. If you need the same element in multiple places, select it once and store the result:

```js
// Bad — selects the same element three times
document.querySelector("#count").textContent = "0";
document.querySelector("#count").classList.add("visible");
document.querySelector("#count").style.color = "green";

// Good — one lookup, reuse the reference
const count = document.querySelector("#count");
count.textContent = "0";
count.classList.add("visible");
count.style.color = "green";
```

For elements you use throughout a file, cache them at the top:

```js
// Cache all your selectors up front
const form = document.querySelector("#task-form");
const input = document.querySelector("#task-input");
const list = document.querySelector("#task-list");
const status = document.querySelector("#status");

// Now use them freely in event handlers and render functions
form.addEventListener("submit", (e) => {
  e.preventDefault();
  const value = input.value.trim();
  // ...
});
```

This pattern is cleaner, faster, and makes it obvious which elements your script depends on.

## Selecting by attribute

Attribute selectors from CSS work in `querySelector` too:

```js
// elements with a specific data attribute
const draggables = document.querySelectorAll("[draggable]");

// inputs of a specific type
const checkboxes = document.querySelectorAll("input[type='checkbox']");

// links pointing to external sites
const external = document.querySelectorAll("a[target='_blank']");
```

## Practice

:::quiz
Q: What does `document.querySelector(".card")` return if no element has the class `card`?
- An empty NodeList
- An empty string
- undefined
- null *
E: `querySelector` returns the first match or `null`. It does not return a list — that is `querySelectorAll`, which would return an empty NodeList.
:::

:::quiz
Q: You need to find all `.tag` elements inside a `.sidebar` element you already have stored in a variable. What is the best call?
- `document.querySelectorAll(".sidebar .tag")`
- `sidebar.querySelectorAll(".tag")` *
- `document.querySelector(".tag")`
- `sidebar.querySelector(".tag")`
E: Calling `querySelectorAll` on the `sidebar` element scopes the search to its descendants. This is cleaner, shorter, and doesn't break if another `.tag` exists elsewhere on the page.
:::

## Recap

- `querySelector` returns the **first** match (or `null`); `querySelectorAll` returns a **static NodeList** of all matches.
- `getElementById` takes a bare id (no `#`) and is slightly faster, but `querySelector` is more versatile.
- A **NodeList** has `forEach` but lacks `map`/`filter` — convert with `Array.from()` or `[...nodeList]`.
- An **HTMLCollection** (from `getElementsByClassName`, `children`) is live and has no `forEach` — prefer `querySelectorAll`.
- Always check for `null` before using a selected element, or throw early to catch bugs.
- **Scope queries** by calling `querySelector` on a parent element instead of `document`.
- **Cache selectors** in variables at the top of your script — avoid repeated DOM lookups.

**Next up:** Changing Content and Attributes — what to do once you've found your elements.

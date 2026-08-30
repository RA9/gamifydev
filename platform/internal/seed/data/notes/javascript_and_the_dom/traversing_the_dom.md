# Traversing the DOM

Sometimes you do not know the exact selector for an element, but you know where it is *relative* to another element — its parent, its first child, the sibling right after it. DOM traversal lets you walk the tree in any direction: up, down, or sideways.

## Walking up — parentElement

`parentElement` gives you the direct parent of an element:

```js
const item = document.querySelector(".todo-item");
const list = item.parentElement; // the <ul> or <ol> containing this item
const wrapper = list.parentElement; // whatever wraps the list
```

You can chain calls to walk multiple levels up:

```js
const grandparent = item.parentElement.parentElement;
```

But chaining is fragile — if the HTML structure changes by one level, it breaks. For walking up to a specific ancestor, `closest()` is the right tool (covered below).

## Walking down — children, firstElementChild, lastElementChild

`children` returns an HTMLCollection of an element's direct child *elements* (no text nodes, no comments):

```js
const nav = document.querySelector("nav");
console.log(nav.children);          // HTMLCollection of child elements
console.log(nav.children.length);   // how many direct children
console.log(nav.children[0]);       // first child element
console.log(nav.children[2]);       // third child element
```

Shortcut properties for common positions:

```js
const list = document.querySelector("ul");

const first = list.firstElementChild; // first <li>
const last = list.lastElementChild;   // last <li>
```

:::tip
There are also `firstChild` and `lastChild` (without "Element"), but those include text nodes (like whitespace between tags). Stick with the `Element` versions — they skip text nodes and give you what you almost always want.
:::

## Walking sideways — nextElementSibling, previousElementSibling

Siblings are elements that share the same parent. You can step forward or backward:

```js
const second = document.querySelector("#item-2");

const next = second.nextElementSibling;     // the element right after
const prev = second.previousElementSibling; // the element right before
```

If there is no next or previous sibling, the property is `null`.

```js
const first = list.firstElementChild;
console.log(first.previousElementSibling); // null — nothing before the first
```

### Walking through all siblings

```js
const list = document.querySelector("ul");
let current = list.firstElementChild;

while (current) {
  console.log(current.textContent);
  current = current.nextElementSibling;
}
```

This walks every child of the list without needing a selector for each one.

## closest() — finding the nearest ancestor

`closest(selector)` walks **up** the tree from an element and returns the first ancestor that matches the selector, or `null` if none match. It also checks the element itself.

```js
const deleteBtn = document.querySelector(".delete-btn");

// Find the nearest parent with class "card"
const card = deleteBtn.closest(".card");

// Find the nearest <form> ancestor
const form = deleteBtn.closest("form");
```

This is incredibly useful in event delegation, where a click might land on a child element and you need to find the meaningful parent:

```js
list.addEventListener("click", (event) => {
  const item = event.target.closest(".todo-item");
  if (!item) return; // click was not on a todo item

  const id = item.dataset.id;
  // handle the click for this item
});
```

:::key
`closest` walks *up* from the element (including itself). `querySelector` walks *down* from the element (its descendants). They search in opposite directions.
:::

## matches() — checking if an element fits a selector

`matches(selector)` returns `true` if the element itself matches the given CSS selector:

```js
const el = document.querySelector(".card");

el.matches(".card");           // true
el.matches(".card.featured");  // true if it also has "featured"
el.matches("div");             // true if the element is a <div>
el.matches(".sidebar");        // false
```

This is useful for filtering or branching inside a loop:

```js
const items = document.querySelectorAll("li");

items.forEach(item => {
  if (item.matches(".urgent")) {
    item.style.fontWeight = "bold";
  }
});
```

## Practical patterns

### Pattern 1: Find and update a specific sibling

A "move up" button in a sortable list:

```js
function moveUp(item) {
  const prev = item.previousElementSibling;
  if (prev) {
    item.parentElement.insertBefore(item, prev);
  }
}
```

### Pattern 2: Walk up to find a data attribute

A nested button inside a card needs the card's id:

```html
<div class="card" data-id="42">
  <div class="card-body">
    <button class="edit-btn">Edit</button>
  </div>
</div>
```

```js
editBtn.addEventListener("click", (event) => {
  const card = event.target.closest(".card");
  const id = card.dataset.id; // "42"
  openEditor(id);
});
```

Without `closest`, you would need fragile chaining: `event.target.parentElement.parentElement.dataset.id`.

### Pattern 3: Get all siblings of an element

```js
function getSiblings(el) {
  return Array.from(el.parentElement.children).filter(child => child !== el);
}

const siblings = getSiblings(document.querySelector("#item-3"));
```

### Pattern 4: Walking down to collect nested content

```js
const section = document.querySelector("#features");
const headings = [];
let child = section.firstElementChild;

while (child) {
  if (child.matches("h3")) {
    headings.push(child.textContent);
  }
  child = child.nextElementSibling;
}
```

## The full traversal map

Starting from any element, here is every direction you can go:

| Direction | Property/Method | Returns |
|---|---|---|
| Up (parent) | `parentElement` | Single element or `null` |
| Up (ancestor) | `closest(selector)` | First matching ancestor or `null` |
| Down (all children) | `children` | HTMLCollection |
| Down (first) | `firstElementChild` | Single element or `null` |
| Down (last) | `lastElementChild` | Single element or `null` |
| Down (query) | `querySelector(selector)` | First match or `null` |
| Sideways (next) | `nextElementSibling` | Single element or `null` |
| Sideways (prev) | `previousElementSibling` | Single element or `null` |
| Self check | `matches(selector)` | `true` or `false` |

## Practice

:::quiz
Q: A click event fires on a `<span>` inside a `<button>` inside a `.card`. What is the cleanest way to find the `.card` from the event?
- `event.target.parentElement.parentElement`
- `document.querySelector(".card")`
- `event.target.closest(".card")` *
- `event.target.matches(".card")`
E: `closest` walks up the tree and returns the first ancestor matching the selector. It works regardless of how deeply nested the click target is, unlike chained `parentElement` calls that break if the nesting changes.
:::

:::quiz
Q: What does `firstElementChild` return if the parent has no child elements?
- An empty string
- An empty array
- undefined
- null *
E: All traversal properties return `null` when there is nothing to find. This is consistent with `querySelector`, `parentElement`, and the sibling properties.
:::

## Recap

- **`parentElement`** moves one level up; **`closest(selector)`** walks up until it finds a match.
- **`children`** returns an HTMLCollection of child elements; **`firstElementChild`** and **`lastElementChild`** jump to the ends.
- **`nextElementSibling`** and **`previousElementSibling`** step sideways between siblings.
- **`closest()`** searches up (including the element itself); **`querySelector()`** searches down. They are opposites.
- **`matches(selector)`** checks if an element itself fits a CSS selector — useful for filtering.
- Use `closest` over chained `parentElement` — it is resilient to HTML structure changes.

**Next up:** Event Delegation and Bubbling — one listener to rule them all.

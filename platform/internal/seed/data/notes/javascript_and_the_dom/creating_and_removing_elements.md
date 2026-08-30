# Creating and Removing Elements

You are not limited to elements that already exist on the page. JavaScript lets you create new nodes, insert them anywhere in the DOM tree, and remove them when they are no longer needed. This lesson covers every tool for building and tearing down the DOM programmatically.

## createElement — building a new element

`document.createElement(tagName)` creates a new element in memory. It does **not** appear on the page until you insert it.

```js
const li = document.createElement("li");
li.textContent = "Buy groceries";
li.classList.add("task");
// The element exists, but it is floating in memory — not visible yet
```

You can set any property on a freshly created element before inserting it: `textContent`, `classList`, `setAttribute`, `dataset`, event listeners — everything.

## Inserting elements into the page

Once built, you need to place the element somewhere in the tree.

### append and prepend

```js
const list = document.querySelector("#todo-list");
const li = document.createElement("li");
li.textContent = "New task";

list.append(li);   // inserts as the LAST child
list.prepend(li);  // inserts as the FIRST child
```

`append` and `prepend` also accept plain strings (they become text nodes):

```js
list.append("Some extra text");
```

You can pass multiple items at once:

```js
const a = document.createElement("li");
a.textContent = "First";
const b = document.createElement("li");
b.textContent = "Second";

list.append(a, b); // both added at the end, in order
```

### before and after — sibling insertion

```js
const ref = document.querySelector("#item-3");
const newItem = document.createElement("li");
newItem.textContent = "Inserted";

ref.before(newItem); // places it immediately BEFORE item-3
ref.after(newItem);  // places it immediately AFTER item-3
```

### appendChild — the classic

`appendChild` does the same job as `append` but only accepts a single node (no strings, no multiple arguments):

```js
list.appendChild(li); // same as list.append(li) for a single element
```

You will see `appendChild` in older code. For new code, prefer `append` — it is more flexible.

:::tip
Here is a quick decision guide:
- **End of parent:** `parent.append(el)`
- **Start of parent:** `parent.prepend(el)`
- **Before a sibling:** `sibling.before(el)`
- **After a sibling:** `sibling.after(el)`
:::

## insertAdjacentHTML — injecting HTML strings

When you have a chunk of HTML as a string and want to insert it without wiping existing content (which `innerHTML =` would do), use `insertAdjacentHTML`:

```js
const list = document.querySelector("#todo-list");

list.insertAdjacentHTML("beforeend", "<li>Appended item</li>");
list.insertAdjacentHTML("afterbegin", "<li>Prepended item</li>");
```

The first argument is the position:

| Position | Where it goes |
|---|---|
| `"beforebegin"` | Before the element itself (as a sibling) |
| `"afterbegin"` | Inside, before the first child |
| `"beforeend"` | Inside, after the last child |
| `"afterend"` | After the element itself (as a sibling) |

:::warning
`insertAdjacentHTML` parses raw HTML, so the same XSS rules apply: never pass user input into it. Use it only with markup you control.
:::

## cloneNode — copying an existing element

`cloneNode` creates a copy of an element. Pass `true` for a **deep** clone (includes all descendants) or `false` for a **shallow** clone (just the element, no children):

```js
const template = document.querySelector(".card-template");

const clone = template.cloneNode(true); // deep copy — includes inner elements
clone.querySelector(".title").textContent = "New Card";
clone.classList.remove("card-template");

document.querySelector("#card-grid").append(clone);
```

:::key
`cloneNode(true)` copies the element and all its children. `cloneNode(false)` copies only the element itself (empty). Event listeners are **not** copied — you must add them to the clone separately.
:::

## remove() — taking an element off the page

Call `.remove()` on any element to pull it out of the DOM:

```js
const item = document.querySelector("#task-7");
item.remove(); // gone from the page
```

The element still exists in memory (if you have a reference to it), so you could re-insert it later. But if nothing references it, the garbage collector reclaims it.

## replaceWith — swapping one element for another

```js
const oldHeading = document.querySelector("h1");
const newHeading = document.createElement("h2");
newHeading.textContent = oldHeading.textContent;

oldHeading.replaceWith(newHeading);
// The <h1> is gone; the <h2> now sits in its place
```

You can also pass a string to replace an element with a text node:

```js
el.replaceWith("Just plain text now");
```

## DocumentFragment — batch inserts for performance

When you need to add many elements, inserting them one by one causes the browser to re-layout the page after each insertion. A `DocumentFragment` lets you build up a batch in memory, then insert everything in one operation:

```js
const fragment = document.createDocumentFragment();

for (let i = 0; i < 100; i++) {
  const li = document.createElement("li");
  li.textContent = `Item ${i + 1}`;
  fragment.append(li);
}

// One DOM insert, one re-layout
document.querySelector("#big-list").append(fragment);
```

After appending, the fragment is empty — its children moved into the DOM.

:::tip
For small lists (under ~50 items), the performance difference is negligible. But when you are rendering hundreds of elements — search results, log entries, data tables — a DocumentFragment avoids visible jank.
:::

## Full pattern — building a dynamic list

Here is everything in action: create elements, set content, batch them in a fragment, and insert:

```js
function renderTasks(tasks) {
  const list = document.querySelector("#task-list");
  list.innerHTML = ""; // clear previous items

  const fragment = document.createDocumentFragment();

  tasks.forEach(task => {
    const li = document.createElement("li");
    li.classList.toggle("done", task.completed);
    li.dataset.id = task.id;

    const text = document.createElement("span");
    text.textContent = task.title;

    const deleteBtn = document.createElement("button");
    deleteBtn.textContent = "Delete";
    deleteBtn.setAttribute("aria-label", `Delete ${task.title}`);

    li.append(text, deleteBtn);
    fragment.append(li);
  });

  list.append(fragment);
}
```

This is the pattern used in real applications: build in memory, assemble, then insert once.

## Practice

:::quiz
Q: What is the difference between `createElement` and `insertAdjacentHTML`?
- `createElement` creates a single element object you can configure before inserting; `insertAdjacentHTML` parses an HTML string directly into the DOM *
- `createElement` is faster because it skips parsing
- `insertAdjacentHTML` creates element objects like `createElement`
- There is no difference
E: `createElement` gives you a reference to a new DOM node you can set up with properties and listeners before inserting. `insertAdjacentHTML` takes a raw HTML string and parses it into the DOM at a specified position. The trade-off is flexibility vs convenience.
:::

:::quiz
Q: You are adding 200 list items in a loop. What should you use to avoid 200 separate page re-layouts?
- innerHTML
- replaceWith
- DocumentFragment *
- cloneNode
E: A `DocumentFragment` lets you build all 200 items in memory and insert them in a single DOM operation, causing just one re-layout instead of 200.
:::

## Recap

- **`createElement(tag)`** makes a new element in memory — nothing appears until you insert it.
- **`append`/`prepend`** add children at the end/start; **`before`/`after`** insert as siblings.
- **`appendChild`** is the older, less flexible version of `append` — prefer `append` in new code.
- **`insertAdjacentHTML(position, html)`** injects HTML strings at four possible positions without clearing existing content.
- **`cloneNode(true)`** deep-copies an element and its descendants (but not event listeners).
- **`remove()`** pulls an element out of the page; **`replaceWith(el)`** swaps one element for another.
- **`DocumentFragment`** batches multiple inserts into one DOM operation for better performance.

**Next up:** Traversing the DOM — walking between parent, child, and sibling elements.

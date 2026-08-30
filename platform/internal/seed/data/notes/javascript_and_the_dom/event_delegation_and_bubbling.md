# Event Delegation and Bubbling

When you have a list of 50 items, you do not need 50 event listeners. Thanks to how events travel through the DOM, one listener on the parent can handle clicks on any child. This lesson explains the event flow and shows you how to use it to write cleaner, faster, and more flexible code.

## How events travel — bubbling

When you click a button, the browser does not only fire the event on that button. The event **bubbles** — it fires on the element you clicked, then on its parent, then on *its* parent, all the way up to `document`.

```html
<div id="container">
  <ul id="list">
    <li><button class="delete-btn">Delete</button></li>
  </ul>
</div>
```

If you click the `<button>`, the `click` event fires on:

1. `button.delete-btn`
2. `li`
3. `ul#list`
4. `div#container`
5. `body`
6. `html`
7. `document`

Every ancestor gets a chance to hear the event. This is called the **bubbling phase**, and it is the default behavior.

## Capturing — the other direction

There is also a **capturing phase** that runs *before* bubbling: the event travels down from `document` to the target element. You rarely need it, but you can opt into it:

```js
document.addEventListener("click", handler, true);
// or
document.addEventListener("click", handler, { capture: true });
```

The full flow: capture down → arrive at target → bubble up. In practice, you almost always work with the bubbling phase.

:::key
Events bubble up by default. A listener on a parent element will hear events from all of its descendants. This is the foundation of event delegation.
:::

## event.target vs event.currentTarget

These two properties answer different questions:

- **`event.target`** — the element the user actually clicked (the deepest element where the event originated).
- **`event.currentTarget`** — the element the listener is attached to.

```js
const list = document.querySelector("#list");

list.addEventListener("click", (event) => {
  console.log(event.target);        // the <button> the user clicked
  console.log(event.currentTarget); // the <ul> where the listener lives
});
```

When you delegate, `event.target` is how you figure out *which* child was interacted with.

## stopPropagation — stopping the bubble

Sometimes you need to prevent an event from reaching parent listeners. `stopPropagation` halts the bubble:

```js
const inner = document.querySelector(".modal-content");

inner.addEventListener("click", (event) => {
  event.stopPropagation(); // click stays here, does not bubble to the backdrop
});

const backdrop = document.querySelector(".modal-backdrop");
backdrop.addEventListener("click", () => {
  closeModal();
});
```

:::warning
Use `stopPropagation` sparingly. It makes events "disappear" for any parent listener, which can cause confusing bugs when other parts of your code depend on hearing those events. Most of the time, a careful check with `closest` or `matches` is better than stopping propagation.
:::

## Event delegation

**Event delegation** means putting one listener on a parent and using `event.target` to figure out which child triggered it. Instead of this:

```js
// One listener per item — does not work for items added later
document.querySelectorAll(".todo-item").forEach(item => {
  item.addEventListener("click", () => handleClick(item));
});
```

You write this:

```js
// One listener on the parent — works for every child, even new ones
const list = document.querySelector("#todo-list");

list.addEventListener("click", (event) => {
  const item = event.target.closest(".todo-item");
  if (!item) return; // click was not on a todo item
  handleClick(item);
});
```

### Why this is better

1. **One listener instead of many.** Fewer listeners means less memory and faster setup.
2. **Handles dynamic content.** Elements added to the list after the listener was set still get handled — no need to rewire listeners.
3. **Cleaner cleanup.** One listener to remove instead of tracking dozens.

:::tip
The pattern is always the same: listen on a stable parent, use `event.target.closest(selector)` to find the meaningful element, and bail out with `if (!el) return` if the click landed elsewhere.
:::

## Delegating with closest()

`closest` is the key to reliable delegation. The user might click on a `<span>` inside a `<button>` inside a `.card`. `event.target` would be the `<span>`, not the `.card`. `closest` walks up to find it:

```js
const grid = document.querySelector("#card-grid");

grid.addEventListener("click", (event) => {
  // Handle delete buttons
  const deleteBtn = event.target.closest(".delete-btn");
  if (deleteBtn) {
    const card = deleteBtn.closest(".card");
    deleteCard(card.dataset.id);
    return;
  }

  // Handle edit buttons
  const editBtn = event.target.closest(".edit-btn");
  if (editBtn) {
    const card = editBtn.closest(".card");
    openEditor(card.dataset.id);
    return;
  }
});
```

One listener on the grid handles both delete and edit actions on any card, including cards added later.

## Delegation for dynamic lists

This is where delegation truly shines. Without it, dynamically created elements need listeners attached every time you render:

```js
// Without delegation — must add listeners every time you add items
function renderItems(items) {
  const list = document.querySelector("#list");
  list.innerHTML = "";

  items.forEach(item => {
    const li = document.createElement("li");
    li.textContent = item.text;

    const btn = document.createElement("button");
    btn.textContent = "Delete";
    btn.addEventListener("click", () => deleteItem(item.id)); // new listener each time
    li.append(btn);
    list.append(li);
  });
}
```

```js
// With delegation — listener set up once, works forever
const list = document.querySelector("#list");

list.addEventListener("click", (event) => {
  const btn = event.target.closest("button");
  if (!btn) return;

  const li = btn.closest("li");
  const id = li.dataset.id;
  deleteItem(id);
});

function renderItems(items) {
  list.innerHTML = "";
  items.forEach(item => {
    const li = document.createElement("li");
    li.dataset.id = item.id;
    li.innerHTML = `<span>${item.text}</span><button>Delete</button>`;
    list.append(li);
  });
}
```

The delegated version is simpler, works for any number of items, and does not create new listener functions on every render.

## Practical patterns

### Pattern 1: Tab switching with delegation

```js
const tabBar = document.querySelector("#tabs");

tabBar.addEventListener("click", (event) => {
  const tab = event.target.closest("[data-tab]");
  if (!tab) return;

  activeTab = tab.dataset.tab;
  render();
});
```

### Pattern 2: Multiple action buttons in a table row

```js
const table = document.querySelector("#users-table");

table.addEventListener("click", (event) => {
  const row = event.target.closest("tr");
  if (!row) return;

  const userId = row.dataset.id;

  if (event.target.closest(".edit-btn"))   editUser(userId);
  if (event.target.closest(".delete-btn")) deleteUser(userId);
  if (event.target.closest(".view-btn"))   viewUser(userId);
});
```

## Practice

:::quiz
Q: You add a new `<li>` to a list after the page loads. A listener directly on each `<li>` was set up on page load. Does the new `<li>` respond to clicks?
- Yes, all elements automatically get listeners
- No, the listener was not added to the new element *
- Yes, because events bubble
- No, because events do not work on dynamic elements
E: Listeners added with `addEventListener` are bound only to the elements that existed at the time. New elements need their own listener — or you use delegation on the parent, which catches clicks from any child, including new ones.
:::

:::quiz
Q: In a delegated listener on a `<ul>`, the user clicks a `<strong>` inside an `<li>`. What does `event.target` point to?
- The `<ul>`
- The `<li>`
- The `<strong>` *
- The `document`
E: `event.target` is always the deepest element where the event originated — the actual element the user clicked. Use `event.target.closest("li")` to find the meaningful ancestor.
:::

## Recap

- Events **bubble up** from the target element through every ancestor to `document`.
- **`event.target`** is the element the user interacted with; **`event.currentTarget`** is where the listener lives.
- **`stopPropagation`** halts the bubble — use it sparingly.
- **Event delegation:** listen on a stable parent, use `event.target.closest(selector)` to find the relevant child.
- Delegation handles **dynamic content** automatically — elements added later are covered by the existing listener.
- One delegated listener replaces many individual listeners, saving memory and reducing complexity.

**Next up:** Form Events and Input Handling — making forms work the right way.

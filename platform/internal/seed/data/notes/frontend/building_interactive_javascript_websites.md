# Building Interactive JavaScript Websites

You now know how to select elements (the DOM) and how to react to the user (events). This lesson ties them together by *building real features* — each with complete, copy-pasteable code — and teaches the single most important pattern in frontend work: **state → render**.

## The state → render pattern

As features grow, it is tempting to poke at the DOM directly everywhere: add a class here, change text there, remove a node somewhere else. That quickly becomes a tangle where the page and your variables disagree.

The cleaner mental model is:

1. **State** — plain JavaScript data (variables, arrays, objects) that is the *single source of truth* about what should be on screen.
2. **Render** — one function that reads the state and updates the DOM to match.
3. **Events** — change the state, then call render again.

:::key
The golden rule: **events change state, then render reads state and rebuilds the UI.** You never update the screen directly from an event handler — you update the data and re-render. The screen becomes a *picture of your data*, and the two can never drift apart.
:::

:::analogy
State is the script of a play; render is the actors performing it. When the story changes, you rewrite the script and the actors perform the new version — you do not run on stage and reposition each actor by hand.
:::

:::project
We will build four small features you would actually ship: a persistent theme toggle, a live character counter, a to-do list driven by an array, and a tabbed interface. Make a single HTML file, paste each example in, and watch them come alive. Read the comments — they explain the *why*, not just the *what*.
:::

## Feature 1: Dark/light theme toggle with localStorage

`localStorage` is a tiny key-value store the browser keeps even after the page closes. We use it to remember the user's theme choice between visits.

```css
body { background: white; color: #111; transition: 0.2s; }
body.dark { background: #111; color: #eee; }
```

```html
<button id="themeBtn">Toggle theme</button>
```

```js
const themeBtn = document.querySelector("#themeBtn");

// On load: apply the saved theme (state lives in localStorage)
const saved = localStorage.getItem("theme");
if (saved === "dark") {
  document.body.classList.add("dark");
}

themeBtn.addEventListener("click", () => {
  // toggle returns true if the class is now present
  const isDark = document.body.classList.toggle("dark");
  // persist the new state so it survives a reload
  localStorage.setItem("theme", isDark ? "dark" : "light");
});
```

Reload the page after toggling — the choice sticks. That is the whole appeal of `localStorage`.

:::warning
`localStorage` only stores **strings**. Numbers and booleans get converted to text, and objects/arrays must be saved with `JSON.stringify` and read back with `JSON.parse`. Storing an object directly gives you the useless string `"[object Object]"`.
:::

## Feature 2: Live character counter

A textarea with a running count and a warning when the user nears the limit. The `input` event fires on every keystroke, so the count is always current.

```html
<textarea id="bio" maxlength="100" placeholder="Tell us about yourself"></textarea>
<p id="count">0 / 100</p>
```

```css
#count.warn { color: crimson; font-weight: bold; }
```

```js
const bio = document.querySelector("#bio");
const count = document.querySelector("#count");
const LIMIT = 100;

function render() {
  const used = bio.value.length;          // read state (the text)
  count.textContent = `${used} / ${LIMIT}`;
  count.classList.toggle("warn", used > LIMIT - 10); // warn near the end
}

bio.addEventListener("input", render);
render(); // run once so the initial "0 / 100" is correct
```

Notice the pattern already: a `render` function that reads the current value and updates the DOM, called both on `input` and once at startup.

:::tip
Always call your render function once when the page loads, not only inside the event handler. Otherwise the UI starts in a stale or blank state until the user interacts. Here it ensures "0 / 100" shows immediately.
:::

## Feature 3: A to-do list (array as the source of truth)

This is the showcase of state → render. The list of tasks is an **array of objects**. Events only ever modify that array, then call `render`, which rebuilds the whole list from scratch.

```html
<form id="todoForm">
  <input id="todoInput" placeholder="Add a task" />
  <button type="submit">Add</button>
</form>
<ul id="todoList"></ul>
```

```css
.todo.done span { text-decoration: line-through; opacity: 0.6; }
```

```js
const form = document.querySelector("#todoForm");
const input = document.querySelector("#todoInput");
const list = document.querySelector("#todoList");

// STATE: the single source of truth
let todos = [];

// RENDER: rebuild the list to match the state
function render() {
  list.innerHTML = ""; // clear the current DOM

  todos.forEach((todo) => {
    const li = document.createElement("li");
    li.className = "todo" + (todo.done ? " done" : "");

    const text = document.createElement("span");
    text.textContent = todo.text;

    const doneBtn = document.createElement("button");
    doneBtn.textContent = todo.done ? "Undo" : "Done";
    doneBtn.addEventListener("click", () => {
      todo.done = !todo.done; // change state
      render();               // re-render
    });

    const delBtn = document.createElement("button");
    delBtn.textContent = "Delete";
    delBtn.addEventListener("click", () => {
      todos = todos.filter((t) => t !== todo); // change state
      render();                                // re-render
    });

    li.append(text, doneBtn, delBtn);
    list.append(li);
  });
}

// EVENT: adding a task changes state, then renders
form.addEventListener("submit", (event) => {
  event.preventDefault();
  const value = input.value.trim();
  if (!value) return; // ignore empty input

  todos.push({ text: value, done: false }); // change state
  input.value = "";                          // clear the box
  render();                                  // re-render
});

render(); // start with an empty (but correct) list
```

Study what *never* happens here: no handler reaches in to cross out one specific `<li>` or surgically remove one node. Every change edits the `todos` array and re-renders. The DOM is always a faithful picture of the array.

:::example
Want to persist the to-dos across reloads? Combine this with `localStorage`:
```js
// save after every change
function save() { localStorage.setItem("todos", JSON.stringify(todos)); }
// load at startup
todos = JSON.parse(localStorage.getItem("todos") || "[]");
```
Call `save()` at the end of `render`, and the load line replaces `let todos = []`. Now your list survives a refresh.
:::

:::warning
Rebuilding the whole list with `innerHTML = ""` is simple and perfect for small lists. For very large or animated lists it can feel heavy, since you are recreating every element each time. That trade-off is exactly what frameworks like React optimize for — but for learning and for most real pages, full re-render is clear and fast enough.
:::

## Feature 4: A tabbed interface

Tabs are a small state machine: exactly one tab is active at a time. The "active tab" is a piece of state; render reflects it. Event delegation keeps it tidy.

```html
<div id="tabs">
  <button data-tab="home" class="tab">Home</button>
  <button data-tab="profile" class="tab">Profile</button>
  <button data-tab="settings" class="tab">Settings</button>
</div>

<div class="panel" data-panel="home">Welcome home.</div>
<div class="panel" data-panel="profile">Your profile.</div>
<div class="panel" data-panel="settings">Adjust settings.</div>
```

```css
.panel { display: none; }
.panel.active { display: block; }
.tab.active { font-weight: bold; border-bottom: 2px solid teal; }
```

```js
const tabBar = document.querySelector("#tabs");
const tabs = document.querySelectorAll(".tab");
const panels = document.querySelectorAll(".panel");

let activeTab = "home"; // STATE: which tab is selected

function render() {
  tabs.forEach((tab) => {
    tab.classList.toggle("active", tab.dataset.tab === activeTab);
  });
  panels.forEach((panel) => {
    panel.classList.toggle("active", panel.dataset.panel === activeTab);
  });
}

// EVENT: one delegated listener for all tabs
tabBar.addEventListener("click", (event) => {
  const tab = event.target.closest(".tab");
  if (!tab) return;
  activeTab = tab.dataset.tab; // change state
  render();                    // re-render
});

render(); // show the default tab on load
```

`tab.dataset.tab` reads the `data-tab` attribute — a clean way to attach custom data to elements. One listener on the bar handles every tab via delegation, and render uses `classList.toggle(name, condition)` to switch panels by simply comparing each element's data to `activeTab`.

:::tip
`classList.toggle("active", someBoolean)` adds the class when the boolean is true and removes it when false. It is the perfect render tool: loop over your elements and toggle based on whether each one matches the current state.
:::

## Practice

:::quiz
Q: In the to-do list, what happens directly when the user clicks "Delete"?
- The matching <li> element is removed from the DOM by hand
- The todos array is filtered to drop that item, then render() rebuilds the list *
- localStorage is cleared
- The form is submitted again
E: Following state → render, the handler changes the *data* (filters the array) and calls render, which rebuilds the DOM to match. The UI is always derived from the array, never edited directly.
:::

:::predict
A user toggles dark mode, then closes and reopens the browser tab. With the theme-toggle code above, what theme do they see, and why?
:::

:::fill
Because `localStorage` only stores strings, you must save an array of to-dos with `JSON.______(todos)` and read it back with `JSON.parse`.
- stringify
E: `JSON.stringify` converts an array or object into a string for storage; `JSON.parse` turns that string back into real data when you load it.
:::

## Recap

- The **state → render** pattern keeps data and UI in sync: events change state, then a render function rebuilds the DOM from that state.
- Keep a **single source of truth** (a variable, array, or object) and never edit the DOM directly from handlers.
- `localStorage` persists data between visits but stores only **strings** — use `JSON.stringify`/`JSON.parse` for objects and arrays.
- The **character counter** shows render-on-input plus a startup render call.
- The **to-do list** uses an array of objects and a full re-render — simple, predictable, and the foundation frameworks build on.
- **Tabs** are a tiny state machine: store the active tab, render with `classList.toggle(name, condition)`, and use event delegation.

**Next up:** organizing larger projects — modules, fetching data from APIs, and rendering dynamic content from the network.

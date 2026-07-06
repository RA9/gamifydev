# JavaScript State and Local Storage

If you want to build interfaces that feel like products instead of demos, you need a clean way to answer two questions:

1. **What is the current state of the UI?**
2. **What should survive a page refresh?**

This lesson covers both.

## What state means

**State** is the data that determines what the user currently sees.

Examples:

- which tab is active
- whether dark mode is enabled
- what tasks exist in a board
- which filter is selected
- what the user typed into a form

If the state changes, the UI should change.

:::key
A reliable frontend app has a **single source of truth** for each part of the interface.
:::

That source of truth should usually be JavaScript data, not random bits of DOM state spread everywhere.

## A simple state object

Instead of storing important values in five different variables, group related state together.

```js
const state = {
  filter: "all",
  theme: "light",
  tasks: [],
};
```

That makes your app easier to inspect, update, and save.

## Render from state

A strong pattern is:

1. state changes
2. render runs
3. DOM updates to match state

```js
const status = document.querySelector("#status");

const state = {
  count: 0,
};

function render() {
  status.textContent = `Count: ${state.count}`;
}

function increment() {
  state.count += 1;
  render();
}
```

This is much safer than changing the DOM from lots of different places with no shared model.

## Derived state

Not every value needs to be stored directly.

Some values can be **derived** from existing state.

```js
const state = {
  tasks: [
    { title: "Ship portfolio", done: false },
    { title: "Fix bug", done: true },
  ],
};

function getOpenCount() {
  return state.tasks.filter((task) => !task.done).length;
}
```

That count does not need its own variable because it can be calculated from the real source of truth.

:::tip
Store the minimum amount of state you need. Derive the rest when possible.
:::

## UI state vs persisted state

Not all state should live forever.

A useful distinction is:

- **UI state** — temporary state for the current session
- **persisted state** — state you want to keep after a refresh

Examples:

### UI state

- whether a modal is open
- which accordion section is expanded
- whether a button is in a loading state

### Persisted state

- theme preference
- saved tasks
- last selected dashboard layout
- recently viewed usernames

This distinction helps you decide what belongs in Local Storage.

## What Local Storage is

`localStorage` is a tiny browser key-value store.

It lets you save strings by key, and the browser keeps them between page reloads.

```js
localStorage.setItem("theme", "dark");
const theme = localStorage.getItem("theme");
console.log(theme); // "dark"
```

You can also remove one value or clear everything.

```js
localStorage.removeItem("theme");
localStorage.clear();
```

## Local Storage only stores strings

This is the most important rule.

If you want to save arrays or objects, convert them to JSON first.

```js
const tasks = [
  { title: "Write docs", done: false },
  { title: "Ship release", done: true },
];

localStorage.setItem("tasks", JSON.stringify(tasks));

const saved = localStorage.getItem("tasks");
const parsedTasks = saved ? JSON.parse(saved) : [];
```

Without `JSON.stringify`, complex data becomes useless text like `[object Object]`.

## A practical save/load pattern

This is a solid starter structure for many small apps.

```js
const state = {
  theme: "light",
  tasks: [],
};

function saveState() {
  localStorage.setItem("app-state", JSON.stringify(state));
}

function loadState() {
  const saved = localStorage.getItem("app-state");
  if (!saved) return;

  const parsed = JSON.parse(saved);
  state.theme = parsed.theme || "light";
  state.tasks = Array.isArray(parsed.tasks) ? parsed.tasks : [];
}
```

Then your startup flow becomes:

1. load saved state
2. render the UI
3. save after meaningful changes

## Be defensive when loading

Stored data can be missing, outdated, or malformed.

For that reason, avoid assuming saved values are always valid.

```js
function loadTheme() {
  const savedTheme = localStorage.getItem("theme");
  return savedTheme === "dark" ? "dark" : "light";
}
```

For larger objects, a `try/catch` is useful.

```js
function loadTasks() {
  try {
    const saved = localStorage.getItem("tasks");
    return saved ? JSON.parse(saved) : [];
  } catch (error) {
    return [];
  }
}
```

That prevents a broken saved value from breaking the whole app.

## Example: theme preference

```js
const state = {
  theme: localStorage.getItem("theme") === "dark" ? "dark" : "light",
};

function renderTheme() {
  document.body.classList.toggle("dark", state.theme === "dark");
}

document.querySelector("#theme-btn").addEventListener("click", () => {
  state.theme = state.theme === "dark" ? "light" : "dark";
  localStorage.setItem("theme", state.theme);
  renderTheme();
});

renderTheme();
```

This is small, but it teaches a professional habit:

- state changes in one place
- persistence happens deliberately
- render updates the UI

## Example: saved filters

Imagine a product page with `all`, `design`, and `dev` filters.

```js
const state = {
  filter: localStorage.getItem("filter") || "all",
};

function setFilter(nextFilter) {
  state.filter = nextFilter;
  localStorage.setItem("filter", nextFilter);
  render();
}
```

Saving little preferences like this makes apps feel more thoughtful.

## What should not go into Local Storage

Do **not** use Local Storage for:

- secrets or API keys
- sensitive user data
- large datasets
- important data that must sync across devices

Local Storage is best for lightweight browser-side preferences and small app state.

:::warning
Anything in Local Storage is visible to JavaScript running on the page. Treat it as convenient storage, not secure storage.
:::

## A clean app flow

For small frontend apps, this structure works well:

```js
const state = loadInitialState();

function render() {
  renderTasks();
  renderFilters();
  renderCounts();
}

function updateTasks(nextTasks) {
  state.tasks = nextTasks;
  saveState();
  render();
}
```

That pattern keeps behavior predictable:

- update state
- save if needed
- render

## Common mistakes

### 1. Saving on every tiny keystroke without a reason

Sometimes that is fine. Often it is noise.

Save at meaningful checkpoints unless the feature clearly needs draft persistence.

### 2. Mixing DOM state and JavaScript state randomly

If the DOM says one thing and your objects say another, debugging gets painful.

Pick a source of truth and stick to it.

### 3. Forgetting to render after state changes

If the data changed but the UI did not, the first question should be: did `render()` run?

## Why this matters for expert-level frontend work

Larger frameworks formalize these same ideas:

- component state
- derived data
- persistence
- one-way rendering

Learning this clearly in vanilla JavaScript gives you stronger instincts later, even when tools change.

:::quiz
Q: Why is `JSON.stringify` needed before saving an array to Local Storage?
- Because Local Storage only stores strings *
- Because arrays cannot be reused in JavaScript
- Because `filter()` requires it
E: Local Storage is string-based. `JSON.stringify` turns arrays and objects into storable text, and `JSON.parse` turns that text back into data.
:::

:::quiz
Q: Which is the better rule for small apps?
- Update the DOM from many handlers and store data wherever it is convenient
- Keep a clear state object, render from it, and persist only what should survive refreshes *
E: A state-first structure is easier to debug, easier to extend, and less likely to drift into inconsistent UI behavior.
:::

## Recap

- State is the data that determines what the UI shows.
- A single source of truth makes interfaces more predictable.
- Some values should be derived instead of stored separately.
- Local Storage is good for small persisted browser-side data like themes, tasks, and preferences.
- Save complex data with `JSON.stringify` and read it back with `JSON.parse`.
- The clean pattern is: update state, save if needed, then render.

**Next up:** you’ll apply this directly in richer projects where state, forms, and persistence all work together.
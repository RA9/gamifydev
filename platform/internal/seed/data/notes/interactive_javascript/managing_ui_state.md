# Managing UI State

As your features grow beyond a single button click, you need a clear mental model for how data flows through your application. This lesson formalizes the state → render pattern and shows you how to design state objects that scale.

## What state really is

**State** is the data that determines what the user sees right now. If the state changes, the UI should change. If the UI looks wrong, the first question is always: what does the state say?

```js
// This state fully describes a task board
const state = {
  filter: "all",       // "all", "active", "done"
  tasks: [
    { id: 1, text: "Ship portfolio", done: false },
    { id: 2, text: "Fix nav bug", done: true },
  ],
};
```

Every piece of interactive behavior — filters, toggles, form values, open/closed panels — should trace back to a state variable.

## Single source of truth

A **single source of truth** means each piece of information exists in exactly one place. The DOM is not that place.

```js
// Bad — state is scattered
// Is the task done? Check the class on the <li>.
// What is the filter? Check which button has .active.
// How many open tasks? Count the <li> elements without .done.

// Good — state is centralized
const state = { filter: "all", tasks: [...] };
// The DOM is rebuilt from state every time something changes.
```

:::key
Your JavaScript state is the truth. The DOM is a *picture* of that truth. If you ever need to know what the app is doing, read the state — not the DOM.
:::

## The update → render cycle

The core loop of any interactive frontend:

1. **Something happens** (click, input, timer, fetch response).
2. **Update the state** with the new data.
3. **Call render** to rebuild the UI from state.

```js
const state = { count: 0 };
const countEl = document.querySelector("#count");

function render() {
  countEl.textContent = `Count: ${state.count}`;
}

document.querySelector("#increment").addEventListener("click", () => {
  state.count += 1; // update state
  render();         // re-render
});

render(); // initial render
```

:::tip
Always call `render()` once at the bottom of your script so the UI starts in a correct state. Without this, the page may show stale HTML until the user's first interaction.
:::

## Derived state — compute, don't store

Some values can be **calculated** from existing state. When that is the case, do not store them separately — derive them in your render function or in a helper:

```js
const state = {
  tasks: [
    { text: "Ship it", done: false },
    { text: "Test it", done: true },
    { text: "Plan it", done: false },
  ],
};

// Derived — not stored in state
function getActiveCount() {
  return state.tasks.filter(t => !t.done).length;
}

function getDoneTasks() {
  return state.tasks.filter(t => t.done);
}

function render() {
  document.querySelector("#active-count").textContent = getActiveCount();
  // ...
}
```

Storing derived values creates a risk: the stored copy and the real data can drift apart. Computing on the fly is simpler and always correct.

:::key
Store the minimum state you need. Derive the rest. If you can compute a value from existing state, do not add another property for it.
:::

## UI state vs persisted state

Not all state should survive a page refresh. A useful distinction:

**UI state** — temporary, for the current session:
- Whether a dropdown is open
- Which step of a multi-step form the user is on
- Whether a tooltip is showing
- Loading and error flags

**Persisted state** — should survive a refresh (stored in `localStorage`, a database, or a URL):
- Theme preference
- Saved tasks or notes
- User's selected filters
- Recently viewed items

```js
// UI-only state — reset every page load
let isMenuOpen = false;
let isLoading = false;

// Persisted state — load from storage on startup
const tasks = JSON.parse(localStorage.getItem("tasks") || "[]");
const theme = localStorage.getItem("theme") || "light";
```

Keeping these separate prevents you from accidentally saving UI noise to storage or forgetting to persist important data.

## State shape design

How you structure your state object matters. Here are common patterns:

### Pattern 1: Simple toggles

```js
const state = {
  isDarkMode: false,
  isSidebarOpen: true,
};
```

### Pattern 2: Active selection from a set

```js
const state = {
  activeTab: "home", // one of: "home", "profile", "settings"
};
```

### Pattern 3: List with selection

```js
const state = {
  items: [
    { id: 1, text: "Item A" },
    { id: 2, text: "Item B" },
  ],
  selectedId: null, // id of the selected item, or null
};
```

### Pattern 4: Filtered list

```js
const state = {
  searchQuery: "",
  activeCategory: "all",
  items: [...],
};

// Derived: the filtered list used for rendering
function getFilteredItems() {
  return state.items
    .filter(item => state.activeCategory === "all" || item.category === state.activeCategory)
    .filter(item => item.title.toLowerCase().includes(state.searchQuery.toLowerCase()));
}
```

## Separating state updates from rendering

As your app grows, isolate state changes into small functions:

```js
function addTask(text) {
  state.tasks.push({ id: Date.now(), text, done: false });
  render();
}

function toggleTask(id) {
  const task = state.tasks.find(t => t.id === id);
  if (task) task.done = !task.done;
  render();
}

function setFilter(filter) {
  state.filter = filter;
  render();
}
```

Each function does one thing: update state, then render.

## A complete example — task manager

```js
const state = {
  filter: "all",
  tasks: JSON.parse(localStorage.getItem("tasks") || "[]"),
};

function save() {
  localStorage.setItem("tasks", JSON.stringify(state.tasks));
}

function getVisibleTasks() {
  if (state.filter === "active") return state.tasks.filter(t => !t.done);
  if (state.filter === "done") return state.tasks.filter(t => t.done);
  return state.tasks;
}

function render() {
  const list = document.querySelector("#list");
  const visible = getVisibleTasks();

  list.innerHTML = visible.map(t =>
    `<li data-id="${t.id}" class="${t.done ? "done" : ""}">
       <span>${t.text}</span>
       <button class="toggle-btn">${t.done ? "Undo" : "Done"}</button>
       <button class="delete-btn">Delete</button>
     </li>`
  ).join("");

  document.querySelector("#count").textContent =
    `${state.tasks.filter(t => !t.done).length} tasks remaining`;
}

// Delegated event listener
document.querySelector("#list").addEventListener("click", (e) => {
  const li = e.target.closest("li");
  if (!li) return;
  const id = Number(li.dataset.id);

  if (e.target.closest(".toggle-btn")) {
    const task = state.tasks.find(t => t.id === id);
    if (task) task.done = !task.done;
    save();
    render();
  }

  if (e.target.closest(".delete-btn")) {
    state.tasks = state.tasks.filter(t => t.id !== id);
    save();
    render();
  }
});

render();
```

## Practice

:::quiz
Q: Your app has a list of 20 tasks and a count showing "5 remaining." Where should the "5" come from?
- A separate `remainingCount` variable in state
- A DOM query counting `.active` elements
- Calculated from the tasks array each time render runs *
- Stored in localStorage
E: The remaining count is derived state — it can be computed from the existing tasks array. Storing it separately risks the count and the array going out of sync.
:::

:::quiz
Q: Which of these should be saved to localStorage?
- Whether a dropdown is currently open
- The user's selected theme preference *
- Whether a loading spinner is showing
- The current value of event.target
E: Theme preference is a meaningful user choice that should persist across sessions. UI-only state like dropdown visibility, loading flags, and event data are temporary and should reset on page load.
:::

## Recap

- **State** is JavaScript data that determines the UI. The DOM is a picture of that data.
- Maintain a **single source of truth** — never rely on the DOM for your app's current state.
- The cycle: **event → update state → render**. Render reads state and rebuilds the DOM.
- **Derived state** is computed from existing state — do not store it separately.
- Separate **UI state** (temporary) from **persisted state** (saved to localStorage or a server).
- Design your **state shape** around your feature: toggles, active selections, lists with filters.
- Isolate state updates into small functions that call `render()` at the end.

**Next up:** Theme Toggles and Preferences — applying state management to a real feature.

# Search and Filter UIs

Search and filter interfaces let users explore data interactively — typing a query, clicking category buttons, sorting results. This lesson builds a complete filter system with live search, category filtering, sort controls, and an empty state.

## The filter state

One object tracks every aspect of how the data is being filtered:

```js
const state = {
  searchQuery: "",
  activeCategory: "all",
  sortBy: "name",        // "name" or "date"
  items: [
    { id: 1, name: "Portfolio site", category: "design", date: "2025-03-15" },
    { id: 2, name: "REST API", category: "dev", date: "2025-06-01" },
    { id: 3, name: "Blog theme", category: "design", date: "2025-01-20" },
    { id: 4, name: "CLI tool", category: "dev", date: "2025-07-10" },
    { id: 5, name: "Brand guide", category: "design", date: "2025-04-05" },
  ],
};
```

## Live search with the input event

```html
<input type="search" id="search" placeholder="Search projects..." />
```

```js
const searchInput = document.querySelector("#search");

searchInput.addEventListener("input", (event) => {
  state.searchQuery = event.target.value;
  render();
});
```

The `input` event fires on every keystroke, giving instant feedback. The render function uses the query to filter the displayed items.

## Category buttons with active state

```html
<div id="filters">
  <button class="filter-btn active" data-category="all">All</button>
  <button class="filter-btn" data-category="design">Design</button>
  <button class="filter-btn" data-category="dev">Dev</button>
</div>
```

```js
const filterBar = document.querySelector("#filters");

filterBar.addEventListener("click", (event) => {
  const btn = event.target.closest(".filter-btn");
  if (!btn) return;

  state.activeCategory = btn.dataset.category;
  render();
});
```

The render function handles updating which button looks active — no need to toggle classes in the event handler.

## Chaining filter + search

The key insight: apply all filters in sequence to produce the visible list. Each filter narrows the results further.

```js
function getFilteredItems() {
  let items = state.items;

  // 1. Category filter
  if (state.activeCategory !== "all") {
    items = items.filter(item => item.category === state.activeCategory);
  }

  // 2. Search filter
  if (state.searchQuery.trim()) {
    const query = state.searchQuery.toLowerCase();
    items = items.filter(item =>
      item.name.toLowerCase().includes(query)
    );
  }

  // 3. Sort
  items = [...items].sort((a, b) => {
    if (state.sortBy === "name") return a.name.localeCompare(b.name);
    if (state.sortBy === "date") return new Date(b.date) - new Date(a.date);
    return 0;
  });

  return items;
}
```

:::key
Filter functions should never modify the original array. Always work on a copy (using `.filter()` which returns a new array) and spread before `.sort()` (which sorts in place). The source data stays intact.
:::

## Sort controls

```html
<select id="sort">
  <option value="name">Sort by Name</option>
  <option value="date">Sort by Date</option>
</select>
```

```js
document.querySelector("#sort").addEventListener("change", (event) => {
  state.sortBy = event.target.value;
  render();
});
```

## The render function

```js
const list = document.querySelector("#results");
const emptyState = document.querySelector("#empty-state");
const resultCount = document.querySelector("#result-count");

function render() {
  const filtered = getFilteredItems();

  // Update result count
  resultCount.textContent = `${filtered.length} project${filtered.length !== 1 ? "s" : ""}`;

  // Update category button active states
  document.querySelectorAll(".filter-btn").forEach(btn => {
    btn.classList.toggle("active", btn.dataset.category === state.activeCategory);
  });

  // Show empty state or results
  if (filtered.length === 0) {
    list.innerHTML = "";
    emptyState.hidden = false;
    emptyState.textContent = state.searchQuery
      ? `No results for "${state.searchQuery}"`
      : "No items in this category";
    return;
  }

  emptyState.hidden = true;
  list.innerHTML = filtered.map(item =>
    `<div class="card" data-id="${item.id}">
       <h3>${item.name}</h3>
       <span class="badge">${item.category}</span>
       <time>${item.date}</time>
     </div>`
  ).join("");
}
```

## Empty state — never leave the user guessing

When filters produce zero results, show a helpful message instead of a blank screen:

```html
<p id="empty-state" hidden></p>
```

:::tip
A good empty state tells the user *why* there are no results and *what to do about it*: "No results for 'xyz' — try a different search" or "No items in Design — create your first project." A blank container is confusing.
:::

## URL query params with URLSearchParams

Persist filter state in the URL so users can share filtered views and use the back button:

```js
function syncStateToURL() {
  const params = new URLSearchParams();
  if (state.searchQuery) params.set("q", state.searchQuery);
  if (state.activeCategory !== "all") params.set("category", state.activeCategory);
  if (state.sortBy !== "name") params.set("sort", state.sortBy);

  const newURL = params.toString()
    ? `${window.location.pathname}?${params}`
    : window.location.pathname;

  history.replaceState(null, "", newURL);
}

function loadStateFromURL() {
  const params = new URLSearchParams(window.location.search);
  state.searchQuery = params.get("q") || "";
  state.activeCategory = params.get("category") || "all";
  state.sortBy = params.get("sort") || "name";

  // Sync the search input
  searchInput.value = state.searchQuery;
}
```

Call `syncStateToURL()` at the end of render, and `loadStateFromURL()` on startup:

```js
loadStateFromURL();
render();
```

:::warning
Use `history.replaceState` (not `pushState`) for filter changes, unless you specifically want every filter tweak to create a new browser history entry. Too many push states makes the back button unusable.
:::

## Debouncing the search

For large datasets or API-backed search, debounce the input handler:

```js
function debounce(fn, delay) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

const handleSearch = debounce((query) => {
  state.searchQuery = query;
  render();
}, 250);

searchInput.addEventListener("input", (e) => handleSearch(e.target.value));
```

## Highlighting search matches

A nice touch: highlight the matching text in results.

```js
function highlightMatch(text, query) {
  if (!query) return text;
  const regex = new RegExp(`(${query})`, "gi");
  return text.replace(regex, "<mark>$1</mark>");
}

// In your render:
`<h3>${highlightMatch(item.name, state.searchQuery)}</h3>`
```

:::warning
This uses `innerHTML`, so only use it with data you control (your own items array), never with raw user-submitted content.
:::

## Practice

:::quiz
Q: The user selects the "Design" category and types "blog" in the search box. How should `getFilteredItems` apply these?
- Apply only the most recent filter
- Apply category OR search, showing results that match either
- Apply category AND search, narrowing results to Design items containing "blog" *
- Ignore the search if a category is selected
E: Filters should chain — each one narrows the results further. First filter by category (Design items only), then filter those results by the search query ("blog"). The user sees only Design items whose name contains "blog."
:::

:::quiz
Q: The filtered result list is empty. What should the UI show?
- Nothing — just an empty container
- A message explaining why there are no results and what the user can try *
- The unfiltered list
- An error message
E: An empty state should tell the user why there are no results (their search or filter produced nothing) and suggest an action (try a different query, reset filters). A blank screen leaves users confused about whether something is broken.
:::

## Recap

- Store filter state in **one object**: `searchQuery`, `activeCategory`, `sortBy`.
- Use the **`input`** event for live search and **`click`** delegation for category buttons.
- **Chain filters** in a `getFilteredItems()` function — each narrows the results further.
- Sort with `.sort()` on a **copy** of the array (spread first to avoid mutating the source).
- Render an **empty state** message when no items match — never leave the screen blank.
- Sync state to **URL query params** with `URLSearchParams` and `history.replaceState`.
- **Debounce** expensive search operations to avoid running on every keystroke.

**Next up:** Building Reusable UI Components — creating functions that return DOM, the step before frameworks.

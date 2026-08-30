# Loading Error and Empty States

Every time you fetch data, the UI goes through at least three possible states: loading (waiting), success (data arrived), or error (something failed). Most beginner code only handles success. Professional code handles all three — plus a fourth: the data arrived, but it is empty.

## The three states of every fetch

```
User clicks "Load" → Loading → Success (render data)
                              → Error (show message)
                              → Empty (show "no results")
```

Your state object should model this explicitly:

```js
const state = {
  items: [],
  isLoading: false,
  error: null,
};
```

:::key
Every fetch-driven UI has at least three states: **loading**, **error**, and **success**. If you only code the success path, your users see a blank screen or a frozen spinner when things go wrong.
:::

## Showing a loading state

Set the loading flag *before* the fetch and clear it when the fetch completes:

```js
async function loadItems() {
  state.isLoading = true;
  state.error = null;
  render();

  try {
    const response = await fetch("/api/items");
    if (!response.ok) throw new Error(`Status: ${response.status}`);
    state.items = await response.json();
  } catch (error) {
    state.error = "Unable to load items. Please try again.";
    console.error(error);
  } finally {
    state.isLoading = false;
    render();
  }
}
```

The `finally` block guarantees the loading flag is cleared whether the fetch succeeds or fails.

### Rendering the loading state

```js
function render() {
  const container = document.querySelector("#content");

  if (state.isLoading) {
    container.innerHTML = `
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Loading items...</p>
      </div>
    `;
    return;
  }

  // ... error and success rendering below
}
```

A simple CSS spinner:

```css
.spinner {
  width: 24px;
  height: 24px;
  border: 3px solid #e0e0e0;
  border-top-color: teal;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
```

:::tip
A spinner or skeleton screen should appear **immediately** when the fetch starts. The user should never see a blank screen and wonder if anything is happening.
:::

## Clearing old data before a new fetch

When refetching (changing a filter, searching again), clear the old data so stale content does not linger:

```js
async function loadItems(category) {
  state.items = [];       // clear old data
  state.isLoading = true;
  state.error = null;
  render();               // shows loading state immediately

  try {
    const response = await fetch(`/api/items?category=${category}`);
    if (!response.ok) throw new Error(`Status: ${response.status}`);
    state.items = await response.json();
  } catch (error) {
    state.error = "Failed to load items.";
  } finally {
    state.isLoading = false;
    render();
  }
}
```

:::warning
If you do not clear old data before rendering the loading state, the user sees outdated results from the previous query mixed with a loading indicator. Clear the array and re-render before starting the fetch.
:::

## Rendering error states

An error message should tell the user *what happened* and *what they can do about it*:

```js
if (state.error) {
  container.innerHTML = `
    <div class="error-state">
      <p class="error-message">${state.error}</p>
      <button id="retry-btn">Try Again</button>
    </div>
  `;

  document.querySelector("#retry-btn").addEventListener("click", loadItems);
  return;
}
```

:::tip
Error messages in the UI should answer two questions: **What happened?** and **What can I do?** Always include a retry button, a suggestion, or a link to help.
:::

## Empty state — data loaded, but there is nothing

This is the fourth state many developers forget. The fetch succeeded, but the array is empty:

```js
if (state.items.length === 0) {
  container.innerHTML = `
    <div class="empty-state">
      <p>No projects yet.</p>
      <p>Create your first project to get started.</p>
      <button id="create-btn">New Project</button>
    </div>
  `;

  document.querySelector("#create-btn").addEventListener("click", openNewProjectForm);
  return;
}
```

### Context-aware empty states

Different reasons for emptiness need different messages:

```js
function getEmptyMessage() {
  if (state.searchQuery) {
    return `No results for "${state.searchQuery}". Try a different search.`;
  }
  if (state.activeFilter !== "all") {
    return `No ${state.activeFilter} items. Try changing the filter.`;
  }
  return "Nothing here yet. Create your first item to get started.";
}
```

## The complete render function

Putting all four states together:

```js
function render() {
  const container = document.querySelector("#content");

  // 1. Loading
  if (state.isLoading) {
    container.innerHTML = renderLoadingSkeleton();
    return;
  }

  // 2. Error
  if (state.error) {
    container.innerHTML = `
      <div class="error-state">
        <p>${state.error}</p>
        <button class="retry-btn">Try Again</button>
      </div>
    `;
    container.querySelector(".retry-btn").addEventListener("click", () => loadItems());
    return;
  }

  // 3. Empty
  if (state.items.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <p>${getEmptyMessage()}</p>
      </div>
    `;
    return;
  }

  // 4. Success — render the data
  container.innerHTML = state.items
    .map(item => `
      <div class="card">
        <h3>${item.name}</h3>
        <p>${item.description}</p>
      </div>
    `)
    .join("");
}
```

:::key
Render in priority order: **loading → error → empty → success**. Each state returns early so only one state renders at a time. This structure makes it impossible to show a spinner and data at the same time, or an error with stale results.
:::

## Practice

:::quiz
Q: A fetch returns an empty array `[]`. What should the UI show?
- A loading spinner
- An error message
- An empty state message explaining there are no items and suggesting an action *
- Nothing — leave the screen blank
E: An empty array is a successful response, not an error. Show a helpful empty state message that tells the user why there is nothing (their filters produced no results, they have not created anything yet) and guides them on what to do next.
:::

:::quiz
Q: You are refetching data after the user changes a filter. What should you do before starting the new fetch?
- Nothing — let the old data stay until the new data arrives
- Clear the old data and show a loading state immediately *
- Show an error message
- Disable the entire page
E: Clearing old data and showing a loading state prevents the user from seeing stale results from the previous filter while the new results load. Set `items = []`, `isLoading = true`, and call render before starting the fetch.
:::

## Recap

- Every fetch UI has at minimum **three states**: loading, error, and success. Add **empty** as a fourth.
- Show a **spinner or skeleton** immediately when the fetch starts — never leave the screen blank.
- **Clear old data** before refetching to prevent stale results from lingering.
- Error messages should tell the user **what happened** and **what to do** — include a retry button.
- Empty states should be **context-aware** — different messages for empty searches, empty filters, and genuinely empty data.
- Render in priority order: **loading → error → empty → success**. Each returns early.
- **Disable buttons** during loading to prevent duplicate requests.

**You now have the full async toolkit:** promises, async/await, error handling, fetch, HTTP, CORS, and UI state management for every phase of a network request.

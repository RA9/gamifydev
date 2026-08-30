# Project: Build a To-Do App

The to-do app is one of the best early frontend projects because it teaches you how to connect state, rendering, events, and persistence into one believable product.

:::project
**Goal:** Build a task manager where users can add tasks, mark them complete, delete them, and keep them saved between page refreshes with localStorage.
:::

## What the app needs

Your finished version should support:

- adding a task from a form
- rendering the current task list
- toggling a task complete/incomplete
- deleting a task
- saving the list in localStorage
- showing a helpful empty state when there are no tasks

## Step 1 — Build the semantic HTML shell

```html
<main class="todo-shell">
  <section class="todo-card">
    <h1>Task Tracker</h1>

    <form id="task-form">
      <label for="task-input">Add a task</label>
      <div class="task-row">
        <input id="task-input" type="text" placeholder="What needs doing?" />
        <button type="submit">Add</button>
      </div>
    </form>

    <ul id="task-list"></ul>
    <p id="empty-state">No tasks yet. Add your first one.</p>
  </section>
</main>
```

That gives you a real, accessible form and a destination for the rendered task list.

## Step 2 — Style the layout

Keep the visuals simple and trustworthy.

```css
body {
  margin: 0;
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f8fafc;
  font-family: system-ui, sans-serif;
}

.todo-card {
  width: min(100% - 2rem, 34rem);
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 18px;
  padding: 2rem;
  box-shadow: 0 14px 36px rgba(15, 23, 42, 0.08);
}

.task-row {
  display: flex;
  gap: 0.75rem;
}

#task-list {
  list-style: none;
  padding: 0;
  margin: 1.25rem 0 0;
}
```

You can add the rest of the polish as the app logic lands.

## Step 3 — Model the data

Use an array of objects.

```js
let tasks = [];
```

Each task should look like:

```js
{
  id: 1,
  text: "Ship landing page",
  done: false
}
```

This is much stronger than storing plain strings, because it gives each task its own identity and completion state.

## Step 4 — Create a render function

The whole app becomes easier when one function redraws the UI from the task array.

```js
function renderTasks() {
  taskList.innerHTML = "";

  if (tasks.length === 0) {
    emptyState.hidden = false;
    return;
  }

  emptyState.hidden = true;

  tasks.forEach(function (task) {
    const item = document.createElement("li");
    item.textContent = task.text;
    taskList.appendChild(item);
  });
}
```

You'll grow this function to include complete/delete controls.

## Step 5 — Add task creation

On form submit:

- prevent default behavior
- read the input value
- ignore empty text
- push a new task object into the array
- save
- render
- clear the input

This is your first full event → state → render loop in the app.

## Step 6 — Add complete and delete actions

Each rendered task should include:

- a clickable complete toggle
- a delete button

A simple pattern is:

- click task row or checkbox → flip `done`
- click delete button → remove that task by `id`
- save and re-render after both

Use `stopPropagation()` if the delete button sits inside a clickable row.

## Step 7 — Persist with localStorage

To keep tasks across refreshes:

```js
function saveTasks() {
  localStorage.setItem("tasks", JSON.stringify(tasks));
}

function loadTasks() {
  const saved = localStorage.getItem("tasks");
  tasks = saved ? JSON.parse(saved) : [];
}
```

Call `loadTasks()` when the app starts, then `renderTasks()`.

:::key
This is the real magic of the project: the UI is driven by your `tasks` array, and localStorage preserves that state between sessions.
:::

## Helpful UX polish

A stronger final version also includes:

- a visual completed state (line-through or muted text)
- a counter like `3 tasks left`
- button hover states
- clear focus styles for keyboard users

## Stretch goals

- filter all / active / completed
- add due dates or priorities
- allow editing a task
- save the chosen filter in localStorage

:::quiz
Q: What makes a to-do app a strong frontend project?
- It is only about colors and layout
- It combines state, rendering, events, and persistence *
- It avoids using arrays or objects
E: A to-do app teaches real frontend architecture because the UI depends on changing data and stored browser state.
:::

## What "done" looks like

A finished version should feel like a small real app: reliable, understandable, and able to recover its state after a refresh.
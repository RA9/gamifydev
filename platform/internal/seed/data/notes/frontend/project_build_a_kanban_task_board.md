# Project: Build a Kanban Task Board

A to-do list teaches state. A Kanban board teaches **state with structure**.

Instead of one flat list, your tasks move through stages like a real product workflow.

:::project
**Goal:** Build a task board with columns such as **Backlog**, **In Progress**, and **Done**. Users should be able to add tasks, move them between columns, persist the board in Local Storage, and clearly see the current status of work.
:::

## Why this is a strong frontend project

This project pushes you beyond beginner CRUD.

You will practise:

- modeling richer state
- rendering multiple derived views from one data source
- handling actions like add, move, edit, and delete
- persisting changes locally
- designing a UI that stays readable as the data grows

That is much closer to the kind of interface work real teams build.

## What the finished board should include

A strong version should support:

- at least 3 columns
- adding a new task
- moving tasks between columns
- deleting tasks
- optional tags, priority, or assignee text
- Local Storage persistence
- responsive layout
- empty states for columns with no tasks

## Step 1 — Build the board shell

Start with a semantic layout.

```html
<main class="board-shell">
  <header class="board-head">
    <div>
      <h1>Project Board</h1>
      <p>Track work from idea to done.</p>
    </div>

    <form id="task-form" class="task-form">
      <label class="sr-only" for="task-title">Task title</label>
      <input id="task-title" type="text" placeholder="Add a new task" />
      <select id="task-status">
        <option value="backlog">Backlog</option>
        <option value="in-progress">In Progress</option>
        <option value="done">Done</option>
      </select>
      <button type="submit">Add task</button>
    </form>
  </header>

  <section class="board-grid">
    <article class="column" data-status="backlog">
      <h2>Backlog</h2>
      <ul id="backlog-list"></ul>
      <p class="empty-state">No backlog tasks yet.</p>
    </article>

    <article class="column" data-status="in-progress">
      <h2>In Progress</h2>
      <ul id="in-progress-list"></ul>
      <p class="empty-state">Nothing in progress right now.</p>
    </article>

    <article class="column" data-status="done">
      <h2>Done</h2>
      <ul id="done-list"></ul>
      <p class="empty-state">No completed tasks yet.</p>
    </article>
  </section>
</main>
```

This gives you a clean destination for derived task lists.

## Step 2 — Style the board like a product surface

Use Grid so the layout can expand and collapse cleanly.

```css
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  background: #f8fafc;
  color: #0f172a;
}

.board-shell {
  width: min(100% - 2rem, 78rem);
  margin: 0 auto;
  padding: 2rem 0 4rem;
}

.board-grid {
  display: grid;
  gap: 1rem;
}

.column {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 18px;
  padding: 1rem;
  min-height: 20rem;
}

@media (min-width: 900px) {
  .board-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
```

Give each card enough spacing so the board remains scannable.

## Step 3 — Model the state

Use one array of task objects.

```js
let tasks = [];
```

Each task can look like this:

```js
{
  id: 1,
  title: "Refine signup validation",
  status: "in-progress",
  priority: "medium"
}
```

The important part is that **status lives in the data**.

That lets one array drive all three columns.

## Step 4 — Render each column from the same source of truth

Instead of keeping separate arrays for each column, derive them.

```js
function getTasksByStatus(status) {
  return tasks.filter((task) => task.status === status);
}
```

Then render each column by filtering the main array.

```js
function renderColumn(status, listElement) {
  const visibleTasks = getTasksByStatus(status);
  listElement.innerHTML = "";

  visibleTasks.forEach((task) => {
    const item = document.createElement("li");
    item.textContent = task.title;
    listElement.appendChild(item);
  });
}
```

This is a core frontend architecture habit: one source of truth, many derived views.

## Step 5 — Add new tasks

On form submit:

- prevent default
- read the title and chosen status
- ignore empty input
- push a new task object into `tasks`
- save
- render
- reset the form

This should feel similar to the to-do app, but with stronger data modeling.

## Step 6 — Add move actions

Each rendered task card should include actions like:

- **Move left**
- **Move right**
- **Delete**

A simple move rule might be:

- `backlog` → `in-progress`
- `in-progress` → `done`
- `done` has no move-right action

You can render action buttons based on the task’s current status.

:::tip
Make the move actions explicit first. Drag and drop is a great stretch goal, but clear buttons are easier to implement, easier to test, and usually more accessible.
:::

## Step 7 — Save the board with Local Storage

Persist the task array so the board survives a refresh.

```js
function saveTasks() {
  localStorage.setItem("kanban-tasks", JSON.stringify(tasks));
}

function loadTasks() {
  const saved = localStorage.getItem("kanban-tasks");
  tasks = saved ? JSON.parse(saved) : [];
}
```

Run `loadTasks()` before the first render.

## Step 8 — Add empty states and counts

Each column should clearly communicate whether it has work.

Good polish includes:

- an empty message when a column has zero tasks
- a count beside the column title
- subtle status colors or badges

This is what makes the board feel intentional rather than just functional.

## Step 9 — Review the UX quality

Check for:

- keyboard-accessible buttons
- visible focus states
- readable spacing on mobile
- stable card layout with long task titles
- no disappearing data after refresh

## Stretch goals

Once the core board works, extend it with:

- drag and drop between columns
- priority labels
- edit task title inline
- saved filters like `all`, `high priority`, or `done only`
- due dates
- assignee avatars or initials

## What expert-level thinking looks like here

The strongest versions of this project do not just move cards around.

They show that you can:

- model data clearly
- derive UI from state cleanly
- keep interactions predictable
- persist data without breaking the flow

That is exactly the kind of thinking that carries into larger apps and frameworks.

:::quiz
Q: Why is one `tasks` array usually better than keeping separate arrays for each board column?
- Because separate arrays cannot be rendered
- Because one source of truth makes filtering, moving, and saving tasks simpler *
- Because Grid only works with one array
E: A single task list with a `status` field is easier to reason about and easier to persist, while each column can be derived with `filter()`.
:::

## What “done” looks like

A finished board should feel like a small workflow product:

- tasks can be created quickly
- movement between stages is clear
- the layout works on laptop and mobile
- data survives refreshes
- the interface stays understandable as the board grows
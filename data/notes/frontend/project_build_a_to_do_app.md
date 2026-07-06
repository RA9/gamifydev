# Project: Build a To-Do App

The to-do app is a classic frontend project because it forces you to combine almost everything that matters:

- semantic HTML
- layout and component styling
- DOM manipulation
- events
- state
- persistence

:::project
You'll build a task manager where users can add tasks, mark them complete, delete them, and keep them saved between page refreshes with `localStorage`.
:::

## Why this project matters

Unlike a static page, a to-do app has real product behavior:

- a list changes over time
- the UI depends on state
- actions have consequences
- the browser remembers data

That makes it one of the best early frontend projects for developing real engineering instincts.

## Define the product before coding

Your version should support these minimum behaviors:

- add a task
- show the task list
- toggle a task complete/incomplete
- delete a task
- persist tasks across refreshes

Optional upgrades later:

- filter by complete / incomplete
- edit task text
- clear completed tasks
- show counts (`3 tasks left`)

## Model the data first

A clean starter model is an array of objects:

```js
let tasks = [
  { id: 1, text: "Ship landing page", done: false },
  { id: 2, text: "Audit form labels", done: true }
];
```

Why objects instead of plain strings?

Because real tasks need more than text. They often need:

- a unique id
- a completion state
- maybe due dates or priority later

## Build the UI shell

A good HTML starting point:

```html
<section class="todo-app">
  <h1>Task Tracker</h1>

  <form id="taskForm">
    <label for="taskInput">Add a task</label>
    <input id="taskInput" name="task" type="text" required />
    <button type="submit">Add task</button>
  </form>

  <ul id="taskList"></ul>
</section>
```

That already gives you:

- a clear heading
- an accessible label
- a form submission path
- a destination for rendered tasks

## Use the render pattern

Instead of manually inserting fragments all over the place, create a single render function.

```js
function renderTasks() {
  taskList.innerHTML = "";

  tasks.forEach(function (task) {
    const item = document.createElement("li");
    item.textContent = task.text;
    taskList.appendChild(item);
  });
}
```

As the project grows, that render function can also:

- apply a `.done` class
- inject delete buttons
- update a task counter
- show empty-state messages

## Wire up the interactions

The app needs a few key event flows.

### Add a task

- user submits the form
- validate non-empty input
- push a new task object into the array
- save
- render

### Toggle a task

- user clicks a task or checkbox
- flip `done`
- save
- render

### Delete a task

- user clicks delete
- remove task from array
- save
- render

That loop — **update state → save → render** — is the heart of the app.

:::key
When you feel lost in a project like this, return to the loop: what changed in state, what should be persisted, and what should the UI look like now?
:::

## Persist with `localStorage`

To keep tasks after a refresh:

```js
function saveTasks() {
  localStorage.setItem("tasks", JSON.stringify(tasks));
}

function loadTasks() {
  const saved = localStorage.getItem("tasks");
  tasks = saved ? JSON.parse(saved) : [];
}
```

Then call `loadTasks()` when the page starts, and `saveTasks()` after every meaningful state update.

## UX details that make it feel real

A stronger to-do app also includes:

- clear empty state text when there are no tasks
- visible completion styling
- buttons with specific labels
- form reset after successful add
- no accidental duplicate whitespace-only tasks

Example empty state idea:

```js
if (tasks.length === 0) {
  taskList.innerHTML = "<li>No tasks yet. Add your first one.</li>";
  return;
}
```

## Acceptance criteria

Your project is complete when:

- tasks can be added through a form
- each task renders into the list
- tasks can be toggled complete
- tasks can be deleted
- the list is restored after page refresh
- the UI stays understandable when the list is empty

:::quiz
Q: Which project pattern best describes a to-do app?
- Static content only
- State-driven UI with persistent browser data *
- Styling without interaction
E: A to-do app works because the interface is driven by task state and that state is persisted in the browser.
:::

## Stretch goals

After the core version works, try:

- adding a checkbox instead of clicking the whole list item
- showing "tasks left" count
- saving filters in `localStorage`
- separating app logic into smaller helper functions

## If you want a full build-along

A dedicated step-by-step build-along version of this project also lives in the **Projects gallery**. Use this lesson as the architecture guide and the gallery version as the hands-on implementation path.

:::warning
Don't jump straight into advanced features before the core loop works. A small, reliable app teaches more than a large, broken one.
:::

## What's next

In **Git and GitHub**, you'll learn how to save, version, and share builds like this one the way real teams do.
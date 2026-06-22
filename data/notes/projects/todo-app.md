# To-Do List App

The to-do app is the "hello world" of real applications — and yours will remember your tasks even after you close the tab. You'll learn arrays, rendering, events, and browser storage all in one build.

:::project
You'll have a single HTML file where you can add tasks, click a task to mark it done (with a line-through), delete tasks, and have everything persist across page refreshes using localStorage.
:::

To follow along you only need a text editor and a web browser. Build everything in one file called `todo.html` with a `<style>` and `<script>` block, refreshing after each step.

## Step 1 — Build the HTML structure

Create the page with a text input, an Add button, and an empty `<ul>` that JavaScript will fill with tasks.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>To-Do List</title>
</head>
<body>
  <div class="app">
    <h1>My Tasks</h1>
    <div class="row">
      <input id="taskInput" type="text" placeholder="What needs doing?" />
      <button id="addBtn">Add</button>
    </div>
    <ul id="list"></ul>
  </div>
</body>
</html>
```

## Step 2 — Add basic styling

Add a `<style>` block so the app is centered and readable. We also prepare a `.done` class to cross out completed tasks and a delete button style for later.

```html
<style>
  body {
    margin: 0; min-height: 100vh;
    display: flex; justify-content: center; padding-top: 60px;
    background: #fafaf9; font-family: system-ui, sans-serif;
  }
  .app { width: 360px; }
  h1 { color: #1c1917; }
  .row { display: flex; gap: 8px; }
  #taskInput {
    flex: 1; padding: 10px; border: 1px solid #d6d3d1; border-radius: 8px; font-size: 15px;
  }
  #addBtn {
    padding: 10px 18px; border: none; border-radius: 8px;
    background: #f59e0b; color: #fff; font-weight: 600; cursor: pointer;
  }
  ul { list-style: none; padding: 0; margin-top: 16px; }
  li {
    display: flex; justify-content: space-between; align-items: center;
    padding: 12px; background: #fff; border: 1px solid #e7e5e4;
    border-radius: 8px; margin-bottom: 8px; cursor: pointer;
  }
  li.done span { text-decoration: line-through; color: #a8a29e; }
  .del {
    border: none; background: #ef4444; color: #fff;
    border-radius: 6px; padding: 4px 10px; cursor: pointer;
  }
</style>
```

## Step 3 — Add a task and render the list

Open a `<script>` block before `</body>`. Keep tasks in an array of objects (text plus a done flag). Write a `render` function that rebuilds the list, and wire the Add button to push a new task and re-render.

```html
<script>
  let tasks = [];

  const input = document.getElementById("taskInput");
  const addBtn = document.getElementById("addBtn");
  const list = document.getElementById("list");

  function render() {
    list.innerHTML = "";
    tasks.forEach((task, index) => {
      const li = document.createElement("li");
      if (task.done) li.classList.add("done");
      li.innerHTML = `<span>${task.text}</span>
                      <button class="del">x</button>`;
      list.appendChild(li);
    });
  }

  function addTask() {
    const text = input.value.trim();
    if (text === "") return;
    tasks.push({ text: text, done: false });
    input.value = "";
    render();
  }

  addBtn.onclick = addTask;
  render();
</script>
```

## Step 4 — Toggle a task complete

Add click handling so clicking a task flips its `done` flag. We attach the handler inside `render` so each list item knows its own index.

```js
function render() {
  list.innerHTML = "";
  tasks.forEach((task, index) => {
    const li = document.createElement("li");
    if (task.done) li.classList.add("done");
    li.innerHTML = `<span>${task.text}</span>
                    <button class="del">x</button>`;

    li.onclick = () => {
      tasks[index].done = !tasks[index].done;
      render();
    };

    list.appendChild(li);
  });
}
```

## Step 5 — Delete a task

Hook up the red delete button. Because the whole `<li>` also toggles, we call `stopPropagation()` so deleting doesn't accidentally toggle the task first.

```js
li.querySelector(".del").onclick = (event) => {
  event.stopPropagation();
  tasks.splice(index, 1);
  render();
};
```

:::tip
Put that block inside the `tasks.forEach` loop in `render`, right before `list.appendChild(li)`. It needs the `li` and `index` from the loop.
:::

## Step 6 — Save and load with localStorage

Finally, make tasks survive a refresh. Save the array as JSON after every change, and load it when the page opens. Add a `save()` call inside `addTask`, the toggle handler, and the delete handler.

```js
function save() {
  localStorage.setItem("tasks", JSON.stringify(tasks));
}

function load() {
  const stored = localStorage.getItem("tasks");
  if (stored) tasks = JSON.parse(stored);
}

// Call save() at the end of addTask(), after the toggle line,
// and after the splice in the delete handler. For example:
//   tasks.push({ text: text, done: false }); ... save();

load();
render();
```

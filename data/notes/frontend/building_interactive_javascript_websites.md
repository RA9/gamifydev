# Building Interactive JavaScript Websites

This is where JavaScript stops feeling like isolated syntax and starts feeling like product development.

Interactive websites are built from a few repeated ideas:

- state changes over time
- the DOM reflects that state
- events trigger updates
- reusable render/update patterns keep code understandable

By the end of this lesson, you should be able to trace a small interactive feature from user action to updated UI.

## Interactivity is a loop, not a trick

Most interactive frontend behavior follows a repeatable loop:

1. the user does something
2. your code receives the event
3. your state changes
4. the UI re-renders to match the new state

That loop scales from a tiny counter to a full app.

:::key
Strong frontend code is not a pile of random DOM edits. It is a loop connecting events, state, and rendering.
:::

## The select → listen → update pattern

At the smallest level, many interactive features look like this:

```js
const button = document.getElementById("likeBtn");
const output = document.getElementById("count");
let likes = 0;

button.addEventListener("click", function () {
  likes = likes + 1;
  output.textContent = likes;
});
```

This is still the core pattern:

- **select** the elements
- **listen** for an event
- **update** state or UI

## Rendering from arrays

As soon as you have lists, rendering becomes more interesting.

```js
const tasks = ["Write intro", "Fix hero spacing", "Test on mobile"];
const list = document.getElementById("taskList");

function renderTasks() {
  list.innerHTML = "";

  tasks.forEach(function (task) {
    const li = document.createElement("li");
    li.textContent = task;
    list.appendChild(li);
  });
}
```

This matters because lots of frontend features are just lists with behavior:

- task lists
- search results
- comments
- navigation menus
- cards

## Events often come from forms and buttons

Two of the most common event sources are:

- buttons with `click`
- forms with `submit`

A form example:

```js
form.addEventListener("submit", function (event) {
  event.preventDefault();
  // validate, update state, render feedback
});
```

That `preventDefault()` call is often crucial because it stops the browser from doing its normal full-page submission behavior.

:::quiz
Q: Why is `event.preventDefault()` often used in a form submit handler?
- To clear all variables
- To stop the browser's default submit/reload behavior *
- To create a new DOM node
E: It lets your JavaScript control the submission flow instead of letting the browser immediately reload or navigate away.
:::

## UI state examples you should recognize

Interactive sites usually need more than one UI mode.

Common states include:

- loading
- empty
- error
- success
- active / inactive
- open / closed

Example:

```js
const state = {
  isMenuOpen: false,
  isSaving: false,
  error: ""
};
```

If you get comfortable modeling these states, frontend problems become much easier to reason about.

## Use classes to reflect state

A helpful pattern is to let CSS own appearance while JavaScript only changes classes or content.

```js
panel.classList.toggle("hidden", !state.isMenuOpen);
button.classList.toggle("is-active", state.isMenuOpen);
```

This is cleaner than stuffing lots of inline style changes into JavaScript.

## Mini feature: live character counter

Imagine a textarea with a 200-character limit.

State and UI flow:

- user types
- input event fires
- code reads the value length
- counter text updates
- maybe the submit button disables if too long

```js
const field = document.getElementById("bio");
const counter = document.getElementById("counter");

field.addEventListener("input", function () {
  counter.textContent = field.value.length + " / 200";
});
```

This is a small but real example of responsive UI behavior.

## Mini feature: toggle panels and menus

Many UIs need open/close behavior.

```js
const toggleBtn = document.getElementById("faqBtn");
const answer = document.getElementById("faqAnswer");

let open = false;

toggleBtn.addEventListener("click", function () {
  open = !open;
  answer.classList.toggle("hidden", !open);
});
```

That same pattern powers:

- accordions
- dropdown menus
- modal panels
- mobile nav

## Render functions keep complexity under control

If your interface has multiple moving parts, a render function helps.

```js
function render() {
  saveButton.disabled = state.isSaving;
  errorBox.textContent = state.error;
  menu.classList.toggle("hidden", !state.isMenuOpen);
}
```

Now your event handlers can focus on changing state, then call `render()`.

This makes code much easier to debug than scattered one-off updates.

:::fill
Q: Complete the line so the menu is hidden when `state.isMenuOpen` is false.
`menu.classList.toggle("hidden", ___);`
- !state.isMenuOpen *
- state.isMenuOpen
- menu
E: If the menu is not open, the `hidden` class should be applied.
:::

## Architecture mindset for small interactive sites

As features grow, use this sequence:

- define the data you need
- write the HTML structure
- select the needed elements
- write event handlers
- centralize repeated UI updates into render logic

That is the path from demo-level code to maintainable frontend code.

## Common mistakes to avoid

- updating the DOM in ten different places with no single pattern
- treating the DOM itself as the source of truth
- forgetting submit behavior on forms
- coupling visual styling directly into JS with lots of `.style.*` changes
- writing giant event handlers that do everything

:::warning
If you can't explain what the current state of the UI is, the code will eventually become hard to trust.
:::

## What good looks like

You should now be able to:

- describe the event → state → render loop
- build small interactive features around state changes
- render lists from arrays
- use form and button events intentionally
- reflect state in the DOM with text and classes
- recognize common UI states like loading, error, and success

## What's next

In **Project: Build a Quiz Game**, you'll use these patterns to build a complete mini-app with data, events, score tracking, and a replay flow.
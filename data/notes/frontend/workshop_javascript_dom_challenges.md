# Workshop: JavaScript DOM Challenges

This workshop is about repetition with purpose.

Instead of one long feature, you'll practice the same frontend loop across several small DOM problems:

> select → listen → update

Each challenge targets a core JavaScript interaction skill you'll reuse in larger projects.

:::project
Complete four mini challenges: a counter, a dark-mode toggle, a live password helper, and a removable list item flow. If you can finish all four, your DOM instincts will be much stronger.
:::

## Challenge 1 — Build a counter

### Goal

Create a button that increases a visible count each time it is clicked.

### Requirements

- one button
- one visible number
- count starts at `0`
- each click adds `1`

### Hints

- store the count in a variable
- select the display element
- update `textContent` after each click

Starter HTML:

```html
<p id="count">0</p>
<button id="increment">Add one</button>
```

## Challenge 2 — Build a theme toggle

### Goal

Toggle a `dark` class on the page when a button is clicked.

### Requirements

- one toggle button
- one class added/removed on `<body>`
- button label should change too (`Dark mode` / `Light mode`)

### Hints

- track whether dark mode is on
- use `classList.toggle()` or add/remove classes directly
- update the button text after the state changes

Starter HTML:

```html
<button id="themeBtn">Dark mode</button>
```

## Challenge 3 — Live password helper

### Goal

As the user types into a password field, show whether the password is strong enough.

### Requirements

- input field
- helper text below it
- helper says something like:
  - `Too short` when fewer than 8 characters
  - `Looks good` when 8 or more characters

Starter HTML:

```html
<label for="password">Password</label>
<input id="password" type="password" />
<p id="passwordHelp"></p>
```

### Hints

- use the `input` event
- read `field.value.length`
- update helper text dynamically

## Challenge 4 — Remove an item from a list

### Goal

Render a list of tasks with a delete button next to each one. Clicking delete removes that item.

### Requirements

- tasks stored in an array
- render function outputs list items
- each item gets a delete button
- delete updates the array and re-renders the list

Starter idea:

```js
const tasks = ["Write homepage", "Fix hero image", "Test contact form"];
```

### Hints

- use `forEach` when rendering
- use `splice()` to remove an item by index
- re-render after deleting

## Workshop review questions

After you complete the challenges, answer these:

- Which challenge only needed one piece of state?
- Which challenge required state plus rendering from an array?
- Which one felt easiest to reason about, and why?

:::quiz
Q: Which challenge most clearly demonstrates the render-from-array pattern?
- Counter
- Theme toggle
- Removable list item flow *
E: When rendering from an array and deleting by index, you're practicing the same pattern used in many mini-apps and product interfaces.
:::

## Debugging reminders

If something doesn't work:

- check your selectors first
- confirm the event is firing
- log the state after it changes
- confirm the render function is actually being called

:::key
Small DOM exercises are not childish practice — they are how you build fast, reliable instincts before moving into bigger app logic.
:::

## What's next

In **Forms, State, and UI Patterns**, you'll turn these interaction moves into more structured UI behavior with validation, success states, and cleaner rendering.
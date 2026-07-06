# DOM Manipulation & Events

JavaScript becomes useful in the browser the moment it can **see the page** and **react to the user**.

That is what the DOM and events give you.

By the end of this lesson, you should be able to:

- select elements from the page
- change text, classes, and structure
- listen for user actions
- connect page updates to real UI state

## What the DOM is

DOM stands for **Document Object Model**.

When the browser reads your HTML, it turns it into a live tree of nodes. JavaScript can then:

- find elements
- read their content
- update their content
- create new elements
- remove elements
- respond to interaction

If HTML is the structure and CSS is the presentation, the DOM is the **JavaScript-accessible version of the page**.

## Selecting elements

You'll do this constantly.

Common approaches:

```js
const title = document.getElementById("title");
const button = document.querySelector(".cta-button");
const cards = document.querySelectorAll(".card");
```

Use:

- `getElementById` when you know the exact id
- `querySelector` when you want the first element matching a CSS selector
- `querySelectorAll` when you want a list of matches

## Updating content

A simple update looks like this:

```js
const message = document.getElementById("message");
message.textContent = "Saved successfully";
```

Useful properties and methods:

- `textContent`
- `innerHTML` (use carefully)
- `classList.add()`
- `classList.remove()`
- `classList.toggle()`
- `setAttribute()`

:::fill
Q: Complete the line so the element gets the `active` class.
`card.classList.___("active");`
- add *
- textContent
- querySelector
E: `classList.add()` is the simplest way to attach a class to an element.
:::

## Creating and inserting elements

You are not limited to elements already in the HTML.

```js
const li = document.createElement("li");
li.textContent = "Ship a cleaner UI";
document.getElementById("tasks").appendChild(li);
```

This is how dynamic lists, notifications, cards, and search results are often rendered.

## Events: how the page reacts to users

An event is something that happens in the browser:

- a click
- typing into an input
- submitting a form
- hovering
- scrolling

Listening for an event looks like this:

```js
const button = document.getElementById("saveBtn");

button.addEventListener("click", function () {
  console.log("Saved");
});
```

That one pattern — **select → listen → update** — powers a huge amount of frontend work.

:::key
Most UI behavior can be described as: something happened, so update the state or the page.
:::

## The event object

Many handlers receive an event object with useful information.

```js
form.addEventListener("submit", function (event) {
  event.preventDefault();
  console.log("Form submitted without page reload");
});
```

Useful event patterns:

- `event.preventDefault()` — stop default browser behavior
- `event.target` — the element that triggered the event

## Class updates are often cleaner than inline styles

Instead of doing this everywhere:

```js
box.style.backgroundColor = "green";
```

Prefer toggling classes:

```js
box.classList.add("is-success");
```

Why?

- CSS stays responsible for appearance
- JavaScript stays responsible for behavior
- your code is easier to maintain

## Mini example — a dismissible alert

HTML:

```html
<div id="notice" class="notice">Profile saved</div>
<button id="closeBtn">Dismiss</button>
```

JavaScript:

```js
const notice = document.getElementById("notice");
const closeBtn = document.getElementById("closeBtn");

closeBtn.addEventListener("click", function () {
  notice.classList.add("hidden");
});
```

That is already a real UI behavior.

## Render from data, not chaos

As soon as an interface gets more complex, it's better to think from data outward.

```js
const tasks = ["Ship hero section", "Fix mobile nav"];
const list = document.getElementById("tasks");

list.innerHTML = "";

tasks.forEach(function (task) {
  const item = document.createElement("li");
  item.textContent = task;
  list.appendChild(item);
});
```

This pattern matters because later your app state might come from:

- form input
- API data
- localStorage
- filters or search queries

## Predict the result

:::predict
```js
const count = document.getElementById("count");
let clicks = 0;

button.addEventListener("click", function () {
  clicks++;
  count.textContent = clicks;
});
```
Q: What happens each time the button is clicked?
- The page reloads
- The visible count increases by 1 *
- The `count` element is deleted
E: Each click increments the `clicks` variable and writes the updated value into the DOM.
:::

## Mini practice — build a toggle

Try building a tiny FAQ interaction.

HTML:

```html
<button id="faqBtn">What does Sprintboard do?</button>
<p id="faqAnswer" class="hidden">It helps teams plan weekly priorities.</p>
```

JavaScript idea:

- select the button and paragraph
- listen for a click on the button
- toggle the `hidden` class on the answer

This is the perfect first DOM exercise because it's small, visible, and real.

## Common mistakes to avoid

- selecting the wrong element
- writing lots of inline styles in JS instead of classes
- forgetting `preventDefault()` on form handlers when needed
- mutating the DOM repeatedly with no plan
- putting all logic directly inside one giant click handler

:::warning
If your JavaScript is manually patching random bits of the page everywhere, you probably need clearer state and cleaner render logic.
:::

## What good looks like

You should now be able to:

- select individual elements and groups of elements
- update text and classes
- create and append new nodes
- attach click and submit listeners
- use the event object intentionally
- think in the pattern: select → listen → update

## What's next

In **Forms, State, and UI Patterns**, you'll move from single interactions to a more important frontend idea: keeping UI behavior organized around state.
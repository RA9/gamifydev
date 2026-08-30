# JavaScript Events

A web page that never reacts is just a poster. Events are what make a page *respond* — to clicks, typing, hovering, and form submissions. In this lesson you will learn to listen for events and run your own code when they happen.

## What is an event?

An **event** is something that happens in the browser: the user clicks a button, types into a box, moves the mouse, submits a form, or the page finishes loading. Your job is to *listen* for the events you care about and run a function — called an **event handler** — when they fire.

:::analogy
Think of an event listener like a doorbell. You install the bell (the listener) on a specific door (an element). When someone presses it (the event), a specific chime plays (your handler). The door does nothing until the bell is pressed — and you decide what the chime is.
:::

## addEventListener

The modern way to listen is `element.addEventListener(type, handler)`. The first argument is the event name as a string; the second is the function to run.

```html
<button id="greet">Say hi</button>
```

```js
const button = document.querySelector("#greet");

button.addEventListener("click", function () {
  alert("Hi there!");
});
```

You can also pass an arrow function, which is the common style today:

```js
button.addEventListener("click", () => {
  console.log("Button was clicked");
});
```

You can attach as many listeners as you like, even for the same event — they all run. That flexibility is why `addEventListener` is preferred over the old `onclick =` attribute approach.

:::warning
Pass the function itself, *not* the result of calling it. `addEventListener("click", handleClick)` is correct. `addEventListener("click", handleClick())` runs `handleClick` immediately and passes its return value (often `undefined`) as the handler — a classic mistake.
:::

## Common event types

You will reach for these constantly:

- `click` — the element was clicked
- `input` — the value of a text box changed (fires on every keystroke)
- `change` — a form control lost focus after its value changed (great for `<select>` and checkboxes)
- `submit` — a form was submitted
- `keydown` — a key was pressed
- `mouseover` / `mouseout` — the pointer entered / left an element
- `DOMContentLoaded` — the HTML has been fully parsed

```js
document.addEventListener("DOMContentLoaded", () => {
  console.log("Page is ready, safe to grab elements now");
});
```

`DOMContentLoaded` matters because if your script runs before the elements exist, `querySelector` returns `null`. Putting setup code inside this listener guarantees the page is ready.

## The event object

Your handler automatically receives an **event object** as its first argument. It carries details about what happened. By convention it is named `event` or just `e`.

```js
const button = document.querySelector("#greet");

button.addEventListener("click", (event) => {
  console.log(event.type);   // "click"
  console.log(event.target); // the element that was clicked
});
```

### event.target

`event.target` is the actual element the event happened on. This is powerful: one handler can find out exactly which element triggered it.

```html
<div id="menu">
  <button>Home</button>
  <button>About</button>
  <button>Contact</button>
</div>
```

```js
const menu = document.querySelector("#menu");

menu.addEventListener("click", (event) => {
  console.log("You clicked:", event.target.textContent);
});
```

One listener on the `<div>` reports which button was clicked. We will build on this idea with *event delegation* below.

## event.preventDefault

Some elements have built-in default behavior. A form submit reloads the page; a link navigates away. `event.preventDefault()` cancels that default so your JavaScript can take over.

```html
<form id="signup">
  <input id="email" type="email" placeholder="you@example.com" />
  <button type="submit">Sign up</button>
</form>
<p id="status"></p>
```

```js
const form = document.querySelector("#signup");
const status = document.querySelector("#status");

form.addEventListener("submit", (event) => {
  event.preventDefault(); // stop the page from reloading

  const email = document.querySelector("#email").value;
  status.textContent = `Thanks! We'll email ${email}.`;
});
```

:::key
Without `event.preventDefault()`, submitting a form reloads the page and your JavaScript result vanishes instantly. Any time you handle a form with JS, calling `preventDefault` on `submit` is almost always step one.
:::

## Reading input values with the input event

The `input` event fires on every keystroke, which is perfect for live feedback. Read the current text from `event.target.value` (or the input's `.value`).

```html
<input id="name" placeholder="Type your name" />
<p>Hello, <span id="echo">stranger</span>!</p>
```

```js
const nameField = document.querySelector("#name");
const echo = document.querySelector("#echo");

nameField.addEventListener("input", (event) => {
  echo.textContent = event.target.value || "stranger";
});
```

As you type, the greeting updates in real time. This "live mirror" pattern is everywhere: search-as-you-type, password strength meters, character counters.

## A click counter

Putting it together — a button that counts its own clicks. Notice the count lives in a variable *outside* the handler so it persists between clicks.

```html
<button id="counter">Clicked 0 times</button>
```

```js
const btn = document.querySelector("#counter");
let count = 0;

btn.addEventListener("click", () => {
  count++;
  btn.textContent = `Clicked ${count} times`;
});
```

:::tip
Variables a handler needs to remember between firings belong *outside* the handler. If you wrote `let count = 0` inside the handler, it would reset to 0 on every click and never grow past 1.
:::

## Event delegation

Imagine a list where items get added and removed all the time. Attaching a listener to each item is fragile — new items added later have no listener. The fix is **event delegation**: attach *one* listener to a stable parent, and use `event.target.closest()` to figure out which child was involved.

`closest(selector)` walks up from the clicked element to find the nearest ancestor (or itself) matching the selector. This handles clicks on nested content gracefully.

```html
<ul id="tasks">
  <li>Buy milk <button class="del">x</button></li>
  <li>Walk dog <button class="del">x</button></li>
</ul>
```

```js
const tasks = document.querySelector("#tasks");

tasks.addEventListener("click", (event) => {
  // Did the click land on a delete button (or inside one)?
  const deleteBtn = event.target.closest(".del");
  if (!deleteBtn) return; // clicked elsewhere — ignore

  const item = deleteBtn.closest("li");
  item.remove();
});
```

Because the listener lives on `#tasks`, it works for items that exist now *and* any added later — no need to wire up new listeners each time.

:::example
Why the `if (!deleteBtn) return` guard? Clicks anywhere in the list bubble up to the parent listener, including on the text. The guard makes sure we only act when a delete button was actually the target.
:::

## Removing listeners

To stop listening, call `removeEventListener` with the *same* event type and the *same* function reference. This is why you must use a named function — you cannot remove an anonymous one.

```js
function handleClick() {
  console.log("clicked once");
  button.removeEventListener("click", handleClick);
}

button.addEventListener("click", handleClick);
```

Here the handler removes itself after the first click, creating a one-time listener.

:::warning
`removeEventListener` only works if you pass the exact same function you added. Two separate arrow functions with identical code are still different objects, so `removeEventListener("click", () => {...})` will silently fail to remove anything.
:::

## Practice

:::quiz
Q: Why do we call `event.preventDefault()` inside a form's submit handler?
- To clear the form fields
- To stop the browser from reloading the page on submit *
- To submit the form faster
- To validate the email automatically
E: A form's default submit behavior reloads (or navigates) the page, which wipes out any result your JavaScript produced. `preventDefault` cancels that so your code stays in control.
:::

:::predict
What appears as you type "Sam" into the box, character by character?
```js
field.addEventListener("input", (e) => {
  output.textContent = e.target.value.toUpperCase();
});
```
:::

:::fill
With event delegation, you put one listener on a parent and use `event.target.______(selector)` to find which matching child was clicked.
- closest
E: `closest` walks up from the event target to the nearest ancestor matching the selector, so a single parent listener can handle clicks on any child, including ones added later.
:::

## Recap

- An **event** is a browser happening (click, input, submit...); a **handler** is the function you run in response.
- Listen with `element.addEventListener(type, handler)` — pass the function, do not call it.
- The **event object** carries details; `event.target` is the element that triggered it.
- `event.preventDefault()` cancels default behavior — essential for form `submit`.
- The `input` event plus `.value` powers live, type-as-you-go feedback.
- **Event delegation** (one parent listener + `event.target.closest`) handles dynamic lists cleanly.
- Remove listeners with `removeEventListener` using the same named function reference.

**Next up:** Building Interactive JavaScript Websites — combining the DOM and events to build real features.

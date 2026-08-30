# Building Reusable UI Components

As your projects grow, you find yourself writing the same card layout, the same button group, the same list item structure in multiple places. This lesson shows you how to build reusable pieces in vanilla JavaScript — functions that accept data and return DOM — and why this thinking maps directly to how frameworks like React work.

## The problem: copy-pasted DOM code

Without reusable components, you end up with repetition:

```js
// Creating a card in three different places
const card1 = document.createElement("div");
card1.className = "card";
card1.innerHTML = `<h3>Project A</h3><p>Description A</p>`;
container.append(card1);

// same code again...
const card2 = document.createElement("div");
card2.className = "card";
card2.innerHTML = `<h3>Project B</h3><p>Description B</p>`;
container.append(card2);
```

If you need to change the card structure, you have to find and fix every copy. That is fragile and boring.

## The solution: functions that return DOM

A component is just a function that takes data and returns an element:

```js
function createCard({ title, description, tag }) {
  const card = document.createElement("div");
  card.className = "card";

  const heading = document.createElement("h3");
  heading.textContent = title;

  const body = document.createElement("p");
  body.textContent = description;

  const badge = document.createElement("span");
  badge.className = "badge";
  badge.textContent = tag;

  card.append(heading, body, badge);
  return card;
}
```

Now creating cards is a one-liner:

```js
container.append(
  createCard({ title: "Portfolio", description: "My personal site", tag: "design" }),
  createCard({ title: "API Server", description: "REST endpoints", tag: "dev" }),
);
```

:::key
A component function takes **data in** and returns **DOM out**. The structure is defined in one place. When you need to change the card layout, you change one function.
:::

## Separating data from presentation

Keep your data as plain objects or arrays, and let component functions handle how it looks:

```js
const projects = [
  { title: "Portfolio", description: "Personal site", tag: "design" },
  { title: "CLI Tool", description: "File converter", tag: "dev" },
  { title: "Blog", description: "Tech articles", tag: "design" },
];

function renderProjects(projects) {
  const container = document.querySelector("#project-grid");
  container.innerHTML = "";

  const fragment = document.createDocumentFragment();
  projects.forEach(project => {
    fragment.append(createCard(project));
  });
  container.append(fragment);
}

renderProjects(projects);
```

The data knows nothing about HTML. The component knows nothing about where the data came from. They are cleanly separated.

## Components with callbacks

Pass event handlers as callbacks to keep components reusable:

```js
function createTodoItem({ id, text, done }, { onToggle, onDelete }) {
  const li = document.createElement("li");
  li.className = done ? "todo-item done" : "todo-item";

  const label = document.createElement("span");
  label.textContent = text;

  const toggleBtn = document.createElement("button");
  toggleBtn.textContent = done ? "Undo" : "Done";
  toggleBtn.addEventListener("click", () => onToggle(id));

  const deleteBtn = document.createElement("button");
  deleteBtn.textContent = "Delete";
  deleteBtn.addEventListener("click", () => onDelete(id));

  li.append(label, toggleBtn, deleteBtn);
  return li;
}
```

The component does not know *how* to toggle or delete — it only calls the callback. This decoupling is the same pattern React uses with props.

:::tip
Passing callbacks keeps components flexible. The same `createTodoItem` could be used in different projects with different toggle and delete implementations.
:::

## Composing small components

Build complex UIs by combining small functions:

```js
function createBadge(text) {
  const span = document.createElement("span");
  span.className = "badge";
  span.textContent = text;
  return span;
}

function createUserCard(user) {
  const card = document.createElement("div");
  card.className = "user-card";

  const name = document.createElement("h3");
  name.textContent = user.name;

  card.append(name, createBadge(user.role));
  return card;
}
```

Change the badge design once, and every user card updates.

## Encapsulating state inside components

Some components need internal state that does not belong in the global app state:

```js
function createAccordion(items) {
  let openIndex = null; // internal state

  const container = document.createElement("div");
  container.className = "accordion";

  function render() {
    container.innerHTML = "";

    items.forEach((item, index) => {
      const section = document.createElement("div");
      section.className = "accordion-section";

      const header = document.createElement("button");
      header.className = "accordion-header";
      header.textContent = item.title;
      header.addEventListener("click", () => {
        openIndex = openIndex === index ? null : index;
        render();
      });

      const body = document.createElement("div");
      body.className = "accordion-body";
      body.textContent = item.content;
      body.hidden = openIndex !== index;

      section.append(header, body);
      container.append(section);
    });
  }

  render();
  return container;
}

// Usage
const faq = createAccordion([
  { title: "What is this?", content: "A learning platform." },
  { title: "Is it free?", content: "Yes, completely free." },
]);
document.querySelector("#faq").append(faq);
```

The accordion manages its own `openIndex` — the parent does not need to know or care which section is expanded.

## A config-driven component pattern

For components with many options, accept a configuration object:

```js
function createButton({ text, variant = "primary", size = "md", onClick }) {
  const btn = document.createElement("button");
  btn.className = `btn btn-${variant} btn-${size}`;
  btn.textContent = text;
  if (onClick) btn.addEventListener("click", onClick);
  return btn;
}

// Usage
container.append(
  createButton({ text: "Save", variant: "primary", onClick: handleSave }),
  createButton({ text: "Cancel", variant: "secondary", onClick: handleCancel }),
  createButton({ text: "Delete", variant: "danger", size: "sm", onClick: handleDelete }),
);
```

Default values in the destructuring keep the API simple for common cases while allowing full customization.

## Why this maps to frameworks

Everything you have learned here is exactly how component-based frameworks work:

| Vanilla JS | React equivalent |
|---|---|
| Function that returns DOM | Function that returns JSX |
| Config object parameter | Props |
| Callback functions (`onClick`) | Event handler props |
| Internal state (`let openIndex`) | `useState` |
| Composing small functions | Composing components |
| Call render after state change | React's automatic re-render |

The mental model is identical. Frameworks add convenience (automatic re-rendering, diffing, declarative syntax), but the architecture — data in, UI out, compose small pieces — is the same pattern you are already using.

:::key
Learning to build components in vanilla JavaScript gives you a framework-ready mindset. When you move to React, Vue, or Svelte, you will recognize the same ideas dressed in different syntax.
:::

## Practice

:::quiz
Q: What is the main advantage of wrapping DOM creation in a function?
- It runs faster than inline DOM code
- The structure is defined once and reused everywhere, so changes happen in one place *
- It automatically handles event delegation
- It prevents XSS attacks
E: A component function centralizes the DOM structure. When you need to add a class, change the layout, or add an attribute, you update one function instead of hunting through every place that creates that element.
:::

:::quiz
Q: A `createTodoItem` function accepts `{ onToggle, onDelete }` callbacks. Why is this better than having the component call `toggleTask` directly?
- It is faster
- It keeps the component reusable — it does not depend on specific app functions *
- Callbacks prevent memory leaks
- Direct function calls do not work in JavaScript
E: Passing callbacks decouples the component from the app logic. The same `createTodoItem` could be used in different projects with different toggle and delete implementations. This is the same principle as React props.
:::

## Recap

- A **component** is a function that takes data and returns a DOM element.
- Separate **data** (plain objects) from **presentation** (component functions).
- Pass **callbacks** for event handling to keep components decoupled from app logic.
- **Compose** small components (avatar, badge) into larger ones (user card).
- **Encapsulate state** inside a component when it is purely internal (like which accordion section is open).
- Use **config objects** with defaults for components that accept many options.
- This is the exact architecture of **React, Vue, and Svelte** — functions that take data and return UI.

**Next up:** Local Storage and Persistence — keeping user data alive between sessions.

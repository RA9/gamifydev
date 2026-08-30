# Modal Dialogs

A modal takes over the user's attention: it appears over everything, demands a response, and blocks interaction with the page behind it. Getting this right means handling visibility, backdrop clicks, Escape, focus trapping, and accessibility — all with clean state management.

## Show and hide with state

The modal's open/closed status is a single boolean:

```js
let isModalOpen = false;

function openModal() {
  isModalOpen = true;
  renderModal();
}

function closeModal() {
  isModalOpen = false;
  renderModal();
}
```

## HTML structure

```html
<div id="modal" class="modal-backdrop" hidden>
  <div class="modal-content" role="dialog" aria-modal="true" aria-labelledby="modal-title">
    <h2 id="modal-title">Confirm Action</h2>
    <p>Are you sure you want to delete this item?</p>
    <div class="modal-actions">
      <button id="modal-cancel">Cancel</button>
      <button id="modal-confirm">Delete</button>
    </div>
  </div>
</div>
```

Key ARIA attributes:
- `role="dialog"` tells assistive tech this is a dialog.
- `aria-modal="true"` signals that interaction with the rest of the page is blocked.
- `aria-labelledby` points to the heading so screen readers announce the dialog's purpose.

## CSS

```css
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}
```

## The render function

```js
const modal = document.querySelector("#modal");

function renderModal() {
  modal.hidden = !isModalOpen;

  if (isModalOpen) {
    document.body.style.overflow = "hidden"; // prevent background scroll
  } else {
    document.body.style.overflow = "";
  }
}
```

:::key
Setting `body.style.overflow = "hidden"` prevents the page behind the modal from scrolling. Always restore it to `""` (empty string, not `"auto"`) when the modal closes so the original CSS value takes effect.
:::

## Closing the modal — three ways

Users expect to close modals by:
1. Clicking a close/cancel button
2. Clicking the backdrop (the dark overlay)
3. Pressing the Escape key

```js
// 1. Button clicks
document.querySelector("#modal-cancel").addEventListener("click", closeModal);

// 2. Backdrop click — close only if the click is on the backdrop itself
modal.addEventListener("click", (event) => {
  if (event.target === modal) {
    closeModal();
  }
});

// 3. Escape key
document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && isModalOpen) {
    closeModal();
  }
});
```

:::tip
For the backdrop click, check `event.target === modal` (the backdrop element), not just any click inside the modal. Clicks on buttons and text inside `.modal-content` should *not* close the modal — they bubble up to the backdrop, so without this check, any click anywhere would close it.
:::

## Focus trapping

When a modal is open, pressing Tab should cycle through elements *inside* the modal, never escaping to the page behind it.

```js
function trapFocus(event) {
  if (event.key !== "Tab") return;

  const focusable = modal.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );
  const first = focusable[0];
  const last = focusable[focusable.length - 1];

  if (event.shiftKey) {
    // Shift+Tab: if on first element, jump to last
    if (document.activeElement === first) {
      event.preventDefault();
      last.focus();
    }
  } else {
    // Tab: if on last element, jump to first
    if (document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  }
}
```

Wire it up in your open/close functions:

```js
function openModal() {
  isModalOpen = true;
  renderModal();
  document.addEventListener("keydown", trapFocus);

  // Move focus into the modal
  const firstFocusable = modal.querySelector(
    'button, [href], input, select, textarea'
  );
  if (firstFocusable) firstFocusable.focus();
}

function closeModal() {
  isModalOpen = false;
  renderModal();
  document.removeEventListener("keydown", trapFocus);
}
```

## Restoring focus on close

When the modal closes, focus should return to the element that opened it. Otherwise keyboard users are lost:

```js
let previouslyFocused = null;

function openModal() {
  previouslyFocused = document.activeElement; // remember where focus was
  isModalOpen = true;
  renderModal();
  document.addEventListener("keydown", trapFocus);

  const firstFocusable = modal.querySelector("button, [href], input");
  if (firstFocusable) firstFocusable.focus();
}

function closeModal() {
  isModalOpen = false;
  renderModal();
  document.removeEventListener("keydown", trapFocus);

  if (previouslyFocused) {
    previouslyFocused.focus(); // return focus to the trigger
  }
}
```

:::key
Save `document.activeElement` before opening the modal and call `.focus()` on it after closing. This is critical for keyboard accessibility — without it, focus disappears into the void after the modal closes.
:::

## The native dialog element

Modern browsers have a built-in `<dialog>` element that handles some of these patterns for you:

```html
<dialog id="native-dialog">
  <h2>Native Dialog</h2>
  <p>This uses the built-in dialog element.</p>
  <button id="close-dialog">Close</button>
</dialog>

<button id="open-dialog">Open Dialog</button>
```

```js
const dialog = document.querySelector("#native-dialog");

document.querySelector("#open-dialog").addEventListener("click", () => {
  dialog.showModal(); // opens as a modal with built-in backdrop
});

document.querySelector("#close-dialog").addEventListener("click", () => {
  dialog.close();
});

// The dialog fires a "close" event
dialog.addEventListener("close", () => {
  console.log("Dialog closed with value:", dialog.returnValue);
});
```

`showModal()` gives you:
- A built-in backdrop (styleable with `::backdrop`)
- Escape to close (by default)
- Focus trapping (by default)
- `aria-modal` behavior (by default)

```css
dialog::backdrop {
  background: rgba(0, 0, 0, 0.5);
}
```

:::tip
If browser support covers your audience, the `<dialog>` element saves you from implementing focus trapping, Escape handling, and backdrop behavior manually. It is the recommended approach for new projects.
:::

## Practice

:::quiz
Q: A user opens a modal and presses Tab until they reach the last focusable element inside it. What should happen on the next Tab press?
- Focus moves to the browser's address bar
- Focus moves to the first element behind the modal
- Focus wraps to the first focusable element inside the modal *
- Nothing happens
E: Focus trapping keeps the Tab cycle inside the modal. When the user tabs past the last element, focus wraps back to the first. This prevents keyboard users from accidentally interacting with the page behind the modal.
:::

:::quiz
Q: After closing a modal, where should focus go?
- The first element on the page
- The body element
- The element that triggered the modal to open *
- The last focused element inside the modal
E: Restoring focus to the trigger element gives keyboard users a predictable location. Save `document.activeElement` before opening and call `.focus()` on it when closing.
:::

## Recap

- Modal state is a **single boolean**. Open and close functions update it and call render.
- Support **three close methods:** button click, backdrop click (`event.target === backdrop`), and Escape key.
- **Trap focus** inside the modal by intercepting Tab/Shift+Tab and wrapping between the first and last focusable elements.
- **Restore focus** to the trigger element when the modal closes — save `document.activeElement` on open.
- Prevent **background scroll** with `overflow: hidden` or a `position: fixed` approach.
- Use ARIA: `role="dialog"`, `aria-modal="true"`, `aria-labelledby`.
- The native `<dialog>` element with `showModal()` provides focus trapping, Escape, and backdrop for free.

**Next up:** Multi-Step Form Flows — walking users through complex input one step at a time.

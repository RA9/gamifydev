# Keyboard Accessibility

Not everyone uses a mouse. Some people use only a keyboard, others use switch devices or voice control — all of which ultimately navigate the page like a keyboard. If your site works with a keyboard, it works for all of them. In this lesson you'll learn tab order, focus management, visible focus styles, and how to build custom components that keyboard users can actually operate.

## Why Keyboard Access Matters

Keyboard accessibility isn't a niche concern. It serves:

- People with **motor disabilities** who can't use a mouse.
- **Power users** who prefer keyboard shortcuts for speed.
- **Screen reader users** — screen readers navigate via keyboard commands.
- Anyone with a **temporary injury** — a broken arm, a trackpad that stopped working.
- **Automated testing tools** that simulate keyboard interactions.

If a button can only be clicked with a mouse, all of these users are locked out.

:::key
The fundamental rule: every interactive element on your page must be reachable and operable using only the keyboard. If you can't Tab to it and activate it with Enter or Space, it's broken.
:::

## Natural Tab Order

When users press **Tab**, the browser moves focus through interactive elements in **document order** — the order they appear in the HTML source. This is called the natural tab order.

These elements are focusable by default:

- `<a href="…">` — links
- `<button>` — buttons
- `<input>`, `<select>`, `<textarea>` — form controls
- `<details>` / `<summary>` — disclosure widgets

These are **not** focusable by default:

- `<div>`, `<span>`, `<p>` — generic elements
- `<img>` — images
- `<h1>` through `<h6>` — headings

If you need a non-interactive element to receive focus, use `tabindex`.

## Understanding `tabindex`

`tabindex` controls whether and when an element can receive focus.

### `tabindex="0"` — Add to Tab Order

Adds the element to the natural tab order, at its position in the DOM. Use this when you're building a custom interactive component from a non-interactive element.

```html
<div role="button" tabindex="0" onclick="doSomething()">
  Custom Button
</div>
```

### `tabindex="-1"` — Programmatically Focusable Only

The element can receive focus via JavaScript (`element.focus()`) but is *not* reachable by Tab. Use this for elements you want to manage focus to programmatically — like moving focus into a modal or to an error message.

```html
<h2 tabindex="-1" id="error-section">Form Errors</h2>

<script>
  // After validation fails, move focus to the error heading
  document.getElementById('error-section').focus();
</script>
```

### Positive `tabindex` — Never Use This

`tabindex="1"`, `tabindex="5"`, or any positive number forces an element to the *front* of the tab order. This sounds useful but creates chaos: the tab order no longer matches the visual order, confusing every keyboard user.

:::warning
Never use positive `tabindex` values. They break the natural reading order and create a maintenance nightmare. If you need an element to come earlier in the tab order, move it earlier in the HTML source.
:::

```html
<!-- Bad: unpredictable tab order -->
<button tabindex="3">Third</button>
<button tabindex="1">First</button>
<button tabindex="2">Second</button>

<!-- Good: let document order do the work -->
<button>First</button>
<button>Second</button>
<button>Third</button>
```

:::quiz
Q: What does `tabindex="-1"` do?
- Removes the element from the page entirely
- Makes the element focusable via JavaScript but not reachable by Tab *
- Adds the element to the tab order at position -1
E: `tabindex="-1"` allows programmatic focusing with `element.focus()` while keeping the element out of the Tab key sequence. It's essential for focus management in modals, error messages, and single-page navigation.
:::

## Visible Focus Styles

Keyboard users need to *see* where they are. The browser shows a focus ring by default, but many developers remove it for aesthetic reasons — breaking keyboard navigation completely.

### The `:focus-visible` Selector

`:focus-visible` applies only when the browser determines that focus should be visible — typically on keyboard navigation, not mouse clicks. This gives you the best of both worlds.

```css
/* Remove the default outline (only safe because we're replacing it) */
button:focus {
  outline: none;
}

/* Add a visible focus ring for keyboard users */
button:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px;
  border-radius: 4px;
}
```

Mouse users won't see the outline (they don't need it — they can see the cursor). Keyboard users will see a clear, high-contrast ring.

:::key
Never write `outline: none` or `outline: 0` without providing a replacement visible focus style via `:focus-visible`. This is the single most common keyboard accessibility violation on the web.
:::

### What a Good Focus Style Looks Like

A good focus indicator should:

- Have a **contrast ratio of at least 3:1** against adjacent colors.
- Be visible on **all backgrounds** (light and dark).
- Be at least **2px thick** — thin outlines are easy to miss.
- Use **outline** rather than border (outline doesn't affect layout).

```css
:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px;
}

/* High contrast for dark backgrounds */
.dark-section :focus-visible {
  outline-color: #93c5fd;
}
```

## Skip Links

A skip link lets keyboard users jump past repetitive navigation to the main content. Without it, a keyboard user has to Tab through every nav link on *every page load*.

```html
<body>
  <a href="#main" class="skip-link">Skip to main content</a>
  <header>
    <nav>
      <!-- 20 navigation links -->
    </nav>
  </header>
  <main id="main" tabindex="-1">
    <h1>Page Content</h1>
  </main>
</body>
```

```css
.skip-link {
  position: absolute;
  top: -100%;
  left: 1rem;
  padding: 0.75rem 1.5rem;
  background: #1f2937;
  color: #ffffff;
  font-weight: bold;
  z-index: 1000;
  border-radius: 0 0 8px 8px;
}

.skip-link:focus {
  top: 0;
}
```

The link is invisible until focused. When a keyboard user presses Tab on page load, the skip link appears first, letting them jump straight to `<main>`.

:::tip
Add `tabindex="-1"` to the `<main>` target so that `element.focus()` works in all browsers. Some older browsers don't move focus to non-interactive elements via fragment links without it.
:::

## Focus Trapping in Modals

When a modal opens, focus must be **trapped inside it** — Tab should cycle through the modal's interactive elements and never escape to the page behind. When the modal closes, focus must **return to the element that opened it**.

The native `<dialog>` element handles this automatically with `showModal()`:

```html
<dialog id="confirm-dialog">
  <h2>Confirm Deletion</h2>
  <p>This cannot be undone.</p>
  <button onclick="this.closest('dialog').close()">Cancel</button>
  <button onclick="deleteAccount()">Delete</button>
</dialog>
```

```javascript
let previouslyFocused = null;

function openModal() {
  previouslyFocused = document.activeElement;
  document.getElementById('confirm-dialog').showModal();
}

function closeModal() {
  document.getElementById('confirm-dialog').close();
  if (previouslyFocused) previouslyFocused.focus();
}
```

:::tip
Prefer native `<dialog>` with `showModal()` — it gives you focus trapping, Escape to close, and `aria-modal` for free. Only build custom focus trapping (intercepting Tab/Shift+Tab on the first and last focusable elements) if you can't use `<dialog>`.

## Roving Tabindex for Custom Widgets

For composite widgets — tab lists, toolbars, menus, tree views — having every item in the tab order is overwhelming. The **roving tabindex** pattern solves this:

- Only the currently active item has `tabindex="0"`.
- All other items have `tabindex="-1"`.
- Arrow keys move the active item; Tab moves *out* of the widget.

```html
<div role="tablist">
  <button role="tab" tabindex="0" aria-selected="true">Tab 1</button>
  <button role="tab" tabindex="-1" aria-selected="false">Tab 2</button>
  <button role="tab" tabindex="-1" aria-selected="false">Tab 3</button>
</div>
```

```javascript
const tabs = document.querySelectorAll('[role="tab"]');
let currentIndex = 0;

tabs.forEach((tab, index) => {
  tab.addEventListener('keydown', (e) => {
    let newIndex = currentIndex;

    if (e.key === 'ArrowRight') newIndex = (currentIndex + 1) % tabs.length;
    if (e.key === 'ArrowLeft') newIndex = (currentIndex - 1 + tabs.length) % tabs.length;

    if (newIndex !== currentIndex) {
      tabs[currentIndex].setAttribute('tabindex', '-1');
      tabs[currentIndex].setAttribute('aria-selected', 'false');

      tabs[newIndex].setAttribute('tabindex', '0');
      tabs[newIndex].setAttribute('aria-selected', 'true');
      tabs[newIndex].focus();

      currentIndex = newIndex;
    }
  });
});
```

One Tab press enters the widget; arrow keys navigate within; the next Tab press exits. This matches native OS behavior and is what screen reader users expect.

## Keyboard-Only Testing Checklist

Before shipping, unplug your mouse and run through this:

- [ ] Can I **Tab** to every link, button, and form control?
- [ ] Is Tab order **logical** and matches the visual layout?
- [ ] Is there a **visible focus indicator** on every focused element?
- [ ] Can I activate buttons with **Enter** and **Space**?
- [ ] Can I follow links with **Enter**?
- [ ] Does the **skip link** appear on the first Tab press?
- [ ] Is focus **trapped inside modals** while they're open?
- [ ] Does focus **return to the trigger** when a modal closes?
- [ ] Can I **close** modals and dropdowns with **Escape**?
- [ ] Do custom widgets (tabs, menus) support **arrow key** navigation?
- [ ] Are there any **focus traps** I can't escape?
- [ ] Can I reach and dismiss all **notifications and alerts**?

:::tip
Spend five minutes tabbing through your site with no mouse once a week during development. This single habit catches more keyboard accessibility issues than any automated scanner.
:::

:::quiz
Q: Why should you never use positive `tabindex` values like `tabindex="5"`?
- They make elements invisible to screen readers
- They break the natural document-order tab sequence, confusing keyboard users *
- They slow down page rendering
E: Positive `tabindex` values force elements to the front of the tab order regardless of their position in the HTML. This creates a confusing, unpredictable tab sequence. Use document order for tab sequence and `tabindex="0"` or `tabindex="-1"` for the rare cases where you need to add or manage focus.
:::

## Recap

- Every interactive element must be reachable and operable via keyboard.
- Use `tabindex="0"` to add elements to the tab order; `tabindex="-1"` for programmatic focus only; never use positive values.
- Provide visible focus styles with `:focus-visible` — never remove the outline without a replacement.
- Implement skip links to let keyboard users bypass navigation.
- Trap focus inside modals and restore it to the trigger on close.
- Use roving tabindex for composite widgets (tabs, toolbars, menus).
- Test with keyboard only — regularly, not just once before launch.

**Next up:** HTML Best Practices and Review — polishing your pages for performance, SEO, and production readiness.

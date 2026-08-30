# ARIA and Screen Readers

ARIA — Accessible Rich Internet Applications — is a set of HTML attributes that add accessibility information to elements when native HTML can't express it alone. It's powerful, essential for complex interfaces, and dangerous if misused. In this lesson you'll learn what ARIA does, when to reach for it, and how to avoid the most common pitfalls.

## What ARIA Actually Does

ARIA doesn't change how an element looks or behaves. It changes what **assistive technology** (screen readers, switch devices, voice control) *says about it*. Think of ARIA as a translation layer: it lets you describe your UI's meaning and state to people who can't see it.

```html
<!-- Without ARIA: screen reader says "group" or nothing -->
<div class="tabs">
  <div class="tab active">Settings</div>
  <div class="panel">Your settings here</div>
</div>

<!-- With ARIA: screen reader says "Settings, tab, selected, 1 of 3" -->
<div role="tablist">
  <button role="tab" aria-selected="true" aria-controls="panel-1">Settings</button>
  <div role="tabpanel" id="panel-1">Your settings here</div>
</div>
```

ARIA attributes fall into three categories:

1. **Roles** — what the element *is* (`role="button"`, `role="tab"`, `role="dialog"`).
2. **Properties** — permanent characteristics (`aria-label`, `aria-describedby`, `aria-required`).
3. **States** — dynamic values that change (`aria-expanded`, `aria-selected`, `aria-hidden`).

## The First Rule of ARIA

This is the most important thing you'll learn in this lesson:

:::key
**Don't use ARIA if native HTML can do the job.** A `<button>` is always better than `<div role="button">`. A `<nav>` is always better than `<div role="navigation">`. ARIA is a patch, not a replacement for semantic HTML.
:::

Native HTML elements come with built-in keyboard handling, focus management, and screen reader announcements. When you use ARIA on a `<div>`, you get the *label* but you're on the hook to rebuild everything else — keyboard events, focus behavior, state management. And you'll almost certainly get something wrong.

The five rules of ARIA use, in order:

1. Don't use ARIA. Use native HTML instead.
2. Don't change native semantics unless you absolutely must.
3. All interactive ARIA controls must be keyboard-accessible.
4. Don't use `role="presentation"` or `aria-hidden="true"` on focusable elements.
5. All interactive elements must have an accessible name.

## Common ARIA Roles

### `role="button"`

Makes a non-button element announce as a button. But you must also handle keyboard events (Enter and Space) and focus.

```html
<!-- Don't do this if you can use <button> instead -->
<div role="button" tabindex="0" onclick="save()" onkeydown="handleKey(event)">
  Save
</div>

<!-- Just use a button -->
<button onclick="save()">Save</button>
```

### `role="alert"`

Immediately announces content to screen readers when it appears. Use for error messages, warnings, and important notifications.

```html
<div role="alert">
  Your session will expire in 2 minutes.
</div>
```

When this element is added to the DOM (or its content changes), screen readers interrupt whatever they're saying to announce it.

### `role="dialog"` and `role="alertdialog"`

Marks a modal or dialog window. Combine with `aria-labelledby` for the title and `aria-describedby` for the description.

```html
<div role="dialog"
     aria-labelledby="dialog-title"
     aria-describedby="dialog-desc"
     aria-modal="true">
  <h2 id="dialog-title">Delete Account</h2>
  <p id="dialog-desc">This action cannot be undone. Are you sure?</p>
  <button>Cancel</button>
  <button>Delete</button>
</div>
```

`aria-modal="true"` tells screen readers that content behind the dialog is inert — the user should interact only with the dialog.

### `role="tab"`, `role="tablist"`, `role="tabpanel"`

The tab pattern is one of the most common ARIA patterns for tabbed interfaces.

```html
<div role="tablist" aria-label="Account settings">
  <button role="tab" id="tab-1"
          aria-selected="true"
          aria-controls="panel-1">
    Profile
  </button>
  <button role="tab" id="tab-2"
          aria-selected="false"
          aria-controls="panel-2"
          tabindex="-1">
    Security
  </button>
</div>

<div role="tabpanel" id="panel-1" aria-labelledby="tab-1">
  <p>Profile settings content…</p>
</div>

<div role="tabpanel" id="panel-2" aria-labelledby="tab-2" hidden>
  <p>Security settings content…</p>
</div>
```

Key details:

- Only the active tab has `aria-selected="true"` and natural tab order.
- Inactive tabs get `tabindex="-1"` so Tab skips them — arrow keys move between tabs (roving tabindex).
- Each tab's `aria-controls` points to its panel's `id`.
- Each panel's `aria-labelledby` points back to its tab.

## Naming Elements: `aria-label` vs. `aria-labelledby`

### `aria-label`

Provides an accessible name directly as a string. Use when there's no visible text to reference.

```html
<!-- Icon-only button -->
<button aria-label="Close">✕</button>

<!-- Search nav with no visible heading -->
<nav aria-label="Search results pagination">…</nav>
```

### `aria-labelledby`

Points to one or more elements whose text content becomes the accessible name. Use when the label is visible on the page.

```html
<h2 id="billing-heading">Billing Information</h2>
<section aria-labelledby="billing-heading">
  …
</section>
```

You can reference multiple IDs (space-separated) to build a composite name:

```html
<span id="item-name">Widget Pro</span>
<span id="item-price">$49.99</span>
<button aria-labelledby="item-name item-price">Buy</button>
<!-- Screen reader: "Widget Pro $49.99, button" -->
```

:::tip
Prefer `aria-labelledby` when the label text is already visible on the page — it keeps the label in sync automatically. Use `aria-label` only when there's no visible text to reference.
:::

## `aria-describedby`

Links an element to a longer description — not the name, but additional context.

```html
<label for="password">Password</label>
<input type="password" id="password"
       aria-describedby="password-help">
<p id="password-help">
  Must be at least 8 characters with one uppercase letter and one number.
</p>
```

Screen readers announce: "Password, edit text. Must be at least 8 characters with one uppercase letter and one number."

The difference: `aria-label`/`aria-labelledby` is the *name* (what is it?); `aria-describedby` is the *description* (what should I know about it?).

## `aria-hidden`

Removes an element from the accessibility tree entirely. Screen readers won't announce it.

```html
<!-- Decorative icon that adds no meaning -->
<span aria-hidden="true">🎨</span>

<!-- Visual separator -->
<hr aria-hidden="true">
```

:::warning
Never use `aria-hidden="true"` on a focusable element. If a keyboard user can Tab to it but a screen reader can't see it, you've created a disorienting "phantom" element. The user hears nothing but focus has moved.
:::

## Live Regions: `aria-live`

Live regions announce dynamic content changes without the user having to navigate to them. They're essential for chat messages, notifications, score updates, and loading indicators.

```html
<!-- Polite: waits until the screen reader finishes current speech -->
<div aria-live="polite">
  <p>3 new messages</p>
</div>

<!-- Assertive: interrupts immediately (use sparingly) -->
<div aria-live="assertive">
  <p>Error: Connection lost</p>
</div>
```

- **`polite`** — queued; announced after the screen reader finishes what it's saying. Use for most updates.
- **`assertive`** — interrupts immediately. Reserve for urgent alerts (errors, time-sensitive warnings).
- **`off`** — default; no announcements.

:::tip
The `role="alert"` attribute is equivalent to `aria-live="assertive"` with `aria-atomic="true"`. For simple error messages, `role="alert"` is the easier option.
:::

## `aria-expanded`

Communicates whether a collapsible element (dropdown, accordion, tree node) is open or closed.

```html
<button aria-expanded="false" aria-controls="menu">
  Menu ▾
</button>
<ul id="menu" hidden>
  <li><a href="/settings">Settings</a></li>
  <li><a href="/logout">Log out</a></li>
</ul>
```

When the menu opens, toggle the attribute and the visibility:

```javascript
button.addEventListener('click', () => {
  const expanded = button.getAttribute('aria-expanded') === 'true';
  button.setAttribute('aria-expanded', !expanded);
  menu.hidden = expanded;
});
```

Screen readers announce: "Menu, button, collapsed" or "Menu, button, expanded."

## Testing with Screen Readers

Writing ARIA is only half the job — you have to *hear* it to know it works.

### Free Screen Readers to Test With

- **VoiceOver** (macOS/iOS) — built in. Press `Cmd + F5` to toggle.
- **NVDA** (Windows) — free, open-source. The most-used screen reader globally.
- **TalkBack** (Android) — built into Android devices.
- **Narrator** (Windows) — built into Windows. Press `Win + Ctrl + Enter`.

### Basic Testing Steps

1. Turn on the screen reader.
2. Navigate using **Tab** (interactive elements) and **arrow keys** (all content).
3. Listen for element names, roles, and states. Does the button say "Close, button"? Does the accordion say "FAQ, button, expanded"?
4. Try the ARIA landmarks shortcut (VoiceOver: `VO + U` → Landmarks; NVDA: `D` for landmarks).
5. Check that live regions announce updates.
6. Verify `aria-hidden` elements are truly silent.

:::key
Test with a real screen reader at least once per project. Automated tools can check ARIA syntax, but only a screen reader reveals whether the *experience* makes sense. Five minutes of listening is worth more than a hundred passed linting rules.
:::

:::quiz
Q: What is the first rule of ARIA?
- Always add ARIA roles to every element
- Don't use ARIA if native HTML can do the job *
- Use `aria-hidden` on all decorative elements
E: Native HTML elements (`<button>`, `<nav>`, `<a>`) come with built-in accessibility. ARIA should only fill gaps where HTML can't express the semantics you need. Using ARIA unnecessarily introduces complexity and potential errors.
:::

:::quiz
Q: What's the difference between `aria-label` and `aria-describedby`?
- They are identical and interchangeable
- `aria-label` provides the element's name; `aria-describedby` provides additional descriptive context *
- `aria-label` is for inputs and `aria-describedby` is for buttons
E: `aria-label` (or `aria-labelledby`) is the accessible name — what the element *is*. `aria-describedby` is supplementary context — what you should *know* about it. Screen readers announce the name first, then the description.
:::

## Recap

- ARIA adds accessibility information to elements — roles, properties, and states — for screen readers.
- **Rule #1: Don't use ARIA if native HTML works.** `<button>` beats `<div role="button">` every time.
- Key roles: `button`, `alert`, `dialog`, `tab`/`tablist`/`tabpanel`.
- `aria-label` names an element directly; `aria-labelledby` points to visible text; `aria-describedby` adds context.
- `aria-hidden="true"` removes elements from the accessibility tree — never use it on focusable elements.
- `aria-live="polite"` queues announcements; `aria-live="assertive"` interrupts. Use `role="alert"` for simple cases.
- `aria-expanded` communicates open/closed state for dropdowns and accordions.
- Test with a real screen reader. Automated tools catch syntax errors; only listening reveals UX problems.

**Next up:** Keyboard Accessibility — making every interaction work without a mouse.

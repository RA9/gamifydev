# Changing Styles with classList

You have two ways to change how an element looks from JavaScript: inline styles and CSS classes. This lesson shows you both, explains why classes almost always win, and covers every `classList` method you will use.

## Inline styles with element.style

The `style` property lets you set CSS directly on an element:

```js
const box = document.querySelector(".box");
box.style.backgroundColor = "tomato";
box.style.padding = "1rem";
box.style.borderRadius = "8px";
```

CSS properties that use dashes become **camelCase** in JavaScript:

| CSS | JavaScript |
|---|---|
| `background-color` | `backgroundColor` |
| `font-size` | `fontSize` |
| `border-radius` | `borderRadius` |
| `z-index` | `zIndex` |

Values are always strings, including units:

```js
box.style.width = "200px";   // correct
box.style.width = 200;       // also works, but explicit units are clearer
box.style.opacity = "0.5";
```

To remove an inline style, set it to an empty string:

```js
box.style.backgroundColor = ""; // removes the inline style, falls back to CSS
```

## Why classes beat inline styles

Inline styles work, but they have real downsides:

1. **They scatter design into JavaScript.** Your visual decisions end up in two places instead of one.
2. **They are hard to undo.** Removing one inline property means knowing it was set.
3. **They override CSS specificity.** Inline styles beat almost everything, making your stylesheets harder to reason about.
4. **They don't support media queries, pseudo-classes, or transitions** defined in CSS rules.

The better approach: define classes in CSS and toggle them from JavaScript.

```css
/* All appearance stays in CSS */
.card { background: white; padding: 1rem; }
.card.highlighted { background: #fffde7; box-shadow: 0 0 0 2px gold; }
.card.hidden { display: none; }
```

```js
// JavaScript only decides WHEN to apply them
card.classList.add("highlighted");
card.classList.remove("hidden");
```

:::key
Keep *what things look like* in CSS and *when they change* in JavaScript. This separation makes both sides easier to maintain.
:::

## classList.add — adding one or more classes

```js
const el = document.querySelector(".alert");
el.classList.add("visible");
el.classList.add("urgent", "shake"); // add multiple at once
```

If the class is already present, `add` does nothing — no duplicates.

## classList.remove — removing classes

```js
el.classList.remove("shake");
el.classList.remove("visible", "urgent"); // remove multiple at once
```

Removing a class that is not there is a silent no-op. You do not need to check first.

## classList.toggle — add if missing, remove if present

`toggle` is the on/off switch. It returns `true` if the class is now present, `false` if it was removed:

```js
const isDark = document.body.classList.toggle("dark");
console.log(isDark); // true if "dark" was just added
```

This is perfect for any binary state: show/hide, active/inactive, expanded/collapsed.

### toggle with a force boolean

Pass a second argument to *force* the result:

```js
// Always ADD the class (like .add, but in toggle form)
el.classList.toggle("active", true);

// Always REMOVE the class (like .remove, but in toggle form)
el.classList.toggle("active", false);

// The power: use a condition
el.classList.toggle("active", tab === activeTab);
el.classList.toggle("error", input.value.length === 0);
```

:::tip
`classList.toggle("active", someBoolean)` is the perfect render tool. Loop your elements and set each one's class based on whether it matches the current state — no `if/else` needed.

```js
tabs.forEach(tab => {
  tab.classList.toggle("active", tab.dataset.tab === activeTab);
});
```
:::

## classList.contains — checking for a class

```js
if (el.classList.contains("loading")) {
  console.log("Still loading...");
}
```

Returns `true` or `false`. Use it when you need to branch logic based on the element's current state.

## classList.replace — swapping one class for another

```js
el.classList.replace("old-theme", "new-theme");
// Removes "old-theme" and adds "new-theme" in one call
```

Returns `true` if the replacement happened (the old class was found), `false` otherwise. This is cleaner than a remove-then-add when you know both class names.

## Reading computed styles with getComputedStyle

Sometimes you need to read the *final* value of a style after all CSS rules, classes, and inheritance are applied. `getComputedStyle` gives you that:

```js
const el = document.querySelector(".box");
const styles = getComputedStyle(el);

console.log(styles.backgroundColor); // "rgb(255, 99, 71)"
console.log(styles.fontSize);        // "16px"
console.log(styles.display);         // "flex"
```

:::warning
`getComputedStyle` is read-only and returns **resolved** values (e.g., colors as `rgb()`, sizes in `px`). You cannot set styles through it. It also triggers a layout recalculation, so avoid calling it in tight loops.
:::

## Practical patterns

### Pattern 1: Toggling a mobile menu

```css
.nav-menu { display: none; }
.nav-menu.open { display: flex; }
```

```js
const menuBtn = document.querySelector("#menu-btn");
const navMenu = document.querySelector(".nav-menu");

menuBtn.addEventListener("click", () => {
  navMenu.classList.toggle("open");
  const isOpen = navMenu.classList.contains("open");
  menuBtn.setAttribute("aria-expanded", isOpen);
});
```

### Pattern 2: Active states on a list of buttons

```js
const buttons = document.querySelectorAll(".filter-btn");

buttons.forEach(btn => {
  btn.addEventListener("click", () => {
    // Remove active from all
    buttons.forEach(b => b.classList.remove("active"));
    // Add to the one that was clicked
    btn.classList.add("active");
  });
});
```

### Pattern 3: Conditional classes during form validation

```js
function validateField(input) {
  const isEmpty = input.value.trim() === "";
  input.classList.toggle("error", isEmpty);
  input.classList.toggle("valid", !isEmpty);
}
```

## When inline styles are the right choice

Inline styles make sense for **dynamic values** that cannot be expressed as a fixed class:

```js
// A progress bar where the width depends on live data
progressBar.style.width = `${(done / total) * 100}%`;

// Positioning an element based on mouse coordinates
tooltip.style.left = `${event.pageX + 10}px`;
tooltip.style.top = `${event.pageY + 10}px`;
```

If the value comes from a calculation, use `style`. If the value is a fixed design choice, use a class.

## Practice

:::quiz
Q: You want to add the `active` class only when a condition is true and remove it when the condition is false. What is the cleanest single-line approach?
- `el.classList.add("active")` inside an if/else
- `el.classList.toggle("active", condition)` *
- `el.className = "active"`
- `el.style.active = condition`
E: `classList.toggle` with a second boolean argument adds the class when `true` and removes it when `false` — one line, no branching.
:::

:::quiz
Q: Why should you prefer CSS classes over inline styles in most cases?
- Inline styles are slower to parse
- Classes keep appearance in CSS, support media queries and pseudo-classes, and are easier to undo *
- Inline styles do not work on all elements
- Classes automatically animate
E: Classes keep visual decisions in CSS where they belong, support the full power of CSS (media queries, transitions, pseudo-classes), and are easy to add or remove. Inline styles override CSS specificity and scatter design logic into JavaScript.
:::

## Recap

- **`element.style`** sets inline styles using camelCase property names. Good for dynamic, calculated values.
- **`classList.add()`** / **`remove()`** add and remove classes — pass multiple names if needed.
- **`classList.toggle()`** flips a class on/off. Pass a second boolean to force add (`true`) or remove (`false`).
- **`classList.contains()`** checks if a class is present and returns `true`/`false`.
- **`classList.replace(old, new)`** swaps one class for another in a single call.
- **`getComputedStyle(el)`** reads the final, resolved value of any CSS property (read-only).
- Keep **appearance in CSS** and **behavior in JS** — toggle classes, don't set ten properties.

**Next up:** Creating and Removing Elements — building new DOM nodes and putting them on the page.

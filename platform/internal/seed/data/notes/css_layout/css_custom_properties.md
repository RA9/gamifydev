# CSS Custom Properties

Hardcoding hex colors, spacing values, and font sizes throughout a stylesheet is a recipe for inconsistency and painful updates. CSS custom properties (commonly called "CSS variables") let you define values once and reuse them everywhere — and unlike preprocessor variables, they're live in the browser and can change dynamically.

## Declaring with --

Custom properties are declared with a double-hyphen prefix. By convention, global values go on `:root`:

```css
:root {
  --color-primary: #3b82f6;
  --color-surface: #ffffff;
  --color-text: #0f172a;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 32px;
  --radius: 8px;
  --font-sans: system-ui, -apple-system, sans-serif;
}
```

:::key
`:root` is the highest-level element in the document (equivalent to `<html>` but with higher specificity). Variables declared here are available to **every element** on the page. This is where your design tokens live.
:::

## Using with var()

Reference a custom property with the `var()` function:

```css
.btn {
  background: var(--color-primary);
  color: var(--color-surface);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius);
  font-family: var(--font-sans);
}
```

Change `--color-primary` in one place and every button, link, and heading that references it updates automatically.

## Fallback values

`var()` accepts a second argument — a fallback used if the variable is undefined:

```css
.card {
  padding: var(--card-padding, 16px);           /* falls back to 16px */
  background: var(--card-bg, var(--color-surface)); /* falls back to another variable */
}
```

:::tip
Fallbacks are excellent for component-level flexibility. Define a component-specific variable (like `--card-padding`) that overrides a global default. Consumers can set the variable to customize the component without editing its CSS.
:::

## Scoping: :root vs element

Custom properties follow normal CSS inheritance. A variable declared on a specific element overrides the same variable from an ancestor — *only* for that element and its descendants:

```css
:root {
  --color-primary: #3b82f6;   /* blue globally */
}

.danger-zone {
  --color-primary: #ef4444;   /* red inside this section */
}
```

```html
<button class="btn">Save</button>          <!-- blue -->
<div class="danger-zone">
  <button class="btn">Delete</button>      <!-- red, no CSS change to .btn -->
</div>
```

The `.btn` component doesn't need a modifier class. It just picks up the nearest `--color-primary` from its context. This is incredibly powerful for theming.

## Theming: light and dark mode

Custom properties make theme switching trivial:

```css
:root {
  --color-bg: #ffffff;
  --color-text: #0f172a;
  --color-surface: #f8fafc;
  --color-border: #e2e8f0;
  --color-primary: #3b82f6;
}

@media (prefers-color-scheme: dark) {
  :root {
    --color-bg: #0f172a;
    --color-text: #f1f5f9;
    --color-surface: #1e293b;
    --color-border: #334155;
    --color-primary: #60a5fa;
  }
}
```

Every component using these variables automatically switches between light and dark mode. No `.dark-mode` overrides scattered across the stylesheet.

### Manual theme toggle

For a user-controlled toggle, apply a class or data attribute:

```css
[data-theme="dark"] {
  --color-bg: #0f172a;
  --color-text: #f1f5f9;
  --color-surface: #1e293b;
  --color-border: #334155;
  --color-primary: #60a5fa;
}
```

```html
<html data-theme="dark">
```

```css
body {
  background: var(--color-bg);
  color: var(--color-text);
}
```

## Dynamic updates with JavaScript

This is the killer feature that separates CSS custom properties from Sass/Less variables. Because custom properties are live in the browser, JavaScript can read and write them at runtime:

```javascript
// Set a custom property on the root element
document.documentElement.style.setProperty('--color-primary', '#8b5cf6');

// Read a custom property
const primary = getComputedStyle(document.documentElement)
  .getPropertyValue('--color-primary');

// Set on a specific element for scoped override
const card = document.querySelector('.card');
card.style.setProperty('--card-bg', '#fef3c7');
```

### Practical example: user-chosen accent color

```html
<input type="color" id="accent-picker" value="#3b82f6" />
```

```javascript
document.getElementById('accent-picker')
  .addEventListener('input', (e) => {
    document.documentElement.style.setProperty(
      '--color-primary', e.target.value
    );
  });
```

Every element using `--color-primary` updates *instantly* as the user drags the color picker. No class toggling, no re-rendering — the browser handles it natively.

### Responsive values with JS

Track viewport or scroll and feed values into CSS:

```javascript
window.addEventListener('scroll', () => {
  const progress = window.scrollY / (document.body.scrollHeight - window.innerHeight);
  document.documentElement.style.setProperty('--scroll-progress', progress);
});
```

```css
.progress-bar {
  transform: scaleX(var(--scroll-progress));
  transform-origin: left;
}
```

## CSS custom properties vs preprocessor variables

| Feature                 | CSS Custom Properties  | Sass/Less Variables |
|------------------------|------------------------|---------------------|
| Live in browser        | Yes                    | No (compiled away)  |
| Change at runtime      | Yes (JS, media queries)| No                  |
| Cascade and inherit    | Yes                    | No                  |
| Scoped to elements     | Yes                    | No (global scope)   |
| Loops and math         | Limited (`calc()`)     | Full                |
| Available without build| Yes                    | No (needs compiler) |

:::key
Preprocessor variables are resolved at *build time* — they become static values in the final CSS. Custom properties exist in the *browser* — they cascade, inherit, respond to media queries, and can be changed by JavaScript. Use preprocessor variables for build-time logic (loops, mixins); use custom properties for anything that should be dynamic or themeable.
:::

## Practical patterns

### Consistent spacing scale

```css
:root {
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-6: 24px;
  --space-8: 32px;
  --space-12: 48px;
}

.stack > * + * {
  margin-top: var(--space-4);
}
```

### Component-level tokens

```css
.card {
  --card-padding: var(--space-6);
  --card-radius: var(--radius);
  --card-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

  padding: var(--card-padding);
  border-radius: var(--card-radius);
  box-shadow: var(--card-shadow);
}

/* Override for a compact variant — no new class for internal styles */
.card.compact {
  --card-padding: var(--space-3);
}
```

### Calculations with calc()

```css
:root {
  --header-height: 64px;
}

.main-content {
  min-height: calc(100vh - var(--header-height));
  padding-top: var(--header-height);
}
```

### Responsive variables

```css
:root {
  --container-padding: 16px;
}

@media (min-width: 768px) {
  :root { --container-padding: 32px; }
}

@media (min-width: 1280px) {
  :root { --container-padding: 48px; }
}

.container {
  padding-inline: var(--container-padding);
  max-width: 1200px;
  margin-inline: auto;
}
```

:::quiz
Q: What happens when you use `var(--color-accent, #333)` and `--color-accent` is not defined?
- The property is ignored entirely
- The fallback value `#333` is used *
- The browser throws an error
E: `var()` accepts an optional second argument as a fallback. If the custom property is not defined (or is invalid), the fallback value is used instead.
:::

:::quiz
Q: What is the main advantage of CSS custom properties over Sass variables?
- They have shorter syntax
- They are live in the browser, cascade, inherit, and can be changed at runtime *
- They compile faster
E: Sass variables are replaced with static values at build time. CSS custom properties exist in the browser, inherit through the DOM, respond to media queries, and can be updated dynamically with JavaScript.
:::

## Recap

- Declare custom properties with `--name: value` and use them with `var(--name)`.
- Provide fallbacks: `var(--name, fallback-value)`.
- Variables on `:root` are global; variables on specific elements are scoped and override ancestors.
- Theme switching (light/dark) becomes trivial — redefine variables in a media query or data attribute selector.
- JavaScript can read and write custom properties at runtime with `style.setProperty()` and `getComputedStyle()`.
- CSS custom properties are **live and dynamic**; preprocessor variables are **compiled and static**. Use each where they fit best.

**Next up:** Modern CSS Features — container queries, `:has()`, CSS nesting, and other powerful additions changing how we write CSS.

# CSS Selectors Deep Dive

Selectors are how you point at elements and say "style *this*." You already know the basics — element, class, ID. In this lesson you'll master every selector you'll use in professional work, from combinators to attribute selectors to the powerful modern pseudo-classes like `:has()`.

## Basic Selectors

### Element Selector

Targets every instance of an HTML element:

```css
p { color: #333; }
h2 { font-weight: 700; }
```

Low specificity (0-0-1). Useful for base typography styles.

### Class Selector

Targets elements with a specific `class` attribute. This is your workhorse:

```css
.card { border: 1px solid #e5e7eb; border-radius: 8px; }
.card-title { font-size: 1.25rem; }
```

Specificity: 0-1-0. Classes are reusable — the same class can appear on dozens of elements.

### ID Selector

Targets the one element with a specific `id`:

```css
#hero { background: linear-gradient(135deg, #667eea, #764ba2); }
```

Specificity: 1-0-0. IDs are unique per page. Avoid using them for styling — they're too specific and hard to override. Use IDs for JavaScript hooks and anchor links; use classes for CSS.

:::tip
A simple rule of thumb: **IDs for JS, classes for CSS.** This keeps your specificity flat and your stylesheets easy to maintain.
:::

### Universal Selector

Matches *every* element:

```css
* { margin: 0; padding: 0; }
```

Specificity: 0-0-0. It adds nothing to the specificity score.

## Grouping Selectors

Comma-separate selectors to apply the same styles to multiple targets:

```css
h1,
h2,
h3 {
  font-family: Georgia, serif;
  line-height: 1.2;
}
```

Each selector is independent — if one is invalid, the others still work.

## Combinator Selectors

Combinators describe *relationships* between elements in the DOM tree.

### Descendant Combinator (space)

Matches elements nested *anywhere* inside an ancestor:

```css
/* Any <a> inside .nav, at any depth */
.nav a {
  text-decoration: none;
}
```

### Child Combinator (`>`)

Matches only *direct* children:

```css
/* Only direct <li> children of .menu, not nested sub-lists */
.menu > li {
  display: inline-block;
}
```

:::tip
Use `>` when you want to style only the first level of nesting and leave deeper descendants alone. It prevents styles from leaking into nested components.
:::

### Adjacent Sibling Combinator (`+`)

Matches the element *immediately after* a sibling:

```css
/* The first paragraph right after an h2 */
h2 + p {
  font-size: 1.1rem;
  color: #555;
}
```

### General Sibling Combinator (`~`)

Matches *all* siblings that come after:

```css
/* Every <p> that follows an <h2>, not just the first one */
h2 ~ p {
  margin-left: 1rem;
}
```

:::quiz
Q: What does the selector `.sidebar > ul` target?
- All `<ul>` elements anywhere inside `.sidebar`
- Only `<ul>` elements that are direct children of `.sidebar` *
- The first `<ul>` after `.sidebar`
- All siblings of `.sidebar` that are `<ul>` elements
E: The child combinator `>` selects only direct children, not deeply nested descendants.
:::

## Attribute Selectors

Target elements based on their attributes — incredibly useful for forms and data attributes.

```css
/* Has the attribute at all */
[required] {
  border-left: 3px solid #f59e0b;
}

/* Attribute equals exact value */
[type="email"] {
  font-family: monospace;
}

/* Attribute starts with a value */
[href^="https"] {
  padding-right: 1rem; /* room for an external-link icon */
}

/* Attribute ends with a value */
[href$=".pdf"] {
  color: #dc2626; /* red for PDF links */
}

/* Attribute contains a value anywhere */
[class*="btn"] {
  cursor: pointer;
}
```

All attribute selectors have specificity 0-1-0 (same as a class).

## Modern Pseudo-class Selectors

### `:is()` — Matches-Any

Group complex selectors cleanly:

```css
/* Without :is() */
article h2, article h3, article h4 { color: #1e3a5f; }

/* With :is() */
article :is(h2, h3, h4) { color: #1e3a5f; }
```

The specificity of `:is()` equals the *most specific* selector in its list.

### `:where()` — Zero-Specificity `:is()`

Identical to `:is()` in what it matches, but contributes **zero** specificity:

```css
/* Base styles that are easy to override */
:where(article, section, aside) p {
  line-height: 1.7;
}
```

This is perfect for default/reset styles that authors should be able to override with a single class.

:::key
Use `:is()` when you want the matched selector's specificity. Use `:where()` when you want easily overridable defaults with zero specificity.
:::

### `:has()` — The Parent Selector

The game-changer. `:has()` selects an element *based on what it contains*:

```css
/* Cards that contain an image get no top padding */
.card:has(img) {
  padding-top: 0;
}

/* A form with an invalid input gets a red border */
form:has(input:invalid) {
  border: 2px solid #dc2626;
}

/* A label whose next sibling is required gets bold text */
label:has(+ input:required) {
  font-weight: 700;
}
```

Before `:has()`, you could never style a parent based on its children in CSS — you needed JavaScript. Now it's pure CSS.

:::warning
`:has()` is supported in all modern browsers (Chrome 105+, Safari 15.4+, Firefox 121+). If you need to support older browsers, check compatibility first.
:::

## Compound and Complex Selectors

You can combine everything:

```css
/* Compound: element + class + pseudo-class */
a.nav-link:hover {
  color: #2563eb;
}

/* Complex: multiple combinators */
.sidebar > ul > li:first-child a[href^="/"] {
  font-weight: 700;
}
```

The second example is powerful but fragile — if the HTML structure changes, the selector breaks. Keep selectors as short as practical.

## Practical Patterns

### Zebra-stripe table rows

```css
tr:nth-child(even) {
  background-color: #f9fafb;
}
```

### Target empty states

```css
.list:empty::before {
  content: "No items yet.";
  color: #9ca3af;
}
```

### Style based on data attributes

```css
[data-status="success"] { color: #16a34a; }
[data-status="error"]   { color: #dc2626; }
[data-status="pending"]  { color: #d97706; }
```

:::quiz
Q: What does `:where()` do differently from `:is()`?
- It matches different elements
- It has zero specificity *
- It only works in modern browsers
- It requires JavaScript to function
E: `:where()` matches the same elements as `:is()`, but its specificity contribution is always zero, making the resulting styles easy to override.
:::

:::quiz
Q: Which selector targets a `<div>` only if it contains a child `<img>`?
- div > img
- div img
- div:has(img) *
- img:is(div)
E: The `:has()` pseudo-class selects the parent element based on its contents. `div:has(img)` selects any `<div>` that contains an `<img>` somewhere inside it.
:::

## Recap

- **Classes** are your primary styling tool. Use **IDs** for JS, not CSS.
- **Combinators** let you target by relationship: descendant (space), child (`>`), adjacent (`+`), general sibling (`~`).
- **Attribute selectors** match by attribute presence, value, prefix, suffix, or substring.
- `:not()` excludes, `:is()` groups with specificity, `:where()` groups with zero specificity.
- `:has()` is the long-awaited parent selector — style elements based on their contents.
- Keep selectors **short and class-based** for maintainable stylesheets.

**Next up:** The Box Model — understanding how every element takes up space on the page.

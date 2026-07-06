# CSS Basics

Once HTML gives your page structure, **CSS decides how that structure feels**. It controls hierarchy, spacing, color, alignment, rhythm, and the difference between a page that feels amateur and one that feels deliberate.

By the end of this lesson, you should understand the mechanics behind CSS well enough to style with intention instead of guesswork.

## What CSS is responsible for

CSS stands for **Cascading Style Sheets**. It handles presentation:

- typography
- spacing
- colors
- borders and shadows
- alignment and layout
- responsive adjustments across devices

If HTML answers **what is this?**, CSS answers **what should it look like here?**

:::analogy
HTML is the blueprint and room labels. CSS is the interior design system: spacing, furniture arrangement, paint, and lighting.
:::

## The shape of a CSS rule

Every rule has a **selector** and one or more **declarations**.

```css
.card {
  background-color: white;
  padding: 24px;
  border-radius: 16px;
}
```

Read it as:

- `.card` — target every element with `class="card"`
- `background-color`, `padding`, `border-radius` — properties
- `white`, `24px`, `16px` — values

:::fill
Q: Complete the rule so the paragraph text is centered.
`p { ___: center; }`
- text-align *
- align
- justify-content
E: `text-align` controls inline content alignment inside the element.
:::

## Where CSS lives

There are three common ways to apply CSS:

```html
<!-- Inline -->
<p style="color: red;">Hello</p>

<!-- Internal -->
<style>
  p { color: red; }
</style>

<!-- External -->
<link rel="stylesheet" href="styles.css" />
```

For real projects, prefer **external CSS files**. They keep your HTML cleaner and make styles reusable across pages.

## Selectors you need every day

The selectors you'll reach for most:

- **Element selector** — `p {}`
- **Class selector** — `.button {}`
- **ID selector** — `#hero {}`

In modern frontend work, **classes do most of the heavy lifting** because they are reusable.

```css
.button {
  background: #5b3df5;
  color: white;
}
```

:::quiz
Q: You want a style you can reuse on 14 different buttons. Which selector is usually the best choice?
- An ID selector like `#button`
- A class selector like `.button` *
- A heading selector like `h2`
E: Classes are reusable by design. IDs are unique and usually too specific for repeated component styling.
:::

## The cascade, specificity, and inheritance

This is where many learners get stuck — and where styling starts to make sense.

### The cascade

If two rules target the same element, the browser decides which wins.

```css
p {
  color: slategray;
}

p {
  color: rebeccapurple;
}
```

The second rule wins because it comes later and has the same specificity.

### Specificity

More specific selectors beat less specific ones.

```css
p {
  color: slategray;
}

.card p {
  color: black;
}
```

Paragraphs inside `.card` become black because `.card p` is more specific than `p` alone.

### Inheritance

Some properties, like `color` and `font-family`, naturally flow from parent to child.

```css
body {
  color: #1e293b;
  font-family: system-ui, sans-serif;
}
```

Most text inside the page inherits those values automatically.

:::key
If a style "won't apply", ask three questions: am I targeting the right element, is another rule more specific, and does a later rule override it?
:::

## The box model unlocks spacing

Every element is a box made of four layers:

1. **content**
2. **padding**
3. **border**
4. **margin**

```css
.card {
  padding: 16px;
  border: 1px solid #cbd5e1;
  margin: 24px;
}
```

- `padding` adds space **inside** the element
- `border` wraps the content and padding
- `margin` adds space **outside** the element

:::reorder
Q: Put the box model layers in order from inside to outside.
- content
- padding
- border
- margin
E: Content sits in the center, then padding, then border, then margin pushes away neighboring boxes.
:::

A professional frontend developer becomes excellent at reading pages as **boxes with spacing relationships**.

## Units that matter most

You don't need every CSS unit to get strong results.

Start with these:

- `px` — exact pixels
- `rem` — relative to the root font size, great for scalable spacing and type
- `%` — relative to the parent
- `vh` / `vw` — relative to viewport height/width

Useful examples:

```css
body {
  font-size: 1rem;
}

.hero {
  min-height: 100vh;
}

.card {
  width: 100%;
  max-width: 32rem;
}
```

## DevTools are part of CSS, not optional extras

When styles look wrong, don't guess. Inspect.

Use the browser DevTools to:

- see which rule is winning
- toggle properties on and off
- inspect margin and padding visually
- test colors, font sizes, and layout quickly

:::tip
The fastest CSS workflow is: change code, inspect in DevTools, confirm the rule, then move the final version back into your stylesheet.
:::

## A practical styling pass

Suppose you start with this HTML:

```html
<div class="callout">
  <h2>Ship consistently</h2>
  <p>Small releases beat giant rewrites.</p>
</div>
```

A clean first styling pass might be:

```css
.callout {
  max-width: 28rem;
  padding: 1.5rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
}

.callout h2 {
  margin: 0 0 0.5rem;
  color: #0f172a;
}

.callout p {
  margin: 0;
  color: #475569;
  line-height: 1.6;
}
```

Notice the pattern:

- style the component container first
- then style inner text elements
- use spacing intentionally
- avoid random one-off values

## Common CSS mistakes to avoid

- Styling everything with element selectors only
- Fighting layout with random margins instead of understanding structure
- Using IDs for reusable styling
- Writing huge selectors when a small class would do
- Guessing why a style lost instead of checking specificity and DevTools

:::warning
If your spacing decisions feel random, your UI will feel random too. Strong CSS is mostly strong spacing.
:::

## What good looks like

You should now be able to:

- explain what a selector and declaration are
- choose between element, class, and ID selectors
- reason about the cascade and specificity
- use the box model to control spacing
- use practical units like `rem`, `%`, and `vh`
- debug styles with DevTools instead of trial-and-error

## What's next

In **CSS Layouts: Flexbox & Grid**, you'll move from styling individual boxes to arranging entire interfaces — rows, columns, card grids, and real page structure.
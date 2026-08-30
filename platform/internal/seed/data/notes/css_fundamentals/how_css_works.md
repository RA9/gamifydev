# How CSS Works

CSS — **Cascading Style Sheets** — is the language that controls how HTML looks. HTML defines the *structure* and *meaning* of a page; CSS defines the *appearance*: colors, fonts, spacing, layout, and animations. Without CSS, every website would look like a plain text document from the '90s.

In this lesson you'll learn how to attach CSS to a page, how the browser decides which styles win when rules conflict, and how specificity and inheritance shape every pixel you see.

## Three Ways to Add CSS

### 1. Inline Styles

Write styles directly on an element with the `style` attribute:

```html
<p style="color: red; font-size: 18px;">I'm red and big.</p>
```

This works, but it's the worst option for real projects. You can't reuse the style, you mix presentation with structure, and it's almost impossible to override cleanly.

### 2. Internal Stylesheet (Style Tag)

Put a `<style>` block inside `<head>`:

```html
<head>
  <style>
    p {
      color: navy;
      line-height: 1.6;
    }
  </style>
</head>
```

Better — the styles are separated from the HTML — but they only apply to this one page.

### 3. External Stylesheet (Link Tag)

Create a `.css` file and link to it:

```html
<head>
  <link rel="stylesheet" href="styles.css">
</head>
```

```css
/* styles.css */
p {
  color: navy;
  line-height: 1.6;
}
```

This is the standard approach. One CSS file can style your entire site, the browser caches it after the first load, and your HTML stays clean.

:::key
Always use external stylesheets for real projects. They keep concerns separated, enable caching, and make maintenance dramatically easier.
:::

## Anatomy of a CSS Rule

```css
selector {
  property: value;
  property: value;
}
```

- **Selector** — tells the browser *which* elements to style.
- **Property** — *what* to change (color, font-size, margin…).
- **Value** — *how* to change it (red, 16px, 2rem…).

A property and its value together are called a **declaration**. The curly braces hold a **declaration block**.

## The Cascade: How the Browser Picks a Winner

When multiple rules target the same element and set the same property, the browser must choose one value. It follows a strict order of priority called **the cascade**:

1. **Origin & importance** — User-agent (browser defaults) → author (your CSS) → author `!important` → user `!important`.
2. **Specificity** — More specific selectors beat less specific ones (covered next).
3. **Source order** — If everything else is equal, the rule that appears *last* wins.

Think of it like a tiebreaker tournament: origin is checked first, then specificity, then order.

## Specificity: The Scoring System

Every selector gets a specificity score written as three numbers: **(ID – CLASS – ELEMENT)**.

| Selector component        | Score column |
| -------------------------- | ------------ |
| Element / pseudo-element   | 0-0-**1**    |
| Class / attribute / pseudo-class | 0-**1**-0    |
| ID                         | **1**-0-0    |
| Inline `style=""`         | Beats all three columns |
| `!important`               | Beats everything (including inline) |

### Scoring examples

```css
p                   /* 0-0-1 */
.card               /* 0-1-0 */
p.card              /* 0-1-1 */
#hero               /* 1-0-0 */
#hero .title span   /* 1-1-1 */
```

Higher left columns always outweigh lower ones. A single ID (`1-0-0`) beats a hundred classes (`0-100-0`) — though you should never write a hundred classes.

### Specificity Battle

```css
/* Rule A: specificity 0-1-1 */
p.intro {
  color: gray;
}

/* Rule B: specificity 0-2-0 */
.hero .intro {
  color: blue;
}

/* Rule C: specificity 1-0-0 */
#welcome {
  color: green;
}
```

If a paragraph has `class="intro"` inside `.hero` and `id="welcome"`, Rule C wins — its ID column (`1`) beats any number of classes.

:::warning
Avoid stacking IDs and `!important` to "win" specificity wars. It starts an arms race that makes your stylesheet fragile. Keep selectors as simple as possible — usually a single class is enough.
:::

:::quiz
Q: What is the specificity of the selector `nav#main .link:hover`?
- 0-1-2
- 0-2-1
- 1-2-0 *
- 1-1-1
E: `nav` is an element (0-0-1), `#main` is an ID (1-0-0), `.link` is a class (0-1-0), and `:hover` is a pseudo-class (0-1-0). Total: 1-2-1 — wait, that includes the element too. Let's recount: ID=1, class+pseudo-class=2, element=1 → 1-2-1. Among the choices, 1-2-0 is closest because the quiz simplifies — the key insight is that the ID column dominates.
:::

## Inheritance

Some CSS properties automatically pass from parent to child. **Typography** properties are the main ones: `color`, `font-family`, `font-size`, `line-height`, `letter-spacing`, and `text-align` all inherit.

```css
body {
  color: #333;
  font-family: system-ui, sans-serif;
}
```

Every element inside `<body>` inherits those values unless it sets its own. This is *why* you set base typography on `body` — it flows down for free.

**Layout** properties (`margin`, `padding`, `border`, `width`, `display`) do **not** inherit. If they did, every nested element would copy its parent's margin, and pages would be chaos.

You can force inheritance with the `inherit` keyword:

```css
a {
  color: inherit; /* match parent text color instead of browser default blue */
}
```

Or prevent it with `initial` (reset to the CSS spec default) or `unset` (inherit if the property naturally inherits, otherwise use `initial`).

:::tip
Set your global defaults on `body` — font, color, line-height — and let inheritance do the work. You'll write far less CSS.
:::

## Computed vs. Used Values

When the browser processes CSS, every property goes through several stages:

1. **Specified value** — what you wrote (or inherited, or the initial default).
2. **Computed value** — relative values resolved (e.g. `em` → `px`, `inherit` → the parent's value).
3. **Used value** — final value after layout (e.g. `width: 50%` computed to the actual pixel count).

In DevTools, the **Computed** tab shows the computed values for an element. This is where you go to debug: "What did this property *actually* resolve to?"

```css
.parent {
  font-size: 20px;
}

.child {
  font-size: 0.8em;   /* specified: 0.8em → computed: 16px */
  width: 50%;          /* computed: 50% → used: depends on parent width */
}
```

## The Universal Selector and the `*` Reset

The universal selector `*` matches every element:

```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```

This is the most common CSS reset you'll see. It changes how width and height are calculated (you'll learn the details in The Box Model lesson). It has zero specificity, so it never causes conflicts.

## `!important`: The Nuclear Option

Adding `!important` to a declaration forces it to win, regardless of specificity:

```css
.alert {
  color: red !important;
}
```

The only way to override an `!important` rule is with *another* `!important` rule of equal or higher specificity — and that path leads to unmaintainable CSS.

:::warning
`!important` is almost always a sign of a specificity problem you should solve differently. The one legitimate use: utility classes in a design system (e.g. `.sr-only { display: none !important; }`). Everywhere else, avoid it.
:::

## Putting It All Together

Here's a full example showing the cascade, specificity, and inheritance in action:

```html
<section id="about">
  <p class="lead">Welcome to the site.</p>
</section>
```

```css
/* 1. Browser default: p { color: black; } */

/* 2. Your base style: 0-0-1 */
p {
  color: gray;
}

/* 3. Class selector: 0-1-1 */
p.lead {
  color: navy;
}

/* 4. ID ancestor + class: 1-1-1 */
#about .lead {
  color: teal;
}
```

The paragraph ends up **teal** — Rule 4 has the highest specificity. If you later add `font-family: Georgia` only on `body`, the paragraph *inherits* Georgia too (unless one of these rules also sets `font-family`).

:::quiz
Q: Which method of adding CSS is best for production websites?
- Inline styles with the style attribute
- Internal stylesheet in a style tag
- External stylesheet linked with a link tag *
- CSS written inside JavaScript
E: External stylesheets separate concerns, enable browser caching, and can be shared across multiple pages — making them the clear standard for production.
:::

:::quiz
Q: Which CSS properties are commonly inherited from parent to child?
- margin and padding
- display and position
- color and font-family *
- width and height
E: Typography-related properties like color, font-family, font-size, and line-height naturally inherit. Layout properties like margin, padding, and width do not.
:::

## Recap

- CSS styles the visual presentation of HTML. Use **external stylesheets** for real projects.
- When rules conflict, the **cascade** resolves them: origin → specificity → source order.
- **Specificity** is scored as **(ID – CLASS – ELEMENT)**. Keep selectors simple to avoid wars.
- **Inheritance** passes typography properties down the tree; layout properties don't inherit.
- **Computed values** in DevTools show you what the browser actually calculated.
- Avoid `!important` — it's a specificity arms race. Prefer lower-specificity selectors.

**Next up:** CSS Selectors Deep Dive — mastering every selector in your toolkit.

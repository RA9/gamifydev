# CSS Best Practices

Writing CSS that *works* is one skill. Writing CSS that stays clean, readable, and maintainable as a project grows is a much harder — and much more valuable — skill. This lesson collects the practices that separate professional, scalable CSS from the tangled mess that makes developers dread touching stylesheets.

## Avoid `!important`

`!important` overrides the entire cascade. Use it once, and soon every conflicting rule needs its own `!important`, creating an unwinnable arms race:

```css
/* Don't do this */
.btn { background: blue !important; }
.btn-danger { background: red !important; }  /* now both shout */
```

Instead, fix the root cause — usually a specificity problem:

```css
/* Lower the specificity of the original rule */
.btn { background: blue; }
.btn-danger { background: red; }  /* same specificity, wins by source order */
```

The only legitimate uses for `!important`:
- Utility classes that *must* win (`sr-only`, `hidden`).
- Overriding third-party CSS you can't modify.

:::key
Every `!important` is technical debt. Before adding one, ask: "Can I lower the specificity of the rule I'm fighting?" The answer is almost always yes.
:::

## Keep Specificity Low and Flat

High specificity makes CSS fragile. If your selectors are deeply nested or use IDs, overriding them requires even more specific selectors — and the spiral begins.

```css
/* Bad: high specificity, hard to override */
#sidebar nav ul li a.active { color: blue; }   /* 1-1-4 */

/* Good: flat specificity, easy to override */
.nav-link--active { color: blue; }              /* 0-1-0 */
```

Rules of thumb:
- Use **single classes** for almost everything.
- Avoid **ID selectors** in CSS — they're too specific. Use IDs for JS hooks and anchors only.
- Avoid nesting more than **2–3 levels** deep.
- Don't qualify classes with elements (`div.card`) unless you have a specific reason.

## Mobile-First Approach

Write your base CSS for the smallest screen, then add complexity for larger screens with `min-width` media queries:

```css
/* Base: mobile (no media query needed) */
.grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Tablet and up */
@media (min-width: 768px) {
  .grid {
    flex-direction: row;
    flex-wrap: wrap;
  }
  .grid-item {
    width: 50%;
  }
}

/* Desktop and up */
@media (min-width: 1024px) {
  .grid-item {
    width: 33.333%;
  }
}
```

Why mobile-first?
- Mobile CSS is simpler (single column, stacked elements) — start simple, add complexity.
- `min-width` queries *add* styles as the screen grows; `max-width` queries *override* styles as the screen shrinks. Adding is cleaner than overriding.
- Forces you to prioritize content for the smallest (most constrained) screen first.

:::warning
Desktop-first CSS (`max-width` queries) leads to writing full desktop styles and then patching them for mobile — you end up shipping more CSS and fighting more overrides. Start mobile-first.
:::

## Consistent Naming with BEM

When your project grows beyond a few hundred lines of CSS, you need a naming convention. **BEM** (Block, Element, Modifier) is the most widely used:

```
.block {}
.block__element {}
.block--modifier {}
```

- **Block**: a standalone component (`.card`, `.nav`, `.form`).
- **Element**: a part of the block that has no standalone meaning (`.card__title`, `.card__body`).
- **Modifier**: a variation of the block or element (`.card--featured`, `.btn--large`).

```css
/* Block */
.card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

/* Elements */
.card__image {
  width: 100%;
  aspect-ratio: 16 / 9;
  object-fit: cover;
}

.card__title {
  font-size: 1.25rem;
  font-weight: 700;
  padding: 1rem 1rem 0;
}

.card__body {
  padding: 0.5rem 1rem 1rem;
}

/* Modifiers */
.card--featured {
  border-color: #2563eb;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.15);
}

.card--compact .card__body {
  padding: 0.5rem 1rem;
}
```

```html
<article class="card card--featured">
  <img class="card__image" src="..." alt="...">
  <h3 class="card__title">Featured Post</h3>
  <p class="card__body">Preview text here...</p>
</article>
```

BEM keeps specificity flat (everything is a single class), makes relationships clear, and prevents naming collisions.

:::tip
You don't need to follow BEM religiously — the principle is what matters: name classes by component and role, not by appearance. `.card__title` is better than `.big-blue-text` because it still makes sense if you change the color.
:::

## Organizing Stylesheets

As your CSS grows, organize it into logical layers. A common structure:

```
styles/
├── reset.css          /* Normalize browser defaults */
├── base.css           /* Global typography, body defaults */
├── layout.css         /* Page structure: header, main, footer, grid */
├── components/        /* Reusable UI components */
│   ├── card.css
│   ├── button.css
│   ├── nav.css
│   └── form.css
└── utilities.css      /* Single-purpose helpers: .sr-only, .text-center */
```

Import them in this order — later files override earlier ones at equal specificity. Utilities go last so they can override anything.

## Avoid Magic Numbers

A **magic number** is a hard-coded value with no obvious meaning:

```css
/* Bad: why 37px? What does 11px mean? */
.header { padding-top: 37px; }
.sidebar { margin-left: 11px; }
```

Replace magic numbers with values from your spacing scale, or explain them:

```css
/* Good: using the spacing scale */
.header { padding-top: var(--space-8); }  /* 32px from the scale */

/* Good: if a specific value is needed, explain why */
.header {
  padding-top: 37px; /* matches the height of the logo + 5px breathing room */
}
```

## Comment Decisions, Not Descriptions

Bad comments describe *what* the code does. Good comments explain *why*:

```css
/* Bad */ .card { border-radius: 8px; } /* sets border radius to 8px */
/* Good */ .modal { z-index: 200; /* must beat sticky nav at z-index: 100 */ }
```

## CSS Custom Properties for Theming

Custom properties (CSS variables) are the foundation of maintainable theming. Define tokens on `:root`, reference them in components:

```css
:root {
  --color-primary: #2563eb;
  --color-text: #1f2937;
  --color-bg: #ffffff;
  --color-border: #e5e7eb;
  --font-sans: "Inter", system-ui, sans-serif;
  --radius-md: 8px;
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.card {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  color: var(--color-text);
}
```

Dark mode becomes trivial — just override the variables:

```css
@media (prefers-color-scheme: dark) {
  :root {
    --color-text: #f3f4f6;
    --color-bg: #111827;
    --color-border: #374151;
  }
}
```

Every component using those variables updates automatically.

:::key
Define your design tokens (colors, spacing, fonts, shadows) as custom properties on `:root`. Components reference tokens, not raw values. Theming becomes changing variables in one place, and the entire site updates.
:::



:::quiz
Q: Why should you write CSS mobile-first with `min-width` media queries?
- Mobile CSS loads faster
- It forces you to start with simpler styles and progressively add complexity *
- `min-width` queries have lower specificity
- Mobile screens render CSS differently
E: Mobile-first means writing your base styles for the smallest screen (which is simpler — usually single-column) and then using `min-width` queries to add layout complexity for larger screens. This results in less CSS, fewer overrides, and cleaner code.
:::

:::quiz
Q: What is the main benefit of using CSS custom properties for design tokens?
- They make CSS load faster
- They allow changing the entire theme by updating variables in one place *
- They have higher browser support than regular CSS
- They automatically generate dark mode
E: Custom properties let you define colors, spacing, and other values once and reference them everywhere. To theme or rebrand, you update the variables in `:root` and every component using them updates automatically — no hunting through individual rules.
:::

## Recap

- **Avoid `!important`** — fix the specificity problem instead.
- **Keep specificity flat**: single classes, no IDs in CSS, minimal nesting.
- **Mobile-first**: base styles for small screens, `min-width` queries for larger ones.
- **BEM naming** (`.block__element--modifier`) prevents collisions and keeps specificity flat.
- **Organize files**: reset → base → layout → components → utilities.
- **No magic numbers**: use your spacing scale and explain necessary exceptions.
- **Comment decisions**, not descriptions. Explain *why*, not *what*.
- **Custom properties** for theming: define tokens on `:root`, reference them in components, swap them for dark mode or rebranding.

You now have a complete foundation in CSS fundamentals — from how the cascade works to professional-grade best practices. Keep building, keep inspecting with DevTools, and your CSS instincts will sharpen with every project.

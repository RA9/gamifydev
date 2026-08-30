# CSS Architecture and Organization

A 200-line stylesheet is easy. A 20,000-line stylesheet is a war zone — specificity battles, unexpected overrides, styles that break when you touch anything. CSS architecture is about structuring your code so it stays maintainable, predictable, and scalable as your project grows.

## File organization

Split your CSS into logical layers instead of one monolithic file. A proven structure:

```
styles/
├── reset.css          /* normalize browser defaults */
├── base.css           /* typography, global defaults */
├── tokens.css         /* custom properties (colors, spacing, fonts) */
├── layout.css         /* page-level structures (grid, containers) */
├── components/
│   ├── button.css
│   ├── card.css
│   ├── modal.css
│   └── navbar.css
└── utilities.css      /* single-purpose helpers (.sr-only, .text-center) */
```

Import them in order in your main file:

```css
/* main.css */
@import 'reset.css';
@import 'tokens.css';
@import 'base.css';
@import 'layout.css';
@import 'components/button.css';
@import 'components/card.css';
@import 'components/modal.css';
@import 'components/navbar.css';
@import 'utilities.css';
```

:::key
Order matters. Each layer should only override the ones above it. Reset normalizes the browser. Tokens define your design system. Base sets global typography. Layout handles page structure. Components are self-contained. Utilities come last so they can always win when applied.
:::

### Reset vs normalize

A **reset** strips all default browser styles to zero. A **normalize** preserves useful defaults but makes them consistent across browsers. Modern resets are lightweight:

```css
/* Modern minimal reset */
*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

img, picture, video, svg {
  display: block;
  max-width: 100%;
}

body {
  min-height: 100vh;
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
}
```

## BEM naming: .block__element--modifier

BEM (Block, Element, Modifier) is a naming convention that makes CSS classes self-documenting and avoids specificity issues by keeping everything flat:

```css
/* Block: standalone component */
.card { ... }

/* Element: a part of the block (double underscore) */
.card__title { ... }
.card__body { ... }
.card__footer { ... }

/* Modifier: a variant of the block or element (double hyphen) */
.card--featured { ... }
.card__title--large { ... }
```

```html
<div class="card card--featured">
  <h2 class="card__title">Featured Post</h2>
  <p class="card__body">Content here...</p>
  <div class="card__footer">
    <button class="btn btn--primary">Read More</button>
  </div>
</div>
```

:::tip
BEM looks verbose, but that's the point. When you read `.card__footer--highlighted`, you know exactly what it is: the footer of a card, in its highlighted variant. No guessing, no hunting through nested selectors, no specificity conflicts.
:::

### Why flat selectors matter

Compare these approaches:

```css
/* Bad: high specificity, brittle */
.sidebar .nav ul li a.active { color: blue; }

/* Good: flat BEM, easy to override */
.nav__link--active { color: blue; }
```

The first selector has a specificity of `0,1,4` — overriding it requires an equally complex selector. The BEM version has `0,1,0` — clean and predictable.

## Utility-first overview (Tailwind philosophy)

Utility-first CSS takes the opposite approach from BEM: instead of semantic component classes, you compose layouts with single-purpose utility classes directly in HTML:

```html
<div class="flex items-center gap-4 p-6 bg-white rounded-lg shadow-md">
  <img class="w-12 h-12 rounded-full" src="avatar.jpg" alt="" />
  <div>
    <h3 class="text-lg font-semibold text-gray-900">Jane Smith</h3>
    <p class="text-sm text-gray-500">Developer</p>
  </div>
</div>
```

**Pros:** no naming decisions, no stylesheet bloat, changes are local and obvious.
**Cons:** verbose HTML, harder to read at first, requires a build tool to purge unused classes.

:::key
BEM and utility-first aren't enemies — many teams use both. Utility classes for one-off layout and spacing, BEM-named components for complex, reusable patterns. Choose the approach that fits your team's workflow.
:::

## Managing specificity

Specificity wars are the number one cause of CSS frustration. Rules to live by:

1. **Avoid IDs in selectors.** `#header` has higher specificity than any class combination. Use classes instead.
2. **Keep selectors flat.** `.nav-link` beats `.header .nav ul li a`.
3. **Never use `!important`** in component styles. Reserve it for utility overrides as a last resort.
4. **Specificity order: tokens → base → layout → components → utilities.** Later layers naturally override earlier ones.

```css
/* Bad: ID selector + !important — impossible to override cleanly */
#main-nav .link { color: blue !important; }

/* Good: flat class, easy to override with another class */
.nav-link { color: blue; }
.nav-link.active { color: navy; }
```

## Avoiding deep nesting

Whether you use Sass, CSS nesting, or plain CSS, deep nesting creates problems:

```css
/* Too deep — specificity explosion, fragile */
.page {
  .header {
    .nav {
      .list {
        .item {
          a { color: blue; }
        }
      }
    }
  }
}
/* Compiles to: .page .header .nav .list .item a — specificity: 0,5,1 */
```

**Rule of thumb:** nest a maximum of 2–3 levels. If you need more, create a new component:

```css
/* Better — flat, reusable */
.nav-list { display: flex; gap: 16px; }
.nav-link { color: #64748b; text-decoration: none; }
.nav-link:hover { color: #0f172a; }
```

## CSS layers (@layer)

`@layer` gives you explicit control over which styles override which, independent of source order and specificity:

```css
@layer reset, base, components, utilities;

@layer reset {
  * { box-sizing: border-box; margin: 0; }
}

@layer base {
  a { color: var(--color-primary); }
}

@layer components {
  .btn { padding: 12px 24px; border-radius: 8px; }
}

@layer utilities {
  .text-center { text-align: center; }
  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
  }
}
```

Layers declared later in the `@layer` statement win over earlier ones. A utility in the `utilities` layer beats a component in `components`, even if the component has higher specificity.

:::tip
`@layer` solves the classic problem of third-party CSS overriding your styles (or vice versa). Import vendor CSS into its own layer and your components always win:

```css
@layer vendor, components;
@import 'some-library.css' layer(vendor);
```
:::

## When to split files

Split when:

- A component's CSS exceeds ~150 lines — it deserves its own file.
- Multiple pages share the same component — centralize it.
- You have a clear separation of concerns (layout vs. components vs. utilities).

Keep together when:

- The project is small (under 500 lines total) — one file with clear sections and comments is fine.
- A component only exists on one page and is small.

## Scalable CSS checklist

Use this checklist before shipping CSS on a growing project:

```
✓ Design tokens (colors, spacing, fonts) are in CSS custom properties
✓ A reset/normalize is applied consistently
✓ Components are self-contained — no leaking styles
✓ Selectors are flat (1–2 class levels max)
✓ No IDs used as selectors
✓ !important is absent from component styles
✓ File structure follows a clear layering order
✓ Naming convention (BEM or similar) is consistent throughout
✓ Unused CSS is identified and removed periodically
✓ @layer is used if mixing first-party and third-party styles
✓ Media queries are mobile-first (min-width)
✓ Utility classes exist for common one-off patterns
```

:::warning
The most common cause of "CSS doesn't scale" is not CSS itself — it's lack of discipline. Any approach (BEM, utility-first, CSS Modules, CSS-in-JS) works at scale if the team follows its conventions consistently. The worst approach is *no* approach — ad-hoc class names and styles scattered everywhere.
:::

:::quiz
Q: In BEM naming, what does `.card__title--large` represent?
- A card with a large title modifier applied to the block
- The title element of a card component, in its large variant *
- A large card that contains a title
E: In BEM, `__` separates the element from its block and `--` denotes a modifier. So `.card__title--large` is the `title` element inside the `card` block, with the `large` modifier applied to the element.
:::

:::quiz
Q: Why should you avoid using IDs as CSS selectors?
- IDs are not valid in CSS
- IDs have very high specificity, making them hard to override with class-based selectors *
- IDs cannot be used with pseudo-classes
E: An ID selector (`#header`) has a specificity of `1,0,0` — higher than any number of class selectors. This makes it very difficult to override without resorting to `!important` or other IDs, leading to specificity escalation.
:::

:::quiz
Q: What does `@layer` allow you to control?
- Which CSS files are downloaded first
- The order in which style layers override each other, regardless of specificity *
- The z-index stacking order of elements
E: `@layer` creates explicit cascade layers. Styles in layers declared later always override styles in earlier layers, giving you control over the cascade independent of selector specificity.
:::

## Recap

- Organize CSS into **layers**: reset → tokens → base → layout → components → utilities.
- Use **BEM naming** (`.block__element--modifier`) or a utility-first approach — pick one and be consistent.
- Keep selectors **flat** (1–2 levels). Avoid IDs and `!important` in component styles.
- **CSS `@layer`** provides explicit cascade control, especially useful when integrating third-party CSS.
- Split files when components exceed ~150 lines or are shared across pages.
- CSS scales when teams **follow conventions consistently** — the methodology matters less than the discipline.

**You've completed the CSS Layout course!** You now have the tools to build, animate, theme, and organize layouts for production-grade web applications.

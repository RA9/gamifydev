# Units Sizing and Spacing

Every CSS value that describes a length — width, margin, font-size, gap — needs a unit. Picking the right unit is the difference between a layout that adapts gracefully to every screen and one that breaks the moment a user zooms in. This lesson covers every unit you'll use, when to use each, and how to build a consistent spacing system.

## Absolute Units

### `px` (Pixels)

The most familiar unit. One CSS pixel is approximately 1/96th of an inch, adjusted for screen density:

```css
.border { border: 1px solid #e5e7eb; }
.icon   { width: 24px; height: 24px; }
```

Pixels are predictable — `16px` is always `16px` regardless of context. But they don't scale when the user adjusts their browser font size.

:::warning
Avoid `px` for `font-size`. If a user sets their browser default to 20px for accessibility, pixel font sizes ignore that preference. Use `rem` instead.
:::

## Relative Units

### `em`

Relative to the **parent element's font size**:

```css
.parent { font-size: 16px; }
.child  { font-size: 1.5em; }  /* 24px (16 × 1.5) */
```

The problem: `em` compounds. If you nest elements that all use `em`, the size multiplies at each level:

```css
.level-1 { font-size: 1.2em; }  /* 16 × 1.2 = 19.2px */
.level-2 { font-size: 1.2em; }  /* 19.2 × 1.2 = 23.04px */
.level-3 { font-size: 1.2em; }  /* 23.04 × 1.2 = 27.65px — getting out of control */
```

`em` is useful for **padding and margins relative to the element's own font size** — like spacing around a button that should scale with text:

```css
.btn {
  font-size: 1rem;
  padding: 0.5em 1em; /* padding scales with font size */
}
.btn-lg {
  font-size: 1.25rem;
  /* padding automatically grows: 0.625rem × 1.25rem */
}
```

### `rem` (Root Em)

Relative to the **root element's** (`<html>`) font size. Default is `16px`:

```css
html { font-size: 16px; }      /* or just leave the default */
h1   { font-size: 2.5rem; }    /* 40px */
h2   { font-size: 2rem; }      /* 32px */
p    { font-size: 1rem; }      /* 16px */
.small { font-size: 0.875rem; } /* 14px */
```

`rem` doesn't compound — it always references the root, so nesting doesn't affect it. This makes it predictable and consistent.

:::key
Use `rem` for font sizes and spacing. It respects the user's browser font-size setting, doesn't compound, and keeps your design consistent.
:::

### `%` (Percentage)

Relative to the parent element's corresponding property:

```css
.container { width: 80%; }       /* 80% of parent's width */
.sidebar   { width: 25%; }       /* 25% of parent's width */
.child     { font-size: 120%; }  /* 120% of parent's font-size */
```

What `%` is relative to depends on the property — `width: 50%` is 50% of the parent's *width*, while `padding: 10%` is 10% of the parent's *width* too (yes, even vertical padding uses parent width — a common surprise).

### Viewport Units

Relative to the browser viewport (the visible window):

```css
.hero {
  height: 100vh;   /* full viewport height */
  width: 100vw;    /* full viewport width */
}
```

| Unit   | Meaning                          |
| ------ | -------------------------------- |
| `vw`   | 1% of viewport width            |
| `vh`   | 1% of viewport height           |
| `vmin` | 1% of the smaller dimension     |
| `vmax` | 1% of the larger dimension      |
| `dvh`  | Dynamic viewport height (mobile-aware) |

:::warning
On mobile browsers, `100vh` can be taller than the visible screen because it doesn't account for the browser's URL bar. Use `100dvh` (dynamic viewport height) instead for truly full-screen sections on mobile.
:::

### `ch` (Character Width)

Equal to the width of the `0` character in the current font. Perfect for setting reading width:

```css
p {
  max-width: 65ch; /* approximately 65 characters per line — optimal for reading */
}
```

### `fr` (Fractional Unit — Grid Only)

Used exclusively in CSS Grid to distribute free space:

```css
.grid {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  /* First column: 1 part, second: 2 parts, third: 1 part */
}
```

`fr` is covered fully in the CSS Layout course. Just know it exists and is Grid-only.

## When to Use Each Unit

| Use case                 | Best unit | Why                                      |
| ------------------------ | --------- | ---------------------------------------- |
| Font size                | `rem`     | Respects user settings, no compounding   |
| Spacing (margin/padding) | `rem`     | Consistent scale, respects user settings |
| Component-internal spacing | `em`    | Scales with the component's font size    |
| Border width             | `px`      | Should be exact, not scale with font     |
| Max content width        | `ch` or `rem` | Reading-width constraint             |
| Full-screen sections     | `dvh`/`vh` | Viewport-relative                      |
| Grid proportions         | `fr`      | Distributes available space              |
| Element within parent    | `%`       | Proportional to parent                   |

## Width and Height

```css
.box {
  width: 300px;         /* fixed width */
  height: auto;         /* default: grow to fit content */

  min-width: 200px;     /* never narrower than 200px */
  max-width: 600px;     /* never wider than 600px */

  min-height: 100px;    /* at least 100px tall */
  max-height: 400px;    /* at most 400px tall, then scroll */
}
```

:::tip
Prefer `max-width` over `width` for flexible layouts. `width: 600px` is a rigid box that overflows on small screens. `max-width: 600px` fills the available space up to 600px, then stops — naturally responsive without media queries.
:::

### `aspect-ratio`

Forces an element to maintain proportional dimensions:

```css
.video-wrapper {
  width: 100%;
  aspect-ratio: 16 / 9;
}

.avatar {
  width: 80px;
  aspect-ratio: 1; /* perfect square */
  border-radius: 50%;
}
```

## Fluid Sizing with `min()`, `max()`, and `clamp()`

These CSS math functions let you create responsive values without media queries.

### `min()`

Uses the **smaller** of two values:

```css
.container {
  width: min(90%, 1200px);
  /* 90% of the viewport, but never more than 1200px */
}
```

### `max()`

Uses the **larger** of two values:

```css
.sidebar {
  width: max(250px, 25%);
  /* 25% of the parent, but never less than 250px */
}
```

### `clamp()` — The Swiss Army Knife

`clamp(minimum, preferred, maximum)` — the preferred value is used unless it goes below the minimum or above the maximum:

```css
h1 {
  font-size: clamp(1.75rem, 4vw, 3rem);
  /* At least 1.75rem, ideally 4% of viewport width, at most 3rem */
}

.container {
  padding: clamp(1rem, 3vw, 3rem);
  /* Padding grows with viewport but stays within bounds */
}
```

This single line replaces multiple media queries:

```css
/* clamp() replaces all of this: */
h1 { font-size: 1.75rem; }
@media (min-width: 768px) { h1 { font-size: 2.5rem; } }
@media (min-width: 1200px) { h1 { font-size: 3rem; } }
```

:::key
`clamp()` is the most powerful tool for fluid responsive design. Use it for font sizes, spacing, and container widths. The pattern `clamp(min-rem, preferred-vw, max-rem)` handles almost every responsive sizing need.
:::

## Building a Spacing Scale

Professional design systems use a consistent spacing scale based on a base unit. The most popular approach is a **4px / 8px grid**:

```css
:root {
  --space-1:  0.25rem;  /* 4px */
  --space-2:  0.5rem;   /* 8px */
  --space-3:  0.75rem;  /* 12px */
  --space-4:  1rem;     /* 16px */
  --space-5:  1.25rem;  /* 20px */
  --space-6:  1.5rem;   /* 24px */
  --space-8:  2rem;     /* 32px */
  --space-10: 2.5rem;   /* 40px */
  --space-12: 3rem;     /* 48px */
  --space-16: 4rem;     /* 64px */
}

.card {
  padding: var(--space-6);         /* 24px */
  margin-bottom: var(--space-8);   /* 32px */
}

.card-title {
  margin-bottom: var(--space-2);   /* 8px */
}
```

:::tip
A spacing scale prevents "magic numbers" — random values like `padding: 13px` that have no relationship to anything else. With a scale, every spacing value is intentional, and the design feels cohesive.
:::

:::quiz
Q: Why is `rem` preferred over `px` for font sizes?
- `rem` renders more smoothly
- `rem` respects the user's browser font-size preference *
- `rem` is faster to compute
- `rem` works in more browsers
E: When a user changes their browser's default font size (often for accessibility), `rem` values scale proportionally because they're relative to the root font size. Pixel values ignore this preference entirely.
:::

:::quiz
Q: What does `clamp(1rem, 5vw, 3rem)` do?
- Sets the value to exactly 5vw
- Uses 1rem on mobile and 3rem on desktop
- Uses 5vw as the value, but never less than 1rem or more than 3rem *
- Sets three different values for three breakpoints
E: `clamp(min, preferred, max)` uses the preferred value (5vw) but clamps it to stay between the minimum (1rem) and maximum (3rem). It's a fluid value that smoothly scales with the viewport.
:::

## Recap

- Use `rem` for font sizes and spacing — it's predictable and respects user settings.
- Use `em` for padding/margins that should scale with an element's own font size.
- Use `px` for borders and tiny details that shouldn't scale.
- Use `dvh`/`vh` for viewport-height sections; prefer `dvh` on mobile.
- Use `ch` for setting reading-width constraints (`max-width: 65ch`).
- Prefer `max-width` over `width` for naturally responsive containers.
- `clamp()` is your best tool for fluid typography and spacing — one line replaces multiple media queries.
- Build a **spacing scale** (4px/8px grid) with custom properties to eliminate magic numbers.

**Next up:** Borders, Shadows, and Effects — adding depth and polish to your designs.

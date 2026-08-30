# Borders Shadows and Effects

CSS provides a rich set of visual tools for adding depth, polish, and personality to your elements. Borders define edges, shadows create depth, filters transform appearance, and clipping shapes break out of the rectangle. This lesson covers each one with practical patterns you'll use constantly.

## Borders

### The `border` Shorthand

```css
.card {
  border: 1px solid #e5e7eb;
}
/* Shorthand: width | style | color */
```

Common border styles: `solid`, `dashed`, `dotted`, `double`, `none`.

### Individual Sides

```css
.section {
  border-bottom: 2px solid #2563eb;  /* underline effect */
}

.table-cell {
  border-top: 1px solid #e5e7eb;
  border-left: 1px solid #e5e7eb;
}
```

### Individual Properties

```css
.box {
  border-width: 1px 2px 1px 2px;   /* top right bottom left */
  border-style: solid;
  border-color: #e5e7eb;
}
```

## `border-radius`

Rounds the corners of an element:

```css
/* All corners equal */
.card { border-radius: 8px; }

/* Each corner separately: top-left, top-right, bottom-right, bottom-left */
.tag { border-radius: 4px 4px 0 0; } /* rounded top, square bottom */

/* Perfect circle (for square elements) */
.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
}

/* Pill shape */
.badge {
  padding: 4px 12px;
  border-radius: 9999px;  /* large value → pill shape */
}
```

:::tip
For consistent design, pick 2–3 border-radius values and stick with them across your project. A common set: `4px` for small elements (inputs, tags), `8px` for cards, `50%` for avatars.
:::

## `box-shadow`

Adds shadow effects around an element's box:

```css
.card {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
/* Syntax: offset-x | offset-y | blur-radius | color */
```

### Shadow Anatomy

```css
box-shadow: 2px 4px 8px 0px rgba(0, 0, 0, 0.15);
/*          ↑    ↑    ↑   ↑    ↑
            x    y  blur spread color
*/
```

- **offset-x** — horizontal shift (positive = right).
- **offset-y** — vertical shift (positive = down).
- **blur-radius** — how soft the edges are (0 = sharp).
- **spread-radius** — grows or shrinks the shadow (optional, default 0).
- **color** — usually semi-transparent black or the brand color.

### Multiple Shadows

Layer shadows for realistic depth. Professional designs use subtle, layered shadows:

```css
.card-elevated {
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.06),
    0 4px 8px rgba(0, 0, 0, 0.08),
    0 12px 24px rgba(0, 0, 0, 0.06);
}
```

Multiple shadows are comma-separated. The first listed is on top.

### Inset Shadows

Add `inset` to draw the shadow *inside* the element:

```css
.input:focus {
  box-shadow: inset 0 0 0 2px #2563eb;
  /* Creates an inner "border" effect without changing layout */
}

.well {
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
  /* Pressed-in / recessed look */
}
```

:::key
Inset `box-shadow` with zero blur and spread is a great trick for focus rings and borders that don't affect layout — unlike real borders, shadows don't add to the element's size.
:::

:::tip
Build a shadow elevation scale with custom properties (`--shadow-sm`, `--shadow-md`, `--shadow-lg`) for consistent depth across your design system.
:::

## `text-shadow`

Adds shadow to text — same syntax as `box-shadow` but without `spread` or `inset`:

```css
.hero-title {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}
```

Use `text-shadow` sparingly — mostly for hero text over images.

## Outline and `outline-offset`

`outline` draws a line *outside* the border without affecting layout:

```css
button:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px;  /* gap between outline and border */
}
```

Unlike borders, outlines:
- Don't take up space in the box model.
- Can't be set per-side.
- Don't follow `border-radius` in older browsers (modern browsers do follow it).
- Are the standard for keyboard focus indicators.

## CSS Filters

The `filter` property applies visual effects to an entire element:

```css
.photo        { filter: grayscale(100%); }
.photo:hover  { filter: grayscale(0%); }   /* color on hover */

.blurred      { filter: blur(4px); }
.bright       { filter: brightness(1.2); }
.dark-overlay { filter: brightness(0.6); }
.contrast     { filter: contrast(1.3); }
.faded        { filter: saturate(0.5); }
.sepia        { filter: sepia(80%); }
```

### `drop-shadow()` vs `box-shadow`

`filter: drop-shadow()` follows the *actual shape* of the element, including transparency in PNGs and SVGs:

```css
/* box-shadow: rectangle around the entire element */
.png-icon { box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); }

/* drop-shadow: follows the icon's shape */
.png-icon { filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2)); }
```

Use `drop-shadow()` for images with transparency. Use `box-shadow` for everything else (it's more performant).

### Combining Filters

Chain multiple filters in one declaration:

```css
.vintage-photo {
  filter: sepia(60%) contrast(1.1) brightness(0.9);
}
```

:::warning
Filters can be expensive for performance, especially `blur()` on large elements. Avoid animating complex filters on elements larger than ~500px. If you need a blurred background, use `backdrop-filter` instead.
:::

## `backdrop-filter`

Applies effects to the area *behind* an element — the content visible through it:

```css
.glass-panel {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
}
```

This is the classic "glassmorphism" or frosted-glass effect. The element must have a semi-transparent background for `backdrop-filter` to be visible.

```css
.sticky-nav {
  position: sticky;
  top: 0;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(8px);
}
```

## `mix-blend-mode`

Controls how an element's colors blend with the content behind it:

```css
.overlay-text {
  mix-blend-mode: multiply;     /* darkens */
}

.light-blend {
  mix-blend-mode: screen;       /* lightens */
}

.dramatic {
  mix-blend-mode: difference;   /* inverts overlapping colors */
}
```

Common modes: `multiply` (darkens), `screen` (lightens), `overlay` (contrast), `difference` (invert).

## `clip-path` Basics

Clips an element to a shape, hiding everything outside:

```css
.avatar { clip-path: circle(50%); }
.arrow  { clip-path: polygon(50% 0%, 0% 100%, 100% 100%); }
.hero   { clip-path: polygon(0 0, 100% 0, 100% 85%, 0 100%); } /* angled divider */
```

`clip-path` is powerful for non-rectangular layouts, diagonal dividers, and creative shapes — all without images.

:::quiz
Q: What's the difference between `box-shadow` and `filter: drop-shadow()`?
- They are identical
- `box-shadow` follows the element's shape including transparency, `drop-shadow` doesn't
- `drop-shadow` follows the element's actual shape including transparency, `box-shadow` is always rectangular *
- `box-shadow` works on text, `drop-shadow` doesn't
E: `box-shadow` creates a rectangular shadow based on the element's box. `filter: drop-shadow()` traces the actual visual shape of the element, including transparent areas in PNGs and SVGs, creating a more natural shadow.
:::

:::quiz
Q: How do you create a frosted-glass effect in CSS?
- Use `filter: blur()` on the element
- Use a semi-transparent background with `backdrop-filter: blur()` *
- Set `opacity: 0.5` and `filter: brightness(1.2)`
- Use `mix-blend-mode: screen`
E: The frosted-glass effect requires `backdrop-filter: blur()` combined with a semi-transparent background. The blur applies to the content *behind* the element, while the semi-transparent background lets that blurred content show through.
:::

## Recap

- `border-radius` rounds corners: `8px` for cards, `50%` for circles, `9999px` for pills.
- Layer multiple `box-shadow`s for realistic depth. Use inset shadows for pressed effects and focus rings.
- Use `text-shadow` sparingly — mostly for hero text over images.
- `filter` applies visual effects: `blur()`, `grayscale()`, `brightness()`, `drop-shadow()`.
- `drop-shadow()` follows transparency; `box-shadow` is always rectangular.
- `backdrop-filter: blur()` creates frosted-glass effects on semi-transparent elements.
- `clip-path` breaks out of rectangles — circles, polygons, and diagonal edges.
- Build a **shadow scale** with custom properties for consistent elevation across your design.

**Next up:** Pseudo-classes and Pseudo-elements — styling element states and creating decorative content without extra HTML.

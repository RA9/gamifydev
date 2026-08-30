# Modern CSS Features

CSS is evolving fast. Features that used to require JavaScript or complex workarounds are now built into the language. This lesson covers the most impactful modern CSS additions — the ones already shipping in browsers and changing how professionals write stylesheets.

## Container queries (@container)

Media queries ask "how wide is the *viewport*?" Container queries ask "how wide is this element's *parent*?" That's a game-changer for reusable components.

### Setup

First, declare a containment context on the parent:

```css
.card-grid {
  container-type: inline-size;  /* track this element's width */
  container-name: cards;        /* optional name */
}
```

Then write queries against it:

```css
@container cards (min-width: 500px) {
  .card {
    display: grid;
    grid-template-columns: 200px 1fr;
    gap: 16px;
  }
}

@container cards (max-width: 499px) {
  .card {
    display: flex;
    flex-direction: column;
  }
}
```

:::key
Container queries let components adapt to their *context*, not the viewport. A card in a narrow sidebar shows a vertical layout; the same card in a wide main area shows a horizontal layout — without the parent knowing anything about the card's internals. This is true component-level responsive design.
:::

## The :has() selector

Called the "parent selector" — something CSS developers wanted for over a decade. `:has()` selects an element based on what it *contains*:

```css
/* Style a card differently when it contains an image */
.card:has(img) {
  padding: 0;
}

/* Style a form group when its input is invalid */
.form-group:has(:invalid) {
  border-left: 3px solid #ef4444;
}

/* Style a label when its sibling checkbox is checked */
label:has(+ input:checked) {
  font-weight: 700;
  color: #3b82f6;
}
```

### Practical: required field indicator

```css
label:has(+ input:required)::after {
  content: ' *';
  color: #ef4444;
}
```

No JavaScript needed. The asterisk appears automatically when the adjacent input has the `required` attribute.

:::tip
`:has()` is incredibly versatile — it can replace a lot of JavaScript-driven class toggling. Think of it as "if this element contains X, style it differently." Check browser support (now in all major browsers as of late 2023) and use it confidently.
:::

## CSS nesting

Write child selectors *inside* their parent, just like Sass — but natively:

```css
.nav {
  display: flex;
  gap: 16px;

  a {
    color: #64748b;
    text-decoration: none;
    transition: color 200ms;

    &:hover {
      color: #0f172a;
    }

    &.active {
      color: #3b82f6;
      font-weight: 600;
    }
  }
}
```

The `&` refers to the parent selector, just like in Sass. For simple descendant selectors (like `a` above), the `&` is implied.

:::warning
Don't nest too deeply — three levels is a practical maximum. Deeply nested CSS creates overly specific selectors that are hard to override and maintain. The rule from Sass still applies: if you need more than 3 levels, rethink your structure.
:::

## Subgrid

When a child element uses `display: grid`, its tracks are independent of the parent grid. Subgrid lets a child *inherit* the parent's tracks, so everything lines up:

```css
.card-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.card {
  display: grid;
  grid-template-rows: subgrid;    /* inherit parent's row tracks */
  grid-row: span 3;               /* occupy 3 of the parent's rows */
}
```

This is invaluable for card layouts where titles, descriptions, and buttons need to align across cards — even when the content lengths vary.

## color-mix() and oklch colors

### color-mix()

Blend two colors directly in CSS:

```css
.btn-hover {
  /* Mix the primary color with black at 20% */
  background: color-mix(in srgb, var(--color-primary) 80%, black);
}

.overlay {
  background: color-mix(in srgb, var(--color-primary) 30%, transparent);
}
```

This replaces the need for manually calculating darker or lighter variants of a color.

### oklch — perceptually uniform colors

`oklch` defines colors by Lightness, Chroma (saturation), and Hue — in a perceptually uniform space. Two colors with the same L value actually *look* equally bright, unlike hex or HSL:

```css
:root {
  --blue: oklch(0.6 0.2 250);
  --green: oklch(0.6 0.2 145);
  --red: oklch(0.6 0.2 25);
  /* Same lightness and chroma — they genuinely look equally vivid */
}
```

:::tip
`oklch` is the best color space for design systems. You can generate consistent palettes by varying only the hue while keeping lightness and chroma fixed. Every color in the palette will look equally bright and saturated.
:::

## Scroll snap

Create carousel-like scroll behavior without JavaScript:

```css
.carousel {
  display: flex;
  overflow-x: auto;
  scroll-snap-type: x mandatory;    /* snap horizontally, always */
  gap: 16px;
}

.carousel-item {
  scroll-snap-align: start;        /* snap to the start edge */
  flex: 0 0 300px;
}
```

Values for `scroll-snap-type`:

- `x mandatory` — always snaps to the nearest item on the x-axis
- `y proximity` — snaps vertically if close enough
- `both mandatory` — snaps in both directions

## Scroll-driven animations

Link animation progress to scroll position — no JavaScript needed:

```css
@keyframes reveal {
  from { opacity: 0; transform: translateY(20px); }
  to   { opacity: 1; transform: translateY(0); }
}

.section {
  animation: reveal linear both;
  animation-timeline: view();          /* tied to scroll visibility */
  animation-range: entry 0% entry 100%;  /* animate during entry */
}
```

The animation plays as the element scrolls into view. `view()` creates a timeline based on the element's visibility in the viewport.

:::warning
Scroll-driven animations are newer than other features on this page. Check browser support before relying on them in production. As of 2024, Chromium browsers have full support; Firefox and Safari are catching up. Always provide a fallback (the element should be visible without animation).
:::

## Logical properties

Logical properties replace physical directions (`left`, `right`, `top`, `bottom`) with flow-relative ones that work correctly in any writing direction:

| Physical          | Logical                      |
|-------------------|------------------------------|
| `margin-left`     | `margin-inline-start`        |
| `margin-right`    | `margin-inline-end`          |
| `margin-top`      | `margin-block-start`         |
| `margin-bottom`   | `margin-block-end`           |
| `padding-left/right` | `padding-inline`          |
| `padding-top/bottom` | `padding-block`            |
| `width`           | `inline-size`                |
| `height`          | `block-size`                 |
| `text-align: left`| `text-align: start`          |

```css
.card {
  margin-block: 16px;      /* top and bottom */
  padding-inline: 24px;    /* left and right (or start/end in RTL) */
  border-inline-start: 3px solid var(--color-primary);
}
```

:::key
If your site might ever support right-to-left languages (Arabic, Hebrew), logical properties handle it automatically. Even if not, they're clearer about intent: `margin-inline: auto` obviously means "center horizontally" more clearly than `margin-left: auto; margin-right: auto`.
:::

## Checking browser support

### @supports in CSS

Test for feature support directly:

```css
/* Fallback layout */
.grid { display: flex; flex-wrap: wrap; }

/* Enhanced layout if grid is supported */
@supports (display: grid) {
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  }
}

@supports (container-type: inline-size) {
  .card-wrapper { container-type: inline-size; }
}
```

### caniuse.com

Before using any feature in production, check **caniuse.com** for current browser support percentages. Look at the global usage stats — if support is above 95%, you're generally safe. Below that, provide a fallback.

:::quiz
Q: What problem do container queries solve that media queries cannot?
- Container queries are faster than media queries
- Container queries let components respond to their parent's size instead of the viewport size *
- Container queries work without CSS
E: Media queries respond to the viewport. Container queries respond to the size of the element's container, making components truly reusable across different layout contexts.
:::

:::quiz
Q: What does the `:has()` selector allow you to do for the first time in CSS?
- Select the next sibling
- Select a parent element based on its children or contents *
- Select elements by their z-index
E: `:has()` is the long-awaited "parent selector." It lets you style an element based on what it contains — for example, styling a `.card` differently when it contains an `img` element.
:::

## Recap

- **Container queries** (`@container`) make components respond to their parent's size, not the viewport.
- **`:has()`** selects elements based on their contents — the CSS "parent selector."
- **CSS nesting** lets you write child selectors inside parents, just like Sass but native.
- **Subgrid** lets child grids inherit parent track sizes for perfect alignment.
- **`color-mix()`** blends colors in CSS; **`oklch`** provides perceptually uniform color definitions.
- **Scroll snap** creates carousel behavior; **scroll-driven animations** tie animation progress to scroll position.
- **Logical properties** (`inline`/`block`) replace physical directions for internationalization-ready CSS.
- Use **`@supports`** for feature detection and **caniuse.com** for browser support data.

**Next up:** CSS Architecture and Organization — structuring your stylesheets for projects that scale.

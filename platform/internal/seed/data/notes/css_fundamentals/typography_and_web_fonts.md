# Typography and Web Fonts

Typography is arguably the most impactful CSS skill you can develop. Text makes up the majority of most web pages, so getting it right — readable, beautiful, consistent — transforms the entire feel of a site. This lesson covers every typographic property and how to load custom fonts.

## `font-family` and Font Stacks

`font-family` specifies which font to use. You always provide a **stack** — a list of fonts in order of preference. The browser uses the first one it finds installed:

```css
body {
  font-family: "Inter", "Helvetica Neue", Arial, sans-serif;
}
```

- `"Inter"` — the preferred font (loaded via @font-face or Google Fonts).
- `"Helvetica Neue"` — fallback on macOS.
- `Arial` — fallback on Windows.
- `sans-serif` — generic family, the ultimate fallback.

:::warning
Always end your font stack with a **generic family** (`serif`, `sans-serif`, `monospace`, `cursive`, or `system-ui`). Without it, the browser picks whatever default it wants if no named font is available.
:::

### The System Font Stack

Use the OS's native UI font — fast (no download) and familiar to users:

```css
body {
  font-family: system-ui, -apple-system, BlinkMacSystemFont,
               "Segoe UI", Roboto, "Helvetica Neue", Arial,
               "Noto Sans", sans-serif;
}
```

This gives you San Francisco on Apple, Segoe UI on Windows, and Roboto on Android — all without loading a single font file. Many professional sites use this approach.

:::tip
If you don't have a brand font requirement, the system font stack is the best default. Zero download cost, great readability, and native feel on every platform.
:::

## Loading Web Fonts

### Google Fonts with `<link>` (Recommended)

Add a `<link>` in your HTML `<head>`:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap"
      rel="stylesheet">
```

The `preconnect` hints tell the browser to open connections early, reducing latency.

### Google Fonts with `@import` (Simpler, Slower)

Add this at the top of your CSS:

```css
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700&display=swap');
```

This is simpler but **blocks rendering** until the import is fully loaded. The `<link>` approach is preferred for performance.

### Self-Hosting with `@font-face`

For maximum control and privacy, host fonts yourself:

```css
@font-face {
  font-family: "Inter";
  src: url("/fonts/inter-regular.woff2") format("woff2");
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: "Inter";
  src: url("/fonts/inter-bold.woff2") format("woff2");
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
```

Use `.woff2` format — it has the best compression and broadest support.

### `font-display: swap`

This is critical for performance. It controls what happens while the custom font is loading:

- `swap` — Show text immediately in the fallback font, swap to the custom font once loaded. Text is always visible.
- `block` — Hide text for up to 3 seconds waiting for the font. Bad for UX.
- `fallback` — Short block period (100ms), then fallback. Doesn't swap if font loads too late.
- `optional` — Very short block, then fallback permanently if font isn't cached. Best for non-essential fonts.

:::key
Always use `font-display: swap` (or `optional`) to prevent invisible text while fonts load. Users should never see a blank page waiting for a font file.
:::

## Font Properties

### `font-weight`

Controls thickness — from 100 (thin) to 900 (black):

```css
.light    { font-weight: 300; }
.regular  { font-weight: 400; } /* normal */
.medium   { font-weight: 500; }
.semibold { font-weight: 600; }
.bold     { font-weight: 700; } /* bold */
```

Only use weights your font actually provides. If you set `font-weight: 600` but only loaded 400 and 700, the browser will *fake* it with synthetic bolding, which looks bad.

### `font-style`

```css
.italic { font-style: italic; }
.normal { font-style: normal; }
```

### `font-size`

Sets the size of the text. Use `rem` for consistency:

```css
body { font-size: 16px; }        /* base size */
h1   { font-size: 2.5rem; }      /* 40px relative to root */
h2   { font-size: 2rem; }        /* 32px */
.small { font-size: 0.875rem; }  /* 14px */
```

We'll cover units in depth in the Units lesson — for now, know that `rem` is the standard for font sizing.

## Line Height

`line-height` controls the space between lines of text. It's the single biggest factor in readability.

```css
body {
  line-height: 1.6; /* unitless — recommended */
}

h1 {
  line-height: 1.2; /* tighter for headings */
}
```

:::key
Always use a **unitless** `line-height` (like `1.6`, not `1.6em` or `26px`). A unitless value is a *multiplier* of the current font size, so it scales correctly with every element that inherits it. A value with units creates a fixed height that causes overlapping text at larger font sizes.
:::

```css
/* Bad: fixed line-height causes problems */
body { font-size: 16px; line-height: 24px; }
h1   { font-size: 48px; } /* Line-height is still 24px → lines overlap! */

/* Good: multiplier scales with font size */
body { font-size: 16px; line-height: 1.5; }
h1   { font-size: 48px; } /* Line-height = 48 × 1.5 = 72px ✓ */
```

## Spacing Properties

### `letter-spacing`

Adjusts space between characters. Useful for uppercase text and headings:

```css
.uppercase-label {
  text-transform: uppercase;
  letter-spacing: 0.05em; /* slight spread for readability */
  font-size: 0.75rem;
}

h1 {
  letter-spacing: -0.02em; /* slightly tighter for large headings */
}
```

### `word-spacing`

Adjusts space between words — rarely needed but occasionally useful:

```css
.spread-text { word-spacing: 0.1em; }
```

## Text Transformation and Decoration

### `text-transform`

Changes the case of text *visually* without changing the HTML:

```css
.uppercase { text-transform: uppercase; }
.capitalize { text-transform: capitalize; } /* First Letter Of Each Word */
.lowercase { text-transform: lowercase; }
```

### `text-decoration`

Controls underlines, overlines, and line-throughs:

```css
a { text-decoration: none; }                         /* remove link underline */
a:hover { text-decoration: underline; }

.strikethrough { text-decoration: line-through; }

/* Modern: control style, color, and thickness */
a {
  text-decoration: underline;
  text-decoration-color: #93c5fd;
  text-decoration-thickness: 2px;
  text-underline-offset: 3px;       /* space between text and underline */
}
```

### `text-align`

Horizontal alignment of text within its container:

```css
.center { text-align: center; }
.right  { text-align: right; }
.left   { text-align: left; }    /* default for LTR languages */
```

## Building a Type Scale

Professional typography uses a consistent set of sizes — a **type scale**. Instead of picking arbitrary sizes, choose a ratio and generate sizes mathematically:

```css
:root {
  --text-xs:   0.75rem;   /* 12px */
  --text-sm:   0.875rem;  /* 14px */
  --text-base: 1rem;      /* 16px */
  --text-lg:   1.125rem;  /* 18px */
  --text-xl:   1.25rem;   /* 20px */
  --text-2xl:  1.5rem;    /* 24px */
  --text-3xl:  1.875rem;  /* 30px */
  --text-4xl:  2.25rem;   /* 36px */
  --text-5xl:  3rem;      /* 48px */
}

h1 { font-size: var(--text-5xl); }
h2 { font-size: var(--text-4xl); }
h3 { font-size: var(--text-3xl); }
body { font-size: var(--text-base); }
.caption { font-size: var(--text-sm); }
```

:::tip
Using custom properties for your type scale means you can adjust every heading on the site by changing one variable. This is how design systems like Tailwind CSS and Material Design work under the hood.
:::

:::quiz
Q: Why should you use a unitless value for `line-height`?
- It renders faster in the browser
- It's the only valid syntax
- It scales proportionally with font size, preventing overlapping text *
- It makes text look bolder
E: A unitless `line-height` like `1.5` is a multiplier of the element's font size. If the font size is 48px, line-height becomes 72px. A fixed value like `24px` would cause 48px text to overlap.
:::

:::quiz
Q: What does `font-display: swap` do?
- It swaps between two different fonts randomly
- It shows text immediately in a fallback font, then swaps to the custom font when loaded *
- It makes the font load faster
- It hides text until the custom font loads
E: `font-display: swap` ensures text is always visible. It renders immediately using the fallback font, then swaps to the custom font once it finishes downloading. This prevents the "flash of invisible text" problem.
:::

## Recap

- Use **font stacks** with a generic family at the end. The **system font stack** is a great zero-cost default.
- Load web fonts with `<link>` (faster) or `@import` (simpler). Self-host with `@font-face` for full control.
- Always set `font-display: swap` to keep text visible while fonts load.
- Use **unitless `line-height`** (e.g. `1.6`) — it scales correctly with font size.
- Add `letter-spacing` to uppercase text and tighten it on large headings.
- Build a **type scale** with CSS custom properties for consistency.
- Set `max-width: 65ch` on body text for optimal reading width.

**Next up:** Units, Sizing, and Spacing — choosing the right measurement for every situation.

# Colors Backgrounds and Gradients

Color is one of the most powerful tools in visual design — it directs attention, communicates meaning, and sets mood. CSS gives you multiple color systems, sophisticated background controls, and gradients that can replace many images entirely. This lesson covers every color tool you'll need.

## Color Formats

### Named Colors

CSS defines 148 named colors like `red`, `tomato`, `cornflowerblue`, and `rebeccapurple`:

```css
.error { color: red; }
.info  { color: dodgerblue; }
```

Named colors are great for quick prototyping but limited for production — you rarely find the exact shade you want.

### Hexadecimal

The most common format in professional CSS. A `#` followed by 6 hex digits (2 each for red, green, blue):

```css
.brand { color: #2563eb; }       /* full 6-digit hex */
.muted { color: #6b7280; }
.dark  { background: #111; }     /* 3-digit shorthand: #111 = #111111 */
```

You can add a 4th pair for alpha (opacity): `#2563eb80` = 50% transparent blue.

### `rgb()` and `rgba()`

Specify red, green, blue as numbers (0–255) or percentages:

```css
.box { color: rgb(37, 99, 235); }
.overlay { background: rgba(0, 0, 0, 0.5); } /* black at 50% opacity */

/* Modern syntax: no commas needed */
.box { color: rgb(37 99 235); }
.overlay { background: rgb(0 0 0 / 0.5); }
```

The modern syntax uses spaces and a `/` before the alpha value. `rgba()` is now just an alias for `rgb()` with alpha.

### `hsl()` and `hsla()`

**Hue, Saturation, Lightness** — the most intuitive color model for humans:

```css
.primary { color: hsl(220, 90%, 54%); }
/* Modern syntax */
.primary { color: hsl(220 90% 54%); }
.primary-light { color: hsl(220 90% 70%); } /* same hue, lighter */
.primary-dark  { color: hsl(220 90% 35%); } /* same hue, darker */
```

- **Hue** — a degree on the color wheel (0–360). 0=red, 120=green, 240=blue.
- **Saturation** — how vivid (0%=gray, 100%=full color).
- **Lightness** — how bright (0%=black, 50%=pure color, 100%=white).

:::tip
HSL is the best format for building color palettes. To get lighter or darker shades of the same color, just change the **lightness** value. To make it more muted, reduce **saturation**. This is why design-focused codebases prefer HSL.
:::

### `currentColor`

A special keyword that references the element's computed `color` value:

```css
.btn {
  color: #2563eb;
  border: 2px solid currentColor;  /* border matches text color */
}
.btn:hover {
  color: #1d4ed8; /* border changes too, automatically */
}
```

This is extremely useful — change `color` once and every `currentColor` reference updates. It works in borders, backgrounds, shadows, SVG fills, and more.

## Opacity

The `opacity` property makes an entire element (and its children) transparent:

```css
.faded { opacity: 0.5; }  /* 50% transparent */
```

Unlike `rgba()`, which only affects the *color* of one property, `opacity` affects *everything* — text, background, borders, children. Use `rgba()`/`hsla()` when you want just the background to be transparent.

```css
/* Only the background is transparent; text stays fully opaque */
.overlay {
  background: rgb(0 0 0 / 0.6);
  color: white;
}
```

## Backgrounds

### `background-color`

Fills the content + padding area with a solid color:

```css
.card { background-color: #f9fafb; }
```

### `background-image`

Sets an image (or gradient) behind the content:

```css
.hero {
  background-image: url('hero.jpg');
}
```

### `background-size`

Controls how the image scales:

```css
.hero {
  background-size: cover;    /* fill the element, cropping if needed */
}
.logo-area {
  background-size: contain;  /* fit the whole image, may leave gaps */
}
.pattern {
  background-size: 60px 60px; /* exact dimensions for tiling patterns */
}
```

- **`cover`** — the image covers the entire element; some parts may be clipped. Best for hero images.
- **`contain`** — the whole image fits inside; there may be empty space. Best for logos.

### `background-position`

Where to anchor the image:

```css
.hero {
  background-position: center;          /* centered both ways */
  background-position: top right;       /* pinned to top-right corner */
  background-position: 50% 30%;         /* 50% from left, 30% from top */
}
```

### `background-repeat`

Whether the image tiles:

```css
.hero    { background-repeat: no-repeat; }   /* single image */
.pattern { background-repeat: repeat; }       /* tile both directions (default) */
.stripe  { background-repeat: repeat-x; }     /* tile horizontally only */
```

### The `background` Shorthand

Combine everything in one declaration:

```css
.hero {
  background: url('hero.jpg') center / cover no-repeat #1a1a2e;
}
/* Breakdown: image position / size repeat fallback-color */
```

:::warning
The shorthand resets any background property you don't include. If you already set `background-color` and then use `background: url(...)`, the color gets reset. When in doubt, use the individual properties.
:::

## Gradients

Gradients are generated images — they're values for `background-image`, not `background-color`.

### `linear-gradient()`

Creates a gradient along a line:

```css
/* Top to bottom (default) */
.box { background: linear-gradient(#667eea, #764ba2); }

/* At an angle */
.box { background: linear-gradient(135deg, #667eea, #764ba2); }

/* To a direction */
.box { background: linear-gradient(to right, #ff6b6b, #feca57); }

/* Multiple stops with positions */
.box {
  background: linear-gradient(
    to right,
    #ff6b6b 0%,
    #feca57 50%,
    #48dbfb 100%
  );
}
```

You can create sharp transitions by placing two stops at the same position (e.g. `#002395 33%, #ffffff 33%`) — this produces stripes with no blend.

### `radial-gradient()`

Creates a gradient radiating from a center point:

```css
.spotlight {
  background: radial-gradient(circle, #fbbf24, #b45309);
}

.ellipse {
  background: radial-gradient(ellipse at top left, #dbeafe, transparent);
}
```

### Multiple Backgrounds

CSS supports layering multiple backgrounds — the first one listed is on top:

```css
.hero {
  background:
    linear-gradient(to bottom, transparent 60%, rgba(0, 0, 0, 0.8)),
    url('hero.jpg') center / cover no-repeat;
}
```

This classic pattern overlays a gradient on a photo — the image shows through the transparent top, and text at the bottom sits on the dark overlay.

:::tip
Overlaying a semi-transparent gradient on a background image is the go-to technique for making text readable over photos. You'll use this pattern constantly.
:::

## `color-mix()`

A modern function that mixes two colors together:

```css
.muted-primary {
  /* Mix 70% primary with 30% white */
  color: color-mix(in srgb, #2563eb 70%, white);
}

.hover-darken {
  /* Mix 80% of the color with 20% black */
  background: color-mix(in srgb, #2563eb 80%, black);
}
```

The `in srgb` part specifies the color space. You can also use `in hsl`, `in oklch`, etc. for different blending behavior.

`color-mix()` is powerful combined with custom properties:

```css
:root {
  --brand: #2563eb;
}
.btn {
  background: var(--brand);
}
.btn:hover {
  background: color-mix(in srgb, var(--brand) 85%, black);
}
```

Now changing `--brand` automatically generates a consistent hover shade.

:::quiz
Q: What's the main advantage of HSL over hex for building color palettes?
- HSL loads faster in the browser
- You can change lightness/saturation intuitively without affecting the base hue *
- HSL supports more colors than hex
- HSL is the only format that supports transparency
E: In HSL, you can create lighter, darker, or more muted variants of a color by adjusting just the lightness or saturation value while keeping the hue constant. With hex, creating related shades requires calculating all three channels.
:::

:::quiz
Q: What does `background-size: cover` do?
- Stretches the image to exactly fill the element, possibly distorting it
- Scales the image to fill the element completely, cropping if needed to maintain aspect ratio *
- Scales the image to fit entirely inside the element, leaving gaps if needed
- Tiles the image to cover the element
E: `cover` scales the image proportionally until it completely fills the element. If the aspect ratios don't match, parts of the image are cropped. `contain` is the alternative that shows the entire image, potentially leaving gaps.
:::

## Recap

- **Hex** is standard, **HSL** is best for palettes, `rgb()` is great with alpha.
- `currentColor` references the element's `color` — use it for borders and icons.
- Use `rgba()`/`hsla()` for transparent backgrounds; `opacity` affects *everything*.
- `background-size: cover` for hero images, `contain` for logos.
- **Gradients** are `background-image` values — use hard stops for stripes, smooth stops for blends.
- **Multiple backgrounds** layer top-to-bottom — gradient-over-photo is a core pattern.
- `color-mix()` generates tints and shades from a single color dynamically.

**Next up:** Typography and Web Fonts — making your text beautiful and readable.

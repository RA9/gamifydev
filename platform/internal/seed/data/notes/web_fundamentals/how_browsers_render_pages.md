# How Browsers Render Pages

You write HTML, CSS, and JavaScript. But how does the browser turn those text files into the visual, interactive page you see on screen? Understanding the rendering pipeline demystifies performance issues, explains why script placement matters, and gives you the knowledge to build faster websites.

## The rendering pipeline overview

When a browser receives an HTML file, it goes through a series of steps to produce pixels on screen:

```
HTML  →  DOM Tree
                  ↘
                    Render Tree  →  Layout  →  Paint  →  Compositing  →  Pixels
                  ↗
CSS   →  CSSOM
```

Each step feeds into the next. Let's walk through them.

## Step 1: Parsing HTML into the DOM

The browser reads the HTML file top to bottom and constructs the **DOM** (Document Object Model) — a tree structure where every HTML element becomes a node:

```html
<html>
  <head>
    <title>My Page</title>
  </head>
  <body>
    <h1>Hello</h1>
    <p>Welcome to my site</p>
  </body>
</html>
```

Becomes this tree:

```
Document
└── html
    ├── head
    │   └── title
    │       └── "My Page"
    └── body
        ├── h1
        │   └── "Hello"
        └── p
            └── "Welcome to my site"
```

The DOM is the browser's internal representation of your page. JavaScript interacts with this tree — when you call `document.querySelector('h1')`, you're reaching into the DOM.

:::key
The DOM is not your HTML file — it's a live, in-memory tree the browser builds *from* your HTML. JavaScript can modify the DOM (add, remove, change elements), and those changes are reflected on screen. The DOM is the bridge between your code and what users see.
:::

## Step 2: Parsing CSS into the CSSOM

While building the DOM, the browser also parses all CSS (from `<link>` tags, `<style>` blocks, and inline styles) into the **CSSOM** (CSS Object Model) — another tree that maps every element to its computed styles:

```css
body { font-family: sans-serif; color: #333; }
h1   { font-size: 2rem; color: #000; }
p    { line-height: 1.6; }
```

The CSSOM resolves inheritance, specificity, and the cascade to determine the *final* computed style for every element. For example, the `<p>` inherits `font-family: sans-serif` and `color: #333` from `body`, even though those aren't declared on `p` directly.

:::tip
Open DevTools → select any element → look at the **Computed** tab. That's the CSSOM's output — every single CSS property resolved to its final value, including inherited and default values. It's the definitive answer to "what style does this element actually have?"
:::

## Step 3: The Render Tree

The browser combines the DOM and CSSOM into the **Render Tree** — a tree of only the elements that will actually be *visible* on screen:

- Elements with `display: none` are **excluded** (they produce no visual output).
- `visibility: hidden` elements are **included** (they're invisible but still take space).
- The `<head>` element and its children are excluded (they're not visual).
- Pseudo-elements (`::before`, `::after`) are added to the render tree even though they're not in the DOM.

```
Render Tree (simplified):
body (font-family: sans-serif)
├── h1 "Hello" (font-size: 2rem, color: #000)
└── p "Welcome" (line-height: 1.6, color: #333)
```

## Step 4: Layout (Reflow)

The browser calculates the **exact position and size** of every element in the render tree. This step answers:

- How wide is this element?
- How tall is this element?
- Where does it sit relative to its parent and siblings?

Layout works from the outside in: the viewport determines the body's width, which determines its children's available space, and so on recursively.

```
Viewport: 1200px wide
└── body: 1200px wide, starts at (0, 0)
    ├── h1: 1200px wide, 48px tall, at (0, 0)
    └── p: 1200px wide, 28px tall, at (0, 48)
```

:::warning
Layout is **expensive**. Changing an element's width, height, margin, padding, or position forces the browser to recalculate the position of that element *and* everything affected by it. This is called a **reflow**. Avoid triggering layout in animations — use `transform` instead of `top`/`left`/`width`.
:::

## Step 5: Paint

Once layout is complete, the browser **fills in the pixels**: text, colors, backgrounds, borders, shadows, images. Each element's visual properties are rasterized into bitmap layers.

Painting is done in a specific order:

1. Background color
2. Background image
3. Border
4. Children (text, nested elements)
5. Outline

Changes that affect appearance but not geometry (like `color`, `background-color`, `visibility`) trigger a **repaint** without a full layout recalculation — cheaper than reflow, but still not free.

## Step 6: Compositing

The browser often splits the page into multiple **layers** (think of them as transparent sheets stacked on top of each other). Compositing is the step where these layers are assembled into the final image you see.

Elements that get their own layer:

- Elements with `transform` or `opacity` animations
- Elements with `will-change: transform`
- Fixed or sticky positioned elements
- Elements with `overflow: scroll`

:::key
Compositing is the *cheapest* rendering step. When you animate `transform` or `opacity`, the browser only needs to move or fade an already-painted layer — no layout, no repaint. That's why these two properties are the key to smooth 60fps animations.
:::

## Why script placement matters

When the HTML parser encounters a `<script>` tag, it **stops parsing** to download and execute the script. This blocks DOM construction:

```html
<head>
  <!-- This blocks everything until the script downloads and runs -->
  <script src="heavy-app.js"></script>
</head>
<body>
  <!-- Nothing below is parsed until the script is done -->
  <h1>Hello</h1>
</body>
```

### Solution 1: Place scripts at the end of `<body>`

```html
<body>
  <h1>Hello</h1>
  <p>Content renders first</p>

  <!-- Scripts load after all HTML is parsed -->
  <script src="app.js"></script>
</body>
```

### Solution 2: Use `defer` or `async`

```html
<head>
  <!-- defer: downloads in parallel, executes after HTML parsing is done -->
  <script src="app.js" defer></script>

  <!-- async: downloads in parallel, executes as soon as it's ready (may interrupt parsing) -->
  <script src="analytics.js" async></script>
</head>
```

| Attribute | Download        | Execute                     | Use when                   |
|-----------|-----------------|-----------------------------|----------------------------|
| (none)    | Blocks parsing  | Immediately, blocks parsing | Almost never               |
| `defer`   | Parallel        | After HTML parsing          | Your main app scripts      |
| `async`   | Parallel        | As soon as ready            | Independent scripts (analytics) |

:::tip
Use `defer` for most scripts — it gives the best of both worlds. The script downloads while the HTML parses, and runs after the DOM is ready. Use `async` only for scripts that don't depend on the DOM or other scripts.
:::

## Render-blocking resources

CSS is also render-blocking: the browser won't paint anything until it has built the CSSOM, because painting without styles would cause a flash of unstyled content (FOUC).

```html
<head>
  <!-- This blocks rendering until the CSS is downloaded and parsed -->
  <link rel="stylesheet" href="styles.css" />
</head>
```

To minimize the impact:

- **Keep CSS files small** — split critical styles from non-critical ones.
- **Inline critical CSS** for above-the-fold content:

```html
<head>
  <style>
    /* Only styles needed for the initial viewport */
    body { font-family: sans-serif; margin: 0; }
    .hero { height: 100vh; display: grid; place-items: center; }
  </style>
  <!-- Full stylesheet loads without blocking the initial paint -->
  <link rel="stylesheet" href="full-styles.css" media="print" onload="this.media='all'" />
</head>
```

## The Critical Rendering Path

The **Critical Rendering Path** (CRP) is the sequence of steps the browser must complete before it can paint the first pixel. Optimizing it means reducing the time to **first paint**:

```
1. Download HTML
2. Parse HTML → build DOM
3. Discover CSS → download → parse → build CSSOM
4. Discover JS → download → execute (if blocking)
5. Combine DOM + CSSOM → Render Tree
6. Layout → Paint → Compositing
7. FIRST PAINT ← this is what users see
```

Every render-blocking resource (CSS, synchronous JS) adds to this path. The goal is to minimize what's needed before that first paint.

### CRP optimization checklist

```
✓ Minimize the number of render-blocking resources
✓ Use defer/async on scripts
✓ Inline critical CSS, defer the rest
✓ Compress and minify CSS and JS files
✓ Use modern image formats (WebP, AVIF)
✓ Set appropriate cache headers
```

:::quiz
Q: What is the DOM?
- The CSS rules applied to a page
- A tree structure the browser builds from HTML, representing every element as a node *
- A JavaScript library for manipulating pages
E: The DOM (Document Object Model) is the browser's internal tree representation of the HTML document. JavaScript interacts with the DOM to read and modify page content.
:::

:::quiz
Q: Why is animating `transform` and `opacity` faster than animating `width` or `top`?
- `transform` and `opacity` use smaller file sizes
- `transform` and `opacity` only require compositing, skipping the expensive layout and paint steps *
- `transform` and `opacity` are newer CSS properties
E: Animating `transform` or `opacity` only requires the compositing step — the browser moves or fades an already-painted layer. Animating `width` or `top` triggers layout (recalculating positions) and paint (redrawing pixels), which is far more expensive.
:::

:::quiz
Q: What does the `defer` attribute on a `<script>` tag do?
- Prevents the script from running entirely
- Downloads the script in parallel with HTML parsing and executes it after parsing is complete *
- Makes the script load synchronously
E: `defer` allows the browser to download the script while continuing to parse HTML. The script executes only after the full HTML document has been parsed, preserving execution order and avoiding parsing delays.
:::

## Recap

- The browser rendering pipeline: **HTML → DOM**, **CSS → CSSOM**, **Render Tree → Layout → Paint → Compositing → Pixels**.
- The **DOM** is a live tree built from HTML; the **CSSOM** resolves all CSS to computed styles.
- The **Render Tree** excludes invisible elements (`display: none`, `<head>`).
- **Layout** calculates positions and sizes; **paint** fills in pixels; **compositing** assembles layers.
- Scripts block HTML parsing unless you use `defer` (recommended) or `async`.
- CSS is render-blocking — minimize it and inline critical styles for faster first paint.
- The **Critical Rendering Path** is the minimum work before the first pixel. Optimize it for fast-loading pages.

**Next up:** Setting Up Your Dev Environment — the tools every web developer needs to start building.

# Mobile-First Workflow

Most web traffic is mobile. If your site doesn't work beautifully on a phone, half your audience is gone before they even see your content. Mobile-first design flips the traditional approach: you start with the smallest screen and *add* complexity as space increases.

## Why mobile-first beats desktop-first

In a desktop-first workflow you build a rich, wide layout and then wrestle it into a phone-sized box — overriding widths, hiding sidebars, and undoing multi-column grids. You end up writing your layout *twice*.

Mobile-first is the opposite: your base CSS is for the narrowest screen, and you progressively *add* layout features as the viewport grows.

```css
/* Desktop-first: undo everything on small screens */
.sidebar { width: 300px; float: left; }
@media (max-width: 768px) {
  .sidebar { width: 100%; float: none; }  /* undo, undo, undo */
}

/* Mobile-first: add layout as space appears */
.sidebar { width: 100%; }
@media (min-width: 768px) {
  .sidebar { width: 300px; }  /* enhance when there's room */
}
```

:::key
Mobile-first means your **base styles are for small screens** with no media query. You then use `min-width` media queries to layer on complexity for larger viewports. Less CSS, fewer overrides, cleaner code.
:::

## min-width media queries

The building block of mobile-first. A `min-width` query activates when the viewport is *at least* that wide:

```css
/* Base: single column, no query needed */
.grid { display: flex; flex-direction: column; gap: 16px; }

/* Tablets and up */
@media (min-width: 640px) {
  .grid { flex-direction: row; flex-wrap: wrap; }
  .grid > * { flex: 1 1 45%; }
}

/* Desktops and up */
@media (min-width: 1024px) {
  .grid > * { flex: 1 1 30%; }
}
```

Each breakpoint *adds* to the previous one. You never reset or undo — you build upward.

## Common breakpoints

There's no universal set, but these cover the vast majority of devices:

| Breakpoint | Typical target          |
|------------|-------------------------|
| 640px      | Large phones / small tablets |
| 768px      | Tablets (portrait)      |
| 1024px     | Tablets (landscape) / small laptops |
| 1280px     | Desktops                |
| 1536px     | Large desktops          |

:::tip
Don't chase exact device widths — there are hundreds of screen sizes. Pick 2–3 breakpoints where your *design* actually breaks and needs adjustment. Let the content guide the breakpoints, not the devices.
:::

## Progressive enhancement in practice

Progressive enhancement means every user gets a working experience, and users with more capable browsers or wider screens get *more*.

### Navigation example

On mobile, a hamburger menu. On desktop, a full horizontal nav:

```css
/* Base: stack links vertically, hide by default */
.nav-links {
  display: none;
  flex-direction: column;
  gap: 8px;
}
.nav-links.open { display: flex; }

.hamburger { display: block; }

/* Desktop: show links inline, hide the hamburger */
@media (min-width: 768px) {
  .nav-links {
    display: flex;
    flex-direction: row;
    gap: 24px;
  }
  .hamburger { display: none; }
}
```

### Card grid example

Single column on phones, two columns on tablets, three on desktops:

```css
.cards {
  display: grid;
  gap: 16px;
  grid-template-columns: 1fr;           /* 1 column */
}

@media (min-width: 640px) {
  .cards { grid-template-columns: repeat(2, 1fr); }   /* 2 columns */
}

@media (min-width: 1024px) {
  .cards { grid-template-columns: repeat(3, 1fr); }   /* 3 columns */
}
```

## Designing for touch targets

Mobile users tap with fingers, not cursors. Small click targets cause frustration and accessibility failures.

```css
.btn {
  min-height: 44px;     /* Apple's recommended minimum */
  min-width: 44px;
  padding: 12px 24px;
  font-size: 16px;      /* prevents iOS zoom on input focus */
}

/* Give links breathing room in nav lists */
.nav-link {
  display: block;
  padding: 12px 16px;   /* generous tap area */
}
```

:::warning
A 44×44px minimum for interactive elements isn't just a suggestion — it's an accessibility guideline (WCAG). Tiny buttons and links are the number one mobile usability complaint. Always test by *actually tapping* on a real phone.
:::

### Input fields

On iOS, inputs with a `font-size` below 16px trigger an automatic zoom that disorients users:

```css
input, select, textarea {
  font-size: 16px;      /* prevents auto-zoom on iOS */
  padding: 12px;
  border-radius: 8px;
}
```

## Responsive images strategy

Images are often the heaviest assets on a page. Sending a 2000px hero image to a 375px phone wastes bandwidth and slows load times.

### The `srcset` approach

Let the browser pick the right image size:

```html
<img
  src="hero-800.jpg"
  srcset="hero-400.jpg 400w,
          hero-800.jpg 800w,
          hero-1600.jpg 1600w"
  sizes="(min-width: 1024px) 50vw, 100vw"
  alt="Hero banner"
/>
```

- `srcset` lists available image files and their widths.
- `sizes` tells the browser how wide the image will display at each breakpoint.
- The browser does the math and downloads only what it needs.

### The `<picture>` element

For art direction — showing a different crop on different screens:

```html
<picture>
  <source media="(min-width: 1024px)" srcset="hero-wide.jpg" />
  <source media="(min-width: 640px)" srcset="hero-medium.jpg" />
  <img src="hero-mobile.jpg" alt="Hero banner" />
</picture>
```

### CSS background images

```css
.hero {
  background-image: url('hero-mobile.jpg');
  background-size: cover;
}

@media (min-width: 1024px) {
  .hero { background-image: url('hero-desktop.jpg'); }
}
```

:::tip
Use modern formats like **WebP** or **AVIF** inside `<picture>` with a JPEG fallback. They're 25–50% smaller at the same quality, which makes a huge difference on mobile connections.
:::

## Testing workflow

Building mobile-first is pointless if you only test on your laptop monitor. A solid testing workflow:

1. **Browser DevTools** — Chrome and Firefox have responsive mode (toggle device toolbar). Test at 375px, 768px, and 1280px as a minimum.
2. **Throttle the network** — In DevTools' Network tab, simulate "Slow 3G" to feel what your users feel.
3. **Test on a real phone** — Emulators miss real-world issues like touch behavior, font rendering, and actual network speed.
4. **Check both orientations** — Portrait and landscape can produce different breakpoint results.

### Quick DevTools checklist

```
✓ Toggle device toolbar (Ctrl+Shift+M / Cmd+Shift+M)
✓ Select a phone preset (iPhone SE, Pixel 7)
✓ Drag the viewport edge to find where your layout breaks
✓ Check at least: 320px, 375px, 768px, 1024px, 1440px
✓ Rotate to landscape
```

## A complete mobile-first page skeleton

Putting it all together — a simple page layout that scales gracefully:

```css
/* === Base (mobile) === */
* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: system-ui, sans-serif;
  line-height: 1.6;
  padding: 16px;
}

.header { padding: 16px 0; }
.main   { margin-top: 24px; }
.sidebar { margin-top: 24px; }

/* === Tablet and up === */
@media (min-width: 768px) {
  body { padding: 24px 32px; }

  .page-layout {
    display: grid;
    grid-template-columns: 1fr 280px;
    gap: 32px;
  }
  .sidebar { margin-top: 0; }
}

/* === Desktop and up === */
@media (min-width: 1280px) {
  body { max-width: 1200px; margin-inline: auto; }
}
```

At mobile sizes everything stacks. At 768px a sidebar appears beside the main content. At 1280px the page caps its width and centers itself. Three layers, zero overrides.

:::quiz
Q: In a mobile-first approach, which type of media query do you use?
- max-width
- min-width *
- exact-width
E: Mobile-first uses `min-width` queries — the base styles target mobile, and each query adds enhancements for progressively wider viewports.
:::

:::quiz
Q: Why should touch targets be at least 44×44 pixels?
- To match the CSS box model default
- Because fingers are imprecise — small targets cause tap errors and accessibility failures *
- To prevent iOS auto-zoom
E: Fingers are much larger and less precise than mouse cursors. The 44×44px minimum (from Apple's HIG and WCAG guidelines) ensures users can reliably tap interactive elements.
:::

## Recap

- **Mobile-first** means base styles for small screens, then `min-width` media queries to enhance for larger viewports.
- Pick breakpoints based on where your **design breaks**, not on specific devices.
- **Progressive enhancement**: every screen gets a working experience; wider screens get richer layouts.
- Make touch targets at least **44×44px** and keep input font size at **16px** to avoid iOS zoom.
- Use `srcset`, `<picture>`, and modern image formats to serve appropriately sized images.
- **Test on real devices**, throttle the network, and check both orientations.

**Next up:** CSS Transitions — making your interfaces feel alive with smooth property changes.

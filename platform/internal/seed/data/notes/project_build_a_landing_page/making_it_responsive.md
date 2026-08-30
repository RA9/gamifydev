# Making it Responsive

Your landing page looks great on a wide desktop monitor. Now open it on a phone — it probably doesn't. Cards overflow, text is tiny, buttons are impossible to tap. This lesson fixes all of that using a mobile-first approach.

## The mobile-first mindset

**Mobile-first** means you write your base CSS for the smallest screen, then add complexity for larger screens with `min-width` media queries.

Why mobile-first?

- Mobile is the harder constraint. Start there and it forces you to prioritize content.
- Adding layout is easier than removing it.
- Most visitors are on phones. Design for them first.

```css
/* Base styles = mobile (single column, full width) */
.feature-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-lg);
}

/* Wider screens get multiple columns */
@media (min-width: 700px) {
  .feature-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

:::key
Mobile-first means your default CSS targets small screens. Media queries with `min-width` add layout for larger screens. This is the opposite of writing desktop CSS and then "fixing" mobile with `max-width` overrides.
:::

## Stacking the feature grid

On mobile, three side-by-side cards are unreadable. Stack them vertically by defaulting to `1fr`, then expanding at your breakpoint.

```css
.feature-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-lg);
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 var(--space-md);
}

@media (min-width: 700px) {
  .feature-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 960px) {
  .feature-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

This gives you:
- **< 700px** — 1 column (phones)
- **700–959px** — 2 columns (tablets)
- **960px+** — 3 columns (desktop)

## Adjusting font sizes with clamp()

Instead of writing media queries for every heading size, use `clamp()` to let the font size scale fluidly between a minimum and maximum.

```css
.hero h1 {
  font-size: clamp(1.75rem, 4vw + 0.5rem, 2.5rem);
}

.features h2 {
  font-size: clamp(1.375rem, 3vw + 0.25rem, 1.75rem);
}
```

The syntax is `clamp(minimum, preferred, maximum)`:
- **Minimum:** never smaller than this (good for readability)
- **Preferred:** scales with the viewport (`vw` units)
- **Maximum:** never larger than this (prevents absurd sizes on ultrawide monitors)

:::tip
`clamp()` replaces most font-size media queries. One line handles phones, tablets, and desktops. The formula `clamp(min, Xvw + Yrem, max)` covers nearly every case.
:::

## Hamburger nav pattern (simplified)

On desktop, horizontal nav links work fine. On mobile, they overflow. A simplified hamburger menu using a checkbox hack (no JavaScript) works well for a landing page.

```html
<header class="site-header">
  <div class="logo">Brewly</div>
  <input type="checkbox" id="nav-toggle" class="nav-toggle" />
  <label for="nav-toggle" class="nav-toggle-label" aria-label="Toggle menu">
    <span></span>
  </label>
  <nav class="nav">
    <a href="#features">Features</a>
    <a href="#testimonials">Testimonials</a>
    <a href="#" class="nav-cta">Sign up</a>
  </nav>
</header>
```

```css
.nav-toggle {
  display: none;    /* hide the checkbox */
}

.nav-toggle-label {
  display: none;    /* hidden on desktop */
  cursor: pointer;
  padding: var(--space-sm);
}

/* The hamburger icon — three lines */
.nav-toggle-label span,
.nav-toggle-label span::before,
.nav-toggle-label span::after {
  display: block;
  width: 24px;
  height: 2px;
  background: var(--color-ink);
  position: relative;
}
.nav-toggle-label span::before,
.nav-toggle-label span::after {
  content: '';
  position: absolute;
}
.nav-toggle-label span::before { top: -7px; }
.nav-toggle-label span::after  { top:  7px; }

/* Mobile: show hamburger, hide nav by default */
@media (max-width: 699px) {
  .nav-toggle-label {
    display: block;
  }

  .nav {
    display: none;
    width: 100%;
    text-align: center;
  }

  .nav a {
    display: block;
    padding: var(--space-sm) 0;
  }

  /* When checkbox is checked, show nav */
  .nav-toggle:checked ~ .nav {
    display: block;
  }
}
```

:::warning
The checkbox hack is great for simple static landing pages. For production sites with complex navigation, use JavaScript for better accessibility — managing `aria-expanded`, focus trapping, and Escape-key handling.
:::

## Testing at common breakpoints

Don't just test at your media query breakpoints. Check these common real-world widths:

| Width   | Device                     |
|---------|----------------------------|
| 320px   | Small phones (iPhone SE)   |
| 375px   | iPhone 12/13/14            |
| 414px   | Larger phones              |
| 768px   | iPad portrait              |
| 1024px  | iPad landscape / small laptop |
| 1280px  | Standard laptop            |
| 1440px  | Large desktop              |

Use DevTools device toolbar (**Cmd/Ctrl + Shift + M**) to drag through these widths. Look for:

- Text that overflows its container
- Buttons too close together or too small to tap
- Images that stretch or distort
- Horizontal scrollbars (a dead giveaway something is too wide)

## Fluid images

Images should never overflow their container. This one rule handles 90% of image responsiveness:

```css
img {
  max-width: 100%;
  height: auto;
  display: block;
}
```

For hero background images, use `background-size: cover` to fill the section regardless of screen size:

```css
.hero {
  background-image: url('hero-bg.jpg');
  background-size: cover;
  background-position: center;
}
```

:::tip
Always set `width` and `height` attributes on `<img>` tags in the HTML. This lets the browser reserve the right amount of space before the image loads, preventing layout shift (CLS).
:::

## Touch target sizes

On mobile, fingers are less precise than mouse pointers. Interactive elements need to be large enough to tap comfortably.

**Minimum touch target: 44×44px** (Apple's guideline) or **48×48px** (Google's recommendation).

```css
/* Ensure buttons are tappable */
.button {
  min-height: 44px;
  padding: 0.75rem 1.5rem;
}

/* Nav links need breathing room on mobile */
@media (max-width: 699px) {
  .nav a {
    display: block;
    padding: var(--space-md) var(--space-sm);
    min-height: 44px;
  }
}

/* Footer links too */
.footer-nav a {
  display: inline-block;
  padding: var(--space-sm);
  min-height: 44px;
}
```

## Responsive section padding

Sections need less horizontal padding on mobile (the screen is narrow) and more vertical padding on desktop (the screen is wide and content needs breathing room).

```css
section {
  padding: var(--space-xl) var(--space-md);
}

@media (min-width: 700px) {
  section {
    padding: var(--space-xl) var(--space-lg);
  }
}

@media (min-width: 960px) {
  section {
    padding: var(--space-2xl) var(--space-xl);
  }
}
```

:::quiz
Q: What does "mobile-first" mean in CSS?
- Writing all your CSS in a mobile app
- Writing base styles for small screens and using min-width media queries to add layout for larger screens *
- Only supporting mobile devices
- Loading a separate stylesheet for phones
E: Mobile-first means your default CSS is the mobile layout. You progressively enhance for wider viewports with `min-width` breakpoints, which is simpler and more maintainable than trying to undo desktop styles.
:::

:::quiz
Q: What does `clamp(1.75rem, 4vw + 0.5rem, 2.5rem)` do for a font size?
- Sets the font to exactly 4vw on all screens
- Scales the font fluidly with the viewport, never going below 1.75rem or above 2.5rem *
- Only applies on screens wider than 1.75rem
- Clamps the text to a single line
E: `clamp()` takes a minimum, a preferred value that scales, and a maximum. The browser picks whichever is appropriate for the current viewport width, giving you fluid typography without media queries.
:::

## Recap

- **Mobile-first:** write base CSS for small screens, enhance with `min-width` media queries.
- **Stack grids** on mobile (1 column), expand to 2 or 3 columns at wider breakpoints.
- **`clamp()`** creates fluid font sizes without multiple media queries.
- A **checkbox-based hamburger nav** works for simple landing pages without JavaScript.
- **Test at real device widths** (320px through 1440px), not just your breakpoints.
- **Fluid images** use `max-width: 100%; height: auto` to prevent overflow.
- **Touch targets** should be at least 44×44px for comfortable tapping.

**Next up:** Adding Transitions and Polish — making the page feel alive with hover effects, scroll reveals, and finishing touches.

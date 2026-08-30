# Styling the Page with CSS

You have the HTML structure and a style guide. Now you'll bring the landing page to life with CSS — applying your color palette, building the layout, and making every section feel polished and intentional.

## Apply the palette with custom properties

Start your stylesheet by translating the style guide into CSS custom properties. Define them once on `:root` and reference them everywhere.

```css
:root {
  --color-primary: #b5651d;
  --color-primary-dark: #7c3f10;
  --color-ink: #2b2118;
  --color-muted: #6b5d4f;
  --color-bg: #fdf8f3;
  --color-card: #ffffff;

  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 2rem;
  --space-xl: 4rem;
  --space-2xl: 8rem;

  --radius: 12px;
  --shadow-card: 0 4px 14px rgba(0, 0, 0, 0.05);
  --shadow-hover: 0 8px 24px rgba(0, 0, 0, 0.1);
}
```

:::key
Custom properties aren't just convenient — they enforce consistency. When you write `var(--color-primary)` instead of `#b5651d`, you can't accidentally use `#b5641d` and wonder why two browns don't match.
:::

## Reset and base typography

Set a clean foundation before styling individual sections.

```css
*, *::before, *::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', system-ui, sans-serif;
  color: var(--color-ink);
  background: var(--color-bg);
  line-height: 1.6;
}

h1, h2, h3 {
  font-family: 'DM Serif Display', serif;
  line-height: 1.2;
}

h1 { font-size: 2.5rem; }
h2 { font-size: 1.75rem; }
h3 { font-size: 1.25rem; }

img {
  max-width: 100%;
  display: block;
}
```

This gives you a consistent type hierarchy across every section.

## Style the hero — full-width, centered text, background

The hero is the first thing visitors see. It should feel spacious, clear, and confident.

```css
.hero {
  text-align: center;
  padding: var(--space-2xl) var(--space-md);
  background: linear-gradient(180deg, #fff6ec, var(--color-bg));
}

.hero h1 {
  max-width: 18ch;
  margin: 0 auto var(--space-md);
  color: var(--color-ink);
}

.hero p {
  max-width: 45ch;
  margin: 0 auto var(--space-lg);
  color: var(--color-muted);
  font-size: 1.125rem;
}
```

:::tip
The `ch` unit is based on the width of the character "0" in the current font. Setting `max-width: 45ch` on a paragraph keeps line length readable — roughly 45–75 characters per line is the sweet spot for comfortable reading.
:::

### The CTA button

The primary button should stand out. Use your brand color with enough padding to make it a comfortable tap target.

```css
.button {
  display: inline-block;
  background: var(--color-primary);
  color: #fff;
  padding: 0.8rem 1.6rem;
  border: none;
  border-radius: var(--radius);
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: background 0.2s ease;
}

.button:hover {
  background: var(--color-primary-dark);
}
```

## Feature grid with CSS Grid

CSS Grid makes equal-column card layouts trivial. Three columns, equal widths, consistent gaps.

```css
.features {
  padding: var(--space-xl) var(--space-lg);
  text-align: center;
}

.features h2 {
  margin-bottom: var(--space-lg);
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-lg);
  max-width: 1000px;
  margin: 0 auto;
}

.feature {
  background: var(--color-card);
  padding: var(--space-lg);
  border-radius: var(--radius);
  box-shadow: var(--shadow-card);
  text-align: left;
}

.feature h3 {
  color: var(--color-primary);
  margin-bottom: var(--space-sm);
}

.feature p {
  color: var(--color-muted);
}
```

:::key
`repeat(3, 1fr)` creates three columns of equal fractional width. `gap` handles all spacing between items — no margin hacks, no `:last-child` overrides.
:::

## Testimonial card with box-shadow

Testimonials build trust. A subtle card with a quote, name, and role is all you need.

```css
.testimonials {
  padding: var(--space-xl) var(--space-lg);
  background: var(--color-card);
}

.testimonial-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--space-lg);
  max-width: 900px;
  margin: 0 auto;
}

.testimonial {
  padding: var(--space-lg);
  border-radius: var(--radius);
  background: var(--color-bg);
  box-shadow: var(--shadow-card);
}

.testimonial blockquote {
  font-style: italic;
  margin-bottom: var(--space-md);
  color: var(--color-ink);
}

.testimonial .author {
  font-weight: 600;
  color: var(--color-primary);
}

.testimonial .role {
  font-size: 0.875rem;
  color: var(--color-muted);
}
```

The `auto-fit` with `minmax` pattern lets testimonial cards wrap naturally — two columns on desktop, one on mobile, with no media query required.

## CTA section with strong contrast

The final CTA block should feel visually distinct from everything above. Use a dark or brand-colored background to create contrast.

```css
.cta-section {
  text-align: center;
  padding: var(--space-2xl) var(--space-md);
  background: var(--color-ink);
  color: #fff;
}

.cta-section h2 {
  margin-bottom: var(--space-md);
  color: #fff;
}

.cta-section p {
  max-width: 40ch;
  margin: 0 auto var(--space-lg);
  color: #d8ccbf;
}

.cta-section .button {
  background: var(--color-primary);
  font-size: 1.125rem;
  padding: 1rem 2rem;
}
```

:::warning
When using light text on a dark background, check contrast for both the heading *and* the body text. It's easy for muted text like `#d8ccbf` on `#2b2118` to fall below the 4.5:1 WCAG ratio.
:::

## Consistent spacing with the scale

Every vertical gap on the page should use values from your spacing scale. Here's a mental model:

| Scale value    | Where to use it                        |
|----------------|----------------------------------------|
| `--space-sm`   | Between a heading and its subtext       |
| `--space-md`   | Between elements inside a card          |
| `--space-lg`   | Between cards, between section heading and content |
| `--space-xl`   | Section padding (top/bottom)            |
| `--space-2xl`  | Hero padding, major visual breaks       |

```css
/* Good: every value comes from the scale */
.features       { padding: var(--space-xl) var(--space-lg); }
.features h2    { margin-bottom: var(--space-lg); }
.feature h3     { margin-bottom: var(--space-sm); }
```

## Typography hierarchy in practice

Make sure a visitor can scan the page and understand its structure purely from text size and weight:

```css
/* Hero: largest, darkest */
.hero h1 { font-size: 2.5rem; font-weight: 700; color: var(--color-ink); }

/* Section headings: medium, still prominent */
.features h2 { font-size: 1.75rem; font-weight: 700; }

/* Card headings: smaller, accented */
.feature h3 { font-size: 1.25rem; color: var(--color-primary); }

/* Body text: base size, muted */
p { font-size: 1rem; color: var(--color-muted); }

/* Small text: captions, meta */
.role { font-size: 0.875rem; }
```

The reader should never wonder which text is most important — the hierarchy makes it obvious.

## Footer

```css
.site-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: var(--space-md);
  padding: var(--space-lg);
  background: var(--color-ink);
  color: #d8ccbf;
  font-size: 0.875rem;
}

.footer-nav a {
  color: #d8ccbf;
  text-decoration: none;
  margin-left: var(--space-md);
}

.footer-nav a:hover {
  color: #fff;
}
```

:::quiz
Q: Why define colors as CSS custom properties instead of typing hex values directly?
- Custom properties make colors render faster in the browser
- They ensure consistency and let you change a color site-wide from one place *
- Browsers require custom properties for dark backgrounds
- They automatically check contrast ratios
E: Custom properties create a single source of truth. Change `--color-primary` in `:root` and every element using it updates. This prevents accidental inconsistencies and makes theming easy.
:::

:::quiz
Q: What does `grid-template-columns: repeat(3, 1fr)` do?
- Creates 3 rows of equal height
- Creates 3 columns that each take one equal fraction of the available width *
- Repeats the grid content 3 times
- Sets the font size to 3 fractional units
E: `1fr` means one fractional unit of the remaining space. Three `1fr` columns split the container into three equal-width columns.
:::

## Recap

- **Custom properties** on `:root` centralize your palette, spacing, shadows, and radii.
- **Base styles** set the type hierarchy and a clean reset before section-specific CSS.
- The **hero** uses centered text, generous padding, and a gradient background.
- The **feature grid** uses `repeat(3, 1fr)` with `gap` for clean card layouts.
- **Testimonial cards** use `box-shadow` and `auto-fit` with `minmax` for flexible wrapping.
- The **CTA section** uses a contrasting dark background to stand out.
- Every spacing value comes from the **spacing scale** — no magic numbers.
- The **typography hierarchy** ensures scannability through deliberate size and weight choices.

**Next up:** Making it Responsive — adapting this layout to work beautifully on phones and tablets.

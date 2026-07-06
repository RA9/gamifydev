# Responsive Design & Media Queries

A page that only works on your laptop is not finished.

Responsive design means your interface adapts gracefully across:

- phones
- tablets
- laptops
- large desktop screens

This lesson is about designing for different viewports without creating separate websites.

## Start mobile-first

A strong default approach is:

1. build a clean mobile layout first
2. add more room and enhancements as the screen gets larger

Why this works:

- mobile forces focus and simplicity
- your base styles stay lean
- larger layouts are added intentionally instead of hacked backward

```css
.hero {
  display: grid;
  gap: 1rem;
}

@media (min-width: 768px) {
  .hero {
    grid-template-columns: 1fr 1fr;
    align-items: center;
  }
}
```

## Flexible sizing beats fixed sizing

Responsive layouts usually depend on flexible values.

Good patterns:

```css
.container {
  width: min(100% - 2rem, 70rem);
  margin-inline: auto;
}

img {
  max-width: 100%;
  height: auto;
}
```

These rules do a lot of work:

- the container never overflows the viewport
- the image never bursts out of its parent
- the layout feels fluid instead of rigid

## What media queries do

A media query says: *when the viewport matches this condition, apply these styles too.*

```css
@media (min-width: 768px) {
  .cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
```

You are not rewriting the whole layout. You are **enhancing** it for more available space.

:::fill
Q: Complete the media query so the styles apply at tablet size and above.
`@media (___: 768px) { ... }`
- min-width *
- max-height
- min-height
E: `min-width` is a common mobile-first breakpoint pattern: apply extra layout when there is at least this much horizontal space.
:::

## Typical things that change at larger sizes

As screen size grows, you might:

- switch stacked content into columns
- increase white space
- enlarge headings slightly
- show more navigation options in a row
- move cards into 2, 3, or 4 columns

Example:

```css
.features {
  display: grid;
  gap: 1rem;
}

@media (min-width: 640px) {
  .features {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 960px) {
  .features {
    grid-template-columns: repeat(4, 1fr);
  }
}
```

## Responsiveness is not just columns

A responsive interface also considers:

- readable line length
- tap target size
- image behavior
- spacing density
- whether content order still makes sense

If a button becomes too small to tap or a paragraph becomes 160 characters wide, the page may be technically responsive but still poor to use.

## Keep media and text under control

Helpful defaults:

```css
img,
picture,
video {
  max-width: 100%;
}

body {
  line-height: 1.6;
}
```

For text, avoid stretches that are too wide to scan comfortably. Containers and max-widths help a lot.

## A practical breakpoint mindset

Don't ask: "What are the official screen sizes?"

Ask:

- When does this layout start to feel cramped?
- When do these cards deserve a second column?
- When does the nav need more space?

Breakpoints should respond to the **content**, not internet folklore.

:::quiz
Q: What is the healthier responsive mindset?
- Set one fixed desktop width and hope it works
- Add breakpoints when the content actually needs them *
- Make separate HTML pages for phone and desktop
E: Strong responsive design is content-driven. Add layout changes when the UI needs them, not just because a random device exists.
:::

## Mini practice — adapt a card section

Start with a simple mobile layout:

```css
.cards {
  display: grid;
  gap: 1rem;
}
```

Then enhance it:

```css
@media (min-width: 700px) {
  .cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .cards {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

You now have:

- one column on small screens
- two columns on medium screens
- three columns on large screens

That is responsive design in practice.

## Common responsive mistakes to avoid

- fixed widths everywhere
- large images with no max-width
- tiny tap targets
- giant desktop navs crammed into phone screens
- too many breakpoints too early
- designing only by dragging the browser manually without checking actual content behavior

:::warning
If you only test a layout at one width, you're not testing the layout — you're testing a screenshot.
:::

## What good looks like

You should now be able to:

- explain what mobile-first means
- use `min-width` media queries
- create fluid containers and media
- change layout from stacked to multi-column at sensible points
- think about readability and usability, not just geometry

## What's next

In **Lab: Build a Responsive Landing Page**, you'll apply all of this to a real frontend build with a hero, feature grid, testimonial section, and call-to-action.
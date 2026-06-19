# CSS Basics

**CSS** (Cascading Style Sheets) controls how your HTML *looks* — colors, fonts, spacing, and layout. HTML is the skeleton; CSS is the clothing.

## How a CSS rule works

A rule has a **selector** (what to style) and a **declaration block** (how to style it):

- Selector: `p` — every paragraph.
- Declaration: `color: blue;` — make the text blue.

Put together: `p { color: blue; font-size: 18px; }`

## Ways to apply CSS

- **Inline** — a `style` attribute on one element (use sparingly).
- **Internal** — a `<style>` block in the `<head>`.
- **External** — a separate `.css` file linked with `<link rel="stylesheet" href="main.css">`. This is the recommended approach.

## Selectors you will use constantly

- **Element**: `h1 { ... }`
- **Class**: `.button { ... }` targets `<button class="button">`. Reusable.
- **ID**: `#header { ... }` targets one unique element.

## The box model

Every element is a box made of four layers, from inside out:

- **Content** — the text or image.
- **Padding** — space *inside* the box, around the content.
- **Border** — a line around the padding.
- **Margin** — space *outside* the box, between it and its neighbors.

Understanding the box model is the key to controlling spacing and layout.

## Useful properties

- **Color**: `color`, `background-color`
- **Text**: `font-size`, `font-weight`, `text-align`
- **Spacing**: `margin`, `padding`
- **Layout**: `display`, and modern tools like **flexbox** (`display: flex`) and **grid** (`display: grid`)

## Key takeaways

- CSS rules pair a selector with declarations.
- Prefer external stylesheets and reusable classes.
- The box model (content, padding, border, margin) governs spacing.

Next up: **Building a Website with HTML and CSS**, where you combine both.

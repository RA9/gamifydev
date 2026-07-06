# Responsive Landing Page

Landing pages are a frontend staple because they test your ability to combine hierarchy, layout, responsiveness, and clear calls-to-action in one coherent build.

:::project
You'll build a responsive product landing page with a header, hero section, feature grid, testimonial block, and final CTA. The page should feel strong on mobile first, then expand cleanly on larger screens.
:::

## Step 1 — Build the page structure

Start with a semantic skeleton and section names you can reason about.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Sprintboard</title>
</head>
<body>
  <header class="site-header"></header>
  <main>
    <section class="hero"></section>
    <section class="features"></section>
    <section class="testimonial"></section>
    <section class="cta"></section>
  </main>
  <footer class="site-footer"></footer>
</body>
</html>
```

Fill each section with real copy before worrying about styling details.

## Step 2 — Add global CSS foundations

Create a `<style>` block or external stylesheet with defaults that make later work easier.

```css
* { box-sizing: border-box; }
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  line-height: 1.6;
  color: #0f172a;
  background: #f8fafc;
}
img {
  max-width: 100%;
  height: auto;
  display: block;
}
.container {
  width: min(100% - 2rem, 70rem);
  margin-inline: auto;
}
```

These rules give you predictable sizing, readable text, and responsive images.

## Step 3 — Build the mobile hero first

Create the hero content as a single-column mobile layout.

```html
<section class="hero container">
  <div class="hero-copy">
    <h1>Plan faster with less process</h1>
    <p>Sprintboard helps small teams turn scattered ideas into focused weekly work.</p>
    <a class="button" href="#signup">Start free</a>
  </div>
  <img src="dashboard-preview.png" alt="Sprintboard dashboard showing weekly priorities and tasks" />
</section>
```

```css
.hero {
  display: grid;
  gap: 2rem;
  padding: 4rem 0;
}
.button {
  display: inline-block;
  padding: 0.875rem 1.25rem;
  border-radius: 999px;
  background: #5b3df5;
  color: white;
  text-decoration: none;
  font-weight: 700;
}
```

## Step 4 — Turn features into a card grid

Add three or four feature cards and lay them out with CSS Grid.

```css
.features-grid {
  display: grid;
  gap: 1rem;
}
.feature-card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  padding: 1.25rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}
```

A responsive upgrade can come later with a media query.

## Step 5 — Add testimonial and CTA contrast

Give the page rhythm by making sections feel different.

```css
.testimonial {
  background: #0f172a;
  color: white;
  border-radius: 1.5rem;
  padding: 2rem;
}

.cta {
  text-align: center;
  padding: 4rem 0;
}
```

A landing page should guide the eye, not feel like one endless block of identical spacing.

## Step 6 — Add responsive breakpoints

At a larger width, upgrade the layout.

```css
@media (min-width: 768px) {
  .hero {
    grid-template-columns: 1.1fr 1fr;
    align-items: center;
  }

  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .features-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
```

The goal is not to invent a new page. It's to let the same page breathe with more space.

## Step 7 — Polish the details and review

Finish with:

- hover states on buttons and links
- consistent section spacing
- strong heading hierarchy
- a footer with simple contact or copyright info

Final review checklist:

- no horizontal overflow
- CTA is visible and obvious
- images stay inside their containers
- the page feels readable on both phone and desktop

When that all works, you've built a real responsive marketing page.
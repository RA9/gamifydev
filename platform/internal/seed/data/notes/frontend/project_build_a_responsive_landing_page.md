# Project: Build a Responsive Landing Page

This project is where your HTML and CSS start to feel like real product work.

:::project
**Goal:** Build a polished marketing page for a fictional product. It should have a header, hero, features, testimonial, CTA, and footer. It must look good on mobile first, then expand into a stronger desktop layout.
:::

## What you'll practise

- semantic page structure
- visual hierarchy
- responsive layout with Grid and Flexbox
- reusable card/button styling
- mobile-first CSS
- QA and polish before shipping

## The page brief

Build a landing page for **Sprintboard**, a planning tool for small teams.

Your page should include:

- a header with brand name and nav links
- a hero section with headline, short copy, and call to action
- a features section with 3–4 cards
- a testimonial section
- a final CTA section
- a footer

## Step 1 — Plan the sections before coding

Before you write code, map the page into regions:

- `<header>`
- `<main>`
  - hero
  - features
  - testimonial
  - CTA
- `<footer>`

This gives you structure before you get distracted by visual detail.

## Step 2 — Build the HTML skeleton

Start with semantic HTML and real content.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Sprintboard</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
  <header class="site-header">
    <nav class="container nav">
      <a class="brand" href="#">Sprintboard</a>
      <div class="nav-links">
        <a href="#features">Features</a>
        <a href="#stories">Stories</a>
        <a href="#signup">Start free</a>
      </div>
    </nav>
  </header>

  <main>
    <section class="hero container">
      <div class="hero-copy">
        <h1>Plan faster with less process</h1>
        <p>Turn scattered ideas into focused weekly work with a board your whole team can understand in minutes.</p>
        <a class="btn btn-primary" href="#signup">Start free</a>
      </div>
      <div class="hero-media" aria-hidden="true">Dashboard preview</div>
    </section>
  </main>

  <footer class="site-footer">
    <div class="container">
      <p>© 2026 Sprintboard</p>
    </div>
  </footer>
</body>
</html>
```

## Step 3 — Add global layout foundations

Create a small design system up front.

```css
:root {
  --bg: #f8fafc;
  --card: #ffffff;
  --ink: #0f172a;
  --muted: #64748b;
  --brand: #5b3df5;
  --brand-dark: #4630c9;
  --line: #e2e8f0;
  --radius: 18px;
}

* { box-sizing: border-box; }
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  line-height: 1.6;
  color: var(--ink);
  background: var(--bg);
}

.container {
  width: min(100% - 2rem, 72rem);
  margin-inline: auto;
}

a { color: inherit; text-decoration: none; }
```

These simple defaults make everything else easier.

## Step 4 — Build the mobile layout first

Start with a single-column layout that reads clearly on a phone.

```css
.hero {
  display: grid;
  gap: 2rem;
  padding: 4rem 0;
}

.hero-media {
  min-height: 220px;
  border-radius: var(--radius);
  background: linear-gradient(135deg, #dbeafe, #ede9fe);
  display: grid;
  place-items: center;
  color: var(--muted);
  border: 1px solid var(--line);
}

.btn {
  display: inline-block;
  padding: 0.85rem 1.2rem;
  border-radius: 999px;
  font-weight: 700;
}

.btn-primary {
  background: var(--brand);
  color: white;
}
```

Make sure the page already feels complete at small widths before you add complexity.

## Step 5 — Add feature cards and section rhythm

Feature cards are a great place to practise reusable components.

```html
<section id="features" class="container section">
  <h2>Why teams like Sprintboard</h2>
  <div class="feature-grid">
    <article class="card">
      <h3>Weekly priorities</h3>
      <p>Keep the whole team focused on what matters this week.</p>
    </article>
    <article class="card">
      <h3>Clear ownership</h3>
      <p>Every task has an owner, status, and due date.</p>
    </article>
    <article class="card">
      <h3>Fast setup</h3>
      <p>Go from empty board to usable workflow in minutes.</p>
    </article>
  </div>
</section>
```

```css
.section { padding: 3rem 0; }
.feature-grid {
  display: grid;
  gap: 1rem;
}
.card {
  background: var(--card);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  padding: 1.25rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}
```

## Step 6 — Add responsive upgrades

Now give the page more structure when more space exists.

```css
@media (min-width: 768px) {
  .nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .hero {
    grid-template-columns: 1.15fr 1fr;
    align-items: center;
  }

  .feature-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

This is the key responsive mindset: keep the same content, then let the layout breathe at larger widths.

## Step 7 — Polish and review like a frontend developer

Before calling it done, check:

- is the CTA still obvious on small screens?
- do any sections feel cramped or too stretched?
- are spacing values consistent?
- do buttons and links have hover states?
- is there any horizontal overflow?

:::quiz
Q: What is usually the strongest starting point for a responsive landing page?
- A complex desktop-only layout
- A clean single-column mobile layout *
- Absolute positioning for every section
E: Mobile-first work keeps the content focused and the layout dependable, then scales up for larger screens.
:::

## Stretch goals

- add a logo strip under the hero
- include a pricing section
- add soft hover lift on cards
- create a dark section for the testimonial or final CTA

## What "done" looks like

A finished version should feel like a real product page, not a classroom exercise. It should be readable, responsive, and confident enough that you'd actually link to it.
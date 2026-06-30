# Building a Website with HTML and CSS

In this hands-on lesson you'll build a complete, responsive landing page from scratch — first the semantic HTML structure, then the CSS that makes it look good on any screen. By the end you'll have a real page you can open in a browser and a clear mental model of how the two languages work together.

We'll build a landing page for a fictional product called **Brewly**, a coffee-subscription service.

## Step 1: Plan the page

Before touching code, sketch the sections. A simple landing page usually has:

- A **header** with a logo and navigation.
- A **hero** — a big headline, a sentence, and a call-to-action button.
- A **features** section — a few cards explaining what's great.
- A **footer** with copyright and links.

:::tip
Always plan structure before styling. If your HTML is organized into clear sections, the CSS practically writes itself.
:::

## Step 2: The HTML skeleton

Every page starts with the same boilerplate. The `<!DOCTYPE html>` tells the browser this is modern HTML, and the `<meta name="viewport">` line is essential for responsive design on phones.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Brewly — Fresh Coffee, Delivered</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
  <!-- Page sections go here -->
</body>
</html>
```

The `<link>` tag connects an external `style.css` file. Keeping CSS in its own file keeps things tidy.

## Step 3: The semantic HTML

"Semantic" means using tags that describe *what content is*, not just how it looks. `<header>`, `<nav>`, `<section>`, and `<footer>` make your page easier to read, better for accessibility, and friendlier to search engines.

### The header and navigation

```html
<header class="site-header">
  <div class="logo">Brewly</div>
  <nav class="nav">
    <a href="#features">Features</a>
    <a href="#pricing">Pricing</a>
    <a href="#" class="nav-cta">Sign up</a>
  </nav>
</header>
```

### The hero

```html
<section class="hero">
  <h1>Fresh coffee, delivered to your door</h1>
  <p>Hand-roasted beans from around the world, on your schedule. Cancel anytime.</p>
  <a href="#" class="button">Start your subscription</a>
</section>
```

### The features section

We use a wrapping `<section>` with three feature "cards" inside.

```html
<section id="features" class="features">
  <h2>Why people love Brewly</h2>
  <div class="feature-grid">
    <article class="feature">
      <h3>Always fresh</h3>
      <p>Roasted within 48 hours of shipping, so every cup tastes alive.</p>
    </article>
    <article class="feature">
      <h3>Flexible plans</h3>
      <p>Weekly, biweekly, or monthly. Pause or cancel whenever you like.</p>
    </article>
    <article class="feature">
      <h3>Global beans</h3>
      <p>Single-origin coffee sourced ethically from trusted farms.</p>
    </article>
  </div>
</section>
```

### The footer

```html
<footer class="site-footer">
  <p>&copy; 2026 Brewly. All rights reserved.</p>
  <nav class="footer-nav">
    <a href="#">Privacy</a>
    <a href="#">Terms</a>
    <a href="#">Contact</a>
  </nav>
</footer>
```

Drop all four sections inside the `<body>` in that order and you have a complete, *unstyled* page. Open it in a browser — it works, it's just plain. CSS comes next.

:::quiz
Q: Why use `<header>` and `<section>` instead of plain `<div>` everywhere?
- They load faster
- They describe the meaning of the content, helping accessibility and SEO *
- They are required by every browser
E: Semantic tags communicate *what* the content is. Screen readers and search engines understand them, and your code is easier to read.
:::

## Step 4: CSS foundations — variables and resets

Now the `style.css` file. Start with a small **reset** so margins are predictable, and define **custom properties** (CSS variables) for your colors and spacing. Defining values once means you change them in one place later.

```css
:root {
  --brand: #b5651d;
  --brand-dark: #7c3f10;
  --ink: #2b2118;
  --muted: #6b5d4f;
  --bg: #fdf8f3;
  --card: #ffffff;
  --space: 1rem;
  --radius: 12px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: system-ui, sans-serif;
  color: var(--ink);
  background: var(--bg);
  line-height: 1.6;
}
```

:::key
`box-sizing: border-box` makes an element's declared width *include* its padding and border. Without it, padding pushes elements wider than you expect. Set it once on `*` and forget about it.
:::

## Step 5: Style the header with flexbox

**Flexbox** lays items out in a row (or column) and spaces them easily. Here we push the logo to the left and the nav to the right with `justify-content: space-between`.

```css
.site-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space) calc(var(--space) * 2);
  background: var(--card);
  border-bottom: 1px solid #eee;
}

.logo {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--brand);
}

.nav a {
  margin-left: var(--space);
  text-decoration: none;
  color: var(--ink);
}

.nav .nav-cta {
  color: var(--brand);
  font-weight: 600;
}
```

## Step 6: Style the hero and the button

Center the hero text and give it breathing room. The reusable `.button` class styles every call-to-action consistently.

```css
.hero {
  text-align: center;
  padding: calc(var(--space) * 5) var(--space);
  background: linear-gradient(180deg, #fff6ec, var(--bg));
}

.hero h1 {
  font-size: 2.5rem;
  max-width: 18ch;
  margin: 0 auto var(--space);
}

.hero p {
  color: var(--muted);
  max-width: 45ch;
  margin: 0 auto calc(var(--space) * 1.5);
}

.button {
  display: inline-block;
  background: var(--brand);
  color: #fff;
  padding: 0.8rem 1.6rem;
  border-radius: var(--radius);
  text-decoration: none;
  font-weight: 600;
  transition: background 0.2s ease;
}

.button:hover {
  background: var(--brand-dark);
}
```

## Step 7: Lay out the features with CSS Grid

**Grid** is perfect for card layouts. `repeat(3, 1fr)` makes three equal columns, and `gap` adds space between them — no margins needed.

```css
.features {
  padding: calc(var(--space) * 4) calc(var(--space) * 2);
  text-align: center;
}

.features h2 {
  margin-bottom: calc(var(--space) * 2);
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space);
  max-width: 1000px;
  margin: 0 auto;
}

.feature {
  background: var(--card);
  padding: calc(var(--space) * 1.5);
  border-radius: var(--radius);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.05);
  text-align: left;
}

.feature h3 {
  color: var(--brand);
  margin-bottom: 0.5rem;
}
```

## Step 8: Style the footer

```css
.site-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: var(--space);
  padding: calc(var(--space) * 2);
  background: var(--ink);
  color: #d8ccbf;
}

.footer-nav a {
  color: #d8ccbf;
  margin-left: var(--space);
  text-decoration: none;
}
```

## Step 9: Make it responsive

On a phone, three side-by-side cards are too cramped. A **media query** applies CSS only below a certain width. Here, when the screen is 700px or narrower, we stack the cards into one column and shrink the headline.

```css
@media (max-width: 700px) {
  .feature-grid {
    grid-template-columns: 1fr;
  }

  .hero h1 {
    font-size: 1.8rem;
  }

  .site-header {
    flex-direction: column;
    gap: 0.5rem;
  }
}
```

:::analogy
A media query is like a dress code that only applies on certain days. "When the screen is small, follow *these* rules instead." The rest of your CSS still applies as normal.
:::

:::quiz
Q: What does a media query like `@media (max-width: 700px)` do?
- Hides the page on small screens
- Applies its CSS rules only when the screen is 700px wide or narrower *
- Makes images load faster
E: Media queries apply styles conditionally based on screen size, which is how you make a layout adapt to phones, tablets, and desktops.
:::

:::predict
You change `grid-template-columns: repeat(3, 1fr)` to `repeat(2, 1fr)` on `.feature-grid`. What happens to the three feature cards on a wide screen?
ANSWER: Two cards sit in the first row and the third wraps to a new row below them, since the grid now only allows two columns.
:::

## Recap

- **Plan the sections** first: header, hero, features, footer.
- Use **semantic HTML** (`<header>`, `<nav>`, `<section>`, `<footer>`) for meaningful, accessible structure.
- Start CSS with a **reset** and **custom properties** (variables) for colors and spacing.
- **Flexbox** arranges items in a row; `justify-content: space-between` pushes them apart.
- **Grid** with `repeat(3, 1fr)` and `gap` builds clean card layouts.
- A reusable `.button` class and `:hover` states keep styling consistent and interactive.
- A **media query** stacks the layout on small screens, making the page responsive.

**Next up:** Project: Build a Profile Card — apply these skills to craft a polished, self-contained component.

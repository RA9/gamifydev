# Building a Website with HTML and CSS

You now have the core pieces:

- HTML for structure
- semantic thinking for meaning
- CSS for presentation
- flexbox and grid for layout
- media queries for responsiveness

This lesson is about combining them like a working frontend developer would: from brief to plan to implementation to QA.

## Real websites are built from decisions, not random code

A common beginner trap is opening a file and immediately styling details.

A stronger workflow is:

1. understand the page goal
2. list the content sections
3. write semantic HTML
4. add layout CSS
5. style components
6. test and refine

:::key
Professionals don't start with colors. They start with structure and content.
:::

## Start from a page brief

Imagine the brief says:

> Build a homepage for a small creative studio. The page needs a hero, a services section, a recent-work section, a testimonial, and a footer with contact details.

Before you write code, translate the brief into page regions:

- header / nav
- hero
- services
- work gallery
- testimonial
- footer

That gives you the structure of the whole build before a single visual detail is chosen.

## Step 1 — scaffold the HTML

Start with the page outline:

```html
<body>
  <header>
    <nav></nav>
  </header>

  <main>
    <section class="hero"></section>
    <section class="services"></section>
    <section class="work"></section>
    <section class="testimonial"></section>
  </main>

  <footer></footer>
</body>
```

Then fill each section with real headings, paragraphs, cards, links, images, and buttons.

## Step 2 — organize your files

A tidy starting structure might be:

```text
studio-site/
├── index.html
├── styles.css
└── images/
```

This seems basic, but good file hygiene becomes more important as projects grow.

## Step 3 — create layout foundations first

In `styles.css`, begin with broad rules:

```css
* {
  box-sizing: border-box;
}

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
```

These rules give you:

- predictable sizing
- sensible typography defaults
- responsive images

## Step 4 — lay out the big sections

Work from large structure to small details.

For example:

```css
.hero {
  display: grid;
  gap: 2rem;
  padding: 4rem 1.25rem;
}

.services-grid,
.work-grid {
  display: grid;
  gap: 1rem;
}

@media (min-width: 768px) {
  .hero {
    grid-template-columns: 1.1fr 1fr;
    align-items: center;
  }

  .services-grid,
  .work-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
```

Only after the page structure feels right should you obsess over shadows, hover states, or exact spacing tokens.

## Step 5 — style reusable components

Real sites repeat patterns. Instead of styling each card from scratch, create reusable component rules.

```css
.card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  padding: 1.25rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
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

This is how design consistency starts.

## Step 6 — check the build like a reviewer

A decent-looking page can still be weak.

Run through this checklist:

- Does the page have one clear `<h1>`?
- Are links and buttons used correctly?
- Do sections have enough spacing to scan easily?
- Does the layout still work on mobile?
- Are images responsive?
- Does the CTA stand out?
- Is the page still readable with no CSS?

:::quiz
Q: Which order usually produces stronger frontend work?
- Color palette → animations → content structure
- Structure → layout → components → polish *
- Hover states → icons → headings
E: Strong builds move from content structure to layout and reusable components, then finish with polish.
:::

## Mini build plan — from brief to finished page

Use this workflow on your next static site:

### Plan

- write down the sections
- decide the content for each section
- list repeated UI patterns (cards, buttons, nav links)

### Build HTML

- create the page skeleton
- add headings, copy, and media
- keep semantics clean

### Build CSS

- add global defaults
- build layout systems
- style components
- add responsive breakpoints

### QA

- resize the viewport
- test button and link clarity
- inspect spacing rhythm
- check image overflow

## Common mistakes to avoid

- styling before content exists
- using different spacing values everywhere
- hard-coding widths that break on smaller screens
- skipping semantic structure because the page "looks fine"
- building sections one at a time with no shared system

:::warning
If every section is built with completely different spacing, typography, and layout logic, the site will feel inconsistent even if each part looks okay in isolation.
:::

## What good looks like

You should now be able to:

- translate a page brief into content sections
- scaffold a semantic multi-section page
- organize files cleanly
- build broad layout systems before styling details
- create reusable components for consistency
- review a static site for structure, responsiveness, and polish

## What's next

In **Project: Build a Profile Card**, you'll zoom back into a smaller UI component and polish it properly — hierarchy, spacing, alignment, and finishing detail.
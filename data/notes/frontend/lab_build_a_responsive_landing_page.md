# Lab: Build a Responsive Landing Page

This is your first HTML/CSS build that should feel like something a real product team might ship.

:::project
You'll build a responsive landing page for a fictional product. It should read well on mobile, grow into a stronger layout on larger screens, and use the HTML/CSS skills you've learned so far.
:::

## The brief

Build a landing page for a product called **Sprintboard** — a lightweight planning tool for small teams.

Your page should include:

- a header with product name and nav links
- a hero section with heading, paragraph, and CTA button
- a feature section with three or four feature cards
- a testimonial section
- a final call-to-action section
- a footer

## Content goals

Your layout should communicate hierarchy clearly:

- one obvious primary heading
- supporting copy that is easy to scan
- cards that feel related but distinct
- a CTA that stands out visually

## Build constraints

Use:

- semantic HTML landmarks
- flexbox and/or grid for layout
- mobile-first CSS
- at least one media query
- responsive images or image placeholders that stay inside their containers

## Recommended build order

### 1. Structure the page in HTML

Start with the sections only.

```html
<header></header>
<main>
  <section class="hero"></section>
  <section class="features"></section>
  <section class="testimonial"></section>
  <section class="cta"></section>
</main>
<footer></footer>
```

Then fill each section with real content.

### 2. Build a clean mobile layout first

On mobile, let the page stack naturally.

A strong mobile default usually means:

- single-column layout
- generous vertical spacing
- easy-to-read text
- buttons wide enough to tap comfortably

### 3. Add visual hierarchy

Use CSS to create separation between sections:

- spacing
- background contrast
- card borders or shadows
- heading size and paragraph rhythm

### 4. Upgrade the layout at larger widths

At a tablet or desktop breakpoint, you might:

- place hero text and hero image side by side
- turn feature cards into a multi-column grid
- tighten or rebalance spacing

Example:

```css
.features-grid {
  display: grid;
  gap: 1rem;
}

@media (min-width: 768px) {
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

### 5. Review like a frontend developer

Check:

- can you read it comfortably at phone width?
- does anything overflow horizontally?
- do images stay inside the layout?
- does the hero still feel balanced on desktop?
- is the CTA still obvious on all screen sizes?

## Suggested HTML starter

```html
<header>
  <nav>
    <a href="/">Sprintboard</a>
    <a href="#features">Features</a>
    <a href="#pricing">Pricing</a>
    <a href="#contact">Contact</a>
  </nav>
</header>

<main>
  <section class="hero">
    <div>
      <h1>Plan faster with less process</h1>
      <p>Sprintboard helps small teams turn messy ideas into clear weekly priorities.</p>
      <a href="#signup">Start free</a>
    </div>
    <img src="dashboard-preview.png" alt="Sprintboard dashboard showing tasks, priorities, and weekly goals" />
  </section>
</main>
```

## Acceptance criteria

Your lab is complete when:

- the page uses semantic sections and landmarks
- the layout starts from a clean mobile version
- at least one section changes layout at a wider breakpoint
- the feature cards use a clear flex or grid layout
- no content overflows the viewport
- the CTA is visually distinct

:::quiz
Q: On a narrow phone screen, what is usually the safest starting layout for major page sections?
- Four columns side by side
- A single stacked column *
- Absolute positioning everywhere
E: A single-column mobile layout is usually the simplest, strongest starting point. Add complexity only when more space exists.
:::

## Stretch goals

If you finish early:

- add a pricing section
- create a logo strip under the hero
- add hover states to buttons and nav links
- use `max-width` containers to improve readability on desktop

## Reflection prompts

Explain:

- Why did you choose flexbox or grid for each section?
- Where did you place your breakpoint, and why?
- What changed between mobile and desktop layouts?

:::key
Responsive design is not about shrinking a desktop page. It's about designing a layout that stays usable and readable as space changes.
:::

## What's next

In **Building a Website with HTML and CSS**, you'll take the next step: planning and assembling a fuller multi-section website with a more professional workflow.
# Semantic HTML & Accessibility

A lot of developers learn tags. Strong frontend developers learn **why the tag choice matters**.

Semantic HTML improves:

- accessibility for screen-reader and keyboard users
- maintainability for your future self
- SEO and machine understanding
- CSS and JavaScript clarity

This lesson is about making better structural decisions, not memorizing more markup.

## What semantic HTML means

A semantic element tells the browser and the developer something about the role of the content.

Examples:

- `<header>` — introductory content for a page or section
- `<nav>` — navigation links
- `<main>` — the primary content of the page
- `<section>` — a thematic grouping of content
- `<article>` — a self-contained piece of content
- `<aside>` — supporting content
- `<footer>` — closing or supporting info
- `<button>` — an action
- `<a>` — navigation to another location

```html
<header>
  <nav>
    <a href="/">Home</a>
    <a href="/work">Work</a>
  </nav>
</header>
<main>
  <article>
    <h1>Designing for the web</h1>
    <p>Great interfaces start with meaningful structure.</p>
  </article>
</main>
```

## Landmarks help users navigate

Screen-reader users can jump between major page landmarks. That only works if you mark them up properly.

A typical page might look like:

- `<header>`
- `<nav>`
- `<main>`
- `<aside>`
- `<footer>`

You do not need to wrap everything in generic `<div>` elements when a more meaningful option exists.

:::quiz
Q: Which element should wrap the page's primary content?
- `<section>`
- `<main>` *
- `<footer>`
E: `<main>` identifies the main content of the document. Screen readers and other tools treat it as a major landmark.
:::

## Headings create a real outline

Headings are not just visual size. They create structure.

Good practice:

- use one clear `<h1>` for the page's main topic
- use `<h2>` for major subsections
- use `<h3>` inside those subsections when needed
- don't skip heading levels just to get a visual effect

```html
<h1>Frontend Portfolio</h1>
<section>
  <h2>Featured Projects</h2>
  <article>
    <h3>Accessible signup form</h3>
  </article>
</section>
```

If you want smaller text, style it with CSS — don't fake structure with the wrong heading level.

## Buttons are not links, and links are not buttons

This mistake shows up everywhere.

Use:

- `<a>` when the user is **going somewhere**
- `<button>` when the user is **doing something**

Examples:

```html
<a href="/pricing">See pricing</a>
<button type="button">Open filters</button>
```

Wrong idea:

```html
<a href="#" onclick="openModal()">Open modal</a>
```

That should be a button, because it triggers an in-page action.

:::quiz
Q: A user clicks something to open a mobile menu without leaving the page. What should it be?
- A link
- A button *
- A heading
E: Opening a menu is an in-page action, so a `<button>` is the correct semantic element.
:::

## Accessible names matter

Interactive elements need labels users can understand.

Bad:

```html
<button>+</button>
```

Better:

```html
<button aria-label="Add task">+</button>
```

Form inputs also need visible labels:

```html
<label for="email">Email address</label>
<input id="email" type="email" />
```

:::fill
Q: Complete the label so it correctly connects to the input.
`<label ___="email">Email address</label>`
- for *
- id
- name
E: The label's `for` value should match the input's `id`.
:::

## Alt text: describe meaning, not pixels

When writing `alt` text, think about the image's job in the page.

Meaningful image:

```html
<img src="team.jpg" alt="Four developers presenting their launch demo on stage" />
```

Decorative image:

```html
<img src="divider.svg" alt="" />
```

Good `alt` text is:

- concise
- relevant
- contextual
- not redundant with nearby text

## Basic keyboard and focus rules

Not every user uses a mouse.

Your interface should support:

- tabbing to links, buttons, and form controls
- visible focus states
- logical document order
- real controls instead of clickable `<div>`s

If you find yourself writing `onclick` on a non-interactive element, stop and ask whether you picked the wrong element.

:::match
Q: Match the element to the job it does best.
- `<nav>` | A landmark for site navigation
- `<button>` | Triggers an in-page action
- `<a>` | Moves the user to a different location
- `<aside>` | Supporting or secondary content
E: Semantic choices make interfaces clearer for both humans and browsers.
:::

## Quick audit checklist

When reviewing your own HTML, ask:

- Does the page have one clear `<h1>`?
- Did I use landmarks like `<main>` and `<nav>`?
- Are buttons used for actions and links for navigation?
- Do meaningful images have helpful `alt` text?
- Does every form field have a label?
- Can the page still make sense if I remove all CSS?

:::key
If a page still reads clearly with no styling, your HTML structure is probably doing its job.
:::

## Mini practice — fix the semantics

Suppose you start with this:

```html
<div class="top">
  <div class="big">My Blog</div>
  <div class="menu">
    <div>Home</div>
    <div>Posts</div>
  </div>
</div>
<div class="content">
  <div class="title">Why forms fail</div>
  <div>Small label mistakes create big usability problems.</div>
  <div onclick="savePost()">Save</div>
</div>
```

A semantic rewrite would move toward:

- `<header>` instead of a generic top wrapper
- `<h1>` for the site title
- `<nav>` with real links for navigation
- `<main>` and `<article>` for content
- `<h2>` or `<h1>` for the article title depending on page structure
- `<button>` for the Save action

That's the skill: **seeing roles, not just boxes**.

## What's next

In **HTML Forms & Validation**, you'll apply the same thinking to forms — one of the places where semantics, labels, and accessibility matter most.
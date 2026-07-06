# HTML Basics

If frontend development is building for the browser, **HTML is the contract** between your ideas and the browser. It defines what content exists, what each piece means, and how the rest of your tools can work with it.

By the end of this lesson, you should be able to turn a rough page brief into clean, meaningful markup — not just throw random tags at the screen.

## What HTML is actually for

HTML stands for **HyperText Markup Language**:

- **HyperText** — content that can connect to other content through links.
- **Markup** — labels wrapped around content to describe what it *is*.
- **Language** — a shared vocabulary and set of rules browsers understand.

HTML is not about appearance. It answers questions like:

- Is this a page heading or a paragraph?
- Is this text a navigation link or a button?
- Is this an image, a list, or a form field?

:::key
A strong frontend developer thinks in **content and meaning first**. Styling and interaction come later.
:::

## The minimum structure every page needs

A real page usually starts with this skeleton:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>My first page</title>
  </head>
  <body>
    <h1>Hello, web!</h1>
    <p>This page has real structure.</p>
  </body>
</html>
```

What each part does:

- `<!DOCTYPE html>` tells the browser to use modern HTML rules.
- `<html lang="en">` wraps the whole document and declares the page language.
- `<head>` holds metadata the user doesn't directly read in the page body.
- `<meta charset="UTF-8">` lets the browser display text reliably.
- `<meta name="viewport" ...>` makes mobile layout behave properly.
- `<title>` sets the browser-tab title.
- `<body>` contains everything visible on the page.

:::reorder
Q: Put the document skeleton in the correct order.
- <!DOCTYPE html>
- <html lang="en">
- <head>
- <body>
- </html>
E: The document type comes first, then the root `<html>` element. Inside it, the page splits into `<head>` metadata and the visible `<body>`.
:::

## Elements, tags, and nesting

The core building block in HTML is the **element**. Most elements have:

- an opening tag
- content
- a closing tag

```html
<p>Ship small, improve fast.</p>
```

Some elements are **empty** because there is no inner content to wrap:

```html
<img src="avatar.png" alt="Portrait of Maya" />
<br />
```

HTML is also **nested**. Elements live inside other elements, creating a tree:

```html
<article>
  <h2>Frontend notes</h2>
  <p>Use <strong>semantic</strong> HTML from day one.</p>
</article>
```

The rule is simple: **close tags in the reverse order you opened them**.

## Attributes: extra information for an element

Attributes configure an element.

```html
<a href="/projects">See my projects</a>
<img src="headshot.jpg" alt="Headshot of Sam" />
```

Common examples:

- `href` — where a link goes
- `src` — where an image comes from
- `alt` — what an image means if it can't be seen
- `class` — a reusable styling hook for CSS and JS
- `id` — a unique identifier

:::fill
Q: Complete the image element with the attribute that tells the browser where the file lives.
`<img ___="avatar.png" alt="Portrait of a frontend developer" />`
- src *
- href
- class
E: `src` points to the image resource. `alt` describes the image for users who can't see it.
:::

## The tags you'll use constantly

These cover a huge amount of everyday frontend work:

- `<h1>` to `<h6>` — headings
- `<p>` — paragraphs
- `<ul>`, `<ol>`, `<li>` — lists
- `<a>` — links
- `<img>` — images
- `<strong>`, `<em>` — emphasis
- `<div>` and `<span>` — generic containers when no more meaningful element fits

```html
<h1>Launch checklist</h1>
<p>Before shipping, check content, layout, and accessibility.</p>
<ul>
  <li>Write real button labels</li>
  <li>Add alt text to images</li>
  <li>Test the page on mobile</li>
</ul>
```

:::match
Q: Match each element to its main job.
- `<a>` | A link to another location
- `<img>` | Embedded image content
- `<ul>` | A bulleted list
- `<h1>` | The main heading of the page
E: When you choose tags by purpose, you make your page easier to style, search, and navigate.
:::

## Text structure matters more than people think

Beginners often focus on tags one by one. Professionals think about **reading order**.

Good HTML usually has:

- one clear page-level `<h1>`
- sections that move from general to specific
- related items grouped in lists
- links with descriptive text

Weak link text:

```html
<a href="/pricing">Click here</a>
```

Better:

```html
<a href="/pricing">View pricing plans</a>
```

The second version makes sense even if the user only hears the link out loud in a screen reader or sees it out of context.

## Images need meaning, not decoration-only markup

Images are not just visual decoration. In HTML, they are content that may need explanation.

```html
<img src="team-photo.jpg" alt="The studio team standing in front of their first product launch banner" />
```

A good `alt` attribute answers: *what would a user miss if this image disappeared?*

Use empty alt text only for purely decorative images:

```html
<img src="sparkles.svg" alt="" />
```

## HTML is a tree the browser can reason about

The browser doesn't see a page the way a human does. It builds a tree from your markup and uses that structure to:

- render content
- apply CSS selectors
- attach JavaScript events
- expose the page to assistive technology

That means sloppy HTML creates problems later in **all three frontend layers**.

:::quiz
Q: Which statement best describes HTML?
- It decides how content looks on screen
- It gives content structure and meaning *
- It makes a page interactive with click handlers
E: HTML defines structure and meaning. CSS handles presentation, and JavaScript handles behavior.
:::

## Practical drill — mark up a content brief

Imagine you're handed this content brief:

- Site title: *Northstar Coffee*
- Intro paragraph about the shop
- A bulleted list of three drinks
- One photo of the café interior
- A link to the menu page

Before you write code, decide the structure:

- the site title should be an `<h1>`
- the intro should be a `<p>`
- the drinks belong in a `<ul>` with three `<li>` items
- the photo should be an `<img>` with useful alt text
- the menu should be an `<a>`

A clean solution might start like this:

```html
<h1>Northstar Coffee</h1>
<p>A neighborhood coffee shop focused on bright espresso, seasonal beans, and calm work-friendly seating.</p>
<ul>
  <li>Flat white</li>
  <li>Honey latte</li>
  <li>Cold brew</li>
</ul>
<img src="cafe-interior.jpg" alt="Warm café interior with wooden tables and plants near the window" />
<a href="/menu">View the full menu</a>
```

That is the mindset you want: **content brief → structure → markup**.

## Common beginner mistakes to avoid

- Using headings just because they are big by default
- Using a `<div>` for everything
- Skipping `alt` text on meaningful images
- Writing vague links like "read more" with no context
- Forgetting the viewport meta tag
- Nesting tags incorrectly

:::warning
If you use HTML purely to chase visual results, your CSS becomes harder, your JavaScript becomes messier, and your accessibility suffers.
:::

## What good looks like

By now, you should be able to:

- explain what HTML is for
- write a valid page skeleton
- choose common text, media, and list elements correctly
- add attributes like `href`, `src`, and `alt`
- turn a page brief into meaningful markup

## What's next

In **Semantic HTML & Accessibility**, you'll level this up: landmarks, heading hierarchy, buttons vs links, and the decisions that separate beginner markup from professional markup.
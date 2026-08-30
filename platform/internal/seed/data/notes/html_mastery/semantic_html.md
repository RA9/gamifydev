# Semantic HTML

Semantic HTML means choosing tags based on what the content *is*, not what it looks like. A `<nav>` says "this is navigation." A `<div>` says nothing. That difference ripples through accessibility, SEO, maintainability, and even your team's ability to read the code six months from now. In this lesson you'll learn every major semantic element and exactly when to reach for each one.

## Why Semantics Matter

Three audiences read your HTML — and none of them care about your CSS class names:

1. **Screen readers** build a navigable outline from your tags. A `<nav>` becomes a landmark a blind user can jump to. A `<div class="nav">` is invisible to that system.
2. **Search engines** weight content by tag. An `<h1>` carries more significance than a `<span class="big-text">`. An `<article>` helps Google identify standalone content.
3. **Developers** scan semantic tags faster than deciphering class names. `<aside>` instantly communicates intent; `<div class="sidebar-wrapper-container">` does not.

:::key
Semantic HTML is not extra work — it's the *right* work. Choosing the correct tag is usually the same effort as choosing a `<div>`, but the benefits in accessibility, SEO, and readability are enormous.
:::

## The Core Landmark Elements

HTML5 introduced landmark elements that describe the *role* of each part of a page. Screen readers expose these as jump-to landmarks, giving users a site-wide map.

### `<header>`

The introductory content for a page or a section — typically a logo, title, and navigation.

```html
<header>
  <a href="/" class="logo">Acme Co</a>
  <nav aria-label="Main navigation">…</nav>
</header>
```

A page can have multiple `<header>` elements — one for the site and one inside an `<article>`, for example.

### `<nav>`

A major block of navigation links. Use it for the site's main menu, a table of contents, or breadcrumbs — not for every group of links.

```html
<nav aria-label="Main navigation">
  <ul>
    <li><a href="/">Home</a></li>
    <li><a href="/docs">Docs</a></li>
  </ul>
</nav>
```

If your page has multiple `<nav>` elements, give each a distinct `aria-label` so screen readers can differentiate them.

### `<main>`

The central content of the page — the content unique to *this* page, excluding headers, footers, sidebars, and nav. There must be **exactly one** per page.

```html
<main>
  <h1>Getting Started with HTML</h1>
  <p>Welcome to the course…</p>
</main>
```

:::warning
Only one `<main>` per page. If you have multiple, screen readers won't know which is the primary content. Duplicating `<main>` is invalid HTML.
:::

### `<footer>`

Closing content for a page or section — copyright, related links, contact info.

```html
<footer>
  <p>&copy; 2025 Acme Co. All rights reserved.</p>
  <nav aria-label="Footer navigation">…</nav>
</footer>
```

Like `<header>`, you can have multiple `<footer>` elements — one for the page and one inside an `<article>`.

### `<article>`

A **self-contained, independently distributable** piece of content. If the content would make sense syndicated in an RSS feed, shared on social media, or embedded on another site, it's an `<article>`.

```html
<article>
  <header>
    <h2>Understanding Flexbox</h2>
    <time datetime="2025-03-15">March 15, 2025</time>
  </header>
  <p>Flexbox is a one-dimensional layout method…</p>
  <footer>
    <p>Written by Jane Doe</p>
  </footer>
</article>
```

Common uses: blog posts, news stories, forum threads, product cards, comments.

### `<section>`

A thematic grouping of content — a chapter, a tab panel, or a distinct segment of a page. It should almost always have a heading.

```html
<section>
  <h2>Pricing</h2>
  <p>Choose the plan that fits your needs.</p>
  …
</section>
```

### `<aside>`

Content tangentially related to the surrounding content — sidebars, pull quotes, ads, related links. It should make sense *removed* from the page without losing the main narrative.

```html
<aside>
  <h3>Related Articles</h3>
  <ul>
    <li><a href="/css-grid">CSS Grid Guide</a></li>
    <li><a href="/flexbox">Flexbox Guide</a></li>
  </ul>
</aside>
```

## The Decision Tree: `<div>` vs. Semantic Elements

When deciding which element to use, walk through these questions:

1. **Is it navigation?** → `<nav>`
2. **Is it the page's primary content?** → `<main>`
3. **Is it self-contained content that could stand alone?** → `<article>`
4. **Is it a thematic grouping with a heading?** → `<section>`
5. **Is it tangential or supplementary content?** → `<aside>`
6. **Is it introductory or closing content?** → `<header>` or `<footer>`
7. **None of the above, but you need a container for styling?** → `<div>`

:::key
`<div>` is the element of last resort. It's a generic container with no semantic meaning. Use it only when no semantic element fits — typically for CSS layout wrappers or JavaScript hooks.
:::

## A Complete Page Structure

Here's how these elements fit together in a typical page:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Blog — Acme Co</title>
</head>
<body>
  <header>
    <a href="/">Acme Co</a>
    <nav aria-label="Main navigation">
      <ul>
        <li><a href="/">Home</a></li>
        <li><a href="/blog" aria-current="page">Blog</a></li>
        <li><a href="/about">About</a></li>
      </ul>
    </nav>
  </header>

  <main>
    <h1>Blog</h1>

    <article>
      <header>
        <h2>Semantic HTML Matters</h2>
        <time datetime="2025-06-01">June 1, 2025</time>
      </header>
      <p>Choosing the right element improves accessibility…</p>
    </article>

    <article>
      <header>
        <h2>CSS Grid Deep Dive</h2>
        <time datetime="2025-05-20">May 20, 2025</time>
      </header>
      <p>Grid layout gives you two-dimensional control…</p>
    </article>
  </main>

  <aside>
    <h2>Popular Tags</h2>
    <ul>
      <li><a href="/tags/html">HTML</a></li>
      <li><a href="/tags/css">CSS</a></li>
    </ul>
  </aside>

  <footer>
    <p>&copy; 2025 Acme Co</p>
  </footer>
</body>
</html>
```

## Heading Hierarchy

Headings create the outline that screen readers and search engines use to understand your page. The rules are simple but critical:

1. **One `<h1>` per page** — the page's main topic.
2. **Never skip levels** — go `h1` → `h2` → `h3`, not `h1` → `h3`.
3. **Nest logically** — an `h3` is a subtopic of the nearest `h2` above it.

```html
<!-- Good: clean hierarchy -->
<h1>Web Development Course</h1>
  <h2>Module 1: HTML</h2>
    <h3>Lesson 1: Document Structure</h3>
    <h3>Lesson 2: Text Elements</h3>
  <h2>Module 2: CSS</h2>
    <h3>Lesson 1: Selectors</h3>

<!-- Bad: skipped levels, multiple h1s -->
<h1>Web Development</h1>
<h1>HTML Module</h1>
<h4>Document Structure</h4>
```

:::tip
Install the HeadingsMap browser extension or use the W3C validator to visualize your heading outline. It's the fastest way to catch skipped levels.
:::

## Landmark Roles and ARIA

Semantic elements automatically have associated ARIA landmark roles:

| Element | Implicit Role |
|---------|--------------|
| `<header>` (page-level) | `banner` |
| `<nav>` | `navigation` |
| `<main>` | `main` |
| `<aside>` | `complementary` |
| `<footer>` (page-level) | `contentinfo` |
| `<section>` (with label) | `region` |
| `<article>` | `article` |

You do **not** need to add `role="navigation"` to a `<nav>` — the role is already built in. Adding it is redundant.

:::warning
Don't add ARIA roles to elements that already have them. `<nav role="navigation">` is harmless but noisy. `<main role="main">` is doubly redundant. Trust the semantic elements to do their job.
:::

## Common Mistakes

1. **Using `<section>` as a generic container.** If it doesn't have a heading and isn't a thematic group, it's probably a `<div>`.
2. **Using `<article>` for everything.** An article is self-contained content. A sidebar widget is an `<aside>`, not an article.
3. **Wrapping everything in `<div>`s and adding classes.** Before writing `<div class="footer">`, ask if `<footer>` does the job.
4. **Skipping `<main>`.** Every page should have one `<main>` — it's the skip-link target and the primary landmark.

:::quiz
Q: Which element should wrap the primary, unique content of a page?
- `<section>`
- `<article>`
- `<main>` *
E: `<main>` marks the dominant content of the page — the part unique to this page, excluding site-wide navigation, headers, and footers. There should be exactly one per page.
:::

:::quiz
Q: When is a `<div>` the right choice?
- When you need a thematic grouping with a heading
- When you need a generic container for styling or scripting with no semantic meaning *
- When you need to mark self-contained content
E: `<div>` is a semantically neutral container. Use it when no semantic element fits — typically for CSS layout wrappers or JavaScript hooks. If the content has meaning, there's usually a better element.
:::

## Recap

- Semantic HTML communicates meaning to screen readers, search engines, and developers.
- Use `<header>`, `<nav>`, `<main>`, `<footer>`, `<article>`, `<section>`, and `<aside>` for their intended purposes.
- Only one `<main>` per page. `<article>` is for self-contained content. `<section>` is for thematic groups with headings.
- `<div>` is the fallback when no semantic element fits — never the default choice.
- Follow a strict heading hierarchy: one `<h1>`, never skip levels.
- Semantic elements carry implicit ARIA roles — don't duplicate them.

**Next up:** ARIA and Screen Readers — adding accessibility information when HTML alone isn't enough.

# HTML Best Practices and Review

You've learned document structure, semantics, accessibility, forms, and media. This final lesson ties it all together with the professional practices that separate polished production HTML from rough drafts — validation, performance, SEO, social sharing, and a comprehensive review checklist.

## HTML Validation

The W3C Markup Validation Service checks your HTML against the specification and catches errors you'd never spot by eye — unclosed tags, invalid attributes, duplicate IDs, and deprecated elements.

**How to use it:**

1. Go to [validator.w3.org](https://validator.w3.org).
2. Paste your URL, upload a file, or paste your HTML directly.
3. Fix the errors. Warnings are worth reading too.

```html
<!-- The validator catches this: -->
<p>Hello <strong>world</p></strong>
<!-- Error: end tag "p" seen before end tag "strong" -->

<!-- Corrected: -->
<p>Hello <strong>world</strong></p>
```

:::tip
Run the validator early and often — not just before launch. Catching a mismatched tag in development takes seconds; finding it after deployment in a specific browser on a specific device takes hours.
:::

Common errors the validator catches:

- Unclosed or misnested tags
- Duplicate `id` attributes
- Missing required attributes (`alt` on `<img>`, `charset` in `<meta>`)
- Invalid or deprecated elements and attributes
- Stray elements in the wrong context (`<div>` inside `<p>`)

## Performance Attributes

HTML gives you several attributes to control how resources load, directly impacting page speed.

### `defer` and `async` on Scripts

```html
<!-- Defer: download in parallel, execute after HTML is parsed, in order -->
<script src="app.js" defer></script>

<!-- Async: download in parallel, execute immediately when ready, any order -->
<script src="analytics.js" async></script>

<!-- Neither: blocks HTML parsing until downloaded and executed -->
<script src="blocking.js"></script>
```

| Attribute | Downloads in parallel? | Executes when? | Order guaranteed? |
|-----------|----------------------|----------------|-------------------|
| (none) | No — blocks parsing | Immediately | Yes |
| `defer` | Yes | After HTML parsed | Yes |
| `async` | Yes | When download finishes | No |

:::key
Use `defer` for scripts that depend on the DOM or each other. Use `async` only for truly independent scripts (analytics, ads) that don't care about order. Never use neither — blocking scripts are the #1 cause of slow page loads.
:::

### `preload` — Fetch Critical Resources Early

Tell the browser to start downloading a resource immediately, before it discovers the need during parsing.

```html
<link rel="preload" href="/fonts/Inter.woff2" as="font"
      type="font/woff2" crossorigin>
<link rel="preload" href="/css/critical.css" as="style">
<link rel="preload" href="/images/hero.webp" as="image">
```

Use `preload` for resources the browser would discover late — fonts referenced in CSS, hero images loaded via background-image, or critical scripts.

### `preconnect` — Warm Up Third-Party Connections

When you know you'll fetch from a third-party domain, `preconnect` resolves DNS, establishes TCP, and negotiates TLS *ahead of time*.

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://cdn.example.com" crossorigin>
```

This saves 100–300ms per connection. Use it for your CDN, font providers, and API endpoints.

### `fetchpriority`

Hint to the browser about how important a resource is relative to others of the same type.

```html
<!-- Hero image: load this first -->
<img src="hero.jpg" alt="…" fetchpriority="high">

<!-- Below-the-fold image: lower priority -->
<img src="footer-graphic.jpg" alt="…" fetchpriority="low" loading="lazy">
```

### `loading="lazy"` — Defer Off-Screen Resources

You learned this in the media lesson. It's worth repeating: lazy-load all below-the-fold images and iframes.

```html
<img src="photo.jpg" alt="…" width="800" height="600" loading="lazy">
<iframe src="…" title="…" loading="lazy"></iframe>
```

## SEO Meta Tags

Search engines read specific meta tags to understand, index, and display your pages.

### Page Description

```html
<meta name="description"
      content="Learn HTML from scratch with practical, project-based lessons.
               Master semantic markup, forms, accessibility, and more.">
```

This appears as the snippet below your title in search results. Keep it under 160 characters, make it compelling, and include relevant keywords naturally.

### Robots

Control how search engines index and follow your page.

```html
<!-- Default: index and follow all links (you don't need to specify this) -->
<meta name="robots" content="index, follow">

<!-- Don't index this page -->
<meta name="robots" content="noindex">

<!-- Index but don't follow links -->
<meta name="robots" content="index, nofollow">
```

### Canonical URL

Tell search engines which URL is the "official" version when the same content is accessible at multiple URLs.

```html
<link rel="canonical" href="https://example.com/blog/semantic-html">
```

This prevents duplicate content penalties when your page is reachable at both `http://` and `https://`, with and without `www`, or with query parameters.

## Open Graph Tags

Open Graph (OG) tags control how your page looks when shared on social media — Facebook, LinkedIn, Slack, Discord, and more.

```html
<meta property="og:title" content="HTML Mastery Course">
<meta property="og:description" content="Master HTML from document structure to advanced accessibility.">
<meta property="og:image" content="https://example.com/images/og-card.png">
<meta property="og:url" content="https://example.com/courses/html-mastery">
<meta property="og:type" content="website">
```

### Twitter/X Cards

```html
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="HTML Mastery Course">
<meta name="twitter:description" content="Master HTML from document structure to advanced accessibility.">
<meta name="twitter:image" content="https://example.com/images/twitter-card.png">
```

:::tip
OG images should be at least 1200×630 pixels for sharp display on high-DPI screens. Test your social cards with Facebook's Sharing Debugger and Twitter's Card Validator before launching.
:::

## Favicons

Modern browsers expect multiple icon formats:

```html
<!-- Standard favicon -->
<link rel="icon" href="/favicon.ico" sizes="32x32">

<!-- SVG favicon (modern browsers) -->
<link rel="icon" href="/favicon.svg" type="image/svg+xml">

<!-- Apple touch icon -->
<link rel="apple-touch-icon" href="/apple-touch-icon.png">

<!-- Web app manifest for PWA icons -->
<link rel="manifest" href="/site.webmanifest">
```

At minimum, provide a 32×32 `.ico` or `.png` and a 180×180 `apple-touch-icon`.

## Structured Data (JSON-LD)

Structured data helps search engines understand your content and display **rich results** — star ratings, recipe cards, FAQ dropdowns, event listings.

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Course",
  "name": "HTML Mastery",
  "description": "Master HTML from document structure to advanced accessibility.",
  "provider": {
    "@type": "Organization",
    "name": "GamifyDev"
  }
}
</script>
```

Common schema types: `Article`, `Course`, `Product`, `FAQPage`, `BreadcrumbList`, `Organization`.

:::tip
Use Google's Rich Results Test to validate your structured data. Invalid JSON-LD won't break your page, but it won't generate rich results either.
:::

## The Production HTML Checklist

Before every launch, walk through this list:

### Document Structure
- [ ] `<!DOCTYPE html>` is the first line
- [ ] `<html lang="…">` has the correct language
- [ ] `<meta charset="UTF-8">` is first in `<head>`
- [ ] `<meta name="viewport">` is set for responsive design
- [ ] `<title>` is unique and descriptive per page
- [ ] Exactly one `<main>` element

### Semantics
- [ ] Heading hierarchy is logical (`h1` → `h2` → `h3`, no skips)
- [ ] Semantic elements used over `<div>` where appropriate
- [ ] `<nav>`, `<header>`, `<footer>`, `<article>`, `<section>` used correctly

### Accessibility
- [ ] Every `<img>` has meaningful `alt` (or `alt=""` for decorative)
- [ ] Every form field has a connected `<label>`
- [ ] Color contrast meets WCAG 4.5:1 minimum
- [ ] Focus styles are visible (`:focus-visible`)
- [ ] Skip link is present
- [ ] ARIA used correctly and only where HTML falls short

### Performance
- [ ] Scripts use `defer` or `async`
- [ ] Critical fonts and images use `preload`
- [ ] Third-party origins use `preconnect`
- [ ] Below-the-fold images use `loading="lazy"`
- [ ] Images have `width` and `height` to prevent CLS

### SEO & Social
- [ ] `<meta name="description">` is set
- [ ] `<link rel="canonical">` points to the right URL
- [ ] Open Graph tags are set with title, description, and image
- [ ] Favicon is present
- [ ] Structured data validates in Rich Results Test

### Validation
- [ ] HTML passes W3C validator with no errors
- [ ] No duplicate `id` attributes
- [ ] No deprecated elements or attributes

:::key
This checklist is your safety net. Print it, bookmark it, or build it into your CI pipeline. The difference between professional and amateur HTML is consistency, and checklists enforce consistency.
:::

:::quiz
Q: What is the difference between `defer` and `async` on a `<script>` tag?
- `defer` blocks parsing; `async` does not
- `defer` executes after HTML is parsed and preserves order; `async` executes as soon as it downloads in any order *
- They are identical but `defer` is newer
E: Both download in parallel without blocking parsing. `defer` waits until the document is fully parsed and executes scripts in order. `async` executes each script the moment it finishes downloading, regardless of order — which is unpredictable for scripts that depend on each other.
:::

:::quiz
Q: Why should you include a `<link rel="canonical">` tag?
- To speed up page loading
- To tell search engines the preferred URL for a page, preventing duplicate content penalties *
- To add a favicon to the page
E: When the same content is available at multiple URLs (http vs https, with vs without www, with query parameters), the canonical tag tells search engines which one to index, consolidating ranking signals.
:::

## Recap

- **Validate** your HTML with the W3C validator regularly — not just before launch.
- Use `defer` for dependent scripts, `async` for independent ones, and never omit both.
- `preload` critical resources and `preconnect` to third-party origins for faster loads.
- Set `<meta name="description">`, `<link rel="canonical">`, and Open Graph tags for SEO and social sharing.
- Provide favicons in multiple formats and structured data for rich search results.
- Use the production checklist before every deployment.

**Congratulations!** You've completed the HTML Mastery course. You now understand document structure, text semantics, links, media, lists, tables, forms, validation, semantic HTML, ARIA, keyboard accessibility, and professional best practices. Time to build something great.

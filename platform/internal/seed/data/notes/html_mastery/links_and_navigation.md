# Links and Navigation

Links are the threads that stitch the web together. Without the humble `<a>` tag, every page would be an island. In this lesson you'll master anchor tags, build accessible navigation systems, and learn the subtle details that separate beginner link markup from professional-grade code.

## The Anchor Tag

The `<a>` (anchor) element creates a hyperlink. The `href` attribute is where it points.

```html
<a href="https://developer.mozilla.org">MDN Web Docs</a>
```

The text between the tags — "MDN Web Docs" — is the **link text**. It's what users see, click, and what screen readers announce. Make it descriptive.

## Types of `href` Values

### Absolute URLs

A full web address including the protocol. Use for links to external sites.

```html
<a href="https://github.com">GitHub</a>
```

### Relative URLs

Paths relative to the current page. Use for internal site links.

```html
<!-- Same directory -->
<a href="about.html">About</a>

<!-- Subdirectory -->
<a href="blog/first-post.html">First Post</a>

<!-- Parent directory -->
<a href="../index.html">Home</a>

<!-- Site root -->
<a href="/contact">Contact</a>
```

:::tip
Root-relative paths (starting with `/`) are the most reliable for internal links. They work regardless of where the current page sits in the directory tree, so you won't break links when pages move.
:::

### Fragment Links (Anchors)

Jump to a specific section within the page by linking to an element's `id`.

```html
<!-- The link -->
<a href="#pricing">Jump to Pricing</a>

<!-- The target (anywhere on the page) -->
<section id="pricing">
  <h2>Pricing</h2>
  …
</section>
```

You can combine fragments with full paths to jump to a section on another page:

```html
<a href="/docs/faq.html#refunds">Refund Policy</a>
```

### Email and Phone Links

```html
<a href="mailto:hello@example.com">Email us</a>
<a href="tel:+15551234567">Call us: (555) 123-4567</a>
```

`mailto:` opens the user's email client. `tel:` triggers a phone call on mobile devices. Always include the country code in `tel:` links for international compatibility.

:::tip
You can pre-fill email fields: `mailto:hello@example.com?subject=Support&body=Hi%20there`. URL-encode spaces as `%20`.
:::

## Opening Links in New Tabs

Use `target="_blank"` to open a link in a new tab. But always pair it with `rel="noreferrer"` (or at minimum `rel="noopener"`).

```html
<a href="https://external-site.com"
   target="_blank"
   rel="noreferrer">
  External Resource
</a>
```

Why `rel="noreferrer"`? Without it, the new page can access your page via `window.opener` — a security risk called **reverse tabnabbing** where the external page could redirect your tab to a phishing site. `noreferrer` also implies `noopener` and prevents the referer header from being sent.

:::warning
Never use `target="_blank"` without `rel="noreferrer"` or `rel="noopener"`. It's a security vulnerability. Many linters and accessibility tools will flag this automatically.
:::

:::key
As a general rule, don't force links to open in new tabs. Let users decide. Reserve `target="_blank"` for cases where the user would lose important state — like a link in the middle of filling out a form.
:::

## Writing Good Link Text

Screen reader users often navigate by pulling up a list of all links on a page. Every link that says "click here" or "read more" becomes meaningless noise.

```html
<!-- Bad -->
<p>To see our plans, <a href="/pricing">click here</a>.</p>

<!-- Good -->
<p>Compare our <a href="/pricing">pricing plans</a>.</p>
```

Rules for link text:

1. **Describe the destination**, not the action of clicking.
2. **Make sense out of context** — the link text alone should tell you where it goes.
3. **Be specific** — "Download the 2025 annual report (PDF, 2.4 MB)" beats "Download."
4. **Don't use the URL as the text** unless the URL itself is the point.

## Skip Links

A **skip link** is a hidden link at the very top of the page that lets keyboard users jump straight to the main content, skipping the navigation.

```html
<body>
  <a href="#main-content" class="skip-link">Skip to main content</a>
  <header>…</header>
  <nav>…</nav>
  <main id="main-content">…</main>
</body>
```

```css
.skip-link {
  position: absolute;
  top: -100%;
  left: 0;
  z-index: 100;
  padding: 0.5rem 1rem;
  background: #1f2937;
  color: #ffffff;
}

.skip-link:focus {
  top: 0;
}
```

The link is off-screen by default and slides into view only when focused via keyboard. This way sighted mouse users never see it, but keyboard users can hit Tab once and skip a long nav menu.

:::key
Skip links are a WCAG 2.1 Level A requirement. Every production site with a navigation bar should have one. It's one of the easiest accessibility wins you can implement.
:::

## The `<nav>` Element

Wrap your primary navigation in `<nav>`. It tells assistive technology "this is a navigation landmark" so screen reader users can jump to it directly.

```html
<nav aria-label="Main navigation">
  <ul>
    <li><a href="/">Home</a></li>
    <li><a href="/products">Products</a></li>
    <li><a href="/about">About</a></li>
    <li><a href="/contact">Contact</a></li>
  </ul>
</nav>
```

If your page has multiple `<nav>` elements (main nav, footer nav, sidebar nav), give each a unique `aria-label` so screen readers can distinguish them.

```html
<nav aria-label="Main navigation">…</nav>
<nav aria-label="Footer navigation">…</nav>
```

:::tip
Navigation links should live in a `<ul>` inside `<nav>`. The list structure tells screen readers how many links there are ("list, 4 items"), which helps users gauge the size of the menu.
:::

## Marking the Current Page

Help users know where they are by marking the current page's link with `aria-current="page"`.

```html
<nav aria-label="Main navigation">
  <ul>
    <li><a href="/" aria-current="page">Home</a></li>
    <li><a href="/products">Products</a></li>
    <li><a href="/about">About</a></li>
  </ul>
</nav>
```

Screen readers will announce "Home, current page, link." You can also style it in CSS:

```css
a[aria-current="page"] {
  font-weight: bold;
  border-bottom: 2px solid currentColor;
}
```

## Breadcrumb Navigation

Breadcrumbs show the user's path through a site hierarchy.

```html
<nav aria-label="Breadcrumb">
  <ol>
    <li><a href="/">Home</a></li>
    <li><a href="/docs">Documentation</a></li>
    <li><a href="/docs/html" aria-current="page">HTML</a></li>
  </ol>
</nav>
```

Use `<ol>` (ordered list) because the order matters — it's a path. The last item gets `aria-current="page"`. Separators like `/` or `>` should be added via CSS `::before` pseudo-elements rather than as text, so screen readers don't read out "slash" between every item.

```css
nav[aria-label="Breadcrumb"] li + li::before {
  content: " / ";
  color: #6b7280;
}
```

## Links vs. Buttons

A common source of confusion: when do you use a link, and when a button?

- **Link (`<a>`)** — navigates the user to a new page or location. Has an `href`.
- **Button (`<button>`)** — performs an action: toggle a menu, submit a form, open a dialog.

```html
<!-- Correct: navigating to a page -->
<a href="/settings">Settings</a>

<!-- Correct: performing an action -->
<button onclick="toggleMenu()">Menu</button>

<!-- Wrong: a link that doesn't go anywhere -->
<a href="#" onclick="toggleMenu()">Menu</a>
```

:::warning
Don't use `<a href="#">` or `<a href="javascript:void(0)">` as buttons. They announce as links to screen readers, they appear in link lists, and the `#` jumps the page to the top. If it triggers an action, it's a `<button>`.
:::

:::quiz
Q: Why should you add `rel="noreferrer"` when using `target="_blank"`?
- To improve SEO rankings
- To prevent the new page from accessing your page via `window.opener` *
- To make the link open faster
E: Without `rel="noreferrer"` (or `rel="noopener"`), the linked page can access `window.opener` and potentially redirect your tab to a malicious site — a vulnerability called reverse tabnabbing.
:::

:::quiz
Q: What's the correct element for a "Delete account" action that doesn't navigate anywhere?
- `<a href="#">Delete account</a>`
- `<button>Delete account</button>` *
- `<a href="javascript:void(0)">Delete account</a>`
E: Actions that don't navigate the user to a new URL should use `<button>`. Links are for navigation. Using `<a>` with `href="#"` or `javascript:void(0)` is a semantic misuse that hurts accessibility.
:::

## Recap

- `<a href>` supports absolute URLs, relative paths, fragment links (`#id`), `mailto:`, and `tel:`.
- Always pair `target="_blank"` with `rel="noreferrer"` to prevent reverse tabnabbing.
- Write descriptive link text that makes sense read out of context.
- Use skip links to let keyboard users bypass navigation — it's a WCAG requirement.
- Wrap navigation in `<nav>` with `aria-label`; use `<ul>` for nav lists and `aria-current="page"` for the active link.
- Use `<ol>` for breadcrumbs and add separators via CSS, not inline text.
- Links navigate; buttons act. Never use `<a href="#">` as a button.

**Next up:** Images and Media — making visuals fast, responsive, and accessible.

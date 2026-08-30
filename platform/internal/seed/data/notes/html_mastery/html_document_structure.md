# HTML Document Structure

Every web page starts with the same skeleton. Get this structure right and everything else — styling, scripting, accessibility — slots into place naturally. Get it wrong and you'll fight mysterious rendering bugs, broken SEO, and confused screen readers. In this lesson you'll build a complete starter template and understand every single line.

## The DOCTYPE

The very first line of any HTML file is the DOCTYPE declaration.

```html
<!DOCTYPE html>
```

This isn't an HTML tag — it's an instruction to the browser: "Treat this document as modern HTML5." Without it, the browser falls back to **quirks mode**, an old compatibility mode that interprets CSS and layout differently. You'll get subtle, maddening bugs.

:::warning
Never omit the DOCTYPE. Quirks mode renders box models, margins, and fonts differently from standards mode. A missing DOCTYPE is one of the hardest bugs to diagnose because the page *almost* works.
:::

## The `<html>` Element

The root element wraps everything on the page. Its most important attribute is `lang`.

```html
<html lang="en">
```

The `lang` attribute tells browsers and screen readers which language the content is in. A screen reader uses it to pick the correct pronunciation engine — without it, a French sentence might be read with English phonetics. Search engines also use it for language-specific results.

:::tip
If your page contains a section in a different language, you can add `lang` to that specific element: `<blockquote lang="fr">`. The attribute is inherited, so you only need it on `<html>` for the primary language.
:::

## The `<head>` Section

The `<head>` contains metadata — information *about* the page that the user doesn't see directly. Here's what belongs inside.

### Character Encoding

```html
<meta charset="UTF-8">
```

This **must** be the first element inside `<head>` (within the first 1024 bytes of the file). UTF-8 supports every character from every language, plus emoji. If you omit it, some browsers will guess the encoding and your special characters may appear garbled.

### Viewport Meta

```html
<meta name="viewport" content="width=device-width, initial-scale=1.0">
```

This tells mobile browsers to set the viewport width to the device's screen width rather than pretending to be a 980px-wide desktop monitor. Without it, your page will look tiny and zoomed-out on phones.

:::key
The viewport meta tag is what makes responsive design possible. If your page looks "zoomed out" on mobile, this tag is almost certainly missing.
:::

### Title

```html
<title>Home — My Website</title>
```

The `<title>` appears in browser tabs, bookmarks, search results, and screen reader announcements. Keep it concise, unique per page, and front-load the meaningful words. "Home — My Website" is better than "My Website — Home" because the important word appears first in a narrow tab.

### Linking Stylesheets

```html
<link rel="stylesheet" href="/css/main.css">
```

Stylesheets go in `<head>` because the browser needs them before it can render the page. If you place a stylesheet in `<body>`, the page may flash unstyled content while the CSS loads.

### Favicon

```html
<link rel="icon" href="/favicon.ico" type="image/x-icon">
```

The tiny icon next to your page's title in the browser tab.

## The `<body>` Section

Everything the user *sees* lives inside `<body>`. A semantic body typically contains landmark elements:

```html
<body>
  <header>…</header>
  <nav>…</nav>
  <main>…</main>
  <footer>…</footer>
</body>
```

We'll cover these landmarks in depth in the Semantic HTML lesson.

## Script Placement

Where you put `<script>` tags matters for performance and correctness.

```html
<!-- Option A: end of body (classic approach) -->
<body>
  <!-- page content -->
  <script src="/js/app.js"></script>
</body>

<!-- Option B: in head with defer (modern approach) -->
<head>
  <script src="/js/app.js" defer></script>
</head>
```

**Option B is preferred today.** The `defer` attribute tells the browser to download the script in parallel while parsing the HTML, then execute it after the document is fully parsed. You get the fastest load time without blocking rendering.

:::tip
`defer` scripts execute in the order they appear, which is predictable. `async` scripts execute as soon as they finish downloading, in any order — use `async` only for independent scripts like analytics that don't depend on each other or on DOM content.
:::

## The Complete Starter Template

Here is a production-ready starter template. Every line earns its place.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <!-- Character encoding — must come first -->
  <meta charset="UTF-8">

  <!-- Responsive viewport -->
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  <!-- Page title (shows in tabs, bookmarks, search) -->
  <title>Page Title — Site Name</title>

  <!-- Page description for search engines -->
  <meta name="description" content="A brief, compelling summary of this page.">

  <!-- Favicon -->
  <link rel="icon" href="/favicon.ico" type="image/x-icon">

  <!-- Stylesheet -->
  <link rel="stylesheet" href="/css/main.css">

  <!-- JavaScript with defer so it runs after parsing -->
  <script src="/js/app.js" defer></script>
</head>
<body>
  <header>
    <a href="/">Site Name</a>
    <nav aria-label="Main navigation">
      <ul>
        <li><a href="/">Home</a></li>
        <li><a href="/about">About</a></li>
        <li><a href="/contact">Contact</a></li>
      </ul>
    </nav>
  </header>

  <main>
    <h1>Page Title</h1>
    <p>Your content starts here.</p>
  </main>

  <footer>
    <p>&copy; 2025 Site Name</p>
  </footer>
</body>
</html>
```

:::key
This template is your starting point for every project. Memorize the order: DOCTYPE → html with lang → head (charset first, then viewport, title, meta, links, scripts with defer) → body with semantic landmarks.
:::

## Common Mistakes

1. **Putting `<meta charset>` after `<title>`** — the browser must know the encoding before it reads any text. Always make charset the first thing in `<head>`.
2. **Multiple `<main>` elements** — a page should have exactly one `<main>`. It marks the primary content, distinct from headers, footers, and sidebars.
3. **Forgetting `lang` on `<html>`** — this silently breaks screen reader pronunciation and hurts internationalization.
4. **Loading scripts without `defer` in `<head>`** — the browser stops parsing HTML to download and execute the script, delaying the page render.

:::quiz
Q: Why must `<meta charset="UTF-8">` appear as the first element inside `<head>`?
- It makes the page load faster
- The browser needs to know the encoding before reading any text content *
- It prevents CSS from loading incorrectly
E: The browser must determine the character encoding before it can interpret any text in the document, including the `<title>`. Placing charset first ensures nothing is misread.
:::

:::quiz
Q: What does the `defer` attribute on a `<script>` tag do?
- Prevents the script from running entirely
- Downloads the script in parallel and executes it after the HTML is fully parsed *
- Makes the script load synchronously
E: `defer` lets the browser download the script without pausing HTML parsing, then runs the script in order once the document is ready — giving you the best of both worlds.
:::

## Recap

- Start every page with `<!DOCTYPE html>` to trigger standards mode.
- Set `lang` on `<html>` for screen readers and search engines.
- In `<head>`: charset first, then viewport, title, meta description, favicon, stylesheets, and deferred scripts.
- The viewport meta tag is essential for responsive design on mobile.
- Place scripts in `<head>` with `defer`, or at the end of `<body>`.
- Use exactly one `<main>` element and build the body with semantic landmarks.

**Next up:** Text Elements and Typography — choosing the right tags for every kind of text.

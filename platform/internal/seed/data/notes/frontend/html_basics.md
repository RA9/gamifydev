# HTML Basics

HTML is the language you use to give a web page its structure and meaning. In this lesson you'll learn how pages are built out of tags, what goes in the document skeleton, and how to write clean, meaningful markup that browsers and people both understand.

## What HTML Is

HTML stands for **HyperText Markup Language**. It's not a programming language (there are no variables or loops here) — it's a *markup* language. That means you take plain text and wrap it in labels that describe what each piece *is*: this is a heading, this is a paragraph, this is a link.

The browser reads your HTML and turns it into the page you see. Your job is to describe the content; the browser handles drawing it.

:::analogy
Think of HTML like the labeled boxes when you move house. You write "KITCHEN" or "BOOKS" on each box. The label doesn't change what's inside, but it tells everyone how to handle it. HTML tags label your content so the browser knows a heading from a paragraph from a link.
:::

## Tags, Elements, and Attributes

Three words you'll hear constantly. Let's pin them down.

A **tag** is the keyword inside angle brackets. Most come in pairs: an opening tag and a closing tag (the closing one has a slash).

```html
<p>This is a paragraph.</p>
```

Here `<p>` is the opening tag and `</p>` is the closing tag.

An **element** is the whole thing together: the opening tag, the content, and the closing tag. The example above is one paragraph element.

An **attribute** adds extra information to an element. Attributes live inside the opening tag as `name="value"` pairs.

```html
<a href="https://example.com">Visit Example</a>
```

Here `href` is the attribute name and `"https://example.com"` is its value.

:::tip
Always wrap attribute values in double quotes. `href="page.html"` is safe; `href=page.html` can break the moment a value contains a space.
:::

Some elements have no content and no closing tag — they're called **void** or self-closing elements, like the line break and the image:

```html
<br>
<img src="cat.jpg" alt="A sleeping cat">
```

## The Document Skeleton

Every HTML page starts from the same bones. Memorize this shape — you'll type it constantly.

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My First Page</title>
  </head>
  <body>
    <h1>Hello, world!</h1>
    <p>Welcome to my page.</p>
  </body>
</html>
```

Let's walk through each part:

- `<!DOCTYPE html>` tells the browser "use modern HTML." It must be the very first line.
- `<html lang="en">` wraps the whole document. The `lang` attribute tells screen readers and search engines the page is in English.
- `<head>` holds information *about* the page that visitors don't see directly.
- `<meta charset="UTF-8">` sets the character encoding so accented letters and emoji display correctly.
- `<meta name="viewport" ...>` makes the page scale properly on phones. Without it, your site looks tiny and zoomed-out on mobile.
- `<title>` is the text shown in the browser tab and in search results.
- `<body>` holds everything the visitor actually sees.

:::warning
Leaving out the viewport meta tag is one of the most common beginner mistakes. Your page might look fine on a laptop and unusable on a phone. Always include it.
:::

## Headings: h1 to h6

Headings create a hierarchy, like a table of contents. There are six levels.

```html
<h1>Chocolate Chip Cookies</h1>
<h2>Ingredients</h2>
<h3>For the dough</h3>
<h2>Steps</h2>
```

`<h1>` is the most important (usually the page's main title) and `<h6>` the least. Use them in order based on *meaning*, not size.

:::key
Use only one `<h1>` per page, and never skip levels for looks. If you want smaller text, change it with CSS later — don't jump from `<h1>` to `<h4>` just because `<h4>` is smaller.
:::

## Paragraphs

The paragraph element holds blocks of text.

```html
<p>HTML is the skeleton of every web page.</p>
<p>Once you know the tags, you can build anything.</p>
```

Note that extra spaces and line breaks in your code are collapsed by the browser. To start a new paragraph, use a new `<p>`, not blank lines.

## Links

Links are what make the web a *web*. You create them with the anchor element `<a>` and its `href` attribute (the destination).

```html
<a href="https://developer.mozilla.org">Read the docs</a>
```

There are two kinds of destinations:

- **Absolute URL** — the full address, used for other websites: `href="https://example.com/about"`
- **Relative URL** — a path relative to the current page, used within your own site: `href="about.html"` or `href="pages/contact.html"`

You can also open a link in a new tab with `target`:

```html
<a href="https://example.com" target="_blank" rel="noopener">Open in new tab</a>
```

:::tip
When you use `target="_blank"`, add `rel="noopener"`. It's a small security measure that prevents the new page from gaining control over your original tab.
:::

## Lists

Two everyday list types: unordered (bullets) and ordered (numbers). Both contain list items, `<li>`.

```html
<ul>
  <li>Eggs</li>
  <li>Milk</li>
  <li>Flour</li>
</ul>

<ol>
  <li>Preheat the oven</li>
  <li>Mix the batter</li>
  <li>Bake for 20 minutes</li>
</ol>
```

Use `<ul>` when order doesn't matter and `<ol>` when it does (steps, rankings).

## Images, Conceptually

Images use the void element `<img>`. The `src` attribute points to the file; the `alt` attribute describes it.

```html
<img src="beach.jpg" alt="A sunny beach with palm trees">
```

The `alt` text is read aloud by screen readers and shown if the image fails to load. It is not optional fluff.

:::warning
Never skip `alt`. For meaningful images, describe what's shown. For purely decorative images, use an empty `alt=""` so screen readers skip it. You'll learn more about this in the Accessibility lesson.
:::

## Block vs Inline Elements

Elements come in two flavors:

- **Block** elements start on a new line and take up the full width available. Examples: `<p>`, `<h1>`, `<ul>`, `<div>`.
- **Inline** elements sit inside a line of text and only take the space they need. Examples: `<a>`, `<img>`, `<span>`, `<strong>`.

```html
<p>This is a paragraph with a <a href="#">link</a> inside it.</p>
```

The link flows inside the paragraph because it's inline. The paragraph itself stacks on its own line because it's block.

## Nesting and Indentation

Elements can live inside other elements — that's **nesting**. The rule: tags must close in the reverse order you opened them.

```html
<p>This text is <strong>very <em>important</em></strong>.</p>
```

Notice `<em>` closes before `<strong>` — they don't cross over. Indenting nested elements with spaces makes this structure easy to read.

:::quiz
Q: What is wrong with `<p>Hello <strong>friend</p></strong>`?
- The `<p>` tag is spelled wrong
- `<strong>` is not a real tag
- The tags overlap instead of nesting — `</strong>` should come before `</p>` *
E: Tags must close in reverse order, so the inner `<strong>` must close before the outer `<p>`.
:::

## Comments

Comments are notes for yourself that the browser ignores.

```html
<!-- This is the main navigation -->
<nav>...</nav>
```

Use them to explain tricky bits, but don't leave secrets in there — anyone can view your page source.

## Semantic Elements (Beat the Div Soup)

A `<div>` is a generic box with no meaning. If you build an entire page out of nothing but `<div>`s, you get what developers call **div soup** — markup that works visually but tells no one what anything *is*.

Semantic elements fix this by naming the parts of a page:

```html
<header>Site logo and intro</header>
<nav>Links to other pages</nav>
<main>
  <article>A self-contained blog post</article>
  <aside>Related links sidebar</aside>
</main>
<footer>Copyright and contact</footer>
```

- `<header>` — top section, often logo and intro
- `<nav>` — navigation links
- `<main>` — the primary content (only one per page)
- `<section>` — a thematic grouping of content
- `<article>` — a standalone piece that makes sense on its own
- `<aside>` — side content like a sidebar
- `<footer>` — bottom section

:::key
Semantic HTML helps screen readers, search engines, and other developers understand your page. A `<nav>` announces "navigation" to assistive tech; a `<div class="nav">` says nothing. Reach for the meaningful tag first, and use `<div>` only when no semantic element fits.
:::

:::quiz
Q: Which element should hold the main, unique content of a page?
- `<div>`
- `<header>`
- `<main>` *
E: `<main>` is the semantic element for the page's primary content, and there should be only one per page.
:::

## A Full Example Page

Here's everything coming together into one small, valid page.

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ada's Recipes</title>
  </head>
  <body>
    <header>
      <h1>Ada's Recipes</h1>
      <nav>
        <a href="index.html">Home</a>
        <a href="about.html">About</a>
      </nav>
    </header>

    <main>
      <article>
        <h2>Chocolate Chip Cookies</h2>
        <p>The best cookies you'll ever bake.</p>
        <img src="cookies.jpg" alt="A plate of golden cookies">
        <h3>Ingredients</h3>
        <ul>
          <li>Flour</li>
          <li>Butter</li>
          <li>Chocolate chips</li>
        </ul>
        <h3>Steps</h3>
        <ol>
          <li>Mix the dough</li>
          <li>Bake for 12 minutes</li>
        </ol>
      </article>
    </main>

    <footer>
      <p>Made with love by Ada.</p>
    </footer>
  </body>
</html>
```

:::example
Copy this into a file called `index.html`, double-click it, and your browser will render a real web page. Change the text, save, and refresh to see your edits live. That fast feedback loop is the heart of frontend development.
:::

## Recap

- HTML uses **tags** to wrap content into **elements**; **attributes** add extra info inside opening tags.
- Every page needs the skeleton: `<!DOCTYPE html>`, `<html>`, `<head>` (with `charset` and `viewport`), `<title>`, and `<body>`.
- Headings `<h1>`–`<h6>` build a hierarchy; use one `<h1>` and don't skip levels.
- Links use `<a href>` (absolute vs relative); lists use `<ul>`/`<ol>`/`<li>`; images need a meaningful `alt`.
- **Block** elements stack; **inline** elements flow within text. Nest tags so they close in reverse order.
- **Semantic** elements (`header`, `nav`, `main`, `section`, `article`, `aside`, `footer`) give meaning and beat div soup.

**Next up:** HTML Forms and Inputs — collecting information from your visitors.

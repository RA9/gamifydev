# Building a Website with HTML and CSS

You've learned the pieces — HTML for structure, CSS for style. Now you'll put them together to build a real, multi-section web page. This is where it finally *feels* like making a website.

By the end of this lesson you'll know how to organise a project's files, structure a page with semantic HTML, and connect a stylesheet that brings it to life.

## From snippets to a real page

A real page isn't one giant blob of HTML — it's a handful of well-named **regions** stacked together: a header at the top, the main content in the middle, a footer at the bottom.

![A typical semantic page layout](/images/lessons/page-layout.svg)

Using elements that *describe their role* — `<header>`, `<main>`, `<section>`, `<aside>`, `<footer>` — is called **semantic HTML**. It makes your page easier to style, better for search engines, and accessible to screen readers out of the box.

:::analogy
Building a page is like organising a house into rooms. You wouldn't pile the kitchen, bedroom, and bathroom into one space. `<header>`, `<main>`, and `<footer>` are your rooms — each with a clear purpose.
:::

## Organising your project files

Even a tiny site has a tidy structure. A common starting point:

```text
my-site/
├── index.html      ← the page
├── styles.css      ← all your CSS
└── images/         ← logos, photos
    └── logo.png
```

Keeping HTML, CSS, and images in their own places keeps the project understandable as it grows.

## Connecting HTML and CSS

The two files meet through one line in the `<head>`:

```html
<head>
  <link rel="stylesheet" href="styles.css" />
</head>
```

That `<link>` tells the browser "also load this stylesheet." Now any rule in `styles.css` applies to your page.

## A complete little page

Here's a real page using semantic structure:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <link rel="stylesheet" href="styles.css" />
  </head>
  <body>
    <header>
      <h1>Sunrise Café</h1>
    </header>
    <main>
      <section>
        <h2>Today's special</h2>
        <p>Freshly baked cinnamon rolls.</p>
      </section>
    </main>
    <footer>
      <p>Open daily, 7am–3pm.</p>
    </footer>
  </body>
</html>
```

And a touch of CSS to style it:

```css
header {
  background: #6740e8;
  color: white;
  padding: 24px;
}

main {
  padding: 24px;
}
```

Notice how the **HTML never mentions colours or spacing** — that's all in the CSS. Keeping structure and style separate means you can redesign the whole look without touching the content.

:::quiz
Q: Which element should wrap the primary content of the page?
- `<header>`
- `<main>` *
- `<footer>`
E: `<main>` holds the page's primary content. `<header>` and `<footer>` are for the top and bottom regions.
:::

:::tip
Style broad, then specific. Set sensible defaults on elements (like `body` font and colour), then use classes for the bits that need to stand out. You'll write far less CSS.
:::

## Layout with flexbox (a first taste)

To place items side by side — like cards in a row — modern CSS uses **flexbox**:

```css
.cards {
  display: flex;
  gap: 16px;
}
```

Add `display: flex` to a container and its children line up in a row with even spacing. It's the workhorse of modern layout; you'll use it constantly.

:::quiz
Q: What connects an external stylesheet to your HTML page?
- A `<style>` element in the body
- A `<link rel="stylesheet">` in the head *
- The `class` attribute
E: The `<link rel="stylesheet" href="...">` tag in the `<head>` loads an external CSS file.
:::

:::key
A real page is **semantic regions** (`<header>`, `<main>`, `<footer>`) structured in HTML and styled by a **separate** stylesheet linked in the `<head>`. Keeping structure and style apart is what makes sites maintainable.
:::

## Talk about it

Explain out loud:

> "Why do we keep HTML and CSS in separate files? And what does semantic HTML give us?"

If you can answer both, you understand the foundation real websites are built on.

## What's next

Your page looks great — but it's still static. Next, **Building Interactive JavaScript Websites** makes it respond to the people using it.

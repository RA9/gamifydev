# HTML Basics

Every website you've ever used — your bank, your favourite social app, this very page — starts life as an HTML document. Before you can style anything or make it interactive, you need a **structure**. That's what HTML gives you.

By the end of this lesson you won't just *know* a few tags — you'll understand what HTML is actually for, so you could explain it to a friend without looking anything up.

## What HTML really is

HTML stands for **HyperText Markup Language**. Let's unpack that name, because it tells you almost everything:

- **HyperText** — text that can link to other text. The "links" you click are the *hyper* part. This is the idea the whole web is built on.
- **Markup** — you take plain content and *mark it up* by wrapping pieces of it in labels that say what each piece *is*: "this is a heading", "this is a paragraph", "this is a link".
- **Language** — it has rules and a vocabulary, like any language.

So HTML is not a programming language — it doesn't make decisions or do maths. It's a **markup language**: a way of describing the *meaning and structure* of content.

:::analogy
HTML is the **skeleton** of a web page. CSS is the *skin and clothes* (how it looks), and JavaScript is the *muscles* (how it moves and reacts). You're learning the skeleton first — because nothing else can stand up without it.
:::

## The three languages of the web

It's worth fixing this mental model early, because every front-end task fits into one of these three buckets:

1. **HTML** — *structure*. What content exists and what it means.
2. **CSS** — *presentation*. Colours, spacing, layout, fonts.
3. **JavaScript** — *behaviour*. What happens when you click, type, or scroll.

When something looks wrong, ask: "is this a structure problem (HTML), a styling problem (CSS), or a behaviour problem (JavaScript)?" That single question will save you hours.

## Elements and tags

The building block of HTML is the **element**. You create an element with **tags** — labels wrapped in angle brackets.

![The anatomy of an HTML element](/images/lessons/html-element-anatomy.svg)

Read the diagram above slowly. An element usually has:

- an **opening tag**, like `<p>`
- some **content** in the middle
- a **closing tag**, like `</p>` — note the slash `/`

Put together, a paragraph element looks like this:

```html
<p>Learning to code is a superpower.</p>
```

A few elements are **empty** — they have no content and no closing tag, because there's nothing to wrap. An image and a line break are the classic examples:

```html
<img src="cat.png" alt="A sleepy cat" />
<br />
```

:::quiz
Q: In the element `<h1>Welcome</h1>`, which part is the **closing** tag?
- `<h1>`
- `</h1>` *
- `Welcome`
E: The closing tag repeats the tag name with a slash in front: `</h1>`. It tells the browser where the element ends.
:::

## A complete (tiny) page

Real pages share the same backbone. Here's the smallest sensible HTML document, and what each line is doing:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>My first page</title>
  </head>
  <body>
    <h1>Hello, world!</h1>
    <p>This is my very first web page.</p>
  </body>
</html>
```

- `<!DOCTYPE html>` — tells the browser "use modern HTML". Always the first line.
- `<html>` — the root; everything lives inside it. `lang="en"` says the page is in English (great for screen readers and search engines).
- `<head>` — information *about* the page that you don't see directly, like its `<title>` (the text on the browser tab).
- `<body>` — everything you actually see on screen.

:::analogy
Think of an HTML document like a letter. The `<head>` is the envelope — address and metadata. The `<body>` is the letter itself — the part the reader reads.
:::

## Attributes: giving elements extra information

Tags can carry **attributes** — name–value pairs that configure the element. You already saw two: `href` on a link and `src`/`alt` on an image.

```html
<a href="https://example.com">Visit Example</a>
```

Here `href` (the value `"https://example.com"`) tells the link *where to go*. Without it, a link doesn't link to anything.

:::tip
The `alt` attribute on an image describes it in words. If the image fails to load — or a blind user is browsing with a screen reader — the `alt` text is what they get. Always write meaningful `alt` text. It's one of the easiest ways to make the web work for everyone.
:::

## Nesting: HTML is a tree

Elements live *inside* other elements. That's called **nesting**, and it turns your page into a tree of parents and children:

```html
<body>
  <article>
    <h2>Why nesting matters</h2>
    <p>The <strong>important</strong> word is bold.</p>
  </article>
</body>
```

Here `<article>` contains an `<h2>` and a `<p>`; the `<p>` contains a `<strong>`. The one rule: **close your tags in the reverse order you opened them** — like nesting boxes. `<p><strong>...</strong></p>` is correct; `<p><strong>...</p></strong>` is not.

## A quick tour of everyday elements

You'll reach for these constantly:

- **Headings**: `<h1>` down to `<h6>` — `<h1>` is the most important, one per page.
- **Paragraph**: `<p>`
- **Emphasis**: `<strong>` (important) and `<em>` (stressed)
- **Lists**: `<ul>` for bullets, `<ol>` for numbers, with `<li>` items inside
- **Link**: `<a>`
- **Image**: `<img>`
- **Generic containers**: `<div>` (block) and `<span>` (inline) for grouping

:::match
Q: Match each element to its job.
- `<h1>` | Main heading
- `<a>` | A link
- `<img>` | An image
- `<ul>` | A bulleted list
E: Choose elements for their meaning — that's semantic HTML.
:::

:::quiz
Q: You want a bulleted list of three hobbies. Which element wraps the whole list?
- `<ol>`
- `<ul>` *
- `<li>`
E: `<ul>` is an *unordered* (bulleted) list. Each hobby goes in its own `<li>` inside it. `<ol>` would give you numbers instead.
:::

:::key
HTML describes **what content is**, not how it looks. Choose tags for their *meaning* — a heading because it's a heading, not because you want big text. This is called *semantic* HTML, and it makes your pages easier to style, easier to search, and accessible by default.
:::

## Talk about it

Before you move on, try explaining this out loud (or to a rubber duck 🦆):

> "What is HTML for, and how is it different from CSS and JavaScript? What are the parts of an element?"

If you can answer that in your own words — without reciting a list — you've genuinely got it. If you stumble, scroll back to the analogy that helped most. Understanding beats memorising every time.

## What's next

You can now structure a page. Next, in **CSS Basics**, you'll make it beautiful — colours, spacing, and layout. When you're ready, take the quick quiz to lock in what you've learned.

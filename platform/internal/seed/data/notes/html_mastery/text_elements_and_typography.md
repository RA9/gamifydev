# Text Elements and Typography

HTML gives you a rich vocabulary for text — far beyond just paragraphs and bold. Every tag carries *meaning*, and choosing the right one helps screen readers, search engines, and future developers understand your content. In this lesson you'll learn when and why to use each text element.

## Headings: `<h1>` through `<h6>`

Headings create the outline of your page. Think of them like a table of contents: `<h1>` is the book title, `<h2>` is a chapter, `<h3>` is a section within a chapter, and so on.

```html
<h1>Baking Fundamentals</h1>
  <h2>Ingredients</h2>
    <h3>Flour Types</h3>
    <h3>Leavening Agents</h3>
  <h2>Techniques</h2>
    <h3>Creaming</h3>
    <h3>Folding</h3>
```

Rules to live by:

- **One `<h1>` per page.** It's the main topic.
- **Never skip levels.** Don't jump from `<h1>` to `<h3>` — screen reader users navigate by heading level and skips break their mental model.
- **Don't pick a heading for its size.** Need smaller text? Use CSS. The heading level is about document structure, not appearance.

:::warning
Using `<h3>` because "it looks the right size" is one of the most common semantic mistakes. Pick the heading that makes logical sense in your outline, then style it with CSS.
:::

## Paragraphs

The `<p>` element is for blocks of running text. Browsers add vertical spacing between paragraphs automatically.

```html
<p>HTML was invented by Tim Berners-Lee in 1991. It remains the
backbone of every web page.</p>

<p>Despite decades of evolution, the core tags — headings,
paragraphs, links, and lists — haven't changed much.</p>
```

Don't use `<br>` to fake paragraph spacing. If content is a new thought, it deserves its own `<p>`.

## Strong vs. Bold: `<strong>` vs. `<b>`

Both render text in bold, but they mean different things.

- **`<strong>`** — this text is **important**. Screen readers can emphasize it (some change tone). Search engines may weigh it more heavily.
- **`<b>`** — visually bold with **no extra importance**. Use it for stylistic bold where the meaning doesn't change — like a product name in a review.

```html
<p><strong>Warning:</strong> Do not delete the config file.</p>
<p>The <b>Raspberry Pi 5</b> ships with 8 GB of RAM.</p>
```

:::key
When in doubt, use `<strong>`. If removing the bold wouldn't change the meaning of the sentence, `<b>` is fine. If the bold signals importance, urgency, or a warning, use `<strong>`.
:::

## Emphasis vs. Italic: `<em>` vs. `<i>`

The same logic applies to italics.

- **`<em>`** — stress emphasis. The word's meaning changes if you read it with vocal emphasis: "I *did* finish it" vs. "I did *finish* it."
- **`<i>`** — an alternate voice or mood with no extra emphasis. Use it for technical terms, foreign words, or thoughts.

```html
<p>You <em>must</em> save before closing.</p>
<p>The French word <i lang="fr">bonjour</i> means hello.</p>
```

Notice the `lang` attribute on the `<i>` tag — this tells screen readers to switch pronunciation.

## Blockquotes and Citations

Use `<blockquote>` for long quotations from external sources, and `<cite>` for the title of the work being referenced.

```html
<blockquote>
  <p>The best way to predict the future is to invent it.</p>
  <footer>— <cite>Alan Kay</cite></footer>
</blockquote>
```

For inline quotes within a sentence, use `<q>`:

```html
<p>As Alan Kay said, <q>The best way to predict the future is
to invent it.</q></p>
```

The browser adds quotation marks around `<q>` automatically — don't add your own.

:::tip
`<cite>` is for titles of works (books, articles, films) or attributions. Don't use it for people's names in running text — only when you're crediting a source.
:::

## Code, Pre, and Kbd

When writing about code or keyboard actions, HTML has you covered.

- **`<code>`** — inline code, like a variable name or short snippet.
- **`<pre>`** — preformatted text, preserving whitespace and line breaks. Usually wraps `<code>` for multi-line blocks.
- **`<kbd>`** — a keyboard key or input.
- **`<samp>`** — sample output from a program.
- **`<var>`** — a mathematical or programming variable.

```html
<p>Run <code>npm install</code> in your terminal.</p>

<pre><code>function greet(name) {
  return `Hello, ${name}!`;
}</code></pre>

<p>Press <kbd>Ctrl</kbd> + <kbd>S</kbd> to save.</p>

<p>The terminal prints <samp>Hello, world!</samp></p>

<p>Let <var>x</var> equal the number of users.</p>
```

:::tip
Wrapping `<code>` inside `<pre>` is the standard pattern for code blocks. `<pre>` preserves formatting; `<code>` marks it as code. Together they give you proper semantics and layout.
:::

## Line Breaks: `<br>` vs. CSS

The `<br>` tag forces a line break. It's appropriate for content where line breaks are part of the meaning — addresses and poetry:

```html
<address>
  123 Elm Street<br>
  Springfield, IL 62701<br>
  United States
</address>

<p>
  Roses are red,<br>
  Violets are blue,<br>
  Semantic HTML,<br>
  Is good for you.
</p>
```

**Do not use `<br>` for spacing.** If you need vertical space between blocks, use CSS `margin` or `padding`. Stacking `<br><br><br>` is a code smell.

:::warning
If you find yourself using multiple `<br>` tags in a row, you almost certainly want a new `<p>`, a list, or CSS margin instead. `<br>` is for breaks within content, not space between sections.
:::

## Other Inline Text Elements

HTML has several niche but valuable inline elements:

### `<abbr>` — Abbreviation

```html
<p>The <abbr title="World Health Organization">WHO</abbr>
issued new guidelines.</p>
```

The `title` attribute gives the full expansion. Some browsers show a dotted underline and a tooltip on hover.

### `<mark>` — Highlight

```html
<p>Your search matched: <mark>accessibility</mark> is important.</p>
```

Like a highlighter pen. Use it for search matches or to draw attention to text relevant to the current context.

### `<small>` — Side Comment

```html
<p>Price: $49.99 <small>(tax not included)</small></p>
```

Represents small print, disclaimers, or legal text. It's not about font size — it's about reduced importance.

### `<sub>` and `<sup>` — Subscript and Superscript

```html
<p>Water is H<sub>2</sub>O.</p>
<p>E = mc<sup>2</sup></p>
<p>See footnote<sup><a href="#fn1">1</a></sup></p>
```

Use them for chemical formulas, mathematical exponents, and footnote references — anywhere subscript or superscript is part of the meaning.

### `<del>` and `<ins>` — Deleted and Inserted Text

```html
<p>Price: <del>$99</del> <ins>$79</ins></p>
```

Screen readers can announce these as "deleted" and "inserted," giving context for corrections or price changes.

### `<time>` — Machine-Readable Dates

```html
<p>Published on <time datetime="2025-03-15">March 15, 2025</time></p>
```

The `datetime` attribute provides a machine-readable format while the user sees the friendly version. Search engines can parse this for structured data.

## Choosing the Right Element: A Decision Guide

Ask these questions in order:

1. **Is it a section title?** → Use the correct heading level (`h1`–`h6`).
2. **Is it a block of prose?** → `<p>`.
3. **Does bold mean "this is important"?** → `<strong>`. Just stylistic? → `<b>`.
4. **Does italic mean "stress this word"?** → `<em>`. Alternate voice? → `<i>`.
5. **Is it a long quote from elsewhere?** → `<blockquote>`. Short inline? → `<q>`.
6. **Is it code?** → `<code>`, wrapped in `<pre>` for blocks.
7. **Is it a keyboard shortcut?** → `<kbd>`.
8. **None of the above?** → Check `<abbr>`, `<mark>`, `<small>`, `<time>`, `<sub>`, `<sup>`, `<del>`, `<ins>`.
9. **Still nothing?** → A plain `<span>` with a class is fine — it carries no meaning, which is exactly right when none of these elements apply.

:::quiz
Q: What's the difference between `<strong>` and `<b>`?
- `<strong>` is bold and `<b>` is italic
- `<strong>` indicates importance; `<b>` is purely visual bold with no added meaning *
- They are identical and interchangeable
E: `<strong>` conveys semantic importance — screen readers may emphasize it and search engines may weight it. `<b>` is bold styling without additional meaning.
:::

:::quiz
Q: Which element is correct for displaying a multi-line code snippet?
- `<p><code>…</code></p>`
- `<pre><code>…</code></pre>` *
- `<code><br>…</code>`
E: `<pre>` preserves whitespace and line breaks, and wrapping `<code>` inside it marks the content as code — the standard combination for code blocks.
:::

## Recap

- Headings (`h1`–`h6`) build your page outline — one `h1`, never skip levels, and don't pick by size.
- `<strong>` means important; `<b>` is visual-only bold. `<em>` means stress emphasis; `<i>` is alternate voice.
- Use `<blockquote>` for long quotes, `<q>` for inline quotes, and `<cite>` for source titles.
- `<pre><code>` is the standard pattern for code blocks; `<kbd>` marks keyboard keys.
- `<br>` is for meaningful line breaks (addresses, poetry), not for spacing — use CSS for that.
- Reach for `<abbr>`, `<mark>`, `<small>`, `<time>`, `<sub>`, `<sup>`, `<del>`, and `<ins>` when they fit.

**Next up:** Links and Navigation — connecting your pages and guiding users around your site.

# CSS Basics

You've built structure with HTML; now it's time to make it *look* like something. CSS (Cascading Style Sheets) is the language that controls colour, spacing, fonts, and layout — turning a plain skeleton of text into a real, designed page.

## What CSS is for

HTML says *what* something is ("this is a heading", "this is a paragraph"). CSS says *how it should look* ("headings are dark blue and bold", "paragraphs have comfortable line spacing"). They're two different jobs, and keeping them separate is one of the most important habits you can build.

:::analogy
If HTML is the skeleton of a page, CSS is the skin, clothes, and styling. The same skeleton can wear a business suit or a hoodie — same structure, totally different look. That's exactly what CSS lets you do: restyle content without rewriting it.
:::

## Anatomy of a CSS rule

Everything in CSS is built from **rules**. A rule has two parts: a **selector** (what to style) and a **declaration block** (how to style it).

```css
h1 {
  color: navy;
  font-size: 32px;
}
```

Let's name every piece, because you'll hear these words constantly:

- `h1` is the **selector** — it picks which elements this rule applies to.
- The `{ }` curly braces hold the **declaration block**.
- `color: navy;` is one **declaration**.
- Inside a declaration, `color` is the **property** and `navy` is the **value**, separated by a colon.
- Every declaration ends with a **semicolon**.

```css
selector {
  property: value;
  property: value;
}
```

:::warning
The single most common beginner mistake is a missing semicolon or a missing closing brace. CSS won't crash — it'll just silently ignore the broken part, leaving you confused about why nothing changed. When a style "doesn't work", check your punctuation first.
:::

## The three ways to include CSS

You can attach CSS to your HTML in three ways. They are not equally good — we'll recommend one strongly.

### 1. Inline styles (avoid)

A `style` attribute right on the element:

```html
<p style="color: red; font-size: 18px;">Hello!</p>
```

This works, but it's the worst option: it mixes structure and style, can't be reused, and is hard to override. Save it for quick one-off experiments only.

### 2. Internal stylesheet (okay for tiny pages)

A `<style>` block inside the `<head>`:

```html
<head>
  <style>
    p { color: red; font-size: 18px; }
  </style>
</head>
```

Better — the styles are in one place — but they only apply to this one HTML file.

### 3. External stylesheet (recommended)

A separate `.css` file, linked from the `<head>`:

```html
<head>
  <link rel="stylesheet" href="styles.css">
</head>
```

```css
/* styles.css */
p {
  color: red;
  font-size: 18px;
}
```

:::key
Use **external stylesheets**. One CSS file can style your entire site, browsers cache it so pages load faster, and your HTML stays clean. This is what professionals do almost all the time.
:::

## Selectors: choosing what to style

Selectors are how you target elements. Here are the ones you'll use daily.

```css
/* Element selector — every paragraph */
p { color: #333; }

/* Class selector — any element with class="card" */
.card { border: 1px solid gray; }

/* ID selector — the one element with id="header" */
#header { background: black; }

/* Grouping — apply the same rule to several selectors */
h1, h2, h3 { font-family: Georgia, serif; }

/* Descendant — <a> elements inside a <nav> */
nav a { text-decoration: none; }
```

```html
<header id="header">...</header>
<div class="card">A card</div>
<nav><a href="#">Home</a></nav>
```

**Classes vs IDs:** a class can be used on many elements; an ID must be unique to one element per page. In practice you'll reach for **classes** almost always — they're reusable and flexible.

### Pseudo-classes

Pseudo-classes target an element in a particular *state* or *position*, using a colon:

```css
/* When the user hovers their mouse over a link */
a:hover { color: orange; }

/* The first <li> inside its parent */
li:first-child { font-weight: bold; }

/* The last item, and every other row */
li:last-child { border-bottom: none; }
tr:nth-child(even) { background: #f5f5f5; }
```

:::example
`:hover` is your first taste of interactivity *without JavaScript*. Buttons that change colour when you point at them, links that underline on hover — all pure CSS.
:::

## The cascade and specificity

What happens when two rules try to style the same element differently? CSS resolves the conflict using two ideas.

**The cascade** means *later rules can override earlier ones* (when specificity is equal). Source order matters.

```css
p { color: blue; }
p { color: green; }   /* wins — it comes later */
```

**Specificity** means *more specific selectors beat less specific ones*, regardless of order. A rough ranking from weakest to strongest:

1. Element selectors (`p`, `h1`) — weakest
2. Class selectors (`.card`) and pseudo-classes (`:hover`)
3. ID selectors (`#header`) — strong
4. Inline `style="..."` — stronger still
5. `!important` — overrides everything (use as a last resort)

```css
p { color: blue; }          /* low specificity */
.intro { color: green; }    /* beats the above for <p class="intro"> */
#lead { color: red; }       /* beats both for <p id="lead"> */
```

:::warning
Don't reach for `!important` to "win" a fight with specificity. It works, but it makes future styles even harder to override and quickly turns a stylesheet into a mess. Prefer a more specific (but reasonable) selector instead.
:::

## Inheritance

Some properties are passed down from a parent element to its children automatically. Set a font on `<body>` and every paragraph inside inherits it:

```css
body {
  font-family: Arial, sans-serif;
  color: #222;
}
```

Now every `<p>`, `<li>`, and `<span>` uses Arial and that text colour without you repeating yourself. **Text-related** properties (`color`, `font-family`, `font-size`, `line-height`) inherit; **box-related** properties (`border`, `margin`, `padding`, `background`) do *not*.

## The box model (in depth)

Here is the idea that explains 90% of "why is there a gap there?" questions. **Every element is a rectangular box** made of four layers, from the inside out:

1. **Content** — the text or image itself.
2. **Padding** — space *inside* the box, between content and border.
3. **Border** — a line around the padding.
4. **Margin** — space *outside* the box, pushing other elements away.

```css
.box {
  width: 200px;
  padding: 20px;
  border: 5px solid black;
  margin: 30px;
}
```

:::analogy
Think of a framed photo on a wall. The **content** is the photo. The **padding** is the white mat around it inside the frame. The **border** is the frame itself. The **margin** is the empty wall space keeping it away from other frames.
:::

You can set each side individually or all at once:

```css
.box {
  margin: 10px;                 /* all four sides */
  margin: 10px 20px;            /* top/bottom, left/right */
  margin: 10px 20px 30px 40px;  /* top, right, bottom, left (clockwise) */
  padding-left: 16px;           /* just one side */
}
```

### box-sizing: border-box (the pro setting)

Here's the gotcha. By default, `width` sets only the **content** width. Padding and border are *added on top*, so a box you declared as `200px` wide can actually render much wider:

```css
.box {
  width: 200px;
  padding: 20px;   /* +40px total */
  border: 5px;     /* +10px total */
  /* Actual rendered width: 250px! */
}
```

The fix — used in virtually every professional codebase — is `box-sizing: border-box`, which makes `width` *include* padding and border:

```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```

:::key
Put that `* { box-sizing: border-box; }` rule at the top of every stylesheet. Now `width: 200px` means the box is exactly 200px wide, padding and all. It makes layouts predictable and saves endless headaches.
:::

## Colours

You can describe colours four common ways:

```css
.a { color: tomato; }                 /* named — 140+ keywords */
.b { color: #ff6347; }                /* hex — same tomato */
.c { color: rgb(255, 99, 71); }       /* red, green, blue 0–255 */
.d { color: rgba(255, 99, 71, 0.5); } /* + alpha (transparency) */
.e { color: hsl(9, 100%, 64%); }      /* hue, saturation, lightness */
```

- **Hex** (`#rrggbb`) is the most common in real code.
- **rgba** is handy when you need transparency.
- **hsl** is the easiest to *tweak by hand* — bump the lightness to get a lighter shade of the same colour.

## Units: px vs em vs rem vs %

Sizes need units, and choosing the right one matters:

```css
.fixed   { font-size: 16px;  }  /* absolute pixels */
.relemt  { font-size: 1.5em; }  /* relative to PARENT font-size */
.relroot { font-size: 1.5rem;}  /* relative to ROOT (<html>) font-size */
.fluid   { width: 50%;       }  /* percentage of the parent */
```

- **px** — absolute and predictable. Great for borders and fine details.
- **em** — relative to the *parent's* font size. Powerful but can compound unexpectedly when nested.
- **rem** — relative to the *root* font size (the `<html>` element). Predictable *and* scalable. The modern favourite for font sizes and spacing.
- **%** — relative to the parent's size. Ideal for fluid widths.

:::tip
A great default: use **rem** for font sizes and spacing, **%** (or later, flex/grid) for widths, and **px** for tiny details like `1px` borders. If a user increases their browser's default font size for accessibility, rem-based layouts scale with them.
:::

## Typography

Text styling is most of what CSS does day to day:

```css
body {
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
  font-size: 1rem;
  font-weight: 400;     /* 400 = normal, 700 = bold */
  line-height: 1.6;     /* space between lines */
  text-align: left;     /* left | center | right | justify */
}

h1 { font-weight: 700; }
```

That comma-separated list in `font-family` is a **font stack**: the browser tries each font in order and uses the first one available, falling back to a generic family (`sans-serif`, `serif`, `monospace`) at the end. Always end with a generic fallback.

:::tip
A `line-height` of around `1.5`–`1.6` (unitless) makes body text far more readable than the cramped default. Unitless line-height scales correctly with font size, so prefer `1.6` over `26px`.
:::

## display: block, inline, inline-block

Every element has a default `display` value that controls how it flows:

```css
.block        { display: block; }        /* full width, stacks vertically */
.inline       { display: inline; }       /* flows in a line, ignores width/height */
.inline-block { display: inline-block; } /* flows in a line, but respects width/height */
```

- **block** elements (`<div>`, `<p>`, `<h1>`) take the full available width and start on a new line.
- **inline** elements (`<a>`, `<span>`, `<strong>`) sit in the text flow and ignore `width`/`height` and top/bottom margins.
- **inline-block** is the best of both: it flows inline but lets you set a width, height, and full padding.

```html
<a href="#" class="button">Click me</a>
```

```css
.button {
  display: inline-block;
  padding: 10px 20px;
  background: navy;
  color: white;
}
```

There's a fourth value, `display: flex`, that unlocks real layout power — and that's exactly where the next lesson begins.

:::quiz
Q: Why do professionals add `* { box-sizing: border-box; }` to their stylesheets?
- It makes all elements display side by side
- It makes an element's declared `width` include its padding and border *
- It removes the default margin from every element
E: With border-box, `width: 200px` stays 200px even after you add padding and a border, making layouts predictable.
:::

:::quiz
Q: Two rules target the same paragraph: `p { color: blue; }` and `.intro { color: green; }`. The paragraph is `<p class="intro">`. What colour is it?
- Blue, because element selectors come first
- Green, because a class selector is more specific than an element selector *
- Red, because conflicting rules cancel out
E: Specificity decides conflicts before source order; a class beats a bare element selector, so green wins.
:::

:::predict
Q: A `<span>` has `display: inline` and you set `width: 300px`. What happens to its width?
- The span becomes exactly 300px wide
- The `width` is ignored and the span stays only as wide as its text *
- The span breaks onto its own line
E: Inline elements ignore `width` and `height`; switch to `inline-block` or `block` if you need to size them.
:::

## Recap

- A CSS rule is a **selector** plus a **declaration block** of `property: value;` pairs.
- Prefer **external stylesheets** over internal or inline styles.
- Key selectors: element, `.class`, `#id`, grouping (`a, b`), descendant (`nav a`), and pseudo-classes like `:hover` and `:first-child`.
- Conflicts resolve by the **cascade** (later wins) and **specificity** (more specific wins); avoid `!important`.
- Text properties **inherit**; box properties don't.
- The **box model** is content → padding → border → margin; set `box-sizing: border-box` everywhere.
- Colours: named, hex, rgb/rgba, hsl. Units: prefer **rem** for type/spacing, **%** for widths, **px** for fine details.
- `display` controls flow: **block** stacks, **inline** flows, **inline-block** flows but is sizable.

**Next up:** CSS Layout with Flexbox — the modern, sane way to arrange boxes in rows and columns and center anything.

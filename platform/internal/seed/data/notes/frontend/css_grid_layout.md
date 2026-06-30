# CSS Grid Layout

CSS Grid is the most powerful layout tool in CSS. Where Flexbox arranges things in a single line, Grid lets you control rows *and* columns at the same time, which makes it perfect for full-page layouts and structured galleries.

## When Grid beats Flexbox: 2-D vs 1-D

The single most useful way to choose between the two:

- **Flexbox is one-dimensional.** It lays items out along *one* axis at a time, a row OR a column. It is brilliant for navbars, button groups, and centering one thing.
- **Grid is two-dimensional.** It controls rows AND columns together, so the items line up into a clean matrix.

:::analogy
Think of Flexbox as people lining up at a coffee shop, one queue, flexible spacing. Grid is the seating chart in a theater, fixed rows and columns where every seat has an address like "Row C, Seat 7."
:::

If your content naturally forms a table-like structure (gallery, dashboard, page skeleton), reach for Grid. If it is a single strip of items that should flow and wrap, Flexbox is simpler.

## Turning on Grid

You make an element a grid container with `display: grid`. Its *direct children* become grid items automatically.

```css
.container {
  display: grid;
}
```

By itself that does almost nothing visible. The magic comes from defining the tracks (the columns and rows).

## Defining columns and rows

`grid-template-columns` defines how many columns there are and how wide each is. `grid-template-rows` does the same for rows.

```css
.container {
  display: grid;
  grid-template-columns: 200px 200px 200px; /* three fixed columns */
  grid-template-rows: 100px 100px;          /* two fixed rows */
}
```

That creates a 3-column, 2-row grid. The number of *values* you list is the number of tracks you get.

```html
<div class="container">
  <div>1</div><div>2</div><div>3</div>
  <div>4</div><div>5</div><div>6</div>
</div>
```

You usually only need to define columns. Rows are created automatically as items wrap onto new lines.

## The `fr` unit (the secret sauce)

Fixed pixel columns do not adapt to screen size. The `fr` unit means "one fraction of the leftover space." It is what makes grids fluid.

```css
.container {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr; /* three equal flexible columns */
}
```

You can mix ratios. `2fr 1fr` makes the first column twice as wide as the second:

```css
.sidebar-layout {
  display: grid;
  grid-template-columns: 1fr 3fr; /* sidebar gets 1 part, content gets 3 */
}
```

And you can mix `fr` with fixed sizes. Here the sidebar is exactly 250px and the content takes whatever remains:

```css
.app {
  display: grid;
  grid-template-columns: 250px 1fr;
}
```

:::key
`fr` distributes the space *left over* after fixed sizes, padding, and gaps are subtracted. That is why `250px 1fr` always fits perfectly, no overflow math required.
:::

## Spacing with `gap`

`gap` puts consistent space *between* tracks without adding awkward margins on the outer edges.

```css
.container {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;              /* same gap for rows and columns */
}
```

You can set row and column gaps separately:

```css
.container {
  gap: 24px 16px; /* row-gap column-gap */
}
```

:::tip
`gap` works in Flexbox too now, but in Grid it is essential. Always prefer `gap` over margins on items, it never leaves a stray edge margin and is far easier to maintain.
:::

## `repeat()` to avoid repetition

Typing `1fr 1fr 1fr 1fr` gets old fast. `repeat()` does it for you:

```css
.container {
  display: grid;
  grid-template-columns: repeat(4, 1fr); /* four equal columns */
}
```

You can combine repeat with other tracks:

```css
.container {
  grid-template-columns: 200px repeat(3, 1fr); /* fixed first, then three flexible */
}
```

## `minmax()` for flexible-but-bounded tracks

`minmax(min, max)` sets a lower and upper bound for a track. It is the key to columns that shrink and grow within sensible limits.

```css
.container {
  display: grid;
  /* each column is at least 150px, at most 1fr of the space */
  grid-template-columns: repeat(3, minmax(150px, 1fr));
}
```

This says "never let a column get narrower than 150px, but happily let it stretch."

## Responsive auto-grids (the one-liner you will reuse forever)

Combine `repeat()`, `auto-fit`, and `minmax()` and you get a grid that *automatically* changes its column count based on available width, with **no media queries**:

```css
.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}
```

Read it as: "fit as many columns as you can, each at least 200px wide; share the leftover space equally." On a phone you get 1 column, on a tablet 2 or 3, on a wide monitor 5 or 6. This single rule replaces a stack of breakpoints.

:::warning
`auto-fit` collapses empty tracks so existing items stretch to fill. `auto-fill` keeps empty phantom tracks, leaving gaps where items *could* go. For most galleries you want `auto-fit`.
:::

## Spanning columns and rows

By default each item fills one cell. To make an item cover multiple tracks, use `grid-column` and `grid-row` with `span`.

```css
.featured {
  grid-column: span 2; /* take up two columns */
}
.tall {
  grid-row: span 2;    /* take up two rows */
}
```

You can also span using explicit line numbers. Grid lines are numbered starting at 1 on the left/top:

```css
.banner {
  grid-column: 1 / 3; /* start at line 1, end at line 3 = spans columns 1 and 2 */
}
```

```html
<div class="gallery">
  <div class="featured">Big photo</div>
  <div>Photo</div>
  <div>Photo</div>
  <div>Photo</div>
</div>
```

## Named areas for page layout

For whole-page skeletons, `grid-template-areas` lets you literally draw the layout in text. You name each region, then place items by name. This is wonderfully readable.

```css
.page {
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-template-rows: auto 1fr auto;
  grid-template-areas:
    "header  header"
    "sidebar main"
    "footer  footer";
  min-height: 100vh;
  gap: 12px;
}

.page > header  { grid-area: header; }
.page > nav     { grid-area: sidebar; }
.page > main    { grid-area: main; }
.page > footer  { grid-area: footer; }
```

```html
<div class="page">
  <header>Logo + nav</header>
  <nav>Sidebar links</nav>
  <main>Page content</main>
  <footer>Footer</footer>
</div>
```

The quotes are rows; each word is a column cell. Repeating a name (like `header header`) makes that area span those cells. A period `.` marks an empty cell.

## Full example: image gallery

```css
.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  grid-auto-rows: 180px;
  gap: 12px;
}
.photo-grid .hero {
  grid-column: span 2;
  grid-row: span 2;
}
.photo-grid div {
  background: #ddd;
  border-radius: 8px;
}
```

```html
<div class="photo-grid">
  <div class="hero">Featured</div>
  <div>1</div><div>2</div><div>3</div>
  <div>4</div><div>5</div><div>6</div>
</div>
```

## Full example: holy-grail page

```css
.layout {
  display: grid;
  grid-template-areas:
    "header header header"
    "nav    main   aside"
    "footer footer footer";
  grid-template-columns: 180px 1fr 220px;
  grid-template-rows: auto 1fr auto;
  min-height: 100vh;
  gap: 10px;
}
.layout > header { grid-area: header; }
.layout > nav    { grid-area: nav; }
.layout > main   { grid-area: main; }
.layout > aside  { grid-area: aside; }
.layout > footer { grid-area: footer; }
```

## Flexbox vs Grid: a quick decision note

- **Use Flexbox** when content flows in one direction and item sizes should drive the layout: navbars, toolbars, tag lists, centering one element.
- **Use Grid** when you need rows and columns to align together: page layouts, dashboards, galleries, forms with label/field columns.
- They **work together**. A common pattern is Grid for the overall page and Flexbox *inside* a grid cell (like a navbar within the header area).

:::quiz
Q: Which property makes columns that automatically adjust their count to fit the screen without media queries?
- grid-template-rows: repeat(3, 1fr)
- grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)) *
- display: flex
E: `auto-fit` plus `minmax()` packs as many columns as fit and stretches them to share leftover space, so the count changes with the viewport on its own.
:::

:::fill
Fill in the blank: The ____ unit represents one fraction of the leftover space in a grid track.
A: fr
:::

:::predict
What will `grid-template-columns: 1fr 2fr;` produce?
A: Two columns where the second is twice as wide as the first; together they fill the full container width.
:::

## Recap

- `display: grid` turns an element into a 2-D grid container; its direct children become grid items.
- Define tracks with `grid-template-columns` / `grid-template-rows`; the `fr` unit splits leftover space.
- Use `gap` for clean spacing, `repeat()` to avoid repetition, and `minmax()` to bound track sizes.
- `repeat(auto-fit, minmax(200px, 1fr))` is a responsive grid in one line.
- Span items with `grid-column: span 2` / `grid-row: span 2`, and lay out pages with named `grid-template-areas`.
- Choose Grid for 2-D structure, Flexbox for 1-D flow, and combine them freely.

**Next up:** Responsive Design, where you will make these layouts adapt beautifully to every screen size.

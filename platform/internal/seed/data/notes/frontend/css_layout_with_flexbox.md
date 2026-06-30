# CSS Layout with Flexbox

You can style individual boxes now — but real pages are about *arranging* boxes: navigation bars, rows of cards, sidebars. Flexbox is the modern CSS tool that makes laying out elements in a row or column easy, flexible, and (finally) sane.

## The problem flexbox solves

For years, putting two boxes side by side or centering something vertically was genuinely painful. Developers abused `float`, fought with `display: inline-block` whitespace bugs, and memorised hacks just to center a box.

Flexbox replaced all of that. With a couple of lines you can align items, distribute space evenly, and adapt to different screen sizes — no hacks.

:::analogy
Think of a flex container as a row of books on a shelf. You can push them all to the left, spread them out evenly, center them, or stand them at the same height. Flexbox gives you simple controls for exactly those everyday arrangements.
:::

## display: flex — container vs items

Flexbox always involves two roles:

- The **flex container** — the parent element you put `display: flex` on.
- The **flex items** — its *direct children*, which now lay out according to flex rules.

```html
<div class="container">
  <div class="item">A</div>
  <div class="item">B</div>
  <div class="item">C</div>
</div>
```

```css
.container {
  display: flex;   /* turns on flexbox for the children */
}
```

The moment you add `display: flex`, those three `<div>`s — which would normally stack vertically — line up in a **row**.

:::key
Flex properties split into two groups: ones you put on the **container** (`justify-content`, `align-items`, `gap`, `flex-direction`, `flex-wrap`) and ones you put on the **items** (`flex`, `align-self`). When something isn't working, first ask: "is this a container property or an item property?"
:::

## Main axis vs cross axis

This is the concept that makes everything else click. A flex container has two axes:

- The **main axis** — the direction items flow (by default, left to right).
- The **cross axis** — perpendicular to the main axis (by default, top to bottom).

`justify-content` aligns items along the **main** axis. `align-items` aligns them along the **cross** axis. Which direction each axis points depends on `flex-direction`.

## flex-direction

Controls which way the main axis runs:

```css
.container { flex-direction: row; }            /* default: left → right */
.container { flex-direction: column; }         /* top → bottom */
.container { flex-direction: row-reverse; }    /* right → left */
.container { flex-direction: column-reverse; } /* bottom → top */
```

:::warning
When you switch to `flex-direction: column`, the main axis becomes vertical — so `justify-content` now controls *vertical* positioning and `align-items` controls *horizontal*. They effectively swap. This trips up everyone at first, so keep the axis idea in mind.
:::

## justify-content (main axis)

Distributes items along the main axis:

```css
.container { justify-content: flex-start; }    /* packed at the start (default) */
.container { justify-content: flex-end; }      /* packed at the end */
.container { justify-content: center; }        /* centered */
.container { justify-content: space-between; }  /* first & last at edges, gaps between */
.container { justify-content: space-around; }  /* equal space around each item */
.container { justify-content: space-evenly; }  /* equal space everywhere */
```

`space-between` is the workhorse — it's how you push a logo to the left and links to the right in a nav bar.

## align-items (cross axis)

Aligns items along the cross axis:

```css
.container { align-items: stretch; }     /* default: fill the cross axis */
.container { align-items: flex-start; }  /* align to the top (in a row) */
.container { align-items: flex-end; }    /* align to the bottom */
.container { align-items: center; }      /* center vertically (in a row) */
.container { align-items: baseline; }    /* align text baselines */
```

## gap

The cleanest way to add space *between* flex items — no fiddly margins:

```css
.container {
  display: flex;
  gap: 16px;          /* same gap in both directions */
  /* gap: 16px 8px;   row-gap then column-gap */
}
```

:::tip
Use `gap` instead of adding `margin-right` to every item and then awkwardly removing it from the last one. `gap` only puts space *between* items, never on the outer edges. It's simpler and bug-free.
:::

## flex-wrap

By default, flex items all try to squeeze onto one line, shrinking to fit. `flex-wrap: wrap` lets them spill onto new lines when there's no room:

```css
.container {
  display: flex;
  flex-wrap: wrap;   /* nowrap (default) | wrap */
  gap: 16px;
}
```

This is essential for responsive card grids: as the window narrows, cards drop to the next row instead of becoming uncomfortably thin.

## align-self (per item)

Override `align-items` for a single item:

```css
.special {
  align-self: flex-end;   /* this one item aligns differently */
}
```

Everything else in the container obeys `align-items`; the one item with `align-self` does its own thing.

## The flex shorthand

Items don't have to be a fixed size — they can *flex* to fill available space. The `flex` property is shorthand for three things:

```css
.item {
  flex: 1;
  /* shorthand for: */
  /* flex-grow: 1;     how much to GROW into extra space */
  /* flex-shrink: 1;   how much to SHRINK when space is tight */
  /* flex-basis: 0;    the starting size before growing/shrinking */
}
```

In plain English:

- **flex-grow** — if there's leftover space, how greedily does this item expand to claim it? `0` means "don't grow"; higher numbers grow more.
- **flex-shrink** — if items overflow, how readily does this item give up space?
- **flex-basis** — the item's ideal size before growing or shrinking kicks in.

Common values you'll actually type:

```css
.item { flex: 1; }      /* grow to share space equally with other flex:1 items */
.item { flex: 2; }      /* take twice as much space as a flex:1 sibling */
.item { flex: 0 0 200px; } /* don't grow, don't shrink, stay 200px (a fixed sidebar) */
.item { flex: 1 1 300px; } /* start at 300px but grow/shrink as needed */
```

:::example
Put `flex: 1` on three sibling items and they each take exactly one-third of the row, automatically, no matter how wide the container is. Resize the window and they stay equal. That's the magic.
:::

## Centering anything — the classic

The thing flexbox is most famous for. To center a child both horizontally *and* vertically:

```html
<div class="hero">
  <h1>Welcome</h1>
</div>
```

```css
.hero {
  display: flex;
  justify-content: center;  /* center on the main (horizontal) axis */
  align-items: center;      /* center on the cross (vertical) axis */
  height: 100vh;            /* full viewport height so there's room to center in */
}
```

:::key
"Center a div" used to be a running joke in web development. With flexbox it's three lines: `display: flex; justify-content: center; align-items: center;`. Memorise this trio — you'll use it constantly.
:::

## Pattern 1: a nav bar (logo left, links right)

```html
<nav class="navbar">
  <div class="logo">GamifyDev</div>
  <ul class="links">
    <li><a href="#">Home</a></li>
    <li><a href="#">Courses</a></li>
    <li><a href="#">About</a></li>
  </ul>
</nav>
```

```css
.navbar {
  display: flex;
  justify-content: space-between; /* logo left, links pushed right */
  align-items: center;           /* vertically centered */
  padding: 16px 24px;
  background: #1a1a2e;
}
.navbar .links {
  display: flex;     /* the <ul> is ALSO a flex container */
  gap: 24px;
  list-style: none;
  margin: 0;
  padding: 0;
}
.navbar a {
  color: white;
  text-decoration: none;
}
```

Notice the nav is a flex container, and the `<ul>` inside it is *also* a flex container for its own links. Nesting flex containers is completely normal and very common.

## Pattern 2: a row of equal cards

```html
<div class="cards">
  <div class="card">Card 1</div>
  <div class="card">Card 2</div>
  <div class="card">Card 3</div>
</div>
```

```css
.cards {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;   /* cards drop to a new row on narrow screens */
}
.card {
  flex: 1 1 250px;   /* start ~250px, grow to fill, shrink if needed */
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}
```

`flex: 1 1 250px` is the responsive-card sweet spot: each card aims for about 250px, shares leftover space equally, and wraps gracefully when the row runs out of room.

## Pattern 3: sidebar + content split

```html
<div class="layout">
  <aside class="sidebar">Menu</aside>
  <main class="content">Main content goes here</main>
</div>
```

```css
.layout {
  display: flex;
  gap: 24px;
}
.sidebar {
  flex: 0 0 240px;  /* fixed 240px: don't grow, don't shrink */
}
.content {
  flex: 1;          /* take all remaining space */
}
```

The sidebar stays a steady 240px while the main content stretches to fill whatever's left. This `flex: 0 0 <fixed>` plus `flex: 1` combo is one of the most useful layout recipes you'll learn.

## When to reach for Grid instead

Flexbox is **one-dimensional** — it lays things out in a single row *or* a single column at a time (wrapping is still essentially one line repeated). That's perfect for nav bars, button groups, and toolbars.

But when you need a real **two-dimensional** layout — rows *and* columns that line up together, like a photo gallery or a page template with header, sidebar, and footer all aligned to a grid — Flexbox starts to fight you. That's the job of CSS Grid, which is exactly the next lesson.

:::tip
A handy rule of thumb: **content in a line → Flexbox. A true grid of rows and columns → Grid.** They also work beautifully together — Grid for the overall page skeleton, Flexbox for aligning items inside each region.
:::

:::quiz
Q: Which property aligns flex items along the **cross axis** (vertically, in a default row)?
- justify-content
- align-items *
- flex-direction
E: `justify-content` handles the main axis; `align-items` handles the cross axis, which is vertical in a default row.
:::

:::quiz
Q: You want a fixed-width 240px sidebar that never grows or shrinks. Which `flex` value fits?
- flex: 1
- flex: 0 0 240px *
- flex: 1 1 240px
E: `flex: 0 0 240px` sets grow and shrink to 0 with a 240px basis, locking the width.
:::

:::predict
Q: A container has `display: flex; justify-content: center; align-items: center;` and one child. Where does the child sit?
- In the top-left corner
- Dead center, both horizontally and vertically *
- Stretched across the full width
E: Centering on both axes is the classic flex trick — the child lands in the middle of the container.
:::

## Recap

- Add `display: flex` to a **container**; its direct children become **flex items**.
- Flexbox has a **main axis** (set by `flex-direction`) and a perpendicular **cross axis**.
- `justify-content` aligns along the main axis; `align-items` along the cross axis.
- Use `gap` for spacing between items, and `flex-wrap: wrap` to let items flow onto new lines.
- The `flex` shorthand (`flex-grow flex-shrink flex-basis`) controls how items size and share space; `flex: 1` shares equally, `flex: 0 0 240px` stays fixed.
- Center anything with `display: flex; justify-content: center; align-items: center;`.
- Reach for **Grid** when you need true two-dimensional row-and-column layouts.

**Next up:** CSS Grid Layout — building two-dimensional page structures where rows and columns line up perfectly.

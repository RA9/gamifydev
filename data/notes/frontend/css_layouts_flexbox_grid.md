# CSS Layouts: Flexbox & Grid

Styling individual components is one skill. Arranging entire interfaces is another.

This lesson is about the two layout systems modern frontend developers rely on constantly:

- **Flexbox** for one-dimensional layout
- **Grid** for two-dimensional layout

Master these and you stop fighting CSS for positioning.

## Think in layout problems first

Before choosing a layout tool, ask:

- Am I arranging items in a **row or column**? → flexbox
- Am I arranging items across **rows and columns together**? → grid

That simple question removes a lot of confusion.

## Flexbox: one dimension at a time

Flexbox is ideal when you want to line items up in a row or a column and control how they share space.

```css
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}
```

High-value flex properties:

- `display: flex`
- `flex-direction`
- `justify-content`
- `align-items`
- `gap`
- `flex-wrap`

Example:

```css
.cards {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}
```

This creates a row of cards that can wrap onto new lines when space runs out.

:::fill
Q: Complete the rule so the container becomes a flex container.
`.nav { display: ___; }`
- flex *
- block
- grid
E: `display: flex` switches the element into flex layout mode.
:::

## The two flexbox axes

Flexbox always works across two axes:

- **main axis** — controlled by `flex-direction`
- **cross axis** — the perpendicular axis

If `flex-direction: row`, then:

- `justify-content` moves items left/right along the row
- `align-items` moves items up/down across the row

```css
.hero-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.75rem;
}
```

:::reorder
Q: Put the steps for centering a box with flexbox in a sensible order.
- display: flex;
- justify-content: center;
- align-items: center;
E: First make the parent a flex container, then center children on both axes.
:::

## Grid: rows and columns together

Grid shines when you need more structured page layout.

```css
.feature-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}
```

Useful grid properties:

- `display: grid`
- `grid-template-columns`
- `grid-template-rows`
- `gap`
- `grid-column`
- `grid-row`

A common responsive pattern:

```css
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(16rem, 1fr));
  gap: 1rem;
}
```

This says:

- create as many columns as fit
- no column should get smaller than `16rem`
- otherwise let columns stretch evenly

## When to choose flexbox vs grid

Use **flexbox** when:

- building nav bars
- aligning buttons
- spacing items inside a card
- creating small UI rows/columns

Use **grid** when:

- building page sections
- laying out card galleries
- creating dashboard panels
- controlling both rows and columns together

:::quiz
Q: You're building a gallery of six cards in responsive columns. Which layout system is usually the better first choice?
- Flexbox
- Grid *
- A table
E: Grid is usually the best fit when you want structured columns and rows for a repeated layout like a card gallery.
:::

## Layout example: a landing page section

HTML:

```html
<section class="hero">
  <div>
    <h1>Design systems for small teams</h1>
    <p>Move faster with a shared UI language.</p>
  </div>
  <img src="hero.png" alt="Product dashboard preview" />
</section>
```

CSS:

```css
.hero {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 2rem;
  align-items: center;
}
```

That two-column grid works well when text and media need equal importance.

## Let content size influence layout

A common beginner mistake is forcing rigid widths everywhere.

Instead of:

```css
.card {
  width: 340px;
}
```

Prefer flexible patterns:

```css
.card {
  width: 100%;
  max-width: 21rem;
}
```

Layout becomes much easier when components can grow and shrink within sensible bounds.

## `gap` beats margin hacks

When spacing items inside flex or grid layouts, prefer `gap`.

```css
.actions {
  display: flex;
  gap: 0.75rem;
}
```

Why this is better than manual margins:

- cleaner code
- no awkward last-item spacing
- easier to change globally

## Mini practice — build three layout patterns

Try building these with HTML and CSS:

1. a nav bar with a logo on the left and links on the right
2. a three-card features section
3. a hero section with text on the left and image on the right

Suggested tools:

- nav bar → flexbox
- features section → grid
- hero section → grid, or flexbox if the structure is simpler

## Common layout mistakes to avoid

- using flexbox for everything, even when grid is clearer
- hard-coding widths that break on smaller screens
- using margins everywhere instead of `gap`
- centering content without understanding the main and cross axes
- trying to "eyeball" layout without first deciding the structure

:::key
The layout question always comes before the CSS question. Decide the relationship between the items first, then choose flexbox or grid to express that relationship.
:::

## What's next

In **Responsive Design & Media Queries**, you'll make these layouts adapt across devices instead of only looking good at one screen size.
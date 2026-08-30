# The Box Model

Every element on a web page is a rectangular box. The **box model** defines how much space that box takes up — and it's built from four layers, like nesting boxes inside each other. Understanding this model is the foundation for every layout decision you'll ever make in CSS.

## The Four Layers

From inside out:

1. **Content** — the actual text, image, or child elements.
2. **Padding** — transparent space between the content and the border.
3. **Border** — a visible (or invisible) edge around the padding.
4. **Margin** — transparent space outside the border, pushing other elements away.

```
┌─────────────────────────── margin ──────────────────────────┐
│  ┌──────────────────────── border ────────────────────────┐  │
│  │  ┌───────────────────── padding ───────────────────┐   │  │
│  │  │  ┌────────────────── content ────────────────┐  │   │  │
│  │  │  │                                           │  │   │  │
│  │  │  │       Your text / images live here        │  │   │  │
│  │  │  │                                           │  │   │  │
│  │  │  └───────────────────────────────────────────┘  │   │  │
│  │  └─────────────────────────────────────────────────┘   │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

```css
.box {
  width: 300px;
  padding: 20px;
  border: 2px solid #333;
  margin: 16px;
}
```

## `content-box` vs `border-box`

This is one of the most important concepts in CSS, and one of the first things that trips up every beginner.

### `content-box` (the default)

By default, `width` and `height` apply **only to the content**. Padding and border are *added on top*:

```css
.box {
  box-sizing: content-box; /* default */
  width: 300px;
  padding: 20px;
  border: 2px solid #333;
}
/* Total visible width = 300 + 20 + 20 + 2 + 2 = 344px */
```

You asked for 300px, but the box takes up 344px. This is confusing, error-prone, and the cause of countless layout bugs.

### `border-box` (what you want)

With `border-box`, `width` and `height` include padding and border:

```css
.box {
  box-sizing: border-box;
  width: 300px;
  padding: 20px;
  border: 2px solid #333;
}
/* Total visible width = 300px (content shrinks to 256px to fit) */
```

The box is exactly 300px wide. The content area shrinks to accommodate the padding and border. This matches how humans think about sizing.

:::key
Always set `border-box` globally. This is the universal reset used by every modern CSS framework and professional codebase:
```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```
:::

:::quiz
Q: With `box-sizing: content-box`, if an element has `width: 200px`, `padding: 10px`, and `border: 5px solid`, how wide is the visible box?
- 200px
- 220px
- 230px *
- 250px
E: With `content-box`, total width = content (200) + left padding (10) + right padding (10) + left border (5) + right border (5) = 230px. This is why `border-box` is preferred — it eliminates this math.
:::

## Padding

Padding creates space *inside* the box, between the content and the border. The background color fills the padding area.

```css
/* All four sides */
.card { padding: 24px; }

/* Vertical | Horizontal */
.card { padding: 16px 24px; }

/* Top | Right | Bottom | Left (clockwise) */
.card { padding: 16px 24px 20px 24px; }

/* Individual sides */
.card {
  padding-top: 16px;
  padding-right: 24px;
  padding-bottom: 20px;
  padding-left: 24px;
}
```

:::tip
Remember the clockwise order with **TRouBLe**: **T**op, **R**ight, **B**ottom, **L**eft.
:::

Padding cannot be negative. If you need to "pull" something closer, use negative margins instead.

## Margin

Margin creates space *outside* the box, pushing neighboring elements away. Margins are always transparent — the parent's background shows through.

```css
.card { margin: 16px; }
.card { margin: 0 auto; }  /* center a block element horizontally */
```

### Margin Collapse

This is the box model's most surprising behavior: **vertical margins collapse**. When two block elements stack vertically, their margins don't add up — the larger margin wins.

```css
.heading { margin-bottom: 24px; }
.paragraph { margin-top: 16px; }
/* Gap between them = 24px, NOT 40px */
```

The 24px and 16px margins overlap, and only the larger one (24px) is used.

Margin collapse happens:
- Between adjacent siblings (vertical only).
- Between a parent and its first/last child (if no padding, border, or content separates them).
- On an empty element's own top and bottom margins.

Margin collapse does **not** happen:
- Horizontally (left/right margins never collapse).
- With flexbox or grid children.
- When there's padding, border, or a gap between the elements.
- On floated or absolutely positioned elements.

:::warning
Margin collapse is the #1 reason spacing looks "wrong" to beginners. If you see unexpected gaps (or missing gaps), check for collapsing margins. Opening DevTools and hovering over elements to see their margin highlights is the fastest diagnostic.
:::

### Negative Margins

Margins can be negative, pulling an element in the opposite direction:

```css
.pullquote {
  margin-left: -40px;
  margin-right: -40px;
  padding: 24px 40px;
  background: #f1f5f9;
}
```

This is a classic pattern: a child that stretches *wider* than its parent. Negative margins are legitimate and useful, but use them sparingly — they can cause overlapping content if you're not careful.

### Centering with `margin: 0 auto`

To horizontally center a block element with a defined width:

```css
.container {
  max-width: 960px;
  margin: 0 auto; /* auto distributes remaining horizontal space equally */
}
```

`auto` on left and right margins splits the available space evenly, centering the element. This only works on block elements with an explicit width or max-width.

## Border

Borders sit between padding and margin:

```css
/* Shorthand: width style color */
.card { border: 1px solid #e5e7eb; }

/* Individual sides */
.card { border-bottom: 2px solid #2563eb; }

/* Individual properties */
.card {
  border-width: 1px;
  border-style: solid;
  border-color: #e5e7eb;
}
```

Common border styles: `solid`, `dashed`, `dotted`, `double`, `none`.

## Outline vs Border

`outline` looks like a border but behaves differently:

| Feature          | `border`                    | `outline`                  |
| ---------------- | --------------------------- | -------------------------- |
| Part of box model | Yes (takes up space)       | No (drawn outside the box) |
| Affects layout    | Yes                        | No                         |
| Per-side control  | Yes (`border-left`, etc.)  | No (all sides only)        |
| Can be offset     | No                         | Yes (`outline-offset`)     |
| Main use          | Visual design              | Focus indicators           |

```css
button:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px; /* gap between outline and border */
}
```

:::tip
Use `outline` for focus styles instead of `border`. Since outlines don't affect layout, toggling them on/off won't cause elements to shift position.
:::

## Reading the DevTools Box Model Diagram

Every browser's DevTools shows a visual box model diagram when you select an element. It displays four nested rectangles with pixel values:

- **Blue** (center) = content size
- **Green** = padding
- **Yellow/Orange** = border
- **Orange/Salmon** = margin

Hover over any section and the browser highlights that layer directly on the page. This is your most reliable tool for debugging spacing issues.

To access it: right-click an element → Inspect → look for the box model diagram in the **Computed** or **Layout** panel.

## `display: inline` and the Box Model

Inline elements (`<span>`, `<a>`, `<strong>`) behave differently:

- `width` and `height` are **ignored**.
- Vertical `padding` and `margin` are **applied but don't push surrounding lines** — they overflow visually.
- Horizontal `padding` and `margin` work normally.

If you need an inline element to respect the full box model, use `display: inline-block` — it stays in the text flow but accepts width, height, and vertical padding/margin.

```css
.badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 9999px;
  background: #dbeafe;
  font-size: 0.75rem;
}
```

:::quiz
Q: Why do modern CSS resets set `box-sizing: border-box` on all elements?
- It makes elements render faster
- It includes padding and border in the declared width/height, making sizing predictable *
- It prevents margin collapse
- It automatically centers elements
E: With `border-box`, when you set `width: 300px`, the element is exactly 300px wide including padding and border. Without it, padding and border are added on top, making the element wider than you specified.
:::

:::quiz
Q: Two stacked block elements have `margin-bottom: 30px` and `margin-top: 20px` respectively. What is the gap between them?
- 50px
- 30px *
- 20px
- 10px
E: Vertical margins collapse — the larger margin wins. The gap is 30px, not 50px. This only happens vertically and doesn't apply inside flexbox or grid containers.
:::

## Recap

- Every element is a box: **content → padding → border → margin**.
- Always use `box-sizing: border-box` so `width` means total visible width.
- Vertical margins **collapse** — the larger one wins. This doesn't happen in flex or grid.
- Use `margin: 0 auto` to center block elements with a set width.
- **Outline** is not part of the box model — use it for focus styles.
- Inline elements ignore `width`/`height` — use `inline-block` if you need the full box model.
- The DevTools box model diagram is your best friend for debugging spacing.

**Next up:** Colors, Backgrounds, and Gradients — bringing your boxes to life with color.

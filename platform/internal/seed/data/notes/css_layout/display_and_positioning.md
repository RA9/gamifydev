# Display and Positioning

Every element on a web page occupies a box — but *how* that box behaves, where it sits, and how it interacts with its neighbors all come down to two CSS properties: `display` and `position`. Master these and you'll understand why elements stack, sit side by side, or overlap.

## Display: the element's outer behavior

The `display` property controls how an element participates in the normal document flow. Every HTML element ships with a default, but you can change it.

### Block elements

Block elements start on a **new line** and stretch to fill the full width of their parent. Think `<div>`, `<p>`, `<h1>`–`<h6>`, `<section>`.

```css
.box {
  display: block;
  width: 300px;        /* you CAN set width/height */
  padding: 16px;
  margin-bottom: 12px;
}
```

### Inline elements

Inline elements sit **within a line of text** and only take up as much width as their content. Think `<span>`, `<a>`, `<strong>`, `<em>`.

```css
.tag {
  display: inline;
  /* width and height are IGNORED on inline elements */
  padding: 2px 8px;   /* horizontal padding works; vertical won't push neighbors */
}
```

:::warning
You **cannot** set `width`, `height`, or vertical `margin` on an inline element. If you need those, switch to `inline-block` or `block`.
:::

### Inline-block

The best of both worlds: the element flows **inline** (sits beside other content) but accepts `width`, `height`, and vertical `margin`/`padding` like a block element.

```css
.badge {
  display: inline-block;
  width: 80px;
  height: 28px;
  line-height: 28px;
  text-align: center;
  background: #e0e7ff;
  border-radius: 4px;
}
```

:::key
**Block** = new line, full width. **Inline** = stays in the text flow, no box dimensions. **Inline-block** = stays in the text flow but behaves like a box you can size.
:::

### display: none vs visibility: hidden

Both hide an element, but in very different ways:

```css
.gone    { display: none; }        /* removed from flow — takes no space */
.invisible { visibility: hidden; } /* hidden but still occupies its space */
```

Use `display: none` when you truly want to remove something (a closed modal). Use `visibility: hidden` when you need the element's space preserved (preventing layout shifts).

## Position: breaking out of normal flow

The `position` property changes *where* an element is placed. Once you set it to anything other than `static`, the offset properties `top`, `right`, `bottom`, and `left` become active.

### static (default)

Elements stack in normal document flow. The offset properties have no effect.

```css
.card { position: static; } /* this is what every element already does */
```

### relative

The element stays in normal flow but you can **nudge** it from its original spot. The space it *would have* occupied is preserved.

```css
.nudged {
  position: relative;
  top: 8px;    /* pushes it 8px DOWN from where it was */
  left: 12px;  /* pushes it 12px RIGHT from where it was */
}
```

### absolute

The element is **removed from normal flow** and positioned relative to its nearest **positioned ancestor** (any ancestor with `position` other than `static`). If there is none, it uses the `<html>` element.

```css
.parent {
  position: relative;   /* establishes the reference frame */
}
.tooltip {
  position: absolute;
  top: 100%;            /* just below the parent */
  left: 0;
  background: #333;
  color: #fff;
  padding: 8px 12px;
  border-radius: 4px;
}
```

:::tip
The most common pattern is wrapping the absolute-positioned child in a `position: relative` parent. Without this, the child flies off to the edge of the page — the classic "where did my element go?" bug.
:::

### fixed

Like `absolute`, but the element is positioned relative to the **viewport** and **stays in place when you scroll**.

```css
.sticky-header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background: #fff;
  z-index: 100;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
```

### sticky

A hybrid: the element scrolls normally until it hits a threshold, then **sticks** in place like `fixed`. It un-sticks when its container scrolls past.

```css
.section-heading {
  position: sticky;
  top: 0;               /* sticks once it reaches the top of the viewport */
  background: #fff;
  padding: 12px 0;
  border-bottom: 1px solid #eee;
}
```

:::warning
`position: sticky` requires the element to have a scrolling ancestor *and* a `top`/`bottom` value. It also won't work if any ancestor has `overflow: hidden`. These are the two most common reasons sticky "just doesn't work."
:::

## z-index and stacking contexts

When positioned elements overlap, `z-index` controls which one is on top. Higher values stack above lower ones.

```css
.dropdown {
  position: absolute;
  z-index: 10;
}
.modal-overlay {
  position: fixed;
  z-index: 50;
}
.modal {
  position: fixed;
  z-index: 100;
}
```

### Stacking contexts

A **stacking context** is a self-contained layer. An element with `z-index: 9999` inside a parent that creates a low-priority stacking context can *still* appear behind elements outside that parent.

New stacking contexts are created by:

- `position: relative/absolute/fixed` + a `z-index` value
- `opacity` less than 1
- `transform`, `filter`, or `will-change`

:::key
`z-index` only competes within the same stacking context. If an element is trapped inside a low-priority context, no amount of `z-index` on it will push it above elements in a higher context. When z-index "doesn't work," check the parent tree for stacking context creators.
:::

## Centering techniques

Centering is one of the most common positioning tasks. Here are the reliable methods, ordered by preference.

### Flexbox centering (recommended)

```css
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
```

### Grid centering (one-liner)

```css
.container {
  display: grid;
  place-items: center;
  height: 100vh;
}
```

### Absolute + transform centering

```css
.parent { position: relative; }
.child {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
```

### Horizontal centering for block elements

```css
.centered-block {
  width: 600px;
  margin-inline: auto;   /* auto left and right margins */
}
```

## Practical example: a notification badge

Combining `relative`, `absolute`, and `z-index` to place a badge on an icon:

```html
<div class="icon-wrapper">
  <img src="bell.svg" alt="Notifications" />
  <span class="badge">3</span>
</div>
```

```css
.icon-wrapper {
  position: relative;
  display: inline-block;
}
.badge {
  position: absolute;
  top: -6px;
  right: -6px;
  background: #ef4444;
  color: #fff;
  font-size: 12px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  z-index: 1;
}
```

The wrapper is `relative` so the badge's `absolute` positioning is anchored to the icon rather than flying off to a page corner.

:::quiz
Q: What is the key difference between `display: none` and `visibility: hidden`?
- `display: none` hides the element but it still takes up space
- `visibility: hidden` removes the element from the flow entirely
- `display: none` removes the element from the flow; `visibility: hidden` hides it but preserves its space *
E: `display: none` collapses the element — it's gone from the layout. `visibility: hidden` makes it invisible but other elements still flow around it as if it's there.
:::

:::quiz
Q: An absolutely positioned element looks for its nearest ______ ancestor to position itself against.
- `static`
- positioned (anything other than `static`) *
- `relative` only
E: `position: absolute` uses the nearest ancestor with `position` set to `relative`, `absolute`, `fixed`, or `sticky`. If none exists, it falls back to the root `<html>` element.
:::

:::quiz
Q: You set `z-index: 9999` on an element but it still appears behind another element. What's the most likely cause?
- The browser ignores z-index values above 999
- The element is inside a stacking context with a lower z-index than the other element's context *
- z-index only works on block elements
E: `z-index` only competes within the same stacking context. A child trapped inside a parent that creates a low-priority stacking context can't escape it, no matter how high the child's z-index is.
:::

## Recap

- `display` controls how an element flows: **block** (new line, full width), **inline** (in the text, no sizing), **inline-block** (in the text, accepts sizing).
- `display: none` removes an element from flow; `visibility: hidden` hides it but preserves its space.
- `position: static` is the default. `relative` nudges from the original spot. `absolute` positions against the nearest positioned ancestor. `fixed` sticks to the viewport. `sticky` switches from scrolling to fixed at a threshold.
- `z-index` controls stacking order *within the same stacking context*. Watch for parent elements that silently create new contexts.
- For centering, prefer flexbox (`justify-content: center; align-items: center`) or grid (`place-items: center`).

**Next up:** Mobile-First Workflow — designing for small screens first and scaling up.

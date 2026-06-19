# CSS Basics

You've built structure with HTML. Now you'll make it *beautiful*. **CSS** — Cascading Style Sheets — controls how your HTML looks: colours, fonts, spacing, and layout.

By the end of this lesson you'll understand how a CSS rule works, where styles come from, and the single most important concept for controlling layout: the box model.

## HTML is the skeleton, CSS is the style

If HTML is the skeleton of a page, CSS is its skin, clothes, and posture. The same HTML can look like a serious bank or a playful game — the difference is entirely CSS.

:::analogy
Think of HTML as a plain house — walls, doors, rooms in the right places. CSS is the paint, furniture, and lighting. You can completely redecorate without moving a single wall.
:::

## How a CSS rule works

Every CSS rule has two parts: a **selector** (what to style) and a **declaration block** (how to style it). Each declaration is a `property: value;` pair.

```css
p {
  color: blue;
  font-size: 18px;
}
```

Read it as: "find every `<p>` element, and make its text blue and 18 pixels tall." `color` and `font-size` are *properties*; `blue` and `18px` are their *values*.

## Three ways to apply CSS

```html
<!-- 1. Inline: on a single element (use sparingly) -->
<p style="color: red;">Hi</p>

<!-- 2. Internal: a <style> block in the page <head> -->
<style>
  p { color: red; }
</style>

<!-- 3. External: a separate file (recommended) -->
<link rel="stylesheet" href="main.css" />
```

The **external** stylesheet is the professional choice: one file styles your whole site, and it's cached so pages load fast.

## Selectors: targeting the right elements

You'll use three selectors constantly:

- **Element** — `p { }` styles *every* paragraph.
- **Class** — `.btn { }` styles any element with `class="btn"`. Reusable, your everyday workhorse.
- **ID** — `#header { }` styles the one element with `id="header"`. Unique.

```css
.button {
  background: #6740e8;
  color: white;
  border-radius: 8px;
}
```

:::quiz
Q: You want a style you can reuse on many buttons across the site. Which selector fits best?
- An ID selector like `#button`
- A class selector like `.button` *
- An element selector like `button`
E: Classes are reusable — apply `class="button"` to as many elements as you like. IDs must be unique, and styling every `<button>` element is often too broad.
:::

## The box model — the key to layout

Here's the concept that unlocks CSS layout. **Every element on a page is a rectangular box**, and that box is made of four layers, from the inside out:

![The CSS box model: content, padding, border, margin](/images/lessons/css-box-model.svg)

- **Content** — the text or image itself.
- **Padding** — space *inside* the box, between the content and the border. Think breathing room.
- **Border** — a line around the padding.
- **Margin** — space *outside* the box, pushing other elements away.

:::analogy
Picture a framed photo on a wall. The **photo** is the content. The **mount** around it is the padding. The **frame** is the border. The **gap** to the next picture is the margin. Same four layers, every time.
:::

A common beginner confusion: *padding vs margin*. Padding grows the space **inside** (the background colour fills it); margin adds space **outside** (it's always transparent).

```css
.card {
  padding: 16px;   /* space inside, around the content */
  border: 2px solid #ddd;
  margin: 24px;    /* space outside, between this card and others */
}
```

:::quiz
Q: You want more space *between* the content and the edge of its coloured box. Which property?
- `margin`
- `padding` *
- `border`
E: Padding is the space inside the box, between content and border — and the background colour fills it. Margin would push *other* elements away instead.
:::

## A few properties you'll reach for

- **Colour**: `color`, `background-color`
- **Text**: `font-size`, `font-weight`, `text-align`
- **Spacing**: `margin`, `padding`
- **Layout**: `display`, and the modern power tools **flexbox** (`display: flex`) and **grid** (`display: grid`)

:::key
A CSS rule pairs a **selector** with **declarations**. Prefer reusable **classes** and external stylesheets. And remember the **box model** — content, padding, border, margin — because almost every layout question comes back to it.
:::

## Talk about it

Explain these out loud:

> "What are the two parts of a CSS rule? And what's the difference between padding and margin?"

If the padding-vs-margin answer comes easily, you've grasped the single most useful idea in CSS.

## What's next

You can now structure (HTML) and style (CSS) a page. Next, **JavaScript Basics** adds *behaviour* — making the page respond to the people using it.

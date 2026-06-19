# HTML Basics

**HTML** (HyperText Markup Language) describes the structure of a web page. You write *elements* using *tags*, and the browser turns them into the page you see.

## Tags and elements

A tag looks like `<p>`. Most elements have an opening and a closing tag, with content in between:

- `<p>Hello world</p>` — a paragraph.
- `<h1>Title</h1>` — the biggest heading. There are `<h1>` through `<h6>`.
- `<a href="https://example.com">Visit</a>` — a link.

Some elements are *empty* (no closing tag), like `<br>` for a line break and `<img>` for an image.

## The skeleton of a page

Every HTML document shares the same basic shape:

- `<!DOCTYPE html>` tells the browser this is modern HTML.
- `<html>` wraps the whole document.
- `<head>` holds information *about* the page (title, links to CSS).
- `<body>` holds everything you actually see.

## Common elements you will use a lot

- **Headings**: `<h1>` ... `<h6>`
- **Text**: `<p>`, `<strong>` (bold), `<em>` (italic)
- **Lists**: `<ul>` with `<li>` items for bullets, `<ol>` for numbered
- **Links**: `<a>`
- **Images**: `<img src="cat.png" alt="A cat">`
- **Containers**: `<div>` and `<span>` for grouping

## Attributes

Tags can carry extra information through *attributes*, written as `name="value"`. The `href` in a link and the `src` and `alt` in an image are attributes. The `alt` text is important — it describes images for screen readers and when images fail to load.

## Key takeaways

- HTML is made of elements written with tags.
- A page has a `<head>` (info) and a `<body>` (visible content).
- Attributes add details, and accessible `alt` text matters.

Next up: **CSS Basics**, where you make this structure look good.

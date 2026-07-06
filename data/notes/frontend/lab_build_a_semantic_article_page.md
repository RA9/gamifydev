# Lab: Build a Semantic Article Page

Time to stop reading about good HTML and actually produce some.

:::project
You'll build a multi-section article page from a content brief. The goal is not visual polish yet — it's **structure**: landmarks, heading hierarchy, media, supporting content, and readable document flow.
:::

## The brief

Create a page for an article called **"How Small Teams Ship Faster"**.

Your page should include:

- a site header with the publication name
- a nav with at least three links
- the article itself with:
  - a main heading
  - a short standfirst / intro paragraph
  - at least three subsections
  - one bulleted list
  - one pull-quote or emphasized point
- a figure with an image and caption
- an author section
- an aside with "Related reads"
- a footer

## Your constraints

Build it with **HTML only**.

That means:

- no CSS needed yet
- no JavaScript needed yet
- no generic wrapper soup unless you truly need it

Use the most meaningful elements you can.

## Suggested structure

You do **not** have to copy this exactly, but this is the level of structure you should aim for:

```html
<body>
  <header>
    <h1>Signal Weekly</h1>
    <nav>
      <a href="/">Home</a>
      <a href="/articles">Articles</a>
      <a href="/about">About</a>
    </nav>
  </header>

  <main>
    <article>
      <h1>How Small Teams Ship Faster</h1>
      <p>Intro paragraph...</p>

      <figure>
        <img src="team.jpg" alt="Small product team sketching ideas on a whiteboard" />
        <figcaption>Small teams often communicate faster and make decisions with less friction.</figcaption>
      </figure>

      <section>
        <h2>Clear ownership</h2>
        <p>...</p>
      </section>

      <section>
        <h2>Tight feedback loops</h2>
        <p>...</p>
      </section>
    </article>

    <aside>
      <h2>Related reads</h2>
      <ul>
        <li><a href="/post/1">Design reviews that stay useful</a></li>
      </ul>
    </aside>
  </main>

  <footer>
    <p>© Signal Weekly</p>
  </footer>
</body>
```

## Build order

### 1. Start with the page skeleton

Set up `<!DOCTYPE html>`, `<html lang="en">`, `<head>`, a meaningful `<title>`, and `<body>`.

### 2. Add the landmarks first

Before you write detailed content, create the large page regions:

- `<header>`
- `<main>`
- `<article>`
- `<aside>`
- `<footer>`

This forces you to think structurally from the start.

### 3. Add the article outline

Write the article heading and subsections with real heading levels.

A good outline might be:

- `<h1>` — article title
- `<h2>` — section 1
- `<h2>` — section 2
- `<h2>` — section 3

### 4. Add supporting content

Use:

- `<figure>` + `<figcaption>` for the image
- `<ul>` + `<li>` for the related reads or checklist
- `<strong>` or `<em>` for emphasis

### 5. Review semantics before moving on

Ask:

- Did I use a button where I really needed a link?
- Does every heading make sense in order?
- Is the image alt text useful?
- If all CSS vanished, would the page still read well?

## Acceptance criteria

Your lab is complete when:

- the page has one clear `<h1>` for the article
- major page regions use landmark elements
- navigation uses real links
- the article is divided into meaningful sections
- the image has useful alt text
- the aside contains genuinely secondary content
- the footer contains supporting or ownership info

:::quiz
Q: Which element is the best fit for a box of supporting "Related reads" links beside the article?
- `<main>`
- `<aside>` *
- `<strong>`
E: `<aside>` is designed for supporting or secondary content related to the main content.
:::

## Stretch goals

If you finish early:

- add a newsletter signup form under the article
- add a publication tagline in the header
- add a second figure or a list of takeaways
- rewrite vague links so each one makes sense out of context

## Reflection prompts

When you're done, explain:

- Why is `<article>` a better choice than a generic `<div>` here?
- Why does the heading order matter?
- Why is the related content not part of the main article flow?

:::key
This lab is about building the habit of thinking in **regions, relationships, and meaning**. That's the foundation every later CSS and JavaScript decision depends on.
:::

## What's next

In **CSS Basics**, you'll start making pages like this feel polished — with intentional spacing, typography, and visual hierarchy.
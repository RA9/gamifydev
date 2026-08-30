# Web Accessibility Basics

Accessibility means building websites everyone can use — including people with disabilities. The good news: most of it comes free when you write clean, semantic HTML. In this lesson you'll learn the core habits that make your pages welcoming to all.

## What Accessibility Is (and Who It Helps)

Accessibility — often shortened to **a11y** (an "a", then 11 letters, then a "y") — is the practice of making sure your site works for people with a wide range of abilities.

Who benefits?

- People who are **blind or low-vision** and use screen readers.
- People who **can't use a mouse** and navigate with a keyboard or switch device.
- People who are **deaf or hard of hearing** and need captions.
- People with **color blindness** who can't distinguish certain hues.
- People with **cognitive or attention differences** who need clear, simple structure.

And here's the secret: accessible sites are better for *everyone*. Captions help in a noisy room. Good contrast helps in bright sunlight. Clear link text helps a rushed, distracted user. You are likely one of those users yourself sometimes.

:::analogy
Think of curb cuts — those little ramps where a sidewalk meets the road. They were built for wheelchairs, but they help people with strollers, suitcases, and delivery carts too. Accessible web design is the same: built for some, better for all.
:::

## Semantic HTML Is the Foundation

You already learned semantic elements in HTML Basics. Here's why they matter so much for a11y: assistive technology *reads the meaning* of your tags. A screen reader can jump straight to the `<nav>`, list all the headings, or skip to `<main>` — but only if you used those real elements.

Compare these two buttons:

```html
<!-- Bad: a div pretending to be a button -->
<div class="btn" onclick="save()">Save</div>

<!-- Good: a real button -->
<button onclick="save()">Save</button>
```

The real `<button>` is automatically focusable, works with the Enter and Space keys, and announces itself as "Save, button" to a screen reader. The `<div>` does none of that unless you painstakingly rebuild it.

:::key
The first rule of accessibility: use the correct HTML element for the job. A `<button>` for buttons, `<a>` for links, `<nav>` for navigation. Semantic HTML gives you accessibility almost for free.
:::

## Meaningful Alt Text

Images need an `alt` attribute describing what they show, because screen reader users can't see them.

```html
<!-- Bad: useless or missing -->
<img src="chart.png">
<img src="chart.png" alt="image">

<!-- Good: describes the meaning -->
<img src="chart.png" alt="Bar chart showing sales doubled from 2024 to 2025">
```

Write `alt` text that conveys the *purpose* of the image, not the file name. Ask yourself: "If this image were gone, what words would I put in its place?"

For **purely decorative** images that add no information (a swirl, a divider), use an empty `alt`:

```html
<img src="decorative-swirl.png" alt="">
```

An empty `alt=""` tells the screen reader "skip this — it's decoration." Leaving `alt` out entirely is different and worse: some screen readers will read the file name aloud instead.

:::warning
Don't start alt text with "Image of..." or "Picture of..." — the screen reader already announces it as an image. Just describe what it shows.
:::

## Logical Heading Order

Headings are how screen reader users skim a page. Many pull up a list of all headings to navigate. So your headings must form a sensible outline, in order, without skipping levels.

```html
<!-- Bad: skips from h1 to h4 for styling reasons -->
<h1>Recipes</h1>
<h4>Cookies</h4>

<!-- Good: a clean outline -->
<h1>Recipes</h1>
<h2>Cookies</h2>
<h3>Chocolate Chip</h3>
<h2>Cakes</h2>
```

Use one `<h1>` for the page topic, then `<h2>` for sections, `<h3>` for sub-sections. Need smaller text? Change the size with CSS — never pick a heading level by how big it looks.

## Labels for Form Fields

You met this in the Forms lesson, and it's worth repeating here because it's one of the biggest a11y wins. Every input needs a label connected by `for`/`id`.

```html
<!-- Bad: no label, placeholder only -->
<input type="text" placeholder="Search">

<!-- Good: a real label -->
<label for="search">Search</label>
<input type="text" id="search" name="search">
```

Without a label, a screen reader announces "edit text, blank" and the user has no idea what to type.

## Keyboard Navigation and Visible Focus

Plenty of people never touch a mouse. They press **Tab** to move between interactive elements and **Enter** or **Space** to activate them. Your site must work this way too.

Two rules:

1. **Everything interactive must be reachable by Tab.** Real `<a>`, `<button>`, and form elements get this automatically. Clickable `<div>`s do not.
2. **The focused element must be visibly highlighted** so keyboard users can see where they are.

Browsers add a focus outline by default. A shockingly common mistake is removing it for looks:

```css
/* Bad: blinds keyboard users */
button:focus { outline: none; }

/* Good: provide a clear, custom focus style instead */
button:focus-visible {
  outline: 3px solid #2563eb;
  outline-offset: 2px;
}
```

:::warning
Never write `outline: none` without replacing it with another visible focus style. Removing the focus ring makes your site unusable for keyboard users — they literally cannot see where they are.
:::

:::quiz
Q: Why should you avoid `outline: none` on focusable elements?
- It makes the page load slower
- Keyboard users lose the visible indicator of where they are *
- It breaks the mouse cursor
E: The focus outline shows keyboard users their current position; removing it without a replacement makes navigation impossible for them.
:::

## Color Contrast

Text must stand out clearly from its background. Low-contrast text (light gray on white) is hard to read for low-vision users — and honestly for everyone.

The WCAG guideline asks for a contrast ratio of at least **4.5:1** for normal text (and 3:1 for large text).

```css
/* Bad: light gray on white, hard to read */
.note { color: #bbbbbb; background: #ffffff; }

/* Good: dark text on white passes easily */
.note { color: #1f2937; background: #ffffff; }
```

:::tip
Don't eyeball it — use a contrast checker. Browser DevTools (covered later in the course) shows the contrast ratio right in the color picker, and free online checkers let you paste two colors. Also: never use color *alone* to convey meaning (e.g. "click the green button") — add a label or icon too, for color-blind users.
:::

## Link Text That Makes Sense Out of Context

Screen reader users often pull up a list of all links on a page. If every link says "click here," that list is useless.

```html
<!-- Bad: meaningless out of context -->
<p>To read our policy, <a href="/privacy">click here</a>.</p>

<!-- Good: the link text describes the destination -->
<p>Read our <a href="/privacy">privacy policy</a>.</p>
```

The link text alone should tell you where it goes.

:::quiz
Q: Which link text is best for accessibility?
- "Click here"
- "Read more"
- "Download the 2025 annual report (PDF)" *
E: Descriptive link text makes sense even when read out of context in a list of links.
:::

## ARIA Basics (Use Sparingly)

ARIA — **Accessible Rich Internet Applications** — is a set of extra attributes (like `role`, `aria-label`, `aria-hidden`) that add accessibility information when plain HTML can't. It's powerful, but easy to misuse.

```html
<!-- An icon-only button needs a name for screen readers -->
<button aria-label="Close dialog">✕</button>

<!-- Hide purely decorative content from screen readers -->
<span aria-hidden="true">🎉</span>
```

The single most important ARIA principle:

:::key
**No ARIA is better than bad ARIA.** A native `<button>` is more reliable than a `<div role="button">` patched up with ARIA. Reach for semantic HTML first, and only add ARIA when there is no HTML element that does the job.
:::

A classic mistake is slapping `role="button"` on a `<div>` instead of just using `<button>`. The `<button>` already has the role, keyboard handling, and focus built in — ARIA on a div gives you only the *label*, leaving you to rebuild everything else by hand (and usually getting it wrong).

## How Screen Readers Work (High Level)

A **screen reader** is software that reads the page aloud (or sends it to a Braille display). The user doesn't see the screen — they *hear* it, one element at a time, and navigate with the keyboard.

The screen reader walks through your HTML in order, announcing what each element *is* and *says*: "Heading level 1, Recipes." "Link, privacy policy." "Email, edit text, required." It builds this from your tags, labels, and alt text.

This is why semantic HTML matters so much: the screen reader can only describe what your markup tells it. Good markup produces a clear spoken page; div soup produces confusing noise.

:::analogy
Imagine someone describing a webpage to you over the phone, reading it top to bottom. If the page is well-structured ("Main heading: Recipes. Navigation with three links. Article: Cookies...") you can picture it. If it's just "text, text, link, text, button button button" you're lost. That phone call is roughly what a screen reader experience is like.
:::

## A Quick Self-Audit Checklist

Before you ship, run through this:

- Can I reach and operate everything with only the **Tab** and **Enter** keys?
- Is there a **visible focus** outline as I tab through?
- Does every image have a meaningful `alt` (or `alt=""` if decorative)?
- Does every form field have a connected `<label>`?
- Do my headings go in order (`h1` → `h2` → `h3`) without skipping?
- Is my text **high-contrast** against its background?
- Does every link make sense **read on its own**?
- Did I use **real** `<button>` and `<a>` elements instead of clickable divs?

:::tip
Try unplugging your mouse and using your own site with the keyboard for two minutes. You'll find accessibility problems faster than any tool can describe them.
:::

## Recap

- Accessibility (**a11y**) makes sites usable by people with disabilities — and better for everyone.
- **Semantic HTML is the foundation**: use real `<button>`, `<a>`, `<nav>`, headings, and labels.
- Give images meaningful `alt` text; use empty `alt=""` for decorative ones.
- Keep headings in logical order, label every form field, and keep color contrast high.
- Ensure full **keyboard access** and a **visible focus** style — never `outline: none` with no replacement.
- Write link text that makes sense out of context; use ARIA sparingly — *no ARIA beats bad ARIA*.

**Next up:** CSS Basics — making your accessible, well-structured pages beautiful.

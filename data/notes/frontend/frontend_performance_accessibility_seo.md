# Frontend Performance, Accessibility & SEO

Shipping a page that looks good is only part of frontend work.

Strong frontend developers also care about whether a page is:

- fast
- accessible
- discoverable
- stable under real-world conditions

This lesson is a practical production-readiness pass over everything you've built so far.

## Performance: make the browser do less work

Performance is usually improved by reducing unnecessary cost.

High-value wins:

- compress and resize images
- avoid shipping oversized media
- lazy-load non-critical images
- keep layout and CSS simple where possible
- avoid unnecessary JavaScript for basic UI

Helpful example:

```html
<img src="team-photo.jpg" alt="Product team during a planning session" loading="lazy" />
```

That tells the browser it can delay loading the image until it is needed.

:::quiz
Q: What does `loading="lazy"` help with on a non-critical image?
- It makes the image accessible to screen readers
- It allows the browser to delay loading the image until needed *
- It turns the image into a background
E: Lazy loading helps reduce initial work so the browser can focus on visible content first.
:::

## Stability matters too

A page can feel slow even when downloads are okay if the layout jumps around.

Common causes of layout shift:

- images with no reserved dimensions
- content injected above the current reading position
- late-loading fonts or banners that push content down

A simple habit helps a lot:

- give media predictable sizing behavior
- keep important layout containers stable

## Accessibility is part of product quality

Accessibility is not a bonus round after launch.

Your baseline should include:

- semantic HTML landmarks
- proper heading order
- labels for all form controls
- meaningful alt text for real images
- visible focus states
- enough color contrast
- buttons for actions, links for navigation

These choices improve the experience for everyone, not only users with permanent disabilities.

## Keyboard and focus checks

A fast accessibility review asks:

- Can I tab through the interactive elements?
- Is the current focus clearly visible?
- Can I reach the main actions without a mouse?
- Does the interface still make sense in reading order?

If the answer is "not really," the page is not ready.

## SEO starts with content structure

Search engines are not magic. They rely heavily on clear structure and metadata.

Important basics:

- a useful `<title>`
- a clear `<h1>`
- descriptive headings
- meaningful link text
- semantic structure
- alt text where appropriate

Weak title:

```html
<title>Home</title>
```

Better title:

```html
<title>Sprintboard — Weekly planning for small product teams</title>
```

The better title tells both humans and search engines what the page is about.

## Content clarity improves both UX and SEO

Helpful patterns:

- descriptive section headings
- short readable paragraphs
- links that say where they lead
- real page copy instead of placeholder filler

This is why semantics, readability, and SEO overlap so much.

:::match
Q: Match each concern to a practical fix.
- Slow image-heavy page | Compress and resize images
- Hard-to-use form | Add labels and clear errors
- Weak search snippet | Write a clearer page title
- Keyboard trap | Use real interactive elements and visible focus states
E: Production quality comes from handling specific risks with practical fixes.
:::

## A practical ship checklist

Before calling a page done, review:

### Performance

- Are images appropriately sized?
- Are there any obvious giant assets?
- Is non-critical media lazy-loaded?

### Accessibility

- Is there one logical `<h1>`?
- Do all inputs have labels?
- Is button/link usage correct?
- Can the page be used with keyboard navigation?

### SEO / discoverability

- Does the page title describe the content clearly?
- Do headings describe real sections?
- Do links make sense out of context?

## Lighthouse and DevTools are your allies

Even without deep tooling knowledge, browser tools can help you spot:

- large assets
- accessibility issues
- contrast problems
- missing labels
- layout behavior across devices

The goal is not chasing perfect scores. The goal is catching obvious quality problems before users do.

## Common mistakes to avoid

- uploading raw giant images straight from a phone
- relying only on color to communicate state
- shipping forms with missing labels
- writing page titles that say nothing specific
- assuming desktop-only testing is enough

:::warning
A page can be visually beautiful and still be slow, confusing, or hard to find. Frontend quality is broader than appearance.
:::

## What good looks like

You should now be able to:

- identify easy performance wins
- run a lightweight accessibility review
- improve page titles, headings, and link clarity for SEO
- think about production readiness, not just implementation

## What's next

In **Next Steps**, you'll map out how to keep growing from strong fundamentals into advanced frontend engineering, frameworks, testing, and professional portfolio work.
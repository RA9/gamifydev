# Adding Transitions and Polish

The layout works, it's responsive, and the content is in place. Now you add the finishing touches that make the page feel *alive* — smooth hover effects, scroll-based reveals, and the small details that signal quality.

## Button hover transitions

A button that changes instantly on hover feels jarring. A smooth transition makes it feel intentional.

```css
.button {
  background: var(--color-primary);
  color: #fff;
  padding: 0.8rem 1.6rem;
  border: none;
  border-radius: var(--radius);
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  transition: background 0.2s ease, transform 0.2s ease;
}

.button:hover {
  background: var(--color-primary-dark);
  transform: translateY(-1px);
}

.button:active {
  transform: translateY(0);
}
```

The `translateY(-1px)` gives a subtle lift on hover. The `:active` state brings it back down on click, creating a press effect.

:::tip
Keep transitions between 150ms and 300ms. Faster than 150ms feels abrupt; slower than 300ms feels sluggish. `0.2s ease` is a safe default for most hover effects.
:::

## Subtle card hover lifts

Feature cards and testimonial cards feel more interactive with a gentle lift on hover.

```css
.feature {
  background: var(--color-card);
  padding: var(--space-lg);
  border-radius: var(--radius);
  box-shadow: var(--shadow-card);
  transition: box-shadow 0.25s ease, transform 0.25s ease;
}

.feature:hover {
  box-shadow: var(--shadow-hover);
  transform: translateY(-4px);
}
```

The effect is simple: the shadow deepens and the card rises slightly. It feels like picking up a card off a table.

:::warning
Don't add hover effects to everything. Cards and buttons benefit from them. Body text and headings do not. Overdoing hover effects makes the page feel twitchy, not polished.
:::

## Nav link hover effects

Give navigation links a subtle underline animation instead of the default browser underline.

```css
.nav a {
  text-decoration: none;
  color: var(--color-ink);
  position: relative;
}

.nav a::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--color-primary);
  transition: width 0.25s ease;
}

.nav a:hover::after {
  width: 100%;
}
```

This creates an underline that slides in from the left on hover — much more refined than a sudden text-decoration toggle.

## Scroll-based reveals with IntersectionObserver

Elements that fade in as you scroll down feel dynamic. The `IntersectionObserver` API watches for elements entering the viewport and triggers a class change.

First, add a CSS class for the hidden and revealed states:

```css
.reveal {
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}

.reveal.visible {
  opacity: 1;
  transform: translateY(0);
}
```

Then add the `reveal` class to the elements you want to animate in HTML:

```html
<article class="feature reveal">...</article>
<article class="feature reveal">...</article>
<article class="feature reveal">...</article>
```

Finally, the JavaScript — short and simple:

```js
const reveals = document.querySelectorAll('.reveal');

const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      entry.target.classList.add('visible');
      observer.unobserve(entry.target); // animate once, then stop watching
    }
  });
}, {
  threshold: 0.15  // trigger when 15% of the element is visible
});

reveals.forEach(el => observer.observe(el));
```

:::key
`IntersectionObserver` is far more performant than listening to the `scroll` event. The browser handles the intersection detection natively instead of running your code on every pixel of scroll.
:::

## Smooth scroll for anchor links

When a visitor clicks "Features" in the nav, the page should glide to that section, not jump.

The simplest approach is one CSS declaration:

```css
html {
  scroll-behavior: smooth;
}
```

For more control (like offsetting for a sticky header), use JavaScript:

```js
document.querySelectorAll('a[href^="#"]').forEach(link => {
  link.addEventListener('click', (e) => {
    e.preventDefault();
    const target = document.querySelector(link.getAttribute('href'));
    if (target) {
      target.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  });
});
```

:::tip
If you have a sticky/fixed header, the smooth scroll will land the section title *behind* the header. Add `scroll-margin-top` to your sections to offset for the header height:

```css
section {
  scroll-margin-top: 80px; /* matches your header height */
}
```
:::

## Adding a favicon

A missing favicon makes the browser tab look unfinished. It's the smallest detail with a disproportionate impact.

1. Create a 32×32px `.ico` or `.png` icon (use favicon.io or realfavicongenerator.net)
2. Place it in your project root
3. Add the link in `<head>`:

```html
<link rel="icon" type="image/png" href="favicon.png" />
```

For maximum compatibility across devices:

```html
<link rel="icon" type="image/png" sizes="32x32" href="favicon-32x32.png" />
<link rel="icon" type="image/png" sizes="16x16" href="favicon-16x16.png" />
<link rel="apple-touch-icon" sizes="180x180" href="apple-touch-icon.png" />
```

## Final typography pass

Before calling the page done, review typography one more time:

- **Line length:** Is body text constrained to ~45–75 characters per line? Use `max-width` in `ch` units.
- **Line height:** Body text at `1.6`, headings at `1.2`.
- **Font weight contrast:** Headings should be visibly bolder than body text.
- **Whitespace:** Do headings have more space above them than below? They should — it groups them with the content they introduce.

```css
/* Space above heading > space below heading */
h2 {
  margin-top: var(--space-xl);
  margin-bottom: var(--space-md);
}
```

## Checking color contrast

Every text color + background color pairing on your page must meet WCAG AA:

- **Normal text:** 4.5:1 minimum contrast ratio
- **Large text (18px+ bold or 24px+):** 3:1 minimum

Use the DevTools color picker — it shows the contrast ratio inline. Or use webaim.org/resources/contrastchecker.

```text
✅ #2b2118 on #fdf8f3 → 13.8:1 (excellent)
✅ #b5651d on #ffffff → 4.5:1 (passes AA)
⚠️ #6b5d4f on #fdf8f3 → 3.9:1 (fails for small text — darken it)
```

:::warning
Muted text colors are the most common contrast failure. That light gray you chose for captions may look elegant but be unreadable for people with low vision. Always verify with a contrast checker.
:::

## Focus states for accessibility

Hover effects are great, but keyboard users navigate with Tab. Make sure focused elements are visually obvious.

```css
.button:focus-visible {
  outline: 3px solid var(--color-primary);
  outline-offset: 3px;
}

.nav a:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
```

`:focus-visible` only shows the outline for keyboard navigation, not mouse clicks — the best of both worlds.

:::quiz
Q: Why use `IntersectionObserver` instead of a scroll event listener for reveal animations?
- It supports more browsers
- The browser handles intersection detection natively, making it far more performant *
- It automatically creates animations without CSS
- Scroll events don't work in mobile browsers
E: Scroll event listeners fire on every pixel of scroll, which can cause performance issues (jank). `IntersectionObserver` uses native browser optimization to detect when elements enter the viewport with minimal performance cost.
:::

:::quiz
Q: What is the minimum contrast ratio for normal-sized text under WCAG AA?
- 2:1
- 3:1
- 4.5:1 *
- 7:1
E: WCAG AA requires at least 4.5:1 for normal text. Large text (18px bold or 24px regular) can pass with 3:1. WCAG AAA requires 7:1 for normal text.
:::

## Recap

- **Button transitions** (0.2s ease) with a subtle `translateY` lift feel intentional and polished.
- **Card hover effects** deepen the shadow and lift the card — apply sparingly.
- **Nav link underline animation** using `::after` replaces the default underline with a sliding effect.
- **IntersectionObserver** triggers `.visible` classes as elements scroll into view — far better than scroll listeners.
- **Smooth scroll** with `scroll-behavior: smooth` or `scrollIntoView()` makes anchor links glide.
- **Favicon** and final **typography pass** prevent the "almost done" feeling.
- **Color contrast** must meet 4.5:1 for normal text (WCAG AA).
- **Focus states** with `:focus-visible` ensure keyboard users can navigate.

**Next up:** Code Review and Ship — validating, auditing, and deploying your finished landing page.

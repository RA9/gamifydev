# CSS Transitions

Static interfaces feel lifeless. When a button changes color *instantly* on hover, it feels jarring. CSS transitions let you smoothly animate property changes over a duration you control — turning abrupt flips into polished interactions with a single line of CSS.

## The four transition properties

A transition needs at least two things: *what* to animate and *how long*:

```css
.btn {
  transition-property: background-color;
  transition-duration: 200ms;
  transition-timing-function: ease;
  transition-delay: 0ms;
}
```

### transition-property

Specifies which CSS properties to animate:

```css
.card { transition-property: transform, box-shadow; }  /* specific */
.card { transition-property: all; }                     /* everything */
```

:::tip
Prefer listing specific properties over `all`. It avoids accidentally animating layout-heavy properties and makes your intent clear to other developers.
:::

### transition-duration

```css
.fast   { transition-duration: 150ms; }   /* micro-interactions */
.normal { transition-duration: 300ms; }   /* standard UI speed */
```

:::key
For UI interactions (hover, focus, toggles), **150–300ms** is the sweet spot. Under 100ms feels instant. Over 500ms feels sluggish. Save longer durations for deliberate effects like page transitions.
:::

### transition-timing-function

Controls the acceleration curve:

```css
.linear    { transition-timing-function: linear; }
.ease      { transition-timing-function: ease; }          /* default */
.ease-out  { transition-timing-function: ease-out; }      /* best for UI */
.ease-both { transition-timing-function: ease-in-out; }
.bouncy    { transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1); }
```

:::tip
`ease-out` is the best default for UI transitions. Elements decelerate as they arrive, which feels natural — like sliding a book across a table.
:::

### transition-delay

```css
.delayed { transition-delay: 100ms; }
```

## The shorthand

```
transition: <property> <duration> <timing-function> <delay>;
```

```css
.btn {
  background: #3b82f6;
  transition: background-color 200ms ease-out;
}
.btn:hover { background: #2563eb; }
```

Multiple transitions:

```css
.card {
  transition: transform 200ms ease-out,
              box-shadow 200ms ease-out;
}
```

## Which properties animate well

### Cheap (GPU-composited, 60fps)

- **`transform`** — translate, scale, rotate, skew
- **`opacity`** — fade in/out

### Expensive (trigger layout or paint)

- `width`, `height`, `margin`, `padding` — triggers **layout**
- `background-color`, `color`, `box-shadow` — triggers **paint**

:::key
Whenever possible, animate **`transform`** and **`opacity`**. Instead of animating `width`, use `transform: scaleX()`. Instead of `top`, use `transform: translateY()`. Same visual result, dramatically better performance.
:::

## Common hover and focus effects

### Button hover

```css
.btn {
  background: #3b82f6;
  color: #fff;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 200ms ease-out,
              transform 150ms ease-out;
}
.btn:hover  { background: #2563eb; transform: translateY(-2px); }
.btn:active { transform: translateY(0); }
```

### Card lift

```css
.card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: transform 200ms ease-out,
              box-shadow 200ms ease-out;
}
.card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}
```

### Focus ring for accessibility

```css
.input {
  border: 2px solid #d1d5db;
  outline: none;
  transition: border-color 150ms ease-out,
              box-shadow 150ms ease-out;
}
.input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}
```

:::warning
Never remove the focus outline (`outline: none`) without providing a visible alternative. Keyboard users rely on focus indicators to know where they are. Removing them creates a serious accessibility barrier.
:::

## Transitions on state changes

Transitions work with any state change — classes toggled by JS, `:focus`, `:checked`:

```css
.toggle-track {
  width: 48px;
  height: 28px;
  background: #d1d5db;
  border-radius: 14px;
  transition: background-color 200ms ease;
}
.toggle-input:checked + .toggle-track {
  background: #22c55e;
}
```

## Performance: GPU compositing

When you stick to `transform` and `opacity`, the browser promotes the element to its own GPU layer:

```css
.animated-element {
  will-change: transform;   /* hints the browser to prepare a layer */
}
```

:::warning
Don't put `will-change` on everything — each promoted layer uses GPU memory. Use it only on elements that animate frequently.
:::

## A polished nav link pattern

```css
.nav-link {
  position: relative;
  color: #64748b;
  text-decoration: none;
  transition: color 200ms ease-out;
}
.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: #3b82f6;
  transform: scaleX(0);
  transform-origin: right;
  transition: transform 250ms ease-out;
}
.nav-link:hover { color: #0f172a; }
.nav-link:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}
```

The underline slides in from the left on hover and out to the right on leave — a professional touch using only CSS transitions.

:::quiz
Q: Which two CSS properties are cheapest to animate because they run on the GPU?
- width and height
- transform and opacity *
- color and background-color
E: `transform` and `opacity` are composited by the GPU and don't trigger layout or paint. Animating layout properties like `width` forces expensive recalculations every frame.
:::

:::quiz
Q: What's the recommended duration range for typical UI transitions?
- 10–50ms
- 150–300ms *
- 1–2 seconds
E: 150–300ms is the sweet spot: fast enough to feel responsive, slow enough for the eye to perceive. Under 100ms is effectively instant; over 500ms feels sluggish.
:::

## Recap

- CSS transitions animate property changes using `transition-property`, `transition-duration`, `transition-timing-function`, and `transition-delay`.
- Shorthand: `transition: property duration timing delay`.
- **150–300ms** with **`ease-out`** is the go-to for most UI interactions.
- Animate **`transform`** and **`opacity`** for GPU-composited, smooth performance.
- Use transitions on `:hover`, `:focus`, `:checked`, and JS-toggled classes.
- Never remove focus indicators without providing a visible replacement.

**Next up:** CSS Animations and Keyframes — multi-step animations that go beyond simple A-to-B transitions.

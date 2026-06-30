# CSS Transitions and Animations

Motion makes an interface feel alive and responsive. CSS gives you two tools: *transitions* for smooth changes between two states, and *animations* for richer, multi-step sequences. You will learn both, plus how to keep them fast and accessible.

## Transitions: smoothing a change between two states

By default, a CSS change is instant. A `transition` tells the browser to *animate* the change over time instead of snapping to it.

The shorthand has four parts: which **property** to animate, the **duration**, the **timing-function** (the speed curve), and an optional **delay**.

```css
.button {
  background: #4f46e5;
  transition: background 0.3s ease 0s;
  /*          property  dur  curve  delay */
}
.button:hover {
  background: #4338ca; /* the transition smooths this change */
}
```

You can transition multiple properties at once by separating them with commas:

```css
.card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
```

:::tip
A transition lives on the *base* state, not the hover state. Putting it on `.button` (not `.button:hover`) means the change animates smoothly *both* on hover and when the mouse leaves. If you put it only on `:hover`, leaving the element will snap back instantly.
:::

## What can and cannot animate

Not every property is animatable. The rule of thumb: properties with a numeric or color value can usually transition (width, opacity, color, transform), while properties that flip between discrete keywords cannot.

- **Animatable:** `opacity`, `color`, `background-color`, `transform`, `width`, `height`, `margin`, `padding`, `box-shadow`, `border-radius`.
- **Not animatable:** `display` (e.g. `none` to `block`), `position` keyword changes, `font-family`.

:::warning
A very common bug: trying to fade something in by transitioning `display: none` to `display: block`. It will not work, `display` jumps instantly. Animate `opacity` (and `visibility`) instead, or use a keyframe animation.
:::

## Timing functions (the speed curve)

The timing-function controls how speed changes over the duration:

- `ease` (the default) starts slow, speeds up, ends slow. Feels natural.
- `linear` moves at a constant speed. Good for spinners.
- `ease-in` starts slow; `ease-out` ends slow; `ease-in-out` does both.
- `cubic-bezier(...)` lets you craft a custom curve for full control.

```css
.a { transition: transform 0.4s ease-out; }
.b { transition: transform 0.4s linear; }
.c { transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1); } /* a playful overshoot */
```

You rarely need to hand-write a cubic-bezier early on, but it is good to know it exists for that satisfying "bounce" feel.

## `transform`: move, scale, and rotate

`transform` repositions or reshapes an element *without* affecting the layout around it. The common functions:

```css
.move   { transform: translate(10px, -5px); } /* shift right & up */
.moveY  { transform: translateY(-8px); }      /* shift up */
.bigger { transform: scale(1.1); }            /* 110% size */
.spin   { transform: rotate(45deg); }         /* rotate */
.combo  { transform: translateY(-4px) scale(1.05); } /* combine, space-separated */
```

### Why transforms and opacity are performant

When you animate `width`, `top`, or `margin`, the browser has to recalculate layout and repaint, which can stutter. But `transform` and `opacity` can be handled by the GPU on a separate layer, so they animate at a smooth 60fps even on phones.

:::key
Whenever you have a choice, animate `transform` and `opacity` instead of `width`, `height`, `top`, `left`, or `margin`. This is the single biggest tip for buttery-smooth motion. Use `transform: translateY(-4px)` instead of changing `margin-top`, for example.
:::

## Hover transitions on buttons and cards

A polished button lift:

```css
.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  background: #4f46e5;
  color: #fff;
  cursor: pointer;
  transition: transform 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
}
.btn:hover {
  transform: translateY(-2px);
  background: #4338ca;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}
.btn:active {
  transform: translateY(0); /* press down on click */
}
```

A card that lifts and brightens on hover:

```css
.card {
  border-radius: 12px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}
.card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.18);
}
```

## Keyframe animations

Transitions only go between two states. For multi-step or looping motion, use `@keyframes` to define stages, then attach them with the `animation` property.

```css
@keyframes pulse {
  0%   { transform: scale(1); }
  50%  { transform: scale(1.15); }
  100% { transform: scale(1); }
}
```

You can use `from` and `to` for simple two-step animations:

```css
@keyframes fadeIn {
  from { opacity: 0; }
  to   { opacity: 1; }
}
```

### The `animation` shorthand

The `animation` shorthand bundles several values:

```css
.badge {
  animation: pulse 1.5s ease-in-out infinite;
  /*         name  dur  curve       iteration-count */
}
```

Useful pieces:

- **name** the `@keyframes` name to run.
- **duration** how long one cycle takes (`1.5s`).
- **iteration-count** how many times: a number or `infinite`.
- **direction** `normal`, `reverse`, or `alternate` (play forward then backward, great for smooth loops).
- **fill-mode** `forwards` keeps the final frame after it ends; `backwards` applies the first frame during a delay; `both` does both.

```css
.notice {
  animation: fadeIn 0.6s ease-out forwards; /* stays visible after finishing */
}
```

:::warning
Without `fill-mode: forwards`, an element often snaps back to its starting style the instant the animation ends. If your faded-in element disappears again, add `forwards`.
:::

## Full example: loading spinner

A spinner is just a partial ring rotating forever, the textbook use for `linear` and `infinite`.

```css
.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;          /* light ring */
  border-top-color: #4f46e5;          /* colored arc */
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
```

```html
<div class="spinner" role="status" aria-label="Loading"></div>
```

## Full example: fade-in-up on entry

Content that gently rises into view as it appears, common for cards and headings.

```css
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-up {
  animation: fadeInUp 0.6s ease-out forwards;
}

/* Stagger several items by delaying each one */
.fade-up:nth-child(2) { animation-delay: 0.1s; }
.fade-up:nth-child(3) { animation-delay: 0.2s; }
```

```html
<div class="fade-up">First</div>
<div class="fade-up">Second</div>
<div class="fade-up">Third</div>
```

This uses only `opacity` and `transform`, so it is GPU-friendly and silky smooth.

## Accessibility: `prefers-reduced-motion`

Some people get dizzy or nauseous from motion. Browsers expose a setting for this, and you should respect it. Wrap or override your animations so they calm down when the user has asked for reduced motion.

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

A more targeted approach is to *add* motion only when it is welcome:

```css
.card { transition: none; }

@media (prefers-reduced-motion: no-preference) {
  .card { transition: transform 0.25s ease; }
}
```

:::key
Respecting `prefers-reduced-motion` is not optional polish, it is a core accessibility practice. Essential feedback (like a focus outline) should still work; it is the decorative movement you tone down.
:::

:::quiz
Q: Why should you animate `transform` and `opacity` instead of `width` or `margin`?
- They are shorter to type
- The browser can hand them to the GPU, avoiding layout recalculation for smooth 60fps *
- They are the only animatable properties
E: `transform` and `opacity` can be composited on the GPU without triggering layout/paint, so they stay smooth even on slower devices.
:::

:::fill
Fill in the blank: To keep an animation's final frame applied after it finishes, set animation-fill-mode to ____.
A: forwards
:::

:::predict
You set `transition: background 0.3s ease;` only inside `.button:hover`, not on `.button`. What happens when the mouse leaves the button?
A: The background snaps back instantly with no animation, because the transition rule only exists while hovering; it should live on the base `.button` so both directions animate.
:::

## Recap

- `transition` smoothly animates a change between two states: property, duration, timing-function, delay.
- Numeric and color properties animate; `display` and similar keyword switches do not, so fade with `opacity` instead.
- Timing functions (`ease`, `linear`, `cubic-bezier`) shape the speed curve.
- `transform` (translate/scale/rotate) and `opacity` are the performant properties, prefer them for motion.
- `@keyframes` plus the `animation` shorthand (name, duration, iteration-count, direction, fill-mode) power spinners, fades, and loops.
- Always honor `prefers-reduced-motion` so your animations stay accessible.

**Next up:** putting transitions, animations, Grid, and responsive design together to build a complete, polished interface.

# CSS Animations and Keyframes

Transitions animate simple A-to-B changes. But what about a spinner that rotates forever, a notification that bounces in with multiple stages, or an element that pulses on a loop? CSS animations and `@keyframes` let you define multi-step, looping, self-starting animations entirely in CSS.

## @keyframes: defining the animation

A `@keyframes` rule describes what happens at each stage. You give it a name and define waypoints:

```css
@keyframes fadeIn {
  from { opacity: 0; }
  to   { opacity: 1; }
}
```

`from` equals `0%` and `to` equals `100%`. For multi-step animations, use percentages:

```css
@keyframes bounce {
  0%   { transform: translateY(0); }
  40%  { transform: translateY(-30px); }
  60%  { transform: translateY(-15px); }
  100% { transform: translateY(0); }
}
```

:::key
`@keyframes` only *defines* the animation — it doesn't apply it. You connect it to an element using the `animation` property. Think of `@keyframes` as the choreography and `animation` as saying "that element should dance."
:::

## Animation properties

### animation-name and animation-duration

```css
.element {
  animation-name: fadeIn;
  animation-duration: 500ms;
}
```

### animation-timing-function

Same curves as transitions — `ease`, `linear`, `ease-in-out`, or `cubic-bezier()`. For step-based animations (sprite sheets), use `steps()`:

```css
.typewriter { animation-timing-function: steps(20); }
```

### animation-iteration-count

```css
.spinner   { animation-iteration-count: infinite; }
.attention { animation-iteration-count: 3; }
```

### animation-direction

```css
.pulse { animation-direction: alternate; }
/* normal → 0% to 100% each cycle */
/* alternate → forward, backward, forward... */
```

`alternate` is perfect for pulsing effects — the animation smoothly reverses instead of snapping back.

### animation-fill-mode

Controls styles *before* and *after* the animation:

```css
.element { animation-fill-mode: forwards; }
/* none     → reverts to original styles (default) */
/* forwards → keeps the last keyframe's styles */
/* both     → applies first keyframe during delay, keeps last after */
```

:::tip
If your element snaps back to its original state after animating, you need `animation-fill-mode: forwards`. This tells it to hold the final keyframe's styles.
:::

## The animation shorthand

```
animation: name duration timing-function delay iteration-count direction fill-mode;
```

In practice, specify only what you need:

```css
.spinner { animation: spin 1s linear infinite; }
.toast   { animation: slideUp 300ms ease-out forwards; }
.pulse   { animation: glow 2s ease-in-out infinite alternate; }
```

## Building a spinner

```css
@keyframes spin {
  to { transform: rotate(360deg); }
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
```

`linear` timing is essential — any easing makes the rotation look jerky.

## Entrance animations

### Fade and slide up

```css
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.toast { animation: slideUp 300ms ease-out forwards; }
```

### Staggered entrance

Use `animation-delay` to make items appear sequentially:

```css
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: translateY(0); }
}

.card {
  opacity: 0;
  animation: fadeInUp 400ms ease-out forwards;
}
.card:nth-child(1) { animation-delay: 0ms; }
.card:nth-child(2) { animation-delay: 100ms; }
.card:nth-child(3) { animation-delay: 200ms; }
```

:::tip
For staggered animations, set the base element to `opacity: 0` and use `fill-mode: forwards` so each card stays visible after its animation finishes.
:::

## A pulsing notification dot

```css
@keyframes pulse {
  0%   { transform: scale(1);   opacity: 1; }
  50%  { transform: scale(1.4); opacity: 0.7; }
  100% { transform: scale(1);   opacity: 1; }
}

.notification-dot {
  width: 10px;
  height: 10px;
  background: #ef4444;
  border-radius: 50%;
  animation: pulse 2s ease-in-out infinite;
}
```

## Respecting prefers-reduced-motion

Some users experience motion sickness or discomfort from animations. Always include this:

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

:::warning
This isn't optional polish — it's an accessibility requirement. Users who enable "Reduce motion" in their system settings have a medical or comfort need. Always include a `prefers-reduced-motion` rule in production CSS.
:::

## Animation performance

The same rules from transitions apply, amplified for continuous animations:

- Animate **`transform`** and **`opacity`** — GPU-composited, smooth at 60fps.
- Avoid animating `width`, `height`, `margin`, `top`, `left` — they trigger layout recalculations every frame.
- Use `will-change: transform` on elements that animate continuously (like spinners).

:::quiz
Q: What does `animation-fill-mode: forwards` do?
- Plays the animation in the forward direction
- Makes the element keep the styles from the last keyframe after the animation ends *
- Automatically starts the animation when the element enters the viewport
E: `forwards` tells the element to retain the computed values from the final keyframe after the animation completes, instead of reverting to its original styles.
:::

:::quiz
Q: Why is `linear` timing used for a loading spinner instead of `ease`?
- `linear` is faster
- `linear` provides constant speed so the rotation looks smooth and continuous *
- `ease` doesn't work with `infinite` animations
E: A spinner needs constant rotational speed. `ease` would slow and speed up each revolution, creating a jerky rotation.
:::

:::quiz
Q: What media query disables animations for motion-sensitive users?
- @media (prefers-color-scheme: dark)
- @media (prefers-reduced-motion: reduce) *
- @media (hover: none)
E: `prefers-reduced-motion: reduce` matches users who enabled "Reduce motion" in their OS, letting you disable or simplify animations for accessibility.
:::

## Recap

- `@keyframes` defines animation stages using `from`/`to` or percentage waypoints.
- Connect keyframes with `animation: name duration timing delay count direction fill-mode`.
- Use `infinite` for continuous animations and `alternate` for back-and-forth effects.
- `animation-fill-mode: forwards` keeps the final keyframe styles.
- Stagger entrances with incremental `animation-delay` values.
- Animate `transform` and `opacity` for best performance.
- **Always** include `prefers-reduced-motion: reduce` for accessibility.

**Next up:** CSS Custom Properties — creating reusable, dynamic design tokens with CSS variables.

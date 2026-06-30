# Project: Build a Profile Card

In this guided project you'll build a polished, responsive **profile card** — the kind of component you see on social apps and team pages. You'll practise semantic HTML, CSS variables, flexbox, hover states, and responsive design, all in one self-contained piece.

:::project
**Goal:** Build a profile card with an avatar block, a name and role, a short bio, a row of stats, action buttons, and social links. It should look great on desktop and adapt cleanly to a phone. Everything lives in two files: `index.html` and `style.css`.
:::

## What we're building

The card has, from top to bottom:

- An **avatar block** (a colored banner with a circular avatar).
- A **name** and **role**.
- A short **bio**.
- A **stat row** (posts, followers, following).
- Two **action buttons** (Follow, Message).
- A row of **social links**.

Let's build it in four steps: structure, base styles, hover states, then responsive.

## Step 1: The HTML structure

Start with the standard boilerplate and link the stylesheet. Notice how each visual chunk gets its own clearly named element — this makes the CSS easy to target.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Profile Card</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
  <main class="card">
    <div class="card-banner"></div>

    <div class="card-avatar">JD</div>

    <div class="card-body">
      <h1 class="card-name">Jordan Diaz</h1>
      <p class="card-role">Frontend Developer</p>
      <p class="card-bio">
        I build accessible, fast websites and teach beginners how the web works.
        Coffee enthusiast and weekend hiker.
      </p>

      <div class="card-stats">
        <div class="stat">
          <span class="stat-num">128</span>
          <span class="stat-label">Posts</span>
        </div>
        <div class="stat">
          <span class="stat-num">8.4k</span>
          <span class="stat-label">Followers</span>
        </div>
        <div class="stat">
          <span class="stat-num">312</span>
          <span class="stat-label">Following</span>
        </div>
      </div>

      <div class="card-actions">
        <button class="btn btn-primary">Follow</button>
        <button class="btn btn-ghost">Message</button>
      </div>

      <div class="card-social">
        <a href="#" aria-label="Twitter">Tw</a>
        <a href="#" aria-label="GitHub">Gh</a>
        <a href="#" aria-label="LinkedIn">In</a>
      </div>
    </div>
  </main>
</body>
</html>
```

We're using text initials (`JD`) and short labels (`Tw`, `Gh`, `In`) instead of images, since this lesson uses no images. The `aria-label` attributes keep the social links accessible to screen readers.

:::tip
The `<main>` element marks the primary content of the page. For a single-component demo like this, wrapping the card in `<main>` is a small accessibility win.
:::

## Step 2: Base styles

Set up variables, center the card on the page, and style each block. The card uses `overflow: hidden` so the banner's color stays inside the rounded corners.

```css
:root {
  --brand: #6c5ce7;
  --brand-dark: #5848c2;
  --ink: #1f2233;
  --muted: #7a7f99;
  --bg: #eef0f8;
  --card: #ffffff;
  --line: #e7e9f3;
  --radius: 18px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: system-ui, sans-serif;
  background: var(--bg);
  color: var(--ink);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.card {
  width: 340px;
  max-width: 100%;
  background: var(--card);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: 0 18px 40px rgba(31, 34, 51, 0.12);
}

.card-banner {
  height: 96px;
  background: linear-gradient(135deg, var(--brand), #a29bfe);
}
```

Now the avatar. We pull it *up* over the banner with a negative margin and center it. `border: 4px solid white` gives the classic "punched out" look.

```css
.card-avatar {
  width: 88px;
  height: 88px;
  margin: -44px auto 0;
  border-radius: 50%;
  background: var(--brand-dark);
  color: #fff;
  font-size: 1.6rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4px solid var(--card);
}
```

Then the text body, stats, buttons, and social row.

```css
.card-body {
  padding: 1rem 1.5rem 1.75rem;
  text-align: center;
}

.card-name {
  font-size: 1.4rem;
}

.card-role {
  color: var(--brand);
  font-weight: 600;
  margin-bottom: 0.75rem;
}

.card-bio {
  color: var(--muted);
  font-size: 0.95rem;
  margin-bottom: 1.25rem;
}

.card-stats {
  display: flex;
  justify-content: space-around;
  padding: 1rem 0;
  border-top: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
  margin-bottom: 1.25rem;
}

.stat {
  display: flex;
  flex-direction: column;
}

.stat-num {
  font-size: 1.2rem;
  font-weight: 700;
}

.stat-label {
  font-size: 0.8rem;
  color: var(--muted);
}

.card-actions {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
}

.btn {
  flex: 1;
  padding: 0.7rem;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.18s ease;
}

.btn-primary {
  background: var(--brand);
  color: #fff;
}

.btn-ghost {
  background: transparent;
  color: var(--brand);
  border-color: var(--brand);
}

.card-social {
  display: flex;
  justify-content: center;
  gap: 0.75rem;
}

.card-social a {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: var(--bg);
  color: var(--muted);
  display: flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  font-size: 0.8rem;
  font-weight: 700;
  transition: all 0.18s ease;
}
```

At this point you have a clean, fully styled card. Open it in a browser and admire it.

## Step 3: Hover states

Interactivity makes a card feel alive. We add hover effects to the buttons, the social icons, and the whole card. The `transition` properties we already set make these changes animate smoothly.

```css
.btn-primary:hover {
  background: var(--brand-dark);
  transform: translateY(-1px);
}

.btn-ghost:hover {
  background: var(--brand);
  color: #fff;
}

.card-social a:hover {
  background: var(--brand);
  color: #fff;
  transform: translateY(-2px);
}

.card:hover {
  box-shadow: 0 24px 55px rgba(31, 34, 51, 0.2);
  transform: translateY(-3px);
}

.card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
```

:::key
`transform: translateY(-2px)` nudges an element up slightly, and a `transition` makes the move smooth. This tiny "lift" on hover is one of the most-used tricks for making interfaces feel responsive and polished.
:::

## Step 4: Make it responsive

The card already has `max-width: 100%`, so it shrinks on narrow screens. We just tidy a couple of details for very small phones — reducing side padding and stacking the action buttons.

```css
@media (max-width: 380px) {
  .card-body {
    padding: 1rem 1rem 1.5rem;
  }

  .card-actions {
    flex-direction: column;
  }
}
```

:::quiz
Q: Why does the avatar use `margin: -44px auto 0`?
- To hide the avatar
- To pull it up over the banner (the negative top margin) and center it horizontally (auto sides) *
- To add space below the banner
E: A negative top margin shifts the element upward so it overlaps the banner; `auto` on the left and right margins centers it.
:::

:::predict
You remove `overflow: hidden` from `.card`. What happens to the banner's rounded top corners?
ANSWER: The banner's colored corners would stick out past the card's rounded corners, since nothing is clipping the child to the card's shape anymore.
:::

## Challenge extensions

Now make it yours. Try these, from easiest to hardest:

1. **Theme it.** Change the `--brand` color variables and watch the whole card re-skin instantly. Add a dark-mode version using `@media (prefers-color-scheme: dark)`.
2. **Add a verified badge.** Place a small circular checkmark element next to the name using absolute positioning.
3. **Online status dot.** Add a small green dot on the corner of the avatar to show the user is online.
4. **Real social icons.** Replace the `Tw`/`Gh`/`In` text with inline SVG icons (still no image files).
5. **Toggle Follow.** Add a tiny bit of JavaScript so clicking "Follow" switches the button text to "Following" and changes its style.
6. **Card grid.** Duplicate the card three times and arrange them in a responsive grid that wraps on small screens.

## Recap

- Built a complete profile card with avatar, name, role, bio, stats, actions, and social links.
- Used **CSS variables** so colors and radius can be changed in one place.
- Centered the card with **flexbox** on `body`, and used `overflow: hidden` to clip the banner.
- Overlapped the avatar with a **negative margin** and centered it with `auto`.
- Added smooth **hover states** using `transition` and `transform: translateY`.
- Made it **responsive** with `max-width: 100%` and a small-screen media query.

**Next up:** Project: Build a Quiz Game — add JavaScript to make a fully interactive, working app.

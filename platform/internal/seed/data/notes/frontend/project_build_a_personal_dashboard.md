# Project: Build a Personal Dashboard

This project pulls together modern frontend skills into one richer build.

:::project
**Goal:** Build a personal dashboard that combines UI state, async data loading, local storage, and component layout. Your dashboard should greet the user, show a focus area, fetch at least one piece of live data, and remember a preference like theme or name.
:::

## What the dashboard should include

A good starter dashboard can contain:

- a greeting section
- a focus card with today's main task
- a quote or weather widget loaded with fetch
- a theme toggle
- a name or preference saved in localStorage

You do **not** need a complex backend. This is about frontend integration.

## Step 1 — Build the layout

Create a simple dashboard shell with a header and a small grid of cards.

```html
<main class="dashboard-shell">
  <header class="dashboard-head">
    <h1 id="greeting">Welcome back</h1>
    <button id="theme-btn">Toggle theme</button>
  </header>

  <section class="dashboard-grid">
    <article class="card" id="focus-card">
      <h2>Today's focus</h2>
      <p id="focus-text">Finish the responsive landing page</p>
    </article>

    <article class="card" id="quote-card">
      <h2>Quote</h2>
      <p id="quote-text">Loading quote…</p>
    </article>
  </section>
</main>
```

## Step 2 — Style it like a real product surface

Use Grid for the dashboard and a class on `body` for theme switching.

```css
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  background: #f8fafc;
  color: #0f172a;
}

.dashboard-shell {
  width: min(100% - 2rem, 72rem);
  margin: 0 auto;
  padding: 2rem 0 4rem;
}

.dashboard-grid {
  display: grid;
  gap: 1rem;
}

.card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 18px;
  padding: 1.25rem;
}

body.dark {
  background: #0f172a;
  color: #e2e8f0;
}
```

Then add responsive columns at a wider breakpoint.

## Step 3 — Load saved preferences

Use localStorage for small personal settings.

Examples:

- saved display name
- saved focus task
- saved theme choice

```js
const savedTheme = localStorage.getItem("dashboard-theme");
if (savedTheme === "dark") {
  document.body.classList.add("dark");
}
```

This makes the dashboard feel personal immediately.

## Step 4 — Add a theme toggle

When the button is clicked:

- toggle the `dark` class on `body`
- save the chosen theme to localStorage

This is a clean example of UI state + persistence.

## Step 5 — Fetch live data

Use `fetch()` to load one piece of remote data.

A good simple example is a random quote API.

```js
async function loadQuote() {
  try {
    const response = await fetch("https://dummyjson.com/quotes/random");
    const data = await response.json();
    quoteText.textContent = data.quote;
  } catch (error) {
    quoteText.textContent = "Could not load a quote right now.";
  }
}
```

This teaches a real production lesson: async work needs both success and failure states.

## Step 6 — Add an editable focus card

Let the user update today's focus and save it locally.

One approach:

- render the current focus text
- include an edit input and Save button
- update the stored value on save
- re-render the visible card text

Now the dashboard combines:

- local state
- persistent state
- remote state

That is much closer to a real frontend application.

## Step 7 — Review the product quality

Check that the app:

- works on phone and desktop
- survives a refresh with saved preferences
- handles API failure gracefully
- keeps the layout stable while loading
- has clear interactive controls

:::quiz
Q: Which combination makes a personal dashboard a richer frontend project than a simple static page?
- Only typography and colors
- Persistent preferences, fetched data, and interactive UI state *
- One heading and no buttons
E: A dashboard becomes more realistic when it combines local persistence, async loading, and multiple interactive UI states.
:::

## Stretch goals

- add a task counter
- save the user's name and personalize the greeting
- add a clock or date widget
- add a second API card
- animate theme transitions carefully

## What "done" looks like

A finished version should feel like a small personal web app — not just a styled page. It should remember something, load something, and respond smoothly to the user.
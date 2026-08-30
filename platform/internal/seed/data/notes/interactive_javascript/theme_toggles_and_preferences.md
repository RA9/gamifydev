# Theme Toggles and Preferences

A theme toggle is one of the simplest features that touches every part of state management: reading a preference, applying it to the UI, saving it, and respecting the user's operating system setting. This lesson builds a professional-grade theme system step by step.

## The basic toggle — classList on body

The simplest approach: define dark styles under a class and toggle it.

```css
body {
  background: #ffffff;
  color: #111111;
}

body.dark {
  background: #111111;
  color: #eeeeee;
}
```

```js
const themeBtn = document.querySelector("#theme-btn");

themeBtn.addEventListener("click", () => {
  document.body.classList.toggle("dark");
});
```

One click, the page flips. But this has no memory — refresh and it resets.

## Saving the preference with localStorage

Store the choice so it survives across page loads:

```js
const themeBtn = document.querySelector("#theme-btn");

// Load saved theme on startup
const savedTheme = localStorage.getItem("theme");
if (savedTheme === "dark") {
  document.body.classList.add("dark");
}

themeBtn.addEventListener("click", () => {
  const isDark = document.body.classList.toggle("dark");
  localStorage.setItem("theme", isDark ? "dark" : "light");
});
```

Now the theme persists. But what about users who have *never* clicked the button? Their operating system has a preference too.

## Detecting system preference with matchMedia

The CSS media query `prefers-color-scheme` tells you what the OS is set to. JavaScript can read it with `matchMedia`:

```js
const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");
console.log(prefersDark.matches); // true if the OS is in dark mode
```

You can also listen for changes — the user might switch their OS theme while your page is open:

```js
prefersDark.addEventListener("change", (event) => {
  console.log("System changed to:", event.matches ? "dark" : "light");
});
```

## The system + manual override pattern

The professional pattern supports three states:
1. **No manual override** — follow the system preference.
2. **User chose light** — always light, ignoring the system.
3. **User chose dark** — always dark, ignoring the system.

```js
const themeBtn = document.querySelector("#theme-btn");
const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");

function getEffectiveTheme() {
  const saved = localStorage.getItem("theme");
  if (saved === "dark" || saved === "light") return saved;
  // No manual choice — follow the system
  return prefersDark.matches ? "dark" : "light";
}

function applyTheme() {
  const theme = getEffectiveTheme();
  document.body.classList.toggle("dark", theme === "dark");
}

// Apply on load
applyTheme();

// React to system changes (only matters when there is no manual override)
prefersDark.addEventListener("change", () => {
  if (!localStorage.getItem("theme")) {
    applyTheme();
  }
});

// Manual toggle
themeBtn.addEventListener("click", () => {
  const current = getEffectiveTheme();
  const next = current === "dark" ? "light" : "dark";
  localStorage.setItem("theme", next);
  applyTheme();
});
```

:::key
The priority chain is: **manual override (localStorage) → system preference (matchMedia) → default**. This gives users control while respecting their OS settings when they have not chosen.
:::

## Updating the button label

A toggle button should tell the user what it will *do*, not what the current state is:

```js
function applyTheme() {
  const theme = getEffectiveTheme();
  document.body.classList.toggle("dark", theme === "dark");

  // Update button to show the opposite action
  themeBtn.textContent = theme === "dark" ? "Switch to Light" : "Switch to Dark";
  themeBtn.setAttribute("aria-label",
    theme === "dark" ? "Switch to light theme" : "Switch to dark theme"
  );
}
```

## Adding a "reset to system" option

Some interfaces offer three choices: light, dark, and auto (follow system). A select or a triple-click cycle:

```js
function cycleTheme() {
  const saved = localStorage.getItem("theme");

  if (!saved) {
    // auto → light
    localStorage.setItem("theme", "light");
  } else if (saved === "light") {
    // light → dark
    localStorage.setItem("theme", "dark");
  } else {
    // dark → auto (remove override)
    localStorage.removeItem("theme");
  }

  applyTheme();
}

themeBtn.addEventListener("click", cycleTheme);
```

:::tip
When offering a "system/auto" option, remove the localStorage key entirely rather than storing `"auto"`. This way, `getEffectiveTheme` naturally falls through to the matchMedia check.
:::

## Smooth transitions on theme change

Without transitions, the theme switch is an abrupt flash. Add a CSS transition to body:

```css
body {
  background: #ffffff;
  color: #111111;
  transition: background-color 0.3s ease, color 0.3s ease;
}

body.dark {
  background: #111111;
  color: #eeeeee;
}
```

:::warning
The transition on body should only apply *after* the initial page load. Otherwise, users see a flash of the wrong theme transitioning to the right one on every page load. One way to prevent this:

```js
// Add a class that enables transitions only after the first render
window.addEventListener("load", () => {
  document.body.classList.add("transitions-ready");
});
```

```css
body.transitions-ready {
  transition: background-color 0.3s ease, color 0.3s ease;
}
```
:::

## Using a data attribute instead of a class

An alternative to `body.dark` is a `data-theme` attribute. This makes it easy to support more than two themes:

```css
[data-theme="dark"] {
  --bg: #111;
  --text: #eee;
}

[data-theme="light"] {
  --bg: #fff;
  --text: #111;
}

body {
  background: var(--bg);
  color: var(--text);
}
```

```js
function applyTheme() {
  const theme = getEffectiveTheme();
  document.documentElement.setAttribute("data-theme", theme);
}
```

This pattern scales to any number of themes (`"solarized"`, `"ocean"`, etc.) and plays well with CSS custom properties.

## Preventing the flash of wrong theme

There is a subtle problem: if your script is in the `<body>` or deferred, the page briefly renders with the wrong theme before JavaScript runs. To fix this, put a tiny inline script in the `<head>`:

```html
<head>
  <script>
    // Runs before the page renders — no flash
    const theme = localStorage.getItem("theme");
    if (theme === "dark" || (!theme && matchMedia("(prefers-color-scheme: dark)").matches)) {
      document.documentElement.classList.add("dark");
    }
  </script>
</head>
```

This blocking script runs before the browser paints, so the page starts in the correct theme.

## Practice

:::quiz
Q: A user has never clicked the theme toggle. Their OS is set to dark mode. What should `getEffectiveTheme()` return?
- "light" (the default)
- "dark" (matching the system) *
- null
- "auto"
E: When there is no manual override in localStorage, the function should fall through to the system preference via `matchMedia("(prefers-color-scheme: dark)")`. Since the OS is in dark mode, it returns "dark".
:::

:::quiz
Q: Why should you avoid putting the theme transition CSS on page load?
- Transitions slow down the page
- Users see the wrong theme flash-transitioning to the correct one on every page load *
- CSS transitions do not work on body elements
- localStorage cannot be read fast enough
E: If transitions are active during the initial theme application, the page visibly animates from the default theme to the saved theme. Adding the transition class only after the first paint prevents this flash.
:::

## Recap

- Toggle themes with **`classList.toggle("dark")`** on `body` or a `data-theme` attribute on `<html>`.
- **Save** the preference to `localStorage`; **load** it on startup.
- Detect the system preference with **`matchMedia("(prefers-color-scheme: dark)")`** and listen for live changes.
- Priority: **manual override → system preference → default**.
- Update the button label to show the **opposite** action (what the click will do).
- Add **CSS transitions** but only after initial load to avoid a theme flash.
- For instant, flash-free theme loading, put a small inline script in `<head>`.

**Next up:** Tab Interfaces — managing which content panel is visible.

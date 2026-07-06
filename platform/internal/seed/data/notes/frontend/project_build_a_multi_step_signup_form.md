# Project: Build a Multi-Step Signup Form

Real products often break long forms into smaller, calmer steps. This project teaches you how to manage that flow with HTML, CSS, and JavaScript.

:::project
**Goal:** Build a multi-step sign-up form with progress, next/back buttons, validation at each step, and a final success screen. You will practise state, rendering, validation, and form UX all at once.
:::

## What the form should do

Your form will have three steps:

1. **Account** — name + email
2. **Security** — password + confirm password
3. **Profile** — role + short bio

Then show a final confirmation screen.

## Step 1 — Build the HTML structure

Start with a form shell and separate step panels.

```html
<main class="signup-shell">
  <section class="signup-card">
    <p id="progress-text">Step 1 of 3</p>
    <div class="progress-bar"><span id="progress-fill"></span></div>

    <form id="signup-form">
      <div class="step" data-step="0">
        <h1>Create your account</h1>
        <label for="name">Full name</label>
        <input id="name" type="text" />
        <p class="error" id="name-error"></p>

        <label for="email">Email</label>
        <input id="email" type="email" />
        <p class="error" id="email-error"></p>
      </div>

      <div class="step hidden" data-step="1">
        <h1>Choose a password</h1>
        <label for="password">Password</label>
        <input id="password" type="password" />
        <p class="error" id="password-error"></p>

        <label for="confirm-password">Confirm password</label>
        <input id="confirm-password" type="password" />
        <p class="error" id="confirm-password-error"></p>
      </div>

      <div class="step hidden" data-step="2">
        <h1>Tell us about you</h1>
        <label for="role">Role</label>
        <input id="role" type="text" />
        <label for="bio">Short bio</label>
        <textarea id="bio" rows="4"></textarea>
      </div>

      <div class="actions">
        <button type="button" id="back-btn">Back</button>
        <button type="button" id="next-btn">Next</button>
        <button type="submit" id="submit-btn" class="hidden">Create account</button>
      </div>
    </form>

    <div id="success" class="hidden">
      <h1>You're in 🎉</h1>
      <p>Your account has been created.</p>
    </div>
  </section>
</main>
```

## Step 2 — Style the card and progress UI

Focus on calm form UX: spacing, readable labels, and a visible progress bar.

```css
body {
  margin: 0;
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f8fafc;
  font-family: system-ui, sans-serif;
}

.signup-card {
  width: min(100% - 2rem, 34rem);
  background: white;
  border-radius: 18px;
  border: 1px solid #e2e8f0;
  padding: 2rem;
  box-shadow: 0 14px 36px rgba(15, 23, 42, 0.08);
}

.progress-bar {
  height: 10px;
  background: #e2e8f0;
  border-radius: 999px;
  overflow: hidden;
  margin-bottom: 1rem;
}

#progress-fill {
  display: block;
  height: 100%;
  width: 33.33%;
  background: #5b3df5;
}
```

## Step 3 — Create your form state

Your JavaScript should track both values and which step is active.

```js
const state = {
  step: 0,
  values: {
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
    role: "",
    bio: ""
  },
  errors: {}
};
```

Keeping everything in one object makes the UI much easier to reason about.

## Step 4 — Render the active step

Write a `render()` function that:

- shows only the current step
- updates the progress text and fill width
- shows Back / Next / Submit at the right times
- writes any current error messages into the page

That pattern turns a messy form flow into a predictable one.

## Step 5 — Validate each step before advancing

Use step-specific rules.

Examples:

- step 1: name required, email must include `@`
- step 2: password length >= 8, confirm password must match
- step 3: optional lighter validation

Only move forward if the current step is valid.

:::tip
Validate the *current* step instead of validating the entire form every time. That's what makes multi-step forms feel calm instead of overwhelming.
:::

## Step 6 — Submit and show success

When the final step is valid:

- prevent the default form submission
- hide the form
- show the success state

You can simulate a real submission first, then later swap it for an API call if you want.

## Step 7 — Polish the experience

A stronger final version should include:

- visible focus states
- disabled buttons when appropriate
- consistent spacing between fields
- no jumping layout when errors appear

:::quiz
Q: Why are multi-step forms often easier to use than one giant form?
- They remove the need for validation
- They reduce cognitive load by revealing smaller chunks of work *
- They work without JavaScript
E: Breaking a flow into steps makes long forms feel more manageable and helps people focus on one decision at a time.
:::

## Stretch goals

- add live validation while typing
- show a step summary before submission
- save draft progress to localStorage
- add a password visibility toggle

## What "done" looks like

A finished version should feel like a real product flow: clear progress, helpful validation, smooth navigation, and a satisfying completion state.
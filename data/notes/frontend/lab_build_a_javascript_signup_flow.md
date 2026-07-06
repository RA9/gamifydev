# Lab: Build a JavaScript Signup Flow

This lab turns form theory into a real interactive flow.

:::project
You'll build a signup experience that validates user input, shows inline errors, disables duplicate submission, and reveals a success message when the form is valid.
:::

## The goal

Build a signup form with:

- full name field
- email field
- password field
- submit button
- inline error messages
- success feedback

Your JavaScript should manage the behavior cleanly instead of scattering random DOM changes everywhere.

## What you're practicing

This lab combines:

- DOM selection
- form submit handling
- validation
- state
- conditional rendering
- feedback states

## Suggested state model

Start with a small state object.

```js
const state = {
  values: {
    name: "",
    email: "",
    password: ""
  },
  errors: {
    name: "",
    email: "",
    password: ""
  },
  isSubmitting: false,
  successMessage: ""
};
```

This makes it much easier to reason about the UI.

## Required behavior

### 1. Prevent default submission

When the user submits the form, your JavaScript should intercept it.

```js
form.addEventListener("submit", function (event) {
  event.preventDefault();
});
```

### 2. Validate the fields

Use rules like:

- name cannot be empty
- email must include `@`
- password must be at least 8 characters

### 3. Show errors inline

Each field should have a visible place where errors can appear.

Example:

```html
<p id="emailError" class="error"></p>
```

### 4. Show success only when valid

If there are no validation errors:

- clear old errors
- show a success message
- optionally reset the form values

### 5. Reflect submission state in the button

While submitting, the button might:

- be disabled
- show `Creating account...`

Even if the submit is simulated, modeling this state is great practice.

## A render-based approach

Try centralizing UI updates into a `render()` function.

```js
function render() {
  nameError.textContent = state.errors.name;
  emailError.textContent = state.errors.email;
  passwordError.textContent = state.errors.password;
  submitBtn.disabled = state.isSubmitting;
  successBox.textContent = state.successMessage;
}
```

Then your submit handler can:

- update the state
- call `render()`

That keeps logic easier to follow.

## Acceptance criteria

Your lab is complete when:

- invalid fields show clear inline messages
- valid input removes those messages
- the form does not reload the page on submit
- success feedback appears after a valid submission
- submit state is reflected in the button or UI

## Stretch goals

If you finish early:

- validate live while the user types
- add password visibility toggle
- add a character counter for the name field
- save draft values to `localStorage`

:::quiz
Q: In this lab, what is the biggest advantage of keeping errors in a state object?
- It makes CSS unnecessary
- It gives you one predictable place to derive the UI from *
- It automatically submits the form
E: When errors live in state, the UI can render from one reliable source of truth instead of scattered conditional DOM edits.
:::

## Reflection prompts

When you're done, explain:

- What values belonged in state?
- What parts of the UI were derived from that state?
- What would get harder if you skipped the state object and edited the DOM ad hoc?

:::key
This lab is not just about forms. It's about learning to model interactive UI as data plus rendering — one of the most important frontend habits you can build.
:::

## What's next

In **Building Interactive JavaScript Websites**, you'll use this same mindset across richer mini-app interactions and multi-step UI behavior.
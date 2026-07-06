# Forms, State, and UI Patterns

Once you can manipulate the DOM, the next big leap is learning to structure UI around **state** instead of scattered DOM updates.

This is the difference between code that *works once* and code that scales.

## What state means in frontend work

State is the data that represents what your UI currently knows.

Examples:

- the current value of an input
- whether a modal is open
- which tab is active
- which tasks are completed
- whether a form is loading, succeeded, or failed

```js
const formState = {
  email: "",
  password: "",
  isSubmitting: false,
  errors: {}
};
```

When the state changes, the UI should reflect it.

## UI should be derived from state

A beginner approach often looks like this:

- update one DOM node here
- hide another there
- set an inline style somewhere else
- forget what the real source of truth is

A stronger approach is:

- keep the important data in one place
- render the UI from that data
- update the data when the user acts

:::key
State is the **source of truth**. The DOM is the visible output of that truth.
:::

## A simple validation flow

Consider a signup form.

State might include:

- field values
- field errors
- submission status

```js
const state = {
  email: "",
  password: "",
  errors: {
    email: "",
    password: ""
  },
  isSubmitting: false
};
```

When the user types, update the values.
When they submit, validate.
When validation fails, update the errors.
When validation passes, start submission.

## Common UI states every app needs

Many interfaces need more than one visible mode.

For example:

- **empty** — no tasks yet
- **loading** — data is being fetched or submitted
- **error** — something went wrong
- **success** — action completed

A polished frontend handles these states deliberately instead of pretending everything always works perfectly.

## Example: disabling submit until valid

You often want the UI to respond to state.

```js
const isValid = email.includes("@") && password.length >= 8;
submitButton.disabled = !isValid;
```

The button is not deciding anything by itself. It is reflecting the current validity state.

:::quiz
Q: In a well-structured UI, what should usually be the source of truth?
- Random DOM text values
- Your application state object or variables *
- CSS selectors
E: UI should be derived from state. The DOM is the rendered result, not the primary source of truth.
:::

## Validation should help, not punish

Good form feedback is:

- specific
- timely
- connected to the correct field
- easy to understand

Weak error:

- `Invalid input`

Better:

- `Password must be at least 8 characters`
- `Please enter a valid email address`

## Pattern: render after updates

A small but powerful pattern:

```js
function render() {
  emailInput.value = state.email;
  passwordInput.value = state.password;
  emailError.textContent = state.errors.email;
  passwordError.textContent = state.errors.password;
  submitButton.disabled = state.isSubmitting;
}
```

Then:

- update `state`
- call `render()`

This keeps UI changes predictable.

## Pattern: event handlers only update state

Instead of putting all visual decisions directly into handlers, let handlers update state.

```js
emailInput.addEventListener("input", function (event) {
  state.email = event.target.value;
  render();
});
```

That keeps responsibilities clear:

- handler updates data
- render updates UI

## Useful UI patterns to practice

### 1. Character counter

- state: current input value
- derived UI: remaining characters count

### 2. Password visibility toggle

- state: `showPassword`
- derived UI: input `type` and button label

### 3. Filter buttons

- state: active filter
- derived UI: highlighted button + filtered list

### 4. Toast or success message

- state: success or error outcome
- derived UI: visible feedback banner

## Predict the UI change

:::predict
```js
const state = { open: false };

button.addEventListener("click", function () {
  state.open = !state.open;
  panel.classList.toggle("hidden", !state.open);
});
```
Q: What happens each time the button is clicked?
- The panel alternates between hidden and visible *
- The panel is permanently removed
- The button stops working after one click
E: Each click flips the boolean state and the `hidden` class reflects that state.
:::

## Mini practice — upgrade a signup form

Take a simple HTML form and add these behaviors:

- show an inline message when email is invalid
- disable the submit button while "submitting"
- show a success message after submission
- clear the form only after a successful result

Even if the submission is fake, these are real UI patterns you'll use everywhere.

## Local persistence is also state management

When you save data to `localStorage`, you're preserving state between page loads.

```js
localStorage.setItem("draftEmail", state.email);
```

That same mindset later scales to larger app state systems and frameworks.

## Common mistakes to avoid

- treating the DOM as the main data store
- duplicating the same data in multiple variables with no plan
- showing generic error messages
- mixing validation, rendering, and submission logic in one huge function
- forgetting that loading and success states are part of the product experience

:::warning
If your UI logic feels hard to reason about, it's often because the state model is unclear — not because the DOM is hard.
:::

## What good looks like

You should now be able to:

- explain what state means in frontend work
- separate state updates from rendering
- model common form and UI states clearly
- derive disabled, error, success, and visibility states from data
- think in reusable interaction patterns instead of one-off hacks

## What's next

In **Building Interactive JavaScript Websites**, you'll combine this thinking with DOM work to build more complete browser experiences.
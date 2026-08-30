# Form Validation and Best Practices

You learned how to build forms in a previous lesson. Now it's time to make them bulletproof. HTML gives you a powerful validation system that works without a single line of JavaScript — and when you need more control, the Constraint Validation API lets you customize everything. In this lesson you'll master both.

## Why Validate on the Client?

Client-side validation gives users **instant feedback** before they hit submit. No round trip to the server, no page reload, no waiting. A required field glows red the moment they try to skip it, and a malformed email gets caught before they wonder why the form didn't work.

:::warning
Client-side validation is for **user experience**, not security. A determined user can bypass any browser validation (DevTools, curl, disabled JavaScript). Always validate again on the server. Think of HTML validation as a helpful guardrail, not a locked gate.
:::

## The `required` Attribute

The simplest validation: the field cannot be empty.

```html
<label for="name">Full name</label>
<input type="text" id="name" name="name" required>
```

If the user tries to submit with this field blank, the browser blocks submission and shows a tooltip like "Please fill out this field." Works on `<input>`, `<textarea>`, and `<select>`.

## Text Length: `minlength` and `maxlength`

Limit how many characters a text field accepts.

```html
<label for="username">Username (3–20 characters)</label>
<input type="text" id="username" name="username"
       required minlength="3" maxlength="20">
```

- **`minlength`** — the field is invalid if the user types fewer characters than this.
- **`maxlength`** — the browser physically prevents typing past this limit.

Note the behavioral difference: `minlength` shows an error on submit, but `maxlength` silently stops input. Some developers pair `maxlength` with a visible character counter so users aren't surprised when their text stops appearing.

## Numeric Limits: `min` and `max`

For `type="number"`, `type="date"`, and `type="range"`, you can set valid boundaries.

```html
<label for="age">Age (must be 13 or older)</label>
<input type="number" id="age" name="age" min="13" max="120">

<label for="event-date">Event date</label>
<input type="date" id="event-date" name="event_date"
       min="2025-01-01" max="2025-12-31">
```

The `step` attribute controls increments:

```html
<!-- Only accept multiples of 5 -->
<input type="number" name="quantity" min="0" max="100" step="5">
```

## Pattern Matching

The `pattern` attribute takes a regular expression that the value must match. Always pair it with a `title` that explains the expected format in plain language.

```html
<label for="zip">ZIP Code</label>
<input type="text" id="zip" name="zip"
       pattern="[0-9]{5}"
       title="Five-digit ZIP code"
       required>

<label for="slug">URL Slug</label>
<input type="text" id="slug" name="slug"
       pattern="[a-z0-9]+(-[a-z0-9]+)*"
       title="Lowercase letters, numbers, and hyphens only">
```

The browser shows the `title` text alongside its validation message, giving the user actionable guidance.

:::tip
The `pattern` regex is automatically anchored — the browser tests the entire value against the pattern as if you wrote `^pattern$`. You don't need to add `^` and `$` yourself.
:::

## Type-Based Validation

Several input types have built-in validation rules:

```html
<!-- Must be a valid email format -->
<input type="email" name="email" required>

<!-- Must be a valid URL -->
<input type="url" name="website" placeholder="https://example.com">

<!-- Brings up a numeric keypad on mobile but no built-in format check -->
<input type="tel" name="phone"
       pattern="[+]?[0-9\s\-()]+"
       title="Phone number with optional country code">
```

- **`type="email"`** checks for `something@something`. It also accepts multiple comma-separated emails if you add `multiple`.
- **`type="url"`** requires a full URL with protocol (`https://…`).
- **`type="tel"`** does **not** validate the format automatically (phone formats vary globally). Pair it with `pattern` for format enforcement.

:::key
`type="email"` and `type="url"` give you free format validation. `type="tel"` only changes the keyboard — it doesn't validate the format. Always add a `pattern` for phone numbers.
:::

## Custom Validation Messages with `setCustomValidity`

The browser's default messages ("Please fill out this field") are functional but generic. You can replace them with your own using the Constraint Validation API.

```html
<form id="signup-form">
  <label for="password">Password</label>
  <input type="password" id="password" name="password"
         required minlength="8">

  <label for="confirm">Confirm Password</label>
  <input type="password" id="confirm" name="confirm" required>

  <button type="submit">Create Account</button>
</form>

<script>
const password = document.getElementById('password');
const confirm = document.getElementById('confirm');

confirm.addEventListener('input', () => {
  if (confirm.value !== password.value) {
    confirm.setCustomValidity('Passwords do not match.');
  } else {
    confirm.setCustomValidity(''); // Clear = valid
  }
});
</script>
```

Call `setCustomValidity('')` (empty string) to mark the field as valid. Any non-empty string makes it invalid and that string becomes the error message.

:::warning
If you call `setCustomValidity('some error')` and forget to clear it later, the field will *always* be invalid — even after the user fixes the problem. Always clear it with an empty string when the input becomes valid.
:::

## The Constraint Validation API

JavaScript gives you full access to the browser's validation system through these properties and methods on form elements:

```javascript
const input = document.getElementById('email');

// Check if the value is valid right now
input.validity.valid        // true/false
input.validity.valueMissing // true if required and empty
input.validity.typeMismatch // true if wrong format for type
input.validity.tooShort     // true if shorter than minlength
input.validity.tooLong      // true if longer than maxlength
input.validity.patternMismatch // true if pattern doesn't match
input.validity.rangeUnderflow  // true if below min
input.validity.rangeOverflow   // true if above max

// Get the browser's validation message
input.validationMessage     // "Please include an '@'..."

// Manually trigger validation UI
input.reportValidity();     // Shows the tooltip

// Check without showing the tooltip
input.checkValidity();      // Returns true/false
```

You can use this to build entirely custom error displays:

```javascript
const form = document.getElementById('signup-form');

form.addEventListener('submit', (event) => {
  // Prevent default browser validation UI
  event.preventDefault();

  // Check each field
  const fields = form.querySelectorAll('input, textarea, select');
  let firstInvalid = null;

  fields.forEach(field => {
    const errorEl = document.getElementById(`${field.id}-error`);
    if (!field.checkValidity()) {
      errorEl.textContent = field.validationMessage;
      errorEl.hidden = false;
      field.setAttribute('aria-invalid', 'true');
      if (!firstInvalid) firstInvalid = field;
    } else {
      errorEl.textContent = '';
      errorEl.hidden = true;
      field.removeAttribute('aria-invalid');
    }
  });

  if (firstInvalid) {
    firstInvalid.focus();
  } else {
    form.submit();
  }
});
```

## The `novalidate` Attribute

Adding `novalidate` to a form disables all built-in browser validation.

```html
<form action="/signup" method="post" novalidate>
  …
</form>
```

Why would you want this? When you're building a **custom validation UI** with JavaScript. You still use the Constraint Validation API to *check* validity, but you take full control of how errors are displayed instead of relying on the browser's default tooltips.

:::tip
Use `novalidate` only when you are implementing a custom error display with JavaScript. If you're not, let the browser do the work — its built-in validation is free and well-tested.
:::

## Accessible Error Display

Validation is useless if users can't find the error messages. Here's the accessible pattern:

```html
<div>
  <label for="email">Email</label>
  <input type="email" id="email" name="email" required
         aria-describedby="email-error">
  <span id="email-error" role="alert" hidden></span>
</div>
```

Key techniques:

1. **`aria-describedby`** links the input to its error message, so screen readers announce the error when the field is focused.
2. **`role="alert"`** makes the error a live region — screen readers announce it immediately when it appears.
3. **`aria-invalid="true"`** (set via JavaScript) explicitly marks the field as invalid.
4. **Focus the first invalid field** after validation so keyboard users land right where the problem is.

```html
<!-- After validation fails -->
<input type="email" id="email" name="email" required
       aria-describedby="email-error"
       aria-invalid="true">
<span id="email-error" role="alert">
  Please enter a valid email address.
</span>
```

:::key
Always connect error messages to their fields with `aria-describedby`, use `role="alert"` so screen readers announce errors immediately, and focus the first invalid field. This is the difference between a form that works and a form that's accessible.
:::

## Validation Timing Best Practices

When you show errors matters as much as how:

- **Don't validate on every keystroke from the start.** Users haven't finished typing yet — flagging "invalid email" after they type the first letter is annoying.
- **Validate on `blur`** (when the user leaves a field) for the first check.
- **Validate on `input`** (each keystroke) *after* the field has been marked invalid, so the error clears the moment the user fixes it.
- **Validate on `submit`** as a final catch-all.

```javascript
field.addEventListener('blur', () => validate(field));
field.addEventListener('input', () => {
  if (field.getAttribute('aria-invalid') === 'true') {
    validate(field); // Only re-check if already invalid
  }
});
```

This pattern — validate on blur, re-validate on input once invalid — gives the smoothest user experience.

:::quiz
Q: What happens when you call `setCustomValidity('Error!')` on an input and never clear it?
- The error is shown once and then disappears
- The field is permanently marked as invalid until you call `setCustomValidity('')` *
- The form submits normally, ignoring the error
E: `setCustomValidity` with a non-empty string makes the field invalid. It stays invalid until you explicitly clear it with an empty string, even if the user changes the value.
:::

:::quiz
Q: Which timing strategy gives the best user experience for form validation?
- Validate every field on every keystroke from the start
- Validate on blur first, then on input only after a field has been marked invalid *
- Only validate on form submission
E: Validating on blur gives feedback after the user moves on. Validating on input *only after* a field is invalid lets the error clear immediately when fixed — without pestering users who are still typing.
:::

## Recap

- `required`, `minlength`/`maxlength`, `min`/`max`, and `pattern` give you powerful validation with zero JavaScript.
- `type="email"` and `type="url"` validate format automatically; `type="tel"` does not — use `pattern`.
- `setCustomValidity()` lets you write custom error messages; always clear with `''` when valid.
- The Constraint Validation API (`validity`, `checkValidity`, `reportValidity`) gives JavaScript full access to the browser's validation engine.
- Use `novalidate` only when building a custom error UI with JavaScript.
- Connect errors with `aria-describedby`, use `role="alert"`, set `aria-invalid`, and focus the first invalid field.
- Validate on blur first, then on input once a field is flagged — never on every keystroke from the start.

**Next up:** Semantic HTML — choosing elements for meaning, not appearance.

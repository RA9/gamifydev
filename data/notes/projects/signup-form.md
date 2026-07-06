# Accessible Signup Form

Forms are one of the highest-value frontend components to learn because they combine semantics, styling, validation, and feedback in one place.

:::project
You'll build a polished signup form with visible labels, strong layout, client-side validation, inline error messages, and a success state.
:::

## Step 1 — Write the semantic form markup

Start with a form that already has good structure before any JavaScript.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Signup Form</title>
</head>
<body>
  <section class="signup-card">
    <h1>Create your account</h1>
    <form id="signupForm" novalidate>
      <label for="name">Full name</label>
      <input id="name" name="name" type="text" />
      <p class="error" id="nameError"></p>

      <label for="email">Email address</label>
      <input id="email" name="email" type="email" />
      <p class="error" id="emailError"></p>

      <label for="password">Password</label>
      <input id="password" name="password" type="password" />
      <p class="error" id="passwordError"></p>

      <button type="submit">Create account</button>
      <p id="successMessage" class="success hidden"></p>
    </form>
  </section>
</body>
</html>
```

## Step 2 — Style the form layout

Add a clean card surface, field spacing, and visible input styling.

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
  width: min(100% - 2rem, 28rem);
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  border: 1px solid #e2e8f0;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
}
label {
  display: block;
  margin-top: 1rem;
  font-weight: 600;
}
input {
  width: 100%;
  padding: 0.875rem;
  margin-top: 0.375rem;
  border: 1px solid #cbd5e1;
  border-radius: 0.75rem;
}
```

## Step 3 — Create a validation helper

In JavaScript, start by reading the current field values and checking them.

```js
function validate(values) {
  const errors = {};

  if (!values.name.trim()) {
    errors.name = "Please enter your full name";
  }

  if (!values.email.includes("@")) {
    errors.email = "Please enter a valid email address";
  }

  if (values.password.length < 8) {
    errors.password = "Password must be at least 8 characters";
  }

  return errors;
}
```

This keeps your rules in one place instead of scattering them across click handlers.

## Step 4 — Handle form submission and show errors

Listen for the submit event, stop the default reload, validate the fields, and write any errors into the page.

```js
const form = document.getElementById("signupForm");
const successMessage = document.getElementById("successMessage");

form.addEventListener("submit", function (event) {
  event.preventDefault();

  const values = {
    name: form.name.value,
    email: form.email.value,
    password: form.password.value
  };

  const errors = validate(values);

  document.getElementById("nameError").textContent = errors.name || "";
  document.getElementById("emailError").textContent = errors.email || "";
  document.getElementById("passwordError").textContent = errors.password || "";
```

At this point, if `errors` has any keys, stop and let the user correct them.

## Step 5 — Add a success state

If validation passes, show a helpful success message and reset the form.

```js
  if (Object.keys(errors).length === 0) {
    successMessage.textContent = "Account created — check your inbox for a welcome email.";
    successMessage.classList.remove("hidden");
    form.reset();
  } else {
    successMessage.classList.add("hidden");
    successMessage.textContent = "";
  }
});
```

This makes the interaction feel complete, not abrupt.

## Step 6 — Polish error and focus states

Add visual feedback so the form feels intentional.

```css
button {
  margin-top: 1.25rem;
  width: 100%;
  padding: 0.9rem;
  border: none;
  border-radius: 999px;
  background: #5b3df5;
  color: white;
  font-weight: 700;
}
.error {
  min-height: 1.2rem;
  margin: 0.35rem 0 0;
  color: #dc2626;
  font-size: 0.9rem;
}
.success {
  margin-top: 1rem;
  color: #15803d;
  font-weight: 600;
}
.hidden {
  display: none;
}
input:focus {
  outline: 2px solid #c4b5fd;
  border-color: #8b5cf6;
}
```

That final focus work matters: accessible forms should be easy to navigate and understand.
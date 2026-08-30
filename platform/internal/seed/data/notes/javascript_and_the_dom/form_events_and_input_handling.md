# Form Events and Input Handling

Forms are the primary way users send data to your application. This lesson covers how to listen for form submissions, read every type of input, validate in real time, and use the modern `FormData` API to collect everything at once.

## The submit event

When a user presses Enter inside a form or clicks a submit button, the browser fires a `submit` event on the `<form>`. By default, the browser reloads the page — you almost always want to prevent that:

```js
const form = document.querySelector("#signup-form");

form.addEventListener("submit", (event) => {
  event.preventDefault(); // stop the page reload
  // handle the data with JavaScript instead
});
```

:::key
Always call `event.preventDefault()` in your submit handler unless you intentionally want the browser to navigate to the form's `action` URL. Without it, the page reloads and your JavaScript state is lost.
:::

## Reading input values

### Text inputs and textareas

```js
const nameInput = document.querySelector("#name");
const bio = document.querySelector("#bio");

console.log(nameInput.value); // whatever the user typed
console.log(bio.value);       // full textarea contents
```

### Checkboxes — .checked

Checkboxes have a `.checked` boolean, not a meaningful `.value`:

```js
const agree = document.querySelector("#agree");
console.log(agree.checked); // true or false
```

### Select dropdowns — .value

```js
const country = document.querySelector("#country");
console.log(country.value); // the value of the selected <option>
```

For multi-selects, you need to read the `selectedOptions`:

```js
const selected = Array.from(country.selectedOptions).map(opt => opt.value);
```

### Radio buttons

Radio buttons sharing the same `name` form a group. Get the checked one:

```js
const plan = document.querySelector("input[name='plan']:checked");
if (plan) console.log(plan.value); // "free", "pro", etc.
```

## input vs change events

Two events track user edits, and they fire at different times:

- **`input`** — fires on *every* keystroke, paste, or edit. Use it for live feedback.
- **`change`** — fires when the user *leaves* the field (blur) or selects a new option. Use it for less frequent updates.

```js
const search = document.querySelector("#search");

// Live filtering — updates with every character
search.addEventListener("input", (event) => {
  filterResults(event.target.value);
});

// Validate only when the user is done typing
const email = document.querySelector("#email");
email.addEventListener("change", (event) => {
  validateEmail(event.target.value);
});
```

:::tip
For `<select>`, `<input type="checkbox">`, and `<input type="radio">`, `input` and `change` fire at the same time (on selection). The difference matters mainly for text fields.
:::

## The FormData API

`FormData` collects every named input in a form into a single object — no manual queries needed:

```html
<form id="profile-form">
  <input name="username" value="ada" />
  <input name="email" value="ada@example.com" />
  <textarea name="bio">Mathematician</textarea>
  <input type="checkbox" name="newsletter" checked />
  <button type="submit">Save</button>
</form>
```

```js
const form = document.querySelector("#profile-form");

form.addEventListener("submit", (event) => {
  event.preventDefault();

  const data = new FormData(form);

  console.log(data.get("username"));    // "ada"
  console.log(data.get("email"));       // "ada@example.com"
  console.log(data.get("bio"));         // "Mathematician"
  console.log(data.has("newsletter"));  // true (checkbox was checked)
});
```

To convert FormData into a plain object:

```js
const obj = Object.fromEntries(data);
// { username: "ada", email: "ada@example.com", bio: "Mathematician", newsletter: "on" }
```

:::warning
`Object.fromEntries` only keeps the *last* value for each name. If you have multiple inputs with the same name (like checkboxes), use `data.getAll("name")` to get an array of all values.
:::

## Real-time validation pattern

A common pattern: validate as the user types (with `input`) and show feedback immediately.

```html
<form id="register">
  <label>
    Username
    <input id="username" name="username" required minlength="3" />
    <span class="hint"></span>
  </label>
  <button type="submit">Register</button>
</form>
```

```css
.hint.error { color: crimson; }
.hint.success { color: green; }
```

```js
const username = document.querySelector("#username");
const hint = document.querySelector(".hint");

username.addEventListener("input", () => {
  const value = username.value.trim();

  if (value.length === 0) {
    hint.textContent = "";
    hint.className = "hint";
  } else if (value.length < 3) {
    hint.textContent = "Too short — at least 3 characters";
    hint.className = "hint error";
  } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
    hint.textContent = "Only letters, numbers, and underscores";
    hint.className = "hint error";
  } else {
    hint.textContent = "Looks good!";
    hint.className = "hint success";
  }
});

const form = document.querySelector("#register");
form.addEventListener("submit", (event) => {
  event.preventDefault();
  const value = username.value.trim();
  if (value.length < 3) {
    username.focus();
    return;
  }
  // proceed with submission
});
```

## Debouncing input

For expensive operations (API calls, heavy filtering), you do not want to run logic on every single keystroke. **Debouncing** waits until the user pauses typing:

```js
function debounce(fn, delay) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

const search = document.querySelector("#search");

const handleSearch = debounce((query) => {
  console.log("Searching for:", query);
  // fetch results from API
}, 300);

search.addEventListener("input", (event) => {
  handleSearch(event.target.value);
});
```

The function only runs 300ms after the user *stops* typing. Every new keystroke resets the timer.

:::tip
300ms is a good default debounce delay. Short enough to feel responsive, long enough to avoid firing on every character of a fast typist. Increase to 500ms for network requests; decrease to 150ms for local filtering.
:::

## Disabling submit during processing

Prevent double submissions by disabling the button:

```js
form.addEventListener("submit", async (event) => {
  event.preventDefault();
  const submitBtn = form.querySelector("button[type='submit']");

  submitBtn.disabled = true;
  submitBtn.textContent = "Saving...";

  try {
    await saveData(new FormData(form));
    submitBtn.textContent = "Saved!";
  } catch {
    submitBtn.textContent = "Save";
    submitBtn.disabled = false;
  }
});
```

## Practice

:::quiz
Q: What is the difference between the `input` and `change` events on a text input?
- `input` fires on every edit; `change` fires when the field loses focus *
- They are the same event with different names
- `change` fires on every edit; `input` fires on blur
- `input` only works on text inputs; `change` works on all elements
E: `input` fires immediately on every keystroke, paste, or other edit. `change` fires only when the user finishes editing and moves to another element (or presses Enter). Use `input` for live feedback and `change` for deferred validation.
:::

:::quiz
Q: You have a form with `<input name="email">`, `<input name="name">`, and a submit button. What does `new FormData(form)` give you?
- An array of input elements
- A key-value collection with entries for "email" and "name" *
- A JSON string
- The form's HTML
E: `FormData` automatically collects every named input in the form into a key-value structure. Use `.get("name")` to read individual values or `Object.fromEntries(data)` to convert to a plain object.
:::

## Recap

- Listen for the **`submit`** event on the `<form>` element, and call `event.preventDefault()` to stop the page reload.
- Read text with **`.value`**, checkboxes with **`.checked`**, selects with **`.value`**, and radio groups with a `:checked` selector.
- **`input`** fires on every edit (great for live feedback); **`change`** fires on blur (good for deferred validation).
- **`FormData`** collects all named inputs at once — use `.get()`, `.has()`, `.getAll()`, or convert with `Object.fromEntries()`.
- **Real-time validation:** listen on `input`, check conditions, and update a hint element.
- **Debounce** expensive operations so they only run after the user pauses typing.
- **Disable** the submit button during processing to prevent double submissions.

**Next up:** The State-Render Pattern — the key to building bigger features without losing control.

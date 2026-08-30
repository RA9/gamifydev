# HTML Forms and Inputs

Forms are how websites *listen*. Every search box, login screen, and checkout page is a form collecting information from you. In this lesson you'll learn to build accessible, validated forms from scratch.

## The `<form>` Element

A form wraps a group of inputs and decides what happens when the user submits.

```html
<form action="/signup" method="post">
  <!-- inputs go here -->
  <button type="submit">Sign up</button>
</form>
```

Two important attributes:

- `action` — the URL where the data is sent.
- `method` — *how* the data is sent. The two common choices are `get` and `post`.

**GET vs POST** is a question every beginner asks:

- `method="get"` appends the data to the URL as a query string (`/search?q=cats`). Use it for searches and anything safe to bookmark or repeat. The data is visible in the address bar.
- `method="post"` sends the data in the request body, not the URL. Use it for sign-ups, logins, and anything that changes data or includes passwords.

:::key
Rule of thumb: if the form *reads* data (a search), use GET. If it *changes* something or carries sensitive info (a password, a payment), use POST.
:::

## Labels and Why They Matter

Every input should have a `<label>`. A label is the human-readable name for a field — but it does far more than look nice. You connect a label to its input by matching the label's `for` attribute to the input's `id`.

```html
<label for="email">Email address</label>
<input type="email" id="email" name="email">
```

Now `for="email"` points to `id="email"`. This connection gives you two big wins:

1. **Clicking the label focuses the input** — a bigger, friendlier click target.
2. **Screen readers announce the label** when the user reaches the field, so they know what to type.

:::warning
A placeholder is NOT a label. Placeholder text vanishes the moment you start typing, leaving people unsure what the field was for. Always use a real `<label>`, even if you also add a placeholder.
:::

You can also wrap the input inside the label, which removes the need for `for`/`id`:

```html
<label>
  Email address
  <input type="email" name="email">
</label>
```

Both styles are fine; the `for`/`id` version is the most common.

## Input Types

The single most powerful attribute on `<input>` is `type`. It changes the keyboard on phones, the validation rules, and the on-screen control.

```html
<input type="text" name="username">
<input type="email" name="email">
<input type="password" name="pw">
<input type="number" name="age">
<input type="tel" name="phone">
<input type="url" name="website">
<input type="date" name="dob">
<input type="checkbox" name="newsletter">
<input type="radio" name="plan" value="free">
<input type="file" name="avatar">
<input type="range" name="volume" min="0" max="100">
<input type="color" name="theme">
```

A quick tour:

- `text` — plain single-line text.
- `email` — expects an email; mobile keyboards show the `@` key, and the browser checks the format.
- `password` — hides characters as dots.
- `number` — numeric input with up/down steppers.
- `tel` — phone numbers; shows a numeric keypad on phones.
- `url` — expects a web address.
- `date` — shows a date picker.
- `checkbox` — an on/off toggle; several can be checked at once.
- `radio` — pick exactly one from a group (more on grouping below).
- `file` — lets the user upload a file.
- `range` — a slider between `min` and `max`.
- `color` — a color picker.

:::tip
Choosing the right `type` is free accessibility and UX. `type="email"` on a phone brings up the `@` key automatically; `type="number"` brings up the number pad. Picking the right type saves your users real effort.
:::

## Placeholders

A placeholder is faint hint text shown inside an empty field.

```html
<label for="name">Full name</label>
<input type="text" id="name" name="name" placeholder="e.g. Ada Lovelace">
```

Use it for an *example* of the expected value — never as the field's only label.

## Textareas

For multi-line text (comments, messages), use `<textarea>` instead of `<input>`. Note it has a separate closing tag.

```html
<label for="bio">About you</label>
<textarea id="bio" name="bio" rows="4" placeholder="Tell us about yourself"></textarea>
```

## Select Menus

A dropdown is a `<select>` containing `<option>` elements.

```html
<label for="country">Country</label>
<select id="country" name="country">
  <option value="">Choose one</option>
  <option value="lr">Liberia</option>
  <option value="us">United States</option>
  <option value="gh">Ghana</option>
</select>
```

The `value` is what gets sent to the server; the text between the tags is what the user sees. The first empty option acts as a "please choose" prompt.

## Buttons

A `<button>` triggers an action. Its `type` matters more than people expect.

```html
<button type="submit">Submit form</button>
<button type="button">Just a button</button>
<button type="reset">Clear form</button>
```

- `type="submit"` sends the form. This is the default inside a form, so a `<button>` with no type *will submit*.
- `type="button"` does nothing on its own — you wire it up with JavaScript later.
- `type="reset"` clears all fields (use sparingly; people hate accidental resets).

:::warning
A common bug: you add a `<button>` inside a form to do something with JavaScript, but forget `type="button"`. The page reloads because the button defaulted to `submit`. Always set the type on buttons inside forms.
:::

## Built-in Validation

The browser can check input *before* it's sent, with zero JavaScript. Just add attributes.

```html
<input type="text" name="username" required minlength="3" maxlength="20">
<input type="number" name="age" min="13" max="120">
<input type="text" name="zip" pattern="[0-9]{5}" title="Five digits">
<input type="email" name="email" required>
```

- `required` — the field can't be empty.
- `minlength` / `maxlength` — limits the number of characters.
- `min` / `max` — limits numeric (or date) values.
- `pattern` — a regular expression the value must match. Pair it with `title` to explain the rule.

If a field fails, the browser blocks submission and shows a message. It's the cheapest validation you'll ever get.

:::tip
Browser validation is a first line of defense for usability, not security. Always re-check data on the server too — a determined user can bypass the browser entirely.
:::

## Grouping with Fieldset and Legend

Related controls — especially radio buttons — belong together visually and semantically. `<fieldset>` draws a group; `<legend>` titles it.

```html
<fieldset>
  <legend>Choose a plan</legend>
  <label><input type="radio" name="plan" value="free"> Free</label>
  <label><input type="radio" name="plan" value="pro"> Pro</label>
  <label><input type="radio" name="plan" value="team"> Team</label>
</fieldset>
```

Notice every radio shares `name="plan"`. That shared name is what makes them **mutually exclusive** — selecting one deselects the others. Each has a distinct `value` so the server knows which was chosen.

:::key
Radios in the same group MUST share the same `name`. If each radio has a different name, the user can select all of them — a classic beginner bug.
:::

:::quiz
Q: What makes a set of radio buttons act as one group where only one can be selected?
- They are inside the same `<fieldset>`
- They all share the same `name` attribute *
- They all have the same `id`
E: The shared `name` ties radios into one group; `id` must be unique and `<fieldset>` is only visual grouping.
:::

## Worked Example: A Complete Sign-Up Form

Let's put it all together into a real, accessible sign-up form.

```html
<form action="/signup" method="post">
  <h2>Create your account</h2>

  <label for="name">Full name</label>
  <input type="text" id="name" name="name" required minlength="2"
         placeholder="Ada Lovelace">

  <label for="email">Email</label>
  <input type="email" id="email" name="email" required
         placeholder="you@example.com">

  <label for="password">Password</label>
  <input type="password" id="password" name="password" required minlength="8">

  <label for="age">Age</label>
  <input type="number" id="age" name="age" min="13" max="120">

  <fieldset>
    <legend>Plan</legend>
    <label><input type="radio" name="plan" value="free" checked> Free</label>
    <label><input type="radio" name="plan" value="pro"> Pro</label>
  </fieldset>

  <label>
    <input type="checkbox" name="newsletter"> Send me the newsletter
  </label>

  <label for="bio">About you</label>
  <textarea id="bio" name="bio" rows="4"></textarea>

  <button type="submit">Sign up</button>
</form>
```

:::example
Every input above has a label tied by `for`/`id` (or wrapping). The password requires 8 characters, the email is validated by type, and the plan radios share a name with one `checked` by default. This is what production-quality form markup looks like.
:::

:::quiz
Q: A `<button>` sits inside a `<form>` with no `type` attribute. What happens when clicked?
- Nothing — it needs JavaScript first
- It submits the form *
- It clears the form
E: Inside a form, a button defaults to `type="submit"`, so it will submit unless you set `type="button"`.
:::

## Accessibility of Labels (One More Time)

It bears repeating: labels are the backbone of accessible forms. A screen reader user moving through your form hears each field's label. With no label, they hear "edit text, blank" — useless. With a label, they hear "Email, edit text." Always label every field.

## Recap

- `<form>` has `action` (where) and `method` (`get` for reads, `post` for changes/passwords).
- Tie every `<label>` to its input with `for`/`id` (or wrap the input) — placeholders are not labels.
- The `type` attribute picks the right control, keyboard, and validation: `text`, `email`, `password`, `number`, `tel`, `url`, `date`, `checkbox`, `radio`, `file`, `range`, `color`.
- Use `<textarea>` for long text, `<select>`/`<option>` for dropdowns, and set `<button type>` deliberately.
- Built-in validation (`required`, `minlength`, `maxlength`, `min`, `max`, `pattern`) is free — but still validate on the server.
- Group related fields with `<fieldset>`/`<legend>`; radios in a group must share the same `name`.

**Next up:** Web Accessibility Basics — making your pages usable by everyone.

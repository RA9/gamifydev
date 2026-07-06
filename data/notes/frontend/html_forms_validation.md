# HTML Forms & Validation

Forms are where frontends stop being brochures and start becoming products.

Login flows, signup pages, checkout screens, search bars, support forms, onboarding — they all depend on form quality. Weak forms confuse people, lose conversions, and create accessibility issues fast.

By the end of this lesson, you should be able to build clean, usable HTML forms before JavaScript even enters the picture.

## What a form is doing

A form collects structured user input.

A basic example:

```html
<form>
  <label for="name">Full name</label>
  <input id="name" name="name" type="text" />

  <label for="email">Email</label>
  <input id="email" name="email" type="email" />

  <button type="submit">Create account</button>
</form>
```

The big ideas are:

- users need to know what each field is for
- the browser needs to know what kind of data is expected
- the form should be usable with keyboard, mouse, and assistive tech

## Labels are mandatory, not optional polish

Every meaningful form control needs a label.

```html
<label for="email">Email address</label>
<input id="email" name="email" type="email" />
```

The label's `for` connects to the input's `id`.

Placeholders are **not** a replacement for labels.

Bad:

```html
<input type="email" placeholder="Email" />
```

Why this is weak:

- placeholder text disappears while typing
- it has poorer accessibility support
- it doesn't clearly persist as the field's name

:::quiz
Q: What is the main job of a form label?
- Add visual decoration
- Connect a human-readable field name to a control *
- Replace the need for an input type
E: Labels tell users and assistive tech what the control is for. They are core structure, not decoration.
:::

## Choose the right input type

The `type` attribute gives the browser useful information.

Common examples:

- `text`
- `email`
- `password`
- `tel`
- `url`
- `number`
- `date`
- `checkbox`
- `radio`

```html
<input id="website" name="website" type="url" />
<input id="password" name="password" type="password" />
```

When you choose the right type, you get better mobile keyboards, better native validation, and better user expectations.

:::fill
Q: Complete the input so the browser understands it should expect an email address.
`<input id="email" name="email" type="___" />`
- email *
- text
- password
E: `type="email"` gives the browser more context and enables basic built-in validation.
:::

## `name` attributes matter too

A common beginner question: "Why do I need `name` if I already have `id`?"

- `id` connects labels and scripts to a specific element
- `name` identifies the field in submitted form data

```html
<input id="email" name="email" type="email" />
```

In real apps, you usually want both.

## Group related controls

Use `fieldset` and `legend` when a set of controls belongs together.

```html
<fieldset>
  <legend>Preferred contact method</legend>

  <label>
    <input type="radio" name="contact" value="email" />
    Email
  </label>

  <label>
    <input type="radio" name="contact" value="phone" />
    Phone
  </label>
</fieldset>
```

This gives screen readers and users clearer context.

## Native validation gets you a lot for free

HTML already has strong built-in validation tools.

Useful attributes:

- `required`
- `minlength`
- `maxlength`
- `min`
- `max`
- `pattern`
- `type="email"`, `type="url"`, etc.

```html
<input
  id="password"
  name="password"
  type="password"
  minlength="8"
  required
/>
```

This doesn't replace all JavaScript validation, but it gives you a solid baseline immediately.

:::key
Use native HTML validation first. Add JavaScript validation to improve the experience, not to replace good form structure.
:::

## Helpful defaults and attributes

A few more high-value tools:

- `autocomplete` helps browsers fill known data faster
- `placeholder` can offer an example, but not replace a label
- `textarea` is better than a text input for longer messages
- `select` works well when there are a fixed number of known options

```html
<label for="message">Project details</label>
<textarea id="message" name="message" rows="5"></textarea>
```

## Submit buttons should be explicit

Inside forms, button types matter.

```html
<button type="submit">Send message</button>
<button type="button">Open help</button>
<button type="reset">Clear form</button>
```

If the button should submit the form, make that intent clear with `type="submit"`.

## Layout and usability guidance

Even before CSS gets advanced, good forms usually:

- keep one label per field
- stack fields in a predictable order
- use concise labels
- avoid asking for unnecessary information
- explain errors clearly
- keep submit labels action-focused (`Create account`, `Send message`)

Weak submit labels:

- `Submit`
- `Go`

Better:

- `Join the waitlist`
- `Book consultation`
- `Save profile`

## Mini practice — build a signup form

A strong starter form might include:

- full name
- email address
- password
- role or plan selection
- a checkbox agreeing to terms
- a submit button

Example skeleton:

```html
<form>
  <label for="full-name">Full name</label>
  <input id="full-name" name="fullName" type="text" required />

  <label for="email">Email address</label>
  <input id="email" name="email" type="email" required autocomplete="email" />

  <label for="password">Password</label>
  <input id="password" name="password" type="password" minlength="8" required />

  <label>
    <input type="checkbox" name="terms" required />
    I agree to the terms
  </label>

  <button type="submit">Create account</button>
</form>
```

## Common form mistakes to avoid

- relying on placeholder text instead of labels
- using the wrong input type
- forgetting `name` attributes
- making every button a submit button by accident
- not grouping related radio buttons or checkboxes
- writing vague or hostile error messaging

:::warning
If a form is frustrating, the user doesn't blame the form. They blame the product.
:::

## What good looks like

You should now be able to:

- connect labels correctly with `for` and `id`
- choose appropriate input types
- use `name` attributes intentionally
- group related controls with `fieldset` and `legend`
- use native validation attributes like `required` and `minlength`
- structure a form that's understandable before any CSS or JavaScript is added

## What's next

In **Lab: Build a Semantic Article Page**, you'll put your HTML skills together in a more realistic content build — multiple regions, media, document structure, and clean semantics from top to bottom.
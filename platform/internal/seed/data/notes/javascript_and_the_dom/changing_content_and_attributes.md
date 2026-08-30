# Changing Content and Attributes

Now that you can find elements, it is time to *change* them. This lesson covers every way to read and update the text, HTML, and attributes on an element — and the security trap you must know about.

## textContent — safe plain text

`textContent` gets or sets the text inside an element. It treats everything as plain text — no HTML parsing.

```js
const heading = document.querySelector("h1");

// Reading
console.log(heading.textContent); // "Welcome"

// Writing
heading.textContent = "Hello, world!";
```

If you assign a string that looks like HTML, the user literally sees the tags:

```js
heading.textContent = "<em>Hello</em>";
// renders as: <em>Hello</em>  (visible tags, not italic text)
```

That is exactly what you want when displaying user-supplied text — no tags will be interpreted, so nothing malicious can run.

## innerHTML — powerful but dangerous

`innerHTML` reads or writes the full HTML markup inside an element. The browser parses whatever you assign:

```js
const list = document.querySelector("#items");
list.innerHTML = "<li>Apples</li><li>Bananas</li>";
// creates real <li> elements inside the list
```

This is convenient for building chunks of markup quickly, but it opens a massive security hole if you are not careful.

:::warning
**Never** put user-supplied text into `innerHTML`. A malicious input like `<img src=x onerror="stealCookies()">` would execute JavaScript on your page. This is called a **Cross-Site Scripting (XSS)** attack.

```js
// DANGEROUS — userInput could contain malicious HTML
container.innerHTML = userInput;

// SAFE — textContent escapes everything
container.textContent = userInput;
```

Rule of thumb: use `innerHTML` only with markup **you** wrote. For anything from a user, use `textContent`.
:::

## innerText vs textContent

Both look similar, but they behave differently:

| | `textContent` | `innerText` |
|---|---|---|
| Speed | Fast | Slower (triggers layout) |
| Hidden elements | Includes them | Skips them |
| CSS awareness | No | Yes (respects `display: none`) |

```js
// <span style="display:none">secret</span><span>visible</span>
el.textContent; // "secretvisible"
el.innerText;   // "visible"
```

:::tip
Use `textContent` unless you specifically need the CSS-aware behavior of `innerText`. It is faster and more predictable.
:::

## Reading input values with .value

Form inputs (`<input>`, `<textarea>`, `<select>`) do not use `textContent`. Their current content lives in the `.value` property:

```js
const nameInput = document.querySelector("#name");

// Reading what the user typed
console.log(nameInput.value);

// Setting it programmatically
nameInput.value = "Ada Lovelace";
```

```js
const textarea = document.querySelector("#bio");
console.log(textarea.value); // the full text in the box

const select = document.querySelector("#country");
console.log(select.value); // the currently selected option's value
```

:::key
For inputs, read `.value`. For display elements, read `.textContent`. Mixing them up is one of the most common DOM bugs — `input.textContent` is always an empty string, and `div.value` is `undefined`.
:::

## getAttribute and setAttribute

Every HTML attribute can be read and written through JavaScript:

```js
const link = document.querySelector("a");

// Reading
const href = link.getAttribute("href");
const target = link.getAttribute("target");

// Writing
link.setAttribute("href", "https://gamifydev.com");
link.setAttribute("target", "_blank");
```

Some attributes also have direct property shortcuts:

```js
const img = document.querySelector("img");
console.log(img.src);  // same as img.getAttribute("src")
console.log(img.alt);  // same as img.getAttribute("alt")

img.src = "/images/new-photo.jpg";
img.alt = "A new photo";
```

## The dataset property — custom data-* attributes

HTML lets you attach custom data to any element with `data-*` attributes. JavaScript reads them through the `dataset` property:

```html
<button data-action="delete" data-item-id="42">Remove</button>
```

```js
const btn = document.querySelector("button");

console.log(btn.dataset.action);  // "delete"
console.log(btn.dataset.itemId);  // "42"  (note: camelCase)
```

The naming rule: `data-item-id` in HTML becomes `dataset.itemId` in JavaScript. Dashes are converted to camelCase.

You can also set data attributes:

```js
btn.dataset.status = "pending";
// The element now has data-status="pending" in the DOM
```

:::tip
`data-*` attributes are the cleanest way to attach metadata to elements — use them to store ids, types, categories, or any other value your JavaScript needs to read from the DOM.
:::

## The hidden attribute

The `hidden` attribute is a built-in way to hide elements. It works like `display: none`:

```js
const message = document.querySelector("#error-message");

// Hide it
message.hidden = true;

// Show it
message.hidden = false;

// Toggle
message.hidden = !message.hidden;
```

It is a boolean property — set it to `true` or `false`, no `setAttribute` needed.

## The disabled attribute

Form elements can be disabled to prevent interaction:

```js
const submitBtn = document.querySelector("#submit");

// Disable during form submission
submitBtn.disabled = true;
submitBtn.textContent = "Submitting...";

// Re-enable after
submitBtn.disabled = false;
submitBtn.textContent = "Submit";
```

This is common for preventing double-clicks on form submissions.

## Putting it together — a practical example

A user profile card that updates from data:

```js
function renderProfile(user) {
  const card = document.querySelector("#profile-card");
  
  card.querySelector(".name").textContent = user.name;
  card.querySelector(".bio").textContent = user.bio;
  card.querySelector("img").src = user.avatar;
  card.querySelector("img").alt = `${user.name}'s avatar`;
  card.querySelector("a").setAttribute("href", user.website);
  card.dataset.userId = user.id;
  
  // Hide the "no profile" message
  document.querySelector("#no-profile").hidden = true;
  card.hidden = false;
}

renderProfile({
  id: 7,
  name: "Ada Lovelace",
  bio: "First programmer",
  avatar: "/img/ada.jpg",
  website: "https://example.com",
});
```

Notice how each property uses the right tool: `textContent` for displayed text, `.src` for the image, `setAttribute` for the link, `dataset` for metadata, and `hidden` for visibility.

## Practice

:::quiz
Q: You are building a comment section and need to display a user's comment on the page. Which is the safe way?
- `commentEl.innerHTML = comment.text`
- `commentEl.textContent = comment.text` *
- `commentEl.innerText = comment.text`
- `commentEl.setAttribute("text", comment.text)`
E: `textContent` treats the string as plain text, so any HTML a user typed is displayed harmlessly. `innerHTML` would parse it and could execute malicious scripts (XSS).
:::

:::quiz
Q: An `<input id="email">` contains "user@example.com". How do you read that value?
- `document.querySelector("#email").textContent`
- `document.querySelector("#email").innerHTML`
- `document.querySelector("#email").value` *
- `document.querySelector("#email").getAttribute("value")`
E: Input elements store the user's current text in `.value`. `textContent` and `innerHTML` would return empty strings, and `getAttribute("value")` returns the *initial* HTML attribute, not the current content.
:::

## Recap

- **`textContent`** reads/writes plain text — safe for user input, fast, and predictable.
- **`innerHTML`** reads/writes HTML markup — powerful but a **XSS** risk with user data.
- **`innerText`** is CSS-aware (skips hidden elements) but slower — prefer `textContent` in most cases.
- **`.value`** is for form inputs (`<input>`, `<textarea>`, `<select>`), not `textContent`.
- **`getAttribute`/`setAttribute`** work on any HTML attribute; many have shortcut properties (`.src`, `.href`, `.alt`).
- **`dataset`** reads `data-*` attributes as camelCase properties — the clean way to attach metadata.
- **`hidden`** and **`disabled`** are boolean properties you set directly (`el.hidden = true`).

**Next up:** Changing Styles with classList — the right way to make elements look different.

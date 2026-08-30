# DevTools: Console Mastery

The Console is where you talk to JavaScript and where JavaScript talks back. Most developers use `console.log` and nothing else. That's like owning a power drill and only using it as a hammer. This lesson covers the full Console toolkit.

## console.log with multiple arguments

You already know `console.log("hello")`. But it accepts multiple arguments, separated by commas, and prints them all on one line:

```js
const user = "Ada";
const score = 42;
console.log("User:", user, "Score:", score);
// → User: Ada Score: 42
```

This is faster and more readable than string concatenation. You can also pass objects directly:

```js
const config = { theme: "dark", lang: "en" };
console.log("Config:", config);
// → Config: {theme: "dark", lang: "en"}
```

The object is interactive in the console — click the arrow to expand it.

:::tip
Always label your logs. `console.log(x)` prints a mystery number. `console.log("price:", x)` tells you what that number means. When you have ten logs firing, labels are the difference between clarity and chaos.
:::

## console.table

`console.table` turns arrays of objects into a sortable, readable table.

```js
const users = [
  { name: "Ada", role: "Engineer", level: 3 },
  { name: "Linus", role: "Maintainer", level: 5 },
  { name: "Grace", role: "Architect", level: 4 },
];

console.table(users);
```

This renders a real table with headers: `(index)`, `name`, `role`, `level`. Click a column header to sort by that column.

You can also pass an array of column names to display only specific fields:

```js
console.table(users, ["name", "level"]);
```

:::key
`console.table` is the fastest way to inspect structured data. Use it whenever you're debugging arrays of objects — API responses, form data, quiz questions, shopping carts.
:::

## console.warn and console.error

These work like `console.log` but with visual emphasis:

```js
console.warn("API rate limit approaching — 90% used");
console.error("Failed to load user profile: 404 Not Found");
```

- **`console.warn`** — yellow icon, yellow background. Use for non-critical issues.
- **`console.error`** — red icon, red background, includes a stack trace. Use for actual errors.

They help you visually separate important messages from routine logs when the console gets busy.

## console.dir

`console.dir` shows an object as an expandable tree, which is especially useful for DOM elements:

```js
console.log(document.body);    // shows the HTML tag
console.dir(document.body);    // shows the JS object with all its properties
```

`console.log` on a DOM element renders it like HTML. `console.dir` reveals the JavaScript object underneath — its properties, methods, event listeners, and computed values.

## console.time and console.timeEnd

Measure how long something takes with precision:

```js
console.time("data-load");

// ... some operation
fetch("/api/users")
  .then(res => res.json())
  .then(data => {
    console.timeEnd("data-load");
    // → data-load: 142.38ms
  });
```

The string label connects the `time` and `timeEnd` calls. You can have multiple timers running simultaneously with different labels.

```js
console.time("render");
console.time("fetch");

// ... fetch completes
console.timeEnd("fetch");    // → fetch: 95ms

// ... render completes
console.timeEnd("render");   // → render: 210ms
```

## console.group and console.groupEnd

Group related logs under a collapsible header:

```js
console.group("User Authentication");
console.log("Token:", token);
console.log("Expires:", expiry);
console.log("Roles:", roles);
console.groupEnd();

console.group("API Response");
console.table(data.items);
console.log("Total:", data.total);
console.groupEnd();
```

Use `console.groupCollapsed` instead of `console.group` if you want the group collapsed by default — useful for verbose debugging output you want available but not in the way.

## Running JavaScript in the console

The console is a live JavaScript environment connected to the current page. You can:

```js
// Query the page
document.title
document.querySelectorAll("img").length
getComputedStyle(document.body).backgroundColor

// Modify the page
document.body.style.background = "lightyellow"
document.querySelector("h1").textContent = "Hacked!"

// Test logic
[1, 2, 3, 4, 5].filter(n => n % 2 === 0)
// → [2, 4]
```

This is powerful for testing CSS selectors, trying out array methods, or verifying DOM structure before writing code in your editor.

:::warning
Code run in the console affects the live page but is never saved. If you discover a fix in the console, copy it to your source file immediately — it vanishes on refresh.
:::

## The $0 shortcut

In the Elements panel, click any element to select it. In the console, `$0` refers to that selected element.

```js
// Click a button in the Elements panel, then in the console:
$0                          // → <button class="cta">Sign up</button>
$0.textContent              // → "Sign up"
$0.classList                // → DOMTokenList ["cta"]
$0.getBoundingClientRect()  // → its position and size
```

`$1` through `$4` reference the previously selected elements, but `$0` is the one you'll use 99% of the time.

## copy() and $$() helpers

These are console-only utilities — they don't work in your source code.

**`copy()`** — copies any value to your clipboard:

```js
copy(document.title)
copy(JSON.stringify(data, null, 2))
copy($$("a").map(a => a.href))
```

**`$$(selector)`** — shorthand for `document.querySelectorAll`, but returns a real array:

```js
$$("img")                    // → array of all images
$$("img").map(img => img.src)  // → array of all image URLs
$$(".btn").length             // → count all buttons
```

Unlike `document.querySelectorAll` which returns a NodeList, `$$` returns an array you can immediately use with `.map`, `.filter`, etc.

:::tip
`copy($$("a").map(a => a.href))` is a one-liner that copies every link URL on the page to your clipboard. Extremely useful for auditing or scraping.
:::

## Reading stack traces

When an error occurs, the console shows a **stack trace** — the path the code took to reach the error:

```
Uncaught TypeError: Cannot read properties of null (reading 'addEventListener')
    at initApp (app.js:45:22)
    at main (app.js:12:3)
    at app.js:8:1
```

Read it **top to bottom**:

1. **Line 1:** the error type and message — `TypeError: Cannot read properties of null`
2. **Line 2:** the function and file where it happened — `initApp` in `app.js` at line 45
3. **Lines 3-4:** the call chain — `main` called `initApp`, which was called from line 8

The top line is where to look first. In this case, something at line 45 is `null` when it shouldn't be.

## Common error messages decoded

| Error | Translation | Likely cause |
|-------|-------------|-------------|
| `Cannot read properties of null` | You tried to use something that doesn't exist | Wrong `getElementById` — element not found or script ran before DOM loaded |
| `x is not a function` | You called something as a function but it isn't one | Typo in function name, or a variable shadows the function |
| `x is not defined` | This variable/function name doesn't exist in scope | Typo, or forgot to declare with `let`/`const` |
| `Unexpected token` | JavaScript hit something it can't parse | Missing bracket, extra comma, or unclosed string |
| `Assignment to constant variable` | You tried to reassign a `const` | Use `let` if the value needs to change |

:::key
Error messages are not insults — they're instructions. Read them literally: what type of error, what happened, and where. The answer is almost always in the message itself.
:::

## Filtering console output

When the console gets noisy, use the filter bar:

- Type text to filter messages containing that text
- Click the level buttons (Errors, Warnings, Info) to show only that severity
- Right-click a message and choose "Hide messages from [source]" to suppress third-party noise

## Clear the console

- Click the 🚫 icon (top-left of console)
- Or type `console.clear()` in the console
- Or press **Cmd + K** (Mac) / **Ctrl + L** (Windows)

:::quiz
Q: What is the advantage of `console.table` over `console.log` for arrays of objects?
- console.table uses less memory
- It renders a sortable, readable table with column headers instead of a collapsed object *
- It saves the data to a file
- It only works in production
E: `console.table` transforms arrays of objects into a visual table format that's immediately scannable and sortable by any column. It makes patterns in data obvious at a glance.
:::

:::quiz
Q: What does `$0` refer to in the DevTools console?
- The first element on the page
- The currently selected element in the Elements panel *
- The first console message
- The document object
E: `$0` is a console shortcut that always refers to whatever element you've selected (clicked on) in the Elements panel. It lets you quickly inspect that element's properties, styles, and dimensions.
:::

## Recap

- **`console.log`** with multiple args and labels is better than string concatenation.
- **`console.table`** renders arrays of objects as sortable tables.
- **`console.warn`** and **`console.error`** provide visual severity levels.
- **`console.dir`** shows DOM elements as JavaScript objects.
- **`console.time`/`timeEnd`** measures execution duration.
- **`console.group`/`groupEnd`** organizes related logs under collapsible headers.
- **`$0`** references the selected element; **`$$()`** returns real arrays; **`copy()`** copies to clipboard.
- **Stack traces** read top-to-bottom: error type, message, file, and line number.
- **Error messages are instructions** — read them literally.

**Next up:** DevTools Sources and Breakpoints — pausing your code mid-execution and inspecting everything.

# Browser DevTools and Debugging

Every modern browser ships with a powerful set of developer tools built right in. Learning to use them turns "my code doesn't work and I have no idea why" into "let me look at exactly what's happening."

## Opening DevTools

You don't install anything — DevTools are already in Chrome, Edge, Firefox, and Safari.

- Press **F12** (Windows/Linux) or **Cmd + Option + I** (Mac).
- Or right-click any element on a page and choose **Inspect**.
- Or use the browser menu: ⋮ → More Tools → Developer Tools.

:::tip
Right-clicking the *exact* thing you care about and choosing **Inspect** jumps you straight to that element in the Elements panel. It's the fastest way in.
:::

DevTools open as a panel docked to the side or bottom of your window. Along the top you'll see tabs: **Elements**, **Console**, **Sources**, **Network**, and more. Let's walk through the ones you'll use most.

## The Elements panel — inspect and live-edit

The Elements panel shows the live HTML of the page (the **DOM**), plus the CSS applied to whatever you select.

Click any element in the tree and the right-hand **Styles** pane shows every CSS rule affecting it. The best part: you can edit it live.

- Double-click a tag, attribute, or text to change it.
- In the Styles pane, click a value to edit it, or click the checkbox to toggle a rule off.
- Add a brand-new property by clicking the empty space inside a rule and typing.

```css
/* Try changing this live in DevTools to see instant results */
.button {
  background: rebeccapurple;
  padding: 12px 20px;
}
```

Nothing you change here is saved — refresh and it's gone. That's exactly what makes it safe for experimenting.

:::key
The **box model** diagram at the bottom of the Styles pane shows the selected element's content, padding, border, and margin as nested boxes with real pixel numbers. When spacing looks "off," this diagram usually tells you why instantly.
:::

## The Console — errors and running JavaScript

The Console is where JavaScript errors show up and where you can run code yourself.

When something breaks, a red message appears here. You can also type any JavaScript and press Enter to run it against the current page:

```js
document.title;            // read the page title
document.querySelectorAll("a").length;  // count the links
2 + 2;                     // yes, it's also a calculator
```

The `console` object has more than just `log`:

```js
console.log("plain message", someVariable);
console.table([{ name: "Ada", age: 36 }, { name: "Linus", age: 54 }]);
console.warn("this is a yellow warning");
console.error("this is a red error");
console.dir(document.body); // inspect an object as an expandable tree
```

`console.table` is a hidden gem — give it an array of objects and it renders a real, sortable table. Great for inspecting data.

:::example
Stuck on why a variable is wrong? Sprinkle `console.log` before and after the line you suspect:

```js
let total = price * quantity;
console.log("price:", price, "quantity:", quantity, "total:", total);
```

If `price` prints as `"10"` (a string) instead of `10` (a number), you've found your bug — string `"10" * 2` behaves differently than you'd hope.
:::

## The Sources panel and breakpoints

`console.log` is great, but **breakpoints** are the real superpower. A breakpoint pauses your code at an exact line so you can look around while it's frozen.

1. Open the **Sources** panel.
2. Find your `.js` file in the file tree on the left.
3. Click a line number — a blue marker appears. That's a breakpoint.
4. Trigger the code (e.g. click the button that runs it).

Execution freezes on that line. Now you can:

- Hover over any variable to see its current value.
- Read the **Scope** pane to see all variables in play.
- Use the step controls: **Step over** (run the next line), **Step into** (go inside a function call), and **Resume** (continue running).

```js
function checkout(cart) {
  let total = 0;
  for (const item of cart) {
    total += item.price; // set a breakpoint here and watch total grow
  }
  return total;
}
```

:::tip
You can also pause from code itself by writing the keyword `debugger;` on its own line. When DevTools is open, the browser stops there as if you'd clicked a breakpoint.
:::

## The Network tab — see every request

The Network tab records everything the page loads: HTML, CSS, images, fonts, and data from APIs.

Open it, then **refresh the page** so it captures the requests. Each row shows:

- The **Name** of the file or request.
- The **Status** code: `200` means OK, `404` means not found, `500` means server error.
- The **Type** and **Size**.
- The **Time** it took.

Click any row to see its headers, the request you sent, and the response you got back (the **Response** and **Preview** tabs).

:::warning
A broken image or a missing stylesheet almost always shows up here as a red `404`. If your CSS "isn't loading," check the Network tab first — the file path in your `<link>` is probably wrong.
:::

## The device toolbar — responsive testing

Click the small phone/tablet icon (or **Cmd/Ctrl + Shift + M**) to enter device mode. You can:

- Pick a device preset (iPhone, Pixel, iPad) or drag to any custom size.
- Rotate between portrait and landscape.
- Throttle the network to simulate slow connections.

This is how you check that your responsive CSS actually works without owning ten phones.

## How to read an error message

Beginners often see a red error and panic. Don't — errors are *instructions*. Read them top to bottom:

```js
Uncaught TypeError: Cannot read properties of null (reading 'addEventListener')
    at app.js:12:18
```

Decode it piece by piece:

- **`TypeError`** — the *kind* of error.
- **`Cannot read properties of null`** — the *meaning*: you tried to use something that was `null` (nothing).
- **`(reading 'addEventListener')`** — *what* you tried to do.
- **`at app.js:12:18`** — *where*: file `app.js`, line 12, column 18.

That stack trace is a map. The top line is usually where it actually broke. In this case, `document.getElementById(...)` returned `null` because the ID didn't match — a classic typo or a script running before the element exists.

:::analogy
A stack trace is like a "you are here" trail of breadcrumbs showing how the code got to the point of failure. Read the top breadcrumb first; that's the scene of the crime.
:::

## A debugging mindset

Tools help, but *how you think* matters more. Work through this loop:

1. **Reproduce** — make the bug happen on purpose. A bug you can't trigger reliably is a bug you can't fix.
2. **Isolate** — narrow it down. Comment out half your code. Does it still break? Now you know which half.
3. **Check your assumptions** — you *think* a variable is a number; print it and confirm. Most bugs hide in a gap between what you assume and what's true.
4. **`console.log` generously** — print values at every step until reality matches expectation.
5. **Rubber-duck it** — explain the broken code out loud, line by line, to a rubber duck (or anyone). You'll often hear your own mistake before you finish the sentence.

:::quiz
Q: In the Network tab, you see your stylesheet listed with a status code of `404`. What does that tell you?
- The CSS loaded fine but has a syntax error
- The browser could not find the file at the path you requested *
- The server crashed while sending the file
E: A `404` means "not found" — the requested file isn't where the page looked for it. The fix is almost always correcting the path in your `<link>` tag.
:::

:::quiz
Q: What is the main advantage of a breakpoint over a `console.log`?
- It permanently saves your variables to a file
- It pauses execution so you can inspect every variable's live value and step through line by line *
- It makes your code run faster
E: Breakpoints freeze the program mid-run, letting you explore the full state interactively and step forward one line at a time, rather than guessing which values to print ahead of time.
:::

:::fill
The keyword you can write directly in your JavaScript to make the browser pause there (when DevTools is open) is `______;`.
- debugger *
- breakpoint
- pause
E: Writing `debugger;` on its own line triggers a breakpoint at that spot whenever DevTools is open.
:::

## Recap

- DevTools open with **F12** or right-click → **Inspect**; everything you need is already in the browser.
- The **Elements** panel lets you inspect and live-edit HTML/CSS and read the box model.
- The **Console** shows errors and runs JavaScript; `console.log`, `console.table`, `console.warn`, and `console.error` are your friends.
- The **Sources** panel + **breakpoints** (or `debugger;`) let you pause and step through code.
- The **Network** tab reveals requests, status codes (`200`, `404`, `500`), and payloads.
- The **device toolbar** simulates phones and tablets for responsive testing.
- Read error messages and stack traces top-down: kind, meaning, and the file:line where it broke.
- Debug with a process: reproduce, isolate, check assumptions, log, and rubber-duck.

**Next up:** Git and GitHub — tracking your changes and sharing your code with the world.

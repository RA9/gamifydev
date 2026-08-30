# Your First Web Page

Time to build something real. In this lesson you'll create a complete web page from scratch — HTML for structure, CSS for style, JavaScript for interactivity. By the end, you'll have a working page running in your browser with live reload.

## Create your project folder

Open a terminal and set up a clean workspace:

```bash
mkdir my-first-page
cd my-first-page
touch index.html
code .
```

## The HTML boilerplate

Every HTML page starts with the same skeleton. Open `index.html`, type `!`, then press **Tab** — Emmet expands it:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>
</head>
<body>

</body>
</html>
```

| Line                          | Purpose                                          |
|-------------------------------|--------------------------------------------------|
| `<!DOCTYPE html>`             | Tells the browser this is an HTML5 document      |
| `<html lang="en">`           | Root element; `lang` helps screen readers        |
| `<meta charset="UTF-8">`     | Character encoding — supports all languages      |
| `<meta name="viewport" ...>` | Makes the page responsive on mobile devices      |
| `<title>`                     | Text shown in the browser tab                    |
| `<head>`                      | Metadata, links to CSS/JS — not visible on page  |
| `<body>`                      | Everything visible on the page goes here         |

:::key
The `<head>` is for the browser — metadata, page title, links to stylesheets and scripts. The `<body>` is for the user — everything they see and interact with.
:::

## Add content to the body

Replace the title and add content:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>My First Page</title>
</head>
<body>
  <h1>Hello, World!</h1>
  <p>This is my first web page. I built it from scratch.</p>

  <h2>Things I want to learn</h2>
  <ul>
    <li>HTML — page structure</li>
    <li>CSS — styling and layout</li>
    <li>JavaScript — interactivity</li>
  </ul>

  <h2>About me</h2>
  <p>I'm learning web development and this is just the beginning.</p>
  <a href="https://example.com">Visit my other site</a>
</body>
</html>
```

Key elements: `<h1>`–`<h6>` for headings, `<p>` for paragraphs, `<ul>`/`<li>` for lists, `<a>` for links.

## Open in the browser with Live Server

In VS Code, right-click `index.html` → **"Open with Live Server"**. Your browser opens and the page **auto-reloads** on every save. This is the workflow you'll use from now on.

:::tip
If "Open with Live Server" doesn't appear, make sure you installed the Live Server extension and opened the **folder** in VS Code (not just the file).
:::

## Create and link a CSS file

Create `style.css` in the same folder and link it in your HTML `<head>`:

```html
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>My First Page</title>
  <link rel="stylesheet" href="style.css">
</head>
```

## Write your first CSS

Open `style.css` and add:

```css
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: system-ui, -apple-system, sans-serif;
  line-height: 1.6;
  color: #1e293b;
  max-width: 680px;
  margin: 0 auto;
  padding: 32px 24px;
}

h1 { font-size: 2.5rem; margin-bottom: 8px; }
h2 { font-size: 1.5rem; margin-top: 32px; margin-bottom: 12px; }
p  { margin-bottom: 16px; }
ul { margin-bottom: 16px; padding-left: 24px; }
li { margin-bottom: 6px; }

a {
  color: #3b82f6;
  text-decoration: none;
}
a:hover { text-decoration: underline; }
```

Save and check your browser. The page transforms — clean typography, centered content, styled links. Same HTML, completely different appearance.

:::key
CSS connects to HTML via the `<link>` tag in `<head>`. CSS *selects* HTML elements (like `h1`, `p`, `a`) and applies styles to them. Change the CSS and the appearance changes — the HTML stays the same.
:::

## Create and link a JavaScript file

Create `script.js` and link it at the **end of `<body>`**:

```html
  <a href="https://example.com">Visit my other site</a>

  <script src="script.js"></script>
</body>
</html>
```

:::warning
Place `<script>` tags at the end of `<body>`, just before `</body>`. This ensures the HTML is fully loaded before JavaScript runs. If the script runs before elements exist, `querySelector` returns `null` and your code breaks.
:::

## Write your first JavaScript

Open `script.js`:

```javascript
const heading = document.querySelector('h1');

heading.addEventListener('click', function () {
  heading.textContent = 'You clicked me! 🎉';
  heading.style.color = '#3b82f6';
});

console.log('JavaScript is connected and running!');
```

Save and check: open DevTools (`F12`) → **Console** tab — you should see the log message. Click the heading — it changes text and color.

## Add a button with interactivity

Add this to the body, before the script tag:

```html
<h2>Counter</h2>
<button id="counter-btn">Clicked 0 times</button>
```

Update `script.js`:

```javascript
const heading = document.querySelector('h1');
const counterBtn = document.getElementById('counter-btn');
let count = 0;

heading.addEventListener('click', function () {
  heading.textContent = 'You clicked me! 🎉';
  heading.style.color = '#3b82f6';
});

counterBtn.addEventListener('click', function () {
  count = count + 1;
  counterBtn.textContent = 'Clicked ' + count + ' times';
});
```

Style the button in `style.css`:

```css
button {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 12px 24px;
  font-size: 1rem;
  border-radius: 8px;
  cursor: pointer;
}
button:hover { background: #2563eb; }
```

You now have a styled, interactive counter — all three languages working together.

## How the three languages connect

```
index.html (structure)
    ├── <link href="style.css">     → CSS styles the elements
    └── <script src="script.js">    → JS adds behavior to elements
```

- **HTML** provides the elements (`<h1>`, `<button>`, `<p>`).
- **CSS** selects those elements and styles them.
- **JavaScript** selects those elements and adds behavior.

## Your final project structure

```
my-first-page/
├── index.html       ← structure (HTML)
├── style.css        ← appearance (CSS)
└── script.js        ← behavior (JavaScript)
```

This three-file pattern is the foundation of every web project.

:::quiz
Q: Where should you place a `<script>` tag for best results?
- In the `<head>` without any attributes
- At the end of `<body>`, just before `</body>` *
- Outside the `<html>` element
E: Placing scripts at the end of `<body>` ensures all HTML elements are loaded before the script runs, preventing errors from trying to select elements that don't exist yet.
:::

:::quiz
Q: What does `<link rel="stylesheet" href="style.css">` do?
- Embeds CSS directly into the HTML
- Tells the browser to download and apply an external CSS file *
- Creates a hyperlink to a CSS tutorial
E: The `<link>` tag in `<head>` connects an external CSS file to the HTML document. The browser downloads `style.css` and applies its rules to the page's elements.
:::

## Recap

- Every HTML page starts with a **boilerplate**: `<!DOCTYPE html>`, `<html>`, `<head>`, `<body>`. Use Emmet's `!` shortcut.
- Add content with semantic elements: `<h1>`, `<p>`, `<ul>`, `<img>`, `<a>`.
- Link CSS with `<link rel="stylesheet" href="style.css">` in `<head>`.
- Link JS with `<script src="script.js"></script>` at the end of `<body>`.
- Use **Live Server** for auto-reloading during development.
- The three-file pattern (`index.html`, `style.css`, `script.js`) is the foundation of all web development.

**Congratulations!** You've built your first web page with HTML, CSS, and JavaScript.

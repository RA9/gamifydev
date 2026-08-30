# Setting Up Your Dev Environment

Before you write a single line of code, you need the right tools. A good development environment makes you faster, catches mistakes early, and removes friction from the build-test-debug cycle. This lesson walks you through setting up everything you need to start building web pages.

## Installing VS Code

**Visual Studio Code** (VS Code) is the most popular code editor for web development — it's free, fast, and endlessly customizable.

1. Go to **code.visualstudio.com**.
2. Download the installer for your operating system (Windows, macOS, or Linux).
3. Run the installer with default settings.
4. Open VS Code — you should see the Welcome tab.

:::tip
During installation on Windows, check the boxes for **"Add to PATH"** and **"Add Open with Code to context menu"**. These let you open VS Code from the terminal and right-click any folder to open it in VS Code — small things that save a lot of time.
:::

## Essential extensions

VS Code's power comes from extensions. Install these three to start:

### 1. Live Server

Launches a local development server that **auto-reloads** your page whenever you save a file. No more manually refreshing the browser.

- Open the Extensions sidebar (`Ctrl+Shift+X` / `Cmd+Shift+X`)
- Search for **"Live Server"** by Ritwick Dey
- Click **Install**

To use it: right-click your `index.html` → **"Open with Live Server"**. Your page opens in the browser and refreshes automatically on save.

### 2. Prettier

An opinionated code formatter that automatically formats your HTML, CSS, and JS on save — consistent indentation, spacing, and style without thinking about it.

- Search for **"Prettier - Code formatter"**
- Click **Install**
- Enable format on save: Settings → search "Format On Save" → check the box

```json
// Settings you'll want (VS Code settings.json):
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode"
}
```

### 3. Emmet (built-in)

Emmet is already included in VS Code — it lets you write HTML shortcuts that expand into full elements:

```
div.container>h1+p*3
```

Press **Tab** and it expands to:

```html
<div class="container">
  <h1></h1>
  <p></p>
  <p></p>
  <p></p>
</div>
```

Common Emmet shortcuts:

| Shortcut          | Expands to                              |
|-------------------|-----------------------------------------|
| `!`               | Full HTML boilerplate                   |
| `div.className`   | `<div class="className"></div>`         |
| `ul>li*5`         | `<ul>` with 5 `<li>` children          |
| `a[href="#"]`     | `<a href="#"></a>`                      |
| `p{Hello world}`  | `<p>Hello world</p>`                   |

:::key
Emmet's `!` shortcut is the fastest way to create an HTML boilerplate. Type `!` in an empty `.html` file and press Tab — you get a complete, valid HTML5 document structure instantly.
:::

## Terminal basics

The **terminal** (also called command line or shell) is how you run commands, navigate files, and later use tools like Git and npm. VS Code has a built-in terminal.

Open it: **View → Terminal** (or `` Ctrl+` `` / `` Cmd+` ``).

### Essential commands

```bash
# See where you are
pwd                    # print working directory

# List files in the current folder
ls                     # macOS / Linux
dir                    # Windows (Command Prompt)

# Move into a folder
cd my-project          # change directory
cd ..                  # go up one level
cd ~                   # go to your home directory

# Create a new folder
mkdir my-project       # make directory

# Create a new file
touch index.html       # macOS / Linux
echo. > index.html     # Windows (Command Prompt)

# Clear the terminal
clear                  # macOS / Linux
cls                    # Windows
```

:::tip
You can drag a folder from your file explorer into the VS Code terminal to paste its full path. This saves time when navigating to deeply nested directories.
:::

### A typical workflow

```bash
cd ~/Documents              # go to your Documents folder
mkdir my-first-site         # create a project folder
cd my-first-site            # move into it
touch index.html            # create the HTML file
touch style.css             # create the CSS file
touch script.js             # create the JS file
code .                      # open the current folder in VS Code
```

## Creating a project folder

Every project should live in its own folder with a clear structure. Even a simple project benefits from organization:

```
my-first-site/
├── index.html         # your main HTML page
├── style.css          # your stylesheet
├── script.js          # your JavaScript
└── images/            # a folder for image files
    └── logo.png
```

:::warning
Keep your project folder in an easy-to-find location like `~/Documents/projects/` or `~/Desktop/projects/`. Avoid deeply nested or system-protected directories — they cause permission issues and make terminal navigation painful.
:::

## Opening a project in VS Code

Three ways to open a project:

1. **File → Open Folder** → select your project folder.
2. **Terminal:** navigate to your project and type `code .` (requires VS Code in PATH).
3. **Drag and drop** the folder onto the VS Code icon.

Always open the **folder**, not individual files. This gives VS Code full project context — IntelliSense, file linking, and the file explorer sidebar all work best with a project folder.

## File naming conventions

Web servers and URLs are picky about file names. Follow these rules to avoid broken links and headaches:

```
✓ lowercase:        style.css         (not Style.CSS)
✓ hyphens:          about-us.html     (not about us.html)
✓ no spaces:        my-page.html      (not my page.html)
✓ descriptive:      contact-form.js   (not stuff.js)
✗ no special chars: résumé.html       (use resume.html)
```

:::key
Always use **lowercase** file names with **hyphens** instead of spaces. URLs are case-sensitive on most servers, so `About.html` and `about.html` are *different files*. Spaces in file names become `%20` in URLs, which is messy and error-prone.
:::

Your main HTML page should always be named `index.html` — web servers look for this file by default when someone visits a directory.

## Browser DevTools intro

Every modern browser includes built-in developer tools. They're your primary debugging, inspection, and performance analysis toolkit.

### Opening DevTools

- **Chrome/Edge:** `F12` or `Ctrl+Shift+I` / `Cmd+Option+I`
- **Firefox:** `F12` or `Ctrl+Shift+I` / `Cmd+Option+I`
- **Safari:** Enable in Preferences → Advanced → "Show Develop menu", then `Cmd+Option+I`

### The key panels

```
Elements (Inspector)  — See and edit the live DOM and CSS
Console              — Run JavaScript, see errors and logs
Network              — Track every HTTP request and response
Sources              — Debug JavaScript with breakpoints
Application          — Inspect storage (cookies, localStorage)
```

### Your first DevTools exercise

1. Open any website in Chrome.
2. Press `F12` to open DevTools.
3. Click the **Elements** tab.
4. Click the selector icon (top-left of DevTools panel).
5. Hover over any element on the page — you'll see its box model highlighted.
6. Click the element — its HTML and CSS appear in the DevTools panel.
7. Try editing a CSS value in the Styles panel — changes happen live, right in the browser.

:::tip
Changes you make in DevTools are **temporary** — they don't modify your actual files. This makes DevTools a safe playground to experiment with CSS changes, test fixes, and debug layouts without any risk.
:::

### The Console panel

The Console is where JavaScript errors appear and where you can run code:

```javascript
// Type these directly into the Console:
console.log('Hello from DevTools!');
document.title = 'I changed the title!';
document.body.style.background = 'lightblue';
```

Get in the habit of checking the Console whenever something looks wrong — JavaScript errors appear here with the file name and line number.

## Recommended VS Code settings

A starter configuration that improves the editing experience:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.tabSize": 2,
  "editor.wordWrap": "on",
  "editor.minimap.enabled": false,
  "files.autoSave": "onFocusChange",
  "emmet.includeLanguages": {
    "javascript": "html"
  }
}
```

Open settings: `Ctrl+Shift+P` / `Cmd+Shift+P` → type **"settings json"** → select **"Preferences: Open User Settings (JSON)"**.

:::quiz
Q: What does the Live Server extension do?
- Compiles your code into a production-ready build
- Launches a local server that auto-reloads your page when you save changes *
- Connects your code to a remote hosting server
E: Live Server runs a local development server and watches your files. When you save a change, it automatically refreshes the browser — eliminating the need to manually reload after every edit.
:::

:::quiz
Q: Why should you use lowercase filenames with hyphens (like `about-us.html`) for web projects?
- It's a JavaScript requirement
- URLs are case-sensitive on most servers, and spaces create messy `%20` encoding in URLs *
- Browsers can't load files with uppercase letters
E: Most web servers are case-sensitive, meaning `About.html` and `about.html` are different files. Spaces become `%20` in URLs. Using lowercase with hyphens is the universal convention that avoids both problems.
:::

## Recap

- Install **VS Code** as your code editor — it's free and extensible.
- Install three essential extensions: **Live Server** (auto-reload), **Prettier** (auto-format), and use built-in **Emmet** (HTML shortcuts).
- Learn basic **terminal commands**: `cd`, `ls`, `mkdir`, `touch`, `pwd`.
- Create a dedicated **project folder** with clear file names (lowercase, hyphens, no spaces).
- Open the **folder** in VS Code, not individual files.
- **Browser DevTools** (F12) let you inspect HTML, edit CSS live, debug JavaScript, and monitor network requests.

**Next up:** Your First Web Page — writing real HTML, CSS, and JavaScript and seeing it in the browser.

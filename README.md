# GamifyDev

**Learn to code by building — gamified, offline-friendly, and free.**

GamifyDev turns learning web development into a game. Follow guided lessons,
build real projects, test yourself with quizzes, and earn XP, streaks, and
badges along the way. It runs entirely in the browser, works offline, and
stores all your progress locally — no account required.

> Live: [gamifydev.vercel.app](https://gamifydev.vercel.app)

---

## Features

- **Guided learning paths** — Frontend, Backend, and Fullstack, each a sequence
  of lessons that unlock as you progress.
- **Rich lessons** — narrative explanations with analogies, worked code
  examples, SVG diagrams, callouts, and inline check-for-understanding quizzes.
  Written to be *understood and explained*, not memorised.
- **Hands-on projects** — build a profile card, a quiz game, a REST API, and a
  full-stack notes app, step by step.
- **Test Yourself quizzes** — timed multiple-choice quizzes across HTML, CSS,
  JavaScript, C, Python, Java, and SQL, with a full answer review (every
  question explained).
- **Gamification** — XP and levels, daily streaks, achievement badges, and a
  progress dashboard with per-language best scores and attempt history.
- **Dark mode** — persistent, with no flash on load.
- **Offline-first PWA** — installable, with the app and content cached for use
  without a connection.

### Learning paths

| Path | Modules | Highlights |
| --- | --- | --- |
| **Frontend** | 12 | History of the Web → HTML/CSS/JS → **Profile Card** & **Quiz Game** projects → Git → Portfolio |
| **Backend** | 6 | Servers & APIs → Python → Databases/SQL → REST APIs → **To-Do API** project → Deploy |
| **Fullstack** | 4 | The full stack → fetch & APIs → **Notes App** project → Career next steps |

The quiz banks hold ~390 questions across 8 categories.

---

## Tech stack

- **Vanilla JavaScript** — no framework. A small hash-based router renders pages
  into `<main>`.
- **Tailwind CSS** — compiled to a static stylesheet (`css/tailwind.css`) via
  the Tailwind CLI. A small design system lives in `css/input.css` and
  `tailwind.config.js`.
- **Dexie / IndexedDB** — all progress (questions, paths, scores, state) is
  stored locally in the browser.
- **PWA** — a service worker (Workbox) plus a web manifest for offline use and
  installability.
- **Vercel** — static hosting with edge middleware (`@vercel/edge`).

---

## Project structure

```text
.
├── index.html            # App shell: header, nav, footer, script includes
├── index.js              # Hash router (#about, #test, #progress, …)
├── css/
│   ├── input.css         # Tailwind directives + design system (source)
│   └── tailwind.css      # Compiled stylesheet (generated)
├── js/
│   ├── home.js           # Landing, onboarding, learning-path selection
│   ├── tys.js            # "Test Yourself" quizzes + results/review
│   ├── modules.js        # Lessons, markdown renderer, lesson quizzes
│   ├── progress.js       # XP/streak/badges dashboard
│   ├── storage.js        # Dexie schema + seeding (questions, paths)
│   ├── icons.js          # SVG icon set, language badges, illustrations
│   ├── theme.js          # Dark-mode toggle, mobile menu
│   └── utils.js          # Helpers (ids, shuffle, escaping)
├── data/
│   ├── app.json          # Learning paths and their modules
│   ├── questions.json    # Quiz banks by category
│   └── notes/<path>/*.md # Lesson content (markdown)
├── images/lessons/*.svg  # Lesson diagrams
└── tailwind.config.js    # Design tokens (brand palette, fonts, animations)
```

---

## Getting started

Requires only Node.js (for the Tailwind build) and any static file server.

```bash
# 1. Install dev dependencies (Tailwind CLI)
npm install

# 2. Build the stylesheet
npm run build:css        # or: npm run watch:css  (rebuild on change)

# 3. Serve the folder with any static server, e.g.
npx serve .
```

Then open the served URL in your browser. Progress is saved in IndexedDB, so it
persists across reloads on the same browser.

> **Styling:** the UI is built from utility classes, so after changing any
> markup or `css/input.css` you must re-run `npm run build:css` to regenerate
> `css/tailwind.css`.

---

## Authoring content

Content is data-driven — you can add lessons and questions without touching the
app code.

**Learning paths** (`data/app.json`) — each path is a list of modules linked by
`next`/`previous`. On load, `storage.js` upserts them into IndexedDB by
`(path, title)`, so adding paths or modules reaches existing users without
resetting their progress.

**Lessons** (`data/notes/<path>/<slug>.md`) — the file name is the module title
slugified (e.g. *"HTML Basics"* → `html_basics.md`). The lesson renderer
supports an extended Markdown:

- Headings, **bold**/*italic*/`code`, ordered & unordered lists, links
- Fenced code blocks: ` ```html … ``` `
- Images → captioned figures: `![caption](/images/lessons/diagram.svg)`
- Callout boxes: `:::tip`, `:::analogy`, `:::warning`, `:::key`, `:::example`,
  `:::project` … `:::`
- Inline quizzes:

  ```text
  :::quiz
  Q: What does HTML structure?
  - Style
  - Content *
  - Servers
  E: HTML describes the structure and meaning of content.
  :::
  ```

  (Mark the correct option with a trailing `*`; `E:` is the explanation shown on
  answer.) All lesson input is escaped — no raw HTML is injected.

**Quiz banks** (`data/questions.json`) — questions grouped by category. New
categories/questions are merged in on load without duplicating existing ones.

---

## Design system

- **Brand:** indigo/violet (`brand`) is the identity; **slate** is the neutral.
- **Colour means something:** green = success, amber = warning, red = error.
  Colour is not used decoratively.
- Tokens (palette, the Nunito display font, soft shadows, rounded corners,
  animations) live in `tailwind.config.js`; reusable component classes
  (`gd-card`, `gd-btn`, `gd-callout`, …) live in `css/input.css`.

---

## Credits

Created by [Carlos S. Nah](https://twitter.com/rademejs). Built for beginners,
especially those learning with limited connectivity.

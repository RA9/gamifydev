# GamifyDev

**Learn to code by building — a gamified, offline-friendly, free learning adventure.**

GamifyDev turns learning to code into a game. You join a "startup," pick a
**world** to explore, and follow a winding quest map through guided lessons —
earning XP, building streaks, unlocking achievements, practising in a real
terminal, and shipping build-along projects. Pixel, your in-app guide, cheers
you on the whole way. It runs entirely in the browser, works offline, and stores
all your progress locally — no account required.

> Live: [gamifydev.vercel.app](https://gamifydev.vercel.app)

---

## Features

- **Six worlds to explore** — Frontend, Backend, Fullstack, **C**, **Java**, and
  **Linux & Bash**. Pick one from the Worlds page and switch any time; your
  progress in each is saved independently.
- **A game-style quest map** — each world is a winding road of lessons that
  unlock as you advance, with a one-time mission briefing from Pixel when you
  arrive.
- **Rich, depth-first lessons** — narrative explanations with analogies, worked
  code examples, on-brand SVG diagrams, callouts, and **interactive exercises**
  (predict-the-output, reorder, fill-in-the-blank, match, and check-for-
  understanding quizzes). Written to be *understood and explained*, not memorised.
- **Build-along projects** — a dedicated Projects gallery with step-by-step
  builds you tick off as you go, finishing with a "Shipped it!" celebration.
  Profile Card, Quiz Game & To-Do App (Frontend), Greeting Script (Bash),
  Number Guessing Game (C), and a Bank Account (Java).
- **Terminal Trainer** — a simulated Bash shell with a virtual filesystem and
  nine hands-on missions, so command-line learning is real practice, not theory.
- **Spaced repetition** — a Leitner-box review system resurfaces questions you've
  answered at growing intervals, so what you learn actually sticks.
- **Daily Standup** — a daily loop with an XP goal, your next "ticket," a daily
  challenge, a weekly sprint review, and freeze-aware streaks.
- **Cross-world achievements** — 12 achievements spanning lessons, worlds,
  projects, terminal missions, streaks, and reviews, with live progress on the
  ones you haven't unlocked yet.
- **Test Yourself quizzes** — timed multiple-choice quizzes across nine
  categories, with a full answer review (every question explained).
- **Pixel, your guide** — a friendly mascot who narrates mission briefings, ship
  celebrations, level-ups, and achievement unlocks.
- **Dark mode** — persistent, with no flash on load.
- **Offline-first PWA** — installable, with the app and content cached for use
  without a connection.

### Worlds

| World | Modules | Highlights |
| --- | --- | --- |
| **Frontend** | 12 | History of the Web → HTML/CSS/JS → **Profile Card** & **Quiz Game** projects → Git → Portfolio |
| **Backend** | 6 | Servers & APIs → Python → Databases/SQL → REST APIs → **To-Do API** → Deploy |
| **Fullstack** | 4 | The full stack → fetch & APIs → **Notes App** → Career next steps |
| **C** | 6 | What is C? → basics → control flow → functions → **pointers & memory** → next steps |
| **Java** | 6 | The JVM → basics → control flow → methods → **objects & classes** → next steps |
| **Linux & Bash** | 6 | What is Linux? → the command line → filesystem → files → **scripting** → next steps |

The quiz banks hold **345 questions across 9 categories** (frontend, javascript,
html, css, c, python, java, sql, linux). Lessons across every world include
on-brand SVG diagrams (compile model, pointer memory, JVM bytecode, class vs.
object, filesystem tree, command anatomy, pipes, and more).

---

## Tech stack

- **Vanilla JavaScript** — no framework. A small hash-based router renders pages
  into `<main>`.
- **Tailwind CSS** — compiled to a static stylesheet (`css/tailwind.css`) via
  the Tailwind CLI. A small design system lives in `css/input.css` and
  `tailwind.config.js`.
- **Dexie / IndexedDB** — all progress (questions, paths, scores, daily streaks,
  spaced-repetition reviews, meta) is stored locally in the browser. Seeding is
  migration-safe: new worlds, lessons, and questions reach existing users
  without resetting their progress.
- **PWA** — a web manifest plus a bundled service worker (`sw.js`) for
  installability and offline caching.
- **Vercel** — static hosting with edge middleware (`@vercel/edge`); Node pinned
  to 24.x.

---

## Project structure

```text
.
├── index.html              # App shell: header, nav, footer, script includes
├── index.js                # Hash router (#today, #worlds, #projects, #terminal, …)
├── css/
│   ├── input.css           # Tailwind directives + design system (source)
│   └── tailwind.css        # Compiled stylesheet (generated)
├── js/
│   ├── home.js             # Landing, onboarding, world preference
│   ├── daily.js            # Daily Standup: streaks, XP goal, challenge, sprint review
│   ├── modules.js          # Lessons, markdown renderer, exercises, quest map, Worlds page
│   ├── projects.js         # Build-along Projects gallery + step-tracked player
│   ├── terminal.js         # Simulated Bash shell + Terminal Trainer missions
│   ├── review.js           # Spaced-repetition (Leitner) review system
│   ├── achievements.js     # Cross-world achievements + unlock celebrations
│   ├── mascot.js           # Pixel: avatar, mission briefings, celebrations
│   ├── progress.js         # XP/level/streak + achievements dashboard
│   ├── tys.js              # "Test Yourself" quizzes + results/review
│   ├── storage.js          # Dexie schema + migration-safe seeding
│   ├── icons.js            # SVG icon set, language badges, illustrations
│   ├── theme.js            # Dark-mode toggle, mobile menu
│   └── utils.js            # Helpers (ids, shuffle, escaping)
├── data/
│   ├── app.json            # Worlds and their modules
│   ├── questions.json      # Quiz banks by category
│   ├── projects.json       # Project metadata (build-along)
│   ├── challenges.json     # Daily challenge pool
│   └── notes/
│       ├── <world>/*.md    # Lesson content (markdown)
│       └── projects/*.md   # Project build-along steps (markdown)
├── images/lessons/*.svg    # Lesson diagrams
└── tailwind.config.js      # Design tokens (brand palette, fonts, animations)
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

Content is data-driven — you can add worlds, lessons, questions, and projects
without touching the app code.

**Worlds** (`data/app.json`) — each world is a list of modules linked by
`next`/`previous`. On load, `storage.js` upserts them into IndexedDB by
`(path, title)`, so adding worlds or modules reaches existing users without
resetting their progress.

**Lessons** (`data/notes/<world>/<slug>.md`) — the file name is the module title
slugified (e.g. *"HTML Basics"* → `html_basics.md`). The lesson renderer supports
an extended Markdown:

- Headings, **bold**/*italic*/`code`, ordered & unordered lists, links
- Fenced code blocks: ` ```html … ``` ` / ` ```c ` / ` ```java ` / ` ```bash `
- Images → captioned figures: `![caption](/images/lessons/diagram.svg)`
- Callout boxes: `:::tip`, `:::analogy`, `:::warning`, `:::key`, `:::example`,
  `:::project` … `:::`
- Interactive exercises: `:::quiz`, `:::predict`, `:::reorder`, `:::fill`,
  `:::match` (mark a correct option with a trailing `*`; `E:` is the explanation
  shown on answer). For example:

  ```text
  :::quiz
  Q: What does HTML structure?
  - Style
  - Content *
  - Servers
  E: HTML describes the structure and meaning of content.
  :::
  ```

  All lesson input is escaped — no raw HTML is injected.

**Quiz banks** (`data/questions.json`) — questions grouped by category. New
categories/questions are merged in on load without duplicating existing ones.
A lesson in a single-topic world (C, Java, Linux) draws its end-of-lesson quiz
from the matching bank.

**Projects** (`data/projects.json` + `data/notes/projects/<id>.md`) — metadata
(title, emoji, world, tech, difficulty, step count) lives in the JSON; the build
itself is a markdown file split into `## Step N — Title` sections that the
build-along player renders one step at a time.

---

## Design system

- **Brand:** indigo/violet (`brand`) is the identity; **slate** is the neutral.
- **Colour means something:** green = success, amber = warning, red = error.
  Colour is not used decoratively.
- Tokens (palette, the Nunito display font, soft shadows, rounded corners,
  animations) live in `tailwind.config.js`; reusable component classes
  (`gd-card`, `gd-btn`, `gd-callout`, `qm-node`, `gd-terminal`, …) live in
  `css/input.css`.
- **Pixel**, the mascot, is drawn as inline SVG (`js/mascot.js`) so celebrations
  and briefings work offline with no image assets.

---

## Credits

Created by [Carlos S. Nah](https://twitter.com/rademejs). Built for beginners,
especially those learning with limited connectivity.

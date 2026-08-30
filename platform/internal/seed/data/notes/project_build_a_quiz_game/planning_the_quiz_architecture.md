# Planning the Quiz Architecture

Before writing code, plan the architecture. A quiz game is simple enough to build in one sitting, but the decisions you make about data, state, and rendering determine whether the code is maintainable or a tangled mess.

## The data model — questions as data, not HTML

The single most important decision: your questions are **data**, stored in a JavaScript array — not hard-coded HTML.

Each question is an object with:
- The question text
- An array of option strings
- The index of the correct option

```js
const questions = [
  {
    question: "What does CSS stand for?",
    options: [
      "Computer Style Sheets",
      "Cascading Style Sheets",
      "Creative Style System",
      "Colorful Style Sheets"
    ],
    correct: 1
  },
  {
    question: "Which HTML tag creates a hyperlink?",
    options: ["<link>", "<href>", "<a>", "<url>"],
    correct: 2
  },
  {
    question: "What property changes text color in CSS?",
    options: ["font-color", "text-color", "color", "foreground"],
    correct: 2
  }
];
```

:::key
Separating data from presentation is a fundamental engineering pattern. When questions live in a JavaScript array, you can add, remove, or reorder them without touching a single line of HTML or CSS. Later, you could even fetch them from an API.
:::

## The state model

State is the set of variables that describe what's happening in the app right now. For a quiz, you need exactly three:

```js
let currentIndex = 0;   // which question we're on
let score = 0;           // how many correct answers
let answered = false;    // has the user answered the current question?
```

These three values, combined with the `questions` array, tell you everything you need to render the correct screen.

:::tip
Before building any interactive JavaScript project, write out your state model first. Ask: "What variables do I need to fully describe the current state of this app?" The answer is usually simpler than you expect.
:::

## UI states — the two screens

The quiz has exactly two screens:

1. **Question screen** — shows the current question, options, progress, and a Next button
2. **Results screen** — shows the final score and a Restart button

```text
┌─────────────────────────────────┐
│  QUESTION SCREEN                │
│                                 │
│  Question 2 of 10               │
│  ┌─────────────────────────┐    │
│  │ What does CSS stand for?│    │
│  └─────────────────────────┘    │
│                                 │
│  [ Computer Style Sheets  ]     │
│  [ Cascading Style Sheets ]     │
│  [ Creative Style System  ]     │
│  [ Colorful Style Sheets  ]     │
│                                 │
│          [Next →]               │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│  RESULTS SCREEN                 │
│                                 │
│  🎉 You scored 8 / 10          │
│                                 │
│  That's 80% — Great job!        │
│                                 │
│        [Restart Quiz]           │
└─────────────────────────────────┘
```

At any given moment, only one screen is visible. You toggle between them based on whether `currentIndex < questions.length`.

## The state → render loop

This is the core pattern that drives the entire app:

1. **Something happens** (user clicks an option, clicks Next, clicks Restart)
2. **Update the state** (increment `score`, advance `currentIndex`, set `answered`)
3. **Re-render the UI** based on the new state

```text
User Action → Update State → Render UI
     ↑                          │
     └──────────────────────────┘
```

```js
function render() {
  if (currentIndex < questions.length) {
    showQuestionScreen();
  } else {
    showResultsScreen();
  }
}
```

Every time state changes, call `render()`. The render function reads the current state and updates the DOM accordingly. This keeps your logic clean — you never manually show/hide individual elements scattered across different event handlers.

:::key
The state → render loop is the foundation of modern UI development. React, Vue, and Svelte all use this pattern. Learning it now with vanilla JavaScript gives you a head start on every framework.
:::

## File structure

Keep it simple — three files:

```text
quiz-game/
├── index.html     ← structure and both screen containers
├── style.css      ← layout, cards, buttons, animations
└── script.js      ← data, state, event handlers, render logic
```

### index.html — the skeleton

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Quiz Game</title>
  <link rel="stylesheet" href="style.css" />
</head>
<body>
  <div class="quiz-container">
    <!-- Question screen -->
    <div id="quiz-screen">
      <p id="progress"></p>
      <h2 id="question"></h2>
      <div id="options"></div>
      <button id="next-btn" disabled>Next</button>
    </div>

    <!-- Results screen (hidden initially) -->
    <div id="results-screen" class="hidden">
      <h2 id="result-message"></h2>
      <p id="score-display"></p>
      <button id="restart-btn">Restart Quiz</button>
    </div>
  </div>

  <script src="script.js"></script>
</body>
</html>
```

Both screens exist in the HTML. JavaScript controls which one is visible at any time.

### The hidden utility class

```css
.hidden {
  display: none;
}
```

Simple, reusable, and clear. Toggle visibility by adding or removing this class.

## Planning the JavaScript flow

Before writing the full implementation, map out the functions you'll need:

```js
// Data
const questions = [ /* ... */ ];

// State
let currentIndex = 0;
let score = 0;
let answered = false;

// Core functions
function renderQuestion()   { /* show current question and options */ }
function handleAnswer(index) { /* check answer, update score, show feedback */ }
function nextQuestion()     { /* advance index, re-render */ }
function showResults()      { /* display final score */ }
function restartQuiz()      { /* reset state, re-render */ }

// Event listeners
document.getElementById('next-btn').addEventListener('click', nextQuestion);
document.getElementById('restart-btn').addEventListener('click', restartQuiz);

// Start
renderQuestion();
```

:::tip
Planning your functions before writing their bodies is like sketching a wireframe before writing CSS. It clarifies the structure and prevents you from solving everything inside one giant event handler.
:::

## Why this architecture matters

This three-file, state-driven architecture might seem like overkill for a small quiz. It isn't. Here's what it gives you:

- **Easy to extend** — adding a timer, categories, or difficulty levels means adding state variables and updating the render function. The structure doesn't change.
- **Easy to debug** — if the wrong question shows, check `currentIndex`. If the score is wrong, check `handleAnswer`. Each function has a clear job.
- **Transferable** — this exact pattern (data + state + render loop) is how React, Vue, and every modern framework works.

:::quiz
Q: Why store quiz questions in a JavaScript array instead of writing them directly in HTML?
- Arrays make the page load faster
- Separating data from presentation makes it easy to add, remove, or reorder questions without touching HTML *
- JavaScript arrays are the only way to create buttons
- HTML doesn't support question marks
E: When data lives in JavaScript, you can change it independently of the presentation layer. This makes the app flexible — you could even load questions from an API later without changing any HTML.
:::

:::quiz
Q: What three state variables are needed to describe the current state of the quiz?
- question, answer, result
- currentIndex, score, and answered *
- html, css, and js
- startTime, endTime, and duration
E: `currentIndex` tracks which question we're on, `score` tracks correct answers, and `answered` tracks whether the user has selected an option for the current question. These three values fully describe the quiz state.
:::

## Recap

- **Data model:** questions are an array of objects, each with `question`, `options`, and `correct` index.
- **State model:** three variables — `currentIndex`, `score`, `answered`.
- **Two UI states:** question screen and results screen. Only one is visible at a time.
- **State → render loop:** user action → update state → re-render the UI.
- **File structure:** `index.html`, `style.css`, `script.js` — keep it simple.
- **Plan functions first:** `renderQuestion`, `handleAnswer`, `nextQuestion`, `showResults`, `restartQuiz`.
- This architecture is the same pattern used by modern frameworks.

**Next up:** Building the Question Engine — implementing `renderQuestion` and generating the UI from data.

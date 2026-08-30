# Code Review and Extensions

Your quiz game works, looks polished, and handles keyboard input. Before shipping, step back and review the code like a professional — then explore how to extend it into an even stronger portfolio piece.

## Review the code structure

A well-structured quiz game should have clear separation:

```text
script.js
├── Data:    questions array
├── State:   currentIndex, score, answered
├── DOM:     cached element references
├── Core:    renderQuestion, handleAnswer, nextQuestion
├── Results: showResults, getMessage, restartQuiz
├── Events:  click and keyboard listeners
└── Init:    renderQuestion() call
```

Open your `script.js` and check:

- **Is data at the top?** The questions array should be the first thing in the file.
- **Is state clearly grouped?** The three state variables should sit together, easy to find.
- **Are DOM references cached once?** No `document.getElementById` calls inside functions that run repeatedly.
- **Does each function do one thing?** `renderQuestion` renders, `handleAnswer` handles answers, `showResults` shows results.

:::key
Good code structure means a stranger could open your file, read the first 20 lines, and understand the architecture. Data → state → DOM → functions → events → init.
:::

## Refactor opportunities

After a first pass, most quiz game code has a few things worth cleaning up.

### Extract magic numbers

```js
// Before: what does 15 mean?
if (timeLeft <= 0) { ... }

// After: named constant
const SECONDS_PER_QUESTION = 15;
if (timeLeft <= 0) { ... }
```

### Replace innerHTML with DOM methods where possible

```js
// innerHTML is fine for clearing:
optionsEl.innerHTML = '';

// But for building elements, createElement is safer and avoids XSS risks:
const button = document.createElement('button');
button.textContent = option;  // textContent is safe; innerHTML is not
```

### Group related DOM references

```js
// Before: scattered variables
const quizScreen = document.getElementById('quiz-screen');
const resultsScreen = document.getElementById('results-screen');
const questionEl = document.getElementById('question');

// After: organized by screen
const dom = {
  quiz: {
    screen: document.getElementById('quiz-screen'),
    progress: document.getElementById('progress'),
    question: document.getElementById('question'),
    options: document.getElementById('options'),
    nextBtn: document.getElementById('next-btn'),
  },
  results: {
    screen: document.getElementById('results-screen'),
    message: document.getElementById('result-message'),
    score: document.getElementById('score-display'),
    restartBtn: document.getElementById('restart-btn'),
  }
};
```

:::tip
Refactoring doesn't change what the code does — it changes how clearly it communicates. If your code already works, refactoring is about making it easier to read and extend.
:::

### Encapsulate state in an object

```js
// Before
let currentIndex = 0;
let score = 0;
let answered = false;

// After
const state = {
  currentIndex: 0,
  score: 0,
  answered: false,
  reset() {
    this.currentIndex = 0;
    this.score = 0;
    this.answered = false;
  }
};
```

Now `restartQuiz` just calls `state.reset()` instead of manually resetting each variable.

## Extension ideas

Each of these adds real functionality and demonstrates a new skill to potential employers.

### 1. Timer per question

Add a countdown timer that auto-advances when time runs out.

```js
const SECONDS_PER_QUESTION = 15;
let timeLeft = SECONDS_PER_QUESTION;
let timer = null;

function startTimer() {
  timeLeft = SECONDS_PER_QUESTION;
  updateTimerDisplay();

  timer = setInterval(() => {
    timeLeft--;
    updateTimerDisplay();
    if (timeLeft <= 0) {
      clearInterval(timer);
      handleAnswer(-1);  // time's up = wrong
    }
  }, 1000);
}

function stopTimer() {
  clearInterval(timer);
}

function updateTimerDisplay() {
  timerEl.textContent = `${timeLeft}s`;
  timerEl.classList.toggle('warning', timeLeft <= 5);
}
```

```css
.timer.warning {
  color: #ef4444;
  animation: pulse 0.5s ease infinite alternate;
}

@keyframes pulse {
  to { transform: scale(1.1); }
}
```

### 2. Question categories

Group questions by topic and let the user choose.

```js
const questionBank = {
  html: [
    { question: "What does HTML stand for?", options: [...], correct: 0 },
    // ...
  ],
  css: [
    { question: "Which property sets text color?", options: [...], correct: 2 },
    // ...
  ],
  javascript: [
    { question: "What does === check?", options: [...], correct: 1 },
    // ...
  ]
};
```

### 3. Difficulty levels

Tag questions with difficulty and let users choose easy, medium, or hard.

```js
const questions = [
  { question: "...", options: [...], correct: 0, difficulty: "easy" },
  { question: "...", options: [...], correct: 2, difficulty: "hard" },
];

const filtered = questions.filter(q => q.difficulty === selectedDifficulty);
```

### 4. High scores with localStorage

Save the best scores so they persist across sessions.

```js
function saveHighScore(score, total) {
  const percentage = Math.round((score / total) * 100);
  const scores = JSON.parse(localStorage.getItem('quizHighScores') || '[]');

  scores.push({
    score: percentage,
    date: new Date().toLocaleDateString(),
    total: total
  });

  // Keep only top 5
  scores.sort((a, b) => b.score - a.score);
  scores.splice(5);

  localStorage.setItem('quizHighScores', JSON.stringify(scores));
}

function loadHighScores() {
  return JSON.parse(localStorage.getItem('quizHighScores') || '[]');
}
```

:::warning
`localStorage` only stores strings. Always `JSON.stringify` when saving and `JSON.parse` when loading. Forgetting this is a classic bug — you'll save `"[object Object]"` instead of your actual data.
:::

### 5. Fetch questions from Open Trivia API

Replace your hard-coded questions with live data from opentdb.com:

```js
async function fetchQuestions(amount = 10, category = 18) {
  const url = `https://opentdb.com/api.php?amount=${amount}&category=${category}&type=multiple`;
  const response = await fetch(url);
  const data = await response.json();

  return data.results.map(q => {
    const options = [...q.incorrect_answers, q.correct_answer];
    // Shuffle options
    options.sort(() => Math.random() - 0.5);
    return {
      question: decodeHTML(q.question),
      options: options.map(decodeHTML),
      correct: options.indexOf(q.correct_answer)
    };
  });
}

function decodeHTML(text) {
  const textarea = document.createElement('textarea');
  textarea.innerHTML = text;
  return textarea.value;
}
```

Category 18 is "Computer Science." This extension demonstrates async/await, API integration, and data transformation — all highly valued skills.

### 6. Shuffle questions

Randomize question order each time:

```js
function shuffle(array) {
  const shuffled = [...array];
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
  }
  return shuffled;
}

// At quiz start:
const activeQuestions = shuffle(questions);
```

## What makes this a strong portfolio project

A quiz game stands out on a portfolio when:

- **It's deployed and working** — live demo link, no console errors
- **It demonstrates state management** — not just DOM manipulation, but a clear data → state → render pattern
- **The code is clean** — readable function names, cached DOM references, no spaghetti
- **It has polish** — progress bar, animations, keyboard support
- **You can explain the architecture** — talk about the state model, the render loop, and why you structured it the way you did
- **It has at least one extension** — a timer, localStorage scores, or API-fetched questions show initiative

:::tip
When writing your portfolio case study for this project, lead with the architecture, not the feature list. "I used a state → render loop to manage quiz state" is more impressive than "it has 10 questions."
:::

## Final code checklist

Before shipping:

- [ ] All questions have 4 options with exactly 1 correct answer
- [ ] Score calculates correctly (test: answer all right, answer all wrong)
- [ ] Restart fully resets (score, progress bar, question index)
- [ ] Keyboard shortcuts work (1-4 and Enter)
- [ ] No console errors at any point during the quiz
- [ ] Responsive on mobile (buttons tappable, text readable)
- [ ] Focus styles visible on all interactive elements
- [ ] Progress bar fills correctly and hits 100% at results

:::quiz
Q: Why is separating data (questions array) from rendering logic a good practice?
- It makes the quiz run faster
- It allows you to change, add, or fetch questions without modifying the rendering code *
- Data arrays use less memory than HTML
- The browser requires this separation
E: When data and rendering are separate, you can swap the data source (hard-coded array, API, localStorage) without touching the render functions. This separation of concerns makes code more maintainable and extensible.
:::

:::quiz
Q: What skill does adding localStorage high scores demonstrate to a potential employer?
- CSS animation expertise
- The ability to persist data across browser sessions using browser storage APIs *
- Server-side database management
- Advanced algorithm design
E: Using localStorage shows you understand client-side data persistence — how to save, retrieve, parse, and manage data that survives page refreshes. It's a practical skill used in real applications for preferences, caches, and offline functionality.
:::

## Recap

- **Review your structure:** data → state → DOM refs → functions → events → init.
- **Refactor** for clarity: extract constants, group DOM references, encapsulate state.
- **Extensions** add real value: timer, categories, difficulty, localStorage scores, API-fetched questions.
- **Portfolio-worthy** means deployed, clean, polished, and explainable.
- Lead your case study with **architecture and decisions**, not just features.

**Congratulations** — you've built a complete, interactive quiz game from scratch. The architecture you learned here — data-driven rendering with a state → render loop — is the foundation of modern frontend development.

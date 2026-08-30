# Building the Question Engine

The architecture is planned. Now you'll implement the question engine — the core of the quiz that reads from the data array and renders each question to the screen.

## Creating the questions array

Start `script.js` with your question data. Use at least 5–8 questions so the quiz feels substantial.

```js
const questions = [
  {
    question: "What does HTML stand for?",
    options: [
      "Hyper Text Markup Language",
      "High Tech Modern Language",
      "Hyper Transfer Markup Language",
      "Home Tool Markup Language"
    ],
    correct: 0
  },
  {
    question: "Which CSS property controls text size?",
    options: ["text-size", "font-style", "font-size", "text-style"],
    correct: 2
  },
  {
    question: "What does the === operator check in JavaScript?",
    options: [
      "Value only",
      "Value and type",
      "Type only",
      "Neither value nor type"
    ],
    correct: 1
  },
  {
    question: "Which HTML element is used for the largest heading?",
    options: ["<heading>", "<h6>", "<head>", "<h1>"],
    correct: 3
  },
  {
    question: "What does DOM stand for?",
    options: [
      "Document Object Model",
      "Display Output Manager",
      "Digital Ordinance Map",
      "Document Order Method"
    ],
    correct: 0
  }
];
```

:::tip
Always use four options per question for consistency. It keeps the layout predictable and the probability of guessing correct at 25% — hard enough to feel meaningful.
:::

## Initialize the state

Below your questions array, set up the state variables:

```js
let currentIndex = 0;
let score = 0;
let answered = false;
```

And grab references to the DOM elements you'll update:

```js
const quizScreen    = document.getElementById('quiz-screen');
const resultsScreen = document.getElementById('results-screen');
const progressEl    = document.getElementById('progress');
const questionEl    = document.getElementById('question');
const optionsEl     = document.getElementById('options');
const nextBtn       = document.getElementById('next-btn');
```

:::key
Cache your DOM references once at the top of the file. Calling `document.getElementById` every time you render is wasteful and clutters your functions. Store the reference once and reuse it.
:::

## The renderQuestion function

This is the heart of the question engine. It reads the current question from the array and builds the UI.

```js
function renderQuestion() {
  const q = questions[currentIndex];

  // Update progress text
  progressEl.textContent = `Question ${currentIndex + 1} of ${questions.length}`;

  // Set the question text
  questionEl.textContent = q.question;

  // Clear previous options
  optionsEl.innerHTML = '';

  // Generate option buttons from data
  q.options.forEach((option, index) => {
    const button = document.createElement('button');
    button.classList.add('option-btn');
    button.textContent = option;
    button.addEventListener('click', () => handleAnswer(index));
    optionsEl.appendChild(button);
  });

  // Reset state for this question
  answered = false;
  nextBtn.disabled = true;
}
```

Walk through what this does:

1. **Reads the current question** from the array using `currentIndex`
2. **Updates the progress indicator** — "Question 1 of 5"
3. **Sets the question heading** text
4. **Clears old option buttons** — essential when advancing to a new question
5. **Creates new buttons** from the options array, attaching a click handler to each
6. **Resets state** — the user hasn't answered yet, so Next is disabled

## Generating option buttons from data

The key technique here is building DOM elements dynamically from an array. Instead of hard-coding four `<button>` elements in HTML, you create them in JavaScript:

```js
q.options.forEach((option, index) => {
  const button = document.createElement('button');
  button.classList.add('option-btn');
  button.textContent = option;
  button.addEventListener('click', () => handleAnswer(index));
  optionsEl.appendChild(button);
});
```

This approach means:
- Adding a fifth option to a question just works — no HTML changes needed
- Each button knows its own index, passed to `handleAnswer`
- The UI always matches the data

:::warning
Always clear `optionsEl.innerHTML = ''` before generating new buttons. If you forget, the previous question's options stay on screen and pile up.
:::

## Displaying progress — "Question 1 of 5"

Progress tells the user where they are. It's one line but makes the experience feel structured:

```js
progressEl.textContent = `Question ${currentIndex + 1} of ${questions.length}`;
```

We add 1 to `currentIndex` because arrays are 0-indexed but humans count from 1. `questions.length` gives the total automatically — no hard-coded number to update when you add questions.

## The Next button — disabled until answered

The Next button should be disabled when a new question appears and enabled only after the user selects an option.

In the HTML:

```html
<button id="next-btn" disabled>Next</button>
```

In `renderQuestion`:

```js
nextBtn.disabled = true;
```

After the user answers (covered in the next lesson):

```js
nextBtn.disabled = false;
```

## The CSS for the question screen

Style the quiz container, question, and option buttons:

```css
.quiz-container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
  font-family: system-ui, sans-serif;
}

#progress {
  font-size: 0.875rem;
  color: #64748b;
  margin-bottom: 0.5rem;
}

#question {
  font-size: 1.375rem;
  margin-bottom: 1.5rem;
  line-height: 1.3;
}

.option-btn {
  display: block;
  width: 100%;
  padding: 0.875rem 1rem;
  margin-bottom: 0.75rem;
  background: #f1f5f9;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.option-btn:hover {
  border-color: #3b82f6;
  background: #eff6ff;
}

#next-btn {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}

#next-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.hidden {
  display: none;
}
```

## Advancing to the next question

The `nextQuestion` function advances `currentIndex` and either renders the next question or shows results:

```js
function nextQuestion() {
  currentIndex++;

  if (currentIndex < questions.length) {
    renderQuestion();
  } else {
    showResults();
  }
}

nextBtn.addEventListener('click', nextQuestion);
```

This is the state → render loop in action: update the state (`currentIndex++`), then render based on the new state.

## Starting the quiz

At the bottom of `script.js`, kick everything off:

```js
renderQuestion();
```

That single call reads `currentIndex` (which is 0), renders the first question, and the quiz is running.

:::quiz
Q: Why do we use `forEach` to generate option buttons instead of writing them directly in HTML?
- forEach is faster than HTML
- Generating buttons from data means the UI always matches the questions array, and adding options requires no HTML changes *
- HTML doesn't support buttons
- forEach automatically styles the buttons
E: When you generate UI from data, the two can never go out of sync. Add a question to the array and the button appears automatically. This is the data-driven rendering pattern used in every modern framework.
:::

:::quiz
Q: Why is the Next button disabled when a new question renders?
- To prevent the user from skipping ahead without answering *
- Because disabled buttons load faster
- The browser requires it for accessibility
- To save memory
E: Disabling the Next button forces the user to select an answer before advancing. This ensures every question is attempted and the score is meaningful.
:::

## Recap

- The **questions array** stores all quiz data as objects with `question`, `options`, and `correct`.
- **State variables** (`currentIndex`, `score`, `answered`) are initialized at the top.
- **DOM references** are cached once for reuse.
- **`renderQuestion()`** reads the current question, builds option buttons dynamically, updates progress, and resets state.
- **Buttons are generated from data** using `forEach` and `createElement` — never hard-coded.
- **Progress** uses `currentIndex + 1` for human-friendly numbering.
- The **Next button** starts disabled and is enabled after answering.
- **`nextQuestion()`** advances the index and either re-renders or shows results.

**Next up:** Handling Answers and Scoring — detecting clicks, checking correctness, and giving visual feedback.

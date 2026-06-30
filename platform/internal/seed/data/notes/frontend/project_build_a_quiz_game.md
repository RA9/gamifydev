# Project: Build a Quiz Game

Time to make something that *does* things. In this guided project you'll build a working multiple-choice quiz game with HTML, CSS, and JavaScript: questions render one at a time, clicking an answer marks it correct or wrong, the score updates, and a final screen lets you restart.

:::project
**Goal:** Build a quiz that stores questions as data, shows one question with clickable options, gives instant feedback, tracks the score, advances with a Next button, and shows a final score screen with a Restart button. Everything lives in `index.html`, `style.css`, and `script.js`.
:::

## How the app works

The whole game is driven by a few pieces of state:

- A `questions` **array** holding all the questions and answers.
- A `current` number tracking which question we're on.
- A `score` number counting correct answers.

The JavaScript reads this state, draws the right screen, and updates the state when you click. This "data drives the display" pattern is the foundation of almost every interactive app.

## Step 1: The HTML

We need a container for the quiz screen and a separate (hidden) container for the results screen. JavaScript will fill these in.

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
  <main class="quiz">
    <!-- Question screen -->
    <div id="quiz-screen">
      <div class="quiz-top">
        <span id="progress">Question 1 of 4</span>
        <span id="score">Score: 0</span>
      </div>

      <h1 id="question">Question text goes here</h1>

      <div id="options" class="options"></div>

      <button id="next-btn" class="btn" disabled>Next</button>
    </div>

    <!-- Results screen (hidden until the end) -->
    <div id="result-screen" class="hidden">
      <h1>Quiz complete!</h1>
      <p id="final-score">You scored 0 / 4</p>
      <button id="restart-btn" class="btn">Restart</button>
    </div>
  </main>

  <script src="script.js"></script>
</body>
</html>
```

The `<div id="options">` is empty on purpose — JavaScript builds the answer buttons. The Next button starts `disabled` so players must answer before moving on.

## Step 2: The CSS

Style the card, the option buttons, and define the helper classes JavaScript will toggle: `.correct`, `.wrong`, and `.hidden`.

```css
:root {
  --brand: #2d8cff;
  --brand-dark: #1f6fd4;
  --ink: #1c2030;
  --muted: #6b7280;
  --bg: #eef2fb;
  --card: #ffffff;
  --good: #16a34a;
  --bad: #dc2626;
  --line: #e3e8f3;
  --radius: 14px;
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: system-ui, sans-serif;
  background: var(--bg);
  color: var(--ink);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.quiz {
  width: 440px;
  max-width: 100%;
  background: var(--card);
  border-radius: var(--radius);
  padding: 1.75rem;
  box-shadow: 0 16px 40px rgba(28, 32, 48, 0.12);
}

.quiz-top {
  display: flex;
  justify-content: space-between;
  color: var(--muted);
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

#question {
  font-size: 1.3rem;
  margin-bottom: 1.25rem;
}

.options {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  margin-bottom: 1.25rem;
}

.option {
  text-align: left;
  padding: 0.85rem 1rem;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #fff;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.15s ease;
}

.option:hover:not(:disabled) {
  border-color: var(--brand);
  background: #f5f9ff;
}

.option:disabled {
  cursor: default;
}

.option.correct {
  border-color: var(--good);
  background: #e9faf0;
  color: var(--good);
  font-weight: 600;
}

.option.wrong {
  border-color: var(--bad);
  background: #fdecec;
  color: var(--bad);
  font-weight: 600;
}

.btn {
  width: 100%;
  padding: 0.85rem;
  border: none;
  border-radius: 10px;
  background: var(--brand);
  color: #fff;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease;
}

.btn:hover:not(:disabled) { background: var(--brand-dark); }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }

#result-screen { text-align: center; }
#final-score { font-size: 1.2rem; color: var(--muted); margin: 1rem 0 1.5rem; }

.hidden { display: none; }
```

## Step 3: The data

Open `script.js`. First, the **questions array** — an array of objects. Each object has the question text, an array of options, and the index of the correct option.

```js
const questions = [
  {
    question: "What does HTML stand for?",
    options: [
      "Hyper Text Markup Language",
      "High Tech Modern Language",
      "Home Tool Markup Language",
      "Hyperlink Text Marking Language",
    ],
    correct: 0,
  },
  {
    question: "Which language styles a web page?",
    options: ["HTML", "CSS", "Python", "SQL"],
    correct: 1,
  },
  {
    question: "Which symbol starts an ID selector in CSS?",
    options: [".", "#", "*", "@"],
    correct: 1,
  },
  {
    question: "What does a function let you do?",
    options: [
      "Store a single number",
      "Reuse a block of instructions",
      "Change the page color",
      "Connect to the internet",
    ],
    correct: 1,
  },
];
```

:::key
Storing your content as **data** (an array of objects) separates *what* the quiz asks from *how* it's displayed. To add a question, you just add an object — you never touch the rendering code.
:::

## Step 4: Grab elements and set up state

Get references to the HTML elements once, and declare the state variables.

```js
const quizScreen = document.getElementById("quiz-screen");
const resultScreen = document.getElementById("result-screen");
const questionEl = document.getElementById("question");
const optionsEl = document.getElementById("options");
const progressEl = document.getElementById("progress");
const scoreEl = document.getElementById("score");
const nextBtn = document.getElementById("next-btn");
const finalScoreEl = document.getElementById("final-score");
const restartBtn = document.getElementById("restart-btn");

let current = 0;
let score = 0;
let answered = false;
```

## Step 5: Render the current question

This function clears the old options and builds fresh buttons for the current question. Each button gets a click handler.

```js
function showQuestion() {
  answered = false;
  nextBtn.disabled = true;

  const q = questions[current];
  questionEl.textContent = q.question;
  progressEl.textContent = `Question ${current + 1} of ${questions.length}`;
  scoreEl.textContent = `Score: ${score}`;

  optionsEl.innerHTML = ""; // clear previous options

  q.options.forEach((optionText, index) => {
    const button = document.createElement("button");
    button.textContent = optionText;
    button.className = "option";
    button.addEventListener("click", () => selectAnswer(index, button));
    optionsEl.appendChild(button);
  });
}
```

We use `forEach` (a loop) to create one button per option. `createElement` and `appendChild` add real elements to the page.

## Step 6: Handle an answer click

When a player clicks an option we: stop further answers, mark the chosen button correct or wrong, reveal the right answer, update the score, and enable Next.

```js
function selectAnswer(index, button) {
  if (answered) return; // ignore extra clicks
  answered = true;

  const q = questions[current];
  const buttons = optionsEl.querySelectorAll(".option");

  buttons.forEach((b) => (b.disabled = true)); // lock all options

  if (index === q.correct) {
    button.classList.add("correct");
    score++;
  } else {
    button.classList.add("wrong");
    buttons[q.correct].classList.add("correct"); // show the right one
  }

  scoreEl.textContent = `Score: ${score}`;
  nextBtn.disabled = false;
}
```

:::analogy
Adding `.correct` or `.wrong` is like a teacher marking your test with a green tick or a red cross. The JavaScript decides the mark; the CSS draws it.
:::

## Step 7: Next button and the results screen

The Next button advances to the next question, or shows the final screen if we've run out.

```js
nextBtn.addEventListener("click", () => {
  current++;
  if (current < questions.length) {
    showQuestion();
  } else {
    showResults();
  }
});

function showResults() {
  quizScreen.classList.add("hidden");
  resultScreen.classList.remove("hidden");
  finalScoreEl.textContent = `You scored ${score} / ${questions.length}`;
}
```

Showing and hiding screens is just toggling the `.hidden` class, which sets `display: none`.

## Step 8: Restart and kick things off

Restart resets the state and shows the first question again. The final line starts the quiz when the page loads.

```js
restartBtn.addEventListener("click", () => {
  current = 0;
  score = 0;
  resultScreen.classList.add("hidden");
  quizScreen.classList.remove("hidden");
  showQuestion();
});

// Start the quiz
showQuestion();
```

That's a complete, working quiz game. Open `index.html` and play through it.

:::quiz
Q: Why is each question stored as an object with a `correct` index instead of hard-coding answers in the HTML?
- It makes the page load faster
- It separates content from display, so adding questions never requires changing the rendering code *
- The browser requires it
E: Keeping content as data means the same rendering code works for any number of questions. To add or edit a question, you only change the data array.
:::

:::predict
A player clicks the same wrong option twice quickly. Because of the `if (answered) return;` line at the top of `selectAnswer`, what happens on the second click?
ANSWER: Nothing — the function returns immediately because `answered` is already `true`, so the score and styling don't change again.
:::

## Challenge extensions

Level up your quiz:

1. **Add a timer.** Give each question 15 seconds with `setInterval`. If time runs out, auto-reveal the answer and disable the options.
2. **Shuffle questions.** Randomize the `questions` array on start (and shuffle each question's options) so no two playthroughs are identical.
3. **Categories.** Add a `category` field to each question and a start screen where the player picks a topic, then filter the array.
4. **Progress bar.** Replace the "Question X of Y" text with a visual bar that fills as the player advances.
5. **High score.** Save the best score in `localStorage` so it persists after the page reloads.
6. **Review screen.** On the results screen, list each question with the player's answer marked right or wrong.

## Recap

- Stored quiz content as a **`questions` array of objects** (text, options, correct index).
- Tracked game **state** with `current`, `score`, and `answered` variables.
- Rendered each question by **looping** over options and creating buttons with `createElement`.
- Handled clicks to **mark correct/wrong**, reveal the answer, lock options, and update the score.
- Used a **Next button** to advance and a **results screen** toggled with the `.hidden` class.
- Added a **Restart** that resets state and replays — and a list of challenges to extend it.

**Next up:** Keep building. Combine everything — HTML structure, CSS styling, and JavaScript behaviour — into your own projects.

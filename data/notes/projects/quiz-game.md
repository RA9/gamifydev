# Quiz Game

Quizzes are a perfect first taste of "real" JavaScript: data, rendering, clicks, and state all working together. You'll build a 3-question multiple-choice quiz that tallies a score and lets you play again.

:::project
You'll have a single HTML file that shows one question at a time with clickable answer options, tracks how many you got right, and displays a final score with a Restart button.
:::

To follow along you only need a text editor and a web browser. Put everything in one file called `quiz.html` — markup, a `<style>` block, and a `<script>` block — saving and refreshing as you go.

## Step 1 — Lay out the HTML structure

Build the shell of the quiz: a container with a spot for the question, an empty options area we'll fill with JavaScript, a Next button, and a hidden score area for the end.

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Quiz Game</title>
</head>
<body>
  <div class="quiz">
    <h1 id="question">Question goes here</h1>
    <div id="options"></div>
    <button id="next" disabled>Next</button>
    <div id="score" class="hidden"></div>
  </div>
</body>
</html>
```

## Step 2 — Add some basic styling

Add a `<style>` block in the `<head>` so the quiz is centered and the option buttons look tappable. The `.hidden` class lets us show and hide elements from JavaScript.

```html
<style>
  body {
    margin: 0; min-height: 100vh;
    display: flex; justify-content: center; align-items: center;
    background: #0f172a; font-family: system-ui, sans-serif;
  }
  .quiz {
    background: #1e293b; color: #e2e8f0;
    padding: 32px; border-radius: 16px; width: 340px;
  }
  #options button {
    display: block; width: 100%; margin: 8px 0; padding: 12px;
    background: #334155; color: #e2e8f0; border: none;
    border-radius: 8px; cursor: pointer; text-align: left; font-size: 15px;
  }
  #options button:hover { background: #475569; }
  #next {
    margin-top: 12px; padding: 10px 24px; border: none; border-radius: 8px;
    background: #6366f1; color: #fff; font-weight: 600; cursor: pointer;
  }
  #next:disabled { opacity: 0.5; cursor: not-allowed; }
  .hidden { display: none; }
</style>
```

## Step 3 — Define the questions as data

Open a `<script>` block right before `</body>`. Store your quiz as an array of objects — each with the question text, a list of options, and the index of the correct answer. Keeping data separate from display logic is a habit that scales to big apps.

```html
<script>
  const questions = [
    {
      q: "What does HTML stand for?",
      options: ["Hyper Text Markup Language", "Hot Mail", "How To Make Lasagna"],
      answer: 0
    },
    {
      q: "Which symbol starts a CSS id selector?",
      options: [".", "#", "*"],
      answer: 1
    },
    {
      q: "Which keyword declares a constant in JavaScript?",
      options: ["var", "let", "const"],
      answer: 2
    }
  ];

  let current = 0;
  let score = 0;
</script>
```

## Step 4 — Render the current question

Inside the same `<script>`, grab the elements and write a `render` function that fills in the question text and creates a button for each option. Call it once to show the first question.

```js
const questionEl = document.getElementById("question");
const optionsEl = document.getElementById("options");
const nextBtn = document.getElementById("next");
const scoreEl = document.getElementById("score");

function render() {
  const item = questions[current];
  questionEl.textContent = item.q;
  optionsEl.innerHTML = "";
  nextBtn.disabled = true;

  item.options.forEach((text, index) => {
    const btn = document.createElement("button");
    btn.textContent = text;
    btn.onclick = () => selectAnswer(index, btn);
    optionsEl.appendChild(btn);
  });
}

render();
```

## Step 5 — Handle clicks and track the score

Now react when an option is clicked: lock the options so they can't change their answer, add to the score if they were right, and enable the Next button to move on.

```js
function selectAnswer(index, btn) {
  const correct = questions[current].answer;
  const buttons = optionsEl.querySelectorAll("button");
  buttons.forEach(b => b.disabled = true);

  if (index === correct) {
    score++;
    btn.style.background = "#16a34a";
  } else {
    btn.style.background = "#dc2626";
    buttons[correct].style.background = "#16a34a";
  }

  nextBtn.disabled = false;
}
```

## Step 6 — Show the final score and restart

Wire up the Next button. If there are more questions, advance and re-render. If not, hide the quiz controls and reveal the final score with a Restart button that resets everything.

```js
nextBtn.onclick = () => {
  current++;
  if (current < questions.length) {
    render();
  } else {
    questionEl.classList.add("hidden");
    optionsEl.classList.add("hidden");
    nextBtn.classList.add("hidden");
    scoreEl.classList.remove("hidden");
    scoreEl.innerHTML =
      `<h2>You scored ${score} / ${questions.length}</h2>
       <button id="restart">Play again</button>`;
    document.getElementById("restart").onclick = () => location.reload();
  }
};
```

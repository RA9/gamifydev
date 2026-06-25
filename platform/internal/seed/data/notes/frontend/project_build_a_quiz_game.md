# Project: Build a Quiz Game

You've learned how JavaScript reacts to clicks and updates the page. Let's prove it by building a real, working **quiz game** — the same kind of interaction this very app uses.

:::project
A multiple-choice quiz: show a question and options, let the player pick, tell them if they're right, keep score, and move to the next question. Pure HTML, CSS, and JavaScript in one file.
:::

## Step 1 — The HTML shell

```html
<div id="game">
  <h2 id="question"></h2>
  <div id="options"></div>
  <p id="score">Score: 0</p>
</div>
```

Three empty slots — for the question, the options, and the score. JavaScript will fill them in.

## Step 2 — The data

A quiz is just a list of questions. We'll store them as an array of objects:

```js
const questions = [
  { q: "What does HTML structure?", options: ["Style", "Content", "Servers"], answer: "Content" },
  { q: "Which language adds behaviour?", options: ["CSS", "HTML", "JavaScript"], answer: "JavaScript" },
  { q: "What centers items easily?", options: ["Flexbox", "Tables", "Borders"], answer: "Flexbox" }
];

let current = 0;
let score = 0;
```

`current` tracks which question we're on; `score` counts correct answers. This is **state** — your data, separate from the display.

## Step 3 — Show a question

This function reads the current question and builds a button for each option:

```js
function showQuestion() {
  const q = questions[current];
  document.getElementById("question").textContent = q.q;

  const box = document.getElementById("options");
  box.innerHTML = ""; // clear old options

  q.options.forEach(function (option) {
    const btn = document.createElement("button");
    btn.textContent = option;
    btn.addEventListener("click", function () { checkAnswer(option); });
    box.appendChild(btn);
  });
}
```

See the **select → listen → update** pattern? We create each button, listen for its click, and add it to the page.

## Step 4 — Check the answer

```js
function checkAnswer(choice) {
  if (choice === questions[current].answer) {
    score++;
    document.getElementById("score").textContent = "Score: " + score;
  }
  current++;
  if (current < questions.length) {
    showQuestion();
  } else {
    document.getElementById("game").innerHTML =
      "<h2>Done! You scored " + score + " / " + questions.length + " 🎉</h2>";
  }
}
```

If the choice matches the answer, bump the score. Then move on — or, when there are no questions left, show the final result.

## Step 5 — Start the game

One line kicks everything off:

```js
showQuestion();
```

Open the file and play. You built a quiz game! 🎮

:::quiz
Q: In this project, what does the `score` variable represent?
- The HTML structure of the page
- The game's state — data the display is built from *
- A CSS style rule
E: `score` (and `current`) is *state*: the data your game tracks. The page is re-rendered from that state each step — the core idea behind every interactive app.
:::

## You built it! 🚀

This is a genuine interactive app: data, events, score-keeping, and a finish screen. The select → listen → update loop you used here powers everything from to-do lists to games.

**Stretch goals:**

- Add your own questions to the array.
- Show "Correct!" or "Wrong!" feedback after each answer.
- Add a "Play again" button that resets `current` and `score`.

:::key
You turned a list of data into a playable game using **state** (`current`, `score`) and the **select → listen → update** loop. That's the blueprint for real JavaScript apps.
:::

## What's next

You're building real things now. Next, **Git and GitHub** shows you how to save and share projects like this one.

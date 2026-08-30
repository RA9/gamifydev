# Building the Results Screen

The user has answered every question. Now the quiz needs to detect that, hide the question screen, and show a results screen with the final score, a message, and a way to play again.

## Detecting the end of questions

In `nextQuestion`, you already have the branching logic:

```js
function nextQuestion() {
  currentIndex++;

  if (currentIndex < questions.length) {
    renderQuestion();
  } else {
    showResults();
  }
}
```

When `currentIndex` reaches `questions.length`, there are no more questions. That's the signal to switch screens.

:::key
This is the state → render loop in action. The state changes (`currentIndex++`), and the render logic checks whether to show another question or the results. No special "end" flag is needed — the data itself tells you when you're done.
:::

## Hiding the quiz screen, showing results

Switch screens by toggling the `.hidden` class:

```js
function showResults() {
  quizScreen.classList.add('hidden');
  resultsScreen.classList.remove('hidden');

  // Display the score
  const percentage = Math.round((score / questions.length) * 100);
  scoreDisplay.textContent = `You scored ${score} out of ${questions.length} (${percentage}%)`;

  // Show an encouraging message
  resultMessage.textContent = getMessage(percentage);
}
```

Grab the results screen DOM references at the top of your file:

```js
const resultsScreen = document.getElementById('results-screen');
const resultMessage = document.getElementById('result-message');
const scoreDisplay  = document.getElementById('score-display');
const restartBtn    = document.getElementById('restart-btn');
```

## Displaying the final score

The score display shows two things: the raw count and the percentage.

```js
const percentage = Math.round((score / questions.length) * 100);
scoreDisplay.textContent = `You scored ${score} out of ${questions.length} (${percentage}%)`;
```

`Math.round` ensures you get clean numbers like 80% instead of 79.9999%.

**Examples:**
- 4 out of 5 → `Math.round((4 / 5) * 100)` → **80%**
- 3 out of 8 → `Math.round((3 / 8) * 100)` → **38%**
- 5 out of 5 → `Math.round((5 / 5) * 100)` → **100%**

## Encouraging message based on score

Give the user feedback that matches their performance. This makes the results screen feel personal, not robotic.

```js
function getMessage(percentage) {
  if (percentage === 100) return "🏆 Perfect score! You nailed every question!";
  if (percentage >= 80)  return "🎉 Excellent work! You really know your stuff.";
  if (percentage >= 60)  return "👍 Good job! A little more practice and you'll master this.";
  if (percentage >= 40)  return "📚 Not bad — review the topics you missed and try again.";
  return "💪 Keep learning! Every expert was once a beginner.";
}
```

:::tip
Always make the lowest-score message encouraging, never insulting. A quiz should motivate the user to learn more, not make them feel bad. "Keep learning!" is better than "You failed."
:::

## Styling the results screen

```css
#results-screen {
  text-align: center;
  padding: 2rem;
}

#result-message {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

#score-display {
  font-size: 1.125rem;
  color: #64748b;
  margin-bottom: 2rem;
}

#restart-btn {
  padding: 0.75rem 2rem;
  background: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s ease;
}

#restart-btn:hover {
  background: #2563eb;
}
```

## The Restart button — resetting state

The Restart button resets all state to its initial values and re-renders the first question.

```js
function restartQuiz() {
  currentIndex = 0;
  score = 0;
  answered = false;

  resultsScreen.classList.add('hidden');
  quizScreen.classList.remove('hidden');

  renderQuestion();
}

restartBtn.addEventListener('click', restartQuiz);
```

Walk through what this does:

1. **Reset state** — `currentIndex` back to 0, `score` back to 0
2. **Switch screens** — hide results, show quiz
3. **Render the first question** — calls `renderQuestion()` which reads `currentIndex` (now 0)

:::key
Restart doesn't create new HTML or reload the page. It resets the state variables and calls the same `renderQuestion()` function that started the quiz. This is the power of the state → render loop — the same render function handles both initial load and restart.
:::

## The complete results flow

Here's how the functions connect:

```text
User clicks last Next
  → nextQuestion()
    → currentIndex++ (now equals questions.length)
    → showResults()
      → hide quiz screen, show results screen
      → calculate percentage
      → display score and message

User clicks Restart
  → restartQuiz()
    → reset state to initial values
    → hide results screen, show quiz screen
    → renderQuestion() (shows question 1 again)
```

## Adding a score breakdown (optional)

For a more detailed results screen, show which questions were right and wrong. This requires tracking answers as you go.

Add an answers array to your state:

```js
let answers = [];  // track what the user chose
```

In `handleAnswer`, record each answer:

```js
answers.push({
  question: questions[currentIndex].question,
  selected: selectedIndex,
  correct: questions[currentIndex].correct,
  isCorrect: selectedIndex === questions[currentIndex].correct
});
```

In `showResults`, render the breakdown:

```js
function renderBreakdown() {
  const breakdownEl = document.getElementById('breakdown');
  breakdownEl.innerHTML = '';

  answers.forEach((a, i) => {
    const div = document.createElement('div');
    div.classList.add('breakdown-item', a.isCorrect ? 'correct' : 'wrong');
    div.innerHTML = `
      <strong>Q${i + 1}:</strong> ${a.question}<br />
      <span>Your answer: ${questions[i].options[a.selected]}</span>
      ${!a.isCorrect ? `<br /><span>Correct: ${questions[i].options[a.correct]}</span>` : ''}
    `;
    breakdownEl.appendChild(div);
  });
}
```

```css
.breakdown-item {
  text-align: left;
  padding: 0.75rem;
  margin-bottom: 0.5rem;
  border-radius: 6px;
  font-size: 0.9rem;
}

.breakdown-item.correct {
  background: #f0fdf4;
  border-left: 3px solid #22c55e;
}

.breakdown-item.wrong {
  background: #fef2f2;
  border-left: 3px solid #ef4444;
}
```

Don't forget to reset `answers = []` in `restartQuiz`.

:::warning
If you add the breakdown feature, the `answers` array must be cleared on restart. Forgetting this means the breakdown accumulates results from multiple quiz runs, showing 10 answers after two 5-question runs.
:::

## The full showResults and restartQuiz

```js
function showResults() {
  quizScreen.classList.add('hidden');
  resultsScreen.classList.remove('hidden');

  const percentage = Math.round((score / questions.length) * 100);
  resultMessage.textContent = getMessage(percentage);
  scoreDisplay.textContent = `You scored ${score} out of ${questions.length} (${percentage}%)`;
}

function restartQuiz() {
  currentIndex = 0;
  score = 0;
  answered = false;

  resultsScreen.classList.add('hidden');
  quizScreen.classList.remove('hidden');

  renderQuestion();
}

restartBtn.addEventListener('click', restartQuiz);
```

:::quiz
Q: What triggers the transition from the question screen to the results screen?
- A timer running out
- The user clicking a special "Finish" button
- currentIndex reaching questions.length after the last nextQuestion call *
- The browser detecting the last question automatically
E: When `currentIndex++` makes it equal to `questions.length`, the `if (currentIndex < questions.length)` check in `nextQuestion` fails, and `showResults()` runs. The data length determines when the quiz ends.
:::

:::quiz
Q: Why does `restartQuiz` call `renderQuestion()` instead of reloading the page?
- Reloading the page is too slow
- renderQuestion reads the reset state and renders correctly, keeping the app in a single page lifecycle *
- The browser doesn't allow page reloads
- renderQuestion is faster than HTML
E: The state → render pattern means the same render function works for any state. After resetting `currentIndex` to 0, calling `renderQuestion()` shows question 1 — no page reload needed. This keeps the app fast and avoids losing any runtime state.
:::

## Recap

- **Detect the end** by checking `currentIndex >= questions.length` in `nextQuestion`.
- **Switch screens** by toggling the `.hidden` class on the quiz and results containers.
- **Calculate percentage** with `Math.round((score / questions.length) * 100)`.
- **Encourage the user** with a message tailored to their score range.
- **Restart** resets all state variables and calls `renderQuestion()` — the same function that started the quiz.
- **Optional breakdown** tracks each answer and shows a detailed right/wrong summary.

**Next up:** Adding Polish and Animations — progress bar, transitions, keyboard support, and responsive layout.

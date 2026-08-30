# Handling Answers and Scoring

The question engine renders questions and options. Now you'll make the quiz interactive — detect which option the user clicks, check if it's correct, show visual feedback, and track the score.

## The click handler on options

Each option button already has a click handler attached in `renderQuestion`:

```js
button.addEventListener('click', () => handleAnswer(index));
```

The `index` parameter tells `handleAnswer` which option was selected (0, 1, 2, or 3). Now implement the function:

```js
function handleAnswer(selectedIndex) {
  if (answered) return;  // prevent double-answering

  answered = true;
  const q = questions[currentIndex];

  // Get all option buttons
  const buttons = optionsEl.querySelectorAll('.option-btn');

  // Check if the selected answer is correct
  if (selectedIndex === q.correct) {
    score++;
    buttons[selectedIndex].classList.add('correct');
  } else {
    buttons[selectedIndex].classList.add('wrong');
    buttons[q.correct].classList.add('correct');  // always show the right answer
  }

  // Disable all buttons
  buttons.forEach(btn => btn.disabled = true);

  // Enable the Next button
  nextBtn.disabled = false;
}
```

:::key
The `if (answered) return` guard at the top is critical. Without it, a fast-clicking user could select multiple options and inflate their score. Always guard against repeated actions in event handlers.
:::

## Comparing selected index to correct index

The comparison is simple — strict equality:

```js
if (selectedIndex === q.correct) {
  // correct!
} else {
  // wrong — show what they picked AND the right answer
}
```

When wrong, you show two pieces of feedback:
1. The selected option turns **red** (wrong)
2. The correct option turns **green** (correct)

This teaches the user immediately — they don't have to guess which one was right.

## Adding .correct / .wrong CSS classes

These classes provide instant visual feedback through color.

```css
.option-btn.correct {
  background: #dcfce7;
  border-color: #22c55e;
  color: #166534;
}

.option-btn.wrong {
  background: #fef2f2;
  border-color: #ef4444;
  color: #991b1b;
}
```

The colors are intentionally soft — green-tinted background for correct, red-tinted for wrong. Strong enough to be clear, gentle enough not to feel harsh.

:::tip
Use background + border + text color together for feedback states, not just a single color change. This ensures the feedback is visible to users with color vision deficiency — the change is structural, not just a hue shift.
:::

## Disabling all options after selection

Once the user picks an answer, all four buttons should become unclickable. This prevents changing the answer and reinforces that the choice is final.

```js
buttons.forEach(btn => {
  btn.disabled = true;
});
```

The CSS for disabled buttons:

```css
.option-btn:disabled {
  cursor: not-allowed;
  opacity: 0.85;
}

.option-btn:disabled:hover {
  border-color: #e2e8f0;
  background: #f1f5f9;
}

/* But keep feedback colors visible even when disabled */
.option-btn.correct:disabled {
  background: #dcfce7;
  border-color: #22c55e;
  opacity: 1;
}

.option-btn.wrong:disabled {
  background: #fef2f2;
  border-color: #ef4444;
  opacity: 1;
}
```

:::warning
Order matters in CSS. The `.correct:disabled` and `.wrong:disabled` rules must come *after* the generic `:disabled` rule, or the generic opacity will override your feedback colors. More specific selectors win, but when specificity is equal, the last rule wins.
:::

## Incrementing the score

The score increments only when the selected answer matches the correct index:

```js
if (selectedIndex === q.correct) {
  score++;
  buttons[selectedIndex].classList.add('correct');
}
```

The score is used later on the results screen to show "You scored 4 out of 5." You don't need to display the running score during the quiz — it can pressure or distract. The progress indicator ("Question 3 of 5") is enough context.

## Enabling the Next button

The Next button starts disabled for each question (set in `renderQuestion`). After the user answers, enable it:

```js
nextBtn.disabled = false;
```

The user now has a clear flow:
1. Read the question
2. Click an option → see feedback
3. Click Next → advance to the next question

## The complete handleAnswer function

Here's the full function with comments:

```js
function handleAnswer(selectedIndex) {
  // Guard: don't allow answering twice
  if (answered) return;
  answered = true;

  const q = questions[currentIndex];
  const buttons = optionsEl.querySelectorAll('.option-btn');

  // Check correctness
  if (selectedIndex === q.correct) {
    score++;
    buttons[selectedIndex].classList.add('correct');
  } else {
    buttons[selectedIndex].classList.add('wrong');
    buttons[q.correct].classList.add('correct');
  }

  // Lock all options
  buttons.forEach(btn => {
    btn.disabled = true;
  });

  // Allow advancing
  nextBtn.disabled = false;
}
```

## Walking through a user interaction

Let's trace what happens when the user answers question 1 incorrectly:

1. **State before:** `currentIndex = 0`, `score = 0`, `answered = false`
2. **User clicks option 2** (index 2)
3. **`handleAnswer(2)` runs:**
   - `answered` is `false`, so we proceed
   - Set `answered = true`
   - `q.correct` is `0`, `selectedIndex` is `2` → not equal → wrong answer
   - `score` stays at `0`
   - Option 2 gets `.wrong` class (red)
   - Option 0 gets `.correct` class (green)
   - All buttons disabled
   - Next button enabled
4. **State after:** `currentIndex = 0`, `score = 0`, `answered = true`

The user sees: their wrong answer highlighted in red, the correct answer highlighted in green, and the Next button is now clickable.

:::key
Tracing through a user interaction step by step — tracking state changes at each point — is how you verify your logic works before running the code. This is a powerful debugging technique called a "desk check."
:::

## Adding a score indicator (optional)

If you want to show the running score during the quiz, add a small element:

```html
<p id="score-indicator"></p>
```

Update it in `handleAnswer`:

```js
const scoreIndicator = document.getElementById('score-indicator');

// Inside handleAnswer, after checking correctness:
scoreIndicator.textContent = `Score: ${score}`;
```

```css
#score-indicator {
  font-size: 0.875rem;
  color: #64748b;
  margin-top: 0.5rem;
}
```

This is optional — some quiz designs show the score, others don't. Both are valid choices.

:::quiz
Q: Why does `handleAnswer` show the correct answer even when the user picks wrong?
- To punish the user for guessing
- To teach the user the right answer immediately, turning each question into a learning moment *
- Because the browser requires it
- To make the quiz easier
E: Showing the correct answer after a wrong selection turns the quiz into a learning tool, not just a test. The user walks away knowing the right answer, which is the whole point of an educational quiz.
:::

:::quiz
Q: What does the `if (answered) return` guard prevent?
- It prevents the quiz from loading
- It prevents the user from clicking multiple options and inflating the score *
- It prevents the Next button from working
- It stops the browser from crashing
E: Without the guard, a user could click multiple options before the buttons are disabled (or click very fast). Each click would potentially increment the score. The guard ensures only the first click counts.
:::

## Recap

- **`handleAnswer(selectedIndex)`** is called when the user clicks an option button.
- **Guard against double-clicks** with `if (answered) return` at the top.
- **Compare** `selectedIndex === q.correct` to determine correctness.
- **Add `.correct`** (green) to the right answer and **`.wrong`** (red) to wrong selections.
- **Disable all buttons** after selection to lock in the answer.
- **Increment `score`** only on correct answers.
- **Enable the Next button** so the user can advance.
- **Desk-check** your logic by tracing through a user interaction step by step.

**Next up:** Building the Results Screen — detecting the end of the quiz and showing the final score.

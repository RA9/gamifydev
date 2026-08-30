# Adding Polish and Animations

The quiz works. Now make it feel good. A progress bar, smooth transitions, keyboard support, and responsive styling turn a functional quiz into a polished one.

## Progress bar that fills as questions advance

A visual progress bar gives users a sense of momentum. Add it above the question text.

```html
<div class="progress-bar">
  <div class="progress-fill" id="progress-fill"></div>
</div>
```

```css
.progress-bar {
  width: 100%;
  height: 6px;
  background: #e2e8f0;
  border-radius: 3px;
  margin-bottom: 1.5rem;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #3b82f6;
  border-radius: 3px;
  width: 0%;
  transition: width 0.4s ease;
}
```

Update the fill width in `renderQuestion`:

```js
const progressFill = document.getElementById('progress-fill');

function renderQuestion() {
  // Update progress bar
  const percent = ((currentIndex) / questions.length) * 100;
  progressFill.style.width = `${percent}%`;

  // ... rest of renderQuestion
}
```

When the quiz ends, fill the bar to 100% in `showResults`:

```js
function showResults() {
  progressFill.style.width = '100%';
  // ... rest of showResults
}
```

:::tip
The progress bar uses `currentIndex` (not `currentIndex + 1`) for the width calculation so it fills gradually. It starts at 0%, reaches 80% at question 5 of 5, and hits 100% when results display.
:::

## Transition on option selection

When an option is selected, animate the color change instead of making it instant:

```css
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
  transition: background 0.3s ease,
              border-color 0.3s ease,
              transform 0.15s ease;
}

.option-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  background: #eff6ff;
  transform: translateY(-1px);
}

.option-btn.correct {
  background: #dcfce7;
  border-color: #22c55e;
  color: #166534;
  transform: scale(1.02);
}

.option-btn.wrong {
  background: #fef2f2;
  border-color: #ef4444;
  color: #991b1b;
}
```

The `transition` property on `.option-btn` makes the class changes animate smoothly. The correct answer gets a subtle `scale(1.02)` to make it "pop" — a small visual reward.

## Card entrance animation

When a new question renders, animate the quiz card sliding in. Use a CSS animation triggered by a class.

```css
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(30px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.quiz-container.animate {
  animation: slideIn 0.3s ease forwards;
}
```

In `renderQuestion`, toggle the animation class:

```js
function renderQuestion() {
  const container = document.querySelector('.quiz-container');

  // Remove and re-add class to retrigger animation
  container.classList.remove('animate');
  // Force reflow so the browser sees the class removal
  void container.offsetWidth;
  container.classList.add('animate');

  // ... rest of renderQuestion
}
```

:::key
The `void container.offsetWidth` line forces a browser reflow. Without it, removing and immediately re-adding the same class happens in the same paint cycle, so the browser skips the animation. This is a standard trick for retriggering CSS animations.
:::

## Keyboard support — number keys and Enter

Power users (and accessibility-conscious developers) expect keyboard navigation. Add support for pressing 1–4 to select an option and Enter to advance.

```js
document.addEventListener('keydown', (e) => {
  // Number keys 1-4 select options
  const num = parseInt(e.key);
  if (num >= 1 && num <= 4 && !answered) {
    const buttons = optionsEl.querySelectorAll('.option-btn');
    if (buttons[num - 1]) {
      handleAnswer(num - 1);
    }
  }

  // Enter key clicks Next
  if (e.key === 'Enter' && !nextBtn.disabled) {
    nextQuestion();
  }
});
```

Add a visual hint so users know keyboard shortcuts exist:

```html
<p class="keyboard-hint">Tip: Press 1-4 to answer, Enter to continue</p>
```

```css
.keyboard-hint {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 1rem;
  text-align: center;
}
```

:::warning
Always test keyboard handlers by pressing keys at different quiz states. Can you press "2" after already answering? Can you press Enter before selecting an option? The `!answered` and `!nextBtn.disabled` guards handle these edge cases.
:::

## Visual option numbering

To make the keyboard shortcuts intuitive, show numbers next to each option:

```js
q.options.forEach((option, index) => {
  const button = document.createElement('button');
  button.classList.add('option-btn');
  button.textContent = `${index + 1}. ${option}`;
  button.addEventListener('click', () => handleAnswer(index));
  optionsEl.appendChild(button);
});
```

Now the user sees:
```
1. Hyper Text Markup Language
2. High Tech Modern Language
3. Hyper Transfer Markup Language
4. Home Tool Markup Language
```

The "1" on screen maps to the "1" key on the keyboard. Intuitive.

## Responsive layout

The quiz should work on any screen. Mobile-first CSS:

```css
.quiz-container {
  max-width: 600px;
  margin: 1rem auto;
  padding: 1.5rem;
  font-family: system-ui, sans-serif;
}

@media (min-width: 600px) {
  .quiz-container {
    margin: 3rem auto;
    padding: 2.5rem;
  }

  #question {
    font-size: 1.5rem;
  }
}

/* Ensure touch targets are large enough on mobile */
.option-btn {
  min-height: 48px;
}

#next-btn,
#restart-btn {
  min-height: 48px;
  width: 100%;
}

@media (min-width: 600px) {
  #next-btn,
  #restart-btn {
    width: auto;
  }
}
```

On mobile:
- The container uses full width with small padding
- Buttons are full-width for easy tapping
- All touch targets meet the 48px minimum

On desktop:
- The container is centered with a 600px max-width
- Buttons shrink to fit their content

## Accessibility check

Run through this checklist before calling the quiz polished:

```text
✅ All buttons have readable text (not just icons)
✅ Keyboard navigation works (Tab, 1-4, Enter)
✅ Color isn't the only indicator of correct/wrong
   (also uses border + background change)
✅ Focus states are visible on all interactive elements
✅ Progress is announced (text + visual bar)
✅ Results message is clear and encouraging
```

Add visible focus styles:

```css
.option-btn:focus-visible {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

#next-btn:focus-visible,
#restart-btn:focus-visible {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}
```

:::tip
Test your quiz using only the keyboard — no mouse at all. Tab to each button, press Enter to activate, use 1–4 for options. If anything feels broken or invisible, fix it before shipping.
:::

## Adding a timer display (optional teaser)

A countdown timer per question adds urgency. Here's a minimal implementation:

```js
let timeLeft = 15;
let timer = null;

function startTimer() {
  timeLeft = 15;
  timerEl.textContent = `⏱ ${timeLeft}s`;

  timer = setInterval(() => {
    timeLeft--;
    timerEl.textContent = `⏱ ${timeLeft}s`;

    if (timeLeft <= 0) {
      clearInterval(timer);
      handleAnswer(-1);  // -1 = no selection (auto-wrong)
    }
  }, 1000);
}
```

Call `startTimer()` at the end of `renderQuestion()` and `clearInterval(timer)` in `handleAnswer()`. This is an extension idea — don't add it unless the core quiz is solid first.

:::quiz
Q: Why is `void container.offsetWidth` needed before re-adding an animation class?
- It makes the animation faster
- It forces a browser reflow so the class removal and re-addition happen in separate paint cycles, retriggering the animation *
- It clears the browser cache
- It's required by the CSS specification
E: Without forcing a reflow, the browser optimizes by batching the class removal and re-addition into one operation, effectively doing nothing. Reading `offsetWidth` forces the browser to recalculate layout, creating the break needed to restart the animation.
:::

:::quiz
Q: Why should the correct/wrong feedback use more than just color changes?
- Because colors are slower to render
- To ensure users with color vision deficiency can still distinguish correct from wrong answers *
- Because CSS doesn't support color on buttons
- To reduce the number of CSS rules needed
E: Relying solely on color excludes users who can't distinguish red from green. Using border changes, background changes, and scale transforms ensures the feedback is perceivable through multiple visual channels.
:::

## Recap

- **Progress bar** with `transition: width 0.4s ease` fills smoothly as questions advance.
- **Option transitions** animate background and border color changes for polished feedback.
- **Card entrance animation** uses `@keyframes slideIn` retriggered with a reflow trick.
- **Keyboard support** (1–4 for options, Enter for Next) makes the quiz fast and accessible.
- **Responsive layout** uses full-width buttons on mobile and a centered container on desktop.
- **Accessibility** requires focus styles, multiple feedback channels beyond color, and full keyboard operability.

**Next up:** Code Review and Extensions — reviewing your code structure, refactoring, and ideas for taking the quiz further.

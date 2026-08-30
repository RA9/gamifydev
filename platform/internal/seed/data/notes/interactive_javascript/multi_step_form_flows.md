# Multi-Step Form Flows

Long forms are overwhelming. Multi-step flows break them into digestible pieces: one step at a time, a progress bar showing how far along the user is, and the ability to go back without losing data. This lesson builds one from scratch using the state → render pattern.

## The state — a step counter and form data

```js
const state = {
  step: 1,
  totalSteps: 3,
  formData: {
    name: "",
    email: "",
    plan: "free",
    cardNumber: "",
  },
  errors: {},
};
```

`step` drives everything: which panel is visible, which progress indicator is highlighted, and which buttons appear.

## HTML structure

```html
<div class="progress-bar">
  <div class="progress-fill" id="progress-fill"></div>
</div>
<p id="step-label">Step 1 of 3</p>

<form id="wizard">
  <div class="step-panel" data-step="1">
    <h2>Your Info</h2>
    <label>Name <input name="name" required /></label>
    <label>Email <input name="email" type="email" required /></label>
  </div>

  <div class="step-panel" data-step="2" hidden>
    <h2>Choose a Plan</h2>
    <label><input type="radio" name="plan" value="free" checked /> Free</label>
    <label><input type="radio" name="plan" value="pro" /> Pro</label>
  </div>

  <div class="step-panel" data-step="3" hidden>
    <h2>Payment</h2>
    <label>Card Number <input name="cardNumber" placeholder="1234 5678 9012 3456" /></label>
  </div>

  <div class="form-actions">
    <button type="button" id="back-btn" hidden>Back</button>
    <button type="button" id="next-btn">Next</button>
    <button type="submit" id="submit-btn" hidden>Submit</button>
  </div>
</form>
```

## CSS for the progress bar

```css
.progress-bar {
  width: 100%;
  height: 6px;
  background: #e0e0e0;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: teal;
  transition: width 0.3s ease;
}

.step-panel[hidden] { display: none; }
```

## The render function

```js
const panels = document.querySelectorAll(".step-panel");
const progressFill = document.querySelector("#progress-fill");
const stepLabel = document.querySelector("#step-label");
const backBtn = document.querySelector("#back-btn");
const nextBtn = document.querySelector("#next-btn");
const submitBtn = document.querySelector("#submit-btn");

function render() {
  // Show the active panel, hide the rest
  panels.forEach(panel => {
    panel.hidden = Number(panel.dataset.step) !== state.step;
  });

  // Update progress bar
  const percent = (state.step / state.totalSteps) * 100;
  progressFill.style.width = `${percent}%`;
  stepLabel.textContent = `Step ${state.step} of ${state.totalSteps}`;

  // Show/hide navigation buttons
  backBtn.hidden = state.step === 1;
  nextBtn.hidden = state.step === state.totalSteps;
  submitBtn.hidden = state.step !== state.totalSteps;
}
```

:::tip
The progress bar width is a simple calculation: `(currentStep / totalSteps) * 100`. Because the width is set as an inline style with a CSS transition, the bar animates smoothly between steps.
:::

## Per-step validation

Each step should validate its own inputs before allowing the user to proceed:

```js
function validateStep(step) {
  state.errors = {};

  if (step === 1) {
    if (!state.formData.name.trim()) {
      state.errors.name = "Name is required";
    }
    if (!state.formData.email.includes("@")) {
      state.errors.email = "Enter a valid email";
    }
  }

  if (step === 2) {
    if (!state.formData.plan) {
      state.errors.plan = "Select a plan";
    }
  }

  if (step === 3) {
    if (state.formData.plan === "pro" && !state.formData.cardNumber.trim()) {
      state.errors.cardNumber = "Card number is required for Pro plan";
    }
  }

  return Object.keys(state.errors).length === 0;
}
```

## Collecting input values

Sync the form inputs into state whenever they change:

```js
const form = document.querySelector("#wizard");

form.addEventListener("input", (event) => {
  const { name, value, type, checked } = event.target;
  if (!name) return;

  if (type === "radio" || type === "checkbox") {
    state.formData[name] = type === "checkbox" ? checked : value;
  } else {
    state.formData[name] = value;
  }
});
```

## Navigation — next and back

```js
nextBtn.addEventListener("click", () => {
  if (!validateStep(state.step)) {
    renderErrors();
    return;
  }
  state.step += 1;
  saveDraft();
  render();
});

backBtn.addEventListener("click", () => {
  state.step -= 1;
  render();
});
```

:::key
The Back button should *never* validate — let the user go back freely to fix earlier fields. Only the Next and Submit buttons need validation.
:::

## Submitting on the final step

```js
form.addEventListener("submit", (event) => {
  event.preventDefault();

  if (!validateStep(state.step)) {
    renderErrors();
    return;
  }

  console.log("Submitting:", state.formData);
  localStorage.removeItem("form-draft"); // clear the saved draft
  // send to server, show success message, etc.
});
```

## Saving drafts to localStorage

Users should not lose their progress if they accidentally close the tab:

```js
function saveDraft() {
  localStorage.setItem("form-draft", JSON.stringify({
    step: state.step,
    formData: state.formData,
  }));
}

function loadDraft() {
  try {
    const saved = localStorage.getItem("form-draft");
    if (!saved) return;

    const draft = JSON.parse(saved);
    state.step = draft.step || 1;
    state.formData = { ...state.formData, ...draft.formData };

    // Restore input values in the DOM
    for (const [name, value] of Object.entries(state.formData)) {
      const input = form.querySelector(`[name="${name}"]`);
      if (!input) continue;
      if (input.type === "radio") {
        const radio = form.querySelector(`[name="${name}"][value="${value}"]`);
        if (radio) radio.checked = true;
      } else {
        input.value = value;
      }
    }
  } catch {
    // corrupt draft — start fresh
  }
}
```

Call `loadDraft()` before the first `render()`:

```js
loadDraft();
render();
```

:::warning
Always wrap `JSON.parse` in a `try/catch` when loading from localStorage. The stored data might be corrupt, from an old version of your form, or manually tampered with. Failing gracefully beats crashing on load.
:::

## Practice

:::quiz
Q: A user fills out step 1 and step 2, then accidentally closes the browser tab. When they reopen the page, what should happen?
- The form starts over from step 1 with empty fields
- The form resumes at step 2 with their data restored from localStorage *
- The form shows step 3 because they completed step 2
- An error message appears
E: The `saveDraft` function saves both the current step and form data to localStorage on every change. The `loadDraft` function restores this state on page load, letting the user continue where they left off.
:::

:::quiz
Q: The user is on step 2 and clicks the Back button. Should the form validate step 2 before going back?
- Yes, always validate before any navigation
- No, Back should navigate freely without validation *
- Only if there are errors
- Only for the first step
E: The Back button should let users move freely to previous steps to review or fix earlier data. Validation only applies when moving *forward* (Next or Submit) to ensure each step is complete before proceeding.
:::

## Recap

- State for a multi-step form: `step` (integer), `formData` (object), `errors` (object).
- **Render** shows/hides panels based on `step` and updates the progress bar width.
- **Progress bar:** width = `(step / totalSteps) * 100` percent, animated with CSS transition.
- **Validate per step** — only block forward navigation, never block going back.
- Sync inputs into `state.formData` with an `input` event listener.
- **Save drafts** to localStorage on every change; load on startup with `try/catch`.
- On final submit, clear the draft from localStorage.

**Next up:** Search and Filter UIs — combining input events, state, and rendering for live data exploration.

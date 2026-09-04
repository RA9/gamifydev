// problems.js — the language picker on a practice problem.
//
// Switching language swaps in that language's starter, but only when the
// editor still holds the starter it was given. Someone who has written code and
// then changes the picker to check what the C skeleton looks like should not
// lose their work to a dropdown.

const editor = document.getElementById("pb-editor");
const code = document.getElementById("pb-code");
const picker = document.getElementById("pb-lang");
const startersEl = document.getElementById("pb-starters");
if (editor && code && picker && startersEl) {
  let starters = {};
  try {
    starters = JSON.parse(startersEl.textContent) || {};
  } catch {
    starters = {};
  }

  // What the box held when the page loaded: either a starter, or a previous
  // submission. Either way, it is the text the learner has not touched.
  let untouched = code.value;
  code.addEventListener("input", () => {
    // Once they type, nothing is safe to replace until they ask for it.
    untouched = null;
  });

  picker.addEventListener("change", () => {
    const lang = picker.value;
    const starter = starters[lang] ?? "";
    editor.setAttribute("lang", lang);
    if (untouched === null) {
      const keep = confirm(
        "Replace your code with the " + lang + " starter?\n\nYour current code will be lost.",
      );
      if (!keep) return;
    }
    code.value = starter;
    untouched = starter;
    // The editor draws its own gutter and line count from this event.
    code.dispatchEvent(new Event("input", { bubbles: true }));
  });
}

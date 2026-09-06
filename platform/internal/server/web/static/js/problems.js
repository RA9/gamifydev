// problems.js — the problem workspace.
//
// Everything here is a convenience over markup that already works: the panes
// have sensible widths, the console shows the cases, and both buttons submit
// through htmx. Nothing below is load-bearing.

const ide = document.getElementById("ide");
if (ide) {
  const editor = document.getElementById("pb-editor");
  const code = document.getElementById("pb-code");
  const picker = document.getElementById("pb-lang");
  const startersEl = document.getElementById("pb-starters");

  let starters = {};
  try {
    starters = JSON.parse(startersEl?.textContent || "{}") || {};
  } catch {
    starters = {};
  }
  const starterFor = (lang) => starters[lang] ?? "";

  // --- language picker ------------------------------------------------------
  // Swapping language swaps the starter, but only while the box still holds
  // the text it was given. Somebody who has written code and then looks at
  // another language's skeleton should not lose their work to a dropdown.
  let untouched = code ? code.value : null;
  code?.addEventListener("input", () => {
    untouched = null;
  });

  function setCode(text) {
    code.value = text;
    untouched = text;
    // The editor draws its gutter and line count from this.
    code.dispatchEvent(new Event("input", { bubbles: true }));
  }

  picker?.addEventListener("change", () => {
    const lang = picker.value;
    editor?.setAttribute("lang", lang);
    if (untouched === null && !confirm(
      "Replace your code with the " + lang + " starter?\n\nYour current code will be lost.",
    )) return;
    setCode(starterFor(lang));
  });

  document.getElementById("ideReset")?.addEventListener("click", () => {
    const lang = picker ? picker.value : Object.keys(starters)[0];
    if (untouched !== null || confirm("Discard your code and start from the starter?")) {
      setCode(starterFor(lang));
    }
  });

  // --- copy buttons ---------------------------------------------------------
  ide.addEventListener("click", (e) => {
    const btn = e.target.closest("[data-copy-btn]");
    if (!btn) return;
    const source = btn.parentElement?.querySelector("[data-copy]");
    if (!source) return;
    navigator.clipboard?.writeText(source.textContent.trim()).then(() => {
      btn.classList.add("is-done");
      setTimeout(() => btn.classList.remove("is-done"), 1200);
    });
  });

  // --- console tabs ---------------------------------------------------------
  // Test cases and output share one body, which htmx replaces wholesale. The
  // tabs scroll to what is already there rather than hiding half of it: a
  // verdict the learner cannot see because a tab is selected is a bug, and the
  // panel is short enough that scrolling is the honest answer.
  const body = document.getElementById("ide-console-body");
  ide.querySelectorAll(".ide-ctab").forEach((tab) => {
    tab.addEventListener("click", () => {
      ide.querySelectorAll(".ide-ctab").forEach((t) => t.classList.toggle("is-on", t === tab));
      const target = tab.dataset.panel === "output"
        ? body?.querySelector(".ide-output, .verdict")
        : body?.querySelector(".case-list, .case-head");
      target?.scrollIntoView({ block: "nearest", behavior: "smooth" });
    });
  });

  // --- draggable split ------------------------------------------------------
  const split = document.getElementById("ideSplit");
  const MIN = 22, MAX = 74; // percent, so neither pane can be dragged shut
  function setLeft(pct) {
    ide.style.setProperty("--ide-left", Math.min(MAX, Math.max(MIN, pct)) + "fr");
  }
  let dragging = false;
  split?.addEventListener("pointerdown", (e) => {
    dragging = true;
    split.setPointerCapture(e.pointerId);
    document.body.style.userSelect = "none";
  });
  split?.addEventListener("pointermove", (e) => {
    if (!dragging) return;
    const box = ide.getBoundingClientRect();
    setLeft(((e.clientX - box.left) / box.width) * 100);
  });
  const stop = (e) => {
    if (!dragging) return;
    dragging = false;
    split.releasePointerCapture?.(e.pointerId);
    document.body.style.userSelect = "";
  };
  split?.addEventListener("pointerup", stop);
  split?.addEventListener("pointercancel", stop);
  // Keyboard, because a divider that only responds to a mouse is not a control.
  split?.addEventListener("keydown", (e) => {
    const step = e.key === "ArrowLeft" ? -3 : e.key === "ArrowRight" ? 3 : 0;
    if (!step) return;
    e.preventDefault();
    const cur = parseFloat(getComputedStyle(ide).getPropertyValue("--ide-left")) || 46;
    setLeft(cur + step);
  });
}

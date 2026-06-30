// step-runner.js — drives an interactive, step-by-step lesson:
//   - live-previews the learner's HTML in a sandboxed iframe
//   - on "Check", runs each step's author-defined checks against the rendered DOM
//   - when all checks pass, records progress and reveals the Next button
//
// Each check is an <li data-test="<js boolean expression>"> where the expression
// is evaluated with `doc` (the preview document) and `code` (the raw editor text)
// in scope. Authors are trusted (same as lesson content), so eval is acceptable.
(function () {
  const lab = document.querySelector(".step-lab");
  if (!lab) return;
  const ta = document.getElementById("stepCode");
  const frame = document.getElementById("stepPreview");
  const checkBtn = document.getElementById("checkBtn");
  const nextBtn = document.getElementById("nextBtn");
  const msg = document.getElementById("stepMsg");
  const items = [...document.querySelectorAll("#stepChecks li")];
  if (!ta || !frame) return;

  const completeURL = lab.dataset.completeUrl;
  let alreadyDone = lab.dataset.completed === "1";

  // --- live preview (debounced) ---
  const renderPreview = () => {
    frame.srcdoc = ta.value;
  };
  let timer = null;
  ta.addEventListener("input", () => {
    clearTimeout(timer);
    timer = setTimeout(renderPreview, 250);
  });
  renderPreview();

  // Render the code into the iframe and resolve once it has loaded.
  const renderAndWait = () =>
    new Promise((resolve) => {
      const done = () => {
        frame.removeEventListener("load", done);
        resolve(frame.contentDocument);
      };
      frame.addEventListener("load", done);
      frame.srcdoc = ta.value;
    });

  const runCheck = (doc, code, test) => {
    try {
      // eslint-disable-next-line no-new-func
      return !!Function("doc", "code", "return (" + test + ");")(doc, code);
    } catch (_) {
      return false;
    }
  };

  const markComplete = () => {
    if (alreadyDone) return;
    alreadyDone = true;
    fetch(completeURL, { method: "POST", credentials: "same-origin" }).catch(() => {});
  };

  checkBtn.addEventListener("click", async () => {
    checkBtn.disabled = true;
    const doc = await renderAndWait();
    const code = ta.value;
    let allPass = true;
    items.forEach((li) => {
      const ok = runCheck(doc, code, li.dataset.test);
      li.classList.toggle("pass", ok);
      li.classList.toggle("fail", !ok);
      const mark = li.querySelector(".check-mark");
      if (mark) mark.textContent = ok ? "✓" : "✗";
      if (!ok) allPass = false;
    });
    checkBtn.disabled = false;

    if (allPass) {
      msg.textContent = "Great — all checks passed! 🎉";
      msg.className = "step-msg is-ok";
      markComplete();
      if (nextBtn) {
        nextBtn.hidden = false;
        nextBtn.focus();
      }
    } else {
      msg.textContent = "Not quite — fix the items marked ✗ and check again.";
      msg.className = "step-msg is-bad";
    }
  });

  // Tab inserts two spaces instead of leaving the editor.
  ta.addEventListener("keydown", (e) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const s = ta.selectionStart, en = ta.selectionEnd;
      ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(en);
      ta.selectionStart = ta.selectionEnd = s + 2;
    }
  });
})();

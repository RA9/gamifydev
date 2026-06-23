// Code Lab — hands-on coding sessions graded by running the learner's code
// against hidden test cases. This is the objective foundation of the
// certification program: you don't just read, you write code that has to pass.

let GD_CHALLENGES = null;
const GD_EDITOR_INSTANCE = { current: null };

async function loadChallenges() {
  if (GD_CHALLENGES) return GD_CHALLENGES;
  try {
    const res = await fetch("./data/coding-challenges.json");
    GD_CHALLENGES = (await res.json()).challenges || [];
  } catch (e) {
    GD_CHALLENGES = [];
  }
  return GD_CHALLENGES;
}

function getChallenge(id) {
  return (GD_CHALLENGES || []).find((c) => c.id === id) || null;
}

async function isChallengeDone(id) {
  if (typeof getMeta !== "function") return false;
  return !!(await getMeta("code:" + id, false));
}

const GD_DIFF_RANK = { "Warm-up": 0, Easy: 1, Medium: 2, Hard: 3 };

// --- Gallery ----------------------------------------------------------------
async function CodeLabPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  const challenges = await loadChallenges();
  let doneCount = 0;

  const cards = await Promise.all(
    challenges.map(async (c) => {
      const done = await isChallengeDone(c.id);
      if (done) doneCount += 1;
      const lang = (c.language || "javascript").toUpperCase();
      return `
      <button type="button" onclick="openChallenge('${c.id}')"
        class="gd-card group text-left flex flex-col gap-2 transition-transform hover:-translate-y-1">
        <div class="flex items-start justify-between gap-2">
          <span class="gd-chip gd-chip-brand !py-0.5 !px-2 !text-[10px]">${lang}</span>
          ${
            done
              ? `<span class="gd-chip gd-chip-grass !text-[10px] !py-0.5">${icon("check", "w-3 h-3")} Solved</span>`
              : `<span class="gd-chip !text-[10px] !py-0.5">${c.difficulty || ""}</span>`
          }
        </div>
        <h3 class="text-base font-extrabold leading-tight">${c.title}</h3>
        <p class="text-xs text-slate-500 leading-snug">${(c.prompt || "").replace(/[`*_#]/g, "").slice(0, 96)}…</p>
        <span class="mt-1 text-xs font-bold text-brand-600 group-hover:text-brand-700">${done ? "Revisit" : "Solve"} →</span>
      </button>`;
    })
  );

  htmlEl.innerHTML = `
    <div class="max-w-5xl mx-auto animate-fade-up">
      <div class="text-center mb-6">
        <span class="gd-chip gd-chip-brand mb-2">${icon("code", "w-3.5 h-3.5")} Write real code</span>
        <h1 class="text-2xl sm:text-3xl font-extrabold">Code Lab</h1>
        <p class="text-slate-500 mt-1 max-w-md mx-auto">Hands-on coding sessions. Write a solution, run it against the tests, and pass them all. Your code is graded for real — no multiple choice.</p>
        <p class="mt-3 text-sm font-bold text-slate-500">${doneCount} / ${challenges.length} solved</p>
      </div>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">${cards.join("")}</div>
    </div>`;
}

function openChallenge(id) {
  window.location.hash = "code/" + id;
}

// --- Challenge player -------------------------------------------------------
async function ChallengePage(htmlEl, id) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  await loadChallenges();
  const c = getChallenge(id);
  if (!c) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <p class="text-slate-600 mb-5">That challenge doesn't exist.</p>
          <a href="#code" class="gd-btn gd-btn-primary">Back to Code Lab</a>
        </div>
      </div>`;
    return;
  }

  const done = await isChallengeDone(id);
  const draft = typeof getMeta === "function" ? await getMeta("code-draft:" + id, null) : null;
  const visibleTests = (c.tests || []).filter((t) => !t.hidden);
  const hiddenCount = (c.tests || []).length - visibleTests.length;

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-4">
      <a href="#code" class="inline-flex items-center gap-1 text-sm font-bold text-brand-600 hover:text-brand-700">${icon("arrowLeft", "w-4 h-4")} Code Lab</a>

      <div class="gd-card">
        <div class="flex items-center gap-2 mb-2">
          <span class="gd-chip gd-chip-brand !py-0.5 !px-2 !text-[10px]">${(c.language || "js").toUpperCase()}</span>
          <span class="gd-chip !py-0.5 !px-2 !text-[10px]">${c.difficulty || ""}</span>
          ${done ? `<span class="gd-chip gd-chip-grass !py-0.5 !px-2 !text-[10px]">${icon("check", "w-3 h-3")} Solved</span>` : ""}
        </div>
        <h1 class="text-xl font-extrabold mb-2">${c.title}</h1>
        <div class="gd-prose text-slate-700">${typeof markdownToHtml === "function" ? markdownToHtml(c.prompt || "") : c.prompt || ""}</div>
      </div>

      <div class="gd-card">
        <div class="flex items-center justify-between mb-2">
          <h2 class="text-sm font-extrabold uppercase tracking-wide text-slate-500">Your solution</h2>
          <button id="code-reset" class="text-xs font-bold text-slate-400 hover:text-slate-600">Reset</button>
        </div>
        <div id="code-editor" class="gd-editor-wrap"></div>
        <div class="mt-3 flex flex-wrap items-center gap-3">
          <button id="code-run" class="gd-btn gd-btn-primary !py-2 !px-5">${icon("play", "w-4 h-4")} Run tests</button>
          <span class="text-xs text-slate-400">${visibleTests.length} sample test${visibleTests.length === 1 ? "" : "s"}${hiddenCount ? ` + ${hiddenCount} hidden` : ""}</span>
        </div>
        <div id="code-results" class="mt-4"></div>
      </div>
    </div>`;

  // Mount editor with the saved draft or the starter code.
  const mount = htmlEl.querySelector("#code-editor");
  if (GD_EDITOR_INSTANCE.current) {
    try { GD_EDITOR_INSTANCE.current.destroy(); } catch (e) {}
  }
  const editor = createCodeEditor(mount, {
    doc: draft != null ? draft : c.starter || "",
    language: c.language,
    onChange: (val) => {
      if (typeof setMeta === "function") setMeta("code-draft:" + id, val);
    },
  });
  GD_EDITOR_INSTANCE.current = editor;

  htmlEl.querySelector("#code-reset").addEventListener("click", () => {
    editor.setValue(c.starter || "");
    if (typeof setMeta === "function") setMeta("code-draft:" + id, c.starter || "");
    htmlEl.querySelector("#code-results").innerHTML = "";
  });

  htmlEl.querySelector("#code-run").addEventListener("click", async () => {
    const runBtn = htmlEl.querySelector("#code-run");
    const resultsEl = htmlEl.querySelector("#code-results");
    runBtn.disabled = true;
    runBtn.classList.add("opacity-60");
    const loadingPy =
      (c.language || "").toLowerCase().startsWith("py") &&
      typeof isPyReady === "function" &&
      !isPyReady();
    resultsEl.innerHTML = `<p class="text-sm font-bold text-slate-400">${
      loadingPy ? "Loading Python… (one-time ~10s download)" : "Running…"
    }</p>`;

    const code = editor.getValue();
    const res = await runChallenge(c, code);

    runBtn.disabled = false;
    runBtn.classList.remove("opacity-60");
    renderResults(resultsEl, res, c);

    if (res.ok && res.results.length && res.results.every((r) => r.pass)) {
      const wasDone = await isChallengeDone(id);
      if (typeof setMeta === "function") await setMeta("code:" + id, true);
      if (!wasDone && typeof pixelCelebrate === "function") {
        pixelCelebrate({
          title: "All tests passed! ✅",
          message: `You solved ${c.title} — your code works against every case, including the hidden ones. That's the real deal.`,
          mood: "celebrate",
          cta: "Onward",
        });
      }
    }
  });

  editor.focus();
}

function renderResults(el, res, challenge) {
  if (res.pending) {
    el.innerHTML = `<div class="gd-callout gd-callout-tip"><div class="gd-callout-body"><p>${res.compileError}</p></div></div>`;
    return;
  }
  if (!res.ok && (res.compileError || res.timeout)) {
    el.innerHTML = `<div class="rounded-2xl border-2 border-rose-200 bg-rose-50 p-3 text-sm">
      <p class="font-extrabold text-rose-700 mb-1">${res.timeout ? "Timed out" : "Couldn't run your code"}</p>
      <pre class="whitespace-pre-wrap font-mono text-xs text-rose-800">${(res.compileError || "").replace(/</g, "&lt;")}</pre>
    </div>`;
    return;
  }

  const passed = res.results.filter((r) => r.pass).length;
  const total = res.results.length;
  const allPass = passed === total && total > 0;

  const rows = res.results
    .map((r) => {
      const detail =
        r.hidden
          ? `<span class="text-xs text-slate-400">hidden test</span>`
          : `<code class="text-xs text-slate-500">${(JSON.stringify(r.args || []) || "").replace(/^\[|\]$/g, "")} → expected ${r.expected}${
              r.pass ? "" : `, got ${r.error ? "error" : r.got}`
            }</code>`;
      return `
      <div class="flex items-center justify-between gap-3 rounded-xl border px-3 py-2 ${
        r.pass ? "border-grass-200 bg-grass-50" : "border-rose-200 bg-rose-50"
      }">
        <div class="flex items-center gap-2 min-w-0">
          <span class="${r.pass ? "text-grass-600" : "text-rose-500"}">${icon(r.pass ? "checkCircle" : "xCircle", "w-4 h-4")}</span>
          <span class="font-bold text-sm text-slate-700 truncate">${r.name || "test"}</span>
        </div>
        ${detail}
      </div>`;
    })
    .join("");

  const logs =
    res.logs && res.logs.length
      ? `<details class="mt-3"><summary class="text-xs font-bold text-slate-400 cursor-pointer">console output (${res.logs.length})</summary><pre class="mt-1 whitespace-pre-wrap font-mono text-xs text-slate-500 bg-slate-50 rounded-xl p-2">${res.logs.join("\n").replace(/</g, "&lt;")}</pre></details>`
      : "";

  el.innerHTML = `
    <div class="mb-2 flex items-center gap-2">
      <span class="text-sm font-extrabold ${allPass ? "text-grass-600" : "text-slate-600"}">${passed} / ${total} tests passed</span>
      ${allPass ? `<span class="gd-chip gd-chip-grass !text-[10px]">${icon("check", "w-3 h-3")} Solved</span>` : ""}
    </div>
    <div class="space-y-1.5">${rows}</div>
    ${logs}`;
}

// import { DB, createStorage, getStorage, updateStorage } from "./storage.js";
// import { randomID } from "./utils.js";

class Note {
  constructor(title, description, current = false, nextNote = null) {
    this.title = title;
    this.description = description;
    this.current = current;
    this.nextNote = nextNote;
  }
}

const createLinkedList = (notes) => {
  const noteMap = new Map();
  notes.forEach((note, index) => {
    noteMap.set(
      note.title,
      new Note(note.title, note.description, note.current)
    );
  });

  notes.forEach((note, index) => {
    if (note.next) {
      noteMap.get(note.title).nextNote = noteMap.get(note.next);
    }
  });

  let headOfLinkedList = null;
  notes.forEach((note, index) => {
    if (!note.previous) {
      headOfLinkedList = noteMap.get(note.title);
    }
  });

  return headOfLinkedList; // Return the head of the linked list
};

async function handleNotePage(title) {
  let state = await DB.states.where("name").equals("general").last();
  if (!state) {
    // Be resilient if the general nav state was never created (e.g. a learner
    // who reached the Journey board directly).
    await createStorage("states", {
      id: randomID(),
      name: "general",
      previous: "scratch",
      current: "note",
      next: "note-quiz",
    });
  } else {
    state.current = "note";
    state.previous = "scratch";
    state.next = "note-quiz";
    await updateStorage("states", state);
  }
  notePage(document.querySelector("main"), title);
}

// Turn a module title into its note filename, e.g.
// "History of the Web" -> "history_of_the_web".
function slugifyTitle(title) {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "");
}

// Fetches a module's markdown note. Returns null (not an error) when the
// lesson file doesn't exist yet, so the caller can show a friendly fallback.
async function fetchNote(pathName, title) {
  try {
    const response = await fetch(
      `data/notes/${pathName.toLowerCase()}/${slugifyTitle(title)}.md`
    );
    if (!response.ok) return null;
    return await response.text();
  } catch (error) {
    console.log("fetchNote error:", error);
    return null;
  }
}

// Minimal, safe Markdown -> HTML renderer. Escapes all input first (no raw
// HTML injection), then applies a small subset: headings, bold/italic/code,
// links, and unordered lists. Deliberately avoids eval / innerHTML of raw md.
// Escapes text for safe insertion as element content.
function mdEscHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
// Escapes text for safe insertion into a double-quoted attribute.
function mdEscAttr(s) {
  return mdEscHtml(s).replace(/"/g, "&quot;");
}

// Inline markdown: code, bold, italic, links, inline images. HTML is escaped
// first so authored lessons can't inject markup.
function mdInline(s) {
  return mdEscHtml(s)
    .replace(/`([^`]+)`/g, '<code class="gd-code">$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
    .replace(/\*([^*]+)\*/g, "<em>$1</em>")
    // Links (but not the "[..](..)" part of an image — keep a non-"!" prefix).
    .replace(
      /(^|[^!])\[([^\]]+)\]\(([^)\s]+)\)/g,
      '$1<a href="$3" target="_blank" rel="noopener" class="font-bold text-brand-600 underline decoration-2 underline-offset-2">$2</a>'
    );
}

// Friendly, callout-box presets keyed by the `:::type` fence.
const MD_CALLOUTS = {
  tip: { cls: "gd-callout-tip", icon: "bulb", label: "Tip" },
  analogy: { cls: "gd-callout-analogy", icon: "puzzle", label: "Think of it like…" },
  warning: { cls: "gd-callout-warning", icon: "alert", label: "Watch out" },
  key: { cls: "gd-callout-key", icon: "key", label: "Key idea" },
  example: { cls: "gd-callout-example", icon: "flask", label: "Example" },
  project: { cls: "gd-callout-project", icon: "rocket", label: "What you'll build" },
};

function renderCallout(type, inner) {
  const c = MD_CALLOUTS[type] || MD_CALLOUTS.tip;
  // Inner content is rendered as full markdown (paragraphs, lists, code).
  return `<aside class="gd-callout ${c.cls}">
    <p class="gd-callout-title">${icon(c.icon, "w-4 h-4 shrink-0")} ${c.label}</p>
    <div class="gd-callout-body">${markdownToHtml(inner)}</div>
  </aside>`;
}

function renderCodeBlock(code, lang) {
  const label = lang ? lang.toUpperCase() : "CODE";
  return `<div class="gd-codeblock">
    <div class="gd-codeblock-head">${mdEscHtml(label)}</div>
    <pre><code>${mdEscHtml(code)}</code></pre>
  </div>`;
}

// Parses a `:::quiz` block into an interactive check-for-understanding widget.
// Syntax:  Q: question / - option (append " *" to mark correct) / E: explanation
function renderInlineQuiz(buf) {
  let question = "";
  let explanation = "";
  const options = [];
  buf.forEach((l) => {
    const t = l.trim();
    if (/^Q:/i.test(t)) question = t.replace(/^Q:\s*/i, "");
    else if (/^E:/i.test(t)) explanation = t.replace(/^E:\s*/i, "");
    else if (/^[-*]\s+/.test(t)) {
      let txt = t.replace(/^[-*]\s+/, "");
      const correct = /\*\s*$/.test(txt);
      txt = txt.replace(/\s*\*\s*$/, "");
      options.push({ txt, correct });
    }
  });

  const optsHtml = options
    .map(
      (o) =>
        `<button type="button" class="lesson-quiz-option gd-option w-full text-left" data-correct="${o.correct}">${mdInline(
          o.txt
        )}</button>`
    )
    .join("");

  return `<div class="lesson-quiz gd-card-sm my-6"${
    explanation ? ` data-explanation="${mdEscAttr(explanation)}"` : ""
  }>
    <p class="gd-label mb-2 flex items-center gap-1.5">${icon("help", "w-4 h-4")} Quick check</p>
    <p class="font-extrabold mb-3">${mdInline(question)}</p>
    <div class="space-y-2">${optsHtml}</div>
    <div class="lesson-quiz-feedback hidden mt-3 text-sm font-bold"></div>
  </div>`;
}

// Wires an inline quiz widget: reveal correct answer + explanation on first click.
function initLessonQuiz(quizEl) {
  const opts = quizEl.querySelectorAll(".lesson-quiz-option");
  const feedback = quizEl.querySelector(".lesson-quiz-feedback");
  const explanation = quizEl.getAttribute("data-explanation");
  let answered = false;

  opts.forEach((btn) => {
    btn.addEventListener("click", () => {
      if (answered) return;
      answered = true;
      const correct = btn.dataset.correct === "true";
      opts.forEach((o) => {
        o.disabled = true;
        if (o.dataset.correct === "true") o.classList.add("lesson-quiz-correct");
      });
      if (!correct) btn.classList.add("lesson-quiz-wrong");
      feedback.classList.remove("hidden");
      feedback.classList.add(correct ? "text-grass-600" : "text-rose-500");
      feedback.innerHTML =
        `<span class="inline-flex items-center gap-1 align-text-bottom">${icon(
          correct ? "checkCircle" : "xCircle",
          "w-4 h-4"
        )} ${correct ? "Correct!" : "Not quite."}</span> ` +
        (explanation ? `<span class="font-normal text-slate-600">${explanation}</span>` : "");
    });
  });
}

// Parses a `:::reorder` block into a tap-to-sequence exercise. The lines are
// authored in the CORRECT order; the widget shuffles them for the learner.
//   <prompt> / - line (in order) / E: explanation
function renderReorder(buf) {
  let prompt = "";
  let explanation = "";
  const lines = [];
  buf.forEach((l) => {
    const t = l.trim();
    if (/^E:/i.test(t)) explanation = t.replace(/^E:\s*/i, "");
    else if (/^[-*]\s/.test(l.trimStart())) lines.push(l.replace(/^\s*[-*]\s/, ""));
    else if (t && !prompt) prompt = t.replace(/^Q:\s*/i, "");
  });

  const items = lines.map((text, index) => ({ text, index }));
  const chips = shuffle(items)
    .map(
      (it) =>
        `<button type="button" class="reorder-chip" data-index="${it.index}">${mdEscHtml(it.text)}</button>`
    )
    .join("");

  return `<div class="lesson-reorder gd-card-sm my-6"${
    explanation ? ` data-explanation="${mdEscAttr(explanation)}"` : ""
  }>
    <p class="gd-label mb-2"><span aria-hidden="true">🔀</span> Put the code in order</p>
    <p class="font-extrabold mb-1">${mdInline(prompt)}</p>
    <p class="text-xs text-slate-500 mb-3">Tap the lines in the right order. Tap a placed line to send it back.</p>
    <div class="reorder-answer space-y-1.5 mb-2 min-h-[2.75rem] rounded-xl border-2 border-dashed border-slate-200 p-2"></div>
    <div class="reorder-source space-y-1.5">${chips}</div>
    <div class="flex gap-2 mt-3">
      <button type="button" class="reorder-check gd-btn gd-btn-primary !py-2 !px-4 !text-sm">Check</button>
      <button type="button" class="reorder-reset gd-btn gd-btn-secondary !py-2 !px-4 !text-sm">Reset</button>
    </div>
    <div class="reorder-feedback hidden mt-3 text-sm font-bold"></div>
  </div>`;
}

// Parses a `:::fill` block into a fill-in-the-blank exercise: a code snippet
// containing `___`, and option chips for the missing token.
//   Q: prompt / `code with ___` / - option (correct marked " *") / E: explanation
function renderFill(buf) {
  let prompt = "";
  let explanation = "";
  let code = "";
  const options = [];
  buf.forEach((l) => {
    const t = l.trim();
    if (/^E:/i.test(t)) explanation = t.replace(/^E:\s*/i, "");
    else if (/^[-*]\s/.test(t)) {
      let txt = t.replace(/^[-*]\s+/, "");
      const correct = /\*\s*$/.test(txt);
      txt = txt.replace(/\s*\*\s*$/, "");
      options.push({ txt, correct });
    } else if (t.includes("___")) code = t.replace(/^`/, "").replace(/`$/, "");
    else if (t && !prompt) prompt = t.replace(/^Q:\s*/i, "");
  });

  const codeHtml = code
    .split("___")
    .map((p) => mdEscHtml(p))
    .join('<span class="fill-slot" data-state="empty">______</span>');

  const optsHtml = options
    .map(
      (o) =>
        `<button type="button" class="fill-option" data-correct="${o.correct}" data-text="${mdEscAttr(
          o.txt
        )}">${mdEscHtml(o.txt)}</button>`
    )
    .join("");

  return `<div class="lesson-fill gd-card-sm my-6"${
    explanation ? ` data-explanation="${mdEscAttr(explanation)}"` : ""
  }>
    <p class="gd-label mb-2"><span aria-hidden="true">✏️</span> Fill in the blank</p>
    <p class="font-extrabold mb-3">${mdInline(prompt)}</p>
    <div class="gd-codeblock"><div class="gd-codeblock-head">CODE</div><pre><code>${codeHtml}</code></pre></div>
    <div class="flex flex-wrap gap-2 mt-3">${optsHtml}</div>
    <div class="lesson-fill-feedback hidden mt-3 text-sm font-bold"></div>
  </div>`;
}

// Parses a `:::predict` block: a fenced code snippet plus multiple-choice
// options for its output. It reuses the inline-quiz markup/interaction
// (initLessonQuiz), so no extra wiring is needed.
//   ```lang ... ``` / Q: prompt (optional) / - option (correct " *") / E: ...
function renderPredict(buf) {
  let prompt = "What does this code print?";
  let explanation = "";
  let lang = "";
  const codeLines = [];
  const options = [];
  let inCode = false;
  buf.forEach((l) => {
    const t = l.trim();
    if (/^```/.test(t)) {
      if (!inCode) {
        inCode = true;
        lang = t.slice(3).trim();
      } else {
        inCode = false;
      }
      return;
    }
    if (inCode) {
      codeLines.push(l);
      return;
    }
    if (/^E:/i.test(t)) explanation = t.replace(/^E:\s*/i, "");
    else if (/^[-*]\s/.test(t)) {
      let txt = t.replace(/^[-*]\s+/, "");
      const correct = /\*\s*$/.test(txt);
      txt = txt.replace(/\s*\*\s*$/, "");
      options.push({ txt, correct });
    } else if (/^Q:/i.test(t)) prompt = t.replace(/^Q:\s*/i, "");
  });

  const optsHtml = options
    .map(
      (o) =>
        `<button type="button" class="lesson-quiz-option gd-option w-full text-left" data-correct="${o.correct}">${mdInline(
          o.txt
        )}</button>`
    )
    .join("");

  return `<div class="lesson-quiz gd-card-sm my-6"${
    explanation ? ` data-explanation="${mdEscAttr(explanation)}"` : ""
  }>
    <p class="gd-label mb-2"><span aria-hidden="true">🔮</span> Predict the output</p>
    <p class="font-extrabold mb-3">${mdInline(prompt)}</p>
    <div class="gd-codeblock"><div class="gd-codeblock-head">${
      lang ? mdEscHtml(lang.toUpperCase()) : "CODE"
    }</div><pre><code>${mdEscHtml(codeLines.join("\n"))}</code></pre></div>
    <div class="space-y-2 mt-3">${optsHtml}</div>
    <div class="lesson-quiz-feedback hidden mt-3 text-sm font-bold"></div>
  </div>`;
}

// Parses a `:::match` block into a tap-to-match pairs exercise.
//   Q: prompt / - left | right / E: explanation
function renderMatch(buf) {
  let prompt = "Match the pairs.";
  let explanation = "";
  const pairs = [];
  buf.forEach((l) => {
    const t = l.trim();
    if (/^E:/i.test(t)) explanation = t.replace(/^E:\s*/i, "");
    else if (/^[-*]\s/.test(t) && t.includes("|")) {
      const body = t.replace(/^[-*]\s+/, "");
      const idx = body.indexOf("|");
      pairs.push({ left: body.slice(0, idx).trim(), right: body.slice(idx + 1).trim() });
    } else if (/^Q:/i.test(t)) prompt = t.replace(/^Q:\s*/i, "");
    else if (t && prompt === "Match the pairs.") prompt = t;
  });

  const lefts = pairs
    .map(
      (p, i) =>
        `<button type="button" class="match-item match-left" data-pair="${i}">${mdInline(p.left)}</button>`
    )
    .join("");
  const rights = shuffle(pairs.map((p, i) => ({ i, right: p.right })))
    .map(
      (r) =>
        `<button type="button" class="match-item match-right" data-pair="${r.i}">${mdInline(r.right)}</button>`
    )
    .join("");

  return `<div class="lesson-match gd-card-sm my-6"${
    explanation ? ` data-explanation="${mdEscAttr(explanation)}"` : ""
  }>
    <p class="gd-label mb-2"><span aria-hidden="true">🔗</span> Match the pairs</p>
    <p class="font-extrabold mb-3">${mdInline(prompt)}</p>
    <div class="grid grid-cols-2 gap-2">
      <div class="space-y-2">${lefts}</div>
      <div class="space-y-2">${rights}</div>
    </div>
    <div class="lesson-match-feedback hidden mt-3 text-sm font-bold"></div>
  </div>`;
}

// Wires a match exercise: select a left, then a right; a correct pair locks
// green, a wrong pair flashes red.
function initMatch(el) {
  const feedback = el.querySelector(".lesson-match-feedback");
  const explanation = el.getAttribute("data-explanation");
  const total = el.querySelectorAll(".match-left").length;
  let selLeft = null;
  let selRight = null;
  let matched = 0;

  const tryMatch = () => {
    if (!selLeft || !selRight) return;
    const a = selLeft;
    const b = selRight;
    selLeft = null;
    selRight = null;
    if (a.dataset.pair === b.dataset.pair) {
      [a, b].forEach((x) => {
        x.classList.remove("match-selected");
        x.classList.add("match-correct");
        x.disabled = true;
      });
      matched++;
      if (matched === total) {
        feedback.classList.remove("hidden");
        feedback.className = "lesson-match-feedback mt-3 text-sm font-bold text-grass-600";
        feedback.innerHTML =
          "✅ All matched! " +
          (explanation ? `<span class="font-normal text-slate-600">${explanation}</span>` : "");
      }
    } else {
      [a, b].forEach((x) => {
        x.classList.remove("match-selected");
        x.classList.add("match-wrong");
      });
      setTimeout(() => [a, b].forEach((x) => x.classList.remove("match-wrong")), 500);
    }
  };

  el.addEventListener("click", (e) => {
    const item = e.target.closest(".match-item");
    if (!item || item.disabled) return;
    const isLeft = item.classList.contains("match-left");
    const sel = isLeft ? selLeft : selRight;
    // toggle off if re-tapping the current selection
    if (sel === item) {
      item.classList.remove("match-selected");
      if (isLeft) selLeft = null;
      else selRight = null;
      return;
    }
    if (sel) sel.classList.remove("match-selected");
    item.classList.add("match-selected");
    if (isLeft) selLeft = item;
    else selRight = item;
    tryMatch();
  });
}

// Wires a fill-in-the-blank: clicking an option fills the slot and reveals
// whether it was right, with the explanation.
function initFill(el) {
  const slot = el.querySelector(".fill-slot");
  const options = el.querySelectorAll(".fill-option");
  const feedback = el.querySelector(".lesson-fill-feedback");
  const explanation = el.getAttribute("data-explanation");
  let answered = false;

  options.forEach((btn) => {
    btn.addEventListener("click", () => {
      if (answered) return;
      answered = true;
      const correct = btn.dataset.correct === "true";
      if (slot) {
        slot.textContent = btn.dataset.text;
        slot.dataset.state = correct ? "correct" : "wrong";
      }
      options.forEach((o) => {
        o.disabled = true;
        if (o.dataset.correct === "true") o.classList.add("fill-option-correct");
      });
      if (!correct) btn.classList.add("fill-option-wrong");
      feedback.classList.remove("hidden");
      feedback.classList.add(correct ? "text-grass-600" : "text-rose-500");
      feedback.innerHTML =
        (correct ? "✅ Correct! " : "❌ Not quite. ") +
        (explanation ? `<span class="font-normal text-slate-600">${explanation}</span>` : "");
    });
  });
}

// Wires a reorder widget: tap to move chips between source and answer, then
// Check validates the order against the original sequence.
function initReorder(el) {
  const source = el.querySelector(".reorder-source");
  const answer = el.querySelector(".reorder-answer");
  const feedback = el.querySelector(".reorder-feedback");
  const explanation = el.getAttribute("data-explanation");
  const clearMarks = () =>
    el
      .querySelectorAll(".reorder-chip")
      .forEach((c) => c.classList.remove("reorder-correct", "reorder-wrong"));

  el.addEventListener("click", (e) => {
    const chip = e.target.closest(".reorder-chip");
    if (chip) {
      (chip.parentElement === source ? answer : source).appendChild(chip);
      clearMarks();
      feedback.classList.add("hidden");
      return;
    }
    if (e.target.closest(".reorder-reset")) {
      el.querySelectorAll(".reorder-answer .reorder-chip").forEach((c) => source.appendChild(c));
      clearMarks();
      feedback.classList.add("hidden");
      return;
    }
    if (e.target.closest(".reorder-check")) {
      const placed = [...answer.querySelectorAll(".reorder-chip")];
      const total = el.querySelectorAll(".reorder-chip").length;
      clearMarks();
      let correct = placed.length === total;
      placed.forEach((c, i) => {
        const ok = Number(c.dataset.index) === i;
        c.classList.add(ok ? "reorder-correct" : "reorder-wrong");
        if (!ok) correct = false;
      });
      feedback.classList.remove("hidden");
      feedback.className =
        "reorder-feedback mt-3 text-sm font-bold " + (correct ? "text-grass-600" : "text-rose-500");
      const msg = correct
        ? "✅ Perfect order!"
        : placed.length < total
        ? "Place all the lines first."
        : "❌ Not quite — try again.";
      feedback.innerHTML =
        msg +
        (correct && explanation
          ? ` <span class="font-normal text-slate-600">${explanation}</span>`
          : "");
    }
  });
}

// Lightweight Markdown -> HTML for lessons. Supports headings, ordered &
// unordered lists, bold/italic/inline-code, links, block images (figures),
// fenced code blocks (```lang), callout boxes (:::tip/analogy/warning/key/
// example) and interactive inline quizzes (:::quiz). All input is escaped, so
// authored lessons cannot inject raw HTML or scripts.
function markdownToHtml(md) {
  const lines = md.split("\n");
  const html = [];
  let i = 0;
  let inUL = false;
  let inOL = false;
  const closeLists = () => {
    if (inUL) {
      html.push("</ul>");
      inUL = false;
    }
    if (inOL) {
      html.push("</ol>");
      inOL = false;
    }
  };

  while (i < lines.length) {
    const line = lines[i].trimEnd();

    // Fenced code block: ```lang ... ```
    if (/^```/.test(line.trim())) {
      closeLists();
      const lang = line.trim().slice(3).trim();
      const buf = [];
      i++;
      while (i < lines.length && !/^```/.test(lines[i].trim())) {
        buf.push(lines[i]);
        i++;
      }
      i++; // skip closing fence
      html.push(renderCodeBlock(buf.join("\n"), lang));
      continue;
    }

    // Fenced callout / quiz: :::type ... :::
    const fence = line.trim().match(/^:::\s*([a-z]+)\s*$/i);
    if (fence) {
      closeLists();
      const type = fence[1].toLowerCase();
      const buf = [];
      i++;
      while (i < lines.length && lines[i].trim() !== ":::") {
        buf.push(lines[i]);
        i++;
      }
      i++; // skip closing :::
      let block;
      if (type === "quiz") block = renderInlineQuiz(buf);
      else if (type === "reorder") block = renderReorder(buf);
      else if (type === "fill") block = renderFill(buf);
      else if (type === "predict") block = renderPredict(buf);
      else if (type === "match") block = renderMatch(buf);
      else block = renderCallout(type, buf.join("\n"));
      html.push(block);
      continue;
    }

    // Standalone image -> figure with caption.
    const img = line.match(/^!\[([^\]]*)\]\(([^)\s]+)\)\s*$/);
    if (img) {
      closeLists();
      const cap = img[1]
        ? `<figcaption class="text-center text-sm text-slate-500 mt-2">${mdInline(
            img[1]
          )}</figcaption>`
        : "";
      html.push(
        `<figure class="my-6"><img src="${mdEscAttr(
          img[2]
        )}" alt="${mdEscAttr(img[1])}" loading="lazy" class="mx-auto rounded-2xl max-w-full" />${cap}</figure>`
      );
      i++;
      continue;
    }

    const heading = line.match(/^(#{1,4})\s+(.*)$/);
    const ul = line.match(/^[-*]\s+(.*)$/);
    const ol = line.match(/^\d+\.\s+(.*)$/);

    if (heading) {
      closeLists();
      const level = heading[1].length;
      const sizes = { 1: "text-3xl", 2: "text-2xl", 3: "text-xl", 4: "text-lg" };
      html.push(
        `<h${level} class="${sizes[level]} font-extrabold mt-7 mb-3">${mdInline(
          heading[2]
        )}</h${level}>`
      );
    } else if (ul) {
      if (inOL) {
        html.push("</ol>");
        inOL = false;
      }
      if (!inUL) {
        html.push('<ul class="list-disc pl-6 mb-4 space-y-1.5">');
        inUL = true;
      }
      html.push(`<li>${mdInline(ul[1])}</li>`);
    } else if (ol) {
      if (inUL) {
        html.push("</ul>");
        inUL = false;
      }
      if (!inOL) {
        html.push('<ol class="list-decimal pl-6 mb-4 space-y-1.5">');
        inOL = true;
      }
      html.push(`<li>${mdInline(ol[1])}</li>`);
    } else if (line.trim() === "") {
      closeLists();
    } else {
      closeLists();
      html.push(`<p class="mb-4 leading-relaxed">${mdInline(line)}</p>`);
    }
    i++;
  }
  closeLists();
  return html.join("\n");
}

async function scratchPage(htmlEl) {
  const state = await DB.states.where("name").equals("general").last();
  const user = (await DB.users.toArray())[0];
  const pathName = await getActivePath();
  const notes = await DB.paths
    .where("path_name")
    .equals(pathName)
    .toArray();

  const scratchNotes = createLinkedList(notes);
  // Look up each module's stored record (for is_completed) by title.
  const recordByTitle = {};
  notes.forEach((n) => (recordByTitle[n.title] = n));

  let currentNote = scratchNotes;
  const scratchNotesArray = [];

  while (currentNote !== null) {
    const record = recordByTitle[currentNote.title] || {};
    let badge = "";
    let button;
    let accent = "border-slate-100";
    let stepIconName = "lock";
    let stepBox = "bg-slate-100 text-slate-400";
    if (record.is_completed) {
      accent = "border-grass-200";
      stepIconName = "check";
      stepBox = "bg-grass-100 text-grass-600";
      badge = `<span class="gd-chip gd-chip-grass mb-3">${icon("check", "w-3.5 h-3.5")} Completed</span>`;
      button = `<button onclick="handleNotePage('${currentNote.title.replace(/'/g, "\\'")}')" class="gd-btn gd-btn-secondary gd-btn-block">Review lesson</button>`;
    } else if (currentNote.current) {
      accent = "border-brand-300 ring-2 ring-brand-200";
      stepIconName = "play";
      stepBox = "bg-brand-100 text-brand-600";
      badge = `<span class="gd-chip gd-chip-brand mb-3">In Progress</span>`;
      button = `<button onclick="handleNotePage()" class="gd-btn gd-btn-primary gd-btn-block">Start lesson</button>`;
    } else {
      badge = `<span class="gd-chip gd-chip-slate mb-3">${icon("lock", "w-3.5 h-3.5")} Locked</span>`;
      button = `<button class="gd-btn gd-btn-secondary gd-btn-block" disabled>Locked</button>`;
    }

    scratchNotesArray.push(`
    <div class="gd-card !p-6 flex flex-col border-2 ${accent}">
      <div class="flex items-start justify-between gap-3 mb-1">
        <div class="grid h-10 w-10 shrink-0 place-items-center rounded-2xl ${stepBox}">${icon(stepIconName, "w-5 h-5")}</div>
        ${badge}
      </div>
      <div class="flex-1">
        <h2 class="text-xl font-extrabold mb-2">${currentNote.title}</h2>
        <p class="text-slate-600 text-sm mb-5">${currentNote.description}</p>
      </div>
      ${button}
    </div>`);

    currentNote = currentNote.nextNote;
  }

  if (state.current === "scratch") {
    const completedCount = notes.filter((n) => n.is_completed).length;
    const pct = notes.length ? Math.round((completedCount / notes.length) * 100) : 0;
    htmlEl.innerHTML = `
      <div class="animate-fade-up space-y-6">
        <div class="gd-card">
          <div class="flex flex-wrap items-center justify-between gap-4">
            <div>
              <span class="gd-chip gd-chip-brand mb-2">${icon("book", "w-3.5 h-3.5")} Your learning path</span>
              <h1 class="text-2xl font-extrabold">Keep building, one lesson at a time</h1>
            </div>
            <div class="text-right">
              <p class="text-2xl font-extrabold text-grass-600">${completedCount}/${notes.length}</p>
              <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Completed</p>
            </div>
          </div>
          <div class="gd-progress h-3 mt-4">
            <div class="gd-progress-fill" style="width: ${pct}%"></div>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          ${scratchNotesArray.join("")}
        </div>
      </div>
        `;
  } else if (state.current === "note") {
    return notePage(document.querySelector("main"));
  }
}

// Returns a path name that actually has modules: the user's preference if it
// has lessons, otherwise falls back to "frontend" (the seeded default path).
async function resolveLessonPath(preference) {
  const count = await DB.paths
    .where("path_name")
    .equals((preference || "").toLowerCase())
    .count();
  return count > 0 ? preference.toLowerCase() : "frontend";
}

// The world the learner is currently playing: the one chosen on the Worlds
// page (persisted in meta), falling back to their onboarding preference.
async function getActivePath() {
  let selected = null;
  if (typeof getMeta === "function") {
    try {
      selected = await getMeta("active-path", null);
    } catch (e) {
      /* ignore */
    }
  }
  if (selected) {
    const count = await DB.paths.where("path_name").equals(selected).count();
    if (count > 0) return selected;
  }
  const user = (await DB.users.toArray())[0];
  return resolveLessonPath(user ? user.preference : "frontend");
}

async function backToModules() {
  // The Journey board is the canonical "see your whole path" view. Lessons are
  // rendered into <main> without changing the hash, so re-render directly when
  // we're already on #journey.
  if (window.location.hash !== "#journey") {
    window.location.hash = "journey";
  } else {
    JourneyPage(document.querySelector("main"));
  }
}

// The Journey board — the learner's path as a vertical roadmap of "tickets",
// with done / current / locked stops and overall progress.
// Renders the path as a game-style world map: a winding road that snakes
// through a themed landscape, with level nodes along it and the travelled
// portion of the road filled in to show progress.
function renderQuestMap(ordered, recordByTitle) {
  const n = ordered.length;
  const W = 400;
  const spacing = 132;
  const topPad = 84;
  const bottomPad = 110;
  const H = topPad + (n - 1) * spacing + bottomPad;
  const clamp = (v, lo, hi) => Math.max(lo, Math.min(hi, v));

  const pts = ordered.map((_, i) => ({
    x: clamp(200 + 120 * Math.sin(i * 0.9 + 0.3), 86, 314),
    y: topPad + i * spacing,
  }));

  // How far the learner has travelled (last completed, or the current node).
  let reachedIndex = -1;
  ordered.forEach((no, i) => {
    const rec = recordByTitle[no.title] || {};
    if (rec.is_completed || no.current) reachedIndex = i;
  });

  const pathThrough = (slice) => {
    if (slice.length < 2) return slice.length ? `M ${slice[0].x} ${slice[0].y}` : "";
    let d = `M ${slice[0].x} ${slice[0].y}`;
    for (let i = 1; i < slice.length; i++) {
      const p = slice[i - 1];
      const c = slice[i];
      const my = (p.y + c.y) / 2;
      d += ` C ${p.x} ${my}, ${c.x} ${my}, ${c.x} ${c.y}`;
    }
    return d;
  };
  const fullD = pathThrough(pts);
  const traveledD = pathThrough(pts.slice(0, Math.max(1, reachedIndex + 1)));

  // Decorative scenery.
  const clouds = `<g fill="#ffffff" opacity="0.9">
      <ellipse cx="86" cy="46" rx="32" ry="15"/><ellipse cx="114" cy="40" rx="24" ry="13"/>
      <ellipse cx="318" cy="86" rx="28" ry="13"/><ellipse cx="342" cy="80" rx="20" ry="11"/>
      <ellipse cx="70" cy="${topPad + (n - 1) * spacing - 30}" rx="26" ry="12"/>
    </g>`;
  const bush = (x, y, s) =>
    `<g transform="translate(${x} ${y}) scale(${s})"><circle r="15" fill="#7ad23e"/><circle cx="13" cy="3" r="11" fill="#9ce26a"/><circle cx="-13" cy="3" r="11" fill="#9ce26a"/></g>`;
  const decorations = pts
    .map((p, i) => (i % 2 ? bush(p.x < 200 ? 350 : 50, p.y + 36, 1) : ""))
    .join("");

  const nodes = ordered
    .map((no, i) => {
      const p = pts[i];
      const rec = recordByTitle[no.title] || {};
      const isProject = /^project/i.test(no.title);
      const done = rec.is_completed;
      const current = no.current;
      const titleEsc = no.title.replace(/'/g, "\\'");

      let badge, glyph, onclick, dis = "", labelCls;
      if (done) {
        badge = "qm-node qm-done";
        glyph = icon("check", "w-6 h-6");
        onclick = `onclick="handleNotePage('${titleEsc}')"`;
        labelCls = "text-slate-700";
      } else if (current) {
        badge = "qm-node qm-current";
        glyph = isProject ? icon("rocket", "w-6 h-6") : icon("play", "w-6 h-6");
        onclick = `onclick="handleNotePage()"`;
        labelCls = "text-brand-700";
      } else {
        badge = "qm-node qm-locked";
        glyph = icon("lock", "w-5 h-5");
        onclick = "";
        dis = "disabled";
        labelCls = "text-slate-400";
      }

      return `<foreignObject x="${p.x - 64}" y="${p.y - 30}" width="128" height="124">
        <div xmlns="http://www.w3.org/1999/xhtml" class="relative flex flex-col items-center">
          <button ${onclick} ${dis} class="${badge}" aria-label="${no.title}">
            ${current ? '<span class="absolute inset-0 rounded-full bg-brand-400/40 animate-ping"></span>' : ""}
            <span class="relative grid h-full w-full place-items-center text-white">${glyph}</span>
          </button>
          <span class="mt-1.5 text-center text-[11px] font-bold leading-tight ${labelCls}">${no.title}${
        isProject ? ' <span aria-hidden="true">🚀</span>' : ""
      }</span>
        </div>
      </foreignObject>`;
    })
    .join("");

  const goalY = topPad + (n - 1) * spacing + 64;

  return `<svg viewBox="0 0 ${W} ${H}" class="w-full h-auto" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="Your quest map">
    <defs>
      <linearGradient id="qm-sky" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" stop-color="#ede9ff"/><stop offset="0.55" stop-color="#f5f3ff"/><stop offset="1" stop-color="#e9fbe1"/>
      </linearGradient>
    </defs>
    <rect width="${W}" height="${H}" rx="28" fill="url(#qm-sky)"/>
    ${clouds}
    ${decorations}
    <path d="${fullD}" fill="none" stroke="#e2ddff" stroke-width="24" stroke-linecap="round" stroke-linejoin="round"/>
    <path d="${traveledD}" fill="none" stroke="#7ad23e" stroke-width="24" stroke-linecap="round" stroke-linejoin="round"/>
    <path d="${fullD}" fill="none" stroke="#ffffff" stroke-width="3" stroke-dasharray="1 14" stroke-linecap="round" opacity="0.8"/>
    <text x="200" y="${goalY}" text-anchor="middle" font-size="34">🏁</text>
    ${nodes}
  </svg>`;
}

// Each path is presented as a selectable "world". Order here drives the grid.
const WORLDS = [
  {
    path: "frontend",
    name: "Frontend",
    emoji: "🎨",
    tagline: "Build what users see and touch",
    blurb: "HTML, CSS & JavaScript — ship real, interactive web pages.",
  },
  {
    path: "backend",
    name: "Backend",
    emoji: "🗄️",
    tagline: "Power apps from behind the scenes",
    blurb: "Servers, databases & APIs — the engine room of every app.",
  },
  {
    path: "fullstack",
    name: "Fullstack",
    emoji: "🔗",
    tagline: "Connect front and back into apps",
    blurb: "Tie it all together into complete, shippable products.",
  },
  {
    path: "c",
    name: "C",
    emoji: "⚙️",
    tagline: "Program close to the metal",
    blurb: "Pointers, memory & control flow — the roots of computing.",
  },
  {
    path: "java",
    name: "Java",
    emoji: "☕",
    tagline: "Robust, portable applications",
    blurb: "Object-oriented thinking that runs everywhere.",
  },
  {
    path: "linux",
    name: "Linux & Bash",
    emoji: "🐧",
    tagline: "Command the machine directly",
    blurb: "The shell, the filesystem & scripting — a developer superpower.",
  },
];

// Set the active world and jump straight to its quest map.
async function selectWorld(pathName) {
  if (typeof setMeta === "function") {
    try {
      await setMeta("active-path", pathName);
    } catch (e) {
      /* ignore */
    }
  }
  window.location.hash = "journey";
}

async function WorldsPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }

  const allPaths = await DB.paths.toArray();
  const stats = {};
  allPaths.forEach((p) => {
    const s = (stats[p.path_name] = stats[p.path_name] || { total: 0, done: 0 });
    s.total += 1;
    if (p.is_completed) s.done += 1;
  });

  let active = "frontend";
  try {
    active = await getActivePath();
  } catch (e) {
    /* default */
  }

  const cards = WORLDS.map((w) => {
    const s = stats[w.path] || { total: 0, done: 0 };
    const pct = s.total ? Math.round((s.done / s.total) * 100) : 0;
    const isActive = w.path === active;
    const label = s.done === 0 ? "Start" : pct === 100 ? "Replay" : "Continue";
    return `
      <button type="button" onclick="selectWorld('${w.path}')"
        class="gd-card group text-left flex flex-col gap-3 transition-transform hover:-translate-y-1 ${
          isActive ? "ring-2 ring-brand-400" : ""
        }">
        <div class="flex items-start justify-between">
          <div class="grid h-14 w-14 place-items-center rounded-2xl bg-brand-100 text-3xl">${w.emoji}</div>
          ${
            isActive
              ? '<span class="gd-chip gd-chip-brand !text-[10px] !py-0.5">▶ Playing</span>'
              : ""
          }
        </div>
        <div>
          <h3 class="text-lg font-extrabold">${w.name}</h3>
          <p class="text-sm text-slate-500">${w.blurb}</p>
        </div>
        <div class="gd-progress h-2 mt-auto"><div class="gd-progress-fill" style="width:${pct}%"></div></div>
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold text-slate-500">${s.done}/${s.total} complete</span>
          <span class="gd-btn gd-btn-primary !py-1.5 !px-4 !text-xs group-hover:brightness-105">${label}</span>
        </div>
      </button>`;
  }).join("");

  htmlEl.innerHTML = `
    <div class="max-w-5xl mx-auto animate-fade-up">
      <div class="text-center mb-7">
        <span class="gd-chip gd-chip-brand mb-2">${icon("map", "w-3.5 h-3.5")} Choose your world</span>
        <h1 class="text-2xl sm:text-3xl font-extrabold">Worlds</h1>
        <p class="text-slate-500 mt-1 max-w-md mx-auto">Each world is a full adventure with its own quest map. Pick one to play — your progress is saved, so switch anytime.</p>
      </div>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">${cards}</div>
    </div>`;
}

async function JourneyPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }

  const pathName = await getActivePath();
  const notes = await DB.paths.where("path_name").equals(pathName).toArray();
  if (!notes.length) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <p class="text-slate-600 mb-5">No lessons are available for this path yet.</p>
          <a href="#test" class="gd-btn gd-btn-primary">Take a quiz instead</a>
        </div>
      </div>`;
    return;
  }

  const recordByTitle = {};
  notes.forEach((n) => (recordByTitle[n.title] = n));

  // Order via the linked list.
  const ordered = [];
  let cursor = createLinkedList(notes);
  while (cursor) {
    ordered.push(cursor);
    cursor = cursor.nextNote;
  }

  const completed = notes.filter((n) => n.is_completed).length;
  const pct = Math.round((completed / notes.length) * 100);
  const pathLabel = pathName.charAt(0).toUpperCase() + pathName.slice(1);

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-6">
      <div class="gd-card">
        <div class="flex flex-wrap items-center justify-between gap-4 mb-4">
          <div>
            <span class="gd-chip gd-chip-brand mb-2">${icon("book", "w-3.5 h-3.5")} ${pathLabel} adventure</span>
            <h1 class="text-2xl font-extrabold">Your quest map</h1>
            <a href="#worlds" class="inline-flex items-center gap-1 mt-1 text-sm font-bold text-brand-600 hover:text-brand-700">${icon("map", "w-3.5 h-3.5")} Switch world</a>
          </div>
          <div class="text-right">
            <p class="text-2xl font-extrabold text-brand-600">${completed}/${notes.length}</p>
            <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Completed</p>
          </div>
        </div>
        <div class="gd-progress h-3"><div class="gd-progress-fill" style="width: ${pct}%"></div></div>
      </div>
      <div class="max-w-md mx-auto">${renderQuestMap(ordered, recordByTitle)}</div>
    </div>`;

  // First time on this path's board: Pixel briefs the mission (once per path).
  if (
    typeof getMeta === "function" &&
    typeof pixelCelebrate === "function" &&
    typeof pathMission === "function"
  ) {
    const key = "mission-" + pathName;
    const seen = await getMeta(key, false);
    if (!seen) {
      await setMeta(key, true);
      const m = pathMission(pathName);
      pixelCelebrate({
        title: `Welcome to the ${m.crew}! 🎉`,
        message: m.mission,
        mood: "happy",
        cta: "Let's get to work",
      });
    }
  }
}

async function notePage(htmlEl, requestedTitle) {
  const state = await DB.states.where("name").equals("general").last();
  if (state.current !== "note") return;

  const user = (await DB.users.toArray())[0];
  const pathName = await getActivePath();
  const current = requestedTitle
    ? await DB.paths
        .where("path_name")
        .equals(pathName)
        .and((c) => c.title === requestedTitle)
        .first()
    : await DB.paths
        .where("path_name")
        .equals(pathName)
        .and((c) => c.current)
        .last();

  // Reviewing an already-completed module (not the active one): read-only.
  const isReview = !!requestedTitle && !!current && !current.current;

  if (!current) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="grid h-14 w-14 mx-auto mb-3 place-items-center rounded-2xl bg-slate-100 text-slate-400">${icon("book", "w-7 h-7")}</div>
          <p class="text-slate-600 mb-5">No lessons are available for this path yet.</p>
          <button id="back-to-modules" class="gd-btn gd-btn-primary">Back to Modules</button>
        </div>
      </div>`;
    document.querySelector("#back-to-modules").addEventListener("click", backToModules);
    return;
  }

  const md = await fetchNote(pathName, current.title);
  const esc = (s) => (s || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

  const body = md
    ? markdownToHtml(md)
    : `<h1 class="text-3xl font-bold mb-4">${esc(current.title)}</h1>
       <p class="mb-4 leading-relaxed">${esc(current.description)}</p>
       <p class="text-gray-500 italic">Full lesson content is coming soon.</p>`;

  const resources = Array.isArray(current.resources) ? current.resources : [];
  const resourcesHtml = resources.length
    ? `<div class="mt-8 border-t border-slate-200 pt-5">
         <h3 class="flex items-center gap-2 text-lg font-extrabold mb-3">${icon("link", "w-5 h-5 text-brand-500")} Resources</h3>
         <ul class="space-y-2">
           ${resources
             .map(
               (r) =>
                 `<li><a href="${esc(r.url)}" target="_blank" rel="noopener" class="flex items-center gap-2 rounded-xl bg-brand-50 px-3 py-2 font-bold text-brand-700 hover:bg-brand-100 transition-colors">${icon("link", "w-4 h-4 shrink-0")}${esc(r.title)}</a></li>`
             )
             .join("")}
         </ul>
       </div>`
    : "";

  const hasNext = !!current.next;
  const quizCategory = moduleQuizCategory(current.title, pathName);
  const quizCount = quizCategory
    ? await DB.questions.where("category").equals(quizCategory).count()
    : 0;
  const hasQuiz = quizCount > 0;

  htmlEl.innerHTML = `
    <article class="max-w-3xl mx-auto gd-card animate-fade-up">
      ${
        isReview
          ? `<span class="gd-chip gd-chip-grass mb-4">${icon("check", "w-3.5 h-3.5")} Completed lesson</span>`
          : `<span class="gd-chip gd-chip-brand mb-4">${icon("book", "w-3.5 h-3.5")} Lesson</span>`
      }
      ${body}
      ${resourcesHtml}
      <div class="mt-8 flex flex-wrap justify-between gap-3 border-t border-slate-200 pt-6">
        <button id="back-to-modules" class="gd-btn gd-btn-secondary">${icon("arrowLeft", "w-4 h-4")} Back to Modules</button>
        ${
          isReview
            ? ""
            : `<button id="complete-continue" class="gd-btn gd-btn-primary">
          ${
            hasQuiz
              ? `Take the Quiz ${icon("arrowRight", "w-4 h-4")}`
              : hasNext
              ? `Mark Complete & Continue ${icon("arrowRight", "w-4 h-4")}`
              : `Finish Path ${icon("trophy", "w-4 h-4")}`
          }
        </button>`
        }
      </div>
    </article>`;

  document.querySelector("#back-to-modules").addEventListener("click", backToModules);

  // Activate any inline check-for-understanding quizzes embedded in the lesson.
  htmlEl.querySelectorAll(".lesson-quiz").forEach(initLessonQuiz);
  htmlEl.querySelectorAll(".lesson-reorder").forEach(initReorder);
  htmlEl.querySelectorAll(".lesson-fill").forEach(initFill);
  htmlEl.querySelectorAll(".lesson-match").forEach(initMatch);

  if (isReview) return;

  document.querySelector("#complete-continue").addEventListener("click", () => {
    if (hasQuiz) {
      lessonQuizPage(htmlEl, current, pathName, quizCategory);
    } else {
      advanceFromModule(current, pathName, htmlEl);
    }
  });
}

// Maps a lesson title to a question-bank category for its check-for-understanding
// quiz. Returns null when no category fits (e.g. history or soft-skill modules).
function moduleQuizCategory(title, pathName) {
  const t = title.toLowerCase();
  if (t.startsWith("project")) return null; // projects end by building, not a quiz
  // Single-topic paths map every lesson to their question bank.
  if (pathName === "c") return "c";
  if (pathName === "java") return "java";
  if (pathName === "linux") return "linux";
  if (t.includes("javascript")) return "javascript";
  if (t.includes("python")) return "python";
  if (t.includes("sql") || t.includes("database")) return "sql";
  if (t.includes("css")) return "css";
  if (t.includes("html")) return "html";
  return null;
}

// Marks a module complete and moves the "current" pointer to the next one.
async function advanceFromModule(current, pathName, htmlEl) {
  await DB.paths.update(current.id, { is_completed: true, current: false });
  const next = current.next
    ? await DB.paths
        .where("path_name")
        .equals(pathName)
        .and((c) => c.title === current.next)
        .first()
    : null;

  // Pixel pops up to celebrate the milestone.
  if (typeof pixelCelebrate === "function") {
    const isProject = /^project/i.test(current.title);
    if (!next) {
      await pixelCelebrate({
        title: "Path complete! 🎉",
        message:
          "You finished the whole path — that's a serious milestone. Take a victory lap, then pick your next adventure.",
        cta: "See my journey",
      });
    } else if (isProject) {
      await pixelCelebrate({
        title: "🚀 You shipped it!",
        message: gdPickSafe([
          "That's a real, working thing you built — exactly what devs do all day. On to the next!",
          "Shipped! Founders everywhere just shed a happy tear. Keep that momentum.",
          "You didn't just learn it — you built it. That's the whole game.",
        ]),
        cta: "Keep building",
      });
    } else {
      await pixelCelebrate({
        title: "Lesson complete!",
        message: gdPickSafe([
          "Nice — one step closer. Onward!",
          "Locked in. Let's keep the streak alive.",
          "That's progress. Pixel approves. 👍",
        ]),
        cta: "Continue",
      });
    }
  }

  if (next) {
    await DB.paths.update(next.id, { current: true });
    notePage(htmlEl); // render the next lesson
  } else {
    backToModules(); // path finished — return to the module list
  }
}

// gdPick lives in mascot.js; fall back gracefully if it isn't loaded.
function gdPickSafe(arr) {
  return typeof gdPick === "function" ? gdPick(arr) : arr[0];
}

const LESSON_QUIZ_SIZE = 5;
const LESSON_QUIZ_PASS = 60;

// A short check-for-understanding quiz tied to a lesson. Reuses the Test
// Yourself option renderer and the answer-review builder.
async function lessonQuizPage(htmlEl, module, pathName, category) {
  const all = await DB.questions.where("category").equals(category).toArray();
  const quizQuestions = shuffle(all).slice(0, Math.min(LESSON_QUIZ_SIZE, all.length));

  const questionsHtml = quizQuestions
    .map(
      (q, i) => `
      <div class="lq-question gd-card-sm mb-3">
        <div class="flex items-start gap-3 mb-3">
          <span class="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-brand-100 text-brand-700 font-extrabold text-sm">${i + 1}</span>
          <p class="text-lg font-bold pt-0.5">${escapeHTMLToEntities(q.details.question)}</p>
        </div>
        <form class="space-y-2">${tysRandomizeOptions(q.details.options).join("")}</form>
      </div>`
    )
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up">
      <div class="text-center mb-6">
        <span class="gd-chip gd-chip-brand mb-3">${icon("puzzle", "w-3.5 h-3.5")} Lesson Quiz</span>
        <h1 class="text-2xl font-extrabold">${escapeHTMLToEntities(module.title)}</h1>
        <p class="text-slate-500 mt-1">Answer these ${quizQuestions.length} questions to check your understanding.</p>
      </div>
      ${questionsHtml}
      <div class="mt-5 flex gap-3">
        <button id="lq-skip" class="gd-btn gd-btn-secondary flex-1">Skip</button>
        <button id="lq-submit" class="gd-btn gd-btn-primary flex-[2]">Submit Quiz</button>
      </div>
    </div>`;

  document
    .querySelector("#lq-skip")
    .addEventListener("click", () => advanceFromModule(module, pathName, htmlEl));

  document.querySelector("#lq-submit").addEventListener("click", async () => {
    const selected = [...htmlEl.querySelectorAll(".lq-question")].map(
      (q) => q.querySelector('input[name="option"]:checked')?.value ?? null
    );
    let numCorrect = 0;
    quizQuestions.forEach((q, i) => {
      if (selected[i] === q.details.answer) numCorrect++;
    });
    const numWrong = quizQuestions.length - numCorrect;
    const score = Math.round((numCorrect / quizQuestions.length) * 100);

    // Record the attempt so it feeds XP, streaks, badges and history.
    await createStorage("scores", {
      id: randomID(),
      test_id: "lesson-" + randomID(),
      score,
      numCorrect,
      numWrong,
      details: { questions: quizQuestions, selectedOptions: selected },
      created_at: new Date(),
    });

    lessonQuizResult(htmlEl, module, pathName, {
      score,
      numCorrect,
      questions: quizQuestions,
      selected,
    });
  });
}

function lessonQuizResult(htmlEl, module, pathName, result) {
  const passed = result.score >= LESSON_QUIZ_PASS;
  const reviewData = {
    details: { questions: result.questions, selectedOptions: result.selected },
  };

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up">
      <div class="gd-card">
        <div class="text-center mb-6">
          <div class="grid h-16 w-16 mx-auto mb-3 place-items-center rounded-2xl animate-pop-in ${
            passed ? "bg-grass-50 text-grass-500" : "bg-amber-50 text-amber-500"
          }">${icon(passed ? "trophy" : "target", "w-8 h-8")}</div>
          <h1 class="text-2xl font-extrabold mb-2">${escapeHTMLToEntities(module.title)}</h1>
          <p class="text-slate-600">
            You scored
            <span class="font-extrabold ${passed ? "text-grass-600" : "text-rose-500"}">${result.score}%</span>
            (${result.numCorrect}/${result.questions.length} correct).
            ${passed ? "Great job!" : "Review the explanations below, then continue or try again."}
          </p>
        </div>
        <div id="lq-review" class="mb-6">${buildTysReview(reviewData)}</div>
        <div class="flex gap-3">
          <button id="lq-retry" class="gd-btn gd-btn-secondary flex-1">${icon("refresh", "w-4 h-4")} Retake Quiz</button>
          <button id="lq-continue" class="gd-btn gd-btn-primary flex-1">
            ${module.next ? `Continue ${icon("arrowRight", "w-4 h-4")}` : `Finish Path ${icon("trophy", "w-4 h-4")}`}
          </button>
        </div>
      </div>
    </div>`;

  document
    .querySelector("#lq-retry")
    .addEventListener("click", () =>
      lessonQuizPage(htmlEl, module, pathName, moduleQuizCategory(module.title))
    );
  document
    .querySelector("#lq-continue")
    .addEventListener("click", () => advanceFromModule(module, pathName, htmlEl));
}

// export { scratchPage };

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
  const pathName = await resolveLessonPath(user.preference);
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
async function JourneyPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }

  const pathName = await resolveLessonPath(user.preference);
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

  const stops = ordered
    .map((note, i) => {
      const rec = recordByTitle[note.title] || {};
      const isProject = /^project/i.test(note.title);
      const last = i === ordered.length - 1;
      const titleEsc = note.title.replace(/'/g, "\\'");

      let node, button, cardRing = "", titleCls = "";
      if (rec.is_completed) {
        node = `<div class="grid h-11 w-11 place-items-center rounded-full bg-grass-500 text-white">${icon("check", "w-5 h-5")}</div>`;
        button = `<button onclick="handleNotePage('${titleEsc}')" class="gd-btn gd-btn-secondary !py-2 !px-4 !text-sm shrink-0">Review</button>`;
      } else if (note.current) {
        node = `<div class="grid h-11 w-11 place-items-center rounded-full bg-brand-500 text-white ring-4 ring-brand-200">${icon("play", "w-5 h-5")}</div>`;
        cardRing = "ring-2 ring-brand-300";
        button = `<button onclick="handleNotePage()" class="gd-btn gd-btn-primary !py-2 !px-4 !text-sm shrink-0">${isProject ? "Build" : "Continue"}</button>`;
      } else {
        node = `<div class="grid h-11 w-11 place-items-center rounded-full bg-slate-100 text-slate-400">${icon("lock", "w-4 h-4")}</div>`;
        titleCls = "text-slate-500";
        button = `<button class="gd-btn gd-btn-secondary !py-2 !px-4 !text-sm shrink-0" disabled>Locked</button>`;
      }

      return `
      <div class="relative pl-16 ${last ? "" : "pb-5"}">
        ${last ? "" : '<span class="absolute left-[21px] top-11 -bottom-1 w-0.5 bg-slate-200"></span>'}
        <div class="absolute left-0 top-1">${node}</div>
        <div class="gd-card-sm ${cardRing} ${rec.is_completed || note.current ? "" : "opacity-75"}">
          <div class="flex items-center justify-between gap-3">
            <div class="min-w-0">
              ${isProject ? '<span class="gd-chip gd-chip-brand mb-1.5 !text-[10px]">🚀 Project</span>' : ""}
              <h3 class="font-extrabold ${titleCls}">${note.title}</h3>
              <p class="text-slate-500 text-sm mt-0.5">${note.description}</p>
            </div>
            ${button}
          </div>
        </div>
      </div>`;
    })
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-5">
      <div class="gd-card">
        <div class="flex flex-wrap items-center justify-between gap-4 mb-4">
          <div>
            <span class="gd-chip gd-chip-brand mb-2">${icon("book", "w-3.5 h-3.5")} ${pathLabel} roadmap</span>
            <h1 class="text-2xl font-extrabold">Your journey</h1>
          </div>
          <div class="text-right">
            <p class="text-2xl font-extrabold text-brand-600">${completed}/${notes.length}</p>
            <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Completed</p>
          </div>
        </div>
        <div class="gd-progress h-3"><div class="gd-progress-fill" style="width: ${pct}%"></div></div>
      </div>
      <div class="pt-1">${stops}</div>
    </div>`;
}

async function notePage(htmlEl, requestedTitle) {
  const state = await DB.states.where("name").equals("general").last();
  if (state.current !== "note") return;

  const user = (await DB.users.toArray())[0];
  const pathName = await resolveLessonPath(user.preference);
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
  const quizCategory = moduleQuizCategory(current.title);
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
function moduleQuizCategory(title) {
  const t = title.toLowerCase();
  if (t.startsWith("project")) return null; // projects end by building, not a quiz
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

  if (next) {
    await DB.paths.update(next.id, { current: true });
    notePage(htmlEl); // render the next lesson
  } else {
    backToModules(); // path finished — return to the module list
  }
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

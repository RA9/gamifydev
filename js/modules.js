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
  const state = await DB.states.where("name").equals("general").last();
  state.current = "note";
  state.previous = "scratch";
  state.next = "note-quiz";
  await updateStorage("states", state);
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
function markdownToHtml(md) {
  const esc = (s) =>
    s
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");

  const inline = (s) =>
    esc(s)
      .replace(/`([^`]+)`/g, '<code class="gd-code">$1</code>')
      .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
      .replace(/\*([^*]+)\*/g, "<em>$1</em>")
      .replace(
        /\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g,
        '<a href="$2" target="_blank" rel="noopener" class="text-blue-600 underline">$1</a>'
      );

  const lines = md.split("\n");
  const html = [];
  let inList = false;
  const closeList = () => {
    if (inList) {
      html.push("</ul>");
      inList = false;
    }
  };

  for (const raw of lines) {
    const line = raw.trimEnd();
    const heading = line.match(/^(#{1,4})\s+(.*)$/);
    const listItem = line.match(/^[-*]\s+(.*)$/);

    if (heading) {
      closeList();
      const level = heading[1].length;
      const sizes = { 1: "text-3xl", 2: "text-2xl", 3: "text-xl", 4: "text-lg" };
      html.push(
        `<h${level} class="${sizes[level]} font-bold mt-6 mb-2">${inline(heading[2])}</h${level}>`
      );
    } else if (listItem) {
      if (!inList) {
        html.push('<ul class="list-disc pl-6 mb-4 space-y-1">');
        inList = true;
      }
      html.push(`<li>${inline(listItem[1])}</li>`);
    } else if (line.trim() === "") {
      closeList();
    } else {
      closeList();
      html.push(`<p class="mb-4 leading-relaxed">${inline(line)}</p>`);
    }
  }
  closeList();
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
    let stepIcon = "🔒";
    if (record.is_completed) {
      accent = "border-grass-200";
      stepIcon = "✓";
      badge = `<span class="gd-chip gd-chip-grass mb-3">✓ Completed</span>`;
      button = `<button onclick="handleNotePage('${currentNote.title.replace(/'/g, "\\'")}')" class="gd-btn gd-btn-secondary gd-btn-block">Review lesson</button>`;
    } else if (currentNote.current) {
      accent = "border-brand-300 ring-2 ring-brand-200";
      stepIcon = "▶";
      badge = `<span class="gd-chip gd-chip-brand mb-3">In Progress</span>`;
      button = `<button onclick="handleNotePage()" class="gd-btn gd-btn-grass gd-btn-block">Start lesson</button>`;
    } else {
      badge = `<span class="gd-chip gd-chip-slate mb-3">🔒 Locked</span>`;
      button = `<button class="gd-btn gd-btn-secondary gd-btn-block" disabled>Locked</button>`;
    }

    scratchNotesArray.push(`
    <div class="gd-card !p-6 flex flex-col border-2 ${accent}">
      <div class="flex items-start justify-between gap-3 mb-1">
        <div class="grid h-10 w-10 shrink-0 place-items-center rounded-2xl bg-slate-100 text-lg font-extrabold">${stepIcon}</div>
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
              <span class="gd-chip gd-chip-brand mb-2">📚 Your learning path</span>
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
  const state = await DB.states.where("name").equals("general").last();
  state.current = "scratch";
  state.previous = "note";
  state.next = null;
  await updateStorage("states", state);
  scratchPage(document.querySelector("main"));
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
          <div class="text-4xl mb-2">🚧</div>
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
         <h3 class="text-lg font-extrabold mb-3">📎 Resources</h3>
         <ul class="space-y-2">
           ${resources
             .map(
               (r) =>
                 `<li><a href="${esc(r.url)}" target="_blank" rel="noopener" class="flex items-center gap-2 rounded-xl bg-brand-50 px-3 py-2 font-bold text-brand-700 hover:bg-brand-100 transition-colors"><span>🔗</span>${esc(r.title)}</a></li>`
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
      ${isReview ? '<span class="gd-chip gd-chip-grass mb-4">✓ Completed lesson</span>' : '<span class="gd-chip gd-chip-brand mb-4">📖 Lesson</span>'}
      ${body}
      ${resourcesHtml}
      <div class="mt-8 flex flex-wrap justify-between gap-3 border-t border-slate-200 pt-6">
        <button id="back-to-modules" class="gd-btn gd-btn-secondary">← Back to Modules</button>
        ${
          isReview
            ? ""
            : `<button id="complete-continue" class="gd-btn gd-btn-grass">
          ${hasQuiz ? "Take the Quiz →" : hasNext ? "Mark Complete & Continue →" : "Finish Path 🎉"}
        </button>`
        }
      </div>
    </article>`;

  document.querySelector("#back-to-modules").addEventListener("click", backToModules);

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
  if (t.includes("javascript")) return "javascript";
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
        <span class="gd-chip gd-chip-brand mb-3">🧩 Lesson Quiz</span>
        <h1 class="text-2xl font-extrabold">${escapeHTMLToEntities(module.title)}</h1>
        <p class="text-slate-500 mt-1">Answer these ${quizQuestions.length} questions to check your understanding.</p>
      </div>
      ${questionsHtml}
      <div class="mt-5 flex gap-3">
        <button id="lq-skip" class="gd-btn gd-btn-secondary flex-1">Skip</button>
        <button id="lq-submit" class="gd-btn gd-btn-grass flex-[2]">Submit Quiz</button>
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
          <div class="text-5xl mb-2">${passed ? "🎉" : "💪"}</div>
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
          <button id="lq-retry" class="gd-btn gd-btn-secondary flex-1">Retake Quiz</button>
          <button id="lq-continue" class="gd-btn gd-btn-grass flex-1">
            ${module.next ? "Continue →" : "Finish Path 🎉"}
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

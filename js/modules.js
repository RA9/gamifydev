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
      .replace(/`([^`]+)`/g, '<code class="bg-gray-100 text-pink-600 px-1 rounded">$1</code>')
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
    if (record.is_completed) {
      badge = `<span class="inline-block text-xs font-semibold text-green-700 bg-green-100 rounded px-2 py-1 mb-2">✓ Completed</span>`;
      button = `<button onclick="handleNotePage('${currentNote.title.replace(/'/g, "\\'")}')" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded inline-block">Review</button>`;
    } else if (currentNote.current) {
      badge = `<span class="inline-block text-xs font-semibold text-blue-700 bg-blue-100 rounded px-2 py-1 mb-2">In Progress</span>`;
      button = `<button onclick="handleNotePage()" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-block">Enroll Now</button>`;
    } else {
      badge = `<span class="inline-block text-xs font-semibold text-gray-500 bg-gray-100 rounded px-2 py-1 mb-2">🔒 Locked</span>`;
      button = `<button class="bg-gray-300 text-white font-bold py-2 px-4 rounded inline-block cursor-not-allowed" disabled>Locked</button>`;
    }

    scratchNotesArray.push(`
    <div class="bg-white rounded-lg shadow p-8 flex flex-col">
      <div class="flex-1">
        ${badge}
        <h2 class="text-2xl font-bold mb-2">${currentNote.title}</h2>
        <p class="text-gray-700 mb-4">${currentNote.description}</p>
      </div>
      ${button}
    </div>`);

    currentNote = currentNote.nextNote;
  }

  if (state.current === "scratch") {
    htmlEl.innerHTML = `
        <div class="grid grid-cols-1 sm:grid-cols-2  lg:grid-cols-3 gap-4">
        ${scratchNotesArray.join("")}
      </div><br /> <br /><br />
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
      <div class="max-w-3xl mx-auto bg-white rounded-lg shadow p-8">
        <p class="text-gray-600 mb-4">No lessons are available for this path yet.</p>
        <button id="back-to-modules" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Back to Modules</button>
      </div><br/><br/><br/>`;
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
    ? `<div class="mt-8 border-t border-gray-200 pt-4">
         <h3 class="text-lg font-bold mb-2">Resources</h3>
         <ul class="list-disc pl-6 space-y-1">
           ${resources
             .map(
               (r) =>
                 `<li><a href="${esc(r.url)}" target="_blank" rel="noopener" class="text-blue-600 underline">${esc(r.title)}</a></li>`
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
    <article class="max-w-3xl mx-auto bg-white rounded-lg shadow p-8 prose-gray">
      ${isReview ? '<span class="inline-block text-xs font-semibold text-green-700 bg-green-100 rounded px-2 py-1 mb-4">✓ Completed lesson</span>' : ""}
      ${body}
      ${resourcesHtml}
      <div class="mt-8 flex flex-wrap justify-between gap-4">
        <button id="back-to-modules" class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">Back to Modules</button>
        ${
          isReview
            ? ""
            : `<button id="complete-continue" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
          ${hasQuiz ? "Take the Quiz" : hasNext ? "Mark Complete & Continue" : "Finish Path"}
        </button>`
        }
      </div>
    </article><br/><br/><br/>`;

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
      <div class="lq-question mb-4 border-b border-gray-100 pb-4">
        <p class="text-lg font-semibold mb-2">${i + 1}. ${escapeHTMLToEntities(q.details.question)}</p>
        <form>${tysRandomizeOptions(q.details.options).join("")}</form>
      </div>`
    )
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto bg-white rounded-lg shadow p-8">
      <span class="inline-block text-xs font-semibold text-indigo-700 bg-indigo-100 rounded px-2 py-1 mb-2">Lesson Quiz</span>
      <h1 class="text-2xl font-bold mb-1">${escapeHTMLToEntities(module.title)}</h1>
      <p class="text-gray-500 mb-6">Answer these ${quizQuestions.length} questions to check your understanding.</p>
      ${questionsHtml}
      <div class="mt-4 flex flex-wrap justify-between gap-4">
        <button id="lq-skip" class="bg-gray-400 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded">Skip</button>
        <button id="lq-submit" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">Submit Quiz</button>
      </div>
    </div><br/><br/><br/>`;

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
    <div class="max-w-3xl mx-auto bg-white rounded-lg shadow p-8">
      <h1 class="text-2xl font-bold mb-2">${escapeHTMLToEntities(module.title)} — Quiz Results</h1>
      <p class="mb-4">You scored
        <span class="font-bold ${passed ? "text-green-600" : "text-red-600"}">${result.score}%</span>
        (${result.numCorrect}/${result.questions.length} correct).
        ${passed ? "Great job! 🎉" : "Review the explanations below, then continue or try again."}
      </p>
      <div id="lq-review" class="mb-6">${buildTysReview(reviewData)}</div>
      <div class="flex flex-wrap justify-between gap-4">
        <button id="lq-retry" class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded">Retake Quiz</button>
        <button id="lq-continue" class="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
          ${module.next ? "Complete & Continue" : "Finish Path"}
        </button>
      </div>
    </div><br/><br/><br/>`;

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

// Spaced repetition (Phase 4).
// A lightweight Leitner system: every answered question becomes a review item
// with a "box". Answer it right and it moves up a box (longer interval); miss
// it and it resets to box 0 (comes back soon). Items resurface in the Daily
// Review when they're due, so the mental model actually sticks.

// Days until the next review, by box.
const REVIEW_INTERVALS = [1, 2, 4, 7, 14, 30];

function reviewDayNum(date) {
  const d = date ? new Date(date) : new Date();
  return Math.floor(new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime() / 86400000);
}

function reviewKey(category, question) {
  return `${category}::${question}`;
}

// Upsert a review item after a question is answered.
async function recordReview(details, category, wasCorrect) {
  if (!details || !details.question) return;
  try {
    const key = reviewKey(category || "unknown", details.question);
    const existing = await DB.reviews.get(key);
    let box = existing ? existing.box : 0;
    box = wasCorrect ? Math.min(box + 1, REVIEW_INTERVALS.length - 1) : 0;
    const today = reviewDayNum();
    await DB.reviews.put({
      key,
      category: category || "unknown",
      details,
      box,
      due: today + REVIEW_INTERVALS[box],
      correct: (existing ? existing.correct || 0 : 0) + (wasCorrect ? 1 : 0),
      wrong: (existing ? existing.wrong || 0 : 0) + (wasCorrect ? 0 : 1),
      lastSeen: today,
    });
  } catch (e) {
    /* ignore */
  }
}

// Record a whole quiz's worth of answers at once.
async function recordReviewBatch(questions, selectedOptions) {
  if (!Array.isArray(questions)) return;
  for (let i = 0; i < questions.length; i++) {
    const q = questions[i];
    const det = q && q.details;
    if (!det) continue;
    const cat = q.category || det.category || "unknown";
    await recordReview(det, cat, selectedOptions[i] === det.answer);
  }
}

// Review items that are due today or earlier, soonest first.
async function getDueReviews(limit = 50) {
  try {
    const today = reviewDayNum();
    const all = await DB.reviews.toArray();
    return all
      .filter((r) => r.due <= today)
      .sort((a, b) => a.due - b.due)
      .slice(0, limit);
  } catch (e) {
    return [];
  }
}

async function dueReviewCount() {
  return (await getDueReviews(999)).length;
}

// The Daily Review session — presents due items one at a time.
async function ReviewPage(htmlEl) {
  const due = await getDueReviews(50);

  if (!due.length) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="text-4xl mb-2">🧠</div>
          <h1 class="text-2xl font-extrabold mb-2">All caught up!</h1>
          <p class="text-slate-600 mb-5">Nothing's due for review right now. Answer more quizzes and the concepts will resurface here on a smart schedule.</p>
          <a href="#today" class="gd-btn gd-btn-primary">Back to Today</a>
        </div>
      </div>`;
    return;
  }

  let i = 0;
  let remembered = 0;

  const renderSummary = () => {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="text-4xl mb-2">🎉</div>
          <h1 class="text-2xl font-extrabold mb-1">Review done!</h1>
          <p class="text-slate-600 mb-5">You remembered <b>${remembered} of ${due.length}</b>. The ones you missed will come back sooner.</p>
          <a href="#today" class="gd-btn gd-btn-primary">Back to Today</a>
        </div>
      </div>`;
  };

  const renderItem = () => {
    if (i >= due.length) return renderSummary();
    const item = due[i];
    const q = item.details;
    const opts = shuffle((q.options || []).slice())
      .map(
        (o) =>
          `<button type="button" class="review-option gd-option w-full text-left" data-correct="${
            o === q.answer
          }">${escapeHTMLToEntities(o)}</button>`
      )
      .join("");

    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card">
          <div class="flex items-center justify-between mb-3">
            <span class="gd-chip gd-chip-brand">${icon("book", "w-3.5 h-3.5")} Daily review</span>
            <span class="text-sm font-bold text-slate-500">${i + 1} / ${due.length}</span>
          </div>
          <p class="font-extrabold text-lg mb-3">${escapeHTMLToEntities(q.question)}</p>
          <div class="space-y-2">${opts}</div>
          <div class="review-feedback hidden mt-3 text-sm font-bold"></div>
          <button class="review-next gd-btn gd-btn-primary gd-btn-block mt-4 hidden">${
            i === due.length - 1 ? "Finish" : "Next"
          }</button>
        </div>
      </div>`;

    const optionEls = htmlEl.querySelectorAll(".review-option");
    const fb = htmlEl.querySelector(".review-feedback");
    const nextBtn = htmlEl.querySelector(".review-next");
    let answered = false;

    optionEls.forEach((btn) => {
      btn.addEventListener("click", async () => {
        if (answered) return;
        answered = true;
        const correct = btn.dataset.correct === "true";
        if (correct) remembered++;
        optionEls.forEach((o) => {
          o.disabled = true;
          if (o.dataset.correct === "true") o.classList.add("lesson-quiz-correct");
        });
        if (!correct) btn.classList.add("lesson-quiz-wrong");
        fb.classList.remove("hidden");
        fb.className = "review-feedback mt-3 text-sm font-bold " + (correct ? "text-grass-600" : "text-rose-500");
        fb.innerHTML =
          (correct ? "✅ Remembered!" : "❌ Worth another look.") +
          (q.explanation
            ? ` <span class="font-normal text-slate-600">${escapeHTMLToEntities(q.explanation)}</span>`
            : "");
        nextBtn.classList.remove("hidden");
        await recordReview(q, item.category, correct);
      });
    });

    nextBtn.addEventListener("click", () => {
      i++;
      renderItem();
    });
  };

  renderItem();
}

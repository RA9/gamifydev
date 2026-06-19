// Progress / gamification dashboard.
// Everything here is derived on the fly from the `scores` table (and `tests`
// for language labels), so there is no schema migration: XP, level, streak and
// badges are all computed from past attempts.

const XP_PER_CORRECT = 10;
const XP_PER_LEVEL = 100;

// Local YYYY-MM-DD key for a date, used for day-based streak counting.
function dayKey(date) {
  const d = new Date(date);
  return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;
}

// Number of consecutive days (ending today or yesterday) with at least one
// completed test. Returns 0 if the most recent attempt is older than yesterday.
function computeStreak(dates) {
  const days = new Set(dates.map(dayKey));
  if (days.size === 0) return 0;

  const today = new Date();
  const cursor = new Date(today.getFullYear(), today.getMonth(), today.getDate());

  // Allow the streak to be "current" if the last activity was today or yesterday.
  if (!days.has(dayKey(cursor))) {
    cursor.setDate(cursor.getDate() - 1);
    if (!days.has(dayKey(cursor))) return 0;
  }

  let streak = 0;
  while (days.has(dayKey(cursor))) {
    streak++;
    cursor.setDate(cursor.getDate() - 1);
  }
  return streak;
}

function computeProgress(scores, languageByTestId = {}) {
  const totalTests = scores.length;
  const totalCorrect = scores.reduce((a, s) => a + (s.numCorrect || 0), 0);
  const totalWrong = scores.reduce((a, s) => a + (s.numWrong || 0), 0);
  const totalQuestions = totalCorrect + totalWrong;
  const xp = totalCorrect * XP_PER_CORRECT;
  const level = Math.floor(xp / XP_PER_LEVEL) + 1;
  const xpIntoLevel = xp % XP_PER_LEVEL;
  const accuracy = totalQuestions
    ? Math.round((totalCorrect / totalQuestions) * 100)
    : 0;
  const perfectCount = scores.filter((s) => s.score === 100).length;

  // Best score per language.
  const bestByLanguage = {};
  scores.forEach((s) => {
    const lang =
      languageByTestId[s.test_id] ||
      s.details?.questions?.[0]?.category ||
      "unknown";
    bestByLanguage[lang] = Math.max(bestByLanguage[lang] ?? 0, s.score || 0);
  });

  const streak = computeStreak(scores.map((s) => s.created_at));

  const badges = [
    { id: "first-steps", icon: "rocket", name: "First Steps", desc: "Complete your first test", earned: totalTests >= 1 },
    { id: "scholar", icon: "cap", name: "Scholar", desc: "Complete 10 tests", earned: totalTests >= 10 },
    { id: "centurion", icon: "medal", name: "Centurion", desc: "Answer 100 questions", earned: totalQuestions >= 100 },
    { id: "perfectionist", icon: "star", name: "Perfectionist", desc: "Score 100% on a test", earned: perfectCount >= 1 },
    { id: "polyglot", icon: "globe", name: "Polyglot", desc: "Test in 3 languages", earned: Object.keys(bestByLanguage).length >= 3 },
    { id: "on-fire", icon: "flame", name: "On Fire", desc: "Reach a 3-day streak", earned: streak >= 3 },
  ];

  return {
    totalTests,
    totalCorrect,
    totalQuestions,
    accuracy,
    perfectCount,
    xp,
    level,
    xpIntoLevel,
    streak,
    bestByLanguage,
    badges,
  };
}

function statCard(label, value, iconName) {
  return `
    <div class="gd-card-sm text-center">
      <div class="grid h-10 w-10 mx-auto mb-2 place-items-center rounded-xl bg-brand-50 text-brand-500">${icon(
        iconName,
        "w-5 h-5"
      )}</div>
      <p class="text-3xl font-extrabold text-slate-800">${value}</p>
      <p class="text-xs font-bold uppercase tracking-wide text-slate-500">${label}</p>
    </div>`;
}

async function ProgressPage(htmlEl) {
  const scores = await DB.scores.orderBy("created_at").reverse().toArray();

  if (scores.length === 0) {
    htmlEl.innerHTML = `
      <div class="max-w-xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="w-44 mx-auto mb-2 animate-float">${illustration("trophy")}</div>
          <h1 class="text-2xl font-extrabold mb-2">Your Progress</h1>
          <p class="text-slate-600 mb-6">No quizzes yet! Complete one to start earning <b>XP</b>, build a <b>streak</b>, and unlock <b>badges</b>.</p>
          <a href="#test" class="gd-btn gd-btn-grass">Take your first quiz</a>
        </div>
      </div>`;
    return;
  }

  const tests = await DB.tests.toArray();
  const languageByTestId = {};
  tests.forEach((t) => (languageByTestId[t.id] = t.language));

  const p = computeProgress(scores, languageByTestId);

  const badgesHTML = p.badges
    .map(
      (b) => `
      <div class="rounded-2xl p-4 text-center border-2 transition-transform hover:-translate-y-0.5 ${
        b.earned
          ? "border-brand-200 bg-brand-50 shadow-card"
          : "border-slate-200 bg-slate-50 opacity-60"
      }" title="${b.desc}">
        <div class="grid h-12 w-12 mx-auto mb-2 place-items-center rounded-2xl ${
          b.earned ? "bg-brand-100 text-brand-600" : "bg-slate-100 text-slate-400"
        }">${icon(b.icon, "w-6 h-6")}</div>
        <p class="text-sm font-extrabold text-slate-800">${b.name}</p>
        <p class="text-xs text-slate-500">${b.desc}</p>
        ${
          b.earned
            ? `<span class="gd-chip gd-chip-grass mt-2 text-[10px]">${icon("check", "w-3 h-3")} Unlocked</span>`
            : `<span class="gd-chip gd-chip-slate mt-2 text-[10px]">${icon("lock", "w-3 h-3")} Locked</span>`
        }
      </div>`
    )
    .join("");

  const languagesHTML = Object.entries(p.bestByLanguage)
    .sort((a, b) => b[1] - a[1])
    .map(
      ([lang, best]) => {
        const barColor =
          best >= 85
            ? "from-grass-400 to-grass-600"
            : best >= 70
            ? "from-amber-400 to-amber-500"
            : "from-rose-400 to-rose-500";
        return `
      <div class="py-2.5">
        <div class="flex items-center justify-between mb-1.5">
          <span class="font-extrabold text-slate-700">${lang.toUpperCase()}</span>
          <span class="font-extrabold text-slate-500 text-sm">${best}%</span>
        </div>
        <div class="gd-progress h-2.5">
          <div class="h-full rounded-full bg-gradient-to-r ${barColor}" style="width: ${best}%"></div>
        </div>
      </div>`;
      }
    )
    .join("");

  const recentHTML = scores
    .slice(0, 10)
    .map((s, i) => {
      const lang =
        languageByTestId[s.test_id] ||
        s.details?.questions?.[0]?.category ||
        "unknown";
      const date = new Date(s.created_at).toLocaleDateString();
      const scoreColor =
        s.score >= 85 ? "text-grass-600" : s.score >= 70 ? "text-amber-500" : "text-rose-500";
      return `
      <div class="border-2 border-slate-100 rounded-2xl mb-2">
        <div class="flex items-center justify-between p-3">
          <div class="flex items-center gap-3">
            <span class="grid h-9 w-9 place-items-center rounded-xl bg-slate-100 text-sm font-extrabold ${scoreColor}">${s.score}</span>
            <div>
              <span class="font-extrabold text-slate-800 block leading-tight">${lang.toUpperCase()}</span>
              <span class="text-xs text-slate-400">${date}</span>
            </div>
          </div>
          <button class="history-review gd-btn gd-btn-secondary !py-1.5 !px-4 !text-xs" data-index="${i}">Review</button>
        </div>
        <div class="history-review-container hidden px-3 pb-3" data-index="${i}"></div>
      </div>`;
    })
    .join("");

  const earnedBadges = p.badges.filter((b) => b.earned).length;

  htmlEl.innerHTML = `
    <div class="max-w-5xl mx-auto space-y-5 animate-fade-up">
      <!-- Level / XP / streak hero -->
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-6 sm:p-8">
        <div class="absolute -top-8 -right-6 h-36 w-36 rounded-full bg-white/10 blur-2xl"></div>
        <div class="relative flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-4">
            <div class="grid h-16 w-16 place-items-center rounded-2xl bg-white/15 text-3xl font-extrabold">${p.level}</div>
            <div>
              <h1 class="text-2xl font-extrabold text-white">Level ${p.level}</h1>
              <p class="flex items-center gap-1.5 text-white/80 font-bold">${icon("gem", "w-4 h-4 text-grass-300")} ${p.xp} XP total</p>
            </div>
          </div>
          <div class="flex gap-3">
            <div class="text-center rounded-2xl bg-white/15 px-4 py-2">
              <p class="flex items-center justify-center gap-1.5 text-2xl font-extrabold">${icon("flame", "w-6 h-6 text-orange-300")} ${p.streak}</p>
              <p class="text-xs font-bold uppercase tracking-wide text-white/80">Streak</p>
            </div>
            <div class="text-center rounded-2xl bg-white/15 px-4 py-2">
              <p class="flex items-center justify-center gap-1.5 text-2xl font-extrabold">${icon("medal", "w-6 h-6 text-amber-300")} ${earnedBadges}</p>
              <p class="text-xs font-bold uppercase tracking-wide text-white/80">Badges</p>
            </div>
          </div>
        </div>
        <div class="relative mt-6">
          <div class="flex justify-between text-sm font-bold text-white/90 mb-1.5">
            <span>Level ${p.level}</span>
            <span>${p.xpIntoLevel} / ${XP_PER_LEVEL} XP to Level ${p.level + 1}</span>
          </div>
          <div class="w-full bg-white/20 rounded-full h-3.5 overflow-hidden">
            <div class="h-full rounded-full bg-gradient-to-r from-grass-300 to-grass-500" style="width: ${p.xpIntoLevel}%"></div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4">
        ${statCard("Tests Taken", p.totalTests, "file")}
        ${statCard("Questions", p.totalQuestions, "help")}
        ${statCard("Accuracy", p.accuracy + "%", "target")}
        ${statCard("Perfect Scores", p.perfectCount, "star")}
      </div>

      <div class="gd-card">
        <h2 class="flex items-center gap-2 text-xl font-extrabold mb-4">${icon("trophy", "w-5 h-5 text-amber-500")} Badges</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3">${badgesHTML}</div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div class="gd-card">
          <h2 class="text-xl font-extrabold mb-3">Best Score by Language</h2>
          ${languagesHTML}
        </div>
        <div class="gd-card">
          <h2 class="text-xl font-extrabold mb-3">Recent Tests</h2>
          ${recentHTML}
        </div>
      </div>
    </div>`;

  // Expand/collapse per-attempt answer review, reusing buildTysReview from tys.js.
  htmlEl.querySelectorAll(".history-review").forEach((btn) => {
    btn.addEventListener("click", () => {
      const i = btn.dataset.index;
      const container = htmlEl.querySelector(
        `.history-review-container[data-index="${i}"]`
      );
      if (container.classList.contains("hidden")) {
        if (!container.dataset.rendered) {
          container.innerHTML =
            typeof buildTysReview === "function"
              ? buildTysReview(scores[i])
              : "<p class='text-sm text-gray-500'>Review unavailable.</p>";
          container.dataset.rendered = "true";
        }
        container.classList.remove("hidden");
        btn.textContent = "Hide";
      } else {
        container.classList.add("hidden");
        btn.textContent = "Review";
      }
    });
  });
}

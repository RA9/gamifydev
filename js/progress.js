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
    { id: "first-steps", icon: "🚀", name: "First Steps", desc: "Complete your first test", earned: totalTests >= 1 },
    { id: "scholar", icon: "🎓", name: "Scholar", desc: "Complete 10 tests", earned: totalTests >= 10 },
    { id: "centurion", icon: "💯", name: "Centurion", desc: "Answer 100 questions", earned: totalQuestions >= 100 },
    { id: "perfectionist", icon: "⭐", name: "Perfectionist", desc: "Score 100% on a test", earned: perfectCount >= 1 },
    { id: "polyglot", icon: "🌐", name: "Polyglot", desc: "Test in 3 languages", earned: Object.keys(bestByLanguage).length >= 3 },
    { id: "on-fire", icon: "🔥", name: "On Fire", desc: "Reach a 3-day streak", earned: streak >= 3 },
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

function statCard(label, value) {
  return `
    <div class="bg-gray-50 border border-gray-200 rounded-lg p-4 text-center">
      <p class="text-3xl font-bold text-gray-800">${value}</p>
      <p class="text-sm text-gray-500">${label}</p>
    </div>`;
}

async function ProgressPage(htmlEl) {
  const scores = await DB.scores.orderBy("created_at").reverse().toArray();

  if (scores.length === 0) {
    htmlEl.innerHTML = `
      <div class="max-w-3xl mx-auto bg-white rounded-lg shadow p-8 text-center">
        <h1 class="text-2xl font-bold mb-2">Your Progress</h1>
        <p class="text-gray-600 mb-6">You haven't taken any tests yet. Complete a quiz to start earning XP, streaks, and badges!</p>
        <a href="#test" class="inline-block bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded">Take a Test</a>
      </div><br/><br/><br/>`;
    return;
  }

  const tests = await DB.tests.toArray();
  const languageByTestId = {};
  tests.forEach((t) => (languageByTestId[t.id] = t.language));

  const p = computeProgress(scores, languageByTestId);

  const badgesHTML = p.badges
    .map(
      (b) => `
      <div class="rounded-lg p-4 text-center border ${
        b.earned
          ? "border-indigo-300 bg-indigo-50"
          : "border-gray-200 bg-gray-50 opacity-50"
      }" title="${b.desc}">
        <div class="text-3xl mb-1 ${b.earned ? "" : "grayscale"}">${b.icon}</div>
        <p class="text-sm font-semibold text-gray-800">${b.name}</p>
        <p class="text-xs text-gray-500">${b.desc}</p>
      </div>`
    )
    .join("");

  const languagesHTML = Object.entries(p.bestByLanguage)
    .sort((a, b) => b[1] - a[1])
    .map(
      ([lang, best]) => `
      <div class="flex items-center justify-between py-2 border-b border-gray-100">
        <span class="font-medium text-gray-700">${lang.toUpperCase()}</span>
        <span class="font-bold ${
          best >= 85 ? "text-green-600" : best >= 70 ? "text-yellow-600" : "text-red-600"
        }">${best}%</span>
      </div>`
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
      return `
      <div class="border border-gray-200 rounded-lg mb-2">
        <div class="flex items-center justify-between p-3">
          <div>
            <span class="font-semibold text-gray-800">${lang.toUpperCase()}</span>
            <span class="text-sm text-gray-500 ml-2">${date}</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="font-bold ${
              s.score >= 85 ? "text-green-600" : s.score >= 70 ? "text-yellow-600" : "text-red-600"
            }">${s.score}%</span>
            <button class="history-review text-sm bg-indigo-600 hover:bg-indigo-700 text-white py-1 px-3 rounded" data-index="${i}">Review</button>
          </div>
        </div>
        <div class="history-review-container hidden px-3 pb-3" data-index="${i}"></div>
      </div>`;
    })
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-5xl mx-auto space-y-4">
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold">Your Progress</h1>
            <p class="text-gray-500">Level ${p.level} · ${p.xp} XP</p>
          </div>
          <div class="text-right">
            <p class="text-2xl font-bold text-orange-500">🔥 ${p.streak}</p>
            <p class="text-sm text-gray-500">day streak</p>
          </div>
        </div>
        <div class="mt-4">
          <div class="flex justify-between text-sm text-gray-500 mb-1">
            <span>Level ${p.level}</span>
            <span>${p.xpIntoLevel} / ${XP_PER_LEVEL} XP to next level</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-3">
            <div class="bg-indigo-600 h-3 rounded-full" style="width: ${p.xpIntoLevel}%"></div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        ${statCard("Tests Taken", p.totalTests)}
        ${statCard("Questions Answered", p.totalQuestions)}
        ${statCard("Accuracy", p.accuracy + "%")}
        ${statCard("Perfect Scores", p.perfectCount)}
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-xl font-bold mb-4">Badges</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3">${badgesHTML}</div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-bold mb-2">Best Score by Language</h2>
          ${languagesHTML}
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-bold mb-2">Recent Tests</h2>
          ${recentHTML}
        </div>
      </div>
    </div><br/><br/><br/>`;

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

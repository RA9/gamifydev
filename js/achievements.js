// Cross-world achievements — a meta-progression layer that rewards the whole
// journey, not just quizzes. Every metric is derived on the fly from existing
// data (paths, projects, scores, terminal progress, SRS reviews), so there is
// no schema change. Earned achievements are remembered in the `meta` store only
// so we know which ones are *new* (for the unlock celebration).

// Gather every metric the achievements need, in one pass.
async function computeAchievementMetrics() {
  const m = {
    lessonsCompleted: 0,
    worldsTouched: 0,
    worldsCompleted: 0,
    projectsCompleted: 0,
    terminalPassed: 0,
    correctAnswers: 0,
    perfectQuizzes: 0,
    xp: 0,
    streak: 0,
    reviewsTracked: 0,
  };

  try {
    const paths = await DB.paths.toArray();
    const byPath = {};
    paths.forEach((p) => {
      const g = (byPath[p.path_name] = byPath[p.path_name] || { total: 0, done: 0 });
      g.total += 1;
      if (p.is_completed) {
        g.done += 1;
        m.lessonsCompleted += 1;
      }
    });
    Object.values(byPath).forEach((g) => {
      if (g.done > 0) m.worldsTouched += 1;
      if (g.total > 0 && g.done >= g.total) m.worldsCompleted += 1;
    });
  } catch (e) {
    /* ignore */
  }

  try {
    if (typeof loadProjects === "function") {
      const projects = await loadProjects();
      for (const p of projects) {
        const done = await getProjectDone(p.id);
        if (done.filter((i) => i < (p.steps || 1)).length >= (p.steps || 1)) {
          m.projectsCompleted += 1;
        }
      }
    }
  } catch (e) {
    /* ignore */
  }

  try {
    if (typeof getMeta === "function") {
      m.terminalPassed = (await getMeta("terminal-progress", 0)) || 0;
    }
  } catch (e) {
    /* ignore */
  }

  try {
    const scores = await DB.scores.toArray();
    m.correctAnswers = scores.reduce((a, s) => a + (s.numCorrect || 0), 0);
    m.perfectQuizzes = scores.filter((s) => s.score === 100).length;
    m.xp = m.correctAnswers * (typeof XP_PER_CORRECT !== "undefined" ? XP_PER_CORRECT : 10);
  } catch (e) {
    /* ignore */
  }

  try {
    if (typeof getStreakInfo === "function") {
      m.streak = (await getStreakInfo()).streak || 0;
    }
  } catch (e) {
    /* ignore */
  }

  try {
    m.reviewsTracked = await DB.reviews.count();
  } catch (e) {
    /* ignore */
  }

  return m;
}

// Each achievement reports how far along it is via have/need, which gives us
// both the earned flag and a motivating "2 / 3" progress label for free.
const ACHIEVEMENTS = [
  { id: "first-lesson", emoji: "🎓", name: "First Lesson", desc: "Finish your first lesson", need: 1, have: (m) => m.lessonsCompleted },
  { id: "quiz-time", emoji: "⚡", name: "Quiz Time", desc: "Answer 10 questions correctly", need: 10, have: (m) => m.correctAnswers },
  { id: "sharpshooter", emoji: "🎯", name: "Sharpshooter", desc: "Score 100% on a quiz", need: 1, have: (m) => m.perfectQuizzes },
  { id: "on-fire", emoji: "🔥", name: "On Fire", desc: "Reach a 3-day streak", need: 3, have: (m) => m.streak },
  { id: "unstoppable", emoji: "🌟", name: "Unstoppable", desc: "Reach a 7-day streak", need: 7, have: (m) => m.streak },
  { id: "explorer", emoji: "🧭", name: "Explorer", desc: "Make progress in 3 worlds", need: 3, have: (m) => m.worldsTouched },
  { id: "conqueror", emoji: "🏆", name: "World Conqueror", desc: "Complete an entire world", need: 1, have: (m) => m.worldsCompleted },
  { id: "builder", emoji: "🛠️", name: "Builder", desc: "Finish your first project", need: 1, have: (m) => m.projectsCompleted },
  { id: "shipyard", emoji: "🚢", name: "Shipyard", desc: "Finish 3 projects", need: 3, have: (m) => m.projectsCompleted },
  { id: "hacker", emoji: "💻", name: "Terminal Hacker", desc: "Clear 5 terminal missions", need: 5, have: (m) => m.terminalPassed },
  { id: "memory", emoji: "🧠", name: "Memory Master", desc: "Track 20 cards in spaced repetition", need: 20, have: (m) => m.reviewsTracked },
  { id: "veteran", emoji: "💎", name: "Veteran", desc: "Earn 1000 XP", need: 1000, have: (m) => m.xp },
];

async function getAchievements() {
  const metrics = await computeAchievementMetrics();
  const list = ACHIEVEMENTS.map((a) => {
    const have = Math.max(0, a.have(metrics) || 0);
    const earned = have >= a.need;
    return {
      id: a.id,
      emoji: a.emoji,
      name: a.name,
      desc: a.desc,
      need: a.need,
      have: Math.min(have, a.need),
      earned,
      pct: Math.min(100, Math.round((have / a.need) * 100)),
    };
  });
  return {
    metrics,
    list,
    earnedCount: list.filter((a) => a.earned).length,
    total: list.length,
  };
}

function achievementsGridHTML(ach) {
  return ach.list
    .map(
      (a) => `
      <div class="rounded-2xl p-4 text-center border-2 transition-transform hover:-translate-y-0.5 ${
        a.earned
          ? "border-brand-200 bg-brand-50 shadow-card"
          : "border-slate-200 bg-slate-50"
      }" title="${a.desc}">
        <div class="text-3xl mb-1 ${a.earned ? "" : "grayscale opacity-40"}">${a.emoji}</div>
        <p class="text-sm font-extrabold text-slate-800">${a.name}</p>
        <p class="text-xs text-slate-500 leading-snug mt-0.5">${a.desc}</p>
        ${
          a.earned
            ? `<span class="gd-chip gd-chip-grass mt-2 text-[10px]">${icon("check", "w-3 h-3")} Unlocked</span>`
            : `<div class="mt-2">
                 <div class="gd-progress h-1.5"><div class="gd-progress-fill" style="width:${a.pct}%"></div></div>
                 <p class="text-[10px] font-bold text-slate-400 mt-1">${a.have} / ${a.need}</p>
               </div>`
        }
      </div>`
    )
    .join("");
}

// Show a Pixel celebration the first time new achievements are earned. On the
// very first run we silently baseline the current set (so existing learners
// don't get a flood), then only genuinely new unlocks celebrate afterwards.
async function maybeCelebrateAchievements() {
  if (typeof getMeta !== "function") return;
  // Don't stack on top of another celebration/briefing overlay — retry later.
  if (document.querySelector(".gd-overlay")) return;
  let ach;
  try {
    ach = await getAchievements();
  } catch (e) {
    return;
  }
  const earnedIds = ach.list.filter((a) => a.earned).map((a) => a.id);
  const seen = await getMeta("ach-seen", null);

  if (seen === null) {
    await setMeta("ach-seen", earnedIds); // baseline existing progress, no popup
    return;
  }
  const fresh = earnedIds.filter((id) => !seen.includes(id));
  if (!fresh.length) return;

  await setMeta("ach-seen", earnedIds);
  const def = ach.list.find((a) => a.id === fresh[0]);
  const extra = fresh.length > 1 ? ` (+${fresh.length - 1} more!)` : "";
  if (typeof pixelCelebrate === "function" && def) {
    pixelCelebrate({
      title: `Achievement unlocked! ${def.emoji}`,
      message: `${def.name} — ${def.desc}.${extra} You're stacking up the wins. Keep going!`,
      mood: "celebrate",
      cta: "Nice!",
    });
  }
}

// Daily engagement loop — the "Daily Standup".
// Startup-adventure framing: the learner has joined a startup; each day starts
// with a standup where they hit a small XP goal and pick up their next ticket.
// Everything here is derived from existing data (scores / paths / users), so
// there is no schema change.

const DAILY_GOAL_XP = 50;

function gdIsToday(d) {
  const x = new Date(d);
  const n = new Date();
  return (
    x.getFullYear() === n.getFullYear() &&
    x.getMonth() === n.getMonth() &&
    x.getDate() === n.getDate()
  );
}

function gdPerCorrect() {
  return typeof XP_PER_CORRECT !== "undefined" ? XP_PER_CORRECT : 10;
}

// XP earned today (from quizzes and lesson quizzes recorded in `scores`).
function computeDailyXp(scores) {
  const today = scores.filter((s) => gdIsToday(s.created_at));
  return today.reduce((a, s) => a + (s.numCorrect || 0), 0) * gdPerCorrect();
}

// --- Streaks that matter: freezes, milestones, calendar ---------------------
// `dayKey(date)` (Y-M-D, 0-indexed month) comes from progress.js.

const MAX_FREEZES = 2;
const STREAK_MILESTONES = [3, 7, 14, 30, 60, 100];

async function getMeta(key, fallback) {
  try {
    const row = await DB.table("meta").get(key);
    return row ? row.value : fallback;
  } catch (e) {
    return fallback;
  }
}
async function setMeta(key, value) {
  try {
    await DB.table("meta").put({ key, value });
  } catch (e) {
    /* ignore */
  }
}
async function getFrozenDayKeys() {
  try {
    return (await DB.table("daily").toArray())
      .filter((r) => r.frozen)
      .map((r) => r.date);
  } catch (e) {
    return [];
  }
}

// Streak counted from a set of "covered" day keys (active or frozen).
function computeStreakFromKeys(keys) {
  const set = keys instanceof Set ? keys : new Set(keys);
  if (!set.size) return 0;
  const n = new Date();
  const cursor = new Date(n.getFullYear(), n.getMonth(), n.getDate());
  if (!set.has(dayKey(cursor))) {
    cursor.setDate(cursor.getDate() - 1);
    if (!set.has(dayKey(cursor))) return 0;
  }
  let streak = 0;
  while (set.has(dayKey(cursor))) {
    streak++;
    cursor.setDate(cursor.getDate() - 1);
  }
  return streak;
}

// The canonical, freeze-aware streak info used across the app.
async function getStreakInfo() {
  const scores = await DB.scores.toArray();
  const active = new Set(scores.map((s) => dayKey(s.created_at)));
  const frozen = new Set(await getFrozenDayKeys());
  const covered = new Set([...active, ...frozen]);
  return {
    streak: computeStreakFromKeys(covered),
    freezes: await getMeta("freezes", MAX_FREEZES),
    active,
    frozen,
    covered,
  };
}

// Run once on load: auto-spend a freeze to bridge a single missed day, and
// grant a freeze at each new 7-day milestone (capped at MAX_FREEZES).
async function maintainStreak() {
  const scores = await DB.scores.toArray();
  if (!scores.length) return;

  const active = new Set(scores.map((s) => dayKey(s.created_at)));
  const frozen = new Set(await getFrozenDayKeys());
  let freezes = await getMeta("freezes", MAX_FREEZES);

  const n = new Date();
  const today = new Date(n.getFullYear(), n.getMonth(), n.getDate());
  const yest = new Date(today);
  yest.setDate(today.getDate() - 1);
  const dby = new Date(today);
  dby.setDate(today.getDate() - 2);
  const covered = (d) => active.has(dayKey(d)) || frozen.has(dayKey(d));

  // Missed yesterday but active the day before -> spend a freeze to bridge it.
  if (!covered(today) && !covered(yest) && covered(dby) && freezes > 0) {
    await DB.table("daily").put({ date: dayKey(yest), frozen: true });
    frozen.add(dayKey(yest));
    freezes -= 1;
    await setMeta("freezes", freezes);
  }

  // Grant a freeze for each new 7-day milestone (capped).
  const streak = computeStreakFromKeys(new Set([...active, ...frozen]));
  const milestone = Math.floor(streak / 7);
  const lastMilestone = await getMeta("lastFreezeMilestone", 0);
  if (milestone > lastMilestone) {
    freezes = Math.min(MAX_FREEZES, freezes + (milestone - lastMilestone));
    await setMeta("freezes", freezes);
    await setMeta("lastFreezeMilestone", milestone);
  }
}

function nextMilestone(streak) {
  return STREAK_MILESTONES.find((m) => m > streak) || null;
}

// --- Sprint review: a weekly recap of the last 7 days ------------------------
const GD_WEEKDAY = ["S", "M", "T", "W", "T", "F", "S"];

function computeWeekReview(scores) {
  const perCorrect = gdPerCorrect();
  const now = new Date();
  const days = [];
  const idx = {};
  for (let i = 6; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth(), now.getDate() - i);
    idx[dayKey(d)] = days.length;
    days.push({ label: GD_WEEKDAY[d.getDay()], xp: 0 });
  }
  let total = 0;
  let quizzes = 0;
  scores.forEach((s) => {
    const k = dayKey(s.created_at);
    if (k in idx) {
      const x = (s.numCorrect || 0) * perCorrect;
      days[idx[k]].xp += x;
      total += x;
      quizzes++;
    }
  });
  return {
    days,
    total,
    quizzes,
    activeDays: days.filter((d) => d.xp > 0).length,
    maxXp: Math.max(1, ...days.map((d) => d.xp)),
  };
}

function weekReviewCard(wr) {
  const bars = wr.days
    .map((d, i) => {
      const h = Math.round((d.xp / wr.maxXp) * 100);
      const isToday = i === wr.days.length - 1;
      const filled = d.xp > 0;
      return `
      <div class="flex flex-col items-center gap-1 flex-1">
        <div class="w-full flex items-end h-20">
          <div class="w-full rounded-t-md ${filled ? "bg-brand-500" : "bg-slate-200"} ${
        isToday ? "ring-2 ring-brand-300" : ""
      }" style="height: ${Math.max(filled ? 8 : 3, h)}%" title="${d.xp} XP"></div>
        </div>
        <span class="text-[10px] font-extrabold text-slate-400">${d.label}</span>
      </div>`;
    })
    .join("");
  return `
    <div class="gd-card">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-extrabold">This week's sprint</h2>
        <span class="text-xs font-bold uppercase tracking-wide text-slate-500">${wr.activeDays}/7 active days</span>
      </div>
      <div class="flex items-end gap-2 mb-4">${bars}</div>
      <div class="flex gap-5 text-sm">
        <div><span class="font-extrabold text-brand-600">${wr.total}</span> <span class="text-slate-500">XP this week</span></div>
        <div><span class="font-extrabold text-brand-600">${wr.quizzes}</span> <span class="text-slate-500">quiz${wr.quizzes === 1 ? "" : "zes"}</span></div>
      </div>
    </div>`;
}

// A 14-day activity calendar: active (filled), frozen (snowflake), missed.
function streakCalendar(covered, frozen) {
  const n = new Date();
  const cells = [];
  for (let i = 13; i >= 0; i--) {
    const d = new Date(n.getFullYear(), n.getMonth(), n.getDate() - i);
    const k = dayKey(d);
    const isToday = i === 0;
    let cls = "bg-slate-100 text-slate-400";
    let inner = String(d.getDate());
    if (frozen.has(k)) {
      cls = "bg-brand-100 text-brand-600";
      inner = "❄";
    } else if (covered.has(k)) {
      cls = "bg-brand-500 text-white";
    }
    cells.push(
      `<div class="grid place-items-center h-7 w-7 rounded-lg text-[11px] font-extrabold ${cls} ${
        isToday ? "ring-2 ring-brand-400 ring-offset-1 ring-offset-white" : ""
      }">${inner}</div>`
    );
  }
  return `<div class="flex flex-wrap gap-1.5">${cells.join("")}</div>`;
}

// The streak card shown on the Daily Standup.
function streakCard(info) {
  const next = nextMilestone(info.streak);
  const freezeText =
    info.freezes > 0
      ? `<span class="inline-flex items-center gap-1 text-brand-600 font-bold">❄ ${info.freezes} streak freeze${info.freezes > 1 ? "s" : ""}</span>`
      : `<span class="text-slate-400 font-bold">No freezes left</span>`;
  return `
    <div class="gd-card">
      <div class="flex items-center justify-between gap-4 mb-4">
        <div class="flex items-center gap-3">
          <div class="text-3xl">🔥</div>
          <div>
            <p class="text-2xl font-extrabold leading-none">${info.streak} day${info.streak === 1 ? "" : "s"}</p>
            <p class="text-xs font-bold uppercase tracking-wide text-slate-500">Current streak</p>
          </div>
        </div>
        <div class="text-right text-sm">${freezeText}</div>
      </div>
      ${streakCalendar(info.covered, info.frozen)}
      <p class="text-sm text-slate-500 mt-3">${
        next
          ? `<b>${next - info.streak} day${next - info.streak === 1 ? "" : "s"}</b> to your ${next}-day milestone.`
          : "You're a streak legend. 🏆"
      } A freeze covers one missed day so a busy day won't break your run.</p>
    </div>`;
}

// --- Launch Week: the guided 7-day "first week on the job" -------------------
// A curated Frontend arc that ends each "ship" day with a real artifact, so a
// learner has built something tangible within their first week.

const LAUNCH_WEEK_FRONTEND = [
  { theme: "Orientation", subtitle: "Get your bearings in web dev", modules: ["History of the Web", "Intro to Programming"] },
  { theme: "Structure", subtitle: "Build the skeleton of a page", modules: ["HTML Basics"] },
  { theme: "Style", subtitle: "Make it beautiful", modules: ["CSS Basics"] },
  { theme: "Build day", subtitle: "Put structure and style together", modules: ["Building a Website with HTML and CSS"] },
  { theme: "Ship it 🚀", subtitle: "Build & ship your first real component", modules: ["Project: Build a Profile Card"] },
  { theme: "Interactivity", subtitle: "Make pages respond to people", modules: ["JavaScript Basics", "Building Interactive JavaScript Websites"] },
  { theme: "Ship again 🚀", subtitle: "Build a working mini-app", modules: ["Project: Build a Quiz Game"] },
];

// Launch Week currently guides the Frontend path (the onboarding default).
async function computeLaunchWeek() {
  const user = (await DB.users.toArray())[0];
  if (!user) return null;
  const pathName = await resolveLessonPath(user.preference);
  if (pathName !== "frontend") return null;

  const notes = await DB.paths.where("path_name").equals("frontend").toArray();
  const rec = {};
  notes.forEach((n) => (rec[n.title] = n));
  const currentTitle = (notes.find((n) => n.current) || {}).title;

  const days = LAUNCH_WEEK_FRONTEND.map((d, i) => {
    const modules = d.modules.map((t) => ({
      title: t,
      done: !!(rec[t] && rec[t].is_completed),
      current: t === currentTitle,
    }));
    return { n: i + 1, theme: d.theme, subtitle: d.subtitle, modules, complete: modules.every((m) => m.done) };
  });

  const firstIncomplete = days.findIndex((d) => !d.complete);
  return {
    days,
    completedDays: days.filter((d) => d.complete).length,
    total: days.length,
    currentDay: firstIncomplete === -1 ? days.length : firstIncomplete + 1,
    allDone: firstIncomplete === -1,
  };
}

// A banner for the Daily Standup linking into Launch Week.
function launchWeekBanner(lw) {
  if (!lw || lw.allDone) return "";
  const day = lw.days[lw.currentDay - 1];
  return `
    <a href="#launch" class="block rounded-3xl bg-gradient-to-r from-brand-600 to-brand-500 text-white shadow-soft p-5 hover:brightness-110 transition">
      <div class="flex items-center justify-between gap-4">
        <div>
          <span class="gd-chip bg-white/15 text-white mb-1">🚀 Launch Week</span>
          <p class="font-extrabold text-lg">Day ${lw.currentDay}: ${day.theme}</p>
          <p class="text-white/85 text-sm">${day.subtitle}</p>
        </div>
        <div class="text-right shrink-0">
          <p class="text-2xl font-extrabold">${lw.completedDays}/${lw.total}</p>
          <p class="text-xs text-white/80 font-bold uppercase tracking-wide">days</p>
        </div>
      </div>
    </a>`;
}

async function LaunchWeekPage(htmlEl) {
  const lw = await computeLaunchWeek();
  if (!lw) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <div class="text-4xl mb-2">🚀</div>
          <h1 class="text-2xl font-extrabold mb-2">Launch Week</h1>
          <p class="text-slate-600 mb-5">Launch Week is the guided first-week plan for the <b>Frontend</b> path. Switch to Frontend to follow it, or keep going on your current journey.</p>
          <a href="#journey" class="gd-btn gd-btn-primary">Go to your journey</a>
        </div>
      </div>`;
    return;
  }

  const pct = Math.round((lw.completedDays / lw.total) * 100);
  const daysHtml = lw.days
    .map((d) => {
      const isCurrent = d.n === lw.currentDay && !lw.allDone;
      if (d.complete) {
        return `
        <div class="gd-card-sm flex items-center gap-3">
          <div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-grass-500 text-white">${icon("check", "w-5 h-5")}</div>
          <div><p class="font-extrabold">Day ${d.n}: ${d.theme}</p><p class="text-xs text-slate-500">Complete</p></div>
        </div>`;
      }
      if (isCurrent) {
        const tasks = d.modules
          .map((m) => {
            const node = m.done
              ? icon("check", "w-4 h-4 text-grass-600")
              : m.current
              ? icon("play", "w-4 h-4 text-brand-600")
              : icon("lock", "w-4 h-4 text-slate-300");
            return `<li class="flex items-center gap-2 py-1"><span class="shrink-0">${node}</span><span class="${m.done ? "line-through text-slate-400" : ""}">${gdEsc(m.title)}</span></li>`;
          })
          .join("");
        return `
        <div class="gd-card ring-2 ring-brand-300">
          <span class="gd-chip gd-chip-brand mb-2">Today's focus · Day ${d.n}</span>
          <h2 class="text-xl font-extrabold">${d.theme}</h2>
          <p class="text-slate-500 text-sm mb-3">${d.subtitle}</p>
          <ul class="mb-4 text-sm font-bold text-slate-700">${tasks}</ul>
          <button id="lw-continue" class="gd-btn gd-btn-primary">Continue Day ${d.n} ${icon("arrowRight", "w-4 h-4")}</button>
        </div>`;
      }
      return `
        <div class="gd-card-sm flex items-center gap-3 opacity-70">
          <div class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-slate-100 text-slate-400">${icon("lock", "w-4 h-4")}</div>
          <div><p class="font-extrabold text-slate-500">Day ${d.n}: ${d.theme}</p><p class="text-xs text-slate-400">${d.subtitle}</p></div>
        </div>`;
    })
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-4">
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-6 sm:p-8">
        <div class="absolute -top-8 -right-6 h-36 w-36 rounded-full bg-white/10 blur-2xl"></div>
        <div class="relative">
          <span class="gd-chip bg-white/15 text-white mb-2">🚀 Launch Week</span>
          <h1 class="text-2xl sm:text-3xl font-extrabold">${lw.allDone ? "Launch Week complete! 🎉" : "Your first week as a dev"}</h1>
          <p class="text-white/85 mt-1">${lw.allDone ? "You built and shipped real projects in 7 days. This is just the start." : "Seven focused days. By the end you'll have built and shipped real things."}</p>
          <div class="mt-4">
            <div class="flex justify-between text-sm font-bold text-white/90 mb-1"><span>${lw.completedDays} of ${lw.total} days</span><span>${pct}%</span></div>
            <div class="w-full bg-white/20 rounded-full h-3"><div class="h-full rounded-full bg-white" style="width: ${pct}%"></div></div>
          </div>
        </div>
      </div>
      <div class="space-y-3">${daysHtml}</div>
    </div>`;

  const cont = document.querySelector("#lw-continue");
  if (cont) {
    cont.addEventListener("click", () => {
      if (typeof handleNotePage === "function") handleNotePage();
    });
  }
}

// Playful startup job titles that level up with the learner.
function roleForLevel(level) {
  if (level >= 8) return "Staff Engineer";
  if (level >= 6) return "Senior Dev";
  if (level >= 4) return "Developer";
  if (level >= 2) return "Junior Dev";
  return "Intern";
}

function gdEsc(s) {
  return (s || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

// The learner's next actionable: continue their current lesson/project, or, if
// the path is finished, nudge a quiz.
async function getNextTicket(user) {
  const pathName = await resolveLessonPath(user.preference);
  const current = await DB.paths
    .where("path_name")
    .equals(pathName)
    .and((p) => p.current)
    .last();

  if (current) {
    const isProject = /^project/i.test(current.title);
    return {
      kind: isProject ? "project" : "lesson",
      title: current.title,
      cta: isProject ? "Start building" : "Continue lesson",
    };
  }
  return { kind: "quiz", title: "Test your skills", cta: "Take a quiz", hash: "test" };
}

// A circular daily-goal progress ring.
function dailyGoalRing(pct) {
  const r = 52;
  const c = 2 * Math.PI * r;
  const off = c * (1 - Math.min(100, pct) / 100);
  return `<svg viewBox="0 0 120 120" class="w-28 h-28 -rotate-90" aria-hidden="true">
    <circle cx="60" cy="60" r="${r}" fill="none" stroke-width="12" class="text-slate-200" stroke="currentColor" />
    <circle cx="60" cy="60" r="${r}" fill="none" stroke-width="12" stroke-linecap="round"
      class="text-brand-500" stroke="currentColor"
      stroke-dasharray="${c.toFixed(1)}" stroke-dashoffset="${off.toFixed(1)}"
      style="transition: stroke-dashoffset .8s ease" />
  </svg>`;
}

// --- Daily Challenge: one bite-sized puzzle per day --------------------------
// A deterministic daily pick from a curated pool. Completing it records a score
// (test_id "challenge-<day>"), so it feeds XP, streak and the sprint review and
// can only be done once per day.
async function getDailyChallenge() {
  let pool = [];
  try {
    pool = await (await fetch("./data/challenges.json")).json();
  } catch (e) {
    return null;
  }
  if (!pool.length) return null;
  const n = new Date();
  const dayNum = Math.floor(
    new Date(n.getFullYear(), n.getMonth(), n.getDate()).getTime() / 86400000
  );
  const challenge = pool[dayNum % pool.length];
  const testId = "challenge-" + dayKey(n);
  let done = null;
  try {
    done = await DB.scores.where("test_id").equals(testId).first();
  } catch (e) {
    /* ignore */
  }
  return { challenge, testId, done: !!done, doneCorrect: done ? (done.numCorrect || 0) > 0 : false };
}

function challengeCard(dc) {
  if (!dc) return "";
  if (dc.done) {
    return `
      <div class="gd-card">
        <span class="gd-chip gd-chip-brand mb-2">${icon("zap", "w-3.5 h-3.5")} Daily challenge</span>
        <p class="font-extrabold">${dc.doneCorrect ? "Nailed it! 🎉" : "Done for today"}</p>
        <p class="text-slate-500 text-sm mt-1">Come back tomorrow for a fresh challenge.</p>
      </div>`;
  }
  const c = dc.challenge;
  const code = c.code
    ? `<div class="gd-codeblock"><div class="gd-codeblock-head">${gdEsc(
        (c.lang || "code").toUpperCase()
      )}</div><pre><code>${gdEsc(c.code)}</code></pre></div>`
    : "";
  const opts = c.options
    .map(
      (o) =>
        `<button type="button" class="challenge-option gd-option w-full text-left" data-correct="${!!o.correct}">${gdEsc(
          o.text
        )}</button>`
    )
    .join("");
  return `
    <div class="gd-card" id="daily-challenge">
      <span class="gd-chip gd-chip-brand mb-2">${icon("zap", "w-3.5 h-3.5")} Daily challenge · +20 XP</span>
      <p class="font-extrabold mb-3">${gdEsc(c.prompt)}</p>
      ${code}
      <div class="space-y-2 mt-3">${opts}</div>
      <div class="challenge-feedback hidden mt-3 text-sm font-bold"></div>
    </div>`;
}

// The Daily Standup landing for returning learners. First-time visitors fall
// back to the marketing home.
async function TodayPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) return HomePage(htmlEl);

  const scores = await DB.scores.toArray();
  const totalCorrect = scores.reduce((a, s) => a + (s.numCorrect || 0), 0);
  const xp = totalCorrect * gdPerCorrect();
  const level = Math.floor(xp / 100) + 1;
  const streakInfo = await getStreakInfo();
  const streak = streakInfo.streak;
  const dailyXp = computeDailyXp(scores);
  const pct = Math.round((dailyXp / DAILY_GOAL_XP) * 100);
  const metGoal = dailyXp >= DAILY_GOAL_XP;
  const ticket = await getNextTicket(user);
  const launchWeek = await computeLaunchWeek();
  const weekReview = computeWeekReview(scores);
  const dailyChallenge = await getDailyChallenge();
  const dueReviews = typeof getDueReviews === "function" ? await getDueReviews(999) : [];

  // The mentor's contextual greeting. "New" = no quiz history and nothing
  // completed yet.
  let completedCount = 0;
  try {
    completedCount = (await DB.paths.toArray()).filter((p) => p.is_completed).length;
  } catch (e) {
    /* ignore */
  }
  const greeting =
    typeof mentorLine === "function"
      ? mentorLine({
          name: user.name,
          isNew: scores.length === 0 && completedCount === 0,
          streak,
          metGoal,
          dailyXp,
          goal: DAILY_GOAL_XP,
        })
      : `Welcome back, ${user.name}!`;

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto space-y-5 animate-fade-up">
      <!-- Standup header: your mentor greets you -->
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-6 sm:p-8">
        <div class="absolute -top-8 -right-6 h-36 w-36 rounded-full bg-white/10 blur-2xl"></div>
        <div class="relative flex items-start gap-4">
          <div class="shrink-0 grid place-items-center h-14 w-14 rounded-2xl bg-white/95 shadow-soft animate-pop-in">${
            typeof mascotSvg === "function" ? mascotSvg("w-11 h-11") : ""
          }</div>
          <div class="min-w-0">
            <span class="gd-chip bg-white/15 text-white mb-2">${
              typeof MASCOT_NAME !== "undefined" ? MASCOT_NAME : "Pixel"
            } · your guide</span>
            <p class="text-xl sm:text-2xl font-extrabold leading-snug">${gdEsc(greeting)}</p>
            <p class="text-white/70 text-sm font-bold mt-2">${roleForLevel(level)} · Level ${level} · 🔥 ${streak}-day streak</p>
          </div>
        </div>
      </div>

      <!-- Daily goal ring -->
      <div class="gd-card flex items-center gap-6">
        <div class="relative shrink-0 grid place-items-center">
          ${dailyGoalRing(pct)}
          <div class="absolute inset-0 grid place-items-center text-center">
            <div>
              <p class="text-xl font-extrabold leading-none">${dailyXp}</p>
              <p class="text-[11px] font-bold text-slate-500">/ ${DAILY_GOAL_XP} XP</p>
            </div>
          </div>
        </div>
        <div class="flex-1">
          <h2 class="text-lg font-extrabold">${metGoal ? "Goal smashed! 🎉" : "Today's goal"}</h2>
          <p class="text-slate-600 text-sm mt-1">${
            metGoal
              ? "You hit today's goal and protected your streak. Push for more, or rest up for tomorrow."
              : `Earn <b>${DAILY_GOAL_XP - dailyXp} more XP</b> today to hit your goal and keep your streak alive.`
          }</p>
        </div>
      </div>

      <!-- Launch Week -->
      ${launchWeekBanner(launchWeek)}

      <!-- Daily review (spaced repetition) -->
      ${
        dueReviews.length
          ? `<a href="#review" class="gd-card flex items-center justify-between gap-4 hover:-translate-y-0.5 transition-transform">
              <div class="flex items-center gap-3">
                <div class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-brand-100 text-brand-600">${icon("book", "w-5 h-5")}</div>
                <div>
                  <p class="font-extrabold leading-tight">Daily review</p>
                  <p class="text-sm text-slate-500">${dueReviews.length} concept${dueReviews.length === 1 ? "" : "s"} due — keep them fresh</p>
                </div>
              </div>
              <span class="gd-btn gd-btn-primary !py-2 !px-4 !text-sm shrink-0">Review</span>
            </a>`
          : ""
      }

      <!-- Daily challenge -->
      ${challengeCard(dailyChallenge)}

      <!-- Streak -->
      ${streakCard(streakInfo)}

      <!-- Next ticket -->
      <div class="gd-card">
        <span class="gd-chip gd-chip-brand mb-3">${icon("file", "w-3.5 h-3.5")} Your next ticket</span>
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 class="text-xl font-extrabold">${gdEsc(ticket.title)}</h2>
            <p class="text-slate-500 text-sm mt-1">${
              ticket.kind === "project"
                ? "A hands-on build — time to ship something."
                : ticket.kind === "quiz"
                ? "Your path is clear — sharpen your skills with a quiz."
                : "Pick up right where you left off."
            }</p>
          </div>
          ${
            ticket.hash
              ? `<a href="#${ticket.hash}" class="gd-btn gd-btn-primary">${ticket.cta} ${icon("arrowRight", "w-4 h-4")}</a>`
              : `<button id="today-continue" class="gd-btn gd-btn-primary">${ticket.cta} ${icon("arrowRight", "w-4 h-4")}</button>`
          }
        </div>
      </div>

      <!-- Quick links -->
      <div class="grid grid-cols-2 gap-4">
        <a href="#progress" class="gd-card-sm flex items-center gap-3 hover:-translate-y-0.5 transition-transform">
          <div class="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-brand-100 text-brand-600">${icon("trophy", "w-5 h-5")}</div>
          <div><p class="font-extrabold leading-tight">Progress</p><p class="text-xs text-slate-500">${xp} XP · Level ${level}</p></div>
        </a>
        <a href="#test" class="gd-card-sm flex items-center gap-3 hover:-translate-y-0.5 transition-transform">
          <div class="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-brand-100 text-brand-600">${icon("zap", "w-5 h-5")}</div>
          <div><p class="font-extrabold leading-tight">Quiz</p><p class="text-xs text-slate-500">Test yourself</p></div>
        </a>
      </div>

      <!-- Sprint review -->
      ${weekReviewCard(weekReview)}
    </div>`;

  const cont = document.querySelector("#today-continue");
  if (cont) {
    cont.addEventListener("click", () => {
      if (typeof handleNotePage === "function") handleNotePage();
    });
  }

  // Wire the daily challenge: answering records a score (feeding XP/streak),
  // then refreshes the standup so the goal ring and "done" state update.
  const challengeEl = document.querySelector("#daily-challenge");
  if (challengeEl && dailyChallenge && !dailyChallenge.done) {
    const opts = challengeEl.querySelectorAll(".challenge-option");
    opts.forEach((btn) => {
      btn.addEventListener("click", async () => {
        if ([...opts].some((o) => o.disabled)) return;
        const c = dailyChallenge.challenge;
        const correct = btn.dataset.correct === "true";
        const correctText = (c.options.find((o) => o.correct) || {}).text;
        opts.forEach((o) => {
          o.disabled = true;
          if (o.dataset.correct === "true") o.classList.add("lesson-quiz-correct");
        });
        if (!correct) btn.classList.add("lesson-quiz-wrong");
        const fb = challengeEl.querySelector(".challenge-feedback");
        fb.classList.remove("hidden");
        fb.className =
          "challenge-feedback mt-3 text-sm font-bold " + (correct ? "text-grass-600" : "text-rose-500");
        fb.innerHTML =
          (correct ? "✅ Correct! +20 XP " : "❌ Not quite. ") +
          `<span class="font-normal text-slate-600">${gdEsc(c.explanation || "")}</span>`;
        try {
          await createStorage("scores", {
            id: randomID(),
            test_id: dailyChallenge.testId,
            score: correct ? 100 : 0,
            numCorrect: correct ? 2 : 0,
            numWrong: correct ? 0 : 1,
            details: {
              questions: [
                {
                  details: {
                    question: c.prompt,
                    options: c.options.map((o) => o.text),
                    answer: correctText,
                    explanation: c.explanation,
                  },
                },
              ],
              selectedOptions: [btn.textContent],
            },
            created_at: new Date(),
          });
        } catch (e) {
          /* ignore */
        }
        setTimeout(() => TodayPage(document.querySelector("main")), 1500);
      });
    });
  }
}

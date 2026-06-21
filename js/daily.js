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

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto space-y-5 animate-fade-up">
      <!-- Standup header -->
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-500 text-white shadow-soft p-6 sm:p-8">
        <div class="absolute -top-8 -right-6 h-36 w-36 rounded-full bg-white/10 blur-2xl"></div>
        <div class="relative">
          <span class="gd-chip bg-white/15 text-white mb-3">${icon("zap", "w-3.5 h-3.5")} Daily standup</span>
          <h1 class="text-2xl sm:text-3xl font-extrabold">Welcome back, ${gdEsc(user.name)}!</h1>
          <p class="text-white/85 font-bold mt-1">${roleForLevel(level)} · Level ${level} · 🔥 ${streak}-day streak</p>
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
    </div>`;

  const cont = document.querySelector("#today-continue");
  if (cont) {
    cont.addEventListener("click", () => {
      if (typeof handleNotePage === "function") handleNotePage();
    });
  }
}

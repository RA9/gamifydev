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
  const streak =
    typeof computeStreak === "function" ? computeStreak(scores.map((s) => s.created_at)) : 0;
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

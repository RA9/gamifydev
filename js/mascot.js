// Pixel — the learner's guide/mentor at the "startup".
// A friendly app-buddy mascot plus a contextual voice that reacts to the
// learner's state. Kept self-contained so other moments (lesson done, level up)
// can reuse it later.

const MASCOT_NAME = "Pixel";

// Each path's startup "mission" — Pixel briefs this the first time a learner
// opens that path's Journey board.
const PATH_STORY = {
  frontend: {
    crew: "Frontend crew",
    mission:
      "You're on the Frontend crew — the team that builds everything users see and touch. Your mission: go from zero to shipping real, interactive web pages. I'll be right here the whole way.",
  },
  backend: {
    crew: "Backend crew",
    mission:
      "Welcome to the Backend crew — the engine room. Your mission: store data, build APIs, and power apps from behind the scenes. Let's make the machine run.",
  },
  fullstack: {
    crew: "Fullstack crew",
    mission:
      "Welcome to the Fullstack crew — you'll work across the whole stack. Your mission: connect front and back into complete, working apps you can actually ship.",
  },
  c: {
    crew: "Systems crew",
    mission:
      "Welcome to the Systems crew — we work close to the metal. C is the language behind operating systems, databases, and game engines. Your mission: master the fundamentals that everything else is built on.",
  },
  java: {
    crew: "Java crew",
    mission:
      "Welcome to the Java crew — builders of robust, portable software that runs everywhere. Your mission: learn the object-oriented thinking that powers banks, Android apps, and huge enterprise systems.",
  },
  linux: {
    crew: "Ops crew",
    mission:
      "Welcome to the Ops crew — the people who run the machines. Linux and the command line are where real developers live. Your mission: get fluent at the shell, move around the filesystem blindfolded, and write your first scripts. Tip: try things live in the Terminal Trainer as you go!",
  },
};

function pathMission(pathName) {
  return PATH_STORY[pathName] || PATH_STORY.frontend;
}

// The mascot avatar as inline SVG, sized via the class string.
function mascotSvg(cls = "w-12 h-12", mood = "happy") {
  const smile =
    mood === "celebrate"
      ? '<path d="M17 27 q7 8 14 0 z" fill="#3d2880"/>' // open grin
      : '<path d="M18 28 q6 5 12 0" stroke="#3d2880" stroke-width="2.6" fill="none" stroke-linecap="round"/>';
  const eyes =
    mood === "wink"
      ? '<circle cx="19" cy="23" r="3" fill="#3d2880"/><path d="M26 23 q3 -2 6 0" stroke="#3d2880" stroke-width="2.6" fill="none" stroke-linecap="round"/>'
      : '<circle cx="19" cy="23" r="3" fill="#3d2880"/><circle cx="29" cy="23" r="3" fill="#3d2880"/>';
  return `<svg viewBox="0 0 48 48" class="${cls}" role="img" aria-label="${MASCOT_NAME}">
    <defs><linearGradient id="mascotg" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#7857f7"/><stop offset="1" stop-color="#5731c8"/>
    </linearGradient></defs>
    <rect x="4" y="6" width="40" height="38" rx="13" fill="url(#mascotg)"/>
    <circle cx="24" cy="5" r="2.5" fill="#b4aaff"/>
    <rect x="3" y="4" width="2.5" height="4" rx="1" fill="#b4aaff" transform="rotate(0 24 4)"/>
    <rect x="11" y="14" width="26" height="22" rx="9" fill="#ffffff"/>
    ${eyes}
    <circle cx="14" cy="28" r="1.8" fill="#c4b5fd"/>
    <circle cx="34" cy="28" r="1.8" fill="#c4b5fd"/>
    ${smile}
  </svg>`;
}

function gdPick(arr) {
  return arr[Math.floor(Math.random() * arr.length)] || arr[0];
}

// A celebratory overlay where Pixel pops up — used at "ship it"/level-up
// moments. Returns a promise that resolves when dismissed, so callers can
// await it before continuing.
function pixelCelebrate({ title, message, mood = "celebrate", cta = "Keep going" }) {
  return new Promise((resolve) => {
    const overlay = document.createElement("div");
    overlay.className =
      "gd-overlay fixed inset-0 z-50 grid place-items-center bg-slate-900/60 backdrop-blur-sm p-4";
    overlay.innerHTML = `
      <div class="gd-card max-w-sm w-full text-center animate-pop-in">
        <div class="w-20 h-20 mx-auto mb-3 animate-float">${
          typeof mascotSvg === "function" ? mascotSvg("w-20 h-20", mood) : ""
        }</div>
        <h2 class="text-2xl font-extrabold mb-1">${title}</h2>
        <p class="text-slate-600 mb-5">${message}</p>
        <button class="gd-btn gd-btn-primary gd-btn-block">${cta}</button>
      </div>`;
    const close = () => {
      overlay.remove();
      resolve();
    };
    overlay.querySelector("button").addEventListener("click", close);
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) close();
    });
    document.body.appendChild(overlay);
    overlay.querySelector("button").focus();
  });
}

// A contextual one-liner from the mentor, based on the learner's state.
// ctx: { name, isNew, streak, metGoal, dailyXp, goal }
function mentorLine(ctx) {
  const name = ctx.name || "there";
  if (ctx.isNew) {
    return `Welcome to the team, ${name}! I'm ${MASCOT_NAME}, your guide. Let's ship your first thing today.`;
  }
  if (ctx.metGoal) {
    return gdPick([
      "Daily goal smashed — you're on fire today. 🔥",
      "Nice work hitting today's goal. Future-you says thanks.",
      "Goal done! Want to push for a bit more, or call it a win?",
    ]);
  }
  if (ctx.streak >= 7) {
    return `${ctx.streak} days straight — that's serious momentum. Let's keep it rolling.`;
  }
  if (ctx.streak >= 3) {
    return `${ctx.streak}-day streak! You're building a real habit now.`;
  }
  if (!ctx.dailyXp) {
    return gdPick([
      `Morning, ${name}! A quick lesson and the day's already moving.`,
      "Ready to ship something today? Let's grab your next ticket.",
      "Five focused minutes is all it takes. Shall we?",
    ]);
  }
  const left = Math.max(0, (ctx.goal || 0) - ctx.dailyXp);
  return gdPick([
    `Good momentum, ${name} — ${left} XP to today's goal.`,
    "Nice progress. You're closing in on today's goal.",
  ]);
}

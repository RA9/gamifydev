// Pixel — the learner's guide/mentor at the "startup".
// A friendly app-buddy mascot plus a contextual voice that reacts to the
// learner's state. Kept self-contained so other moments (lesson done, level up)
// can reuse it later.

const MASCOT_NAME = "Pixel";

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

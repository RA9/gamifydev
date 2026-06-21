// Central icon + illustration system.
//
// The whole UI is built from JS template strings, so instead of pulling in an
// icon font or external sprite (bad for an offline-first PWA) we keep a small
// set of inline SVGs here:
//   icon(name, cls)      -> a Lucide-style line icon, inheriting `currentColor`
//   langBadge(key, cls)  -> a language monogram chip in the language's brand color
//   illustration(name)   -> a larger decorative SVG scene for heroes/empty states
//
// Line icons share one stroke style for a consistent, "proper icon" look.

// Inner markup for each line icon (viewBox 0 0 24 24, stroke = currentColor).
const GD_ICONS = {
  menu: '<path d="M4 6h16M4 12h16M4 18h16"/>',
  x: '<path d="M18 6 6 18M6 6l12 12"/>',
  sun: '<circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/>',
  moon: '<path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9z"/>',
  flame: '<path d="M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z"/>',
  gem: '<path d="M6 3h12l4 6-10 13L2 9z"/><path d="M11 3 8 9l4 13 4-13-3-6M2 9h20"/>',
  zap: '<path d="M13 2 3 14h9l-1 8 10-12h-9l1-8z"/>',
  trophy: '<path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/><path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/><path d="M18 2H6v7a6 6 0 0 0 12 0V2z"/>',
  medal: '<path d="M7.21 15 2.66 7.14a2 2 0 0 1 .13-2.2L4.4 2.8A2 2 0 0 1 6 2h12a2 2 0 0 1 1.6.8l1.6 2.14a2 2 0 0 1 .14 2.2L16.79 15"/><path d="M11 12 5.12 2.2M13 12l5.88-9.8M8 7h8"/><circle cx="12" cy="17" r="5"/><path d="M12 18v-2h-.5"/>',
  book: '<path d="M12 7v14"/><path d="M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z"/>',
  map: '<path d="M14.106 5.553a2 2 0 0 0 1.788 0l3.659-1.83A1 1 0 0 1 21 4.619v12.764a1 1 0 0 1-.553.894l-4.553 2.277a2 2 0 0 1-1.788 0l-4.212-2.106a2 2 0 0 0-1.788 0l-3.659 1.83A1 1 0 0 1 3 19.381V6.618a1 1 0 0 1 .553-.894l4.553-2.277a2 2 0 0 1 1.788 0zM15 5.764v15M9 3.236v15"/>',
  target: '<circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/>',
  rocket: '<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="M12 15l-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/>',
  cap: '<path d="M21.42 10.922a1 1 0 0 0-.019-1.838L12.83 5.18a2 2 0 0 0-1.66 0L2.6 9.08a1 1 0 0 0 0 1.832l8.57 3.908a2 2 0 0 0 1.66 0z"/><path d="M22 10v6"/><path d="M6 12.5V16a6 3 0 0 0 12 0v-3.5"/>',
  puzzle: '<path d="M15.39 4.39a1 1 0 0 0 1.68-.474 2.5 2.5 0 1 1 3.014 3.015 1 1 0 0 0-.474 1.68l1.683 1.682a2.414 2.414 0 0 1 0 3.414L19.61 19.39a1 1 0 0 1-1.68-.474 2.5 2.5 0 1 0-3.014 3.015 1 1 0 0 1 .474 1.68l-1.683 1.682a2.414 2.414 0 0 1-3.414 0l-.005-.005a1 1 0 0 1 .474-1.68 2.5 2.5 0 1 0-3.014-3.015 1 1 0 0 1-1.68.474l-1.683-1.682a2.414 2.414 0 0 1 0-3.414l1.683-1.682a1 1 0 0 0 .474-1.68 2.5 2.5 0 1 1 3.014-3.015 1 1 0 0 0 1.68-.474l1.683-1.682a2.414 2.414 0 0 1 3.414 0z"/>',
  code: '<path d="m16 18 6-6-6-6M8 6l-6 6 6 6"/>',
  chat: '<path d="M7.9 20A9 9 0 1 0 4 16.1L2 22z"/>',
  sparkles: '<path d="M9.937 15.5A2 2 0 0 0 8.5 14.063l-6.135-1.582a.5.5 0 0 1 0-.962L8.5 9.936A2 2 0 0 0 9.937 8.5l1.582-6.135a.5.5 0 0 1 .963 0L14.063 8.5A2 2 0 0 0 15.5 9.937l6.135 1.581a.5.5 0 0 1 0 .964L15.5 14.063a2 2 0 0 0-1.437 1.437l-1.582 6.135a.5.5 0 0 1-.963 0z"/><path d="M20 3v4M22 5h-4M4 17v2M5 18H3"/>',
  check: '<path d="M20 6 9 17l-5-5"/>',
  checkCircle: '<path d="M21.801 10A10 10 0 1 1 17 3.335"/><path d="m9 11 3 3L22 4"/>',
  xCircle: '<circle cx="12" cy="12" r="10"/><path d="m15 9-6 6M9 9l6 6"/>',
  lock: '<rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>',
  play: '<path d="M6 3 20 12 6 21z"/>',
  arrowRight: '<path d="M5 12h14M12 5l7 7-7 7"/>',
  arrowLeft: '<path d="M19 12H5M12 19l-7-7 7-7"/>',
  chevrons: '<path d="m6 17 5-5-5-5M13 17l5-5-5-5"/>',
  mail: '<rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>',
  pin: '<path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0z"/><circle cx="12" cy="10" r="3"/>',
  clock: '<circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>',
  globe: '<circle cx="12" cy="12" r="10"/><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20M2 12h20"/>',
  link: '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>',
  file: '<path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7z"/><path d="M14 2v4a2 2 0 0 0 2 2h4M16 13H8M16 17H8M10 9H8"/>',
  help: '<circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><path d="M12 17h.01"/>',
  star: '<path d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.123 2.123 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.123 2.123 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.122 2.122 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.122 2.122 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.122 2.122 0 0 0 1.597-1.16z"/>',
  refresh: '<path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M3 21v-5h5"/>',
  database: '<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.657 4.03 3 9 3s9-1.343 9-3V5"/><path d="M3 12c0 1.657 4.03 3 9 3s9-1.343 9-3"/>',
  gamepad: '<path d="M6 11h4M8 9v4M15 12h.01M18 10h.01"/><rect width="20" height="12" x="2" y="6" rx="6"/>',
  twitter: '<path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"/>',
  send: '<path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"/><path d="m21.854 2.147-10.94 10.939"/>',
  bulb: '<path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"/><path d="M9 18h6M10 22h4"/>',
  alert: '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3z"/><path d="M12 9v4M12 17h.01"/>',
  key: '<circle cx="7.5" cy="15.5" r="5.5"/><path d="m21 2-9.6 9.6M15.5 7.5l3 3L22 7l-3-3"/>',
  flask: '<path d="M14 2v6a2 2 0 0 0 .245.96l5.51 10.08A2 2 0 0 1 18 22H6a2 2 0 0 1-1.755-2.96l5.51-10.08A2 2 0 0 0 10 8V2"/><path d="M6.453 15h11.094M8.5 2h7"/>',
};

// Icons whose look reads better solid-filled than outlined.
const GD_FILLED = new Set(["zap", "play"]);

function icon(name, cls = "w-5 h-5") {
  const inner = GD_ICONS[name];
  if (!inner) return "";
  const filled = GD_FILLED.has(name);
  const paint = filled
    ? 'fill="currentColor"'
    : 'fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"';
  return `<svg class="${cls}" viewBox="0 0 24 24" ${paint} aria-hidden="true" focusable="false">${inner}</svg>`;
}

// ---------------------------------------------------------------------------
// Standard brand palette (mirrors tailwind.config.js tokens) + language colors.
// Languages get a monogram badge in their official brand color — consistent and
// reliable across all of them, without fragile per-logo SVG paths.
// ---------------------------------------------------------------------------
const GD_BRAND = {
  primary: "#6740e8", // brand-600
  success: "#46a302", // grass-600
  warning: "#f59e0b", // amber-500
  danger: "#f43f5e", // rose-500
};

const GD_LANG = {
  html: { abbr: "&lt;/&gt;", fg: "#E34F26" },
  css: { abbr: "CSS", fg: "#1572B6" },
  javascript: { abbr: "JS", fg: "#A16207", bg: "#F7DF1E" },
  c: { abbr: "C", fg: "#0F5FAB" },
  python: { abbr: "Py", fg: "#3776AB" },
  java: { abbr: "Jv", fg: "#E76F00" },
  sql: { abbr: "SQL", fg: "#0E7490" },
};

// A square monogram chip for a programming language. Uses the brand colour for
// every language (the monogram itself distinguishes them) so the palette stays
// on-brand instead of a rainbow of per-language colours.
function langBadge(key, cls = "h-12 w-12 text-base") {
  const m = GD_LANG[key] || { abbr: (key || "?").slice(0, 2).toUpperCase() };
  const brand = "#6740e8"; // brand-600
  return `<span class="grid place-items-center rounded-2xl font-extrabold leading-none ${cls}"
    style="background:${brand}1f;color:${brand}">${m.abbr}</span>`;
}

// ---------------------------------------------------------------------------
// Illustrations — larger decorative SVG scenes built from the brand palette.
// ---------------------------------------------------------------------------
function illustration(name, cls = "w-full h-auto") {
  const scenes = {
    // Friendly "learning app" scene: a code window with a progress ring,
    // floating XP/star accents and soft blobs.
    hero: `
      <svg class="${cls}" viewBox="0 0 420 320" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <circle cx="210" cy="160" r="150" fill="#ffffff" opacity="0.10"/>
        <circle cx="330" cy="80" r="46" fill="#b4aaff" opacity="0.55"/>
        <circle cx="80" cy="250" r="30" fill="#b4aaff" opacity="0.45"/>
        <g filter="url(#s)">
          <rect x="78" y="86" width="244" height="168" rx="20" fill="#ffffff"/>
        </g>
        <rect x="78" y="86" width="244" height="40" rx="20" fill="#f1f0ff"/>
        <rect x="78" y="106" width="244" height="20" fill="#f1f0ff"/>
        <circle cx="100" cy="106" r="5" fill="#937dff"/>
        <circle cx="118" cy="106" r="5" fill="#b4aaff"/>
        <circle cx="136" cy="106" r="5" fill="#937dff"/>
        <rect x="100" y="146" width="70" height="10" rx="5" fill="#c7cbff"/>
        <rect x="100" y="166" width="120" height="10" rx="5" fill="#e2e8f0"/>
        <rect x="100" y="186" width="96" height="10" rx="5" fill="#e2e8f0"/>
        <rect x="100" y="212" width="150" height="12" rx="6" fill="#e7e5ff"/>
        <rect x="100" y="212" width="96" height="12" rx="6" fill="#937dff"/>
        <g transform="translate(250 150)">
          <circle r="40" fill="#f2f1ff"/>
          <path d="M0 -40 A40 40 0 1 1 -34 20" stroke="#6740e8" stroke-width="10" stroke-linecap="round" fill="none"/>
          <path d="M-14 0 l9 9 l19 -19" stroke="#6740e8" stroke-width="8" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
        </g>
        <path d="M348 150 l7 14 16 2 -12 11 3 16 -14 -8 -14 8 3 -16 -12 -11 16 -2z" fill="#b4aaff"/>
        <path d="M52 110 l5 10 11 1 -8 8 2 11 -10 -5 -10 5 2 -11 -8 -8 11 -1z" fill="#7857f7" opacity="0.85"/>
        <defs>
          <filter id="s" x="58" y="76" width="284" height="208" filterUnits="userSpaceOnUse">
            <feDropShadow dx="0" dy="12" stdDeviation="14" flood-color="#281c5a" flood-opacity="0.18"/>
          </filter>
        </defs>
      </svg>`,
    // Empty-state trophy with sparkles, for the progress dashboard.
    trophy: `
      <svg class="${cls}" viewBox="0 0 240 200" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <circle cx="120" cy="100" r="84" fill="#f1f0ff"/>
        <path d="M84 56h72v22a36 36 0 0 1-72 0z" fill="#b4aaff"/>
        <path d="M84 56H66a18 18 0 0 0 18 18M156 56h18a18 18 0 0 1-18 18" stroke="#6740e8" stroke-width="8" fill="none" stroke-linecap="round"/>
        <rect x="112" y="112" width="16" height="20" fill="#6740e8"/>
        <rect x="96" y="132" width="48" height="14" rx="6" fill="#6740e8"/>
        <rect x="88" y="146" width="64" height="12" rx="6" fill="#7857f7"/>
        <path d="M120 64l4 8 9 1-6 6 1 9-8-4-8 4 1-9-6-6 9-1z" fill="#ffffff"/>
        <path d="M188 60l3 6 7 1-5 5 1 7-6-3-6 3 1-7-5-5 7-1z" fill="#b4aaff"/>
        <path d="M44 96l3 6 7 1-5 5 1 7-6-3-6 3 1-7-5-5 7-1z" fill="#b4aaff"/>
      </svg>`,
  };
  return scenes[name] || "";
}

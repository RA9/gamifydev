// Build-Along Projects — the "now go build something real" track.
//
// Projects are guided, step-by-step builds. Metadata lives in
// data/projects.json; the step content is a markdown file per project in
// data/notes/projects/<id>.md, structured as an intro followed by
// "## Step N — Title" sections (one focused step each). The player parses that
// file into steps and lets the learner tick each one off, with progress saved
// in the `meta` store (via getMeta/setMeta from daily.js).

let GD_PROJECTS = null; // cached metadata
const GD_PROJECT_CACHE = {}; // id -> parsed { title, introHtml, steps }

async function loadProjects() {
  if (GD_PROJECTS) return GD_PROJECTS;
  try {
    const res = await fetch("./data/projects.json");
    GD_PROJECTS = (await res.json()).projects || [];
  } catch (e) {
    GD_PROJECTS = [];
  }
  return GD_PROJECTS;
}

function getProjectMeta(id) {
  return (GD_PROJECTS || []).find((p) => p.id === id) || null;
}

// Parse a project markdown file into an intro + ordered steps.
function parseProjectMd(md) {
  const lines = md.split("\n");
  let title = "Project";
  let start = 0;
  if (lines[0] && lines[0].startsWith("# ")) {
    title = lines[0].slice(2).trim();
    start = 1;
  }
  const stepHeading = /^##\s+Step\s+\d+\s*[—–-]\s*(.*)$/i;
  // Find where the first step begins.
  let firstStep = lines.length;
  for (let i = start; i < lines.length; i++) {
    if (stepHeading.test(lines[i])) {
      firstStep = i;
      break;
    }
  }
  const introMd = lines.slice(start, firstStep).join("\n").trim();
  const introHtml =
    introMd && typeof markdownToHtml === "function" ? markdownToHtml(introMd) : "";

  const steps = [];
  let cur = null;
  for (let i = firstStep; i < lines.length; i++) {
    const m = lines[i].match(stepHeading);
    if (m) {
      if (cur) steps.push(cur);
      cur = { title: (m[1] || "").trim() || `Step ${steps.length + 1}`, body: [] };
    } else if (cur) {
      cur.body.push(lines[i]);
    }
  }
  if (cur) steps.push(cur);

  steps.forEach((s) => {
    const body = s.body.join("\n").trim();
    s.bodyHtml =
      body && typeof markdownToHtml === "function" ? markdownToHtml(body) : "";
    delete s.body;
  });

  return { title, introHtml, steps };
}

async function getParsedProject(id) {
  if (GD_PROJECT_CACHE[id]) return GD_PROJECT_CACHE[id];
  try {
    const res = await fetch(`./data/notes/projects/${id}.md`);
    if (!res.ok) return null;
    const md = await res.text();
    const parsed = parseProjectMd(md);
    GD_PROJECT_CACHE[id] = parsed;
    return parsed;
  } catch (e) {
    return null;
  }
}

// --- Progress (persisted in meta) -------------------------------------------
async function getProjectDone(id) {
  if (typeof getMeta !== "function") return [];
  const arr = await getMeta("proj:" + id, []);
  return Array.isArray(arr) ? arr : [];
}
async function setProjectDone(id, arr) {
  if (typeof setMeta === "function") await setMeta("proj:" + id, arr);
}

// --- Gallery ----------------------------------------------------------------
async function ProjectsPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  const projects = await loadProjects();

  const worldName = (w) =>
    typeof WORLDS !== "undefined" && WORLDS.find((x) => x.path === w)
      ? WORLDS.find((x) => x.path === w).name
      : w.charAt(0).toUpperCase() + w.slice(1);

  const cards = await Promise.all(
    projects.map(async (p) => {
      const done = await getProjectDone(p.id);
      const total = p.steps || 1;
      const completed = done.filter((i) => i < total).length;
      const pct = Math.round((completed / total) * 100);
      const isDone = completed >= total;
      const label = completed === 0 ? "Start building" : isDone ? "Revisit" : "Continue";
      const tech = (p.tech || [])
        .map((t) => `<span class="gd-chip !py-0.5 !px-2 !text-[10px]">${t}</span>`)
        .join("");
      return `
      <button type="button" onclick="openProject('${p.id}')"
        class="gd-card group text-left flex flex-col gap-3 transition-transform hover:-translate-y-1">
        <div class="flex items-start justify-between gap-2">
          <div class="grid h-14 w-14 place-items-center rounded-2xl bg-brand-100 text-3xl">${p.emoji || "🛠️"}</div>
          ${
            isDone
              ? `<span class="gd-chip gd-chip-grass !text-[10px] !py-0.5">${icon("check", "w-3 h-3")} Shipped</span>`
              : `<span class="gd-chip !text-[10px] !py-0.5">${p.difficulty || ""}</span>`
          }
        </div>
        <div>
          <h3 class="text-lg font-extrabold">${p.title}</h3>
          <p class="text-sm text-slate-500">${p.tagline || ""}</p>
        </div>
        <div class="flex flex-wrap gap-1.5">${tech}<span class="gd-chip gd-chip-brand !py-0.5 !px-2 !text-[10px]">${worldName(p.world)}</span></div>
        <div class="gd-progress h-2 mt-auto"><div class="gd-progress-fill" style="width:${pct}%"></div></div>
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold text-slate-500">${completed}/${total} steps</span>
          <span class="gd-btn gd-btn-primary !py-1.5 !px-4 !text-xs">${label}</span>
        </div>
      </button>`;
    })
  );

  htmlEl.innerHTML = `
    <div class="max-w-5xl mx-auto animate-fade-up">
      <div class="text-center mb-7">
        <span class="gd-chip gd-chip-brand mb-2">${icon("rocket", "w-3.5 h-3.5")} Build something real</span>
        <h1 class="text-2xl sm:text-3xl font-extrabold">Projects</h1>
        <p class="text-slate-500 mt-1 max-w-md mx-auto">Theory sticks when you build. Each project is a guided, step-by-step build — tick off steps as you go and ship something you can show off.</p>
      </div>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">${cards.join("")}</div>
    </div>`;
}

function openProject(id) {
  window.location.hash = "project/" + id;
}

// --- Build-along player ------------------------------------------------------
async function ProjectPlayerPage(htmlEl, id) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  await loadProjects();
  const meta = getProjectMeta(id);
  const parsed = await getParsedProject(id);

  if (!meta || !parsed || !parsed.steps.length) {
    htmlEl.innerHTML = `
      <div class="max-w-2xl mx-auto animate-fade-up">
        <div class="gd-card text-center">
          <p class="text-slate-600 mb-5">This project isn't available yet.</p>
          <a href="#projects" class="gd-btn gd-btn-primary">Back to Projects</a>
        </div>
      </div>`;
    return;
  }

  const steps = parsed.steps;
  const done = await getProjectDone(id);
  const doneSet = new Set(done.filter((i) => i < steps.length));
  const completed = doneSet.size;
  const pct = Math.round((completed / steps.length) * 100);
  const allDone = completed >= steps.length;

  const stepCards = steps
    .map((s, i) => {
      const isDone = doneSet.has(i);
      const num = isDone
        ? `<span class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-grass-500 text-white">${icon("check", "w-5 h-5")}</span>`
        : `<span class="grid h-9 w-9 shrink-0 place-items-center rounded-full bg-brand-100 text-brand-700 font-extrabold">${i + 1}</span>`;
      return `
      <div class="gd-card ${isDone ? "opacity-75" : ""}" id="pstep-${i}">
        <div class="flex items-center gap-3 mb-3">
          ${num}
          <h2 class="text-lg font-extrabold flex-1">${s.title}</h2>
        </div>
        <div class="gd-prose">${s.bodyHtml}</div>
        <div class="mt-4 flex justify-end">
          <button type="button" onclick="toggleProjectStep('${id}', ${i})"
            class="gd-btn ${isDone ? "gd-btn-secondary" : "gd-btn-primary"} !py-2 !px-4 !text-sm">
            ${isDone ? "✓ Completed — undo" : "Mark step complete"}
          </button>
        </div>
      </div>`;
    })
    .join("");

  const techTags = (meta.tech || [])
    .map((t) => `<span class="gd-chip !py-0.5 !px-2 !text-[10px]">${t}</span>`)
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-5">
      <a href="#projects" class="inline-flex items-center gap-1 text-sm font-bold text-brand-600 hover:text-brand-700">${icon("arrowLeft", "w-4 h-4")} All projects</a>

      <div class="gd-card">
        <div class="flex items-start gap-4">
          <div class="grid h-16 w-16 shrink-0 place-items-center rounded-2xl bg-brand-100 text-4xl">${meta.emoji || "🛠️"}</div>
          <div class="flex-1">
            <h1 class="text-2xl font-extrabold">${meta.title}</h1>
            <p class="text-slate-500 text-sm mt-0.5">${meta.outcome || meta.tagline || ""}</p>
            <div class="flex flex-wrap gap-1.5 mt-2">${techTags}<span class="gd-chip !py-0.5 !px-2 !text-[10px]">${meta.difficulty || ""}</span></div>
          </div>
        </div>
        ${parsed.introHtml ? `<div class="gd-prose mt-4 pt-4 border-t border-slate-200/70">${parsed.introHtml}</div>` : ""}
        <div class="mt-4">
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-xs font-bold uppercase tracking-wide text-slate-500">Build progress</span>
            <span class="text-sm font-extrabold text-brand-600">${completed}/${steps.length} steps</span>
          </div>
          <div class="gd-progress h-3"><div class="gd-progress-fill" style="width:${pct}%"></div></div>
        </div>
      </div>

      ${
        allDone
          ? `<div class="gd-card text-center bg-grass-50 border border-grass-200">
               <p class="text-3xl mb-1">🚀</p>
               <h2 class="text-xl font-extrabold text-grass-700">You shipped ${meta.title}!</h2>
               <p class="text-slate-600 mt-1">Every step done. Go run it, tweak it, and make it yours — then pick your next build.</p>
               <a href="#projects" class="gd-btn gd-btn-primary mt-4 inline-block">Pick another project</a>
             </div>`
          : ""
      }

      <div class="space-y-4">${stepCards}</div>
    </div>`;
}

// Toggle a step's completion, persist it, and re-render. Fires a one-time
// celebration the moment the final step lands.
async function toggleProjectStep(id, idx) {
  await loadProjects();
  const meta = getProjectMeta(id);
  const parsed = await getParsedProject(id);
  if (!parsed) return;
  const total = parsed.steps.length;

  const done = await getProjectDone(id);
  const set = new Set(done.filter((i) => i < total));
  const wasAllDone = set.size >= total;
  if (set.has(idx)) set.delete(idx);
  else set.add(idx);
  const arr = [...set].sort((a, b) => a - b);
  await setProjectDone(id, arr);

  const nowAllDone = arr.length >= total;

  await ProjectPlayerPage(document.querySelector("main"), id);

  if (nowAllDone && !wasAllDone && typeof pixelCelebrate === "function") {
    pixelCelebrate({
      title: "Shipped it! 🚀",
      message: `You built ${meta ? meta.title : "your project"} step by step. That's a real thing you made — be proud of it.`,
      mood: "celebrate",
      cta: "Heck yeah",
    });
  } else {
    // Keep the just-toggled step in view after re-render.
    const el = document.getElementById("pstep-" + idx);
    if (el) el.scrollIntoView({ block: "center", behavior: "smooth" });
  }
}

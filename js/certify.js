// Certification — the spine of the program. A certificate is earned by the full
// gauntlet: finish the track's lessons, solve coding sessions, ship the capstone
// project, and pass a timed final assessment (logical MCQs + live code
// challenges) at or above the pass threshold. Everything is derived from data
// already tracked locally; only the assessment score and earned-date are stored.

let GD_CERTS = null;
const GD_ASSESS = { current: null, editors: {}, timer: null, endsAt: 0 };

async function loadCertifications() {
  if (GD_CERTS) return GD_CERTS;
  try {
    const res = await fetch("./data/certifications.json");
    GD_CERTS = (await res.json()).certifications || [];
  } catch (e) {
    GD_CERTS = [];
  }
  return GD_CERTS;
}

function getCert(id) {
  return (GD_CERTS || []).find((c) => c.id === id) || null;
}

// Compute live status of every requirement for a certificate.
async function computeCertStatus(cert) {
  const req = cert.requirements || {};

  // Lessons in the track's world.
  let lessonsDone = 0,
    lessonsTotal = 0;
  try {
    const rows = await DB.paths.where("path_name").equals(cert.world).toArray();
    lessonsTotal = rows.length;
    lessonsDone = rows.filter((r) => r.is_completed).length;
  } catch (e) {}
  const lessonsOk = !req.lessons || (lessonsTotal > 0 && lessonsDone >= lessonsTotal);

  // Coding sessions solved (any).
  let codingDone = 0;
  const needCoding = req.codingSessions || 0;
  try {
    if (typeof loadChallenges === "function") {
      const all = await loadChallenges();
      for (const c of all) {
        if (await isChallengeDone(c.id)) codingDone += 1;
      }
    }
  } catch (e) {}
  const codingOk = codingDone >= needCoding;

  // Capstone project complete.
  let capstoneOk = !req.capstone;
  let capstoneTitle = "";
  if (req.capstone) {
    try {
      await loadProjects();
      const meta = getProjectMeta(req.capstone);
      capstoneTitle = meta ? meta.title : req.capstone;
      const done = await getProjectDone(req.capstone);
      capstoneOk = meta && done.filter((i) => i < (meta.steps || 1)).length >= (meta.steps || 1);
    } catch (e) {}
  }

  // Assessment best score.
  let assessScore = 0;
  try {
    assessScore = (typeof getMeta === "function" && (await getMeta("cert-assess:" + cert.id, 0))) || 0;
  } catch (e) {}
  const assessOk = !req.assessment || assessScore >= cert.passThreshold;

  // Prereqs (everything except the exam itself) gate the assessment.
  const eligibleForAssessment = lessonsOk && codingOk && capstoneOk;
  const earned = lessonsOk && codingOk && capstoneOk && assessOk;

  return {
    lessons: { done: lessonsDone, total: lessonsTotal, ok: lessonsOk },
    coding: { done: codingDone, need: needCoding, ok: codingOk },
    capstone: { ok: capstoneOk, title: capstoneTitle, id: req.capstone },
    assessment: { score: assessScore, ok: assessOk, threshold: cert.passThreshold },
    eligibleForAssessment,
    earned,
  };
}

// --- Certification hub ------------------------------------------------------
async function CertifyPage(htmlEl) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  const certs = await loadCertifications();

  const cards = await Promise.all(
    certs.map(async (cert) => {
      const s = await computeCertStatus(cert);
      const row = (label, ok, detail) => `
        <li class="flex items-center gap-3 py-1.5">
          <span class="${ok ? "text-grass-600" : "text-slate-300"}">${icon(ok ? "checkCircle" : "lock", "w-5 h-5")}</span>
          <span class="flex-1 text-sm ${ok ? "font-bold text-slate-700" : "text-slate-500"}">${label}</span>
          <span class="text-xs font-bold ${ok ? "text-grass-600" : "text-slate-400"}">${detail}</span>
        </li>`;

      let cta;
      if (s.earned) {
        cta = `<a href="#certificate/${cert.id}" class="gd-btn gd-btn-primary gd-btn-block">${icon("medal", "w-4 h-4")} View your certificate</a>`;
      } else if (s.eligibleForAssessment) {
        cta = `<a href="#assess/${cert.id}" class="gd-btn gd-btn-primary gd-btn-block">${s.assessment.score ? "Retake" : "Take"} the final assessment</a>`;
      } else {
        cta = `<button class="gd-btn gd-btn-secondary gd-btn-block opacity-70" disabled>Finish the steps above to unlock the exam</button>`;
      }

      return `
      <div class="gd-card max-w-2xl mx-auto">
        <div class="flex items-start gap-4 mb-4">
          <div class="grid h-16 w-16 shrink-0 place-items-center rounded-2xl bg-brand-100 text-4xl">${cert.emoji || "🎓"}</div>
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <h2 class="text-xl font-extrabold">${cert.title}</h2>
              ${s.earned ? `<span class="gd-chip gd-chip-grass !text-[10px]">${icon("check", "w-3 h-3")} Earned</span>` : ""}
            </div>
            <p class="text-sm text-slate-500 mt-0.5">${cert.blurb || ""}</p>
          </div>
        </div>
        <ul class="border-t border-b border-slate-200/70 divide-y divide-slate-100 mb-4">
          ${row(`Complete the ${cert.world.charAt(0).toUpperCase() + cert.world.slice(1)} lessons`, s.lessons.ok, `${s.lessons.done}/${s.lessons.total}`)}
          ${row(`Solve coding sessions`, s.coding.ok, `${s.coding.done}/${s.coding.need}`)}
          ${row(`Ship the capstone${s.capstone.title ? `: ${s.capstone.title}` : ""}`, s.capstone.ok, s.capstone.ok ? "done" : "to do")}
          ${row(`Pass the final assessment (≥ ${s.assessment.threshold}%)`, s.assessment.ok, s.assessment.score ? `${s.assessment.score}%` : "—")}
        </ul>
        ${cta}
      </div>`;
    })
  );

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up">
      <div class="text-center mb-6">
        <span class="gd-chip gd-chip-brand mb-2">${icon("medal", "w-3.5 h-3.5")} Earn your credential</span>
        <h1 class="text-2xl sm:text-3xl font-extrabold">Certification</h1>
        <p class="text-slate-500 mt-1 max-w-md mx-auto">No shortcuts. Finish the lessons, write real code, ship the capstone, and pass a graded exam — then the certificate is yours.</p>
      </div>
      <div class="space-y-5">${cards.join("")}</div>
    </div>`;
}

// --- Final assessment -------------------------------------------------------
async function AssessmentPage(htmlEl, id) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  await loadCertifications();
  const cert = getCert(id);
  if (!cert) {
    htmlEl.innerHTML = `<div class="max-w-2xl mx-auto gd-card text-center"><p class="mb-4">Unknown certification.</p><a href="#certify" class="gd-btn gd-btn-primary">Back</a></div>`;
    return;
  }
  const status = await computeCertStatus(cert);
  if (!status.eligibleForAssessment) {
    window.location.hash = "certify";
    return;
  }

  const a = cert.assessment || {};
  // Build the MCQ set from the logical-challenge bank.
  const all = await DB.questions.toArray();
  const pool = all.filter(
    (q) => q.category === a.mcqCategory && q.details && q.details.difficulty === "challenge"
  );
  const mcqs = shuffle(pool).slice(0, a.mcqCount || 8).map((q) => ({
    question: q.details.question,
    options: shuffle(q.details.options.slice()),
    answer: q.details.answer,
  }));

  // Build the code challenge set.
  await loadChallenges();
  const codes = (a.codeIds || []).map((cid) => getChallenge(cid)).filter(Boolean);

  GD_ASSESS.current = { cert, mcqs, codes, selected: new Array(mcqs.length).fill(null) };
  GD_ASSESS.editors = {};

  const mcqHTML = mcqs
    .map(
      (q, i) => `
      <div class="gd-card-sm" data-mcq="${i}">
        <div class="flex items-start gap-3 mb-3">
          <span class="grid h-7 w-7 shrink-0 place-items-center rounded-full bg-brand-100 text-brand-700 font-extrabold text-xs">${i + 1}</span>
          <p class="font-bold whitespace-pre-wrap">${typeof mdInline === "function" ? mdInline(q.question) : q.question}</p>
        </div>
        <div class="space-y-2">
          ${q.options
            .map(
              (opt) =>
                `<button type="button" class="assess-opt gd-option w-full text-left" data-q="${i}" data-opt="${escapeHTMLToEntities(
                  opt
                )}">${escapeHTMLToEntities(opt)}</button>`
            )
            .join("")}
        </div>
      </div>`
    )
    .join("");

  const codeHTML = codes
    .map(
      (c, i) => `
      <div class="gd-card-sm" data-code="${c.id}">
        <div class="flex items-center gap-2 mb-2">
          <span class="gd-chip gd-chip-brand !py-0.5 !px-2 !text-[10px]">${(c.language || "js").toUpperCase()}</span>
          <h3 class="font-extrabold">Coding task ${i + 1}: ${c.title}</h3>
        </div>
        <div class="gd-prose text-sm text-slate-700 mb-2">${typeof markdownToHtml === "function" ? markdownToHtml(c.prompt || "") : c.prompt}</div>
        <div class="gd-editor-wrap" id="assess-editor-${c.id}"></div>
      </div>`
    )
    .join("");

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-4">
      <div class="gd-card sticky top-[72px] z-20 flex items-center justify-between">
        <div>
          <h1 class="text-lg font-extrabold">${cert.title} — Final Assessment</h1>
          <p class="text-xs text-slate-500">${mcqs.length} questions + ${codes.length} coding tasks · pass at ${cert.passThreshold}%</p>
        </div>
        <div class="text-right">
          <p id="assess-timer" class="text-2xl font-extrabold text-brand-600 tabular-nums">${a.durationMin || 20}:00</p>
          <p class="text-[10px] font-bold uppercase tracking-wide text-slate-400">time left</p>
        </div>
      </div>

      <div class="space-y-3">
        <h2 class="text-sm font-extrabold uppercase tracking-wide text-slate-500">Part 1 · Reasoning</h2>
        ${mcqHTML}
      </div>

      ${codes.length ? `<div class="space-y-3"><h2 class="text-sm font-extrabold uppercase tracking-wide text-slate-500">Part 2 · Write the code</h2>${codeHTML}</div>` : ""}

      <div id="assess-result"></div>
      <button id="assess-submit" class="gd-btn gd-btn-primary gd-btn-block !py-3">Submit assessment</button>
      <a href="#certify" class="block text-center text-sm font-bold text-slate-400 hover:text-slate-600">Cancel</a>
    </div>`;

  // MCQ selection.
  htmlEl.querySelectorAll(".assess-opt").forEach((btn) => {
    btn.addEventListener("click", () => {
      const qi = Number(btn.dataset.q);
      htmlEl
        .querySelectorAll(`.assess-opt[data-q="${qi}"]`)
        .forEach((b) => b.classList.remove("gd-option-selected"));
      btn.classList.add("gd-option-selected");
      GD_ASSESS.current.selected[qi] = btn.dataset.opt;
    });
  });

  // Mount code editors.
  codes.forEach((c) => {
    const mountEl = htmlEl.querySelector("#assess-editor-" + CSS.escape(c.id));
    if (mountEl) {
      GD_ASSESS.editors[c.id] = createCodeEditor(mountEl, { doc: c.starter || "", language: c.language });
    }
  });

  htmlEl.querySelector("#assess-submit").addEventListener("click", () => gradeAssessment(htmlEl));

  // Timer.
  startAssessTimer(htmlEl, (a.durationMin || 20) * 60);
}

function startAssessTimer(htmlEl, totalSeconds) {
  if (GD_ASSESS.timer) clearInterval(GD_ASSESS.timer);
  let remaining = totalSeconds;
  const el = htmlEl.querySelector("#assess-timer");
  const tick = () => {
    const m = Math.floor(remaining / 60);
    const s = remaining % 60;
    if (el) el.textContent = `${m}:${String(s).padStart(2, "0")}`;
    if (remaining <= 0) {
      clearInterval(GD_ASSESS.timer);
      gradeAssessment(htmlEl, true);
      return;
    }
    remaining -= 1;
  };
  tick();
  GD_ASSESS.timer = setInterval(tick, 1000);
}

async function gradeAssessment(htmlEl, autoSubmitted = false) {
  if (!GD_ASSESS.current) return;
  if (GD_ASSESS.timer) clearInterval(GD_ASSESS.timer);
  const { cert, mcqs, codes, selected } = GD_ASSESS.current;
  const submitBtn = htmlEl.querySelector("#assess-submit");
  if (submitBtn) {
    submitBtn.disabled = true;
    submitBtn.classList.add("opacity-60");
    submitBtn.textContent = "Grading…";
  }

  // MCQ points.
  let mcqCorrect = 0;
  mcqs.forEach((q, i) => {
    if (selected[i] === q.answer) mcqCorrect += 1;
  });

  // Code points — run each against its tests.
  let codePass = 0;
  for (const c of codes) {
    const editor = GD_ASSESS.editors[c.id];
    const code = editor ? editor.getValue() : "";
    const res = await runChallenge(c, code);
    if (res.ok && res.results.length && res.results.every((r) => r.pass)) codePass += 1;
  }

  const total = mcqs.length + codes.length;
  const points = mcqCorrect + codePass;
  const score = total ? Math.round((points / total) * 100) : 0;
  const passed = score >= cert.passThreshold;

  // Persist the best score.
  let best = score;
  if (typeof getMeta === "function") {
    const prev = (await getMeta("cert-assess:" + cert.id, 0)) || 0;
    best = Math.max(prev, score);
    await setMeta("cert-assess:" + cert.id, best);
  }

  // If this completes the gauntlet, stamp the earned date once.
  let nowEarned = false;
  if (passed) {
    const status = await computeCertStatus(cert);
    if (status.earned && typeof getMeta === "function") {
      const existing = await getMeta("cert-earned:" + cert.id, null);
      if (!existing) {
        await setMeta("cert-earned:" + cert.id, new Date().toISOString());
        nowEarned = true;
      }
    }
  }

  const resultEl = htmlEl.querySelector("#assess-result");
  if (resultEl) {
    resultEl.innerHTML = `
      <div class="rounded-2xl border-2 p-5 text-center ${passed ? "border-grass-300 bg-grass-50" : "border-rose-200 bg-rose-50"}">
        <p class="text-4xl font-extrabold ${passed ? "text-grass-600" : "text-rose-500"}">${score}%</p>
        <p class="font-bold text-slate-700 mt-1">${mcqCorrect}/${mcqs.length} reasoning · ${codePass}/${codes.length} coding</p>
        <p class="mt-2 ${passed ? "text-grass-700" : "text-slate-600"} font-bold">
          ${
            passed
              ? nowEarned
                ? "You passed — and earned the certificate! 🎉"
                : "You passed the assessment! ✅"
              : `You need ${cert.passThreshold}% to pass.${autoSubmitted ? " (Time ran out.)" : ""} Review and try again.`
          }
        </p>
        <div class="mt-4 flex flex-wrap justify-center gap-3">
          ${passed ? `<a href="#certificate/${cert.id}" class="gd-btn gd-btn-primary">View certificate</a>` : `<a href="#assess/${cert.id}" class="gd-btn gd-btn-primary">Try again</a>`}
          <a href="#certify" class="gd-btn gd-btn-secondary">Back to certification</a>
        </div>
      </div>`;
    resultEl.scrollIntoView({ behavior: "smooth", block: "center" });
  }
  if (submitBtn) submitBtn.classList.add("hidden");

  if (nowEarned && typeof pixelCelebrate === "function") {
    pixelCelebrate({
      title: "Certified! 🎓",
      message: `You earned the ${cert.title} certificate — lessons, code, capstone, and a graded exam, all done. This is a real milestone. Go show it off.`,
      mood: "celebrate",
      cta: "See my certificate",
    });
  }
}

// --- Certificate ------------------------------------------------------------
function certVerifyId(name, certId, dateISO) {
  const hash = (s) => {
    let h = 5381;
    for (let i = 0; i < s.length; i++) h = ((h * 33) ^ s.charCodeAt(i)) >>> 0;
    return h.toString(36).toUpperCase();
  };
  return "GD-" + hash(name + "|" + certId).slice(0, 4) + "-" + hash(certId + "|" + dateISO).slice(0, 4);
}

function certificateSvg(name, cert, dateISO, verifyId) {
  const date = new Date(dateISO).toLocaleDateString(undefined, { year: "numeric", month: "long", day: "numeric" });
  const esc = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  return `<svg xmlns="http://www.w3.org/2000/svg" width="900" height="620" viewBox="0 0 900 620" role="img" aria-label="${esc(cert.title)} certificate for ${esc(name)}">
  <defs>
    <linearGradient id="cg" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#7857f7"/><stop offset="1" stop-color="#5731c8"/></linearGradient>
    <style>.t{font-family:Nunito,ui-rounded,system-ui,sans-serif}</style>
  </defs>
  <rect width="900" height="620" rx="24" fill="#ffffff"/>
  <rect x="14" y="14" width="872" height="592" rx="18" fill="none" stroke="url(#cg)" stroke-width="4"/>
  <rect x="28" y="28" width="844" height="564" rx="12" fill="none" stroke="#e7e5ff" stroke-width="2"/>
  <text x="450" y="96" text-anchor="middle" class="t" font-size="22" font-weight="800" fill="#6740e8" letter-spacing="3">GAMIFYDEV</text>
  <text x="450" y="150" text-anchor="middle" class="t" font-size="18" font-weight="700" fill="#94a3b8" letter-spacing="2">CERTIFICATE OF COMPLETION</text>
  <text x="450" y="210" text-anchor="middle" class="t" font-size="20" fill="#475569">This certifies that</text>
  <text x="450" y="280" text-anchor="middle" class="t" font-size="52" font-weight="800" fill="#1e293b">${esc(name)}</text>
  <line x1="250" y1="305" x2="650" y2="305" stroke="#e7e5ff" stroke-width="2"/>
  <text x="450" y="350" text-anchor="middle" class="t" font-size="20" fill="#475569">has successfully completed the</text>
  <text x="450" y="404" text-anchor="middle" class="t" font-size="36" font-weight="800" fill="#6740e8">${esc(cert.title)}</text>
  <text x="450" y="440" text-anchor="middle" class="t" font-size="18" fill="#64748b">track — lessons, coding sessions, a capstone project, and a graded assessment.</text>
  <text x="180" y="540" text-anchor="middle" class="t" font-size="16" font-weight="700" fill="#334155">${esc(date)}</text>
  <text x="180" y="562" text-anchor="middle" class="t" font-size="12" fill="#94a3b8">Date earned</text>
  <text x="720" y="540" text-anchor="middle" class="t" font-size="16" font-weight="700" fill="#334155">${esc(verifyId)}</text>
  <text x="720" y="562" text-anchor="middle" class="t" font-size="12" fill="#94a3b8">Verification ID</text>
  <circle cx="450" cy="540" r="34" fill="url(#cg)"/>
  <text x="450" y="551" text-anchor="middle" font-size="30">🏆</text>
</svg>`;
}

async function CertificatePage(htmlEl, id) {
  const user = (await DB.users.toArray())[0];
  if (!user) {
    window.location.hash = "";
    return;
  }
  await loadCertifications();
  const cert = getCert(id);
  const status = cert ? await computeCertStatus(cert) : null;
  if (!cert || !status || !status.earned) {
    window.location.hash = "certify";
    return;
  }
  let dateISO = (typeof getMeta === "function" && (await getMeta("cert-earned:" + id, null))) || new Date().toISOString();
  const name = user.name || "Learner";
  const verifyId = certVerifyId(name, id, dateISO);
  const svg = certificateSvg(name, cert, dateISO, verifyId);

  htmlEl.innerHTML = `
    <div class="max-w-3xl mx-auto animate-fade-up space-y-5">
      <a href="#certify" class="inline-flex items-center gap-1 text-sm font-bold text-brand-600 hover:text-brand-700">${icon("arrowLeft", "w-4 h-4")} Certification</a>
      <div class="gd-card overflow-hidden">
        <div class="rounded-2xl overflow-hidden shadow-card">${svg}</div>
        <div class="mt-4 flex flex-wrap justify-center gap-3">
          <button id="cert-download" class="gd-btn gd-btn-primary">${icon("file", "w-4 h-4")} Download certificate</button>
          <a href="#progress" class="gd-btn gd-btn-secondary">My progress</a>
        </div>
        <p class="text-center text-xs text-slate-400 mt-3">Verification ID: <span class="font-mono font-bold">${verifyId}</span></p>
      </div>
    </div>`;

  htmlEl.querySelector("#cert-download").addEventListener("click", () => {
    const blob = new Blob([svg], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `gamifydev-${id}-certificate.svg`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  });
}

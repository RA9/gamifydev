// playground.js — posts the editor's Python to /api/run and shows the result.
(function () {
  const root = document.querySelector(".pg");
  if (!root || root.dataset.enabled !== "1") return;
  const ta = document.getElementById("pgCode");
  const stdin = document.getElementById("pgStdin");
  const runBtn = document.getElementById("pgRun");
  const out = document.getElementById("pgOut");
  const meta = document.getElementById("pgMeta");
  const status = document.getElementById("pgStatus");
  if (!ta || !runBtn) return;

  const esc = (s) =>
    String(s).replace(/[&<>]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;" }[c]));

  // Tab inserts two spaces.
  ta.addEventListener("keydown", (e) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const s = ta.selectionStart, en = ta.selectionEnd;
      ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(en);
      ta.selectionStart = ta.selectionEnd = s + 2;
    }
  });
  // Cmd/Ctrl+Enter runs.
  ta.addEventListener("keydown", (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === "Enter") { e.preventDefault(); run(); }
  });

  const render = (data) => {
    let html = "";
    if (data.stdout) html += esc(data.stdout);
    if (data.stderr) html += '<span class="pg-err">' + esc(data.stderr) + "</span>";
    if (!html) html = '<span class="console-empty">(no output)</span>';
    if (data.truncated) html += '<span class="console-empty">\n… output truncated.</span>';
    out.innerHTML = html;
    const bits = [];
    if (data.timed_out) bits.push("timed out");
    else bits.push("exit " + data.exit_code);
    if (typeof data.duration_ms === "number") bits.push(data.duration_ms + " ms");
    meta.textContent = bits.join(" · ");
  };

  async function run() {
    runBtn.disabled = true;
    status.textContent = "Running…";
    out.innerHTML = '<span class="console-empty">Running on the server…</span>';
    meta.textContent = "";
    try {
      const resp = await fetch("/api/run", {
        method: "POST",
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code: ta.value, stdin: stdin ? stdin.value : "" }),
      });
      const data = await resp.json().catch(() => ({}));
      if (!resp.ok) {
        out.innerHTML = '<span class="pg-err">' + esc(data.error || ("Request failed (" + resp.status + ")")) + "</span>";
        meta.textContent = "";
      } else {
        render(data);
      }
    } catch (e) {
      out.innerHTML = '<span class="pg-err">Network error: ' + esc(e.message || e) + "</span>";
    } finally {
      runBtn.disabled = false;
      status.textContent = "";
    }
  }

  runBtn.addEventListener("click", run);
})();

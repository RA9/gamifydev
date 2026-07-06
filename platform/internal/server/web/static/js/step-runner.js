// step-runner.js — drives an interactive, step-by-step lesson in two modes:
//
//   HTML/CSS (data-lang="html"): the editor's HTML is live-previewed in a
//     same-origin iframe; each check is a JS boolean expression evaluated in the
//     parent with `doc` (the preview document) and `code` (raw source) in scope.
//
//   JavaScript (data-lang="js"): the learner's code runs inside a sandboxed
//     (scripts-only, opaque-origin) iframe that cannot reach the parent or its
//     cookies. A harness captures console output, runs the checks *inside* the
//     iframe — where the learner's functions/variables, `logs`, `code`, and
//     `document` are in scope — and posts the results back.
//
// Checks are authored (trusted), so eval of their test expressions is fine.
(function () {
  const lab = document.querySelector(".step-lab");
  if (!lab) return;
  const ta = document.getElementById("stepCode");
  const frame = document.getElementById("stepPreview");
  const consoleEl = document.getElementById("stepConsole");
  const checkBtn = document.getElementById("checkBtn");
  const nextBtn = document.getElementById("nextBtn");
  const msg = document.getElementById("stepMsg");
  const items = [...document.querySelectorAll("#stepChecks li")];
  if (!ta) return; // `frame` is absent in Python mode (no preview iframe)

  const lang = (lab.dataset.lang === "js" || lab.dataset.lang === "python") ? lab.dataset.lang : "html";
  const completeURL = lab.dataset.completeUrl;
  let done = lab.dataset.completed === "1";
  let runSeq = 0; // makes each run's srcdoc unique so the iframe always reloads

  const markComplete = () => {
    if (done) return;
    done = true;
    fetch(completeURL, { method: "POST", credentials: "same-origin" }).catch(() => {});
  };

  const applyResults = (results) => {
    let all = items.length > 0;
    items.forEach((li, i) => {
      const ok = !!results[i];
      li.classList.toggle("pass", ok);
      li.classList.toggle("fail", !ok);
      const mark = li.querySelector(".check-mark");
      if (mark) mark.textContent = ok ? "✓" : "✗";
      if (!ok) all = false;
    });
    if (all) {
      msg.textContent = "Great — all checks passed! 🎉";
      msg.className = "step-msg is-ok";
      markComplete();
      if (nextBtn) { nextBtn.hidden = false; nextBtn.focus(); }
    } else {
      msg.textContent = "Not quite — fix the items marked ✗ and check again.";
      msg.className = "step-msg is-bad";
    }
  };

  const escapeHtml = (s) =>
    String(s).replace(/[&<>]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;" }[c]));

  // Tab inserts two spaces (both modes).
  ta.addEventListener("keydown", (e) => {
    if (e.key === "Tab") {
      e.preventDefault();
      const s = ta.selectionStart, en = ta.selectionEnd;
      ta.value = ta.value.slice(0, s) + "  " + ta.value.slice(en);
      ta.selectionStart = ta.selectionEnd = s + 2;
    }
  });

  // ---------------------------------------------------------------- HTML mode
  if (lang === "html") {
    const renderPreview = () => { frame.srcdoc = ta.value; };
    let timer = null;
    ta.addEventListener("input", () => { clearTimeout(timer); timer = setTimeout(renderPreview, 250); });
    renderPreview();

    const renderAndWait = () =>
      new Promise((resolve) => {
        const done = () => { frame.removeEventListener("load", done); resolve(frame.contentDocument); };
        frame.addEventListener("load", done);
        frame.srcdoc = ta.value + "\n<!--gd" + ++runSeq + "-->";
      });

    checkBtn.addEventListener("click", async () => {
      checkBtn.disabled = true;
      const doc = await renderAndWait();
      const code = ta.value;
      const results = items.map((li) => {
        try {
          // eslint-disable-next-line no-new-func
          return !!Function("doc", "code", "return (" + li.dataset.test + ");")(doc, code);
        } catch (_) { return false; }
      });
      checkBtn.disabled = false;
      applyResults(results);
    });
    return;
  }

  // ------------------------------------------------------------------ JS mode
  const scaffold = lab.dataset.scaffold || "";

  const showConsole = (logs, error) => {
    if (!consoleEl) return;
    let html = (logs || []).map((l) => '<div class="console-line">' + escapeHtml(l) + "</div>").join("");
    if (error) html += '<div class="console-line console-error">⚠ ' + escapeHtml(error) + "</div>";
    if (!html) {
      const how = lang === "python" ? "print(…)" : "console.log(…)";
      html = '<div class="console-empty">No output — use ' + how + " to print something.</div>";
    }
    consoleEl.innerHTML = html;
  };

  const finish = (data) => {
    checkBtn.disabled = false;
    showConsole(data.logs, data.error);
    applyResults(data.results || []);
  };
  const timedOut = () => {
    checkBtn.disabled = false;
    msg.textContent = "Your code didn't finish — check for an infinite loop, then try again.";
    msg.className = "step-msg is-bad";
  };

  // --------------------------------------------------------------- Python mode
  // Runs the learner's Python via Pyodide (WASM) in a Web Worker. Checks are
  // Python expressions evaluated in the learner's namespace, with `_out`
  // (printed output) and `_code` available. The worker is killed if it hangs.
  if (lang === "python") {
    const PYBASE = "https://cdn.jsdelivr.net/pyodide/v0.26.4/full/";
    const WORKER_SRC =
      "let pyReady=null;" +
      "function getPy(){if(!pyReady){importScripts('" + PYBASE + "pyodide.js');pyReady=loadPyodide({indexURL:'" + PYBASE + "'})}return pyReady}" +
      "self.onmessage=async function(e){" +
      "var code=e.data.code,tests=e.data.tests,py;" +
      "try{py=await getPy()}catch(err){self.postMessage({type:'gd-check',results:tests.map(function(){return false}),logs:[],error:'Could not load Python: '+(err&&err.message||err)});return}" +
      "var logs=[];py.setStdout({batched:function(s){logs.push(s)}});py.setStderr({batched:function(s){logs.push(s)}});" +
      "var ns=py.toPy({});var error=null,results=[];" +
      "try{py.runPython(code,{globals:ns})}catch(err){error=String(err&&err.message||err);" +
      "var el=error.split('\\n').filter(Boolean);var mine=el.filter(function(x){return x.indexOf('<exec>')!==-1});" +
      "error=(mine.length?mine.join('\\n')+'\\n':'')+el[el.length-1]}" +
      "try{ns.set('_out',logs.join('\\n'));ns.set('_code',code)}catch(_){}" +
      "results=tests.map(function(t){try{return !!py.runPython('bool('+t+')',{globals:ns})}catch(e){return false}});" +
      "try{ns.destroy()}catch(_){}" +
      "self.postMessage({type:'gd-check',results:results,logs:logs,error:error})};";
    const workerURL = URL.createObjectURL(new Blob([WORKER_SRC], { type: "text/javascript" }));
    let worker = null;
    let pyLoaded = false;

    checkBtn.addEventListener("click", () => {
      checkBtn.disabled = true;
      if (!pyLoaded && consoleEl) {
        consoleEl.innerHTML = '<div class="console-empty">Loading Python… the first run downloads it (about 10&nbsp;MB), which can take a minute or two on a slow connection. Later runs are instant.</div>';
      }
      const code = ta.value;
      const tests = items.map((li) => li.dataset.test);
      if (!worker) worker = new Worker(workerURL);
      const w = worker;
      // Cold start pulls Pyodide (~10MB) from the CDN, so the first run gets a
      // very generous timeout; once loaded, runs are near-instant and an 8s cap
      // catches an accidental infinite loop.
      const wasLoaded = pyLoaded;
      const timer = setTimeout(() => {
        w.terminate();
        worker = null;
        checkBtn.disabled = false;
        if (wasLoaded) {
          timedOut();
        } else {
          msg.textContent = "Python is taking too long to load — check your connection and try again.";
          msg.className = "step-msg is-bad";
        }
      }, wasLoaded ? 8000 : 180000);
      w.onmessage = (e) => {
        clearTimeout(timer);
        pyLoaded = true;
        if (e.data && e.data.type === "gd-check") finish(e.data);
      };
      w.onerror = () => { clearTimeout(timer); w.terminate(); worker = null; timedOut(); };
      w.postMessage({ code, tests });
    });
    return;
  }

  if (scaffold) {
    // --- DOM lab: run in a sandboxed iframe (needs `document`) ---
    const escScript = (s) => String(s).replace(/<\/(script)/gi, "<\\/$1");
    const buildSandbox = (code, tests) => {
      const harness =
        "window.__logs=[];window.__err=null;" +
        "(function(){var o=console.log;console.log=function(){" +
        "window.__logs.push(Array.prototype.map.call(arguments,function(x){" +
        "try{return typeof x==='object'?JSON.stringify(x):String(x)}catch(e){return String(x)}}).join(' '));" +
        "o.apply(console,arguments)};})();" +
        "window.onerror=function(m){window.__err=String(m);return false};";
      const runner =
        "(function(){var logs=window.__logs;var error=window.__err;" +
        "var code=" + escScript(JSON.stringify(code)) + ";" +
        "var tests=" + escScript(JSON.stringify(tests)) + ";" +
        "var results=tests.map(function(t){try{return !!eval('('+t+')')}catch(e){return false}});" +
        "parent.postMessage({type:'gd-check',results:results,logs:logs,error:error},'*');})();";
      return (
        "<!doctype html><html><body>" + scaffold +
        "<script>" + harness + "<\/script>" +
        "<script>\n" + escScript(code) + "\n<\/script>" +
        "<script>" + runner + "<\/script>" +
        "</body></html>"
      );
    };
    checkBtn.addEventListener("click", () => {
      checkBtn.disabled = true;
      const code = ta.value;
      const tests = items.map((li) => li.dataset.test);
      const onMsg = (e) => {
        if (!e || !e.data || e.data.type !== "gd-check") return;
        window.removeEventListener("message", onMsg);
        clearTimeout(timer);
        finish(e.data);
      };
      window.addEventListener("message", onMsg);
      const timer = setTimeout(() => { window.removeEventListener("message", onMsg); timedOut(); }, 5000);
      frame.srcdoc = buildSandbox(code, tests) + "<!--gd" + ++runSeq + "-->";
    });
    return;
  }

  // --- Pure-logic JS: run in a Web Worker so an infinite loop can be killed ---
  // Learner code + checks run in one eval (so checks see the learner's top-level
  // let/const); the worker has no DOM, and is terminated if it exceeds the timeout.
  const WORKER_SRC =
    "self.onmessage=function(e){" +
    "var code=e.data.code,tests=e.data.tests,logs=[];" +
    "self.console={log:function(){logs.push(Array.prototype.map.call(arguments,function(x){try{return typeof x==='object'?JSON.stringify(x):String(x)}catch(_){return String(x)}}).join(' '))}};" +
    "var checks=tests.map(function(t){return '(function(){try{return !!('+t+')}catch(e){return false}})()'}).join(',');" +
    "var program=code+'\\n;globalThis.__gdResults=['+checks+'];';" +
    "var error=null,results=[];" +
    "try{eval(program);results=globalThis.__gdResults||[]}catch(err){error=String(err&&err.message||err);results=tests.map(function(){return false})}" +
    "self.postMessage({type:'gd-check',results:results,logs:logs,error:error})};";
  const workerURL = URL.createObjectURL(new Blob([WORKER_SRC], { type: "text/javascript" }));

  checkBtn.addEventListener("click", () => {
    checkBtn.disabled = true;
    const code = ta.value;
    const tests = items.map((li) => li.dataset.test);
    const worker = new Worker(workerURL);
    const timer = setTimeout(() => { worker.terminate(); timedOut(); }, 3000);
    worker.onmessage = (e) => {
      clearTimeout(timer);
      worker.terminate();
      if (e.data && e.data.type === "gd-check") finish(e.data);
    };
    worker.onerror = () => { clearTimeout(timer); worker.terminate(); timedOut(); };
    worker.postMessage({ code, tests });
  });
})();

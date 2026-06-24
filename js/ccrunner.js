// C runner — executes learner C in JSCPP (a C/C++ interpreter in pure JS,
// vendored offline as js/external/jscpp.bundle.js). C challenges use the classic
// stdin -> stdout autograder model: the program reads input and prints output,
// and each test compares the trimmed stdout to the expected output. Runs inside
// a Web Worker so a runaway loop can be terminated without freezing the page.

const GD_CC = { worker: null };

function gdCcUrl() {
  const base = (self.location && self.location.href) || location.href;
  return new URL("js/external/jscpp.bundle.js", base).href;
}

function gdCreateCcWorker() {
  const src = `
    var __loaded = false;
    self.onmessage = function (e) {
      var m = e.data;
      if (!__loaded) {
        try { importScripts(m.jscppUrl); __loaded = true; }
        catch (err) { self.postMessage({ ok: false, compileError: 'Failed to load the C runtime.', results: [], logs: [] }); return; }
      }
      var code = m.code, tests = m.tests || [];
      var results = [], firstErr = null, allErr = true;
      for (var i = 0; i < tests.length; i++) {
        var t = tests[i], out = '', err = null;
        try {
          JSCPP.run(code, t.input || '', { stdio: { write: function (s) { out += s; } } });
        } catch (ex) { err = String(ex && ex.message || ex); }
        if (err) { if (firstErr === null) firstErr = err; } else { allErr = false; }
        var got = (out || '').trim();
        var exp = String(t.expected == null ? '' : t.expected).trim();
        results.push({ name: t.name, pass: !err && got === exp, got: got, expected: exp, error: err, hidden: !!t.hidden, args: t.hidden ? undefined : [t.input] });
      }
      if (allErr && firstErr) { self.postMessage({ ok: false, compileError: firstErr, results: [], logs: [] }); return; }
      self.postMessage({ ok: true, results: results, logs: [] });
    };
  `;
  const blob = new Blob([src], { type: "application/javascript" });
  return new Worker(URL.createObjectURL(blob));
}

function runCChallenge(code, tests, timeoutMs = 6000) {
  return new Promise((resolve) => {
    const w = GD_CC.worker || (GD_CC.worker = gdCreateCcWorker());
    let done = false;
    const finish = (p) => {
      if (done) return;
      done = true;
      resolve(p);
    };
    const timer = setTimeout(() => {
      try { w.terminate(); } catch (e) {}
      GD_CC.worker = null;
      finish({
        ok: false,
        timeout: true,
        compileError: "Your code took too long to run (possible infinite loop).",
        results: [],
        logs: [],
      });
    }, timeoutMs);
    const onmsg = (e) => {
      clearTimeout(timer);
      w.removeEventListener("message", onmsg);
      finish(e.data);
    };
    w.addEventListener("message", onmsg);
    try {
      w.postMessage({ type: "run", code, tests, jscppUrl: gdCcUrl() });
    } catch (err) {
      clearTimeout(timer);
      GD_CC.worker = null;
      finish({ ok: false, compileError: String(err && err.message || err), results: [], logs: [] });
    }
  });
}

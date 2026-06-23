// Sandboxed code runner + test grader.
//
// Learner code is executed in a throwaway Web Worker created from a Blob URL,
// so it runs off the main thread, can't touch the DOM, and can be hard-killed
// if it hangs (infinite-loop protection). The worker runs the code, then calls
// the entry function with each test's args and deep-compares the result to the
// expected value. This is the objective half of grading: code either passes the
// hidden tests or it doesn't.
//
// Languages other than JavaScript (Python via Pyodide, C/Java via WASM) plug in
// behind the same runChallenge() interface in later phases; for now non-JS
// languages report a friendly "runtime loading soon" result so the editor and
// authoring flow can be built and tested today.

const JS_WORKER_SOURCE = `
self.onmessage = function (e) {
  const { code, entry, tests, consoleCapture } = e.data;
  const logs = [];
  const realLog = function () {};
  // Sandboxed console: collect output, never touch the page.
  const sandboxConsole = {
    log: (...a) => logs.push(a.map(String).join(' ')),
    error: (...a) => logs.push(a.map(String).join(' ')),
    warn: (...a) => logs.push(a.map(String).join(' ')),
    info: (...a) => logs.push(a.map(String).join(' ')),
  };

  function deepEqual(a, b) {
    if (a === b) return true;
    if (typeof a !== typeof b) return false;
    if (a && b && typeof a === 'object') {
      const ka = Object.keys(a), kb = Object.keys(b);
      if (Array.isArray(a) !== Array.isArray(b)) return false;
      if (ka.length !== kb.length) return false;
      return ka.every((k) => deepEqual(a[k], b[k]));
    }
    return false;
  }

  let fn;
  try {
    // Define the learner's code, then hand back the entry function.
    const factory = new Function('console', code + '\\n; return typeof ' + entry + ' === "function" ? ' + entry + ' : undefined;');
    fn = factory(sandboxConsole);
  } catch (err) {
    self.postMessage({ ok: false, compileError: String(err && err.message || err), results: [], logs });
    return;
  }
  if (typeof fn !== 'function') {
    self.postMessage({ ok: false, compileError: 'Could not find a function named "' + entry + '". Make sure it is defined.', results: [], logs });
    return;
  }

  const results = (tests || []).map((t) => {
    try {
      const out = fn.apply(null, t.args || []);
      const pass = deepEqual(out, t.expected);
      return { name: t.name, pass, got: safe(out), expected: safe(t.expected), hidden: !!t.hidden, args: t.hidden ? undefined : (t.args || []) };
    } catch (err) {
      return { name: t.name, pass: false, error: String(err && err.message || err), expected: safe(t.expected), hidden: !!t.hidden, args: t.hidden ? undefined : (t.args || []) };
    }
  });
  self.postMessage({ ok: true, results, logs });

  function safe(v) {
    try { return JSON.stringify(v); } catch (e) { return String(v); }
  }
};
`;

// Run a JavaScript challenge in a worker with a timeout.
function runJsChallenge(code, entry, tests, timeoutMs = 4000) {
  return new Promise((resolve) => {
    let worker;
    let done = false;
    const blob = new Blob([JS_WORKER_SOURCE], { type: "application/javascript" });
    const url = URL.createObjectURL(blob);
    const finish = (payload) => {
      if (done) return;
      done = true;
      try { worker && worker.terminate(); } catch (e) {}
      URL.revokeObjectURL(url);
      resolve(payload);
    };
    const timer = setTimeout(
      () =>
        finish({
          ok: false,
          timeout: true,
          compileError: "Your code took too long to run (possible infinite loop).",
          results: [],
          logs: [],
        }),
      timeoutMs
    );
    try {
      worker = new Worker(url);
      worker.onmessage = (e) => {
        clearTimeout(timer);
        finish(e.data);
      };
      worker.onerror = (e) => {
        clearTimeout(timer);
        finish({ ok: false, compileError: String(e.message || "Runtime error"), results: [], logs: [] });
      };
      worker.postMessage({ code, entry, tests });
    } catch (err) {
      clearTimeout(timer);
      finish({ ok: false, compileError: String(err.message || err), results: [], logs: [] });
    }
  });
}

// Unified entry. Routes by language; non-JS runtimes arrive in later phases.
async function runChallenge(challenge, code) {
  const lang = (challenge.language || "javascript").toLowerCase();
  if (lang === "javascript" || lang === "js") {
    return runJsChallenge(code, challenge.entry, challenge.tests || []);
  }
  if ((lang === "python" || lang === "py") && typeof runPythonChallenge === "function") {
    return runPythonChallenge(code, challenge.entry, challenge.tests || []);
  }
  return {
    ok: false,
    pending: true,
    compileError: `The ${lang.toUpperCase()} runtime is being added in an upcoming update. JavaScript and Python challenges run fully today.`,
    results: [],
    logs: [],
  };
}

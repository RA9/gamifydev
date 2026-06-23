// Python runner — executes learner Python in Pyodide (CPython compiled to WASM).
//
// Pyodide is lazy-loaded from the CDN inside a dedicated Web Worker on first use
// (a one-time ~10MB download, then browser-cached). Running it in a worker keeps
// the page responsive and — crucially — lets us terminate a runaway program
// (e.g. `while True:`) the same way the JS runner does. The worker returns the
// same result shape as the JS runner so coding.js/certify.js render it unchanged.

const GD_PY_INDEX = "https://cdn.jsdelivr.net/pyodide/v0.26.4/full/";
const GD_PY = { worker: null, ready: null };

// The grading harness, authored here (single escaping layer) and sent to the
// worker. Uses single-quoted Python and rstrip() with no args so there are no
// escape sequences to mangle. Reads __entry / __tests_json from globals and
// returns a JSON string with the same shape as the JS runner.
const GD_PY_HARNESS = [
  "import json, sys, io",
  "def __safe(v):",
  "    try: return json.dumps(v)",
  "    except Exception: return json.dumps(str(v))",
  "__logs = []",
  "__results = []",
  "__fn = globals().get(__entry)",
  "if not callable(__fn):",
  "    __payload = {'ok': False, 'compileError': 'Define a function named ' + str(__entry) + '.', 'results': [], 'logs': []}",
  "else:",
  "    for __t in json.loads(__tests_json):",
  "        __old = sys.stdout",
  "        try:",
  "            sys.stdout = io.StringIO()",
  "            __r = __fn(*__t['args'])",
  "            __cap = sys.stdout.getvalue()",
  "            sys.stdout = __old",
  "            if __cap.strip(): __logs.append(__cap.rstrip())",
  "            __results.append({'name': __t.get('name'), 'pass': bool(__r == __t['expected']), 'got': __safe(__r), 'expected': __safe(__t['expected']), 'hidden': bool(__t.get('hidden', False)), 'args': (None if __t.get('hidden') else __t.get('args'))})",
  "        except Exception as __e:",
  "            sys.stdout = __old",
  "            __results.append({'name': __t.get('name'), 'pass': False, 'error': str(__e), 'expected': __safe(__t['expected']), 'hidden': bool(__t.get('hidden', False)), 'args': (None if __t.get('hidden') else __t.get('args'))})",
  "    __payload = {'ok': True, 'compileError': None, 'results': __results, 'logs': __logs}",
  "json.dumps(__payload)",
].join("\n");

function gdCreatePyWorker() {
  const src = `
    self.onmessage = async (e) => {
      const msg = e.data;
      if (msg.type === 'init') {
        try {
          importScripts(msg.indexURL + 'pyodide.js');
          self.pyodide = await loadPyodide({ indexURL: msg.indexURL });
          self.postMessage({ type: 'ready' });
        } catch (err) {
          self.postMessage({ type: 'initError', error: String(err && err.message || err) });
        }
        return;
      }
      if (msg.type === 'run') {
        const py = self.pyodide;
        const { code, entry, tests, harness } = msg;
        let compileError = null;
        try {
          py.runPython(code);
        } catch (err) {
          compileError = String(err && err.message || err);
        }
        if (compileError) {
          self.postMessage({ ok: false, compileError, results: [], logs: [] });
          return;
        }
        try {
          py.globals.set('__entry', entry);
          py.globals.set('__tests_json', JSON.stringify(tests || []));
          const out = await py.runPythonAsync(harness);
          self.postMessage(JSON.parse(out));
        } catch (err) {
          self.postMessage({ ok: false, compileError: String(err && err.message || err), results: [], logs: [] });
        }
      }
    };
  `;
  const blob = new Blob([src], { type: "application/javascript" });
  return new Worker(URL.createObjectURL(blob));
}

// Resolve once Pyodide has finished loading in the worker.
function ensurePyReady() {
  if (GD_PY.ready) return GD_PY.ready;
  GD_PY.ready = new Promise((resolve, reject) => {
    let w;
    try {
      w = gdCreatePyWorker();
    } catch (err) {
      reject(err);
      return;
    }
    GD_PY.worker = w;
    const to = setTimeout(() => {
      reject(new Error("Python runtime took too long to load (no connection?)."));
    }, 60000);
    w.addEventListener("message", function onmsg(e) {
      if (!e.data) return;
      if (e.data.type === "ready") {
        clearTimeout(to);
        resolve(w);
      } else if (e.data.type === "initError") {
        clearTimeout(to);
        GD_PY.worker = null;
        GD_PY.ready = null;
        reject(new Error(e.data.error));
      }
    });
    w.postMessage({ type: "init", indexURL: GD_PY_INDEX });
  });
  return GD_PY.ready;
}

// True once Pyodide is loaded (so the UI can show a one-time loading hint).
function isPyReady() {
  return !!(GD_PY.worker && GD_PY.ready);
}

function runPythonChallenge(code, entry, tests, timeoutMs = 8000) {
  return new Promise(async (resolve) => {
    let w;
    try {
      w = await ensurePyReady();
    } catch (err) {
      resolve({
        ok: false,
        compileError:
          String(err && err.message || err) +
          " Python runs in your browser via a one-time download — check your connection and try again.",
        results: [],
        logs: [],
      });
      return;
    }
    let done = false;
    const finish = (p) => {
      if (done) return;
      done = true;
      resolve(p);
    };
    const timer = setTimeout(() => {
      // Likely an infinite loop — kill the worker so it can't freeze anything,
      // and reset so the next run rebuilds a fresh interpreter.
      try { w.terminate(); } catch (e) {}
      GD_PY.worker = null;
      GD_PY.ready = null;
      finish({
        ok: false,
        timeout: true,
        compileError: "Your code took too long to run (possible infinite loop).",
        results: [],
        logs: [],
      });
    }, timeoutMs);
    const onmsg = (e) => {
      if (e.data && (e.data.type === "ready" || e.data.type === "initError")) return;
      clearTimeout(timer);
      w.removeEventListener("message", onmsg);
      finish(e.data);
    };
    w.addEventListener("message", onmsg);
    w.postMessage({ type: "run", code, entry, tests, harness: GD_PY_HARNESS });
  });
}

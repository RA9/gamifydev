package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/runner"
)

// resultMarker prefixes the JSON line the check harness prints, so we can find
// it in the program's stdout amid the learner's own output.
const resultMarker = "__GD_RESULT__"

// handleStepRun executes a server-side lab step: it runs the learner's Python
// (a Flask/FastAPI app or plain program) plus the step's authored checks inside
// the sandbox, and returns per-check pass/fail. Checks live server-side, so —
// unlike the client-run labs — a learner can't read or tamper with them.
func (s *Server) handleStepRun(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context()) // guaranteed by requireAuth
	writeJSON := func(status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(http.StatusBadRequest, map[string]string{"error": "bad step"})
		return
	}
	step, err := s.st.GetStep(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if step.Lang != "pyserver" {
		writeJSON(http.StatusBadRequest, map[string]string{"error": "this step is not a server-side lab"})
		return
	}
	if !s.exec.Enabled() {
		writeJSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Server-side labs are turned off on this instance.",
		})
		return
	}
	if wait, ok := s.runlim.allow(u.ID); !ok {
		w.Header().Set("Retry-After", "2")
		writeJSON(http.StatusTooManyRequests, map[string]any{
			"error": "You're checking too quickly — give it a second.", "retry_after_ms": wait.Milliseconds(),
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxProgramBytes+2*1024)
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Code) > maxProgramBytes {
		writeJSON(http.StatusBadRequest, map[string]string{"error": "Could not read your code."})
		return
	}

	var checks []stepCheck
	_ = json.Unmarshal([]byte(step.Checks), &checks)
	tests := make([]string, len(checks))
	for i, c := range checks {
		tests[i] = c.Test
	}

	rel, ok := s.runlim.acquire(r.Context())
	if !ok {
		writeJSON(http.StatusServiceUnavailable, map[string]string{"error": "The runner is busy — try again in a moment."})
		return
	}
	defer rel()

	res, err := s.exec.Run(r.Context(), runner.Request{Code: buildPyHarness(body.Code, tests)})
	if err != nil {
		writeJSON(http.StatusBadGateway, map[string]string{"error": "The runner had a problem: " + err.Error()})
		return
	}

	results, logs := parseHarnessOutput(res.Stdout, len(tests))
	errText := ""
	if res.TimedOut {
		errText = "Your code didn't finish in time — make sure you're not calling app.run() (the checker runs your app for you)."
	} else if strings.TrimSpace(res.Stderr) != "" {
		errText = res.Stderr
	}
	writeJSON(http.StatusOK, map[string]any{
		"results":  results,
		"logs":     logs,
		"error":    errText,
		"timedOut": res.TimedOut,
	})
}

// parseHarnessOutput pulls the results line out of stdout and returns it plus
// the learner-visible console lines (everything else). If the marker is absent
// (e.g. the program crashed before the harness ran), every check is false.
func parseHarnessOutput(stdout string, n int) ([]bool, []string) {
	results := make([]bool, n)
	var logs []string
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, resultMarker) {
			var got []bool
			if err := json.Unmarshal([]byte(line[len(resultMarker):]), &got); err == nil {
				for i := range results {
					if i < len(got) {
						results[i] = got[i]
					}
				}
			}
			continue
		}
		if line != "" {
			logs = append(logs, line)
		}
	}
	return results, logs
}

// buildPyHarness wraps the learner's program with a check harness. The harness
// neutralizes accidental server starts (app.run / uvicorn.run block forever),
// exposes a `client` test client for whatever web app the code defines, then
// evaluates each authored check and prints the results as a JSON line.
//
// The check expressions are authored (trusted) content; the learner's code is
// untrusted but runs inside the sandbox, so evaluating them together is safe.
func buildPyHarness(code string, tests []string) string {
	testsJSON, _ := json.Marshal(tests)
	var b strings.Builder
	b.WriteString("import json as _gd_json\n")
	// Stop a stray app.run()/uvicorn.run() from blocking the whole run.
	b.WriteString("try:\n    import flask as _gd_flask\n    _gd_flask.Flask.run = lambda *a, **k: None\nexcept Exception:\n    pass\n")
	b.WriteString("try:\n    import uvicorn as _gd_uv\n    _gd_uv.run = lambda *a, **k: None\nexcept Exception:\n    pass\n")
	b.WriteString("\n# ---- learner program ----\n")
	b.WriteString(code)
	b.WriteString("\n# ---- GamifyDev check harness ----\n")
	b.WriteString("def _gd_make_client():\n")
	b.WriteString("    app = globals().get('app')\n")
	b.WriteString("    if app is None:\n        return None\n")
	b.WriteString("    if hasattr(app, 'test_client'):\n        return app.test_client()\n") // Flask / WSGI
	b.WriteString("    try:\n        from starlette.testclient import TestClient\n        return TestClient(app)\n    except Exception:\n        return None\n")
	b.WriteString("try:\n    client = _gd_make_client()\nexcept Exception:\n    client = None\n")
	// Expose the learner's source so checks can inspect it (e.g. testing labs
	// that confirm an `assert` was actually written).
	codeJSON, _ := json.Marshal(code)
	b.WriteString("_code = ")
	b.Write(codeJSON)
	b.WriteString("\n_gd_tests = ")
	b.Write(testsJSON)
	b.WriteString("\n_gd_results = []\n")
	b.WriteString("for _gd_t in _gd_tests:\n")
	b.WriteString("    try:\n        _gd_results.append(bool(eval(_gd_t)))\n    except Exception:\n        _gd_results.append(False)\n")
	b.WriteString("print(\"" + resultMarker + "\" + _gd_json.dumps(_gd_results))\n")
	return b.String()
}

// Package pyharness wraps a learner's Python program with an authored check
// harness and reads the results back out.
//
// Extracted from the server-side lab runner so lesson steps and assignment
// checkpoints are graded by exactly the same code. Two harnesses that drifted
// apart would mean a learner's lab passing and their checkpoint failing on
// identical logic.
package pyharness

import (
	"encoding/json"
	"strings"
)

// ResultMarker prefixes the JSON line the harness prints, so the caller can
// find it among the learner's own output.
const ResultMarker = "__GD_RESULT__"

// parseHarnessOutput pulls the results line out of stdout and returns it plus
// the learner-visible console lines (everything else). If the marker is absent
// (e.g. the program crashed before the harness ran), every check is false.
func ParseOutput(stdout string, n int) ([]bool, []string) {
	results := make([]bool, n)
	var logs []string
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, ResultMarker) {
			var got []bool
			if err := json.Unmarshal([]byte(line[len(ResultMarker):]), &got); err == nil {
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
func Build(code string, tests []string) string {
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
	b.WriteString("print(\"" + ResultMarker + "\" + _gd_json.dumps(_gd_results))\n")
	return b.String()
}

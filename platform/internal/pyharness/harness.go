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

// ParseOutput pulls the result lines out of stdout and returns them plus the
// learner-visible console lines (everything else).
//
// Each check reports on its own line as it finishes, so a run that dies partway
// — a timeout on the check that feeds it a large input, say — still returns
// everything proven up to that point. A check that never reported is false,
// which is also what happens when the marker is absent entirely because the
// program crashed before the harness ran.
func ParseOutput(stdout string, n int) ([]bool, []string) {
	results := make([]bool, n)
	var logs []string
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, ResultMarker) {
			var pair []any
			if err := json.Unmarshal([]byte(line[len(ResultMarker):]), &pair); err == nil && len(pair) == 2 {
				idx, iok := pair[0].(float64)
				val, vok := pair[1].(bool)
				if iok && vok && int(idx) >= 0 && int(idx) < n {
					results[int(idx)] = val
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

// evalLoop is the tail of every harness: evaluate each authored expression and
// report it immediately.
//
// Shared by both builders so the two can't drift into reporting results in
// different shapes, which ParseOutput would then read wrong for one of them.
// flush=True matters — without it a killed process loses its buffer, which is
// exactly the run whose partial results are most worth having.
const evalLoop = `for _gd_i, _gd_t in enumerate(_gd_tests):
    try:
        _gd_r = bool(eval(_gd_t))
    except Exception:
        _gd_r = False
    print("` + ResultMarker + `" + _gd_json.dumps([_gd_i, _gd_r]), flush=True)
`

// Build wraps the learner's program with a check harness. The harness
// neutralizes accidental server starts (app.run / uvicorn.run block forever),
// exposes a `client` test client for whatever web app the code defines, then
// evaluates each authored check, reporting each one as it finishes.
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
	b.WriteString("\n")
	b.WriteString(evalLoop)
	return b.String()
}

// BuildAssertions builds a Python program that binds a set of values as
// literals and evaluates authored expressions against them.
//
// This is how a compiled language is checked. A C program leaves no namespace
// to inspect, so its stdout, stderr and exit code are captured, bound here as
// `_out`, `_err`, `_exit` and `_code`, and the authored expressions assert over
// those — `_out.strip() == "9"` rather than `largest(3,9,4) == 9`.
//
// Reusing Python for the assertions keeps one expression language across the
// platform: an author who can write a lesson-step check can write a C
// checkpoint check without learning anything new.
//
// `strs` and `ints` are separate because an exit code bound as a string would
// silently make `_exit == 0` false.
// AssertGroup is one program run's results and the checks made against them.
type AssertGroup struct {
	Strings map[string]string
	Ints    map[string]int
	Tests   []string
}

// BuildGroupedAssertions evaluates the checks for many runs in one program.
//
// Grading a submission runs the learner's program once per test input, and each
// of those used to be followed by its own Python process just to decide whether
// the output was right. That doubled the number of sandboxed processes for work
// with no reason to be spread across them: the groups do not depend on each
// other, and every value each one needs is already known by the time any of them
// is checked. One process now settles all of them.
//
// Results carry a running index across every group, so the caller maps them back
// to checks with the same ParseOutput as before.
func BuildGroupedAssertions(groups []AssertGroup) string {
	var b strings.Builder
	b.WriteString("import json as _gd_json\n")
	base := 0
	for _, g := range groups {
		for name, value := range g.Strings {
			enc, _ := json.Marshal(value)
			b.WriteString(name + " = ")
			b.Write(enc)
			b.WriteString("\n")
		}
		for name, value := range g.Ints {
			b.WriteString(name + " = " + itoa(value) + "\n")
		}
		testsJSON, _ := json.Marshal(g.Tests)
		b.WriteString("_gd_tests = ")
		b.Write(testsJSON)
		b.WriteString("\n_gd_base = " + itoa(base) + "\n")
		b.WriteString(groupedLoop)
		base += len(g.Tests)
	}
	return b.String()
}

// groupedLoop is evalLoop with an offset, so one program can report results for
// several runs without the indices colliding.
const groupedLoop = `for _gd_i, _gd_t in enumerate(_gd_tests):
    try:
        _gd_r = bool(eval(_gd_t))
    except Exception:
        _gd_r = False
    print("` + ResultMarker + `" + _gd_json.dumps([_gd_base + _gd_i, _gd_r]), flush=True)
`

func BuildAssertions(strs map[string]string, ints map[string]int, tests []string) string {
	var b strings.Builder
	b.WriteString("import json as _gd_json\n")
	for name, value := range strs {
		enc, _ := json.Marshal(value)
		b.WriteString(name + " = ")
		b.Write(enc)
		b.WriteString("\n")
	}
	for name, value := range ints {
		b.WriteString(name + " = " + itoa(value) + "\n")
	}
	testsJSON, _ := json.Marshal(tests)
	b.WriteString("_gd_tests = ")
	b.Write(testsJSON)
	b.WriteString("\n")
	b.WriteString(evalLoop)
	return b.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var d [20]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		d[i] = '-'
	}
	return string(d[i:])
}

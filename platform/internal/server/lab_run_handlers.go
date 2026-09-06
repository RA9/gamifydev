package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/pyharness"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/store"
)

// handleStepRun executes a server-side Python, C, or shell lab step plus its
// authored checks inside the sandbox. Checks live server-side, so unlike the
// client-run labs a learner cannot read or tamper with them.
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
	if step.Lang != "pyserver" && step.Lang != runner.LangC && step.Lang != runner.LangShell {
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

	rel, ok := s.runlim.acquire(r.Context())
	if !ok {
		writeJSON(http.StatusServiceUnavailable, map[string]string{"error": "The runner is busy — try again in a moment."})
		return
	}
	defer rel()

	if step.Lang == runner.LangC || step.Lang == runner.LangShell {
		gradeChecks := make([]store.Check, len(checks))
		for i, check := range checks {
			gradeChecks[i] = store.Check{ID: int64(i + 1), Sort: i, Label: check.Text, Test: check.Test, Points: 1}
		}
		grade, err := jobs.Grade(r.Context(), s.exec, step.Lang, body.Code, nil, gradeChecks)
		if err != nil {
			writeJSON(http.StatusBadGateway, map[string]string{"error": "The runner had a problem: " + err.Error()})
			return
		}
		results := make([]bool, len(gradeChecks))
		for i, check := range gradeChecks {
			results[i] = grade.Passed[check.ID]
		}
		errText := ""
		if grade.TimedOut {
			errText = "Your code didn't finish in time — check for an infinite loop."
		} else if grade.CompileFailed {
			errText = "Your program did not compile. Review the compiler output below."
		} else if grade.Crashed {
			errText = "Your program exited before the checks completed."
		}
		logs := []string{}
		if strings.TrimSpace(grade.Output) != "" {
			logs = append(logs, grade.Output)
		}
		writeJSON(http.StatusOK, map[string]any{
			"results": results, "logs": logs, "error": errText, "timedOut": grade.TimedOut,
		})
		return
	}

	tests := make([]string, len(checks))
	for i, check := range checks {
		tests[i] = check.Test
	}
	res, err := s.exec.Run(r.Context(), runner.Request{Code: pyharness.Build(body.Code, tests)})
	if err != nil {
		writeJSON(http.StatusBadGateway, map[string]string{"error": "The runner had a problem: " + err.Error()})
		return
	}

	results, logs := pyharness.ParseOutput(res.Stdout, len(tests))
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

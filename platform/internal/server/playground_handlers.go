package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/runner"
)

// maxProgramBytes caps the size of a submitted program.
const maxProgramBytes = 64 * 1024

// handlePlayground renders the server-side Python scratchpad.
func (s *Server) handlePlayground(w http.ResponseWriter, r *http.Request) {
	starter := "# This Python runs on the server, so things the browser labs can't do\n" +
		"# work here: input(), the filesystem, and the full standard library.\n\n" +
		"name = input(\"What's your name? \")\n" +
		"print(f\"Hello, {name}! Welcome to server-side Python.\")\n"
	s.render(w, r, "playground.html", ViewData{
		Title: "Playground",
		Data: map[string]any{
			"starter": starter,
			"enabled": s.exec.Enabled(),
			"kind":    s.exec.Kind(),
		},
	})
}

// handleRunCode executes a submitted program and returns the result as JSON.
func (s *Server) handleRunCode(w http.ResponseWriter, r *http.Request) {
	u := auth.CurrentUser(r.Context()) // guaranteed by requireAuth
	writeJSON := func(status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}

	if !s.exec.Enabled() {
		writeJSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Server-side code execution is turned off on this instance.",
		})
		return
	}
	if wait, ok := s.runlim.allow(u.ID); !ok {
		w.Header().Set("Retry-After", "2")
		writeJSON(http.StatusTooManyRequests, map[string]any{
			"error": "You're running code too quickly — give it a second.", "retry_after_ms": wait.Milliseconds(),
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxProgramBytes+2*1024)
	var req runner.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(http.StatusBadRequest, map[string]string{"error": "Could not read your code."})
		return
	}
	if len(req.Code) > maxProgramBytes {
		writeJSON(http.StatusRequestEntityTooLarge, map[string]string{"error": "That program is too large."})
		return
	}

	// Bound how long we're willing to hold the connection, independent of the
	// executor's own per-run timeout.
	rel, ok := s.runlim.acquire(r.Context())
	if !ok {
		writeJSON(http.StatusServiceUnavailable, map[string]string{"error": "The runner is busy — try again in a moment."})
		return
	}
	defer rel()

	res, err := s.exec.Run(r.Context(), req)
	if err != nil {
		writeJSON(http.StatusBadGateway, map[string]string{"error": "The runner had a problem: " + err.Error()})
		return
	}
	writeJSON(http.StatusOK, res)
}

// runLimiter throttles code execution: a per-user minimum interval between runs
// plus a global cap on how many run at once.
type runLimiter struct {
	mu          sync.Mutex
	last        map[int64]time.Time
	minInterval time.Duration
	slots       chan struct{}
}

func newRunLimiter(maxConcurrent int, minIntervalMs int) *runLimiter {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	return &runLimiter{
		last:        map[int64]time.Time{},
		minInterval: time.Duration(minIntervalMs) * time.Millisecond,
		slots:       make(chan struct{}, maxConcurrent),
	}
}

// allow enforces the per-user minimum interval; it records "now" when it says
// yes. On a no it returns how long the caller should wait.
func (l *runLimiter) allow(userID int64) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if last, ok := l.last[userID]; ok {
		if wait := l.minInterval - now.Sub(last); wait > 0 {
			return wait, false
		}
	}
	l.last[userID] = now
	return 0, true
}

// acquire takes a global concurrency slot, waiting until one frees up or the
// request is cancelled. The returned func releases it.
func (l *runLimiter) acquire(ctx context.Context) (func(), bool) {
	select {
	case l.slots <- struct{}{}:
		return func() { <-l.slots }, true
	case <-ctx.Done():
		return func() {}, false
	}
}

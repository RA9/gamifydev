// Command execd is a tiny, standalone HTTP service that executes untrusted
// Python and returns its output. It is the sandbox tier: run it in its own
// locked-down container (see Dockerfile.execd) with bubblewrap available, and
// point the main app at it with CODE_EXEC=remote + CODE_EXEC_SANDBOX_URL.
//
// Keeping execution in a separate service means the main app image stays lean
// (distroless, no interpreter) and a sandbox escape lands in a disposable
// container, not the app that holds the database credentials.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/RA9/gamifydev/platform/internal/runner"
)

func main() {
	addr := envOr("EXECD_ADDR", ":9090")
	token := os.Getenv("EXECD_TOKEN")
	if token == "" && os.Getenv("EXECD_ALLOW_NOAUTH") != "1" {
		log.Fatal("execd: set EXECD_TOKEN to a shared secret (or EXECD_ALLOW_NOAUTH=1 to run without auth)")
	}

	// execd IS the hardened container, so unsandboxed local is accepted here —
	// but we still log loudly when bubblewrap isn't providing isolation.
	exec, err := runner.New(runner.Config{
		Mode:        "local",
		PythonPath:  envOr("PYTHON", "python3"),
		AllowUnsafe: true,
		Limits:      runner.DefaultLimits(),
	})
	if err != nil {
		log.Fatalf("execd: %v", err)
	}

	// Bound how many runs execute at once so a burst can't fork-bomb the box.
	maxConc, _ := strconv.Atoi(envOr("EXECD_MAX_CONCURRENCY", "4"))
	if maxConc < 1 {
		maxConc = 1
	}
	slots := make(chan struct{}, maxConc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /run", func(w http.ResponseWriter, r *http.Request) {
		if token != "" && r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 256*1024)
		var req runner.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		select {
		case slots <- struct{}{}:
			defer func() { <-slots }()
		case <-r.Context().Done():
			http.Error(w, "busy", http.StatusServiceUnavailable)
			return
		}
		res, err := exec.Run(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	if runner.SandboxAvailable() {
		log.Printf("execd: bubblewrap sandbox active (network off, read-only root)")
	} else {
		log.Printf("execd: WARNING — no bubblewrap on PATH; runs are NOT namespace-isolated. " +
			"Only run this where the container itself is the security boundary.")
	}
	log.Printf("execd listening on %s (max concurrency %d)", addr, maxConc)
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

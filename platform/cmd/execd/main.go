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
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/RA9/gamifydev/platform/internal/runner"
)

func main() {
	// `execd selftest` runs the containment probe once and exits 0 (sandboxed)
	// or 1 (not) — for CI, health checks, and manual validation on a real host.
	if len(os.Args) > 1 && os.Args[1] == "selftest" {
		runSelftest()
		return
	}

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
		CCPath:      envOr("CC", "cc"),
		AllowUnsafe: true,
		Limits:      runner.DefaultLimits(),
	})
	if err != nil {
		log.Fatalf("execd: %v", err)
	}

	// Prove, at boot, that untrusted code is actually contained. This is the
	// difference between "we configured a sandbox" and "the sandbox works on
	// THIS host/kernel". If a sandbox is required (the container default) and
	// the probe can't confirm containment, refuse to serve — fail closed.
	requireSandbox := envOr("EXECD_REQUIRE_SANDBOX", "0") == "1"
	if envOr("EXECD_SELFTEST", "1") == "1" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		rep, perr := runner.Probe(ctx, exec)
		cerr := runner.ProbeC(ctx, exec)
		cancel()
		if cerr != nil {
			log.Fatalf("execd: C toolchain self-test FAILED — refusing to start without required C execution: %v", cerr)
		}
		switch {
		case perr == nil && rep.Sandboxed:
			log.Printf("execd: sandbox self-test PASSED — network blocked, host FS hidden, bwrap mount active, C toolchain ready")
		case requireSandbox:
			log.Fatalf("execd: sandbox self-test FAILED and EXECD_REQUIRE_SANDBOX=1 — refusing to serve untrusted code unsandboxed. Reasons: %v (raw: %s)", rep.Reasons, rep.Raw)
		default:
			log.Printf("execd: WARNING — sandbox self-test did NOT confirm containment: %v. "+
				"Set EXECD_REQUIRE_SANDBOX=1 to make this fatal. Do not expose to the public internet like this.", rep.Reasons)
		}
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

	if !runner.SandboxAvailable() {
		log.Printf("execd: note — no bubblewrap on PATH; see the self-test result above for containment status")
	}
	log.Printf("execd listening on %s (max concurrency %d)", addr, maxConc)
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

// runSelftest runs the containment probe and exits with a status code so it can
// gate a container's readiness/health check: 0 = sandboxed, 1 = not.
func runSelftest() {
	exec, err := runner.New(runner.Config{Mode: "local", PythonPath: envOr("PYTHON", "python3"), CCPath: envOr("CC", "cc"), AllowUnsafe: true})
	if err != nil {
		log.Fatalf("selftest: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rep, err := runner.Probe(ctx, exec)
	cErr := runner.ProbeC(ctx, exec)
	out, _ := json.MarshalIndent(struct {
		Sandbox runner.SandboxReport `json:"sandbox"`
		CReady  bool                 `json:"c_ready"`
	}{Sandbox: rep, CReady: cErr == nil}, "", "  ")
	_, _ = os.Stdout.Write(append(out, '\n'))
	if err != nil || !rep.Sandboxed || cErr != nil {
		if cErr != nil {
			log.Printf("selftest: %v", cErr)
		}
		os.Exit(1)
	}
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

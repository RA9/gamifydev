// GamifyDev platform — Go server (server-rendered HTML + htmx + tan-compose,
// Turso/libSQL storage). Run: `go run .`  (uses a local SQLite file by default).
package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/RA9/gamifydev/platform/internal/email"
	"github.com/RA9/gamifydev/platform/internal/jobs"
	"github.com/RA9/gamifydev/platform/internal/rdb"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/seed"
	"github.com/RA9/gamifydev/platform/internal/server"
	"github.com/RA9/gamifydev/platform/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	loadDotEnv(".env")
	addr := envOr("ADDR", ":8080")
	// Railway (and most PaaS) inject the port to bind via $PORT.
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	// Local dev: a SQLite file. Production: set DATABASE_URL to your Turso URL,
	// e.g. libsql://<db>.turso.io?authToken=<token>
	dsn := envOr("DATABASE_URL", "gamifydev.db")
	// Secure cookies behind HTTPS. Railway always serves the public domain over
	// HTTPS, so enable them automatically there.
	secure := os.Getenv("SECURE_COOKIES") == "1" || os.Getenv("RAILWAY_ENVIRONMENT") != ""

	st, err := store.Open(dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Seed the embedded course content on a fresh database (or when SEED=1 is
	// set to force a re-seed). Non-fatal: a content issue shouldn't down the app.
	if os.Getenv("SEED") == "1" {
		if r, err := seed.Run(ctx, st); err != nil {
			log.Printf("seed (forced): %v", err)
		} else {
			log.Printf("seed (forced): %d courses, %d lessons, %d assignments, %d paths", r.Courses, r.Lessons, r.Assignments, r.Paths)
		}
	} else if r, seeded, err := seed.RunIfEmpty(ctx, st); err != nil {
		log.Printf("seed: %v", err)
	} else if seeded {
		log.Printf("seed: fresh database — %d courses, %d lessons, %d assignments, %d paths", r.Courses, r.Lessons, r.Assignments, r.Paths)
	}

	// Redis is optional — connect only if REDIS_URL is configured. Future
	// features (caching, rate limiting, live competition pub/sub) will use it.
	var rc *redis.Client
	if url := os.Getenv("REDIS_URL"); url != "" {
		rc, err = rdb.Open(ctx, url)
		if err != nil {
			log.Printf("redis: %v (continuing without it)", err)
		} else {
			log.Printf("redis: connected")
			defer rc.Close()
		}
	} else {
		log.Printf("redis: REDIS_URL not set, running without Redis")
	}

	// Email is optional. When SMTP is unconfigured, invite links are shown to the
	// admin instead of being emailed.
	mailer := email.New()
	if mailer.Configured() {
		log.Printf("email: SMTP configured")
	} else {
		log.Printf("email: SMTP not configured — invite links will be shown to the admin")
	}

	// Server-side Python execution (the Playground and future server-run labs).
	// Disabled unless CODE_EXEC is set; secure-by-default rules live in runner.New.
	exec, err := runner.New(runner.Config{
		Mode:         os.Getenv("CODE_EXEC"), // ""(off) | local | remote
		PythonPath:   os.Getenv("PYTHON"),
		SandboxURL:   os.Getenv("CODE_EXEC_SANDBOX_URL"),
		SandboxToken: os.Getenv("CODE_EXEC_SANDBOX_TOKEN"),
		Deployed:     os.Getenv("RAILWAY_ENVIRONMENT") != "",
		AllowUnsafe:  os.Getenv("CODE_EXEC_UNSAFE") == "1",
	})
	if err != nil {
		log.Fatalf("code exec: %v", err)
	}
	log.Printf("code exec: mode=%s enabled=%v", exec.Kind(), exec.Enabled())

	// If we're running code in-process (local mode), prove the sandbox actually
	// contains it before serving. runner.New already blocks local-in-prod
	// without bwrap; this also catches "bwrap is present but not isolating"
	// (e.g. user namespaces disabled). Remote mode self-tests inside execd.
	if exec.Enabled() && exec.Kind() == "local" {
		pctx, pcancel := context.WithTimeout(ctx, 15*time.Second)
		rep, _ := runner.Probe(pctx, exec)
		pcancel()
		deployed := os.Getenv("RAILWAY_ENVIRONMENT") != ""
		unsafe := os.Getenv("CODE_EXEC_UNSAFE") == "1"
		switch {
		case rep.Sandboxed:
			log.Printf("code exec: sandbox self-test PASSED")
		case deployed && !unsafe:
			log.Fatalf("code exec: sandbox self-test FAILED in a deployed environment — refusing to run untrusted code unsandboxed. Reasons: %v", rep.Reasons)
		default:
			log.Printf("code exec: WARNING — no verified sandbox (%v); fine for local dev, never for public traffic", rep.Reasons)
		}
	}

	// Scheduled work (standup windows, activity rollup, sanction expiry). The
	// runner takes a per-job database lease, so running more than one instance
	// is safe. Disable with JOBS=off for one-shot/CLI-style runs.
	runner := jobs.New(st, os.Getenv("INSTANCE_ID"))
	if os.Getenv("JOBS") == "off" {
		log.Printf("jobs: disabled (JOBS=off)")
		runner = nil
	} else {
		jobs.Register(runner, st)
	}

	srv, err := server.New(st, rc, mailer, exec, runner, secure)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	// One context governs background work; cancelled on SIGINT/SIGTERM so the
	// runner finishes its in-flight job instead of being killed mid-write.
	bgCtx, stopBG := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopBG()
	if runner != nil {
		runner.Start(bgCtx)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Serve until a signal arrives, then drain.
	serveErr := make(chan error, 1)
	go func() {
		log.Printf("GamifyDev platform listening on %s (db: %s)", addr, dsn)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		log.Fatal(err)
	case <-bgCtx.Done():
		log.Printf("shutting down…")
	}

	shutCtx, cancelShut := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShut()
	if err := httpServer.Shutdown(shutCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	if runner != nil {
		runner.Stop() // waits for the current job to finish
	}
	log.Printf("stopped")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// loadDotEnv reads simple KEY=VALUE lines from path (if present) into the
// process environment. Real environment variables always win — a var already
// set (e.g. by the shell, systemd, or a PaaS) is never overwritten — so .env
// is purely a local-dev convenience, not a config layer that can mask deploy
// config. Blank lines and lines starting with # are ignored; values may be
// wrapped in matching single or double quotes.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if key == "" {
			continue
		}
		if _, set := os.LookupEnv(key); !set {
			os.Setenv(key, val)
		}
	}
}

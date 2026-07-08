// GamifyDev platform — Go server (server-rendered HTML + htmx + tan-compose,
// Turso/libSQL storage). Run: `go run .`  (uses a local SQLite file by default).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RA9/gamifydev/platform/internal/email"
	"github.com/RA9/gamifydev/platform/internal/rdb"
	"github.com/RA9/gamifydev/platform/internal/runner"
	"github.com/RA9/gamifydev/platform/internal/seed"
	"github.com/RA9/gamifydev/platform/internal/server"
	"github.com/RA9/gamifydev/platform/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
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

	srv, err := server.New(st, rc, mailer, exec, secure)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("GamifyDev platform listening on %s (db: %s)", addr, dsn)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

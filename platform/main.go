// GamifyDev platform — Go server (server-rendered HTML + htmx + tan-compose,
// Turso/libSQL storage). Run: `go run .`  (uses a local SQLite file by default).
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RA9/gamifydev/platform/internal/server"
	"github.com/RA9/gamifydev/platform/internal/store"
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

	srv, err := server.New(st, secure)
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

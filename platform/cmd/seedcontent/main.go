// Command seedcontent (idempotently) seeds the embedded course content into a
// database. The app also auto-seeds a fresh DB on boot, so this is mainly for
// manually (re)seeding a specific database, e.g. production:
//
//	go run ./cmd/seedcontent                                  # local gamifydev.db
//	DATABASE_URL=libsql://...?authToken=... go run ./cmd/seedcontent   # Turso
package main

import (
	"context"
	"log"
	"os"

	"github.com/RA9/gamifydev/platform/internal/seed"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "gamifydev.db"
	}
	st, err := store.Open(dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer st.Close()

	ctx := context.Background()
	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	r, err := seed.Run(ctx, st)
	if err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Printf("seeded %d courses, %d lessons, %d assignments", r.Courses, r.Lessons, r.Assignments)
}

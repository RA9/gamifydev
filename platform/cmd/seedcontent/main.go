// Command seedcontent idempotently synchronizes the embedded course content
// into a database. The server performs the same synchronization on every boot;
// this command is useful for running it without starting the HTTP server:
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
	log.Printf("seeded %d courses, %d lessons, %d assignments, %d paths, %d diagnostic items, %d problems",
		r.Courses, r.Lessons, r.Assignments, r.Paths, r.Items, r.Problems)
}

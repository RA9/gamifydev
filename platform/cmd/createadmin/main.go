// Command createadmin creates or promotes the operator account without exposing
// an admin bootstrap path over HTTP.
//
// Usage:
//
//	ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD='...' go run ./cmd/createadmin
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func main() {
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := os.Getenv("ADMIN_PASSWORD")
	name := strings.TrimSpace(os.Getenv("ADMIN_NAME"))
	if name == "" {
		name = "GamifyDev Admin"
	}
	if !strings.Contains(email, "@") || len(password) < 12 {
		log.Fatal("ADMIN_EMAIL must be valid and ADMIN_PASSWORD must be at least 12 characters")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "gamifydev.db"
	}
	st, err := store.Open(dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer st.Close()
	ctx := context.Background()
	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	u, err := st.GetUserByEmail(ctx, email)
	if err == nil {
		if err := st.SetUserRole(ctx, u.ID, "admin"); err != nil {
			log.Fatalf("promote account: %v", err)
		}
		log.Printf("promoted %s to admin", email)
		return
	}
	if !errors.Is(err, store.ErrNotFound) {
		log.Fatalf("find account: %v", err)
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	if _, err := st.CreateUser(ctx, email, hash, name, "admin"); err != nil {
		log.Fatalf("create admin: %v", err)
	}
	log.Printf("created admin account %s", email)
}

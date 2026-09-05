package seed

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

func TestEnsureFirstAdminRequiresStrongPasswordForCreation(t *testing.T) {
	for _, password := range []string{"", "too-short"} {
		t.Run(password, func(t *testing.T) {
			ctx, st := newAdminTestStore(t)
			created, err := EnsureFirstAdmin(ctx, st, password)
			if !errors.Is(err, ErrAdminPasswordRequired) {
				t.Fatalf("error = %v, want ErrAdminPasswordRequired", err)
			}
			if created {
				t.Fatal("created = true after rejected password")
			}
			if _, err := st.GetUserByEmail(ctx, FirstAdminEmail); !errors.Is(err, store.ErrNotFound) {
				t.Fatalf("lookup after rejection = %v, want store.ErrNotFound", err)
			}
		})
	}
}

func TestEnsureFirstAdminIsIdempotentAndPreservesPassword(t *testing.T) {
	ctx, st := newAdminTestStore(t)
	const originalPassword = "a-strong-original-password"

	created, err := EnsureFirstAdmin(ctx, st, originalPassword)
	if err != nil {
		t.Fatalf("first ensure: %v", err)
	}
	if !created {
		t.Fatal("first ensure did not report creation")
	}

	first, err := st.GetUserByEmail(ctx, FirstAdminEmail)
	if err != nil {
		t.Fatalf("get created admin: %v", err)
	}
	if !first.IsAdmin() {
		t.Fatalf("created role = %q, want admin", first.Role)
	}
	if !auth.CheckPassword(first.PasswordHash, originalPassword) {
		t.Fatal("created password does not match original password")
	}

	created, err = EnsureFirstAdmin(ctx, st, "a-different-strong-password")
	if err != nil {
		t.Fatalf("second ensure: %v", err)
	}
	if created {
		t.Fatal("second ensure reported creating a duplicate")
	}

	users, err := st.ListUsers(ctx)
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("user count = %d, want 1", len(users))
	}
	if users[0].ID != first.ID {
		t.Fatalf("user ID changed from %d to %d", first.ID, users[0].ID)
	}
	if users[0].PasswordHash != first.PasswordHash {
		t.Fatal("second ensure replaced the existing password hash")
	}
	if !auth.CheckPassword(users[0].PasswordHash, originalPassword) {
		t.Fatal("original password stopped working after second ensure")
	}
	if auth.CheckPassword(users[0].PasswordHash, "a-different-strong-password") {
		t.Fatal("second ensure reset the password")
	}
}

func TestEnsureFirstAdminPromotesExistingUserWithoutPassword(t *testing.T) {
	ctx, st := newAdminTestStore(t)
	const learnerPassword = "existing-learner-password"
	hash, err := auth.HashPassword(learnerPassword)
	if err != nil {
		t.Fatalf("hash learner password: %v", err)
	}
	learner, err := st.CreateUser(ctx, FirstAdminEmail, hash, "Existing Learner", "learner")
	if err != nil {
		t.Fatalf("create learner: %v", err)
	}

	created, err := EnsureFirstAdmin(ctx, st, "")
	if err != nil {
		t.Fatalf("ensure existing user: %v", err)
	}
	if created {
		t.Fatal("ensure reported creating an existing user")
	}

	promoted, err := st.GetUserByEmail(ctx, FirstAdminEmail)
	if err != nil {
		t.Fatalf("get promoted user: %v", err)
	}
	if promoted.ID != learner.ID {
		t.Fatalf("promoted user ID = %d, want %d", promoted.ID, learner.ID)
	}
	if !promoted.IsAdmin() {
		t.Fatalf("promoted role = %q, want admin", promoted.Role)
	}
	if promoted.PasswordHash != hash || !auth.CheckPassword(promoted.PasswordHash, learnerPassword) {
		t.Fatal("promotion changed the existing password")
	}
}

func newAdminTestStore(t *testing.T) (context.Context, *store.Store) {
	t.Helper()
	ctx := context.Background()
	st, err := store.Open(filepath.Join(t.TempDir(), "admin.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return ctx, st
}

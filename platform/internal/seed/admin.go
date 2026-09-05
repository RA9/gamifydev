package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/RA9/gamifydev/platform/internal/auth"
	"github.com/RA9/gamifydev/platform/internal/store"
)

const (
	FirstAdminEmail = "cnah27@gmail.com"
	firstAdminName  = "GamifyDev Admin"
)

var ErrAdminPasswordRequired = errors.New("ADMIN_PASSWORD must be at least 12 characters to create the first administrator")

// EnsureFirstAdmin creates the fixed bootstrap administrator when it does not
// exist, or promotes the existing account without changing its password. It is
// safe to call on every process start.
func EnsureFirstAdmin(ctx context.Context, st *store.Store, password string) (created bool, err error) {
	u, err := st.GetUserByEmail(ctx, FirstAdminEmail)
	if err == nil {
		if !u.IsAdmin() {
			if err := st.SetUserRole(ctx, u.ID, "admin"); err != nil {
				return false, fmt.Errorf("promote first administrator: %w", err)
			}
		}
		return false, nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return false, fmt.Errorf("find first administrator: %w", err)
	}
	if len(password) < 12 {
		return false, ErrAdminPasswordRequired
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return false, fmt.Errorf("hash first administrator password: %w", err)
	}
	if _, err := st.CreateUser(ctx, FirstAdminEmail, hash, firstAdminName, "admin"); err == nil {
		return true, nil
	} else {
		// Another process may have created the account after our initial lookup.
		// Resolve that race by treating the now-existing account exactly as a
		// subsequent boot rather than attempting another insert.
		u, lookupErr := st.GetUserByEmail(ctx, FirstAdminEmail)
		if lookupErr != nil {
			return false, fmt.Errorf("create first administrator: %w", err)
		}
		if !u.IsAdmin() {
			if promoteErr := st.SetUserRole(ctx, u.ID, "admin"); promoteErr != nil {
				return false, fmt.Errorf("promote concurrently created first administrator: %w", promoteErr)
			}
		}
		return false, nil
	}
}

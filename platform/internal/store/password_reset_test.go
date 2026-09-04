package store

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestResetTokenRoundTrip(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")

	if err := st.CreatePasswordReset(ctx, uid, "tok-abc", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("create: %v", err)
	}
	u, err := st.UserByResetToken(ctx, "tok-abc")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if u.ID != uid {
		t.Fatalf("token resolved to user %d, want %d", u.ID, uid)
	}
}

func TestTheRawTokenIsNeverStored(t *testing.T) {
	// The whole point of hashing: a dump of this table must not yield anything
	// that can be replayed as a reset link.
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	const token = "sekrit-token-value"

	if err := st.CreatePasswordReset(ctx, uid, token, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("create: %v", err)
	}
	var stored string
	if err := st.db.QueryRow(`SELECT token_hash FROM password_resets WHERE user_id = ?`, uid).Scan(&stored); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if strings.Contains(stored, token) || stored == token {
		t.Fatalf("the raw token is in the database: %q", stored)
	}
	if stored != HashResetToken(token) {
		t.Fatalf("stored value is not the token's hash")
	}
}

func TestLookupDoesNotConsumeTheToken(t *testing.T) {
	// Opening the link and reloading the page must not spend it.
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	_ = st.CreatePasswordReset(ctx, uid, "tok", time.Now().Add(time.Hour))

	for i := 0; i < 3; i++ {
		if _, err := st.UserByResetToken(ctx, "tok"); err != nil {
			t.Fatalf("lookup %d: %v", i+1, err)
		}
	}
	if _, err := st.ResetPassword(ctx, "tok", "newhash"); err != nil {
		t.Fatalf("reset after lookups: %v", err)
	}
}

func TestATokenWorksExactlyOnce(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	_ = st.CreatePasswordReset(ctx, uid, "tok", time.Now().Add(time.Hour))

	if _, err := st.ResetPassword(ctx, "tok", "hash-one"); err != nil {
		t.Fatalf("first use: %v", err)
	}
	_, err := st.ResetPassword(ctx, "tok", "hash-two")
	if !errors.Is(err, ErrResetSpent) {
		t.Fatalf("second use returned %v, want ErrResetSpent", err)
	}
	var hash string
	_ = st.db.QueryRow(`SELECT password_hash FROM users WHERE id = ?`, uid).Scan(&hash)
	if hash != "hash-one" {
		t.Fatalf("password = %q; the replayed token changed it", hash)
	}
}

func TestAnExpiredTokenIsRefused(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	_ = st.CreatePasswordReset(ctx, uid, "stale", time.Now().Add(-time.Minute))

	if _, err := st.UserByResetToken(ctx, "stale"); err == nil {
		t.Fatal("an expired token resolved to a user")
	}
	if _, err := st.ResetPassword(ctx, "stale", "nope"); !errors.Is(err, ErrResetSpent) {
		t.Fatalf("expired token returned %v, want ErrResetSpent", err)
	}
}

func TestRequestingAgainInvalidatesTheEarlierLink(t *testing.T) {
	// Three clicks on "forgot password" must not leave three working keys to
	// the account sitting in an inbox.
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")

	_ = st.CreatePasswordReset(ctx, uid, "first", time.Now().Add(time.Hour))
	_ = st.CreatePasswordReset(ctx, uid, "second", time.Now().Add(time.Hour))

	if _, err := st.ResetPassword(ctx, "first", "hash"); !errors.Is(err, ErrResetSpent) {
		t.Fatalf("the superseded link still worked (%v)", err)
	}
	if _, err := st.ResetPassword(ctx, "second", "hash"); err != nil {
		t.Fatalf("the newest link failed: %v", err)
	}
}

func TestResetRevokesEverySession(t *testing.T) {
	// Resetting is what someone does when their account was taken. Leaving the
	// intruder's session alive would hand them the account back.
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	other := mkUser(t, st, "someone-else@example.com")

	exp := time.Now().Add(time.Hour)
	_ = st.CreateSession(ctx, "victim-session", uid, exp)
	_ = st.CreateSession(ctx, "intruder-session", uid, exp)
	_ = st.CreateSession(ctx, "bystander-session", other, exp)
	_ = st.CreatePasswordReset(ctx, uid, "tok", exp)

	if _, err := st.ResetPassword(ctx, "tok", "newhash"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	for _, tok := range []string{"victim-session", "intruder-session"} {
		if _, err := st.UserBySession(ctx, tok); err == nil {
			t.Errorf("session %q survived the reset", tok)
		}
	}
	// ...but only that user's sessions.
	if _, err := st.UserBySession(ctx, "bystander-session"); err != nil {
		t.Errorf("an unrelated user was signed out: %v", err)
	}
}

func TestAnUnknownTokenIsRefused(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	if _, err := st.ResetPassword(ctx, "never-issued", "hash"); !errors.Is(err, ErrResetSpent) {
		t.Fatalf("unknown token returned %v, want ErrResetSpent", err)
	}
}

func TestPruneDropsOnlyDeadRows(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	uid := mkUser(t, st, "learner@example.com")
	live := mkUser(t, st, "live@example.com")

	_ = st.CreatePasswordReset(ctx, uid, "old", time.Now().Add(-48*time.Hour))
	_ = st.CreatePasswordReset(ctx, live, "fresh", time.Now().Add(time.Hour))

	n, err := st.PruneExpiredResets(ctx, 24*time.Hour)
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if n != 1 {
		t.Fatalf("pruned %d rows, want 1", n)
	}
	if _, err := st.UserByResetToken(ctx, "fresh"); err != nil {
		t.Fatalf("prune removed a live token: %v", err)
	}
}

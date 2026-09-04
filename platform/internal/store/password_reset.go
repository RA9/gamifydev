package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

// HashResetToken is the one-way transform between the token a learner receives
// and the row that records it.
//
// SHA-256 rather than bcrypt deliberately: the token is 32 bytes of CSPRNG
// output, so there is no low-entropy secret to slow an attacker down over —
// the only job here is that the stored value can't be replayed as a token.
func HashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreatePasswordReset records a reset request and returns nothing — the caller
// already holds the token, and it is never readable again from here.
//
// Any earlier unused token for the same user is expired first. Otherwise a
// learner who clicks "forgot password" three times leaves three working keys to
// their account lying in three inboxes.
func (s *Store) CreatePasswordReset(ctx context.Context, userID int64, token string, expires time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx,
		`UPDATE password_resets SET used_at = datetime('now')
		 WHERE user_id = ? AND used_at IS NULL`, userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO password_resets (user_id, token_hash, expires_at) VALUES (?, ?, ?)`,
		userID, HashResetToken(token), expires.UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	return tx.Commit()
}

// UserByResetToken returns the user a live token belongs to, so the reset form
// can be shown before the new password is chosen. It does not consume the
// token: a learner who opens the link and then reloads the page must not find
// it already spent.
func (s *Store) UserByResetToken(ctx context.Context, token string) (*User, error) {
	return s.scanUser(s.db.QueryRowContext(ctx, `
		SELECT `+userColsPrefixed+`
		FROM password_resets pr JOIN users u ON u.id = pr.user_id
		WHERE pr.token_hash = ? AND pr.used_at IS NULL AND pr.expires_at > ?`,
		HashResetToken(token), time.Now().UTC().Format(time.RFC3339)))
}

// ErrResetSpent means the token was already used, expired, or never existed.
// One error for all three: the difference is not something the person holding
// the link can act on differently, and telling them would confirm which valid
// tokens exist.
var ErrResetSpent = errors.New("store: reset link is no longer valid")

// ResetPassword consumes a token and sets the new password in one transaction,
// then deletes every session the user has.
//
// The three belong together. Signing out everywhere is the point of a reset for
// anyone whose account was taken: leaving the intruder's session alive would
// hand them the new password's account too. And consuming inside the same
// transaction is what stops two submissions of the same link from both landing.
func (s *Store) ResetPassword(ctx context.Context, token, passwordHash string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck

	var id, userID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id, user_id FROM password_resets
		WHERE token_hash = ? AND used_at IS NULL AND expires_at > ?`,
		HashResetToken(token), time.Now().UTC().Format(time.RFC3339)).Scan(&id, &userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrResetSpent
	} else if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE password_resets SET used_at = datetime('now') WHERE id = ?`, id); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET password_hash = ?, updated_at = datetime('now') WHERE id = ?`,
		passwordHash, userID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM sessions WHERE user_id = ?`, userID); err != nil {
		return 0, err
	}
	return userID, tx.Commit()
}

// PruneExpiredResets drops reset rows that are long dead, so the table doesn't
// grow without bound on a public signup form.
func (s *Store) PruneExpiredResets(ctx context.Context, olderThan time.Duration) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM password_resets WHERE expires_at < ?`,
		time.Now().Add(-olderThan).UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// Invitation is a pending or accepted invite for someone to join with a role.
type Invitation struct {
	ID          int64
	Email       string
	Name        string
	Role        string
	Token       string
	InvitedBy   sql.NullInt64
	AcceptedAt  sql.NullString
	ExpiresAt   string
	CreatedAt   string
	InviterName string // joined from users (empty if unknown)
}

// Accepted reports whether the invitation has already been used.
func (iv Invitation) Accepted() bool { return iv.AcceptedAt.Valid && iv.AcceptedAt.String != "" }

// Expired reports whether the invitation is past its expiry.
func (iv Invitation) Expired() bool {
	t, err := time.Parse(time.RFC3339, iv.ExpiresAt)
	return err == nil && time.Now().After(t)
}

type rowScanner interface{ Scan(dest ...any) error }

const invitationSelect = `
SELECT i.id, i.email, i.name, i.role, i.token, i.invited_by, i.accepted_at, i.expires_at, i.created_at,
       COALESCE(u.name, '')
FROM invitations i
LEFT JOIN users u ON u.id = i.invited_by`

func scanInvitation(sc rowScanner) (*Invitation, error) {
	var iv Invitation
	err := sc.Scan(&iv.ID, &iv.Email, &iv.Name, &iv.Role, &iv.Token, &iv.InvitedBy,
		&iv.AcceptedAt, &iv.ExpiresAt, &iv.CreatedAt, &iv.InviterName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &iv, nil
}

// CreateInvitation records an invite and returns it.
func (s *Store) CreateInvitation(ctx context.Context, email, name, role, token string, invitedBy int64, expires time.Time) (*Invitation, error) {
	var inviter sql.NullInt64
	if invitedBy > 0 {
		inviter = sql.NullInt64{Int64: invitedBy, Valid: true}
	}
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO invitations (email, name, role, token, invited_by, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?) RETURNING id`,
		strings.ToLower(strings.TrimSpace(email)), name, role, token, inviter,
		expires.UTC().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		return nil, err
	}
	return s.GetInvitationByID(ctx, id)
}

func (s *Store) GetInvitationByID(ctx context.Context, id int64) (*Invitation, error) {
	return scanInvitation(s.db.QueryRowContext(ctx, invitationSelect+` WHERE i.id = ?`, id))
}

func (s *Store) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	return scanInvitation(s.db.QueryRowContext(ctx, invitationSelect+` WHERE i.token = ?`, token))
}

// ListPendingInvitations returns unaccepted invitations, newest first.
func (s *Store) ListPendingInvitations(ctx context.Context) ([]Invitation, error) {
	rows, err := s.db.QueryContext(ctx, invitationSelect+` WHERE i.accepted_at IS NULL ORDER BY i.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Invitation
	for rows.Next() {
		iv, err := scanInvitation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *iv)
	}
	return out, rows.Err()
}

// AcceptInvitation marks an invitation as used.
func (s *Store) AcceptInvitation(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE invitations SET accepted_at = datetime('now') WHERE token = ?`, token)
	return err
}

// DeleteInvitation removes an invitation (revoke).
func (s *Store) DeleteInvitation(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM invitations WHERE id = ?`, id)
	return err
}

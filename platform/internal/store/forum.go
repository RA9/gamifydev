package store

import (
	"context"
	"database/sql"
	"errors"
)

type Thread struct {
	ID           int64
	UserID       int64
	AuthorName   string
	AuthorRole   string
	Title        string
	Body         string
	Category     string
	Pinned       bool
	Locked       bool
	CreatedAt    string
	ReplyCount   int
	LastActivity string
}

type Reply struct {
	ID         int64
	ThreadID   int64
	UserID     int64
	AuthorName string
	AuthorRole string
	Body       string
	CreatedAt  string
}

func (s *Store) CreateThread(ctx context.Context, userID int64, title, body, category string) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO forum_threads (user_id, title, body, category) VALUES (?, ?, ?, ?)`,
		userID, title, body, category)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ListThreads returns threads with author, reply count, and last-activity time,
// pinned first then most-recently-active.
func (s *Store) ListThreads(ctx context.Context) ([]Thread, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.id, t.user_id, u.name, u.role, t.title, t.category, t.pinned, t.locked, t.created_at,
		       (SELECT COUNT(*) FROM forum_replies r WHERE r.thread_id = t.id) AS reply_count,
		       COALESCE((SELECT MAX(r.created_at) FROM forum_replies r WHERE r.thread_id = t.id), t.created_at) AS last_activity
		FROM forum_threads t JOIN users u ON u.id = t.user_id
		ORDER BY t.pinned DESC, last_activity DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Thread
	for rows.Next() {
		var t Thread
		var pinned, locked int
		if err := rows.Scan(&t.ID, &t.UserID, &t.AuthorName, &t.AuthorRole, &t.Title, &t.Category,
			&pinned, &locked, &t.CreatedAt, &t.ReplyCount, &t.LastActivity); err != nil {
			return nil, err
		}
		t.Pinned, t.Locked = pinned == 1, locked == 1
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) GetThread(ctx context.Context, id int64) (*Thread, error) {
	var t Thread
	var pinned, locked int
	err := s.db.QueryRowContext(ctx, `
		SELECT t.id, t.user_id, u.name, u.role, t.title, t.body, t.category, t.pinned, t.locked, t.created_at
		FROM forum_threads t JOIN users u ON u.id = t.user_id WHERE t.id = ?`, id).
		Scan(&t.ID, &t.UserID, &t.AuthorName, &t.AuthorRole, &t.Title, &t.Body, &t.Category, &pinned, &locked, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	t.Pinned, t.Locked = pinned == 1, locked == 1
	return &t, nil
}

func (s *Store) ListReplies(ctx context.Context, threadID int64) ([]Reply, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT r.id, r.thread_id, r.user_id, u.name, u.role, r.body, r.created_at
		FROM forum_replies r JOIN users u ON u.id = r.user_id
		WHERE r.thread_id = ? ORDER BY r.created_at`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Reply
	for rows.Next() {
		var r Reply
		if err := rows.Scan(&r.ID, &r.ThreadID, &r.UserID, &r.AuthorName, &r.AuthorRole, &r.Body, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) CreateReply(ctx context.Context, threadID, userID int64, body string) error {
	if _, err := s.db.ExecContext(ctx,
		`INSERT INTO forum_replies (thread_id, user_id, body) VALUES (?, ?, ?)`, threadID, userID, body); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `UPDATE forum_threads SET updated_at = datetime('now') WHERE id = ?`, threadID)
	return err
}

func (s *Store) CountThreads(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM forum_threads`).Scan(&n)
	return n, err
}

// --- Moderation -------------------------------------------------------------

func (s *Store) SetThreadPinned(ctx context.Context, id int64, pinned bool) error {
	_, err := s.db.ExecContext(ctx, `UPDATE forum_threads SET pinned = ? WHERE id = ?`, boolToInt(pinned), id)
	return err
}

func (s *Store) SetThreadLocked(ctx context.Context, id int64, locked bool) error {
	_, err := s.db.ExecContext(ctx, `UPDATE forum_threads SET locked = ? WHERE id = ?`, boolToInt(locked), id)
	return err
}

func (s *Store) DeleteThread(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM forum_threads WHERE id = ?`, id)
	return err
}

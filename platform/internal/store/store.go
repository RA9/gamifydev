// Package store is the data layer: it opens the database (local SQLite file in
// dev, Turso/libSQL in production — both speak the same SQL), runs embedded
// migrations, and exposes typed query methods. No ORM; just database/sql.
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql" // driver "libsql" (remote Turso)
	_ "modernc.org/sqlite"                               // driver "sqlite" (local file, pure Go)
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var ErrNotFound = errors.New("not found")

type Store struct {
	db *sql.DB
}

// Open connects using a DSN. A libsql:// or http(s):// DSN uses the Turso
// driver; anything else is treated as a local SQLite file path.
func Open(dsn string) (*Store, error) {
	driver := "sqlite"
	conn := dsn
	if strings.HasPrefix(dsn, "libsql://") || strings.HasPrefix(dsn, "http://") || strings.HasPrefix(dsn, "https://") {
		driver = "libsql"
	} else {
		// Local file: enable foreign keys + WAL for sane concurrent reads.
		if !strings.Contains(conn, "?") {
			conn += "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
		}
	}
	db, err := sql.Open(driver, conn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Migrate applies any embedded migrations not yet recorded, in filename order.
func (s *Store) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TEXT NOT NULL DEFAULT (datetime('now')))`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		var exists int
		_ = s.db.QueryRowContext(ctx, `SELECT 1 FROM schema_migrations WHERE name = ?`, name).Scan(&exists)
		if exists == 1 {
			continue
		}
		body, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		if _, err := s.db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if _, err := s.db.ExecContext(ctx, `INSERT INTO schema_migrations (name) VALUES (?)`, name); err != nil {
			return err
		}
	}
	return nil
}

// --- Models -----------------------------------------------------------------

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	Role         string
	Bio          string
	AvatarURL    string
	CreatedAt    string
}

func (u User) IsAdmin() bool  { return u.Role == "admin" }
func (u User) CanGrade() bool { return u.Role == "admin" || u.Role == "grader" }

type Session struct {
	Token     string
	UserID    int64
	ExpiresAt time.Time
}

// --- User queries -----------------------------------------------------------

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, name, role string) (*User, error) {
	// Use INSERT ... RETURNING to get the new row in a single round-trip. This
	// avoids relying on LastInsertId(), which the libsql/Turso HTTP driver does
	// not reliably populate (it can return 0, making a follow-up lookup fail
	// even though the INSERT succeeded).
	return s.scanUser(s.db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, name, role) VALUES (?, ?, ?, ?)
		 RETURNING `+userCols,
		strings.ToLower(strings.TrimSpace(email)), passwordHash, name, role))
}

func (s *Store) scanUser(row *sql.Row) (*User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.Bio, &u.AvatarURL, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

const userCols = `id, email, password_hash, name, role, bio, avatar_url, created_at`

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.scanUser(s.db.QueryRowContext(ctx,
		`SELECT `+userCols+` FROM users WHERE email = ?`, strings.ToLower(strings.TrimSpace(email))))
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return s.scanUser(s.db.QueryRowContext(ctx, `SELECT `+userCols+` FROM users WHERE id = ?`, id))
}

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+userCols+` FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.Bio, &u.AvatarURL, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (s *Store) SetUserRole(ctx context.Context, id int64, role string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET role = ?, updated_at = datetime('now') WHERE id = ?`, role, id)
	return err
}

// --- Session queries --------------------------------------------------------

func (s *Store) CreateSession(ctx context.Context, token string, userID int64, expires time.Time) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expires.UTC().Format(time.RFC3339))
	return err
}

// UserBySession returns the user for a valid, unexpired session token.
func (s *Store) UserBySession(ctx context.Context, token string) (*User, error) {
	var u User
	var expires string
	err := s.db.QueryRowContext(ctx,
		`SELECT u.id, u.email, u.password_hash, u.name, u.role, u.bio, u.avatar_url, u.created_at, s.expires_at
		 FROM sessions s JOIN users u ON u.id = s.user_id WHERE s.token = ?`, token).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.Bio, &u.AvatarURL, &u.CreatedAt, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t, perr := time.Parse(time.RFC3339, expires); perr == nil && time.Now().After(t) {
		_ = s.DeleteSession(ctx, token)
		return nil, ErrNotFound
	}
	return &u, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = ?`, token)
	return err
}

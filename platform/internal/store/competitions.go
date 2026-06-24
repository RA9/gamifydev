package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Competition struct {
	ID         int64
	Slug       string
	Title      string
	Prompt     string
	Language   string
	Points     int
	StartsAt   time.Time
	EndsAt     time.Time
	Published  bool
	EntryCount int
}

// Status is "upcoming", "live", or "ended" relative to now.
func (c Competition) Status() string {
	now := time.Now()
	switch {
	case now.Before(c.StartsAt):
		return "upcoming"
	case now.After(c.EndsAt):
		return "ended"
	default:
		return "live"
	}
}

func (c Competition) IsLive() bool { return c.Status() == "live" }

type Entry struct {
	ID        int64
	UserID    int64
	UserName  string
	Code      string
	Score     int
	CreatedAt string
	UpdatedAt string
}

func parseTime(s string) time.Time {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func (s *Store) CreateCompetition(ctx context.Context, c Competition) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO competitions (slug, title, prompt, language, points, starts_at, ends_at, published)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Slug, c.Title, c.Prompt, c.Language, c.Points,
		c.StartsAt.UTC().Format(time.RFC3339), c.EndsAt.UTC().Format(time.RFC3339), boolToInt(c.Published))
	return err
}

const compCols = `c.id, c.slug, c.title, c.prompt, c.language, c.points, c.starts_at, c.ends_at, c.published`

func scanComp(sc interface{ Scan(...any) error }) (*Competition, error) {
	var c Competition
	var starts, ends string
	var pub int
	if err := sc.Scan(&c.ID, &c.Slug, &c.Title, &c.Prompt, &c.Language, &c.Points, &starts, &ends, &pub); err != nil {
		return nil, err
	}
	c.StartsAt, c.EndsAt, c.Published = parseTime(starts), parseTime(ends), pub == 1
	return &c, nil
}

func (s *Store) ListCompetitions(ctx context.Context, includeUnpublished bool) ([]Competition, error) {
	q := `SELECT ` + compCols + `,
		   (SELECT COUNT(*) FROM competition_entries e WHERE e.competition_id = c.id) AS entry_count
		   FROM competitions c`
	if !includeUnpublished {
		q += ` WHERE c.published = 1`
	}
	q += ` ORDER BY c.starts_at DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Competition
	for rows.Next() {
		var c Competition
		var starts, ends string
		var pub int
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Prompt, &c.Language, &c.Points, &starts, &ends, &pub, &c.EntryCount); err != nil {
			return nil, err
		}
		c.StartsAt, c.EndsAt, c.Published = parseTime(starts), parseTime(ends), pub == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetCompetitionBySlug(ctx context.Context, slug string) (*Competition, error) {
	c, err := scanComp(s.db.QueryRowContext(ctx, `SELECT `+compCols+` FROM competitions c WHERE c.slug = ?`, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

// UpsertEntry creates or updates a learner's entry for a competition.
func (s *Store) UpsertEntry(ctx context.Context, competitionID, userID int64, code string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO competition_entries (competition_id, user_id, code, updated_at)
		VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT(competition_id, user_id) DO UPDATE SET code = excluded.code, updated_at = datetime('now')`,
		competitionID, userID, code)
	return err
}

func (s *Store) GetEntry(ctx context.Context, competitionID, userID int64) (*Entry, error) {
	var e Entry
	err := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, code, score, created_at, updated_at FROM competition_entries
		 WHERE competition_id = ? AND user_id = ?`, competitionID, userID).
		Scan(&e.ID, &e.UserID, &e.Code, &e.Score, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &e, err
}

// Leaderboard ranks entries by score (desc) then earliest submission.
func (s *Store) Leaderboard(ctx context.Context, competitionID int64) ([]Entry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT e.id, e.user_id, u.name, e.score, e.created_at, e.updated_at
		FROM competition_entries e JOIN users u ON u.id = e.user_id
		WHERE e.competition_id = ?
		ORDER BY e.score DESC, e.created_at ASC`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.UserID, &e.UserName, &e.Score, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// ListEntriesWithCode returns all entries (with code) for staff scoring.
func (s *Store) ListEntriesWithCode(ctx context.Context, competitionID int64) ([]Entry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT e.id, e.user_id, u.name, e.code, e.score, e.created_at
		FROM competition_entries e JOIN users u ON u.id = e.user_id
		WHERE e.competition_id = ? ORDER BY e.score DESC, e.created_at ASC`, competitionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.UserID, &e.UserName, &e.Code, &e.Score, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) SetEntryScore(ctx context.Context, entryID int64, score int) error {
	_, err := s.db.ExecContext(ctx, `UPDATE competition_entries SET score = ? WHERE id = ?`, score, entryID)
	return err
}

func (s *Store) CountCompetitions(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM competitions`).Scan(&n)
	return n, err
}

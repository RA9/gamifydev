package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	ID          int64
	Slug        string
	Title       string
	Excerpt     string
	Body        string
	CoverURL    string
	AuthorID    sql.NullInt64
	AuthorName  string
	Published   bool
	PublishedAt string
	CreatedAt   string
}

const postCols = `p.id, p.slug, p.title, p.excerpt, p.body, p.cover_url, p.author_id, p.published,
	COALESCE(p.published_at,''), p.created_at, COALESCE(u.name,'')`

func scanPost(sc interface{ Scan(...any) error }) (*Post, error) {
	var p Post
	var pub int
	if err := sc.Scan(&p.ID, &p.Slug, &p.Title, &p.Excerpt, &p.Body, &p.CoverURL, &p.AuthorID,
		&pub, &p.PublishedAt, &p.CreatedAt, &p.AuthorName); err != nil {
		return nil, err
	}
	p.Published = pub == 1
	return &p, nil
}

// ListPosts returns posts ordered newest-first. When publishedOnly, drafts are
// excluded and ordering uses the publish date.
func (s *Store) ListPosts(ctx context.Context, publishedOnly bool) ([]Post, error) {
	q := `SELECT ` + postCols + ` FROM blog_posts p LEFT JOIN users u ON u.id = p.author_id`
	if publishedOnly {
		q += ` WHERE p.published = 1 ORDER BY COALESCE(p.published_at, p.created_at) DESC`
	} else {
		q += ` ORDER BY p.created_at DESC`
	}
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, rows.Err()
}

// GetPublishedPostBySlug returns a published post for the public site.
func (s *Store) GetPublishedPostBySlug(ctx context.Context, slug string) (*Post, error) {
	p, err := scanPost(s.db.QueryRowContext(ctx,
		`SELECT `+postCols+` FROM blog_posts p LEFT JOIN users u ON u.id = p.author_id
		 WHERE p.slug = ? AND p.published = 1`, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *Store) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	p, err := scanPost(s.db.QueryRowContext(ctx,
		`SELECT `+postCols+` FROM blog_posts p LEFT JOIN users u ON u.id = p.author_id WHERE p.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *Store) CreatePost(ctx context.Context, p Post) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO blog_posts (slug, title, excerpt, body, cover_url, author_id, published, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Slug, p.Title, p.Excerpt, p.Body, p.CoverURL, p.AuthorID, boolToInt(p.Published), nullStr(p.PublishedAt))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdatePost(ctx context.Context, id int64, p Post) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE blog_posts SET slug=?, title=?, excerpt=?, body=?, cover_url=?, published=?, published_at=?, updated_at=datetime('now')
		WHERE id=?`,
		p.Slug, p.Title, p.Excerpt, p.Body, p.CoverURL, boolToInt(p.Published), nullStr(p.PublishedAt), id)
	return err
}

func (s *Store) DeletePost(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM blog_posts WHERE id = ?`, id)
	return err
}

func (s *Store) CountPosts(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM blog_posts`).Scan(&n)
	return n, err
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

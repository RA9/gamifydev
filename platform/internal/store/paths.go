package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// Path is a career path: an ordered curriculum of courses toward a job outcome.
type Path struct {
	ID          int64
	Slug        string
	Title       string
	Tagline     string
	Description string
	Emoji       string
	Level       string
	Sort        int
	Published   bool
	CourseCount int // populated by listing queries
	LessonCount int // total lessons across the path's courses
}

const pathCols = `p.id, p.slug, p.title, p.tagline, p.description, p.emoji, p.level, p.sort, p.published,
	(SELECT COUNT(*) FROM path_courses pc WHERE pc.path_id = p.id) AS course_count,
	(SELECT COALESCE(SUM((SELECT COUNT(*) FROM lessons l WHERE l.course_id = pc.course_id)), 0)
	   FROM path_courses pc WHERE pc.path_id = p.id) AS lesson_count`

func scanPath(sc interface{ Scan(...any) error }) (*Path, error) {
	var p Path
	var pub int
	if err := sc.Scan(&p.ID, &p.Slug, &p.Title, &p.Tagline, &p.Description, &p.Emoji, &p.Level,
		&p.Sort, &pub, &p.CourseCount, &p.LessonCount); err != nil {
		return nil, err
	}
	p.Published = pub == 1
	return &p, nil
}

// ListPaths returns paths (with course/lesson counts) ordered for display.
func (s *Store) ListPaths(ctx context.Context, includeUnpublished bool) ([]Path, error) {
	q := `SELECT ` + pathCols + ` FROM paths p`
	if !includeUnpublished {
		q += ` WHERE p.published = 1`
	}
	q += ` ORDER BY p.sort, p.title`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Path
	for rows.Next() {
		p, err := scanPath(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, rows.Err()
}

func (s *Store) GetPathBySlug(ctx context.Context, slug string) (*Path, error) {
	p, err := scanPath(s.db.QueryRowContext(ctx, `SELECT `+pathCols+` FROM paths p WHERE p.slug = ?`, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *Store) GetPathByID(ctx context.Context, id int64) (*Path, error) {
	p, err := scanPath(s.db.QueryRowContext(ctx, `SELECT `+pathCols+` FROM paths p WHERE p.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

// PathCourses returns a path's courses (with lesson counts) in curriculum order.
func (s *Store) PathCourses(ctx context.Context, pathID int64) ([]Course, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.slug, c.title, c.emoji, c.tagline, c.description, c.sort, c.published,
		       (SELECT COUNT(*) FROM lessons l WHERE l.course_id = c.id) AS lesson_count
		FROM path_courses pc JOIN courses c ON c.id = pc.course_id
		WHERE pc.path_id = ?
		ORDER BY pc.sort, c.title`, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Course
	for rows.Next() {
		var c Course
		var pub int
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Emoji, &c.Tagline, &c.Description, &c.Sort, &pub, &c.LessonCount); err != nil {
			return nil, err
		}
		c.Published = pub == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

// PathCourseIDs returns the set of course ids assigned to a path (for the admin
// editor to mark which courses are included).
func (s *Store) PathCourseIDs(ctx context.Context, pathID int64) (map[int64]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT course_id, sort FROM path_courses WHERE path_id = ?`, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]int{}
	for rows.Next() {
		var id int64
		var sort int
		if err := rows.Scan(&id, &sort); err != nil {
			return nil, err
		}
		out[id] = sort
	}
	return out, rows.Err()
}

func (s *Store) CreatePath(ctx context.Context, p Path) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO paths (slug, title, tagline, description, emoji, level, sort, published)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		p.Slug, p.Title, p.Tagline, p.Description, p.Emoji, p.Level, p.Sort, boolToInt(p.Published)).Scan(&id)
	return id, err
}

func (s *Store) UpdatePath(ctx context.Context, id int64, p Path) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE paths SET slug=?, title=?, tagline=?, description=?, emoji=?, level=?, published=?,
		 updated_at=datetime('now') WHERE id=?`,
		p.Slug, p.Title, p.Tagline, p.Description, p.Emoji, p.Level, boolToInt(p.Published), id)
	return err
}

// UpsertPath inserts or updates a path by slug (used by the seeder); returns id.
func (s *Store) UpsertPath(ctx context.Context, p Path) (int64, error) {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO paths (slug, title, tagline, description, emoji, level, sort, published, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(slug) DO UPDATE SET
			title=excluded.title, tagline=excluded.tagline, description=excluded.description,
			emoji=excluded.emoji, level=excluded.level, sort=excluded.sort, published=excluded.published,
			updated_at=datetime('now')`,
		p.Slug, p.Title, p.Tagline, p.Description, p.Emoji, p.Level, p.Sort, boolToInt(p.Published))
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.db.QueryRowContext(ctx, `SELECT id FROM paths WHERE slug = ?`, p.Slug).Scan(&id)
	return id, err
}

func (s *Store) DeletePath(ctx context.Context, id int64) error {
	// Foreign keys are disabled, so remove memberships explicitly.
	if _, err := s.db.ExecContext(ctx, `DELETE FROM path_courses WHERE path_id = ?`, id); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `DELETE FROM paths WHERE id = ?`, id)
	return err
}

// SetPathCourses replaces a path's course list with the given course ids, in the
// given order (index becomes sort).
func (s *Store) SetPathCourses(ctx context.Context, pathID int64, courseIDs []int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(ctx, `DELETE FROM path_courses WHERE path_id = ?`, pathID); err != nil {
		return err
	}
	for i, cid := range courseIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO path_courses (path_id, course_id, sort) VALUES (?, ?, ?)`, pathID, cid, i); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// SetPathCoursesBySlug resolves course slugs to ids and assigns them (used by the
// seeder). Unknown slugs are skipped.
func (s *Store) SetPathCoursesBySlug(ctx context.Context, pathID int64, slugs []string) error {
	var ids []int64
	for _, sl := range slugs {
		c, err := s.GetCourseBySlug(ctx, strings.TrimSpace(sl))
		if err != nil {
			continue
		}
		ids = append(ids, c.ID)
	}
	return s.SetPathCourses(ctx, pathID, ids)
}

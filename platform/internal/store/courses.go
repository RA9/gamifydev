package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

const (
	ModeFun         = "fun"
	ModeTheoretical = "theoretical"
	ModePractical   = "practical"

	LanguageNone = "none"

	LanguagePolicyMixed          = "mixed"
	LanguagePolicyCOnly          = "c_only"
	LanguagePolicyCWithException = "c_with_exception"
)

type Course struct {
	ID                int64
	Slug              string
	Title             string
	Emoji             string
	Tagline           string
	Description       string
	Sort              int
	Published         bool
	PrimaryLanguage   string
	LanguagePolicy    string
	LanguageException string
	LessonCount       int // populated by listing queries
}

type Lesson struct {
	ID              int64
	CourseID        int64
	Slug            string
	Title           string
	Summary         string
	Body            string
	VideoURL        string
	AudioURL        string
	Sort            int
	Section         string // heading this lesson sits under ("" = ungrouped)
	Kind            string // theory | lab | workshop | project
	WorkloadMinutes int
	Mode            string // fun | theoretical | practical
	Language        string // none for language-neutral reading
}

var ErrInvalidCourseLanguagePolicy = errors.New("invalid course language policy")

func normalizeCourseMetadata(c *Course) error {
	if c.LanguagePolicy == "" {
		c.LanguagePolicy = LanguagePolicyMixed
	}
	switch c.LanguagePolicy {
	case LanguagePolicyMixed, LanguagePolicyCOnly:
		c.LanguageException = ""
	case LanguagePolicyCWithException:
		if strings.TrimSpace(c.LanguageException) == "" {
			return ErrInvalidCourseLanguagePolicy
		}
	default:
		return ErrInvalidCourseLanguagePolicy
	}
	return nil
}

func normalizeLessonMetadata(l *Lesson) {
	if l.WorkloadMinutes <= 0 {
		l.WorkloadMinutes = 30
	}
	if l.Mode == "" {
		if l.Kind == "lab" || l.Kind == "workshop" || l.Kind == "project" {
			l.Mode = ModePractical
		} else {
			l.Mode = ModeTheoretical
		}
	}
	if l.Language == "" {
		l.Language = LanguageNone
	}
}

// UpsertCourse inserts or updates a course by slug; returns its id.
func (s *Store) UpsertCourse(ctx context.Context, c Course) (int64, error) {
	if err := normalizeCourseMetadata(&c); err != nil {
		return 0, err
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO courses
			(slug, title, emoji, tagline, description, sort, published,
			 primary_language, language_policy, language_exception, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(slug) DO UPDATE SET
			title=excluded.title, emoji=excluded.emoji, tagline=excluded.tagline,
			description=excluded.description, sort=excluded.sort, published=excluded.published,
			primary_language=excluded.primary_language, language_policy=excluded.language_policy,
			language_exception=excluded.language_exception, updated_at=datetime('now')`,
		c.Slug, c.Title, c.Emoji, c.Tagline, c.Description, c.Sort, boolToInt(c.Published),
		c.PrimaryLanguage, c.LanguagePolicy, c.LanguageException)
	if err != nil {
		return 0, err
	}
	var id int64
	err = s.db.QueryRowContext(ctx, `SELECT id FROM courses WHERE slug = ?`, c.Slug).Scan(&id)
	return id, err
}

// UpsertLesson inserts or updates a lesson by (course_id, slug).
func (s *Store) UpsertLesson(ctx context.Context, l Lesson) error {
	normalizeLessonMetadata(&l)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO lessons
			(course_id, slug, title, summary, body, video_url, audio_url, sort,
			 section, kind, workload_minutes, mode, language, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(course_id, slug) DO UPDATE SET
			title=excluded.title, summary=excluded.summary, body=excluded.body,
			video_url=excluded.video_url, audio_url=excluded.audio_url, sort=excluded.sort,
			section=excluded.section, kind=excluded.kind,
			workload_minutes=excluded.workload_minutes, mode=excluded.mode,
			language=excluded.language, updated_at=datetime('now')`,
		l.CourseID, l.Slug, l.Title, l.Summary, l.Body, l.VideoURL, l.AudioURL, l.Sort,
		l.Section, l.Kind, l.WorkloadMinutes, l.Mode, l.Language)
	return err
}

// ListCourses returns published courses (with lesson counts) ordered for display.
func (s *Store) ListCourses(ctx context.Context, includeUnpublished bool) ([]Course, error) {
	q := `SELECT c.id, c.slug, c.title, c.emoji, c.tagline, c.description, c.sort, c.published,
		       c.primary_language, c.language_policy, c.language_exception,
		       (SELECT COUNT(*) FROM lessons l WHERE l.course_id = c.id) AS lesson_count
		  FROM courses c`
	if !includeUnpublished {
		q += ` WHERE c.published = 1`
	}
	q += ` ORDER BY c.sort, c.title`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Course
	for rows.Next() {
		var c Course
		var pub int
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Emoji, &c.Tagline, &c.Description, &c.Sort, &pub,
			&c.PrimaryLanguage, &c.LanguagePolicy, &c.LanguageException, &c.LessonCount); err != nil {
			return nil, err
		}
		c.Published = pub == 1
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetCourseBySlug(ctx context.Context, slug string) (*Course, error) {
	var c Course
	var pub int
	err := s.db.QueryRowContext(ctx,
		`SELECT id, slug, title, emoji, tagline, description, sort, published,
		        primary_language, language_policy, language_exception
		 FROM courses WHERE slug = ?`, slug).
		Scan(&c.ID, &c.Slug, &c.Title, &c.Emoji, &c.Tagline, &c.Description, &c.Sort, &pub,
			&c.PrimaryLanguage, &c.LanguagePolicy, &c.LanguageException)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c.Published = pub == 1
	return &c, nil
}

// ListLessons returns a course's lessons in order (without bodies, for menus).
func (s *Store) ListLessons(ctx context.Context, courseID int64) ([]Lesson, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, course_id, slug, title, summary, sort, section, kind,
		        workload_minutes, mode, language
		 FROM lessons WHERE course_id = ? ORDER BY sort`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Slug, &l.Title, &l.Summary, &l.Sort,
			&l.Section, &l.Kind, &l.WorkloadMinutes, &l.Mode, &l.Language); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *Store) GetLesson(ctx context.Context, courseID int64, slug string) (*Lesson, error) {
	var l Lesson
	err := s.db.QueryRowContext(ctx,
		`SELECT id, course_id, slug, title, summary, body, video_url, audio_url, sort,
		        section, kind, workload_minutes, mode, language
		 FROM lessons WHERE course_id = ? AND slug = ?`, courseID, slug).
		Scan(&l.ID, &l.CourseID, &l.Slug, &l.Title, &l.Summary, &l.Body, &l.VideoURL,
			&l.AudioURL, &l.Sort, &l.Section, &l.Kind, &l.WorkloadMinutes, &l.Mode, &l.Language)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

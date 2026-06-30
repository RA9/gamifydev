package store

import (
	"context"
	"database/sql"
	"errors"
)

// --- Courses (admin authoring) ----------------------------------------------

func (s *Store) GetCourseByID(ctx context.Context, id int64) (*Course, error) {
	var c Course
	var pub int
	err := s.db.QueryRowContext(ctx,
		`SELECT id, slug, title, emoji, tagline, description, sort, published FROM courses WHERE id = ?`, id).
		Scan(&c.ID, &c.Slug, &c.Title, &c.Emoji, &c.Tagline, &c.Description, &c.Sort, &pub)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c.Published = pub == 1
	return &c, nil
}

func (s *Store) CreateCourse(ctx context.Context, c Course) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO courses (slug, title, emoji, tagline, description, sort, published)
		 VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		c.Slug, c.Title, c.Emoji, c.Tagline, c.Description, c.Sort, boolToInt(c.Published)).Scan(&id)
	return id, err
}

func (s *Store) UpdateCourse(ctx context.Context, id int64, c Course) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE courses SET slug=?, title=?, emoji=?, tagline=?, description=?, published=?, updated_at=datetime('now')
		 WHERE id=?`,
		c.Slug, c.Title, c.Emoji, c.Tagline, c.Description, boolToInt(c.Published), id)
	return err
}

// --- Lessons (admin authoring) ----------------------------------------------

func (s *Store) GetLessonByID(ctx context.Context, id int64) (*Lesson, error) {
	var l Lesson
	err := s.db.QueryRowContext(ctx,
		`SELECT id, course_id, slug, title, summary, body, video_url, audio_url, sort, section, kind FROM lessons WHERE id = ?`, id).
		Scan(&l.ID, &l.CourseID, &l.Slug, &l.Title, &l.Summary, &l.Body, &l.VideoURL, &l.AudioURL, &l.Sort, &l.Section, &l.Kind)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *Store) CreateLesson(ctx context.Context, l Lesson) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO lessons (course_id, slug, title, summary, body, video_url, audio_url, sort, section, kind)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.CourseID, l.Slug, l.Title, l.Summary, l.Body, l.VideoURL, l.AudioURL, l.Sort, l.Section, l.Kind)
	return err
}

func (s *Store) UpdateLesson(ctx context.Context, id int64, l Lesson) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE lessons SET slug=?, title=?, summary=?, body=?, video_url=?, audio_url=?, sort=?, section=?, kind=?, updated_at=datetime('now')
		 WHERE id=?`,
		l.Slug, l.Title, l.Summary, l.Body, l.VideoURL, l.AudioURL, l.Sort, l.Section, l.Kind, id)
	return err
}

func (s *Store) DeleteLesson(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM lessons WHERE id = ?`, id)
	return err
}

// NextLessonSort returns the next sort index for a new lesson in a course.
func (s *Store) NextLessonSort(ctx context.Context, courseID int64) int {
	var n sql.NullInt64
	_ = s.db.QueryRowContext(ctx, `SELECT MAX(sort) FROM lessons WHERE course_id = ?`, courseID).Scan(&n)
	if n.Valid {
		return int(n.Int64) + 1
	}
	return 0
}

// --- Assignments (admin authoring) ------------------------------------------

func (s *Store) GetAssignmentByID(ctx context.Context, id int64) (*Assignment, error) {
	a, err := scanAssignment(s.db.QueryRowContext(ctx,
		`SELECT `+assignmentCols+` FROM assignments a LEFT JOIN courses c ON c.id = a.course_id WHERE a.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

func (s *Store) CreateAssignment(ctx context.Context, a Assignment) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO assignments (course_id, slug, title, language, prompt, starter, max_points, published, sort, required, pass_points)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.CourseID, a.Slug, a.Title, a.Language, a.Prompt, a.Starter, a.MaxPoints, boolToInt(a.Published), a.Sort,
		boolToInt(a.Required), a.PassPoints)
	return err
}

func (s *Store) UpdateAssignment(ctx context.Context, id int64, a Assignment) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE assignments SET course_id=?, slug=?, title=?, language=?, prompt=?, starter=?, max_points=?, published=?, required=?, pass_points=?, updated_at=datetime('now')
		 WHERE id=?`,
		a.CourseID, a.Slug, a.Title, a.Language, a.Prompt, a.Starter, a.MaxPoints, boolToInt(a.Published),
		boolToInt(a.Required), a.PassPoints, id)
	return err
}

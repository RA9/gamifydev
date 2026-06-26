-- Career paths: a curated, ordered curriculum of courses aimed at a job outcome
-- (e.g. "Full-Stack Developer"). A path groups existing courses; learners follow
-- it top to bottom. Path -> Course -> Lesson.

CREATE TABLE IF NOT EXISTS paths (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  slug        TEXT NOT NULL UNIQUE,
  title       TEXT NOT NULL,
  tagline     TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  emoji       TEXT NOT NULL DEFAULT '',
  level       TEXT NOT NULL DEFAULT 'Beginner',
  sort        INTEGER NOT NULL DEFAULT 0,
  published   INTEGER NOT NULL DEFAULT 0,
  created_at  TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Ordered membership of courses within a path.
CREATE TABLE IF NOT EXISTS path_courses (
  path_id   INTEGER NOT NULL REFERENCES paths(id) ON DELETE CASCADE,
  course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  sort      INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (path_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_path_courses_path ON path_courses(path_id, sort);

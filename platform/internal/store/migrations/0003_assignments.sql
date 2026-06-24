-- Coding assignments and their learner submissions. Submissions are graded by
-- a human (admin/grader) — the core mentor-feedback loop.

CREATE TABLE IF NOT EXISTS assignments (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id   INTEGER REFERENCES courses(id) ON DELETE SET NULL,
    slug        TEXT NOT NULL UNIQUE,
    title       TEXT NOT NULL,
    language    TEXT NOT NULL DEFAULT 'text',
    prompt      TEXT NOT NULL DEFAULT '',
    starter     TEXT NOT NULL DEFAULT '',
    max_points  INTEGER NOT NULL DEFAULT 100,
    published   INTEGER NOT NULL DEFAULT 1,
    sort        INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS submissions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    assignment_id INTEGER NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code          TEXT NOT NULL DEFAULT '',
    note          TEXT NOT NULL DEFAULT '',           -- learner's note to the grader
    status        TEXT NOT NULL DEFAULT 'submitted'   -- submitted | graded | returned
                  CHECK (status IN ('submitted','graded','returned')),
    score         INTEGER,
    feedback      TEXT NOT NULL DEFAULT '',
    graded_by     INTEGER REFERENCES users(id) ON DELETE SET NULL,
    graded_at     TEXT,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status, created_at);
CREATE INDEX IF NOT EXISTS idx_submissions_user ON submissions(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_submissions_assignment ON submissions(assignment_id);

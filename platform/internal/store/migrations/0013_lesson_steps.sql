-- Interactive, step-by-step lessons (FreeCodeCamp-style). A lesson with steps is
-- rendered as a guided editor: each step shows an instruction + starter code, the
-- learner edits, and author-defined checks verify the result before advancing.

CREATE TABLE IF NOT EXISTS lesson_steps (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  lesson_id   INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
  sort        INTEGER NOT NULL DEFAULT 0,
  instruction TEXT NOT NULL DEFAULT '',   -- markdown shown for the step
  starter     TEXT NOT NULL DEFAULT '',   -- code the editor starts with
  checks      TEXT NOT NULL DEFAULT '[]', -- JSON: [{"text":"...","test":"<js bool expr>"}]
  created_at  TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_lesson_steps_lesson ON lesson_steps(lesson_id, sort);

-- Per-learner completion of a step.
CREATE TABLE IF NOT EXISTS step_progress (
  user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  step_id      INTEGER NOT NULL REFERENCES lesson_steps(id) ON DELETE CASCADE,
  completed_at TEXT NOT NULL DEFAULT (datetime('now')),
  PRIMARY KEY (user_id, step_id)
);

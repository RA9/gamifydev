-- Fixture files for a checkpoint.
--
-- A Linux checkpoint is usually "do something with these files", and a data
-- exercise in any language often needs input to read. Without fixtures those
-- checkpoints cannot exist at all.
--
-- Files are written into the sandbox's scratch directory before the learner's
-- program runs, and vanish with it. Names are validated in the runner: a plain
-- filename only, so a fixture cannot be used to write outside the sandbox.
CREATE TABLE IF NOT EXISTS assignment_files (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  assignment_id INTEGER NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
  sort          INTEGER NOT NULL DEFAULT 0,
  name          TEXT NOT NULL,
  content       TEXT NOT NULL DEFAULT '',
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_files_assignment ON assignment_files(assignment_id, sort);
CREATE UNIQUE INDEX IF NOT EXISTS idx_files_name ON assignment_files(assignment_id, name);

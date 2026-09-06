-- Proctoring for the placement diagnostic.
--
-- The count lives on the attempt because that is what the rule is about — two
-- strikes inside one sitting — and the log lives beside it because a decision
-- that voids somebody's admissions test should be reviewable afterwards by a
-- human. A bare counter can tell a learner they cheated; it cannot tell an
-- administrator what it saw, or let them overturn it.
ALTER TABLE assessment_attempts ADD COLUMN violations INTEGER NOT NULL DEFAULT 0;

-- Empty for an ordinary attempt; 'cheating' for one the proctor closed. Kept as
-- a reason rather than a boolean so a future rule (a second invigilator, a
-- withdrawn sitting) has somewhere to say what happened.
ALTER TABLE assessment_attempts ADD COLUMN voided_reason TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS attempt_violations (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  attempt_id INTEGER NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
  -- What the browser reported: 'copy', 'hidden', 'blur', 'print', 'capture'.
  kind       TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_attempt_violations ON attempt_violations(attempt_id, id);

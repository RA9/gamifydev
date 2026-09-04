-- Automated checks on assignments — step 1 of the checkpoint grading design.
--
-- Checkpoints are the assignments that gate progression between courses, and
-- until now they were graded entirely by hand. This adds the half of the
-- judgment a machine can make ("does it work?") so a human is only needed for
-- the half it cannot ("is it any good?", "is something wrong here?").
--
-- The shape deliberately mirrors `lesson_steps.checks`, which already works:
-- an authored boolean expression evaluated after the learner's program runs.
-- The server-side lab harness that evaluates them is reused as-is.
--
-- Constraint worth stating plainly: the sandbox executes PYTHON ONLY. An
-- assignment in another language can carry checks for documentation, but they
-- cannot be run, and it falls back to mentor grading. `auto_gradable` records
-- that decision per assignment rather than leaving it implicit.

CREATE TABLE IF NOT EXISTS assignment_checks (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  assignment_id INTEGER NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
  sort          INTEGER NOT NULL DEFAULT 0,
  -- What the learner is told this check verifies. Written for them, not for us.
  label         TEXT NOT NULL,
  -- A Python boolean expression evaluated in the namespace left behind by the
  -- learner's program, with `_code` (their source) and `client` (a test client
  -- for any web app they defined) also in scope.
  test          TEXT NOT NULL,
  -- Hidden checks still gate, but their label is withheld until after a pass,
  -- so a learner cannot reverse-engineer the full spec from the failure list.
  hidden        INTEGER NOT NULL DEFAULT 0,
  -- Weight toward the score. Zero means the check is advisory: it shows in the
  -- results but does not affect pass/fail.
  points        INTEGER NOT NULL DEFAULT 1,
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_checks_assignment ON assignment_checks(assignment_id, sort);

-- Per-check outcome for one submission. Kept per row rather than as a blob so
-- the admin can see which check fails most often — that is the signal a
-- checkpoint is badly worded rather than the learners being wrong.
CREATE TABLE IF NOT EXISTS submission_checks (
  submission_id INTEGER NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
  check_id      INTEGER NOT NULL REFERENCES assignment_checks(id) ON DELETE CASCADE,
  passed        INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (submission_id, check_id)
);

CREATE INDEX IF NOT EXISTS idx_subchecks_check ON submission_checks(check_id, passed);

-- Whether this assignment's checks can actually be executed. Set when an
-- assignment is Python and has at least one check; the runner still has to be
-- enabled on the instance for it to mean anything.
ALTER TABLE assignments ADD COLUMN auto_gradable INTEGER NOT NULL DEFAULT 0;

-- Submissions gain the machine's verdict, kept separate from the human's.
--
--   checks_ran_at  — when the harness last ran against this submission
--   checks_passed  — how many checks passed, and of how many that count
--
-- `status` and `score` keep their existing meaning, so the progression gate in
-- AssignmentPassState is untouched: auto-grading works by *setting* them, not
-- by adding a second parallel notion of passing.
ALTER TABLE submissions ADD COLUMN checks_ran_at TEXT;
ALTER TABLE submissions ADD COLUMN checks_passed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE submissions ADD COLUMN checks_total   INTEGER NOT NULL DEFAULT 0;
-- The harness's stdout/stderr, so a learner sees why a check failed rather than
-- just that it did.
ALTER TABLE submissions ADD COLUMN checks_output TEXT NOT NULL DEFAULT '';

-- Phase 2 of the cohort program: the placement diagnostic and the routing
-- decision it produces.
--
-- The decision recorded in decisions/03 of the PRD shapes this schema: the test
-- ROUTES (picks a path, proposes exemptions) but never by itself lets a learner
-- past a course — that still requires the course's checkpoint assignment. So
-- `placement_results` stores a *proposal*, and nothing here grants credit.

-- --------------------------------------------------------------------------
-- An assessment is a named bank of items. Today there is exactly one
-- ('placement'), but the shape supports per-path diagnostics later.
CREATE TABLE IF NOT EXISTS assessments (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  slug         TEXT NOT NULL UNIQUE,
  title        TEXT NOT NULL DEFAULT '',
  kind         TEXT NOT NULL DEFAULT 'placement',
  time_limit_s INTEGER NOT NULL DEFAULT 1800,
  -- How many items to sample per topic for one sitting. Sampling (rather than
  -- serving the whole bank in a fixed order) is the leak mitigation: two
  -- learners rarely see the same paper.
  per_topic    INTEGER NOT NULL DEFAULT 4,
  published    INTEGER NOT NULL DEFAULT 1,
  created_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

-- --------------------------------------------------------------------------
-- The item bank. `topic` is what makes routing possible: scores are computed
-- per topic, and a topic maps to the courses it can exempt.
--
-- `answer` is the 0-based index into `options` (a JSON array of strings).
-- Answers live server-side only and are never rendered into the page.
CREATE TABLE IF NOT EXISTS assessment_items (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  assessment_id INTEGER NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
  topic         TEXT NOT NULL,
  difficulty    INTEGER NOT NULL DEFAULT 1,   -- 1 easy .. 3 hard
  prompt        TEXT NOT NULL,
  code          TEXT NOT NULL DEFAULT '',     -- optional snippet shown above the options
  lang          TEXT NOT NULL DEFAULT '',     -- language for the snippet's highlighting
  options       TEXT NOT NULL DEFAULT '[]',   -- JSON array of option strings
  answer        INTEGER NOT NULL DEFAULT 0,   -- index into options
  explanation   TEXT NOT NULL DEFAULT '',     -- shown on the result page
  retired       INTEGER NOT NULL DEFAULT 0,   -- rotate items out without deleting history
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_items_bank ON assessment_items(assessment_id, topic, retired);

-- --------------------------------------------------------------------------
-- One sitting. `expires_at` is written at start so the time limit is enforced
-- server-side; a client clock cannot buy extra time.
CREATE TABLE IF NOT EXISTS assessment_attempts (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  assessment_id INTEGER NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
  started_at    TEXT NOT NULL DEFAULT (datetime('now')),
  expires_at    TEXT NOT NULL,
  submitted_at  TEXT,
  score         INTEGER,          -- overall percentage, 0-100
  topic_scores  TEXT NOT NULL DEFAULT '{}',  -- JSON: {"topic": pct}
  abandoned     INTEGER NOT NULL DEFAULT 0   -- expired without a submission
);

CREATE INDEX IF NOT EXISTS idx_attempts_user ON assessment_attempts(user_id, started_at DESC);

-- The item set drawn for this attempt, in presentation order. Recorded at start
-- so a reload shows the same paper, and so scoring cannot be influenced by a
-- client resubmitting a different set of questions.
CREATE TABLE IF NOT EXISTS assessment_attempt_items (
  attempt_id INTEGER NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
  item_id    INTEGER NOT NULL REFERENCES assessment_items(id) ON DELETE CASCADE,
  sort       INTEGER NOT NULL DEFAULT 0,
  response   INTEGER,          -- chosen option index; NULL = unanswered
  correct    INTEGER,          -- 1/0, written at scoring time
  PRIMARY KEY (attempt_id, item_id)
);

-- --------------------------------------------------------------------------
-- The routing decision produced from an attempt.
--
-- `exemptions` is a JSON array of course slugs the learner may skip in the
-- schedule. Per the PRD this is a proposal about *scheduling*, not credit, and
-- exempted content stays readable.
CREATE TABLE IF NOT EXISTS placement_results (
  id                   INTEGER PRIMARY KEY AUTOINCREMENT,
  attempt_id           INTEGER NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
  user_id              INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recommended_path_id  INTEGER REFERENCES paths(id) ON DELETE SET NULL,
  foundations_required INTEGER NOT NULL DEFAULT 1,
  exemptions           TEXT NOT NULL DEFAULT '[]',
  created_at           TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_placement_user ON placement_results(user_id, created_at DESC);

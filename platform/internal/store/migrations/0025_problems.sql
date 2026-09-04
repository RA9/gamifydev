-- The practice problem bank: DSA exercises anyone can attempt, signed in or
-- not.
--
-- Deliberately separate from assignments. An assignment is coursework that
-- gates a path; a problem is open practice that gates nothing. Sharing a table
-- would mean one careless query letting problem points open a course.
CREATE TABLE IF NOT EXISTS problems (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  slug          TEXT NOT NULL UNIQUE,
  title         TEXT NOT NULL,
  difficulty    TEXT NOT NULL DEFAULT 'easy' CHECK (difficulty IN ('easy','medium','hard')),
  topic         TEXT NOT NULL DEFAULT '',
  statement     TEXT NOT NULL DEFAULT '',
  time_limit_ms INTEGER NOT NULL DEFAULT 5000,
  sort          INTEGER NOT NULL DEFAULT 0,
  published     INTEGER NOT NULL DEFAULT 1,
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

-- One starter per language the problem accepts. The set of rows here is also
-- what decides which languages the editor offers.
CREATE TABLE IF NOT EXISTS problem_starters (
  problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  language   TEXT NOT NULL,
  code       TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (problem_id, language)
);

-- Same shape as assignment_checks, so a problem's tests map straight onto
-- store.Check and are judged by exactly the code that grades checkpoints.
CREATE TABLE IF NOT EXISTS problem_tests (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  sort       INTEGER NOT NULL DEFAULT 0,
  label      TEXT NOT NULL DEFAULT '',
  test       TEXT NOT NULL DEFAULT '',
  stdin      TEXT NOT NULL DEFAULT '',
  hidden     INTEGER NOT NULL DEFAULT 0,
  points     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_problem_tests_problem ON problem_tests(problem_id, sort);

-- A guest's identity, held in a cookie so an anonymous visitor's solved set
-- survives a page reload and can be handed to an account later.
--
-- The token is stored hashed, like a password reset: it is a bearer credential,
-- and there is no reason for the database to hold a replayable copy.
CREATE TABLE IF NOT EXISTS guest_sessions (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  token_hash   TEXT NOT NULL UNIQUE,
  claimed_by   INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at   TEXT NOT NULL DEFAULT (datetime('now')),
  last_seen_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_guest_sessions_hash ON guest_sessions(token_hash);

-- Exactly one of user_id / guest_id identifies the author. A guest's rows are
-- rewritten to user_id when they claim the session by signing up.
CREATE TABLE IF NOT EXISTS problem_submissions (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  user_id    INTEGER REFERENCES users(id) ON DELETE CASCADE,
  guest_id   INTEGER REFERENCES guest_sessions(id) ON DELETE CASCADE,
  language   TEXT NOT NULL DEFAULT 'python',
  code       TEXT NOT NULL DEFAULT '',
  verdict    TEXT NOT NULL DEFAULT 'pending',
  passed     INTEGER NOT NULL DEFAULT 0,
  total      INTEGER NOT NULL DEFAULT 0,
  output     TEXT NOT NULL DEFAULT '',
  -- The first visible test the program got wrong. A bare "9/12" tells someone
  -- their code is broken without telling them what it got wrong, which is the
  -- one thing they need; naming the requirement turns a score into a lead.
  -- Hidden tests never appear here — that is what makes them hidden.
  failed_label TEXT NOT NULL DEFAULT '',
  runtime_ms INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  judged_at  TEXT,
  CHECK ((user_id IS NULL) <> (guest_id IS NULL))
);
CREATE INDEX IF NOT EXISTS idx_problem_subs_pending ON problem_submissions(verdict, id);
CREATE INDEX IF NOT EXISTS idx_problem_subs_user ON problem_submissions(user_id, problem_id);
CREATE INDEX IF NOT EXISTS idx_problem_subs_guest ON problem_submissions(guest_id, problem_id);

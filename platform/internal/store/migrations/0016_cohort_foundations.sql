-- Phase 1 of the cohort program: the substrate every later phase builds on.
--
--   * account_status  — may this account participate, and until when
--   * enrollments     — a learner's relationship to a path (and later, a cohort)
--   * activity_days   — one row per learner per active day (the attendance substrate)
--   * job_runs/locks  — the scheduled-job runner's history and single-flight leases
--
-- Cohort, standup, assessment and sanction tables arrive in later phases. The
-- columns those phases need on `enrollments` (cohort_id) are added then, so this
-- migration stays honest about what actually exists today.

-- --------------------------------------------------------------------------
-- Account participation state.
--
-- One row per user; absent row means "active" so existing users need no backfill.
-- `state` is the account-level gate (a suspension stops participation entirely);
-- enrollment state below is the per-path gate.
CREATE TABLE IF NOT EXISTS account_status (
  user_id      INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  state        TEXT NOT NULL DEFAULT 'active',  -- active | suspended | banned
  reason       TEXT NOT NULL DEFAULT '',
  effective_at TEXT NOT NULL DEFAULT (datetime('now')),
  expires_at   TEXT,                            -- NULL = indefinite
  updated_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_account_status_expiry
  ON account_status(state, expires_at);

-- --------------------------------------------------------------------------
-- Enrollment: the learner's relationship to a path.
--
-- States (transitions enforced in Go, see internal/store/enrollment.go):
--   unplaced  -> placed              diagnostic completed, path decided
--   placed    -> active              cohort started
--   active    -> paused | dropped | completed
--   paused    -> active | dropped    deferral: rejoin a later cohort
--
-- A learner has at most one non-terminal enrollment at a time; the partial
-- unique index below enforces that in the database rather than by convention.
CREATE TABLE IF NOT EXISTS enrollments (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  path_id    INTEGER REFERENCES paths(id) ON DELETE SET NULL,
  state      TEXT NOT NULL DEFAULT 'unplaced',
  reason     TEXT NOT NULL DEFAULT '',   -- why it ended, for dropped/paused
  placed_at  TEXT,
  started_at TEXT,
  ended_at   TEXT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_enrollments_user  ON enrollments(user_id, state);
CREATE INDEX IF NOT EXISTS idx_enrollments_path  ON enrollments(path_id, state);

-- At most one live enrollment per user. Terminal states (dropped, completed)
-- are excluded so history accumulates.
CREATE UNIQUE INDEX IF NOT EXISTS idx_enrollments_one_live
  ON enrollments(user_id)
  WHERE state IN ('unplaced', 'placed', 'active', 'paused');

-- --------------------------------------------------------------------------
-- Activity ledger: one row per learner per day they did something real.
--
-- Materialized by the `activity:rollup` job from signals that already exist
-- (completed steps, graded submissions). `sources` is a comma-separated set of
-- the signal kinds seen that day, kept for debugging why a day counted.
--
-- `day` is a UTC date string (YYYY-MM-DD). Cohort-local attendance in phase 5
-- resolves against the cohort's timezone band; this table stays timezone-neutral.
CREATE TABLE IF NOT EXISTS activity_days (
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  day        TEXT NOT NULL,
  sources    TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  PRIMARY KEY (user_id, day)
);

CREATE INDEX IF NOT EXISTS idx_activity_days_day ON activity_days(day);

-- --------------------------------------------------------------------------
-- Scheduled job bookkeeping.
--
-- job_locks is a lease table: a runner claims a job name until `expires_at`, so
-- two app instances never run the same job on the same tick. Leases expire on
-- their own, so a crashed runner does not wedge the schedule.
CREATE TABLE IF NOT EXISTS job_locks (
  name       TEXT PRIMARY KEY,
  holder     TEXT NOT NULL,
  acquired_at TEXT NOT NULL DEFAULT (datetime('now')),
  expires_at TEXT NOT NULL
);

-- job_runs is the history: what ran, how long it took, and what it did or failed
-- to do. Surfaced in the admin UI so the schedule is observable.
CREATE TABLE IF NOT EXISTS job_runs (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  started_at  TEXT NOT NULL DEFAULT (datetime('now')),
  finished_at TEXT,
  ok          INTEGER,           -- NULL while running, 1 ok, 0 failed
  detail      TEXT NOT NULL DEFAULT '',
  duration_ms INTEGER
);

CREATE INDEX IF NOT EXISTS idx_job_runs_name ON job_runs(name, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_job_runs_recent ON job_runs(started_at DESC);

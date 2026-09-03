-- Phase 5 of the cohort program: attendance, sanctions, and appeals.
--
-- This is the phase the PRD deliberately put last and deliberately ships inert
-- (§07). Everything here records what *would* happen; nothing is applied until
-- the thresholds have been calibrated against a real cohort's attendance. The
-- `shadow` column is what makes that possible, and it defaults to 1.
--
-- The governing decision (§03) is that the offense is going dark, not being
-- absent. A learner who files a notice is excused. A learner who vanishes
-- without a word accrues an unexcused absence. Removal is a DROP — they may
-- reapply — not a ban; bans are reserved for conduct.

-- --------------------------------------------------------------------------
-- "I'll be out" — the flow that converts an absence into an excused one.
--
-- Filed ahead of time or shortly after; both are fine, because the behaviour
-- being taught is telling your team, not predicting your life.
CREATE TABLE IF NOT EXISTS absence_notices (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  cohort_id  INTEGER REFERENCES cohorts(id) ON DELETE SET NULL,
  from_day   TEXT NOT NULL,           -- YYYY-MM-DD inclusive
  to_day     TEXT NOT NULL,           -- YYYY-MM-DD inclusive
  reason     TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_notices_user ON absence_notices(user_id, from_day, to_day);

-- --------------------------------------------------------------------------
-- Resolved daily state, one row per learner per SCHEDULED day.
--
-- Only scheduled days are recorded. Weekends and rest days carry no work, so
-- counting them as absences would punish the rest the schedule deliberately
-- builds in.
--
--   present  — posted standup, or completed scheduled work, that day
--   excused  — covered by an absence notice
--   absent   — silence on a day work was due
CREATE TABLE IF NOT EXISTS attendance_marks (
  user_id   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  cohort_id INTEGER NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
  day       TEXT NOT NULL,
  state     TEXT NOT NULL,            -- present | excused | absent
  reason    TEXT NOT NULL DEFAULT '',
  resolved_at TEXT NOT NULL DEFAULT (datetime('now')),
  PRIMARY KEY (user_id, day)
);

CREATE INDEX IF NOT EXISTS idx_attendance_cohort ON attendance_marks(cohort_id, day);
CREATE INDEX IF NOT EXISTS idx_attendance_state  ON attendance_marks(user_id, state, day);

-- --------------------------------------------------------------------------
-- Sanctions.
--
-- `shadow = 1` means "this is what the rule would have done" — recorded for
-- calibration, with no effect on the learner. Only a sanction with shadow = 0
-- has been applied. Nothing flips that column except an operator deliberately
-- enabling enforcement.
CREATE TABLE IF NOT EXISTS sanctions (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  cohort_id    INTEGER REFERENCES cohorts(id) ON DELETE SET NULL,
  kind         TEXT NOT NULL,          -- drop | suspend
  reason       TEXT NOT NULL DEFAULT '',
  window_from  TEXT NOT NULL DEFAULT '',  -- the absence run that triggered it
  window_to    TEXT NOT NULL DEFAULT '',
  shadow       INTEGER NOT NULL DEFAULT 1,
  applied_at   TEXT,                   -- NULL while shadow
  expires_at   TEXT,                   -- for suspensions
  created_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_sanctions_user   ON sanctions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sanctions_shadow ON sanctions(shadow, created_at DESC);

-- One sanction per learner per absence run, so re-running the evaluator does not
-- pile up duplicates for the same stretch of silence.
CREATE UNIQUE INDEX IF NOT EXISTS idx_sanctions_once
  ON sanctions(user_id, window_from, window_to);

-- --------------------------------------------------------------------------
-- Appeals. Present from the first migration rather than retrofitted: once the
-- platform makes consequential decisions about people, a route to contest them
-- is part of the feature, not a follow-up.
CREATE TABLE IF NOT EXISTS appeals (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  sanction_id   INTEGER NOT NULL REFERENCES sanctions(id) ON DELETE CASCADE,
  user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  body          TEXT NOT NULL DEFAULT '',
  state         TEXT NOT NULL DEFAULT 'open',   -- open | upheld | granted
  decided_by    INTEGER REFERENCES users(id) ON DELETE SET NULL,
  decided_at    TEXT,
  decision_note TEXT NOT NULL DEFAULT '',
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_appeals_one ON appeals(sanction_id);
CREATE INDEX IF NOT EXISTS idx_appeals_open ON appeals(state, created_at);

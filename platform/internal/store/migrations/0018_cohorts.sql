-- Phase 3 of the cohort program: groups, mentors, and the async daily standup.
--
-- Two decisions from the PRD shape this schema:
--
--   * Cohorts are sized for decay (§03). They form at 6–8 and are expected to
--     settle toward four, so nothing here assumes a fixed size — and a cohort of
--     one is a supported state, not an error, which is also the cold-start
--     answer.
--   * The standup is async with a daily window (§04). A fixed synchronous time
--     excludes a large share of a global audience, so each cohort carries a
--     timezone band and every standup stores its own open/close instants.

-- --------------------------------------------------------------------------
-- Bot mentors are real user rows so they can author standup posts through the
-- same tables as anyone else. The flag keeps them out of learner listings,
-- counts and sign-in.
ALTER TABLE users ADD COLUMN is_bot INTEGER NOT NULL DEFAULT 0;

-- --------------------------------------------------------------------------
-- A cohort is a group running one path together on a schedule.
--
-- `tz_band` is what makes the async standup workable: members share a band, so
-- the daily window overlaps everyone's waking hours.
CREATE TABLE IF NOT EXISTS cohorts (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  path_id     INTEGER NOT NULL REFERENCES paths(id) ON DELETE CASCADE,
  name        TEXT NOT NULL DEFAULT '',
  state       TEXT NOT NULL DEFAULT 'active',   -- active | completed | archived
  tz_band     TEXT NOT NULL DEFAULT 'europe_africa',
  starts_on   TEXT NOT NULL,                    -- YYYY-MM-DD
  created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_cohorts_open ON cohorts(state, path_id, tz_band);

-- Membership. Mentors (human or bot) sit in the same table as learners with a
-- different role, so "who is in this cohort" is one query.
--
-- left_at marks departure without deleting history: a dropped learner's standup
-- posts stay attributable.
CREATE TABLE IF NOT EXISTS cohort_members (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  cohort_id  INTEGER NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role       TEXT NOT NULL DEFAULT 'learner',   -- learner | mentor
  joined_at  TEXT NOT NULL DEFAULT (datetime('now')),
  left_at    TEXT
);

CREATE INDEX IF NOT EXISTS idx_members_cohort ON cohort_members(cohort_id, left_at);
CREATE INDEX IF NOT EXISTS idx_members_user   ON cohort_members(user_id, left_at);

-- One active membership per user per cohort.
CREATE UNIQUE INDEX IF NOT EXISTS idx_members_one_active
  ON cohort_members(cohort_id, user_id) WHERE left_at IS NULL;

-- --------------------------------------------------------------------------
-- Enrollments gain their cohort and the band used to place them. Added here
-- rather than in 0016 so the earlier migration stayed honest about what existed.
ALTER TABLE enrollments ADD COLUMN cohort_id INTEGER REFERENCES cohorts(id) ON DELETE SET NULL;
ALTER TABLE enrollments ADD COLUMN tz_band TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_enrollments_queue
  ON enrollments(state, path_id, tz_band) WHERE cohort_id IS NULL;

-- --------------------------------------------------------------------------
-- One standup per cohort per day. The window instants are stored rather than
-- derived at read time so that changing a band's hours later cannot silently
-- rewrite whether a past post counted as on time.
CREATE TABLE IF NOT EXISTS standups (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  cohort_id  INTEGER NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
  day        TEXT NOT NULL,          -- YYYY-MM-DD, the cohort's local day
  opens_at   TEXT NOT NULL,          -- UTC instant
  closes_at  TEXT NOT NULL,          -- UTC instant
  prompt     TEXT NOT NULL DEFAULT '',
  closed     INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_standups_day ON standups(cohort_id, day);
CREATE INDEX IF NOT EXISTS idx_standups_open ON standups(closed, closes_at);

-- A learner's post. The three fields are the standard standup shape; `blockers`
-- is the one the mentor acts on.
CREATE TABLE IF NOT EXISTS standup_entries (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  standup_id INTEGER NOT NULL REFERENCES standups(id) ON DELETE CASCADE,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  yesterday  TEXT NOT NULL DEFAULT '',
  today      TEXT NOT NULL DEFAULT '',
  blockers   TEXT NOT NULL DEFAULT '',
  posted_at  TEXT NOT NULL DEFAULT (datetime('now')),
  late       INTEGER NOT NULL DEFAULT 0   -- posted after closes_at
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_entries_one ON standup_entries(standup_id, user_id);
CREATE INDEX IF NOT EXISTS idx_entries_user ON standup_entries(user_id, posted_at DESC);

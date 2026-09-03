-- Phase 4 of the cohort program: the schedule. This is the change that actually
-- ends self-pacing — a cohort's path is expanded into dated work, and the
-- dashboard leads with what is due today rather than wherever you left off.
--
-- Two design notes:
--
--   * The schedule belongs to the COHORT, not the learner. One shared calendar
--     is what makes a group a cohort — everyone is on the same day, so standup
--     and peer help have common ground. Per-learner exemptions from the
--     diagnostic are applied as an overlay at read time (nothing is due for you
--     on a course you were exempted from), never by forking the schedule.
--
--   * Lesson completion did not exist before this migration. Only individual
--     lesson STEPS were tracked, so a theory lesson had no way to be finished
--     and "what's due today" could never be satisfied.

-- --------------------------------------------------------------------------
-- Lesson-level completion.
--
-- Written two ways: explicitly for a reading lesson ("mark as done"), and
-- automatically when the last step of a step-based lesson is completed.
CREATE TABLE IF NOT EXISTS lesson_progress (
  user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  lesson_id    INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
  completed_at TEXT NOT NULL DEFAULT (datetime('now')),
  PRIMARY KEY (user_id, lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_lesson_progress_user ON lesson_progress(user_id, completed_at);

-- --------------------------------------------------------------------------
-- The cohort's dated plan. One row per scheduled item.
--
-- `day_index` is 0-based from the cohort's start; `due_on` is the calendar date
-- it resolves to. Both are stored: the index survives a start-date change, and
-- the date is what every query actually filters on.
--
-- `sprint` groups days into weeks so the UI can say "Week 2 of 9" without
-- recomputing it. Weekends carry no rows — a schedule that demands seven days a
-- week is one nobody sustains, and it would make the absence rules unforgiving.
CREATE TABLE IF NOT EXISTS cohort_schedule (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  cohort_id     INTEGER NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
  day_index     INTEGER NOT NULL,
  sprint        INTEGER NOT NULL DEFAULT 1,
  due_on        TEXT NOT NULL,            -- YYYY-MM-DD
  kind          TEXT NOT NULL,            -- lesson | checkpoint
  course_id     INTEGER REFERENCES courses(id) ON DELETE CASCADE,
  lesson_id     INTEGER REFERENCES lessons(id) ON DELETE CASCADE,
  assignment_id INTEGER REFERENCES assignments(id) ON DELETE CASCADE,
  sort          INTEGER NOT NULL DEFAULT 0,
  created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_schedule_cohort ON cohort_schedule(cohort_id, due_on);
CREATE INDEX IF NOT EXISTS idx_schedule_day    ON cohort_schedule(cohort_id, day_index, sort);

-- One row per item per cohort — re-materializing must not duplicate the plan.
CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_lesson_once
  ON cohort_schedule(cohort_id, lesson_id) WHERE lesson_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_checkpoint_once
  ON cohort_schedule(cohort_id, assignment_id) WHERE assignment_id IS NOT NULL;

-- Marks a cohort as having had its plan built, so the job can skip it cheaply
-- rather than re-walking the whole path every run.
ALTER TABLE cohorts ADD COLUMN scheduled_at TEXT;

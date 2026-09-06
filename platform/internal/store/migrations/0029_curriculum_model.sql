-- Phase 3: make curriculum intent explicit and connect practice to courses.
-- Workload is measured in estimated learner minutes. Mode is one of
-- fun/theoretical/practical and is used to enforce the program's 15/30/55 mix.

ALTER TABLE courses ADD COLUMN primary_language TEXT NOT NULL DEFAULT '';
ALTER TABLE courses ADD COLUMN language_policy TEXT NOT NULL DEFAULT 'mixed'
  CHECK (language_policy IN ('mixed','c_only','c_with_exception'));
ALTER TABLE courses ADD COLUMN language_exception TEXT NOT NULL DEFAULT ''
  CHECK (language_policy != 'c_with_exception' OR length(trim(language_exception)) > 0);

ALTER TABLE lessons ADD COLUMN workload_minutes INTEGER NOT NULL DEFAULT 30 CHECK (workload_minutes > 0);
ALTER TABLE lessons ADD COLUMN mode TEXT NOT NULL DEFAULT 'theoretical' CHECK (mode IN ('fun','theoretical','practical'));
ALTER TABLE lessons ADD COLUMN language TEXT NOT NULL DEFAULT 'none';

ALTER TABLE assignments ADD COLUMN workload_minutes INTEGER NOT NULL DEFAULT 60 CHECK (workload_minutes > 0);
ALTER TABLE assignments ADD COLUMN mode TEXT NOT NULL DEFAULT 'practical' CHECK (mode IN ('fun','theoretical','practical'));

ALTER TABLE problems ADD COLUMN workload_minutes INTEGER NOT NULL DEFAULT 30 CHECK (workload_minutes > 0);
ALTER TABLE problems ADD COLUMN mode TEXT NOT NULL DEFAULT 'practical' CHECK (mode IN ('fun','theoretical','practical'));
ALTER TABLE problems ADD COLUMN language TEXT NOT NULL DEFAULT 'python';

CREATE TABLE IF NOT EXISTS course_problems (
  course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  sort INTEGER NOT NULL DEFAULT 0,
  required INTEGER NOT NULL DEFAULT 1 CHECK (required IN (0, 1)),
  PRIMARY KEY (course_id, problem_id)
);

CREATE INDEX IF NOT EXISTS idx_course_problems_order
  ON course_problems(course_id, sort);

-- The C checkpoints use new slugs rather than rewriting historical Python
-- assignments in place. Preserve their submissions and checks as archived rows.
UPDATE assignments
SET published = 0, required = 0
WHERE slug IN ('ds-hash-map','cx-pair-sum','algo-shortest-path','hcw-twos-complement');

ALTER TABLE cohort_schedule ADD COLUMN problem_id INTEGER REFERENCES problems(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_problem_once
  ON cohort_schedule(cohort_id, problem_id) WHERE problem_id IS NOT NULL;

-- Existing schedules are materialized snapshots. Rebuild Foundations once after
-- this migration so active cohorts drop Java and receive the new labs/problems.
DELETE FROM cohort_schedule
WHERE cohort_id IN (
  SELECT co.id FROM cohorts co
  JOIN paths p ON p.id = co.path_id
  WHERE p.slug = 'cs-foundations'
);
UPDATE cohorts
SET scheduled_at = NULL
WHERE path_id IN (SELECT id FROM paths WHERE slug = 'cs-foundations');

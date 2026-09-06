-- Phase 4 pacing enforcement.
--
-- A checkpoint has exactly one definition of "passed". Once any submission
-- meets the configured threshold it stays passed; a later retry cannot revoke
-- earned progress. Schedules, course gates, and path completion read this view.
CREATE VIEW IF NOT EXISTS passed_checkpoints AS
SELECT sub.user_id, sub.assignment_id, MIN(sub.id) AS submission_id
FROM submissions sub
JOIN assignments a ON a.id = sub.assignment_id
WHERE sub.status = 'graded'
  AND (a.pass_points = 0 OR (sub.score IS NOT NULL AND sub.score >= a.pass_points))
GROUP BY sub.user_id, sub.assignment_id;

-- Enrollment completion already uses ended_at. Cohorts need their own immutable
-- completion timestamp so a completed group is distinguishable from an archive.
ALTER TABLE cohorts ADD COLUMN completed_at TEXT;
CREATE INDEX IF NOT EXISTS idx_cohorts_completion ON cohorts(state, completed_at);

-- Assignments can act as a course checkpoint that gates progression.
--   required    = the learner must pass this to advance to the next course in a path
--   pass_points = minimum score to count as passed (0 = any graded submission passes)

ALTER TABLE assignments ADD COLUMN required INTEGER NOT NULL DEFAULT 0;
ALTER TABLE assignments ADD COLUMN pass_points INTEGER NOT NULL DEFAULT 0;

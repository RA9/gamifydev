-- Group lessons into named sections within a course, and give each lesson a
-- type so the curriculum reads like a structured syllabus.
--   section = the heading a lesson sits under (e.g. "HTML", "CSS"); "" = ungrouped
--   kind    = theory | lab | workshop | project (shown as a badge)

ALTER TABLE lessons ADD COLUMN section TEXT NOT NULL DEFAULT '';
ALTER TABLE lessons ADD COLUMN kind TEXT NOT NULL DEFAULT 'theory';

-- Optional base HTML for a JS step. For DOM labs, this markup is rendered in the
-- sandbox first and the learner's JavaScript runs against it (so they can select,
-- change, and wire events on existing elements). Ignored for html steps.

ALTER TABLE lesson_steps ADD COLUMN scaffold TEXT NOT NULL DEFAULT '';

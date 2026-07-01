-- A step's language selects how it runs and how its checks are evaluated:
--   html = live preview in a same-origin iframe; checks read the DOM via `doc`
--   js   = the learner's JS runs in a sandboxed (scripts-only) iframe; checks run
--          inside it with the learner's functions, `logs` (console output), `code`
--          (source), and `document` in scope, and post results back.

ALTER TABLE lesson_steps ADD COLUMN lang TEXT NOT NULL DEFAULT 'html';

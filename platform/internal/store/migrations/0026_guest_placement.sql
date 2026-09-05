-- Make placement diagnostics available to both signed-in learners and guests.
-- Every attempt and result has exactly one owner; guest-owned rows are rewritten
-- to user-owned rows atomically when a passing guest creates an account.

CREATE TABLE assessment_attempts_v26 (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id       INTEGER REFERENCES users(id) ON DELETE CASCADE,
  guest_id      INTEGER REFERENCES guest_sessions(id) ON DELETE CASCADE,
  assessment_id INTEGER NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
  started_at    TEXT NOT NULL DEFAULT (datetime('now')),
  expires_at    TEXT NOT NULL,
  submitted_at  TEXT,
  score         INTEGER,
  topic_scores  TEXT NOT NULL DEFAULT '{}',
  abandoned     INTEGER NOT NULL DEFAULT 0,
  CHECK ((user_id IS NULL) <> (guest_id IS NULL))
);

INSERT INTO assessment_attempts_v26
  (id, user_id, guest_id, assessment_id, started_at, expires_at,
   submitted_at, score, topic_scores, abandoned)
SELECT id, user_id, NULL, assessment_id, started_at, expires_at,
       submitted_at, score, topic_scores, abandoned
FROM assessment_attempts;

DROP TABLE assessment_attempts;
ALTER TABLE assessment_attempts_v26 RENAME TO assessment_attempts;

CREATE INDEX idx_attempts_user
  ON assessment_attempts(user_id, started_at DESC);
CREATE INDEX idx_attempts_guest
  ON assessment_attempts(guest_id, started_at DESC);

-- A sitting owns an immutable copy of every drawn item. The bank row id remains
-- the response key, but no read or grade depends on that mutable bank row.
CREATE TABLE assessment_attempt_items_v26 (
  attempt_id  INTEGER NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
  item_id     INTEGER NOT NULL,
  sort        INTEGER NOT NULL DEFAULT 0,
  response    INTEGER,
  correct     INTEGER,
  topic       TEXT NOT NULL DEFAULT '',
  difficulty  INTEGER NOT NULL DEFAULT 1,
  prompt      TEXT NOT NULL DEFAULT '',
  code        TEXT NOT NULL DEFAULT '',
  lang        TEXT NOT NULL DEFAULT '',
  options     TEXT NOT NULL DEFAULT '[]',
  answer      INTEGER NOT NULL DEFAULT 0,
  explanation TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (attempt_id, item_id)
);

INSERT INTO assessment_attempt_items_v26
  (attempt_id, item_id, sort, response, correct, topic, difficulty,
   prompt, code, lang, options, answer, explanation)
SELECT ai.attempt_id, ai.item_id, ai.sort, ai.response, ai.correct,
       COALESCE(i.topic, ''), COALESCE(i.difficulty, 1),
       COALESCE(i.prompt, ''), COALESCE(i.code, ''), COALESCE(i.lang, ''),
       COALESCE(i.options, '[]'), COALESCE(i.answer, 0),
       COALESCE(i.explanation, '')
FROM assessment_attempt_items ai
LEFT JOIN assessment_items i ON i.id = ai.item_id;

DROP TABLE assessment_attempt_items;
ALTER TABLE assessment_attempt_items_v26 RENAME TO assessment_attempt_items;

-- The pass bit is an explicit policy decision. Existing account-owned results
-- are grandfathered as passing: those learners were already admitted under the
-- previous routing-only contract and must not be locked out retroactively.
CREATE TABLE placement_results_v26 (
  id                   INTEGER PRIMARY KEY AUTOINCREMENT,
  attempt_id           INTEGER NOT NULL UNIQUE REFERENCES assessment_attempts(id) ON DELETE CASCADE,
  user_id              INTEGER REFERENCES users(id) ON DELETE CASCADE,
  guest_id             INTEGER REFERENCES guest_sessions(id) ON DELETE CASCADE,
  recommended_path_id  INTEGER REFERENCES paths(id) ON DELETE SET NULL,
  foundations_required INTEGER NOT NULL DEFAULT 1,
  passed               INTEGER NOT NULL DEFAULT 0 CHECK (passed IN (0, 1)),
  exemptions           TEXT NOT NULL DEFAULT '[]',
  created_at           TEXT NOT NULL DEFAULT (datetime('now')),
  CHECK ((user_id IS NULL) <> (guest_id IS NULL))
);

INSERT INTO placement_results_v26
  (id, attempt_id, user_id, guest_id, recommended_path_id,
   foundations_required, passed, exemptions, created_at)
SELECT id, attempt_id, user_id, NULL, recommended_path_id,
       foundations_required, 1, exemptions, created_at
FROM (
  SELECT pr.*,
         ROW_NUMBER() OVER (
           PARTITION BY attempt_id
           ORDER BY datetime(created_at) DESC, id DESC
         ) AS result_rank
  FROM placement_results pr
) AS ranked
WHERE result_rank = 1;

DROP TABLE placement_results;
ALTER TABLE placement_results_v26 RENAME TO placement_results;

CREATE INDEX idx_placement_user
  ON placement_results(user_id, created_at DESC);
CREATE INDEX idx_placement_guest
  ON placement_results(guest_id, created_at DESC);

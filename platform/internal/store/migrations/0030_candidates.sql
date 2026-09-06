-- Candidate identity for the placement diagnostic.
--
-- Until now a sitting was owned by a guest cookie, which means the retry
-- cooldown, the three-attempt limit and a sitting voided for cheating all reset
-- the moment somebody opened a private window. The gate was fitted to the one
-- thing the candidate controls completely.
--
-- A candidate is a person as identified by their email address, and the gating
-- reads through it. The email is not verified, so this does not stop somebody
-- determined to invent a second identity — it stops the accidental and the
-- casual bypass, and it gives verification somewhere to attach later.
CREATE TABLE IF NOT EXISTS candidates (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  -- Normalized: lower-cased and stripped of any +tag, so alice+1@x.com and
  -- Alice@x.com are one person rather than two fresh sets of attempts.
  email      TEXT NOT NULL UNIQUE,
  -- As typed, because that is the address to actually write to.
  email_raw  TEXT NOT NULL,
  name       TEXT NOT NULL,
  -- Set once they hold an account, so a later sitting sees their enrollment.
  user_id    INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  last_seen_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_candidates_user ON candidates(user_id);

ALTER TABLE assessment_attempts ADD COLUMN candidate_id INTEGER REFERENCES candidates(id);
ALTER TABLE placement_results  ADD COLUMN candidate_id INTEGER REFERENCES candidates(id);
CREATE INDEX IF NOT EXISTS idx_attempts_candidate ON assessment_attempts(candidate_id, started_at);
CREATE INDEX IF NOT EXISTS idx_results_candidate ON placement_results(candidate_id, created_at);

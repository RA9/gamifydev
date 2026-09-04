-- Password reset tokens.
--
-- token_hash, not the token: a reset token is a bearer credential for taking
-- over an account, so the database stores only its SHA-256. A dump of this
-- table then yields nothing usable, which is not true of invitations(token)
-- and is the difference worth having here.
--
-- Rows are kept after use rather than deleted so a consumed token stays
-- consumed, and so a support question ("did that reset get used?") has an
-- answer.
CREATE TABLE IF NOT EXISTS password_resets (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  used_at    TEXT,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_password_resets_hash ON password_resets(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_resets_user ON password_resets(user_id);

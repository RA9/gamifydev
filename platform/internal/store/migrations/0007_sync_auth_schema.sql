-- Ensure the sessions table exists. On a database whose users table predates the
-- platform (so migration 0001's CREATE TABLE IF NOT EXISTS users was a no-op),
-- its companion sessions table may never have been created.

CREATE TABLE IF NOT EXISTS sessions (
  token      TEXT PRIMARY KEY,
  user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);

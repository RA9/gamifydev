-- Sync the auth schema on databases whose `users` table predates the platform
-- (e.g. a Turso/libSQL database previously used by another app). When `users`
-- already exists, migration 0001's CREATE TABLE IF NOT EXISTS is a no-op, so the
-- table can be missing the platform's columns (registration then fails with
-- "no such column: password_hash") and a companion table like `sessions` may
-- never have been created.
--
-- This migration is idempotent: CREATE ... IF NOT EXISTS no-ops where objects
-- already exist, and the migration runner skips "duplicate column" errors (SQLite
-- has no ADD COLUMN IF NOT EXISTS), so each ALTER adds the column where missing
-- and is silently skipped where it already exists.

-- 1. Ensure the core auth tables/indexes exist (no-op on a healthy database).
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL DEFAULT '',
    name          TEXT NOT NULL DEFAULT '',
    role          TEXT NOT NULL DEFAULT 'learner',
    bio           TEXT NOT NULL DEFAULT '',
    avatar_url    TEXT NOT NULL DEFAULT '',
    created_at    TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);

-- 2. Add any columns missing from a pre-existing `users` table. Constant
--    defaults are required: SQLite rejects expression defaults like
--    datetime('now') in ALTER TABLE ADD COLUMN.
ALTER TABLE users ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'learner';
ALTER TABLE users ADD COLUMN bio TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN avatar_url TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN created_at TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN updated_at TEXT NOT NULL DEFAULT '';

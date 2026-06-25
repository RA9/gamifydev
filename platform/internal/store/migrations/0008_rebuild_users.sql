-- Rebuild the users table to the exact canonical shape.
--
-- Some databases (e.g. a Turso/libSQL database previously used by another app)
-- have a users table carrying legacy columns the platform never populates — for
-- example a NOT NULL `fullname` — so INSERTs fail with
-- "NOT NULL constraint failed: users.fullname". Migration 0007 added the
-- platform's columns, but SQLite cannot drop columns or relax a NOT NULL
-- constraint in place, so the only fix is to recreate the table.
--
-- This selects only the platform columns (which 0007 guarantees exist on every
-- database) into a fresh table and drops the old one along with any legacy
-- columns. id and email are preserved.
--
-- The migration runner applies this with foreign_keys=OFF and
-- legacy_alter_table=ON (set on the connection, outside the transaction), so the
-- RENAME leaves child tables' "REFERENCES users" pointing at the rebuilt table
-- and the old table can be dropped without tripping foreign keys.

ALTER TABLE users RENAME TO users_rebuild_0008;

CREATE TABLE users (
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

INSERT INTO users (id, email, password_hash, name, role, bio, avatar_url, created_at, updated_at)
SELECT id, email, password_hash, name, role, bio, avatar_url, created_at, updated_at
FROM users_rebuild_0008;

DROP TABLE users_rebuild_0008;

-- Rebuild the users table to the exact canonical shape.
--
-- Some databases (e.g. a Turso/libSQL database previously used by another app)
-- have a users table carrying legacy columns the platform never populates — for
-- example a NOT NULL `fullname` — so INSERTs fail with
-- "NOT NULL constraint failed: users.fullname". Migration 0007 added the
-- platform's columns, but SQLite cannot drop columns or relax a NOT NULL
-- constraint in place, so the only fix is to recreate the table.
--
-- Strategy: build the canonical table under a temporary name, copy the platform
-- columns (which 0007 guarantees exist on every database), drop the old table,
-- then rename the new one into place. Child tables reference `users` by name the
-- whole time and never reference the temporary name, so no foreign key is ever
-- rewritten — after the final rename they bind to the rebuilt table. This needs
-- only foreign_keys=OFF (the runner sets it), which Turso/libSQL allows, unlike
-- PRAGMA legacy_alter_table. id and email are preserved.

CREATE TABLE users_new_0008 (
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

INSERT INTO users_new_0008 (id, email, password_hash, name, role, bio, avatar_url, created_at, updated_at)
SELECT id, email, password_hash, name, role, bio, avatar_url, created_at, updated_at
FROM users;

DROP TABLE users;

ALTER TABLE users_new_0008 RENAME TO users;

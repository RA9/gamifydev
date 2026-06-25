-- Recreate the users table to the canonical schema.
--
-- Some databases (e.g. a Turso/libSQL database previously used by another app)
-- have a users table with incompatible legacy columns — for example a NOT NULL
-- `fullname` — so INSERTs fail ("NOT NULL constraint failed: users.fullname").
-- SQLite cannot drop a column or relax a NOT NULL constraint in place, and legacy
-- rows can't authenticate under the platform's bcrypt auth anyway, so the table
-- is dropped and recreated empty. The runner applies migrations with foreign
-- keys off, so dropping a table that child tables reference is safe; after the
-- table is recreated, those "REFERENCES users" foreign keys resolve to it again.

DROP TABLE IF EXISTS users;

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

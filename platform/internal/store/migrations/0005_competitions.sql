-- Coding competitions: timed contests with a leaderboard. One entry per user
-- per competition (updatable while the contest is live); scored by staff.

CREATE TABLE IF NOT EXISTS competitions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    slug       TEXT NOT NULL UNIQUE,
    title      TEXT NOT NULL,
    prompt     TEXT NOT NULL DEFAULT '',
    language   TEXT NOT NULL DEFAULT 'text',
    points     INTEGER NOT NULL DEFAULT 100,
    starts_at  TEXT NOT NULL,
    ends_at    TEXT NOT NULL,
    published  INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS competition_entries (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    competition_id INTEGER NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    user_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code           TEXT NOT NULL DEFAULT '',
    score          INTEGER NOT NULL DEFAULT 0,
    created_at     TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at     TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (competition_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_entries_comp ON competition_entries(competition_id, score);

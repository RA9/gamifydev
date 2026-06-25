# Database Migrations

Plain-SQL schema migrations for the GamifyDev platform. The same SQL runs against
a local SQLite file in development and Turso/libSQL in production.

## How it works

Migrations run automatically on every server start (`store.Migrate`, called from
`main.go`). They are embedded into the binary with `//go:embed migrations/*.sql`,
so there are no external files to ship. The runner:

- creates a `schema_migrations` table to track what has been applied;
- loads all `.sql` files from this directory and sorts them by filename;
- skips migrations already recorded in `schema_migrations`;
- runs each remaining migration **inside a transaction** — if a statement fails,
  the whole migration is rolled back and startup aborts;
- **skips "already exists" / "duplicate column" errors** so idempotent statements
  (notably `ALTER TABLE ADD COLUMN`) can run safely against a database that
  already has the change.

That last point is what keeps an older or shared database in sync. SQLite has no
`ADD COLUMN IF NOT EXISTS`, so a "sync" migration simply `ALTER TABLE ADD
COLUMN`s the expected columns: it adds them where missing and no-ops where they
already exist. See `0007_sync_auth_schema.sql` for an example.

## File naming

```
{version}_{description}.sql
```

- Version numbers are **zero-padded to 4 digits** (`0001`, `0002`, …) to match
  the existing files and keep filename sort order correct.
- Version numbers must be unique and must not be reused.
- Use underscores in the description; the extension must be `.sql`.

Examples:

- `0001_init.sql`
- `0006_blog.sql`
- `0007_sync_auth_schema.sql`

## Creating a migration

1. Add a new file with the next version number, e.g.
   `0008_add_user_streaks.sql`.
2. Write standard SQLite-compatible SQL (see do's/don'ts below).
3. Rebuild/restart the server — it applies automatically on boot. To apply to a
   specific database without a full deploy:
   ```bash
   DATABASE_URL='libsql://<db>.turso.io?authToken=<token>' go run ./cmd/seedcontent
   ```
   (the seed command runs `Migrate` before seeding), or set `SEED=1` on the
   Railway service for one boot.

## Do / Don't

**Do**

- Use `IF NOT EXISTS` on `CREATE TABLE` / `CREATE INDEX`.
- Use `IF EXISTS` on `DROP`.
- For NOT NULL `ADD COLUMN`, provide a **constant** default — SQLite rejects
  expression defaults like `datetime('now')` in `ALTER TABLE ADD COLUMN`.
- Keep migrations small and focused; add indexes for foreign keys and
  frequently-queried columns.

**Don't**

- Don't edit a migration after it has been applied anywhere — add a new one.
- Don't delete migration files or reuse version numbers.
- Don't rely on database-specific syntax beyond what SQLite/libSQL share.

## SQLite limitations

SQLite can't drop or retype columns directly. For those, use the rebuild pattern:

```sql
ALTER TABLE old_table RENAME TO old_table_backup;
CREATE TABLE old_table ( ... new schema ... );
INSERT INTO old_table (col_a, col_b) SELECT col_a, col_b FROM old_table_backup;
DROP TABLE old_table_backup;
```

Note: foreign keys cannot be added to an existing table, and `DROP TABLE` may
need foreign-key enforcement disabled depending on the connection.

## Checking status

```sql
SELECT name, applied_at FROM schema_migrations ORDER BY name;
```

## Rollback

There is no automatic rollback. To undo a change, write a new migration with a
higher version number that reverses it.

## Production notes

- Back up the database before deploying schema changes.
- Watch the server logs on boot: each applied migration logs
  `migrate: applied <name>`, and skipped idempotent statements log
  `migrate <name>: skipping (...)`.

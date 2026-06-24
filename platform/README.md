# GamifyDev Platform (Go)

The next-generation GamifyDev: a server-rendered Go application — a full coding
academy with courses (text/image/video/audio), human-graded coding assignments,
coding competitions, and a community forum.

## Stack

- **Go** (`net/http`, `html/template`) — server-rendered HTML, single binary
  (templates, static assets, and migrations are embedded via `embed`).
- **[tan-compose](https://ra9.github.io/tan-compose/)** — interactive client-side
  web components, dropped into pages where rich interactivity is needed.
- **htmx** — server-driven partial updates (e.g. live role changes in the admin
  user table).
- **Turso / libSQL** — the database. Dev uses a local SQLite file; production
  uses Turso. Same SQL either way; no ORM, plain migrations in
  `internal/store/migrations`.

## Run locally

```bash
cd platform
go run .            # serves http://localhost:8080 using a local gamifydev.db
```

The **first account you register becomes the admin** (so the admin portal is
reachable); everyone after is a learner. Roles: `admin`, `grader`, `learner`.

## Configuration

See `.env.example`. Key vars: `DATABASE_URL` (local file or `libsql://…` Turso
URL), `ADDR`, `SECURE_COOKIES`.

## Layout

```
main.go                         # wiring: db, migrations, server
internal/store/                 # data layer (database/sql) + SQL migrations
internal/auth/                  # password hashing, sessions, request context
internal/server/                # routes, middleware, handlers, rendering
  web/templates/                # html/template pages + shared layout (embedded)
  web/static/                   # css / js assets (embedded)
```

## Status

Foundation complete: accounts + sessions + roles, unified layout, and the admin
portal with user management (htmx live role changes). Next: course authoring &
delivery with media, the assignment submission + mentor grading workflow,
competitions, and the forum.

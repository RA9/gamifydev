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
URL), `ADDR`, `SECURE_COOKIES`. In production the server also honours `PORT`
(it binds `:$PORT` when set).

## Deploy to Railway

The app is containerised (`Dockerfile`) and reads `PORT` from Railway. Railway's
filesystem is ephemeral, so production uses **Turso** for the database.

1. **Create a Turso database** and grab its URL + token:
   ```bash
   curl -sSfL https://get.tur.so/install.sh | bash   # if you don't have the CLI
   turso auth signup
   turso db create gamifydev
   turso db show gamifydev --url        # -> libsql://gamifydev-<org>.turso.io
   turso db tokens create gamifydev     # -> the auth token
   ```

2. **Create the Railway service** from this GitHub repo. Because the app lives in
   a subdirectory, set the service's **Root Directory** to `platform` — Railway
   then uses `platform/Dockerfile` and `platform/railway.json` automatically.

3. **Set the service variables** (Railway → Variables):
   - `DATABASE_URL` = `libsql://<your-db>.turso.io?authToken=<token>`
   - `PORT` is provided by Railway; `SECURE_COOKIES` is auto-enabled there.

4. **Deploy.** Railway builds the image, runs migrations on boot, and serves on
   the generated domain (health-checked at `/healthz`). The **first account you
   register becomes the admin**.

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

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
   - *(Optional)* add a **Redis** service in Railway and set `REDIS_URL` to its
     connection string (e.g. `${{Redis.REDIS_URL}}`). The app runs without it
     today; future features (caching, rate limiting, live competitions) use it.

4. **Deploy.** Railway builds the image, runs migrations on boot, and serves on
   the generated domain (health-checked at `/healthz`). The **first account you
   register becomes the admin**.

## Layout

```
main.go                         # wiring: db, migrations, server
internal/store/                 # data layer (database/sql) + SQL migrations
internal/auth/                  # password hashing, sessions, request context
internal/rdb/                   # optional Redis client (REDIS_URL)
internal/server/                # routes, middleware, handlers, rendering
  web/templates/                # html/template pages + shared layout (embedded)
  web/static/                   # css / js assets (embedded)
```

## Accounts and sign-in

The first account to register becomes the admin, so a fresh deploy is
reachable without touching the database.

**Password reset** is at `/forgot`. Tokens are stored as a SHA-256 hash, last an
hour, work once, and requesting a new one invalidates the old. Completing a
reset deletes every session that user has — the point of a reset is usually that
someone else has one. When `SMTP_HOST` is unset the link goes to the server log
for an operator to retrieve; it is never shown in the browser, since the form
accepts any address.

**Failed sign-ins and reset requests are throttled** — five inside fifteen
minutes trips a fifteen-minute cool-off. Two keys are counted: the email being
tried and the caller's address, so neither one account under sustained attack
nor one host spraying many accounts gets through. A correct password clears the
count. `X-Forwarded-For` is trusted only when the app knows it is behind a proxy
(the signal that also enables Secure cookies).

Both `/login` and `/forgot` answer identically whether or not an account exists,
including timing — a miss burns the same bcrypt work a real comparison would, so
the form can't be used to enumerate who is registered.

## Checkpoints

A checkpoint is an assignment marked **required**: passing it opens the next
course in the path. Each carries a list of *checks*, authored at
`/admin/assignments/{id}`, and each check is a Python boolean expression
whatever the assignment's language:

- **Python** — the expression runs in the namespace the learner's program left
  behind, so it can call their functions directly (`sum_list([]) == 0`), plus
  `_code` for their source.
- **C and shell** — there is no namespace to inspect, so the program runs first
  and the expression asserts over what it produced: `_out`, `_err`, `_exit`,
  `_in` (the stdin that check supplied) and `_code`. Checks sharing a stdin run
  together, so one program run serves all of them.
- **Anything else** (Java, JavaScript) has no sandbox yet and falls to a mentor.
  The admin form says which of these applies before you write a single check.

A check worth 1+ point must pass for the submission to be graded; one worth 0 is
advisory. Marking a check **hidden** withholds its wording until the learner
passes it, which is what stops the spec being read off the failure list — and is
how a checkpoint catches an answer that was hardcoded rather than computed.

**Fixture files** are written into the sandbox beside the program and thrown
away after, which is what makes "do something with these files" possible. They
are shown on the assignment page too, so a learner can read what they're being
asked to process.

Checks are run by the `checks:run` job, not inline on submit, so a cohort all
submitting the same evening is metered rather than stampeding the runner. Each
check reports as it finishes, so a submission killed by the time limit still
shows what it had passed.

A gate holds back **every course after it in the path**, so a mid-path gate has
to be one the sandbox can grade. A mentor-graded gate at position two puts one
person's availability in front of everything behind it. Where a course's
language has no sandbox, mark its checkpoint `elective` — it stays available for
practice and feedback without blocking anyone. `TestNoMidPathGateWaitsOnAHuman`
enforces this; the last course in a path is exempt, since its gate holds nothing
back.

Seeded checkpoints live in `internal/seed/seed.go`. Every one of them has a
reference solution in `internal/seed/checkpoints_test.go` that must pass all its
checks, and a plausible wrong solution that must fail a named one — a gate
nothing can fail, or that nothing can pass, is caught by `go test` rather than
by a learner.

## Status

Foundation complete: accounts + sessions + roles, unified layout, and the admin
portal with user management (htmx live role changes). Next: course authoring &
delivery with media, the assignment submission + mentor grading workflow,
competitions, and the forum.

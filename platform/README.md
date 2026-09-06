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

Learners take and pass the public placement test before creating an account.
Public signup always creates a `learner`; it can never create an administrator.
On startup, the server creates the first operator account for
`cnah27@gmail.com`. Set a password of at least 12 characters before the first
boot:

```bash
ADMIN_PASSWORD='use-a-long-password' go run .
```

The startup seed is idempotent. Later boots ensure the account still has the
`admin` role without creating a duplicate or resetting its password, so
`ADMIN_PASSWORD` is only required while that account does not exist. The
`cmd/createadmin` operator command remains available to create or promote a
different account explicitly:

```bash
ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD='use-a-long-password' go run ./cmd/createadmin
```

Roles are `admin`, `grader`, and `learner`; staff invitations are available from
the admin portal after bootstrap.

## Configuration

See `.env.example`. Key vars: `DATABASE_URL` (local file or `libsql://…` Turso
URL), `ADMIN_PASSWORD` (at least 12 characters for first-admin creation), `ADDR`,
and `SECURE_COOKIES`. In production the server also honours `PORT` (it binds
`:$PORT` when set).

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
   - `ADMIN_PASSWORD` = a strong password of at least 12 characters for
     `cnah27@gmail.com`; startup fails clearly if the account must be created and
     this is missing or too short.
   - `PORT` is provided by Railway; `SECURE_COOKIES` is auto-enabled there.
   - *(Optional)* add a **Redis** service in Railway and set `REDIS_URL` to its
     connection string (e.g. `${{Redis.REDIS_URL}}`). The app runs without it
     today; future features (caching, rate limiting, live competitions) use it.

4. **Deploy.** Railway builds the image, runs migrations, ensures the first
   administrator, synchronizes seed content, and serves on the generated domain
   (health-checked at `/healthz`). Public registration never grants staff roles.

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

Account creation is the final step of placement: a guest must hold a passing,
unclaimed placement result before `GET` or `POST /register` succeeds. The result,
assessment history, any guest problem submissions, and an unplaced enrollment are
transferred to the new learner atomically. Enrollment then binds the learner to
the exact generated result: its path, exemption snapshot, and selected timezone
region are committed together. The snapshot—not any later placement row—drives
the cohort schedule, and queued, active, or paused enrollments lock placement
retakes so an active schedule cannot be rewritten.

The fixed first administrator is ensured during startup; additional staff
accounts are created with `cmd/createadmin` or invited by an existing
administrator.

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

## Curriculum model

Every lesson, guided problem, and checkpoint carries an estimated workload and a
learning mode: `fun`, `theoretical`, or `practical`. The seeded Computer Science
Foundations path totals 2,000 minutes and is enforced at **15% fun, 30%
theoretical, and 55% practical**. Labs are step-based lessons embedded in each
course, while selected problems from the public practice bank are linked directly
from the course they reinforce.

Foundations uses **C** for all general programming material, labs, integrated
problems, and checkpoints. The only language exception is the Linux course,
where command-line labs and its operations checkpoint use shell because shell
interaction is itself the subject being taught. Java is not part of the
Foundations path. Seed policy tests verify the workload ratio, integrations,
checkpoint languages, step languages, and Markdown code fences.

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

## Tests

```bash
go test ./...
```

The policy packages (`placement`, `cohort`, `schedule`, `attendance`) are unit
tested against their rules. `runner` and `pyharness` are tested against real
sandboxed execution, and `seed` proves every shipped checkpoint is solvable.

`internal/server` has smoke tests: they drive the real mux, real templates and a
real seeded database. Two of them read the routing table out of `server.go` and
assert that **every** non-public route turns away a guest and every `/admin`
route turns away a learner — so a handler registered without its middleware
fails the build, including one added long after this was written. That is the
one class of bug in the HTTP layer that nothing else would catch.

They are smoke tests, not thorough ones: they check that pages come back and
that the locked doors are locked, not that every branch inside a handler is
right.

## Status

Foundation complete: accounts + sessions + roles, unified layout, and the admin
portal with user management (htmx live role changes). Next: course authoring &
delivery with media, the assignment submission + mentor grading workflow,
competitions, and the forum.

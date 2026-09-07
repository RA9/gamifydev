# runner — server-side code execution

This package runs short, **untrusted** programs server-side so the platform can
offer work the in-browser Pyodide sandbox can't: real `input()`, a real
filesystem scratch, and (in the sandboxed container) installed packages. It
powers the `/playground` page, checkpoint grading, and the practice problem
judge.

## Languages

| Lang | How it runs | Notes |
| --- | --- | --- |
| `python` | interpreted, `python3 -I -B` | no compile step |
| `c` | `cc -std=c11 -O0 -Wall`, then the binary | warnings are shown; they are teaching material |
| `go` | `go build`, then the binary | standard library only — the sandbox has no network, so an import needing a download fails on the proxy rather than on the code |
| `shell` | `bash main.sh` | for the Linux course |

Two things make compiled languages workable rather than merely possible.

**The build cache is shared across runs.** Go compiles the standard library on a
cold cache — around twenty seconds, against roughly one when warm — so a
per-run cache would time out every Go submission. Both caches are keyed by a
hash of the build inputs, so one submission cannot make another's compile
resolve to something it planted. Call `WarmGo` at boot; otherwise the first Go
submission after a deploy pays for the cold build and dies on its own limit.

**The compiler gets its own time budget.** The learner's limit is about their
algorithm; a compile has nothing to do with them. `buildAllowance` adds a bounded
extra to the wall clock for compiled languages, so a problem whose limit suits a
correct Python answer does not fail the same answer in Go on build time alone.

Running code a stranger typed is dangerous. This package treats that seriously:
it is **disabled by default** and never pretends a bare subprocess is a security
boundary.

## Modes

Selected by the `CODE_EXEC` env var on the main app:

| `CODE_EXEC` | What runs | Isolation | Use it for |
|-------------|-----------|-----------|------------|
| unset / `off` | nothing (`ErrDisabled`) | — | default |
| `local` | a subprocess on the app host | rlimits + timeout + process-group kill, **plus** a bubblewrap namespace sandbox *if `bwrap` is on PATH* | local dev; self-hosting where `bwrap` is present |
| `remote` | forwarded to the `execd` service | whatever `execd`'s container provides | **production** |

### Secure-by-default rules (`runner.New`)

- Unknown/empty mode → disabled executor (never an error).
- `remote` with no `CODE_EXEC_SANDBOX_URL` → startup error.
- `local` **in a deployed environment** (`RAILWAY_ENVIRONMENT` set) **without**
  `bwrap` → startup error, unless `CODE_EXEC_UNSAFE=1` is set to acknowledge the
  risk explicitly.

## Every run is bounded

- **Wall-clock timeout** (default 5s, max 10s); the whole process group is
  SIGKILLed on expiry, so `while True: pass` dies cleanly.
- **CPU / memory / file-size / process-count** caps via `ulimit` (Linux).
- **Output cap** per stream (64 KB); excess is dropped and flagged `Truncated`.
- **Scrubbed environment** — no inherited secrets — and a fresh per-run temp
  dir that's deleted afterward.
- **`python -I -B`** (isolated mode, no bytecode, no ambient `PYTHON*` / cwd on
  `sys.path`).

The HTTP endpoint (`POST /api/run`) adds auth (login required), a per-user
minimum interval, a global concurrency cap, and a 64 KB program-size limit.

## Boot self-test (proves the sandbox actually works)

Configuring bwrap isn't the same as bwrap *working* — unprivileged user
namespaces can be disabled by the host kernel or a container runtime, in which
case bwrap silently degrades. So `execd` runs a **containment self-test at
boot**: it executes a probe through its own sandbox and checks that outbound
network is blocked, the host filesystem (`/etc/hosts`) is hidden, and the run is
under the bwrap mount. With `EXECD_REQUIRE_SANDBOX=1` (the container default) a
failed self-test is **fatal — execd refuses to serve** rather than run untrusted
code unsandboxed. It fails closed: a timeout or unparseable probe never reads as
safe. The same probe backs `execd selftest` (exit 0 = sandboxed, 1 = not),
wired as the container `HEALTHCHECK`, so an unhealthy `execd` means "the sandbox
isn't containing code on this host."

## One-command deploy (execution enabled)

`docker-compose.yml` stands up the app plus a sandboxed `execd`, wired with a
shared secret:

```bash
echo "EXEC_TOKEN=$(openssl rand -hex 32)" > .env
echo "DATABASE_URL=libsql://<db>.turso.io?authToken=<t>" >> .env
docker compose up --build
```

The app comes up with `CODE_EXEC=remote` pointed at `execd`; `execd` only
becomes healthy once its self-test passes. If your host has unprivileged user
namespaces disabled, `execd` stays unhealthy (fail closed) — enable them
(`sysctl kernel.unprivileged_userns_clone=1`) or run `execd` under gVisor/Kata.

## What `bwrap` adds (and why remote is the prod answer)

Limits and timeouts stop resource abuse, but a *plain* subprocess can still read
the host filesystem and reach the network. `bwrap` (bubblewrap) closes that: it
runs Python in an unprivileged namespace with **no network**, a **read-only
root**, a private `/tmp`, and only the per-run scratch dir writable.

The production image (`distroless/static`) has no Python at all, so `local` mode
can't even run there. The intended topology is:

```
main app (distroless)  --CODE_EXEC=remote-->  execd (Dockerfile.execd: python + bwrap)
```

The blast radius of a sandbox escape is then a disposable `execd` container, not
the app that holds the database credentials.

### Deploy `execd`

Build `Dockerfile.execd` and run it locked down (see the header of that file for
the exact `docker run` flags: `--read-only`, `--cap-drop ALL`,
`--security-opt no-new-privileges`, `--pids-limit`, `--memory`, `--cpus`). Set:

- `EXECD_TOKEN` — shared secret; the app sends it as `Authorization: Bearer …`.
- `EXECD_MAX_CONCURRENCY` — simultaneous runs (default 4).

Then on the app: `CODE_EXEC=remote`, `CODE_EXEC_SANDBOX_URL=http://execd:9090`,
`CODE_EXEC_SANDBOX_TOKEN=$EXECD_TOKEN`.

> If the host disables unprivileged user namespaces, `bwrap` can't isolate and
> `execd` logs a warning. There, run the `execd` container itself under a
> stronger boundary (gVisor, Kata, or a dedicated VM).

## Threat model, briefly

In scope: resource exhaustion (CPU/mem/time/output/pids), secret exfiltration
from the environment, filesystem tampering outside the scratch dir, and — with
`bwrap` or `remote` — network egress and host filesystem reads. Out of scope for
plain `local` mode without `bwrap` (hence the dev-only gating): kernel exploits,
network access, and host filesystem reads. Do not expose plain `local` mode to
the public internet.

# runner — server-side Python execution

This package runs short, **untrusted** Python snippets server-side so the
platform can offer labs the in-browser Pyodide sandbox can't: real `input()`, a
real filesystem scratch, and (in the sandboxed container) installed packages.
It powers the `/playground` page and is the foundation for future server-run
labs (Phase 3: web/servers).

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

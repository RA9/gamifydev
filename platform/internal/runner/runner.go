// Package runner executes short, untrusted snippets of Python server-side and
// returns their output. It exists so the platform can offer labs that the
// in-browser Pyodide sandbox cannot — real stdin, a real filesystem scratch,
// and (via the standalone executor service) installed packages.
//
// SECURITY, READ THIS FIRST.
//
// Running code a stranger typed is one of the most dangerous things a web app
// can do. This package does NOT pretend a bare subprocess is a security
// boundary. It offers three modes, chosen by config, and is DISABLED by
// default:
//
//   - "off"    — every Run returns ErrDisabled. This is the default.
//   - "local"  — run in a subprocess on THIS host. Hardened (isolated Python,
//     scrubbed env, wall-clock + CPU + memory + output limits, whole
//     process-group kill) and, when the `bwrap` (bubblewrap) binary
//     is present, wrapped in an unprivileged namespace sandbox with
//     NO network and a read-only root. Without bwrap it is only a
//     convenience for local development, NOT safe for the public
//     internet, and New() refuses to enable it in a deployed
//     environment unless CODE_EXEC_UNSAFE=1 is set explicitly.
//   - "remote" — forward each request to a standalone executor service
//     (cmd/execd) that runs in its own hardened container. This is
//     the intended production topology; the main app image stays
//     lean (distroless, no Python) and the blast radius of a sandbox
//     escape is a throwaway container, not the app.
//
// See internal/runner/README.md for the deployment guide and threat model.
package runner

import (
	"context"
	"errors"
)

// ErrDisabled is returned by every Run when execution is turned off.
var ErrDisabled = errors.New("server-side code execution is disabled")

// Request is one execution: a program plus the stdin it should receive.
type Request struct {
	Code      string `json:"code"`
	Stdin     string `json:"stdin,omitempty"`
	TimeoutMs int    `json:"timeout_ms,omitempty"` // wall-clock; clamped to limits
}

// Result is the outcome of a run. Error carries harness failures (the sandbox
// itself failing); a program that raises a Python exception is a successful
// run with a non-zero ExitCode and a populated Stderr.
type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	TimedOut   bool   `json:"timed_out"`
	Truncated  bool   `json:"truncated"`
	Error      string `json:"error,omitempty"`
}

// Executor runs a Request and reports the Result.
type Executor interface {
	Run(ctx context.Context, req Request) (Result, error)
	// Kind is "off", "local", or "remote" — for logging and health.
	Kind() string
	// Enabled reports whether Run can actually execute code.
	Enabled() bool
}

// Limits bound every run regardless of what the caller asks for.
type Limits struct {
	DefaultTimeoutMs int   // used when a Request doesn't specify one
	MaxTimeoutMs     int   // hard ceiling on wall-clock time
	MaxOutputBytes   int   // per stream (stdout, stderr); excess is truncated
	MaxMemoryBytes   int64 // address-space cap for the child (local mode)
	MaxProcesses     int   // nproc cap for the child (local mode)
}

// DefaultLimits are deliberately small — these run interactively.
func DefaultLimits() Limits {
	return Limits{
		DefaultTimeoutMs: 5000,
		MaxTimeoutMs:     10000,
		MaxOutputBytes:   64 * 1024,
		MaxMemoryBytes:   256 * 1024 * 1024,
		MaxProcesses:     64,
	}
}

// Config selects and parameterizes the executor.
type Config struct {
	Mode         string // "off" | "local" | "remote"
	PythonPath   string // local mode; defaults to "python3" on PATH
	SandboxURL   string // remote mode; the execd service base URL
	SandboxToken string // remote mode; shared bearer secret
	Deployed     bool   // true in prod (e.g. RAILWAY_ENVIRONMENT set)
	AllowUnsafe  bool   // CODE_EXEC_UNSAFE=1 — permit unsandboxed local in prod
	Limits       Limits
}

// New builds the executor described by cfg, applying the secure-by-default
// rules. It returns an error only for a misconfiguration the operator must
// fix (e.g. remote mode with no URL, or unsandboxed local in production
// without an explicit override); an unknown or empty mode yields a disabled
// executor, never an error.
func New(cfg Config) (Executor, error) {
	if (cfg.Limits == Limits{}) {
		cfg.Limits = DefaultLimits()
	}
	switch cfg.Mode {
	case "", "off", "disabled":
		return disabled{}, nil
	case "remote":
		if cfg.SandboxURL == "" {
			return nil, errors.New("runner: mode=remote needs CODE_EXEC_SANDBOX_URL")
		}
		return newRemote(cfg), nil
	case "local":
		if cfg.PythonPath == "" {
			cfg.PythonPath = "python3"
		}
		sandboxed := bwrapAvailable()
		if cfg.Deployed && !sandboxed && !cfg.AllowUnsafe {
			return nil, errors.New("runner: mode=local in a deployed environment without a bwrap sandbox is unsafe; " +
				"run the standalone execd service and use mode=remote, or set CODE_EXEC_UNSAFE=1 to accept the risk")
		}
		return newLocal(cfg, sandboxed), nil
	default:
		return nil, errors.New("runner: unknown mode " + cfg.Mode)
	}
}

// disabled is the default executor: it runs nothing.
type disabled struct{}

func (disabled) Run(context.Context, Request) (Result, error) { return Result{}, ErrDisabled }
func (disabled) Kind() string                                 { return "off" }
func (disabled) Enabled() bool                                { return false }

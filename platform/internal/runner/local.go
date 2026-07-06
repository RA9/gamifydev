package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// localExecutor runs Python in a subprocess on this host.
type localExecutor struct {
	python    string
	limits    Limits
	sandboxed bool // true when wrapping with bwrap (real namespace isolation)
}

func newLocal(cfg Config, sandboxed bool) *localExecutor {
	return &localExecutor{python: cfg.PythonPath, limits: cfg.Limits, sandboxed: sandboxed}
}

func (e *localExecutor) Kind() string  { return "local" }
func (e *localExecutor) Enabled() bool { return true }

// SandboxAvailable reports whether a real namespace sandbox (bubblewrap) backs
// local execution on this host. Callers (e.g. cmd/execd) use it to warn loudly
// when they'd be running untrusted code without isolation.
func SandboxAvailable() bool { return bwrapAvailable() }

// bwrapAvailable reports whether the bubblewrap binary is on PATH. bwrap gives
// unprivileged namespace isolation (no network, read-only root, private /tmp),
// which is what turns "local" mode from a dev convenience into a real sandbox.
func bwrapAvailable() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	_, err := exec.LookPath("bwrap")
	return err == nil
}

func (e *localExecutor) Run(ctx context.Context, req Request) (Result, error) {
	if strings.TrimSpace(req.Code) == "" {
		return Result{Error: "no code"}, nil
	}
	// Per-run scratch directory, wiped afterwards. All the program's file I/O
	// and any bytecode it writes stay inside here.
	dir, err := os.MkdirTemp("", "gd-run-*")
	if err != nil {
		return Result{}, fmt.Errorf("scratch dir: %w", err)
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "main.py"), []byte(req.Code), 0o600); err != nil {
		return Result{}, fmt.Errorf("write program: %w", err)
	}

	timeout := time.Duration(clampTimeout(req.TimeoutMs, e.limits)) * time.Millisecond
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	name, args := e.command(dir, timeout)
	cmd := exec.Command(name, args...)
	if !e.sandboxed {
		cmd.Dir = dir // bwrap sets its own chdir; plain mode runs in the scratch dir
	}
	cmd.Env = e.childEnv(dir)
	cmd.Stdin = strings.NewReader(req.Stdin)
	outCap := &capWriter{max: e.limits.MaxOutputBytes}
	errCap := &capWriter{max: e.limits.MaxOutputBytes}
	cmd.Stdout, cmd.Stderr = outCap, errCap
	setProcAttr(cmd) // new process group so we can kill the whole tree

	start := time.Now()
	if err := cmd.Start(); err != nil {
		return Result{Error: "could not start interpreter: " + err.Error()}, nil
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	var res Result
	select {
	case <-rctx.Done():
		killGroup(cmd) // SIGKILL the whole group; a while True: pass dies here
		<-done
		res.TimedOut = true
		res.ExitCode = -1
	case werr := <-done:
		res.ExitCode = exitCode(werr)
	}
	res.DurationMs = time.Since(start).Milliseconds()
	res.Stdout = outCap.String()
	res.Stderr = errCap.String()
	res.Truncated = outCap.truncated || errCap.truncated
	if res.TimedOut && res.Stderr == "" {
		res.Stderr = fmt.Sprintf("Killed: exceeded the %d ms time limit.", timeout.Milliseconds())
	}
	return res, nil
}

// command builds the argv. With bwrap we get a locked-down namespace; without
// it we fall back to a POSIX shell that applies ulimits before exec-ing Python.
func (e *localExecutor) command(dir string, timeout time.Duration) (string, []string) {
	cpuSecs := int(timeout/time.Second) + 1
	if e.sandboxed {
		// bubblewrap: read-only system dirs, private proc/dev, fresh tmpfs,
		// the scratch dir bound writable at /box, and no network at all.
		args := []string{
			"--ro-bind-try", "/usr", "/usr",
			"--ro-bind-try", "/bin", "/bin",
			"--ro-bind-try", "/lib", "/lib",
			"--ro-bind-try", "/lib64", "/lib64",
			"--ro-bind-try", "/etc/alternatives", "/etc/alternatives",
			"--proc", "/proc",
			"--dev", "/dev",
			"--tmpfs", "/tmp",
			"--bind", dir, "/box",
			"--chdir", "/box",
			"--unshare-all", // net, pid, ipc, uts, cgroup, user
			"--die-with-parent",
			"--new-session",
			"--clearenv",
			"--setenv", "HOME", "/box",
			"--setenv", "PATH", "/usr/bin:/bin",
			"--setenv", "PYTHONDONTWRITEBYTECODE", "1",
			"--setenv", "PYTHONUNBUFFERED", "1",
			"--setenv", "LC_ALL", "C.UTF-8",
			e.python, "-I", "-B", "main.py",
		}
		return "bwrap", args
	}
	// Plain fallback (dev only when not sandboxed). ulimit caps CPU time,
	// address space, file size, and process count before Python starts.
	var b strings.Builder
	b.WriteString("ulimit -t " + strconv.Itoa(cpuSecs))
	if runtime.GOOS == "linux" {
		// macOS's /bin/sh can't set these, so only do so on Linux.
		fmt.Fprintf(&b, "; ulimit -v %d", e.limits.MaxMemoryBytes/1024)
		fmt.Fprintf(&b, "; ulimit -u %d", e.limits.MaxProcesses)
		b.WriteString("; ulimit -f 8192")
	}
	b.WriteString("; exec " + shellQuote(e.python) + " -I -B main.py")
	return "/bin/sh", []string{"-c", b.String()}
}

// childEnv is a deliberately tiny environment — no inherited secrets.
func (e *localExecutor) childEnv(dir string) []string {
	return []string{
		"PATH=/usr/bin:/bin",
		"HOME=" + dir,
		"TMPDIR=" + dir,
		"PYTHONDONTWRITEBYTECODE=1",
		"PYTHONUNBUFFERED=1",
		"LC_ALL=C.UTF-8",
	}
}

func clampTimeout(ms int, l Limits) int {
	if ms <= 0 {
		ms = l.DefaultTimeoutMs
	}
	if ms > l.MaxTimeoutMs {
		ms = l.MaxTimeoutMs
	}
	return ms
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// exitCode extracts a process exit code from Wait's error.
func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return -1
}

// capWriter collects up to max bytes and drops the rest, noting truncation.
type capWriter struct {
	max       int
	buf       []byte
	truncated bool
}

func (w *capWriter) Write(p []byte) (int, error) {
	if room := w.max - len(w.buf); room > 0 {
		if len(p) > room {
			w.buf = append(w.buf, p[:room]...)
			w.truncated = true
		} else {
			w.buf = append(w.buf, p...)
		}
	} else if len(p) > 0 {
		w.truncated = true
	}
	return len(p), nil // always report full consumption so the pipe keeps draining
}

func (w *capWriter) String() string { return string(w.buf) }

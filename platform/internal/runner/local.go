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

// localExecutor runs a learner program in a subprocess on this host.
//
// Python is interpreted directly. C is compiled first and the resulting binary
// executed — both steps inside the same sandbox invocation, because the
// compiler is running on untrusted input too.
type localExecutor struct {
	python    string
	cc        string
	shell     string
	limits    Limits
	sandboxed bool // true when wrapping with bwrap (real namespace isolation)
}

func newLocal(cfg Config, sandboxed bool) *localExecutor {
	cc := cfg.CCPath
	if cc == "" {
		cc = "cc"
	}
	sh := cfg.ShellPath
	if sh == "" {
		// The Linux course teaches bash; fall back to POSIX sh where it is
		// absent rather than failing every shell submission.
		if _, err := exec.LookPath("bash"); err == nil {
			sh = "bash"
		} else {
			sh = "sh"
		}
	}
	return &localExecutor{python: cfg.PythonPath, cc: cc, shell: sh, limits: cfg.Limits, sandboxed: sandboxed}
}

// writeFixtures lays the authored fixture files into the scratch directory.
//
// Names are restricted to a plain filename. Authored content is trusted, but a
// path like "../../etc/passwd" in a fixture would escape the scratch directory
// before the sandbox ever starts — the one place a fixture could do damage.
func writeFixtures(dir string, files map[string]string) error {
	for name, content := range files {
		if name == "" || name != filepath.Base(name) || strings.HasPrefix(name, ".") {
			return fmt.Errorf("invalid fixture filename %q", name)
		}
		switch name {
		case "main.py", "main.c", "main.sh", "prog", "build.log":
			return fmt.Errorf("fixture %q would overwrite the program", name)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			return fmt.Errorf("write fixture %s: %w", name, err)
		}
	}
	return nil
}

// compileFailMarker is written to stderr by the C build script when the
// compiler rejects the program, so the caller can tell a build failure from a
// crash without guessing at exit codes.
const compileFailMarker = "__GD_COMPILE_FAILED__"

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

	lang, ok := NormalizeLang(req.Lang)
	if !ok {
		return Result{Error: "unsupported language: " + req.Lang}, nil
	}
	srcName := "main.py"
	switch lang {
	case LangC:
		srcName = "main.c"
	case LangShell:
		srcName = "main.sh"
	}
	if err := os.WriteFile(filepath.Join(dir, srcName), []byte(req.Code), 0o600); err != nil {
		return Result{}, fmt.Errorf("write program: %w", err)
	}
	if err := writeFixtures(dir, req.Files); err != nil {
		return Result{Error: err.Error()}, nil
	}

	timeout := time.Duration(clampTimeout(req.TimeoutMs, e.limits)) * time.Millisecond
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	name, args := e.command(dir, timeout, lang)
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
	if strings.HasPrefix(res.Stderr, compileFailMarker) {
		res.CompileFailed = true
		res.Stderr = strings.TrimSpace(strings.TrimPrefix(res.Stderr, compileFailMarker))
	}
	res.Truncated = outCap.truncated || errCap.truncated
	if res.TimedOut && res.Stderr == "" {
		res.Stderr = fmt.Sprintf("Killed: exceeded the %d ms time limit.", timeout.Milliseconds())
	}
	return res, nil
}

// command builds the argv. With bwrap we get a locked-down namespace; without
// it we fall back to a POSIX shell that applies ulimits before exec-ing.
//
// Both languages run through /bin/sh so the same limits apply either way, and
// so C can compile-then-exec in a single trip through the sandbox rather than
// two.
func (e *localExecutor) command(dir string, timeout time.Duration, lang string) (string, []string) {
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
			"/bin/sh", "-c", e.script(lang, cpuSecs, false),
		}
		return "bwrap", args
	}
	// Plain fallback (dev only when not sandboxed).
	return "/bin/sh", []string{"-c", e.script(lang, cpuSecs, true)}
}

// script is the shell program that runs the learner's code.
//
// withUlimits applies resource caps on the plain path; under bwrap the
// namespace already bounds the process, and the compiler legitimately needs
// more headroom than a learner's program does.
func (e *localExecutor) script(lang string, cpuSecs int, withUlimits bool) string {
	var b strings.Builder
	if withUlimits {
		b.WriteString("ulimit -t " + strconv.Itoa(cpuSecs))
		if runtime.GOOS == "linux" {
			// macOS's /bin/sh can't set these, so only do so on Linux.
			fmt.Fprintf(&b, "; ulimit -v %d", e.limits.MaxMemoryBytes/1024)
			fmt.Fprintf(&b, "; ulimit -u %d", e.limits.MaxProcesses)
			b.WriteString("; ulimit -f 8192")
		}
		b.WriteString("; ")
	}
	if lang == LangC {
		// Compile first. On failure emit the marker followed by the compiler's
		// own diagnostics — for a beginner those are the most useful thing on
		// the screen, so they are passed through rather than summarised.
		//
		// -std=c11 pins the dialect; -O0 keeps line numbers honest in a crash;
		// warnings are deliberately left on and shown, because they are
		// teaching material in a C course.
		b.WriteString(shellQuote(e.cc) + " -std=c11 -O0 -Wall -o prog main.c 2>build.log")
		b.WriteString(" || { printf '%s' " + shellQuote(compileFailMarker) + " >&2; cat build.log >&2; exit 97; }")
		b.WriteString("; cat build.log >&2") // warnings from a successful build
		b.WriteString("; exec ./prog")
		return b.String()
	}
	if lang == LangShell {
		// The script is the program. It runs in the scratch dir, so any files it
		// creates stay inside the sandbox and vanish with it.
		b.WriteString("exec " + shellQuote(e.shell) + " main.sh")
		return b.String()
	}
	b.WriteString("exec " + shellQuote(e.python) + " -I -B main.py")
	return b.String()
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

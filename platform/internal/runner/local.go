package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	golang    string
	gocache   string
	bincache  string
	shell     string
	limits    Limits
	sandboxed bool // true when wrapping with bwrap (real namespace isolation)
}

func newLocal(cfg Config, sandboxed bool) *localExecutor {
	cc := cfg.CCPath
	if cc == "" {
		cc = "cc"
	}
	golang := cfg.GoPath
	if golang == "" {
		golang = "go"
	}
	// Resolve the toolchain to an absolute path. The child runs with a
	// deliberately tiny PATH of /usr/bin:/bin, and Go is rarely installed
	// there — /usr/local/go/bin on most Linux boxes, a Homebrew prefix on a
	// Mac. Go locates its own GOROOT relative to the binary, so the absolute
	// path is all it needs. (cc and python3 are left as names because they do
	// live on that PATH.)
	if abs, err := exec.LookPath(golang); err == nil {
		golang = abs
	}
	// One build cache for every run, rather than a fresh one per submission.
	//
	// A cold cache makes Go compile the standard library — measured at ~20s
	// here against ~1s warm — so a per-run cache would time out every Go
	// submission on a limit any reasonable problem would set.
	//
	// Sharing it across untrusted programs is safe enough to be worth it: Go
	// keys the cache by a hash of the actual build inputs, so one submission
	// cannot make another's compile resolve to something it planted. What it
	// can do is fill the disk, which is why this lives under the system temp
	// dir where the usual cleanup applies.
	gocache := cfg.GoCache
	if gocache == "" {
		gocache = filepath.Join(os.TempDir(), "gamifydev-gocache")
	}
	_ = os.MkdirAll(gocache, 0o700)
	// Compiled programs are cached by a hash of their own source.
	//
	// Grading one submission runs it once per test input, and each of those was
	// a fresh compile — six cases meant six identical builds of the same file.
	// The binary is a pure function of the source, so the second case onward can
	// reuse the first one's. Keyed by content, like the Go cache above, for the
	// same reason: a submission cannot make another's lookup find its binary.
	bincache := filepath.Join(os.TempDir(), "gamifydev-bincache")
	_ = os.MkdirAll(bincache, 0o700)
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
	return &localExecutor{
		python: cfg.PythonPath, cc: cc, golang: golang, gocache: gocache,
		bincache: bincache, shell: sh, limits: cfg.Limits, sandboxed: sandboxed,
	}
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
		case "main.py", "main.c", "main.sh", "main.go", "go.mod", "prog", "build.log":
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
	case LangGo:
		srcName = "main.go"
	}
	if err := os.WriteFile(filepath.Join(dir, srcName), []byte(req.Code), 0o600); err != nil {
		return Result{}, fmt.Errorf("write program: %w", err)
	}
	if err := writeFixtures(dir, req.Files); err != nil {
		return Result{Error: err.Error()}, nil
	}

	// The learner's budget is about their algorithm. A compiled language spends
	// time before the program starts that has nothing to do with them — and a
	// problem whose limit is tuned so a correct Python answer fits would fail
	// the same answer in Go purely on build time. So the compiler gets its own
	// allowance on top.
	timeout := time.Duration(clampTimeout(req.TimeoutMs, e.limits)) * time.Millisecond
	// Only when this program actually has to be built. A cached binary runs
	// with the learner's budget and nothing more, which keeps a runaway program
	// on a short leash instead of inheriting a compiler's allowance it never
	// used.
	key := buildKey(lang, req.Code)
	if _, err := os.Stat(filepath.Join(e.bincache, key)); err != nil {
		timeout += buildAllowance(lang)
	}
	rctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	name, args := e.command(dir, timeout, lang, key)
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
func (e *localExecutor) command(dir string, timeout time.Duration, lang, key string) (string, []string) {
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
			"--bind", e.gocache, e.gocache,
			"--bind", e.bincache, e.bincache,
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
			"/bin/sh", "-c", e.script(lang, key, cpuSecs, true),
		}
		return "bwrap", args
	}
	// Plain fallback (dev only when not sandboxed).
	return "/bin/sh", []string{"-c", e.script(lang, key, cpuSecs, true)}
}

// script is the shell program that runs the learner's code.
//
// withUlimits applies resource caps inside either execution path. Bubblewrap
// isolates namespaces and the filesystem, but it does not cap memory, process
// count, CPU time, or file growth on its own.
func (e *localExecutor) script(lang, key string, cpuSecs int, withUlimits bool) string {
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
		cached := shellQuote(filepath.Join(e.bincache, key))
		b.WriteString("if [ -f " + cached + " ]; then cp " + cached + " prog && chmod +x prog; else ")
		b.WriteString(shellQuote(e.cc) + " -std=c11 -O0 -Wall -o prog main.c 2>build.log")
		b.WriteString(" || { printf '%s' " + shellQuote(compileFailMarker) + " >&2; cat build.log >&2; exit 97; }; ")
		b.WriteString("cat build.log >&2; cp prog " + cached + ".tmp$$ 2>/dev/null && mv " + cached + ".tmp$$ " + cached + " 2>/dev/null; fi")
		b.WriteString("; exec ./prog")
		return b.String()
	}
	if lang == LangGo {
		// A module file, written here rather than asked of the learner: a bare
		// `go build` outside a module fails with advice about `go mod init`,
		// which is not the lesson. GOFLAGS pins it to the local module and
		// GOPROXY=off makes a stray import fail as a missing dependency instead
		// of hanging on a network the sandbox does not have.
		cached := shellQuote(filepath.Join(e.bincache, key))
		b.WriteString("if [ -f " + cached + " ]; then cp " + cached + " prog && chmod +x prog; else ")
		b.WriteString("printf 'module solution\ngo 1.21\n' > go.mod; ")
		b.WriteString("GOCACHE=" + shellQuote(e.gocache) + " GOPATH=$HOME/go GOPROXY=off GOFLAGS=-mod=mod ")
		b.WriteString(shellQuote(e.golang) + " build -o prog main.go 2>build.log")
		b.WriteString(" || { printf '%s' " + shellQuote(compileFailMarker) + " >&2; cat build.log >&2; exit 97; }; ")
		b.WriteString("cat build.log >&2; cp prog " + cached + ".tmp$$ 2>/dev/null && mv " + cached + ".tmp$$ " + cached + " 2>/dev/null; fi")
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

// buildKey names the cached binary for a program: a hash of the language and
// the source, so only an identical program finds it.
func buildKey(lang, code string) string {
	sum := sha256.Sum256([]byte(lang + "\x00" + code))
	return hex.EncodeToString(sum[:])
}

// buildAllowance is the extra wall clock a compiled language may spend before
// its program starts running.
//
// Go's is the larger because it links the standard library even from a warm
// cache, and because several submissions compiling at once contend for the same
// machine. It is bounded rather than generous: a program that never terminates
// still dies at the limit plus this, not at leisure.
func buildAllowance(lang string) time.Duration {
	switch lang {
	case LangGo:
		return 20 * time.Second
	case LangC:
		return 15 * time.Second
	}
	return 0
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

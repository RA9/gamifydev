package runner

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// localExec builds a local executor for tests, skipping if python3 is absent.
func localExec(t *testing.T) *localExecutor {
	t.Helper()
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not on PATH")
	}
	return newLocal(Config{PythonPath: py, Limits: DefaultLimits()}, false)
}

func TestLocalHelloWorld(t *testing.T) {
	e := localExec(t)
	res, err := e.Run(context.Background(), Request{Code: "print('hello', 1 + 2)"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.ExitCode != 0 || res.TimedOut {
		t.Fatalf("unexpected: %+v", res)
	}
	if strings.TrimSpace(res.Stdout) != "hello 3" {
		t.Fatalf("stdout = %q, want %q", res.Stdout, "hello 3")
	}
}

func TestLocalStdin(t *testing.T) {
	e := localExec(t)
	code := "import sys\nname = sys.stdin.readline().strip()\nprint('hi ' + name)"
	res, err := e.Run(context.Background(), Request{Code: code, Stdin: "Ada\n"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(res.Stdout) != "hi Ada" {
		t.Fatalf("stdout = %q", res.Stdout)
	}
}

func TestLocalExceptionExitCode(t *testing.T) {
	e := localExec(t)
	res, err := e.Run(context.Background(), Request{Code: "raise ValueError('boom')"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got %+v", res)
	}
	if !strings.Contains(res.Stderr, "ValueError") {
		t.Fatalf("stderr missing traceback: %q", res.Stderr)
	}
}

func TestLocalInfiniteLoopIsKilled(t *testing.T) {
	e := localExec(t)
	e.limits.MaxTimeoutMs = 1500
	start := time.Now()
	res, err := e.Run(context.Background(), Request{Code: "while True:\n    pass", TimeoutMs: 1000})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !res.TimedOut {
		t.Fatalf("expected timeout, got %+v", res)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("kill took too long: %v", elapsed)
	}
}

func TestLocalOutputTruncated(t *testing.T) {
	e := localExec(t)
	e.limits.MaxOutputBytes = 100
	res, err := e.Run(context.Background(), Request{Code: "print('x' * 10000)"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !res.Truncated {
		t.Fatalf("expected truncation")
	}
	if len(res.Stdout) > 100 {
		t.Fatalf("stdout not capped: %d bytes", len(res.Stdout))
	}
}

func TestNewSecureByDefault(t *testing.T) {
	e, err := New(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if e.Enabled() || e.Kind() != "off" {
		t.Fatalf("default should be disabled, got kind=%s enabled=%v", e.Kind(), e.Enabled())
	}
	if _, err := e.Run(context.Background(), Request{Code: "print(1)"}); err != ErrDisabled {
		t.Fatalf("disabled Run should return ErrDisabled, got %v", err)
	}
}

func TestNewRejectsUnsafeLocalInProd(t *testing.T) {
	// Deployed + no bwrap + no override → must refuse (this host has no bwrap).
	if bwrapAvailable() {
		t.Skip("host has bwrap; the unsafe-in-prod guard doesn't trigger")
	}
	if _, err := New(Config{Mode: "local", Deployed: true}); err == nil {
		t.Fatal("expected New to refuse unsandboxed local in a deployed env")
	}
	// With the explicit override it should be allowed.
	if _, err := New(Config{Mode: "local", Deployed: true, AllowUnsafe: true}); err != nil {
		t.Fatalf("override should allow it: %v", err)
	}
}

func TestNewRemoteNeedsURL(t *testing.T) {
	if _, err := New(Config{Mode: "remote"}); err == nil {
		t.Fatal("remote mode without a URL should error")
	}
	if _, err := New(Config{Mode: "remote", SandboxURL: "http://localhost:9999"}); err != nil {
		t.Fatalf("remote with URL should be fine: %v", err)
	}
}

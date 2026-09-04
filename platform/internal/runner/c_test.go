package runner

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

// newCTestExecutor builds a local, unsandboxed executor for tests, skipping if
// no C compiler is on PATH.
func newCTestExecutor(t *testing.T) Executor {
	t.Helper()
	if _, err := exec.LookPath("cc"); err != nil {
		t.Skip("no C compiler on PATH")
	}
	e, err := New(Config{Mode: "local", PythonPath: "python3", CCPath: "cc", AllowUnsafe: true})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	return e
}

func TestCRunsAndCapturesStdout(t *testing.T) {
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{Lang: LangC, Code: `
#include <stdio.h>
int main(void) { printf("hello from C\n"); return 0; }
`})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.CompileFailed {
		t.Fatalf("valid program reported a compile failure: %s", res.Stderr)
	}
	if !strings.Contains(res.Stdout, "hello from C") {
		t.Fatalf("stdout = %q", res.Stdout)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit = %d, want 0", res.ExitCode)
	}
}

func TestCReadsStdin(t *testing.T) {
	// The shape most beginner C checkpoints take: read input, print an answer.
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{Lang: LangC, Stdin: "3 9 4\n", Code: `
#include <stdio.h>
int main(void) {
    int a, b, c;
    if (scanf("%d %d %d", &a, &b, &c) != 3) return 1;
    int m = a;
    if (b > m) m = b;
    if (c > m) m = c;
    printf("%d\n", m);
    return 0;
}
`})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(res.Stdout) != "9" {
		t.Fatalf("stdout = %q, want 9 (stderr: %s)", res.Stdout, res.Stderr)
	}
}

func TestCCompileFailureIsDistinctFromACrash(t *testing.T) {
	// A learner needs to know "this doesn't build" is a different problem from
	// "this ran and died" — they are different lessons.
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{Lang: LangC, Code: `
#include <stdio.h>
int main(void) { printf("missing semicolon") return 0; }
`})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !res.CompileFailed {
		t.Fatal("a syntax error was not reported as a compile failure")
	}
	// The compiler's own diagnostics are the most useful thing on screen.
	if !strings.Contains(res.Stderr, "main.c") {
		t.Fatalf("compiler diagnostics missing from stderr: %q", res.Stderr)
	}
	// The marker itself must never leak into what the learner reads.
	if strings.Contains(res.Stderr, compileFailMarker) {
		t.Fatalf("internal marker leaked into stderr: %q", res.Stderr)
	}
}

func TestCRuntimeCrashIsNotACompileFailure(t *testing.T) {
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{Lang: LangC, Code: `
#include <stdio.h>
int main(void) { int *p = 0; *p = 1; printf("unreachable\n"); return 0; }
`})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.CompileFailed {
		t.Fatal("a segfault was misreported as a compile failure")
	}
	if res.ExitCode == 0 {
		t.Fatalf("a null dereference exited 0")
	}
}

func TestCInfiniteLoopIsKilled(t *testing.T) {
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{
		Lang: LangC, TimeoutMs: 1500,
		Code: "int main(void) { for (;;) { } }",
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !res.TimedOut {
		t.Fatal("an infinite loop was not killed by the timeout")
	}
}

func TestPythonStillWorksAndIsTheDefault(t *testing.T) {
	// Adding C must not change what an empty Lang means for every existing
	// caller (step labs, the probe, assignment checks).
	e := newCTestExecutor(t)
	for _, lang := range []string{"", LangPython} {
		res, err := e.Run(context.Background(), Request{Lang: lang, Code: `print("py ok")`})
		if err != nil {
			t.Fatalf("lang %q: %v", lang, err)
		}
		if !strings.Contains(res.Stdout, "py ok") {
			t.Fatalf("lang %q: stdout = %q (stderr %q)", lang, res.Stdout, res.Stderr)
		}
	}
}

func TestUnknownLanguageIsRejected(t *testing.T) {
	// Silently treating an unknown language as Python would grade a Java
	// submission by running it through the wrong interpreter.
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{Lang: "java", Code: "class X {}"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Error == "" {
		t.Fatal("an unsupported language was accepted")
	}
	if _, ok := NormalizeLang("java"); ok {
		t.Fatal("NormalizeLang accepted java")
	}
}

func TestShellRuns(t *testing.T) {
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{
		Lang:  LangShell,
		Code:  "echo \"lines: $(wc -l < data.txt | tr -d ' ')\"",
		Files: map[string]string{"data.txt": "alpha\nbeta\ngamma\n"},
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(res.Stdout) != "lines: 3" {
		t.Fatalf("stdout = %q (stderr %q)", res.Stdout, res.Stderr)
	}
}

func TestShellReadsStdinAndReportsExit(t *testing.T) {
	e := newCTestExecutor(t)
	res, err := e.Run(context.Background(), Request{
		Lang:  LangShell,
		Stdin: "one\ntwo\nthree\n",
		Code:  "wc -l | tr -d ' '\nexit 3",
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(res.Stdout) != "3" {
		t.Fatalf("stdout = %q", res.Stdout)
	}
	// A script's exit status is often the thing being checked.
	if res.ExitCode != 3 {
		t.Fatalf("exit = %d, want 3", res.ExitCode)
	}
}

func TestFixturesCannotEscapeTheScratchDir(t *testing.T) {
	// Authored content is trusted, but a traversing filename would be written
	// before the sandbox even starts — so it is rejected outright.
	e := newCTestExecutor(t)
	for _, bad := range []string{"../escape.txt", "sub/dir.txt", "/etc/passwd", "main.py"} {
		res, err := e.Run(context.Background(), Request{
			Lang: LangShell, Code: "echo hi", Files: map[string]string{bad: "x"},
		})
		if err != nil {
			t.Fatalf("run %q: %v", bad, err)
		}
		if res.Error == "" {
			t.Fatalf("fixture name %q was accepted", bad)
		}
	}
}

func TestFixturesDoNotLeakBetweenRuns(t *testing.T) {
	// Each run gets a fresh scratch dir; a file written by one submission must
	// not be visible to the next.
	e := newCTestExecutor(t)
	if _, err := e.Run(context.Background(), Request{
		Lang: LangShell, Code: "echo secret > leaked.txt",
	}); err != nil {
		t.Fatalf("first run: %v", err)
	}
	res, err := e.Run(context.Background(), Request{
		Lang: LangShell, Code: "cat leaked.txt 2>/dev/null || echo clean",
	})
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !strings.Contains(res.Stdout, "clean") {
		t.Fatalf("a file leaked between runs: %q", res.Stdout)
	}
}

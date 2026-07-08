package runner

import (
	"context"
	"testing"
)

// stubExecutor returns a fixed stdout, so we can test Probe's verdict logic for
// environments we can't reproduce on the test host (e.g. a real bwrap sandbox).
type stubExecutor struct {
	stdout   string
	timedOut bool
}

func (s stubExecutor) Run(context.Context, Request) (Result, error) {
	return Result{Stdout: s.stdout, TimedOut: s.timedOut}, nil
}
func (stubExecutor) Kind() string  { return "local" }
func (stubExecutor) Enabled() bool { return true }

func TestProbeVerdictSandboxed(t *testing.T) {
	// Signals a healthy bwrap sandbox: cwd is /box, network unreachable,
	// /etc/hosts not present.
	e := stubExecutor{stdout: `__GD_PROBE__{"cwd":"/box","etc_hosts":false,"net":false,"pid":2}` + "\n"}
	rep, err := Probe(context.Background(), e)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Sandboxed {
		t.Fatalf("expected sandboxed, got %+v", rep)
	}
	if !rep.UnderBwrapMount || !rep.NetworkBlocked || !rep.HostFSHidden {
		t.Fatalf("signals wrong: %+v", rep)
	}
}

func TestProbeVerdictNotSandboxed(t *testing.T) {
	cases := map[string]string{
		"network reachable":     `__GD_PROBE__{"cwd":"/box","etc_hosts":false,"net":true,"pid":2}`,
		"host fs visible":       `__GD_PROBE__{"cwd":"/box","etc_hosts":true,"net":false,"pid":2}`,
		"not under bwrap mount": `__GD_PROBE__{"cwd":"/tmp/gd-run-1","etc_hosts":false,"net":false,"pid":900}`,
	}
	for name, out := range cases {
		rep, _ := Probe(context.Background(), stubExecutor{stdout: out + "\n"})
		if rep.Sandboxed {
			t.Fatalf("%s: expected NOT sandboxed, got %+v", name, rep)
		}
		if len(rep.Reasons) == 0 {
			t.Fatalf("%s: expected a reason", name)
		}
	}
}

func TestProbeFailsClosed(t *testing.T) {
	// No verdict line, or a timeout, must never read as sandboxed.
	for _, e := range []stubExecutor{
		{stdout: "some noise but no marker\n"},
		{stdout: "", timedOut: true},
		{stdout: "__GD_PROBE__not-json\n"},
	} {
		rep, _ := Probe(context.Background(), e)
		if rep.Sandboxed {
			t.Fatalf("fail-closed violated for %+v -> %+v", e, rep)
		}
	}
}

func TestProbeRealExecutorNotSandboxedHere(t *testing.T) {
	// On this dev host there's no bwrap, so a real run must report NOT sandboxed
	// (cwd is a temp dir, not /box). This guards against the probe ever
	// false-positively declaring an unsandboxed host safe.
	e := localExec(t) // skips if python3 missing
	rep, err := Probe(context.Background(), e)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Sandboxed {
		t.Fatalf("unsandboxed host reported as sandboxed: %+v", rep)
	}
	if rep.UnderBwrapMount {
		t.Fatalf("cwd unexpectedly under /box on a no-bwrap host: %s", rep.Raw)
	}
}

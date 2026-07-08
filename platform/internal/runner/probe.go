package runner

import (
	"context"
	"encoding/json"
	"strings"
)

// probeProgram runs inside the executor and reports what the sandbox actually
// allows. It is the ground truth for "is untrusted code really contained here":
// it checks the working directory (our bwrap config chdir's to /box), whether
// the host's /etc/hosts is visible, and whether outbound network works.
const probeProgram = `import os, json, socket
r = {}
r["cwd"] = os.getcwd()
r["etc_hosts"] = os.path.exists("/etc/hosts")
try:
    s = socket.create_connection(("1.1.1.1", 53), timeout=2)
    s.close()
    r["net"] = True
except Exception:
    r["net"] = False
r["pid"] = os.getpid()
print("__GD_PROBE__" + json.dumps(r))
`

const probeMarker = "__GD_PROBE__"

// SandboxReport is the verdict of a self-test run. Sandboxed is true only when
// every containment signal held; anything unexpected fails closed.
type SandboxReport struct {
	Sandboxed       bool     `json:"sandboxed"`
	UnderBwrapMount bool     `json:"under_bwrap_mount"` // cwd is /box → our bwrap ran
	NetworkBlocked  bool     `json:"network_blocked"`   // couldn't reach the internet
	HostFSHidden    bool     `json:"host_fs_hidden"`    // /etc/hosts not visible
	Reasons         []string `json:"reasons,omitempty"` // why it's NOT sandboxed
	Raw             string   `json:"raw,omitempty"`     // probe's JSON, for logs
}

// Probe runs the containment self-test through the given executor and reports
// whether untrusted code is actually sandboxed. It fails closed: a timeout,
// harness error, or unparseable output all yield Sandboxed=false.
func Probe(ctx context.Context, e Executor) (SandboxReport, error) {
	res, err := e.Run(ctx, Request{Code: probeProgram, TimeoutMs: 8000})
	if err != nil {
		return SandboxReport{Reasons: []string{"probe could not run: " + err.Error()}}, err
	}
	if res.TimedOut {
		return SandboxReport{Reasons: []string{"probe timed out"}}, nil
	}
	var line string
	for _, l := range strings.Split(res.Stdout, "\n") {
		if strings.HasPrefix(l, probeMarker) {
			line = l[len(probeMarker):]
			break
		}
	}
	if line == "" {
		return SandboxReport{Reasons: []string{"probe produced no verdict; stderr: " + strings.TrimSpace(res.Stderr)}}, nil
	}
	var raw struct {
		Cwd      string `json:"cwd"`
		EtcHosts bool   `json:"etc_hosts"`
		Net      bool   `json:"net"`
		PID      int    `json:"pid"`
	}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return SandboxReport{Reasons: []string{"probe output unparseable"}, Raw: line}, nil
	}

	rep := SandboxReport{
		UnderBwrapMount: strings.HasPrefix(raw.Cwd, "/box"),
		NetworkBlocked:  !raw.Net,
		HostFSHidden:    !raw.EtcHosts,
		Raw:             line,
	}
	if !rep.UnderBwrapMount {
		rep.Reasons = append(rep.Reasons, "not running under the bwrap mount (cwd="+raw.Cwd+")")
	}
	if !rep.NetworkBlocked {
		rep.Reasons = append(rep.Reasons, "outbound network is reachable")
	}
	if !rep.HostFSHidden {
		rep.Reasons = append(rep.Reasons, "the host filesystem (/etc/hosts) is visible")
	}
	rep.Sandboxed = rep.UnderBwrapMount && rep.NetworkBlocked && rep.HostFSHidden
	return rep, nil
}

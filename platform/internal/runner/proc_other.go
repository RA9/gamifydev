//go:build !unix

package runner

import "os/exec"

// Non-unix fallback (e.g. Windows dev machines): no process-group semantics, so
// kill just the process. Production runs on Linux, which uses proc_unix.go.
func setProcAttr(*exec.Cmd) {}

func killGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}

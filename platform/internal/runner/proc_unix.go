//go:build unix

package runner

import (
	"os/exec"
	"syscall"
)

// setProcAttr puts the child in its own process group so that killGroup can
// take down the interpreter and anything it spawned in one signal.
func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killGroup SIGKILLs the child's whole process group (negative PID).
func killGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

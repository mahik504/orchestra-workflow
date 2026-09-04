//go:build !windows

package runner

import (
	"os/exec"
	"syscall"
)

// prepareCommandPlatform sets Setpgid: true on Unix so the child process
// becomes leader of a new process group.
func prepareCommandPlatform(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

// killProcessTree terminates the entire process group on Unix using SIGKILL to -PGID.
func killProcessTree(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil || cmd.Process.Pid <= 0 {
		return nil
	}
	pid := cmd.Process.Pid

	// On Unix, when Setpgid is true, PGID == PID of the root child process.
	// Sending SIGKILL to -pid kills all processes in that process group.
	err := syscall.Kill(-pid, syscall.SIGKILL)
	if err != nil {
		// Fallback to killing the root process if process group kill fails
		_ = cmd.Process.Kill()
		return err
	}
	return nil
}

//go:build windows

package runner

import (
	"os/exec"
	"strconv"
)

// prepareCommandPlatform configures platform-specific command attributes on Windows.
func prepareCommandPlatform(cmd *exec.Cmd) {
	// On Windows, child processes are spawned normally.
}

// killProcessTree terminates the command process and all its child/descendant processes on Windows.
func killProcessTree(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil || cmd.Process.Pid <= 0 {
		return nil
	}
	pid := cmd.Process.Pid

	// 1. Forcefully terminate the process and all descendant child processes via taskkill /F /T /PID
	taskkillPath, err := exec.LookPath("taskkill")
	if err != nil {
		taskkillPath = "taskkill"
	}
	killCmd := exec.Command(taskkillPath, "/F", "/T", "/PID", strconv.Itoa(pid))
	_ = killCmd.Run()

	// 2. Ensure root process handle is terminated
	_ = cmd.Process.Kill()
	return nil
}

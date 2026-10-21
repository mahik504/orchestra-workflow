package benchmark

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// runHost starts the operator executable in a fresh process.
// It is not a host manager and it does not select a model.
func runHost(ctx context.Context, exe, workspace, prompt, gocache, modcache, modelLabel string, timeout time.Duration) (exit int, timedOut bool, wall time.Duration, stdoutBytes int, stderrBytes int, err error) {
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	hctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	stdoutPath := filepath.Join(workspace, "host-stdout.txt")
	stderrPath := filepath.Join(workspace, "host-stderr.txt")
	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return -1, false, 0, 0, 0, err
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return -1, false, 0, 0, 0, err
	}
	defer stderrFile.Close()
	cmd := exec.CommandContext(hctx, exe)
	cmd.Dir = workspace
	cmd.Env = hostEnv(workspace, prompt, gocache, modcache, modelLabel)
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.WaitDelay = 500 * time.Millisecond
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			killTree(cmd.Process.Pid)
		}
		return nil
	}
	start := time.Now()
	err = cmd.Run()
	wall = time.Since(start)
	_ = stdoutFile.Close()
	_ = stderrFile.Close()
	stdoutBytes = fileSize(stdoutPath)
	stderrBytes = fileSize(stderrPath)
	if errors.Is(hctx.Err(), context.DeadlineExceeded) {
		return -1, true, wall, stdoutBytes, stderrBytes, err
	}
	if err == nil {
		return 0, false, wall, stdoutBytes, stderrBytes, nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), false, wall, stdoutBytes, stderrBytes, err
	}
	return -1, false, wall, stdoutBytes, stderrBytes, err
}

func fileSize(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return int(info.Size())
}

func hostEnv(workspace, prompt, gocache, modcache, modelLabel string) []string {
	deny := map[string]bool{
		"ORCHESTRA_HOME":              true,
		"ORCHESTRA_MEMORY_PATH":       true,
		"ORCHESTRA_OVERLAY_PATH":      true,
		"ORCHESTRA_BENCH_MODEL_LABEL": true,
		"ORCHESTRA_BENCH_HOST":        true,
		"ORCHESTRA_BENCH_REAL":        true,
		"ORCHESTRA_BENCH_WORKSPACE":   true,
		"ORCHESTRA_BENCH_PROMPT":      true,
		"GOCACHE":                     true,
		"GOMODCACHE":                  true,
	}
	var env []string
	for _, kv := range os.Environ() {
		key, _, ok := strings.Cut(kv, "=")
		if !ok || deny[key] {
			continue
		}
		env = append(env, kv)
	}
	env = append(env,
		"ORCHESTRA_BENCH_WORKSPACE="+workspace,
		"ORCHESTRA_BENCH_PROMPT="+prompt,
		"GOCACHE="+gocache,
		"GOMODCACHE="+modcache,
		"GOPROXY=off",
	)
	if strings.TrimSpace(modelLabel) != "" {
		env = append(env, "ORCHESTRA_BENCH_MODEL_LABEL="+modelLabel)
	}
	return env
}

func killTree(pid int) {
	if pid <= 0 {
		return
	}
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).Run()
		return
	}
	_ = exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
}

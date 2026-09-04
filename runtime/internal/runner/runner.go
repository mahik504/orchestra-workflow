package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

// Common Runner Errors
var (
	ErrCommandTimeout        = errors.New("command execution timed out")
	ErrCommandNotFound       = errors.New("executable binary not found in PATH")
	ErrInvalidExecutableName = errors.New("invalid or unsafe executable name")
	ErrSecurityViolation     = errors.New("command violates security policy")
	ErrWorkingDirectory      = errors.New("invalid working directory")
	ErrOutputLimitExceeded   = errors.New("command output exceeded maximum allowed size")
)

// CommandRunner abstracts subprocess execution with security and timeout controls.
type CommandRunner interface {
	Run(ctx context.Context, cmd Command) (*RunResult, error)
	LookPath(file string) (string, error)
}

// Command encapsulates subprocess invocation parameters.
type Command struct {
	Name           string        `json:"name"`
	Args           []string      `json:"args"`
	Dir            string        `json:"dir"`
	Env            []string      `json:"env,omitempty"`
	Timeout        time.Duration `json:"timeout,omitempty"`
	Stdout         io.Writer     `json:"-"`
	Stderr         io.Writer     `json:"-"`
	MaxOutputBytes int64         `json:"max_output_bytes,omitempty"`
	StripANSI      bool          `json:"strip_ansi,omitempty"`
}

// RunResult captures the complete execution outcome.
type RunResult struct {
	Command         string        `json:"command"`
	Args            []string      `json:"args"`
	Dir             string        `json:"dir"`
	ExitCode        int           `json:"exit_code"`
	Stdout          string        `json:"stdout"`
	Stderr          string        `json:"stderr"`
	CombinedOutput  string        `json:"combined_output"`
	Duration        time.Duration `json:"duration"`
	PID             int           `json:"pid"`
	TimedOut        bool          `json:"timed_out"`
	OutputTruncated bool          `json:"output_truncated"`
}

// OSCommandRunner is the production implementation of CommandRunner enforcing direct OS execution.
type OSCommandRunner struct {
	DefaultTimeout  time.Duration
	MaxDefaultBytes int64
}

// NewOSCommandRunner creates a standard OS command runner.
func NewOSCommandRunner() *OSCommandRunner {
	return &OSCommandRunner{
		DefaultTimeout:  5 * time.Minute,
		MaxDefaultBytes: 10 * 1024 * 1024, // 10MB
	}
}

// LookPath locates an executable in PATH, handling extensions and validating safety.
func (r *OSCommandRunner) LookPath(file string) (string, error) {
	// 1. Check raw file for null bytes and control characters BEFORE trimming
	if strings.ContainsRune(file, '\x00') || strings.ContainsAny(file, "\r\n") {
		return "", fmt.Errorf("%w: executable path contains control characters", ErrSecurityViolation)
	}
	if strings.ContainsAny(file, "|;&$`><") {
		return "", fmt.Errorf("%w: executable path contains prohibited shell characters", ErrSecurityViolation)
	}

	trimmed := strings.TrimSpace(file)
	if trimmed == "" {
		return "", ErrInvalidExecutableName
	}

	path, err := exec.LookPath(trimmed)
	if err != nil {
		return "", fmt.Errorf("%w: %s (%v)", ErrCommandNotFound, trimmed, err)
	}
	return path, nil
}

// Run executes the command with timeout, output capture, and security checks.
// Direct execution is strictly enforced (never via a shell).
func (r *OSCommandRunner) Run(ctx context.Context, cmd Command) (*RunResult, error) {
	// 1. Sanitize Executable Name on raw input BEFORE trimming
	if strings.ContainsRune(cmd.Name, '\x00') || strings.ContainsAny(cmd.Name, "\r\n") {
		return nil, fmt.Errorf("%w: executable name contains control characters", ErrSecurityViolation)
	}
	// Prohibit shell operators in executable name to prevent injection tricks
	if strings.ContainsAny(cmd.Name, "|;&$`><") {
		return nil, fmt.Errorf("%w: executable contains prohibited shell characters", ErrSecurityViolation)
	}

	execName := strings.TrimSpace(cmd.Name)
	if execName == "" {
		return nil, ErrInvalidExecutableName
	}

	// 2. Sanitize Arguments
	for i, arg := range cmd.Args {
		if strings.ContainsRune(arg, '\x00') {
			return nil, fmt.Errorf("%w: argument %d contains null bytes", ErrSecurityViolation, i)
		}
	}

	// 2b. Sanitize Environment Variables (scan for null bytes and control characters)
	for i, envVar := range cmd.Env {
		if strings.ContainsRune(envVar, '\x00') {
			return nil, fmt.Errorf("%w: environment variable %d contains null bytes", ErrSecurityViolation, i)
		}
		if strings.ContainsAny(envVar, "\r\n") {
			return nil, fmt.Errorf("%w: environment variable %d contains control characters", ErrSecurityViolation, i)
		}
	}

	// 3. Validate Working Directory & Quarantine Boundary
	// Security boundary enforcement MUST precede filesystem I/O checks
	targetDir := cmd.Dir
	if targetDir == "" {
		if cwd, err := os.Getwd(); err == nil {
			targetDir = cwd
		}
	}
	if targetDir != "" {
		if err := resources.CheckQuarantineBoundary(targetDir); err != nil {
			// Wrap with %w so errors.Is(err, resources.ErrQuarantinedPath) evaluates to true
			return nil, fmt.Errorf("%w: working directory violates quarantine boundary: %w", ErrSecurityViolation, err)
		}
	}
	if cmd.Dir != "" {
		info, err := os.Stat(cmd.Dir)
		if err != nil || !info.IsDir() {
			return nil, fmt.Errorf("%w: directory '%s' does not exist or is not a directory", ErrWorkingDirectory, cmd.Dir)
		}
	}

	// 4. Setup Timeout & Context
	timeout := cmd.Timeout
	if timeout <= 0 {
		timeout = r.DefaultTimeout
	}
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	maxBytes := cmd.MaxOutputBytes
	if maxBytes <= 0 {
		maxBytes = r.MaxDefaultBytes
	}

	// 5. Build exec.Cmd (Direct invocation without shell wrapper)
	execCmd := exec.CommandContext(execCtx, execName, cmd.Args...)
	if cmd.Dir != "" {
		execCmd.Dir = cmd.Dir
	}
	if len(cmd.Env) > 0 {
		execCmd.Env = append(os.Environ(), cmd.Env...)
	}

	// Platform-specific process tree setup (process groups on Unix)
	prepareCommandPlatform(execCmd)

	// Custom Cancel hook to terminate the entire process tree on both Windows and Unix
	execCmd.Cancel = func() error {
		return killProcessTree(execCmd)
	}

	// WaitDelay ensures child processes that hold pipes open do not hang indefinitely
	execCmd.WaitDelay = 2 * time.Second

	// 6. Setup Output Buffers and Streaming
	var stdoutBuf, stderrBuf limitedBuffer
	stdoutBuf.limit = maxBytes
	stderrBuf.limit = maxBytes

	var stdoutWriter io.Writer = &stdoutBuf
	if cmd.Stdout != nil {
		stdoutWriter = io.MultiWriter(&stdoutBuf, cmd.Stdout)
	}
	execCmd.Stdout = stdoutWriter

	var stderrWriter io.Writer = &stderrBuf
	if cmd.Stderr != nil {
		stderrWriter = io.MultiWriter(&stderrBuf, cmd.Stderr)
	}
	execCmd.Stderr = stderrWriter

	start := time.Now()
	err := execCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start command %s: %w", execName, err)
	}

	pid := 0
	if execCmd.Process != nil {
		pid = execCmd.Process.Pid
	}

	waitErr := execCmd.Wait()
	duration := time.Since(start)

	timedOut := false
	if execCtx.Err() == context.DeadlineExceeded || (ctx.Err() == nil && waitErr != nil && strings.Contains(waitErr.Error(), "killed")) {
		timedOut = true
	}

	exitCode := 0
	if waitErr != nil {
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	outStr := stdoutBuf.String()
	errStr := stderrBuf.String()
	combStr := outStr
	if errStr != "" {
		if combStr != "" {
			combStr += "\n" + errStr
		} else {
			combStr = errStr
		}
	}

	res := &RunResult{
		Command:         execName,
		Args:            cmd.Args,
		Dir:             cmd.Dir,
		ExitCode:        exitCode,
		Stdout:          outStr,
		Stderr:          errStr,
		CombinedOutput:  combStr,
		Duration:        duration,
		PID:             pid,
		TimedOut:        timedOut,
		OutputTruncated: stdoutBuf.truncated || stderrBuf.truncated,
	}

	if timedOut {
		return res, ErrCommandTimeout
	}
	if waitErr != nil && exitCode != 0 {
		return res, fmt.Errorf("command %s exited with code %d: %s", execName, exitCode, strings.TrimSpace(errStr))
	}

	return res, nil
}

// limitedBuffer caps output capture to prevent OOM
type limitedBuffer struct {
	buf       bytes.Buffer
	limit     int64
	truncated bool
}

func (b *limitedBuffer) Write(p []byte) (n int, err error) {
	if b.limit <= 0 {
		return b.buf.Write(p)
	}
	if int64(b.buf.Len()) >= b.limit {
		b.truncated = true
		return len(p), nil
	}
	remaining := b.limit - int64(b.buf.Len())
	if int64(len(p)) > remaining {
		b.truncated = true
		n, _ = b.buf.Write(p[:remaining])
		return len(p), nil
	}
	return b.buf.Write(p)
}

func (b *limitedBuffer) String() string {
	return b.buf.String()
}

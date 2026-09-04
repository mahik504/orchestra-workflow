package runner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// MockCommandRunner enables deterministic unit testing without spawning OS processes.
type MockCommandRunner struct {
	mu           sync.Mutex
	RecordedCmds []Command
	Handlers     map[string]func(cmd Command) (*RunResult, error)
	DefaultResp  *RunResult
	DefaultErr   error
	LookPathMap  map[string]string
}

// NewMockCommandRunner creates an initialized mock command runner.
func NewMockCommandRunner() *MockCommandRunner {
	return &MockCommandRunner{
		RecordedCmds: make([]Command, 0),
		Handlers:     make(map[string]func(cmd Command) (*RunResult, error)),
		LookPathMap: map[string]string{
			"npm":  "C:\\Program Files\\nodejs\\npm.cmd",
			"pnpm": "C:\\Program Files\\nodejs\\pnpm.cmd",
			"yarn": "C:\\Program Files\\nodejs\\yarn.cmd",
			"git":  "C:\\Program Files\\Git\\bin\\git.exe",
			"npx":  "C:\\Program Files\\nodejs\\npx.cmd",
			"node": "C:\\Program Files\\nodejs\\node.exe",
			"bun":  "C:\\Users\\User\\.bun\\bin\\bun.exe",
			"bunx": "C:\\Users\\User\\.bun\\bin\\bunx.exe",
		},
		DefaultResp: &RunResult{
			ExitCode: 0,
			Stdout:   "OK",
			Duration: 5 * time.Millisecond,
		},
	}
}

// LookPath returns the mocked executable path or ErrCommandNotFound.
func (m *MockCommandRunner) LookPath(file string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	trimmed := strings.TrimSpace(file)
	if trimmed == "" {
		return "", ErrInvalidExecutableName
	}
	if strings.ContainsRune(trimmed, '\x00') || strings.ContainsAny(trimmed, "\r\n") {
		return "", fmt.Errorf("%w: executable path contains control characters", ErrSecurityViolation)
	}

	if path, ok := m.LookPathMap[trimmed]; ok {
		return path, nil
	}
	// Case-insensitive lookup fallback
	for k, v := range m.LookPathMap {
		if strings.EqualFold(k, trimmed) {
			return v, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrCommandNotFound, trimmed)
}

// Run executes the mocked command, recording invocation and matching handlers.
func (m *MockCommandRunner) Run(ctx context.Context, cmd Command) (*RunResult, error) {
	m.mu.Lock()
	m.RecordedCmds = append(m.RecordedCmds, cmd)
	m.mu.Unlock()

	// Check context cancellation
	if ctx.Err() != nil {
		return &RunResult{
			Command:  cmd.Name,
			Args:     cmd.Args,
			Dir:      cmd.Dir,
			ExitCode: -1,
			TimedOut: true,
		}, ErrCommandTimeout
	}

	fullCmd := cmd.Name
	if len(cmd.Args) > 0 {
		fullCmd += " " + strings.Join(cmd.Args, " ")
	}

	m.mu.Lock()
	var handler func(cmd Command) (*RunResult, error)

	// 1. Exact match
	if h, ok := m.Handlers[fullCmd]; ok {
		handler = h
	} else {
		// 2. Prefix match
		for k, h := range m.Handlers {
			if strings.HasPrefix(fullCmd, k) {
				handler = h
				break
			}
		}
	}
	defaultResp := m.DefaultResp
	defaultErr := m.DefaultErr
	m.mu.Unlock()

	if handler != nil {
		return handler(cmd)
	}

	if defaultErr != nil {
		return nil, defaultErr
	}

	resp := *defaultResp
	resp.Command = cmd.Name
	resp.Args = cmd.Args
	resp.Dir = cmd.Dir
	return &resp, nil
}

// RegisterHandler registers a custom handler for an exact or prefix command string.
func (m *MockCommandRunner) RegisterHandler(pattern string, h func(cmd Command) (*RunResult, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Handlers[pattern] = h
}

// LastCommand returns the most recently recorded command.
func (m *MockCommandRunner) LastCommand() (Command, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.RecordedCmds) == 0 {
		return Command{}, false
	}
	return m.RecordedCmds[len(m.RecordedCmds)-1], true
}

// Reset clears all recorded commands.
func (m *MockCommandRunner) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RecordedCmds = make([]Command, 0)
}

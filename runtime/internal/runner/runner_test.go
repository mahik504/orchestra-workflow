package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestOSCommandRunner_LookPath(t *testing.T) {
	r := NewOSCommandRunner()

	// Valid command lookup
	path, err := r.LookPath("git")
	if err != nil {
		t.Logf("git not found in PATH: %v (skipping positive assertion)", err)
	} else if path == "" {
		t.Errorf("expected non-empty path for git")
	}

	// Empty command lookup
	_, err = r.LookPath("")
	if !errors.Is(err, ErrInvalidExecutableName) {
		t.Errorf("expected ErrInvalidExecutableName, got %v", err)
	}

	// Null byte security violation
	_, err = r.LookPath("git\x00")
	if !errors.Is(err, ErrSecurityViolation) {
		t.Errorf("expected ErrSecurityViolation for null byte, got %v", err)
	}

	// Non-existent command
	_, err = r.LookPath("non_existent_binary_xyz_12345")
	if !errors.Is(err, ErrCommandNotFound) {
		t.Errorf("expected ErrCommandNotFound, got %v", err)
	}
}

func TestOSCommandRunner_SecuritySanitization(t *testing.T) {
	r := NewOSCommandRunner()
	ctx := context.Background()

	// Empty executable name
	_, err := r.Run(ctx, Command{Name: ""})
	if !errors.Is(err, ErrInvalidExecutableName) {
		t.Errorf("expected ErrInvalidExecutableName, got %v", err)
	}

	// Shell metacharacters in executable name
	for _, badName := range []string{"cmd;rm", "sh|cat", "app&echo", "foo`bar`", "test$var"} {
		_, err := r.Run(ctx, Command{Name: badName})
		if !errors.Is(err, ErrSecurityViolation) {
			t.Errorf("expected ErrSecurityViolation for %q, got %v", badName, err)
		}
	}

	// Null byte in arguments
	_, err = r.Run(ctx, Command{
		Name: "git",
		Args: []string{"version", "foo\x00bar"},
	})
	if !errors.Is(err, ErrSecurityViolation) {
		t.Errorf("expected ErrSecurityViolation for null byte in arg, got %v", err)
	}
}

func TestOSCommandRunner_QuarantineBoundaryEnforcement(t *testing.T) {
	r := NewOSCommandRunner()
	ctx := context.Background()

	// Working dir inside quarantined skills_library
	quarantinedDir := filepath.Join(os.TempDir(), "test_skills_library_violation")
	_ = os.MkdirAll(quarantinedDir, 0755)
	defer os.RemoveAll(quarantinedDir)

	_, err := r.Run(ctx, Command{
		Name: "git",
		Args: []string{"version"},
		Dir:  quarantinedDir,
	})
	if !errors.Is(err, ErrSecurityViolation) {
		t.Errorf("expected ErrSecurityViolation for quarantined working dir, got %v", err)
	}
}

func TestOSCommandRunner_LiveExecution(t *testing.T) {
	r := NewOSCommandRunner()
	ctx := context.Background()

	// Run git version (or go version)
	res, err := r.Run(ctx, Command{
		Name: "git",
		Args: []string{"version"},
	})
	if err != nil {
		t.Skipf("git command failed or not present: %v", err)
	}

	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if !strings.Contains(strings.ToLower(res.Stdout), "git version") {
		t.Errorf("expected stdout to contain 'git version', got %q", res.Stdout)
	}
	if res.Duration <= 0 {
		t.Errorf("expected non-zero duration")
	}
}

func TestOSCommandRunner_Timeout(t *testing.T) {
	r := NewOSCommandRunner()
	ctx := context.Background()

	// Command with 1ms timeout
	_, err := r.Run(ctx, Command{
		Name:    "git",
		Args:    []string{"log", "--all"},
		Timeout: 1 * time.Nanosecond,
	})
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
}

func TestOSCommandRunner_OutputTruncation(t *testing.T) {
	r := NewOSCommandRunner()
	ctx := context.Background()

	res, err := r.Run(ctx, Command{
		Name:           "git",
		Args:           []string{"version"},
		MaxOutputBytes: 4, // truncate after 4 bytes
	})
	if err != nil {
		t.Skipf("git command failed: %v", err)
	}

	if len(res.Stdout) > 4 {
		t.Errorf("expected max 4 bytes output, got %d bytes: %q", len(res.Stdout), res.Stdout)
	}
	if !res.OutputTruncated {
		t.Errorf("expected OutputTruncated to be true")
	}
}

func TestMockCommandRunner_Deterministic(t *testing.T) {
	m := NewMockCommandRunner()
	ctx := context.Background()

	// Register specific handler
	m.RegisterHandler("npm install --save gsap", func(cmd Command) (*RunResult, error) {
		return &RunResult{
			Command:  cmd.Name,
			Args:     cmd.Args,
			ExitCode: 0,
			Stdout:   "added 1 package",
		}, nil
	})

	res, err := m.Run(ctx, Command{
		Name: "npm",
		Args: []string{"install", "--save", "gsap"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout != "added 1 package" {
		t.Errorf("expected 'added 1 package', got %q", res.Stdout)
	}

	// Verify command recording
	last, ok := m.LastCommand()
	if !ok {
		t.Fatalf("expected command to be recorded")
	}
	if last.Name != "npm" || len(last.Args) != 3 {
		t.Errorf("unexpected recorded command: %v", last)
	}

	// Reset
	m.Reset()
	if _, ok := m.LastCommand(); ok {
		t.Errorf("expected empty records after Reset")
	}
}

func TestMockCommandRunner_Concurrency(t *testing.T) {
	m := NewMockCommandRunner()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _ = m.Run(ctx, Command{
				Name: "pnpm",
				Args: []string{"add", "test-pkg"},
			})
			_, _ = m.LookPath("pnpm")
		}(i)
	}
	wg.Wait()

	if len(m.RecordedCmds) != 50 {
		t.Errorf("expected 50 recorded commands, got %d", len(m.RecordedCmds))
	}
}

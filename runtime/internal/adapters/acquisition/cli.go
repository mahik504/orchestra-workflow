package acquisition

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// CLIAdapter implements AcquisitionAdapter for ephemeral on-demand CLI executions
type CLIAdapter struct {
	runner         runner.CommandRunner
	defaultTimeout time.Duration
}

// NewCLIAdapter creates an initialized CLIAdapter
func NewCLIAdapter(r runner.CommandRunner) *CLIAdapter {
	if r == nil {
		r = runner.NewOSCommandRunner()
	}
	timeout := 60 * time.Second
	if os.Getenv("CI") != "" {
		timeout = 5 * time.Second
	}
	return &CLIAdapter{
		runner:         r,
		defaultTimeout: timeout,
	}
}

func (a *CLIAdapter) Name() string {
	return "cli"
}

func (a *CLIAdapter) CanHandle(method string) bool {
	lower := strings.ToLower(strings.TrimSpace(method))
	return lower == "cli" || lower == "npx"
}

// validateArgs enforces shell safety and delegates anti-global enforcement to canonical CheckGlobalInstallSafety
func (a *CLIAdapter) validateArgs(args []string, dest string) error {
	for _, arg := range args {
		if strings.ContainsAny(arg, "|;&$`\n\r><") {
			return fmt.Errorf("%w: argument %q contains prohibited shell metacharacters", ErrCommandInjectionRisk, arg)
		}
	}
	return CheckGlobalInstallSafety(args, dest)
}

// detectRunner determines the ephemeral package runner (pnpm dlx, bunx, yarn dlx, or npx --yes)
func (a *CLIAdapter) detectRunner(workingDir string) (binary string, prefixArgs []string) {
	if workingDir != "" {
		if _, err := os.Stat(filepath.Join(workingDir, "pnpm-lock.yaml")); err == nil {
			if _, err := a.runner.LookPath("pnpm"); err == nil {
				return "pnpm", []string{"dlx"}
			}
		}
		if _, err := os.Stat(filepath.Join(workingDir, "bun.lockb")); err == nil {
			if _, err := a.runner.LookPath("bunx"); err == nil {
				return "bunx", []string{}
			}
		}
		if _, err := os.Stat(filepath.Join(workingDir, "yarn.lock")); err == nil {
			if _, err := a.runner.LookPath("yarn"); err == nil {
				return "yarn", []string{"dlx"}
			}
		}
	}

	// Fallback to npx with non-interactive flag
	return "npx", []string{"--yes"}
}

// Acquire executes the CLI tool ephemerally without modifying package.json
func (a *CLIAdapter) Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error) {
	start := time.Now()

	if res == nil || res.ID == "" {
		return nil, ErrResourceNotFound
	}
	if !a.CanHandle(res.AcquisitionMethod) {
		return nil, fmt.Errorf("%w: resource '%s' requires acquisition method '%s', not cli/npx", ErrResourceNotAllowed, res.ID, res.AcquisitionMethod)
	}

	// Working directory validation
	if dest != "" {
		if err := resources.CheckQuarantineBoundary(dest); err != nil {
			return nil, fmt.Errorf("%w: destination '%s' violates quarantine: %w", ErrQuarantineViolation, dest, err)
		}
	}

	// Resolve package / command to execute
	pkgSpec := res.ID
	if alias, ok := DefaultPackageAliases[res.ID]; ok {
		pkgSpec = alias
	}

	// In live OS runner mode, project-scoped CLI commands require a project with package.json
	if _, isMock := a.runner.(*runner.MockCommandRunner); !isMock {
		if dest != "" {
			pkgJSONPath := filepath.Join(dest, "package.json")
			if _, err := os.Stat(pkgJSONPath); err != nil {
				if os.IsNotExist(err) {
					return nil, fmt.Errorf("%w: %s", ErrPackageJSONNotFound, pkgJSONPath)
				}
			}
		}
	}

	// Enforce shell metacharacter and anti-global policy on destination and package specifier
	if err := a.validateArgs([]string{pkgSpec}, dest); err != nil {
		return nil, err
	}

	runnerBin, prefixArgs := a.detectRunner(dest)

	execArgs := append([]string{}, prefixArgs...)
	execArgs = append(execArgs, pkgSpec)

	// Double-check full execution arguments against anti-global filter
	if err := CheckGlobalInstallSafety(execArgs, dest); err != nil {
		return nil, err
	}

	cmd := runner.Command{
		Name:    runnerBin,
		Args:    execArgs,
		Dir:     dest,
		Timeout: a.defaultTimeout,
	}

	runRes, err := a.runner.Run(ctx, cmd)
	if err != nil {
		stdout := ""
		stderr := ""
		if runRes != nil {
			stdout = runRes.Stdout
			stderr = runRes.Stderr
		}
		return nil, fmt.Errorf("ephemeral CLI execution failed for %s: %w (stdout: %s, stderr: %s)", pkgSpec, err, stdout, stderr)
	}

	// Compute hash of execution output for tamper-evident provenance
	outBytes := []byte(runRes.Stdout)
	if len(outBytes) == 0 {
		outBytes = []byte(pkgSpec + ":" + runRes.CombinedOutput)
	}
	h := sha256.Sum256(outBytes)
	outHash := hex.EncodeToString(h[:])

	return &AcquisitionResult{
		ResourceID:        res.ID,
		AdapterName:       a.Name(),
		AcquisitionMethod: "cli",
		SourceURL:         res.CanonicalURL,
		PackageName:       pkgSpec,
		VersionOrSHA:      outHash[:16],
		SHA256Hash:        outHash,
		InstalledPath:     "", // Ephemeral: does not pollute workspace disk
		ResolvedTarget:    pkgSpec,
		AlreadyInstalled:  false,
		Ephemeral:         true,
		Duration:          time.Since(start),
		ExecutedCommand:   cmd.Name + " " + strings.Join(cmd.Args, " "),
		Output:            runRes.Stdout,
		Metadata: map[string]string{
			"runner":    runnerBin,
			"ephemeral": "true",
		},
	}, nil
}

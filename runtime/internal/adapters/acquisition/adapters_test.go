package acquisition

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// Helper to create a temporary test project workspace
func createTestProject(t *testing.T, pkgJSONContent string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "orchestra_acq_test_*")
	if err != nil {
		t.Fatalf("failed to create temp test dir: %v", err)
	}

	if pkgJSONContent != "" {
		pkgPath := filepath.Join(dir, "package.json")
		if err := os.WriteFile(pkgPath, []byte(pkgJSONContent), 0644); err != nil {
			t.Fatalf("failed to write package.json: %v", err)
		}
	}

	return dir
}

// 1. Adversarial Anti-Global Install Enforcement
func TestNPMAdapter_AntiGlobalInstallBlocked(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewNPMAdapter(mockRunner, nil, nil)
	ctx := context.Background()

	res := &resources.Resource{
		ID:                "gsap",
		Name:              "GSAP",
		AcquisitionMethod: "npm",
	}

	// 1.1 Banned CLI flags in arguments
	bannedFlags := []string{
		"-g",
		"--global",
		"-global",
		"--location=global",
		"-g=true",
		"--global=true",
		"global",
		"--prefix=/usr/local",
		"--prefix=C:\\Windows",
	}

	for _, flag := range bannedFlags {
		err := CheckGlobalInstallSafety([]string{flag}, "C:\\projects\\my-app")
		if !errors.Is(err, ErrGlobalInstallBlocked) {
			t.Errorf("expected ErrGlobalInstallBlocked for flag %q, got %v", flag, err)
		}
	}

	// 1.2 System-wide destination paths
	systemDests := []string{
		"",
		"C:\\",
		"C:\\Windows",
		"C:\\Windows\\System32",
		"C:\\Program Files",
		"C:\\Program Files (x86)",
		"/usr",
		"/usr/local",
		"/etc",
		"C:\\Users\\User\\AppData\\Roaming\\npm",
	}

	for _, dest := range systemDests {
		_, err := adapter.Acquire(ctx, res, dest)
		if !errors.Is(err, ErrGlobalInstallBlocked) {
			t.Errorf("expected ErrGlobalInstallBlocked for system destination %q, got %v", dest, err)
		}
	}

	// Ensure no command was ever dispatched to runner
	if len(mockRunner.RecordedCmds) > 0 {
		t.Errorf("security violation: runner executed %d commands during global install attempts", len(mockRunner.RecordedCmds))
	}
}

// 2. Project-Scoped Conditional Installation
func TestNPMAdapter_ConditionalInstall(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewNPMAdapter(mockRunner, nil, nil)
	ctx := context.Background()

	res := &resources.Resource{
		ID:                "gsap",
		Name:              "GSAP",
		AcquisitionMethod: "npm",
	}

	// Case A: Package already declared and installed in node_modules
	pkgJSON := `{"name": "test-app", "dependencies": {"gsap": "^3.12.5"}}`
	workdirA := createTestProject(t, pkgJSON)
	defer os.RemoveAll(workdirA)

	// Create fake node_modules/gsap/package.json
	modDir := filepath.Join(workdirA, "node_modules", "gsap")
	_ = os.MkdirAll(modDir, 0755)
	_ = os.WriteFile(filepath.Join(modDir, "package.json"), []byte(`{"name": "gsap", "version": "3.12.5"}`), 0644)

	resA, err := adapter.Acquire(ctx, res, workdirA)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !resA.AlreadyInstalled {
		t.Errorf("expected AlreadyInstalled to be true")
	}
	if resA.VersionOrSHA != "3.12.5" {
		t.Errorf("expected version 3.12.5, got %q", resA.VersionOrSHA)
	}
	if len(mockRunner.RecordedCmds) != 0 {
		t.Errorf("expected zero commands executed for already installed package, got %d", len(mockRunner.RecordedCmds))
	}

	// Case B: Missing package -> should trigger project-scoped installation
	workdirB := createTestProject(t, `{"name": "test-app", "dependencies": {}}`)
	defer os.RemoveAll(workdirB)

	resB, err := adapter.Acquire(ctx, res, workdirB)
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if resB.AlreadyInstalled {
		t.Errorf("expected AlreadyInstalled to be false")
	}
	if len(mockRunner.RecordedCmds) != 1 {
		t.Fatalf("expected 1 command executed, got %d", len(mockRunner.RecordedCmds))
	}

	lastCmd, _ := mockRunner.LastCommand()
	// Should be strictly project-scoped (e.g. pnpm add or npm install --save)
	if lastCmd.Dir != workdirB {
		t.Errorf("expected command working dir %q, got %q", workdirB, lastCmd.Dir)
	}
	for _, arg := range lastCmd.Args {
		if arg == "-g" || arg == "--global" {
			t.Errorf("detected global flag in project install command: %v", lastCmd.Args)
		}
	}
}

// 3. Package Manager Auto-Detection
func TestNPMAdapter_PackageManagerDetection(t *testing.T) {
	ctx := context.Background()

	res := &resources.Resource{
		ID:                "lenis",
		Name:              "Lenis",
		AcquisitionMethod: "npm",
	}

	// 3.1 pnpm detection via pnpm-lock.yaml
	t.Run("pnpm-lock.yaml", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		adapter := NewNPMAdapter(mockRunner, nil, nil)
		workdir := createTestProject(t, `{"name": "test-app"}`)
		defer os.RemoveAll(workdir)
		_ = os.WriteFile(filepath.Join(workdir, "pnpm-lock.yaml"), []byte("lockfileVersion: 5.4"), 0644)

		_, err := adapter.Acquire(ctx, res, workdir)
		if err != nil {
			t.Fatalf("Acquire failed: %v", err)
		}
		cmd, _ := mockRunner.LastCommand()
		if cmd.Name != "pnpm" || len(cmd.Args) == 0 || cmd.Args[0] != "add" {
			t.Errorf("expected pnpm add, got %s %v", cmd.Name, cmd.Args)
		}
	})

	// 3.2 yarn detection via yarn.lock
	t.Run("yarn.lock", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		adapter := NewNPMAdapter(mockRunner, nil, nil)
		workdir := createTestProject(t, `{"name": "test-app"}`)
		defer os.RemoveAll(workdir)
		_ = os.WriteFile(filepath.Join(workdir, "yarn.lock"), []byte("# yarn lockfile"), 0644)

		_, err := adapter.Acquire(ctx, res, workdir)
		if err != nil {
			t.Fatalf("Acquire failed: %v", err)
		}
		cmd, _ := mockRunner.LastCommand()
		if cmd.Name != "yarn" || len(cmd.Args) == 0 || cmd.Args[0] != "add" {
			t.Errorf("expected yarn add, got %s %v", cmd.Name, cmd.Args)
		}
	})

	// 3.3 packageManager field in package.json
	t.Run("packageManager field", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		adapter := NewNPMAdapter(mockRunner, nil, nil)
		workdir := createTestProject(t, `{"name": "test-app", "packageManager": "pnpm@8.15.0"}`)
		defer os.RemoveAll(workdir)

		_, err := adapter.Acquire(ctx, res, workdir)
		if err != nil {
			t.Fatalf("Acquire failed: %v", err)
		}
		cmd, _ := mockRunner.LastCommand()
		if cmd.Name != "pnpm" {
			t.Errorf("expected pnpm from packageManager field, got %s", cmd.Name)
		}
	})
}

// 4. Resource Catalog Validation & Alias Resolution
func TestNPMAdapter_CatalogValidation(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	// 4.1 Alias resolution (r3f -> @react-three/fiber)
	adapter := NewNPMAdapter(mockRunner, nil, nil)
	workdir := createTestProject(t, `{"name": "test-app"}`)
	defer os.RemoveAll(workdir)

	resR3F := &resources.Resource{
		ID:                "r3f",
		Name:              "React Three Fiber",
		AcquisitionMethod: "npm",
	}

	result, err := adapter.Acquire(ctx, resR3F, workdir)
	if err != nil {
		t.Fatalf("Acquire r3f failed: %v", err)
	}
	if result.PackageName != "@react-three/fiber" {
		t.Errorf("expected alias @react-three/fiber, got %s", result.PackageName)
	}

	// 4.2 Non-NPM method rejected
	resGit := &resources.Resource{
		ID:                "my-repo",
		AcquisitionMethod: "git",
	}
	_, err = adapter.Acquire(ctx, resGit, workdir)
	if !errors.Is(err, ErrResourceNotAllowed) {
		t.Errorf("expected ErrResourceNotAllowed for git resource, got %v", err)
	}

	// 4.3 Missing package.json in destination
	emptyDir, _ := os.MkdirTemp("", "empty_dir_*")
	defer os.RemoveAll(emptyDir)
	_, err = adapter.Acquire(ctx, resR3F, emptyDir)
	if !errors.Is(err, ErrPackageJSONNotFound) {
		t.Errorf("expected ErrPackageJSONNotFound, got %v", err)
	}
}

// 5. GitAdapter: Shallow Clone, Commit SHA Verification, and Concurrency
func TestGitAdapter_ShallowCloneAndPinning(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	pinnedSHA := "abc1234567890abcdef1234567890abcdef12345"

	// Mock git rev-parse HEAD returning expected pinned SHA
	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{
			ExitCode: 0,
			Stdout:   pinnedSHA + "\n",
		}, nil
	})

	locker := NewHybridLocker(t.TempDir())
	adapter := NewGitAdapter(mockRunner, locker)
	adapter.SetDefaultOptions(GitCloneOptions{
		Depth:     1,
		PinnedSHA: pinnedSHA,
	})

	workdir, _ := os.MkdirTemp("", "git_test_*")
	defer os.RemoveAll(workdir)
	dest := filepath.Join(workdir, "cloned_repo")

	res := &resources.Resource{
		ID:                "taste-design",
		Name:              "Taste Design",
		SourceRepository:  "https://github.com/example/taste-design.git",
		AcquisitionMethod: "git",
	}

	result, err := adapter.Acquire(ctx, res, dest)
	if err != nil {
		t.Fatalf("Git Acquire failed: %v", err)
	}

	if result.VersionOrSHA != pinnedSHA {
		t.Errorf("expected commit SHA %s, got %s", pinnedSHA, result.VersionOrSHA)
	}

	// Check that clone was shallow (--depth 1)
	firstCmd := mockRunner.RecordedCmds[0]
	hasDepth := false
	for i, arg := range firstCmd.Args {
		if arg == "--depth" && i+1 < len(firstCmd.Args) && firstCmd.Args[i+1] == "1" {
			hasDepth = true
			break
		}
	}
	if !hasDepth {
		t.Errorf("expected git clone to have --depth 1, got args: %v", firstCmd.Args)
	}
}

func TestGitAdapter_CommitSHAMismatchRejection(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	expectedSHA := "1111111111111111111111111111111111111111"
	actualWrongSHA := "9999999999999999999999999999999999999999"

	// Mock git rev-parse HEAD returning WRONG SHA
	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{
			ExitCode: 0,
			Stdout:   actualWrongSHA + "\n",
		}, nil
	})

	locker := NewHybridLocker(t.TempDir())
	adapter := NewGitAdapter(mockRunner, locker)
	adapter.SetDefaultOptions(GitCloneOptions{
		Depth:     1,
		PinnedSHA: expectedSHA,
	})

	workdir, _ := os.MkdirTemp("", "git_mismatch_*")
	defer os.RemoveAll(workdir)
	dest := filepath.Join(workdir, "tampered_repo")

	res := &resources.Resource{
		ID:                "untrusted-repo",
		SourceRepository:  "https://github.com/example/untrusted.git",
		AcquisitionMethod: "git",
	}

	_, err := adapter.Acquire(ctx, res, dest)
	if !errors.Is(err, ErrCommitSHAMismatch) {
		t.Errorf("expected ErrCommitSHAMismatch, got %v", err)
	}

	// Verify directory was removed on mismatch
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Errorf("expected cloned directory to be wiped on SHA mismatch, but it still exists")
	}
}

func TestGitAdapter_ConcurrencyLocking(t *testing.T) {
	locker := NewHybridLocker(t.TempDir())
	ctx := context.Background()

	unlock1, err := locker.Lock(ctx, "res-lock-test")
	if err != nil {
		t.Fatalf("first lock failed: %v", err)
	}

	var wg sync.WaitGroup
	secondLockAcquired := false

	wg.Add(1)
	go func() {
		defer wg.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		unlock2, err := locker.Lock(timeoutCtx, "res-lock-test")
		if err == nil {
			secondLockAcquired = true
			unlock2()
		}
	}()

	wg.Wait()
	if secondLockAcquired {
		t.Errorf("concurrency violation: second lock acquired before first unlocked")
	}

	unlock1()
}

// 6. CLIAdapter: Anti-Global & Ephemeral Execution
func TestCLIAdapter_AntiGlobalAndEphemeral(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewCLIAdapter(mockRunner)
	ctx := context.Background()

	res := &resources.Resource{
		ID:                "skillui",
		Name:              "SkillUI CLI",
		AcquisitionMethod: "npx",
	}

	workdir, _ := os.MkdirTemp("", "cli_test_*")
	defer os.RemoveAll(workdir)

	// Valid ephemeral execution
	result, err := adapter.Acquire(ctx, res, workdir)
	if err != nil {
		t.Fatalf("CLI Acquire failed: %v", err)
	}

	if !result.Ephemeral {
		t.Errorf("expected Ephemeral: true")
	}
	if result.InstalledPath != "" {
		t.Errorf("expected empty InstalledPath for ephemeral execution, got %q", result.InstalledPath)
	}

	cmd, _ := mockRunner.LastCommand()
	if cmd.Name != "npx" || len(cmd.Args) == 0 || cmd.Args[0] != "--yes" {
		t.Errorf("expected npx --yes, got %s %v", cmd.Name, cmd.Args)
	}

	// Injection rejection
	badRes := &resources.Resource{
		ID:                "malicious;rm -rf /",
		AcquisitionMethod: "npx",
	}
	_, err = adapter.Acquire(ctx, badRes, workdir)
	if !errors.Is(err, ErrCommandInjectionRisk) {
		t.Errorf("expected ErrCommandInjectionRisk for shell metacharacters, got %v", err)
	}
}

// 7. WebAdapter: HTTP Fetch, SSRF, SHA256 & Offline Fallback
func TestWebAdapter_FetchAndSHA256(t *testing.T) {
	// Create test HTTP server
	expectedContent := `{"inspiration": "awwwards-top-design", "palette": ["#000", "#fff"]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, expectedContent)
	}))
	defer ts.Close()

	adapter := NewWebAdapter(false)
	adapter.SetAllowPrivateIP(true) // Allow loopback for httptest
	ctx := context.Background()

	res := &resources.Resource{
		ID:                "test-web-ref",
		CanonicalURL:      ts.URL,
		AcquisitionMethod: "web_fetch",
	}

	result, err := adapter.Acquire(ctx, res, "")
	if err != nil {
		t.Fatalf("Web Acquire failed: %v", err)
	}

	if result.Output != expectedContent {
		t.Errorf("expected output %q, got %q", expectedContent, result.Output)
	}
	if result.SHA256Hash == "" {
		t.Errorf("expected non-empty SHA256Hash")
	}

	// Scheme check (reject file://)
	resFile := &resources.Resource{
		ID:                "file-ref",
		CanonicalURL:      "file:///etc/passwd",
		AcquisitionMethod: "web_fetch",
	}
	_, err = adapter.Acquire(ctx, resFile, "")
	if !errors.Is(err, ErrUnsupportedURLScheme) {
		t.Errorf("expected ErrUnsupportedURLScheme for file://, got %v", err)
	}
}

func TestWebAdapter_OfflineFixtureFallback(t *testing.T) {
	adapter := NewWebAdapter(true) // offline mode enabled
	ctx := context.Background()

	// awwwards is registered in CuratedSourceFixtures
	res := &resources.Resource{
		ID:                "awwwards",
		CanonicalURL:      "https://www.awwwards.com",
		AcquisitionMethod: "web_fetch",
	}

	result, err := adapter.Acquire(ctx, res, "")
	if err != nil {
		t.Fatalf("Offline fixture fallback failed: %v", err)
	}

	if result.Metadata["source"] != "offline_fixture" {
		t.Errorf("expected offline_fixture source, got %q", result.Metadata["source"])
	}
	if result.SHA256Hash == "" {
		t.Errorf("expected valid SHA256 hash for fixture")
	}
}

// 8. MCPAdapter: Validation & Anti-Global Blocking
func TestMCPAdapter_ConfigValidation(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewMCPAdapter(mockRunner)

	// 8.1 Banned servers rejected
	bannedServers := []string{"higgsfield", "magicui", "21st", "open-design"}
	for _, srv := range bannedServers {
		err := adapter.ValidateConfig(srv, &MCPServerConfig{Command: "npx"})
		if !errors.Is(err, ErrMCPRejected) {
			t.Errorf("expected ErrMCPRejected for banned server %q, got %v", srv, err)
		}
	}

	// 8.2 Global install blocked in MCP config
	globalCfg := &MCPServerConfig{
		Command: "npm",
		Args:    []string{"install", "-g", "@playwright/mcp"},
	}
	err := adapter.ValidateConfig("playwright", globalCfg)
	if !errors.Is(err, ErrGlobalInstallBlocked) {
		t.Errorf("expected ErrGlobalInstallBlocked for -g in args, got %v", err)
	}

	// 8.3 Shell injection blocked
	injectionCfg := &MCPServerConfig{
		Command: "npx; rm -rf /",
	}
	err = adapter.ValidateConfig("playwright", injectionCfg)
	if !errors.Is(err, ErrCommandInjectionRisk) {
		t.Errorf("expected ErrCommandInjectionRisk for shell operator, got %v", err)
	}

	// 8.4 Valid config passes
	validCfg := &MCPServerConfig{
		Command: "npx",
		Args:    []string{"-y", "@upstash/context7-mcp"},
	}
	err = adapter.ValidateConfig("context7", validCfg)
	if err != nil {
		t.Errorf("expected valid config to pass, got %v", err)
	}
}

// 9. AdapterRegistry integration
func TestAdapterRegistry(t *testing.T) {
	reg := NewAdapterRegistry()
	mockRunner := runner.NewMockCommandRunner()

	reg.RegisterAdapter(NewNPMAdapter(mockRunner, nil, nil))
	reg.RegisterAdapter(NewGitAdapter(mockRunner, nil))
	reg.RegisterAdapter(NewCLIAdapter(mockRunner))
	reg.RegisterAdapter(NewWebAdapter(true))
	reg.RegisterAdapter(NewMCPAdapter(mockRunner))

	methods := []string{"npm", "git", "cli", "npx", "web_fetch", "mcp_install"}
	for _, m := range methods {
		adapter, err := reg.GetAdapterForMethod(m)
		if err != nil {
			t.Errorf("expected adapter for method %q, got error: %v", m, err)
		} else if !adapter.CanHandle(m) {
			t.Errorf("adapter %s reported CanHandle(%q) = false", adapter.Name(), m)
		}
	}
}

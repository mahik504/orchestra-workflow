package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/orchestra-v3/internal/adapters/acquisition"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// Helper to create a temporary test project workspace with a package.json
func createTestProject(t *testing.T, pkgJSONContent string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "orch_challenger_m3_*")
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

// =========================================================================
// TEST SUITE 1: ANTI-GLOBAL INSTALL DEFENSES
// =========================================================================

// TestChallenger_GlobalInstallFlags_SafetyCheck empirically tests CheckGlobalInstallSafety
func TestChallenger_GlobalInstallFlags_SafetyCheck(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		dest        string
		shouldBlock bool
	}{
		{"flag -g", []string{"-g"}, "C:\\safe\\dir", true},
		{"flag --global", []string{"--global"}, "C:\\safe\\dir", true},
		{"flag -global", []string{"-global"}, "C:\\safe\\dir", true},
		{"flag --location=global", []string{"--location=global"}, "C:\\safe\\dir", true},
		{"flag --location=\"global\"", []string{"--location=\"global\""}, "C:\\safe\\dir", true},
		{"flag --location='global'", []string{"--location='global'"}, "C:\\safe\\dir", true},
		{"flag --location global (2 tokens)", []string{"--location", "global"}, "C:\\safe\\dir", true},
		{"flag -g=true", []string{"-g=true"}, "C:\\safe\\dir", true},
		{"flag --global=true", []string{"--global=true"}, "C:\\safe\\dir", true},
		{"flag -g=false", []string{"-g=false"}, "C:\\safe\\dir", true},
		{"flag --global=false", []string{"--global=false"}, "C:\\safe\\dir", true},
		{"flag -g=1", []string{"-g=1"}, "C:\\safe\\dir", true},
		{"flag --global=1", []string{"--global=1"}, "C:\\safe\\dir", true},
		{"flag -g=yes", []string{"-g=yes"}, "C:\\safe\\dir", true},
		{"keyword global single", []string{"global"}, "C:\\safe\\dir", true},
		{"yarn global add (3 tokens)", []string{"yarn", "global", "add"}, "C:\\safe\\dir", true},
		{"yarn global add (1 token)", []string{"yarn global add"}, "C:\\safe\\dir", true},
		{"--prefix=/usr", []string{"--prefix=/usr"}, "C:\\safe\\dir", true},
		{"-prefix=/usr", []string{"-prefix=/usr"}, "C:\\safe\\dir", true},
		{"--prefix=/usr/local", []string{"--prefix=/usr/local"}, "C:\\safe\\dir", true},
		{"--prefix=C:\\Windows", []string{"--prefix=C:\\Windows"}, "C:\\safe\\dir", true},
		{"--prefix /usr (2 tokens)", []string{"--prefix", "/usr"}, "C:\\safe\\dir", true},
		{"--prefix C:\\Windows (2 tokens)", []string{"--prefix", "C:\\Windows"}, "C:\\safe\\dir", true},
		{"dest root /", nil, "/", true},
		{"dest root C:\\", nil, "C:\\", true},
		{"dest /usr", nil, "/usr", true},
		{"dest /usr/local", nil, "/usr/local", true},
		{"dest C:\\Windows", nil, "C:\\Windows", true},
		{"dest C:\\Program Files", nil, "C:\\Program Files", true},
		{"dest npm AppData", nil, "C:\\Users\\User\\AppData\\Roaming\\npm", true},
		{"dest empty", nil, "", true},
		{"legitimate project-scoped pnpm add", []string{"add", "gsap"}, "C:\\safe\\dir", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := acquisition.CheckGlobalInstallSafety(tc.args, tc.dest)
			if tc.shouldBlock {
				if err == nil {
					t.Errorf("[VULNERABILITY] Expected ErrGlobalInstallBlocked for %s, but got nil (args=%v, dest=%s)", tc.name, tc.args, tc.dest)
				} else if !errors.Is(err, acquisition.ErrGlobalInstallBlocked) {
					t.Errorf("[DEFECT] Expected ErrGlobalInstallBlocked for %s, got: %v", tc.name, err)
				} else {
					t.Logf("[PASS] Correctly blocked %s: %v", tc.name, err)
				}
			} else {
				if err != nil {
					t.Errorf("[FALSE POSITIVE] Expected no error for %s, got: %v", tc.name, err)
				}
			}
		})
	}
}

// TestChallenger_CLIAdapter_GlobalFlags empirically tests CLIAdapter against all global flags
func TestChallenger_CLIAdapter_GlobalFlags(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := acquisition.NewCLIAdapter(mockRunner)
	ctx := context.Background()

	workdir := createTestProject(t, `{"name": "cli-test"}`)
	defer os.RemoveAll(workdir)

	globalVariants := []struct {
		name string
		id   string
	}{
		{"flag -g", "-g"},
		{"flag --global", "--global"},
		{"flag -global", "-global"},
		{"flag --location=global", "--location=global"},
		{"flag -g=true", "-g=true"},
		{"flag --global=true", "--global=true"},
		{"keyword global", "global"},
		{"yarn global add", "yarn global add"},
		{"--prefix=/usr", "--prefix=/usr"},
		{"-prefix=/usr", "-prefix=/usr"},
		{"--prefix=/usr/local", "--prefix=/usr/local"},
		{"--prefix=C:\\Windows", "--prefix=C:\\Windows"},
	}

	for _, gv := range globalVariants {
		t.Run(gv.name, func(t *testing.T) {
			mockRunner.Reset()
			res := &resources.Resource{
				ID:                gv.id,
				AcquisitionMethod: "cli",
			}

			_, err := adapter.Acquire(ctx, res, workdir)
			if err == nil {
				t.Errorf("[CRITICAL VULNERABILITY] CLIAdapter dispatched command without error for global flag %q! Executed: %v", gv.id, mockRunner.RecordedCmds)
			} else if !errors.Is(err, acquisition.ErrGlobalInstallBlocked) {
				t.Errorf("[DEFECT] CLIAdapter rejected %q but returned %v instead of ErrGlobalInstallBlocked", gv.id, err)
			} else {
				t.Logf("[PASS] CLIAdapter blocked %q: %v", gv.id, err)
			}

			if len(mockRunner.RecordedCmds) > 0 {
				t.Errorf("[CRITICAL] Command was dispatched to runner despite global flag %q: %v", gv.id, mockRunner.RecordedCmds)
			}
		})
	}
}

// TestChallenger_NPMAdapter_GlobalFlags empirically tests NPMAdapter against all global flags
func TestChallenger_NPMAdapter_GlobalFlags(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := acquisition.NewNPMAdapter(mockRunner, nil, nil)
	ctx := context.Background()

	workdir := createTestProject(t, `{"name": "npm-test", "dependencies": {}}`)
	defer os.RemoveAll(workdir)

	globalVariants := []struct {
		name string
		id   string
	}{
		{"flag -g", "-g"},
		{"flag --global", "--global"},
		{"flag -global", "-global"},
		{"flag --location=global", "--location=global"},
		{"flag -g=true", "-g=true"},
		{"flag --global=true", "--global=true"},
		{"keyword global", "global"},
		{"yarn global add", "yarn global add"},
		{"--prefix=/usr", "--prefix=/usr"},
		{"-prefix=/usr", "-prefix=/usr"},
		{"--prefix=/usr/local", "--prefix=/usr/local"},
		{"--prefix=C:\\Windows", "--prefix=C:\\Windows"},
	}

	for _, gv := range globalVariants {
		t.Run(gv.name, func(t *testing.T) {
			mockRunner.Reset()
			res := &resources.Resource{
				ID:                gv.id,
				AcquisitionMethod: "npm",
			}

			_, err := adapter.Acquire(ctx, res, workdir)
			if err == nil {
				t.Errorf("[CRITICAL VULNERABILITY] NPMAdapter allowed global install argument %q! Executed: %v", gv.id, mockRunner.RecordedCmds)
			} else if !errors.Is(err, acquisition.ErrGlobalInstallBlocked) {
				t.Errorf("[DEFECT] NPMAdapter rejected %q but returned %v instead of ErrGlobalInstallBlocked", gv.id, err)
			} else {
				t.Logf("[PASS] NPMAdapter blocked %q: %v", gv.id, err)
			}

			if len(mockRunner.RecordedCmds) > 0 {
				t.Errorf("[CRITICAL] Command dispatched to runner despite global flag %q: %v", gv.id, mockRunner.RecordedCmds)
			}
		})
	}
}

// =========================================================================
// TEST SUITE 2: COMMAND RUNNER INJECTION & METAHARACTERS
// =========================================================================

func TestChallenger_OSCommandRunner_InjectionDefenses(t *testing.T) {
	r := runner.NewOSCommandRunner()
	ctx := context.Background()

	// 2.1 Metacharacters in Executable Name
	metacharNames := []string{
		"git;rm",
		"sh|cat",
		"app&calc",
		"var$TEST",
		"foo`id`",
		"out>file",
		"in<file",
		"git\x00evil",
		"git\r\ncmd",
	}

	for _, name := range metacharNames {
		t.Run("exec_name_"+name, func(t *testing.T) {
			_, err := r.Run(ctx, runner.Command{Name: name})
			if !errors.Is(err, runner.ErrSecurityViolation) {
				t.Errorf("expected ErrSecurityViolation for executable %q, got: %v", name, err)
			} else {
				t.Logf("[PASS] Executable %q correctly rejected: %v", name, err)
			}
		})
	}

	// 2.2 Boundary Control Characters in Executable Name (TrimSpace bypass check)
	boundaryControlNames := []string{
		"git\n",
		"git\r",
		"\ngit",
		"\rgit",
	}

	for _, name := range boundaryControlNames {
		t.Run("boundary_control_"+name, func(t *testing.T) {
			_, err := r.Run(ctx, runner.Command{Name: name})
			// The contract says control characters in executable name must trigger ErrSecurityViolation.
			// Because strings.TrimSpace runs first, these characters are stripped and execution proceeds!
			if !errors.Is(err, runner.ErrSecurityViolation) {
				t.Errorf("[VULNERABILITY] Boundary control char in %q was stripped by TrimSpace, did not return ErrSecurityViolation: %v", name, err)
			} else {
				t.Logf("[PASS] Boundary control char in %q rejected", name)
			}
		})
	}

	// 2.3 Null bytes in arguments
	nullArgTests := []struct {
		name string
		args []string
	}{
		{"trailing null", []string{"version", "foo\x00"}},
		{"leading null", []string{"\x00version"}},
		{"embedded null", []string{"commit", "-m", "hello\x00world"}},
		{"null byte alone", []string{"\x00"}},
	}

	for _, tc := range nullArgTests {
		t.Run("null_arg_"+tc.name, func(t *testing.T) {
			_, err := r.Run(ctx, runner.Command{
				Name: "git",
				Args: tc.args,
			})
			if !errors.Is(err, runner.ErrSecurityViolation) {
				t.Errorf("expected ErrSecurityViolation for null byte in args %v, got: %v", tc.args, err)
			} else {
				t.Logf("[PASS] Null byte in args %v rejected: %v", tc.args, err)
			}
		})
	}

	// 2.4 Shell metacharacters in arguments (Direct OS execution test)
	t.Run("shell_metacharacters_in_args_direct_execution", func(t *testing.T) {
		tempDir := t.TempDir()
		sentinelFile := filepath.Join(tempDir, "pwned.txt")

		// If a shell was invoked, "; touch <sentinelFile>" would create the sentinel file.
		_, _ = r.Run(ctx, runner.Command{
			Name: "git",
			Args: []string{"version", ";", "touch", sentinelFile, "&", "echo", "pwned"},
			Dir:  tempDir,
		})

		if _, err := os.Stat(sentinelFile); err == nil {
			t.Errorf("[CRITICAL INJECTION VULNERABILITY] Shell metacharacters executed! Sentinel file %s was created!", sentinelFile)
		} else {
			t.Logf("[PASS] Shell metacharacters treated as discrete literal arguments; no command injection occurred.")
		}
	})

	// 2.5 Null byte and injection in cmd.Env
	t.Run("env_null_byte_sanitization", func(t *testing.T) {
		cmd := runner.Command{
			Name: "git",
			Args: []string{"version"},
			Env:  []string{"EVIL=\x00MALICIOUS"},
		}

		_, err := r.Run(ctx, cmd)
		if !errors.Is(err, runner.ErrSecurityViolation) {
			t.Errorf("[DEFECT] Null byte in cmd.Env was not rejected by CommandRunner with ErrSecurityViolation. Got: %v", err)
		} else {
			t.Logf("[PASS] Null byte in cmd.Env properly rejected with ErrSecurityViolation")
		}
	})

	// 2.6 Control characters and newlines in cmd.Env
	t.Run("env_control_chars_sanitization", func(t *testing.T) {
		cmd := runner.Command{
			Name: "git",
			Args: []string{"version"},
			Env:  []string{"EVIL=hello\r\nworld"},
		}

		_, err := r.Run(ctx, cmd)
		if !errors.Is(err, runner.ErrSecurityViolation) {
			t.Errorf("[DEFECT] Control chars in cmd.Env were not rejected with ErrSecurityViolation. Got: %v", err)
		} else {
			t.Logf("[PASS] Control chars in cmd.Env properly rejected")
		}
	})

	// 2.7 LookPath control character bypass
	lookPathTests := []string{
		"git\n",
		"git\r",
		"\ngit",
		"\rgit",
		"git\x00evil",
	}
	for _, lp := range lookPathTests {
		t.Run("lookpath_control_"+lp, func(t *testing.T) {
			_, err := r.LookPath(lp)
			if !errors.Is(err, runner.ErrSecurityViolation) {
				t.Errorf("[VULNERABILITY] LookPath(%q) did not return ErrSecurityViolation! Got: %v", lp, err)
			} else {
				t.Logf("[PASS] LookPath(%q) correctly rejected with ErrSecurityViolation", lp)
			}
		})
	}
}

// =========================================================================
// TEST SUITE 3: WORKING DIRECTORY QUARANTINE ENFORCEMENT
// =========================================================================

func TestChallenger_QuarantineEnforcement(t *testing.T) {
	r := runner.NewOSCommandRunner()
	ctx := context.Background()

	// 3.1 Existing directory matching quarantine boundary
	quarantineTempDir := filepath.Join(os.TempDir(), "orch_test_skills_library_challenger")
	_ = os.MkdirAll(quarantineTempDir, 0755)
	defer os.RemoveAll(quarantineTempDir)

	_, err := r.Run(ctx, runner.Command{
		Name: "git",
		Args: []string{"version"},
		Dir:  quarantineTempDir,
	})

	if err == nil {
		t.Fatalf("[CRITICAL VULNERABILITY] CommandRunner executed in quarantined directory %s without error!", quarantineTempDir)
	}

	t.Logf("Quarantine execution returned error: %v", err)

	// Verify error chain: MUST unwrap to resources.ErrQuarantinedPath
	if !errors.Is(err, resources.ErrQuarantinedPath) {
		t.Errorf("[DEFECT - BROKEN ERROR CHAIN] Command error does NOT match resources.ErrQuarantinedPath via errors.Is! Error: %v", err)
	} else {
		t.Logf("[PASS] Command error unwraps to resources.ErrQuarantinedPath")
	}

	// 3.2 Non-existent directory matching quarantine pattern
	// Check whether quarantine boundary is checked BEFORE os.Stat
	nonExistentQuarantine := filepath.Join(os.TempDir(), "nonexistent_skills_library_dir", "sub")
	_, errNonExist := r.Run(ctx, runner.Command{
		Name: "git",
		Args: []string{"version"},
		Dir:  nonExistentQuarantine,
	})

	if errors.Is(errNonExist, runner.ErrWorkingDirectory) {
		t.Errorf("[DEFECT - INVERTED CHECK ORDER] Quarantine check was BYPASSED because directory does not exist; returned ErrWorkingDirectory instead of quarantine violation")
	} else if errors.Is(errNonExist, resources.ErrQuarantinedPath) || errors.Is(errNonExist, runner.ErrSecurityViolation) {
		t.Logf("[PASS] Quarantine boundary validated before filesystem stat check")
	}

	// 3.3 CheckQuarantineBoundary direct validation
	traversalPaths := []string{
		"skills_library",
		"SKILLS_LIBRARY",
		"Skills_Library",
		"skills-library",
		"curated_catalog/quarantine",
		"skills~1",
		"C:\\Users\\mockuser\\.gemini\\config\\skills_library",
	}

	for _, p := range traversalPaths {
		t.Run("quarantine_check_"+p, func(t *testing.T) {
			checkErr := resources.CheckQuarantineBoundary(p)
			if checkErr == nil {
				t.Errorf("[CRITICAL QUARANTINE BYPASS] CheckQuarantineBoundary returned nil for %q!", p)
			} else if !errors.Is(checkErr, resources.ErrQuarantinedPath) {
				t.Errorf("Expected ErrQuarantinedPath for %q, got: %v", p, checkErr)
			} else {
				t.Logf("[PASS] CheckQuarantineBoundary correctly rejected %q", p)
			}
		})
	}
}

// =========================================================================
// TEST SUITE 4: SUPPLY CHAIN INTEGRITY & SSRF DEFENSES
// =========================================================================

func TestChallenger_SupplyChainAndWebDefenses(t *testing.T) {
	ctx := context.Background()

	// 4.1 Git Shallow Clone & Commit SHA Verification
	t.Run("git_commit_sha_mismatch_wipes_disk", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		actualTamperedSHA := "badbeef1234567890badbeef1234567890badbeef"
		expectedPinnedSHA := "cafebabe1234567890cafebabe1234567890cafeb"

		mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
			return &runner.RunResult{
				ExitCode: 0,
				Stdout:   actualTamperedSHA + "\n",
			}, nil
		})

		locker := acquisition.NewHybridLocker(t.TempDir())
		gitAdapter := acquisition.NewGitAdapter(mockRunner, locker)
		gitAdapter.SetDefaultOptions(acquisition.GitCloneOptions{
			Depth:     1,
			PinnedSHA: expectedPinnedSHA,
		})

		dest := filepath.Join(t.TempDir(), "cloned_repo")
		res := &resources.Resource{
			ID:                "test-git-res",
			SourceRepository:  "https://github.com/example/repo.git",
			AcquisitionMethod: "git",
		}

		_, err := gitAdapter.Acquire(ctx, res, dest)
		if !errors.Is(err, acquisition.ErrCommitSHAMismatch) {
			t.Errorf("Expected ErrCommitSHAMismatch, got: %v", err)
		}
		if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
			t.Errorf("[SECURITY VULNERABILITY] Destination directory was NOT wiped following SHA mismatch!")
		} else {
			t.Logf("[PASS] Directory wiped following SHA mismatch")
		}
	})

	// 4.2 Web Adapter SSRF Protection against Loopback and Metadata Addresses
	t.Run("web_ssrf_private_ip_blocked", func(t *testing.T) {
		webAdapter := acquisition.NewWebAdapter(false) // live mode with private IP blocking
		ssrfTargets := []string{
			"http://127.0.0.1:8080/secret",
			"http://169.254.169.254/latest/meta-data/",
			"http://10.0.0.1/admin",
			"http://192.168.1.1/router",
		}

		for _, target := range ssrfTargets {
			res := &resources.Resource{
				ID:                "ssrf-target",
				CanonicalURL:      target,
				AcquisitionMethod: "web_fetch",
			}
			_, err := webAdapter.Acquire(ctx, res, "")
			if !errors.Is(err, acquisition.ErrSSRFDetected) {
				t.Errorf("Expected ErrSSRFDetected for %s, got: %v", target, err)
			} else {
				t.Logf("[PASS] Blocked SSRF to %s: %v", target, err)
			}
		}
	})

	// 4.3 Web Adapter 10MB Payload Limit
	t.Run("web_payload_limit_enforced", func(t *testing.T) {
		// Server emitting 11MB
		oversizedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chunk := strings.Repeat("A", 1024*1024) // 1MB chunk
			for i := 0; i < 11; i++ {
				_, _ = fmt.Fprint(w, chunk)
			}
		}))
		defer oversizedServer.Close()

		webAdapter := acquisition.NewWebAdapter(false)
		webAdapter.SetAllowPrivateIP(true) // allow loopback for test server
		res := &resources.Resource{
			ID:                "oversized-web",
			CanonicalURL:      oversizedServer.URL,
			AcquisitionMethod: "web_fetch",
		}

		_, err := webAdapter.Acquire(ctx, res, "")
		if !errors.Is(err, acquisition.ErrPayloadTooLarge) {
			t.Errorf("Expected ErrPayloadTooLarge for >10MB payload, got: %v", err)
		} else {
			t.Logf("[PASS] Web adapter rejected payload exceeding 10MB: %v", err)
		}
	})
}

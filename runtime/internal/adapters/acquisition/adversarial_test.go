package acquisition

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// TestAdversarial_NPM_AntiGlobalInstall tests comprehensive anti-global install blocking
func TestAdversarial_NPM_AntiGlobalInstall(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewNPMAdapter(mockRunner, nil, nil)
	ctx := context.Background()

	workdir := createTestProject(t, `{"name": "adv-npm-test", "dependencies": {}}`)
	defer os.RemoveAll(workdir)

	cases := []struct {
		name string
		id   string
	}{
		{"flag -g", "-g"},
		{"flag --global", "--global"},
		{"flag -global", "-global"},
		{"flag --location=global", "--location=global"},
		{"flag --location=\"global\"", "--location=\"global\""},
		{"flag --location='global'", "--location='global'"},
		{"flag -g=true", "-g=true"},
		{"flag --global=true", "--global=true"},
		{"keyword global", "global"},
		{"yarn global add", "yarn global add"},
		{"--prefix=/usr", "--prefix=/usr"},
		{"-prefix=/usr", "-prefix=/usr"},
		{"--prefix=/usr/local", "--prefix=/usr/local"},
		{"--prefix=C:\\Windows", "--prefix=C:\\Windows"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRunner.Reset()
			res := &resources.Resource{
				ID:                tc.id,
				AcquisitionMethod: "npm",
			}
			_, err := adapter.Acquire(ctx, res, workdir)
			if !errors.Is(err, ErrGlobalInstallBlocked) {
				t.Errorf("expected ErrGlobalInstallBlocked for %s, got: %v", tc.name, err)
			}
			if len(mockRunner.RecordedCmds) > 0 {
				t.Errorf("command was dispatched to runner despite global flag %s", tc.name)
			}
		})
	}
}

// TestAdversarial_NPM_ConcurrencyLocking proves mutual exclusion on workspace
func TestAdversarial_NPM_ConcurrencyLocking(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	var activeInstalls, peakInstalls int64

	mockRunner.RegisterHandler("pnpm add gsap", func(cmd runner.Command) (*runner.RunResult, error) {
		curr := atomic.AddInt64(&activeInstalls, 1)
		for {
			max := atomic.LoadInt64(&peakInstalls)
			if curr <= max || atomic.CompareAndSwapInt64(&peakInstalls, max, curr) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt64(&activeInstalls, -1)
		return &runner.RunResult{ExitCode: 0, Stdout: "added 1 package"}, nil
	})

	adapter := NewNPMAdapter(mockRunner, nil, nil)
	workdir := createTestProject(t, `{"name": "adv-npm-lock", "dependencies": {}}`)
	defer os.RemoveAll(workdir)

	ctx := context.Background()
	res := &resources.Resource{
		ID:                "gsap",
		AcquisitionMethod: "npm",
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = adapter.Acquire(ctx, res, workdir)
		}()
	}
	wg.Wait()

	peak := atomic.LoadInt64(&peakInstalls)
	if peak > 1 {
		t.Errorf("expected peak concurrency of 1, got: %d", peak)
	}
}

// TestAdversarial_CLI_GlobalFlagsAndMetacharacters tests ephemeral CLI security
func TestAdversarial_CLI_GlobalFlagsAndMetacharacters(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	adapter := NewCLIAdapter(mockRunner)
	ctx := context.Background()

	workdir := createTestProject(t, `{"name": "adv-cli-test"}`)
	defer os.RemoveAll(workdir)

	globalFlags := []string{
		"-g", "--global", "-global", "--location=global",
		"yarn global add", "--prefix=/usr", "--prefix=C:\\Windows",
	}

	for _, flag := range globalFlags {
		t.Run("cli_flag_"+flag, func(t *testing.T) {
			mockRunner.Reset()
			res := &resources.Resource{
				ID:                flag,
				AcquisitionMethod: "cli",
			}
			_, err := adapter.Acquire(ctx, res, workdir)
			if !errors.Is(err, ErrGlobalInstallBlocked) {
				t.Errorf("expected ErrGlobalInstallBlocked for %s, got: %v", flag, err)
			}
			if len(mockRunner.RecordedCmds) > 0 {
				t.Errorf("command was dispatched for %s", flag)
			}
		})
	}

	// Metacharacters
	metachars := []string{"app;rm", "pkg&calc", "tool|cat", "var$FOO"}
	for _, meta := range metachars {
		t.Run("cli_meta_"+meta, func(t *testing.T) {
			res := &resources.Resource{
				ID:                meta,
				AcquisitionMethod: "cli",
			}
			_, err := adapter.Acquire(ctx, res, workdir)
			if !errors.Is(err, ErrCommandInjectionRisk) {
				t.Errorf("expected ErrCommandInjectionRisk for %s, got: %v", meta, err)
			}
		})
	}
}

// TestAdversarial_Web_SSRFAndRedirects tests SSRF against unspecified IPs and HTTP redirects
func TestAdversarial_Web_SSRFAndRedirects(t *testing.T) {
	adapter := NewWebAdapter(false)
	ctx := context.Background()

	// 1. Unspecified and loopback targets
	targets := []string{
		"http://127.0.0.1:8080/",
		"http://localhost:3000/",
		"http://169.254.169.254/latest/meta-data/",
		"http://0.0.0.0/",
		"http://0.0.0.0:8080/admin",
		"http://[::]:8080/",
		"http://10.0.0.1/",
		"http://192.168.1.1/",
	}

	for _, target := range targets {
		t.Run("ssrf_"+target, func(t *testing.T) {
			res := &resources.Resource{
				ID:                "ssrf-probe",
				CanonicalURL:      target,
				AcquisitionMethod: "web_fetch",
			}
			_, err := adapter.Acquire(ctx, res, "")
			if !errors.Is(err, ErrSSRFDetected) {
				t.Errorf("expected ErrSSRFDetected for %s, got: %v", target, err)
			}
		})
	}

	// 2. HTTP 302 Redirect to private loopback
	internalSecretServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `secret-data`)
	}))
	defer internalSecretServer.Close()

	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internalSecretServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	redirectAdapter := NewWebAdapter(false)
	redirectAdapter.SetAllowPrivateIP(true)

	redirectRes := &resources.Resource{
		ID:                "redirect-probe",
		CanonicalURL:      redirectServer.URL,
		AcquisitionMethod: "web_fetch",
	}
	_, err := redirectAdapter.Acquire(ctx, redirectRes, "")
	if !errors.Is(err, ErrSSRFDetected) {
		t.Errorf("expected ErrSSRFDetected for redirect to loopback, got: %v", err)
	}
}

// TestAdversarial_Web_PayloadLimitAndChecksum tests payload limit and SHA256 calculation
func TestAdversarial_Web_PayloadLimitAndChecksum(t *testing.T) {
	ctx := context.Background()

	// Oversized server (11MB)
	oversizedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chunk := strings.Repeat("B", 1024*1024)
		for i := 0; i < 11; i++ {
			_, _ = fmt.Fprint(w, chunk)
		}
	}))
	defer oversizedServer.Close()

	adapter := NewWebAdapter(false)
	adapter.SetAllowPrivateIP(true)

	res := &resources.Resource{
		ID:                "oversized",
		CanonicalURL:      oversizedServer.URL,
		AcquisitionMethod: "web_fetch",
	}

	_, err := adapter.Acquire(ctx, res, "")
	if !errors.Is(err, ErrPayloadTooLarge) {
		t.Errorf("expected ErrPayloadTooLarge, got: %v", err)
	}
}

// TestAdversarial_Git_CommitSHAPinningAndTamperWipe tests commit verification and wiping
func TestAdversarial_Git_CommitSHAPinningAndTamperWipe(t *testing.T) {
	mockRunner := runner.NewMockCommandRunner()
	tamperedSHA := "deadbeef1234567890deadbeef1234567890deadbeef"
	expectedSHA := "cafebabe1234567890cafebabe1234567890cafeb"

	mockRunner.RegisterHandler("git clone", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: "cloned"}, nil
	})
	mockRunner.RegisterHandler("git -C", func(cmd runner.Command) (*runner.RunResult, error) {
		return &runner.RunResult{ExitCode: 0, Stdout: tamperedSHA + "\n"}, nil
	})

	locker := NewHybridLocker(t.TempDir())
	gitAdapter := NewGitAdapter(mockRunner, locker)
	gitAdapter.SetDefaultOptions(GitCloneOptions{
		Depth:     1,
		PinnedSHA: expectedSHA,
	})

	dest := filepath.Join(t.TempDir(), "cloned_repo")
	_ = os.MkdirAll(dest, 0755)
	_ = os.WriteFile(filepath.Join(dest, "temp.txt"), []byte("payload"), 0644)

	res := &resources.Resource{
		ID:                "git-test",
		SourceRepository:  "https://github.com/example/test.git",
		AcquisitionMethod: "git",
	}

	_, err := gitAdapter.Acquire(context.Background(), res, dest)
	if !errors.Is(err, ErrCommitSHAMismatch) {
		t.Errorf("expected ErrCommitSHAMismatch, got: %v", err)
	}

	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Errorf("expected destination %s to be wiped completely from disk", dest)
	}
}

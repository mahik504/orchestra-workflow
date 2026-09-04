package acquisition

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// GitCloneOptions configures git clone operations
type GitCloneOptions struct {
	Depth       int           // Default: 1
	BranchOrTag string        // Optional branch or tag
	PinnedSHA   string        // Expected commit SHA
	LockTimeout time.Duration // Default: 30s
}

// GitAdapter implements AcquisitionAdapter for Git repositories
type GitAdapter struct {
	gitBinary   string
	runner      runner.CommandRunner
	locker      ResourceLocker
	defaultOpts GitCloneOptions
}

// NewGitAdapter creates an initialized GitAdapter
func NewGitAdapter(r runner.CommandRunner, locker ResourceLocker) *GitAdapter {
	if r == nil {
		r = runner.NewOSCommandRunner()
	}
	if locker == nil {
		locker = NewHybridLocker(filepath.Join(os.TempDir(), "orchestra_git_locks"))
	}
	return &GitAdapter{
		gitBinary: "git",
		runner:    r,
		locker:    locker,
		defaultOpts: GitCloneOptions{
			Depth:       1,
			LockTimeout: 30 * time.Second,
		},
	}
}

func (g *GitAdapter) Name() string {
	return "git"
}

func (g *GitAdapter) CanHandle(method string) bool {
	return strings.EqualFold(method, "git")
}

// SetDefaultOptions updates the default clone options
func (g *GitAdapter) SetDefaultOptions(opts GitCloneOptions) {
	g.defaultOpts = opts
}

// Acquire performs shallow clone, commit SHA verification, locking, and quarantine enforcement
func (g *GitAdapter) Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error) {
	start := time.Now()

	if res == nil || res.ID == "" {
		return nil, ErrResourceNotFound
	}
	if !g.CanHandle(res.AcquisitionMethod) {
		return nil, fmt.Errorf("%w: resource '%s' requires acquisition method '%s', not git", ErrResourceNotAllowed, res.ID, res.AcquisitionMethod)
	}

	repoURL := res.SourceRepository
	if repoURL == "" {
		repoURL = res.CanonicalURL
	}
	if strings.TrimSpace(repoURL) == "" {
		return nil, fmt.Errorf("%w: repository URL is empty for resource %s", ErrInvalidRepositoryURL, res.ID)
	}

	// 1. Quarantine & Safety Pre-checks
	if err := resources.CheckQuarantineBoundary(repoURL); err != nil {
		return nil, fmt.Errorf("%w: repository URL violates quarantine: %v", ErrQuarantineViolation, err)
	}
	if err := resources.CheckQuarantineBoundary(dest); err != nil {
		return nil, fmt.Errorf("%w: destination '%s' violates quarantine: %v", ErrQuarantineViolation, dest, err)
	}
	if strings.HasPrefix(repoURL, "-") {
		return nil, fmt.Errorf("%w: repository URL %q starts with a dash (flag injection risk)", ErrInvalidRepositoryURL, repoURL)
	}

	// Determine expected/pinned SHA
	pinnedSHA := g.defaultOpts.PinnedSHA
	// If the resource specifies a version that looks like a 40-char or 7-char git SHA, use it
	if len(res.Rationale) >= 7 && isHexString(res.Rationale) {
		pinnedSHA = res.Rationale
	}

	// 2. Concurrency Control via Resource Locking
	unlock, err := g.locker.Lock(ctx, res.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrResourceLockTimeout, err)
	}
	defer unlock()

	// 3. Idempotency & Cache Verification
	gitDir := filepath.Join(dest, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		actualSHA, err := g.getHEADCommitSHA(ctx, dest)
		if err == nil && actualSHA != "" {
			if pinnedSHA == "" || strings.EqualFold(actualSHA, pinnedSHA) {
				contentHash, _ := g.computeDirectoryChecksum(dest)
				return &AcquisitionResult{
					ResourceID:        res.ID,
					AdapterName:       g.Name(),
					AcquisitionMethod: "git",
					SourceURL:         repoURL,
					VersionOrSHA:      actualSHA,
					SHA256Hash:        contentHash,
					InstalledPath:     dest,
					ResolvedTarget:    dest,
					AlreadyInstalled:  true,
					Duration:          time.Since(start),
					Metadata: map[string]string{
						"cache_hit":  "true",
						"commit_sha": actualSHA,
					},
				}, nil
			}
			// SHA mismatch in existing directory: remove and re-clone
			_ = os.RemoveAll(dest)
		}
	}

	// 4. Construct Shallow Clone Command
	depth := g.defaultOpts.Depth
	if depth <= 0 {
		depth = 1
	}

	args := []string{"clone", "--depth", fmt.Sprintf("%d", depth)}
	if g.defaultOpts.BranchOrTag != "" {
		args = append(args, "--branch", g.defaultOpts.BranchOrTag)
	}
	// Safety delimiter -- to prevent flag injection from repoURL or dest
	args = append(args, "--", repoURL, dest)

	cmd := runner.Command{
		Name:    g.gitBinary,
		Args:    args,
		Timeout: 2 * time.Minute,
	}

	runRes, err := g.runner.Run(ctx, cmd)
	if err != nil {
		stdout := ""
		stderr := ""
		if runRes != nil {
			stdout = runRes.Stdout
			stderr = runRes.Stderr
		}
		return nil, fmt.Errorf("git clone failed for %s: %w (stdout: %s, stderr: %s)", repoURL, err, stdout, stderr)
	}

	// 5. Post-Clone Commit SHA Verification
	actualSHA, err := g.getHEADCommitSHA(ctx, dest)
	if err != nil {
		_ = os.RemoveAll(dest)
		return nil, fmt.Errorf("failed to verify cloned git repository: %w", err)
	}

	if pinnedSHA != "" && !strings.EqualFold(actualSHA, pinnedSHA) {
		// Wipe directory immediately to remove unverified code
		_ = os.RemoveAll(dest)
		return nil, fmt.Errorf("%w: expected SHA %s, but cloned HEAD is %s", ErrCommitSHAMismatch, pinnedSHA, actualSHA)
	}

	// 6. Compute Directory Content Checksum
	contentHash, err := g.computeDirectoryChecksum(dest)
	if err != nil {
		contentHash = actualSHA
	}

	return &AcquisitionResult{
		ResourceID:        res.ID,
		AdapterName:       g.Name(),
		AcquisitionMethod: "git",
		SourceURL:         repoURL,
		VersionOrSHA:      actualSHA,
		SHA256Hash:        contentHash,
		InstalledPath:     dest,
		ResolvedTarget:    dest,
		AlreadyInstalled:  false,
		Duration:          time.Since(start),
		ExecutedCommand:   cmd.Name + " " + strings.Join(cmd.Args, " "),
		Output:            runRes.Stdout,
		Metadata: map[string]string{
			"commit_sha": actualSHA,
			"depth":      fmt.Sprintf("%d", depth),
		},
	}, nil
}

// getHEADCommitSHA queries the repository for the current HEAD commit hash
func (g *GitAdapter) getHEADCommitSHA(ctx context.Context, repoDir string) (string, error) {
	cmd := runner.Command{
		Name:    g.gitBinary,
		Args:    []string{"-C", repoDir, "rev-parse", "HEAD"},
		Timeout: 15 * time.Second,
	}
	res, err := g.runner.Run(ctx, cmd)
	if err != nil {
		return "", err
	}
	sha := strings.TrimSpace(res.Stdout)
	if len(sha) < 7 {
		return "", fmt.Errorf("invalid commit SHA received: %q", sha)
	}
	return sha, nil
}

// computeDirectoryChecksum generates a deterministic SHA256 of all files in dir (excluding .git)
func (g *GitAdapter) computeDirectoryChecksum(dir string) (string, error) {
	h := sha256.New()
	var filePaths []string

	err := filepath.Walk(dir, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err == nil {
			filePaths = append(filePaths, rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(filePaths)

	for _, rel := range filePaths {
		abs := filepath.Join(dir, rel)
		f, err := os.Open(abs)
		if err != nil {
			continue
		}
		_, _ = fmt.Fprintf(h, "file:%s\n", filepath.ToSlash(rel))
		_, _ = io.Copy(h, f)
		_ = f.Close()
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func isHexString(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

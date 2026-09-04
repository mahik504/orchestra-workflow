package acquisition

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

// Universal Acquisition Errors
var (
	ErrGlobalInstallBlocked   = errors.New("global installation is strictly blocked by policy: only project-scoped or on-demand installations allowed")
	ErrQuarantineViolation    = errors.New("resource violates quarantine boundary")
	ErrCommitSHAMismatch      = errors.New("git commit SHA verification failed")
	ErrInvalidRepositoryURL   = errors.New("invalid or unsafe repository URL")
	ErrResourceLockTimeout    = errors.New("timed out acquiring resource lock")
	ErrUnsupportedURLScheme   = errors.New("unsupported URL scheme: only http and https are allowed")
	ErrSSRFDetected           = errors.New("target URL resolves to blocked private, loopback, or metadata IP address")
	ErrPayloadTooLarge        = errors.New("web response payload exceeds maximum allowed size")
	ErrContentHashMismatch    = errors.New("content SHA256 hash does not match expected hash")
	ErrCommandInjectionRisk   = errors.New("command arguments contain prohibited shell metacharacters")
	ErrPackageJSONNotFound    = errors.New("package.json not found in target project workspace")
	ErrPackageManagerNotFound = errors.New("no supported package manager (pnpm, yarn, npm) found in PATH")
	ErrResourceNotAllowed     = errors.New("resource is not authorized for this acquisition method")
	ErrResourceRejected       = errors.New("resource status is REJECTED in canonical registry")
	ErrResourceNotFound       = errors.New("resource ID not found in canonical registry")
	ErrInvalidDestinationPath = errors.New("invalid or unsafe project destination path")
	ErrInstallationFailed     = errors.New("resource package installation failed")
	ErrQuarantinedDestination = errors.New("destination violates skill quarantine boundary")
	ErrMCPRejected            = errors.New("MCP server rejected by security or studio policy")
	ErrMCPBinaryNotFound      = errors.New("MCP server command binary not found in PATH")
)

// AcquisitionAdapter defines the canonical interface for all acquisition adapters.
type AcquisitionAdapter interface {
	Name() string
	CanHandle(method string) bool
	Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error)
}

// AcquisitionResult contains detailed telemetry and artifacts from resource acquisition.
type AcquisitionResult struct {
	ResourceID        string            `json:"resource_id"`
	AdapterName       string            `json:"adapter_name"`
	AcquisitionMethod string            `json:"acquisition_method"`
	SourceURL         string            `json:"source_url,omitempty"`
	PackageName       string            `json:"package_name,omitempty"`
	VersionOrSHA      string            `json:"version_or_sha,omitempty"`
	SHA256Hash        string            `json:"sha256_hash,omitempty"`
	InstalledPath     string            `json:"installed_path,omitempty"`
	ResolvedTarget    string            `json:"resolved_target,omitempty"`
	AlreadyInstalled  bool              `json:"already_installed"`
	Duration          time.Duration     `json:"duration"`
	ExecutedCommand   string            `json:"executed_command,omitempty"`
	Output            string            `json:"output,omitempty"`
	Ephemeral         bool              `json:"ephemeral"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// ResourceLocker abstracts concurrency locking to prevent racing parallel acquisitions.
type ResourceLocker interface {
	Lock(ctx context.Context, key string) (func(), error)
}

// HybridLocker combines in-memory per-key mutexes with filesystem locks and lease expiration.
type HybridLocker struct {
	mu       sync.Mutex
	inMem    map[string]chan struct{}
	lockDir  string
	leaseTTL time.Duration
}

// NewHybridLocker constructs a HybridLocker rooted at lockDir.
func NewHybridLocker(lockDir string) *HybridLocker {
	if lockDir == "" {
		lockDir = filepath.Join(os.TempDir(), "orchestra_resource_locks")
	}
	_ = os.MkdirAll(lockDir, 0755)
	return &HybridLocker{
		inMem:    make(map[string]chan struct{}),
		lockDir:  lockDir,
		leaseTTL: 5 * time.Minute,
	}
}

// Lock acquires both memory and filesystem locks for the given resource key.
func (h *HybridLocker) Lock(ctx context.Context, key string) (func(), error) {
	// 1. In-memory semaphore acquisition with context awareness
	h.mu.Lock()
	sem, exists := h.inMem[key]
	if !exists {
		sem = make(chan struct{}, 1)
		h.inMem[key] = sem
	}
	h.mu.Unlock()

	select {
	case sem <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// 2. Filesystem lock with retry and lease recovery
	keyHash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
	lockPath := filepath.Join(h.lockDir, keyHash+".lock")

	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		deadline = time.Now().Add(30 * time.Second)
	}

	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()

	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "pid=%d time=%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
			_ = f.Close()

			return func() {
				_ = os.Remove(lockPath)
				<-sem
			}, nil
		}

		// Check for stale lock
		if info, statErr := os.Stat(lockPath); statErr == nil {
			if time.Since(info.ModTime()) > h.leaseTTL {
				_ = os.Remove(lockPath)
				continue
			}
		}

		select {
		case <-ctx.Done():
			<-sem
			return nil, ctx.Err()
		case t := <-ticker.C:
			if t.After(deadline) {
				<-sem
				return nil, ErrResourceLockTimeout
			}
		}
	}
}

// AdapterRegistry manages the active collection of acquisition adapters.
type AdapterRegistry struct {
	mu       sync.RWMutex
	adapters []AcquisitionAdapter
}

// NewAdapterRegistry creates an empty adapter registry.
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make([]AcquisitionAdapter, 0),
	}
}

// RegisterAdapter adds an adapter to the registry.
func (r *AdapterRegistry) RegisterAdapter(a AcquisitionAdapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters = append(r.adapters, a)
}

// GetAdapterForMethod returns the first registered adapter that handles the given acquisition method.
func (r *AdapterRegistry) GetAdapterForMethod(method string) (AcquisitionAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.adapters {
		if a.CanHandle(method) {
			return a, nil
		}
	}
	return nil, fmt.Errorf("no acquisition adapter registered for method %q", method)
}

// ListAdapters returns a slice of all registered adapters.
func (r *AdapterRegistry) ListAdapters() []AcquisitionAdapter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]AcquisitionAdapter, len(r.adapters))
	copy(res, r.adapters)
	return res
}

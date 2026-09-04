package acquisition

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

// Provenance Errors
var (
	ErrProvenanceEntryInvalid = errors.New("provenance entry failed validation")
	ErrProvenanceNotFound     = errors.New("resource ID not found in provenance store")
	ErrProvenanceStoreCorrupt = errors.New("provenance ledger file is corrupted or invalid")
	ErrIntegrityCheckFailed   = errors.New("provenance integrity verification failed")
)

// ProvenanceEntry models the mandatory 9-field provenance ledger record.
type ProvenanceEntry struct {
	ResourceID          string `json:"resource_id"`
	AcquisitionMethod   string `json:"acquisition_method"`
	SourceURL           string `json:"source_url"`
	VersionOrSHA        string `json:"version_or_sha"`
	SHA256Hash          string `json:"sha256_hash"`
	InstalledPath       string `json:"installed_path"`
	Timestamp           string `json:"timestamp"`
	JustificationTaskID string `json:"justification_task_id"`
	IsQuarantined       bool   `json:"is_quarantined"`
}

// ProvenanceDocument represents the persisted JSON structure at .orchestra/provenance.json.
type ProvenanceDocument struct {
	SchemaVersion string            `json:"schema_version"`
	LastUpdated   string            `json:"last_updated"`
	Entries       []ProvenanceEntry `json:"entries"`
}

// IntegrityIssue details a discrepancy identified during disk verification.
type IntegrityIssue struct {
	ResourceID  string `json:"resource_id"`
	Severity    string `json:"severity"`   // "ERROR" | "WARNING"
	IssueType   string `json:"issue_type"` // "MISSING_FILE" | "HASH_MISMATCH" | "QUARANTINE_BREACH"
	Expected    string `json:"expected"`
	Actual      string `json:"actual"`
	Description string `json:"description"`
}

// IntegrityReport summarizes the comprehensive health and integrity check of all acquired resources.
type IntegrityReport struct {
	TotalChecked     int              `json:"total_checked"`
	PassedCount      int              `json:"passed_count"`
	FailedCount      int              `json:"failed_count"`
	AllValid         bool             `json:"all_valid"`
	Issues           []IntegrityIssue `json:"issues,omitempty"`
	VerificationTime string           `json:"verification_time"`
}

// ProvenanceStore provides thread-safe management and persistence of .orchestra/provenance.json.
type ProvenanceStore struct {
	workspaceRoot string
	storePath     string
	entries       map[string]ProvenanceEntry // keyed by resource_id (case-insensitive)
	order         []string                   // preserve insertion order
	mu            sync.RWMutex
}

// NewProvenanceStore constructs and loads a ProvenanceStore for the given workspaceRoot.
func NewProvenanceStore(workspaceRoot string) (*ProvenanceStore, error) {
	if workspaceRoot == "" {
		workspaceRoot = "."
	}

	orchestraDir := filepath.Join(workspaceRoot, ".orchestra")
	if err := resources.CheckQuarantineBoundary(orchestraDir); err != nil {
		return nil, fmt.Errorf("orchestra directory violates quarantine: %w", err)
	}

	_ = os.MkdirAll(orchestraDir, 0755)
	storePath := filepath.Join(orchestraDir, "provenance.json")

	store := &ProvenanceStore{
		workspaceRoot: workspaceRoot,
		storePath:     storePath,
		entries:       make(map[string]ProvenanceEntry),
		order:         make([]string, 0),
	}

	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %v", ErrProvenanceStoreCorrupt, err)
	}

	return store, nil
}

// load parses existing provenance.json if present on disk.
func (s *ProvenanceStore) load() error {
	data, err := os.ReadFile(s.storePath)
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return nil
	}

	var doc ProvenanceDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, entry := range doc.Entries {
		key := strings.ToLower(entry.ResourceID)
		s.entries[key] = entry
		s.order = append(s.order, entry.ResourceID)
	}

	return nil
}

// ValidateEntry ensures all 9 fields conform to schema requirements.
func ValidateEntry(entry *ProvenanceEntry) error {
	if strings.TrimSpace(entry.ResourceID) == "" {
		return fmt.Errorf("%w: resource_id is required", ErrProvenanceEntryInvalid)
	}
	if strings.TrimSpace(entry.AcquisitionMethod) == "" {
		return fmt.Errorf("%w: acquisition_method is required", ErrProvenanceEntryInvalid)
	}
	if strings.TrimSpace(entry.SourceURL) == "" {
		return fmt.Errorf("%w: source_url is required", ErrProvenanceEntryInvalid)
	}
	if strings.TrimSpace(entry.VersionOrSHA) == "" {
		return fmt.Errorf("%w: version_or_sha is required", ErrProvenanceEntryInvalid)
	}
	if strings.TrimSpace(entry.SHA256Hash) == "" {
		return fmt.Errorf("%w: sha256_hash is required", ErrProvenanceEntryInvalid)
	}
	if strings.TrimSpace(entry.JustificationTaskID) == "" {
		return fmt.Errorf("%w: justification_task_id is required", ErrProvenanceEntryInvalid)
	}
	return nil
}

// Record validates, checks quarantine boundary, and persists a ProvenanceEntry atomically.
func (s *ProvenanceStore) Record(entry ProvenanceEntry) error {
	if err := ValidateEntry(&entry); err != nil {
		return err
	}

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	// Automatic Quarantine Detection: if installed path points to banned skills_library, flag it
	if entry.InstalledPath != "" {
		if err := resources.CheckQuarantineBoundary(entry.InstalledPath); err != nil {
			entry.IsQuarantined = true
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := strings.ToLower(entry.ResourceID)
	if _, exists := s.entries[key]; !exists {
		s.order = append(s.order, entry.ResourceID)
	}
	s.entries[key] = entry

	return s.persistLocked()
}

// persistLocked writes the ledger document to disk atomically via temporary file rename.
func (s *ProvenanceStore) persistLocked() error {
	doc := ProvenanceDocument{
		SchemaVersion: "1.0.0",
		LastUpdated:   time.Now().UTC().Format(time.RFC3339),
		Entries:       make([]ProvenanceEntry, 0, len(s.order)),
	}

	for _, id := range s.order {
		key := strings.ToLower(id)
		if entry, ok := s.entries[key]; ok {
			doc.Entries = append(doc.Entries, entry)
		}
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize provenance document: %w", err)
	}

	tmpPath := s.storePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary provenance file: %w", err)
	}

	if err := os.Rename(tmpPath, s.storePath); err != nil {
		// On Windows, fallback to replace if rename fails
		_ = os.Remove(s.storePath)
		if err := os.Rename(tmpPath, s.storePath); err != nil {
			return fmt.Errorf("failed to atomically rename provenance file: %w", err)
		}
	}

	return nil
}

// GetByResourceID retrieves a recorded provenance entry by its resource ID.
func (s *ProvenanceStore) GetByResourceID(id string) (*ProvenanceEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := strings.ToLower(strings.TrimSpace(id))
	entry, found := s.entries[key]
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrProvenanceNotFound, id)
	}
	return &entry, nil
}

// ListAll returns all recorded provenance entries in insertion order.
func (s *ProvenanceStore) ListAll() ([]ProvenanceEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]ProvenanceEntry, 0, len(s.order))
	for _, id := range s.order {
		key := strings.ToLower(id)
		if entry, ok := s.entries[key]; ok {
			res = append(res, entry)
		}
	}
	return res, nil
}

// ListQuarantined returns all entries flagged as quarantined.
func (s *ProvenanceStore) ListQuarantined() ([]ProvenanceEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []ProvenanceEntry
	for _, id := range s.order {
		key := strings.ToLower(id)
		if entry, ok := s.entries[key]; ok && entry.IsQuarantined {
			res = append(res, entry)
		}
	}
	return res, nil
}

// VerifyIntegrity inspects all recorded entries against the actual disk state.
func (s *ProvenanceStore) VerifyIntegrity(workspaceRoot string) (*IntegrityReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if workspaceRoot == "" {
		workspaceRoot = s.workspaceRoot
	}

	report := &IntegrityReport{
		TotalChecked:     len(s.entries),
		AllValid:         true,
		Issues:           make([]IntegrityIssue, 0),
		VerificationTime: time.Now().UTC().Format(time.RFC3339),
	}

	for _, id := range s.order {
		key := strings.ToLower(id)
		entry := s.entries[key]

		// 1. Ephemeral resource check (installed_path is empty)
		if entry.InstalledPath == "" {
			if len(entry.SHA256Hash) < 8 {
				report.AllValid = false
				report.FailedCount++
				report.Issues = append(report.Issues, IntegrityIssue{
					ResourceID:  entry.ResourceID,
					Severity:    "ERROR",
					IssueType:   "HASH_MISMATCH",
					Expected:    "Valid SHA256 hex string",
					Actual:      entry.SHA256Hash,
					Description: "Ephemeral resource has missing or malformed SHA256 hash",
				})
			} else {
				report.PassedCount++
			}
			continue
		}

		// 2. Resolve target path
		targetPath := entry.InstalledPath
		if !filepath.IsAbs(targetPath) {
			targetPath = filepath.Join(workspaceRoot, targetPath)
		}

		// 3. Quarantine Boundary Audit
		if err := resources.CheckQuarantineBoundary(targetPath); err != nil && !entry.IsQuarantined {
			report.AllValid = false
			report.FailedCount++
			report.Issues = append(report.Issues, IntegrityIssue{
				ResourceID:  entry.ResourceID,
				Severity:    "ERROR",
				IssueType:   "QUARANTINE_BREACH",
				Expected:    "Safe workspace path outside skills_library",
				Actual:      targetPath,
				Description: fmt.Sprintf("Recorded path breaches skill quarantine boundary: %v", err),
			})
			continue
		}

		// 4. File / Directory Existence Check
		info, err := os.Stat(targetPath)
		if err != nil {
			report.AllValid = false
			report.FailedCount++
			report.Issues = append(report.Issues, IntegrityIssue{
				ResourceID:  entry.ResourceID,
				Severity:    "ERROR",
				IssueType:   "MISSING_FILE",
				Expected:    targetPath,
				Actual:      "File not found on disk",
				Description: fmt.Sprintf("Target installed path does not exist: %v", err),
			})
			continue
		}

		// 5. Content Hash Check
		var actualHash string
		if info.IsDir() {
			dirHash, dirErr := computeDirectoryHash(targetPath)
			if dirErr == nil {
				actualHash = dirHash
			}
			// Also allow package.json match if recorded entry matches package.json
			if actualHash != entry.SHA256Hash {
				pkgPath := filepath.Join(targetPath, "package.json")
				if pkgData, pkgErr := os.ReadFile(pkgPath); pkgErr == nil {
					h := sha256.Sum256(pkgData)
					pkgHash := hex.EncodeToString(h[:])
					if strings.EqualFold(pkgHash, entry.SHA256Hash) {
						actualHash = pkgHash
					}
				}
			}
		} else {
			f, err := os.Open(targetPath)
			if err != nil {
				report.AllValid = false
				report.FailedCount++
				report.Issues = append(report.Issues, IntegrityIssue{
					ResourceID:  entry.ResourceID,
					Severity:    "ERROR",
					IssueType:   "MISSING_FILE",
					Expected:    targetPath,
					Actual:      err.Error(),
					Description: "Failed opening target file for hash verification",
				})
				continue
			}
			h := sha256.New()
			_, _ = io.Copy(h, f)
			_ = f.Close()
			actualHash = hex.EncodeToString(h.Sum(nil))
		}

		// If entry hash is provided and doesn't match
		if entry.SHA256Hash != "" && actualHash != "" && !strings.EqualFold(actualHash, entry.SHA256Hash) {
			report.AllValid = false
			report.FailedCount++
			report.Issues = append(report.Issues, IntegrityIssue{
				ResourceID:  entry.ResourceID,
				Severity:    "ERROR",
				IssueType:   "HASH_MISMATCH",
				Expected:    entry.SHA256Hash,
				Actual:      actualHash,
				Description: fmt.Sprintf("Content hash %s does not match recorded hash %s", actualHash, entry.SHA256Hash),
			})
			continue
		}

		report.PassedCount++
	}

	// 6. Scan workspace for unlisted/backdoor files
	recordedPaths := make(map[string]bool)
	recordedDirs := make([]string, 0)
	for _, id := range s.order {
		entry := s.entries[strings.ToLower(id)]
		if entry.InstalledPath != "" {
			p := entry.InstalledPath
			if !filepath.IsAbs(p) {
				p = filepath.Join(workspaceRoot, p)
			}
			cleanP := filepath.Clean(p)
			recordedPaths[strings.ToLower(cleanP)] = true
			if fi, err := os.Stat(cleanP); err == nil && fi.IsDir() {
				recordedDirs = append(recordedDirs, strings.ToLower(cleanP))
			}
		}
	}

	_ = filepath.Walk(workspaceRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == ".orchestra" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		cleanP := filepath.Clean(p)
		lowerP := strings.ToLower(cleanP)

		if recordedPaths[lowerP] {
			return nil
		}
		for _, dir := range recordedDirs {
			if strings.HasPrefix(lowerP, dir+string(filepath.Separator)) {
				return nil
			}
		}

		baseName := strings.ToLower(info.Name())
		switch baseName {
		case "package.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "bun.lockb",
			"go.mod", "go.sum", "tsconfig.json", "readme.md", "license", ".gitignore":
			return nil
		}

		rel, _ := filepath.Rel(workspaceRoot, p)
		report.AllValid = false
		report.FailedCount++
		report.Issues = append(report.Issues, IntegrityIssue{
			ResourceID:  "unlisted",
			Severity:    "ERROR",
			IssueType:   "UNLISTED_FILE",
			Expected:    "Only provenance-recorded files in workspace",
			Actual:      rel,
			Description: fmt.Sprintf("Unlisted file detected in workspace: %s", rel),
		})
		return nil
	})

	return report, nil
}

// computeDirectoryHash generates a deterministic SHA256 hash of all directory files
func computeDirectoryHash(dir string) (string, error) {
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

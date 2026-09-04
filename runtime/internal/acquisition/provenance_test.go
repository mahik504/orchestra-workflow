package acquisition

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestProvenanceStore_RecordAndRetrieve(t *testing.T) {
	tempDir := t.TempDir()

	store, err := NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create provenance store: %v", err)
	}

	entry := ProvenanceEntry{
		ResourceID:          "gsap",
		AcquisitionMethod:   "npm",
		SourceURL:           "https://github.com/greensock/GSAP",
		VersionOrSHA:        "3.12.5",
		SHA256Hash:          "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		InstalledPath:       filepath.Join("node_modules", "gsap"),
		Timestamp:           time.Now().UTC().Format(time.RFC3339),
		JustificationTaskID: "task-premium-landing",
		IsQuarantined:       false,
	}

	if err := store.Record(entry); err != nil {
		t.Fatalf("failed to record entry: %v", err)
	}

	// Verify retrieval by ID
	retrieved, err := store.GetByResourceID("gsap")
	if err != nil {
		t.Fatalf("failed to get entry by ID: %v", err)
	}
	if retrieved.ResourceID != "gsap" || retrieved.VersionOrSHA != "3.12.5" {
		t.Errorf("mismatched retrieved entry: %+v", retrieved)
	}

	// Verify persistence file exists on disk
	ledgerPath := filepath.Join(tempDir, ".orchestra", "provenance.json")
	if _, statErr := os.Stat(ledgerPath); statErr != nil {
		t.Errorf("provenance.json not found on disk: %v", statErr)
	}

	// Reload from new store instance to test disk persistence
	reloadedStore, err := NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatalf("failed to reload store: %v", err)
	}
	allEntries, err := reloadedStore.ListAll()
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}
	if len(allEntries) != 1 || allEntries[0].ResourceID != "gsap" {
		t.Errorf("expected 1 reloaded entry for gsap, got %v", allEntries)
	}
}

func TestProvenanceStore_ValidationErrors(t *testing.T) {
	tempDir := t.TempDir()
	store, _ := NewProvenanceStore(tempDir)

	validEntry := ProvenanceEntry{
		ResourceID:          "lenis",
		AcquisitionMethod:   "npm",
		SourceURL:           "https://github.com/studio-freight/lenis",
		VersionOrSHA:        "1.0.0",
		SHA256Hash:          "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		InstalledPath:       "node_modules/lenis",
		JustificationTaskID: "task-1",
	}

	// Missing ResourceID
	bad := validEntry
	bad.ResourceID = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing ResourceID, got %v", err)
	}

	// Missing AcquisitionMethod
	bad = validEntry
	bad.AcquisitionMethod = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing AcquisitionMethod, got %v", err)
	}

	// Missing SourceURL
	bad = validEntry
	bad.SourceURL = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing SourceURL, got %v", err)
	}

	// Missing VersionOrSHA
	bad = validEntry
	bad.VersionOrSHA = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing VersionOrSHA, got %v", err)
	}

	// Missing SHA256Hash
	bad = validEntry
	bad.SHA256Hash = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing SHA256Hash, got %v", err)
	}

	// Missing JustificationTaskID
	bad = validEntry
	bad.JustificationTaskID = ""
	if err := store.Record(bad); !errors.Is(err, ErrProvenanceEntryInvalid) {
		t.Errorf("expected ErrProvenanceEntryInvalid for missing JustificationTaskID, got %v", err)
	}
}

func TestProvenanceStore_QuarantineAutoDetection(t *testing.T) {
	tempDir := t.TempDir()
	store, _ := NewProvenanceStore(tempDir)

	entry := ProvenanceEntry{
		ResourceID:          "quarantined-item",
		AcquisitionMethod:   "manual",
		SourceURL:           "https://example.com/item",
		VersionOrSHA:        "1.0",
		SHA256Hash:          "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		InstalledPath:       filepath.Join("some", "path", "skills_library", "tool"),
		JustificationTaskID: "task-quarantine-test",
		IsQuarantined:       false,
	}

	if err := store.Record(entry); err != nil {
		t.Fatalf("failed to record entry: %v", err)
	}

	retrieved, _ := store.GetByResourceID("quarantined-item")
	if !retrieved.IsQuarantined {
		t.Errorf("expected IsQuarantined to be automatically set to true for skills_library path")
	}

	quarantinedList, err := store.ListQuarantined()
	if err != nil {
		t.Fatalf("failed to list quarantined: %v", err)
	}
	if len(quarantinedList) != 1 || quarantinedList[0].ResourceID != "quarantined-item" {
		t.Errorf("expected 1 quarantined entry, got %v", quarantinedList)
	}
}

func TestProvenanceStore_VerifyIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	store, _ := NewProvenanceStore(tempDir)

	// Create a real file on disk
	relPath := filepath.Join("assets", "spec.json")
	fullPath := filepath.Join(tempDir, relPath)
	_ = os.MkdirAll(filepath.Dir(fullPath), 0755)

	content := []byte(`{"version": "1.0", "name": "spec"}`)
	_ = os.WriteFile(fullPath, content, 0644)
	h := sha256.Sum256(content)
	correctHash := hex.EncodeToString(h[:])

	entry := ProvenanceEntry{
		ResourceID:          "spec-file",
		AcquisitionMethod:   "web_fetch",
		SourceURL:           "https://example.com/spec.json",
		VersionOrSHA:        "1.0",
		SHA256Hash:          correctHash,
		InstalledPath:       relPath,
		JustificationTaskID: "task-verify",
	}
	_ = store.Record(entry)

	// 1. Valid integrity check
	report, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatalf("integrity check failed: %v", err)
	}
	if !report.AllValid || report.PassedCount != 1 || report.FailedCount != 0 {
		t.Errorf("expected all valid report, got: %+v", report)
	}

	// 2. Tampered file -> Hash Mismatch
	_ = os.WriteFile(fullPath, []byte(`{"tampered": true}`), 0644)
	reportTampered, _ := store.VerifyIntegrity(tempDir)
	if reportTampered.AllValid || reportTampered.FailedCount != 1 {
		t.Errorf("expected hash mismatch failure, got: %+v", reportTampered)
	}
	if len(reportTampered.Issues) == 0 || reportTampered.Issues[0].IssueType != "HASH_MISMATCH" {
		t.Errorf("expected HASH_MISMATCH issue type, got: %+v", reportTampered.Issues)
	}

	// 3. Deleted file -> Missing File
	_ = os.Remove(fullPath)
	reportDeleted, _ := store.VerifyIntegrity(tempDir)
	if reportDeleted.AllValid || reportDeleted.FailedCount != 1 {
		t.Errorf("expected missing file failure, got: %+v", reportDeleted)
	}
	if len(reportDeleted.Issues) == 0 || reportDeleted.Issues[0].IssueType != "MISSING_FILE" {
		t.Errorf("expected MISSING_FILE issue type, got: %+v", reportDeleted.Issues)
	}
}

func TestProvenanceStore_ConcurrentWrites(t *testing.T) {
	tempDir := t.TempDir()
	store, _ := NewProvenanceStore(tempDir)

	var wg sync.WaitGroup
	count := 30

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			entry := ProvenanceEntry{
				ResourceID:          filepath.Base(tempDir) + "-" + hex.EncodeToString([]byte{byte(idx)}),
				AcquisitionMethod:   "npm",
				SourceURL:           "https://example.com/pkg",
				VersionOrSHA:        "1.0.0",
				SHA256Hash:          "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				InstalledPath:       "node_modules/pkg",
				JustificationTaskID: "task-concurrent",
			}
			_ = store.Record(entry)
		}(i)
	}

	wg.Wait()

	all, err := store.ListAll()
	if err != nil {
		t.Fatalf("failed to list all after concurrency: %v", err)
	}
	if len(all) != count {
		t.Errorf("expected %d entries, got %d", count, len(all))
	}
}

func TestVerifyIntegrity_NPMPackageJSONNotDirectoryHash(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewProvenanceStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	pkgDir := filepath.Join(tempDir, "node_modules", "@react-three", "drei")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}
	pkgJSON := []byte(`{"name":"@react-three/drei","version":"10.7.8"}`)
	if err := os.WriteFile(filepath.Join(pkgDir, "package.json"), pkgJSON, 0644); err != nil {
		t.Fatal(err)
	}
	// Extra file so a full directory walk cannot equal the package.json hash.
	if err := os.WriteFile(filepath.Join(pkgDir, "index.js"), []byte("module.exports = {}"), 0644); err != nil {
		t.Fatal(err)
	}

	pkgHash := sha256.Sum256(pkgJSON)
	pkgHex := hex.EncodeToString(pkgHash[:])
	ident := sha256.Sum256([]byte("drei@10.7.8"))
	identHex := hex.EncodeToString(ident[:])
	if pkgHex == identHex {
		t.Fatal("package.json hash unexpectedly equals identity hash")
	}

	rel, _ := filepath.Rel(tempDir, pkgDir)
	entry := ProvenanceEntry{
		ResourceID:          "drei",
		AcquisitionMethod:   "npm",
		SourceURL:           "https://github.com/pmndrs/drei",
		VersionOrSHA:        "10.7.8",
		SHA256Hash:          identHex,
		InstalledPath:       rel,
		JustificationTaskID: "task-drei",
	}
	if err := store.Record(entry); err != nil {
		t.Fatal(err)
	}

	report, err := store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if !report.AllValid || report.FailedCount != 0 {
		t.Fatalf("npm identity hash should verify when package.json exists, got %+v", report)
	}

	entry.SHA256Hash = pkgHex
	if err := store.Record(entry); err != nil {
		t.Fatal(err)
	}
	report, err = store.VerifyIntegrity(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if !report.AllValid || report.FailedCount != 0 {
		t.Fatalf("npm package.json hash should verify, got %+v", report)
	}
}

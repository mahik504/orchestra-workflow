package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/orchestra-v3/internal/acquisition"
	acqAdapters "github.com/user/orchestra-v3/internal/adapters/acquisition"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

func TestImplementStage_Execute_WithAcquisitionAndProvenance(t *testing.T) {
	workdir := t.TempDir()

	// Setup minimal package.json in workdir
	pkgPath := filepath.Join(workdir, "package.json")
	_ = os.WriteFile(pkgPath, []byte(`{"name": "test-project", "version": "1.0.0"}`), 0644)

	mockRunner := runner.NewMockCommandRunner()
	reg := acqAdapters.NewAdapterRegistry()
	reg.RegisterAdapter(acqAdapters.NewNPMAdapter(mockRunner, nil, nil))
	reg.RegisterAdapter(acqAdapters.NewGitAdapter(mockRunner, nil))
	reg.RegisterAdapter(acqAdapters.NewCLIAdapter(mockRunner))
	reg.RegisterAdapter(acqAdapters.NewWebAdapter(true)) // offline fixture mode
	reg.RegisterAdapter(acqAdapters.NewMCPAdapter(mockRunner))

	stage := NewImplementStageWithRegistry(reg)

	taskReq := &TaskRequest{
		ID:            "task-test-implement",
		WorkspaceRoot: workdir,
		Type:          "DESIGN",
	}

	taskCtx := &TaskContext{
		Ctx:  context.Background(),
		Task: taskReq,
		Classification: &ClassificationData{
			Archetype: "premium-website",
			ResolvedRoutes: []resources.CapabilityRoute{
				{
					CapabilityID:   "motion",
					Implementation: []string{"gsap", "playwright"},
					QA:             []string{"playwright"},
				},
			},
		},
		StageResults:  make(map[StageName]*StageResult),
		ArtifactPaths: make(map[string]string),
	}

	res, err := stage.Execute(taskCtx)
	if err != nil {
		t.Fatalf("ImplementStage.Execute failed: %v", err)
	}

	if res.Status != StatusCompleted {
		t.Errorf("expected StatusCompleted, got %s", res.Status)
	}

	// 1. Verify provenance ledger exists
	provPath := filepath.Join(workdir, ".orchestra", "provenance.json")
	if _, statErr := os.Stat(provPath); statErr != nil {
		t.Fatalf("provenance.json does not exist: %v", statErr)
	}

	// 2. Read provenance ledger and verify entries
	provStore, err := acquisition.NewProvenanceStore(workdir)
	if err != nil {
		t.Fatalf("failed to open provenance store: %v", err)
	}

	entries, err := provStore.ListAll()
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 entries (gsap, playwright), got %d", len(entries))
	}

	// Verify task justification ID
	for _, entry := range entries {
		if entry.JustificationTaskID != taskReq.ID {
			t.Errorf("entry %s has wrong JustificationTaskID: expected %s, got %s",
				entry.ResourceID, taskReq.ID, entry.JustificationTaskID)
		}
		if entry.SHA256Hash == "" {
			t.Errorf("entry %s has empty SHA256Hash", entry.ResourceID)
		}
	}

	// 3. Verify modified files in ImplementationData
	if taskCtx.Implementation == nil {
		t.Fatalf("taskCtx.Implementation is nil")
	}

	foundProvInModified := false
	for _, f := range taskCtx.Implementation.ModifiedFiles {
		if stringsEqualNorm(f.Path, provPath) {
			foundProvInModified = true
			if f.SHA256 == "" {
				t.Errorf("expected non-empty checksum for provenance.json in ModifiedFiles")
			}
			break
		}
	}
	if !foundProvInModified {
		t.Errorf("expected provenance.json to be included in ModifiedFiles")
	}
}

func TestImplementStage_DryRun(t *testing.T) {
	workdir := t.TempDir()

	stage := NewImplementStage()

	taskReq := &TaskRequest{
		ID:            "task-dry-run",
		WorkspaceRoot: workdir,
		DryRun:        true,
	}

	taskCtx := &TaskContext{
		Ctx:  context.Background(),
		Task: taskReq,
		Classification: &ClassificationData{
			ResolvedRoutes: []resources.CapabilityRoute{
				{
					Implementation: []string{"gsap"},
				},
			},
		},
		StageResults:  make(map[StageName]*StageResult),
		ArtifactPaths: make(map[string]string),
	}

	res, err := stage.Execute(taskCtx)
	if err != nil {
		t.Fatalf("dry-run execute failed: %v", err)
	}

	if res.Status != StatusCompleted {
		t.Errorf("expected StatusCompleted for dry run, got %s", res.Status)
	}

	provStore, err := acquisition.NewProvenanceStore(workdir)
	if err != nil {
		t.Fatalf("failed opening provenance store: %v", err)
	}

	entry, err := provStore.GetByResourceID("gsap")
	if err != nil {
		t.Fatalf("expected dry-run entry for gsap: %v", err)
	}
	if entry.VersionOrSHA != "dry-run-v1" {
		t.Errorf("expected dry-run-v1 version, got %s", entry.VersionOrSHA)
	}
}

func stringsEqualNorm(p1, p2 string) bool {
	return filepath.Clean(p1) == filepath.Clean(p2)
}

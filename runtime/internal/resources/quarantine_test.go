package resources

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckQuarantineBoundary(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "clean workspace path",
			path:        "C:/dev/orchestra-workflow/registries/resources.json",
			expectError: false,
		},
		{
			name:        "clean private workspace path",
			path:        "C:\\dev\\my-workspace\\registries\\design-resource-graph.json",
			expectError: false,
		},
		{
			name:        "quarantined skills_library path",
			path:        "C:/Users/mockuser/.gemini/config/skills_library/some-skill/SKILL.md",
			expectError: true,
		},
		{
			name:        "quarantined skills_library backslash path",
			path:        "C:\\Users\\mockuser\\.gemini\\config\\skills_library\\another-skill",
			expectError: true,
		},
		{
			name:        "case-insensitive quarantined path",
			path:        "C:/Users/mockuser/.gemini/config/SKILLS_LIBRARY/skill/SKILL.md",
			expectError: true,
		},
		{
			name:        "hyphenated skills-library path",
			path:        "/home/user/.gemini/config/skills-library/tool",
			expectError: true,
		},
		{
			name:        "curated_catalog quarantine path",
			path:        "C:/projects/orchestra/curated_catalog/quarantine/bad-skill",
			expectError: true,
		},
		{
			name:        "8.3 short name SKILLS~1 path",
			path:        "C:/Users/mockuser/.gemini/config/SKILLS~1/skill/SKILL.md",
			expectError: true,
		},
		{
			name:        "8.3 short name SKILLS~2 path",
			path:        "C:/Users/mockuser/.gemini/config/SKILLS~2/archive-skill",
			expectError: true,
		},
		{
			name:        "8.3 short name skills~3 path",
			path:        "/var/data/skills~3/bad",
			expectError: true,
		},
		{
			name:        "8.3 short name CURATE~1 path",
			path:        "C:/projects/CURATE~1/quarantine/bad-skill",
			expectError: true,
		},
		{
			name:        "8.3 short name CURATE~2 path",
			path:        "C:/projects/CURATE~2/quarantine/bad-skill",
			expectError: true,
		},
		{
			name:        "URL percent-encoded skills_library (%73kills_library)",
			path:        "file:///C:/Users/mockuser/.gemini/config/%73kills_library/skill/SKILL.md",
			expectError: true,
		},
		{
			name:        "URL percent-encoded underscore (skills%5flibrary)",
			path:        "C:/Users/mockuser/.gemini/config/skills%5flibrary/skill/SKILL.md",
			expectError: true,
		},
		{
			name:        "URL percent-encoded 8.3 (SKILLS%7e1)",
			path:        "file:///C:/Users/mockuser/.gemini/config/SKILLS%7e1/skill",
			expectError: true,
		},
		{
			name:        "quarantined skills_archive path",
			path:        "C:/Users/mockuser/.gemini/config/skills_archive_6-02_044222/old-skill",
			expectError: true,
		},
		{
			name:        "hyphenated skills-archive path",
			path:        "/home/user/.gemini/config/skills-archive/old-skill",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckQuarantineBoundary(tc.path)
			if tc.expectError && err == nil {
				t.Errorf("expected error for path '%s', got nil", tc.path)
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error for path '%s': %v", tc.path, err)
			}
			if tc.expectError && err != nil && !errors.Is(err, ErrQuarantinedPath) {
				t.Errorf("expected ErrQuarantinedPath, got: %v", err)
			}
		})
	}
}

func TestQuarantineBoundaryRejection(t *testing.T) {
	tempDir := t.TempDir()
	quarantineSubdir := filepath.Join(tempDir, "skills_library")
	if err := os.MkdirAll(quarantineSubdir, 0755); err != nil {
		t.Fatalf("failed to create quarantine test dir: %v", err)
	}

	fakeCatalog := filepath.Join(quarantineSubdir, "resources.json")
	if err := os.WriteFile(fakeCatalog, []byte("[]"), 0644); err != nil {
		t.Fatalf("failed to write fake catalog: %v", err)
	}

	fakeGraph := filepath.Join(quarantineSubdir, "graph.json")
	if err := os.WriteFile(fakeGraph, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write fake graph: %v", err)
	}

	// 1. LoadResourceCatalog must reject quarantined path
	_, err := LoadResourceCatalog(fakeCatalog)
	if err == nil {
		t.Errorf("expected LoadResourceCatalog to reject quarantined path, but it succeeded")
	} else if !errors.Is(err, ErrQuarantinedPath) {
		t.Errorf("expected ErrQuarantinedPath, got: %v", err)
	}

	// 2. LoadDesignGraph must reject quarantined path
	_, err = LoadDesignGraph(fakeGraph)
	if err == nil {
		t.Errorf("expected LoadDesignGraph to reject quarantined path, but it succeeded")
	} else if !errors.Is(err, ErrQuarantinedPath) {
		t.Errorf("expected ErrQuarantinedPath, got: %v", err)
	}

	// 3. Registry.LoadFromJSON must reject quarantined path
	reg := NewRegistry()
	err = reg.LoadFromJSON(fakeCatalog)
	if err == nil {
		t.Errorf("expected Registry.LoadFromJSON to reject quarantined path, but it succeeded")
	} else if !errors.Is(err, ErrQuarantinedPath) {
		t.Errorf("expected ErrQuarantinedPath, got: %v", err)
	}
}

func TestQuarantineStateAudit(t *testing.T) {
	tempDir := t.TempDir()
	workspaceRoot := filepath.Join(tempDir, "workspace")
	quarantineRoot := filepath.Join(tempDir, "external", "skills_library")

	if err := os.MkdirAll(workspaceRoot, 0755); err != nil {
		t.Fatalf("failed to create workspace dir: %v", err)
	}
	if err := os.MkdirAll(quarantineRoot, 0755); err != nil {
		t.Fatalf("failed to create quarantine dir: %v", err)
	}

	// Create 3 dummy skill dirs
	for _, name := range []string{"skill1", "skill2", "skill3"} {
		_ = os.MkdirAll(filepath.Join(quarantineRoot, name), 0755)
	}

	status, err := AuditQuarantineState(workspaceRoot, quarantineRoot)
	if err != nil {
		t.Fatalf("AuditQuarantineState failed: %v", err)
	}

	if !status.QuarantineDirectoryExists {
		t.Errorf("expected quarantine directory to exist")
	}
	if status.QuarantinedCount != 3 {
		t.Errorf("expected quarantined count 3, got %d", status.QuarantinedCount)
	}
	if !status.IsStrictlyIsolated {
		t.Errorf("expected strict isolation for external quarantine directory")
	}

	// Now test internal quarantine directory (inside workspace root) which should trigger a violation
	internalQuarantine := filepath.Join(workspaceRoot, "skills_library")
	_ = os.MkdirAll(internalQuarantine, 0755)
	badStatus, err := AuditQuarantineState(workspaceRoot, internalQuarantine)
	if err != nil {
		t.Fatalf("AuditQuarantineState failed: %v", err)
	}
	if badStatus.IsStrictlyIsolated {
		t.Errorf("expected isolation violation when quarantine resides inside workspace")
	}
	if len(badStatus.ViolationsFound) == 0 {
		t.Errorf("expected violation messages to be recorded")
	}
}

func TestLiveRegistriesQuarantineIsolation(t *testing.T) {
	workflowRoot := filepath.Join("..", "..", "..")
	resourcesPath := filepath.Join(workflowRoot, "registries", "resources.json")
	graphPath := filepath.Join(workflowRoot, "registries", "design-resource-graph.json")

	// 1. Raw file content must have zero mentions of skills_library
	resBytes, err := os.ReadFile(resourcesPath)
	if err != nil {
		t.Fatalf("failed to read resources.json: %v", err)
	}
	if strings.Contains(strings.ToLower(string(resBytes)), "skills_library") {
		t.Fatalf("SECURITY VIOLATION: resources.json contains reference to quarantined skills_library!")
	}

	graphBytes, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read design-resource-graph.json: %v", err)
	}
	if strings.Contains(strings.ToLower(string(graphBytes)), "skills_library") {
		t.Fatalf("SECURITY VIOLATION: design-resource-graph.json contains reference to quarantined skills_library!")
	}

	// 2. Parsed catalog check
	cat, err := LoadResourceCatalog(resourcesPath)
	if err != nil {
		t.Fatalf("failed to load catalog: %v", err)
	}
	violations := VerifyCatalogQuarantineClean(cat)
	if len(violations) > 0 {
		t.Fatalf("SECURITY VIOLATION: catalog resources contain quarantined URLs: %v", violations)
	}
}

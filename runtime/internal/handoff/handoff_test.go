package handoff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandoffLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestra-handoff-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample file tracked by Agent A
	sampleFile := filepath.Join(tempDir, "feature.go")
	if err := os.WriteFile(sampleFile, []byte("package main\n\nfunc Run() {}\n"), 0644); err != nil {
		t.Fatalf("Failed to write sample file: %v", err)
	}

	hashA, err := ComputeFileHash(sampleFile)
	if err != nil {
		t.Fatalf("Failed to hash sample file: %v", err)
	}

	// Agent A (e.g. Antigravity) writes Handoff State v1
	stateA := &HandoffState{
		SessionID:   "sess-12345",
		Version:     1,
		SourceAgent: "antigravity",
		TargetAgent: "cursor",
		ActiveTasks: []string{"task-001"},
		PlanURI:     "docs/plan.md",
		ChangedFiles: []FileChecksum{
			{Path: "feature.go", SHA256: hashA},
		},
		CompletedSteps: []string{"step-1-scaffold"},
		PendingSteps:   []string{"step-2-implement", "step-3-verify"},
	}

	if err := WriteState(stateA, tempDir); err != nil {
		t.Fatalf("Agent A failed to write state: %v", err)
	}

	// Agent B (e.g. Cursor) reads Handoff State
	stateB, err := ReadState(tempDir)
	if err != nil {
		t.Fatalf("Agent B failed to read state: %v", err)
	}

	if stateB.Version != 1 || stateB.TargetAgent != "cursor" {
		t.Errorf("Agent B read incorrect state: %+v", stateB)
	}

	// Check conflicts before work
	conflicts, err := DetectConflicts(stateB, tempDir)
	if err != nil {
		t.Fatalf("Failed to detect conflicts: %v", err)
	}
	if len(conflicts) != 0 {
		t.Errorf("Expected 0 conflicts, got %v", conflicts)
	}

	// Agent B completes step 2 and increments version to v2
	stateB.Version = 2
	stateB.SourceAgent = "cursor"
	stateB.TargetAgent = "antigravity"
	stateB.CompletedSteps = append(stateB.CompletedSteps, "step-2-implement")
	stateB.PendingSteps = []string{"step-3-verify"}

	if err := WriteState(stateB, tempDir); err != nil {
		t.Fatalf("Agent B failed to update state to v2: %v", err)
	}

	// Read state v2
	finalState, err := ReadState(tempDir)
	if err != nil {
		t.Fatalf("Failed to read v2 state: %v", err)
	}
	if finalState.Version != 2 || len(finalState.CompletedSteps) != 2 {
		t.Errorf("Unexpected v2 state: %+v", finalState)
	}
}

func TestConflictDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestra-conflict-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sampleFile := filepath.Join(tempDir, "shared.go")
	if err := os.WriteFile(sampleFile, []byte("original content"), 0644); err != nil {
		t.Fatalf("Failed to write sample file: %v", err)
	}

	origHash, _ := ComputeFileHash(sampleFile)

	state := &HandoffState{
		SessionID: "sess-conflict",
		Version:   1,
		ChangedFiles: []FileChecksum{
			{Path: "shared.go", SHA256: origHash},
		},
	}

	// Modify file externally behind agent's back
	if err := os.WriteFile(sampleFile, []byte("clobbered content from other tool"), 0644); err != nil {
		t.Fatalf("Failed to clobber file: %v", err)
	}

	conflicts, err := DetectConflicts(state, tempDir)
	if err != nil {
		t.Fatalf("Failed conflict detection: %v", err)
	}

	if len(conflicts) == 0 {
		t.Errorf("Expected conflict to be detected, but got 0 conflicts")
	}
}

func TestFailureRecoveryAndResume(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestra-recovery-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	state := &HandoffState{
		SessionID:      "sess-recover",
		Version:        1,
		CompletedSteps: []string{"step-1"},
		PendingSteps:   []string{"step-2", "step-3"},
		FailureRecovery: &RecoveryPoint{
			FailedStep:     "step-2",
			ErrorReason:    "compiler error in template",
			CanResume:      true,
			ResumeFromStep: "step-2",
		},
	}

	if err := WriteState(state, tempDir); err != nil {
		t.Fatalf("Failed to write recovery state: %v", err)
	}

	loaded, err := ReadState(tempDir)
	if err != nil {
		t.Fatalf("Failed to read state: %v", err)
	}

	if loaded.FailureRecovery == nil || !loaded.FailureRecovery.CanResume || loaded.FailureRecovery.ResumeFromStep != "step-2" {
		t.Errorf("Failure recovery metadata did not persist correctly: %+v", loaded.FailureRecovery)
	}
}

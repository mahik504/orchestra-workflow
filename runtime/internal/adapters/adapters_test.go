package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/resources"
)

func TestHostAdapters_Contract(t *testing.T) {
	adapters := []HostAdapter{
		&CursorAdapter{},
		&ClaudeAdapter{},
		&AntigravityAdapter{},
	}

	expectedNames := map[string]bool{
		"cursor":      true,
		"claude":      true,
		"antigravity": true,
	}

	tmpDir := t.TempDir()

	for _, a := range adapters {
		name := a.Name()
		if !expectedNames[name] {
			t.Errorf("Unexpected adapter name: %s", name)
		}

		envs := a.SupportedEnvironments()
		if len(envs) == 0 {
			t.Errorf("Adapter %s returned empty environments", name)
		}

		skillsPath := a.GetActiveSkillsPath(tmpDir)
		if skillsPath == "" {
			t.Errorf("Adapter %s returned empty skills path", name)
		}

		// Verify ExecutePlan
		err := a.ExecutePlan(&handoff.HandoffState{SessionID: "test-sess"})
		if err != nil {
			t.Errorf("Adapter %s ExecutePlan returned error: %v", name, err)
		}

		// Verify GenerateConfig
		err = a.GenerateConfig(tmpDir)
		if err != nil {
			t.Errorf("Adapter %s GenerateConfig returned error: %v", name, err)
		}
	}
}

func TestHostSyncEngine_VerifyParity_IsolatedTemp(t *testing.T) {
	fakeHome := t.TempDir()
	engine := NewHostSyncEngine(fakeHome)

	// Set up fake cursor, claude, antigravity with 3 canonical skills
	testSkills := CanonicalActiveSkills[:3]
	hostPaths := []string{
		filepath.Join(fakeHome, ".cursor", "skills"),
		filepath.Join(fakeHome, ".claude", "skills"),
		filepath.Join(fakeHome, ".gemini", "config", "skills"),
	}

	for _, hp := range hostPaths {
		for _, s := range testSkills {
			skillDir := filepath.Join(hp, s)
			_ = os.MkdirAll(skillDir, 0755)
			_ = os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Skill "+s), 0644)
		}
	}

	report, err := engine.VerifyParity(fakeHome)
	if err != nil {
		t.Fatalf("VerifyParity failed: %v", err)
	}

	if report.CanonicalSkillCount != 30 {
		t.Errorf("Expected 30 canonical skills, got %d", report.CanonicalSkillCount)
	}

	// 27 missing from the 30
	if report.IsParityComplete {
		t.Errorf("Expected parity to be false because only 3 skills are installed")
	}

	for _, h := range []HostType{HostCursor, HostClaudeCode, HostAntigravity} {
		if report.HostSkillCounts[h] != 3 {
			t.Errorf("Expected host %s to have 3 skills, got %d", h, report.HostSkillCounts[h])
		}
	}

	if !report.ByteIdentical {
		t.Errorf("Expected byte identity across the 3 matching skills")
	}
}

func TestHostSyncEngine_QuarantineRejection(t *testing.T) {
	fakeHome := t.TempDir()
	quarantineHome := filepath.Join(fakeHome, "skills_library")
	engine := NewHostSyncEngine(".")

	// SyncAll with quarantined path must fail
	err := engine.SyncAll(quarantineHome, "all")
	if err == nil {
		t.Fatalf("Expected quarantine error when syncing to skills_library, got nil")
	}
	if !errors.Is(err, resources.ErrQuarantinedPath) {
		t.Errorf("Expected error wrapping ErrQuarantinedPath, got: %v", err)
	}
}

func TestComputeSkillDirHash(t *testing.T) {
	tmp1 := t.TempDir()
	tmp2 := t.TempDir()

	skill1 := filepath.Join(tmp1, "test-skill")
	skill2 := filepath.Join(tmp2, "test-skill")

	_ = os.MkdirAll(skill1, 0755)
	_ = os.MkdirAll(skill2, 0755)

	_ = os.WriteFile(filepath.Join(skill1, "SKILL.md"), []byte("Hello World"), 0644)
	_ = os.WriteFile(filepath.Join(skill2, "SKILL.md"), []byte("Hello World"), 0644)

	h1, err := ComputeSkillDirHash(skill1)
	if err != nil {
		t.Fatalf("ComputeSkillDirHash 1 failed: %v", err)
	}
	h2, err := ComputeSkillDirHash(skill2)
	if err != nil {
		t.Fatalf("ComputeSkillDirHash 2 failed: %v", err)
	}

	if h1 != h2 {
		t.Errorf("Hashes must match for identical files: %s != %s", h1, h2)
	}

	// Mutate skill2
	_ = os.WriteFile(filepath.Join(skill2, "SKILL.md"), []byte("Hello World Modified"), 0644)
	h3, err := ComputeSkillDirHash(skill2)
	if err != nil {
		t.Fatalf("ComputeSkillDirHash 3 failed: %v", err)
	}

	if h1 == h3 {
		t.Errorf("Hashes must not match for modified files")
	}
}

func TestLiveMachineHostParity(t *testing.T) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot resolve UserHomeDir, skipping live check")
	}

	engine := NewHostSyncEngine(".")
	report, err := engine.VerifyParity(userHome)
	if err != nil {
		t.Fatalf("Live VerifyParity failed: %v", err)
	}

	// If Cursor, Claude, and Gemini skills exist on this machine, verify 30-skill parity
	if report.HostSkillCounts[HostCursor] == 30 && report.HostSkillCounts[HostClaudeCode] == 30 && report.HostSkillCounts[HostAntigravity] == 30 {
		if !report.ByteIdentical {
			t.Errorf("Expected 100%% byte-identical skills across live hosts, mismatches: %v", report.MismatchDetails)
		}
		if len(report.QuarantineViolations) > 0 {
			t.Errorf("Quarantine violations detected: %v", report.QuarantineViolations)
		}
	}
}

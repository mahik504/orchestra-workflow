package router

import (
	"testing"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

func setupTestRegistry() *resources.Registry {
	reg := resources.NewRegistry()
	reg.Capabilities["superpowers-planning"] = &resources.Capability{
		ID:                 "superpowers-planning",
		Name:               "Superpowers Planning",
		TokenContextWeight: 1500,
	}
	reg.Capabilities["taste-skill"] = &resources.Capability{
		ID:                 "taste-skill",
		Name:               "Taste Skill",
		TokenContextWeight: 2000,
	}
	reg.Capabilities["impeccable"] = &resources.Capability{
		ID:                 "impeccable",
		Name:               "Impeccable Design",
		TokenContextWeight: 1800,
	}
	reg.Capabilities["semgrep-adapter"] = &resources.Capability{
		ID:                 "semgrep-adapter",
		Name:               "Semgrep Security Adapter",
		TokenContextWeight: 1200,
	}
	return reg
}

func TestVisualTaskRouting(t *testing.T) {
	reg := setupTestRegistry()
	r := NewRouter(reg)

	task := &classifier.Task{
		ID:             "task-visual",
		Type:           "DESIGN",
		RequiresVisual: true,
	}

	plan := r.Compose(task)

	if !plan.RequiresHumanGate {
		t.Errorf("Expected visual task to require human gate, got false")
	}

	hasTaste := false
	hasImpeccable := false
	for _, cap := range plan.SelectedCapabilities {
		if cap.ID == "taste-skill" {
			hasTaste = true
		}
		if cap.ID == "impeccable" {
			hasImpeccable = true
		}
	}

	if !hasTaste || !hasImpeccable {
		t.Errorf("Expected taste-skill and impeccable to be selected for visual task")
	}

	if len(plan.ExecutionDirectives) < 3 {
		t.Errorf("Expected at least 3 execution directives, got %d", len(plan.ExecutionDirectives))
	}
}

func TestBackendBugRouting(t *testing.T) {
	reg := setupTestRegistry()
	r := NewRouter(reg)

	task := &classifier.Task{
		ID:               "task-backend-bug",
		Type:             "BUGFIX",
		RequiresVisual:   false,
		RequiresSecurity: false,
	}

	plan := r.Compose(task)

	if plan.RequiresHumanGate {
		t.Errorf("Expected backend bug task NOT to require human visual gate")
	}

	for _, cap := range plan.SelectedCapabilities {
		if cap.ID == "taste-skill" || cap.ID == "impeccable" {
			t.Errorf("Backend bug task should NOT activate visual capabilities, found %s", cap.ID)
		}
	}
}

func TestMixedApplicationRouting(t *testing.T) {
	reg := setupTestRegistry()
	r := NewRouter(reg)

	task := &classifier.Task{
		ID:               "task-mixed",
		Type:             "FEATURE",
		RequiresVisual:   true,
		RequiresSecurity: true,
	}

	plan := r.Compose(task)

	hasVisual := false
	hasSecurity := false
	for _, cap := range plan.SelectedCapabilities {
		if cap.ID == "taste-skill" {
			hasVisual = true
		}
		if cap.ID == "semgrep-adapter" {
			hasSecurity = true
		}
	}

	if !hasVisual || !hasSecurity {
		t.Errorf("Expected mixed application to compose both visual and security capabilities")
	}
}

func TestUnknownTechnologyRouting(t *testing.T) {
	reg := setupTestRegistry()
	r := NewRouter(reg)

	task := &classifier.Task{
		ID:                 "task-unknown-tech",
		Type:               "FEATURE",
		SuggestedResources: []string{"unindexed-wasm-runtime", "quantum-crypto-sdk"},
	}

	plan := r.Compose(task)

	if len(plan.GapResearchNeeded) != 2 {
		t.Errorf("Expected 2 gap technologies detected, got %d", len(plan.GapResearchNeeded))
	}

	hasGapDirective := false
	for _, d := range plan.ExecutionDirectives {
		if d.CapabilityID == "capability-gap-research" {
			hasGapDirective = true
		}
	}

	if !hasGapDirective {
		t.Errorf("Expected capability-gap-research directive to be generated")
	}
}

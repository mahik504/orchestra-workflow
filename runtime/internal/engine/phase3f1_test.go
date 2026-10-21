package engine

import (
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/memory"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

func TestPhase3F1_RequiredUnprobedUsesHumanGate(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "gate", Objective: "Fix parser crash", DeclaredType: "BUGFIX",
		Acceptance: []string{"parser returns the node"},
	})
	acc, bar, retry := strings.Join(m.Acceptance, "|"), m.QualityBar, m.Budget.RetryBudget
	ctx := &TaskContext{
		Mission: m,
		Task: &TaskRequest{
			ID: m.ID, WorkspaceRoot: t.TempDir(), Type: "BUGFIX", RawRequest: m.Objective,
			Parameters: map[string]any{
				"playwright_available": false,
				"capability_mcp": []capability.MCPCandidate{{
					ID: "needed-tool", Health: "OPTIONAL", Required: true,
				}},
			},
		},
	}
	err := selectCapabilities(ctx)
	if err != ErrHumanGateRequired {
		t.Fatalf("got %v", err)
	}
	if strings.Join(m.Acceptance, "|") != acc || m.QualityBar != bar || m.Budget.RetryBudget != retry {
		t.Fatal("human gate mutated mission authority")
	}
	if ctx.Capability == nil || !strings.Contains(ctx.Capability.HumanDetail, "needed-tool") {
		t.Fatalf("decision %+v", ctx.Capability)
	}
}

func TestPhase3F1_BrainMissStillCompiles(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "miss", Objective: "Fix parser crash", DeclaredType: "BUGFIX",
		Acceptance: []string{"parser returns the node"}, QualityBar: mission.BarPremium,
	})
	acc, bar, retry := strings.Join(m.Acceptance, "|"), m.QualityBar, m.Budget.RetryBudget
	for _, mode := range []string{"unavailable", "empty", "malformed", "missing"} {
		raw := []byte(nil)
		if mode == "malformed" {
			raw = []byte("{")
		}
		if mode == "empty" {
			raw = []byte("[]")
		}
		got := memory.LoadDurable(raw, mode, "job", nil)
		if len(got.Selected) != 0 {
			t.Fatalf("%s enriched %d", mode, len(got.Selected))
		}
		ctx := &TaskContext{
			Mission: m,
			Task: &TaskRequest{
				ID: m.ID, WorkspaceRoot: t.TempDir(), Type: "BUGFIX", RawRequest: m.Objective,
				Parameters: map[string]any{"playwright_available": false},
			},
		}
		if err := selectCapabilities(ctx); err != nil {
			t.Fatalf("%s select: %v", mode, err)
		}
		compileBuilderCapsule(ctx)
		if ctx.Capsule == nil || ctx.Capsule.Role != mission.RoleBuilder {
			t.Fatalf("%s capsule missing", mode)
		}
		if !ctx.Capability.Has(capability.IDCodingWorker) {
			t.Fatalf("%s dropped coding.worker", mode)
		}
		if strings.Join(m.Acceptance, "|") != acc || m.QualityBar != bar || m.Budget.RetryBudget != retry {
			t.Fatalf("%s mutated mission", mode)
		}
	}
}

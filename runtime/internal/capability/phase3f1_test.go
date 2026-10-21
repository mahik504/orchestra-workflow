package capability

import (
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

func TestPhase3F1_HealthDoesNotDefineRequirement(t *testing.T) {
	cases := []struct {
		name     string
		health   string
		required bool
		human    bool
		selected bool
	}{
		{"optional unprobed", "OPTIONAL", false, false, false},
		{"required unprobed", "OPTIONAL", true, true, false},
		{"required blank", "", true, true, false},
		{"required auth", "AUTH_REQUIRED", true, true, false},
		{"required broken", "BROKEN", true, true, false},
		{"required disabled", "DISABLED", true, true, false},
		{"optional broken", "BROKEN", false, false, false},
		{"required healthy", "HEALTHY", true, false, true},
		{"optional healthy", "HEALTHY", false, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := mission.FromSeed(mission.Seed{ID: "mcp", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
			pw := false
			before, _ := missionSnapshot(m)
			d, err := Select(Input{
				Mission: m, PlaywrightAvailable: &pw,
				MCP: []MCPCandidate{{ID: "needed-tool", Health: tc.health, Required: tc.required}},
			})
			after, _ := missionSnapshot(m)
			if before != after {
				t.Fatal("mcp health changed the mission")
			}
			if tc.human {
				if !IsHumanStop(err) {
					t.Fatalf("want human stop, got %v", err)
				}
				if !strings.Contains(d.HumanDetail, "needed-tool") {
					t.Fatalf("detail %s", d.HumanDetail)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			got := false
			for _, id := range d.Tools {
				if id == "needed-tool" {
					got = true
				}
			}
			if got != tc.selected {
				t.Fatalf("selected=%v want %v tools=%v", got, tc.selected, d.Tools)
			}
		})
	}
}

func missionSnapshot(m *mission.Mission) (string, int) {
	return m.QualityBar + "|" + strings.Join(m.Acceptance, "|"), m.Budget.RetryBudget
}

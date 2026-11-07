package benchmark

import "testing"

func TestComparisonSuiteMakesNoClaim(t *testing.T) {
	suite := ComparisonSuite()
	if len(suite) != 6 {
		t.Fatalf("tasks = %d", len(suite))
	}
	want := map[string]bool{
		"requirement_completion":    false,
		"visual_quality":            false,
		"interaction_correctness":   false,
		"code_quality":              false,
		"architectural_correctness": false,
		"performance":               false,
		"accessibility":             false,
		"regression_count":          false,
		"time":                      false,
		"tool_cost":                 false,
		"token_cost":                false,
	}
	enabledPilot := 0
	for _, task := range EnabledTasks() {
		enabledPilot++
		if task.ID == "c1-premium-landing" {
			t.Fatal("comparison task enabled in the pilot")
		}
	}
	if enabledPilot != 2 {
		t.Fatalf("pilot enabled tasks = %d, want 2", enabledPilot)
	}
	for _, task := range suite {
		if task.Enabled {
			t.Fatalf("%s is enabled without a paired run", task.ID)
		}
		if task.Claim != ClaimNone {
			t.Fatalf("%s claim %s", task.ID, task.Claim)
		}
		for _, measure := range task.Measures {
			if _, ok := want[measure]; ok {
				want[measure] = true
			}
		}
	}
	for measure, seen := range want {
		if !seen {
			t.Fatalf("missing measure %s", measure)
		}
	}
}

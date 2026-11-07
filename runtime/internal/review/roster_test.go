package review

import "testing"

func TestTypoDoesNotSummonThePanel(t *testing.T) {
	got := Select(Input{Task: "typo", DiffLines: 2})
	if len(got) != 1 || got[0].Role != "builder" {
		t.Fatalf("%+v", got)
	}
}

func TestFlagshipUsesIndependentJudges(t *testing.T) {
	got := Select(Input{Task: "flagship", Capability: "premium-website", DiffLines: 200})
	lenses := map[string]bool{}
	for _, item := range got {
		lenses[item.Lens] = true
		if item.Role == "judge" {
			if err := EvidenceOK(item, nil, "looks done"); err == nil {
				t.Fatalf("empty evidence accepted for %s", item.Lens)
			}
		}
	}
	for _, lens := range []string{"requirements", "implement", "visual", "performance", "simplify"} {
		if !lenses[lens] {
			t.Fatalf("missing %s in %+v", lens, got)
		}
	}
}

func TestVisualJudgeNeedsAScreenshot(t *testing.T) {
	judge := Assignment{Role: "judge", Lens: "visual", Evidence: "screenshots"}
	if err := EvidenceOK(judge, []string{"builder says it looks good"}, "builder says it looks good"); err == nil {
		t.Fatal("accepted a builder impression")
	}
	if err := EvidenceOK(judge, []string{".orchestra/runs/example/desktop.png"}, "builder says it looks good"); err != nil {
		t.Fatal(err)
	}
}

func TestBackendSkipsTheVisualJudge(t *testing.T) {
	got := Select(Input{Task: "backend", Capability: "b2b-portal", TouchesAuth: true, DiffLines: 90})
	for _, item := range got {
		if item.Lens == "visual" {
			t.Fatal("visual judge on a backend task")
		}
	}
	if !hasLens(got, "architecture") || !hasLens(got, "security") || !hasLens(got, "simplify") {
		t.Fatalf("%+v", got)
	}
}

func hasLens(items []Assignment, lens string) bool {
	for _, item := range items {
		if item.Lens == lens {
			return true
		}
	}
	return false
}

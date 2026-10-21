package mission

import "testing"

func TestInferVisualDoesNotUseQualityBar(t *testing.T) {
	flag := FromSeed(Seed{ID: "f", Objective: "Award-winning 3D portfolio showcase", DeclaredType: "DESIGN", QualityBar: BarStandard})
	if flag.VisualImportance != VisualFlagship {
		t.Fatalf("language must mark flagship even on STANDARD, got %s", flag.VisualImportance)
	}
	crud := FromSeed(Seed{ID: "c", Objective: "Add user CRUD endpoints", DeclaredType: "FEATURE", QualityBar: BarExperimental})
	if crud.VisualImportance == VisualFlagship {
		t.Fatalf("EXPERIMENTAL CRUD must not become flagship via the bar, type=%s vis=%s", crud.Type, crud.VisualImportance)
	}
	settings := FromSeed(Seed{ID: "s", Objective: "Build a settings screen for profile CRUD", DeclaredType: "FEATURE", QualityBar: BarPremium})
	if settings.Type != TypeProductUI || settings.VisualImportance != VisualProduct {
		t.Fatalf("settings UI: type=%s vis=%s", settings.Type, settings.VisualImportance)
	}
	if settings.Budget.VisualReview != VisualReviewBrowser {
		t.Fatalf("review %s", settings.Budget.VisualReview)
	}
	if settings.Budget.IndependentJudge == JudgeRequired {
		t.Fatal("settings UI must not require a judge")
	}
}

func TestCompileJudgeStripsWriteTools(t *testing.T) {
	m := validMission()
	cap, _, err := Compile(CompileRequest{
		Mission:      m,
		Role:         RoleJudge,
		ToolsAllowed: []string{"read_file", "write_file", "apply_patch", "shell"},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range cap.ToolsAllowed {
		if tool == "write_file" || tool == "shell" || tool == "apply_patch" {
			t.Fatalf("judge tools: %v", cap.ToolsAllowed)
		}
	}
}

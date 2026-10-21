package visual

import "testing"

func TestDecideJudgeOrdinaryOff(t *testing.T) {
	d := DecideJudge(JudgeInput{Class: ClassOrdinary, Policy: "conditional"})
	if d.Activated || d.Mode != "off" {
		t.Fatalf("%+v", d)
	}
}

func TestDecideJudgeProductConditionalStaysOff(t *testing.T) {
	d := DecideJudge(JudgeInput{Class: ClassProductUI, Policy: "off"})
	if d.Activated {
		t.Fatalf("ordinary CRUD UI must not pay for a judge: %+v", d)
	}
}

func TestDecideJudgeHighRiskRequired(t *testing.T) {
	d := DecideJudge(JudgeInput{Class: ClassFlagship, SecurityRelease: true})
	if !d.Activated || d.Mode != "required" {
		t.Fatalf("%+v", d)
	}
	d = DecideJudge(JudgeInput{Class: ClassFlagship, Policy: "conditional", PriorVisualFail: true})
	if !d.Activated {
		t.Fatalf("prior visual fail should activate conditional judge: %+v", d)
	}
	d = DecideJudge(JudgeInput{Class: ClassFlagship, Policy: "conditional"})
	if d.Activated {
		t.Fatalf("clear flagship with no trigger stays off: %+v", d)
	}
}

func TestJudgeCapsuleStripsWriteTools(t *testing.T) {
	got := JudgeCapsuleTools([]string{"read_file", "write_file", "apply_patch", "shell"})
	for _, tname := range got {
		if tname == "write_file" || tname == "apply_patch" || tname == "shell" {
			t.Fatalf("write tool leaked: %v", got)
		}
	}
}

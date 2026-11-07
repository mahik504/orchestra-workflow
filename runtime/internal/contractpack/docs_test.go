package contractpack

import "testing"

func TestBugCreatesNoDesignPack(t *testing.T) {
	docs := []string{"DESIGN.md", "3D_SCENE.md", "API.md"}
	if Required("bug", docs) != nil {
		t.Fatal("a bug opened the design pack")
	}
	if Required("control", docs) != nil {
		t.Fatal("one control opened the design pack")
	}
}

func TestPremiumKeepsTheDomainPack(t *testing.T) {
	got := Required("flagship", []string{"DESIGN.md", "SPEC.md", "MOTION.md"})
	if len(got) != 3 || got[0] != "DESIGN.md" {
		t.Fatalf("%v", got)
	}
}

package mission

import "testing"

func TestOwnedFilesConflictAcrossMissions(t *testing.T) {
	left := &Capsule{MissionID: "a", OwnedFiles: []string{"src/app.tsx", "README.md"}}
	right := &Capsule{MissionID: "b", OwnedFiles: []string{"src/app.tsx"}}
	same := &Capsule{MissionID: "a", OwnedFiles: []string{"README.md"}}
	if got := OwnershipConflict([]*Capsule{left, same}); len(got) != 0 {
		t.Fatalf("same mission conflicted: %v", got)
	}
	got := OwnershipConflict([]*Capsule{left, right})
	if len(got) != 1 || got[0] == "" {
		t.Fatalf("conflicts %v", got)
	}
	windowsPath := &Capsule{MissionID: "b", OwnedFiles: []string{`src\app.tsx`}}
	slash := OwnershipConflict([]*Capsule{left, windowsPath})
	if len(slash) != 1 {
		t.Fatalf("slash conflicts %v", slash)
	}
	otherCase := &Capsule{MissionID: "c", OwnedFiles: []string{"src/App.tsx"}}
	folded := OwnershipConflict([]*Capsule{left, otherCase})
	if len(folded) != 1 {
		t.Fatalf("case conflicts %v", folded)
	}
}

func TestCompileCopiesOwnedFiles(t *testing.T) {
	m := FromSeed(Seed{
		ID: "job", Objective: "Fix parser crash", DeclaredType: "BUGFIX",
		Acceptance: []string{"parser returns the node"}, QualityBar: BarStandard,
	})
	capsule, _, err := Compile(CompileRequest{Mission: m, Files: []string{"runtime/internal/mission/capsule.go"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(capsule.OwnedFiles) != 1 || capsule.OwnedFiles[0] != "runtime/internal/mission/capsule.go" {
		t.Fatalf("owned %v", capsule.OwnedFiles)
	}
}

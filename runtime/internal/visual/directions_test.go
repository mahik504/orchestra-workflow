package visual

import "testing"

func TestFlagshipDirectionsPolicy(t *testing.T) {
	if err := ValidateExploration(ClassOrdinary, DefaultFlagshipSet("x")); err == nil {
		t.Fatal("ordinary work must not generate directions")
	}
	if err := ValidateExploration(ClassProductUI, DefaultFlagshipSet("dashboard")); err == nil {
		t.Fatal("normal UI must not generate 3-5 directions")
	}
	dirs := DefaultFlagshipSet("portfolio")
	if err := ValidateExploration(ClassFlagship, dirs); err != nil {
		t.Fatal(err)
	}
	if len(dirs) < 3 || len(dirs) > 5 {
		t.Fatalf("len %d", len(dirs))
	}
	cosmetic := []Direction{
		{ID: "a", Concept: "blue", VisualGrammar: "cards", Typography: "sans", Composition: "centered hero", Interaction: "click cta"},
		{ID: "b", Concept: "green", VisualGrammar: "cards", Typography: "sans", Composition: "centered hero", Interaction: "click cta"},
		{ID: "c", Concept: "purple gradient", VisualGrammar: "cards", Typography: "sans", Composition: "centered hero", Interaction: "click cta"},
	}
	if err := ValidateExploration(ClassFlagship, cosmetic); err == nil {
		t.Fatal("cosmetic-only set must be rejected")
	}
	_, err := SelectOne(dirs, "")
	if err == nil {
		t.Fatal("empty selection must require a human/design decision")
	}
	sel, err := SelectOne(dirs, dirs[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	sealed, err := Seal(*sel)
	if err != nil || !sealed.Sealed || sealed.DirectionID != dirs[1].ID {
		t.Fatalf("seal: %+v %v", sealed, err)
	}
	if !Mutated(sealed, dirs[0]) {
		t.Fatal("switching direction after seal is a mutation")
	}
	if Mutated(sealed, *sel) {
		t.Fatal("same direction is not a mutation")
	}
}

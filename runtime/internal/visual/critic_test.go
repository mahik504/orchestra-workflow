package visual

import "testing"

func TestCriticNotABlacklist(t *testing.T) {
	got := Critique(CriticInput{
		Class:         ClassFlagship,
		HasGradient:   true,
		HasRounded:    true,
		HasDark:       true,
		HasCommonFont: true,
		HasFamiliarUI: true,
	})
	if len(got) != 0 {
		t.Fatalf("common techniques must not fail by themselves: %+v", got)
	}
}

func TestCriticActionableNotGeneric(t *testing.T) {
	got := Critique(CriticInput{
		Class:           ClassFlagship,
		CenteredHero:    true,
		NoProductAnchor: true,
		Surface:         "hero",
	})
	if len(got) == 0 {
		t.Fatal("expected an actionable composition finding")
	}
	if got[0].Explanation == "This looks generic." {
		t.Fatal("finding must not be a vague generic slur")
	}
	if got[0].Action == "" {
		t.Fatal("finding must be actionable")
	}
}

func TestCriticOrdinaryNotBlocked(t *testing.T) {
	got := Critique(CriticInput{Class: ClassOrdinary, CenteredHero: true, NoProductAnchor: true})
	if len(got) != 0 {
		t.Fatalf("ordinary work must not be blocked: %+v", got)
	}
}

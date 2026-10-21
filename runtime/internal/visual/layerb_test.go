package visual

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLayerBRequiresRenderedEvidence(t *testing.T) {
	rep := Evaluate(LayerBInput{Class: ClassFlagship, Review: "deep", SourceOnly: true})
	if rep.Passed {
		t.Fatal("source-only must not pass Layer B")
	}
	rep = Evaluate(LayerBInput{Class: ClassFlagship, Review: "deep", DOMOnly: true})
	if rep.Passed {
		t.Fatal("DOM-only must not pass Layer B")
	}
	rep = Evaluate(LayerBInput{Class: ClassFlagship, Review: "browser", Assertion: "looks good"})
	if rep.Passed {
		t.Fatal("assertion without screenshots must not pass")
	}
}

func TestLayerBStructuredFindingsFromScreenshots(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "desktop.png")
	if err := os.WriteFile(p, TinyPNG, 0644); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(LayerBInput{
		Class:  ClassFlagship,
		Review: "deep",
		Captures: []Capture{{
			Viewport: Viewport{Name: "desktop", Width: 1440, Height: 900},
			Route:    "/", ScreenshotPath: p, Passed: true,
		}},
		Critic: CriticInput{CenteredHero: true, NoProductAnchor: true, Surface: "hero"},
	})
	if rep.Passed {
		t.Fatal("expected visual findings")
	}
	if len(rep.Findings) == 0 || rep.Findings[0].Action == "" {
		t.Fatalf("structured findings required: %+v", rep)
	}
}

func TestLayerBOrdinaryNotRequired(t *testing.T) {
	rep := Evaluate(LayerBInput{Class: ClassOrdinary, Review: "none"})
	if !rep.Passed {
		t.Fatal(rep.Reason)
	}
}

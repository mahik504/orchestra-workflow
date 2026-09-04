package classifier

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/orchestra-v3/internal/resources"
)

func loadRealGraph(t *testing.T) *resources.DesignResourceGraph {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "..", "registries", "design-resource-graph.json"),
		filepath.Join("..", "..", "registries", "design-resource-graph.json"),
		`C:\projects\orchestra-workflow\registries\design-resource-graph.json`,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			g, err := resources.LoadDesignGraph(c)
			if err != nil {
				t.Fatalf("load graph: %v", err)
			}
			return g
		}
	}
	t.Fatalf("could not find design-resource-graph.json")
	return nil
}

func TestProbeBriefs(t *testing.T) {
	if os.Getenv("ORCHESTRA_PROBE") == "" {
		t.Skip("set ORCHESTRA_PROBE=1 to print classifier scores")
	}
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	briefs := []string{
		"Build a landing page for a friend's coffee roastery, they want it to feel expensive",
		"The login form throws a 500 when the email has a plus sign",
		"I want a portfolio site that also sells prints, with a checkout and an admin area to manage orders",
		"Internal tool for our team to see which servers are down right now, live",
		"A reading app for arXiv papers with proper footnotes and math",
		"Add haptic press feedback to the buttons in our existing Expo app",
		"Audit our API for injection and check the dependencies",
		"Extract the design tokens from https://example.com and write them up",
		"Build a scheduling dashboard for a school with attendance charts",
	}

	for _, b := range briefs {
		br := c.ClassifyBrief(b, Options{})
		fmt.Printf("\n=== %s\n", b)
		fmt.Printf("  type=%s archetype=%s cap=%s bar=%s lab=%v visual=%v research=%s verify=%s\n",
			br.Type, br.Archetype, br.CapabilityID, br.QualityBar, br.DesignLabRequired, br.RequiresVisual, br.ResearchDepth, br.VerifyDepth)
		for i, s := range br.Selected {
			if i > 2 {
				break
			}
			fmt.Printf("    sel %-22s %.2f tags=%v\n", s.CapabilityID, s.Score, s.MatchedTags)
		}
		if br.Ambiguous {
			fmt.Printf("  ASK: %s\n", br.ClarifyingQuestion)
			br.ResolveSilence()
			fmt.Printf("  SILENCE -> %s (%s)\n", br.CapabilityID, br.AssumptionNote)
		}
	}
}

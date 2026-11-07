package classifier

import (
	"strings"
	"testing"
)

func TestV33_SeventeenCapabilitiesHaveTriggerAndSkip(t *testing.T) {
	g := loadRealGraph(t)
	if len(g.Capabilities) != 26 {
		t.Fatalf("capabilities = %d, want 26", len(g.Capabilities))
	}
	for _, id := range []string{
		"interaction-components",
		"ow-reference-forensics",
		"ow-golden-still",
		"ow-fresh-review",
		"ow-showcase",
	} {
		cap, ok := g.Capabilities[id]
		if !ok {
			t.Fatalf("missing %s", id)
		}
		if len(cap.TriggerConditions) == 0 || len(cap.SkipConditions) == 0 {
			t.Errorf("%s missing trigger or skip", id)
		}
	}
}

func TestV33_ClassifySmoke(t *testing.T) {
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	cases := []struct {
		name    string
		brief   string
		wantCap string
		wantLab bool
	}{
		{"premium website", "Build a premium landing page for a studio brand, award-winning marketing site", "premium-website", true},
		{"3d portfolio", "Build a 3D WebGL portfolio with R3F camera orbits and custom shaders", "3d-portfolio", true},
		{"saas dashboard", "Build a scheduling dashboard for a school with attendance charts and KPIs", "saas-dashboard", false},
		{"mobile app", "Build a mobile onboarding flow in Expo for iOS with touch gestures", "mobile-app", true},
		{"research paper", "Write a methods paper with citations about capability graphs, no UI", "ow-research-paper", false},
		{"security", "Audit our API for injection and check the dependencies", "security-audit", false},
		{"interaction control", "Add a magnetic gradient button using Kokonut after inspecting existing components", "interaction-components", true},
		{"tint reference", "Use this screenshot as a reference, keep the structure, change orange to light blue", "reverse-engineering", false},
		{"game-like 3d", "Build an interactive 3D WebGL product UI with R3F, Three.js, and Drei — cinematic presentation", "3d-portfolio", true},
		{"visual forensics", "Run ow-reference-forensics on a named-live-site: composition-forensics, mechanism-vs-identity, write REFERENCE_BENCHMARK.md", "ow-reference-forensics", true},
		{"golden still", "Run ow-golden-still for a desktop-still and mobile-still opening-scene-still PASS/FAIL only", "ow-golden-still", true},
		{"fresh review", "Run ow-fresh-review independent-review builder-blind-spot on a completed phase for pre-release-review", "ow-fresh-review", false},
		{"final showcase", "Record ow-showcase brag walkthrough-video ship-assets after production-ready green ledger", "ow-showcase", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := c.ClassifyBrief(tc.brief, Options{})
			if b.CapabilityID != tc.wantCap {
				t.Errorf("cap = %q (%s), want %q; top=%v", b.CapabilityID, b.ArchetypeReason, tc.wantCap, summarize(b))
			}
			if b.DesignLabRequired != tc.wantLab {
				t.Errorf("design lab = %v (%s), want %v", b.DesignLabRequired, b.DesignLabReason, tc.wantLab)
			}
		})
	}
}

func TestV33_BackendDoesNotArmLabOr3D(t *testing.T) {
	g := loadRealGraph(t)
	b := NewClassifierWithGraph(g).ClassifyBrief(
		"Add Postgres indexes and Row Level Security on a Supabase schema, no UI",
		Options{},
	)
	if b.DesignLabRequired {
		t.Fatalf("backend/database armed Design Lab: %s / %s", b.CapabilityID, b.DesignLabReason)
	}
	if b.CapabilityID == "3d-portfolio" || b.CapabilityID == "premium-website" {
		t.Fatalf("backend routed to %s", b.CapabilityID)
	}
}

func summarize(b *Brief) string {
	var parts []string
	n := 3
	if len(b.Selected) < n {
		n = len(b.Selected)
	}
	for i := 0; i < n; i++ {
		parts = append(parts, b.Selected[i].CapabilityID)
	}
	return strings.Join(parts, ",")
}

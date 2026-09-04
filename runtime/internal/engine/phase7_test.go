package engine

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/user/orchestra-v3/internal/adapters/acquisition"
	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/research"
)

// Phase 7 briefs. Each one must resolve differently enough that a restaurant,
// a school SaaS, and a 3D portfolio cannot collapse into the same graph.
type phase7Brief struct {
	name    string
	request string
	wantCap string
	wantLab bool
	wantBar string
}

func phase7Briefs() []phase7Brief {
	return []phase7Brief{
		{
			name:    "1 backend bug",
			request: "The login form throws a 500 when the email has a plus sign",
			wantCap: "",
			wantLab: false,
			wantBar: classifier.BarStandard,
		},
		{
			name:    "2 school saas",
			request: "Build a scheduling dashboard for a school with attendance charts",
			wantCap: "saas-dashboard",
			wantLab: false,
			wantBar: classifier.BarStandard,
		},
		{
			name:    "3 restaurant premium",
			request: "Build a landing page for a friend's coffee roastery, they want it to feel expensive",
			wantCap: "premium-website",
			wantLab: true,
			wantBar: classifier.BarPremium,
		},
		{
			name:    "4 3d portfolio",
			request: "Build a 3D WebGL portfolio with R3F camera orbits and custom shaders",
			wantCap: "3d-portfolio",
			wantLab: true,
			wantBar: classifier.BarExperimental,
		},
		{
			name:    "5 research paper",
			request: "Write a methods paper with citations about capability graphs, no UI",
			wantCap: "research-paper",
			wantLab: false,
			wantBar: classifier.BarStandard,
		},
		{
			name:    "6 unnamed libraries",
			request: "Build a distinctive restaurant website. Make it feel expensive. Do not pick a library for me.",
			wantCap: "premium-website",
			wantLab: true,
			wantBar: classifier.BarPremium,
		},
		{
			name:    "7 ambiguous shop portfolio",
			request: "I want a portfolio site that also sells prints, with a checkout and an admin area to manage orders",
			wantCap: "premium-website", // first pick; silence falls to b2b-portal
			wantLab: true,
			wantBar: classifier.BarPremium,
		},
	}
}

func TestPhase7_SevenBriefs(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)

	for _, tc := range phase7Briefs() {
		t.Run(tc.name, func(t *testing.T) {
			b := cls.ClassifyBrief(tc.request, classifier.Options{})
			if tc.name == "7 ambiguous shop portfolio" {
				if !b.Ambiguous {
					t.Fatalf("expected one clarifying question, got a silent pick of %s (%s)", b.CapabilityID, b.ArchetypeReason)
				}
				if strings.Count(b.ClarifyingQuestion, "?") != 1 {
					t.Errorf("expected exactly one question, got: %s", b.ClarifyingQuestion)
				}
				if b.CapabilityID != "premium-website" {
					t.Errorf("first pick = %q, want premium-website", b.CapabilityID)
				}
				b.ResolveSilence()
				if !b.Assumed || !strings.Contains(b.AssumptionNote, "no response") {
					t.Errorf("silence fallback missing assume log: %q", b.AssumptionNote)
				}
				if b.CapabilityID != "b2b-portal" {
					t.Errorf("silence picked %s, want b2b-portal (lower risk)", b.CapabilityID)
				}
				return
			}
			if b.CapabilityID != tc.wantCap {
				t.Errorf("capability = %q, want %q (%s)", b.CapabilityID, tc.wantCap, b.ArchetypeReason)
			}
			if b.DesignLabRequired != tc.wantLab {
				t.Errorf("design lab = %v, want %v (%s)", b.DesignLabRequired, tc.wantLab, b.DesignLabReason)
			}
			if b.QualityBar != tc.wantBar {
				t.Errorf("quality bar = %q, want %q", b.QualityBar, tc.wantBar)
			}
		})
	}
}

func TestPhase7_BackendBugLoadsNoExpoNo3D(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)
	b := cls.ClassifyBrief("The login form throws a 500 when the email has a plus sign", classifier.Options{})
	if b.DesignLabRequired {
		t.Fatal("backend bug armed the Design Lab")
	}
	if b.CapabilityID != "" {
		if route, ok := graph.ResolveCapabilityRoute(b.CapabilityID, nil); ok {
			blob := strings.ToLower(strings.Join(route.AllResourceIDs, " "))
			for _, banned := range []string{"expo", "r3f", "shadergradient", "threeui"} {
				if strings.Contains(blob, banned) {
					t.Errorf("backend bug route loaded %s: %v", banned, route.AllResourceIDs)
				}
			}
		}
	}
}

func TestPhase7_ThreeGraphsDiffer(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)

	ids := []string{"saas-dashboard", "premium-website", "3d-portfolio"}
	briefs := map[string]string{
		"saas-dashboard":   "Build a scheduling dashboard for a school with attendance charts",
		"premium-website":  "Build a landing page for a friend's coffee roastery, they want it to feel expensive",
		"3d-portfolio":     "Build a 3D WebGL portfolio with R3F camera orbits and custom shaders",
	}

	sets := map[string]string{}
	for _, id := range ids {
		b := cls.ClassifyBrief(briefs[id], classifier.Options{})
		if b.CapabilityID != id {
			t.Fatalf("%s classified as %s", id, b.CapabilityID)
		}
		route, ok := graph.ResolveCapabilityRoute(id, nil)
		if !ok {
			t.Fatalf("no route for %s", id)
		}
		sets[id] = strings.Join(route.AllResourceIDs, "|")
		t.Logf("%s resources (%d): %s", id, len(route.AllResourceIDs), strings.Join(route.AllResourceIDs, ", "))
	}
	if sets["saas-dashboard"] == sets["premium-website"] ||
		sets["saas-dashboard"] == sets["3d-portfolio"] ||
		sets["premium-website"] == sets["3d-portfolio"] {
		t.Fatal("restaurant / school SaaS / 3D portfolio resolved to the same resource graph")
	}
}

func TestPhase7_3DCardNamesR3FAndShaderGradient(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	route, ok := graph.ResolveCapabilityRoute("3d-portfolio", []string{"3d", "r3f", "shader"})
	if !ok {
		t.Fatal("3d-portfolio missing from graph")
	}
	blob := strings.Join(append(append([]string{}, route.Implementation...), route.AllResourceIDs...), " ")
	if !strings.Contains(blob, "r3f") {
		t.Errorf("3D route does not name r3f: %v", route.Implementation)
	}
	if !strings.Contains(blob, "shadergradient") {
		t.Errorf("3D route does not name shadergradient: %v", route.Implementation)
	}
}

func TestPhase7_PaperLoadsOrchestraDocs(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)
	b := cls.ClassifyBrief("Write a methods paper with citations about capability graphs, no UI", classifier.Options{})
	if b.CapabilityID != "research-paper" {
		t.Fatalf("got %s, want research-paper", b.CapabilityID)
	}
	if b.DesignLabRequired {
		t.Error("a paper is not visual work")
	}
	route, ok := graph.ResolveCapabilityRoute("research-paper", nil)
	if !ok {
		t.Fatal("research-paper missing")
	}
	found := false
	for _, id := range route.AllResourceIDs {
		if id == "orchestra-docs" {
			found = true
		}
	}
	if !found {
		t.Errorf("research-paper route missing orchestra-docs: %v", route.AllResourceIDs)
	}
}

func TestPhase7_UnnamedLibrariesStillSurfaceResearchDesignMotion(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)
	req := "Build a distinctive restaurant website. Make it feel expensive. Do not pick a library for me."
	b := cls.ClassifyBrief(req, classifier.Options{})
	if b.CapabilityID != "premium-website" {
		t.Fatalf("got %s, want premium-website", b.CapabilityID)
	}
	// The brief named no libraries. The graph still has to produce classes.
	route, ok := graph.ResolveCapabilityRoute("premium-website", b.Selected[0].MatchedTags)
	if !ok {
		t.Fatal("premium-website missing")
	}

	has := func(ids []string, needle string) bool {
		for _, id := range ids {
			if strings.Contains(id, needle) {
				return true
			}
		}
		return false
	}
	if !has(route.DiscoveryDomains, "visual_research") && !has(route.DiscoveryResources, "awwwards") {
		t.Errorf("unnamed premium brief did not surface research: domains=%v res=%v", route.DiscoveryDomains, route.DiscoveryResources)
	}
	if !has(route.SynthesisDomains, "design_synthesis") && !has(route.SynthesisResources, "taste-design") {
		t.Errorf("unnamed premium brief did not surface design: domains=%v res=%v", route.SynthesisDomains, route.SynthesisResources)
	}
	if !has(route.Implementation, "gsap") && !has(route.Implementation, "motion") && !has(route.AllResourceIDs, "gsap") {
		t.Errorf("unnamed premium brief did not surface motion: impl=%v all=%v", route.Implementation, route.AllResourceIDs)
	}

	// Catalog presence is not usage: the full catalog is much larger than this route.
	if cat != nil && len(route.AllResourceIDs) >= len(cat.All()) {
		t.Errorf("route loaded the whole catalog (%d) instead of a selected set", len(route.AllResourceIDs))
	}
}

func TestPhase7_UnknownLibraryIsSurfacedNotForced(t *testing.T) {
	_, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	cls := classifier.NewClassifierWithGraph(graph)
	b := cls.ClassifyBrief(
		"Build a landing page using pixi.js for the hero canvas",
		classifier.Options{},
	)
	found := false
	for _, u := range b.UnknownTechnology {
		if u == "pixi.js" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected pixi.js in unknown technology, got %v (cap=%s)", b.UnknownTechnology, b.CapabilityID)
	}
}

func TestPhase7_MultiSourceOneWorld(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	route, ok := graph.ResolveCapabilityRoute("premium-website", []string{"landing-page"})
	if !ok {
		t.Fatal("premium-website missing")
	}
	coord := research.NewResearchCoordinator(cat, graph)
	res, err := coord.Coordinate(context.Background(), &research.ResearchRequest{
		Task: &classifier.Task{
			ID:             "phase7-research",
			RawRequest:     "Build a landing page for a friend's coffee roastery",
			Type:           classifier.TypeFeature,
			RequiresVisual: true,
		},
		Route:            route,
		ProjectOutputDir: workdir,
		Options: research.SelectionOptions{
			MinSources:       2,
			MaxSources:       4,
			OfflineBenchmark: true,
		},
	})
	if err != nil {
		t.Fatalf("research: %v", err)
	}
	if len(res.SelectedSources) < 2 {
		t.Fatalf("expected at least 2 named references, got %d", len(res.SelectedSources))
	}
	names := map[string]bool{}
	for _, s := range res.SelectedSources {
		if s == nil || s.Resource == nil || s.Resource.ID == "" {
			t.Fatal("selected a nameless reference")
		}
		names[s.Resource.ID] = true
		t.Logf("reference: %s (%s)", s.Resource.ID, s.Resource.Name)
	}
	if len(names) < 2 {
		t.Fatal("references collapsed to a single source")
	}
	// One coherent world: one display face, one accent. A collage would ship four display fonts.
	display := 0
	accents := map[string]bool{}
	for _, ty := range res.SynthesizedTypography {
		if ty.Role == "display" {
			display++
		}
	}
	for _, p := range res.SynthesizedPalette {
		if p.Role == "--color-accent-primary" {
			accents[strings.ToLower(p.Hex)] = true
		}
	}
	if display != 1 {
		t.Errorf("expected one display face, got %d (collage)", display)
	}
	if len(accents) != 1 {
		t.Errorf("expected one accent, got %d (collage)", len(accents))
	}
}

func TestPhase7_AcquisitionScopes(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	awwwards, ok := cat.FindByID("awwwards")
	if !ok {
		t.Fatal("awwwards missing from catalog")
	}
	if awwwards.AcquisitionMethod != "web_fetch" || awwwards.RuntimeMethod != "on_demand_research" {
		t.Errorf("reference awwwards should be on-demand fetch, got acq=%s runtime=%s", awwwards.AcquisitionMethod, awwwards.RuntimeMethod)
	}

	gsap, ok := cat.FindByID("gsap")
	if !ok {
		t.Fatal("gsap missing from catalog")
	}
	if gsap.AcquisitionMethod != "npm" || gsap.RuntimeMethod != "project_scoped_install" {
		t.Errorf("impl library gsap should be project-scoped npm, got acq=%s runtime=%s", gsap.AcquisitionMethod, gsap.RuntimeMethod)
	}

	err := acquisition.CheckGlobalInstallSafety([]string{"install", "-g", "gsap"}, workdir)
	if err == nil {
		t.Fatal("global npm install was not blocked")
	}

	err = acquisition.CheckGlobalInstallSafety([]string{"install", "--save", "gsap"}, workdir)
	if err != nil {
		t.Errorf("project-scoped install was blocked: %v", err)
	}

	// Premium route may be expensive. Cost is recorded, not used as a fail.
	route, _ := graph.ResolveCapabilityRoute("premium-website", nil)
	if route == nil || len(route.AllResourceIDs) == 0 {
		t.Fatal("premium-website has no resources to cost")
	}
	t.Logf("premium-website resource count (cost proxy) = %d; not a failure", len(route.AllResourceIDs))
}

func TestPhase7_UsedIsNotTheWholeCatalog(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	route, ok := graph.ResolveCapabilityRoute("premium-website", nil)
	if !ok {
		t.Fatal("premium-website missing")
	}
	all := cat.All()
	if len(all) == 0 {
		t.Fatal("empty catalog")
	}
	if len(route.AllResourceIDs) == 0 {
		t.Fatal("route selected nothing")
	}
	if len(route.AllResourceIDs) >= len(all) {
		t.Fatalf("used %d resources, catalog has %d — selected set is not smaller than listed set", len(route.AllResourceIDs), len(all))
	}
}

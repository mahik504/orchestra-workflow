package resources

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleResourcesJSON = `[
  {
    "id": "gsap",
    "name": "GSAP",
    "canonical_url": "https://gsap.com/",
    "source_type": "github_repository",
    "source_repository": "https://github.com/greensock/GSAP.git",
    "category": ["motion", "frontend"],
    "representation": "dependency",
    "routing_tags": ["motion", "microinteraction", "timeline", "scroll"],
    "acquisition_method": "npm",
    "runtime_method": "project_scoped_install",
    "status": "ACTIVE",
    "trigger_conditions": ["complex timeline animation", "orchestrated motion"],
    "avoid_conditions": ["simple hover transition"]
  },
  {
    "id": "awwwards",
    "name": "Awwwards",
    "canonical_url": "https://www.awwwards.com/",
    "source_type": "web_reference",
    "category": ["visual_research", "design_reference"],
    "representation": "research_source",
    "routing_tags": ["premium", "award-winning", "layout", "typography"],
    "acquisition_method": "web_fetch",
    "runtime_method": "on_demand_research",
    "status": "ACTIVE",
    "trigger_conditions": ["premium visual quality", "landing page", "creative direction"],
    "avoid_conditions": ["basic admin screens", "backend bugs"]
  },
  {
    "id": "lenis",
    "name": "Lenis",
    "canonical_url": "https://lenis.darkroom.engineering/",
    "source_type": "github_repository",
    "source_repository": "https://github.com/darkroomengineering/lenis.git",
    "category": ["motion", "scroll"],
    "representation": "dependency",
    "routing_tags": ["smooth scroll", "parallax", "motion"],
    "acquisition_method": "npm",
    "runtime_method": "project_scoped_install",
    "status": "ACTIVE"
  },
  {
    "id": "r3f",
    "name": "React Three Fiber",
    "canonical_url": "https://docs.pmnd.rs/react-three-fiber",
    "source_type": "github_repository",
    "source_repository": "https://github.com/pmndrs/react-three-fiber.git",
    "category": ["3d", "webgl"],
    "representation": "dependency",
    "routing_tags": ["3d", "canvas", "webgl", "threejs"],
    "acquisition_method": "npm",
    "runtime_method": "project_scoped_install",
    "status": "ACTIVE"
  }
]`

const sampleDesignGraphJSON = `{
  "$schema": "./schemas/design-resource-graph.schema.json",
  "version": "3.0.0",
  "description": "Test design graph",
  "domains": {
    "visual_research": ["awwwards", "jiro", "cari", "designmd"],
    "design_synthesis": ["taste-design", "impeccable", "emil-design-eng"],
    "motion": ["gsap", "lenis"],
    "webgl": ["r3f", "shadergradient"]
  },
  "capabilities": {
    "premium-website": {
      "name": "Premium Website",
      "description": "Premium creative site",
      "primary_archetype": "creative_showcase",
      "trigger_tags": ["premium", "landing-page"],
      "discovery": ["visual_research"],
      "synthesis": ["design_synthesis"],
      "implementation": ["motion"],
      "optional_extensions": ["webgl"],
      "qa": ["playwright", "lighthouse"]
    },
    "3d-portfolio": {
      "name": "3D Portfolio",
      "description": "Interactive 3D showcase",
      "primary_archetype": "spatial_experience",
      "trigger_tags": ["3d", "portfolio", "webgl"],
      "discovery": ["visual_research"],
      "synthesis": ["design_synthesis"],
      "implementation": ["motion", "webgl"],
      "qa": ["playwright"]
    }
  }
}`

func createTempFile(t *testing.T, dir, filename, content string) string {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file %s: %v", path, err)
	}
	return path
}

func TestLoadResourceCatalog_Success(t *testing.T) {
	tempDir := t.TempDir()
	path := createTempFile(t, tempDir, "resources.json", sampleResourcesJSON)

	cat, err := LoadResourceCatalog(path)
	if err != nil {
		t.Fatalf("LoadResourceCatalog returned unexpected error: %v", err)
	}

	if cat.Count() != 4 {
		t.Errorf("expected count 4, got %d", cat.Count())
	}

	all := cat.All()
	if len(all) != 4 {
		t.Errorf("expected 4 items in All(), got %d", len(all))
	}

	gsap, ok := cat.FindByID("gsap")
	if !ok || gsap == nil {
		t.Fatalf("failed to find resource 'gsap'")
	}
	if gsap.Name != "GSAP" {
		t.Errorf("expected name 'GSAP', got %s", gsap.Name)
	}
	if gsap.CanonicalURL != "https://gsap.com/" {
		t.Errorf("expected canonical URL 'https://gsap.com/', got %s", gsap.CanonicalURL)
	}
	if gsap.AcquisitionMethod != "npm" {
		t.Errorf("expected acquisition method 'npm', got %s", gsap.AcquisitionMethod)
	}
	if gsap.RuntimeMethod != "project_scoped_install" {
		t.Errorf("expected runtime method 'project_scoped_install', got %s", gsap.RuntimeMethod)
	}
	if len(gsap.RoutingTags) != 4 {
		t.Errorf("expected 4 routing tags for gsap, got %d", len(gsap.RoutingTags))
	}
}

func TestResourceCatalog_LookupsAndIndexing(t *testing.T) {
	tempDir := t.TempDir()
	path := createTempFile(t, tempDir, "resources.json", sampleResourcesJSON)

	cat, err := LoadResourceCatalog(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Case-insensitive FindByID
	if res, ok := cat.FindByID("GSAP"); !ok || res.ID != "gsap" {
		t.Errorf("expected case-insensitive lookup for GSAP to succeed")
	}
	if _, ok := cat.FindByID("non-existent"); ok {
		t.Errorf("expected non-existent ID lookup to fail")
	}

	// 2. FindByTag
	motionRes := cat.FindByTag("motion")
	if len(motionRes) != 2 { // gsap and lenis
		t.Errorf("expected 2 resources tagged 'motion', got %d", len(motionRes))
	}

	// 3. FindByTags (Multi-tag union with deduplication)
	multiRes := cat.FindByTags([]string{"motion", "webgl", "typography"})
	if len(multiRes) != 4 { // gsap, lenis, r3f, awwwards
		t.Errorf("expected 4 unique resources for multi-tag lookup, got %d", len(multiRes))
	}

	// 4. FindByCategory
	frontends := cat.FindByCategory("frontend")
	if len(frontends) != 1 || frontends[0].ID != "gsap" {
		t.Errorf("expected 1 resource in category 'frontend', got %d", len(frontends))
	}

	// 5. FindByAcquisitionMethod
	npmRes := cat.FindByAcquisitionMethod("npm")
	if len(npmRes) != 3 { // gsap, lenis, r3f
		t.Errorf("expected 3 resources with acquisition method 'npm', got %d", len(npmRes))
	}
	webRes := cat.FindByAcquisitionMethod("web_fetch")
	if len(webRes) != 1 || webRes[0].ID != "awwwards" {
		t.Errorf("expected 1 resource with acquisition method 'web_fetch', got %d", len(webRes))
	}

	// 6. FindByRepresentation
	depRes := cat.FindByRepresentation("dependency")
	if len(depRes) != 3 {
		t.Errorf("expected 3 dependencies, got %d", len(depRes))
	}
	researchRes := cat.FindByRepresentation("research_source")
	if len(researchRes) != 1 {
		t.Errorf("expected 1 research_source, got %d", len(researchRes))
	}

	// 7. Tags listing
	tags := cat.Tags()
	if len(tags) == 0 {
		t.Errorf("expected indexed tags list to be non-empty")
	}
}

func TestLoadDesignGraph_Success(t *testing.T) {
	tempDir := t.TempDir()
	path := createTempFile(t, tempDir, "design-resource-graph.json", sampleDesignGraphJSON)

	graph, err := LoadDesignGraph(path)
	if err != nil {
		t.Fatalf("LoadDesignGraph returned unexpected error: %v", err)
	}

	if len(graph.Domains) != 4 {
		t.Errorf("expected 4 domains, got %d", len(graph.Domains))
	}

	if len(graph.Capabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(graph.Capabilities))
	}

	// Domain resources check
	visRes := graph.GetDomainResources("visual_research")
	if len(visRes) != 4 {
		t.Errorf("expected 4 visual_research resources, got %d", len(visRes))
	}
	if graph.GetDomainResources("non_existent") != nil {
		t.Errorf("expected nil for non_existent domain")
	}

	caps := graph.AllCapabilities()
	if len(caps) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(caps))
	}

	doms := graph.AllDomains()
	if len(doms) != 4 {
		t.Errorf("expected 4 domains, got %d", len(doms))
	}
}

func TestDesignResourceGraph_ResolveCapabilities(t *testing.T) {
	tempDir := t.TempDir()
	path := createTempFile(t, tempDir, "design-resource-graph.json", sampleDesignGraphJSON)

	graph, err := LoadDesignGraph(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Direct capability ID match
	routes := graph.ResolveCapabilities([]string{"premium-website"})
	if len(routes) != 1 {
		t.Fatalf("expected 1 route for 'premium-website', got %d", len(routes))
	}
	r := routes[0]
	if r.CapabilityName != "premium-website" {
		t.Errorf("expected capability 'premium-website', got %s", r.CapabilityName)
	}
	if len(r.DiscoveryResources) != 4 { // awwwards, jiro, cari, designmd
		t.Errorf("expected 4 discovery resources, got %d", len(r.DiscoveryResources))
	}
	if len(r.SynthesisResources) != 3 { // taste-design, impeccable, emil-design-eng
		t.Errorf("expected 3 synthesis resources, got %d", len(r.SynthesisResources))
	}
	if len(r.Implementation) != 2 { // gsap, lenis
		t.Errorf("expected 2 implementation resources, got %d", len(r.Implementation))
	}
	if len(r.OptionalExtensions) != 2 { // r3f, shadergradient
		t.Errorf("expected 2 optional extensions, got %d", len(r.OptionalExtensions))
	}
	if len(r.QA) != 2 { // playwright, lighthouse
		t.Errorf("expected 2 QA resources, got %d", len(r.QA))
	}

	// 2. Trigger tag match
	tagRoutes := graph.ResolveCapabilities([]string{"landing-page"})
	if len(tagRoutes) != 1 || tagRoutes[0].CapabilityName != "premium-website" {
		t.Errorf("expected trigger tag 'landing-page' to resolve 'premium-website'")
	}

	// 3. Short tag / substring match (e.g. "3d")
	shortTagRoutes := graph.ResolveCapabilities([]string{"3d"})
	if len(shortTagRoutes) != 1 || shortTagRoutes[0].CapabilityName != "3d-portfolio" {
		t.Errorf("expected '3d' tag to resolve '3d-portfolio', got %v", shortTagRoutes)
	}

	// 4. Domain match (e.g. "webgl" activates both capabilities referencing webgl)
	webglRoutes := graph.ResolveCapabilities([]string{"webgl"})
	if len(webglRoutes) != 2 {
		t.Errorf("expected 'webgl' domain match to resolve 2 capabilities, got %d", len(webglRoutes))
	}

	// 5. Empty tags
	emptyRoutes := graph.ResolveCapabilities([]string{})
	if len(emptyRoutes) != 0 {
		t.Errorf("expected empty routes for empty tags, got %d", len(emptyRoutes))
	}
}

func TestBackwardsCompatibility_Registry(t *testing.T) {
	tempDir := t.TempDir()
	path := createTempFile(t, tempDir, "resources.json", sampleResourcesJSON)

	cat, err := LoadResourceCatalog(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reg := cat.ToRegistry()
	if reg == nil || len(reg.Capabilities) != 4 {
		t.Fatalf("expected registry with 4 capabilities, got %v", reg)
	}

	gsapCap, ok := reg.Capabilities["gsap"]
	if !ok || gsapCap == nil {
		t.Fatalf("failed to find 'gsap' capability in converted registry")
	}
	if gsapCap.Category != CategorySpecialist {
		t.Errorf("expected CategorySpecialist for gsap, got %s", gsapCap.Category)
	}
	if !gsapCap.RuntimeDependency {
		t.Errorf("expected RuntimeDependency to be true for gsap")
	}
	if gsapCap.InstallationMethod != "npm" {
		t.Errorf("expected InstallationMethod 'npm', got %s", gsapCap.InstallationMethod)
	}

	// Test direct Registry.LoadResourceCatalog
	newReg := NewRegistry()
	if err := newReg.LoadResourceCatalog(path); err != nil {
		t.Fatalf("LoadResourceCatalog on Registry failed: %v", err)
	}
	if len(newReg.Capabilities) != 4 {
		t.Errorf("expected 4 capabilities in newReg, got %d", len(newReg.Capabilities))
	}
}

func TestWindowsBOMHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Prepend UTF-8 BOM \xef\xbb\xbf
	bomBytes := []byte{0xef, 0xbb, 0xbf}
	catalogContent := append(bomBytes, []byte(sampleResourcesJSON)...)
	graphContent := append(bomBytes, []byte(sampleDesignGraphJSON)...)

	catPath := filepath.Join(tempDir, "bom_resources.json")
	if err := os.WriteFile(catPath, catalogContent, 0644); err != nil {
		t.Fatalf("failed to write BOM catalog: %v", err)
	}

	graphPath := filepath.Join(tempDir, "bom_graph.json")
	if err := os.WriteFile(graphPath, graphContent, 0644); err != nil {
		t.Fatalf("failed to write BOM graph: %v", err)
	}

	// 1. Verify LoadResourceCatalog handles BOM
	cat, err := LoadResourceCatalog(catPath)
	if err != nil {
		t.Fatalf("LoadResourceCatalog failed with BOM: %v", err)
	}
	if cat.Count() != 4 {
		t.Errorf("expected 4 resources from BOM file, got %d", cat.Count())
	}

	// 2. Verify LoadDesignGraph handles BOM
	graph, err := LoadDesignGraph(graphPath)
	if err != nil {
		t.Fatalf("LoadDesignGraph failed with BOM: %v", err)
	}
	if len(graph.Capabilities) != 2 {
		t.Errorf("expected 2 capabilities from BOM file, got %d", len(graph.Capabilities))
	}

	// 3. Verify Registry.LoadFromJSON handles BOM
	reg := NewRegistry()
	if err := reg.LoadFromJSON(catPath); err != nil {
		t.Fatalf("Registry.LoadFromJSON failed with BOM: %v", err)
	}
	if len(reg.Capabilities) != 4 {
		t.Errorf("expected 4 capabilities in reg from BOM file, got %d", len(reg.Capabilities))
	}
}

func TestLiveRegistryFiles(t *testing.T) {
	// Locate repository registries directory relative to runtime/internal/resources
	workflowRoot := filepath.Join("..", "..", "..")
	resourcesPath := filepath.Join(workflowRoot, "registries", "resources.json")
	graphPath := filepath.Join(workflowRoot, "registries", "design-resource-graph.json")

	// 1. Validate resources.json has all 126 canonical resources
	cat, err := LoadResourceCatalog(resourcesPath)
	if err != nil {
		t.Fatalf("Failed to load live resources.json: %v", err)
	}
	if cat.Count() != 126 {
		t.Errorf("expected exactly 126 resources in canonical registry, got %d", cat.Count())
	}

	// Verify required canonical resources exist
	requiredResources := []string{
		"gsap",
		"lenis",
		"r3f",
		"awwwards",
		"godly",
		"refero",
		"jiro",
		"cari",
		"designmd",
		"taste-design",
		"impeccable",
		"emil-design-eng",
		"pretext",
		"trig-js",
		"playwright",
		"strix",
		"semgrep-adapter",
	}

	for _, id := range requiredResources {
		if _, ok := cat.FindByID(id); !ok {
			t.Errorf("required canonical resource %q was not found in live resources.json", id)
		}
	}

	// 2. Validate design-resource-graph.json has 20 domains and 11 capabilities
	graph, err := LoadDesignGraph(graphPath)
	if err != nil {
		t.Fatalf("Failed to load live design-resource-graph.json: %v", err)
	}

	if len(graph.Domains) < 20 {
		t.Errorf("expected at least 20 domains in live design-resource-graph.json, got %d", len(graph.Domains))
	}
	if len(graph.Capabilities) < 11 {
		t.Errorf("expected at least 11 capabilities in live design-resource-graph.json, got %d", len(graph.Capabilities))
	}

	// Verify all 7 required archetypes from prompt are present
	requiredCapabilities := []string{
		"premium-website",
		"3d-portfolio",
		"operator-hud",
		"b2b-portal",
		"academic-reader",
		"micro-interactions",
		"physics-canvas",
	}

	for _, capID := range requiredCapabilities {
		route, ok := graph.ResolveCapabilityRoute(capID, []string{capID})
		if !ok {
			t.Errorf("required capability %q was not found in live graph", capID)
			continue
		}
		if len(route.DiscoveryResources) == 0 && len(route.DiscoveryDomains) == 0 && capID != "security-audit" {
			t.Errorf("capability %q has empty discovery phase", capID)
		}
		if len(route.SynthesisResources) == 0 && len(route.SynthesisDomains) == 0 {
			t.Errorf("capability %q has empty synthesis phase", capID)
		}
		if len(route.Implementation) == 0 {
			t.Errorf("capability %q has empty implementation phase", capID)
		}
		if len(route.QA) == 0 {
			t.Errorf("capability %q has empty QA phase", capID)
		}
	}
}

func TestCapabilityRoute_QA_DomainExpansion(t *testing.T) {
	graph := &DesignResourceGraph{
		Domains: map[string][]string{
			"qa_domain": {"playwright", "lighthouse"},
		},
		Capabilities: map[string]*CapabilityPhaseDefinition{
			"test-qa-cap": {
				Name:             "Test QA Cap",
				PrimaryArchetype: "quality_assurance",
				Discovery:        []string{"res-1"},
				Synthesis:        []string{"res-2"},
				Implementation:   []string{"res-3"},
				QA:               []string{"qa_domain", "custom-qa-tool"},
			},
			"test-empty-qa": {
				Name:             "Test Empty QA",
				PrimaryArchetype: "minimal",
				Discovery:        []string{"res-1"},
				Synthesis:        []string{"res-2"},
				Implementation:   []string{"res-3"},
				QA:               nil,
			},
		},
	}

	// 1. Verify domain expansion in QA
	route, ok := graph.ResolveCapabilityRoute("test-qa-cap", []string{"test"})
	if !ok {
		t.Fatalf("expected capability route to resolve")
	}
	expectedQA := []string{"playwright", "lighthouse", "custom-qa-tool"}
	if len(route.QA) != len(expectedQA) {
		t.Fatalf("expected %d QA resources after domain expansion, got %d: %v", len(expectedQA), len(route.QA), route.QA)
	}
	for i, exp := range expectedQA {
		if route.QA[i] != exp {
			t.Errorf("expected route.QA[%d]=%q, got %q", i, exp, route.QA[i])
		}
	}

	// Verify AllResourceIDs includes expanded QA resources
	hasPlaywright := false
	for _, id := range route.AllResourceIDs {
		if id == "playwright" {
			hasPlaywright = true
			break
		}
	}
	if !hasPlaywright {
		t.Errorf("expected route.AllResourceIDs to contain 'playwright', got %v", route.AllResourceIDs)
	}

	// 2. Verify nil guard: empty QA produces non-nil slice
	emptyRoute, ok := graph.ResolveCapabilityRoute("test-empty-qa", []string{"test"})
	if !ok {
		t.Fatalf("expected empty QA capability route to resolve")
	}
	if emptyRoute.QA == nil {
		t.Errorf("expected emptyRoute.QA to be non-nil empty slice, got nil")
	}
}


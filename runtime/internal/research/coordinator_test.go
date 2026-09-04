package research

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

func setupTestCatalogAndGraph(t *testing.T) (*resources.ResourceCatalog, *resources.DesignResourceGraph) {
	t.Helper()

	candidates := []string{
		filepath.Join("..", "..", "..", "registries"),
		filepath.Join("..", "..", "registries"),
		filepath.Join("registries"),
		`C:\projects\orchestra-workflow\registries`,
	}

	var regDir string
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			regDir = c
			break
		}
	}
	if regDir == "" {
		t.Fatalf("Could not locate registries directory in candidates: %v", candidates)
	}

	catPath := filepath.Join(regDir, "resources.json")
	graphPath := filepath.Join(regDir, "design-resource-graph.json")

	cat, err := resources.LoadResourceCatalog(catPath)
	if err != nil {
		t.Fatalf("Failed to load catalog from %s: %v", catPath, err)
	}

	graph, err := resources.LoadDesignGraph(graphPath)
	if err != nil {
		t.Fatalf("Failed to load design graph from %s: %v", graphPath, err)
	}

	return cat, graph
}

func TestResearchCoordinator_HighVisual_MultiSourceSelection(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:                "task-premium-visual-01",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"creative", "landing-page", "award-winning", "layout"},
	}

	routes := graph.ResolveCapabilities([]string{"premium-website"})
	if len(routes) == 0 {
		t.Fatalf("Expected route for premium-website")
	}
	route := &routes[0]

	opts := SelectionOptions{
		MinSources:       2,
		MaxSources:       4,
		IncludeBookmarks: true,
	}

	selected, err := coord.SelectReferences(context.Background(), task, route, opts)
	if err != nil {
		t.Fatalf("SelectReferences failed: %v", err)
	}

	if len(selected) < 2 {
		t.Errorf("Expected at least 2 selected sources for high-visual task, got %d", len(selected))
	}

	// Verify diverse families selected
	families := make(map[string]bool)
	for _, s := range selected {
		families[s.Family] = true
	}
	if len(families) < 2 {
		t.Errorf("Expected at least 2 distinct archetype families for diversity, got %d", len(families))
	}

	// Verify top score is reasonable
	if selected[0].Score <= 0.5 {
		t.Errorf("Expected top scored source to have composite score > 0.5, got %.2f", selected[0].Score)
	}
}

func TestResearchCoordinator_ScoringAndDiversity(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:                "task-scoring-test",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"typography", "aesthetic-history"},
	}

	opts := SelectionOptions{
		MinSources:       2,
		MaxSources:       4,
		IncludeBookmarks: true,
	}

	selected, err := coord.SelectReferences(context.Background(), task, nil, opts)
	if err != nil {
		t.Fatalf("SelectReferences failed: %v", err)
	}

	// Ensure ACTIVE sources score higher than BOOKMARK when domains are equal
	for _, s := range selected {
		if s.Resource.Status == "ACTIVE" && s.Score < 0.6 {
			t.Errorf("Active source %s scored unexpectedly low: %.2f", s.Resource.ID, s.Score)
		}
	}
}

func TestResearchCoordinator_NonVisualBypass(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:             "task-backend-bugfix",
		Type:           "BUGFIX",
		RequiresVisual: false,
	}

	req := &ResearchRequest{
		Task:    task,
		Options: SelectionOptions{MinSources: 2, MaxSources: 4},
	}

	res, err := coord.Coordinate(context.Background(), req)
	if err != nil {
		t.Fatalf("Coordinate failed on non-visual task: %v", err)
	}

	if res.TotalSourcesQueried != 0 {
		t.Errorf("Expected 0 sources queried for non-visual task, got %d", res.TotalSourcesQueried)
	}
	if len(res.SelectedSources) != 0 {
		t.Errorf("Expected 0 selected sources for non-visual task, got %d", len(res.SelectedSources))
	}
}

func TestResearchCoordinator_ReferenceLogGeneration(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	tempDir, err := os.MkdirTemp("", "orchestra_research_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	task := &classifier.Task{
		ID:                "task-reflog-test",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"creative", "showcase"},
	}

	routes := graph.ResolveCapabilities([]string{"premium-website"})
	if len(routes) == 0 {
		t.Fatalf("Expected routes for premium-website")
	}

	req := &ResearchRequest{
		Task:             task,
		Route:            &routes[0],
		Options:          SelectionOptions{MinSources: 2, MaxSources: 4, IncludeBookmarks: true},
		ProjectOutputDir: tempDir,
	}

	res, err := coord.Coordinate(context.Background(), req)
	if err != nil {
		t.Fatalf("Coordinate failed: %v", err)
	}

	expectedLogPath := filepath.Join(tempDir, "reference-log.md")
	if res.ReferenceLogPath != expectedLogPath {
		t.Errorf("Expected ReferenceLogPath %s, got %s", expectedLogPath, res.ReferenceLogPath)
	}

	content, err := os.ReadFile(expectedLogPath)
	if err != nil {
		t.Fatalf("Failed to read generated reference log from %s: %v", expectedLogPath, err)
	}
	logStr := string(content)

	// Validate required sections
	requiredSections := []string{
		"# Visual Research Reference Log",
		"## 1. Executive Direction & Creative Hypothesis",
		"## 2. Curated Reference Sources & Algorithmic Scoring",
		"## 3. Direct Citation Anchors",
		"## 4. Discovered Visual Motifs & Layout Lessons",
		"## 5. Synthesized Color Palette Tokens",
		"## 6. Typography Pairing & Hierarchy Tokens",
		"## 7. Kinetic & Interaction Dynamics",
		"## 8. Negative Guardrails & Anti-Pattern Checklist",
	}

	for _, sec := range requiredSections {
		if !strings.Contains(logStr, sec) {
			t.Errorf("Missing required section in reference-log.md: %s", sec)
		}
	}

	// Validate anti-pattern check: Inter must NOT be display font
	if strings.Contains(logStr, "| `display` | **Inter**") {
		t.Errorf("Anti-pattern violation detected: Inter font used as display font!")
	}

	// Validate minimum color tokens present
	if !strings.Contains(logStr, "--color-bg-base") || !strings.Contains(logStr, "--color-accent-primary") {
		t.Errorf("Missing essential color palette tokens in reference-log.md")
	}
}

func TestResearchCoordinator_QuarantineEnforcement(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	quarantinedPath := filepath.Join("C:", "Users", "mockuser", ".gemini", "config", "skills_library", "reference-log.md")
	res := &ResearchResult{
		TaskID:      "test-quarantine",
		Archetype:   "creative_showcase",
		QualityBar:  "premium",
		GeneratedAt: time.Now().UTC(),
		SynthesizedPalette: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14"},
			{Role: "--color-surface-elevated", Hex: "#141A24"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B"},
			{Role: "--color-text-headline", Hex: "#F0F4F8"},
		},
		SelectedSources: []*ScoredReference{
			{Resource: &resources.Resource{ID: "awwwards", Name: "Awwwards"}},
			{Resource: &resources.Resource{ID: "jiro-design", Name: "Jiro"}},
		},
	}

	err := coord.GenerateReferenceLog(context.Background(), res, quarantinedPath)
	if err == nil {
		t.Fatalf("Expected GenerateReferenceLog to reject quarantined skills_library path, but got nil")
	}
}

func TestResearchCoordinator_OfflineBenchmarkMode(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)
	coord.OfflineMode = true

	task := &classifier.Task{
		ID:                "task-offline-benchmark",
		Type:              "BENCHMARK",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"agriculture", "ttb-agro", "precision-farming"},
	}

	routes := graph.ResolveCapabilities([]string{"saas-dashboard", "operator-hud"})
	var route *resources.CapabilityRoute
	if len(routes) > 0 {
		route = &routes[0]
	}

	req := &ResearchRequest{
		Task:  task,
		Route: route,
		Options: SelectionOptions{
			MinSources:       2,
			MaxSources:       3,
			OfflineBenchmark: true,
		},
	}

	res, err := coord.Coordinate(context.Background(), req)
	if err != nil {
		t.Fatalf("Offline benchmark coordinate failed: %v", err)
	}

	if len(res.SelectedSources) < 2 {
		t.Errorf("Expected at least 2 sources in offline benchmark, got %d", len(res.SelectedSources))
	}

	if len(res.SynthesizedPalette) < 4 {
		t.Errorf("Expected at least 4 synthesized palette tokens in offline benchmark, got %d", len(res.SynthesizedPalette))
	}

	if len(res.SynthesizedTypography) < 3 {
		t.Errorf("Expected at least 3 typography tokens in offline benchmark, got %d", len(res.SynthesizedTypography))
	}
}



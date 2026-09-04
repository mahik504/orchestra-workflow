package research

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

// ---------------------------------------------------------------------------
// 1. ATTEMPTED QUARANTINE EVASION TESTS
// ---------------------------------------------------------------------------

func TestAdversarial_Quarantine_DirectBannedPaths(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	bannedSubstrings := []string{
		`C:\Users\mockuser\.gemini\config\skills_library\output`,
		`C:/Users/mockuser/.gemini/config/skills_library/output`,
		`C:\Users\mockuser\.gemini\config\SKILLS_LIBRARY\reference-log.md`,
		`C:\projects\orchestra-workflow\skills-library\log.md`,
		`C:\projects\curated_catalog\quarantine\log.md`,
		`./relative/path/skills_library/log.md`,
		`skills_library`,
	}

	dummyResult := &ResearchResult{
		TaskID:      "task-quarantine-test",
		Archetype:   "creative_showcase",
		QualityBar:  "premium",
		GeneratedAt: time.Now().UTC(),
		SelectedSources: []*ScoredReference{
			{Resource: &resources.Resource{ID: "awwwards", Name: "Awwwards"}},
			{Resource: &resources.Resource{ID: "jiro-design", Name: "Jiro"}},
		},
		SynthesizedPalette: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14"},
			{Role: "--color-surface-elevated", Hex: "#141A24"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B"},
			{Role: "--color-text-headline", Hex: "#F0F4F8"},
		},
	}

	for _, bannedPath := range bannedSubstrings {
		t.Run("Path_"+filepath.Base(bannedPath), func(t *testing.T) {
			// Direct call to GenerateReferenceLog
			err := coord.GenerateReferenceLog(context.Background(), dummyResult, bannedPath)
			if err == nil {
				t.Fatalf("Expected GenerateReferenceLog to reject banned path %q, but succeeded", bannedPath)
			}
			if !errors.Is(err, resources.ErrQuarantinedPath) && !strings.Contains(err.Error(), "quarantine") {
				t.Errorf("Expected ErrQuarantinedPath error, got: %v", err)
			}

			// Call via Coordinate with ProjectOutputDir
			req := &ResearchRequest{
				Task: &classifier.Task{
					ID:             "task-quarantine-coord",
					Type:           "DESIGN",
					RequiresVisual: true,
				},
				Options: SelectionOptions{
					MinSources: 2,
					MaxSources: 4,
				},
				ProjectOutputDir: bannedPath,
			}
			_, coordErr := coord.Coordinate(context.Background(), req)
			if coordErr == nil {
				t.Fatalf("Expected Coordinate to reject banned ProjectOutputDir %q, but succeeded", bannedPath)
			}
			if !errors.Is(coordErr, resources.ErrQuarantinedPath) && !strings.Contains(coordErr.Error(), "quarantine") {
				t.Errorf("Expected quarantine error from Coordinate, got: %v", coordErr)
			}
		})
	}
}

func TestAdversarial_Quarantine_8dot3Aliases(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	aliases := []string{
		`C:\Users\mockuser\.gemini\config\SKILLS~1\reference-log.md`,
		`C:\Users\mockuser\.gemini\config\skills~1\reference-log.md`,
		`C:\projects\CURATE~1\reference-log.md`,
		`C:\projects\curate~1\reference-log.md`,
		`./SKILLS~1/reference-log.md`,
	}

	dummyResult := &ResearchResult{
		TaskID:      "task-8dot3-test",
		Archetype:   "creative_showcase",
		QualityBar:  "premium",
		GeneratedAt: time.Now().UTC(),
		SelectedSources: []*ScoredReference{
			{Resource: &resources.Resource{ID: "awwwards", Name: "Awwwards"}},
			{Resource: &resources.Resource{ID: "jiro-design", Name: "Jiro"}},
		},
		SynthesizedPalette: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14"},
			{Role: "--color-surface-elevated", Hex: "#141A24"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B"},
			{Role: "--color-text-headline", Hex: "#F0F4F8"},
		},
	}

	for _, alias := range aliases {
		t.Run("Alias_"+filepath.Base(alias), func(t *testing.T) {
			err := coord.GenerateReferenceLog(context.Background(), dummyResult, alias)
			if err == nil {
				t.Fatalf("Expected 8.3 alias %q to be strictly rejected, but succeeded", alias)
			}
			if !errors.Is(err, resources.ErrQuarantinedPath) && !strings.Contains(err.Error(), "quarantine") {
				t.Errorf("Expected ErrQuarantinedPath error, got: %v", err)
			}
		})
	}
}

func TestAdversarial_Quarantine_NTFSJunctionResolution(t *testing.T) {
	// Empirically create a real NTFS junction pointing to a quarantined directory structure
	baseTemp, err := os.MkdirTemp("", "orch_quarantine_adv_*")
	if err != nil {
		t.Fatalf("Failed to create base temp dir: %v", err)
	}
	defer os.RemoveAll(baseTemp)

	// Target directory with quarantined name
	quarantinedTarget := filepath.Join(baseTemp, "skills_library")
	if err := os.MkdirAll(quarantinedTarget, 0755); err != nil {
		t.Fatalf("Failed to create target quarantined dir: %v", err)
	}

	// Junction link directory with innocuous name
	innocuousJunction := filepath.Join(baseTemp, "innocuous_design_cache")

	// Create NTFS junction using Windows mklink /J
	cmd := exec.Command("cmd", "/c", "mklink", "/J", innocuousJunction, quarantinedTarget)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("NTFS Junction creation skipped (environment may not support mklink): %v, output: %s", err, string(out))
		return
	}
	defer func() {
		_ = exec.Command("cmd", "/c", "rmdir", innocuousJunction).Run()
	}()

	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	dummyResult := &ResearchResult{
		TaskID:      "task-junction-test",
		Archetype:   "creative_showcase",
		QualityBar:  "premium",
		GeneratedAt: time.Now().UTC(),
		SelectedSources: []*ScoredReference{
			{Resource: &resources.Resource{ID: "awwwards", Name: "Awwwards"}},
			{Resource: &resources.Resource{ID: "jiro-design", Name: "Jiro"}},
		},
		SynthesizedPalette: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14"},
			{Role: "--color-surface-elevated", Hex: "#141A24"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B"},
			{Role: "--color-text-headline", Hex: "#F0F4F8"},
		},
	}

	// Test 1: File directly inside junction link
	junctionFile := filepath.Join(innocuousJunction, "reference-log.md")
	err = coord.GenerateReferenceLog(context.Background(), dummyResult, junctionFile)
	if err == nil {
		t.Fatalf("CRITICAL SECURITY FAILURE: Junction pointing to skills_library was NOT rejected! File written to: %s", junctionFile)
	}
	if !errors.Is(err, resources.ErrQuarantinedPath) && !strings.Contains(err.Error(), "quarantine") {
		t.Errorf("Expected ErrQuarantinedPath error on junction, got: %v", err)
	}

	// Test 2: Nested path inside junction where leaf directory does not exist yet
	nestedJunctionFile := filepath.Join(innocuousJunction, "subfolder", "nested-log.md")
	nestedErr := coord.GenerateReferenceLog(context.Background(), dummyResult, nestedJunctionFile)
	if nestedErr == nil {
		t.Fatalf("CRITICAL SECURITY FAILURE: Nested path through junction was NOT rejected!")
	}
	if !errors.Is(nestedErr, resources.ErrQuarantinedPath) && !strings.Contains(nestedErr.Error(), "quarantine") {
		t.Errorf("Expected ErrQuarantinedPath error on nested junction, got: %v", nestedErr)
	}
}

func TestAdversarial_Quarantine_CandidateInjection(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:                "task-malicious-candidates",
		Type:              "DESIGN",
		RequiresVisual:    true,
		SuggestedResources: []string{"skills_library", "skills~1", "curated_catalog/quarantine"},
	}

	route := &resources.CapabilityRoute{
		CapabilityID:       "premium-website",
		DiscoveryResources: []string{"skills_library", "skills~1", "awwwards"},
		DiscoveryDomains:   []string{"visual_research"},
	}

	opts := SelectionOptions{
		MinSources: 2,
		MaxSources: 4,
	}

	selected, err := coord.SelectReferences(context.Background(), task, route, opts)
	if err != nil {
		t.Fatalf("SelectReferences unexpected error: %v", err)
	}

	for _, s := range selected {
		lowerID := strings.ToLower(s.Resource.ID)
		if strings.Contains(lowerID, "skills_library") || strings.Contains(lowerID, "skills~1") || strings.Contains(lowerID, "quarantine") {
			t.Fatalf("CRITICAL SECURITY LEAK: Quarantined candidate leaked into selected sources: %s", s.Resource.ID)
		}
	}

	if len(selected) < 2 {
		t.Errorf("Expected at least 2 clean sources to be selected despite malicious injection, got %d", len(selected))
	}
}

// ---------------------------------------------------------------------------
// 2. OFFLINE RESILIENCE & ZERO NETWORK TESTS
// ---------------------------------------------------------------------------

func TestAdversarial_OfflineResilience_ZeroNetworkSimulation(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)
	coord.OfflineMode = true

	// Strict deadline of 50ms to guarantee zero blocking network I/O
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	task := &classifier.Task{
		ID:                "task-zero-network",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"landing-page", "creative", "awwwards"},
	}

	req := &ResearchRequest{
		Task: task,
		Options: SelectionOptions{
			MinSources:       2,
			MaxSources:       4,
			OfflineBenchmark: true,
		},
	}

	start := time.Now()
	res, err := coord.Coordinate(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Offline coordinate failed: %v", err)
	}

	if elapsed > 30*time.Millisecond {
		t.Errorf("Offline coordinate took too long (%v > 30ms); suspect blocking network call", elapsed)
	}

	if len(res.SelectedSources) < 2 {
		t.Errorf("Expected >= 2 sources in offline mode, got %d", len(res.SelectedSources))
	}

	if len(res.Findings) < 2 {
		t.Errorf("Expected >= 2 findings in offline mode, got %d", len(res.Findings))
	}

	// Verify all findings originate from verified offline fixtures
	for _, f := range res.Findings {
		if len(f.KeyTakeaways) == 0 {
			t.Errorf("Finding for %s has empty KeyTakeaways in offline mode", f.SourceID)
		}
		if len(f.Citations) == 0 {
			t.Errorf("Finding for %s has empty Citations in offline mode", f.SourceID)
		}
	}
}

func TestAdversarial_OfflineResilience_CuratedFixturesIntegrity(t *testing.T) {
	requiredFixtures := []string{
		"awwwards",
		"jiro-design",
		"cari-institute",
		"awesome-design-md",
		"godly-design",
		"refero-design",
	}

	for _, fixtureID := range requiredFixtures {
		f, ok := CuratedSourceFixtures[fixtureID]
		if !ok {
			t.Fatalf("CuratedSourceFixtures missing mandatory offline fixture: %s", fixtureID)
		}

		if f.SourceID != fixtureID {
			t.Errorf("Fixture ID mismatch: %s != %s", f.SourceID, fixtureID)
		}
		if f.SourceName == "" {
			t.Errorf("Fixture %s missing SourceName", fixtureID)
		}
		if !strings.HasPrefix(f.SourceURL, "https://") {
			t.Errorf("Fixture %s SourceURL must be valid HTTPS URL, got: %s", fixtureID, f.SourceURL)
		}
		if len(f.KeyTakeaways) < 2 {
			t.Errorf("Fixture %s has insufficient KeyTakeaways (%d < 2)", fixtureID, len(f.KeyTakeaways))
		}
		if len(f.Citations) == 0 {
			t.Errorf("Fixture %s has 0 Citations", fixtureID)
		}
	}
}

func TestAdversarial_OfflineResilience_RichVisualProfiles(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:                "task-rich-visual",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"editorial", "high-aesthetic"},
	}

	routes := graph.ResolveCapabilities([]string{"premium-website"})
	if len(routes) == 0 {
		t.Fatalf("Expected routes for premium-website")
	}

	req := &ResearchRequest{
		Task:  task,
		Route: &routes[0],
		Options: SelectionOptions{
			MinSources: 2,
			MaxSources: 4,
		},
	}

	res, err := coord.Coordinate(context.Background(), req)
	if err != nil {
		t.Fatalf("Coordinate failed: %v", err)
	}

	// 1. Color Palette Tokens Verification
	if len(res.SynthesizedPalette) < 6 {
		t.Errorf("Expected at least 6 synthesized palette tokens, got %d", len(res.SynthesizedPalette))
	}

	paletteRoles := make(map[string]PaletteToken)
	for _, p := range res.SynthesizedPalette {
		paletteRoles[p.Role] = p
	}

	requiredRoles := []string{
		"--color-bg-base",
		"--color-surface-elevated",
		"--color-accent-primary",
		"--color-text-headline",
		"--color-text-muted",
		"--color-border-subtle",
	}
	for _, role := range requiredRoles {
		token, exists := paletteRoles[role]
		if !exists {
			t.Errorf("Missing required color token role: %s", role)
		}
		if token.Hex == "" || !strings.HasPrefix(token.Hex, "#") {
			t.Errorf("Token %s has invalid hex color: %s", role, token.Hex)
		}
		if token.Contrast == "" {
			t.Errorf("Token %s missing contrast ratio specification", role)
		}
	}

	// Verify dark matte base is #0B0E14 (not pure #000000)
	if paletteRoles["--color-bg-base"].Hex == "#000000" {
		t.Errorf("Anti-pattern violation: --color-bg-base must NOT be pure #000000")
	}

	// 2. Typography Triad Hierarchy Verification
	if len(res.SynthesizedTypography) < 4 {
		t.Errorf("Expected at least 4 typography hierarchy tokens, got %d", len(res.SynthesizedTypography))
	}

	typoRoles := make(map[string]TypographyToken)
	for _, typo := range res.SynthesizedTypography {
		typoRoles[typo.Role] = typo
	}

	for _, role := range []string{"display", "headline", "body", "mono"} {
		typo, ok := typoRoles[role]
		if !ok {
			t.Errorf("Missing typography hierarchy role: %s", role)
			continue
		}
		if typo.FontFamily == "" {
			t.Errorf("Typography %s has empty FontFamily", role)
		}
		if typo.SizeClamp == "" {
			t.Errorf("Typography %s missing SizeClamp formula", role)
		}
	}

	// Check that display font is NOT Inter or Space Grotesk
	displayFont := strings.ToLower(typoRoles["display"].FontFamily)
	if strings.Contains(displayFont, "inter") || strings.Contains(displayFont, "space grotesk") {
		t.Errorf("Anti-pattern violation: display font cannot be Inter or Space Grotesk, got: %s", typoRoles["display"].FontFamily)
	}

	// 3. Interaction Dynamics Verification
	if res.InteractionRules.MotionCurve == "" {
		t.Errorf("Missing MotionCurve in InteractionRules")
	}
	if !strings.Contains(res.InteractionRules.ActivePress, "scale(0.98)") {
		t.Errorf("Expected active press feedback to specify scale(0.98), got: %s", res.InteractionRules.ActivePress)
	}
	if len(res.InteractionRules.ProhibitedProps) == 0 {
		t.Errorf("ProhibitedProps must list layout-triggering properties")
	}

	// 4. Anti-Patterns Checklist Verification
	if len(res.BannedAntiPatterns) < 5 {
		t.Errorf("Expected at least 5 banned anti-patterns, got %d", len(res.BannedAntiPatterns))
	}
}

// ---------------------------------------------------------------------------
// 3. DIVERSITY ENFORCEMENT & ARCHETYPE FAMILIES TESTS
// ---------------------------------------------------------------------------

func TestAdversarial_Diversity_SingleSourceSuggestionRejection(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	// User suggests only a single source for a high-visual task
	task := &classifier.Task{
		ID:                "task-single-source",
		Type:              "DESIGN",
		RequiresVisual:    true,
		SuggestedResources: []string{"awwwards"},
	}

	route := &resources.CapabilityRoute{
		CapabilityID:       "premium-website",
		PrimaryArchetype:   "creative_showcase",
		DiscoveryResources: []string{"awwwards"}, // only 1 resource provided
	}

	opts := SelectionOptions{
		MinSources: 2,
		MaxSources: 4,
	}

	selected, err := coord.SelectReferences(context.Background(), task, route, opts)
	if err != nil {
		t.Fatalf("SelectReferences failed: %v", err)
	}

	// Coordinator MUST enforce >= 2 distinct sources
	if len(selected) < 2 {
		t.Fatalf("Diversity enforcement failure: expected >= 2 sources for high-visual task, got %d", len(selected))
	}

	if selected[0].Resource.ID == selected[1].Resource.ID {
		t.Fatalf("Diversity enforcement failure: duplicate sources selected: %s and %s", selected[0].Resource.ID, selected[1].Resource.ID)
	}

	// Verify distinct families
	families := make(map[string]bool)
	for _, s := range selected {
		families[s.Family] = true
	}
	if len(families) < 2 {
		t.Errorf("Diversity enforcement failure: expected >= 2 distinct families, got %d: %v", len(families), families)
	}
}

func TestAdversarial_Diversity_CrossFamilyDistribution(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	task := &classifier.Task{
		ID:                "task-all-families",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"creative", "typography", "taxonomy", "shaders"},
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

	families := make(map[string]int)
	for _, s := range selected {
		families[s.Family]++
		if s.Role == "" {
			t.Errorf("Selected source %s has empty assigned Role", s.Resource.ID)
		}
	}

	if len(families) < 2 {
		t.Errorf("Expected multi-family representation, found only %d families: %v", len(families), families)
	}

	// Verify role assignment matches archetype family
	for _, s := range selected {
		switch s.Family {
		case FamilyAestheticBenchmark:
			if s.Role != "Aesthetic Benchmark & Creative Direction" {
				t.Errorf("Role mismatch for family %s: got %q", s.Family, s.Role)
			}
		case FamilyMovementTaxonomy:
			if s.Role != "Aesthetic Taxonomy & Design Contracts" {
				t.Errorf("Role mismatch for family %s: got %q", s.Family, s.Role)
			}
		case FamilySpecialistEcho:
			if s.Role != "Specialist Domain Echo & Pattern Reference" {
				t.Errorf("Role mismatch for family %s: got %q", s.Family, s.Role)
			}
		}
	}
}

func TestAdversarial_Diversity_MinSourceViolationGate(t *testing.T) {
	cat, graph := setupTestCatalogAndGraph(t)
	coord := NewResearchCoordinator(cat, graph)

	tempDir, err := os.MkdirTemp("", "orch_gate_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// High-visual result with ONLY 1 source (violates policy)
	singleSourceResult := &ResearchResult{
		TaskID:      "task-gate-violation",
		Archetype:   "creative_showcase",
		QualityBar:  "premium",
		GeneratedAt: time.Now().UTC(),
		SelectedSources: []*ScoredReference{
			{Resource: &resources.Resource{ID: "awwwards", Name: "Awwwards"}},
		},
		SynthesizedPalette: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14"},
			{Role: "--color-surface-elevated", Hex: "#141A24"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B"},
			{Role: "--color-text-headline", Hex: "#F0F4F8"},
		},
	}

	outPath := filepath.Join(tempDir, "reference-log.md")
	err = coord.GenerateReferenceLog(context.Background(), singleSourceResult, outPath)
	if err == nil {
		t.Fatalf("Expected GenerateReferenceLog to REJECT single-source premium result, but succeeded")
	}
	if !strings.Contains(err.Error(), "mandate at least 2 curated sources") {
		t.Errorf("Expected 2-source gate rejection message, got: %v", err)
	}
}

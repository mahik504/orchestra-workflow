package router

import (
	"strings"
	"testing"

	"github.com/user/orchestra-v3/internal/classifier"
)

// ---------------------------------------------------------------------------
// 1. ALL 10 DESIGN ARCHETYPES TOKEN WEIGHT SCALING & BOUNDS
// ---------------------------------------------------------------------------

func TestAdversarial_Router_All10DesignArchetypesTokenScaling(t *testing.T) {
	cat, graph := setupLiveCatalogAndGraph(t)
	r := NewRouterWithGraph(nil, cat, graph)

	type ArchetypeTestDef struct {
		ID                    string
		Name                  string
		RequiresVisual        bool
		RequiresSec           bool
		Keywords              []string
		MinTokens             float64
		MaxTokens             float64
		ExpectedDirectivesMin int
	}

	testCases := []ArchetypeTestDef{
		{
			ID:                    "standard-baseline-bugfix",
			Name:                  "Standard Non-Visual Bugfix",
			RequiresVisual:        false,
			RequiresSec:           false,
			Keywords:              []string{"bugfix", "backend", "database-migration"},
			MinTokens:             1500,
			MaxTokens:             3000,
			ExpectedDirectivesMin: 1, // superpowers-planning
		},
		{
			ID:                    "security-audit",
			Name:                  "Automated Security Audit & Hardening",
			RequiresVisual:        false,
			RequiresSec:           true,
			Keywords:              []string{"security-audit", "pentest", "vulnerability-scan"},
			MinTokens:             2500,
			MaxTokens:             15000,
			ExpectedDirectivesMin: 2, // superpowers + semgrep
		},
		{
			ID:                    "premium-website",
			Name:                  "Premium Creative Website",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"premium-website", "landing-page", "award-winning"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3, // superpowers + taste + impeccable + route
		},
		{
			ID:                    "3d-portfolio",
			Name:                  "Interactive 3D WebGL Portfolio",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"3d-portfolio", "webgl", "threejs", "r3f"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "operator-hud",
			Name:                  "Operator HUD & Mission Control",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"operator-hud", "hud", "telemetry", "dark-crimson"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "b2b-portal",
			Name:                  "B2B Enterprise Portal & SaaS",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"b2b-portal", "b2b", "portal", "enterprise"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "academic-reader",
			Name:                  "Academic Reader & Research Paper Viewer",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"academic-reader", "research-paper", "reading"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "micro-interactions",
			Name:                  "Sensory Micro-Interactions & Gestures",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"micro-interactions", "spring-physics", "tactile"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "physics-canvas",
			Name:                  "Interactive Physics Canvas & 2D Simulation",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"physics-canvas", "canvas-game", "simulation"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "saas-dashboard",
			Name:                  "Modern SaaS Analytics Dashboard",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"saas-dashboard", "analytics", "kpi-charts"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
		{
			ID:                    "mobile-app",
			Name:                  "Native-Feeling Mobile Web Experience",
			RequiresVisual:        true,
			RequiresSec:           false,
			Keywords:              []string{"mobile-app", "expo", "touch-native"},
			MinTokens:             8000,
			MaxTokens:             50000,
			ExpectedDirectivesMin: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.ID, func(t *testing.T) {
			task := &classifier.Task{
				ID:                "test-" + tc.ID,
				Type:              "FEATURE",
				RequiresVisual:    tc.RequiresVisual,
				RequiresSecurity:  tc.RequiresSec,
				ExtractedKeywords: tc.Keywords,
			}

			plan := r.Compose(task)

			// Validate token bounds
			if plan.EstimatedTokenCost < tc.MinTokens {
				t.Errorf("Archetype %s token cost %.0f is below minimum bound %.0f", tc.ID, plan.EstimatedTokenCost, tc.MinTokens)
			}
			if plan.EstimatedTokenCost > tc.MaxTokens {
				t.Errorf("Archetype %s token cost %.0f exceeds maximum bound %.0f", tc.ID, plan.EstimatedTokenCost, tc.MaxTokens)
			}

			// Validate minimum directive count
			if len(plan.ExecutionDirectives) < tc.ExpectedDirectivesMin {
				t.Errorf("Archetype %s has insufficient execution directives (%d < %d)", tc.ID, len(plan.ExecutionDirectives), tc.ExpectedDirectivesMin)
			}

			// Validate 8-stage pipeline stage string
			expectedStage := "Discover -> Classify -> Research -> Synthesize -> Design System -> Implement -> Visual QA -> Iterate"
			if plan.PipelineStage != expectedStage {
				t.Errorf("Archetype %s returned incorrect pipeline stage string %q", tc.ID, plan.PipelineStage)
			}

			// Validate human gate rules
			if tc.RequiresVisual && !plan.RequiresHumanGate {
				t.Errorf("Archetype %s requires visual design but human gate was not triggered", tc.ID)
			}
			if !tc.RequiresVisual && plan.RequiresHumanGate {
				t.Errorf("Archetype %s does not require visual design but human gate was triggered", tc.ID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 2. TOKEN DEDUPLICATION & IDEMPOTENCY
// ---------------------------------------------------------------------------

func TestAdversarial_Router_TokenDeduplication(t *testing.T) {
	cat, graph := setupLiveCatalogAndGraph(t)
	r := NewRouterWithGraph(nil, cat, graph)

	// Task with heavily redundant tags resolving to overlapping resources
	redundantTask := &classifier.Task{
		ID:             "task-redundant-tags",
		Type:           "DESIGN",
		RequiresVisual: true,
		ExtractedKeywords: []string{
			"premium-website",
			"landing-page",
			"agency",
			"award-winning",
			"portfolio-hero",
			"high-aesthetic",
			"visual_research",
		},
		SuggestedResources: []string{
			"awwwards",
			"awwwards",
			"taste-design",
			"impeccable",
			"gsap",
			"gsap",
		},
	}

	cleanTask := &classifier.Task{
		ID:                "task-clean-tags",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"premium-website"},
	}

	planRedundant := r.Compose(redundantTask)
	planClean := r.Compose(cleanTask)

	// Since both resolve the same capability resources and duplicates are deduplicated,
	// the token costs should be identical or extremely close (no duplicate resource addition).
	if planRedundant.EstimatedTokenCost != planClean.EstimatedTokenCost {
		t.Errorf("Deduplication failure: redundant tags inflated token cost (clean=%.0f vs redundant=%.0f)",
			planClean.EstimatedTokenCost, planRedundant.EstimatedTokenCost)
	}

	// Verify idempotency: run multiple times, cost must remain strictly identical
	for i := 0; i < 5; i++ {
		repeatPlan := r.Compose(redundantTask)
		if repeatPlan.EstimatedTokenCost != planRedundant.EstimatedTokenCost {
			t.Fatalf("Idempotency failure on run %d: got %.0f, expected %.0f", i, repeatPlan.EstimatedTokenCost, planRedundant.EstimatedTokenCost)
		}
	}
}

// ---------------------------------------------------------------------------
// 3. MONOTONIC SCALING ACROSS COMPLEXITY TIERS
// ---------------------------------------------------------------------------

func TestAdversarial_Router_MonotonicComplexityScaling(t *testing.T) {
	cat, graph := setupLiveCatalogAndGraph(t)
	r := NewRouterWithGraph(nil, cat, graph)

	// Tier 1: Minimal Bugfix
	t1 := &classifier.Task{
		ID:             "tier-1-bugfix",
		Type:           "BUGFIX",
		RequiresVisual: false,
	}
	p1 := r.Compose(t1)

	// Tier 2: Security Audit
	t2 := &classifier.Task{
		ID:               "tier-2-security",
		Type:             "AUDIT",
		RequiresSecurity: true,
	}
	p2 := r.Compose(t2)

	// Tier 3: Standard Visual Task
	t3 := &classifier.Task{
		ID:                "tier-3-visual",
		Type:              "FEATURE",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"b2b-portal"},
	}
	p3 := r.Compose(t3)

	// Tier 4: High-Complexity Spatial WebGL Showcase
	t4 := &classifier.Task{
		ID:                "tier-4-spatial",
		Type:              "DESIGN",
		RequiresVisual:    true,
		ExtractedKeywords: []string{"3d-portfolio", "webgl", "spatial-ui"},
	}
	p4 := r.Compose(t4)

	// Tier 5: Composite Mission Control with Telemetry + Security Audit
	t5 := &classifier.Task{
		ID:                "tier-5-composite",
		Type:              "DESIGN",
		RequiresVisual:    true,
		RequiresSecurity:  true,
		ExtractedKeywords: []string{"operator-hud", "telemetry", "terminal-ui", "security-audit"},
	}
	p5 := r.Compose(t5)

	t.Logf("Empirical Token Costs: Tier1=%.0f, Tier2=%.0f, Tier3=%.0f, Tier4=%.0f, Tier5=%.0f",
		p1.EstimatedTokenCost, p2.EstimatedTokenCost, p3.EstimatedTokenCost, p4.EstimatedTokenCost, p5.EstimatedTokenCost)

	// Assert strict monotonic progression: Tier 1 < Tier 2 < Tier 3 <= Tier 4 < Tier 5
	if p1.EstimatedTokenCost >= p2.EstimatedTokenCost {
		t.Errorf("Scaling violation: Tier 1 (%.0f) >= Tier 2 (%.0f)", p1.EstimatedTokenCost, p2.EstimatedTokenCost)
	}
	if p2.EstimatedTokenCost >= p3.EstimatedTokenCost {
		t.Errorf("Scaling violation: Tier 2 (%.0f) >= Tier 3 (%.0f)", p2.EstimatedTokenCost, p3.EstimatedTokenCost)
	}
	if p3.EstimatedTokenCost > p4.EstimatedTokenCost {
		t.Errorf("Scaling violation: Tier 3 (%.0f) > Tier 4 (%.0f)", p3.EstimatedTokenCost, p4.EstimatedTokenCost)
	}
	if p4.EstimatedTokenCost >= p5.EstimatedTokenCost {
		t.Errorf("Scaling violation: Tier 4 (%.0f) >= Tier 5 (%.0f)", p4.EstimatedTokenCost, p5.EstimatedTokenCost)
	}
}

// ---------------------------------------------------------------------------
// 4. NIL CATALOG & GRAPH FALLBACK RESILIENCE
// ---------------------------------------------------------------------------

func TestAdversarial_Router_NilCatalogAndGraphFallback(t *testing.T) {
	// Standalone router with NO live registries loaded (pure fallback mode)
	r := NewRouterWithGraph(nil, nil, nil)

	// Standard task: baseline superpowers (1500)
	taskStandard := &classifier.Task{ID: "fallback-standard", Type: "BUGFIX"}
	planStandard := r.Compose(taskStandard)
	if planStandard.EstimatedTokenCost != 1500 {
		t.Errorf("Expected fallback standard token cost == 1500, got %.0f", planStandard.EstimatedTokenCost)
	}

	// Visual task: baseline (1500) + taste (2000) + impeccable (1800) = 5300
	taskVisual := &classifier.Task{ID: "fallback-visual", Type: "DESIGN", RequiresVisual: true}
	planVisual := r.Compose(taskVisual)
	expectedVisual := 1500.0 + 2000.0 + 1800.0
	if planVisual.EstimatedTokenCost != expectedVisual {
		t.Errorf("Expected fallback visual token cost == %.0f, got %.0f", expectedVisual, planVisual.EstimatedTokenCost)
	}

	// Visual + Security: 5300 + semgrep (1200) = 6500
	taskMixed := &classifier.Task{ID: "fallback-mixed", Type: "FEATURE", RequiresVisual: true, RequiresSecurity: true}
	planMixed := r.Compose(taskMixed)
	expectedMixed := expectedVisual + 1200.0
	if planMixed.EstimatedTokenCost != expectedMixed {
		t.Errorf("Expected fallback mixed token cost == %.0f, got %.0f", expectedMixed, planMixed.EstimatedTokenCost)
	}

	// Manifest generation without panicking
	manifest := planMixed.GenerateExecutionManifest()
	if !strings.Contains(manifest, "Orchestra V3 Capability Execution Manifest") {
		t.Errorf("Failed to generate valid execution manifest in fallback mode")
	}
}

package research

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

// Coordinator defines the interface for multi-source design research coordination
type Coordinator interface {
	SelectReferences(ctx context.Context, task *classifier.Task, route *resources.CapabilityRoute, opts SelectionOptions) ([]*ScoredReference, error)
	Coordinate(ctx context.Context, req *ResearchRequest) (*ResearchResult, error)
	GenerateReferenceLog(ctx context.Context, res *ResearchResult, outputPath string) error
}

// ResearchCoordinator coordinates research across curated visual sources
type ResearchCoordinator struct {
	Catalog     *resources.ResourceCatalog
	Graph       *resources.DesignResourceGraph
	OfflineMode bool
}

// NewResearchCoordinator instantiates a ResearchCoordinator with loaded registries
func NewResearchCoordinator(cat *resources.ResourceCatalog, graph *resources.DesignResourceGraph) *ResearchCoordinator {
	return &ResearchCoordinator{
		Catalog:     cat,
		Graph:       graph,
		OfflineMode: true, // Default to offline resilience using canonical registries & fixtures
	}
}

// SelectReferences performs candidate discovery, filtering, and multi-factor scoring
func (c *ResearchCoordinator) SelectReferences(
	ctx context.Context,
	task *classifier.Task,
	route *resources.CapabilityRoute,
	opts SelectionOptions,
) ([]*ScoredReference, error) {
	if opts.MinSources <= 0 {
		opts.MinSources = 2
	}
	if opts.MaxSources <= 0 {
		opts.MaxSources = 4
	}

	// 1. Gather Candidate Resource IDs
	candidateIDs := make(map[string]bool)
	primaryDomainMap := make(map[string]bool)

	if route != nil {
		for _, id := range route.DiscoveryResources {
			candidateIDs[id] = true
		}
		for _, dom := range route.DiscoveryDomains {
			primaryDomainMap[dom] = true
			if c.Graph != nil {
				for _, id := range c.Graph.GetDomainResources(dom) {
					candidateIDs[id] = true
				}
			}
		}
	}

	// Ensure standard visual research candidates are considered if pool is small
	if len(candidateIDs) < opts.MinSources && c.Graph != nil {
		for _, id := range c.Graph.GetDomainResources("visual_research") {
			candidateIDs[id] = true
		}
	}

	// Canonical default pool fallback if still empty
	if len(candidateIDs) == 0 {
		defaults := []string{"awwwards", "jiro-design", "cari-institute", "awesome-design-md"}
		for _, d := range defaults {
			candidateIDs[d] = true
		}
	}

	// 2. Filter Candidates and Resolve Resources
	var scoredCandidates []*ScoredReference
	selectedFamilies := make(map[string]bool)

	taskTags := make([]string, 0)
	if task != nil {
		taskTags = append(taskTags, task.ExtractedKeywords...)
		if route != nil {
			taskTags = append(taskTags, route.MatchedTags...)
		}
	}

	for id := range candidateIDs {
		// Verify quarantine on ID and any derived path
		if err := resources.CheckQuarantineBoundary(id); err != nil {
			continue
		}

		var res *resources.Resource
		if c.Catalog != nil {
			if r, found := c.Catalog.FindByID(id); found {
				res = r
			}
		}

		// Fallback minimal resource if catalog doesn't index it yet
		if res == nil {
			res = &resources.Resource{
				ID:             id,
				Name:           id,
				Representation: "research_source",
				Status:         "ACTIVE",
				RoutingTags:    []string{"visual-research"},
			}
		}

		// Check quarantine on canonical URL and source repository
		if res.CanonicalURL != "" {
			if err := resources.CheckQuarantineBoundary(res.CanonicalURL); err != nil {
				continue
			}
		}
		if res.SourceRepository != "" {
			if err := resources.CheckQuarantineBoundary(res.SourceRepository); err != nil {
				continue
			}
		}

		// Filter out rejected or non-viable status
		statusUpper := strings.ToUpper(res.Status)
		if statusUpper == "REJECTED" || statusUpper == "ARCHIVED" || statusUpper == "DEPRECATED" {
			continue
		}
		if statusUpper == "BOOKMARK" && !opts.IncludeBookmarks {
			continue
		}

		// 3. Multi-Factor Scoring Formula:
		// Score = 0.40 * S_domain + 0.25 * S_tag + 0.20 * S_status + 0.15 * S_diversity

		// S_domain
		sDomain := 0.20
		if route != nil {
			for _, dr := range route.DiscoveryResources {
				if dr == id {
					sDomain = 1.0
					break
				}
			}
		}
		if sDomain < 1.0 {
			for dom := range primaryDomainMap {
				if c.Graph != nil {
					for _, r := range c.Graph.GetDomainResources(dom) {
						if r == id {
							sDomain = 0.85
							break
						}
					}
				}
				if sDomain >= 0.85 {
					break
				}
			}
		}
		if sDomain < 0.85 {
			relatedDomains := []string{"interaction_research", "editorial_research", "hud_research", "saas_research", "webgl_research"}
			for _, rdom := range relatedDomains {
				if c.Graph != nil {
					for _, r := range c.Graph.GetDomainResources(rdom) {
						if r == id {
							sDomain = 0.70
							break
						}
					}
				}
				if sDomain >= 0.70 {
					break
				}
			}
		}

		// S_tag: Jaccard-like tag overlap
		sTag := 0.3
		if len(taskTags) > 0 && len(res.RoutingTags) > 0 {
			matches := 0
			for _, tt := range taskTags {
				ttLower := strings.ToLower(tt)
				for _, rt := range res.RoutingTags {
					if strings.Contains(strings.ToLower(rt), ttLower) || strings.Contains(ttLower, strings.ToLower(rt)) {
						matches++
						break
					}
				}
			}
			sTag = float64(matches) / 2.0
			if sTag > 1.0 {
				sTag = 1.0
			}
			if sTag == 0 {
				sTag = 0.2
			}
		}

		// S_status: Authority Weight
		sStatus := 0.50
		switch statusUpper {
		case "ACTIVE", "CORE":
			sStatus = 1.0
		case "BOOKMARK":
			sStatus = 0.75
		case "CURATED_OPTIONAL":
			sStatus = 0.60
		case "REFERENCE":
			sStatus = 0.50
		}

		// S_diversity: Cross-Source Diversity Bonus
		family, hasFamily := SourceFamilyMap[id]
		if !hasFamily {
			family = FamilyDefault
		}
		sDiversity := 0.0
		if !selectedFamilies[family] {
			sDiversity = 0.15
		}

		compositeScore := (0.40 * sDomain) + (0.25 * sTag) + (0.20 * sStatus) + sDiversity

		role := "Visual Direction & Layout Inspiration"
		if family == FamilyAestheticBenchmark {
			role = "Aesthetic Benchmark & Creative Direction"
		} else if family == FamilyMovementTaxonomy {
			role = "Aesthetic Taxonomy & Design Contracts"
		} else if family == FamilySpecialistEcho {
			role = "Specialist Domain Echo & Pattern Reference"
		}

		scoredCandidates = append(scoredCandidates, &ScoredReference{
			Resource:    res,
			Score:       compositeScore,
			Domain:      "visual_research",
			Family:      family,
			MatchReason: fmt.Sprintf("Domain: %.2f, TagMatch: %.2f, StatusWeight: %.2f, DiversityBonus: %.2f", sDomain, sTag, sStatus, sDiversity),
			Role:        role,
		})
	}

	// 4. Sort by composite score descending
	sort.Slice(scoredCandidates, func(i, j int) bool {
		return scoredCandidates[i].Score > scoredCandidates[j].Score
	})

	// 5. Select diverse top N sources
	var selected []*ScoredReference
	for _, cand := range scoredCandidates {
		selected = append(selected, cand)
		selectedFamilies[cand.Family] = true
		if len(selected) >= opts.MaxSources {
			break
		}
	}

	// 6. Minimum Source Enforcement for High-Visual Tasks
	if len(selected) < opts.MinSources {
		defaultIDs := []string{"awwwards", "jiro-design"}
		for _, defID := range defaultIDs {
			alreadySelected := false
			for _, s := range selected {
				if s.Resource.ID == defID {
					alreadySelected = true
					break
				}
			}
			if !alreadySelected {
				var defRes *resources.Resource
				if c.Catalog != nil {
					defRes, _ = c.Catalog.FindByID(defID)
				}
				if defRes == nil {
					defRes = &resources.Resource{
						ID:             defID,
						Name:           defID,
						Representation: "research_source",
						Status:         "ACTIVE",
					}
				}
				selected = append(selected, &ScoredReference{
					Resource:    defRes,
					Score:       0.85,
					Domain:      "visual_research",
					Family:      SourceFamilyMap[defID],
					MatchReason: "Minimum source enforcement default",
					Role:        "Baseline Visual Direction",
				})
				if len(selected) >= opts.MinSources {
					break
				}
			}
		}
	}

	return selected, nil
}

// Coordinate queries selected sources and synthesizes visual decisions into a ResearchResult
func (c *ResearchCoordinator) Coordinate(ctx context.Context, req *ResearchRequest) (*ResearchResult, error) {
	if req == nil {
		return nil, fmt.Errorf("research request cannot be nil")
	}

	taskID := "task-unspecified"
	archetype := "standard_feature"
	qualityBar := "standard"
	if req.Task != nil {
		taskID = req.Task.ID
		qualityBar = "premium"
	}
	if req.Route != nil && req.Route.PrimaryArchetype != "" {
		archetype = req.Route.PrimaryArchetype
	}

	// Non-visual task bypass: if task does not require visual design and has no gap technologies, return minimal result
	if req.Task != nil && !req.Task.RequiresVisual && len(req.Task.SuggestedResources) == 0 {
		return &ResearchResult{
			TaskID:              taskID,
			Archetype:           archetype,
			QualityBar:          qualityBar,
			TotalSourcesQueried: 0,
			SelectedSources:     []*ScoredReference{},
			Findings:            []*ReferenceFinding{},
			GeneratedAt:         time.Now().UTC(),
		}, nil
	}

	// 1. Select references via multi-factor scoring
	selected, err := c.SelectReferences(ctx, req.Task, req.Route, req.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to select research references: %w", err)
	}

	// 2. Query each selected reference (using CuratedSourceFixtures for offline resilience)
	var findings []*ReferenceFinding
	var allPalettes []PaletteToken
	var allTypography []TypographyToken
	var allMotifs []VisualMotif

	for _, ref := range selected {
		finding, found := CuratedSourceFixtures[ref.Resource.ID]
		if !found {
			// Generate deterministic finding from resource metadata
			finding = &ReferenceFinding{
				SourceID:   ref.Resource.ID,
				SourceName: ref.Resource.Name,
				SourceURL:  ref.Resource.CanonicalURL,
				Category:   "curated-reference",
				KeyTakeaways: []string{
					fmt.Sprintf("Follow architectural conventions specified in %s (%s)", ref.Resource.Name, ref.Resource.CanonicalURL),
					"Maintain strict alignment with canonical capability routing",
				},
				Citations: []string{ref.Resource.CanonicalURL},
			}
		}
		findings = append(findings, finding)
		allPalettes = append(allPalettes, finding.ExtractedPalettes...)
		allTypography = append(allTypography, finding.ExtractedTypography...)
		allMotifs = append(allMotifs, finding.VisualMotifs...)
	}

	// 3. Synthesize Unified Color Palette
	// Priority: dark matte base (#0B0E14), elevated surface (#141A24), single calibrated accent (<80% sat), high contrast text
	synthesizedPalette := []PaletteToken{
		{Role: "--color-bg-base", Hex: "#0B0E14", HSL: "hsl(220, 20%, 6%)", Contrast: "17.4:1 against text (AAA)", SourceID: "awwwards"},
		{Role: "--color-surface-elevated", Hex: "#141A24", HSL: "hsl(216, 28%, 11%)", Contrast: "13.8:1 against text (AAA)", SourceID: "awwwards"},
		{Role: "--color-accent-primary", Hex: "#FF4B4B", HSL: "hsl(0, 100%, 65%)", Contrast: "4.8:1 against bg (AA)", SourceID: "awwwards"},
		{Role: "--color-text-headline", Hex: "#F0F4F8", HSL: "hsl(210, 33%, 96%)", Contrast: "17.4:1 against bg-base (AAA)", SourceID: "awwwards"},
		{Role: "--color-text-muted", Hex: "#8B9BB4", HSL: "hsl(216, 21%, 63%)", Contrast: "6.2:1 against bg-base (AA)", SourceID: "awwwards"},
		{Role: "--color-border-subtle", Hex: "#232D3F", HSL: "hsl(218, 28%, 19%)", Contrast: "2.5:1 decorative", SourceID: "awwwards"},
	}
	// Override accent if specialized finding provides one
	for _, p := range allPalettes {
		if p.Role == "--color-accent-primary" && p.Hex != "" {
			synthesizedPalette[2] = p
			break
		}
	}

	// 4. Synthesize Typographic Hierarchy
	// Triad: Distinctive Display Serif + Clean Geometric Sans + Monospace figures
	synthesizedTypography := []TypographyToken{
		{
			Role:          "display",
			FontFamily:    "Instrument Serif",
			Fallback:      "serif",
			SizeClamp:     "clamp(2.75rem, 6vw + 1rem, 5.5rem)",
			Weight:        "400",
			LineHeight:    "1.05",
			LetterSpacing: "-0.03em",
			SourceID:      "awwwards",
		},
		{
			Role:          "headline",
			FontFamily:    "Plus Jakarta Sans",
			Fallback:      "sans-serif",
			SizeClamp:     "clamp(1.5rem, 3vw + 0.5rem, 2.5rem)",
			Weight:        "600",
			LineHeight:    "1.2",
			LetterSpacing: "-0.02em",
			SourceID:      "awwwards",
		},
		{
			Role:          "body",
			FontFamily:    "Plus Jakarta Sans",
			Fallback:      "sans-serif",
			SizeClamp:     "clamp(1rem, 1vw + 0.5rem, 1.125rem)",
			Weight:        "400",
			LineHeight:    "1.6",
			LetterSpacing: "0",
			SourceID:      "awwwards",
		},
		{
			Role:          "mono",
			FontFamily:    "JetBrains Mono",
			Fallback:      "monospace",
			SizeClamp:     "0.875rem",
			Weight:        "500",
			LineHeight:    "1.4",
			LetterSpacing: "-0.01em",
			SourceID:      "awwwards",
		},
	}

	// 5. Interaction Dynamics
	interactionRules := InteractionDynamics{
		MotionCurve:     "cubic-bezier(0.16, 1, 0.3, 1)",
		SpringConfig:    "stiffness: 300, damping: 30, mass: 1",
		ActivePress:     "transform: scale(0.98); transition: transform 0.1s ease-out",
		ProhibitedProps: []string{"width", "height", "top", "left", "padding", "margin"},
		SourceID:        "jiro-design",
	}

	// 6. Gather Anti-Patterns
	bannedAntiPatterns := []string{
		"NO Inter or Space Grotesk in creative/display headlines",
		"NO neon purple or cosmic AI gradient backgrounds",
		"NO generic 3-equal-cards feature grid rows",
		"NO pure black (#000000) flat background surfaces",
		"NO horizontal overflow on mobile viewports (<390px)",
		"NO animation of layout-triggering properties (width, height, top, padding)",
	}
	if req.Route != nil && len(req.Route.AntiPatterns) > 0 {
		bannedAntiPatterns = append(bannedAntiPatterns, req.Route.AntiPatterns...)
	}

	result := &ResearchResult{
		TaskID:                taskID,
		Archetype:             archetype,
		QualityBar:            qualityBar,
		TotalSourcesQueried:   len(selected),
		SelectedSources:       selected,
		Findings:              findings,
		SynthesizedPalette:    synthesizedPalette,
		SynthesizedTypography: synthesizedTypography,
		SelectedMotifs:        allMotifs,
		InteractionRules:      interactionRules,
		BannedAntiPatterns:    bannedAntiPatterns,
		GeneratedAt:           time.Now().UTC(),
	}

	// 7. Write reference-log.md if output directory specified
	if req.ProjectOutputDir != "" {
		logPath := filepath.Join(req.ProjectOutputDir, "reference-log.md")
		if err := c.GenerateReferenceLog(ctx, result, logPath); err != nil {
			return nil, fmt.Errorf("failed to generate reference-log.md: %w", err)
		}
		result.ReferenceLogPath = logPath
	}

	return result, nil
}

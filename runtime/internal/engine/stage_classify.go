package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

type ClassifyStage struct{}

func NewClassifyStage() *ClassifyStage {
	return &ClassifyStage{}
}

func (s *ClassifyStage) Name() StageName {
	return StageClassify
}

func (s *ClassifyStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *ClassifyStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	// 1. Tag Aggregation & Normalization
	tagSet := make(map[string]bool)
	var allTags []string
	addTag := func(t string) {
		trimmed := strings.ToLower(strings.TrimSpace(t))
		if trimmed != "" && !tagSet[trimmed] {
			tagSet[trimmed] = true
			allTags = append(allTags, trimmed)
		}
	}

	if ctx.Discovery != nil {
		for _, t := range ctx.Discovery.DetectedTags {
			addTag(t)
		}
	}
	for _, t := range ctx.Task.Tags {
		addTag(t)
	}
	if ctx.Task.ArchetypeHint != "" {
		addTag(ctx.Task.ArchetypeHint)
	}

	// 2. Resolve Capability Routes from Graph
	var resolvedRoutes []resources.CapabilityRoute
	if ctx.Graph != nil && len(allTags) > 0 {
		resolvedRoutes = ctx.Graph.ResolveCapabilities(allTags)
	}

	// Fallback to explicit ArchetypeHint if not resolved
	if len(resolvedRoutes) == 0 && ctx.Graph != nil && ctx.Task.ArchetypeHint != "" {
		if route, ok := ctx.Graph.ResolveCapabilityRoute(ctx.Task.ArchetypeHint, []string{ctx.Task.ArchetypeHint}); ok {
			resolvedRoutes = append(resolvedRoutes, *route)
		}
	}

	// Determine Primary Archetype
	archetype := "standard-feature"
	if len(resolvedRoutes) > 0 && resolvedRoutes[0].PrimaryArchetype != "" {
		archetype = resolvedRoutes[0].PrimaryArchetype
	} else if ctx.Task.ArchetypeHint != "" {
		archetype = ctx.Task.ArchetypeHint
	} else if len(allTags) > 0 {
		archetype = allTags[0]
	}

	// 3. Determine Visual & Security Requirements
	reqLower := strings.ToLower(ctx.Task.RawRequest)
	requiresVisual := ctx.Task.Type == "DESIGN"
	requiresSecurity := false

	visualKeywords := []string{
		"ui", "frontend", "design", "landing", "3d", "creative", "showcase",
		"portfolio", "hud", "dashboard", "portal", "mobile", "responsive",
		"agro", "typography", "palette", "animation", "motion",
	}
	for _, kw := range visualKeywords {
		if strings.Contains(reqLower, kw) || tagSet[kw] {
			requiresVisual = true
			break
		}
	}

	securityKeywords := []string{
		"security", "audit", "auth", "login", "token", "secret", "leak",
		"vulnerability", "sanitize", "sast",
	}
	for _, kw := range securityKeywords {
		if strings.Contains(reqLower, kw) || tagSet[kw] {
			requiresSecurity = true
			break
		}
	}

	// Explicit override for backend bugfix
	if ctx.Task.Type == "BUGFIX" && !strings.Contains(reqLower, "ui") && !strings.Contains(reqLower, "frontend") && !strings.Contains(reqLower, "css") {
		requiresVisual = false
	}

	// 4. Human Approval Gate Evaluation
	requiresHumanGate := false
	gateReason := ""

	if requiresVisual {
		if ctx.Task.SkipVisualGate {
			requiresHumanGate = false
			gateReason = "Bypassed visual gate via TaskRequest.SkipVisualGate"
		} else if ctx.Task.UserOverride != nil && (ctx.Task.UserOverride.SkipVisualGate || ctx.Task.UserOverride.ForceBypassGate) {
			requiresHumanGate = false
			gateReason = "Bypassed visual gate via user override"
		} else {
			requiresHumanGate = true
			gateReason = fmt.Sprintf("High-impact visual task (Archetype: %s) requires design laboratory approval before coding", archetype)
		}
	}

	// 5. Gap Technologies Detection
	var gapTechnologies []string
	for _, resID := range ctx.Task.SuggestedResources {
		found := false
		if ctx.Catalog != nil {
			if _, ok := ctx.Catalog.FindByID(resID); ok {
				found = true
			}
		}
		if !found {
			gapTechnologies = append(gapTechnologies, resID)
		}
	}

	ctx.Classification = &ClassificationData{
		Archetype:         archetype,
		NormalizedTags:    allTags,
		ResolvedRoutes:    resolvedRoutes,
		RequiresVisual:    requiresVisual,
		RequiresSecurity:  requiresSecurity,
		RequiresHumanGate: requiresHumanGate,
		GateReason:        gateReason,
		GapTechnologies:   gapTechnologies,
	}

	return &StageResult{
		StageName: StageClassify,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Classification,
	}, nil
}

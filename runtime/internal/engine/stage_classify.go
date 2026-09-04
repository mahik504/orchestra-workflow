package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/verify"
)

// ClassifyStage runs the real classifier: it weighs every capability row in the
// graph against the request, records why the losing routes were declined, and
// arms the Design Lab gate before anything can be written.
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

	// 1. Gather every tag we know about before scoring.
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

	// 2. Score every capability row.
	cls := classifier.NewClassifierWithGraph(ctx.Graph)
	brief := cls.ClassifyBrief(ctx.Task.RawRequest, classifier.Options{
		TaskID:        ctx.Task.ID,
		ExtraTags:     allTags,
		ArchetypeHint: ctx.Task.ArchetypeHint,
		SkipLab:       ctx.Task.SkipVisualGate,
		DeclaredType:  ctx.Task.Type,
	})

	// 3. Nobody is here to answer a clarifying question inside a pipeline run,
	//    so an ambiguous brief resolves to the lower-risk route and says so.
	if brief.Ambiguous {
		brief.ResolveSilence()
	}

	// 4. Expand the chosen routes into full execution routes.
	var resolvedRoutes []resources.CapabilityRoute
	if ctx.Graph != nil {
		for _, cand := range brief.Selected {
			if route, ok := ctx.Graph.ResolveCapabilityRoute(cand.CapabilityID, cand.MatchedTags); ok {
				resolvedRoutes = append(resolvedRoutes, *route)
			}
		}
		// Keep the assumed or hinted route first, whatever the raw score said.
		for i, r := range resolvedRoutes {
			if r.CapabilityID == brief.CapabilityID && i != 0 {
				resolvedRoutes[0], resolvedRoutes[i] = resolvedRoutes[i], resolvedRoutes[0]
				break
			}
		}
		if len(resolvedRoutes) == 0 && ctx.Task.ArchetypeHint != "" {
			if route, ok := ctx.Graph.ResolveCapabilityRoute(ctx.Task.ArchetypeHint, []string{ctx.Task.ArchetypeHint}); ok {
				resolvedRoutes = append(resolvedRoutes, *route)
			}
		}
	}

	archetype := brief.Archetype
	if archetype == "" {
		archetype = "standard-feature"
	}

	// 5. Arm the gate. RequiresHumanGate now means "a Design Lab is owed",
	//    which is a narrower and more honest claim than "this looks visual".
	requiresHumanGate := brief.DesignLabRequired
	gateReason := brief.DesignLabReason
	if requiresHumanGate {
		if ctx.Task.SkipVisualGate {
			requiresHumanGate = false
			gateReason = "bypassed via TaskRequest.SkipVisualGate"
		} else if ctx.Task.UserOverride != nil && (ctx.Task.UserOverride.SkipVisualGate || ctx.Task.UserOverride.ForceBypassGate) {
			requiresHumanGate = false
			gateReason = "bypassed via user override"
		} else {
			gateReason = fmt.Sprintf("%s bar on %s: %s", brief.QualityBar, archetype, brief.DesignLabReason)
		}
	}
	brief.DesignLabRequired = requiresHumanGate
	brief.DesignLabReason = gateReason

	// 6. Resources the request names but the catalog has never heard of.
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
	gapTechnologies = append(gapTechnologies, brief.UnknownTechnology...)

	var declined []classifier.Candidate
	for _, c := range brief.Considered {
		if c.Declined {
			declined = append(declined, c)
		}
	}

	ctx.Classification = &ClassificationData{
		Archetype:         archetype,
		NormalizedTags:    allTags,
		ResolvedRoutes:    resolvedRoutes,
		RequiresVisual:    brief.RequiresVisual,
		RequiresSecurity:  brief.RequiresSecurity,
		RequiresHumanGate: requiresHumanGate,
		GateReason:        gateReason,
		GapTechnologies:   gapTechnologies,
		Brief:             brief,
		QualityBar:        brief.QualityBar,
		Platform:          brief.Platform,
		ResearchDepth:     brief.ResearchDepth,
		VerifyDepth:       brief.VerifyDepth,
		DeclinedRoutes:    declined,
	}

	// 7. Lock frontend writes until a direction is approved.
	ctx.DesignLab = verify.NewDesignLab(brief, ctx.Task.WorkspaceRoot)

	return &StageResult{
		StageName: StageClassify,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Classification,
	}, nil
}

package engine

import (
	"path/filepath"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/research"
	"github.com/user/orchestra-v3/internal/resources"
)

type ResearchStage struct{}

func NewResearchStage() *ResearchStage {
	return &ResearchStage{}
}

func (s *ResearchStage) Name() StageName {
	return StageResearch
}

func (s *ResearchStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	if ctx.Classification != nil && !ctx.Classification.RequiresVisual && len(ctx.Classification.GapTechnologies) == 0 {
		return true, "Non-visual task with zero unindexed dependencies skips visual research"
	}
	return false, ""
}

func (s *ResearchStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	// Instantiate ResearchCoordinator if not present
	if ctx.ResearchCoord == nil {
		ctx.ResearchCoord = research.NewResearchCoordinator(ctx.Catalog, ctx.Graph)
	}

	taskModel := &classifier.Task{
		ID:                 ctx.Task.ID,
		RawRequest:         ctx.Task.RawRequest,
		Type:               ctx.Task.Type,
		RequiresVisual:     ctx.Classification.RequiresVisual,
		RequiresSecurity:   ctx.Classification.RequiresSecurity,
		ExtractedKeywords:  ctx.Classification.NormalizedTags,
		SuggestedResources: ctx.Task.SuggestedResources,
	}

	var primaryRoute *resources.CapabilityRoute
	if len(ctx.Classification.ResolvedRoutes) > 0 {
		primaryRoute = &ctx.Classification.ResolvedRoutes[0]
	}

	outDir := filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra")
	req := &research.ResearchRequest{
		Task:             taskModel,
		Route:            primaryRoute,
		ProjectOutputDir: outDir,
		Options: research.SelectionOptions{
			MinSources:       2,
			MaxSources:       4,
			IncludeBookmarks: true,
			OfflineBenchmark: true,
		},
	}

	res, err := ctx.ResearchCoord.Coordinate(ctx.Ctx, req)
	if err != nil {
		return &StageResult{
			StageName: StageResearch,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	var refEntries []research.ScoredReference
	for _, s := range res.SelectedSources {
		if s != nil {
			refEntries = append(refEntries, *s)
		}
	}

	var colors []string
	for _, p := range res.SynthesizedPalette {
		colors = append(colors, p.Hex)
	}

	var typos []string
	for _, t := range res.SynthesizedTypography {
		typos = append(typos, t.FontFamily)
	}

	var motifs []string
	for _, m := range res.SelectedMotifs {
		motifs = append(motifs, m.Description)
	}

	ctx.Research = &ResearchData{
		ReferenceEntries:   refEntries,
		ReferenceLogPath:   res.ReferenceLogPath,
		ColorInspirations:  colors,
		TypographyPairings: typos,
		MotionPatterns:     motifs,
		SynthesizedResult:  res,
	}

	var artifacts []string
	if res.ReferenceLogPath != "" {
		ctx.AddArtifact("reference-log.md", res.ReferenceLogPath)
		artifacts = append(artifacts, res.ReferenceLogPath)
	}

	return &StageResult{
		StageName: StageResearch,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Research,
		Artifacts: artifacts,
	}, nil
}

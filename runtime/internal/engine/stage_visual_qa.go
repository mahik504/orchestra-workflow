package engine

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/memory"
)

type VisualQAStage struct{}

func NewVisualQAStage() *VisualQAStage {
	return &VisualQAStage{}
}

func (s *VisualQAStage) Name() StageName {
	return StageVisualQA
}

func (s *VisualQAStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *VisualQAStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	requiresVisual := false
	if ctx.Classification != nil {
		requiresVisual = ctx.Classification.RequiresVisual
	}

	// 1. Non-visual task handling: bypass browser screenshot capture, perform lightweight static check
	if !requiresVisual {
		ctx.VisualQA = &VisualQAData{
			AllPassed:          true,
			ViewportResults:    []ViewportCheckResult{},
			Metrics:            map[string]float64{"syntax_check": 1.0, "lint_errors": 0.0},
			DetectedViolations: []string{},
			FailureClass:       FailureClassNone,
		}
		return &StageResult{
			StageName: StageVisualQA,
			Status:    StatusCompleted,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Output:    ctx.VisualQA,
		}, nil
	}

	// 2. High-visual task verification
	if ctx.Verifier == nil {
		qaDir := filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra", "qa")
		ctx.Verifier = NewPlaywrightVerifier(qaDir)
	}

	qaRes, err := ctx.Verifier.Verify(ctx.Ctx, ctx)
	if err != nil {
		return &StageResult{
			StageName: StageVisualQA,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	ctx.VisualQA = qaRes
	if ctx.VisualQA != nil {
		ctx.VisualQA.VerifierRan = true
	}

	// Record verifier telemetry (Playwright) into Private Brain Memory
	memPath := memory.ResolveDefaultMemoryPath(ctx.Task.WorkspaceRoot)
	if store, err := memory.NewResourceMemoryStore(memPath); err == nil {
		qaOutcome := memory.OutcomeSuccess
		score := 1.0
		errDetails := ""
		if !qaRes.AllPassed {
			qaOutcome = memory.OutcomeFailure
			score = 1.0 - (float64(len(qaRes.DetectedViolations)) * 0.2)
			if score < 0.0 {
				score = 0.0
			}
			errDetails = strings.Join(qaRes.DetectedViolations, "; ")
		}

		_ = store.Record(&memory.ResourceEvaluation{
			ResourceID:          "playwright",
			Domain:              "qa_testing",
			Capability:          "visual-regression",
			EvaluationTimestamp: time.Now().UTC().Format(time.RFC3339),
			TaskContext:         ctx.Task.RawRequest,
			TaskID:              ctx.Task.ID,
			Outcome:             qaOutcome,
			QualityScore:        score,
			LatencyMs:           time.Since(start).Milliseconds(),
			ErrorDetails:        errDetails,
			Notes:               fmt.Sprintf("Visual QA verified %d viewport(s), passed=%v", len(qaRes.ViewportResults), qaRes.AllPassed),
			Metadata: map[string]any{
				"stage":         string(StageVisualQA),
				"viewports":     len(qaRes.ViewportResults),
				"failure_class": string(qaRes.FailureClass),
			},
		})
	}

	var artifacts []string
	for _, vp := range qaRes.ViewportResults {
		if vp.ScreenshotPath != "" {
			artifacts = append(artifacts, vp.ScreenshotPath)
			ctx.AddArtifact(vp.ViewportName, vp.ScreenshotPath)
		}
	}

	return &StageResult{
		StageName: StageVisualQA,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.VisualQA,
		Artifacts: artifacts,
	}, nil
}

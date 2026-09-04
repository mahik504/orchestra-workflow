package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/memory"
)

type IterateStage struct{}

func NewIterateStage() *IterateStage {
	return &IterateStage{}
}

func (s *IterateStage) Name() StageName {
	return StageIterate
}

func (s *IterateStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *IterateStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	if ctx.VisualQA == nil || ctx.VisualQA.AllPassed {
		ctx.Iteration.FinalVerdict = "PASSED"
		recordIterateMemory(ctx, "PASSED", time.Since(start).Milliseconds())
		return &StageResult{
			StageName: StageIterate,
			Status:    StatusCompleted,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Output:    ctx.Iteration,
		}, nil
	}

	// Visual QA has violations: evaluate iteration budget
	violations := ctx.VisualQA.DetectedViolations
	failureClass := ctx.VisualQA.FailureClass

	// Check if max iterations reached
	if ctx.Iteration.CurrentIteration >= ctx.Iteration.MaxIterations {
		ctx.Iteration.FinalVerdict = "EXHAUSTED"
		recordIterateMemory(ctx, "EXHAUSTED", time.Since(start).Milliseconds())

		// Record RecoveryPoint in handoff state
		resumeStep := "Stage 6 (Implement)"
		if failureClass == FailureClassTokenStyle {
			resumeStep = "Stage 5 (Design System)"
		}

		state, err := handoff.ReadState(ctx.Task.WorkspaceRoot)
		if err != nil || state == nil {
			state = &handoff.HandoffState{
				SessionID: ctx.Task.ID,
				Version:   3,
			}
		}
		state.FailureRecovery = &handoff.RecoveryPoint{
			FailedStep:     "Visual QA",
			ErrorReason:    fmt.Sprintf("Visual QA violations persisted after %d iterations: %s", ctx.Iteration.CurrentIteration, strings.Join(violations, "; ")),
			CanResume:      true,
			ResumeFromStep: resumeStep,
		}
		_ = handoff.WriteState(state, ctx.Task.WorkspaceRoot)

		return &StageResult{
			StageName: StageIterate,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     ErrMaxIterationsExceeded,
			Output:    ctx.Iteration,
		}, ErrMaxIterationsExceeded
	}

	// Budget remaining: increment iteration and plan loopback
	ctx.Iteration.CurrentIteration++

	var loopTarget StageName
	var correctiveDirective string

	switch failureClass {
	case FailureClassTokenStyle:
		loopTarget = StageDesignSystem
		correctiveDirective = fmt.Sprintf("Iteration %d: Remediate token/style violations: %s", ctx.Iteration.CurrentIteration, strings.Join(violations, "; "))
	case FailureClassLayoutCode:
		loopTarget = StageImplement
		correctiveDirective = fmt.Sprintf("Iteration %d: Remediate layout/code violations: %s", ctx.Iteration.CurrentIteration, strings.Join(violations, "; "))
	default:
		loopTarget = StageImplement
		correctiveDirective = fmt.Sprintf("Iteration %d: Address detected defects: %s", ctx.Iteration.CurrentIteration, strings.Join(violations, "; "))
	}

	ctx.Iteration.LoopTargetStage = loopTarget
	ctx.Iteration.CorrectiveFeedback = append(ctx.Iteration.CorrectiveFeedback, correctiveDirective)
	ctx.Iteration.History = append(ctx.Iteration.History, fmt.Sprintf("Iteration %d: Looping to %s due to %s", ctx.Iteration.CurrentIteration, loopTarget, failureClass))

	return &StageResult{
		StageName: StageIterate,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Iteration,
	}, nil
}

func recordIterateMemory(ctx *TaskContext, verdict string, latencyMs int64) {
	if ctx.Implementation == nil || len(ctx.Implementation.AcquiredResources) == 0 {
		return
	}
	memPath := memory.ResolveDefaultMemoryPath(ctx.Task.WorkspaceRoot)
	store, err := memory.NewResourceMemoryStore(memPath)
	if err != nil {
		return
	}

	outcome := memory.OutcomeSuccess
	score := 1.0
	if verdict != "PASSED" {
		outcome = memory.OutcomeFailure
		score = 0.5
	}

	primaryCap := "design-implementation"
	if ctx.Classification != nil && len(ctx.Classification.ResolvedRoutes) > 0 {
		primaryCap = ctx.Classification.ResolvedRoutes[0].CapabilityID
	}

	now := time.Now().UTC().Format(time.RFC3339)
	for _, resID := range ctx.Implementation.AcquiredResources {
		domain := "component_library"
		if ctx.Catalog != nil {
			if r, found := ctx.Catalog.FindByID(resID); found && len(r.Category) > 0 {
				domain = strings.ToLower(r.Category[0])
			}
		}

		_ = store.Record(&memory.ResourceEvaluation{
			ResourceID:          resID,
			Domain:              domain,
			Capability:          primaryCap,
			EvaluationTimestamp: now,
			TaskContext:         ctx.Task.RawRequest,
			TaskID:              ctx.Task.ID,
			Outcome:             outcome,
			QualityScore:        score,
			LatencyMs:           latencyMs,
			Notes:               fmt.Sprintf("Stage 8 Iterate verdict: %s after %d iteration(s)", verdict, ctx.Iteration.CurrentIteration),
			Metadata: map[string]any{
				"stage":      string(StageIterate),
				"verdict":    verdict,
				"iterations": ctx.Iteration.CurrentIteration,
			},
		})
	}
}

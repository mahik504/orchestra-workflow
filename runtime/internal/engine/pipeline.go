package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/memory"
	"github.com/user/orchestra-v3/internal/research"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/router"
)

type PipelineStatus string

const (
	PipelineStatusSuccess PipelineStatus = "SUCCESS"
	PipelineStatusFailed  PipelineStatus = "FAILED"
	PipelineStatusGated   PipelineStatus = "GATED_WAITING_APPROVAL"
	PipelineStatusPartial PipelineStatus = "PARTIAL_COMPLETION"
)

// DesignExecutionResult captures full execution telemetry and outputs of the 8-stage pipeline
type DesignExecutionResult struct {
	TaskID             string                      `json:"task_id"`
	Status             PipelineStatus              `json:"status"`
	Archetype          string                      `json:"archetype"`
	Classification     *ClassificationData         `json:"classification,omitempty"`
	ActiveCapabilities []string                    `json:"active_capabilities"`
	ResolvedRoutes     []resources.CapabilityRoute `json:"resolved_routes"`
	IterationCount     int                         `json:"iteration_count"`
	TotalDuration      time.Duration               `json:"total_duration"`
	StageResults       []*StageResult              `json:"stage_results"`
	DesignMDPath       string                      `json:"design_md_path,omitempty"`
	ReferenceLogPath   string                      `json:"reference_log_path,omitempty"`
	HandoffStatePath   string                      `json:"handoff_state_path,omitempty"`
	ModifiedFiles      []string                    `json:"modified_files,omitempty"`
	AcquiredResources  []string                    `json:"acquired_resources,omitempty"`
	VisualQAPassed     bool                        `json:"visual_qa_passed"`
	Screenshots        []string                    `json:"screenshots,omitempty"`
	QAMetrics          map[string]float64          `json:"qa_metrics,omitempty"`
	Violations         []string                    `json:"violations,omitempty"`
	StagesExecuted     []StageName                 `json:"stages_executed"`
	Error              error                       `json:"error,omitempty"`
}

// DesignPipeline manages and orchestrates the 8-stage design-first execution lifecycle
type DesignPipeline struct {
	Catalog       *resources.ResourceCatalog
	Graph         *resources.DesignResourceGraph
	Router        *router.Router
	ResearchCoord *research.ResearchCoordinator
	Verifier      VisualVerifier
	Stages        []Stage
	MaxIterations int
}

// NewDesignPipeline constructs an initialized DesignPipeline with all 8 stages
func NewDesignPipeline(
	cat *resources.ResourceCatalog,
	graph *resources.DesignResourceGraph,
	coord *research.ResearchCoordinator,
	verifier VisualVerifier,
) *DesignPipeline {
	if coord == nil && cat != nil && graph != nil {
		coord = research.NewResearchCoordinator(cat, graph)
	}

	rtr := router.NewRouterWithGraph(nil, cat, graph)

	p := &DesignPipeline{
		Catalog:       cat,
		Graph:         graph,
		Router:        rtr,
		ResearchCoord: coord,
		Verifier:      verifier,
		MaxIterations: 3,
	}

	p.Stages = []Stage{
		NewDiscoverStage(),
		NewClassifyStage(),
		NewResearchStage(),
		NewSynthesizeStage(),
		NewDesignSystemStage(),
		NewImplementStage(),
		NewVisualQAStage(),
		NewIterateStage(),
	}

	return p
}

// Plan executes a dry-run through Discover, Classify, and Synthesize to emit an Execution Manifest
func (p *DesignPipeline) Plan(ctx context.Context, req *TaskRequest) (*router.CompositionPlan, error) {
	req.DryRun = true
	taskCtx := NewTaskContext(ctx, req, p.Catalog, p.Graph, p.Router, p.ResearchCoord, p.Verifier)

	// Execute Discover
	if res, err := p.Stages[0].Execute(taskCtx); err != nil {
		return nil, fmt.Errorf("discover failed: %w", err)
	} else {
		taskCtx.SetStageResult(p.Stages[0].Name(), res)
	}

	// Execute Classify
	if res, err := p.Stages[1].Execute(taskCtx); err != nil {
		return nil, fmt.Errorf("classify failed: %w", err)
	} else {
		taskCtx.SetStageResult(p.Stages[1].Name(), res)
	}

	// Execute Synthesize
	if res, err := p.Stages[3].Execute(taskCtx); err != nil && err != ErrHumanGateRequired {
		return nil, fmt.Errorf("synthesize failed: %w", err)
	} else {
		taskCtx.SetStageResult(p.Stages[3].Name(), res)
	}

	// Build router.CompositionPlan
	primaryArch := taskCtx.Classification.Archetype
	plan := &router.CompositionPlan{
		PrimaryArchetype:    primaryArch,
		ExecutionDirectives: taskCtx.Synthesis.ActiveDirectives,
		EstimatedTokenCost:  taskCtx.Synthesis.TokenContextCost,
		RequiresHumanGate:   taskCtx.Classification.RequiresHumanGate,
		ApprovalReason:      taskCtx.Classification.GateReason,
		PipelineStage:       "Discover -> Classify -> Research -> Synthesize -> Design System -> Implement -> Visual QA -> Iterate",
		GapResearchNeeded:   taskCtx.Classification.GapTechnologies,
	}

	if p.Router != nil {
		// Populate selected capabilities from resolved routes
		for _, route := range taskCtx.Classification.ResolvedRoutes {
			plan.SelectedCapabilities = append(plan.SelectedCapabilities, &resources.Capability{
				ID:             route.CapabilityID,
				Name:           route.Name,
				Category:       resources.CategorySpecialist,
				CapabilityDesc: route.PrimaryArchetype,
				Status:         "ACTIVE",
			})
		}
	}

	return plan, nil
}

// Execute runs the full 8-stage closed-loop pipeline
func (p *DesignPipeline) Execute(ctx context.Context, req *TaskRequest) (*DesignExecutionResult, error) {
	startTime := time.Now()
	if req.MaxIterations <= 0 {
		req.MaxIterations = p.MaxIterations
	}

	taskCtx := NewTaskContext(ctx, req, p.Catalog, p.Graph, p.Router, p.ResearchCoord, p.Verifier)

	var stageResultsList []*StageResult

	// stageIndexMap maps StageName to stage slice index
	stageIndexMap := make(map[StageName]int)
	for i, stage := range p.Stages {
		stageIndexMap[stage.Name()] = i
	}

	// Execution loop
	currentStageIdx := 0
	for currentStageIdx < len(p.Stages) {
		stage := p.Stages[currentStageIdx]

		// Check skip condition
		skip, reason := stage.ShouldSkip(taskCtx)
		if skip {
			stageRes := &StageResult{
				StageName:  stage.Name(),
				Status:     StatusSkipped,
				StartTime:  time.Now(),
				EndTime:    time.Now(),
				SkipReason: reason,
			}
			taskCtx.SetStageResult(stage.Name(), stageRes)
			stageResultsList = append(stageResultsList, stageRes)
			currentStageIdx++
			continue
		}

		// Execute stage
		stageRes, err := stage.Execute(taskCtx)
		if stageRes != nil {
			taskCtx.SetStageResult(stage.Name(), stageRes)
			stageResultsList = append(stageResultsList, stageRes)
		}

		if err != nil {
			if err == ErrHumanGateRequired {
				return p.buildResult(taskCtx, PipelineStatusGated, stageResultsList, startTime, err), err
			}
			failRes := p.buildResult(taskCtx, PipelineStatusFailed, stageResultsList, startTime, err)
			_ = p.RecordPipelineMemory(taskCtx, failRes)
			return failRes, err
		}

		// Handle closed-loop return from Stage 8 (Iterate)
		if stage.Name() == StageIterate {
			if taskCtx.Iteration != nil && taskCtx.Iteration.LoopTargetStage != "" {
				targetStage := taskCtx.Iteration.LoopTargetStage
				taskCtx.Iteration.LoopTargetStage = "" // reset
				if targetIdx, ok := stageIndexMap[targetStage]; ok {
					currentStageIdx = targetIdx // Jump to target stage (e.g. Stage 5 or Stage 6)
					continue
				}
			}
		}

		currentStageIdx++
	}

	result := p.buildResult(taskCtx, PipelineStatusSuccess, stageResultsList, startTime, nil)
	_ = p.RecordPipelineMemory(taskCtx, result)
	return result, nil
}

func (p *DesignPipeline) buildResult(
	taskCtx *TaskContext,
	status PipelineStatus,
	stages []*StageResult,
	startTime time.Time,
	err error,
) *DesignExecutionResult {
	res := &DesignExecutionResult{
		TaskID:         taskCtx.Task.ID,
		Status:         status,
		TotalDuration:  time.Since(startTime),
		StageResults:   stages,
		StagesExecuted: taskCtx.StagesExecuted,
		Error:          err,
	}

	if taskCtx.Classification != nil {
		res.Archetype = taskCtx.Classification.Archetype
		res.Classification = taskCtx.Classification
		res.ResolvedRoutes = taskCtx.Classification.ResolvedRoutes
		for _, r := range taskCtx.Classification.ResolvedRoutes {
			res.ActiveCapabilities = append(res.ActiveCapabilities, r.CapabilityID)
		}
	}

	if taskCtx.Iteration != nil {
		res.IterationCount = taskCtx.Iteration.CurrentIteration
	}

	if taskCtx.Research != nil {
		res.ReferenceLogPath = taskCtx.Research.ReferenceLogPath
	}

	if taskCtx.DesignSystem != nil {
		res.DesignMDPath = taskCtx.DesignSystem.DesignMDPath
	}

	if taskCtx.Implementation != nil {
		res.HandoffStatePath = taskCtx.Implementation.HandoffStatePath
		res.AcquiredResources = taskCtx.Implementation.AcquiredResources
		for _, f := range taskCtx.Implementation.ModifiedFiles {
			res.ModifiedFiles = append(res.ModifiedFiles, f.Path)
		}
	}

	if taskCtx.VisualQA != nil {
		res.VisualQAPassed = taskCtx.VisualQA.AllPassed
		res.QAMetrics = taskCtx.VisualQA.Metrics
		res.Violations = taskCtx.VisualQA.DetectedViolations
		for _, vp := range taskCtx.VisualQA.ViewportResults {
			if vp.ScreenshotPath != "" {
				res.Screenshots = append(res.Screenshots, vp.ScreenshotPath)
			}
		}
	}

	return res
}

// RecordPipelineMemory persists evaluations into private brain memory for all active resources
func (p *DesignPipeline) RecordPipelineMemory(taskCtx *TaskContext, res *DesignExecutionResult) error {
	if taskCtx == nil || res == nil {
		return nil
	}
	memPath := memory.ResolveDefaultMemoryPath(taskCtx.Task.WorkspaceRoot)
	store, err := memory.NewResourceMemoryStore(memPath)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	requiresVisual := taskCtx.Classification != nil && taskCtx.Classification.RequiresVisual
	verifierRan := taskCtx.VisualQA != nil && taskCtx.VisualQA.VerifierRan
	if requiresVisual && verifierRan {
		qaOutcome := memory.OutcomeSuccess
		qaScore := 1.0
		errDetails := ""
		if !taskCtx.VisualQA.AllPassed {
			qaOutcome = memory.OutcomeFailure
			violationsCount := len(taskCtx.VisualQA.DetectedViolations)
			qaScore = 1.0 - (float64(violationsCount) * 0.2)
			if qaScore < 0.0 {
				qaScore = 0.0
			}
			errDetails = strings.Join(taskCtx.VisualQA.DetectedViolations, "; ")
		}

		qaStageRes := taskCtx.GetStageResult(StageVisualQA)
		latency := int64(0)
		if qaStageRes != nil {
			latency = qaStageRes.Duration.Milliseconds()
		}

		_ = store.Record(&memory.ResourceEvaluation{
			ResourceID:          "playwright",
			Domain:              "qa_testing",
			Capability:          "visual-regression",
			EvaluationTimestamp: now,
			TaskContext:         taskCtx.Task.RawRequest,
			TaskID:              taskCtx.Task.ID,
			Outcome:             qaOutcome,
			QualityScore:        qaScore,
			LatencyMs:           latency,
			ErrorDetails:        errDetails,
			Notes:               fmt.Sprintf("Visual QA pass status: %v, viewports: %d", taskCtx.VisualQA.AllPassed, len(taskCtx.VisualQA.ViewportResults)),
			Metadata: map[string]any{
				"stage":         string(StageVisualQA),
				"failure_class": string(taskCtx.VisualQA.FailureClass),
				"verifier_ran":  true,
			},
		})
	}

	taskOutcome := memory.OutcomeSuccess
	qualityScore := 1.0
	if res.Status != PipelineStatusSuccess {
		taskOutcome = memory.OutcomeFailure
		qualityScore = 0.5
	}

	primaryCap := "design-implementation"
	if len(res.ActiveCapabilities) > 0 {
		primaryCap = res.ActiveCapabilities[0]
	}

	recorded := map[string]bool{}
	recordOne := func(resID, why, installed, command, verification string, outcome memory.Outcome, score float64) {
		resID = strings.TrimSpace(resID)
		if resID == "" || recorded[resID] {
			return
		}
		recorded[resID] = true
		domain := "component_library"
		if p.Catalog != nil {
			if r, found := p.Catalog.FindByID(resID); found && len(r.Category) > 0 {
				domain = strings.ToLower(r.Category[0])
			}
		}
		meta := map[string]any{
			"archetype":    res.Archetype,
			"status":       string(res.Status),
			"why_selected": why,
			"verification": verification,
		}
		if installed != "" {
			meta["installed_path"] = installed
		}
		if command != "" {
			meta["executed_command"] = command
		}
		_ = store.Record(&memory.ResourceEvaluation{
			ResourceID:          resID,
			Domain:              domain,
			Capability:          primaryCap,
			EvaluationTimestamp: now,
			TaskContext:         taskCtx.Task.RawRequest,
			TaskID:              taskCtx.Task.ID,
			Outcome:             outcome,
			QualityScore:        score,
			LatencyMs:           res.TotalDuration.Milliseconds(),
			Notes:               fmt.Sprintf("Pipeline %s after %d iteration(s)", res.Status, res.IterationCount),
			Metadata:            meta,
		})
	}

	whyFor := func(id string) string {
		if taskCtx.Classification != nil {
			for _, line := range taskCtx.Classification.OverlayActivations {
				if line == id {
					return "overlay trigger matched this task"
				}
			}
		}
		return "acquired during implement stage"
	}
	installedFor := func(id string) string {
		if taskCtx.Implementation != nil && taskCtx.Implementation.InstalledPaths != nil {
			return taskCtx.Implementation.InstalledPaths[id]
		}
		return ""
	}
	cmdFor := func(id string) string {
		if taskCtx.Implementation != nil && taskCtx.Implementation.AcquireCommands != nil {
			return taskCtx.Implementation.AcquireCommands[id]
		}
		return ""
	}

	for _, resID := range res.AcquiredResources {
		o, s := taskOutcome, qualityScore
		if installedFor(resID) != "" {
			o, s = memory.OutcomeSuccess, 1.0
		}
		recordOne(resID, whyFor(resID), installedFor(resID), cmdFor(resID), "acquired in implement stage", o, s)
	}
	if taskCtx.Classification != nil {
		for _, resID := range taskCtx.Classification.OverlayActivations {
			recordOne(resID, whyFor(resID), installedFor(resID), cmdFor(resID), "overlay resource activated for this task", taskOutcome, qualityScore)
		}
	}

	return nil
}

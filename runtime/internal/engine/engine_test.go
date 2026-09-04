package engine

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/resources"
)

func setupEngineFixtures(t *testing.T) (*resources.ResourceCatalog, *resources.DesignResourceGraph, string) {
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

	cat, err := resources.LoadResourceCatalog(filepath.Join(regDir, "resources.json"))
	if err != nil {
		t.Fatalf("Failed to load catalog: %v", err)
	}

	graph, err := resources.LoadDesignGraph(filepath.Join(regDir, "design-resource-graph.json"))
	if err != nil {
		t.Fatalf("Failed to load design graph: %v", err)
	}

	tempWorkdir, err := os.MkdirTemp("", "orchestra_engine_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp workdir: %v", err)
	}

	return cat, graph, tempWorkdir
}

func TestPipeline_SequentialStageExecution(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:             "task-seq-01",
		RawRequest:     "Build modern portfolio showcase with animations",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected status SUCCESS, got %s", res.Status)
	}

	expectedStages := []StageName{
		StageDiscover,
		StageClassify,
		StageResearch,
		StageSynthesize,
		StageDesignSystem,
		StageImplement,
		StageVisualQA,
		StageIterate,
	}

	if len(res.StagesExecuted) != len(expectedStages) {
		t.Fatalf("Expected %d stages executed, got %d: %v", len(expectedStages), len(res.StagesExecuted), res.StagesExecuted)
	}

	for i, stage := range expectedStages {
		if res.StagesExecuted[i] != stage {
			t.Errorf("Stage index %d: expected %s, got %s", i, stage, res.StagesExecuted[i])
		}
	}
}

func TestPipeline_StandardTask_BypassesVisualQA(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:            "task-backend-bug",
		RawRequest:    "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir,
		Type:          "BUGFIX",
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Standard task pipeline failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS status, got %s", res.Status)
	}

	// Verify Research Stage was skipped
	researchSkipped := false
	for _, sr := range res.StageResults {
		if sr.StageName == StageResearch && sr.Status == StatusSkipped {
			researchSkipped = true
			if sr.SkipReason == "" {
				t.Errorf("Expected non-empty SkipReason for skipped Research stage")
			}
			break
		}
	}
	if !researchSkipped {
		t.Errorf("Expected Research stage to be skipped for backend bug task")
	}

	// Verify no screenshots generated
	if len(res.Screenshots) != 0 {
		t.Errorf("Expected 0 screenshots for non-visual task, got %d", len(res.Screenshots))
	}
}

func TestPipeline_PremiumVisualTask_FullExecution(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:             "task-premium-visual",
		RawRequest:     "Award-winning 3D spatial portfolio with GSAP animations",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Premium visual pipeline failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS, got %s", res.Status)
	}

	// Verify reference-log.md exists and contains sections
	if res.ReferenceLogPath == "" {
		t.Errorf("Expected reference-log.md path to be populated")
	} else {
		content, err := os.ReadFile(res.ReferenceLogPath)
		if err != nil {
			t.Fatalf("Failed to read reference-log.md: %v", err)
		}
		if !strings.Contains(string(content), "# Visual Research Reference Log") {
			t.Errorf("reference-log.md missing title header")
		}
	}

	// Verify DESIGN.md exists and contains tokens
	if res.DesignMDPath == "" {
		t.Errorf("Expected DESIGN.md path to be populated")
	} else {
		content, err := os.ReadFile(res.DesignMDPath)
		if err != nil {
			t.Fatalf("Failed to read DESIGN.md: %v", err)
		}
		if !strings.Contains(string(content), "# DESIGN.md — Design System Contract") {
			t.Errorf("DESIGN.md missing title header")
		}
		if !strings.Contains(string(content), "--bg-base") {
			t.Errorf("DESIGN.md missing color tokens")
		}
	}

	// Verify screenshots generated
	if len(res.Screenshots) < 3 {
		t.Errorf("Expected at least 3 viewport screenshots, got %d", len(res.Screenshots))
	}

	// Verify handoff state exists
	if res.HandoffStatePath == "" {
		t.Errorf("Expected HandoffStatePath to be populated")
	} else {
		state, err := handoff.ReadState(workdir)
		if err != nil {
			t.Fatalf("Failed to read handoff state: %v", err)
		}
		if state.SessionID != req.ID {
			t.Errorf("Expected task ID %s in handoff state, got %s", req.ID, state.SessionID)
		}
	}
}

func TestPipeline_HumanApprovalGate(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// 1. Without SkipVisualGate -> must halt with ErrHumanGateRequired
	reqGated := &TaskRequest{
		ID:             "task-gate-test",
		RawRequest:     "Design interactive landing page with bold typography",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: false,
	}

	resGated, err := pipeline.Execute(context.Background(), reqGated)
	if err != ErrHumanGateRequired {
		t.Fatalf("Expected ErrHumanGateRequired, got err=%v, status=%s", err, resGated.Status)
	}
	if resGated.Status != PipelineStatusGated {
		t.Errorf("Expected status GATED_WAITING_APPROVAL, got %s", resGated.Status)
	}

	// 2. With SkipVisualGate -> proceeds to completion
	reqApproved := &TaskRequest{
		ID:             "task-gate-approved",
		RawRequest:     "Design interactive landing page with bold typography",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	resApproved, err := pipeline.Execute(context.Background(), reqApproved)
	if err != nil {
		t.Fatalf("Approved task failed: %v", err)
	}
	if resApproved.Status != PipelineStatusSuccess {
		t.Errorf("Expected status SUCCESS, got %s", resApproved.Status)
	}
}

func TestPipeline_IterationSelfHealing_LayoutCode(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.SimulateMobileOverflow = true
	mockVerifier.FailUntilIteration = 1 // Fails on iteration 0, heals on iteration 1!

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)

	req := &TaskRequest{
		ID:             "task-heal-layout",
		RawRequest:     "Build mobile dashboard with graphs",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  3,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Self-healing pipeline failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS after self-healing, got %s", res.Status)
	}

	if res.IterationCount < 1 {
		t.Errorf("Expected at least 1 iteration loop, got %d", res.IterationCount)
	}
}

func TestPipeline_IterationSelfHealing_TokenStyle(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.SimulateTokenStyleDefect = true
	mockVerifier.FailUntilIteration = 1 // Fails iteration 0, heals on iteration 1

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)

	req := &TaskRequest{
		ID:             "task-heal-token",
		RawRequest:     "Build high-contrast dark theme showcase",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  3,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Self-healing pipeline failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS after self-healing, got %s", res.Status)
	}

	if res.IterationCount < 1 {
		t.Errorf("Expected at least 1 iteration loop, got %d", res.IterationCount)
	}
}

func TestPipeline_MaxIterationsExceeded(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.AlwaysFail = true
	mockVerifier.SimulateMobileOverflow = true

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)
	pipeline.MaxIterations = 2

	req := &TaskRequest{
		ID:             "task-persistent-failure",
		RawRequest:     "Build layout with permanent overflow defect",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  2,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrMaxIterationsExceeded {
		t.Fatalf("Expected ErrMaxIterationsExceeded, got err=%v", err)
	}

	if res.Status != PipelineStatusFailed {
		t.Errorf("Expected status FAILED, got %s", res.Status)
	}

	if res.IterationCount != 2 {
		t.Errorf("Expected exactly 2 iterations executed, got %d", res.IterationCount)
	}

	// Verify recovery point recorded in handoff
	state, err := handoff.ReadState(workdir)
	if err != nil {
		t.Fatalf("Failed to read handoff state after failure: %v", err)
	}
	if state.FailureRecovery == nil {
		t.Fatalf("Expected recovery point recorded in handoff state")
	}
	if state.FailureRecovery.FailedStep != "Visual QA" {
		t.Errorf("Expected failed step 'Visual QA', got %s", state.FailureRecovery.FailedStep)
	}
}

func TestPipeline_PlanDryRun(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:            "task-plan-test",
		RawRequest:    "Build award-winning creative landing page with 3D WebGL",
		WorkspaceRoot: workdir,
		Type:          "DESIGN",
	}

	plan, err := pipeline.Plan(context.Background(), req)
	if err != nil {
		t.Fatalf("Pipeline.Plan failed: %v", err)
	}

	if plan.PrimaryArchetype == "" {
		t.Errorf("Expected non-empty PrimaryArchetype in plan")
	}

	if plan.EstimatedTokenCost <= 1500 {
		t.Errorf("Expected token cost to accumulate capabilities, got %.0f", plan.EstimatedTokenCost)
	}

	if !plan.RequiresHumanGate {
		t.Errorf("Expected high-visual task to require human gate in plan")
	}

	manifest := plan.GenerateExecutionManifest()
	if !strings.Contains(manifest, "Orchestra V3 Capability Execution Manifest") {
		t.Errorf("Manifest missing title header")
	}
}

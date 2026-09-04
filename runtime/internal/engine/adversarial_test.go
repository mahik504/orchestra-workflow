package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/user/orchestra-v3/internal/handoff"
)

// OscillatingVerifier alternates failure classes between LAYOUT_CODE and TOKEN_STYLE
type OscillatingVerifier struct {
	mu          sync.Mutex
	callCount   int
	HealAfter   int // if > 0, passes when callCount >= HealAfter
	AlwaysFail  bool
	LastFailure string
}

func NewOscillatingVerifier(healAfter int) *OscillatingVerifier {
	return &OscillatingVerifier{
		HealAfter:  healAfter,
		AlwaysFail: healAfter <= 0,
	}
}

func (o *OscillatingVerifier) Name() string {
	return "oscillating_adversarial_verifier"
}

func (o *OscillatingVerifier) Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	defer func() { o.callCount++ }()

	if !o.AlwaysFail && o.HealAfter > 0 && o.callCount >= o.HealAfter {
		o.LastFailure = FailureClassNone
		return &VisualQAData{
			AllPassed: true,
			ViewportResults: []ViewportCheckResult{
				{ViewportName: "desktop", Width: 1440, Height: 900, Passed: true},
				{ViewportName: "tablet", Width: 768, Height: 1024, Passed: true},
				{ViewportName: "mobile", Width: 390, Height: 844, Passed: true},
			},
			Metrics:      map[string]float64{"contrast_ratio": 14.0},
			FailureClass: FailureClassNone,
		}, nil
	}

	// Even calls: Mobile horizontal overflow -> LAYOUT_CODE
	// Odd calls: Pure black / typography anti-pattern -> TOKEN_STYLE
	if o.callCount%2 == 0 {
		o.LastFailure = FailureClassLayoutCode
		return &VisualQAData{
			AllPassed: false,
			ViewportResults: []ViewportCheckResult{
				{ViewportName: "desktop", Width: 1440, Height: 900, Passed: true},
				{ViewportName: "tablet", Width: 768, Height: 1024, Passed: true},
				{ViewportName: "mobile", Width: 390, Height: 844, HasOverflow: true, MaxScrollWidth: 440, Passed: false},
			},
			Metrics:            map[string]float64{"scroll_excess": 50.0},
			DetectedViolations: []string{fmt.Sprintf("Call %d: Mobile viewport horizontal overflow (440px > 390px)", o.callCount)},
			FailureClass:       FailureClassLayoutCode,
		}, nil
	}

	o.LastFailure = FailureClassTokenStyle
	return &VisualQAData{
		AllPassed: false,
		ViewportResults: []ViewportCheckResult{
			{ViewportName: "desktop", Width: 1440, Height: 900, Passed: false},
			{ViewportName: "tablet", Width: 768, Height: 1024, Passed: false},
			{ViewportName: "mobile", Width: 390, Height: 844, Passed: false},
		},
		Metrics:            map[string]float64{"contrast_ratio": 2.1},
		DetectedViolations: []string{fmt.Sprintf("Call %d: Banned pure black (#000000) token in surface background", o.callCount)},
		FailureClass:       FailureClassTokenStyle,
	}, nil
}

// -----------------------------------------------------------------------------
// Focus Area 1: State Transition Invariants
// -----------------------------------------------------------------------------

func TestAdversarial_StateTransition_SequentialExecution(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:             "task-state-invariants",
		RawRequest:     "Design high-end luxury e-commerce landing page",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	expectedOrder := []StageName{
		StageDiscover,
		StageClassify,
		StageResearch,
		StageSynthesize,
		StageDesignSystem,
		StageImplement,
		StageVisualQA,
		StageIterate,
	}

	if len(res.StagesExecuted) != len(expectedOrder) {
		t.Fatalf("Expected exactly %d stages executed, got %d: %v", len(expectedOrder), len(res.StagesExecuted), res.StagesExecuted)
	}

	for i, expected := range expectedOrder {
		if res.StagesExecuted[i] != expected {
			t.Errorf("Stage index %d invariant violation: expected %s, got %s", i, expected, res.StagesExecuted[i])
		}
	}
}

func TestAdversarial_StateTransition_EarlyHaltOnFailure(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// High-visual task without SkipVisualGate must halt at StageSynthesize
	req := &TaskRequest{
		ID:             "task-gated-halt",
		RawRequest:     "Build 3D interactive WebGL showcase",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: false,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrHumanGateRequired {
		t.Fatalf("Expected ErrHumanGateRequired, got: %v", err)
	}

	if res.Status != PipelineStatusGated {
		t.Errorf("Expected status %s, got %s", PipelineStatusGated, res.Status)
	}

	// Invariant: downstream stages (DesignSystem, Implement, VisualQA, Iterate) MUST NOT have executed!
	bannedStages := map[StageName]bool{
		StageDesignSystem: true,
		StageImplement:    true,
		StageVisualQA:     true,
		StageIterate:      true,
	}

	for _, executed := range res.StagesExecuted {
		if bannedStages[executed] {
			t.Errorf("Stage invariant broken: stage %s executed despite human gate halt!", executed)
		}
	}

	// Verify handoff state was NOT created
	handoffPath := filepath.Join(workdir, ".orchestra", "handoff", "state.json")
	if _, err := os.Stat(handoffPath); err == nil {
		t.Errorf("Handoff state was prematurely created despite human gate block")
	}
}

func TestAdversarial_StateTransition_DirectStageExecutionWithoutPrerequisites(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	req := &TaskRequest{
		ID:            "task-direct-stage",
		RawRequest:    "Adversarial direct stage execution",
		WorkspaceRoot: workdir,
	}

	// 1. StageResearch executed on raw context without Classify
	t.Run("Research_Without_Classify", func(t *testing.T) {
		taskCtx := NewTaskContext(context.Background(), req, cat, graph, nil, nil, nil)
		stage := NewResearchStage()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("CONFIRMED: Direct ResearchStage.Execute without Classify panicked as expected: %v", r)
			}
		}()

		_, err := stage.Execute(taskCtx)
		t.Logf("ResearchStage executed without panic, err: %v", err)
	})

	// 2. StageIterate executed on raw context without VisualQA
	t.Run("Iterate_Without_VisualQA", func(t *testing.T) {
		taskCtx := NewTaskContext(context.Background(), req, cat, graph, nil, nil, nil)
		stage := NewIterateStage()

		res, err := stage.Execute(taskCtx)
		if err != nil {
			t.Fatalf("IterateStage returned error: %v", err)
		}
		if taskCtx.Iteration.FinalVerdict != "PASSED" {
			t.Errorf("Expected default PASSED when VisualQA is nil, got %s", taskCtx.Iteration.FinalVerdict)
		}
		t.Logf("IterateStage result status: %s, verdict: %s", res.Status, taskCtx.Iteration.FinalVerdict)
	})
}

// -----------------------------------------------------------------------------
// Focus Area 2: Rapid Max Iteration Exhaustion & Recovery Point Generation
// -----------------------------------------------------------------------------

func TestAdversarial_RapidIterationExhaustion_MaxIter1(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.AlwaysFail = true
	mockVerifier.SimulateMobileOverflow = true

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)

	req := &TaskRequest{
		ID:             "task-rapid-exhaust-1",
		RawRequest:     "Build mobile app with persistent overflow",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  1,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrMaxIterationsExceeded {
		t.Fatalf("Expected ErrMaxIterationsExceeded, got: %v", err)
	}

	if res.Status != PipelineStatusFailed {
		t.Errorf("Expected status %s, got %s", PipelineStatusFailed, res.Status)
	}

	if res.IterationCount != 1 {
		t.Errorf("Expected iteration count 1, got %d", res.IterationCount)
	}

	// Verify Recovery Point in handoff state
	state, err := handoff.ReadState(workdir)
	if err != nil {
		t.Fatalf("Failed to read handoff state: %v", err)
	}
	if state.FailureRecovery == nil {
		t.Fatalf("Expected non-nil FailureRecovery in handoff state")
	}
	if !state.FailureRecovery.CanResume {
		t.Errorf("Expected CanResume=true, got false")
	}
	if state.FailureRecovery.FailedStep != "Visual QA" {
		t.Errorf("Expected FailedStep='Visual QA', got %s", state.FailureRecovery.FailedStep)
	}
	if state.FailureRecovery.ResumeFromStep != "Stage 6 (Implement)" {
		t.Errorf("Expected ResumeFromStep='Stage 6 (Implement)', got %s", state.FailureRecovery.ResumeFromStep)
	}
}

func TestAdversarial_RapidIterationExhaustion_TokenStyleDefect(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.AlwaysFail = true
	mockVerifier.SimulateTokenStyleDefect = true

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)

	req := &TaskRequest{
		ID:             "task-rapid-exhaust-token",
		RawRequest:     "Build showcase with illegal black theme",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  2,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrMaxIterationsExceeded {
		t.Fatalf("Expected ErrMaxIterationsExceeded, got: %v", err)
	}

	if res.IterationCount != 2 {
		t.Errorf("Expected iteration count 2, got %d", res.IterationCount)
	}

	// Verify Recovery Point directs resumption to Stage 5 (Design System)
	state, err := handoff.ReadState(workdir)
	if err != nil {
		t.Fatalf("Failed to read handoff state: %v", err)
	}
	if state.FailureRecovery == nil {
		t.Fatalf("Expected non-nil FailureRecovery")
	}
	if state.FailureRecovery.ResumeFromStep != "Stage 5 (Design System)" {
		t.Errorf("Expected ResumeFromStep='Stage 5 (Design System)', got %s", state.FailureRecovery.ResumeFromStep)
	}
}

func TestAdversarial_RapidIterationStress_50Cycles(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.AlwaysFail = true
	mockVerifier.SimulateMobileOverflow = true

	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)

	for i := 0; i < 50; i++ {
		taskWorkdir := filepath.Join(workdir, fmt.Sprintf("cycle_%d", i))
		req := &TaskRequest{
			ID:             fmt.Sprintf("task-stress-%d", i),
			RawRequest:     "Adversarial rapid exhaustion cycle",
			WorkspaceRoot:  taskWorkdir,
			Type:           "DESIGN",
			SkipVisualGate: true,
			MaxIterations:  1,
		}

		res, err := pipeline.Execute(context.Background(), req)
		if err != ErrMaxIterationsExceeded {
			t.Fatalf("Cycle %d: expected ErrMaxIterationsExceeded, got: %v", i, err)
		}
		if res.Status != PipelineStatusFailed {
			t.Fatalf("Cycle %d: expected status FAILED, got %s", i, res.Status)
		}

		state, err := handoff.ReadState(taskWorkdir)
		if err != nil || state == nil || state.FailureRecovery == nil {
			t.Fatalf("Cycle %d: failed to read valid recovery point from state.json", i)
		}
	}
}

// -----------------------------------------------------------------------------
// Focus Area 3: Malformed & Adversarial Task Requests
// -----------------------------------------------------------------------------

func TestAdversarial_MalformedTask_EmptyFields(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// Completely empty task request
	req := &TaskRequest{
		WorkspaceRoot: workdir,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Engine failed on empty task request: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS on empty task request fallback, got %s", res.Status)
	}

	if res.Archetype == "" {
		t.Errorf("Expected default archetype on empty request, got empty string")
	}
}

func TestAdversarial_MalformedTask_UnicodeAttacks(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// Zalgo, null bytes, RTL override, multi-byte UTF-8, emoji spam
	zalgo := "T̷͓̈h̶̯̄e̶̯͛ ̷͇̌q̸̳̽u̷̮͌i̶̙̒c̶̫̏k̷̙͘ ̸͙̓b̷̬͐r̴͎̍o̷̟͝w̷̱̔ñ̷͕ ̴̜̅f̷̺̈́ó̵̧x̷͔͝"
	rtl := "\u202Ereversed_text_override\u202C"
	emojis := "🚀🔥💻✨🎨🛡️💥"
	massivePayload := strings.Repeat("Awwwards-winning 3D showcase with GSAP animations! ", 2000) // ~100KB

	payload := zalgo + " " + rtl + " " + emojis + "\n" + massivePayload

	req := &TaskRequest{
		ID:             "task-unicode-attack-🛡️",
		RawRequest:     payload,
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Engine crashed or failed on unicode attack payload: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS status on unicode payload, got %s", res.Status)
	}

	if res.Archetype == "" {
		t.Errorf("Expected classified archetype, got empty")
	}
}

func TestAdversarial_MalformedTask_InjectionPayloads(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	injections := []string{
		"'; DROP TABLE users; --",
		"<script>alert('xss')</script>",
		"%s%s%s%s%s%s%n%x%d",
		"| rm -rf / |",
		"$(cat /etc/passwd)",
		"../../../../../../windows/system32/cmd.exe",
	}

	for i, inj := range injections {
		subWorkdir := filepath.Join(workdir, fmt.Sprintf("inj_%d", i))
		req := &TaskRequest{
			ID:             fmt.Sprintf("task-inj-%d", i),
			RawRequest:     fmt.Sprintf("Build dashboard with %s", inj),
			WorkspaceRoot:  subWorkdir,
			Type:           "DESIGN",
			SkipVisualGate: true,
		}

		res, err := pipeline.Execute(context.Background(), req)
		if err != nil {
			t.Fatalf("Injection string '%s' caused unexpected execution failure: %v", inj, err)
		}
		if res.Status != PipelineStatusSuccess {
			t.Errorf("Expected SUCCESS, got %s for payload %s", res.Status, inj)
		}
	}
}

func TestAdversarial_MalformedTask_NilSafety(t *testing.T) {
	cat, graph, _ := setupEngineFixtures(t)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// Subtest 1: pipeline.Execute with nil request
	t.Run("Nil_TaskRequest_Execute", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("CONFIRMED: pipeline.Execute(ctx, nil) panicked as expected with nil TaskRequest: %v", r)
			}
		}()
		_, _ = pipeline.Execute(context.Background(), nil)
	})

	// Subtest 2: pipeline.Plan with nil request
	t.Run("Nil_TaskRequest_Plan", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("CONFIRMED: pipeline.Plan(ctx, nil) panicked as expected with nil TaskRequest: %v", r)
			}
		}()
		_, _ = pipeline.Plan(context.Background(), nil)
	})

	// Subtest 3: pipeline.Execute with nil context
	t.Run("Nil_Context_Execute", func(t *testing.T) {
		req := &TaskRequest{
			ID:            "task-nil-ctx",
			RawRequest:    "Standard request",
			WorkspaceRoot: t.TempDir(),
		}
		defer func() {
			if r := recover(); r != nil {
				t.Logf("CONFIRMED: pipeline.Execute(nil, req) panicked with nil Context: %v", r)
			}
		}()
		res, err := pipeline.Execute(nil, req)
		t.Logf("Execute with nil context completed: res=%v, err=%v", res, err)
	})
}

func TestAdversarial_MalformedTask_QuarantineBreachAttempt(t *testing.T) {
	cat, graph, _ := setupEngineFixtures(t)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	bannedPaths := []string{
		`C:\Users\mockuser\.gemini\config\skills_library`,
		`C:\Users\mockuser\.gemini\config\skills_library\subfolder`,
		`workspace/skills_library/test`,
		`relative/path/with/skills-library/breach`,
	}

	for _, banned := range bannedPaths {
		req := &TaskRequest{
			ID:            "task-quarantine-attack",
			RawRequest:    "Build component using quarantined files",
			WorkspaceRoot: banned,
			Type:          "DESIGN",
		}

		res, err := pipeline.Execute(context.Background(), req)
		if err == nil {
			t.Fatalf("SECURITY VIOLATION: Quarantine path '%s' was NOT blocked!", banned)
		}
		if !strings.Contains(err.Error(), "quarantine") && !strings.Contains(err.Error(), "strictly forbidden") {
			t.Errorf("Expected quarantine error message, got: %v", err)
		}
		if res.Status != PipelineStatusFailed {
			t.Errorf("Expected status FAILED for quarantined path, got: %s", res.Status)
		}
	}
}

// -----------------------------------------------------------------------------
// Focus Area 4: Multi-Viewport Failure Detection & Oscillation
// -----------------------------------------------------------------------------

func TestAdversarial_MultiViewport_OscillationStage5And6(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	// Oscillating verifier alternates between LAYOUT_CODE and TOKEN_STYLE indefinitely
	oscVerifier := NewOscillatingVerifier(0) // never heals

	pipeline := NewDesignPipeline(cat, graph, nil, oscVerifier)

	req := &TaskRequest{
		ID:             "task-oscillation-test",
		RawRequest:     "Build responsive interactive showcase with animations",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  4, // Should cycle 4 times then halt cleanly
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrMaxIterationsExceeded {
		t.Fatalf("Expected ErrMaxIterationsExceeded after 4 oscillations, got: %v", err)
	}

	if res.Status != PipelineStatusFailed {
		t.Errorf("Expected status FAILED, got %s", res.Status)
	}

	if res.IterationCount != 4 {
		t.Errorf("Expected exactly 4 iterations executed, got %d", res.IterationCount)
	}

	// Verify that the pipeline jumped back and forth between Stage 5 and Stage 6
	var designSysCount, implementCount, visualQACount, iterateCount int
	for _, stage := range res.StagesExecuted {
		switch stage {
		case StageDesignSystem:
			designSysCount++
		case StageImplement:
			implementCount++
		case StageVisualQA:
			visualQACount++
		case StageIterate:
			iterateCount++
		}
	}

	t.Logf("Oscillation telemetry: DesignSystem=%d, Implement=%d, VisualQA=%d, Iterate=%d",
		designSysCount, implementCount, visualQACount, iterateCount)

	if iterateCount != 5 { // Initial + 4 iteration re-evaluations
		t.Errorf("Expected 5 Iterate stage evaluations, got %d", iterateCount)
	}

	if visualQACount != 5 {
		t.Errorf("Expected 5 VisualQA evaluations, got %d", visualQACount)
	}

	// Both Stage 5 (DesignSystem) and Stage 6 (Implement) must have been re-executed due to oscillation!
	if designSysCount <= 1 {
		t.Errorf("Expected StageDesignSystem to be re-executed (>1), got %d", designSysCount)
	}
	if implementCount <= 1 {
		t.Errorf("Expected StageImplement to be re-executed (>1), got %d", implementCount)
	}

	// Check final state.json recovery point
	state, err := handoff.ReadState(workdir)
	if err != nil || state == nil || state.FailureRecovery == nil {
		t.Fatalf("Failed to read valid RecoveryPoint from handoff state after oscillation")
	}

	t.Logf("RecoveryPoint after oscillation: FailedStep='%s', ResumeFromStep='%s', ErrorReason='%s'",
		state.FailureRecovery.FailedStep, state.FailureRecovery.ResumeFromStep, state.FailureRecovery.ErrorReason)

	if state.FailureRecovery.FailedStep != "Visual QA" {
		t.Errorf("Expected FailedStep='Visual QA', got %s", state.FailureRecovery.FailedStep)
	}
}

func TestAdversarial_MultiViewport_OscillationWithEventualHealing(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	// Oscillates twice (fails call 0 and 1), then heals on call 2 (iteration 2)
	oscVerifier := NewOscillatingVerifier(2)

	pipeline := NewDesignPipeline(cat, graph, nil, oscVerifier)

	req := &TaskRequest{
		ID:             "task-oscillation-healing",
		RawRequest:     "Build responsive showcase that self-heals after 2 iterations",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  4,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Self-healing oscillation failed: %v", err)
	}

	if res.Status != PipelineStatusSuccess {
		t.Errorf("Expected SUCCESS after self-healing oscillation, got %s", res.Status)
	}

	if res.IterationCount != 2 {
		t.Errorf("Expected IterationCount=2, got %d", res.IterationCount)
	}
}

// -----------------------------------------------------------------------------
// Focus Area 5: Determinism, Concurrency & Memory Leak Verification
// -----------------------------------------------------------------------------

func TestAdversarial_Determinism_50Runs(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	var baselineArchetype string
	var baselineTokens float64
	var baselineRoutes []string

	for i := 0; i < 50; i++ {
		subWorkdir := filepath.Join(workdir, fmt.Sprintf("det_%d", i))
		req := &TaskRequest{
			ID:             fmt.Sprintf("task-det-%d", i),
			RawRequest:     "Award-winning 3D WebGL spatial portfolio with GSAP animations",
			WorkspaceRoot:  subWorkdir,
			Type:           "DESIGN",
			SkipVisualGate: true,
		}

		res, err := pipeline.Execute(context.Background(), req)
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}

		if i == 0 {
			baselineArchetype = res.Archetype
			for _, r := range res.ResolvedRoutes {
				baselineRoutes = append(baselineRoutes, r.CapabilityID)
			}
			plan, _ := pipeline.Plan(context.Background(), req)
			baselineTokens = plan.EstimatedTokenCost
		} else {
			if res.Archetype != baselineArchetype {
				t.Fatalf("Determinism broken at run %d: expected archetype %s, got %s", i, baselineArchetype, res.Archetype)
			}
			var currentRoutes []string
			for _, r := range res.ResolvedRoutes {
				currentRoutes = append(currentRoutes, r.CapabilityID)
			}
			if len(currentRoutes) != len(baselineRoutes) {
				t.Fatalf("Determinism broken at run %d: routes count mismatch (%d vs %d)", i, len(currentRoutes), len(baselineRoutes))
			}
			plan, _ := pipeline.Plan(context.Background(), req)
			if plan.EstimatedTokenCost != baselineTokens {
				t.Fatalf("Determinism broken at run %d: token cost mismatch (%.0f vs %.0f)", i, plan.EstimatedTokenCost, baselineTokens)
			}
		}
	}
}

func TestAdversarial_Concurrency_20Pipelines(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	concurrency := 20
	var wg sync.WaitGroup
	errCh := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			subDir := filepath.Join(workdir, fmt.Sprintf("conc_%d", idx))
			req := &TaskRequest{
				ID:             fmt.Sprintf("task-conc-%d", idx),
				RawRequest:     "Concurrent e-commerce dashboard with charts",
				WorkspaceRoot:  subDir,
				Type:           "DESIGN",
				SkipVisualGate: true,
			}

			res, err := pipeline.Execute(context.Background(), req)
			if err != nil {
				errCh <- fmt.Errorf("goroutine %d failed: %w", idx, err)
				return
			}
			if res.Status != PipelineStatusSuccess {
				errCh <- fmt.Errorf("goroutine %d expected SUCCESS, got %s", idx, res.Status)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("Concurrent execution error: %v", err)
	}
}

func TestAdversarial_MemoryLeak_100Cycles(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	// Warmup
	for i := 0; i < 5; i++ {
		subDir := filepath.Join(workdir, fmt.Sprintf("warmup_%d", i))
		req := &TaskRequest{
			ID:             "warmup",
			RawRequest:     "Warmup task",
			WorkspaceRoot:  subDir,
			SkipVisualGate: true,
		}
		_, _ = pipeline.Execute(context.Background(), req)
	}

	runtime.GC()
	var mBefore runtime.MemStats
	runtime.ReadMemStats(&mBefore)

	// Run 100 pipeline cycles
	for i := 0; i < 100; i++ {
		subDir := filepath.Join(workdir, fmt.Sprintf("mem_%d", i))
		req := &TaskRequest{
			ID:             fmt.Sprintf("task-mem-%d", i),
			RawRequest:     "Memory profiling showcase",
			WorkspaceRoot:  subDir,
			Type:           "DESIGN",
			SkipVisualGate: true,
		}
		res, err := pipeline.Execute(context.Background(), req)
		if err != nil || res.Status != PipelineStatusSuccess {
			t.Fatalf("Memory cycle %d failed", i)
		}
	}

	runtime.GC()
	var mAfter runtime.MemStats
	runtime.ReadMemStats(&mAfter)

	heapGrowthBytes := int64(mAfter.HeapAlloc) - int64(mBefore.HeapAlloc)
	t.Logf("Memory Stats: HeapAlloc Before=%d KB, After=%d KB, Growth=%d KB, TotalAlloc=%d KB, NumGC=%d",
		mBefore.HeapAlloc/1024, mAfter.HeapAlloc/1024, heapGrowthBytes/1024, mAfter.TotalAlloc/1024, mAfter.NumGC)

	// Ensure heap growth after 100 executions is reasonably bounded (< 25 MB)
	if heapGrowthBytes > 25*1024*1024 {
		t.Errorf("Excessive heap growth detected: %d bytes (> 25MB)", heapGrowthBytes)
	}
}

// -----------------------------------------------------------------------------
// Focus Area 6: State JSON Integrity & Schema Conformity
// -----------------------------------------------------------------------------

func TestAdversarial_HandoffState_JSONIntegrity(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)

	req := &TaskRequest{
		ID:             "task-json-integrity",
		RawRequest:     "Build spatial audio visualizer",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	statePath := filepath.Join(workdir, ".orchestra", "handoff", "state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("Could not read state.json: %v", err)
	}

	var stateMap map[string]any
	if err := json.Unmarshal(data, &stateMap); err != nil {
		t.Fatalf("state.json is malformed JSON: %v", err)
	}

	requiredKeys := []string{"session_id", "version", "source_agent", "target_agent", "active_tasks", "changed_files", "completed_steps", "timestamp"}
	for _, k := range requiredKeys {
		if _, ok := stateMap[k]; !ok {
			t.Errorf("Missing required key '%s' in state.json", k)
		}
	}

	if stateMap["session_id"] != req.ID {
		t.Errorf("Expected session_id=%s, got %v", req.ID, stateMap["session_id"])
	}
	if fmt.Sprintf("%v", stateMap["version"]) != "3" {
		t.Errorf("Expected version=3, got %v", stateMap["version"])
	}
	_ = res
}

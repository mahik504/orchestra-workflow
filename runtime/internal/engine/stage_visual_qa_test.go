package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/orchestra-v3/internal/memory"
)

type mockVisualVerifier struct {
	result *VisualQAData
	err    error
}

func (m *mockVisualVerifier) Name() string {
	return "mock-verifier"
}

func (m *mockVisualVerifier) Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func TestVisualQAStage_NonVisualTask(t *testing.T) {
	stage := NewVisualQAStage()
	taskCtx := &TaskContext{
		Task: &TaskRequest{
			ID:            "test-non-visual",
			RawRequest:    "Create backend API route",
			WorkspaceRoot: t.TempDir(),
		},
		Classification: &ClassificationData{
			RequiresVisual: false,
		},
		StageResults: make(map[StageName]*StageResult),
	}

	res, err := stage.Execute(taskCtx)
	if err != nil {
		t.Fatalf("VisualQAStage failed on non-visual task: %v", err)
	}
	if res.Status != StatusCompleted {
		t.Errorf("Expected StatusCompleted, got %s", res.Status)
	}
	if taskCtx.VisualQA == nil || !taskCtx.VisualQA.AllPassed {
		t.Errorf("Expected VisualQA to pass for non-visual task")
	}
}

func TestVisualQAStage_VisualTaskWithMemoryRecording(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "memory", "resource-memory.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", memFile)

	mockVerifier := &mockVisualVerifier{
		result: &VisualQAData{
			AllPassed: true,
			ViewportResults: []ViewportCheckResult{
				{ViewportName: "desktop", Passed: true, ScreenshotPath: filepath.Join(tmpDir, "desktop.png")},
				{ViewportName: "mobile", Passed: true, ScreenshotPath: filepath.Join(tmpDir, "mobile.png")},
			},
			Metrics:            map[string]float64{"visual_score": 1.0},
			DetectedViolations: []string{},
			FailureClass:       FailureClassNone,
		},
	}

	stage := NewVisualQAStage()
	taskCtx := &TaskContext{
		Ctx: context.Background(),
		Task: &TaskRequest{
			ID:            "test-visual-task",
			RawRequest:    "Create luxury landing page with hero cards",
			WorkspaceRoot: tmpDir,
		},
		Classification: &ClassificationData{
			RequiresVisual: true,
		},
		Verifier:      mockVerifier,
		StageResults:  make(map[StageName]*StageResult),
		ArtifactPaths: make(map[string]string),
	}

	res, err := stage.Execute(taskCtx)
	if err != nil {
		t.Fatalf("VisualQAStage failed: %v", err)
	}
	if res.Status != StatusCompleted {
		t.Errorf("Expected completed status, got: %s", res.Status)
	}

	// Verify memory file was created and contains playwright telemetry
	store, err := memory.NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to open recorded memory: %v", err)
	}

	agg, found := store.GetAggregate("playwright")
	if !found {
		t.Fatalf("Expected playwright aggregate in memory store")
	}
	if agg.TotalEvaluations != 1 || agg.SuccessCount != 1 {
		t.Errorf("Unexpected playwright metrics: %+v", agg)
	}
}

func TestIterateStage_RecordsMemoryOnVerdict(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "memory", "resource-memory.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", memFile)

	stage := NewIterateStage()
	taskCtx := &TaskContext{
		Ctx: context.Background(),
		Task: &TaskRequest{
			ID:            "test-iterate-task",
			RawRequest:    "Interactive 3D WebGL showcase",
			WorkspaceRoot: tmpDir,
		},
		Classification: &ClassificationData{
			RequiresVisual: true,
		},
		VisualQA: &VisualQAData{
			AllPassed: true,
		},
		Implementation: &ImplementationData{
			AcquiredResources: []string{"three", "react-bits"},
		},
		Iteration: &IterationData{
			CurrentIteration: 1,
			MaxIterations:    3,
		},
		StageResults: make(map[StageName]*StageResult),
	}

	res, err := stage.Execute(taskCtx)
	if err != nil {
		t.Fatalf("IterateStage failed: %v", err)
	}
	if res.Status != StatusCompleted {
		t.Errorf("Expected StatusCompleted, got %s", res.Status)
	}
	if taskCtx.Iteration.FinalVerdict != "PASSED" {
		t.Errorf("Expected verdict PASSED, got %s", taskCtx.Iteration.FinalVerdict)
	}

	// Verify acquired resources were recorded in memory
	store, err := memory.NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to load memory: %v", err)
	}

	for _, resID := range []string{"three", "react-bits"} {
		agg, found := store.GetAggregate(resID)
		if !found {
			t.Errorf("Expected aggregate for %s", resID)
		} else if agg.SuccessCount != 1 {
			t.Errorf("Expected 1 success for %s, got %d", resID, agg.SuccessCount)
		}
	}
}

func TestPipeline_RecordPipelineMemory(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "pipeline-memory.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", memFile)

	pipeline := &DesignPipeline{}
	taskCtx := &TaskContext{
		Task: &TaskRequest{
			ID:            "task-pipeline-test",
			RawRequest:    "Awwwards portfolio",
			WorkspaceRoot: tmpDir,
		},
		VisualQA: &VisualQAData{
			AllPassed:          false,
			DetectedViolations: []string{"horizontal overflow on mobile"},
			FailureClass:       FailureClassLayoutCode,
		},
		StageResults: make(map[StageName]*StageResult),
	}
	taskCtx.SetStageResult(StageVisualQA, &StageResult{
		StageName: StageVisualQA,
		Duration:  350 * time.Millisecond,
	})

	execRes := &DesignExecutionResult{
		TaskID:            "task-pipeline-test",
		Status:            PipelineStatusSuccess,
		Archetype:         "premium-website",
		TotalDuration:     1200 * time.Millisecond,
		AcquiredResources: []string{"gsap"},
	}

	err := pipeline.RecordPipelineMemory(taskCtx, execRes)
	if err != nil {
		t.Fatalf("RecordPipelineMemory failed: %v", err)
	}

	store, err := memory.NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to load recorded store: %v", err)
	}

	if _, ok := store.GetAggregate("playwright"); ok {
		t.Fatal("non-visual VisualQA must not record playwright")
	}

	gsapAgg, ok := store.GetAggregate("gsap")
	if !ok {
		t.Fatalf("Expected gsap aggregate")
	}
	if gsapAgg.SuccessCount != 1 || gsapAgg.LastOutcome != memory.OutcomeSuccess {
		t.Errorf("Expected gsap success, got: %+v", gsapAgg)
	}
}

func TestPipeline_RecordPipelineMemory_AcquiredSuccessOnFailedPipeline(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "pipeline-acquired.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", memFile)

	pipeline := &DesignPipeline{}
	taskCtx := &TaskContext{
		Task: &TaskRequest{
			ID:            "task-drei-integrity",
			RawRequest:    "Install the drei helper library into this Node fixture so the package resolves",
			WorkspaceRoot: tmpDir,
		},
		Classification: &ClassificationData{
			OverlayActivations: []string{"drei"},
		},
		Implementation: &ImplementationData{
			AcquiredResources: []string{"drei"},
			InstalledPaths:    map[string]string{"drei": filepath.Join(tmpDir, "node_modules", "@react-three", "drei")},
			AcquireCommands:   map[string]string{"drei": "pnpm add @react-three/drei"},
		},
		StageResults: make(map[StageName]*StageResult),
	}

	err := pipeline.RecordPipelineMemory(taskCtx, &DesignExecutionResult{
		TaskID:            "task-drei-integrity",
		Status:            PipelineStatusFailed,
		Archetype:         "standard-feature",
		AcquiredResources: []string{"drei"},
		TotalDuration:     800 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("RecordPipelineMemory failed: %v", err)
	}

	store, err := memory.NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatal(err)
	}
	agg, ok := store.GetAggregate("drei")
	if !ok {
		t.Fatal("expected drei aggregate")
	}
	if agg.LastOutcome != memory.OutcomeSuccess || agg.SuccessCount != 1 {
		t.Fatalf("acquired install path must record success even when later pipeline stages fail, got %+v", agg)
	}
}

func TestPipeline_RecordPipelineMemory_PlaywrightOnlyWhenVerifierRan(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "pipeline-memory-visual.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", memFile)

	pipeline := &DesignPipeline{}
	taskCtx := &TaskContext{
		Task: &TaskRequest{
			ID:            "task-visual-mem",
			RawRequest:    "Awwwards portfolio",
			WorkspaceRoot: tmpDir,
		},
		Classification: &ClassificationData{RequiresVisual: true},
		VisualQA: &VisualQAData{
			AllPassed:          false,
			DetectedViolations: []string{"horizontal overflow on mobile"},
			FailureClass:       FailureClassLayoutCode,
			VerifierRan:        true,
		},
		StageResults: make(map[StageName]*StageResult),
	}
	taskCtx.SetStageResult(StageVisualQA, &StageResult{
		StageName: StageVisualQA,
		Duration:  350 * time.Millisecond,
	})
	if err := pipeline.RecordPipelineMemory(taskCtx, &DesignExecutionResult{
		TaskID: "task-visual-mem", Status: PipelineStatusSuccess,
	}); err != nil {
		t.Fatal(err)
	}
	store, err := memory.NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatal(err)
	}
	pwAgg, ok := store.GetAggregate("playwright")
	if !ok || pwAgg.FailureCount != 1 {
		t.Fatalf("expected playwright failure when verifier ran, got %+v ok=%v", pwAgg, ok)
	}
}

func init() {
	// Clean env during unit test if needed
	_ = os.Setenv("ORCHESTRA_ISOLATION_MODE", "test")
}

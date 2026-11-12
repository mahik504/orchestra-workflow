package engine

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func TestIterateNotApplicableDoesNotRecordSuccess(t *testing.T) {
	dir := t.TempDir()
	mem := filepath.Join(dir, "resource-memory.json")
	t.Setenv("ORCHESTRA_MEMORY_PATH", mem)
	ctx := &TaskContext{
		Task:           &TaskRequest{ID: "skip", WorkspaceRoot: dir, RawRequest: "secret-prompt-do-not-store"},
		VisualQA:       &VisualQAData{Verdict: CheckNotApplicable},
		Implementation: &ImplementationData{AcquiredResources: []string{"playwright"}},
		Iteration:      &IterationData{},
	}
	if _, err := NewIterateStage().Execute(ctx); err != nil {
		t.Fatal(err)
	}
	if ctx.Iteration.FinalVerdict != CheckNotApplicable {
		t.Fatalf("verdict %s", ctx.Iteration.FinalVerdict)
	}
	if _, err := os.Stat(mem); err == nil {
		raw, _ := os.ReadFile(mem)
		if strings.Contains(string(raw), "success") || strings.Contains(string(raw), "secret-prompt") {
			t.Fatalf("recorded a skip: %s", raw)
		}
	}
}

func TestPlanStageNamesOmitSkippedResearch(t *testing.T) {
	ctx := &TaskContext{StageResults: map[StageName]*StageResult{}}
	if got := planStageNames(ctx); got != "Discover -> Classify -> Synthesize" {
		t.Fatal(got)
	}
	ctx.SetStageResult(StageResearch, &StageResult{StageName: StageResearch, Status: StatusCompleted})
	if got := planStageNames(ctx); got != "Discover -> Classify -> Research -> Synthesize" {
		t.Fatal(got)
	}
}

type recordingRunner struct {
	got visual.LayerARequest
}

func (r *recordingRunner) Name() string { return "recording" }

func (r *recordingRunner) Run(_ context.Context, req visual.LayerARequest) (*visual.LayerAReport, error) {
	r.got = req
	return &visual.LayerAReport{Owner: "recording", Passed: true}, nil
}

func TestVerifierUsesSelectedViewports(t *testing.T) {
	runner := &recordingRunner{}
	verifier := NewPlaywrightVerifier(t.TempDir())
	verifier.Runner = runner
	selected, err := visual.SelectViewports("mobile")
	if err != nil {
		t.Fatal(err)
	}
	verifier.Viewports = selected
	_, err = verifier.Verify(context.Background(), &TaskContext{
		Task:    &TaskRequest{WorkspaceRoot: t.TempDir(), RawRequest: "landing page"},
		Mission: &mission.Mission{VisualImportance: mission.VisualProduct, Objective: "landing page"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.got.Viewports) != 1 || runner.got.Viewports[0].Name != "mobile" {
		t.Fatalf("viewports %+v", runner.got.Viewports)
	}
}

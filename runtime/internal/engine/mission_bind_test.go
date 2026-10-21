package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

func TestPipeline_MissionStateAndCheapBudget(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)
	req := &TaskRequest{
		ID:            "task-3c-api",
		RawRequest:    "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir,
		Type:          "BUGFIX",
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("pipeline: %v", err)
	}
	if res.MissionPath == "" || res.RunStatePath == "" {
		t.Fatalf("missing mission/state paths: %+v", res)
	}
	m, err := mission.ReadMissionFile(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mission.Validate(m); err != nil {
		t.Fatal(err)
	}
	if m.Type != mission.TypeDebug {
		t.Fatalf("type=%s", m.Type)
	}
	if !mission.CheapBudget(m.Budget) {
		t.Fatalf("debug mission should be cheap: %+v (%s)", m.Budget, m.BudgetDerivation)
	}
	if leak := mission.ProviderLeak([]byte(m.BudgetDerivation + string(m.Budget.Reasoning))); leak != "" {
		t.Fatalf("provider leak %q", leak)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Status != mission.StatusSucceeded {
		t.Fatalf("status %s", st.Status)
	}
	why, err := mission.ContinueWithoutConversation(st, m)
	if err == nil {
		t.Fatalf("terminal run should not continue, got %s", why)
	}

	restored, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if restored.CurrentStage == "" {
		t.Fatalf("current stage lost on restore")
	}
}

func TestPipeline_CapsuleCompiledAndScoped(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID:             "task-3c-capsule",
		RawRequest:     "Build mobile dashboard with graphs",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
	}
	if _, err := pipeline.Execute(context.Background(), req); err != nil {
		t.Fatalf("pipeline: %v", err)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st.CapsuleScope == nil {
		t.Fatal("capsule scope missing")
	}
	if st.CapsuleScope.CapsuleBytes <= 0 {
		t.Fatal("capsule bytes not measured")
	}
	m, err := mission.ReadMissionFile(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	cap, report, err := mission.Compile(mission.CompileRequest{
		Mission:             m,
		Role:                mission.RoleBuilder,
		Files:               []string{"DESIGN.md", "README.md"},
		UnrelatedFiles:      []string{"README.md"},
		AllSkills:           []string{"hallmark", "gsap-core"},
		Skills:              m.Capabilities,
		WholeProjectContext: []byte("AGENTS.md\n" + string(make([]byte, 8000))),
		ConversationHistory: []string{"old chat"},
		MCPCatalog:          []string{"playwright", "context7"},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range cap.Files {
		if f == "README.md" {
			t.Fatal("unrelated file in capsule")
		}
	}
	if report.Ratio >= 1 && report.WholeProjectBytes > 0 {
		t.Fatalf("capsule not smaller than whole-project dump")
	}
}

func TestPipeline_HumanGateStopsAndRecordsState(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)
	req := &TaskRequest{
		ID:             "task-3c-gate",
		RawRequest:     "Design interactive landing page with bold typography",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: false,
	}
	_, err := pipeline.Execute(context.Background(), req)
	if err != ErrHumanGateRequired {
		t.Fatalf("expected ErrHumanGateRequired, got %v", err)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st.PendingHuman == nil || !st.PendingHuman.Blocking {
		t.Fatalf("expected blocking human decision, got %+v", st.PendingHuman)
	}
	if ok, _ := st.CanContinue(); ok {
		t.Fatal("automatic execution continued past human gate")
	}
}

func TestPipeline_RetryHistoryRecordedOnVisualLoop(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.SimulateMobileOverflow = true
	mockVerifier.FailUntilIteration = 1
	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)
	req := &TaskRequest{
		ID:             "task-3c-retry-hist",
		RawRequest:     "Build mobile dashboard with graphs",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  3,
	}
	if _, err := pipeline.Execute(context.Background(), req); err != nil {
		t.Fatalf("pipeline: %v", err)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.FailureHistory) == 0 {
		t.Fatal("expected retry history on visual self-heal")
	}
	if st.FailureHistory[0].FailureClass == "" {
		t.Fatal("failure class missing")
	}
	if _, err := os.Stat(filepath.Join(mission.RunDir(workdir, req.ID), "state.json")); err != nil {
		t.Fatal(err)
	}
}

func TestPipeline_HealedIterateDoesNotLeaveHumanBlock(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.SimulateMobileOverflow = true
	mockVerifier.FailUntilIteration = 2
	pipeline := NewDesignPipeline(cat, graph, nil, mockVerifier)
	req := &TaskRequest{
		ID:             "task-3c1-heal-no-block",
		RawRequest:     "Build mobile dashboard with graphs",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  3,
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("pipeline: %v", err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("status %s", res.Status)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st.PendingHuman != nil {
		t.Fatalf("healed iterate left pending_human: %+v", st.PendingHuman)
	}
	if st.Status != mission.StatusSucceeded {
		t.Fatalf("status %s", st.Status)
	}
	if len(st.FailureHistory) < 2 {
		t.Fatalf("expected history for two QA misses, got %d", len(st.FailureHistory))
	}
}

func TestPipeline_SentinelsUnchangedAfterBind(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, nil)
	gated, err := pipeline.Execute(context.Background(), &TaskRequest{
		ID:             "task-3c1-sentinel-gate",
		RawRequest:     "Design interactive landing page with bold typography",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: false,
	})
	if err != ErrHumanGateRequired {
		t.Fatalf("ErrHumanGateRequired identity changed: %v", err)
	}
	if gated.Status != PipelineStatusGated {
		t.Fatalf("status %s", gated.Status)
	}

	mockVerifier := NewMockVisualVerifier()
	mockVerifier.AlwaysFail = true
	mockVerifier.SimulateMobileOverflow = true
	failPipe := NewDesignPipeline(cat, graph, nil, mockVerifier)
	failPipe.MaxIterations = 2
	exhausted, err := failPipe.Execute(context.Background(), &TaskRequest{
		ID:             "task-3c1-sentinel-max",
		RawRequest:     "Build layout with permanent overflow defect",
		WorkspaceRoot:  workdir,
		Type:           "DESIGN",
		SkipVisualGate: true,
		MaxIterations:  2,
	})
	if err != ErrMaxIterationsExceeded {
		t.Fatalf("ErrMaxIterationsExceeded identity changed: %v", err)
	}
	if exhausted.Status != PipelineStatusFailed {
		t.Fatalf("status %s", exhausted.Status)
	}
}

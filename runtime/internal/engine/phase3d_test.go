package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func TestFlowA_APIDebugStaysCheap(t *testing.T) {
	visual.ResetScreenshotEvaluatorForTest()
	defer visual.ResetScreenshotEvaluatorForTest()
	visual.RegisterScreenshotEvaluator(nopShot{})
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewPlaywrightVerifier(filepath.Join(workdir, "qa")))
	req := &TaskRequest{
		ID: "flow-a-api", RawRequest: "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir, Type: "BUGFIX",
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("status %s", res.Status)
	}
	m, err := mission.ReadMissionFile(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if m.Budget.VisualReview != mission.VisualReviewNone {
		t.Fatalf("visual_review=%s", m.Budget.VisualReview)
	}
	if m.VisualImportance == mission.VisualFlagship {
		t.Fatal("API debug must not be flagship")
	}
	if !mission.CheapBudget(m.Budget) {
		t.Fatalf("must stay cheap: %+v", m.Budget)
	}
	if visual.ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("ordinary API must not invoke screenshot eval")
	}
	st, _ := mission.ReadRunState(workdir, req.ID)
	for _, o := range st.Outputs {
		if filepath.Ext(o.Path) == ".png" {
			t.Fatalf("ordinary API must not write screenshots: %+v", o)
		}
	}
}

func TestFlowB_CRUDUINoDirectionsNoJudge(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "flow-b", Objective: "Build a settings screen for profile CRUD",
		DeclaredType: "FEATURE", QualityBar: mission.BarPremium,
	})
	if m.Type != mission.TypeProductUI {
		t.Fatalf("type %s", m.Type)
	}
	if m.VisualImportance != mission.VisualProduct {
		t.Fatalf("importance %s", m.VisualImportance)
	}
	if m.Budget.VisualReview != mission.VisualReviewBrowser {
		t.Fatalf("review %s", m.Budget.VisualReview)
	}
	if m.Budget.IndependentJudge == mission.JudgeRequired {
		t.Fatal("CRUD UI must not require independent judge")
	}
	if err := visual.ValidateExploration(visual.ClassProductUI, visual.DefaultFlagshipSet(m.Objective)); err == nil {
		t.Fatal("CRUD UI must not generate 3-5 directions")
	}
}

func TestFlowC_ResponsiveRunningApp(t *testing.T) {
	requireLivePlaywright(t)
	html, err := os.ReadFile(filepath.Join("..", "visual", "testdata", "fixture.html"))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(html)
	}))
	defer srv.Close()
	dir := t.TempDir()
	ctxTO, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	v := &PlaywrightVerifier{OutputDir: dir, Runner: visual.NewPlaywrightRunner()}
	ctx := &TaskContext{
		Ctx:            ctxTO,
		Task:           &TaskRequest{ID: "flow-c", RawRequest: "Responsive settings page", WorkspaceRoot: dir, Parameters: map[string]any{"app_url": srv.URL}},
		Mission:        mission.FromSeed(mission.Seed{ID: "flow-c", Objective: "Responsive settings page", DeclaredType: "FEATURE"}),
		Classification: &ClassificationData{RequiresVisual: true},
		RunState:       mission.NewRunState(mission.FromSeed(mission.Seed{ID: "flow-c", Objective: "Responsive settings page"})),
	}
	data, err := v.Verify(ctxTO, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if data.LayerA != nil && data.LayerA.PlaywrightMissing {
		t.Fatalf("BLOCKED: %s", data.LayerA.EnvironmentError)
	}
	if data.AllPassed {
		t.Fatal("mobile overflow must fail independently of desktop")
	}
	if data.CanonicalOwner != visual.CanonicalOwner {
		t.Fatalf("owner %s", data.CanonicalOwner)
	}
	var desk, mob bool
	for _, vp := range data.ViewportResults {
		if vp.ScreenshotPath != "" {
			b, rerr := os.ReadFile(vp.ScreenshotPath)
			if rerr != nil || !visual.IsPNG(b) || visual.IsMockScreenshot(b) || len(b) <= len(visual.TinyPNG) {
				t.Fatalf("flow-c %s is not a real PNG: path=%s err=%v size=%d", vp.ViewportName, vp.ScreenshotPath, rerr, len(b))
			}
			t.Logf("flow-c viewport=%s bytes=%d overflow=%v path=%s", vp.ViewportName, len(b), vp.HasOverflow, vp.ScreenshotPath)
		}
		if vp.ViewportName == "desktop" {
			desk = true
			if vp.HasOverflow {
				t.Fatal("desktop must not inherit mobile overflow")
			}
		}
		if vp.ViewportName == "mobile" {
			mob = true
			if !vp.HasOverflow {
				t.Fatal("mobile overflow missing")
			}
		}
	}
	if !desk || !mob {
		t.Fatalf("desktop+mobile evidence required: %+v", data.ViewportResults)
	}
	if data.LayerA != nil {
		t.Logf("flow-c browser=%s %s startup_ms=%d screenshot_ms=%d closed=%v",
			data.LayerA.BrowserName, data.LayerA.BrowserVersion, data.LayerA.StartupMS, data.LayerA.ScreenshotMS, data.LayerA.BrowserClosed)
	}
}

func TestFlowD_FlagshipDirectionsThenLayers(t *testing.T) {
	requireLivePlaywright(t)
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	html := []byte("<!doctype html><html><body><h1 data-qa=hero>Flagship</h1></body></html>")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(html)
	}))
	defer srv.Close()
	v := &PlaywrightVerifier{OutputDir: filepath.Join(workdir, "qa"), Runner: visual.NewPlaywrightRunner()}
	pipeline := NewDesignPipeline(cat, graph, nil, v)
	dirs := visual.DefaultFlagshipSet("portfolio")
	req := &TaskRequest{
		ID: "flow-d", RawRequest: "Award-winning 3D portfolio showcase",
		WorkspaceRoot: workdir, Type: "DESIGN", SkipVisualGate: true,
		Parameters: map[string]any{
			"app_url": srv.URL, "selected_direction_id": dirs[2].ID, "directions": dirs,
		},
	}
	ctxTO, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	res, err := pipeline.Execute(ctxTO, req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("status %s", res.Status)
	}
	// Reconstruct exploration from a dedicated bind check
	tc := NewTaskContext(context.Background(), req, cat, graph, pipeline.Router, pipeline.ResearchCoord, v)
	attachMission(tc)
	if err := bindFlagshipExploration(tc); err != nil {
		t.Fatal(err)
	}
	if tc.Exploration == nil || tc.Exploration.Sealed == nil || tc.Exploration.SelectedID != dirs[2].ID {
		t.Fatalf("sealed direction: %+v", tc.Exploration)
	}
	if visual.Mutated(tc.Exploration.Sealed, dirs[0]) == false {
		t.Fatal("contract must not silently accept another direction")
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	evPath := filepath.Join(mission.RunDir(workdir, req.ID), "visual-evidence.json")
	raw, err := os.ReadFile(evPath)
	if err != nil {
		t.Fatalf("durable visual evidence missing: %v", err)
	}
	if bytes.Contains(raw, []byte("MOCK_SCREENSHOT_")) {
		t.Fatal("durable evidence must not record mock screenshots")
	}
	t.Logf("flow-d visual-evidence bytes=%d outputs=%d", len(raw), len(st.Outputs))
	_ = st
}

func TestFlowE_VisualHealChangesImplementation(t *testing.T) {
	requireLivePlaywright(t)
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	overflowHTML, err := os.ReadFile(filepath.Join("..", "visual", "testdata", "fixture.html"))
	if err != nil {
		t.Fatal(err)
	}
	fixedHTML := []byte("<!doctype html><html><body><h1 data-qa=hero>Healed dashboard</h1></body></html>")
	var gen int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if atomic.LoadInt32(&gen) == 0 {
			_, _ = w.Write(overflowHTML)
			return
		}
		_, _ = w.Write(fixedHTML)
	}))
	defer srv.Close()
	inner := &PlaywrightVerifier{OutputDir: t.TempDir(), Runner: visual.NewPlaywrightRunner()}
	v := &flipOnFailVerifier{inner: inner, gen: &gen}
	pipeline := NewDesignPipeline(cat, graph, nil, v)
	ctxTO, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	res, err := pipeline.Execute(ctxTO, &TaskRequest{
		ID: "flow-e", RawRequest: "Build mobile dashboard with graphs",
		WorkspaceRoot: workdir, Type: "DESIGN", SkipVisualGate: true, MaxIterations: 3,
		Parameters: map[string]any{"app_url": srv.URL},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("%s", res.Status)
	}
	if res.IterationCount < 1 {
		t.Fatal("expected a changed retry")
	}
	st, _ := mission.ReadRunState(workdir, "flow-e")
	if st.PendingHuman != nil && st.PendingHuman.Blocking {
		t.Fatal("successful heal must not leave a blocking human gate")
	}
	t.Logf("flow-e iterations=%d gen=%d", res.IterationCount, atomic.LoadInt32(&gen))
}

type nopShot struct{}

func (nopShot) Evaluate(ctx context.Context, req visual.ScreenshotEvalRequest) (*visual.ScreenshotEvalResult, error) {
	return &visual.ScreenshotEvalResult{EvaluatorType: visual.EvaluatorScreenshot, Result: "pass"}, nil
}

type flipOnFailVerifier struct {
	inner *PlaywrightVerifier
	gen   *int32
}

func (f *flipOnFailVerifier) Name() string { return f.inner.Name() }

func (f *flipOnFailVerifier) Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error) {
	data, err := f.inner.Verify(ctx, taskCtx)
	if err == nil && data != nil && !data.AllPassed {
		atomic.StoreInt32(f.gen, 1)
	}
	return data, err
}

func requireLivePlaywright(t *testing.T) {
	t.Helper()
	probe := visual.ProbePlaywright()
	if !probe.Available {
		t.Skipf("Playwright library or Chromium is not installed: %s", probe.Error)
	}
}

func TestFlowF_IdenticalFailureStops(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	mock := NewMockVisualVerifier()
	mock.AlwaysFail = true
	mock.SimulateMobileOverflow = true
	pipeline := NewDesignPipeline(cat, graph, nil, mock)
	pipeline.MaxIterations = 2
	_, err := pipeline.Execute(context.Background(), &TaskRequest{
		ID: "flow-f", RawRequest: "Build layout with permanent overflow defect",
		WorkspaceRoot: workdir, Type: "DESIGN", SkipVisualGate: true, MaxIterations: 2,
	})
	if !errors.Is(err, ErrMaxIterationsExceeded) && err != ErrMaxIterationsExceeded {
		t.Fatalf("sentinel identity: %v", err)
	}
}

func TestPlaywrightVerifierFailClosedNoApp(t *testing.T) {
	v := NewPlaywrightVerifier(t.TempDir())
	data, err := v.Verify(context.Background(), &TaskContext{
		Task:           &TaskRequest{ID: "no-app", RawRequest: "Build a dashboard", WorkspaceRoot: t.TempDir()},
		Mission:        mission.FromSeed(mission.Seed{ID: "no-app", Objective: "Build a dashboard", DeclaredType: "DESIGN"}),
		Classification: &ClassificationData{RequiresVisual: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.AllPassed {
		t.Fatal("must not pass without a running app")
	}
	if data.CanonicalOwner != visual.CanonicalOwner {
		t.Fatalf("owner %s", data.CanonicalOwner)
	}
}

func TestJudgeCapsuleScoped(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "judge-scope", Objective: "Award-winning 3D portfolio showcase",
		DeclaredType: "DESIGN", QualityBar: mission.BarPremium,
	})
	ctx := &TaskContext{
		Mission:  m,
		Task:     &TaskRequest{ID: m.ID, RawRequest: m.Objective, WorkspaceRoot: t.TempDir()},
		VisualQA: &VisualQAData{Evidence: []visual.Evidence{{ScreenshotPath: "qa/desktop.png", Layer: "A", Viewport: "desktop"}}},
	}
	compileJudgeCapsule(ctx, ctx.VisualQA)
	if ctx.JudgeCapsule == nil || ctx.JudgeScope == nil {
		t.Fatal("judge capsule missing")
	}
	if ctx.JudgeCapsule.Role != mission.RoleJudge {
		t.Fatalf("role %s", ctx.JudgeCapsule.Role)
	}
	raw, _ := json.Marshal(ctx.JudgeCapsule)
	blob := string(raw)
	if strings.Contains(blob, "WHOLE PROJECT DUMP") || strings.Contains(blob, "full conversation") ||
		strings.Contains(blob, "playwright-mcp") || strings.Contains(blob, "hallmark") {
		t.Fatalf("judge capsule leaked context: %s", blob)
	}
	for _, tool := range ctx.JudgeCapsule.ToolsAllowed {
		if tool == "write_file" || tool == "shell" {
			t.Fatalf("write tool in judge: %v", ctx.JudgeCapsule.ToolsAllowed)
		}
	}
	t.Logf("judge_capsule_bytes=%d whole_project_bytes=%d ratio=%.4f excluded=%v",
		ctx.JudgeScope.CapsuleBytes, ctx.JudgeScope.WholeProjectBytes, ctx.JudgeScope.Ratio, ctx.JudgeScope.Excluded)
	if ctx.JudgeScope.WholeProjectBytes > 0 && ctx.JudgeScope.Ratio >= 1 {
		t.Fatalf("judge capsule not smaller than project dump")
	}
}

func TestDurableVisualEvidenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	m := mission.FromSeed(mission.Seed{ID: "ev", Objective: "Build a dashboard", DeclaredType: "DESIGN"})
	st := mission.NewRunState(m)
	ctx := &TaskContext{
		Task:     &TaskRequest{ID: m.ID, WorkspaceRoot: dir},
		Mission:  m,
		RunState: st,
		VisualQA: &VisualQAData{
			CanonicalOwner: visual.CanonicalOwner,
			Evidence: []visual.Evidence{{
				What: "layer-a-capture", Route: "/", Surface: "home", Viewport: "desktop",
				Layer: "A", Passed: true, ScreenshotPath: filepath.Join(dir, "desktop.png"),
			}},
			FinalResult: "pass",
		},
	}
	if err := os.WriteFile(filepath.Join(dir, "desktop.png"), visual.TinyPNG, 0644); err != nil {
		t.Fatal(err)
	}
	persistVisualEvidence(ctx)
	if _, err := mission.WriteRunState(dir, st); err != nil {
		t.Fatal(err)
	}
	got, err := mission.ReadRunState(dir, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Outputs) == 0 {
		t.Fatal("evidence must survive durable state")
	}
	if _, err := os.Stat(filepath.Join(mission.RunDir(dir, m.ID), "visual-evidence.json")); err != nil {
		t.Fatal(err)
	}
}

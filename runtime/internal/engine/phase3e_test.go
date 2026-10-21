package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/handoff"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func TestPhase3E_FlowA_APIDebugConsumed(t *testing.T) {
	visual.ResetScreenshotEvaluatorForTest()
	defer visual.ResetScreenshotEvaluatorForTest()
	visual.RegisterScreenshotEvaluator(nopShot{})
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewPlaywrightVerifier(t.TempDir()))
	req := &TaskRequest{
		ID: "p3e-a", RawRequest: "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir, Type: "BUGFIX",
		Parameters: map[string]any{"discovered_models": []string{"grok-4.7"}},
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("status %s", res.Status)
	}
	researchSkipped := false
	for _, sr := range res.StageResults {
		if sr.StageName == StageResearch && sr.Status == StatusSkipped {
			researchSkipped = true
		}
	}
	if !researchSkipped {
		t.Fatal("research.worker absent must skip research")
	}
	d, err := capability.ReadDecision(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.Has(capability.IDResearchWorker) || d.Has(capability.IDBrowserPlaywright) || d.Has(capability.IDJudgeIndependent) || d.Has(capability.IDDesignWorker) {
		t.Fatalf("ordinary extra caps: %+v", d.Selected)
	}
	if !d.Has(capability.IDCodingWorker) {
		t.Fatal("coding.worker missing")
	}
	m, err := mission.ReadMissionFile(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	rawM, _ := json.Marshal(m)
	if leak := mission.ProviderLeak(rawM); leak != "" {
		t.Fatalf("mission leak %s", leak)
	}
	if strings.Contains(strings.ToLower(string(rawM)), "grok-4.7") {
		t.Fatal("mission contained opaque model label")
	}
	ctx := &TaskContext{Task: req, Mission: m, Capability: d, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	if ctx.Capsule == nil || ctx.Capsule.Role != mission.RoleBuilder {
		t.Fatal("coding.worker not consumed")
	}
	rawC, _ := json.Marshal(ctx.Capsule)
	if leak := mission.ProviderLeak(rawC); leak != "" {
		t.Fatalf("capsule leak %s", leak)
	}
	if strings.Contains(string(rawC), "grok-4.7") {
		t.Fatal("capsule contained opaque model label")
	}
	if ctx.JudgeCapsule != nil {
		t.Fatal("judge capsule must be nil")
	}
	if ctx.Exploration != nil && len(ctx.Exploration.Candidates) > 0 {
		t.Fatal("ordinary must not generate directions")
	}
	if visual.ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("ordinary must not invoke screenshot eval")
	}
	if res.VisualQAPassed == false && len(res.Screenshots) > 0 {
		t.Fatal("ordinary wrote screenshots")
	}
	capBytes, _ := os.Stat(filepath.Join(mission.RunDir(workdir, req.ID), "capability.json"))
	t.Logf("flow-a capsule_bytes=%d capability_bytes=%d", len(rawC), capBytes.Size())
	if capBytes.Size() > 20000 {
		t.Fatalf("capability.json too large: %d", capBytes.Size())
	}
}

func TestPhase3E_FlowB_ProductUIConsumed(t *testing.T) {
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
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID: "p3e-b", RawRequest: "Build a settings screen for profile CRUD",
		WorkspaceRoot: workdir, Type: "FEATURE", SkipVisualGate: true, MaxIterations: 1,
		Parameters: map[string]any{
			"app_url": srv.URL,
			"capability_skills": []string{
				"superpowers-planning", "gsap-core", "impeccable", "taste-design", "react-best-practices",
			},
		},
	}
	res, err := pipeline.Execute(context.Background(), req)
	_ = res
	d, rerr := capability.ReadDecision(workdir, req.ID)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if !d.Has(capability.IDCodingWorker) || !d.Has(capability.IDBrowserPlaywright) || !d.Has(capability.IDVisualStructural) {
		t.Fatalf("product ui caps: %+v", d.Selected)
	}
	if d.Has(capability.IDDesignWorker) {
		t.Fatal("product UI must not select design.worker")
	}
	m, _ := mission.ReadMissionFile(workdir, req.ID)
	ctx := &TaskContext{Task: req, Mission: m, Capability: d, RunState: mission.NewRunState(m), Ctx: context.Background()}
	compileBuilderCapsule(ctx)
	if ctx.Capsule == nil || ctx.Capsule.Role != mission.RoleBuilder {
		t.Fatal("coding.worker not consumed by builder capsule")
	}
	if err := bindFlagshipExploration(ctx); err != nil {
		t.Fatal(err)
	}
	if ctx.Exploration != nil {
		t.Fatalf("product UI exploration must be nil: %+v", ctx.Exploration)
	}
	ctxTO, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	ctx.Ctx = ctxTO
	ctx.Classification = &ClassificationData{RequiresVisual: true}
	v := NewPlaywrightVerifier(t.TempDir())
	data, verr := v.Verify(ctxTO, ctx)
	if verr != nil {
		t.Fatal(verr)
	}
	if data.LayerA != nil && data.LayerA.PlaywrightMissing {
		t.Fatalf("BLOCKED: %s", data.LayerA.EnvironmentError)
	}
	if len(data.ViewportResults) == 0 {
		t.Fatal("Layer A must capture viewports")
	}
	var desk, mob bool
	for _, vp := range data.ViewportResults {
		if vp.ScreenshotPath == "" {
			continue
		}
		b, rerr := os.ReadFile(vp.ScreenshotPath)
		if rerr != nil || !visual.IsPNG(b) || visual.IsMockScreenshot(b) {
			t.Fatalf("not real PNG %s", vp.ViewportName)
		}
		t.Logf("flow-b viewport=%s bytes=%d", vp.ViewportName, len(b))
		if vp.ViewportName == "desktop" {
			desk = true
		}
		if vp.ViewportName == "mobile" {
			mob = true
		}
	}
	if !desk || !mob {
		t.Fatal("desktop+mobile Chromium captures required")
	}
	if data.LayerB == nil || data.LayerB.EvaluatorType == visual.EvaluatorNone {
		t.Fatal("structural Layer B must run")
	}
	if data.LayerB.EvaluatorType == "multimodal" {
		t.Fatal("evaluator_type must not be multimodal")
	}
	if ctx.JudgeCapsule != nil {
		t.Fatal("product UI judge capsule must stay nil unless selected")
	}
	rawC, _ := json.Marshal(ctx.Capsule)
	t.Logf("flow-b capsule_bytes=%d err=%v status=%v", len(rawC), err, pipelineStatus(res))
}

func pipelineStatus(res *DesignExecutionResult) string {
	if res == nil {
		return ""
	}
	return string(res.Status)
}

func TestPhase3E_FlowC_FlagshipConsumed(t *testing.T) {
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
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID: "p3e-c", RawRequest: "Build a flagship portfolio landing with unusual motion",
		WorkspaceRoot: workdir, Type: "DESIGN", SkipVisualGate: true,
		Parameters: map[string]any{"app_url": srv.URL},
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != PipelineStatusSuccess {
		t.Fatalf("status %s", res.Status)
	}
	d, err := capability.ReadDecision(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !d.Has(capability.IDDesignWorker) || !d.Has(capability.IDCodingWorker) {
		t.Fatalf("flagship caps %+v", d.Selected)
	}
	m, _ := mission.ReadMissionFile(workdir, req.ID)
	researchRan := false
	researchSkipped := false
	for _, sr := range res.StageResults {
		if sr.StageName == StageResearch && sr.Status == StatusCompleted {
			researchRan = true
		}
		if sr.StageName == StageResearch && sr.Status == StatusSkipped {
			researchSkipped = true
		}
	}
	if m.Budget.Research == mission.ResearchNone {
		if !researchSkipped {
			t.Fatal("research.worker absent must skip")
		}
	} else if !researchRan {
		t.Fatal("research.worker selected must run the existing coordinator")
	}
	ctx := &TaskContext{Task: req, Mission: m, Capability: d, RunState: mission.NewRunState(m)}
	if err := bindFlagshipExploration(ctx); err != nil {
		t.Fatal(err)
	}
	if ctx.Exploration == nil || ctx.Exploration.Sealed == nil {
		t.Fatal("design.worker must seal one contract")
	}
	compileBuilderCapsule(ctx)
	if ctx.Capsule == nil {
		t.Fatal("coding.worker not consumed")
	}
	if d.Has(capability.IDJudgeIndependent) {
		compileJudgeCapsule(ctx, nil)
		if ctx.JudgeCapsule == nil || ctx.JudgeCapsule.Role != mission.RoleJudge {
			t.Fatal("judge.independent selected but capsule nil")
		}
	} else if ctx.JudgeCapsule != nil {
		t.Fatal("judge capsule present without selection")
	}
	ctxTO, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	ctx.Ctx = ctxTO
	ctx.Classification = &ClassificationData{RequiresVisual: true}
	v := NewPlaywrightVerifier(t.TempDir())
	data, verr := v.Verify(ctxTO, ctx)
	if verr != nil {
		t.Fatal(verr)
	}
	if data.LayerA != nil && data.LayerA.PlaywrightMissing {
		t.Fatalf("BLOCKED: %s", data.LayerA.EnvironmentError)
	}
	rawC, _ := json.Marshal(ctx.Capsule)
	st, _ := os.Stat(filepath.Join(mission.RunDir(workdir, req.ID), "capability.json"))
	t.Logf("flow-c capsule_bytes=%d capability_bytes=%d judge=%v", len(rawC), st.Size(), d.Has(capability.IDJudgeIndependent))
}

func TestPhase3E_FlowD_MultimodalFallbackStructural(t *testing.T) {
	visual.ResetScreenshotEvaluatorForTest()
	defer visual.ResetScreenshotEvaluatorForTest()
	m := mission.FromSeed(mission.Seed{
		ID: "p3e-d", Objective: "Build a flagship portfolio landing with unusual motion",
		DeclaredType: "DESIGN", QualityBar: mission.BarPremium,
	})
	pw, shot := true, false
	d, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw, ScreenshotEvalAvailable: &shot})
	if err != nil {
		t.Fatal(err)
	}
	if d.Has(capability.IDVisualMultimodal) {
		t.Fatal("multimodal must stay unselected")
	}
	if !d.Has(capability.IDVisualStructural) {
		t.Fatal("structural fallback missing")
	}
	ctx := &TaskContext{
		Task:           &TaskRequest{ID: m.ID, RawRequest: m.Objective, WorkspaceRoot: t.TempDir()},
		Mission:        m,
		Capability:     &d,
		Classification: &ClassificationData{RequiresVisual: true},
		RunState:       mission.NewRunState(m),
	}
	dir := t.TempDir()
	png := filepath.Join(dir, "desktop.png")
	if err := os.WriteFile(png, visual.TinyPNG, 0644); err != nil {
		t.Fatal(err)
	}
	rep := visual.Evaluate(visual.LayerBInput{
		Class:  visual.ClassFlagship,
		Review: string(mission.VisualReviewDeep),
		Captures: []visual.Capture{{
			Viewport:       visual.Viewport{Name: "desktop", Width: 1440, Height: 900},
			ScreenshotPath: png, Passed: true,
		}},
	})
	if rep.EvaluatorType == visual.EvaluatorScreenshot || rep.EvaluatorType == "multimodal" {
		t.Fatalf("evaluator_type %s", rep.EvaluatorType)
	}
	if visual.ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("multimodal fallback must not invoke screenshot eval")
	}
	_ = ctx
}

func TestPhase3E_FlowE_BrowserUnavailableHumanGate(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewPlaywrightVerifier(t.TempDir()))
	req := &TaskRequest{
		ID: "p3e-e", RawRequest: "Build a settings screen for profile CRUD",
		WorkspaceRoot: workdir, Type: "FEATURE", SkipVisualGate: true,
		Parameters: map[string]any{"playwright_available": false},
	}
	res, err := pipeline.Execute(context.Background(), req)
	if err != ErrHumanGateRequired {
		t.Fatalf("sentinel %v", err)
	}
	if res == nil || res.Status != PipelineStatusGated {
		t.Fatalf("status %+v", res)
	}
	st, err := mission.ReadRunState(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if st.PendingHuman == nil || !st.PendingHuman.Blocking {
		t.Fatal("pending_human blocking required")
	}
	if !strings.Contains(st.PendingHuman.Question, capability.IDBrowserPlaywright) {
		t.Fatalf("question %s", st.PendingHuman.Question)
	}
	if st.Attempt > 1 {
		t.Fatalf("attempt looped: %d", st.Attempt)
	}
	classifyLoops := 0
	for _, sr := range res.StageResults {
		if sr.StageName == StageClassify {
			classifyLoops++
		}
	}
	if classifyLoops != 1 {
		t.Fatalf("classify ran %d times", classifyLoops)
	}
}

func TestPhase3E_FlowF_HostSelectionOnly(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID: "p3e-f-both", RawRequest: "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir, Type: "BUGFIX",
		Parameters: map[string]any{
			"capability_hosts": []capability.HostSnapshot{
				{ID: "cursor", Available: true},
				{ID: "claude", Available: true},
			},
		},
	}
	if _, err := pipeline.Execute(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	d, err := capability.ReadDecision(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.HostLabel != "cursor" {
		t.Fatalf("host %s", d.HostLabel)
	}
	mraw, _ := os.ReadFile(filepath.Join(mission.RunDir(workdir, req.ID), "mission.json"))
	if leak := mission.ProviderLeak(mraw); leak != "" {
		t.Fatalf("mission leak %s", leak)
	}
	m, _ := mission.ReadMissionFile(workdir, req.ID)
	ctx := &TaskContext{Task: req, Mission: m, Capability: d, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	craw, _ := json.Marshal(ctx.Capsule)
	if leak := mission.ProviderLeak(craw); leak != "" {
		t.Fatalf("capsule leak %s", leak)
	}

	workdir2 := t.TempDir()
	pipeline2 := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req2 := &TaskRequest{
		ID: "p3e-f-fallback", RawRequest: "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir2, Type: "BUGFIX",
		Parameters: map[string]any{
			"capability_hosts": []capability.HostSnapshot{
				{ID: "cursor", Available: false},
				{ID: "claude", Available: true},
			},
		},
	}
	if _, err := pipeline2.Execute(context.Background(), req2); err != nil {
		t.Fatal(err)
	}
	d2, err := capability.ReadDecision(workdir2, req2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d2.HostLabel != "claude" {
		t.Fatalf("fallback host %s", d2.HostLabel)
	}

	cat3, graph3, workdir3 := setupEngineFixtures(t)
	defer os.RemoveAll(workdir3)
	pipe3 := NewDesignPipeline(cat3, graph3, nil, NewMockVisualVerifier())
	req3 := &TaskRequest{
		ID: "p3e-f-handoff", RawRequest: "Fix parser crash",
		WorkspaceRoot: workdir3, Type: "BUGFIX",
	}
	if _, err := pipe3.Execute(context.Background(), req3); err != nil {
		t.Fatal(err)
	}
	st, err := handoff.ReadState(workdir3)
	if err != nil {
		t.Fatal(err)
	}
	if st.SourceAgent != "antigravity" || st.TargetAgent != "cursor" {
		t.Fatalf("handoff literals changed: %s -> %s", st.SourceAgent, st.TargetAgent)
	}
}

func TestPhase3E_TamperBrowserRemovedFailsVisualQA(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "tamper", Objective: "Build a settings screen for profile CRUD", DeclaredType: "FEATURE",
	})
	d := capability.Decision{
		MissionID: m.ID,
		Selected: []capability.Item{
			{ID: capability.IDCodingWorker, Category: capability.CatCodingWorker, Strength: capability.StrengthRequired, Available: true},
		},
	}
	v := NewPlaywrightVerifier(t.TempDir())
	data, err := v.Verify(context.Background(), &TaskContext{
		Task:           &TaskRequest{ID: m.ID, RawRequest: m.Objective, WorkspaceRoot: t.TempDir(), Parameters: map[string]any{"app_url": "http://127.0.0.1:1"}},
		Mission:        m,
		Capability:     &d,
		Classification: &ClassificationData{RequiresVisual: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if data.AllPassed {
		t.Fatal("tampered decision must not AllPassed")
	}
	if data.LayerA != nil && len(data.ViewportResults) > 0 {
		t.Fatal("tamper must not run Layer A")
	}
}

func TestPhase3E_HostileCatalogStaysOutOfCapsule(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "hostile", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	var skills []string
	for i := 0; i < 44; i++ {
		skills = append(skills, "allowlisted-skill-"+itoa3e(i))
	}
	dump := bytes.Repeat([]byte("RESOURCE DUMP LINE\n"), 8000)
	var resources []string
	for i := 0; i < 168; i++ {
		resources = append(resources, "catalog-row-"+itoa3e(i))
	}
	pw := false
	d, err := capability.Select(capability.Input{
		Mission: m, PlaywrightAvailable: &pw, Skills: skills, Resources: resources,
		MCP: []capability.MCPCandidate{
			{ID: "playwright-mcp", Health: "HEALTHY"},
			{ID: "cursor-browser", Health: "HEALTHY"},
			{ID: "stitch", Health: "HEALTHY"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := &TaskContext{Mission: m, Capability: &d, Task: &TaskRequest{ID: m.ID}, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	if ctx.Capsule == nil {
		t.Fatal("capsule missing")
	}
	raw, _ := json.Marshal(ctx.Capsule)
	blob := string(raw)
	if len(ctx.Capsule.Skills) > capability.MaxSkills || len(ctx.Capsule.Resources) > capability.MaxResources || len(ctx.Capsule.ToolsAllowed) > capability.MaxMCPTools {
		t.Fatalf("capsule lists too large: skills=%d resources=%d tools=%d", len(ctx.Capsule.Skills), len(ctx.Capsule.Resources), len(ctx.Capsule.ToolsAllowed))
	}
	if strings.Count(blob, "allowlisted-skill-") == 44 {
		t.Fatal("full skill dump in capsule")
	}
	if strings.Contains(blob, "playwright-mcp") && strings.Contains(blob, "cursor-browser") && strings.Contains(blob, "stitch") {
		t.Fatal("full MCP name list in capsule")
	}
	t.Logf("hostile capsule_bytes=%d dump_bytes=%d ratio=%.4f", len(raw), len(dump), float64(len(raw))/float64(len(dump)))
	if len(raw) > 20000 {
		t.Fatalf("capsule not compact: %d", len(raw))
	}
	if float64(len(raw))/float64(len(dump)) >= 0.2 {
		t.Fatal("capsule not far under dump")
	}
}

func TestPhase3E_OpaqueModelDoesNotChangeSelection(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)
	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID: "p3e-opaque", RawRequest: "Fix nil pointer dereference in parser",
		WorkspaceRoot: workdir, Type: "BUGFIX",
		Parameters: map[string]any{"discovered_models": []string{"grok-4.7"}},
	}
	if _, err := pipeline.Execute(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	d, err := capability.ReadDecision(workdir, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	mraw, _ := os.ReadFile(filepath.Join(mission.RunDir(workdir, req.ID), "mission.json"))
	if strings.Contains(strings.ToLower(string(mraw)), "grok") {
		t.Fatal("mission json contained grok")
	}
	m, _ := mission.ReadMissionFile(workdir, req.ID)
	ctx := &TaskContext{Task: req, Mission: m, Capability: d, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	craw, _ := json.Marshal(ctx.Capsule)
	if strings.Contains(strings.ToLower(string(craw)), "grok") {
		t.Fatal("capsule contained grok")
	}
	pw := false
	alt, err := capability.Select(capability.Input{
		Mission: m, PlaywrightAvailable: &pw, DiscoveredModels: []string{"other-model-9"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if capIDs(d) != capIDs(&alt) {
		t.Fatalf("opaque label changed ids: %s vs %s", capIDs(d), capIDs(&alt))
	}
}

func TestPhase3E_CodingWorkerOmittedFailsImplement(t *testing.T) {
	ctx := &TaskContext{
		Capability: &capability.Decision{Selected: nil},
		Capsule:    nil,
	}
	if err := assertCodingWorkerConsumed(ctx); err == nil {
		t.Fatal("omitted coding.worker must fail implement")
	}
}

func capIDs(d *capability.Decision) string {
	if d == nil {
		return ""
	}
	var s []string
	for _, it := range d.Selected {
		if it.Category == capability.CatSkill || it.Category == capability.CatResource || it.Category == capability.CatMCPTool {
			continue
		}
		s = append(s, it.ID)
	}
	return strings.Join(s, ",")
}

func itoa3e(i int) string {
	if i == 0 {
		return "0"
	}
	n := i
	var d []byte
	for n > 0 {
		d = append([]byte{byte('0' + n%10)}, d...)
		n /= 10
	}
	return string(d)
}

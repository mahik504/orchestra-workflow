package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

// VisualVerifier defines the abstraction for multi-viewport UI verification.
// Canonical Layer A owner is the harness Playwright runner, not Playwright MCP
// or the Cursor browser. Those remain host adapters.
type VisualVerifier interface {
	Name() string
	Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error)
}

// PlaywrightVerifier is the canonical Layer A owner for orchestra runs.
type PlaywrightVerifier struct {
	OutputDir string
	Runner    visual.Runner
}

func NewPlaywrightVerifier(outputDir string) *PlaywrightVerifier {
	if outputDir == "" {
		outputDir = filepath.Join(".orchestra", "qa")
	}
	return &PlaywrightVerifier{OutputDir: outputDir, Runner: visual.NewPlaywrightRunner()}
}

func (v *PlaywrightVerifier) Name() string {
	return visual.CanonicalOwner
}

func (v *PlaywrightVerifier) Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error) {
	outDir := v.OutputDir
	if outDir == "" {
		if taskCtx != nil && taskCtx.Task != nil && taskCtx.Task.WorkspaceRoot != "" {
			outDir = filepath.Join(taskCtx.Task.WorkspaceRoot, ".orchestra", "qa")
		} else {
			outDir = filepath.Join(".orchestra", "qa")
		}
	}
	if err := resources.CheckQuarantineBoundary(outDir); err != nil {
		return nil, err
	}
	_ = os.MkdirAll(outDir, 0755)

	class := designClassOf(taskCtx)
	policy := visual.PolicyFor(class)
	review := policy.VisualReview
	judgePolicy := policy.IndependentJudge
	if taskCtx != nil && taskCtx.Mission != nil {
		if taskCtx.Mission.Budget.VisualReview != "" {
			review = string(taskCtx.Mission.Budget.VisualReview)
		}
		if taskCtx.Mission.Budget.IndependentJudge != "" {
			judgePolicy = string(taskCtx.Mission.Budget.IndependentJudge)
		}
	}

	data := &VisualQAData{
		AllPassed:          false,
		ViewportResults:    []ViewportCheckResult{},
		Metrics:            map[string]float64{},
		DetectedViolations: []string{},
		FailureClass:       FailureClassNone,
		DesignClass:        string(class),
		CanonicalOwner:     visual.CanonicalOwner,
	}

	if review == string(mission.VisualReviewNone) && class == visual.ClassOrdinary {
		return skippedVisual(class), nil
	}

	if taskCtx != nil && taskCtx.Capability != nil {
		wantsBrowser := review == string(mission.VisualReviewBrowser) || review == string(mission.VisualReviewDeep)
		if wantsBrowser && !taskCtx.Capability.Has(capability.IDBrowserPlaywright) {
			data.AllPassed = false
			data.Verdict = CheckFail
			data.FinalResult = "fail"
			data.FailureClass = FailureClassFatal
			data.DetectedViolations = append(data.DetectedViolations, "required browser.playwright was not selected; visual QA must not pass")
			data.Findings = append(data.Findings, visual.Finding{
				Type: "evidence", Layer: "A",
				Explanation: "required browser.playwright was not selected; visual QA must not pass",
				Action:      "Restore browser.playwright on the capability decision before Layer A.",
			})
			return data, nil
		}
		if !taskCtx.Capability.Has(capability.IDBrowserPlaywright) {
			return skippedVisual(class), nil
		}
	}

	req := visual.LayerARequest{
		BaseURL:   appURL(taskCtx),
		OutputDir: outDir,
		Routes:    appRoutes(taskCtx),
		Viewports: visual.DefaultViewports(visual.NeedsTablet(objectiveOf(taskCtx), paramsOf(taskCtx))),
		Surface:   surfaceOf(taskCtx),
	}
	if req.BaseURL == "" {
		req.BaseURL = strings.TrimSpace(os.Getenv("ORCHESTRA_APP_URL"))
	}

	runner := v.Runner
	if runner == nil {
		runner = visual.NewPlaywrightRunner()
	}
	aRep, err := runner.Run(ctx, req)
	if err != nil {
		return nil, err
	}
	if aRep == nil {
		aRep = &visual.LayerAReport{Owner: visual.CanonicalOwner, Passed: false, EnvironmentError: "empty Layer A report"}
	}
	data.LayerA = aRep
	if data.Metrics == nil {
		data.Metrics = map[string]float64{}
	}
	data.Metrics["layer_a_startup_ms"] = float64(aRep.StartupMS)
	data.Metrics["layer_a_screenshot_ms"] = float64(aRep.ScreenshotMS)
	data.Findings = append(data.Findings, aRep.Findings...)
	for _, c := range aRep.Captures {
		if c.ScreenshotPath != "" {
			if b, err := os.ReadFile(c.ScreenshotPath); err == nil && visual.IsMockScreenshot(b) {
				aRep.Passed = false
				data.Findings = append(data.Findings, visual.Finding{
					Type: "evidence", Layer: "A", Viewport: c.Viewport.Name, Route: c.Route,
					Explanation:  "MOCK_SCREENSHOT placeholders are not rendered evidence.",
					Action:       "Capture a real PNG from the running application.",
					EvidencePath: c.ScreenshotPath,
				})
			}
		}
		data.ViewportResults = append(data.ViewportResults, ViewportCheckResult{
			ViewportName:   c.Viewport.Name,
			Width:          c.Viewport.Width,
			Height:         c.Viewport.Height,
			ScreenshotPath: c.ScreenshotPath,
			HasOverflow:    c.Overflow,
			MaxScrollWidth: c.ScrollWidth,
			Passed:         c.Passed,
		})
		data.Evidence = append(data.Evidence, visual.Evidence{
			What: "layer-a-capture", Route: c.Route, Surface: req.Surface,
			Viewport: c.Viewport.Name, Layer: "A", Passed: c.Passed,
			ScreenshotPath: c.ScreenshotPath, Findings: nil,
		})
	}

	needB := (review == string(mission.VisualReviewBrowser) || review == string(mission.VisualReviewDeep)) && class != visual.ClassOrdinary
	runStructural := needB
	runScreenshot := needB
	if taskCtx != nil && taskCtx.Capability != nil {
		runStructural = taskCtx.Capability.Has(capability.IDVisualStructural)
		runScreenshot = taskCtx.Capability.Has(capability.IDVisualScreenshot)
		needB = runStructural || runScreenshot
	}
	if needB {
		inB := visual.LayerBInput{
			Class:    class,
			Review:   review,
			Captures: aRep.Captures,
			Surface:  req.Surface,
			Critic:   criticFromTask(taskCtx, class, req.Surface, aRep),
		}
		var bRep *visual.LayerBReport
		if runScreenshot {
			bRep = visual.ApplyScreenshotEval(ctx, inB, objectiveOf(taskCtx), contractPathOf(taskCtx), nil)
		} else {
			bRep = visual.Evaluate(inB)
		}
		data.LayerB = bRep
		if bRep != nil {
			data.Findings = append(data.Findings, bRep.Findings...)
			if data.Metrics == nil {
				data.Metrics = map[string]float64{}
			}
			data.Metrics["layer_b_capsule_bytes"] = float64(bRep.CapsuleBytes)
			if bRep.Invoked {
				data.Metrics["screenshot_evaluator_invocations"] = 1
			}
			for _, c := range aRep.Captures {
				data.Evidence = append(data.Evidence, visual.Evidence{
					What: "layer-b-eval", Route: c.Route, Surface: req.Surface,
					Viewport: c.Viewport.Name, Layer: "B", Passed: bRep.Passed,
					ScreenshotPath: c.ScreenshotPath, EvaluatorType: bRep.EvaluatorType,
					Confidence: bRep.Confidence,
				})
			}
		}
	}

	judgeIn := visual.JudgeInput{
		Class:               class,
		Policy:              judgePolicy,
		VisualNovelty:       noveltyOf(taskCtx),
		AcceptanceAmbiguous: false,
		PriorVisualFail:     priorVisualFail(taskCtx),
		Confidence:          confidenceOf(taskCtx),
		MultiSurface:        multiSurface(taskCtx),
		SecurityRelease:     taskCtx != nil && taskCtx.Mission != nil && taskCtx.Mission.RequiresSecurity,
		HumanRequested:      judgeForced(taskCtx),
	}
	jd := visual.DecideJudge(judgeIn)
	data.JudgeMode = jd.Mode
	data.JudgeActivated = jd.Activated
	if taskCtx != nil && taskCtx.Capability != nil {
		data.JudgeActivated = taskCtx.Capability.Has(capability.IDJudgeIndependent)
		if data.JudgeActivated {
			compileJudgeCapsule(taskCtx, data)
			if taskCtx.JudgeCapsule == nil {
				data.Findings = append(data.Findings, visual.Finding{
					Type: "evidence", Layer: "B",
					Explanation: "judge.independent selected but judge capsule is nil",
					Action:      "compile the scoped RoleJudge capsule when the decision selects judge.independent",
				})
			}
		}
	} else if jd.Activated && taskCtx != nil {
		compileJudgeCapsule(taskCtx, data)
	}

	for _, f := range data.Findings {
		data.DetectedViolations = append(data.DetectedViolations, f.Explanation)
	}
	aOK := aRep.Passed && aRep.EnvironmentError == ""
	bOK := data.LayerB == nil || data.LayerB.Passed
	judgeOK := true
	if taskCtx != nil && taskCtx.Capability != nil && taskCtx.Capability.Has(capability.IDJudgeIndependent) && taskCtx.JudgeCapsule == nil {
		judgeOK = false
	}
	data.AllPassed = aOK && bOK && judgeOK && len(data.Findings) == 0
	if !data.AllPassed {
		data.FailureClass = visual.MapFailureClass(data.Findings)
		if data.FailureClass == "NONE" {
			data.FailureClass = FailureClassLayoutCode
		}
		data.FinalResult = "fail"
		if aRep != nil && aRep.EnvironmentError != "" {
			data.Verdict = CheckBlocked
		} else {
			data.Verdict = CheckFail
		}
	} else {
		data.FailureClass = FailureClassNone
		data.FinalResult = "pass"
		data.Verdict = CheckPass
	}
	if taskCtx != nil && taskCtx.Iteration != nil {
		data.RetryCount = taskCtx.Iteration.CurrentIteration
		data.ChangesMade = append([]string{}, taskCtx.Iteration.CorrectiveFeedback...)
	}
	for i := range data.Evidence {
		data.Evidence[i].RetryCount = data.RetryCount
		data.Evidence[i].Changes = data.ChangesMade
		data.Evidence[i].FinalResult = data.FinalResult
		data.Evidence[i].Findings = data.Findings
	}
	return data, nil
}

func skippedVisual(class visual.Class) *VisualQAData {
	return &VisualQAData{
		AllPassed:          false,
		Verdict:            CheckNotApplicable,
		ViewportResults:    []ViewportCheckResult{},
		Metrics:            map[string]float64{},
		DetectedViolations: []string{},
		FailureClass:       FailureClassNone,
		DesignClass:        string(class),
		FinalResult:        "skipped_ordinary",
	}
}

// visualClosed is true when iteration must not treat the check as a failure.
// A skipped check is closed and is not a pass.
func visualClosed(qa *VisualQAData) bool {
	if qa == nil {
		return true
	}
	switch qa.Verdict {
	case CheckPass, CheckNotApplicable:
		return true
	case CheckFail, CheckBlocked, CheckNotRun:
		return false
	default:
		return qa.AllPassed
	}
}

func designClassOf(taskCtx *TaskContext) visual.Class {
	in := visual.ClassifyInput{}
	if taskCtx != nil && taskCtx.Mission != nil {
		in.Type = taskCtx.Mission.Type
		in.Objective = taskCtx.Mission.Objective
		in.VisualImportance = taskCtx.Mission.VisualImportance
		in.RequiresVisual = taskCtx.Mission.VisualImportance != mission.VisualNone
	} else if taskCtx != nil && taskCtx.Task != nil {
		in.Objective = taskCtx.Task.RawRequest
		in.Type = taskCtx.Task.Type
		in.Tags = taskCtx.Task.Tags
	}
	if taskCtx != nil && taskCtx.Classification != nil {
		in.RequiresVisual = in.RequiresVisual || taskCtx.Classification.RequiresVisual
	}
	return visual.Classify(in)
}

func appURL(taskCtx *TaskContext) string {
	if taskCtx == nil || taskCtx.Task == nil || taskCtx.Task.Parameters == nil {
		return ""
	}
	if u, ok := taskCtx.Task.Parameters["app_url"].(string); ok {
		return strings.TrimSpace(u)
	}
	return ""
}

func appRoutes(taskCtx *TaskContext) []string {
	if taskCtx != nil && taskCtx.Task != nil && taskCtx.Task.Parameters != nil {
		switch v := taskCtx.Task.Parameters["routes"].(type) {
		case []string:
			if len(v) > 0 {
				return v
			}
		case []any:
			var out []string
			for _, x := range v {
				if s, ok := x.(string); ok && s != "" {
					out = append(out, s)
				}
			}
			if len(out) > 0 {
				return out
			}
		}
	}
	return []string{"/"}
}

func paramsOf(taskCtx *TaskContext) map[string]any {
	if taskCtx != nil && taskCtx.Task != nil {
		return taskCtx.Task.Parameters
	}
	return nil
}

func objectiveOf(taskCtx *TaskContext) string {
	if taskCtx != nil && taskCtx.Mission != nil {
		return taskCtx.Mission.Objective
	}
	if taskCtx != nil && taskCtx.Task != nil {
		return taskCtx.Task.RawRequest
	}
	return ""
}

func surfaceOf(taskCtx *TaskContext) string {
	if taskCtx != nil && taskCtx.Mission != nil && len(taskCtx.Mission.AffectedSurfaces) > 0 {
		return taskCtx.Mission.AffectedSurfaces[0]
	}
	return ""
}

func noveltyOf(taskCtx *TaskContext) string {
	if taskCtx != nil && taskCtx.Mission != nil {
		return taskCtx.Mission.Novelty
	}
	return ""
}

func confidenceOf(taskCtx *TaskContext) string {
	if taskCtx != nil && taskCtx.Mission != nil {
		return taskCtx.Mission.Confidence
	}
	return ""
}

func multiSurface(taskCtx *TaskContext) bool {
	return taskCtx != nil && taskCtx.Mission != nil && len(taskCtx.Mission.AffectedSurfaces) > 1
}

func priorVisualFail(taskCtx *TaskContext) bool {
	if taskCtx == nil || taskCtx.RunState == nil {
		return false
	}
	for _, h := range taskCtx.RunState.FailureHistory {
		low := strings.ToLower(h.PreviousFailure + " " + string(h.FailureClass))
		if strings.Contains(low, "visual") || strings.Contains(low, "overflow") || strings.Contains(low, "layer") {
			return true
		}
	}
	return false
}

func judgeForced(taskCtx *TaskContext) bool {
	if taskCtx == nil || taskCtx.Task == nil || taskCtx.Task.Parameters == nil {
		return false
	}
	v, ok := taskCtx.Task.Parameters["independent_judge"]
	if !ok {
		return false
	}
	s := strings.ToLower(fmt.Sprint(v))
	return s == "required" || s == "true" || s == "on"
}

func contractPathOf(taskCtx *TaskContext) string {
	if taskCtx != nil && taskCtx.DesignSystem != nil && taskCtx.DesignSystem.DesignMDPath != "" {
		return taskCtx.DesignSystem.DesignMDPath
	}
	return ""
}

func criticFromTask(taskCtx *TaskContext, class visual.Class, surface string, aRep *visual.LayerAReport) visual.CriticInput {
	in := visual.CriticInput{Class: class, Surface: surface}
	if aRep != nil && len(aRep.Captures) > 0 {
		in.ScreenshotPath = aRep.Captures[0].ScreenshotPath
		in.Route = aRep.Captures[0].Route
		in.Viewport = aRep.Captures[0].Viewport.Name
	}
	if taskCtx != nil && taskCtx.Task != nil && taskCtx.Task.Parameters != nil {
		if v, ok := taskCtx.Task.Parameters["critic"].(map[string]any); ok {
			in.CenteredHero = boolOf(v["centered_hero"])
			in.NoProductAnchor = boolOf(v["no_product_anchor"])
			in.GenericType = boolOf(v["generic_type"])
			in.NoHierarchy = boolOf(v["no_hierarchy"])
			in.NoSpatialPlan = boolOf(v["no_spatial_plan"])
			in.NoInteraction = boolOf(v["no_interaction"])
			in.ArbitraryMotion = boolOf(v["arbitrary_motion"])
			in.HasGradient = boolOf(v["has_gradient"])
			in.HasRounded = boolOf(v["has_rounded"])
			in.HasDark = boolOf(v["has_dark"])
			in.HasCommonFont = boolOf(v["has_common_font"])
			in.HasFamiliarUI = boolOf(v["has_familiar_ui"])
		}
	}
	return in
}

func boolOf(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1" || t == "yes"
	default:
		return false
	}
}

func persistVisualEvidence(ctx *TaskContext) {
	if ctx == nil || ctx.RunState == nil || ctx.VisualQA == nil {
		return
	}
	for _, ev := range ctx.VisualQA.Evidence {
		if ev.ScreenshotPath == "" {
			continue
		}
		ctx.RunState.Outputs = append(ctx.RunState.Outputs, mission.EvidenceRef{
			Name: ev.Layer + "-" + ev.Viewport,
			Path: ev.ScreenshotPath,
		})
	}
	if ctx.Task == nil || ctx.Task.WorkspaceRoot == "" || ctx.Mission == nil {
		return
	}
	bundle := visual.Bundle{
		Owner:       ctx.VisualQA.CanonicalOwner,
		Class:       visual.Class(ctx.VisualQA.DesignClass),
		LayerA:      ctx.VisualQA.LayerA,
		LayerB:      ctx.VisualQA.LayerB,
		Findings:    ctx.VisualQA.Findings,
		Items:       ctx.VisualQA.Evidence,
		RetryCount:  ctx.VisualQA.RetryCount,
		Changes:     ctx.VisualQA.ChangesMade,
		FinalResult: ctx.VisualQA.FinalResult,
		JudgeMode:   ctx.VisualQA.JudgeMode,
	}
	raw, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return
	}
	dir := mission.RunDir(ctx.Task.WorkspaceRoot, ctx.Mission.ID)
	_ = os.MkdirAll(dir, 0755)
	_ = os.WriteFile(filepath.Join(dir, "visual-evidence.json"), raw, 0644)
}

// MockVisualVerifier provides deterministic testing for self-healing loops.
type MockVisualVerifier struct {
	SimulateMobileOverflow   bool
	SimulateTokenStyleDefect bool
	FailUntilIteration       int
	AlwaysFail               bool
}

func NewMockVisualVerifier() *MockVisualVerifier {
	return &MockVisualVerifier{}
}

func (m *MockVisualVerifier) Name() string {
	return visual.OwnerMockTest
}

func (m *MockVisualVerifier) Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error) {
	currentIter := 0
	if taskCtx != nil && taskCtx.Iteration != nil {
		currentIter = taskCtx.Iteration.CurrentIteration
	}

	viewports := []struct {
		name   string
		width  int
		height int
	}{
		{"desktop", 1440, 900},
		{"mobile", 390, 844},
	}

	var vpResults []ViewportCheckResult
	var violations []string
	var findings []visual.Finding
	failureClass := FailureClassNone

	shouldFail := false
	if m.AlwaysFail {
		shouldFail = true
	} else if m.FailUntilIteration > 0 {
		shouldFail = currentIter < m.FailUntilIteration
	} else if m.SimulateMobileOverflow || m.SimulateTokenStyleDefect {
		shouldFail = true
	}

	for _, vp := range viewports {
		hasOverflow := false
		maxScroll := vp.width
		passed := true
		if shouldFail && vp.name == "mobile" && m.SimulateMobileOverflow {
			hasOverflow = true
			maxScroll = 415
			passed = false
			violations = append(violations, "Mobile horizontal overflow detected: scrollWidth 415px exceeds viewport 390px")
			failureClass = FailureClassLayoutCode
			findings = append(findings, visual.Finding{
				Type: "overflow", Layer: "A", Viewport: "mobile", Trivial: true,
				Explanation: "Mobile horizontal overflow detected: scrollWidth 415px exceeds viewport 390px",
				Action:      "Constrain the overflowing element on mobile and re-run Layer A.",
			})
		}
		vpResults = append(vpResults, ViewportCheckResult{
			ViewportName:   vp.name,
			Width:          vp.width,
			Height:         vp.height,
			ScreenshotPath: filepath.Join(".orchestra", "qa", fmt.Sprintf("%s.png", vp.name)),
			HasOverflow:    hasOverflow,
			MaxScrollWidth: maxScroll,
			Passed:         passed,
		})
	}

	if shouldFail && m.SimulateTokenStyleDefect {
		violations = append(violations, "Typography pairing is not a deliberate product-specific decision.")
		failureClass = FailureClassTokenStyle
		findings = append(findings, visual.Finding{
			Type: "typography", Layer: "B",
			Explanation: "Typography pairing is not a deliberate product-specific decision.",
			Action:      "Specify a display + body pairing with a reason and apply it to the rendered hero.",
		})
	}

	allPassed := len(violations) == 0
	final := "pass"
	verdict := CheckPass
	if !allPassed {
		final = "fail"
		verdict = CheckFail
	}
	return &VisualQAData{
		AllPassed:          allPassed,
		Verdict:            verdict,
		ViewportResults:    vpResults,
		Metrics:            map[string]float64{"cls": 0.02},
		DetectedViolations: violations,
		FailureClass:       failureClass,
		Findings:           findings,
		CanonicalOwner:     visual.OwnerMockTest,
		FinalResult:        final,
		RetryCount:         currentIter,
	}, nil
}

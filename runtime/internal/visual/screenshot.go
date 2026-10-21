package visual

import (
	"context"
	"encoding/json"
	"sync"
)

const (
	EvaluatorNone       = "none"
	EvaluatorStructural = "structural"
	EvaluatorScreenshot = "screenshot"
)

// ScreenshotEvaluator is an optional, provider-neutral adapter. Core never
// names a model. Invoked only when visual_review is browser|deep, class is
// product_ui or flagship, and an adapter is registered.
type ScreenshotEvaluator interface {
	Evaluate(ctx context.Context, req ScreenshotEvalRequest) (*ScreenshotEvalResult, error)
}

type ScreenshotEvalRequest struct {
	Objective      string   `json:"objective,omitempty"`
	ContractPath   string   `json:"contract_path,omitempty"`
	Criteria       []string `json:"criteria,omitempty"`
	ScreenshotRefs []string `json:"screenshot_refs,omitempty"`
	Surface        string   `json:"surface,omitempty"`
	Route          string   `json:"route,omitempty"`
	Viewport       string   `json:"viewport,omitempty"`
}

type ScreenshotEvalResult struct {
	EvaluatorType string    `json:"evaluator_type"`
	Findings      []Finding `json:"findings,omitempty"`
	Confidence    string    `json:"confidence,omitempty"`
	Result        string    `json:"result,omitempty"`
}

type EvalCapsuleReport struct {
	CapsuleBytes      int      `json:"capsule_bytes"`
	WholeProjectBytes int      `json:"whole_project_bytes,omitempty"`
	Ratio             float64  `json:"ratio,omitempty"`
	Excluded          []string `json:"excluded,omitempty"`
}

var (
	screenshotMu          sync.Mutex
	screenshotEval        ScreenshotEvaluator
	screenshotInvocations int
)

func RegisterScreenshotEvaluator(e ScreenshotEvaluator) {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	screenshotEval = e
}

func ResetScreenshotEvaluatorForTest() {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	screenshotEval = nil
	screenshotInvocations = 0
}

func ActiveScreenshotEvaluator() ScreenshotEvaluator {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	return screenshotEval
}

func ScreenshotEvaluatorInvocations() int {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	return screenshotInvocations
}

func recordScreenshotInvocation() {
	screenshotMu.Lock()
	screenshotInvocations++
	screenshotMu.Unlock()
}

func WantsScreenshotEval(class Class, review string) bool {
	if class == ClassOrdinary {
		return false
	}
	return review == "browser" || review == "deep"
}

func CompileScreenshotEvalCapsule(req ScreenshotEvalRequest, wholeProject []byte) ([]byte, EvalCapsuleReport) {
	payload := ScreenshotEvalRequest{
		Objective:      req.Objective,
		ContractPath:   req.ContractPath,
		Criteria:       append([]string{}, req.Criteria...),
		ScreenshotRefs: append([]string{}, req.ScreenshotRefs...),
		Surface:        req.Surface,
		Route:          req.Route,
		Viewport:       req.Viewport,
	}
	raw, _ := json.Marshal(payload)
	rep := EvalCapsuleReport{
		CapsuleBytes: len(raw),
		Excluded:     []string{"conversation_history", "whole_project_context", "mcp_catalog", "unrelated_skills"},
	}
	if len(wholeProject) > 0 {
		rep.WholeProjectBytes = len(wholeProject)
		rep.Ratio = float64(len(raw)) / float64(len(wholeProject))
	}
	return raw, rep
}

func ApplyScreenshotEval(ctx context.Context, in LayerBInput, objective, contractPath string, criteria []string) *LayerBReport {
	base := Evaluate(in)
	if base == nil {
		base = &LayerBReport{Passed: false, EvaluatorType: EvaluatorStructural}
	}
	if in.Class == ClassOrdinary || in.Review == "" || in.Review == "none" {
		base.EvaluatorType = EvaluatorNone
		return base
	}
	if base.EvaluatorType == "" {
		base.EvaluatorType = EvaluatorStructural
	}
	if !WantsScreenshotEval(in.Class, in.Review) {
		return base
	}
	eval := ActiveScreenshotEvaluator()
	if eval == nil {
		base.EvaluatorNote = "no screenshot adapter registered; Layer B stays structural"
		return base
	}
	refs := make([]string, 0, len(in.Captures))
	route, viewport := "", ""
	for _, c := range in.Captures {
		if c.ScreenshotPath != "" {
			refs = append(refs, c.ScreenshotPath)
		}
		if route == "" {
			route = c.Route
			viewport = c.Viewport.Name
		}
	}
	req := ScreenshotEvalRequest{
		Objective:      objective,
		ContractPath:   contractPath,
		Criteria:       criteria,
		ScreenshotRefs: refs,
		Surface:        in.Surface,
		Route:          route,
		Viewport:       viewport,
	}
	capsule, scope := CompileScreenshotEvalCapsule(req, nil)
	base.CapsuleBytes = scope.CapsuleBytes
	_ = capsule
	recordScreenshotInvocation()
	res, err := eval.Evaluate(ctx, req)
	base.EvaluatorType = EvaluatorScreenshot
	base.Invoked = true
	if err != nil {
		base.Passed = false
		base.Reason = "screenshot evaluator failed: " + err.Error()
		base.Findings = append(base.Findings, Finding{
			Type: "evidence", Layer: "B", Surface: in.Surface,
			Explanation: base.Reason,
			Action:      "Repair the screenshot evaluator adapter or fall back to structural Layer B.",
		})
		return base
	}
	if res == nil {
		return base
	}
	if res.Confidence != "" {
		base.Confidence = res.Confidence
	}
	if len(res.Findings) > 0 {
		base.Findings = append(base.Findings, res.Findings...)
	}
	if res.Result == "fail" || len(res.Findings) > 0 {
		base.Passed = false
		if base.Reason == "" || base.Reason == "Layer B accepted rendered evidence; no actionable anti-slop findings" {
			base.Reason = "screenshot evaluator reported visual issues"
		}
	}
	if res.Result == "pass" && len(base.Findings) == 0 {
		base.Passed = true
		base.Reason = "screenshot evaluator accepted rendered evidence"
	}
	base.Passed = base.Passed && len(base.Findings) == 0
	return base
}

package visual

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type stubShot struct {
	fail bool
}

func (s stubShot) Evaluate(ctx context.Context, req ScreenshotEvalRequest) (*ScreenshotEvalResult, error) {
	res := &ScreenshotEvalResult{
		EvaluatorType: EvaluatorScreenshot,
		Confidence:    "medium",
		Result:        "pass",
	}
	if s.fail {
		res.Result = "fail"
		res.Findings = []Finding{{
			Type: "composition", Layer: "B", Surface: req.Surface, Route: req.Route, Viewport: req.Viewport,
			Explanation:  "hero lacks a product-specific visual object",
			Action:       "Give the hero one named product object, then re-evaluate.",
			EvidencePath: firstRef(req.ScreenshotRefs),
		}}
	}
	return res, nil
}

func firstRef(refs []string) string {
	if len(refs) == 0 {
		return ""
	}
	return refs[0]
}

func TestScreenshotEvaluatorProviderNeutral(t *testing.T) {
	b, err := os.ReadFile("screenshot.go")
	if err != nil {
		t.Fatal(err)
	}
	low := bytes.ToLower(b)
	for _, needle := range []string{"gpt-4", "gpt-5", "claude-3", "claude-sonnet", "gemini-1", "gemini-2"} {
		if bytes.Contains(low, []byte(needle)) {
			t.Fatalf("core screenshot evaluator must not hardcode model %s", needle)
		}
	}
}

func TestScreenshotEvalCapsuleExcludesDump(t *testing.T) {
	caps, rep := CompileScreenshotEvalCapsule(ScreenshotEvalRequest{
		Objective:      "Review flagship hero",
		ContractPath:   "design.md",
		Criteria:       []string{"one focal point"},
		ScreenshotRefs: []string{"qa/desktop.png", "qa/mobile.png"},
		Surface:        "hero",
		Route:          "/",
		Viewport:       "desktop",
	}, bytes.Repeat([]byte("WHOLE PROJECT DUMP\nplaywright-mcp\nhallmark\nfull conversation\n"), 2000))
	blob := string(caps)
	if strings.Contains(blob, "WHOLE PROJECT DUMP") || strings.Contains(blob, "playwright-mcp") ||
		strings.Contains(blob, "hallmark") || strings.Contains(blob, "full conversation") {
		t.Fatalf("eval capsule leaked excluded context: %s", blob)
	}
	if rep.WholeProjectBytes == 0 || rep.Ratio >= 1 {
		t.Fatalf("capsule must be smaller than project dump: %+v", rep)
	}
	t.Logf("screenshot_eval_capsule_bytes=%d whole_project_bytes=%d ratio=%.4f excluded=%v",
		rep.CapsuleBytes, rep.WholeProjectBytes, rep.Ratio, rep.Excluded)
}

func TestScreenshotEvalOnlyWhenPolicyAndAdapter(t *testing.T) {
	ResetScreenshotEvaluatorForTest()
	defer ResetScreenshotEvaluatorForTest()
	dir := t.TempDir()
	p := filepath.Join(dir, "desktop.png")
	if err := os.WriteFile(p, TinyPNG, 0644); err != nil {
		t.Fatal(err)
	}
	in := LayerBInput{
		Class:  ClassFlagship,
		Review: "deep",
		Captures: []Capture{{
			Viewport: Viewport{Name: "desktop", Width: 1440, Height: 900},
			Route:    "/", ScreenshotPath: p, Passed: true,
		}},
		Surface: "hero",
	}
	rep := ApplyScreenshotEval(context.Background(), in, "flagship hero", "design.md", []string{"focal point"})
	if rep.EvaluatorType != EvaluatorStructural {
		t.Fatalf("without adapter expected structural, got %s", rep.EvaluatorType)
	}
	if !strings.Contains(rep.EvaluatorNote, "no screenshot adapter") {
		t.Fatalf("expected adapter-missing note: %s", rep.EvaluatorNote)
	}
	if ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("must not invoke a missing adapter")
	}

	RegisterScreenshotEvaluator(stubShot{fail: true})
	rep = ApplyScreenshotEval(context.Background(), in, "flagship hero", "design.md", []string{"focal point"})
	if ScreenshotEvaluatorInvocations() != 1 {
		t.Fatalf("invocations=%d", ScreenshotEvaluatorInvocations())
	}
	if rep.EvaluatorType != EvaluatorScreenshot {
		t.Fatalf("type %s", rep.EvaluatorType)
	}
	if rep.Passed {
		t.Fatal("screenshot evaluator findings must fail Layer B")
	}
	if rep.Confidence == "" || rep.CapsuleBytes == 0 {
		t.Fatalf("confidence and capsule size required: %+v", rep)
	}
}

func TestOrdinaryNeverInvokesScreenshotEval(t *testing.T) {
	ResetScreenshotEvaluatorForTest()
	defer ResetScreenshotEvaluatorForTest()
	RegisterScreenshotEvaluator(stubShot{})
	rep := ApplyScreenshotEval(context.Background(), LayerBInput{
		Class: ClassOrdinary, Review: "none",
	}, "fix nil pointer", "", nil)
	if ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("ordinary API/debug must not invoke image eval")
	}
	if rep.EvaluatorType != EvaluatorNone {
		t.Fatalf("evaluator_type=%s", rep.EvaluatorType)
	}
}

func TestEvidenceStoresPathsNotRawImages(t *testing.T) {
	ev := Evidence{
		What: "layer-a-capture", Layer: "A", Viewport: "desktop",
		ScreenshotPath: "qa/desktop.png", EvaluatorType: EvaluatorStructural, Passed: true,
	}
	js, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(js, TinyPNG) || bytes.Contains(js, []byte("MOCK_SCREENSHOT_")) {
		t.Fatal("evidence must not embed image bytes")
	}
	if !bytes.Contains(js, []byte("qa/desktop.png")) {
		t.Fatal("evidence must store screenshot path")
	}
}

package visual

import (
	"fmt"
	"os"
)

type LayerBReport struct {
	Passed        bool      `json:"passed"`
	Findings      []Finding `json:"findings,omitempty"`
	Reason        string    `json:"reason,omitempty"`
	EvaluatorType string    `json:"evaluator_type,omitempty"`
	EvaluatorNote string    `json:"evaluator_note,omitempty"`
	Confidence    string    `json:"confidence,omitempty"`
	Invoked       bool      `json:"invoked,omitempty"`
	CapsuleBytes  int       `json:"capsule_bytes,omitempty"`
}

type LayerBInput struct {
	Class      Class
	Review     string // none | browser | deep
	Captures   []Capture
	Surface    string
	SourceOnly bool
	DOMOnly    bool
	Assertion  string
	Critic     CriticInput
}

// Evaluate is screenshot-capable perception on Layer A captures.
// It cannot pass on source, DOM, static HTML, or a developer assertion.
func Evaluate(in LayerBInput) *LayerBReport {
	rep := &LayerBReport{Passed: false, EvaluatorType: EvaluatorStructural}
	if in.Review == "" || in.Review == "none" {
		rep.Passed = true
		rep.EvaluatorType = EvaluatorNone
		rep.Reason = "visual_review is none; Layer B not required"
		return rep
	}
	if in.Class == ClassOrdinary {
		rep.Passed = true
		rep.EvaluatorType = EvaluatorNone
		rep.Reason = "ordinary work is not blocked by subjective visual criteria"
		return rep
	}
	if in.SourceOnly || in.DOMOnly || in.Assertion != "" && len(in.Captures) == 0 {
		rep.Reason = "Layer B cannot pass on source code, DOM, static HTML, or assertion without rendered screenshots"
		rep.Findings = []Finding{{
			Type:        "evidence",
			Layer:       "B",
			Surface:     in.Surface,
			Explanation: rep.Reason,
			Action:      "Run Layer A against the live app and pass those PNG captures to Layer B.",
		}}
		return rep
	}
	if len(in.Captures) == 0 {
		rep.Reason = "Layer B requires Layer A screenshot evidence"
		rep.Findings = []Finding{{
			Type:        "evidence",
			Layer:       "B",
			Surface:     in.Surface,
			Explanation: rep.Reason,
			Action:      "Capture desktop and mobile screenshots from the running application, then re-evaluate.",
		}}
		return rep
	}
	for _, c := range in.Captures {
		if c.ScreenshotPath == "" {
			rep.Reason = "a Layer A capture is missing its screenshot path"
			rep.Findings = append(rep.Findings, Finding{
				Type: "evidence", Layer: "B", Surface: in.Surface, Route: c.Route, Viewport: c.Viewport.Name,
				Explanation: rep.Reason,
				Action:      "Re-run Layer A so every viewport has a real PNG.",
			})
			return rep
		}
		b, err := os.ReadFile(c.ScreenshotPath)
		if err != nil || !AcceptableScreenshot(b) {
			rep.Reason = fmt.Sprintf("screenshot %s is not rendered evidence", c.ScreenshotPath)
			rep.Findings = append(rep.Findings, Finding{
				Type: "evidence", Layer: "B", Surface: in.Surface, Route: c.Route, Viewport: c.Viewport.Name,
				Explanation:  rep.Reason,
				Action:       "Replace mock or missing files with Layer A PNG captures from the running app.",
				EvidencePath: c.ScreenshotPath,
			})
			return rep
		}
	}
	criticIn := in.Critic
	if criticIn.Class == "" {
		criticIn.Class = in.Class
	}
	if criticIn.Surface == "" {
		criticIn.Surface = in.Surface
	}
	if criticIn.ScreenshotPath == "" && len(in.Captures) > 0 {
		criticIn.ScreenshotPath = in.Captures[0].ScreenshotPath
	}
	rep.Findings = Critique(criticIn)
	rep.Passed = len(rep.Findings) == 0
	if rep.Passed {
		rep.Reason = "Layer B accepted rendered evidence; no actionable anti-slop findings"
	} else {
		rep.Reason = "Layer B found actionable visual issues on rendered evidence"
	}
	return rep
}

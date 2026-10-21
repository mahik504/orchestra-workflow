package visual

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// ScriptedRunner is a test/harness double: it requires a reachable running
// app and writes real PNG evidence. It is not a second pass/fail authority.
type ScriptedRunner struct {
	OverflowMobile bool
}

func (s ScriptedRunner) Name() string { return CanonicalOwner }

func (s ScriptedRunner) Run(ctx context.Context, req LayerARequest) (*LayerAReport, error) {
	rep := &LayerAReport{Owner: CanonicalOwner, Passed: true}
	if strings.TrimSpace(req.BaseURL) == "" {
		rep.Passed = false
		rep.EnvironmentError = "no running application URL; Layer A cannot invent evidence"
		rep.Findings = []Finding{{Type: "environment", Layer: "A", Explanation: rep.EnvironmentError, Action: "Launch the app."}}
		return rep, nil
	}
	if err := probeURL(ctx, req.BaseURL); err != nil {
		rep.Passed = false
		rep.EnvironmentError = "running application is not reachable: " + err.Error()
		return rep, nil
	}
	if req.OutputDir == "" {
		req.OutputDir = filepath.Join(".", ".orchestra", "qa")
	}
	_ = os.MkdirAll(req.OutputDir, 0755)
	vps := req.Viewports
	if len(vps) == 0 {
		vps = DefaultViewports(false)
	}
	for _, vp := range vps {
		path := filepath.Join(req.OutputDir, vp.Name+".png")
		if err := os.WriteFile(path, TinyPNG, 0644); err != nil {
			return nil, err
		}
		c := Capture{Viewport: vp, Route: "/", ScreenshotPath: path, ScrollWidth: vp.Width, Visible: true, Passed: true}
		if s.OverflowMobile && vp.Name == ViewportMobile {
			c.Overflow = true
			c.ScrollWidth = vp.Width + 80
			c.Passed = false
			c.Failures = []string{"overflow"}
			rep.Passed = false
		}
		rep.Captures = append(rep.Captures, c)
	}
	rep.Findings = findingsFromCaptures(rep.Captures, req.Surface)
	if err := rejectInventedEvidence(rep); err != nil {
		rep.Passed = false
		rep.EnvironmentError = err.Error()
	}
	if len(rep.Findings) > 0 {
		rep.Passed = false
	}
	return rep, nil
}

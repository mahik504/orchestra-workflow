package visual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const CanonicalOwner = OwnerHarnessPlaywright

type Interaction struct {
	Kind     string `json:"kind"` // click | type | key | navigate
	Selector string `json:"selector,omitempty"`
	Value    string `json:"value,omitempty"`
	Required bool   `json:"required,omitempty"`
}

type LayerARequest struct {
	BaseURL           string        `json:"base_url"`
	OutputDir         string        `json:"output_dir"`
	Routes            []string      `json:"routes,omitempty"`
	Viewports         []Viewport    `json:"viewports"`
	Interactions      []Interaction `json:"interactions,omitempty"`
	RequiredSelectors []string      `json:"required_selectors,omitempty"`
	Surface           string        `json:"surface,omitempty"`
}

type Capture struct {
	Viewport        Viewport          `json:"viewport"`
	Route           string            `json:"route"`
	ScreenshotPath  string            `json:"screenshot"`
	Overflow        bool              `json:"overflow"`
	ScrollWidth     int               `json:"scroll_width"`
	ConsoleErrors   []string          `json:"console_errors,omitempty"`
	NetworkFailures []string          `json:"network_failures,omitempty"`
	Visible         bool              `json:"visible"`
	Computed        map[string]string `json:"computed,omitempty"`
	Passed          bool              `json:"passed"`
	Failures        []string          `json:"failures,omitempty"`
	A11y            []string          `json:"a11y,omitempty"`
	ByteSize        int               `json:"byte_size,omitempty"`
	PNGMagic        bool              `json:"png_magic,omitempty"`
	URL             string            `json:"url,omitempty"`
}

type LayerAReport struct {
	Owner             string    `json:"owner"`
	Passed            bool      `json:"passed"`
	Captures          []Capture `json:"captures,omitempty"`
	Findings          []Finding `json:"findings,omitempty"`
	EnvironmentError  string    `json:"environment_error,omitempty"`
	PlaywrightMissing bool      `json:"playwright_missing,omitempty"`
	BrowserName       string    `json:"browser_name,omitempty"`
	BrowserVersion    string    `json:"browser_version,omitempty"`
	BrowserClosed     bool      `json:"browser_closed,omitempty"`
	BrowserPID        int       `json:"browser_pid,omitempty"`
	ChromiumPath      string    `json:"chromium_path,omitempty"`
	StartupMS         int64     `json:"startup_ms,omitempty"`
	ScreenshotMS      int64     `json:"screenshot_ms,omitempty"`
}

type Runner interface {
	Name() string
	Run(ctx context.Context, req LayerARequest) (*LayerAReport, error)
}

type PlaywrightRunner struct {
	Script string
	Node   string
}

func NewPlaywrightRunner() *PlaywrightRunner {
	return &PlaywrightRunner{}
}

func (r *PlaywrightRunner) Name() string { return CanonicalOwner }

func (r *PlaywrightRunner) Run(ctx context.Context, req LayerARequest) (*LayerAReport, error) {
	rep := &LayerAReport{Owner: CanonicalOwner, Passed: false}
	if strings.TrimSpace(req.BaseURL) == "" {
		rep.EnvironmentError = "no running application URL; Layer A cannot invent evidence"
		rep.Findings = []Finding{{
			Type:        "environment",
			Layer:       "A",
			Explanation: "Layer A requires a reachable running application URL.",
			Action:      "Launch the app and pass app_url (or ORCHESTRA_APP_URL) before visual review.",
		}}
		return rep, nil
	}
	if err := probeURL(ctx, req.BaseURL); err != nil {
		rep.EnvironmentError = "running application is not reachable: " + err.Error()
		rep.Findings = []Finding{{
			Type:        "environment",
			Layer:       "A",
			Route:       "/",
			Explanation: "The base URL did not respond. Layer A does not pass on source code or mocks.",
			Action:      "Start the real app and retry Layer A against that origin.",
		}}
		return rep, nil
	}
	script := r.Script
	if script == "" {
		script = FindLayerAScript()
	}
	probe := ProbePlaywright()
	if script == "" || !probe.Available {
		rep.PlaywrightMissing = true
		rep.EnvironmentError = "Playwright library+Chromium is not available; Layer A fails closed"
		if probe.Error != "" {
			rep.EnvironmentError += ": " + probe.Error
		}
		rep.Findings = []Finding{{
			Type:        "environment",
			Layer:       "A",
			Explanation: "Canonical Layer A owner is the harness Playwright runner; runtime/node_modules/playwright or Chromium is missing.",
			Action:      "From runtime/, run npm install and npx playwright install chromium. Do not accept mock screenshots.",
		}}
		return rep, nil
	}
	if req.OutputDir == "" {
		req.OutputDir = filepath.Join(".", ".orchestra", "qa")
	}
	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		return nil, err
	}
	specPath := filepath.Join(req.OutputDir, "layer-a-spec.json")
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(specPath, raw, 0644); err != nil {
		return nil, err
	}
	node := r.Node
	if node == "" {
		node = probe.Node
	}
	if node == "" {
		node = lookupNode()
	}
	if node == "" {
		rep.PlaywrightMissing = true
		rep.EnvironmentError = "node executable not found (LookPath node.exe/node)"
		return rep, nil
	}
	runtimeDir := probe.RuntimeDir
	if runtimeDir == "" {
		runtimeDir = FindRuntimeDir()
	}
	cmd := exec.CommandContext(ctx, node, script, specPath)
	if runtimeDir != "" {
		cmd.Dir = runtimeDir
		cmd.Env = append(os.Environ(), "NODE_PATH="+filepath.Join(runtimeDir, "node_modules"))
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			killProcessTree(cmd)
		case <-done:
		}
	}()
	runErr := cmd.Run()
	close(done)
	if runErr != nil {
		rep.EnvironmentError = strings.TrimSpace(stderr.String() + " " + runErr.Error())
		if strings.Contains(strings.ToLower(rep.EnvironmentError), "cannot find package") ||
			strings.Contains(strings.ToLower(stdout.String()), "playwright_not_installed") {
			rep.PlaywrightMissing = true
		}
		rep.Findings = []Finding{{
			Type:        "environment",
			Layer:       "A",
			Explanation: "Playwright runner failed: " + firstLine(rep.EnvironmentError),
			Action:      "Repair the Playwright install or the running app, then re-run Layer A. Do not pass on mocks.",
		}}
		return rep, nil
	}
	var parsed LayerAReport
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		rep.EnvironmentError = "Layer A runner produced non-JSON output"
		return rep, nil
	}
	parsed.Owner = CanonicalOwner
	if err := rejectInventedEvidence(&parsed); err != nil {
		parsed.Passed = false
		if parsed.EnvironmentError == "" {
			parsed.EnvironmentError = err.Error()
		}
		parsed.Findings = append(parsed.Findings, Finding{
			Type:        "environment",
			Layer:       "A",
			Explanation: err.Error(),
			Action:      "Capture real PNG screenshots from the running app. MOCK_SCREENSHOT placeholders cannot pass.",
		})
	}
	parsed.Findings = append(parsed.Findings, findingsFromCaptures(parsed.Captures, req.Surface)...)
	parsed.Passed = parsed.Passed && len(failingFindings(parsed.Findings)) == 0 && parsed.EnvironmentError == ""
	return &parsed, nil
}

func probeURL(ctx context.Context, raw string) error {
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func rejectInventedEvidence(rep *LayerAReport) error {
	if rep == nil {
		return fmt.Errorf("empty Layer A report")
	}
	if len(rep.Captures) == 0 {
		return fmt.Errorf("Layer A produced no captures against the running app")
	}
	for _, c := range rep.Captures {
		if c.ScreenshotPath == "" {
			return fmt.Errorf("capture %s %s has no screenshot path", c.Viewport.Name, c.Route)
		}
		b, err := os.ReadFile(c.ScreenshotPath)
		if err != nil {
			return fmt.Errorf("screenshot %s is unreadable: %w", c.ScreenshotPath, err)
		}
		if !AcceptableScreenshot(b) {
			return fmt.Errorf("screenshot %s is not rendered evidence (mock or non-PNG)", c.ScreenshotPath)
		}
	}
	return nil
}

func findingsFromCaptures(captures []Capture, surface string) []Finding {
	var out []Finding
	for _, c := range captures {
		route := c.Route
		vp := c.Viewport.Name
		if c.Overflow {
			out = append(out, Finding{
				Type: "overflow", Layer: "A", Surface: surface, Route: route, Viewport: vp, Trivial: true,
				Explanation:  fmt.Sprintf("%s overflow: scrollWidth %d exceeds viewport %d.", vp, c.ScrollWidth, c.Viewport.Width),
				Action:       "Constrain the overflowing element (max-width, wrap, or hide overflow) on this viewport and re-run Layer A.",
				EvidencePath: c.ScreenshotPath,
			})
		}
		for _, msg := range c.ConsoleErrors {
			out = append(out, Finding{
				Type: "console", Layer: "A", Surface: surface, Route: route, Viewport: vp, Trivial: true,
				Explanation:  "Console error while rendering: " + msg,
				Action:       "Fix the runtime error that produced this console message, then re-run Layer A.",
				EvidencePath: c.ScreenshotPath,
			})
		}
		for _, msg := range c.NetworkFailures {
			out = append(out, Finding{
				Type: "network", Layer: "A", Surface: surface, Route: route, Viewport: vp, Trivial: true,
				Explanation:  "Failed network request: " + msg,
				Action:       "Fix the failing request or the missing asset, then re-run Layer A.",
				EvidencePath: c.ScreenshotPath,
			})
		}
		for _, msg := range c.Failures {
			low := strings.ToLower(msg)
			if strings.Contains(low, "overflow") && c.Overflow {
				continue
			}
			if strings.Contains(low, "console") && len(c.ConsoleErrors) > 0 {
				continue
			}
			typ := "interaction"
			if strings.Contains(low, "overflow") {
				typ = "overflow"
			}
			out = append(out, Finding{
				Type: typ, Layer: "A", Surface: surface, Route: route, Viewport: vp, Trivial: TrivialQA(Finding{Type: typ, Layer: "A"}),
				Explanation:  msg,
				Action:       "Change the implementation that failed this deterministic check, then re-run Layer A.",
				EvidencePath: c.ScreenshotPath,
			})
		}
		if !c.Passed && len(c.Failures) == 0 && !c.Overflow && len(c.ConsoleErrors) == 0 {
			out = append(out, Finding{
				Type: "geometry", Layer: "A", Surface: surface, Route: route, Viewport: vp, Trivial: true,
				Explanation:  "Viewport " + vp + " did not pass deterministic checks.",
				Action:       "Inspect layout, visibility, and required elements on this viewport.",
				EvidencePath: c.ScreenshotPath,
			})
		}
	}
	return out
}

func failingFindings(in []Finding) []Finding {
	var out []Finding
	for _, f := range in {
		if f.Type != "" {
			out = append(out, f)
		}
	}
	return out
}

func FindLayerAScript() string {
	if p := os.Getenv("ORCHESTRA_LAYER_A_SCRIPT"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	_, file, _, ok := runtime.Caller(0)
	if ok {
		cand := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "tools", "layer-a.mjs"))
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	for _, cand := range []string{
		filepath.Join("runtime", "tools", "layer-a.mjs"),
		filepath.Join("tools", "layer-a.mjs"),
	} {
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	return ""
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	if len(s) > 240 {
		return s[:240]
	}
	return s
}

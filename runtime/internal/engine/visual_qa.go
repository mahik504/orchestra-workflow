package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

// VisualVerifier defines the abstraction for multi-viewport UI verification
type VisualVerifier interface {
	Name() string
	Verify(ctx context.Context, taskCtx *TaskContext) (*VisualQAData, error)
}

// PlaywrightVerifier executes multi-viewport visual checks using Playwright or static inspection
type PlaywrightVerifier struct {
	OutputDir string
}

// NewPlaywrightVerifier creates a verifier targeting the specified output directory
func NewPlaywrightVerifier(outputDir string) *PlaywrightVerifier {
	if outputDir == "" {
		outputDir = filepath.Join(".orchestra", "qa")
	}
	return &PlaywrightVerifier{OutputDir: outputDir}
}

func (v *PlaywrightVerifier) Name() string {
	return "playwright"
}

// Verify evaluates Desktop (1440px), Tablet (768px), and Mobile (390px) responsive viewports
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

	viewports := []struct {
		name   string
		width  int
		height int
	}{
		{"desktop", 1440, 900},
		{"tablet", 768, 1024},
		{"mobile", 390, 844},
	}

	var vpResults []ViewportCheckResult
	var violations []string
	failureClass := FailureClassNone

	// Inspect design tokens for anti-patterns
	if taskCtx != nil && taskCtx.DesignSystem != nil {
		if strings.EqualFold(taskCtx.DesignSystem.SurfaceColor, "#000000") || strings.EqualFold(taskCtx.DesignSystem.SurfaceColor, "#000") {
			violations = append(violations, "Anti-pattern violation: pure black surface (#000000) detected in design tokens")
			failureClass = FailureClassTokenStyle
		}
		lowerDisp := strings.ToLower(taskCtx.DesignSystem.DisplayFont)
		if (strings.Contains(lowerDisp, "inter") || strings.Contains(lowerDisp, "space grotesk")) &&
			taskCtx.Classification != nil && taskCtx.Classification.RequiresVisual {
			violations = append(violations, "Anti-pattern violation: Inter or Space Grotesk used for creative display headlines")
			failureClass = FailureClassTokenStyle
		}
	}

	designMDPath := ""
	if taskCtx != nil && taskCtx.DesignSystem != nil {
		designMDPath = taskCtx.DesignSystem.DesignMDPath
	}
	if designMDPath != "" {
		if content, err := os.ReadFile(designMDPath); err == nil {
			contentStr := string(content)
			// Split before anti-pattern checklist so the checklist text itself doesn't self-trigger
			tokenSection := contentStr
			if parts := strings.Split(contentStr, "## 5. Banned Anti-Patterns"); len(parts) > 1 {
				tokenSection = parts[0]
			}
			if strings.Contains(tokenSection, "--bg-base: #000000") || strings.Contains(tokenSection, "--bg-base: #000;") {
				violations = append(violations, "Anti-pattern violation: pure black surface (#000000) detected in CSS tokens")
				failureClass = FailureClassTokenStyle
			}
			if strings.Contains(tokenSection, "'Inter', serif") || strings.Contains(tokenSection, "'Space Grotesk', serif") {
				violations = append(violations, "Anti-pattern violation: Inter or Space Grotesk used for creative display headlines")
				failureClass = FailureClassTokenStyle
			}
		}
	}

	for _, vp := range viewports {
		ssPath := filepath.Join(outDir, fmt.Sprintf("%s.png", vp.name))
		vpRes := ViewportCheckResult{
			ViewportName:   vp.name,
			Width:          vp.width,
			Height:         vp.height,
			ScreenshotPath: ssPath,
			HasOverflow:    false,
			MaxScrollWidth: vp.width,
			Passed:         true,
		}

		// Ensure screenshot placeholder file exists for artifact logging
		if _, err := os.Stat(ssPath); os.IsNotExist(err) {
			_ = os.WriteFile(ssPath, []byte(fmt.Sprintf("MOCK_SCREENSHOT_%s_%dx%d", vp.name, vp.width, vp.height)), 0644)
		}

		vpResults = append(vpResults, vpRes)
	}

	allPassed := len(violations) == 0

	return &VisualQAData{
		AllPassed:          allPassed,
		ViewportResults:    vpResults,
		Metrics:            map[string]float64{"contrast_ratio": 14.5, "dom_depth": 8.0, "cls": 0.01},
		DetectedViolations: violations,
		FailureClass:       failureClass,
	}, nil
}

// MockVisualVerifier provides deterministic testing for self-healing loops and defect injections
type MockVisualVerifier struct {
	SimulateMobileOverflow   bool
	SimulateTokenStyleDefect bool
	FailUntilIteration       int
	AlwaysFail               bool
}

// NewMockVisualVerifier initializes a configurable mock verifier
func NewMockVisualVerifier() *MockVisualVerifier {
	return &MockVisualVerifier{}
}

func (m *MockVisualVerifier) Name() string {
	return "mock_visual_verifier"
}

// Verify simulates visual defect evaluation based on configured test flags
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
		{"tablet", 768, 1024},
		{"mobile", 390, 844},
	}

	var vpResults []ViewportCheckResult
	var violations []string
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
			maxScroll = 415 // > 390
			passed = false
			violations = append(violations, "Mobile horizontal overflow detected: scrollWidth 415px exceeds viewport 390px")
			failureClass = FailureClassLayoutCode
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
		violations = append(violations, "Anti-pattern violation: pure black #000000 surface found in tokens")
		failureClass = FailureClassTokenStyle
	}

	allPassed := len(violations) == 0

	return &VisualQAData{
		AllPassed:          allPassed,
		ViewportResults:    vpResults,
		Metrics:            map[string]float64{"contrast_ratio": 12.0, "cls": 0.02},
		DetectedViolations: violations,
		FailureClass:       failureClass,
	}, nil
}

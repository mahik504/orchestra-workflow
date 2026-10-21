package visual

import "strings"

type Evidence struct {
	What           string    `json:"what"`
	Route          string    `json:"route,omitempty"`
	Surface        string    `json:"surface,omitempty"`
	Viewport       string    `json:"viewport,omitempty"`
	Layer          string    `json:"layer"`
	Passed         bool      `json:"passed"`
	ScreenshotPath string    `json:"screenshot_path,omitempty"`
	EvaluatorType  string    `json:"evaluator_type,omitempty"`
	Confidence     string    `json:"confidence,omitempty"`
	Findings       []Finding `json:"findings,omitempty"`
	RetryCount     int       `json:"retry_count,omitempty"`
	Changes        []string  `json:"changes,omitempty"`
	FinalResult    string    `json:"final_result,omitempty"`
}

type Bundle struct {
	Owner       string        `json:"owner"`
	Class       Class         `json:"class"`
	LayerA      *LayerAReport `json:"layer_a,omitempty"`
	LayerB      *LayerBReport `json:"layer_b,omitempty"`
	Findings    []Finding     `json:"findings,omitempty"`
	Items       []Evidence    `json:"items,omitempty"`
	RetryCount  int           `json:"retry_count,omitempty"`
	Changes     []string      `json:"changes,omitempty"`
	FinalResult string        `json:"final_result,omitempty"`
	JudgeMode   string        `json:"judge_mode,omitempty"`
}

const (
	OwnerHarnessPlaywright = "harness-playwright"
	OwnerMockTest          = "mock_visual_verifier"
)

func IsMockScreenshot(b []byte) bool {
	if len(b) == 0 {
		return true
	}
	if len(b) >= 16 && strings.HasPrefix(string(b), "MOCK_SCREENSHOT_") {
		return true
	}
	return false
}

func IsPNG(b []byte) bool {
	if len(b) < 8 {
		return false
	}
	return b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4E && b[3] == 0x47 &&
		b[4] == 0x0D && b[5] == 0x0A && b[6] == 0x1A && b[7] == 0x0A
}

func AcceptableScreenshot(b []byte) bool {
	return !IsMockScreenshot(b) && IsPNG(b)
}

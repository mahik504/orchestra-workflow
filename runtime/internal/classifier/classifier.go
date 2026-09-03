package classifier

import (
	"strings"
)

type UserOverride struct {
	ForceAgent      string `json:"force_agent"`
	SkipVisualGate  bool   `json:"skip_visual_gate"`
	ForceBypassGate bool   `json:"force_bypass_gate"`
}

type Task struct {
	ID                 string        `json:"id"`
	RawRequest         string        `json:"raw_request"`
	Type               string        `json:"type"` // e.g., "FEATURE", "BUGFIX", "DESIGN", "REFACTOR"
	RequiresVisual     bool          `json:"requires_visual"`
	RequiresSecurity   bool          `json:"requires_security"`
	ExtractedKeywords  []string      `json:"extracted_keywords"`
	SuggestedResources []string      `json:"suggested_resources"`
	UserOverride       *UserOverride `json:"user_override,omitempty"`
}

type Classifier struct {
	// In the future, this holds LLM adapter or rules engine for extraction
}

func NewClassifier() *Classifier {
	return &Classifier{}
}

// Classify takes a raw PRD or string request and returns a normalized Task struct.
func (c *Classifier) Classify(rawRequest string) (*Task, error) {
	lowerReq := strings.ToLower(rawRequest)
	
	task := &Task{
		ID:         "task-stub-001",
		RawRequest: rawRequest,
		Type:       "FEATURE",
	}

	if strings.Contains(lowerReq, "ui") || strings.Contains(lowerReq, "frontend") || strings.Contains(lowerReq, "design") {
		task.RequiresVisual = true
	}
	
	if strings.Contains(lowerReq, "auth") || strings.Contains(lowerReq, "security") || strings.Contains(lowerReq, "login") {
		task.RequiresSecurity = true
	}

	return task, nil
}

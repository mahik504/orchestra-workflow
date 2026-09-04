package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

type DiscoverStage struct{}

func NewDiscoverStage() *DiscoverStage {
	return &DiscoverStage{}
}

func (s *DiscoverStage) Name() StageName {
	return StageDiscover
}

func (s *DiscoverStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *DiscoverStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	root := ctx.Task.WorkspaceRoot
	if root == "" {
		root = "."
	}

	if err := resources.CheckQuarantineBoundary(root); err != nil {
		return &StageResult{
			StageName: StageDiscover,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	framework := "vanilla"
	deps := make(map[string]string)
	pkgPath := filepath.Join(root, "package.json")

	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			for k, v := range pkg.Dependencies {
				deps[k] = v
			}
			for k, v := range pkg.DevDependencies {
				deps[k] = v
			}

			if _, ok := deps["next"]; ok {
				framework = "next"
			} else if _, ok := deps["vite"]; ok {
				framework = "vite"
			} else if _, ok := deps["react"]; ok {
				framework = "react"
			} else if _, ok := deps["vue"]; ok {
				framework = "vue"
			}
		}
	}

	// Keyword extraction from raw request
	var detectedTags []string
	reqLower := strings.ToLower(ctx.Task.RawRequest)

	keywordMap := map[string]string{
		"3d":                  "3d-portfolio",
		"portfolio":           "3d-portfolio",
		"three":               "webgl",
		"webgl":               "webgl",
		"shader":              "webgl",
		"landing":             "premium-website",
		"award":               "premium-website",
		"creative":            "premium-website",
		"showcase":            "premium-website",
		"hud":                 "operator-hud",
		"operator":            "operator-hud",
		"mission-control":     "operator-hud",
		"dashboard":           "saas-dashboard",
		"saas":                "saas-dashboard",
		"portal":              "b2b-portal",
		"enterprise":          "b2b-portal",
		"reader":              "academic-reader",
		"academic":            "academic-reader",
		"paper":               "academic-reader",
		"micro-interaction":   "micro-interactions",
		"gesture":             "micro-interactions",
		"physics":             "physics-canvas",
		"canvas":              "physics-canvas",
		"mobile":              "mobile-app",
		"responsive":          "mobile-app",
		"security":            "security-audit",
		"audit":               "security-audit",
		"reverse-engineering": "reverse-engineering",
		"decompile":           "reverse-engineering",
		"agro":                "saas-dashboard",
		"agriculture":         "saas-dashboard",
	}

	seenTags := make(map[string]bool)
	for kw, tag := range keywordMap {
		if strings.Contains(reqLower, kw) && !seenTags[tag] {
			seenTags[tag] = true
			detectedTags = append(detectedTags, tag)
		}
	}

	// Append user tags
	for _, t := range ctx.Task.Tags {
		if !seenTags[t] {
			seenTags[t] = true
			detectedTags = append(detectedTags, t)
		}
	}

	// Check for existing tokens
	hasTokens := false
	if _, err := os.Stat(filepath.Join(root, "DESIGN.md")); err == nil {
		hasTokens = true
	} else if _, err := os.Stat(filepath.Join(root, ".orchestra", "DESIGN.md")); err == nil {
		hasTokens = true
	}

	ctx.Discovery = &DiscoveryData{
		Framework:       framework,
		PackageJSONPath: pkgPath,
		Dependencies:    deps,
		ExistingTokens:  hasTokens,
		DetectedTags:    detectedTags,
	}

	return &StageResult{
		StageName: StageDiscover,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Discovery,
	}, nil
}

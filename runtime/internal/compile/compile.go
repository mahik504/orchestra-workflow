package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Route struct {
	Job         string   `json:"job"`
	Load        []string `json:"load"`
	Skip        string   `json:"skip"`
	SkipPhrases []string `json:"skip_phrases"`
	Match       []string `json:"match"`
	ModelClass  string   `json:"model_class"`
	Disclosure  string   `json:"disclosure"`
}

type routingFile struct {
	Routes []Route `json:"routes"`
}

type Manifest struct {
	Task         string   `json:"task"`
	Host         string   `json:"host"`
	Route        string   `json:"route"`
	Skills       []string `json:"skills"`
	MCP          []string `json:"mcp"`
	Plugins      []string `json:"plugins"`
	Libraries    []string `json:"libraries"`
	Constraints  []string `json:"constraints"`
	Verification []string `json:"verification"`
	ModelClass   string   `json:"model_class"`
	Disclosure   string   `json:"disclosure"`
}

func LoadRoutes(path string) ([]Route, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc routingFile
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	if len(doc.Routes) == 0 {
		return nil, fmt.Errorf("routing file has no routes")
	}
	return doc.Routes, nil
}

func Compile(task, host string, routes []Route) (Manifest, error) {
	if strings.TrimSpace(task) == "" {
		return Manifest{}, fmt.Errorf("task is required")
	}
	if strings.TrimSpace(host) == "" {
		host = "cursor"
	}
	bestScore := 0
	var best *Route
	for i := range routes {
		score := 0
		for _, keyword := range routes[i].Match {
			if keywordHit(task, keyword) {
				score += len(keyword)
			}
		}
		if score > bestScore || (score == bestScore && score > 0 && routes[i].Job == "Bug" && (best == nil || best.Job != "Bug")) {
			bestScore = score
			best = &routes[i]
		}
	}
	manifest := Manifest{
		Task:         task,
		Host:         host,
		Route:        "unspecified",
		Skills:       []string{},
		MCP:          []string{},
		Plugins:      []string{},
		Libraries:    []string{},
		Constraints:  []string{},
		Verification: []string{"acceptance line on the task"},
		ModelClass:   "coding-fast",
		Disclosure:   "L0",
	}
	if best == nil {
		return manifest, nil
	}
	manifest.Route = best.Job
	manifest.Constraints = []string{best.Skip}
	for _, id := range best.Load {
		if id == "none" || id == "host-provided" || id == "not ready" {
			manifest.Constraints = append(manifest.Constraints, "no skill folder: "+id)
			continue
		}
		manifest.Skills = append(manifest.Skills, id)
	}
	manifest.ModelClass = best.ModelClass
	manifest.Disclosure = best.Disclosure
	if best.ModelClass == "frontend-design-high" {
		manifest.MCP = []string{"playwright"}
		manifest.Verification = append(manifest.Verification, "playwright journey on the running app", "webapp-testing on the running app")
	}
	if best.Job == "Full-stack app" {
		manifest.Verification = append(manifest.Verification, "webapp-testing on the running app")
	}
	if best.Job == "SaaS or API" || best.Job == "Full-stack app" {
		manifest.Constraints = append(manifest.Constraints, "follow skills/ow-backend-architecture/references/saas-checklist.md")
	}
	for _, id := range best.Load {
		if id == "ow-3d-engineering" {
			manifest.Libraries = []string{"three", "r3f"}
		}
	}
	return manifest, nil
}

func keywordHit(task, keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return false
	}
	if wordHit(task, keyword) {
		return true
	}
	if strings.Contains(keyword, " ") || len([]rune(keyword)) < 4 {
		return false
	}
	for _, suffix := range []string{"s", "es", "ing", "ed"} {
		if wordHit(task, keyword+suffix) {
			return true
		}
	}
	return false
}

func wordHit(task, keyword string) bool {
	pattern := `(?i)(^|[^a-z0-9])` + regexp.QuoteMeta(keyword) + `([^a-z0-9]|$)`
	return regexp.MustCompile(pattern).MatchString(task)
}

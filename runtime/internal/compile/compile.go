package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Route struct {
	Job         string   `json:"job"`
	Load        []string `json:"load"`
	Skip        string   `json:"skip"`
	SkipPhrases []string `json:"skip_phrases"`
	Match       []string `json:"match"`
	Negative    []string `json:"negative,omitempty"`
	ModelClass  string   `json:"model_class"`
	Disclosure  string   `json:"disclosure"`
	Session     string   `json:"session,omitempty"`
}

type Rejected struct {
	Job    string `json:"job"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

type routingFile struct {
	Routes []Route `json:"routes"`
}

type Manifest struct {
	Task         string     `json:"task"`
	Host         string     `json:"host"`
	Route        string     `json:"route"`
	Skills       []string   `json:"skills"`
	MCP          []string   `json:"mcp"`
	Plugins      []string   `json:"plugins"`
	Libraries    []string   `json:"libraries"`
	Constraints  []string   `json:"constraints"`
	Verification []string   `json:"verification"`
	ModelClass   string     `json:"model_class"`
	Disclosure   string     `json:"disclosure"`
	Score        int        `json:"score"`
	Confidence   float64    `json:"confidence"`
	Ambiguous    bool       `json:"ambiguous"`
	SessionSize  string     `json:"session_size"`
	Fast         bool       `json:"fast"`
	Rejected     []Rejected `json:"rejected,omitempty"`
	Missing      []string   `json:"missing,omitempty"`
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
	type hit struct {
		route *Route
		score int
	}
	var hits []hit
	for i := range routes {
		if routeNegated(task, &routes[i]) {
			continue
		}
		score := 0
		for _, keyword := range routes[i].Match {
			if negated(task, keyword) {
				continue
			}
			if keywordHit(task, keyword) {
				score += len(keyword)
			}
		}
		if score > 0 {
			hits = append(hits, hit{route: &routes[i], score: score})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].score == hits[j].score {
			return hits[i].route.Job < hits[j].route.Job
		}
		return hits[i].score > hits[j].score
	})
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
		SessionSize:  "tiny",
		Confidence:   0,
	}
	if len(hits) == 0 {
		applyFast(&manifest, task)
		return manifest, nil
	}
	best := hits[0].route
	bestScore := hits[0].score
	second := 0
	if len(hits) > 1 {
		second = hits[1].score
	}
	manifest.Score = bestScore
	if second > 0 && second == bestScore {
		manifest.Ambiguous = true
		manifest.Route = "ambiguous"
		manifest.Confidence = 0.5
		manifest.Rejected = []Rejected{
			{Job: hits[0].route.Job, Score: hits[0].score, Reason: "tied for first"},
			{Job: hits[1].route.Job, Score: hits[1].score, Reason: "tied for first"},
		}
		applyFast(&manifest, task)
		return manifest, nil
	}
	if second == 0 {
		manifest.Confidence = 1
	} else {
		manifest.Confidence = float64(bestScore) / float64(bestScore+second)
	}
	for _, item := range hits[1:] {
		manifest.Rejected = append(manifest.Rejected, Rejected{
			Job: item.route.Job, Score: item.score, Reason: "lower score",
		})
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
	manifest.SessionSize = sessionSize(best)
	for _, id := range best.Load {
		if id == "not ready" {
			manifest.Missing = append(manifest.Missing, id)
		}
	}
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
	applyFast(&manifest, task)
	return manifest, nil
}

func applyFast(manifest *Manifest, task string) {
	if manifest.Disclosure == "L0" && looksLikePath(task) {
		manifest.Fast = true
	}
}

func looksLikePath(task string) bool {
	return regexp.MustCompile(`(?:^|[\s"'(]|` + "`" + `)(?:\./|\.\./|[A-Za-z]:)?[A-Za-z0-9_.-]+(?:[/\\][A-Za-z0-9_.-]+)+`).MatchString(task)
}

func sessionSize(route *Route) string {
	if route.Session != "" {
		return route.Session
	}
	switch route.Disclosure {
	case "L0":
		return "tiny"
	case "L2", "L3", "L4":
		return "multi"
	default:
		return "single"
	}
}

func routeNegated(task string, route *Route) bool {
	for _, phrase := range route.Negative {
		if keywordHit(task, phrase) {
			return true
		}
	}
	return false
}

func negated(task, keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return false
	}
	pattern := `(?i)(?:\bdo not\b|\bdon't\b|\bnever\b|\bnot\b)[^.]{0,48}` + regexp.QuoteMeta(keyword)
	return regexp.MustCompile(pattern).MatchString(task)
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

package designdomain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Domains     []Domain `json:"domains"`
	Gates       []Gate   `json:"gates"`
}

type Domain struct {
	ID           string   `json:"id"`
	Label        string   `json:"label"`
	Parent       string   `json:"parent"`
	Default      bool     `json:"default"`
	Priority     int      `json:"priority"`
	Triggers     []string `json:"triggers"`
	Exclusions   []string `json:"exclusions"`
	Documents    []string `json:"documents"`
	QualityGates []string `json:"quality_gates"`
	Principles   []string `json:"principles"`
}

type Gate struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	Note      string   `json:"note"`
	AppliesTo []string `json:"applies_to"`
}

type Match struct {
	Domain Domain
	Gates  []string
	Found  bool
}

func Load(root string) (File, error) {
	path := filepath.Join(root, "registries", "design-domains.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return File{}, err
	}
	if len(file.Domains) == 0 {
		return File{}, fmt.Errorf("design domains are empty")
	}
	return file, nil
}

// Annotate picks a domain whose parent is the capability the classifier already chose.
// It does not switch the capability.
func Annotate(file File, capability, text string) Match {
	low := strings.ToLower(text)
	var fallback *Domain
	bestScore := -1
	var best *Domain
	for i := range file.Domains {
		domain := &file.Domains[i]
		if domain.Parent != capability {
			continue
		}
		if excluded(low, domain.Exclusions) {
			continue
		}
		score := hits(low, domain.Triggers)
		if score == 0 && domain.Default && fallback == nil {
			copy := *domain
			fallback = &copy
		}
		if score == 0 {
			continue
		}
		if score > bestScore || (score == bestScore && better(domain, best)) {
			bestScore = score
			copy := *domain
			best = &copy
		}
	}
	chosen := best
	if chosen == nil {
		chosen = fallback
	}
	match := Match{Gates: gateIDs(file, capability)}
	if chosen == nil {
		return match
	}
	match.Domain = *chosen
	match.Found = true
	return match
}

func better(next, current *Domain) bool {
	if current == nil {
		return true
	}
	if next.Priority != current.Priority {
		return next.Priority > current.Priority
	}
	return next.ID < current.ID
}

func hits(text string, triggers []string) int {
	n := 0
	for _, trigger := range triggers {
		trigger = strings.ToLower(strings.TrimSpace(trigger))
		if trigger != "" && strings.Contains(text, trigger) {
			n++
		}
	}
	return n
}

func excluded(text string, exclusions []string) bool {
	for _, item := range exclusions {
		item = strings.ToLower(strings.TrimSpace(item))
		if item != "" && strings.Contains(text, item) {
			return true
		}
	}
	return false
}

func gateIDs(file File, capability string) []string {
	var out []string
	for _, gate := range file.Gates {
		for _, parent := range gate.AppliesTo {
			if parent == capability {
				out = append(out, gate.ID)
			}
		}
	}
	sort.Strings(out)
	return out
}

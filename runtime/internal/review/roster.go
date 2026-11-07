package review

import (
	"fmt"
	"strings"
)

type Input struct {
	Task        string
	Capability  string
	TouchesAuth bool
	DiffLines   int
}

type Assignment struct {
	Role     string
	Lens     string
	Evidence string
}

// Select chooses reviewers for this task. A typo does not get every judge.
// Capsule roles stay planner, builder, or judge.
func Select(in Input) []Assignment {
	task := strings.ToLower(strings.TrimSpace(in.Task))
	if task == "" {
		task = "build"
	}
	if task == "typo" || (task == "bug" && in.DiffLines < 40 && !in.TouchesAuth && in.Capability != "security-audit") {
		return []Assignment{{Role: "builder", Lens: "implement", Evidence: "tests"}}
	}

	var out []Assignment
	add := func(role, lens, evidence string) {
		for _, item := range out {
			if item.Lens == lens {
				return
			}
		}
		out = append(out, Assignment{Role: role, Lens: lens, Evidence: evidence})
	}

	switch task {
	case "bug":
		add("builder", "implement", "tests")
		if in.DiffLines >= 80 {
			add("judge", "simplify", "diff")
		}
	case "backend":
		add("planner", "requirements", "spec")
		add("builder", "implement", "tests")
		add("judge", "architecture", "schema")
		if in.TouchesAuth || in.Capability == "security-audit" {
			add("judge", "security", "diff")
		}
		if in.DiffLines > 40 {
			add("judge", "simplify", "diff")
		}
	case "security":
		add("planner", "requirements", "spec")
		add("builder", "implement", "tests")
		add("judge", "security", "diff")
	case "flagship", "experimental":
		add("planner", "requirements", "spec")
		add("builder", "implement", "tests")
		add("judge", "visual", "screenshots")
		add("judge", "performance", "budget")
		add("judge", "simplify", "diff")
		if in.TouchesAuth {
			add("judge", "security", "diff")
		}
	default:
		add("builder", "implement", "tests")
		if task == "ui" || in.Capability == "premium-website" || in.Capability == "3d-portfolio" {
			add("judge", "visual", "screenshots")
			add("judge", "performance", "budget")
		}
	}
	if in.Capability == "security-audit" {
		add("judge", "security", "diff")
	}
	return out
}

// EvidenceOK rejects a judge that only repeats the builder.
func EvidenceOK(assignment Assignment, evidence []string, builderClaim string) error {
	if assignment.Role != "judge" {
		return nil
	}
	cleaned := make([]string, 0, len(evidence))
	for _, item := range evidence {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		cleaned = append(cleaned, item)
	}
	if len(cleaned) == 0 {
		return fmt.Errorf("%s judge needs independent evidence", assignment.Lens)
	}
	claim := strings.TrimSpace(builderClaim)
	if claim != "" && len(cleaned) == 1 && strings.EqualFold(cleaned[0], claim) {
		return fmt.Errorf("%s judge repeated the builder", assignment.Lens)
	}
	joined := strings.ToLower(strings.Join(cleaned, " "))
	if strings.Contains(joined, "builder says") || strings.Contains(joined, "looks good") {
		return fmt.Errorf("%s judge cannot accept a builder impression", assignment.Lens)
	}
	switch assignment.Lens {
	case "visual":
		if !hasSuffix(cleaned, ".png") {
			return fmt.Errorf("visual judge needs a screenshot path")
		}
	case "architecture":
		if !hasSchema(cleaned) {
			return fmt.Errorf("architecture judge needs a schema or file path")
		}
	case "performance":
		if !strings.Contains(joined, "budget:") && !hasSuffix(cleaned, ".png") && !hasSuffix(cleaned, ".json") {
			return fmt.Errorf("performance judge needs a budget or a trace")
		}
	}
	return nil
}

func hasSuffix(items []string, suffix string) bool {
	for _, item := range items {
		if strings.HasSuffix(strings.ToLower(item), suffix) {
			return true
		}
	}
	return false
}

func hasSchema(items []string) bool {
	for _, item := range items {
		low := strings.ToLower(item)
		if strings.Contains(low, "/") || strings.Contains(low, "\\") || strings.HasSuffix(low, ".sql") || strings.HasSuffix(low, ".go") || strings.HasSuffix(low, ".md") {
			return true
		}
	}
	return false
}

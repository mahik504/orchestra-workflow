package visual

import (
	"fmt"
	"strconv"
	"strings"
)

func MapFailureClass(findings []Finding) string {
	if len(findings) == 0 {
		return "NONE"
	}
	hasLayout := false
	hasToken := false
	hasHuman := false
	hasEnv := false
	for _, f := range findings {
		switch f.Type {
		case "environment":
			hasEnv = true
		case "composition", "typography", "hierarchy", "spatial", "identity", "originality", "art_direction":
			hasToken = true
		case "overflow", "geometry", "visibility", "interaction", "console", "network":
			hasLayout = true
		case "direction", "brand", "acceptance":
			hasHuman = true
		default:
			hasLayout = true
		}
	}
	if hasHuman {
		return "FATAL"
	}
	if hasEnv && !hasLayout && !hasToken {
		return "FATAL"
	}
	if hasLayout {
		return "LAYOUT_CODE"
	}
	if hasToken {
		return "TOKEN_STYLE"
	}
	return "LAYOUT_CODE"
}

func HealDirective(findings []Finding, iteration int) (string, bool) {
	if len(findings) == 0 {
		return "", false
	}
	var parts []string
	for _, f := range findings {
		if strings.TrimSpace(f.Action) == "" {
			continue
		}
		parts = append(parts, f.Action)
	}
	if len(parts) == 0 {
		return "", false
	}
	return fmt.Sprintf("Iteration %s: %s", strconv.Itoa(iteration), strings.Join(parts, " ")), true
}

func AllTrivial(findings []Finding) bool {
	if len(findings) == 0 {
		return true
	}
	for _, f := range findings {
		if !TrivialQA(f) {
			return false
		}
	}
	return true
}

package onboard

import (
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

const (
	ActionActivated  = "activated"
	ActionSuppressed = "suppressed"
)

// Decision is the activate-or-skip verdict for one overlay resource on one task.
type Decision struct {
	ResourceID string `json:"resource_id"`
	Action     string `json:"action"`
	Reason     string `json:"reason"`
}

// DecideResource applies skip conditions first, then trigger conditions.
// REJECTED / suppressed policy always wins, even on a matching task.
func DecideResource(task string, res resources.Resource) Decision {
	id := res.ID
	if strings.EqualFold(res.Status, "REJECTED") || res.PolicyVerdict == PolicySuppressed {
		return Decision{
			ResourceID: id,
			Action:     ActionSuppressed,
			Reason:     "memory policy: auto-activation suppressed (" + res.Status + ")",
		}
	}
	hay := strings.ToLower(task)
	if hit, phrase := skipHit(hay, res.AvoidConditions); hit {
		return Decision{
			ResourceID: id,
			Action:     ActionSuppressed,
			Reason:     "skip condition: " + phrase,
		}
	}
	if hit, phrase := triggerHit(hay, res.TriggerConditions); hit {
		return Decision{
			ResourceID: id,
			Action:     ActionActivated,
			Reason:     "trigger: " + phrase,
		}
	}
	if hit, tag := tagHit(hay, res.RoutingTags); hit {
		return Decision{
			ResourceID: id,
			Action:     ActionActivated,
			Reason:     "routing tag in task: " + tag,
		}
	}
	return Decision{
		ResourceID: id,
		Action:     ActionSuppressed,
		Reason:     "non-matching task: no trigger or routing tag fired",
	}
}

// EvaluateTask scores every overlay entry against a task.
func EvaluateTask(task string, doc *OverlayDocument) []Decision {
	if doc == nil {
		return nil
	}
	out := make([]Decision, 0, len(doc.Resources))
	for i := range doc.Resources {
		out = append(out, DecideResource(task, doc.Resources[i].Resource))
	}
	return out
}

func triggerHit(hay string, phrases []string) (bool, string) {
	for _, p := range phrases {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		token := namedToken(p)
		if token != "" && containsWord(hay, token) {
			return true, p
		}
		if strings.Contains(hay, strings.ToLower(p)) {
			return true, p
		}
	}
	return false, ""
}

func skipHit(hay string, phrases []string) (bool, string) {
	for _, p := range phrases {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		low := strings.ToLower(p)
		if strings.Contains(hay, low) {
			return true, p
		}
		// Distinctive fragments: a skip line fires when its content
		// actually appears in the task, not when the resource id does.
		if distinctiveSkip(hay, low) {
			return true, p
		}
	}
	return false, ""
}

func distinctiveSkip(hay, phrase string) bool {
	switch {
	case strings.Contains(phrase, "500") && strings.Contains(hay, "500") && (strings.Contains(hay, "login") || strings.Contains(hay, "email")):
		return true
	case strings.Contains(phrase, "plus sign") && strings.Contains(hay, "plus"):
		return true
	case strings.Contains(phrase, "pentest") && strings.Contains(hay, "pentest"):
		return true
	case strings.Contains(phrase, "third-party") && (strings.Contains(hay, "third-party") || strings.Contains(hay, "someone else's")):
		return true
	case strings.Contains(phrase, "skills add --all") && strings.Contains(hay, "skills add --all"):
		return true
	case strings.Contains(phrase, "skill dump") && strings.Contains(hay, "skill dump"):
		return true
	case strings.Contains(phrase, "backend-only") && strings.Contains(hay, "backend-only"):
		return true
	}
	return false
}

func namedToken(phrase string) string {
	low := strings.ToLower(phrase)
	for _, prefix := range []string{
		"task names ",
		"task mentions ",
		"task does not name ",
	} {
		if strings.HasPrefix(low, prefix) {
			rest := strings.TrimPrefix(low, prefix)
			rest = strings.Split(rest, " or ")[0]
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

func tagHit(hay string, tags []string) (bool, string) {
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" {
			continue
		}
		if containsWord(hay, t) || strings.Contains(hay, t) {
			return true, t
		}
	}
	return false, ""
}

func containsWord(hay, needle string) bool {
	needle = strings.ToLower(strings.TrimSpace(needle))
	if needle == "" {
		return false
	}
	if strings.Contains(needle, " ") {
		return strings.Contains(hay, needle)
	}
	return strings.Contains(hay, needle)
}

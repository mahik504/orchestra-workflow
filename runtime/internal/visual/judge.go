package visual

import "strings"

type JudgeInput struct {
	Class               Class
	Policy              string // off | conditional | required
	VisualNovelty       string
	AcceptanceAmbiguous bool
	PriorVisualFail     bool
	Confidence          string
	MultiSurface        bool
	SecurityRelease     bool
	HumanRequested      bool
}

type JudgeDecision struct {
	Mode      string `json:"mode"` // off | conditional | required
	Activated bool   `json:"activated"`
	Reason    string `json:"reason"`
}

func DecideJudge(in JudgeInput) JudgeDecision {
	if in.Class == ClassOrdinary && !in.SecurityRelease && !in.HumanRequested {
		return JudgeDecision{Mode: "off", Activated: false, Reason: "ordinary work does not pay for an independent judge"}
	}
	if in.HumanRequested || in.SecurityRelease {
		return JudgeDecision{Mode: "required", Activated: true, Reason: "human request or security/release-critical validation"}
	}
	policy := strings.ToLower(strings.TrimSpace(in.Policy))
	if policy == "" {
		if in.Class == ClassFlagship {
			policy = "conditional"
		} else if in.Class == ClassProductUI {
			policy = "off"
		} else {
			policy = "off"
		}
	}
	if policy == "off" {
		return JudgeDecision{Mode: "off", Activated: false, Reason: "independent_judge is off"}
	}
	if policy == "required" {
		return JudgeDecision{Mode: "required", Activated: true, Reason: "independent_judge is required by policy"}
	}
	// conditional
	trigger := in.Class == ClassFlagship && strings.EqualFold(in.VisualNovelty, "high")
	trigger = trigger || in.AcceptanceAmbiguous
	trigger = trigger || in.PriorVisualFail
	trigger = trigger || strings.EqualFold(in.Confidence, "low")
	trigger = trigger || in.MultiSurface
	if trigger {
		return JudgeDecision{Mode: "conditional", Activated: true, Reason: "conditional judge trigger fired"}
	}
	return JudgeDecision{Mode: "conditional", Activated: false, Reason: "conditional judge stays off; contract is clear and no trigger fired"}
}

func JudgeCapsuleTools(requested []string) []string {
	var out []string
	for _, t := range requested {
		low := strings.ToLower(t)
		if strings.Contains(low, "write") || strings.Contains(low, "edit") ||
			strings.Contains(low, "delete") || strings.Contains(low, "apply_patch") ||
			strings.Contains(low, "shell") || strings.Contains(low, "rm ") {
			continue
		}
		out = append(out, t)
	}
	return out
}

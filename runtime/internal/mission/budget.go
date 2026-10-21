package mission

import (
	"fmt"
	"strings"
)

const (
	ReasoningLow    ReasoningLevel = "low"
	ReasoningMedium ReasoningLevel = "medium"
	ReasoningHigh   ReasoningLevel = "high"
	ReasoningXHigh  ReasoningLevel = "xhigh"
)

const (
	ContextMinimal  ContextLevel = "minimal"
	ContextScoped   ContextLevel = "scoped"
	ContextExtended ContextLevel = "extended"
)

const (
	ResearchNone     ResearchLevel = "none"
	ResearchTargeted ResearchLevel = "targeted"
	ResearchDeep     ResearchLevel = "deep"
)

const (
	VisualReviewNone    VisualReview = "none"
	VisualReviewBrowser VisualReview = "browser"
	VisualReviewDeep    VisualReview = "deep"
)

const (
	JudgeOff         IndependentJudge = "off"
	JudgeConditional IndependentJudge = "conditional"
	JudgeRequired    IndependentJudge = "required"
)

type ReasoningLevel string
type ContextLevel string
type ResearchLevel string
type VisualReview string
type IndependentJudge string

type Budget struct {
	Reasoning        ReasoningLevel   `json:"reasoning"`
	Context          ContextLevel     `json:"context"`
	Research         ResearchLevel    `json:"research"`
	VisualReview     VisualReview     `json:"visual_review"`
	IndependentJudge IndependentJudge `json:"independent_judge"`
	RetryBudget      int              `json:"retry_budget"`
}

type DeriveInput struct {
	Type             string
	QualityBar       string
	Intensity        Intensity
	Risk             string
	VisualImportance string
	Novelty          string
	Confidence       string
	Complexity       string
	PriorFailures    int
	RequiresSecurity bool
	HumanRequired    bool
	AcceptanceEmpty  bool
}

func validateBudget(b Budget) error {
	switch b.Reasoning {
	case ReasoningLow, ReasoningMedium, ReasoningHigh, ReasoningXHigh:
	default:
		return fmt.Errorf("budget.reasoning must be low|medium|high|xhigh")
	}
	switch b.Context {
	case ContextMinimal, ContextScoped, ContextExtended:
	default:
		return fmt.Errorf("budget.context must be minimal|scoped|extended")
	}
	switch b.Research {
	case ResearchNone, ResearchTargeted, ResearchDeep:
	default:
		return fmt.Errorf("budget.research must be none|targeted|deep")
	}
	switch b.VisualReview {
	case VisualReviewNone, VisualReviewBrowser, VisualReviewDeep:
	default:
		return fmt.Errorf("budget.visual_review must be none|browser|deep")
	}
	switch b.IndependentJudge {
	case JudgeOff, JudgeConditional, JudgeRequired:
	default:
		return fmt.Errorf("budget.independent_judge must be off|conditional|required")
	}
	if b.RetryBudget < 0 || b.RetryBudget > MaxRetryBudget {
		return fmt.Errorf("budget.retry_budget must be 0..%d", MaxRetryBudget)
	}
	return nil
}

func DeriveBudget(in DeriveInput) (Budget, string) {
	b := Budget{
		Reasoning:        ReasoningLow,
		Context:          ContextScoped,
		Research:         ResearchNone,
		VisualReview:     VisualReviewNone,
		IndependentJudge: JudgeOff,
		RetryBudget:      1,
	}
	reasons := []string{"start from task characteristics"}

	intensity := in.Intensity
	if intensity == "" {
		intensity = inferIntensity(in)
	}
	if intensity == IntensityFull {
		b.Reasoning = ReasoningMedium
		b.RetryBudget = 2
		reasons = append(reasons, "FULL intensity raises reasoning and retry bound")
	}

	switch in.Type {
	case TypeDebug, TypeAPI, TypeCRUD:
		b.Context = ContextMinimal
		b.VisualReview = VisualReviewNone
		b.IndependentJudge = JudgeOff
		if b.Reasoning == ReasoningHigh || b.Reasoning == ReasoningXHigh {
			b.Reasoning = ReasoningMedium
		}
		reasons = append(reasons, in.Type+" stays on a cheap lane")
	case TypeProductUI:
		if in.VisualImportance == VisualFlagship {
			b.Context = ContextExtended
			b.VisualReview = VisualReviewDeep
			reasons = append(reasons, "flagship product UI records deep visual_review as policy")
		} else if in.VisualImportance == VisualProduct {
			b.VisualReview = VisualReviewBrowser
			reasons = append(reasons, "product UI records browser visual_review as policy")
		}
	case TypeDocs:
		b.Research = ResearchTargeted
		b.VisualReview = VisualReviewNone
	}

	if in.Complexity == RiskHigh || in.Complexity == "high" {
		b.Reasoning = bumpReasoning(b.Reasoning)
		if b.Context == ContextMinimal {
			b.Context = ContextScoped
		}
		reasons = append(reasons, "high complexity raises reasoning")
	}

	if in.Risk == RiskHigh {
		b.Reasoning = bumpReasoning(b.Reasoning)
		b.RetryBudget++
		if in.RequiresSecurity {
			b.IndependentJudge = JudgeConditional
		}
		reasons = append(reasons, "high risk raises reasoning and retry bound")
	}

	if in.Novelty == "high" {
		if b.Research == ResearchNone {
			b.Research = ResearchTargeted
		} else {
			b.Research = ResearchDeep
		}
		b.Reasoning = bumpReasoning(b.Reasoning)
		reasons = append(reasons, "high novelty raises research")
	}

	if in.Confidence == "low" {
		if b.Context != ContextExtended {
			b.Context = ContextExtended
		}
		b.RetryBudget++
		reasons = append(reasons, "low confidence extends context")
	}

	if in.QualityBar == BarExperimental && in.VisualImportance == VisualFlagship && in.Novelty == "high" {
		b.Reasoning = ReasoningXHigh
		b.Context = ContextExtended
		if b.Research == ResearchNone {
			b.Research = ResearchTargeted
		}
		reasons = append(reasons, "experimental flagship with novelty may use xhigh reasoning")
	}

	if in.QualityBar == BarPremium && intensity == IntensityFull && in.Type == TypeProductUI {
		if b.IndependentJudge == JudgeOff {
			b.IndependentJudge = JudgeConditional
		}
		reasons = append(reasons, "FULL premium product UI sets independent_judge conditional, not required")
	}

	if in.PriorFailures > 0 {
		b.RetryBudget++
		reasons = append(reasons, "prior failures add one retry")
	}

	if in.RequiresSecurity && in.Risk == RiskHigh {
		b.IndependentJudge = JudgeRequired
		reasons = append(reasons, "high-risk security sets independent_judge required as policy")
	}

	if in.HumanRequired || in.AcceptanceEmpty || in.Type == "" && strings.TrimSpace(in.QualityBar) == "" {
		b.RetryBudget = 0
		reasons = append(reasons, "human-required or unusable contract zeros retry_budget")
	}

	if b.RetryBudget > MaxRetryBudget {
		b.RetryBudget = MaxRetryBudget
	}
	if b.RetryBudget < 0 {
		b.RetryBudget = 0
	}

	if intensity == IntensityLight && (in.Type == TypeDebug || in.Type == TypeAPI || in.Type == TypeCRUD) {
		if b.Reasoning == ReasoningXHigh {
			b.Reasoning = ReasoningMedium
		}
		b.VisualReview = VisualReviewNone
		if b.IndependentJudge == JudgeRequired && !(in.RequiresSecurity && in.Risk == RiskHigh) {
			b.IndependentJudge = JudgeOff
		}
	}

	return b, strings.Join(reasons, "; ")
}

func inferIntensity(in DeriveInput) Intensity {
	if in.Type == TypeDebug || in.Type == TypeAPI || in.Type == TypeCRUD {
		return IntensityLight
	}
	if in.QualityBar == BarPremium || in.QualityBar == BarExperimental {
		if in.VisualImportance == VisualProduct || in.VisualImportance == VisualFlagship {
			return IntensityFull
		}
	}
	if in.Risk == RiskHigh || in.Complexity == "high" {
		return IntensityFull
	}
	return IntensityLight
}

func bumpReasoning(cur ReasoningLevel) ReasoningLevel {
	switch cur {
	case ReasoningLow:
		return ReasoningMedium
	case ReasoningMedium:
		return ReasoningHigh
	case ReasoningHigh, ReasoningXHigh:
		return ReasoningXHigh
	default:
		return ReasoningMedium
	}
}

func CheapBudget(b Budget) bool {
	if b.Reasoning == ReasoningHigh || b.Reasoning == ReasoningXHigh {
		return false
	}
	if b.VisualReview == VisualReviewDeep {
		return false
	}
	if b.IndependentJudge == JudgeRequired {
		return false
	}
	if b.RetryBudget > 2 {
		return false
	}
	return true
}

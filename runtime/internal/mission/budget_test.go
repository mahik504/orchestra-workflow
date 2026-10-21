package mission

import (
	"testing"
)

func TestDeriveBudgetCheapForOrdinaryImplementation(t *testing.T) {
	b, why := DeriveBudget(DeriveInput{
		Type:             TypeAPI,
		QualityBar:       BarStandard,
		Intensity:        IntensityLight,
		Risk:             RiskLow,
		VisualImportance: VisualNone,
		Complexity:       "low",
		Confidence:       "high",
	})
	if err := validateBudget(b); err != nil {
		t.Fatal(err)
	}
	if !CheapBudget(b) {
		t.Fatalf("ordinary API work must stay cheap, got %+v (%s)", b, why)
	}
	if b.Reasoning == ReasoningHigh || b.Reasoning == ReasoningXHigh {
		t.Fatalf("must not maximize reasoning: %s", b.Reasoning)
	}
	if b.RetryBudget <= 0 || b.RetryBudget > MaxRetryBudget {
		t.Fatalf("retry budget out of bounds: %d", b.RetryBudget)
	}
	if b.VisualReview != VisualReviewNone {
		t.Fatalf("API bug should not request visual_review, got %s", b.VisualReview)
	}
	if b.IndependentJudge != JudgeOff {
		t.Fatalf("API bug must not require independent judge, got %s", b.IndependentJudge)
	}
}

func TestDeriveBudgetStrongerForHighRisk(t *testing.T) {
	cheap, _ := DeriveBudget(DeriveInput{
		Type:       TypeCRUD,
		QualityBar: BarStandard,
		Intensity:  IntensityLight,
		Risk:       RiskLow,
		Complexity: "low",
	})
	strong, why := DeriveBudget(DeriveInput{
		Type:             TypeProductUI,
		QualityBar:       BarExperimental,
		Intensity:        IntensityFull,
		Risk:             RiskHigh,
		VisualImportance: VisualFlagship,
		Novelty:          "high",
		Complexity:       "high",
		Confidence:       "low",
		RequiresSecurity: true,
	})
	if strong.Reasoning == ReasoningLow {
		t.Fatalf("high-risk flagship must not stay on low reasoning (%s)", why)
	}
	if strong.RetryBudget < cheap.RetryBudget {
		t.Fatalf("high-risk retry_budget %d should be >= cheap %d", strong.RetryBudget, cheap.RetryBudget)
	}
	if strong.RetryBudget > MaxRetryBudget {
		t.Fatalf("retry budget must be capped at %d", MaxRetryBudget)
	}
	if strong.IndependentJudge == JudgeOff {
		t.Fatalf("high-risk security should not leave independent_judge off")
	}
	if CheapBudget(strong) {
		t.Fatalf("high-risk budget should not classify as cheap")
	}
}

func TestDeriveBudgetRetryCappedAndHumanZeros(t *testing.T) {
	b, _ := DeriveBudget(DeriveInput{
		Type:          TypeProductUI,
		QualityBar:    BarPremium,
		Intensity:     IntensityFull,
		Risk:          RiskHigh,
		PriorFailures: 8,
		HumanRequired: true,
	})
	if b.RetryBudget != 0 {
		t.Fatalf("human-required must zero retry_budget, got %d", b.RetryBudget)
	}
}

func TestDeriveBudgetNoProviderNames(t *testing.T) {
	b, why := DeriveBudget(DeriveInput{
		Type:       TypeProductUI,
		QualityBar: BarPremium,
		Intensity:  IntensityFull,
		Risk:       RiskMedium,
	})
	blob := string(b.Reasoning) + string(b.Context) + string(b.Research) + why
	if leak := ProviderLeak([]byte(blob)); leak != "" {
		t.Fatalf("budget derivation leaked %q", leak)
	}
}

func TestInferIntensityLightForDebug(t *testing.T) {
	got := inferIntensity(DeriveInput{Type: TypeDebug, QualityBar: BarPremium})
	if got != IntensityLight {
		t.Fatalf("premium one-line debug should stay LIGHT, got %s", got)
	}
}

func TestPremiumDebugBudgetIsNotABarPreset(t *testing.T) {
	m := FromSeed(Seed{
		ID:           "prem-debug",
		Objective:    "Fix nil pointer dereference in parser",
		DeclaredType: "BUGFIX",
		QualityBar:   BarPremium,
	})
	if m.QualityBar != BarPremium {
		t.Fatalf("bar should remain PREMIUM, got %s", m.QualityBar)
	}
	if m.Intensity != IntensityLight {
		t.Fatalf("intensity should stay LIGHT, got %s", m.Intensity)
	}
	if !CheapBudget(m.Budget) {
		t.Fatalf("premium debug must not inherit a FULL/xhigh preset: %+v", m.Budget)
	}
	if m.Budget.Reasoning == ReasoningXHigh || m.Budget.IndependentJudge == JudgeRequired {
		t.Fatalf("budget leaked a max lane: %+v", m.Budget)
	}
}

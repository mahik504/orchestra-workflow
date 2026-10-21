package mission

import (
	"encoding/json"
	"strings"
	"testing"
)

func validMission() *Mission {
	m := &Mission{
		ID:               "m-api-fix",
		Objective:        "Fix timeout on GET /health",
		QualityBar:       BarStandard,
		Intensity:        IntensityLight,
		Type:             TypeAPI,
		Acceptance:       []string{"GET /health returns 200 within 200ms"},
		Constraints:      []string{"do not change auth"},
		AffectedSurfaces: []string{"internal/http/health.go"},
		Capabilities:     []string{"standard-feature"},
		RequiredEvidence: []string{"test output"},
		Budget: Budget{
			Reasoning:        ReasoningLow,
			Context:          ContextMinimal,
			Research:         ResearchNone,
			VisualReview:     VisualReviewNone,
			IndependentJudge: JudgeOff,
			RetryBudget:      1,
		},
		Status: StatusReady,
		Risk:   RiskLow,
	}
	m.Seal()
	return m
}

func TestMissionValidRoundTrip(t *testing.T) {
	m := validMission()
	if err := Validate(m); err != nil {
		t.Fatalf("valid mission rejected: %v", err)
	}
	raw, err := Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got, err := Unmarshal(raw)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != m.ID || got.Objective != m.Objective {
		t.Fatalf("round trip lost identity: %+v", got)
	}
	if err := got.AcceptancePreserved(); err != nil {
		t.Fatalf("acceptance not preserved: %v", err)
	}
	if got.Acceptance[0] != m.Acceptance[0] {
		t.Fatalf("acceptance text changed")
	}
}

func TestMissionInvalid(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Mission)
	}{
		{"empty id", func(m *Mission) { m.ID = "" }},
		{"empty objective", func(m *Mission) { m.Objective = " " }},
		{"bad bar", func(m *Mission) { m.QualityBar = "LIGHT" }},
		{"empty acceptance", func(m *Mission) { m.Acceptance = nil; m.Seal() }},
		{"retry too high", func(m *Mission) { m.Budget.RetryBudget = 99 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := validMission()
			tc.mut(m)
			if err := Validate(m); err == nil {
				t.Fatalf("expected invalid mission")
			}
		})
	}
}

func TestMissionRejectsGiantObjective(t *testing.T) {
	m := validMission()
	m.Objective = strings.Repeat("x", MaxObjectiveBytes+1)
	if err := Validate(m); err == nil {
		t.Fatalf("giant objective must be rejected")
	}
}

func TestMissionRejectsGiantConstraintDump(t *testing.T) {
	m := validMission()
	m.Constraints = []string{strings.Repeat("AGENTS.md dump ", MaxObjectiveBytes/10)}
	if err := Validate(m); err == nil {
		t.Fatalf("giant constraints must be rejected")
	}
}

func TestMissionAllowsProviderNameInObjective(t *testing.T) {
	m := validMission()
	m.Objective = "Fix timeout; do not send this job to claude"
	m.Seal()
	if err := Validate(m); err != nil {
		t.Fatalf("user objective text should not fail structural provider scan: %v", err)
	}
}

func TestMissionAcceptanceImmutable(t *testing.T) {
	m := validMission()
	m.Acceptance = []string{"quietly rewrite the definition of done"}
	if err := m.AcceptancePreserved(); err == nil {
		t.Fatalf("mutated acceptance must fail preservation check")
	}
	if _, err := Marshal(m); err == nil {
		t.Fatalf("marshal must refuse mutated acceptance")
	}
}

func TestMissionJSONHasNoProviderLeak(t *testing.T) {
	m := validMission()
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if leak := ProviderLeak(m.structuralJSON()); leak != "" {
		t.Fatalf("structural leak %q in %s", leak, raw)
	}
}

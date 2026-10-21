package mission

import (
	"strings"
)

var providerTokens = []string{
	"claude",
	"gemini",
	"grok",
	"opendesign",
	"antigravity",
	"cursor-grok",
	"anthropic",
	"openai",
}

func ProviderLeak(raw []byte) string {
	low := strings.ToLower(string(raw))
	for _, tok := range providerTokens {
		if strings.Contains(low, tok) {
			return tok
		}
	}
	return ""
}

type Seed struct {
	ID               string
	Objective        string
	WorkspaceRoot    string
	DeclaredType     string
	QualityBar       string
	SkipVisualGate   bool
	MaxIterations    int
	RequiresVisual   bool
	RequiresSecurity bool
	Tags             []string
	Acceptance       []string
	Surfaces         []string
	Capabilities     []string
	Constraints      []string
}

func FromSeed(seed Seed) *Mission {
	m := &Mission{
		SchemaVersion:    SchemaVersion,
		ID:               seed.ID,
		Objective:        strings.TrimSpace(seed.Objective),
		QualityBar:       seed.QualityBar,
		RepoPath:         seed.WorkspaceRoot,
		Constraints:      append([]string{}, seed.Constraints...),
		Acceptance:       append([]string{}, seed.Acceptance...),
		AffectedSurfaces: append([]string{}, seed.Surfaces...),
		Capabilities:     append([]string{}, seed.Capabilities...),
		Autonomy:         AutonomyAuto,
		Run: RunMetadata{
			WorkspaceRoot: seed.WorkspaceRoot,
			Attempt:       1,
			CurrentWorker: string(RoleBuilder),
		},
		RequiresSecurity: seed.RequiresSecurity,
	}
	if m.ID == "" {
		m.ID = "mission"
	}
	if m.QualityBar == "" {
		m.QualityBar = BarStandard
	}
	m.Type = mapDeclaredType(seed.DeclaredType, seed.RequiresVisual, seed.Tags, seed.Objective)
	m.VisualImportance = inferVisual(m.Type, seed.RequiresVisual, seed.Objective)
	m.Risk = inferRisk(m.Type, seed.RequiresSecurity, m.QualityBar)
	m.Complexity = inferComplexity(m.Type, m.QualityBar, seed.Objective)
	m.Novelty = "low"
	m.Confidence = "medium"
	m.Intensity = inferIntensity(DeriveInput{
		Type:             m.Type,
		QualityBar:       m.QualityBar,
		Risk:             m.Risk,
		VisualImportance: m.VisualImportance,
		Complexity:       m.Complexity,
	})
	if len(m.Acceptance) == 0 {
		m.Acceptance = defaultAcceptance(m)
	}
	m.RequiredEvidence = defaultEvidence(m)
	in := DeriveInput{
		Type:             m.Type,
		QualityBar:       m.QualityBar,
		Intensity:        m.Intensity,
		Risk:             m.Risk,
		VisualImportance: m.VisualImportance,
		Novelty:          m.Novelty,
		Confidence:       m.Confidence,
		Complexity:       m.Complexity,
		RequiresSecurity: seed.RequiresSecurity,
		HumanRequired:    false,
		AcceptanceEmpty:  len(m.Acceptance) == 0,
	}
	b, why := DeriveBudget(in)
	if seed.MaxIterations > 0 && seed.MaxIterations < b.RetryBudget {
		b.RetryBudget = seed.MaxIterations
	}
	m.Budget = b
	m.BudgetDerivation = why
	m.Seal()
	return m
}

type RefineInput struct {
	QualityBar        string
	Capabilities      []string
	RequiresVisual    bool
	RequiresSecurity  bool
	RequiresHumanGate bool
	GateReason        string
	ResearchDepth     string
	Surfaces          []string
}

func Refine(m *Mission, in RefineInput) {
	if m == nil {
		return
	}
	sealed := CloneAcceptance(m.Acceptance)
	fp := m.AcceptanceFingerprint
	if in.QualityBar != "" {
		m.QualityBar = in.QualityBar
	}
	if len(in.Capabilities) > 0 {
		m.Capabilities = uniqueTrim(in.Capabilities)
	}
	if len(in.Surfaces) > 0 {
		m.AffectedSurfaces = uniqueTrim(append(m.AffectedSurfaces, in.Surfaces...))
	}
	m.RequiresSecurity = in.RequiresSecurity
	if m.Type == TypeDebug && in.RequiresVisual {
		m.Type = TypeProductUI
	}
	if in.RequiresVisual {
		inferred := inferVisual(m.Type, true, m.Objective)
		if m.VisualImportance == VisualNone || inferred == VisualFlagship {
			m.VisualImportance = inferred
		}
	}
	if in.RequiresHumanGate {
		m.Autonomy = AutonomyAuto
	}
	m.Intensity = inferIntensity(DeriveInput{
		Type:             m.Type,
		QualityBar:       m.QualityBar,
		Risk:             m.Risk,
		VisualImportance: m.VisualImportance,
		Complexity:       m.Complexity,
		RequiresSecurity: m.RequiresSecurity,
	})
	if in.ResearchDepth != "" {
		m.Novelty = mapResearchDepth(in.ResearchDepth)
	}
	b, why := DeriveBudget(DeriveInput{
		Type:             m.Type,
		QualityBar:       m.QualityBar,
		Intensity:        m.Intensity,
		Risk:             m.Risk,
		VisualImportance: m.VisualImportance,
		Novelty:          m.Novelty,
		Confidence:       m.Confidence,
		Complexity:       m.Complexity,
		RequiresSecurity: m.RequiresSecurity,
	})
	m.Budget = b
	m.BudgetDerivation = why
	m.Acceptance = sealed
	m.AcceptanceFingerprint = fp
	m.Seal()
	m.Acceptance = sealed
	m.AcceptanceFingerprint = fp
}

func mapDeclaredType(declared string, visual bool, tags []string, objective string) string {
	d := strings.ToLower(strings.TrimSpace(declared))
	hay := strings.ToLower(objective + " " + strings.Join(tags, " "))
	switch d {
	case "bugfix", "bug":
		return TypeDebug
	case "refactor":
		return TypeAPI
	case "research", "docs", "documentation":
		return TypeDocs
	case "security":
		return TypeAPI
	case "design":
		return TypeProductUI
	case "feature":
		if visual || containsAny(hay, []string{"ui", "landing", "portfolio", "page", "dashboard", "settings", "form", "screen", "responsive"}) {
			return TypeProductUI
		}
		if containsAny(hay, []string{"crud", "form", "table"}) {
			return TypeCRUD
		}
		if containsAny(hay, []string{"api", "endpoint", "backend"}) {
			return TypeAPI
		}
		return TypeCRUD
	}
	if visual || containsAny(hay, []string{"landing", "portfolio", "visual", "ui"}) {
		return TypeProductUI
	}
	if containsAny(hay, []string{"nil pointer", "fix", "bug", "crash"}) {
		return TypeDebug
	}
	if containsAny(hay, []string{"api", "endpoint"}) {
		return TypeAPI
	}
	return TypeCRUD
}

func inferVisual(typ string, requiresVisual bool, objective string) string {
	if typ != TypeProductUI && !requiresVisual {
		return VisualNone
	}
	hay := strings.ToLower(objective)
	if containsAny(hay, []string{
		"portfolio", "award", "3d", "flagship", "showcase", "art-direct", "art directed",
		"brand experience", "hiring site", "unusual interactive", "visually distinctive",
		"webgl", "shader",
	}) {
		return VisualFlagship
	}
	if requiresVisual || typ == TypeProductUI {
		return VisualProduct
	}
	return VisualNone
}

func inferRisk(typ string, security bool, bar string) string {
	if security || bar == BarExperimental {
		return RiskHigh
	}
	if typ == TypeProductUI || bar == BarPremium {
		return RiskMedium
	}
	return RiskLow
}

func inferComplexity(typ, bar, objective string) string {
	hay := strings.ToLower(objective)
	if bar == BarExperimental || containsAny(hay, []string{"3d", "webgl", "shader"}) {
		return "high"
	}
	if typ == TypeProductUI || bar == BarPremium {
		return "medium"
	}
	return "low"
}

func mapResearchDepth(depth string) string {
	switch strings.ToUpper(strings.TrimSpace(depth)) {
	case "DEEP":
		return "high"
	case "LIGHT", "TARGETED":
		return "medium"
	default:
		return "low"
	}
}

func defaultAcceptance(m *Mission) []string {
	switch m.Type {
	case TypeDebug:
		return []string{"The reported defect is fixed and existing tests still pass."}
	case TypeAPI:
		return []string{"The API behavior matches the stated contract and mechanical tests pass."}
	case TypeCRUD:
		return []string{"The requested behavior works on the named surfaces without unrelated refactors."}
	case TypeDocs:
		return []string{"The documentation answers the stated question and names its sources."}
	default:
		return []string{"The implemented surfaces match the stated objective and required evidence is produced."}
	}
}

func defaultEvidence(m *Mission) []string {
	switch m.Type {
	case TypeProductUI:
		return []string{"named files changed", "running app or screenshot paths when verification runs"}
	case TypeDebug, TypeAPI, TypeCRUD:
		return []string{"named files changed", "test or build output"}
	default:
		return []string{"named files changed"}
	}
}

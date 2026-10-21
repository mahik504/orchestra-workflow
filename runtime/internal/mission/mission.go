package mission

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	SchemaVersion     = 1
	MaxRetryBudget    = 4
	MaxObjectiveBytes = 16384
)

const (
	IntensityLight Intensity = "LIGHT"
	IntensityFull  Intensity = "FULL"
)

const (
	BarStandard     = "STANDARD"
	BarPremium      = "PREMIUM"
	BarExperimental = "EXPERIMENTAL"
)

const (
	TypeProductUI = "product_ui"
	TypeCRUD      = "crud"
	TypeDebug     = "debug"
	TypeAPI       = "api"
	TypeDocs      = "docs"
)

const (
	RiskLow    = "low"
	RiskMedium = "medium"
	RiskHigh   = "high"
)

const (
	VisualNone     = "none"
	VisualProduct  = "product"
	VisualFlagship = "flagship"
)

const (
	AutonomyAuto          = "auto"
	AutonomyHumanRequired = "human_required"
)

const (
	StatusDraft        MissionStatus = "draft"
	StatusReady        MissionStatus = "ready"
	StatusRunning      MissionStatus = "running"
	StatusWaitingHuman MissionStatus = "waiting_human"
	StatusRetrying     MissionStatus = "retrying"
	StatusSucceeded    MissionStatus = "succeeded"
	StatusFailed       MissionStatus = "failed"
	StatusExhausted    MissionStatus = "exhausted"
)

type Intensity string
type MissionStatus string

type Mission struct {
	SchemaVersion         int           `json:"schema_version"`
	ID                    string        `json:"id"`
	Objective             string        `json:"objective"`
	QualityBar            string        `json:"quality_bar"`
	Intensity             Intensity     `json:"intensity"`
	Type                  string        `json:"type"`
	Constraints           []string      `json:"constraints,omitempty"`
	Acceptance            []string      `json:"acceptance"`
	AcceptanceFingerprint string        `json:"acceptance_fingerprint"`
	AffectedSurfaces      []string      `json:"affected_surfaces,omitempty"`
	Capabilities          []string      `json:"capabilities,omitempty"`
	RequiredEvidence      []string      `json:"required_evidence,omitempty"`
	HumanNeeded           []string      `json:"human_needed,omitempty"`
	Budget                Budget        `json:"budget"`
	Status                MissionStatus `json:"status"`
	RepoPath              string        `json:"repo_path,omitempty"`
	Risk                  string        `json:"risk,omitempty"`
	VisualImportance      string        `json:"visual_importance,omitempty"`
	Autonomy              string        `json:"autonomy,omitempty"`
	Novelty               string        `json:"novelty,omitempty"`
	Confidence            string        `json:"confidence,omitempty"`
	Complexity            string        `json:"complexity,omitempty"`
	RequiresSecurity      bool          `json:"requires_security,omitempty"`
	Run                   RunMetadata   `json:"run"`
	BudgetDerivation      string        `json:"budget_derivation,omitempty"`
}

type RunMetadata struct {
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	WorkspaceRoot string `json:"workspace_root,omitempty"`
	CurrentStage  string `json:"current_stage,omitempty"`
	CurrentWorker string `json:"current_worker,omitempty"`
	Attempt       int    `json:"attempt"`
}

func DefaultHumanNeeded() []string {
	return []string{
		"ambiguous requirement affecting product behavior",
		"missing critical credential or authorization",
		"destructive action requiring approval",
		"specification or contract conflict",
		"acceptance criteria cannot be interpreted safely",
		"repeated bounded failure requiring human direction",
		"identity or brand decision",
		"creative direction lock (stills, named-ref pick)",
		"flagship visual direction selection",
		"security-sensitive change",
		"irreversible deploy",
		"material model, tool, or host decision missing",
	}
}

func FingerprintAcceptance(acceptance []string) string {
	joined := strings.Join(acceptance, "\n")
	sum := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(sum[:])
}

func (m *Mission) Seal() {
	if m.SchemaVersion == 0 {
		m.SchemaVersion = SchemaVersion
	}
	if m.HumanNeeded == nil {
		m.HumanNeeded = DefaultHumanNeeded()
	}
	if m.Acceptance == nil {
		m.Acceptance = []string{}
	}
	m.AcceptanceFingerprint = FingerprintAcceptance(m.Acceptance)
	now := time.Now().UTC().Format(time.RFC3339)
	if m.Run.CreatedAt == "" {
		m.Run.CreatedAt = now
	}
	m.Run.UpdatedAt = now
	if m.Status == "" {
		m.Status = StatusReady
	}
	if m.Run.Attempt <= 0 {
		m.Run.Attempt = 1
	}
}

func (m *Mission) AcceptancePreserved() error {
	if m == nil {
		return fmt.Errorf("mission is nil")
	}
	got := FingerprintAcceptance(m.Acceptance)
	if m.AcceptanceFingerprint == "" {
		return fmt.Errorf("mission has no sealed acceptance fingerprint")
	}
	if got != m.AcceptanceFingerprint {
		return fmt.Errorf("acceptance criteria were mutated after seal")
	}
	return nil
}

func Validate(m *Mission) error {
	if m == nil {
		return fmt.Errorf("mission is nil")
	}
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("mission id is required")
	}
	if strings.ContainsAny(m.ID, `/\`) {
		return fmt.Errorf("mission id must not contain path separators")
	}
	obj := strings.TrimSpace(m.Objective)
	if obj == "" {
		return fmt.Errorf("mission objective is required")
	}
	if len(m.Objective) > MaxObjectiveBytes {
		return fmt.Errorf("mission objective exceeds %d bytes (refusing whole-project dumps)", MaxObjectiveBytes)
	}
	if dumpSize(m.Constraints) > MaxObjectiveBytes || dumpSize(m.RequiredEvidence) > MaxObjectiveBytes || dumpSize(m.HumanNeeded) > MaxObjectiveBytes {
		return fmt.Errorf("mission field exceeds %d bytes (refusing whole-project dumps)", MaxObjectiveBytes)
	}
	switch m.QualityBar {
	case BarStandard, BarPremium, BarExperimental:
	default:
		return fmt.Errorf("quality bar must be STANDARD, PREMIUM, or EXPERIMENTAL")
	}
	switch m.Intensity {
	case IntensityLight, IntensityFull:
	default:
		return fmt.Errorf("intensity must be LIGHT or FULL")
	}
	switch m.Type {
	case TypeProductUI, TypeCRUD, TypeDebug, TypeAPI, TypeDocs, "":
	default:
		return fmt.Errorf("unknown mission type %q", m.Type)
	}
	if len(m.Acceptance) == 0 {
		return fmt.Errorf("at least one acceptance criterion is required")
	}
	for _, line := range m.Acceptance {
		if strings.TrimSpace(line) == "" {
			return fmt.Errorf("acceptance criteria must not contain empty entries")
		}
	}
	if err := validateBudget(m.Budget); err != nil {
		return err
	}
	if m.Autonomy == AutonomyHumanRequired && m.Budget.RetryBudget != 0 {
		return fmt.Errorf("human-required missions must have retry_budget 0")
	}
	if leak := ProviderLeak(m.structuralJSON()); leak != "" {
		return fmt.Errorf("provider-specific token %q leaked into mission contract", leak)
	}
	return nil
}

func (m *Mission) structuralJSON() []byte {
	if m == nil {
		return nil
	}
	snap := *m
	snap.Objective = ""
	snap.Constraints = nil
	snap.Acceptance = nil
	snap.RequiredEvidence = nil
	snap.HumanNeeded = nil
	b, _ := json.Marshal(snap)
	return b
}

func Marshal(m *Mission) ([]byte, error) {
	if err := Validate(m); err != nil {
		return nil, err
	}
	if err := m.AcceptancePreserved(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(m, "", "  ")
}

func Unmarshal(raw []byte) (*Mission, error) {
	var m Mission
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("invalid mission json: %w", err)
	}
	if err := Validate(&m); err != nil {
		return nil, err
	}
	if err := m.AcceptancePreserved(); err != nil {
		return nil, err
	}
	return &m, nil
}

func CloneAcceptance(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

func dumpSize(ss []string) int {
	n := 0
	for _, s := range ss {
		n += len(s)
	}
	return n
}

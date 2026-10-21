package memory

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

const (
	ClassPreference           = "preference"
	ClassArchitectureDecision = "architecture_decision"
	ClassProjectFact          = "project_fact"
	ClassWorkflow             = "workflow"
	ClassFailureLesson        = "failure_lesson"
	ClassResourceKnowledge    = "resource_knowledge"
	ClassEvaluationHistory    = "evaluation_history"

	StatusCurrent         = "current"
	StatusHistorical      = "historical"
	StatusStale           = "stale"
	StatusUnverified      = "unverified"
	StatusLegacyAggregate = "legacy_aggregate_without_payload"
	StatusSuperseded      = "superseded"

	MaxDurableFacts = 8

	WriteUserPreference       = "user_preference"
	WriteArchitectureDecision = "architecture_decision"
	WriteValidatedWorkflow    = "validated_workflow"
	WriteFailureLesson        = "failure_lesson"
	WriteVerifiedResource     = "verified_resource"
	WriteProjectFact          = "project_fact"
	WriteChat                 = "chat"
	WriteToolCall             = "tool_call"
	WriteRetry                = "retry"
	WriteTransientError       = "transient_error"
	WriteScreenshot           = "screenshot"
	WriteModelAvailability    = "model_availability"
	WriteSuccessfulCommand    = "successful_command"

	PrivateFixtureMark = "PRIVATE_BRAIN_FIXTURE"
)

// Fact is a durable Brain note. It is not runtime state.
type Fact struct {
	ID           string `json:"id,omitempty"`
	Class        string `json:"class,omitempty"`
	Text         string `json:"text,omitempty"`
	Source       string `json:"source,omitempty"`
	AsOf         string `json:"as_of,omitempty"`
	Status       string `json:"status,omitempty"`
	Scope        string `json:"scope,omitempty"`
	LastVerified string `json:"last_verified,omitempty"`
	SupersededBy string `json:"superseded_by,omitempty"`
	PublicSafe   bool   `json:"public_safe,omitempty"`
}

type SelectInput struct {
	Candidates     []Fact
	Scope          string
	IncludeHistory bool
	Live           map[string]string
}

type SelectResult struct {
	Selected      []Fact `json:"selected,omitempty"`
	Rejected      int    `json:"rejected"`
	StaleExcluded int    `json:"stale_excluded"`
	Bytes         int    `json:"bytes"`
	MaxFacts      int    `json:"max_facts"`
}

func AdmitDurable(reason string) bool {
	switch strings.TrimSpace(reason) {
	case WriteUserPreference, WriteArchitectureDecision, WriteValidatedWorkflow,
		WriteFailureLesson, WriteVerifiedResource, WriteProjectFact:
		return true
	default:
		return false
	}
}

// AuthoritativeRuntime returns the live value. A Brain string never replaces it.
func AuthoritativeRuntime(live, brain string) string {
	if strings.TrimSpace(live) != "" {
		return live
	}
	return ""
}

// SummarizeRun keeps a lesson and a path reference. It refuses runtime dumps.
func SummarizeRun(missionID, lesson, pathRef string) (Fact, error) {
	blob := lesson + "\n" + pathRef
	for _, banned := range []string{"CAPSULE_DUMP", "RETRY_LOG", "PNG_BYTES", "RUN_DIR_LISTING"} {
		if strings.Contains(blob, banned) {
			return Fact{}, fmt.Errorf("brain summary cannot carry runtime state")
		}
	}
	if strings.Count(pathRef, "\n") > 1 {
		return Fact{}, fmt.Errorf("brain summary cannot carry a run directory listing")
	}
	return Fact{
		ID:     missionID,
		Class:  ClassFailureLesson,
		Text:   strings.TrimSpace(lesson),
		Source: "runtime-summary",
		Status: StatusCurrent,
		Scope:  pathRef,
	}, nil
}

func RejectPrivateLeak(text string, publicSafe bool) error {
	if strings.Contains(text, PrivateFixtureMark) && !publicSafe {
		return fmt.Errorf("private brain content cannot enter a public artifact")
	}
	return nil
}

func RepoCopy(text string) bool {
	for _, marker := range []string{"AGENTS.md FULL COPY", "ARCHITECTURE.md FULL COPY", "REGISTRY_FULL_COPY"} {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}

// SelectDurable returns at most MaxDurableFacts current facts for a scope.
// Stale, superseded, and legacy aggregates are excluded unless IncludeHistory is set.
// A fact whose id is already present in Live is rejected; runtime wins.
func SelectDurable(in SelectInput) SelectResult {
	out := SelectResult{MaxFacts: MaxDurableFacts}
	for _, f := range in.Candidates {
		if in.Scope != "" && f.Scope != "" && f.Scope != in.Scope {
			out.Rejected++
			continue
		}
		if !in.IncludeHistory && (f.Status == StatusStale || f.Status == StatusSuperseded || f.Status == StatusLegacyAggregate) {
			out.StaleExcluded++
			out.Rejected++
			continue
		}
		if in.Live != nil {
			if _, ok := in.Live[f.ID]; ok {
				out.Rejected++
				continue
			}
		}
		if len(out.Selected) >= MaxDurableFacts {
			out.Rejected++
			continue
		}
		out.Selected = append(out.Selected, f)
	}
	raw, _ := json.Marshal(out.Selected)
	out.Bytes = len(raw)
	return out
}

// ConflictWinner names the authority that beats a Brain fact.
// Order: live runtime, current docs, verified registry, brain history, unverified.
func ConflictWinner(id string, live, docs map[string]string, registryVerified map[string]bool, fact Fact) string {
	if live != nil && strings.TrimSpace(live[id]) != "" {
		return "runtime"
	}
	if docs != nil && strings.TrimSpace(docs[id]) != "" {
		return "docs"
	}
	if registryVerified != nil && registryVerified[id] {
		return "registry"
	}
	switch fact.Status {
	case StatusCurrent, StatusHistorical:
		return "brain"
	case StatusSuperseded:
		return "superseded"
	default:
		return "unverified"
	}
}

type RegistryRow struct {
	ID           string `json:"id"`
	LastVerified string `json:"last_verified"`
}

func SplitVerification(rows []RegistryRow) (verified, unverified int) {
	for _, row := range rows {
		if strings.TrimSpace(row.LastVerified) != "" {
			verified++
		} else {
			unverified++
		}
	}
	return verified, unverified
}

func ResourceVerified(lastVerified string) bool {
	return strings.TrimSpace(lastVerified) != ""
}

// LoadDurable reads an already-fetched Brain payload. Unavailable, missing,
// empty, and malformed payloads yield no facts and no error. Brain cannot
// stop the caller.
func LoadDurable(raw []byte, mode, scope string, live map[string]string) SelectResult {
	empty := SelectResult{MaxFacts: MaxDurableFacts}
	switch strings.TrimSpace(mode) {
	case "unavailable", "missing", "empty":
		return empty
	}
	if len(bytesTrim(raw)) == 0 {
		return empty
	}
	var facts []Fact
	if err := json.Unmarshal(raw, &facts); err != nil {
		return empty
	}
	return SelectDurable(SelectInput{Candidates: facts, Scope: scope, Live: live})
}

func bytesTrim(raw []byte) []byte {
	return []byte(strings.TrimSpace(string(raw)))
}

// DurableDecisions is the Brain → mission-scoped context step.
// Catalogs, conversation, dumps, and provider labels never become decisions.
func DurableDecisions(in SelectInput) (decisions []string, result SelectResult) {
	var clean []Fact
	dropped := 0
	for _, f := range in.Candidates {
		if f.Source == WriteModelAvailability || !capsuleSafe(f.Text) {
			dropped++
			continue
		}
		clean = append(clean, f)
	}
	in.Candidates = clean
	result = SelectDurable(in)
	result.Rejected += dropped
	for _, f := range result.Selected {
		decisions = append(decisions, f.Text)
	}
	return decisions, result
}

// BindCapsule compiles a worker capsule from bounded Brain decisions.
// mission.Compile is unchanged. A Brain miss compiles with no decisions.
func BindCapsule(m *mission.Mission, in SelectInput, extra mission.CompileRequest) (*mission.Capsule, *mission.ScopeReport, SelectResult, error) {
	decisions, result := DurableDecisions(in)
	extra.Mission = m
	if extra.Role == "" {
		extra.Role = mission.RoleBuilder
	}
	extra.Decisions = append(append([]string{}, extra.Decisions...), decisions...)
	c, rep, err := mission.Compile(extra)
	return c, rep, result, err
}

func capsuleSafe(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}
	if mission.ProviderLeak([]byte(text)) != "" {
		return false
	}
	for _, mark := range []string{
		"CONVERSATION_HISTORY", "MCP_CATALOG", "SKILL_CATALOG", "RESOURCE_CATALOG",
		"FULL_BRAIN", "RUN_DUMP", "CAPSULE_DUMP", "RETRY_LOG", "PNG_BYTES",
		PrivateFixtureMark, "FULL_PRIVATE_DOCUMENT",
	} {
		if strings.Contains(text, mark) {
			return false
		}
	}
	return true
}

// Retrieve copies a fact. It does not change provenance.
func Retrieve(f Fact) Fact {
	return f
}

// VerifyFact is the only status upgrade. An empty check date is refused.
func VerifyFact(f Fact, checkedOn string) (Fact, error) {
	if strings.TrimSpace(checkedOn) == "" {
		return f, fmt.Errorf("verification requires an explicit check date")
	}
	f.Status = StatusCurrent
	f.LastVerified = checkedOn
	return f, nil
}

// SerializePublic is the public-artifact gate. Private text, private paths,
// full private documents, and runtime dumps stay out. A repository path
// that is not a copy is allowed.
func SerializePublic(text string, publicSafe bool) (string, error) {
	if err := RejectPrivateLeak(text, publicSafe); err != nil {
		return "", err
	}
	if RepoCopy(text) {
		return "", fmt.Errorf("public artifact cannot copy a repository document")
	}
	low := strings.ToLower(strings.ReplaceAll(text, "\\", "/"))
	privateVault := "orchestra" + "-brain"
	if !publicSafe && (strings.Contains(low, privateVault+"/") || strings.Contains(low, privateVault+"\\")) {
		return "", fmt.Errorf("private brain path cannot enter a public artifact")
	}
	for _, mark := range []string{"FULL_PRIVATE_DOCUMENT", "RUN_DUMP", "CAPSULE_DUMP", "RETRY_LOG", "PNG_BYTES"} {
		if strings.Contains(text, mark) {
			return "", fmt.Errorf("runtime or private dump cannot enter a public artifact")
		}
	}
	return text, nil
}

// SurfaceClass is a duplicate-surface label. Actions are KEEP or MIGRATE only.
type SurfaceClass struct {
	Name   string
	Class  string
	Action string
	Note   string
}

func DuplicateSurfaces() []SurfaceClass {
	return []SurfaceClass{
		{"PlaywrightVerifier", "CANONICAL", "KEEP", "Layer A harness owner"},
		{"catalog playwright", "ADAPTER", "KEEP", "library knowledge"},
		{"host-stack Playwright MCP", "ADAPTER", "KEEP", "not Layer A"},
		{"playwright-cli", "ADAPTER", "KEEP", "catalog row"},
		{"browser-use", "OPTIONAL", "KEEP", "do not promote"},
		{"agent-browser", "OPTIONAL", "KEEP", "do not promote"},
		{"Context7 MCP", "CANONICAL", "KEEP", "library-doc lookup"},
		{"Context7 skill-list label", "DUPLICATE", "MIGRATE", "point the note at the one MCP"},
		{"gsap-core skill", "CANONICAL", "KEEP", "motion implementation owner"},
		{"catalog gsap", "ADAPTER", "KEEP", "library knowledge"},
		{"motion-dev", "OPTIONAL", "KEEP", "reference"},
		{"motion-primitives", "OPTIONAL", "KEEP", "reference"},
		{"motionsites-ai", "OPTIONAL", "KEEP", "reference"},
		{"shadcn MCP", "OPTIONAL", "KEEP", "not core, not repaired"},
		{"catalog shadcn-ui", "OPTIONAL", "KEEP", "not a second shadcn server"},
		{"GitHub host-stack plugin", "OPTIONAL", "KEEP", "no catalog row"},
	}
}

func Stances() map[string]string {
	return map[string]string{
		"shadcn":     "OPTIONAL",
		"opendesign": "RETIRED",
		"jcode":      "DEFER",
		"grok-4.7":   "OPAQUE_LABEL",
	}
}

package memory

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

func TestPhase3F1_BrainToCapsule(t *testing.T) {
	var facts []Fact
	for i := 0; i < 10; i++ {
		facts = append(facts, Fact{
			ID: "keep", Class: ClassProjectFact, Scope: "job", Status: StatusCurrent,
			Text: "check the parser contract",
		})
	}
	facts = append(facts,
		Fact{ID: "old", Scope: "job", Status: StatusStale, Text: "STALE_FACT_MARKER"},
		Fact{ID: "other", Scope: "parked", Status: StatusCurrent, Text: "UNRELATED_PROJECT_FACT"},
		Fact{ID: "layer-a", Scope: "job", Status: StatusHistorical, Text: "brain says the owner changed"},
		Fact{ID: "chat", Scope: "job", Status: StatusCurrent, Text: "CONVERSATION_HISTORY user said hello"},
		Fact{ID: "mcp", Scope: "job", Status: StatusCurrent, Text: "MCP_CATALOG every server"},
		Fact{ID: "skills", Scope: "job", Status: StatusCurrent, Text: "SKILL_CATALOG forty four"},
		Fact{ID: "res", Scope: "job", Status: StatusCurrent, Text: "RESOURCE_CATALOG one six eight"},
		Fact{ID: "full", Scope: "job", Status: StatusCurrent, Text: "FULL_BRAIN preferences career"},
		Fact{ID: "run", Scope: "job", Status: StatusCurrent, Text: "RUN_DUMP mission json"},
		Fact{ID: "grok", Scope: "job", Status: StatusCurrent, Source: WriteModelAvailability, Text: "prefer grok-4.7"},
		Fact{ID: "claude", Scope: "job", Status: StatusCurrent, Source: WriteModelAvailability, Text: "prefer claude"},
		Fact{ID: "gemini", Scope: "job", Status: StatusCurrent, Source: WriteModelAvailability, Text: "prefer gemini"},
	)
	m := mission.FromSeed(mission.Seed{
		ID: "job", Objective: "Fix parser crash", DeclaredType: "BUGFIX",
		Acceptance: []string{"parser returns the node"}, QualityBar: mission.BarPremium,
	})
	bar, acc, retry := m.QualityBar, strings.Join(m.Acceptance, "|"), m.Budget.RetryBudget
	whole := bytes.Repeat([]byte("X"), 136000)
	c, rep, sel, err := BindCapsule(m, SelectInput{
		Candidates: facts,
		Scope:      "job",
		Live:       map[string]string{"layer-a": "PlaywrightVerifier"},
	}, mission.CompileRequest{
		ConversationHistory: []string{"CONVERSATION_HISTORY user said hello"},
		MCPCatalog:          []string{"MCP_CATALOG every server"},
		AllSkills:           []string{"allowed-skill"},
		Skills:              []string{"allowed-skill", "not-allowed"},
		WholeProjectContext: whole,
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.QualityBar != bar || strings.Join(m.Acceptance, "|") != acc || m.Budget.RetryBudget != retry {
		t.Fatal("capsule bind changed mission authority")
	}
	if len(sel.Selected) != MaxDurableFacts {
		t.Fatalf("selected facts %d, cap %d", len(sel.Selected), MaxDurableFacts)
	}
	if len(c.Decisions) == 0 || len(c.Decisions) > MaxDurableFacts {
		t.Fatalf("decisions %d", len(c.Decisions))
	}
	if sel.StaleExcluded < 1 {
		t.Fatalf("stale excluded %d", sel.StaleExcluded)
	}
	raw, _ := json.Marshal(c)
	for _, banned := range []string{
		"STALE_FACT_MARKER", "UNRELATED_PROJECT_FACT", "owner changed", "CONVERSATION_HISTORY",
		"MCP_CATALOG", "SKILL_CATALOG", "RESOURCE_CATALOG", "FULL_BRAIN", "RUN_DUMP",
		"grok", "claude", "gemini",
	} {
		if strings.Contains(strings.ToLower(string(raw)), banned) {
			t.Fatalf("capsule contained %s (%d bytes)", banned, len(raw))
		}
	}
	if len(raw) >= len(whole)/10 {
		t.Fatalf("capsule %d bytes is not far under %d", len(raw), len(whole))
	}
	joined := strings.Join(rep.Excluded, ",")
	for _, want := range []string{"conversation_history", "mcp_catalog", "whole_project_context"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("excluded %s missing from %s", want, joined)
		}
	}
	t.Logf("phase3f1 available=%d selected=%d capsule_decisions=%d stale_excluded=%d rejected=%d capsule_bytes=%d whole_project_bytes=%d",
		len(facts), len(sel.Selected), len(c.Decisions), sel.StaleExcluded, sel.Rejected, len(raw), len(whole))
}

func TestPhase3F1_BrainFailureDoesNotBlock(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "bf", Objective: "Fix parser crash", DeclaredType: "BUGFIX",
		Acceptance: []string{"parser returns the node"}, QualityBar: mission.BarPremium,
	})
	bar, acc, retry := m.QualityBar, strings.Join(m.Acceptance, "|"), m.Budget.RetryBudget
	staleRaw, _ := json.Marshal([]Fact{{ID: "old", Scope: "job", Status: StatusStale, Text: "STALE_FACT_MARKER"}})
	modes := []struct {
		name string
		raw  []byte
	}{
		{"unavailable", nil},
		{"empty", []byte("[]")},
		{"malformed", []byte("{")},
		{"stale", staleRaw},
		{"missing", nil},
	}
	for _, mode := range modes {
		got := LoadDurable(mode.raw, mode.name, "job", nil)
		if len(got.Selected) != 0 {
			t.Fatalf("%s selected %d", mode.name, len(got.Selected))
		}
		c, _, _, err := BindCapsule(m, SelectInput{Scope: "job"}, mission.CompileRequest{})
		if err != nil {
			t.Fatalf("%s compile: %v", mode.name, err)
		}
		if c.Role != mission.RoleBuilder {
			t.Fatalf("%s role %s", mode.name, c.Role)
		}
		if m.QualityBar != bar || strings.Join(m.Acceptance, "|") != acc || m.Budget.RetryBudget != retry {
			t.Fatalf("%s mutated mission", mode.name)
		}
	}
}

func TestPhase3F1_PublicSerialize(t *testing.T) {
	if _, err := SerializePublic(PrivateFixtureMark+" career note", false); err == nil {
		t.Fatal("private fact crossed")
	}
	if _, err := SerializePublic("Brain remembers. Runtime executes.", false); err != nil {
		t.Fatal(err)
	}
	if _, err := SerializePublic(PrivateFixtureMark+" marked public", true); err != nil {
		t.Fatal(err)
	}
	ref, err := SerializePublic("ARCHITECTURE.md", false)
	if err != nil || strings.Contains(ref, "FULL COPY") {
		t.Fatal(err)
	}
	if _, err := SerializePublic(`C:\projects\orchestra`+"-brain"+`\Preferences.md`, false); err == nil {
		t.Fatal("private path crossed")
	}
	if _, err := SerializePublic("FULL_PRIVATE_DOCUMENT preferences body", false); err == nil {
		t.Fatal("full private document crossed")
	}
	if _, err := SerializePublic("RUN_DUMP .orchestra/runs/m1/mission.json", false); err == nil {
		t.Fatal("runtime dump crossed")
	}
}

func TestPhase3F1_ProvenanceStaysPut(t *testing.T) {
	facts := []Fact{
		{ID: "u", Status: StatusUnverified, Scope: "job", Text: "unverified note"},
		{ID: "s", Status: StatusStale, Scope: "job", Text: "stale note"},
		{ID: "h", Status: StatusHistorical, Scope: "job", Text: "historical note"},
		{ID: "l", Status: StatusLegacyAggregate, Scope: "job", Text: "legacy note"},
	}
	before := append([]Fact{}, facts...)
	_ = SelectDurable(SelectInput{Candidates: facts, Scope: "job"})
	for i := range facts {
		got := Retrieve(facts[i])
		if got.Status != before[i].Status || facts[i].Status != before[i].Status {
			t.Fatalf("retrieval upgraded %s to %s", before[i].Status, got.Status)
		}
		if _, err := VerifyFact(facts[i], ""); err == nil {
			t.Fatal("blank verification upgraded status")
		}
		if facts[i].Status != before[i].Status {
			t.Fatal("refused verification still wrote status")
		}
	}
	up, err := VerifyFact(facts[0], "2026-09-22")
	if err != nil {
		t.Fatal(err)
	}
	if up.Status != StatusCurrent || up.LastVerified != "2026-09-22" {
		t.Fatalf("explicit verify %+v", up)
	}
	if facts[0].Status != StatusUnverified {
		t.Fatal("verify mutated the stored fact")
	}
}

func TestPhase3F1_LegacyCannotBecomeCurrent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "resource-memory.json")
	doc := ResourceMemoryDocument{
		SchemaVersion:     "1.0.0",
		TotalEvaluations:  2648,
		AggregateSum:      2648,
		LedgerStatus:      StatusLegacyAggregate,
		EvaluationPayload: "absent",
		Resources: map[string]*ResourceAggregate{
			"gsap": {ResourceID: "gsap", TotalEvaluations: 2648, SuccessCount: 2648},
		},
		Evaluations: []ResourceEvaluation{},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}
	store, err := NewResourceMemoryStore(path)
	if err != nil {
		t.Fatal(err)
	}
	total, _, _, _ := store.GetSummary()
	if total != 0 || store.CurrentEvidenceCount() != 0 {
		t.Fatalf("summary treated 2648 as current: total=%d evidence=%d", total, store.CurrentEvidenceCount())
	}
	sum, status := store.LegacyAggregate()
	if sum != 2648 || status != StatusLegacyAggregate {
		t.Fatalf("legacy %d %s", sum, status)
	}
	store.doc.TotalEvaluations = 2648
	total, _, _, _ = store.GetSummary()
	if total != 0 {
		t.Fatalf("rewritten document total became current evidence: %d", total)
	}
}

func TestPhase3F1_RegistryCountsStayPut(t *testing.T) {
	catalog := filepath.Join("..", "..", "..", "registries", "resources.json")
	raw, err := os.ReadFile(catalog)
	if err != nil {
		t.Fatal(err)
	}
	var rows []RegistryRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	verified, unverified := SplitVerification(rows)
	if len(rows) != 208 || verified != 54 || unverified != 154 {
		t.Fatalf("rows=%d verified=%d unverified=%d", len(rows), verified, unverified)
	}
	if verified == unverified {
		t.Fatal("unverified counted as verified")
	}
	dated := 0
	checkedOn := map[string]string{
		"watermelon-ui":                 "2026-11-05",
		"archify":                       "2026-11-05",
		"understand-anything":           "2026-11-05",
		"markitdown":                    "2026-11-05",
		"data-engineer-handbook":        "2026-11-05",
		"storybook":                     "2026-11-05",
		"chromatic":                     "2026-11-05",
		"grilling":                      "2026-11-07",
		"domain-modeling":               "2026-11-07",
		"to-tickets":                    "2026-11-07",
		"constraint-driven-development": "2026-11-07",
		"receiving-code-review":         "2026-11-07",
	}
	for _, row := range rows {
		if ResourceVerified(row.LastVerified) {
			dated++
			want := "2026-09-06"
			if day, ok := checkedOn[row.ID]; ok {
				want = day
			}
			if row.LastVerified != want {
				t.Fatalf("%s date %s", row.ID, row.LastVerified)
			}
		}
	}
	if dated != verified {
		t.Fatalf("dated %d verified %d", dated, verified)
	}
}

func TestPhase3F1_ModelLabelIsNotArchitecture(t *testing.T) {
	labels := []string{"grok-4.7", "claude", "gemini", "gpt-x"}
	var got []string
	for _, label := range labels {
		decisions, _ := DurableDecisions(SelectInput{Candidates: []Fact{
			{ID: "job", Scope: "job", Status: StatusCurrent, Text: "check the parser contract"},
			{ID: "model", Scope: "job", Status: StatusCurrent, Source: WriteModelAvailability, Text: "prefer " + label},
		}, Scope: "job"})
		got = append(got, strings.Join(decisions, "|"))
	}
	for i := 1; i < len(got); i++ {
		if got[i] != got[0] {
			t.Fatalf("label changed decisions: %q vs %q", got[0], got[i])
		}
	}
	if strings.Contains(got[0], "grok") || strings.Contains(got[0], "claude") || strings.Contains(got[0], "gemini") || strings.Contains(got[0], "gpt") {
		t.Fatalf("vendor label entered decisions: %s", got[0])
	}
	if Stances()["grok-4.7"] != "OPAQUE_LABEL" {
		t.Fatal("grok stance changed")
	}
}

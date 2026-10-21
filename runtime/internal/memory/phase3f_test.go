package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase3F_AdmitDurable(t *testing.T) {
	allow := []string{
		WriteUserPreference, WriteArchitectureDecision, WriteValidatedWorkflow,
		WriteFailureLesson, WriteVerifiedResource, WriteProjectFact,
	}
	deny := []string{
		WriteChat, WriteToolCall, WriteRetry, WriteTransientError,
		WriteScreenshot, WriteModelAvailability, WriteSuccessfulCommand,
	}
	for _, reason := range allow {
		if !AdmitDurable(reason) {
			t.Fatalf("expected admit %s", reason)
		}
	}
	for _, reason := range deny {
		if AdmitDurable(reason) {
			t.Fatalf("expected deny %s", reason)
		}
	}
}

func TestPhase3F_OwnershipRuntimeWins(t *testing.T) {
	live := "browser"
	brain := "none"
	if got := AuthoritativeRuntime(live, brain); got != live {
		t.Fatalf("runtime lost: %s", got)
	}
	if got := AuthoritativeRuntime("", brain); got != "" {
		t.Fatalf("empty live fell through to brain: %q", got)
	}
	if _, err := SummarizeRun("m1", "CAPSULE_DUMP builder json", ".orchestra/runs/m1/mission.json"); err == nil {
		t.Fatal("capsule dump was admitted")
	}
	if _, err := SummarizeRun("m1", "lesson", "a\nb\nc"); err == nil {
		t.Fatal("run listing was admitted")
	}
	fact, err := SummarizeRun("m1", "spec conflict stopped the run", ".orchestra/runs/m1/mission.json")
	if err != nil {
		t.Fatal(err)
	}
	if fact.Class != ClassFailureLesson || strings.Contains(fact.Text, "mission.json") {
		t.Fatalf("summary should be a lesson, got %+v", fact)
	}
}

func TestPhase3F_SelectCapAndStale(t *testing.T) {
	var facts []Fact
	for i := 0; i < 12; i++ {
		facts = append(facts, Fact{ID: "f", Class: ClassProjectFact, Text: "current fact", Status: StatusCurrent, Scope: "flagship"})
	}
	facts = append(facts, Fact{ID: "old", Text: "stale", Status: StatusStale, Scope: "flagship"})
	facts = append(facts, Fact{ID: "gone", Text: "replaced", Status: StatusSuperseded, Scope: "flagship", SupersededBy: "f"})
	facts = append(facts, Fact{ID: "hist", Text: "2648", Status: StatusLegacyAggregate, Scope: "flagship"})
	facts = append(facts, Fact{ID: "other", Text: "elsewhere", Status: StatusCurrent, Scope: "other"})
	got := SelectDurable(SelectInput{Candidates: facts, Scope: "flagship"})
	if len(got.Selected) != MaxDurableFacts {
		t.Fatalf("selected %d, want %d", len(got.Selected), MaxDurableFacts)
	}
	if got.StaleExcluded < 3 {
		t.Fatalf("stale excluded %d", got.StaleExcluded)
	}
	if got.Bytes <= 0 || got.Bytes > 8000 {
		t.Fatalf("selected bytes %d", got.Bytes)
	}
	hist := SelectDurable(SelectInput{Candidates: facts, Scope: "flagship", IncludeHistory: true})
	if hist.StaleExcluded != 0 {
		t.Fatalf("history mode still excluded %d", hist.StaleExcluded)
	}
}

func TestPhase3F_ConflictOrder(t *testing.T) {
	fact := Fact{ID: "layer-a", Text: "playwright mcp is layer a", Status: StatusHistorical}
	live := map[string]string{"layer-a": "PlaywrightVerifier"}
	docs := map[string]string{"layer-a": "PlaywrightVerifier harness"}
	reg := map[string]bool{"layer-a": true}
	if ConflictWinner("layer-a", live, docs, reg, fact) != "runtime" {
		t.Fatal("live runtime did not win")
	}
	if ConflictWinner("layer-a", nil, docs, reg, fact) != "docs" {
		t.Fatal("docs did not beat registry")
	}
	if ConflictWinner("layer-a", nil, nil, reg, fact) != "registry" {
		t.Fatal("registry did not beat brain")
	}
	if ConflictWinner("layer-a", nil, nil, nil, fact) != "brain" {
		t.Fatal("historical brain did not remain when nothing newer exists")
	}
	unknown := Fact{ID: "layer-a", Status: StatusUnverified}
	if ConflictWinner("layer-a", nil, nil, nil, unknown) != "unverified" {
		t.Fatal("unverified counted as knowledge")
	}
	selected := SelectDurable(SelectInput{
		Candidates: []Fact{fact},
		Live:       live,
	})
	if len(selected.Selected) != 0 {
		t.Fatal("conflicting brain fact was retrieved")
	}
}

func TestPhase3F_PublicPrivate(t *testing.T) {
	private := "PRIVATE_BRAIN_FIXTURE career note"
	if err := RejectPrivateLeak(private, false); err == nil {
		t.Fatal("private fixture entered a public artifact")
	}
	if err := RejectPrivateLeak("Brain remembers. Runtime executes.", false); err != nil {
		t.Fatal(err)
	}
	if err := RejectPrivateLeak(private, true); err != nil {
		t.Fatal("explicit public-safe marker was rejected")
	}
	if !RepoCopy("AGENTS.md FULL COPY of the public file") {
		t.Fatal("full repo copy was not detected")
	}
	if RepoCopy("see docs/3.3.2/ORCHESTRA_3_3_2_ARCHITECTURE.md") {
		t.Fatal("a path reference was treated as a copy")
	}
}

func TestPhase3F_ContextBytes(t *testing.T) {
	var facts []Fact
	for i := 0; i < 30; i++ {
		facts = append(facts, Fact{
			ID:     "p",
			Class:  ClassPreference,
			Text:   strings.Repeat("word ", 40),
			Status: StatusCurrent,
			Scope:  "method",
		})
	}
	facts = append(facts, Fact{ID: "s", Text: "old", Status: StatusStale, Scope: "method"})
	got := SelectDurable(SelectInput{Candidates: facts, Scope: "method"})
	const vaultDump = 136000
	if got.Bytes >= vaultDump/10 {
		t.Fatalf("retrieval %d bytes is not far under the %d-byte dump", got.Bytes, vaultDump)
	}
	if len(got.Selected) != MaxDurableFacts {
		t.Fatalf("cap %d", len(got.Selected))
	}
	if got.StaleExcluded != 1 {
		t.Fatalf("stale %d", got.StaleExcluded)
	}
	t.Logf("phase3f retrieval selected=%d bytes=%d rejected=%d stale_excluded=%d cap=%d",
		len(got.Selected), got.Bytes, got.Rejected, got.StaleExcluded, got.MaxFacts)
}

func TestPhase3F_LegacyAggregateSplit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "resource-memory.json")
	doc := ResourceMemoryDocument{
		SchemaVersion:    "1.0.0",
		TotalEvaluations: 2648,
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
	total, succ, fail, rate := store.GetSummary()
	if total != 0 || succ != 0 || fail != 0 || rate != 0 {
		t.Fatalf("current summary treated legacy as evidence: %d %d %d %f", total, succ, fail, rate)
	}
	if store.CurrentEvidenceCount() != 0 {
		t.Fatal("payload length was not zero")
	}
	sum, status := store.LegacyAggregate()
	if sum != 2648 || status != StatusLegacyAggregate {
		t.Fatalf("legacy = %d %s", sum, status)
	}
	if err := store.Record(&ResourceEvaluation{
		ResourceID: "gsap", Outcome: OutcomeSuccess, QualityScore: 1,
	}); err != nil {
		t.Fatal(err)
	}
	sum, status = store.LegacyAggregate()
	if sum != 2648 || status != StatusLegacyAggregate {
		t.Fatalf("record erased legacy: %d %s", sum, status)
	}
	if store.CurrentEvidenceCount() != 1 {
		t.Fatalf("current evidence %d", store.CurrentEvidenceCount())
	}
	total, _, _, _ = store.GetSummary()
	if total != 1 {
		t.Fatalf("summary total %d", total)
	}
}

func TestPhase3F_DuplicateMatrixAndStance(t *testing.T) {
	seen := map[string]bool{}
	for _, s := range DuplicateSurfaces() {
		if s.Action != "KEEP" && s.Action != "MIGRATE" {
			t.Fatalf("%s action %s", s.Name, s.Action)
		}
		if s.Class == "" {
			t.Fatalf("%s missing class", s.Name)
		}
		seen[s.Name] = true
	}
	for _, name := range []string{"PlaywrightVerifier", "Context7 MCP", "gsap-core skill", "shadcn MCP"} {
		if !seen[name] {
			t.Fatalf("missing surface %s", name)
		}
	}
	st := Stances()
	if st["shadcn"] != "OPTIONAL" || st["opendesign"] != "RETIRED" || st["jcode"] != "DEFER" || st["grok-4.7"] != "OPAQUE_LABEL" {
		t.Fatalf("stances %+v", st)
	}
}

func TestPhase3F_MissingLastVerifiedIsUnverified(t *testing.T) {
	verified, unverified := SplitVerification([]RegistryRow{
		{ID: "dated", LastVerified: "2026-09-06"},
		{ID: "blank", LastVerified: ""},
		{ID: "space", LastVerified: "   "},
	})
	if verified != 1 || unverified != 2 {
		t.Fatalf("verified=%d unverified=%d", verified, unverified)
	}
	if ResourceVerified("") || ResourceVerified("  ") {
		t.Fatal("missing date counted as verified")
	}
	catalog := filepath.Join("..", "..", "..", "registries", "resources.json")
	raw, err := os.ReadFile(catalog)
	if err != nil {
		t.Fatal(err)
	}
	var rows []RegistryRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	v, u := SplitVerification(rows)
	if v+u != len(rows) {
		t.Fatalf("split %d+%d != %d", v, u, len(rows))
	}
	if u == 0 {
		t.Fatal("checker reported every row verified")
	}
	for _, row := range rows {
		if strings.TrimSpace(row.LastVerified) == "" && ResourceVerified(row.LastVerified) {
			t.Fatalf("%s missing date reported verified", row.ID)
		}
	}
}

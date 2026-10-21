package benchmark

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/memory"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

func TestMain(m *testing.M) {
	if os.Getenv("ORCHESTRA_BENCH_FAKE") == "1" {
		os.Exit(fakeHostMain())
	}
	dir, err := os.MkdirTemp("", "orchestra-bench-mem-*")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	_ = os.Setenv("ORCHESTRA_MEMORY_PATH", filepath.Join(dir, "resource-memory.json"))
	_ = os.Setenv("ORCHESTRA_OVERLAY_PATH", filepath.Join(dir, "added-resources.json"))
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func loadCatalogs(t *testing.T) (*resources.ResourceCatalog, *resources.DesignResourceGraph) {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "..", "registries"),
		filepath.Join("..", "..", "registries"),
		`C:\projects\orchestra-workflow\registries`,
	}
	var regDir string
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			regDir = c
			break
		}
	}
	if regDir == "" {
		t.Fatal("registries not found")
	}
	cat, err := resources.LoadResourceCatalog(filepath.Join(regDir, "resources.json"))
	if err != nil {
		t.Fatal(err)
	}
	graph, err := resources.LoadDesignGraph(filepath.Join(regDir, "design-resource-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	return cat, graph
}

func fixtureDir(t *testing.T, id string) string {
	t.Helper()
	p := filepath.Join("..", "..", "benchmarks", "fixtures", id)
	if _, err := os.Stat(p); err != nil {
		t.Fatal(err)
	}
	return p
}

func taskByID(id string) Task {
	for _, task := range Manifest() {
		if task.ID == id {
			return task
		}
	}
	return Task{}
}

func TestPhase3G_ExpansionGateStaysShut(t *testing.T) {
	for _, task := range Manifest() {
		if task.Class == "T1" || task.Class == "T3" {
			if !task.Enabled {
				t.Fatalf("%s should be enabled", task.ID)
			}
			continue
		}
		if task.Enabled {
			t.Fatalf("%s enabled during the pilot", task.ID)
		}
	}
	if !ExpansionOpen(true, true, true, true) {
		t.Fatal("gate should be able to open later")
	}
	for _, task := range Manifest() {
		if task.Class != "T1" && task.Class != "T3" && task.Enabled {
			t.Fatal("opening the gate must not enable classes in this pass")
		}
	}
}

func TestPhase3G_ScoreIgnoresJudgment(t *testing.T) {
	rec, err := ScoreWorktree(fixtureDir(t, "t3-debug"), &Judgment{Note: "looks finished", Positive: true})
	if err != nil {
		t.Fatal(err)
	}
	if !rec.HardFail || rec.TestFailures == 0 {
		t.Fatalf("judgment overrode a test failure: %+v", rec.Deterministic)
	}
	ok, err := ScoreWorktree(fixtureDir(t, "t1-api"), &Judgment{Positive: false})
	if err != nil {
		t.Fatal(err)
	}
	if ok.HardFail {
		t.Fatalf("t1 hard fail: %+v", ok.Deterministic)
	}
}

func TestPhase3G_RecoveryStopsIdenticalLoop(t *testing.T) {
	rec, err := RecoverSeeded(mission.Budget{RetryBudget: 2})
	if err != nil {
		t.Fatal(err)
	}
	if rec.RecoveryClass != string(mission.ClassTest) || rec.RecoveryChange == string(mission.ChangeNone) {
		t.Fatalf("recovery %+v", rec)
	}
	if !rec.IdenticalRetryStopped || rec.Retries != 1 {
		t.Fatalf("loop not stopped: %+v", rec)
	}
}

func TestPhase3G_ModelLabelDoesNotChangeSelection(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "neu", Objective: "Add a small API endpoint that validates input", DeclaredType: "FEATURE",
	})
	pw := false
	var first []string
	for _, label := range []string{"grok-4.7", "claude-opus", "gemini-2.5", "gpt-x"} {
		d, err := capability.Select(capability.Input{
			Mission: m, PlaywrightAvailable: &pw, DiscoveredModels: []string{label},
		})
		if err != nil {
			t.Fatal(err)
		}
		var ids []string
		for _, it := range d.Selected {
			ids = append(ids, it.ID)
		}
		joined := strings.Join(ids, ",")
		if first == nil {
			first = ids
		} else if strings.Join(first, ",") != joined {
			t.Fatalf("%s changed ids to %s", label, joined)
		}
		raw := []byte(joined)
		if strings.Contains(strings.ToLower(string(raw)), "grok") || strings.Contains(joined, "claude") || strings.Contains(joined, "gemini") {
			t.Fatalf("label entered ids %s", joined)
		}
	}
}

func TestPhase3G_ScopePollution(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "t1-api", Objective: "Add a small API endpoint that validates input", DeclaredType: "FEATURE",
	})
	cap, _, sel, err := memory.BindCapsule(m, memory.SelectInput{
		Scope: "t1-api",
		Candidates: []memory.Fact{
			{ID: "other", Scope: "t9-reverse", Status: memory.StatusCurrent, Text: "UNRELATED_MEMORY"},
			{ID: "here", Scope: "t1-api", Status: memory.StatusCurrent, Text: "validate empty input"},
		},
	}, mission.CompileRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if sel.Rejected < 1 {
		t.Fatal("unrelated fact was not rejected")
	}
	blob := strings.Join(cap.Decisions, "\n")
	if strings.Contains(blob, "UNRELATED_MEMORY") {
		t.Fatal("unrelated memory entered the capsule")
	}
}

func TestPhase3G_PilotPairIsNotAClaim(t *testing.T) {
	rec := PairedHostSmoke("grok-4.7")
	if rec.Phase != PhasePilot || rec.Status != "unavailable" || rec.Claim != ClaimNone || rec.Trial != 1 {
		t.Fatalf("%+v", rec)
	}
	if rec.Tokens != TokensUnavailable {
		t.Fatal("invented tokens")
	}
}

func TestPhase3G_Traces(t *testing.T) {
	cat, graph := loadCatalogs(t)
	t.Setenv("ORCHESTRA_MEMORY_PATH", filepath.Join(t.TempDir(), "missing-brain.json"))
	for _, id := range []string{"t1-api", "t3-debug"} {
		var hashes []string
		var ids []string
		for trial := 0; trial < 2; trial++ {
			rec, err := RunTrace(context.Background(), TraceInput{
				Catalog:    cat,
				Graph:      graph,
				FixtureDir: fixtureDir(t, id),
				Task:       taskByID(id),
				ModelLabel: "grok-4.7",
				WorkRoot:   t.TempDir(),
			})
			if err != nil {
				t.Fatalf("%s trial %d: %v", id, trial, err)
			}
			if rec.Status != "trace_only" || rec.Claim != ClaimNone || rec.Phase != PhasePilot {
				t.Fatalf("%s status %+v", id, rec)
			}
			if rec.Tokens != TokensUnavailable || rec.TaskExecution != TaskCostUnavailable {
				t.Fatal("invented cost")
			}
			if rec.MissionBytes <= 0 || rec.CapabilityBytes <= 0 || rec.CapsuleBytes <= 0 || rec.OrchestrationBytes < rec.MissionBytes+rec.CapabilityBytes {
				t.Fatalf("%s bytes mission=%d capability=%d capsule=%d", id, rec.MissionBytes, rec.CapabilityBytes, rec.CapsuleBytes)
			}
			if rec.BrainFacts != 0 || rec.Directions != 0 || rec.BrowserCalls != 0 || rec.ResearchCalls != 0 || rec.JudgeCalls != 0 || rec.ScreenshotEvals != 0 {
				t.Fatalf("%s overhead brain=%d dir=%d browser=%d research=%d judge=%d shots=%d", id, rec.BrainFacts, rec.Directions, rec.BrowserCalls, rec.ResearchCalls, rec.JudgeCalls, rec.ScreenshotEvals)
			}
			if rec.Intensity != string(mission.IntensityLight) {
				t.Fatalf("%s intensity %s", id, rec.Intensity)
			}
			hashes = append(hashes, rec.FixtureHash)
			ids = append(ids, strings.Join(rec.CapabilityIDs, ","))
			t.Logf("%s trial %d wall_ms=%d mission=%d capability=%d capsule=%d ids=%s", id, trial, rec.WallMS, rec.MissionBytes, rec.CapabilityBytes, rec.CapsuleBytes, ids[len(ids)-1])
		}
		if hashes[0] != hashes[1] || ids[0] != ids[1] {
			t.Fatalf("%s not reproducible hash %v ids %v", id, hashes, ids)
		}
	}
}

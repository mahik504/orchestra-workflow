package benchmark

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fakeHostMain() int {
	root := os.Getenv("ORCHESTRA_BENCH_WORKSPACE")
	if root == "" {
		return 2
	}
	switch os.Getenv("ORCHESTRA_BENCH_FAKE_MODE") {
	case "fix":
		return rewriteFile(root, "sum.go", "a + b + 1", "a + b")
	case "cheat":
		return rewriteFile(root, "sum_test.go", "got != 3", "got != 4")
	case "extra":
		if rewriteFile(root, "sum.go", "a + b + 1", "a + b") != 0 {
			return 1
		}
		err := os.WriteFile(filepath.Join(root, "extra.go"), []byte("package t3debug\n"), 0644)
		if err != nil {
			return 1
		}
		return 0
	case "sleep":
		time.Sleep(20 * time.Second)
		return 0
	case "exit2":
		return 2
	case "leak":
		return rewriteFile(root, "sum.go", "package t3debug", "package t3debug // "+os.Getenv("ORCHESTRA_BENCH_MODEL_LABEL"))
	case "dirtybase":
		raw, err := os.ReadFile(os.Getenv("ORCHESTRA_BENCH_PROMPT"))
		if err != nil {
			return 1
		}
		if strings.Contains(string(raw), "Capsule file:") {
			return 0
		}
		if err := os.MkdirAll(filepath.Join(root, ".orchestra"), 0755); err != nil {
			return 1
		}
		return 0
	case "touch":
		return rewriteFile(root, "sum.go", "package t3debug", "package t3debug // touched")
	case "fixorch":
		raw, err := os.ReadFile(os.Getenv("ORCHESTRA_BENCH_PROMPT"))
		if err != nil {
			return 1
		}
		if !strings.Contains(string(raw), "Capsule file:") {
			return 0
		}
		return rewriteFile(root, "sum.go", "a + b + 1", "a + b")
	case "outside":
		target := os.Getenv("ORCHESTRA_BENCH_OUTSIDE")
		if target == "" {
			return 1
		}
		if err := os.WriteFile(target, []byte("x"), 0644); err != nil {
			return 1
		}
		return 0
	case "t1impl":
		return writeT1Done(root, true)
	case "t1code":
		return writeT1Done(root, false)
	default:
		return 0
	}
}

func writeT1Done(root string, withTest bool) int {
	code := "package t1api\n\nimport \"errors\"\n\nfunc Sum(a, b int) (int, error) {\n\tif a < 0 || b < 0 {\n\t\treturn 0, errors.New(\"invalid\")\n\t}\n\treturn a + b, nil\n}\n"
	if err := os.WriteFile(filepath.Join(root, "sum.go"), []byte(code), 0644); err != nil {
		return 1
	}
	if !withTest {
		return 0
	}
	testPath := filepath.Join(root, "sum_test.go")
	raw, err := os.ReadFile(testPath)
	if err != nil {
		return 1
	}
	if strings.Contains(string(raw), "TestSumRejectsNegative") {
		return 0
	}
	extra := "\nfunc TestSumRejectsNegative(t *testing.T) {\n\tif _, err := Sum(-1, 2); err == nil {\n\t\tt.Fatal(\"expected an error\")\n\t}\n}\n"
	if err := os.WriteFile(testPath, append(raw, []byte(extra)...), 0644); err != nil {
		return 1
	}
	return 0
}

func rewriteFile(root, name, old, new string) int {
	path := filepath.Join(root, name)
	raw, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(raw), old) {
		return 1
	}
	next := strings.Replace(string(raw), old, new, 1)
	if err := os.WriteFile(path, []byte(next), 0644); err != nil {
		return 1
	}
	return 0
}

func pairInput(t *testing.T, id, mode string, timeout time.Duration) PairInput {
	t.Helper()
	t.Setenv("ORCHESTRA_BENCH_FAKE", "1")
	t.Setenv("ORCHESTRA_BENCH_FAKE_MODE", mode)
	cat, graph := loadCatalogs(t)
	root := t.TempDir()
	work := filepath.Join(root, "work")
	if err := os.MkdirAll(work, 0755); err != nil {
		t.Fatal(err)
	}
	return PairInput{
		Catalog:        cat,
		Graph:          graph,
		FixtureDir:     fixtureDir(t, id),
		Task:           taskByID(id),
		WorkRoot:       work,
		Executable:     os.Args[0],
		Timeout:        timeout,
		HostLabel:      "fake",
		BaselineLabel:  "grok-4.7",
		OrchestraLabel: "grok-4.7",
		ResultPath:     filepath.Join(root, "pair.json"),
	}
}

func TestPhase3G1_Unavailable(t *testing.T) {
	work := filepath.Join(t.TempDir(), "work")
	res, err := RunPair(context.Background(), PairInput{Task: Task{ID: "t1-api"}, WorkRoot: work})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "unavailable" || res.Claim != ClaimNone || res.Tokens != TokensUnavailable {
		t.Fatalf("%+v", res)
	}
	if _, err := os.Stat(work); !os.IsNotExist(err) {
		t.Fatalf("worktree remained: %v", err)
	}
}

func TestPhase3G1_LabelMismatch(t *testing.T) {
	res, err := RunPair(context.Background(), PairInput{
		Task:           Task{ID: "t1-api"},
		Executable:     "missing-host",
		BaselineLabel:  "grok-4.7",
		OrchestraLabel: "other-label",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "invalid_pair" || res.Claim != ClaimNone {
		t.Fatalf("%+v", res)
	}
}

func TestPhase3G1_T1Noop(t *testing.T) {
	in := pairInput(t, "t1-api", "noop", 30*time.Second)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone || !res.BaselineIsolated || !res.WorktreeRemoved {
		t.Fatalf("status %s isolated %v removed %v", res.Status, res.BaselineIsolated, res.WorktreeRemoved)
	}
	if _, err := os.Stat(in.WorkRoot); !os.IsNotExist(err) {
		t.Fatal("worktree still present")
	}
	for _, rec := range []Record{res.Baseline, res.Orchestra} {
		build, bok := checkNamed(rec, DimBuild)
		tests, tok := checkNamed(rec, DimTests)
		behavior, behOK := checkNamed(rec, DimRequiredBehavior)
		accept, aok := checkNamed(rec, DimAcceptance)
		scope, sok := checkNamed(rec, DimScope)
		if !bok || !tok || !behOK || !aok || !sok || !rec.HardFail {
			t.Fatalf("expected a hard fail: %+v", rec.Deterministic)
		}
		if !build.Pass || !tests.Pass || behavior.Pass || !accept.Pass || !scope.Pass || rec.TaskCompleted || rec.HostEdited {
			t.Fatalf("noop should fail required_behavior only: %+v edited=%v done=%v", rec.Deterministic, rec.HostEdited, rec.TaskCompleted)
		}
	}
	if res.Orchestra.Intensity != "LIGHT" || res.Orchestra.Directions != 0 || res.Orchestra.BrowserCalls != 0 || res.Orchestra.BrainFacts != 0 {
		t.Fatalf("overhead %+v", res.Orchestra)
	}
	if res.Orchestra.MissionBytes <= 0 || res.Orchestra.CapsuleBytes <= 0 {
		t.Fatalf("bytes %+v", res.Orchestra)
	}
	if res.Baseline.Tokens != TokensUnavailable || res.Orchestra.HostRetries != TokensUnavailable {
		t.Fatal("invented host cost")
	}
	raw, err := os.ReadFile(in.ResultPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"claim": "none"`) || !strings.Contains(string(raw), `"status": "compared"`) {
		t.Fatalf("result %s", raw)
	}
}

func TestPhase3G1_T3Fix(t *testing.T) {
	in := pairInput(t, "t3-debug", "fix", 30*time.Second)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || !res.BaselineIsolated || res.Claim != ClaimNone {
		t.Fatalf("%s %v", res.Status, res.BaselineIsolated)
	}
	if res.Baseline.HardFail || res.Orchestra.HardFail {
		t.Fatalf("base %+v orch %+v", res.Baseline.Deterministic, res.Orchestra.Deterministic)
	}
	if res.Baseline.Mode != "off" || res.Orchestra.Mode != "on" {
		t.Fatalf("modes %s %s", res.Baseline.Mode, res.Orchestra.Mode)
	}
	for _, rec := range []Record{res.Baseline, res.Orchestra} {
		if rec.Artifact == nil || len(rec.Artifact.Changed) == 0 {
			t.Fatalf("missing diff %+v", rec.Artifact)
		}
	}
}

func TestPhase3G1_T3Cheat(t *testing.T) {
	in := pairInput(t, "t3-debug", "cheat", 30*time.Second)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" {
		t.Fatal(res.Status)
	}
	tests, ok := checkNamed(res.Orchestra, DimTests)
	accept, aok := checkNamed(res.Orchestra, DimAcceptance)
	if !ok || !aok || !tests.Pass || accept.Pass || !res.Orchestra.HardFail {
		t.Fatalf("tests %+v accept %+v hard %v", tests, accept, res.Orchestra.HardFail)
	}
	ApplyScore(&res.Orchestra, res.Orchestra.Deterministic, &Judgment{Positive: true, Note: "looks finished"})
	if !res.Orchestra.HardFail {
		t.Fatal("judgment cleared a cheated test")
	}
}

func TestPhase3G1_T3ExtraFile(t *testing.T) {
	in := pairInput(t, "t3-debug", "extra", 30*time.Second)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	scope, ok := checkNamed(res.Orchestra, DimScope)
	tests, tok := checkNamed(res.Orchestra, DimTests)
	if !ok || !tok || scope.Pass || !tests.Pass || !res.Orchestra.HardFail {
		t.Fatalf("scope %+v tests %+v hard %v art %+v", scope, tests, res.Orchestra.HardFail, res.Orchestra.Artifact)
	}
}

func TestPhase3G1_Timeout(t *testing.T) {
	in := pairInput(t, "t3-debug", "sleep", 400*time.Millisecond)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "timeout" || res.Claim != ClaimNone {
		t.Fatalf("%+v", res.Status)
	}
	if res.Baseline.Artifact == nil || res.Orchestra.Artifact == nil || !res.Baseline.Artifact.TimedOut || !res.Orchestra.Artifact.TimedOut {
		t.Fatalf("base %+v orch %+v", res.Baseline.Artifact, res.Orchestra.Artifact)
	}
	for _, line := range res.Limitations {
		if strings.Contains(line, "not a win") {
			return
		}
	}
	t.Fatal("timeout was treated as a win")
}

func TestPhase3G1_Prompts(t *testing.T) {
	base := baselinePrompt("Add a small API sum", []string{"negative input returns an error"})
	if strings.Contains(strings.ToLower(base), "grok") || strings.Contains(base, "capsule") || strings.Contains(strings.ToLower(base), "winner") {
		t.Fatal(base)
	}
	on := orchestraPrompt("Fix Sum", []string{"Sum(1, 2) returns 3"}, "mission.json", "capsule.json", "capability.json")
	if !strings.Contains(on, "capsule.json") || strings.Contains(strings.ToLower(on), "grok-4.7") {
		t.Fatal(on)
	}
	if !textHasLabel("model grok-4.7", "grok-4.7") || textHasLabel("coding.worker", "grok-4.7") {
		t.Fatal("label scanner")
	}
}

func TestPhase3G1_ExpansionStaysShut(t *testing.T) {
	for _, task := range Manifest() {
		if task.Class != "T1" && task.Class != "T3" && task.Enabled {
			t.Fatalf("%s enabled", task.ID)
		}
	}
}

func TestPhase3G1_RealHostPilot(t *testing.T) {
	if os.Getenv("ORCHESTRA_BENCH_REAL") != "1" {
		t.Skip("real host pair is explicit; default tests use the fake host")
	}
	exe := os.Getenv("ORCHESTRA_BENCH_HOST")
	label := os.Getenv("ORCHESTRA_BENCH_MODEL_LABEL")
	if exe == "" {
		for _, id := range []string{"t1-api", "t3-debug"} {
			res, err := RunPair(context.Background(), PairInput{Task: taskByID(id), WorkRoot: filepath.Join(t.TempDir(), "work")})
			if err != nil || res.Status != "unavailable" || res.Claim != ClaimNone {
				t.Fatalf("%s %+v %v", id, res, err)
			}
		}
		return
	}
	cat, graph := loadCatalogs(t)
	for _, id := range []string{"t1-api", "t3-debug"} {
		root := t.TempDir()
		res, err := RunPair(context.Background(), PairInput{
			Catalog:        cat,
			Graph:          graph,
			FixtureDir:     fixtureDir(t, id),
			Task:           taskByID(id),
			WorkRoot:       filepath.Join(root, "work"),
			Executable:     exe,
			Timeout:        10 * time.Minute,
			HostLabel:      "cursor",
			BaselineLabel:  label,
			OrchestraLabel: label,
			ResultPath:     filepath.Join(root, id+".json"),
		})
		if err != nil {
			t.Fatal(err)
		}
		if res.Claim != ClaimNone || res.Trial != 1 || res.Phase != PhasePilot {
			t.Fatalf("%s %+v", id, res)
		}
		t.Logf("%s status=%s isolated=%v", id, res.Status, res.BaselineIsolated)
	}
}

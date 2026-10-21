package benchmark

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPhase3G23_T1NoopNotCompleted(t *testing.T) {
	in := pairInput(t, "t1-api", "noop", 30*time.Second)
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone {
		t.Fatalf("status %s claim %s", res.Status, res.Claim)
	}
	for _, rec := range []Record{res.Baseline, res.Orchestra} {
		behavior, ok := checkNamed(rec, DimRequiredBehavior)
		if !ok || behavior.Pass || rec.TaskCompleted || rec.HostEdited {
			t.Fatalf("noop done=%v edited=%v behavior=%+v", rec.TaskCompleted, rec.HostEdited, behavior)
		}
		if rec.Artifact != nil {
			for _, name := range rec.Artifact.Changed {
				if name == oracleFileName {
					t.Fatal("oracle file was scored as an edit")
				}
			}
		}
	}
	raw, err := os.ReadFile(in.ResultPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), oracleFileName) || strings.Contains(string(raw), "TestBenchOracle") {
		t.Fatal("oracle leaked into the result")
	}
}

func TestPhase3G23_T1ImplCompleted(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t1-api", "t1impl", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone {
		t.Fatalf("status %s claim %s", res.Status, res.Claim)
	}
	for _, rec := range []Record{res.Baseline, res.Orchestra} {
		tests, tok := checkNamed(rec, DimTests)
		behavior, bok := checkNamed(rec, DimRequiredBehavior)
		scope, sok := checkNamed(rec, DimScope)
		if !tok || !bok || !sok || !tests.Pass || !behavior.Pass || !scope.Pass || rec.HardFail || !rec.TaskCompleted {
			t.Fatalf("impl tests=%+v behavior=%+v scope=%+v hard=%v done=%v", tests, behavior, scope, rec.HardFail, rec.TaskCompleted)
		}
		if !changedPath(rec, "sum.go") || !changedPath(rec, "sum_test.go") {
			t.Fatalf("changed %+v", rec.Artifact)
		}
	}
}

func TestPhase3G23_T1NoTestNotCompleted(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t1-api", "t1code", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone {
		t.Fatalf("status %s claim %s", res.Status, res.Claim)
	}
	for _, rec := range []Record{res.Baseline, res.Orchestra} {
		behavior, ok := checkNamed(rec, DimRequiredBehavior)
		if !ok || !behavior.Pass || rec.HardFail || rec.TaskCompleted || !rec.HostEdited {
			t.Fatalf("code-only behavior=%+v hard=%v done=%v edited=%v", behavior, rec.HardFail, rec.TaskCompleted, rec.HostEdited)
		}
		if rec.Artifact == nil || rec.Artifact.TimedOut || rec.Artifact.HostExit != 0 || changedPath(rec, "sum_test.go") {
			t.Fatalf("artifact %+v", rec.Artifact)
		}
	}
}

func TestPhase3G23_OracleNotLeaked(t *testing.T) {
	if _, err := os.Stat(filepath.Join(fixtureDir(t, "t1-api"), oracleFileName)); !os.IsNotExist(err) {
		t.Fatal("oracle shipped in the fixture")
	}
	root := t.TempDir()
	if err := CopyTree(fixtureDir(t, "t1-api"), root); err != nil {
		t.Fatal(err)
	}
	task, acc, err := ReadTask(root)
	if err != nil {
		t.Fatal(err)
	}
	base := baselinePrompt(task, acc)
	on := orchestraPrompt(task, acc, "mission.json", "capsule.json", "capability.json")
	for _, text := range []string{base, on} {
		if strings.Contains(text, "TestBenchOracle") || strings.Contains(text, oracleFileName) {
			t.Fatal("oracle text is in a prompt")
		}
	}
	if err := os.WriteFile(filepath.Join(root, "prompt.txt"), []byte(base), 0644); err != nil {
		t.Fatal(err)
	}
	if err := commitSnapshot(root); err != nil {
		t.Fatal(err)
	}
	policy, err := PolicyFor("t1-api")
	if err != nil {
		t.Fatal(err)
	}
	before, err := protectHashes(root, policy.Protect)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := ScoreHost(root, policy, before, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, oracleFileName)); !os.IsNotExist(err) {
		t.Fatal("oracle file remained after scoring")
	}
	if changedPath(rec, oracleFileName) {
		t.Fatal("oracle was scored as an edit")
	}
	behavior, ok := checkNamed(rec, DimRequiredBehavior)
	tests, tok := checkNamed(rec, DimTests)
	if !ok || !tok || behavior.Pass || !tests.Pass {
		t.Fatalf("starting tree tests=%+v behavior=%+v", tests, behavior)
	}
}

func TestPhase3G23_TimeoutArmNotCompleted(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "sleep", 400*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "timeout" || res.Claim != ClaimNone {
		t.Fatalf("status %s", res.Status)
	}
	if res.Baseline.TaskCompleted || res.Orchestra.TaskCompleted {
		t.Fatal("timed out arm was completed")
	}
	if res.Baseline.Artifact == nil || !res.Baseline.Artifact.TimedOut {
		t.Fatalf("artifact %+v", res.Baseline.Artifact)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sum.go"), []byte("package t3debug\n\nfunc Sum(a, b int) int {\n\treturn a + b\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	art := *res.Baseline.Artifact
	art.TimedOut = true
	art.HostExit = -1
	passing := passingRecord(&art)
	if taskCompleted("t3-debug", dir, passing) {
		t.Fatal("timeout counted as completion")
	}
	art.TimedOut = false
	art.HostExit = 0
	if !taskCompleted("t3-debug", dir, passing) {
		t.Fatal("a finished arm should count once the seeded defect is gone")
	}
}

func TestPhase3G23_ExitArmNotCompleted(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "exit2", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "host_error" || res.Claim != ClaimNone {
		t.Fatalf("status %s", res.Status)
	}
	if res.Baseline.TaskCompleted || res.Orchestra.TaskCompleted {
		t.Fatal("non-zero exit was completed")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sum.go"), []byte("package t3debug\n\nfunc Sum(a, b int) int {\n\treturn a + b\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	art := Artifact{HostExit: 2}
	if taskCompleted("t3-debug", dir, passingRecord(&art)) {
		t.Fatal("non-zero exit counted as completion")
	}
	art.HostExit = 0
	if !taskCompleted("t3-debug", dir, passingRecord(&art)) {
		t.Fatal("exit 0 should count once the seeded defect is gone")
	}
}

func TestPhase3G23_HostEnvDeniesHome(t *testing.T) {
	t.Setenv("ORCHESTRA_HOME", "bench-home-marker")
	env := hostEnv("ws", "prompt.txt", "cache", "mod", "label")
	for _, kv := range env {
		key, _, _ := strings.Cut(kv, "=")
		if key == "ORCHESTRA_HOME" {
			t.Fatal("ORCHESTRA_HOME was passed to the host")
		}
	}
}

func TestPhase3G23_BrainGuardFromHome(t *testing.T) {
	home := t.TempDir()
	ledger := filepath.Join(home, "memory", "resource-memory.json")
	if err := os.MkdirAll(filepath.Dir(ledger), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ledger, []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ORCHESTRA_HOME", home)
	t.Setenv("ORCHESTRA_MEMORY_PATH", "")
	guards := realGuards("")
	var sawHome, sawLedger bool
	for _, g := range guards {
		if g.Class == "brain" && g.Path == home {
			sawHome = true
		}
		if g.Class == "brain" && g.Path == ledger {
			sawLedger = true
		}
	}
	if !sawHome || !sawLedger {
		t.Fatal("brain guards missing")
	}
	snaps := snapshotGuards(guards)
	if err := os.WriteFile(ledger, []byte("{}\n\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if guardsChanged(snaps) != "brain" {
		t.Fatal("ledger write was not caught")
	}
}

func TestPhase3G23_T3InputsUnchanged(t *testing.T) {
	path := filepath.Join("..", "..", "benchmarks", "results", "3g2-t3-debug.json")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skip("saved t3 result is not in this checkout")
	}
	if err != nil {
		t.Fatal(err)
	}
	var saved struct {
		FixtureHash    string `json:"fixture_hash"`
		TaskHash       string `json:"task_hash"`
		AcceptanceHash string `json:"acceptance_hash"`
	}
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatal(err)
	}
	got, err := HashTree(fixtureDir(t, "t3-debug"))
	if err != nil {
		t.Fatal(err)
	}
	if got != saved.FixtureHash {
		t.Fatalf("fixture hash changed")
	}
	task, acc, err := ReadTask(fixtureDir(t, "t3-debug"))
	if err != nil {
		t.Fatal(err)
	}
	if shortHash(task) != saved.TaskHash || shortHash(strings.Join(acc, "\n")) != saved.AcceptanceHash {
		t.Fatal("task or acceptance hash changed")
	}
	policy, err := PolicyFor("t3-debug")
	if err != nil || policy.Oracle != "" {
		t.Fatal("t3 gained an oracle")
	}
}

func TestPhase3G23_RealT1(t *testing.T) {
	realPilot(t, "3g23", []string{"t1-api"})
}

func TestPhase3G24_RealT3(t *testing.T) {
	realPilot(t, "3g24", []string{"t3-debug"})
}

func passingRecord(art *Artifact) Record {
	return Record{
		Deterministic: []Check{
			{Name: DimBuild, Pass: true},
			{Name: DimTests, Pass: true},
			{Name: DimRequiredBehavior, Pass: true},
			{Name: DimAcceptance, Pass: true},
			{Name: DimScope, Pass: true},
		},
		Artifact: art,
	}
}

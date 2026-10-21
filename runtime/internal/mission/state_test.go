package mission

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStateRoundTripAndContinue(t *testing.T) {
	dir := t.TempDir()
	m := validMission()
	m.RepoPath = dir
	st := NewRunState(m)
	st.CurrentStage = "Implement"
	st.CurrentWorker = string(RoleBuilder)
	st.Attempt = 1
	st.Outputs = []EvidenceRef{{Name: "test-log", Path: "tmp/out.txt"}}
	st.NextAction = "run tests"
	path, err := WriteRunState(dir, st)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadRunState(dir, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.MissionID != m.ID || got.CurrentStage != "Implement" || got.Attempt != 1 {
		t.Fatalf("restore mismatch %+v", got)
	}
	if len(got.Outputs) != 1 {
		t.Fatalf("outputs lost")
	}
	why, err := ContinueWithoutConversation(got, m)
	if err != nil {
		t.Fatalf("continue: %v", err)
	}
	if why != "run tests" {
		t.Fatalf("next action %q", why)
	}
}

func TestMissionFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	m := validMission()
	if _, err := WriteMissionFile(dir, m); err != nil {
		t.Fatal(err)
	}
	got, err := ReadMissionFile(dir, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.AcceptanceFingerprint != m.AcceptanceFingerprint {
		t.Fatalf("fingerprint lost")
	}
}

func TestContinueStopsOnHuman(t *testing.T) {
	m := validMission()
	st := NewRunState(m)
	st.PendingHuman = &HumanDecision{Reason: "destructive deploy", Blocking: true}
	st.Status = StatusWaitingHuman
	if ok, _ := st.CanContinue(); ok {
		t.Fatalf("human-needed state continued")
	}
	if _, err := ContinueWithoutConversation(st, m); err == nil {
		t.Fatalf("expected continue error")
	}
}

func TestFullProgressStuckDetection(t *testing.T) {
	dir := t.TempDir()
	m := validMission()
	m.Intensity = IntensityFull
	st := NewRunState(m)
	for i := 0; i < 3; i++ {
		if err := WriteProgress(dir, st, IntensityFull, "still waiting on same step"); err != nil {
			t.Fatal(err)
		}
	}
	if !st.Stuck {
		t.Fatalf("expected stuck")
	}
	raw, err := os.ReadFile(filepath.Join(RunDir(dir, m.ID), "progress.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatalf("progress.txt empty")
	}
	m.Intensity = IntensityLight
	st2 := NewRunState(m)
	if err := WriteProgress(dir, st2, IntensityLight, "light should skip file"); err != nil {
		t.Fatal(err)
	}
}

func TestFromSeedDoesNotDumpContext(t *testing.T) {
	m := FromSeed(Seed{
		ID:            "seed-1",
		Objective:     "Fix nil pointer in parser",
		DeclaredType:  "BUGFIX",
		WorkspaceRoot: "/tmp/app",
	})
	if err := Validate(m); err != nil {
		t.Fatal(err)
	}
	if m.Type != TypeDebug {
		t.Fatalf("type %s", m.Type)
	}
	if m.Intensity != IntensityLight {
		t.Fatalf("intensity %s", m.Intensity)
	}
	if !CheapBudget(m.Budget) {
		t.Fatalf("debug seed should be cheap: %+v", m.Budget)
	}
	if m.AcceptanceFingerprint == "" {
		t.Fatalf("unsealed")
	}
}

func TestInterruptedContinueDoesNotNeedConversation(t *testing.T) {
	dir := t.TempDir()
	m := validMission()
	st := NewRunState(m)
	st.Status = StatusRunning
	st.CurrentStage = "Implement"
	st.CurrentWorker = string(RoleBuilder)
	st.Attempt = 2
	st.NextAction = "retry_with_approach"
	st.FailureHistory = []AttemptRecord{{
		Attempt:         1,
		FailureClass:    ClassImplementation,
		PreviousFailure: "build failed: does not compile",
		RecoveryAction:  ActionRetry,
		Change:          ChangeApproach,
		ResultingState:  StatusRetrying,
	}}
	if _, err := WriteRunState(dir, st); err != nil {
		t.Fatal(err)
	}
	got, err := ReadRunState(dir, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	why, err := ContinueWithoutConversation(got, m)
	if err != nil {
		t.Fatalf("continue: %v", err)
	}
	if why != "retry_with_approach" {
		t.Fatalf("next %q", why)
	}
	raw, _ := json.Marshal(got)
	if strings.Contains(string(raw), "yesterday we talked") || strings.Contains(string(raw), "AGENTS.md") {
		t.Fatalf("continuation state carried conversation/project dump")
	}
}

package resources

import (
	"path/filepath"
	"testing"
)

func TestProjectContractDoesNotInventUsage(t *testing.T) {
	path := filepath.Join("..", "..", "..", "registries", "resources.json")
	cat, err := LoadResourceCatalog(path)
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	if len(cat.allResources) < 10 {
		t.Fatalf("catalog too small: %d", len(cat.allResources))
	}
	for _, res := range cat.allResources {
		got, err := ProjectContract(res)
		if err != nil {
			t.Fatalf("%s: %v", res.ID, err)
		}
		if got.Lifecycle != LifeDiscovered && !hasRunEvidence(got.UsageEvidence) {
			t.Fatalf("%s lifecycle %s without run evidence", res.ID, got.Lifecycle)
		}
		if got.UsageEvidence != UsageNoneRecorded && !hasRunEvidence(got.UsageEvidence) {
			t.Fatalf("%s usage %q is not a recorded run", res.ID, got.UsageEvidence)
		}
		if got.Verification == "verified" && !hasRunEvidence(got.UsageEvidence) {
			t.Fatalf("%s verified without evidence", res.ID)
		}
	}
}

func TestPromoteCannotSkipToSuccessful(t *testing.T) {
	if _, err := Promote(LifeDiscovered, LifeSuccessful); err == nil {
		t.Fatal("discovered must not jump to successful")
	}
	next, err := Promote(LifeDiscovered, LifeInstalled)
	if err != nil || next != LifeInstalled {
		t.Fatalf("one step: %s %v", next, err)
	}
}

func TestExecutedLifecycleNeedsEvidence(t *testing.T) {
	res := &Resource{
		ID:                "example-tool",
		Name:              "Example",
		TriggerConditions: []string{"when named"},
		SkipConditions:    []string{"when absent"},
		LifecycleStage:    LifeSuccessful,
		UsageEvidence:     "the registry lists it",
	}
	if _, err := ProjectContract(res); err == nil {
		t.Fatal("a registry sentence is not run evidence")
	}
	res.UsageEvidence = "run: 2026-11-07 example"
	if _, err := ProjectContract(res); err != nil {
		t.Fatal(err)
	}
}

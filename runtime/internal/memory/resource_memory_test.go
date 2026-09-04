package memory

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/user/orchestra-v3/internal/resources"
)

func TestNewResourceMemoryStore_Initialization(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "test-memory.json")

	store, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	if store.doc.TotalEvaluations != 0 {
		t.Errorf("Expected 0 evaluations, got %d", store.doc.TotalEvaluations)
	}
	if len(store.doc.Resources) != 0 {
		t.Errorf("Expected 0 resources, got %d", len(store.doc.Resources))
	}
}

func TestRecordEvaluation_Aggregations(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "test-memory.json")

	store, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	// 1. Record Success
	err = store.Record(&ResourceEvaluation{
		ResourceID:   "playwright",
		Domain:       "qa_testing",
		Capability:   "visual-regression",
		TaskID:       "task-1",
		Outcome:      OutcomeSuccess,
		QualityScore: 1.0,
		LatencyMs:    500,
	})
	if err != nil {
		t.Fatalf("Record success failed: %v", err)
	}

	// 2. Record Second Success
	err = store.Record(&ResourceEvaluation{
		ResourceID:   "playwright",
		Domain:       "qa_testing",
		Capability:   "visual-regression",
		TaskID:       "task-2",
		Outcome:      OutcomeSuccess,
		QualityScore: 0.9,
		LatencyMs:    700,
	})
	if err != nil {
		t.Fatalf("Record second success failed: %v", err)
	}

	// 3. Record Failure
	err = store.Record(&ResourceEvaluation{
		ResourceID:   "playwright",
		Domain:       "qa_testing",
		Capability:   "visual-regression",
		TaskID:       "task-3",
		Outcome:      OutcomeFailure,
		QualityScore: 0.5,
		LatencyMs:    900,
		ErrorDetails: "viewport overflow",
	})
	if err != nil {
		t.Fatalf("Record failure failed: %v", err)
	}

	agg, ok := store.GetAggregate("playwright")
	if !ok {
		t.Fatalf("Aggregate for playwright not found")
	}

	if agg.TotalEvaluations != 3 {
		t.Errorf("Expected 3 evals, got %d", agg.TotalEvaluations)
	}
	if agg.SuccessCount != 2 {
		t.Errorf("Expected 2 successes, got %d", agg.SuccessCount)
	}
	if agg.FailureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", agg.FailureCount)
	}

	expectedRate := 2.0 / 3.0
	if agg.SuccessRate < expectedRate-0.001 || agg.SuccessRate > expectedRate+0.001 {
		t.Errorf("Expected success rate ~%.3f, got %.3f", expectedRate, agg.SuccessRate)
	}

	expectedAvgLatency := (500.0 + 700.0 + 900.0) / 3.0
	if agg.AverageLatencyMs != expectedAvgLatency {
		t.Errorf("Expected avg latency %.1f, got %.1f", expectedAvgLatency, agg.AverageLatencyMs)
	}

	expectedAvgScore := (1.0 + 0.9 + 0.5) / 3.0
	if agg.AverageQualityScore < expectedAvgScore-0.001 || agg.AverageQualityScore > expectedAvgScore+0.001 {
		t.Errorf("Expected avg score %.3f, got %.3f", expectedAvgScore, agg.AverageQualityScore)
	}

	if agg.LastOutcome != OutcomeFailure {
		t.Errorf("Expected last outcome failure, got %s", agg.LastOutcome)
	}

	// Test GetSummary
	total, succ, fail, rate := store.GetSummary()
	if total != 3 || succ != 2 || fail != 1 {
		t.Errorf("Summary mismatch: total=%d, succ=%d, fail=%d", total, succ, fail)
	}
	if rate < expectedRate-0.001 || rate > expectedRate+0.001 {
		t.Errorf("Summary rate mismatch: expected %.3f, got %.3f", expectedRate, rate)
	}

	// Test ListEvaluations
	evals := store.ListEvaluations(10, "playwright")
	if len(evals) != 3 {
		t.Errorf("Expected 3 evals, got %d", len(evals))
	}
	// Reverse chronological: most recent first
	if evals[0].TaskID != "task-3" {
		t.Errorf("Expected most recent eval task-3, got %s", evals[0].TaskID)
	}
}

func TestQuarantineEnforcement(t *testing.T) {
	quarantinePath := `C:\Users\mockuser\.gemini\config\skills_library\resource-memory.json`
	_, err := NewResourceMemoryStore(quarantinePath)
	if err == nil {
		t.Fatalf("Expected quarantine error when initializing in skills_library, got nil")
	}
	if !errors.Is(err, resources.ErrQuarantinedPath) {
		t.Errorf("Expected error wrapping ErrQuarantinedPath, got %v", err)
	}
}

func TestThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "concurrent-memory.json")

	store, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	var wg sync.WaitGroup
	workers := 10
	iterations := 5

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = store.Record(&ResourceEvaluation{
					ResourceID:   "react-bits",
					Domain:       "component_library",
					Capability:   "rich-components",
					TaskID:       "concurrent-task",
					Outcome:      OutcomeSuccess,
					QualityScore: 1.0,
					LatencyMs:    100,
				})
			}
		}(w)
	}

	wg.Wait()

	agg, ok := store.GetAggregate("react-bits")
	if !ok {
		t.Fatalf("react-bits aggregate missing")
	}
	expectedTotal := workers * iterations
	if agg.TotalEvaluations != expectedTotal {
		t.Errorf("Expected %d total evals, got %d", expectedTotal, agg.TotalEvaluations)
	}
}

func TestReloadPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "persist-memory.json")

	store1, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Init 1 failed: %v", err)
	}

	_ = store1.Record(&ResourceEvaluation{
		ResourceID:   "awwwards",
		Domain:       "visual_design",
		Capability:   "premium-editorial-web",
		TaskID:       "persist-task",
		Outcome:      OutcomeSuccess,
		QualityScore: 1.0,
		LatencyMs:    300,
	})

	// Reload from disk
	store2, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	agg, ok := store2.GetAggregate("awwwards")
	if !ok || agg.TotalEvaluations != 1 {
		t.Fatalf("Reloaded store missing aggregate data: %+v", agg)
	}
	if agg.SuccessCount != 1 || agg.AverageLatencyMs != 300.0 {
		t.Errorf("Reloaded data mismatch: %+v", agg)
	}
}

func TestInvalidEvaluationInputs(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "invalid-input-memory.json")

	store, err := NewResourceMemoryStore(memFile)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Nil eval
	if err := store.Record(nil); !errors.Is(err, ErrInvalidEvaluation) {
		t.Errorf("Expected ErrInvalidEvaluation for nil, got: %v", err)
	}

	// Missing resource_id
	if err := store.Record(&ResourceEvaluation{Outcome: OutcomeSuccess}); err == nil {
		t.Errorf("Expected error for empty resource_id")
	}

	// Invalid outcome
	if err := store.Record(&ResourceEvaluation{ResourceID: "foo", Outcome: "maybe"}); err == nil {
		t.Errorf("Expected error for invalid outcome")
	}
}

// A real workspace's memory file must parse and stay internally consistent.
// It asserts nothing about which resources are present: an empty memory file is
// the correct state for a fresh install, and seeding it would be fabricated
// learning rather than recorded outcomes.
func TestLiveWorkspaceMemoryFileValid(t *testing.T) {
	memPath := os.Getenv("ORCHESTRA_MEMORY_PATH")
	if memPath == "" {
		t.Skip("ORCHESTRA_MEMORY_PATH not set, skipping live workspace check")
	}
	if _, err := os.Stat(memPath); os.IsNotExist(err) {
		t.Skipf("no memory file at %s, skipping", memPath)
	}

	store, err := NewResourceMemoryStore(memPath)
	if err != nil {
		t.Fatalf("Failed to load workspace memory: %v", err)
	}

	for _, agg := range store.ListAggregates() {
		if agg.TotalEvaluations < 0 {
			t.Errorf("resource %s has negative evaluation count", agg.ResourceID)
		}
	}
}

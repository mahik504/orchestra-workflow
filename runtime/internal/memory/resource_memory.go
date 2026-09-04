package memory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

var (
	// ErrInvalidEvaluation is returned when required evaluation fields are missing
	ErrInvalidEvaluation = errors.New("invalid evaluation entry: missing required fields")
	// ErrResourceNotFound is returned when querying a resource not in the memory store
	ErrResourceNotFound = errors.New("resource not found in memory store")
	// ErrMemoryCorrupted is returned when memory JSON cannot be decoded
	ErrMemoryCorrupted = errors.New("resource memory JSON corrupted or invalid")
)

// Outcome represents the success or failure outcome of a resource evaluation
type Outcome string

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
)

// ResourceEvaluation records a single execution verdict and telemetry data point
type ResourceEvaluation struct {
	EvaluationID        string         `json:"evaluation_id"`
	ResourceID          string         `json:"resource_id"`
	Domain              string         `json:"domain"`
	Capability          string         `json:"capability"`
	EvaluationTimestamp string         `json:"evaluation_timestamp"`
	TaskContext         string         `json:"task_context"`
	TaskID              string         `json:"task_id"`
	Outcome             Outcome        `json:"outcome"`
	QualityScore        float64        `json:"quality_score"`
	LatencyMs           int64          `json:"latency_ms"`
	ErrorDetails        string         `json:"error_details,omitempty"`
	Notes               string         `json:"notes,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
}

// ResourceAggregate maintains rolling statistical performance metrics for a resource
type ResourceAggregate struct {
	ResourceID          string   `json:"resource_id"`
	Domain              string   `json:"domain,omitempty"`
	Capability          string   `json:"capability,omitempty"`
	TotalEvaluations    int      `json:"total_evaluations"`
	SuccessCount        int      `json:"success_count"`
	FailureCount        int      `json:"failure_count"`
	SuccessRate         float64  `json:"success_rate"`
	AverageLatencyMs    float64  `json:"average_latency_ms"`
	AverageQualityScore float64  `json:"average_quality_score"`
	LastUsedTimestamp   string   `json:"last_used_timestamp"`
	LastOutcome         Outcome  `json:"last_outcome"`
	Tags                []string `json:"tags,omitempty"`
}

// ResourceMemoryDocument represents the root JSON schema structure
type ResourceMemoryDocument struct {
	Schema           string                        `json:"$schema,omitempty"`
	SchemaVersion    string                        `json:"schema_version"`
	LastUpdated      string                        `json:"last_updated"`
	TotalEvaluations int                           `json:"total_evaluations"`
	Resources        map[string]*ResourceAggregate `json:"resources"`
	Evaluations      []ResourceEvaluation          `json:"evaluations"`
}

// ResourceMemoryStore provides thread-safe, durable, atomic persistence for resource memory
type ResourceMemoryStore struct {
	filePath string
	doc      ResourceMemoryDocument
	mu       sync.RWMutex
}

// ResolveDefaultMemoryPath resolves the prioritized path to resource-memory.json.
//
// Resolution is relative to the caller's workspace so that a fresh clone writes
// memory into that clone's own workspace. Point ORCHESTRA_HOME at a private
// workspace to keep resource memory outside the repository.
func ResolveDefaultMemoryPath(workspaceRoot string) string {
	if envPath := os.Getenv("ORCHESTRA_MEMORY_PATH"); envPath != "" {
		return filepath.Clean(envPath)
	}

	var candidates []string
	if home := os.Getenv("ORCHESTRA_HOME"); home != "" {
		candidates = append(candidates, filepath.Join(home, "memory", "resource-memory.json"))
	}
	candidates = append(candidates,
		filepath.Join(workspaceRoot, "memory", "resource-memory.json"),
		filepath.Join(workspaceRoot, ".orchestra", "memory", "resource-memory.json"),
		filepath.Join(workspaceRoot, "..", "orchestra-workspace", "memory", "resource-memory.json"),
	)

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return filepath.Clean(c)
		}
	}

	// Nothing on disk yet. Keep Orchestra's own state under .orchestra/ rather
	// than dropping a memory/ directory into somebody else's repository.
	return filepath.Clean(filepath.Join(workspaceRoot, ".orchestra", "memory", "resource-memory.json"))
}

// NewResourceMemoryStore loads or initializes a ResourceMemoryStore at filePath
func NewResourceMemoryStore(filePath string) (*ResourceMemoryStore, error) {
	if strings.TrimSpace(filePath) == "" {
		filePath = ResolveDefaultMemoryPath(".")
	}

	if err := resources.CheckQuarantineBoundary(filePath); err != nil {
		return nil, fmt.Errorf("memory store path violates quarantine: %w", err)
	}

	store := &ResourceMemoryStore{
		filePath: filepath.Clean(filePath),
		doc: ResourceMemoryDocument{
			Schema:           "https://orchestra.workflow/schemas/resource-memory.schema.json",
			SchemaVersion:    "1.0.0",
			LastUpdated:      time.Now().UTC().Format(time.RFC3339),
			TotalEvaluations: 0,
			Resources:        make(map[string]*ResourceAggregate),
			Evaluations:      make([]ResourceEvaluation, 0),
		},
	}

	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %v", ErrMemoryCorrupted, err)
	}

	return store, nil
}

func (s *ResourceMemoryStore) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // Strip UTF-8 BOM

	var doc ResourceMemoryDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}

	if doc.Resources == nil {
		doc.Resources = make(map[string]*ResourceAggregate)
	}
	if doc.Evaluations == nil {
		doc.Evaluations = make([]ResourceEvaluation, 0)
	}

	s.doc = doc
	return nil
}

// Record appends an evaluation record and updates running historical aggregations
func (s *ResourceMemoryStore) Record(eval *ResourceEvaluation) error {
	if eval == nil {
		return ErrInvalidEvaluation
	}

	resID := strings.ToLower(strings.TrimSpace(eval.ResourceID))
	if resID == "" {
		return fmt.Errorf("%w: resource_id is required", ErrInvalidEvaluation)
	}
	if eval.Outcome != OutcomeSuccess && eval.Outcome != OutcomeFailure {
		return fmt.Errorf("%w: outcome must be 'success' or 'failure'", ErrInvalidEvaluation)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	eval.ResourceID = resID
	if eval.EvaluationTimestamp == "" {
		eval.EvaluationTimestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if eval.EvaluationID == "" {
		eval.EvaluationID = fmt.Sprintf("eval-%d-%s", time.Now().UnixNano(), resID)
	}

	// Clamp quality score
	if eval.QualityScore < 0.0 {
		eval.QualityScore = 0.0
	} else if eval.QualityScore > 1.0 {
		eval.QualityScore = 1.0
	}

	// Update aggregate
	agg, exists := s.doc.Resources[resID]
	if !exists {
		agg = &ResourceAggregate{
			ResourceID: resID,
			Domain:     eval.Domain,
			Capability: eval.Capability,
		}
		s.doc.Resources[resID] = agg
	}

	if agg.Domain == "" && eval.Domain != "" {
		agg.Domain = eval.Domain
	}
	if agg.Capability == "" && eval.Capability != "" {
		agg.Capability = eval.Capability
	}

	agg.TotalEvaluations++
	if eval.Outcome == OutcomeSuccess {
		agg.SuccessCount++
	} else {
		agg.FailureCount++
	}

	agg.SuccessRate = float64(agg.SuccessCount) / float64(agg.TotalEvaluations)

	// Running statistical averages
	n := float64(agg.TotalEvaluations)
	agg.AverageLatencyMs = (agg.AverageLatencyMs*(n-1.0) + float64(eval.LatencyMs)) / n
	agg.AverageQualityScore = (agg.AverageQualityScore*(n-1.0) + eval.QualityScore) / n
	agg.LastUsedTimestamp = eval.EvaluationTimestamp
	agg.LastOutcome = eval.Outcome

	// Append evaluation
	s.doc.Evaluations = append(s.doc.Evaluations, *eval)
	s.doc.TotalEvaluations = len(s.doc.Evaluations)
	s.doc.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	return s.saveLocked()
}

// saveLocked performs an atomic, Windows-safe file write
func (s *ResourceMemoryStore) saveLocked() error {
	if err := resources.CheckQuarantineBoundary(s.filePath); err != nil {
		return fmt.Errorf("cannot persist memory: %w", err)
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed creating directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(s.doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed encoding memory JSON: %w", err)
	}

	tmpFile := fmt.Sprintf("%s.tmp.%d", s.filePath, time.Now().UnixNano())
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed writing tmp memory file: %w", err)
	}

	// Atomic replace with Windows-safe fallback
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		_ = os.Remove(s.filePath)
		if err := os.Rename(tmpFile, s.filePath); err != nil {
			// Direct write fallback
			if err := os.WriteFile(s.filePath, data, 0644); err != nil {
				_ = os.Remove(tmpFile)
				return fmt.Errorf("failed writing memory file: %w", err)
			}
			_ = os.Remove(tmpFile)
		}
	}

	return nil
}

// GetAggregate returns the aggregated metrics for a specific resource
func (s *ResourceMemoryStore) GetAggregate(resourceID string) (*ResourceAggregate, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agg, exists := s.doc.Resources[strings.ToLower(resourceID)]
	if !exists {
		return nil, false
	}
	cp := *agg
	return &cp, true
}

// ListAggregates returns all resource aggregates sorted alphabetically by ID
func (s *ResourceMemoryStore) ListAggregates() []*ResourceAggregate {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*ResourceAggregate, 0, len(s.doc.Resources))
	for _, agg := range s.doc.Resources {
		cp := *agg
		result = append(result, &cp)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ResourceID < result[j].ResourceID
	})

	return result
}

// ListEvaluations returns recorded evaluations in reverse chronological order
func (s *ResourceMemoryStore) ListEvaluations(limit int, resourceIDFilter string) []ResourceEvaluation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := strings.ToLower(strings.TrimSpace(resourceIDFilter))
	var result []ResourceEvaluation

	for i := len(s.doc.Evaluations) - 1; i >= 0; i-- {
		e := s.doc.Evaluations[i]
		if filter != "" && e.ResourceID != filter {
			continue
		}
		result = append(result, e)
		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result
}

// GetSummary returns macro-level totals and overall success rate
func (s *ResourceMemoryStore) GetSummary() (total int, successes int, failures int, avgSuccessRate float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total = s.doc.TotalEvaluations
	for _, agg := range s.doc.Resources {
		successes += agg.SuccessCount
		failures += agg.FailureCount
	}
	if total > 0 {
		avgSuccessRate = float64(successes) / float64(total)
	}
	return
}

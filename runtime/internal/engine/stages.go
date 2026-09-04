package engine

import (
	"errors"
	"time"
)

// Standard Engine Errors
var (
	ErrHumanGateRequired    = errors.New("human approval gate required before writing code")
	ErrMaxIterationsExceeded = errors.New("maximum QA iteration loops exceeded without resolution")
	ErrQuarantinedPath       = errors.New("quarantined path violation")
)

// StageName identifies one of the 8 canonical stages in the design-first pipeline
type StageName string

const (
	StageDiscover     StageName = "Discover"
	StageClassify     StageName = "Classify"
	StageResearch     StageName = "Research"
	StageSynthesize   StageName = "Synthesize"
	StageDesignSystem StageName = "DesignSystem"
	StageImplement    StageName = "Implement"
	StageVisualQA     StageName = "VisualQA"
	StageIterate      StageName = "Iterate"
)

// StageStatus represents the execution state of a stage
type StageStatus string

const (
	StatusPending   StageStatus = "PENDING"
	StatusRunning   StageStatus = "RUNNING"
	StatusCompleted StageStatus = "COMPLETED"
	StatusSkipped   StageStatus = "SKIPPED"
	StatusFailed    StageStatus = "FAILED"
)

// StageResult captures the timing, artifacts, and output of a single stage run
type StageResult struct {
	StageName  StageName     `json:"stage_name"`
	Status     StageStatus   `json:"status"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	Output     any           `json:"output,omitempty"`
	Error      error         `json:"error,omitempty"`
	SkipReason string        `json:"skip_reason,omitempty"`
	Artifacts  []string      `json:"artifacts,omitempty"`
}

// Stage defines the execution contract for each modular pipeline phase
type Stage interface {
	Name() StageName
	Execute(ctx *TaskContext) (*StageResult, error)
	ShouldSkip(ctx *TaskContext) (bool, string)
}

package handoff

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileChecksum struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type HandoffState struct {
	SessionID       string         `json:"session_id"`
	Version         int            `json:"version"`
	Timestamp       string         `json:"timestamp"`
	SourceAgent     string         `json:"source_agent"` // e.g. "antigravity"
	TargetAgent     string         `json:"target_agent"` // e.g. "cursor"
	ActiveTasks     []string       `json:"active_tasks"`
	PlanURI         string         `json:"plan_uri"`
	ChangedFiles    []FileChecksum `json:"changed_files"`
	CompletedSteps  []string       `json:"completed_steps"`
	PendingSteps    []string       `json:"pending_steps"`
	FailureRecovery *RecoveryPoint `json:"failure_recovery,omitempty"`
}

type RecoveryPoint struct {
	FailedStep     string `json:"failed_step"`
	ErrorReason    string `json:"error_reason"`
	CanResume      bool   `json:"can_resume"`
	ResumeFromStep string `json:"resume_from_step"`
}

// ComputeFileHash calculates SHA256 of a local file
func ComputeFileHash(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}

// WriteState serializes the versioned handoff state
func WriteState(state *HandoffState, workdir string) error {
	if state.Timestamp == "" {
		state.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	handoffDir := filepath.Join(workdir, ".orchestra", "handoff")
	if err := os.MkdirAll(handoffDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(handoffDir, "state.json")
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// ReadState deserializes the handoff state
func ReadState(workdir string) (*HandoffState, error) {
	path := filepath.Join(workdir, ".orchestra", "handoff", "state.json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no handoff state found at %s: %w", path, err)
	}

	var state HandoffState
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, fmt.Errorf("corrupt handoff state: %w", err)
	}

	return &state, nil
}

// DetectConflicts checks whether tracked files were modified out-of-band
func DetectConflicts(state *HandoffState, baseDir string) ([]string, error) {
	var conflicts []string
	for _, f := range state.ChangedFiles {
		fullPath := filepath.Join(baseDir, f.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			conflicts = append(conflicts, fmt.Sprintf("%s was deleted externally", f.Path))
			continue
		}
		currentHash, err := ComputeFileHash(fullPath)
		if err != nil {
			return nil, err
		}
		if currentHash != f.SHA256 {
			conflicts = append(conflicts, fmt.Sprintf("%s content changed externally (expected %s, got %s)", f.Path, f.SHA256[:8], currentHash[:8]))
		}
	}
	return conflicts, nil
}

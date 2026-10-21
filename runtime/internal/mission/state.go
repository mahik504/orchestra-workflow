package mission

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type EvidenceRef struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type HumanDecision struct {
	Reason   string `json:"reason"`
	Question string `json:"question,omitempty"`
	Blocking bool   `json:"blocking"`
}

type RunState struct {
	SchemaVersion  int             `json:"schema_version"`
	MissionID      string          `json:"mission_id"`
	Status         MissionStatus   `json:"status"`
	CurrentStage   string          `json:"current_stage,omitempty"`
	CurrentWorker  string          `json:"current_worker,omitempty"`
	Attempt        int             `json:"attempt"`
	Budget         Budget          `json:"budget"`
	FailureHistory []AttemptRecord `json:"failure_history,omitempty"`
	Outputs        []EvidenceRef   `json:"outputs,omitempty"`
	PendingHuman   *HumanDecision  `json:"pending_human,omitempty"`
	NextAction     string          `json:"next_action,omitempty"`
	CapsuleScope   *ScopeReport    `json:"capsule_scope,omitempty"`
	Progress       []string        `json:"progress,omitempty"`
	Stuck          bool            `json:"stuck,omitempty"`
	UpdatedAt      string          `json:"updated_at"`
}

func NewRunState(m *Mission) *RunState {
	if m == nil {
		return nil
	}
	return &RunState{
		SchemaVersion: SchemaVersion,
		MissionID:     m.ID,
		Status:        m.Status,
		CurrentWorker: string(RoleBuilder),
		Attempt:       m.Run.Attempt,
		Budget:        m.Budget,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		NextAction:    "execute",
	}
}

func RunDir(workspace, missionID string) string {
	id := SanitizeID(missionID)
	return filepath.Join(workspace, ".orchestra", "runs", id)
}

func SanitizeID(id string) string {
	id = strings.TrimSpace(id)
	id = strings.ReplaceAll(id, `\`, "_")
	id = strings.ReplaceAll(id, `/`, "_")
	if id == "" {
		return "unnamed"
	}
	return id
}

func WriteMissionFile(workspace string, m *Mission) (string, error) {
	if m == nil {
		return "", fmt.Errorf("mission is nil")
	}
	dir := RunDir(workspace, m.ID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	raw, err := Marshal(m)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "mission.json")
	return path, os.WriteFile(path, raw, 0644)
}

func ReadMissionFile(workspace, missionID string) (*Mission, error) {
	path := filepath.Join(RunDir(workspace, missionID), "mission.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no mission file at %s: %w", path, err)
	}
	return Unmarshal(raw)
}

func WriteRunState(workspace string, st *RunState) (string, error) {
	if st == nil {
		return "", fmt.Errorf("run state is nil")
	}
	st.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	dir := RunDir(workspace, st.MissionID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	raw, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "state.json")
	return path, os.WriteFile(path, raw, 0644)
}

func ReadRunState(workspace, missionID string) (*RunState, error) {
	path := filepath.Join(RunDir(workspace, missionID), "state.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no run state at %s: %w", path, err)
	}
	var st RunState
	if err := json.Unmarshal(raw, &st); err != nil {
		return nil, fmt.Errorf("corrupt run state: %w", err)
	}
	return &st, nil
}

func WriteProgress(workspace string, st *RunState, intensity Intensity, line string) error {
	if st == nil || intensity != IntensityFull {
		return nil
	}
	st.AppendProgress(line)
	dir := RunDir(workspace, st.MissionID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "progress.txt")
	var b strings.Builder
	for _, l := range st.Progress {
		b.WriteString(l)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}

func (s *RunState) AppendProgress(line string) {
	if s == nil {
		return
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	s.Progress = append(s.Progress, line)
	if len(s.Progress) >= 3 {
		n := len(s.Progress)
		if s.Progress[n-1] == s.Progress[n-2] && s.Progress[n-2] == s.Progress[n-3] {
			s.Stuck = true
		}
	}
}

func (s *RunState) CanContinue() (bool, string) {
	if s == nil {
		return false, "no run state"
	}
	if s.PendingHuman != nil && s.PendingHuman.Blocking {
		return false, "pending human decision: " + s.PendingHuman.Reason
	}
	if s.Status == StatusWaitingHuman {
		return false, "status is waiting_human"
	}
	if s.Stuck {
		return false, "FULL progress is stuck"
	}
	if s.Status == StatusSucceeded || s.Status == StatusFailed || s.Status == StatusExhausted {
		return false, "run already terminal"
	}
	if s.NextAction == "" {
		return true, "resume current stage"
	}
	return true, s.NextAction
}

func ContinueWithoutConversation(st *RunState, m *Mission) (string, error) {
	if st == nil || m == nil {
		return "", fmt.Errorf("mission and run state are required to continue")
	}
	ok, why := st.CanContinue()
	if !ok {
		return "", fmt.Errorf("%s", why)
	}
	return why, nil
}

package mission

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	RolePlanner WorkerRole = "planner"
	RoleBuilder WorkerRole = "builder"
	RoleJudge   WorkerRole = "judge"
)

type WorkerRole string

type FailurePolicy struct {
	RetryBudget          int            `json:"retry_budget"`
	IdenticalReplayStops bool           `json:"identical_replay_stops"`
	HumanClasses         []FailureClass `json:"human_classes,omitempty"`
}

type Capsule struct {
	MissionID          string          `json:"mission_id"`
	Role               WorkerRole      `json:"role"`
	Objective          string          `json:"objective"`
	Constraints        []string        `json:"constraints,omitempty"`
	Surfaces           []string        `json:"surfaces,omitempty"`
	Files              []string        `json:"files,omitempty"`
	OwnedFiles         []string        `json:"owned_files,omitempty"`
	Skills             []string        `json:"skills,omitempty"`
	Resources          []string        `json:"resources,omitempty"`
	ExpectedOutput     string          `json:"expected_output,omitempty"`
	Acceptance         []string        `json:"acceptance"`
	Evidence           []string        `json:"evidence,omitempty"`
	Budget             Budget          `json:"budget"`
	FailurePolicy      FailurePolicy   `json:"failure_policy"`
	PreviousAttempt    *AttemptRecord  `json:"previous_attempt,omitempty"`
	KnownFailures      []FailureRecord `json:"known_failures,omitempty"`
	Decisions          []string        `json:"decisions,omitempty"`
	DesignContractPath string          `json:"design_contract_path,omitempty"`
	ToolsAllowed       []string        `json:"tools_allowed,omitempty"`
}

type ScopeReport struct {
	CapsuleBytes      int      `json:"capsule_bytes"`
	FileRefs          int      `json:"file_refs"`
	SkillRefs         int      `json:"skill_refs"`
	ConstraintRefs    int      `json:"constraint_refs"`
	Excluded          []string `json:"excluded,omitempty"`
	WholeProjectBytes int      `json:"whole_project_bytes,omitempty"`
	Ratio             float64  `json:"ratio,omitempty"`
}

type CompileRequest struct {
	Mission             *Mission
	Role                WorkerRole
	Files               []string
	Skills              []string
	Resources           []string
	ToolsAllowed        []string
	ExpectedOutput      string
	Previous            *AttemptRecord
	KnownFailures       []FailureRecord
	Decisions           []string
	DesignContractPath  string
	WholeProjectContext []byte
	ConversationHistory []string
	MCPCatalog          []string
	AllSkills           []string
	UnrelatedFiles      []string
}

func Compile(req CompileRequest) (*Capsule, *ScopeReport, error) {
	if req.Mission == nil {
		return nil, nil, fmt.Errorf("capsule requires a mission")
	}
	if err := Validate(req.Mission); err != nil {
		return nil, nil, err
	}
	role := req.Role
	if role == "" {
		role = RoleBuilder
	}
	switch role {
	case RolePlanner, RoleBuilder, RoleJudge:
	default:
		return nil, nil, fmt.Errorf("worker role must be planner, builder, or judge")
	}

	skills := intersect(req.Skills, req.AllSkills)
	if len(req.AllSkills) == 0 {
		skills = uniqueTrim(req.Skills)
	}
	files := uniqueTrim(req.Files)
	unrelated := map[string]bool{}
	for _, u := range req.UnrelatedFiles {
		unrelated[u] = true
	}
	scopedFiles := make([]string, 0, len(files))
	for _, f := range files {
		if unrelated[f] {
			continue
		}
		scopedFiles = append(scopedFiles, f)
	}

	c := &Capsule{
		MissionID:      req.Mission.ID,
		Role:           role,
		Objective:      req.Mission.Objective,
		Constraints:    append([]string{}, req.Mission.Constraints...),
		Surfaces:       append([]string{}, req.Mission.AffectedSurfaces...),
		Files:          scopedFiles,
		OwnedFiles:     append([]string{}, scopedFiles...),
		Skills:         skills,
		Resources:      uniqueTrim(req.Resources),
		ExpectedOutput: req.ExpectedOutput,
		Acceptance:     CloneAcceptance(req.Mission.Acceptance),
		Evidence:       append([]string{}, req.Mission.RequiredEvidence...),
		Budget:         req.Mission.Budget,
		FailurePolicy: FailurePolicy{
			RetryBudget:          req.Mission.Budget.RetryBudget,
			IdenticalReplayStops: true,
			HumanClasses:         []FailureClass{ClassHuman, ClassSpecification},
		},
		PreviousAttempt:    req.Previous,
		KnownFailures:      append([]FailureRecord{}, req.KnownFailures...),
		Decisions:          uniqueTrim(req.Decisions),
		DesignContractPath: req.DesignContractPath,
		ToolsAllowed:       uniqueTrim(req.ToolsAllowed),
	}
	if role == RoleJudge {
		c.ToolsAllowed = filterJudgeTools(c.ToolsAllowed)
		if c.ExpectedOutput == "" {
			c.ExpectedOutput = "independent evaluation findings; no writes"
		}
	}

	raw, err := json.Marshal(c)
	if err != nil {
		return nil, nil, err
	}
	if leak := ProviderLeak(raw); leak != "" && !strings.Contains(strings.ToLower(req.Mission.Objective), leak) {
		return nil, nil, fmt.Errorf("provider-specific token %q leaked into worker capsule", leak)
	}

	rep := &ScopeReport{
		CapsuleBytes:   len(raw),
		FileRefs:       len(c.Files),
		SkillRefs:      len(c.Skills),
		ConstraintRefs: len(c.Constraints),
		Excluded:       []string{},
	}
	if len(req.ConversationHistory) > 0 {
		rep.Excluded = append(rep.Excluded, "conversation_history")
	}
	if len(req.MCPCatalog) > 0 {
		rep.Excluded = append(rep.Excluded, "mcp_catalog")
	}
	if len(req.AllSkills) > len(c.Skills) {
		rep.Excluded = append(rep.Excluded, "unrelated_skills")
	}
	if len(req.UnrelatedFiles) > 0 {
		rep.Excluded = append(rep.Excluded, "unrelated_files")
	}
	if len(req.WholeProjectContext) > 0 {
		rep.Excluded = append(rep.Excluded, "whole_project_context")
		rep.WholeProjectBytes = len(req.WholeProjectContext)
		if rep.WholeProjectBytes > 0 {
			rep.Ratio = float64(rep.CapsuleBytes) / float64(rep.WholeProjectBytes)
		}
	}
	return c, rep, nil
}

func filterJudgeTools(in []string) []string {
	var out []string
	for _, t := range in {
		low := strings.ToLower(t)
		if strings.Contains(low, "write") || strings.Contains(low, "edit") ||
			strings.Contains(low, "delete") || strings.Contains(low, "apply_patch") ||
			strings.Contains(low, "shell") {
			continue
		}
		out = append(out, t)
	}
	return out
}

func OwnershipConflict(capsules []*Capsule) []string {
	owner := map[string]string{}
	var conflicts []string
	for _, capsule := range capsules {
		if capsule == nil {
			continue
		}
		for _, file := range capsule.OwnedFiles {
			key := ownedKey(file)
			prev, seen := owner[key]
			if seen && prev != capsule.MissionID {
				conflicts = append(conflicts, file+" claimed by "+prev+" and "+capsule.MissionID)
				continue
			}
			owner[key] = capsule.MissionID
		}
	}
	return conflicts
}

func ownedKey(file string) string {
	file = strings.TrimSpace(file)
	file = strings.ReplaceAll(file, "\\", "/")
	return strings.ToLower(file)
}

func uniqueTrim(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func intersect(want, universe []string) []string {
	if len(universe) == 0 {
		return uniqueTrim(want)
	}
	allow := map[string]bool{}
	for _, s := range universe {
		allow[strings.TrimSpace(s)] = true
	}
	var out []string
	seen := map[string]bool{}
	for _, s := range want {
		s = strings.TrimSpace(s)
		if s == "" || !allow[s] || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

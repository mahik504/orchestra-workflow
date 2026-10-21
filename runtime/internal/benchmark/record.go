package benchmark

// Record is one comparison cell. Phase pilot is not an effectiveness claim.
// Extra trials and task classes fit the same struct. There is no overall score.
type Record struct {
	BenchmarkID           string    `json:"benchmark_id"`
	TaskClass             string    `json:"task_class"`
	FixtureID             string    `json:"fixture_id"`
	FixtureHash           string    `json:"fixture_hash"`
	Clean                 bool      `json:"clean"`
	HostLabel             string    `json:"host_label,omitempty"`
	ModelLabel            string    `json:"model_label,omitempty"`
	Mode                  string    `json:"mode"`
	Intensity             string    `json:"intensity,omitempty"`
	Phase                 string    `json:"phase"`
	Trial                 int       `json:"trial"`
	Start                 string    `json:"start,omitempty"`
	End                   string    `json:"end,omitempty"`
	Status                string    `json:"status"`
	Claim                 string    `json:"claim"`
	AcceptanceFailures    int       `json:"acceptance_failures"`
	TestFailures          int       `json:"test_failures"`
	ScopeFailures         int       `json:"scope_failures,omitempty"`
	Retries               int       `json:"retries"`
	HumanGates            int       `json:"human_gates"`
	UnnecessaryGates      int       `json:"unnecessary_gates"`
	WallMS                int64     `json:"wall_ms"`
	MissionBytes          int       `json:"mission_bytes"`
	CapabilityBytes       int       `json:"capability_bytes"`
	CapsuleBytes          int       `json:"capsule_bytes"`
	BrowserCalls          int       `json:"browser_calls"`
	ResearchCalls         int       `json:"research_calls"`
	JudgeCalls            int       `json:"judge_calls"`
	BrainFacts            int       `json:"brain_facts"`
	ScreenshotEvals       int       `json:"screenshot_evals"`
	Directions            int       `json:"directions"`
	OrchestrationBytes    int       `json:"orchestration_bytes"`
	TaskExecution         string    `json:"task_execution"`
	Tokens                string    `json:"tokens"`
	HostRetries           string    `json:"host_retries,omitempty"`
	HostEdited            bool      `json:"host_edited,omitempty"`
	TaskCompleted         bool      `json:"task_completed,omitempty"`
	EndingHash            string    `json:"ending_hash,omitempty"`
	TaskHash              string    `json:"task_hash,omitempty"`
	AcceptanceHash        string    `json:"acceptance_hash,omitempty"`
	PromptHash            string    `json:"prompt_hash,omitempty"`
	Deterministic         []Check   `json:"deterministic,omitempty"`
	Artifact              *Artifact `json:"artifact,omitempty"`
	Judgment              *Judgment `json:"judgment,omitempty"`
	HardFail              bool      `json:"hard_fail"`
	CapabilityIDs         []string  `json:"capability_ids,omitempty"`
	Limitations           []string  `json:"limitations,omitempty"`
	RecoveryClass         string    `json:"recovery_class,omitempty"`
	RecoveryChange        string    `json:"recovery_change,omitempty"`
	IdenticalRetryStopped bool      `json:"identical_retry_stopped,omitempty"`
}

// Artifact is the measured diff. It is not a second verdict and not an overall score.
type Artifact struct {
	DiffStat       string   `json:"diff_stat,omitempty"`
	DiffBytes      int      `json:"diff_bytes"`
	Changed        []string `json:"changed,omitempty"`
	OutsideAllow   []string `json:"outside_allow,omitempty"`
	NewDeps        bool     `json:"new_deps"`
	HostExit       int      `json:"host_exit"`
	TimedOut       bool     `json:"timed_out"`
	HostWallMS     int64    `json:"host_wall_ms"`
	StdoutBytes    int      `json:"stdout_bytes,omitempty"`
	StderrBytes    int      `json:"stderr_bytes,omitempty"`
	HarnessChanged bool     `json:"harness_changed"`
}

package engine

import (
	"context"
	"sync"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/research"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/router"
	"github.com/user/orchestra-v3/internal/verify"
)

// Failure classification constants for closed-loop self-healing
const (
	FailureClassNone       = "NONE"
	FailureClassTokenStyle = "TOKEN_STYLE"
	FailureClassLayoutCode = "LAYOUT_CODE"
	FailureClassFatal      = "FATAL"
)

// TaskRequest defines the intake contract for the design-first engine
type TaskRequest struct {
	ID                 string                   `json:"id"`
	RawRequest         string                   `json:"raw_request"`
	WorkspaceRoot      string                   `json:"workspace_root"`
	Type               string                   `json:"type"` // "FEATURE", "DESIGN", "BUGFIX", "BENCHMARK"
	ArchetypeHint      string                   `json:"archetype_hint,omitempty"`
	Tags               []string                 `json:"tags,omitempty"`
	SuggestedResources []string                 `json:"suggested_resources,omitempty"`
	UserOverride       *classifier.UserOverride `json:"user_override,omitempty"`
	MaxIterations      int                      `json:"max_iterations"` // Default: 3
	Timeout            time.Duration            `json:"timeout"`        // Default: 10m
	SkipVisualGate     bool                     `json:"skip_visual_gate"`
	DryRun             bool                     `json:"dry_run"`
	Parameters         map[string]any           `json:"parameters,omitempty"`
}

// DiscoveryData contains project and environment findings from Stage 1
type DiscoveryData struct {
	Framework       string            `json:"framework"` // e.g. "react", "vite", "next", "vanilla"
	PackageJSONPath string            `json:"package_json_path"`
	Dependencies    map[string]string `json:"dependencies"`
	ExistingTokens  bool              `json:"existing_tokens"`
	DetectedTags    []string          `json:"detected_tags"`
}

// ClassificationData contains capability mappings and governance gates from Stage 2
type ClassificationData struct {
	Archetype           string                          `json:"archetype"`
	NormalizedTags      []string                        `json:"normalized_tags"`
	ResolvedRoutes      []resources.CapabilityRoute     `json:"resolved_routes"`
	RequiresVisual      bool                            `json:"requires_visual"`
	RequiresSecurity    bool                            `json:"requires_security"`
	RequiresHumanGate   bool                            `json:"requires_human_gate"`
	GateReason          string                          `json:"gate_reason"`
	GapTechnologies     []string                        `json:"gap_technologies"`
	RecommendedAgent    router.AllocationRecommendation `json:"recommended_agent"`
	Brief               *classifier.Brief               `json:"brief,omitempty"`
	QualityBar          string                          `json:"quality_bar"`
	Platform            string                          `json:"platform"`
	ResearchDepth       string                          `json:"research_depth"`
	VerifyDepth         string                          `json:"verify_depth"`
	DeclinedRoutes      []classifier.Candidate          `json:"declined_routes,omitempty"`
	OverlayActivations  []string                        `json:"overlay_activations,omitempty"`
	OverlaySuppressions []string                        `json:"overlay_suppressions,omitempty"`
}

// ResearchData contains visual inspirations, palettes, typography from Stage 3
type ResearchData struct {
	ReferenceEntries   []research.ScoredReference `json:"reference_entries"`
	ReferenceLogPath   string                     `json:"reference_log_path"`
	ColorInspirations  []string                   `json:"color_inspirations"`
	TypographyPairings []string                   `json:"typography_pairings"`
	MotionPatterns     []string                   `json:"motion_patterns"`
	GapDocumentation   map[string]string          `json:"gap_documentation"`
	SynthesizedResult  *research.ResearchResult   `json:"synthesized_result,omitempty"`
}

// SynthesisData contains actionable rules and token weights from Stage 4
type SynthesisData struct {
	ActionableRules    []string                               `json:"actionable_rules"`
	AntiPatterns       []string                               `json:"anti_patterns"`
	VerificationChecks []string                               `json:"verification_checks"`
	TokenContextCost   float64                                `json:"token_context_cost"`
	DirectionApproved  bool                                   `json:"direction_approved"`
	ActiveDirectives   []*router.CapabilityExecutionDirective `json:"active_directives"`
}

// DesignSystemData contains concrete tokens, styles, and DESIGN.md path from Stage 5
type DesignSystemData struct {
	DesignMDPath  string            `json:"design_md_path"`
	DisplayFont   string            `json:"display_font"`
	BodyFont      string            `json:"body_font"`
	MonoFont      string            `json:"mono_font"`
	PrimaryAccent string            `json:"primary_accent"` // Saturation < 80%
	SurfaceColor  string            `json:"surface_color"`  // No pure #000000
	SpringPhysics map[string]string `json:"spring_physics"` // damping, stiffness, mass
	LayoutRules   []string          `json:"layout_rules"`
	CSSVariables  map[string]string `json:"css_variables,omitempty"`
}

// ImplementationData contains modified files, acquired packages from Stage 6
type ImplementationData struct {
	TargetAgent       string                 `json:"target_agent"`
	AcquiredResources []string               `json:"acquired_resources"`
	ModifiedFiles     []handoff.FileChecksum `json:"modified_files"`
	HandoffStatePath  string                 `json:"handoff_state_path"`
	BuildOutput       string                 `json:"build_output"`
	BuildPassed       bool                   `json:"build_passed"`
	InstalledPaths    map[string]string      `json:"installed_paths,omitempty"`
	AcquireCommands   map[string]string      `json:"acquire_commands,omitempty"`
}

// ViewportCheckResult defines QA output for a specific responsive viewport
type ViewportCheckResult struct {
	ViewportName       string   `json:"viewport_name"` // "desktop", "tablet", "mobile"
	Width              int      `json:"width"`
	Height             int      `json:"height"`
	ScreenshotPath     string   `json:"screenshot_path"`
	HasOverflow        bool     `json:"has_overflow"`
	MaxScrollWidth     int      `json:"max_scroll_width"`
	Passed             bool     `json:"passed"`
	ContrastViolations []string `json:"contrast_violations,omitempty"`
	AntiPatternMatches []string `json:"anti_pattern_matches,omitempty"`
}

// VisualQAData captures verification verdict and detected defects from Stage 7
type VisualQAData struct {
	AllPassed          bool                  `json:"all_passed"`
	ViewportResults    []ViewportCheckResult `json:"viewport_results"`
	Metrics            map[string]float64    `json:"metrics"`
	DetectedViolations []string              `json:"detected_violations"`
	FailureClass       string                `json:"failure_class"` // "NONE", "TOKEN_STYLE", "LAYOUT_CODE", "FATAL"
	VerifierRan        bool                  `json:"verifier_ran,omitempty"`
}

// IterationData tracks self-healing loop cycles and feedback in Stage 8
type IterationData struct {
	CurrentIteration   int       `json:"current_iteration"`
	MaxIterations      int       `json:"max_iterations"`
	History            []string  `json:"history"`
	CorrectiveFeedback []string  `json:"corrective_feedback"`
	FinalVerdict       string    `json:"final_verdict"` // "PASSED", "EXHAUSTED", "GATED"
	LoopTargetStage    StageName `json:"loop_target_stage,omitempty"`
}

// TaskContext carries mutable pipeline state, loaded registries, and artifacts across stages
type TaskContext struct {
	Ctx           context.Context
	Task          *TaskRequest
	Catalog       *resources.ResourceCatalog
	Graph         *resources.DesignResourceGraph
	Router        *router.Router
	ResearchCoord *research.ResearchCoordinator
	Verifier      VisualVerifier

	// DesignLab is the write-blocking gate. While it is PENDING, the implement
	// stage refuses to write any file a browser would render.
	DesignLab *verify.DesignLab

	// Stage-Specific Data Payloads
	Discovery      *DiscoveryData
	Classification *ClassificationData
	Research       *ResearchData
	Synthesis      *SynthesisData
	DesignSystem   *DesignSystemData
	Implementation *ImplementationData
	VisualQA       *VisualQAData
	Iteration      *IterationData

	// Execution History & Telemetry
	StageResults   map[StageName]*StageResult
	ArtifactPaths  map[string]string // artifact name -> absolute path
	StagesExecuted []StageName

	mu sync.RWMutex
}

// NewTaskContext initializes a clean TaskContext for a task
func NewTaskContext(
	ctx context.Context,
	task *TaskRequest,
	cat *resources.ResourceCatalog,
	graph *resources.DesignResourceGraph,
	rtr *router.Router,
	coord *research.ResearchCoordinator,
	verifier VisualVerifier,
) *TaskContext {
	if task.MaxIterations <= 0 {
		task.MaxIterations = 3
	}
	if task.WorkspaceRoot == "" {
		task.WorkspaceRoot = "."
	}

	return &TaskContext{
		Ctx:           ctx,
		Task:          task,
		Catalog:       cat,
		Graph:         graph,
		Router:        rtr,
		ResearchCoord: coord,
		Verifier:      verifier,
		StageResults:  make(map[StageName]*StageResult),
		ArtifactPaths: make(map[string]string),
		Iteration: &IterationData{
			CurrentIteration:   0,
			MaxIterations:      task.MaxIterations,
			History:            []string{},
			CorrectiveFeedback: []string{},
			FinalVerdict:       "PENDING",
		},
	}
}

// SetStageResult safely records a stage result and appends to execution order
func (c *TaskContext) SetStageResult(name StageName, res *StageResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.StageResults[name] = res
	c.StagesExecuted = append(c.StagesExecuted, name)
}

// GetStageResult retrieves the result for a given stage
func (c *TaskContext) GetStageResult(name StageName) *StageResult {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.StageResults[name]
}

// AddArtifact registers a produced artifact path
func (c *TaskContext) AddArtifact(name, path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ArtifactPaths == nil {
		c.ArtifactPaths = make(map[string]string)
	}
	c.ArtifactPaths[name] = path
}

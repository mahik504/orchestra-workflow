package engine

import (
	"fmt"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/classifier"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func attachMission(ctx *TaskContext) {
	if ctx == nil || ctx.Task == nil {
		return
	}
	if ctx.Mission == nil {
		id := ctx.Task.ID
		if id == "" {
			id = "mission"
		}
		ctx.Mission = mission.FromSeed(mission.Seed{
			ID:             id,
			Objective:      ctx.Task.RawRequest,
			WorkspaceRoot:  ctx.Task.WorkspaceRoot,
			DeclaredType:   ctx.Task.Type,
			SkipVisualGate: ctx.Task.SkipVisualGate,
			MaxIterations:  ctx.Task.MaxIterations,
			Tags:           ctx.Task.Tags,
		})
	}
	if ctx.RunState == nil {
		ctx.RunState = mission.NewRunState(ctx.Mission)
	}
}

func refineMissionFromClassification(ctx *TaskContext) {
	if ctx == nil || ctx.Mission == nil || ctx.Classification == nil {
		return
	}
	var caps []string
	for _, r := range ctx.Classification.ResolvedRoutes {
		caps = append(caps, r.CapabilityID)
	}
	mission.Refine(ctx.Mission, mission.RefineInput{
		QualityBar:        ctx.Classification.QualityBar,
		Capabilities:      caps,
		RequiresVisual:    ctx.Classification.RequiresVisual,
		RequiresSecurity:  ctx.Classification.RequiresSecurity,
		RequiresHumanGate: ctx.Classification.RequiresHumanGate,
		GateReason:        ctx.Classification.GateReason,
		ResearchDepth:     ctx.Classification.ResearchDepth,
	})
	if ctx.RunState != nil {
		ctx.RunState.Budget = ctx.Mission.Budget
	}
}

func selectCapabilities(ctx *TaskContext) error {
	if ctx == nil || ctx.Mission == nil {
		return nil
	}
	d, err := capability.Select(capabilityInput(ctx))
	ctx.Capability = &d
	if ctx.Task != nil {
		_, _ = capability.WriteDecision(ctx.Task.WorkspaceRoot, ctx.Mission.ID, d)
	}
	if err != nil {
		if capability.IsHumanStop(err) {
			return ErrHumanGateRequired
		}
		return err
	}
	return nil
}

func capabilityInput(ctx *TaskContext) capability.Input {
	in := capability.Input{Mission: ctx.Mission}
	if ctx.Task != nil {
		in.Resources = append(in.Resources, ctx.Task.SuggestedResources...)
		if ctx.Task.Parameters != nil {
			if v, ok := ctx.Task.Parameters["playwright_available"].(bool); ok {
				in.PlaywrightAvailable = &v
			}
			if hosts, ok := ctx.Task.Parameters["capability_hosts"].([]capability.HostSnapshot); ok {
				in.Hosts = hosts
			}
			if mcp, ok := ctx.Task.Parameters["capability_mcp"].([]capability.MCPCandidate); ok {
				in.MCP = mcp
			}
			if models, ok := ctx.Task.Parameters["discovered_models"].([]string); ok {
				in.DiscoveredModels = models
			}
			if skills, ok := ctx.Task.Parameters["capability_skills"].([]string); ok {
				in.Skills = append(in.Skills, skills...)
			}
			if shot, ok := ctx.Task.Parameters["screenshot_eval_available"].(bool); ok {
				in.ScreenshotEvalAvailable = &shot
			}
		}
	}
	if ctx.Mission != nil {
		in.Skills = append(in.Skills, ctx.Mission.Capabilities...)
	}
	if ctx.Classification != nil {
		in.Resources = append(in.Resources, ctx.Classification.OverlayActivations...)
		for _, r := range ctx.Classification.ResolvedRoutes {
			in.Skills = append(in.Skills, r.CapabilityID)
		}
		if ctx.Classification.Brief != nil {
			brief := ctx.Classification.Brief
			in.Skills = append(in.Skills, brief.Skills...)
			in.Trusted = append(in.Trusted, brief.Skills...)
			if ctx.Graph != nil && brief.CapabilityID != "" {
				if cap := ctx.Graph.Capabilities[brief.CapabilityID]; cap != nil {
					mcp := classifier.ResolveMCP(brief.CapabilityID, cap, ctx.Catalog, in.MCP, brief.RawRequest)
					for _, c := range mcp {
						in.Trusted = append(in.Trusted, c.ID)
					}
					in.MCP = mergeMCP(mcp, in.MCP)
				}
			}
		}
	}
	if in.PlaywrightAvailable == nil && ctx.Verifier != nil {
		if _, ok := ctx.Verifier.(*PlaywrightVerifier); !ok {
			avail := true
			in.PlaywrightAvailable = &avail
		}
	}
	if len(in.Hosts) == 0 {
		in.Hosts = []capability.HostSnapshot{
			{ID: "cursor", Available: true},
			{ID: "claude", Available: true},
		}
	}
	return in
}

func mergeMCP(route, extra []capability.MCPCandidate) []capability.MCPCandidate {
	seen := map[string]bool{}
	var out []capability.MCPCandidate
	for _, c := range route {
		id := strings.TrimSpace(c.ID)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, c)
	}
	for _, c := range extra {
		id := strings.TrimSpace(c.ID)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, c)
	}
	return out
}

func compileBuilderCapsule(ctx *TaskContext) {
	if ctx == nil || ctx.Mission == nil {
		return
	}
	if ctx.Capability != nil && !ctx.Capability.Has(capability.IDCodingWorker) {
		return
	}
	var files []string
	if ctx.DesignSystem != nil && ctx.DesignSystem.DesignMDPath != "" {
		files = append(files, ctx.DesignSystem.DesignMDPath)
	}
	if ctx.Discovery != nil && ctx.Discovery.PackageJSONPath != "" {
		files = append(files, ctx.Discovery.PackageJSONPath)
	}
	files = append(files, ctx.Mission.AffectedSurfaces...)
	var previous *mission.AttemptRecord
	var known []mission.FailureRecord
	if ctx.RunState != nil && len(ctx.RunState.FailureHistory) > 0 {
		last := ctx.RunState.FailureHistory[len(ctx.RunState.FailureHistory)-1]
		previous = &last
		for _, h := range ctx.RunState.FailureHistory {
			known = append(known, mission.FailureRecord{Class: h.FailureClass, Note: h.PreviousFailure, Stage: ctx.RunState.CurrentStage})
		}
	}
	contract := ""
	if ctx.DesignSystem != nil {
		contract = ctx.DesignSystem.DesignMDPath
	}
	skills := ctx.Mission.Capabilities
	var resources []string
	var tools []string
	if ctx.Capability != nil {
		skills = ctx.Capability.IDs(capability.CatSkill)
		resources = ctx.Capability.IDs(capability.CatResource)
		tools = ctx.Capability.IDs(capability.CatMCPTool)
	}
	cap, report, err := mission.Compile(mission.CompileRequest{
		Mission:            ctx.Mission,
		Role:               mission.RoleBuilder,
		Files:              files,
		Skills:             skills,
		Resources:          resources,
		ToolsAllowed:       tools,
		ExpectedOutput:     "scoped implementation for the mission objective",
		Previous:           previous,
		KnownFailures:      known,
		DesignContractPath: contract,
		UnrelatedFiles:     []string{"AGENTS.md", "README.md", "ARCHITECTURE.md"},
	})
	if err != nil {
		return
	}
	ctx.Capsule = cap
	if ctx.RunState != nil {
		ctx.RunState.CapsuleScope = report
		ctx.RunState.CurrentWorker = string(mission.RoleBuilder)
	}
}

func assertCodingWorkerConsumed(ctx *TaskContext) error {
	if ctx == nil || ctx.Capability == nil {
		return nil
	}
	if !ctx.Capability.Has(capability.IDCodingWorker) {
		return fmt.Errorf("required coding.worker omitted from capability decision")
	}
	if ctx.Capsule != nil && ctx.Capsule.Role != mission.RoleBuilder {
		return fmt.Errorf("required coding.worker was not consumed by the builder capsule")
	}
	return nil
}

func persistMission(ctx *TaskContext, note string) error {
	if ctx == nil || ctx.Task == nil || ctx.Mission == nil {
		return nil
	}
	if ctx.RunState == nil {
		ctx.RunState = mission.NewRunState(ctx.Mission)
	}
	ctx.Mission.Run.WorkspaceRoot = ctx.Task.WorkspaceRoot
	if len(ctx.StagesExecuted) > 0 {
		ctx.RunState.CurrentStage = string(ctx.StagesExecuted[len(ctx.StagesExecuted)-1])
		ctx.Mission.Run.CurrentStage = ctx.RunState.CurrentStage
	}
	ctx.Mission.Budget = ctx.RunState.Budget
	ctx.Mission.Run.Attempt = ctx.RunState.Attempt
	ctx.Mission.Status = ctx.RunState.Status
	if ctx.RunState.Status == "" {
		ctx.RunState.Status = mission.StatusRunning
		ctx.Mission.Status = mission.StatusRunning
	}
	_, _ = mission.WriteMissionFile(ctx.Task.WorkspaceRoot, ctx.Mission)
	_, _ = mission.WriteRunState(ctx.Task.WorkspaceRoot, ctx.RunState)
	_ = mission.WriteProgress(ctx.Task.WorkspaceRoot, ctx.RunState, ctx.Mission.Intensity, note)
	if ctx.Capability != nil {
		if _, err := capability.WriteDecision(ctx.Task.WorkspaceRoot, ctx.Mission.ID, *ctx.Capability); err != nil {
			return err
		}
	}
	return nil
}

func ensureCapabilityAuthoritative(ctx *TaskContext) error {
	if ctx == nil || ctx.Capability == nil || ctx.Task == nil || ctx.Mission == nil {
		return nil
	}
	onDisk, err := capability.ReadDecision(ctx.Task.WorkspaceRoot, ctx.Mission.ID)
	if err != nil {
		return nil
	}
	if !capability.AuthorityEqual(*onDisk, *ctx.Capability) {
		return capability.ErrDecisionMutated
	}
	return nil
}

func recoverStageError(ctx *TaskContext, stage StageName, err error) (bool, error) {
	if ctx == nil || ctx.Mission == nil || ctx.RunState == nil || err == nil {
		return false, err
	}
	fail := mission.Classify(mission.ClassifyInput{
		Err:   err,
		Stage: string(stage),
	})
	dec := mission.Decide(mission.DecideInput{
		Failure: fail,
		Budget:  ctx.RunState.Budget,
		History: ctx.RunState.FailureHistory,
		Attempt: ctx.RunState.Attempt,
	})
	state := mission.StatusFailed
	if dec.Retry {
		state = mission.StatusRetrying
	} else if dec.Human {
		state = mission.StatusWaitingHuman
		question := fail.What
		if ctx.Capability != nil && ctx.Capability.HumanDetail != "" {
			question = ctx.Capability.HumanDetail
		}
		ctx.RunState.PendingHuman = &mission.HumanDecision{
			Reason:   dec.Reason,
			Question: question,
			Blocking: true,
		}
		ctx.RunState.Status = state
	} else {
		ctx.RunState.Status = state
	}
	ctx.RunState.FailureHistory = mission.RecordAttempt(ctx.RunState.FailureHistory, ctx.RunState.Attempt, fail, dec, state)
	ctx.RunState.NextAction = dec.NextAction
	if dec.AdjustedBudget != nil {
		ctx.RunState.Budget = *dec.AdjustedBudget
		ctx.Mission.Budget = *dec.AdjustedBudget
	}
	persistMission(ctx, "failure "+string(fail.Class))
	if dec.Retry {
		ctx.RunState.Attempt++
		ctx.Mission.Run.Attempt = ctx.RunState.Attempt
		ctx.RunState.Status = mission.StatusRetrying
		return true, nil
	}
	return false, err
}

func recordIterateFailure(ctx *TaskContext, message, visualClass string) {
	if ctx == nil || ctx.RunState == nil {
		return
	}
	fail := mission.Classify(mission.ClassifyInput{
		Message: message,
		Stage:   string(StageVisualQA),
		Signal:  string(mission.ClassTest),
	})
	if visualClass == FailureClassFatal {
		sig := string(mission.ClassHuman)
		low := strings.ToLower(message)
		if strings.Contains(low, "playwright") || strings.Contains(low, "environment") || strings.Contains(low, "running application") {
			sig = string(mission.ClassEnvironment)
		}
		fail = mission.Classify(mission.ClassifyInput{
			Message: message,
			Stage:   string(StageVisualQA),
			Signal:  sig,
		})
	}
	dec := mission.Decide(mission.DecideInput{
		Failure: fail,
		Budget:  ctx.RunState.Budget,
		History: ctx.RunState.FailureHistory,
		Attempt: ctx.RunState.Attempt,
	})
	exhausted := ctx.Iteration != nil && ctx.Iteration.CurrentIteration >= ctx.Iteration.MaxIterations
	state := mission.StatusRetrying
	if exhausted {
		state = mission.StatusExhausted
		if dec.Human || fail.HumanRequired {
			state = mission.StatusWaitingHuman
			ctx.RunState.PendingHuman = &mission.HumanDecision{Reason: dec.Reason, Blocking: true}
		}
	}
	ctx.RunState.FailureHistory = mission.RecordAttempt(ctx.RunState.FailureHistory, ctx.RunState.Attempt, fail, dec, state)
	ctx.RunState.NextAction = dec.NextAction
}

func bindFlagshipExploration(ctx *TaskContext) error {
	if ctx == nil {
		return nil
	}
	if ctx.Capability != nil && !ctx.Capability.Has(capability.IDDesignWorker) {
		ctx.Exploration = nil
		return nil
	}
	class := designClassOf(ctx)
	dirs := directionsFromTask(ctx)
	if class != visual.ClassFlagship {
		if err := visual.ValidateExploration(class, dirs); err != nil {
			return err
		}
		ctx.Exploration = &visual.Exploration{Class: class}
		return nil
	}
	if len(dirs) == 0 {
		dirs = visual.DefaultFlagshipSet(objectiveOf(ctx))
	}
	if err := visual.ValidateExploration(class, dirs); err != nil {
		return err
	}
	selectedID := selectedDirectionID(ctx)
	skip := ctx.Task != nil && (ctx.Task.SkipVisualGate || ctx.Task.DryRun)
	if selectedID == "" {
		if skip {
			selectedID = dirs[0].ID
		} else {
			ctx.Exploration = &visual.Exploration{Class: class, Candidates: dirs}
			return ErrHumanGateRequired
		}
	}
	sel, err := visual.SelectOne(dirs, selectedID)
	if err != nil {
		return err
	}
	sealed, err := visual.Seal(*sel)
	if err != nil {
		return err
	}
	ctx.Exploration = &visual.Exploration{
		Class:      class,
		Candidates: dirs,
		SelectedID: selectedID,
		Sealed:     sealed,
		Waived:     skip && selectedDirectionID(ctx) == "",
	}
	return nil
}

func directionsFromTask(ctx *TaskContext) []visual.Direction {
	if ctx == nil || ctx.Task == nil || ctx.Task.Parameters == nil {
		return nil
	}
	raw, ok := ctx.Task.Parameters["directions"]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []visual.Direction:
		return v
	case []any:
		var out []visual.Direction
		for _, item := range v {
			if d, ok := item.(visual.Direction); ok {
				out = append(out, d)
			}
		}
		return out
	default:
		return nil
	}
}

func selectedDirectionID(ctx *TaskContext) string {
	if ctx == nil || ctx.Task == nil || ctx.Task.Parameters == nil {
		return ""
	}
	if s, ok := ctx.Task.Parameters["selected_direction_id"].(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func compileJudgeCapsule(ctx *TaskContext, qa *VisualQAData) {
	if ctx == nil || ctx.Mission == nil {
		return
	}
	if ctx.Capability != nil && !ctx.Capability.Has(capability.IDJudgeIndependent) {
		return
	}
	var files []string
	if qa != nil {
		for _, ev := range qa.Evidence {
			if ev.ScreenshotPath != "" {
				files = append(files, ev.ScreenshotPath)
			}
		}
	}
	if ctx.DesignSystem != nil && ctx.DesignSystem.DesignMDPath != "" {
		files = append(files, ctx.DesignSystem.DesignMDPath)
	}
	cap, report, err := mission.Compile(mission.CompileRequest{
		Mission:             ctx.Mission,
		Role:                mission.RoleJudge,
		Files:               files,
		Skills:              []string{"ow-reference-forensics"},
		AllSkills:           []string{"ow-reference-forensics", "hallmark", "taste-design", "impeccable", "gsap-core"},
		ExpectedOutput:      "structured Layer B/judge findings against Layer A captures",
		ToolsAllowed:        []string{"read_file", "write_file", "apply_patch", "shell"},
		UnrelatedFiles:      []string{"AGENTS.md", "README.md", "ARCHITECTURE.md"},
		ConversationHistory: []string{"full conversation must not enter the judge capsule"},
		MCPCatalog:          []string{"playwright-mcp", "cursor-browser", "stitch"},
		WholeProjectContext: []byte(strings.Repeat("WHOLE PROJECT DUMP\n", 4000)),
	})
	if err != nil {
		return
	}
	ctx.JudgeCapsule = cap
	ctx.JudgeScope = report
}

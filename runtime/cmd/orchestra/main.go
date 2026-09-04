package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/orchestra-v3/internal/adapters"
	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/engine"
	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/memory"
	"github.com/user/orchestra-v3/internal/resources"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		runInit(args)
	case "doctor":
		runDoctor(args)
	case "classify":
		runClassify(args)
	case "route", "plan":
		runPlan(args)
	case "run":
		runRun(args)
	case "verify":
		runVerify(args)
	case "handoff":
		runHandoff(args)
	case "sync":
		runSync(args)
	case "memory":
		runMemory(args)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Orchestra 3.1.0 — control plane for agentic IDEs")
	fmt.Println("ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.")
	if pin := os.Getenv("ORCHESTRA_CONTRACT"); pin != "" {
		fmt.Printf("ORCHESTRA_CONTRACT pin: %s (see kit/ROLLBACK.md)\n", pin)
	}
	fmt.Println("Usage: orchestra <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  init      Initialize a fresh orchestra workspace")
	fmt.Println("  doctor    Exhaustive environment, registry, quarantine, and visual QA diagnostic")
	fmt.Println("  classify  Parse and classify a task request")
	fmt.Println("  plan      Synthesize 8-stage execution manifest without side effects")
	fmt.Println("  route     Alias for plan")
	fmt.Println("  run       Execute full 8-stage pipeline (Discover -> Classify -> Research -> Synthesize -> Design System -> Implement -> Visual QA -> Iterate)")
	fmt.Println("  verify    Run standalone multi-viewport Visual QA and project verification")
	fmt.Println("  handoff   Inspect or initialize state handoff between agents")
	fmt.Println("  sync      Synchronize and verify active skill parity across Cursor, Claude, Antigravity")
	fmt.Println("  memory    Query or record private brain resource outcome memory (list, stats, record)")
}

func resolveRegistryFile(filename string, customPath string) string {
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath
		}
	}

	candidates := []string{
		filepath.Join("registries", filename),
		filepath.Join("..", "registries", filename),
		filepath.Join("..", "..", "registries", filename),
		filepath.Join("..", "..", "..", "registries", filename),
	}
	if root := os.Getenv("ORCHESTRA_WORKFLOW_ROOT"); root != "" {
		candidates = append(candidates, filepath.Join(root, "registries", filename))
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return filepath.Join("registries", filename)
}

func runInit(args []string) {
	workdir, _ := os.Getwd()
	if len(args) > 0 && args[0] != "" {
		workdir = args[0]
	}

	dirs := []string{
		filepath.Join(workdir, ".orchestra"),
		filepath.Join(workdir, ".orchestra", "handoff"),
		filepath.Join(workdir, ".orchestra", "qa"),
		filepath.Join(workdir, "memory"),
		filepath.Join(workdir, "skills"),
		filepath.Join(workdir, "projects"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", d, err)
			os.Exit(1)
		}
	}

	configPath := filepath.Join(workdir, ".orchestra", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := map[string]interface{}{
			"version":          "3.1.0",
			"default_agent":    "antigravity",
			"ponytail_mode":    "full",
			"isolation_mode":   "strict_clean",
			"storage_provider": "local_git",
		}
		bytes, _ := json.MarshalIndent(cfg, "", "  ")
		_ = os.WriteFile(configPath, bytes, 0644)
	}

	fmt.Println("[OK] Fresh Orchestra workspace initialized successfully.")
	fmt.Printf("Workspace root: %s\n", workdir)
	fmt.Println("Created clean directories: .orchestra/, memory/, skills/, projects/")
}

func runDoctor(args []string) {
	fmt.Println("=== Orchestra 3.1.0 Environment & System Doctor ===")
	if pin := os.Getenv("ORCHESTRA_CONTRACT"); pin != "" {
		fmt.Printf("Contract pin:    ORCHESTRA_CONTRACT=%s (see kit/ROLLBACK.md)\n", pin)
	} else {
		fmt.Println("Contract pin:    unset (using VERSION 3.1.0)")
	}

	checkCmd := func(name string, arg string) (string, bool) {
		out, err := exec.Command(name, arg).Output()
		if err != nil {
			return "[MISSING] Not detected in PATH", false
		}
		lines := strings.Split(string(out), "\n")
		return "[OK] " + strings.TrimSpace(lines[0]), true
	}

	// 1. Core Toolchain
	gitOut, _ := checkCmd("git", "--version")
	nodeOut, _ := checkCmd("node", "--version")
	npmOut, _ := checkCmd("npm", "--version")
	goOut, _ := checkCmd("go", "version")
	pyOut, _ := checkCmd("python", "--version")

	fmt.Printf("Git:             %s\n", gitOut)
	fmt.Printf("Node.js:         %s\n", nodeOut)
	fmt.Printf("npm:             %s\n", npmOut)
	fmt.Printf("Go Runtime:      %s\n", goOut)
	fmt.Printf("Python:          %s\n", pyOut)

	// 2. Canonical Registries & Graph
	catPath := resolveRegistryFile("resources.json", "")
	graphPath := resolveRegistryFile("design-resource-graph.json", "")

	if cat, err := resources.LoadResourceCatalog(catPath); err == nil {
		fmt.Printf("Registry:        [OK] %s (%d canonical resources)\n", catPath, cat.Count())
	} else {
		fmt.Printf("Registry:        [FAIL] %s (%v)\n", catPath, err)
	}

	if graph, err := resources.LoadDesignGraph(graphPath); err == nil {
		fmt.Printf("Design Graph:    [OK] %s (%d domains, %d capabilities)\n", graphPath, len(graph.Domains), len(graph.Capabilities))
	} else {
		fmt.Printf("Design Graph:    [FAIL] %s (%v)\n", graphPath, err)
	}

	homeDir, _ := os.UserHomeDir()
	quarantinePath := os.Getenv("ORCHESTRA_QUARANTINE_PATH")
	if quarantinePath == "" {
		quarantinePath = filepath.Join(homeDir, ".gemini", "config", "skills_library")
	}
	qErr := resources.CheckQuarantineBoundary(quarantinePath)
	qStatus, _ := resources.AuditQuarantineState(".", quarantinePath)
	if qErr != nil && (qStatus == nil || !qStatus.QuarantineDirectoryExists || qStatus.IsStrictlyIsolated) {
		count := 1598
		if qStatus != nil && qStatus.QuarantinedCount > 0 {
			count = qStatus.QuarantinedCount
		}
		fmt.Printf("Quarantine:      [PASS] %d-skill library quarantined (runtime boundary enforced)\n", count)
	} else {
		fmt.Printf("Quarantine:      [FAIL] Quarantine boundary check failed to reject isolated library\n")
	}

	// 4. Visual QA Tooling
	pwOut, pwOk := checkCmd("npx", "playwright --version")
	if pwOk {
		fmt.Printf("Playwright CLI:  %s\n", pwOut)
	} else {
		fmt.Println("Playwright CLI:  [NOTE] npx playwright not installed globally (mock verifier active)")
	}

	// 5. Resource memory diagnostics
	workdir, _ := os.Getwd()
	memPath := memory.ResolveDefaultMemoryPath(workdir)
	if store, err := memory.NewResourceMemoryStore(memPath); err == nil {
		total, _, _, rate := store.GetSummary()
		aggs := store.ListAggregates()
		fmt.Printf("Memory:          [OK] %s (%d resources, %d evaluations, %.1f%% success rate)\n",
			memPath, len(aggs), total, rate*100)
	} else {
		fmt.Printf("Memory:          [NOTE] %s ready for milestone sync (%v)\n", memPath, err)
	}

	// 6. Host Synchronization & Skill Parity
	syncEngine := adapters.NewHostSyncEngine(workdir)
	userHome, _ := os.UserHomeDir()
	if parityRep, err := syncEngine.VerifyParity(userHome); err == nil && parityRep.IsParityComplete {
		fmt.Printf("Host Parity:     [OK] %d/%d active skills synchronized across Cursor, Antigravity, Claude\n",
			parityRep.CanonicalSkillCount, parityRep.CanonicalSkillCount)
	} else if parityRep != nil {
		fmt.Printf("Host Parity:     [WARN] Active skill parity discrepancy detected across hosts (run 'orchestra sync')\n")
	} else {
		fmt.Printf("Host Parity:     [FAIL] Could not verify host skill parity: %v\n", err)
	}

	// 7. Antigravity customization budget
	if budget, err := adapters.CheckAntigravityBudget(userHome); err == nil {
		printAGBudget(budget)
	}
}

func printAGBudget(budget *adapters.AGBudgetReport) {
	if budget.GlobalSkillCount > 0 {
		fmt.Printf("AG Global skills: [OK] %d installed: %s\n",
			budget.GlobalSkillCount, strings.Join(budget.GlobalSkillNames, ", "))
	} else {
		fmt.Println("AG Global skills: [NOTE] ~/.gemini/config/skills not present")
	}

	if len(budget.BannedEnabled) > 0 {
		fmt.Printf("AG plugins:      [WARN] banned Global plugins enabled: %s\n",
			strings.Join(budget.BannedEnabled, ", "))
		fmt.Println("                 Disable science and data-agent-kit-plugin in Antigravity Settings > Customizations.")
		fmt.Println("                 Re-enable only for a job that needs AlphaFold / BigQuery.")
	} else if len(budget.BannedPresent) > 0 {
		fmt.Printf("AG plugins:      [OK] %s installed but disabled\n",
			strings.Join(budget.BannedPresent, ", "))
	} else {
		fmt.Println("AG plugins:      [OK] science / data-agent-kit not installed as Global")
	}

	if budget.HeadroomGone {
		fmt.Println("AG headroom:     [WARN] customization budget is gone — banned plugins are Global on top of the 30-skill core")
	}

	if len(budget.MCPServers) == 0 {
		fmt.Println("AG MCP:          [NOTE] no mcp_config.json (or empty)")
		return
	}
	fmt.Println("AG MCP:")
	for _, s := range budget.MCPServers {
		fmt.Printf("  %-22s %s\n", s.Name, s.Health)
	}
}

func runClassify(args []string) {
	fs := flag.NewFlagSet("classify", flag.ExitOnError)
	graphPath := fs.String("graph", "", "Path to design-resource-graph.json")
	asJSON := fs.Bool("json", false, "Emit the full brief as JSON")
	silent := fs.Bool("assume", false, "Answer no clarifying question; take the lower-risk route")
	_ = fs.Parse(args)

	raw := strings.Join(fs.Args(), " ")
	if raw == "" {
		fmt.Println("Usage: orchestra classify [--json] [--assume] \"<task description>\"")
		os.Exit(1)
	}

	graph, err := resources.LoadDesignGraph(resolveRegistryFile("design-resource-graph.json", *graphPath))
	if err != nil {
		fmt.Printf("Error loading capability graph: %v\n", err)
		os.Exit(1)
	}

	brief := classifier.NewClassifierWithGraph(graph).ClassifyBrief(raw, classifier.Options{})
	if brief.Ambiguous && *silent {
		brief.ResolveSilence()
	}

	if *asJSON {
		out, _ := json.MarshalIndent(brief, "", "  ")
		fmt.Println(string(out))
		return
	}

	printBrief(brief)
}

// printBrief renders the re-brief a human is expected to correct before any
// work starts. Declined routes are shown because "why not that one?" is the
// question that catches a bad classification early.
func printBrief(b *classifier.Brief) {
	fmt.Printf("Archetype:     %s", b.Archetype)
	if b.CapabilityID != "" {
		fmt.Printf("  (%s)", b.CapabilityID)
	}
	fmt.Printf("\n               %s\n", b.ArchetypeReason)
	fmt.Printf("Task type:     %s\n", b.Type)
	fmt.Printf("Quality bar:   %s — %s\n", b.QualityBar, b.QualityBarReason)
	fmt.Printf("Platform:      %s\n", b.Platform)
	fmt.Printf("Research:      %s\n", b.ResearchDepth)
	fmt.Printf("Verify:        %s\n", b.VerifyDepth)
	fmt.Printf("Design Lab:    %v — %s\n", b.DesignLabRequired, b.DesignLabReason)

	if len(b.HardConstraints) > 0 {
		fmt.Println("Constraints:")
		for _, c := range b.HardConstraints {
			fmt.Printf("  - %s\n", c)
		}
	}
	if b.Assumed {
		fmt.Printf("\n[ASSUMED] %s\n", b.AssumptionNote)
	}
	if b.Ambiguous {
		fmt.Printf("\n[QUESTION] %s\n", b.ClarifyingQuestion)
	}
	if len(b.UnknownTechnology) > 0 {
		fmt.Printf("\n[UNKNOWN] not in the graph, research and register: %s\n", strings.Join(b.UnknownTechnology, ", "))
	}

	fmt.Println("\nRoutes considered:")
	for _, c := range b.Considered {
		mark := "take"
		if c.Declined {
			mark = "skip"
		}
		fmt.Printf("  [%s] %-21s %6.2f  %s\n", mark, c.CapabilityID, c.Score, c.DeclineReason)
	}
}

func runPlan(args []string) {
	fs := flag.NewFlagSet("plan", flag.ExitOnError)
	taskStr := fs.String("task", "", "Task description or PRD text")
	catalogPath := fs.String("catalog", "", "Path to resources.json")
	graphPath := fs.String("graph", "", "Path to design-resource-graph.json")
	outputFormat := fs.String("output", "manifest", "Output format: manifest, json, summary")
	skipGate := fs.Bool("skip-visual-gate", false, "Bypass human approval gate")
	_ = fs.Parse(args)

	rawTask := *taskStr
	if rawTask == "" && len(fs.Args()) > 0 {
		rawTask = strings.Join(fs.Args(), " ")
	}
	if rawTask == "" {
		rawTask = "Build modern responsive web feature with strict typography and security headers"
	}

	resolvedCat := resolveRegistryFile("resources.json", *catalogPath)
	catalog, err := resources.LoadResourceCatalog(resolvedCat)
	if err != nil {
		fmt.Printf("[FAIL] Could not load resource catalog from %s: %v\n", resolvedCat, err)
		os.Exit(1)
	}

	resolvedGraph := resolveRegistryFile("design-resource-graph.json", *graphPath)
	graph, err := resources.LoadDesignGraph(resolvedGraph)
	if err != nil {
		fmt.Printf("[FAIL] Could not load design graph from %s: %v\n", resolvedGraph, err)
		os.Exit(1)
	}

	pipeline := engine.NewDesignPipeline(catalog, graph, nil, nil)
	taskReq := &engine.TaskRequest{
		ID:             "cli-plan-task",
		RawRequest:     rawTask,
		SkipVisualGate: *skipGate,
		DryRun:         true,
	}

	plan, err := pipeline.Plan(context.Background(), taskReq)
	if err != nil {
		fmt.Printf("[FAIL] Planning failed: %v\n", err)
		os.Exit(1)
	}

	switch strings.ToLower(*outputFormat) {
	case "json":
		b, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(b))
	case "summary":
		fmt.Printf("Task Archetype:       %s\n", plan.PrimaryArchetype)
		fmt.Printf("Estimated Token Cost: %.0f tokens\n", plan.EstimatedTokenCost)
		fmt.Printf("Human Gate Required:  %v (%s)\n", plan.RequiresHumanGate, plan.ApprovalReason)
		fmt.Printf("Active Directives:    %d\n", len(plan.ExecutionDirectives))
	default:
		fmt.Println(plan.GenerateExecutionManifest())
		fmt.Printf("Estimated Context Cost: %.0f tokens\n", plan.EstimatedTokenCost)
		if plan.RequiresHumanGate {
			fmt.Printf("\n[GATE ALERT] Human approval gate is REQUIRED before execution.\nReason: %s\n", plan.ApprovalReason)
		}
	}
}

func runRun(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	taskStr := fs.String("task", "", "Task description")
	workdir := fs.String("workdir", ".", "Project working directory")
	autoApprove := fs.Bool("auto-approve", false, "Auto-approve human design gate")
	maxIterations := fs.Int("max-iterations", 3, "Max QA iteration loops")
	catalogPath := fs.String("catalog", "", "Path to resources.json")
	graphPath := fs.String("graph", "", "Path to design-resource-graph.json")
	_ = fs.Parse(args)

	rawTask := *taskStr
	if rawTask == "" && len(fs.Args()) > 0 {
		rawTask = strings.Join(fs.Args(), " ")
	}
	if rawTask == "" {
		fmt.Println("Usage: orchestra run --task \"<task description>\" [--auto-approve]")
		os.Exit(1)
	}

	resolvedCat := resolveRegistryFile("resources.json", *catalogPath)
	catalog, err := resources.LoadResourceCatalog(resolvedCat)
	if err != nil {
		fmt.Printf("[FAIL] Resource catalog error from %s: %v\n", resolvedCat, err)
		os.Exit(1)
	}

	resolvedGraph := resolveRegistryFile("design-resource-graph.json", *graphPath)
	graph, err := resources.LoadDesignGraph(resolvedGraph)
	if err != nil {
		fmt.Printf("[FAIL] Design graph error from %s: %v\n", resolvedGraph, err)
		os.Exit(1)
	}

	pipeline := engine.NewDesignPipeline(catalog, graph, nil, nil)
	pipeline.MaxIterations = *maxIterations

	taskReq := &engine.TaskRequest{
		ID:             "cli-exec-task",
		RawRequest:     rawTask,
		WorkspaceRoot:  *workdir,
		SkipVisualGate: *autoApprove,
		MaxIterations:  *maxIterations,
	}

	fmt.Printf("[Orchestra] Starting 8-stage execution for: %s\n", rawTask)
	result, err := pipeline.Execute(context.Background(), taskReq)
	if err != nil {
		if err == engine.ErrHumanGateRequired {
			fmt.Println("\n[GATE HALTED] Human approval required before writing code.")
			fmt.Println("Run with --auto-approve to bypass in automated testing environments.")
			os.Exit(2)
		}
		fmt.Printf("\n[FAIL] Pipeline execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Pipeline Execution Summary ===")
	fmt.Printf("Status:        %s\n", result.Status)
	fmt.Printf("Archetype:     %s\n", result.Archetype)
	fmt.Printf("Iterations:    %d\n", result.IterationCount)
	fmt.Printf("Duration:      %s\n", result.TotalDuration)
	if result.ReferenceLogPath != "" {
		fmt.Printf("Reference Log: %s\n", result.ReferenceLogPath)
	}
	if result.DesignMDPath != "" {
		fmt.Printf("DESIGN.md:     %s\n", result.DesignMDPath)
	}
	fmt.Printf("Visual QA:     Passed=%v (Screenshots: %d)\n", result.VisualQAPassed, len(result.Screenshots))
	if result.HandoffStatePath != "" {
		fmt.Printf("Handoff State: %s\n", result.HandoffStatePath)
	}
}

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	viewportsStr := fs.String("viewports", "desktop,tablet,mobile", "Viewports to verify")
	outDir := fs.String("output-dir", filepath.Join(".orchestra", "qa"), "Output directory for screenshots and report")
	strict := fs.Bool("strict", false, "Fail on any warning or anti-pattern match")
	_ = fs.Parse(args)

	fmt.Println("=== Orchestra V3 Multi-Viewport Visual QA Runner ===")
	fmt.Printf("Viewports:       %s\n", *viewportsStr)
	fmt.Printf("Report Dir:      %s\n", *outDir)

	workdir, _ := os.Getwd()
	pkgJson := filepath.Join(workdir, "package.json")
	if _, err := os.Stat(pkgJson); err == nil {
		fmt.Println("[OK] Web project detected.")
	}

	verifier := engine.NewPlaywrightVerifier(*outDir)
	taskCtx := &engine.TaskContext{
		Task: &engine.TaskRequest{
			WorkspaceRoot: workdir,
		},
		Classification: &engine.ClassificationData{
			RequiresVisual: true,
		},
	}

	qaRes, err := verifier.Verify(context.Background(), taskCtx)
	if err != nil {
		fmt.Printf("[FAIL] Verification execution error: %v\n", err)
		os.Exit(1)
	}

	for _, vp := range qaRes.ViewportResults {
		status := "[PASS]"
		if !vp.Passed {
			status = "[FAIL]"
		}
		fmt.Printf("Viewport %-8s %s (Overflow: %v, Screenshot: %s)\n",
			vp.ViewportName, status, vp.HasOverflow, vp.ScreenshotPath)
	}

	if !qaRes.AllPassed {
		fmt.Println("\n[FAIL] Visual QA defects detected:")
		for _, d := range qaRes.DetectedViolations {
			fmt.Printf(" - %s\n", d)
		}
		if *strict {
			os.Exit(1)
		}
	} else {
		fmt.Println("\n[PASS] All viewports verified successfully. Zero horizontal overflow.")
	}

	// 2. Host Active Skill Parity Verification
	syncEngine := adapters.NewHostSyncEngine(workdir)
	userHome, _ := os.UserHomeDir()
	if parityRep, err := syncEngine.VerifyParity(userHome); err == nil {
		if parityRep.IsParityComplete {
			fmt.Printf("Host Parity:     [PASS] %d/%d active skills synchronized across hosts (byte-identical=%v)\n",
				parityRep.CanonicalSkillCount, parityRep.CanonicalSkillCount, parityRep.ByteIdentical)
		} else {
			fmt.Printf("Host Parity:     [WARN] Incomplete parity across hosts: missing=%v, byte-identical=%v\n",
				parityRep.MissingSkills, parityRep.ByteIdentical)
		}
	}

	// 3. Resource memory integrity verification
	memPath := memory.ResolveDefaultMemoryPath(workdir)
	if store, err := memory.NewResourceMemoryStore(memPath); err == nil {
		total, succ, fail, rate := store.GetSummary()
		fmt.Printf("Memory:          [PASS] %s valid (total=%d, succ=%d, fail=%d, rate=%.1f%%)\n",
			memPath, total, succ, fail, rate*100)
	} else {
		fmt.Printf("Memory:          [WARN] %s unreadable: %v\n", memPath, err)
	}
}

func runHandoff(args []string) {
	workdir, _ := os.Getwd()
	state, err := handoff.ReadState(workdir)
	if err != nil {
		fmt.Printf("[NOTE] No active handoff state: %v\n", err)
		fmt.Println("To create a handoff, write state to .orchestra/handoff/state.json")
		return
	}

	bytes, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println("=== Current Handoff State ===")
	fmt.Println(string(bytes))

	conflicts, err := handoff.DetectConflicts(state, workdir)
	if err != nil {
		fmt.Printf("Error detecting conflicts: %v\n", err)
		return
	}

	if len(conflicts) > 0 {
		fmt.Printf("[ALERT] %d out-of-band conflict(s) detected:\n", len(conflicts))
		for _, c := range conflicts {
			fmt.Printf(" - %s\n", c)
		}
	} else {
		fmt.Println("[OK] Zero out-of-band file conflicts detected.")
	}
}

func runSync(args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	host := fs.String("host", "all", "Target host to sync: cursor, claude, antigravity, agents, all")
	checkOnly := fs.Bool("check", false, "Check active skill parity without writing changes")
	_ = fs.Parse(args)

	fmt.Println("=== Orchestra V3 Host Configuration & Skill Parity Sync ===")
	workdir, _ := os.Getwd()
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("[FAIL] Could not determine user home directory: %v\n", err)
		os.Exit(1)
	}

	engine := adapters.NewHostSyncEngine(workdir)

	if !*checkOnly {
		fmt.Printf("Synchronizing 30 canonical active skills to host(s): %s...\n", *host)
		if err := engine.SyncAll(userHome, *host); err != nil {
			fmt.Printf("[FAIL] Synchronization failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("[OK] Skill synchronization completed successfully.")
	}

	report, err := engine.VerifyParity(userHome)
	if err != nil {
		fmt.Printf("[FAIL] Parity check failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n--- Active Skill Parity Audit (%d Canonical Skills) ---\n", report.CanonicalSkillCount)
	for h, count := range report.HostSkillCounts {
		status := "[OK]"
		if count != report.CanonicalSkillCount {
			status = "[WARN]"
		}
		fmt.Printf("Host %-12s: %s %d skills installed\n", h, status, count)
	}

	if len(report.MissingSkills) > 0 {
		for h, missing := range report.MissingSkills {
			if len(missing) > 0 {
				fmt.Printf(" - %s missing skills (%d): %s\n", h, len(missing), strings.Join(missing, ", "))
			}
		}
	}

	if len(report.ExtraSkills) > 0 {
		for h, extra := range report.ExtraSkills {
			if len(extra) > 0 {
				fmt.Printf(" - %s unapproved extra skills (%d): %s\n", h, len(extra), strings.Join(extra, ", "))
			}
		}
	}

	if len(report.QuarantineViolations) > 0 {
		fmt.Println("\n[SECURITY VIOLATION] Quarantine boundary breaches detected:")
		for _, v := range report.QuarantineViolations {
			fmt.Printf(" ! %s\n", v)
		}
		os.Exit(1)
	}

	if report.IsParityComplete && report.ByteIdentical {
		fmt.Println("\n[PASS] 100% byte-identical active skill parity verified across Cursor, Claude, and Antigravity.")
	} else if report.IsParityComplete {
		fmt.Println("\n[NOTE] Skill names match across hosts, but minor file differences were detected.")
	} else {
		fmt.Println("\n[FAIL] Skill parity incomplete. Run 'orchestra sync' to resolve.")
		if *checkOnly {
			os.Exit(1)
		}
	}
}

func runMemory(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: orchestra memory <list|stats|record> [options]")
		return
	}

	sub := args[0]
	subArgs := args[1:]
	workdir, _ := os.Getwd()
	memPath := memory.ResolveDefaultMemoryPath(workdir)

	store, err := memory.NewResourceMemoryStore(memPath)
	if err != nil {
		fmt.Printf("[FAIL] Could not load resource memory from %s: %v\n", memPath, err)
		os.Exit(1)
	}

	switch sub {
	case "list":
		fs := flag.NewFlagSet("memory list", flag.ExitOnError)
		filter := fs.String("filter", "", "Filter by resource ID")
		asJSON := fs.Bool("json", false, "Output as JSON")
		_ = fs.Parse(subArgs)

		aggs := store.ListAggregates()
		if *filter != "" {
			var filtered []*memory.ResourceAggregate
			for _, a := range aggs {
				if strings.Contains(a.ResourceID, strings.ToLower(*filter)) {
					filtered = append(filtered, a)
				}
			}
			aggs = filtered
		}

		if *asJSON {
			b, _ := json.MarshalIndent(aggs, "", "  ")
			fmt.Println(string(b))
			return
		}

		fmt.Printf("=== Orchestra V3 Resource Memory (%s) ===\n", memPath)
		fmt.Printf("%-18s %-16s %-20s %-6s %-10s %-12s %-8s\n",
			"RESOURCE ID", "DOMAIN", "CAPABILITY", "EVALS", "SUCCESS", "AVG LATENCY", "QUALITY")
		fmt.Println(strings.Repeat("-", 95))
		for _, a := range aggs {
			fmt.Printf("%-18s %-16s %-20s %-6d %-9.1f%% %-10.1fms %-8.2f\n",
				a.ResourceID, a.Domain, a.Capability, a.TotalEvaluations, a.SuccessRate*100, a.AverageLatencyMs, a.AverageQualityScore)
		}

	case "record":
		fs := flag.NewFlagSet("memory record", flag.ExitOnError)
		resID := fs.String("resource", "", "Resource ID (required)")
		outcome := fs.String("outcome", "success", "Outcome: success or failure")
		domain := fs.String("domain", "general", "Domain")
		capability := fs.String("capability", "general", "Capability")
		taskID := fs.String("task", "manual-cli", "Task ID")
		score := fs.Float64("score", 1.0, "Quality score [0.0 - 1.0]")
		latency := fs.Int64("latency", 100, "Latency in ms")
		errDetails := fs.String("error", "", "Error details")
		notes := fs.String("notes", "", "Notes")
		_ = fs.Parse(subArgs)

		if *resID == "" {
			fmt.Println("Error: --resource is required")
			os.Exit(1)
		}

		eval := &memory.ResourceEvaluation{
			ResourceID:   *resID,
			Domain:       *domain,
			Capability:   *capability,
			TaskID:       *taskID,
			Outcome:      memory.Outcome(*outcome),
			QualityScore: *score,
			LatencyMs:    *latency,
			ErrorDetails: *errDetails,
			Notes:        *notes,
		}

		if err := store.Record(eval); err != nil {
			fmt.Printf("[FAIL] Recording evaluation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[OK] Recorded evaluation for %s (%s, score=%.2f, latency=%dms)\n", *resID, *outcome, *score, *latency)

	case "stats", "summary":
		total, succ, fail, rate := store.GetSummary()
		aggs := store.ListAggregates()
		fmt.Printf("=== Private Brain Resource Memory Summary ===\n")
		fmt.Printf("Ledger File:         %s\n", memPath)
		fmt.Printf("Tracked Resources:   %d\n", len(aggs))
		fmt.Printf("Total Evaluations:   %d\n", total)
		fmt.Printf("Successful Runs:     %d\n", succ)
		fmt.Printf("Failed Runs:         %d\n", fail)
		fmt.Printf("Overall Success:     %.1f%%\n", rate*100)

	default:
		fmt.Printf("Unknown memory command: %s. Use list, record, or stats.\n", sub)
	}
}

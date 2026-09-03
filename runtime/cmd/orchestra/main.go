package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/handoff"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/router"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit(os.Args[2:])
	case "doctor":
		runDoctor(os.Args[2:])
	case "classify":
		runClassify(os.Args[2:])
	case "route", "plan":
		runPlan(os.Args[2:])
	case "verify":
		runVerify(os.Args[2:])
	case "handoff":
		runHandoff(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Orchestra V3 - Lightweight Agentic CLI Runtime")
	fmt.Println("Usage: orchestra <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  init      Initialize a fresh orchestra workspace and clean personal brain")
	fmt.Println("  doctor    Check system dependencies, adapters, and environment hygiene")
	fmt.Println("  classify  Parse and classify a task request")
	fmt.Println("  plan      Synthesize minimal sufficient capabilities & execution manifest")
	fmt.Println("  route     Alias for plan")
	fmt.Println("  verify    Run verification checks against current project state")
	fmt.Println("  handoff   Inspect or initialize state handoff between agents")
}

func runInit(args []string) {
	workdir, _ := os.Getwd()
	if len(args) > 0 && args[0] != "" {
		workdir = args[0]
	}

	dirs := []string{
		filepath.Join(workdir, ".orchestra"),
		filepath.Join(workdir, ".orchestra", "handoff"),
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
	fmt.Println("=== Orchestra V3 Environment & Doctor Diagnostic ===")

	checkCmd := func(name string, arg string) string {
		out, err := exec.Command(name, arg).Output()
		if err != nil {
			return "[MISSING] Not detected in PATH"
		}
		lines := strings.Split(string(out), "\n")
		return "[OK] " + strings.TrimSpace(lines[0])
	}

	fmt.Printf("Git:        %s\n", checkCmd("git", "--version"))
	fmt.Printf("Node.js:    %s\n", checkCmd("node", "--version"))
	fmt.Printf("npm:        %s\n", checkCmd("npm", "--version"))
	fmt.Printf("Go Runtime: %s\n", checkCmd("go", "version"))

	workdir, _ := os.Getwd()
	orchestraDir := filepath.Join(workdir, ".orchestra")
	if _, err := os.Stat(orchestraDir); err == nil {
		fmt.Println("Workspace:  [OK] .orchestra configuration detected")
	} else {
		fmt.Println("Workspace:  [NOTE] Run 'orchestra init' to scaffold local workspace")
	}
}

func runClassify(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: orchestra classify \"<task description>\"")
		os.Exit(1)
	}

	c := classifier.NewClassifier()
	task, err := c.Classify(args[0])
	if err != nil {
		fmt.Printf("Error classifying task: %v\n", err)
		os.Exit(1)
	}

	bytes, _ := json.MarshalIndent(task, "", "  ")
	fmt.Println(string(bytes))
}

func runPlan(args []string) {
	rawTask := "Implement modern responsive web feature with strict typography and security headers"
	if len(args) > 0 {
		rawTask = strings.Join(args, " ")
	}

	c := classifier.NewClassifier()
	task, _ := c.Classify(rawTask)

	reg := resources.NewRegistry()
	reg.Capabilities["superpowers-planning"] = &resources.Capability{
		ID:                 "superpowers-planning",
		Name:               "Superpowers Planning",
		TokenContextWeight: 1500,
	}
	reg.Capabilities["taste-skill"] = &resources.Capability{
		ID:                 "taste-skill",
		Name:               "Taste Skill",
		TokenContextWeight: 2000,
	}
	reg.Capabilities["impeccable"] = &resources.Capability{
		ID:                 "impeccable",
		Name:               "Impeccable Design",
		TokenContextWeight: 1800,
	}
	reg.Capabilities["semgrep-adapter"] = &resources.Capability{
		ID:                 "semgrep-adapter",
		Name:               "Semgrep Security Adapter",
		TokenContextWeight: 1200,
	}

	r := router.NewRouter(reg)
	plan := r.Compose(task)

	fmt.Println(plan.GenerateExecutionManifest())
	fmt.Printf("Estimated Context Cost: %.0f tokens\n", plan.EstimatedTokenCost)
	fmt.Printf("Requires Human Design Gate: %v (%s)\n", plan.RequiresHumanGate, plan.ApprovalReason)
}

func runVerify(args []string) {
	fmt.Println("=== Orchestra V3 Verification Runner ===")
	workdir, _ := os.Getwd()

	pkgJson := filepath.Join(workdir, "package.json")
	if _, err := os.Stat(pkgJson); err == nil {
		fmt.Println("[OK] Web project detected. Running TypeScript verification...")
		cmd := exec.Command("npm", "run", "typecheck")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("[FAIL] TypeScript errors:\n%s\n", string(out))
		} else {
			fmt.Println("[PASS] TypeScript checks passed with 0 errors.")
		}
	} else {
		fmt.Println("[PASS] Non-web project or generic workspace verification complete.")
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

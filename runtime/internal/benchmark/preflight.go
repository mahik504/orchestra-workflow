package benchmark

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PreflightResult is the operator check before a real pair. It stores no paths and no transcript.
type PreflightResult struct {
	OK            bool   `json:"ok"`
	Version       string `json:"version,omitempty"`
	AgentBase     string `json:"agent_base,omitempty"`
	SlugResolved  bool   `json:"slug_resolved"`
	SmokeWrote    bool   `json:"smoke_wrote"`
	SupportsForce bool   `json:"supports_force"`
	SupportsTrust bool   `json:"supports_trust"`
	Blocker       string `json:"blocker,omitempty"`
}

// PreflightAgent checks the absolute agent executable. It does not edit PATH.
// A failure is a blocker. It does not try a second model.
func PreflightAgent(ctx context.Context) PreflightResult {
	agent := strings.TrimSpace(os.Getenv("ORCHESTRA_BENCH_AGENT"))
	if agent == "" {
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			agent = filepath.Join(local, "cursor-agent", "agent.cmd")
		}
	}
	res := PreflightResult{AgentBase: safeBase(agent)}
	if agent == "" {
		res.Blocker = "agent executable was not found"
		return res
	}
	if _, err := os.Stat(agent); err != nil {
		res.Blocker = "agent executable was not found"
		return res
	}
	label := strings.TrimSpace(os.Getenv("ORCHESTRA_BENCH_MODEL_LABEL"))
	if label == "" {
		res.Blocker = "model label is unset"
		return res
	}
	verCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	ver, err := exec.CommandContext(verCtx, agent, "--version").CombinedOutput()
	cancel()
	if err != nil {
		res.Blocker = "agent --version failed"
		return res
	}
	res.Version = safeVersion(string(ver))
	helpCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	help, err := exec.CommandContext(helpCtx, agent, "--help").CombinedOutput()
	cancel()
	if err != nil {
		res.Blocker = "agent --help failed"
		return res
	}
	helpText := string(help)
	if !flagListed(helpText, "--print") {
		res.Blocker = "agent help does not list --print"
		return res
	}
	res.SupportsForce = flagListed(helpText, "--force")
	res.SupportsTrust = flagListed(helpText, "--trust")
	modelsCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	models, err := exec.CommandContext(modelsCtx, agent, "models").CombinedOutput()
	cancel()
	if err != nil {
		res.Blocker = "agent models failed"
		return res
	}
	if !slugListed(string(models), label) {
		res.Blocker = "model label is not in the agent model list"
		return res
	}
	res.SlugResolved = true
	wrote, smokeErr := smokeWrite(ctx, agent, label, res.SupportsForce, res.SupportsTrust)
	if smokeErr != nil || !wrote {
		res.Blocker = "headless write smoke failed"
		return res
	}
	res.SmokeWrote = true
	res.OK = true
	return res
}

func flagListed(help, flag string) bool {
	return strings.Contains(help, flag)
}

func slugListed(models, label string) bool {
	label = strings.ToLower(strings.TrimSpace(label))
	if label == "" {
		return false
	}
	for _, field := range strings.Fields(models) {
		if strings.ToLower(field) == label {
			return true
		}
	}
	return false
}

func safeBase(path string) string {
	base := filepath.Base(path)
	if base == "." || base == string(filepath.Separator) || strings.Contains(base, `\`) {
		return ""
	}
	return base
}

func safeVersion(raw string) string {
	line := strings.TrimSpace(raw)
	if i := strings.IndexAny(line, "\r\n"); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	if strings.Contains(line, `\`) || strings.Contains(line, "Users") || strings.Contains(line, "@") {
		return "resolved"
	}
	if len(line) > 80 {
		return "resolved"
	}
	return line
}

func smokeWrite(ctx context.Context, agent, label string, force, trust bool) (bool, error) {
	dir, err := os.MkdirTemp("", "orch-bench-smoke-")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(dir)
	args := []string{"--print", "--workspace", dir, "--model", label}
	if force {
		args = append(args, "--force")
	}
	if trust {
		args = append(args, "--trust")
	}
	args = append(args, "--", "Create one file named marker.txt whose contents are the word ok. Do not create any other file.")
	sctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(sctx, agent, args...)
	cmd.Dir = dir
	_, err = cmd.CombinedOutput()
	if _, statErr := os.Stat(filepath.Join(dir, "marker.txt")); statErr != nil {
		if err != nil {
			return false, err
		}
		return false, statErr
	}
	return true, nil
}

func realGuards(fixture string) []Guard {
	var guards []Guard
	if fixture != "" {
		guards = append(guards, Guard{Class: "fixture", Path: fixture})
	}
	if root := methodRoot(); root != "" {
		guards = append(guards, Guard{Class: "method_repo", Path: root})
	}
	if brain := strings.TrimSpace(os.Getenv("ORCHESTRA_MEMORY_PATH")); brain != "" {
		if _, err := os.Stat(brain); err == nil {
			guards = append(guards, Guard{Class: "brain", Path: brain})
		}
	}
	if home := strings.TrimSpace(os.Getenv("ORCHESTRA_HOME")); home != "" {
		if info, err := os.Stat(home); err == nil && info.IsDir() {
			guards = append(guards, Guard{Class: "brain", Path: home})
			addFile(&guards, "brain", filepath.Join(home, "memory", "resource-memory.json"))
		}
	}
	if app := os.Getenv("APPDATA"); app != "" {
		addFile(&guards, "cursor_config", filepath.Join(app, "Cursor", "User", "settings.json"))
		addFile(&guards, "mcp_config", filepath.Join(app, "Cursor", "User", "mcp.json"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		addFile(&guards, "mcp_config", filepath.Join(home, ".cursor", "mcp.json"))
	}
	return guards
}

func addFile(guards *[]Guard, class, path string) {
	if _, err := os.Stat(path); err == nil {
		*guards = append(*guards, Guard{Class: class, Path: path})
	}
}

func methodRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "VERSION")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

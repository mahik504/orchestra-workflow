package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/orchestra-v3/internal/handoff"
)

// AgentAdapter defines the legacy interface for portable agent execution
type AgentAdapter interface {
	Name() string
	SupportedEnvironments() []string
	ExecutePlan(state *handoff.HandoffState) error
}

// HostAdapter defines the full interface for host integration, configuration, and skill paths
type HostAdapter interface {
	AgentAdapter
	GetActiveSkillsPath(userHome string) string
	GenerateConfig(destDir string) error
}

// CursorAdapter handles execution via Cursor IDE
type CursorAdapter struct{}

func (a *CursorAdapter) Name() string { return "cursor" }
func (a *CursorAdapter) SupportedEnvironments() []string { return []string{"vscode", "cursor"} }
func (a *CursorAdapter) ExecutePlan(state *handoff.HandoffState) error {
	fmt.Println("[Adapter: Cursor] Writing context and instructing IDE...")
	return nil
}
func (a *CursorAdapter) GetActiveSkillsPath(userHome string) string {
	return filepath.Join(userHome, ".cursor", "skills")
}
func (a *CursorAdapter) GenerateConfig(destDir string) error {
	rulesPath := filepath.Join(destDir, ".cursorrules")
	if _, err := os.Stat(rulesPath); err == nil {
		return nil
	}
	content := `# Cursor Rules — Orchestra V3 Architecture & Capability Contract
You are executing tasks within the Orchestra V3 Resource Ecosystem.
Orchestra V3 operates on an 8-stage capability pipeline:
Discover -> Classify -> Research -> Synthesize -> Design System -> Implement -> Visual QA -> Iterate.
`
	return os.WriteFile(rulesPath, []byte(content), 0644)
}

// ClaudeAdapter handles execution via Claude Code CLI
type ClaudeAdapter struct{}

func (a *ClaudeAdapter) Name() string { return "claude" }
func (a *ClaudeAdapter) SupportedEnvironments() []string { return []string{"cli", "terminal"} }
func (a *ClaudeAdapter) ExecutePlan(state *handoff.HandoffState) error {
	fmt.Println("[Adapter: Claude Code] Generating CLAUDE.md context and terminal execution packet...")
	return nil
}
func (a *ClaudeAdapter) GetActiveSkillsPath(userHome string) string {
	return filepath.Join(userHome, ".claude", "skills")
}
func (a *ClaudeAdapter) GenerateConfig(destDir string) error {
	claudePath := filepath.Join(destDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		return nil
	}
	content := `# Claude Code Agentic System Matrix — Orchestra V3
You are Claude Code operating in the Orchestra V3 framework.
`
	return os.WriteFile(claudePath, []byte(content), 0644)
}

// AntigravityAdapter handles execution via Google Antigravity Agent
type AntigravityAdapter struct{}

func (a *AntigravityAdapter) Name() string { return "antigravity" }
func (a *AntigravityAdapter) SupportedEnvironments() []string { return []string{"cli", "ide"} }
func (a *AntigravityAdapter) ExecutePlan(state *handoff.HandoffState) error {
	fmt.Println("[Adapter: Antigravity] Dispatching JSON handoff state to Antigravity runtime...")
	return nil
}
func (a *AntigravityAdapter) GetActiveSkillsPath(userHome string) string {
	return filepath.Join(userHome, ".gemini", "config", "skills")
}
func (a *AntigravityAdapter) GenerateConfig(destDir string) error {
	mpPath := filepath.Join(destDir, "kit", "antigravity", "MASTER-PROMPT.md")
	if _, err := os.Stat(mpPath); err == nil {
		return nil
	}
	return nil
}

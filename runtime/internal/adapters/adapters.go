package adapters

import (
	"fmt"
	"github.com/user/orchestra-v3/internal/handoff"
)

// AgentAdapter defines the interface for portable agent integration
type AgentAdapter interface {
	Name() string
	SupportedEnvironments() []string
	ExecutePlan(state *handoff.HandoffState) error
}

// CursorAdapter handles execution via Cursor IDE
type CursorAdapter struct{}

func (a *CursorAdapter) Name() string { return "cursor" }
func (a *CursorAdapter) SupportedEnvironments() []string { return []string{"vscode", "cursor"} }
func (a *CursorAdapter) ExecutePlan(state *handoff.HandoffState) error {
	fmt.Println("[Adapter: Cursor] Writing context and instructing IDE...")
	// Implementation would trigger IDE specific files e.g. .cursorrules or open files
	return nil
}

// AntigravityAdapter handles execution via Antigravity Agent
type AntigravityAdapter struct{}

func (a *AntigravityAdapter) Name() string { return "antigravity" }
func (a *AntigravityAdapter) SupportedEnvironments() []string { return []string{"cli", "ide"} }
func (a *AntigravityAdapter) ExecutePlan(state *handoff.HandoffState) error {
	fmt.Println("[Adapter: Antigravity] Dispatching JSON handoff state to Antigravity runtime...")
	// Integration with AGY CLI or memory injection
	return nil
}

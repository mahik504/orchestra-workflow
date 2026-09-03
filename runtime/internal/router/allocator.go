package router

import (
	"github.com/user/orchestra-v3/internal/classifier"
)

type AgentType string

const (
	AgentCursor      AgentType = "Cursor"      // Implementation
	AgentAntigravity AgentType = "Antigravity" // Research/Visual Redesign
	AgentClaudeCode  AgentType = "ClaudeCode"  // Hostile Review
)

type AllocationRecommendation struct {
	PrimaryAgent AgentType
	Model        string
	Mode         string
	Effort       string
	Reason       string
}

type Allocator struct{}

func NewAllocator() *Allocator {
	return &Allocator{}
}

func (a *Allocator) Allocate(task *classifier.Task) AllocationRecommendation {
	if task.RequiresVisual {
		return AllocationRecommendation{
			PrimaryAgent: AgentAntigravity,
			Model:        "claude-3-5-sonnet",
			Mode:         "architect",
			Effort:       "high",
			Reason:       "Task requires visual redesign or heavy research, best suited for Antigravity.",
		}
	}
	
	if task.RequiresSecurity {
		return AllocationRecommendation{
			PrimaryAgent: AgentClaudeCode,
			Model:        "claude-3-7-sonnet",
			Mode:         "hostile-review",
			Effort:       "high",
			Reason:       "Task involves security or requires hostile review, best suited for Claude Code.",
		}
	}
	
	// Default to Cursor for implementation tasks
	return AllocationRecommendation{
		PrimaryAgent: AgentCursor,
		Model:        "claude-3-5-sonnet",
		Mode:         "implement",
		Effort:       "medium",
		Reason:       "Standard implementation task, routing to Cursor.",
	}
}

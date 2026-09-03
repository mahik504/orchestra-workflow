package router

import (
	"fmt"
	"strings"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

type Router struct {
	Registry *resources.Registry
}

func NewRouter(reg *resources.Registry) *Router {
	return &Router{Registry: reg}
}

type CapabilityExecutionDirective struct {
	CapabilityID          string   `json:"capability_id"`
	ActionableRules       []string `json:"actionable_rules"`
	AntiPatterns          []string `json:"anti_patterns"`
	VerificationChecklist []string `json:"verification_checklist"`
}

type CompositionPlan struct {
	SelectedCapabilities []*resources.Capability         `json:"selected_capabilities"`
	ExecutionDirectives  []*CapabilityExecutionDirective `json:"execution_directives"`
	GapResearchNeeded    []string                        `json:"gap_research_needed,omitempty"`
	EstimatedTokenCost   float64                         `json:"estimated_token_cost"`
	RequiresHumanGate    bool                            `json:"requires_human_gate"`
	ApprovalReason       string                          `json:"approval_reason"`
	PipelineStage        string                          `json:"pipeline_stage"` // "RETRIEVAL -> ANALYSIS -> APPLICATION -> VERIFICATION"
}

// Compose selects minimal sufficient capabilities AND synthesizes actionable execution contracts
func (r *Router) Compose(task *classifier.Task) *CompositionPlan {
	plan := &CompositionPlan{
		PipelineStage: "RETRIEVAL -> ANALYSIS -> APPLICATION -> VERIFICATION",
	}

	// 1. Baseline Methodology
	if cap, ok := r.Registry.Capabilities["superpowers-planning"]; ok {
		plan.SelectedCapabilities = append(plan.SelectedCapabilities, cap)
		plan.EstimatedTokenCost += cap.TokenContextWeight
		_ = cap.LoadDetails()

		plan.ExecutionDirectives = append(plan.ExecutionDirectives, &CapabilityExecutionDirective{
			CapabilityID: "superpowers-planning",
			ActionableRules: []string{
				"Deconstruct task into independent testable milestones before coding",
				"Execute Phase 0 (Specification & Contract Lock) prior to modifying production code",
				"Enforce adversarial verification gates at each milestone",
			},
			AntiPatterns: []string{
				"No coding without approved design/contract specification",
				"No merging without explicit multi-viewport verification",
			},
			VerificationChecklist: []string{
				"Does an approved specification document exist?",
				"Are milestone deliverables independently verified?",
			},
		})
	}

	// 2. Visual & Design Governance
	if task.RequiresVisual {
		if cap, ok := r.Registry.Capabilities["taste-skill"]; ok {
			plan.SelectedCapabilities = append(plan.SelectedCapabilities, cap)
			_ = cap.LoadDetails()

			plan.ExecutionDirectives = append(plan.ExecutionDirectives, &CapabilityExecutionDirective{
				CapabilityID: "taste-skill",
				ActionableRules: []string{
					"Enforce strict typographic pairing: distinctive Display Serif + Clean Sans + Monospace figures",
					"Asymmetric compositions over repetitive 3-column card rows",
					"Single calibrated accent with saturation under 80%",
					"High-density data must use tabular lining and monospace numbers",
				},
				AntiPatterns: []string{
					"NO Inter or Space Grotesk in creative/premium contexts",
					"NO neon gradients, cosmic glows, or AI purple glows",
					"NO generic 3-equal-cards feature rows",
					"NO pure black (#000000) surfaces",
				},
				VerificationChecklist: []string{
					"Grep confirmation: 0 instances of Inter or Space Grotesk in built assets",
					"Viewport inspection: 0 centered 3-card monotony rows",
					"Contrast ratio: 4.5:1 minimum on all text",
				},
			})
		}

		if cap, ok := r.Registry.Capabilities["impeccable"]; ok {
			plan.SelectedCapabilities = append(plan.SelectedCapabilities, cap)
			_ = cap.LoadDetails()

			plan.ExecutionDirectives = append(plan.ExecutionDirectives, &CapabilityExecutionDirective{
				CapabilityID: "impeccable",
				ActionableRules: []string{
					"Maintain strict layout containment and responsive spacing clamps",
					"Hardware-accelerated CSS transforms (scale, opacity) exclusively",
					"Tactile button feedback on :active (scale 0.98)",
				},
				AntiPatterns: []string{
					"NO scale(0) entry transitions",
					"NO animation of layout-triggering properties (width, height, top, padding)",
				},
				VerificationChecklist: []string{
					"Mobile drawer handles viewport resize and locks body scroll",
					"No horizontal overflow on mobile viewports (<390px)",
				},
			})
		}

		if task.UserOverride != nil && (task.UserOverride.SkipVisualGate || task.UserOverride.ForceBypassGate) {
			plan.RequiresHumanGate = false
			plan.ApprovalReason = "Bypassed visual gate via user override"
		} else {
			plan.RequiresHumanGate = true
			plan.ApprovalReason = fmt.Sprintf("High-impact visual task requires design laboratory approval with %d active directives", len(plan.ExecutionDirectives))
		}
	}

	// 3. Security & Code Integrity
	if task.RequiresSecurity {
		if cap, ok := r.Registry.Capabilities["semgrep-adapter"]; ok {
			plan.SelectedCapabilities = append(plan.SelectedCapabilities, cap)
			_ = cap.LoadDetails()

			plan.ExecutionDirectives = append(plan.ExecutionDirectives, &CapabilityExecutionDirective{
				CapabilityID: "semgrep-adapter",
				ActionableRules: []string{
					"Audit all external API routes for input sanitization and honeypot spam traps",
					"Ensure zero client-side leakage of serverless environment secrets",
				},
				AntiPatterns: []string{
					"NO hardcoded API keys or credentials in client bundles",
				},
				VerificationChecklist: []string{
					"Serverless route /api/inquiry parses and sanitizes inputs server-side",
				},
			})
		}
	}

	// 4. Capability Gap Detection & Research Trigger
	for _, res := range task.SuggestedResources {
		if _, exists := r.Registry.Capabilities[res]; !exists {
			plan.GapResearchNeeded = append(plan.GapResearchNeeded, res)
		}
	}
	if len(plan.GapResearchNeeded) > 0 {
		plan.ExecutionDirectives = append(plan.ExecutionDirectives, &CapabilityExecutionDirective{
			CapabilityID: "capability-gap-research",
			ActionableRules: []string{
				fmt.Sprintf("Trigger deep research for %d unindexed technologies: %s", len(plan.GapResearchNeeded), strings.Join(plan.GapResearchNeeded, ", ")),
				"Synthesize official docs, version constraints, and known gotchas prior to synthesis",
			},
			AntiPatterns: []string{
				"Do NOT guess API signatures of unindexed libraries",
			},
			VerificationChecklist: []string{
				"Verify official SDK documentation has been fetched and indexed",
			},
		})
	}

	return plan
}

// GenerateExecutionManifest formats the active directives as a markdown manifest for downstream agents
func (p *CompositionPlan) GenerateExecutionManifest() string {
	var sb strings.Builder
	sb.WriteString("# Orchestra V3 Capability Execution Manifest\n\n")
	sb.WriteString(fmt.Sprintf("**Pipeline Stage:** %s\n", p.PipelineStage))
	sb.WriteString(fmt.Sprintf("**Active Directives:** %d\n\n", len(p.ExecutionDirectives)))

	for _, d := range p.ExecutionDirectives {
		sb.WriteString(fmt.Sprintf("## Capability: %s\n", d.CapabilityID))
		sb.WriteString("### Mandatory Actionable Rules:\n")
		for _, r := range d.ActionableRules {
			sb.WriteString(fmt.Sprintf("- %s\n", r))
		}
		sb.WriteString("### Banned Anti-Patterns:\n")
		for _, a := range d.AntiPatterns {
			sb.WriteString(fmt.Sprintf("- %s\n", a))
		}
		sb.WriteString("### Verification Checklist:\n")
		for _, v := range d.VerificationChecklist {
			sb.WriteString(fmt.Sprintf("- [ ] %s\n", v))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/router"
)

type SynthesizeStage struct{}

func NewSynthesizeStage() *SynthesizeStage {
	return &SynthesizeStage{}
}

func (s *SynthesizeStage) Name() StageName {
	return StageSynthesize
}

func (s *SynthesizeStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *SynthesizeStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	var actionableRules []string
	var antiPatterns []string
	var verificationChecks []string
	var directives []*router.CapabilityExecutionDirective

	// 1. Baseline Engineering Directives
	superpowersDir := &router.CapabilityExecutionDirective{
		CapabilityID: "superpowers-planning",
		ActionableRules: []string{
			"Deconstruct task into independent testable milestones before coding",
			"Execute Phase 0 (Contract & Spec Lock) prior to modifying production code",
			"Enforce adversarial verification gates at each milestone",
		},
		AntiPatterns: []string{
			"No coding without approved design/contract specification",
			"No merging without explicit multi-viewport verification",
		},
		VerificationChecklist: []string{
			"Does an approved specification or DESIGN.md contract exist?",
			"Are milestone deliverables independently verified?",
		},
	}
	directives = append(directives, superpowersDir)
	actionableRules = append(actionableRules, superpowersDir.ActionableRules...)
	antiPatterns = append(antiPatterns, superpowersDir.AntiPatterns...)
	verificationChecks = append(verificationChecks, superpowersDir.VerificationChecklist...)

	// 2. Visual & Architectural Directives
	if ctx.Classification != nil && ctx.Classification.RequiresVisual {
		tasteDir := &router.CapabilityExecutionDirective{
			CapabilityID: "taste-skill",
			ActionableRules: []string{
				"Enforce strict typographic pairing: distinctive Display Serif + Clean Sans + Monospace figures",
				"Asymmetric compositions over repetitive 3-column card rows",
				"Single calibrated accent with saturation under 80%",
				"High-density data must use tabular lining and monospace numbers",
			},
			AntiPatterns: []string{
				"NO Inter or Space Grotesk in creative/display headlines",
				"NO neon gradients, cosmic glows, or AI purple glows",
				"NO generic 3-equal-cards feature rows",
				"NO pure black (#000000) surfaces",
			},
			VerificationChecklist: []string{
				"Grep confirmation: 0 instances of Inter or Space Grotesk in creative headlines",
				"Viewport inspection: 0 centered 3-card monotony rows",
				"Contrast ratio: 4.5:1 minimum on all text",
			},
		}
		directives = append(directives, tasteDir)
		actionableRules = append(actionableRules, tasteDir.ActionableRules...)
		antiPatterns = append(antiPatterns, tasteDir.AntiPatterns...)
		verificationChecks = append(verificationChecks, tasteDir.VerificationChecklist...)

		impeccableDir := &router.CapabilityExecutionDirective{
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
		}
		directives = append(directives, impeccableDir)
		actionableRules = append(actionableRules, impeccableDir.ActionableRules...)
		antiPatterns = append(antiPatterns, impeccableDir.AntiPatterns...)
		verificationChecks = append(verificationChecks, impeccableDir.VerificationChecklist...)
	}

	// 3. Capability Route Directives & Anti-Patterns
	if ctx.Classification != nil {
		for _, route := range ctx.Classification.ResolvedRoutes {
			if len(route.AntiPatterns) > 0 {
				antiPatterns = append(antiPatterns, route.AntiPatterns...)
			}
			if len(route.QA) > 0 {
				for _, qa := range route.QA {
					verificationChecks = append(verificationChecks, fmt.Sprintf("Run %s verification check", qa))
				}
			}
		}
	}

	// 4. Security Directives
	if ctx.Classification != nil && ctx.Classification.RequiresSecurity {
		secDir := &router.CapabilityExecutionDirective{
			CapabilityID: "semgrep-adapter",
			ActionableRules: []string{
				"Audit all external API routes for input sanitization and honeypot spam traps",
				"Ensure zero client-side leakage of serverless environment secrets",
			},
			AntiPatterns: []string{
				"NO hardcoded API keys or credentials in client bundles",
			},
			VerificationChecklist: []string{
				"Serverless routes parse and sanitize inputs server-side",
			},
		}
		directives = append(directives, secDir)
		actionableRules = append(actionableRules, secDir.ActionableRules...)
		antiPatterns = append(antiPatterns, secDir.AntiPatterns...)
		verificationChecks = append(verificationChecks, secDir.VerificationChecklist...)
	}

	// 5. Gap Technologies Research Trigger
	if ctx.Classification != nil && len(ctx.Classification.GapTechnologies) > 0 {
		gapDir := &router.CapabilityExecutionDirective{
			CapabilityID: "capability-gap-research",
			ActionableRules: []string{
				fmt.Sprintf("Synthesize official documentation for unindexed technologies: %s", strings.Join(ctx.Classification.GapTechnologies, ", ")),
				"Do NOT guess API signatures of unindexed libraries",
			},
			AntiPatterns: []string{
				"No hallucinated method names or incompatible major versions",
			},
			VerificationChecklist: []string{
				"Verify official SDK documentation has been fetched and indexed",
			},
		}
		directives = append(directives, gapDir)
		actionableRules = append(actionableRules, gapDir.ActionableRules...)
		antiPatterns = append(antiPatterns, gapDir.AntiPatterns...)
		verificationChecks = append(verificationChecks, gapDir.VerificationChecklist...)
	}

	// 6. Cumulative Token Context Cost Calculation
	totalTokenCost := 1500.0 // Baseline superpowers-planning
	seenResources := make(map[string]bool)

	if ctx.Classification != nil && ctx.Catalog != nil {
		for _, route := range ctx.Classification.ResolvedRoutes {
			for _, resID := range route.AllResourceIDs {
				if !seenResources[resID] {
					seenResources[resID] = true
					if res, found := ctx.Catalog.FindByID(resID); found {
						w := res.TokenContextWeight
						if w <= 0 {
							w = res.TokenWeight
						}
						if w <= 0 {
							w = 800.0
						}
						totalTokenCost += w
					} else {
						totalTokenCost += 500.0
					}
				}
			}
		}
	}

	// 7. Human Approval Gate Verification
	directionApproved := true
	if ctx.Classification != nil && ctx.Classification.RequiresHumanGate {
		if !ctx.Task.SkipVisualGate && !ctx.Task.DryRun {
			// Enforce Gate Halting
			return &StageResult{
				StageName: StageSynthesize,
				Status:    StatusFailed,
				StartTime: start,
				EndTime:   time.Now(),
				Duration:  time.Since(start),
				Error:     ErrHumanGateRequired,
			}, ErrHumanGateRequired
		}
		if ctx.Task.DryRun {
			directionApproved = false
		}
	}

	ctx.Synthesis = &SynthesisData{
		ActionableRules:    actionableRules,
		AntiPatterns:       antiPatterns,
		VerificationChecks: verificationChecks,
		TokenContextCost:   totalTokenCost,
		DirectionApproved:  directionApproved,
		ActiveDirectives:   directives,
	}

	return &StageResult{
		StageName: StageSynthesize,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Synthesis,
	}, nil
}

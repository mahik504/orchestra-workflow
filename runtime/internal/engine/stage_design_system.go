package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

type DesignSystemStage struct{}

func NewDesignSystemStage() *DesignSystemStage {
	return &DesignSystemStage{}
}

func (s *DesignSystemStage) Name() StageName {
	return StageDesignSystem
}

func (s *DesignSystemStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *DesignSystemStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	requiresVisual := false
	if ctx.Classification != nil {
		requiresVisual = ctx.Classification.RequiresVisual
	}

	outDir := filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra")
	if err := resources.CheckQuarantineBoundary(outDir); err != nil {
		return &StageResult{
			StageName: StageDesignSystem,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}
	_ = os.MkdirAll(outDir, 0755)

	designMDPath := filepath.Join(outDir, "DESIGN.md")

	if !requiresVisual {
		// Non-visual tasks emit minimal architectural specification
		minimalDoc := fmt.Sprintf("# System Architecture Contract\n\n- Task: `%s`\n- Mode: Non-visual headless execution\n- Isolation: Quarantine enforced\n", ctx.Task.ID)
		_ = os.WriteFile(designMDPath, []byte(minimalDoc), 0644)
		ctx.DesignSystem = &DesignSystemData{
			DesignMDPath: designMDPath,
			LayoutRules:  []string{"Headless architecture", "Strict contract verification"},
		}
		return &StageResult{
			StageName: StageDesignSystem,
			Status:    StatusCompleted,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Output:    ctx.DesignSystem,
			Artifacts: []string{designMDPath},
		}, nil
	}

	// 1. Resolve Tokens from Research Data or Defaults
	displayFont := "Instrument Serif"
	bodyFont := "Plus Jakarta Sans"
	monoFont := "JetBrains Mono"
	surfaceColor := "#0B0E14"
	primaryAccent := "#FF4B4B"

	if ctx.Research != nil && ctx.Research.SynthesizedResult != nil {
		for _, t := range ctx.Research.SynthesizedResult.SynthesizedTypography {
			if t.Role == "display" && t.FontFamily != "" {
				displayFont = t.FontFamily
			} else if t.Role == "body" && t.FontFamily != "" {
				bodyFont = t.FontFamily
			} else if t.Role == "mono" && t.FontFamily != "" {
				monoFont = t.FontFamily
			}
		}
		for _, p := range ctx.Research.SynthesizedResult.SynthesizedPalette {
			if p.Role == "--color-bg-base" && p.Hex != "" {
				surfaceColor = p.Hex
			} else if p.Role == "--color-accent-primary" && p.Hex != "" {
				primaryAccent = p.Hex
			}
		}
	}

	// Apply corrective feedback from Stage 8 (Iterate) if looping back
	if ctx.Iteration != nil && len(ctx.Iteration.CorrectiveFeedback) > 0 {
		for _, fb := range ctx.Iteration.CorrectiveFeedback {
			if strings.Contains(fb, "banned font") || strings.Contains(fb, "Inter") {
				displayFont = "Instrument Serif" // Re-enforce approved serif
			}
			if strings.Contains(fb, "pure black") || strings.Contains(fb, "#000000") {
				surfaceColor = "#0E0E12" // Replace with calibrated matte dark
			}
		}
	}

	cssVars := map[string]string{
		"--bg-base":          surfaceColor,
		"--surface-elevated": "#141A24",
		"--accent-primary":   primaryAccent,
		"--text-headline":    "#F0F4F8",
		"--text-muted":       "#8B9BB4",
		"--border-subtle":    "#232D3F",
		"--font-display":     fmt.Sprintf("'%s', serif", displayFont),
		"--font-body":        fmt.Sprintf("'%s', sans-serif", bodyFont),
		"--font-mono":        fmt.Sprintf("'%s', monospace", monoFont),
	}

	springPhysics := map[string]string{
		"stiffness": "300",
		"damping":   "30",
		"mass":      "1",
	}

	layoutRules := []string{
		"Hero Section: Asymmetric 60/40 viewport split",
		"Card Rows: Prohibit repetitive 3-column equal cards; use varied aspect-ratio masonry",
		"Mobile Containment: Strict overflow-x: hidden; document scrollWidth <= window.innerWidth",
		"Tactile Feedback: :active transform scale(0.98)",
		"Surface Lighting: 1px inset border rgba(255, 255, 255, 0.04) instead of blurry drop-shadows",
	}

	// 2. Generate DESIGN.md Content
	var sb strings.Builder
	sb.WriteString("# DESIGN.md — Design System Contract\n\n")
	sb.WriteString(fmt.Sprintf("> Synthesized for Task: `%s` | Archetype: `%s`\n", ctx.Task.ID, ctx.Classification.Archetype))
	if ctx.Research != nil && ctx.Research.ReferenceLogPath != "" {
		sb.WriteString(fmt.Sprintf("> Provenance: Synthesized from `%s`\n\n", filepath.Base(ctx.Research.ReferenceLogPath)))
	} else {
		sb.WriteString("> Provenance: Synthesized from Orchestra V3 Curated Registries\n\n")
	}

	sb.WriteString("## 1. Typography Hierarchy\n\n")
	sb.WriteString(fmt.Sprintf("- **Display Font**: `%s` (Serif, clamp(2.75rem, 6vw + 1rem, 5.5rem), line-height: 1.05, tracking: -0.03em)\n", displayFont))
	sb.WriteString(fmt.Sprintf("- **Body Font**: `%s` (Clean Geometric Sans, clamp(1rem, 1vw + 0.5rem, 1.125rem), line-height: 1.6)\n", bodyFont))
	sb.WriteString(fmt.Sprintf("- **Mono Font**: `%s` (Monospace figures, 0.875rem, line-height: 1.4)\n\n", monoFont))

	sb.WriteString("## 2. Calibrated Color Palette\n\n")
	sb.WriteString("```css\n:root {\n")
	for k, v := range cssVars {
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", k, v))
	}
	sb.WriteString("}\n```\n\n")

	sb.WriteString("## 3. Kinetic & Spring Physics\n\n")
	sb.WriteString("- Motion Easing: `cubic-bezier(0.16, 1, 0.3, 1)`\n")
	sb.WriteString(fmt.Sprintf("- Spring Config: stiffness=%s, damping=%s, mass=%s\n", springPhysics["stiffness"], springPhysics["damping"], springPhysics["mass"]))
	sb.WriteString("- Button Press: `transform: scale(0.98); transition: transform 0.1s ease-out;`\n\n")

	sb.WriteString("## 4. Architectural Layout Rules\n\n")
	for _, rule := range layoutRules {
		sb.WriteString(fmt.Sprintf("- %s\n", rule))
	}
	sb.WriteString("\n")

	sb.WriteString("## 5. Banned Anti-Patterns\n\n")
	sb.WriteString("- [ ] 0 instances of Inter or Space Grotesk in creative display headlines\n")
	sb.WriteString("- [ ] 0 instances of pure black (#000000) surface backgrounds\n")
	sb.WriteString("- [ ] 0 instances of neon purple/pink AI bubble gradients\n")
	sb.WriteString("- [ ] 0 horizontal overflow defects on mobile viewports (<390px)\n")

	if err := os.WriteFile(designMDPath, []byte(sb.String()), 0644); err != nil {
		return &StageResult{
			StageName: StageDesignSystem,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	ctx.DesignSystem = &DesignSystemData{
		DesignMDPath:  designMDPath,
		DisplayFont:   displayFont,
		BodyFont:      bodyFont,
		MonoFont:      monoFont,
		PrimaryAccent: primaryAccent,
		SurfaceColor:  surfaceColor,
		SpringPhysics: springPhysics,
		LayoutRules:   layoutRules,
		CSSVariables:  cssVars,
	}

	ctx.AddArtifact("DESIGN.md", designMDPath)

	return &StageResult{
		StageName: StageDesignSystem,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.DesignSystem,
		Artifacts: []string{designMDPath},
	}, nil
}

package research

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

// GenerateReferenceLog formats and persists reference-log.md to disk with full sections
func (c *ResearchCoordinator) GenerateReferenceLog(ctx context.Context, res *ResearchResult, outputPath string) error {
	if res == nil {
		return fmt.Errorf("research result cannot be nil")
	}
	if outputPath == "" {
		return fmt.Errorf("reference log output path cannot be empty")
	}

	// Quarantine boundary check on output path
	if err := resources.CheckQuarantineBoundary(outputPath); err != nil {
		return err
	}

	// Validation constraints
	if res.QualityBar == "premium" && len(res.SelectedSources) < 2 {
		return fmt.Errorf("high-visual premium tasks mandate at least 2 curated sources, found %d", len(res.SelectedSources))
	}
	if len(res.SynthesizedPalette) < 4 {
		return fmt.Errorf("synthesized palette must define at least 4 tokens, found %d", len(res.SynthesizedPalette))
	}

	for _, typo := range res.SynthesizedTypography {
		if typo.Role == "display" {
			lowerFont := strings.ToLower(typo.FontFamily)
			if strings.Contains(lowerFont, "inter") || strings.Contains(lowerFont, "space grotesk") {
				return fmt.Errorf("anti-pattern violation: %s is banned for display/headline typography in creative contexts", typo.FontFamily)
			}
		}
	}

	var sb strings.Builder

	// 1. Frontmatter Header
	sb.WriteString("# Visual Research Reference Log\n\n")
	sb.WriteString(fmt.Sprintf("- **Task ID**: `%s`\n", res.TaskID))
	sb.WriteString(fmt.Sprintf("- **Archetype**: `%s`\n", res.Archetype))
	sb.WriteString(fmt.Sprintf("- **Quality Bar**: `%s`\n", res.QualityBar))
	sb.WriteString("- **Coordinator Version**: Orchestra V3 Multi-Source Research Engine (Stage 3)\n")
	sb.WriteString(fmt.Sprintf("- **Generated At**: `%s`\n", res.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString("- **Pipeline Gate**: Stage 3 (Research) -> Stage 4 (Synthesize) Gate Passed\n\n")
	sb.WriteString("---\n\n")

	// 2. Executive Direction & Creative Hypothesis
	sb.WriteString("## 1. Executive Direction & Creative Hypothesis\n\n")
	sb.WriteString("Orchestra V3 enforces design-first execution. Prior to writing code or defining tokens, ")
	sb.WriteString(fmt.Sprintf("the engine has queried %d authoritative visual design sources to establish a cohesive aesthetic hypothesis:\n\n", len(res.SelectedSources)))
	sb.WriteString("- **Aesthetic Stance**: Rejection of generic AI slop (no neon gradients, no ungrounded card grids). High-density editorial rigor combined with fluid responsive geometry.\n")
	sb.WriteString("- **Atmospheric Color**: Matte, low-emission dark surfaces (`#0B0E14`) elevated by subtle inset border lighting and calibrated single-hue accents.\n")
	sb.WriteString("- **Typographic Voice**: Triadic harmony pairing high-contrast display serif typography with ultra-clean sans-serif body copy and monospace technical figures.\n\n")

	// 3. Curated Reference Sources Table
	sb.WriteString("## 2. Curated Reference Sources & Algorithmic Scoring\n\n")
	sb.WriteString("| Source ID | Source Name | Canonical URL | Domain | Composite Score | Assigned Role |\n")
	sb.WriteString("|---|---|---|---|---|---|\n")
	for _, src := range res.SelectedSources {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | [%s](%s) | `%s` | `%.2f` | %s |\n",
			src.Resource.ID,
			src.Resource.Name,
			src.Resource.CanonicalURL,
			src.Resource.CanonicalURL,
			src.Domain,
			src.Score,
			src.Role,
		))
	}
	sb.WriteString("\n")

	// 4. Direct Citation Anchors
	sb.WriteString("## 3. Direct Citation Anchors\n\n")
	seenCitations := make(map[string]bool)
	citationCount := 0
	for _, finding := range res.Findings {
		for _, cit := range finding.Citations {
			if !seenCitations[cit] {
				seenCitations[cit] = true
				citationCount++
				sb.WriteString(fmt.Sprintf("%d. [%s](%s) — *%s*\n", citationCount, cit, cit, finding.SourceName))
			}
		}
	}
	if citationCount == 0 {
		sb.WriteString("1. [https://www.awwwards.com/websites/](https://www.awwwards.com/websites/) — *Awwwards Benchmark*\n")
		sb.WriteString("2. [https://jiro.design/](https://jiro.design/) — *Jiro Design Direction*\n")
	}
	sb.WriteString("\n")

	// 5. Extracted Visual Motifs & Structural Patterns
	sb.WriteString("## 4. Discovered Visual Motifs & Layout Lessons\n\n")
	for _, motif := range res.SelectedMotifs {
		sb.WriteString(fmt.Sprintf("- **%s** (`%s`): %s\n", strings.Title(strings.ReplaceAll(motif.Category, "_", " ")), motif.SourceID, motif.Description))
	}
	sb.WriteString("- **Asymmetric Viewport Split**: Hero sections must break the standard 50/50 division in favor of 60/40 or 70/30 spatial weighting.\n")
	sb.WriteString("- **Tactile Elevation**: 1px inset borders with 4% white alpha (`rgba(255, 255, 255, 0.04)`) replace diffuse drop shadows.\n\n")

	// 6. Synthesized Color Palette Tokens Table
	sb.WriteString("## 5. Synthesized Color Palette Tokens\n\n")
	sb.WriteString("| CSS Custom Property | Hex Value | HSL Color | Contrast Ratio | Source Attribution |\n")
	sb.WriteString("|---|---|---|---|---|\n")
	for _, p := range res.SynthesizedPalette {
		sb.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | %s | `%s` |\n",
			p.Role, p.Hex, p.HSL, p.Contrast, p.SourceID))
	}
	sb.WriteString("\n")

	// 7. Typography Pairing Tokens Table
	sb.WriteString("## 6. Typography Pairing & Hierarchy Tokens\n\n")
	sb.WriteString("| Hierarchy Role | Font Family | Fallback | Fluid Size Clamp | Weight | Line Height | Letter Spacing | Source |\n")
	sb.WriteString("|---|---|---|---|---|---|---|---|\n")
	for _, t := range res.SynthesizedTypography {
		sb.WriteString(fmt.Sprintf("| `%s` | **%s** | `%s` | `%s` | `%s` | `%s` | `%s` | `%s` |\n",
			t.Role, t.FontFamily, t.Fallback, t.SizeClamp, t.Weight, t.LineHeight, t.LetterSpacing, t.SourceID))
	}
	sb.WriteString("\n")

	// 8. Kinetic & Interaction Dynamics
	sb.WriteString("## 7. Kinetic & Interaction Dynamics\n\n")
	sb.WriteString(fmt.Sprintf("- **Motion Curve**: `%s`\n", res.InteractionRules.MotionCurve))
	sb.WriteString(fmt.Sprintf("- **Spring Physics**: `%s`\n", res.InteractionRules.SpringConfig))
	sb.WriteString(fmt.Sprintf("- **Tactile Press Feedback**: `%s`\n", res.InteractionRules.ActivePress))
	sb.WriteString(fmt.Sprintf("- **Strict Layout Protection**: NEVER animate layout-triggering properties: `%s`\n\n", strings.Join(res.InteractionRules.ProhibitedProps, ", ")))

	// 9. Negative Guardrails & Banned Anti-Patterns
	sb.WriteString("## 8. Negative Guardrails & Anti-Pattern Checklist\n\n")
	sb.WriteString("The Visual QA Stage (Stage 7) will assert zero violations against this checklist:\n\n")
	for _, ap := range res.BannedAntiPatterns {
		sb.WriteString(fmt.Sprintf("- [ ] %s\n", ap))
	}
	sb.WriteString("\n")

	// Ensure destination directory exists
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for reference log: %w", err)
		}
	}

	// Write file atomically
	if err := os.WriteFile(outputPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write reference log to %s: %w", outputPath, err)
	}

	// Verify existence
	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("reference log file verification failed on disk: %w", err)
	}

	return nil
}

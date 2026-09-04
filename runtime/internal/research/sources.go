package research

import (
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
)

// Archetype Families for cross-source diversity calculation
const (
	FamilyAestheticBenchmark = "aesthetic_benchmark"
	FamilyLayoutComposition  = "layout_composition"
	FamilyMovementTaxonomy   = "movement_taxonomy"
	FamilySpecialistEcho     = "specialist_echo"
	FamilyDefault            = "general_reference"
)

// SourceFamilyMap categorizes known curated resources into archetype families
var SourceFamilyMap = map[string]string{
	"awwwards":          FamilyAestheticBenchmark,
	"godly-design":      FamilyAestheticBenchmark,
	"jiro-design":       FamilyLayoutComposition,
	"siteinspire":       FamilyLayoutComposition,
	"land-book":         FamilyLayoutComposition,
	"cari-institute":    FamilyMovementTaxonomy,
	"awesome-design-md": FamilyMovementTaxonomy,
	"refero-design":     FamilySpecialistEcho,
	"k95-studio":        FamilySpecialistEcho,
	"noth-in":           FamilySpecialistEcho,
	"united-carriers":   FamilySpecialistEcho,
	"peter-oravec":      FamilySpecialistEcho,
	"shaders-com":       FamilySpecialistEcho,
	"lax-space":         FamilySpecialistEcho,
	"pryzm-design":      FamilySpecialistEcho,
}

// PaletteToken represents an extracted or synthesized design token for color
type PaletteToken struct {
	Role     string `json:"role"`     // e.g. "--color-bg-base", "--color-surface-elevated", "--color-accent-primary"
	Hex      string `json:"hex"`      // e.g. "#0B0E14"
	HSL      string `json:"hsl"`      // e.g. "hsl(220, 20%, 6%)"
	Contrast string `json:"contrast"` // e.g. "16.2:1 (AAA)"
	SourceID string `json:"source_id"`
}

// TypographyToken represents an extracted or synthesized typographic hierarchy token
type TypographyToken struct {
	Role          string `json:"role"`           // "display", "headline", "body", "mono"
	FontFamily    string `json:"font_family"`    // e.g. "Instrument Serif", "Plus Jakarta Sans"
	Fallback      string `json:"fallback"`       // "serif", "sans-serif", "monospace"
	SizeClamp     string `json:"size_clamp"`     // e.g. "clamp(2.5rem, 5vw + 1rem, 5rem)"
	Weight        string `json:"weight"`         // "400", "500", "600", "700"
	LineHeight    string `json:"line_height"`    // "1.1", "1.5"
	LetterSpacing string `json:"letter_spacing"` // "-0.03em", "0"
	SourceID      string `json:"source_id"`
}

// VisualMotif captures a concrete layout or visual pattern discovered during research
type VisualMotif struct {
	Category    string `json:"category"` // "hero_composition", "grid_structure", "surface_depth", "border_treatment"
	Description string `json:"description"`
	SourceID    string `json:"source_id"`
	SourceURL   string `json:"source_url"`
}

// InteractionDynamics defines kinetic motion curves, physics, and touch constraints
type InteractionDynamics struct {
	MotionCurve     string   `json:"motion_curve"`
	SpringConfig    string   `json:"spring_config"`
	ActivePress     string   `json:"active_press"`
	ProhibitedProps []string `json:"prohibited_props"`
	SourceID        string   `json:"source_id"`
}

// ReferenceFinding aggregates raw research findings for a single source
type ReferenceFinding struct {
	SourceID            string            `json:"source_id"`
	SourceName          string            `json:"source_name"`
	SourceURL           string            `json:"source_url"`
	Category            string            `json:"category"`
	KeyTakeaways        []string          `json:"key_takeaways"`
	ExtractedPalettes   []PaletteToken    `json:"extracted_palettes"`
	ExtractedTypography []TypographyToken `json:"extracted_typography"`
	VisualMotifs        []VisualMotif     `json:"visual_motifs"`
	Citations           []string          `json:"citations"`
}

// ScoredReference couples a canonical Resource with its algorithmic rank and match rationale
type ScoredReference struct {
	Resource    *resources.Resource `json:"resource"`
	Score       float64             `json:"score"`
	Domain      string              `json:"domain"`
	Family      string              `json:"family"`
	MatchReason string              `json:"match_reason"`
	Role        string              `json:"role"`
}

// SelectionOptions controls research candidate filtering and constraints
type SelectionOptions struct {
	MaxSources        int      `json:"max_sources"`
	MinSources        int      `json:"min_sources"`
	IncludeBookmarks  bool     `json:"include_bookmarks"`
	AllowExperimental bool     `json:"allow_experimental"`
	OfflineBenchmark  bool     `json:"offline_benchmark"`
	RequiredDomains   []string `json:"required_domains"`
}

// ResearchRequest packages the intake payload for the ResearchCoordinator
type ResearchRequest struct {
	Task             *classifier.Task           `json:"task"`
	Route            *resources.CapabilityRoute `json:"route"`
	Options          SelectionOptions           `json:"options"`
	ProjectOutputDir string                     `json:"project_output_dir"`
}

// ResearchResult represents the complete multi-source synthesis produced by Stage 3
type ResearchResult struct {
	TaskID                string              `json:"task_id"`
	Archetype             string              `json:"archetype"`
	QualityBar            string              `json:"quality_bar"`
	TotalSourcesQueried   int                 `json:"total_sources_queried"`
	SelectedSources       []*ScoredReference  `json:"selected_sources"`
	Findings              []*ReferenceFinding `json:"findings"`
	SynthesizedPalette    []PaletteToken      `json:"synthesized_palette"`
	SynthesizedTypography []TypographyToken   `json:"synthesized_typography"`
	SelectedMotifs        []VisualMotif       `json:"selected_motifs"`
	InteractionRules      InteractionDynamics `json:"interaction_rules"`
	BannedAntiPatterns    []string            `json:"banned_anti_patterns"`
	ReferenceLogPath      string              `json:"reference_log_path"`
	GeneratedAt           time.Time           `json:"generated_at"`
}

// CuratedSourceFixtures provides pre-indexed, verified profiles for offline resilience and fast determinism
var CuratedSourceFixtures = map[string]*ReferenceFinding{
	"awwwards": {
		SourceID:   "awwwards",
		SourceName: "Awwwards Benchmark",
		SourceURL:  "https://www.awwwards.com/",
		Category:   "award-winning-creative",
		KeyTakeaways: []string{
			"Asymmetric grid layouts disrupt predictable 3-column monotony and command visual authority",
			"Calibrated high-contrast dark surfaces (#0B0E14) prevent eye fatigue compared to harsh #000000",
			"Single vibrant accent with saturation strictly under 80% creates focused visual tension",
			"Fluid typography scaling via CSS clamp() preserves intentional editorial hierarchy across viewports",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0B0E14", HSL: "hsl(220, 20%, 6%)", Contrast: "17.4:1 against text", SourceID: "awwwards"},
			{Role: "--color-surface-elevated", Hex: "#141A24", HSL: "hsl(216, 28%, 11%)", Contrast: "13.8:1 against text", SourceID: "awwwards"},
			{Role: "--color-accent-primary", Hex: "#FF4B4B", HSL: "hsl(0, 100%, 65%)", Contrast: "4.8:1 against bg", SourceID: "awwwards"},
			{Role: "--color-text-headline", Hex: "#F0F4F8", HSL: "hsl(210, 33%, 96%)", Contrast: "17.4:1 against bg-base", SourceID: "awwwards"},
			{Role: "--color-text-muted", Hex: "#8B9BB4", HSL: "hsl(216, 21%, 63%)", Contrast: "6.2:1 against bg-base", SourceID: "awwwards"},
			{Role: "--color-border-subtle", Hex: "#232D3F", HSL: "hsl(218, 28%, 19%)", Contrast: "2.5:1 decorative", SourceID: "awwwards"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "display", FontFamily: "Instrument Serif", Fallback: "serif", SizeClamp: "clamp(2.75rem, 6vw + 1rem, 5.5rem)", Weight: "400", LineHeight: "1.05", LetterSpacing: "-0.03em", SourceID: "awwwards"},
			{Role: "headline", FontFamily: "Plus Jakarta Sans", Fallback: "sans-serif", SizeClamp: "clamp(1.5rem, 3vw + 0.5rem, 2.5rem)", Weight: "600", LineHeight: "1.2", LetterSpacing: "-0.02em", SourceID: "awwwards"},
			{Role: "body", FontFamily: "Plus Jakarta Sans", Fallback: "sans-serif", SizeClamp: "clamp(1rem, 1vw + 0.5rem, 1.125rem)", Weight: "400", LineHeight: "1.6", LetterSpacing: "0", SourceID: "awwwards"},
			{Role: "mono", FontFamily: "JetBrains Mono", Fallback: "monospace", SizeClamp: "0.875rem", Weight: "500", LineHeight: "1.4", LetterSpacing: "-0.01em", SourceID: "awwwards"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "hero_composition", Description: "Asymmetric 60/40 split with massive display headline and offset floating artifact canvas", SourceID: "awwwards", SourceURL: "https://www.awwwards.com/"},
			{Category: "surface_depth", Description: "Subtle 1px inset borders with 4% white alpha over dark matte surface; no blurry drop-shadows", SourceID: "awwwards", SourceURL: "https://www.awwwards.com/"},
		},
		Citations: []string{"https://www.awwwards.com/websites/", "https://www.awwwards.com/sites/editorial-showcase"},
	},
	"jiro-design": {
		SourceID:   "jiro-design",
		SourceName: "Jiro Design Direction",
		SourceURL:  "https://jiro.design/",
		Category:   "creative-layouts",
		KeyTakeaways: []string{
			"Non-standard viewport geometry with bold structural borders and negative whitespace anchors",
			"Micro-interactions should use tactile spring physics (stiffness: 300, damping: 30) instead of linear transitions",
			"Active press states must feel physical with scale(0.98) feedback",
			"Strictly prohibit animating layout-triggering properties (width, height, top, padding)",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0D1117", HSL: "hsl(215, 28%, 7%)", Contrast: "16.8:1 against text", SourceID: "jiro-design"},
			{Role: "--color-surface-elevated", Hex: "#161B22", HSL: "hsl(215, 21%, 11%)", Contrast: "14.2:1 against text", SourceID: "jiro-design"},
			{Role: "--color-accent-primary", Hex: "#E05A47", HSL: "hsl(7, 72%, 58%)", Contrast: "5.1:1 against bg", SourceID: "jiro-design"},
			{Role: "--color-border-subtle", Hex: "#30363D", HSL: "hsl(215, 12%, 21%)", Contrast: "3.1:1 decorative", SourceID: "jiro-design"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "display", FontFamily: "Playfair Display", Fallback: "serif", SizeClamp: "clamp(2.5rem, 5vw + 1rem, 5rem)", Weight: "500", LineHeight: "1.1", LetterSpacing: "-0.025em", SourceID: "jiro-design"},
			{Role: "body", FontFamily: "Geist Sans", Fallback: "sans-serif", SizeClamp: "1rem", Weight: "400", LineHeight: "1.55", LetterSpacing: "-0.01em", SourceID: "jiro-design"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "grid_structure", Description: "Continuous thin grid lines defining structural panels with corner crosshair accents", SourceID: "jiro-design", SourceURL: "https://jiro.design/"},
			{Category: "interaction", Description: "Immediate hover reveal with 0.15s ease-out transform and border color shift", SourceID: "jiro-design", SourceURL: "https://jiro.design/"},
		},
		Citations: []string{"https://jiro.design/", "https://jiro.design/archive"},
	},
	"cari-institute": {
		SourceID:   "cari-institute",
		SourceName: "Consumer Aesthetics Research Institute (CARI)",
		SourceURL:  "https://cari.institute/",
		Category:   "aesthetic-taxonomy",
		KeyTakeaways: []string{
			"Avoid generic modern AI tropes (neon gradient bubbles, purple blobs, ungrounded floating spheres)",
			"Contextualize visual language within authentic historical design movements (Swiss International, Cyberpunk, Neo-Brutalism)",
			"High-density information design demands tabular figures and rigorous mathematical alignment",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0A0D12", HSL: "hsl(218, 29%, 5%)", Contrast: "18.0:1 against text", SourceID: "cari-institute"},
			{Role: "--color-accent-primary", Hex: "#00E5FF", HSL: "hsl(186, 100%, 50%)", Contrast: "12.5:1 against bg", SourceID: "cari-institute"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "display", FontFamily: "Cabinet Grotesk", Fallback: "sans-serif", SizeClamp: "clamp(2.5rem, 5vw + 1rem, 4.75rem)", Weight: "800", LineHeight: "1.0", LetterSpacing: "-0.04em", SourceID: "cari-institute"},
			{Role: "mono", FontFamily: "Fira Code", Fallback: "monospace", SizeClamp: "0.85rem", Weight: "400", LineHeight: "1.5", LetterSpacing: "0", SourceID: "cari-institute"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "editorial_taxonomy", Description: "Typographic classification index tags with monospace bracket notations [01/REF]", SourceID: "cari-institute", SourceURL: "https://cari.institute/"},
		},
		Citations: []string{"https://cari.institute/aesthetics", "https://cari.institute/movement-index"},
	},
	"awesome-design-md": {
		SourceID:   "awesome-design-md",
		SourceName: "Awesome DESIGN.md Repository",
		SourceURL:  "https://github.com/VoltAgent/awesome-design-md",
		Category:   "design-tokens",
		KeyTakeaways: []string{
			"Machine-verifiable DESIGN.md contracts lock typography, color palette, and spring constants before coding",
			"WCAG AAA contrast (>7:1) for body text and AA (>4.5:1) for small labels prevents accessibility failures",
			"Mobile viewports (<390px) must strictly maintain zero horizontal overflow (scrollWidth <= innerWidth)",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0E1116", HSL: "hsl(218, 22%, 7%)", Contrast: "16.5:1 against text", SourceID: "awesome-design-md"},
			{Role: "--color-surface-elevated", Hex: "#171C24", HSL: "hsl(217, 22%, 12%)", Contrast: "13.9:1 against text", SourceID: "awesome-design-md"},
			{Role: "--color-text-headline", Hex: "#F2F5F8", HSL: "hsl(210, 30%, 96%)", Contrast: "16.5:1 against bg", SourceID: "awesome-design-md"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "display", FontFamily: "Fraunces", Fallback: "serif", SizeClamp: "clamp(2.5rem, 5vw + 1rem, 5rem)", Weight: "600", LineHeight: "1.1", LetterSpacing: "-0.02em", SourceID: "awesome-design-md"},
			{Role: "body", FontFamily: "Inter Tight", Fallback: "sans-serif", SizeClamp: "1rem", Weight: "400", LineHeight: "1.5", LetterSpacing: "0", SourceID: "awesome-design-md"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "token_architecture", Description: "Semantic token hierarchy mapping core primitive variables to scoped component variables", SourceID: "awesome-design-md", SourceURL: "https://github.com/VoltAgent/awesome-design-md"},
		},
		Citations: []string{"https://github.com/VoltAgent/awesome-design-md", "https://github.com/VoltAgent/awesome-design-md/blob/main/DESIGN.md"},
	},
	"godly-design": {
		SourceID:   "godly-design",
		SourceName: "Godly Design Inspiration",
		SourceURL:  "https://godly.design/",
		Category:   "award-winning-ui",
		KeyTakeaways: []string{
			"High-density minimal hero layouts with micro-tag badges and prominent typography",
			"Sophisticated motion pacing with staggered entrance delays (50ms per item)",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0C0E12", HSL: "hsl(220, 14%, 6%)", Contrast: "17.0:1", SourceID: "godly-design"},
			{Role: "--color-accent-primary", Hex: "#3B82F6", HSL: "hsl(217, 91%, 60%)", Contrast: "4.9:1", SourceID: "godly-design"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "display", FontFamily: "Instrument Serif", Fallback: "serif", SizeClamp: "clamp(2.5rem, 5.5vw, 5.2rem)", Weight: "400", LineHeight: "1.08", LetterSpacing: "-0.03em", SourceID: "godly-design"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "hero_composition", Description: "Center-aligned stark display title with dual secondary CTA chips", SourceID: "godly-design", SourceURL: "https://godly.design/"},
		},
		Citations: []string{"https://godly.design/"},
	},
	"refero-design": {
		SourceID:   "refero-design",
		SourceName: "Refero Design Patterns",
		SourceURL:  "https://styles.refero.design/",
		Category:   "product-ui",
		KeyTakeaways: []string{
			"Production SaaS UI patterns prioritizing cognitive clarity and rapid task completion",
			"Tabular alignment for metrics cards with delta badges (+12.4%) in monospace figures",
		},
		ExtractedPalettes: []PaletteToken{
			{Role: "--color-bg-base", Hex: "#0F141C", HSL: "hsl(216, 29%, 8%)", Contrast: "15.8:1", SourceID: "refero-design"},
			{Role: "--color-surface-elevated", Hex: "#1A222F", HSL: "hsl(216, 29%, 15%)", Contrast: "12.5:1", SourceID: "refero-design"},
		},
		ExtractedTypography: []TypographyToken{
			{Role: "headline", FontFamily: "Plus Jakarta Sans", Fallback: "sans-serif", SizeClamp: "1.5rem", Weight: "600", LineHeight: "1.3", LetterSpacing: "-0.02em", SourceID: "refero-design"},
			{Role: "mono", FontFamily: "JetBrains Mono", Fallback: "monospace", SizeClamp: "0.875rem", Weight: "500", LineHeight: "1.4", LetterSpacing: "0", SourceID: "refero-design"},
		},
		VisualMotifs: []VisualMotif{
			{Category: "kpi_card", Description: "Compact metric cards with muted label, high-contrast numeric value, and trend indicator", SourceID: "refero-design", SourceURL: "https://styles.refero.design/"},
		},
		Citations: []string{"https://styles.refero.design/"},
	},
}

package visual

import (
	"fmt"
	"strings"
)

// Direction is a lightweight flagship exploration artifact. It is not a DESIGN.md
// and not a 23-card survey card.
type Direction struct {
	ID                         string   `json:"id"`
	Concept                    string   `json:"concept"`
	VisualGrammar              string   `json:"visual_grammar"`
	Typography                 string   `json:"typography"`
	Composition                string   `json:"composition"`
	Interaction                string   `json:"interaction"`
	Motion                     string   `json:"motion,omitempty"`
	Strengths                  []string `json:"strengths,omitempty"`
	Risks                      []string `json:"risks,omitempty"`
	ImplementationImplications string   `json:"implementation_implications,omitempty"`
}

type SealedContract struct {
	DirectionID string    `json:"direction_id"`
	Direction   Direction `json:"direction"`
	Sealed      bool      `json:"sealed"`
}

type Exploration struct {
	Class      Class           `json:"class"`
	Candidates []Direction     `json:"candidates,omitempty"`
	SelectedID string          `json:"selected_id,omitempty"`
	Sealed     *SealedContract `json:"sealed,omitempty"`
	Waived     bool            `json:"waived,omitempty"`
}

func ValidateExploration(class Class, dirs []Direction) error {
	if class != ClassFlagship {
		if len(dirs) > 0 {
			return fmt.Errorf("direction exploration is flagship-only; %s work must not generate 3-5 directions", class)
		}
		return nil
	}
	n := len(dirs)
	if n < 3 || n > 5 {
		return fmt.Errorf("flagship work requires 3-5 lightweight directions, got %d", n)
	}
	for i, d := range dirs {
		if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.Concept) == "" {
			return fmt.Errorf("direction %d is missing id or concept", i)
		}
		if strings.TrimSpace(d.VisualGrammar) == "" || strings.TrimSpace(d.Composition) == "" {
			return fmt.Errorf("direction %q is missing visual grammar or composition", d.ID)
		}
	}
	return structurallyDistinct(dirs)
}

func structurallyDistinct(dirs []Direction) error {
	type key struct{ grammar, typeDir, composition, interaction string }
	seen := map[key]string{}
	cosmeticOnly := true
	for _, d := range dirs {
		k := key{
			grammar:     normalize(d.VisualGrammar),
			typeDir:     normalize(d.Typography),
			composition: normalize(d.Composition),
			interaction: normalize(d.Interaction),
		}
		if prev, ok := seen[k]; ok {
			return fmt.Errorf("directions %q and %q differ only cosmetically (color/gradient/radius/type tweaks are not enough)", prev, d.ID)
		}
		seen[k] = d.ID
		if !cosmeticLanguage(d) {
			cosmeticOnly = false
		}
	}
	if cosmeticOnly {
		return fmt.Errorf("directions must differ structurally (composition, grammar, interaction, density, motion), not merely by color or radius")
	}
	return nil
}

func cosmeticLanguage(d Direction) bool {
	blob := strings.ToLower(d.Concept + " " + d.VisualGrammar + " " + d.Composition + " " + d.Interaction)
	structural := containsAny(blob, []string{
		"asymmetric", "split", "editorial", "spatial", "scroll", "immersive",
		"dense", "sparse", "navigat", "grid", "column", "overlay", "stage",
		"interaction", "motion", "typograph", "hierarchy",
	})
	return !structural
}

func SelectOne(dirs []Direction, selectedID string) (*Direction, error) {
	if strings.TrimSpace(selectedID) == "" {
		return nil, fmt.Errorf("direction selection is a human/design decision; no id provided")
	}
	var found *Direction
	for i := range dirs {
		if dirs[i].ID == selectedID {
			found = &dirs[i]
			break
		}
	}
	if found == nil {
		return nil, fmt.Errorf("selected direction %q is not in the exploration set", selectedID)
	}
	return found, nil
}

func Seal(selected Direction) (*SealedContract, error) {
	if strings.TrimSpace(selected.ID) == "" {
		return nil, fmt.Errorf("cannot seal an empty direction")
	}
	return &SealedContract{
		DirectionID: selected.ID,
		Direction:   selected,
		Sealed:      true,
	}, nil
}

func Mutated(sealed *SealedContract, next Direction) bool {
	if sealed == nil || !sealed.Sealed {
		return false
	}
	return sealed.DirectionID != next.ID ||
		normalize(sealed.Direction.Composition) != normalize(next.Composition) ||
		normalize(sealed.Direction.VisualGrammar) != normalize(next.VisualGrammar)
}

// DefaultFlagshipSet returns four structurally different directions for tests
// and for pipeline fixtures. Conductors may replace these with brief-specific ones.
func DefaultFlagshipSet(objective string) []Direction {
	_ = objective
	return []Direction{
		{
			ID:                         "dir-editorial-split",
			Concept:                    "Editorial split-stage: product artifact owns one half, type the other.",
			VisualGrammar:              "Asymmetric two-pane studio with a single hairline and one accent.",
			Typography:                 "Display serif headlines, grotesque body, tabular figures.",
			Composition:                "Left artifact stage / right stacked narrative. Not a centered hero.",
			Interaction:                "Hover reveals the artifact's working parts; click drills one scene.",
			Motion:                     "Entrance on the artifact only; reduced-motion falls back to static crop.",
			Strengths:                  []string{"clear identity", "product-specific anchor"},
			Risks:                      []string{"needs a real artifact, not a stock render"},
			ImplementationImplications: "One hero canvas plus a reading column; no card grid.",
		},
		{
			ID:                         "dir-dense-console",
			Concept:                    "Operator console: dense instrument panel as the brand, not a marketing page.",
			VisualGrammar:              "High-density modular instruments, monospace readouts, restrained chrome.",
			Typography:                 "Grotesque UI + monospace telemetry. No display serif.",
			Composition:                "Full-bleed work surface with a slim command rail. Hierarchy via weight and grid, not hero type.",
			Interaction:                "Keyboard-first commands; pointer is secondary.",
			Motion:                     "Value ticks only. No page-level parallax.",
			Strengths:                  []string{"distinct from SaaS marketing", "interaction-specific"},
			Risks:                      []string{"can read as admin HUD if type is generic"},
			ImplementationImplications: "CSS grid instruments; skip marketing sections.",
		},
		{
			ID:                         "dir-scroll-essay",
			Concept:                    "Long-scroll essay: chapters of the product story with pinned scenes.",
			VisualGrammar:              "Vertical chapters, large type, one illustration language.",
			Typography:                 "High-contrast serif with a narrow sans for UI chrome.",
			Composition:                "Single column of chapters; each chapter has a distinct spatial rule.",
			Interaction:                "Scroll-tied scene changes; chapter skip nav.",
			Motion:                     "Pinned scene swaps. Honor prefers-reduced-motion by stacking stills.",
			Strengths:                  []string{"narrative identity", "motion with a job"},
			Risks:                      []string{"length; mobile chapter collapse must be designed"},
			ImplementationImplications: "Scroll-driven sections; no 3-card feature row.",
		},
		{
			ID:                         "dir-spatial-object",
			Concept:                    "Spatial object: the product is a 3D/physical thing you orbit, then read.",
			VisualGrammar:              "Object-first stage, sparse UI, material lighting.",
			Typography:                 "Small UI captions; the object carries identity, not the font.",
			Composition:                "Centered object with offset captions. Empty space is structural.",
			Interaction:                "Drag to orbit; click a part to open a spec sheet.",
			Motion:                     "Inertia on orbit; captions fade, they do not bounce.",
			Strengths:                  []string{"unmistakable product metaphor"},
			Risks:                      []string{"needs a real object and a low-end fallback"},
			ImplementationImplications: "R3F/Spline object plus a fallback still; captions as HTML overlay.",
		},
	}
}

func normalize(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(s))), " ")
}

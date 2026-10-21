package visual

import "strings"

type Finding struct {
	Type         string `json:"type"`
	Layer        string `json:"layer"` // A | B | critic | judge
	Surface      string `json:"surface,omitempty"`
	Route        string `json:"route,omitempty"`
	Viewport     string `json:"viewport,omitempty"`
	Explanation  string `json:"explanation"`
	Action       string `json:"action"`
	EvidencePath string `json:"evidence_path,omitempty"`
	Trivial      bool   `json:"trivial,omitempty"`
}

type CriticInput struct {
	Class           Class
	Surface         string
	Route           string
	Viewport        string
	Description     string
	HasGradient     bool
	HasRounded      bool
	HasDark         bool
	HasCommonFont   bool
	HasFamiliarUI   bool
	CenteredHero    bool
	NoProductAnchor bool
	GenericType     bool
	NoHierarchy     bool
	NoSpatialPlan   bool
	NoInteraction   bool
	ArbitraryMotion bool
	ScreenshotPath  string
}

// Critique is an actionable anti-slop critic. Gradients, rounded cards, dark
// surfaces, common fonts, and familiar UI patterns are not automatic fails.
func Critique(in CriticInput) []Finding {
	if in.Class == ClassOrdinary {
		return nil
	}
	var out []Finding
	base := Finding{Layer: "critic", Surface: in.Surface, Route: in.Route, Viewport: in.Viewport, EvidencePath: in.ScreenshotPath}

	if in.CenteredHero && in.NoProductAnchor {
		f := base
		f.Type = "composition"
		f.Explanation = "The hero uses a standard centered headline + CTA composition with no product-specific visual anchor."
		f.Action = "Introduce a distinctive spatial or interactive element that communicates the product's core concept instead of a generic centered stack."
		out = append(out, f)
	}
	if in.GenericType {
		f := base
		f.Type = "typography"
		f.Explanation = "Typography is present but not a deliberate pairing: display and body do not establish a product-specific voice."
		f.Action = "Specify a display + body pairing with a reason (editorial vs console vs object captions) and apply it to the rendered hero and chrome."
		out = append(out, f)
	}
	if in.NoHierarchy {
		f := base
		f.Type = "hierarchy"
		f.Explanation = "Multiple elements compete at the same visual weight; there is no one focal point a stranger would read first."
		f.Action = "Reduce competing weights so one product idea leads; demote secondary claims by size, contrast, or position."
		out = append(out, f)
	}
	if in.NoSpatialPlan {
		f := base
		f.Type = "spatial"
		f.Explanation = "Spacing looks like a default stack of sections rather than a planned spatial grammar."
		f.Action = "Define a spatial rule (split stage, dense console, chapter scroll, or object stage) and space sections from that rule, not from equal card padding."
		out = append(out, f)
	}
	if in.NoInteraction && in.Class == ClassFlagship {
		f := base
		f.Type = "interaction"
		f.Explanation = "The surface is visually decorated but has no interaction model that is specific to this product."
		f.Action = "Name one product-specific interaction (orbit, command, chapter skip, artifact reveal) and implement that, not generic hover chrome."
		out = append(out, f)
	}
	if in.ArbitraryMotion {
		f := base
		f.Type = "motion"
		f.Explanation = "Motion is present without a job: decoration rather than revealing hierarchy, state, or the product object."
		f.Action = "Keep only motion that reveals structure or state; drop unearned page-level motion."
		out = append(out, f)
	}

	// Presence of common techniques is recorded only when paired with missing decisions.
	if len(out) == 0 && (in.HasGradient || in.HasRounded || in.HasDark || in.HasCommonFont || in.HasFamiliarUI) {
		return nil
	}
	_ = strings.TrimSpace(in.Description)
	return out
}

func (f Finding) Signature() string {
	return strings.Join([]string{f.Type, f.Layer, f.Surface, f.Route, f.Viewport, f.Explanation}, "|")
}

func TrivialQA(f Finding) bool {
	if f.Trivial {
		return true
	}
	switch f.Type {
	case "overflow", "console", "network", "visibility", "geometry", "interaction":
		return f.Layer == "A"
	default:
		return false
	}
}

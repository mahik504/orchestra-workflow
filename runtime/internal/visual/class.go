package visual

import "strings"

// Design class is orthogonal to quality bar and to LIGHT/FULL intensity.
const (
	ClassOrdinary  Class = "ordinary"
	ClassProductUI Class = "product_ui"
	ClassFlagship  Class = "flagship"
)

type Class string

type ClassifyInput struct {
	Type           string
	Objective      string
	Tags           []string
	RequiresVisual bool
	// VisualImportance is an already-recorded mission field; empty means infer.
	VisualImportance string
}

type Policy struct {
	Class            Class
	Directions       bool
	LayerA           bool
	LayerB           bool
	Judge            string // off | conditional | required
	VisualReview     string // none | browser | deep
	IndependentJudge string
}

func Classify(in ClassifyInput) Class {
	if in.VisualImportance == "flagship" {
		return ClassFlagship
	}
	hay := strings.ToLower(in.Objective + " " + strings.Join(in.Tags, " "))
	if flagshipLanguage(hay) {
		return ClassFlagship
	}
	typ := strings.ToLower(strings.TrimSpace(in.Type))
	if in.VisualImportance == "product" || in.RequiresVisual || typ == "product_ui" || productUILanguage(hay) {
		return ClassProductUI
	}
	return ClassOrdinary
}

func PolicyFor(c Class) Policy {
	switch c {
	case ClassFlagship:
		return Policy{
			Class:            c,
			Directions:       true,
			LayerA:           true,
			LayerB:           true,
			Judge:            "conditional",
			VisualReview:     "deep",
			IndependentJudge: "conditional",
		}
	case ClassProductUI:
		return Policy{
			Class:            c,
			Directions:       false,
			LayerA:           true,
			LayerB:           false,
			Judge:            "conditional",
			VisualReview:     "browser",
			IndependentJudge: "off",
		}
	default:
		return Policy{
			Class:            ClassOrdinary,
			Directions:       false,
			LayerA:           false,
			LayerB:           false,
			Judge:            "off",
			VisualReview:     "none",
			IndependentJudge: "off",
		}
	}
}

func flagshipLanguage(hay string) bool {
	return containsAny(hay, []string{
		"portfolio",
		"flagship",
		"showcase",
		"award",
		"art-direct",
		"art directed",
		"brand experience",
		"hiring site",
		"unusual interactive",
		"visually distinctive",
		"visually novel",
		"3d",
		"webgl",
		"shader",
	})
}

func productUILanguage(hay string) bool {
	return containsAny(hay, []string{
		"dashboard",
		"settings screen",
		"settings page",
		"responsive",
		"landing",
		"frontend",
		"ui",
		"form",
		"page",
		"screen",
	})
}

func containsAny(hay string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(hay, n) {
			return true
		}
	}
	return false
}

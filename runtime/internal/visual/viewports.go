package visual

import "strings"

type Viewport struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

const (
	ViewportDesktop = "desktop"
	ViewportMobile  = "mobile"
	ViewportTablet  = "tablet"
)

func DefaultViewports(includeTablet bool) []Viewport {
	out := []Viewport{
		{Name: ViewportDesktop, Width: 1440, Height: 900},
		{Name: ViewportMobile, Width: 390, Height: 844},
	}
	if includeTablet {
		out = append(out, Viewport{Name: ViewportTablet, Width: 768, Height: 1024})
	}
	return out
}

func NeedsTablet(objective string, parameters map[string]any) bool {
	if parameters != nil {
		if v, ok := parameters["tablet"]; ok {
			switch t := v.(type) {
			case bool:
				return t
			case string:
				return t == "true" || t == "1" || t == "yes"
			}
		}
		if v, ok := parameters["viewports"]; ok {
			if s, ok := v.(string); ok && strings.Contains(strings.ToLower(s), "tablet") {
				return true
			}
		}
	}
	return strings.Contains(strings.ToLower(objective), "tablet")
}

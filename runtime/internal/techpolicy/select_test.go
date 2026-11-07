package techpolicy

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestStacksStaySmall(t *testing.T) {
	file := loadPolicy(t)
	if len(file.Libraries) < 20 {
		t.Fatalf("libraries = %d", len(file.Libraries))
	}
	cases := []struct {
		name         string
		text         string
		capability   string
		architecture string
		want         []string
		reject       []string
	}{
		{
			name:         "static",
			text:         "A semantic HTML static page with no framework",
			capability:   "premium-website",
			architecture: "html-css",
			want:         []string{"html-css"},
			reject:       []string{"react", "next", "gsap", "r3f", "redis"},
		},
		{
			name:         "scroll",
			text:         "A landing page with scroll-driven pinned sections",
			capability:   "premium-website",
			architecture: "unspecified",
			want:         []string{"gsap"},
			reject:       []string{"motion-react", "lenis", "r3f", "framer-site"},
		},
		{
			name:         "hover",
			text:         "Add a hover transition and a drag gesture in React",
			capability:   "premium-website",
			architecture: "react-typescript",
			want:         []string{"motion-react", "react", "typescript"},
			reject:       []string{"gsap", "lenis", "framer-site"},
		},
		{
			name:         "scene",
			text:         "Build a 3D WebGL portfolio with R3F camera orbits and a custom shader",
			capability:   "3d-portfolio",
			architecture: "r3f",
			want:         []string{"r3f", "three", "react", "typescript", "drei", "glsl"},
			reject:       []string{"theatre", "react-postprocessing", "r3f-scroll-rig", "lenis", "framer-site", "redis"},
		},
		{
			name:         "saas",
			text:         "Build a full-stack SaaS with accounts and billing",
			capability:   "saas-dashboard",
			architecture: "next",
			want:         []string{"next", "react", "typescript", "postgres"},
			reject:       []string{"redis", "bullmq", "supabase", "fastapi", "express", "vite"},
		},
		{
			name:         "fastapi",
			text:         "An inference API in FastAPI for a model",
			capability:   "b2b-portal",
			architecture: "fastapi",
			want:         []string{"fastapi"},
			reject:       []string{"express", "next", "node", "redis", "react"},
		},
		{
			name:         "framer",
			text:         "The deliverable is a Framer site",
			capability:   "premium-website",
			architecture: "framer-site",
			want:         []string{"framer-site"},
			reject:       []string{"react", "gsap", "motion-react", "next"},
		},
		{
			name:         "bug",
			text:         "Fix Sum. The test expects a+b.",
			capability:   "bug",
			architecture: "unchanged",
			reject:       []string{"react", "postgres", "gsap", "fastapi"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Select(file, tc.text, tc.capability)
			if got.Conflict != "" {
				t.Fatalf("conflict: %s", got.Conflict)
			}
			if got.Architecture != tc.architecture {
				t.Fatalf("architecture %s, selected %v", got.Architecture, got.Selected)
			}
			for _, id := range tc.want {
				if !has(got.Selected, id) {
					t.Fatalf("missing %s in %v", id, got.Selected)
				}
			}
			for _, id := range tc.reject {
				if has(got.Selected, id) {
					t.Fatalf("selected %s in %v", id, got.Selected)
				}
			}
			if len(got.Selected) >= len(file.Libraries) {
				t.Fatal("selected every library")
			}
		})
	}
}

func TestScrollAndHoverCanShareAPage(t *testing.T) {
	file := loadPolicy(t)
	got := Select(file, "Scroll-driven chapters plus a hover transition in React", "premium-website")
	if !has(got.Selected, "gsap") || !has(got.Selected, "motion-react") {
		t.Fatalf("selected %v", got.Selected)
	}
	if len(got.Constraints) == 0 {
		t.Fatal("expected a scroll-owner constraint")
	}
}

func TestTwoHTTPStacksConflict(t *testing.T) {
	file := loadPolicy(t)
	got := Select(file, "Use Express and FastAPI for the same API", "b2b-portal")
	if got.Conflict == "" || got.Architecture != "unspecified" {
		t.Fatalf("conflict %q architecture %s", got.Conflict, got.Architecture)
	}
}

func loadPolicy(t *testing.T) File {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	root := filepath.Join(filepath.Dir(here), "..", "..", "..")
	file, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func has(list []string, id string) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}

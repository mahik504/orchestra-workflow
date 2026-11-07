package routing

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReceiptNamesTheDecision(t *testing.T) {
	root := repoRoot(t)
	text, err := Explain(root, Input{
		Text:        "A landing page with scroll-driven pinned sections",
		Capability:  "premium-website",
		Task:        "flagship",
		RouteSkills: []string{"impeccable", "taste-design", "gsap-core", "gsap-scrolltrigger", "gsap-timeline"},
		NamedSkills: []string{"systematic-debugging", "missing-skill"},
		KnownSkills: []string{"impeccable", "taste-design", "gsap-core", "gsap-scrolltrigger", "gsap-timeline", "systematic-debugging", "ow-3d-engineering"},
		SelectedMCP: []string{"playwright"},
		Plugins:     nil,
		Acceptance:  []string{"the pinned section scrubs and the page still reads without motion"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Load: impeccable, taste-design, gsap-core, gsap-scrolltrigger, gsap-timeline, systematic-debugging") {
		t.Fatalf("named skill did not win\n%s", text)
	}
	if !strings.Contains(text, "Missing: missing-skill (missing; do not install globally; wait for approval)") {
		t.Fatalf("missing skill was not reported\n%s", text)
	}
	if !strings.Contains(text, "Libraries: gsap") {
		t.Fatalf("library\n%s", text)
	}
	if strings.Contains(text, "lenis") {
		t.Fatal("lenis was selected without being asked")
	}
	if strings.Contains(text, "Model capability") || strings.Contains(text, "Skip:") || strings.Contains(text, "Quality gates") {
		t.Fatalf("empty categories were printed\n%s", text)
	}
	if !strings.Contains(text, "MCP: playwright") {
		t.Fatal("mcp")
	}
}

func TestNamedSceneSkillLosesToTheSceneOwner(t *testing.T) {
	root := repoRoot(t)
	got, err := Build(root, Input{
		Text:        "One R3F scene",
		Capability:  "3d-portfolio",
		Task:        "experimental",
		RouteSkills: []string{"ow-3d-engineering"},
		NamedSkills: []string{"creative-3d-web"},
		KnownSkills: []string{"ow-3d-engineering", "creative-3d-web"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if contains(got.Load, "creative-3d-web") {
		t.Fatal("creative-3d-web loaded beside ow-3d-engineering")
	}
	found := false
	for _, line := range got.Skip {
		if strings.Contains(line, "creative-3d-web") {
			found = true
		}
	}
	if !found {
		t.Fatalf("skip %v", got.Skip)
	}
}

func TestBugReceiptStaysSmall(t *testing.T) {
	root := repoRoot(t)
	got, err := Build(root, Input{
		Text:        "Fix Sum. The test expects a+b.",
		Capability:  "premium-website",
		Task:        "bug",
		SelectedMCP: []string{"playwright"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Libraries) != 0 {
		t.Fatalf("libraries %v", got.Libraries)
	}
	if len(got.Documents) != 0 {
		t.Fatalf("documents %v", got.Documents)
	}
	if len(got.Gates) != 1 || got.Gates[0] != "tests" {
		t.Fatalf("gates %v", got.Gates)
	}
	if got.Architecture != "unchanged" {
		t.Fatal(got.Architecture)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Join(filepath.Dir(here), "..", "..", "..")
}

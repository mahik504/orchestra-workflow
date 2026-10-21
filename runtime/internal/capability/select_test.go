package capability

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func TestSelect_OrdinaryOmitsResearchAndBrowser(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "ord", Objective: "Fix nil pointer dereference in parser", DeclaredType: "BUGFIX"})
	pw := false
	d, err := Select(Input{Mission: m, PlaywrightAvailable: &pw, Skills: []string{"superpowers-planning", "gsap-core"}})
	if err != nil {
		t.Fatal(err)
	}
	if !d.Has(IDCodingWorker) {
		t.Fatal("coding.worker required")
	}
	if d.Has(IDResearchWorker) || d.Has(IDBrowserPlaywright) || d.Has(IDDesignWorker) || d.Has(IDJudgeIndependent) {
		t.Fatalf("ordinary selected extra: %+v", d.Selected)
	}
	for _, s := range d.Skills {
		if s == "gsap-core" {
			t.Fatal("visual skill leaked onto ordinary")
		}
	}
}

func TestSelect_BrowserUnavailableHumanStop(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "ui", Objective: "Build a settings screen for profile CRUD", DeclaredType: "FEATURE"})
	pw := false
	d, err := Select(Input{Mission: m, PlaywrightAvailable: &pw})
	if !IsHumanStop(err) {
		t.Fatalf("want human stop, got %v", err)
	}
	if d.HumanDetail == "" || !strings.Contains(d.HumanDetail, IDBrowserPlaywright) {
		t.Fatalf("human detail: %s", d.HumanDetail)
	}
}

func TestSelect_MultimodalFallsBackToStructural(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "flag", Objective: "Build a flagship portfolio landing with unusual motion",
		DeclaredType: "DESIGN", QualityBar: mission.BarPremium,
	})
	pw := true
	shot := false
	d, err := Select(Input{Mission: m, PlaywrightAvailable: &pw, ScreenshotEvalAvailable: &shot})
	if err != nil {
		t.Fatal(err)
	}
	if !d.Has(IDVisualStructural) {
		t.Fatal("structural must be selected")
	}
	if d.Has(IDVisualMultimodal) {
		t.Fatal("unavailable multimodal must not be selected")
	}
	if d.Has(IDVisualScreenshot) {
		t.Fatal("screenshot must not be selected without adapter")
	}
	if len(d.Fallbacks) == 0 {
		t.Fatal("fallback hops required")
	}
	found := false
	for _, f := range d.Fallbacks {
		if f.Requested == IDVisualMultimodal && f.Selected == IDVisualStructural {
			found = true
		}
	}
	if !found {
		t.Fatalf("multimodal→structural hop missing: %+v", d.Fallbacks)
	}
}

func TestSelect_OptionalBrokenMCPSkipped(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "ord", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	d, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		MCP: []MCPCandidate{{ID: "some-optional", Health: "BROKEN"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range d.Tools {
		if id == "some-optional" {
			t.Fatal("broken optional MCP selected")
		}
	}
}

func TestSelect_RequiredAuthMCPHumanStop(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "ord", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	d, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		MCP: []MCPCandidate{{ID: "needed-tool", Health: "AUTH_REQUIRED", Required: true}},
	})
	if !IsHumanStop(err) {
		t.Fatalf("want human stop, got %v", err)
	}
	if !strings.Contains(d.HumanDetail, "needed-tool") {
		t.Fatalf("detail %s", d.HumanDetail)
	}
}

func TestSelect_HostPreferenceSelectionOnly(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "h", Objective: "Fix parser", DeclaredType: "BUGFIX"})
	pw := false
	both, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		Hosts: []HostSnapshot{{ID: "cursor", Available: true}, {ID: "claude", Available: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if both.HostLabel != "cursor" {
		t.Fatalf("host %s", both.HostLabel)
	}
	fallback, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		Hosts: []HostSnapshot{{ID: "cursor", Available: false}, {ID: "claude", Available: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if fallback.HostLabel != "claude" {
		t.Fatalf("host %s", fallback.HostLabel)
	}
}

func TestSelect_OpaqueDiscoveredModelIgnored(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "m", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	a, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		DiscoveredModels: []string{"grok-4.7"},
		Skills:           []string{"superpowers-planning"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		DiscoveredModels: []string{"other-model-9"},
		Skills:           []string{"superpowers-planning"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ids(a) != ids(b) {
		t.Fatalf("opaque model label changed selection: %s vs %s", ids(a), ids(b))
	}
	blob, _ := json.Marshal(a)
	if strings.Contains(strings.ToLower(string(blob)), "grok") {
		t.Fatal("decision leaked opaque model label")
	}
}

func TestSelect_HostileCatalogCapped(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "ord", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	var skills []string
	for i := 0; i < 44; i++ {
		skills = append(skills, "skill-pack-"+strconv.Itoa(i))
	}
	skills = append(skills, "gsap-core", "impeccable", "taste-design")
	var resources []string
	for i := 0; i < 80; i++ {
		resources = append(resources, "resource-"+strconv.Itoa(i))
	}
	var mcp []MCPCandidate
	for i := 0; i < 20; i++ {
		mcp = append(mcp, MCPCandidate{ID: "mcp-tool-" + strconv.Itoa(i), Health: "HEALTHY"})
	}
	d, err := Select(Input{Mission: m, PlaywrightAvailable: &pw, Skills: skills, Resources: resources, MCP: mcp})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Skills) > MaxSkills || len(d.Resources) > MaxResources || len(d.Tools) > MaxMCPTools {
		t.Fatalf("caps exceeded skills=%d resources=%d tools=%d", len(d.Skills), len(d.Resources), len(d.Tools))
	}
	joined := strings.Join(d.Skills, ",") + strings.Join(d.Resources, ",") + strings.Join(d.Tools, ",")
	if strings.Contains(joined, "gsap-core") {
		t.Fatal("visual dump survived ordinary filter")
	}
	if len(d.Skills) == 44 || strings.Count(joined, "skill-pack-") == 44 {
		t.Fatal("full skill dump survived")
	}
}

func TestSelect_SourceHasNoVendorLiterals(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		low := strings.ToLower(string(b))
		for _, tok := range []string{"grok", "jcode", "workeradapter", "worker.optional", "gemini", "gpt-", "opus", "sonnet", "openai"} {
			if strings.Contains(low, tok) {
				t.Errorf("%s contains %q", name, tok)
			}
		}
	}
}

func TestWriteDecision_BesideMission(t *testing.T) {
	dir := t.TempDir()
	d := Decision{MissionID: "cap-write", HostLabel: "cursor"}
	path, err := WriteDecision(dir, "cap-write", d)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "capability.json" {
		t.Fatalf("path %s", path)
	}
	got, err := ReadDecision(dir, "cap-write")
	if err != nil {
		t.Fatal(err)
	}
	if got.HostLabel != "cursor" {
		t.Fatalf("%+v", got)
	}
}

func ids(d Decision) string {
	var s []string
	for _, it := range d.Selected {
		if it.Category == CatSkill || it.Category == CatResource || it.Category == CatMCPTool {
			continue
		}
		s = append(s, it.ID)
	}
	return strings.Join(s, ",")
}

func TestSelect_JudgeFollowsDecideJudge(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "j", Objective: "Fix parser", DeclaredType: "BUGFIX"})
	pw := false
	off, err := Select(Input{Mission: m, PlaywrightAvailable: &pw})
	if err != nil {
		t.Fatal(err)
	}
	if off.Has(IDJudgeIndependent) {
		t.Fatal("ordinary judge must stay off")
	}
	on, err := Select(Input{
		Mission: m, PlaywrightAvailable: &pw,
		Judge: &visual.JudgeInput{Class: visual.ClassFlagship, Policy: "required", SecurityRelease: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !on.Has(IDJudgeIndependent) {
		t.Fatal("required judge must be selected")
	}
}

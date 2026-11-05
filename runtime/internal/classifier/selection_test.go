package classifier

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

type selCase struct {
	bucket      string
	name        string
	request     string
	wantCap     string
	wantSkills  []string
	snapshot    []capability.MCPCandidate
	wantMCP     []string
	wantTools   []string
	exclude     []string
	fallback    string
	humanGate   bool
	extraSkills []string
}

func TestSelectionSuite(t *testing.T) {
	g := loadRealGraph(t)
	cat := loadRealCatalog(t)
	cls := NewClassifierWithGraph(g)
	pw := true
	shot := false

	cases := []selCase{
		{bucket: "engineering", name: "small api", request: "Add a small API sum that validates two integers and returns an error for invalid input. Include a test.", wantCap: "", wantSkills: []string{"systematic-debugging"}, exclude: []string{"interaction-components", "typescript-review", "deployed-errors", "slides", "gsap-core"}, extraSkills: []string{"gsap-core"}},
		{bucket: "engineering", name: "typescript boundary", request: "The job is a TypeScript review of types at a boundary", wantCap: "typescript-review", wantSkills: []string{"orchestra-typescript"}, snapshot: []capability.MCPCandidate{{ID: "context7", Health: "HEALTHY"}}, wantMCP: []string{"context7"}, wantTools: []string{"context7"}, exclude: []string{"playwright", "sentry"}},
		{bucket: "engineering", name: "zod edge", request: "The brief names Zod or a typed parser at the edge of the program", wantCap: "typescript-review", wantSkills: []string{"orchestra-typescript"}},
		{bucket: "engineering", name: "architecture diagram", request: "The deliverable is one architecture diagram", wantCap: "architecture-diagram", wantSkills: []string{"diagram-generator", "pretty-mermaid"}, exclude: []string{"orchestra-video", "gsap-core"}},
		{bucket: "engineering", name: "next speed", request: "The job is Next.js speed or a request waterfall", wantCap: "nextjs-speed", wantSkills: []string{"react-best-practices"}, snapshot: []capability.MCPCandidate{{ID: "context7", Health: "HEALTHY"}, {ID: "sentry", Health: "HEALTHY"}}, wantMCP: []string{"context7"}, wantTools: []string{"context7"}, exclude: []string{"sentry", "gsap-core"}},
		{bucket: "engineering", name: "injection audit", request: "Audit our API for injection and check the dependencies", wantCap: "security-audit", wantSkills: []string{"ship-safe"}, exclude: []string{"orchestra-observability"}},
		{bucket: "engineering", name: "simplify parser", request: "Rename the helper and simplify the token walk in the parser", wantCap: "", exclude: []string{"interaction-components", "premium-website", "slides"}},
		{bucket: "engineering", name: "school dashboard", request: "Build a scheduling dashboard for a school with attendance charts and KPIs", wantCap: "saas-dashboard", wantSkills: []string{"backend-architecture", "test-driven-development"}, exclude: []string{"impeccable", "orchestra-marketing"}},
		{bucket: "engineering", name: "package size", request: "The brief asks about a published package size", wantCap: "nextjs-speed", wantSkills: []string{"react-best-practices"}},
		{bucket: "engineering", name: "second diagram skip", request: "The request wants a second diagram system for the README", wantCap: "", exclude: []string{"architecture-diagram"}},

		{bucket: "frontend", name: "roastery", request: "Build a landing page for a friend's coffee roastery, they want it to feel expensive", wantCap: "premium-website", wantSkills: []string{"impeccable", "taste-design", "visual-golden-still", "gsap-core", "gsap-scrolltrigger", "gsap-timeline"}, snapshot: []capability.MCPCandidate{{ID: "playwright", Health: "HEALTHY"}, {ID: "stitch", Health: "HEALTHY"}, {ID: "context7", Health: "HEALTHY"}}, wantMCP: []string{"playwright", "stitch"}, wantTools: []string{"playwright", "stitch"}, exclude: []string{"context7", "orchestra-marketing"}},
		{bucket: "frontend", name: "r3f scene", request: "Build a 3D WebGL portfolio with R3F camera orbits and custom shaders", wantCap: "3d-portfolio", wantSkills: []string{"r3f-shaders"}, exclude: []string{"creative-3d-web", "gsap-core"}},
		{bucket: "frontend", name: "one control", request: "Add a magnetic gradient button using Kokonut after inspecting existing components", wantCap: "interaction-components", exclude: []string{"gsap-core", "impeccable"}},
		{bucket: "frontend", name: "cinematic", request: "Build an interactive 3D WebGL product UI with R3F, Three.js, and Drei — cinematic presentation", wantCap: "3d-portfolio", wantSkills: []string{"r3f-shaders"}, exclude: []string{"slides", "orchestra-presentation"}},
		{bucket: "frontend", name: "marketing site", request: "Build a premium landing page for a studio brand, award-winning marketing site", wantCap: "premium-website", wantSkills: []string{"impeccable", "taste-design", "visual-golden-still", "gsap-core", "gsap-scrolltrigger", "gsap-timeline"}, exclude: []string{"marketing-copy", "orchestra-marketing"}},

		{bucket: "research", name: "methods paper", request: "Write a methods paper with citations about capability graphs, no UI", wantCap: "research-paper", wantSkills: []string{"research-paper", "humanizer"}, snapshot: []capability.MCPCandidate{{ID: "context7", Health: "HEALTHY"}}, wantMCP: []string{"context7"}, wantTools: []string{"context7"}},
		{bucket: "research", name: "said recently", request: "The question is what people said recently in public", wantCap: "recent-research", wantSkills: []string{"orchestra-research"}, exclude: []string{"research-paper"}},
		{bucket: "research", name: "dated source", request: "Public discussions from the last month need a dated source", wantCap: "recent-research", wantSkills: []string{"orchestra-research"}},
		{bucket: "research", name: "thesis", request: "The job is producing a paper, thesis, or methods writeup", wantCap: "research-paper", wantSkills: []string{"research-paper", "humanizer"}, exclude: []string{"slides"}},
		{bucket: "research", name: "arxiv reader", request: "Build an arXiv reading app with a paper viewer and no writing", wantCap: "academic-reader", exclude: []string{"research-paper", "orchestra-research"}},

		{bucket: "documents", name: "slide deck", request: "The deliverable is a slide deck", wantCap: "slides", wantSkills: []string{"orchestra-docs", "orchestra-presentation"}, exclude: []string{"rendered-video", "research-paper"}},
		{bucket: "documents", name: "speaker notes", request: "The brief asks for speaker notes on slides", wantCap: "slides", wantSkills: []string{"orchestra-docs", "orchestra-presentation"}},
		{bucket: "documents", name: "tailor cv", request: "The job is to tailor a CV from notes already on disk", wantCap: "career-notes", wantSkills: []string{"orchestra-career"}, exclude: []string{"slides"}},
		{bucket: "documents", name: "page title", request: "The job is the page title or the on-page wording", wantCap: "marketing-copy", wantSkills: []string{"orchestra-marketing", "humanizer"}, exclude: []string{"premium-website", "impeccable"}},
		{bucket: "documents", name: "positioning", request: "The brief asks for a positioning line", wantCap: "marketing-copy", wantSkills: []string{"orchestra-marketing", "humanizer"}},

		{bucket: "debug", name: "login 500", request: "The login form throws a 500 when the email has a plus sign", wantCap: "", wantSkills: []string{"systematic-debugging"}, exclude: []string{"premium-website", "deployed-errors"}},
		{bucket: "debug", name: "supabase debug", request: "Debug a Supabase database issue in our schema", wantCap: "", wantSkills: []string{"systematic-debugging", "supabase"}, exclude: []string{"saas-dashboard"}},
		{bucket: "debug", name: "wrong total", request: "The report query is broken and returns the wrong total", wantCap: "", wantSkills: []string{"systematic-debugging"}},
		{bucket: "debug", name: "supabase notes", request: "Add Postgres indexes on a Supabase schema, no UI", wantSkills: []string{"supabase"}, exclude: []string{"premium-website", "3d-portfolio", "slides"}},
		{bucket: "debug", name: "local bug", request: "The job is a local bug with no deployment and the handler throws", wantCap: "", wantSkills: []string{"systematic-debugging"}, exclude: []string{"deployed-errors"}},

		{bucket: "deploy", name: "production error", request: "A deployed application is throwing a production error", wantCap: "deployed-errors", wantSkills: []string{"ship-safe", "orchestra-observability"}, snapshot: []capability.MCPCandidate{{ID: "sentry", Health: "HEALTHY"}}, wantMCP: []string{"sentry"}, wantTools: []string{"sentry"}, exclude: []string{"systematic-debugging", "playwright"}},
		{bucket: "deploy", name: "shipped app", request: "The brief names a production error in a shipped app", wantCap: "deployed-errors", wantSkills: []string{"ship-safe", "orchestra-observability"}},
		{bucket: "deploy", name: "showcase", request: "Record final-showcase brag walkthrough-video ship-assets after production-ready green ledger", wantCap: "final-showcase", wantSkills: []string{"orchestra-ship"}, exclude: []string{"rendered-video", "orchestra-video", "gsap-core"}},
		{bucket: "deploy", name: "rendered mp4", request: "The deliverable is a rendered mp4", wantCap: "rendered-video", wantSkills: []string{"orchestra-video"}, exclude: []string{"gsap-core", "final-showcase", "slides"}},
		{bucket: "deploy", name: "hyperframes", request: "The brief names HyperFrames for an mp4 file", wantCap: "rendered-video", wantSkills: []string{"orchestra-video"}},

		{bucket: "ambiguous", name: "cinematic not slides", request: "Build an interactive 3D WebGL product UI with R3F, Three.js, and Drei — cinematic presentation", wantCap: "3d-portfolio", wantSkills: []string{"r3f-shaders"}, exclude: []string{"slides"}},
		{bucket: "ambiguous", name: "award site not copy", request: "Build a premium landing page for a studio brand, award-winning marketing site", wantCap: "premium-website", exclude: []string{"marketing-copy"}},
		{bucket: "ambiguous", name: "paper not slides", request: "Write speaker notes on slides for a methods paper", wantCap: "research-paper", wantSkills: []string{"research-paper", "humanizer"}, exclude: []string{"slides"}},
		{bucket: "ambiguous", name: "scroll not video", request: "A showcase landing page with scroll-driven chapters", wantCap: "premium-website", exclude: []string{"rendered-video", "orchestra-video"}},
		{bucket: "ambiguous", name: "redesign not copy", request: "A stranger will see a visual redesign of this launch page, and it should feel like a distinctive look", wantCap: "premium-website", exclude: []string{"marketing-copy"}},

		{bucket: "unhealthy", name: "playwright broken", request: "Build a landing page for a friend's coffee roastery, they want it to feel expensive", wantCap: "premium-website", wantSkills: []string{"impeccable", "taste-design", "visual-golden-still", "gsap-core", "gsap-scrolltrigger", "gsap-timeline"}, snapshot: []capability.MCPCandidate{{ID: "playwright", Health: "BROKEN", Fallback: "playwright-cli"}}, wantMCP: []string{"playwright"}, fallback: "playwright-cli"},
		{bucket: "unhealthy", name: "playwright required", request: "Build a landing page for a friend's coffee roastery, they want it to feel expensive", wantCap: "premium-website", snapshot: []capability.MCPCandidate{{ID: "playwright", Health: "AUTH_REQUIRED", Required: true}}, wantMCP: []string{"playwright"}, humanGate: true},
		{bucket: "unhealthy", name: "context7 disabled", request: "The question is what people said recently in public", wantCap: "recent-research", wantSkills: []string{"orchestra-research"}, snapshot: []capability.MCPCandidate{{ID: "context7", Health: "DISABLED", Fallback: "use the docs already on disk"}}, wantMCP: []string{"context7"}, fallback: "use the docs already on disk"},
		{bucket: "unhealthy", name: "sentry broken", request: "A deployed application is throwing a production error", wantCap: "deployed-errors", wantSkills: []string{"ship-safe", "orchestra-observability"}, snapshot: []capability.MCPCandidate{{ID: "sentry", Health: "BROKEN", Fallback: "read the host log"}}, wantMCP: []string{"sentry"}, fallback: "read the host log"},
		{bucket: "unhealthy", name: "github not on this route", request: "The job is a TypeScript review of types at a boundary", wantCap: "typescript-review", wantSkills: []string{"orchestra-typescript"}, snapshot: []capability.MCPCandidate{{ID: "github", Health: "HEALTHY"}, {ID: "context7", Health: "PURCHASE_PENDING", Fallback: "docs already in the repo"}}, wantMCP: []string{"context7"}, fallback: "docs already in the repo", exclude: []string{"github"}},
	}

	if len(cases) != 45 {
		t.Fatalf("suite has %d cases, want 45", len(cases))
	}

	passed, failed := 0, 0
	for _, tc := range cases {
		tc := tc
		t.Run(tc.bucket+"/"+tc.name, func(t *testing.T) {
			if runSelCase(t, cls, g, cat, pw, shot, tc) {
				passed++
			} else {
				failed++
			}
		})
	}
	fmt.Printf("selection suite passed=%d failed=%d\n", passed, failed)
	if failed != 0 {
		t.Fatalf("selection suite failed=%d", failed)
	}
}

func runSelCase(t *testing.T, cls *Classifier, g *resources.DesignResourceGraph, cat *resources.ResourceCatalog, pw, shot bool, tc selCase) bool {
	t.Helper()
	ok := true
	b := cls.ClassifyBrief(tc.request, Options{})
	if b.CapabilityID != tc.wantCap {
		t.Errorf("capability = %q, want %q (%s)", b.CapabilityID, tc.wantCap, b.ArchetypeReason)
		ok = false
	}
	if tc.wantSkills != nil && !sameSet(b.Skills, tc.wantSkills) {
		t.Errorf("skills = %v, want %v", b.Skills, tc.wantSkills)
		ok = false
	}
	var cap *resources.CapabilityWorkflow
	if b.CapabilityID != "" && g != nil {
		cap = g.Capabilities[b.CapabilityID]
	}
	mcp := ResolveMCP(b.CapabilityID, cap, cat, tc.snapshot, tc.request)
	gotMCP := mcpIDs(mcp)
	if !sameSet(gotMCP, tc.wantMCP) {
		t.Errorf("mcp candidates = %v, want %v", gotMCP, tc.wantMCP)
		ok = false
	}
	for _, bad := range tc.exclude {
		if b.CapabilityID == bad || containsStr(b.Skills, bad) || containsStr(gotMCP, bad) {
			t.Errorf("excluded %q is present cap=%s skills=%v mcp=%v", bad, b.CapabilityID, b.Skills, gotMCP)
			ok = false
		}
	}

	m := mission.FromSeed(mission.Seed{ID: "sel-" + tc.name, Objective: tc.request})
	skills := append(append([]string{}, b.Skills...), tc.extraSkills...)
	trusted := append([]string{}, b.Skills...)
	for _, c := range mcp {
		trusted = append(trusted, c.ID)
	}
	d, err := capability.Select(capability.Input{
		Mission:                 m,
		Skills:                  skills,
		Trusted:                 trusted,
		MCP:                     mcp,
		PlaywrightAvailable:     &pw,
		ScreenshotEvalAvailable: &shot,
	})
	if tc.humanGate {
		if !capability.IsHumanStop(err) {
			t.Errorf("want human gate, got %v", err)
			ok = false
		}
	} else if err != nil {
		t.Errorf("select: %v", err)
		ok = false
	}
	if !tc.humanGate && !sameSet(d.Tools, tc.wantTools) {
		t.Errorf("tools = %v, want %v", d.Tools, tc.wantTools)
		ok = false
	}
	if tc.fallback != "" && !fallbackHas(d, tc.fallback) {
		t.Errorf("fallback %q missing in %+v", tc.fallback, d.Fallbacks)
		ok = false
	}
	if !tc.humanGate && tc.fallback == "" {
		for _, f := range d.Fallbacks {
			if f.Requested != "" && containsStr(gotMCP, f.Requested) {
				t.Errorf("unexpected fallback %+v", f)
				ok = false
			}
		}
	}
	for _, bad := range tc.exclude {
		if containsStr(d.Skills, bad) || containsStr(d.Tools, bad) {
			t.Errorf("selected excluded %q skills=%v tools=%v", bad, d.Skills, d.Tools)
			ok = false
		}
	}
	if tc.name == "small api" {
		dir := t.TempDir()
		if _, werr := capability.WriteDecision(dir, m.ID, d); werr != nil {
			t.Errorf("write decision: %v", werr)
			ok = false
		} else {
			raw, rerr := filepath.Abs(dir)
			_ = raw
			got, rerr := capability.ReadDecision(dir, m.ID)
			if rerr != nil {
				t.Errorf("read decision: %v", rerr)
				ok = false
			}
			blob := strings.Join(got.Skills, " ") + " " + strings.Join(got.Tools, " ")
			if strings.Contains(blob, "sk-") || strings.Contains(blob, "ghp_") {
				t.Errorf("capability record contains a token")
				ok = false
			}
			if !containsStr(got.Skills, "systematic-debugging") {
				t.Errorf("capability.json skills = %v", got.Skills)
				ok = false
			}
		}
	}
	return ok
}

func loadRealCatalog(t *testing.T) *resources.ResourceCatalog {
	t.Helper()
	p := filepath.Join("..", "..", "..", "registries", "resources.json")
	cat, err := resources.LoadResourceCatalog(p)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	return cat
}

func sameSet(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	g := append([]string{}, got...)
	w := append([]string{}, want...)
	sort.Strings(g)
	sort.Strings(w)
	for i := range g {
		if g[i] != w[i] {
			return false
		}
	}
	return true
}

func mcpIDs(cands []capability.MCPCandidate) []string {
	var out []string
	for _, c := range cands {
		out = append(out, c.ID)
	}
	return out
}

func containsStr(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func fallbackHas(d capability.Decision, text string) bool {
	for _, f := range d.Fallbacks {
		if f.Reason == text || strings.Contains(f.Reason, text) {
			return true
		}
	}
	return false
}

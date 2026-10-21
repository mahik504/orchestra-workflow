package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

func TestPhase3E1_DecisionMutationFailsAndFileStays(t *testing.T) {
	m := mission.FromSeed(mission.Seed{
		ID: "p3e1-auth", Objective: "Build a settings screen for profile CRUD", DeclaredType: "FEATURE",
	})
	pw := true
	d, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw})
	if err != nil {
		t.Fatal(err)
	}
	if !d.Has(capability.IDBrowserPlaywright) {
		t.Fatal("browser.playwright required")
	}
	dir := t.TempDir()
	path, err := capability.WriteDecision(dir, m.ID, d)
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	mutated := d
	var kept []capability.Item
	for _, it := range mutated.Selected {
		if it.ID != capability.IDBrowserPlaywright {
			kept = append(kept, it)
		}
	}
	mutated.Selected = kept
	if _, err := capability.WriteDecision(dir, m.ID, mutated); err != capability.ErrDecisionMutated {
		t.Fatalf("mutation must fail acceptance, got %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, after) {
		t.Fatal("capability.json changed after a rejected mutation")
	}
	got, err := capability.ReadDecision(dir, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Has(capability.IDBrowserPlaywright) {
		t.Fatal("authoritative decision lost browser.playwright")
	}
	ctx := &TaskContext{
		Task:           &TaskRequest{ID: m.ID, RawRequest: m.Objective, WorkspaceRoot: dir},
		Mission:        m,
		Capability:     &mutated,
		Classification: &ClassificationData{RequiresVisual: true},
	}
	if err := ensureCapabilityAuthoritative(ctx); err != capability.ErrDecisionMutated {
		t.Fatalf("downstream check: %v", err)
	}
	data, err := NewPlaywrightVerifier(t.TempDir()).Verify(context.Background(), ctx)
	if err != nil {
		t.Fatal(err)
	}
	if data.AllPassed {
		t.Fatal("mutated decision must not pass visual QA")
	}
}

func TestPhase3E1_SideListsDoNotLeakPastSelected(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "p3e1-side", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	d, err := capability.Select(capability.Input{
		Mission: m, PlaywrightAvailable: &pw,
		Skills:    []string{"superpowers-planning"},
		Resources: []string{"stdlib"},
		MCP:       []capability.MCPCandidate{{ID: "repo-read", Health: "HEALTHY"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	d.Skills = append(d.Skills, "FULL-SKILL-CATALOG")
	d.Resources = append(d.Resources, "FULL-RESOURCE-CATALOG")
	d.Tools = append(d.Tools, "FULL-MCP-CATALOG")
	ctx := &TaskContext{Mission: m, Capability: &d, Task: &TaskRequest{ID: m.ID}, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	raw, _ := json.Marshal(ctx.Capsule)
	blob := string(raw)
	if strings.Contains(blob, "FULL-SKILL-CATALOG") || strings.Contains(blob, "FULL-RESOURCE-CATALOG") || strings.Contains(blob, "FULL-MCP-CATALOG") {
		t.Fatalf("side slices leaked into capsule: %s", blob)
	}
	if !sameStrings(ctx.Capsule.Skills, d.IDs(capability.CatSkill)) ||
		!sameStrings(ctx.Capsule.Resources, d.IDs(capability.CatResource)) ||
		!sameStrings(ctx.Capsule.ToolsAllowed, d.IDs(capability.CatMCPTool)) {
		t.Fatalf("capsule lists != selected ids skills=%v resources=%v tools=%v", ctx.Capsule.Skills, ctx.Capsule.Resources, ctx.Capsule.ToolsAllowed)
	}
}

func TestPhase3E1_MCPSnapshotConsumedWithoutExecutor(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "p3e1-mcp", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	okDecision, err := capability.Select(capability.Input{
		Mission: m, PlaywrightAvailable: &pw,
		MCP: []capability.MCPCandidate{
			{ID: "repo-read", Health: "HEALTHY"},
			{ID: "optional-broken", Health: "BROKEN"},
			{ID: "optional-auth", Health: "AUTH_REQUIRED"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !containsStr(okDecision.Tools, "repo-read") {
		t.Fatalf("available tool not selected: %+v", okDecision.Tools)
	}
	if containsStr(okDecision.Tools, "optional-broken") || containsStr(okDecision.Tools, "optional-auth") {
		t.Fatalf("unavailable optional tools selected: %+v", okDecision.Tools)
	}
	ctx := &TaskContext{Mission: m, Capability: &okDecision, Task: &TaskRequest{ID: m.ID}, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	if !containsStr(ctx.Capsule.ToolsAllowed, "repo-read") {
		t.Fatalf("selected tool did not reach capsule: %v", ctx.Capsule.ToolsAllowed)
	}
	if containsStr(ctx.Capsule.ToolsAllowed, "optional-broken") {
		t.Fatal("broken tool reached capsule")
	}

	dir := t.TempDir()
	gated := &TaskContext{
		Mission: m,
		Task: &TaskRequest{
			ID: m.ID, WorkspaceRoot: dir,
			Parameters: map[string]any{
				"playwright_available": false,
				"capability_mcp": []capability.MCPCandidate{
					{ID: "repo-read", Health: "HEALTHY"},
					{ID: "creds", Health: "AUTH_REQUIRED", Required: true},
				},
			},
		},
		RunState: mission.NewRunState(m),
	}
	if err := selectCapabilities(gated); err != ErrHumanGateRequired {
		t.Fatalf("required auth tool: %v", err)
	}
	if gated.Capability == nil || !strings.Contains(gated.Capability.HumanDetail, "creds") {
		t.Fatalf("human detail: %+v", gated.Capability)
	}
	src, err := os.ReadFile(filepath.Join("..", "capability", "select.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(src), "mcp.json") {
		t.Fatal("capability selector references mcp.json")
	}
}

func TestPhase3E1_EvidenceStaysCompact(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "p3e1-ev", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	var skills []string
	for i := 0; i < 44; i++ {
		skills = append(skills, "allowlisted-skill-"+itoa3e(i))
	}
	var resources []string
	for i := 0; i < 168; i++ {
		resources = append(resources, "catalog-row-"+itoa3e(i))
	}
	var mcp []capability.MCPCandidate
	for i := 0; i < 20; i++ {
		mcp = append(mcp, capability.MCPCandidate{ID: "mcp-catalog-" + itoa3e(i), Health: "HEALTHY"})
	}
	d, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw, Skills: skills, Resources: resources, MCP: mcp})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path, err := capability.WriteDecision(dir, m.ID, d)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{
		"mission_id": true, "selected": true, "fallbacks": true,
		"rejected": true, "unavailable": true, "host_label": true,
	}
	for k := range keys {
		if !allowed[k] {
			t.Fatalf("unexpected evidence key %q", k)
		}
	}
	blob := string(raw)
	for _, banned := range []string{"FULL PROJECT", "conversation_history", "WHOLE PROJECT", "allowlisted-skill-43", "catalog-row-167", "mcp-catalog-19"} {
		if strings.Contains(blob, banned) {
			t.Fatalf("evidence contains %q (%d bytes)", banned, len(raw))
		}
	}
	if strings.Count(blob, "allowlisted-skill-") >= 44 || strings.Count(blob, "catalog-row-") >= 168 || strings.Count(blob, "mcp-catalog-") >= 20 {
		t.Fatalf("catalog dump in evidence (%d bytes)", len(raw))
	}
	dump := bytes.Repeat([]byte("RESOURCE DUMP LINE\n"), 8000)
	t.Logf("hostile capability_bytes=%d dump_bytes=%d", len(raw), len(dump))
	if len(raw) > 20000 || float64(len(raw))/float64(len(dump)) >= 0.2 {
		t.Fatalf("evidence not compact: %d", len(raw))
	}
	ctx := &TaskContext{Mission: m, Capability: &d, Task: &TaskRequest{ID: m.ID}, RunState: mission.NewRunState(m)}
	compileBuilderCapsule(ctx)
	craw, _ := json.Marshal(ctx.Capsule)
	t.Logf("hostile capsule_bytes=%d", len(craw))
	if len(craw) > 20000 {
		t.Fatalf("capsule not compact: %d", len(craw))
	}
}

func TestPhase3E1_FallbackHopsStayHonest(t *testing.T) {
	visual.ResetScreenshotEvaluatorForTest()
	defer visual.ResetScreenshotEvaluatorForTest()
	m := mission.FromSeed(mission.Seed{
		ID: "p3e1-fb", Objective: "Build a flagship portfolio landing with unusual motion",
		DeclaredType: "DESIGN", QualityBar: mission.BarPremium,
	})
	pw, shot := true, false
	d, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw, ScreenshotEvalAvailable: &shot})
	if err != nil {
		t.Fatal(err)
	}
	if !hop(d, capability.IDVisualMultimodal, capability.IDVisualStructural) {
		t.Fatalf("missing multimodal→structural: %+v", d.Fallbacks)
	}
	if !hop(d, capability.IDVisualScreenshot, capability.IDVisualStructural) {
		t.Fatalf("missing screenshot→structural: %+v", d.Fallbacks)
	}
	if d.Has(capability.IDVisualScreenshot) || d.Has(capability.IDVisualMultimodal) {
		t.Fatalf("unavailable vision selected: %+v", d.Selected)
	}
	for _, it := range d.Selected {
		if it.ID == capability.IDVisualStructural && (it.Category != capability.CatVisualStructural || strings.Contains(strings.ToLower(it.Reason), "multimodal")) {
			t.Fatalf("structural labeled as vision: %+v", it)
		}
	}
	if visual.ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("select must not execute a screenshot adapter")
	}
	dir := t.TempDir()
	png := filepath.Join(dir, "desktop.png")
	if err := os.WriteFile(png, visual.TinyPNG, 0644); err != nil {
		t.Fatal(err)
	}
	rep := visual.Evaluate(visual.LayerBInput{
		Class: visual.ClassFlagship, Review: string(mission.VisualReviewDeep),
		Captures: []visual.Capture{{Viewport: visual.Viewport{Name: "desktop", Width: 1440, Height: 900}, ScreenshotPath: png, Passed: true}},
	})
	if rep.EvaluatorType != visual.EvaluatorStructural {
		t.Fatalf("evaluator_type %s", rep.EvaluatorType)
	}

	shot = true
	withShot, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw, ScreenshotEvalAvailable: &shot})
	if err != nil {
		t.Fatal(err)
	}
	if !hop(withShot, capability.IDVisualMultimodal, capability.IDVisualScreenshot) {
		t.Fatalf("missing multimodal→screenshot: %+v", withShot.Fallbacks)
	}
	if visual.ScreenshotEvaluatorInvocations() != 0 {
		t.Fatal("screenshot capability must not claim an execution")
	}
	applied := visual.ApplyScreenshotEval(context.Background(), visual.LayerBInput{
		Class: visual.ClassFlagship, Review: string(mission.VisualReviewDeep),
		Captures: []visual.Capture{{Viewport: visual.Viewport{Name: "desktop", Width: 1440, Height: 900}, ScreenshotPath: png, Passed: true}},
	}, m.Objective, "", nil)
	if applied.EvaluatorType == "multimodal" || applied.Invoked {
		t.Fatalf("nil adapter claimed a vision pass: %+v", applied)
	}

	pw = false
	_, err = capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw})
	if !capability.IsHumanStop(err) {
		t.Fatalf("browser fallback must stop, got %v", err)
	}
}

func TestPhase3E1_OpaqueLabelsDoNotChangeClass(t *testing.T) {
	m := mission.FromSeed(mission.Seed{ID: "p3e1-model", Objective: "Fix parser crash", DeclaredType: "BUGFIX"})
	pw := false
	labels := [][]string{
		{"grok-4.7"},
		{"claude-opus"},
		{"gemini-2.5"},
		{"gpt-x"},
		nil,
	}
	var first string
	for _, label := range labels {
		d, err := capability.Select(capability.Input{Mission: m, PlaywrightAvailable: &pw, DiscoveredModels: label})
		if err != nil {
			t.Fatal(err)
		}
		got := capIDs(&d)
		if first == "" {
			first = got
		} else if got != first {
			t.Fatalf("label %v changed selection: %s vs %s", label, got, first)
		}
		raw, _ := json.Marshal(d)
		low := strings.ToLower(string(raw))
		for _, tok := range []string{"grok", "claude", "gemini", "gpt"} {
			if strings.Contains(low, tok) {
				t.Fatalf("decision leaked %s from %v", tok, label)
			}
		}
	}
}

func hop(d capability.Decision, from, to string) bool {
	for _, f := range d.Fallbacks {
		if f.Requested == from && f.Selected == to {
			return true
		}
	}
	return false
}

func containsStr(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package mission

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCompileCapsuleScopesContext(t *testing.T) {
	m := validMission()
	whole := []byte(strings.Repeat("PROJECT DUMP AGENTS.md SKILLS MCP\n", 4000))
	cap, report, err := Compile(CompileRequest{
		Mission:             m,
		Role:                RoleBuilder,
		Files:               []string{"internal/http/health.go", "README.md"},
		UnrelatedFiles:      []string{"README.md"},
		Skills:              []string{"ow-conductor", "gsap-core"},
		AllSkills:           []string{"ow-conductor", "gsap-core", "hallmark", "supabase"},
		Resources:           []string{"playwright"},
		ExpectedOutput:      "patched handler plus test evidence",
		ToolsAllowed:        []string{"go-test"},
		ConversationHistory: []string{"yesterday we talked about fonts", "unrelated chat"},
		MCPCatalog:          []string{"playwright", "context7", "stitch"},
		WholeProjectContext: whole,
		Decisions:           []string{"keep auth unchanged"},
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if cap.MissionID != m.ID {
		t.Fatalf("capsule mission id")
	}
	if cap.Role != RoleBuilder {
		t.Fatalf("role %s", cap.Role)
	}
	for _, f := range cap.Files {
		if f == "README.md" {
			t.Fatalf("unrelated file leaked into capsule")
		}
	}
	if len(cap.Files) != 1 || cap.Files[0] != "internal/http/health.go" {
		t.Fatalf("files = %v", cap.Files)
	}
	for _, s := range cap.Skills {
		if s == "hallmark" || s == "supabase" {
			t.Fatalf("unrelated skill leaked: %v", cap.Skills)
		}
	}
	raw, _ := json.Marshal(cap)
	blob := string(raw)
	if strings.Contains(blob, "yesterday we talked") {
		t.Fatalf("conversation history copied into capsule")
	}
	if strings.Contains(blob, "context7") {
		t.Fatalf("mcp catalog copied into capsule")
	}
	if strings.Contains(blob, "PROJECT DUMP") {
		t.Fatalf("whole-project dump copied into capsule")
	}
	if report.CapsuleBytes <= 0 {
		t.Fatalf("capsule size not measured")
	}
	if report.WholeProjectBytes != len(whole) {
		t.Fatalf("whole project bytes")
	}
	if report.Ratio >= 1 {
		t.Fatalf("capsule should be smaller than whole-project context, ratio=%f", report.Ratio)
	}
	t.Logf("capsule_bytes=%d whole_project_bytes=%d ratio=%.4f excluded=%v", report.CapsuleBytes, report.WholeProjectBytes, report.Ratio, report.Excluded)
	foundHist := false
	for _, e := range report.Excluded {
		if e == "conversation_history" {
			foundHist = true
		}
	}
	if !foundHist {
		t.Fatalf("excluded list missing conversation_history: %v", report.Excluded)
	}
}

func TestCompileCapsuleProviderNeutral(t *testing.T) {
	m := validMission()
	cap, _, err := Compile(CompileRequest{
		Mission:        m,
		Role:           RoleBuilder,
		ExpectedOutput: "diff",
		Skills:         []string{"ow-conductor"},
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(cap)
	if leak := ProviderLeak(raw); leak != "" {
		t.Fatalf("capsule leaked %q: %s", leak, raw)
	}
	if cap.Role == "claude" || cap.Role == "cursor" {
		t.Fatalf("role must stay planner/builder/judge")
	}
}

func TestCompileCapsuleRequiresMission(t *testing.T) {
	_, _, err := Compile(CompileRequest{Role: RoleBuilder})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCompileIgnoresExtraIrrelevantContext(t *testing.T) {
	m := validMission()
	cap, report, err := Compile(CompileRequest{
		Mission:             m,
		Role:                RoleBuilder,
		Files:               []string{"internal/http/health.go"},
		Skills:              []string{"ow-conductor"},
		AllSkills:           []string{"ow-conductor", "hallmark", "supabase", "gsap-core"},
		UnrelatedFiles:      []string{"STACK.md", "CHANGELOG.md"},
		ConversationHistory: []string{"long forgotten thread about fonts"},
		MCPCatalog:          []string{"playwright", "context7", "stitch", "mobbin"},
		WholeProjectContext: []byte("WHOLE REPO " + strings.Repeat("x", 8000)),
		ExpectedOutput:      "patch",
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(cap)
	blob := string(raw)
	for _, banned := range []string{"hallmark", "supabase", "STACK.md", "CHANGELOG.md", "forgotten thread", "context7", "mobbin", "WHOLE REPO"} {
		if strings.Contains(blob, banned) {
			t.Fatalf("irrelevant %q leaked into capsule", banned)
		}
	}
	if report.Ratio >= 1 {
		t.Fatalf("ratio %f", report.Ratio)
	}
}

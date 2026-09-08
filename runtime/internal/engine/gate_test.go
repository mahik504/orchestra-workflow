package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/classifier"
	"github.com/mahik504/orchestra-workflow/runtime/internal/verify"
)

// The gate has to stop the pipeline before implementation, not merely warn.
// Before this existed, a PREMIUM visual brief could reach the implement stage
// with no approved direction as long as it was a dry run.
func TestGate_BlocksImplementationWithoutApproval(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())

	req := &TaskRequest{
		ID:            "task-gate-locked",
		RawRequest:    "Build a landing page for a client's studio, it needs to feel expensive",
		WorkspaceRoot: workdir,
		Type:          "DESIGN",
		DryRun:        true,
	}

	_, err := pipeline.Execute(context.Background(), req)
	if err == nil {
		t.Fatal("pipeline completed with no approved direction")
	}

	var gateErr *verify.ErrGateNotCleared
	if !errors.Is(err, ErrHumanGateRequired) && !errors.As(err, &gateErr) {
		t.Fatalf("expected a gate error, got: %v", err)
	}

	// Nothing a browser renders should exist in the workspace.
	if found := findFrontendFiles(t, workdir); len(found) > 0 {
		t.Errorf("gate was not cleared but frontend files were written: %v", found)
	}
}

// A STANDARD brief owes no lab and must not be blocked by one.
func TestGate_StandardWorkIsNotGated(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())

	req := &TaskRequest{
		ID:            "task-gate-standard",
		RawRequest:    "Internal tool for our team to see which servers are down right now, live",
		WorkspaceRoot: workdir,
		Type:          "FEATURE",
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("STANDARD work was blocked: %v", err)
	}
	if res.Classification == nil {
		t.Fatal("no classification produced")
	}
	if res.Classification.RequiresHumanGate {
		t.Errorf("STANDARD work armed the gate: %s", res.Classification.GateReason)
	}
	if res.Classification.QualityBar != "STANDARD" {
		t.Errorf("quality bar = %q, want STANDARD", res.Classification.QualityBar)
	}
}

// The brief the human approves has to be on the record, including the routes we
// turned down.
func TestClassify_RecordsBriefAndDeclinedRoutes(t *testing.T) {
	cat, graph, workdir := setupEngineFixtures(t)
	defer os.RemoveAll(workdir)

	pipeline := NewDesignPipeline(cat, graph, nil, NewMockVisualVerifier())
	req := &TaskRequest{
		ID:             "task-brief",
		RawRequest:     "Build a scheduling dashboard for a school with attendance charts",
		WorkspaceRoot:  workdir,
		Type:           "FEATURE",
		SkipVisualGate: true,
	}

	res, err := pipeline.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}

	c := res.Classification
	if c.Brief == nil {
		t.Fatal("no brief recorded on the classification")
	}
	if c.Brief.CapabilityID != "saas-dashboard" {
		t.Errorf("capability = %q, want saas-dashboard (%s)", c.Brief.CapabilityID, c.Brief.ArchetypeReason)
	}
	if len(c.DeclinedRoutes) == 0 {
		t.Error("no declined routes recorded; every row should be answered for")
	}
	for _, d := range c.DeclinedRoutes {
		if d.DeclineReason == "" {
			t.Errorf("%s declined with no reason", d.CapabilityID)
		}
	}
	if c.ResearchDepth == "" || c.VerifyDepth == "" {
		t.Errorf("brief did not set research/verify depth: %q / %q", c.ResearchDepth, c.VerifyDepth)
	}
}

func TestGate_ContractApprovalDoesNotUnlockProductPage(t *testing.T) {
	workdir := t.TempDir()
	lab := verify.NewDesignLab(&classifier.Brief{
		TaskID:            "engine-contract",
		DesignLabRequired: true,
		DesignLabReason:   "PREMIUM",
	}, workdir)
	if err := lab.SkipSurvey("pasted DESIGN.md fixture"); err != nil {
		t.Fatal(err)
	}
	dir := verify.Direction{
		ID: "a", Concept: "Paper field",
		Typography: "Newsreader + IBM Plex Sans", TypographySrc: "editorial newspapers",
		ColorWorld: "Warm paper", ColorSrc: "Japanese stationery",
		LayoutLanguage: "asymmetric columns", ComponentKit: "custom",
		MotionEngine: "CSS", MotionWhy: "one engine, reduced-motion first",
		ThreeD: "no", Shader: "no", LogoMethod: "wordmark", IconSystem: "line",
		Stack: []string{"astro", "tailwind"},
	}
	if err := lab.OfferContract(dir); err != nil {
		t.Fatal(err)
	}
	if err := lab.ApproveContract("a", "operator"); err != nil {
		t.Fatal(err)
	}
	if lab.Cleared() {
		t.Fatal("contract approval cleared the product write lock")
	}
	if err := lab.GuardWrite(filepath.Join(workdir, "app", "page.tsx")); err == nil {
		t.Fatal("CONTRACT_APPROVED allowed product page.tsx")
	}
	still := filepath.Join(lab.GoldenStillsDir(), "opening.tsx")
	if err := os.MkdirAll(filepath.Dir(still), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := lab.GuardWrite(still); err != nil {
		t.Fatalf("golden-stills renderer blocked: %v", err)
	}
}

func TestGate_StillsApprovalUnlocksProductPage(t *testing.T) {
	workdir := t.TempDir()
	lab := verify.NewDesignLab(&classifier.Brief{
		TaskID:            "engine-stills",
		DesignLabRequired: true,
		DesignLabReason:   "PREMIUM",
	}, workdir)
	if err := lab.SkipSurvey("pasted DESIGN.md fixture"); err != nil {
		t.Fatal(err)
	}
	dir := verify.Direction{
		ID: "a", Concept: "Paper field",
		Typography: "Newsreader + IBM Plex Sans", TypographySrc: "editorial newspapers",
		ColorWorld: "Warm paper", ColorSrc: "Japanese stationery",
		LayoutLanguage: "asymmetric columns", ComponentKit: "custom",
		MotionEngine: "CSS", MotionWhy: "one engine, reduced-motion first",
		ThreeD: "no", Shader: "no", LogoMethod: "wordmark", IconSystem: "line",
		Stack: []string{"astro", "tailwind"},
	}
	if err := lab.OfferContract(dir); err != nil {
		t.Fatal(err)
	}
	if err := lab.ApproveContract("a", "operator"); err != nil {
		t.Fatal(err)
	}
	desk := filepath.Join(lab.GoldenStillsDir(), "desktop.png")
	mob := filepath.Join(lab.GoldenStillsDir(), "mobile.png")
	if err := os.MkdirAll(lab.GoldenStillsDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(desk, []byte("desk"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mob, []byte("mob"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := lab.ApproveStills("operator", desk, mob, "desktop and mobile opening scenes pass"); err != nil {
		t.Fatalf("ApproveStills: %v", err)
	}
	if !lab.Cleared() {
		t.Fatal("stills approval did not clear the gate")
	}
	if err := lab.GuardWrite(filepath.Join(workdir, "app", "page.tsx")); err != nil {
		t.Fatalf("product page.tsx still blocked after stills: %v", err)
	}
}

func findFrontendFiles(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var walk func(dir string, items []os.DirEntry)
	walk = func(dir string, items []os.DirEntry) {
		for _, e := range items {
			p := dir + string(os.PathSeparator) + e.Name()
			if e.IsDir() {
				if sub, err := os.ReadDir(p); err == nil {
					walk(p, sub)
				}
				continue
			}
			if verify.IsFrontendPath(p) {
				out = append(out, p)
			}
		}
	}
	walk(root, entries)
	return out
}

package onboard

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/orchestra-v3/internal/memory"
	"github.com/user/orchestra-v3/internal/resources"
)

func TestExtractURLFromIntent(t *testing.T) {
	intent := "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability. https://github.com/example/example-resource"
	got := ExtractURL(intent)
	if got != "https://github.com/example/example-resource" {
		t.Fatalf("ExtractURL = %q", got)
	}
}

func TestInferKinds(t *testing.T) {
	skill := Infer(Inspection{
		URL: "https://github.com/acme/cool-skill", NormalizedURL: "https://github.com/acme/cool-skill",
		Host: "github.com", Owner: "acme", Repo: "cool-skill", StatusCode: 200, Reachable: true, HasSkillMD: true,
	}, "", "url_submit")
	if skill.Kind != KindSkill || skill.Resource.Representation != "skill" {
		t.Fatalf("skill: kind=%s rep=%s", skill.Kind, skill.Resource.Representation)
	}

	dep := Infer(Inspection{
		URL: "https://github.com/acme/lib", NormalizedURL: "https://github.com/acme/lib",
		Host: "github.com", Owner: "acme", Repo: "lib", StatusCode: 200, Reachable: true, HasPackageJSON: true,
	}, "", "url_submit")
	if dep.Kind != KindDependency || dep.InstallScope != ScopeProject {
		t.Fatalf("dep: kind=%s scope=%s", dep.Kind, dep.InstallScope)
	}

	mcp := Infer(Inspection{
		URL: "https://github.com/acme/tools-mcp", NormalizedURL: "https://github.com/acme/tools-mcp",
		Host: "github.com", Owner: "acme", Repo: "tools-mcp", StatusCode: 200, Reachable: true, HasMCPManifest: true,
	}, "", "url_submit")
	if mcp.Kind != KindMCP {
		t.Fatalf("mcp kind=%s", mcp.Kind)
	}

	ptr := Infer(Inspection{
		URL: "https://github.com/example/example-resource", NormalizedURL: "https://github.com/example/example-resource",
		Host: "github.com", Owner: "example", Repo: "example-resource", StatusCode: 404, Reachable: false,
	}, "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability.", "user_intent")
	if ptr.Kind != KindReference {
		t.Fatalf("404 github should be reference, got %s", ptr.Kind)
	}
	if ptr.InstallScope != ScopeOnDemand {
		t.Fatalf("404 should be ON_DEMAND, got %s", ptr.InstallScope)
	}
	if ptr.Resource.AcquisitionMethod != "git" {
		t.Fatalf("github acquisition want git, got %s", ptr.Resource.AcquisitionMethod)
	}
	for _, tag := range ptr.Resource.RoutingTags {
		if tag == "example" || tag == "resource" {
			t.Fatalf("weak routing tag leaked: %s", tag)
		}
	}
	if len(ptr.Resource.TriggerConditions) == 0 || len(ptr.Resource.AvoidConditions) == 0 {
		t.Fatalf("triggers/skips missing")
	}
}

func TestInferCopiesNpmPackageName(t *testing.T) {
	prop := Infer(Inspection{
		URL: "https://github.com/pmndrs/drei", NormalizedURL: "https://github.com/pmndrs/drei",
		Host: "github.com", Owner: "pmndrs", Repo: "drei", StatusCode: 200, Reachable: true,
		HasPackageJSON: true, PackageName: "@react-three/drei",
	}, "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability.", "user_intent")
	if prop.Kind != KindDependency || prop.InstallScope != ScopeProject {
		t.Fatalf("drei: kind=%s scope=%s", prop.Kind, prop.InstallScope)
	}
	if prop.Resource.NpmPackage != "@react-three/drei" {
		t.Fatalf("npm package: %s", prop.Resource.NpmPackage)
	}
	if prop.Resource.AcquisitionMethod != "npm" {
		t.Fatalf("acq=%s", prop.Resource.AcquisitionMethod)
	}
}

func TestParseNpmPackageName(t *testing.T) {
	if got := parseNpmPackageName(`{ "name": "@react-three/drei", "version": "1.0.0" }`); got != "@react-three/drei" {
		t.Fatalf("got %q", got)
	}
}

func TestMatchVsMismatchAndMemoryPolicy(t *testing.T) {
	insp := Inspection{
		URL: "https://github.com/example/example-resource", NormalizedURL: "https://github.com/example/example-resource",
		Host: "github.com", Owner: "example", Repo: "example-resource", StatusCode: 404, Reachable: false,
	}
	prop := Infer(insp, "whenever the task requires its capability", "user_intent")
	doc := &OverlayDocument{Resources: []OverlayEntry{{
		Kind: prop.Kind, InstallScope: prop.InstallScope, RoutingPolicy: "active",
		Resource: prop.Resource, Inspection: insp,
	}}}

	match := EvaluateTask("Use example-resource for the capability this repository provides.", doc)
	if len(match) != 1 || match[0].Action != ActionActivated {
		t.Fatalf("matching task should activate: %+v", match)
	}
	mismatch := EvaluateTask("The login form throws a 500 when the email has a plus sign.", doc)
	if len(mismatch) != 1 || mismatch[0].Action != ActionSuppressed {
		t.Fatalf("non-matching task should suppress: %+v", mismatch)
	}

	tmp := t.TempDir()
	store, err := memory.NewResourceMemoryStore(filepath.Join(tmp, "resource-memory.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(&memory.ResourceEvaluation{
		ResourceID: prop.Resource.ID, Outcome: memory.OutcomeSuccess, QualityScore: 1, TaskID: "ok",
	}); err != nil {
		t.Fatal(err)
	}
	ApplyMemoryPolicy(doc, store)
	if EvaluateTask("Use example-resource here", doc)[0].Action != ActionActivated {
		t.Fatal("success should keep activation")
	}
	if err := store.Record(&memory.ResourceEvaluation{
		ResourceID: prop.Resource.ID, Outcome: memory.OutcomeFailure, QualityScore: 0.1, TaskID: "fail",
		ErrorDetails: "git ls-remote failed",
	}); err != nil {
		t.Fatal(err)
	}
	ApplyMemoryPolicy(doc, store)
	d := EvaluateTask("Use example-resource here", doc)[0]
	if d.Action != ActionSuppressed {
		t.Fatalf("failure should suppress even on a matching task, got %s (%s)", d.Action, d.Reason)
	}
	if doc.Resources[0].RoutingPolicy != "suppressed" {
		t.Fatalf("policy want suppressed, got %s", doc.Resources[0].RoutingPolicy)
	}
}

func TestLifecycleDoesNotEditPublicCatalog(t *testing.T) {
	tmp := t.TempDir()
	public := filepath.Join(tmp, "resources.json")
	original := []byte(`[{"id":"gsap","name":"GSAP","canonical_url":"https://gsap.com","source_type":"web_reference","category":["MOTION"],"representation":"dependency","routing_tags":["motion"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`)
	if err := os.WriteFile(public, original, 0644); err != nil {
		t.Fatal(err)
	}

	overlayPath := filepath.Join(tmp, "added-resources.json")
	memPath := filepath.Join(tmp, "resource-memory.json")
	workflow := findWorkflowRoot(t)

	rep, err := RunLifecycle(Options{
		URL:           "https://github.com/example/example-resource",
		Intent:        "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability. https://github.com/example/example-resource",
		Origin:        "user_intent",
		OverlayPath:   overlayPath,
		MemoryPath:    memPath,
		PublicCatalog: public,
		WorkflowRoot:  workflow,
		InspectFn: func(u string) Inspection {
			return Inspection{
				URL: u, NormalizedURL: "https://github.com/example/example-resource",
				Host: "github.com", Owner: "example", Repo: "example-resource",
				StatusCode: 404, Reachable: false, Error: "injected 404",
			}
		},
		GitProbeFn: func(source string) GitProbe {
			return GitProbe{Command: "git ls-remote --heads " + source, OK: false, ExitCode: 128, Output: "repository not found"}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !rep.CatalogUnchanged {
		t.Fatal("public catalog was modified")
	}
	after, _ := os.ReadFile(public)
	if string(after) != string(original) {
		t.Fatal("public catalog bytes changed")
	}
	if len(rep.Steps) != 15 {
		t.Fatalf("want 15 steps, got %d", len(rep.Steps))
	}
	for _, s := range rep.Steps {
		if s.Status == "FAIL" {
			t.Errorf("step %d %s failed: %s", s.N, s.Name, s.Detail)
		}
	}
	if rep.MatchAfterSuccess != ActionActivated {
		t.Fatalf("after success want activated, got %s", rep.MatchAfterSuccess)
	}
	if rep.MatchAfterFailure != ActionSuppressed {
		t.Fatalf("after failure want suppressed, got %s", rep.MatchAfterFailure)
	}

	doc, err := LoadOverlay(overlayPath)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Find("example-resource") == nil {
		t.Fatal("overlay missing example-resource")
	}

	store, err := memory.NewResourceMemoryStore(memPath)
	if err != nil {
		t.Fatal(err)
	}
	evals := store.ListEvaluations(0, "example-resource")
	if len(evals) != 1 {
		t.Fatalf("want 1 failure evaluation, got %d", len(evals))
	}
	if evals[0].Outcome != memory.OutcomeFailure {
		t.Fatalf("Pass A must only record failure, got %s", evals[0].Outcome)
	}

	cat, err := resources.LoadResourceCatalog(public)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cat.FindByID("example-resource"); ok {
		t.Fatal("example-resource leaked into the public catalog file")
	}
	if err := MergeIntoCatalog(cat, doc); err != nil {
		t.Fatal(err)
	}
	if _, ok := cat.FindByID("example-resource"); !ok {
		t.Fatal("overlay merge should make the resource visible to the live catalog")
	}

	if _, err := os.Stat(filepath.Join(tmp, "lifecycle-pass-a", "03-inspection.json")); err != nil {
		t.Fatalf("pass A inspection artifact missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "lifecycle-pass-a", "15-failure-evaluation.json")); err != nil {
		t.Fatalf("pass A failure artifact missing: %v", err)
	}

	raw, _ := json.Marshal(rep)
	if strings.Contains(strings.ToLower(string(raw)), "reinforcement learning") && !strings.Contains(rep.NotRL, "No training loop") {
		t.Fatal("must not describe this as RL")
	}
}

func TestAddFromIntentWritesOverlayOnly(t *testing.T) {
	tmp := t.TempDir()
	overlayPath := filepath.Join(tmp, "added-resources.json")
	entry, _, err := AddFromIntent(Options{
		Intent:      "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability. https://github.com/example/example-resource",
		OverlayPath: overlayPath,
		InspectFn: func(u string) Inspection {
			return Inspection{URL: u, NormalizedURL: u, Host: "github.com", Owner: "example", Repo: "example-resource", StatusCode: 404}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if entry.Origin != "user_intent" {
		t.Fatalf("origin=%s", entry.Origin)
	}
	if _, err := os.Stat(overlayPath); err != nil {
		t.Fatal(err)
	}
}

func TestArchitectureFilesExist(t *testing.T) {
	root := findWorkflowRoot(t)
	facts := VerifyArchitecture(root)
	if len(facts) == 0 {
		t.Fatal("no facts")
	}
	for _, f := range facts {
		if !f.OK {
			t.Errorf("%s: %s (%s)", f.Claim, f.Evidence, f.Path)
		}
	}
}

func TestInspectGitHubPackageName(t *testing.T) {
	prev := HTTPClient
	HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		u := req.URL.String()
		code := 404
		body := "not found"
		switch {
		case strings.Contains(u, "github.com/pmndrs/drei") && !strings.Contains(u, "raw.githubusercontent.com"):
			code, body = 200, "<html><title>pmndrs/drei</title></html>"
		case strings.Contains(u, "raw.githubusercontent.com/pmndrs/drei/HEAD/package.json"):
			code, body = 200, `{"name":"@react-three/drei","version":"10.7.8"}`
		}
		return &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})}
	defer func() { HTTPClient = prev }()

	insp := Inspect("https://github.com/pmndrs/drei")
	if !insp.Reachable || !insp.HasPackageJSON {
		t.Fatalf("inspect: reachable=%v pkg=%v err=%s", insp.Reachable, insp.HasPackageJSON, insp.Error)
	}
	if insp.PackageName != "@react-three/drei" {
		t.Fatalf("package name: %q", insp.PackageName)
	}
	prop := Infer(insp, "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability.", "user_intent")
	if prop.Resource.NpmPackage != "@react-three/drei" || prop.Resource.AcquisitionMethod != "npm" {
		t.Fatalf("infer npm=%s method=%s", prop.Resource.NpmPackage, prop.Resource.AcquisitionMethod)
	}

	tmp := t.TempDir()
	public := filepath.Join(tmp, "resources.json")
	original := []byte(`[{"id":"r3f","name":"R3F","canonical_url":"https://github.com/pmndrs/react-three-fiber","source_type":"github_repository","category":["3D"],"representation":"dependency","routing_tags":["r3f"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`)
	if err := os.WriteFile(public, original, 0644); err != nil {
		t.Fatal(err)
	}
	overlayPath := filepath.Join(tmp, "added-resources.json")
	_, _, err := AddFromIntent(Options{
		Intent:      "Add this GitHub repository to Orchestra and make it available whenever the task requires its capability. https://github.com/pmndrs/drei",
		OverlayPath: overlayPath,
		InspectFn:   func(string) Inspection { return insp },
	})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(public)
	if string(after) != string(original) {
		t.Fatal("public catalog must stay unchanged")
	}
	doc, err := LoadOverlay(overlayPath)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Find("drei") == nil {
		t.Fatal("overlay missing drei")
	}
	match := EvaluateTask("Install the drei helper library into this Node fixture so the package resolves", doc)
	if len(match) != 1 || match[0].Action != ActionActivated {
		t.Fatalf("matching task should activate drei: %+v", match)
	}
	mismatch := EvaluateTask("The login form throws a 500 when the email has a plus sign.", doc)
	if len(mismatch) != 1 || mismatch[0].Action != ActionSuppressed {
		t.Fatalf("mismatch should skip drei: %+v", mismatch)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestResolveOverlayPath_HomeBeatsWorkspaceSeed(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "brain")
	clone := filepath.Join(tmp, "method")
	if err := os.MkdirAll(filepath.Join(clone, "memory"), 0755); err != nil {
		t.Fatal(err)
	}
	seed := filepath.Join(clone, "memory", "added-resources.json")
	if err := os.WriteFile(seed, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ORCHESTRA_OVERLAY_PATH", "")
	t.Setenv("ORCHESTRA_HOME", home)
	got := ResolveOverlayPath(clone)
	want := filepath.Join(home, "memory", "added-resources.json")
	if got != filepath.Clean(want) {
		t.Fatalf("HOME set: got %s want %s", got, want)
	}

	t.Setenv("ORCHESTRA_HOME", "")
	got = ResolveOverlayPath(clone)
	want = filepath.Join(clone, ".orchestra", "memory", "added-resources.json")
	if got != filepath.Clean(want) {
		t.Fatalf("HOME unset: got %s want %s", got, want)
	}
}

func findWorkflowRoot(t *testing.T) string {
	t.Helper()
	dir, _ := os.Getwd()
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "registries")); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("workflow root not found")
	return ""
}

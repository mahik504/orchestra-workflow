package compile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type fixture struct {
	Task       string   `json:"task"`
	Host       string   `json:"host"`
	Route      string   `json:"route"`
	Skills     []string `json:"skills"`
	ModelClass string   `json:"model_class"`
	Disclosure string   `json:"disclosure"`
}

func TestCompileFixtures(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join("..", "..", "testdata", "compile", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) < 12 {
		t.Fatalf("fixtures %d, want at least 12", len(matches))
	}
	for _, name := range matches {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		var want fixture
		if err := json.Unmarshal(raw, &want); err != nil {
			t.Fatal(err)
		}
		got, err := Compile(want.Task, want.Host, routes)
		if err != nil {
			t.Fatal(err)
		}
		if got.Route != want.Route {
			t.Fatalf("%s route %q, want %q", name, got.Route, want.Route)
		}
		if !reflect.DeepEqual(got.Skills, want.Skills) {
			t.Fatalf("%s skills %v, want %v", name, got.Skills, want.Skills)
		}
		if got.ModelClass != want.ModelClass || got.Disclosure != want.Disclosure {
			t.Fatalf("%s class %s/%s", name, got.ModelClass, got.Disclosure)
		}
		if got.Host != want.Host {
			t.Fatalf("%s host %s", name, got.Host)
		}
		if got.Route == "One 3D scene" && !reflect.DeepEqual(got.Libraries, []string{"three", "r3f"}) {
			t.Fatalf("%s libraries %v", name, got.Libraries)
		}
	}
}

func TestCompileReviewCases(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		task  string
		route string
	}{
		{"fix the production crash", "Bug"},
		{"bug in the api", "Bug"},
		{"the api is crashing", "Bug"},
		{"add citations", "Paper"},
		{"one control for the date", "One control"},
		{"read the pdf", "Docs on disk"},
		{"reverse this apk", "Binary or APK"},
		{"Fix the bug in the landing page contact form", "Bug"},
		{"Write the README for the landing page repo", "README and post"},
		{"Write a research paper on landing page conversion", "Paper"},
	}
	for _, item := range cases {
		got, err := Compile(item.task, "cursor", routes)
		if err != nil {
			t.Fatal(err)
		}
		if got.Route != item.route {
			t.Fatalf("%q route %q, want %q", item.task, got.Route, item.route)
		}
	}
	control, err := Compile("one control for the date", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if len(control.Skills) != 0 {
		t.Fatalf("control skills %v", control.Skills)
	}
	pdf, err := Compile("read the pdf", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if len(pdf.Skills) != 0 {
		t.Fatalf("pdf skills %v", pdf.Skills)
	}
}

func TestBareAppDoesNotStealShip(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Compile("rename the app icon", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if got.Route == "Ship a deployed app" {
		t.Fatalf("bare app stole the ship route: %+v", got)
	}
}

func TestFrontendAndFullStackNameWebappTesting(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	award, err := Compile("Build an award-style agency site", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(award.Verification, "\n"), "webapp-testing on the running app") {
		t.Fatalf("award verification %v", award.Verification)
	}
	full, err := Compile("Build the full-stack app", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if full.Route != "Full-stack app" {
		t.Fatalf("route %s", full.Route)
	}
	if !strings.Contains(strings.Join(full.Verification, "\n"), "webapp-testing on the running app") {
		t.Fatalf("full-stack verification %v", full.Verification)
	}
	saas, err := Compile("Design the SaaS API for billing", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if saas.Route != "SaaS or API" {
		t.Fatalf("saas route %s", saas.Route)
	}
	if !strings.Contains(strings.Join(saas.Constraints, "\n"), "saas-checklist.md") {
		t.Fatalf("saas constraints %v", saas.Constraints)
	}
}

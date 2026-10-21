package visual

import "testing"

func TestClassifyOrdinaryAPI(t *testing.T) {
	got := Classify(ClassifyInput{Type: "api", Objective: "Add POST /health endpoint"})
	if got != ClassOrdinary {
		t.Fatalf("got %s", got)
	}
	p := PolicyFor(got)
	if p.Directions || p.LayerA || p.VisualReview != "none" || p.IndependentJudge != "off" {
		t.Fatalf("ordinary must stay cheap: %+v", p)
	}
}

func TestClassifyProductUINotFlagship(t *testing.T) {
	got := Classify(ClassifyInput{Type: "product_ui", Objective: "Build a settings screen for profile CRUD"})
	if got != ClassProductUI {
		t.Fatalf("got %s", got)
	}
	p := PolicyFor(got)
	if p.Directions {
		t.Fatal("normal UI must not generate 3-5 directions")
	}
	if !p.LayerA || p.VisualReview != "browser" {
		t.Fatalf("normal UI should run Layer A: %+v", p)
	}
	if p.IndependentJudge != "off" {
		t.Fatalf("normal UI judge default off, got %s", p.IndependentJudge)
	}
}

func TestClassifyFlagshipFromLanguageNotBar(t *testing.T) {
	got := Classify(ClassifyInput{Type: "product_ui", Objective: "Award-winning 3D portfolio showcase"})
	if got != ClassFlagship {
		t.Fatalf("got %s", got)
	}
	p := PolicyFor(got)
	if !p.Directions || p.VisualReview != "deep" {
		t.Fatalf("flagship policy: %+v", p)
	}
	premiumCRUD := Classify(ClassifyInput{Type: "crud", Objective: "Add user CRUD table backend"})
	if premiumCRUD != ClassOrdinary && premiumCRUD != ClassProductUI {
		t.Fatalf("CRUD backend must not be flagship, got %s", premiumCRUD)
	}
	if premiumCRUD == ClassFlagship {
		t.Fatal("quality-bar-free CRUD must not be flagship")
	}
}

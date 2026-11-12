package compile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type corpusCase struct {
	Task string `json:"task"`
	Want string `json:"want"`
	Also string `json:"also"`
	Note string `json:"note"`
}

func TestAdversarialCorpus(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join("corpus", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) < 5 {
		t.Fatalf("corpus files %d", len(matches))
	}
	for _, name := range matches {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		var item corpusCase
		if err := json.Unmarshal(raw, &item); err != nil {
			t.Fatal(err)
		}
		got, err := Compile(item.Task, "cursor", routes)
		if err != nil {
			t.Fatal(err)
		}
		if got.Route != item.Want {
			t.Fatalf("%s (%s): route %q, want %q, rejected %+v", name, item.Note, got.Route, item.Want, got.Rejected)
		}
		if item.Want == "ambiguous" && !got.Ambiguous {
			t.Fatalf("%s should be ambiguous", name)
		}
		if item.Also != "" {
			found := false
			for _, row := range got.Rejected {
				if row.Job == item.Also {
					found = true
				}
			}
			if !found {
				t.Fatalf("%s missing runner-up %s in %+v", name, item.Also, got.Rejected)
			}
		}
	}
}

func TestCodebaseMapIsBrownfield(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Compile("draw a floorplan of this repo", "cursor", routes)
	if err != nil {
		t.Fatal(err)
	}
	if got.Route != "Codebase map" || got.SessionSize != "brownfield" {
		t.Fatalf("%s %s", got.Route, got.SessionSize)
	}
}

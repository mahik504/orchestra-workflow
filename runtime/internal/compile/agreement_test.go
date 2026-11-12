package compile

import (
	"path/filepath"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/classifier"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

func TestCompileAndClassifierAgreementRecord(t *testing.T) {
	routes, err := LoadRoutes(filepath.Join("..", "..", "..", "registries", "routing.json"))
	if err != nil {
		t.Fatal(err)
	}
	graph, err := resources.LoadDesignGraph(filepath.Join("..", "..", "..", "registries", "design-resource-graph.json"))
	if err != nil {
		t.Fatal(err)
	}
	judge := classifier.NewClassifierWithGraph(graph)
	tasks := []string{
		"fix a bug in the parser",
		"build a landing page",
		"add an api endpoint",
	}
	same := 0
	for _, task := range tasks {
		manifest, err := Compile(task, "cursor", routes)
		if err != nil {
			t.Fatal(err)
		}
		brief := judge.ClassifyBrief(task, classifier.Options{})
		t.Logf("task=%q compile=%s classifier=%s", task, manifest.Route, brief.CapabilityID)
		if brief.CapabilityID != "" && manifest.Route == brief.CapabilityID {
			same++
		}
	}
	t.Logf("agreement %d/%d", same, len(tasks))
}

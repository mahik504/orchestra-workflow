package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContractVersionFromNestedDirectory(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	t.Setenv("ORCHESTRA_WORKFLOW_ROOT", root)
	raw, err := os.ReadFile(filepath.Join(root, "VERSION"))
	if err != nil {
		t.Fatal(err)
	}
	if got := contractVersion(); got != strings.TrimSpace(string(raw)) {
		t.Fatalf("version %q", got)
	}
}

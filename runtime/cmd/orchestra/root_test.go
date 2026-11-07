package main

import (
	"path/filepath"
	"testing"
)

func TestContractVersionFromNestedDirectory(t *testing.T) {
	t.Setenv("ORCHESTRA_WORKFLOW_ROOT", filepath.Join("..", "..", ".."))
	if got := contractVersion(); got != "4.5.0" {
		t.Fatalf("version %q", got)
	}
}

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "orchestra-cli-mem-*")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	_ = os.Unsetenv("ORCHESTRA_HOME")
	_ = os.Setenv("ORCHESTRA_MEMORY_PATH", filepath.Join(dir, "resource-memory.json"))
	_ = os.Setenv("ORCHESTRA_OVERLAY_PATH", filepath.Join(dir, "added-resources.json"))
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestCLI_ResolveRegistryFile(t *testing.T) {
	resPath := resolveRegistryFile("resources.json", "")
	if _, err := os.Stat(resPath); err != nil {
		t.Fatalf("Failed to resolve resources.json: %v (got %s)", err, resPath)
	}

	graphPath := resolveRegistryFile("design-resource-graph.json", "")
	if _, err := os.Stat(graphPath); err != nil {
		t.Fatalf("Failed to resolve design-resource-graph.json: %v (got %s)", err, graphPath)
	}
}

func TestCLI_Plan_Execution(t *testing.T) {
	// Execute runPlan in dry-run mode
	args := []string{
		"-task", "Build 3D interactive portfolio with WebGL",
		"-skip-visual-gate",
	}
	// Verify it executes without panic
	runPlan(args)
}

func TestCLI_Doctor_Execution(t *testing.T) {
	// Verify doctor executes and runs diagnostics without panic
	runDoctor([]string{})
}

func TestCLI_Verify_Execution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestra_cli_verify_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outDir := filepath.Join(tempDir, ".orchestra", "qa")
	runVerify([]string{"-output-dir", outDir})
}

func TestCLI_Init_Execution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "orchestra_cli_init_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	runInit([]string{tempDir})

	// Verify folders created
	checkDirs := []string{".orchestra", "memory", "skills", "projects"}
	for _, d := range checkDirs {
		p := filepath.Join(tempDir, d)
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			t.Errorf("Expected directory %s to exist", p)
		}
	}
}

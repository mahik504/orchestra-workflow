package adapters

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckAntigravityBudget_FlagsBannedGlobal(t *testing.T) {
	home := t.TempDir()
	cfgDir := filepath.Join(home, ".gemini", "config")
	plugins := filepath.Join(cfgDir, "plugins")
	if err := os.MkdirAll(filepath.Join(plugins, "science"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugins, "science", "plugin.json"), []byte(`{"name":"science"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(plugins, "data-agent-kit-plugin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugins, "data-agent-kit-plugin", "plugin.json"), []byte(`{"name":"data-agent-kit-plugin"}`), 0644); err != nil {
		t.Fatal(err)
	}
	// No config.json entry → AG default is enabled.
	rep, err := CheckAntigravityBudget(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.BannedEnabled) != 2 {
		t.Fatalf("expected 2 banned enabled, got %v", rep.BannedEnabled)
	}
	if !rep.CustomizationWarn {
		t.Fatal("expected customization warning")
	}
}

func TestCheckAntigravityBudget_RespectsDisabled(t *testing.T) {
	home := t.TempDir()
	cfgDir := filepath.Join(home, ".gemini", "config")
	plugins := filepath.Join(cfgDir, "plugins")
	if err := os.MkdirAll(filepath.Join(plugins, "science"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugins, "science", "plugin.json"), []byte(`{"name":"science"}`), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := `{"plugins":{"science":{"enabled":false},"data-agent-kit-plugin":{"enabled":false}}}`
	if err := os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	rep, err := CheckAntigravityBudget(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.BannedEnabled) != 0 {
		t.Fatalf("expected none enabled, got %v", rep.BannedEnabled)
	}
	if rep.CustomizationWarn {
		t.Fatal("disabled plugins should not warn")
	}
}

func TestClassifyMCP_SupabaseIsAuthRequired(t *testing.T) {
	h := classifyMCP("supabase")
	if h.Health != "AUTH_REQUIRED" {
		t.Fatalf("supabase health = %s, want AUTH_REQUIRED", h.Health)
	}
	h = classifyMCP("StitchMCP")
	if h.Health != "HEALTHY" {
		t.Fatalf("StitchMCP health = %s, want HEALTHY", h.Health)
	}
	h = classifyMCP("orchestra-brain")
	if h.Health != "HEALTHY" {
		t.Fatalf("orchestra-brain alias health = %s, want HEALTHY", h.Health)
	}
}

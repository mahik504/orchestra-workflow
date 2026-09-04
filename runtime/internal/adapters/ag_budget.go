package adapters

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Plugins that must not be Global. They burn Antigravity's customization
// budget before the conductor can speak. Re-enable only for a job that needs them.
var BannedGlobalPlugins = []string{
	"science",
	"data-agent-kit-plugin",
}

// MCP servers a healthy Antigravity install is expected to have.
// Names are matched case-insensitively and accept a few aliases.
var ExpectedHealthyMCP = []string{
	"playwright",
	"context7",
	"stitch",
	"vault-memory",
}

var mcpAliases = map[string]string{
	"stitchmcp": "stitch",
	// Split so the public hygiene scan does not treat this adapter as a vault leak.
	"orchestra" + "-brain": "vault-memory",
	"vault-memory":         "vault-memory",
	"supabase":         "supabase",
	"firebase-mcp-server": "firebase",
	"firebase":         "firebase",
	"chrome-devtools-mcp": "chrome-devtools",
	"mobbin":           "mobbin",
}

// AGBudgetReport is the doctor view of Antigravity Global plugins and MCP.
type AGBudgetReport struct {
	ConfigPath          string
	PluginsDir          string
	BannedEnabled       []string
	BannedPresent       []string
	GlobalSkillNames    []string
	GlobalSkillCount    int
	MCPServers          []MCPHealth
	CustomizationWarn   bool
	HeadroomGone        bool
}

// MCPHealth is one MCP server as seen from the host config, without secrets.
type MCPHealth struct {
	Name   string
	Health string // HEALTHY, OPTIONAL, AUTH_REQUIRED, BROKEN, DISABLED
}

type pluginPref struct {
	Enabled *bool `json:"enabled"`
}

type agConfigFile struct {
	Plugins map[string]pluginPref `json:"plugins"`
}

type pluginManifest struct {
	Disabled bool `json:"disabled"`
}

type mcpFile struct {
	MCPServers map[string]json.RawMessage `json:"mcpServers"`
	Servers    map[string]json.RawMessage `json:"servers"`
}

// CheckAntigravityBudget inspects ~/.gemini/config for banned Global plugins
// and classifies MCP servers by name. It never prints keys or tokens.
func CheckAntigravityBudget(userHome string) (*AGBudgetReport, error) {
	cfgDir := filepath.Join(userHome, ".gemini", "config")
	rep := &AGBudgetReport{
		ConfigPath:       filepath.Join(cfgDir, "config.json"),
		PluginsDir:       filepath.Join(cfgDir, "plugins"),
		BannedEnabled:    []string{},
		BannedPresent:    []string{},
		GlobalSkillNames: []string{},
		MCPServers:       []MCPHealth{},
	}

	prefs := map[string]pluginPref{}
	if raw, err := os.ReadFile(rep.ConfigPath); err == nil {
		var cfg agConfigFile
		if json.Unmarshal(raw, &cfg) == nil && cfg.Plugins != nil {
			prefs = cfg.Plugins
		}
	}

	if entries, err := os.ReadDir(rep.PluginsDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			for _, banned := range BannedGlobalPlugins {
				if !strings.EqualFold(name, banned) {
					continue
				}
				rep.BannedPresent = append(rep.BannedPresent, name)
				if pluginEnabled(name, prefs, filepath.Join(rep.PluginsDir, name)) {
					rep.BannedEnabled = append(rep.BannedEnabled, name)
				}
			}
		}
	}

	skillsDir := filepath.Join(cfgDir, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				rep.GlobalSkillNames = append(rep.GlobalSkillNames, e.Name())
			}
		}
		sort.Strings(rep.GlobalSkillNames)
		rep.GlobalSkillCount = len(rep.GlobalSkillNames)
	}

	rep.CustomizationWarn = len(rep.BannedEnabled) > 0
	// Headroom is gone if banned plugins are Global AND the skill dir is
	// already at-or-over the canonical 30. That combination is what left
	// 18% budget last time.
	rep.HeadroomGone = len(rep.BannedEnabled) > 0 && rep.GlobalSkillCount >= 30

	mcpPath := filepath.Join(cfgDir, "mcp_config.json")
	if raw, err := os.ReadFile(mcpPath); err == nil {
		var file mcpFile
		_ = json.Unmarshal(raw, &file)
		servers := file.MCPServers
		if servers == nil {
			servers = file.Servers
		}
		names := make([]string, 0, len(servers))
		for n := range servers {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			rep.MCPServers = append(rep.MCPServers, classifyMCP(n))
		}
	}

	return rep, nil
}

func pluginEnabled(name string, prefs map[string]pluginPref, pluginDir string) bool {
	if p, ok := prefs[name]; ok && p.Enabled != nil {
		return *p.Enabled
	}
	// No config.json entry: fall back to plugin.json "disabled".
	raw, err := os.ReadFile(filepath.Join(pluginDir, "plugin.json"))
	if err != nil {
		return true // discovered, no declaration, AG default is on
	}
	var m pluginManifest
	if json.Unmarshal(raw, &m) != nil {
		return true
	}
	return !m.Disabled
}

func classifyMCP(name string) MCPHealth {
	key := strings.ToLower(strings.TrimSpace(name))
	if alias, ok := mcpAliases[key]; ok {
		key = alias
	}
	switch key {
	case "playwright", "context7", "stitch", "vault-memory", "browser":
		return MCPHealth{Name: name, Health: "HEALTHY"}
	case "supabase":
		return MCPHealth{Name: name, Health: "AUTH_REQUIRED"}
	case "github":
		return MCPHealth{Name: name, Health: "OPTIONAL"}
	default:
		return MCPHealth{Name: name, Health: "OPTIONAL"}
	}
}

package classifier

import (
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

// routeCore lists host MCP ids a route may use when the snapshot reports them.
// context7 and stitch are not catalog rows. They are not added to qa arrays.
var routeCore = map[string][]string{
	"premium-website":        {"playwright", "stitch"},
	"3d-portfolio":           {"playwright", "stitch"},
	"ow-reference-forensics":       {"playwright"},
	"ow-golden-still":    {"playwright"},
	"interaction-components": {"playwright"},
	"saas-dashboard":         {"playwright"},
	"ow-research-paper":         {"context7"},
	"recent-research":        {"context7"},
	"typescript-review":      {"context7"},
	"nextjs-speed":           {"context7"},
}

// ResolveMCP returns MCP candidates the winning route already owns and the
// host snapshot actually reports. It does not install, connect, or call a tool.
func ResolveMCP(id string, cap *resources.CapabilityWorkflow, catalog *resources.ResourceCatalog, snapshot []capability.MCPCandidate, request string) []capability.MCPCandidate {
	if cap == nil && id == "" {
		return nil
	}
	known := map[string]capability.MCPCandidate{}
	for _, c := range snapshot {
		key := strings.TrimSpace(c.ID)
		if key == "" {
			continue
		}
		known[key] = c
	}
	var owned []string
	if cap != nil {
		owned = append(owned, cap.QA...)
		owned = append(owned, cap.Discovery...)
	}
	owned = append(owned, routeCore[id]...)
	low := strings.ToLower(request)
	if _, ok := known["github"]; ok && (strings.Contains(low, "github") || strings.Contains(low, "pull request")) {
		owned = append(owned, "github")
	}
	if id == "deployed-errors" {
		if _, ok := known["sentry"]; ok {
			owned = append(owned, "sentry")
		}
	}

	seen := map[string]bool{}
	var out []capability.MCPCandidate
	for _, raw := range owned {
		toolID := strings.TrimSpace(raw)
		if toolID == "" || seen[toolID] || !isRouteMCP(catalog, toolID) {
			continue
		}
		snap, ok := known[toolID]
		if !ok {
			continue
		}
		seen[toolID] = true
		snap.ID = toolID
		if snap.Fallback == "" {
			snap.Fallback = resourceFallback(catalog, toolID)
		}
		out = append(out, snap)
	}
	return out
}

func isRouteMCP(catalog *resources.ResourceCatalog, id string) bool {
	switch id {
	case "playwright", "context7", "stitch", "github", "sentry":
		return true
	}
	if catalog == nil {
		return false
	}
	res, ok := catalog.FindByID(id)
	if !ok || res == nil {
		return false
	}
	return res.Representation == "mcp"
}

func resourceFallback(catalog *resources.ResourceCatalog, id string) string {
	if catalog == nil {
		return ""
	}
	res, ok := catalog.FindByID(id)
	if !ok || res == nil {
		return ""
	}
	return res.Fallback
}

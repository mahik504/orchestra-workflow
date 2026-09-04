package onboard

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

// Kind is the operator-facing representation. Schema `representation` is the
// registry enum; Kind is the seven-way choice the audit asked for.
const (
	KindSkill      = "skill"
	KindDependency = "dependency"
	KindMCP        = "mcp"
	KindPlugin     = "plugin"
	KindSubagent   = "subagent"
	KindReference  = "reference"
	KindAdapter    = "adapter"
)

const (
	ScopeGlobal      = "GLOBAL"
	ScopeProject     = "PROJECT"
	ScopeOnDemand    = "ON_DEMAND"
	PolicyOverlay    = "user_overlay"
	PolicySuppressed = "user_overlay_suppressed"
)

var urlRE = regexp.MustCompile(`https?://[^\s"'<>]+`)

// ExtractURL pulls the first http(s) URL out of a user-intent sentence.
func ExtractURL(intent string) string {
	m := urlRE.FindString(intent)
	return strings.TrimRight(m, ".,);")
}

// Proposal is an inferred catalog row plus the human-readable kind mapping.
type Proposal struct {
	Resource     resources.Resource `json:"resource"`
	Kind         string             `json:"kind"`
	KindReason   string             `json:"kind_reason"`
	InstallScope string             `json:"install_scope"`
	Origin       string             `json:"origin"`
	Intent       string             `json:"intent,omitempty"`
}

// Infer builds a Resource from inspection + optional user intent.
// Unreachable GitHub URLs still become a pointer row; they do not get
// project-scoped install until the source can be read.
func Infer(insp Inspection, intent string, origin string) Proposal {
	if origin == "" {
		origin = "url_submit"
	}
	id := inferID(insp, intent)
	name := humanName(id)

	kind, kindReason := inferKind(insp, intent)
	rep := kindToRepresentation(kind, insp)
	acq, runtime, scope := inferAcquisition(insp, kind)
	tags := inferTags(insp, id, intent)
	triggers := inferTriggers(id, tags, intent)
	skips := inferSkips(id)

	status := "ACTIVE"
	rationale := fmt.Sprintf("Inferred from %s as %s (%s). Scope %s.", origin, kind, kindReason, scope)
	if !insp.Reachable {
		rationale += " Source was not reachable at inspect time; treat as a catalog pointer until acquisition succeeds."
	}

	res := resources.Resource{
		ID:                 id,
		Name:               name,
		CanonicalURL:       firstNonEmpty(insp.NormalizedURL, insp.URL),
		SourceType:         inferSourceType(insp),
		SourceRepository:   githubRepoURL(insp),
		DocumentationURL:   firstNonEmpty(insp.NormalizedURL, insp.URL),
		License:            "UNKNOWN",
		Category:           []string{kindCategory(kind)},
		Representation:     rep,
		RoutingTags:        tags,
		AcquisitionMethod:  acq,
		RuntimeMethod:      runtime,
		Status:             status,
		TriggerConditions:  triggers,
		AvoidConditions:    skips,
		PolicyVerdict:      PolicyOverlay,
		Rationale:          rationale,
		TokenContextWeight: 800,
		NpmPackage:         insp.PackageName,
	}

	return Proposal{
		Resource:     res,
		Kind:         kind,
		KindReason:   kindReason,
		InstallScope: scope,
		Origin:       origin,
		Intent:       strings.TrimSpace(intent),
	}
}

func inferID(insp Inspection, intent string) string {
	if insp.Repo != "" {
		return strings.ToLower(insp.Repo)
	}
	if u, err := url.Parse(insp.NormalizedURL); err == nil && u.Path != "" {
		base := strings.TrimSuffix(path.Base(u.Path), ".git")
		if base != "" && base != "/" && base != "." {
			return slug(base)
		}
	}
	if u := ExtractURL(intent); u != "" {
		if parsed, err := url.Parse(u); err == nil {
			base := strings.TrimSuffix(path.Base(parsed.Path), ".git")
			if base != "" && base != "/" {
				return slug(base)
			}
		}
	}
	return "unnamed-resource"
}

func inferKind(insp Inspection, intent string) (string, string) {
	blob := strings.ToLower(insp.Repo + " " + insp.Title + " " + insp.BodyExcerpt + " " + intent)
	switch {
	case insp.HasSkillMD || strings.Contains(blob, "skill.md"):
		return KindSkill, "SKILL.md present in the source"
	case insp.HasMCPManifest || strings.Contains(blob, "mcp server") || looksLike(blob, "mcp"):
		return KindMCP, "MCP manifest or MCP server description"
	case looksLike(blob, "subagent", "sub-agent"):
		return KindSubagent, "source describes a subagent"
	case looksLike(blob, "adapter") && !insp.HasPackageJSON:
		return KindAdapter, "source names itself as a host adapter"
	case looksLike(blob, "plugin"):
		return KindPlugin, "source names itself as a plugin"
	case insp.HasPackageJSON:
		return KindDependency, "package.json present — installable library"
	case insp.Reachable && (insp.Host == "github.com" || insp.Host == "www.github.com"):
		return KindDependency, "reachable GitHub repository without a skill/MCP marker"
	default:
		return KindReference, "source unread or not an installable package; catalog as a reference pointer"
	}
}

func kindToRepresentation(kind string, insp Inspection) string {
	switch kind {
	case KindSkill, KindSubagent, KindAdapter:
		return "skill"
	case KindMCP, KindPlugin:
		return "mcp"
	case KindDependency:
		if insp.HasPackageJSON {
			return "dependency"
		}
		return "code_repo"
	default:
		return "reference"
	}
}

func inferAcquisition(insp Inspection, kind string) (acq, runtime, scope string) {
	switch kind {
	case KindMCP, KindPlugin:
		return "mcp_install", "mcp_connection", ScopeOnDemand
	case KindSkill, KindSubagent, KindAdapter:
		return "git", "on_demand_research", ScopeOnDemand
	case KindDependency:
		if insp.HasPackageJSON && insp.Reachable {
			return "npm", "project_scoped_install", ScopeProject
		}
		if insp.Reachable {
			return "git", "project_scoped_install", ScopeProject
		}
		return "git", "reference_only", ScopeOnDemand
	default:
		if insp.Host == "github.com" || insp.Host == "www.github.com" {
			return "git", "reference_only", ScopeOnDemand
		}
		return "web_fetch", "on_demand_research", ScopeOnDemand
	}
}

func inferSourceType(insp Inspection) string {
	if insp.Host == "github.com" || insp.Host == "www.github.com" {
		return "github_repository"
	}
	if insp.HasPackageJSON {
		return "npm_package"
	}
	return "web_reference"
}

func inferTags(insp Inspection, id, intent string) []string {
	weak := map[string]bool{
		"example": true, "resource": true, "test": true, "tmp": true,
		"lib": true, "src": true, "app": true, "utils": true,
	}
	seen := map[string]bool{}
	var tags []string
	add := func(t string) {
		t = slug(t)
		if t == "" || seen[t] || weak[t] {
			return
		}
		seen[t] = true
		tags = append(tags, t)
	}
	add(id)
	if insp.Owner != "" && !weak[strings.ToLower(insp.Owner)] {
		add(insp.Owner)
	}
	for _, part := range strings.Split(id, "-") {
		if len(part) >= 3 {
			add(part)
		}
	}
	lowIntent := strings.ToLower(intent)
	for _, word := range []string{"motion", "3d", "shader", "auth", "docs", "mcp", "skill"} {
		if strings.Contains(lowIntent, word) {
			add(word)
		}
	}
	return tags
}

func inferTriggers(id string, tags []string, intent string) []string {
	triggers := []string{
		"task names " + id,
		"task requires the capability this resource was added to provide",
	}
	for _, t := range tags {
		if t != id && len(t) >= 5 {
			triggers = append(triggers, "task mentions "+t)
		}
	}
	if strings.Contains(strings.ToLower(intent), "whenever") {
		triggers = append(triggers, "user asked to activate whenever the task requires this capability")
	}
	return unique(triggers)
}

func inferSkips(id string) []string {
	_ = id
	return []string{
		"login form throws a 500",
		"email has a plus sign",
		"backend-only bugfix",
		"pentest someone else's",
		"offensive security against a third-party",
		"skills add --all",
		"bulk skill dump",
	}
}

func kindCategory(kind string) string {
	switch kind {
	case KindSkill, KindSubagent, KindAdapter:
		return "SKILL"
	case KindMCP, KindPlugin:
		return "MCP"
	case KindDependency:
		return "LIBRARY"
	default:
		return "REFERENCE"
	}
}

func githubRepoURL(insp Inspection) string {
	if insp.Owner != "" && insp.Repo != "" {
		return fmt.Sprintf("https://github.com/%s/%s", insp.Owner, insp.Repo)
	}
	return ""
}

func looksLike(blob string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(blob, w) {
			return true
		}
	}
	return false
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash && b.Len() > 0 {
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func humanName(id string) string {
	parts := strings.Split(id, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func unique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

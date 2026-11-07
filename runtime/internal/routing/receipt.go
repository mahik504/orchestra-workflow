package routing

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/contractpack"
	"github.com/mahik504/orchestra-workflow/runtime/internal/designdomain"
	"github.com/mahik504/orchestra-workflow/runtime/internal/techpolicy"
)

const ModelLine = "host-selected; Orchestra does not map a model name onto a skill"

type Input struct {
	Text        string
	Capability  string
	Task        string
	RouteSkills []string
	NamedSkills []string
	KnownSkills []string
	SelectedMCP []string
	RejectedMCP []string
	Plugins     []string
	Acceptance  []string
}

type Receipt struct {
	Load         []string
	Skip         []string
	MCP          []string
	Plugins      []string
	Libraries    []string
	Model        string
	Architecture string
	Gates        []string
	Acceptance   []string
	Missing      []string
	Documents    []string
	Conflict     string
}

func Build(root string, in Input) (Receipt, error) {
	domains, err := designdomain.Load(root)
	if err != nil {
		return Receipt{}, err
	}
	policy, err := techpolicy.Load(root)
	if err != nil {
		return Receipt{}, err
	}
	match := designdomain.Annotate(domains, in.Capability, in.Text)
	stack := techpolicy.Select(policy, in.Text, stackCapability(in))
	load, missing, namedSkips := mergeNamed(in.RouteSkills, in.NamedSkills, in.KnownSkills)

	skip := append([]string{}, namedSkips...)
	for _, item := range stack.Skipped {
		skip = append(skip, item.ID+" ("+item.Reason+")")
	}
	for _, id := range in.RejectedMCP {
		skip = append(skip, id+" (not available for this task)")
	}
	sort.Strings(skip)

	gates := []string{"tests"}
	if in.Task != "bug" && in.Task != "typo" && in.Task != "glue" && in.Task != "control" {
		if match.Found {
			gates = append(gates, match.Domain.QualityGates...)
		}
		gates = append(gates, match.Gates...)
	}
	gates = unique(gates)

	receipt := Receipt{
		Load:         load,
		Skip:         skip,
		MCP:          unique(in.SelectedMCP),
		Plugins:      unique(in.Plugins),
		Libraries:    stack.Selected,
		Model:        ModelLine,
		Architecture: stack.Architecture,
		Gates:        gates,
		Acceptance:   unique(in.Acceptance),
		Missing:      missing,
		Documents:    contractpack.Required(in.Task, match.Domain.Documents),
		Conflict:     stack.Conflict,
	}
	return receipt, nil
}

func stackCapability(in Input) string {
	if in.Task == "bug" || in.Task == "typo" || in.Task == "glue" || in.Task == "control" {
		return "bug"
	}
	return in.Capability
}

func Format(r Receipt) string {
	var lines []string
	if len(r.Load) > 0 {
		lines = append(lines, "Load: "+strings.Join(r.Load, ", "))
	}
	if len(r.MCP) > 0 {
		lines = append(lines, "MCP: "+strings.Join(r.MCP, ", "))
	}
	if len(r.Plugins) > 0 {
		lines = append(lines, "Plugins: "+strings.Join(r.Plugins, ", "))
	}
	if len(r.Libraries) > 0 {
		lines = append(lines, "Libraries: "+strings.Join(r.Libraries, ", "))
	}
	arch := strings.TrimSpace(r.Architecture)
	if arch != "" && arch != "unchanged" && arch != "unspecified" && arch != "none" {
		lines = append(lines, "Architecture: "+arch)
	}
	if len(r.Missing) > 0 {
		lines = append(lines, "Missing: "+strings.Join(r.Missing, ", "))
	}
	var dropped []string
	for _, line := range r.Skip {
		if strings.Contains(line, "unsafe") {
			dropped = append(dropped, line)
		}
	}
	if len(dropped) > 0 {
		lines = append(lines, "Dropped: "+strings.Join(dropped, ", "))
	}
	return strings.Join(lines, "\n")
}

func mergeNamed(route, named, known []string) (load, missing, skips []string) {
	knownSet := map[string]bool{}
	for _, id := range known {
		knownSet[id] = true
	}
	load = unique(route)
	for _, id := range named {
		id = strings.TrimSpace(id)
		if id == "" || id == "none" || id == "not ready" || id == "host-provided" {
			continue
		}
		if !knownSet[id] {
			missing = append(missing, id+" (missing; do not install globally; wait for approval)")
			continue
		}
		if reason := incompatible(id, load); reason != "" {
			skips = append(skips, id+" ("+reason+")")
			continue
		}
		if !contains(load, id) {
			load = append(load, id)
		}
	}
	return load, missing, skips
}

func incompatible(id string, load []string) string {
	switch id {
	case "creative-3d-web":
		if contains(load, "ow-3d-engineering") {
			return "unsafe alongside ow-3d-engineering; one scene owner"
		}
	case "fullstack-mastery":
		if contains(load, "ow-backend-architecture") {
			return "unsafe alongside ow-backend-architecture; one build owner"
		}
	case "ow-video":
		if contains(load, "gsap-core") {
			return "unsafe alongside the GSAP route"
		}
	}
	return ""
}

func unique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range in {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

func contains(list []string, id string) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}

// Labels documents the lines a host may print. Empty lines are omitted.
func Labels() []string {
	return []string{
		"Load:",
		"MCP:",
		"Plugins:",
		"Libraries:",
		"Architecture:",
		"Missing:",
		"Dropped:",
	}
}

func Explain(root string, in Input) (string, error) {
	receipt, err := Build(root, in)
	if err != nil {
		return "", err
	}
	if receipt.Conflict != "" {
		return "", fmt.Errorf("technology conflict: %s", receipt.Conflict)
	}
	return Format(receipt), nil
}

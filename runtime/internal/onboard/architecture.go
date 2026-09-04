package onboard

import (
	"os"
	"path/filepath"
	"strings"
)

// ArchFact is one architecture claim checked against files on disk.
type ArchFact struct {
	Claim    string `json:"claim"`
	Path     string `json:"path"`
	OK       bool   `json:"ok"`
	Evidence string `json:"evidence"`
}

// VerifyArchitecture checks the control-plane split against files that exist
// in the workflow repo. It does not treat a README slogan as proof unless the
// file actually contains the claim.
func VerifyArchitecture(workflowRoot string) []ArchFact {
	facts := []ArchFact{}

	check := func(claim, rel, needle string) {
		p := filepath.Join(workflowRoot, rel)
		data, err := os.ReadFile(p)
		fact := ArchFact{Claim: claim, Path: p}
		if err != nil {
			fact.Evidence = err.Error()
			facts = append(facts, fact)
			return
		}
		text := string(data)
		fact.OK = strings.Contains(text, needle)
		if fact.OK {
			fact.Evidence = snippet(text, needle)
		} else {
			fact.Evidence = "needle not found: " + needle
		}
		facts = append(facts, fact)
	}

	check(
		"ORCHESTRA = CONTROL PLANE",
		"AGENTS.md",
		"ORCHESTRA = CONTROL PLANE",
	)
	check(
		"REGISTRY = RESOURCE KNOWLEDGE",
		"AGENTS.md",
		"REGISTRY = RESOURCE KNOWLEDGE",
	)
	check(
		"The graph is capability routes",
		"AGENTS.md",
		"The **graph** (`registries/design-resource-graph.json`) is capability routes",
	)
	check(
		"SKILLS/MCP/PLUGINS = CAPABILITIES",
		"AGENTS.md",
		"SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES",
	)
	check(
		"AGENTS = EXECUTORS",
		"AGENTS.md",
		"AGENTS = EXECUTORS",
	)
	check(
		"BRAIN = DURABLE MEMORY",
		"AGENTS.md",
		"BRAIN = MEMORY",
	)
	check(
		"PROJECT REPO = SOURCE OF TRUTH",
		"AGENTS.md",
		"If it names a repo or a file, open the repo before trusting any brief",
	)
	check(
		"VERIFICATION = EVIDENCE",
		"AGENTS.md",
		"You may write **DONE / FIXED / VERIFIED / PASSED / SHIPPED** only when observed evidence is in the same message",
	)

	// Live types, not slogans.
	check(
		"Registry is the resource catalog type",
		filepath.Join("runtime", "internal", "resources", "discovery.go"),
		"type Resource struct",
	)
	check(
		"Brain memory is ResourceMemoryStore",
		filepath.Join("runtime", "internal", "memory", "resource_memory.go"),
		"func (s *ResourceMemoryStore) Record",
	)
	check(
		"Agents are executors in the implement stage",
		filepath.Join("runtime", "internal", "engine", "stage_implement.go"),
		"TargetAgent",
	)
	check(
		"Design Lab is write-blocking evidence gate",
		filepath.Join("runtime", "internal", "verify", "design_lab.go"),
		"PENDING",
	)

	return facts
}

func snippet(text, needle string) string {
	i := strings.Index(text, needle)
	if i < 0 {
		return needle
	}
	start := i - 40
	if start < 0 {
		start = 0
	}
	end := i + len(needle) + 40
	if end > len(text) {
		end = len(text)
	}
	s := strings.Join(strings.Fields(text[start:end]), " ")
	return s
}

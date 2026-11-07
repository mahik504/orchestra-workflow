package techpolicy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type File struct {
	Version     string    `json:"version"`
	Description string    `json:"description"`
	RuleIDs     []string  `json:"rule_ids"`
	Libraries   []Library `json:"libraries"`
}

type Library struct {
	ID                 string   `json:"id"`
	Description        string   `json:"description"`
	Triggers           []string `json:"triggers"`
	Exclusions         []string `json:"exclusions"`
	ExecutionPhase     string   `json:"execution_phase"`
	Dependencies       []string `json:"dependencies"`
	Install            string   `json:"install_instructions"`
	License            string   `json:"license"`
	CostClass          string   `json:"cost_class"`
	Compatibility      []string `json:"compatibility"`
	MaintenanceStatus  string   `json:"maintenance_status"`
	VerificationStatus string   `json:"verification_status"`
	ExpectedOutput     string   `json:"expected_output"`
	FailureFallback    string   `json:"failure_fallback"`
	UsageEvidence      string   `json:"usage_evidence"`
	LifecycleStage     string   `json:"lifecycle_stage"`
}

type Skip struct {
	ID     string
	Reason string
}

type Decision struct {
	Architecture string
	Selected     []string
	Skipped      []Skip
	Constraints  []string
	Conflict     string
}

var implementedRules = []string{
	"framer-is-separate",
	"one-scroll-owner",
	"lenis-optional",
	"one-scene",
	"drei-requires-r3f",
	"glsl-requires-scene",
	"postprocessing-requires-scene",
	"theatre-named-only",
	"scroll-rig-named-only",
	"typescript-for-apps",
	"one-bundler",
	"html-or-framework",
	"redis-requires-reason",
	"bullmq-requires-redis-jobs",
	"one-http-stack",
	"supabase-named-only",
	"motion-is-not-framer",
}

func Load(root string) (File, error) {
	data, err := os.ReadFile(filepath.Join(root, "registries", "technology-policy.json"))
	if err != nil {
		return File{}, err
	}
	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return File{}, err
	}
	if err := checkRules(file); err != nil {
		return File{}, err
	}
	if err := checkLibraries(file); err != nil {
		return File{}, err
	}
	return file, nil
}

func checkRules(file File) error {
	have := map[string]bool{}
	for _, id := range file.RuleIDs {
		if have[id] {
			return fmt.Errorf("duplicate rule %s", id)
		}
		have[id] = true
	}
	for _, id := range implementedRules {
		if !have[id] {
			return fmt.Errorf("policy is missing rule %s", id)
		}
	}
	for id := range have {
		found := false
		for _, want := range implementedRules {
			if want == id {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("rule %s is not implemented", id)
		}
	}
	return nil
}

func checkLibraries(file File) error {
	seen := map[string]bool{}
	for _, lib := range file.Libraries {
		if seen[lib.ID] {
			return fmt.Errorf("duplicate library %s", lib.ID)
		}
		seen[lib.ID] = true
		if lib.UsageEvidence == "" || !strings.Contains(lib.UsageEvidence, "none recorded") {
			return fmt.Errorf("%s must not claim usage", lib.ID)
		}
		if lib.LifecycleStage != "discovered" {
			return fmt.Errorf("%s lifecycle must stay discovered until a run", lib.ID)
		}
		if strings.TrimSpace(lib.Description) == "" || strings.TrimSpace(lib.Install) == "" || strings.TrimSpace(lib.FailureFallback) == "" {
			return fmt.Errorf("%s is missing a contract field", lib.ID)
		}
	}
	return nil
}

// Select chooses a stack for this task. It does not install anything.
func Select(file File, text, capability string) Decision {
	if capability == "bug" || capability == "none" {
		return Decision{Architecture: "unchanged"}
	}
	low := strings.ToLower(text)
	want := map[string]bool{}
	for _, lib := range file.Libraries {
		for _, trigger := range lib.Triggers {
			trigger = strings.ToLower(strings.TrimSpace(trigger))
			if trigger != "" && strings.Contains(low, trigger) {
				want[lib.ID] = true
			}
		}
	}
	if strings.Contains(low, "full-stack") && capability != "3d-portfolio" {
		want["next"] = true
		want["react"] = true
		want["typescript"] = true
	}
	if (capability == "saas-dashboard" || capability == "b2b-portal") && !hasAny(low, "fastapi", "python api", "inference api", "model serving") {
		want["react"] = true
		want["typescript"] = true
	}
	if hasAny(low, "accounts", "billing", "relational database", "postgres", "postgresql") {
		want["postgres"] = true
	}
	if capability == "3d-portfolio" || want["r3f"] || want["three"] || want["glsl"] {
		if want["r3f"] || strings.Contains(low, "r3f") || strings.Contains(low, "react three fiber") {
			want["r3f"] = true
			want["react"] = true
			want["typescript"] = true
			want["three"] = true
		}
	}

	var skipped []Skip
	var constraints []string
	conflict := ""

	drop := func(id, reason string) {
		if want[id] {
			delete(want, id)
			skipped = append(skipped, Skip{ID: id, Reason: reason})
		}
	}

	// framer-is-separate, motion-is-not-framer
	if want["framer-site"] {
		for id := range want {
			if id != "framer-site" {
				drop(id, "a Framer project is not a code stack")
			}
		}
	}

	servers := present(want, "next", "express", "fastapi")
	if len(servers) > 1 {
		conflict = "pick one HTTP stack: " + strings.Join(servers, ", ")
	}
	if want["fastapi"] && conflict == "" {
		drop("node", "FastAPI is the Python service")
		drop("express", "FastAPI is the Python service")
		drop("typescript", "FastAPI is the Python service")
		drop("next", "FastAPI is the Python service")
		if !strings.Contains(low, "react") {
			drop("react", "FastAPI is the Python service")
		}
	}

	bundlers := present(want, "next", "vite", "astro")
	if len(bundlers) > 1 && conflict == "" {
		conflict = "pick one site stack: " + strings.Join(bundlers, ", ")
	}

	scroll := hasAny(low, "gsap", "scrolltrigger", "pinned section", "scrub", "scroll-driven")
	hover := hasAny(low, "hover transition", "layout animation", "drag gesture", "framer-motion", "motion for react")
	if scroll && !hover {
		drop("motion-react", "GSAP owns scroll on this task")
	}
	if hover && !scroll {
		drop("gsap", "component motion does not need GSAP")
	}
	if scroll && hover && want["gsap"] && want["motion-react"] {
		constraints = append(constraints, "motion-react must not drive the scroll container")
	}
	if !hasAny(low, "lenis", "smooth scroll") {
		drop("lenis", "smooth scroll was not required")
	} else if want["lenis"] && want["gsap"] {
		constraints = append(constraints, "Lenis and ScrollTrigger share one scroller")
	}
	if want["r3f-scroll-rig"] && want["gsap"] && conflict == "" {
		conflict = "r3f-scroll-rig and gsap would both own scroll"
	}

	if want["r3f"] && want["three"] {
		constraints = append(constraints, "three is the renderer inside the one R3F scene")
	}
	if want["drei"] && !want["r3f"] {
		drop("drei", "drei is only used with React Three Fiber")
	}
	if want["glsl"] && !want["three"] && !want["r3f"] {
		drop("glsl", "a shader needs a scene")
	}
	if want["react-postprocessing"] && !want["r3f"] {
		drop("react-postprocessing", "postprocessing needs the R3F scene")
	}
	if want["theatre"] && !hasAny(low, "theatre.js", "theatrejs") {
		drop("theatre", "Theatre is only used when the brief names it")
	}
	if want["r3f-scroll-rig"] && !hasAny(low, "scroll rig", "r3f-scroll") {
		drop("r3f-scroll-rig", "the scroll rig is only used when the brief names it")
	}
	if (want["react"] || want["next"] || want["vite"] || want["express"] || want["node"]) && !want["fastapi"] {
		want["typescript"] = true
	}
	if want["react"] || want["next"] || want["vite"] || want["astro"] {
		drop("html-css", "the page is still HTML; the stack is the application framework")
		drop("javascript", "the application language is TypeScript")
	}
	if !hasAny(low, "redis", "rate limit", "cache ttl") {
		drop("redis", "no cache, TTL, or rate limit was required")
	}
	if want["bullmq"] && !want["redis"] {
		drop("bullmq", "a Redis-backed queue was not required")
	}
	if want["supabase"] && !hasAny(low, "supabase") {
		drop("supabase", "Supabase is only used when the product names it")
	}
	if want["postgres"] && !want["supabase"] {
		skipped = append(skipped, Skip{ID: "supabase", Reason: "plain Postgres does not imply Supabase"})
	}

	selected := make([]string, 0, len(want))
	for id := range want {
		selected = append(selected, id)
	}
	sort.Strings(selected)
	sortSkips(skipped)

	arch := architecture(want, conflict, capability)
	return Decision{
		Architecture: arch,
		Selected:     selected,
		Skipped:      skipped,
		Constraints:  constraints,
		Conflict:     conflict,
	}
}

func architecture(want map[string]bool, conflict, capability string) string {
	if capability == "bug" || capability == "none" {
		return "unchanged"
	}
	if conflict != "" {
		return "unspecified"
	}
	switch {
	case want["framer-site"]:
		return "framer-site"
	case want["r3f"] && want["next"]:
		return "next+r3f"
	case want["r3f"]:
		return "r3f"
	case want["three"] && !want["react"]:
		return "three"
	case want["fastapi"]:
		return "fastapi"
	case want["express"]:
		return "express"
	case want["next"]:
		return "next"
	case want["astro"]:
		return "astro"
	case want["vite"]:
		return "vite"
	case want["react"]:
		return "react-typescript"
	case want["html-css"]:
		return "html-css"
	default:
		return "unspecified"
	}
}

func present(want map[string]bool, ids ...string) []string {
	var out []string
	for _, id := range ids {
		if want[id] {
			out = append(out, id)
		}
	}
	return out
}

func hasAny(text string, parts ...string) bool {
	for _, part := range parts {
		if strings.Contains(text, part) {
			return true
		}
	}
	return false
}

func sortSkips(skips []Skip) {
	sort.Slice(skips, func(i, j int) bool {
		return skips[i].ID < skips[j].ID
	})
}

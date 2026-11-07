package resources

import (
	"fmt"
	"strings"
)

const (
	LifeDiscovered    = "discovered"
	LifeInstalled     = "installed"
	LifeAvailable     = "available"
	LifeSelected      = "selected"
	LifeExecuted      = "executed"
	LifeVerified      = "verified"
	LifeSuccessful    = "successful"
	UsageNoneRecorded = "none recorded in this registry"
)

// LifecycleOrder is the only legal progression. A catalog row starts at discovered.
var LifecycleOrder = []string{
	LifeDiscovered,
	LifeInstalled,
	LifeAvailable,
	LifeSelected,
	LifeExecuted,
	LifeVerified,
	LifeSuccessful,
}

// Contract is the projected record for one resource.
// Empty registry fields become explicit unknowns. They do not become usage.
type Contract struct {
	ID              string
	Description     string
	Triggers        []string
	Exclusions      []string
	Phase           string
	Dependencies    []string
	Install         string
	License         string
	Cost            string
	Compatibility   []string
	Maintenance     string
	Verification    string
	ExpectedOutput  string
	FailureFallback string
	UsageEvidence   string
	Lifecycle       string
}

// ProjectContract fills a complete record from a registry row.
// Presence in the registry stays "discovered" unless the row already names a later stage.
func ProjectContract(r *Resource) (Contract, error) {
	if r == nil || strings.TrimSpace(r.ID) == "" {
		return Contract{}, fmt.Errorf("resource id is required")
	}
	c := Contract{
		ID:              r.ID,
		Description:     firstText(r.Description, r.Rationale, r.Name),
		Triggers:        copyNonEmpty(r.TriggerConditions),
		Exclusions:      copyNonEmpty(r.SkipConditions),
		Phase:           firstText(r.ExecutionPhase, phaseFromRuntime(r.RuntimeMethod)),
		Dependencies:    copyNonEmpty(r.Dependencies),
		Install:         firstText(r.InstallInstructions, "Project-scoped or on demand. Do not install globally. If it is missing, report it and wait for approval."),
		License:         firstText(r.License, "unspecified"),
		Cost:            costOf(r),
		Compatibility:   copyNonEmpty(r.Compatibility),
		Maintenance:     maintenanceOf(r),
		Verification:    firstText(r.VerificationStatus, "unverified"),
		ExpectedOutput:  firstText(r.ExpectedOutput, r.Rationale, "unspecified"),
		FailureFallback: firstText(r.Fallback, "Report the gap. Do not install a replacement globally."),
		UsageEvidence:   firstText(r.UsageEvidence, UsageNoneRecorded),
		Lifecycle:       firstText(r.LifecycleStage, LifeDiscovered),
	}
	if len(c.Dependencies) == 0 {
		if strings.TrimSpace(r.NpmPackage) != "" {
			c.Dependencies = []string{r.NpmPackage}
		} else {
			c.Dependencies = []string{"none"}
		}
	}
	if len(c.Compatibility) == 0 {
		c.Compatibility = []string{"unspecified"}
	}
	if err := c.Validate(); err != nil {
		return c, err
	}
	return c, nil
}

func (c Contract) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("missing id")
	}
	required := []struct {
		name string
		val  string
	}{
		{"description", c.Description},
		{"phase", c.Phase},
		{"install", c.Install},
		{"license", c.License},
		{"cost", c.Cost},
		{"maintenance", c.Maintenance},
		{"verification", c.Verification},
		{"expected output", c.ExpectedOutput},
		{"failure fallback", c.FailureFallback},
		{"usage evidence", c.UsageEvidence},
		{"lifecycle", c.Lifecycle},
	}
	for _, field := range required {
		if strings.TrimSpace(field.val) == "" {
			return fmt.Errorf("%s: missing %s", c.ID, field.name)
		}
	}
	if len(c.Triggers) == 0 || len(c.Exclusions) == 0 || len(c.Dependencies) == 0 || len(c.Compatibility) == 0 {
		return fmt.Errorf("%s: triggers, exclusions, dependencies, and compatibility are required", c.ID)
	}
	if !knownLifecycle(c.Lifecycle) {
		return fmt.Errorf("%s: unknown lifecycle %q", c.ID, c.Lifecycle)
	}
	switch c.Cost {
	case "free", "paid", "freemium", "unspecified":
	default:
		return fmt.Errorf("%s: unknown cost %q", c.ID, c.Cost)
	}
	switch c.Verification {
	case "unverified", "schema_only", "verified":
	default:
		return fmt.Errorf("%s: unknown verification %q", c.ID, c.Verification)
	}
	if c.Lifecycle == LifeSuccessful || c.Lifecycle == LifeVerified || c.Lifecycle == LifeExecuted {
		if !hasRunEvidence(c.UsageEvidence) {
			return fmt.Errorf("%s: lifecycle %s requires run evidence", c.ID, c.Lifecycle)
		}
	}
	return nil
}

// Promote moves one step. A catalog row cannot jump to successful.
func Promote(from, to string) (string, error) {
	if !knownLifecycle(from) || !knownLifecycle(to) {
		return from, fmt.Errorf("unknown lifecycle")
	}
	fi, ti := lifecycleIndex(from), lifecycleIndex(to)
	if ti == fi+1 {
		return to, nil
	}
	return from, fmt.Errorf("cannot move lifecycle from %s to %s", from, to)
}

func knownLifecycle(stage string) bool {
	return lifecycleIndex(stage) >= 0
}

func lifecycleIndex(stage string) int {
	for i, name := range LifecycleOrder {
		if name == stage {
			return i
		}
	}
	return -1
}

func hasRunEvidence(evidence string) bool {
	e := strings.ToLower(evidence)
	return strings.Contains(e, "run:") || strings.Contains(e, "evidence:")
}

func firstText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func copyNonEmpty(in []string) []string {
	out := make([]string, 0, len(in))
	for _, item := range in {
		if strings.TrimSpace(item) != "" {
			out = append(out, strings.TrimSpace(item))
		}
	}
	return out
}

func phaseFromRuntime(method string) string {
	switch strings.TrimSpace(method) {
	case "on_demand_research", "reference_only":
		return "research"
	case "project_scoped_install":
		return "implementation"
	case "mcp_connection":
		return "tool"
	case "global_active_skill", "curated_catalog_lazy_load":
		return "route"
	case "on_demand_cli":
		return "tool"
	default:
		return "unspecified"
	}
}

func costOf(r *Resource) string {
	if r.CostClass != "" {
		return r.CostClass
	}
	status := strings.ToUpper(r.Status)
	auth := strings.ToUpper(r.Auth)
	if status == "PURCHASE_PENDING" || auth == "PURCHASE_PENDING" {
		return "paid"
	}
	return "unspecified"
}

func maintenanceOf(r *Resource) string {
	if r.MaintenanceStatus != "" {
		return r.MaintenanceStatus
	}
	switch strings.ToUpper(r.Status) {
	case "OBSOLETE":
		return "obsolete"
	case "ARCHIVED":
		return "unmaintained"
	default:
		return "unspecified"
	}
}

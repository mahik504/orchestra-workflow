package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Resource represents a canonical machine-readable resource or capability
// as defined in registries/resources.json.
type Resource struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	CanonicalURL       string   `json:"canonical_url"`
	SourceType         string   `json:"source_type"`                 // "github_repository", "web_reference", "npm_package", "cli_binary"
	SourceRepository   string   `json:"source_repository,omitempty"` // Git repository URL if applicable
	DocumentationURL   string   `json:"documentation_url,omitempty"`
	DocURL             string   `json:"doc_url,omitempty"`
	License            string   `json:"license,omitempty"`
	Category           []string `json:"category"`
	Representation     string   `json:"representation"`     // "dependency", "research_source", "skill", "cli", "mcp"
	RoutingTags        []string `json:"routing_tags"`       // e.g. ["motion", "microinteraction", "timeline"]
	AcquisitionMethod  string   `json:"acquisition_method"` // "npm", "git", "web_fetch", "pip"
	RuntimeMethod      string   `json:"runtime_method"`     // "project_scoped_install", "on_demand_research", "global_active_skill", "on_demand_cli"
	Status             string   `json:"status"`             // "ACTIVE", "CURATED_OPTIONAL", "REJECTED", "BOOKMARK", "DEFERRED", "CORE", "REFERENCE", "ARCHIVED", "LEFTOVER"
	TriggerConditions  []string `json:"trigger_conditions,omitempty"`
	AvoidConditions    []string `json:"avoid_conditions,omitempty"`
	PolicyVerdict      string   `json:"policy_verdict,omitempty"`
	Rationale          string   `json:"rationale,omitempty"`
	TokenContextWeight float64  `json:"token_context_weight,omitempty"`
	TokenWeight        float64  `json:"token_weight,omitempty"`
}

// ToCapability converts a Resource to the backwards-compatible Capability model.
func (r *Resource) ToCapability() *Capability {
	cat := CategoryCore
	if len(r.Category) > 0 {
		switch strings.ToUpper(r.Category[0]) {
		case "MOTION", "WEBGL", "3D", "VISUAL_RESEARCH", "DESIGN_DIRECTION", "TYPOGRAPHY", "FRONTEND":
			cat = CategorySpecialist
		case "VERIFICATION", "VISUAL_QA":
			cat = CategoryVerificationAdapter
		case "RESEARCH", "REFERENCE", "INSPIRATION":
			cat = CategoryReference
		case "EXPERIMENTAL":
			cat = CategoryExperimental
		case "REJECTED":
			cat = CategoryRejected
		default:
			cat = CategorySpecialist
		}
	}

	weight := r.TokenContextWeight
	if weight <= 0 {
		weight = r.TokenWeight
	}
	if weight <= 0 {
		weight = 1000
	}

	docURL := r.DocumentationURL
	if docURL == "" {
		docURL = r.DocURL
	}

	return &Capability{
		ID:                    r.ID,
		Name:                  r.Name,
		Repository:            r.SourceRepository,
		Category:              cat,
		CapabilityDesc:        strings.Join(r.RoutingTags, ", "),
		SupportedEnvironments: []string{"web", "node"},
		InstallationMethod:    r.AcquisitionMethod,
		RuntimeDependency:     r.Representation == "dependency",
		TokenContextWeight:    weight,
		ActivationConditions:  r.TriggerConditions,
		Alternatives:          r.AvoidConditions,
		Provenance:            r.CanonicalURL,
		Rationale:             r.Rationale,
		Status:                r.Status,
		IsLoaded:              true,
	}
}

// ResourceCatalog manages the in-memory indexed collection of resources loaded from resources.json.
type ResourceCatalog struct {
	resourcesByID       map[string]*Resource
	resourcesByTag      map[string][]*Resource
	resourcesByCategory map[string][]*Resource
	resourcesByAcq      map[string][]*Resource
	resourcesByRep      map[string][]*Resource
	aliases             map[string]*Resource
	allResources        []*Resource
	sourcePath          string
}

// LoadResourceCatalog loads and indexes resources from a JSON array file (e.g. registries/resources.json).
func LoadResourceCatalog(path string) (*ResourceCatalog, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("resource catalog path cannot be empty")
	}

	if err := CheckQuarantineBoundary(path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource catalog from %s: %w", path, err)
	}

	// Strip UTF-8 Byte Order Mark (BOM) if present (common on Windows)
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return nil, fmt.Errorf("resource catalog file %s is empty", path)
	}

	if !strings.HasPrefix(trimmed, "[") {
		return nil, fmt.Errorf("invalid resource catalog in %s: expected JSON array root", path)
	}

	var rawList []*Resource
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse resource catalog JSON from %s: %w", path, err)
	}

	catalog := &ResourceCatalog{
		resourcesByID:       make(map[string]*Resource),
		resourcesByTag:      make(map[string][]*Resource),
		resourcesByCategory: make(map[string][]*Resource),
		resourcesByAcq:      make(map[string][]*Resource),
		resourcesByRep:      make(map[string][]*Resource),
		aliases:             make(map[string]*Resource),
		allResources:        make([]*Resource, 0, len(rawList)),
		sourcePath:          path,
	}

	for i, res := range rawList {
		if res == nil {
			continue
		}
		res.ID = strings.TrimSpace(res.ID)
		if res.ID == "" {
			return nil, fmt.Errorf("invalid resource at index %d in %s: resource id is required", i, path)
		}

		// Ensure documentation fields sync
		if res.DocumentationURL == "" && res.DocURL != "" {
			res.DocumentationURL = res.DocURL
		} else if res.DocURL == "" && res.DocumentationURL != "" {
			res.DocURL = res.DocumentationURL
		}

		// Quarantine checks on fields
		if err := CheckQuarantineBoundary(res.CanonicalURL); err != nil {
			return nil, fmt.Errorf("resource %s has quarantined canonical_url: %w", res.ID, err)
		}
		if err := CheckQuarantineBoundary(res.SourceRepository); err != nil {
			return nil, fmt.Errorf("resource %s has quarantined source_repository: %w", res.ID, err)
		}
		if err := CheckQuarantineBoundary(res.DocumentationURL); err != nil {
			return nil, fmt.Errorf("resource %s has quarantined documentation_url: %w", res.ID, err)
		}

		lowerID := strings.ToLower(res.ID)
		if _, exists := catalog.resourcesByID[lowerID]; exists {
			return nil, fmt.Errorf("duplicate resource id %q found in %s", res.ID, path)
		}

		catalog.resourcesByID[lowerID] = res
		catalog.allResources = append(catalog.allResources, res)

		// Register shorthand aliases
		if strings.HasSuffix(lowerID, "-design") {
			catalog.aliases[strings.TrimSuffix(lowerID, "-design")] = res
		}
		if strings.HasSuffix(lowerID, "-studio") {
			catalog.aliases[strings.TrimSuffix(lowerID, "-studio")] = res
		}
		if strings.HasSuffix(lowerID, "-institute") {
			catalog.aliases[strings.TrimSuffix(lowerID, "-institute")] = res
		}
		if strings.HasSuffix(lowerID, "-effect") {
			catalog.aliases[strings.TrimSuffix(lowerID, "-effect")] = res
		}
		if lowerID == "semgrep" {
			catalog.aliases["semgrep-adapter"] = res
		}
		if lowerID == "awesome-design-md" {
			catalog.aliases["designmd"] = res
		}
		if lowerID == "react-bits" {
			catalog.aliases["reactbits"] = res
		}
		if lowerID == "apple-playing-haptics" {
			catalog.aliases["apple-haptics"] = res
		}
		if lowerID == "android-vibration-effect" {
			catalog.aliases["android-vibration"] = res
		}
		if lowerID == "material-3-typography" {
			catalog.aliases["material3-typography"] = res
		}
		if lowerID == "not-your-type" {
			catalog.aliases["notyourtype"] = res
		}

		// Index by routing tags
		for _, tag := range res.RoutingTags {
			normTag := strings.ToLower(strings.TrimSpace(tag))
			if normTag != "" {
				catalog.resourcesByTag[normTag] = append(catalog.resourcesByTag[normTag], res)
			}
		}

		// Index by categories
		for _, cat := range res.Category {
			normCat := strings.ToLower(strings.TrimSpace(cat))
			if normCat != "" {
				catalog.resourcesByCategory[normCat] = append(catalog.resourcesByCategory[normCat], res)
			}
		}

		// Index by acquisition method
		if acq := strings.ToLower(strings.TrimSpace(res.AcquisitionMethod)); acq != "" {
			catalog.resourcesByAcq[acq] = append(catalog.resourcesByAcq[acq], res)
		}

		// Index by representation
		if rep := strings.ToLower(strings.TrimSpace(res.Representation)); rep != "" {
			catalog.resourcesByRep[rep] = append(catalog.resourcesByRep[rep], res)
		}
	}

	return catalog, nil
}

// FindByID returns a resource by its ID or known alias (case-insensitive lookup).
func (c *ResourceCatalog) FindByID(id string) (*Resource, bool) {
	if c == nil {
		return nil, false
	}
	normID := strings.ToLower(strings.TrimSpace(id))
	if res, ok := c.resourcesByID[normID]; ok {
		return res, true
	}
	if res, ok := c.aliases[normID]; ok {
		return res, true
	}
	return nil, false
}

// FindByTag returns all resources tagged with the given tag (case-insensitive).
func (c *ResourceCatalog) FindByTag(tag string) []*Resource {
	if c == nil {
		return nil
	}
	normTag := strings.ToLower(strings.TrimSpace(tag))
	items := c.resourcesByTag[normTag]
	result := make([]*Resource, len(items))
	copy(result, items)
	return result
}

// FindByTags returns a deduplicated union of resources matching any of the given tags.
func (c *ResourceCatalog) FindByTags(tags []string) []*Resource {
	if c == nil || len(tags) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var result []*Resource

	for _, tag := range tags {
		normTag := strings.ToLower(strings.TrimSpace(tag))
		for _, res := range c.resourcesByTag[normTag] {
			if !seen[res.ID] {
				seen[res.ID] = true
				result = append(result, res)
			}
		}
	}

	return result
}

// FindByCategory returns all resources in the specified category.
func (c *ResourceCatalog) FindByCategory(category string) []*Resource {
	if c == nil {
		return nil
	}
	normCat := strings.ToLower(strings.TrimSpace(category))
	items := c.resourcesByCategory[normCat]
	result := make([]*Resource, len(items))
	copy(result, items)
	return result
}

// FindByAcquisitionMethod returns resources with the given acquisition method (e.g. "npm", "git").
func (c *ResourceCatalog) FindByAcquisitionMethod(method string) []*Resource {
	if c == nil {
		return nil
	}
	normMethod := strings.ToLower(strings.TrimSpace(method))
	items := c.resourcesByAcq[normMethod]
	result := make([]*Resource, len(items))
	copy(result, items)
	return result
}

// FindByRepresentation returns resources with the given representation (e.g. "dependency", "skill", "cli").
func (c *ResourceCatalog) FindByRepresentation(rep string) []*Resource {
	if c == nil {
		return nil
	}
	normRep := strings.ToLower(strings.TrimSpace(rep))
	items := c.resourcesByRep[normRep]
	result := make([]*Resource, len(items))
	copy(result, items)
	return result
}

// All returns all resources in the catalog in loaded order.
func (c *ResourceCatalog) All() []*Resource {
	if c == nil {
		return nil
	}
	result := make([]*Resource, len(c.allResources))
	copy(result, c.allResources)
	return result
}

// Count returns the total number of resources in the catalog.
func (c *ResourceCatalog) Count() int {
	if c == nil {
		return 0
	}
	return len(c.allResources)
}

// Tags returns a sorted list of all unique indexed tags.
func (c *ResourceCatalog) Tags() []string {
	if c == nil {
		return nil
	}
	tags := make([]string, 0, len(c.resourcesByTag))
	for tag := range c.resourcesByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// ToRegistry converts the entire catalog into a backwards-compatible Registry.
func (c *ResourceCatalog) ToRegistry() *Registry {
	reg := NewRegistry()
	if c == nil {
		return reg
	}
	for _, res := range c.allResources {
		reg.Capabilities[res.ID] = res.ToCapability()
	}
	return reg
}

// CapabilityWorkflow defines the multi-phase workflow associated with a capability in the design graph.
type CapabilityWorkflow struct {
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	PrimaryArchetype   string   `json:"primary_archetype,omitempty"`
	TriggerTags        []string `json:"trigger_tags,omitempty"`
	QualityBar         string   `json:"quality_bar,omitempty"`
	RiskRank           int      `json:"risk_rank,omitempty"`
	Platform           string   `json:"platform,omitempty"`
	TriggerConditions  []string `json:"trigger_conditions,omitempty"`
	SkipConditions     []string `json:"skip_conditions,omitempty"`
	Discovery          []string `json:"discovery"`
	ReverseEngineering []string `json:"reverse_engineering,omitempty"`
	Synthesis          []string `json:"synthesis"`
	Implementation     []string `json:"implementation"`
	OptionalExtensions []string `json:"optional_extensions,omitempty"`
	QA                 []string `json:"qa"`
	AntiPatterns       []string `json:"anti_patterns,omitempty"`
}

// CapabilityPhaseDefinition is an alias for CapabilityWorkflow for full compatibility.
type CapabilityPhaseDefinition = CapabilityWorkflow

// DesignResourceGraph models the capability-to-resource graph loaded from design-resource-graph.json.
type DesignResourceGraph struct {
	Schema       string                                `json:"$schema,omitempty"`
	Version      string                                `json:"version,omitempty"`
	Description  string                                `json:"description,omitempty"`
	Domains      map[string][]string                   `json:"domains"`
	Capabilities map[string]*CapabilityPhaseDefinition `json:"capabilities"`
	sourcePath   string
}

// CapabilityRoute represents the fully resolved execution route for a capability.
type CapabilityRoute struct {
	CapabilityName     string   `json:"capability_name"`
	CapabilityID       string   `json:"capability_id,omitempty"`
	Name               string   `json:"name,omitempty"`
	PrimaryArchetype   string   `json:"primary_archetype,omitempty"`
	QualityBar         string   `json:"quality_bar,omitempty"`
	RiskRank           int      `json:"risk_rank,omitempty"`
	Platform           string   `json:"platform,omitempty"`
	TriggerConditions  []string `json:"trigger_conditions,omitempty"`
	SkipConditions     []string `json:"skip_conditions,omitempty"`
	MatchedTags        []string `json:"matched_tags,omitempty"`
	DiscoveryDomains   []string `json:"discovery_domains"`
	DiscoveryResources []string `json:"discovery_resources"`
	ReverseEngineering []string `json:"reverse_engineering,omitempty"`
	SynthesisDomains   []string `json:"synthesis_domains"`
	SynthesisResources []string `json:"synthesis_resources"`
	Implementation     []string `json:"implementation_resources"`
	OptionalExtensions []string `json:"optional_extensions,omitempty"`
	QA                 []string `json:"qa"`
	AntiPatterns       []string `json:"anti_patterns,omitempty"`
	AllResourceIDs     []string `json:"all_resource_ids"`
}

// LoadDesignGraph loads and parses design-resource-graph.json.
func LoadDesignGraph(path string) (*DesignResourceGraph, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("design graph path cannot be empty")
	}

	if err := CheckQuarantineBoundary(path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read design graph from %s: %w", path, err)
	}

	// Strip UTF-8 Byte Order Mark (BOM) if present (common on Windows)
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var graph DesignResourceGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("failed to parse design graph JSON from %s: %w", path, err)
	}

	if graph.Domains == nil {
		graph.Domains = make(map[string][]string)
	}
	if graph.Capabilities == nil {
		graph.Capabilities = make(map[string]*CapabilityPhaseDefinition)
	}
	graph.sourcePath = path

	return &graph, nil
}

// GetDomainResources returns the list of resource IDs belonging to the given domain.
func (g *DesignResourceGraph) GetDomainResources(domain string) []string {
	if g == nil || g.Domains == nil {
		return nil
	}
	items, ok := g.Domains[domain]
	if !ok {
		return nil
	}
	result := make([]string, len(items))
	copy(result, items)
	return result
}

// expandEntries expands domain names into their constituent resource IDs, or preserves raw resource IDs.
func (g *DesignResourceGraph) expandEntries(entries []string) (domains []string, resources []string) {
	seenRes := make(map[string]bool)
	seenDom := make(map[string]bool)

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if domRes, isDomain := g.Domains[entry]; isDomain {
			if !seenDom[entry] {
				seenDom[entry] = true
				domains = append(domains, entry)
			}
			for _, r := range domRes {
				if !seenRes[r] {
					seenRes[r] = true
					resources = append(resources, r)
				}
			}
		} else {
			if !seenRes[entry] {
				seenRes[entry] = true
				resources = append(resources, entry)
			}
		}
	}
	return domains, resources
}

// ResolveCapabilityRoute resolves a single capability by name into a full CapabilityRoute.
func (g *DesignResourceGraph) ResolveCapabilityRoute(name string, matchedTags []string) (*CapabilityRoute, bool) {
	if g == nil || g.Capabilities == nil {
		return nil, false
	}

	workflow, ok := g.Capabilities[name]
	if !ok {
		return nil, false
	}

	discDom, discRes := g.expandEntries(workflow.Discovery)
	synDom, synRes := g.expandEntries(workflow.Synthesis)
	_, implRes := g.expandEntries(workflow.Implementation)
	_, optRes := g.expandEntries(workflow.OptionalExtensions)
	_, revRes := g.expandEntries(workflow.ReverseEngineering)
	_, qaRes := g.expandEntries(workflow.QA)

	allSeen := make(map[string]bool)
	var allRes []string
	combineLists := func(list []string) {
		for _, r := range list {
			if !allSeen[r] {
				allSeen[r] = true
				allRes = append(allRes, r)
			}
		}
	}

	combineLists(discRes)
	combineLists(revRes)
	combineLists(synRes)
	combineLists(implRes)
	combineLists(optRes)
	combineLists(qaRes)

	if qaRes == nil {
		qaRes = []string{}
	}

	route := &CapabilityRoute{
		CapabilityName:     name,
		CapabilityID:       name,
		Name:               workflow.Name,
		PrimaryArchetype:   workflow.PrimaryArchetype,
		QualityBar:         workflow.QualityBar,
		RiskRank:           workflow.RiskRank,
		Platform:           workflow.Platform,
		TriggerConditions:  workflow.TriggerConditions,
		SkipConditions:     workflow.SkipConditions,
		MatchedTags:        matchedTags,
		DiscoveryDomains:   discDom,
		DiscoveryResources: discRes,
		ReverseEngineering: revRes,
		SynthesisDomains:   synDom,
		SynthesisResources: synRes,
		Implementation:     implRes,
		OptionalExtensions: optRes,
		QA:                 qaRes,
		AntiPatterns:       workflow.AntiPatterns,
		AllResourceIDs:     allRes,
	}

	return route, true
}

// ResolveCapabilities matches input tags against graph capabilities and returns resolved routes.
func (g *DesignResourceGraph) ResolveCapabilities(tags []string) []CapabilityRoute {
	if g == nil || g.Capabilities == nil || len(tags) == 0 {
		return []CapabilityRoute{}
	}

	matchedCapabilities := make(map[string][]string) // capabilityName -> matchedTags

	for _, rawTag := range tags {
		normTag := strings.ToLower(strings.TrimSpace(rawTag))
		if normTag == "" {
			continue
		}

		// 1. Direct match on capability name
		for capName, def := range g.Capabilities {
			lowerCap := strings.ToLower(capName)
			if lowerCap == normTag || strings.ReplaceAll(lowerCap, "-", " ") == normTag || strings.ReplaceAll(lowerCap, " ", "-") == normTag {
				matchedCapabilities[capName] = append(matchedCapabilities[capName], rawTag)
				continue
			}

			// Check trigger tags
			for _, trig := range def.TriggerTags {
				if strings.EqualFold(trig, normTag) {
					matchedCapabilities[capName] = append(matchedCapabilities[capName], rawTag)
					break
				}
			}
		}

		// 2. Keyword/substring match against capability name
		for capName := range g.Capabilities {
			parts := strings.Split(strings.ToLower(capName), "-")
			for _, part := range parts {
				if part == normTag && len(part) >= 2 {
					matchedCapabilities[capName] = append(matchedCapabilities[capName], rawTag)
				}
			}
		}

		// 3. Domain match: if tag matches a domain name, activate capabilities referencing that domain across any phase
		if _, isDomain := g.Domains[normTag]; isDomain {
			for capName, wf := range g.Capabilities {
				totalLen := len(wf.Discovery) + len(wf.Synthesis) + len(wf.Implementation) + len(wf.OptionalExtensions) + len(wf.ReverseEngineering) + len(wf.QA)
				allPhases := make([]string, 0, totalLen)
				allPhases = append(allPhases, wf.Discovery...)
				allPhases = append(allPhases, wf.Synthesis...)
				allPhases = append(allPhases, wf.Implementation...)
				allPhases = append(allPhases, wf.OptionalExtensions...)
				allPhases = append(allPhases, wf.ReverseEngineering...)
				allPhases = append(allPhases, wf.QA...)
				for _, d := range allPhases {
					if strings.ToLower(d) == normTag {
						matchedCapabilities[capName] = append(matchedCapabilities[capName], rawTag)
						break
					}
				}
			}
		}
	}

	var routes []CapabilityRoute
	var capNames []string
	for name := range matchedCapabilities {
		capNames = append(capNames, name)
	}
	sort.Strings(capNames)

	for _, name := range capNames {
		route, ok := g.ResolveCapabilityRoute(name, matchedCapabilities[name])
		if ok {
			routes = append(routes, *route)
		}
	}

	return routes
}

// AllCapabilities returns a sorted list of all capability names defined in the graph.
func (g *DesignResourceGraph) AllCapabilities() []string {
	if g == nil || g.Capabilities == nil {
		return nil
	}
	names := make([]string, 0, len(g.Capabilities))
	for name := range g.Capabilities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// AllDomains returns a sorted list of all domain names defined in the graph.
func (g *DesignResourceGraph) AllDomains() []string {
	if g == nil || g.Domains == nil {
		return nil
	}
	names := make([]string, 0, len(g.Domains))
	for name := range g.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type CapabilityCategory string

const (
	CategoryCore                CapabilityCategory = "CORE"
	CategorySpecialist          CapabilityCategory = "SPECIALIST"
	CategoryVerificationAdapter CapabilityCategory = "VERIFICATION ADAPTER"
	CategoryExecutionAdapter    CapabilityCategory = "EXECUTION ADAPTER"
	CategoryReference           CapabilityCategory = "REFERENCE"
	CategoryExperimental        CapabilityCategory = "EXPERIMENTAL"
	CategoryRejected            CapabilityCategory = "REJECTED"
)

type Capability struct {
	ID                    string             `json:"id"`
	Name                  string             `json:"name"`
	Repository            string             `json:"repository,omitempty"`
	Category              CapabilityCategory `json:"category"`
	CapabilityDesc        string             `json:"capability"`
	License               string             `json:"license,omitempty"`
	Maturity              string             `json:"maturity,omitempty"`
	MaintenanceSignal     string             `json:"maintenance_signal,omitempty"`
	SupportedEnvironments []string           `json:"supported_environments,omitempty"`
	InstallationMethod    string             `json:"installation_method,omitempty"`
	RuntimeDependency     bool               `json:"runtime_dependency,omitempty"`
	TokenContextWeight    float64            `json:"token_context_weight,omitempty"`
	ActivationConditions  []string           `json:"activation_conditions,omitempty"`
	Conflicts             []string           `json:"conflicts,omitempty"`
	Alternatives          []string           `json:"alternatives,omitempty"`
	IntegrationMode       string             `json:"integration_mode,omitempty"`
	VerificationMethod    string             `json:"verification_method,omitempty"`
	Provenance            string             `json:"provenance,omitempty"`
	Rationale             string             `json:"rationale,omitempty"`
	Status                string             `json:"status"`
	
	// Lazy Loading
	LazyLoadPath          string             `json:"-"`
	IsLoaded              bool               `json:"-"`
	RawContent            []byte             `json:"-"`
}

func (c *Capability) LoadDetails() error {
	if c.IsLoaded {
		return nil
	}
	if c.LazyLoadPath == "" {
		return fmt.Errorf("no path specified for lazy loading capability %s", c.ID)
	}
	
	bytes, err := os.ReadFile(c.LazyLoadPath)
	if err != nil {
		return err
	}
	
	c.RawContent = bytes
	c.IsLoaded = true
	return nil
}

type Registry struct {
	Capabilities map[string]*Capability
}

func NewRegistry() *Registry {
	return &Registry{
		Capabilities: make(map[string]*Capability),
	}
}

// LoadFromJSON loads a capability or resource catalog from a JSON file and adds it to the registry.
// If the target file is a JSON array (like resources.json), it automatically imports all catalog items.
func (r *Registry) LoadFromJSON(filepath string) error {
	if err := CheckQuarantineBoundary(filepath); err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	// Strip UTF-8 Byte Order Mark (BOM) if present (common on Windows)
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "[") {
		return r.LoadResourceCatalog(filepath)
	}

	var cap Capability
	if err := json.Unmarshal(data, &cap); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filepath, err)
	}

	if cap.ID == "" {
		return fmt.Errorf("capability ID is required in %s", filepath)
	}

	cap.LazyLoadPath = filepath
	r.Capabilities[cap.ID] = &cap
	return nil
}

// LoadResourceCatalog loads resources from a JSON array catalog (e.g. registries/resources.json)
// and imports all items into r.Capabilities, preserving backwards compatibility with existing consumers.
func (r *Registry) LoadResourceCatalog(filepath string) error {
	cat, err := LoadResourceCatalog(filepath)
	if err != nil {
		return err
	}
	r.ImportCatalog(cat)
	return nil
}

// ImportCatalog populates the registry with capabilities converted from the given ResourceCatalog.
func (r *Registry) ImportCatalog(cat *ResourceCatalog) {
	if r == nil || cat == nil {
		return
	}
	for _, res := range cat.All() {
		r.Capabilities[res.ID] = res.ToCapability()
	}
}

// ImportResource converts a Resource and adds it into r.Capabilities.
func (r *Registry) ImportResource(res *Resource) {
	if r == nil || res == nil {
		return
	}
	r.Capabilities[res.ID] = res.ToCapability()
}


package resources

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// LoadFromJSON loads a capability from a JSON file and adds it to the registry
func (r *Registry) LoadFromJSON(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	var cap Capability
	if err := json.Unmarshal(bytes, &cap); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filepath, err)
	}

	if cap.ID == "" {
		return fmt.Errorf("capability ID is required in %s", filepath)
	}

	cap.LazyLoadPath = filepath
	r.Capabilities[cap.ID] = &cap
	return nil
}

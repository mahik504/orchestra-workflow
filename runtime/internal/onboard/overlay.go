package onboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/memory"
	"github.com/user/orchestra-v3/internal/resources"
)

// OverlayDocument is the Brain-side catalog of user-added resources.
// It is not the public registries/resources.json.
type OverlayDocument struct {
	SchemaVersion string         `json:"schema_version"`
	Updated       string         `json:"updated"`
	Mechanism     string         `json:"mechanism"`
	Resources     []OverlayEntry `json:"resources"`
}

// OverlayEntry is one user-added resource plus inspect/policy metadata.
type OverlayEntry struct {
	Origin              string             `json:"origin"`
	Intent              string             `json:"intent,omitempty"`
	Kind                string             `json:"kind"`
	KindReason          string             `json:"kind_reason"`
	InstallScope        string             `json:"install_scope"`
	RoutingPolicy       string             `json:"routing_policy"`
	FutureRoutingEffect string             `json:"future_routing_effect"`
	Inspection          Inspection         `json:"inspection"`
	Resource            resources.Resource `json:"resource"`
	Learning            map[string]string  `json:"learning,omitempty"`
}

const overlayMechanism = "experience → structured evaluation → memory → routing update → future selection"

// ResolveOverlayPath picks added-resources.json the same way memory picks
// resource-memory.json: exact env, then ORCHESTRA_HOME, then workspace .orchestra/.
func ResolveOverlayPath(workspaceRoot string) string {
	if envPath := os.Getenv("ORCHESTRA_OVERLAY_PATH"); envPath != "" {
		return filepath.Clean(envPath)
	}
	if home := os.Getenv("ORCHESTRA_HOME"); home != "" {
		return filepath.Clean(filepath.Join(home, "memory", "added-resources.json"))
	}
	if strings.TrimSpace(workspaceRoot) == "" {
		workspaceRoot = "."
	}
	return filepath.Clean(filepath.Join(workspaceRoot, ".orchestra", "memory", "added-resources.json"))
}

// LoadOverlay reads an overlay document, or returns empty if the file is missing.
func LoadOverlay(path string) (*OverlayDocument, error) {
	if strings.TrimSpace(path) == "" {
		path = ResolveOverlayPath(".")
	}
	if err := resources.CheckQuarantineBoundary(path); err != nil {
		return nil, err
	}
	doc := &OverlayDocument{
		SchemaVersion: "1.0.0",
		Mechanism:     overlayMechanism,
		Resources:     []OverlayEntry{},
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return doc, nil
		}
		return nil, err
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(data, doc); err != nil {
		return nil, fmt.Errorf("overlay JSON: %w", err)
	}
	if doc.Resources == nil {
		doc.Resources = []OverlayEntry{}
	}
	if doc.Mechanism == "" {
		doc.Mechanism = overlayMechanism
	}
	return doc, nil
}

// SaveOverlay writes the overlay atomically.
func SaveOverlay(path string, doc *OverlayDocument) error {
	if err := resources.CheckQuarantineBoundary(path); err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("overlay document is nil")
	}
	doc.SchemaVersion = "1.0.0"
	doc.Mechanism = overlayMechanism
	doc.Updated = time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s.tmp.%d", path, time.Now().UnixNano())
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(path)
		if err2 := os.Rename(tmp, path); err2 != nil {
			_ = os.WriteFile(path, data, 0644)
			_ = os.Remove(tmp)
		}
	}
	return nil
}

// UpsertEntry replaces or appends by resource id.
func (d *OverlayDocument) UpsertEntry(entry OverlayEntry) {
	if d == nil {
		return
	}
	id := strings.ToLower(strings.TrimSpace(entry.Resource.ID))
	for i := range d.Resources {
		if strings.ToLower(d.Resources[i].Resource.ID) == id {
			d.Resources[i] = entry
			return
		}
	}
	d.Resources = append(d.Resources, entry)
}

func (d *OverlayDocument) Find(id string) *OverlayEntry {
	if d == nil {
		return nil
	}
	want := strings.ToLower(strings.TrimSpace(id))
	for i := range d.Resources {
		if strings.ToLower(d.Resources[i].Resource.ID) == want {
			return &d.Resources[i]
		}
	}
	return nil
}

// MergeIntoCatalog copies overlay rows into the live catalog. Public
// registries/resources.json is not written.
func MergeIntoCatalog(catalog *resources.ResourceCatalog, doc *OverlayDocument) error {
	if catalog == nil || doc == nil {
		return nil
	}
	for i := range doc.Resources {
		res := doc.Resources[i].Resource
		if err := catalog.Upsert(&res); err != nil {
			return err
		}
	}
	return nil
}

// ApplyMemoryPolicy updates overlay routing from recorded evaluations.
// This is not a training loop. Last failure suppresses auto-activation;
// last success restores ACTIVE. Public catalog rows are not touched.
func ApplyMemoryPolicy(doc *OverlayDocument, store *memory.ResourceMemoryStore) {
	if doc == nil {
		return
	}
	for i := range doc.Resources {
		entry := &doc.Resources[i]
		id := entry.Resource.ID
		var agg *memory.ResourceAggregate
		if store != nil {
			if found, ok := store.GetAggregate(id); ok {
				agg = found
			}
		}
		switch {
		case agg != nil && agg.LastOutcome == memory.OutcomeFailure:
			entry.RoutingPolicy = "suppressed"
			if agg.SuccessCount > 0 {
				entry.Resource.Status = "CURATED_OPTIONAL"
			} else {
				entry.Resource.Status = "REJECTED"
			}
			entry.Resource.PolicyVerdict = PolicySuppressed
			entry.FutureRoutingEffect = "auto-activation suppressed after recorded failure; re-add or record a later success to restore"
		case agg != nil && agg.LastOutcome == memory.OutcomeSuccess:
			entry.RoutingPolicy = "active"
			entry.Resource.Status = "ACTIVE"
			entry.Resource.PolicyVerdict = PolicyOverlay
			entry.FutureRoutingEffect = "keep auto-activating when trigger conditions match"
		default:
			if entry.RoutingPolicy == "" {
				entry.RoutingPolicy = "active"
			}
			if entry.FutureRoutingEffect == "" {
				entry.FutureRoutingEffect = "activate when trigger conditions match; skip otherwise"
			}
			entry.Resource.PolicyVerdict = PolicyOverlay
		}
	}
}

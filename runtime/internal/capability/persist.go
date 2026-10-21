package capability

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

// ErrDecisionMutated means a later stage changed the selected set, strength,
// fallback chain, or availability after Select. The file written at Select stays.
var ErrDecisionMutated = errors.New("acceptance criteria violated: capability decision mutated after select")

func WriteDecision(workspace, missionID string, d Decision) (string, error) {
	dir := mission.RunDir(workspace, missionID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "capability.json")
	if prev, err := os.ReadFile(path); err == nil {
		var old Decision
		if json.Unmarshal(prev, &old) == nil {
			if !AuthorityEqual(old, d) {
				return path, ErrDecisionMutated
			}
			return path, nil
		}
	}
	raw, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return path, os.WriteFile(path, raw, 0644)
}

func AuthorityEqual(a, b Decision) bool {
	return authorityKey(a) == authorityKey(b)
}

func authorityKey(d Decision) string {
	var b strings.Builder
	b.WriteString(d.MissionID)
	b.WriteByte('|')
	b.WriteString(d.HostLabel)
	b.WriteByte('|')
	writeItems(&b, d.Selected)
	b.WriteByte('|')
	writeItems(&b, d.Rejected)
	b.WriteByte('|')
	writeItems(&b, d.Unavailable)
	b.WriteByte('|')
	var hops []string
	for _, f := range d.Fallbacks {
		hops = append(hops, f.Requested+"->"+f.Selected)
	}
	sort.Strings(hops)
	b.WriteString(strings.Join(hops, ","))
	return b.String()
}

func writeItems(b *strings.Builder, items []Item) {
	var parts []string
	for _, it := range items {
		parts = append(parts, it.ID+"/"+it.Category+"/"+it.Strength+"/"+fmt.Sprintf("%v", it.Available))
	}
	sort.Strings(parts)
	b.WriteString(strings.Join(parts, ","))
}

func ReadDecision(workspace, missionID string) (*Decision, error) {
	path := filepath.Join(mission.RunDir(workspace, missionID), "capability.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no capability file at %s: %w", path, err)
	}
	var d Decision
	if err := json.Unmarshal(raw, &d); err != nil {
		return nil, fmt.Errorf("corrupt capability file: %w", err)
	}
	return &d, nil
}

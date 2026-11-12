package adapters

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

// HostType represents an AI agent host IDE/environment
type HostType string

const (
	HostCursor      HostType = "cursor"
	HostClaudeCode  HostType = "claude"
	HostAntigravity HostType = "antigravity"
	HostAgents      HostType = "agents"
	HostJcode       HostType = "jcode"
	HostCodex       HostType = "codex"
)

func init() {
	if ids := skillIDsFromRegistry(); len(ids) > 0 {
		CanonicalActiveSkills = ids
	}
}

func skillIDsFromRegistry() []string {
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}
	for i := 0; i < 8; i++ {
		raw, err := os.ReadFile(filepath.Join(dir, "registries", "skills.json"))
		if err == nil {
			var doc struct {
				Skills []struct {
					ID string `json:"id"`
				} `json:"skills"`
			}
			if json.Unmarshal(raw, &doc) == nil && len(doc.Skills) > 0 {
				ids := make([]string, 0, len(doc.Skills))
				for _, row := range doc.Skills {
					if row.ID != "" {
						ids = append(ids, row.ID)
					}
				}
				sort.Strings(ids)
				return ids
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil
}

// CanonicalActiveSkills is the skill id list from registries/skills.json.
// The slice below is the fallback when that file cannot be read.
var CanonicalActiveSkills = []string{
	"animate",
	"ci-security-scanning-with-strix",
	"ow-diagrams",
	"eas-app-stores",
	"emil-design-eng",
	"expo-data-fetching",
	"expo-dev-client",
	"expo-native-ui",
	"expo-project-structure",
	"expo-router",
	"expo-upgrade",
	"fix-security-vulnerabilities-with-strix",
	"gsap-core",
	"hallmark",
	"impeccable",
	"ow-conductor",
	"ow-docs",
	"ow-ship",
	"ow-vault",
	"penetration-testing-with-strix",
	"postgres-best-practices",
	"pretty-mermaid",
	"react-best-practices",
	"refero-design",
	"review-animations",
	"ow-semgrep",
	"ow-ship-safe",
	"ow-skill-security-check",
	"stitch-code-to-design",
	"stitch-extract-design-md",
	"stitch-extract-static-html",
	"stitch-generate-design",
	"stitch-manage-design-system",
	"stitch-react-components",
	"stitch-upload-to-stitch",
	"supabase",
	"taste-design",
	"web-design-guidelines",
}

// ParityReport contains the results of cross-host active skill parity check
type ParityReport struct {
	IsParityComplete     bool                  `json:"is_parity_complete"`
	CanonicalSkillCount  int                   `json:"canonical_skill_count"`
	HostSkillCounts      map[HostType]int      `json:"host_skill_counts"`
	MissingSkills        map[HostType][]string `json:"missing_skills"`
	ExtraSkills          map[HostType][]string `json:"extra_skills"`
	QuarantineViolations []string              `json:"quarantine_violations"`
	ByteIdentical        bool                  `json:"byte_identical"`
	MismatchDetails      []string              `json:"mismatch_details"`
	HostsPresent         int                   `json:"hosts_present"`
	AbsentHosts          []HostType            `json:"absent_hosts,omitempty"`
}

// HostSyncEngine manages generation, validation, and synchronization of host configurations
type HostSyncEngine struct {
	WorkspaceRoot string
}

// NewHostSyncEngine initializes a HostSyncEngine rooted at workspaceRoot
func NewHostSyncEngine(workspaceRoot string) *HostSyncEngine {
	if workspaceRoot == "" {
		workspaceRoot = "."
	}
	return &HostSyncEngine{WorkspaceRoot: workspaceRoot}
}

// ResolveHostDirs returns the map of active host directories for a user home.
// The list is registries/host-dirs.json, shared with ow-sync.
func (e *HostSyncEngine) ResolveHostDirs(userHome string) map[HostType]string {
	out := map[HostType]string{}
	for name, parts := range loadHostSegments() {
		segs := append([]string{userHome}, parts...)
		out[HostType(name)] = filepath.Join(segs...)
	}
	return out
}

func loadHostSegments() map[string][]string {
	fallback := map[string][]string{
		"cursor":      {".cursor", "skills"},
		"claude":      {".claude", "skills"},
		"antigravity": {".gemini", "config", "skills"},
		"agents":      {".agents", "skills"},
		"jcode":       {".jcode", "skills"},
		"codex":       {".codex", "skills"},
	}
	dir, err := os.Getwd()
	if err != nil {
		return fallback
	}
	for i := 0; i < 8; i++ {
		raw, err := os.ReadFile(filepath.Join(dir, "registries", "host-dirs.json"))
		if err == nil {
			var doc struct {
				Hosts map[string][]string `json:"hosts"`
			}
			if json.Unmarshal(raw, &doc) == nil && len(doc.Hosts) > 0 {
				return doc.Hosts
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return fallback
}

// VerifyParity audits installed skill directories across hosts and validates 100% byte parity
func (e *HostSyncEngine) VerifyParity(userHome string) (*ParityReport, error) {
	report := &ParityReport{
		IsParityComplete:     true,
		CanonicalSkillCount:  len(CanonicalActiveSkills),
		HostSkillCounts:      make(map[HostType]int),
		MissingSkills:        make(map[HostType][]string),
		ExtraSkills:          make(map[HostType][]string),
		QuarantineViolations: make([]string, 0),
		ByteIdentical:        true,
		MismatchDetails:      make([]string, 0),
	}

	hostDirs := e.ResolveHostDirs(userHome)

	canonicalSet := make(map[string]bool)
	for _, s := range CanonicalActiveSkills {
		canonicalSet[s] = true
	}

	// Track hashes: skillName -> (host -> sha256)
	skillHashes := make(map[string]map[HostType]string)
	activeHosts := make([]HostType, 0)

	for host, dir := range hostDirs {
		if _, err := os.Stat(dir); err != nil {
			report.AbsentHosts = append(report.AbsentHosts, host)
			report.IsParityComplete = false
			continue
		}
		activeHosts = append(activeHosts, host)

		// Check quarantine boundary
		if err := resources.CheckQuarantineBoundary(dir); err != nil {
			report.QuarantineViolations = append(report.QuarantineViolations, fmt.Sprintf("%s points to quarantined path: %v", dir, err))
			report.IsParityComplete = false
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed reading host dir %s: %w", dir, err)
		}

		installed := make(map[string]bool)
		for _, entry := range entries {
			if entry.IsDir() {
				installed[entry.Name()] = true
			}
		}

		report.HostSkillCounts[host] = len(installed)

		// Check for missing canonical skills
		for _, expected := range CanonicalActiveSkills {
			if !installed[expected] {
				report.MissingSkills[host] = append(report.MissingSkills[host], expected)
				report.IsParityComplete = false
			} else {
				// Compute hash of the installed skill
				skillPath := filepath.Join(dir, expected)
				h, err := ComputeSkillDirHash(skillPath)
				if err == nil {
					if skillHashes[expected] == nil {
						skillHashes[expected] = make(map[HostType]string)
					}
					skillHashes[expected][host] = h
				}
			}
		}

		// Check for unapproved skills
		for actual := range installed {
			if !canonicalSet[actual] {
				report.ExtraSkills[host] = append(report.ExtraSkills[host], actual)
				report.IsParityComplete = false
			}
		}
	}

	report.HostsPresent = len(activeHosts)
	if len(activeHosts) == 0 {
		report.IsParityComplete = false
	}

	// Verify byte-identical parity across hosts for each skill
	for _, skill := range CanonicalActiveSkills {
		hashes := skillHashes[skill]
		if hashes == nil || len(hashes) <= 1 {
			continue
		}

		var firstHash string
		var firstHost HostType
		for hst, h := range hashes {
			if firstHash == "" {
				firstHash = h
				firstHost = hst
			} else if h != firstHash {
				report.ByteIdentical = false
				report.IsParityComplete = false
				report.MismatchDetails = append(report.MismatchDetails,
					fmt.Sprintf("skill '%s' hash mismatch between %s (%s) and %s (%s)", skill, firstHost, firstHash[:8], hst, h[:8]))
			}
		}
	}

	return report, nil
}

// AuditActiveSkills performs an exhaustive audit of active host directories and quarantine isolation
func (e *HostSyncEngine) AuditActiveSkills(userHome string) (*ParityReport, error) {
	report, err := e.VerifyParity(userHome)
	if err != nil {
		return nil, err
	}

	// Ensure no active skill path has symlink or junction into quarantined paths
	hostDirs := e.ResolveHostDirs(userHome)
	for host, dir := range hostDirs {
		if _, statErr := os.Stat(dir); statErr != nil {
			continue
		}
		for _, skill := range CanonicalActiveSkills {
			skillPath := filepath.Join(dir, skill)
			if err := resources.CheckQuarantineBoundary(skillPath); err != nil {
				report.QuarantineViolations = append(report.QuarantineViolations,
					fmt.Sprintf("[%s] skill '%s' violates quarantine: %v", host, skill, err))
				report.IsParityComplete = false
			}
		}
	}

	return report, nil
}

// SyncAll delegates to ow-sync. It does not copy over a local edit by itself.
func (e *HostSyncEngine) SyncAll(userHome string, targetHost string) error {
	if err := resources.CheckQuarantineBoundary(userHome); err != nil {
		return fmt.Errorf("userHome violates quarantine: %w", err)
	}
	script := filepath.Join(e.WorkspaceRoot, "runtime", "tools", "ow-sync.mjs")
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("ow-sync missing: %w", err)
	}
	args := []string{script, "--apply"}
	host := strings.ToLower(strings.TrimSpace(targetHost))
	if host == "gemini" {
		host = "antigravity"
	}
	if host != "" && host != "all" {
		args = append(args, "--host", host)
	}
	cmd := exec.Command("node", args...)
	cmd.Dir = e.WorkspaceRoot
	cmd.Env = append(os.Environ(), "OW_SYNC_HOME="+userHome)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ow-sync: %w: %s", err, out)
	}
	return nil
}

// ComputeSkillDirHash recursively hashes all files within a directory in alphabetical order
func ComputeSkillDirHash(dirPath string) (string, error) {
	var fileList []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, relErr := filepath.Rel(dirPath, path)
			if relErr != nil {
				return relErr
			}
			fileList = append(fileList, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(fileList)

	hasher := sha256.New()
	for _, rel := range fileList {
		hasher.Write([]byte(rel))
		fullPath := filepath.Join(dirPath, filepath.FromSlash(rel))
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return "", err
		}
		hasher.Write(data)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

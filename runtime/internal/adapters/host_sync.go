package adapters

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

// HostType represents an AI agent host IDE/environment
type HostType string

const (
	HostCursor      HostType = "cursor"
	HostClaudeCode  HostType = "claude"
	HostAntigravity HostType = "antigravity"
	HostAgents      HostType = "agents"
)

// CanonicalActiveSkills defines the authoritative 30-skill active working set
var CanonicalActiveSkills = []string{
	"animate",
	"ci-security-scanning-with-strix",
	"eas-app-stores",
	"emil-design-eng",
	"expo-data-fetching",
	"expo-dev-client",
	"expo-native-ui",
	"expo-project-structure",
	"expo-router",
	"expo-upgrade",
	"fix-security-vulnerabilities-with-strix",
	"impeccable",
	"orchestra-conductor",
	"orchestra-docs",
	"orchestra-ship",
	"orchestra-vault",
	"penetration-testing-with-strix",
	"review-animations",
	"semgrep-adapter",
	"ship-safe",
	"stitch-code-to-design",
	"stitch-extract-design-md",
	"stitch-extract-static-html",
	"stitch-generate-design",
	"stitch-manage-design-system",
	"stitch-react-components",
	"stitch-upload-to-stitch",
	"superpowers-planning",
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

// ResolveHostDirs returns the map of active host directories for a user home
func (e *HostSyncEngine) ResolveHostDirs(userHome string) map[HostType]string {
	return map[HostType]string{
		HostCursor:      filepath.Join(userHome, ".cursor", "skills"),
		HostClaudeCode:  filepath.Join(userHome, ".claude", "skills"),
		HostAntigravity: filepath.Join(userHome, ".gemini", "config", "skills"),
		HostAgents:      filepath.Join(userHome, ".agents", "skills"),
	}
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
		if _, err := os.Stat(dir); os.IsNotExist(err) {
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

// SyncAll synchronizes canonical skills from the workspace repository into target hosts
func (e *HostSyncEngine) SyncAll(userHome string, targetHost string) error {
	if err := resources.CheckQuarantineBoundary(userHome); err != nil {
		return fmt.Errorf("userHome violates quarantine: %w", err)
	}

	sourceDir := filepath.Join(e.WorkspaceRoot, "skills")
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		// Fallback to Cursor skills or Gemini skills as source
		cand1 := filepath.Join(userHome, ".cursor", "skills")
		cand2 := filepath.Join(userHome, ".gemini", "config", "skills")
		if _, err := os.Stat(cand1); err == nil {
			sourceDir = cand1
		} else if _, err := os.Stat(cand2); err == nil {
			sourceDir = cand2
		} else {
			return fmt.Errorf("no valid skill source directory found (checked %s, %s, %s)", sourceDir, cand1, cand2)
		}
	}

	if err := resources.CheckQuarantineBoundary(sourceDir); err != nil {
		return fmt.Errorf("source skills directory violates quarantine: %w", err)
	}

	hostDirs := e.ResolveHostDirs(userHome)

	var targets []HostType
	switch strings.ToLower(targetHost) {
	case "cursor":
		targets = []HostType{HostCursor}
	case "claude":
		targets = []HostType{HostClaudeCode}
	case "antigravity", "gemini":
		targets = []HostType{HostAntigravity}
	case "agents":
		targets = []HostType{HostAgents}
	case "all", "":
		targets = []HostType{HostCursor, HostClaudeCode, HostAntigravity, HostAgents}
	default:
		return fmt.Errorf("unknown target host: %s", targetHost)
	}

	for _, host := range targets {
		destDir := hostDirs[host]
		if err := resources.CheckQuarantineBoundary(destDir); err != nil {
			return fmt.Errorf("target directory %s violates quarantine: %w", destDir, err)
		}

		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed creating host dir %s: %w", destDir, err)
		}

		for _, skill := range CanonicalActiveSkills {
			srcSkill := filepath.Join(sourceDir, skill)
			destSkill := filepath.Join(destDir, skill)

			if _, err := os.Stat(srcSkill); os.IsNotExist(err) {
				continue
			}

			if err := copyDir(srcSkill, destSkill); err != nil {
				return fmt.Errorf("failed copying skill %s to %s: %w", skill, destSkill, err)
			}
		}
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

// copyDir recursively copies files from src to dst
func copyDir(src string, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

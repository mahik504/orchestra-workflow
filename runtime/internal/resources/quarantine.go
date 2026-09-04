package resources

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrQuarantinedPath is returned when any loader attempts to read from or reference the quarantined library
	ErrQuarantinedPath = errors.New("access to quarantined skills_library is strictly forbidden by policy")
)

// BannedPathSubstrings defines path fragments that trigger immediate quarantine rejection
var BannedPathSubstrings = []string{
	"skills_library",
	"skills-library",
	"skills_archive",
	"skills-archive",
	"curated_catalog/quarantine",
	"skills~", // Wildcard prefix matching SKILLS~1, SKILLS~2, etc.
	"curate~",
}

// CheckQuarantineBoundary inspects an incoming filesystem path or URL and rejects it if it points to quarantined paths
func CheckQuarantineBoundary(path string) error {
	lookupPath := path

	// 1. URL unescape to defeat percent-encoding evasion (%73kills_library, %5F, etc.)
	if unescaped, err := url.PathUnescape(lookupPath); err == nil {
		lookupPath = unescaped
	}

	normalized := strings.ToLower(strings.ReplaceAll(lookupPath, "\\", "/"))
	for _, banned := range BannedPathSubstrings {
		if strings.Contains(normalized, banned) {
			return fmt.Errorf("%w: path '%s' violates quarantine boundaries", ErrQuarantinedPath, path)
		}
	}

	// If this is a URL scheme (e.g. http://, https://, but not file://), do not perform filesystem symlink resolution
	lowerPath := strings.ToLower(lookupPath)
	if strings.HasPrefix(lowerPath, "file:///") {
		lookupPath = strings.TrimPrefix(lookupPath, "file:///")
	} else if strings.HasPrefix(lowerPath, "file://") {
		lookupPath = strings.TrimPrefix(lookupPath, "file://")
	} else if strings.Contains(lookupPath, "://") {
		return nil
	}

	fsPath := filepath.FromSlash(lookupPath)

	// Resolve symlinks and NTFS directory junctions if the path exists on disk
	if resolved, err := filepath.EvalSymlinks(fsPath); err == nil {
		normResolved := strings.ToLower(filepath.ToSlash(resolved))
		for _, banned := range BannedPathSubstrings {
			if strings.Contains(normResolved, banned) {
				return fmt.Errorf("%w: resolved path '%s' points to quarantined '%s'", ErrQuarantinedPath, path, resolved)
			}
		}
	} else {
		// If the leaf doesn't exist yet, evaluate parent directories to detect junctions/symlinks in the directory tree
		curr := filepath.Dir(fsPath)
		for curr != "" && curr != "." && curr != filepath.Dir(curr) {
			if resolvedParent, err := filepath.EvalSymlinks(curr); err == nil {
				normParent := strings.ToLower(filepath.ToSlash(resolvedParent))
				for _, banned := range BannedPathSubstrings {
					if strings.Contains(normParent, banned) {
						return fmt.Errorf("%w: parent path '%s' points to quarantined '%s'", ErrQuarantinedPath, curr, resolvedParent)
					}
				}
				break
			}
			curr = filepath.Dir(curr)
		}
	}

	return nil
}

// QuarantineStatus contains diagnostic metrics regarding skill library quarantine
type QuarantineStatus struct {
	QuarantineDirectoryExists bool     `json:"quarantine_directory_exists"`
	QuarantinePath            string   `json:"quarantine_path"`
	QuarantinedCount          int      `json:"quarantined_count"`
	IsStrictlyIsolated        bool     `json:"is_strictly_isolated"`
	ViolationsFound           []string `json:"violations_found"`
}

// AuditQuarantineState inspects the environment to verify that the 1,598-skill library remains quarantined
func AuditQuarantineState(workspaceRoot string, quarantineRoot string) (*QuarantineStatus, error) {
	status := &QuarantineStatus{
		QuarantinePath:     quarantineRoot,
		ViolationsFound:    make([]string, 0),
		IsStrictlyIsolated: true,
	}

	info, err := os.Stat(quarantineRoot)
	if err == nil && info.IsDir() {
		status.QuarantineDirectoryExists = true
		entries, err := os.ReadDir(quarantineRoot)
		if err == nil {
			status.QuarantinedCount = len(entries)
		}
	}

	// Verify that quarantineRoot is not within workspaceRoot
	if workspaceRoot != "" && quarantineRoot != "" {
		rel, err := filepath.Rel(workspaceRoot, quarantineRoot)
		if err == nil && !strings.HasPrefix(rel, "..") {
			status.IsStrictlyIsolated = false
			status.ViolationsFound = append(status.ViolationsFound, "Quarantine directory resides inside active workspace tree")
		}
	}

	return status, nil
}

// VerifyCatalogQuarantineClean ensures no resources in the catalog reference quarantined paths
func VerifyCatalogQuarantineClean(cat *ResourceCatalog) []string {
	var violations []string
	if cat == nil {
		return violations
	}
	for _, res := range cat.All() {
		if err := CheckQuarantineBoundary(res.CanonicalURL); err != nil {
			violations = append(violations, fmt.Sprintf("resource %s canonical_url: %v", res.ID, err))
		}
		if err := CheckQuarantineBoundary(res.SourceRepository); err != nil {
			violations = append(violations, fmt.Sprintf("resource %s source_repository: %v", res.ID, err))
		}
		if err := CheckQuarantineBoundary(res.DocumentationURL); err != nil {
			violations = append(violations, fmt.Sprintf("resource %s documentation_url: %v", res.ID, err))
		}
	}
	return violations
}

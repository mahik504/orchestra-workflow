package benchmark

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func HashTree(root string) (string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	h := sha256.New()
	for _, rel := range files {
		_, _ = io.WriteString(h, rel)
		_, _ = io.WriteString(h, "\n")
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return "", err
		}
		_, _ = h.Write(raw)
		_, _ = io.WriteString(h, "\n")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func TreeClean(root string) bool {
	bad := false
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		base := info.Name()
		if base == ".orchestra" || base == "mission.json" || base == "capability.json" || base == "packet.json" {
			bad = true
		}
		return nil
	})
	return !bad
}

var baselineForbidden = []string{
	".orchestra",
	"mission.json",
	"capability.json",
	"capsule.json",
	"resource-memory.json",
}

func BaselineIsolated(root string) error {
	var hits []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		name := info.Name()
		for _, ban := range baselineForbidden {
			if name == ban {
				hits = append(hits, name)
			}
		}
		return nil
	})
	raw, err := os.ReadFile(filepath.Join(root, "packet.json"))
	if err == nil {
		text := string(raw)
		for _, ban := range []string{"mission.json", "capability.json", "Brain fact", ".orchestra/runs"} {
			if strings.Contains(text, ban) && strings.Contains(text, `"orchestration": "off"`) {
				if ban != "mission.json" || strings.Contains(text, `"mission_path"`) {
					hits = append(hits, "packet:"+ban)
				}
			}
		}
	}
	if len(hits) > 0 {
		return fmt.Errorf("baseline is contaminated: %s", strings.Join(hits, ", "))
	}
	return nil
}

package handoff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Receipt is the executor's record that a capsule was run.
// PREPARED is the default. EXECUTED requires the capsule hash and a zero exit.
// VERIFIED also requires a PASS verdict.
type Receipt struct {
	CapsuleHash  string   `json:"capsule_hash"`
	Host         string   `json:"host"`
	Commands     []string `json:"commands"`
	ExitCodes    []int    `json:"exit_codes"`
	ChangedFiles []string `json:"changed_files"`
	Verdict      string   `json:"verdict,omitempty"`
	Status       string   `json:"status"`
}

func (r Receipt) StatusFor() string {
	ready := strings.TrimSpace(r.CapsuleHash) != "" && len(r.ExitCodes) > 0
	for _, code := range r.ExitCodes {
		if code != 0 {
			ready = false
		}
	}
	if ready && r.Verdict == "PASS" {
		return "VERIFIED"
	}
	if ready {
		return "EXECUTED"
	}
	return "PREPARED"
}

func WriteReceipt(root string, receipt Receipt) (string, error) {
	receipt.Status = receipt.StatusFor()
	dir := filepath.Join(root, ".orchestra", "handoff")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "receipt.json")
	raw, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, append(raw, '\n'), 0600); err != nil {
		return "", err
	}
	return path, nil
}

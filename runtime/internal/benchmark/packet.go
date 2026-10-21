package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Packet struct {
	Orchestration  string   `json:"orchestration"`
	Task           string   `json:"task"`
	Acceptance     []string `json:"acceptance"`
	FixtureID      string   `json:"fixture_id"`
	FixtureHash    string   `json:"fixture_hash"`
	TimeBudget     string   `json:"time_budget"`
	HumanPolicy    string   `json:"human_policy"`
	ModelLabel     string   `json:"model_label,omitempty"`
	MissionPath    string   `json:"mission_path,omitempty"`
	CapabilityPath string   `json:"capability_path,omitempty"`
	CapsulePath    string   `json:"capsule_path,omitempty"`
	Intensity      string   `json:"intensity,omitempty"`
}

func ReadTask(root string) (string, []string, error) {
	task, err := os.ReadFile(filepath.Join(root, "task.txt"))
	if err != nil {
		return "", nil, err
	}
	accRaw, err := os.ReadFile(filepath.Join(root, "acceptance.txt"))
	if err != nil {
		return "", nil, err
	}
	var acc []string
	for _, line := range strings.Split(string(accRaw), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			acc = append(acc, line)
		}
	}
	return strings.TrimSpace(string(task)), acc, nil
}

func WritePacket(dir string, p Packet) error {
	if p.Orchestration == "off" {
		p.MissionPath = ""
		p.CapabilityPath = ""
		p.CapsulePath = ""
		p.Intensity = ""
	}
	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	if p.Orchestration == "off" && (strings.Contains(string(raw), "mission_path") || strings.Contains(string(raw), "capability_path") || strings.Contains(string(raw), "capsule_path")) {
		return fmt.Errorf("baseline packet gained an orchestra path")
	}
	return os.WriteFile(filepath.Join(dir, "packet.json"), raw, 0644)
}

package onboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

func writeProofJSON(dir, name string, v any) error {
	if strings.TrimSpace(dir) == "" {
		return nil
	}
	if err := resources.CheckQuarantineBoundary(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), data, 0644)
}

func writeProofText(dir, name, text string) error {
	if strings.TrimSpace(dir) == "" {
		return nil
	}
	if err := resources.CheckQuarantineBoundary(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), []byte(text), 0644)
}

func proof(dir, name string, v any) {
	if err := writeProofJSON(dir, name, v); err != nil {
		fmt.Fprintf(os.Stderr, "lifecycle artifact %s: %v\n", name, err)
	}
}

package mcpserver

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeCompileAndProvenance(t *testing.T) {
	routes := filepath.Join("..", "..", "..", "registries", "routing.json")
	var in bytes.Buffer
	write := func(msg any) {
		raw, err := json.Marshal(msg)
		if err != nil {
			t.Fatal(err)
		}
		in.Write(raw)
		in.WriteByte('\n')
	}
	write(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize"})
	write(map[string]any{"jsonrpc": "2.0", "method": "notifications/initialized"})
	write(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
	write(map[string]any{
		"jsonrpc": "2.0", "id": 3, "method": "tools/call",
		"params": map[string]any{"name": "compile", "arguments": map[string]any{"task": "The failing test is a parser crash", "host": "cursor"}},
	})
	write(map[string]any{
		"jsonrpc": "2.0", "id": 4, "method": "tools/call",
		"params": map[string]any{"name": "provenance_check", "arguments": map[string]any{
			"resource_id": "demo", "acquisition_method": "npm", "source_url": "https://orchestra.internal/resources/demo",
			"kind": "REAL_ARTIFACT", "version_or_sha": "1.0.0", "sha256_hash": "abc", "justification_task_id": "t",
		}},
	})
	var out bytes.Buffer
	if err := Serve(&in, &out, routes); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("responses %d\n%s", len(lines), out.String())
	}
	var compiled map[string]any
	if err := json.Unmarshal([]byte(lines[2]), &compiled); err != nil {
		t.Fatal(err)
	}
	result := compiled["result"].(map[string]any)
	content := result["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(content, "Bug") {
		t.Fatalf("compile text %s", content)
	}
	if !strings.Contains(lines[3], "synthetic source url") {
		t.Fatalf("provenance %s", lines[3])
	}
}

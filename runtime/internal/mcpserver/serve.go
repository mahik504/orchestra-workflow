package mcpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mahik504/orchestra-workflow/runtime/internal/acquisition"
	"github.com/mahik504/orchestra-workflow/runtime/internal/compile"
)

// Serve speaks one JSON-RPC message per line on stdin and stdout.
// Experimental. Hosts can call compile and the provenance check without a second conductor.
func Serve(r io.Reader, w io.Writer, routesPath string) error {
	routes, err := compile.LoadRoutes(routesPath)
	if err != nil {
		return err
	}
	in := bufio.NewScanner(r)
	in.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	out := bufio.NewWriter(w)
	for in.Scan() {
		line := in.Bytes()
		if len(line) == 0 {
			continue
		}
		var req map[string]any
		if err := json.Unmarshal(line, &req); err != nil {
			if err := writeMsg(out, map[string]any{
				"jsonrpc": "2.0",
				"error":   map[string]any{"code": -32700, "message": "parse error"},
			}); err != nil {
				return err
			}
			continue
		}
		id, hasID := req["id"]
		method, _ := req["method"].(string)
		if !hasID {
			continue
		}
		result, callErr := dispatch(method, req["params"], routes)
		msg := map[string]any{"jsonrpc": "2.0", "id": id}
		if callErr != nil {
			msg["error"] = map[string]any{"code": -32000, "message": callErr.Error()}
		} else {
			msg["result"] = result
		}
		if err := writeMsg(out, msg); err != nil {
			return err
		}
	}
	return in.Err()
}

func dispatch(method string, params any, routes []compile.Route) (any, error) {
	switch method {
	case "initialize":
		return map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "orchestra", "version": "experimental"},
		}, nil
	case "tools/list":
		return map[string]any{"tools": []any{
			map[string]any{"name": "compile", "description": "Compile a task into a route manifest. Does not install anything."},
			map[string]any{"name": "provenance_check", "description": "Validate one provenance row. A synthetic orchestra.internal URL fails."},
		}}, nil
	case "tools/call":
		return callTool(params, routes)
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func callTool(params any, routes []compile.Route) (any, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var call struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(raw, &call); err != nil {
		return nil, err
	}
	switch call.Name {
	case "compile":
		task, _ := call.Arguments["task"].(string)
		host, _ := call.Arguments["host"].(string)
		if host == "" {
			host = "cursor"
		}
		manifest, err := compile.Compile(task, host, routes)
		if err != nil {
			return nil, err
		}
		body, err := json.Marshal(manifest)
		if err != nil {
			return nil, err
		}
		return map[string]any{"content": []any{map[string]any{"type": "text", "text": string(body)}}}, nil
	case "provenance_check":
		entry := acquisition.ProvenanceEntry{
			ResourceID:          str(call.Arguments["resource_id"]),
			AcquisitionMethod:   str(call.Arguments["acquisition_method"]),
			SourceURL:           str(call.Arguments["source_url"]),
			VersionOrSHA:        str(call.Arguments["version_or_sha"]),
			SHA256Hash:          str(call.Arguments["sha256_hash"]),
			JustificationTaskID: str(call.Arguments["justification_task_id"]),
			Kind:                str(call.Arguments["kind"]),
		}
		if err := acquisition.ValidateEntry(&entry); err != nil {
			return map[string]any{"content": []any{map[string]any{"type": "text", "text": err.Error()}}, "isError": true}, nil
		}
		return map[string]any{"content": []any{map[string]any{"type": "text", "text": "ok"}}}, nil
	default:
		return nil, fmt.Errorf("unknown tool %s", call.Name)
	}
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func writeMsg(w *bufio.Writer, msg any) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return w.Flush()
}

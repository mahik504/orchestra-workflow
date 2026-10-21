package visual

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const playwrightChromiumProbe = `import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import path from "node:path";
import fs from "node:fs";
const require = createRequire(pathToFileURL(path.join(process.cwd(), "package.json")).href);
const { chromium } = require("playwright");
const p = chromium.executablePath();
if (!p || !fs.existsSync(p)) { process.stderr.write("chromium missing"); process.exit(2); }
process.stdout.write(p);
`

type PlaywrightProbe struct {
	Available  bool   `json:"available"`
	RuntimeDir string `json:"runtime_dir,omitempty"`
	Node       string `json:"node,omitempty"`
	Library    string `json:"library,omitempty"`
	Chromium   string `json:"chromium,omitempty"`
	Error      string `json:"error,omitempty"`
}

func PlaywrightAvailable() bool {
	return ProbePlaywright().Available
}

func ProbePlaywright() PlaywrightProbe {
	p := PlaywrightProbe{}
	rt := FindRuntimeDir()
	p.RuntimeDir = rt
	if rt == "" {
		p.Error = "runtime directory not found (expected runtime/package.json next to tools/layer-a.mjs)"
		return p
	}
	lib := filepath.Join(rt, "node_modules", "playwright")
	p.Library = lib
	if _, err := os.Stat(lib); err != nil {
		p.Error = "playwright library missing under runtime/node_modules"
		return p
	}
	node := lookupNode()
	p.Node = node
	if node == "" {
		p.Error = "node executable not on PATH (LookPath node.exe/node)"
		return p
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "--input-type=module", "-e", playwrightChromiumProbe)
	cmd.Dir = rt
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		p.Error = strings.TrimSpace(stderr.String() + " " + err.Error())
		return p
	}
	p.Chromium = strings.TrimSpace(string(out))
	if p.Chromium == "" {
		p.Error = "chromium executable path empty"
		return p
	}
	if _, err := os.Stat(p.Chromium); err != nil {
		p.Error = "chromium path does not exist: " + p.Chromium
		return p
	}
	p.Available = true
	return p
}

func lookupNode() string {
	names := []string{"node"}
	if runtime.GOOS == "windows" {
		names = []string{"node.exe", "node"}
	}
	for _, n := range names {
		if p, err := exec.LookPath(n); err == nil {
			return p
		}
	}
	return ""
}

func FindRuntimeDir() string {
	if script := FindLayerAScript(); script != "" {
		// runtime/tools/layer-a.mjs -> runtime
		dir := filepath.Dir(filepath.Dir(script))
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
	}
	_, file, _, ok := runtime.Caller(0)
	if ok {
		cand := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
		if _, err := os.Stat(filepath.Join(cand, "package.json")); err == nil {
			return cand
		}
	}
	for _, cand := range []string{"runtime", "."} {
		if _, err := os.Stat(filepath.Join(cand, "package.json")); err == nil {
			if _, err := os.Stat(filepath.Join(cand, "tools", "layer-a.mjs")); err == nil {
				abs, _ := filepath.Abs(cand)
				return abs
			}
		}
	}
	return ""
}

func killProcessTree(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil || cmd.Process.Pid <= 0 {
		return
	}
	pid := cmd.Process.Pid
	if runtime.GOOS == "windows" {
		taskkill, err := exec.LookPath("taskkill.exe")
		if err != nil {
			taskkill = "taskkill"
		}
		_ = exec.Command(taskkill, "/F", "/T", "/PID", strconv.Itoa(pid)).Run()
	}
	_ = cmd.Process.Kill()
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	if runtime.GOOS == "windows" {
		out, err := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid), "/NH").Output()
		if err != nil {
			return false
		}
		s := strings.ToLower(string(out))
		if strings.Contains(s, "no tasks") {
			return false
		}
		return strings.Contains(s, strconv.Itoa(pid))
	}
	return exec.Command("kill", "-0", strconv.Itoa(pid)).Run() == nil
}

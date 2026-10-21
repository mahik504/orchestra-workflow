package benchmark

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPhase3G2_ArgvContract(t *testing.T) {
	for _, help := range []string{"full", "print"} {
		t.Run(help, func(t *testing.T) {
			log := runWrapper(t, help, "opaque-label", "ship-the-sum")
			lines := nonEmptyLines(log)
			if len(lines) < 2 {
				t.Fatalf("argv log %q", log)
			}
			invoke := lines[len(lines)-1]
			for _, want := range []string{"--print", "--workspace", "--model", "opaque-label", "--", promptPointer} {
				if !strings.Contains(invoke, want) {
					t.Fatalf("missing %s in %s", want, invoke)
				}
			}
			if strings.Contains(invoke, "ship-the-sum") {
				t.Fatalf("prompt text leaked into argv: %s", invoke)
			}
			for _, banned := range []string{"--plan", "--worktree", "--add-dir", "--approve-mcps", "mcp.json", "orchestra-workflow", "api_key", "sk-"} {
				if strings.Contains(invoke, banned) {
					t.Fatalf("banned %s in %s", banned, invoke)
				}
			}
			home, _ := os.UserHomeDir()
			if home != "" {
				for _, field := range strings.Fields(invoke) {
					if field == home {
						t.Fatalf("user profile was an argument")
					}
				}
			}
			if help == "full" {
				if !strings.Contains(invoke, "--force") || !strings.Contains(invoke, "--trust") {
					t.Fatalf("supported flags missing: %s", invoke)
				}
			} else if strings.Contains(invoke, "--force") || strings.Contains(invoke, "--trust") {
				t.Fatalf("unsupported flags present: %s", invoke)
			}
		})
	}
}

const promptPointer = "Read prompt.txt in the workspace root and complete the task it describes, including its acceptance lines."

func TestPhase3G2_LauncherKeepsArgs(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "argv.txt")
	ps1 := filepath.Join(dir, "log.ps1")
	ps1Body := "if ($args.Count -gt 0 -and $args[0] -eq '--help') { Write-Output '--print --force --trust --workspace --model'; exit 0 }\r\n" +
		"Add-Content -LiteralPath '" + logPath + "' -Value ($args -join ' ')\r\n"
	if err := os.WriteFile(ps1, []byte(ps1Body), 0644); err != nil {
		t.Fatal(err)
	}
	agent := filepath.Join(dir, "agent.cmd")
	cmdBody := "@echo off\r\nsetlocal enabledelayedexpansion\r\n" +
		"%SystemRoot%\\System32\\WindowsPowerShell\\v1.0\\powershell.exe -NoProfile -ExecutionPolicy Bypass -File \"" + ps1 + "\" %*\r\n"
	if err := os.WriteFile(agent, []byte(cmdBody), 0644); err != nil {
		t.Fatal(err)
	}
	ws := filepath.Join(dir, "ws")
	if err := os.MkdirAll(ws, 0755); err != nil {
		t.Fatal(err)
	}
	prompt := "Fix Sum.\r\n\r\nAcceptance:\r\n- Sum(1, 2) returns 3\r\n"
	if err := os.WriteFile(filepath.Join(ws, "prompt.txt"), []byte(prompt), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(wrapperCmd(t))
	cmd.Env = append(os.Environ(),
		"ORCHESTRA_BENCH_AGENT="+agent,
		"ORCHESTRA_BENCH_WORKSPACE="+ws,
		"ORCHESTRA_BENCH_PROMPT="+filepath.Join(ws, "prompt.txt"),
		"ORCHESTRA_BENCH_MODEL_LABEL=opaque-label",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wrapper: %v %s", err, out)
	}
	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	for _, want := range []string{"--model", "opaque-label", "--", promptPointer} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "Acceptance:") || strings.Contains(got, "Fix Sum.") {
		t.Fatalf("prompt body reached the launcher: %s", got)
	}
}

func TestPhase3G2_SlugIsExact(t *testing.T) {
	text := "grok-4.7-high - Grok 4.7\n"
	if slugListed(text, "grok-4.7") {
		t.Fatal("a longer id is not the operator label")
	}
	if !slugListed(text, "grok-4.7-high") {
		t.Fatal("exact id was missed")
	}
}

func TestPhase3G2_PreflightMissing(t *testing.T) {
	t.Setenv("ORCHESTRA_BENCH_AGENT", filepath.Join(t.TempDir(), "missing.cmd"))
	t.Setenv("ORCHESTRA_BENCH_MODEL_LABEL", "opaque-label")
	got := PreflightAgent(context.Background())
	if got.OK || got.Blocker == "" || strings.Contains(got.Blocker, `\`) {
		t.Fatalf("%+v", got)
	}
}

func TestPhase3G2_Fingerprints(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t1-api", "noop", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || !res.PromptEquivalent || res.Claim != ClaimNone {
		t.Fatalf("%+v", res)
	}
	if res.TaskHash == "" || res.AcceptanceHash == "" || res.BaselinePromptHash == "" || res.OrchestraPromptHash == "" {
		t.Fatalf("missing hashes %+v", res)
	}
	if res.BaselinePromptHash == res.OrchestraPromptHash {
		t.Fatal("prompt hashes should differ by the harness references")
	}
	if res.Baseline.TaskHash != res.Orchestra.TaskHash || res.Baseline.AcceptanceHash != res.Orchestra.AcceptanceHash {
		t.Fatal("task or acceptance diverged")
	}
	if res.Baseline.HostEdited || res.Baseline.TaskCompleted || res.Orchestra.HostEdited || res.Orchestra.TaskCompleted {
		t.Fatalf("t1 noop base edited=%v done=%v orch edited=%v done=%v", res.Baseline.HostEdited, res.Baseline.TaskCompleted, res.Orchestra.HostEdited, res.Orchestra.TaskCompleted)
	}
}

func TestPhase3G2_NonZeroExit(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "exit2", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "host_error" || res.Claim != ClaimNone {
		t.Fatalf("%+v", res.Status)
	}
	if !hasLimitation(res, "not a win") {
		t.Fatal(res.Limitations)
	}
}

func TestPhase3G2_BaselineContamination(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "dirtybase", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "invalid_pair" || res.BaselineIsolated || res.Claim != ClaimNone {
		t.Fatalf("status %s isolated %v", res.Status, res.BaselineIsolated)
	}
}

func TestPhase3G2_LabelLeak(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "leak", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "invalid_pair" || res.Claim != ClaimNone || !hasLimitation(res, "model label") {
		t.Fatalf("%s %v", res.Status, res.Limitations)
	}
}

func TestPhase3G2_TouchIsNotCompletion(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "touch", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone {
		t.Fatalf("%s", res.Status)
	}
	if !res.Baseline.HostEdited || res.Baseline.TaskCompleted || !res.Orchestra.HostEdited || res.Orchestra.TaskCompleted {
		t.Fatalf("base edited=%v done=%v orch edited=%v done=%v", res.Baseline.HostEdited, res.Baseline.TaskCompleted, res.Orchestra.HostEdited, res.Orchestra.TaskCompleted)
	}
}

func TestPhase3G2_ClaimStaysNone(t *testing.T) {
	res, err := RunPair(context.Background(), pairInput(t, "t3-debug", "fixorch", 30*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "compared" || res.Claim != ClaimNone {
		t.Fatalf("%s %s", res.Status, res.Claim)
	}
	if res.Baseline.TaskCompleted || !res.Orchestra.TaskCompleted {
		t.Fatalf("base done=%v orch done=%v", res.Baseline.TaskCompleted, res.Orchestra.TaskCompleted)
	}
	if res.Orchestra.HostEdited != true {
		t.Fatal("orchestra edit was not recorded")
	}
}

func TestPhase3G2_ExternalWrite(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(outside, []byte("{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ORCHESTRA_BENCH_OUTSIDE", outside)
	in := pairInput(t, "t1-api", "outside", 30*time.Second)
	in.Guards = []Guard{{Class: "cursor_config", Path: outside}}
	res, err := RunPair(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "invalid_pair" || res.ExternalClass != "cursor_config" || res.Claim != ClaimNone {
		t.Fatalf("%s %s", res.Status, res.ExternalClass)
	}
	raw, err := os.ReadFile(in.ResultPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), outside) {
		t.Fatal("result stored a protected path")
	}
}

func TestPhase3G2_RealHostPilot(t *testing.T) {
	realPilot(t, "3g2", []string{"t1-api", "t3-debug"})
}

func realPilot(t *testing.T, prefix string, ids []string) {
	if os.Getenv("ORCHESTRA_BENCH_REAL") != "1" {
		t.Skip("real host pilot is opt-in")
	}
	dir := filepath.Join("..", "..", "benchmarks", "results")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	pf := PreflightAgent(context.Background())
	raw, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, prefix+"-preflight.json"), raw, 0644); err != nil {
		t.Fatal(err)
	}
	if !pf.OK {
		t.Logf("preflight unavailable: %s", pf.Blocker)
		for _, id := range ids {
			writeUnavailable(t, dir, prefix, id, pf.Blocker)
		}
		return
	}
	host := wrapperCmd(t)
	label := os.Getenv("ORCHESTRA_BENCH_MODEL_LABEL")
	var t1Status string
	for _, id := range ids {
		if id == "t3-debug" && (t1Status == "unavailable" || t1Status == "invalid_pair") {
			writeUnavailable(t, dir, prefix, id, "stopped after T1 "+t1Status)
			t.Logf("%s skipped after T1 %s", id, t1Status)
			return
		}
		cat, graph := loadCatalogs(t)
		root := t.TempDir()
		in := PairInput{
			Catalog:        cat,
			Graph:          graph,
			FixtureDir:     fixtureDir(t, id),
			Task:           taskByID(id),
			WorkRoot:       filepath.Join(root, "work"),
			Executable:     host,
			Timeout:        10 * time.Minute,
			HostLabel:      "cursor",
			BaselineLabel:  label,
			OrchestraLabel: label,
			ResultPath:     filepath.Join(dir, prefix+"-"+id+".json"),
			SharedProfile:  true,
			Guards:         realGuards(fixtureDir(t, id)),
		}
		if err := os.MkdirAll(in.WorkRoot, 0755); err != nil {
			t.Fatal(err)
		}
		res, err := RunPair(context.Background(), in)
		if err != nil {
			t.Fatal(err)
		}
		if res.Claim != ClaimNone {
			t.Fatalf("%s claim %s", id, res.Claim)
		}
		if id == "t1-api" {
			t1Status = res.Status
		}
		t.Logf("%s status %s", id, res.Status)
	}
}

func writeUnavailable(t *testing.T, dir, prefix, id, why string) {
	t.Helper()
	task := taskByID(id)
	res := PairResult{
		BenchmarkID: "3g1-" + id,
		TaskClass:   task.Class,
		FixtureID:   id,
		Phase:       PhasePilot,
		Trial:       1,
		Status:      "unavailable",
		Claim:       ClaimNone,
		Tokens:      TokensUnavailable,
		Limitations: []string{"pilot observation", "claim none", "one trial", why},
	}
	if err := writePair(filepath.Join(dir, prefix+"-"+id+".json"), res); err != nil {
		t.Fatal(err)
	}
}

func runWrapper(t *testing.T, help, label, prompt string) string {
	t.Helper()
	dir := t.TempDir()
	agent := filepath.Join(dir, "agent.cmd")
	logPath := filepath.Join(dir, "argv.txt")
	script := "@echo off\r\necho %*>>\"" + logPath + "\"\r\n" +
		"if /I \"%~1\"==\"--help\" (\r\n" +
		"  if /I \"%ORCHESTRA_FAKE_HELP%\"==\"print\" (\r\n" +
		"    echo --print --workspace --model\r\n" +
		"  ) else (\r\n" +
		"    echo --print --force --trust --workspace --model\r\n" +
		"  )\r\n" +
		"  exit /b 0\r\n" +
		")\r\n" +
		"exit /b 0\r\n"
	if err := os.WriteFile(agent, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	ws := filepath.Join(dir, "ws")
	if err := os.MkdirAll(ws, 0755); err != nil {
		t.Fatal(err)
	}
	promptPath := filepath.Join(ws, "prompt.txt")
	if err := os.WriteFile(promptPath, []byte(prompt), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(wrapperCmd(t))
	cmd.Env = append(os.Environ(),
		"ORCHESTRA_BENCH_AGENT="+agent,
		"ORCHESTRA_BENCH_WORKSPACE="+ws,
		"ORCHESTRA_BENCH_PROMPT="+promptPath,
		"ORCHESTRA_BENCH_MODEL_LABEL="+label,
		"ORCHESTRA_FAKE_HELP="+help,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wrapper: %v %s", err, out)
	}
	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func wrapperCmd(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "benchmarks", "cursor-host.cmd"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	return path
}

func nonEmptyLines(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func hasLimitation(res PairResult, part string) bool {
	for _, line := range res.Limitations {
		if strings.Contains(line, part) {
			return true
		}
	}
	return false
}

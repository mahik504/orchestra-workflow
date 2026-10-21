package benchmark

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ScopePolicy struct {
	Protect []string
	Allow   []string
	Oracle  string
}

func PolicyFor(fixtureID string) (ScopePolicy, error) {
	switch fixtureID {
	case "t1-api":
		return ScopePolicy{
			Protect: []string{"go.mod", "task.txt", "acceptance.txt"},
			Allow:   []string{"sum.go", "sum_test.go"},
			Oracle:  t1Oracle,
		}, nil
	case "t3-debug":
		return ScopePolicy{
			Protect: []string{"sum_test.go", "task.txt", "acceptance.txt", "go.mod"},
			Allow:   []string{"sum.go"},
		}, nil
	default:
		return ScopePolicy{}, fmt.Errorf("no scope policy for %s", fixtureID)
	}
}

func hashFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func protectHashes(root string, names []string) (map[string]string, error) {
	out := make(map[string]string, len(names))
	for _, name := range names {
		sum, err := hashFile(filepath.Join(root, name))
		if err != nil {
			return nil, err
		}
		out[name] = sum
	}
	return out, nil
}

func harnessPath(rel string) bool {
	rel = filepath.ToSlash(rel)
	switch rel {
	case "packet.json", "prompt.txt", "benchmark.test.exe", "host-stdout.txt", "host-stderr.txt", oracleFileName:
		return true
	}
	return rel == ".git" || strings.HasPrefix(rel, ".git/") || rel == ".orchestra" || strings.HasPrefix(rel, ".orchestra/")
}

func ScoreHost(root string, policy ScopePolicy, before map[string]string, judgment *Judgment) (Record, error) {
	art, err := diffArtifact(root)
	if err != nil {
		return Record{}, err
	}
	var outside []string
	allow := map[string]bool{}
	for _, name := range policy.Allow {
		allow[name] = true
	}
	for _, name := range art.Changed {
		if harnessPath(name) {
			continue
		}
		if !allow[name] {
			outside = append(outside, name)
		}
		if name == "go.mod" {
			art.NewDeps = true
		}
	}
	art.OutsideAllow = outside

	var checks []Check
	testsPass := false
	buildOK := false
	build := exec.Command("go", "test", "-c", "-o", filepath.Join(root, "benchmark.test.exe"))
	build.Dir = root
	build.Env = scoreEnv(root)
	if err := build.Run(); err != nil {
		checks = append(checks, Check{Name: DimBuild, Pass: false})
		checks = append(checks, Check{Name: DimTests, Pass: false})
	} else {
		buildOK = true
		_ = os.Remove(filepath.Join(root, "benchmark.test.exe"))
		checks = append(checks, Check{Name: DimBuild, Pass: true})
		test := exec.Command("go", "test", ".")
		test.Dir = root
		test.Env = scoreEnv(root)
		out, err := test.CombinedOutput()
		testsPass = err == nil && !strings.Contains(string(out), "FAIL")
		checks = append(checks, Check{Name: DimTests, Pass: testsPass})
	}
	oraclePass := true
	if policy.Oracle != "" && buildOK {
		var oerr error
		oraclePass, oerr = runOracle(root, policy.Oracle)
		if oerr != nil {
			return Record{}, oerr
		}
	}
	behavior := testsPass
	if policy.Oracle != "" {
		behavior = testsPass && oraclePass
	}
	checks = append(checks, Check{Name: DimRequiredBehavior, Pass: behavior})
	accept := true
	for _, name := range policy.Protect {
		sum, err := hashFile(filepath.Join(root, name))
		if err != nil || sum != before[name] {
			accept = false
			break
		}
	}
	checks = append(checks, Check{Name: DimAcceptance, Pass: accept})
	checks = append(checks, Check{Name: DimScope, Pass: len(outside) == 0})
	rec := Record{
		Tokens:        TokensUnavailable,
		TaskExecution: TaskCostUnavailable,
		HostRetries:   TokensUnavailable,
		Claim:         ClaimNone,
		Phase:         PhasePilot,
		Artifact:      &art,
	}
	ApplyScore(&rec, checks, judgment)
	return rec, nil
}

func diffArtifact(root string) (Artifact, error) {
	stat, err := gitOut(root, "diff", "--stat")
	if err != nil {
		return Artifact{}, err
	}
	status, err := gitOut(root, "status", "--porcelain")
	if err != nil {
		return Artifact{}, err
	}
	full, err := gitOut(root, "diff")
	if err != nil {
		return Artifact{}, err
	}
	var names []string
	for _, line := range strings.Split(status, "\n") {
		if len(strings.TrimSpace(line)) < 4 {
			continue
		}
		name := line[3:]
		if i := strings.LastIndex(name, " -> "); i >= 0 {
			name = name[i+4:]
		}
		name = filepath.ToSlash(strings.TrimSpace(name))
		if name == "" || harnessPath(name) {
			continue
		}
		names = append(names, name)
		if strings.HasPrefix(line, "??") {
			raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
			if err == nil {
				full += string(raw)
			}
		}
	}
	return Artifact{DiffStat: strings.TrimSpace(stat), DiffBytes: len(full), Changed: names}, nil
}

func gitOut(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w (%s)", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func commitSnapshot(dir string) error {
	steps := [][]string{
		{"init"},
		{"add", "-A"},
		{"-c", "user.name=bench", "-c", "user.email=bench@localhost", "-c", "commit.gpgsign=false", "commit", "-m", "start"},
	}
	for _, args := range steps {
		if _, err := gitOut(dir, args...); err != nil {
			return err
		}
	}
	return nil
}

func scoreEnv(root string) []string {
	cache := filepath.Join(root, ".score-cache")
	mod := filepath.Join(root, ".score-mod")
	_ = os.MkdirAll(cache, 0755)
	_ = os.MkdirAll(mod, 0755)
	return hostEnv(root, filepath.Join(root, "prompt.txt"), cache, mod, "")
}

func baselinePrompt(task string, acc []string) string {
	var b strings.Builder
	b.WriteString(task)
	b.WriteString("\n\nAcceptance:\n")
	for _, line := range acc {
		b.WriteString("- ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func orchestraPrompt(task string, acc []string, missionPath, capsulePath, capabilityPath string) string {
	text := baselinePrompt(task, acc)
	text += "\nImplement the capsule.\n"
	text += "Mission file: " + missionPath + "\n"
	text += "Capsule file: " + capsulePath + "\n"
	text += "Capability file: " + capabilityPath + "\n"
	return text
}

func textHasLabel(text, label string) bool {
	label = strings.TrimSpace(label)
	if label == "" {
		return false
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(label))
}

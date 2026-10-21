package benchmark

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ScoreWorktree(root string, judgment *Judgment) (Record, error) {
	var checks []Check
	build := exec.Command("go", "test", "-c", "-o", filepath.Join(root, "benchmark.test.exe"))
	build.Dir = root
	if err := build.Run(); err != nil {
		checks = append(checks, Check{Name: DimBuild, Pass: false})
		checks = append(checks, Check{Name: DimTests, Pass: false})
	} else {
		checks = append(checks, Check{Name: DimBuild, Pass: true})
		_ = os.Remove(filepath.Join(root, "benchmark.test.exe"))
		test := exec.Command("go", "test", ".")
		test.Dir = root
		out, err := test.CombinedOutput()
		checks = append(checks, Check{Name: DimTests, Pass: err == nil && !strings.Contains(string(out), "FAIL")})
	}
	task, acc, err := ReadTask(root)
	if err != nil {
		checks = append(checks, Check{Name: DimAcceptance, Pass: false})
		checks = append(checks, Check{Name: DimRequiredBehavior, Pass: false})
	} else {
		ok := task != "" && len(acc) > 0
		checks = append(checks, Check{Name: DimAcceptance, Pass: ok})
		checks = append(checks, Check{Name: DimRequiredBehavior, Pass: ok})
	}
	rec := Record{Tokens: TokensUnavailable, TaskExecution: TaskCostUnavailable, Claim: ClaimNone, Phase: PhasePilot}
	ApplyScore(&rec, checks, judgment)
	return rec, nil
}

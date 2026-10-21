package benchmark

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const oracleFileName = "zz_bench_oracle_test.go"

// t1Oracle is the harness check for t1-api. It is written only while scoring,
// after the diff is taken, and then removed. The host never sees it.
const t1Oracle = `package t1api

import "testing"

func TestBenchOracleValid(t *testing.T) {
	got, err := Sum(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 {
		t.Fatalf("got %d", got)
	}
}

func TestBenchOracleNegative(t *testing.T) {
	if _, err := Sum(-1, 2); err == nil {
		t.Fatal("expected an error")
	}
	if _, err := Sum(2, -1); err == nil {
		t.Fatal("expected an error")
	}
}
`

func runOracle(root, src string) (bool, error) {
	path := filepath.Join(root, oracleFileName)
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		return false, err
	}
	defer os.Remove(path)
	test := exec.Command("go", "test", "-run", "^TestBenchOracle", ".")
	test.Dir = root
	test.Env = scoreEnv(root)
	out, err := test.CombinedOutput()
	pass := err == nil && !strings.Contains(string(out), "FAIL")
	return pass, nil
}

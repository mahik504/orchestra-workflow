package handoff

import (
	"os"
	"strings"
	"testing"
)

func TestReceiptStatus(t *testing.T) {
	if (Receipt{}).StatusFor() != "PREPARED" {
		t.Fatal("empty receipt")
	}
	if (Receipt{CapsuleHash: "abc", ExitCodes: []int{1}}).StatusFor() != "PREPARED" {
		t.Fatal("nonzero exit")
	}
	if (Receipt{CapsuleHash: "abc", ExitCodes: []int{0}}).StatusFor() != "EXECUTED" {
		t.Fatal("zero exit")
	}
	got := Receipt{CapsuleHash: "abc", ExitCodes: []int{0}, Verdict: "PASS"}.StatusFor()
	if got != "VERIFIED" {
		t.Fatal(got)
	}
}

func TestWriteReceipt(t *testing.T) {
	dir := t.TempDir()
	path, err := WriteReceipt(dir, Receipt{CapsuleHash: "abc", Host: "cursor", Commands: []string{"go test"}, ExitCodes: []int{0}})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"status": "EXECUTED"`) {
		t.Fatalf("%s", raw)
	}
}

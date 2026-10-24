package classifier

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadKind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "idea.md")
	body := "---\nkind: hiring-cv\n---\n\n# Idea\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := ReadKind(path); got != "hiring-cv" {
		t.Fatalf("ReadKind = %q", got)
	}
	if got := ReadKind(filepath.Join(dir, "missing.md")); got != "" {
		t.Fatalf("missing file = %q", got)
	}
}

func TestApplyKindCollegeStaysStandard(t *testing.T) {
	c := NewClassifier()
	dir := t.TempDir()
	path := filepath.Join(dir, "idea.md")
	if err := os.WriteFile(path, []byte("---\nkind: college\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	b := c.ClassifyBrief("build a small page", Options{IdeaPath: path})
	if b.QualityBar != BarStandard {
		t.Fatalf("college bar = %q (%s)", b.QualityBar, b.QualityBarReason)
	}
}

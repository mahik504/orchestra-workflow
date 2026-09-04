package resources

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// CHALLENGE SUITE 1: BOM & Encoding Adversarial Tests
// ============================================================================

func TestChallenger_BOM_UTF8_ValidAndEdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	validArrayJSON := `[{"id":"res-bom-1","name":"BOM Test 1","canonical_url":"https://example.com/1","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag1"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`
	validGraphJSON := `{"domains":{"dom-1":["res-bom-1"]},"capabilities":{"cap-1":{"name":"Cap 1","primary_archetype":"premium-website","discovery":["dom-1"],"synthesis":["dom-1"],"implementation":["res-bom-1"],"qa":["res-bom-1"]}}}`

	// 1. Standard UTF-8 BOM
	bomUTF8 := []byte("\xef\xbb\xbf")

	catPath := filepath.Join(tempDir, "catalog_utf8_bom.json")
	if err := os.WriteFile(catPath, append(bomUTF8, []byte(validArrayJSON)...), 0644); err != nil {
		t.Fatalf("failed to write catalog_utf8_bom: %v", err)
	}

	graphPath := filepath.Join(tempDir, "graph_utf8_bom.json")
	if err := os.WriteFile(graphPath, append(bomUTF8, []byte(validGraphJSON)...), 0644); err != nil {
		t.Fatalf("failed to write graph_utf8_bom: %v", err)
	}

	cat, err := LoadResourceCatalog(catPath)
	if err != nil {
		t.Fatalf("LoadResourceCatalog failed on UTF-8 BOM: %v", err)
	}
	if cat.Count() != 1 {
		t.Errorf("expected 1 resource, got %d", cat.Count())
	}

	graph, err := LoadDesignGraph(graphPath)
	if err != nil {
		t.Fatalf("LoadDesignGraph failed on UTF-8 BOM: %v", err)
	}
	if len(graph.Capabilities) != 1 {
		t.Errorf("expected 1 capability, got %d", len(graph.Capabilities))
	}

	// 2. Registry.LoadFromJSON with UTF-8 BOM
	reg := NewRegistry()
	if err := reg.LoadFromJSON(catPath); err != nil {
		t.Fatalf("Registry.LoadFromJSON failed on UTF-8 BOM: %v", err)
	}
	if len(reg.Capabilities) != 1 {
		t.Errorf("expected 1 capability in registry, got %d", len(reg.Capabilities))
	}

	// 3. UTF-8 BOM on empty content (BOM only)
	emptyBOMPath := filepath.Join(tempDir, "empty_bom.json")
	_ = os.WriteFile(emptyBOMPath, bomUTF8, 0644)

	_, err = LoadResourceCatalog(emptyBOMPath)
	if err == nil {
		t.Errorf("expected error on BOM-only empty file, got nil")
	} else {
		t.Logf("Empirical result (BOM-only file error): %v", err)
	}

	// 4. UTF-8 BOM followed by whitespace only
	wsBOMPath := filepath.Join(tempDir, "whitespace_bom.json")
	_ = os.WriteFile(wsBOMPath, append(bomUTF8, []byte("   \r\n\t  ")...), 0644)
	_, err = LoadResourceCatalog(wsBOMPath)
	if err == nil {
		t.Errorf("expected error on BOM-with-whitespace-only file, got nil")
	} else {
		t.Logf("Empirical result (BOM-with-whitespace error): %v", err)
	}
}

func TestChallenger_BOM_CorruptedHeaders(t *testing.T) {
	tempDir := t.TempDir()
	validArrayJSON := `[{"id":"res-bom-corrupt","name":"Corrupt BOM","canonical_url":"https://example.com/c","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag1"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`

	testCases := []struct {
		name        string
		headerBytes []byte
		description string
	}{
		{"1-byte truncation", []byte("\xef"), "truncated 1st byte of UTF-8 BOM"},
		{"2-byte truncation", []byte("\xef\xbb"), "truncated 2 bytes of UTF-8 BOM"},
		{"inverted 3rd byte", []byte("\xef\xbb\xbe"), "invalid 3rd byte"},
		{"bogus binary header", []byte("\x00\x01\x02\x03"), "arbitrary binary junk"},
		{"double UTF-8 BOM", []byte("\xef\xbb\xbf\xef\xbb\xbf"), "double BOM sequence"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(tempDir, fmt.Sprintf("corrupt_%s.json", strings.ReplaceAll(tc.name, " ", "_")))
			content := append(tc.headerBytes, []byte(validArrayJSON)...)
			if err := os.WriteFile(path, content, 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// LoadResourceCatalog should NOT panic and should cleanly return an error
			_, err := LoadResourceCatalog(path)
			if err == nil {
				t.Errorf("[%s] expected error for corrupted header %s, but LoadResourceCatalog succeeded", tc.name, tc.description)
			} else {
				t.Logf("[%s] LoadResourceCatalog rejected as expected: %v", tc.name, err)
			}

			// LoadDesignGraph should also cleanly return an error without panic
			_, err = LoadDesignGraph(path)
			if err == nil {
				t.Errorf("[%s] expected error for corrupted header %s, but LoadDesignGraph succeeded", tc.name, tc.description)
			} else {
				t.Logf("[%s] LoadDesignGraph rejected as expected: %v", tc.name, err)
			}
		})
	}
}

func encodeUTF16LE(s string) []byte {
	bom := []byte{0xff, 0xfe} // UTF-16 LE BOM
	runes := []rune(s)
	buf := new(bytes.Buffer)
	buf.Write(bom)
	for _, r := range runes {
		_ = binary.Write(buf, binary.LittleEndian, uint16(r))
	}
	return buf.Bytes()
}

func encodeUTF16BE(s string) []byte {
	bom := []byte{0xfe, 0xff} // UTF-16 BE BOM
	runes := []rune(s)
	buf := new(bytes.Buffer)
	buf.Write(bom)
	for _, r := range runes {
		_ = binary.Write(buf, binary.BigEndian, uint16(r))
	}
	return buf.Bytes()
}

func TestChallenger_BOM_UTF16(t *testing.T) {
	tempDir := t.TempDir()
	validArrayJSON := `[{"id":"res-utf16","name":"UTF16 Test","canonical_url":"https://example.com/u16","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag1"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`

	// 1. UTF-16 Little Endian
	utf16leBytes := encodeUTF16LE(validArrayJSON)
	lePath := filepath.Join(tempDir, "catalog_utf16_le.json")
	if err := os.WriteFile(lePath, utf16leBytes, 0644); err != nil {
		t.Fatalf("failed to write UTF-16 LE file: %v", err)
	}

	_, err := LoadResourceCatalog(lePath)
	t.Logf("Empirical result LoadResourceCatalog(UTF-16 LE): err=%v", err)
	if err == nil {
		t.Errorf("expected error or parse failure for UTF-16 LE since Go JSON expects UTF-8, but got nil")
	}

	// 2. UTF-16 Big Endian
	utf16beBytes := encodeUTF16BE(validArrayJSON)
	bePath := filepath.Join(tempDir, "catalog_utf16_be.json")
	if err := os.WriteFile(bePath, utf16beBytes, 0644); err != nil {
		t.Fatalf("failed to write UTF-16 BE file: %v", err)
	}

	_, err = LoadResourceCatalog(bePath)
	t.Logf("Empirical result LoadResourceCatalog(UTF-16 BE): err=%v", err)
	if err == nil {
		t.Errorf("expected error or parse failure for UTF-16 BE since Go JSON expects UTF-8, but got nil")
	}
}

// ============================================================================
// CHALLENGE SUITE 2: Quarantine Evasion & Symlink Adversarial Tests
// ============================================================================

func TestChallenger_Quarantine_RelativePathVariations(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"parent directory traversal", "../skills_library"},
		{"multi-level traversal", "../../skills_library/some-skill"},
		{"current dir traversal", "./foo/../skills_library/sub"},
		{"deep traversal", "a/b/c/../../../../skills_library"},
		{"traversal with backslashes", "..\\..\\skills_library\\sub"},
		{"mixed slashes traversal", "..\\../skills_library/foo\\bar"},
		{"mixed dots and slashes", "././skills_library/./"},
		{"hyphenated variant traversal", "../skills-library/tool"},
		{"curated quarantine traversal", "../curated_catalog/quarantine/sub"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckQuarantineBoundary(tc.path)
			if err == nil {
				t.Fatalf("SECURITY VIOLATION: CheckQuarantineBoundary failed to reject path '%s'", tc.path)
			}
			if !errors.Is(err, ErrQuarantinedPath) {
				t.Errorf("expected ErrQuarantinedPath wrap, got: %v", err)
			}
		})
	}
}

func TestChallenger_Quarantine_CaseVariations(t *testing.T) {
	caseVariations := []string{
		"SKILLS_LIBRARY",
		"Skills_Library",
		"sKiLlS_lIbRaRy",
		"sKiLlS-lIbRaRy",
		"SKILLS-LIBRARY",
		"cUrAtEd_CaTaLoG/qUaRaNtInE",
		"CURATED_CATALOG/QUARANTINE",
	}

	for _, p := range caseVariations {
		t.Run(p, func(t *testing.T) {
			fullPath := filepath.Join("C:", "Users", "mockuser", ".gemini", "config", p, "skill", "SKILL.md")
			err := CheckQuarantineBoundary(fullPath)
			if err == nil {
				t.Fatalf("SECURITY VIOLATION: Case variation '%s' bypassed quarantine!", fullPath)
			}
			if !errors.Is(err, ErrQuarantinedPath) {
				t.Errorf("expected ErrQuarantinedPath wrap for '%s', got: %v", fullPath, err)
			}
		})
	}
}

func TestChallenger_Quarantine_URLFieldQuarantine(t *testing.T) {
	tempDir := t.TempDir()

	testResources := []struct {
		name        string
		jsonContent string
		expectErr   bool
	}{
		{
			name: "quarantined CanonicalURL",
			jsonContent: `[{"id":"res-q-1","name":"Q1","canonical_url":"file:///C:/Users/mockuser/.gemini/config/skills_library/foo","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectErr:   true,
		},
		{
			name: "quarantined SourceRepository",
			jsonContent: `[{"id":"res-q-2","name":"Q2","canonical_url":"https://example.com/ok","source_repository":"https://github.com/org/skills_library.git","source_type":"github_repository","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag"],"acquisition_method":"git","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectErr:   true,
		},
		{
			name: "quarantined DocumentationURL",
			jsonContent: `[{"id":"res-q-3","name":"Q3","canonical_url":"https://example.com/ok","documentation_url":"https://docs.example.com/curated_catalog/quarantine/test","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectErr:   true,
		},
	}

	for _, tc := range testResources {
		t.Run(tc.name, func(t *testing.T) {
			p := filepath.Join(tempDir, fmt.Sprintf("%s.json", tc.name))
			if err := os.WriteFile(p, []byte(tc.jsonContent), 0644); err != nil {
				t.Fatalf("failed to write test json: %v", err)
			}

			_, err := LoadResourceCatalog(p)
			if tc.expectErr && err == nil {
				t.Fatalf("SECURITY VIOLATION: LoadResourceCatalog did not reject quarantined field in %s", tc.name)
			}
			if tc.expectErr && !errors.Is(err, ErrQuarantinedPath) {
				t.Errorf("expected ErrQuarantinedPath error, got: %v", err)
			}
		})
	}
}

func TestChallenger_Quarantine_SymlinkOrJunctionAnalysis(t *testing.T) {
	tempDir := t.TempDir()

	realQuarantineDir := filepath.Join(tempDir, "skills_library")
	if err := os.MkdirAll(realQuarantineDir, 0755); err != nil {
		t.Fatalf("failed to create real quarantine dir: %v", err)
	}

	quarantinedFile := filepath.Join(realQuarantineDir, "resources.json")
	payload := `[{"id":"quarantined-res","name":"Infected","canonical_url":"https://example.com","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["tag"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`
	if err := os.WriteFile(quarantinedFile, []byte(payload), 0644); err != nil {
		t.Fatalf("failed to write quarantined file: %v", err)
	}

	// Create an innocent junction path
	junctionPath := filepath.Join(tempDir, "innocent_junction")
	fileThroughJunction := filepath.Join(junctionPath, "resources.json")

	// Create NTFS junction using cmd /c mklink /J
	cmd := exec.Command("cmd", "/c", "mklink", "/J", junctionPath, realQuarantineDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("mklink /J not available (%v: %s), skipping live junction test", err, string(output))
		return
	}
	t.Logf("Created NTFS junction: %s -> %s (%s)", junctionPath, realQuarantineDir, strings.TrimSpace(string(output)))

	// 1. Test CheckQuarantineBoundary on junction path
	boundaryErr := CheckQuarantineBoundary(fileThroughJunction)
	t.Logf("Empirical CheckQuarantineBoundary(%s) = %v", fileThroughJunction, boundaryErr)

	// 2. Test LoadResourceCatalog on junction path
	cat, loadErr := LoadResourceCatalog(fileThroughJunction)
	t.Logf("Empirical LoadResourceCatalog(%s) = cat:%v, err:%v", fileThroughJunction, cat != nil, loadErr)

	if loadErr == nil {
		t.Logf("FINDING (SECURITY EVASION): LoadResourceCatalog successfully loaded quarantined file via NTFS junction without returning ErrQuarantinedPath! (Cause: CheckQuarantineBoundary uses lexical substring search without filepath.EvalSymlinks resolution)")
	} else if errors.Is(loadErr, ErrQuarantinedPath) {
		t.Logf("SUCCESS: LoadResourceCatalog blocked junction to quarantine.")
	}
}

// ============================================================================
// CHALLENGE SUITE 3: Robustness & Adversarial Edge Cases
// ============================================================================

func TestChallenger_Robustness_MalformedInputs(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name        string
		content     string
		expectError bool
		errorSubstr string
	}{
		{
			name:        "empty catalog file",
			content:     "",
			expectError: true,
			errorSubstr: "is empty",
		},
		{
			name:        "whitespace only catalog",
			content:     "   \n\t  \r\n",
			expectError: true,
			errorSubstr: "is empty",
		},
		{
			name:        "object instead of array for catalog",
			content:     `{"id":"not-an-array"}`,
			expectError: true,
			errorSubstr: "expected JSON array root",
		},
		{
			name:        "primitive number instead of array",
			content:     `42`,
			expectError: true,
			errorSubstr: "expected JSON array root",
		},
		{
			name:        "primitive boolean",
			content:     `true`,
			expectError: true,
			errorSubstr: "expected JSON array root",
		},
		{
			name:        "truncated array JSON",
			content:     `[{"id":"res-1", "name": "incomplete"`,
			expectError: true,
			errorSubstr: "failed to parse",
		},
		{
			name:        "trailing comma in array",
			content:     `[{"id":"res-1","name":"N","canonical_url":"u","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"},]`,
			expectError: true,
			errorSubstr: "failed to parse",
		},
		{
			name:        "null item in array",
			content:     `[null, {"id":"res-ok","name":"N","canonical_url":"u","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}, null]`,
			expectError: false, // null items are skipped gracefully by discovery.go
		},
		{
			name:        "resource with empty ID",
			content:     `[{"id":"","name":"No ID","canonical_url":"u","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectError: true,
			errorSubstr: "resource id is required",
		},
		{
			name:        "resource with whitespace-only ID",
			content:     `[{"id":"   \t ","name":"Whitespace ID","canonical_url":"u","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectError: true,
			errorSubstr: "resource id is required",
		},
		{
			name:        "duplicate ID (exact case)",
			content:     `[{"id":"dup-1","name":"A","canonical_url":"u1","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"},{"id":"dup-1","name":"B","canonical_url":"u2","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectError: true,
			errorSubstr: "duplicate resource id",
		},
		{
			name:        "duplicate ID (case-insensitive clash)",
			content:     `[{"id":"case-id","name":"A","canonical_url":"u1","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"},{"id":"CASE-ID","name":"B","canonical_url":"u2","source_type":"npm_package","category":["FRONTEND"],"representation":"dependency","routing_tags":["t"],"acquisition_method":"npm","runtime_method":"project_scoped_install","status":"ACTIVE"}]`,
			expectError: true,
			errorSubstr: "duplicate resource id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(tempDir, fmt.Sprintf("%s.json", strings.ReplaceAll(tc.name, " ", "_")))
			if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			cat, err := LoadResourceCatalog(path)
			if tc.expectError {
				if err == nil {
					t.Fatalf("[%s] expected error, but LoadResourceCatalog succeeded", tc.name)
				}
				if tc.errorSubstr != "" && !strings.Contains(err.Error(), tc.errorSubstr) {
					t.Errorf("[%s] expected error containing %q, got %q", tc.name, tc.errorSubstr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("[%s] unexpected error: %v", tc.name, err)
				}
				if cat.Count() != 1 {
					t.Errorf("[%s] expected 1 valid resource, got %d", tc.name, cat.Count())
				}
			}
		})
	}
}

func TestChallenger_Robustness_DesignGraphMalformed(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name        string
		content     string
		expectError bool
	}{
		{"empty file", "", true},
		{"array instead of object", `[]`, true},
		{"corrupted JSON", `{"domains": {`, true},
		{"empty domains and capabilities object", `{"domains":{},"capabilities":{}}`, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(tempDir, fmt.Sprintf("%s.json", strings.ReplaceAll(tc.name, " ", "_")))
			if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			_, err := LoadDesignGraph(path)
			if tc.expectError && err == nil {
				t.Fatalf("[%s] expected error, but LoadDesignGraph succeeded", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Fatalf("[%s] unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestChallenger_Stress_LargeCatalog(t *testing.T) {
	// Stress test with 5,000 generated resources
	const resourceCount = 5000
	tempDir := t.TempDir()
	largeCatPath := filepath.Join(tempDir, "large_catalog.json")

	buf := new(bytes.Buffer)
	buf.WriteString("[\n")
	for i := 0; i < resourceCount; i++ {
		if i > 0 {
			buf.WriteString(",\n")
		}
		item := fmt.Sprintf(`  {
    "id": "stress-res-%05d",
    "name": "Stress Resource %d",
    "canonical_url": "https://example.com/res/%d",
    "source_type": "npm_package",
    "category": ["SPECIALIST", "CATEGORY_%d"],
    "representation": "dependency",
    "routing_tags": ["tag_%d", "common_tag", "category_%d"],
    "acquisition_method": "npm",
    "runtime_method": "project_scoped_install",
    "status": "ACTIVE"
  }`, i, i, i, i%10, i%20, i%10)
		buf.WriteString(item)
	}
	buf.WriteString("\n]")

	if err := os.WriteFile(largeCatPath, buf.Bytes(), 0644); err != nil {
		t.Fatalf("failed to write large catalog: %v", err)
	}

	start := time.Now()
	cat, err := LoadResourceCatalog(largeCatPath)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("failed to load 5,000-item catalog: %v", err)
	}
	if cat.Count() != resourceCount {
		t.Errorf("expected %d resources, got %d", resourceCount, cat.Count())
	}
	t.Logf("Empirical result: Loaded 5,000 resources in %v (%.2f ms)", duration, float64(duration.Microseconds())/1000.0)

	// Verify indexing lookups under scale
	res, ok := cat.FindByID("stress-res-02499")
	if !ok || res == nil {
		t.Fatalf("failed to find stress-res-02499")
	}

	commonItems := cat.FindByTag("common_tag")
	if len(commonItems) != resourceCount {
		t.Errorf("expected %d common_tag items, got %d", resourceCount, len(commonItems))
	}

	catItems := cat.FindByCategory("category_3")
	if len(catItems) != resourceCount/10 {
		t.Errorf("expected %d category_3 items, got %d", resourceCount/10, len(catItems))
	}
}

// ============================================================================
// CHALLENGE SUITE 4: Slice Backing Array Mutation & Concurrency Stress
// ============================================================================

func TestChallenger_SliceBackingArrayMutation_Investigation(t *testing.T) {
	// Adversarial test for line 599 in discovery.go:
	// allPhases := append(append(append(append(append(wf.Discovery, wf.Synthesis...), wf.Implementation...), wf.OptionalExtensions...), wf.ReverseEngineering...), wf.QA...)
	// If wf.Discovery has capacity > length, append will overwrite elements in its backing array!

	graph := &DesignResourceGraph{
		Domains: map[string][]string{
			"test-dom": {"res-1"},
		},
		Capabilities: map[string]*CapabilityPhaseDefinition{
			"test-cap": {
				Name:             "Test Cap",
				PrimaryArchetype: "premium-website",
				// Create a slice with excess capacity
				Discovery:          make([]string, 1, 10),
				Synthesis:          []string{"synth-1", "synth-2"},
				Implementation:     []string{"impl-1"},
				OptionalExtensions: []string{"opt-1"},
				ReverseEngineering: []string{"rev-1"},
				QA:                 []string{"qa-1"},
			},
		},
	}
	graph.Capabilities["test-cap"].Discovery[0] = "initial-discovery"

	// Verify initial state
	if cap(graph.Capabilities["test-cap"].Discovery) != 10 {
		t.Fatalf("expected cap 10, got %d", cap(graph.Capabilities["test-cap"].Discovery))
	}

	// Call ResolveCapabilities with a matching domain tag
	_ = graph.ResolveCapabilities([]string{"test-dom"})

	// Check if Discovery's underlying array beyond length 1 was mutated:
	// In Go, wf.Discovery has len=1. Slicing beyond len:
	expandedSlice := graph.Capabilities["test-cap"].Discovery[:5]
	t.Logf("Discovery backing array inspection: %v", expandedSlice)

	// If append(wf.Discovery, ...) reused the backing array, expandedSlice[1] will be "synth-1"!
	if expandedSlice[1] == "synth-1" {
		t.Logf("FINDING (MUTATION RISK): ResolveCapabilities append modifies the backing array of wf.Discovery when cap > len! Backing array element [1]=%q", expandedSlice[1])
	} else {
		t.Logf("Backing array was not mutated (allocated new buffer).")
	}
}

func TestChallenger_Concurrency_LiveRegistriesUnderLoad(t *testing.T) {
	workflowRoot := filepath.Join("..", "..", "..")
	resourcesPath := filepath.Join(workflowRoot, "registries", "resources.json")
	graphPath := filepath.Join(workflowRoot, "registries", "design-resource-graph.json")

	cat, err := LoadResourceCatalog(resourcesPath)
	if err != nil {
		t.Fatalf("failed to load resources catalog: %v", err)
	}

	graph, err := LoadDesignGraph(graphPath)
	if err != nil {
		t.Fatalf("failed to load design graph: %v", err)
	}

	const goroutines = 50
	const iterationsPerGoroutine = 500

	var totalOps int64
	var errorCount int64

	knownIDs := []string{
		"gsap", "threejs", "awwwards-directory", "framer-motion", "tailwind-css",
		"semgrep", "semgrep-adapter", "godly", "godly-design", "react-bits",
		"lenis", "playwright", "lucide-icons",
	}

	knownTags := []string{
		"motion", "3d", "canvas", "typography", "haptics", "verification",
		"premium", "design-system", "unknown-tag-xyz",
	}

	var wg sync.WaitGroup
	startSignal := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			<-startSignal

			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			for j := 0; j < iterationsPerGoroutine; j++ {
				op := rng.Intn(6)
				switch op {
				case 0:
					// FindByID
					id := knownIDs[rng.Intn(len(knownIDs))]
					res, _ := cat.FindByID(id)
					if res != nil && res.ID == "" {
						atomic.AddInt64(&errorCount, 1)
					}
				case 1:
					// FindByTag / FindByTags
					tag := knownTags[rng.Intn(len(knownTags))]
					items := cat.FindByTag(tag)
					if len(items) > 0 && items[0] == nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 2:
					// ResolveCapabilities
					tag := knownTags[rng.Intn(len(knownTags))]
					routes := graph.ResolveCapabilities([]string{tag})
					if len(routes) > 0 && routes[0].CapabilityName == "" {
						atomic.AddInt64(&errorCount, 1)
					}
				case 3:
					// ResolveCapabilityRoute
					capName := "premium-website"
					if rng.Intn(2) == 0 {
						capName = "3d-portfolio"
					}
					route, ok := graph.ResolveCapabilityRoute(capName, []string{"motion"})
					if ok && route == nil {
						atomic.AddInt64(&errorCount, 1)
					}
				case 4:
					// Tags & All
					tags := cat.Tags()
					if len(tags) == 0 {
						atomic.AddInt64(&errorCount, 1)
					}
				case 5:
					// Concurrent Loaders
					if j%50 == 0 {
						c, e := LoadResourceCatalog(resourcesPath)
						if e != nil || c.Count() != 126 {
							atomic.AddInt64(&errorCount, 1)
						}
					}
				}
				atomic.AddInt64(&totalOps, 1)
			}
		}(i)
	}

	startTime := time.Now()
	close(startSignal)
	wg.Wait()
	duration := time.Since(startTime)

	ops := atomic.LoadInt64(&totalOps)
	errs := atomic.LoadInt64(&errorCount)

	t.Logf("Empirical Concurrency Results:")
	t.Logf("  Goroutines: %d", goroutines)
	t.Logf("  Total Operations: %d", ops)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Throughput: %.2f ops/sec", float64(ops)/duration.Seconds())
	t.Logf("  Errors/Inconsistencies: %d", errs)

	if errs > 0 {
		t.Fatalf("Concurrency stress test encountered %d errors/data corruption issues!", errs)
	}
}

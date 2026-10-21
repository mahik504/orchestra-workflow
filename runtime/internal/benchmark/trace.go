package benchmark

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/acquisition"
	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/engine"
	"github.com/mahik504/orchestra-workflow/runtime/internal/memory"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

func CopyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

type TraceInput struct {
	Catalog    *resources.ResourceCatalog
	Graph      *resources.DesignResourceGraph
	FixtureDir string
	Task       Task
	ModelLabel string
	WorkRoot   string
}

func RunTrace(ctx context.Context, in TraceInput) (Record, error) {
	start := time.Now()
	if !in.Task.Enabled {
		return Record{}, fmt.Errorf("%s is behind the expansion gate", in.Task.ID)
	}
	srcHash, err := HashTree(in.FixtureDir)
	if err != nil {
		return Record{}, err
	}
	if !TreeClean(in.FixtureDir) {
		return Record{}, fmt.Errorf("fixture %s is not clean", in.Task.ID)
	}
	base := filepath.Join(in.WorkRoot, "baseline")
	orch := filepath.Join(in.WorkRoot, "orchestra")
	if err := CopyTree(in.FixtureDir, base); err != nil {
		return Record{}, err
	}
	if err := CopyTree(in.FixtureDir, orch); err != nil {
		return Record{}, err
	}
	baseHash, err := HashTree(base)
	if err != nil {
		return Record{}, err
	}
	orchHash, err := HashTree(orch)
	if err != nil {
		return Record{}, err
	}
	if baseHash != orchHash || baseHash != srcHash {
		return Record{}, fmt.Errorf("starting snapshots differ")
	}
	task, acc, err := ReadTask(base)
	if err != nil {
		return Record{}, err
	}
	task2, acc2, err := ReadTask(orch)
	if err != nil {
		return Record{}, err
	}
	if task != task2 || strings.Join(acc, "|") != strings.Join(acc2, "|") {
		return Record{}, fmt.Errorf("task text differs")
	}
	if err := WritePacket(base, Packet{
		Orchestration: "off",
		Task:          task,
		Acceptance:    acc,
		FixtureID:     in.Task.ID,
		FixtureHash:   baseHash,
		TimeBudget:    "same",
		HumanPolicy:   "existing human gate only",
		ModelLabel:    in.ModelLabel,
	}); err != nil {
		return Record{}, err
	}
	if err := BaselineIsolated(base); err != nil {
		return Record{}, err
	}

	if err := declareFixture(orch, in.Task.ID); err != nil {
		return Record{}, err
	}
	brain := memory.LoadDurable(nil, "missing", in.Task.ID, nil)
	req := &engine.TaskRequest{
		ID:            in.Task.ID,
		RawRequest:    task,
		WorkspaceRoot: orch,
		Type:          in.Task.Type,
		MaxIterations: 1,
		Parameters: map[string]any{
			"playwright_available": false,
		},
	}
	if in.ModelLabel != "" {
		req.Parameters["discovered_models"] = []string{in.ModelLabel}
	}
	pipeline := engine.NewDesignPipeline(in.Catalog, in.Graph, nil, engine.NewMockVisualVerifier())
	res, err := pipeline.Execute(ctx, req)
	if err != nil {
		return Record{}, err
	}
	rec := Record{
		BenchmarkID:   "3g-pilot-" + in.Task.ID,
		TaskClass:     in.Task.Class,
		FixtureID:     in.Task.ID,
		FixtureHash:   srcHash,
		Clean:         true,
		ModelLabel:    in.ModelLabel,
		Mode:          "on",
		Phase:         PhasePilot,
		Trial:         1,
		Start:         start.UTC().Format(time.RFC3339),
		End:           time.Now().UTC().Format(time.RFC3339),
		Status:        "trace_only",
		Claim:         ClaimNone,
		WallMS:        time.Since(start).Milliseconds(),
		BrainFacts:    len(brain.Selected),
		TaskExecution: TaskCostUnavailable,
		Tokens:        TokensUnavailable,
		Limitations:   []string{"trace_only: pipeline does not write the product"},
	}
	if res != nil {
		rec.WallMS = res.TotalDuration.Milliseconds()
		for _, sr := range res.StageResults {
			if sr.StageName == engine.StageResearch && sr.Status != engine.StatusSkipped {
				rec.ResearchCalls++
			}
		}
		rec.ScreenshotEvals = len(res.Screenshots)
		rec.BrowserCalls = len(res.Screenshots)
	}
	m, err := mission.ReadMissionFile(orch, in.Task.ID)
	if err != nil {
		return rec, err
	}
	rec.Intensity = string(m.Intensity)
	if m.Budget.VisualReview != mission.VisualReviewNone {
		rec.BrowserCalls++
	}
	if m.Budget.IndependentJudge != mission.JudgeOff {
		rec.JudgeCalls++
	}
	rawM, err := os.ReadFile(filepath.Join(mission.RunDir(orch, in.Task.ID), "mission.json"))
	if err != nil {
		return rec, err
	}
	rec.MissionBytes = len(rawM)
	if in.ModelLabel != "" && strings.Contains(strings.ToLower(string(rawM)), strings.ToLower(in.ModelLabel)) {
		return rec, fmt.Errorf("model label entered the mission")
	}
	dec, err := capability.ReadDecision(orch, in.Task.ID)
	if err != nil {
		return rec, err
	}
	rawC, err := os.ReadFile(filepath.Join(mission.RunDir(orch, in.Task.ID), "capability.json"))
	if err != nil {
		return rec, err
	}
	rec.CapabilityBytes = len(rawC)
	if in.ModelLabel != "" && strings.Contains(strings.ToLower(string(rawC)), strings.ToLower(in.ModelLabel)) {
		return rec, fmt.Errorf("model label entered the capability decision")
	}
	for _, it := range dec.Selected {
		rec.CapabilityIDs = append(rec.CapabilityIDs, it.ID)
		switch it.ID {
		case capability.IDDesignWorker:
			rec.Directions++
		case capability.IDJudgeIndependent:
			rec.JudgeCalls++
		case capability.IDBrowserPlaywright:
			rec.BrowserCalls++
		case capability.IDResearchWorker:
			rec.ResearchCalls++
		}
	}
	if st, err := mission.ReadRunState(orch, in.Task.ID); err == nil && st.CapsuleScope != nil {
		rec.CapsuleBytes = st.CapsuleScope.CapsuleBytes
	}
	rec.OrchestrationBytes = rec.MissionBytes + rec.CapabilityBytes + rec.CapsuleBytes
	if err := WritePacket(orch, Packet{
		Orchestration:  "on",
		Task:           task,
		Acceptance:     acc,
		FixtureID:      in.Task.ID,
		FixtureHash:    srcHash,
		TimeBudget:     "same",
		HumanPolicy:    "existing human gate only",
		ModelLabel:     in.ModelLabel,
		MissionPath:    filepath.Join(mission.RunDir(orch, in.Task.ID), "mission.json"),
		CapabilityPath: filepath.Join(mission.RunDir(orch, in.Task.ID), "capability.json"),
		Intensity:      rec.Intensity,
	}); err != nil {
		return rec, err
	}
	after, err := HashTree(base)
	if err != nil {
		return rec, err
	}
	pkt, err := os.ReadFile(filepath.Join(base, "packet.json"))
	if err != nil {
		return rec, err
	}
	if after == baseHash {
		return rec, fmt.Errorf("baseline hash ignored the packet")
	}
	_ = pkt
	bare, err := hashWithout(base, "packet.json")
	if err != nil {
		return rec, err
	}
	if bare != srcHash {
		return rec, fmt.Errorf("baseline fixture changed")
	}
	return rec, nil
}

func declareFixture(workspace, taskID string) error {
	store, err := acquisition.NewProvenanceStore(workspace)
	if err != nil {
		return err
	}
	return filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".orchestra" {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() == "packet.json" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(raw)
		hash := hex.EncodeToString(sum[:])
		rel, err := filepath.Rel(workspace, path)
		if err != nil {
			return err
		}
		return store.Record(acquisition.ProvenanceEntry{
			ResourceID:          "fixture-" + strings.ReplaceAll(filepath.ToSlash(rel), "/", "-"),
			AcquisitionMethod:   "fixture",
			SourceURL:           "file://fixture/" + filepath.ToSlash(rel),
			VersionOrSHA:        hash[:12],
			SHA256Hash:          hash,
			InstalledPath:       path,
			JustificationTaskID: taskID,
		})
	})
}

func hashWithout(root, skip string) (string, error) {
	tmp, err := os.MkdirTemp("", "bench-hash-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	if err := CopyTree(root, tmp); err != nil {
		return "", err
	}
	if err := os.Remove(filepath.Join(tmp, skip)); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	return HashTree(tmp)
}

// PairedHostSmoke records one pilot cell. No host is invoked. The cell is not an effectiveness result.
func PairedHostSmoke(model string) Record {
	return Record{
		BenchmarkID:   "3g-pilot-pair",
		TaskClass:     "T1",
		FixtureID:     "t1-api",
		ModelLabel:    model,
		Mode:          "pair",
		Phase:         PhasePilot,
		Trial:         1,
		Status:        "unavailable",
		Claim:         ClaimNone,
		Tokens:        TokensUnavailable,
		TaskExecution: TaskCostUnavailable,
		Limitations: []string{
			"pilot smoke only",
			"no host executed a paired implementation",
			"not an effectiveness result",
		},
	}
}

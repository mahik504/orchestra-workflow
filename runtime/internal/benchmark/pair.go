package benchmark

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/capability"
	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
)

// PairResult is one pilot comparison. Claim none is not an effectiveness result.
type PairResult struct {
	BenchmarkID         string   `json:"benchmark_id"`
	TaskClass           string   `json:"task_class"`
	FixtureID           string   `json:"fixture_id"`
	FixtureHash         string   `json:"fixture_hash"`
	HostLabel           string   `json:"host_label,omitempty"`
	ModelLabel          string   `json:"model_label,omitempty"`
	Phase               string   `json:"phase"`
	Trial               int      `json:"trial"`
	Status              string   `json:"status"`
	Claim               string   `json:"claim"`
	Order               string   `json:"order"`
	Tokens              string   `json:"tokens"`
	BaselineIsolated    bool     `json:"baseline_isolated"`
	TaskHash            string   `json:"task_hash,omitempty"`
	AcceptanceHash      string   `json:"acceptance_hash,omitempty"`
	BaselinePromptHash  string   `json:"baseline_prompt_hash,omitempty"`
	OrchestraPromptHash string   `json:"orchestra_prompt_hash,omitempty"`
	PromptEquivalent    bool     `json:"prompt_equivalent"`
	GoVersion           string   `json:"go_version,omitempty"`
	AgentBase           string   `json:"agent_base,omitempty"`
	SharedProfile       string   `json:"shared_profile,omitempty"`
	ExternalClass       string   `json:"external_class,omitempty"`
	Baseline            Record   `json:"baseline"`
	Orchestra           Record   `json:"orchestra"`
	Limitations         []string `json:"limitations,omitempty"`
	WorktreeRemoved     bool     `json:"worktree_removed"`
}

// Guard is one protected path outside the fixture copy. A change invalidates the pair.
type Guard struct {
	Class string
	Path  string
}

const sharedProfileNote = "shared_cursor_profile: both arms use the same authenticated Cursor profile and inherited MCP and rules configuration. This does not isolate model-only effects."

type PairInput struct {
	Catalog        *resources.ResourceCatalog
	Graph          *resources.DesignResourceGraph
	FixtureDir     string
	Task           Task
	WorkRoot       string
	Executable     string
	Timeout        time.Duration
	HostLabel      string
	BaselineLabel  string
	OrchestraLabel string
	ResultPath     string
	SharedProfile  bool
	Guards         []Guard
}

// RunPair compares one baseline process with one Orchestra process.
// An empty executable is unavailable. This does not call a model by itself.
func RunPair(ctx context.Context, in PairInput) (res PairResult, err error) {
	res = PairResult{
		BenchmarkID: "3g1-" + in.Task.ID,
		TaskClass:   in.Task.Class,
		FixtureID:   in.Task.ID,
		HostLabel:   in.HostLabel,
		ModelLabel:  in.BaselineLabel,
		Phase:       PhasePilot,
		Trial:       1,
		Claim:       ClaimNone,
		Order:       "baseline_then_orchestra",
		Tokens:      TokensUnavailable,
		GoVersion:   goVersion(),
		Limitations: []string{"pilot observation", "claim none", "one trial", "two pairs are not a superiority claim"},
	}
	if agent := os.Getenv("ORCHESTRA_BENCH_AGENT"); agent != "" {
		res.AgentBase = filepath.Base(agent)
	}
	if in.SharedProfile {
		res.SharedProfile = sharedProfileNote
		res.Limitations = append(res.Limitations, sharedProfileNote)
	}
	defer func() {
		if in.WorkRoot != "" {
			if rmErr := os.RemoveAll(in.WorkRoot); rmErr == nil {
				res.WorktreeRemoved = true
			}
		}
		if in.ResultPath != "" {
			_ = writePair(in.ResultPath, res)
		}
	}()
	if in.BaselineLabel != in.OrchestraLabel {
		res.Limitations = append(res.Limitations, "model labels differ")
		res.Status = "invalid_pair"
		return res, nil
	}
	if strings.TrimSpace(in.Executable) == "" {
		res.Limitations = append(res.Limitations, "no host executable")
		res.Status = "unavailable"
		return res, nil
	}
	if in.HostLabel == "" {
		res.HostLabel = "cursor"
	}
	trace, err := RunTrace(ctx, TraceInput{
		Catalog:    in.Catalog,
		Graph:      in.Graph,
		FixtureDir: in.FixtureDir,
		Task:       in.Task,
		WorkRoot:   in.WorkRoot,
	})
	if err != nil {
		return res, err
	}
	res.FixtureHash = trace.FixtureHash
	base := filepath.Join(in.WorkRoot, "baseline")
	orch := filepath.Join(in.WorkRoot, "orchestra")
	capsuleRel, err := exportCapsule(orch, in.Task.ID)
	if err != nil {
		return res, err
	}
	task, acc, err := ReadTask(base)
	if err != nil {
		return res, err
	}
	basePrompt := baselinePrompt(task, acc)
	missionRel := filepath.ToSlash(filepath.Join(".orchestra", "runs", mission.SanitizeID(in.Task.ID), "mission.json"))
	capRel := filepath.ToSlash(filepath.Join(".orchestra", "runs", mission.SanitizeID(in.Task.ID), "capability.json"))
	orchPrompt := orchestraPrompt(task, acc, missionRel, capsuleRel, capRel)
	if err := os.WriteFile(filepath.Join(base, "prompt.txt"), []byte(basePrompt), 0644); err != nil {
		return res, err
	}
	if err := os.WriteFile(filepath.Join(orch, "prompt.txt"), []byte(orchPrompt), 0644); err != nil {
		return res, err
	}
	res.TaskHash = shortHash(task)
	res.AcceptanceHash = shortHash(strings.Join(acc, "\n"))
	res.BaselinePromptHash = shortHash(basePrompt)
	res.OrchestraPromptHash = shortHash(orchPrompt)
	res.PromptEquivalent = strings.HasPrefix(orchPrompt, basePrompt) && strings.Contains(orchPrompt, "Capsule file:") && !strings.Contains(basePrompt, "Capsule file:")
	if !res.PromptEquivalent {
		res.Limitations = append(res.Limitations, "task text or acceptance diverged")
		res.Status = "invalid_pair"
		return res, nil
	}
	if err := WritePacket(orch, Packet{
		Orchestration:  "on",
		Task:           task,
		Acceptance:     acc,
		FixtureID:      in.Task.ID,
		FixtureHash:    trace.FixtureHash,
		TimeBudget:     "same",
		HumanPolicy:    "existing human gate only",
		MissionPath:    missionRel,
		CapabilityPath: capRel,
		CapsulePath:    capsuleRel,
		Intensity:      trace.Intensity,
	}); err != nil {
		return res, err
	}
	leaked, err := labelInTrees(base, orch, in.Task.ID, in.BaselineLabel)
	if err != nil {
		return res, err
	}
	if leaked {
		res.Limitations = append(res.Limitations, "model label entered a host file")
		res.Status = "invalid_pair"
		return res, nil
	}
	policy, err := PolicyFor(in.Task.ID)
	if err != nil {
		return res, err
	}
	baseBefore, err := protectHashes(base, policy.Protect)
	if err != nil {
		return res, err
	}
	orchBefore, err := protectHashes(orch, policy.Protect)
	if err != nil {
		return res, err
	}
	if err := commitSnapshot(base); err != nil {
		return res, err
	}
	if err := commitSnapshot(orch); err != nil {
		return res, err
	}
	baseCache, orchCache, err := armCaches(in.WorkRoot)
	if err != nil {
		return res, err
	}
	guards := snapshotGuards(in.Guards)
	baseMark := harnessMark(base)
	baseRun := invoke(ctx, in, base, baseCache, in.Timeout)
	baseHarness := harnessMark(base) != baseMark
	baseEnd, err := HashTree(base)
	if err != nil {
		return res, err
	}
	baseRec, err := ScoreHost(base, policy, baseBefore, nil)
	if err != nil {
		return res, err
	}
	orchMark := harnessMark(orch)
	orchRun := invoke(ctx, in, orch, orchCache, in.Timeout)
	orchHarness := harnessMark(orch) != orchMark
	orchEnd, err := HashTree(orch)
	if err != nil {
		return res, err
	}
	orchRec, err := ScoreHost(orch, policy, orchBefore, nil)
	if err != nil {
		return res, err
	}
	baseRec.EndingHash = baseEnd
	orchRec.EndingHash = orchEnd
	baseRec.TaskHash = res.TaskHash
	orchRec.TaskHash = res.TaskHash
	baseRec.AcceptanceHash = res.AcceptanceHash
	orchRec.AcceptanceHash = res.AcceptanceHash
	baseRec.PromptHash = res.BaselinePromptHash
	orchRec.PromptHash = res.OrchestraPromptHash
	attachRun(&baseRec, baseRun)
	attachRun(&orchRec, orchRun)
	if baseRec.Artifact != nil {
		baseRec.Artifact.HarnessChanged = baseHarness
	}
	if orchRec.Artifact != nil {
		orchRec.Artifact.HarnessChanged = orchHarness
	}
	baseRec.HostEdited = hostEdited(baseRec)
	orchRec.HostEdited = hostEdited(orchRec)
	baseRec.TaskCompleted = taskCompleted(in.Task.ID, base, baseRec)
	orchRec.TaskCompleted = taskCompleted(in.Task.ID, orch, orchRec)
	baseRec.Mode = "off"
	orchRec.Mode = "on"
	baseRec.FixtureHash = trace.FixtureHash
	baseRec.FixtureID = in.Task.ID
	baseRec.TaskClass = in.Task.Class
	baseRec.HostLabel = res.HostLabel
	baseRec.ModelLabel = in.BaselineLabel
	copyTraceMetrics(&orchRec, trace)
	orchRec.HostLabel = res.HostLabel
	orchRec.ModelLabel = in.OrchestraLabel
	res.Baseline = baseRec
	res.Orchestra = orchRec
	if class := guardsChanged(guards); class != "" {
		res.ExternalClass = class
		res.BaselineIsolated = false
		res.Limitations = append(res.Limitations, "protected path changed: "+class)
		res.Status = "invalid_pair"
		return res, nil
	}
	baseLeak, err := treeHasLabel(base, in.BaselineLabel)
	if err != nil {
		return res, err
	}
	orchLeak, err := treeHasLabel(orch, in.OrchestraLabel)
	if err != nil {
		return res, err
	}
	if baseLeak || orchLeak {
		res.Limitations = append(res.Limitations, "model label entered a host file")
		res.Status = "invalid_pair"
		return res, nil
	}
	if err := BaselineIsolated(base); err != nil {
		res.BaselineIsolated = false
		res.Limitations = append(res.Limitations, "baseline gained orchestra artifacts")
		res.Status = "invalid_pair"
		return res, nil
	}
	res.BaselineIsolated = true
	switch {
	case baseRun.notFound || orchRun.notFound:
		res.Limitations = append(res.Limitations, "host executable was not started")
		res.Status = "unavailable"
		return res, nil
	case baseRun.timedOut || orchRun.timedOut:
		res.Limitations = append(res.Limitations, "host timed out; the other arm is not a win")
		res.Status = "timeout"
		return res, nil
	case baseRun.exit != 0 || orchRun.exit != 0:
		res.Limitations = append(res.Limitations, "host exited non-zero; the other arm is not a win")
		res.Status = "host_error"
		return res, nil
	default:
		res.Status = "compared"
		res.Claim = ClaimNone
		return res, nil
	}
}

type hostRun struct {
	exit        int
	timedOut    bool
	notFound    bool
	wall        time.Duration
	stdoutBytes int
	stderrBytes int
}

func invoke(ctx context.Context, in PairInput, workspace string, caches [2]string, timeout time.Duration) hostRun {
	prompt := filepath.Join(workspace, "prompt.txt")
	code, timedOut, wall, stdoutBytes, stderrBytes, err := runHost(ctx, in.Executable, workspace, prompt, caches[0], caches[1], in.BaselineLabel, timeout)
	run := hostRun{exit: code, timedOut: timedOut, wall: wall, stdoutBytes: stdoutBytes, stderrBytes: stderrBytes}
	var ee *exec.Error
	if errors.As(err, &ee) {
		run.notFound = true
	}
	return run
}

func attachRun(rec *Record, run hostRun) {
	if rec.Artifact == nil {
		rec.Artifact = &Artifact{}
	}
	rec.Artifact.HostExit = run.exit
	rec.Artifact.TimedOut = run.timedOut
	rec.Artifact.HostWallMS = run.wall.Milliseconds()
	rec.Artifact.StdoutBytes = run.stdoutBytes
	rec.Artifact.StderrBytes = run.stderrBytes
	rec.HostRetries = TokensUnavailable
	rec.Tokens = TokensUnavailable
	rec.TaskExecution = TaskCostUnavailable
}

func copyTraceMetrics(dst *Record, src Record) {
	dst.MissionBytes = src.MissionBytes
	dst.CapabilityBytes = src.CapabilityBytes
	dst.CapsuleBytes = src.CapsuleBytes
	dst.OrchestrationBytes = src.OrchestrationBytes
	dst.CapabilityIDs = append([]string{}, src.CapabilityIDs...)
	dst.Intensity = src.Intensity
	dst.BrainFacts = src.BrainFacts
	dst.Directions = src.Directions
	dst.BrowserCalls = src.BrowserCalls
	dst.ResearchCalls = src.ResearchCalls
	dst.JudgeCalls = src.JudgeCalls
	dst.ScreenshotEvals = src.ScreenshotEvals
	dst.WallMS = src.WallMS
	dst.FixtureHash = src.FixtureHash
	dst.FixtureID = src.FixtureID
	dst.TaskClass = src.TaskClass
	if len(src.Limitations) > 0 {
		dst.Limitations = append(dst.Limitations, src.Limitations...)
	}
}

func armCaches(root string) (base, orch [2]string, err error) {
	for _, arm := range []string{"baseline", "orchestra"} {
		cache := filepath.Join(root, ".cache-"+arm)
		mod := filepath.Join(root, ".mod-"+arm)
		if err = os.MkdirAll(cache, 0755); err != nil {
			return base, orch, err
		}
		if err = os.MkdirAll(mod, 0755); err != nil {
			return base, orch, err
		}
		if arm == "baseline" {
			base = [2]string{cache, mod}
		} else {
			orch = [2]string{cache, mod}
		}
	}
	return base, orch, nil
}

func exportCapsule(workspace, taskID string) (string, error) {
	m, err := mission.ReadMissionFile(workspace, taskID)
	if err != nil {
		return "", err
	}
	dec, err := capability.ReadDecision(workspace, taskID)
	if err != nil {
		return "", err
	}
	cap, _, err := mission.Compile(mission.CompileRequest{
		Mission:        m,
		Role:           mission.RoleBuilder,
		Skills:         dec.IDs(capability.CatSkill),
		Resources:      dec.IDs(capability.CatResource),
		ToolsAllowed:   dec.IDs(capability.CatMCPTool),
		ExpectedOutput: "scoped implementation for the mission objective",
		UnrelatedFiles: []string{"AGENTS.md", "README.md", "ARCHITECTURE.md"},
	})
	if err != nil {
		return "", err
	}
	raw, err := json.MarshalIndent(cap, "", "  ")
	if err != nil {
		return "", err
	}
	rel := filepath.ToSlash(filepath.Join(".orchestra", "runs", mission.SanitizeID(taskID), "capsule.json"))
	path := filepath.Join(workspace, filepath.FromSlash(rel))
	if err := os.WriteFile(path, raw, 0644); err != nil {
		return "", err
	}
	return rel, nil
}

func labelInTrees(base, orch, taskID, label string) (bool, error) {
	if strings.TrimSpace(label) == "" {
		return false, nil
	}
	baseFiles := []string{"prompt.txt", "task.txt", "acceptance.txt", "packet.json"}
	orchFiles := append([]string{}, baseFiles...)
	run := filepath.ToSlash(filepath.Join(".orchestra", "runs", mission.SanitizeID(taskID)))
	orchFiles = append(orchFiles, run+"/mission.json", run+"/capability.json", run+"/capsule.json")
	for _, rel := range baseFiles {
		raw, err := os.ReadFile(filepath.Join(base, filepath.FromSlash(rel)))
		if err != nil {
			return false, err
		}
		if textHasLabel(string(raw), label) {
			return true, nil
		}
	}
	for _, rel := range orchFiles {
		raw, err := os.ReadFile(filepath.Join(orch, filepath.FromSlash(rel)))
		if err != nil {
			return false, err
		}
		if textHasLabel(string(raw), label) {
			return true, nil
		}
	}
	return false, nil
}

func shortHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func goVersion() string {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func hostEdited(rec Record) bool {
	return rec.Artifact != nil && len(rec.Artifact.Changed) > 0
}

func taskCompleted(id, root string, rec Record) bool {
	if rec.Artifact == nil || rec.Artifact.TimedOut || rec.Artifact.HostExit != 0 {
		return false
	}
	tests, okT := checkNamed(rec, DimTests)
	scope, okS := checkNamed(rec, DimScope)
	acc, okA := checkNamed(rec, DimAcceptance)
	if !okT || !okS || !okA || !tests.Pass || !scope.Pass || !acc.Pass || rec.HardFail {
		return false
	}
	if id == "t1-api" {
		behavior, okB := checkNamed(rec, DimRequiredBehavior)
		if !okB || !behavior.Pass || !changedPath(rec, "sum_test.go") {
			return false
		}
	}
	if id == "t3-debug" {
		raw, err := os.ReadFile(filepath.Join(root, "sum.go"))
		if err != nil || strings.Contains(string(raw), "a + b + 1") {
			return false
		}
	}
	return true
}

func changedPath(rec Record, name string) bool {
	if rec.Artifact == nil {
		return false
	}
	for _, path := range rec.Artifact.Changed {
		if path == name {
			return true
		}
	}
	return false
}

func harnessMark(root string) string {
	prompt := "missing"
	if sum, err := hashFile(filepath.Join(root, "prompt.txt")); err == nil {
		prompt = sum
	}
	orch := "missing"
	dir := filepath.Join(root, ".orchestra")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		sum, err := HashTree(dir)
		if err != nil {
			orch = "hash-error"
		} else {
			orch = sum
		}
	}
	return prompt + "|" + orch
}

func treeHasLabel(root, label string) (bool, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return false, nil
	}
	want := strings.ToLower(label)
	found := false
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}
		name := info.Name()
		if info.IsDir() {
			switch name {
			case ".git", ".score-cache", ".score-mod":
				return filepath.SkipDir
			}
			return nil
		}
		switch name {
		case "host-stdout.txt", "host-stderr.txt", "benchmark.test.exe":
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(string(raw)), want) {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil && !errors.Is(err, filepath.SkipAll) {
		return false, err
	}
	return found, nil
}

type guardSnap struct {
	class string
	path  string
	sum   string
}

func snapshotGuards(gs []Guard) []guardSnap {
	out := make([]guardSnap, 0, len(gs))
	for _, g := range gs {
		out = append(out, guardSnap{class: g.Class, path: g.Path, sum: fingerprintPath(g.Path)})
	}
	return out
}

func guardsChanged(snaps []guardSnap) string {
	for _, s := range snaps {
		if fingerprintPath(s.path) != s.sum {
			return s.class
		}
	}
	return ""
}

func fingerprintPath(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "missing"
	}
	if info.IsDir() {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			raw, err := gitOut(path, "status", "--porcelain")
			if err != nil {
				return "git-error"
			}
			return shortHash(raw)
		}
		sum, err := HashTree(path)
		if err != nil {
			return "hash-error"
		}
		return sum
	}
	sum, err := hashFile(path)
	if err != nil {
		return "hash-error"
	}
	return sum
}

func writePair(path string, res PairResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0644)
}

func checkNamed(rec Record, name string) (Check, bool) {
	for _, c := range rec.Deterministic {
		if c.Name == name {
			return c, true
		}
	}
	return Check{}, false
}

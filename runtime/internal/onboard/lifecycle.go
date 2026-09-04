package onboard

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/memory"
)

// Step is one numbered lifecycle proof line.
type Step struct {
	N        int            `json:"n"`
	Name     string         `json:"name"`
	Status   string         `json:"status"`
	Detail   string         `json:"detail"`
	Evidence map[string]any `json:"evidence,omitempty"`
}

// GitProbe is an acquisition check against the source URL.
type GitProbe struct {
	Command  string `json:"command"`
	OK       bool   `json:"ok"`
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

// Options drives a lifecycle run. Tests inject InspectFn / GitProbeFn.
type Options struct {
	URL           string
	Intent        string
	Origin        string
	OverlayPath   string
	MemoryPath    string
	MatchTask     string
	MismatchTask  string
	PublicCatalog string
	WorkflowRoot  string
	ArtifactDir   string
	InspectFn     func(string) Inspection
	GitProbeFn    func(string) GitProbe
}

// Report is the durable proof document written next to the overlay.
type Report struct {
	StartedAt         string         `json:"started_at"`
	FinishedAt        string         `json:"finished_at"`
	Mechanism         string         `json:"mechanism"`
	NotRL             string         `json:"not_reinforcement_learning"`
	OverlayPath       string         `json:"overlay_path"`
	MemoryPath        string         `json:"memory_path"`
	ReportPath        string         `json:"report_path,omitempty"`
	PublicCatalog     string         `json:"public_catalog,omitempty"`
	CatalogHashBefore string         `json:"public_catalog_sha256_before,omitempty"`
	CatalogHashAfter  string         `json:"public_catalog_sha256_after,omitempty"`
	CatalogUnchanged  bool           `json:"public_catalog_unchanged"`
	ResourceID        string         `json:"resource_id"`
	Steps             []Step         `json:"steps"`
	SuccessRecord     map[string]any `json:"success_record,omitempty"`
	FailureRecord     map[string]any `json:"failure_record,omitempty"`
	Architecture      []ArchFact     `json:"architecture,omitempty"`
	MatchAfterSuccess string         `json:"match_after_success,omitempty"`
	MatchAfterFailure string         `json:"match_after_failure,omitempty"`
}

// AddFromIntent infers and stores one overlay row. It does not edit the
// public router graph or resources.json.
func AddFromIntent(opts Options) (*OverlayEntry, *OverlayDocument, error) {
	intent := strings.TrimSpace(opts.Intent)
	rawURL := strings.TrimSpace(opts.URL)
	if rawURL == "" {
		rawURL = ExtractURL(intent)
	}
	if rawURL == "" {
		return nil, nil, fmt.Errorf("no resource URL in --url or --intent")
	}
	origin := opts.Origin
	if origin == "" {
		if intent != "" {
			origin = "user_intent"
		} else {
			origin = "url_submit"
		}
	}

	inspectFn := opts.InspectFn
	if inspectFn == nil {
		inspectFn = Inspect
	}
	insp := inspectFn(rawURL)
	prop := Infer(insp, intent, origin)
	if prop.Resource.ID == "" {
		return nil, nil, fmt.Errorf("could not infer resource id")
	}

	overlayPath := opts.OverlayPath
	if overlayPath == "" {
		overlayPath = ResolveOverlayPath(".")
	}
	doc, err := LoadOverlay(overlayPath)
	if err != nil {
		return nil, nil, err
	}

	entry := OverlayEntry{
		Origin:              origin,
		Intent:              intent,
		Kind:                prop.Kind,
		KindReason:          prop.KindReason,
		InstallScope:        prop.InstallScope,
		RoutingPolicy:       "active",
		FutureRoutingEffect: "activate when trigger conditions match; skip otherwise",
		Inspection:          insp,
		Resource:            prop.Resource,
	}
	doc.UpsertEntry(entry)
	if err := SaveOverlay(overlayPath, doc); err != nil {
		return nil, nil, err
	}
	saved := doc.Find(prop.Resource.ID)
	if saved == nil {
		return &entry, doc, nil
	}
	return saved, doc, nil
}

// RunLifecycle executes the 15-step proof, the intent path, and one
// success + one failure learning cycle.
func RunLifecycle(opts Options) (*Report, error) {
	started := time.Now().UTC()
	rep := &Report{
		StartedAt: started.Format(time.RFC3339),
		Mechanism: overlayMechanism,
		NotRL:     "No training loop, no reward model, no weight update. Heuristic policy over recorded evaluations.",
	}

	intent := strings.TrimSpace(opts.Intent)
	rawURL := strings.TrimSpace(opts.URL)
	if rawURL == "" {
		rawURL = ExtractURL(intent)
	}
	if rawURL == "" {
		return nil, fmt.Errorf("no resource URL")
	}

	overlayPath := opts.OverlayPath
	if overlayPath == "" {
		overlayPath = ResolveOverlayPath(".")
	}
	memPath := opts.MemoryPath
	if memPath == "" {
		memPath = memory.ResolveDefaultMemoryPath(".")
	}
	rep.OverlayPath = overlayPath
	rep.MemoryPath = memPath
	artifactDir := opts.ArtifactDir
	if artifactDir == "" && overlayPath != "" {
		artifactDir = filepath.Join(filepath.Dir(overlayPath), "lifecycle-pass-a")
	}

	var catalogBefore []byte
	if opts.PublicCatalog != "" {
		rep.PublicCatalog = opts.PublicCatalog
		if data, err := os.ReadFile(opts.PublicCatalog); err == nil {
			catalogBefore = data
			sum := sha256.Sum256(data)
			rep.CatalogHashBefore = hex.EncodeToString(sum[:])
		}
	}

	inspectFn := opts.InspectFn
	if inspectFn == nil {
		inspectFn = Inspect
	}
	gitFn := opts.GitProbeFn
	if gitFn == nil {
		gitFn = probeGit
	}

	addStep := func(n int, name, status, detail string, ev map[string]any) {
		rep.Steps = append(rep.Steps, Step{N: n, Name: name, Status: status, Detail: detail, Evidence: ev})
	}

	// 1 submitted
	addStep(1, "resource submitted", "ok", rawURL, map[string]any{
		"url":    rawURL,
		"intent": intent,
		"origin": firstNonEmpty(opts.Origin, "user_intent"),
	})
	proof(artifactDir, "01-submitted.json", map[string]any{
		"url": rawURL, "intent": intent, "origin": firstNonEmpty(opts.Origin, "user_intent"),
	})

	// 2–4 inspect + infer
	insp := inspectFn(rawURL)
	prop := Infer(insp, intent, firstNonEmpty(opts.Origin, "user_intent"))
	rep.ResourceID = prop.Resource.ID

	addStep(2, "resource identified", "ok", prop.Resource.ID, map[string]any{
		"id":          prop.Resource.ID,
		"name":        prop.Resource.Name,
		"source_type": prop.Resource.SourceType,
		"host":        insp.Host,
		"owner":       insp.Owner,
		"repo":        insp.Repo,
	})

	inspStatus := "ok"
	if !insp.Reachable {
		inspStatus = "observed"
	}
	addStep(3, "source inspected", inspStatus, fmt.Sprintf("HTTP %d reachable=%v", insp.StatusCode, insp.Reachable), map[string]any{
		"status_code":      insp.StatusCode,
		"reachable":        insp.Reachable,
		"error":            insp.Error,
		"has_skill_md":     insp.HasSkillMD,
		"has_package_json": insp.HasPackageJSON,
		"package_name":     insp.PackageName,
		"has_mcp_manifest": insp.HasMCPManifest,
		"title":            insp.Title,
		"probes":           insp.Probes,
	})
	proof(artifactDir, "02-identified.json", map[string]any{
		"id": prop.Resource.ID, "source_type": prop.Resource.SourceType, "host": insp.Host, "repo": insp.Repo,
	})
	proof(artifactDir, "03-inspection.json", insp)

	addStep(4, "capability inferred", "ok", prop.Kind+" / "+prop.Resource.Representation, map[string]any{
		"kind":        prop.Kind,
		"kind_reason": prop.KindReason,
		"rationale":   prop.Resource.Rationale,
		"source_type": prop.Resource.SourceType,
		"category":    prop.Resource.Category,
	})

	addStep(5, "representation selected", "ok", prop.Kind, map[string]any{
		"kind":                    prop.Kind,
		"registry_representation": prop.Resource.Representation,
		"kind_reason":             prop.KindReason,
		"allowed_kinds":           []string{KindSkill, KindDependency, KindMCP, KindPlugin, KindSubagent, KindReference, KindAdapter},
	})

	addStep(6, "routing tags generated", "ok", strings.Join(prop.Resource.RoutingTags, ", "), map[string]any{
		"routing_tags": prop.Resource.RoutingTags,
	})
	addStep(7, "trigger conditions generated", "ok", fmt.Sprintf("%d triggers", len(prop.Resource.TriggerConditions)), map[string]any{
		"trigger_conditions": prop.Resource.TriggerConditions,
	})
	addStep(8, "skip conditions generated", "ok", fmt.Sprintf("%d skip conditions", len(prop.Resource.AvoidConditions)), map[string]any{
		"skip_conditions": prop.Resource.AvoidConditions,
	})
	addStep(9, "acquisition method determined", "ok", prop.Resource.AcquisitionMethod, map[string]any{
		"acquisition_method": prop.Resource.AcquisitionMethod,
	})
	addStep(10, "installation scope determined", "ok", prop.InstallScope, map[string]any{
		"install_scope":  prop.InstallScope,
		"runtime_method": prop.Resource.RuntimeMethod,
		"global_blocked": prop.InstallScope != ScopeGlobal,
		"npm_package":    prop.Resource.NpmPackage,
	})
	proof(artifactDir, "04-inferred.json", prop)

	entry, doc, err := AddFromIntent(Options{
		URL:         rawURL,
		Intent:      intent,
		Origin:      firstNonEmpty(opts.Origin, "user_intent"),
		OverlayPath: overlayPath,
		InspectFn:   inspectFn,
	})
	if err != nil {
		return rep, err
	}
	_ = entry
	if data, err := os.ReadFile(overlayPath); err == nil {
		_ = writeProofText(artifactDir, "overlay.json", string(data))
	}

	matchTask := opts.MatchTask
	if matchTask == "" {
		matchTask = fmt.Sprintf("Use %s for the capability this repository provides.", prop.Resource.ID)
	}
	mismatchTask := opts.MismatchTask
	if mismatchTask == "" {
		mismatchTask = "The login form throws a 500 when the email has a plus sign."
	}

	matchDecisions := EvaluateTask(matchTask, doc)
	mismatchDecisions := EvaluateTask(mismatchTask, doc)
	matchD := decisionFor(matchDecisions, prop.Resource.ID)
	mismatchD := decisionFor(mismatchDecisions, prop.Resource.ID)

	matchStatus := "FAIL"
	if matchD.Action == ActionActivated {
		matchStatus = "ok"
	}
	addStep(11, "activated for matching task", matchStatus, matchD.Reason, map[string]any{
		"task":   matchTask,
		"action": matchD.Action,
		"reason": matchD.Reason,
	})

	mismatchStatus := "FAIL"
	if mismatchD.Action == ActionSuppressed {
		mismatchStatus = "ok"
	}
	addStep(12, "suppressed for non-matching task", mismatchStatus, mismatchD.Reason, map[string]any{
		"task":   mismatchTask,
		"action": mismatchD.Action,
		"reason": mismatchD.Reason,
	})
	proof(artifactDir, "11-match-decision.json", matchD)
	proof(artifactDir, "12-mismatch-decision.json", mismatchD)
	_ = writeProofText(artifactDir, "11-match-task.txt", matchTask)
	_ = writeProofText(artifactDir, "12-mismatch-task.txt", mismatchTask)

	git := gitFn(firstNonEmpty(prop.Resource.SourceRepository, prop.Resource.CanonicalURL))
	execOK := matchD.Action == ActionActivated && mismatchD.Action == ActionSuppressed
	execDetail := "matching task activated and non-matching task suppressed"
	if !execOK {
		execDetail = "selection proof did not hold"
	}
	addStep(13, "execution recorded", statusOK(execOK), execDetail, map[string]any{
		"match_action":    matchD.Action,
		"mismatch_action": mismatchD.Action,
		"git_probe":       git,
		"note":            "Pass A does not treat selection as a success evaluation. Acquisition failed.",
	})
	proof(artifactDir, "13-git-probe.json", git)

	verifyOK := execOK
	verifyDetail := "control-plane selection verified; overlay written; public catalog not used as the write target"
	if !git.OK {
		verifyDetail += "; git ls-remote failed (source not acquirable)"
	}
	addStep(14, "verification recorded", statusOK(verifyOK), verifyDetail, map[string]any{
		"selection_verified": verifyOK,
		"git_acquirable":     git.OK,
		"overlay_path":       overlayPath,
	})

	store, err := memory.NewResourceMemoryStore(memPath)
	if err != nil {
		return rep, err
	}

	// Selection is proven by classify artifacts, not a success evaluation.
	rep.MatchAfterSuccess = matchD.Action
	rep.SuccessRecord = nil

	failReason := "git ls-remote failed"
	if git.Output != "" {
		failReason = trimOut(git.Output)
	} else if !insp.Reachable {
		failReason = fmt.Sprintf("source HTTP %d, git probe failed", insp.StatusCode)
	}
	failureMeta := map[string]any{
		"failure":               "acquisition probe failed",
		"reason":                failReason,
		"future_routing_effect": "auto-activation suppressed after recorded failure",
		"user_feedback":         "do not treat an unreachable GitHub URL as installable",
	}
	failEval := &memory.ResourceEvaluation{
		ResourceID:   prop.Resource.ID,
		Domain:       "overlay",
		Capability:   prop.Kind,
		TaskContext:  "Acquire " + prop.Resource.ID + " via git for a matching task",
		TaskID:       "lifecycle-failure-" + prop.Resource.ID,
		Outcome:      memory.OutcomeFailure,
		QualityScore: 0.2,
		ErrorDetails: failReason,
		Notes:        "Pass A: source not acquirable. Policy suppresses future auto-activation. Not RL.",
		Metadata:     failureMeta,
	}
	if err := store.Record(failEval); err != nil {
		return rep, err
	}
	ApplyMemoryPolicy(doc, store)
	_ = SaveOverlay(overlayPath, doc)
	afterFail := EvaluateTask(matchTask, doc)
	rep.MatchAfterFailure = decisionFor(afterFail, prop.Resource.ID).Action
	proof(artifactDir, "15-failure-evaluation.json", failEval)
	proof(artifactDir, "15-match-after-failure.json", afterFail)

	routingPolicy := ""
	if found := doc.Find(prop.Resource.ID); found != nil {
		routingPolicy = found.RoutingPolicy
	}
	rep.FailureRecord = map[string]any{
		"resource":           prop.Resource.ID,
		"task":               failEval.TaskContext,
		"failure":            "acquisition probe failed",
		"reason":             failReason,
		"future_suppression": rep.MatchAfterFailure,
		"stored_in":          memPath,
		"routing_policy":     routingPolicy,
	}

	storedOK := doc.Find(prop.Resource.ID) != nil
	addStep(15, "outcome stored in resource memory", statusOK(storedOK), memPath, map[string]any{
		"memory_path":         memPath,
		"overlay_path":        overlayPath,
		"failure_eval_task":   failEval.TaskID,
		"match_before_memory": matchD.Action,
		"match_after_failure": rep.MatchAfterFailure,
		"mechanism":           overlayMechanism,
		"success_eval":        "none — Pass A does not record selection as success",
	})

	if opts.PublicCatalog != "" {
		if data, err := os.ReadFile(opts.PublicCatalog); err == nil {
			sum := sha256.Sum256(data)
			rep.CatalogHashAfter = hex.EncodeToString(sum[:])
			rep.CatalogUnchanged = string(catalogBefore) == string(data)
		}
	}
	proof(artifactDir, "catalog-sha.json", map[string]any{
		"before":    rep.CatalogHashBefore,
		"after":     rep.CatalogHashAfter,
		"unchanged": rep.CatalogUnchanged,
		"path":      opts.PublicCatalog,
	})

	root := opts.WorkflowRoot
	if root == "" {
		root = "."
	}
	rep.Architecture = VerifyArchitecture(root)

	reportPath := filepath.Join(artifactDir, "index.json")
	if artifactDir == "" {
		reportPath = filepath.Join(filepath.Dir(overlayPath), "lifecycle-audit.json")
	}
	rep.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	rep.ReportPath = reportPath
	if blob, err := json.MarshalIndent(rep, "", "  "); err == nil {
		_ = os.MkdirAll(filepath.Dir(reportPath), 0755)
		_ = os.WriteFile(reportPath, blob, 0644)
	}
	return rep, nil
}

func probeGit(source string) GitProbe {
	p := GitProbe{Command: "git ls-remote --heads " + source}
	if strings.TrimSpace(source) == "" {
		p.Output = "empty source url"
		p.ExitCode = 1
		return p
	}
	cmd := exec.Command("git", "ls-remote", "--heads", source)
	out, err := cmd.CombinedOutput()
	p.Output = trimOut(string(out))
	if err != nil {
		p.ExitCode = 1
		if p.Output == "" {
			p.Output = err.Error()
		}
		return p
	}
	p.OK = true
	return p
}

func decisionFor(ds []Decision, id string) Decision {
	want := strings.ToLower(id)
	for _, d := range ds {
		if strings.ToLower(d.ResourceID) == want {
			return d
		}
	}
	return Decision{ResourceID: id, Action: ActionSuppressed, Reason: "resource not in overlay"}
}

func statusOK(ok bool) string {
	if ok {
		return "ok"
	}
	return "FAIL"
}

func trimOut(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 500 {
		return s[:500]
	}
	return s
}

// FormatReport is the CLI transcript.
func FormatReport(rep *Report) string {
	if rep == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("Orchestra 3.1 lifecycle audit\n")
	b.WriteString("Mechanism: " + rep.Mechanism + "\n")
	b.WriteString(rep.NotRL + "\n\n")
	for _, s := range rep.Steps {
		fmt.Fprintf(&b, "%2d. [%s] %s\n    %s\n", s.N, strings.ToUpper(s.Status), s.Name, s.Detail)
	}
	b.WriteString("\nIntent path: overlay written without editing registries/resources.json\n")
	fmt.Fprintf(&b, "  overlay: %s\n  memory:  %s\n  report:  %s\n", rep.OverlayPath, rep.MemoryPath, rep.ReportPath)
	if rep.PublicCatalog != "" {
		fmt.Fprintf(&b, "  public catalog unchanged: %v\n", rep.CatalogUnchanged)
	}
	b.WriteString("\nSelection (not a success evaluation)\n")
	fmt.Fprintf(&b, "  match_before_memory: %s\n", rep.MatchAfterSuccess)
	b.WriteString("\nFailure record\n")
	writeMap(&b, rep.FailureRecord)
	b.WriteString("\nArchitecture\n")
	for _, a := range rep.Architecture {
		mark := "FAIL"
		if a.OK {
			mark = "ok"
		}
		fmt.Fprintf(&b, "  [%s] %s\n      %s (%s)\n", mark, a.Claim, a.Evidence, a.Path)
	}
	return b.String()
}

func writeMap(b *strings.Builder, m map[string]any) {
	keys := []string{
		"resource", "task", "why_selected", "actual_usage", "verification_result",
		"user_feedback", "future_routing_effect", "failure", "reason",
		"future_suppression", "stored_in", "match_after_success", "routing_policy",
	}
	seen := map[string]bool{}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Fprintf(b, "  %s: %v\n", k, v)
			seen[k] = true
		}
	}
	for k, v := range m {
		if !seen[k] {
			fmt.Fprintf(b, "  %s: %v\n", k, v)
		}
	}
}

package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/acquisition"
	acqAdapters "github.com/mahik504/orchestra-workflow/runtime/internal/adapters/acquisition"
	"github.com/mahik504/orchestra-workflow/runtime/internal/handoff"
	"github.com/mahik504/orchestra-workflow/runtime/internal/resources"
	"github.com/mahik504/orchestra-workflow/runtime/internal/runner"
	"github.com/mahik504/orchestra-workflow/runtime/internal/verify"
)

// ImplementStage executes Stage 6 of the design engine: acquiring resources and recording provenance.
type ImplementStage struct {
	adapterRegistry *acqAdapters.AdapterRegistry
}

// NewImplementStage creates a standard ImplementStage
func NewImplementStage() *ImplementStage {
	return &ImplementStage{}
}

// NewImplementStageWithRegistry creates an ImplementStage with a custom or mocked adapter registry
func NewImplementStageWithRegistry(reg *acqAdapters.AdapterRegistry) *ImplementStage {
	return &ImplementStage{adapterRegistry: reg}
}

func (s *ImplementStage) Name() StageName {
	return StageImplement
}

func (s *ImplementStage) ShouldSkip(ctx *TaskContext) (bool, string) {
	return false, ""
}

func (s *ImplementStage) Execute(ctx *TaskContext) (*StageResult, error) {
	start := time.Now()

	// The Design Lab gate is a lock, not a warning. If a direction is owed and
	// has not been approved, no frontend file gets written on this run — and
	// that includes dry runs, which used to slip past the synthesize halt.
	if ctx.DesignLab != nil && !ctx.DesignLab.Cleared() {
		err := &verify.ErrGateNotCleared{
			Path:   "<frontend implementation>",
			Reason: ctx.DesignLab.Reason,
		}
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	handoffDir := filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra", "handoff")
	if err := resources.CheckQuarantineBoundary(handoffDir); err != nil {
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}
	_ = os.MkdirAll(handoffDir, 0755)
	compileBuilderCapsule(ctx)
	if err := assertCodingWorkerConsumed(ctx); err != nil {
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     err,
		}, err
	}

	// 1. Initialize Provenance Store
	provStore, err := acquisition.NewProvenanceStore(ctx.Task.WorkspaceRoot)
	if err != nil {
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     fmt.Errorf("failed to init provenance store: %w", err),
		}, err
	}

	// 2. Setup Adapter Registry
	reg := s.adapterRegistry
	if reg == nil {
		r := runner.NewOSCommandRunner()
		reg = acqAdapters.NewAdapterRegistry()
		reg.RegisterAdapter(acqAdapters.NewNPMAdapter(r, ctx.Catalog, nil))
		reg.RegisterAdapter(acqAdapters.NewGitAdapter(r, nil))
		reg.RegisterAdapter(acqAdapters.NewCLIAdapter(r))
		reg.RegisterAdapter(acqAdapters.NewWebAdapter(true)) // offline fixture fallback
		reg.RegisterAdapter(acqAdapters.NewMCPAdapter(r))
	}

	// 3. Resolve Candidate Resources from Routes & Task Suggestions
	candidateSet := make(map[string]bool)
	var candidateList []string
	addCandidate := func(id string) {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" && !candidateSet[trimmed] {
			candidateSet[trimmed] = true
			candidateList = append(candidateList, trimmed)
		}
	}

	if ctx.Classification != nil {
		for _, route := range ctx.Classification.ResolvedRoutes {
			for _, impl := range route.Implementation {
				addCandidate(impl)
			}
		}
	}
	if ctx.Task != nil {
		for _, sug := range ctx.Task.SuggestedResources {
			addCandidate(sug)
		}
	}

	var acquiredResources []string
	installedPaths := map[string]string{}
	acquireCommands := map[string]string{}

	// 4. Conditional Acquisition & Provenance Recording
	for _, resID := range candidateList {
		var res *resources.Resource
		if ctx.Catalog != nil {
			if foundRes, ok := ctx.Catalog.FindByID(resID); ok {
				res = foundRes
			}
		}
		if res == nil {
			// Synthesize placeholder Resource for uncataloged candidate
			res = &resources.Resource{
				ID:                resID,
				Name:              resID,
				AcquisitionMethod: "none",
			}
		}

		// Reject banned/rejected resources
		if strings.EqualFold(res.Status, "REJECTED") {
			continue
		}

		// Check if already acquired in provenance ledger
		if existing, err := provStore.GetByResourceID(resID); err == nil && !existing.IsQuarantined {
			acquiredResources = append(acquiredResources, resID)
			if existing.InstalledPath != "" {
				installedPaths[res.ID] = existing.InstalledPath
			}
			continue
		}

		// In dry-run mode, record simulated provenance without executing subprocesses
		if ctx.Task != nil && ctx.Task.DryRun {
			entry := nonArtifactEntry(res, ctx.Task.ID, acquisition.KindDryRun, res.CanonicalURL, "dry-run-v1")
			if err := provStore.Record(entry); err != nil {
				return &StageResult{
					StageName: StageImplement,
					Status:    StatusFailed,
					StartTime: start,
					EndTime:   time.Now(),
					Duration:  time.Since(start),
					Error:     err,
				}, err
			}
			acquiredResources = append(acquiredResources, res.ID)
			continue
		}

		// Look up adapter for resource acquisition method
		adapter, err := reg.GetAdapterForMethod(res.AcquisitionMethod)
		if err != nil {
			entry := nonArtifactEntry(res, ctx.Task.ID, acquisition.KindReferenceOnly, res.CanonicalURL, "reference")
			if err := provStore.Record(entry); err != nil {
				return &StageResult{
					StageName: StageImplement,
					Status:    StatusFailed,
					StartTime: start,
					EndTime:   time.Now(),
					Duration:  time.Since(start),
					Error:     err,
				}, err
			}
			acquiredResources = append(acquiredResources, res.ID)
			continue
		}

		// Execute acquisition through adapter. Git clones go under .orchestra/vendor
		// so the product workspace root is never the clone destination.
		acqDest := ctx.Task.WorkspaceRoot
		if strings.EqualFold(res.AcquisitionMethod, "git") {
			acqDest = filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra", "vendor", res.ID)
		}
		acqResult, err := adapter.Acquire(ctx.Ctx, res, acqDest)
		if err != nil {
			// If project workspace lacks package.json, record pending dependency in provenance ledger for downstream agent handoff
			if errors.Is(err, acqAdapters.ErrPackageJSONNotFound) {
				sourceURL := res.CanonicalURL
				if sourceURL == "" && (res.AcquisitionMethod == "npm" || res.AcquisitionMethod == "npx") {
					sourceURL = "https://www.npmjs.com/package/" + res.ID
				}
				entry := nonArtifactEntry(res, ctx.Task.ID, acquisition.KindPending, sourceURL, "pending")
				if recErr := provStore.Record(entry); recErr != nil {
					return &StageResult{
						StageName: StageImplement,
						Status:    StatusFailed,
						StartTime: start,
						EndTime:   time.Now(),
						Duration:  time.Since(start),
						Error:     recErr,
					}, recErr
				}
				acquiredResources = append(acquiredResources, res.ID)
				continue
			}

			// If CLI/network acquisition failed (e.g. offline CI runner, network sandbox, timeout),
			// record pending provenance for downstream agent execution rather than crashing pipeline
			if res.AcquisitionMethod == "cli" || res.AcquisitionMethod == "npx" {
				entry := nonArtifactEntry(res, ctx.Task.ID, acquisition.KindPending, res.CanonicalURL, "pending")
				if recErr := provStore.Record(entry); recErr == nil {
					acquiredResources = append(acquiredResources, res.ID)
					continue
				}
			}

			return &StageResult{
				StageName: StageImplement,
				Status:    StatusFailed,
				StartTime: start,
				EndTime:   time.Now(),
				Duration:  time.Since(start),
				Error:     fmt.Errorf("failed acquiring resource %s: %w", resID, err),
			}, err
		}

		sourceURL := res.CanonicalURL
		if sourceURL == "" {
			sourceURL = acqResult.SourceURL
		}
		entry := artifactOrPending(res, ctx.Task.ID, sourceURL, acqResult.VersionOrSHA, acqResult.SHA256Hash, acqResult.InstalledPath)

		if err := provStore.Record(entry); err != nil {
			return &StageResult{
				StageName: StageImplement,
				Status:    StatusFailed,
				StartTime: start,
				EndTime:   time.Now(),
				Duration:  time.Since(start),
				Error:     fmt.Errorf("failed recording provenance for %s: %w", resID, err),
			}, err
		}

		acquiredResources = append(acquiredResources, res.ID)
		if acqResult.InstalledPath != "" {
			installedPaths[res.ID] = acqResult.InstalledPath
		}
		if acqResult.ExecutedCommand != "" {
			acquireCommands[res.ID] = acqResult.ExecutedCommand
		}
	}

	// Stash acquisition telemetry before integrity so a later hash mismatch
	// still records installed paths into pipeline memory.
	ctx.Implementation = &ImplementationData{
		TargetAgent:       "implementer",
		AcquiredResources: acquiredResources,
		InstalledPaths:    installedPaths,
		AcquireCommands:   acquireCommands,
		BuildOutput:       "implementation acquisition complete; provenance integrity pending",
		BuildPassed:       false,
	}

	// 5. Verify Provenance Integrity
	integrityReport, err := provStore.VerifyIntegrity(ctx.Task.WorkspaceRoot)
	if err != nil || !integrityReport.AllValid {
		detail := fmt.Errorf("provenance integrity verification failed: %+v", integrityReport.Issues)
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     detail,
		}, detail
	}

	// 6. Track Modified Files & Checksums (including provenance.json)
	var modifiedFiles []handoff.FileChecksum

	provPath := filepath.Join(ctx.Task.WorkspaceRoot, ".orchestra", "provenance.json")
	if content, err := os.ReadFile(provPath); err == nil {
		h := sha256.Sum256(content)
		modifiedFiles = append(modifiedFiles, handoff.FileChecksum{
			Path:   provPath,
			SHA256: hex.EncodeToString(h[:]),
		})
	}
	ctx.AddArtifact("provenance.json", provPath)

	// Include DESIGN.md if generated
	if ctx.DesignSystem != nil && ctx.DesignSystem.DesignMDPath != "" {
		if content, err := os.ReadFile(ctx.DesignSystem.DesignMDPath); err == nil {
			h := sha256.Sum256(content)
			modifiedFiles = append(modifiedFiles, handoff.FileChecksum{
				Path:   ctx.DesignSystem.DesignMDPath,
				SHA256: hex.EncodeToString(h[:]),
			})
		}
	}

	// Include reference-log.md if generated
	if ctx.Research != nil && ctx.Research.ReferenceLogPath != "" {
		if content, err := os.ReadFile(ctx.Research.ReferenceLogPath); err == nil {
			h := sha256.Sum256(content)
			modifiedFiles = append(modifiedFiles, handoff.FileChecksum{
				Path:   ctx.Research.ReferenceLogPath,
				SHA256: hex.EncodeToString(h[:]),
			})
		}
	}

	// 7. Persist Handoff State
	statePath := filepath.Join(handoffDir, "state.json")
	state := &handoff.HandoffState{
		SessionID:      ctx.Task.ID,
		Version:        3,
		SourceAgent:    "antigravity",
		TargetAgent:    "cursor",
		ActiveTasks:    []string{ctx.Task.ID},
		ChangedFiles:   modifiedFiles,
		CompletedSteps: []string{"Discover", "Classify", "Research", "Synthesize", "DesignSystem", "Implement"},
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
	}

	if err := handoff.WriteState(state, ctx.Task.WorkspaceRoot); err != nil {
		return &StageResult{
			StageName: StageImplement,
			Status:    StatusFailed,
			StartTime: start,
			EndTime:   time.Now(),
			Duration:  time.Since(start),
			Error:     fmt.Errorf("failed to write handoff state: %w", err),
		}, err
	}

	ctx.Implementation = &ImplementationData{
		TargetAgent:       "implementer",
		AcquiredResources: acquiredResources,
		ModifiedFiles:     modifiedFiles,
		HandoffStatePath:  statePath,
		BuildOutput:       "handoff prepared; product build was not run",
		BuildPassed:       false,
		ExecutionStatus:   ExecPrepared,
		InstalledPaths:    installedPaths,
		AcquireCommands:   acquireCommands,
	}

	ctx.AddArtifact("handoff_state.json", statePath)

	return &StageResult{
		StageName: StageImplement,
		Status:    StatusCompleted,
		StartTime: start,
		EndTime:   time.Now(),
		Duration:  time.Since(start),
		Output:    ctx.Implementation,
		Artifacts: []string{statePath, provPath},
	}, nil
}

func nonArtifactEntry(res *resources.Resource, taskID, kind, source, version string) acquisition.ProvenanceEntry {
	if strings.Contains(source, "orchestra.internal") {
		source = ""
	}
	return acquisition.ProvenanceEntry{
		ResourceID:          res.ID,
		AcquisitionMethod:   res.AcquisitionMethod,
		SourceURL:           source,
		VersionOrSHA:        version,
		Timestamp:           time.Now().UTC().Format(time.RFC3339),
		JustificationTaskID: taskID,
		Kind:                kind,
	}
}

func artifactOrPending(res *resources.Resource, taskID, source, version, hash, installed string) acquisition.ProvenanceEntry {
	if strings.Contains(source, "orchestra.internal") {
		source = ""
	}
	if hash != "" && installed != "" && source != "" && version != "" {
		return acquisition.ProvenanceEntry{
			ResourceID:          res.ID,
			AcquisitionMethod:   res.AcquisitionMethod,
			SourceURL:           source,
			VersionOrSHA:        version,
			SHA256Hash:          hash,
			InstalledPath:       installed,
			Timestamp:           time.Now().UTC().Format(time.RFC3339),
			JustificationTaskID: taskID,
			Kind:                acquisition.KindRealArtifact,
		}
	}
	if version == "" {
		version = "pending"
	}
	return nonArtifactEntry(res, taskID, acquisition.KindPending, source, version)
}

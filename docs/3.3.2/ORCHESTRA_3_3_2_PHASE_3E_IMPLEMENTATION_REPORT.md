# PHASE 3E IMPLEMENTATION REPORT

Orchestra 3.3.2 capability gateway. **PASS WITH LIMITATIONS.** VERSION on disk is still `3.3.1`. No commit, tag, or 3F.

Method: `C:\projects\orchestra-workflow`. Package: `runtime/internal/capability`. Engine hook after `refineMissionFromClassification`. Evidence: `.orchestra/runs/<id>/capability.json`.

---

## Verdict

The gateway selects a minimum set from the existing mission and budget, and each selected required id is consumed by an existing owner (builder capsule, research skip, design bind, Layer A, Layer B, or judge capsule). There is no `WorkerAdapter`, no host-execution router, no model registry, and no 1jehuang/jcode symbol in the capability package.

`go test ./internal/...` from `runtime/` passed. `node runtime/tools/check-host-stack.js` — 44 skills, 0 problems. `node runtime/tools/hygiene.js` — VERSION `3.3.1`, 0 failures.

## Consumption table

| Selected id | Consumer | Proof |
| --- | --- | --- |
| `coding.worker` | `compileBuilderCapsule` → `RoleBuilder` | Flow A/B/C capsules; omitted-id implement fail |
| `research.worker` | `ResearchStage.ShouldSkip` when `Capability` is set | Flow A skip; Flow C skip unless `Budget.Research != none` |
| `design.worker` | `bindFlagshipExploration` | Flow B `Exploration` nil; Flow C one sealed contract |
| `browser.playwright` | harness `PlaywrightVerifier` | Flow B real Chromium; tamper must not `AllPassed`; Flow E human gate |
| `visual.structural` | Layer B `Evaluate` | Flow B/C/D `evaluator_type` structural |
| `visual.screenshot` | adapter `Evaluate` only if selected and registered | Flow A invocations 0 with adapter registered |
| `visual.multimodal` | fallback record only | Flow D; `evaluator_type` never `multimodal` |
| `judge.independent` | `compileJudgeCapsule` | selected → role judge; not selected → nil |
| skill / resource / tool | capsule lists = decision lists | hostile dump test; caps 8 / 8 / 4 |

Tamper: product-UI budget wants browser, `browser.playwright` removed from the decision → visual QA does not `AllPassed` and does not run Layer A.

## Flow F — selection-only

Two compatible host snapshots (cursor, claude). Both available → evidence `host_label=cursor`. Cursor unavailable → `host_label=claude`. Mission JSON and capsule JSON have no provider/model token. Handoff `SourceAgent` / `TargetAgent` remain `antigravity` → `cursor`. `ExecutePlan` is not called. Execution routing stays outside 3E.

## Capsule and evidence bytes

3C scale was capsule ~716 bytes vs whole-project fixture ~136000.

| Case | Capsule bytes | `capability.json` bytes | Notes |
| --- | --- | --- | --- |
| Flow A ordinary | 542 | 265 | no PNG |
| Flow B product UI | 705 | — | Layer A desktop 8751 / mobile 5611 (same Chromium fixture class as 3D.1) |
| Flow C flagship | 678 | 2023 | judge selected by `DecideJudge` on this fixture |
| Hostile catalog | 852 | — | dump 152000 bytes; ratio 0.0056 |

Both stay far under the dump.

## Flows

- **A** API/debug: research skipped, no browser, screenshot eval invocations 0 with adapter registered, judge nil, no directions, no PNG.
- **B** product UI: coding worker consumed, `design.worker` absent, real Chromium Layer A on `visual/testdata/fixture.html`, structural Layer B, zero directions. `SkipVisualGate` waiver unchanged.
- **C** flagship: `design.worker` seals one contract, coding worker consumed, research follows budget, real Chromium, judge capsule only because `DecideJudge` selected it.
- **D** multimodal unavailable → structural; not screenshot, not multimodal.
- **E** `playwright_available=false` while `visual_review=browser` → `GATED_WAITING_APPROVAL`, question names `browser.playwright`, classify once, attempt does not loop. Chromium was not uninstalled.
- **F** selection-only host label as above.

3D.1 regression: `go test ./internal/...` includes `TestLayerASmokePlaywright` / Flow C Chromium. Flow B this phase logged desktop 8751 / mobile 5611 PNG bytes.

## Honest limits

- Live MCP `BROKEN` is not produced by today’s `classifyMCP` scanner (3F). Tests inject `BROKEN` / `AUTH_REQUIRED`.
- Screenshot and multimodal adapters are absent on this PC. Screenshot eval runs only when a test registers one. Multimodal has no executor.
- Cursor model discovery is not in the runtime. Opaque `DiscoveredModels` on `Input` (test label `grok-4.7`) is ignored; capability package source has no `grok` literal.
- 1jehuang/jcode is not installed and is not referenced from Go. Stance remains KEEP Orchestra / ADAPT later (not 3E) / DEFER swarms-embeddings-Firefox / REJECT a second conductor.
- OpenDesign stays retired. Ban lists were not edited.
- Implement handoff literals were not changed. `HostAdapter.ExecutePlan` remains a stub.

## Complexity boundary

One new package. No new stage, second budget, registry file, model list, `WorkerAdapter`, multimodal runner, or host-execution router. 3A–3D.1 stay locked.

## Stop

3E is implemented. Do not start 3F, Gate 0, Portfolio, VERSION bump, or a commit unless the operator asks.

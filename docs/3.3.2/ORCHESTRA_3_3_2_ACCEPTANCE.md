# Orchestra 3.3.2 — acceptance

Measurable definition of done. Phase 3A / 3A.1 were spec. **3C runtime is locked.** **3D runtime acceptance is live.** **3E capability gateway is live.** **3F Brain reconciliation is live** (section below). **3G is a measurement pilot** on T1 and T3. It is not an effectiveness claim. 3H–3I remain later.

Architecture fails if A–I products are not materially better than **representative 3.3.1 behavior** without unjustified cost or complexity. Cleaner docs are not success.

---

## 3C runtime DoD (this phase)

Pass when all of the following are true. Commands: `go test ./internal/mission/ ./internal/engine/` from `runtime/`.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | A mission can be represented and validated | `TestMissionValidRoundTrip` / `TestMissionInvalid` |
| 2 | Worker receives a compact scoped capsule | `TestCompileCapsuleScopesContext` — unrelated files, skills, chat, MCP catalog, whole-project dump excluded |
| 3 | Execution budget is derived, not universally maximized | `TestDeriveBudgetCheapForOrdinaryImplementation` vs `TestDeriveBudgetStrongerForHighRisk` |
| 4 | Failure classes are explicit | `TestClassifyRepresentativeFailures` |
| 5 | Recovery is bounded and failure-aware | `TestRecoveryRetryableWithinBound`, `TestRecoveryBudgetExhaustedStops` |
| 6 | Human-required failures stop automatic execution | `TestRecoveryHumanNeededStops`, `TestClassifyHumanNeededSituations`, `TestPipeline_HumanGateStopsAndRecordsState`, `TestContinueStopsOnHuman` |
| 7 | Run state survives handoff / interruption / retry | `TestRunStateRoundTripAndContinue`, `TestInterruptedContinueDoesNotNeedConversation`, `TestPipeline_MissionStateAndCheapBudget` |
| 8 | Retry history is inspectable | `TestRecoveryIdenticalFailuresTerminate`, `TestRecoveryDifferentFailureAfterRetryStillRetries`, `TestPipeline_RetryHistoryRecordedOnVisualLoop` |
| 9 | Provider-specific details do not leak into core contracts | `TestMissionJSONHasNoProviderLeak`, `TestCompileCapsuleProviderNeutral`, `TestDeriveBudgetNoProviderNames`, `TestMissionAllowsProviderNameInObjective` |
| 10 | Existing Orchestra architecture remains intact | 8-stage engine, one conductor, STANDARD/PREMIUM/EXPERIMENTAL, two-lock Design Lab still pass `./internal/engine`; `TestPipeline_SentinelsUnchangedAfterBind` |
| 11 | Existing tests continue passing | engine + classifier + verify suites |
| 12 | New tests demonstrate the mechanisms | `runtime/internal/mission/*_test.go`, `runtime/internal/engine/mission_bind_test.go` |
| 13 | No 3D–3I mechanism implemented | no Layer A/B judge runtime, no capability gateway, no MCP/Brain/OpenDesign/Portfolio edits in 3C |
| 14 | Capsule/context scope is measurably smaller than whole-project context | `ScopeReport.Ratio < 1` when whole-project bytes are provided (`TestCompileCapsuleScopesContext`, `TestCompileIgnoresExtraIrrelevantContext`) |
| 15 | Complexity stays inside the 3C budget | only mission, capsule, budget, failure class, bounded recovery, filesystem run state |
| 16 | Host conductor parity | live Cursor/Gemini/Agents/Claude/jcode `orchestra-conductor` SHA-256 matches `skills/orchestra-conductor/SKILL.md` after 3C.1 one-skill copy |

Limitations that do **not** fail 3C: Healed iterate must not leave blocking `pending_human` (`TestPipeline_HealedIterateDoesNotLeaveHumanBlock`). 3C recorded `visual_review` / `independent_judge` as policy; 3D consumes them.

## 3D runtime DoD

Pass when all of the following are true. Commands: `go test ./internal/...` from `runtime/`, plus `node runtime/tools/check-host-stack.js` and `node runtime/tools/hygiene.js` from the repo root.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Ordinary functional work skips directions and visual review | `TestFlowA_APIDebugStaysCheap`, `TestClassifyOrdinaryAPI` |
| 2 | Normal UI: no 3–5 directions; Layer A policy on; judge not required | `TestFlowB_CRUDUINoDirectionsNoJudge`, `TestClassifyProductUINotFlagship` |
| 3 | Flagship language, not quality bar, triggers directions | `TestInferVisualDoesNotUseQualityBar`, `TestClassifyFlagshipFromLanguageNotBar` |
| 4 | 3–5 structurally different directions; one sealed contract | `TestFlagshipDirectionsPolicy`, `TestFlowD_FlagshipDirectionsThenLayers` |
| 5 | Anti-slop is actionable, not a style blacklist | `TestCriticNotABlacklist`, `TestCriticActionableNotGeneric`, `TestCriticOrdinaryNotBlocked` |
| 6 | Canonical Layer A owner is the harness; mocks cannot pass | `TestPlaywrightVerifierFailClosedNoApp`, `TestLayerARejectsMockScreenshot`, `TestLayerAFailClosedWithoutURL`, `TestLayerASmokePlaywright` |
| 7 | Running-app responsive evidence; desktop ≠ mobile | `TestFlowC_ResponsiveRunningApp`, `TestLayerAResponsiveRealBrowser`, `TestLayerAResponsiveFixtureEvidence` |
| 8 | Layer B requires rendered PNG evidence | `TestLayerBRequiresRenderedEvidence`, `TestLayerBStructuredFindingsFromScreenshots`, `TestScreenshotEvalOnlyWhenPolicyAndAdapter` |
| 9 | Independent judge is adaptive and scoped | `TestDecideJudgeOrdinaryOff`, `TestDecideJudgeHighRiskRequired`, `TestJudgeCapsuleScoped` |
| 10 | Visual failure → structured heal; identical fail stops | `TestFlowE_VisualHealChangesImplementation`, `TestFlowF_IdenticalFailureStops` |
| 11 | Sentinels unchanged | `TestPipeline_SentinelsUnchangedAfterBind`, `TestFlowF_IdenticalFailureStops` |
| 12 | Durable visual evidence | `TestDurableVisualEvidenceRoundTrip` |
| 13 | Golden stills remain pre-implementation | existing `verify` Design Lab tests; 3D does not treat Layer A shots as stills |
| 14 | 3F–3I not implemented | no MCP/Brain overlay rewrite, OpenDesign repair, VERSION bump |

Limitation that does **not** silently pass: 3D ran without a project-local Playwright library (wrong `npx playwright --version` probe). **3D.1** installs `runtime/package.json` + Chromium and requires live pixel tests (`TestLayerASmokePlaywright`, `TestLayerAResponsiveRealBrowser`). If Chromium cannot launch after that install, those tests **BLOCKED**-fail rather than mock-pass.

## 3D.1 hardening DoD

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Project-local `playwright` library + Chromium | `runtime/package.json`, `PlaywrightAvailable` / `ProbePlaywright` |
| 2 | Real Layer A pixels; MOCK_SCREENSHOT fails | `TestLayerASmokePlaywright`, `TestLayerARejectsMockScreenshot` |
| 3 | Desktop 1440 ≠ mobile 390 overflow | `TestLayerAResponsiveRealBrowser`, `TestFlowC_ResponsiveRunningApp` |
| 4 | Layer B structural default; screenshot eval optional, no hardcoded model | `TestScreenshotEvalOnlyWhenPolicyAndAdapter`, `TestScreenshotEvaluatorProviderNeutral`, `TestOrdinaryNeverInvokesScreenshotEval` |
| 5 | 1jehuang/jcode not installed | [ORCHESTRA_3_3_2_JCODE_EVAL.md](ORCHESTRA_3_3_2_JCODE_EVAL.md) |
| 6 | 3E not started in 3D.1; VERSION unchanged | 3E is a later phase; `VERSION` still 3.3.1 |

## 3E capability gateway DoD

Pass when all of the following are true. Commands: `go test ./internal/...` from `runtime/`, plus `node runtime/tools/check-host-stack.js` and `node runtime/tools/hygiene.js` from the repo root.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Minimum set is actually consumed | `TestPhase3E_FlowA_APIDebugConsumed`, `TestPhase3E_FlowB_ProductUIConsumed`, `TestPhase3E_FlowC_FlagshipConsumed` |
| 2 | Tamper: required `browser.playwright` removed must not pass | `TestPhase3E_TamperBrowserRemovedFailsVisualQA` |
| 3 | Browser required + unavailable → blocking human, no loop | `TestPhase3E_FlowE_BrowserUnavailableHumanGate` |
| 4 | Multimodal unavailable falls back to structural; no fake vision | `TestPhase3E_FlowD_MultimodalFallbackStructural` |
| 5 | Flow F is selection-only; handoff literals unchanged | `TestPhase3E_FlowF_HostSelectionOnly` |
| 6 | Hostile skill/resource/MCP dump stays out of the capsule | `TestPhase3E_HostileCatalogStaysOutOfCapsule` |
| 7 | Opaque discovered-model labels do not change selection or leak | `TestPhase3E_OpaqueModelDoesNotChangeSelection`, `TestSelect_OpaqueDiscoveredModelIgnored`, `TestSelect_SourceHasNoVendorLiterals` |
| 8 | No `WorkerAdapter` / no 1jehuang/jcode in the capability package | `TestSelect_SourceHasNoVendorLiterals` |
| 9 | VERSION unchanged | `VERSION` still `3.3.1`; hygiene |

Limitations that do **not** fail 3E: screenshot and multimodal adapters are absent on this PC; Cursor model discovery is not in the runtime. Injected MCP probes, including `BROKEN`, are 3F. The budget checker still does not spawn servers.

## 3F Brain reconciliation DoD

Pass when all of the following are true. Commands: `go test ./internal/...` from `runtime/`, plus `node runtime/tools/check-host-stack.js` and `node runtime/tools/hygiene.js` from the repo root.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Runtime owns mission, budget, capsule, runs, capability, visual evidence, retries | `TestPhase3F_OwnershipRuntimeWins` |
| 2 | Write allow-list and deny-list | `TestPhase3F_AdmitDurable` |
| 3 | Cap 8; stale, superseded, and legacy aggregates drop unless history is requested | `TestPhase3F_SelectCapAndStale` |
| 4 | Conflict order: runtime, docs, registry, brain, unverified | `TestPhase3F_ConflictOrder` |
| 5 | Private fixture cannot enter a public artifact | `TestPhase3F_PublicPrivate` |
| 6 | Legacy 2648 with empty payload is not current evidence and survives the next Record | `TestPhase3F_LegacyAggregateSplit` |
| 7 | Unprobed MCP is OPTIONAL; a probe can emit all five states | `TestClassifyMCP_UnprobedIsOptional`, `TestClassifyMCP_ProbeStates` |
| 8 | Missing `last_verified` is unverified | `TestPhase3F_MissingLastVerifiedIsUnverified` |
| 9 | Retrieval stays far under a ~136000-byte dump | `TestPhase3F_ContextBytes` |
| 10 | VERSION unchanged | `VERSION` still `3.3.1`; hygiene |

Limitations that do **not** fail 3F: the 126 unverified catalog rows are not stamped; shadcn is not probed; OpenDesign is not repaired; 1jehuang/jcode is not installed; 3G is not started.

## 3F.1 hardening DoD

Pass when the 3F invariants hold on the existing path. No new memory engine.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Health does not define requirement; unprobed cannot satisfy a required tool | `TestPhase3F1_HealthDoesNotDefineRequirement`, `TestPhase3F1_RequiredUnprobedUsesHumanGate` |
| 2 | Bounded Brain facts can enter a capsule; catalogs, history, and dumps cannot | `TestPhase3F1_BrainToCapsule` |
| 3 | Unavailable, empty, malformed, stale, or missing Brain does not stop the run | `TestPhase3F1_BrainFailureDoesNotBlock`, `TestPhase3F1_BrainMissStillCompiles` |
| 4 | Public serialization rejects private facts, private paths, full private documents, and run dumps | `TestPhase3F1_PublicSerialize` |
| 5 | Retrieval does not upgrade provenance | `TestPhase3F1_ProvenanceStaysPut` |
| 6 | 2648 stays a legacy aggregate | `TestPhase3F1_LegacyCannotBecomeCurrent` |
| 7 | 168 / 42 / 126 and missing dates stay unverified | `TestPhase3F1_RegistryCountsStayPut` |
| 8 | Opaque model labels do not change decisions | `TestPhase3F1_ModelLabelIsNotArchitecture` |

---

## 3G pilot DoD

Pass when the harness can measure T1 and T3. A pass is not evidence that Orchestra improves output. Rows A–E stay 3H. Rows F / G / H / I map to later task classes; this pilot runs only the T1 and T3 fixtures.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Two traces from the same fixture hash match | `TestPhase3G_Traces` |
| 2 | Baseline and Orchestra start equivalent; Orchestra artifacts stay out of the baseline | `RunTrace` integrity checks |
| 3 | LIGHT stays cheap on T1 and T3 | zero directions, browser, screenshots, research, judge, and Brain facts |
| 4 | Unnecessary capabilities stay off | selected id is `coding.worker` |
| 5 | Seeded T3 failure retries once, then stops | `TestPhase3G_RecoveryStopsIdenticalLoop` |
| 6 | A positive judgment cannot pass a failed test | `TestPhase3G_ScoreIgnoresJudgment` |
| 7 | Context is bytes; tokens stay unavailable | mission, capability, and capsule byte fields |
| 8 | T2 and T4–T10 stay disabled | `TestPhase3G_ExpansionGateStaysShut` |

The paired host cell is `unavailable`. Claim stays `none`.

---

## 3G.1 host pilot DoD

Pass when a fresh host can be compared on T1 and T3, or the cell is explicitly `unavailable`. A pass is not an effectiveness claim.

| # | Criterion | Evidence |
| --- | --- | --- |
| 1 | Empty host executable stays `unavailable` | `TestPhase3G1_Unavailable` |
| 2 | Different model labels are `invalid_pair` | `TestPhase3G1_LabelMismatch` |
| 3 | Baseline stays free of Orchestra artifacts | `TestPhase3G1_T1Noop`, `TestPhase3G1_T3Fix` |
| 4 | Editing the T3 test cannot pass | `TestPhase3G1_T3Cheat` |
| 5 | An extra T3 file is a scope failure | `TestPhase3G1_T3ExtraFile` |
| 6 | A timeout is not a win for the other arm | `TestPhase3G1_Timeout` |
| 7 | T2 and T4–T10 stay disabled | `TestPhase3G1_ExpansionStaysShut` |

The real pair runs only with `ORCHESTRA_BENCH_REAL=1`. Default `go test` uses a fake host.

---

## Baseline comparison (3.3.1 vs 3.3.2)

For each A–I class that is actually run in 3H, record 3.3.1 and 3.3.2 on:

| Axis | What to record |
| --- | --- |
| Output quality | Human Keep/Rework/Kill plus Layer A/B evidence where visual |
| Reliability | Completion without silent skip of required gates |
| Speed | Wall time to an accepted artifact |
| Token / context cost | Approximate tokens / context dumps (order of magnitude is enough) |
| Human intervention | Gates hit, unsolicited questions, stills/direction picks |
| Failure count | Classified failures |
| Retry count | Retries that changed something vs identical replays (replays are a fail) |
| Completion rate | Finished vs abandoned |
| Implementation complexity | New runtime concepts / files added to Orchestra itself |
| Operational overhead | Extra hosts, MCPs, skills loaded, eval agents spawned |

No invented numeric targets. Use evidence from those runs.

**Anti-pattern (not an improvement unless the operator explicitly accepts the cost):** ~10% prettier, ~4× slower, ~2× more expensive.

A change is beneficial when outcomes improve on the axes that matter for that class **and** cost/complexity do not jump without a reason recorded on the mission.

“Materially better” is **not** “at least one axis up, others ignored.”

---

## Process / parity DoD (3B+)

| Gate | Pass when |
| --- | --- |
| Host / Brain / docs parity | Live Cursor + Gemini conductor skills match the repo contract. `~\.gemini\GEMINI.md` matches. Brain overlay docs no longer say “3.2 contract.” `.cursorrules` / README teach 3 translations, not skip. |
| Skill allowlist | One source. `sync-ides.ps1` and `kit/antigravity/install.ps1` install the same names. Live hosts have the four job-loaded skills. |
| Registry integrity | Every row has honest verification state. Overlay `total_evaluations` matches payload or is reset. GitHub MCP has a row or is absent. Broken tools marked broken. Duplicates have one owner. |
| Routing | Mission records model, reasoning, host, MCP/skill/resource recs. No silent bad provider. No silent quality downgrade. OpenDesign not selected as conductor. Budget derived, not a one-size preset. |
| Recovery | Failures classified. Retry only with remaining budget **and a meaningful change**. FULL long path has stuck detection. LIGHT does not hang in eval theatre. Independent judge is not spawned merely because intensity is FULL. |
| Design | Two locks still required for FULL product UI. Named refs = 3 translations. Flagship may use 3–5 **lightweight** directions; **one** implementation. CRUD/API/debug: zero extra directions. |
| Visual QA | No mock-screenshot pass. Layer A (Playwright) owns objective pass/fail. Layer B owns perception when `visual_review` is browser/deep. Stills remain pre-implementation. One browser owner. |
| Functional QA | Mechanical tests still gate LIGHT. API/CRUD/debug do not require Design Lab. |
| Security | Human gate on security-sensitive and irreversible deploy. Independent security eval when risk warrants. No secrets in the vault. |
| Complexity | No DAG OS, no fourth quality bar, no always-on Magic UI, no OpenDesign-as-core, no bulk skill/MCP install, no always-on independent judge. |
| Hygiene | 3I: leftover 3.3.1/3.2/3.1 pieces classified KEEP / MIGRATE / DELETE / ARCHIVE; no known dead runtime modules, stale duplicate routing, or misleading active docs. Classification in 3A.1 is spec-only — no deletion in this phase. |
| Safety lock | 3.3.1 KEEP / FIX / DEFER / REMOVE / REPLACE / INVESTIGATE respected. One conductor preserved. |

---

## Output benchmarks A–I

Each benchmark is a **real product class** on this PC, not a synthetic prompt.

| ID | Product class | What “better” means vs 3.3.1 |
| --- | --- | --- |
| **A** Portfolio | Flagship creative; named visual refs; Elaine Yu as bar not clone | Lightweight distinct directions (not five palettes, not five sites). Stills lock. Layer A + Layer B. Anti-slop catches indecision, not gradients. Independent judge only if triggers fire. |
| **B** Cinematic | Motion-heavy product UI | FULL path. Motion is a direction axis, not a library dump. Layer A on the running app; Layer B for motion perception when budgeted. |
| **C** Parked | Parked product; treat as “if revived” | Same **path** as flagship creative. Do not build in 3.3.2. |
| **D** Operator HUD | Existing product bar | Reliability + one browser owner. No new always-on skills. |
| **E** Hiring flagship | Flagship first | FULL. Gateway does not pick a weak executor. Fresh-context / independent judge **if triggers fire**, not automatically. |
| **F** CRUD | Glue / tables / forms | LIGHT. No Design Lab tax. Mechanical tests suffice. Faster and cheaper than FULL theatre. |
| **G** Debug | One failing behavior | LIGHT. Capsule with files + known failure. One classified retry that changes something, then human. |
| **H** API | Backend contract | LIGHT. Mechanical + optional architecture eval if risk high. No golden stills. No Layer B. |
| **I** Responsive UI | Breakpoints, overflow, tap targets | Layer A (computed styles, widths, overflow). Completeness (empty/error/loading). Layer B only if art direction is in scope. Not a second browser MCP. |

Class C stays parked. The C row is a **path acceptance**, not a build order.

---

## Improvement × benchmark map

Expected effect. Q = quality, S = speed, R = reliability, C = complexity (↓ is good). Cost must not silently explode.

| # | Improvement | A Port | B Motion | C Parked* | D HUD | E Flagship | F CRUD | G Debug | H API | I Resp |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | Layer A Playwright, no mocks | Q+ R+ | Q+ R+ | Q+ R+ | Q+ R+ | Q+ R+ | — | — | — | Q+ R+ |
| 2 | Host/docs/skill parity | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ |
| 3 | Adaptive independent judge | Q+ C↓ vs always-on | Q+ if triggered | Q+ if triggered | Q+ if triggered | Q+ if triggered | S+ (off) | S+ (off) | S+ (off) | Q+ if triggered |
| 4 | Worker capsule | Q+ R+ | Q+ | Q+ | R+ | R+ | S+ R+ | S+ R+ | R+ | R+ |
| 5 | Failure class + retry-must-change | R+ | R+ | R+ | R+ | R+ | R+ S+ | R+ S+ | R+ | R+ |
| 6 | Anti-slop critic / 2026 clichés | Q+ | Q+ | Q+ | Q+ | Q+ | — | — | — | Q+ |
| 7 | Lightweight 3–5 directions (flagship) | Q+ C↓ vs 5 builds | Q+ | Q+ | — | Q+ | — | — | — | — |
| 8 | Mission file / immutable acceptance | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ |
| 9 | Capability → runtime recommendation | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ |
| 10 | Registry / overlay honesty | R+ C↓ | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ |
| 11 | Specified human gate | Q+ R+ | Q+ | Q+ | R+ | R+ | R+ (less guess) | R+ | R+ | Q+ |
| 12 | LIGHT intensity | — | — | — | S+ on glue | S+ on glue | **S+ C↓** | **S+ C↓** | **S+ C↓** | S+ if mechanical |
| 13 | One browser owner + Layer B when budgeted | C↓ R+ | C↓ Q+ | C↓ | C↓ | C↓ | C↓ | C↓ | — | C↓ R+ |
| 14 | Optional OD worker (later) | Q? if used | Q? | Q? | — | — | **must stay off** | off | off | off |
| 15 | Memory classes | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ | R+ |
| 16 | Execution budget | C↓ | C↓ vs max lane | C↓ | C↓ | C↓ | **S+ C↓** | **S+ C↓** | **S+ C↓** | C↓ |

`*` Class C path only. Do not un-park that product for 3.3.2.

**Must move the needle:** 1, 2, 3, 4, 5, 6, 7 on A/B/E; 12 and 16 on F/G/H. If those land as docs-only, 3.3.2 failed even if every other file is tidy.

**Must not regress:** LIGHT work getting Design Lab; FULL work skipping stills; OpenDesign becoming a second studio; mock screenshots returning; FULL always spawning an independent judge; 3–5 full implementations; Playwright grading art direction.

---

## Mechanism matrix

Every major 3.3.2 mechanism. Six complexity questions are folded into problem / evidence / complement / cost / failure / removal.

| Mechanism | Problem | Expected benefit | Evidence | Cost | Acceptance test | Failure behavior | Removal condition |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Mission file | Run truth lives in chat sludge | Immutable acceptance; cheaper handoff | Nonstop / Spec Kit artifacts; Orchestra already has gates | Small file I/O | Mission exists on STANDARD+ runs; acceptance not silently rewritten | Human gate if missing acceptance | Drop if a single chat is enough for LIGHT-only shops |
| Intensity LIGHT/FULL | Ceremony on glue; skipped stills on flagship | Right-sized path | minhtran lanes; Anthropic eval simplification | Conductor classify | LIGHT skips Design Lab; FULL cannot skip stills on product UI | Mis-classify → human or classified retry | Merge back into bars only if intensity never diverges |
| Execution budget | One-size FULL lane | QUALITY × EVIDENCE / COST | Anthropic later removed always-on eval | Derive + record | API bug stays cheap; portfolio not forced through max eval | If budget ignored, treat as harness fail | Remove fields models no longer need |
| Worker capsule | Brain dump into workers | Fewer wrong files | Tenet `compile_context` | Packager | Capsule has files/objective/failures only | Too-fat capsule → context class failure | Collapse into prompt if models ignore extras |
| Evaluator policy + conditional judge | Builder grades itself **or** always-on GAN | Independent verdict only when warranted | cwc no-write; ~20x cost if always-on | Extra context when on | Judge off on F/G/H; on only when a listed trigger fires | Conflicting eval → judge or human | Turn off when self-eval is empirically honest |
| Two-layer visual QA | Mocks; Playwright as taste; duplicate browsers | Objective + perception, one browser | Canary Layer A; screenshot eval Layer B | Browser run + optional Layer B | Mocks cannot pass; Layer A geometry pass ≠ Layer B pass | Layer A fail = product/test; Layer B fail = visual | Layer B removable if humans stills suffice |
| Lightweight 3–5 directions | Five palettes **or** five full builds | Divergence before commitment | Design Lab already two-lock | Cheap sketches + one build | Flagship: 3–5 sketches, one impl; CRUD: zero | Skip pick → human gate | Skip exploration when contract is pasted |
| Failure class + retry-must-change | Blind loops | Bounded recovery | Tenet / Ralph / Zenith | Taxonomy + N | Identical replay fails the harness; N exhausted → human | Class recorded on each miss | Shrink taxonomy if unused |
| Runtime recommendation | Silent weak/expensive pick | Right model/host/tool | twaldin RunSpec idea | Rec on mission | Recorded rec; human if material choice missing | Ambiguous → `ErrHumanGateRequired` | Collapse if one host forever |
| Registry / memory classes | Fake 2648 evals; unverified counted live | Honest routing | Freeze audit | Hygiene | Overlay count matches payload; unverified not claimed verified | INVESTIGATE overlay before delete | Memory plugin still deferred |
| Human-needed gate | Conductor assumes | Stop on identity/direction/destroy/security/deploy/tools | `ErrHumanGateRequired` exists | Latency | Listed situations halt | Missing gate on stills = fail | Do not remove; may shorten the list |
| Token ladder | Rediscovery | Lower tokens | Progressive disclosure already policy | Discipline | Capsule ≠ catalog dump | Context class if dump detected | Keep even if models cheapen |
| Repo hygiene classes | Dead/stale modules in “active” tree | Active architecture only | Freeze: mock QA, skip docs, 30-name install | 3I pass | KEEP/MIGRATE/DELETE/ARCHIVE assigned; no misleading active docs | Unknown leftover → INVESTIGATE, not silent DELETE | N/A at release |

---

## Release hygiene (3I gate — classify now, delete later)

For every leftover 3.3.1 / 3.2 / 3.1 component:

| Class | Meaning |
| --- | --- |
| **KEEP** | Actively required by 3.3.2 |
| **MIGRATE** | Needed function in an old form; fold into 3.3.2 |
| **DELETE** | Dead, duplicated, obsolete, transitional; no active value. Evidence required. |
| **ARCHIVE** | Useful history; must not participate in runtime or active workflow |

Do not keep a module because it is old. Do not delete because it looks old.

**Release requires:** no known dead runtime modules; no stale duplicate routing; no obsolete version-specific instructions in the active workflow; no misleading docs; no broken references; no unused scaffolding kept only for history; no duplicate systems where one canonical owner is enough.

### Seed classifications (spec only — **do not delete files in 3A.1**)

| Item | Class | Why |
| --- | --- | --- |
| Two-lock Design Lab, 8-stage engine, one conductor | KEEP | Active OS |
| `ErrHumanGateRequired` | KEEP | Extend when-list only |
| Job-loaded four skills (in-repo) | KEEP | Hosts must catch up (3B) |
| Named-ref 3 translations in engine/AGENTS | KEEP | Docs lag is the bug |
| Real-app Playwright **policy** | KEEP | Implementation FIX |
| `MOCK_SCREENSHOT_*` success path | DELETE (3D) | Not verification |
| `.cursorrules` / README named-website skip | MIGRATE (3B) | Teach 3 translations |
| `ARCHITECTURE.md` stack-card gate prose | DELETE/fix (3B) | Stills replaced it |
| `kit/antigravity/install.ps1` 30 names | MIGRATE (3B) | Same allowlist as `sync-ides.ps1` |
| Duplicate Playwright / Context7 / motion surfaces | MIGRATE (3D/3F) | One owner each |
| Overlay `total_evaluations: 2648` vs empty array | INVESTIGATE then MIGRATE or DELETE counter | Honesty |
| `docs/3.3.1/CAPABILITY_AUDIT.md` | ARCHIVE or KEEP as reference | Must not contradict 3.3.2 active docs |
| GSD/ECC/Karpathy/always-on Magic UI as core | KEEP REJECTED | Not reintroduced |
| OpenDesign as second studio | KEEP banned | Optional lab only |

3A.1 does **not** clean the repository.

---

## Explicit non-goals (fail if done)

- Replacing the 3.3.1 OS.
- Tagging `v3.3.2` in 3A / 3A.1.
- Running `kit/sync-ides.ps1` as part of 3A / 3A.1.
- Installing skill packs or extra MCP/plugins.
- Building the parked product, Portfolio, or Gate 0 as a 3.3.2 vehicle.
- Pixel-cloning Elaine Yu or any named ref.
- Treating FULL as always-on independent judge.
- Implementing 3–5 complete design directions.
- Deleting runtime files during 3A.1.

# Orchestra 3.3.2 — reality audit

Status tracker for [ORCHESTRA_3_3_2_ACCEPTANCE.md](ORCHESTRA_3_3_2_ACCEPTANCE.md). Decisions stay in [ORCHESTRA_3_3_2_RECONCILIATION.md](ORCHESTRA_3_3_2_RECONCILIATION.md). This file does not replace either one.

Read on 2026-10-15. VERSION is `3.3.1`. No suite was re-run. Named tests were confirmed present. The last full `go test ./...` from `runtime/` passed on 2026-10-10, before `TestPhase3G24_RealT3` existed. That opt-in test skips unless `ORCHESTRA_BENCH_REAL=1`. It does not change 3C–3F behavior. Focused 3G checks passed the same day.

`PASS` means the named test exists and that 2026-10-10 run covered its package. It does not mean a fresh run today. `claim` for output quality stays `none`.

## Verdicts

- Engineering: 3C, 3D, 3D.1, 3E, 3F, and 3F.1 acceptance rows that name a test are `PASS`, except the live-host hash row, which is `UNVERIFIED`. Phase-scope locks that later phases were allowed to lift are `NOT APPLICABLE` as current-tree invariants.
- Output: nothing in this repo shows that Orchestra caused a better website. Anecdotes stay anecdotes.
- Version: `VERSION` is `3.3.1`. A later prerelease tag does not change this verdict. Nothing in 3G shows this tree beats older workflows at frontend or 3D work. Stable `v3.3.2` waits for phase 3I.

## 3C

Evidence lives in `runtime/internal/mission` and `runtime/internal/engine`.

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 1 | PASS | `TestMissionValidRoundTrip`, `TestMissionInvalid` | none | none |
| 2 | PASS | `TestCompileCapsuleScopesContext` | none | none |
| 3 | PASS | `TestDeriveBudgetCheapForOrdinaryImplementation`, `TestDeriveBudgetStrongerForHighRisk` | none | none |
| 4 | PASS | `TestClassifyRepresentativeFailures` | none | none |
| 5 | PASS | `TestRecoveryRetryableWithinBound`, `TestRecoveryBudgetExhaustedStops` | none | none |
| 6 | PASS | `TestRecoveryHumanNeededStops`, `TestClassifyHumanNeededSituations`, `TestPipeline_HumanGateStopsAndRecordsState`, `TestContinueStopsOnHuman` | none | none |
| 7 | PASS | `TestRunStateRoundTripAndContinue`, `TestInterruptedContinueDoesNotNeedConversation`, `TestPipeline_MissionStateAndCheapBudget` | none | none |
| 8 | PASS | `TestRecoveryIdenticalFailuresTerminate`, `TestRecoveryDifferentFailureAfterRetryStillRetries`, `TestPipeline_RetryHistoryRecordedOnVisualLoop` | none | none |
| 9 | PASS | `TestMissionJSONHasNoProviderLeak`, `TestCompileCapsuleProviderNeutral`, `TestDeriveBudgetNoProviderNames`, `TestMissionAllowsProviderNameInObjective` | none | none |
| 10 | PASS | `TestPipeline_SentinelsUnchangedAfterBind` | none | none |
| 11 | PASS | 2026-10-10 `go test ./...` | not re-run on 2026-10-15 | re-run only if runtime code changes |
| 12 | PASS | `runtime/internal/mission/*_test.go`, `runtime/internal/engine/mission_bind_test.go` | none | none |
| 13 | NOT APPLICABLE | 3C scope lock. 3D, 3E, and 3F are later acceptance sections | reading it as a current ban would contradict those sections | leave the historical row; do not delete later phases |
| 14 | PASS | `TestCompileCapsuleScopesContext`, `TestCompileIgnoresExtraIrrelevantContext` | none | none |
| 15 | PASS | mission, capsule, budget, failure class, recovery, run state | none | none |
| 16 | UNVERIFIED | `runtime/tools/check-host-stack.js` compares `registries/host-stack.json` to `skills/` only | live Cursor, Gemini, Agents, Claude, and jcode conductor hashes were not compared | hash the live conductor copies before any host sync |

Healed-iterate limitation: `TestPipeline_HealedIterateDoesNotLeaveHumanBlock` exists. `PASS` on the same evidence rule.

## 3D

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 1 | PASS | `TestFlowA_APIDebugStaysCheap`, `TestClassifyOrdinaryAPI` | none | none |
| 2 | PASS | `TestFlowB_CRUDUINoDirectionsNoJudge`, `TestClassifyProductUINotFlagship` | none | none |
| 3 | PASS | `TestInferVisualDoesNotUseQualityBar`, `TestClassifyFlagshipFromLanguageNotBar` | none | none |
| 4 | PASS | `TestFlagshipDirectionsPolicy`, `TestFlowD_FlagshipDirectionsThenLayers` | none | none |
| 5 | PASS | `TestCriticNotABlacklist`, `TestCriticActionableNotGeneric`, `TestCriticOrdinaryNotBlocked` | none | none |
| 6 | PASS | `TestPlaywrightVerifierFailClosedNoApp`, `TestLayerARejectsMockScreenshot`, `TestLayerAFailClosedWithoutURL`, `TestLayerASmokePlaywright` | live pixels were not re-shot today | none unless Layer A code changes |
| 7 | PASS | `TestFlowC_ResponsiveRunningApp`, `TestLayerAResponsiveRealBrowser`, `TestLayerAResponsiveFixtureEvidence` | none | none |
| 8 | PASS | `TestLayerBRequiresRenderedEvidence`, `TestLayerBStructuredFindingsFromScreenshots`, `TestScreenshotEvalOnlyWhenPolicyAndAdapter` | none | none |
| 9 | PASS | `TestDecideJudgeOrdinaryOff`, `TestDecideJudgeHighRiskRequired`, `TestJudgeCapsuleScoped` | none | none |
| 10 | PASS | `TestFlowE_VisualHealChangesImplementation`, `TestFlowF_IdenticalFailureStops` | none | none |
| 11 | PASS | `TestPipeline_SentinelsUnchangedAfterBind`, `TestFlowF_IdenticalFailureStops` | none | none |
| 12 | PASS | `TestDurableVisualEvidenceRoundTrip` | none | none |
| 13 | PASS | Design Lab tests in `runtime/internal/verify`; Layer A shots are not stills | none | none |
| 14 | NOT APPLICABLE | 3D scope lock. 3F is a later section | same as 3C row 13 | leave it |

## 3D.1

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 1 | PASS | `runtime/package.json`; `ProbePlaywright` | none | none |
| 2 | PASS | `TestLayerASmokePlaywright`, `TestLayerARejectsMockScreenshot` | none | none |
| 3 | PASS | `TestLayerAResponsiveRealBrowser`, `TestFlowC_ResponsiveRunningApp` | none | none |
| 4 | PASS | `TestScreenshotEvalOnlyWhenPolicyAndAdapter`, `TestScreenshotEvaluatorProviderNeutral`, `TestOrdinaryNeverInvokesScreenshotEval` | none | none |
| 5 | PASS | [ORCHESTRA_3_3_2_JCODE_EVAL.md](ORCHESTRA_3_3_2_JCODE_EVAL.md) says the third-party harness is not installed | none | none |
| 6 | NOT APPLICABLE | 3E is a later section. `VERSION` is still `3.3.1` | the "3E not started" clause is historical | leave it |

## 3E

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 1 | PASS | `TestPhase3E_FlowA_APIDebugConsumed`, `TestPhase3E_FlowB_ProductUIConsumed`, `TestPhase3E_FlowC_FlagshipConsumed` | none | none |
| 2 | PASS | `TestPhase3E_TamperBrowserRemovedFailsVisualQA` | none | none |
| 3 | PASS | `TestPhase3E_FlowE_BrowserUnavailableHumanGate` | none | none |
| 4 | PASS | `TestPhase3E_FlowD_MultimodalFallbackStructural` | none | none |
| 5 | PASS | `TestPhase3E_FlowF_HostSelectionOnly` | none | none |
| 6 | PASS | `TestPhase3E_HostileCatalogStaysOutOfCapsule` | none | none |
| 7 | PASS | `TestPhase3E_OpaqueModelDoesNotChangeSelection`, `TestSelect_OpaqueDiscoveredModelIgnored`, `TestSelect_SourceHasNoVendorLiterals` | none | none |
| 8 | PASS | `TestSelect_SourceHasNoVendorLiterals` | none | none |
| 9 | PASS | `VERSION` is `3.3.1` | none | none |

Screenshot and multimodal adapters stay absent, as the acceptance limitations allow. That absence is not a 3E fail.

## 3F and 3F.1

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 3F 1 | PASS | `TestPhase3F_OwnershipRuntimeWins` | none | none |
| 3F 2 | PASS | `TestPhase3F_AdmitDurable` | none | none |
| 3F 3 | PASS | `TestPhase3F_SelectCapAndStale` | none | none |
| 3F 4 | PASS | `TestPhase3F_ConflictOrder` | none | none |
| 3F 5 | PASS | `TestPhase3F_PublicPrivate` | none | none |
| 3F 6 | PASS | `TestPhase3F_LegacyAggregateSplit` | none | none |
| 3F 7 | PASS | `TestClassifyMCP_UnprobedIsOptional`, `TestClassifyMCP_ProbeStates` | none | none |
| 3F 8 | PASS | `TestPhase3F_MissingLastVerifiedIsUnverified` | none | none |
| 3F 9 | PASS | `TestPhase3F_ContextBytes` | none | none |
| 3F 10 | PASS | `VERSION` is `3.3.1` | none | none |
| 3F.1 1 | PASS | `TestPhase3F1_HealthDoesNotDefineRequirement`, `TestPhase3F1_RequiredUnprobedUsesHumanGate` | none | none |
| 3F.1 2 | PASS | `TestPhase3F1_BrainToCapsule` | none | none |
| 3F.1 3 | PASS | `TestPhase3F1_BrainFailureDoesNotBlock`, `TestPhase3F1_BrainMissStillCompiles` | none | none |
| 3F.1 4 | PASS | `TestPhase3F1_PublicSerialize` | none | none |
| 3F.1 5 | PASS | `TestPhase3F1_ProvenanceStaysPut` | none | none |
| 3F.1 6 | PASS | `TestPhase3F1_LegacyCannotBecomeCurrent` | none | none |
| 3F.1 7 | PASS | `TestPhase3F1_RegistryCountsStayPut` | none | none |
| 3F.1 8 | PASS | `TestPhase3F1_ModelLabelIsNotArchitecture` | none | none |

Unverified catalog rows, an unprobed UI kit, and OpenDesign stay out of scope, as the 3F limitations say.

## 3G and 3G.1

| # | Status | Evidence | Gap | Next action |
| --- | --- | --- | --- | --- |
| 3G 1 | PASS | `TestPhase3G_Traces` | none | none |
| 3G 2 | PASS | `RunTrace` checks; `TestPhase3G1_T1Noop`, `TestPhase3G1_T3Fix` | none | none |
| 3G 3 | PASS | traces assert zero directions, browser, screenshots, research, judge, and Brain facts | none | none |
| 3G 4 | PASS | selected id `coding.worker` on the real Orchestra arms | none | none |
| 3G 5 | PASS | `TestPhase3G_RecoveryStopsIdenticalLoop` | none | none |
| 3G 6 | PASS | `TestPhase3G_ScoreIgnoresJudgment` | none | none |
| 3G 7 | PASS | byte fields set; tokens `unavailable` in the result JSON | none | none |
| 3G 8 | PASS | `TestPhase3G_ExpansionGateStaysShut`, `TestPhase3G1_ExpansionStaysShut` | none | none |
| 3G.1 1–7 | PASS | `TestPhase3G1_Unavailable`, `TestPhase3G1_LabelMismatch`, `TestPhase3G1_T1Noop`, `TestPhase3G1_T3Fix`, `TestPhase3G1_T3Cheat`, `TestPhase3G1_T3ExtraFile`, `TestPhase3G1_Timeout`, `TestPhase3G1_ExpansionStaysShut` | none | none |

The acceptance line that says the paired host cell is `unavailable` is stale. Real pairs exist under `runtime/benchmarks/results/`. Do not treat that sentence as the result. Do not edit the historical reports to match it.

## 3G evidence

`claim` is `none` on every saved pair. Tokens and host retries are `unavailable`. Both arms share one Cursor profile, including the MCP configuration and rules the CLI already loads. That is a confounder.

Corrected T1, `3g23-t1-api.json`: `compared`. Both arms exited 0, changed `sum.go` and `sum_test.go`, and were `task_completed`. Walls 166561 ms and 172559 ms. Diff bytes 592 and 1027. The earlier `3g2-t1-api.json` completion flags are superseded; that pair was `timeout`.

First T3, `3g2-t3-debug.json`: `compared`. Both arms changed only `sum.go`, 183 diff bytes, walls 117337 ms and 133935 ms. The host inherited `ORCHESTRA_HOME`.

Clean T3, `3g24-t3-debug.json`: `compared`. Both arms again changed only `sum.go`, 183 diff bytes. Walls 104157 ms and 129875 ms. `ORCHESTRA_HOME` was withheld. Fixture, task, and acceptance hashes match the first T3 row.

These tasks cannot score frontend or 3D quality. T3 tells both arms to change the return value. Matching diffs and a slower Orchestra arm are a cost observation on that bugfix, not a verdict on websites.

## Friends' sites

Status: `UNVERIFIED`.

This method repo has no prompt, repository, screenshot, or run log that ties a shipped site to Orchestra 2.0 or 3.0 versus a raw host. `runtime/testdata/phase7-qa/` holds fixture PNGs for a harness test. They are not that comparison.

Missing before any causal claim: the brief, the host, the Orchestra version, the prompt, the product repo, and a raw-host control of the same brief. Do not invent scores. Do not run a five-version bake-off from this audit.

## Process gates

| Gate | Status | Evidence | Next action |
| --- | --- | --- | --- |
| Repo skill allowlist | PASS | `kit/antigravity/install.ps1` reads `registries/host-stack.json`. The reconciliation note about 30 hardcoded names is stale. | leave the installer; do not treat the old note as current |
| Live host conductor parity | UNVERIFIED | `check-host-stack.js` does not hash installed conductor copies | hash before a host sync |
| Classes A–I output benchmarks | NOT APPLICABLE | acceptance says those runs are 3H | do not start 3H from this file |
| Release tag | NOT APPLICABLE | acceptance defers `v3.3.2` | do not tag |

## Remaining work

One gap blocks the acceptance rule that 3.3.2 must be materially better than representative 3.3.1 behavior. The mechanisms above can run. No paired visual artifact shows they improve a site.

| Item | Dependency | Pass test | Blocker |
| --- | --- | --- | --- |
| One visual comparison, not a new phase | operator names an existing run, or later approves one brief | same brief, same host and model, raw host versus this tree, functional checks and visual notes recorded separately, `claim` still `none` until a repeat | no such artifact is in this repo |

Do not implement that comparison here. Do not enable T2. Do not start 3H. Do not add a model router. Do not bump VERSION.

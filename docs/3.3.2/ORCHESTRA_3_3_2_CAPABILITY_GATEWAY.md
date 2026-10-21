# Orchestra 3.3.2 — capability gateway (3E)

Minimum sufficient set, actually consumed. Not a second conductor, not a host-execution router, not a model registry.

Source of truth: `runtime/internal/capability.Select`. Evidence: `.orchestra/runs/<id>/capability.json` beside `mission.json`. The decision blob does not live inside `mission.json`.

---

## Purpose

Given this mission and this budget, return the smallest capability set that can meet acceptance. Every selected **required** id has a downstream consumer. A selected required id that the pipeline ignores is a failing test.

```
Mission + Budget → Select → TaskContext.Capability → existing consumer
```

## Ownership

| Owner | What it decides / does |
| --- | --- |
| Graph router `Compose` | Candidate skill/resource ids (already narrowed) |
| Host-stack JSON | Static host / optional MCP state for live snapshots |
| **`capability.Select`** | Minimum set, strength, fallback hops, catalog cap, selection-only host label |
| `mission.Compile` | Worker capsule from the decision’s narrowed lists |
| Layer A `PlaywrightVerifier` | Pixels, only when `browser.playwright` is selected and available |
| Layer B structural / optional screenshot adapter | Perception, only when those ids are selected |
| `HostAdapter` / `AgentAdapter` | Host execution (3E does not call `ExecutePlan`) |
| Implement handoff | Still `SourceAgent: antigravity`, `TargetAgent: cursor` |

3E does **not** own: a new stage, a second budget, a registry file, a model list, MCP health scanning (`classifyMCP` / `BROKEN` is 3F), OpenDesign, 1jehuang/jcode install.

## Capability model

Sparse view (filled fields only): `id`, `category`, `strength` (`required|preferred|optional|fallback|unavailable|human_required`), `available`, `reason`.

Categories Select may emit:

| id | When |
| --- | --- |
| `coding.worker` | Implementation. Always required. |
| `research.worker` | Only if `Budget.Research` is `targeted` or `deep`. No new budget bump. |
| `design.worker` | Only if `visual.PolicyFor` says directions (flagship). Gateway does not invent directions. |
| `browser.playwright` | `visual_review` is `browser` or `deep`. Availability from `visual.ProbePlaywright()` unless a test override is set. Not `npx playwright --version`. Not Playwright MCP. |
| `visual.structural` | Review browser/deep and class is not ordinary. |
| `visual.screenshot` | Preferred on those same tasks; available only if `visual.ActiveScreenshotEvaluator()` is non-nil. |
| `visual.multimodal` | Preferred only on flagship + `deep`. No adapter in 3E → unavailable. |
| `judge.independent` | Only if `visual.DecideJudge` activates. Ordinary stays off. |
| `skill` / `resource` / `mcp.tool` | Only ids already on `Input` after the filter below. |

There is no `worker.optional`. There is no `WorkerAdapter` interface.

## Consumption (required)

| Selected id | Consumer |
| --- | --- |
| `coding.worker` | `compileBuilderCapsule` emits `RoleBuilder`. Implement fails if the decision omitted this id. Does not call `HostAdapter.ExecutePlan`. |
| `research.worker` | `ResearchStage.ShouldSkip` skips when `TaskContext.Capability` is set and this id is absent. When selected, the existing research coordinator runs. `Capability == nil` keeps the pre-3E skip. |
| `design.worker` | `bindFlagshipExploration` runs only when selected. Product UI leaves `Exploration` nil. Flagship with the id selected seals one contract. |
| `browser.playwright` | Visual QA calls harness `PlaywrightVerifier` only when selected and available. Ordinary writes no PNG. Tamper: product-UI budget wants browser, id removed → visual QA must not `AllPassed`. |
| `visual.structural` | Structural critic / `Evaluate` only when selected. Ordinary stays `evaluator_type=none`. |
| `visual.screenshot` | Adapter `Evaluate` only when selected **and** adapter non-nil. Selected but nil adapter: fallback recorded, structural runs, invocation count 0. Not selected: invocations stay 0 even if a test adapter is registered. |
| `visual.multimodal` | **Intentionally unexecuted.** Consumer is the fallback record into screenshot (if adapter) or structural. `evaluator_type` is never `multimodal`. |
| `judge.independent` | `compileJudgeCapsule` only when selected. Selected → role judge, no write tools. Not selected → `JudgeCapsule` nil. |
| skill / resource / tool ids | `Capsule.Skills`, `Resources`, `ToolsAllowed` equal the decision’s narrowed lists. |

Do not add a fake runner, fake browser, or fake vision call so a row looks wired.

## Minimum set and fallback

Permitted fallback chain: `visual.multimodal` → `visual.screenshot` → `visual.structural`. Each hop is recorded. Structural is never labeled vision.

Prohibited: `browser.playwright` required and unavailable → `ErrHumanGateRequired` with human-detail (what was required, why, what was unavailable, fallback prohibited). No silent skip, no retry. Same sentinel as Design Lab. `recoverStageError` uses the capability human-detail as `PendingHuman.Question`.

Optional broken MCP: skip and record. Required broken or `AUTH_REQUIRED` MCP: same human stop.

## Provider neutrality

Mission JSON and capsule JSON contain no provider or model token (`ProviderLeak`). Host labels (`cursor`, `claude`) live on `capability.json` evidence only.

`Select` does not rank Claude vs Cursor vs Gemini. It does not hardcode vendor model ids. Opaque `Input.DiscoveredModels` strings are ignored; selection matches capability class only. Swapping the opaque label does not change selected ids.

A future host adapter may supply model metadata. 3E does not scrape Cursor’s model picker and does not add a model registry.

## Catalog filter

Ownership: existing router or host-stack **candidate ids** → `Select` → narrowed ids → `mission.Compile`.

The gateway does not open the 44-skill allowlist, the ~168-row resource catalog, or the MCP example to discover new ids. Filter: drop duplicates; drop ids not relevant to selected categories; drop unavailable MCP (`BROKEN`, `AUTH_REQUIRED`, `PURCHASE_PENDING`, `DISABLED`); cap skills ≤ 8, resources ≤ 8, mcp tools ≤ 4; empty candidates stay empty.

Live MCP input is static host-stack state already in JSON, or a test snapshot. 3E does not open `mcp.json`. 3E does not teach `classifyMCP` to emit `BROKEN` (3F). Tests inject those health values.

## Visual split

Layer A/B stay the owners of pixels and perception. The gateway only decides whether those paths run.

Canonical Layer A remains `PlaywrightVerifier` + `runtime/tools/layer-a.mjs` + project-local Playwright library + Chromium. Availability is `ProbePlaywright()`, not `npx playwright --version`, not Playwright MCP in example config.

## Selection-only host note

3E does **not** route execution to a host. `ExecutePlan` is a stub. Implement handoff literals stay `antigravity` → `cursor`. Changing them would be a second worker router.

Flow F checks capability-level host selection only: two compatible snapshots; both available → evidence host label `cursor` (deterministic preference matching today’s implement target); cursor snapshot unavailable → the other snapshot. Execution routing stays outside 3E.

## 1jehuang/jcode (doc only)

Not a Go type, import, string match, binary check, or `exec` in the capability package or 3E engine hook.

Stance recorded from [ORCHESTRA_3_3_2_JCODE_EVAL.md](ORCHESTRA_3_3_2_JCODE_EVAL.md):

| Mechanism | Stance |
| --- | --- |
| Orchestra conductor, mission, budget, Stage 8, Design Lab, harness Layer A | **KEEP** |
| Headless `jcode run` as a later optional worker | **ADAPT later** means **not this phase** |
| Swarms / embeddings / Firefox Agent Bridge | **DEFER** |
| Second conductor / mission / retry / state | **REJECT** |

No reinstall. Distinct from OmniRoute host `jcode` on the host-stack list.

## OpenDesign

Retired as a second studio. Banned MCP name. Provider-leak token. 3E does not wire it, rank it, or repair it.

## Human gate

Existing `ErrHumanGateRequired`. No new gate type. `Decide` retry rules unchanged. `human_needed` does not loop.

## Model-discovery limitation

Runtime does not expose Cursor’s model picker. 3E will not scrape it. Adding a model does not require editing `Select`. An opaque adapter label in a test (including a current Cursor model id) is ignored by class matching.

## Evidence

`capability.json` is the authority for the run. `WriteDecision` refuses a later write whose selected set, strength, availability, fallback hops, rejected rows, or host label differ, and returns an acceptance failure (`ErrDecisionMutated`). The pipeline checks that before the next stage. There is no second store.

The file holds only mission id, selected, fallbacks, rejected, unavailable, and the selection-only host label. Skill, resource, and tool id lists used by the capsule stay in memory and are read back from `selected`. Rejected rows are capped so a hostile candidate dump cannot land in the file. No conversation history, full mission, full capsule, or provider/model dump.

## Complexity boundary

One package. No new stage. No second budget. No registry file. No model list. No `WorkerAdapter`. No multimodal runner. No host-execution router. VERSION stays `3.3.1`. 3F is next when the operator says so.

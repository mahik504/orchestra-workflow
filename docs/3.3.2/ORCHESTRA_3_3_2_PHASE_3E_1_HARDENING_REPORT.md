# PHASE 3E.1 HARDENING REPORT

Orchestra 3.3.2 capability gateway, acceptance hardening only. **PASS WITH LIMITATIONS.** VERSION is still `3.3.1`. No commit, push, tag, or 3F.

---

## Verdict

The decision written at `Select` stays authoritative. A later change to the selected set, strength, availability, or fallback hops is refused, the file is left as first written, and acceptance fails. Narrowed skill, resource, and tool ids are what the worker capsule receives. No second store, no new gate type, no MCP executor, no model registry.

## Files modified

- `runtime/internal/capability/select.go` — evidence fields stay off the JSON file; capsule ids come from `selected`; rejected rows capped; fallback hops recorded
- `runtime/internal/capability/persist.go` — `WriteDecision` refuses a mutated authority write (`ErrDecisionMutated`)
- `runtime/internal/capability/select_test.go` — source scan includes more model-id tokens
- `runtime/internal/engine/mission_bind.go` — capsule lists taken from selected categories; persist returns the mutation error; pre-stage authority check
- `runtime/internal/engine/pipeline.go` — authority check before each stage; persist failure fails the run
- `runtime/internal/engine/phase3e1_test.go` — new
- `docs/3.3.2/ORCHESTRA_3_3_2_CAPABILITY_GATEWAY.md` — evidence authority paragraph
- Vault `Preferences.md` — one 3E.1 line

Not modified: Portfolio, Gate 0, OpenDesign, JCODE install/config, MCP configs, `VERSION`, release tags, `classifyMCP`.

## Tests added

| Test | What it locks |
| --- | --- |
| `TestPhase3E1_DecisionMutationFailsAndFileStays` | remove required `browser.playwright` after write; file unchanged; visual QA does not pass |
| `TestPhase3E1_SideListsDoNotLeakPastSelected` | extra skill/resource/tool slices do not enter the capsule |
| `TestPhase3E1_MCPSnapshotConsumedWithoutExecutor` | healthy tool reaches the capsule; optional broken and auth-required do not; required auth uses `ErrHumanGateRequired` |
| `TestPhase3E1_EvidenceStaysCompact` | evidence keys and hostile byte size |
| `TestPhase3E1_FallbackHopsStayHonest` | multimodal → screenshot → structural hops; no fake vision pass; browser miss stays a human stop |
| `TestPhase3E1_OpaqueLabelsDoNotChangeClass` | `grok-4.7`, `claude-opus`, `gemini-2.5`, `gpt-x`, and an empty label select the same ids |

## Capability immutability

`capability.json` is the existing evidence file, not a second store. `WriteDecision` no-ops when the authority matches and returns `ErrDecisionMutated` when it does not, without rewriting the file. The pipeline checks before the next stage and on persist. The regression removes `browser.playwright` from the in-memory decision: the error is the acceptance failure, the file still contains `browser.playwright`, and visual QA does not `AllPassed`.

`ErrDecisionMutated` reuses the existing acceptance-criteria wording. It is not a new gate type. `Decide` retry rules were not changed.

## Skill / resource / tool boundary

Trace stays candidate ids → `Select` → `Decision.Selected` → `mission.Compile` → builder capsule. The capability package does not open `resources.json`, the skill allowlist, or `mcp.json`. Caps remain skills ≤ 8, resources ≤ 8, mcp tools ≤ 4. Duplicates were already dropped. Unavailable MCP health values are not selected. Appending `FULL-SKILL-CATALOG`, `FULL-RESOURCE-CATALOG`, or `FULL-MCP-CATALOG` onto the side slices does not change the capsule. The capsule lists equal the selected ids for those categories.

## MCP

Injected snapshot, not the live scanner:

- `repo-read` / `HEALTHY` → selected and present on `Capsule.ToolsAllowed`
- `optional-broken` / `BROKEN` → skipped
- `optional-auth` / `AUTH_REQUIRED` → skipped
- required `AUTH_REQUIRED` → `ErrHumanGateRequired` (existing sentinel) and the human detail names the tool

No MCP execution engine. `select.go` does not reference `mcp.json`. `classifyMCP` was not taught to emit `BROKEN`.

## Evidence size

Hostile candidate list (44 skills, 168 resources, 20 MCP ids) against an ordinary mission:

| Artifact | Bytes |
| --- | --- |
| `capability.json` | 5797 |
| worker capsule | 934 |
| dump used as the scale | 152000 |

Top-level keys are only `mission_id`, `selected`, `fallbacks`, `rejected`, `unavailable`, `host_label`. The file does not contain `allowlisted-skill-43`, `catalog-row-167`, or `mcp-catalog-19`. Capsule 934 is the same class as the earlier ~852 byte result. Evidence is larger than the capsule because capped rejected rows are part of the contract, and it stays far under the dump.

## Fallback

With no screenshot adapter: hops `visual.multimodal` → `visual.structural` and `visual.screenshot` → `visual.structural` are recorded. Structural is category `visual.structural` and is not described as multimodal. `evaluator_type` is `structural`. Screenshot invocation count stays 0.

With the availability flag set and still no adapter execution: hop `visual.multimodal` → `visual.screenshot` is recorded, and `ApplyScreenshotEval` does not set `invoked` and does not report multimodal.

`browser.playwright` required and unavailable still returns the human stop (`IsHumanStop` / `ErrHumanGateRequired`). No vision model and no screenshot adapter were added.

## Provider / model neutrality

Same ordinary mission, labels `grok-4.7`, `claude-opus`, `gemini-2.5`, `gpt-x`, and none: the same capability ids. Those strings do not appear in the decision JSON. Non-test capability sources have no `grok`, `gemini`, `gpt-`, `opus`, `sonnet`, or `openai` literal. Host labels `cursor` and `claude` remain the Flow F selection-only preference, not model ids. No model registry and no Cursor picker scrape.

## Chromium regression

`TestLayerASmokePlaywright`: Chromium `153.0.8010.12`, desktop **1440×900** PNG **8751** bytes, mobile **390×844** PNG **5611** bytes, overflow only on mobile. `TestPhase3E_FlowB_ProductUIConsumed` logged the same 8751 / 5611 pair. Layer A owner is still the harness Playwright runner. MCP Playwright is not an authority. Mock screenshots are not an authority.

## Full tests

`go test ./internal/...` from `runtime/`: all packages `ok` (engine 53.382s, visual 5.864s). Exit 0.

`node runtime/tools/check-host-stack.js`: 44 skills, 0 problems.

`node runtime/tools/hygiene.js`: 725 files, VERSION `3.3.1`, 0 failures.

## Remaining limitations

- Live MCP `BROKEN` is still not produced by `classifyMCP`. 3E.1 only accepts an injected snapshot.
- No screenshot adapter and no multimodal adapter are registered on this PC. Those paths stay unexecuted.
- Cursor model discovery is still not in the runtime.
- 1jehuang/jcode is still not installed and is not referenced from the capability package.
- Rejected evidence is capped at 16 rows so a hostile dump cannot fill the file. Rows past that cap are omitted rather than summarized one by one.

## 3F readiness

**Not started.** 3E.1 does not implement the MCP health scanner, registry honesty, or Brain overlay work. The gateway can already consume an injected health snapshot and human-gate a required unavailable tool. That is the handoff 3F would build on. Do not start 3F in this pass.

## Stop

No commit. No push. No tag. No VERSION bump.

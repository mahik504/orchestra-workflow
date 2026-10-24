# Orchestra 3.3.2 — Phase 3F.1 hardening report

**Verdict: PASS WITH LIMITATIONS.**

3F.1 did not redesign 3F. VERSION stays `3.3.1`. No commit, no tag, no 3G.

Brain is an enhancement. It is not a runtime single point of failure.

## Invariants

- Runtime executes.
- Brain remembers.
- Registry verifies.
- Capability Gateway selects.
- Human decides ambiguous or high-impact conflicts.
- Brain cannot override runtime truth.
- Brain failure cannot stop runtime execution.
- Health status cannot redefine mission requirements.
- Historical evidence cannot become current evidence merely through aggregation.
- Provider or model identity does not enter core architecture.

## MCP health and requirement

`HEALTHY` is the only health that can satisfy a tool. `OPTIONAL` and a blank health are unprobed. `Required` on the candidate is a separate flag.

| Case | Result |
| --- | --- |
| Optional + unprobed (`OPTIONAL`) | Skipped. No human gate. Not selected. |
| Required + unprobed | `ErrHumanStop`, mapped by `selectCapabilities` to `ErrHumanGateRequired` |
| Required + blank health | Same human gate |
| Required + `AUTH_REQUIRED` | Same human gate |
| Required + `BROKEN` | Same human gate |
| Required + `DISABLED` | Same human gate |
| Optional + `BROKEN` | Skipped. No human gate. |
| Required + `HEALTHY` | Selected. No human gate. |
| Optional + `HEALTHY` | Selected when the id is relevant. |

`TestPhase3F1_HealthDoesNotDefineRequirement` and `TestPhase3F1_RequiredUnprobedUsesHumanGate` passed. No new MCP manager. `mcp.json` was not edited. No server was spawned. The 3E gateway still selects the minimum set. Only the availability check changed: unprobed no longer counts as available.

## Brain → mission → capsule

`BindCapsule` calls `SelectDurable`, then the existing `mission.Compile`. Compile was not redesigned.

`TestPhase3F1_BrainToCapsule` measured:

| | |
| --- | --- |
| Facts available | 22 |
| Facts selected | 8 (cap) |
| Capsule decisions after duplicate collapse | 1 |
| Stale excluded | 1 |
| Rejected | 14 (unsafe catalogs and dumps, one unrelated project, one live conflict, two over the cap, plus the stale fact) |
| Capsule bytes | 484 |
| Whole-project bytes offered | 136000 |

The capsule did not contain the stale fact, the unrelated project fact, the conflicting Brain sentence, conversation history, the MCP catalog, the skill catalog, the resource catalog, a full Brain blob, a run dump, or grok / claude / gemini. Those inputs were passed on the compile request and recorded as excluded.

## Brain failure

`LoadDurable` returns no facts and no error for unavailable, empty, malformed, stale, and missing payloads. It does not write a fallback store.

`TestPhase3F1_BrainFailureDoesNotBlock` and `TestPhase3F1_BrainMissStillCompiles` passed. After each failure mode the mission kept its acceptance, quality bar, and retry budget. `selectCapabilities` still selected `coding.worker`. `compileBuilderCapsule` still produced a builder capsule.

## Public / private serialization

`SerializePublic` is the gate on text that would leave as a public artifact. `TestPhase3F1_PublicSerialize` passed.

| Input | Result |
| --- | --- |
| Private fixture fact | Rejected |
| Ordinary public sentence | Allowed |
| Same fixture marked public-safe | Allowed |
| `docs/3.3.2/ORCHESTRA_3_3_2_ARCHITECTURE.md` | Kept as a reference |
| Private vault preferences file | Rejected |
| Full private document marker | Rejected |
| Runtime run dump | Rejected |

## Provenance

`Retrieve` copies a fact and does not change status. `TestPhase3F1_ProvenanceStaysPut` passed for unverified, stale, historical, and legacy aggregate. A blank verification date is refused and does not write. `VerifyFact` with `2026-09-22` returns a new current fact and leaves the stored fact unverified.

## Legacy aggregate

`TestPhase3F1_LegacyCannotBecomeCurrent` passed. Payload length is 0. `GetSummary` returns 0. Legacy sum stays 2648 with status `legacy_aggregate_without_payload`. Writing 2648 back into the document total does not make `GetSummary` return 2648. No rows were invented.

## Resource freshness

`TestPhase3F1_RegistryCountsStayPut` passed against `registries/resources.json`: 168 rows, 42 with `last_verified` (all `2026-09-06`), 126 without. A missing date is unverified. Unverified is not verified. No dates were stamped.

## Provider neutrality

`TestPhase3F1_ModelLabelIsNotArchitecture` passed. Facts sourced as model availability for grok-4.7, claude, gemini, and an arbitrary `gpt-x` label produced the same project decision. Those labels did not enter the decision text. `grok-4.7` remains `OPAQUE_LABEL`. No model registry, ranking, or vendor preference was added. `Select` still ignores `DiscoveredModels`.

## Chromium

`TestLayerASmokePlaywright` passed. Chromium 153.0.8010.12. Desktop 1440×900, 8751 bytes. Mobile 390×844, 5611 bytes. Layer A remains the harness. Playwright MCP was not promoted.

## Full regression

| Check | Result |
| --- | --- |
| `go test ./internal/...` | pass (engine 96.9s, visual 6.4s) |
| `check-host-stack.js` | 44 skills, 0 problems |
| `hygiene.js` | VERSION 3.3.1, 0 failures |
| `VERSION` | 3.3.1 |

The memory measurement test was re-run after the suite so the selected-count assertion is included. It passed.

## Scope

Not changed: Portfolio, Gate 0, OpenDesign, JCODE, MCP configuration, VERSION, release tags, 3G, catalog row bodies, host-stack server commands.

Changed: `filterMCP` availability (unprobed is not available), Brain bind/load/serialize helpers, 3F.1 tests, the invariant paragraph in `ORCHESTRA_3_3_2_BRAIN.md`, the 3F.1 acceptance rows, one implementation-plan sentence, and one Preferences line.

## Limitations

- Unprobed servers stay `OPTIONAL`. Nothing in this pass spawned a process to prove `HEALTHY`.
- The 126 catalog rows are still unverified.
- shadcn was not probed and was not repaired.
- OpenDesign stays retired. 1jehuang/jcode stays deferred.
- Brain enrichment is an explicit bind. The pipeline does not read the vault on every run.
- Identical selected sentences collapse inside the existing capsule compiler, so 8 selected facts became 1 decision line in the measurement.

## 3G readiness

3G is not started. The Brain boundary is no longer a reason to hold the test phase: health is separate from requirement, a Brain miss does not stop a run, and 2648 cannot be read as current evidence. This pass does not authorize 3G. Stop here.

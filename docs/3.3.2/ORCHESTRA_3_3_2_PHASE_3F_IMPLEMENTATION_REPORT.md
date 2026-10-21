# Orchestra 3.3.2 — Phase 3F implementation report

**Verdict: PASS WITH LIMITATIONS.**

Brain remembers. Runtime executes. Registry verifies. Human decides ambiguous or high-impact conflicts.

3F is reconciliation and honesty. It is not a memory platform. VERSION stays `3.3.1`. No commit, no tag, no 3G.

## What changed

Runtime state still lives in `.orchestra/runs/<id>/` (`mission.json`, `state.json`, `capability.json`, visual evidence). Brain notes are a capped list of durable facts. A live runtime value wins. A missing catalog date stays unverified. A config name is not an MCP health probe.

The private overlay number 2648 is the sum of 14 per-resource counters. The evaluation payload is empty. That sum is now `legacy_aggregate_without_payload`. Current evidence is the payload length, which is 0. The next recorded evaluation freezes 2648 and does not replace it with the new length.

## Files created

- `docs/3.3.2/ORCHESTRA_3_3_2_BRAIN.md`
- `runtime/internal/memory/reconcile.go`
- `runtime/internal/memory/phase3f_test.go`
- this report

## Files modified

Method repo:

- `runtime/internal/memory/resource_memory.go`
- `runtime/internal/adapters/ag_budget.go`
- `runtime/internal/adapters/ag_budget_test.go`
- `runtime/tools/doctor-ag.js`
- `registries/schemas/resource-memory.schema.json`
- `docs/3.3.2/ORCHESTRA_3_3_2_IMPLEMENTATION_PLAN.md`
- `docs/3.3.2/ORCHESTRA_3_3_2_ARCHITECTURE.md`
- `docs/3.3.2/ORCHESTRA_3_3_2_ACCEPTANCE.md`
- `docs/3.3.2/ORCHESTRA_3_3_2_DECISIONS.md`
- `docs/3.3.2/ORCHESTRA_3_3_2_RECONCILIATION.md`
- `skills/orchestra-conductor/SKILL.md`
- `kit/antigravity/skills/orchestra-conductor/SKILL.md`

Private Brain:

- `routes.md` (3.3.2 architecture pointer)
- `START HERE.md` (same pointer)
- `memory/decisions.md` (one current row; older rows left as written)
- `memory/resource-memory.json` (status fields only)
- `Preferences.md` (one 3F line in each existing Orchestra list)

Not modified: `VERSION`, `mcp.json`, catalog row bodies, host-stack server commands, vault `runtime/`, OpenDesign, JCODE, Portfolio, Gate 0.

## Ownership

`TestPhase3F_OwnershipRuntimeWins` passed.

| Owner | What it keeps |
| --- | --- |
| Runtime | Mission, budget, worker capsule, run directory, capability decision, browser and visual evidence, retries |
| Brain | A lesson plus an optional path reference, after the write allow-list |
| Registry | `last_verified` |
| Human gate | Ambiguous or high-impact conflict (`ErrHumanGateRequired`, unchanged) |

`AuthoritativeRuntime` returns the live value. An empty live value does not fall through to a Brain string. A summary that contains `CAPSULE_DUMP`, `RETRY_LOG`, `PNG_BYTES`, or a multi-line run listing is rejected. A one-line lesson with a path reference is accepted.

## Read / write

`TestPhase3F_AdmitDurable` passed.

Admitted: explicit user preference, finalized architecture decision, workflow validated more than once, reusable failure lesson, verified resource, durable project fact.

Refused: chat turn, tool call, retry, transient error, screenshot, temporary model availability, one successful command.

`SelectDurable` reads an already-loaded list. Cap 8. It does not open the vault and it does not replace `mission.Compile`.

## Provenance and conflicts

`TestPhase3F_ConflictOrder` and `TestPhase3F_SelectCapAndStale` passed.

Winner order:

1. Live runtime
2. Current docs
3. Verified registry
4. Brain history
5. Unverified memory

A fact whose id is already in the live map is not retrieved. Stale, superseded, and `legacy_aggregate_without_payload` facts are excluded unless history was requested. An unverified fact does not count as verified.

## Public / private

`TestPhase3F_PublicPrivate` passed.

The fixture marker `PRIVATE_BRAIN_FIXTURE` is rejected from a public artifact unless that text is explicitly marked public-safe. A path such as `docs/3.3.2/ORCHESTRA_3_3_2_ARCHITECTURE.md` is a reference. A blob marked `AGENTS.md FULL COPY` is a copy and is detected. The private ledger was not copied into the method seed. Tests do not contain real preference or career text.

## Duplicate surfaces

`TestPhase3F_DuplicateMatrixAndStance` passed. Actions are only KEEP or MIGRATE. No DELETE, no ARCHIVE, no file removal.

| Surface | Class | Action |
| --- | --- | --- |
| PlaywrightVerifier | CANONICAL | KEEP |
| Catalog playwright | ADAPTER | KEEP |
| Host-stack Playwright MCP | ADAPTER | KEEP (not Layer A) |
| playwright-cli | ADAPTER | KEEP |
| browser-use, agent-browser | OPTIONAL | KEEP |
| Context7 MCP | CANONICAL | KEEP |
| Second Context7 skill-list label | DUPLICATE | MIGRATE the note to the one MCP |
| gsap-core skill | CANONICAL | KEEP |
| Catalog gsap | ADAPTER | KEEP |
| motion-dev, motion-primitives, motionsites-ai | OPTIONAL | KEEP |
| shadcn MCP | OPTIONAL | KEEP |
| Catalog shadcn-ui | OPTIONAL | KEEP |
| GitHub host-stack plugin | OPTIONAL | KEEP, no new catalog row |

## Stance

| Item | Status |
| --- | --- |
| shadcn MCP | OPTIONAL. Not repaired. Not `BROKEN` without a probe. |
| OpenDesign + Antigravity | RETIRED |
| 1jehuang/jcode | DEFER. No install, no exec, no Go import. |
| grok-4.7 | OPAQUE_LABEL. Not architecture, not a registry, not current availability. |

## Legacy aggregate vs current evidence

`TestPhase3F_LegacyAggregateSplit` passed.

| View | Value |
| --- | --- |
| Private `total_evaluations` | 2648 (historical document field, not rewritten to 0) |
| `aggregate_sum` | 2648 |
| `evaluation_payload` | absent |
| `status` | `legacy_aggregate_without_payload` |
| `evaluations` length | 0 |
| `GetSummary` / `CurrentEvidenceCount` on that file | 0 |
| After one new `Record` | legacy sum still 2648; current evidence 1 |

No evaluation rows were invented. The public method seed was not given this ledger.

## MCP

`TestClassifyMCP_UnprobedIsOptional` and `TestClassifyMCP_ProbeStates` passed.

| Probe | Health |
| --- | --- |
| None | OPTIONAL (playwright, context7, stitch, vault-memory, supabase, github included) |
| Observed + reachable | HEALTHY |
| Observed, nothing else | OPTIONAL |
| Missing credentials | AUTH_REQUIRED |
| Process missing / refused / invalid | BROKEN |
| Turned off | DISABLED |

`CheckAntigravityBudget` and `doctor-ag.js` pass a nil probe. They do not read `mcp.json` and they do not spawn servers. The capability gateway was not rewritten.

## Registry freshness

`TestPhase3F_MissingLastVerifiedIsUnverified` passed against `registries/resources.json`.

| | Count |
| --- | --- |
| Rows | 168 |
| `last_verified` present | 42 (all `2026-09-06`) |
| Missing date | 126 |

A blank date is unverified. No row was stamped.

## Context efficiency

`TestPhase3F_ContextBytes`:

| | |
| --- | --- |
| Selected facts | 8 |
| Selected bytes | 2225 |
| Rejected | 23 |
| Stale excluded | 1 |
| Cap | 8 |
| 3C hostile dump scale | ~136000 bytes |

2225 is well under a tenth of that dump.

## Non-regression

| Check | Result |
| --- | --- |
| `go test ./internal/...` | pass (engine 62.7s, visual 8.1s, memory and capability included) |
| `TestLayerASmokePlaywright` | Chromium 153.0.8010.12 |
| Desktop PNG | 1440×900, 8751 bytes, real PNG |
| Mobile PNG | 390×844, 5611 bytes, real PNG |
| `check-host-stack.js` | 44 skills, 0 problems |
| `hygiene.js` | VERSION 3.3.1, 0 failures |
| `VERSION` | 3.3.1 |

3E selection, provider neutrality, capability evidence, fallback, and the human gate were not rewritten. Their packages passed in the full suite.

## Limitations

- The 126 unverified catalog rows were not checked and were not dated.
- No live process probe. Unprobed servers stay OPTIONAL, including ones that used to be stamped HEALTHY by name.
- shadcn was not installed and was not marked BROKEN.
- OpenDesign was not repaired. 1jehuang/jcode was not installed.
- Brain retrieval is an in-memory helper. There is no index, embedding, or second store.
- 3G is not started.

## Stop

Phase 3F is ready for review. No commit, no push, no tag, no VERSION bump.

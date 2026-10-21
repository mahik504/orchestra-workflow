# Orchestra 3.3.2 — Brain and overlay

**Brain remembers. Runtime executes. Registry verifies. Human decides ambiguous or high-impact conflicts.**

3F reconciles the private Brain with the 3.3.2 runtime. It does not add a second engine, a memory database, or a model registry.

## Ownership

| Thing | Owner |
| --- | --- |
| Mission, budget, worker capsule, `.orchestra/runs/<id>/`, capability decision, browser and visual evidence, retries | Runtime |
| Durable preference, architecture decision, project fact, validated workflow, failure lesson, verified resource note | Brain |
| `last_verified` on a catalog row | Registry |
| Ambiguous or high-impact conflict | Existing human gate (`ErrHumanGateRequired`) |

A Brain note may cite a run path. It does not store the mission file, the capsule, a screenshot, or the retry log. Live runtime state wins over every note.

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

MCP health (`HEALTHY`, `OPTIONAL`, `AUTH_REQUIRED`, `BROKEN`, `DISABLED`) is not the mission's required flag. `OPTIONAL` and a blank health are unprobed. Only `HEALTHY` can satisfy a tool. A required tool that is unprobed, `AUTH_REQUIRED`, `BROKEN`, or `DISABLED` uses the existing human gate. An optional unprobed or broken tool is skipped. Retrieval does not upgrade provenance. An explicit verification date is the only status change.

Conflict order, highest first:

1. Live runtime state
2. Current `docs/3.3.2` and the VERSION-era code
3. A registry row that has `last_verified`
4. Brain historical knowledge, including a legacy aggregate
5. Unverified memory

A preference never overrides one conductor, the Layer A owner, the human gate, or the rule that runtime state stays in the run directory.

## Memory classes

`preference`, `architecture_decision`, `project_fact`, `workflow`, `failure_lesson`, `resource_knowledge`, `evaluation_history`.

Provenance on a durable fact: `source`, `as_of`, `status` (`current` | `historical` | `stale` | `unverified` | `legacy_aggregate_without_payload` | `superseded`), `scope`, `last_verified` when someone actually checked, `superseded_by`.

`evaluation_history` may cite the legacy aggregate. It is not current runtime state.

## Write

Admit a durable note only for: an explicit user preference, a finalized architecture decision, a workflow validated more than once, a reusable failure lesson, a resource with a real `last_verified`, or a durable project fact.

Do not admit: a chat turn, a tool call, a retry, a transient error, a screenshot, temporary model availability, or one successful command. Do not infer a preference.

## Read

`SelectDurable` reads an already-loaded candidate list. Cap 8. Drop stale, superseded, and legacy-aggregate facts unless history was requested. Scope filter applies. A fact whose id is already in the live map is rejected.

Never open the vault from the runtime. Never pass `Preferences.md` whole. Never substitute a Brain blob for `mission.Compile`.

## Overlay 2648

The private ledger's `total_evaluations: 2648` equals the sum of the per-resource counters. `evaluations` is empty. The rows are not reconstructable.

That sum is `legacy_aggregate_without_payload`. Current APIs (`GetSummary`, `CurrentEvidenceCount`) use the payload length, which is 0 until a new evaluation is stored. The next `Record` freezes the historical sum in `aggregate_sum` and does not replace it with the new payload length. Do not invent rows. Do not set the historical sum to 0.

The public method seed stays empty. Do not copy the private ledger into the method repo.

## MCP health

States: `HEALTHY`, `OPTIONAL`, `AUTH_REQUIRED`, `BROKEN`, `DISABLED`.

A config name is not a probe. Unprobed is `OPTIONAL`, including Playwright MCP and a server named supabase. `HEALTHY`, `AUTH_REQUIRED`, `BROKEN`, and `DISABLED` come only from an injected probe. `classifyMCP` does not read `mcp.json`, does not spawn a process, and is not a new MCP manager. Layer A stays `PlaywrightVerifier`.

## Registry freshness

A missing `last_verified` is unverified. Do not stamp the unverified catalog rows in this phase.

## Duplicate surfaces

| Surface | Class | Action |
| --- | --- | --- |
| `PlaywrightVerifier` | CANONICAL | KEEP |
| Catalog `playwright` | ADAPTER | KEEP |
| Host-stack Playwright MCP | ADAPTER | KEEP (not Layer A) |
| `playwright-cli` | ADAPTER | KEEP |
| `browser-use`, `agent-browser` | OPTIONAL | KEEP (do not promote) |
| Context7 MCP | CANONICAL | KEEP |
| Second "Context7" label in a host skill list | DUPLICATE | MIGRATE the note to the one MCP. Do not delete the server. Do not invent a catalog row. |
| `gsap-core` skill | CANONICAL | KEEP (motion implementation) |
| Catalog `gsap` | ADAPTER | KEEP |
| `motion-dev`, `motion-primitives`, `motionsites-ai` | OPTIONAL | KEEP |
| shadcn MCP | OPTIONAL | KEEP. Not core. Not repaired. Not `BROKEN` without a probe. |
| Catalog `shadcn-ui` | OPTIONAL | KEEP |
| GitHub host-stack plugin | OPTIONAL | KEEP. No catalog row. Not `HEALTHY` by name. |

No `DELETE` or `ARCHIVE` in 3F. No speculative cleanup.

## Stance

| Item | 3F status |
| --- | --- |
| shadcn MCP | OPTIONAL. Freeze "broken" was not re-proven. |
| OpenDesign + Antigravity | Retired. Do not repair. |
| 1jehuang/jcode | DEFER. No install, no exec, no Go import. |
| Grok 4.7 and other opaque model labels | Not architecture, not a registry, not current availability. Same capability class selects the same ids. |

## Public and private

Public artifacts get this policy, the honesty code, and tests. Private Brain gets the stale pointer fixes, one current decision row, and the status fields on its ledger.

Do not copy private notes, preference bodies, or the 2648 ledger into the method repo. Tests use the fixture marker `PRIVATE_BRAIN_FIXTURE`, not real private text. A public architecture id may be referenced by path. Do not paste `AGENTS.md`, `ARCHITECTURE.md`, or `resources.json` into a Brain summary.

## Retrieval size

A capped select of 8 facts must stay far under a vault or project dump (the 3C hostile dump is about 136000 bytes). Log selected count, selected bytes, rejected count, and stale excluded. Do not chase a magic token constant.

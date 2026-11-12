# Architecture — Orchestra 4.7.0

The contract pin is 4.6.0. The engine keeps a mission, a worker capsule, a derived budget, and a browser check on a running app. Those sit around the eight stages. They do not replace them. A model name does not pick a skill. 4.1 adds the selection receipt, design domains, and technology rules on the same loop. 4.2 adds five skills and keeps the same conductor.

Orchestra is a **control plane**. Your agent is the executor.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

The **registry** is the resource catalog. The **graph** is capability routes (17 rows, each with trigger and skip). Host rules and IDE customizations are adapters. They never implement a competing orchestration policy.

## The eight stages

```mermaid
flowchart TD
  discover[1 Discover: inspect workspace and existing tooling]
  classify[2 Classify: archetype, quality bar, platform]
  research[3 Research: walk the capability graph]
  synth[4 Synthesize: tokens, type, motion, rules]
  design[5 Design System: DESIGN.md contract then golden stills]
  gate{Human gate}
  implement[6 Implement: acquire scoped resources, write code]
  qa[7 Visual QA: multi-viewport, console, contrast]
  iterate[8 Iterate: route defects back]

  discover --> classify --> research --> synth --> design --> gate
  gate -->|approved| implement --> qa --> iterate
  gate -->|edit| research
  iterate -->|layout defect| implement
  iterate -->|token defect| design
```

Backend, research, and documentation tasks run the same spine with stages 4–5 collapsed into a technical plan and stage 7 replaced by tests and static analysis.

## Mission execution (3C)

The eight stages stay the executor. Around them the engine now keeps a compact, provider-neutral control plane:

```
MISSION → CONTEXT SELECTION → WORKER CAPSULE → BUDGET → EXECUTION → EVIDENCE / RESULT → FAILURE CLASSIFICATION → RECOVERY OR HUMAN GATE
```

- **Mission** — run-scoped JSON (`runtime/internal/mission`). Objective, quality bar, LIGHT/FULL intensity, immutable acceptance, derived execution budget. Not a Brain dump.
- **Worker capsule** — one assignment: role (planner/builder/judge), named files, skill ids, constraints, acceptance, prior failures. No conversation history, MCP catalog, or unrelated skills.
- **Budget** — derived per task (`reasoning`, `context`, `research`, `visual_review`, `independent_judge`, `retry_budget` 0–4). Ordinary API/debug stays cheap. **3D** consumes `visual_review` / `independent_judge`: harness Layer A/B when budgeted; independent judge only when policy triggers.
- **Failure class** — environment, dependency, provider/tool, implementation, test, specification, human-needed, timeout/resource. Non-zero exit is not automatically “implementation.”
- **Recovery** — retry only with remaining budget **and** a change (input/context/strategy/tool/approach/environment). Identical signatures stop. Specification and human-needed never auto-retry.
- **Durable state** — `.orchestra/runs/<id>/mission.json` + `state.json` (FULL also writes `progress.txt`). Enough to continue without replaying chat. Existing `.orchestra/handoff/state.json` is unchanged.

Quality bars remain STANDARD / PREMIUM / EXPERIMENTAL. Intensity is not a fourth bar.

## Five subsystems

1. **Registries** (`registries/`) — machine-readable resource catalog and the capability graph, validated by JSON Schema. JSON is canonical.
2. **Engine** (`runtime/internal/engine`) — the eight-stage state machine, plus the 3C mission/capsule/budget/failure bind (`runtime/internal/mission`).
3. **Acquisition adapters** (`runtime/internal/adapters/acquisition`) — npm, git, CLI, MCP, and web reference fetch. Scoped `GLOBAL` / `PROJECT` / `ON_DEMAND`. Global package installs are blocked.
4. **Host adapters** (`runtime/internal/adapters`) — keep Cursor, Antigravity, and Claude Code on one capability contract.
5. **Resource memory** — durable outcomes per resource. Only real executed jobs write rows.

## Gates

| Gate | Rule |
| --- | --- |
| Design Lab | `PREMIUM` / `EXPERIMENTAL` visual work blocks product frontend writes until `CONTRACT_APPROVED` and `ApproveStills`. Golden stills are pre-implementation evidence. `STANDARD` is exempt unless requested. |
| Real-app verification | Launch and exercise the application. Source inspection is not verification. |
| Correctness review | Bugs, missing behavior, security. |
| Simplify review | Duplication, needless abstraction, dead complexity. Separate pass. |
| Evidence-first | `DONE` / `VERIFIED` / `PASSED` require observed output in the same message. |
| Human-needed (3C) | Stop automatic execution on ambiguous product behavior, missing credentials, destructive actions, spec/contract conflict, uninterpretable acceptance, identity/brand, stills/direction pick, security-sensitive change, irreversible deploy, material tool/host gap, or repeated bounded failure. `ErrHumanGateRequired` is the existing sentinel. |

## Resource lifecycle

A run may move one step at a time:

```
discovered → installed → available → selected → executed → verified → successful
```

The older names still describe the same idea: acquired is installed, used is executed. A row in `registries/resources.json` stays **discovered**. Projecting a contract does not mark it executed. `Promote` refuses a jump from discovered to successful.

Presence in a registry is not usage. A resource counts as used only when a run recorded evidence (`run:` or `evidence:`).

## Selection receipt (4.1)

When a task loads skills, an MCP, a plugin, or a library, the conductor names those and leaves the other lines off. Design domains annotate the capability the classifier already chose. Technology rules keep one scroll owner, one HTTP stack, and one bundler. Redis waits for a cache, TTL, or rate-limit requirement. The pieces are `runtime/internal/routing`, `designdomain`, `techpolicy`, `review`, and `contractpack`. Protocols: `protocols/EXPLAINABLE_ROUTING.md`, `protocols/TECHNOLOGY_SELECTION.md`.

A design pack is `templates/contracts/`. A bug, a typo, or one control does not generate it. Review lenses are chosen per task. A visual judge needs a screenshot path and cannot repeat the builder. Comparison tasks for six product classes are specified and disabled. They do not score this release.

## Quality bars

- `STANDARD` — efficient product work. No slop, no lab overhead.
- `PREMIUM` — multi-source research, Design Lab, visual QA, iteration.
- `EXPERIMENTAL` — 3D, shaders, WebGL, novel interaction. Always with a low-end fallback.

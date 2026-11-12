# Performance and context 5.5

Recorded 2026-11-12. Bytes are measurements of this tree. Tokens and milliseconds are estimates or NOT_RUN. Nothing here is a claim that Orchestra is faster or cheaper than another harness.

## What was measured

PowerShell `Get-Item` length and `Get-Content -Raw` length. On this tree the two numbers matched, so the files are a single byte per character (no UTF-16, no extra CR in the count).

| File | Bytes |
| --- | ---: |
| `AGENTS.md` | 2049 |
| `CLAUDE.md` | 1242 |
| `.cursorrules` | 1117 |
| `kit/antigravity/GEMINI.md` | 1138 |
| `kit/antigravity/ALWAYS-ON.md` | 940 |
| `kit/codex/AGENTS.md` | 552 |
| `kit/opencode/AGENTS.md` | 558 |
| `kit/kimi/AGENTS.md` | 550 |
| `docs/ROUTING.md` | 4010 |
| `docs/CONDUCTOR.md` | 19467 |
| `ARCHITECTURE.md` | 6797 |

`docs/CONTEXT-BUDGET.md` still says `AGENTS.md` is 1,847 characters and all canonical skill descriptions are 25,614 characters, with `ow-*` descriptions at 7,303 and none over 300. Those three skill figures were not re-measured. `runtime/tools/measure-context.mjs` writes `docs/CONTEXT-BUDGET.md`. Re-running it would dirty a product doc. The 1,847 figure is stale. The skill totals are REPORTED from that doc.

This Cursor session’s workspace rules included `AGENTS.md` (2,049 bytes), `CLAUDE.md` (1,242), and `.cursorrules` (1,117). The only rule file with `alwaysApply: true` is `.cursor/rules/000-core.mdc`, measured at 543 bytes. The other rule files attach by glob or stay requestable. A host that injects every adapter adds the sizes together. A host that injects one adapter does not.

`docs/CONDUCTOR.md` is the long procedure `skills/ow-conductor/SKILL.md` points at, and its title is still Orchestra 4.6.0. If a model follows that pointer, the file is about 19 KB before any skill body. A rough token estimate is length divided by 4, about 4,900 tokens. That divisor is a convention, not a tokenizer run. Label it ESTIMATE.

## What the compiler does

`compile.Compile` walks 36 routes and scores keyword lengths. That work is in-process string matching. Wall-clock time was NOT_RUN. There is no model call inside the compiler. The model call happens when the host decides to read the selected `SKILL.md`.

Disclosure levels in `docs/CONTEXT-BUDGET.md`:

| Level | What the worker is supposed to see |
| --- | --- |
| L0 | Task and acceptance |
| L1 | The route row |
| L2 | The selected `SKILL.md` |
| L3 | `references/` beside that skill |
| L4 | Tool output from the run |

The budget doc says a STANDARD job that names the file stays at L0, and the conductor fast path does not print a brief. `docs/CONDUCTOR.md` says the same for mode `fast`. `compile.go` sets `Disclosure` from the winning route. It does not look for a file path in the task. A bug route is `L0` because the route says so, not because the compiler noticed a path. A task that matches no keyword is `L0` as well, with confidence left at 1 (`BASELINE` F-01). The fast path is prose plus a route field. It is not a separate code path.

## Where time and tokens go

Observed structure, not a trace:

- **Always-on text.** The adapter files above, plus `000-core`, plus whatever the host adds. Duplicating the same contract into `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, and `GEMINI.md` costs maintenance on every edit. It costs tokens only when one session loads more than one of them. This session did.
- **The conductor document.** 19,467 bytes, stale pin. This is the largest single instruction file measured.
- **Stub skills.** A `stub: true` file is small. The cost is a model that then improvises, or a model that opens `docs/CONDUCTOR.md` because the stub points there. Small files are not the waste. The pointer to the stale long doc is.
- **Skill description catalog.** If a host loads every description (the budget doc’s 25,614 characters), that is a standing tax. The route design is supposed to load one row. Whether Cursor’s skill list injects every description is a host behavior. UNKNOWN for a fresh session. This session’s skill list in the system prompt was large. That list is the host’s installed skills, which include Orchestra skills and others. It was not counted.
- **Subagents.** This research session launched five. Three returned. They are optional. A bugfix that names a file should not launch them. The compiler does not forbid it.
- **Design pipeline.** `orchestra plan` with an empty task invents a frontend sentence (`main.go` lines 403–405) and builds the eight-stage pipeline. That is a latency trap for a missing argument.
- **Verification.** `orchestra verify` and visual QA are the right cost when a UI changed. They are the wrong cost when the verdict is `NOT_APPLICABLE` and a memory row is still written as success (F-06). The time cost of a skipped check is small. The false lesson in memory is the real cost.
- **Web acquire fallback.** A network error that returns an `offline_fixture` saves a retry and spends a later turn on fake content. That is a latency win and a correctness loss (F-05).

## Model class

Routes set `model_class` (`coding-fast`, `coding-precise`, `reasoning-high`, `frontend-design-high`, `security-high`, `documentation`, `research-web`). The host chooses the model. Orchestra cannot switch Gemini Flash and Claude Opus unless that host exposes a switch. `model_class` is a hint on the manifest. Treating it as a router would claim a control the binary does not have.

A deterministic script is the right tool when the check is a test, a schema, or a hash. `PromoteSkill`, `check-versions.js`, and the hook denylist are that kind of tool. A second model is justified when the task is a review that must not repeat the builder. The contract already says the visual judge cannot be the builder. No trial here shows that a second model pays for itself.

## Fast path to keep

A task that names a file and is a bug, a typo, or one control:

- Do not write a brief.
- Do not open the design lab.
- Do not load `docs/CONDUCTOR.md`.
- Load the one skill the route names, or none for “one control”.
- Run the test the task names.
- Skip the review panel unless the change is security-sensitive or the user asked.

`compile` should emit `fast: true` and confidence below 1 when nothing matched, so a host that only reads the JSON does not invent a premium loop. That emit is a recommendation. It is not in the binary today.

## What was not measured

Routing latency, retrieval latency, time to first useful token, tool-call counts, model-call counts, and dollar cost. Any number for those would be invented. The eval document says how to measure them. The result column is `not_run`.

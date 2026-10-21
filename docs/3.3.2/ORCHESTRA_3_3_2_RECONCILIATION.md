# Orchestra 3.3.2 — reconciliation

Primary Phase 3A document. Maps external mechanisms and internal freeze items onto KEEP / ADAPT / COMBINE / DEFER / REJECT / FIX IN 3.3.2 / REMOVE / REPLACE / INVESTIGATE.

Phase **3A.1** (docs only): FULL is not a synonym for independent judge; visual QA is two-layer; 3–5 directions are lightweight.

Scores 0–5: **Out** real output · **Des** design/creative · **Rel** reliability · **Ctx** context efficiency · **Spd** speed · **Rec** recovery · **Port** portability · **Mnt** maintainability · **Cx** complexity cost (higher = worse) · **Dup** duplication risk (higher = worse).

---

## Freeze items (current Orchestra)

| Item | Classification | Why |
| --- | --- | --- |
| One conductor | **KEEP** | Control plane. Not a swarm. |
| Control plane vs capabilities vs executors vs Brain vs registry | **KEEP** | Correct split. |
| STANDARD / PREMIUM / EXPERIMENTAL bars | **KEEP** | Do not add a fourth **bar** named LIGHT. Intensity is orthogonal. |
| Two-lock Design Lab (`CONTRACT_APPROVED` then `ApproveStills`) | **KEEP** | Implemented in `runtime/internal/verify/design_lab.go`. |
| Golden stills as pre-implementation visual evidence | **KEEP** | Not a substitute for live-app QA. |
| Real-app Playwright policy | **KEEP** as policy · **FIX IN 3.3.2** as Layer A | `visual_qa.go` can write `MOCK_SCREENSHOT_*`. Playwright is not the art director. |
| Keep / Rework / Kill / Park | **KEEP** | Human judgment. |
| Discover broadly / activate selectively | **KEEP** | Registry posture. |
| Job-loaded skills (`visual-forensics`, `visual-golden-still`, `fresh-context-review`, `final-showcase`) | **KEEP** · **FIX IN 3.3.2** host parity | In-repo and on `sync-ides` allowlist (44). Live hosts 40. |
| Named visual refs → exactly 3 translations, not skip | **KEEP** in engine/AGENTS · **FIX IN 3.3.2** in `.cursorrules` / README / adapters | AGENTS reversed 3.3.0 skip; docs lag. |
| Pasted DESIGN.md skips survey only | **KEEP** | Documented in AGENTS 3.3.1. |
| `kit/sync-ides.ps1` as canonical 44-skill install | **KEEP** | Do not run it in 3A. |
| `kit/antigravity/install.ps1` hardcoded 30 names | **FIX IN 3.3.2** | Must consume the same allowlist. |
| Live conductor skills still 3.3.0 | **FIX IN 3.3.2** | Repo 3.3.1 vs hosts 3.3.0. |
| Brain `START HERE.md` / `routes.md` still “3.2 contract” | **DONE in 3F** | Pointers now name `docs/3.3.2/ORCHESTRA_3_3_2_ARCHITECTURE.md`. Overlay stays overlay. |
| Catalog 168 / 42 verified / 126 unverified | **FIX IN 3.3.2** honesty · **DEFER** bulk verification | Do not inflate `last_verified`. |
| Brain `resource-memory.json` 2648 vs empty `evaluations` | **DONE in 3F** | Sum of the per-resource counters. Payload absent. Status `legacy_aggregate_without_payload`. Current evidence is the payload length. |
| GitHub MCP optional, no catalog row | **DONE in 3F** | Stays a host-stack plugin. No new catalog row. A name is not `HEALTHY`. |
| Duplicate Playwright / Context7 / motion surfaces | **DONE in 3F** classification | One Layer A owner: `PlaywrightVerifier`. Other surfaces KEEP. The extra Context7 label MIGRATE as a note. No deletes. |
| Magic UI OPTIONAL vs engine ban of always-on `magicui` | **KEEP** ban · **FIX IN 3.3.2** docs/catalog wording | Same posture as 3.3.1. |
| shadcn MCP broken at freeze | **INVESTIGATE** | Repair or mark broken; do not pretend it works. |
| OpenDesign as second studio vs Stitch | **KEEP banned** | Optional isolated lab only. Claude primary; Cursor Agent `cursor-grok-4.6-high` secondary. OD×Antigravity **retired**. |
| OpenDesign wired into Orchestra | **REJECT** | Not in 3.3.2 core. |
| `ARCHITECTURE.md` “stack card” gate | **FIX IN 3.3.2** | Stale vs stills lock. |
| 8-stage engine | **KEEP** | Add mission/capsule/failure around it; do not replace. |
| `ErrHumanGateRequired` | **KEEP** · **FIX IN 3.3.2** specify when | Identity, direction, destructive, ambiguous scope, security, irreversible deploy, missing acceptance. |
| GSD / ECC dumps / Karpathy packs / always-on Magic UI as core | **KEEP REJECTED** | `docs/3.3.1/CAPABILITY_AUDIT.md`. |
| Git tag still `v3.3.0` | **DEFER** | Tag 3.3.2 only after operator says so. No 3A tag. |
| VERSION bump to 3.3.2 | **DEFER** | Not this phase. |

---

## External mechanisms

| Source | Mechanism | Problem it hits | Evidence | Output | Cx/Dup | Decision |
| --- | --- | --- | --- | --- | --- | --- |
| Anthropic 2026 harness | Independent evaluator; Playwright on running app | Builder grades itself; mock screenshots | Anthropic article; cwc repo; later **removed** always-on eval | High when **triggered**, not on every FULL | Cx 4 if always-on | **ADAPT** independence + live Layer A. **REJECT** always-on 3-agent theatre and FULL=judge. |
| Execution budget (derived) | reasoning/context/research/visual_review/judge/retry | One-size FULL lane; waste | Anthropic eval simplification; LIGHT vs FULL | High cost control | Cx 1 | **ADAPT** as routing policy. **REJECT** hardcoded presets. |
| Two-layer visual QA | Layer A Playwright + Layer B perception | Playwright as taste; mock QA | Canary (harness owns A); screenshot eval (B) | High visual | Cx 2 if one browser | **ADAPT**. **REJECT** second browser stack. |
| Anthropic 2025 harness | Initializer + `progress.txt` + git | Lost context on long runs | Official article | High on FULL | Cx 2 | **ADAPT** for FULL only. |
| cwc-long-running-agents | Evaluator has no write tools | Same-model rubber stamp | Public repo | High | Cx 2 | **ADAPT** |
| Codex / Canary Lab | Harness owns Playwright **Layer A** verdict | Agent self-greens tests | Canary README; Codex harness | High reliability | Cx 2 | **ADAPT** Layer A ownership. Layer B is separate. |
| Tenet | `compile_context` capsule | Workers get Brain dumps | Tenet MCP surface | High context | Cx 5 if DAG imported | **ADAPT** capsule. **REJECT** DAG OS / SQLite jobs. |
| Tenet | Independent critics + retry classes | Blind loops | Tenet heartbeat/orphan | High recovery | Cx 3 | **ADAPT** taxonomy + budget. |
| Spec Kit | Compact spec → plan → tasks | Vague STANDARD+ work | GitHub spec-kit | Medium | Dup 4 if second conductor | **ADAPT** artifacts. **REJECT** slash zoo / constitution. |
| Deep Suite | Skippable brainstorm; independent review | Forced ceremony on glue | deep-loop skip flags | Medium | Dup 4 | **COMBINE** skip + existing `fresh-context-review`. **REJECT** plugin suite. |
| Ralph | Fresh instance per story; max iterations; stuck | Premature “done”; infinite loops | snarktank + zenith comparison | Medium FULL | Cx 3 | **ADAPT** budget + stuck. **REJECT** unattended PRD default. |
| Zenith | Stop when original gaps closed | Harness that never ends | Zenith study (not re-run) | Medium | Cx 5 if imported | **ADAPT** stop rule. **REJECT** multi-day MCP OS. |
| Nonstop | Immutable acceptance board | Mid-run rewriting DoD | seolcoding repo | Medium | Dup 3 | **ADAPT** (mission file). |
| frontend-design skill | One signature; 2026 clichés as calibration | Inter/#000 greps; cream-serif default | anthropics/skills SKILL.md | High design | Dup 2 | **ADAPT** as critic. **REJECT** skill dump. |
| no-slop-ui / ui-craft | Structural questions + optional scanners | Taste adjectives; false bans | GitHub READMEs | Medium design | Dup 3 | **ADAPT** questions. **REJECT** banned-pattern law. |
| Get Design Done | Verify vs brief/tokens | Design Lab overlap | GDD repos | Medium | Dup 4 | **ADAPT** verify-vs-brief. **REJECT** GDD install. |
| twaldin/harness | RunSpec → RunResult | Silent bad provider pick | CLI wrapper | Medium routing | Cx 4 if 13 CLIs | **ADAPT** capability→runtime rec. **REJECT** wrapping banned CLIs. |
| minhtran lanes | tiny / normal / high-risk | One intensity for all | Hook-enforced lanes | High speed on LIGHT | Cx 1 | **ADAPT** as LIGHT vs FULL. |
| cc-sdd / SDD workers | Direct implement when small; clean worker context | SDD zoo on CRUD | gotalab/cc-sdd | High speed LIGHT | Dup 3 | **ADAPT** clean workers. **REJECT** `/sdd-*`. |
| HarnessXP | Cross-repo contract sync | API drift | mebusw/HarnessXP | Low today | Cx 4 | **DEFER** until multi-repo products exist. |
| Memory harnesses | Measured memory; provenance | Fake 2648 evals | kenhuangus + others | Medium | Cx 4 | **DEFER** plugins. **KEEP** explicit Brain writes + classes. |
| Spexcode | Spec nodes + real-surface eval | Extra product | spexcode.net | Medium later | Cx 5 | **DEFER** |
| ShioriCode | Session GUI / reconnect | Host already owns this | shiorihq | Low | Dup 5 | **REJECT** as Orchestra product. |
| Interactive preview IDEs | Click-to-source | Human visual speed | OpenSail / Kombai | Medium later | Cx 4 | **DEFER**. KEEP stills + Playwright. |
| Pixel-match / Figma loops | SSIM to mock | Clone risk (Elaine Yu bar) | Imugi / pixelprobe | Harmful if DoD | Cx 3 | **REJECT** as quality bar. |
| Agent Harness “v10” | Name collision | Cannot vendor a ghost | OLAV / HarnessAudit / others | n/a | Dup 5 | **INVESTIGATE** only. **REJECT** import. |
| GSD / ECC / Karpathy packs | Ceremony + dumps | Already rejected 3.3.1 | CAPABILITY_AUDIT | Negative | Cx 5 | **KEEP REJECTED** |
| Agentic Harness layers | Identity/rules/skills/hooks/memory | Orchestra already has this | sevenschulte | Low delta | Dup 4 | **KEEP** existing split. **REJECT** install. |

---

## Complexity budget (3.3.2)

Spend complexity only where output moves. Prefer adaptive routing.

Every adopted mechanism must answer: (1) output problem, (2) evidence of improvement, (3) replace or complement, (4) operational cost, (5) failure mode introduced, (6) removable if stronger models make it unnecessary. Full answers: [ORCHESTRA_3_3_2_ACCEPTANCE.md](ORCHESTRA_3_3_2_ACCEPTANCE.md) mechanism matrix.

| Budget item | Allowed | Not allowed |
| --- | --- | --- |
| New runtime concepts | Mission, capsule, failure class, evaluator policy, intensity LIGHT/FULL, **execution budget** | DAG, SQLite job OS, fourth quality bar, model zoo, one-size budget presets |
| New agents | Roles (planner/builder/judge). Independent judge **conditional on triggers** — FULL is not sufficient | Always-on three-agent GAN; FULL = judge |
| New skills | None as always-on. Job-load existing four; maybe inject critic heuristics into `visual-forensics` later | Skill packs, Magic UI core, GDD/GSD commands, always-on critics |
| New MCP | Collapse duplicates. Optional OD **not** in host-stack by default | Catalog inflation, GitHub MCP as required, second browser |
| New docs | This folder + conductor/docs parity in 3B | Parallel constitutions, second AGENTS |
| Cross-repo | Lightweight contract mention only | HarnessXP-style fleet sync |
| Visual | One Playwright owner (Layer A); Layer B on A’s screenshots; stills stay lock 2; lightweight 3–5 directions then **one** build | Second browser MCP; pixel-clone DoD; mock screenshots as pass; Playwright-as-taste; 3–5 implementations |
| Research / eval | Token ladder; Layer B / independent judge when budgeted | Always-on research; always-on Layer B; always-on judge |

If a change does not move A–I product output **without unjustified cost** (see ACCEPTANCE), it is out of budget.

### Six questions — adopted mechanisms (short)

| Mechanism | Problem | Evidence | Complements / replaces | Cost | Failure mode | Removal |
| --- | --- | --- | --- | --- | --- | --- |
| Execution budget | Max-spend FULL | Anthropic dropped always-on eval | Complements intensity; does not replace bars | Derive + record | Ignored budget = waste | Drop fields models no longer need |
| Conditional independent judge | Self-grade vs ~20x GAN | cwc; Anthropic 2026 then simplification | Complements mechanical + Layer A | Extra context when on | Trigger-miss → weak UI ships | Off when self-eval is honest |
| Layer A + Layer B | Mocks; taste-as-Playwright | Canary; screenshot QA | Replaces mock `visual_qa`; stills stay | Browser + optional B | A green / B red confusion | B removable; A is not |
| Lightweight 3–5 directions | Five palettes or five sites | Two-lock lab already | Complements contract; replaces 23-card named-ref path | Cheap sketches | Skipping pick | Skip when DESIGN.md pasted |
| Capsule + token ladder | Brain dumps | Tenet compile_context | Replaces ad-hoc worker prompts | Packager | Fat capsule | Collapse if unused |
| Retry-must-change | Blind loops | Ralph / Tenet / Zenith | Replaces unbounded retry | Taxonomy + N | Wrong class | Shrink unused classes |

---

## Output-impact ranking of adopted mechanisms

1. Layer A live Playwright (kills mock QA). Perception is Layer B.
2. Host/docs/skill parity with 3.3.1 (four missing skills; skip-vs-translate docs).
3. Adaptive evaluator (independent only when **triggers** fire — not because FULL).
4. Worker capsule + token ladder.
5. Failure class + retry-must-change.
6. Anti-slop as absence-of-decision + 2026 cliché calibration.
7. 3–5 **lightweight** directions on flagship, then one implementation.
8. Mission file + derived execution budget.
9. Capability → model/reasoning/host/tool recommendation.
10. Registry honesty (verified vs unverified; empty overlay).
11. Human-needed gate specified (including material tool/model gap).
12. LIGHT intensity (CRUD/debug without Design Lab tax).
13. Collapse duplicate browsers; Layer B when budgeted.
14. Optional OpenDesign worker **later**, not wired.
15. Memory classes + 3I hygiene KEEP/MIGRATE/DELETE/ARCHIVE.

Items 1–7 are the ones that can fail the architecture if they do not show up in A–I benchmarks **without** a 4× cost anti-pattern.

---

## Score matrix (0–5)

Axes: **Out** real output · **Des** design · **Rel** reliability · **Ctx** context · **Spd** speed · **Rec** recovery · **Port** portability · **Mnt** maintainability · **Cx** complexity cost (higher = worse) · **Dup** duplication risk (higher = worse).

| Mechanism | Out | Des | Rel | Ctx | Spd | Rec | Port | Mnt | Cx | Dup | Decision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Anthropic independent eval + live Playwright | 5 | 4 | 5 | 3 | 1 | 4 | 4 | 4 | 4 | 2 | ADAPT (not always-on) |
| Anthropic progress.txt + git | 4 | 1 | 5 | 5 | 3 | 5 | 5 | 5 | 2 | 1 | ADAPT FULL only |
| cwc evaluator no-write | 4 | 2 | 5 | 4 | 2 | 3 | 4 | 5 | 2 | 2 | ADAPT |
| Canary/Codex harness-owned Playwright | 5 | 3 | 5 | 4 | 3 | 4 | 4 | 4 | 2 | 2 | ADAPT |
| Tenet compile_context capsule | 4 | 2 | 4 | 5 | 4 | 3 | 3 | 4 | 2 | 2 | ADAPT |
| Tenet DAG / SQLite / MCP OS | 3 | 1 | 4 | 2 | 2 | 5 | 2 | 2 | 5 | 4 | REJECT |
| Tenet critics + retry/heartbeat | 4 | 2 | 5 | 3 | 3 | 5 | 3 | 3 | 3 | 2 | ADAPT |
| Spec Kit artifacts | 3 | 2 | 4 | 4 | 3 | 3 | 5 | 4 | 2 | 3 | ADAPT compact |
| Spec Kit slash zoo / constitution | 2 | 1 | 3 | 2 | 2 | 2 | 3 | 2 | 4 | 5 | REJECT |
| Deep Suite skippable + independent review | 3 | 2 | 4 | 3 | 3 | 3 | 3 | 3 | 4 | 4 | COMBINE skip only |
| Ralph fresh loop + max iter + stuck | 3 | 1 | 4 | 5 | 2 | 4 | 4 | 4 | 3 | 3 | ADAPT budget |
| Ralph unattended PRD default | 2 | 0 | 2 | 3 | 2 | 2 | 3 | 3 | 4 | 3 | REJECT |
| Zenith stop discipline | 3 | 1 | 4 | 3 | 3 | 4 | 3 | 3 | 2 | 2 | ADAPT |
| Zenith multi-day MCP OS | 2 | 1 | 3 | 2 | 1 | 4 | 2 | 2 | 5 | 4 | REJECT |
| Nonstop immutable acceptance | 3 | 1 | 5 | 4 | 3 | 3 | 4 | 5 | 2 | 3 | ADAPT |
| frontend-design critic | 4 | 5 | 3 | 5 | 4 | 2 | 5 | 5 | 1 | 2 | ADAPT as text |
| no-slop / ui-craft questions | 3 | 4 | 3 | 4 | 4 | 2 | 4 | 4 | 2 | 3 | ADAPT questions |
| Banned-pattern law | 2 | 2 | 2 | 4 | 4 | 1 | 4 | 3 | 3 | 3 | REJECT |
| Get Design Done verify-vs-brief | 3 | 4 | 3 | 3 | 3 | 2 | 3 | 3 | 3 | 4 | ADAPT overlap only |
| twaldin RunSpec idea | 3 | 1 | 4 | 4 | 3 | 3 | 4 | 4 | 2 | 2 | ADAPT rec only |
| twaldin wrap every CLI | 2 | 0 | 3 | 3 | 2 | 3 | 4 | 2 | 4 | 5 | REJECT |
| LIGHT / tiny lane | 4 | 1 | 4 | 5 | 5 | 3 | 5 | 5 | 1 | 1 | ADAPT |
| SDD clean workers | 3 | 2 | 4 | 5 | 4 | 3 | 4 | 4 | 2 | 3 | ADAPT |
| `/sdd-*` zoo | 2 | 1 | 2 | 2 | 2 | 2 | 3 | 2 | 4 | 5 | REJECT |
| HarnessXP fleet sync | 1 | 0 | 3 | 3 | 2 | 3 | 2 | 3 | 4 | 3 | DEFER |
| Memory plugins | 2 | 0 | 3 | 3 | 2 | 2 | 3 | 2 | 4 | 4 | DEFER |
| Spexcode product | 3 | 2 | 4 | 4 | 2 | 4 | 2 | 3 | 5 | 3 | DEFER |
| ShioriCode GUI | 1 | 0 | 3 | 2 | 2 | 4 | 2 | 2 | 4 | 5 | REJECT |
| Click-to-source preview IDE | 2 | 3 | 2 | 3 | 3 | 2 | 2 | 3 | 4 | 4 | DEFER |
| Pixel-clone / SSIM DoD | 1 | 0 | 2 | 3 | 2 | 1 | 3 | 3 | 3 | 2 | REJECT |
| Agent Harness “v10” import | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 5 | 5 | REJECT / INVESTIGATE |
| GSD / ECC / Karpathy dumps | 1 | 1 | 2 | 1 | 1 | 2 | 2 | 1 | 5 | 5 | KEEP REJECTED |
| Agentic Harness install | 2 | 1 | 3 | 3 | 2 | 2 | 3 | 3 | 3 | 4 | KEEP existing layers |

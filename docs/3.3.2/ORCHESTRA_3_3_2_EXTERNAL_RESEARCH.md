# Orchestra 3.3.2 — external research

Phase 3A. Research only. Not live. Not a toolbox expansion.

Fetched 2026-09-15 from current GitHub/docs. Mechanisms, not architectures. Stars and README size are not evidence.

Source snapshot of Orchestra: public clone `C:\projects\orchestra-workflow` HEAD `a9a5051`, `VERSION` 3.3.1. Live hosts on this PC remain 3.3.0 (see reconciliation).

Score legend (0–5) used in the companion matrix: output, design, reliability, context, speed, recovery, portability, maintainability, **complexity cost**, **duplication risk**. Higher complexity/duplication is worse.

---

## Method

1. Audit the current Orchestra OS (implemented vs documented vs stale).
2. Fetch each named system’s current repo or official article.
3. Extract one transferable mechanism, or say the source is dead/ambiguous.
4. Reject whole-system imports.
5. Prefer systems that showed real running apps, Playwright verdicts, or measured cost/quality — not architecture diagrams.

---

## Orchestra freeze facts this research sits on

Already true on disk. Do not rediscover.

- Repo contract is 3.3.1. Live Cursor/Gemini conductor skills and `~\.gemini\GEMINI.md` are 3.3.0. Four job-loaded skills exist in-repo and are missing live.
- Engine two-lock Design Lab is implemented. `.cursorrules` / README still teach the 3.3.0 named-website skip.
- `runtime/internal/engine/visual_qa.go` can write `MOCK_SCREENSHOT_*` placeholders. That is not visual verification.
- Catalog 168 rows; 42 `last_verified` (all 2026-09-06); 126 unverified. Brain overlay `total_evaluations: 2648` with empty `evaluations` array.
- OpenDesign is banned as a second studio vs Stitch. Isolated lab (Claude primary, Cursor Agent `cursor-grok-4.6-high` secondary) stays outside Orchestra.
- 3.3.1 already rejected GSD/ECC dumps, Karpathy skill packs, always-on Magic UI (`docs/3.3.1/CAPABILITY_AUDIT.md`).

---

## Named systems

### 1. Tenet — ADAPT capsule + critics + retry budget; REJECT DAG OS

- Source: https://github.com/JeiKeiLim/tenet (also https://github.com/zoytsa/tenet)
- Evidence: shipping CLI/MCP; SQLite job state; heartbeat/orphan recovery; independent code/test/Playwright critics; `tenet_compile_context` for workers. Real 12h-cycle claim; not independently re-run here.
- Useful: compiled worker context; independent critics with the same spec as the worker; failure/retry with a budget; stuck/heartbeat.
- Reject: importing a DAG job manager, SQLite orchestrator, or 19 MCP tools.

### 2. Anthropic long-running harness — ADAPT independence + live-app eval; REJECT always-on three-agent theatre

- Sources:
  - https://www.anthropic.com/engineering/harness-design-long-running-apps (2026)
  - https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents (2025)
  - https://github.com/anthropics/cwc-long-running-agents
- Evidence: planner → generator → evaluator; evaluator uses Playwright on the running app; self-eval praised mediocre work. Reported ~20x cost vs single-agent. Later work **removed** forced per-sprint evaluation when the model was strong enough. That simplification is as important as the original GAN loop.
- Useful: builder must not grade itself; browser eval on the live app; structured artifacts between sessions; sprint contract only when the spec is high-level.
- Reject: forcing planner+evaluator on CRUD; infinite eval loops; treating $100+ runs as the default.

### 3. Agentic Harness (sevenschulte) — KEEP as a map of layers Orchestra already has

- Source: https://github.com/sevenschulte/agentic-harness
- Evidence: six layers — identity (`AGENTS.md`), rules, skills, hooks, curated memory, scheduled agents. Small reference repo, not a production OS.
- Useful: identity vs rules vs skills vs memory vs hooks stay separate; memory updates explicitly (files do not self-write).
- Reject: installing this repo; scheduled unattended agents.

### 4. HarnessXP — DEFER cross-repo contracts

- Source: https://github.com/mebusw/HarnessXP
- Evidence: one orchestrator repo + `shared/` sync to fight API drift across N product repos. Solves a problem Orchestra does not have yet (no multi-repo product fleet in production).
- Useful later: ownership + contract path + sync check.
- Reject now: building a distributed orchestrator.

### 5. Deep Suite — ADAPT skippable phases + independent review; REJECT as core install

- Source: https://github.com/Sungmin-Cho/claude-deep-suite (+ deep-work / deep-review / deep-loop)
- Evidence: Brainstorm → Research → Spec → Plan → Implement → Test with skip flags; separate reviewer so the author does not sign its own diff; receipts as evidence.
- Useful: skippable brainstorm; plan approval gate; independent review plugin pattern (Orchestra already has `fresh-context-review`).
- Reject: installing Deep Suite; mutation-test theatre on LIGHT work. GSD remains REJECTED as core (archived `gsd-build/get-shit-done` → `open-gsd/gsd-core`; pattern already in Orchestra’s loop).

### 6. GitHub Spec Kit — ADAPT compact spec/plan/tasks; REJECT slash-command zoo

- Source: https://github.com/github/spec-kit
- Evidence: specify → plan → tasks → implement → converge as markdown artifacts; optional clarify/analyze gates; generic agent integration.
- Useful: intent before stack; tasks grouped by user story; converge is append-only (does not silently rewrite code).
- Reject: `/speckit.*` as a second conductor; constitution file competing with `AGENTS.md`.

### 7. Agent Harness v10 — INVESTIGATE / DEFER (name is ambiguous)

No single public repo titled “Agent Harness v10 Native/Portable/Audited” matched cleanly. Closest:

- OLAV Agent Harness claims v0.10.0 audit logging: https://docs.olavai.com/concepts/agent-harness/
- HarnessAudit (UCSB) audits whether native harnesses obey tool/resource boundaries: https://github.com/UCSB-AI/HarnessAudit
- `ar27111994/agent-harness` is a cross-host **asset supply chain**, not an execution OS.

Decision: do not invent a v10 product. Steal only **choose rigor from risk** (tiny vs high-risk lanes) if a later source is pinned. `minhtran3124/agent-harness` has a useful `tiny|normal|high-risk` lane with hook enforcement — ADAPT the lane idea, REJECT the rest.

### 8. Zenith — ADAPT stopping discipline; REJECT multi-day MCP orchestrator

- Source: https://github.com/Intelligent-Internet/zenith
- Evidence: premature completion is the dominant failure; compared five harnesses; Ralph is a strong simple baseline but expensive and stop-rule-weak; Zenith claims better mean rank at <½ Ralph cost (study, not re-run here).
- Useful: gap-finding vs original requirement; independent testers; explicit stop decision.
- Reject: days-long MCP/ACP mission OS.

### 9. Ralph — ADAPT fresh-context loop + progress file + max iterations; REJECT unattended PRD loops as default

- Sources: https://github.com/snarktank/ralph , https://github.com/harrymunro/ralph-wiggum , https://github.com/Kiriketsuki/ralph-loop
- Evidence: each iteration is a new instance; memory is git + `prd.json` + `progress.txt`; max iterations; some variants add stuck detection and a QA agent that cannot read source.
- Useful: one story per context window; progress file; iteration budget; optional source-blind QA.
- Reject: fire-and-forget loops on PREMIUM visual work (needs human stills). Orchestra already forbids silent skip of Design Lab.

### 10. Harness (twaldin) — ADAPT RunSpec/RunResult idea; REJECT wrapping every CLI

- Source: https://github.com/twaldin/harness
- Evidence: one `RunSpec → RunResult` for many headless CLIs; model name normalization; timeout; tokens/cost when the CLI reports them.
- Useful: capability → runtime invocation contract; timeouts; structured results.
- Reject: adding 13 CLI adapters Orchestra will not run (Aider, SWE-agent, Crush, Kilo — several already banned in host policy).

### 11. Harness Skills — COMBINE with existing progressive load

- Evidence: Agent Skills `SKILL.md` standard + progressive disclosure (sevenschulte; OpenMind7 guide). Orchestra already job-loads four thin skills and forbids loading every kit.
- Useful: load matching skill text, not the whole registry; inject compact rules into workers (Gentleman-Programming SDD orchestrator does this well).
- Reject: bulk skill install; a second skill standard.

### 12. Harness SDD — ADAPT clean-context workers; REJECT phase zoo

- Closest live sources: GitHub Spec Kit; https://github.com/gotalab/cc-sdd ; Gentleman-Programming `sdd-orchestrator.md` (explore → spec → design → tasks → apply → verify; workers get fresh context and compact rules, not the whole registry).
- “Inception / planner / driller / orchestrator” did not resolve to one canonical repo. Treat as the SDD pattern, not a product to vendor.
- Useful: workers do not inherit the orchestrator’s context dump; size alone must not force SDD.
- Reject: `/sdd-*` command surface as Orchestra’s UI.

### 13. Memory Harness — DEFER measured memory; KEEP explicit Brain writes

- Sources: https://github.com/kenhuangus/agent-memory-harness (benchmarks Haiku+memory vs Opus); https://github.com/Chachamaru127/harness-mem ; https://github.com/rohitg00/agentmemory
- Evidence: memory quality can be measured; auto-extracting memory servers add MCP surface and stale-fact risk.
- Useful: provenance/confidence/timestamp on memories; explicit save (sevenschulte).
- Reject: 54-tool memory MCP; silent observation compression into Brain.

### 14. ShioriCode — REJECT as an Orchestra product; KEEP reconnect as a host concern

- Source: https://github.com/shiorihq/ShioriCode
- Evidence: desktop/web broker for long Codex sessions, timeline, reconnect. A UI around a provider CLI — not a control plane.
- Orchestra should not grow a session GUI. Cursor/Claude/Codex already own reconnect.

### 15. Nonstop Agent — ADAPT initializer vs coding session; KEEP immutable acceptance

- Source: https://github.com/seolcoding/nonstop-agent
- Evidence: implements Anthropic’s two-agent pattern; `feature_list.json` descriptions/steps are immutable; git + progress files; command allowlists.
- Useful: first session sets the board; later sessions implement one item; do not rewrite acceptance mid-run.
- Reject: generating a 24/7 Claude daemon.

### 16. Spexcode — DEFER spec/code pairing; ADAPT eval-on-product-surface later

- Source: https://github.com/shuxueshuxue/spexcode — https://spexcode.net/
- Evidence: spec nodes in git; workers in worktrees; human merges; evals against the real product. Strong idea, extra product to operate.
- Useful later: spec and code commit together; eval brackets on the running app.
- Reject now: Spex as a required runtime.

### 17. Get Design Done — ADAPT brief → explore → plan → design → verify; REJECT GDD install

- Sources: https://github.com/hegemonart/get-design-done ; https://github.com/0martinez/get-design-done
- Evidence: design workflow with tokens, local design knowledge, verify against brief. Overlaps Orchestra Design Lab + `DESIGN.md`.
- Useful: verify against the brief and tokens, not taste adjectives; keep a local design-knowledge file per product.
- Reject: another `/gdd:*` command pack. Orchestra already has two locks.

### 18. frontend-craft / ui-craft — ADAPT detection heuristics; REJECT CI slop scorer as core

- Source: https://github.com/educlopez/ui-craft (closest “frontend-craft”)
- Evidence: persistent brief/tokens; deterministic anti-slop scanner (43 rules). Useful for LIGHT gates; dangerous if it becomes a blacklist that forbids legitimate cream/serif when the brief asks.
- Useful: detect `transition: all`, missing empty/error states, placeholder copy, purple-gradient tells — as **signals**, not laws.
- Reject: 0–100 taste scores as Definition of Done.

### 19. no-slop-ui — ADAPT checklist questions; REJECT banned-pattern dump as always-on

- Source: https://github.com/LeoStehlik/no-slop-ui
- Evidence: concrete frontend rules; also ships `banned-patterns.md`. Orchestra’s own `visual_qa.go` already bans Inter/Space Grotesk and `#000` via string match — too crude, and it contradicts Anthropic’s 2026 skill (cream-serif is a cliché **unless the brief asks**).
- Useful: structural questions (hierarchy, density, one primary CTA).
- Reject: “never gradients / never cards / never purple” as policy.

### 20. Anthropic frontend-design skill — ADAPT as critic text; REJECT always-on skill dump

- Source: https://github.com/anthropics/skills/blob/main/skills/frontend-design/SKILL.md
- Evidence: current skill names three 2026 clichés (cream+#F4F1EA+terracotta serif; near-black + acid accent; broadsheet hairlines). Brief wins if it asks for one of them. Spend boldness in **one** place; remove one accessory.
- Directly relevant to the fictional Northline pages in OpenDesign Task 3 (cream paper + serif is exactly cliché 1 unless chosen on purpose).
- Useful: one signature; product-specific choices; clichés as calibration not bans.
- Reject: installing anthropics/skills wholesale (Karpathy/skill-dump policy).

### 21. Interactive preview / design-engineer UIs — DEFER host features; ADAPT computed-style evidence

- Sources: Tesslate OpenSail Design Engineer; Kombai browser; Visionaire `assert_visual`; Pixelprobe computed-style diffs; Imugi pixel-loop (Figma match).
- Evidence: click-to-source and live CSS inspectors improve human visual decisions. Pixel-match loops risk cloning a reference (forbidden for Elaine Yu / any distinctive site).
- Useful: Playwright `getComputedStyle` + overflow/contrast as **evidence** attached to QA; human stills remain the pre-implementation lock.
- Reject: Figma/pixel-perfect clone loops; a second in-Orchestra visual IDE.

### 22. Design Engineer (completeness + computed style) — ADAPT completeness gates; REJECT pixel-clone as DoD

- Combine Get Design Done verify, Visionaire assertions, Pixelprobe breakpoints.
- Useful: completeness (empty/error/loading/focus), computed spacing/type, responsive sweep.
- Reject: SSIM-to-mockup as the quality bar for original work.

---

## Additional systems with practical output

| System | URL | Why it matters | Decision |
| --- | --- | --- | --- |
| Anthropic cwc-long-running-agents | https://github.com/anthropics/cwc-long-running-agents | Evaluator subagent with **no Write/Edit**; `/goal` loop | ADAPT |
| Codex long-running harness | https://github.com/LongWeihan/codex-long-running-harness | Executable browser/API/DB checks; small trial (3/3) | ADAPT eval ownership |
| Canary Lab | https://github.com/ferterahadi/canary-lab | Harness runs Playwright and owns the verdict; agent cannot self-green | ADAPT |
| Visionaire | https://github.com/mi60dev/visionaire-engine | Deterministic geometry assertions, not LLM taste | ADAPT later for PREMIUM |
| minhtran agent-harness lanes | https://github.com/minhtran3124/agent-harness | `tiny\|normal\|high-risk` with hook-enforced blast radius | ADAPT intensity |
| cc-sdd | https://github.com/gotalab/cc-sdd | Direct implement path when discovery says small | KEEP (matches LIGHT) |

Reject as import even if useful to read: OpenHands, OmniRoute, Kilo, skill marketplaces, MCP server catalogs, Magic UI always-on, OpenDesign as conductor.

---

## What this research does **not** justify

- Replacing Orchestra with Tenet, GSD, Deep Suite, Spec Kit, Zenith, or Ralph.
- A model zoo or 13-CLI gateway.
- Pixel-matching a competitor’s site.
- Loading all skills into every run.
- Treating mock screenshots as QA.
- Wiring OpenDesign into the conductor.

---

## 3A.1 addendum (mechanisms only — scores unchanged)

Does not re-rank the 22 systems. Confirms how 3.3.2 should use them:

- **Layer A** = harness-owned Playwright verdict on the running app (Canary Lab / Codex executable eval). Objective checks only. The builder must not self-green Layer A.
- **Layer B** = screenshot-capable **perception** of captures produced by Layer A (plus pre-implementation stills as the intended direction). This is not a second browser. Playwright does not grade art direction.
- **Independent evaluation stays adaptive.** Anthropic’s 2026 harness showed builder self-eval is weak; later work **removed** forced per-sprint eval when the model was strong enough. Orchestra: FULL is not “always judge.” Triggers on / off live in ARCHITECTURE. Goal is QUALITY × EVIDENCE / COST, not more agents.
- **Directions before commitment** matches Design Lab + Spec Kit intent: cheap artifacts first, expensive implementation once. 3–5 lightweight directions, one contract, one build.
- Still reject: DAG/SQLite job OS, skill dumps, always-on Magic UI, OpenDesign as conductor, star-count READMEs as evidence.

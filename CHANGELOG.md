# Changelog

## [3.3.0] - Richer graph, four hosts, two-stage Design Lab, free-resource honesty
*Evolve 3.2. Same OS. Curated 40 skills. Contract stays open. Version stays 3.3.0.*

- Four-host sync: Cursor, Antigravity, Claude Code, jcode (`~/.jcode/skills`). Allowlist from `registries/host-stack.json`.
- Design Lab: 23 short cards → one full `DESIGN.md`. Named override outranks the graph. Custom paste wins.
- Catalog: GetLayers `PURCHASE_PENDING`. Refero MCP AUTH_REQUIRED / PRO; Refero public pages FREE. 21st OPTIONAL_FREE_MCP (2 installs/day). Magic UI FREE MCP + OSS. Tailkit MCP PAID; public Tailkit RESOURCE.
- 13th capability `interaction-components` (one route across kits, not a skill per library). Tint/adapt on `reverse-engineering`.
- Free research: Godly, Land-book, Lapa, SaaSFrame. OSS 3D: three.js, R3F, drei, postprocessing, scroll-rig. Prompt catalogs as references. devops-skills quarantined (no 88 dump).
- Docs: RESOURCE_INVENTORY A–Q, RESOURCE_HEALTH, DESIGN_RESOURCE_ROUTING, manual-setup 9-field, diagrams `01-control-plane`…`08-resource-combination`.
- Global skills added (wrappers, not dumps): refero-design, hallmark, unlazy, react-best-practices, gsap-core, supabase, postgres-best-practices, diagram-generator, pretty-mermaid, skill-security-check.
- Job packs in `templates/packs/`. README comparison is architectural judgment with estimated token ranges — no fake A/B.
- Rollback: `v3.2.0` (`9d6900d`) and the 3.1 zip. Uncommitted 3.2 hygiene (YAML stubs, Go module rename) lands here.

## [3.2.0] - Intelligence planes on the same OS
*Evolve 3.1. Do not rebuild. Keep the 30 skills. Contract stays open.*

- Research: Scrapling is the primary public crawler; Firecrawl on-demand; Crawl4AI catalog fallback. Playwright stays on **our** app.
- Design: Designer SOS (S.A.V.E., six-pass, Keep/Kill/Park, Vault/Live/Anatomy); named source tiers; GetLayers on premium/3D/shader only.
- Quality: Unlazy-style verification ledger; Taste → DESIGN.md → Impeccable → Playwright → SOS → one consultant packet.
- Catalog (not global install): `playwright-cli`, `crawl4ai`, `serena`, `jcode`. Optional telemetry fields on rows. Cheap trigger matching in `AGENTS.md`.
- Ingest: graph miss → research → one named capability → wait for **update**. `templates/custom-skill.md`.
- Executors: jcode as a fifth host adapter under `AGENTS.md`. Does not replace Cursor / Antigravity / Claude Code.
- Taste: on-disk Stitch wrapper vs upstream v2 experimental — do not `npx skills add` the pack.
- README rewritten for the evolve story. Overlay + Antigravity adapters identify as 3.2.0.
- Post-tag hygiene (working tree after `v3.2.0` / `9d6900d`): drop duplicate YAML catalogs; rename Go module to `github.com/mahik504/orchestra-workflow/runtime`; genericize personal from-strings.
- Catalog: Soup is CURATED_OPTIONAL for ML fine-tune jobs (not a 31st skill). Anti-Slop AJ + blader/humanizer + Graphify catalogued; none dumped into the 30.
- jcode on this operator PC: OmniRoute default; `auto/best-coding` / `auto/best-reasoning` / `auto/best-fast` / `auto/best-free` aliases. Bare `jcode` is enough after PATH.

## [3.1.0] - Control plane
*One contract. Design Lab is a gate. Evidence-first completion.*

- Overlay, Cursor, Antigravity, and Claude Code share the same 3.1.0 rule.
- Public registry stripped of personal catalog rows. MCP health is explicit.
- `orchestra doctor` warns if Antigravity science or data-agent-kit plugins are Global.

**Routing.** The classifier was a keyword stub returning `task-stub-001`; it now scores every
capability row against the brief and produces a structured re-brief (type, archetype, quality bar,
platform, research depth, verify depth, hard constraints). Every row in
`design-resource-graph.json` gained `trigger_conditions`, `skip_conditions`, `quality_bar`, and
`risk_rank` — a route that can never decline itself is a default in disguise. Declined routes are
reported with the skip condition that fired. Two close routes produce exactly one question; silence
takes the lower `risk_rank` and logs `assumed <capability>, no response`.

**Design Lab is a lock, not a warning.** `verify.DesignLab` blocks writes to anything a browser
renders while the gate is `PENDING`, including on dry runs, which previously slipped past the
synthesize halt. Directions must number 2 or 3 and name a source for typography, colour, and motion.
Rejections persist to `.orchestra/design-lab/rejected-directions.json`, fingerprinted by stack, so a
renamed rejected direction cannot be re-offered. Bypass is allowed but requires a note.
See `protocols/DESIGN_LAB_PROTOCOL.md`.

**No second conductor.** Deleted the unreferenced `internal/kernel` and `internal/planner` packages,
which duplicated the pipeline's approval and execution logic.

**Isolation.** A fresh clone keeps resource memory at `.orchestra/memory/resource-memory.json`.
The public method repo does not ship a live `memory/resource-memory.json`. Set `ORCHESTRA_HOME`
to a private Brain so overlay and outcomes resolve there instead of the clone.

**Phase 7 evidence.** Twelve capability rows. `research-paper` loads `orchestra-docs` and does not arm the Design Lab. Restaurant / school SaaS / 3D portfolio resolve to three different resource graphs. A premium brief that names no libraries still surfaces research, design, and motion. `pixi.js` is reported as unknown technology rather than forced into a route. References acquire `on_demand` via `web_fetch`; GSAP is `project_scoped_install`; global npm stays blocked. Live fixture screenshots at 1440×900, 768×1024, 390×844; click wrote a hold; heading contrast 14.43:1; no horizontal overflow.

**User-added resources.** `orchestra add --intent "..."` inspects a URL, infers kind
(skill / dependency / MCP / plugin / subagent / reference / adapter), and writes a Brain overlay
at `memory/added-resources.json`. It does not edit `registries/resources.json` or the capability
graph. Matching tasks activate the overlay row; non-matching tasks skip it. Recorded outcomes in
`memory/resource-memory.json` update overlay routing: last failure suppresses auto-activation.
That loop is experience → evaluation → memory → routing update → future selection. It is not
reinforcement learning. `orchestra lifecycle` prints the 15-step proof for one URL.

**README.** Rewritten against the code. The previous version cited capability and domain names that
do not exist in the graph, and described research as querying live design indexes when
`ResearchCoordinator` defaults to offline fixtures. CI and Hygiene badges sit on the README.
The eight-stage architecture diagram from `ARCHITECTURE.md` is inlined there.

**Doctor.** Playwright is checked as `npx playwright --version` (split argv). Missing Playwright
is `[OPTIONAL]` — visual jobs use a project install or the agent browser. Global npm stays blocked.

**Engine install.** `kit/install-local-engine.ps1` / `.sh` persist `ORCHESTRA_HOME` as a user env
and `go install` `orchestra` onto `~/go/bin`. Chat still uses the markdown contract without the
binary. `docs/getting-started.md` treats `classify` / `plan` as optional.

**Catalog.** `build-your-own-x` rationale is a public bookmark sentence. It does not name private notes.
`drei` (`@react-three/drei`) is a public catalog row, project-scoped npm, so a clone does not need the author's overlay.
Public catalog slimed to the named design stack plus a small refuse map. `liquid-glass` points at dashersw/liquid-glass-js. `vgpu-shaders` points at Vercel `vgpu`. `scroll-world` is optional on 3D routes. School SaaS no longer loads Awwwards or museums.

**Host stack.** `registries/host-stack.json` names the 30 skills, MCP templates, and marketplace plugins per host.
`kit/bootstrap.ps1` / `.sh` is the clone front door: pick hosts, copy skills and adapters, never overwrite a live MCP config.

## [3.0.0] - The Castle Pass Release
*Orchestra V3: Capability router and Go execution layer.*

- **Core Engine Rewrite**: Replaced V1/V2 monolithic Python scripts with a compiled, isolated Go binary (`orchestra.exe`).
- **4-Stage Capability Pipeline**: Strict `Retrieval -> Analysis -> Execution -> Verification` routing model that evaluates tasks before agent allocation.
- **Lazy Loading Context**: Reduced typical execution prompts from 40,000+ tokens to ~1,500 tokens using JIT capability injection.
- **Cryptography-Backed Handoffs**: Added `internal/handoff` system utilizing SHA256 file checksums and `state.json` versioning to prevent Cursor/Antigravity collision and silent data corruption.
- **Clean Workspace Generation**: Added `orchestra init` to template private local Brains safely isolated from the public repository.
- **Adversarial Playwright/Lighthouse Integration**: Enforced visual QA workflows directly out-of-the-box via programmatic verification steps.
- **Clean-Clone Tested**: Absolute separation between this public workflow repository and any private workspace it is pointed at.

## [2.0.0] - Orchestra Workflow v2
*Capability router release.*

- Protocols: design, typography, reverse engineering, visual QA, security, research, sensory, prompt brief
- Registries with CORE / SPECIALIST / OPTIONAL / EXPERIMENTAL / REJECTED
- Activation card before high-visual UI
- Init + skill-install scripts
- Adapters: Cursor, Antigravity, Claude, Gemini, Codex, OpenCode, Hermes

## [1.0.0] - Orchestra Workflow v1
*Jump table, dual-track import, Cursor + Antigravity packets.*

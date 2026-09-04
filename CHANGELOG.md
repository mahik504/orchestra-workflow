# Changelog

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

**Isolation.** A fresh clone now keeps resource memory at `.orchestra/memory/resource-memory.json`
instead of creating a `memory/` directory in the host repository.

**Phase 7 evidence.** Twelve capability rows. `research-paper` loads `orchestra-docs` and does not arm the Design Lab. Restaurant / school SaaS / 3D portfolio resolve to three different resource graphs. A premium brief that names no libraries still surfaces research, design, and motion. `pixi.js` is reported as unknown technology rather than forced into a route. References acquire `on_demand` via `web_fetch`; GSAP is `project_scoped_install`; global npm stays blocked. Live fixture screenshots at 1440×900, 768×1024, 390×844; click wrote a hold; heading contrast 14.43:1; no horizontal overflow.

**README.** Rewritten against the code. The previous version cited capability and domain names that
do not exist in the graph, and described research as querying live design indexes when
`ResearchCoordinator` defaults to offline fixtures.

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

# Changelog

## [3.0.0] - The Castle Pass Release
*Orchestra V3: Capability router and Go execution layer.*

- **Core Engine Rewrite**: Replaced V1/V2 monolithic Python scripts with a compiled, isolated Go binary (`orchestra.exe`).
- **4-Stage Capability Pipeline**: Strict `Retrieval -> Analysis -> Execution -> Verification` routing model that evaluates tasks before agent allocation.
- **Lazy Loading Context**: Reduced typical execution prompts from 40,000+ tokens to ~1,500 tokens using JIT capability injection.
- **Cryptography-Backed Handoffs**: Added `internal/handoff` system utilizing SHA256 file checksums and `state.json` versioning to prevent Cursor/Antigravity collision and silent data corruption.
- **Clean Workspace Generation**: Added `orchestra init` to template private local Brains safely isolated from the public repository.
- **Adversarial Playwright/Lighthouse Integration**: Enforced visual QA workflows directly out-of-the-box via programmatic verification steps.
- **Clean-Clone Tested**: Absolute separation of public `orchestra-workflow` and private `orchestra-brain` data boundaries.

## [2.0.0] - Orchestra Workflow v2
*Capability router release.*

- Protocols: design, typography, reverse engineering, visual QA, security, research, sensory, prompt brief
- Registries with CORE / SPECIALIST / OPTIONAL / EXPERIMENTAL / REJECTED
- Activation card before high-visual UI
- Init + skill-install scripts
- Adapters: Cursor, Antigravity, Claude, Gemini, Codex, OpenCode, Hermes

## [1.0.0] - Orchestra Workflow v1
*Jump table, dual-track import, Cursor + Antigravity packets.*

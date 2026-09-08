# Host capability map — Orchestra 3.3.1

Sync means **the same contract**, not the same plugin list.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

The curated skill core is identical on Cursor, Antigravity, Claude Code, and jcode. Host extras stay on the host that owns them.

| Capability | Cursor | Antigravity | Claude Code | jcode |
| --- | --- | --- | --- | --- |
| Contract | overlay `AGENTS.md` + `.cursorrules` | `kit/antigravity/MASTER-PROMPT.md` | `CLAUDE.md` + `~/.claude/CLAUDE.md` | `AGENTS.md` + `~/.jcode/AGENTS.md` |
| Skill dir | `~/.cursor/skills` | `~/.gemini/config/skills` | `~/.claude/skills` | `~/.jcode/skills` |
| Browser / visual QA | Cursor browser tools | Playwright MCP | Playwright MCP | packet / Playwright if configured |
| 2D screens | Stitch MCP | StitchMCP | Stitch MCP if configured | — |
| Library docs | Context7 | Context7 | Context7 | — |
| Vault IO | filesystem MCP | filesystem MCP | filesystem MCP | — |
| Mail | Cursor Gmail plugin | — do not install | — | — |
| OmniRoute | — | — | — | OpenAI-compatible default_provider |

Rules:

- Do not install Gmail on Antigravity. Do not install Firebase on Cursor "to match."
- A server in `AUTH_REQUIRED` or `PURCHASE_PENDING` is not active.
- Host extras are OPTIONAL capabilities. The allowlist is the contract.
- After `kit/bootstrap.ps1` (or `kit/sync-ides.ps1`), hosts identify as Orchestra 3.3.1.

Auth clicks: `docs/manual-setup/`. Rollback: `ORCHESTRA_CONTRACT` (see `kit/ROLLBACK.md`).

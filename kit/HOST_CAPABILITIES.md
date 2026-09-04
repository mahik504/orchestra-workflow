# Host capability map — Orchestra 3.1.0

Sync means **the same contract**, not the same plugin list.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

The 30-skill core is identical on Cursor, Antigravity, and Claude Code. Host extras stay on the host that owns them.

| Capability | Cursor | Antigravity | Claude Code |
| --- | --- | --- | --- |
| Contract | overlay `AGENTS.md` + `.cursorrules` | `kit/antigravity/MASTER-PROMPT.md` | `CLAUDE.md` |
| Skill dir | `~/.cursor/skills` | `~/.gemini/config/skills` | `~/.claude/skills` |
| Browser / visual QA | Cursor browser tools | Playwright MCP | Playwright MCP |
| 2D screens | Stitch MCP | StitchMCP | Stitch MCP if configured |
| Library docs | Context7 | Context7 | Context7 |
| Vault IO | filesystem MCP | filesystem MCP (local name) | filesystem MCP |
| Mail | Cursor Gmail plugin | — do not install | — |
| Calendar / Drive | Cursor Google plugins | — do not install | — |
| Payments | Cursor Stripe plugin | — do not install | — |
| Cloud data | — | Firebase / Supabase when AUTH is done for that job | — |
| Science / BigQuery packs | — | **off as Global**; on only for that job | — |
| Design references | Mobbin MCP if present | Mobbin if present | — |

Rules:

- Do not install Gmail on Antigravity. Do not install Firebase on Cursor "to match."
- A server in `AUTH_REQUIRED` is not active.
- Host extras are OPTIONAL capabilities. The 30-skill core is the contract.
- After `kit/sync-ides.ps1`, a new Cursor session and the AG kit both identify as Orchestra 3.1.0.

Rollback: `ORCHESTRA_CONTRACT` (see `kit/ROLLBACK.md`). A bad sync is not a git revert of the vault.

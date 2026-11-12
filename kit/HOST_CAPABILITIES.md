# Host capability map — Orchestra 5.5.0

Sync means **the same contract**, not the same plugin list.

**ORCHESTRA = CONTROL PLANE. SKILLS / MCPs / PLUGINS / LIBRARIES = CAPABILITIES. AGENTS = EXECUTORS. BRAIN = MEMORY. REGISTRY = RESOURCE KNOWLEDGE.**

The curated skill core is identical on Cursor, Antigravity, Claude Code, and jcode. Host extras stay on the host that owns them.

| Capability | Cursor | Antigravity | Claude Code | jcode |
| --- | --- | --- | --- | --- |
| Contract | overlay `AGENTS.md` + `.cursorrules` | `kit/antigravity/MASTER-PROMPT.md` | `CLAUDE.md` + `~/.claude/CLAUDE.md` | `AGENTS.md` + `~/.jcode/AGENTS.md` |
| Skill dir | `~/.cursor/skills` | desktop `~/.gemini/config/skills` | `~/.claude/skills` | `~/.jcode/skills` |
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
- After `kit/bootstrap.ps1` (or `kit/sync-ides.ps1`), hosts identify as the version in `VERSION`.

## Antigravity desktop and CLI

Desktop skills sync to `~/.gemini/config/skills`. On this machine `agy.exe --help` lists plugins and agents and does not list a skills command. `~/.gemini/skills` and `~/.gemini/antigravity-cli/skills` are absent. The CLI path is opt-in: `node runtime/tools/ow-sync.mjs --host antigravity-cli`. It is not in the default host map, and hooks are not installed for it. Do not treat desktop and CLI as the same host.

## Native first

Use the host feature when it already does the job. A copied pack that reimplements it stays off.

| Job | Cursor | Claude Code | Antigravity | jcode |
| --- | --- | --- | --- | --- |
| Code review | Bugbot, or one review subagent | `/code-review` | the visual QA pass | the packet's reviewer |
| Worktrees | a git worktree, then `using-git-worktrees` | `claude -w` | — | — |
| Subagents | Task subagents | subagents | — | one worker packet |
| Hooks | `.cursor/hooks.json` | `settings.json` hooks | — | — |
| Goal | the acceptance line in the brief | `/goal` | the brief | the packet acceptance line |
| Memory | the vault | the vault | the vault | the vault |

Claude Code auto-memory is on for some project folders on this PC. `~/.claude/settings.json` has no memory key, so this is the host default, not an Orchestra setting. Those notes are a second memory beside the vault. Orchestra does not read or write them. The vault stays the only long-term memory.

Auth clicks: `docs/manual-setup/`. Rollback: `ORCHESTRA_CONTRACT` (see `kit/ROLLBACK.md`).

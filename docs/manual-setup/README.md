# Manual setup — Orchestra 3.3.0

Only actions a human must take. Never paste keys into git, the vault, or a prompt you will forward.

Clones keep optional adapters. Personal provider tokens and live MCP JSON stay on **your** machine.

| Host | File | Skill dir |
| --- | --- | --- |
| [Cursor](CURSOR.md) | `.cursorrules` | `~/.cursor/skills` |
| [Antigravity](ANTIGRAVITY.md) | `kit/antigravity/MASTER-PROMPT.md` | `~/.gemini/config/skills` |
| [Claude Code](CLAUDE-CODE.md) | `CLAUDE.md` + `~/.claude/CLAUDE.md` | `~/.claude/skills` |
| [jcode](JCODE.md) | `templates/jcode-omniroute-packet.md` | `~/.jcode/skills` |

After every workflow skill change: `kit/sync-ides.ps1`.

Say **skip orchestra** (session) or **skip the lab** (one visual gate). Plain text. Not a slash command.

Status labels used below:

- **FREE NOW** — usable without purchase
- **AUTH REQUIRED** — configured possible; you must sign in
- **PAID BUT RESOURCE-AVAILABLE** — MCP/paid gated; public pages still registered
- **PURCHASE PENDING** — GetLayers. Do not Connect.

Example MCP JSON (placeholders only): `mcp_config.example.json`.

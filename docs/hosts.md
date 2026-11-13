# Hosts

Orchestra 5.5.0 syncs one skill set. The adapters are rendered from `registries/routing.json`.

| Host | Skills | Instructions | Commands |
| --- | --- | --- | --- |
| Cursor | `~/.cursor/skills` | `.cursorrules` and `AGENTS.md` | `.cursor/commands/ow-*.md` |
| Claude Code | `~/.claude/skills` | `CLAUDE.md` | Paste `commands/ow-*.md`. Claude Code has its own slash-command folder when the human adds one. |
| Antigravity | `~/.gemini/config/skills` | `kit/antigravity/GEMINI.md` | The master prompt. No slash-command folder in this repo. |
| Codex | `~/.codex/skills` | `kit/codex/AGENTS.md` | `AGENTS.md` is the instruction file. |
| OpenCode | same agents folder, or the host's skill dir if it has one | `kit/opencode/AGENTS.md` | Markdown prompt. |
| Kimi | the host's skill dir when it exists | `kit/kimi/AGENTS.md` | Markdown prompt. |
| jcode | `~/.jcode/skills` | `AGENTS.md` | Markdown prompt. |

`commands/` is the neutral source. `.cursor/commands/` is the Cursor copy. `agents/` is a reviewer prompt any of these hosts can paste. A host without an agent file format still reads the markdown.

## Sync

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File kit/sync-ides.ps1 -Sync
```

The first command prints a plan. `-Sync` writes. A skill you edited locally is left in place. `-Prune` exits 2. User-installed skills stay.

`creative-3d-web`, `fullstack-mastery`, and the old ids in `registries/skill-aliases.json` are not current skills. `-Sync` moves those folders to the backup directory. A skill you installed yourself, under a name that is not an old Orchestra id, stays.

## Global instruction files

`ow-sync` writes `~/.claude/CLAUDE.md`, `~/.gemini/GEMINI.md`, and `~/.codex/AGENTS.md` when the file is already an Orchestra adapter or is missing. It copies the previous file into the backup directory first. A file that does not look like Orchestra is left alone. `GEMINI.md` must stay under 24KB. `kit/bootstrap.ps1` still stamps a product repo between `<!-- orchestra:start -->` and `<!-- orchestra:end -->`.

Limitless stays optional. The model class list is `registries/model-classes.json`. It does not name a vendor model.

## MCP on this release

`ow-sync` adds `shadcn` (`npx shadcn@latest mcp`) and `magic-ui` (`npx -y @magicuidesign/mcp@latest`) to `~/.claude/settings.json` when that file already exists. It backs the file up first and prints the server names. The template is `templates/mcp/optional-servers.json`. Keyed servers stay unset until you paste a key. Java is required before REA can run.

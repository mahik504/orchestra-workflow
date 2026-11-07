# Hosts

Orchestra 4.5.0 syncs one skill set. The adapters are rendered from `registries/routing.json`.

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

`creative-3d-web` and `fullstack-mastery` are not in `registries/skills.json`. Sync does not copy them. If a host still has a folder from an older copy, remove it by hand or with `node runtime/tools/ow-sync.mjs --retire <id> --apply` only when `.ow-manifest.json` records that Orchestra installed that hash.

## Global instruction files

`kit/bootstrap.ps1` stamps `AGENTS.md`, `CLAUDE.md`, and `.cursorrules` between `<!-- orchestra:start -->` and `<!-- orchestra:end -->`. An unmarked file is left alone. Back up `~/.claude/CLAUDE.md` and `~/.gemini/GEMINI.md` before a manual replace. This release does not overwrite those two global files by itself.

Limitless stays optional. The model class list is `registries/model-classes.json`. It does not name a vendor model.

## MCP on this release

On 2026-11-08 the Cursor user config and the Antigravity user config were copied to a `.bak-4.5.0` file beside each config, then `shadcn` (`npx shadcn@latest mcp`) and `magic-ui` (`npx -y @magicuidesign/mcp@latest`) were added. The template is `templates/mcp/optional-servers.json`. Keyed servers were left unset. Java was not on `PATH`, so REA was not installed. The global `CLAUDE.md` and `GEMINI.md` were not edited.

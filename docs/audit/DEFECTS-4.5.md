# Defects still live at v4.5.0 (`f4208a3`)

Checked 2026-11-08 on this machine. CI and Hygiene were green on that commit. These are the gaps that green check did not cover.

- `~/.cursor/skills` still has 29 pre-rename folders (`orchestra-*`, `prompt-forge`, `r3f-shaders`, `backend-architecture`, and the rest of the alias list) beside the `ow-*` copies. The same set is on the other hosts. Sync never retires an old name.
- `~/.gemini/GEMINI.md` is still the 4.1.0 adapter and names `visual-forensics`, `r3f-shaders`, and `backend-architecture`.
- `~/.claude/CLAUDE.md` is a full-file copy. It still says the allowlist cap is 80.
- `.cursorrules` and `CLAUDE.md` tell the agent to follow `fullstack-mastery`, `creative-3d-web`, and `r3f-shaders`. `AGENTS.md` forbids those beside the current skills.
- `runtime/internal/adapters/host_sync.go` `SyncAll` copies over a local edit. `ow-sync.mjs` does not. Limitless `scripts/orchestra.js` copies `CLAUDE.md` over the global file and copies every folder under `skills/`, including untracked ones.
- Limitless treats the adapter as synced only when the file contains the text `Orchestra 4.5` or `Orchestra 4.1`.
- The README compile diagram uses a Mermaid node id `class`. Mermaid reserves that word, so the diagram does not render.
- The c1 landing trials ran in bare temp folders. The Orchestra arm had no `AGENTS.md`. Both arms could see the global skills. The record measures a prompt difference.
- `actions/checkout` and `actions/setup-node` are pinned to SHAs that GitHub warns still target Node.js 20.
- `docs/audit/FILE-INVENTORY.json` counts the tree at `60cf9e7`. `runtime/tools/rename-ow.mjs` and `append-public-resources.mjs` are one-shot scripts left in the tree.
- Three host skills stay as conflicts because their bytes differ from this tree and the manifest does not own them: `agent-reach`, `eas-app-stores`, `pretty-mermaid`.

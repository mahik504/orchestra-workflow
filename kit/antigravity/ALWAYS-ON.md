# Antigravity — stay on Orchestra 3.1.0

`AGENTS.md` in the workflow clone (and the vault overlay) is the contract. This file does not invent a second loop. If this note and `AGENTS.md` disagree, `AGENTS.md` wins.

Every new Antigravity chat:

1. Open the **app repo**. Add Folder the private workspace (`ORCHESTRA_HOME`).
2. Fill and paste `MASTER-PROMPT.md` once (`MODE=specialist` if Cursor is conducting, `MODE=conductor` if Antigravity is alone).
3. The `orchestra-conductor` skill should already be in `~/.gemini/config/skills/`. Do not `npx skills add --all`.

You do **not** need `orchestra classify` for Orchestra to be on. That binary scores a brief and enforces the Design Lab write lock. The markdown contract is the default.

Set a persistent `ORCHESTRA_HOME` (see `kit/install-local-engine.ps1` / `.sh`) so `orchestra doctor` binds memory to the Brain, not the clone's `.orchestra/`.

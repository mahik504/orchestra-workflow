# Antigravity orchestra kit

This folder is what you send a friend who has **Antigravity Pro** and no Cursor.

Zip **only this folder**. Never the private vault. Never `mcp_config.json` (secrets). Never `memory/career.md`.

Public template (workflow + this kit, no career/unpublished briefs): [mahik504/orchestra-workflow](https://github.com/mahik504/orchestra-workflow). Private backup stays [<your-private-vault>](https://github.com/example/your-private-vault).

## Honest limit

A prompt cannot log into Google, create API keys, or install Docker by itself. After paste, Antigravity can install skills and write MCP config. **You** still paste a Stitch key (or install Stitch from the MCP Store) and sign in to GitHub in the MCP Store. Then restart Antigravity.

## Kinds

Every product idea has `kind`: `college` | `personal` | `hiring-cv`. Do not mix college assignments into hiring-cv work.

## You (author) — specialist mode

Cursor stays conductor. In Antigravity: open the **app repo** + add folder `<ORCHESTRA_HOME>`. Paste `MASTER-PROMPT.md` with `MODE=specialist`. Read `STACK.md` for what is installed vs refused.

## Friend — conductor mode

1. Copy this `kit/antigravity` folder to them (zip). Do **not** send your whole vault (`career.md`, product ideas, Preferences with your notes, `STACK.md`).
2. They pick a vault path, e.g. `<ORCHESTRA_HOME>`.
3. Antigravity: open their **app folder**, Add Folder the vault (create it if missing).
4. Paste `MASTER-PROMPT.md` with `MODE=conductor` and their vault path.
5. They put **their** Stitch key into MCP config. Never yours.

## Do not install (friend or us)

OpenHands, Dify, Langflow, Coolify, Maxun, OpenCode, Kilo, OmniRoute, addyosmani/agent-skills, vercel/obra/frontend-design dumps. Optional Claude Code is a **separate CLI** think box, not an IDE plugin.

Vault backup for the author is **private git** (12h). The public template is `orchestra-workflow`. Friend vaults: clone the public template or this zip; never reuse the private remote.

## Do not share

Stitch API keys, GitHub PATs, or `mcp_config.json` from this machine.

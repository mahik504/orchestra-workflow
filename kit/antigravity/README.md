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

## Customization budget (required)

Antigravity has a **token budget for Global customizations**. Science and data-engineering plugin packs as Global will consume it before Orchestra can speak.

**Keep Global:** the Orchestra 30 skills, Stitch, Expo, design (taste / impeccable / emil), ship-safe / Strix / semgrep, and tiny Antigravity builtins.

**Disable as Global:**

1. Open Antigravity → **Settings → Customizations** (or the Plugins list).
2. Find **science** and **data-agent-kit** / `data-agent-kit-plugin`.
3. Turn each **off**. Re-enable only for a job that needs those skills.
4. Restart Antigravity.

The same preference lives in `~/.gemini/config/config.json` under `plugins.<name>.enabled`. `false` wins over the plugin default.

**MCP keep as HEALTHY:** stitch, vault-memory, playwright, context7.

**MCP AUTH_REQUIRED until you sign in for a job:** supabase, and any other host-connected cloud. Listing them is not "active."

```
orchestra doctor
# or, without Go:
powershell -NoProfile -ExecutionPolicy Bypass -File kit/orchestra-doctor.ps1
```

Doctor must list the Global skill names and warn if science or data-agent-kit come back as Global.

## Do not share

Stitch API keys, GitHub PATs, or `mcp_config.json` from this machine.
